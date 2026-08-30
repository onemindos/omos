package main

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

// om-overwatch-state — the fusion spine of om-overwatch.
//
// Subscribes to every sensor feed on the OneMind bus (:4222), fuses each report
// into a single Entity world-state record, persists it in the OVERWATCH_STATE
// JetStream KV bucket (keyed by entity_id), and republishes the fused state on
// ent.air.{id}.state. Every other organ (C2 displays, interceptor) reads from
// here. No embedded NATS, no SQLite, no web UI — just the spine.
//
// Inputs:
//   det.air.contact.{id}   — RF/acoustic detections (om-overwatch-sentinel)
//   det.air.visual.{id}    — visual confirmations   (om-overwatch-vision)
//   tel.drone.{id}.{stream}— friendly UAV telemetry (om-overwatch-relay)
//   det.cot.{id}           — TAK CoT tracks (om-overwatch-tak, human affiliation)
// Output:
//   ent.air.{id}.state     — fused entity state (also KV: OVERWATCH_STATE/{id})

const (
	kvBucket      = "OVERWATCH_STATE"
	hbSubject     = "agents.hb.overwatch.state.zeus"
	staleAfter    = 30 * time.Second // no update -> mark lost
	sweepInterval = 10 * time.Second
	// assocRadiusM is the cross-sensor association radius. A detection with a
	// NEW id that lands within this distance of an existing active air track is
	// fused into it (confirmation) instead of spawning a duplicate blip. Kept
	// tight enough that genuinely distinct drones (typically >1 km apart) stay
	// separate, loose enough to absorb monocular-vision range error.
	assocRadiusM = 500.0
)

// contactReport is the permissive shape accepted from detection sensors.
type contactReport struct {
	ID          string   `json:"id"`
	EntityType  string   `json:"entity_type"`
	Affiliation string   `json:"affiliation"`
	Name        string   `json:"name"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Altitude    *float64 `json:"altitude"`
	Heading     *float64 `json:"heading"`
	Velocity    *float64 `json:"velocity"`
	Confidence  float64  `json:"confidence"`
	Source      string   `json:"source"`
}

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = "nats://127.0.0.1:4222"
	}

	nc, err := nats.Connect(url, nats.Name("om-overwatch-state"),
		nats.MaxReconnects(-1), nats.ReconnectWait(2*time.Second))
	if err != nil {
		log.Fatal("NATS connect:", err)
	}
	defer nc.Close()
	log.Printf("om-overwatch-state: connected to OneMind bus at %s", url)

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("JetStream:", err)
	}

	kv, err := js.KeyValue(kvBucket)
	if err != nil {
		kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:      kvBucket,
			Description: "om-overwatch fused entity world-state (keyed by entity_id)",
			History:     5,
			Storage:     nats.FileStorage,
		})
		if err != nil {
			log.Fatal("KV bucket:", err)
		}
		log.Printf("om-overwatch-state: created KV bucket %s", kvBucket)
	} else {
		log.Printf("om-overwatch-state: opened KV bucket %s", kvBucket)
	}

	spine := &Spine{nc: nc, kv: kv}

	// Detection feeds — RF/acoustic and visual.
	nc.Subscribe("det.air.contact.>", func(m *nats.Msg) {
		spine.ingestContact(m, EntityTypeAircraftUnknown, AffiliationHostile)
	})
	nc.Subscribe("det.air.visual.>", func(m *nats.Msg) {
		spine.ingestContact(m, EntityTypeAircraftUnknown, AffiliationUnknown)
	})
	// Friendly UAV telemetry — subject: tel.drone.{id}.{stream}
	nc.Subscribe("tel.drone.>", func(m *nats.Msg) {
		spine.ingestTelemetry(m)
	})
	// TAK CoT tracks — subject: det.cot.{id} (om-overwatch-tak). Human-tagged
	// affiliation from the operator's ATAK is authoritative (locks affiliation).
	nc.Subscribe("det.cot.>", func(m *nats.Msg) {
		spine.ingestCoT(m)
	})
	// Manual entities from om-overwatch-api — the spine stays the single KV writer.
	nc.Subscribe("ent.manual.upsert", func(m *nats.Msg) {
		spine.manualUpsert(m)
	})
	nc.Subscribe("ent.manual.delete", func(m *nats.Msg) {
		spine.manualDelete(m)
	})

	log.Println("om-overwatch-state: fusion active (det.air.contact.>, det.air.visual.>, tel.drone.>, det.cot.>, ent.manual.*)")

	go spine.sweepStale()

	// Heartbeat.
	nc.Publish(hbSubject, []byte(`{"status":"up"}`))
	ticker := time.NewTicker(15 * time.Second)
	for range ticker.C {
		nc.Publish(hbSubject, []byte(`{"t":"`+time.Now().UTC().Format(time.RFC3339)+`","status":"fusing"}`))
	}
}

// Spine holds the shared NATS + KV handles and the fusion logic.
type Spine struct {
	nc *nats.Conn
	kv nats.KeyValue
}

// lastToken returns the final dot-delimited token of a subject.
func lastToken(subject string) string {
	parts := strings.Split(subject, ".")
	return parts[len(parts)-1]
}

// idFromSubject pulls the entity id token. For det.air.contact.{id} the id is
// the last token. For tel.drone.{id}.{stream} the id is the 3rd token.
func idFromSubject(subject string, index int) string {
	parts := strings.Split(subject, ".")
	if index < 0 || index >= len(parts) {
		return ""
	}
	return parts[index]
}

// load fetches the current entity from KV, or nil if absent.
func (s *Spine) load(id string) *Entity {
	entry, err := s.kv.Get(id)
	if err != nil {
		return nil
	}
	var e Entity
	if json.Unmarshal(entry.Value(), &e) != nil {
		return nil
	}
	return &e
}

// save persists the entity to KV and republishes fused state.
func (s *Spine) save(e *Entity) {
	e.UpdatedAt = time.Now().UTC()
	data, err := json.Marshal(e)
	if err != nil {
		return
	}
	if _, err := s.kv.Put(e.EntityID, data); err != nil {
		log.Printf("kv put %s: %v", e.EntityID, err)
		return
	}
	s.nc.Publish("ent.air."+e.EntityID+".state", data)
}

// ingestContact fuses a detection report into entity state.
//
// Fusion has three cases:
//  1. an entity with this id already exists  → normal same-track update
//  2. a NEW id lands within assocRadiusM of an existing active air track
//     → proximity MERGE (cross-sensor confirmation): the report corroborates
//       the existing track rather than spawning a duplicate blip
//  3. otherwise → a brand-new track
//
// This is what turns two sensors seeing the same drone (acoustic TRK-… and
// vision CUID) into ONE fused entity with sources=[sentinel,vision] and a
// confidence boosted by their independent agreement.
func (s *Spine) ingestContact(m *nats.Msg, defType, defAffil string) {
	var r contactReport
	if json.Unmarshal(m.Data, &r) != nil {
		return
	}
	id := r.ID
	if id == "" {
		id = lastToken(m.Subject)
	}
	if id == "" {
		return
	}

	now := time.Now().UTC()
	src := r.Source
	if src == "" {
		src = idFromSubject(m.Subject, 2) // "contact" or "visual"
	}

	e := s.load(id)
	merged := false
	if e == nil && r.Latitude != nil && r.Longitude != nil {
		// No track with this id — try to confirm an existing nearby air track
		// before creating a new one.
		if near := s.associateNearby(*r.Latitude, *r.Longitude, r.Affiliation); near != nil {
			e = near
			merged = true
		}
	}
	if e == nil {
		e = &Entity{
			EntityID:    id,
			EntityType:  defType,
			Affiliation: defAffil,
			Status:      StatusActive,
			Priority:    PriorityHigh, // unclassified air contact defaults to high
			FirstSeen:   now,
		}
	}

	// Entity type: adopt specific classifications, but never downgrade a
	// classified track back to "unknown" just because another sensor is coarse.
	if r.EntityType != "" && (e.EntityType == "" || e.EntityType == EntityTypeAircraftUnknown || r.EntityType != EntityTypeAircraftUnknown) {
		e.EntityType = r.EntityType
	}
	// Respect human affiliation authority: a machine sensor never overwrites an
	// affiliation a human asserted via TAK CoT or the operator API.
	if r.Affiliation != "" && !e.AffiliationLocked {
		e.Affiliation = r.Affiliation
	}
	if r.Name != "" {
		e.Name = r.Name
	}
	// Confidence: an independent NEW sensor corroborating a track raises
	// certainty (probabilistic OR); a same-sensor update just tracks the latest.
	newSource := !e.hasSource(src)
	if r.Confidence > 0 {
		if newSource && e.Confidence > 0 {
			e.Confidence = 1 - (1-e.Confidence)*(1-r.Confidence)
		} else {
			e.Confidence = r.Confidence
		}
	}
	// Position: a same-id sensor owns its track's position. A proximity merge
	// (confirmation from a different, possibly coarser sensor) must NOT reposition
	// the primary track — it only corroborates identity/classification/confidence.
	if !merged && r.Latitude != nil && r.Longitude != nil {
		if e.Position == nil {
			e.Position = &Position{}
		}
		e.Position.Latitude = *r.Latitude
		e.Position.Longitude = *r.Longitude
		if r.Altitude != nil {
			e.Position.Altitude = *r.Altitude
		}
		if r.Heading != nil {
			e.Position.Heading = *r.Heading
		}
		if r.Velocity != nil {
			e.Position.Velocity = *r.Velocity
		}
	}
	e.addSource(src)
	e.Status = StatusActive
	e.IsLive = true
	e.LastSeen = now
	if merged {
		log.Printf("om-overwatch-state: %s confirmed by %s → sources=%v conf=%.2f", e.EntityID, src, e.Sources, e.Confidence)
	}
	s.save(e)
}

// associateNearby finds the closest active AIR track within assocRadiusM of the
// given point, for cross-sensor confirmation. Returns nil if none qualifies.
//
// Safety rule: it will not fuse into a track whose affiliation a human has
// LOCKED to a different value — better to show an extra blip than to hide a
// hostile inside an operator-declared-friendly track.
//
// O(n) over KV keys per detection; fine at counter-UAS track counts. A spatial
// index is the obvious optimization if the COP ever grows large.
func (s *Spine) associateNearby(lat, lon float64, affil string) *Entity {
	keys, err := s.kv.Keys()
	if err != nil {
		return nil
	}
	var best *Entity
	bestDist := assocRadiusM
	for _, k := range keys {
		e := s.load(k)
		if e == nil || e.Position == nil || e.Status == StatusLost {
			continue
		}
		if !strings.HasPrefix(e.EntityType, "aircraft") {
			continue // air-to-air association only
		}
		if e.AffiliationLocked && affil != "" && affil != e.Affiliation {
			continue // don't hide a differently-classified contact in a locked track
		}
		if d := haversineM(lat, lon, e.Position.Latitude, e.Position.Longitude); d <= bestDist {
			bestDist = d
			best = e
		}
	}
	return best
}

// haversineM returns great-circle distance in metres between two lat/lon points.
func haversineM(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6_371_000.0
	rad := math.Pi / 180.0
	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * R * math.Asin(math.Sqrt(a))
}

// ingestTelemetry fuses friendly UAV telemetry (tel.drone.{id}.{stream}).
func (s *Spine) ingestTelemetry(m *nats.Msg) {
	id := idFromSubject(m.Subject, 2)
	if id == "" {
		return
	}
	now := time.Now().UTC()
	e := s.load(id)
	if e == nil {
		e = &Entity{
			EntityID:    id,
			EntityType:  EntityTypeAircraftMultirotor,
			Affiliation: AffiliationFriendly,
			Status:      StatusActive,
			Priority:    PriorityNormal,
			FirstSeen:   now,
		}
	}
	// Best-effort position extraction from common telemetry field names.
	var body map[string]interface{}
	if json.Unmarshal(m.Data, &body) == nil {
		if pos := extractPos(body); pos != nil {
			e.Position = pos
		}
	}
	e.addSource("mavlink")
	// Friendly UAV telemetry asserts friendly — unless a human has locked a
	// different affiliation (e.g. a captured/spoofed airframe tagged hostile).
	if !e.AffiliationLocked {
		e.Affiliation = AffiliationFriendly
	}
	e.Status = StatusActive
	e.IsLive = true
	e.LastSeen = now
	s.save(e)
}

// ingestCoT fuses a TAK CoT track (det.cot.{id}) from om-overwatch-tak. CoT
// carries human-asserted affiliation (a-f/h/n/u), which is authoritative: it
// locks affiliation so downstream machine sensors cannot override the
// operator's friend/foe call. CoT tracks are live (staleness-managed), not
// pinned. Reuses the permissive contactReport shape.
func (s *Spine) ingestCoT(m *nats.Msg) {
	var r contactReport
	if json.Unmarshal(m.Data, &r) != nil {
		return
	}
	id := r.ID
	if id == "" {
		id = lastToken(m.Subject)
	}
	if id == "" {
		return
	}
	now := time.Now().UTC()
	e := s.load(id)
	if e == nil {
		e = &Entity{
			EntityID:    id,
			EntityType:  EntityTypeSensorPlatform, // overridden below if provided
			Affiliation: AffiliationUnknown,
			Status:      StatusActive,
			Priority:    PriorityNormal,
			FirstSeen:   now,
		}
	}
	if r.EntityType != "" {
		e.EntityType = r.EntityType
	}
	if r.Affiliation != "" {
		e.Affiliation = r.Affiliation
		e.AffiliationLocked = true // human-in-the-loop: CoT affiliation is final
	}
	if r.Name != "" {
		e.Name = r.Name
	}
	if r.Latitude != nil && r.Longitude != nil {
		if e.Position == nil {
			e.Position = &Position{}
		}
		e.Position.Latitude = *r.Latitude
		e.Position.Longitude = *r.Longitude
		if r.Altitude != nil {
			e.Position.Altitude = *r.Altitude
		}
		if r.Heading != nil {
			e.Position.Heading = *r.Heading
		}
		if r.Velocity != nil {
			e.Position.Velocity = *r.Velocity
		}
	}
	e.addSource("cot")
	e.Status = StatusActive
	e.IsLive = true
	e.LastSeen = now
	s.save(e)
}

// manualUpsert applies a full Entity supplied by om-overwatch-api. Manual
// entities are pinned (exempt from staleness sweep) and preserve their
// operator-set fields. The spine remains the single writer of KV.
func (s *Spine) manualUpsert(m *nats.Msg) {
	var in Entity
	if json.Unmarshal(m.Data, &in) != nil || in.EntityID == "" {
		return
	}
	now := time.Now().UTC()
	e := s.load(in.EntityID)
	if e == nil {
		e = &Entity{EntityID: in.EntityID, FirstSeen: now}
	}
	if in.Name != "" {
		e.Name = in.Name
	}
	if in.EntityType != "" {
		e.EntityType = in.EntityType
	}
	if in.Affiliation != "" {
		e.Affiliation = in.Affiliation
		e.AffiliationLocked = true // operator-set affiliation is authoritative
	}
	if in.Priority != "" {
		e.Priority = in.Priority
	}
	if in.Position != nil {
		e.Position = in.Position
	}
	if in.Metadata != nil {
		e.Metadata = in.Metadata
	}
	if e.Status == "" {
		e.Status = StatusActive
	}
	e.addSource("manual")
	e.Pinned = true // operator-created entities persist until explicitly deleted
	e.IsLive = true
	e.LastSeen = now
	s.save(e)
	log.Printf("om-overwatch-state: manual upsert %s (%s)", e.EntityID, e.EntityType)
}

// manualDelete removes an operator-created entity from KV.
func (s *Spine) manualDelete(m *nats.Msg) {
	id := strings.TrimSpace(string(m.Data))
	if id == "" {
		// also accept {"id":"..."}
		var body struct {
			ID string `json:"id"`
		}
		if json.Unmarshal(m.Data, &body) == nil {
			id = body.ID
		}
	}
	if id == "" {
		return
	}
	if err := s.kv.Delete(id); err != nil {
		log.Printf("om-overwatch-state: manual delete %s: %v", id, err)
		return
	}
	s.nc.Publish("ent.air."+id+".deleted", []byte(`{"entity_id":"`+id+`","deleted":true}`))
	log.Printf("om-overwatch-state: manual delete %s", id)
}

// extractPos pulls lat/lon/alt from a loosely-typed telemetry map. Handles both
// flat payloads ({"lat":..}) and enveloped payloads from om-overwatch-relay
// ({"message_type":..,"data":{"lat":..}}).
func extractPos(body map[string]interface{}) *Position {
	// Merge a nested "data" envelope (om-overwatch-relay format) over the top level.
	scope := body
	if inner, ok := body["data"].(map[string]interface{}); ok {
		scope = inner
	}
	get := func(keys ...string) (float64, bool) {
		for _, k := range keys {
			if v, ok := scope[k]; ok {
				if f, ok := v.(float64); ok {
					return f, true
				}
			}
		}
		return 0, false
	}
	lat, okLat := get("lat", "latitude")
	lon, okLon := get("lon", "lng", "longitude")
	if !okLat || !okLon {
		return nil
	}
	p := &Position{Latitude: lat, Longitude: lon}
	if alt, ok := get("alt", "altitude", "relative_alt"); ok {
		p.Altitude = alt
	}
	if hdg, ok := get("heading", "hdg", "yaw"); ok {
		p.Heading = hdg
	}
	if vel, ok := get("velocity", "groundspeed", "speed"); ok {
		p.Velocity = vel
	}
	return p
}

// sweepStale marks entities lost when they stop reporting.
func (s *Spine) sweepStale() {
	ticker := time.NewTicker(sweepInterval)
	for range ticker.C {
		keys, err := s.kv.Keys()
		if err != nil {
			continue
		}
		now := time.Now().UTC()
		for _, k := range keys {
			e := s.load(k)
			if e == nil || e.Status == StatusLost || e.Pinned {
				continue // pinned (manual) entities never go stale
			}
			if now.Sub(e.LastSeen) > staleAfter {
				e.Status = StatusLost
				e.IsLive = false
				s.save(e)
				log.Printf("om-overwatch-state: %s marked lost (stale %.0fs)", k, now.Sub(e.LastSeen).Seconds())
			}
		}
	}
}
