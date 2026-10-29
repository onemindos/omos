package projector

// The gates of cot-projection.md §3 (G1 to G8; G9 budget is in engine.go) and the AOI index
// built from KV `geofences` and KV `mission-state`.

import (
	"encoding/json"
	"math"
	"strings"
	"time"
)

// Gate reasons (counted per update and logged with removals).
const (
	ReasonDisabled    = "g1_disabled"
	ReasonNotAlive    = "g2_not_alive"
	ReasonNoType      = "g3_no_type"
	ReasonPerson      = "g4_person"
	ReasonOnTAK       = "g5_already_on_tak"
	ReasonEgressOff   = "g6_egress_off"
	ReasonMarking     = "g7_marking"
	ReasonNotRelevant = "g8_not_relevant"
	ReasonCap         = "g9_cap"
	ReasonDeleted     = "deleted"
)

// Ranks for the cap order (§3.3); lower is kept first.
const (
	RankOwn     = 1
	RankTasked  = 2
	RankPinned  = 3
	RankStarred = 4
	RankAOI     = 5
)

// aoiKinds are the zone kinds that act as an AOI by themselves (§3.2).
var aoiKinds = map[string]bool{"aoi": true, "keep_in": true, "keep_out": true, "hazard": true, "control_area": true}

type zone struct {
	id     string
	kind   string
	active bool
	area   *area
}

// mission and task facts the relevance tier reads.
type missionFacts struct {
	activeAOI map[string]bool // zone ids that are the aoi of an active mission
	assigned  map[string]bool // "{entity_type}.{entity_id}" assigned to an active mission
	targeted  map[string]bool // entity ids that are a target of a non-terminal task
	tasked    map[string]bool // "{entity_type}.{entity_id}" assignee of a non-terminal task
}

type aoiIndex struct {
	zones    map[string]*zone
	missions map[string]missionRec // key -> parsed mission-state record
	facts    missionFacts
	active   []*area // recomputed on change
}

func newAOIIndex() *aoiIndex {
	x := &aoiIndex{zones: map[string]*zone{}, missions: map[string]missionRec{}}
	x.rebuild()
	return x
}

// geofenceRec is a KV `geofences` value: a zone entity record (entity.md §4.1, the mirror is
// `ent.zone.*`) or the minimal form {zone_id, geo, active}. Proposal A6-P2 asks the owner to
// fix this shape with D1.
type geofenceRec struct {
	ZoneID    string `json:"zone_id"`
	EntityID  string `json:"entity_id"`
	ID        string `json:"id"`
	Active    *bool  `json:"active"`
	Live      *bool  `json:"live"`
	Lifecycle string `json:"lifecycle"`
	Geo       *Geo   `json:"geo"`
}

// setZone parses a geofences value. ok is false for an unreadable value (counted, ignored).
func (x *aoiIndex) setZone(key string, data []byte) (ok bool) {
	var r geofenceRec
	if err := json.Unmarshal(data, &r); err != nil || r.Geo == nil {
		delete(x.zones, key)
		x.rebuild()
		return false
	}
	a, err := geoArea(r.Geo)
	if err != nil {
		delete(x.zones, key)
		x.rebuild()
		return false
	}
	active := (r.Active == nil || *r.Active) && (r.Live == nil || *r.Live)
	switch r.Lifecycle {
	case "expired", "lost", "deleted", "removed":
		active = false
	}
	id := key
	for _, c := range []string{r.ZoneID, r.EntityID, r.ID} {
		if c != "" {
			id = c
			break
		}
	}
	x.zones[key] = &zone{id: id, kind: strings.ToLower(r.Geo.Kind), active: active, area: a}
	x.rebuild()
	return true
}

func (x *aoiIndex) delZone(key string) {
	delete(x.zones, key)
	x.rebuild()
}

// missionRec is the part of a mission-state record (mission-task.md §3.1, §3.2) read here.
type missionRec struct {
	State       string `json:"state"`
	Rollup      string `json:"rollup"`
	AOI         string `json:"aoi"`
	Assignments []struct {
		To string `json:"to"`
	} `json:"assignments"`
	Assignee *struct {
		EntityType string `json:"entity_type"`
		EntityID   string `json:"entity_id"`
	} `json:"assignee"`
	TargetSnapshots []struct {
		EntityID string `json:"entity_id"`
	} `json:"target_snapshots"`
	Spec struct {
		Target *struct {
			EntityID string `json:"entity_id"`
		} `json:"target"`
		TargetID string `json:"target_id"`
	} `json:"spec"`
}

func (x *aoiIndex) setMission(key string, data []byte) bool {
	var r missionRec
	if err := json.Unmarshal(data, &r); err != nil {
		delete(x.missions, key)
		x.rebuild()
		return false
	}
	x.missions[key] = r
	x.rebuild()
	return true
}

func (x *aoiIndex) delMission(key string) {
	delete(x.missions, key)
	x.rebuild()
}

var terminalRollup = map[string]bool{"completed": true, "failed": true, "cancelled": true}

func (x *aoiIndex) rebuild() {
	f := missionFacts{activeAOI: map[string]bool{}, assigned: map[string]bool{}, targeted: map[string]bool{}, tasked: map[string]bool{}}
	for key, m := range x.missions {
		switch {
		case strings.HasPrefix(key, "msn-"):
			if m.State != "active" {
				continue
			}
			if m.AOI != "" {
				f.activeAOI[m.AOI] = true
			}
			for _, a := range m.Assignments {
				if a.To != "" {
					f.assigned[a.To] = true
				}
			}
		case strings.HasPrefix(key, "tsk-"):
			if terminalRollup[m.Rollup] {
				continue
			}
			if m.Assignee != nil && m.Assignee.EntityID != "" {
				f.tasked[m.Assignee.EntityType+"."+m.Assignee.EntityID] = true
			}
			for _, t := range m.TargetSnapshots {
				if t.EntityID != "" {
					f.targeted[t.EntityID] = true
				}
			}
			if m.Spec.Target != nil && m.Spec.Target.EntityID != "" {
				f.targeted[m.Spec.Target.EntityID] = true
			}
			if m.Spec.TargetID != "" {
				f.targeted[m.Spec.TargetID] = true
			}
		}
	}
	x.facts = f
	x.active = x.active[:0]
	for _, z := range x.zones {
		if z.active && (aoiKinds[z.kind] || f.activeAOI[z.id]) {
			x.active = append(x.active, z.area)
		}
	}
}

// inAOI reports whether the entity point, or any vertex of its geo, is inside an active AOI,
// and the distance in metres to the nearest containing AOI centroid (cap order).
func (x *aoiIndex) inAOI(e *Entity) (bool, float64) {
	if len(x.active) == 0 {
		return false, 0
	}
	alt, hasAlt := e.Altitude()
	pts := []pt{}
	if lat, lon, ok := e.Position(); ok {
		pts = append(pts, pt{lon, lat})
	}
	if e.Geo != nil {
		pts = append(pts, geoPoints(e.Geo)...)
	}
	best, hit := math.Inf(1), false
	for _, a := range x.active {
		for _, p := range pts {
			if a.contains(p[1], p[0], alt, hasAlt) {
				hit = true
				if d := haversineM(p[1], p[0], a.cLat, a.cLon); d < best {
					best = d
				}
				break
			}
		}
	}
	return hit, best
}

// decision is the result of the gates for one entity.
type decision struct {
	ok          bool
	reason      string // first failing gate
	cotType     string
	downgraded  bool
	rank        int
	dist        float64
	onlyOutside bool // failed only the relevance tier (hysteresis applies)
}

var aliveLifecycle = map[string]bool{"active": true, "stale": true}

var takInputSubjects = map[string]bool{"tak.cot.out": true, "tak.cot.mesh.rx": true}

// evaluate runs G1 to G8 in order.
func evaluate(e *Entity, cfg *Config, x *aoiIndex, enabled bool, now time.Time) decision {
	if !enabled {
		return decision{reason: ReasonDisabled}
	}
	// G2: alive and located.
	if !e.Live || !aliveLifecycle[e.Lifecycle] || !e.hasTrait("locatable") {
		return decision{reason: ReasonNotAlive}
	}
	if _, _, ok := e.Position(); !ok {
		return decision{reason: ReasonNotAlive}
	}
	// G4 before G3, so a person-class record is always reported as such (both stop it).
	if e.IsPersonClass() || e.EntityType == "station" {
		return decision{reason: ReasonPerson}
	}
	// G3: type allowed.
	t, down, why := CotType(e)
	if why != "" {
		return decision{reason: ReasonNoType}
	}
	// G5: not already on TAK.
	staleMs := int64(StaleSeconds(e, cfg.OtherInterval)) * 1000
	for _, s := range e.Sources {
		if takInputSubjects[s.Subject] && s.LastSeenT > now.UnixMilli()-staleMs {
			return decision{reason: ReasonOnTAK}
		}
	}
	// G6: egress allowed.
	if !e.Indicators.EgressTAK {
		return decision{reason: ReasonEgressOff}
	}
	// G7: marking (§3.4), fail closed on anything unknown.
	if !markingAllowed(e, cfg.ReleaseGroups) {
		return decision{reason: ReasonMarking}
	}
	// G8: relevance tier (§3.2; mesh profile §6.2).
	d := decision{ok: true, cotType: t, downgraded: down}
	key := e.Key()
	switch {
	case e.IsOwnAsset():
		d.rank = RankOwn
	case e.CurrentTask != nil || x.facts.targeted[e.EntityID] || x.facts.tasked[key]:
		d.rank = RankTasked
	case x.facts.assigned[key]:
		d.rank = RankTasked
	case e.IsPinnedType() && (cfg.Target == TargetServer || e.EntityType == "zone" || cfg.MeshWorldFeeds):
		d.rank = RankPinned
	case e.Indicators.Starred && (cfg.Target == TargetServer || cfg.MeshWorldFeeds):
		d.rank = RankStarred
	default:
		in, dist := x.inAOI(e)
		if in && cfg.Target == TargetMesh && !cfg.MeshWorldFeeds && !x.inMissionAOI(e) {
			in = false // mesh: no world feeds unless MESH_WORLD_FEEDS (§6.2)
		}
		if !in {
			return decision{reason: ReasonNotRelevant, onlyOutside: true, cotType: t}
		}
		d.rank, d.dist = RankAOI, dist
	}
	return d
}

// inMissionAOI is the mesh profile's "active-mission entities" test: inside the AOI zone of an
// active mission.
func (x *aoiIndex) inMissionAOI(e *Entity) bool {
	lat, lon, ok := e.Position()
	if !ok {
		return false
	}
	alt, hasAlt := e.Altitude()
	for _, z := range x.zones {
		if z.active && x.facts.activeAOI[z.id] && z.area.contains(lat, lon, alt, hasAlt) {
			return true
		}
	}
	return false
}

func markingAllowed(e *Entity, groups []string) bool {
	switch strings.ToUpper(strings.TrimSpace(e.Marking.TLP)) {
	case "CLEAR", "GREEN":
		return true
	case "AMBER", "AMBER+STRICT":
		for _, r := range e.Marking.ReleasableTo {
			for _, g := range groups {
				if r == g {
					return true
				}
			}
		}
	}
	return false // RED, missing or unknown: never
}
