package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type Workload struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	State     string    `json:"state"`
	Type      string    `json:"type"`
	Lifecycle string    `json:"lifecycle"`
	StartTime time.Time `json:"start_time"`
	Node      string    `json:"node"`
	LastSeen  time.Time `json:"last_seen"`
}

type NodeInfo struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Version  string    `json:"version"`
	LastSeen time.Time `json:"last_seen"`
}

type BusEvent struct {
	Subject string          `json:"subject"`
	Data    json.RawMessage `json:"data"`
	Time    time.Time       `json:"time"`
}

type State struct {
	mu        sync.RWMutex
	Workloads map[string]*Workload `json:"workloads"`
	Nodes     map[string]*NodeInfo `json:"nodes"`
	Events    []BusEvent           `json:"events"`
	SubjectHz map[string]int       `json:"subject_hz"`
}

var state = &State{
	Workloads: make(map[string]*Workload),
	Nodes:     make(map[string]*NodeInfo),
	Events:    []BusEvent{},
	SubjectHz: make(map[string]int),
}

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://127.0.0.1:4222"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3301"
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal("NATS connect:", err)
	}
	defer nc.Close()
	log.Printf("om-nex-ui connected to %s", natsURL)

	// Subscribe to NEX events
	nc.Subscribe("$NEX.>", func(msg *nats.Msg) {
		state.mu.Lock()
		defer state.mu.Unlock()
		var data json.RawMessage = msg.Data
		state.Events = append([]BusEvent{{Subject: msg.Subject, Data: data, Time: time.Now()}}, state.Events...)
		if len(state.Events) > 200 {
			state.Events = state.Events[:200]
		}
		// Parse workload events
		var evt map[string]interface{}
		if json.Unmarshal(msg.Data, &evt) == nil {
			if id, ok := evt["workload_id"].(string); ok && id != "" {
				name, _ := evt["workload_name"].(string)
				wl := state.Workloads[id]
				if wl == nil {
					wl = &Workload{ID: id, Name: name}
					state.Workloads[id] = wl
				}
				if s, ok := evt["state"].(string); ok {
					wl.State = s
				}
				wl.LastSeen = time.Now()
			}
			if nodeID, ok := evt["node_id"].(string); ok && nodeID != "" {
				node := state.Nodes[nodeID]
				if node == nil {
					node = &NodeInfo{ID: nodeID}
					state.Nodes[nodeID] = node
				}
				if v, ok := evt["version"].(string); ok {
					node.Version = v
				}
				node.LastSeen = time.Now()
			}
		}
	})

	// Subscribe to all agent heartbeats
	nc.Subscribe("agents.hb.>", func(msg *nats.Msg) {
		state.mu.Lock()
		state.SubjectHz[msg.Subject]++
		state.mu.Unlock()
	})

	// Reset Hz counters every second
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			state.mu.Lock()
			state.SubjectHz = make(map[string]int)
			state.mu.Unlock()
		}
	}()

	// HTTP server
	mux := http.NewServeMux()

	// API
	mux.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		state.mu.RLock()
		defer state.mu.RUnlock()
		json.NewEncoder(w).Encode(state)
	})

	// Dashboard HTML
	mux.HandleFunc("/", serveDashboard)

	log.Printf("om-nex-ui dashboard at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func serveDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, dashHTML)
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>OneMind NEX Dashboard</title>
<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { background: #0a0e1a; color: #e2e8f0; font-family: 'JetBrains Mono', 'Courier New', monospace; }
.header { background: #0d1117; border-bottom: 1px solid #1a2744; padding: 12px 24px; display: flex; align-items: center; justify-content: space-between; }
.logo { font-size: 18px; font-weight: 700; color: #60a5fa; letter-spacing: 2px; }
.logo span { color: #f59e0b; }
.status-dot { width: 8px; height: 8px; border-radius: 50%; background: #22c55e; box-shadow: 0 0 6px #22c55e; display: inline-block; margin-right: 6px; animation: pulse 2s infinite; }
@keyframes pulse { 0%,100%{opacity:1} 50%{opacity:0.5} }
.grid { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16px; padding: 20px; }
.card { background: #0d1117; border: 1px solid #1a2744; border-radius: 8px; padding: 16px; }
.card-title { font-size: 11px; text-transform: uppercase; letter-spacing: 2px; color: #60a5fa; margin-bottom: 12px; }
.workload-row { display: flex; align-items: center; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #1a2744; }
.workload-row:last-child { border-bottom: none; }
.wl-name { color: #f1f5f9; font-size: 13px; }
.wl-id { color: #475569; font-size: 10px; }
.badge { padding: 2px 8px; border-radius: 4px; font-size: 10px; font-weight: 600; text-transform: uppercase; }
.badge-running { background: #052e16; color: #22c55e; border: 1px solid #166534; }
.badge-stopped { background: #450a0a; color: #f87171; border: 1px solid #7f1d1d; }
.badge-pending { background: #1c1917; color: #a8a29e; border: 1px solid #44403c; }
.event-row { padding: 4px 0; border-bottom: 1px solid #1a2744; font-size: 11px; }
.event-subject { color: #60a5fa; }
.event-time { color: #475569; float: right; }
.big-number { font-size: 36px; font-weight: 700; color: #f59e0b; }
.big-label { font-size: 11px; color: #64748b; margin-top: 4px; }
.stat-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.stat-box { text-align: center; }
.full-row { grid-column: 1 / -1; }
.events-list { max-height: 300px; overflow-y: auto; }
.events-list::-webkit-scrollbar { width: 4px; }
.events-list::-webkit-scrollbar-track { background: #0d1117; }
.events-list::-webkit-scrollbar-thumb { background: #1a2744; }
.node-row { padding: 6px 0; border-bottom: 1px solid #1a2744; display: flex; justify-content: space-between; }
.refresh-info { font-size: 10px; color: #475569; }
.stack-pill { display: inline-block; padding: 3px 10px; border-radius: 12px; font-size: 10px; margin: 3px; background: #0f172a; border: 1px solid #1e3a5f; color: #93c5fd; }
.stack-pill.active { background: #052e16; border-color: #166534; color: #4ade80; }
</style>
</head>
<body>
<div class="header">
  <div class="logo">ONE<span>MIND</span> NEX</div>
  <div><span class="status-dot"></span><span style="font-size:12px;color:#94a3b8">LIVE · NATS fabric connected</span></div>
  <div class="refresh-info" id="refreshInfo">refreshing...</div>
</div>
<div class="grid">
  <!-- Stats row -->
  <div class="card">
    <div class="card-title">Workloads</div>
    <div class="stat-grid">
      <div class="stat-box"><div class="big-number" id="wlRunning">0</div><div class="big-label">Running</div></div>
      <div class="stat-box"><div class="big-number" style="color:#f87171" id="wlStopped">0</div><div class="big-label">Stopped</div></div>
    </div>
  </div>
  <div class="card">
    <div class="card-title">NEX Nodes</div>
    <div class="stat-grid">
      <div class="stat-box"><div class="big-number" id="nodeCount">0</div><div class="big-label">Active Nodes</div></div>
      <div class="stat-box"><div class="big-number" style="color:#a78bfa" id="eventCount">0</div><div class="big-label">Bus Events</div></div>
    </div>
  </div>
  <div class="card">
    <div class="card-title">Active Stack</div>
    <div id="stackPills"></div>
  </div>
  <!-- Workloads -->
  <div class="card full-row">
    <div class="card-title">Live NEX Workloads <span style="color:#475569;font-size:10px">(from $NEX.> subjects)</span></div>
    <div id="workloadList"></div>
  </div>
  <!-- Nodes -->
  <div class="card">
    <div class="card-title">NEX Nodes</div>
    <div id="nodeList"></div>
  </div>
  <!-- Events -->
  <div class="card" style="grid-column: 2 / -1">
    <div class="card-title">Live Bus Events <span style="color:#475569;font-size:10px">(last 200)</span></div>
    <div class="events-list" id="eventList"></div>
  </div>
</div>
<script>
const STACK = [
  {name:'nex', active:true}, {name:'om-task', active:true},
  {name:'om-meta', active:true}, {name:'om-streaming-bridge', active:true},
  {name:'orbit.go', active:true}, {name:'callout.go', active:false},
  {name:'agenttap', active:true}, {name:'bento', active:true},
  {name:'rethink_conn', active:true}, {name:'om-nex-ui', active:true},
];

async function refresh() {
  try {
    const r = await fetch('/api/state');
    const d = await r.json();
    
    const wls = Object.values(d.workloads || {});
    const running = wls.filter(w => w.state === 'running' || w.state === '').length;
    document.getElementById('wlRunning').textContent = running;
    document.getElementById('wlStopped').textContent = wls.filter(w => w.state === 'stopped').length;
    document.getElementById('nodeCount').textContent = Object.keys(d.nodes || {}).length;
    document.getElementById('eventCount').textContent = (d.events || []).length;
    
    // Workloads
    const wlEl = document.getElementById('workloadList');
    if (wls.length === 0) {
      wlEl.innerHTML = '<div style="color:#475569;font-size:12px;padding:8px">No workload events yet — events appear as $NEX.> messages arrive. Use: nats workload ls</div>';
    } else {
      wlEl.innerHTML = wls.map(w => {
        const badge = w.state === 'stopped' ? 'badge-stopped' : w.state ? 'badge-' + w.state : 'badge-running';
        return '<div class="workload-row"><div><div class="wl-name">' + w.name + '</div><div class="wl-id">' + w.id + '</div></div><span class="badge ' + badge + '">' + (w.state || 'running') + '</span></div>';
      }).join('');
    }
    
    // Nodes
    const nodes = Object.values(d.nodes || {});
    const nodeEl = document.getElementById('nodeList');
    nodeEl.innerHTML = nodes.length === 0 
      ? '<div style="color:#475569;font-size:12px;padding:8px">Listening on $NEX.> for node events...</div>'
      : nodes.map(n => '<div class="node-row"><span style="color:#93c5fd">' + (n.name || n.id?.substring(0,16)) + '</span><span class="badge badge-running">online</span></div>').join('');
    
    // Events
    const events = d.events || [];
    document.getElementById('eventList').innerHTML = events.slice(0,50).map(e => 
      '<div class="event-row"><span class="event-subject">' + e.subject + '</span><span class="event-time">' + new Date(e.time).toLocaleTimeString() + '</span></div>'
    ).join('') || '<div style="color:#475569;font-size:12px;padding:8px">Waiting for $NEX.> traffic...</div>';
    
    // Stack pills
    document.getElementById('stackPills').innerHTML = STACK.map(s => 
      '<span class="stack-pill' + (s.active ? ' active' : '') + '">' + s.name + '</span>'
    ).join('');
    
    document.getElementById('refreshInfo').textContent = 'updated ' + new Date().toLocaleTimeString();
  } catch(e) { document.getElementById('refreshInfo').textContent = 'fetch error: ' + e.message; }
}
refresh();
setInterval(refresh, 2000);
</script>
</body>
</html>`
