---
id: so-01-a2a-readme-20260430
domain: SO
area: "01"
subdomain: "01 Intelligence / A2A - (AI Agent Frameworks)"
type: index
status: active
agent: Legacy
tags: [domain/so, area/01, folder/a2a, type/index, status/active]
ai_summary: "OneMind multi-agent system — design history, Letta architecture, voice stack, and live Letta strategy."
updated: 2026-04-30
---
# A2A — AI Agent Frameworks

> *The OneMind multi-agent design system — from initial roster design through Letta architecture to the live OpenClaw-backed voice stack.*

**Domain**: SO → 01 Intelligence → **A2A** ← you are here

---

## Root Files

| File | What It Is |
|------|-----------|
| `LETTA-STRATEGY.md` | **Current strategy.** Letta's role as OpenClaw's backend brain — archival memory service, background agents, sleep-time compute, multi-model experimentation. Also covers Letta Code as the dev tool for building OpenClaw. Includes `letta-memory` skill code and codex upload script. |
| `Training Proactive and Personalized LLM Agents.md` | arXiv clip (2511.02208) — failed import, body is empty. Pending re-import or deletion. |

---

## Subfolders

### `agent-roster/` — Phase 1: The 12-Agent Design (Nov 2025)
Original OneMind agent architecture — Legacy as master orchestrator over 12 specialized agents across four pillars (HP / GE / LE / IT). Clean domain separation, autonomy levels, coordination patterns. **Superseded by the live 9-agent OpenClaw roster**, but the pillar structure and role definitions remain canonical reference.

| File | Content |
|------|---------|
| `00 - Agent Roster Overview.md` | Full system architecture, coordination patterns, autonomy framework, implementation roadmap |
| `00 - Legacy Agent.md` | Legacy master orchestrator — memory access, morning/evening routines, decision escalation framework |
| `01–11` numbered profiles | One file per agent: HP (Intelligence, Performance, Finance), GE (Founder Assistant, Capital Systems), LE (Home Command, Homestead Ops, Family Wealth), IT (Infrastructure, Robotics) |
| `Agent Naming Convention.md` | Domain prefix system (HP/GE/LE/IT), quick reference matrix, autonomy key |
| `LEGACY AI - REFINED AGENT ROSTER.md` | 551-line detailed spec — skills, memory blocks, tools, and autonomous behaviors per agent |

---

### `letta-architecture/` — Phase 2: Letta Cloud Build Plan (Jan 2026)
Full technical implementation of the four-pillar system using Letta Cloud — 41 agents total (Mother AI → 4 conductors → 37 specialists), OMCS coordinate system, N8N cron workflows, human-in-the-loop approval flows, MCP server config. **Built before OpenClaw became the live runtime.**

| File | Content |
|------|---------|
| `ONE-MIND-LETTA-AGENT-ARCHITECTURE-v2.md` | 1,312-line master spec — Command Layer (MotherAI + Protocol Guardian + Approval Manager + Memory Curator), all 4 pillar conductors + 37 specialists, shared memory blocks, Todoist project structure, N8N cron jobs, human-in-loop approval flow, tool patterns |
| `Letta AI Agents Guide.md` | Letta primer — stateful agents, core/archival/recall memory, Python SDK quickstart |

---

### `voice/` — Phase 3+4: Voice Architecture (2026)
Everything related to giving OneMind a voice — from research through proposed architecture to implementation plans.

| File | Content |
|------|---------|
| `OPENCLAW-VOICE-RESEARCH.md` | Deep research doc — VisionClaw (Gemini Live pattern), openclaw-nerve web UI (752 stars), 3 Flutter voice approaches (Talk mode / Gemini Live / Hybrid), full VoiceService state machine blueprint, universal device matrix, Quick Wins deployable today |
| `VOICE-FIRST-SWARM-ARCHITECTURE.md` | Proposed swarm design — Legacy as swarm consciousness, MotherAgent + AgentHandoff + Teams already built but disconnected. 5-phase wiring plan. "The power is already built. The plumbing is disconnected." |
| `VOICE_MULTIMODAL_PLAN.md` | LiveKit-based voice/video/screen-share plan (Jan 2026) — Silero VAD + Deepgram STT + Cartesia TTS pipeline. ~80% of code already exists. Predates OpenClaw voice pivot; partially superseded by `OPENCLAW-VOICE-RESEARCH.md`. |

---

### `integrations/` — Data Layer

| File | Content |
|------|---------|
| `BRIDGE-V2-ARCHITECTURE.md` | Bridge v2 — simplified LifeOS ↔ Learning Machine connector. Entity wrapper removed; Bridge now does analytics only (PI/LI/WI intelligence scores, pattern detection, day predictions). Full PostgreSQL schema, API endpoints, data flow with triggers, migration steps. ~47% code reduction from v1. |

---

## Design Evolution

```
Nov 2025  →  agent-roster/    12-agent HP/GE/LE/IT pillar design
Jan 2026  →  letta-architecture/  Full Letta Cloud build (41 agents, OMCS coordinates)
Jan 2026  →  integrations/    Bridge v2 data layer
2026      →  voice/           Voice research + swarm + LiveKit plans
Now       →  LETTA-STRATEGY.md  OpenClaw is live. Letta = backend brain.
```
