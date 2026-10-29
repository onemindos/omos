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
// GeoChat on the "All Chat Rooms" group (224.10.10.1:17012) travels the same
// subjects. Optional Commo-compatible mesh encryption (AES-256-CBC with
// HMAC-SHA256) and peer-hosted data package transfer (/getfile) are supported.
//
// With MESH_TO_TAK_SERVER=true (default false) it also copies received CoT to
// tak.cot.inject, so mesh devices appear on TAK Server, except OneMind's own
// output (uids OM.* and INTEL.*) and gateway presence (omos-gw-*).
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

const buildVersion = "0.3.0"

const (
	subjectRx     = "tak.cot.mesh.rx"
	subjectTx     = "tak.cot.mesh.tx"
	subjectInject = "tak.cot.inject"

	// headerDestUID on a tak.cot.mesh.tx message sends it as a directed
	// message to that contact's endpoint instead of broadcasting it.
	headerDestUID = "Tak-Dest-Uid"

	gatewayUIDPrefix = "omos-gw-"
)

// ── Config ────────────────────────────────────────────────────────────────────

type config struct {
	natsURL, natsUser, natsPassword string
	gatewayID                       string
	meshGroup                       string
	chatGroup                       string
	meshIface                       string
	meshTTL                         int
	authKey, cipherKey              string
	forceXML                        bool
	directPort                      int
	udpPort                         int
	toTakServer                     bool
	callsign, cotType, endpointHost string
	lat, lon                        float64
	datapkgDir                      string
	datapkgPort                     int
	datapkgMax                      int64
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
		chatGroup:    getenv("MESH_CHAT_GROUP", "224.10.10.1:17012"),
		meshIface:    os.Getenv("MESH_INTERFACE"),
		meshTTL:      mustInt("MESH_TTL", "1"),
		authKey:      os.Getenv("MESH_AUTH_KEY"),
		cipherKey:    os.Getenv("MESH_CIPHER_KEY"),
		directPort:   mustInt("DIRECT_TCP_PORT", "4242"),
		udpPort:      mustInt("UDP_LISTEN_PORT", "0"),
		toTakServer:  mustBool("MESH_TO_TAK_SERVER", "false"),
		callsign:     os.Getenv("GATEWAY_CALLSIGN"),
		cotType:      getenv("GATEWAY_COT_TYPE", "a-f-G-E-S"),
		endpointHost: os.Getenv("GATEWAY_ENDPOINT_HOST"),
		lat:          mustFloat("GATEWAY_LAT", "0"),
		lon:          mustFloat("GATEWAY_LON", "0"),
		datapkgDir:   os.Getenv("DATAPKG_DIR"),
		datapkgPort:  mustInt("DATAPKG_HTTP_PORT", "0"),
		datapkgMax:   int64(mustInt("DATAPKG_MAX_BYTES", "52428800")),
	}
	if c.chatGroup == "off" {
		c.chatGroup = ""
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

// injectAllowed implements CP-2 (cot-projection contract §7): with
// MESH_TO_TAK_SERVER=true, the copy to tak.cot.inject skips OneMind's own
// projected CoT (OM.*, INTEL.*), which must pass the hub's gates instead, and
// the presence SA of this or any other OneMind gateway.
func injectAllowed(ev *Event) bool {
	return !strings.HasPrefix(ev.UID, "OM.") && !strings.HasPrefix(ev.UID, "INTEL.") &&
		!strings.HasPrefix(ev.UID, gatewayUIDPrefix)
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	cfg := loadConfig()
	id := subjectToken(cfg.gatewayID)
	uid := gatewayUIDPrefix + id
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
	var chat *net.UDPAddr
	if cfg.chatGroup != "" {
		if chat, err = net.ResolveUDPAddr("udp4", cfg.chatGroup); err != nil {
			log.Fatalf("[%s] MESH_CHAT_GROUP: %v", id, err)
		}
	}
	var crypto *MeshCrypto
	if cfg.authKey != "" || cfg.cipherKey != "" {
		if crypto, err = NewMeshCrypto(cfg.authKey, cfg.cipherKey); err != nil {
			log.Fatalf("[%s] MESH_AUTH_KEY/MESH_CIPHER_KEY: %v", id, err)
		}
		log.Printf("[%s] mesh encryption on: plaintext mesh and directed messages are dropped", id)
	}
	mesh, err := NewMesh(MeshConfig{
		Group: group, ChatGroup: chat, Interface: iface, TTL: cfg.meshTTL,
		ContactUID: uid, ForceXML: cfg.forceXML, Crypto: crypto,
	})
	if err != nil {
		log.Fatalf("[%s] %v", id, err)
	}
	txIface := "system default"
	if ti, err := mesh.TransmitInterface(); err == nil && ti != nil {
		txIface = ti.Name
	}
	log.Printf("[%s] mesh %s, chat %s, transmit interface %s", id, cfg.meshGroup, cfg.chatGroup, txIface)

	var rx, tx, errCount atomic.Int64
	var lastErr atomic.Value
	lastErr.Store("")
	onErr := func(err error) {
		errCount.Add(1)
		lastErr.Store(err.Error())
		log.Printf("[%s] %v", id, err)
	}

	// ── Data packages ────────────────────────────────────────────────────────
	var dp *DataPackages
	if cfg.datapkgDir != "" {
		dp = &DataPackages{Dir: cfg.datapkgDir, MaxBytes: cfg.datapkgMax, SelfUID: uid, OnErr: onErr,
			SendAck: func(ack *Event, senderUID, from string) error {
				// Commo routes the ack like any CoT to the sender: directed to
				// its advertised endpoint, else its source address on 4242.
				ep := mesh.Endpoint(senderUID)
				if ep == "" {
					if host, _, err := net.SplitHostPort(from); err == nil {
						ep = net.JoinHostPort(host, "4242")
					}
				}
				if err := mesh.SendDirect(ack, ep); err == nil {
					return nil
				}
				return mesh.Send(ack)
			}}
		if err := dp.Init(); err != nil {
			log.Fatalf("[%s] DATAPKG_DIR: %v", id, err)
		}
		log.Printf("[%s] data packages: inbox %s/inbox, outbox %s/outbox", id, cfg.datapkgDir, cfg.datapkgDir)
	}
	var httpSrv *http.Server
	if dp != nil && cfg.datapkgPort > 0 {
		httpSrv = &http.Server{Addr: fmt.Sprintf(":%d", cfg.datapkgPort), Handler: dp,
			ReadHeaderTimeout: 10 * time.Second}
		go func() {
			if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("[%s] DATAPKG_HTTP_PORT: %v", id, err)
			}
		}()
		log.Printf("[%s] serving outbox on http :%d/getfile", id, cfg.datapkgPort)
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
		if cfg.toTakServer && injectAllowed(r.Event) {
			if err := nc.Publish(subjectInject, xmlBytes); err != nil {
				onErr(fmt.Errorf("publish %s: %w", subjectInject, err))
			}
		}
		rx.Add(1)
		if dp != nil {
			dp.Handle(r)
		}
	}

	go mesh.Run(publish, onErr)

	var directLn net.Listener
	if cfg.directPort > 0 {
		if directLn, err = net.Listen("tcp4", fmt.Sprintf(":%d", cfg.directPort)); err != nil {
			log.Fatalf("[%s] direct TCP :%d: %v", id, cfg.directPort, err)
		}
		observe := func(ev *Event, ctl *TakControl, v uint64) { mesh.observe(ev, ctl, v, time.Now()) }
		go ServeDirect(directLn, mesh.DecodeInbound, observe, publish, onErr)
		log.Printf("[%s] directed messages on TCP :%d", id, cfg.directPort)
	}

	var udpConn *net.UDPConn
	if cfg.udpPort > 0 {
		if udpConn, err = net.ListenUDP("udp4", &net.UDPAddr{Port: cfg.udpPort}); err != nil {
			log.Fatalf("[%s] UDP :%d: %v", id, cfg.udpPort, err)
		}
		go ServeUDP(udpConn, mesh.DecodeInbound, publish, onErr)
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
		send := func() error { return mesh.Send(ev) }
		if dest := m.Header.Get(headerDestUID); dest != "" {
			send = func() error {
				ep := mesh.Endpoint(dest)
				if ep == "" {
					return fmt.Errorf("no endpoint known for contact %s", dest)
				}
				return mesh.SendDirect(ev, ep)
			}
		}
		if err := send(); err != nil {
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
			meta := map[string]any{
				"gateway": id, "version": buildVersion, "mesh_group": cfg.meshGroup,
				"chat_group": cfg.chatGroup, "tx_interface": txIface, "encrypted": crypto != nil,
				"to_tak_server": cfg.toTakServer,
				"proto_version": mesh.Version(), "contacts": mesh.Contacts(),
				"rx": rx.Load(), "tx": tx.Load(), "errors": errCount.Load(), "last_error": lastErr.Load(),
			}
			if dp != nil {
				meta["datapkg_rx"] = dp.Received.Load()
				meta["datapkg_failed"] = dp.Failed.Load()
			}
			beat, _ := json.Marshal(map[string]any{"uid": uid, "type": "fleet.edge", "t": now.UnixMilli(), "meta": meta})
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
	if httpSrv != nil {
		_ = httpSrv.Close()
	}
	if err := nc.Drain(); err != nil {
		log.Printf("[%s] NATS drain: %v", id, err)
	}
}
