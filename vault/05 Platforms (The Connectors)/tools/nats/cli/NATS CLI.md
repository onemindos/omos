---
title: "CLI — nats (NATS CLI)"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/nats, topic/messaging]
ai_summary: "NATS CLI deep reference — pub/sub messaging, JetStream streams/consumers, server monitoring. Connects to onemind-nats container on port 4222."
---

# nats — NATS CLI

> Message bus between OpenClaw agents and services on `onemind-main`.
> Container: `onemind-nats` | Port: `4222` (client), `8222` (monitoring)
> Install: `brew install nats-io/nats-tools/nats`

---

## Connect

```bash
# Connect to VPS NATS (via Tailscale — preferred)
nats --server nats://100.77.52.93:4222 server ping

# Connect via public IP
nats --server nats://104.248.4.212:4222 server ping

# If auth is set
nats --server nats://user:password@100.77.52.93:4222 server ping

# Save a context for quick access
nats context add onemind-main --server nats://100.77.52.93:4222
nats context select onemind-main
```

---

## Pub/Sub

```bash
# Subscribe to a subject (blocks, prints incoming)
nats sub "codex.events.>"

# Publish a message
nats pub "codex.events.test" "hello from cli"

# Wildcard subscriptions
nats sub "agents.>"          # all agent messages
nats sub "agents.legacy.*"   # legacy agent specific subjects
```

---

## JetStream (Streams & Consumers)

```bash
# List streams
nats stream list

# Create a stream
nats stream create CODEX_EVENTS \
  --subjects "codex.events.>" \
  --storage file \
  --retention limits \
  --max-age 24h

# Stream info
nats stream info CODEX_EVENTS

# Purge stream
nats stream purge CODEX_EVENTS

# Delete stream
nats stream delete CODEX_EVENTS

# Create push consumer
nats consumer create CODEX_EVENTS event-processor \
  --pull \
  --deliver all \
  --ack explicit

# List consumers for a stream
nats consumer list CODEX_EVENTS

# Consume messages
nats consumer next CODEX_EVENTS event-processor
```

---

## Server Monitoring

```bash
# Ping server
nats server ping

# Server info
nats server info

# Server stats
nats server report connections

# View via HTTP monitoring endpoint
curl http://100.77.52.93:8222/varz        # server stats
curl http://100.77.52.93:8222/connz       # connections
curl http://100.77.52.93:8222/subsz       # subscriptions
curl http://100.77.52.93:8222/jsz         # JetStream info
```

---

## Key-Value Store

```bash
# Create KV bucket
nats kv add onemind-config --ttl 24h

# Put value
nats kv put onemind-config key "value"

# Get value
nats kv get onemind-config key

# List keys
nats kv ls onemind-config

# Watch for changes
nats kv watch onemind-config
```

---

## Docker (VPS)

```bash
# Check container
docker ps | grep onemind-nats

# NATS logs
docker logs onemind-nats -f

# NATS config
docker exec onemind-nats cat /etc/nats/nats.conf
```

---

## OneMind Agent Subjects

| Subject | Used By | Purpose |
|---------|---------|---------|
| `agents.legacy.*` | Legacy | Orchestration events |
| `codex.events.>` | All agents | Vault change events |
| `pulse.heartbeat` | legacy-pulse | Service heartbeat |


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
