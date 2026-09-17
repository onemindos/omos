---
id: recommended-agent-architecture-and-homestead-scenarios
title: "Recommended Agent Architecture & Homestead Scenarios"
type: note
source: apple-notes
imported: 2026-03-07
status: staged
tags:
  - apple-notes-import
  - openclaw
  - architecture
  - agents
  - homestead
  - codex
---

# Recommended Agent Architecture & Homestead Scenarios

## 4. Recommended Agent Architecture

Based on your screenshots (Route Method vs Terminal Method) and how OpenClaw natively works:

**Use the Route Method — it's already your architecture.**

OpenClaw's system IS the Route Method:

- Shared project (one gateway)
- Sub-agents share the same toolset (with `allow` filters)
- Quick setup — just create workspace folders
- Specialists (HP/LE/GE) in one project, like departments in the same office

The Terminal Method (separate workspaces, full isolation) would mean running 4 separate OpenClaw gateways. That's overkill — you'd need 4x the resources, 4x the tokens, and lose the native `sessions_send` inter-agent communication.

### Sub-Agent Spawning — How Each Agent Gets 5 Workers

OpenClaw supports this natively via `subAgents` config. Here's how it should work:

```
LEGACY (orchestrator, gpt-5.2)
├── spawn: researcher    (gpt-4o-mini, web_search + web_fetch)
├── spawn: file-worker   (gpt-4o-mini, read_file + exec)
├── spawn: code-reviewer (gpt-4o-mini, exec + github)
├── spawn: briefing-bot  (gpt-5-nano, memory + cron)
└── spawn: dispatcher    (gpt-5-nano, sessions_send)

HP (coach, gpt-4o)
├── spawn: meal-planner   (gpt-4o-mini, web_search + memory)
├── spawn: workout-log    (gpt-5-nano, memory + exec)
├── spawn: habit-tracker  (gpt-5-nano, memory + cron)
├── spawn: finance-scan   (gpt-4o-mini, web_fetch)
└── spawn: sleep-analyst  (gpt-5-nano, memory)

LE (steward, gpt-4o)
├── spawn: home-monitor   (gpt-5-nano, nodes + wireclaw)
├── spawn: family-cal     (gpt-4o-mini, calendar + memory)
├── spawn: maintenance    (gpt-5-nano, memory + cron)
├── spawn: estate-planner (gpt-4o-mini, web_search + memory)
└── spawn: guest-manager  (gpt-5-nano, nodes + locks)

GE (advisor, gpt-5.2)
├── spawn: market-scan    (gpt-4o-mini, web_search + web_fetch)
├── spawn: revenue-track  (gpt-4o-mini, memory + exec)
├── spawn: content-draft  (gpt-4o, memory)
├── spawn: deal-analyst   (gpt-4o-mini, web_search)
└── spawn: community-mod  (gpt-5-nano, discord + memory)
```

---

## Real Homestead Use Cases (Live Scenarios)

### Morning (7am) — Legacy orchestrates

1. Legacy's heartbeat fires → reads HEARTBEAT.md
2. Spawns `briefing-bot` sub-agent → compiles calendar + weather + priorities
3. `sessions_send` to HP: "What's Zeus's workout scheduled for today?"
4. `sessions_send` to LE: "Any homestead alerts overnight?"
5. `sessions_send` to GE: "Any revenue notifications?"
6. Legacy synthesizes all 3 responses into morning brief → sends to Discord/WhatsApp

### Greenhouse Alert (anytime) — LE handles

1. WireClaw sensor: `wireclaw.greenhouse.zone-a.soil < 25%`
2. NATS → Legacy's nexus-bridge picks up event
3. Legacy routes to LE: "Zone A soil critically dry"
4. LE spawns `home-monitor` sub-agent → checks FarmBot status, last irrigation time
5. LE decides: trigger irrigation via HA → reports result back to Legacy
6. Legacy notifies Zeus: "Greenhouse Zone A irrigation triggered, soil was at 22%"

### Guest Arrival — LE + HP coordinate

1. Zeus says on iOS: "Legacy, prepare the guest house for Mom arriving Friday"
2. Legacy delegates to LE via `sessions_send`
3. LE spawns `guest-manager`: check lock status, set temp to 72F, verify WiFi
4. LE spawns `maintenance`: verify last cleaning date, check supplies
5. LE reports to Legacy: "Guest house prepped: lock set, AC running, last cleaned Monday"
6. Legacy routes to HP: "Note: Mom visiting this weekend, adjust meal plan"
7. HP spawns `meal-planner`: plan meals for 2 extra people

### Business Intelligence — GE crunches

1. Zeus asks: "Legacy, what's the greenhouse ROI this month?"
2. Legacy delegates to GE
3. GE spawns `revenue-track`: pull yield data from FarmBot logs
4. GE spawns `market-scan`: check current market prices for crops
5. GE compiles: "40lbs harvested, market value $2.10/lb = $84 revenue. Energy cost: $12. Water: $3. Net: $69. Annualized ROI: 22%"

---

## 6. Codex — Yes, Keep Local + Git Sync

Your instinct is right. The Codex should:

1. **Live on your Mac** as the primary Obsidian vault
2. **Sync to GitHub** via git (already happening)
3. **Sync to VPS** via Syncthing over Tailscale → `/home/onemind/codex/`
4. **Legacy reads it** via `memorySearch` + `read_file` on the VPS copy

This gives you local-first editing in Obsidian + Legacy has full filesystem access to the vault on the VPS. The `memorySearch` config needs `extraPaths` pointing to the Codex:

```json
"memorySearch": {
    "extraPaths": ["/home/onemind/codex"]
}
```
