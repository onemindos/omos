# NATS Subscriber Build — Session Handoff

> Context compaction from a previous session. Treat as reference only.

---

## Active Task

User confirmed the config-driven subscriber architecture and said "build it." Legacy asked one clarification: run on Atlas or MBP? User has not yet answered. Build is ready to start the moment that's confirmed.

---

## Goal

Build out OneMind OS NATS mesh — universal control plane for all machines, SaaS, AI agents, robotics, drones, devices. Build autonomous subscriber layer so world events trigger actions without Legacy in the loop.

---

## Constraints & Preferences

- User is Zeus
- Context: home-office / active / online, Mac in clamshell 24/7
- Prefer live tool validation over theoretical answers
- Keep things concrete and actionable
- Overwatch cron job: every 1m, still paused until explicitly enabled
- Subscriber architecture must be config-driven via NATS KV (not hardcoded)
- Subscriber rules must be changeable live without restarts or code changes
- Legacy manages subscriber config via natural language → KV writes
- Subscriber runs on Atlas (likely) — pending final confirmation from Zeus

---

## Completed Actions

1. CONFIRMED NATS MCP tools live — 6 tools active: `nats_kv_get`, `nats_kv_put`, `nats_kv_keys`, `nats_publish`, `nats_request`, `nats_stream_read`
2. READ KV bucket zeus-context — keys: location=home-office, focus_mode=active, status=online
3. READ all JetStream streams — AGENT_EVENTS, COMMS, SYSTEM, TASKS, SAAS_EVENTS all empty
4. PUBLISHED test event to `events.legacy.comms_check` — confirmed published
5. REQUEST `cmd.ha.service.call` — confirmed `{"ok": true}`, HA is live
6. QUERIED `$JS.API.STREAM.LIST` — 11 streams total
7. READ KV_zeus-context stream — 8 messages, location history: home-office → farm → gym → home-office
8. READ KV_agent-state stream — 2 messages: legacy.status=idle, hermes.status=idle
9. READ KV_config — nats.version=2.14.1, setup.date=2026-05-21
10. QUERIED `$JS.API.INFO` — JetStream live: 11 streams, 0 consumers, 769 bytes
11. CREATED cron job "Overwatch" — job_id: `3637ca23b2d2`, passive NATS surveillance, Telegram alert on notable events
12. PAUSED Overwatch — state: paused, enabled: false
13. QUERIED HA states via `cmd.ha.states.list` — 183 entities total
14. DISCOVERED plug/switch entities — saved full state dump to `/tmp/ha_states.json` (127105 bytes)
15. UPDATED Overwatch schedule — every 30m → every 1m, still paused
16. TURNED ON switch.plug_mini_ef5a — `{"ok": true}`
17. TURNED OFF switch.plug_mini_ef5a — `{"ok": true}`
18. CYCLED switch.plug_mini_ef5a ON/OFF x3 — all 6 calls confirmed ok, plug ends OFF
19. SEARCHED session history for NATS/NGS setup — found May 21 Synadia research session
20. READ nats-server skill — confirmed architecture: Atlas Lightsail, systemd, JetStream, Tailscale-only, 6 MCP tools

---

## Active State

### NATS Mesh
- JetStream: live, 11 streams, 0 consumers, 769 bytes, 0 errors
- KV buckets with data: zeus-context (8 msgs), agent-state (2 msgs), config (2 msgs)
- KV buckets empty: connections, schedules, triggers
- Event streams all empty: AGENT_EVENTS, COMMS, SYSTEM, TASKS, SAAS_EVENTS
- KV bucket "rules" — does not exist yet, will be created during subscriber build

KV zeus-context: location=home-office, focus_mode=active, status=online
KV agent-state: legacy.status=idle, hermes.status=idle

### Home Assistant — 183 Entities
- Zeus iPhone: 15 Harold Ave Greenwich CT 06830, WiFi: onemindos, battery 75%
- Meter Plus 7616: 74.3°F / 53% humidity
- Meter Plus CEFF: 73.76°F / 52% humidity
- Verizon Router External IP: 74.108.23.135
- switch.plug_mini_a886 = on, 0.0W
- switch.plug_mini_ef5a = off, 0.0W (last controlled)
- switch.plug_mini_6826 = unavailable (offline)
- HA Core: 2026.4.4 installed / 2026.5.3 available

Fabric agents: CC, PI, OpenClaw, Hermes — 4/4 online

---

## In Progress

Subscriber build — waiting on Zeus to confirm run location (Atlas vs MBP) before writing code.

---

## Blocked

- switch.plug_mini_6826 unavailable — offline, needs SwitchBot hub reconnect
- HA MCP server (mcp_home_assistant_search_entities) was flapping — workaround via NATS + terminal
- Subscriber build blocked on: Atlas vs MBP confirmation from Zeus

---

## Key Decisions

- Subscriber architecture: config-driven via NATS KV bucket "rules" — rules readable/writable live, no restarts needed
- Multiple subscribers run simultaneously, each watches different subject namespace, each has own KV rules
- Mac (clamshell 24/7) is valid primary subscriber host; Atlas is fallback for travel/remote
- Queue groups prevent double-firing when both Mac + Atlas subscribers run simultaneously
- Legacy manages subscriber rules via natural language → KV writes (no SSH, no code changes)
- All device control goes through NATS `cmd.ha.service.call` — scales to drones, robots, cloud
- Overwatch runs every 1m (minimum scheduler interval); left paused intentionally
- For true 30s polling loops, Atlas/Mac watchdog script (no LLM per tick) is better than cron
- NATS MCP tools for agent/decision logic; NATS CLI for shell scripts and polling loops
- NATS is native (not a wrapper) — MCP tools talk directly to NATS server, no middleware

---

## Resolved Questions

| Question | Answer |
|----------|--------|
| Can leaf nodes run on anything? | Yes — MacBook, cloud, Pi, Jetson, drone, phone, microcontroller |
| Is JetStream live? | Yes — confirmed with stats |
| Does Legacy have active KV subscriptions? | No — 0 consumers, reads on-demand only |
| What's live in HA? | 183 entities confirmed |
| Can Legacy control plugs via NATS? | Yes — confirmed ON/OFF/x3 of EF5A this session |
| NATS vs direct HTTP? | Similar single-device latency; NATS wins at scale: parallel fan-out, works across firewalls, uniform interface, audit trail |
| What is an Atlas subscriber? | 24/7 Python process on Atlas/Mac connected to NATS, catches events and executes protocols without LLM involvement — world events trigger actions autonomously |
| Can subscriber run on Mac? | Yes — Mac in clamshell is 24/7, valid primary host. Run on both Mac + Atlas with queue groups for failover |
| Why is subscriber more powerful than Legacy-in-loop? | You are no longer the trigger — world events (geofence, motion, battery, CPU) trigger immediate automated responses. Legacy becomes advisor not hands |
| What if I want to change subscriber behavior? | Rules stored in NATS KV bucket. Update a KV key → subscriber picks up change live. No code edit, no restart, no SSH |
| What if I want multiple subscribers? | Each is an independent named process watching different subjects (home.sub, security.sub, health.sub, comms.sub). All run simultaneously |
| nats.swift — app or library? | Swift Package library, embedded in custom iOS app |
| NGS super-cluster? | Synadia's managed global NATS backbone. Current setup is self-hosted. Research saved to `/Users/onemind/Documents/onemind/synadia_research_summary.md` |
| Are we doing this the native way? | Yes — native NATS leaf node + JetStream, MCP tools talk directly to NATS server, no middleware |
| NATS MCP vs NATS CLI? | MCP for agent decision logic and agent-to-agent comms; CLI for shell scripts and sub-minute polling loops on Atlas |

---

## Pending User Asks

- Where should the subscriber run — Atlas or MBP? (Legacy defaulted to Atlas as the right answer; Zeus has not confirmed)

---

## Relevant Files

- `/tmp/ha_states.json` — full HA entity state dump, 127105 bytes, 183 entities
- `/Users/onemind/Documents/onemind/synadia_research_summary.md` — Synadia NGS research, 300+ lines, all plans/pricing/capabilities

---

## Remaining Work

- **NEXT:** Confirm subscriber host (Atlas vs MBP), then build:
  - `subscriber.py` — universal daemon, NATS KV watch, action registry
  - KV bucket "rules" — seed with example rules (on_left_home, on_low_battery, on_motion_2am, battery_threshold)
  - Systemd unit template — parameterized by subscriber name
  - Initial 4 subscribers: home.sub, security.sub, health.sub, comms.sub
- Define 3 mode protocols: Incognito, High Alert, Watch Me
- Wire Telegram → NATS → mode activation path
- Enable Overwatch when ready: `cronjob resume 3637ca23b2d2`
- Decide whether to upgrade HA Core 2026.4.4 → 2026.5.3
- Fix unavailable SwitchBot plug 6826
- Scaffold iOS app with nats.swift
- Consider NGS migration for always-on backbone
- Configure lights, climate, additional switch integrations in HA

---

## Critical Context

| Key | Value |
|-----|-------|
| Overwatch job_id | `3637ca23b2d2` — paused, schedule: every 1m, toolset: terminal |
| NATS version | 2.14.1, self-hosted on Atlas (Lightsail), systemd daemon, Tailscale-only |
| NATS bind | Tailscale interface only — never exposed publicly |
| HA NATS bridge subjects | `cmd.ha.service.call`, `cmd.ha.states.list` |
| Home WiFi SSID | onemindos |
| Zeus iPhone entity prefix | `sensor.zeus_*` |
| External IP | 74.108.23.135 |
| Atlas server | AWS Lightsail |
| SwitchBot entity IDs | switch.plug_mini_a886, switch.plug_mini_ef5a, switch.plug_mini_6826 |

### Planned KV Bucket "rules" Keys
- `rules.on_left_home`
- `rules.on_low_battery`
- `rules.on_motion_2am`
- `rules.battery_threshold`

### Planned Action Registry Mappings
- lock_doors → `cmd.ha.service.call`
- lights_off → `cmd.ha.service.call`
- thermostat_away → `cmd.ha.service.call`
- notify → `cmd.comms.telegram`
- high_alert → `cmd.mode.security.alert`
- record → `cmd.camera.record`

### Planned Subscriber Names + Subjects
- home.sub → `events.zeus.*`
- security.sub → `events.motion.*`
- health.sub → `events.device.*`
- comms.sub → `cmd.mode.*`
