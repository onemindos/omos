// Package bus holds the NATS helpers shared by the agentops components: connecting, and
// publishing a fact with a JetStream ack where a stream binds the subject, or a core
// publish where none does yet (an owner-proposed stream binding that has not merged).
package bus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Connect opens a NATS connection that keeps reconnecting forever.
func Connect(url, user, password, name string) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name(name),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
	}
	if user != "" {
		opts = append(opts, nats.UserInfo(user, password))
	}
	return nats.Connect(url, opts...)
}

// Publish modes.
const (
	ModeAuto      = "auto"      // JetStream publish; core publish if no stream binds the subject
	ModeJetStream = "jetstream" // JetStream publish only; an unbound subject is an error
	ModeCore      = "core"      // core publish only
)

// Publisher publishes facts. It is safe for concurrent use.
type Publisher struct {
	NC   *nats.Conn
	JS   jetstream.JetStream
	Mode string

	Stored   atomic.Int64 // acked by a stream
	Core     atomic.Int64 // published without a stream (auto fallback or core mode)
	Failures atomic.Int64

	mu    sync.Mutex
	bound map[string]boundEntry // subject -> does a stream bind it (auto mode)
}

type boundEntry struct {
	ok bool
	at time.Time
}

// streamBinds asks the server (once a minute per subject) whether a stream captures subject.
// A JetStream publish to an unbound subject that has core subscribers would wait for an ack
// that never comes, so auto mode must know before it publishes.
func (p *Publisher) streamBinds(ctx context.Context, subject string) (bool, error) {
	p.mu.Lock()
	if e, ok := p.bound[subject]; ok && time.Since(e.at) < time.Minute {
		p.mu.Unlock()
		return e.ok, nil
	}
	p.mu.Unlock()
	_, err := p.JS.StreamNameBySubject(ctx, subject)
	ok := err == nil
	if err != nil && !errors.Is(err, jetstream.ErrStreamNotFound) {
		return false, err
	}
	p.mu.Lock()
	if p.bound == nil || len(p.bound) > 10000 {
		p.bound = map[string]boundEntry{}
	}
	p.bound[subject] = boundEntry{ok: ok, at: time.Now()}
	p.mu.Unlock()
	return ok, nil
}

// NewPublisher returns a publisher in the given mode (auto when empty).
func NewPublisher(nc *nats.Conn, mode string) (*Publisher, error) {
	if mode == "" {
		mode = ModeAuto
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	return &Publisher{NC: nc, JS: js, Mode: mode}, nil
}

// Publish sends data on subject. msgID becomes Nats-Msg-Id so a retry or replay is
// deduplicated by the stream. It reports whether a stream stored the message.
func (p *Publisher) Publish(ctx context.Context, subject string, data []byte, msgID string, hdr nats.Header) (bool, error) {
	m := nats.NewMsg(subject)
	m.Data = data
	for k, vs := range hdr {
		for _, v := range vs {
			m.Header.Add(k, v)
		}
	}
	if msgID != "" {
		m.Header.Set(nats.MsgIdHdr, msgID)
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	core := p.Mode == ModeCore
	if p.Mode == ModeAuto {
		ok, err := p.streamBinds(ctx, subject)
		if err != nil {
			p.Failures.Add(1)
			return false, err
		}
		core = !ok
	}
	if core {
		if err := p.NC.PublishMsg(m); err != nil {
			p.Failures.Add(1)
			return false, err
		}
		p.Core.Add(1)
		return false, nil
	}
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		_, err := p.JS.PublishMsg(ctx, m)
		if err == nil {
			p.Stored.Add(1)
			return true, nil
		}
		lastErr = err
		if errors.Is(err, jetstream.ErrNoStreamResponse) || errors.Is(err, nats.ErrNoResponders) {
			break // no stream binds the subject (jetstream mode makes that an error)
		}
		select {
		case <-ctx.Done():
			p.Failures.Add(1)
			return false, lastErr
		case <-time.After(time.Duration(200*(attempt+1)) * time.Millisecond):
		}
	}
	p.Failures.Add(1)
	return false, lastErr
}

// Request sends a JSON request and returns the reply bytes.
func Request(ctx context.Context, nc *nats.Conn, subject string, body []byte, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	m, err := nc.RequestWithContext(ctx, subject, body)
	if err != nil {
		return nil, err
	}
	// NATS micro services answer errors with headers and an empty body.
	if code := m.Header.Get("Nats-Service-Error-Code"); code != "" {
		return nil, fmt.Errorf("%s: service error %s: %s", subject, code, m.Header.Get("Nats-Service-Error"))
	}
	return m.Data, nil
}

// Follow delivers every new message (and, with lastPerSubject, the last stored one per
// subject first) from stream to fn until ctx ends. When the stream does not exist it falls
// back to a core subscription on the filters, which sees only new messages.
func Follow(ctx context.Context, nc *nats.Conn, js jetstream.JetStream, stream string, filters []string, lastPerSubject bool, fn func(subject string, data []byte, hdr nats.Header)) (usedStream bool, err error) {
	if stream != "" {
		if s, serr := js.Stream(ctx, stream); serr == nil {
			pol := jetstream.DeliverNewPolicy
			if lastPerSubject {
				pol = jetstream.DeliverLastPerSubjectPolicy
			}
			c, cerr := s.OrderedConsumer(ctx, jetstream.OrderedConsumerConfig{FilterSubjects: filters, DeliverPolicy: pol})
			if cerr == nil {
				cc, err := c.Consume(func(m jetstream.Msg) {
					fn(m.Subject(), m.Data(), m.Headers())
				})
				if err == nil {
					go func() { <-ctx.Done(); cc.Stop() }()
					return true, nil
				}
			}
		}
	}
	var subs []*nats.Subscription
	for _, f := range filters {
		sub, err := nc.Subscribe(f, func(m *nats.Msg) { fn(m.Subject, m.Data, m.Header) })
		if err != nil {
			for _, s := range subs {
				_ = s.Unsubscribe()
			}
			return false, err
		}
		subs = append(subs, sub)
	}
	go func() {
		<-ctx.Done()
		for _, s := range subs {
			_ = s.Unsubscribe()
		}
	}()
	return false, nil
}
