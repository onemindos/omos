---
type: strategy-doc
topic: agent-model-tiers
status: active
updated: 2026-05-11
---
# One Mind Agent Model Strategy

Provider: **GitHub Copilot** (primary) — flat-rate, no token caps per agent.

---

## Fleet — 6 Agents

| Agent | ID | Node | Model | Role |
|-------|-----|------|-------|------|
| Legacy | main | omos-core | Sonnet 4.6 | Orchestrator — front door, routes all requests |
| Zenith | zenith | omos-core | Sonnet 4.6 | HP domain — body, mind, money, spirit, intelligence |
| Unity | unity | omos-core | Sonnet 4.6 | LE domain — family, home, land, estate, generations |
| Eternal | eternal | omos-core | Sonnet 4.6 | GE domain — empire, brand, capital, ventures, wealth |
| Sovereignty | sovereignty | omos-core | Sonnet 4.6 | SO domain — infra, security, automation, dev, research |
| Sentinel | sentinel | omos-sentinel | Haiku 4.5 | External watchdog + failover — monitors Core from outside |

---

## Model Rationale

**Sonnet 4.6** for all 5 domain agents on omos-core:
- Super-agents handle routed requests, scheduled operations, AND load skill packs
- Multi-step reasoning, tool calling, financial/security decisions require Sonnet
- No Haiku shortcuts for agents with domain authority

**Haiku 4.5** for Sentinel only:
- Sentinel does one thing: watch omos-core and activate failover
- Lightweight ping checks, status reads, alert dispatch — Haiku is sufficient
- Minimum footprint = maximum uptime

---

## Skills Architecture

The old 18 specialist agents are now SKILLS loaded on demand by domain agents.

| Skill Set | Loaded By | Skills |
|-----------|-----------|--------|
| HP skills | Zenith | spartan, oracle, identity, librarian, pulse |
| LE skills | Unity | haven, helix, eden, nomad |
| GE skills | Eternal | empire, empire-ops, creed, trinity, apex, finance |
| SO skills | Sovereignty | grid, guardian, guardian-ops, overwatch, eagleeye |

---

## Cost Rule

Sonnet for domain authority. Haiku for watchdog-only roles. Skill depth replaces agent count.
