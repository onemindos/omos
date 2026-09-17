---
name: pattern-recognition
description: "Identifies recurring patterns across Zeus's life data — behaviors, outcomes, decisions, health metrics, mood, and performance. Use when looking for trends, diagnosing why something keeps happening, spotting cycles in productivity, health, or finances, or doing a retrospective analysis. Activates on: pattern, trend, recurring, cycle, why does this keep happening, what's the pattern, retrospective, analysis, correlation, behavioral loop, habit loop, identify, diagnose."
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Pattern Recognition 🔍 — Skill

Scans across One Mind data to surface recurring patterns — behavioral loops, performance cycles, trigger-response sequences, and outcome predictors. Where Correlation finds links between variables, Pattern Recognition finds the shape of how things unfold over time.

---

## When to Activate

| Signal | Examples |
| ------ | -------- |
| Direct request | "what's the pattern here", "why does this keep happening" |
| Retrospective | Monthly review, quarterly OKR, annual recap |
| Diagnosis mode | "I keep failing at X" — what's the loop? |
| Streak analysis | "my workouts dropped off — when and why?" |
| Decision review | "what patterns do I see in my best/worst decisions?" |

---

## What to Look For

### Behavioral Loops
Identify trigger → behavior → outcome chains that repeat. Three or more occurrences = a pattern worth naming.

| Loop Type | Example |
| --------- | ------- |
| **Productive loop** | Good sleep → morning workout → high focus day → evening satisfaction |
| **Destructive loop** | Late night → skipped workout → low energy → poor decisions → guilt |
| **Avoidance loop** | Stressful situation → distraction → delayed action → compounded stress |
| **Growth loop** | New skill → early success → increased confidence → deeper practice |

### Performance Cycles
Look for rhythms in output, energy, and mood across:
- **Daily** — time-of-day peaks and valleys
- **Weekly** — which days consistently outperform or underperform
- **Monthly** — mid-month slumps, end-of-month sprints
- **Seasonal** — annual cycles in energy, mood, motivation

### Decision Patterns
- Decisions made under pressure vs. clarity — what's the outcome difference?
- Time-of-day or emotional state at decision points
- Which domain (HP/LE/GE/SO) generates the most deferred decisions?
- Recurring decision types that could be systematized

### Outcome Patterns
- What inputs consistently precede your best output days?
- What always precedes a bad week?
- Which tasks/projects complete on time vs. always slip?

---

## Analysis Method

1. **Define the window** — what time period? (last 30 days, last quarter, last year)
2. **Identify the variable** — what are we tracking? (energy, output, mood, finances, fitness)
3. **Pull the data** — vault logs, Todoist completion rates, Oura data, financial records
4. **Look for 3+ occurrences** — one is noise, two is coincidence, three is a pattern
5. **Name it** — give the pattern a label so it's recognizable when it recurs
6. **Prescribe** — if destructive, break it. If productive, protect it and replicate it.

---

## Output Format

When a pattern is identified, report it as:

```
PATTERN: [Name]
Type: [Productive / Destructive / Neutral / Cyclical]
Trigger: [What starts it]
Sequence: [What happens]
Outcome: [Where it ends]
Frequency: [How often / when]
Recommendation: [Protect it / Break it / Monitor it]
```

---

## Data Sources

| Source | What It Reveals |
| ------ | --------------- |
| Vault daily logs | Activity, mood, energy, tasks completed |
| Todoist completion rates | Execution patterns by domain and task type |
| Oura / biometrics | Sleep, HRV, readiness cycles |
| Financial records | Spending patterns, revenue cycles |
| Evening reviews | Self-reported energy, mood, wins |
| Annual recaps (LE 54) | Year-over-year behavioral shifts |

---

## Loaded by

Zenith (HP) or Sovereignty (SO). Triggered when Zeus asks about patterns, trends, cycles, recurring issues, or wants a behavioral analysis.
