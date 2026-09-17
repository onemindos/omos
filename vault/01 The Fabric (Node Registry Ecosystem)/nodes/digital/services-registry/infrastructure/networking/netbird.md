---
title: "NetBird"
type: reference
status: active
tags: [platform, infrastructure, networking, mesh, vpn, dns, proxy, zero-trust]
---
# NetBird

> Unified networking platform: WireGuard mesh + DNS + reverse proxy + zero-trust ACLs. Self-hosted. Replaces Caddy, Unbound, Tailscale.

## Role in OneMind

NetBird is the **single networking layer** for the entire fleet:
- Mesh connectivity between all nodes (NAT traversal)
- DNS resolution (`*.omos.lan`) for all services
- Reverse proxy (HTTP/TCP/UDP routing by hostname)
- Zero-trust access control (group-based policies)
- IoT/edge node enrollment (setup keys)
- TLS certificate management

## Self-Hosted Stack (Docker on MBP)
- Management server
- Signal server
- TURN/Relay (coturn)
- Dashboard (nginx)

## Features Used

| Feature | Purpose |
|---------|---------|
| Custom DNS Zones | `*.omos.lan` resolution across all peers |
| Wildcard DNS Labels | Service-level DNS |
| DNS Round-Robin | Load distribution |
| Built-in Reverse Proxy (HTTP L7) | Route `hermes.omos.lan` → localhost:8787 |
| TCP/UDP L4 Proxy | Non-HTTP services |
| Auto TLS Certificates | HTTPS for all proxied services |
| Setup Keys | Headless IoT/edge enrollment |
| Group-based ACLs | Zero-trust network segmentation |
| Posture Checks | Device security verification |
| SSO/OIDC | Identity-aware access |
| Networks (routing peers) | Bridge to non-NetBird devices on LAN |

## Service Routes (previously Caddy)

| Domain | Backend | Service |
|--------|---------|---------|
| `hermes.omos.lan` | :8787 | Hermes WebUI |
| `hermes-dash.omos.lan` | :9119 | Hermes Dashboard |
| `dashboard.omos.lan` | :3010 | Homepage |
| `nats.omos.lan` | :8222 | NATS Monitoring |
| `nats-dash.omos.lan` | :4321 | NATS Dashboard |
| `uptime.omos.lan` | :3001 | Uptime Kuma |
| `netbird.omos.lan` | :8081 | NetBird Dashboard |
| `agent.omos.lan` | :3300 | OneMind Dashboard |
| `fleet.omos.lan` | :3002 | Fleet View |
| `legacy.omos.lan` | :3000 | Legacy UI |
| `netdata.omos.lan` | :19999 | Netdata |

## Fleet Enrollment

```bash
# New node joins the mesh:
netbird up --setup-key <KEY>

# That's it — node gets:
# - Encrypted mesh connectivity
# - DNS resolution for all *.omos.lan
# - Access per ACL policies
```

## Replaces
- ~~Caddy~~ (reverse proxy)
- ~~Unbound~~ (DNS resolver)
- ~~Tailscale~~ (mesh VPN)

## Integration
- **NATS**: Messaging layer rides on NetBird mesh
- **Cloudflare Tunnels**: External/public access (NetBird = internal)
- **Docker Compose**: Containers on base node, NetBird routes to them
