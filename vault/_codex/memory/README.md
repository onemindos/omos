# OneMind Memory System

> 3-tier persistent memory for agent continuity across sessions.

## Tiers

| Tier | Purpose | Location |
|------|---------|----------|
| 🔨 **Work** | Active tasks, session state, current projects | `work/` |
| 📚 **Learning** | Feedback signals, patterns, lessons | `learning/` |
| 🤝 **Relationship** | How Zeus works, preferences, trust levels | `relationship/` |

## How Agents Use This

### Session Start
1. Load `work/current-focus.json` for active context
2. Scan recent `learning/signals.jsonl` entries
3. Read `relationship/interaction-notes.md`

### During Session
- Update `work/` with task progress
- Log feedback to `learning/signals.jsonl`

### Session End
- Archive completed work
- Extract lessons → `learning/signals.jsonl`
- Update `relationship/interaction-notes.md` if new preferences discovered

## Consolidation (Weekly)

During × Multiply reviews:
1. Process `signals.jsonl` → update `patterns.md`
2. Archive old signals to `learning/archive/`
3. Review and refine `interaction-notes.md`
