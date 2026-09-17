---
title: "Installation"
section: "07-DEPLOYMENT-OPS"
type: note
tags: [nats, installation, docker, brew, binary, systemd, launchd]
source: https://docs.nats.io/running-a-nats-service/introduction/installation
created: 2026-06-14
---

# Installation

NATS is a **single static binary** — installation is mostly "put the binary somewhere and run it." Here are all the paths.

See also: [[NATS/07-DEPLOYMENT-OPS/configuration-reference|Configuration Reference]] · [[NATS/09-CLI-TOOLING/README|CLI & Tooling]]

---

## The Two Binaries

| Binary | What it is |
|--------|------------|
| **`nats-server`** | The server itself (this is NATS) |
| **`nats`** | The CLI client/admin tool (separate repo `nats-io/natscli`) — see [[NATS/09-CLI-TOOLING/README|CLI]] |

Don't confuse them: `nats-server` runs the bus; `nats` talks to it.

---

## macOS (Homebrew) — OneMind's path

```shell
brew install nats-server      # the server
brew install nats-io/nats-tools/nats   # the CLI (or: brew install nats)
nats-server --version
```

> OneMind runs the **binary via launchd** (no Docker) per the sovereignty/locality doctrine. The server is supervised by a launchd plist; see the `nats-server` ops skill for the exact plist and start/stop/reload commands.

---

## Linux (binary)

```shell
# Download a release binary
curl -sf https://binaries.nats.dev/nats-io/nats-server/v2@latest | sh
# or via package managers (deb/rpm) from the releases page
sudo install -m755 nats-server /usr/local/bin/
```

### systemd unit
```ini
[Unit]
Description=NATS Server
After=network-online.target

[Service]
ExecStart=/usr/local/bin/nats-server -c /etc/nats/nats-server.conf
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```
`systemctl enable --now nats` then `systemctl reload nats` for SIGHUP config reloads.

---

## Docker

```shell
docker run -d --name nats -p 4222:4222 -p 8222:8222 nats:latest
# With JetStream + config:
docker run -d --name nats \
  -p 4222:4222 -p 8222:8222 \
  -v $PWD/nats.conf:/etc/nats/nats-server.conf \
  -v nats-data:/data \
  nats:latest -c /etc/nats/nats-server.conf -js
```
Official image: **`nats`** on Docker Hub (Alpine + scratch variants). For clusters, see the Docker Compose / Kubernetes docs.

---

## Go install (from source)

```shell
go install github.com/nats-io/nats-server/v2@latest
go install github.com/nats-io/natscli/nats@latest
```

---

## First Run

```shell
nats-server                      # runs on :4222, no auth, Core only
nats-server -js                  # enable JetStream (file store in a temp dir)
nats-server -js -sd /var/lib/nats   # JetStream with an explicit store dir
nats-server -c nats-server.conf  # run from a config file (production)
nats-server -DV                  # debug + verbose logging (troubleshooting)
```

Verify: `nats-server --version`, then `curl localhost:8222/varz` (if monitoring enabled) or `nats server check` with the CLI.

---

## Ports (defaults)

| Port | Purpose |
|------|---------|
| **4222** | Client connections |
| **6222** | Cluster routes |
| **7222** | Gateways (supercluster) |
| **7422** | Leaf node connections |
| **8222** | HTTP monitoring |
| **1883** | MQTT |
| **8080** | WebSocket |

Open only what you use; keep cluster/route/monitoring ports off the public internet.

## Sources
- Installation — https://docs.nats.io/running-a-nats-service/introduction/installation
- Running the server — https://docs.nats.io/running-a-nats-service/introduction/running
- Docker — https://docs.nats.io/running-a-nats-service/nats_docker
- Releases — https://github.com/nats-io/nats-server/releases
