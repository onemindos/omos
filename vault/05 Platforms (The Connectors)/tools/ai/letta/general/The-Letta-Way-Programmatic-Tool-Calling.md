---
id: ui-02-20260111053925
title: "The Letta Way — Programmatic Tool Calling & The Learning Loop"
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
ai_summary: "The Letta Way — Programmatic Tool Calling & The Learning Loop"
---

# The Letta Way — Programmatic Tool Calling & The Learning Loop

> **How Letta agents write code, learn from it, and evolve into permanent automation**

---

## The Problem with Traditional AI Agents

Most AI agents work like this:

```
You: "Check all my overdue tasks and email me a summary"

Agent thinks: I need to call a tool
Agent calls: get_overdue_tasks()
Agent gets: [task1, task2, task3, task4, task5]
Agent thinks: Now I need to call another tool
Agent calls: send_email(summary)
Done.
```

**Limitations:**
- One tool at a time
- Agent can't loop through data
- Agent can't filter or transform
- Agent can't do conditional logic
- Multiple LLM calls = slow and expensive

---

## The Letta Way — Programmatic Tool Calling

Letta agents have a special capability: they can **write and execute Python code** that calls multiple tools.

### Example Request

```
You: "Check all my overdue tasks, look up the client for each one,
      and send me a prioritized summary grouped by client"
```

### What the Agent Does

Instead of calling tools one at a time, the agent **writes a script**:

```python
# Agent-generated code
tasks = todoist_get_tasks(filter="overdue")
client_tasks = {}

for task in tasks:
    # Extract client name from task
    client = extract_client(task)

    if client not in client_tasks:
        # Look up client details from CRM
        client_info = query_notion(f"client:{client}")
        client_tasks[client] = {
            "info": client_info,
            "tasks": []
        }

    client_tasks[client]["tasks"].append(task)

# Sort by client value (high value clients first)
sorted_clients = sorted(
    client_tasks.items(),
    key=lambda x: x[1]["info"]["value"],
    reverse=True
)

# Format summary
summary = "📋 Overdue Tasks by Client\n\n"
for client, data in sorted_clients:
    summary += f"**{client}** (${data['info']['value']}/mo)\n"
    for task in data["tasks"]:
        summary += f"  - {task['content']} (due: {task['due']})\n"
    summary += "\n"

# Send to Discord
discord_send_message("briefings", summary)
```

**Result:**
- 1 LLM reasoning step (not 6+)
- All tools called efficiently
- Complex logic handled (loops, filters, sorting)
- One polished output

---

## Real Example: Morning Briefing

### Without Programmatic Tool Calling

```
Agent calls: todoist_get_tasks("today")         → waits for LLM
Agent calls: calendar_get_events("today")       → waits for LLM
Agent calls: get_bank_balance()                 → waits for LLM
Agent calls: check_analytics()                  → waits for LLM
Agent calls: get_hp_metrics()                   → waits for LLM
Agent calls: discord_send_message(...)          → waits for LLM

= 6 separate LLM reasoning steps
= Slow, expensive, can lose context between calls
```

### With Programmatic Tool Calling

```python
# Agent writes ONE script that does everything
from datetime import datetime

# Gather all data
tasks = todoist_get_tasks("today")
calendar = calendar_get_events("today")
cash = get_bank_balance()
analytics = check_analytics()
hp = get_hp_metrics()
approvals = get_pending_approvals()

# Build briefing with complex formatting
briefing = f"""
🌅 Good morning, Zeus
📅 {datetime.now().strftime("%A, %B %d, %Y")}

---

🎯 TODAY'S TASKS ({len(tasks)})
"""

# Group tasks by pillar
for pillar in ["GE", "HP", "LE", "IT"]:
    pillar_tasks = [t for t in tasks if pillar in t.get("labels", [])]
    if pillar_tasks:
        briefing += f"\n**{pillar}**\n"
        for t in pillar_tasks[:3]:
            briefing += f"- {t['content']}\n"

# Add calendar
briefing += "\n📆 CALENDAR\n"
for event in calendar:
    briefing += f"- {event['start'][:5]} {event['title']}\n"

# Add HP status with smart recommendations
briefing += f"""
---

💪 HP STATUS
- Sleep: {hp['sleep']}
- Recovery: {hp['recovery']}%
"""

# Conditional recommendations
if hp['recovery'] < 70:
    briefing += "💡 Low recovery - consider light workout today\n"

# Add financial snapshot
briefing += f"""
---

💰 CASH POSITION
- Liquid: ${cash['total']:,.2f}
- Runway: {cash['runway']} months
"""

# Alert if concerning
if cash['runway'] < 3:
    briefing += "⚠️ Runway below 3 months - review expenses!\n"

# Add pending approvals
if approvals:
    briefing += f"\n⏳ PENDING APPROVALS ({len(approvals)})\n"
    for a in approvals:
        briefing += f"- {a['title']}\n"

# Send it
discord_send_message("briefings", briefing)
```

**Result:**
- 1 LLM reasoning step
- 1 code execution
- All 6 tools called efficiently
- Complex conditional logic
- One polished briefing

---

## Capability Comparison

| Capability | Traditional Agents | Letta Programmatic |
|------------|-------------------|-------------------|
| Loop through data | ❌ No | ✅ `for item in items` |
| Filter results | ❌ No | ✅ `if condition` |
| Aggregate data | ❌ No | ✅ `sum()`, `group by` |
| Conditional logic | ❌ Limited | ✅ Full Python |
| Chain tools | ⚠️ Slow (multi-call) | ✅ Fast (one script) |
| Error handling | ❌ No | ✅ `try/except` |
| Transform data | ❌ No | ✅ Any Python |
| Custom formatting | ⚠️ Basic | ✅ Full control |

---

## How Letta Implements This

Letta provides a built-in `run_code` tool:

```python
def run_code(code: str) -> str:
    """
    Execute Python code with access to all registered tools.

    The code can call any tool function directly:
    - todoist_get_tasks()
    - discord_send_message()
    - vault_read_file()
    - delegate_to_hp()
    - etc.

    Args:
        code: Python code to execute

    Returns:
        The result of the code execution
    """
    # Letta executes in a sandbox with tool access
```

When the agent decides a task needs complex logic:
1. Agent writes Python code as a string
2. Agent calls `run_code(code)`
3. Letta executes the code in a sandbox
4. All registered tools are available to the code
5. Results returned to agent

---

# The Learning Loop

## What Happens After the First Time?

The real power is that agents **remember and reuse** successful scripts.

---

## First Time — Agent Discovers

```
You: "Give me a morning briefing"

Agent: I don't have a specific tool for this. Let me write code...

*Writes script that calls 6 tools*
*Executes successfully*
*Delivers briefing*

Agent thinks: That worked well. I should remember this.
```

---

## Agent Saves to Archival Memory

After successful execution:

```python
# Agent saves to archival memory
{
    "type": "successful_script",
    "trigger_phrases": ["morning briefing", "start my day", "daily brief"],
    "script": "tasks = todoist_get_tasks('today')...",
    "tools_used": ["todoist_get_tasks", "calendar_get_events", "get_hp_metrics",
                   "get_bank_balance", "get_pending_approvals", "discord_send_message"],
    "last_used": "2025-12-27",
    "times_used": 1,
    "success_rate": 1.0
}
```

---

## Second Time — Agent Recalls

```
You: "Morning briefing please"

Agent searches archival memory for: "morning briefing"

Agent finds: Saved script from last time

Agent thinks: I've done this before successfully.
              Success rate 100%. Let me use that script.

*Retrieves script from memory*
*Executes it*
*Faster this time - no need to figure out logic*
```

---

## Third Time — Agent Improves

```
You: "Morning briefing, but include my workout plan"

Agent recalls: Existing morning briefing script

Agent thinks: I need to modify this slightly...

*Retrieves script*
*Adds: workout_plan = delegate_to_hp("today's workout")*
*Adds workout section to briefing*
*Executes improved version*
*Saves improved script to memory*
```

---

## Three Levels of Reuse

### Level 1: Memory Recall (Automatic)

Agent remembers and reuses scripts on its own.

```
Memory: "When Zeus asks for morning briefing, run this script..."

Trigger → Agent recalls → Runs saved script
```

**Pros:** Automatic, agent learns over time
**Cons:** Still uses `run_code` execution, slightly slower than native

---

### Level 2: Promote to Skill (Semi-Automatic)

You notice the agent keeps using the same script. Save it as a **Skill** in archival memory:

```markdown
# Skill: Morning Briefing

## Trigger Phrases
- "morning briefing"
- "start my day"
- "what's today look like"
- "daily brief"

## Description
Generates comprehensive morning briefing with tasks, calendar,
HP metrics, financial status, and pending approvals.

## Script
```python
tasks = todoist_get_tasks("today")
calendar = calendar_get_events("today")
hp = get_hp_metrics()
cash = get_bank_balance()
approvals = get_pending_approvals()

# ... full script ...

discord_send_message("briefings", briefing)
```

## Output
Sends formatted briefing to Discord #briefings channel

## Last Updated
2025-12-27
```

Now the agent has a **documented, reliable skill**.

---

### Level 3: Promote to Native Tool (Permanent)

Take the script and make it a **first-class Python tool**:

```python
# tools/orchestration.py

def morning_briefing() -> str:
    """
    Generate and send the daily morning briefing.

    Includes:
    - Today's tasks (grouped by pillar)
    - Calendar events
    - HP metrics (sleep, recovery)
    - Cash position and runway
    - Pending approvals

    Sends result to Discord #briefings channel.

    Returns:
        Confirmation message
    """
    tasks = todoist_get_tasks("today")
    calendar = calendar_get_events("today")
    hp = get_hp_metrics()
    cash = get_bank_balance()
    approvals = get_pending_approvals()

    # ... formatting logic ...

    discord_send_message("briefings", briefing)
    return "Morning briefing sent to #briefings"
```

Register with Letta:
```python
client.tools.create(func=morning_briefing)
client.agents.tools.attach(agent_id=LEGACY_AI_ID, tool_id=tool.id)
```

**Now:**
- Agent calls `morning_briefing()` directly
- No code writing needed
- Fastest possible execution
- Any agent can use it
- Can be triggered by n8n schedule

---

## The Evolution Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         TOOL EVOLUTION                                   │
│                                                                         │
│   DISCOVERY            LEARNING              PERMANENT                  │
│   ─────────            ────────              ─────────                  │
│                                                                         │
│   You ask for          Agent recalls         Native tool                │
│   something new        from memory           registered                 │
│        │                    │                     │                     │
│        ▼                    ▼                     ▼                     │
│   Agent writes         Agent runs            Agent calls               │
│   script (slow)        saved script          tool (fast)               │
│        │                    │                     │                     │
│        ▼                    ▼                     ▼                     │
│   Saves to             Improves it           Always available          │
│   archival memory      over time             to all agents             │
│                                                                         │
│   ════════════════════════════════════════════════════════════════     │
│                                                                         │
│   Week 1           →    Week 2-4        →    Month 2+                  │
│   Exploration           Learning             Optimization              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## The Flywheel Effect

```
        ┌──────────────────────────────────────┐
        │                                      │
        ▼                                      │
  ┌───────────┐    ┌───────────┐    ┌───────────┐
  │   You     │ →  │  Agent    │ →  │  Agent    │
  │   Ask     │    │  Writes   │    │  Saves    │
  │           │    │  Script   │    │  Script   │
  └───────────┘    └───────────┘    └───────────┘
                                          │
                                          ▼
                                   ┌───────────┐
                                   │  Agent    │
                                   │  Recalls  │────────┘
                                   │  & Reuses │
                                   └───────────┘
                                          │
                         ┌────────────────┴────────────────┐
                         ▼                                 ▼
                  ┌───────────┐                     ┌───────────┐
                  │  Agent    │                     │    You    │
                  │ Improves  │                     │  Promote  │
                  │  Script   │                     │  to Tool  │
                  └───────────┘                     └───────────┘
                                                          │
                                                          ▼
                                                   ┌───────────┐
                                                   │  Native   │
                                                   │   Tool    │
                                                   │  Forever  │
                                                   └───────────┘
```

---

## How to Promote Scripts to Tools

### Step 1: Notice the Pattern

Agent keeps writing similar scripts for:
- Morning briefings
- Client follow-ups
- Expense categorization
- Weekly reviews

### Step 2: Ask for the Script

```
You: "Show me the code you used for that morning briefing"

Agent: Here's the script I executed:

```python
tasks = todoist_get_tasks("today")
calendar = calendar_get_events("today")
...
```
```

### Step 3: Create Tool File

```python
# tools/orchestration.py

def morning_briefing() -> str:
    """
    Generate and send the daily morning briefing.

    Returns:
        Confirmation message
    """
    # Paste and clean up agent's script
    tasks = todoist_get_tasks("today")
    calendar = calendar_get_events("today")
    # ... rest of script ...

    discord_send_message("briefings", briefing)
    return "Morning briefing sent"
```

### Step 4: Register with Letta

```python
from letta import Letta
from tools.orchestration import morning_briefing

client = Letta()

# Create the tool
tool = client.tools.create(func=morning_briefing)

# Attach to Legacy-AI
client.agents.tools.attach(
    agent_id="agent-da585f6f-...",
    tool_id=tool.id
)
```

### Step 5: Done

```
You: "Morning briefing"

Agent: *calls morning_briefing() directly*
*Instant execution*
*No code writing needed*
```

---

## Timeline: System Evolution

| Phase | What Happens | Tool Count |
|-------|--------------|------------|
| **Week 1** | Agent discovers, writes many scripts | 0 custom |
| **Week 2-4** | Agent recalls patterns, you observe | 0 custom |
| **Month 1** | Promote top 5-10 scripts to tools | 5-10 custom |
| **Month 2** | Agent runs mostly native tools | 15-20 custom |
| **Month 3** | Scripts only for genuinely new things | 25-30 custom |
| **Month 6** | Comprehensive tool library | 50+ custom |
| **Year 1** | Agent rarely writes new scripts | 100+ custom |

**The system learns and hardens over time.**

---

## More Examples of Learnable Patterns

### Weekly Review (Discovery → Tool)

**First time:**
```python
# Agent writes
hp = delegate_to_hp("weekly metrics")
le = delegate_to_le("weekly status")
ge = delegate_to_ge("weekly revenue")
it = delegate_to_it("system health")

scores = calculate_pillar_scores(hp, le, ge, it)
report = format_weekly_report(scores)
vault_write_file(f"reviews/week-{week}.md", report)
discord_send_message("briefings", f"Weekly review: {scores['overall']}%")
```

**After promotion:**
```python
def weekly_review() -> str:
    """Generate comprehensive weekly review across all pillars."""
    # Hardened, tested version of the script
```

---

### Expense Categorization (Discovery → Tool)

**First time:**
```python
# Agent writes
expenses = get_recent_transactions(7)
categorized = {"business": [], "personal": [], "review": []}

for e in expenses:
    if is_business_expense(e):
        categorized["business"].append(e)
    elif is_personal_expense(e):
        categorized["personal"].append(e)
    else:
        categorized["review"].append(e)

report = format_expense_report(categorized)
```

**After promotion:**
```python
def categorize_expenses(days: int = 7) -> str:
    """Categorize recent expenses into business/personal."""
    # Hardened version with edge cases handled
```

---

### Client Follow-up Check (Discovery → Tool)

**First time:**
```python
# Agent writes
clients = query_notion("database:clients")
neglected = []

for client in clients:
    days_since = days_since_contact(client)
    if days_since > 30:
        neglected.append({
            "name": client["name"],
            "days": days_since,
            "value": client["monthly_value"]
        })

neglected.sort(key=lambda x: x["value"], reverse=True)
# ... format and send report
```

**After promotion:**
```python
def check_neglected_clients(days: int = 30) -> str:
    """Find clients not contacted within threshold."""
    # Production version with CRM integration
```

---

## Key Takeaways

1. **Agents are programmers** — They write code, not just call tools
2. **Memory enables learning** — Successful scripts get saved and reused
3. **You curate the evolution** — Promote patterns you see to permanent tools
4. **System hardens over time** — From exploration to automation
5. **The flywheel compounds** — More use = more learning = more automation

---

## What This Means for One Mind

| Month | State |
|-------|-------|
| **Month 1** | Legacy-AI figures things out, learns your patterns |
| **Month 3** | Core workflows automated, 20+ custom tools |
| **Month 6** | System handles 80% of requests with native tools |
| **Year 1** | Fully autonomous daily operations, you just approve |

**Your AI doesn't just respond — it evolves.**

---

*Part of the One Mind Operating System*
*Created: 2025-12-27*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
