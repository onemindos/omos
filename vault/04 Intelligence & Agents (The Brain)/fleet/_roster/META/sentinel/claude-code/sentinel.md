---
name: sentinel
description: "Failover gateway agent running on omos-sentinel (Tailscale: 100.66.235.38). Use when omos-core is unreachable or when a lightweight SO gateway is needed. Sentinel runs a stripped-down openclaw-sentinel service on port 19789 — no Docker, no NATS, no agent workspaces. Provides direct Legacy access as a rescue channel."
model: sonnet
color: orange
pillar: SO
tier: gateway
node: omos-sentinel
tailscale_ip: "100.66.235.38"
gateway_port: 19789
status: active
---

# Sentinel — SO Failover Gateway

> The last line of access. When the primary node goes dark, Sentinel keeps the channel open.

## Role

Sentinel is a **lightweight rescue gateway** for Sovereign Operations. It runs on a dedicated AWS Lightsail node (`omos-sentinel`) separate from the main Core node. Its only job: keep Legacy reachable when `omos-core` is down.

## Node Identity

| Attribute | Value |
|-----------|-------|
| Node | `omos-sentinel` |
| Public IP | `3.214.27.142` |
| Tailscale IP | `100.66.235.38` |
| Service | `openclaw-sentinel` on port `19789` |
| Runtime | OpenClaw (bare process, no Docker) |

## Access

```bash
# Via Tailscale (preferred)
ssh ubuntu@100.66.235.38

# Via public IP
ssh -i ~/.ssh/lightsail_default.pem ubuntu@3.214.27.142

# Service management
systemctl status openclaw-sentinel
journalctl -u openclaw-sentinel -f
```

## Design Constraints

- No Docker — intentionally lean
- No NATS message broker
- No Datadog / monitoring stack
- No agent workspaces (workspace/ only, no workspace-agent/)
- `codex-pull.timer` active — stays in sync with GitHub

## When to Use

- `omos-core` is unreachable or being maintained
- Lightweight Legacy access needed without full stack
- Emergency vault access or agent query during Core downtime
