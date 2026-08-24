package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// om-overwatch-interceptor — organ 7/7, the effector. Closes the kill chain:
// detect (sentinel/vision) → fuse (spine) → display (tak) → INTERCEPT.
//
// This is the NATS-native completion of SentinelX (a ROS2/PX4 interceptor whose
// guidance was left as a TODO upstream). It adopts SentinelX's ontology + phase/
// authority state machine + C2 command protocol verbatim, implements the missing
// midcourse + terminal proportional-navigation guidance, and simulates point-mass
// flight so the whole engagement closes on the bus. On real hardware, a DDS↔NATS
// bridge on the interceptor's Linux companion swaps the sim for px4_proxy telemetry.
//
// Contract (one bus :4222):
//   IN  int.command.{id}   {"type":"ASSIGN_TARGET|START_INTERCEPT|ABORT",
//                           "mission_id","target_id", target_lat/lon/alt, target_vn/ve/vu}
//   IN  ent.air.{tgt}.state (spine) — live target track; dead-reckoned between updates
//   OUT int.ack.{id} int.phase.{id} int.target.{id} int.health.{id} int.result.{id}
//   OUT tel.drone.{id}     — interceptor's own telemetry (friendly) → spine → COP/ATAK
//   KV  OVERWATCH_INTERCEPT/{id}
//   HB  agents.hb.overwatch.interceptor.zeus (15s)

const (
	subjCommand  = "int.command.*"
	subjTgtState = "ent.air.*.state"
	hbSubject    = "agents.hb.overwatch.interceptor.zeus"
	kvBucket     = "OVERWATCH_INTERCEPT"
	simHz        = 20
)

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func envF(k string, def float64) float64 {
	if v := os.Getenv(k); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

type command struct {
	Type      string  `json:"type"`
	MissionID string  `json:"mission_id"`
	TargetID  string  `json:"target_id"`
	TargetLat float64 `json:"target_lat"`
	TargetLon float64 `json:"target_lon"`
	TargetAlt float64 `json:"target_alt"`
	TargetVN  float64 `json:"target_vn"` // north m/s
	TargetVE  float64 `json:"target_ve"` // east  m/s
	TargetVU  float64 `json:"target_vu"` // up    m/s
	// Optional per-mission launch site (defense point). Falls back to the
	// organ's configured base when omitted.
	BaseLat float64 `json:"base_lat"`
	BaseLon float64 `json:"base_lon"`
	BaseAlt float64 `json:"base_alt"`
}

type manager struct {
	nc     *nats.Conn
	kv     nats.KeyValue
	origin geoOrigin
	params guidanceParams

	mu   sync.Mutex
	eng  map[string]*Engagement
	last map[string]Phase // last-published phase (for change detection)
}

func main() {
	url := env("NATS_URL", "nats://127.0.0.1:4222")
	nc, err := nats.Connect(url, nats.Name("om-overwatch-interceptor"),
		nats.MaxReconnects(-1), nats.ReconnectWait(2*time.Second))
	if err != nil {
		log.Fatal("NATS connect:", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("JetStream:", err)
	}
	kv, err := js.KeyValue(kvBucket)
	if err != nil {
		kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:      kvBucket,
			Description: "om-overwatch interceptor engagement status",
			History:     1,
		})
		if err != nil {
			log.Fatal("KV bucket:", err)
		}
	}

	m := &manager{
		nc: nc, kv: kv,
		origin: geoOrigin{
			lat0: envF("INT_BASE_LAT", 44.3900),
			lon0: envF("INT_BASE_LON", 26.0900),
			alt0: envF("INT_BASE_ALT", 0),
		},
		params: defaultParams(),
		eng:    map[string]*Engagement{},
		last:   map[string]Phase{},
	}
	log.Printf("om-overwatch-interceptor: connected to OneMind bus at %s (base %.4f,%.4f)",
		url, m.origin.lat0, m.origin.lon0)

	if _, err := nc.Subscribe(subjCommand, m.handleCommand); err != nil {
		log.Fatal("subscribe int.command:", err)
	}
	if _, err := nc.Subscribe(subjTgtState, m.handleTargetState); err != nil {
		log.Fatal("subscribe ent.air.*.state:", err)
	}
	log.Printf("om-overwatch-interceptor: armed — listening on %s (targets via %s)", subjCommand, subjTgtState)

	go m.heartbeat()
	m.simLoop() // blocks forever
}

func (m *manager) getOrCreate(id string) *Engagement {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.eng[id]
	if !ok {
		e = newEngagement(id, m.origin, m.params)
		m.eng[id] = e
	}
	return e
}

func (m *manager) handleCommand(msg *nats.Msg) {
	id := strings.TrimPrefix(msg.Subject, "int.command.")
	if id == "" || strings.Contains(id, ".") {
		return
	}
	var c command
	if err := json.Unmarshal(msg.Data, &c); err != nil {
		log.Printf("[int] %s: bad command: %v", id, err)
		return
	}
	e := m.getOrCreate(id)

	switch strings.ToUpper(c.Type) {
	case "ASSIGN_TARGET":
		if c.BaseLat != 0 || c.BaseLon != 0 {
			e.setOrigin(geoOrigin{lat0: c.BaseLat, lon0: c.BaseLon, alt0: c.BaseAlt})
		}
		ok := e.assignTarget(c.MissionID, c.TargetID, c.TargetLat, c.TargetLon, c.TargetAlt, c.TargetVN, c.TargetVE, c.TargetVU)
		m.ack(id, c, ok, cond(ok, "target assigned; ready for START_INTERCEPT", "rejected: cannot assign after launch"))
		if ok {
			log.Printf("[int] %s: ASSIGN_TARGET mission=%s target=%s @ %.4f,%.4f", id, c.MissionID, c.TargetID, c.TargetLat, c.TargetLon)
		}
	case "START_INTERCEPT":
		ok := e.startIntercept()
		m.ack(id, c, ok, cond(ok, "START_INTERCEPT accepted; interceptor launched", "rejected: no target or already launched"))
		if ok {
			log.Printf("[int] %s: LAUNCH → midcourse toward %s", id, e.TargetID)
		}
	case "ABORT", "HIT":
		e.abort()
		m.ack(id, c, true, "ABORT accepted; safing")
		log.Printf("[int] %s: ABORT", id)
	default:
		m.ack(id, c, false, "unsupported command")
	}
}

// handleTargetState refines any engagement locked onto this entity id.
func (m *manager) handleTargetState(msg *nats.Msg) {
	// subject: ent.air.{id}.state
	parts := strings.Split(msg.Subject, ".")
	if len(parts) < 4 {
		return
	}
	tgtID := strings.Join(parts[2:len(parts)-1], ".")
	var ent struct {
		Position *struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Altitude  float64 `json:"altitude"`
		} `json:"position"`
	}
	if json.Unmarshal(msg.Data, &ent) != nil || ent.Position == nil {
		return
	}
	m.mu.Lock()
	var targets []*Engagement
	for _, e := range m.eng {
		if e.TargetID == tgtID && e.Launched {
			targets = append(targets, e)
		}
	}
	m.mu.Unlock()
	for _, e := range targets {
		e.updateTarget(ent.Position.Latitude, ent.Position.Longitude, ent.Position.Altitude)
	}
}

func (m *manager) simLoop() {
	dt := 1.0 / float64(simHz)
	tick := time.NewTicker(time.Duration(float64(time.Second) / float64(simHz)))
	defer tick.Stop()
	var n uint64
	for range tick.C {
		n++
		m.mu.Lock()
		engs := make([]*Engagement, 0, len(m.eng))
		for _, e := range m.eng {
			engs = append(engs, e)
		}
		m.mu.Unlock()

		for _, e := range engs {
			done := e.step(dt)
			snap := e.snapshot()

			// Terminal: finalize exactly once, then retire the engagement so a
			// spent interceptor stops broadcasting (otherwise it lingers on the
			// COP/ATAK forever). The final phase + result are published, and the
			// interceptor track is removed from the world-state.
			if done || snap.Phase.terminal() {
				m.publishPhase(e, snap) // final phase (SUCCESS/FAILED) → int.phase + KV
				m.publishResult(e, snap)
				if snap.Launched {
					// Drop the expended interceptor track from spine → COP/ATAK.
					m.nc.Publish("ent.manual.delete", []byte(e.ID))
				}
				m.mu.Lock()
				delete(m.eng, e.ID)
				delete(m.last, e.ID)
				m.mu.Unlock()
				continue
			}

			// Active flight: own-force telemetry at 10 Hz, phase on change (+1 Hz keepalive).
			if n%2 == 0 {
				m.publishTelemetry(e)
			}
			m.mu.Lock()
			prev, seen := m.last[e.ID]
			m.last[e.ID] = snap.Phase
			m.mu.Unlock()
			if !seen || prev != snap.Phase || n%uint64(simHz) == 0 {
				m.publishPhase(e, snap)
			}
		}
	}
}

// publishTelemetry emits the interceptor's own position as friendly telemetry
// (spine ingests tel.drone.* → COP/ATAK sees it) plus the internal target estimate.
func (m *manager) publishTelemetry(e *Engagement) {
	snap := e.snapshot()
	if !snap.Launched || snap.Phase.terminal() {
		return // pre-launch or spent — nothing to broadcast
	}
	// tel.drone.{id} — friendly own-force track (stable id ⇒ no proximity merge).
	tel := map[string]any{
		"id":          e.ID,
		"entity_type": "aircraft_multirotor",
		"affiliation": "friendly",
		"name":        "Interceptor " + e.ID,
		"latitude":    snap.Lat,
		"longitude":   snap.Lon,
		"altitude":    snap.Alt,
		"heading":     snap.Heading,
		"velocity":    snap.Speed,
		"confidence":  1.0,
		"source":      "interceptor",
	}
	if b, err := json.Marshal(tel); err == nil {
		_ = m.nc.Publish("tel.drone."+e.ID, b)
	}
	// int.target.{id} — InternalTargetEstimate (range/closing).
	est := map[string]any{
		"stamp": time.Now().UTC().Format(time.RFC3339Nano), "interceptor_id": e.ID,
		"target_id": e.TargetID, "valid": snap.Range > 0,
		"range_m": round1(snap.Range), "closing_speed_mps": round1(snap.Closing),
		"confidence": cond2(snap.SeekerLck, 0.9, cond2(snap.SeekerDet, 0.75, 0.5)),
	}
	if b, err := json.Marshal(est); err == nil {
		_ = m.nc.Publish("int.target."+e.ID, b)
	}
}

func (m *manager) publishPhase(e *Engagement, snap snapshot) {
	ph := map[string]any{
		"stamp": time.Now().UTC().Format(time.RFC3339Nano), "interceptor_id": e.ID,
		"mission_id": e.MissionID, "target_id": e.TargetID,
		"phase": uint8(snap.Phase), "phase_name": snap.Phase.String(),
		"authority": uint8(snap.Authority), "launched": snap.Launched,
		"terminal_phase":  snap.Phase == PhaseTerminalHoming,
		"seeker_detected": snap.SeekerDet, "seeker_locked": snap.SeekerLck,
		"range_m": round1(snap.Range), "closing_speed_mps": round1(snap.Closing),
		"latitude": snap.Lat, "longitude": snap.Lon, "altitude": round1(snap.Alt),
		"speed_mps": round1(snap.Speed), "heading_deg": round1(snap.Heading),
	}
	b, _ := json.Marshal(ph)
	_ = m.nc.Publish("int.phase."+e.ID, b)
	_, _ = m.kv.Put(e.ID, b)
}

func (m *manager) publishResult(e *Engagement, snap snapshot) {
	outcome := "FAILED"
	if snap.Phase == PhaseInterceptSuccess {
		outcome = "SUCCESS"
	} else if snap.Phase == PhaseIdle {
		outcome = "RTB" // returned to base after abort
	}
	res := map[string]any{
		"stamp": time.Now().UTC().Format(time.RFC3339Nano), "interceptor_id": e.ID,
		"mission_id": e.MissionID, "target_id": e.TargetID,
		"outcome": outcome, "min_range_m": round1(snap.MinRange), "health": snap.Health,
	}
	b, _ := json.Marshal(res)
	_ = m.nc.Publish("int.result."+e.ID, b)
	log.Printf("[int] %s: RESULT %s (min-range %.1f m) vs %s", e.ID, outcome, snap.MinRange, e.TargetID)
}

func (m *manager) ack(id string, c command, ok bool, message string) {
	ack := map[string]any{
		"stamp": time.Now().UTC().Format(time.RFC3339Nano), "interceptor_id": id,
		"mission_id": c.MissionID, "command": c.Type, "accepted": ok, "message": message,
	}
	b, _ := json.Marshal(ack)
	_ = m.nc.Publish("int.ack."+id, b)
}

func (m *manager) heartbeat() {
	tick := time.NewTicker(15 * time.Second)
	defer tick.Stop()
	for {
		m.mu.Lock()
		active := 0
		for _, e := range m.eng {
			s := e.snapshot()
			if s.Launched && !s.Phase.terminal() {
				active++
			}
		}
		total := len(m.eng)
		m.mu.Unlock()
		hb := map[string]any{"organ": "interceptor",
			"status": cond(active > 0, "engaging", "ready"),
			"active": active, "total": total}
		b, _ := json.Marshal(hb)
		_ = m.nc.Publish(hbSubject, b)
		<-tick.C
	}
}

func cond(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}
func cond2(b bool, t, f float64) float64 {
	if b {
		return t
	}
	return f
}
func round1(f float64) float64 { return float64(int64(f*10+0.5)) / 10 }
