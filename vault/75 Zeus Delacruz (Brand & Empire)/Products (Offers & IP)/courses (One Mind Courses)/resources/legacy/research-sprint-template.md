# Research Sprint Template

> *A focused research sprint: define scope, run it, capture outputs.*
> *Use when you need deep research on a topic in a defined time window.*

---

```yaml
---
type: research-sprint
domain: 06
name: "[Sprint Name]"
agent: Trinity
status: active
started: YYYY-MM-DD
deadline: YYYY-MM-DD
owner: Zeus
---
```

---

## Sprint: [Name]

**Research question:** [The specific question this sprint answers]

**Why this matters now:** [What decision or goal does this research feed?]

**Success criteria:** [How will you know this sprint is complete? What artifact does it produce?]

---

## Scope

**In scope:**
- [Specific subtopic 1]
- [Specific subtopic 2]
- [Specific subtopic 3]

**Out of scope:**
- [What NOT to research — important for staying focused]

**Time box:** [X hours or X days — hard cap]

---

## Research Sources

Define the sources Trinity should use for this sprint:

| Source Type | Specific Sources | Priority |
|-------------|-----------------|----------|
| Primary docs | [Official docs, whitepapers, specs] | High |
| News/current | [Miniflux feeds, specific sites] | Medium |
| Academic | [Arxiv, specific journals] | If needed |
| Community | [GitHub issues, forums, Discord] | Low |
| My vault | [Related notes in SO/06] | First |

---

## Research Brief (Trinity Prompt)

Copy this to Trinity to kick off the sprint:

```
Trinity 🔺, run research sprint: [Sprint Name]

Research question: [FILL]

Scope:
- In: [FILL]
- Out: [FILL]

Use these sources in priority order: [FILL]

Time box: [X hours]

Deliverables:
1. Executive summary (300 words max) answering the research question
2. Key findings (5-10 bullet points — most important discoveries)
3. Evidence gaps (what you couldn't find / needs human research)
4. Recommended next step (one clear action based on findings)
5. Source list (everything you used)

Save output to: SO/06 Research/sprints/[YYYYMMDD]-[sprint-name].md
```

---

## Sprint Output

Fill in after Trinity completes the research:

### Executive Summary

[300 words answering the research question]

---

### Key Findings

1. [Finding 1 — specific and supported by evidence]
2. [Finding 2]
3. [Finding 3]
4. [Finding 4]
5. [Finding 5]

---

### Evidence Gaps

- [What wasn't findable — where is the knowledge incomplete?]
- [What needs primary research vs. secondary sources?]

---

### Recommended Next Step

[One clear action based on the findings. Specific enough to assign and execute.]

---

### Sources

| Source | Type | Reliability | Used For |
|--------|------|-------------|---------|
| [FILL] | | | |

---

## Close

- [ ] Output saved to vault
- [ ] Key findings added to domain 06 knowledge base
- [ ] Any action items created in Todoist
- [ ] Memory updated with findings context

---

*Part of the One Mind OS Intelligence & Research — Course 8*
