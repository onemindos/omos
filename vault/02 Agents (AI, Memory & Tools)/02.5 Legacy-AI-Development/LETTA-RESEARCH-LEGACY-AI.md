---
id: root-00-20260113225954
title: "LETTA RESEARCH LEGACY AI"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub, platform/letta]
ai_summary: "LETTA RESEARCH LEGACY AI"
---
# LETTA DEEP DIVE: Complete Research Findings for Legacy AI Architecture

> Research compiled for designing Legacy AI as a central "mother agent" with shared memory across sub-agents

---

## Table of Contents
1. [The Core Problem Letta Solves](#1-the-core-problem-letta-solves)
2. [Memory Blocks - Full Technical Details](#2-memory-blocks---full-technical-details)
3. [Shared Memory Blocks](#3-shared-memory-blocks---your-architectures-foundation)
4. [Multi-Agent Communication](#4-multi-agent-communication---full-details)
5. [Compaction - Context Overflow Handling](#5-compaction---how-context-overflow-is-handled)
6. [Archival Memory - Unlimited Storage](#6-archival-memory---unlimited-storage)
7. [Sleep-Time Agents - Reducing Latency](#7-sleep-time-agents---reducing-legacy-ais-latency)
8. [Context Engineering](#8-context-engineering---agent-self-management)
9. [Agent State & Persistence](#9-agent-state--persistence)
10. [Explicit Limitations Found](#10-explicit-limitations-found)
11. [Recommendations for Legacy AI Architecture](#11-recommendations-for-your-legacy-ai-architecture)
12. [What to Test Empirically](#12-what-you-should-test-empirically)

---

## 1. THE CORE PROBLEM LETTA SOLVES

From the **MemGPT research paper** (the academic foundation of Letta):

> "Large language models are constrained by limited context windows, hindering their utility in tasks like extended conversations and document analysis."

**Letta's Solution:** "Virtual context management" inspired by operating systems - treating the LLM like a CPU with hierarchical memory (RAM, disk, etc.) where data moves between tiers automatically.

---

## 2. MEMORY BLOCKS - FULL TECHNICAL DETAILS

### Structure of a Memory Block
```typescript
{
  label: string,        // Unique identifier (e.g., "persona", "user_info")
  description: string,  // Tells agent what this block is FOR
  value: string,        // The actual content
  limit: number,        // Character limit YOU define
  read_only: boolean    // Optional - prevents agent modifications
}
```

### How Memory Blocks Work Internally
- **Prepended to agent's prompt in XML-like format**
- **Always visible** in the context window (no retrieval needed)
- Agent sees them every single turn
- Letta automatically manages the formatting

### Code Example - Creating Blocks
```typescript
// Create a standalone block
const block = await client.blocks.create({
  label: "organization",
  description: "Information about the user's organization",
  value: "Organization: Letta\nIndustry: AI\nSize: 50 employees",
  limit: 4000,
  read_only: false
});

// Create agent with memory blocks
const agent = await client.agents.create({
  memory_blocks: [
    {
      label: "persona",
      value: "I am Legacy AI, the central intelligence...",
      limit: 5000
    },
    {
      label: "human",
      value: "User preferences and context...",
      limit: 3000
    }
  ],
  block_ids: [block.id],  // Attach existing blocks
  model: "openai/gpt-4o-mini"
});
```

### Critical Finding: NO EXPLICIT LIMIT ON NUMBER OF BLOCKS

**The documentation does NOT specify a maximum number of memory blocks per agent.**

**BUT** - here's the real constraint:

> All memory blocks compete for space in the **context window** of your underlying LLM.

| Model | Context Window | Practical Block Budget |
|-------|---------------|----------------------|
| GPT-4o | 128K tokens (~500K chars) | ~150K-200K chars for blocks |
| GPT-4o-mini | 128K tokens | ~150K-200K chars for blocks |
| Claude 3.5 Sonnet | 200K tokens (~800K chars) | ~250K-300K chars for blocks |
| Claude 3 Opus | 200K tokens | ~250K-300K chars for blocks |

**Rule of thumb:** Reserve 30-40% of context for memory blocks, 40% for conversation, 20-30% for system/tools.

### Block Operations Available

| Operation | What It Does |
|-----------|-------------|
| Create | Make a new block |
| Attach | Connect block to agent |
| Detach | Disconnect from agent |
| Retrieve | Get block contents |
| List | See all blocks |
| Update | Change block value (replaces entirely) |
| Delete | Remove block |

### CRITICAL LIMITATION: Last Write Wins

> "Updating `value` completely replaces existing content"

If two agents write to the same shared block simultaneously, the last write wins. No merge, no conflict resolution.

---

## 3. SHARED MEMORY BLOCKS - Your Architecture's Foundation

### How Sharing Works
```typescript
// Create a shared block
const sharedState = await client.blocks.create({
  label: "global_state",
  description: "Shared state visible to all agents in the system",
  value: "Current focus: Task planning\nActive agents: 5",
  limit: 10000
});

// Attach to Legacy AI (mother agent)
const legacyAI = await client.agents.create({
  memory_blocks: [...],
  block_ids: [sharedState.id]
});

// Attach SAME block to sub-agent
const financeAgent = await client.agents.create({
  memory_blocks: [...],
  block_ids: [sharedState.id]  // Same block ID!
});
```

### What This Enables
- **Instant context sharing** - When Finance Agent updates the block, Legacy AI sees it next turn
- **Coordination primitive** - Agents can communicate state through shared blocks
- **Read-only policies** - Share rules/guidelines that sub-agents can't modify

### Limitation: Not Real-Time

Changes propagate on the **next agent turn**, not instantly mid-conversation. There's no WebSocket-style push notification.

---

## 4. MULTI-AGENT COMMUNICATION - Full Details

### Three Communication Methods

**1. Async Messaging (Fire and Forget)**
```python
send_message_to_agent_async(
    agent_id="finance_agent_123",
    message="Update the Q4 budget analysis"
)
# Returns immediately, doesn't wait for response
```

**2. Sync Messaging (Wait for Reply)**
```python
response = send_message_to_agent_and_wait_for_reply(
    agent_id="finance_agent_123",
    message="What's the current budget status?"
)
# Blocks until Finance Agent responds
```

**3. Group Broadcast (Supervisor Pattern)**
```python
responses = send_message_to_agents_matching_all_tags(
    tags=["sub-agent", "active"],
    message="Report your current status to Legacy"
)
# Returns list of responses from all matching agents
```

### CRITICAL RECOMMENDATION FROM DOCS:

> "Attach only one messaging tool to prevent agent confusion"

Don't give an agent both async AND sync messaging - it will get confused about which to use.

---

## 5. COMPACTION - How Context Overflow is Handled

When conversation history exceeds the context window, Letta **automatically summarizes** older messages.

### Default Settings
```python
{
    "mode": "sliding_window",           # Keep recent, summarize old
    "sliding_window_percentage": 0.7,   # Keep 70% recent
    "clip_chars": 2000,                 # Summary max length
    "model": "agent's main model"       # Uses same model by default
}
```

### Customization Options
```python
compaction_settings = {
    "model": "openai/gpt-4o-mini",      # Use cheaper model for summaries
    "mode": "sliding_window",            # or "all" to summarize everything
    "sliding_window_percentage": 0.4,    # More aggressive summarization
    "clip_chars": 3000                   # Longer summaries
}
```

### Two Modes

| Mode | Behavior |
|------|----------|
| `sliding_window` | Keeps recent messages, summarizes older ones |
| `all` | Summarizes the entire conversation history |

### Impact on Your Architecture

If Legacy AI has long conversations, older context gets summarized. This means:
- Recent sub-agent updates stay intact
- Older updates become summaries
- Important persistent info should be in **memory blocks** (not conversation)

---

## 6. ARCHIVAL MEMORY - Unlimited Storage

### Key Differences from Memory Blocks

| Memory Blocks | Archival Memory |
|--------------|-----------------|
| Always in context window | Queried on-demand |
| Counts against context limit | **Unlimited storage** |
| Fast, always available | Requires search query |
| For active working memory | For reference/historical data |

### How Archival Memory Works
```python
# Agent stores information
archival_memory_insert(
    content="User prefers dark mode interfaces and minimalist design",
    tags=["user_preferences", "ui"]
)

# Agent searches later
results = archival_memory_search(
    query="what are user's UI preferences",
    tags=["user_preferences"],
    page=0
)
```

### Semantic Search

> "Searches understand conceptual meaning, not just exact keywords"

You can search for "UI preferences" and find content about "dark mode" because it understands the semantic relationship.

### SDK Access (Programmatic)
```python
# Insert via SDK
client.agents.passages.insert(
    agent_id=agent.id,
    content="Finance report Q4 2024: Revenue up 15%",
    tags=["finance", "quarterly_report", "2024"]
)

# Search via SDK
results = client.agents.passages.search(
    agent_id=agent.id,
    query="Q4 financial performance",
    tags=["finance"]
)
```

---

## 7. SLEEP-TIME AGENTS - Reducing Legacy AI's Latency

### What They Are

**Experimental feature** - Background agents that update shared memory blocks asynchronously.

### How It Works
```typescript
const agent = await client.agents.create({
  enable_sleeptime: true,  // Enables the feature
  memory_blocks: [
    { value: "", label: "human" },
    { value: "You are Legacy AI...", label: "persona" },
    { value: "", label: "learned_context" }  // Sleep-time agent fills this
  ]
});
```

When enabled, Letta automatically creates:
1. **Primary agent** (Legacy AI) - handles conversations
2. **Sleep-time agent** - runs in background every N steps

### Configuration
```typescript
{
  enable_sleeptime: true,
  sleeptime_agent_frequency: 5  // Update every 5 messages (default)
}
```

### Performance Trade-offs

| Frequency | Latency Impact | Token Cost |
|-----------|---------------|------------|
| Every 3 steps | Lower latency, higher cost | High |
| Every 5 steps (default) | Balanced | Medium |
| Every 10 steps | Higher latency, lower cost | Low |

> "Frequent updates can be token-intensive with diminishing returns"

**Recommendation:** Keep frequency at 5-10 for most use cases.

### Sleep-Time Agent Tools
- `memory_rethink` - Reflect on and reorganize memory
- `memory_finish_edits` - Signal completion of memory updates

---

## 8. CONTEXT ENGINEERING - Agent Self-Management

### Built-in Memory Tools (Default)

Every agent gets these tools automatically:
- `memory_insert` - Add content to a memory block
- `memory_replace` - Replace content in a memory block

### Disabling Self-Management
```typescript
const agent = await client.agents.create({
  include_base_tools: false  // Agent can't modify its own memory
});
```

### Cross-Agent Memory Modification

Agents can modify OTHER agents' memory blocks via custom tools:
```python
def update_legacy_context(agent_state, new_info: str):
    """Custom tool that updates Legacy AI's memory"""
    client.blocks.update(
        block_id="legacy_shared_context",
        value=new_info
    )
```

---

## 9. AGENT STATE & PERSISTENCE

### What Gets Persisted
- All memory blocks
- Complete message history
- Agent configuration
- Tool assignments
- Archival memory

### Storage Backend
- **PostgreSQL** database
- Data stored at `~/.letta/.persist/pgdata` (self-hosted)
- Checkpointed after each interaction

### Key Architecture Point

> "No concept of threads, just continuous agent interactions"

Unlike ChatGPT conversations, Letta agents have ONE continuous history. Every message ever sent is part of the same thread (until compacted).

---

## 10. EXPLICIT LIMITATIONS FOUND

### Documented Limitations

| Limitation | Details |
|------------|---------|
| Context window | Bound by underlying LLM's limit |
| Block updates | Last write wins (no merge) |
| Compaction | Loses detail from old messages |
| Sleep-time | Experimental, may be unstable |
| Real-time sync | Not available (next-turn propagation) |

### NOT Documented (Gaps Found)

| Missing Info | Impact |
|--------------|--------|
| Max blocks per agent | Unknown - test empirically |
| Max shared block connections | Unknown |
| Multi-agent message latency | Not specified |
| API rate limits (self-hosted) | Not specified |
| Concurrent block write behavior | Only "last write wins" mentioned |
| Memory block indexing performance | Not specified |

---

## 11. RECOMMENDATIONS FOR YOUR LEGACY AI ARCHITECTURE

### Architecture Pattern: Hub with Shared Context

```
┌─────────────────────────────────────────────────────┐
│                    LEGACY AI                         │
│              (Mother Agent - GPT-4o)                 │
│                                                      │
│  Memory Blocks:                                      │
│  ├── [persona] Legacy's identity (5K chars)         │
│  ├── [user_context] About you (3K chars)            │
│  ├── [global_state] SHARED - system state (10K)     │
│  ├── [agent_reports] SHARED - sub-agent updates     │
│  └── [active_tasks] SHARED - current work (5K)      │
│                                                      │
│  Archival Memory: Historical data, logs, docs       │
└─────────────────────────────────────────────────────┘
          │ Shared Blocks │
          ▼               ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Finance     │    │ Health      │    │ Learning    │
│ Agent       │    │ Agent       │    │ Agent       │
├─────────────┤    ├─────────────┤    ├─────────────┤
│ [persona]   │    │ [persona]   │    │ [persona]   │
│ [domain]    │    │ [domain]    │    │ [domain]    │
│ [SHARED]    │    │ [SHARED]    │    │ [SHARED]    │
│ global_state│    │ global_state│    │ global_state│
│ agent_rpts  │    │ agent_rpts  │    │ agent_rpts  │
└─────────────┘    └─────────────┘    └─────────────┘
```

### Block Strategy for Real-Time Context

| Block | Owner | Shared With | Purpose |
|-------|-------|-------------|---------|
| `legacy_persona` | Legacy | None | Legacy's identity, can't be changed |
| `user_context` | Legacy | All (read-only) | Info about you |
| `global_state` | Legacy | All (read-write) | Current system state |
| `agent_reports` | Sub-agents | Legacy (read) | Sub-agents write status here |
| `active_tasks` | Legacy | All (read-only) | What's being worked on |
| `[domain]_context` | Each sub-agent | Legacy (read) | Domain-specific info |

### Communication Pattern

```python
# Sub-agent reports to Legacy via shared block
def report_to_legacy(agent_state, status: str):
    current = client.blocks.get(block_id="agent_reports").value
    updated = f"{current}\n[{datetime.now()}] Finance: {status}"
    client.blocks.update(block_id="agent_reports", value=updated)

# Legacy queries sub-agent directly when needed
response = send_message_to_agent_and_wait_for_reply(
    agent_id="finance_agent",
    message="Give me detailed Q4 analysis"
)
```

### Context Budget (128K token model)

```
Total Context: ~500,000 characters

Memory Blocks Budget: ~150,000 chars (30%)
├── legacy_persona:     5,000 chars
├── user_context:       3,000 chars
├── global_state:      20,000 chars
├── agent_reports:     50,000 chars (rolling, compacted)
├── active_tasks:       5,000 chars
├── domain summaries:  67,000 chars (across all domains)

Conversation Budget: ~200,000 chars (40%)
├── Recent messages preserved
├── Older messages compacted to summaries

System/Tools Budget: ~150,000 chars (30%)
├── System prompt
├── Tool definitions
├── Response generation
```

### Performance Optimizations

**1. Enable Sleep-Time for Legacy**
```typescript
enable_sleeptime: true,
sleeptime_agent_frequency: 5
```

**2. Use Archival Memory for History**
- Don't bloat memory blocks with historical data
- Store completed tasks, old reports in archival

**3. Implement Report Rotation**
```python
# When agent_reports gets too large
if len(agent_reports.value) > 40000:
    # Archive old reports
    client.agents.passages.insert(
        agent_id=legacy.id,
        content=agent_reports.value,
        tags=["agent_reports", "archived"]
    )
    # Reset block with summary
    client.blocks.update(
        block_id="agent_reports",
        value="[Previous reports archived]\n"
    )
```

**4. Tag Sub-Agents for Broadcast**
```typescript
// Create sub-agents with consistent tags
const financeAgent = await client.agents.create({
    tags: ["sub-agent", "domain-finance", "active"]
});

// Legacy can broadcast to all
send_message_to_agents_matching_all_tags(
    tags: ["sub-agent", "active"],
    message: "Status report required"
)
```

---

## 12. WHAT YOU SHOULD TEST EMPIRICALLY

Since Letta doesn't document these limits, you should test:

1. **Max blocks before performance degrades** - Start with 10, increase to 20, 30...
2. **Shared block update latency** - How fast do changes propagate?
3. **Concurrent write behavior** - What happens with simultaneous updates?
4. **Compaction quality** - Does it preserve critical information?
5. **Sleep-time stability** - It's experimental, monitor for issues

---

## BOTTOM LINE

**Your architecture vision is fully supported by Letta.** The platform was literally designed for this - shared memory blocks enable exactly the "mother agent with extensions" pattern you described.

**Key constraints to watch:**
1. Total memory blocks must fit in context window
2. No real-time push - updates propagate next turn
3. Last write wins on shared blocks
4. Sleep-time agents are experimental

---

## Sources
- Letta Documentation: https://docs.letta.com
- MemGPT Research Paper: https://arxiv.org/abs/2310.08560
- Letta GitHub: https://github.com/letta-ai/letta
