# OneMind Memory Architecture
## Unified Memory System — Design Synthesis

> Combining patterns from: OpenClaw native memory, Letta Context Repositories, OpenViking Context Database, ClawVault structured memory, MemFS filesystem paradigm, and the Codex vault — into one coherent, AI-agent-agnostic memory architecture for OneMind OS.

---

## 1. The Problem Space

Every AI agent system re-invents memory with different tradeoffs:

| System | Strength | Weakness |
|--------|----------|----------|
| **OpenClaw native** | Markdown-as-truth, hybrid search (BM25+vector), temporal decay, session transcripts | Flat structure, no graph, no tiered loading |
| **Letta Context Repos** | Git-versioned, progressive disclosure, memory swarms, sleep-time compute | Tied to Letta runtime, git overhead for high-frequency writes |
| **OpenViking** | Filesystem paradigm (`viking://`), L0/L1/L2 tiered loading, directory recursive retrieval | Server-dependent, Python-only, heavy infrastructure |
| **ClawVault** | 8 primitives, knowledge graph, session resilience, Obsidian-native | CLI-heavy, no native tiered loading, no multi-agent coordination |
| **MemFS** | In-memory filesystem API, Node.js native `fs` interface, snapshot/restore | Pure utility — no search, no semantics, no persistence by default |
| **Codex (Obsidian)** | Human-readable, git-synced, rich linking, portable, owned | No native agent APIs, no vector search, no tiered loading |

**The insight**: No single system covers all needs. But they all converge on the same core primitives.

---

## 2. Five Convergent Patterns

Across all six systems, five patterns emerge repeatedly:

### Pattern 1: Filesystem as Memory Primitive
Every system that scales uses filesystem semantics — not database rows or opaque blobs.

- **OpenClaw**: `MEMORY.md` + `memory/YYYY-MM-DD.md` (Markdown files in workspace)
- **Letta**: Context Repositories = git-backed local filesystem, agents use `bash` and scripts
- **OpenViking**: `viking://` virtual filesystem protocol with `ls`, `find`, `tree`
- **ClawVault**: Structured vault directories (`decisions/`, `lessons/`, `people/`, `tasks/`)
- **MemFS**: Literal in-memory `fs` API implementation
- **Codex**: Obsidian vault = folders + Markdown files

**Conclusion**: The memory primitive is **files and directories**, not database records. Files are human-readable, git-friendly, tool-accessible, and agent-familiar.

### Pattern 2: Tiered Context Loading (Progressive Disclosure)
Loading everything into context wastes tokens and adds noise. Every mature system tiers its context.

- **OpenViking**: L0 (abstract ~100 tokens) → L1 (overview ~2k) → L2 (full content)
- **Letta**: `system/` directory always loaded; other files loaded by frontmatter description → full read
- **OpenClaw**: `MEMORY.md` always loaded; `memory/*.md` searched via `memory_search`; session transcripts optional
- **ClawVault**: `MEMORY.md` = boot context; vault = deep store; context profiles tune retrieval
- **Codex**: Folder READMEs → file frontmatter → full document

**Conclusion**: Memory needs at least **four tiers**: always-loaded pinned context (system/), navigable orientation (folder READMEs), discoverable summaries (frontmatter descriptions indexed), and deep content (full file on-demand).

### Pattern 3: Hybrid Search (Semantic + Keyword + Structure)
Pure vector search fails on exact tokens; pure keyword fails on paraphrases. Structure adds navigability.

- **OpenClaw**: BM25 + vector weighted merge, MMR diversity, temporal decay
- **OpenViking**: Vector retrieval → directory positioning → recursive drill-down
- **ClawVault**: BM25 + semantic embeddings + cross-encoder reranking + RRF + knowledge graph traversal
- **Letta**: Filesystem + frontmatter descriptions + agent-written scripts for custom queries
- **Codex**: Obsidian search (keyword) + dataview queries + wiki-link graph

**Conclusion**: Optimal retrieval combines **vector similarity + keyword matching + structural navigation** (directory hierarchy, graph links).

### Pattern 4: Session Resilience and Context Death Prevention
The #1 failure mode of agent memory: context death — losing everything when the window fills.

- **OpenClaw**: Auto-compaction with memory flush, `reserveTokensFloor`, pre-compaction dump to disk
- **Letta**: Git-backed so memory persists across sessions, sleep-time reflection
- **ClawVault**: `checkpoint`/`recover`, `wake`/`sleep` lifecycle, context death detection hooks
- **OpenViking**: Automatic session management, memory extraction at session end

**Conclusion**: Memory systems need **pre-compaction dump** (flush to persistent store before context death), **checkpoint/recover** semantics, and **session lifecycle hooks**.

### Pattern 5: Self-Improvement Loops (Reflection + Consolidation)
Memory isn't just recording — it's learning.

- **Letta**: Sleep-time compute — agents reflect during downtime, rewrite memory state
- **OpenViking**: Automatic session management extracts user preferences and agent experience
- **ClawVault**: `reflect` command, weekly promotion, archival cycle
- **OpenClaw**: Memory flush before compaction (basic), no native reflection loop

**Conclusion**: A complete memory system needs **reflection** (periodic review of recent experience), **consolidation** (promoting important information), **defragmentation** (reorganizing memory structure), and **initialization** (bootstrapping from existing data).

---

## 3. The Unified Architecture

### 3.1 Storage Layer: The Memory Filesystem (MemFS)

The foundation is a portable Markdown filesystem — the Codex is already this. The MemFS layer abstracts it so any agent runtime can mount it.

```
onemind-memory/                          ← "MemFS root" (= the codex + agent memory)
├── system/                              ← L0: Always loaded into context (boot)
│   ├── IDENTITY.md                      ← Who I am (from SOUL.md)
│   ├── USER.md                          ← Who Zeus is
│   ├── FOCUS.md                         ← Current priorities, active projects
│   └── DECISIONS.md                     ← Key standing decisions
├── memory/                              ← L1: Agent working memory
│   ├── daily/                           ← Append-only daily logs
│   │   ├── 2026-03-18.md
│   │   └── ...
│   ├── curated/                         ← Promoted durable memories
│   │   ├── MEMORY.md                    ← Executive summary (OpenClaw boot file)
│   │   ├── people.md                    ← Key people and relationships
│   │   ├── preferences.md               ← User preferences, habits
│   │   └── technical.md                 ← Technical decisions, stack choices
│   ├── sessions/                        ← Session transcripts (indexed)
│   │   └── YYYY-MM-DD-<session-id>.jsonl
│   ├── lessons/                         ← What went wrong/right (ClawVault pattern)
│   ├── handoffs/                        ← Cross-session continuity
│   └── checkpoints/                     ← Crash recovery state
├── knowledge/                           ← L2: Deep knowledge store (= Codex domains)
│   ├── 00-24-ui/                        ← The Brain (tech, infra, architecture)
│   ├── 25-49-hp/                        ← Self (health, skills, finance)
│   ├── 50-74-le/                        ← Legacy (family, home, estate)
│   └── 75-99-ge/                        ← Wealth (business, ventures)
├── agents/                              ← Per-agent memory partitions
│   ├── legacy/
│   │   ├── skills/                      ← Learned capabilities
│   │   ├── instructions/                ← Operating procedures
│   │   └── experience/                  ← Task execution history
│   ├── spartan/
│   ├── oracle/
│   └── .../
├── resources/                           ← Ingested external content
│   ├── repos/                           ← Code repositories (L0/L1/L2 processed)
│   ├── docs/                            ← External documentation
│   └── web/                             ← Web captures
├── .memindex/                           ← Index state (not human-edited)
│   ├── graph-index.json                 ← Knowledge graph from wiki-links
│   ├── embeddings.sqlite                ← Vector embeddings cache
│   ├── bm25.sqlite                      ← Full-text search index
│   ├── facts.json                       ← Extracted structured facts
│   └── entity-graph.json                ← Entity relationship graph
└── .memconfig/                          ← Memory system config
    ├── tiers.json                       ← L0/L1/L2 tier definitions
    ├── profiles.json                    ← Context profiles (planning, incident, etc.)
    └── policies.json                    ← Write policies, retention, decay
```

### 3.2 Tiered Loading Engine

Borrowing OpenViking's L0/L1/L2 model and applying it to the Codex filesystem:

```
┌─────────────────────────────────────────────────────────┐
│                   CONTEXT WINDOW                         │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │  L0 — ALWAYS LOADED (~500 tokens per file)       │   │
│  │  system/IDENTITY.md + system/USER.md +           │   │
│  │  system/FOCUS.md + system/DECISIONS.md            │   │
│  │  + memory/curated/MEMORY.md                       │   │
│  └──────────────────────────────────────────────────┘   │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │  L1 — NAVIGABLE (file tree + frontmatter index)  │   │
│  │  Agent sees directory structure + descriptions    │   │
│  │  Uses memory_search to find relevant L1 summaries │   │
│  └──────────────────────────────────────────────────┘   │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │  L2 — ON-DEMAND (full file reads)                │   │
│  │  Agent uses memory_get/read to load full content  │   │
│  │  Only when L1 indicates relevance                 │   │
│  └──────────────────────────────────────────────────┘   │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

**How each tier works:**

| Tier | What | When Loaded | Token Budget |
|------|------|-------------|--------------|
| System | Pinned context (`_codex/system/`) | Every prompt (pre-loaded) | ~500 tokens |
| L0 | Boot context, identity, current focus | Every prompt | ~2,000 tokens |
| L1 | File tree + frontmatter `ai_summary` descriptions | Every prompt (compact form) | ~1,500 tokens |
| L2 | Full file content | On-demand via search/read | Dynamic (fills remaining budget) |

Each file in the memory filesystem gets automatic L0/L1 generation:
- **L0** = first line of frontmatter `description:` field (~1 sentence)
- **L1** = full frontmatter + first section of file (or auto-generated summary)
- **L2** = the full file

### 3.3 Hybrid Retrieval Pipeline

Combining the best of OpenClaw, ClawVault, and OpenViking search:

```
Query
  │
  ├─→ BM25 Keyword Search (exact tokens, IDs, code symbols)
  │     └─→ top N × candidateMultiplier results
  │
  ├─→ Vector Semantic Search (paraphrases, concepts)
  │     └─→ top N × candidateMultiplier results
  │
  ├─→ Graph Traversal (wiki-links, entity relationships)
  │     └─→ connected nodes from relevant hits
  │
  └─→ Directory-Aware Positioning (OpenViking pattern)
        └─→ score boost for results in relevant directory subtree
  │
  ▼
Weighted Merge
  │  vectorWeight × vectorScore
  │+ textWeight   × textScore
  │+ graphWeight  × graphScore
  │+ dirWeight    × directoryRelevance
  │
  ▼
Post-Processing Pipeline
  │
  ├─→ Temporal Decay (halfLife: 30 days for daily notes)
  │     └─→ evergreen files exempt (MEMORY.md, curated/*)
  │
  ├─→ MMR Re-ranking (λ=0.7 diversity/relevance balance)
  │     └─→ remove near-duplicate snippets
  │
  └─→ Context Profile Filter (planning/incident/handoff/default)
        └─→ apply profile-specific weights and source preferences
  │
  ▼
Top-K Results → Agent Context
```

### 3.4 Session Lifecycle (Context Death Prevention)

Combining OpenClaw's memory flush with ClawVault's checkpoint/recover:

```
Session Start
  │
  ├─→ Load L0 (system/ + MEMORY.md)
  ├─→ Load L1 (file tree + descriptions index)
  ├─→ Check for recovery state (checkpoints/last-checkpoint.json)
  │     └─→ If exists: inject recovery context + alert
  └─→ Load recent daily memory (today + yesterday)
  │
  ▼
Running Session
  │
  ├─→ Periodic checkpoint (every N messages or on demand)
  │     └─→ Save: working-on, focus, open threads, tool state
  │
  ├─→ Approaching compaction threshold
  │     └─→ Memory Flush Turn (OpenClaw pattern)
  │           ├─→ Silent agentic turn: "store durable memories now"
  │           ├─→ Agent writes to memory/daily/YYYY-MM-DD.md
  │           └─→ Agent writes key decisions to memory/curated/
  │
  ├─→ Compaction triggers
  │     └─→ Context Engine compacts (summarize older history)
  │           └─→ If lossless-claw: DAG/vector compaction
  │           └─→ If legacy: built-in summarization
  │
  └─→ Session End
        ├─→ Final checkpoint (handoffs/latest.md)
        ├─→ Sleep-time reflection trigger (async)
        └─→ Summary written to memory/daily/YYYY-MM-DD.md
```

### 3.5 Reflection & Consolidation Loop (Sleep-Time Compute)

Borrowing from Letta's sleep-time compute and OpenViking's self-iteration:

```
┌──────────────────────────────────────────────────────────────┐
│                    REFLECTION CYCLE                            │
│                                                                │
│  ┌──────────┐     ┌──────────────┐     ┌───────────────┐     │
│  │ Session   │────▶│  Sleep-Time  │────▶│ Consolidation │     │
│  │ Transcripts│    │  Reflection  │     │ & Promotion   │     │
│  │ + Daily   │    │  (subagent)  │     │               │     │
│  │ Logs      │    └──────────────┘     └───────────────┘     │
│  └──────────┘           │                      │              │
│                         ▼                      ▼              │
│              ┌─────────────────┐    ┌──────────────────┐     │
│              │ Extract:        │    │ Promote to:      │     │
│              │ • Preferences   │    │ • curated/       │     │
│              │ • Decisions     │    │ • FOCUS.md       │     │
│              │ • Lessons       │    │ • DECISIONS.md   │     │
│              │ • Facts         │    │ • Codex domain   │     │
│              │ • Relationships │    │   folders         │     │
│              └─────────────────┘    └──────────────────┘     │
│                                                                │
│  Periodic Defragmentation:                                     │
│  ┌──────────────────────────────────────────────────────┐    │
│  │ • Merge duplicate memories                           │    │
│  │ • Split over-large files into focused documents      │    │
│  │ • Update graph index                                 │    │
│  │ • Archive stale daily logs (>90 days → compressed)   │    │
│  │ • Refresh L0/L1 summaries for changed files          │    │
│  └──────────────────────────────────────────────────────┘    │
│                                                                │
└──────────────────────────────────────────────────────────────┘
```

### 3.6 Multi-Agent Memory Architecture

For OneMind's 9-agent setup, memory must be both shared and partitioned:

```
┌─────────────────────────────────────────────────────┐
│              SHARED MEMORY (Codex)                    │
│  system/ + knowledge/ + resources/                    │
│  All 9 agents can READ via memory_search + read      │
│  Writes go through guardrails (domain folders only)  │
└──────────────────────┬──────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   ┌────▼────┐   ┌────▼────┐   ┌────▼────┐
   │ Legacy  │   │ Spartan │   │ Oracle  │  ... (per-agent)
   │ memory/ │   │ memory/ │   │ memory/ │
   │ ├ daily │   │ ├ daily │   │ ├ daily │
   │ ├ curat │   │ ├ curat │   │ ├ curat │
   │ └ sess  │   │ └ sess  │   │ └ sess  │
   └─────────┘   └─────────┘   └─────────┘

   Letta pattern: Each agent gets its own git worktree
   for concurrent reflection without blocking.

   Merge strategy: Agent writes → per-agent branch →
   auto-merge to main → conflicts resolved by Legacy
```

---

## 4. Implementation Mapping: What to Use Where

### 4.1 Agent-Agnostic Layer (Codex — owned by Zeus, wired to any runtime)

| Component | Implementation | Status |
|-----------|---------------|--------|
| **Filesystem** | Obsidian vault in monorepo (`codex/`) | ✅ Exists |
| **L0/L1 summaries** | Frontmatter `description:` + folder READMEs | 🔧 Enhance |
| **Knowledge graph** | Obsidian wiki-links + ClawVault `graph-index.json` | 🆕 Add |
| **Git versioning** | Already git-tracked in monorepo | ✅ Exists |
| **Sync** | `onemind-pull.timer` (VPS) + Obsidian Git (local) | ✅ Exists |

### 4.2 OpenClaw Layer (Legacy's runtime — wired to OpenClaw)

| Component | Implementation | Status |
|-----------|---------------|--------|
| **Memory files** | `MEMORY.md` + `memory/YYYY-MM-DD.md` | ✅ Exists |
| **Vector search** | `memorySearch` with OpenAI embeddings, hybrid BM25+vector | ✅ Exists |
| **Extra paths** | `memorySearch.extraPaths → codex/` | ✅ Configured |
| **Session transcripts** | `memorySearch.experimental.sessionMemory` | 🔧 Enable |
| **Context engine** | `plugins.slots.contextEngine` | 🆕 Evaluate lossless-claw |
| **Memory plugin** | `plugins.slots.memory = "memory-core"` | ✅ Active |
| **Memory flush** | `compaction.memoryFlush.enabled` | ✅ Active |
| **QMD backend** | `memory.backend = "qmd"` | 🔧 Evaluate |
| **Multimodal** | Gemini embedding 2 for images/audio | 🆕 Future |
| **ClawVault hooks** | `openclaw hooks install clawvault` | 🆕 Install (native plugin since v3.3) |
| **OpenViking plugin** | `openclaw-memory-plugin` (52% accuracy, 91% token savings) | 🆕 Evaluate |

### 4.3 Enhancement Priorities

| Priority | Enhancement | Source Pattern | Impact |
|----------|------------|---------------|--------|
| **P0** | Enable session memory search | OpenClaw native | Recall recent conversations |
| **P0** | Add L0 frontmatter to all codex files | OpenViking L0/L1/L2 | Better navigability |
| **P1** | Install ClawVault for structured memory | ClawVault | Session resilience, knowledge graph |
| **P1** | Evaluate lossless-claw context engine | OpenClaw plugin | Better compaction than legacy |
| **P1** | Implement reflection cycle | Letta sleep-time | Memory self-improvement |
| **P2** | Evaluate OpenViking as memory plugin | OpenViking | 52% accuracy improvement |
| **P2** | Add temporal decay config | OpenClaw native | Better recency handling |
| **P2** | Enable MMR re-ranking | OpenClaw native | Reduce duplicate results |
| **P3** | Build multi-agent memory partitioning | Letta worktrees | Concurrent agent memory |
| **P3** | Add multimodal memory (images/audio) | OpenClaw Gemini | Index visual/audio notes |

---

## 5. Future Memory Structures — Emerging Patterns

### 5.1 Memory as a Compiled Artifact
Current: memory is raw files. Future: memory is a **compiled knowledge base** with pre-computed summaries, pre-built indexes, and pre-resolved relationships — like a compiled binary vs source code.

### 5.2 Memory Swarms (Concurrent Reflection)
Letta's pattern of multiple subagents working in git worktrees to process memory concurrently. Each agent reflects on a slice, changes merge back. This is the path to scaling reflection.

### 5.3 Skill Memory (Learning from Doing)
Not just episodic memory (what happened) but procedural memory (how to do things). ClawVault's 8 primitives + Letta's skill learning. Agents remember not just facts but **capabilities and strategies that worked**.

### 5.4 Predictive Loading (Anticipatory Context)
Instead of loading memory on-demand when the agent asks, **predict what memory will be needed** based on the conversation trajectory and pre-load it. Like CPU cache prefetching for context windows.

### 5.5 Federated Memory (Cross-Agent Knowledge)
Agents don't just have private memory — they contribute to a **shared knowledge graph** where insights from Spartan about workout routines are available to Haven for scheduling, and Oracle's financial analysis feeds Empire's business decisions. The Codex already serves this role; the enhancement is making it agent-writable with conflict resolution.

### 5.6 Memory Aging and Forgetting
Not everything should be remembered forever. Temporal decay is the start, but true memory aging includes:
- **Active forgetting**: removing outdated information
- **Compression**: reducing verbose memories to their essence
- **Archival**: moving cold memories to cheaper/deeper storage
- **Contradiction resolution**: when new memory conflicts with old

### 5.7 Memory Provenance and Trust
As agents write more memory, provenance becomes critical:
- Who wrote this memory?
- When was it written?
- What was the confidence level?
- Has it been verified by a human?
- Has it been contradicted by newer information?

---

## 6. Configuration Blueprint

### Proposed `openclaw.json` memory config:

```json
{
  "agents": {
    "defaults": {
      "memorySearch": {
        "provider": "openai",
        "model": "text-embedding-3-small",
        "fallback": "local",
        "extraPaths": ["/home/onemind/onemind-os/codex"],
        "query": {
          "hybrid": {
            "enabled": true,
            "vectorWeight": 0.7,
            "textWeight": 0.3,
            "candidateMultiplier": 4,
            "mmr": { "enabled": true, "lambda": 0.7 },
            "temporalDecay": { "enabled": true, "halfLifeDays": 30 }
          }
        },
        "experimental": { "sessionMemory": true },
        "sources": ["memory", "sessions"],
        "cache": { "enabled": true, "maxEntries": 50000 },
        "sync": { "watch": true }
      },
      "compaction": {
        "memoryFlush": {
          "enabled": true,
          "softThresholdTokens": 4000
        }
      }
    }
  },
  "plugins": {
    "slots": {
      "memory": "memory-core",
      "contextEngine": "lossless-claw"
    },
    "entries": {
      "lossless-claw": { "enabled": true }
    }
  }
}
```

---

## 7. The Vision: OneMind Memory Stack

```
┌──────────────────────────────────────────────────────────────┐
│                     AGENT RUNTIMES                            │
│   OpenClaw (Legacy)  ·  Letta  ·  Agno  ·  Custom           │
└──────────────┬───────────────────────────┬───────────────────┘
               │                           │
┌──────────────▼───────────────────────────▼───────────────────┐
│                    MEMORY API LAYER                            │
│  search(query) · read(path) · write(path) · reflect()        │
│  checkpoint() · recover() · consolidate()                     │
│  Implemented as: OpenClaw tools, MCP server, REST API,        │
│  or direct filesystem access                                  │
└──────────────┬───────────────────────────┬───────────────────┘
               │                           │
┌──────────────▼──────────┐  ┌─────────────▼───────────────────┐
│   RETRIEVAL ENGINE      │  │   CONTEXT ENGINE                 │
│   Hybrid Search         │  │   (lossless-claw / legacy)       │
│   BM25 + Vector + Graph │  │   Ingest → Assemble → Compact   │
│   + Dir-Aware + Decay   │  │   + Session lifecycle hooks      │
│   + MMR + Profiles      │  │   + Memory flush                 │
└──────────────┬──────────┘  └─────────────┬───────────────────┘
               │                           │
┌──────────────▼───────────────────────────▼───────────────────┐
│                    MEMORY FILESYSTEM (MemFS)                   │
│  system/  ·  memory/  ·  knowledge/  ·  agents/  ·  resources/│
│  Markdown-native · Git-versioned · Obsidian-compatible         │
│  L0/L1/L2 tiered · Wiki-linked · Frontmatter-described        │
└──────────────┬───────────────────────────────────────────────┘
               │
┌──────────────▼───────────────────────────────────────────────┐
│                    INDEX LAYER                                 │
│  SQLite vectors · BM25 full-text · Knowledge graph             │
│  Entity graph · Fact store · Embedding cache                   │
└──────────────┬───────────────────────────────────────────────┘
               │
┌──────────────▼───────────────────────────────────────────────┐
│                    PERSISTENCE                                │
│  Git repo (codex/) · VPS filesystem · NATS for cross-runtime  │
│  Sync: onemind-pull.timer + onemind-watch.service                  │
└──────────────────────────────────────────────────────────────┘
```

**The key principle**: The Codex (Markdown vault) is the **owned, portable, agent-agnostic** memory substrate. Everything above it — search engines, context engines, plugins — is **interchangeable wiring** that can be swapped without losing the memory itself.

You can change from OpenClaw to Letta to a future system. The Codex remains. The files remain. The knowledge remains. That's the memory system.

---

## 8. What Each System Offers — Side-by-Side

| Capability | OpenClaw Native | Letta | OpenViking | ClawVault | MemFS | Codex (Obsidian) |
| --- | --- | --- | --- | --- | --- | --- |
| **Storage format** | Markdown files in workspace | Git-backed filesystem | `viking://` virtual FS (server) | Markdown vault + `.clawvault/` | In-memory Node.js `fs` API | Obsidian Markdown vault |
| **Search** | Hybrid BM25+vector, MMR, temporal decay | Agent-written scripts + frontmatter | Vector + directory recursive retrieval | BM25 + semantic + RRF + graph | None | Obsidian search + Dataview |
| **Tiered loading** | MEMORY.md always → search for rest | `system/` always → frontmatter → full read | L0 abstract → L1 overview → L2 full | MEMORY.md boot → vault deep | N/A | README per folder → full doc |
| **Session resilience** | Memory flush before compaction | Git versioning, persistent across sessions | Auto session management | checkpoint/recover/wake/sleep | Snapshot/restore | N/A (no sessions) |
| **Reflection/learning** | Basic (memory flush only) | Sleep-time compute, memory swarms | Auto-extract preferences + experience | reflect command, weekly promote | N/A | Manual review |
| **Graph/linking** | None | None | Directory hierarchy traversal | Wiki-link knowledge graph + entity graph | None | Wiki-links, tags, Dataview |
| **Multi-agent** | Per-agent workspace + shared extraPaths | Git worktrees for concurrent reflection | Per-agent context partitions | Single agent focus | N/A | Shared vault (all agents read) |
| **Context engine** | Plugin slot (legacy/lossless-claw) | Built-in MemFS | Built-in L0/L1/L2 engine | None (relies on host) | N/A | N/A |
| **Compaction** | Summarize older history, DAG (lossless-claw) | Agent manages own context repo | Auto-compress, extract long-term memory | None built-in | N/A | N/A |
| **Obsidian compat** | No (separate workspace dir) | No | No | Yes (designed for it) | No | **IS** Obsidian |
| **Human-editable** | Yes (Markdown) | Yes (files) | No (server-processed) | Yes (Markdown) | No (in-memory) | Yes (Markdown) |
| **Multimodal** | Gemini embedding 2 (images/audio) | No | VLM for image understanding | No | No | Images as assets |
| **Portability** | Tied to OpenClaw workspace format | Tied to Letta Code runtime | Tied to OpenViking server | Portable (npm CLI) | Portable (npm lib) | **Fully portable** |
| **License** | Proprietary (OpenClaw) | Open source (Apache 2.0) | Open source (Apache 2.0) | MIT | Apache 2.0 | Your files |

---

## 9. What We Are Doing With Each

### What We USE Today

| System | Role in OneMind | How |
| --- | --- | --- |
| **OpenClaw** | Legacy's runtime — all 9 agents run here | `openclaw.json` config, `workspaces/legacy/`, systemd on VPS |
| **Codex** | Source of truth for all knowledge | `codex/` in monorepo, synced to VPS via `onemind-pull.timer` |
| **OpenClaw memory-core** | Vector search over workspaces + codex | `memorySearch.extraPaths → /home/onemind/onemind-os/codex` |

### What We EVALUATE Next

| System | Why | What It Would Replace/Add |
| --- | --- | --- |
| **ClawVault** | Session resilience, knowledge graph, Obsidian-native | Adds structured memory on top of OpenClaw, installs as hook pack |
| **lossless-claw** | Better compaction than legacy summarization | Replaces `plugins.slots.contextEngine: "legacy"` |
| **OpenViking** | 52% accuracy improvement, 91% token savings | Replaces `plugins.slots.memory: "memory-core"` as memory plugin |
| **Letta patterns** | Sleep-time reflection, memory defrag | Implement as lobster routines or heartbeat tasks, not Letta runtime |

### What We DON'T Use (and Why)

| System | Why Not |
| --- | --- |
| **Letta runtime** | Legacy stays in OpenClaw. We borrow Letta's patterns (git-versioned memory, sleep-time reflection) but implement them in OpenClaw |
| **MemFS (raw)** | Pure utility library. The pattern (filesystem = memory) is what matters, not the npm package itself |
| **OpenViking server** | Too heavy — runs a separate Python server + Go AGFS. We'd use just the OpenClaw memory plugin, not the full server |

---

## 10. Codex Development Strategy

### 10.1 The Repo Question: Separate or Together?

**Recommendation: Keep it in the monorepo.** Here's the analysis:

#### Arguments for separate repo

- Codex syncs independently (agents write to it, humans edit in Obsidian)
- Lighter clone for devices that only need the vault
- Cleaner git history (codex commits don't mix with code commits)

#### Arguments for keeping together (STRONGER)

- **Infra is tightly coupled** — `docker-compose.yml`, NATS configs, deploy scripts all reference codex paths
- **Agent workspaces reference codex** — TOOLS.md, AGENTS.md point at codex paths
- **One sync mechanism** — `onemind-pull.timer` already pulls the whole monorepo; splitting adds a second sync
- **Scripts need both** — shell scripts, Python tools, deploy configs all live alongside the codex
- **0-99 covers everything** — the codex IS the documentation layer for the entire system. Splitting it means every doc has two possible homes
- **Atomic commits** — when you change infra + update the codex doc about it, that's one commit

#### The real problem you're identifying

The monorepo has **loose files at root** that should be IN the codex:

```
CURRENT ROOT (messy)                    → WHERE THEY BELONG
─────────────────────                   ─────────────────────
1PASSWORD-CLEANUP.md                    → codex/08 Infrastructure/
ACCESS.md                              → codex/11 Security/
BLUEPRINT.md                           → codex/00 Framework/
CHANGELOG.md                           → KEEP at root (standard)
CODE.md                                → codex/10 Development/
DEVELOPMENT.md                         → codex/10 Development/
FEATURES.md                            → codex/10 Development/
GITHUB-RESEARCH.md                     → codex/12 Research/
GPU-ANALYSIS.md                        → codex/12 Research/
open-source-embedding-models-2026.md   → codex/12 Research/
PRODUCTION-PLAN.md                     → codex/10 Development/
SWARMCLAW.md                           → codex/01 Command/
STATUS.md                              → codex/08 Infrastructure/
TYPESCRIPT-CROSS-PLATFORM.md           → codex/10 Development/
roadmap/**                             → codex/12 Research/ (or codex domain folders)
infra/*.md                             → codex/08 Infrastructure/
scripts/                               → KEEP (runtime, not docs)
```

#### The clean split

```
onemindos-v0.5/
├── codex/                    ← ALL knowledge, docs, plans, research (0-99)
│   ├── 00-24 SO/             ← includes infra docs, agent docs, scripts docs
│   ├── 25-49 HP/
│   ├── 50-74 LE/
│   ├── 75-99 GE/
│   └── ...
├── openclaw/                 ← RUNTIME: Legacy's OpenClaw workspace + skills + deploy
│   ├── workspaces/           ← Agent workspace configs (SOUL.md, TOOLS.md, skills/)
│   └── deploy/               ← Systemd units, env templates, deploy scripts
├── infra/                    ← RUNTIME: docker-compose, NATS configs, VPS scripts
│   ├── docker-compose.yml    ← Runtime config (not docs)
│   ├── nats/                 ← Runtime configs
│   └── scripts/              ← Shell scripts that run things
├── scripts/                  ← RUNTIME: utility scripts
├── AGENTS.md                 ← Keep at root (agent instructions for this repo)
├── CLAUDE.md                 ← Keep at root (Claude instructions for this repo)
├── README.md                 ← Keep at root (standard)
└── CHANGELOG.md              ← Keep at root (standard)
```

**Rule**: If it's a **document/plan/reference** → goes in codex. If it **runs or configures something** → stays in repo root.

### 10.2 How the Codex Attaches to Any Agent Runtime

The codex is agent-agnostic by design. Here's how each runtime connects:

```
┌─────────────────────────────────────────────────────────────────┐
│                    THE CODEX (Markdown vault)                    │
│                    /home/onemind/onemind-os/codex/              │
│                                                                  │
│  Contains: knowledge, docs, plans, research, preferences,        │
│  decisions, people, SOPs, reference material (00-99)             │
│  Format: Markdown + frontmatter + wiki-links + tags              │
│  Sync: Git (source of truth) → VPS (onemind-pull every 1 min)     │
└──────────┬──────────────┬───────────────┬───────────────────────┘
           │              │               │
    ┌──────▼──────┐ ┌────▼──────┐ ┌──────▼──────┐
    │  OPENCLAW   │ │  LETTA    │ │ AGNO/OTHER  │
    │  (current)  │ │  (future) │ │  (future)   │
    │             │ │           │ │             │
    │ Attach via: │ │ Attach:   │ │ Attach:     │
    │ extraPaths  │ │ /memfs    │ │ MCP server  │
    │ in memory   │ │ clone →   │ │ REST API    │
    │ Search cfg  │ │ local fs  │ │ direct fs   │
    └─────────────┘ └───────────┘ └─────────────┘
```

**OpenClaw (current — all 9 agents use this today)**:

```json
{
  "agents": {
    "defaults": {
      "memorySearch": {
        "extraPaths": ["/home/onemind/onemind-os/codex"]
      }
    }
  }
}
```

This single config line gives every OpenClaw agent read access to the entire codex via `memory_search`. They can also `read` files directly. Write access is controlled by the workspace `AGENTS.md` guardrails (domain folders only, new notes to inbox).

**Letta Code (if we ever add it)**:

```bash
# Symlink codex into Letta's context repo
ln -s /home/onemind/onemind-os/codex ~/.letta/agents/legacy/memfs/knowledge
```

Letta agents would see `knowledge/00-24-ui/...` in their filesystem.

**Any future runtime (MCP server pattern)**:

Build a thin MCP server that exposes:
- `codex_search(query)` → hybrid search over codex files
- `codex_read(path)` → read a codex file
- `codex_write(path, content)` → write to codex (with guardrails)
- `codex_tree(path, depth)` → list directory structure

Any agent that speaks MCP can use the codex. This is the universal adapter.

### 10.3 What Lives in OpenClaw vs What's Universal

| Layer | What | OpenClaw-Specific? | Universal (Codex)? |
| --- | --- | --- | --- |
| **Agent identity** | SOUL.md, IDENTITY.md | ✅ OpenClaw workspace format | ❌ |
| **Agent instructions** | AGENTS.md, TOOLS.md, BOOT.md | ✅ OpenClaw workspace format | ❌ |
| **Agent skills** | skills/ (Node.js packages) | ✅ OpenClaw skill format | ❌ |
| **Agent memory (working)** | memory/YYYY-MM-DD.md, MEMORY.md | ✅ OpenClaw memory format | ❌ |
| **Session transcripts** | sessions/*.jsonl | ✅ OpenClaw session format | ❌ |
| **Memory search config** | memorySearch in openclaw.json | ✅ OpenClaw config | ❌ |
| **Context engine** | lossless-claw plugin | ✅ OpenClaw plugin | ❌ |
| **Knowledge base** | 0-99 domain folders | ❌ | ✅ **Codex** |
| **Decisions & preferences** | What Zeus prefers, standing rules | ❌ | ✅ **Codex** |
| **People & relationships** | Contact info, relationship notes | ❌ | ✅ **Codex** |
| **SOPs & protocols** | How to do things | ❌ | ✅ **Codex** |
| **Research & plans** | Analysis, architecture docs | ❌ | ✅ **Codex** |
| **Infra documentation** | VPS setup, service docs | ❌ | ✅ **Codex** |
| **Life management** | Health, finance, family, home | ❌ | ✅ **Codex** |
| **Business** | Ventures, entrepreneurship | ❌ | ✅ **Codex** |
| **Scripts & tools** | Shell, Python utilities | ✅ Runtime (stays in repo) | ❌ |
| **Deploy configs** | systemd, docker-compose, env | ✅ Runtime (stays in repo) | ❌ |

**The principle**: OpenClaw owns the **agent runtime** (how agents think, remember, and act). The Codex owns the **knowledge** (what they know about). The codex is the "brain's library" — the runtime is the "brain's wiring."

### 10.4 Getting the Full Codex to All 9 Agents

**Current state**: Already working! The config already does this:

```json
{
  "agents": {
    "defaults": {
      "memorySearch": {
        "extraPaths": ["/home/onemind/onemind-os/codex"]
      }
    }
  }
}
```

Because this is under `agents.defaults`, ALL 9 agents (Legacy, Spartan, Oracle, Eden, Empire, Haven, Guardian, Grid, Trinity) get the full codex indexed. They all share the same vector index of the entire codex.

**What to enhance**:

1. **Per-agent domain weighting** — Spartan should prioritize `25-49 HP/27 Body/` results, Empire should prioritize `75-99 GE/`. This can be done with:
   - Per-agent `extraPaths` that list their domain folders first
   - Or context profiles that boost domain-relevant results

2. **Write guardrails per agent** — Already defined in each agent's workspace AGENTS.md:
   - Spartan writes to `27 Body/` and `28 Mastery/`
   - Oracle writes to `25 Identity/` and `30 Finance/`
   - Eden writes to `50-74 LE/` farming subfolders
   - etc.

3. **Shared curated memory** — Top-level `MEMORY.md` equivalent for cross-agent knowledge:
   - Each agent writes their domain insights to codex
   - All agents can read all domains
   - Legacy coordinates and resolves conflicts

### 10.5 What to Collapse Into the Codex

These files/folders currently at repo root should move into the codex:

| Current Location | Move To | Reason |
| --- | --- | --- |
| `roadmap/**` | `codex/12 Research/roadmap/` | Planning docs, not runtime |
| `BLUEPRINT.md` | `codex/00 Framework/` | Architecture doc |
| `FEATURES.md` | `codex/10 Development/` | Feature planning |
| `PRODUCTION-PLAN.md` | `codex/10 Development/` | Deployment planning |
| `CODE.md` | `codex/10 Development/` | Code conventions |
| `DEVELOPMENT.md` | `codex/10 Development/` | Dev setup guide |
| `STATUS.md` | `codex/08 Infrastructure/` | System status |
| `GPU-ANALYSIS.md` | `codex/12 Research/` | Research doc |
| `GITHUB-RESEARCH.md` | `codex/12 Research/` | Research doc |
| `open-source-embedding-models-2026.md` | `codex/12 Research/` | Research doc |
| `TYPESCRIPT-CROSS-PLATFORM.md` | `codex/10 Development/` | Architecture doc |
| `1PASSWORD-CLEANUP.md` | `codex/11 Security/` | Security doc |
| `ACCESS.md` | `codex/11 Security/` | Access doc |
| `SWARMCLAW.md` | `codex/01 Command/` | Dashboard doc |
| `infra/README.md` | `codex/08 Infrastructure/` | Infra doc |
| `infra/STATUS-REPORT.md` | `codex/08 Infrastructure/` | Infra doc |
| `infra/VPS-MAIN.md` | `codex/08 Infrastructure/` | VPS doc |
| `infra/VPS-SENTINEL.md` | `codex/08 Infrastructure/` | VPS doc |

**Keep at repo root** (these run things, not just document things):

| Keep | Reason |
| --- | --- |
| `AGENTS.md` | AI agent instructions for this repo |
| `CLAUDE.md` | Claude-specific instructions for this repo |
| `README.md` | Standard repo readme |
| `CHANGELOG.md` | Standard changelog |
| `openclaw/` | Runtime (workspaces, skills, deploy configs) |
| `infra/docker-compose.yml` | Runtime config |
| `infra/nats/` | Runtime configs |
| `infra/scripts/` | Executable scripts |
| `scripts/` | Executable scripts |
| `.claude/skills/` | Claude skill definitions (runtime) |

### 10.6 Codex Development Roadmap

```
Phase 1: CONSOLIDATE (now)
├── Move scattered docs into codex domain folders
├── Add frontmatter (description, tags) to every codex file
├── Ensure every folder has a README.md (L1 navigability)
└── Clean up root — only runtime + standard files remain

Phase 2: ENHANCE SEARCH (next)
├── Enable OpenClaw session memory search
├── Enable temporal decay + MMR re-ranking
├── Evaluate ClawVault for knowledge graph overlay
└── Add context profiles (planning, incident, etc.)

Phase 3: SMART LOADING (later)
├── Generate L0/L1 summaries for all codex files
├── Build directory-aware retrieval (OpenViking pattern)
├── Evaluate lossless-claw or OpenViking plugin
└── Per-agent domain weighting

Phase 4: SELF-IMPROVING (future)
├── Implement sleep-time reflection (lobster routine)
├── Auto-consolidation of daily logs → curated memory
├── Memory defragmentation cycle
└── Cross-agent knowledge federation
```

---

## 11. Summary Decision Matrix

| Question | Answer |
| --- | --- |
| **Separate codex repo?** | No. Keep in monorepo. Infra + agent configs + codex are tightly coupled. |
| **What moves into codex?** | All docs, plans, research, reference material. Only runtime configs and executable scripts stay outside. |
| **How do all agents access codex?** | Already done — `memorySearch.extraPaths` under `agents.defaults` gives all 9 agents the full codex. |
| **What's OpenClaw-specific?** | Agent identity (SOUL.md), skills (Node.js), session format, memory format, context engine plugins. |
| **What's universal (codex)?** | All knowledge (0-99), decisions, preferences, people, SOPs, research, plans. Portable Markdown. |
| **How does codex attach to other runtimes?** | Symlink (Letta), MCP server (any future runtime), or direct filesystem mount. The Markdown files are the interface. |
| **Best memory plugin?** | Start with OpenClaw native (memory-core + hybrid search). Evaluate OpenViking plugin for accuracy gains. Add ClawVault for graph + session resilience. |
| **Best context engine?** | Evaluate lossless-claw to replace legacy summarization compaction. |

---

## 12. Making the Codex Native — The Two-Layer Architecture

The confusion about "plugin vs skill vs codex folder" dissolves with one insight: **two layers, each with a clear job.**

### The Two Layers

```
┌─────────────────────────────────────────────────────────────────────┐
│                        LAYER 1: THE CODEX                            │
│               (universal, agent-agnostic, the PRODUCT)               │
│                                                                      │
│   codex/_codex/                                                      │
│   ├── INTERFACE.md      ← "Here's how to use this vault"            │
│   ├── CONVENTIONS.md    ← Naming, tags, frontmatter rules           │
│   ├── system/           ← PINNED CONTEXT (always loaded)            │
│   ├── skills/           ← INSTRUCTIONS (markdown, no code)          │
│   │   ├── capture.md    ← How to capture new information            │
│   │   ├── organize.md   ← How to file into 00-99 domains           │
│   │   ├── search.md     ← How to search (4-tier loading)            │
│   │   ├── reflect.md    ← How to do reflection cycles               │
│   │   ├── consolidate.md← How to merge daily notes → knowledge      │
│   │   ├── maintain.md   ← How to keep the vault healthy             │
│   │   ├── init.md       ← How to bootstrap a vault / onboard agent  │
│   │   └── defrag.md     ← How to restructure cluttered domains      │
│   ├── profiles/         ← Context modes                             │
│   │   ├── planning.md   ← Pre-load planning context                 │
│   │   ├── incident.md   ← Pre-load incident response context        │
│   │   └── review.md     ← Pre-load weekly review context            │
│   └── templates/        ← Note templates for consistency            │
│       ├── daily.md      ├── project.md   ├── decision.md            │
│       └── person.md     └── review.md    └── sop.md                 │
│                                                                      │
│   Lives IN the vault. Travels WITH the vault. Readable by humans     │
│   AND any AI. This IS the product.                                   │
└──────────────────────────────┬──────────────────────────────────────┘
                               │
                               │ reads instructions from
                               │
┌──────────────────────────────▼──────────────────────────────────────┐
│                     LAYER 2: RUNTIME ADAPTERS                        │
│            (per-runtime, implements the instructions as tools)        │
│                                                                      │
│   OpenClaw:  skills/codex-1.0.0/     ← Node.js skill with tools     │
│   Claude:    .claude/skills/codex/   ← SKILL.md for Claude Code     │
│   Cursor:    .cursor/rules/codex.mdc ← Rules file for Cursor        │
│   Universal: MCP server              ← codex-mcp for any MCP client │
│                                                                      │
│   Each adapter reads _codex/ to know the rules, then executes.       │
│   Switch runtimes? The instructions stay. Only the adapter changes.  │
└─────────────────────────────────────────────────────────────────────┘
```

**Why this is better than "just a skill" or "just a plugin":**

| Approach | Problem |
| --- | --- |
| OpenClaw skill only | Tied to OpenClaw. Switch runtimes, lose everything. |
| Plugin only | Runtime-specific, doesn't travel with the vault. |
| Codex folder only (instructions, no code) | Any agent can read it, but no executable tools. |
| **Both layers** | Instructions travel with vault. Tools execute in any runtime. |

### Layer 1 Deep Dive: `_codex/` — The Vault's Self-Description

The `_codex/` folder is the codex's own operating manual. It lives inside the vault so it's always present — whether you open it in Obsidian, clone it on a new machine, or point an AI agent at it.

**Why `_codex/` not `.codex/`?** Obsidian hides dotfiles. The underscore prefix keeps it visible in Obsidian but signals "system" to users and agents.

#### `_codex/INTERFACE.md` — The Universal Entry Point

This is the ONE file any AI agent reads first. It tells them:
- What this vault is (a CODEX-structured knowledge base)
- How it's organized (00-99 domains across 3 domains + 1 operating layer)
- Where to find instructions (`_codex/skills/`)
- What conventions to follow (`_codex/CONVENTIONS.md`)
- How to search effectively (`_codex/skills/search.md`)
- What guardrails exist (don't create top-level folders, inbox for new items)

Any runtime adapter's first action is: `read("codex/_codex/INTERFACE.md")`.

#### `_codex/skills/` — Agent-Agnostic Operations (The Key Innovation)

These are **pure markdown instructions** — no code, no runtime dependencies. They describe operations any intelligent agent (human or AI) can perform on the vault:

| Skill File | What It Teaches | Best Feature It Borrows |
| --- | --- | --- |
| `capture.md` | How to ingest new info into the vault | Letta's universal capture |
| `organize.md` | How to file items into 00-99 domains | CODEX method step O |
| `search.md` | How to find things (4-tier: system → L0 → L1 → L2) | OpenViking L0/L1/L2 + Letta system/ |
| `reflect.md` | How to do reflection cycles (daily, weekly, monthly) | Letta's sleep-time compute |
| `consolidate.md` | How to merge scattered notes into curated knowledge | ClawVault's fact extraction |
| `maintain.md` | How to keep the vault healthy (orphan detection, tag cleanup, dedup) | ClawVault's session resilience |
| `init.md` | How to bootstrap a new vault or onboard an agent | Letta's memory initialization |
| `defrag.md` | How to restructure cluttered domains (split, merge, reorganize) | Letta's memory defragmentation |

**Example — `_codex/skills/search.md`** (borrowing OpenViking's tiered loading):

```markdown
# Codex Search — Tiered Loading Protocol

When searching the codex, use progressive disclosure to save context:

## Level 0: Orientation (cheapest)
Read the domain README first.
- `codex/00-24 SO/README.md` → overview of all tech/infra topics
- Scan folder names to narrow down

## Level 1: Discovery
Read the target folder's README or scan file names.
- `codex/00-24 SO/08 Infrastructure/README.md`
- List files in the folder

## Level 2: Deep Read
Only read full files when you've confirmed they're relevant.
- Read the specific file
- If it has frontmatter `related:` links, follow those too

## Search Commands (for runtimes that support search)
- Hybrid search: combine keyword + semantic for best results
- Always filter by domain folder first when the domain is known
- Check `tags:` in frontmatter for cross-domain discovery
```

**The power**: These instructions work whether read by a human, an OpenClaw agent, a Claude Code session, a GPT agent, or any future AI. The vault teaches whoever opens it how to use it.

#### `_codex/profiles/` — Context Modes

Pre-built "loading plans" for common situations. An agent reads a profile to know what to pre-load:

```markdown
# Profile: Weekly Review
# Load this context when doing a weekly review

## Pre-load these files:
- codex/25-49 HP/29 Life Systems/weekly-review-template.md
- codex/00-24 SO/00 Framework/active-goals.md
- Last 7 daily notes from the agent's memory/

## Focus domains:
- 29 Life Systems (daily ops review)
- 30 Finance (spending review)
- 28 Mastery (skill progress)

## Actions:
1. Follow _codex/skills/reflect.md
2. Generate weekly summary
3. Update active goals
```

#### `_codex/templates/` — Consistent Note Creation

Markdown templates with frontmatter. When creating new notes, agents use these instead of inventing formats:

```markdown
# Template: Decision Record
---
id: {domain-code}-{YYYYMMDD}-{slug}
type: decision
status: active
created: {date}
decided_by: zeus
tags: [type/decision]
---
# {Title}

## Context
{What situation prompted this decision?}

## Decision
{What was decided?}

## Rationale
{Why this option over alternatives?}

## Consequences
{What changes as a result?}
```

### Layer 2 Deep Dive: The OpenClaw `codex` Skill

This is the runtime adapter for OpenClaw. It reads `_codex/` for its instructions and provides executable tools:

```
openclaw/workspaces/legacy/skills/codex-1.0.0/
├── SKILL.md          ← Discovery: when to activate
├── index.js          ← Tools: codex_capture, codex_search, codex_organize, etc.
└── package.json
```

**What the skill provides (tools):**

| Tool | What It Does | Reads From |
| --- | --- | --- |
| `codex_capture` | Create a new note in the inbox with proper frontmatter | `_codex/skills/capture.md` + `_codex/templates/` |
| `codex_search` | Tiered search: system → L0 → L1 → L2 | `_codex/skills/search.md` |
| `codex_organize` | Move a note from inbox to the correct 00-99 domain | `_codex/skills/organize.md` |
| `codex_reflect` | Run a reflection cycle (daily, weekly, monthly) | `_codex/skills/reflect.md` + `_codex/profiles/` |
| `codex_consolidate` | Merge scattered daily notes into curated knowledge | `_codex/skills/consolidate.md` |
| `codex_health` | Audit vault: orphans, missing frontmatter, broken links | `_codex/skills/maintain.md` |
| `codex_tree` | Show folder structure at specified depth | Built-in |

**Why a skill and not a plugin?** In OpenClaw:
- **Skills** = tools the agent can call (functions). The agent decides when.
- **Plugins** = hooks that run automatically (context engine, memory backend). The system decides when.

The codex needs **both**:
- **Skill** (`codex-1.0.0`) — agent-initiated operations (search, capture, organize)
- **Plugin** (future) — automatic operations (index new files, session resilience, reflection triggers)

For now, the skill is the right starting point. Plugin behaviors can run as lobster routines or heartbeat tasks until the plugin API is mature enough.

### Where the Best Features Live

| Feature | Lives In | Why There |
| --- | --- | --- |
| **Tiered loading** (OpenViking + Letta) | `_codex/skills/search.md` (instructions) | Agent-agnostic 4-tier protocol. Any AI can follow it. |
| **Hybrid search** (OpenClaw) | Runtime adapter (skill `index.js`) | Requires runtime search API. Can't be markdown. |
| **Knowledge graph** (ClawVault) | Codex itself (wiki-links + tags + frontmatter `related:`) | Graph is IN the content. Obsidian visualizes it. Agents traverse it. |
| **Session resilience** (ClawVault) | `_codex/skills/maintain.md` (instructions) + runtime (plugin/lobster) | Instructions are universal, execution is runtime-specific. |
| **Sleep-time reflection** (Letta) | `_codex/skills/reflect.md` (instructions) + lobster routine (execution) | Reflection protocol is universal. Trigger mechanism is runtime-specific. |
| **Fact extraction** (ClawVault) | `_codex/skills/consolidate.md` (instructions) | Any AI can follow "extract facts from daily notes." |
| **Progressive disclosure** (skills spec) | `_codex/INTERFACE.md` + folder READMEs | Every folder has a README = L0/L1 by default. |
| **Templates** | `_codex/templates/` | Travel with the vault. Any agent uses them. |
| **Domain weighting** | Runtime config (`openclaw.json`) | Per-agent preference. Not universal. |
| **Context profiles** | `_codex/profiles/` | Universal. Any agent can read "what to pre-load for X." |

**The principle**: If it's **knowledge or instructions** → vault (`_codex/`). If it **requires an API or runtime** → adapter (skill/plugin/MCP).

---

## 13. Making It a Product — OneMind Codex for Everyone

### The Vision

The OneMind Codex is not just Zeus's personal vault. It's a **framework + template** anyone can use — from a beginner with Obsidian to an advanced user running 9 AI agents.

### The Product Layers (Beginner → Advanced)

```
LAYER 0: FRAMEWORK (free, universal)
├── The CODEX Method (C.O.D.E.×)
├── The 00-99 domain numbering system
├── The 3 domains + 1 operating layer (SO / HP / LE / GE)
└── Published as: book, blog, README

LAYER 1: STARTER VAULT (free, git template)
├── Pre-built _codex/ folder with all skills, profiles, templates
├── Empty 00-99 domain folders with READMEs
├── Obsidian config (.obsidian/) tuned for the framework
├── Working out of the box — open in Obsidian, start using
└── Published as: GitHub template repo "onemind-codex"

LAYER 2: AI-READY VAULT (free, git template + adapters)
├── Everything from Layer 1
├── Runtime adapters (OpenClaw skill, Claude skill, Cursor rules, MCP server)
├── Pre-configured for agent access
└── Published as: same repo with adapters/ directory

LAYER 3: FULL ONEMIND OS (personal, this monorepo)
├── The codex (Layer 2)
├── OpenClaw runtime (9 agents, skills, NATS)
├── Infrastructure (docker-compose, VPS, services)
├── The complete Zeus implementation
└── Not published as product — this is the reference implementation
```

### The Open-Source Structure

```
onemind-codex/                          ← PUBLIC TEMPLATE REPO
├── README.md                           ← "Get started in 5 minutes"
├── FRAMEWORK.md                        ← The CODEX method explained
├── LICENSE                             ← MIT
├── _codex/                             ← The vault's brain
│   ├── INTERFACE.md                    ← Universal entry point
│   ├── CONVENTIONS.md                  ← Rules for this vault
│   ├── skills/                         ← Agent-agnostic operations
│   │   ├── capture.md
│   │   ├── organize.md
│   │   ├── search.md
│   │   ├── reflect.md
│   │   ├── consolidate.md
│   │   └── maintain.md
│   ├── profiles/                       ← Context modes
│   │   ├── planning.md
│   │   ├── review.md
│   │   └── incident.md
│   └── templates/                      ← Note templates
│       ├── daily.md
│       ├── decision.md
│       ├── person.md
│       ├── project.md
│       ├── sop.md
│       └── review.md
├── .obsidian/                          ← Pre-configured Obsidian
│   ├── app.json
│   ├── plugins/                        ← Essential plugins pre-installed
│   └── themes/
├── 00-24 SO/                           ← Domains (with README in each)
│   ├── README.md
│   ├── 00 Framework/
│   ├── 01 Command/
│   └── ...
├── 25-49 HP/
├── 50-74 LE/
├── 75-99 GE/
└── adapters/                           ← Runtime integration examples
    ├── README.md                       ← "How to connect your AI"
    ├── openclaw/                       ← OpenClaw skill (Node.js)
    │   ├── codex-1.0.0/
    │   │   ├── SKILL.md
    │   │   ├── index.js
    │   │   └── package.json
    │   └── README.md                   ← "Copy skill to your workspace"
    ├── claude/                         ← Claude Code integration
    │   ├── SKILL.md
    │   └── README.md
    ├── cursor/                         ← Cursor rules
    │   ├── codex.mdc
    │   └── README.md
    └── mcp/                            ← Universal MCP server
        ├── server.js                   ← MCP server for any client
        ├── package.json
        └── README.md
```

### What Makes It Beginner-Friendly

1. **Clone and open** — `git clone`, open folder in Obsidian, done
2. **READMEs everywhere** — every folder explains itself
3. **Templates built in** — `_codex/templates/` so new notes have structure from day one
4. **No code required** — Level 1 works with just Obsidian, no agents needed
5. **`_codex/skills/` are readable by humans** — even without AI, you can follow the protocols manually
6. **Progressive adoption** — start with Obsidian → add an agent when ready → add more agents later

### What Makes It Advanced

1. **`_codex/` is the agent interface** — any AI knows how to operate on the vault
2. **Runtime adapters included** — OpenClaw, Claude, Cursor, MCP out of the box
3. **Self-maintaining** — reflection, consolidation, health checks as codex skills
4. **Knowledge graph built into the content** — wiki-links, tags, `related:` frontmatter
5. **Context profiles** — modal awareness for different situations
6. **Fully git-backed** — version history, branching, collaboration, sync across devices

### The Git Strategy

```
onemind-codex (template)              ← The public product
    ↓ "Use this template" (GitHub)
my-codex (user's fork)                ← User's personal vault
    ↓ git clone
~/Documents/my-codex/                 ← Local Obsidian vault
    ↓ git push
GitHub (backup + sync)                ← Cross-device sync
    ↓ agent reads
AI agent (via adapter)                ← Agent has full context
```

Users own their vault. Git provides versioning, backup, and sync. No SaaS dependency. No lock-in. If Obsidian dies tomorrow, the markdown still works everywhere.

### How Zeus's Monorepo Relates to the Product

```
onemindos-v0.5/                       ← Zeus's private monorepo
├── codex/                            ← Zeus's PERSONAL codex
│   ├── _codex/                       ← SAME structure as the template
│   ├── 00-24 SO/                     ← Filled with Zeus's real data
│   └── ...
├── openclaw/                         ← Zeus's runtime (not in product)
├── infra/                            ← Zeus's infra (not in product)
└── ...

onemind-codex/                        ← The PUBLIC product (separate repo)
├── _codex/                           ← SAME _codex/ as Zeus's
├── 00-24 SO/                         ← Empty with READMEs only
├── adapters/                         ← Integration examples
└── ...
```

The `_codex/` folder is the shared DNA. Zeus's personal codex and the public template share the same operating system. When Zeus improves `_codex/skills/reflect.md`, that improvement can be pushed to the public template.

---

## 14. Implementation Plan — What to Build Now

### Phase 0: Create `_codex/` Inside Current Codex ~~(this week)~~ ✅ DONE

```text
codex/_codex/
├── INTERFACE.md           ← Universal agent entry point
├── CONVENTIONS.md         ← Naming, frontmatter (ai_summary required), git commits, self-org
├── system/                ← NEW — Pinned always-loaded context
│   └── README.md          ← Rules for what goes here
├── skills/
│   ├── capture.md         ← How to add new information
│   ├── organize.md        ← How to file into 00-99
│   ├── search.md          ← 4-tier loading protocol (System → L0 → L1 → L2)
│   ├── reflect.md         ← Reflection routines (daily/weekly/monthly + sleep-time)
│   ├── consolidate.md     ← Merge daily notes → knowledge
│   ├── maintain.md        ← Vault health operations (8-point audit)
│   ├── init.md            ← Vault bootstrap + agent onboarding
│   └── defrag.md          ← Structural reorganization (backup → restructure → verify)
├── profiles/
│   ├── planning.md
│   ├── review.md
│   └── incident.md
└── templates/
    ├── daily.md
    ├── decision.md
    ├── person.md
    ├── project.md
    ├── sop.md
    └── review.md
```

### Phase 1: Build the OpenClaw `codex` Skill ~~(next)~~ ✅ DONE

```text
openclaw/workspaces/legacy/skills/codex-1.0.0/
├── SKILL.md          ← Skill discovery + instructions
├── index.js          ← codex_capture, codex_search, codex_organize,
│                        codex_reflect, codex_consolidate, codex_health, codex_tree
└── package.json
```

The skill reads `_codex/` at boot to know the vault's conventions, then provides tools that follow those conventions.

### Phase 2: Add Adapters for Other Runtimes (pending)

- Claude Code: `.claude/skills/codex/SKILL.md` that points at `_codex/`
- Cursor: `.cursor/rules/codex.mdc` that references `_codex/`
- MCP server: `adapters/mcp/server.js` that exposes codex operations via MCP protocol

### Phase 3: Extract the Public Template (when ready)

- Strip personal content from codex
- Keep `_codex/`, empty domain folders with READMEs, adapters, Obsidian config
- Publish as `onemind-codex` GitHub template repo
- Write the getting-started guide

---

*Created: 2026-03-18*
*Updated: 2026-03-18*
*Sources: OpenClaw docs, Letta blog, OpenViking, ClawVault, MemFS, MemSearch, OneMind Codex*
*Domain: 12 Research (R&D)*

---

## 15. Additional Systems Evaluated

After the initial 6-system synthesis, we evaluated 4 more systems:

### 15.1 MemSearch (zilliztech/memsearch) — 929 stars

**What it is**: Python library for markdown-first agent memory, built by Zilliz (Milvus creators). Inspired by OpenClaw's memory model.

**Key features**:
- **Hybrid search**: Dense vectors + BM25 + Reciprocal Rank Fusion (RRF), powered by Milvus
- **Content-hash dedup**: SHA-256 hashing prevents storing the same content twice
- **Live file watcher**: Auto-indexes new/changed Markdown files via filesystem watcher
- **LLM summarization/compaction**: Automatically summarizes and compacts memory
- **Local embedding**: ONNX models (no API key needed) + OpenAI, Google, Voyage, Ollama
- **Multiple backends**: Milvus Lite (embedded), Milvus Server, Zilliz Cloud
- **Claude Code plugin**: Native integration available
- **Framework integrations**: LangChain, LangGraph, LlamaIndex, CrewAI

**Verdict**: EVALUATE LATER — Strong hybrid search engine. Worth testing if OpenClaw native search + ClawVault proves insufficient for codex-wide retrieval. The ONNX local embedding (no API key) is particularly interesting for self-hosted setups.

### 15.2 ClawVault v3.3-3.4 (Updated Assessment)

Since our initial evaluation, ClawVault has shipped significant updates:

**New in v3.3-3.4**:
- **Native OpenClaw plugin**: `openclaw.plugin.json` — installable via `openclaw hooks install clawvault`
- **In-process hybrid search**: BM25 + semantic + cross-encoder reranking (no external server needed)
- **Python SDK**: `clawvault-py` for programmatic access
- **Background workers**: 4 autonomous agents running inside the vault:
  - **Curator**: Promotes important information to permanent storage
  - **Janitor**: Cleans up stale/duplicate content
  - **Distiller**: Extracts facts and summaries from verbose content
  - **Surveyor**: Maps relationships and builds knowledge graph
- **Tailscale + WebDAV sync**: Cross-device vault synchronization
- **466 tests, 595 stars** (growing)

**Verdict**: USE — Install as OpenClaw hooks. The native plugin support means it integrates without custom wiring. The background workers (Curator/Janitor/Distiller/Surveyor) handle consolidation and maintenance automatically — complementing our manual `_codex/skills/consolidate.md` and `maintain.md`.

**Installation**: `openclaw hooks install clawvault && openclaw hooks enable clawvault && systemctl restart openclaw`

### 15.3 MemFS (streamich/memfs) — Updated Assessment

**What it is**: Pure in-memory Node.js `fs` API + browser File System API adapters. 2k stars, Apache 2.0.

**Verdict**: SKIP — This is a testing/browser utility library, not an agent memory system. The *pattern* (filesystem = memory primitive) was already absorbed into our architecture in §2 Pattern 1. The npm package itself adds nothing.

### 15.4 Letta Context Repositories — Structural Patterns

**What it is**: Blog post (Feb 2026) describing git-backed memory, local filesystem clones, progressive disclosure via frontmatter, memory swarms (concurrent subagents in git worktrees), and memory skills (init, sleep-time reflection, defragmentation).

**Verdict**: BORROW patterns, don't use runtime. Key patterns already absorbed into `_codex/skills/`. See §16 for the 5 structural improvements we extracted and implemented.

### 15.5 Final Verdict — What We Do With Each System

| System | Verdict | Action |
| --- | --- | --- |
| **OpenClaw native** | USE | Already running. Enhance config (session memory, temporal decay, MMR). |
| **Codex (Obsidian)** | USE | Already the vault. Enhanced with `_codex/` self-description layer. |
| **ClawVault v3.4** | USE | Install as OpenClaw hooks. Background workers handle auto-maintenance. |
| **Letta patterns** | BORROW | Sleep-time reflection, defrag, init, system/ pinned context — all implemented in `_codex/skills/`. |
| **MemSearch** | EVALUATE LATER | If OpenClaw + ClawVault search is insufficient, test MemSearch's Milvus-backed hybrid search. |
| **OpenViking plugin** | EVALUATE LATER | 52% accuracy gain claimed. Test against OpenClaw native search on real codex queries. |
| **lossless-claw** | EVALUATE LATER | Better compaction than legacy summarization. Test when context window pressure becomes a bottleneck. |
| **MemFS** | SKIP | Pattern absorbed. Package not needed. |
| **Letta runtime** | SKIP | Legacy stays in OpenClaw. We borrow Letta's design patterns, not its runtime. |

---

## 16. Structural Improvements from Letta Context Repositories

After close analysis of Letta's Context Repositories design, we identified and implemented 5 structural gaps in our `_codex/` layer:

### 16.1 Pinned Context: `_codex/system/`

**Letta pattern**: Files in `system/` are ALWAYS fully loaded into the agent's system prompt.

**Our implementation**: Created `_codex/system/` directory with rules:
- All files here are pre-loaded at the start of every agent interaction (~2000 token budget)
- Agents can promote/demote files to/from `system/` as they learn what's most relevant
- Owner can lock files with `pinned: true` in frontmatter
- Contains: active goals, identity context, critical reference data

**Impact**: Eliminates the most common search pattern — looking for "always-needed" context. Before, agents had to search for identity/goals every time. Now it's pre-loaded.

### 16.2 Required `ai_summary` in Frontmatter

**Letta pattern**: Every file has a description visible at directory level (frontmatter), enabling progressive disclosure without opening files.

**Our implementation**: Moved `ai_summary` from "recommended" to **required** in `_codex/CONVENTIONS.md`. Every file must have a ≤20-word summary in frontmatter.

**Impact**: The primary progressive disclosure mechanism. At L1 (Discovery tier), agents can read summaries of all files in a domain without opening any of them. Without this, agents must open every file to know what's inside.

### 16.3 Agent Self-Organization Permissions

**Letta pattern**: Agents can reorganize the file hierarchy, update frontmatter, and move files in/out of `system/` to control what's pinned to context as they learn over time.

**Our implementation**: Added explicit permission model to `_codex/CONVENTIONS.md`:

**Agents CAN**:
- Reorganize files within domains (move, rename, create subfolders)
- Update frontmatter (fix missing fields, improve `ai_summary`, update `modified:`)
- Manage `_codex/system/` (promote/demote files)
- Create subfolders when a domain clusters into 3+ subtopics with 5+ files
- Split oversized files (>300 lines)
- Fix broken `[[wiki-links]]`

**Agents CANNOT**:
- Create new top-level domain group folders (00-99 numbering is fixed)
- Delete files without owner confirmation (archive instead)
- Unpin `_codex/system/` files marked `pinned: true`
- Modify `_codex/` skills, templates, or conventions without explicit instruction

**Impact**: Agents are now empowered to maintain and improve the vault structure over time, not just read from it.

### 16.4 Memory Defragmentation Skill

**Letta pattern**: "Backs up the agent's memory filesystem, then launches a subagent that reorganizes files, splitting large files, merging duplicates, and restructuring into a clean hierarchy of 15-25 focused files."

**Our implementation**: Created `_codex/skills/defrag.md` — distinct from consolidation:

| Consolidate | Defragment |
| --- | --- |
| Merges scattered notes on same topic | Restructures the entire folder hierarchy |
| Operates on individual notes | Operates on whole domains or vault |
| Preserves existing structure | May create/remove subfolders |
| Bottom-up (note-by-note) | Top-down (architecture-first) |
| Safe to run anytime | Requires git backup first |

**Process**: Backup (git tag) → Audit target → Design new structure → Execute (merge → split → move → archive) → Update `system/` → Verify + commit

**Target**: 15-25 focused files per domain (Letta's recommendation).

### 16.5 Memory Initialization Skill

**Letta pattern**: "Bootstraps new agents by exploring the codebase and reviewing historical conversation data using concurrent subagents that work in git worktrees to create the initial hierarchical memory structure."

**Our implementation**: Created `_codex/skills/init.md` covering:
- **New vault init**: Skeleton → seed domains → import existing knowledge → populate `system/` → commit
- **Agent onboarding**: Read INTERFACE.md → load `system/` → orientation scan → create agent profile → register
- **Parallel bootstrapping**: For large imports (1000+ files), partition data → process chunks independently → merge → dedup (Letta fan-out pattern adapted for non-git-worktree environments)

### 16.6 Additional Improvements

**Git commit message convention** — Added to `_codex/CONVENTIONS.md`. Every vault change gets a descriptive commit using action prefixes: `capture:`, `update:`, `organize:`, `consolidate:`, `defrag:`, `reflect:`, `init:`, `archive:`, `system:`.

**Filetree-as-navigation** — Documented in `_codex/INTERFACE.md`. The folder hierarchy and file names serve as a navigational table of contents. Agents always have access to the directory listing — descriptive folder names (numbered domains) and descriptive file names (lowercase-hyphenated) make the filetree scannable without opening any files.

**4-tier search** — Updated `_codex/skills/search.md` from 3 tiers to 4:

| Tier | Name | Token Cost | What |
| --- | --- | --- | --- |
| **System** | Pinned context | 0 (pre-loaded) | `_codex/system/` files — always in context |
| **L0** | Orientation | ~200 tokens | Read domain group README |
| **L1** | Discovery | ~500 tokens | Scan file names + `ai_summary` frontmatter |
| **L2** | Deep read | Full file | Open specific files confirmed relevant |

### 16.7 Updated `_codex/` Structure (Current State)

```
codex/_codex/
├── INTERFACE.md              ← Universal agent entry point (updated: system/, self-org, filetree nav)
├── CONVENTIONS.md            ← Rules (updated: ai_summary required, git commits, self-org permissions)
├── system/                   ← NEW — Pinned always-loaded context
│   └── README.md             ← Rules for what goes here
├── skills/
│   ├── search.md             ← Updated: 4-tier (System → L0 → L1 → L2)
│   ├── capture.md            ← How to add new information
│   ├── organize.md           ← How to file into 00-99
│   ├── reflect.md            ← Reflection cycles (daily/weekly/monthly + sleep-time)
│   ├── consolidate.md        ← Merge scattered notes → curated knowledge
│   ├── maintain.md           ← Vault health audit (8-point checklist)
│   ├── init.md               ← NEW — Vault bootstrap + agent onboarding
│   └── defrag.md             ← NEW — Structural reorganization
├── templates/
│   ├── daily.md              ├── project.md       ├── decision.md
│   ├── person.md             ├── review.md        └── sop.md
└── profiles/
    ├── planning.md           ├── review.md        └── incident.md
```

**Total**: 20 files in `_codex/`, covering the complete vault operating system.
