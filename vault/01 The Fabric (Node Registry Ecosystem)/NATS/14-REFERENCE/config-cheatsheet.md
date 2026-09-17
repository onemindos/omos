---
title: "Config Cheatsheet"
section: "14-REFERENCE"
type: note
tags: [nats, config, cheatsheet, server, conf, reference, jetstream, cluster, tls]
source: https://docs.nats.io/running-a-nats-service/configuration
created: 2026-06-14
---

# Config Cheatsheet

Every major `nats-server.conf` block at a glance. NATS config is a superset of JSON (HOCON-like: unquoted keys, comments, `include`, env vars). Copy a block, fill the blanks.

See also: [[NATS/07-DEPLOYMENT-OPS/configuration-reference|Deployment: Configuration]] · [[NATS/14-REFERENCE/limits-and-defaults|Limits & Defaults]] · [[NATS/14-REFERENCE/cli-cheatsheet|CLI Cheatsheet]]

---

## Top-level / connectivity

```text
server_name: "hub-1"            # unique per server (required in clusters)
host: "0.0.0.0"
port: 4222                      # client port
http: "0.0.0.0:8222"            # monitoring endpoint (or https:)
# https_port, no_auth_user, max_connections, max_payload, max_control_line ...

max_payload: 1MB                # single-message cap (default 1MB)
max_connections: 64K
write_deadline: "10s"

# Logging
debug: false
trace: false
logtime: true
log_file: "/var/log/nats/nats.log"

# Operational
pid_file: "/var/run/nats.pid"
include "./auth.conf"           # split config across files
```

---

## JetStream

```text
jetstream {
  store_dir: "/var/lib/nats/jetstream"
  max_memory_store: 1GB         # RAM cap for memory streams
  max_file_store: 100GB         # disk cap for file streams
  domain: "hub"                 # name this JS domain (for leaf bridging)
}
# minimal form: `jetstream {}` or `jetstream: enabled`
```
See [[NATS/02-JETSTREAM/jetstream-overview|JetStream]] · [[NATS/06-TOPOLOGY-SCALING/domains|Domains]].

---

## Cluster (server ↔ server, full mesh)

```text
cluster {
  name: "omos-hub"              # all members share this name
  listen: "0.0.0.0:6222"        # route port
  routes: [
    nats-route://hub-2:6222
    nats-route://hub-3:6222
  ]
  # authorization { user: ..., password: ... }   # route auth
  # pool_size, no_advertise ...
}
```
3 or 5 servers for JetStream RAFT quorum. See [[NATS/06-TOPOLOGY-SCALING/clustering|Clustering]].

---

## Gateways (cluster ↔ cluster = supercluster)

```text
gateway {
  name: "us-east"              # = this cluster's name
  listen: "0.0.0.0:7222"
  gateways: [
    { name: "us-west", urls: ["nats-gateway://uswest-1:7222"] }
    { name: "eu",      urls: ["nats-gateway://eu-1:7222"] }
  ]
}
```
See [[NATS/06-TOPOLOGY-SCALING/superclusters-gateways|Superclusters]].

---

## Leaf nodes (the edge)

```text
# On the HUB — accept leaf connections:
leafnodes {
  listen: "0.0.0.0:7422"
}

# On the LEAF — connect up to the hub:
leafnodes {
  remotes: [
    {
      urls: ["nats-leaf://hub.example:7422"]
      credentials: "/etc/nats/leaf.creds"
      # account: "APP"
    }
  ]
}
```
See [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]].

---

## TLS

```text
tls {
  cert_file: "/etc/nats/server-cert.pem"
  key_file:  "/etc/nats/server-key.pem"
  ca_file:   "/etc/nats/ca.pem"
  verify:    true              # require + verify client certs (mTLS)
  # verify_and_map: true       # map cert identity → NATS user
  timeout:   2
}
```
Separate `tls{}` blocks can live inside `cluster`, `gateway`, `leafnodes`, `websocket`, `mqtt`. See [[NATS/05-SECURITY/tls-encryption|TLS]].

---

## Auth — simple

```text
# Single user
authorization {
  user: "app"
  password: "s3cret"           # or token: "..."
  timeout: 2
}

# Multiple users + permissions
authorization {
  users: [
    { user: "pub", password: "...", permissions: {
        publish: ["orders.>"], subscribe: ["_INBOX.>"] } }
    { user: "sub", password: "...", permissions: {
        subscribe: ["orders.>"] } }
  ]
}
```

## Auth — accounts (multi-tenant)

```text
accounts {
  APP: {
    users: [ { user: "app", password: "..." } ]
    jetstream: enabled
    exports: [ { service: "svc.echo" } ]
  }
  MONITOR: {
    users: [ { user: "mon", password: "..." } ]
    imports: [ { service: { account: "APP", subject: "svc.echo" } } ]
  }
}
system_account: "SYS"
```
See [[NATS/05-SECURITY/accounts-multitenancy|Accounts]].

## Auth — decentralized (operator/JWT)

```text
operator: "/etc/nats/operator.jwt"
system_account: "ADxxxx..."     # the $SYS account public key
resolver: {
  type: full                    # built-in NATS resolver
  dir: "/etc/nats/jwt"
  allow_delete: false
  interval: "2m"
}
# resolver_preload { ACCT_PUBKEY: "<jwt>" }   # for a few static accounts
```
See [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS]].

---

## MQTT (requires JetStream)

```text
jetstream {}
mqtt {
  listen: "0.0.0.0:1883"
  # tls { ... }   # → 8883
}
```
See [[NATS/11-ADVANCED/mqtt|MQTT]].

## WebSocket

```text
websocket {
  listen: "0.0.0.0:8080"
  tls { cert_file: "...", key_file: "..." }   # wss
  # no_tls: true            # localhost/dev only
  compression: true
  # allowed_origins: ["https://command.onemind.local"]
}
```
See [[NATS/11-ADVANCED/websocket|WebSocket]].

---

## Reload & validate

```bash
nats-server -c nats.conf -t        # -t = validate config & exit
nats-server --signal reload        # hot-reload most settings (or: kill -SIGHUP <pid>)
```

Many settings hot-reload via `SIGHUP` (auth, logging, TLS, limits); some (ports, clustering topology) need a restart.

## OneMind Relevance
- Our hub config is essentially: top-level + `jetstream{}` + `leafnodes{listen}` (to accept edge nodes) + TLS + decentralized auth (`operator`/`resolver`). Edge nodes invert the leaf block (`remotes`).
- Keep auth in an `include`d file (`auth.conf`) out of the main config and out of git — secrets reference 1Password, never inline.
- `nats-server -t` in CI before any deploy — config typos are an outage.

## Sources
- Configuration — https://docs.nats.io/running-a-nats-service/configuration
- Config reload — https://docs.nats.io/running-a-nats-service/configuration/cfgfile
