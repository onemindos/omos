---
name: OpenClaw CLI — Crestodian
updated: 2026-05-11
---

# `openclaw crestodian`

Crestodian is OpenClaw's local setup, repair, and configuration helper. Designed to stay reachable when the normal agent path is broken — broken config, gateway down, no agent configured yet.

**Official docs:** [docs.openclaw.ai/cli/crestodian](https://docs.openclaw.ai/cli/crestodian.md)

> Running `openclaw` with **no command** starts Crestodian in an interactive terminal.

---

## Commands

```bash
openclaw                                          # Start Crestodian (no command = Crestodian)
openclaw crestodian                               # Same, explicit
openclaw crestodian --json
openclaw crestodian --message "models"
openclaw crestodian --message "validate config"
openclaw crestodian --message "set default model openai/gpt-5.5" --yes
openclaw onboard --modern                         # Crestodian as conversational onboarding
```

---

## Inside the Crestodian TUI

```bash
status                                            # Full inventory: config, agents, model, gateway
health
doctor
doctor fix
validate config
setup
setup workspace ~/Projects/work model openai/gpt-5.5
config set gateway.port 19001
config set-ref gateway.auth.token env OPENCLAW_GATEWAY_TOKEN
gateway status
restart gateway
agents
create agent work workspace ~/Projects/work
models
set default model openai/gpt-5.5
plugins list
plugins search slack
plugin install clawhub:openclaw-codex-app-server
audit
talk to work agent                                # Switch to normal agent TUI
quit
```

Return to Crestodian from the normal TUI:
```
/crestodian
/crestodian restart gateway
```

---

## What Crestodian Can Do Without Normal Agent Path

**Read-only (always available):**
- Show overview, list agents, list/search plugins
- Show model/backend status
- Run status or health checks
- Check gateway reachability
- Run doctor without fixes
- Validate config

**Persistent operations (require approval):**
- Write config / run `config set`
- Set SecretRef values via `config set-ref`
- Run setup/onboarding
- Change default model
- Start/stop/restart gateway
- Create agents
- Install/uninstall plugins
- Run doctor repairs

Applied writes are audited to: `~/.openclaw/audit/crestodian.jsonl`

---

## Message Rescue Mode

Crestodian also works as a **message-channel rescue** when the normal agent is down but a trusted channel (Telegram, Discord) still receives commands.

```
You (trusted DM): /crestodian status
OpenClaw: Crestodian rescue mode. Gateway reachable: no. Config valid: no.
You: /crestodian restart gateway
OpenClaw: Plan: restart the Gateway. Reply /crestodian yes to apply.
You: /crestodian yes
OpenClaw: Applied. Audit entry written.
```

Config:
```json
{
  "crestodian": {
    "rescue": {
      "enabled": "auto",
      "ownerDmOnly": true
    }
  }
}
```

`enabled` values: `"auto"` (default — allow when YOLO + no sandbox), `false` (never), `true` (always allow when owner checks pass).

---

## Safe Startup

Crestodian starts even when:
- `openclaw.json` is missing or invalid
- The Gateway is down
- Plugin command registration is unavailable
- No agent has been configured yet

---

## One Mind Use Case

When omos-core's Legacy agent is unresponsive, send `/crestodian status` from Telegram to get a rescue status check without needing SSH. If the gateway needs a restart: `/crestodian restart gateway` → `/crestodian yes`.
