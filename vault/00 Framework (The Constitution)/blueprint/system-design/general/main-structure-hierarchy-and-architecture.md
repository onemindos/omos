---
id: main-structure-hierarchy-and-architecture
title: "Main Structure — Hierarchy & Final Architecture"
type: note
source: apple-notes
imported: 2026-03-07
modified: 2026-04-30
status: staged
tags:
  - apple-notes-import
  - onemind-framework
  - architecture
  - vault-structure
  - openclaw
  - agents
  - cognitive-framework
---

# Main Structure — Hierarchy & Final Architecture

## The Hierarchy (Final Answer)

```
┌─────────────────────────────────────────────────────┐
│                     ONE MIND                        │
│          (The Meta-Cognitive Layer / Life OS)        │
│                                                     │
│  "The unified system where human and AI think       │
│   as one."                                          │
│                                                     │
│  This is:                                           │
│  • The vault name                                   │
│  • The framework name                               │
│  • The product you teach to community               │
│  • The thing that competes with PARA                │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │
│  │    ZEUS     │  │  LEGACY AI  │  │LEGACY CODEX │ │
│  │             │  │             │  │             │ │
│  │   Human     │  │ Artificial  │  │   Shared    │ │
│  │  Cognition  │  │  Cognition  │  │  Knowledge  │ │
│  │             │  │             │  │             │ │
│  │ • Decisions │  │ • Agents    │  │ • HP Path   │ │
│  │ • Values    │  │ • Tools     │  │ • LE Path   │ │
│  │ • Intuition │  │ • Automation│  │ • GE Path   │ │
│  │ • Direction │  │ • Memory    │  │ • IT Path   │ │
│  │             │  │             │  │             │ │
│  │The Commander│  │The Enhancer │  │ The Library │ │
│  └─────────────┘  └─────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────┘
```

> **Legacy AI is your enhancement, your suit, your Jarvis.** But it's not the TOP. The top is **One Mind** — the unified system where YOU + LEGACY AI operate together using the CODEX.

---

## Final Architecture Summary

```
┌═══════════════════════════════════════════════════════┐
│                   ONE MIND SYSTEM                     │
├═══════════════════════════════════════════════════════┤
│                                                       │
│  ┌─────────────────────────────────────────────────┐  │
│  │              INTELLIGENCE LAYER                 │  │
│  │                                                 │  │
│  │         ┌─────────────────────┐                 │  │
│  │         │     LEGACY AI       │                 │  │
│  │         │   (Orchestrator)    │                 │  │
│  │         │                     │                 │  │
│  │         │   Memory Blocks:    │                 │  │
│  │         │   • Core Identity   │                 │  │
│  │         │   • HP Context      │                 │  │
│  │         │   • LE Context      │                 │  │
│  │         │   • GE Context      │                 │  │
│  │         │   • IT Context      │                 │  │
│  │         │   • Agent Registry  │                 │  │
│  │         └────────┬────────────┘                 │  │
│  │                  │                              │  │
│  │        ┌─────────┼─────────┐                    │  │
│  │        ▼         ▼         ▼                    │  │
│  │  ┌──────────┐┌──────────┐┌──────────┐          │  │
│  │  │ Spartan  ││  Eden    ││  Grid    │          │  │
│  │  │(HP Agent)││(LE Agent)││(SO Infra)│          │  │
│  │  └──────────┘└──────────┘└──────────┘          │  │
│  │                                                 │  │
│  │  DOMAIN SPECIALISTS (share LEGACY core memory via NATS)│  │
│  └─────────────────────┬───────────────────────────┘  │
│                        │                              │
│                        ▼                              │
│  ┌─────────────────────────────────────────────────┐  │
│  │               ACTION LAYER                      │  │
│  │         (OpenClaw MCP + NATS JetStream)         │  │
│  │                                                 │  │
│  │  ┌────────┐┌────────┐┌────────┐┌────────┐┌────┐│  │
│  │  │Todoist ││Dropbox ││ Email  ││Home Ast││Cust││  │
│  │  │ Sync   ││ Vault  ││ Tools  ││ Tools  ││MCPs││  │
│  │  └────────┘└────────┘└────────┘└────────┘└────┘│  │
│  └─────────────────────┬───────────────────────────┘  │
│                        │                              │
│                        ▼                              │
│  ┌─────────────────────────────────────────────────┐  │
│  │             KNOWLEDGE LAYER                     │  │
│  │            (Obsidian + Dropbox)                  │  │
│  │                                                 │  │
│  │  ┌────────────────────────────────────────────┐  │  │
│  │  │             One Mind Codex                │  │  │
│  │  │  Command/ │ SO/ │ HP/ │ LE/ │ GE/ │.archive│  │  │
│  │  └────────────────────────────────────────────┘  │  │
│  │                      ▲                          │  │
│  │                      │                          │  │
│  │                 Dropbox Sync                    │  │
│  │           ┌──────────┼──────────┐               │  │
│  │        Windows      Mac      Mobile             │  │
│  └─────────────────────────────────────────────────┘  │
│                                                       │
└═══════════════════════════════════════════════════════┘
```

---

## Area Folder Structure

```
[Pillar]/[Area]/
│
├── _Index.md                    # Area overview
│   │
│   │  ## [Area Name]
│   │  **Pillar:** [HP/LE/GE/IT]
│   │  **Status:** Active
│   │  **Agent:** [Yes/No]
│   │
│   │  ### Current Focus
│   │  [What you're working on now]
│   │
│   │  ### Metrics
│   │  [Key numbers]
│   │
│   │  ### Active Projects
│   │  - [[Project 1]]
│   │  - [[Project 2]]
│
├── Knowledge/                   # Reference material
│   ├── Protocols/
│   ├── Frameworks/
│   └── SOPs/
│
├── Projects/                    # Active work
│   ├── [Project-Name].md
│   └── _Completed/
│
├── Tasks/                       # Syncs to Todoist
│   ├── Active.md
│   ├── Backlog.md
│   └── Recurring.md
│
└── Workflows/                   # Automations & routines
    ├── [Workflow-1].md
    └── [Workflow-2].md
```

---

## Cognitive Framework

The `00-Core` folder contains deep philosophical/cognitive architecture:

- Cognitive frameworks (human + artificial cognition)
- Decision-making protocols
- "One Mind" governance layer
- Long-term civilization design thinking
- Solar Punk / Cyber Punk vision concepts
- OpenClaw agent governance (L0–L3 autonomy)

---

## Summary

This is a **personal operating system** combining:

- **Knowledge Management** (Second Brain / PKM)
- **AI Agent Infrastructure** (OpenClaw — LEGACY + 8 domain agents on NATS JetStream)
- **Life/Task Management** (Todoist integration, HIVE Capture pipeline)
- **Philosophical Framework** (Cognitive architecture, solar-punk governance)
