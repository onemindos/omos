package projector

// The projection engine: per-entity state, hysteresis, cap (G9, §3.3), per-uid intervals,
// heartbeats and the global ceiling (§5.3, §6.1). It has no NATS code and is driven from one
// goroutine (Service.loop), so it needs no locks except for the counters.

import (
	"log/slog"
	"sort"
	"sync"
	"time"
)

// Counters are read by the HTTP status handler from another goroutine.
type Counters struct {
	mu sync.Mutex
	m  map[string]uint64
}

func (c *Counters) Inc(k string) { c.Add(k, 1) }

func (c *Counters) Add(k string, n uint64) {
	c.mu.Lock()
	if c.m == nil {
		c.m = map[string]uint64{}
	}
	c.m[k] += n
	c.mu.Unlock()
}

func (c *Counters) Get(k string) uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.m[k]
}

func (c *Counters) Snapshot() map[string]uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]uint64, len(c.m))
	for k, v := range c.m {
		out[k] = v
	}
	return out
}

type track struct {
	e         *Entity
	dec       decision
	admitted  bool // counted against MAX_PROJECTED and scheduled for sending
	projected bool // a CoT for this uid is live on the target (sent and not removed)
	lastSent  time.Time
	lastType  string
	exitAt    time.Time // hysteresis deadline after leaving every AOI
	gone      bool      // deleted from ENT; forget after the removal
}

type pending struct {
	removal bool
	due     time.Time
	prio    int
	dist    float64
}

// Publisher sends one CoT for an entity.
type Publisher func(b Built, e *Entity) error

// Engine holds the projection state.
type Engine struct {
	cfg      Config
	log      *slog.Logger
	pub      Publisher
	C        *Counters
	enabled  bool
	idx      *aoiIndex
	tracks   map[string]*track
	pend     map[string]*pending
	admitted int
	sent     []time.Time // send times within the last second (global ceiling ring)
	csent    []time.Time // svc.cot.send publishes within the last second
	lastCap  time.Time
	lastHB   time.Time
}

// NewEngine builds an engine; pub is called from Tick only.
func NewEngine(cfg Config, log *slog.Logger, pub Publisher) *Engine {
	return &Engine{cfg: cfg, log: log, pub: pub, C: &Counters{}, enabled: true, idx: newAOIIndex(),
		tracks: map[string]*track{}, pend: map[string]*pending{}}
}

func (g *Engine) interval(e *Entity) time.Duration {
	if e.IsOwnAsset() {
		return g.cfg.OwnInterval
	}
	return g.cfg.OtherInterval
}

// Upsert applies one entity record.
func (g *Engine) Upsert(e *Entity, now time.Time) {
	k := e.Key()
	t := g.tracks[k]
	if t == nil {
		t = &track{}
		g.tracks[k] = t
	}
	t.e, t.gone = e, false
	g.apply(k, t, now)
	if !t.admitted && !t.projected && !e.Live {
		delete(g.tracks, k) // dead and never on TAK: nothing to remember
		delete(g.pend, k)
	}
}

// Delete handles an ENT delete marker or purge for "{type}.{id}".
func (g *Engine) Delete(key string, now time.Time) {
	t := g.tracks[key]
	if t == nil {
		return
	}
	t.gone = true
	if t.admitted {
		g.unadmit(key, t, ReasonDeleted, now)
	}
	if !t.projected {
		delete(g.tracks, key)
		delete(g.pend, key)
	}
}

// SetEnabled is G1 (cmd.projector.enable / disable).
func (g *Engine) SetEnabled(on bool, now time.Time) {
	if g.enabled == on {
		return
	}
	g.enabled = on
	g.reevaluateAll(now)
}

// Zone and mission changes re-run the gates for every entity (§3.2 "Recompute").
func (g *Engine) SetZone(key string, data []byte, now time.Time) {
	if !g.idx.setZone(key, data) {
		g.C.Inc("geofence_unreadable")
	}
	g.reevaluateAll(now)
}

func (g *Engine) DelZone(key string, now time.Time) { g.idx.delZone(key); g.reevaluateAll(now) }

func (g *Engine) SetMission(key string, data []byte, now time.Time) {
	if !g.idx.setMission(key, data) {
		g.C.Inc("mission_unreadable")
	}
	g.reevaluateAll(now)
}

func (g *Engine) DelMission(key string, now time.Time) { g.idx.delMission(key); g.reevaluateAll(now) }

func (g *Engine) reevaluateAll(now time.Time) {
	for k, t := range g.tracks {
		if !t.gone {
			g.apply(k, t, now)
		}
	}
}

func (g *Engine) apply(k string, t *track, now time.Time) {
	d := evaluate(t.e, &g.cfg, g.idx, g.enabled, now)
	t.dec = d
	if d.downgraded {
		g.C.Inc("civilian_downgrade")
	}
	if d.ok {
		t.exitAt = time.Time{}
		if !t.admitted && g.admitted < g.cfg.MaxProjected {
			t.admitted = true
			g.admitted++
		}
		if t.admitted {
			g.schedule(k, t, now, false)
		}
		return
	}
	g.C.Inc("blocked_" + d.reason)
	if t.admitted && d.onlyOutside && g.cfg.AOIExitGrace > 0 {
		if t.exitAt.IsZero() {
			t.exitAt = now.Add(g.cfg.AOIExitGrace)
		}
		g.schedule(k, t, now, false) // still projected during the grace period
		return
	}
	if t.admitted {
		g.unadmit(k, t, d.reason, now)
	}
}

func (g *Engine) unadmit(k string, t *track, reason string, now time.Time) {
	t.admitted = false
	t.exitAt = time.Time{}
	g.admitted--
	if t.projected {
		g.pend[k] = &pending{removal: true, due: now, prio: 0}
		g.log.Debug("removal queued", "uid", "OM."+k, "reason", reason)
		g.C.Inc("removal_" + reason)
	} else {
		delete(g.pend, k)
	}
}

func (g *Engine) schedule(k string, t *track, now time.Time, heartbeat bool) {
	due := now
	if !t.lastSent.IsZero() {
		if next := t.lastSent.Add(g.interval(t.e)); next.After(due) {
			due = next
		}
	}
	if p := g.pend[k]; p != nil && !p.removal && !heartbeat {
		p.prio, p.dist = t.dec.rank, t.dec.dist // keep the earlier due time
		return
	}
	g.pend[k] = &pending{due: due, prio: t.dec.rank, dist: t.dec.dist}
}

// Tick runs the timers and sends what the ceiling allows.
func (g *Engine) Tick(now time.Time) {
	// Hysteresis: entities whose grace ran out.
	for k, t := range g.tracks {
		if t.admitted && !t.exitAt.IsZero() && !now.Before(t.exitAt) {
			d := evaluate(t.e, &g.cfg, g.idx, g.enabled, now)
			t.dec = d
			if d.ok {
				t.exitAt = time.Time{}
			} else {
				g.unadmit(k, t, d.reason, now)
			}
		}
	}
	if now.Sub(g.lastCap) >= g.cfg.CapEvery {
		g.lastCap = now
		g.enforceCap(now)
	}
	if now.Sub(g.lastHB) >= time.Second {
		g.lastHB = now
		g.heartbeats(now)
	}
	g.flush(now)
}

// enforceCap keeps the MAX_PROJECTED best entities (§3.3) when more are eligible.
func (g *Engine) enforceCap(now time.Time) {
	type cand struct {
		k string
		t *track
	}
	var elig []cand
	for k, t := range g.tracks {
		if t.dec.ok || (t.admitted && !t.exitAt.IsZero()) {
			elig = append(elig, cand{k, t})
		}
	}
	if len(elig) <= g.cfg.MaxProjected {
		for _, c := range elig { // room for all: admit any that waited
			if !c.t.admitted {
				c.t.admitted = true
				g.admitted++
				g.schedule(c.k, c.t, now, false)
			}
		}
		return
	}
	rank := func(t *track) (int, float64) {
		r := t.dec.rank
		if r == 0 { // in its grace period: keep its last rank class at the back
			r = RankAOI
		}
		return r, t.dec.dist
	}
	sort.Slice(elig, func(i, j int) bool {
		ri, di := rank(elig[i].t)
		rj, dj := rank(elig[j].t)
		if ri != rj {
			return ri < rj
		}
		if di != dj {
			return di < dj
		}
		return elig[i].k < elig[j].k
	})
	stopped := 0
	for i, c := range elig {
		if i < g.cfg.MaxProjected {
			if !c.t.admitted {
				c.t.admitted = true
				g.admitted++
				g.schedule(c.k, c.t, now, false)
			}
			continue
		}
		stopped++
		if c.t.admitted {
			g.unadmit(c.k, c.t, ReasonCap, now)
		}
	}
	g.C.Add("cap_stopped", uint64(stopped))
	g.log.Info("cap reached", "max_projected", g.cfg.MaxProjected, "eligible", len(elig), "stopped", stopped)
}

// heartbeats re-send unchanged projected entities every min(stale/2, 300 s), except stale ones.
func (g *Engine) heartbeats(now time.Time) {
	for k, t := range g.tracks {
		if !t.admitted || !t.projected || g.pend[k] != nil || t.e.Lifecycle == "stale" {
			continue
		}
		hb := time.Duration(StaleSeconds(t.e, g.interval(t.e))) * time.Second / 2
		if hb > g.cfg.HeartbeatMax {
			hb = g.cfg.HeartbeatMax
		}
		if now.Sub(t.lastSent) >= hb {
			g.schedule(k, t, now, true)
			g.C.Inc("heartbeat")
		}
	}
}

func (g *Engine) budget(now time.Time) int {
	cut := now.Add(-time.Second)
	i := 0
	for i < len(g.sent) && !g.sent[i].After(cut) {
		i++
	}
	g.sent = g.sent[i:]
	return g.cfg.MaxRatePerSec - len(g.sent)
}

// clientBudget reports whether one svc.cot.send publish fits now: under the client ceiling
// (COT_SEND_MAX_PER_SEC) and under the global ceiling it shares with projected entities.
func (g *Engine) clientBudget(now time.Time) bool {
	cut := now.Add(-time.Second)
	i := 0
	for i < len(g.csent) && !g.csent[i].After(cut) {
		i++
	}
	g.csent = g.csent[i:]
	return len(g.csent) < g.cfg.Send.MaxPerSec && g.budget(now) > 0
}

// noteClientSend counts a client publish in both windows.
func (g *Engine) noteClientSend(now time.Time) {
	g.csent = append(g.csent, now)
	g.sent = append(g.sent, now)
	g.C.Inc("client_sent")
}

// flush sends due items, removals first, then by cap order, within the global ceiling: at most
// MAX_RATE_PER_SEC messages in any one-second window. Items over the ceiling stay queued and
// are rebuilt from the newest record when sent, so only the newest per uid goes out.
func (g *Engine) flush(now time.Time) {
	free := g.budget(now)
	if free <= 0 || len(g.pend) == 0 {
		return
	}
	type item struct {
		k string
		p *pending
	}
	due := make([]item, 0, 64)
	for k, p := range g.pend {
		if !p.due.After(now) {
			due = append(due, item{k, p})
		}
	}
	if len(due) == 0 {
		return
	}
	if len(due) > free {
		sort.Slice(due, func(i, j int) bool {
			a, b := due[i].p, due[j].p
			if a.prio != b.prio {
				return a.prio < b.prio
			}
			if a.dist != b.dist {
				return a.dist < b.dist
			}
			if !a.due.Equal(b.due) {
				return a.due.Before(b.due)
			}
			return due[i].k < due[j].k
		})
		g.C.Add("ceiling_deferred", uint64(len(due)-free))
	}
	for _, it := range due {
		if free <= 0 {
			break
		}
		t := g.tracks[it.k]
		delete(g.pend, it.k)
		if t == nil {
			continue
		}
		if g.send(it.k, t, it.p.removal, now) {
			free--
			g.sent = append(g.sent, now)
		}
	}
}

func (g *Engine) send(k string, t *track, removal bool, now time.Time) bool {
	ct := t.dec.cotType
	if removal || ct == "" {
		ct = t.lastType
	}
	if ct == "" {
		return false
	}
	b, err := Build(t.e, ct, now, buildOpts{target: g.cfg.Target, interval: g.interval(t.e),
		maxVertices: g.cfg.MaxVertices, maxBytes: g.cfg.MeshMaxBytes, removal: removal})
	if err != nil {
		g.C.Inc("build_error")
		g.log.Warn("cot not built", "uid", "OM."+k, "err", err)
		return false
	}
	if err := g.pub(b, t.e); err != nil {
		g.C.Inc("publish_error")
		g.log.Warn("publish failed", "uid", b.UID, "err", err)
		if !removal {
			g.pend[k] = &pending{due: now.Add(time.Second), prio: t.dec.rank, dist: t.dec.dist}
		} else {
			g.pend[k] = &pending{removal: true, due: now.Add(time.Second)}
		}
		return false
	}
	if removal {
		t.projected = false
		g.C.Inc("removals_sent")
		if t.gone || !t.admitted {
			if t.gone || !t.e.Live {
				delete(g.tracks, k)
			}
		}
		return true
	}
	t.projected, t.lastSent, t.lastType = true, now, ct
	g.C.Inc("cot_sent")
	return true
}

// Stats is the status summary.
func (g *Engine) Stats() map[string]any {
	return map[string]any{"entities": len(g.tracks), "admitted": g.admitted, "pending": len(g.pend),
		"aoi_active": len(g.idx.active), "enabled": g.enabled}
}
