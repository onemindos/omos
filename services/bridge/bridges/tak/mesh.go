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
	"strings"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
)

const (
	ourMaxVersion   = 1
	controlInterval = 60 * time.Second // rule 1: TakControl at least every 60s
	controlTimeout  = 2 * time.Minute  // rule 3c
	contactTimeout  = 5 * time.Minute  // forget contacts we no longer hear
	echoWindow      = 10 * time.Second // drop our own datagrams looped back
	maxDirectBytes  = 256 << 10        // largest directed TCP message accepted
	maxDirectConns  = 16               // directed TCP connections handled at once
	directTimeout   = 10 * time.Second
)

// Delivery says how a message reached the gateway.
type Delivery string

const (
	DeliveryMesh   Delivery = "mesh"   // multicast group
	DeliveryDirect Delivery = "direct" // TCP send-once to the gateway
	DeliveryUDP    Delivery = "udp"    // unicast UDP input
	DeliveryChat   Delivery = "chat"   // GeoChat multicast group
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
	endpoint    string // "host:port:proto" from the contact's SA, if any
}

// MeshConfig configures a Mesh.
type MeshConfig struct {
	Group      *net.UDPAddr   // SA multicast group, e.g. 239.2.3.1:6969
	ChatGroup  *net.UDPAddr   // GeoChat "All Chat Rooms" group, e.g. 224.10.10.1:17012; nil disables
	Interface  *net.Interface // join and transmit on this interface; nil = system default
	TTL        int            // multicast TTL; 0 means 1, as Commo uses
	ContactUID string         // our identity in TakControl
	ForceXML   bool           // never send protobuf (legacy-only meshes)
	Crypto     *MeshCrypto    // mesh encryption; nil = plaintext
}

// groupConn is one joined multicast group.
type groupConn struct {
	conn     *net.UDPConn
	group    *net.UDPAddr
	delivery Delivery
}

// Mesh is a TAK mesh participant.
type Mesh struct {
	cfg  MeshConfig
	recv []groupConn
	// send is an unconnected socket whose multicast interface is set with
	// IP_MULTICAST_IF, so MESH_INTERFACE applies to transmit as well as
	// receive. Without it, the kernel picks the default route's interface.
	send *net.UDPConn

	mu       sync.Mutex
	contacts map[string]*contact
	sent     map[[32]byte]time.Time
	version  uint32 // last computed broadcast version
}

// NewMesh joins the SA group (and the chat group, if set) and opens the
// transmit socket on the configured interface.
func NewMesh(cfg MeshConfig) (*Mesh, error) {
	m := &Mesh{cfg: cfg, contacts: map[string]*contact{}, sent: map[[32]byte]time.Time{}}
	groups := []groupConn{{group: cfg.Group, delivery: DeliveryMesh}}
	if cfg.ChatGroup != nil {
		groups = append(groups, groupConn{group: cfg.ChatGroup, delivery: DeliveryChat})
	}
	for _, g := range groups {
		conn, err := net.ListenMulticastUDP("udp4", cfg.Interface, g.group)
		if err != nil {
			m.Close()
			return nil, fmt.Errorf("join %s: %w", g.group, err)
		}
		_ = conn.SetReadBuffer(1 << 20)
		g.conn = conn
		m.recv = append(m.recv, g)
	}
	send, err := openTransmit(cfg.Interface, cfg.TTL)
	if err != nil {
		m.Close()
		return nil, err
	}
	m.send = send
	m.version = m.computeVersion(time.Now())
	return m, nil
}

// openTransmit opens the multicast transmit socket. The interface is applied
// with IP_MULTICAST_IF; loopback stays on so peers on the same host (another
// gateway, a test peer) hear us, and our own copy is dropped by isEcho.
func openTransmit(iface *net.Interface, ttl int) (*net.UDPConn, error) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero})
	if err != nil {
		return nil, fmt.Errorf("open transmit socket: %w", err)
	}
	pc := ipv4.NewPacketConn(conn)
	if iface != nil {
		if err := pc.SetMulticastInterface(iface); err != nil {
			conn.Close()
			return nil, fmt.Errorf("transmit on %s: %w", iface.Name, err)
		}
	}
	if ttl <= 0 {
		ttl = 1
	}
	if err := pc.SetMulticastTTL(ttl); err != nil {
		conn.Close()
		return nil, fmt.Errorf("multicast TTL %d: %w", ttl, err)
	}
	if err := pc.SetMulticastLoopback(true); err != nil {
		conn.Close()
		return nil, fmt.Errorf("multicast loopback: %w", err)
	}
	return conn, nil
}

// TransmitInterface is the interface the transmit socket sends multicast on,
// as the kernel reports it (nil when the system default is used).
func (m *Mesh) TransmitInterface() (*net.Interface, error) {
	return ipv4.NewPacketConn(m.send).MulticastInterface()
}

func (m *Mesh) Close() {
	for _, g := range m.recv {
		g.conn.Close()
	}
	if m.send != nil {
		m.send.Close()
	}
}

// open decrypts an inbound message when mesh encryption is on.
func (m *Mesh) open(b []byte) ([]byte, error) {
	if m.cfg.Crypto == nil {
		return b, nil
	}
	return m.cfg.Crypto.Decrypt(b)
}

// seal encrypts an outbound message when mesh encryption is on.
func (m *Mesh) seal(b []byte) ([]byte, error) {
	if m.cfg.Crypto == nil {
		return b, nil
	}
	return m.cfg.Crypto.Encrypt(b)
}

// DecodeInbound decrypts (if configured) and decodes one mesh, directed or
// unicast message. With encryption on, plaintext is rejected, as in Commo.
func (m *Mesh) DecodeInbound(b []byte) (*Event, *TakControl, uint64, error) {
	plain, err := m.open(b)
	if err != nil {
		return nil, nil, 0, err
	}
	return Decode(plain)
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
	// takcontrol.proto: contactUid may be omitted when the TakControl is paired
	// with the sender's own CotEvent.
	ctlUID := ""
	if ctl != nil {
		ctlUID = ctl.ContactUID
		if ctlUID == "" && ev != nil {
			ctlUID = ev.UID
		}
	}
	if ctlUID != "" && ctlUID != m.cfg.ContactUID {
		c := update(ctlUID)
		c.maxVersion = controlVersion(ctl.MaxProtoVersion) // 3b
		c.lastControl = now
	}
	// SA from a mesh participant carries a contact endpoint. Markers and other
	// relayed events do not identify a client, so they are not tracked.
	if ev != nil && ev.UID != m.cfg.ContactUID && ev.Detail.Contact != nil && ev.Detail.Contact.Endpoint != "" {
		var c *contact
		if ctlUID != ev.UID {
			c = update(ev.UID)
		} else {
			c = m.contacts[ev.UID]
		}
		if c != nil {
			c.endpoint = ev.Detail.Contact.Endpoint
		}
	}
}

// controlVersion reads a TakControl version field. takcontrol.proto: "If not
// filled in (reads as 0), version 1 is assumed"; README rule 3b: TakControl
// never carries version 0. Python takproto, for example, attaches a TakControl
// with only contactUid set to every GeoChat message.
func controlVersion(v uint32) uint32 {
	if v == 0 {
		return 1
	}
	return v
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

// Endpoint returns the contact endpoint ("host:port:proto") last advertised
// by uid, or "" if unknown.
func (m *Mesh) Endpoint(uid string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c := m.contacts[uid]; c != nil {
		return c.endpoint
	}
	return ""
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

func (m *Mesh) write(b []byte, to *net.UDPAddr) error {
	b, err := m.seal(b)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(b)
	m.mu.Lock()
	m.sent[sum] = time.Now()
	m.mu.Unlock()
	_, err = m.send.WriteToUDP(b, to)
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
	}}), m.cfg.Group)
}

// IsChat reports whether ev is a GeoChat message (b-t-f, b-t-f-d, ...).
func IsChat(ev *Event) bool { return ev.Type == "b-t-f" || strings.HasPrefix(ev.Type, "b-t-f-") }

// encode renders ev in the negotiated protocol version.
func (m *Mesh) encode(ev *Event) []byte {
	if m.Version() >= 1 {
		return EncodeMesh(&TakMessage{Event: ev})
	}
	return ev.MarshalCotXML(true)
}

// Send broadcasts an event in the negotiated protocol version. GeoChat goes
// to the chat group when one is configured, everything else to the SA group.
func (m *Mesh) Send(ev *Event) error {
	to := m.cfg.Group
	if m.cfg.ChatGroup != nil && IsChat(ev) {
		to = m.cfg.ChatGroup
	}
	return m.write(m.encode(ev), to)
}

// SendDirect delivers ev to one contact endpoint ("host:port" or
// "host:port:tcp") as a directed message: connect, send, close.
func (m *Mesh) SendDirect(ev *Event, endpoint string) error {
	addr, err := directAddr(endpoint)
	if err != nil {
		return err
	}
	b, err := m.seal(m.encode(ev))
	if err != nil {
		return err
	}
	c, err := net.DialTimeout("tcp4", addr, directTimeout)
	if err != nil {
		return err
	}
	defer c.Close()
	_ = c.SetWriteDeadline(time.Now().Add(directTimeout))
	_, err = c.Write(b)
	return err
}

// directAddr turns a TAK endpoint into a TCP dial address. Only TCP
// endpoints can take directed messages; "*:-1:stcp" means "via a server".
func directAddr(endpoint string) (string, error) {
	parts := strings.Split(endpoint, ":")
	switch {
	case len(parts) == 2:
	case len(parts) == 3 && parts[2] == "tcp":
	default:
		return "", fmt.Errorf("endpoint %q is not host:port:tcp", endpoint)
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil || port < 1 || port > 65535 || parts[0] == "" || parts[0] == "*" {
		return "", fmt.Errorf("endpoint %q is not host:port:tcp", endpoint)
	}
	return net.JoinHostPort(parts[0], parts[1]), nil
}

// Run reads the multicast groups until the sockets close, calling fn for every
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

	var wg sync.WaitGroup
	for _, g := range m.recv {
		wg.Add(1)
		go func(g groupConn) {
			defer wg.Done()
			m.readGroup(g, fn, onErr)
		}(g)
	}
	wg.Wait()
}

func (m *Mesh) readGroup(g groupConn, fn func(Received), onErr func(error)) {
	buf := make([]byte, 65536)
	for {
		n, from, err := g.conn.ReadFromUDP(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			onErr(fmt.Errorf("%s read: %w", g.delivery, err))
			continue
		}
		now := time.Now()
		data := buf[:n]
		if m.isEcho(data, now) {
			continue
		}
		ev, ctl, version, err := m.DecodeInbound(data)
		if err != nil {
			onErr(fmt.Errorf("%s decode from %s: %w", g.delivery, from, err))
			continue
		}
		m.observe(ev, ctl, version, now)
		m.recheck(now)
		if ev != nil {
			fn(Received{Event: ev, Version: version, Delivery: g.delivery, From: from.String()})
		}
	}
}

// decodeFunc decodes one inbound message (Decode, or Mesh.DecodeInbound when
// mesh encryption may be on).
type decodeFunc func([]byte) (*Event, *TakControl, uint64, error)

// ServeDirect accepts directed TCP messages (connect, send one message, close),
// the way TAK clients reach a contact's "host:port:tcp" endpoint. At most
// maxDirectConns are read at once; further connections are closed at once, so
// a LAN flood holds at most maxDirectConns*maxDirectBytes of buffers.
func ServeDirect(ln net.Listener, decode decodeFunc, observe func(*Event, *TakControl, uint64), fn func(Received), onErr func(error)) {
	sem := make(chan struct{}, maxDirectConns)
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			onErr(fmt.Errorf("direct accept: %w", err))
			continue
		}
		select {
		case sem <- struct{}{}:
		default:
			onErr(fmt.Errorf("direct: %d connections busy, closed %s", maxDirectConns, conn.RemoteAddr()))
			conn.Close()
			continue
		}
		go func(c net.Conn) {
			defer func() { <-sem }()
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
			ev, ctl, version, err := decode(data)
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
func ServeUDP(conn *net.UDPConn, decode decodeFunc, fn func(Received), onErr func(error)) {
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
		ev, _, version, err := decode(buf[:n])
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
