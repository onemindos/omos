# CODEX-INTEGRATION.md — The Full Architecture

> **OneMind OS: 3 layers, one brain. How the Codex, runtime, and Mission Control connect.**

---

## 1. The SQLite Question — Definitive Answer

**Short answer: Yes, OpenClaw uses SQLite — but not the way you're thinking.**

OpenClaw's SQLite lives at `~/.openclaw/memory/<agentId>.sqlite` and is used exclusively for:
- **Vector embedding index** of MEMORY.md + daily logs + any `extraPaths` you configure
- **BM25 full-text search** table alongside vectors
- **Embedding cache** (so unchanged chunks aren't re-embedded on every sync)
- Optionally: session transcript indexing (experimental)

It is **not** a general-purpose app database. You cannot point a Next.js Mission Control at it. But here's the beautiful part:

> The `memorySearch.extraPaths` config lets you add **any directory of Markdown files** to that same SQLite index. Point it at `/home/onemind/codex/` and Legacy gets native BM25+vector semantic search over the **entire Codex vault** — zero custom code, zero extra skills, just config.

This is already done in `openclaw-vps.json`.

---

## 2. The Three-Layer Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│  LAYER 1 — CODEX (Knowledge)                                        │
│  The brain's long-term memory. Zeus's life. Single source of truth. │
│                                                                     │
│  Source: Mac (Obsidian editor)                                      │
│  Sync:   git pull cron every 5 min (GitHub → VPS)                  │
│  VPS:    /home/onemind/codex/ (read-only for all agents)           │
│  Backup: GitHub via Obsidian Git plugin (also the sync source)      │
│  Access: memory_search (native OpenClaw) + read_file                │
└─────────────────────────────────────────────────────────────────────┘
                              ↕ git pull (cron, every 5 min)
┌─────────────────────────────────────────────────────────────────────┐
│  LAYER 2 — OS RUNTIME (OpenClaw)                                    │
│  Legacy's brain. Runs on VPS. Everything routes through here.       │
│                                                                     │
│  Gateway: port 18789, VPS 134.122.8.135 (Tailscale: 100.104.61.29) │
│  Agents:  legacy (GPT-4o) | hp | le | ge (Claude)                  │
│  Memory:  MEMORY.md + memory/YYYY-MM-DD.md + Codex (via extraPaths) │
│  Channels: BlueBubbles (iMessage) | Discord | Voice (ElevenLabs)   │
└─────────────────────────────────────────────────────────────────────┘
                              ↕ convex-sync skill (future)
┌─────────────────────────────────────────────────────────────────────┐
│  LAYER 3 — MISSION CONTROL (Convex + Next.js)                       │
│  The dashboard Zeus sees. Real-time. OpenClaw ecosystem native.     │
│                                                                     │
│  Stack:   Convex (reactive DB) + Next.js (UI)                       │
│  Legacy writes to Convex via convex-sync skill on every turn        │
│  Zeus sees: tasks, agent status, schedule, memory cards, today log  │
│  Apply:   convex.dev/claw/apply (free credits for Claw projects)    │
└─────────────────────────────────────────────────────────────────────┘
                              ↕ Cloudflare R2 + Codex index
┌─────────────────────────────────────────────────────────────────────┐
│  LAYER 4 — MEDIA (R2 + Codex Index)                                 │
│  Raw files in R2. Metadata in Codex. Legacy reads index, not blobs. │
│                                                                     │
│  Storage: Cloudflare R2 ($0.015/GB/mo, stable CDN URLs)            │
│  Index:   Codex 54 Media Index/ (markdown files with R2 URLs)       │
│  Heavy:   scp/rclone copies important PDFs to VPS → nano-pdf reads  │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 3. How Legacy Accesses the Codex (3 Tiers)

### Tier 1: Semantic Search (use for 90% of Codex queries)

`memory_search` is built into OpenClaw. It now indexes the entire Codex.

```
memory_search "morning protocol"
→ returns: HP/27 Body/morning-protocol.md:12-45 (score: 0.94)
→ snippet: "Wake at 5:30. Cold shower. Journal. 10m meditation..."

memory_search "kids school Ash"
→ returns: LE/51 Family/ash-profile.md:1-20 (score: 0.91)
→ snippet: "Ash, age 9. School: Jefferson Elementary. Bus pickup 7:45am..."
```

Features active (configured in `openclaw-vps.json`):
- **Hybrid BM25+vector**: exact ID/code matches + semantic understanding
- **MMR**: removes near-duplicate results across daily notes
- **Temporal decay**: recent notes outrank older ones (30-day half-life)
- **Embedding cache**: 50k entries, no re-embedding unchanged chunks

### Tier 2: Direct File Read (use when you have the path)

```
read_file /home/onemind/codex/75-99 GE (Generational Endeavor)/75 Ventures/community-playbook.md
```

Use when `memory_search` returns a path you want to read in full.

### Tier 3: Obsidian Skill (for structured vault operations)

```
obsidian search "topic"          → searches note titles
obsidian create "note title"     → creates in Inbox for triage
obsidian move "note" "LE/51"     → moves to correct domain
```

---

## 4. Real-Life Codex Use Cases

### Family & Kids (Legacy Evolution — 50-74 LE)

```
Zeus: "What time does Ash get out of school?"
Legacy: memory_search "Ash school schedule"
→ reads LE/51 Family/ash-profile.md
→ replies: "Ash gets out at 3:15pm from Jefferson Elementary. Bus home unless you pick up."
```

### Health & Protocols (Holistic Performance — 25-49 HP)

```
Zeus: "What are my macros today?"
Legacy: memory_search "nutrition macros targets"
→ reads HP/27 Body/nutrition.md
→ replies: "Your targets: 2,800 kcal, 220g protein, 280g carbs, 80g fat. Today's tracking in Cronometer."
```

### Recipes (anywhere in Codex)

```
Zeus: "Make me a shopping list for the chicken tikka recipe"
Legacy: memory_search "chicken tikka recipe"
→ finds LE/53 Home/recipes/chicken-tikka.md
→ reads it, generates shopping list
```

### Business Protocols (Generational Endeavor — 75-99 GE)

```
Zeus: "What's the onboarding flow for new Circle members?"
Legacy: memory_search "Circle member onboarding"
→ reads GE/79 Community/onboarding-protocol.md
→ walks Zeus through the steps
```

### Infrastructure (Sovereign Operations — 00-24 SO)

```
Zeus: "What's the n8n webhook URL for Todoist sync?"
Legacy: memory_search "n8n todoist webhook"
→ reads UI/09 Automation/todoist-sync.md
→ returns exact URL + auth headers
```

---

## 5. Setup: Git Pull Sync (Codex → VPS)

> Simplest possible setup. Obsidian Git already pushes to GitHub. VPS just pulls. Zero config, zero conflicts.
>
> **Why no conflicts:** VPS is READ-ONLY. Only Mac writes to the Codex. VPS never commits.

### How it works

```
Mac (Obsidian) → edit notes → Obsidian Git plugin auto-commits → GitHub
                                                                     ↓
                                              VPS cron: git pull every 5 min
                                                                     ↓
                                              /home/onemind/codex/ (all agents read from here)
```

### One-time setup on VPS

```bash
# SSH into VPS
ssh openclaw@100.104.61.29

# Clone the Codex once
git clone https://github.com/Zeus-Delacruz/OneMind-Codex.git /home/onemind/codex

# Set up auto-pull every 5 minutes
crontab -e
# Add this line:
*/5 * * * * cd /home/onemind/codex && git pull --ff-only >> /home/onemind/codex-sync.log 2>&1
```

That's it. Done. No pairing, no devices, no UI to manage.

### Verify sync is working

```bash
# Check last pull
tail -5 /home/onemind/codex-sync.log
# Should show "Already up to date." or commits pulled

# Check files are there
ls /home/onemind/codex/
# → 00-24 SO (Sovereign Operations)/  25-49 HP/  50-74 LE/  75-99 GE/  ONEMIND-CODEX.md
```

### Notes on obsidian-cli

The `obsidian` skill is **optional**. You don't need it for agents to read/search the Codex. Only install it if you specifically want agents to create or move notes inside the vault structure. `memory_search` + `read_file` cover 99% of use cases.



---

## 6. Config Changes Applied

### `openclaw-vps.json` — what was added

Under `agents.defaults`:

```json
"memorySearch": {
  "enabled": true,
  "extraPaths": [
    { "path": "/home/onemind/codex", "pattern": "**/*.md" }
  ],
  "query": {
    "hybrid": {
      "enabled": true,
      "vectorWeight": 0.7,
      "textWeight": 0.3,
      "mmr": { "enabled": true, "lambda": 0.7 },
      "temporalDecay": { "enabled": true, "halfLifeDays": 30 }
    }
  },
  "cache": { "enabled": true, "maxEntries": 50000 }
}
```

**Effect**: All 4 agents (legacy, hp, le, ge) now have semantic search over the entire Codex via the native `memory_search` tool. No custom code needed.

### `workspace/AGENTS.md` — what was added

- OneMind Codex section with domain map
- memory_search usage guide with real examples
- Domain file paths by question type
- "Never invent facts" rule → search first

### `workspace/TOOLS.md` — what was added

- Codex Access section moved to top
- memory_search examples by domain
- read_file with real path patterns
- obsidian skill clarification (create/move, not search)
- Domain quick reference table

---

## 7. Mission Control — The Full Plan

### Why Convex (not OpenClaw's SQLite)

OpenClaw's SQLite is an **internal vector index** — not queryable from outside the process. Mission Control is a Next.js app running separately. It needs a real-time database it can subscribe to. Convex is purpose-built for this and is the OpenClaw ecosystem's official partner.

**Apply for Claw program credits**: `convex.dev/claw/apply`
- Free up to 1 year Professional plan (~$25/mo value)
- Project clearly qualifies (OpenClaw-native AI system)

### Architecture

```
[Legacy agent turn]
    ↓ after every response
[convex-sync skill] → writes to Convex
    ├── Tasks table: create/update tasks per domain
    ├── Memory table: today's log summary card
    ├── AgentStatus table: which agent, what topic, when
    ├── Schedule table: upcoming events from cron jobs
    └── ContextLog table: what Zeus talked about, linked Codex note

[Next.js Mission Control] → subscribes to Convex
    ├── /tasks: 4-domain kanban, all agents as assignees
    ├── /memory: daily journal cards + MEMORY.md summary
    ├── /team: legacy | hp | le | ge — live status
    ├── /schedule: upcoming events + heartbeat results
    └── /office: real-time "who Legacy is talking to"
```

### convex-sync Skill (to be built)

Location: `~/.openclaw/workspace/skills/convex-sync/`

```javascript
// convex-sync/index.js
const { ConvexHttpClient } = require("convex/browser");

const client = new ConvexHttpClient(process.env.CONVEX_URL);

module.exports = {
  name: "convex-sync",
  description: "Write structured data to Mission Control (Convex)",
  tools: {
    sync_task: async ({ title, domain, priority, agentId, codexPath }) => {
      await client.mutation("tasks:upsert", { title, domain, priority, agentId, codexPath });
      return `Task synced: ${title}`;
    },
    sync_memory_card: async ({ date, summary, keyFacts }) => {
      await client.mutation("memory:upsertCard", { date, summary, keyFacts });
      return `Memory card synced for ${date}`;
    },
    sync_agent_status: async ({ agentId, status, topic, startedAt }) => {
      await client.mutation("agents:updateStatus", { agentId, status, topic, startedAt });
      return `Agent status updated: ${agentId} → ${status}`;
    }
  }
};
```

**Building it**: Give Legacy the prompt below → she uses `skill-creator` to build it.

### Prompt for Legacy to Build Mission Control

```
Build a Mission Control system for OneMind OS. Stack: Convex + Next.js.

Phase 1 — Convex schema + convex-sync skill:
1. Set up Convex project at convex.dev
2. Define schema: tasks, memory_cards, agent_status, schedule, context_log tables
3. Create convex-sync OpenClaw skill at ~/.openclaw/workspace/skills/convex-sync/
4. Skill tools: sync_task, sync_memory_card, sync_agent_status, sync_schedule_event
5. Auth: CONVEX_URL + CONVEX_DEPLOY_KEY in ~/.openclaw/env

Phase 2 — Next.js Mission Control UI:
1. Scaffold at ~/mission-control/ (Next.js 14 + Convex + Tailwind + shadcn/ui)
2. Pages: /tasks (4-domain kanban), /memory (daily cards), /team (agent status), /schedule, /office
3. Real-time via Convex useQuery hooks (no polling)
4. Deploy to VPS via PM2 or Docker, port 3030
5. Auth: simple token gate from MISSION_CONTROL_TOKEN env var

Use skill-creator to build the OpenClaw skill, coding-agent for the Next.js app.
Reference Mission Control HQ (missioncontrolhq.ai) for UI inspiration.
```

---

## 8. Media Files Strategy

### Problem

Obsidian vaults + media files = disaster. Binary blobs in git destroy performance. Use R2 for storage, Codex for metadata.

### Solution

```
Raw files → Cloudflare R2 ($0.015/GB/mo, CDN speeds, stable URLs)
Metadata  → Codex 54 Media Index/ (markdown with R2 URLs + YAML frontmatter)
Heavy PDFs → scp/rclone to VPS → nano-pdf skill reads locally
```

### Codex Media Index Structure

```
50-74 LE (Legacy Evolution)/
└── 54 Media Index/          ← NEW folder to create
    ├── README.md             ← explains the index system
    ├── photos/
    │   ├── family/           ← one .md per album
    │   └── travel/
    ├── videos/
    ├── documents/
    └── audio/
```

Each entry looks like:

```markdown
---
type: photo-album
date: 2024-12-25
title: Christmas 2024 — Family
people: [ash, zeus, kai]
location: Miami home
r2_url: https://pub-xxx.r2.dev/photos/christmas-2024/
cover: christmas-2024-cover.jpg
count: 47
---

Christmas 2024 at the Miami house. Ash got the Nintendo. Kai's first Christmas walking.
```

Legacy can then:
```
memory_search "Christmas 2024 photos"
→ finds 54 Media Index/photos/family/christmas-2024.md
→ returns r2_url for Zeus's browser or full metadata
```

---

## 9. Production Deployment Checklist

### VPS (134.122.8.135 / Tailscale: 100.104.61.29)

```bash
# 1. Clone Codex + set up auto-pull (replaces Syncthing entirely)
git clone https://github.com/Zeus-Delacruz/OneMind-Codex.git /home/onemind/codex
(crontab -l; echo "*/5 * * * * cd /home/onemind/codex && git pull --ff-only >> /home/onemind/codex-sync.log 2>&1") | crontab -

# 2. Deploy openclaw-vps.json
cp openclaw/docs/config/openclaw-vps.json /home/onemind/.openclaw/openclaw.json

# 3. Deploy workspace files
cp -r openclaw/workspace/* /home/onemind/.openclaw/workspace-legacy/
cp openclaw/workspace/SOUL.md /home/onemind/.openclaw/workspace/
cp openclaw/workspace/MEMORY.md /home/onemind/.openclaw/workspace/

# 4. Restart OpenClaw
pm2 restart openclaw-gateway
# Or: openclaw gateway

# 5. Verify Codex is indexed (wait ~60s for initial index build)
# Ask Legacy: "Search my codex for morning protocol"
# Should return results from HP domain within 2 seconds
```

> obsidian-cli is OPTIONAL. Skip unless agents need to create/move notes inside the vault.

### External Services

| Service | Priority | Action |
|---------|----------|--------|
| Convex | High | Apply at convex.dev/claw/apply |
| Cloudflare R2 | Medium | Create bucket, set public access |
| Railway NATS | Active | Already deployed (railway.nats.toml) |
| TimescaleDB | Active | docker-compose.yml |

---

## 10. The Domain Agent Memory Isolation Pattern

Each domain agent (hp, le, ge) has its own workspace AND should only search its domain:

```json
// In their agent configs, scope memorySearch:
{
  "id": "hp",
  "memorySearch": {
    "extraPaths": [
      { "path": "/home/onemind/codex/25-49 HP (Holistic Performance)", "pattern": "**/*.md" }
    ]
  }
}
```

This means:
- `hp` agent only searches HP domain notes → focused, fast, no noise
- `le` agent only searches LE domain → family context only
- `ge` agent only searches GE domain → business context only
- `legacy` (main) searches the entire Codex → full orchestrator context

> Add domain-scoped extraPaths to hp/le/ge agents in `openclaw-vps.json` for a future upgrade.

---

## 11. Sequence: What Happens When Zeus Asks Something

```
Zeus: "Hey Legacy, what time does Ash's soccer practice start this week?"

1. BlueBubbles → OpenClaw gateway (port 18789)
2. Legacy agent session starts
3. BOOT.md loaded → sets context
4. MEMORY.md read → curated Zeus facts
5. memory/2026-02-26.md read → today's context
6. Legacy receives message
7. Invokes memory_search "Ash soccer practice schedule"
   → SQLite hybrid search (BM25+vector) over MEMORY.md + daily logs + /home/onemind/codex/
   → Returns: LE/51 Family/ash-profile.md:34-41 (score: 0.93)
   → Snippet: "Soccer: Tuesday + Thursday 4pm, Saturday 9am. Field: Brickell Park."
8. Legacy replies: "Ash has soccer Tuesday and Thursday at 4pm, and Saturday mornings at 9am at Brickell Park. Want me to put a reminder in your calendar?"
9. (future) convex-sync skill logs this interaction to Mission Control
10. Response sent via BlueBubbles → Zeus's iMessage
```

Total: ~2-4 seconds. Zero API calls to external knowledge stores. Codex is local to VPS.
