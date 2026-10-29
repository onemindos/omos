// Package presence turns agents into map entities (FEATURES I1). It follows every agent
// heartbeat on `agents.hb.{agent}.{owner}.{name}` (Synadia Agent Protocol §8, stored by the
// `AGENTS` stream) plus the in-process agents of this binary, and publishes each live agent
// as a doctrine §9.2 geo-envelope on `tel.fleet.agent.{id}.position` (doctrine §9.1 class
// `agent`, type `fleet.agent` from §9.3, stored by `TEL`). The entity service (omos PR #552,
// adapter `Fleet`) turns that into `ent.agent.{entity_id}.state`. Nothing here writes `ent.*`:
// the entity service stays the only writer (onemind-xr docs/contracts/entity.md §2).
//
// An agent has no real position. It is shown at an anchor: the position it reports in its
// heartbeat's `onemind` block, else its configured anchor, else the site anchor. The record
// says so in `meta.position_kind` and carries a large `ce` so nobody mistakes it for a fix.
package presence

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/onemindos/omos/services/agentops/internal/agentproto"
	"github.com/onemindos/omos/services/agentops/internal/bus"
)

// Anchor is a fixed display position.
type Anchor struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Agent is one tracked agent.
type Agent struct {
	UID        string // geo-envelope uid and subject token: agent-{agent}-{owner}-{name}
	Ref        string // agent/owner/name
	Callsign   string
	InstanceID string
	Interval   time.Duration
	LastSeen   time.Time
	State      agentproto.OneMind
	Prompt     string // the prompt subject of agent/owner/name
}

// Online reports whether the agent beat within 3 intervals (protocol §8.2).
func (a *Agent) Online(now time.Time) bool {
	iv := a.Interval
	if iv <= 0 {
		iv = 30 * time.Second
	}
	return now.Sub(a.LastSeen) <= 3*iv
}

var badToken = regexp.MustCompile(`[^a-z0-9_-]+`)

// UIDFor returns the uid of agent/owner/name: lowercase, NATS-token safe, at most 64 chars.
func UIDFor(agent, owner, name string) string {
	s := "agent-" + strings.Trim(badToken.ReplaceAllString(strings.ToLower(agent+"-"+owner+"-"+name), "-"), "-")
	if len(s) > 64 {
		s = s[:64]
	}
	return strings.TrimRight(s, "-")
}

// Config configures the tracker.
type Config struct {
	Site         Anchor            // default anchor
	Anchors      map[string]Anchor // per agent ref (agent/owner/name)
	Callsigns    map[string]string // per agent ref
	PublishEvery time.Duration     // default 30 s
	ForgetAfter  time.Duration     // stop publishing an offline agent after this; default 10 min
	NodeID       string            // meta.node_id
}

// Tracker holds the agents. Safe for concurrent use.
type Tracker struct {
	cfg Config
	now func() time.Time

	mu     sync.Mutex
	agents map[string]*Agent // by uid
	last   map[string]string // uid -> last published state digest
	lastT  map[string]time.Time
}

// New returns a tracker.
func New(cfg Config, now func() time.Time) *Tracker {
	if now == nil {
		now = time.Now
	}
	if cfg.PublishEvery <= 0 {
		cfg.PublishEvery = 30 * time.Second
	}
	if cfg.ForgetAfter <= 0 {
		cfg.ForgetAfter = 10 * time.Minute
	}
	return &Tracker{cfg: cfg, now: now, agents: map[string]*Agent{}, last: map[string]string{}, lastT: map[string]time.Time{}}
}

// OnHeartbeat records one `agents.hb.{agent}.{owner}.{name}` message.
func (t *Tracker) OnHeartbeat(subject string, data []byte) (*Agent, error) {
	tok := strings.Split(subject, ".")
	if len(tok) != 5 || tok[0] != "agents" || tok[1] != "hb" {
		return nil, fmt.Errorf("not a heartbeat subject: %s", subject)
	}
	var hb agentproto.Heartbeat
	if err := json.Unmarshal(data, &hb); err != nil {
		return nil, err
	}
	// §8.3: agent and owner must match the subject; the name is the 5th token.
	if hb.Agent != tok[2] || hb.Owner != tok[3] {
		return nil, fmt.Errorf("heartbeat body %s/%s does not match subject %s", hb.Agent, hb.Owner, subject)
	}
	ref := tok[2] + "/" + tok[3] + "/" + tok[4]
	a := &Agent{UID: UIDFor(tok[2], tok[3], tok[4]), Ref: ref, InstanceID: hb.InstanceID,
		Interval: time.Duration(hb.IntervalS * float64(time.Second)), LastSeen: t.now()}
	if a.Interval <= 0 {
		a.Interval = 30 * time.Second
	}
	if hb.OneMind != nil {
		a.State = *hb.OneMind
	}
	a.Callsign = t.callsign(ref, hb, a.State)
	a.Prompt = agentproto.Subject("prompt", tok[2], tok[3], tok[4])
	t.mu.Lock()
	t.agents[a.UID] = a
	t.mu.Unlock()
	return a, nil
}

func (t *Tracker) callsign(ref string, hb agentproto.Heartbeat, st agentproto.OneMind) string {
	if c := t.cfg.Callsigns[ref]; c != "" {
		return c
	}
	if st.Callsign != "" {
		return st.Callsign
	}
	if hb.Name != "" {
		return strings.ToUpper(hb.Name)
	}
	return strings.ToUpper(hb.Agent)
}

// Get returns a copy of one agent.
func (t *Tracker) Get(uid string) (Agent, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	a, ok := t.agents[uid]
	if !ok {
		return Agent{}, false
	}
	return *a, true
}

// List returns copies of every agent, sorted by uid.
func (t *Tracker) List() []Agent {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Agent, 0, len(t.agents))
	for _, a := range t.agents {
		out = append(out, *a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UID < out[j].UID })
	return out
}

// AnchorOf returns where an agent is drawn and how that position was chosen.
func (t *Tracker) AnchorOf(a Agent) (Anchor, string) {
	if a.State.Lat != nil && a.State.Lon != nil {
		return Anchor{Lat: *a.State.Lat, Lon: *a.State.Lon}, "reported"
	}
	if an, ok := t.cfg.Anchors[a.Ref]; ok {
		return an, "anchor"
	}
	return t.cfg.Site, "site"
}

// Envelope is the doctrine §9.2 geo-envelope for one agent.
type Envelope struct {
	UID      string         `json:"uid"`
	Type     string         `json:"type"`
	Lat      float64        `json:"lat"`
	Lon      float64        `json:"lon"`
	HAE      float64        `json:"hae"`
	Course   float64        `json:"course"`
	Speed    float64        `json:"speed"`
	Callsign string         `json:"callsign"`
	Stale    int            `json:"stale"`
	T        int64          `json:"t"`
	CE       float64        `json:"ce"`
	LE       float64        `json:"le"`
	Meta     map[string]any `json:"meta"`
}

// EnvelopeFor builds the geo-envelope of one agent at now.
func (t *Tracker) EnvelopeFor(a Agent, now time.Time) Envelope {
	an, kind := t.AnchorOf(a)
	online := a.Online(now)
	mode := a.State.Mode
	if mode == "" {
		mode = "idle"
	}
	if !online {
		mode = "offline"
	}
	meta := map[string]any{
		"agent_ref":        a.Ref,
		"mode":             mode,
		"connection":       map[bool]string{true: "online", false: "offline"}[online],
		"position_kind":    kind,
		"last_heartbeat_t": a.LastSeen.UnixMilli(),
		"protocol":         "synadia-agent-protocol/" + agentproto.ProtocolVersion,
	}
	if t.cfg.NodeID != "" {
		meta["node_id"] = t.cfg.NodeID
	}
	if a.State.Activity != "" {
		meta["activity"] = a.State.Activity
	}
	if a.State.Role != "" {
		meta["role"] = a.State.Role
	}
	if ct := a.State.CurrentTask; ct != nil && ct.TaskID != "" {
		meta["current_task"] = map[string]any{"task_id": ct.TaskID, "state": ct.State, "summary": ct.Summary}
	} else {
		meta["current_task"] = nil
	}
	if a.Prompt != "" {
		meta["prompt_subject"] = a.Prompt
	}
	stale := int((3 * a.Interval).Seconds())
	if stale < 90 {
		stale = 90
	}
	return Envelope{UID: a.UID, Type: "fleet.agent", Lat: an.Lat, Lon: an.Lon, HAE: 0, Course: -1, Speed: -1,
		Callsign: a.Callsign, Stale: stale, T: now.UnixMilli(), CE: 50, LE: 9999999, Meta: meta}
}

// Due returns the envelopes to publish now: a changed state, or the periodic refresh.
// Offline agents are published until ForgetAfter, then dropped.
func (t *Tracker) Due(now time.Time, force bool) []Envelope {
	t.mu.Lock()
	defer t.mu.Unlock()
	var out []Envelope
	for uid, a := range t.agents {
		if now.Sub(a.LastSeen) > t.cfg.ForgetAfter {
			delete(t.agents, uid)
			delete(t.last, uid)
			delete(t.lastT, uid)
			continue
		}
		env := t.EnvelopeFor(*a, now)
		digest := fmt.Sprintf("%v|%v|%v|%v|%v|%v", env.Meta["mode"], env.Meta["activity"], env.Meta["current_task"], env.Lat, env.Lon, env.Callsign)
		if force || digest != t.last[uid] || now.Sub(t.lastT[uid]) >= t.cfg.PublishEvery {
			t.last[uid], t.lastT[uid] = digest, now
			out = append(out, env)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UID < out[j].UID })
	return out
}

// Subject returns `tel.fleet.agent.{uid}.position`.
func Subject(uid string) string { return "tel.fleet.agent." + uid + ".position" }

// Run follows heartbeats and publishes envelopes until ctx ends.
func (t *Tracker) Run(ctx context.Context, nc *nats.Conn, pub *bus.Publisher, agentsStream string) error {
	kick := make(chan struct{}, 1)
	poke := func() {
		select {
		case kick <- struct{}{}:
		default:
		}
	}
	if _, err := bus.Follow(ctx, nc, pub.JS, agentsStream, []string{"agents.hb.>"}, true, func(subj string, data []byte, _ nats.Header) {
		if _, err := t.OnHeartbeat(subj, data); err == nil {
			poke()
		}
	}); err != nil {
		return err
	}
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()
	publish := func() {
		now := t.now()
		for _, env := range t.Due(now, false) {
			b, _ := json.Marshal(env)
			id := fmt.Sprintf("%s:%d", env.UID, env.T)
			_, _ = pub.Publish(ctx, Subject(env.UID), b, id, nil)
		}
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-kick:
			publish()
		case <-tick.C:
			publish()
		}
	}
}
