---
type: comparison
domain: SO/02
tags: [openclaw, lobster, workflows, automation, efficiency]
created: 2026-05-12
updated: 2026-05-12
---

# Lobster vs Manual Orchestration

Comparison of workflow execution approaches in OpenClaw: manual LLM orchestration vs Lobster deterministic pipelines.

## Before Lobster (manual LLM orchestration)

**What happened when you said "run one mind status":**

```
Turn 1: Read skill file → realize I need to gather data
  → exec: bash status-gather.sh | python3 status-render.py
  → Wait for result
  
Turn 2: Got result → realize I need HTML too
  → exec: bash status-gather.sh | python3 status-html.py > status.html
  → Wait for result
  
Turn 3: Got both → show you the status text
  → Ask: "Should I post this to #briefings?"
  
Turn 4: You approve
  → message: send to Discord channel 1503277459565445201
  → Wait for result
  
Turn 5: Confirm delivery
```

**Cost:**
- **5 LLM turns** (back and forth)
- **4 tool calls** (read, exec, exec, message)
- **~10,000 tokens** (rough estimate for context + responses)
- **Manual approval handling** (I had to ask, wait, then act)

**Problems:**
- ❌ I had to remember each step
- ❌ Gateway restart would lose progress
- ❌ No built-in approval gate (I had to improvise)
- ❌ Each step costs tokens for orchestration
- ❌ Approval was informal ("Should I post this?")
- ❌ No resumable state if something failed

---

## With Lobster (deterministic pipeline)

**What happens when you say "run one mind status":**

```
Turn 1: Call lobster with workflow file
  → Lobster runs all 5 steps internally
  → Pauses at approval gate
  → Returns: status + resume token
  
Turn 2: You approve (or I approve on your behalf)
  → lobster: resume with token + approve=true
  → Returns: final output
```

**Cost:**
- **2 LLM turns** (request → approval → done)
- **2 tool calls** (lobster run, lobster resume)
- **~3,000 tokens** (much less orchestration)
- **Built-in approval gate** (pause/resume is native)

**Benefits:**
- ✅ Single declarative workflow file
- ✅ Steps run in-process (faster, no subprocess overhead per step)
- ✅ Resume tokens survive gateway restarts
- ✅ Approval is a first-class primitive
- ✅ Workflow is auditable (see exactly what will run)
- ✅ Pipeline data piping (`$gather.stdout` → render)

---

## Side-by-Side Comparison

| Dimension | Manual (Before) | Lobster (After) |
|-----------|----------------|-----------------|
| **LLM turns** | 5 | 2 |
| **Tool calls** | 4 | 2 |
| **Token cost** | ~10K | ~3K |
| **Orchestration** | LLM decides each step | Declared in `.lobster` file |
| **Approval** | Informal (LLM asks) | Built-in gate with resume token |
| **Resumable** | ❌ No | ✅ Yes (token survives restarts) |
| **Deterministic** | ❌ LLM might vary | ✅ Always same steps |
| **Auditable** | ❌ Implicit in prompt | ✅ Explicit in workflow file |
| **Data piping** | Manual (run gather twice) | Native (`$step.stdout`) |

---

## Real Example: One Mind Status

### Manual orchestration (before)

```bash
# I ran gather twice (wasteful):
bash status-gather.sh | python3 status-render.py
bash status-gather.sh | python3 status-html.py > status.html

# Then posted manually
message(action=send, target=1503277459565445201, message="...")
```

**Issues:**
- Gathered data twice (inefficient)
- No structured approval
- Manual coordination across multiple exec calls
- No state if interrupted

### Lobster workflow (after)

```yaml
name: onemind-status
description: "ONE MIND STATUS — infra, sync, tasks, vault, crons, sentinel."

steps:
  - id: gather
    command: bash /home/ubuntu/.openclaw/workflows/status-gather.sh

  - id: render
    stdin: $gather.stdout  # Reuses gather output
    command: python3 /home/ubuntu/.openclaw/workflows/status-render.py

  - id: html
    stdin: $gather.stdout  # Reuses gather output
    command: |
      python3 /home/ubuntu/.openclaw/workflows/status-html.py > /home/ubuntu/.openclaw/canvas/status.html
      echo "Dashboard written"

  - id: approve
    approval: required
    prompt: "ONE MIND STATUS ready. Post text board to #briefings and include dashboard URL?"

  - id: deliver
    command: |
      cat << 'EOF' > /tmp/status-message.txt
      $render.stdout
      
      📊 http://100.123.126.34:8099/status.html
      EOF
      cat /tmp/status-message.txt
    condition: $approve.approved
```

**Benefits:**
- Single gather pass, piped to multiple consumers
- Built-in approval checkpoint
- Resumable if gateway restarts
- Self-documenting workflow
- ~70% token reduction for recurring runs

---

## When to Use Each

### Manual Orchestration Is Better For:

- **Simple one-shot tasks** ("check disk space")
- **Exploration/debugging** ("try this, see what happens")
- **User wants to see each step** before proceeding
- **Workflow is too dynamic** (steps depend on complex LLM decisions)

### Lobster Wins For:

- ✅ **Repeatable workflows** (morning brief, status checks, reports)
- ✅ **Approval gates needed** (before posting, sending emails, financial actions)
- ✅ **Long-running tasks** that survive restarts
- ✅ **Token efficiency** (recurring workflows, scheduled tasks)
- ✅ **Audit requirements** (compliance, change control, transparency)
- ✅ **Data reuse** (one gather → multiple consumers)

---

## Efficiency Gains

**For recurring workflows like One Mind Status:**

| Metric | Manual | Lobster | Improvement |
|--------|--------|---------|-------------|
| Turns per run | 5 | 2 | **60% reduction** |
| Tool calls | 4 | 2 | **50% reduction** |
| Tokens | ~10K | ~3K | **70% reduction** |
| Data gathering | 2x | 1x | **50% reduction** |
| Resumability | ❌ | ✅ | Infinite improvement |

**Annual savings** (assuming 365 daily runs):
- **2,555,000 tokens saved** per year
- **At $3/1M input tokens**: **~$7.67/year** (modest but adds up across all workflows)
- **More importantly**: deterministic, auditable, resumable execution

---

## Bottom Line

Lobster moves orchestration from **LLM tokens → runtime engine**.

For recurring workflows, it's:
- **3x more token efficient**
- **Infinitely more reliable** (resumable state)
- **Fully auditable** (workflow file = documentation)
- **Approval-native** (no improvising gates)

**Rule of thumb:** If you'll run it more than 3 times, write a Lobster workflow.

---

## Related Files

- Workflow: `~/.openclaw/workflows/onemind-status.lobster`
- Gather script: `~/.openclaw/workflows/status-gather.sh`
- Render script: `~/.openclaw/workflows/status-render.py`
- HTML script: `~/.openclaw/workflows/status-html.py`
- Skill: `~/.openclaw/plugin-skills/onemind-status/SKILL.md`

## Documentation

- [Lobster docs](https://docs.openclaw.ai/tools/lobster)
- [Task Flow docs](https://docs.openclaw.ai/automation/taskflow)
- [Automation overview](https://docs.openclaw.ai/automation)
