---
title: "Interface — onemind-web (SwarmClaw)"
type: reference
status: active
created: 2026-04-12
owner: zeus
tags: [type/reference, domain/SO, domain/openclaw, interface/web]
ai_summary: "onemind-web / SwarmClaw web UI — URL, auth, Docker deploy, capabilities, dev commands."
---

# Interface: onemind-web (SwarmClaw)

**Status:** ✅ Running — full-featured web UI
**Type:** Web application
**Access:** `http://100.77.52.93:3456` (Tailscale only)

---

## Access

```
URL:         http://100.77.52.93:3456
WebSocket:   ws://100.77.52.93:3457
Auth:        Access key (auto-generated on first run, stored in /opt/onemind-web/data/)
Credentials: 1Password → SO vault → OneMind Web Access Key
Network:     Tailscale mesh only — not publicly accessible
```

---

## Capabilities

- Full session management and workspace switching
- All 9 agents accessible via UI
- Extension (plugin) management
- Tool configuration per agent
- Real-time streaming responses
- File upload, image support
- Agent profiles and persona settings
- Model switching per session
- 3D office visualization (React Three Fiber) — `src/office/`

---

## Start / Stop

```bash
# Production (on VPS):
cd /opt/onemind-web
docker compose up -d         # start
docker compose down          # stop
docker compose logs -f       # live logs

# Local dev (Mac):
cd ~/Documents/onemind-web
npm run dev                  # port 3456, Turbopack
```

---

## Notes

- Production runs as Docker container on the VPS — see [DEPLOYMENT.md](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/openclaw/DEPLOYMENT.md)
- Source code: `onemind-web` repo (`~/Documents/onemind-web`)
- Full architecture in `onemind-web/AGENTS.md` and `onemind-web/CLAUDE.md`
- Intentionally kept `swarmclaw` references in internal keys — do not rename (breaks deployments)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
