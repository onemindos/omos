# NATS 3-Node Cluster — One Mind Fabric

**Status:** Design doc — pre-implementation  
**Created:** 2026-06-17  
**Author:** Zeus Delacruz + Legacy  
**Location:** Vault root (temp) → move to `00-24 SO (Sovereign Operations)/03 Infrastructure (The Grid)/` when ready to implement

---

## Why This Exists

Your MBP is currently a **single point of failure** for the entire One Mind NATS fabric. Everything — ha-bridge, n8n, Hermes agents, agenttap, all automation — connects to one `nats-server` process on one machine. If the MBP crashes, sleeps, loses power, or is taken for repair, **the entire nervous system goes silent.**

That is unacceptable before you demo this to anyone serious, and especially before you go public with a platform.

**This doc describes the fix: a 3-node JetStream cluster.**

---

## How JetStream Clustering Works

NATS uses **Raft consensus** for JetStream clustering. Raft is a voting protocol:

- Every write to JetStream requires a **majority** of nodes to acknowledge it before it's confirmed
- Majority = ⌊N/2⌋ + 1

| Cluster Size | Majority Needed | Can Lose | Verdict |
|---|---|---|---|
| 1 (today) | 1 | 0 | SPOF — one failure = down |
| 2 (the trap) | 2 | 0 | Worse than 1 — two things can kill you |
| **3 (answer)** | **2** | **1** | **Lose any one node — cluster keeps running ✅** |
| 5 | 3 | 2 | Overkill until fleet is large |

**3 is the minimum viable HA configuration.** Never run 2.

---

## The Topology

```
┌─────────────────────────────────────────────────────────────────┐
│                   One Mind NATS Cluster                          │
│                    domain: onemind                               │
├──────────────────┬──────────────────┬───────────────────────────┤
│   Node A         │   Node B         │   Node C                  │
│   MBP M2 Max     │   Second machine │   Tiebreaker              │
│   PRIMARY HUB    │   REPLICA HUB    │   VOTE-ONLY               │
│                  │                  │                            │
│   Full load      │   Full replica   │   Lightweight              │
│   Primary writes │   Takes over if  │   Just needs to be online  │
│                  │   Node A dies    │   to cast deciding vote     │
└──────────────────┴──────────────────┴───────────────────────────┘
         │                  │                    │
         └──────── Tailscale private network ────┘
                    (cluster gossip port 6222)
```

---

## Node Options

### Node A — MBP M2 Max (already running)
No changes to hardware. Config update only.

### Node B — Second Machine Options

| Option | Cost | Notes |
|---|---|---|
| **Mac Mini (used)** | ~$200 | Best — same OS, same stack, silent, low power |
| **Mac Mini (new M4)** | ~$600 | Future-proof, worth it if budget allows |
| **Linux server / NUC** | ~$150–300 | Works — NATS is cross-platform |
| **Raspberry Pi 5 (8GB)** | ~$80 | Capable as Node B if it has SSD, not SD card |

**Recommendation:** Used Mac Mini. Same OS, same brew ecosystem, silent, ~15W idle.

### Node C — Tiebreaker Options

| Option | Cost | Notes |
|---|---|---|
| **HA box (if beefy enough)** | $0 | Check specs first — Pi-class = tiebreaker only |
| **Raspberry Pi 4/5** | ~$50–80 | Tiebreaker only — low storage fine |
| **$5/mo VPS (Hetzner/Vultr)** | ~$60/yr | Remote over Tailscale — adds latency to votes but works |

**Cluster gossip (port 6222) MUST stay on LAN or Tailscale. Never expose to public internet.**

---

## Configuration

Each node runs `nats-server` with a `cluster {}` block. All three nodes point at each other.

### Node A (MBP) — `nats-server.conf`

```hcl
server_name: onemind-node-a

jetstream {
  domain: onemind
  store_dir: /Users/onemind/.nats/jetstream
}

cluster {
  name: onemind-cluster
  listen: 0.0.0.0:6222
  routes: [
    nats-route://NODE_B_TAILSCALE_IP:6222
    nats-route://NODE_C_TAILSCALE_IP:6222
  ]
}

# Client connections (unchanged)
listen: 0.0.0.0:4222

# WebSocket
websocket {
  listen: 0.0.0.0:8080
  no_tls: true
}

# MQTT bridge
mqtt {
  listen: 0.0.0.0:1883
}

# Leaf node port (for edge devices)
leafnodes {
  listen: 0.0.0.0:7422
}
```

### Node B — `nats-server.conf`

```hcl
server_name: onemind-node-b

jetstream {
  domain: onemind
  store_dir: /data/nats/jetstream
}

cluster {
  name: onemind-cluster
  listen: 0.0.0.0:6222
  routes: [
    nats-route://NODE_A_TAILSCALE_IP:6222
    nats-route://NODE_C_TAILSCALE_IP:6222
  ]
}

listen: 0.0.0.0:4222
```

### Node C (Tiebreaker) — `nats-server.conf`

```hcl
server_name: onemind-node-c

# Minimal JetStream — just enough to vote
jetstream {
  domain: onemind
  store_dir: /data/nats/jetstream
  max_memory: 256MB
  max_file: 2GB
}

cluster {
  name: onemind-cluster
  listen: 0.0.0.0:6222
  routes: [
    nats-route://NODE_A_TAILSCALE_IP:6222
    nats-route://NODE_B_TAILSCALE_IP:6222
  ]
}

listen: 0.0.0.0:4222
```

---

## Client Configuration

All NATS clients (ha-bridge, n8n, Hermes agents, omos-bridge processes) list **all three node URLs**. NATS client auto-discovers topology and reconnects to survivors.

```bash
# .env / environment variable
NATS_URL=nats://NODE_A_IP:4222,nats://NODE_B_IP:4222,nats://NODE_C_IP:4222
```

In code (Go):
```go
nc, err := nats.Connect(
    "nats://NODE_A_IP:4222,nats://NODE_B_IP:4222,nats://NODE_C_IP:4222",
    nats.MaxReconnects(-1),
)
```

---

## Stream Replication

After cluster is up, set critical streams to R3 (replicated to all 3 nodes):

```bash
# Set existing streams to R3
nats stream edit AGENT_EVENTS --replicas 3
nats stream edit SERVICE_EVENTS --replicas 3
nats stream edit SYSTEM_EVENTS --replicas 3
nats stream edit HUMAN_EVENTS --replicas 3

# Verify
nats stream info AGENT_EVENTS
# Look for: Replication: 3
```

---

## Failure Scenarios

### MBP dies at 3am

```
Node B + Node C vote → Node B elected leader in < 500ms
All clients auto-reconnect (they have all 3 URLs)
Zero data loss — R3 streams have copies on B and C
MBP comes back → rejoins cluster → syncs automatically
You wake up → nothing was lost
```

### MBP + Node C both die

```
Node B alone = no majority (1 of 3)
JetStream pauses NEW writes (safety — won't accept data it can't replicate)
Existing data is intact on Node B
Core NATS pub/sub still works (only JetStream persistence pauses)
Recover: bring either node back → majority restored → writes resume
```

### All 3 nodes die

```
Total outage — same as today except data is preserved on all 3 when they come back
Recovery: start any node → starts as standalone → others rejoin → cluster reforms
```

---

## Implementation Checklist

- [ ] Acquire Node B hardware (Mac Mini or equivalent)
- [ ] Acquire Node C hardware (Pi 5 or VPS)
- [ ] Add all 3 to Tailscale network
- [ ] Record Tailscale IPs for all 3 nodes
- [ ] Update MBP `nats-server.conf` with cluster block
- [ ] Configure Node B `nats-server.conf`
- [ ] Configure Node C `nats-server.conf`
- [ ] Start all 3 NATS servers
- [ ] Verify cluster: `nats server list` — should show 3 nodes
- [ ] Update all client `NATS_URL` env vars to include all 3
- [ ] Set critical streams to R3
- [ ] Test failover: stop Node A, verify clients reconnect to Node B
- [ ] Update launchd plists on MBP to use new config
- [ ] Document Tailscale IPs in 1Password vault

---

## Verification Commands

```bash
# Check cluster health
nats server list
nats server report jetstream

# Check stream replication
nats stream info AGENT_EVENTS

# Watch cluster leader elections (run this before killing a node)
nats server report connections

# Simulate failover test
# 1. On MBP: sudo launchctl stop io.nats.server
# 2. Watch: nats server list — Node B should appear as leader
# 3. Verify clients reconnected
# 4. Restart: sudo launchctl start io.nats.server
# 5. Verify Node A rejoins
```

---

## Notes

- Cluster gossip port **6222 stays internal** — Tailscale only, never exposed
- Client port 4222 can stay LAN — same as today
- MQTT, WebSocket, and leaf node ports only needed on Node A (the primary edge ingress)
- Node B and C are cluster members only — they don't need MQTT/WS unless you want full HA on those ports too
- JetStream domain stays `onemind` across all 3 nodes

---

*Move this doc to `00-24 SO (Sovereign Operations)/03 Infrastructure (The Grid)/NATS-3-NODE-CLUSTER.md` when implementation begins.*
