---
name: correlation
description: "Personal analytics engine — finds statistical correlations across all of Zeus's life data. Mirrors exist.io's full attribute system. Use when asking 'what affects my X', 'does Y improve when I do Z', 'what predicts a good/bad day', 'run a correlation', or 'what's connected to my mood/sleep/energy/productivity'. Replaces exist.io. Activates on: correlation, what affects, what causes, what predicts, exist.io, personal analytics, data analysis, input-output, what's connected, when I do X does Y change, mood, sleep, energy, steps, HRV, stress, tasks, weight, caffeine, alcohol, workouts."
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Correlation 📊 — Skill

Personal analytics engine. Full replacement for exist.io. Finds statistically significant relationships across Zeus's life data using exist.io's same attribute structure — 12 categories, every attribute correlated against every other.

> "Exist.io found correlations in your data. This skill finds them in *your vault* — with your context, on your schedule, owned by you."

---

## When to Activate

| Signal | Examples |
|--------|---------|
| Direct request | "what affects my sleep?", "does exercise improve my mood?" |
| Diagnostic | "why are my energy levels low this week?" |
| Optimization | "what inputs predict my highest output days?" |
| exist.io replacement | "run a correlation analysis on my data" |
| Review sessions | Monthly / quarterly — what's actually moving the needle? |
| Hypothesis testing | "I think caffeine is hurting my sleep — confirm or rule out" |

---

## Attribute Library (Exist.io Structure)

All attributes tracked by Zeus, mapped to One Mind data sources. Every attribute can be correlated against any other.

### 😊 Mood
| Attribute | Scale | Source |
|-----------|-------|--------|
| Mood | 1–9 | Evening review / manual |
| Energy level | 1–9 | Evening review / manual |
| Stress level | 1–9 | Evening review / manual |
| Daily note | Text | Evening review |

### 😴 Sleep
| Attribute | Unit | Source |
|-----------|------|--------|
| Time asleep | hours | Oura Ring |
| Time in bed | hours | Oura Ring |
| Deep sleep | hours | Oura Ring |
| REM sleep | hours | Oura Ring |
| Light sleep | hours | Oura Ring |
| Bedtime | time | Oura Ring |
| Wake time | time | Oura Ring |
| Awakenings | count | Oura Ring |
| Sleep score | 0–100 | Oura Ring |

### 🏃 Activity
| Attribute | Unit | Source |
|-----------|------|--------|
| Steps | count | Apple Health / Oura |
| Active minutes | minutes | Apple Health |
| Active energy burned | kcal | Apple Health |
| Stand hours | hours | Apple Watch |

### 💪 Workouts
| Attribute | Unit | Source |
|-----------|------|--------|
| Workouts | count | Strava / Spartan logs |
| Workout duration | minutes | Strava |
| Workout distance | km | Strava |
| Workout type | category | Strava / manual |

### 🧠 Productivity
| Attribute | Unit | Source |
|-----------|------|--------|
| Tasks completed | count | Todoist |
| Words written | count | Vault / manual |
| Productive time | hours | manual / Toggl |
| Distracting time | hours | manual |
| Mobile screen time | hours | iOS Screen Time |
| Commits | count | GitHub |
| Emails received | count | manual |
| Emails sent | count | manual |

### 📅 Events
| Attribute | Unit | Source |
|-----------|------|--------|
| Calendar events | count | Calendar |
| Time in meetings | hours | Calendar |

### 💰 Finance
| Attribute | Unit | Source |
|-----------|------|--------|
| Money spent | $ | Finance logs / manual |
| Revenue | $ | Finance logs / manual |

### 🍽️ Food & Drink
| Attribute | Unit | Source |
|-----------|------|--------|
| Calories in | kcal | Manual / MyFitnessPal |
| Protein | g | Manual |
| Carbohydrates | g | Manual |
| Fat | g | Manual |
| Sugar | g | Manual |
| Fibre | g | Manual |
| Caffeine | mg | Manual |
| Coffees | count | Manual |
| Water | ml | Manual |
| Alcoholic drinks | count | Manual |

### 🫀 Health & Body
| Attribute | Unit | Source |
|-----------|------|--------|
| Weight | kg/lbs | Manual / Withings |
| Body fat % | % | Manual |
| Resting heart rate | bpm | Oura |
| HRV | ms | Oura |
| Readiness score | 0–100 | Oura |
| Meditation time | minutes | Manual |
| Average heart rate | bpm | Apple Watch |

### 🎵 Media
| Attribute | Unit | Source |
|-----------|------|--------|
| Articles read | count | Manual / Instapaper |
| Words read | count | Manual |
| Pages read | count | Manual |
| Tracks played | count | Spotify / Last.fm |
| Time watching TV | hours | Manual |

### 🌦️ Weather
| Attribute | Unit | Source |
|-----------|------|--------|
| Max temp | °F/°C | Apple Weather |
| Min temp | °F/°C | Apple Weather |
| Precipitation | mm | Apple Weather |
| Humidity | % | Apple Weather |
| Day length | hours | Apple Weather |

### 🏷️ Custom Tags (Binary — happened or didn't)
| Tag | What it marks |
|-----|--------------|
| `alcohol` | Had a drink today |
| `late-night` | Stayed up past 1am |
| `travel` | Away from home |
| `social` | Meaningful social interaction |
| `no-phone-morning` | No phone first hour |
| `morning-routine` | Full morning routine completed |
| `sunlight` | 20+ min direct sunlight |
| `cold-exposure` | Cold shower / plunge |
| *(add custom tags as needed)* | — |

---

## How Correlation Analysis Works

Mirrors exist.io's engine exactly:

1. **Pair attributes** — every attribute is tested against every other (e.g. "steps vs mood", "alcohol vs sleep", "HRV vs tasks completed")
2. **Set minimum sample size** — 30+ data points for statistical relevance
3. **Calculate correlation coefficient** — r value: direction (positive/negative) + strength
4. **Apply lag testing** — some effects appear next-day or 2–3 days later (e.g. alcohol → HRV impact is +1 day)
5. **Filter for significance** — surface only relationships with r > 0.25 or strong qualitative signal
6. **Rank by strength** — strongest correlations first, with prescription

---

## Correlation Strength Scale

| r value | Strength | Meaning |
|---------|---------|---------|
| < 0.20 | Negligible | No meaningful relationship |
| 0.20–0.39 | Weak | Slight tendency, worth monitoring |
| 0.40–0.59 | Moderate | Real relationship, actionable |
| 0.60–0.79 | Strong | Consistent relationship, high confidence |
| > 0.80 | Very strong | Near-certain relationship |

---

## Output Format

For each finding:

```
CORRELATION: [Attribute A] → [Attribute B]
Direction: Positive ↑ / Negative ↓
Strength: Negligible / Weak / Moderate / Strong / Very Strong (r = X.XX)
Lag: Same-day / +1 day / +2–3 days
Finding: "On days Zeus [A is high/low], [B] is [N]% higher/lower"
Sample size: N data points over [time period]
Confidence: Low / Medium / High
Prescription: [Specific action — protect it, break it, test it further]
```

---

## Standard Correlation Queries

Run these on demand or at monthly/quarterly reviews:

### Weekly Snapshot
- Top 3 positive correlations for mood this week
- What hurt sleep most this week?
- Energy level trend vs. workout days

### Monthly Analysis
- What are the top 5 inputs that predict Zeus's best output days?
- What is the #1 thing that tanks sleep quality?
- Alcohol → HRV lag analysis
- Stress level → spending correlation
- Exercise frequency → mood over 7-day rolling window

### Quarterly Deep Dive
- Full attribute × attribute correlation matrix
- Which day of week shows highest / lowest performance on each key metric?
- What inputs precede Zeus's best creative/strategic work?
- What correlates most strongly with financial discipline?
- Year-over-year: are correlations shifting? (e.g. caffeine sensitivity changing)

### Hypothesis Testing
When Zeus has a hunch — run it:
- "Test: cold exposure → next-day HRV"
- "Test: no alcohol → sleep quality improvement"
- "Test: morning sunlight → energy level"
- "Confirm or rule out: late screen time → worse sleep"

---

## Loaded by

Zenith (HP) or Sovereignty (SO).

**Triggers:** correlation, exist.io, what affects, what causes, what predicts, personal analytics, data analysis, mood analysis, sleep analysis, when I do X does Y change, run an analysis, what's connected.

---

*Replaces: exist.io ($6.99/month)*
*Data stays in: One Mind vault — private, owned, no subscription*
