---
id: skill-direct
title: "Skill: Direct"
type: skill
status: active
created: 2026-03-27
modified: 2026-03-27
owner: legacy
tags: [type/skill, domain/codex, stage/direct]
---

# Skill: Direct

> Prioritize captured and organized items, assign them to the right agent, and set execution order.

## CODEX Stage

**D — Direct**: The planning/prioritization step between Organize and Execute.

## Purpose

After content is captured (C) and organized (O), it sits in the right domain but has no priority, no owner, and no execution plan. Direct fixes that. This skill evaluates queued items against active goals, assigns work to the correct agent, and determines what gets done first.

## Triggers

- New items have been organized but lack priority or assignment
- Active goals have changed and existing items need re-ranking
- Agent asks "what should I do next?"
- Weekly or daily planning sessions
- Backlog exceeds threshold and needs triage

## Workflows

### 1. Priority Assessment

```
INPUT: Unranked items in any domain
PROCESS:
  1. Read `_codex/system/active-goals.md`
  2. Score each item against active goals (high/medium/low alignment)
  3. Factor in urgency (deadline proximity) and effort (quick win vs. deep work)
  4. Assign priority: P1 (do now), P2 (do soon), P3 (do later), P4 (delegate/defer)
OUTPUT: Items tagged with priority level
```

### 2. Agent Assignment

```
INPUT: Prioritized items
PROCESS:
  1. Identify which quadrant the item belongs to
  2. Map to responsible agent:
     - UI (00-24) → Sentinel (UI Mother AI)
     - HP (25-49) → Coach (HP Agent)
     - LE (50-74) → Haven (LE Agent)
     - GE (75-99) → Empire (GE Agent)
  3. Check if the item requires cross-quadrant coordination
  4. If cross-quadrant: Sentinel coordinates, specialist executes
OUTPUT: Items assigned to specific agent with clear ownership
```

### 3. Execution Queue

```
INPUT: Prioritized, assigned items
PROCESS:
  1. Group by agent
  2. Order by priority within each agent's queue
  3. Identify dependencies (item B requires item A first)
  4. Set execution windows (morning deep work, afternoon admin, etc.)
  5. Cap each agent's active queue at 3-5 items
OUTPUT: Per-agent execution queue, ready for Execute stage
```

### 4. Goal Alignment Check

```
INPUT: Current queues + active-goals.md
PROCESS:
  1. Compare time allocation across quadrants
  2. Flag if >60% of work is in one quadrant (imbalance)
  3. Flag if P1 items don't map to any active goal (drift)
  4. Recommend rebalancing if needed
OUTPUT: Alignment report with adjustment recommendations
```

## Decision Framework

| Factor | Weight | How to Evaluate |
|--------|--------|-----------------|
| Goal alignment | High | Does this directly advance an active goal? |
| Urgency | High | Is there a hard deadline within 48 hours? |
| Effort | Medium | Quick win (<15 min) vs. deep work (>1 hour)? |
| Dependencies | Medium | Does anything else block on this? |
| Quadrant balance | Low | Is this quadrant already overloaded? |

## Rules

- Never queue more than 5 P1 items per agent at once
- Always check active-goals.md before assigning priority — stale goals create stale queues
- Cross-quadrant items default to Sentinel for coordination
- Items with no goal alignment drop to P3 or P4 automatically
- Re-run priority assessment when active goals change
