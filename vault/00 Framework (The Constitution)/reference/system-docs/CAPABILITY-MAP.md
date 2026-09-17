---
title: "Capability Map"
type: reference
status: active
created: 2026-04-17
updated: 2026-04-17
owner: zeus
tags: [type/reference, domain/SO, status/active]
ai_summary: "Master inventory of every capability in the OneMind OS stack. Every concept, system, and feature mapped to its codex location, build status, and owner domain. Single source of truth for what exists, what's partial, and what's missing."
---

# Capability Map

> The complete inventory of every capability in the OneMind OS.
> Status: ✅ built & documented | ⚠️ exists but undocumented/partial | ❌ missing entirely
> Updated every time a new system is built or a gap is closed.

---

## How to Read This

Each row is a **capability** — a thing the system can do or a concept it embodies.

| Column | Meaning |
|--------|---------|
| Capability | The word/concept |
| Domain | Where it lives in the codex (SO/HP/LE/GE or cross-cutting) |
| Status | ✅ / ⚠️ / ❌ |
| Codex Path | The actual file or folder |
| Notes | What exists, what's missing, what to do next |

---

## SO — Sovereign Operations (Operating Layer)

### SO/00 — Framework & Doctrine

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **domains** | ✅ | `CODEX-FRAMEWORK.md`, `AGENTS.md` | SO/HP/LE/GE 3+1 model — fully defined |
| **areas** | ✅ | Same as domains | PARA parallel — "areas" = domains in this system |
| **systems** | ✅ | `ONEMIND-CODEX.md` | The OS metaphor — SO as platform, HP/LE/GE as apps |
| **principles** | ⚠️ | `_codex/system/` (beliefs file) | Context file exists, no dedicated principles doctrine |
| **governance** | ⚠️ | `_codex/system/` + user access levels L0–L4 | User autonomy levels defined; agent autonomy rules thin |
| **paths** | ⚠️ | `SO/00 Framework/architecture/PATHS-SYSTEM.md` | Doc exists — not integrated into daily operation |
| **elevation** | ❌ | — | Autonomy escalation + personal growth arc — concept only, nothing built |
| **vision** | ⚠️ | `SO/00 Framework/core-doctrine/VISION.md` | Life vision doc exists; personal vision system not systematized |
| **missions** | ⚠️ | Referenced in `CODEX-FRAMEWORK.md` | Concept defined; no tracking or review system |
| **atlas** | ❌ | — | Visual map of domains/entities/relationships — strong concept, nothing built. Future: onemind-web world view |
| **launch pad** | ❌ | — | Daily mission kick-off surface — needs design. Concept: one place to start every day |

---

### SO/01 — Intelligence (AI Brain)

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **agents** | ✅ | `SO/03 Agents/` + 9 workspace files on VPS | Legacy, Spartan, Oracle, Eden, Haven, Empire, Guardian, Grid, Trinity — all populated |
| **models** | ⚠️ | `machines/vps.md` + openclaw.json on VPS | Ollama models configured; no single codex doc listing all models + roles |
| **memory** | ✅ | `_codex/memory/` | 3-tier: work / learning / relationship — fully built |
| **learning** | ✅ | `_codex/memory/learning/` | Signal capture + feedback loop — built |
| **reasoning** | ❌ | — | Chain-of-thought, qwen3 `/think` mode — model capability, no routing doc or agent config |
| **skills** | ✅ | `_codex/skills/` + `.claude/skills/` | 30+ skill packs — major system, well documented |
| **tools** | ⚠️ | Scattered across openclaw.json + skills | Running but no single inventory of all tools available to agents |
| **knowledge** | ✅ | Entire `_codex/` system | The vault IS the knowledge layer |
| **rag** | ⚠️ | `SO/01 Infrastructure/` (pgvector + postgres) | Infrastructure exists on VPS; pipeline from vault → embeddings → query not built |
| **semantic** | ⚠️ | VPS: pgvector + nomic-embed-text + memorySearch | Memory search active; full semantic pipeline not wired |
| **loras** | ❌ | — | Fine-tuned model adapters — not a current priority, nothing implemented |
| **multimodal** | ❌ | — | Vision/audio input to agents — model capability (qwen2.5-vl etc), no routing or doc |
| **tracing** | ❌ | — | Agent run observability (Langfuse/OpenTelemetry) — not implemented. High value gap |
| **guardrails** | ⚠️ | `_system/security/` | Security validation gate exists; agent output guardrails thin |
| **hivemind** | ⚠️ | Trinity agent (`SO/03 Agents/`) + `MULTI-USER.md` | Trinity coordinates agents; no hivemind doctrine or routing rules |
| **mcp** | ⚠️ | `openclaw.json` (servers: context7, github, codex-vault) | 3 MCP servers running on all agents; zero codex documentation |
| **world models** | ⚠️ | `SO/05 Environments/` entity registry | Entity registry started; live data not wired; digital twin concept only |

---

### SO/02 — Command (Control Center)

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **command** | ✅ | `SO/02 Command (Dashboard & Navigation)/` | Domain exists |
| **navigation** | ✅ | `_codex/INTERFACE.md` | Universal agent + human entry point |
| **dashboard** | ⚠️ | `onemind-web` (running on VPS:3456) | Web UI running; not linked to codex dashboard concept |
| **insights** | ❌ | — | AI-generated patterns from metrics/journal/tasks — high-value gap, nothing built |
| **activity feed** | ❌ | — | Live stream of agent events, completions, system changes — not built |

---

### SO/03 — Agents (Execution)

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **sessions management** | ✅ | `onemind-web/src/lib/server/runtime/session-run-manager.ts` | Full lifecycle: dedup, coalescing, heartbeat preemption, execution lock — not codex-documented |
| **runs** | ✅ | `onemind-web` runtime | Session runs fully implemented — not codex-documented |
| **execution engine** | ✅ | `onemind-web/src/lib/server/runtime/` | `enqueueSessionRun()` is the entry point — not codex-documented |
| **hitl** | ⚠️ | `onemind-web` (`durable_wait` terminal tool) | Human-in-the-loop approval gates built in code; no codex protocol |
| **heartbeat** | ⚠️ | `onemind-web` runtime | Periodic agent turns between messages — built, undocumented |
| **remote execution** | ⚠️ | VPS SSH + systemd | Happening constantly (onemind-watch, deploy scripts); no formal runbook |
| **guardrails** | ⚠️ | `_system/security/patterns.yaml` | Pre-execution validation exists; needs expansion |

---

### SO/05 — Environments (Physical World)

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **vps** | ✅ | `SO/05 Environments/machines/vps.md` | Documented — 104.248.4.212 / 100.77.52.93 Tailscale |
| **infrastructure** | ✅ | `SO/01 Infrastructure/` domain | Domain exists; detailed docs sparse |
| **spatial intel** | ⚠️ | `FIND-ZEUS-MODE.md` + Bouncie skill + HA skill | Pieces exist; not unified into one routing system |
| **world models** | ⚠️ | `SO/05 Environments/` | Entity registry built; live data connections missing |

---

### SO/06 — Protocols

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **protocols** | ✅ | `SO/06 Protocols (SOPs)/` | Domain exists |
| **routines** | ⚠️ | Morning briefing + evening review skills exist | No unified routine tracker or protocol |

---

### SO/09 — Automation (Event Layer)

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **workflows** | ✅ | `SO/09 Automation/` + n8n running on VPS:5678 | n8n active; flows not inventoried in codex |
| **hooks** | ✅ | `_system/hooks/` | 10 lifecycle hooks mapped to CODEX stages — built |
| **cron jobs** | ⚠️ | n8n schedules + systemd timers (onemind-pull.timer) | Running; no central schedule inventory |
| **event driven** | ⚠️ | inotifywait (onemind-watch) + NATS + n8n | All event subsystems running; no unified event architecture doc |
| **webhooks** | ⚠️ | n8n handles incoming; no outbound registry | No central webhook registry; each integration manages its own |
| **schedules** | ⚠️ | n8n + systemd timers | Running; no codex-documented schedule inventory |
| **nats** | ⚠️ | VPS port 4222/9222 | Messaging backbone running; zero codex documentation |
| **bus routing** | ⚠️ | NATS subject routing | Architecture exists; uncharted — no subject map |
| **gateway** | ✅ | `openclaw-gateway` systemd service (VPS:18789) | Running; `SO/02 Intelligence/openclaw/` docs light |
| **signals** | ✅ | `.claude/skills/signals/SKILL.md` | Alert/signal skill built |
| **realtime** | ⚠️ | NATS + WebSocket in onemind-web | Infrastructure exists; not documented or fully leveraged |
| **fastapi** | ❌ | — | Python REST framework — technology choice, not yet deployed anywhere |
| **mem fs** | ❌ | — | RAM-backed working memory for agents — architecture decision not made |

---

### Interfaces Layer (cross-cutting SO/01)

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **interfaces** | ✅ | `SO/02 Intelligence/openclaw/interfaces/` | 8 interface files: Telegram, Discord, Voice, Web, iOS/macOS app, Gateway API, OneMind app |
| **voice** | ✅ | `.claude/skills/vapi-voice-ai/SKILL.md` | Vapi voice AI — call handling, voice agents |
| **multichannel** | ⚠️ | Interface files exist; unified routing undocumented | Each channel has a file; cross-channel routing rules not written |
| **notifications** | ⚠️ | `onemind-web/doc/features/NOTIFICATION-SYSTEM.md` | Spec exists (partial); no implementation |
| **teams** | ❌ | — | Unclear if MS Teams integration or agent-team coordination — needs clarification before building |
| **multimodal** | ❌ | — | Vision/audio input — model capability exists, no routing or interface doc |

---

## Knowledge Layer (`_codex/`)

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **capture** | ✅ | `.claude/skills/codex-capture/SKILL.md` | Full CODEX "C" pipeline — URL/text/image/voice → structured note |
| **inbox** | ✅ | `SO/24 Inbox (Queue)/` | 8,197 files — exists and active, needs triage |
| **notes** | ✅ | All domain folders | Core vault content type |
| **references** | ✅ | `reference` entity type in `_codex/templates/` | One of 17 entity types |
| **tags** | ✅ | `_codex/CONVENTIONS.md` | Full tag taxonomy defined: type/, domain/, status/, area/ |
| **archive** | ⚠️ | `_archive/` folder + `_legacy/` folders | Exists; no formal archival policy or review cadence |
| **mem fs** | ❌ | — | Ephemeral working memory distinct from vault — not defined |

---

## HP — Holistic Performance (YOUR SELF)

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **habits** | ⚠️ | `HP/` domain | Domain exists; no dedicated habit tracking or streak system |
| **goals** | ⚠️ | Referenced in `CODEX-FRAMEWORK.md` | Entity type defined (`goal` template); no active goal tracking system |
| **metrics** | ⚠️ | `.claude/skills/life-score-intelligence/SKILL.md` | PI/WI/LI scores defined; data collection not automated |
| **insights** | ❌ | — | AI-generated pattern recognition from HP data — not built |
| **elevation** | ❌ | — | Personal growth arc, autonomy escalation — concept only |
| **learning** | ✅ | `_codex/memory/learning/` | Signal capture + feedback loops |
| **routines** | ⚠️ | Morning briefing + evening review skills | Skills exist; no routine tracker or protocol |
| **calendar** | ✅ | `.claude/skills/zoho/SKILL.md` | Zoho Calendar covered |
| **agenda** | ✅ | `.claude/skills/morning-briefing/SKILL.md` | Daily agenda via morning briefing |
| **tasks** | ✅ | `.claude/skills/todoist/SKILL.md` + `.claude/skills/clickup/SKILL.md` | Two task systems — Todoist + ClickUp |
| **projects** | ⚠️ | `project` template exists in `_codex/templates/` | Template defined; no active project tracking system |
| **principles** | ⚠️ | `_codex/system/` (beliefs file) | Beliefs in context file; no dedicated principles document |
| **vision** | ⚠️ | `SO/00 Framework/core-doctrine/VISION.md` | Framework vision exists; personal life vision not captured |

---

## GE — Generational Empire (YOUR EMPIRE)

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **missions** | ⚠️ | Referenced in `CODEX-FRAMEWORK.md` | Concept defined; no active mission tracking |
| **projects** | ⚠️ | Template + domain folder | Same as HP — template only |
| **paths** | ⚠️ | `SO/00 Framework/architecture/PATHS-SYSTEM.md` | System designed; not operationalized |
| **insights** | ❌ | — | AI-synthesized business/career intelligence — not built |

---

## Cross-Cutting (All Domains)

| Capability | Status | Codex Path | Notes |
|------------|--------|-----------|-------|
| **mcp** | ⚠️ | `openclaw.json` on VPS | context7 + github + codex-vault MCP servers active on all agents; zero codex doc |
| **tracing** | ❌ | — | Agent observability — no Langfuse, no OpenTelemetry, no audit trail UI |
| **loras** | ❌ | — | Local model fine-tuning — not a priority, placeholder only |

---

## Status Summary

| Status | Count | Examples |
|--------|-------|---------|
| ✅ Built & documented | 25 | agents, memory, skills, capture, interfaces, voice, hooks, workflows, tags, calendar, tasks |
| ⚠️ Exists but partial/undocumented | 32 | nats, rag, mcp, sessions, guardrails, hivemind, habits, goals, metrics, realtime, webhooks, bus routing |
| ❌ Missing entirely | 13 | tracing, loras, mem fs, activity feed, launch pad, insights, atlas, elevation, fastapi, multimodal, teams, reasoning, mcp docs |

---

## The 13 Gaps — Priority Order

| Priority | Capability | Why It Matters | Effort |
|----------|------------|---------------|--------|
| **1** | **tracing** | Blind to agent failures — can't debug what you can't see | Medium — Langfuse OSS or OTEL |
| **2** | **insights** | HP/GE data exists, nobody synthesizes it | Medium — n8n + agent query |
| **3** | **mcp docs** | 3 MCP servers running, no one knows what they do | Small — write one doc |
| **4** | **activity feed** | No way to know what agents did while you were away | Medium — aggregate NATS events |
| **5** | **launch pad** | Missing "start of day" surface — gap in daily operation | Small — design + build in onemind-web |
| **6** | **nats docs** | Core messaging backbone — completely undocumented | Small — write architecture doc |
| **7** | **rag pipeline** | Infrastructure ready, pipeline not wired | Medium — script + pgvector |
| **8** | **atlas** | No visual world map | Large — onemind-web feature |
| **9** | **reasoning** | qwen3 thinking mode available, no agent routing | Small — openclaw config |
| **08** | **multimodal** | Models support vision, nothing routes images to agents | Small — interface config |
| **11** | **elevation** | Growth arc needs definition before building | Small — define concept first |
| **12** | **fastapi** | Only needed when RAG pipeline requires a Python service | Small — build when needed |
| **13** | **loras** | Not a current priority | Large — defer |

---

## Related Files

| File | What |
|------|------|
| `AGENTS.md` | Universal agent instructions + full system overview |
| `CODEX-FRAMEWORK.md` | Method and taxonomy definition |
| `_codex/CONVENTIONS.md` | Naming, tagging, frontmatter rules |
| `_codex/INTERFACE.md` | Navigation entry point |
| `SO/05 Environments/README.md` | Entity registry — people, places, machines, assets |
| `SO/02 Intelligence/openclaw/interfaces/README.md` | All interface docs |
| `SO/02 Intelligence/openclaw/MULTI-USER.md` | Multi-user routing architecture |
| `_system/hooks/` | Lifecycle hooks mapped to CODEX stages |
| `onemind-web/doc/` | All onemind-web feature specs and architecture docs |
