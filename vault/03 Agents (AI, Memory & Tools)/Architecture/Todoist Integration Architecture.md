---
id: ui-02-20260111053925
title: "Todoist Integration Architecture"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Todoist Integration Architecture"
---

┌─────────────────────────────────────────────────────────────┐
│                    TODOIST STRUCTURE                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Projects (Mirror Pillars + Areas)                         │
│  ─────────────────────────────────                         │
│  📁 HP                                                      │
│     └── 📁 Fitness                                          │
│     └── 📁 Nutrition                                        │
│                                                             │
│  📁 LE                                                      │
│     └── 📁 Homestead                                        │
│     └── 📁 Family                                           │
│                                                             │
│  📁 GE                                                      │
│     └── 📁 Legacy AI Community                              │
│     └── 📁 Consulting                                       │
│                                                             │
│  📁 IT                                                      │
│     └── 📁 Infrastructure                                   │
│     └── 📁 AI Systems                                       │
│                                                             │
│  Labels (For Filtering)                                     │
│  ──────────────────────                                     │
│  @today  @this-week  @waiting  @delegated                   │
│  @legacy-ai  @hp-agent  @le-agent  @ge-agent  @it-agent     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   n8n Sync      │
                    │   Workflows     │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
      ┌─────────────┐ ┌───────────┐ ┌─────────────┐
      │ Push tasks  │ │Pull tasks │ │ Sync status │
      │ to Todoist  │ │from vault │ │  both ways  │
      └─────────────┘ └───────────┘ └─────────────┘
      
      


Command/Zeus/Today.md
─────────────────────
## Today's Focus

### HP
- [ ] Morning stack [[HP/Fitness/Workflows/Morning Protocol]]
- [ ] Log meals [[HP/Nutrition/Tasks/Daily Tracking]]

### LE  
- [ ] Chicken rounds [[LE/Homestead/Tasks/Daily Chores]]
- [ ] Family meeting 6pm [[LE/Family/Projects/Weekly Rhythm]]

### GE
- [ ] Community post [[GE/Legacy-AI-Community/Tasks/Content Calendar]]

### IT
- [ ] n8n backup [[IT/Infrastructure/Tasks/Maintenance]]


### How I (Legacy) Work With You

**You say:** "What's the status on the homestead?"

**I do:**

1. Read `LE/Homestead/_Index.md` for current status
2. Scan `LE/Homestead/Projects/` for active work
3. Check `LE/Homestead/Tasks/` for pending items
4. Pull relevant `Knowledge/` if you need SOPs

**Everything I need is in ONE place: `LE/Homestead/`**


┌─────────────────────────────────────────────────────────┐
│                      COMMAND                            │
│                   Zeus + Legacy                         │
│              (Orchestration Layer)                      │
└───────────────────────┬─────────────────────────────────┘
                        │
         ┌──────────────┼──────────────┬──────────────┐
         ▼              ▼              ▼              ▼
    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
    │HP Agent │    │LE Agent │    │GE Agent │    │Grid Agent │
    └────┬────┘    └────┬────┘    └────┬────┘    └────┬────┘
         │              │              │              │
    ┌────┴────┐    ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
    │ HP/     │    │ LE/     │    │ GE/     │    │ IT/     │
    │ ├Fitness│    │ ├Homstd │    │ ├Cmmnty │    │ ├Infra  │
    │ ├Nutrit │    │ ├Family │    │ ├Conslt │    │ ├DevOps │
    │ ├Sleep  │    │ ├Values │    │ ├Invest │    │ ├AI-Sys │
    │ └Mental │    │ └Eductn │    │ └R.Est  │    │ └Securt │
    └─────────┘    └─────────┘    └─────────┘    └─────────┘