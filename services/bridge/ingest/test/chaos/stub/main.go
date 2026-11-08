// chaos-stub imitates two upstreams for the omos-ingest chaos test, with synthetic data:
//   - GET /alerts/active: an NWS-shaped FeatureCollection listing every alert generated in
//     the last -window seconds, at -rate alerts/s (ETag, If-None-Match, 503 when down);
//   - GET /sse: an event stream at -sse-rate events/s that replays from Last-Event-ID.
//
// /ctl/outage?on=1|0 makes both return 503; /ctl/malformed?on=1|0 mixes in records that
// fail the schema. Every emitted id is logged as JSON lines to -log so the verifier can
// prove no loss and no duplicates.
package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type item struct {
	n    int
	at   time.Time
	bad  bool
	json []byte
}

type gen struct {
	mu    sync.Mutex
	items []item // nws: ring of recent alerts; sse: full history (bounded)
	next  int
}

var (
	outage, malformed atomic.Bool
	logMu             sync.Mutex
	logW              *bufio.Writer
)

func emit(src string, n int, bad bool) {
	logMu.Lock()
	defer logMu.Unlock()
	fmt.Fprintf(logW, `{"t":%q,"src":%q,"id":%q,"bad":%t}`+"\n", time.Now().UTC().Format(time.RFC3339Nano), src, idFor(src, n), bad)
}

func idFor(src string, n int) string {
	if src == "nws" {
		return fmt.Sprintf("nws:urn:oid:chaos.%d", n)
	}
	return fmt.Sprintf("sse:e%d", n)
}

func nwsFeature(n int, now time.Time, bad bool) []byte {
	sev := "Minor"
	if bad {
		sev = "Bogus" // fails the weather-alert schema enum
	}
	f := map[string]any{"id": fmt.Sprintf("https://stub/alerts/urn:oid:chaos.%d", n), "type": "Feature", "geometry": nil,
		"properties": map[string]any{"id": fmt.Sprintf("urn:oid:chaos.%d", n), "areaDesc": "Chaos County", "sent": now.Format(time.RFC3339),
			"effective": now.Format(time.RFC3339), "expires": now.Add(time.Hour).Format(time.RFC3339), "status": "Test", "messageType": "Alert",
			"category": "Met", "severity": sev, "certainty": "Likely", "urgency": "Expected", "event": "Test Warning",
			"headline": fmt.Sprintf("synthetic alert %d", n), "affectedZones": []string{}, "geocode": map[string]any{"UGC": []string{"XXC001"}}}}
	b, _ := json.Marshal(f)
	return b
}

func sseData(n int, bad bool) []byte {
	lat := 10.0 + float64(n%70)
	if bad {
		lat = 999 // fails the generic-record schema
	}
	return []byte(fmt.Sprintf(`{"id":"e%d","lat":%g,"lon":20.5}`, n, lat))
}

func main() {
	addr := flag.String("listen", "127.0.0.1:38180", "listen address")
	rate := flag.Int("rate", 50, "nws alerts per second")
	window := flag.Duration("window", 30*time.Second, "how long an alert stays listed")
	sseRate := flag.Int("sse-rate", 50, "sse events per second")
	logPath := flag.String("log", "emitted.jsonl", "emitted id log")
	flag.Parse()
	f, err := os.Create(*logPath)
	if err != nil {
		log.Fatal(err)
	}
	logW = bufio.NewWriterSize(f, 1<<16)
	go func() {
		for range time.Tick(500 * time.Millisecond) {
			logMu.Lock()
			_ = logW.Flush()
			logMu.Unlock()
		}
	}()

	nws, sse := &gen{}, &gen{}
	tick := func(g *gen, src string, perSec int, mk func(int, time.Time, bool) []byte, keep func([]item, time.Time) []item) {
		for now := range time.Tick(time.Second / time.Duration(perSec)) {
			bad := malformed.Load() && g.next%5 == 0
			g.mu.Lock()
			g.next++
			n := g.next
			g.items = keep(append(g.items, item{n: n, at: now, bad: bad, json: mk(n, now, bad)}), now)
			g.mu.Unlock()
			emit(src, n, bad)
		}
	}
	go tick(nws, "nws", *rate, nwsFeature, func(it []item, now time.Time) []item {
		i := 0
		for i < len(it) && now.Sub(it[i].at) > *window {
			i++
		}
		return it[i:]
	})
	go tick(sse, "sse", *sseRate, func(n int, _ time.Time, bad bool) []byte { return sseData(n, bad) }, func(it []item, _ time.Time) []item {
		if len(it) > 200000 {
			return it[len(it)-200000:]
		}
		return it
	})

	http.HandleFunc("/ctl/outage", func(w http.ResponseWriter, r *http.Request) { outage.Store(r.URL.Query().Get("on") == "1") })
	http.HandleFunc("/ctl/malformed", func(w http.ResponseWriter, r *http.Request) { malformed.Store(r.URL.Query().Get("on") == "1") })
	http.HandleFunc("/alerts/active", func(w http.ResponseWriter, r *http.Request) {
		if outage.Load() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "upstream down", http.StatusServiceUnavailable)
			return
		}
		nws.mu.Lock()
		items := append([]item(nil), nws.items...)
		nws.mu.Unlock()
		h := sha1.New()
		for _, it := range items {
			fmt.Fprint(h, it.n, ",")
		}
		tag := fmt.Sprintf(`W/"%x"`, h.Sum(nil))
		if r.Header.Get("If-None-Match") == tag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("ETag", tag)
		w.Header().Set("Content-Type", "application/geo+json")
		bw := bufio.NewWriter(w)
		bw.WriteString(`{"type":"FeatureCollection","features":[`)
		for i := len(items) - 1; i >= 0; i-- { // newest first, as api.weather.gov lists
			bw.Write(items[i].json)
			if i > 0 {
				bw.WriteByte(',')
			}
		}
		bw.WriteString(`]}`)
		_ = bw.Flush()
	})
	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		if outage.Load() {
			http.Error(w, "upstream down", http.StatusServiceUnavailable)
			return
		}
		from, _ := strconv.Atoi(r.Header.Get("Last-Event-ID"))
		w.Header().Set("Content-Type", "text/event-stream")
		fl := w.(http.Flusher)
		for {
			if outage.Load() {
				return // drop the stream: the client must resume with Last-Event-ID
			}
			sse.mu.Lock()
			var batch []item
			for _, it := range sse.items {
				if it.n > from {
					batch = append(batch, it)
				}
			}
			sse.mu.Unlock()
			for _, it := range batch {
				if _, err := fmt.Fprintf(w, "id: %d\ndata: %s\n\n", it.n, it.json); err != nil {
					return
				}
				from = it.n
			}
			fl.Flush()
			select {
			case <-r.Context().Done():
				return
			case <-time.After(100 * time.Millisecond):
			}
		}
	})
	log.Printf("chaos-stub on %s (nws %d/s window %v, sse %d/s)", *addr, *rate, *window, *sseRate)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
