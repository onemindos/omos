---
id: ui-02-20260111053925
title: "Letta + Todoist + GitHub: Complete Architecture Research"
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
ai_summary: "Letta + Todoist + GitHub: Complete Architecture Research"
---

# Letta + Todoist + GitHub: Complete Architecture Research

> **Date:** December 18, 2025  
> **Purpose:** Deep research into Letta Cloud capabilities for One Mind architecture  
> **Goal:** Minimize n8n dependency while maximizing Letta-native orchestration

---

## Executive Summary

After extensive research into Letta's capabilities, the recommended architecture uses **Letta Cloud as the primary orchestration layer** with n8n reduced to minimal cron triggers. Letta's programmatic tool calling, sleep-time compute, and MCP integration make it capable of handling most workflows natively.

**Key Decision:** Use **Todoist** (not TaskNotes) for task management due to MCP integration support and mobile UX superiority.

---

## 1. Letta Cloud Capabilities Deep Dive

### 1.1 Programmatic Tool Calling (Dec 2025)

Letta agents can **write and execute Python scripts that chain multiple tools** using `run_code_with_tools`.

```python
# Agent can write this dynamically:
def morning_routine(tools):
    tasks = tools.todoist_get_tasks(filter="today")
    calendar = tools.gcal_get_events(date="today")
    vault = tools.letta_fs_search("protocols/morning")
    
    briefing = tools.compose_briefing(tasks, calendar, vault)
    tools.send_telegram(briefing)
    return briefing
```

**Key Insight:** The agent becomes its own orchestrator—no n8n needed for most workflows.

### 1.2 Sleep-Time Compute (Apr 2025)

Two agents share memory blocks:
- **Primary Agent** (fast model): Handles real-time conversations
- **Sleep-Time Agent** (strong model): Reorganizes memory during idle time

```
┌─────────────────────────────────────────────────────┐
│                  SHARED MEMORY                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │
│  │zeus_profile │  │daily_state  │  │protocols    │ │
│  └─────────────┘  └─────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────┘
         ▲                                    ▲
         │                                    │
    ┌────┴────┐                          ┌────┴────┐
    │ PRIMARY │  ←── same memory ───►   │ SLEEP   │
    │ AGENT   │                          │ AGENT   │
    │ (fast)  │                          │ (strong)│
    └─────────┘                          └─────────┘
         │                                    │
    Real-time                           Background
    Conversations                       Processing
```

**Benefits:**
- Agent improves autonomously over time
- Memory gets reorganized and optimized
- Insights extracted from conversation history
- No active usage required for improvement

### 1.3 Skill Learning (Dec 2025)

Skills are stored as `.md` files and improve with experience:

```markdown
# Skill: Morning Briefing
## When to use
- Time is between 5:00-7:00 AM
- User hasn't received briefing today

## Steps
1. Get today's calendar events
2. Get today's Todoist tasks  
3. Check weather for Haven
4. Compose briefing with priorities
5. Send via preferred channel

## Learned Preferences
- Zeus prefers bullet points over paragraphs
- Include estimated time for each task
- Flag anything blocking GE pillar work
```

**36.8% performance improvement** on benchmarks with skill learning enabled.

### 1.4 Letta Filesystem (Jul 2025)

Native vault access with semantic search:

| Tool | Purpose |
|------|---------|
| `open(path)` | Read file contents |
| `grep(pattern, path)` | Text search across files |
| `semantic_search(query)` | AI-powered content search |
| `list_files(folder)` | Directory listing |

**Score:** 74% on LoCoMo long-context benchmark

### 1.5 Memory Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    LETTA MEMORY SYSTEM                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              IN-CONTEXT MEMORY                       │   │
│  │  (Always visible to agent, ~4K tokens)               │   │
│  │                                                      │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ │   │
│  │  │ zeus_profile │ │ daily_state  │ │ protocols    │ │   │
│  │  │              │ │              │ │ _cache       │ │   │
│  │  │ - Identity   │ │ - Today date │ │ - Morning    │ │   │
│  │  │ - 4 Pillars  │ │ - Active ctx │ │ - Evening    │ │   │
│  │  │ - Preferences│ │ - Energy lvl │ │ - Decision   │ │   │
│  │  └──────────────┘ └──────────────┘ └──────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                 │
│                           ▼                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              ARCHIVAL MEMORY                         │   │
│  │  (Unlimited, retrieved via search)                   │   │
│  │                                                      │   │
│  │  • Conversation history                              │   │
│  │  • Learned patterns                                  │   │
│  │  • Project contexts                                  │   │
│  │  • Decision rationales                               │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                 │
│                           ▼                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              EXTERNAL MEMORY                         │   │
│  │  (Via MCP/Tools)                                     │   │
│  │                                                      │   │
│  │  • Todoist (tasks)                                   │   │
│  │  • GitHub (vault files)                              │   │
│  │  • Obsidian (via filesystem)                         │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Memory Tools:**
- `memory_replace(label, value)` - Update memory block
- `memory_insert(label, content)` - Append to memory block
- `memory_rethink(label)` - Trigger reorganization
- `archival_insert(content)` - Save to long-term storage
- `archival_search(query)` - Retrieve from long-term storage

---

## 2. MCP Integration Architecture

### 2.1 Available MCP Servers

| Server | Purpose | Tools Provided |
|--------|---------|----------------|
| `@modelcontextprotocol/server-github` | Vault sync | file CRUD, search, commits |
| `Todoist MCP` | Task management | create, update, complete, query |
| `@modelcontextprotocol/server-memory` | External memory | store, retrieve, search |
| `@modelcontextprotocol/server-fetch` | Web access | fetch URLs, extract content |

### 2.2 Letta MCP Manager

```python
# Letta's MCP integration
from letta import MCPServerManager

# Create MCP server connection
mcp = client.mcp.create_server(
    name="todoist",
    server_type="sse",  # or "stdio"
    endpoint="https://todoist-mcp.zeusops.dev/sse"
)

# List available tools
tools = client.mcp.list_server_tools(server_id=mcp.id)

# Execute tool
result = client.mcp.execute_tool(
    server_id=mcp.id,
    tool_name="todoist_create_task",
    arguments={"content": "Review morning briefing", "project": "HP"}
)
```

---

## 3. TaskNotes vs Todoist Analysis

### 3.1 Comparison Matrix

| Feature | TaskNotes | Todoist | Winner |
|---------|-----------|---------|--------|
| **Vault Integration** | Native (lives in vault) | Via MCP | TaskNotes |
| **Mobile UX** | Requires Obsidian app | Native app, excellent | Todoist |
| **API Access** | HTTP localhost:8080 | REST API + MCP | Todoist |
| **Natural Language** | NLP parsing built-in | Natural language input | Tie |
| **Offline Support** | Full (local files) | Syncs when online | TaskNotes |
| **Webhooks** | Full support | Via MCP | Tie |
| **External Access** | Requires tunnel/VPN | Cloud-native | Todoist |
| **Recurring Tasks** | Basic | Advanced | Todoist |
| **Filters/Views** | Limited | Powerful | Todoist |
| **Collaboration** | None | Full | Todoist |

### 3.2 Decision: Use Todoist

**Reasons:**
1. MCP integration already exists for Letta
2. Superior mobile experience (Zeus is mobile-heavy)
3. Cloud-native API (no tunneling needed)
4. Advanced recurring task support
5. Collaboration features for future team growth

**Mitigation for vault integration:**
- Sync completed tasks to vault daily via Letta
- Use archival memory to store task context
- Generate weekly/monthly task reports in vault

---

## 4. Recommended Architecture

### 4.1 System Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         ONE MIND ARCHITECTURE                            │
│                         (Letta-Native Design)                            │
└─────────────────────────────────────────────────────────────────────────┘

                              ┌──────────────┐
                              │    ZEUS      │
                              │  (Human)     │
                              └──────┬───────┘
                                     │
            ┌────────────────────────┼────────────────────────┐
            │                        │                        │
            ▼                        ▼                        ▼
    ┌───────────────┐      ┌───────────────┐      ┌───────────────┐
    │   Telegram    │      │   Todoist     │      │   Voice       │
    │   (chat)      │      │   (mobile)    │      │   (future)    │
    └───────┬───────┘      └───────┬───────┘      └───────┬───────┘
            │                      │                      │
            └──────────────────────┼──────────────────────┘
                                   │
                                   ▼
                    ┌──────────────────────────────┐
                    │     LETTA CLOUD AGENT        │
                    │     "one-mind-zeus"          │
                    │                              │
                    │  ┌────────────────────────┐ │
                    │  │     MEMORY BLOCKS      │ │
                    │  │  • zeus_profile        │ │
                    │  │  • daily_state         │ │
                    │  │  • protocols_cache     │ │
                    │  └────────────────────────┘ │
                    │                              │
                    │  ┌────────────────────────┐ │
                    │  │        TOOLS           │ │
                    │  │  • run_code_with_tools │ │
                    │  │  • memory_*            │ │
                    │  │  • archival_*          │ │
                    │  │  • MCP connections     │ │
                    │  └────────────────────────┘ │
                    └──────────────┬───────────────┘
                                   │
         ┌─────────────────────────┼─────────────────────────┐
         │                         │                         │
         ▼                         ▼                         ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   TODOIST MCP   │     │   GITHUB MCP    │     │  LETTA FS       │
│                 │     │                 │     │                 │
│  • Tasks CRUD   │     │  • Vault files  │     │  • Protocol     │
│  • Projects     │     │  • Commits      │     │    files        │
│  • Filters      │     │  • Issues       │     │  • Reference    │
│  • Labels       │     │  • PRs          │     │    docs         │
└─────────────────┘     └─────────────────┘     └─────────────────┘
         │                         │                         │
         └─────────────────────────┼─────────────────────────┘
                                   │
                                   ▼
                    ┌──────────────────────────────┐
                    │         n8n                   │
                    │    (Minimal Role)             │
                    │                              │
                    │  • Cron: Morning 5:30 AM     │
                    │  • Cron: Evening 9:00 PM     │
                    │  • Webhook: Calendar events  │
                    │  • Webhook: SMS gateway      │
                    └──────────────────────────────┘
```

### 4.2 Data Flow Example: Morning Briefing

```
5:30 AM
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│ n8n Cron Trigger                                                 │
│ POST letta.zeusops.dev/agents/one-mind-zeus/messages            │
│ Body: {"message": "TRIGGER: morning_briefing"}                  │
└─────────────────────────────────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│ Letta Agent receives trigger                                     │
│                                                                  │
│ Agent thinks: "Morning briefing requested. I'll use my          │
│ run_code_with_tools to orchestrate this."                       │
│                                                                  │
│ Agent writes and executes:                                       │
│ ┌────────────────────────────────────────────────────────────┐  │
│ │ def morning_briefing(tools):                                │  │
│ │     # 1. Get today's tasks                                  │  │
│ │     tasks = tools.todoist_get_tasks(filter="today")         │  │
│ │                                                             │  │
│ │     # 2. Get calendar                                       │  │
│ │     cal = tools.gcal_today()                                │  │
│ │                                                             │  │
│ │     # 3. Get protocol                                       │  │
│ │     protocol = tools.letta_fs_open("protocols/morning.md")  │  │
│ │                                                             │  │
│ │     # 4. Compose briefing                                   │  │
│ │     briefing = f"""                                         │  │
│ │     🌅 Good morning Zeus!                                   │  │
│ │                                                             │  │
│ │     📅 Today: {len(cal)} events                             │  │
│ │     ✅ Tasks: {len(tasks)} items                            │  │
│ │                                                             │  │
│ │     🎯 Top 3 Priorities:                                    │  │
│ │     {format_priorities(tasks, cal)}                         │  │
│ │     """                                                     │  │
│ │                                                             │  │
│ │     # 5. Send notification                                  │  │
│ │     tools.send_telegram(briefing)                           │  │
│ │                                                             │  │
│ │     # 6. Update memory                                      │  │
│ │     tools.memory_replace("daily_state", {                   │  │
│ │         "date": today(),                                    │  │
│ │         "briefing_sent": True,                              │  │
│ │         "top_priorities": extract_priorities(tasks)         │  │
│ │     })                                                      │  │
│ │                                                             │  │
│ │     return briefing                                         │  │
│ └────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│ Zeus receives Telegram message                                   │
│ 🌅 Good morning Zeus!                                           │
│ 📅 Today: 4 events                                              │
│ ✅ Tasks: 7 items                                               │
│ 🎯 Top 3 Priorities:                                            │
│    1. [GE] Review pitch deck - 2hrs                             │
│    2. [HP] Gym session - 1hr                                    │
│    3. [LE] Deploy n8n update - 30min                            │
└─────────────────────────────────────────────────────────────────┘
```

---

## 5. n8n Minimal Workflows

### 5.1 Reduced Responsibility

| Before (Heavy n8n) | After (Minimal n8n) |
|--------------------|---------------------|
| Task creation logic | Just trigger Letta |
| Calendar parsing | Just trigger Letta |
| Notification formatting | Just trigger Letta |
| Decision trees | Letta handles |
| State management | Letta memory |

### 5.2 Required n8n Workflows

```yaml
# 1. Morning Briefing Trigger
- name: "Morning Briefing"
  trigger: Cron "30 5 * * *"
  action: POST to Letta
  body: {"message": "TRIGGER: morning_briefing"}

# 2. Evening Review Trigger  
- name: "Evening Review"
  trigger: Cron "0 21 * * *"
  action: POST to Letta
  body: {"message": "TRIGGER: evening_review"}

# 3. Calendar Webhook Relay
- name: "Calendar Event"
  trigger: Google Calendar webhook
  action: POST to Letta
  body: {"message": "CALENDAR_EVENT: {{event_data}}"}

# 4. SMS Gateway (optional)
- name: "SMS Inbound"
  trigger: Twilio webhook
  action: POST to Letta
  body: {"message": "SMS from {{from}}: {{body}}"}
```

---

## 6. Implementation Phases

### Phase 1: Letta Cloud Setup (Week 1)

```python
from letta import Letta

client = Letta(api_key="your-letta-cloud-key")

# Create the One Mind agent
agent = client.agents.create(
    name="one-mind-zeus",
    model="anthropic/claude-sonnet-4-5",
    embedding="openai/text-embedding-3-small",
    memory_blocks=[
        {
            "label": "zeus_profile",
            "value": """
            Zeus Delacruz - One Mind System Operator
            
            4 PILLARS:
            - HP (Holistic Performance): Health, energy, wellbeing
            - LE (Legacy Evolution): Personal growth, learning
            - GE (Generational Endeavor): Business, wealth
            - IT (Innovative Technology): Systems, automation
            
            PREFERENCES:
            - Concise communication (bullet points)
            - Morning focus time sacred
            - Batch similar tasks
            - Sunday planning ritual
            """
        },
        {
            "label": "daily_state",
            "value": """
            Current Date: [auto-updated]
            Last Briefing: None
            Active Context: None
            Energy Level: Unknown
            Focus Block: None
            """
        },
        {
            "label": "protocols_cache",
            "value": """
            MORNING: Check calendar → Review tasks → Set 3 priorities
            EVENING: Log completions → Plan tomorrow → Reflect
            DECISION: Consult 4 pillars → Check protocols → Execute
            CAPTURE: Inbox first → Process in batch → Route to pillar
            """
        }
    ],
    tools=["run_code_with_tools"]
)

print(f"Agent created: {agent.id}")
```

### Phase 2: MCP Connections (Week 1-2)

```python
# Connect Todoist MCP
todoist_mcp = client.mcp.create_server(
    name="todoist",
    server_type="sse",
    endpoint="https://todoist-mcp.example.com/sse",
    config={"api_key": "your-todoist-api-key"}
)

# Connect GitHub MCP  
github_mcp = client.mcp.create_server(
    name="github",
    server_type="stdio",
    command="npx",
    args=["-y", "@modelcontextprotocol/server-github"],
    env={"GITHUB_TOKEN": "your-github-token"}
)

# Attach tools to agent
client.agents.attach_mcp_tools(
    agent_id=agent.id,
    mcp_server_ids=[todoist_mcp.id, github_mcp.id]
)
```

### Phase 3: Letta Filesystem Setup (Week 2)

```python
# Create folders for vault content
folders = [
    "protocols",      # Morning, evening, decision protocols
    "pillars",        # HP, LE, GE, IT context
    "projects",       # Active project briefs
    "reference"       # Key documents
]

for folder in folders:
    client.files.create_folder(
        agent_id=agent.id,
        path=folder
    )

# Upload protocol files
protocols = {
    "protocols/morning.md": "...",
    "protocols/evening.md": "...",
    "protocols/decision.md": "...",
    "pillars/hp.md": "...",
    "pillars/le.md": "...",
    "pillars/ge.md": "...",
    "pillars/it.md": "..."
}

for path, content in protocols.items():
    client.files.upload(
        agent_id=agent.id,
        path=path,
        content=content
    )
```

### Phase 4: Custom Tools (Week 2-3)

```python
# Notification tool
def send_notification(channel: str, message: str) -> str:
    """
    Send notification to Zeus via specified channel.
    
    Args:
        channel: "telegram" | "sms" | "email"
        message: Content to send
    
    Returns:
        Confirmation of delivery
    """
    if channel == "telegram":
        # Use Telegram Bot API
        requests.post(
            f"https://api.telegram.org/bot{BOT_TOKEN}/sendMessage",
            json={"chat_id": ZEUS_CHAT_ID, "text": message}
        )
    elif channel == "sms":
        # Use Twilio
        twilio_client.messages.create(
            to=ZEUS_PHONE,
            from_=TWILIO_NUMBER,
            body=message
        )
    return f"Sent via {channel}"

# Register with Letta
client.tools.create(
    agent_id=agent.id,
    name="send_notification",
    source_code=inspect.getsource(send_notification),
    description="Send notification to Zeus"
)
```

### Phase 5: n8n Cron Setup (Week 3)

```json
{
  "nodes": [
    {
      "name": "Cron",
      "type": "n8n-nodes-base.cron",
      "parameters": {
        "triggerTimes": {
          "item": [{"hour": 5, "minute": 30}]
        }
      }
    },
    {
      "name": "Letta Trigger",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "https://api.letta.com/agents/{{AGENT_ID}}/messages",
        "headers": {
          "Authorization": "Bearer {{LETTA_API_KEY}}"
        },
        "body": {
          "message": "TRIGGER: morning_briefing"
        }
      }
    }
  ]
}
```

---

## 7. Why This Architecture Works

### 7.1 Benefits

| Benefit | Description |
|---------|-------------|
| **Agent Autonomy** | Letta handles orchestration, not n8n |
| **Stateful Memory** | Context persists across conversations |
| **Self-Improving** | Sleep-time compute optimizes over time |
| **Skill Learning** | Agent learns Zeus's preferences |
| **Native Vault Access** | Letta Filesystem replaces sync complexity |
| **Reduced Maintenance** | Fewer n8n workflows to manage |
| **Mobile-First** | Todoist + Telegram for on-the-go |

### 7.2 Trade-offs Accepted

| Trade-off | Mitigation |
|-----------|------------|
| Todoist not in vault | Daily sync to archival memory |
| Cloud dependency | Letta Cloud is production-grade |
| API costs | Sleep-time uses efficient models |
| Learning curve | Start with core workflows only |

---

## 8. Quick Reference

### Letta Cloud URLs
- Dashboard: https://cloud.letta.com
- API: https://api.letta.com
- Docs: https://docs.letta.com

### Key Endpoints
```
POST /agents/{id}/messages     - Send message to agent
GET  /agents/{id}/memory       - Get memory state
POST /agents/{id}/tools        - Create custom tool
POST /mcp/servers              - Connect MCP server
```

### Memory Block Labels
- `zeus_profile` - Identity and preferences
- `daily_state` - Current context
- `protocols_cache` - Quick-access protocols

### MCP Servers Needed
1. Todoist MCP (task management)
2. GitHub MCP (vault sync)
3. Fetch MCP (web access)

---

## 9. Next Actions

- [ ] Sign up for Letta Cloud
- [ ] Create `one-mind-zeus` agent with memory blocks
- [ ] Deploy/connect Todoist MCP server
- [ ] Upload protocol files to Letta Filesystem
- [ ] Create `send_notification` custom tool
- [ ] Configure n8n cron workflows (morning/evening)
- [ ] Test end-to-end morning briefing flow

---

## References

- [Letta Programmatic Tool Calling](https://www.letta.com/blog/programmatic-tool-calling-with-any-llm)
- [Letta Skill Learning](https://www.letta.com/blog/skill-learning)
- [Letta Sleep-Time Compute](https://www.letta.com/blog/sleep-time-compute)
- [Letta Filesystem](https://www.letta.com/blog/letta-filesystem)
- [MCP Servers Registry](https://github.com/modelcontextprotocol/servers)
- [Letta Documentation](https://docs.letta.com)
