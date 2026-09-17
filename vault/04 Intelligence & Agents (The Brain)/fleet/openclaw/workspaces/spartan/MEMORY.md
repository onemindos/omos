---
type: memory-architecture
agent: spartan
status: active
updated: 2026-05-03
---
# MEMORY.md — Spartan ⚔️

## Memory Architecture

Spartan maintains body-domain memory — training history, nutrition patterns, biometric trends, recovery state, and progression tracking.

## Tiers

### Tier 1: Working Memory (session-scoped)
- Today's training plan and completion status
- Current nutrition targets and intake
- Active recovery protocols
- Pending meal prep or supplement reminders

### Tier 2: Episodic Memory (persistent)
- Training log (exercise, sets, reps, weight, RPE, date)
- Nutrition log (meals, macros, calories, adherence)
- Body composition history (weight, measurements, photos)
- Sleep quality data and patterns
- Injury/recovery history

### Tier 3: Semantic Memory (knowledge base)
- Zeus's training preferences and response to different modalities
- Effective programming patterns (what produces results)
- Nutrition principles that work for Zeus's body type
- Recovery protocols and their effectiveness
- Periodization knowledge and current phase

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Today's plan | OpenClaw workspace memory | Daily |
| Training log | Codex vault HP 28 Body | Permanent |
| Body metrics | Codex vault HP 28 Body | Permanent |
| Program design | Codex vault HP 31 Mastery | Permanent |

## Memory Rules

1. **Track progressive overload** — Every session should reference the last session's numbers
2. **Never lose training data** — A missed log is lost intelligence
3. **Correlate sleep → performance** — Track how recovery affects output
4. **Adapt to feedback** — If Zeus says "that was too much", adjust future programming
5. **Periodize memory** — Know where we are in the training cycle (volume, intensity, deload)

## What Spartan Remembers

- Every workout logged (exercise selection, volume, intensity)
- Personal records (1RM, rep maxes, time PRs)
- Injury history and movement limitations
- What training splits Zeus responds to best
- Nutrition patterns that correlate with good training days
- Sleep/recovery trends and their impact on performance
- Current training phase and when to transition
