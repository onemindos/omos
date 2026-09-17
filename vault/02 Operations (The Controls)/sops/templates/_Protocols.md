---
id: ui-04-20260111053926
title: "Protocols — Reusable Plays Registry"
type: procedure
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/04]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Protocols — Reusable Plays Registry"
---

# Protocols — Reusable Plays Registry

> *"Build once, run forever."*

---

## Purpose

Protocols houses the **registry of reusable procedures**. Every repeatable process gets documented here, versioned, and made available for execution across the Codex.

---

## Protocol Types

| Type | Code | Purpose | Lifespan |
|------|------|---------|----------|
| **Protocol** | PRT | Reusable step-by-step procedures | Forever (versioned) |
| **Mission** | MIS | Time-boxed projects with end dates | Temporary (archived) |
| **Routine** | RTN | Recurring scheduled protocol runs | Ongoing (cadenced) |
| **Run** | RUN | Single execution instance | One-time (logged) |

---

## Protocol Registry Structure

```
Protocols/
├── GE/
│   ├── PRT-GE-Marketing-ContentPublish-001.md
│   ├── PRT-GE-Sales-DemoCall-001.md
│   └── PRT-GE-Legal-ContractReview-001.md
├── HP/
│   ├── PRT-HP-Performance-MorningStack-001.md
│   ├── PRT-HP-Training-WorkoutLog-001.md
│   └── PRT-HP-Finance-MonthlyReview-001.md
├── LE/
│   ├── PRT-LE-Family-SundayRitual-001.md
│   ├── PRT-LE-Homestead-GardenMaint-001.md
│   └── PRT-LE-Wealth-QuarterlyReview-001.md
└── IT/
    ├── PRT-IT-Cybersecurity-BackupRestore-001.md
    ├── PRT-IT-Automation-WorkflowDeploy-001.md
    └── PRT-IT-HomeLab-ServerMaint-001.md
```

---

## Protocol Template

```markdown
# [PRT-ID] Protocol Name

**Version:** 1.0  
**Owner:** [Name]  
**Overlay Kits:** [SEC] [OPS] [TRAIN]  
**Last Updated:** YYYY-MM-DD

---

## Purpose
[One sentence: what outcome does this protocol produce?]

## Trigger
[What event starts this protocol?]

## Prerequisites
- [ ] Prerequisite 1
- [ ] Prerequisite 2

## Steps
1. **Step 1 Name**
   - Detail
   - Detail

2. **Step 2 Name**
   - Detail
   - Detail

3. **Step 3 Name**
   - Detail
   - Detail

## Outputs
- [ ] Output 1
- [ ] Output 2

## Signals (Metrics)
- Lead: [what predicts success]
- Lag: [what confirms success]

## Rollback
[What to do if this fails]

## Related
- [[Related Protocol 1]]
- [[Related Protocol 2]]

---

## Changelog
| Version | Date | Change |
|---------|------|--------|
| 1.0 | YYYY-MM-DD | Initial version |
```

---

## Mission Template

```markdown
# [MIS-ID] Mission Name

**Status:** 🟡 In Progress | 🟢 Complete | 🔴 Blocked  
**Owner:** [Name]  
**Timeline:** Q1 2025 (Start: YYYY-MM-DD → End: YYYY-MM-DD)  
**Overlay Kits:** [OPS] [FIN]

---

## Objective
[One sentence: what does done look like?]

## Key Results
1. [ ] KR1: [Measurable outcome]
2. [ ] KR2: [Measurable outcome]
3. [ ] KR3: [Measurable outcome]

## Milestones
| Week | Milestone | Status |
|------|-----------|--------|
| 1-2 | [Milestone 1] | ⬜ |
| 3-4 | [Milestone 2] | ⬜ |
| 5-6 | [Milestone 3] | ⬜ |

## Resources Required
- Time: [hours/week]
- Budget: [$amount]
- Tools: [list]
- People: [who]

## Risks
| Risk | Impact | Mitigation |
|------|--------|------------|
| [Risk 1] | High | [Plan] |

## Related Protocols
- [[PRT-xxx]]

---

## Log
| Date | Update |
|------|--------|
| YYYY-MM-DD | [What happened] |
```

---

## Routine Template

```markdown
# [RTN-ID] Routine Name

**Cadence:** Daily | Weekly | Monthly | Quarterly  
**Owner:** [Name]  
**Protocol:** [[PRT-xxx]]  
**Overlay Kits:** [SYS] [OPS]

---

## Schedule
- **Day/Time:** [When it runs]
- **Duration:** [How long]
- **Reminder:** [How you're prompted]

## Checklist
- [ ] Step from protocol
- [ ] Step from protocol
- [ ] Step from protocol

## Last 5 Runs
| Date | Status | Notes |
|------|--------|-------|
| YYYY-MM-DD | ✅ | |
| YYYY-MM-DD | ✅ | |
| YYYY-MM-DD | ⚠️ | [Issue] |
```

---

## Run Log Template

```markdown
# [RUN-ID] Run Log

**Protocol:** [[PRT-xxx]]  
**Date:** YYYY-MM-DD  
**Executor:** [Name]

---

## Execution
- Start: HH:MM
- End: HH:MM
- Duration: [X min]

## Results
- [ ] Output 1: [Achieved/Not]
- [ ] Output 2: [Achieved/Not]

## Issues
[Any blockers or deviations]

## Learnings
[What to improve in protocol]

## Evidence
[Links to proof/artifacts]
```

---

## Active Protocols Index

| ID | Name | Pillar | Kits | Owner |
|----|------|--------|------|-------|
| PRT-IT-Cybersecurity-BackupRestore-001 | Backup & Restore | IT | SEC,OPS | Zeus |
| PRT-HP-Performance-MorningStack-001 | Morning Stack | HP | SYS,OPS | Zeus |
| PRT-LE-Family-SundayRitual-001 | Sunday Family Ritual | LE | COMM,SYS | Zeus |
| PRT-GE-Marketing-ContentPublish-001 | Content Publishing | GE | COMM,OPS | Zeus |

---

## Navigation

- Parent: [[01 Zeus/_Dashboard]]
- Related: [[00 Framework/ONEMIND-CODEX]], [[01 Zeus/Atlas/_Atlas]], [[04 Metrics/_Metrics]]

---

**Version:** 1.1
**Last Updated:** 2025-12-26
