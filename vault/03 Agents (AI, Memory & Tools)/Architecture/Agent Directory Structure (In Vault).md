---
id: ui-02-20260111053925
title: "Agent Directory Structure (In Vault)"
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
ai_summary: "Agent Directory Structure (In Vault)"
---

Command/
├── Legacy/
│   ├── Charter.md                    # Legacy AI operating rules
│   │
│   ├── Agents/
│   │   ├── _Registry.md              # Master agent directory
│   │   │
│   │   ├── Legacy-AI/                # ORCHESTRATOR
│   │   │   ├── Config.md             # Letta config reference
│   │   │   ├── System-Prompt.md      # Core prompt
│   │   │   ├── Memory-Blocks.md      # What memory blocks attached
│   │   │   ├── Tools.md              # n8n workflows available
│   │   │   └── Logs/
│   │   │
│   │   ├── HP-Agent/                 # PILLAR AGENT
│   │   │   ├── Config.md
│   │   │   ├── System-Prompt.md
│   │   │   ├── Memory-Blocks.md
│   │   │   ├── Tools.md              # HP-specific n8n workflows
│   │   │   ├── Folder-Access.md      # Which folders agent can read/write
│   │   │   └── Logs/
│   │   │
│   │   ├── LE-Agent/
│   │   │   ├── Config.md
│   │   │   ├── System-Prompt.md
│   │   │   ├── Memory-Blocks.md
│   │   │   ├── Tools.md
│   │   │   ├── Folder-Access.md
│   │   │   ├── Sub-Agents/           # Area agents under LE
│   │   │   │   └── Homestead-Agent/
│   │   │   │       ├── Config.md
│   │   │   │       ├── System-Prompt.md
│   │   │   │       ├── Memory-Blocks.md
│   │   │   │       ├── Tools.md
│   │   │   │       └── Folder-Access.md
│   │   │   └── Logs/
│   │   │
│   │   ├── GE-Agent/
│   │   │   └── [same structure]
│   │   │
│   │   └── IT-Agent/
│   │       └── [same structure]
│   │
│   ├── Coordination/
│   │   ├── Cross-Pillar-Projects/
│   │   └── A2A-Protocols.md          # How agents communicate
│   │
│   └── n8n/
│       ├── Workflow-Registry.md      # All available n8n workflows
│       ├── MCP-Endpoints.md          # MCP server configs
│       └── Tool-Documentation/