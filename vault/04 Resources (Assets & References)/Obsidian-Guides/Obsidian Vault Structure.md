---
id: ui-06-20260111053926
title: "Obsidian Vault Structure"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/06]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Obsidian Vault Structure"
---


Legacy/
│
├── Personal/                    # Your pillars
│   ├── Holistic Performance/
│   ├── Generational Endeavor/
│   ├── Legacy Evolution/
│   └── Innovative Technology/
│       └── Legacy AI Ops/       # Operational notes about running Legacy AI
│           ├── Daily Monitoring.md
│           ├── Incident Log.md
│           └── Performance Metrics.md
│
├── Legacy AI/                   # Legacy AI's brain structure
│   │
│   ├── 00 - Codex/             # Foundation (her constitution)
│   │   ├── README.md           # What the Codex is
│   │   ├── Core Principles.md
│   │   ├── Decision Framework.md
│   │   └── Evolution Log.md
│   │
│   ├── 01 - Architecture/      # System design (how she's built)
│   │   ├── System Overview.md
│   │   ├── Memory Architecture.md
│   │   ├── Agent Framework.md
│   │   ├── Hybrid Modes.md
│   │   ├── Tool Integration.md
│   │   └── Data Flow.md
│   │
│   ├── 02 - Development/       # Active work (what's being built)
│   │   ├── Active Sprint.md    # Current focus
│   │   ├── Phase 2A - Production.md
│   │   ├── Phase 2B - Edge Planning.md
│   │   ├── Phase 3 - Multi-Agent.md
│   │   ├── Technical Decisions.md
│   │   └── Implementation Notes/
│   │
│   ├── 03 - Operations/        # Running the system (how to deploy/maintain)
│   │   ├── Deployment Guide.md
│   │   ├── Troubleshooting.md
│   │   ├── Backup & Recovery.md
│   │   ├── Monitoring Setup.md
│   │   └── Runbooks/
│   │
│   ├── 04 - Research/          # Learning & exploration
│   │   ├── Framework Evaluations/
│   │   ├── Technology Comparisons/
│   │   ├── Experiments/
│   │   └── Papers & Articles/
│   │
│   └── 05 - Knowledge Base/    # Accumulated wisdom
│       ├── Learnings.md        # What worked/didn't work
│       ├── Patterns.md         # Reusable patterns
│       ├── API References/
│       └── Best Practices/
│
└── _Templates/                 # Obsidian templates
    ├── Daily Note.md
    ├── Project Plan.md
    ├── Architecture Doc.md
    └── Decision Log.md
    
    
    
    # Separate Repos
/legacy-vault/                    # Private Obsidian vault
  ├── .git (separate)
  ├── /00-inbox/
  ├── /01-research/
  ├── /02-architecture/          # Architecture decisions
  ├── /03-operations/
  ├── /04-knowledge/
  └── .obsidian/

/Legacy-AI-Hybrid/               # Your GitHub repo (public)
  ├── .git (separate)
  ├── /docs/                     # PUBLIC technical docs only
  ├── /src/
  ├── /infrastructure/
  └── README.md