package ingest

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Source is one configured feed bound to the runtime.
type Source struct {
	cfg     SourceConfig
	rt      *Runtime
	parser  Parser
	schema  *jsonschema.Schema
	fetcher *Fetcher
	// Once makes a poll source stop after one successful fetch (bulk CronJob mode).
	Once bool
}

func NewSource(rt *Runtime, cfg SourceConfig) (*Source, error) {
	name := cfg.Parser
	if name == "" && cfg.Kind == KindRSS {
		name = "rss"
	}
	p, ok := LookupParser(name)
	if !ok {
		return nil, fmt.Errorf("source %s: unknown parser %q (have %v)", cfg.ID, name, ParserNames())
	}
	sch, err := rt.Schemas.Get(cfg.Schema)
	if err != nil {
		return nil, fmt.Errorf("source %s: %w", cfg.ID, err)
	}
	s := &Source{cfg: cfg, rt: rt, parser: p, schema: sch}
	s.fetcher = NewFetcher(cfg, func(code int) {
		rt.Metrics.Requests.WithLabelValues(cfg.ID, strconv.Itoa(code)).Inc()
	})
	rt.Board.Add(cfg.ID, cfg.Kind, cfg.StaleAfter.Duration)
	if s.fetcher.Quota != nil {
		rt.Metrics.Quota.WithLabelValues(cfg.ID).Set(float64(s.fetcher.Quota.Remaining()))
	}
	return s, nil
}

func (s *Source) Config() SourceConfig { return s.cfg }
func (s *Source) Fetcher() *Fetcher    { return s.fetcher }

// Run runs the adapter for the source's delivery kind until ctx ends. Call it through
// Runtime.RunOwned so exactly one instance runs each source.
func (s *Source) Run(ctx context.Context) {
	s.rt.Log.Info("source start", "source", s.cfg.ID, "kind", s.cfg.Kind, "url", redactURL(s.cfg.URL))
	switch s.cfg.Kind {
	case KindRESTPoll, KindRSS, KindFile:
		s.runPoll(ctx)
	case KindSSE:
		s.runSSE(ctx)
	case KindWebSocket:
		s.runWebSocket(ctx)
	}
}

func (s *Source) fail(err error) {
	s.rt.Metrics.FetchErrors.WithLabelValues(s.cfg.ID, ErrorReason(err)).Inc()
	s.rt.Board.Fail(s.cfg.ID, err)
	s.rt.Log.Warn("source error", "source", s.cfg.ID, "reason", ErrorReason(err), "err", err.Error())
}

// nextDelay is the poll spacing: the interval after a success, and capped exponential
// backoff after consecutive failures (429/5xx/transport), or the quota reset time.
func (s *Source) nextDelay(fails int, err error) time.Duration {
	var qe *QuotaError
	if errors.As(err, &qe) {
		return max(time.Until(qe.ResetAt), s.cfg.Interval.Duration)
	}
	var se *StatusError
	if errors.As(err, &se) && se.RetryAfter > 0 {
		return max(se.RetryAfter, s.cfg.Interval.Duration)
	}
	d := s.cfg.Interval.Duration
	for i := 0; i < fails && d < s.cfg.MaxBackoff.Duration; i++ {
		d *= 2
	}
	return min(d, max(s.cfg.MaxBackoff.Duration, s.cfg.Interval.Duration))
}

// ---------------------------------------------------------------- poll: REST, RSS, file/bulk

// PollOnce does one fetch-parse-publish-checkpoint cycle. It returns (records, new, err).
// The checkpoint (ETag, Last-Modified, seen set) is written only after every record of
// the fetch is acked; a failed checkpoint read skips the fetch entirely.
func (s *Source) PollOnce(ctx context.Context) (int, int, error) {
	cp, err := s.rt.LoadCheckpoint(ctx, s.cfg.ID)
	if err != nil {
		return 0, 0, err
	}
	res, err := s.fetcher.Get(ctx, s.cfg.URL, Conditional{ETag: cp.ETag, LastModified: cp.LastModified})
	if s.fetcher.Quota != nil {
		s.rt.Metrics.Quota.WithLabelValues(s.cfg.ID).Set(float64(s.fetcher.Quota.Remaining()))
	}
	if err != nil {
		return 0, 0, err
	}
	if res.NotModified() {
		return 0, 0, nil
	}
	recs, perr := s.parser(res.Body, s.cfg.Options)
	if perr != nil {
		// The whole body is unusable: dead-letter it once and keep the old checkpoint.
		if rerr := s.rt.Reject(ctx, s.cfg.ID, "parse: "+perr.Error(), res.Body); rerr != nil {
			return 0, 0, rerr
		}
		return 0, 0, ParseError(perr)
	}
	n, err := s.rt.Emit(ctx, s, recs, &cp, EmitOptions{SuppressRepeatRejects: true})
	if err != nil {
		return len(recs), n, err // not checkpointed: the next poll refetches, dedup absorbs it
	}
	cp.ETag, cp.LastModified = res.ETag, res.LastModified
	cp.Seen = keepCurrentLast(cp.Seen, recs)
	if len(recs) > s.cfg.SeenMax {
		s.rt.Log.Warn("snapshot larger than seenMax; raise seenMax or rely on Nats-Msg-Id", "source", s.cfg.ID, "records", len(recs), "seenMax", s.cfg.SeenMax)
	}
	_ = s.rt.SaveCheckpoint(ctx, s.cfg.ID, &cp, s.cfg.SeenMax)
	return len(recs), n, nil
}

// keepCurrentLast moves the keys of the current snapshot to the end of the seen ring, so
// trimming the ring drops only items the upstream no longer lists.
func keepCurrentLast(seen []string, recs []Record) []string {
	cur := make(map[string]struct{}, len(recs))
	var tail []string
	for _, r := range recs {
		if r.ID == "" || r.Err != nil {
			continue
		}
		if _, ok := cur[r.ID]; !ok {
			cur[r.ID] = struct{}{}
			tail = append(tail, r.ID)
		}
	}
	out := make([]string, 0, len(seen))
	for _, k := range seen {
		if _, ok := cur[k]; !ok {
			out = append(out, k)
		}
	}
	have := map[string]struct{}{}
	for _, k := range seen {
		have[k] = struct{}{}
	}
	for _, k := range tail {
		if _, ok := have[k]; ok { // only keys that were actually acked
			out = append(out, k)
		}
	}
	return out
}

func (s *Source) runPoll(ctx context.Context) {
	fails := 0
	for ctx.Err() == nil {
		count, published, err := s.PollOnce(ctx)
		switch {
		case err == nil:
			fails = 0
			s.rt.Board.OK(s.cfg.ID, count, published)
			if published > 0 {
				s.rt.Log.Info("published", "source", s.cfg.ID, "records", count, "new", published)
			}
			if s.Once {
				return
			}
		case ctx.Err() != nil:
			return
		case errors.Is(err, errCheckpointRead):
			s.fail(err) // skip this poll; never refetch as if there were no state
			fails = 0
		default:
			fails++
			s.fail(err)
		}
		if sleepCtx(ctx, s.nextDelay(fails, err)) != nil {
			return
		}
	}
}

func redactURL(u string) string {
	if i := indexAny(u, "?"); i >= 0 {
		return u[:i] + "?…"
	}
	return u
}

func indexAny(s, chars string) int {
	for i := 0; i < len(s); i++ {
		for j := 0; j < len(chars); j++ {
			if s[i] == chars[j] {
				return i
			}
		}
	}
	return -1
}
