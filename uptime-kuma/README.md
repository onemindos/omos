# uptime-kuma — OneMind Watchdog

Status-monitoring watchdog for the OneMind stack. Lives on an independent DO
Droplet (SFO3) — deliberately **not** on the DOKS hub (nyc3) so it can detect
hub outages.

## Droplet

| Field | Value |
|-------|-------|
| Name | `onemind-watchdog` |
| DO Droplet ID | 582414995 |
| IP | `146.190.143.195` |
| Region | SFO3 (independent from hub nyc3) |
| Size | s-1vcpu-512mb-10gb ($4/mo) |
| SSH key | `zeus-mac-2026` (id_ed25519) |

## Run

```sh
docker run -d \
  --name uptime-kuma \
  --restart unless-stopped \
  -p 3001:3001 \
  -v uptime-kuma:/app/data \
  louislam/uptime-kuma:1
```

## Access

Web UI: `http://146.190.143.195:3001/`  
(Port 3001 is open on the droplet — access from trusted IPs only or add firewall/CF tunnel)

## Monitors (as of 2026-07-05)

| Monitor | Type | Target | Notes |
|---------|------|--------|-------|
| TAK Server (device) | port | `152.42.155.95:8089` | External LB, CoT SSL port |
| WebTAK | http | `https://152.42.155.95:8446` | 403 = alive (auth required) |
| Grafana | http | `https://grafana.onemindos.dev/api/health` | 302 = alive (CF Access) |
| NATS client | — | internal placeholder | NATS not externally exposed; re-wire when cluster push-monitor is in place |
| NATS monitoring | — | internal placeholder | Same |
| ClickHouse | — | internal placeholder | Still Mac-local; update post-migration |
| GeoServer | — | internal placeholder | ClusterIP only; update when CF tunnel route added |
| Murmur (voice) | — | pending migration | Update after Murmur migrates to DOKS |

## Notifications

**TODO: add Telegram notification** (bot token + chat ID needed from Zeus).
In kuma UI: Settings → Notifications → Add → Telegram → enter bot token + chat ID → test.
Then attach to all monitors.

## Data migration

Migrated from Mac (`/Users/zeus/omos-core/services/omos-observ`).
SQL dump preserved at `omos-infra/uptime-kuma/migration-data/kuma-20260705.sql` (gitignored — large).

## Rebuild from scratch

```sh
doctl compute droplet create onemind-watchdog \
  --region sfo3 --size s-1vcpu-512mb-10gb \
  --image ubuntu-24-04-x64 --ssh-keys <key-id> --wait
ssh root@<ip> 'curl -fsSL https://get.docker.com | sh'
# then run the docker run above and import the sql dump
```
