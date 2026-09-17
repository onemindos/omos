---
name: ikigai-purpose-discovery
description: "Everyday purpose alignment skill — checks daily work against ikigai map, detects drift from passion/mission/vocation, recommends course corrections, and runs structured discovery sessions when needed. Use when Zeus needs alignment check, feels off-track, is choosing between priorities, wants purpose-based decision support, or needs a full ikigai discovery session. Activates on: ikigai, purpose, alignment, what should I focus on, am I on the right path, purpose check, priorities, drift, meaning, mission alignment."
allowed-tools:
  - Bash(curl:*)
  - Bash(python3:*)
  - Read
  - Write
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Ikigai — Purpose Alignment Skill

> "Your Ikigai isn't a destination you discover once. It's a compass you check daily."

Two modes: **everyday alignment** (check work against purpose) and **deep discovery** (map the full ikigai from scratch).

---

## When to Activate

| Signal | Mode | Examples |
|--------|------|----------|
| Quick alignment check | Everyday | "am I on track?", "purpose check", "is this aligned?" |
| Priority decision | Everyday | "should I do X or Y?", "what should I focus on today?" |
| Drift detection | Everyday | "I feel off", "something's not right", "I'm scattered" |
| Full discovery | Discovery | "run ikigai", "help me find my purpose", "ikigai session" |
| Annual/quarterly review | Discovery | During life planning sessions |

---

## Mode 1: Everyday Alignment (Default)

### Quick Alignment Check (2 min)
When Zeus asks for a purpose check or seems adrift:

1. **Read** the latest ikigai map from `25-49 HP/25 Identity/`
2. **Ask**: "What did you spend most of your time on today/this week?"
3. **Map** the activity against the four circles:
   - Does it touch what you love? ✅/❌
   - Does it use what you're good at? ✅/❌
   - Does it serve what the world needs? ✅/❌
   - Does it generate or lead to income? ✅/❌
4. **Score**: Count ✅s. 4/4 = deep alignment. 2/4 = partial. 0-1/4 = drift alert.
5. **Recommend**: If <3/4, suggest one concrete adjustment to bring the day closer to center.

### Priority Filter
When Zeus has competing priorities:
1. Score each option against the 4 circles (0-4)
2. The option scoring highest is most ikigai-aligned
3. Present the scoring transparently — Zeus decides, but the data is clear

### Drift Detection
Proactively flag when patterns suggest drift:
- Multiple days in a row scoring <2/4
- Tasks dominated by "paid for" but missing "love" = burnout risk
- Tasks dominated by "love" but missing "paid for" = sustainability risk
- Surface the pattern, not just the symptom

### Weekly Alignment Summary (part of weekly review)
```
🎯 Ikigai Alignment — Week of {date}
Alignment score: X/4 average
Strongest circle: {what you're good at}
Weakest circle: {what the world needs}
Drift alert: {yes/no — pattern description}
Recommendation: {one adjustment for next week}
```

---

## Mode 2: Deep Discovery Session

Run when no ikigai map exists, or during annual/quarterly reviews.

```markdown
---
title: "Ikigai Discovery — {date}"
type: review
status: active
created: {date}
owner: zeus
tags: [type/review, domain/hp-25, topic/ikigai, topic/purpose]
ai_summary: "Ikigai discovery session mapping passion, skills, needs, and market fit"
---

# Ikigai Discovery — {date}

## Context
{What prompted this session, current life season}

## The Four Circles

### 🔥 What I Love
- {bullet list of responses}

### 💪 What I'm Good At
- {bullet list of responses}

### 🌍 What the World Needs
- {bullet list of responses}

### 💰 What I Can Be Paid For
- {bullet list of responses}

## Intersections

### Passion (Love + Good At)
{synthesis}

### Mission (Love + World Needs)
{synthesis}

### Profession (Good At + Paid For)
{synthesis}

### Vocation (World Needs + Paid For)
{synthesis}

## 🎯 IKIGAI

{The core insight — where all four overlap}

## Action Items
- [ ] {concrete next step 1}
- [ ] {concrete next step 2}
- [ ] {concrete next step 3}

## Reflection
{Agent's honest assessment of alignment gaps and opportunities}
```

---

## Integration Points

- **Vault**: Save note to HP 25 Identity
- **Todoist**: Optionally create action items as tasks in relevant projects
- **Life Score**: Ikigai alignment feeds into PI (Performance Intelligence) score
- **Annual Review**: Reference during yearly planning and quarterly OKR sessions

---

## Agent Behavior

- Be genuinely curious, not formulaic
- Push deeper when answers are surface-level ("Tell me more about that...")
- Notice contradictions between circles and surface them gently
- Don't force an answer if there's no clear IKIGAI overlap — sometimes the insight is "you're not there yet, and here's what's missing"
- Keep the session to 15-25 minutes maximum
- Always end with concrete action items, not just abstract insights


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
