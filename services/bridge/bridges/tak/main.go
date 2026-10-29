// SPDX-License-Identifier: Apache-2.0

// Package main is the OneMind TAK mesh gateway. It joins a TAK mesh the way
// ATAK, iTAK, WinTAK and VR-TAK do through Common Commo: multicast SA, TCP
// directed messages, and TAK Protocol version negotiation (XML v0 or
// protobuf v1). Everything it hears goes onto NATS as CoT XML; anything
// published to the transmit subject goes out on the mesh.
//
//	mesh  ──► tak.cot.mesh.rx              (stream COT, deduplicated across gateways)
//	mesh  ◄── tak.cot.mesh.tx              (CoT XML; request/reply gets "ok" or an error)
//	          tel.fleet.edge.{id}.status   (gateway health every 30s)
//
// With MESH_TO_TAK_SERVER=true it also copies received CoT to tak.cot.inject,
// so mesh devices appear on TAK Server and every client connected to it.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

const buildVersion = "0.2.0"

const (
	subjectRx     = "tak.cot.mesh.rx"
	subjectTx     = "tak.cot.mesh.tx"
	subjectInject = "tak.cot.inject"
)

// ── Config ────────────────────────────────────────────────────────────────────

type config struct {
	natsURL, natsUser, natsPassword string
	gatewayID                       string
	meshGroup                       string
	meshIface                       string
	forceXML                        bool
	directPort                      int
	udpPort                         int
	toTakServer                     bool
	callsign, cotType, endpointHost string
	lat, lon                        float64
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustInt(key, def string) int {
	v, err := strconv.Atoi(getenv(key, def))
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return v
}

func mustFloat(key, def string) float64 {
	v, err := strconv.ParseFloat(getenv(key, def), 64)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return v
}

func mustBool(key, def string) bool {
	v, err := strconv.ParseBool(getenv(key, def))
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return v
}

func loadConfig() config {
	c := config{
		natsURL:      getenv("NATS_URL", "nats://localhost:4222"),
		natsUser:     os.Getenv("NATS_USER"),
		natsPassword: os.Getenv("NATS_PASSWORD"),
		gatewayID:    getenv("GATEWAY_ID", "tak-mesh"),
		meshGroup:    getenv("MESH_GROUP", "239.2.3.1:6969"),
		meshIface:    os.Getenv("MESH_INTERFACE"),
		directPort:   mustInt("DIRECT_TCP_PORT", "4242"),
		udpPort:      mustInt("UDP_LISTEN_PORT", "0"),
		toTakServer:  mustBool("MESH_TO_TAK_SERVER", "false"),
		callsign:     os.Getenv("GATEWAY_CALLSIGN"),
		cotType:      getenv("GATEWAY_COT_TYPE", "a-f-G-E-S"),
		endpointHost: os.Getenv("GATEWAY_ENDPOINT_HOST"),
		lat:          mustFloat("GATEWAY_LAT", "0"),
		lon:          mustFloat("GATEWAY_LON", "0"),
	}
	switch p := getenv("MESH_PROTOCOL", "auto"); p {
	case "auto":
	case "xml":
		c.forceXML = true
	default:
		log.Fatalf("invalid MESH_PROTOCOL %q (auto or xml)", p)
	}
	return c
}

var unsafeChars = regexp.MustCompile(`[^a-zA-Z0-9\-_]`)

// subjectToken makes s safe to use as one NATS subject token.
func subjectToken(s string) string { return unsafeChars.ReplaceAllString(s, "_") }

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	cfg := loadConfig()
	id := subjectToken(cfg.gatewayID)
	uid := "omos-gw-" + id
	log.Printf("[%s] starting v%s: mesh %s, NATS %s", id, buildVersion, cfg.meshGroup, cfg.natsURL)

	opts := []nats.Option{
		nats.Name("tak-mesh-" + id),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) { log.Printf("[%s] NATS disconnected: %v", id, err) }),
		nats.ReconnectHandler(func(nc *nats.Conn) { log.Printf("[%s] NATS reconnected to %s", id, nc.ConnectedUrl()) }),
	}
	if cfg.natsUser != "" {
		opts = append(opts, nats.UserInfo(cfg.natsUser, cfg.natsPassword))
	}
	nc, err := nats.Connect(cfg.natsURL, opts...)
	if err != nil {
		log.Fatalf("[%s] NATS connect failed: %v", id, err)
	}

	group, err := net.ResolveUDPAddr("udp4", cfg.meshGroup)
	if err != nil {
		log.Fatalf("[%s] MESH_GROUP: %v", id, err)
	}
	var iface *net.Interface
	if cfg.meshIface != "" {
		if iface, err = net.InterfaceByName(cfg.meshIface); err != nil {
			log.Fatalf("[%s] MESH_INTERFACE: %v", id, err)
		}
	}
	mesh, err := NewMesh(MeshConfig{Group: group, Interface: iface, ContactUID: uid, ForceXML: cfg.forceXML})
	if err != nil {
		log.Fatalf("[%s] %v", id, err)
	}

	var rx, tx, errCount atomic.Int64
	var lastErr atomic.Value
	lastErr.Store("")
	onErr := func(err error) {
		errCount.Add(1)
		lastErr.Store(err.Error())
		log.Printf("[%s] %v", id, err)
	}

	// ── Mesh → NATS ──────────────────────────────────────────────────────────
	publish := func(r Received) {
		xmlBytes := r.Event.MarshalCotXML(false)
		msg := nats.NewMsg(subjectRx)
		msg.Data = xmlBytes
		msg.Header.Set(nats.MsgIdHdr, dedupID(r.Event))
		msg.Header.Set("Tak-Gateway", id)
		msg.Header.Set("Tak-Delivery", string(r.Delivery))
		msg.Header.Set("Tak-Proto", strconv.FormatUint(r.Version, 10))
		msg.Header.Set("Tak-From", r.From)
		if r.Event.Detail.Extensions > 0 {
			// takproto README: tell consumers the message was not fully decoded.
			msg.Header.Set("Tak-Undecoded-Extensions", strconv.Itoa(r.Event.Detail.Extensions))
		}
		if err := nc.PublishMsg(msg); err != nil {
			onErr(fmt.Errorf("publish %s: %w", subjectRx, err))
			return
		}
		if cfg.toTakServer {
			if err := nc.Publish(subjectInject, xmlBytes); err != nil {
				onErr(fmt.Errorf("publish %s: %w", subjectInject, err))
			}
		}
		rx.Add(1)
	}

	go mesh.Run(publish, onErr)

	var directLn net.Listener
	if cfg.directPort > 0 {
		if directLn, err = net.Listen("tcp4", fmt.Sprintf(":%d", cfg.directPort)); err != nil {
			log.Fatalf("[%s] direct TCP :%d: %v", id, cfg.directPort, err)
		}
		observe := func(ev *Event, ctl *TakControl, v uint64) { mesh.observe(ev, ctl, v, time.Now()) }
		go ServeDirect(directLn, observe, publish, onErr)
		log.Printf("[%s] directed messages on TCP :%d", id, cfg.directPort)
	}

	var udpConn *net.UDPConn
	if cfg.udpPort > 0 {
		if udpConn, err = net.ListenUDP("udp4", &net.UDPAddr{Port: cfg.udpPort}); err != nil {
			log.Fatalf("[%s] UDP :%d: %v", id, cfg.udpPort, err)
		}
		go ServeUDP(udpConn, publish, onErr)
		log.Printf("[%s] unicast CoT on UDP :%d", id, cfg.udpPort)
	}

	// ── NATS → mesh ──────────────────────────────────────────────────────────
	sub, err := nc.Subscribe(subjectTx, func(m *nats.Msg) {
		reply := func(s string) {
			if m.Reply != "" {
				_ = m.Respond([]byte(s))
			}
		}
		ev, err := ParseCotXML(m.Data)
		if err != nil {
			onErr(fmt.Errorf("%s: %w", subjectTx, err))
			reply("error: " + err.Error())
			return
		}
		if err := mesh.Send(ev); err != nil {
			onErr(fmt.Errorf("mesh send uid=%s: %w", ev.UID, err))
			reply("error: " + err.Error())
			return
		}
		tx.Add(1)
		reply("ok")
	})
	if err != nil {
		log.Fatalf("[%s] subscribe %s: %v", id, subjectTx, err)
	}

	// ── Presence and health ──────────────────────────────────────────────────
	endpoint := ""
	if cfg.endpointHost != "" && cfg.directPort > 0 {
		endpoint = fmt.Sprintf("%s:%d:tcp", cfg.endpointHost, cfg.directPort)
	}
	statusSubject := fmt.Sprintf("tel.fleet.edge.%s.status", id)
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for now := range ticker.C {
			if cfg.callsign != "" {
				if err := mesh.Send(presence(uid, cfg.callsign, cfg.cotType, endpoint, cfg.lat, cfg.lon, now)); err != nil {
					onErr(fmt.Errorf("presence: %w", err))
				}
			}
			beat, _ := json.Marshal(map[string]any{
				"uid": uid, "type": "fleet.edge", "t": now.UnixMilli(),
				"meta": map[string]any{
					"gateway": id, "version": buildVersion, "mesh_group": cfg.meshGroup,
					"proto_version": mesh.Version(), "contacts": mesh.Contacts(),
					"rx": rx.Load(), "tx": tx.Load(), "errors": errCount.Load(), "last_error": lastErr.Load(),
				},
			})
			if err := nc.Publish(statusSubject, beat); err != nil {
				log.Printf("[%s] status publish: %v", id, err)
			}
		}
	}()
	log.Printf("[%s] ready: %s -> mesh, mesh -> %s", id, subjectTx, subjectRx)

	// ── Shutdown ─────────────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	log.Printf("[%s] received %s, shutting down", id, <-quit)
	ticker.Stop()
	_ = sub.Unsubscribe()
	mesh.Close()
	if directLn != nil {
		directLn.Close()
	}
	if udpConn != nil {
		udpConn.Close()
	}
	if err := nc.Drain(); err != nil {
		log.Printf("[%s] NATS drain: %v", id, err)
	}
}
