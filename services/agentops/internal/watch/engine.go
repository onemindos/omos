package watch

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/onemindos/omos/services/agentops/internal/feeds"
	"github.com/onemindos/omos/services/agentops/internal/geo"
)

// Entity is the subset of an entity.md record the rules read; Raw keeps the whole record
// for threshold paths.
type Entity struct {
	ID        string
	Type      string
	Callsign  string
	Lat, Lon  float64
	HasPos    bool
	T         int64
	Seq       int64
	Lifecycle string
	Live      bool
	Subject   string
	Raw       map[string]any
}

// ParseEntity decodes `ent.{type}.{id}.state`.
func ParseEntity(subject string, data []byte) (*Entity, error) {
	tok := strings.Split(subject, ".")
	if len(tok) != 4 || tok[0] != "ent" || tok[3] != "state" {
		return nil, fmt.Errorf("not an entity subject: %s", subject)
	}
	raw := map[string]any{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	e := &Entity{ID: tok[2], Type: tok[1], Subject: subject, Raw: raw, Live: true}
	if s, ok := raw["callsign"].(string); ok {
		e.Callsign = s
	}
	if s, ok := raw["lifecycle"].(string); ok {
		e.Lifecycle = s
	}
	if b, ok := raw["live"].(bool); ok {
		e.Live = b
	}
	e.T = int64(num(raw["t"]))
	e.Seq = int64(num(raw["seq"]))
	if loc, ok := raw["location"].(map[string]any); ok {
		la, lo := loc["lat"], loc["lon"]
		if _, ok1 := la.(float64); ok1 {
			if _, ok2 := lo.(float64); ok2 {
				e.Lat, e.Lon, e.HasPos = la.(float64), lo.(float64), true
			}
		}
	}
	if !e.HasPos {
		if la, ok := raw["lat"].(float64); ok {
			if lo, ok := raw["lon"].(float64); ok && !(la == 0 && lo == 0) {
				e.Lat, e.Lon, e.HasPos = la, lo, true
			}
		}
	}
	return e, nil
}

func num(v any) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

// Field returns a numeric value at a dotted path.
func (e *Entity) Field(path string) (float64, bool) {
	var cur any = e.Raw
	for _, p := range strings.Split(path, ".") {
		m, ok := cur.(map[string]any)
		if !ok {
			return 0, false
		}
		cur = m[p]
	}
	f, ok := cur.(float64)
	return f, ok && !math.IsNaN(f)
}

// Trigger is one fired rule: the runner turns it into an alert.
type Trigger struct {
	Watch     *Watch
	AlertID   string
	Kind      string // alert_kind: geofence_exit, geofence_enter, asset_stale, asset_lost, threshold, feed_stale ...
	Title     string
	AboutType string
	AboutID   string
	Subject   string
	Evidence  map[string]any
	Lat, Lon  *float64
	T         time.Time
}

// GeofenceEvent is a crossing for doctrine §9.1 `evt.map.geofence.{id}.{enter|exit}`.
type GeofenceEvent struct {
	Zone       string
	Transition string // enter | exit
	EntityType string
	EntityID   string
	Callsign   string
	Lat, Lon   float64
	T          time.Time
	WatchID    string
}

type wkey struct{ watch, about string }

type wstate struct {
	active    bool
	inside    *bool
	lastAlert time.Time
}

// Engine evaluates the rules. It is safe for concurrent use.
type Engine struct {
	mu       sync.Mutex
	watches  []*Watch
	now      func() time.Time
	zones    map[string]geo.Zone
	badZones map[string]string
	ents     map[string]*Entity // by entity id, only entities some rule selects
	st       map[wkey]*wstate
	alerts   map[string][]time.Time // per watch, last hour
	feeds    *feeds.Board

	Suppressed map[string]int // reason -> count (cooldown, hourly_cap)
}

// NewEngine returns an engine over the enabled watches.
func NewEngine(ws []Watch, board *feeds.Board, now func() time.Time) *Engine {
	if now == nil {
		now = time.Now
	}
	e := &Engine{now: now, zones: map[string]geo.Zone{}, badZones: map[string]string{}, ents: map[string]*Entity{},
		st: map[wkey]*wstate{}, alerts: map[string][]time.Time{}, feeds: board, Suppressed: map[string]int{}}
	for i := range ws {
		if !ws[i].Disabled {
			w := ws[i]
			e.watches = append(e.watches, &w)
		}
	}
	return e
}

// Watches returns the enabled watches.
func (e *Engine) Watches() []*Watch { return e.watches }

func selects(w *Watch, ent *Entity) bool {
	ok := false
	for _, t := range w.Rule.EntityTypes {
		if t == ent.Type {
			ok = true
		}
	}
	if !ok {
		return false
	}
	if len(w.Rule.EntityIDs) == 0 {
		return true
	}
	for _, id := range w.Rule.EntityIDs {
		if strings.HasSuffix(id, "*") {
			if strings.HasPrefix(ent.ID, strings.TrimSuffix(id, "*")) {
				return true
			}
		} else if id == ent.ID {
			return true
		}
	}
	return false
}

func (e *Engine) state(k wkey) *wstate {
	s := e.st[k]
	if s == nil {
		s = &wstate{}
		e.st[k] = s
	}
	return s
}

// edge records the rule condition for (watch, about) and returns true when an alert must
// be raised: a rising edge that passes the cooldown and the hourly cap.
func (e *Engine) edge(w *Watch, about string, cond bool, now time.Time) bool {
	s := e.state(wkey{w.ID, about})
	was := s.active
	s.active = cond
	if !cond || was {
		return false
	}
	if !s.lastAlert.IsZero() && now.Sub(s.lastAlert) < time.Duration(w.Rate.CooldownS*float64(time.Second)) {
		e.Suppressed["cooldown"]++
		return false
	}
	cut := now.Add(-time.Hour)
	h := e.alerts[w.ID]
	for len(h) > 0 && h[0].Before(cut) {
		h = h[1:]
	}
	if len(h) >= w.Rate.MaxAlertsPerHour {
		e.alerts[w.ID] = h
		e.Suppressed["hourly_cap"]++
		return false
	}
	e.alerts[w.ID] = append(h, now)
	s.lastAlert = now
	return true
}

// OnZone records a geofences KV value (nil value = deleted).
func (e *Engine) OnZone(key string, value []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if value == nil {
		delete(e.zones, key)
		delete(e.badZones, key)
		return nil
	}
	z, err := geo.ParseZone(key, value)
	if err != nil {
		e.badZones[key] = err.Error()
		delete(e.zones, key)
		return err
	}
	delete(e.badZones, key)
	e.zones[key] = z
	if z.ID != key {
		e.zones[z.ID] = z
	}
	return nil
}

// ZoneProblems lists zones some watch names that cannot be evaluated.
func (e *Engine) ZoneProblems() map[string]string {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := map[string]string{}
	for _, w := range e.watches {
		if w.Rule.Kind != KindZone {
			continue
		}
		if msg, bad := e.badZones[w.Rule.Zone]; bad {
			out[w.Rule.Zone] = msg
		} else if _, ok := e.zones[w.Rule.Zone]; !ok {
			out[w.Rule.Zone] = "zone not found in KV geofences"
		}
	}
	return out
}

// OnEntity evaluates the event-driven rules (zone, threshold, and stale on deletion) for one
// entity record.
func (e *Engine) OnEntity(subject string, data []byte) ([]Trigger, []GeofenceEvent, error) {
	ent, err := ParseEntity(subject, data)
	if err != nil {
		return nil, nil, err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.now()
	var trig []Trigger
	var gev []GeofenceEvent
	interested := false
	for _, w := range e.watches {
		if w.Rule.Kind == KindFeedStale || !selects(w, ent) {
			continue
		}
		interested = true
		k := wkey{w.ID, ent.ID}
		switch w.Rule.Kind {
		case KindZone:
			z, ok := e.zones[w.Rule.Zone]
			if !ok || !z.Live || !ent.HasPos || !ent.Live {
				continue
			}
			in := z.Contains(ent.Lat, ent.Lon)
			s := e.state(k)
			if s.inside != nil && *s.inside != in {
				tr := "exit"
				if in {
					tr = "enter"
				}
				gev = append(gev, GeofenceEvent{Zone: z.ID, Transition: tr, EntityType: ent.Type, EntityID: ent.ID,
					Callsign: ent.Callsign, Lat: ent.Lat, Lon: ent.Lon, T: now, WatchID: w.ID})
			}
			s.inside = &in
			cond := (w.Rule.Alert == "outside" && !in) || (w.Rule.Alert == "inside" && in)
			if e.edge(w, ent.ID, cond, now) {
				kind, verb := "geofence_exit", "is outside"
				if w.Rule.Alert == "inside" {
					kind, verb = "geofence_enter", "is inside"
				}
				lat, lon := ent.Lat, ent.Lon
				trig = append(trig, e.newTrigger(w, ent, kind, fmt.Sprintf("%s %s zone %s (%s)", name(ent), verb, z.ID, z.Kind), now,
					map[string]any{"zone": z.ID, "zone_kind": z.Kind, "lat": round6(lat), "lon": round6(lon), "seq": ent.Seq, "subject": subject}, &lat, &lon))
			}
		case KindThreshold:
			v, ok := e.thresholdValue(w, ent)
			if !ok {
				continue
			}
			s := e.state(k)
			cond := compare(w.Rule.Op, v, w.Rule.Value)
			if s.active && !cond && w.Rule.Clear != nil && compare(w.Rule.Op, v, *w.Rule.Clear) {
				cond = true // hysteresis: stay active until the clear level is passed
			}
			if e.edge(w, ent.ID, cond, now) {
				var lat, lon *float64
				if ent.HasPos {
					la, lo := ent.Lat, ent.Lon
					lat, lon = &la, &lo
				}
				trig = append(trig, e.newTrigger(w, ent, "threshold", fmt.Sprintf("%s %s is %s %s %s", name(ent), w.Rule.Field, trimNum(v), opWord(w.Rule.Op), trimNum(w.Rule.Value)), now,
					map[string]any{"field": w.Rule.Field, "value": v, "threshold": w.Rule.Value, "op": w.Rule.Op, "seq": ent.Seq, "subject": subject}, lat, lon))
			}
		}
	}
	if interested {
		if ent.Live {
			e.ents[ent.ID] = ent
		} else {
			// A final live:false record: a stale watch sees it as lost at once.
			trig = append(trig, e.sweepOne(ent, now, true)...)
			delete(e.ents, ent.ID)
		}
	}
	return trig, gev, nil
}

func (e *Engine) thresholdValue(w *Watch, ent *Entity) (float64, bool) {
	return ent.Field(w.Rule.Field)
}

// Sweep evaluates the time-driven rules (stale entities, feed health).
func (e *Engine) Sweep() []Trigger {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.now()
	var out []Trigger
	ids := make([]string, 0, len(e.ents))
	for id := range e.ents {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		out = append(out, e.sweepOne(e.ents[id], now, false)...)
	}
	if e.feeds != nil {
		gaps := e.feeds.Gaps()
		bySrc := map[string]feeds.Gap{}
		for _, g := range gaps {
			bySrc[g.Source] = g
		}
		for _, w := range e.watches {
			if w.Rule.Kind != KindFeedStale {
				continue
			}
			srcs := w.Rule.Sources
			if len(srcs) == 0 {
				for s := range bySrc {
					srcs = append(srcs, s)
				}
				for k := range e.st {
					if k.watch == w.ID && strings.HasPrefix(k.about, "feed-") {
						srcs = append(srcs, strings.TrimPrefix(k.about, "feed-"))
					}
				}
				sort.Strings(srcs)
			}
			done := map[string]bool{}
			for _, src := range srcs {
				if done[src] {
					continue
				}
				done[src] = true
				g, bad := bySrc[src]
				if bad && len(w.Rule.Reasons) > 0 && !contains(w.Rule.Reasons, g.Reason) {
					bad = false
				}
				about := FeedID(src)
				if e.edge(w, about, bad, now) {
					ev := map[string]any{"source": src, "reason": g.Reason, "detail": g.Detail, "subject": "tel.feed." + src + ".status"}
					if g.LastError != "" {
						ev["last_error"] = g.LastError
					}
					out = append(out, Trigger{Watch: w, AlertID: NewID("alr"), Kind: "feed_" + g.Reason,
						Title: fmt.Sprintf("feed %s is %s: %s", src, g.Reason, g.Detail), AboutType: "feed", AboutID: about,
						Subject: "tel.feed." + src + ".status", Evidence: ev, T: now})
				}
			}
		}
	}
	return out
}

func (e *Engine) sweepOne(ent *Entity, now time.Time, deleted bool) []Trigger {
	var out []Trigger
	for _, w := range e.watches {
		if w.Rule.Kind != KindStale || !selects(w, ent) {
			continue
		}
		age := now.Sub(time.UnixMilli(ent.T)).Seconds()
		lost := deleted || !ent.Live || ent.Lifecycle == "lost" || ent.Lifecycle == "expired"
		cond := lost || ent.Lifecycle == "stale" || (ent.T > 0 && age > w.Rule.MaxAgeS)
		if e.edge(w, ent.ID, cond, now) {
			kind, title := "asset_stale", fmt.Sprintf("%s has not updated for %s (limit %s)", name(ent), dur(age), dur(w.Rule.MaxAgeS))
			if lost {
				kind, title = "asset_lost", fmt.Sprintf("%s is lost (lifecycle %s, last update %s ago)", name(ent), orStr(ent.Lifecycle, "deleted"), dur(age))
			}
			var lat, lon *float64
			if ent.HasPos {
				la, lo := ent.Lat, ent.Lon
				lat, lon = &la, &lo
			}
			out = append(out, e.newTrigger(w, ent, kind, title, now,
				map[string]any{"age_s": math.Round(age), "max_age_s": w.Rule.MaxAgeS, "lifecycle": ent.Lifecycle, "seq": ent.Seq, "subject": ent.Subject}, lat, lon))
		}
	}
	return out
}

func (e *Engine) newTrigger(w *Watch, ent *Entity, kind, title string, now time.Time, ev map[string]any, lat, lon *float64) Trigger {
	return Trigger{Watch: w, AlertID: NewID("alr"), Kind: kind, Title: title, AboutType: ent.Type, AboutID: ent.ID,
		Subject: ent.Subject, Evidence: ev, Lat: lat, Lon: lon, T: now}
}

// FeedID is the alert `about` id of a feed: feed-{source}, a valid subject token.
func FeedID(source string) string {
	s := "feed-" + strings.Trim(badTok.ReplaceAllString(strings.ToLower(source), "-"), "-")
	if len(s) > 64 {
		s = s[:64]
	}
	return s
}

var badTok = regexp.MustCompile(`[^a-z0-9_-]+`)

func compare(op string, v, ref float64) bool {
	switch op {
	case "lt":
		return v < ref
	case "le":
		return v <= ref
	case "gt":
		return v > ref
	case "ge":
		return v >= ref
	}
	return false
}

func opWord(op string) string {
	return map[string]string{"lt": "below", "le": "at or below", "gt": "above", "ge": "at or above"}[op]
}

func name(e *Entity) string {
	if e.Callsign != "" {
		return e.Callsign + " (" + e.ID + ")"
	}
	return e.ID
}

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func orStr(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

func round6(f float64) float64 { return math.Round(f*1e6) / 1e6 }

func trimNum(f float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", f), "0"), ".")
}

func dur(s float64) string { return (time.Duration(s) * time.Second).Round(time.Second).String() }

const crockford = "0123456789abcdefghjkmnpqrstvwxyz"

// NewID returns "{prefix}-{ulid}" with a lowercase ULID (entity.md §3 id format).
func NewID(prefix string) string {
	var b [16]byte
	ms := uint64(time.Now().UnixMilli())
	for i := 5; i >= 0; i-- {
		b[i] = byte(ms)
		ms >>= 8
	}
	_, _ = rand.Read(b[6:])
	// 128 bits -> 26 base32 chars
	out := make([]byte, 26)
	var acc uint64
	bits := 0
	idx := 0
	// Encode with 2 leading pad bits so 130 bits map onto 26 chars.
	acc, bits = 0, 2
	for _, c := range b {
		acc = acc<<8 | uint64(c)
		bits += 8
		for bits >= 5 {
			bits -= 5
			out[idx] = crockford[(acc>>uint(bits))&31]
			idx++
		}
	}
	return prefix + "-" + string(out[:idx])
}
