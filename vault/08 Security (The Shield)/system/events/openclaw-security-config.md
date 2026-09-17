# OpenClaw Security Configuration
**Last updated:** 2026-04-26 | **Maintained by:** Legacy ⚡

---

## Exec Security Tiers

All agents are tiered by trust level. Configuration lives in `~/.openclaw/openclaw.json`.

| Agent | Level | Justification |
|-------|-------|---------------|
| ⚡ Legacy (`main`) | `full` | Central orchestrator — needs unrestricted shell for all cross-domain ops |
| ⚙️ Grid | `full` | VPS infra ops — systemctl, docker, cron, disk, networking |
| 🛡️ Guardian | `full` | Security hardening — needs root-level access to audit and fix |
| 🔺 Trinity | `full` | Meta-orchestration — coordinates multi-agent pipelines |
| 💪 Spartan | `allowlist + ask` | Domain agent — no direct infra needs; approval gate for any shell |
| 🔮 Oracle | `allowlist + ask` | Domain agent — finance/mind; approval gate for any shell |
| 🏡 Haven | `allowlist + ask` | Domain agent — home/family; approval gate for any shell |
| 💰 Empire | `allowlist + ask` | Domain agent — business/brand; approval gate for any shell |
| 🌱 Eden | `allowlist + ask` | Domain agent — environment; approval gate for any shell |

**Rule:** `full` = unrestricted shell. `allowlist + ask` = blocked unless whitelisted, prompts Zeus for approval on unknown commands.

---

## Gateway Security

| Setting | Value | Notes |
|---------|-------|-------|
| `controlUi.allowedOrigins` | Tailscale + localhost only | Wildcard `*` removed 2026-04-26 |
| `gateway.bind` | `lan` (0.0.0.0:18789) | Port should only be reachable via Tailscale or SSH tunnel |
| `gateway.auth.mode` | `token` | Auth token required for all API calls |

**Allowed CORS origins:**
- `http://localhost:18789`
- `http://127.0.0.1:18789`
- `http://100.98.25.10:18789` (VPS Tailscale IP)
- `http://onemindos-cloud-1.tail717ba0.ts.net:18789` (Tailscale hostname)

---

## Discord Channel Security

| Setting | Value |
|---------|-------|
| `dmPolicy` | `allowlist` |
| `groupPolicy` | `allowlist` |
| `allowFrom` | Zeus Discord ID only (`1024395579742228570`) |

No external users can trigger agents via Discord. Zeus-only.

---

## Periodic Security Checklist

Run quarterly or after any major config change:

- [ ] `openclaw status` — check for new audit criticals
- [ ] Review `controlUi.allowedOrigins` — no wildcards
- [ ] Review `groupPolicy` / `dmPolicy` — no `open`
- [ ] Review which agents have `exec=full` — justify each one
- [ ] Check `gateway.bind` — confirm port 18789 not exposed publicly
- [ ] Rotate gateway auth token if compromised
- [ ] Review `env` block in openclaw.json — rotate any stale API keys

---

## Change Log

| Date | Change | By |
|------|--------|----|
| 2026-04-26 | Wildcard CORS removed → Tailscale-only | Legacy ⚡ |
| 2026-04-26 | Discord groupPolicy: open → allowlist | Legacy ⚡ |
| 2026-04-26 | Exec tiering applied to all 9 agents | Legacy ⚡ |


---

## Related

- [[00-24 SO (Sovereign Operations)/04 Security (Cyber & Digital Shield) (The Shield)/04 Security|04 Security]]
