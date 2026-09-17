---
title: "MQTT Support"
section: "11-ADVANCED"
type: note
tags: [nats, mqtt, iot, qos, jetstream, devices, sensors]
source: https://docs.nats.io/running-a-nats-service/configuration/mqtt
created: 2026-06-14
---

# MQTT Support

NATS speaks **MQTT 3.1.1 natively** — an MQTT device connects straight to `nats-server` (no bridge, no adapter) and its messages land on the **same bus** as native NATS clients. This makes NATS a drop-in broker for IoT fleets and a migration path off standalone MQTT brokers (Mosquitto, EMQX).

See also: [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]] · [[NATS/02-JETSTREAM/jetstream-overview|JetStream]] · [[NATS/11-ADVANCED/websocket|WebSocket]]

---

## Why it exists

MQTT is the lingua franca of cheap IoT hardware (ESP32, sensors, gateways). Rather than run a *separate* MQTT broker and bridge it, NATS lets those devices publish/subscribe directly, while your services consume the same data as ordinary NATS subjects. One bus, two protocols.

---

## Enabling it

MQTT **requires JetStream** — QoS state, sessions, and retained messages are persisted in streams. You also must give the account JetStream.

```text
jetstream {}                 # required — MQTT state lives here

mqtt {
  listen: "0.0.0.0:1883"     # standard MQTT port
  # tls { ... }              # → MQTTS on 8883 in practice
  # max_ack_pending: 1024
}
```

A device then connects to `mqtt://host:1883` exactly as it would to any MQTT broker.

---

## Topic ↔ subject mapping

MQTT topics and NATS subjects use different separators; NATS translates automatically:

| MQTT | NATS | Note |
|------|------|------|
| `/` (topic level sep) | `.` (subject token sep) | `sensors/barn/temp` ↔ `sensors.barn.temp` |
| `+` (single-level wildcard) | `*` | subscribe `sensors/+/temp` |
| `#` (multi-level wildcard) | `>` | subscribe `sensors/#` |

So an MQTT publish to `sensors/barn/temp` is consumable by a native NATS subscriber on `sensors.barn.temp` (or `sensors.*.temp`, or `sensors.>`). Mind leading slashes and empty levels — an MQTT topic starting with `/` produces a leading empty token.

---

## QoS & sessions

| MQTT feature | NATS behavior |
|--------------|---------------|
| **QoS 0** (at most once) | Core-style delivery, no persistence |
| **QoS 1** (at least once) | Backed by **JetStream** — PUBACK after store; redelivery on reconnect |
| **QoS 2** (exactly once) | **Not supported** — QoS 2 downgrades; design for idempotency |
| **Retained messages** | Stored (JetStream) and delivered to new subscribers |
| **Persistent sessions** (`cleanSession=false`) | Session + subscriptions persisted; offline queueing for QoS 1 |
| **Last Will (LWT)** | Supported |

The big caveat to remember: **no QoS 2**. For exactly-once semantics, lean on JetStream dedup (`Nats-Msg-Id`) on the NATS side and make consumers idempotent. See [[NATS/02-JETSTREAM/dedup-and-ordering|Dedup & Ordering]].

---

## Architecture pattern — edge MQTT into the fabric

```
ESP32 / sensor ──MQTT 1883──▶ leaf nats-server (barn) ──leaf 7422──▶ hub
   (QoS 1)                      │ JetStream (local)                   │
                               buffers offline, syncs on reconnect
```

Run MQTT on a **leaf node** at each site so devices connect locally (low latency, survives hub outage), and the leaf's JetStream **store-and-forwards** to the hub. This is the canonical IoT edge topology. See [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]].

---

## Gotchas
- **JetStream is mandatory** — MQTT silently needs it; an account without JS can't serve MQTT QoS 1/retained/sessions.
- **No QoS 2** — plan around it.
- **Topic/subject separator** — `/` vs `.` translation can surprise you with leading-slash topics; test the mapping.
- **Max payload** still applies (default 1 MB) — same `max_payload` cap as native clients.
- **Auth** — MQTT username/password maps to NATS auth; you can also use TLS client certs.

---

## OneMind Relevance (high)
- Off-the-shelf IoT hardware (greenhouse sensors, barn devices, future drone telemetry) can join the fabric **without custom NATS firmware** — they just speak MQTT to the nearest leaf.
- QoS 1 + leaf JetStream = sensor readings survive a hub outage and sync on reconnect — exactly the local-first, intermittently-connected edge behavior the fleet is built around.
- Lets us standardize on **one bus** (NATS) while still buying MQTT-native devices. No separate Mosquitto to operate.

## Sources
- MQTT config — https://docs.nats.io/running-a-nats-service/configuration/mqtt
- MQTT in NATS (overview) — https://docs.nats.io/running-a-nats-service/configuration/mqtt/mqtt_config
- JetStream — https://docs.nats.io/nats-concepts/jetstream
