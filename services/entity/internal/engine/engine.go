// Package engine is the fusion core: identity (alias map, plausibility split), per-source
// monotonic time, per-entity estimator, per-attribute conflict resolution, decomposed
// confidence, lifecycle and change-only publishing (entity.md §5; fusion-and-tracking.md
// Recommendation). It is pure: no NATS. The bus layer feeds it observations and publishes
// what it returns.
package engine

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/onemindos/omos/services/entity/internal/fusion"
	"github.com/onemindos/omos/services/entity/internal/identity"
	"github.com/onemindos/omos/services/entity/internal/model"
)

// Drop and reject reasons, counted by the service.
const (
	DropOlderSourceT    = "older_source_t"
	DropNoKeys          = "no_keys"
	DropExpiredOnArrive = "expired_on_arrival"
	DropTentative       = "tentative"
	RejectPersonHostile = "sl9_person_hostile"
	CountClockSkew      = "clock_skew_replaced"
	CountOutOfOrder     = "out_of_order_kinematics"
	CountSplit          = "duplicate_key_split"
	CountHostileAsset   = "hostile_downgraded"
	CountNoLocationYet  = "no_location_yet"
)

// indexKeys drops, for a split entity (`{id}-x2`), the shared key that names the original, so
// the alias map keeps routing that key to the original entity.
func indexKeys(id string, keys []model.Alias) []model.Alias {
	out := keys[:0:0]
	for _, k := range keys {
		if c := identity.Canonical(k); c != id && strings.HasPrefix(id, c+"-x") {
			continue
		}
		out = append(out, k)
	}
	return out
}

func hasTrait(traits []string, t string) bool {
	for _, x := range traits {
		if x == t {
			return true
		}
	}
	return false
}

// Config tunes the engine.
type Config struct {
	// CoreMode is the interim mode without the ENT stream: no restore is possible, so a new
	// entity's seq starts at the current epoch ms instead of 1, which keeps seq increasing
	// across restarts for core subscribers.
	CoreMode     bool
	SkewMs       int64              // future source_t tolerance; entity.md §5 says 5 s
	HeartbeatMax int64              // ms; entity.md §5 says at most 120 s
	Priors       map[string]float64 // per-integration source_prior overrides
	Links        map[string]string  // declared links "type:id" -> entity_id
	ConfirmHits  int                // M of M-of-N for det tracks (fusion-and-tracking.md Stage 3)
	// OwnReleasableTo is marking.releasable_to for own simulated fleet assets: a `sim-*` id
	// whose every source is a `tel.fleet.*` bridge (entity.md §3 example `["homestead"]`).
	// Any other record, including one a feed or CoT source contributes to, keeps `[]`.
	OwnReleasableTo []string
}

// DefaultConfig returns the contract defaults.
func DefaultConfig() Config {
	return Config{SkewMs: 5000, HeartbeatMax: 120000, ConfirmHits: 3, OwnReleasableTo: []string{"homestead"}}
}

// Publish is one record to publish on `ent.{entity_type}.{entity_id}.state`.
type Publish struct {
	Subject string
	MsgID   string // `{entity_id}:{seq}` (entity.md §4.3)
	TTL     string // Nats-TTL header value: "{n}s" or "never"
	Entity  model.Entity
}

// Reject is a dead-lettered input (stop-line SL9 and friends).
type Reject struct {
	Reason string
	Obs    *model.Observation
}

// Result is what one observation produced.
type Result struct {
	Publish *Publish
	Reject  *Reject
	Drop    string
}

type pubSnap struct {
	lat, lon  float64
	sig       string
	wallT     int64
	lifecycle string
}

type attr struct {
	tier model.Tier
	t    int64
}

type track struct {
	rec       model.Entity
	est       fusion.Estimator
	hits      int
	newestT   int64 // newest accepted source_t
	haeT      int64
	attrs     map[string]attr
	sources   map[string]*model.Source
	srcTier   map[string]model.Tier
	srcPrior  map[string]float64
	srcSim    map[string]bool
	srcEgress map[string]*bool
	srcTLP    map[string]string
	last      *pubSnap
	pending   *Publish // a publish that failed; retried by Sweep
	// connBefore is health.connection as it was before Sweep forced it to "offline" when the
	// entity went stale; restored by the next fleet-* observation that does not state a
	// connection (a CoT source resuming does not restore it).
	connBefore any
	// dead marks an expired track kept as a tombstone until its final `live:false` record is
	// acked (Published), so a failed final publish is retried instead of lost.
	dead bool
}

// Engine holds every live entity. Safe for concurrent use.
type Engine struct {
	mu       sync.Mutex
	cfg      Config
	now      func() int64
	res      *identity.Resolver
	tracks   map[string]*track
	counters map[string]uint64
}

// New returns an engine. now returns epoch ms (injected so tests are deterministic).
func New(cfg Config, now func() int64) *Engine {
	if cfg.SkewMs == 0 {
		cfg.SkewMs = 5000
	}
	if cfg.HeartbeatMax == 0 {
		cfg.HeartbeatMax = 120000
	}
	if cfg.ConfirmHits == 0 {
		cfg.ConfirmHits = 3
	}
	return &Engine{cfg: cfg, now: now, res: identity.NewResolver(cfg.Links),
		tracks: map[string]*track{}, counters: map[string]uint64{}}
}

func (e *Engine) count(k string) { e.counters[k]++ }

// Count adds to a named counter (used by the bus layer for skips).
func (e *Engine) Count(k string) {
	e.mu.Lock()
	e.count(k)
	e.mu.Unlock()
}

// Counters returns a copy of the counters.
func (e *Engine) Counters() map[string]uint64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make(map[string]uint64, len(e.counters)+1)
	for k, v := range e.counters {
		out[k] = v
	}
	out["entities"] = uint64(len(e.tracks))
	return out
}

// Apply fuses one observation.
func (e *Engine) Apply(o *model.Observation) Result {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.now()
	if o.RecvT == 0 {
		o.RecvT = now
	}
	if o.SourceT > now+e.cfg.SkewMs {
		o.SourceT = o.RecvT
		e.count(CountClockSkew)
	}
	// Stop-line SL9: no hostile class for people (JADC2-MAP §5.3 test 9; ONTOLOGY §3.2).
	if (o.Person || o.EntityType == "person") && o.Disposition == "hostile" {
		e.count(RejectPersonHostile)
		return Result{Reject: &Reject{Reason: RejectPersonHostile, Obs: o}}
	}
	// Person safety: a detection never makes a person entity (entity.md §5).
	if o.EntityType == "person" && (o.DataType == "det" || !o.Person) {
		o.EntityType, o.Kind, o.Template = "track", "person", "track"
	}
	var hostileNote *model.Conflict
	if o.Disposition == "hostile" {
		hostileNote = &model.Conflict{Attribute: "disposition", Kept: "suspicious", Rejected: "hostile",
			Sources: []string{o.Integration}}
		o.Disposition = "suspicious"
		e.count(CountHostileAsset)
	}

	id, ok := e.res.Resolve(o.Keys)
	if !ok {
		e.count(DropNoKeys)
		return Result{Drop: DropNoKeys}
	}
	id, split := e.plausible(id, o)

	tr := e.tracks[id]
	var prevSeq uint64
	if tr != nil && tr.dead {
		// The entity reappears while its final record is unacked: the new live record supersedes
		// that record. seq continues so `{entity_id}:{seq}` never repeats.
		prevSeq = tr.rec.Seq
		delete(e.tracks, id)
		e.res.Forget(id)
		tr = nil
	}
	if tr == nil {
		if exp := e.expiryFor(o.EntityType, o.Kind, o.StaleS, o.SourceT, o.ExpireT); exp != nil && *exp <= now {
			e.count(DropExpiredOnArrive)
			return Result{Drop: DropExpiredOnArrive}
		}
		tr = e.newTrack(id, o, now)
		if prevSeq > tr.rec.Seq {
			tr.rec.Seq = prevSeq
		}
		e.tracks[id] = tr
		if split {
			tr.rec.Conflicts = append(tr.rec.Conflicts, model.Conflict{Attribute: "identity",
				Kept: id, Rejected: "duplicate_key", Sources: []string{o.Integration}})
		}
	}

	key := o.Integration + "|" + o.SourceID
	src := tr.sources[key]
	if src != nil && o.SourceT <= src.SourceT {
		e.count(DropOlderSourceT)
		return Result{Drop: DropOlderSourceT}
	}
	newSource := src == nil
	if newSource {
		src = &model.Source{Integration: o.Integration, SourceID: o.SourceID, DataType: o.DataType, FirstSeenT: o.RecvT}
		tr.sources[key] = src
	}
	src.Subject, src.Stream, src.Seq = o.Subject, o.Stream, o.StreamSeq
	src.LastSeenT, src.SourceT, src.Components = o.RecvT, o.SourceT, o.Components
	tr.srcTier[key] = o.Tier
	src.Tier = o.Tier.String()
	tr.srcPrior[key] = sourcePrior(o.Integration, e.cfg.Priors)
	tr.srcSim[key] = o.Simulated
	tr.srcEgress[key] = o.EgressTAK
	tr.srcTLP[key] = o.Marking

	e.res.Register(id, indexKeys(id, o.Keys))
	tr.mergeAliases(o)
	if hostileNote != nil {
		tr.addConflict(*hostileNote)
	}
	tr.setAttrs(o)
	if o.SourceT > tr.newestT {
		tr.newestT = o.SourceT
	}
	if o.HasPos {
		e.kinematics(tr, o)
	}
	for k, v := range o.Meta {
		if tr.rec.Meta == nil {
			tr.rec.Meta = map[string]any{}
		}
		tr.rec.Meta[k] = v
	}
	tr.rec.Status = mergeMap(tr.rec.Status, o.Status)
	tr.rec.Health = mergeMap(tr.rec.Health, o.Health)
	// Only the fleet bridge knows the link state: a CoT report about the same asset (tak-mesh
	// relaying it) does not prove the bridge is back, so the connection stays "offline" until
	// a fleet-* integration reports again.
	if tr.connBefore != nil && strings.HasPrefix(o.Integration, "fleet-") {
		if _, ok := o.Health["connection"]; !ok {
			tr.rec.Health["connection"] = tr.connBefore
		}
		tr.connBefore = nil
	}
	if o.HasCapabilities && o.Simulated && strings.HasPrefix(o.Integration, "fleet-") {
		tr.rec.Capabilities = append([]string{}, o.Capabilities...)
	}
	if o.Emergency {
		tr.rec.Indicators.Emergency = true
	}
	if o.CotSecurity != nil {
		tr.rec.CotSecurity = o.CotSecurity
	}
	stale := o.StaleS
	if stale <= 0 {
		stale = defaultStale(tr.rec.EntityType)
	}
	tr.rec.Stale = stale
	tr.rec.ExpiryT = e.expiryFor(tr.rec.EntityType, tr.rec.Kind, o.StaleS, tr.newestT, o.ExpireT)
	tr.rec.Lifecycle = "active"
	if o.DataType == "det" && tr.hits < e.cfg.ConfirmHits {
		tr.rec.Lifecycle = "tentative"
	}
	e.derive(tr, now)

	if tr.rec.Lifecycle == "tentative" {
		e.count(DropTentative)
		return Result{Drop: DropTentative}
	}
	if tr.rec.Location == nil && hasTrait(tr.rec.Traits, "locatable") {
		e.count(CountNoLocationYet) // a status before the first position: wait for a fix
		return Result{}
	}
	if !e.shouldPublish(tr, now, newSource) {
		return Result{}
	}
	return Result{Publish: e.emit(tr, now)}
}

// plausible implements the AIS duplicate-key split: if the implied speed from the entity's
// last fix exceeds the class limit, the observation goes to `{id}-x2` (or -x3), never averaged.
func (e *Engine) plausible(id string, o *model.Observation) (string, bool) {
	lim := maxSpeed(o.EntityType)
	if !o.HasPos || math.IsInf(lim, 1) {
		return id, false
	}
	ok := func(cand string) (exists, plausible bool) {
		tr := e.tracks[cand]
		if tr == nil || tr.rec.Location == nil {
			return tr != nil, true
		}
		dt := math.Abs(float64(o.SourceT-tr.rec.Location.MeasuredT)) / 1000
		d := fusion.DistanceM(tr.rec.Location.Lat, tr.rec.Location.Lon, o.Lat, o.Lon)
		return true, d/math.Max(dt, 1) <= lim
	}
	if _, p := ok(id); p {
		return id, false
	}
	e.count(CountSplit)
	for n := 2; n <= 3; n++ {
		cand := fmt.Sprintf("%s-x%d", id, n)
		if exists, p := ok(cand); p || !exists {
			return cand, true
		}
	}
	return fmt.Sprintf("%s-x3", id), true
}

func (e *Engine) newTrack(id string, o *model.Observation, now int64) *track {
	seq := uint64(0)
	if e.cfg.CoreMode {
		seq = uint64(now) - 1
	}
	docType := o.DocType
	if docType == "" {
		docType = "onemind." + o.EntityType
	}
	tr := &track{
		attrs: map[string]attr{}, sources: map[string]*model.Source{}, srcTier: map[string]model.Tier{},
		srcPrior: map[string]float64{}, srcSim: map[string]bool{}, srcEgress: map[string]*bool{},
		srcTLP: map[string]string{},
	}
	tr.rec = model.Entity{
		UID: id, Type: docType, EntityID: id, EntityType: o.EntityType, V: model.ContractVersion,
		Seq: seq, Live: true, Lifecycle: "active", CreatedT: now, Template: o.Template,
		Traits: append([]string(nil), o.Traits...), Aliases: model.Aliases{IDs: []model.Alias{}},
		Relationships: []map[string]any{}, Sources: []model.Source{},
	}
	if tr.rec.Template == "" {
		tr.rec.Template = "track"
	}
	if tr.rec.Traits == nil {
		tr.rec.Traits = []string{"locatable"}
	}
	if o.Template != "geo" && o.Template != "spi" {
		tr.est = fusion.NewCV(fusion.ProcessNoise(o.EntityType, o.Kind))
	}
	return tr
}

func (t *track) mergeAliases(o *model.Observation) {
	for _, k := range o.Keys {
		dup := false
		for _, a := range t.rec.Aliases.IDs {
			if a == k {
				dup = true
				break
			}
		}
		if !dup {
			t.rec.Aliases.IDs = append(t.rec.Aliases.IDs, k)
		}
	}
}

func (t *track) addConflict(c model.Conflict) {
	for _, x := range t.rec.Conflicts {
		if x.Attribute == c.Attribute && fmt.Sprint(x.Rejected) == fmt.Sprint(c.Rejected) && fmt.Sprint(x.Kept) == fmt.Sprint(c.Kept) {
			return
		}
	}
	t.rec.Conflicts = append(t.rec.Conflicts, c)
	if n := len(t.rec.Conflicts); n > 10 {
		t.rec.Conflicts = t.rec.Conflicts[n-10:]
	}
}

// setAttrs resolves classification and naming per attribute: human, then self-report, then
// sensor; newest source time within a tier; never downgrade to unknown (entity.md §5).
func (t *track) setAttrs(o *model.Observation) {
	unknown := func(v string) bool { return v == "unknown" || strings.HasSuffix(v, ".unknown") }
	set := func(name, val string, dst *string) {
		if val == "" || (unknown(val) && *dst != "" && !unknown(*dst)) {
			return // classification never downgrades to unknown
		}
		cur, has := t.attrs[name]
		if has && unknown(*dst) && !unknown(val) {
			has = false // a specific class replaces unknown without a conflict
		}
		if has && *dst != "" && *dst != val && (o.Tier > cur.tier || (o.Tier == cur.tier && o.SourceT < cur.t)) {
			t.addConflict(model.Conflict{Attribute: name, Kept: *dst, Rejected: val, Sources: []string{o.Integration}})
			return
		}
		if has && *dst != "" && *dst != val && o.Tier < cur.tier {
			t.addConflict(model.Conflict{Attribute: name, Kept: val, Rejected: *dst, Sources: []string{o.Integration}})
		}
		*dst = val
		t.attrs[name] = attr{tier: o.Tier, t: o.SourceT}
	}
	set("callsign", o.Callsign, &t.rec.Callsign)
	set("disposition", o.Disposition, &t.rec.Disposition)
	set("kind", o.Kind, &t.rec.Kind)
	set("cot_type", o.CotType, &t.rec.CotType)
	set("environment", o.Environment, &t.rec.Environment)
	if o.DocType != "" {
		set("type", o.DocType, &t.rec.Type)
	}
	if t.rec.Callsign != "" {
		name := t.rec.Callsign
		t.rec.Aliases.Name = &name
	}
}

// kinematics updates position. A sensor-tier report corroborates but does not move a track
// that has a fresh self-report or human source (F3 overwatch rule).
func (e *Engine) kinematics(t *track, o *model.Observation) {
	t.hits++
	sigma := o.CE
	if sigma <= 0 {
		sigma = defaultSigma(t.rec.EntityType, t.rec.Kind)
	}
	le := o.LE
	if le <= 0 {
		le = sigma
	}
	if t.est == nil { // geo/event: the point is the location
		t.rec.Location = &model.Location{Lat: round(o.Lat, 7), Lon: round(o.Lon, 7), HAEM: o.HAE, MeasuredT: o.SourceT}
		t.rec.Uncertainty = &model.Uncertainty{CEM: sigma, LEM: le}
		t.rec.Geo = &model.Geo{Shape: "point", Kind: "hazard",
			Coords: map[string]any{"type": "Point", "coordinates": []float64{round(o.Lon, 7), round(o.Lat, 7)}}}
		return
	}
	if o.Tier == model.TierSensor {
		for k, s := range t.sources {
			if t.srcTier[k] < model.TierSensor && s.SourceT >= o.SourceT-int64(t.rec.Stale)*1000 && t.rec.Location != nil {
				return
			}
		}
	}
	m := fusion.Measurement{T: o.SourceT, Lat: o.Lat, Lon: o.Lon, SigmaM: sigma}
	if o.Speed != nil && o.Course != nil {
		m.HasVel = true
		m.VE, m.VN = fusion.VelocityFrom(*o.Speed, *o.Course)
		m.SigmaV = math.Max(1, *o.Speed*0.05)
	}
	if err := t.est.Update(m); err != nil {
		e.count(CountOutOfOrder)
		return
	}
	if o.HAE != nil && o.SourceT >= t.haeT {
		t.haeT = o.SourceT
		h := *o.HAE
		t.rec.HAE = &h
	}
	s := t.est.State()
	sp, crs := s.SpeedCourse()
	ce := round(s.CEM(), 2)
	t.rec.Location = &model.Location{
		Lat: round(s.Lat, 7), Lon: round(s.Lon, 7), HAEM: t.rec.HAE,
		CourseDeg: ptr(round(crs, 1)), SpeedMps: ptr(round(sp, 2)),
		VelENU: []float64{round(s.VE, 3), round(s.VN, 3), 0}, MeasuredT: s.T,
	}
	t.rec.Uncertainty = &model.Uncertainty{
		CEM: ce, LEM: round(le, 2),
		PosENUCov: []float64{round(s.PosCov[0], 3), round(s.PosCov[1], 3), 0, round(s.PosCov[2], 3), 0, round(le*le, 3)},
		VelENUCov: []float64{round(s.VelCov[0], 4), round(s.VelCov[1], 4), 0, round(s.VelCov[2], 4), 0, 0},
	}
}

func (e *Engine) expiryFor(entityType, kind string, stale int, newestT, declared int64) *int64 {
	if entityType == "place" || entityType == "structure" || entityType == "zone" || entityType == "team" {
		return nil
	}
	exp := newestT + int64(ttlSeconds(entityType, kind, stale))*1000
	if declared > 0 {
		exp = declared
	}
	return &exp
}

// derive recomputes the fields that depend on all sources: geo-envelope copies, projector
// properties, sources[], provenance, indicators, marking, confidence.
func (e *Engine) derive(t *track, now int64) {
	r := &t.rec
	r.Pinned = r.ExpiryT == nil
	if l := r.Location; l != nil {
		r.Lat, r.Lon = l.Lat, l.Lon
		r.Course, r.Speed = l.CourseDeg, l.SpeedMps
	}
	r.T = t.newestT
	keys := make([]string, 0, len(t.sources))
	for k := range t.sources {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	r.Sources = r.Sources[:0]
	var primary string
	bestTier, bestT := model.Tier(99), int64(-1)
	sim, allCot, allFleet, anyFleet, tlp := false, true, len(keys) > 0, false, ""
	egressSet := true
	active := map[string]float64{}
	for _, k := range keys {
		s := t.sources[k]
		r.Sources = append(r.Sources, *s)
		if tier := t.srcTier[k]; tier < bestTier || (tier == bestTier && s.SourceT > bestT) {
			bestTier, bestT, primary = tier, s.SourceT, s.Integration
		}
		sim = sim || t.srcSim[k]
		allCot = allCot && s.DataType == "cot"
		allFleet = allFleet && s.DataType == "geo" && strings.HasPrefix(s.Integration, "fleet-")
		anyFleet = anyFleet || strings.HasPrefix(s.Integration, "fleet-")
		tlp = stricterTLP(tlp, t.srcTLP[k])
		if eg := t.srcEgress[k]; eg != nil && !*eg && s.DataType != "cot" {
			egressSet = false
		}
		if t.newestT-s.SourceT <= int64(r.Stale)*3000 {
			if p := t.srcPrior[k]; p > active[s.Integration] {
				active[s.Integration] = p
			}
		}
	}
	if tlp == "" {
		tlp = "AMBER"
	}
	if r.Kind == "person" || r.EntityType == "person" {
		tlp = stricterTLP(tlp, "AMBER+STRICT")
	}
	ownAsset := allFleet && (r.EntityType == "drone" || r.EntityType == "robot" || r.EntityType == "vehicle")
	rel := []string{}
	if ownAsset && sim && strings.HasPrefix(strings.ToLower(r.EntityID), "sim-") {
		rel = append(rel, e.cfg.OwnReleasableTo...)
	}
	r.Marking = model.Marking{TLP: tlp, ReleasableTo: rel}
	if anyFleet && r.Health != nil {
		r.Health["overall"] = fleetOverall(r.Health, r.Lifecycle)
	}
	if ownAsset {
		r.CotType = fleetCotType(r.EntityType, r.Kind, r.Environment) // robotics-cot-profile.md §2
	}

	how := "m-g"
	switch {
	case sim:
		how = "m-s"
	case len(active) > 1:
		how = "m-f"
	case bestTier == model.TierHuman:
		how = "h-e"
	}
	r.Provenance = model.Provenance{Primary: primary, How: how, SourceT: t.newestT}

	// cot-projection.md §3.1 "Who sets egress_tak".
	egress := true
	switch {
	case len(r.EntityID) > 4 && r.EntityID[:4] == "sim-":
		egress = true
	case r.EntityType == "person" || r.EntityType == "agent" || r.EntityType == "team" || r.EntityType == "node" || r.Kind == "person":
		egress = false
	case allCot:
		egress = false // sourced only from TAK: already on TAK
	default:
		egress = egressSet
	}
	r.Indicators.Simulated = sim
	r.Indicators.EgressTAK = egress
	// Capabilities only on a simulated, taskable sim-* asset (adapt.FleetCapabilities); any
	// other record, including a restored one that no longer qualifies, advertises none.
	if !sim || !strings.HasPrefix(strings.ToLower(r.EntityID), "sim-") || !hasTrait(r.Traits, "taskable") {
		r.Capabilities = nil
	}

	// Confidence (entity.md §7): parts kept so the UI and Hermes can explain it.
	prior, miss := 0.0, 1.0
	ints := make([]string, 0, len(active))
	for in := range active {
		ints = append(ints, in)
	}
	sort.Strings(ints)
	for _, in := range ints {
		p := active[in]
		prior = math.Max(prior, p)
		miss *= 1 - p
	}
	corr := math.Min(0.99, 1-miss)
	ce := 0.0
	if r.Uncertainty != nil {
		ce = r.Uncertainty.CEM
	}
	tq := math.Min(1, float64(t.hits)/3) / (1 + ce/1000)
	if t.est == nil {
		tq = math.Min(1, float64(t.hits))
	}
	decay := map[string]float64{"stale": 0.5, "lost": 0.2, "expired": 0}[r.Lifecycle]
	if decay == 0 && r.Lifecycle != "expired" {
		decay = 1
	}
	r.ConfidenceParts = model.ConfidenceParts{SourcePrior: round(prior, 4), TrackQuality: round(tq, 4), Corroboration: round(corr, 4)}
	r.Confidence = round(math.Min(0.99, corr*(0.5+0.5*tq))*decay, 4)

	r.Properties = model.Properties{Callsign: r.Callsign, Altitude: r.HAE, Speed: r.Speed, Course: r.Course}
	if u := r.Uncertainty; u != nil {
		r.Properties.CE, r.Properties.LE = ptr(u.CEM), ptr(u.LEM)
	}
	if r.Template == "track" {
		q := int(math.Round(tq * 15))
		var lastMeas int64
		if r.Location != nil {
			lastMeas = r.Location.MeasuredT
		}
		r.Tracked = &model.Tracked{Quality: q, Hits: t.hits, LastMeasT: lastMeas}
	}
}

func mergeMap(dst, src map[string]any) map[string]any {
	if len(src) == 0 {
		return dst
	}
	if dst == nil {
		dst = map[string]any{}
	}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (t *track) signature() string {
	r := &t.rec
	return fmt.Sprint(r.EntityType, r.Type, r.Callsign, r.Disposition, r.Kind, r.Environment, r.CotType,
		r.Indicators, r.Marking.TLP, r.Lifecycle, len(r.Sources), r.Status, r.Health["battery_pct"],
		r.Health["connection"], r.Capabilities, r.Marking.ReleasableTo)
}

func heartbeatMs(stale int, max int64, pinned bool) int64 {
	if pinned {
		return 302400000 // 3.5 d (BUS-DESIGN §4.1)
	}
	hb := int64(stale) * 500
	if hb > max {
		hb = max
	}
	if hb < 1000 {
		hb = 1000
	}
	return hb
}

func (e *Engine) shouldPublish(t *track, now int64, newSource bool) bool {
	if t.last == nil || newSource {
		return true
	}
	if t.signature() != t.last.sig {
		return true
	}
	if t.rec.Location != nil {
		thr := 1.0
		if t.rec.Uncertainty != nil && t.rec.Uncertainty.CEM > thr {
			thr = t.rec.Uncertainty.CEM
		}
		if fusion.DistanceM(t.last.lat, t.last.lon, t.rec.Location.Lat, t.rec.Location.Lon) > thr {
			return true
		}
	}
	return now-t.last.wallT >= heartbeatMs(t.rec.Stale, e.cfg.HeartbeatMax, t.rec.Pinned)
}

// emit bumps seq and returns the publish. The record is deep-copied through JSON-safe values.
func (e *Engine) emit(t *track, now int64) *Publish {
	t.rec.Seq++
	p := &Publish{
		Subject: "ent." + t.rec.EntityType + "." + t.rec.EntityID + ".state",
		MsgID:   t.rec.EntityID + ":" + strconv.FormatUint(t.rec.Seq, 10),
		TTL:     ttlHeader(t.rec, now),
		Entity:  clone(t.rec),
	}
	t.last = &pubSnap{sig: t.signature(), wallT: now, lifecycle: t.rec.Lifecycle}
	if t.rec.Location != nil {
		t.last.lat, t.last.lon = t.rec.Location.Lat, t.rec.Location.Lon
	}
	t.pending = nil
	e.count("published")
	return p
}

func ttlHeader(r model.Entity, now int64) string {
	if !r.Live {
		return "60s"
	}
	if r.ExpiryT == nil {
		return "never"
	}
	s := (*r.ExpiryT - now + 999) / 1000
	if s < 1 {
		s = 1
	}
	return strconv.FormatInt(s, 10) + "s"
}

// clone copies the slices and maps that the engine keeps mutating.
func clone(r model.Entity) model.Entity {
	c := r
	c.Traits = append([]string(nil), r.Traits...)
	c.Capabilities = append([]string(nil), r.Capabilities...)
	c.Aliases.IDs = append([]model.Alias(nil), r.Aliases.IDs...)
	c.Sources = append([]model.Source(nil), r.Sources...)
	c.Conflicts = append([]model.Conflict(nil), r.Conflicts...)
	c.Meta = mergeMap(nil, r.Meta)
	c.Status = mergeMap(nil, r.Status)
	c.Health = mergeMap(nil, r.Health)
	if r.Location != nil {
		l := *r.Location
		c.Location = &l
	}
	if r.Uncertainty != nil {
		u := *r.Uncertainty
		c.Uncertainty = &u
	}
	if r.Tracked != nil {
		tk := *r.Tracked
		c.Tracked = &tk
	}
	return c
}

// Published confirms that p was stored. A tombstone whose final record this was is deleted.
func (e *Engine) Published(p *Publish) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if t := e.tracks[p.Entity.EntityID]; t != nil && t.dead && t.rec.Seq == p.Entity.Seq {
		e.drop(p.Entity.EntityID)
	}
}

func (e *Engine) drop(id string) {
	delete(e.tracks, id)
	e.res.Forget(id)
}

// PublishFailed keeps a publish for retry by Sweep with the same seq and Nats-Msg-Id, so a
// retry that did reach the stream is deduplicated.
func (e *Engine) PublishFailed(p *Publish) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if t := e.tracks[p.Entity.EntityID]; t != nil && t.rec.Seq == p.Entity.Seq {
		t.pending = p
	}
}

// Sweep advances lifecycles at time now: stale, lost, expiry (final `live:false` record),
// heartbeats and failed-publish retries. Results are in entity_id order. The caller reports
// each result with Published or PublishFailed before the next Sweep; an expired tombstone
// with neither reported is treated as delivered and deleted.
func (e *Engine) Sweep() []*Publish {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.now()
	ids := make([]string, 0, len(e.tracks))
	for id := range e.tracks {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var out []*Publish
	for _, id := range ids {
		t := e.tracks[id]
		if t.pending != nil {
			out = append(out, t.pending)
			continue
		}
		if t.dead {
			e.drop(id)
			continue
		}
		if t.rec.Lifecycle == "tentative" || t.last == nil {
			// never published: expire silently, publish nothing
			if t.rec.ExpiryT != nil && now >= *t.rec.ExpiryT {
				e.drop(id)
			}
			continue
		}
		if t.rec.ExpiryT != nil && now >= *t.rec.ExpiryT {
			t.rec.Live, t.rec.Lifecycle = false, "expired"
			e.derive(t, now)
			out = append(out, e.emit(t, now))
			t.dead = true // deleted by Published once the final record is acked
			e.count("expired")
			continue
		}
		age := now - t.newestT
		lc := "active"
		switch {
		case age > int64(t.rec.Stale)*3000:
			lc = "lost"
		case age > int64(t.rec.Stale)*1000:
			lc = "stale"
		}
		t.rec.Lifecycle = lc
		// A bridge that stopped reporting is no longer "online" (the policy service refuses
		// commands unless it reads "online"); entity.md health.connection enum.
		if c, ok := t.rec.Health["connection"]; ok && lc != "active" && c != "offline" {
			t.connBefore, t.rec.Health["connection"] = c, "offline"
		}
		e.derive(t, now)
		if t.last == nil || t.signature() != t.last.sig ||
			now-t.last.wallT >= heartbeatMs(t.rec.Stale, e.cfg.HeartbeatMax, t.rec.Pinned) {
			out = append(out, e.emit(t, now))
		}
	}
	return out
}

// Restore rebuilds one entity from its last `ENT` record (entity.md §5 "Restart"), so seq,
// per-source source_t and aliases continue where the previous process stopped.
func (e *Engine) Restore(r model.Entity) {
	if !r.Live || r.EntityID == "" {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	t := &track{
		rec: r, attrs: map[string]attr{}, sources: map[string]*model.Source{}, srcTier: map[string]model.Tier{},
		srcPrior: map[string]float64{}, srcSim: map[string]bool{}, srcEgress: map[string]*bool{}, srcTLP: map[string]string{},
	}
	t.newestT = r.Provenance.SourceT
	for i := range r.Sources {
		s := r.Sources[i]
		k := s.Key()
		t.sources[k] = &s
		t.srcPrior[k] = sourcePrior(s.Integration, e.cfg.Priors)
		t.srcSim[k] = r.Indicators.Simulated
		t.srcTLP[k] = r.Marking.TLP
		if tier, ok := model.ParseTier(s.Tier); ok {
			t.srcTier[k] = tier
		} else { // record written before sources[].tier existed
			t.srcTier[k] = model.TierSelfReport
			if s.Integration == r.Provenance.Primary && r.Provenance.How == "h-e" {
				t.srcTier[k] = model.TierHuman
			}
		}
		if s.SourceT > t.newestT {
			t.newestT = s.SourceT
		}
	}
	t.hits = len(r.Sources)
	if r.Tracked != nil {
		t.hits = r.Tracked.Hits
	}
	if r.Template != "geo" && r.Template != "spi" {
		t.est = fusion.NewCV(fusion.ProcessNoise(r.EntityType, r.Kind))
		if l := r.Location; l != nil {
			m := fusion.Measurement{T: l.MeasuredT, Lat: l.Lat, Lon: l.Lon, SigmaM: 5}
			if r.Uncertainty != nil && r.Uncertainty.CEM > 0 {
				m.SigmaM = r.Uncertainty.CEM
			}
			if len(l.VelENU) >= 2 {
				m.HasVel, m.VE, m.VN, m.SigmaV = true, l.VelENU[0], l.VelENU[1], 1
			}
			_ = t.est.Update(m)
			if r.HAE != nil {
				t.haeT = l.MeasuredT
			}
		}
	}
	best := model.TierSensor // attributes were set by the best tier among the sources
	for _, tier := range t.srcTier {
		if tier < best {
			best = tier
		}
	}
	if len(t.srcTier) == 0 {
		best = model.TierSelfReport
	}
	for _, name := range []string{"callsign", "disposition", "kind", "cot_type", "environment", "type"} {
		t.attrs[name] = attr{tier: best, t: t.newestT}
	}
	t.last = &pubSnap{sig: t.signature(), wallT: e.now(), lifecycle: r.Lifecycle}
	if r.Location != nil {
		t.last.lat, t.last.lon = r.Location.Lat, r.Location.Lon
	}
	e.tracks[r.EntityID] = t
	e.res.Register(r.EntityID, indexKeys(r.EntityID, r.Aliases.IDs))
	e.count("restored")
}

// Get returns the current record of one entity.
func (e *Engine) Get(id string) (model.Entity, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	t := e.tracks[id]
	if t == nil || t.dead || t.rec.Lifecycle == "tentative" || t.last == nil {
		return model.Entity{}, false
	}
	return clone(t.rec), true
}

// Query returns published entities filtered by type and bounding box [lon0, lat0, lon1, lat1],
// in entity_id order, with entity_id greater than after (the paging cursor), at most limit
// (0 = no limit). more reports whether further matches follow the last one returned.
func (e *Engine) Query(types []string, bbox []float64, after string, limit int) (out []model.Entity, more bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	want := map[string]bool{}
	for _, t := range types {
		want[t] = true
	}
	ids := make([]string, 0, len(e.tracks))
	for id := range e.tracks {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		t := e.tracks[id]
		if id <= after || t.dead || t.rec.Lifecycle == "tentative" || t.last == nil {
			continue
		}
		if len(want) > 0 && !want[t.rec.EntityType] {
			continue
		}
		if len(bbox) == 4 && !(t.rec.Lon >= bbox[0] && t.rec.Lat >= bbox[1] && t.rec.Lon <= bbox[2] && t.rec.Lat <= bbox[3]) {
			continue
		}
		if limit > 0 && len(out) >= limit {
			return out, true
		}
		out = append(out, clone(t.rec))
	}
	return out, false
}
