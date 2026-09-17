---
title: "Networking Architecture"
type: architecture
status: active
created: 2026-05-29
modified: 2026-05-29
owner: zeus
tags: [architecture, networking, netbird, dns, proxy]
---
# Networking Architecture

## Domain: `*.omos`

All internal services are accessible at `<service>.omos:8888` from any device on the mesh.

## Stack

| Layer | Tool | Purpose |
|-------|------|---------|
| Mesh VPN | **NetBird** (self-hosted) | Connects all nodes: MBP, phone, office, IoT |
| DNS | **Unbound** | Wildcard `*.omos` → 192.168.1.182 (set-and-forget) |
| Reverse Proxy | **Caddy** | Routes `<service>.omos:8888` → `localhost:<port>` |
| External Access | **Cloudflare Tunnels** | Public internet → selected services |
| Messaging | **NATS** | Inter-node commands + telemetry |

## How It Works

```
Any device (phone/laptop/remote node)
  → NetBird mesh (encrypted WireGuard)
    → DNS query: hermes.omos
      → NetBird nameserver group → Unbound (192.168.1.182:5335)
        → Wildcard: *.omos → 192.168.1.182
  → HTTP request: hermes.omos:8888
    → Caddy → localhost:8787 (Hermes WebUI)
```

## Active Services

| URL | Port | Service |
|-----|------|---------|
| hermes.omos:8888 | 8787 | Hermes WebUI |
| hermes-dash.omos:8888 | 9119 | Hermes Dashboard |
| dashboard.omos:8888 | 3010 | Homepage |
| nats.omos:8888 | 8222 | NATS Monitoring |
| nats-dash.omos:8888 | 4321 | NATS Dashboard |
| uptime.omos:8888 | 3001 | Uptime Kuma |
| netbird.omos:8888 | 8081 | NetBird Dashboard |
| agent.omos:8888 | 3300 | OneMind Dashboard |
| fleet.omos:8888 | 3002 | Fleet View |
| legacy.omos:8888 | 3000 | Legacy UI |
| netdata.omos:8888 | 19999 | Netdata |
| langfuse.omos:8888 | 3004 | Langfuse |

## Adding a New Service

1. Start the service on any port
2. Add a route to `/opt/homebrew/etc/caddy/Caddyfile`:
   ```
   @myservice host myservice.omos
   handle @myservice {
       reverse_proxy localhost:<port>
   }
   ```
3. Restart Caddy: `brew services restart caddy`
4. Done — DNS is wildcard, no changes needed

## Config Files

- Unbound: `/opt/homebrew/etc/unbound/unbound.conf` (wildcard — don't touch)
- Caddy: `/opt/homebrew/etc/caddy/Caddyfile`
- NetBird: self-hosted at net.onemindos.dev (Docker compose in ~/netbird/)

## Connected Peers

| Peer | IP | Role |
|------|-----|------|
| omos-mbp | 100.65.90.213 | Base (all services) |
| omos-dell | 100.65.18.96 | Office node |
| omos-ha | 100.65.136.250 | Home Assistant |
| omos-iphone | 100.65.56.131 | Mobile |

## Future

- Replace Caddy with Traefik for Docker auto-discovery (zero-config new services)
- Or deploy NetBird reverse proxy when cross-location TLS is needed
- Port 80 blocked by macOS system service — use pfctl redirect or Traefik if portless URLs needed
