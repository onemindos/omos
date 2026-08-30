package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

// om-http-gateway: HTTP webhook → NATS bridge
// Any HTTP POST to /pub/{subject.path} publishes to that NATS subject
// GET /sub/{subject} returns SSE stream of that subject
// Used for: ADS-B feeds, weather APIs, world events, external webhooks

var nc *nats.Conn

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://127.0.0.1:4222"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	var err error
	nc, err = nats.Connect(natsURL, nats.Name("om-http-gateway"))
	if err != nil {
		log.Fatal("NATS connect:", err)
	}
	defer nc.Close()
	log.Printf("om-http-gateway: connected to NATS, listening on :%s", port)

	mux := http.NewServeMux()

	// POST /pub/tel.world.adsb → publish to tel.world.adsb
	mux.HandleFunc("/pub/", handlePublish)

	// GET /sub/tel.world.adsb → SSE stream
	mux.HandleFunc("/sub/", handleSubscribe)

	// GET /health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status": "ok", "nats": nc.IsConnected(), "t": time.Now(),
		})
	})

	// Info UI
	mux.HandleFunc("/", serveInfo)

	nc.Publish("agents.hb.gateway.http.zeus", []byte(`{"status":"up","port":"`+port+`"}`))
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func handlePublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	subject := strings.TrimPrefix(r.URL.Path, "/pub/")
	subject = strings.ReplaceAll(subject, "/", ".")
	if subject == "" {
		http.Error(w, "subject required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := nc.Publish(subject, body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("om-http-gateway: published %d bytes to %s", len(body), subject)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"published": subject, "bytes": len(body)})
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	subject := strings.TrimPrefix(r.URL.Path, "/sub/")
	subject = strings.ReplaceAll(subject, "/", ".")
	if subject == "" {
		http.Error(w, "subject required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Fprintf(w, "data: %s\n\n", string(msg.Data))
		flusher.Flush()
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer sub.Unsubscribe()

	log.Printf("om-http-gateway: SSE client connected to %s", subject)
	<-r.Context().Done()
	log.Printf("om-http-gateway: SSE client disconnected from %s", subject)
}

func serveInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>om-http-gateway</title>
<style>body{background:#0a0e1a;color:#e2e8f0;font-family:monospace;padding:40px}
code{background:#1a2744;padding:4px 8px;border-radius:4px;color:#60a5fa}</style></head>
<body><h1 style="color:#60a5fa">om-http-gateway</h1>
<p>HTTP ↔ NATS bridge for OneMind OS</p>
<h3>Publish (POST)</h3>
<code>POST /pub/tel.world.adsb</code> → publishes body to <code>tel.world.adsb</code><br><br>
<code>POST /pub/evt.world.weather</code> → publishes to <code>evt.world.weather</code><br><br>
<h3>Subscribe (SSE)</h3>
<code>GET /sub/tel.drone.wing-alpha.gps</code> → SSE stream of drone GPS<br><br>
<code>GET /sub/tel.robot.spot-alpha.odom</code> → SSE stream of robot odometry<br><br>
<h3>Health</h3>
<code>GET /health</code> → JSON status<br>
<hr><p style="color:#475569">Connected to NATS: %v | Port: %s</p>
</body></html>`, nc.IsConnected(), os.Getenv("PORT"))
}
