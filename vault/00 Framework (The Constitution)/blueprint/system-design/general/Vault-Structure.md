---
id: arch-00-vault-20260430
title: "OneMind Codex — Vault Structure"
type: architecture
domain: SO
area: "00 Framework"
agent: LEGACY
status: active
created: 2026-04-30
modified: 2026-04-30
owner: zeus
tags: [type/architecture, status/active, domain/so, system/codex, layer/structure]
related:
  - system-primitives.md
  - Cross-Pillar-Routing.md
  - OPENCLAW-FRAMEWORK-STRUCTURES.md
todoist_project_id: "6g6P9mmQ3Pf36qjG"
ai_summary: "Canonical reference for the OneMind Codex folder structure — four domains (SO/HP/LE/GE), numbering convention (00-99), hub-node pattern, and universal archive."
---

# OneMind Codex — Vault Structure

**Version**: 3.0 (SO/HP/LE/GE model)
**Last Updated**: 2026-04-30
**Purpose**: Canonical reference for how the Codex is organized

---

## Overview

The OneMind Codex is a Git-backed Obsidian vault organized into **four domains** covering the full range of Zeus's life, systems, and legacy. Every file belongs to one domain. Every domain has a number range.

```
Domain         Range    Focus
──────────────────────────────────────────────────────────
SO  Sovereign Operations  00-24   The Brain — systems, AI, doctrine, infra
HP  Holistic Performance  25-49   Self Path — health, mind, skills, finance
LE  Legacy Evolution      50-74   Legacy Path — family, home, heritage, estate
GE  Generational Empire 75-99   Wealth Path — business, revenue, ventures
```

---

## Top-Level Structure

```
onemind-codex/
│
├── 00-24 SO (Sovereign Operations)/   ← The operating layer
│   ├── 00 Framework (Doctrine)/       ← Rules, architecture, standards
│   ├── 01 Intelligence (AI & Memory)/ ← Agent configs, OpenClaw, memory
│   ├── 02 Infrastructure/             ← VPS, networking, DevOps
│   ├── 03 Command/                    ← Daily ops, inbox, agenda
│   ├── 04 Security/                   ← Guardrails, access, certs
│   ├── 05 Development/                ← Apps, code, Flutter, APIs
│   └── 06 Research/                   ← Experiments, ML, exploration
│
├── 25-49 HP (Holistic Performance)/   ← Self optimization
│   ├── [25-29] Physical/              ← Health, fitness, nutrition
│   ├── [30-34] Mental/                ← Mind, habits, journaling
│   ├── [35-39] Finance/               ← Personal money, credit, investing
│   ├── [40-44] Skills/                ← Learning, certifications, knowledge
│   └── [45-49] Identity/              ← Purpose, values, relationships
│
├── 50-74 LE (Legacy Evolution)/       ← Family and heritage
│   ├── [50-54] Family/                ← Household, children, relationships
│   ├── [55-59] Home/                  ← Property, homestead, environment
│   ├── [60-64] Wealth/                ← Family estate, investments, trusts
│   ├── [65-69] Heritage/              ← Culture, traditions, memories
│   └── [70-74] Legal/                 ← Estate planning, family legal
│
├── 75-99 GE (Generational Empire)/  ← Business and wealth creation
│   ├── [75-79] Ventures/              ← Business entities, products
│   ├── [80-84] Finance/               ← Business capital, banking, credit
│   ├── [85-89] Operations/            ← Processes, SOPs, team
│   ├── [90-94] Growth/                ← Marketing, sales, community
│   └── [95-99] Strategy/              ← Long-term planning, positioning
│
└── .archive/                          ← Universal archive (hidden from graph)
    ├── 00-framework/                  ← Archived SO/Framework docs
    ├── 01-intelligence/               ← Archived agent/AI docs
    └── [domain]/                      ← Archived docs from any domain
```

---

## SO Domain Detail (00-24)

The Sovereign Operations domain is the brain of the system — it governs all other domains.

```
00-24 SO (Sovereign Operations)/
│
├── 00 Framework (Doctrine)/
│   ├── architecture/          ← System design documents (this file lives here)
│   ├── core-doctrine/         ← Vision, principles, constitution
│   └── reference/             ← Standards, CODEX-STANDARD, tag taxonomy
│
├── 01 Intelligence (AI, Memory & Tools)/
│   ├── agents/                ← OpenClaw agent configs (LEGACY, Spartan, etc.)
│   ├── openclaw/              ← OpenClaw platform docs and specs
│   ├── memory/                ← Memory models, Letta blocks
│   └── tools/                 ← MCP tools, integrations
│
├── 02 Infrastructure/
│   ├── vps/                   ← omos-core, omos-sentinel configs
│   ├── networking/            ← Tailscale, DNS, firewall
│   └── docker/                ← Compose files, container configs
│
├── 03 Command/
│   ├── Inbox/                 ← Universal capture landing zone
│   ├── Agenda/                ← Daily/weekly ops
│   └── Metrics/               ← KPIs, dashboards
│
├── 04 Security/
│   ├── access/                ← SSH keys, secrets management
│   ├── guardrails/            ← Agent autonomy limits (L0-L3)
│   └── compliance/            ← Audit logs, policies
│
├── 05 Development/
│   ├── onemind-app/           ← Flutter mobile/web app
│   ├── backend/               ← API, services
│   └── automation/            ← N8N workflows
│
└── 06 Research/
    ├── experiments/           ← Active R&D
    └── learning/              ← Technical reading notes
```

---

## Numbering Convention

Folders use a **two-digit prefix** that places them within their domain range:

```
SO:  00-24  (e.g., 00 Framework, 01 Intelligence, 02 Infrastructure...)
HP:  25-49  (e.g., 25 Physical, 26 Nutrition, 30 Habits...)
LE:  50-74  (e.g., 50 Family, 55 Homestead, 60 Wealth...)
GE:  75-99  (e.g., 75 Ventures, 80 Finance, 90 Growth...)
```

Files use the ID format: `[domain]-[area]-[YYYYMMDD]-[slug]`
Example: `arch-00-20260430-vault-structure`

---

## Hub Node Pattern

Every folder has a **README.md** that serves as the hub node:

```markdown
---
id: hub-[folder]
type: hub
---

# [Folder Name]

> One-line purpose statement.

## What's Here
- Brief description of subfolders/files

## Key Documents
- [[key-doc-1]]
- [[key-doc-2]]

## See Also
- [[related-doc-in-another-domain]]
```

---

## Universal Archive

All retired/superseded content moves to a **single universal archive**:

```
onemind-codex/.archive/
```

Rules:
- Dot-prefix hides it from Obsidian graph by default
- Organized by source domain subfolder
- Never delete — always move here first
- No distributed `archive/` subfolders inside domain folders

---

## Tag Taxonomy

```
domain/so   domain/hp   domain/le   domain/ge
type/architecture   type/concept   type/procedure   type/reference
type/hub    type/log    type/template
status/active   status/staged   status/archived
system/openclaw   system/hive   system/paths   system/codex
agent/legacy   agent/spartan   agent/oracle   agent/eden   agent/empire
layer/foundation   layer/intelligence   layer/execution   layer/command
```

---

## VPS Infrastructure

| Server | Tailscale IP | Role |
|--------|-------------|------|
| omos-core | 100.123.126.34 | Primary — OpenClaw, NATS, LEGACY agent |
| omos-sentinel | 100.66.235.38 | Secondary — monitoring, failover |

Codex is synced to VPS at `/home/onemind/onemind-os/codex/`

---

## Key Files Reference

| File | Location | Purpose |
|------|----------|---------|
| CODEX-STANDARD.md | `00 Framework/reference/` | The law — frontmatter spec, file types, tags |
| system-primitives.md | `00 Framework/architecture/` | 55 canonical system concepts |
| OPENCLAW-FRAMEWORK-STRUCTURES.md | `00 Framework/architecture/` | Agent roster, autonomy levels, NATS map |
| HIVE-CAPTURE-SYSTEM.md | `00 Framework/architecture/` | 14-source capture pipeline |
| FRAMEWORK-INTEGRATIONS.md | `00 Framework/architecture/` | Why this structure (Maslow, PARA, etc.) |

---

## See Also

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/blueprint/system-design/Cross-Pillar-Routing]] — How to link content across domains
- [[system-primitives.md]] — The 55 primitives that define the OS
- [[OPENCLAW-FRAMEWORK-STRUCTURES.md]] — Agent + infrastructure doctrine

---

*v3.0 — 2026-04-30 | Replaces v1 4-pillar (GE/HP/LE/IT) structure with current SO/HP/LE/GE (00-99) model*
