// Package model holds the entity record (onemind-xr docs/contracts/entity.md, v1) and the
// normalised observation that every input adapter produces (fusion-and-tracking.md Stage 0).
//
// JSON field names are the contract's, exactly. Required fields have no omitempty; fields the
// contract allows to be null are pointers without omitempty; optional fields use omitempty.
package model

// ContractVersion is the `v` field of entity.v1.
const ContractVersion = 1

// Entity is one `ent.{entity_type}.{entity_id}.state` record (entity.md §4.1).
type Entity struct {
	// Group A: doctrine §9.2 geo-envelope.
	UID      string         `json:"uid"`
	Type     string         `json:"type"`
	Lat      float64        `json:"lat"`
	Lon      float64        `json:"lon"`
	HAE      *float64       `json:"hae,omitempty"`
	Course   *float64       `json:"course,omitempty"`
	Speed    *float64       `json:"speed,omitempty"`
	Callsign string         `json:"callsign,omitempty"`
	Stale    int            `json:"stale,omitempty"`
	T        int64          `json:"t"`
	Meta     map[string]any `json:"meta,omitempty"`

	// Group B: fields the TAK projector reads today.
	EntityID   string     `json:"entity_id"`
	EntityType string     `json:"entity_type"`
	Confidence float64    `json:"confidence"`
	Properties Properties `json:"properties"`

	// Group C: OneMind components.
	V               int              `json:"v"`
	Seq             uint64           `json:"seq"`
	Live            bool             `json:"live"`
	Lifecycle       string           `json:"lifecycle"`
	CreatedT        int64            `json:"created_t"`
	ExpiryT         *int64           `json:"expiry_t"`
	Pinned          bool             `json:"pinned"`
	Kind            string           `json:"kind,omitempty"`
	Template        string           `json:"template"`
	Traits          []string         `json:"traits"`
	Disposition     string           `json:"disposition,omitempty"`
	Environment     string           `json:"environment,omitempty"`
	CotType         string           `json:"cot_type,omitempty"`
	Aliases         Aliases          `json:"aliases"`
	Location        *Location        `json:"location,omitempty"`
	Uncertainty     *Uncertainty     `json:"uncertainty,omitempty"`
	Geo             *Geo             `json:"geo"`
	Status          map[string]any   `json:"status,omitempty"`
	Health          map[string]any   `json:"health,omitempty"`
	CurrentTask     map[string]any   `json:"current_task"`
	Tracked         *Tracked         `json:"tracked"`
	Correlation     map[string]any   `json:"correlation"`
	Relationships   []map[string]any `json:"relationships"`
	Route           map[string]any   `json:"route"`
	Indicators      Indicators       `json:"indicators"`
	Marking         Marking          `json:"marking"`
	CotSecurity     map[string]any   `json:"cot_security"`
	Provenance      Provenance       `json:"provenance"`
	Sources         []Source         `json:"sources"`
	ConfidenceParts ConfidenceParts  `json:"confidence_parts"`
	Conflicts       []Conflict       `json:"conflicts,omitempty"`
	Overrides       []map[string]any `json:"overrides,omitempty"`
}

// Properties is the projector block, copied from the components (entity.md §4.1 group B).
type Properties struct {
	Callsign string   `json:"callsign,omitempty"`
	Altitude *float64 `json:"altitude,omitempty"`
	Speed    *float64 `json:"speed,omitempty"`
	Course   *float64 `json:"course,omitempty"`
	CE       *float64 `json:"ce,omitempty"`
	LE       *float64 `json:"le,omitempty"`
}

// Alias is one identifier of an entity in another system.
type Alias struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Aliases is the `aliases` component.
type Aliases struct {
	Name *string `json:"name"`
	IDs  []Alias `json:"ids"`
}

// Location is the `location` component.
type Location struct {
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	HAEM      *float64  `json:"hae_m,omitempty"`
	AGLM      *float64  `json:"agl_m,omitempty"`
	CourseDeg *float64  `json:"course_deg,omitempty"`
	SpeedMps  *float64  `json:"speed_mps,omitempty"`
	VelENU    []float64 `json:"vel_enu_mps,omitempty"`
	AttitudeQ []float64 `json:"attitude_q,omitempty"`
	MeasuredT int64     `json:"measured_t"`
}

// Uncertainty is the `uncertainty` component; ce_m/le_m are 1-sigma.
type Uncertainty struct {
	CEM          float64        `json:"ce_m"`
	LEM          float64        `json:"le_m"`
	PosENUCov    []float64      `json:"pos_enu_cov"`
	VelENUCov    []float64      `json:"vel_enu_cov"`
	ErrorEllipse map[string]any `json:"error_ellipse"`
}

// Geo is the `geo` component for template=geo.
type Geo struct {
	Shape  string         `json:"shape"`
	Coords map[string]any `json:"coords"`
	Kind   string         `json:"kind"`
}

// Tracked is the `tracked` component for template=track.
type Tracked struct {
	Quality   int            `json:"quality"`
	Hits      int            `json:"hits"`
	LastMeasT int64          `json:"last_meas_t"`
	LOB       map[string]any `json:"lob"`
}

// Indicators is the `indicators` component.
type Indicators struct {
	Simulated bool `json:"simulated"`
	Exercise  bool `json:"exercise"`
	Emergency bool `json:"emergency"`
	Starred   bool `json:"starred"`
	EgressTAK bool `json:"egress_tak"`
}

// Marking is the release marking.
type Marking struct {
	TLP          string   `json:"tlp"`
	ReleasableTo []string `json:"releasable_to"`
}

// Provenance names the primary contributor.
type Provenance struct {
	Primary string `json:"primary"`
	How     string `json:"how"`
	SourceT int64  `json:"source_t"`
}

// Source is one contributing source with its stream pointer (the history pointer).
type Source struct {
	Integration string   `json:"integration"`
	SourceID    string   `json:"source_id"`
	DataType    string   `json:"data_type"`
	Subject     string   `json:"subject"`
	Stream      string   `json:"stream"`
	Seq         uint64   `json:"seq"`
	FirstSeenT  int64    `json:"first_seen_t"`
	LastSeenT   int64    `json:"last_seen_t"`
	SourceT     int64    `json:"source_t"`
	Components  []string `json:"components"`
	// Tier is this source's attribute authority ("human", "self_report", "sensor"), kept so a
	// restart restores conflict tiers exactly. Additive field (sources items allow extra keys).
	Tier string `json:"tier,omitempty"`
}

// Key identifies the source for the per-source monotonic rule.
func (s Source) Key() string { return s.Integration + "|" + s.SourceID }

// ConfidenceParts are the three explained parts of `confidence` (entity.md §7).
type ConfidenceParts struct {
	SourcePrior   float64 `json:"source_prior"`
	TrackQuality  float64 `json:"track_quality"`
	Corroboration float64 `json:"corroboration"`
}

// Conflict records one per-attribute resolution.
type Conflict struct {
	Attribute string   `json:"attribute"`
	Kept      any      `json:"kept"`
	Rejected  any      `json:"rejected"`
	Sources   []string `json:"sources"`
}

// Tier orders attribute authority (entity.md §5 "Conflicts").
type Tier int

const (
	TierHuman      Tier = 0 // operator or TAK-declared assertion
	TierSelfReport Tier = 1 // cooperative self-report (ADS-B, AIS, MAVLink, PLI)
	TierSensor     Tier = 2 // sensor estimate or third-party report
)

var tierNames = [...]string{"human", "self_report", "sensor"}

// String returns the name stored in Source.Tier.
func (t Tier) String() string {
	if t >= 0 && int(t) < len(tierNames) {
		return tierNames[t]
	}
	return ""
}

// ParseTier reads a Source.Tier name.
func ParseTier(s string) (Tier, bool) {
	for i, n := range tierNames {
		if n == s {
			return Tier(i), true
		}
	}
	return 0, false
}

// Observation is one normalised input message (fusion-and-tracking.md Stage 0).
type Observation struct {
	Integration string // feed or bridge id, e.g. "feed-aircraft", "tak-eud", "world-air"
	SourceID    string // native id inside that integration
	DataType    string // cot | geo | intel | det
	Subject     string
	Stream      string
	StreamSeq   uint64

	Keys        []Alias // hard identity keys, most specific first
	EntityType  string
	Kind        string
	DocType     string // doctrine §9.3 dot-path
	Callsign    string
	CotType     string
	Disposition string // may be "hostile" here; the engine never publishes it
	Environment string
	Template    string
	Traits      []string
	Tier        Tier
	How         string

	HasPos bool
	Lat    float64
	Lon    float64
	HAE    *float64
	Course *float64
	Speed  *float64
	CE     float64 // 1-sigma metres; 0 = unknown, use the class default
	LE     float64

	SourceT int64 // epoch ms from the payload
	RecvT   int64 // epoch ms when the service received it
	StaleS  int   // seconds the source says the reading stays fresh; 0 = unknown
	ExpireT int64 // absolute expiry the source declares (intel `expires`); 0 = none

	Meta        map[string]any
	Status      map[string]any
	Health      map[string]any
	CotSecurity map[string]any
	Marking     string // TLP the adapter derives from the input path
	EgressTAK   *bool  // nil = engine default
	Simulated   bool
	Emergency   bool
	Person      bool // the record asserts a person (TAK EUD); see SL9
	Components  []string
}
