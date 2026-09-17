---
title: "Limits & Defaults"
section: "14-REFERENCE"
type: note
tags: [nats, limits, defaults, ports, payload, timeouts, jetstream, reference]
source: https://docs.nats.io/running-a-nats-service/configuration
created: 2026-06-14
---

# Limits & Defaults

Every port, cap, timeout, and default value in one place. **Version-sensitive** — these reflect the NATS server v2.14.x line (2026-06-14). When in doubt, confirm with `nats-server --help` and `/varz`.

See also: [[NATS/14-REFERENCE/config-cheatsheet|Config Cheatsheet]] · [[NATS/14-REFERENCE/monitoring-endpoints|Monitoring Endpoints]] · [[NATS/07-DEPLOYMENT-OPS/configuration-reference|Configuration]]

---

## Ports

| Port | Purpose | Config |
|------|---------|--------|
| **4222** | Client connections | `port` |
| **6222** | Cluster routes (server↔server) | `cluster.listen` |
| **7222** | Gateways (cluster↔cluster) | `gateway.listen` |
| **7422** | Leaf node connections | `leafnodes.listen` |
| **8222** | HTTP monitoring | `http` / `https` |
| **1883** | MQTT (convention) | `mqtt.listen` |
| **8883** | MQTTS (convention) | `mqtt.listen` + tls |
| **8080** | WebSocket (convention) | `websocket.listen` |

> Ports above 4222/6222/7222/7422/8222 are the **defaults/conventions**; MQTT/WS ports are conventional (you set them). Mnemonic: client **4**222, routes **6**222, gateway **7**2**22**, leaf **74**22.

---

## Message & connection limits

| Setting | Default | Meaning |
|---------|---------|---------|
| `max_payload` | **1 MB** | Max size of a single message. Hard cap; clients learn it via INFO. |
| `max_pending` (client-side) | 64 MB | Client write buffer before slow-consumer. |
| `max_control_line` | 4 KB | Max protocol control line length. |
| `max_connections` | **64K** | Max simultaneous client connections. |
| `max_subscriptions` | 0 (unlimited) | Per-connection subscription cap (account-settable). |
| `ping_interval` | **2 min** | Server→client PING cadence. |
| `ping_max` (`max_outstanding_pings`) | **2** | Missed PONGs before disconnect. |
| `write_deadline` | **10 s** | Max time to flush to a client before slow-consumer. |
| `auth timeout` | 1–2 s | Time allowed to authenticate after connect. |
| `tls timeout` | 2 s | TLS handshake timeout. |

> **`max_payload` caps at 64 MB** even if you raise it — and large messages hurt latency. For big blobs use the **Object Store** (chunked), not raw payloads. See [[NATS/03-KV-OBJECT-STORE/object-store|Object Store]].

---

## Timeouts & intervals (client defaults — vary by SDK)

| Setting | Typical default | Meaning |
|---------|----------------|---------|
| Connect timeout | 2 s | TCP+handshake connect. |
| Reconnect wait | 2 s | Delay between reconnect attempts. |
| Max reconnect attempts | 60 (or unlimited) | Before giving up. |
| Request timeout | none (you set it) | `request()` must pass a timeout. |
| Drain timeout | 30 s | Graceful drain window. |
| Flush timeout | 10 s | `flush()` wait. |

Exact defaults differ per client — check your SDK ([[NATS/08-CLIENTS-SDKS/README|Section 08]]).

---

## JetStream limits

**Account-level** (set via config/JWT, enforced by server):

| Limit | Default | Meaning |
|-------|---------|---------|
| `max_memory` (`max_memory_store`) | server-set | Total RAM for memory streams. |
| `max_file` (`max_file_store`) | server-set | Total disk for file streams. |
| `max_streams` | unlimited (-1) | Streams per account. |
| `max_consumers` | unlimited (-1) | Consumers per account. |
| `max_ack_pending` | 1000 (consumer) | Unacked messages in flight per consumer. |

**Stream-level** (per stream config):

| Limit | Meaning |
|-------|---------|
| `max_msgs` | Max messages retained (-1 = unlimited). |
| `max_bytes` | Max total bytes. |
| `max_age` | Max message age (TTL). |
| `max_msg_size` | Per-message cap within the stream. |
| `max_msgs_per_subject` | Cap per subject (KV uses this for history). |
| `replicas` | **1** (set 3/5 for HA). |
| `duplicate_window` | **2 min** | Dedup window for `Nats-Msg-Id`. |

**Consumer-level:**

| Setting | Default | Meaning |
|---------|---------|---------|
| `ack_policy` | explicit (or none/all) | How acks work. |
| `ack_wait` | **30 s** | Redelivery if no ack. |
| `max_deliver` | -1 (unlimited) | Max redelivery attempts. |
| `replay_policy` | instant | instant or original timing. |
| `max_ack_pending` | 1000 | In-flight unacked cap. |

See [[NATS/02-JETSTREAM/streams|Streams]] · [[NATS/02-JETSTREAM/consumers|Consumers]].

---

## Replica counts (JetStream RAFT)

| Replicas | Tolerates | Use |
|----------|-----------|-----|
| **R1** | no failures | dev, drop-tolerant, single node |
| **R3** | 1 server loss | **standard production** HA |
| **R5** | 2 server losses | critical, larger clusters |

Always **odd** (quorum = ⌊N/2⌋+1). R3 needs ≥3 servers; R5 needs ≥5. See [[NATS/02-JETSTREAM/jetstream-clustering|Clustering & Replication]].

---

## Storage

| | Memory | File |
|--|--------|------|
| Speed | fastest | fast (SSD) |
| Durability | lost on restart | survives restart |
| Capacity | RAM-bound | disk-bound |
| Use | ephemeral, cache | durable streams (default) |

`store_dir` default varies by packaging; always set it explicitly in production.

---

## Defaults you'll override in production
- `max_payload: 1MB` — fine for most; never raise casually (use Object Store).
- `replicas: 1` — **set to 3** for anything durable.
- `http` monitoring — **enable it** (off by default) but lock it to a private network.
- `store_dir` — **set explicitly** to a known, backed-up volume.
- `server_name` — **set explicitly** (required for clusters; helpful everywhere).
- auth — default is **no auth**; never run open in production.

---

## OneMind Relevance
- **Memorize the ports** — 4222/6222/7222/7422/8222 are the fleet's fixed plan; leaf **7422** is the one every edge node dials home on.
- **R3 is the durable default** once the hub is more than one server; until then R1 is honest (don't fake HA on a single node).
- **1 MB payload cap** shapes design: telemetry frames stay small; firmware images / models / large blobs go through the **Object Store**, chunked — never as giant messages.
- **`duplicate_window` 2 min** is the exactly-once budget — set `Nats-Msg-Id` and process within it. See [[NATS/02-JETSTREAM/dedup-and-ordering|Dedup & Ordering]].
- Default **no-auth** is the thing to never ship — every fleet server runs decentralized JWT auth from day one.

## Sources
- Configuration — https://docs.nats.io/running-a-nats-service/configuration
- JetStream model — https://docs.nats.io/nats-concepts/jetstream
- Limits — https://docs.nats.io/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes
