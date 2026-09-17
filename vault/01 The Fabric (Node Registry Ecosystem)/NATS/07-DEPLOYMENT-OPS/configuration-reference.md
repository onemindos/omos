---
title: "Configuration Reference"
section: "07-DEPLOYMENT-OPS"
type: note
tags: [nats, configuration, config-file, reference, limits]
source: https://docs.nats.io/running-a-nats-service/configuration
created: 2026-06-14
---

# Configuration Reference

The `nats-server.conf` surface. NATS config is a superset of JSON (HOCON-like): unquoted keys, comments, `include`, env-var interpolation, durations (`2m`, `1h`), and sizes (`1GB`).

See also: [[NATS/07-DEPLOYMENT-OPS/installation|Installation]] · [[NATS/05-SECURITY/README|Security]] · [[NATS/06-TOPOLOGY-SCALING/README|Topology]]

---

## Top-Level Keys

```text
# ---- Identity / networking ----
server_name: "omos-mbp"
host: "0.0.0.0"
port: 4222
http_port: 8222                 # monitoring (or https_port for TLS)

# ---- Limits ----
max_connections: 64K
max_payload: 1MB                # max single message (default ~1,048,576 bytes)
max_pending: 64MB               # per-connection outbound buffer before slow-consumer
max_control_line: 4KB
max_subscriptions: 0            # 0 = unlimited
ping_interval: "2m"
ping_max: 2
write_deadline: "10s"

# ---- Logging ----
debug: false
trace: false
logtime: true
log_file: "/var/log/nats/nats.log"
log_size_limit: 100MB

# ---- JetStream ----
jetstream {
  store_dir: "/var/lib/nats/jetstream"
  max_memory_store: 1GB
  max_file_store: 50GB
  domain: "hub"                 # optional; see Domains note
}

# ---- Security (see Section 05) ----
authorization { ... }           # or accounts { } + system_account
tls { ... }

# ---- Topology (see Section 06) ----
cluster   { ... }               # routes (6222)
gateway   { ... }               # supercluster (7222)
leafnodes { ... }               # edge (7422)

# ---- Protocol gateways (see Section 11) ----
mqtt      { ... }               # 1883
websocket { ... }               # 8080

# ---- Includes ----
include "./accounts.conf"
```

---

## Key Limits Explained

| Key | Default | Meaning |
|-----|---------|---------|
| **`max_payload`** | ~1 MB (1,048,576) | Largest single message (payload + headers). Raising costs memory per message; use Object Store for blobs |
| **`max_pending`** | 64 MB | Per-connection outbound queue; exceed → **slow consumer**, messages dropped for that sub |
| **`max_connections`** | 64K | Cap on simultaneous client connections |
| **`max_control_line`** | 4 KB | Max length of a protocol control line (subject + metadata) |
| **`write_deadline`** | 10 s | How long the server waits writing to a slow client before flagging it |
| **`max_subscriptions`** | 0 (∞) | Per-connection subscription cap |

JetStream-specific:
| Key | Meaning |
|-----|---------|
| **`store_dir`** | Where file-backed streams live (put on fast disk) |
| **`max_file_store`** | Cap on total file-store bytes |
| **`max_memory_store`** | Cap on total memory-store bytes |
| **`domain`** | JetStream domain name (see [[NATS/06-TOPOLOGY-SCALING/domains|Domains]]) |

---

## Env Vars & Includes

```text
# Interpolate environment variables
authorization { token: $NATS_TOKEN }

# Split config across files
include "./security.conf"
include "./jetstream.conf"
```

Great for keeping secrets out of the main file and composing large configs.

---

## Reloading Config

```shell
nats-server --signal reload          # SIGHUP: reload in place
# or
kill -HUP $(pgrep -f nats-server)
```

**What reloads via SIGHUP** (no restart): most limits, logging, TLS cert files, config-based users/permissions, authorization.
**What needs a FULL RESTART** (verified pitfalls): `system_account`, account topology changes, JetStream `domain`, and some listener/port changes. Don't expect SIGHUP to swap the system account.

---

## Validating Config

```shell
nats-server -c nats-server.conf -t   # -t = parse/validate config and exit
```
Always validate before reloading a production server.

---

## Minimal Production Skeleton

```text
server_name: "omos-mbp"
port: 4222
http_port: 8222
jetstream { store_dir: "/var/lib/nats/js", max_file_store: 50GB }
accounts {
  SYS:  { users: [{ user: sys,  password: $SYS_PW  }] }
  CORE: { jetstream: enabled, users: [{ user: core, password: $CORE_PW }] }
}
system_account: SYS
```

## Sources
- Configuration — https://docs.nats.io/running-a-nats-service/configuration
- Config reload / signals — https://docs.nats.io/running-a-nats-service/nats_admin/signals
- Limits — https://docs.nats.io/running-a-nats-service/configuration#limits
- JetStream config — https://docs.nats.io/running-a-nats-service/configuration/resource_management
