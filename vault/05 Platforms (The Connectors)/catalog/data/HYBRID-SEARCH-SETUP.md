  # BM25 + Vector Hybrid Search — Reference Guide

> **How Legacy searches your entire Codex vault in <2 seconds with zero custom code.**

---

## TL;DR

OpenClaw has **native hybrid search** built into `memory_search`. You don't install anything. You just point it at a folder of markdown files and it indexes them automatically.

```
┌─────────────────────────────────────────────────────────┐
│  YOUR QUERY: "What time does Ash get out of school?"    │
└───────────────────────┬─────────────────────────────────┘
                        │
          ┌─────────────▼──────────────┐
          │     memory_search          │
          │  (OpenClaw native tool)    │
          └─────────────┬──────────────┘
                        │
         ┌──────────────┼──────────────┐
         │                             │
    ┌────▼────┐                  ┌─────▼─────┐
    │  BM25   │                  │  Vector   │
    │ keyword │                  │ semantic  │
    │ match   │                  │ match     │
    └────┬────┘                  └─────┬─────┘
         │ "Ash" + "school"            │ meaning: child + education
         │ exact term hits             │ + schedule
         │                             │
    ┌────▼──────────────────────▼──────┐
    │    Reciprocal Rank Fusion        │
    │    + MMR deduplication           │
    │    + Temporal decay (recent ↑)   │
    └────────────────┬─────────────────┘
                     │
    ┌────────────────▼─────────────────────────────────┐
    │ Result: LE/51 Family/ash-profile.md (score: 0.93)│
    │ "Ash gets out at 3:15pm from Jefferson Elementary"│
    └──────────────────────────────────────────────────┘
```

---

## 1. Where Does the Codex Go?

The Codex vault lives at **`/home/onemind/codex/`** on the VPS/server where OpenClaw runs.

```
/home/onemind/
├── .openclaw/
│   ├── openclaw.json          ← config (has memorySearch.extraPaths)
│   ├── workspace-legacy/      ← SOUL.md, AGENTS.md, TOOLS.md
│   ├── memory/                ← daily logs
│   │   └── legacy/
│   │       └── agent/
│   │           └── memory.sqlite  ← SQLite index (auto-created)
│   └── env                    ← API keys
│
└── codex/                     ← YOUR CODEX VAULT (cloned from GitHub)
    ├── 00-24 SO (Sovereign Operations)/
    ├── 25-49 HP (Holistic Performance)/
    ├── 50-74 LE (Legacy Evolution)/
    ├── 75-99 GE (Generational Empire)/
    └── ONEMIND-CODEX.md
```

### How It Gets There

```
Mac (Obsidian editor)
  → saves note
  → Obsidian Git plugin auto-commits + pushes to GitHub
  → VPS cron pulls every 5 minutes:
    */5 * * * * cd /home/onemind/codex && git pull --ff-only

Mac ──push──→ GitHub ──pull──→ VPS:/home/onemind/codex/
                                    ↓
                              OpenClaw reads + indexes
```

**VPS is read-only.** Only your Mac writes to the Codex. Zero conflicts.

---

## 2. The Configuration (One Line)

In `openclaw.json` under `agents.defaults`:

```json
"memorySearch": {
  "enabled": true,
  "extraPaths": ["/home/onemind/codex"],
  "query": {
    "hybrid": {
      "enabled": true,
      "vectorWeight": 0.7,
      "textWeight": 0.3,
      "candidateMultiplier": 4,
      "mmr": { "enabled": true, "lambda": 0.7 },
      "temporalDecay": { "enabled": true, "halfLifeDays": 30 }
    },
    "cache": { "enabled": true, "maxEntries": 50000 }
  }
}
```

**Critical format note:** `extraPaths` must be **plain strings**:
- ✅ `["/home/onemind/codex"]`
- ❌ `[{"path": "/home/onemind/codex", "pattern": "**/*.md"}]`

### Domain-Scoped Agents

Each agent can be scoped to search only its domain:

| Agent | `extraPaths` | What It Searches |
|-------|-------------|-----------------|
| `legacy` (main) | `["/home/onemind/codex"]` | Entire vault |
| `hp` | `["/home/onemind/codex/25-49 HP (Holistic Performance)"]` | Health, fitness, personal |
| `le` | `["/home/onemind/codex/50-74 LE (Legacy Evolution)"]` | Family, home, legacy |
| `ge` | `["/home/onemind/codex/75-99 GE (Generational Empire)"]` | Business, revenue |

---

## 3. How BM25 Works (Keyword Search)

BM25 (Best Matching 25) is a statistical ranking algorithm — think Google circa 2005.

**For each document, it scores:**
1. **Term Frequency (TF):** How often the query words appear in this doc
2. **Inverse Document Frequency (IDF):** Rare words score higher
   - "NATS" in 3/500 docs → high IDF → very discriminating
   - "the" in 480/500 docs → low IDF → nearly worthless
3. **Document length normalization:** Shorter docs mentioning the term rank higher

**Formula:**
$$\text{BM25}(D, Q) = \sum_{i=1}^{n} \text{IDF}(q_i) \cdot \frac{f(q_i, D) \cdot (k_1 + 1)}{f(q_i, D) + k_1 \cdot (1 - b + b \cdot \frac{|D|}{\text{avgdl}})}$$

Where:
- $f(q_i, D)$ = frequency of term $q_i$ in document $D$
- $|D|$ = document length
- $\text{avgdl}$ = average document length across corpus
- $k_1$ = 1.2, $b$ = 0.75 (standard tuning parameters)

**Strengths:** Fast, exact keyword matching, finds unique identifiers
**Weakness:** "automobile" ≠ "car" — no semantic understanding

---

## 4. How Vector Search Works (Semantic Search)

Text gets converted to dense vectors (arrays of ~1536 numbers) by an embedding model:

```
"Ash has soccer practice Tuesday" → [0.023, -0.187, 0.445, ... 1536 dims]
"Kid's football training midweek" → [0.019, -0.192, 0.438, ...]

Cosine similarity ≈ 0.91 (very similar!)
Even though they share ZERO exact words.
```

**How it works:**
1. Each markdown chunk gets embedded into a vector via the LLM's embedding model
2. Vectors are stored in SQLite (FTS5 + vector columns)
3. Your query is also embedded
4. Cosine similarity finds the closest vectors

**Strengths:** Understands synonyms, paraphrases, conceptual similarity
**Weakness:** Can miss exact unique identifiers like "agent-12847"

---

## 5. How Hybrid Combines Them

Both search methods run in parallel, then results are fused:

```
Query: "how does Legacy handle morning routines"

BM25 results (keyword):                 Vector results (semantic):
1. HEARTBEAT.md (has "morning")          1. HP/habits.md (conceptually about routines)
2. cron-jobs.md (has "morning")          2. HEARTBEAT.md (about scheduled behavior)
3. HP/habits.md (has "routine")          3. proactive-agent.md (automated actions)

Reciprocal Rank Fusion (RRF):
→ HEARTBEAT.md — ranked high in BOTH → #1
→ HP/habits.md — ranked in both → #2
→ cron-jobs.md — BM25 only → #3
→ proactive-agent.md — vector only → #4
```

**RRF formula:**
$$\text{RRF}(d) = \sum_{r \in R} \frac{1}{k + \text{rank}_r(d)}$$

Where $k = 60$ (constant), $R$ = set of ranking lists.

### Additional Features (Configured)

| Feature | What It Does | Config |
|---------|-------------|--------|
| **vectorWeight: 0.7** | Semantic search gets 70% weight | `query.hybrid.vectorWeight` |
| **textWeight: 0.3** | Keyword search gets 30% weight | `query.hybrid.textWeight` |
| **MMR (λ=0.7)** | Removes near-duplicate results | `query.hybrid.mmr` |
| **Temporal decay (30d)** | Recent notes rank higher (30-day half-life) | `query.hybrid.temporalDecay` |
| **Embedding cache (50K)** | Unchanged chunks aren't re-embedded | `query.cache` |
| **candidateMultiplier: 4** | Retrieves 4x candidates, then re-ranks | `query.hybrid.candidateMultiplier` |

---

## 6. What Gets Indexed

Everything under `extraPaths` + the agent's own memory:

| Source | Location | What |
|--------|----------|------|
| MEMORY.md | `~/.openclaw/workspace-legacy/MEMORY.md` | Curated facts about Zeus |
| Daily logs | `~/.openclaw/memory/YYYY-MM-DD.md` | Today's conversation context |
| **Codex vault** | `/home/onemind/codex/**/*.md` | **All 4 domains, every .md file** |

**File types indexed:** Only `.md` (Markdown) files.
**Not indexed:** Images, PDFs, binaries.
**Storage:** `~/.openclaw/agents/legacy/agent/memory.sqlite` (auto-created)

---

## 7. Setup Steps (One-Time)

### Step 1: Add deploy key to GitHub

Go to: `https://github.com/Zeus-Delacruz/onemindcodex-v0.4/settings/keys`

Add this key (already generated on VPS):
```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBilpHxjjyPbX3U6XrToqyj0HvrDoxzx6H34h83QEGE8 openclaw-vps-codex-readonly
```

### Step 2: Clone Codex on VPS

```bash
ssh root@134.122.8.135
sudo -u onemind git clone git@github-codex:Zeus-Delacruz/onemindcodex-v0.4.git /home/onemind/codex
```

### Step 3: Set up auto-pull cron

```bash
sudo -u onemind crontab -e
# Add:
*/5 * * * * cd /home/onemind/codex && git pull --ff-only >> /home/onemind/codex-sync.log 2>&1
```

### Step 4: Verify memorySearch config

The `openclaw.json` should already have the `memorySearch` block under `agents.defaults`. If not, add it per Section 2 above.

### Step 5: Restart OpenClaw

```bash
sudo systemctl restart openclaw
# Wait ~60 seconds for initial index build
```

### Step 6: Test

Ask Legacy: "Search my codex for morning protocol"
- Should return results from the HP domain
- Response time: <2 seconds (after initial index)

---

## 8. How Legacy Uses It (Examples)

```
Zeus: "What time does Ash get out of school?"
Legacy → memory_search "Ash school schedule"
→ LE/51 Family/ash-profile.md:34-41 (score: 0.93)
→ "Ash gets out at 3:15pm from Jefferson Elementary"

Zeus: "What are my macros today?"
Legacy → memory_search "nutrition macros targets"
→ HP/27 Body/nutrition.md (score: 0.91)
→ "2,800 kcal, 220g protein, 280g carbs, 80g fat"

Zeus: "What's the onboarding flow for Circle?"
Legacy → memory_search "Circle member onboarding"
→ GE/79 Community/onboarding-protocol.md (score: 0.88)

Zeus: "What's the n8n webhook URL?"
Legacy → memory_search "n8n todoist webhook URL"
→ UI/09 Automation/todoist-sync.md (score: 0.85)
```

---

## 9. Troubleshooting

| Problem | Cause | Fix |
|---------|-------|-----|
| No results from Codex | Codex not cloned | `ls /home/onemind/codex/` — should show 4 domain folders |
| Stale results | Cron not running | `crontab -l` as openclaw user — verify `*/5` line |
| Slow first search | Initial indexing | Wait 60-90s after restart. ~50K chunks take 30-60s. |
| "extraPaths" rejected | Wrong format | Must be `["/path"]` not `[{"path": "/path"}]` |
| Results missing recent notes | Sync lag | Check `tail /home/onemind/codex-sync.log` |

---

## 10. Architecture Diagram

```
┌──────────────┐         ┌───────────────┐        ┌──────────────────────┐
│  Mac/Obsidian │ ──git──→│    GitHub     │ ──git→ │  VPS: /home/openclaw │
│  (you edit)   │  push   │  (onemindco  │  pull  │  /codex/             │
│               │         │   dex-v0.4)  │  cron  │  ├── 00-24 SO/       │
└──────────────┘         └───────────────┘   5min  │  ├── 25-49 HP/       │
                                                    │  ├── 50-74 LE/       │
                                                    │  └── 75-99 GE/       │
                                                    └──────────┬───────────┘
                                                               │ extraPaths
                                                    ┌──────────▼───────────┐
                                                    │  OpenClaw SQLite     │
                                                    │  memory.sqlite       │
                                                    │  ├── BM25 FTS table  │
                                                    │  ├── vector embeddings│
                                                    │  └── embedding cache │
                                                    └──────────┬───────────┘
                                                               │ memory_search
                                                    ┌──────────▼───────────┐
                                                    │  Legacy agent        │
                                                    │  "Ash soccer at 4pm" │
                                                    └──────────────────────┘
```

---

*Last updated: 2026-03-01*
*Source: openclaw.json memorySearch config, CODEX-INTEGRATION.md, TOOLS.md*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
