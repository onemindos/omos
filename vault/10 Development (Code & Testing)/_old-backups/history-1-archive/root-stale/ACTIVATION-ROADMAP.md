---
id: root-00-20260113225954
title: "ACTIVATION ROADMAP"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "ACTIVATION ROADMAP"
---
# OneMind-Codex Activation Roadmap

> **Status:** READY FOR EXECUTION
> **Created:** 2026-01-10
> **Purpose:** Activate all OneMind systems for operational use

---

## Executive Summary

The OneMind-Codex is **85% built but 0% activated**. This document provides a phased roadmap to bring all systems online.

### Current State

| Component | Design | Implementation | Activation |
|-----------|--------|----------------|------------|
| Vault Structure | 100% | 90% | 80% |
| Agent Architecture | 100% | 85% | 0% |
| Letta Tools | 100% | 85% | 0% |
| Memory Blocks | 100% | 20% | 0% |
| N8N Workflows | 80% | 50% | 0% |
| Obsidian Plugin | 100% | 90% | 0% |

### Blockers Identified

1. **Security Issues** (CRITICAL)
   - Home Assistant token hardcoded in `letta-tools/tools/home_assistant/home_assistant_tool.py`
   - GitHub token reference in `letta-tools/tools/core/legacy_tools.py`

2. **Missing Configuration**
   - No `.env` file with actual API keys
   - Zoho OAuth tokens not generated
   - Letta agents not created

3. **Incomplete Setup**
   - System prompts directory empty (`02.3.18 System Prompts/`)
   - Memory blocks not populated in Letta
   - No N8N workflows deployed

---

## Phase 1: Security & Configuration (Day 1-2)

### 1.1 Fix Security Issues

```bash
# Files to fix:
# 1. letta-tools/tools/home_assistant/home_assistant_tool.py
#    - Move HA_TOKEN to environment variable
#    - Replace hardcoded token with os.getenv("HA_TOKEN")

# 2. letta-tools/tools/core/legacy_tools.py
#    - Verify no hardcoded tokens
#    - Ensure all secrets use os.getenv()
```

### 1.2 Create Production .env File

```bash
cd /Users/zeuslegacy/Documents/OneMind-Codex/letta-tools
cp .env.example .env
# Edit with actual credentials
```

**Required Environment Variables:**

| Variable | Source | Priority |
|----------|--------|----------|
| `LETTA_API_KEY` | Letta Cloud Dashboard | HIGH |
| `TODOIST_API_TOKEN` | Todoist Settings > Integrations | HIGH |
| `GITHUB_TOKEN` | GitHub Settings > Developer > PAT | HIGH |
| `DISCORD_API_URL` | Your Discord bot server | HIGH |
| `HA_TOKEN` | Home Assistant > Profile > Long-Lived Token | MEDIUM |
| `ZOHO_CLIENT_ID` | Zoho API Console | MEDIUM |
| `ZOHO_CLIENT_SECRET` | Zoho API Console | MEDIUM |
| `ZOHO_CRM_REFRESH_TOKEN` | OAuth flow | MEDIUM |
| `ZOHO_DESK_REFRESH_TOKEN` | OAuth flow | MEDIUM |
| `ZOHO_DESK_ORG_ID` | Zoho Desk Settings | MEDIUM |
| `ZOHO_CALENDAR_REFRESH_TOKEN` | OAuth flow | MEDIUM |
| `CIRCLE_API_TOKEN` | Circle.so Admin | LOW |
| `CIRCLE_COMMUNITY_ID` | Circle.so Admin | LOW |
| `TELEGRAM_BOT_TOKEN` | @BotFather | LOW |
| `TELEGRAM_CHAT_ID` | Your chat ID | LOW |

### 1.3 Verify Tool Installation

```bash
cd letta-tools
pip install -r requirements.txt
python -c "import requests, todoist_api_python, yaml, letta; print('All dependencies OK')"
```

---

## Phase 2: Agent Deployment (Day 2-3)

### 2.1 Create Letta Agents

```bash
cd letta-tools/scripts

# Create all 5 agents (LEGACY + HP + LE + GE + IT)
python letta_agent_setup.py --all

# Or create individually:
python letta_agent_setup.py --agent legacy
python letta_agent_setup.py --agent hp
python letta_agent_setup.py --agent le
python letta_agent_setup.py --agent ge
python letta_agent_setup.py --agent it
```

### 2.2 Attach Tools to Agents

```bash
# Attach Discord tools (requires Discord bot running)
python attach_discord_tools.py

# Upload additional tools to Letta Cloud
# - Follow TOOL-UPLOAD-CHECKLIST.md in 07.4 letta/
```

### 2.3 Configure Memory Blocks

The memory architecture requires these blocks per agent:

| Block | Agent | Size Limit | Content |
|-------|-------|------------|---------|
| Zeus Profile | All | 8,000 chars | Identity, values, preferences |
| Current Context | All | 5,000 chars | Today's priorities, active focus |
| System State | All | 3,000 chars | Agent status, health metrics |
| HP Context | HP | 5,000 chars | Health goals, habits, metrics |
| LE Context | LE | 5,000 chars | Family state, home, legacy plans |
| GE Context | GE | 5,000 chars | Business metrics, deals, revenue |
| IT Context | IT | 5,000 chars | Infrastructure, projects, security |

---

## Phase 3: Integration Activation (Day 3-5)

### 3.1 Todoist Structure Setup

```bash
cd letta-tools/scripts
python todoist_setup.py
# Or
python setup_todoist_structure.py
```

**Expected Structure:**

```
Todoist Projects:
├── 00-24 SO Inbox
├── 25-49 HP Tasks
├── 50-74 LE Tasks
├── 75-99 GE Tasks
└── [Sections by sub-domain]
```

### 3.2 Discord Bot Setup

Requirements:
- Discord bot server running on port 8100
- Webhook URLs configured for notifications
- Approval channel created for HITL

```bash
# Verify Discord bot is accessible
curl http://localhost:8100/health
```

### 3.3 Home Assistant Connection

```bash
# Test Home Assistant connection
curl -X GET "https://1mwywqa9zvz1hhh9lofv5wytww4nswti.ui.nabu.casa/api/" \
  -H "Authorization: Bearer $HA_TOKEN"
```

### 3.4 Obsidian Plugin Build

```bash
cd .obsidian/plugins/letta-ai-agent
npm install
npm run build
# Restart Obsidian to load plugin
```

---

## Phase 4: Automation Workflows (Day 5-7)

### 4.1 N8N Workflow Deployment

Deploy these core workflows:

| Workflow | Purpose | Trigger |
|----------|---------|---------|
| Morning Briefing | Daily brief to Discord | Cron 5:30 AM |
| Evening Review | Daily summary | Cron 9:00 PM |
| Task Sync | Todoist ↔ Vault | Webhook |
| Memory Consolidation | Update agent memory | Cron 2:00 AM |
| Agent Health Check | Monitor all agents | Cron every 15 min |
| Approval Webhook | HITL handling | Webhook |

### 4.2 Configure Webhooks

| Source | Destination | Purpose |
|--------|-------------|---------|
| GitHub | N8N | Vault change triggers |
| Todoist | N8N | Task completion sync |
| Discord | N8N | Approval responses |
| Letta | N8N | Agent action hooks |

---

## Phase 5: Testing & Validation (Day 7-10)

### 5.1 Agent Communication Tests

```bash
# Test each agent responds correctly
# 1. Send test message to LEGACY
# 2. Verify routing to HP/LE/GE/Grid agents
# 3. Test cross-agent communication
```

### 5.2 Tool Execution Tests

| Tool | Test Command | Expected Result |
|------|--------------|-----------------|
| Todoist | Create test task | Task appears in Todoist |
| Discord | Send notification | Message in channel |
| GitHub | Read vault file | File content returned |
| Home Assistant | Get device states | JSON response |
| Zoho CRM | List leads | Lead array returned |

### 5.3 Workflow Tests

- [ ] Morning briefing arrives at configured time
- [ ] Task creation syncs to vault
- [ ] Approval requests appear in Discord
- [ ] Agent health check runs without errors

---

## Phase 6: Go-Live Checklist

### Pre-Launch

- [ ] All security issues fixed
- [ ] Environment variables configured
- [ ] 5 agents created in Letta
- [ ] All tools attached to agents
- [ ] Memory blocks populated
- [ ] Todoist structure created
- [ ] Discord bot running
- [ ] N8N workflows deployed
- [ ] Obsidian plugin installed

### Launch Day

- [ ] Enable Morning Briefing workflow
- [ ] Enable Evening Review workflow
- [ ] Test first HITL approval
- [ ] Verify agent responses
- [ ] Document any issues

### Post-Launch (Week 1)

- [ ] Daily check of agent health
- [ ] Review Morning/Evening briefings
- [ ] Tune agent personalities as needed
- [ ] Add missing skills
- [ ] Optimize workflows

---

## Quick Reference Commands

### Letta Agent Management

```bash
# List all agents
python letta_agent_setup.py --list

# Create specific agent
python letta_agent_setup.py --agent legacy

# Load skills into agent
python letta_agent_setup.py --load-skills
```

### Tool Testing

```bash
# Test Todoist
python -c "from tools.todoist.todoist_tools import get_tasks; print(get_tasks())"

# Test Discord
python -c "from tools.discord.discord_tools import send_discord_notification; send_discord_notification('Test', 'Hello')"

# Test Home Assistant
python -c "from tools.home_assistant.home_assistant_tool import home_assistant; print(home_assistant('summary'))"
```

### Obsidian Plugin

```bash
cd .obsidian/plugins/letta-ai-agent
npm run dev    # Development mode
npm run build  # Production build
```

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     HUMAN INTERFACE                               │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐        │
│  │   Obsidian    │  │    Todoist    │  │    Discord    │        │
│  │    Vault      │  │    Tasks      │  │     Chat      │        │
│  └───────┬───────┘  └───────┬───────┘  └───────┬───────┘        │
└──────────┼──────────────────┼──────────────────┼────────────────┘
           │                  │                  │
           ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                      N8N SYNC LAYER                               │
│  Morning Brief │ Evening Review │ Task Sync │ Memory Update       │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│                     LETTA AI LAYER                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    LEGACY (UI Agent)                         ││
│  │              Central Orchestrator + Tech                     ││
│  └─────────────────────────┬───────────────────────────────────┘│
│                            │                                     │
│       ┌────────────────────┼────────────────────┐               │
│       ▼                    ▼                    ▼               │
│  ┌─────────┐         ┌─────────┐         ┌─────────┐           │
│  │   HP    │         │   LE    │         │   GE    │           │
│  │  Agent  │         │  Agent  │         │  Agent  │           │
│  │ (Self)  │         │(Family) │         │(Business)│           │
│  └─────────┘         └─────────┘         └─────────┘           │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    SHARED MEMORY                             ││
│  │  Zeus Profile │ Current Context │ System State              ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│                    TOOL INTEGRATIONS                              │
│  GitHub │ Todoist │ Zoho │ Circle │ Home Assistant │ Discord    │
└─────────────────────────────────────────────────────────────────┘
```

---

## Estimated Timeline

| Phase | Duration | Effort | Dependencies |
|-------|----------|--------|--------------|
| Phase 1: Security & Config | 1-2 days | 2-4 hours | API keys ready |
| Phase 2: Agent Deployment | 1-2 days | 3-5 hours | Phase 1 complete |
| Phase 3: Integration | 2-3 days | 4-6 hours | Phase 2 complete |
| Phase 4: Automation | 2-3 days | 4-6 hours | Phase 3 complete |
| Phase 5: Testing | 3-4 days | 4-8 hours | All phases complete |
| Phase 6: Go-Live | 1 day | 2-3 hours | Testing passed |

**Total: 10-14 days to full activation**

---

## Support Resources

- **CLAUDE.md** — Repository overview for AI assistants
- **ONEMIND-CODEX.md** — Master system blueprint (1,900+ lines)
- **letta-tools/README.md** — Tool documentation
- **07.4 letta/TOOL-UPLOAD-CHECKLIST.md** — Tool deployment checklist

---

*"HUMAN (Zeus) + AI (LEGACY) = ONE MIND"*
