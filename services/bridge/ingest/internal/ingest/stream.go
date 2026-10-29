package ingest

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/coder/websocket"
)

// ---------------------------------------------------------------- SSE

// SSEEvent is one dispatched server-sent event (WHATWG HTML §9.2.6).
type SSEEvent struct {
	ID    string
	Event string
	Data  string
	Retry time.Duration
}

// ReadSSE parses an event stream and calls fn for each event. Lines are capped at
// maxLine bytes so a hostile stream cannot grow memory without bound.
func ReadSSE(r *bufio.Scanner, fn func(SSEEvent) error) error {
	var ev SSEEvent
	var data []string
	for r.Scan() {
		line := r.Text()
		if line == "" {
			if len(data) > 0 {
				ev.Data = strings.Join(data, "\n")
				if err := fn(ev); err != nil {
					return err
				}
			}
			id := ev.ID // the last event id persists across events
			ev, data = SSEEvent{ID: id}, nil
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		field, value, _ := strings.Cut(line, ":")
		value = strings.TrimPrefix(value, " ")
		switch field {
		case "id":
			if !strings.ContainsRune(value, 0) {
				ev.ID = value
			}
		case "event":
			ev.Event = value
		case "data":
			data = append(data, value)
		case "retry":
			if ms, err := strconv.Atoi(value); err == nil {
				ev.Retry = time.Duration(ms) * time.Millisecond
			}
		}
	}
	return r.Err()
}

// runSSE keeps one stream open, resuming from the checkpointed Last-Event-ID so a
// reconnect or restart leaves no gap (ingestion-platform.md finding 2). The cursor is
// checkpointed after acks, at most once per second, so KV writes stay bounded under load.
func (s *Source) runSSE(ctx context.Context) {
	fails := 0
	for ctx.Err() == nil {
		err := s.sseOnce(ctx)
		if ctx.Err() != nil {
			return
		}
		if err == nil {
			err = errors.New("sse stream ended")
			fails = 0
		} else {
			fails++
		}
		s.fail(err)
		if sleepCtx(ctx, min(time.Duration(1<<min(fails, 6))*time.Second, s.cfg.MaxBackoff.Duration)) != nil {
			return
		}
	}
}

func (s *Source) sseOnce(ctx context.Context) error {
	cp, err := s.rt.LoadCheckpoint(ctx, s.cfg.ID)
	if err != nil {
		return err
	}
	if err := s.fetcher.Quota.Take(); err != nil {
		return err
	}
	if err := s.fetcher.Limiter.Wait(ctx); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.cfg.URL, nil)
	if err != nil {
		return err
	}
	for k, v := range s.cfg.Headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Accept", "text/event-stream")
	if cp.Cursor != "" {
		req.Header.Set("Last-Event-ID", cp.Cursor)
	}
	resp, err := (&http.Client{}).Do(req) // no client timeout: the stream is long-lived
	if err != nil {
		s.fetcher.observe(0)
		return err
	}
	defer resp.Body.Close()
	s.fetcher.observe(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return &StatusError{resp.StatusCode}
	}
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 64<<10), int(min(s.cfg.MaxBodyBytes, 4<<20)))
	lastSave := time.Now()
	pending := false
	defer func() {
		if pending {
			_ = s.rt.SaveCheckpoint(context.WithoutCancel(ctx), s.cfg.ID, cp, s.cfg.SeenMax)
		}
	}()
	return ReadSSE(sc, func(ev SSEEvent) error {
		recs, perr := s.parser([]byte(ev.Data), s.cfg.Options)
		if perr != nil {
			recs = []Record{{Raw: []byte(ev.Data), Err: fmt.Errorf("parse: %v", perr)}}
		}
		n, err := s.rt.Emit(ctx, s, recs, &cp, EmitOptions{})
		if err != nil {
			return err
		}
		if ev.ID != "" {
			cp.Cursor = ev.ID
		}
		pending = true
		s.rt.Board.OK(s.cfg.ID, len(recs), n)
		if time.Since(lastSave) >= time.Second {
			_ = s.rt.SaveCheckpoint(ctx, s.cfg.ID, cp, s.cfg.SeenMax)
			lastSave, pending = time.Now(), false
		}
		return nil
	})
}

// ---------------------------------------------------------------- WebSocket

// runWebSocket reads frames into a bounded queue and publishes from it on another
// goroutine, so the reader keeps answering pings while the bus is slow or down. When the
// queue is full the oldest frame is dropped and counted (ingest_dropped_total): a policy
// for position streams, which have no replay. Memory is bounded by QueueMax frames of at
// most MaxBodyBytes each (capped at 1 MiB).
func (s *Source) runWebSocket(ctx context.Context) {
	q := NewDropQueue[[]byte](s.cfg.QueueMax)
	depth := s.rt.Metrics.QueueDepth.WithLabelValues(s.cfg.ID)
	dropped := s.rt.Metrics.Dropped.WithLabelValues(s.cfg.ID)
	go func() {
		var cp Checkpoint // streams without replay keep no seen set; Nats-Msg-Id dedups
		for {
			raw, ok := q.Pop(ctx)
			if !ok {
				return
			}
			depth.Set(float64(q.Len()))
			recs, perr := s.parser(raw, s.cfg.Options)
			if perr != nil {
				recs = []Record{{Raw: raw, Err: fmt.Errorf("parse: %v", perr)}}
			}
			cp.Seen = cp.Seen[:0]
			n, err := s.rt.Emit(ctx, s, recs, &cp, EmitOptions{})
			if err != nil {
				return
			}
			s.rt.Board.OK(s.cfg.ID, len(recs), n)
		}
	}()
	fails := 0
	for ctx.Err() == nil {
		err := s.wsOnce(ctx, q, func() { dropped.Inc() }, func() { depth.Set(float64(q.Len())) })
		if ctx.Err() != nil {
			return
		}
		if err == nil {
			err = errors.New("websocket closed")
		}
		fails++
		s.fail(err)
		if sleepCtx(ctx, min(time.Duration(1<<min(fails, 6))*time.Second, s.cfg.MaxBackoff.Duration)) != nil {
			return
		}
	}
}

func (s *Source) wsOnce(ctx context.Context, q *DropQueue[[]byte], onDrop, onPush func()) error {
	if err := s.fetcher.Quota.Take(); err != nil {
		return err
	}
	if err := s.fetcher.Limiter.Wait(ctx); err != nil {
		return err
	}
	h := http.Header{}
	for k, v := range s.cfg.Headers {
		h.Set(k, v)
	}
	conn, resp, err := websocket.Dial(ctx, s.cfg.URL, &websocket.DialOptions{HTTPHeader: h})
	if resp != nil {
		s.fetcher.observe(resp.StatusCode)
	}
	if err != nil {
		return err
	}
	defer conn.CloseNow()
	conn.SetReadLimit(min(s.cfg.MaxBodyBytes, 1<<20))
	if s.cfg.Subscribe != "" {
		if err := conn.Write(ctx, websocket.MessageText, []byte(s.cfg.Subscribe)); err != nil {
			return err
		}
	}
	for {
		_, raw, err := conn.Read(ctx)
		if err != nil {
			return err
		}
		if q.Push(raw) {
			onDrop()
		}
		onPush()
	}
}
