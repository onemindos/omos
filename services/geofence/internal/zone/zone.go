// Package zone defines the KV `geofences` value this service writes and how a
// `cmd.map.geofence.{id}.{verb}` command turns into one.
//
// The value is the shape both existing readers already accept (proposal D1-P1):
//   - the TAK projector (PR #563, A6-P2): key = zone id; `zone_id`, `active`, `live`,
//     `lifecycle`, `geo{shape, kind, coords (GeoJSON), radius_m, alt_min_m, alt_max_m}`;
//   - the policy service (PR #554, POL-9): `entity_id`, `live`, `geo{shape polygon|circle,
//     coords, radius_m, kind}`, `properties.applies_to`.
//
// keep_in and keep_out zones are limited to polygon and circle, because the policy service
// fails closed on any other keep-in/keep-out shape (POL-9) and would block every movement.
package zone

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/onemindos/omos/services/geofence/internal/geo"
)

// Kinds (ONTOLOGY.md zone row: `geo.kind`).
const (
	KindAOI         = "aoi"
	KindKeepIn      = "keep_in"
	KindKeepOut     = "keep_out"
	KindHazard      = "hazard"
	KindControlArea = "control_area"
)

var kinds = map[string]bool{KindAOI: true, KindKeepIn: true, KindKeepOut: true, KindHazard: true, KindControlArea: true}

// IDRe is a zone or entity id that is one NATS token.
var IDRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,62}$`)

var entityIDRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_:-]{0,127}$`)
var typeRe = regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`)

// Refusal codes (asset-command.md §6.4 where one fits; the rest are this service's).
const (
	CodeSchemaInvalid    = "schema_invalid"
	CodeUnsupportedShape = "unsupported_shape"
	CodeForbidden        = "forbidden_action"
	CodeExists           = "exists"
	CodeNotFound         = "not_found"
	CodeConflict         = "conflict"
	CodeUnsigned         = "unsigned"
	CodeBadSignature     = "bad_signature"
	CodeExpired          = "expired"
	CodeNotApproved      = "not_approved"
	CodeBadSubject       = "bad_subject"
	CodeUnavailable      = "unavailable"
)

// Refusal is why a command was not applied.
type Refusal struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (r *Refusal) Error() string { return r.Code + ": " + r.Message }

func refuse(code, format string, a ...any) *Refusal {
	return &Refusal{Code: code, Message: fmt.Sprintf(format, a...)}
}

// Geo is the record's `geo` component.
type Geo struct {
	Shape   string          `json:"shape"`
	Kind    string          `json:"kind"`
	Coords  json.RawMessage `json:"coords"`
	RadiusM float64         `json:"radius_m,omitempty"`
	AltMinM *float64        `json:"alt_min_m,omitempty"`
	AltMaxM *float64        `json:"alt_max_m,omitempty"`
}

// Properties carries the policy service's `applies_to`.
type Properties struct {
	AppliesTo []string `json:"applies_to,omitempty"`
}

// Record is one KV `geofences` value (services/schemas/map/geofence.v1.schema.json).
type Record struct {
	V          int        `json:"v"`
	ZoneID     string     `json:"zone_id"`
	EntityID   string     `json:"entity_id"`
	EntityType string     `json:"entity_type"`
	Name       string     `json:"name,omitempty"`
	Active     bool       `json:"active"`
	Live       bool       `json:"live"`
	Lifecycle  string     `json:"lifecycle"`
	Pinned     bool       `json:"pinned"`
	Geo        Geo        `json:"geo"`
	Watch      bool       `json:"watch"`
	WatchTypes []string   `json:"watch_types,omitempty"`
	Properties Properties `json:"properties"`
	CreatedT   int64      `json:"created_t"`
	UpdatedT   int64      `json:"updated_t"`
	CreatedBy  string     `json:"created_by,omitempty"`
	UpdatedBy  string     `json:"updated_by,omitempty"`
	ActionID   string     `json:"action_id,omitempty"`
	CommandID  string     `json:"command_id,omitempty"`
	Simulated  bool       `json:"simulated"`
}

// Params is the command `data` (the policy action parameters of zone.create / zone.update;
// services/schemas/actions/zone.*.v1.schema.json). `polygon` + `kind` is the doctrine §9.7
// shorthand (`cmd.map.geofence.{id}.create {polygon: [[lon,lat],...], watch: true}`).
type Params struct {
	ZoneID     string      `json:"zone_id,omitempty"`
	Name       *string     `json:"name,omitempty"`
	Geo        *Geo        `json:"geo,omitempty"`
	Polygon    [][]float64 `json:"polygon,omitempty"`
	Kind       string      `json:"kind,omitempty"`
	Watch      *bool       `json:"watch,omitempty"`
	WatchTypes *[]string   `json:"watch_types,omitempty"`
	AppliesTo  *[]string   `json:"applies_to,omitempty"`
	Active     *bool       `json:"active,omitempty"`
	Revision   *uint64     `json:"revision,omitempty"`
	Reason     string      `json:"reason,omitempty"`
}

// DefaultWatch is whether crossings of a zone kind are watched when the writer does not say:
// keep-in, keep-out and hazard zones are watched; AOIs and control areas only filter what TAK
// receives (A6) and can be large, so they are not.
func DefaultWatch(kind string) bool {
	return kind == KindKeepIn || kind == KindKeepOut || kind == KindHazard
}

// Meta is what the command envelope says about who and why.
type Meta struct {
	Actor     string
	ActionID  string
	CommandID string
	Simulated bool
	NowMs     int64
}

// Build makes the record for create (prev == nil) or update (prev != nil).
func Build(id string, prev *Record, p Params, m Meta) (Record, *Refusal) {
	if !IDRe.MatchString(id) {
		return Record{}, refuse(CodeBadSubject, "zone id %q is not a lowercase NATS token", id)
	}
	if p.ZoneID != "" && p.ZoneID != id {
		return Record{}, refuse(CodeSchemaInvalid, "data.zone_id %q does not match the subject id %q", p.ZoneID, id)
	}
	var r Record
	if prev != nil {
		r = *prev
		r.WatchTypes = append([]string(nil), prev.WatchTypes...)
		r.Properties.AppliesTo = append([]string(nil), prev.Properties.AppliesTo...)
	} else {
		r = Record{V: 1, ZoneID: id, EntityID: id, EntityType: "zone", Active: true, Pinned: true,
			CreatedT: m.NowMs, CreatedBy: m.Actor}
	}
	g := p.Geo
	if g == nil && len(p.Polygon) > 0 {
		coords, _ := json.Marshal(map[string]any{"type": "Polygon", "coordinates": [][][]float64{p.Polygon}})
		g = &Geo{Shape: geo.ShapePolygon, Kind: p.Kind, Coords: coords}
	}
	if g != nil {
		ng := *g
		if ng.Kind == "" && prev != nil {
			ng.Kind = prev.Geo.Kind
		}
		r.Geo = ng
	} else if prev == nil {
		return Record{}, refuse(CodeSchemaInvalid, "create needs geo (or polygon and kind)")
	} else if p.Kind != "" {
		r.Geo.Kind = p.Kind
	}
	if !kinds[r.Geo.Kind] {
		return Record{}, refuse(CodeSchemaInvalid, "geo.kind %q is not one of aoi, keep_in, keep_out, hazard, control_area", r.Geo.Kind)
	}
	if (r.Geo.Kind == KindKeepIn || r.Geo.Kind == KindKeepOut) &&
		r.Geo.Shape != geo.ShapePolygon && r.Geo.Shape != geo.ShapeCircle {
		return Record{}, refuse(CodeUnsupportedShape,
			"%s zones must be polygon or circle (the policy service cannot evaluate %q and would deny every movement)", r.Geo.Kind, r.Geo.Shape)
	}
	if _, err := geo.Parse(r.Geo.Shape, r.Geo.Coords, r.Geo.RadiusM, r.Geo.AltMinM, r.Geo.AltMaxM); err != nil {
		return Record{}, refuse(CodeSchemaInvalid, "geo: %v", err)
	}
	if r.Geo.Shape == geo.ShapePolygon {
		r.Geo.Coords = closePolygon(r.Geo.Coords)
	}
	if r.Geo.Shape != geo.ShapeCircle && r.Geo.Shape != geo.ShapeCorridor {
		r.Geo.RadiusM = 0
	}
	if p.Name != nil {
		if len(*p.Name) > 120 {
			return Record{}, refuse(CodeSchemaInvalid, "name is longer than 120 characters")
		}
		r.Name = *p.Name
	}
	if p.Active != nil {
		r.Active = *p.Active
	}
	switch {
	case p.Watch != nil:
		r.Watch = *p.Watch
	case prev == nil || (g != nil && g.Kind != "" && g.Kind != prev.Geo.Kind) || (p.Kind != "" && prev != nil && p.Kind != prev.Geo.Kind):
		r.Watch = DefaultWatch(r.Geo.Kind)
	}
	if p.WatchTypes != nil {
		wt, rf := checkTypes(*p.WatchTypes)
		if rf != nil {
			return Record{}, rf
		}
		r.WatchTypes = wt
	}
	if p.AppliesTo != nil {
		if len(*p.AppliesTo) > 100 {
			return Record{}, refuse(CodeSchemaInvalid, "applies_to lists more than 100 entities")
		}
		ids := append([]string(nil), (*p.AppliesTo)...)
		for _, a := range ids {
			if !entityIDRe.MatchString(a) {
				return Record{}, refuse(CodeSchemaInvalid, "applies_to id %q is malformed", a)
			}
		}
		sort.Strings(ids)
		r.Properties.AppliesTo = ids
	}
	r.Live = r.Active
	r.Lifecycle = "active"
	if !r.Active {
		r.Lifecycle = "inactive"
	}
	r.UpdatedT, r.UpdatedBy, r.ActionID, r.CommandID, r.Simulated = m.NowMs, m.Actor, m.ActionID, m.CommandID, m.Simulated
	return r, nil
}

// checkTypes validates watch_types. `person` is refused: zone crossings of people are not
// watched (onemind-xr JADC2-MAP.md §5 stop line; the projector excludes people too, G4).
func checkTypes(in []string) ([]string, *Refusal) {
	if len(in) > 32 {
		return nil, refuse(CodeSchemaInvalid, "watch_types lists more than 32 types")
	}
	out := append([]string(nil), in...)
	for _, t := range out {
		if !typeRe.MatchString(t) {
			return nil, refuse(CodeSchemaInvalid, "watch_types entry %q is malformed", t)
		}
		if t == "person" {
			return nil, refuse(CodeForbidden, "zones do not watch people (civilian stop line)")
		}
	}
	sort.Strings(out)
	return out, nil
}

func closePolygon(coords json.RawMessage) json.RawMessage {
	var gj struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	}
	if json.Unmarshal(coords, &gj) != nil {
		return coords
	}
	b, err := json.Marshal(map[string]any{"type": gj.Type, "coordinates": geo.Closed(gj.Coordinates)})
	if err != nil {
		return coords
	}
	return b
}

// Zone is a parsed record with its resolved area.
type Zone struct {
	Record
	Area *geo.Area
	Rev  uint64
}

// Parse reads a KV value written by this service (or by hand in the same shape).
func Parse(key string, value []byte, rev uint64) (*Zone, error) {
	var r Record
	if err := json.Unmarshal(value, &r); err != nil {
		return nil, fmt.Errorf("zone %s: %w", key, err)
	}
	if r.ZoneID == "" {
		r.ZoneID = r.EntityID
	}
	if r.ZoneID == "" {
		r.ZoneID = key
	}
	if !kinds[r.Geo.Kind] {
		return nil, fmt.Errorf("zone %s: unknown kind %q", key, r.Geo.Kind)
	}
	a, err := geo.Parse(r.Geo.Shape, r.Geo.Coords, r.Geo.RadiusM, r.Geo.AltMinM, r.Geo.AltMaxM)
	if err != nil {
		return nil, fmt.Errorf("zone %s: %w", key, err)
	}
	// Same activity rule as the projector: a missing `active` or `live` means true.
	var flags struct {
		Active *bool `json:"active"`
		Live   *bool `json:"live"`
	}
	_ = json.Unmarshal(value, &flags)
	r.Active = (flags.Active == nil || *flags.Active) && (flags.Live == nil || *flags.Live)
	switch r.Lifecycle {
	case "expired", "lost", "deleted", "removed", "inactive":
		r.Active = false
	}
	return &Zone{Record: r, Area: a, Rev: rev}, nil
}
