# Network — OneMind OS

## Mesh VPN: NetBird (Self-Hosted)

Replaced Tailscale on 2026-05-26. Self-hosted management server on omos-mbp via Docker, accessible at `net.onemindos.dev` through Cloudflare Tunnel.

### Devices

| Device | NetBird IP | DNS Name | LAN IP | OS | Persistent |
|--------|-----------|----------|--------|-----|------------|
| omos-mbp | 100.65.90.213 | omos-mbp-90-213.mesh | 192.168.1.182 | macOS 26.5 | LaunchDaemon |
| omos-dell | 100.65.18.96 | omos-dell.mesh | 192.168.1.181 | Ubuntu 26.04 | systemd |
| omos-ha | 100.65.136.250 | omos-ha.mesh | 192.168.1.180 | Alpine (HA OS) | s6-overlay |
| omos-iphone | 100.65.56.131 | omos-iphone.mesh | dynamic | iOS 26.5 | iOS app |

### DNS

Domain: `.mesh`
Example: `ssh zeus@omos-dell.mesh`

### Management

- **Dashboard:** https://net.onemindos.dev
- **API:** http://localhost:33073/api (local) or https://net.onemindos.dev/api (remote)
- **Auth:** Embedded IdP (Dex) — `me@zeusdelacruz.com`
- **Config dir:** `/Users/onemind/netbird/`

### Architecture

```
Internet (Cloudflare Tunnel) ─── net.onemindos.dev
                                       │
                                 nginx (8081)
                                   │   │   │
                          gRPC(h2) │   │   │ HTTP
                                   v   v   v
                              netbird-server (Docker)
                              ├── Management + Signal
                              ├── Relay (WebSocket)
                              ├── STUN (3478/udp)
                              └── Embedded IdP (Dex)
                                       │
                    ┌──────────────────┼──────────────────┐
                    │                  │                   │
              localhost:33073    0.0.0.0:33080       3478/udp
              (Mac client)      (LAN fallback)      (STUN)
```

### SSH Access

| Device | Command |
|--------|---------|
| Dell | `ssh zeus@omos-dell.mesh` |
| HA | `ssh hassio@omos-ha.mesh` |

### Key Files (Mac)

| File | Purpose |
|------|---------|
| `/Users/onemind/netbird/docker-compose.yml` | Docker stack |
| `/Users/onemind/netbird/nginx.conf` | Reverse proxy (h2c + WS + gRPC) |
| `/Users/onemind/netbird/config.yaml` | Server config |
| `/Users/onemind/netbird/dashboard.env` | Dashboard env |
| `/Library/LaunchDaemons/netbird.plist` | Mac client service |
| `~/.cloudflared/config.yml` | Cloudflare Tunnel |

### Critical Settings

- `http2 on;` in nginx.conf — required for gRPC Signal Exchange
- `protocol: http2` in cloudflared config — h2mux kills WebSocket streams
- Mac client uses `--foreground-mode` — daemon gRPC socket broken in 0.71.4 on ARM
- HA uses `--disable-dns` — prevents resolv.conf hijack in container

---

## LAN

| Device | LAN IP | MAC |
|--------|--------|-----|
| Router/Gateway | 192.168.1.1 | — |
| omos-ha (Home Assistant) | 192.168.1.180 | — |
| omos-dell (Linux server) | 192.168.1.181 | — |
| omos-mbp (Mac) | 192.168.1.182 | — |

---

## Cloudflare Tunnel

Tunnel ID: `1e7cfeee-5157-4e27-84bd-8b54dacbf231`
Config: `~/.cloudflared/config.yml`
Service: `brew services` (LaunchAgent, auto-starts on login)

| Hostname | Target | Purpose |
|----------|--------|---------|
| net.onemindos.dev | localhost:8081 | NetBird (management + dashboard) |
| files.onemindos.dev | localhost:9200 | oCIS file storage |
| status.onemindos.dev | localhost:3001 | Uptime Kuma |

---

## Tailscale (Deprecated)

Removed 2026-05-27. Previously used for mesh VPN. Devices that were on Tailscale but not yet on NetBird:

| Device | Tailscale IP | Status | Notes |
|--------|-------------|--------|-------|
| omos-pad | 100.67.182.83 | Offline 29d | Android tablet |
| omos-tab | 100.81.60.35 | Offline 28d | Android tablet |
| omos-tv | 100.117.247.62 | Offline 12d | Android TV |

These can be added to NetBird when powered on using a setup key + the NetBird Android app.
