// Package agentproto implements the parts of the Synadia Agent Protocol 0.3.0 (Draft,
// github.com/synadia-ai/synadia-agent-sdk-docs core-protocol.md) that agentops needs:
//
//   - Prompt: a caller of `agents.prompt.{agent}.{owner}.{name}` (§5 request, §6 response
//     chunks, empty-payload terminator, §6.6 inactivity timeout, §9 error headers);
//   - Beater: the §8 heartbeat on `agents.hb.{agent}.{owner}.{name}` and the §8.7 status
//     endpoint, with an optional `onemind` block (unknown fields are allowed by §8.3) that
//     carries the agent's map state: callsign, activity, current task, position;
//   - Serve: a minimal prompt responder (one ack chunk, one response chunk, terminator),
//     used by tests and by local stand-ins for an agent.
//
// Mid-stream queries (§7) are never answered: a query chunk is a permission prompt, and an
// approval given over the bus is not a verified human approval (onemind-xr
// docs/research/agents-that-act.md §8). Prompt returns ErrQueryRefused instead.
package agentproto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// ProtocolVersion is the spec version this package follows.
const ProtocolVersion = "0.3"

// Errors returned by Prompt.
var (
	ErrNoAgent      = errors.New("agent not reachable (no responders)")
	ErrInactivity   = errors.New("agent stream inactivity timeout")
	ErrQueryRefused = errors.New("agent asked a mid-stream query; agentops never answers queries")
)

// AgentError is a §9 error response.
type AgentError struct {
	Code    string
	Message string
}

func (e *AgentError) Error() string { return "agent error " + e.Code + ": " + e.Message }

// Subject builds `agents.{verb}.{agent}.{owner}.{name}`.
func Subject(verb, agent, owner, name string) string {
	return "agents." + verb + "." + agent + "." + owner + "." + name
}

// Prompt sends envelope (it must marshal to an object with a non-empty "prompt" string) and
// returns the concatenated response text. inactivity is the §6.6 per-chunk timeout; ctx
// bounds the whole call.
func Prompt(ctx context.Context, nc *nats.Conn, subject string, envelope any, inactivity time.Duration) (string, error) {
	body, err := json.Marshal(envelope)
	if err != nil {
		return "", err
	}
	if inactivity <= 0 {
		inactivity = 60 * time.Second
	}
	inbox := nc.NewRespInbox()
	sub, err := nc.SubscribeSync(inbox)
	if err != nil {
		return "", err
	}
	defer func() { _ = sub.Unsubscribe() }()
	if err := nc.PublishRequest(subject, inbox, body); err != nil {
		return "", err
	}
	var text strings.Builder
	for {
		wait := inactivity
		if dl, ok := ctx.Deadline(); ok {
			if rem := time.Until(dl); rem < wait {
				wait = rem
			}
		}
		if wait <= 0 {
			return text.String(), ErrInactivity
		}
		m, err := sub.NextMsg(wait)
		if err != nil {
			if errors.Is(err, nats.ErrTimeout) {
				return text.String(), ErrInactivity
			}
			if errors.Is(err, nats.ErrNoResponders) {
				return "", ErrNoAgent
			}
			return text.String(), err
		}
		// A no-responders status arrives as an empty message with a 503 status header.
		if len(m.Data) == 0 && m.Header.Get("Status") == "503" {
			return "", ErrNoAgent
		}
		if code := m.Header.Get("Nats-Service-Error-Code"); code != "" {
			msg := m.Header.Get("Nats-Service-Error")
			var eb struct {
				Message string `json:"message"`
			}
			if json.Unmarshal(m.Data, &eb) == nil && eb.Message != "" {
				msg = eb.Message
			}
			return text.String(), &AgentError{Code: code, Message: msg}
		}
		if len(m.Data) == 0 && len(m.Header) == 0 {
			return text.String(), nil // §6.5 terminator
		}
		var ch struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}
		if json.Unmarshal(m.Data, &ch) != nil {
			continue // §6.2: every chunk must be JSON; ignore what is not
		}
		switch ch.Type {
		case "response":
			var s string
			if json.Unmarshal(ch.Data, &s) == nil {
				text.WriteString(s)
				continue
			}
			var o struct {
				Text string `json:"text"`
			}
			if json.Unmarshal(ch.Data, &o) == nil {
				text.WriteString(o.Text)
			}
		case "query":
			return text.String(), ErrQueryRefused
		}
		// status chunks and unknown types: ignored (§6.4, §6.6)
	}
}

// Task is the agent's current task as shown on the map.
type Task struct {
	TaskID  string `json:"task_id"`
	State   string `json:"state"`
	Summary string `json:"summary,omitempty"`
}

// OneMind is the optional `onemind` block of a heartbeat: what the agent wants shown as its
// entity. Every field is optional.
type OneMind struct {
	Callsign    string   `json:"callsign,omitempty"`
	Mode        string   `json:"mode,omitempty"`     // idle | busy | degraded | offline
	Activity    string   `json:"activity,omitempty"` // free text, at most 200 characters
	Role        string   `json:"role,omitempty"`
	CurrentTask *Task    `json:"current_task,omitempty"`
	Lat         *float64 `json:"lat,omitempty"`
	Lon         *float64 `json:"lon,omitempty"`
}

// Heartbeat is the §8.3 payload plus the OneMind block.
type Heartbeat struct {
	Agent           string         `json:"agent"`
	Owner           string         `json:"owner"`
	Session         string         `json:"session,omitempty"`
	InstanceID      string         `json:"instance_id"`
	TS              string         `json:"ts"`
	IntervalS       float64        `json:"interval_s"`
	ProtocolVersion string         `json:"protocol_version,omitempty"`
	Name            string         `json:"name,omitempty"`
	OneMind         *OneMind       `json:"onemind,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

// Beater publishes the heartbeat of one in-process agent and answers its status endpoint.
type Beater struct {
	nc       *nats.Conn
	agent    string
	owner    string
	name     string
	interval time.Duration
	display  string

	mu    sync.Mutex
	state OneMind
	sub   *nats.Subscription
	onSet func(OneMind)
}

// NewBeater prepares a heartbeat for agent/owner/name.
func NewBeater(nc *nats.Conn, agent, owner, name, display string, interval time.Duration) *Beater {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Beater{nc: nc, agent: agent, owner: owner, name: name, display: display, interval: interval}
}

// Ref returns "agent/owner/name".
func (b *Beater) Ref() string { return b.agent + "/" + b.owner + "/" + b.name }

// OnSet registers a callback run after every Set (the presence tracker uses it so an
// in-process agent's change reaches the map without waiting for the next beat).
func (b *Beater) OnSet(fn func(OneMind)) {
	b.mu.Lock()
	b.onSet = fn
	b.mu.Unlock()
}

// Set replaces the agent's OneMind state and publishes a heartbeat at once.
func (b *Beater) Set(s OneMind) {
	if len(s.Activity) > 200 {
		s.Activity = s.Activity[:200]
	}
	b.mu.Lock()
	b.state = s
	fn := b.onSet
	b.mu.Unlock()
	b.beat()
	if fn != nil {
		fn(s)
	}
}

// State returns the current OneMind state.
func (b *Beater) State() OneMind {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

func (b *Beater) payload() []byte {
	b.mu.Lock()
	st := b.state
	b.mu.Unlock()
	hb := Heartbeat{
		Agent: b.agent, Owner: b.owner, Session: b.name,
		InstanceID: fmt.Sprintf("%s-%s-%s", b.agent, b.owner, b.name),
		TS:         time.Now().UTC().Format(time.RFC3339), IntervalS: b.interval.Seconds(),
		ProtocolVersion: ProtocolVersion, Name: b.display, OneMind: &st,
	}
	out, _ := json.Marshal(hb)
	return out
}

func (b *Beater) beat() {
	_ = b.nc.Publish(Subject("hb", b.agent, b.owner, b.name), b.payload())
}

// Run publishes the heartbeat every interval and serves the status endpoint until ctx ends.
func (b *Beater) Run(ctx context.Context) error {
	sub, err := b.nc.QueueSubscribe(Subject("status", b.agent, b.owner, b.name), "agents", func(m *nats.Msg) {
		_ = m.Respond(b.payload())
	})
	if err != nil {
		return err
	}
	b.mu.Lock()
	b.sub = sub
	b.mu.Unlock()
	b.beat()
	t := time.NewTicker(b.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			_ = sub.Unsubscribe()
			return nil
		case <-t.C:
			b.beat()
		}
	}
}

// Handler answers one prompt. sender is the request's prompt text and the decoded envelope.
type Handler func(ctx context.Context, prompt string, envelope map[string]any) (string, error)

// Serve answers prompts on subject (queue group "agents", §3.3) with one ack, one response
// chunk and the terminator. It returns the subscription.
func Serve(nc *nats.Conn, subject string, h Handler) (*nats.Subscription, error) {
	return nc.QueueSubscribe(subject, "agents", func(m *nats.Msg) {
		if m.Reply == "" {
			return
		}
		ack, _ := json.Marshal(map[string]string{"type": "status", "data": "ack"})
		_ = nc.Publish(m.Reply, ack)
		env := map[string]any{}
		prompt := ""
		trimmed := strings.TrimLeft(string(m.Data), " \t\r\n")
		if strings.HasPrefix(trimmed, "{") {
			if json.Unmarshal(m.Data, &env) != nil {
				respondErr(nc, m.Reply, "400", "invalid envelope")
				return
			}
			prompt, _ = env["prompt"].(string)
		} else {
			prompt = string(m.Data)
		}
		if prompt == "" {
			respondErr(nc, m.Reply, "400", "empty prompt")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		text, err := h(ctx, prompt, env)
		if err != nil {
			respondErr(nc, m.Reply, "500", err.Error())
			return
		}
		chunk, _ := json.Marshal(map[string]string{"type": "response", "data": text})
		_ = nc.Publish(m.Reply, chunk)
		_ = nc.Publish(m.Reply, nil)
	})
}

func respondErr(nc *nats.Conn, reply, code, msg string) {
	em := nats.NewMsg(reply)
	em.Header.Set("Nats-Service-Error-Code", code)
	em.Header.Set("Nats-Service-Error", msg)
	_ = nc.PublishMsg(em)
	_ = nc.Publish(reply, nil)
}
