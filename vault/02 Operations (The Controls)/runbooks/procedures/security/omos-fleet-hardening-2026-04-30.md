# Fleet Hardening — 2026-04-30

## Summary
Session: "Supercharging OpenClaw + Fleet Unification"

## Completed

### OpenClaw v2026.4.27
- Updated on: omos-core, omos-sentinel
- Service: `openclaw-gateway.service` (systemd user)

### Tailscale SSH Mesh
- All fleet nodes: `tailscale set --ssh --accept-risk=lose-ssh`
- ACL: `action: accept` (check mode off)
- Mac (zeus-mac) added: Remote Login enabled, fleet key in `~/.ssh/authorized_keys`

### omos-sentinel — Fleet Medic Agent
- **Path**: `/opt/sentinel/sentinel.py`
- **Service**: `sentinel.service` (systemd user, ubuntu)
- **Port**: 9119
- **Monitors**: omos-core, omos-titan, zeus-mac (3 active nodes)
- **Telegram**: Two-way — alerts to topic 38, commands from topic 105 (Sentinel)
- **Commands**: `/status`, `/check`, `/heal`, `/reboot`, `/logs`, `/help`
- **Deps**: aiohttp, requests (system pip)

### Postiz on omos-titan
- **URL**: `http://100.103.197.23:4007`
- **Stack**: Docker Compose — postiz, postgres, redis, temporal, elasticsearch, temporal-ui, temporal-admin-tools, spotlight
- **Config**: `~/postiz/.env` + `~/postiz/docker-compose.yaml`
- **JWT**: stored in .env on server
- **Status**: HTTP 307 (redirects to login) — confirmed healthy

## Fleet Reference

| Machine | Tailscale IP | User | Status |
|---|---|---|---|
| omos-core | 100.123.126.34 | ubuntu | Active — OpenClaw primary |
| omos-sentinel | 100.66.235.38 | ubuntu | Active — Sentinel medic |
| omos-titan | 100.103.197.23 | ubuntu | Active — Postiz |
| zeus-mac | 100.114.142.55 | zeuslegacy | Active — dev machine |
| omos-dell | 100.66.84.71 | zeus | Active |
| omos-ha | 100.64.25.12 | root | Offline — LAN only |

## Pending (manual)
- Destroy DO droplet 567310143 (needs fresh DO token)
- Remove dead Tailscale node `omos-core-1` at login.tailscale.com/admin/machines


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/operations (Ops Procedures)/security (Security Ops)/security|security]]
