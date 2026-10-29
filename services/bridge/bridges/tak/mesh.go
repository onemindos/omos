// SPDX-License-Identifier: Apache-2.0

package main

// Mesh networking as TAK clients do it (takproto/README.txt, "Mesh Network
// Protocol Negotiation"): SA and broadcasts go to a UDP multicast group, one
// message per datagram; directed messages arrive as TCP connect-send-close.
// The gateway tracks the protocol versions of every contact it hears and
// broadcasts in the highest version all of them can decode.

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	ourMaxVersion   = 1
	controlInterval = 60 * time.Second // rule 1: TakControl at least every 60s
	controlTimeout  = 2 * time.Minute  // rule 3c
	contactTimeout  = 5 * time.Minute  // forget contacts we no longer hear
	echoWindow      = 10 * time.Second // drop our own datagrams looped back
	maxDirectBytes  = 1 << 20          // largest directed TCP message accepted
	directTimeout   = 10 * time.Second
)

// Delivery says how a message reached the gateway.
type Delivery string

const (
	DeliveryMesh   Delivery = "mesh"   // multicast group
	DeliveryDirect Delivery = "direct" // TCP send-once to the gateway
	DeliveryUDP    Delivery = "udp"    // unicast UDP input
)

// Received is one decoded message heard by the gateway.
type Received struct {
	Event    *Event
	Version  uint64 // 0 = XML, 1 = protobuf
	Delivery Delivery
	From     string
}

type contact struct {
	maxVersion  uint32 // highest version this contact can decode
	lastControl time.Time
	lastVersion uint32 // version of the most recent message
	lastSeen    time.Time
}

// MeshConfig configures a Mesh.
type MeshConfig struct {
	Group      *net.UDPAddr   // SA multicast group, e.g. 239.2.3.1:6969
	Interface  *net.Interface // nil = system default
	ContactUID string         // our identity in TakControl
	ForceXML   bool           // never send protobuf (legacy-only meshes)
}

// Mesh is a TAK mesh participant.
type Mesh struct {
	cfg  MeshConfig
	recv *net.UDPConn
	send *net.UDPConn

	mu       sync.Mutex
	contacts map[string]*contact
	sent     map[[32]byte]time.Time
	version  uint32 // last computed broadcast version
}

// NewMesh joins the multicast group.
func NewMesh(cfg MeshConfig) (*Mesh, error) {
	recv, err := net.ListenMulticastUDP("udp4", cfg.Interface, cfg.Group)
	if err != nil {
		return nil, fmt.Errorf("join %s: %w", cfg.Group, err)
	}
	_ = recv.SetReadBuffer(1 << 20)
	send, err := net.DialUDP("udp4", nil, cfg.Group)
	if err != nil {
		recv.Close()
		return nil, fmt.Errorf("dial %s: %w", cfg.Group, err)
	}
	m := &Mesh{cfg: cfg, recv: recv, send: send, contacts: map[string]*contact{}, sent: map[[32]byte]time.Time{}}
	m.version = m.computeVersion(time.Now())
	return m, nil
}

func (m *Mesh) Close() {
	m.recv.Close()
	m.send.Close()
}

// Decode turns one datagram or directed payload into an Event. ctl is set
// when the message carried a TakControl. ev is nil for control-only messages.
func Decode(b []byte) (ev *Event, ctl *TakControl, version uint64, err error) {
	if !IsTakProtocol(b) {
		ev, err = ParseCotXML(b)
		return ev, nil, 0, err
	}
	version, payload, err := DecodeMeshHeader(b)
	if err != nil {
		return nil, nil, 0, err
	}
	if version != 1 {
		return nil, nil, version, fmt.Errorf("unsupported TAK protocol version %d", version)
	}
	msg, err := UnmarshalTakMessage(payload)
	if err != nil {
		return nil, nil, version, err
	}
	return msg.Event, msg.Control, version, nil
}

// observe applies negotiation rules 3a-3b for one received message.
func (m *Mesh) observe(ev *Event, ctl *TakControl, version uint64, now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	update := func(uid string) *contact {
		c := m.contacts[uid]
		if c == nil {
			c = &contact{maxVersion: uint32(version)} // 3a
			m.contacts[uid] = c
		}
		c.lastVersion = uint32(version)
		c.lastSeen = now
		return c
	}
	if ctl != nil && ctl.ContactUID != "" && ctl.ContactUID != m.cfg.ContactUID {
		c := update(ctl.ContactUID)
		c.maxVersion = ctl.MaxProtoVersion // 3b
		c.lastControl = now
	}
	// SA from a mesh participant carries a contact endpoint. Markers and other
	// relayed events do not identify a client, so they are not tracked.
	if ev != nil && ev.UID != m.cfg.ContactUID && ev.Detail.Contact != nil && ev.Detail.Contact.Endpoint != "" {
		if ctl == nil || ctl.ContactUID != ev.UID {
			update(ev.UID)
		}
	}
}

// computeVersion applies rules 3c and 4. Caller holds mu or owns m.
func (m *Mesh) computeVersion(now time.Time) uint32 {
	if m.cfg.ForceXML {
		return 0
	}
	v := uint32(ourMaxVersion)
	for uid, c := range m.contacts {
		if now.Sub(c.lastSeen) > contactTimeout {
			delete(m.contacts, uid)
			continue
		}
		max := c.maxVersion
		if now.Sub(c.lastControl) > controlTimeout {
			max = c.lastVersion // 3c
		}
		if max < v {
			v = max
		}
	}
	return v
}

// Version is the protocol version broadcasts currently use.
func (m *Mesh) Version() uint32 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.version
}

// Contacts is the number of mesh clients currently tracked.
func (m *Mesh) Contacts() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.contacts)
}

// recheck recomputes the broadcast version and, per rule 5, announces a change.
func (m *Mesh) recheck(now time.Time) {
	m.mu.Lock()
	v := m.computeVersion(now)
	changed := v != m.version
	if changed {
		log.Printf("[mesh] broadcast protocol version %d -> %d (%d contacts)", m.version, v, len(m.contacts))
	}
	m.version = v
	m.mu.Unlock()
	if changed {
		if err := m.SendControl(); err != nil {
			log.Printf("[mesh] send TakControl: %v", err)
		}
	}
}

func (m *Mesh) write(b []byte) error {
	sum := sha256.Sum256(b)
	m.mu.Lock()
	m.sent[sum] = time.Now()
	m.mu.Unlock()
	_, err := m.send.Write(b)
	return err
}

// isEcho reports whether b is a datagram we sent recently.
func (m *Mesh) isEcho(b []byte, now time.Time) bool {
	sum := sha256.Sum256(b)
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, t := range m.sent {
		if now.Sub(t) > echoWindow {
			delete(m.sent, k)
		}
	}
	_, ok := m.sent[sum]
	return ok
}

// SendControl broadcasts our TakControl. Rule 1 says it goes out as the
// lowest protocol version above 0 even when the mesh has fallen back to XML.
func (m *Mesh) SendControl() error {
	if m.cfg.ForceXML {
		return nil
	}
	return m.write(EncodeMesh(&TakMessage{Control: &TakControl{
		MinProtoVersion: 1, MaxProtoVersion: ourMaxVersion, ContactUID: m.cfg.ContactUID,
	}}))
}

// Send broadcasts an event in the negotiated protocol version.
func (m *Mesh) Send(ev *Event) error {
	if m.Version() >= 1 {
		return m.write(EncodeMesh(&TakMessage{Event: ev}))
	}
	return m.write(ev.MarshalCotXML(true))
}

// Run reads the multicast group until the socket closes, calling fn for every
// event heard. It also sends TakControl on schedule.
func (m *Mesh) Run(fn func(Received), onErr func(error)) {
	go func() {
		t := time.NewTicker(controlInterval)
		defer t.Stop()
		if err := m.SendControl(); err != nil {
			onErr(fmt.Errorf("send TakControl: %w", err))
		}
		for range t.C {
			m.recheck(time.Now())
			if err := m.SendControl(); err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				onErr(fmt.Errorf("send TakControl: %w", err))
			}
		}
	}()

	buf := make([]byte, 65536)
	for {
		n, from, err := m.recv.ReadFromUDP(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			onErr(fmt.Errorf("mesh read: %w", err))
			continue
		}
		now := time.Now()
		data := buf[:n]
		if m.isEcho(data, now) {
			continue
		}
		ev, ctl, version, err := Decode(data)
		if err != nil {
			onErr(fmt.Errorf("mesh decode from %s: %w", from, err))
			continue
		}
		m.observe(ev, ctl, version, now)
		m.recheck(now)
		if ev != nil {
			fn(Received{Event: ev, Version: version, Delivery: DeliveryMesh, From: from.String()})
		}
	}
}

// ServeDirect accepts directed TCP messages (connect, send one message, close),
// the way TAK clients reach a contact's "host:port:tcp" endpoint.
func ServeDirect(ln net.Listener, observe func(*Event, *TakControl, uint64), fn func(Received), onErr func(error)) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			onErr(fmt.Errorf("direct accept: %w", err))
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			_ = c.SetReadDeadline(time.Now().Add(directTimeout))
			data, err := io.ReadAll(io.LimitReader(c, maxDirectBytes+1))
			if err != nil {
				onErr(fmt.Errorf("direct read from %s: %w", c.RemoteAddr(), err))
				return
			}
			if len(data) > maxDirectBytes {
				onErr(fmt.Errorf("direct message from %s over %d bytes", c.RemoteAddr(), maxDirectBytes))
				return
			}
			ev, ctl, version, err := Decode(data)
			if err != nil {
				onErr(fmt.Errorf("direct decode from %s: %w", c.RemoteAddr(), err))
				return
			}
			if observe != nil {
				observe(ev, ctl, version)
			}
			if ev != nil {
				fn(Received{Event: ev, Version: version, Delivery: DeliveryDirect, From: c.RemoteAddr().String()})
			}
		}(conn)
	}
}

// ServeUDP reads unicast UDP CoT, such as an ATAK "output" pointed at the gateway.
func ServeUDP(conn *net.UDPConn, fn func(Received), onErr func(error)) {
	buf := make([]byte, 65536)
	for {
		n, from, err := conn.ReadFromUDP(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			onErr(fmt.Errorf("udp read: %w", err))
			continue
		}
		ev, _, version, err := Decode(buf[:n])
		if err != nil {
			onErr(fmt.Errorf("udp decode from %s: %w", from, err))
			continue
		}
		if ev != nil {
			fn(Received{Event: ev, Version: version, Delivery: DeliveryUDP, From: from.String()})
		}
	}
}

// presence builds the gateway's own SA, so TAK clients see the gateway on
// the map and can message it directly when it has an endpoint.
func presence(uid, callsign, cotType, endpoint string, lat, lon float64, now time.Time) *Event {
	ms := uint64(now.UnixMilli())
	return &Event{
		Type: cotType, UID: uid, How: "h-g-i-g-o",
		SendTime: ms, StartTime: ms, StaleTime: ms + uint64((3 * controlInterval).Milliseconds()),
		Lat: lat, Lon: lon, HAE: 9999999, CE: 9999999, LE: 9999999,
		Detail: Detail{
			Contact: &Contact{Callsign: callsign, Endpoint: endpoint},
			Takv:    &Takv{Device: "OneMind gateway", Platform: "omos-tak-mesh", OS: "linux", Version: buildVersion},
		},
	}
}

// dedupID identifies one CoT report so two gateways hearing the same SA
// produce one JetStream message (Nats-Msg-Id, 2m duplicate window).
func dedupID(ev *Event) string {
	sum := sha256.Sum256([]byte(ev.UID + "|" + ev.Type + "|" + strconv.FormatUint(ev.SendTime, 10) + "|" +
		strconv.FormatFloat(ev.Lat, 'f', -1, 64) + "|" + strconv.FormatFloat(ev.Lon, 'f', -1, 64)))
	return fmt.Sprintf("cot-%x", sum[:16])
}
