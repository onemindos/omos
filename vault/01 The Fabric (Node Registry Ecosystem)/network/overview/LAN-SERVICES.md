---
title: "LAN Services — omos.lan DNS Stack"
type: reference
status: active
created: 2026-05-27
updated: 2026-05-27
owner: zeus
tags: [type/reference, domain/SO, domain/02-infrastructure, networking, dns, caddy, unbound, netbird]
ai_summary: "Complete reference for the omos.lan internal DNS stack — Unbound + Caddy + NetBird nameserver group. All services accessible at http://<service>.omos.lan:8888 from Mac, Dell, and iPhone."
---

# LAN Services — omos.lan DNS Stack

> All OneMind OS services on the MBP are accessible at `http://<service>.omos.lan:8888` from any device on the NetBird mesh (Mac, Dell, iPhone).

---

## Why `.lan` Not `.dev`

`.dev` is owned by Google and is HSTS-preloaded — all browsers silently force HTTPS for any `.dev` URL. Since Caddy runs plain HTTP on port 8888 (no TLS), `.dev` breaks in every browser even when DNS resolves correctly. `.lan` has no HSTS preloading and works cleanly over HTTP.

---

## Architecture

```
Browser on any device
    │  http://hermes.omos.lan:8888
    ▼
NetBird DNS proxy (100.65.255.254:53)   ← pushed to all peers via nameserver group
    │  omos.lan → 192.168.1.182
    ▼
Unbound (192.168.1.182:5335)            ← resolves *.omos.lan → 192.168.1.182
    │
    ▼
Caddy (:8888, all interfaces)           ← routes by Host: header
    │
    ▼
Local service (e.g. localhost:8787)
```

**Key design choices:**
- DNS A records point to **LAN IP** (192.168.1.182), not the NetBird IP (100.65.90.213) — macOS userspace WireGuard cannot self-route traffic through its own utun interface
- Caddy binds on `*:8888` (not a specific interface) — Host-header matching handles all virtual hosts
- Caddy uses port 8888 (not 80/443) — avoids requiring root privileges

---

## Service Directory

| URL | Service | Local Port |
|-----|---------|-----------|
| `http://dashboard.omos.lan:8888` | Homepage Dashboard | 3010 |
| `http://hermes.omos.lan:8888` | Hermes WebUI | 8787 |
| `http://hermes-dash.omos.lan:8888` | Hermes Dashboard | 9119 |
| `http://paperclip.omos.lan:8888` | Paperclip AI | 3101 |
| `http://nats.omos.lan:8888` | NATS Monitor | 8222 |
| `http://nats-dash.omos.lan:8888` | NATS Visual Dashboard | 4321 |
| `http://uptime.omos.lan:8888` | Uptime Kuma | 3001 |
| `http://ocis.omos.lan:8888` | ownCloud OCIS | 9200 |
| `http://netbird.omos.lan:8888` | NetBird Admin UI | 8081 |
| `http://agent.omos.lan:8888` | Synadia Agent Web UI | 3300 |
| `http://fleet.omos.lan:8888` | OneMind Fleet Dashboard | 3002 |
| `http://legacy.omos.lan:8888` | Legacy Next.js UI | 3000 |
| `http://netdata.omos.lan:8888` | Netdata Metrics | 19999 |

---

## Component Config Files

| Component | Config Path |
|-----------|------------|
| Unbound | `/opt/homebrew/etc/unbound/unbound.conf` |
| Caddy | `/opt/homebrew/etc/caddy/Caddyfile` |
| NetBird client | `/Library/LaunchDaemons/netbird.plist` |
| Homepage services | `/Users/onemind/homepage/config/services.yaml` |

---

## Unbound

**Binary:** `/opt/homebrew/bin/unbound` (Homebrew)  
**Port:** 5335 (not 53 — avoids conflict with system resolver)  
**Listens on:** 127.0.0.1, 192.168.1.182 (LAN), 100.65.90.213 (NetBird)  
**Access control:** 127.0.0.0/8, 192.168.1.0/24, 100.64.0.0/10

Unbound runs as the current user (`username: ""`, `chroot: ""`). Started manually — not yet a brew service.

**Start:**
```bash
unbound -c /opt/homebrew/etc/unbound/unbound.conf
```

**Check running:**
```bash
lsof -i :5335 -P -n | grep LISTEN
```

**Test resolution:**
```bash
dig hermes.omos.lan @127.0.0.1 -p 5335 +short
# → 192.168.1.182
```

---

## Caddy

**Binary:** `/opt/homebrew/Cellar/caddy/2.11.3/bin/caddy` (Homebrew)  
**Port:** 8888 (all interfaces)  
**Mode:** HTTP only (`auto_https off`)

Caddy is **allowlisted in the macOS Application Firewall** — required for other devices to reach it on port 8888.

**Start:**
```bash
caddy start --config /opt/homebrew/etc/caddy/Caddyfile
```

**Stop:**
```bash
caddy stop
```

**Reload config without restart** (only works when admin endpoint is enabled — currently off):
```bash
caddy reload --config /opt/homebrew/etc/caddy/Caddyfile
# If that fails: pkill -x caddy && caddy start --config /opt/homebrew/etc/caddy/Caddyfile
```

**Note on Host header rewriting:** Some services (e.g. Hermes Dashboard) reject any `Host:` header other than `127.0.0.1`. For those, Caddy must rewrite the header upstream:
```
reverse_proxy localhost:9119 {
    header_up Host "127.0.0.1:9119"
}
```

**Re-add firewall rule after Caddy update (if binary path changes):**
```bash
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add /opt/homebrew/Cellar/caddy/2.11.3/bin/caddy
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --unblockapp /opt/homebrew/Cellar/caddy/2.11.3/bin/caddy
```

---

## NetBird Nameserver Group

NetBird pushes the `omos.lan` DNS resolver to all peers automatically.

**Nameserver group:** `omos-lan` (ID: `d8bgp796891c73f9m0cg`)  
**Nameserver:** `192.168.1.182:5335` (UDP)  
**Domain:** `omos.lan`  
**Groups:** All  
**Stored in:** NetBird SQLite DB at `netbird-server:/var/lib/netbird/store.db`, table `name_server_groups`

When a peer connects to the NetBird mesh, it receives this config and adds a supplemental resolver entry via `scutil`:
```
domain   : omos.lan
nameserver[0] : 100.65.255.254   ← NetBird DNS proxy
port     : 53
```

NetBird's DNS proxy (100.65.255.254) forwards `omos.lan` queries to Unbound at 192.168.1.182:5335.

**Verify on any device:**
```bash
netbird status | grep Nameserver
# → Nameservers: 1/1 Available

scutil --dns | grep -A4 "omos.lan"
```

---

## Peer Devices

| Device | NetBird IP | LAN IP | DNS Active |
|--------|-----------|--------|-----------|
| omos-mbp | 100.65.90.213 | 192.168.1.182 | ✅ |
| omos-dell | 100.65.18.96 | 192.168.1.181 | ✅ (reconnect after changes) |
| omos-iphone | 100.65.56.131 | 192.168.1.175 | ✅ (toggle NetBird in app) |

**To push DNS config updates to peers:**
1. Update `store.db`, restart `netbird-server` container
2. Each peer: `netbird down && netbird up` (or toggle app on iPhone)

---

## Troubleshooting

**Links don't load in browser (Mac):**
1. Verify Unbound is running: `lsof -i :5335 -P -n | grep LISTEN`
2. Verify Caddy is running: `pgrep caddy`
3. Test DNS via proxy: `dig hermes.omos.lan @100.65.255.254 +short`
4. Test curl (bypasses browser): `curl -s -o /dev/null -w "%{http_code}" http://hermes.omos.lan:8888`
5. If curl works but browser doesn't — check browser isn't forcing HTTPS

**Links don't load from Dell/iPhone:**
1. Verify firewall allows Caddy: `/usr/libexec/ApplicationFirewall/socketfilterfw --getappblocked /opt/homebrew/Cellar/caddy/2.11.3/bin/caddy`
2. Verify peers received DNS: `netbird status | grep Nameserver`
3. Test from Dell: `curl http://hermes.omos.lan:8888`

**After MBP reboot:**
- Unbound and Caddy do NOT auto-start yet (not configured as launch agents)
- Run: `unbound -c /opt/homebrew/etc/unbound/unbound.conf & caddy start --config /opt/homebrew/etc/caddy/Caddyfile`

---

## Related

- [[netbird/netbird-iot-fleet-foundation]] — NetBird mesh VPN overview
- [[dns/DOMAIN-REGISTRY]] — Public domain portfolio
- [[overview/NETWORKING]] — OpenClaw / VPS / Tailscale networking
- NetBird management UI: `http://netbird.omos.lan:8888`
