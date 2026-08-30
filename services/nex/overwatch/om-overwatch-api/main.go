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

// om-overwatch-api — typed entity CRUD for om-overwatch (organ 8/8).
//
// The manual-input organ. Lets atoc and external tools create/edit/delete
// entities directly — no-fly zones, waypoints, manually-dropped tracks, renames.
//
// Architecture: the API never writes KV directly. om-overwatch-state stays the
// single KV writer. Writes publish intents the spine consumes:
//   POST   /api/v1/entities        -> publish ent.manual.upsert  (spine writes KV)
//   DELETE /api/v1/entities/{id}   -> publish ent.manual.delete  (spine deletes KV)
// Reads go straight to the KV bucket (read-only, no contention):
//   GET    /api/v1/entities        -> list all from OVERWATCH_STATE
//   GET    /api/v1/entities/{id}   -> one entity
//
// Manual entities are pinned by the spine (exempt from the staleness sweep).

const (
	kvBucket  = "OVERWATCH_STATE"
	hbSubject = "agents.hb.overwatch.api.zeus"
)

type server struct {
	nc *nats.Conn
	kv nats.KeyValue
}

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = "nats://127.0.0.1:4222"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8091"
	}

	nc, err := nats.Connect(url, nats.Name("om-overwatch-api"),
		nats.MaxReconnects(-1), nats.ReconnectWait(2*time.Second))
	if err != nil {
		log.Fatal("NATS connect:", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("JetStream:", err)
	}
	// The spine creates the bucket; wait for it if the API booted first.
	var kv nats.KeyValue
	for i := 0; i < 15; i++ {
		kv, err = js.KeyValue(kvBucket)
		if err == nil {
			break
		}
		log.Printf("om-overwatch-api: waiting for %s bucket (spine)...", kvBucket)
		time.Sleep(2 * time.Second)
	}
	if kv == nil {
		log.Fatal("KV bucket unavailable:", err)
	}

	s := &server{nc: nc, kv: kv}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/entities", s.handleEntities)
	mux.HandleFunc("/api/v1/entities/", s.handleEntityByID)
	mux.HandleFunc("/health", s.handleHealth)

	log.Printf("om-overwatch-api: HTTP CRUD on :%s, bus %s", port, url)
	nc.Publish(hbSubject, []byte(`{"status":"up"}`))
	go heartbeat(nc)

	if err := http.ListenAndServe(":"+port, cors(mux)); err != nil {
		log.Fatal(err)
	}
}

func heartbeat(nc *nats.Conn) {
	t := time.NewTicker(15 * time.Second)
	for range t.C {
		nc.Publish(hbSubject, []byte(`{"t":"`+time.Now().UTC().Format(time.RFC3339)+`","status":"serving"}`))
	}
}

// entityInput is the typed create/update payload.
type entityInput struct {
	EntityID    string                 `json:"entity_id"`
	Name        string                 `json:"name"`
	EntityType  string                 `json:"entity_type"`
	Affiliation string                 `json:"affiliation"`
	Priority    string                 `json:"priority"`
	Position    *position              `json:"position,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude,omitempty"`
	Heading   float64 `json:"heading,omitempty"`
	Velocity  float64 `json:"velocity,omitempty"`
}

// validEntityTypes mirrors the spine ontology.
var validEntityTypes = map[string]bool{
	"aircraft_fixed_wing": true, "aircraft_multirotor": true, "aircraft_vtol": true,
	"aircraft_helicopter": true, "aircraft_unknown": true, "ground_vehicle": true,
	"surface_vessel": true, "sensor_platform": true, "operator_station": true,
	"interceptor": true, "waypoint": true, "no_fly_zone": true, "geofence": true,
}

func (s *server) handleEntities(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listEntities(w, r)
	case http.MethodPost:
		s.upsertEntity(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *server) handleEntityByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/entities/")
	if id == "" {
		http.Error(w, "missing entity id", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		entry, err := s.kv.Get(id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(entry.Value())
	case http.MethodDelete:
		if err := s.nc.Publish("ent.manual.delete", []byte(id)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]any{"entity_id": id, "deleted": true})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *server) upsertEntity(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var in entityInput
	if err := json.Unmarshal(body, &in); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if in.EntityType == "" || !validEntityTypes[in.EntityType] {
		http.Error(w, "invalid or missing entity_type", http.StatusBadRequest)
		return
	}
	if in.EntityID == "" {
		in.EntityID = fmt.Sprintf("man-%d", time.Now().UnixNano())
	}
	if in.Affiliation == "" {
		in.Affiliation = "friendly"
	}
	if in.Priority == "" {
		in.Priority = "normal"
	}
	// Publish the intent; the spine performs the KV write.
	payload, _ := json.Marshal(in)
	if err := s.nc.Publish("ent.manual.upsert", payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"entity_id": in.EntityID, "accepted": true})
}

func (s *server) listEntities(w http.ResponseWriter, r *http.Request) {
	keys, err := s.kv.Keys()
	if err != nil {
		// empty bucket returns an error in nats.go; treat as empty list
		writeJSON(w, http.StatusOK, []any{})
		return
	}
	out := make([]json.RawMessage, 0, len(keys))
	for _, k := range keys {
		if entry, err := s.kv.Get(k); err == nil {
			out = append(out, json.RawMessage(entry.Value()))
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "nats": s.nc.IsConnected()})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}
