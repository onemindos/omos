---
id: upgrade-10-reports
title: "Upgrade 10: Report Generation"
type: blueprint
status: deployed
created: 2026-03-27
deployed: 2026-03-27
owner: zeus
tags: [type/blueprint, status/deployed, domain/00-framework]
destination: "_codex/skills/report/"
---

# Upgrade 10: Report Generation

> ✅ **DEPLOYED** to `_codex/skills/report/` — commit `28356c78` (2026-03-27)

> **Source**: PAI's TELOS WriteReport + CreateNarrativePoints workflows
> **Destination**: `_codex/skills/report/`
> **Why**: OneMind has daily/weekly/monthly reviews (in reflect.md) but no structured report generation. PAI generates McKinsey-quality reports and dashboards from personal context. We adapt this for quarterly life reviews, domain health reports, and family legacy reports.

---

## Report Types

### 1. Quarterly Life Review
**Frequency**: Every 3 months
**Scope**: All 4 domain groups
**Format**: Markdown report

```markdown
# Quarterly Life Review — Q1 2026

## Executive Summary
{AI-generated 3-sentence overview of the quarter}

## Domain Group Scores
| Group | Health | Trend | Key Win | Key Gap |
|-------|--------|-------|---------|---------|
| SO (Operations) | 🟢 85% | ↑ | Infrastructure stable | Security audit overdue |
| HP (Performance) | 🟡 65% | → | Consistent workouts | Sleep still inconsistent |
| LE (Legacy) | 🟢 80% | ↑ | Family meetings weekly | Estate planning stalled |
| GE (Endeavor) | 🟡 60% | ↑ | Revenue growing | Product not shipped |

## Highlights
{Top 5 accomplishments}

## Lessons Learned
{Extracted from memory/learning/signals.jsonl}

## Next Quarter Priorities
{Based on active-goals.md + gap analysis}
```

### 2. Domain Health Report
**Frequency**: Monthly or on-demand
**Scope**: Single domain or domain group
**Format**: Markdown with metrics

```markdown
# Domain Health: 03 Agents (AI, Memory & Tools)

## Status: 🟢 Healthy
## Files: 47 | Last Updated: 2026-03-27 | Active Projects: 3

## Recent Activity
{Last 30 days of changes in this domain}

## Knowledge Gaps
{Files with status/gap tags or stubs}

## Connections
{Most-linked-to files, cross-domain references}

## Recommendations
{AI-suggested improvements}
```

### 3. Family Legacy Report
**Frequency**: Annually
**Scope**: LE domain group (50-74) + relevant GE domains
**Format**: Comprehensive markdown document

```markdown
# Family Legacy Report — 2026

## Family Identity
{From _codex/system/identity.md + LE domain files}

## Milestones This Year
{Extracted from daily/weekly logs in LE domains}

## Knowledge Captured
{Count and summary of LE domain files created}

## Financial Health
{From GE domain 82 Investments + 30 Finance}

## Systems Status
{Family meeting cadence, shared calendars, SOPs}

## Legacy Score
{Composite metric from LE domains}

## Next Year Goals
{Projected from active-goals.md}
```

### 4. Narrative Points (McKinsey-Style)
**Frequency**: On-demand (for pitches, presentations, content)
**Scope**: Any topic
**Format**: Structured narrative bullets

PAI generates McKinsey-style narrative points in an 8-section arc:
1. **Opening** — Hook the audience
2. **Bad News** — The problem (status quo pain)
3. **Evidence** — Data supporting the problem
4. **Stakes** — What happens if nothing changes
5. **Pivot** — The solution (OneMind)
6. **Requirements** — What it takes
7. **Vision** — The future with the solution
8. **Close** — Call to action

## Skill Pack Structure

```
_codex/skills/report/
├── README.md           ← Overview
├── SKILL.md            ← How agents generate reports
├── INSTALL.md          ← Self-installation
├── VERIFY.md           ← Verification
└── templates/
    ├── quarterly-review.md
    ├── domain-health.md
    ├── family-legacy.md
    └── narrative-points.md
```

## Data Sources

Reports pull from:
- `_codex/system/` — Identity, goals, mission
- `_codex/memory/` — Signals, patterns, work history
- Domain folders — File counts, recent changes, status tags
- Git history — Commit frequency, contributor activity
- Frontmatter — Status fields, tags, dates

## How to Deploy

1. Create `_codex/skills/report/` as a skill pack
2. Copy report templates
3. Agents can generate reports on command: "Generate my quarterly review"
4. Schedule quarterly reports via automation (domain 09)

## PAI Equivalent

PAI's TELOS pack generates full Next.js dashboards (shadcn/ui + Aceternity UI) and professional reports (custom typography, Tokyo Night theme). We adapt the concept to markdown reports (Obsidian-native) instead of web apps. When we build the web platform, these can graduate to interactive dashboards.
