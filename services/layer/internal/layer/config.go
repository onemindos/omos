// Package layer is the OneMind layer service (FEATURES A5): it follows the ENT stream, keeps
// the latest render state of every entity in memory, writes layer manifests and density cells
// to KV `world-state` and site status to KV `fleet-state`, and serves `svc.layer.list`,
// `svc.layer.meta`, `svc.layer.snapshot`, `svc.layer.watch` and `svc.layer.unwatch`.
//
// Contracts: onemind-xr docs/contracts/layer-snapshot.md and layer-delta.md (v1), BUS-DESIGN
// D3 (bucket purposes) and D7 (clients read layers through svc.layer.*).
package layer

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Tiers (layer-snapshot.md §1).
const (
	TierSmall  = "small"
	TierMedium = "medium"
	TierLarge  = "large"
	TierStatic = "static"
)

// CellsCfg says which density-cell zooms a layer keeps and below which camera zoom a
// snapshot or watch returns cells instead of entities.
type CellsCfg struct {
	Zooms     []int `json:"zooms"`
	BelowZoom int   `json:"below_zoom"`
}

// ObjectRef points at a static layer's object in the Object Store (layer-snapshot.md §3.4).
type ObjectRef struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
	SHA256 string `json:"sha256"`
	Bytes  int64  `json:"bytes"`
}

// Def is one layer definition.
type Def struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Tier        string   `json:"tier"`
	EntityTypes []string `json:"entity_types"`
	// CotOnlyTypes are entity types that belong to this layer only when one of the entity's
	// sources has data_type "cot" (TAK shapes appear as `zone` in the `tak` layer).
	CotOnlyTypes []string `json:"cot_only_types,omitempty"`
	// Kinds and SourceSubjects narrow the entity types: when either is set, an entity belongs
	// only if its `kind` is in Kinds or one of its sources' `subject` starts with a prefix in
	// SourceSubjects. The entity service puts feed data under shared ONTOLOGY types (NWS
	// alerts and quakes are `event`, satellites and trains are `vehicle`), so the subject
	// type alone cannot tell a hazard from a conflict report or a train from a fleet rover.
	Kinds          []string `json:"kinds,omitempty"`
	SourceSubjects []string `json:"source_subjects,omitempty"`
	// ExcludeKinds drops entities of these kinds even when the type matches.
	ExcludeKinds  []string   `json:"exclude_kinds,omitempty"`
	Cap           int        `json:"cap"`
	Cells         *CellsCfg  `json:"cells,omitempty"`
	LimitS        int        `json:"limit_s"`
	Feeds         []string   `json:"feeds"`
	DefaultStaleS int        `json:"default_stale_s"`
	Marking       string     `json:"marking"`
	Style         string     `json:"style"`
	Object        *ObjectRef `json:"object,omitempty"`
}

var (
	layerIDRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,31}$`)
	tokenRe   = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)
	// A source subject prefix is one or more literal tokens, each ending in a dot.
	subjPrefixRe = regexp.MustCompile(`^([A-Za-z0-9_-]{1,64}\.){1,8}$`)
)

// Filters returns the ENT subject filters of the layer (one per entity type).
func (d *Def) Filters() []string {
	var out []string
	for _, t := range d.EntityTypes {
		out = append(out, "ent."+t+".>")
	}
	for _, t := range d.CotOnlyTypes {
		out = append(out, "ent."+t+".>")
	}
	return out
}

// Filter is the single `filter` string of the manifest: the one filter when the layer has
// one entity type, otherwise the whole ENT filter; `filters` lists the exact set.
func (d *Def) Filter() string {
	f := d.Filters()
	if len(f) == 1 {
		return f[0]
	}
	if len(f) == 0 {
		return ""
	}
	return "ent.*.*.state"
}

// HasType reports whether entity type t can be in the layer.
func (d *Def) HasType(t string) (member bool, cotOnly bool) {
	for _, x := range d.EntityTypes {
		if x == t {
			return true, false
		}
	}
	for _, x := range d.CotOnlyTypes {
		if x == t {
			return true, true
		}
	}
	return false, false
}

// Qualifies reports whether an entity of a member type, with this kind and these source
// subjects, passes the layer's kind and source-subject rules.
func (d *Def) Qualifies(kind string, subjects []string) bool {
	for _, k := range d.ExcludeKinds {
		if k == kind {
			return false
		}
	}
	if len(d.Kinds) == 0 && len(d.SourceSubjects) == 0 {
		return true
	}
	for _, k := range d.Kinds {
		if k == kind {
			return true
		}
	}
	for _, p := range d.SourceSubjects {
		for _, s := range subjects {
			if strings.HasPrefix(s, p) {
				return true
			}
		}
	}
	return false
}

// CellZoom picks the cell zoom for a camera zoom: the largest configured zoom not above it,
// else the smallest configured zoom.
func (d *Def) CellZoom(zoom int, zooms []int) int {
	best, min := -1, -1
	for _, z := range zooms {
		if min < 0 || z < min {
			min = z
		}
		if z <= zoom && z > best {
			best = z
		}
	}
	if best >= 0 {
		return best
	}
	return min
}

// Validate checks a definition.
func (d *Def) Validate() error {
	if !layerIDRe.MatchString(d.ID) {
		return fmt.Errorf("layer id %q must match %s", d.ID, layerIDRe)
	}
	switch d.Tier {
	case TierSmall, TierMedium, TierLarge:
		if len(d.EntityTypes)+len(d.CotOnlyTypes) == 0 {
			return fmt.Errorf("layer %s: dynamic layer needs entity_types", d.ID)
		}
	case TierStatic:
	default:
		return fmt.Errorf("layer %s: tier %q", d.ID, d.Tier)
	}
	for _, t := range append(append([]string{}, d.EntityTypes...), d.CotOnlyTypes...) {
		if !tokenRe.MatchString(t) {
			return fmt.Errorf("layer %s: entity type %q is not a subject token", d.ID, t)
		}
	}
	for _, k := range append(append([]string{}, d.Kinds...), d.ExcludeKinds...) {
		if !tokenRe.MatchString(k) {
			return fmt.Errorf("layer %s: kind %q is not a token", d.ID, k)
		}
	}
	for _, p := range d.SourceSubjects {
		if !subjPrefixRe.MatchString(p) {
			return fmt.Errorf("layer %s: source subject prefix %q must be literal tokens ending in a dot", d.ID, p)
		}
	}
	if d.Cap <= 0 {
		return fmt.Errorf("layer %s: cap must be > 0", d.ID)
	}
	if d.Tier == TierLarge && (d.Cells == nil || len(d.Cells.Zooms) == 0) {
		return fmt.Errorf("layer %s: large layers need cells", d.ID)
	}
	if d.LimitS <= 0 {
		return fmt.Errorf("layer %s: limit_s must be > 0", d.ID)
	}
	for _, f := range d.Feeds {
		if !tokenRe.MatchString(f) {
			return fmt.Errorf("layer %s: feed %q is not a subject token", d.ID, f)
		}
	}
	return nil
}

// DefaultLayers is the layer table of layer-snapshot.md §1, mapped onto the ONTOLOGY types
// and kinds the entity service (omos #552) actually writes to ENT: NWS alerts, quakes and
// fires are `event` (kinds weather_alert, earthquake, wildfire, or the record's own type for
// other intel), satellites are `vehicle` kind `satellite` (norad-*), trains are `vehicle` kind
// `train` (gtfs-*). The contract's `hazard`, `satellite`, `train` and `station` types do not
// exist on ENT. Counts, TTLs and limits there are estimates; caps here are the contract's
// (large layers 5,000 per AOI request) and a whole-layer cap for small and medium layers.
// Feeds list only collectors that publish tel.feed.{source}.status today or in an open PR
// (omos-ingest: nws-alerts); a layer with no listed feed reports `fresh` while ENT is readable.
func DefaultLayers() []Def {
	cells := &CellsCfg{Zooms: []int{2, 4, 6}, BelowZoom: 7}
	return []Def{
		// Feed vehicles (satellites, trains) are `vehicle` too; they have their own layers.
		{ID: "fleet", Title: "Fleet", Tier: TierSmall, EntityTypes: []string{"drone", "robot", "vehicle", "node", "edge"}, ExcludeKinds: []string{"satellite", "train"}, Cap: 20000, LimitS: 60, DefaultStaleS: 60, Marking: "AMBER+STRICT", Style: "fleet.default.v1"},
		{ID: "tak", Title: "TAK", Tier: TierSmall, EntityTypes: []string{"person", "place"}, CotOnlyTypes: []string{"zone"}, Cap: 20000, LimitS: 300, DefaultStaleS: 300, Marking: "AMBER", Style: "tak.default.v1"},
		{ID: "zones", Title: "Zones", Tier: TierSmall, EntityTypes: []string{"zone"}, Cap: 20000, LimitS: 3600, DefaultStaleS: 0, Marking: "AMBER", Style: "zones.default.v1"},
		// Hazard events: by kind, or by the intel/world domain of a source subject (contract
		// §1 feeds). Other events (conflict, chokepoint) stay out.
		{ID: "hazards", Title: "Hazards", Tier: TierMedium, EntityTypes: []string{"event"},
			Kinds:          []string{"weather_alert", "earthquake", "wildfire", "flood", "storm", "volcano", "tsunami"},
			SourceSubjects: []string{"intel.weather.", "intel.seismic.", "intel.fire.", "intel.disaster.", "intel.natural.", "evt.world.fire.", "evt.world.seismic."},
			Cap:            20000, LimitS: 900, Feeds: []string{"nws-alerts"}, DefaultStaleS: 3600, Marking: "CLEAR", Style: "hazards.default.v1"},
		{ID: "tracks", Title: "Tracks", Tier: TierMedium, EntityTypes: []string{"track"}, Cap: 20000, LimitS: 30, DefaultStaleS: 30, Marking: "AMBER+STRICT", Style: "tracks.default.v1"},
		{ID: "space", Title: "Space", Tier: TierMedium, EntityTypes: []string{"vehicle"}, Kinds: []string{"satellite"}, Cap: 20000, LimitS: 300, DefaultStaleS: 120, Marking: "CLEAR", Style: "space.default.v1"},
		{ID: "transit", Title: "Transit", Tier: TierMedium, EntityTypes: []string{"vehicle"}, Kinds: []string{"train"}, Cap: 20000, LimitS: 300, DefaultStaleS: 120, Marking: "CLEAR", Style: "transit.default.v1"},
		// RF stations are ONTOLOGY `sensor` entities. No entity-service adapter emits APRS or
		// Meshtastic yet, so this layer stays empty until one does.
		{ID: "rf", Title: "RF stations", Tier: TierMedium, EntityTypes: []string{"sensor"}, Kinds: []string{"aprs", "meshtastic", "rf"}, Cap: 20000, LimitS: 3600, DefaultStaleS: 1800, Marking: "CLEAR", Style: "rf.default.v1"},
		{ID: "cameras", Title: "Cameras", Tier: TierSmall, EntityTypes: []string{"camera"}, Cap: 20000, LimitS: 300, DefaultStaleS: 300, Marking: "AMBER+STRICT", Style: "cameras.default.v1"},
		{ID: "air", Title: "Aircraft", Tier: TierLarge, EntityTypes: []string{"aircraft"}, Cap: 5000, Cells: cells, LimitS: 60, DefaultStaleS: 30, Marking: "CLEAR", Style: "air.default.v1"},
		{ID: "sea", Title: "Vessels", Tier: TierLarge, EntityTypes: []string{"vessel"}, Cap: 5000, Cells: &CellsCfg{Zooms: []int{2, 4, 6}, BelowZoom: 7}, LimitS: 1800, DefaultStaleS: 1800, Marking: "CLEAR", Style: "sea.default.v1"},
	}
}

// LoadLayers reads a JSON array of Def from path, or returns DefaultLayers when path is empty.
func LoadLayers(path string) ([]Def, error) {
	if path == "" {
		return DefaultLayers(), nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var defs []Def
	if err := json.Unmarshal(b, &defs); err != nil {
		return nil, fmt.Errorf("layers file: %w", err)
	}
	return defs, ValidateLayers(defs)
}

// ValidateLayers checks every definition and that ids are unique.
func ValidateLayers(defs []Def) error {
	seen := map[string]bool{}
	for i := range defs {
		if err := defs[i].Validate(); err != nil {
			return err
		}
		if seen[defs[i].ID] {
			return fmt.Errorf("duplicate layer id %q", defs[i].ID)
		}
		seen[defs[i].ID] = true
	}
	return nil
}
