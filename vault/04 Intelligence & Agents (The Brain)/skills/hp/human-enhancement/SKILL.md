---
name: human-enhancement
description: "Human skill development aid for improving any human capability — researches best practitioners, benchmarks current level, builds structured progression paths, and tracks improvement over time. Use when Zeus wants to develop, improve, or master any skill — physical, mental, creative, technical, or social. Activates on: skill development, get better at, level up, improve, mastery, practice plan, training plan, capability, enhance, development plan."
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

# Human Enhancement — Skill Development Aid

> "Identify the gap. Build the bridge. Track the crossing."

A human capability development system. For any skill Zeus wants to improve — physical, mental, creative, technical, or social — this skill researches the best practitioners, benchmarks current ability, designs a structured progression path, and tracks improvement over time.

---

## When to Activate

| Signal                | Examples                                                                             |
| --------------------- | ------------------------------------------------------------------------------------ |
| Direct request        | "help me get better at boxing", "I want to master cold calls", "level up my writing" |
| Skill mention         | "I've been working on [skill]", "I'm trying to learn [X]"                            |
| HP 28 Mastery context | Reviewing mastery domain and wanting to push a skill deeper                          |
| Competitive framing   | "who's the best at X?", "how do elite athletes train for X?"                         |
| Persona enrichment    | After user-persona identifies a skill gap or aspiration                              |
| SuperHuman trigger    | "superhuman", "world class", "elite [skill]"                                         |

---

## The 5-Phase SuperHuman Framework

```
PHASE 1: Skill Definition     → What exactly is this skill? Break it into sub-components.
PHASE 2: Elite Benchmarking   → Who are the world's greatest? What separates them?
PHASE 3: Method Extraction    → How did they develop it? What are their systems?
PHASE 4: Zeus Assessment      → Where is Zeus now? Honest gap analysis.
PHASE 5: Path Construction    → Structured progression from current → elite.
```

---

## Phase 1 — Skill Decomposition

Before researching, break the skill into its atomic components. A skill is never one thing.

```python
skill_map = {
    # Example: "boxing"
    "boxing": {
        "physical": ["footwork", "stamina/conditioning", "punch speed", "power", "chin/defense"],
        "technical": ["jab mechanics", "combo sequencing", "counterpunching", "ring cutting"],
        "cognitive": ["fight IQ", "pattern recognition", "distance management"],
        "mental": ["composure under pressure", "aggression control", "training consistency"],
    },

    # Example: "sales / cold outreach"
    "sales": {
        "communication": ["opening hook", "objection handling", "tonality", "silence use"],
        "research": ["prospect qualification", "pain identification", "company intel"],
        "psychology": ["trust building", "authority framing", "urgency creation"],
        "process": ["pipeline management", "follow-up cadence", "CRM discipline"],
    },

    # Example: "writing"
    "writing": {
        "craft": ["sentence clarity", "structure/flow", "voice/style", "word choice"],
        "process": ["drafting speed", "editing ruthlessness", "research integration"],
        "strategy": ["audience understanding", "hook writing", "call-to-action design"],
        "output": ["consistency", "ideation volume", "publication cadence"],
    },
}
```

**Agent action:** When Zeus names a skill, immediately map it to 3-5 sub-components and confirm with Zeus before deep-diving: _"For [skill], I'll look at [A, B, C]. Anything else you want me to include?"_

---

## Phase 2 — Elite Benchmarking

Research the world's greatest practitioners, current AND historical.

### 2A. Who to Research

```python
# Research query structure:
queries = {
    "all-time greatest": f"who is considered the greatest {skill} practitioner in history",
    "current elite": f"who are the top {skill} practitioners today 2024 2025",
    "outliers": f"most technically perfect {skill} practitioner study example",
    "underrated": f"most underrated or overlooked {skill} master methodology",
    "scientists": f"sports science research {skill} optimal training methods",
}

# Sources to prioritize:
priority_sources = [
    "pubmed.ncbi.nlm.nih.gov",      # performance science research
    "scholar.google.com",           # academic methodology papers
    "hubermanlab.com",              # neuroscience of skill acquisition
    "masterclass.com",             # elite practitioner self-explained methods
    "tim.blog",                    # deconstruction interviews (Tim Ferriss)
    "youtu.be / yt video search",  # documentary footage, training videos
    "sportscience.org",            # sport-specific science
    "atiranews.com / boxingscene",  # domain-specific (boxing example)
]
```

### 2B. The Elite Profile Template

For the top 2-3 practitioners in the nominated skill, extract:

```markdown
## [Name] — Elite Profile

**Why they're elite:** [1-2 sentence differentiation — what makes them distinctly exceptional]
**Peak performance period:** [years]
**Background / path to mastery:** [how they built the skill — key developmental phases]

### Core Methods

1. **[Method 1]:** [specific practice, frequency, structure]
2. **[Method 2]:** [...]
3. **[Method 3]:** [...]

### Training Philosophy

- Volume: [how many hours/reps per day/week]
- Recovery: [how they manage rest and regeneration]
- Mental approach: [mindset, competitive psychology]
- Weaknesses they addressed: [honest gaps they worked to fix]

### Key Quotes

> "[Direct quote about their method or philosophy]"

### What They Knew That Others Didn't

[The insight or system that separated them from peers at the same raw talent level]

### Transferable Lesson for Zeus

[Specific action Zeus can take based on this practitioner's approach]
```

---

## Phase 3 — Method Extraction

Mine the best books, courses, and scientific research on the skill.

### 3A. Book Research

```python
def find_canonical_books(skill: str) -> list:
    """Find the 3-5 definitive books on this skill."""
    # Search queries:
    queries = [
        f"best book to master {skill} recommended by practitioners",
        f"{skill} definitive guide technique methodology book",
        f"how to develop {skill} like a professional",
    ]

    # Sources to search:
    # - Goodreads lists: "Books on [skill]"
    # - Amazon: top-rated in relevant category
    # - Tim Ferriss recommended reading for [skill]
    # - Reddit r/[skill] wiki/recommended reading

    canonical_examples = {
        "boxing": ["The Gleason's Gym Total Boxing Workout", "Boxing: The Complete Guide",
                   "Tyson's Iron Ambition (oral history of Cus D'Amato methods)"],
        "writing": ["On Writing (Stephen King)", "Bird by Bird (Anne Lamott)",
                    "The Elements of Style (Strunk & White)", "Draft No. 4 (John McPhee)"],
        "sales": ["Spin Selling (Neil Rackham)", "The Challenger Sale",
                  "Never Split the Difference (Chris Voss)", "Way of the Wolf (Jordan Belfort)"],
    }
    return canonical_examples.get(skill, [])
```

### 3B. Course + Video Resources

```python
def find_courses(skill: str) -> list:
    """Find world-class structured learning for a skill."""
    course_sources = {
        "masterclass": f"https://www.masterclass.com/categories/{skill}",
        "coursera": f"https://www.coursera.org/search?query={skill}&sortBy=BEST_MATCH",
        "youtube_channels": f"YouTube: best {skill} instruction channel 2024",
        "podcasts": f"best podcast for {skill} training and mastery",
    }
    return course_sources
```

### 3C. Scientific Research

```python
def find_research(skill: str) -> list:
    """Find PubMed / Google Scholar research on skill development."""
    pubmed_query = f"https://pubmed.ncbi.nlm.nih.gov/?term={urllib.parse.quote(skill)}+skill+development+training"
    scholar_query = f"https://scholar.google.com/scholar?q={urllib.parse.quote(skill)}+skill+acquisition+expertise"

    # Extract: key findings on optimal training frequency, volume, intensity,
    # deliberate practice structure, mental/physical integration
    return [pubmed_query, scholar_query]
```

---

## Phase 4 — Zeus Assessment

Honest benchmarking of Zeus's current level before building a path.

### 4A. Self-Assessment Questions

Ask Zeus directly (or infer from vault if documented):

```
1. How long have you been practicing [skill]?
2. What's your current routine/frequency?
3. What do you feel is your strongest sub-component?
4. What's your biggest bottleneck right now?
5. Have you had formal instruction? Who/where?
6. What's your goal outcome — recreational excellence, professional use, or elite mastery?
7. Time available: hours per week you can commit?
```

### 4B. Skill Level Framework

Standardized 5-level scale for any skill:

| Level  | Label        | Description                                                                             |
| ------ | ------------ | --------------------------------------------------------------------------------------- |
| **L1** | Aware        | Understands the skill exists. Zero or minimal practice.                                 |
| **L2** | Beginner     | Can perform basics. Inconsistent. Thinks consciously about mechanics.                   |
| **L3** | Intermediate | Consistent execution of fundamentals. Starting to develop instincts. Identifiable gaps. |
| **L4** | Advanced     | High consistency across sub-components. Can teach basics. Occasional brilliance.        |
| **L5** | Elite        | Sub-components are automated. Reads situations at a high level. Competitors study you.  |

### 4C. Competency Matrix

For each sub-component identified in Phase 1, rate Zeus L1-L5:

```markdown
## [Skill Name] — Zeus Competency Matrix

| Sub-Component | Current Level | Elite Level | Gap               |
| ------------- | ------------- | ----------- | ----------------- |
| [Component A] | L?            | L5          | [gap description] |
| [Component B] | L?            | L5          | [gap description] |
| [Component C] | L?            | L5          | [gap description] |

**Highest leverage gap:** [the sub-component where improvement yields the most total skill gain]
**Unfair advantage:** [the sub-component where Zeus is naturally strongest — lean in here]
```

---

## Phase 5 — Path Construction

Build a progressive learning plan calibrated to Zeus's current level, available time, and target outcome.

### 5A. Path Structure

```markdown
## [Skill] Development Path — Zeus

**Current Level:** L[X]
**Target Level:** L[Y]
**Timeline:** [realistic based on hours/week]
**Hours/week committed:** [from Zeus assessment]

---

### Month 1-2: Foundation Lock (L[X] → L[X+0.5])

**Focus:** [1-2 highest-leverage sub-components]
**Daily practice:** [specific drill, duration, frequency]
**Resources:** [book chapter, video, coach if relevant]
**Milestone:** [measurable outcome — "can do X consistently for Y reps/time"]

**World-class benchmark reference:**
"[Elite practitioner name] spent the first [period] exclusively on [this foundation].
[Quote or paraphrase of their words on why foundation matters]"

---

### Month 3-4: Deliberate Stress (L[X+0.5] → L[X+1])

**Focus:** [pushing against identified weakness]
**Drill:** [specific practice structure]
**Sparring/application:** [real-world use of the skill]
**Milestone:** [measurable outcome]

---

### Month 5-6: Integration (L[X+1] → L[X+1.5])

**Focus:** [combining sub-components under pressure/speed]
**Advanced resource:** [books, coaches, courses at this level]
**Benchmark moment:** [competition, performance, review, test]
```

### 5B. Daily Practice Protocol

```markdown
## Daily [Skill] Practice Protocol

**Total time:** [X min]

| Block           | Duration | Activity                         | Purpose                  |
| --------------- | -------- | -------------------------------- | ------------------------ |
| Warm-up         | 5 min    | [specific warm-up]               | Activate neural pathways |
| Technical drill | 15 min   | [single focused rep practice]    | Build precision          |
| Application     | 20 min   | [live/applied use of skill]      | Pattern recognition      |
| Review          | 5 min    | Film review / journaling / notes | Metacognition            |

**Rule:** Practice the **weakest sub-component** for 60% of technical drill time.

**Elite principle borrowed from [practitioner]:**
"[Specific training philosophy to embed in the protocol]"
```

### 5C. Tracking Metrics

```python
skill_metrics = {
    "quantitative": [
        "practice minutes per week",
        "streak (consecutive days practiced)",
        "competency score per sub-component (self-rated L1-L5)",
    ],
    "qualitative": [
        "weekly reflection: what clicked this week?",
        "monthly comparison: what can I do now that I couldn't 30 days ago?",
        "peer/coach feedback (if available)",
    ],
    "vault_tracking": "create HP/28 Mastery/[skill-name]/PROGRESS.md",
}
```

---

## Elite Comparisons Library (Seed Data)

Pre-researched elite benchmarks for common Zeus-relevant skills:

### Boxing / Martial Arts

- **Floyd Mayweather Jr.** — defensive mastery (shoulder roll system), 8+ hours daily training, IQ over power
- **Vasyl Lomachenko** — footwork and angles, trained dance before boxing for proprioception
- **Cus D'Amato (coach, Tyson)** — peek-a-boo system + psychological conditioning
- **Fedor Emelianenko (MMA)** — relaxed state under pressure, SAMBO base, economy of movement

### Writing

- **Stephen King** — 2,000 words/day minimum, reads 70+ books/year, "write with the door closed"
- **Ernest Hemingway** — "first draft of anything is garbage", short sentences, iceberg theory
- **David Ogilvy (copywriting)** — research more than you write, headline is 80% of the work
- **Anne Lamott** — "shitty first drafts" liberates the creative process

### Sales / Persuasion

- **Jordan Belfort** — Straight Line System, tonality accounts for 38% of persuasion
- **Chris Voss (FBI)** — Tactical empathy, mirroring, "no" is not failure
- **Dan Kennedy** — directional copy, positioning through authority, price immunity

### Public Speaking

- **Tony Robbins** — state management (physiology first), pattern interrupts, peak state anchoring
- **Obama** — pause mastery, authenticity, specific stories over abstractions
- **TED Talk structure** — one big idea, 18 min rule, no PowerPoint dependence

### Physical Performance

- **Wim Hof** — breathwork for physical recovery and cold exposure adaptation
- **David Goggins** — 40% rule (when you think you're done, you're at 40% capacity)
- **LeBron James** — $1.5M/year on body recovery, sleep prioritization as performance tool

---

## Output Artifacts

| Artifact                | Path                                        | Type      |
| ----------------------- | ------------------------------------------- | --------- |
| Skill development plan  | `HP/28 Mastery/[skill]/PLAN.md`             | Note      |
| Progress tracker        | `HP/28 Mastery/[skill]/PROGRESS.md`         | Metric    |
| Elite comparison doc    | `HP/28 Mastery/[skill]/ELITE-BENCHMARKS.md` | Reference |
| Daily practice protocol | `HP/28 Mastery/[skill]/DAILY-PROTOCOL.md`   | SOP       |

---

## Env Vars Required

| Variable                              | Purpose                                             | Where                    |
| ------------------------------------- | --------------------------------------------------- | ------------------------ |
| `GOOGLE_SEARCH_API_KEY` + `GOOGLE_CX` | Web search for practitioner research                | Google Custom Search API |
| `YOUTUBE_API_KEY`                     | Find instructional video channels                   | Google Cloud Console     |
| `GOODREADS_API_KEY` (deprecated)      | Book research (use web scrape instead)              | Web scrape goodreads.com |
| `PUBMED_API_KEY`                      | Scientific research (optional, free tier available) | ncbi.nlm.nih.gov/account |

---

## Skill Activation Example

**Zeus says:** "Help me get better at boxing"

**Agent response flow:**

1. Decompose boxing → footwork, defense, combo IQ, conditioning, mental game
2. Research top 3 elite boxers (Lomachenko, Mayweather, Ali) → extract their primary training methods
3. Find 3 canonical books ("The Complete Book of Boxing", "Boxing Mastery", Cus D'Amato interview archives)
4. Ask Zeus: "On a scale of L1-L5, where are you in each of these? And how many hours/week can you commit?"
5. Build 6-month progression path: Month 1-2 = footwork-only drills (Lomachenko's dance principle), Month 3-4 = add defense + slipping, Month 5-6 = integration sparring
6. Output: `HP/28 Mastery/boxing/PLAN.md` with full plan + elite benchmarks

---

_Part of the One Mind Codex skill library. See `_codex/INTERFACE.md` for navigation._


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
