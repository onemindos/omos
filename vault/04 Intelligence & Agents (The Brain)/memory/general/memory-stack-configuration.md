---
title: "OpenClaw Memory Stack Configuration"
type: sop
status: active
created: 2026-04-12
updated: 2026-04-12
owner: zeus
tags: [type/sop, domain/so, domain/agents, status/active]
ai_summary: "Complete configuration reference for the 5-layer OpenClaw memory stack: lossless-claw (LCM context engine), QMD (vector search backend), dreaming (memory consolidation), active-memory, and memory-wiki. Includes exact JSON paths, verified field names, and how the layers interact."
---

# OpenClaw Memory Stack — Configuration Reference

> VPS: `root@104.248.4.212` | Config: `~/.openclaw/openclaw.json` | Activated: 2026-04-12

---

## Architecture Overview

Five memory layers stack on top of each other. Each operates independently but they feed one another.

```
LAYER 5 — memory-wiki (bridge mode)
  ↑ reads dream reports, daily notes, MEMORY.md from memory-core
  ↑ compiles into obsidian-compatible wiki vault

LAYER 4 — Dreaming (memory consolidation)
  ↑ nightly sweep: light → REM → deep phases
  ↑ deep phase writes durable entries to MEMORY.md
  ↑ runs inside memory-core plugin at 03:00 UTC by default

LAYER 3 — QMD (vector search backend)
  ↑ indexes the full codex vault (/home/onemind/onemind-os/codex/)
  ↑ provides semantic search across all markdown files
  ↑ replaces builtin SQLite memory search

LAYER 2 — active-memory (recent context injection)
  ↑ queries session history before each turn
  ↑ injects compact preference summary into system prompt
  ↑ enabled for all 9 VPS agents

LAYER 1 — lossless-claw (context compaction)
  ↑ replaces OpenClaw's default sliding-window truncation
  ↑ DAG-based summarization — nothing is ever lost
  ↑ contextEngine slot owner (exclusive)
```

---

## Complete `openclaw.json` Config

The full config at `~/.openclaw/openclaw.json` after activation:

```json
{
  "agents": {
    "defaults": {
      "model": { "primary": "github-copilot/gpt-4o" },
      "compaction": { "mode": "safeguard" },
      "memorySearch": {
        "extraPaths": ["/home/onemind/onemind-os/codex/"]
      }
    }
  },
  "memory": {
    "backend": "qmd",
    "qmd": {
      "update": { "onBoot": true },
      "sessions": { "enabled": true }
    }
  },
  "plugins": {
    "slots": {
      "contextEngine": "lossless-claw"
    },
    "entries": {
      "active-memory": {
        "enabled": true,
        "config": {
          "agents": ["legacy", "spartan", "oracle", "eden", "haven", "empire", "guardian", "grid", "trinity"],
          "queryMode": "recent",
          "promptStyle": "preference-only",
          "timeoutMs": 15000,
          "maxSummaryChars": 220,
          "logging": true
        }
      },
      "lossless-claw": {
        "enabled": true,
        "config": {
          "freshTailCount": 64,
          "leafChunkTokens": 80000,
          "contextThreshold": 0.75,
          "incrementalMaxDepth": 1,
          "summaryModel": "github-copilot/gpt-4o"
        }
      },
      "memory-core": {
        "enabled": true,
        "config": {
          "dreaming": { "enabled": true }
        }
      },
      "memory-wiki": {
        "enabled": true,
        "config": {
          "vaultMode": "bridge",
          "vault": {
            "path": "/home/onemind/onemind-os/codex/_codex/memory/wiki",
            "renderMode": "obsidian"
          },
          "bridge": {
            "enabled": true,
            "readMemoryArtifacts": true,
            "indexDreamReports": true,
            "indexDailyNotes": true,
            "indexMemoryRoot": true,
            "followMemoryEvents": true
          },
          "context": {
            "includeCompiledDigestPrompt": true
          }
        }
      }
    }
  }
}
```

---

## Layer Details

### 1. lossless-claw (v0.8.2)

**What it does:** Replaces OpenClaw's default sliding-window context truncation with a DAG-based summarization system. Every message is persisted to SQLite; nothing is ever lost.

**Install:** `openclaw plugins install @martian-engineering/lossless-claw`
**Location:** `/root/.openclaw/extensions/lossless-claw`
**Slot:** `plugins.slots.contextEngine: "lossless-claw"` (exclusive)

**Key config fields:**

| Field | Value | Meaning |
|-------|-------|---------|
| `freshTailCount` | 64 | Recent messages protected from compaction |
| `leafChunkTokens` | 80000 | Source tokens per leaf before summarization triggers |
| `contextThreshold` | 0.75 | % of context window that triggers compaction |
| `incrementalMaxDepth` | 1 | How deep per-turn incremental compaction goes (1 = one condensed pass) |
| `summaryModel` | `github-copilot/gpt-4o` | Model used for summarization |

**Agent commands:** `/lcm` (status), `/lcm doctor` (scan broken summaries), `/lossless` (alias)

**DB location:** `~/.openclaw/lcm.db`

---

### 2. active-memory

**What it does:** Before each agent turn, queries recent session history for relevant preferences and injects a compact summary into the system prompt.

**Config path:** `plugins.entries["active-memory"].config`

**Key config fields:**

| Field | Value | Meaning |
|-------|-------|---------|
| `agents` | all 9 | Which agents participate |
| `queryMode` | `recent` | Query strategy (recent / message / full) |
| `promptStyle` | `preference-only` | Injection style: balanced / strict / contextual / recall-heavy / precision-heavy / preference-only |
| `timeoutMs` | 15000 | Max wait for memory query |
| `maxSummaryChars` | 220 | Max injected chars |
| `logging` | true | Log active-memory events |

---

### 3. QMD — Query Memory Daemon (v2.1.0)

**What it does:** Vector search backend across the full codex vault. Replaces OpenClaw's builtin SQLite memory search with semantic search over `/home/onemind/onemind-os/codex/`.

**Binary:** `/usr/bin/qmd` (installed via `npm install -g @tobilu/qmd`)
**Activation:** `memory.backend: "qmd"` at top level of openclaw.json

**How the codex path is wired:**
`agents.defaults.memorySearch.extraPaths: ["/home/onemind/onemind-os/codex/"]` — this path is automatically included as a custom QMD collection when `memory.backend == "qmd"`. No duplicate config needed in `memory.qmd.paths`.

**Key config fields in `memory.qmd`:**

| Field | Value | Meaning |
|-------|-------|---------|
| `update.onBoot` | true | Sync collections on gateway startup |
| `sessions.enabled` | true | Also index session transcripts |

**Index location:** `{agentStateDir}/qmd/xdg-cache/qmd/index.sqlite`

**Note:** First boot triggers QMD embedding run (downloads model if not cached, ~2GB). Subsequent boots are incremental.

---

### 4. Dreaming (memory consolidation)

**What it does:** Nightly background sweep (default 03:00 UTC) that runs three phases:
- **Light** — quick surface-level pattern detection
- **REM** — identifies recurring patterns across sessions
- **Deep** — writes durable entries to `MEMORY.md` (the only phase that persists to disk)

**Activation:** `plugins.entries["memory-core"].config.dreaming.enabled: true`

**Control commands:** `/dreaming on`, `/dreaming off`, `/dreaming status`

**Output files:**
- `MEMORY.md` — permanent memory (deep phase writes here)
- `DREAMS.md` — human-readable summaries and diary entries

**Location of these files:** Inside each agent's workspace at `/home/onemind/onemind-os/openclaw/workspaces/{agent}/`

---

### 5. memory-wiki (bridge mode)

**What it does:** Reads `memory-core` output (MEMORY.md, dream reports, daily notes) and compiles a structured wiki vault in Obsidian-compatible format. Agents can then reference this wiki for context injection.

**Vault path:** `/home/onemind/onemind-os/codex/_codex/memory/wiki`

**Config fields:**

| Field | Value | Meaning |
|-------|-------|---------|
| `vaultMode` | `bridge` | Read from memory-core, compile into wiki (not direct write) |
| `vault.path` | see above | Where the compiled wiki lives |
| `vault.renderMode` | `obsidian` | Format compatible with Obsidian wiki-links |
| `bridge.readMemoryArtifacts` | true | Include MEMORY.md |
| `bridge.indexDreamReports` | true | Include dream reports |
| `bridge.indexDailyNotes` | true | Include daily notes |
| `bridge.indexMemoryRoot` | true | Include memory root folder |
| `bridge.followMemoryEvents` | true | Update wiki when memory-core writes new events |
| `context.includeCompiledDigestPrompt` | true | Inject compiled digest into agent context |

---

## Installation History

| Step | What | When | Status |
|------|------|------|--------|
| Backup config | `openclaw.json.bak.20260412-090153` | 2026-04-12 | ✅ |
| Install lossless-claw | `openclaw plugins install @martian-engineering/lossless-claw` | 2026-04-12 | ✅ v0.8.2 |
| Install QMD | `npm install -g @tobilu/qmd` | 2026-04-12 | ✅ v2.1.0 |
| Write config | All 5 systems, verified | 2026-04-12 | ✅ |
| Create wiki dir | `_codex/memory/wiki/` on VPS | 2026-04-12 | ✅ |
| Restart service | `systemctl restart openclaw` | 2026-04-12 | ✅ active |

---

## How the Layers Interact

```
New conversation turn
  ↓
active-memory queries session history → injects preference summary into system prompt
  ↓
lossless-claw assembles context (recent raw msgs + DAG summaries) → fits in context window
  ↓
Agent processes turn using: injected preferences + lossless context + QMD semantic search
  ↓
After turn: lossless-claw compacts if threshold hit (contextThreshold: 0.75)
  ↓
Nightly: dreaming runs → writes to MEMORY.md → memory-wiki picks up → compiles wiki
  ↓
QMD boot: re-indexes updated codex vault → better semantic search on next turn
```

---

## Troubleshooting

### Check service status
```bash
ssh root@104.248.4.212 'systemctl status openclaw --no-pager -l | head -20'
```

### Check recent logs
```bash
ssh root@104.248.4.212 'journalctl -u openclaw -n 50 --no-pager'
```

### Verify plugins loaded
```bash
ssh root@104.248.4.212 'openclaw plugins list'
```

### Check lossless-claw status (from agent chat)
```
/lcm
/lcm doctor
```

### Check dreaming status (from agent chat)
```
/dreaming status
```

### Restart service
```bash
ssh root@104.248.4.212 'systemctl restart openclaw'
```

### Config backup location
```
/root/.openclaw/openclaw.json.bak.20260412-090153
```

---

## Config Key Reference

Quick-look table of all config paths:

| What | JSON path |
|------|-----------|
| QMD backend switch | `memory.backend` |
| QMD boot sync | `memory.qmd.update.onBoot` |
| QMD session indexing | `memory.qmd.sessions.enabled` |
| Codex path for QMD | `agents.defaults.memorySearch.extraPaths[]` |
| Compaction trigger | `agents.defaults.compaction.mode` |
| Context engine slot | `plugins.slots.contextEngine` |
| active-memory agents | `plugins.entries["active-memory"].config.agents` |
| lossless-claw model | `plugins.entries["lossless-claw"].config.summaryModel` |
| lossless-claw threshold | `plugins.entries["lossless-claw"].config.contextThreshold` |
| Dreaming on/off | `plugins.entries["memory-core"].config.dreaming.enabled` |
| Wiki vault path | `plugins.entries["memory-wiki"].config.vault.path` |
| Wiki mode | `plugins.entries["memory-wiki"].config.vaultMode` |


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
