// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// ClientConfig configures the TAK Server streaming connection.
type ClientConfig struct {
	Addr string      // host:port of the streaming input (TAK Server 8089, TLS)
	TLS  *tls.Config // client certificate and server verification
	UID  string      // this client's CoT uid (pings and negotiation)

	// Protocol is "auto" (accept a v1 offer from the server) or "xml" (stay on v0).
	Protocol string

	// Keep-alive, after Common Commo (streamingsocketmanagement.cpp): ping after
	// IdleBeforePing without receiving anything, repeat every PingEvery, and drop
	// the connection after ResetAfter of silence.
	IdleBeforePing, PingEvery, ResetAfter time.Duration
	// NegotiateTimeout is how long to wait for t-x-takp-r (spec: at least 1 minute).
	NegotiateTimeout time.Duration
	DialTimeout      time.Duration
	BackoffMin       time.Duration
	BackoffMax       time.Duration

	InjectQueue  int           // bounded inject queue length
	InjectMaxAge time.Duration // queued injects older than this are dropped, not sent late
}

func (c *ClientConfig) defaults() {
	if c.Protocol == "" {
		c.Protocol = "auto"
	}
	set := func(d *time.Duration, v time.Duration) {
		if *d <= 0 {
			*d = v
		}
	}
	set(&c.IdleBeforePing, 15*time.Second)
	set(&c.PingEvery, 4500*time.Millisecond)
	set(&c.ResetAfter, 25*time.Second)
	set(&c.NegotiateTimeout, 60*time.Second)
	set(&c.DialTimeout, 10*time.Second)
	set(&c.BackoffMin, time.Second)
	set(&c.BackoffMax, 30*time.Second)
	set(&c.InjectMaxAge, 30*time.Second)
	if c.InjectQueue <= 0 {
		c.InjectQueue = 1024
	}
}

// Stats are cumulative counters, safe to read at any time.
type Stats struct {
	Connects, Disconnects   atomic.Int64
	Received, ParseErrors   atomic.Int64
	Injected, InjectDropped atomic.Int64
	Pings, Negotiations     atomic.Int64
	Connected               atomic.Bool
	Proto                   atomic.Int64 // 0 = XML, 1 = TAK Protocol v1
	LastError               atomic.Value // string
	ServerVersion           atomic.Value // string, from t-x-takp-v
}

type injectItem struct {
	ev *Event
	at time.Time
}

// Client is a TAK Server streaming client: it receives the CoT its certificate's
// groups allow and writes injected CoT, in XML or negotiated TAK Protocol v1.
type Client struct {
	cfg ClientConfig

	// OnEvent gets every received CoT event that is not stream plumbing, with
	// the XML to publish (the original bytes on a v0 stream). Called from the
	// reader goroutine; it must not block for long.
	OnEvent func(ev *Event, xmlBytes []byte, proto int)
	// OnState is called when the connection goes up (true) or down (false).
	OnState func(connected bool)

	inject   chan injectItem
	forceXML atomic.Bool // set after a negotiation timeout: later sessions stay on XML
	Stats    Stats
}

func NewClient(cfg ClientConfig) *Client {
	cfg.defaults()
	c := &Client{cfg: cfg, inject: make(chan injectItem, cfg.InjectQueue)}
	c.Stats.LastError.Store("")
	c.Stats.ServerVersion.Store("")
	return c
}

// Inject queues one event for TAK Server. When the queue is full the oldest
// item is dropped (and counted), so a stalled server never blocks NATS.
func (c *Client) Inject(ev *Event) {
	it := injectItem{ev: ev, at: time.Now()}
	for {
		select {
		case c.inject <- it:
			return
		default:
		}
		select {
		case <-c.inject:
			c.Stats.InjectDropped.Add(1)
		default:
		}
	}
}

func (c *Client) setErr(err error) {
	c.Stats.LastError.Store(err.Error())
}

// Run connects and reconnects until ctx is done.
func (c *Client) Run(ctx context.Context) {
	backoff := c.cfg.BackoffMin
	for ctx.Err() == nil {
		start := time.Now()
		err := c.session(ctx)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			c.setErr(err)
			log.Printf("tak-client: session to %s ended: %v", c.cfg.Addr, err)
		}
		if time.Since(start) > time.Minute {
			backoff = c.cfg.BackoffMin // the last session was healthy
		}
		wait := backoff/2 + rand.N(backoff/2+1) // jitter
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
		if backoff *= 2; backoff > c.cfg.BackoffMax {
			backoff = c.cfg.BackoffMax
		}
	}
}

// writer state shared between the reader and writer goroutines of one session.
type session struct {
	c        *Client
	conn     net.Conn
	lastRx   atomic.Int64 // unix nanos of the last received event
	offer    chan struct{}
	result   chan bool
	mu       sync.Mutex // serialises writes
	proto    int        // owned by the writer goroutine
	shutdown chan struct{}
}

func (c *Client) session(ctx context.Context) error {
	d := &net.Dialer{Timeout: c.cfg.DialTimeout}
	raw, err := d.DialContext(ctx, "tcp", c.cfg.Addr)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	conn := tls.Client(raw, c.cfg.TLS)
	hctx, cancel := context.WithTimeout(ctx, c.cfg.DialTimeout)
	err = conn.HandshakeContext(hctx)
	cancel()
	if err != nil {
		raw.Close()
		return fmt.Errorf("TLS handshake: %w", err)
	}
	s := &session{c: c, conn: conn, offer: make(chan struct{}, 1), result: make(chan bool, 1), shutdown: make(chan struct{})}
	s.lastRx.Store(time.Now().UnixNano())

	c.Stats.Connects.Add(1)
	c.Stats.Connected.Store(true)
	c.Stats.Proto.Store(0)
	log.Printf("tak-client: connected to %s (TLS %s)", c.cfg.Addr, tls.VersionName(conn.ConnectionState().Version))
	if c.OnState != nil {
		c.OnState(true)
	}
	defer func() {
		c.Stats.Connected.Store(false)
		c.Stats.Disconnects.Add(1)
		if c.OnState != nil {
			c.OnState(false)
		}
	}()

	errc := make(chan error, 2)
	go func() { errc <- s.readLoop() }()
	go func() { errc <- s.writeLoop(ctx) }()
	select {
	case err = <-errc:
	case <-ctx.Done():
		err = ctx.Err()
	}
	close(s.shutdown)
	conn.Close()
	<-errc // the other goroutine exits once the connection is closed
	return err
}

func (s *session) write(b []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := s.conn.Write(b)
	return err
}

func (s *session) encode(ev *Event, proto int) []byte {
	if proto == 1 {
		return EncodeStream(&TakMessage{Event: ev})
	}
	return streamXML(ev)
}

// writeLoop owns the protocol version for writes, pings, negotiation and injects.
func (s *session) writeLoop(ctx context.Context) error {
	c := s.c
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	var lastPing time.Time
	var negotiating bool
	var negDeadline time.Time
	for {
		// While negotiating, the spec forbids sending XML: injects wait in the queue.
		var injectCh chan injectItem
		if !negotiating {
			injectCh = c.inject
		}
		select {
		case <-s.shutdown:
			return nil
		case <-ctx.Done():
			return nil
		case <-s.offer:
			if negotiating || s.proto != 0 {
				continue
			}
			if err := s.write(protoRequest(c.cfg.UID, time.Now())); err != nil {
				return fmt.Errorf("write t-x-takp-q: %w", err)
			}
			c.Stats.Negotiations.Add(1)
			negotiating, negDeadline = true, time.Now().Add(c.cfg.NegotiateTimeout)
			log.Printf("tak-client: server offered TAK Protocol v1; requested it")
		case ok := <-s.result:
			if !negotiating {
				continue
			}
			negotiating = false
			if ok {
				s.proto = 1
				c.Stats.Proto.Store(1)
				log.Printf("tak-client: switched to TAK Protocol v1 streaming")
			} else {
				log.Printf("tak-client: server declined TAK Protocol v1; staying on XML")
			}
		case it := <-injectCh:
			if time.Since(it.at) > c.cfg.InjectMaxAge {
				c.Stats.InjectDropped.Add(1)
				continue
			}
			if err := s.write(s.encode(it.ev, s.proto)); err != nil {
				// Put it back so the next session sends it (if still fresh).
				c.Inject(it.ev)
				return fmt.Errorf("write inject: %w", err)
			}
			c.Stats.Injected.Add(1)
		case now := <-tick.C:
			if negotiating && now.After(negDeadline) {
				c.forceXML.Store(true)
				return errors.New("no t-x-takp-r within the negotiation timeout; reconnecting on XML")
			}
			idle := now.Sub(time.Unix(0, s.lastRx.Load()))
			if idle > c.cfg.ResetAfter {
				return fmt.Errorf("nothing received for %s; resetting connection", idle.Round(time.Second))
			}
			if !negotiating && idle > c.cfg.IdleBeforePing && now.Sub(lastPing) >= c.cfg.PingEvery {
				ping := controlEvent(c.cfg.UID+"-ping", typePing, "", now, 20*time.Second)
				ping.How = "h-g-i-g-o"
				if err := s.write(s.encode(ping, s.proto)); err != nil {
					return fmt.Errorf("write ping: %w", err)
				}
				c.Stats.Pings.Add(1)
				lastPing = now
			}
		}
	}
}

// readLoop parses the inbound stream, switching framing after a v1 agreement.
func (s *session) readLoop() error {
	c := s.c
	br := bufio.NewReaderSize(s.conn, 64*1024)
	xr := newXMLReader(br)
	for {
		raw, err := xr.Next()
		if err != nil {
			return readErr(err)
		}
		s.lastRx.Store(time.Now().UnixNano())
		ev, err := ParseCotXML(raw)
		if err != nil {
			c.Stats.ParseErrors.Add(1)
			c.setErr(fmt.Errorf("parse CoT from TAK: %w", err))
			continue
		}
		switch ev.Type {
		case typeProtoOffer:
			ok, sv := offersV1(ev)
			if sv != "" {
				c.Stats.ServerVersion.Store(sv)
				log.Printf("tak-client: TAK Server version %s", sv)
			}
			if ok && c.cfg.Protocol == "auto" && !c.forceXML.Load() {
				select {
				case s.offer <- struct{}{}:
				default:
				}
			}
			continue
		case typeProtoResponse:
			accepted := responseAccepted(ev)
			select {
			case s.result <- accepted:
			default:
			}
			if accepted {
				// Every byte after this event is v1 framing, including any
				// already buffered past </event>.
				v1 := bufio.NewReaderSize(io.MultiReader(bytes.NewReader(xr.Leftover()), br), 64*1024)
				return s.readV1(v1)
			}
			continue
		}
		if isControlType(ev.Type) {
			continue
		}
		c.Stats.Received.Add(1)
		if c.OnEvent != nil {
			c.OnEvent(ev, raw, 0)
		}
	}
}

func (s *session) readV1(br *bufio.Reader) error {
	c := s.c
	for {
		m, err := readStreamMessage(br)
		if err != nil {
			if errors.Is(err, errEventTooLarge) || errors.Is(err, errBadFrame) {
				return err
			}
			if isIOErr(err) {
				return readErr(err)
			}
			// A payload that does not decode: count it and keep the stream.
			c.Stats.ParseErrors.Add(1)
			c.setErr(fmt.Errorf("decode v1 from TAK: %w", err))
			continue
		}
		s.lastRx.Store(time.Now().UnixNano())
		if m.Event == nil || isControlType(m.Event.Type) {
			continue
		}
		c.Stats.Received.Add(1)
		if c.OnEvent != nil {
			c.OnEvent(m.Event, m.Event.MarshalCotXML(false), 1)
		}
	}
}

func isIOErr(err error) bool {
	var ne net.Error
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) || errors.As(err, &ne)
}

func readErr(err error) error {
	if errors.Is(err, io.EOF) {
		return errors.New("TAK Server closed the stream")
	}
	return fmt.Errorf("read: %w", err)
}
