# Task Platform Research
*2026-05-27*

## Context

Evaluating a self-hosted replacement for Todoist as the universal task layer for all 6 agents (Legacy, Zenith, Unity, Eternal, Sovereignty, Sentinel) coordinating house, office, barn, and personal domains.

**Hard requirements:**
- Self-hosted (own the infrastructure)
- Native iOS app with push notifications + offline access
- REST API with API key auth (agent access)
- Webhooks (event-driven, not polling)
- Calendar view
- Reminders
- Multi-user / service accounts per agent

---

## Full Platform Scorecard

| Platform | REST API | Webhooks | iOS + Push | Calendar | Lightweight | Verdict |
|---|---|---|---|---|---|---|
| **Plane CE** | YES (180+ endpoints) | PARTIAL (5 types) | PARTIAL (no confirmed offline) | YES | NO (12 containers, 6-8 GB RAM) | Heavy, CE API not contractually guaranteed |
| **OpenProject CE** | YES (full, CE guaranteed) | YES (HMAC-signed) | PARTIAL (beta app, no push yet) | YES | MEDIUM (9 containers / 2 GB AIO) | Best API, heavy compose |
| **NocoDB** | YES (excellent, auto-gen) | YES (conditional) | NO | YES | YES (1 container, ~350 MB) | No iOS app |
| **Planka** | YES (v2.0+) | PARTIAL (reliability issues) | NO | NO (planned) | YES (2 containers, ~300 MB) | No iOS, no calendar |
| **Leantime** | NO (JSON-RPC only) | NO | NO | YES | YES | 3 hard fails |
| **Baserow** | YES | YES | NO | UNKNOWN | MEDIUM (~700 MB) | No iOS app |
| **Super Productivity** | PARTIAL (localhost-only) | NO | YES (best iOS+offline) | NO | YES | Wrong architecture for multi-agent |
| **AppFlowy** | NO | NO | YES (iOS+offline) | YES | MEDIUM | No public REST API |
| **Huly** | PARTIAL (TypeScript SDK) | NO | NO | YES | NO (8 GB RAM recommended) | No mobile app |
| **Vikunja** | — | — | — | — | — | ELIMINATED — poor experience |
| **WeKan** | PARTIAL | NO | NO | NO | YES | Maintainer warns of data corruption |
| **Taiga** | — | — | — | — | — | DEAD — team pivoted to Tenzu |
| **Focalboard** | — | — | — | — | — | ABANDONED — repo unmaintained |
| **Gitea/Forgejo** | YES (best API) | YES (best webhooks) | NO | NO | YES | No iOS app, wrong tool type |
| **Kanboard** | PARTIAL (JSON-RPC) | PARTIAL | NO | YES | YES | No iOS app, maintenance mode |
| **Nextcloud Deck** | YES | NO | PARTIAL | NO | NO (requires Nextcloud) | Performance warnings from maintainer |
| **Trilium** | PARTIAL (notes-only) | NO | NO | NO | YES | Notes tool, not task manager |
| **Logseq** | NO | NO | PARTIAL (alpha) | NO | NO | PKM tool, wrong category |
| **Outline** | PARTIAL (docs-only) | PARTIAL | NO | NO | YES | Wiki tool, wrong category |
| **Taskwarrior** | NO | NO | NO | NO | YES | CLI-only |

---

## The Core Finding

**The full requirement set does not exist in a single lightweight self-hosted product.**

Every genuinely lightweight option (< 1 GB RAM) fails on the native iOS app requirement. The structural reason: maintaining an App Store app requires Apple Developer infrastructure that small open-source teams can't sustain. Only tools with dedicated teams (Plane, OpenProject) have iOS apps — and both are heavy.

---

## Viable Paths

### Path A — NocoDB + ntfy.sh (Recommended for agent-first use)

```
VPS omos-core
├── NocoDB        (~350 MB) — REST API, webhooks, calendar, task schema
└── ntfy.sh       (~10 MB)  — iOS push notification bridge
```

- Agents hit NocoDB REST API for task CRUD
- NocoDB conditional webhooks → agent orchestration endpoints
- ntfy.sh has real native iOS app with APNs push — agents POST to ntfy topics, iPhone receives push
- **Gap:** No native iOS task management UI — tasks managed via NocoDB web interface in Safari or iOS Shortcuts

**Total: 2 containers, ~360 MB RAM**

---

### Path B — Planka + ntfy.sh

```
VPS omos-core
├── Planka + Postgres  (~300 MB) — Kanban UI, REST API, API keys
└── ntfy.sh            (~10 MB)  — iOS push notifications
```

- Better human-facing task UI than NocoDB (Trello-style board)
- Webhooks exist but have documented reliability issues (v2.0 regressions)
- No calendar view (on roadmap, not shipped)

**Total: 3 containers, ~310 MB RAM**

---

### Path C — OpenProject CE (AIO container)

```
VPS omos-core
└── OpenProject AIO  (~1.5-2 GB) — full feature set, single container
```

- Checks every requirement including iOS app (beta), CalDAV, confirmed REST API in CE
- `.deb` package install on Ubuntu = zero Docker overhead, `systemctl` management
- iOS app is beta — push notifications listed as upcoming
- HAL+JSON API format adds agent coding complexity

**Total: 1 container or 0 (bare metal), ~2 GB RAM**

---

## Habit Tracking

Keep separate from the task layer. Different data model, different agent logic. Habits are streaks/patterns; tasks are discrete completions. Options:
- Habitica (self-hosted, has API)
- Custom SQLite on VPS (cleanest for direct agent writes)

---

## Agent Architecture Mapping

```
Mac (local, always-on)
├── Hermes  → Legacy (orchestrator) ─────────────────┐
└── OpenClaw → Zenith / Unity / Eternal /             │  HTTP over Tailscale
               Sovereignty / Sentinel ────────────────┤
                                                      ↓
VPS omos-core (100.98.25.10)
└── Task Platform
    ├── API tokens: one per agent (scoped)
    ├── Projects: House / Office / Barn / SO / HP / LE / GE
    └── Webhooks → Mac trigger endpoints
```

---

## Decision Pending

**Next step:** Decide which matters more — RAM footprint or native iOS task management experience.

- Choose **NocoDB + ntfy** if agent API reliability and low footprint are the priority
- Choose **OpenProject AIO** if a native iOS task experience is required
