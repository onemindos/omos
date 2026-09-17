# Codex Skill: Report Generation

How to generate structured reports from vault data. Reports pull from domains, memory, git history, and frontmatter to create actionable intelligence.

---

## Report Types

### 1. Quarterly Life Review

**Trigger**: End of quarter, or "generate quarterly review"
**Scope**: All 4 domain groups

**Steps**:
1. Read `_codex/system/active-goals.md` for the quarter's priorities
2. Scan each domain group for recent activity (git log, file dates)
3. Read `_codex/memory/learning/patterns.md` for consolidated insights
4. Read `_codex/memory/learning/signals.jsonl` for the quarter's feedback
5. For each domain group, assess:
   - Health score (🟢 active / 🟡 moderate / 🔴 neglected)
   - Trend (↑ improving / → steady / ↓ declining)
   - Key win and key gap
6. Generate report using `templates/quarterly-review.md`
7. Save to `00-24 SO/00 Framework (Doctrine)/reviews/`

### 2. Domain Health Report

**Trigger**: Monthly, or "check domain health for [domain]"
**Scope**: Single domain or domain group

**Steps**:
1. Count files in the domain
2. Check last modified dates
3. Look for files with `status: gap` or `status: stale` tags
4. Count cross-domain links (backlinks)
5. Check for active projects
6. Generate report using `templates/domain-health.md`

### 3. Family Legacy Report

**Trigger**: Annually, or "generate legacy report"
**Scope**: LE domains (50-74) + relevant GE domains

**Steps**:
1. Read `_codex/system/mission.md` and `beliefs.md` for family context
2. Scan all LE domains for milestones, changes, new files
3. Pull financial data from GE domain 82 (Investments) and HP domain 30 (Finance)
4. Assess family system health (meetings, calendars, SOPs)
5. Generate report using `templates/family-legacy.md`

### 4. Narrative Points (McKinsey-Style)

**Trigger**: "Generate narrative points for [topic]"
**Scope**: Any topic

**Steps**:
1. Gather all relevant context on the topic
2. Structure into 8-section narrative arc:
   - **Opening** — Hook
   - **Bad News** — The problem
   - **Evidence** — Data supporting the problem
   - **Stakes** — What happens if nothing changes
   - **Pivot** — The solution
   - **Requirements** — What it takes
   - **Vision** — The future
   - **Close** — Call to action
3. Generate concise bullet points for each section
4. Output as presentation-ready content

---

## Data Sources

| Source | What It Provides |
|--------|-----------------|
| `_codex/system/` | Identity, goals, mission, beliefs |
| `_codex/memory/` | Signals, patterns, work history |
| Domain folders | File counts, content, status tags |
| Git history | Commit frequency, recent changes |
| Frontmatter | Status, dates, tags, ai_summary |


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
