# mac-leaf — OneMind Mac Leaf Node

Zeus's MacBook Pro runs a full NATS server that bridges to the DOKS cluster via Tailscale.

## Architecture

```
Mac (zeus-mac)
├── nats-server (port 4222) — full server with JetStream
│   ├── JetStream store → ~/nats-data/
│   └── Leaf connection → omos-nats.tail717ba0.ts.net:7422 (DOKS)
├── NEX node (zeus-mac) — runs local workloads
│   ├── om-streaming-bridge
│   ├── om-callout
│   ├── om-http-gateway
│   ├── om-bridge-mavlink
│   ├── om-bridge-constellation
│   └── om-nex-ui
└── Services
    ├── om_presence_v2.py
    ├── heynats (port 5000)
    └── nats-dashboard (port 5173)
```

## Disaster Recovery — New Mac Setup

### Prerequisites
```bash
brew install nats-server nats-io/nats-tools/nats nex
```

### Steps
1. Get credentials from **1Password → "omos-mac-leaf-node"**
2. Copy `nats-server.conf` to `~/nats-server.conf`
3. Replace `LEAF_USER` and `LEAF_PASSWORD` with values from 1Password
4. Create `~/nats-data/` directory
5. Install launchd agents:
```bash
cp launchd/*.plist ~/Library/LaunchAgents/
launchctl load ~/Library/LaunchAgents/com.nats.server.plist
launchctl load ~/Library/LaunchAgents/com.onemind.natsdata-guard.plist
launchctl load ~/Library/LaunchAgents/omos.heynats.plist
launchctl load ~/Library/LaunchAgents/omos.nats-dashboard.plist
```
6. Verify leaf connection:
```bash
nats server info --server nats://localhost:4222
nats server ls
```

## Files
| File | Purpose |
|---|---|
| `nats-server.conf` | Sanitized config (credentials replaced with placeholders) |
| `launchd/com.nats.server.plist` | NATS server autostart |
| `launchd/com.onemind.natsdata-guard.plist` | .DS_Store cleaner for nats-data/ |
| `launchd/omos.heynats.plist` | HeyNATS web UI autostart |
| `launchd/omos.nats-dashboard.plist` | NATS dashboard autostart |

## Credentials Location
**1Password → 00-24 SO (Sovereign Operations) → "omos-mac-leaf-node"**
Contains: LEAF_USER, LEAF_PASSWORD, full nats-server.conf with real credentials
