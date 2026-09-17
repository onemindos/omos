---
title: "Interface — Discord"
type: reference
status: active
created: 2026-04-12
modified: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO, domain/openclaw, interface/discord]
ai_summary: "Discord as an OpenClaw interface — live config, agent channel map, operational categories, streaming, approvals, presence, voice."
---

# Interface: Discord

**Status:** ✅ Active — primary desktop interface
**Bot:** `@OneMind OS` (`1453426935442047191`) on server `OnemindOS` (guild `1453411552181883027`)
**Node:** omos-core (`100.123.126.34`) — ubuntu user service, port 18789
**Bot token:** 1Password → `op://aaeyris5z44pdqspk277u2oppi/qwomupt2iajig4twcv6ijdkyj4/token`
**Authorized user:** Zeus Discord ID `1024395579742228570`

## Live Status (2026-05-11)

| Item | Status | Notes |
|------|--------|-------|
| Bot connected | ✅ | `@OneMind OS` resolved, zero unresolved channels |
| Agent channels | ✅ 6 resolved | legacy, zenith, unity, eternal, sovereignty, sentinel |
| Utility channels | ✅ 3 | notifications, agent-heartbeat, morning-brief (outbound only) |
| Operational channels | ✅ 11 | approvals, inbox, directives, reports, health, pipeline |
| Stale channel IDs | ✅ Cleaned | 17 old IDs removed — no more startup noise |
| Message Content Intent | ⚠️ Limited | Works fine; only required for 100+ server bots |
| Exec approvals | 🔶 Disabled | Requires `gateway.auth.mode=token` — currently `none`. Chat routing unaffected. |
| Voice model | ✅ Active | `gpt-realtime-1.5` deployed on `onemind-voice.openai.azure.com` (East US 2) |
| Talk mode | ✅ Active | Azure OpenAI Realtime `gpt-realtime-1.5` via `onemind-voice.openai.azure.com` |
| streaming | ✅ | `progress` mode, toolProgress enabled |
| autoPresence | ✅ | Bot status reflects node health |

---

## Known Issues

**1. Exec approvals disabled** — `execApprovals.enabled: false` because `gateway.auth.mode=none` conflicts with the approval sub-client which needs device auth. To enable: set `gateway.auth.mode=token`, pair all clients, then re-enable `execApprovals.enabled: true`. Until then, high-risk agent actions go unapproved — agents still ask in chat but no button UI.

**2. Discord intents limited** — Message Content Intent limited for bots under 100 servers. Works fine. Only matters at scale.

**3. Voice `/vc join` ready to test** — `voice.model=gpt-realtime-1.5` configured, deployment live on `onemind-voice.openai.azure.com`. Test: `/vc join channel:1478851927621111808` from Discord.

---

## Live Config (as of 2026-05-11)

```json
{
  "channels": {
    "discord": {
      "enabled": true,
      "allowFrom": ["1024395579742228570"],
      "ackReaction": "⚡",
      "replyToMode": "first",
      "streaming": {
        "mode": "progress",
        "progress": { "toolProgress": true, "maxLines": 8 }
      },
      "autoPresence": {
        "enabled": true,
        "intervalMs": 30000,
        "healthyText": "One Mind OS — Online",
        "degradedText": "One Mind OS — Degraded",
        "exhaustedText": "One Mind OS — {reason}"
      },
      "execApprovals": {
        "enabled": false,
        "approvers": ["1024395579742228570"],
        "cleanupAfterResolve": true,
        "target": "channel"
      },
      "voice": {
        "model": "gpt-realtime-1.5"
      },
      "threadBindings": {
        "enabled": true,
        "spawnSessions": true,
        "defaultSpawnContext": "isolated",
        "idleHours": 24,
        "maxAgeHours": 0
      },
      "groupPolicy": "allowlist",
      "guilds": {
        "1453411552181883027": {
          "channels": {
            "1478604243509448705": { "autoThread": true, "autoArchiveDuration": 10080 },
            "1482355265092190319": { "autoThread": true, "autoArchiveDuration": 10080 },
            "1503272334990118952": { "autoThread": true, "autoArchiveDuration": 10080 },
            "1503272335921123369": { "autoThread": true, "autoArchiveDuration": 10080 },
            "1503272336978087986": { "autoThread": true, "autoArchiveDuration": 10080 },
            "1503272338274390076": { "autoThread": true, "autoArchiveDuration": 10080 },
            "1478604241013833859": { "enabled": false },
            "1478604235494133801": { "enabled": false },
            "1478604238287536148": { "enabled": false },
            "1503277446475157625": {}, "1503277448924500061": {}, "1503277451826954300": {},
            "1503277454318637157": {}, "1503277456939814913": {}, "1503277459565445201": {},
            "1503277462686011403": {}, "1503277465529745438": {}, "1503277468340060230": {},
            "1503277470982475876": {}, "1503277473440333978": {}
          }
        }
      }
    }
  }
}
```

---

## Agent → Channel Map

| Agent | Channel | Channel ID |
|-------|---------|------------|
| Legacy (main) | #legacy-chat | `1478604243509448705` |
| Zenith (HP) | #zenith-chat | `1503272334990118952` |
| Unity (LE) | #unity-chat | `1503272335921123369` |
| Eternal (GE) | #eternal-chat | `1503272336978087986` |
| Sovereignty (SO) | #sovereignty-chat | `1503272338274390076` |
| Sentinel (omos-sentinel) | #sentinel-chat | `1482355265092190319` |

**Utility channels (outbound only — inbound disabled):**

| Channel | Channel ID | Purpose |
|---------|------------|---------|
| #notifications | `1478604241013833859` | System alerts, webhook targets |
| #agent-heartbeat | `1478604235494133801` | Agent health reports |
| #morning-brief | `1478604238287536148` | Legacy daily 7AM briefing |

---

## Operational Channel Groups

### COMMAND & CONTROL (category `1503277359586082847`)

| Channel | Channel ID | Purpose |
|---------|------------|---------|
| #approvals | `1503277446475157625` | Exec approval requests — `execApprovals.target` |
| #inbox | `1503277448924500061` | Incoming captures queued for triage |
| #directives | `1503277451826954300` | Standing orders and mission-level directives |

### REPORTING (category `1503277390057439292`)

| Channel | Channel ID | Purpose |
|---------|------------|---------|
| #daily-reports | `1503277454318637157` | Agent daily summaries (cron output) |
| #weekly-reports | `1503277456939814913` | Agent weekly summaries |
| #briefings | `1503277459565445201` | Morning brief and strategic briefings |

### SYSTEM HEALTH (category `1503277392603644167`)

| Channel | Channel ID | Purpose |
|---------|------------|---------|
| #system-status | `1503277462686011403` | Node health, uptime, connectivity |
| #errors | `1503277465529745438` | Agent errors, exceptions, failure alerts |
| #cron-runs | `1503277468340060230` | Cron job execution logs |

### WORK PIPELINE (category `1503277395220758559`)

| Channel | Channel ID | Purpose |
|---------|------------|---------|
| #active-work | `1503277470982475876` | Tasks currently in progress |
| #completed | `1503277473440333978` | Finished tasks and deliverables log |

---

## Capabilities

- **Streaming:** `progress` mode — live tool progress shown as agents work, max 8 lines
- **Threads:** Every message auto-creates a thread (7-day archive, 24h session idle expiry)
- **Slash commands:** Auto-registered — `/model`, `/session`, `/focus`, `/unfocus`, `/diagnostics`, etc.
- **Approval UI:** Native Discord buttons in `#approvals` — approve/deny with one click, auto-cleaned after resolve
- **Auto-presence:** Bot status reflects node health in real time
- **Ack reaction:** ⚡ appears on your message while the agent is processing
- **Reply threading:** Quote-reply on first message in each exchange
- **Voice:** Azure OpenAI Realtime via `/vc` commands — `agent-proxy` mode, `gpt-realtime-1.5` default, `gpt-realtime-mini` fallback
- **Rich embeds:** Webhooks to utility channels for structured alerts/scorecards

---

## Voice

**Status:** Config documented — pending activation on omos-core.

Discord voice runs through `@discordjs/voice` (Opus codec, 48kHz). The bot joins a voice channel, receives per-user Opus streams, and routes through the realtime provider.

### Provider: Azure OpenAI Realtime — `onemind-voice.openai.azure.com` (East US 2)

Dedicated voice resource with `gpt-realtime-1.5` and `gpt-realtime-mini` deployed. Key: `AZURE_OPENAI_VOICE_KEY` in `gateway.systemd.env`.

Discord voice schema only allows `enabled`, `model`, `autoJoin`, `daveEncryption`, and timing keys — Azure routing is handled by the gateway's `talk` config, not per-channel.

```json
"channels": {
  "discord": {
    "voice": {
      "enabled": true,
      "model": "gpt-realtime-1.5"
    }
  }
}
```

The `talk` section holds the full Azure config pointing to `onemind-voice.openai.azure.com`. See `TALK-MODE-CONFIG.md`.

### Voice Channel

| Channel | ID | Purpose |
|---------|-----|---------|
| legacy-voice | `1478851927621111808` | Primary voice channel for Legacy |

### Commands

```
/vc join channel:1478851927621111808   — join legacy-voice
/vc leave                              — leave voice
/vc status                             — check voice status
```

### Requirements

- `ffmpeg` and `ffprobe` installed on omos-core ✅
- `gpt-realtime-1.5` deployed on `onemind-voice.openai.azure.com` ✅
- `AZURE_OPENAI_VOICE_KEY` in `gateway.systemd.env` ✅

### Cost (Azure Realtime)

| Model | Est. cost/min |
|-------|--------------|
| `gpt-realtime-1.5` (default) | ~$0.21–0.40 |
| `gpt-realtime-mini` (fallback) | ~$0.09–0.21 |

Single-user (Zeus only) at ~30 min/day = ~$3–12/day — covered by Azure credits.

---

## Owner Commands

Available to Zeus only (gated by `commands.ownerAllowFrom`):

```
/diagnostics       — system health check
/export-trajectory — export session history
/config set        — update config live
/session           — inspect/manage sessions
/focus             — bind thread to session
```

---

## Related

- [[CHANNELS]] — Full channel map for both platforms
- [[00-24 SO (Sovereign Operations)/05 Operations (The Controls)/channels/channels/telegram/telegram]] — Telegram interface config
- [[MESSAGING]] — Bot registry
- [[omos-core]] — Primary node
