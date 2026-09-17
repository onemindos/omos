# QMD + Codex Memory Architecture
**Domain:** SO — 01 Intelligence (AI, Memory & Tools)
**Maintained by:** ⚡ Legacy
**Last updated:** 2026-04-18
**Audience:** Zeus Delacruz + One Mind OS students

---

## The Core Idea

> **The Codex is the brain. QMD is what makes it searchable. The agents are what act on it.**

Before this system existed, AI agents had no persistent memory. Every conversation started blank. The agent knew nothing about you unless you told it again — every single time.

One Mind OS solves this with a three-layer architecture:

```
YOUR LIFE (notes, docs, plans, logs)
        ↓  written as markdown
THE CODEX (structured knowledge vault)
        ↓  indexed automatically
QMD (semantic memory engine)
        ↓  searched on every recall
AGENTS (act on your behalf with real context)
```

The result: agents that **remember**, **reason from context**, and **act with continuity** — just like a human executive assistant who has read every file in your office.

---

## Layer 1 — The Codex

The Codex is a **plain markdown knowledge vault** organized by life domain. It is not an app. It is not a database. It is structured text files that any intelligence — human or AI — can read and act on.

### Structure

```
/root/onemind-codex/
├── 00-24 SO (Sovereign Operations)    ← AI, infrastructure, protocols
├── 25-49 HP (Holistic Performance)    ← Zeus's body, mind, money, skills
├── 50-74 LE (Legacy Evolution)        ← Family, home, relationships, community
└── 75-99 GE (Generational Empire)  ← Business, brand, ventures, wealth
```

Each domain contains numbered sub-domains (e.g., `27 Body`, `30 Finance`, `75 Brand`). Everything in your life has a place.

### Why Markdown?

- **Universal** — any tool can read it. Obsidian, VS Code, GitHub, terminal, AI.
- **Permanent** — markdown files outlive every app. Notion dies. Roam dies. Markdown lives.
- **AI-native** — language models were trained on text. Markdown is the closest thing to their native format.
- **Version controlled** — every change tracked in git. Full history forever.

### The CODEX Method

```
C — Capture    Put it in the codex (don't leave it in your head or a chat)
O — Organize   Route it to the right domain/sub-domain
D — Direct     Tell agents what to do with it
E — Execute    Agents act on it
X — Multiply   Results build on each other over time
```

---

## Layer 2 — QMD (Query Memory Database)

QMD is OpenClaw's built-in **semantic memory engine**. It runs silently in the background and does two things:

1. **Indexes** your codex files into a searchable database
2. **Searches** that database when an agent needs to recall something

### How It Works

```
File saved/changed in codex
        ↓
QMD file watcher detects the change (within seconds)
        ↓
QMD re-indexes the file (updates its database)
        ↓
Agent runs memory_search("training program")
        ↓
QMD searches indexed database
        ↓
Returns relevant snippets with file paths + line numbers
        ↓
Agent uses that context to respond accurately
```

### Three Search Modes

QMD supports three progressively deeper search modes:

| Mode | How It Works | Best For |
|------|-------------|----------|
| **Lexical** | Keyword matching (like grep but smarter) | Exact terms, names, dates |
| **Vector** | Semantic/meaning matching via embeddings | Concepts, topics, themes |
| **HyDE** | Hypothetical Document Embedding — AI imagines what the answer looks like, then finds matching content | Deep reasoning, abstract questions |

In practice, QMD runs all three and combines the results for maximum recall.

### The Database

QMD uses **SQLite + vector embeddings** stored locally on the VPS. No cloud. No third-party API. Your data stays on your machine.

- Index location: `~/.local/share/qmd/` (per agent)
- Shared ML models: `~/.cache/qmd/models/` (one copy, all agents share it)
- Auto-rebuilds if files change
- Handles 15,000+ files without performance issues

---

## Layer 3 — Agent Memory Scoping

Each agent's QMD is configured to index only the codex domain relevant to their work. This is intentional.

### Why Scope It?

Without scoping — every agent searches all 15,000+ files. A search for "training plan" by Spartan would return noise from the SO infrastructure domain, GE business docs, and LE family files.

With scoping — Spartan only searches HP. Empire only searches GE. Every search is fast and returns relevant results.

### Domain Scoping Map

```
Agent          QMD searches this codex path
──────────────────────────────────────────────────────────────
⚡ Legacy      FULL codex + all 8 agent workspaces (orchestrator)
💪 Spartan     25-49 HP (Holistic Performance)
🔮 Oracle      25-49 HP (Holistic Performance)
🏡 Haven       50-74 LE (Legacy Evolution)
🌱 Eden        50-74 LE (Legacy Evolution)
💰 Empire      75-99 GE (Generational Empire)
🛡️ Guardian    SO/04 Security (Cyber)
⚙️ Grid        SO/03 Infrastructure + SO/03 Command
🔺 Trinity     SO/00 Framework + SO/03 Command
```

### Two Access Layers (Important Distinction)

QMD scoping only affects **automatic memory recall**. It does not restrict file access.

| Layer | What It Controls | Scope |
|-------|-----------------|-------|
| **QMD memory search** | What agents recall automatically | Domain-scoped (their lane) |
| **codex-vault MCP** | Manual file read/write | Full codex — when directed |

Think of it this way:
- **QMD** = instinctive memory. Always on. Searches their lane automatically.
- **codex-vault** = manual file access. Full reach. Only used when you or Legacy directs them.

---

## The Full Loop — How It All Connects

```
Zeus writes a note in Obsidian (Mac)
        ↓
Git commit + push (onemind-watch service, 3s debounce)
        ↓
GitHub (onemind-codex repo — private)
        ↓
onemind-pull.timer pulls to VPS (every 10 seconds)
        ↓
QMD file watcher detects change
        ↓
QMD re-indexes the file
        ↓
Next time agent runs memory_search → file is in the index
        ↓
Agent responds with full context of what Zeus wrote
```

**End-to-end latency: ~15-30 seconds** from Zeus writing a note to it being searchable by any agent.

---

## What This Enables

### Before One Mind OS
- Every AI conversation starts blank
- You re-explain your context every time
- AI gives generic answers
- No continuity between sessions

### After One Mind OS
- Agents wake up knowing your domains
- Spartan knows your current training program, injury history, goals
- Oracle knows your financial position, active legal matters, learning stack
- Empire knows your brand strategy, product roadmap, revenue targets
- Haven knows your family situation, home projects, upcoming events

**The agents don't just assist. They operate from the same context you do.**

---

## Sync Architecture

Two systemd services keep the codex in sync:

| Service | What It Does | Interval |
|---------|-------------|----------|
| `onemind-watch.service` | VPS → GitHub. Watches for file changes via inotifywait. Auto git add/commit/push. | 3 second debounce |
| `onemind-pull.timer` | GitHub → VPS. Pulls if remote has new commits. | Every 10 seconds |

This means:
- Zeus edits on Mac → pushes to GitHub → VPS pulls within 10s → QMD indexes within seconds
- Legacy writes to codex on VPS → onemind-watch pushes within 3s → Mac syncs on next pull

Both directions are automatic. The codex is always in sync.

---

## Configuration Reference

QMD + codex paths are configured in `~/.openclaw/openclaw.json`:

```json
// Global default — applies to all agents unless overridden
"agents": {
  "defaults": {
    "memorySearch": {
      "extraPaths": ["/root/onemind-codex"]
    }
  },
  "list": [
    {
      "id": "spartan",
      "memorySearch": {
        "extraPaths": [
          "/root/onemind-codex/25-49 HP (Holistic Performance)"
        ]
      }
    }
    // ... one entry per agent
  ]
}
```

To add a new path to an agent's memory — add it to `extraPaths`.
To give an agent cross-domain recall — add multiple paths.
Legacy already has the full codex root, so it sees everything.

---

## Key Principles for Students

1. **The codex is not an app — it's a protocol.** Any AI, any tool, any future system can plug into it.

2. **Write everything down.** An agent can only recall what exists in the codex. If it's not written, it doesn't exist for the AI.

3. **Domain discipline matters.** Put files in the right domain folder. QMD scoping only works if your content is organized correctly.

4. **QMD is automatic.** You don't trigger it. Write the file, QMD handles the rest.

5. **Agents + codex = continuity.** The goal is a system that compounds over time. Every note, every plan, every decision you write becomes permanent context that every agent can act on — today and 10 years from now.

---

## Related Documents

- `agent-skills-architecture.md` — Tool and skill access for all 9 agents
- `CODEX-FRAMEWORK.md` (codex root) — The CODEX method and One Mind OS doctrine
- `ONEMIND-CODEX.md` (codex root) — Full codex structure reference
- `~/.openclaw/openclaw.json` — Live QMD configuration
