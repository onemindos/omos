// Package engine is the geofence watcher and alert logic (FEATURES D1, D6), free of NATS so it
// can be tested on its own. The bus layer feeds it zones (KV `geofences`), entities
// (`ent.*.*.state`) and feed status (`tel.feed.{source}.status`) and publishes what it returns.
//
// Crossings: for every watched, active zone that applies to an entity, a change of
// inside/outside is one `evt.map.geofence.{zone}.{enter|exit}` fact. Membership starts as
// "outside", so the first fix inside a zone is an `enter` (cause `first_fix`); every inside
// membership therefore has an `enter` on the MISSION stream, and replaying those facts on
// start rebuilds the membership exactly (restart idempotence).
//
// Alerts: one alert per condition with a deterministic id; it is raised once, updated when
// its level or reason changes, and resolved when the condition clears. Alert facts are
// `evt.so.alert.{raised|updated|resolved}`; replaying them on start rebuilds the open set, so
// a restart never raises an open alert twice.
package engine

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/onemindos/omos/services/geofence/internal/geo"
	"github.com/onemindos/omos/services/geofence/internal/zone"
)

// Alert kinds, levels and states (ONTOLOGY.md alert row: `alert_kind`, `level`
// advisory|caution|warning, `state` open|acknowledged|resolved).
const (
	KindGeofence   = "geofence"
	KindStaleAsset = "stale_asset"
	KindFeedOutage = "feed_outage"

	LevelAdvisory = "advisory"
	LevelCaution  = "caution"
	LevelWarning  = "warning"

	StateOpen     = "open"
	StateResolved = "resolved"
)

// Crossing causes.
const (
	CauseCrossing    = "crossing"
	CauseFirstFix    = "first_fix"
	CauseZoneChange  = "zone_change"
	CauseZoneRemoved = "zone_removed"
	CauseEntityGone  = "entity_gone"
)

// Config tunes the engine.
type Config struct {
	Source       string          // envelope source and alert raised_by
	AssetTypes   map[string]bool // own assets: stale alerts and default keep-in scope
	FeedSilence  time.Duration   // no tel.feed status for this long = outage (silent)
	StaleDefault int             // seconds, for an asset record without `stale`
	CellDeg      float64
	Now          func() time.Time
}

// DefaultAssetTypes are the commandable classes (doctrine §9.1).
func DefaultAssetTypes() map[string]bool {
	return map[string]bool{"drone": true, "robot": true, "vehicle": true}
}

// nonMoving entity types never cross a zone.
var nonMoving = map[string]bool{"zone": true, "place": true, "structure": true, "event": true, "alert": true, "team": true}

// Entity is the part of an entity.v1 record the engine reads.
type Entity struct {
	EntityID   string
	EntityType string
	Kind       string
	Template   string
	Callsign   string
	Lat, Lon   float64
	HAE        *float64
	HasPos     bool
	Live       bool
	Lifecycle  string
	Seq        uint64
	T          int64
	Stale      int
	Simulated  bool
}

// ParseEntity reads an `ent.*.*.state` record.
func ParseEntity(b []byte) (Entity, error) {
	var r struct {
		EntityID   string   `json:"entity_id"`
		EntityType string   `json:"entity_type"`
		Kind       string   `json:"kind"`
		Template   string   `json:"template"`
		Callsign   string   `json:"callsign"`
		Lat        float64  `json:"lat"`
		Lon        float64  `json:"lon"`
		HAE        *float64 `json:"hae"`
		Live       bool     `json:"live"`
		Lifecycle  string   `json:"lifecycle"`
		Seq        uint64   `json:"seq"`
		T          int64    `json:"t"`
		Stale      int      `json:"stale"`
		Traits     []string `json:"traits"`
		Indicators struct {
			Simulated bool `json:"simulated"`
		} `json:"indicators"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return Entity{}, err
	}
	if r.EntityID == "" || r.EntityType == "" {
		return Entity{}, fmt.Errorf("record has no entity_id or entity_type")
	}
	e := Entity{EntityID: r.EntityID, EntityType: r.EntityType, Kind: r.Kind, Template: r.Template,
		Callsign: r.Callsign, Lat: r.Lat, Lon: r.Lon, HAE: r.HAE, Live: r.Live, Lifecycle: r.Lifecycle,
		Seq: r.Seq, T: r.T, Stale: r.Stale, Simulated: r.Indicators.Simulated}
	locatable := false
	for _, t := range r.Traits {
		if t == "locatable" {
			locatable = true
		}
	}
	// entity.md §4.1: non-located records carry lat 0, lon 0 and no `locatable` trait.
	e.HasPos = (locatable || r.Lat != 0 || r.Lon != 0) && r.Lat >= -90 && r.Lat <= 90 && r.Lon >= -180 && r.Lon <= 180
	return e, nil
}

// LatLon is a position.
type LatLon struct {
	Lat float64  `json:"lat"`
	Lon float64  `json:"lon"`
	HAE *float64 `json:"hae,omitempty"`
}

// About names what an alert is about.
type About struct {
	EntityType string `json:"entity_type,omitempty"`
	EntityID   string `json:"entity_id,omitempty"`
	Feed       string `json:"feed,omitempty"`
}

// Alert is one alert (services/schemas/alert/alert.v1.schema.json).
type Alert struct {
	V         int     `json:"v"`
	AlertID   string  `json:"alert_id"`
	Kind      string  `json:"alert_kind"`
	Level     string  `json:"level"`
	State     string  `json:"state"`
	About     About   `json:"about"`
	ZoneID    string  `json:"zone_id,omitempty"`
	ZoneKind  string  `json:"zone_kind,omitempty"`
	Reason    string  `json:"reason"`
	Message   string  `json:"message"`
	RaisedBy  string  `json:"raised_by"`
	OpenedT   int64   `json:"opened_t"`
	UpdatedT  int64   `json:"updated_t"`
	ResolvedT *int64  `json:"resolved_t,omitempty"`
	Seq       uint64  `json:"seq"`
	Location  *LatLon `json:"location,omitempty"`
	Simulated bool    `json:"simulated"`
}

// Crossing is the data of `evt.map.geofence.{zone}.{enter|exit}`.
type Crossing struct {
	V          int      `json:"v"`
	ZoneID     string   `json:"zone_id"`
	ZoneKind   string   `json:"zone_kind,omitempty"`
	ZoneName   string   `json:"zone_name,omitempty"`
	ZoneRev    uint64   `json:"zone_rev,omitempty"`
	Event      string   `json:"event"`
	Cause      string   `json:"cause"`
	EntityID   string   `json:"entity_id"`
	EntityType string   `json:"entity_type,omitempty"`
	Callsign   string   `json:"callsign,omitempty"`
	Lat        *float64 `json:"lat,omitempty"`
	Lon        *float64 `json:"lon,omitempty"`
	HAE        *float64 `json:"hae,omitempty"`
	EntityT    int64    `json:"entity_t,omitempty"`
	EntitySeq  uint64   `json:"entity_seq,omitempty"`
	Simulated  bool     `json:"simulated"`
}

// Output is one message to publish, in order.
type Output struct {
	Subject string
	MsgID   string // Nats-Msg-Id; deterministic for the same input
	Data    any
}

type entState struct {
	e        Entity
	seenWall time.Time
}

type feedState struct {
	lastWall time.Time
	enabled  bool
}

// Engine holds zones, entities, memberships, open alerts and feed freshness.
type Engine struct {
	mu     sync.Mutex
	cfg    Config
	zones  map[string]*zone.Zone
	idx    *geo.Index
	ents   map[string]*entState
	inside map[string]map[string]bool // zone id -> entity id -> inside
	known  map[string]Crossing        // last crossing per "zone|entity" (for exits without a record)
	alerts map[string]*Alert          // open alerts
	seqs   map[string]uint64          // last seq per alert id
	feeds  map[string]*feedState
	count  map[string]int64
}

// New returns an engine.
func New(cfg Config) *Engine {
	if cfg.Source == "" {
		cfg.Source = "omos-geofence"
	}
	if cfg.AssetTypes == nil {
		cfg.AssetTypes = DefaultAssetTypes()
	}
	if cfg.FeedSilence <= 0 {
		cfg.FeedSilence = 90 * time.Second
	}
	if cfg.StaleDefault <= 0 {
		cfg.StaleDefault = 60
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	return &Engine{cfg: cfg, zones: map[string]*zone.Zone{}, idx: geo.NewIndex(cfg.CellDeg),
		ents: map[string]*entState{}, inside: map[string]map[string]bool{}, known: map[string]Crossing{},
		alerts: map[string]*Alert{}, seqs: map[string]uint64{}, feeds: map[string]*feedState{}, count: map[string]int64{}}
}

func (g *Engine) nowMs() int64 { return g.cfg.Now().UnixMilli() }

// Counters returns a copy of the event counters (metrics).
func (g *Engine) Counters() map[string]int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make(map[string]int64, len(g.count)+3)
	for k, v := range g.count {
		out[k] = v
	}
	out["zones"], out["entities"], out["alerts_open"] = int64(len(g.zones)), int64(len(g.ents)), int64(len(g.alerts))
	return out
}

// ---- ids ----

var b32 = base32.StdEncoding.WithPadding(base32.NoPadding)

// token keeps an id usable inside another id; long or odd ids are hashed.
func token(s string) string {
	if len(s) <= 48 && zone.IDRe.MatchString(strings.ToLower(s)) && s == strings.ToLower(s) {
		return s
	}
	h := sha256.Sum256([]byte(s))
	return "h" + strings.ToLower(b32.EncodeToString(h[:]))[:24]
}

// GeofenceAlertID is the alert id for an entity and a zone.
func GeofenceAlertID(zoneID, entityID string) string {
	return "alr-gf-" + token(zoneID) + "-" + token(entityID)
}

// StaleAlertID is the alert id for an asset that stopped reporting.
func StaleAlertID(entityID string) string { return "alr-stale-" + token(entityID) }

// FeedAlertID is the alert id for a feed outage.
func FeedAlertID(source string) string { return "alr-feed-" + token(source) }

// ---- zones ----

// applies reports whether a zone watches an entity.
func (g *Engine) applies(z *zone.Zone, e Entity) bool {
	if e.EntityType == "person" || (e.EntityType == "track" && e.Kind == "person") {
		return false // civilian stop line: zones never watch people
	}
	if nonMoving[e.EntityType] {
		return false
	}
	if len(z.Properties.AppliesTo) > 0 {
		for _, id := range z.Properties.AppliesTo {
			if id == e.EntityID {
				return true
			}
		}
		return false
	}
	if len(z.WatchTypes) > 0 {
		for _, t := range z.WatchTypes {
			if t == e.EntityType {
				return true
			}
		}
		return false
	}
	if z.Geo.Kind == zone.KindKeepIn {
		return g.cfg.AssetTypes[e.EntityType]
	}
	return true
}

func watched(z *zone.Zone) bool { return z != nil && z.Active && z.Watch }

// violation is whether being inside (or not) a zone is something to alert on.
func violation(kind string, in bool) bool {
	if kind == zone.KindKeepIn {
		return !in
	}
	return in
}

func alertLevel(kind string) string {
	switch kind {
	case zone.KindKeepIn, zone.KindKeepOut:
		return LevelWarning
	case zone.KindHazard:
		return LevelCaution
	}
	return LevelAdvisory
}

func alertReason(kind string) string {
	switch kind {
	case zone.KindKeepIn:
		return "left_keep_in"
	case zone.KindKeepOut:
		return "entered_keep_out"
	case zone.KindHazard:
		return "entered_hazard"
	}
	return "entered_zone"
}

// SetZone adds or replaces a zone and re-evaluates every entity against it.
func (g *Engine) SetZone(z *zone.Zone) []Output {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.zones[z.ZoneID] = z
	g.idx.Put(z.ZoneID, z.Area)
	g.count["zone_set"]++
	return g.reevaluateZone(z.ZoneID, CauseZoneChange)
}

// DeleteZone removes a zone: entities inside it exit (cause zone_removed) and its alerts resolve.
func (g *Engine) DeleteZone(id string) []Output {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.zones[id]; !ok && len(g.inside[id]) == 0 {
		return g.resolveZoneAlerts(id, "zone_removed")
	}
	delete(g.zones, id)
	g.idx.Remove(id)
	g.count["zone_deleted"]++
	return g.reevaluateZone(id, CauseZoneRemoved)
}

// reevaluateZone re-runs one zone for every entity that it may affect.
func (g *Engine) reevaluateZone(id, cause string) []Output {
	ids := make([]string, 0, len(g.ents)+len(g.inside[id]))
	seen := map[string]bool{}
	for eid := range g.ents {
		ids, seen[eid] = append(ids, eid), true
	}
	for eid := range g.inside[id] {
		if !seen[eid] {
			ids = append(ids, eid)
		}
	}
	sort.Strings(ids)
	var out []Output
	for _, eid := range ids {
		out = append(out, g.evalZoneEntity(id, eid, cause)...)
	}
	return append(out, g.resolveZoneAlertsIfGone(id)...)
}

func (g *Engine) resolveZoneAlertsIfGone(id string) []Output {
	if z, ok := g.zones[id]; ok && watched(z) {
		return nil
	}
	reason := "zone_inactive"
	if _, ok := g.zones[id]; !ok {
		reason = "zone_removed"
	}
	return g.resolveZoneAlerts(id, reason)
}

func (g *Engine) resolveZoneAlerts(id, reason string) []Output {
	var ids []string
	for aid, a := range g.alerts {
		if a.Kind == KindGeofence && a.ZoneID == id {
			ids = append(ids, aid)
		}
	}
	sort.Strings(ids)
	var out []Output
	for _, aid := range ids {
		out = append(out, g.closeAlert(aid, reason)...)
	}
	return out
}

// evalZoneEntity evaluates one (zone, entity) pair and returns the crossing and alert facts.
func (g *Engine) evalZoneEntity(zid, eid, cause string) []Output {
	z := g.zones[zid]
	es := g.ents[eid]
	if es == nil && watched(z) && cause != CauseEntityGone {
		// Not read from ENT yet (zones load first on start): keep what the replay says until
		// the entity arrives or ReconcileEntities finds it gone.
		return nil
	}
	prev := g.inside[zid][eid]
	present := es != nil && es.e.Live && es.e.HasPos && es.e.Lifecycle != "expired"
	in, applies := false, false
	if es != nil && watched(z) && g.applies(z, es.e) {
		applies = true
		if present {
			alt, hasAlt := 0.0, es.e.HAE != nil
			if hasAlt {
				alt = *es.e.HAE
			}
			in = z.Area.Contains(es.e.Lat, es.e.Lon, alt, hasAlt)
		}
	}
	var out []Output
	if in != prev {
		c := cause
		if c == "" {
			switch {
			case !in && !present:
				c = CauseEntityGone
			default:
				c = CauseCrossing
			}
		}
		out = append(out, g.crossing(zid, z, eid, es, in, c))
	}
	// Alert: open while a watched zone that applies is violated by a present entity.
	aid := GeofenceAlertID(zid, eid)
	if applies && present && violation(z.Geo.Kind, in) {
		e := es.e
		hae := e.HAE
		name := z.Name
		if name == "" {
			name = zid
		}
		label := e.Callsign
		if label == "" {
			label = e.EntityID
		}
		verb := map[bool]string{true: "is inside", false: "is outside"}[in]
		out = append(out, g.openAlert(Alert{AlertID: aid, Kind: KindGeofence, Level: alertLevel(z.Geo.Kind),
			About: About{EntityType: e.EntityType, EntityID: e.EntityID}, ZoneID: zid, ZoneKind: z.Geo.Kind,
			Reason:   alertReason(z.Geo.Kind),
			Message:  fmt.Sprintf("%s %s %s zone %s", label, verb, strings.ReplaceAll(z.Geo.Kind, "_", "-"), name),
			Location: &LatLon{Lat: e.Lat, Lon: e.Lon, HAE: hae}, Simulated: e.Simulated})...)
	} else if _, open := g.alerts[aid]; open {
		reason := "returned"
		switch {
		case z == nil:
			reason = "zone_removed"
		case !watched(z):
			reason = "zone_inactive"
		case !present:
			reason = "entity_gone"
		case !applies:
			reason = "zone_change"
		}
		out = append(out, g.closeAlert(aid, reason)...)
	}
	return out
}

func (g *Engine) crossing(zid string, z *zone.Zone, eid string, es *entState, in bool, cause string) Output {
	ev := "exit"
	if in {
		ev = "enter"
		if g.inside[zid] == nil {
			g.inside[zid] = map[string]bool{}
		}
		g.inside[zid][eid] = true
	} else {
		delete(g.inside[zid], eid)
		if len(g.inside[zid]) == 0 {
			delete(g.inside, zid)
		}
	}
	c := Crossing{V: 1, ZoneID: zid, Event: ev, Cause: cause, EntityID: eid}
	if z != nil {
		c.ZoneKind, c.ZoneName, c.ZoneRev = z.Geo.Kind, z.Name, z.Rev
	} else if k, ok := g.known[zid+"|"+eid]; ok {
		c.ZoneKind, c.ZoneName = k.ZoneKind, k.ZoneName
	}
	if es != nil {
		e := es.e
		c.EntityType, c.Callsign, c.EntityT, c.EntitySeq, c.Simulated = e.EntityType, e.Callsign, e.T, e.Seq, e.Simulated
		if e.HasPos {
			lat, lon := e.Lat, e.Lon
			c.Lat, c.Lon, c.HAE = &lat, &lon, e.HAE
		}
	} else if k, ok := g.known[zid+"|"+eid]; ok {
		c.EntityType, c.Callsign, c.Simulated = k.EntityType, k.Callsign, k.Simulated
	}
	g.known[zid+"|"+eid] = c
	if !in {
		delete(g.known, zid+"|"+eid)
	}
	g.count["crossing_"+ev]++
	return Output{Subject: "evt.map.geofence." + zid + "." + ev,
		MsgID: fmt.Sprintf("gf:%s:%s:%s:%d:%d:%s", zid, eid, ev, c.EntitySeq, c.ZoneRev, cause), Data: c}
}

// ---- alerts ----

func (g *Engine) openAlert(a Alert) []Output {
	now := g.nowMs()
	cur := g.alerts[a.AlertID]
	if cur != nil {
		if cur.Level == a.Level && cur.Reason == a.Reason {
			cur.Location = a.Location // follow the entity silently
			return nil
		}
		cur.Level, cur.Reason, cur.Message, cur.Location, cur.UpdatedT = a.Level, a.Reason, a.Message, a.Location, now
		g.seqs[a.AlertID]++
		cur.Seq = g.seqs[a.AlertID]
		g.count["alert_updated"]++
		return []Output{g.alertOut(cur, "updated")}
	}
	a.V, a.State, a.RaisedBy, a.OpenedT, a.UpdatedT = 1, StateOpen, g.cfg.Source, now, now
	g.seqs[a.AlertID]++
	a.Seq = g.seqs[a.AlertID]
	g.alerts[a.AlertID] = &a
	g.count["alert_raised"]++
	return []Output{g.alertOut(&a, "raised")}
}

func (g *Engine) closeAlert(id, reason string) []Output {
	cur := g.alerts[id]
	if cur == nil {
		return nil
	}
	now := g.nowMs()
	delete(g.alerts, id)
	cur.State, cur.Reason, cur.UpdatedT, cur.ResolvedT = StateResolved, reason, now, &now
	g.seqs[id]++
	cur.Seq = g.seqs[id]
	g.count["alert_resolved"]++
	return []Output{g.alertOut(cur, "resolved")}
}

func (g *Engine) alertOut(a *Alert, ev string) Output {
	cp := *a
	return Output{Subject: "evt.so.alert." + ev, MsgID: fmt.Sprintf("%s:%d", a.AlertID, a.Seq), Data: cp}
}

// OpenAlerts returns the open alerts sorted by id.
func (g *Engine) OpenAlerts() []Alert {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]Alert, 0, len(g.alerts))
	for _, a := range g.alerts {
		out = append(out, *a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].AlertID < out[j].AlertID })
	return out
}

// ---- entities ----

// SetEntity applies one `ent` record.
func (g *Engine) SetEntity(e Entity) []Output {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.count["entity_records"]++
	es := g.ents[e.EntityID]
	if es == nil {
		es = &entState{}
		g.ents[e.EntityID] = es
	}
	es.e, es.seenWall = e, g.cfg.Now()
	out := g.evalEntity(e.EntityID, "")
	out = append(out, g.staleEval(es)...)
	if !e.Live {
		delete(g.ents, e.EntityID)
	}
	return out
}

// RemoveEntity handles an `ENT` delete marker (the record expired by TTL).
func (g *Engine) RemoveEntity(entityID string) []Output {
	g.mu.Lock()
	defer g.mu.Unlock()
	es := g.ents[entityID]
	if es == nil {
		return g.evalEntity(entityID, CauseEntityGone)
	}
	es.e.Live, es.e.Lifecycle = false, "expired"
	out := g.evalEntity(entityID, "")
	out = append(out, g.staleEval(es)...)
	delete(g.ents, entityID)
	return out
}

// evalEntity runs every zone that may matter for one entity: zones containing it now, zones
// it was inside, keep-in zones (which alert when outside) and zones with an open alert.
func (g *Engine) evalEntity(eid, cause string) []Output {
	cand := map[string]bool{}
	if es := g.ents[eid]; es != nil && es.e.HasPos {
		alt, hasAlt := 0.0, es.e.HAE != nil
		if hasAlt {
			alt = *es.e.HAE
		}
		for _, id := range g.idx.Containing(es.e.Lat, es.e.Lon, alt, hasAlt) {
			cand[id] = true
		}
	}
	for zid, m := range g.inside {
		if m[eid] {
			cand[zid] = true
		}
	}
	for zid, z := range g.zones {
		if z.Geo.Kind == zone.KindKeepIn && watched(z) {
			cand[zid] = true
		}
	}
	for _, a := range g.alerts {
		if a.Kind == KindGeofence && a.About.EntityID == eid {
			cand[a.ZoneID] = true
		}
	}
	ids := make([]string, 0, len(cand))
	for id := range cand {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var out []Output
	for _, zid := range ids {
		out = append(out, g.evalZoneEntity(zid, eid, cause)...)
	}
	return out
}

func (g *Engine) staleLimit(e Entity) time.Duration {
	s := e.Stale
	if s <= 0 {
		s = g.cfg.StaleDefault
	}
	// The entity service republishes an active asset at least every min(stale/2, 120 s);
	// no record for 1.5 x stale + 5 s means the record itself is not being refreshed.
	return time.Duration(float64(s)*1.5*float64(time.Second)) + 5*time.Second
}

// staleEval raises, updates or resolves the stale alert of an own asset from its lifecycle.
func (g *Engine) staleEval(es *entState) []Output {
	e := es.e
	if !g.cfg.AssetTypes[e.EntityType] {
		return nil
	}
	lc := e.Lifecycle
	if !e.Live {
		lc = "expired"
	}
	switch lc {
	case "stale":
		return g.openStale(e, LevelCaution, "stale", "has not reported within its stale time")
	case "lost":
		return g.openStale(e, LevelWarning, "lost", "is lost (no report for three stale periods)")
	case "expired":
		return g.openStale(e, LevelWarning, "expired", "expired: its last record is gone")
	}
	return g.closeAlert(StaleAlertID(e.EntityID), "fresh")
}

func (g *Engine) openStale(e Entity, level, reason, what string) []Output {
	label := e.Callsign
	if label == "" {
		label = e.EntityID
	}
	var loc *LatLon
	if e.HasPos {
		loc = &LatLon{Lat: e.Lat, Lon: e.Lon, HAE: e.HAE}
	}
	return g.openAlert(Alert{AlertID: StaleAlertID(e.EntityID), Kind: KindStaleAsset, Level: level,
		About: About{EntityType: e.EntityType, EntityID: e.EntityID}, Reason: reason,
		Message: fmt.Sprintf("%s %s %s", e.EntityType, label, what), Location: loc, Simulated: e.Simulated})
}

// ---- feeds ----

// FeedStatus applies one `tel.feed.{source}.status` message. It reads the omos-ingest form
// (PR #553: `stale`, `fetchedAt`, `staleAfterSec`, `lastError`), the Redpanda Connect form
// (PR #558: same names) and the legacy feed-intel form (`feed_id`, `enabled`,
// `poll_interval`, `last_poll_ts`, `last_error`).
func (g *Engine) FeedStatus(source string, b []byte) []Output {
	var s struct {
		Stale         *bool    `json:"stale"`
		Enabled       *bool    `json:"enabled"`
		LastError     *string  `json:"lastError"`
		LastErrorPy   *string  `json:"last_error"`
		FetchedAt     *string  `json:"fetchedAt"`
		StaleAfterSec *float64 `json:"staleAfterSec"`
		PollInterval  *float64 `json:"poll_interval"`
		LastPollTS    any      `json:"last_poll_ts"`
	}
	if err := json.Unmarshal(b, &s); err != nil {
		g.mu.Lock()
		g.count["feed_status_unreadable"]++
		g.mu.Unlock()
		return nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	now := g.cfg.Now()
	g.count["feed_status"]++
	f := g.feeds[source]
	if f == nil {
		f = &feedState{}
		g.feeds[source] = f
	}
	f.lastWall = now
	f.enabled = s.Enabled == nil || *s.Enabled
	id := FeedAlertID(source)
	if !f.enabled {
		return g.closeAlert(id, "disabled")
	}
	lastErr := ""
	switch {
	case s.LastError != nil:
		lastErr = *s.LastError
	case s.LastErrorPy != nil:
		lastErr = *s.LastErrorPy
	}
	stale := false
	switch {
	case s.Stale != nil:
		stale = *s.Stale
	case s.LastPollTS != nil && s.PollInterval != nil:
		if t, ok := parseTS(s.LastPollTS); ok {
			limit := math.Max(3**s.PollInterval, 60)
			stale = now.Sub(t).Seconds() > limit
		}
	case s.FetchedAt != nil && s.StaleAfterSec != nil:
		if t, err := time.Parse(time.RFC3339Nano, *s.FetchedAt); err == nil {
			stale = now.Sub(t).Seconds() > *s.StaleAfterSec
		}
	}
	if !stale {
		return g.closeAlert(id, "fresh")
	}
	msg := fmt.Sprintf("feed %s is stale", source)
	if lastErr != "" {
		msg += ": " + lastErr
	}
	return g.openAlert(Alert{AlertID: id, Kind: KindFeedOutage, Level: LevelWarning, About: About{Feed: source},
		Reason: "stale", Message: msg})
}

func parseTS(v any) (time.Time, bool) {
	switch t := v.(type) {
	case float64:
		if t > 1e12 { // milliseconds
			return time.UnixMilli(int64(t)), true
		}
		return time.Unix(int64(t), int64((t-math.Floor(t))*1e9)), true
	case string:
		for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.999999", "2006-01-02 15:04:05"} {
			if p, err := time.Parse(layout, t); err == nil {
				return p, true
			}
		}
	}
	return time.Time{}, false
}

// ---- ticks ----

// Tick raises alerts that come from silence: assets with no new record and feeds with no
// status within their limits.
func (g *Engine) Tick() []Output {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := g.cfg.Now()
	var out []Output
	eids := make([]string, 0, len(g.ents))
	for id := range g.ents {
		eids = append(eids, id)
	}
	sort.Strings(eids)
	for _, id := range eids {
		es := g.ents[id]
		if !g.cfg.AssetTypes[es.e.EntityType] || es.e.Lifecycle == "stale" || es.e.Lifecycle == "lost" {
			continue
		}
		if now.Sub(es.seenWall) > g.staleLimit(es.e) {
			out = append(out, g.openStale(es.e, LevelCaution, "no_update", "has no new entity record")...)
		}
	}
	fids := make([]string, 0, len(g.feeds))
	for id := range g.feeds {
		fids = append(fids, id)
	}
	sort.Strings(fids)
	for _, id := range fids {
		f := g.feeds[id]
		if !f.enabled || now.Sub(f.lastWall) <= g.cfg.FeedSilence {
			continue
		}
		out = append(out, g.openAlert(Alert{AlertID: FeedAlertID(id), Kind: KindFeedOutage, Level: LevelWarning,
			About: About{Feed: id}, Reason: "silent",
			Message: fmt.Sprintf("feed %s sent no status for %s", id, now.Sub(f.lastWall).Round(time.Second))})...)
	}
	return out
}

// ---- restore ----

// RestoreCrossing replays one `evt.map.geofence.{zone}.{enter|exit}` fact.
func (g *Engine) RestoreCrossing(c Crossing) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if c.ZoneID == "" || c.EntityID == "" {
		return
	}
	key := c.ZoneID + "|" + c.EntityID
	if c.Event == "enter" {
		if g.inside[c.ZoneID] == nil {
			g.inside[c.ZoneID] = map[string]bool{}
		}
		g.inside[c.ZoneID][c.EntityID] = true
		g.known[key] = c
	} else {
		delete(g.inside[c.ZoneID], c.EntityID)
		if len(g.inside[c.ZoneID]) == 0 {
			delete(g.inside, c.ZoneID)
		}
		delete(g.known, key)
	}
	g.count["restored_crossings"]++
}

// RestoreAlert replays one `evt.so.alert.*` fact.
func (g *Engine) RestoreAlert(a Alert) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if a.AlertID == "" {
		return
	}
	if a.Seq > g.seqs[a.AlertID] {
		g.seqs[a.AlertID] = a.Seq
	} else if a.Seq < g.seqs[a.AlertID] {
		return // older than what we have
	}
	if a.State == StateResolved {
		delete(g.alerts, a.AlertID)
	} else {
		cp := a
		g.alerts[a.AlertID] = &cp
	}
	g.count["restored_alerts"]++
}

// Restored finishes the alert replay: feeds with an open outage alert are tracked again, so
// silence keeps them open and a fresh status resolves them.
func (g *Engine) Restored() {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := g.cfg.Now()
	for _, a := range g.alerts {
		if a.Kind == KindFeedOutage && a.About.Feed != "" && g.feeds[a.About.Feed] == nil {
			g.feeds[a.About.Feed] = &feedState{lastWall: now, enabled: true}
		}
	}
}

// ReconcileZones runs after the first full read of KV `geofences`: memberships and alerts of
// zones deleted while the service was down are closed.
func (g *Engine) ReconcileZones() []Output {
	g.mu.Lock()
	defer g.mu.Unlock()
	gone := map[string]bool{}
	for zid := range g.inside {
		if _, ok := g.zones[zid]; !ok {
			gone[zid] = true
		}
	}
	for _, a := range g.alerts {
		if a.Kind == KindGeofence {
			if _, ok := g.zones[a.ZoneID]; !ok {
				gone[a.ZoneID] = true
			}
		}
	}
	ids := make([]string, 0, len(gone))
	for id := range gone {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var out []Output
	for _, id := range ids {
		out = append(out, g.reevaluateZone(id, CauseZoneRemoved)...)
	}
	return out
}

// ReconcileEntities runs after the first full read of `ENT`: memberships and geofence alerts
// of entities that expired while the service was down are closed, and the stale alert of an
// asset that is gone becomes `expired`.
func (g *Engine) ReconcileEntities() []Output {
	g.mu.Lock()
	defer g.mu.Unlock()
	gone := map[string]bool{}
	for _, m := range g.inside {
		for eid := range m {
			if g.ents[eid] == nil {
				gone[eid] = true
			}
		}
	}
	var staleGone []*Alert
	for _, a := range g.alerts {
		if a.About.EntityID == "" || g.ents[a.About.EntityID] != nil {
			continue
		}
		switch a.Kind {
		case KindGeofence:
			gone[a.About.EntityID] = true
		case KindStaleAsset:
			staleGone = append(staleGone, a)
		}
	}
	ids := make([]string, 0, len(gone))
	for id := range gone {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var out []Output
	for _, eid := range ids {
		out = append(out, g.evalEntity(eid, CauseEntityGone)...)
	}
	sort.Slice(staleGone, func(i, j int) bool { return staleGone[i].AlertID < staleGone[j].AlertID })
	for _, a := range staleGone {
		e := Entity{EntityID: a.About.EntityID, EntityType: a.About.EntityType, Simulated: a.Simulated}
		if a.Location != nil {
			e.Lat, e.Lon, e.HAE, e.HasPos = a.Location.Lat, a.Location.Lon, a.Location.HAE, true
		}
		out = append(out, g.openStale(e, LevelWarning, "expired", "expired: its last record is gone")...)
	}
	return out
}
