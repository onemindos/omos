package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
)

// om-overwatch-tak — the C2 bridge between the fusion spine and TAK/ATAK.
//
// The TAK Server's Java NatsBridgePlugin already owns all TAK transport (TLS
// :8089, certs, the TAK protocol) and exposes CoT as plain XML on NATS:
//   tak.cot.out    — every CoT the TAK server sees (ATAK devices, operators)
//   tak.cot.inject — CoT we publish → pushed to every ATAK client
//
// This organ is a pure translation + fusion layer:
//
//   INBOUND  tak.cot.out ──parse CoT──► det.cot.{id} ──► spine fuses
//            Human-tagged affiliation (a-f/h/n/u) becomes authoritative.
//
//   OUTBOUND ent.air.{id}.state ──render CoT──► tak.cot.inject ──► ATAK
//            Operators finally see the full fused sensor picture (drone
//            detections, friendly UAV telemetry, etc.) on their tablets.
//
// Loop prevention: CoT-sourced entities are never re-emitted to TAK (ATAK
// already has them), and inbound skips our own "OW." re-broadcasts.

const (
	subjCoTOut    = "tak.cot.out"      // ATAK → us
	subjCoTInject = "tak.cot.inject"   // us → ATAK
	subjCoTIngest = "det.cot."         // us → spine (per-entity: det.cot.{id})
	subjEntState  = "ent.air.*.state"    // spine → us
	subjEntDelete = "ent.air.*.deleted"  // spine → us (entity removed)
	hbSubject     = "agents.hb.overwatch.tak.zeus"
	cotStaleFor   = 90 * time.Second
)

type bridge struct {
	nc          *nats.Conn
	inbound     atomic.Uint64
	outbound    atomic.Uint64
	skippedIn   atomic.Uint64
	skippedOut  atomic.Uint64
}

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = "nats://127.0.0.1:4222"
	}
	nc, err := nats.Connect(url, nats.Name("om-overwatch-tak"),
		nats.MaxReconnects(-1), nats.ReconnectWait(2*time.Second))
	if err != nil {
		log.Fatal("NATS connect:", err)
	}
	defer nc.Close()
	log.Printf("om-overwatch-tak: connected to OneMind bus at %s", url)

	b := &bridge{nc: nc}

	// INBOUND: ATAK CoT → spine
	if _, err := nc.Subscribe(subjCoTOut, b.handleInbound); err != nil {
		log.Fatal("subscribe tak.cot.out:", err)
	}
	// OUTBOUND: fused entity state → ATAK
	if _, err := nc.Subscribe(subjEntState, b.handleOutbound); err != nil {
		log.Fatal("subscribe ent.air.*.state:", err)
	}
	// OUTBOUND: entity deletions → ATAK tombstone (remove the marker promptly)
	if _, err := nc.Subscribe(subjEntDelete, b.handleDelete); err != nil {
		log.Fatal("subscribe ent.air.*.deleted:", err)
	}
	log.Printf("om-overwatch-tak: bridge active (IN %s → det.cot.{id}, OUT %s → %s)",
		subjCoTOut, subjEntState, subjCoTInject)

	// Heartbeat + periodic stats.
	nc.Publish(hbSubject, []byte(`{"status":"up"}`))
	ticker := time.NewTicker(15 * time.Second)
	for range ticker.C {
		log.Printf("om-overwatch-tak: alive | in=%d (skip %d) out=%d (skip %d)",
			b.inbound.Load(), b.skippedIn.Load(), b.outbound.Load(), b.skippedOut.Load())
		nc.Publish(hbSubject, []byte(`{"t":"`+time.Now().UTC().Format(time.RFC3339)+
			`","status":"bridging"}`))
	}
}

// handleInbound parses a CoT event from ATAK and republishes it to the spine
// as a det.cot.{id} contact report carrying human-asserted affiliation.
func (b *bridge) handleInbound(m *nats.Msg) {
	ev, err := ParseCoT(m.Data)
	if err != nil {
		b.skippedIn.Add(1)
		return
	}
	// Ignore our own re-broadcasts — no feedback loop.
	if ev.IsOverwatchOrigin() {
		b.skippedIn.Add(1)
		return
	}
	id, report := ev.ToContactReport()
	if id == "" {
		b.skippedIn.Add(1)
		return
	}
	data, err := json.Marshal(report)
	if err != nil {
		return
	}
	if err := b.nc.Publish(subjCoTIngest+id, data); err != nil {
		log.Printf("publish det.cot.%s: %v", id, err)
		return
	}
	n := b.inbound.Add(1)
	if n%20 == 1 {
		log.Printf("om-overwatch-tak: IN  %s → det.cot.%s (%s, %s) @ %.4f,%.4f",
			ev.UID, id, report.EntityType, report.Affiliation, ev.Point.Lat, ev.Point.Lon)
	}
}

// handleOutbound renders a fused entity as CoT and injects it into TAK so it
// appears on every ATAK client. CoT-sourced tracks are not echoed back.
func (b *bridge) handleOutbound(m *nats.Msg) {
	var e Entity
	if json.Unmarshal(m.Data, &e) != nil || e.EntityID == "" {
		return
	}
	// Never echo a track that originated from TAK — ATAK already displays it.
	if e.hasSource("cot") {
		b.skippedOut.Add(1)
		return
	}
	// Only broadcast tracks with a position (nothing to plot otherwise).
	if e.Position == nil {
		b.skippedOut.Add(1)
		return
	}
	cot, err := RenderCoT(&e, cotStaleFor)
	if err != nil {
		log.Printf("render cot %s: %v", e.EntityID, err)
		return
	}
	if err := b.nc.Publish(subjCoTInject, cot); err != nil {
		log.Printf("publish tak.cot.inject: %v", err)
		return
	}
	n := b.outbound.Add(1)
	if n%20 == 1 {
		log.Printf("om-overwatch-tak: OUT %s (%s, %s, src=%s) → ATAK @ %.4f,%.4f",
			e.EntityID, e.EntityType, e.Affiliation, strings.Join(e.Sources, "+"),
			e.Position.Latitude, e.Position.Longitude)
	}
}

// handleDelete turns a spine entity-deletion into an ATAK tombstone so the
// marker is removed from every client's map promptly (rather than lingering
// until its stale window expires). Fired e.g. when a spent interceptor retires.
func (b *bridge) handleDelete(m *nats.Msg) {
	var d struct {
		EntityID string `json:"entity_id"`
	}
	if json.Unmarshal(m.Data, &d) != nil || d.EntityID == "" {
		return
	}
	cot, err := RenderCoTDelete(d.EntityID)
	if err != nil {
		log.Printf("render delete cot %s: %v", d.EntityID, err)
		return
	}
	if err := b.nc.Publish(subjCoTInject, cot); err != nil {
		log.Printf("publish delete tak.cot.inject: %v", err)
		return
	}
	log.Printf("om-overwatch-tak: DEL %s → ATAK tombstone", d.EntityID)
}
