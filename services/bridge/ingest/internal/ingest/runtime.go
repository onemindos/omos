package ingest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func bytesReader(b []byte) io.Reader { return bytes.NewReader(b) }

// Runtime is the shared collector core for one process.
type Runtime struct {
	NC       *nats.Conn
	JS       jetstream.JetStream
	Locks    jetstream.KeyValue
	State    jetstream.KeyValue
	Board    *StatusBoard
	Metrics  *Metrics
	Schemas  *Schemas
	Log      *slog.Logger
	Instance string
	LeaseTTL time.Duration
	// AckTimeout bounds one publish attempt; MaxPublishBackoff caps the retry spacing.
	AckTimeout        time.Duration
	MaxPublishBackoff time.Duration
	Now               func() time.Time

	leaseMu sync.Mutex
	// fenced lists the sources run through RunOwned; held is true while this instance
	// holds the lease. SaveCheckpoint refuses to write a fenced source's checkpoint
	// unless the lease is held here and in the bucket.
	fenced map[string]bool
	held   map[string]bool
}

func (r *Runtime) setHeld(key string, v bool) {
	r.leaseMu.Lock()
	defer r.leaseMu.Unlock()
	if r.fenced == nil {
		r.fenced, r.held = map[string]bool{}, map[string]bool{}
	}
	r.fenced[key] = true
	r.held[key] = v
}

// leaseState reports whether src runs under a lease and whether this instance holds it.
func (r *Runtime) leaseState(src string) (fenced, held bool) {
	r.leaseMu.Lock()
	defer r.leaseMu.Unlock()
	return r.fenced[src], r.held[src]
}

func (r *Runtime) now() time.Time {
	if r.Now != nil {
		return r.Now()
	}
	return time.Now()
}

// BindBuckets looks up the lease and checkpoint buckets. In the cluster they are NACK
// KeyValue CRs (platform/nats/streams/kv-ingest-*.yaml); create=true is for local runs and
// tests only.
func (r *Runtime) BindBuckets(ctx context.Context, cfg *Config, create bool) error {
	var err error
	if create {
		if r.Locks, err = r.JS.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: cfg.LockBucket, History: 1, TTL: cfg.LeaseTTL.Duration, MaxBytes: 8 << 20}); err != nil {
			return err
		}
		r.State, err = r.JS.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: cfg.StateBucket, History: 1, MaxBytes: 64 << 20})
		return err
	}
	if r.Locks, err = r.JS.KeyValue(ctx, cfg.LockBucket); err != nil {
		return fmt.Errorf("bucket %s: %w", cfg.LockBucket, err)
	}
	if r.State, err = r.JS.KeyValue(ctx, cfg.StateBucket); err != nil {
		return fmt.Errorf("bucket %s: %w", cfg.StateBucket, err)
	}
	return nil
}

// ---------------------------------------------------------------- publish with ack

// ErrPermanentPublish wraps a publish error that no retry can fix: the message is larger
// than the server's max_payload or the stream's max_msg_size, or the user may not publish
// on the subject. Emit dead-letters such a record instead of blocking its feed.
var ErrPermanentPublish = errors.New("publish cannot succeed")

// permanentPublishError classifies err. A publish permission violation reaches the client
// asynchronously (the JetStream publish itself just times out), so it is read from the
// connection's last error for this subject.
// It returns the error to report, or nil when a retry may succeed.
func (r *Runtime) permanentPublishError(subject string, err error) error {
	if errors.Is(err, nats.ErrMaxPayload) || errors.Is(err, jetstream.ErrStreamMessageExceedsMaximum) || errors.Is(err, nats.ErrPermissionViolation) {
		return err
	}
	var je jetstream.JetStreamError
	if errors.As(err, &je) && je.APIError() != nil && je.APIError().ErrorCode == jetstream.JSErrCodeStreamMessageExceedsMaximum {
		return err
	}
	if r.NC != nil {
		if le := r.NC.LastError(); le != nil && errors.Is(le, nats.ErrPermissionViolation) &&
			strings.Contains(strings.ToLower(le.Error()), strings.ToLower(`publish to "`+subject+`"`)) {
			return le
		}
	}
	return nil
}

// Publish sends one message with Nats-Msg-Id and waits for the JetStream PubAck. It
// retries with capped exponential backoff until it is acked or ctx ends, so a record is
// never dropped: a NATS outage, a missing stream or a full stream (discard new) all block
// here, and the caller does not checkpoint. Errors no retry can fix return at once,
// wrapped in ErrPermanentPublish. It reports whether the ack was a duplicate.
func (r *Runtime) Publish(ctx context.Context, src, subject, msgID string, data []byte) (bool, error) {
	backoff := 100 * time.Millisecond
	maxB := r.MaxPublishBackoff
	if maxB == 0 {
		maxB = 10 * time.Second
	}
	ackTO := r.AckTimeout
	if ackTO == 0 {
		ackTO = 5 * time.Second
	}
	for {
		msg := nats.NewMsg(subject)
		msg.Header.Set(jetstream.MsgIDHeader, msgID)
		msg.Data = data
		pctx, cancel := context.WithTimeout(ctx, ackTO)
		ack, err := r.JS.PublishMsg(pctx, msg)
		cancel()
		if err == nil {
			if ack.Duplicate {
				r.Metrics.Duplicates.WithLabelValues(src).Inc()
			} else {
				r.Metrics.Published.WithLabelValues(src).Inc()
			}
			return ack.Duplicate, nil
		}
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
		if perr := r.permanentPublishError(subject, err); perr != nil {
			return false, fmt.Errorf("%w: %v", ErrPermanentPublish, perr)
		}
		r.Metrics.Retries.WithLabelValues(src).Inc()
		r.Log.Warn("publish retry", "source", src, "subject", subject, "msgId", msgID, "err", err.Error(), "backoff", backoff.String())
		if err := sleepCtx(ctx, backoff); err != nil {
			return false, err
		}
		backoff = min(backoff*2, maxB)
	}
}

// Reject publishes a dead-letter record on evt.sys.{source}.rejected with the reason.
// The Nats-Msg-Id is unique per occurrence (hash plus time), so identical bad frames on
// a stream are all recorded (ingestion-platform.md finding 8).
func (r *Runtime) Reject(ctx context.Context, src, reason string, raw []byte) error {
	r.Metrics.Rejected.WithLabelValues(src).Inc()
	rs := string(raw)
	if len(rs) > 4096 {
		rs = rs[:4096]
	}
	now := r.now()
	id := fmt.Sprintf("rej:%s:%016x:%d", src, fnv64(raw), now.UnixNano())
	subj := RejectSubject(src)
	env := NewEnvelope(src, subj, id, now, now, map[string]any{"source": src, "reason": reason, "raw": rs, "rawBytes": len(raw)})
	b, _ := json.Marshal(env)
	r.Log.Warn("record rejected", "source", src, "reason", reason)
	_, err := r.Publish(ctx, src, subj, id, b)
	return err
}

func fnv64(b []byte) uint64 {
	h := uint64(14695981039346656037)
	for _, x := range b {
		h ^= uint64(x)
		h *= 1099511628211
	}
	return h
}

// ---------------------------------------------------------------- checkpoints (after ack)

// Checkpoint is what a source needs to resume without loss or replay.
type Checkpoint struct {
	ETag         string   `json:"etag,omitempty"`
	LastModified string   `json:"lastModified,omitempty"`
	Cursor       string   `json:"cursor,omitempty"`
	Seen         []string `json:"seen,omitempty"` // bounded ring of recent dedup keys
	UpdatedAt    string   `json:"updatedAt,omitempty"`

	// Rev is the KV revision this checkpoint was read or last written at (0: none).
	// SaveCheckpoint writes only on top of it, so an instance that lost its lease cannot
	// overwrite a newer checkpoint written by the new owner.
	Rev uint64 `json:"-"`
	// index mirrors Seen for lookups; it lives as long as the Checkpoint, so a stream
	// does not rebuild it on every event.
	index map[string]struct{}
}

func (cp *Checkpoint) seenIndex() map[string]struct{} {
	if cp.index == nil {
		cp.index = make(map[string]struct{}, len(cp.Seen))
		for _, k := range cp.Seen {
			cp.index[k] = struct{}{}
		}
	}
	return cp.index
}

func (cp *Checkpoint) addSeen(k string) {
	cp.seenIndex()[k] = struct{}{}
	cp.Seen = append(cp.Seen, k)
}

// trimSeen drops the oldest keys so at most n remain.
func (cp *Checkpoint) trimSeen(n int) {
	if len(cp.Seen) <= n {
		return
	}
	drop := cp.Seen[:len(cp.Seen)-n]
	if cp.index != nil {
		for _, k := range drop {
			delete(cp.index, k)
		}
	}
	cp.Seen = append([]string(nil), cp.Seen[len(cp.Seen)-n:]...)
}

// ErrCheckpointFenced means a checkpoint was not written because this instance no longer
// holds the source's lease, or another instance wrote the checkpoint since it was read.
var ErrCheckpointFenced = errors.New("checkpoint not written: lease or checkpoint changed")

// LoadCheckpoint never treats "could not read" as "no state": that would republish
// everything the upstream still lists (ingestion-platform.md finding 4). Only a missing
// key is a fresh start.
func (r *Runtime) LoadCheckpoint(ctx context.Context, src string) (Checkpoint, error) {
	var cp Checkpoint
	e, err := r.State.Get(ctx, src)
	if errors.Is(err, jetstream.ErrKeyNotFound) {
		return cp, nil
	}
	if err != nil {
		return cp, fmt.Errorf("%w: %v", errCheckpointRead, err)
	}
	if err := json.Unmarshal(e.Value(), &cp); err != nil {
		return cp, fmt.Errorf("%w: %v", errCheckpointRead, err)
	}
	cp.Rev = e.Revision()
	return cp, nil
}

// SaveCheckpoint is called only after every record of a fetch is acked. A failed save is
// safe: the next fetch republishes and Nats-Msg-Id or the seen set absorbs it.
// It trims cp.Seen in place to seenMax, so long-lived streams keep bounded memory.
//
// The write is fenced twice: a source run through RunOwned must still hold its lease
// (here and in the lock bucket), and the KV write is a compare-and-set on cp.Rev. A
// fenced write returns ErrCheckpointFenced; the caller reloads the checkpoint.
func (r *Runtime) SaveCheckpoint(ctx context.Context, src string, cp *Checkpoint, seenMax int) error {
	if seenMax > 0 {
		cp.trimSeen(seenMax)
	}
	cp.UpdatedAt = r.now().UTC().Format(time.RFC3339)
	b, _ := json.Marshal(*cp)
	for len(b) > 512<<10 && len(cp.Seen) > 0 { // stay well under the 1 MiB max payload
		cp.trimSeen(len(cp.Seen) / 2)
		b, _ = json.Marshal(*cp)
	}
	if fenced, held := r.leaseState(src); fenced {
		ok := held
		if ok {
			e, err := r.Locks.Get(ctx, src)
			ok = err == nil && string(e.Value()) == r.Instance
		}
		if !ok {
			r.Metrics.Checkpoints.WithLabelValues(src, "fenced").Inc()
			r.Log.Warn("checkpoint not written: lease not held", "source", src)
			return ErrCheckpointFenced
		}
	}
	var rev uint64
	var err error
	if cp.Rev == 0 {
		rev, err = r.State.Create(ctx, src, b)
	} else {
		rev, err = r.State.Update(ctx, src, b, cp.Rev)
	}
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyExists) {
			r.Metrics.Checkpoints.WithLabelValues(src, "fenced").Inc()
			r.Log.Warn("checkpoint not written: another instance wrote it since it was read", "source", src)
			return fmt.Errorf("%w: %v", ErrCheckpointFenced, err)
		}
		r.Metrics.Checkpoints.WithLabelValues(src, "error").Inc()
		r.Log.Warn("checkpoint failed; next fetch refetches and dedup absorbs it", "source", src, "err", err.Error())
		return err
	}
	cp.Rev = rev
	r.Metrics.Checkpoints.WithLabelValues(src, "ok").Inc()
	return nil
}

// ---------------------------------------------------------------- emit

// EmitOptions control dedup of rejects for snapshot sources.
type EmitOptions struct {
	// SuppressRepeatRejects records a rejected item's hash in the seen set, so a bad item
	// that a snapshot upstream keeps listing is dead-lettered once, not on every poll.
	SuppressRepeatRejects bool
}

// Emit validates, dead-letters, skips seen keys, wraps in the envelope and publishes each
// record with an ack. It appends acked keys to cp.Seen. It returns the number of new
// records stored. On error the caller must not checkpoint. A record whose publish can
// never succeed (ErrPermanentPublish) is dead-lettered with the reason and marked seen.
func (r *Runtime) Emit(ctx context.Context, s *Source, recs []Record, cp *Checkpoint, opt EmitOptions) (int, error) {
	seen := cp.seenIndex()
	fetchedAt := r.now()
	n := 0
	for _, rec := range recs {
		reason := ""
		switch {
		case rec.Err != nil:
			reason = rec.Err.Error()
		case rec.ID == "":
			reason = "record has no id"
		default:
			if err := Validate(s.schema, rec.Data); err != nil {
				reason = "schema: " + oneLine(err.Error())
			}
		}
		if reason != "" {
			raw := rec.Raw
			if raw == nil {
				raw, _ = json.Marshal(rec.Data)
			}
			rk := fmt.Sprintf("rej:%016x", fnv64(raw))
			if opt.SuppressRepeatRejects {
				if _, dup := seen[rk]; dup {
					continue
				}
			}
			if err := r.Reject(ctx, s.cfg.ID, reason, raw); err != nil {
				return n, err
			}
			if opt.SuppressRepeatRejects {
				cp.addSeen(rk)
			}
			continue
		}
		if _, dup := seen[rec.ID]; dup {
			continue
		}
		subj := rec.Subject
		if subj == "" {
			subj = s.cfg.Subject
		}
		if err := CheckDataSubject(subj); err != nil {
			if err := r.Reject(ctx, s.cfg.ID, err.Error(), rec.Raw); err != nil {
				return n, err
			}
			continue
		}
		b, err := json.Marshal(NewEnvelope(s.cfg.ID, subj, rec.ID, rec.TS, fetchedAt, rec.Data))
		if err != nil {
			return n, err
		}
		dup, err := r.Publish(ctx, s.cfg.ID, subj, rec.ID, b)
		if errors.Is(err, ErrPermanentPublish) {
			if err := r.Reject(ctx, s.cfg.ID, oneLine(err.Error()), b); err != nil {
				return n, err
			}
			cp.addSeen(rec.ID)
			continue
		}
		if err != nil {
			return n, err
		}
		if !dup {
			n++
		}
		cp.addSeen(rec.ID)
	}
	return n, nil
}

func oneLine(s string) string {
	s = strings.ReplaceAll(s, "\n", "; ")
	if len(s) > 500 {
		s = s[:500]
	}
	return s
}

// ---------------------------------------------------------------- single owner (KV lease)

// RunOwned runs fn only while this instance holds the lease for key. The lease is a KV
// entry created with Create (fails if held) and renewed with Update(revision) every
// LeaseTTL/3; the bucket TTL expires it if the holder dies.
//
// A renewal that fails while the bus is unreachable is retried, but only until the lease
// could have expired in the bucket: fn is cancelled once safeWindow (4/5 of LeaseTTL) has
// passed since the last renewal that succeeded was sent. A leader cut off alone therefore
// stops before another instance can take the lease, and after reconnecting it competes
// for the lease again with Create. The NATS client must not buffer publishes while
// disconnected (nats.ReconnectBufSize(-1)), or a stopped leader would replay them.
func (r *Runtime) RunOwned(ctx context.Context, key string, fn func(context.Context)) {
	renewEvery := max(r.LeaseTTL/3, 100*time.Millisecond)
	safeWindow := r.LeaseTTL * 4 / 5
	r.setHeld(key, false)
	for ctx.Err() == nil {
		sent := time.Now()
		rev, err := r.Locks.Create(ctx, key, []byte(r.Instance))
		if err != nil {
			gctx, gcancel := context.WithTimeout(ctx, renewEvery)
			e, gerr := r.Locks.Get(gctx, key)
			gcancel()
			if gerr == nil && string(e.Value()) == r.Instance {
				// Our own entry: only this instance writes its name. Renew it before use so
				// the safe window starts now.
				sent = time.Now()
				uctx, ucancel := context.WithTimeout(ctx, renewEvery)
				rev, err = r.Locks.Update(uctx, key, []byte(r.Instance), e.Revision())
				ucancel()
			}
			if err != nil {
				if sleepCtx(ctx, renewEvery) != nil {
					return
				}
				continue
			}
		}
		r.Log.Info("lease acquired", "source", key)
		r.setHeld(key, true)
		r.Board.SetLeader(key, true)
		lctx, cancel := context.WithCancel(ctx)
		var expired atomic.Bool
		// The fence runs on its own timer, so a renewal call that blocks cannot delay it.
		fence := time.AfterFunc(safeWindow-time.Since(sent), func() {
			expired.Store(true)
			r.setHeld(key, false)
			cancel()
		})
		extend := func(sentAt time.Time) {
			if !expired.Load() {
				fence.Reset(safeWindow - time.Since(sentAt))
			}
		}
		done := make(chan struct{})
		go func() { defer close(done); fn(lctx) }()
		t := time.NewTicker(renewEvery)
	renew:
		for {
			select {
			case <-t.C:
				if expired.Load() {
					break renew
				}
				// Bound each renewal so a dead bus cannot stall the loop past the TTL.
				sentAt := time.Now()
				rctx, rcancel := context.WithTimeout(lctx, renewEvery)
				nrev, err := r.Locks.Update(rctx, key, []byte(r.Instance), rev)
				rcancel()
				if err == nil {
					rev = nrev
					extend(sentAt)
					continue
				}
				if r.NC.Status() != nats.CONNECTED || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, nats.ErrTimeout) || errors.Is(err, context.Canceled) {
					r.Log.Warn("lease renew failed (bus unreachable); source stops if it does not recover within the lease TTL", "source", key, "err", err.Error())
					continue
				}
				gctx, gcancel := context.WithTimeout(lctx, renewEvery)
				e, gerr := r.Locks.Get(gctx, key)
				gcancel()
				if gerr == nil && string(e.Value()) == r.Instance {
					rev = e.Revision() // our entry at a newer revision; renew on the next tick
					continue
				}
				if gerr != nil && !errors.Is(gerr, jetstream.ErrKeyNotFound) {
					continue // transient read error; the fence still bounds how long fn runs
				}
				r.Log.Warn("lease lost", "source", key, "err", err.Error())
				break renew
			case <-done:
				break renew
			case <-ctx.Done():
				break renew
			}
		}
		t.Stop()
		fence.Stop()
		r.setHeld(key, false)
		cancel()
		<-done
		r.Board.SetLeader(key, false)
		if expired.Load() {
			r.Log.Warn("lease not renewed within its TTL; source stopped, competing for the lease again", "source", key)
		}
		if ctx.Err() != nil {
			// Clean shutdown: hand the lease over now instead of after the TTL.
			dctx, dcancel := context.WithTimeout(context.Background(), 2*time.Second)
			_ = r.Locks.Delete(dctx, key, jetstream.LastRevision(rev))
			dcancel()
			return
		}
		if sleepCtx(ctx, renewEvery) != nil {
			return
		}
	}
}

// ---------------------------------------------------------------- freshness heartbeat

// Heartbeat publishes tel.feed.{id}.status for every source this instance leads. It is a
// core publish: a missed heartbeat is not data loss, and TEL stores it (tel.feed.>).
func (r *Runtime) Heartbeat(ctx context.Context, every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-t.C:
		case <-ctx.Done():
			return
		}
		r.PublishStatus()
	}
}

func (r *Runtime) PublishStatus() {
	for _, s := range r.Board.Snapshot() {
		if !s.Leader {
			continue
		}
		b, _ := json.Marshal(s)
		_ = r.NC.Publish(StatusSubject(s.Source), b)
	}
}

// IsCheckpointReadError reports whether a poll was skipped because the checkpoint could
// not be read (as opposed to a missing checkpoint, which is a fresh start).
func IsCheckpointReadError(err error) bool { return errors.Is(err, errCheckpointRead) }

// IsParseError reports a whole-body parse failure.
func IsParseError(err error) bool { return errors.Is(err, errParse) }
