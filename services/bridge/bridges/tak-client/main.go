// SPDX-License-Identifier: Apache-2.0

// Package main is the OneMind TAK Server client bridge (tak-client). It joins
// TAK Server's streaming input (8089, mutual TLS) the way an ATAK device does,
// and connects that stream to the authority bus:
//
//	TAK Server ──► tak.cot.out      (CoT XML, acked JetStream publish to stream COT, Nats-Msg-Id dedup)
//	TAK Server ◄── tak.cot.inject   (queue group shared with NatsBridgePlugin, so each message is injected once)
//
// It negotiates TAK Protocol v1 when the server offers it (XML otherwise),
// keeps the connection alive with Commo-style pings, reconnects with backoff,
// and drops CoT that TAK echoes back after it was published to tak.cot.inject.
// It replaces services/takserver/cot-tap (see README.md, "Retiring cot-tap").
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

const buildVersion = "0.1.0"

const (
	subjectOut    = "tak.cot.out"
	subjectInject = "tak.cot.inject"
)

type config struct {
	natsURL, natsUser, natsPassword string
	clientID                        string
	takAddr                         string
	certFile, keyFile, caFile       string
	serverName                      string
	insecureSkipVerify              bool
	protocol                        string
	injectEnabled                   bool
	injectQueue                     string
	echoSuppress                    bool
	echoTTL                         time.Duration
	outJetStream                    bool
	httpAddr                        string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key, def string) (bool, error) {
	v, err := strconv.ParseBool(getenv(key, def))
	if err != nil {
		return false, fmt.Errorf("%s: %w", key, err)
	}
	return v, nil
}

func loadConfig() (config, error) {
	c := config{
		natsURL:      getenv("NATS_URL", "nats://127.0.0.1:4222"),
		natsUser:     os.Getenv("NATS_USER"),
		natsPassword: os.Getenv("NATS_PASSWORD"),
		clientID:     getenv("TAK_CLIENT_ID", "omos-tak-client"),
		takAddr:      getenv("TAK_ADDR", "takserver-messaging-service.takserver.svc.cluster.local:8089"),
		certFile:     os.Getenv("TAK_CERT_FILE"),
		keyFile:      os.Getenv("TAK_KEY_FILE"),
		caFile:       os.Getenv("TAK_CA_FILE"),
		serverName:   os.Getenv("TAK_SERVER_NAME"),
		protocol:     getenv("TAK_PROTOCOL", "auto"),
		injectQueue:  getenv("INJECT_QUEUE", "tak-cot-inject"),
		httpAddr:     getenv("HTTP_ADDR", ":8080"),
	}
	var err error
	if c.insecureSkipVerify, err = envBool("TAK_INSECURE_SKIP_VERIFY", "false"); err != nil {
		return c, err
	}
	if c.injectEnabled, err = envBool("INJECT_ENABLED", "true"); err != nil {
		return c, err
	}
	if c.echoSuppress, err = envBool("ECHO_SUPPRESS", "true"); err != nil {
		return c, err
	}
	if c.outJetStream, err = envBool("OUT_JETSTREAM", "true"); err != nil {
		return c, err
	}
	if c.echoTTL, err = time.ParseDuration(getenv("ECHO_TTL", "2m")); err != nil {
		return c, fmt.Errorf("ECHO_TTL: %w", err)
	}
	if c.protocol != "auto" && c.protocol != "xml" {
		return c, fmt.Errorf("TAK_PROTOCOL must be auto or xml, got %q", c.protocol)
	}
	if c.certFile == "" || c.keyFile == "" {
		return c, errors.New("TAK_CERT_FILE and TAK_KEY_FILE are required (TAK Server streaming needs a client certificate)")
	}
	if c.caFile == "" && !c.insecureSkipVerify {
		return c, errors.New("TAK_CA_FILE is required unless TAK_INSECURE_SKIP_VERIFY=true")
	}
	return c, nil
}

// tlsConfig builds the client TLS config. With TAK_SERVER_NAME set, the server
// certificate must match that name. Without it, the chain is verified against
// TAK_CA_FILE but the hostname is not, as Common Commo does for TAK Server
// (TAK Server certificates rarely name the in-cluster Service).
func tlsConfig(c config) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(c.certFile, c.keyFile)
	if err != nil {
		return nil, fmt.Errorf("client certificate: %w", err)
	}
	tc := &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12}
	if c.insecureSkipVerify {
		log.Printf("tak-client: WARNING TAK_INSECURE_SKIP_VERIFY=true, the TAK Server certificate is not checked")
		tc.InsecureSkipVerify = true
		return tc, nil
	}
	pem, err := os.ReadFile(c.caFile)
	if err != nil {
		return nil, fmt.Errorf("TAK_CA_FILE: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pem) {
		return nil, errors.New("TAK_CA_FILE: no PEM certificates found")
	}
	if c.serverName != "" {
		tc.RootCAs, tc.ServerName = pool, c.serverName
		return tc, nil
	}
	tc.InsecureSkipVerify = true // replaced by the chain-only check below
	tc.VerifyConnection = func(cs tls.ConnectionState) error {
		if len(cs.PeerCertificates) == 0 {
			return errors.New("TAK Server sent no certificate")
		}
		inter := x509.NewCertPool()
		for _, ic := range cs.PeerCertificates[1:] {
			inter.AddCert(ic)
		}
		_, err := cs.PeerCertificates[0].Verify(x509.VerifyOptions{
			Roots: pool, Intermediates: inter, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		})
		return err
	}
	return tc, nil
}

// bridge connects one Client to NATS.
type bridge struct {
	cfg    config
	nc     *nats.Conn
	js     nats.JetStreamContext
	client *Client
	echo   *echoFilter

	mu       sync.Mutex
	queueSub *nats.Subscription

	published, publishErrors, echoDropped, injectRejected atomic.Int64
}

// onState joins the inject queue group only while TAK is connected, so
// NatsBridgePlugin (same queue group) takes every message while we are down.
func (b *bridge) onState(up bool) {
	if !b.cfg.injectEnabled {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if up && b.queueSub == nil {
		sub, err := b.nc.QueueSubscribe(subjectInject, b.cfg.injectQueue, b.onInject)
		if err != nil {
			log.Printf("tak-client: subscribe %s: %v", subjectInject, err)
			return
		}
		b.queueSub = sub
		log.Printf("tak-client: injecting from %s (queue %s)", subjectInject, b.cfg.injectQueue)
	}
	if !up && b.queueSub != nil {
		_ = b.queueSub.Unsubscribe()
		b.queueSub = nil
	}
}

// onInject validates one inject message and queues it for TAK Server. The
// event is re-rendered from its parsed form, so one NATS message can never
// carry a second, smuggled event.
func (b *bridge) onInject(m *nats.Msg) {
	ev, err := ParseCotXML(m.Data)
	if err != nil {
		b.injectRejected.Add(1)
		b.client.setErr(fmt.Errorf("%s: %w", subjectInject, err))
		return
	}
	if isControlType(ev.Type) {
		b.injectRejected.Add(1)
		return
	}
	b.client.Inject(ev)
}

// onSeenInject feeds the echo filter from every tak.cot.inject message.
func (b *bridge) onSeenInject(m *nats.Msg) {
	if ev, err := ParseCotXML(m.Data); err == nil {
		b.echo.Remember(ev, time.Now())
	}
}

func (b *bridge) onEvent(ev *Event, xmlBytes []byte, proto int) {
	if b.echo != nil && b.echo.IsEcho(ev, time.Now()) {
		b.echoDropped.Add(1)
		return
	}
	msg := nats.NewMsg(subjectOut)
	msg.Data = xmlBytes
	msg.Header.Set(nats.MsgIdHdr, dedupID(ev))
	msg.Header.Set("Tak-Client", b.cfg.clientID)
	msg.Header.Set("Tak-Proto", strconv.Itoa(proto))
	msg.Header.Set("Tak-Server", b.cfg.takAddr)
	if ev.Detail.Extensions > 0 {
		msg.Header.Set("Tak-Undecoded-Extensions", strconv.Itoa(ev.Detail.Extensions))
	}
	if b.js == nil {
		if err := b.nc.PublishMsg(msg); err != nil {
			b.publishFailed(ev, err)
			return
		}
		b.published.Add(1)
		return
	}
	b.publishAcked(msg, ev, 0)
}

// publishAcked publishes with a JetStream ack and retries a missing ack with
// backoff; the Nats-Msg-Id makes a retry of a stored message a no-op.
func (b *bridge) publishAcked(msg *nats.Msg, ev *Event, attempt int) {
	f, err := b.js.PublishMsgAsync(msg)
	if err != nil {
		b.retryOrFail(msg, ev, attempt, err)
		return
	}
	go func() {
		select {
		case <-f.Ok():
			b.published.Add(1)
		case err := <-f.Err():
			b.retryOrFail(msg, ev, attempt, err)
		}
	}()
}

func (b *bridge) retryOrFail(msg *nats.Msg, ev *Event, attempt int, err error) {
	if attempt >= 3 {
		b.publishFailed(ev, err)
		return
	}
	time.AfterFunc(time.Duration(250<<attempt)*time.Millisecond, func() {
		m := nats.NewMsg(msg.Subject)
		m.Data, m.Header = msg.Data, msg.Header
		b.publishAcked(m, ev, attempt+1)
	})
}

func (b *bridge) publishFailed(ev *Event, err error) {
	b.publishErrors.Add(1)
	b.client.setErr(fmt.Errorf("publish %s uid=%s: %w", subjectOut, ev.UID, err))
	log.Printf("tak-client: publish %s uid=%s failed: %v", subjectOut, ev.UID, err)
}

func (b *bridge) ready() bool {
	return b.nc.IsConnected() && b.client.Stats.Connected.Load()
}

func (b *bridge) metrics(w http.ResponseWriter, _ *http.Request) {
	s := &b.client.Stats
	bool01 := func(v bool) int {
		if v {
			return 1
		}
		return 0
	}
	var sb strings.Builder
	line := func(name, help, typ string, v any) {
		fmt.Fprintf(&sb, "# HELP %s %s\n# TYPE %s %s\n%s %v\n", name, help, name, typ, name, v)
	}
	line("tak_client_tak_connected", "1 while the TAK Server stream is up", "gauge", bool01(s.Connected.Load()))
	line("tak_client_nats_connected", "1 while NATS is connected", "gauge", bool01(b.nc.IsConnected()))
	line("tak_client_protocol_version", "0 = XML, 1 = TAK Protocol v1", "gauge", s.Proto.Load())
	line("tak_client_connects_total", "TAK Server connections made", "counter", s.Connects.Load())
	line("tak_client_received_total", "CoT events received from TAK Server", "counter", s.Received.Load())
	line("tak_client_parse_errors_total", "Inbound CoT that did not parse", "counter", s.ParseErrors.Load())
	line("tak_client_published_total", "CoT published to tak.cot.out (acked)", "counter", b.published.Load())
	line("tak_client_publish_errors_total", "tak.cot.out publishes that failed after retries", "counter", b.publishErrors.Load())
	line("tak_client_echo_dropped_total", "Inbound CoT dropped as an echo of tak.cot.inject", "counter", b.echoDropped.Load())
	line("tak_client_injected_total", "CoT written to TAK Server from tak.cot.inject", "counter", s.Injected.Load())
	line("tak_client_inject_dropped_total", "Injects dropped (queue full or too old)", "counter", s.InjectDropped.Load())
	line("tak_client_inject_rejected_total", "tak.cot.inject messages that were not valid CoT", "counter", b.injectRejected.Load())
	line("tak_client_pings_total", "Keep-alive pings sent", "counter", s.Pings.Load())
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	_, _ = w.Write([]byte(sb.String()))
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("tak-client: %v", err)
	}
	tc, err := tlsConfig(cfg)
	if err != nil {
		log.Fatalf("tak-client: %v", err)
	}
	log.Printf("tak-client v%s: TAK %s (protocol %s), NATS %s", buildVersion, cfg.takAddr, cfg.protocol, cfg.natsURL)

	opts := []nats.Option{
		nats.Name("tak-client-" + cfg.clientID),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
		nats.RetryOnFailedConnect(true),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) { log.Printf("tak-client: NATS disconnected: %v", err) }),
		nats.ReconnectHandler(func(nc *nats.Conn) { log.Printf("tak-client: NATS reconnected to %s", nc.ConnectedUrl()) }),
	}
	if cfg.natsUser != "" {
		opts = append(opts, nats.UserInfo(cfg.natsUser, cfg.natsPassword))
	}
	nc, err := nats.Connect(cfg.natsURL, opts...)
	if err != nil {
		log.Fatalf("tak-client: NATS: %v", err)
	}
	defer nc.Drain()

	b := newBridge(cfg, nc, tc)
	if err := b.start(); err != nil {
		log.Fatalf("tak-client: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok\n")) })
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if !b.ready() {
			http.Error(w, "TAK Server or NATS not connected", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/metrics", b.metrics)
	srv := &http.Server{Addr: cfg.httpAddr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("tak-client: HTTP %s: %v", cfg.httpAddr, err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	b.client.Run(ctx)
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdown)
	log.Printf("tak-client: stopped")
}

func newBridge(cfg config, nc *nats.Conn, tc *tls.Config) *bridge {
	b := &bridge{cfg: cfg, nc: nc}
	b.client = NewClient(ClientConfig{Addr: cfg.takAddr, TLS: tc, UID: cfg.clientID, Protocol: cfg.protocol})
	b.client.OnEvent = b.onEvent
	b.client.OnState = b.onState
	if cfg.echoSuppress {
		b.echo = newEchoFilter(cfg.echoTTL, 200000)
	}
	return b
}

// start sets up JetStream and the echo subscription; the inject queue
// subscription follows the TAK connection (onState).
func (b *bridge) start() error {
	if b.cfg.outJetStream {
		js, err := b.nc.JetStream(nats.PublishAsyncMaxPending(256), nats.PublishAsyncTimeout(5*time.Second))
		if err != nil {
			return fmt.Errorf("JetStream: %w", err)
		}
		b.js = js
	}
	if b.echo != nil {
		if _, err := b.nc.Subscribe(subjectInject, b.onSeenInject); err != nil {
			return fmt.Errorf("subscribe %s (echo filter): %w", subjectInject, err)
		}
	}
	return nil
}
