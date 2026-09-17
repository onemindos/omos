---
title: "OpenClaw — Event Hooks"
type: reference
status: active
created: 2026-05-10
modified: 2026-05-10
owner: zeus
tags: [type/reference, domain/SO-02, openclaw/hooks]
---

# OpenClaw — Event Hooks

Hooks are event-driven scripts that fire on specific OpenClaw gateway events. Installed at `~/.openclaw/hooks/` on omos-core.

---

## Active Hooks

### `onemind-session-logger`

| Field | Value |
|-------|-------|
| **Type** | Internal hook |
| **Status** | ✅ Enabled |
| **Install path** | `~/.openclaw/hooks/onemind-session-logger/` |
| **Source path** | `/home/ubuntu/onemind/_system/openclaw/hooks/onemind-session-logger/` |
| **Installed** | 2026-05-10 |
| **Handler** | `handler.js` |

**What it does:** Logs every agent session to the vault at `/home/ubuntu/onemind/_codex/memory/sessions/`. Each session becomes a dated `.md` file that feeds into the memory pipeline. (The `_codex/` path is a vault-internal folder, not the deprecated "codex" project name.)

**Config in `openclaw.json`:**
```json
"hooks": {
  "internal": {
    "enabled": true,
    "entries": {
      "onemind-session-logger": { "enabled": true }
    }
  }
}
```

---

## Managing Hooks

```bash
# List all hooks
ssh ubuntu@100.123.126.34 "openclaw hooks list"

# Check hook status
ssh ubuntu@100.123.126.34 "openclaw hooks status onemind-session-logger"

# View hook logs
ssh ubuntu@100.123.126.34 "ls ~/.openclaw/hooks/onemind-session-logger/"

# Reload after changes
ssh ubuntu@100.123.126.34 "kill -HUP \$(pgrep -f 'openclaw.*gateway')"
```

Source code: `_system/openclaw/hooks/onemind-session-logger/handler.js`

---

## Related

- [[../cron (Scheduled Jobs)/CRON|CRON]] — Scheduled jobs
- [[../plugin-skills (Custom Skills)/PLUGIN-SKILLS|PLUGIN-SKILLS]] — Agent skills
