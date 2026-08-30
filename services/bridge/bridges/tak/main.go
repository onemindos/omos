// SPDX-License-Identifier: Apache-2.0

// Package main implements the TAK/CoT bridge for the One Mind OS (omos) platform.
// It listens for Cursor-on-Target (CoT) XML messages from ATAK/WinTAK/OTS over UDP,
// translates them into the omos One Mind envelope, and publishes to NATS.
// It also subscribes to command subjects and forwards commands back to a TAK server.
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

// ── CoT XML structures ────────────────────────────────────────────────────────

// CotEvent is the top-level CoT XML element.
type CotEvent struct {
	XMLName xml.Name `xml:"event"`
	Version string   `xml:"version,attr"`
	UID     string   `xml:"uid,attr"`
	Type    string   `xml:"type,attr"`
	How     string   `xml:"how,attr"`
	Time    string   `xml:"time,attr"`
	Start   string   `xml:"start,attr"`
	Stale   string   `xml:"stale,attr"`
	Point   CotPoint `xml:"point"`
	Detail  CotDetail
}

// CotPoint holds geographic position.
type CotPoint struct {
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
	HAE string `xml:"hae,attr"`
	CE  string `xml:"ce,attr"`
	LE  string `xml:"le,attr"`
}

// CotDetail wraps the inner detail element. We parse callsign from it
// and also capture the raw XML bytes.
type CotDetail struct {
	Contact *CotContact `xml:"contact"`
	InnerXML []byte
}

// UnmarshalXML for CotDetail captures both the raw bytes and the contact child.
func (d *CotDetail) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	// Read the inner raw XML by buffering a token stream.
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)

	depth := 0
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			// Check for contact element.
			if t.Name.Local == "contact" {
				contact := &CotContact{}
				for _, attr := range t.Attr {
					if attr.Name.Local == "callsign" {
						contact.Callsign = attr.Value
					}
				}
				d.Contact = contact
			}
			if err2 := enc.EncodeToken(t); err2 != nil {
				return err2
			}
		case xml.EndElement:
			if depth == 0 {
				// This is our closing </detail> — stop.
				goto done
			}
			depth--
			if err2 := enc.EncodeToken(t); err2 != nil {
				return err2
			}
		default:
			if err2 := enc.EncodeToken(t); err2 != nil {
				return err2
			}
		}
	}
done:
	enc.Flush()
	d.InnerXML = buf.Bytes()
	return nil
}

// CotContact holds contact info inside detail.
type CotContact struct {
	Callsign string `xml:"callsign,attr"`
}

// ── One Mind payload structures ───────────────────────────────────────────────

// CotPayload is what gets embedded in the omos envelope.
type CotPayload struct {
	UID       string  `json:"uid"`
	CotType   string  `json:"cot_type"`
	How       string  `json:"how"`
	Callsign  string  `json:"callsign"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	HAE       float64 `json:"hae"`
	CE        float64 `json:"ce"`
	LE        float64 `json:"le"`
	StaleAt   int64   `json:"stale_at"`  // unix ms
	DetailRaw string  `json:"detail_raw"` // raw XML of <detail>…</detail>
}

// EnvelopePayload wraps the domain-specific message type and data.
type EnvelopePayload struct {
	MsgType string      `json:"msg_type"`
	Data    interface{} `json:"data"`
}

// Envelope is the standard One Mind message envelope.
type Envelope struct {
	ID      string          `json:"id"`       // timestamp nanos as string
	TS      int64           `json:"ts"`       // unix ms
	Source  string          `json:"source"`
	Entity  string          `json:"entity"`
	Domain  string          `json:"domain"`
	Payload EnvelopePayload `json:"payload"`
}

// ── Config ────────────────────────────────────────────────────────────────────

type config struct {
	natsURL       string
	bridgeName    string
	listenPort    int
	serverHost    string
	serverPort    int
}

func loadConfig() config {
	c := config{
		natsURL:    getenv("NATS_URL", "nats://localhost:4222"),
		bridgeName: getenv("BRIDGE_NAME", "tak-bridge"),
		serverHost: getenv("TAK_SERVER_HOST", "localhost"),
	}
	var err error
	if c.listenPort, err = strconv.Atoi(getenv("TAK_LISTEN_PORT", "8087")); err != nil {
		log.Fatalf("invalid TAK_LISTEN_PORT: %v", err)
	}
	if c.serverPort, err = strconv.Atoi(getenv("TAK_SERVER_PORT", "8087")); err != nil {
		log.Fatalf("invalid TAK_SERVER_PORT: %v", err)
	}
	return c
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ── Helpers ───────────────────────────────────────────────────────────────────

var unsafeChars = regexp.MustCompile(`[^a-zA-Z0-9\-_]`)

// sanitizeUID replaces any characters that are not alphanumeric, dash, or
// underscore with an underscore so the UID is safe to embed in a NATS subject.
func sanitizeUID(uid string) string {
	return unsafeChars.ReplaceAllString(uid, "_")
}

// parseFloat parses a string to float64; returns 0.0 on error.
func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return v
}

// parseCoT decodes raw CoT XML bytes into a CotPayload.
func parseCoT(raw []byte) (*CotPayload, error) {
	var evt CotEvent
	if err := xml.Unmarshal(raw, &evt); err != nil {
		return nil, fmt.Errorf("xml decode: %w", err)
	}

	// Parse stale time → unix ms.
	var staleAt int64
	if t, err := time.Parse(time.RFC3339, evt.Stale); err == nil {
		staleAt = t.UnixMilli()
	}

	// Extract callsign.
	callsign := ""
	if evt.Detail.Contact != nil {
		callsign = evt.Detail.Contact.Callsign
	}

	// Reconstruct <detail>…</detail> raw XML.
	detailRaw := "<detail>" + string(evt.Detail.InnerXML) + "</detail>"

	p := &CotPayload{
		UID:       evt.UID,
		CotType:   evt.Type,
		How:       evt.How,
		Callsign:  callsign,
		Lat:       parseFloat(evt.Point.Lat),
		Lon:       parseFloat(evt.Point.Lon),
		HAE:       parseFloat(evt.Point.HAE),
		CE:        parseFloat(evt.Point.CE),
		LE:        parseFloat(evt.Point.LE),
		StaleAt:   staleAt,
		DetailRaw: detailRaw,
	}
	return p, nil
}

// buildEnvelope wraps a CotPayload into the standard One Mind envelope.
func buildEnvelope(p *CotPayload) ([]byte, string, error) {
	now := time.Now()
	sanUID := sanitizeUID(p.UID)
	subject := fmt.Sprintf("omos.field.%s.telemetry", sanUID)

	env := Envelope{
		ID:     strconv.FormatInt(now.UnixNano(), 10),
		TS:     now.UnixMilli(),
		Source: "tak-bridge",
		Entity: fmt.Sprintf("field.%s", sanUID),
		Domain: "field",
		Payload: EnvelopePayload{
			MsgType: "cot",
			Data:    p,
		},
	}
	b, err := json.Marshal(env)
	if err != nil {
		return nil, "", fmt.Errorf("json marshal: %w", err)
	}
	return b, subject, nil
}

// ── UDP helpers ───────────────────────────────────────────────────────────────

// buildCommandCoT creates a minimal CoT XML event from a raw command JSON payload.
// For now we wrap the command payload in a <detail> element.
func buildCommandCoT(data []byte) []byte {
	now := time.Now().UTC().Format(time.RFC3339)
	stale := time.Now().UTC().Add(5 * time.Minute).Format(time.RFC3339)
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(
		`<?xml version="1.0"?>`+
			`<event version="2.0" uid="omos-cmd-%d" type="t-x-d-d" how="m-g" time="%s" start="%s" stale="%s">`,
		time.Now().UnixNano(), now, now, stale,
	))
	buf.WriteString(`<point lat="0" lon="0" hae="0" ce="9999999.0" le="9999999.0"/>`)
	buf.WriteString(`<detail>`)
	buf.Write(data)
	buf.WriteString(`</detail>`)
	buf.WriteString(`</event>`)
	return buf.Bytes()
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	cfg := loadConfig()
	log.Printf("[%s] starting — listen UDP :%d, NATS %s", cfg.bridgeName, cfg.listenPort, cfg.natsURL)

	// ── Connect to NATS ──────────────────────────────────────────────────────
	nc, err := nats.Connect(cfg.natsURL,
		nats.Name(cfg.bridgeName),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("[%s] NATS disconnected: %v", cfg.bridgeName, err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("[%s] NATS reconnected to %s", cfg.bridgeName, nc.ConnectedUrl())
		}),
	)
	if err != nil {
		log.Fatalf("[%s] NATS connect failed: %v", cfg.bridgeName, err)
	}
	defer nc.Drain()
	log.Printf("[%s] connected to NATS at %s", cfg.bridgeName, nc.ConnectedUrl())

	// ── UDP listener (inbound CoT from ATAK/OTS) ─────────────────────────────
	udpAddr := fmt.Sprintf("0.0.0.0:%d", cfg.listenPort)
	laddr, err := net.ResolveUDPAddr("udp4", udpAddr)
	if err != nil {
		log.Fatalf("[%s] resolve UDP addr %s: %v", cfg.bridgeName, udpAddr, err)
	}
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		log.Fatalf("[%s] listen UDP %s: %v", cfg.bridgeName, udpAddr, err)
	}
	defer conn.Close()
	log.Printf("[%s] listening for CoT on UDP %s", cfg.bridgeName, udpAddr)

	// errorSubject publishes errors to the bridge error subject.
	errorSubject := fmt.Sprintf("omos.bridge.%s.errors", cfg.bridgeName)
	publishError := func(msg string) {
		payload, _ := json.Marshal(map[string]interface{}{
			"source": cfg.bridgeName,
			"ts":     time.Now().UnixMilli(),
			"error":  msg,
		})
		if pubErr := nc.Publish(errorSubject, payload); pubErr != nil {
			log.Printf("[%s] failed to publish error: %v", cfg.bridgeName, pubErr)
		}
	}

	// ── UDP receive loop ──────────────────────────────────────────────────────
	go func() {
		buf := make([]byte, 65536)
		for {
			n, raddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				// Check if the connection was closed intentionally.
				if strings.Contains(err.Error(), "use of closed network connection") {
					return
				}
				log.Printf("[%s] UDP read error: %v", cfg.bridgeName, err)
				continue
			}
			raw := make([]byte, n)
			copy(raw, buf[:n])

			go func(data []byte, remote *net.UDPAddr) {
				payload, err := parseCoT(data)
				if err != nil {
					log.Printf("[%s] CoT parse error from %s: %v", cfg.bridgeName, remote, err)
					publishError(fmt.Sprintf("cot parse from %s: %v", remote, err))
					return
				}

				envBytes, subject, err := buildEnvelope(payload)
				if err != nil {
					log.Printf("[%s] envelope build error: %v", cfg.bridgeName, err)
					publishError(fmt.Sprintf("envelope build: %v", err))
					return
				}

				if err := nc.Publish(subject, envBytes); err != nil {
					log.Printf("[%s] NATS publish error to %s: %v", cfg.bridgeName, subject, err)
					publishError(fmt.Sprintf("nats publish to %s: %v", subject, err))
					return
				}
				log.Printf("[%s] published uid=%s callsign=%q → %s", cfg.bridgeName, payload.UID, payload.Callsign, subject)
			}(raw, raddr)
		}
	}()

	// ── NATS command subscriber (omos.field.*.command → TAK server) ──────────
	cmdSubject := "omos.field.*.command"
	takServerAddr := fmt.Sprintf("%s:%d", cfg.serverHost, cfg.serverPort)
	raddr, err := net.ResolveUDPAddr("udp4", takServerAddr)
	if err != nil {
		log.Fatalf("[%s] resolve TAK server addr %s: %v", cfg.bridgeName, takServerAddr, err)
	}

	sub, err := nc.Subscribe(cmdSubject, func(msg *nats.Msg) {
		cotXML := buildCommandCoT(msg.Data)

		udpConn, err := net.DialUDP("udp4", nil, raddr)
		if err != nil {
			log.Printf("[%s] dial TAK server %s: %v", cfg.bridgeName, takServerAddr, err)
			publishError(fmt.Sprintf("dial tak server: %v", err))
			return
		}
		defer udpConn.Close()

		if _, err := udpConn.Write(cotXML); err != nil {
			log.Printf("[%s] send command to TAK server: %v", cfg.bridgeName, err)
			publishError(fmt.Sprintf("send cmd to tak: %v", err))
			return
		}
		log.Printf("[%s] forwarded command from subject %s → TAK %s", cfg.bridgeName, msg.Subject, takServerAddr)
	})
	if err != nil {
		log.Fatalf("[%s] NATS subscribe to %s failed: %v", cfg.bridgeName, cmdSubject, err)
	}
	defer sub.Unsubscribe()
	log.Printf("[%s] subscribed to %s", cfg.bridgeName, cmdSubject)

	// ── Heartbeat every 30s ───────────────────────────────────────────────────
	healthSubject := fmt.Sprintf("omos.bridge.%s.health", cfg.bridgeName)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			beat, _ := json.Marshal(map[string]interface{}{
				"source": cfg.bridgeName,
				"ts":     time.Now().UnixMilli(),
				"status": "ok",
			})
			if err := nc.Publish(healthSubject, beat); err != nil {
				log.Printf("[%s] heartbeat publish error: %v", cfg.bridgeName, err)
			}
		}
	}()

	// ── Graceful shutdown on SIGTERM / SIGINT ─────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	log.Printf("[%s] received signal %s — shutting down", cfg.bridgeName, sig)

	// Close UDP listener so the receive goroutine exits cleanly.
	conn.Close()

	// Drain NATS (flushes outbound and unsubscribes).
	if err := nc.Drain(); err != nil {
		log.Printf("[%s] NATS drain error: %v", cfg.bridgeName, err)
	}
	log.Printf("[%s] shutdown complete", cfg.bridgeName)
}
