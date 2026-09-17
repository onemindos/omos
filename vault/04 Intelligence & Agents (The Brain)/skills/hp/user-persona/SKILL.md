---
name: user-persona
description: "Builds and maintains a living psychographic + lifestyle profile for Zeus — zodiac, personality frameworks, food preferences, local restaurant discovery, and behavioral patterns — then enriches USER.md and the vault. Use when Zeus says 'build my persona', 'research my profile', 'what does my zodiac say', 'find pasta restaurants', 'enrich my context', or when agents notice a knowledge gap."
allowed-tools:
  - Bash(curl:*)
  - Bash(python3:*)
  - Bash(jq:*)
  - Read
  - Write
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# User Persona Development

> "The agent that knows you best, serves you best."

This skill builds and maintains a living psychographic + lifestyle profile for Zeus.
Unlike a static form, it actively researches, cross-references, and enriches context
over time — turning surface-level preferences into deep behavioral intelligence.

---

## When to Activate

| Signal                | Examples                                                                         |
| --------------------- | -------------------------------------------------------------------------------- |
| Direct request        | "build my persona", "enrich my user profile", "what does my zodiac say about me" |
| Preference mentioned  | "I like pasta", "I'm an Aries", "I hate mornings", "I prefer dark roast"         |
| Local discovery       | "find me restaurants near Greenwich", "best spots to eat in NYC"                 |
| Gap in context        | Agent says "I don't know your food preferences" or "what's your birthday?"       |
| Periodic refresh      | Monthly background cron trigger                                                  |
| New domain onboarding | Setting up HP 25 Identity or HP 26 Mind for the first time                       |

---

## Persona Architecture

The full persona lives at:

- **`_codex/system/USER-PERSONA.md`** — rich enriched profile (agent context file)
- **`_codex/system/USER.md`** — quick reference identity card (shorter, always loaded)
- **HP 25 Identity / HP 26 Mind** — domain-level depth (goals, beliefs, traits)

---

## Phase 1 — Context Harvest

Read existing context before researching. Pull from:

```python
import os, subprocess

VAULT = "/Users/zeuslegacy/Documents/onemind-codex"
VPS_SHARED = "/home/onemind/onemind-os/openclaw/shared"

sources = [
    f"{VAULT}/_codex/system/USER.md",               # VPS agent identity card
    f"{VAULT}/CLAUDE.md",                             # Session history + known facts
    f"{VAULT}/AGENTS.md",                             # Known about Zeus
    f"{VAULT}/_codex/memory/work/",                  # Recent active memory
]

# Also check VPS shared USER.md if available
# scp or read via ssh: ssh root@104.248.4.212 'cat /home/onemind/onemind-os/openclaw/shared/USER.md'
```

**Extract known facts:**

- Birth date / zodiac sign
- Location (Greenwich CT primary, NYC secondary)
- Food preferences (noted from any prior conversation)
- Physical practices (workout history, sport, diet)
- Cultural background
- Communication style
- Known annoyances / strong preferences

---

## Phase 2 — Zodiac & Archetype Research

When birth date is known, research the full astrological + psychological profile.

### 2A. Zodiac Deep Dive

```python
import urllib.request, json

# Zeus is an Aries (March 21 – April 19)
# Research sources (web fetch):
zodiac_sources = [
    "https://www.astro.com/astrology/in_signs_e.htm",        # detailed sign profiles
    "https://cafeastrology.com/aries.html",                   # deep behavioral patterns
    "https://astrostyle.com/zodiac-signs/aries/",             # modern applied astrology
]

# Key extraction targets:
aries_profile = {
    "core_traits": ["pioneering", "direct", "competitive", "impatient", "passionate"],
    "strengths": ["natural leadership", "courage", "initiative", "energy"],
    "challenges": ["impulsiveness", "short fuse", "difficulty finishing"],
    "work_style": "sprint-focused, dislikes bureaucracy, leads by example",
    "communication": "direct, blunt, hates filler",
    "decision_style": "fast, gut-driven, course-corrects later",
    "motivation": "competition, being first, proving capability",
    "life_challenges": "sustaining long projects, patience, delegation",
    "energy_pattern": "high-burst, needs recovery windows",
    "compatibility_work": ["Leo", "Sagittarius", "Gemini"],
    "shadow": "can dominate conversations, dismisses slow thinkers",
}
```

### 2B. Extended Frameworks (Research if Birth Time Known)

| Framework            | Purpose                        | Source                          |
| -------------------- | ------------------------------ | ------------------------------- |
| **Rising sign**      | How others perceive Zeus       | Requires birth time + location  |
| **Moon sign**        | Emotional processing style     | Requires birth time             |
| **MBTI**             | Cognitive function map         | Ask Zeus or infer from behavior |
| **Enneagram**        | Core motivation + fear         | Ask Zeus or infer               |
| **Human Design**     | Energy type + authority        | Requires birth time + location  |
| **Big Five (OCEAN)** | Research-validated trait model | Infer from history              |

```python
# Example: fetch Aries full profile from cafeastrology
def fetch_zodiac_profile(sign: str) -> dict:
    url = f"https://cafeastrology.com/{sign.lower()}.html"
    req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
    with urllib.request.urlopen(req, timeout=10) as r:
        html = r.read().decode('utf-8', errors='replace')
    # Extract key sections: personality, strengths, weaknesses, love, career, health
    # Return structured dict with extracted text
    return {"sign": sign, "raw_html": html[:5000]}  # trim for LLM processing
```

---

## Phase 3 — Food & Restaurant Intelligence

**Core principle:** User mentions a food preference → enrich with local discovery.

### 3A. Preference Capture

When Zeus mentions a food, log it immediately:

```python
food_prefs = {
    "loves": ["pasta", "dark roast coffee"],  # mentioned explicitly
    "avoids": [],  # capture as mentioned
    "diet_style": "unknown — check HP 27 Fitness",
    "cuisine_types": ["Italian"],  # inferred from pasta preference
}
```

### 3B. Restaurant Discovery — Greenwich CT + NYC

Use Google Maps Places API or Yelp Fusion API:

```python
import urllib.request, json

GOOGLE_MAPS_KEY = os.environ.get("GOOGLE_MAPS_API_KEY")

def find_restaurants(cuisine: str, location: str, radius_m: int = 5000) -> list:
    """Find top-rated restaurants near a location for a given cuisine."""
    # Google Places Text Search
    query = f"best {cuisine} restaurant in {location}"
    encoded_query = urllib.parse.quote(query)
    url = (
        f"https://maps.googleapis.com/maps/api/place/textsearch/json"
        f"?query={encoded_query}&key={GOOGLE_MAPS_KEY}&type=restaurant"
    )
    req = urllib.request.Request(url)
    with urllib.request.urlopen(req, timeout=10) as r:
        data = json.loads(r.read())
    results = data.get("results", [])[:5]
    return [{
        "name": r.get("name"),
        "address": r.get("formatted_address"),
        "rating": r.get("rating"),
        "price_level": r.get("price_level"),  # 1-4 ($–$$$$)
        "user_ratings_total": r.get("user_ratings_total"),
        "maps_url": f"https://www.google.com/maps/place/?q=place_id:{r.get('place_id')}",
    } for r in results]

# Example calls:
greenwich_pasta = find_restaurants("pasta Italian", "Greenwich CT")
nyc_pasta       = find_restaurants("pasta Italian", "New York City")
```

**Yelp alternative (if Yelp API key available):**

```bash
curl -s "https://api.yelp.com/v3/businesses/search?term=pasta+italian&location=Greenwich+CT&sort_by=rating&limit=5" \
  -H "Authorization: Bearer $YELP_API_KEY" | jq '.businesses[] | {name, url, rating, review_count, price}'
```

### 3C. Recipe Discovery

For home cooking preferences:

```python
def find_recipes(ingredient: str, style: str = "family") -> list:
    """Search for top recipes using an ingredient."""
    # Option 1: Spoonacular API
    url = f"https://api.spoonacular.com/recipes/search?query={ingredient}&number=5&apiKey={os.environ.get('SPOONACULAR_API_KEY')}"

    # Option 2: Web search via brave/google
    query = f"best {style} {ingredient} recipe site:seriouseats.com OR site:bonappetit.com OR site:allrecipes.com"

    # Option 3: Scrape NYT Cooking search
    url = f"https://cooking.nytimes.com/search?q={urllib.parse.quote(ingredient)}"

    return []  # return list of {title, url, source, time, difficulty}
```

---

## Phase 4 — Behavioral Pattern Analysis

Cross-reference behavior patterns from the vault to validate/enrich the profile:

```python
patterns_to_analyze = {
    "communication_style": "read CLAUDE.md + USER.md conversation tone section",
    "work_windows": "check HP domain for workout times, daily notes for active hours",
    "decision_velocity": "infer from completed vs in-progress tasks in Todoist",
    "stress_signals": "look for inbox overload, missed reviews, volatile task patterns",
    "energy_peaks": "HP 27 Fitness notes, workout time patterns",
    "learning_style": "HP 28 Mastery notes, course completion patterns",
}
```

---

## Phase 5 — Persona Profile Assembly

Compile everything into a structured `USER-PERSONA.md`:

```markdown
---
title: "Zeus — User Persona Profile"
type: reference
status: active
created: { DATE }
updated: { DATE }
owner: zeus
tags: [type/reference, domain/HP, status/active]
ai_summary:
  "Deep psychographic profile: Aries, direct communicator, sprint-worker,
  Italian food, Greenwich CT + NYC base. Use for personalized suggestions."
---

# Zeus — User Persona Profile

## Identity Core

- **Name:** Zeus Delacruz
- **Archetype:** The Architect — building systems that outlast him
- **Birth sign:** Aries ♈ (fire sign, cardinal)
- **Location base:** Greenwich, CT (home) + New York City (work/social)

## Psychological Profile

### Aries Traits (Validated)

- **Strengths:** Pioneering, direct, high initiative, natural leader
- **Challenges:** Impatience with slow processes, difficulty sustaining long projects
- **Work style:** Sprint-focused, high output in bursts, needs clear targets
- **Communication:** Direct, no filler, rewards honesty
- **Decision style:** Fast and gut-driven, iterates quickly

## Food & Local Intelligence

### Preferences

| Category | Preference                                |
| -------- | ----------------------------------------- |
| Cuisine  | Italian (pasta), [add more as discovered] |
| Coffee   | Dark roast                                |
| Diet     | TBD — check HP 27                         |

### Top Pasta Restaurants — Greenwich CT

{GOOGLE_RESULTS_GREENWICH}

### Top Pasta Restaurants — NYC

{GOOGLE_RESULTS_NYC}

### Saved Recipes

{RECIPE_LIST}

## Energy & Lifestyle Patterns

- **Peak energy:** [from HP workout logs]
- **Low energy:** [from evening review patterns]
- **Communication window:** Direct messages preferred; batch briefings in morning

## Personality Frameworks

| Framework    | Result                    | Source |
| ------------ | ------------------------- | ------ |
| MBTI         | TBD — ask Zeus            |        |
| Enneagram    | TBD — infer from behavior |        |
| Human Design | TBD — requires birth time |        |
| Big Five     | TBD                       |        |

## Agent Personalization Rules

Derived from this profile, agents should:

1. **Skip pleasantries** — get to the point immediately
2. **Lead with the recommendation** — context/reasoning comes after
3. **Use sprint framing** — "in the next 30 min..." not "over the coming weeks..."
4. **Local context first** — Greenwich and NYC for any location-based suggestions
5. **Respect energy cycles** — don't suggest complex tasks at low-energy times
6. **Compete, don't comfort** — Aries responds to challenge more than reassurance
```

---

## Phase 6 — Approval & Save Flow

**NEVER auto-save persona updates.** Always show Zeus first.

```python
# Step 1: Generate the updated persona section
# Step 2: Present a diff or preview to Zeus
# Step 3: Wait for explicit approval ("yes", "save it", "looks good")
# Step 4: Write to vault + sync to VPS shared/USER-PERSONA.md

def save_with_approval(content: str, path: str):
    """
    Present content to user → wait for approval → save.
    For cron/background runs, save to /tmp/persona-draft.md
    and NATS-publish zeus.approval.request.persona for Legacy to surface.
    """
    # Publish to NATS for async approval
    import subprocess
    payload = {"type": "persona_update", "preview": content[:500], "full_path": "/tmp/persona-draft.md"}
    # nats pub zeus.approval.request.persona '{...}'
    pass
```

### Cron Background Mode

When running as a scheduled background job (e.g., monthly):

```bash
# N8N trigger or cron: first of each month at 9am
# 1. Run Phase 1-4 (research only)
# 2. Save draft to /tmp/persona-draft-{DATE}.md
# 3. Publish NATS subject: zeus.approval.request.persona
# 4. Legacy surfaces: "Persona draft ready — 3 new restaurant recs + zodiac update. Approve?"
# 5. Zeus says yes → saved to vault
```

---

## Background Enrichment Triggers

| Event                  | Enrichment Action                                                                           |
| ---------------------- | ------------------------------------------------------------------------------------------- |
| Zeus mentions a food   | → Research top local spots + recipes for that food                                          |
| Zeus mentions a person | → Add to relationship context (LE domain)                                                   |
| Zeus mentions a place  | → Enrich location intelligence                                                              |
| Zeus completes a skill | → Update HP 28 Mastery profile                                                              |
| Monthly cron           | → Full refresh: cross-check all known prefs, new restaurant openings, seasonal menu updates |
| New domain opened      | → Trigger relevant section of persona build                                                 |

---

## Output Artifacts

| Artifact          | Path                                                       | Audience                  |
| ----------------- | ---------------------------------------------------------- | ------------------------- |
| Full persona      | `_codex/system/USER-PERSONA.md`                            | All agents                |
| VPS persona       | `/home/onemind/onemind-os/openclaw/shared/USER-PERSONA.md` | VPS agents (Legacy, etc.) |
| Quick reference   | `_codex/system/USER.md` (update summary section)           | Always-loaded context     |
| Restaurant list   | `HP/25 Identity/local-intelligence/restaurants.md`         | Browsable domain note     |
| Recipe collection | `HP/25 Identity/local-intelligence/recipes.md`             | Browsable domain note     |

---

## Env Vars Required

| Variable              | Purpose                             | Where to Get                        |
| --------------------- | ----------------------------------- | ----------------------------------- |
| `GOOGLE_MAPS_API_KEY` | Places API for restaurant discovery | Google Cloud Console (Maps enabled) |
| `YELP_API_KEY`        | Yelp restaurant search (fallback)   | developer.yelp.com                  |
| `SPOONACULAR_API_KEY` | Recipe discovery                    | spoonacular.com/food-api            |

Store in 1Password → `SO 11 Security / API Keys` vault.

---

## Known Zeus Persona Seeds

Pre-loaded from session history — agents should treat these as confirmed until Zeus updates:

| Fact                | Value                                 | Confidence                  |
| ------------------- | ------------------------------------- | --------------------------- |
| Zodiac sign         | Aries ♈                              | High — mentioned by Zeus    |
| Primary location    | Greenwich, CT (home)                  | High                        |
| Secondary location  | New York City                         | High                        |
| Food: loves         | Pasta (Italian)                       | High — mentioned explicitly |
| Communication style | Direct, no filler, no pleasantries    | High                        |
| Work style          | Sprint-based, high initiative         | High                        |
| Energy              | High-burst, needs recovery            | Medium                      |
| MBTI / Enneagram    | Unknown — ask or infer                | Low                         |
| Birth time          | Unknown — needed for full astro chart | Low                         |

---

_Part of the One Mind Codex skill library. See `_codex/INTERFACE.md` for navigation._


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
