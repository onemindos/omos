---
title: "OpenClaw Native Memory Architecture"
type: doctrine
status: active
created: 2026-04-08
owner: zeus
tags: [type/doctrine, domain/so-03, status/active]
ai_summary: "Complete architecture plan for OpenClaw native memory stack: native memory plugin, dreaming (Light/Deep/REM), memory-wiki, MCP servers, and codex integration."
---

# OpenClaw Native Memory Architecture

> Implemented April 8, 2026. This is the canonical reference for how Legacy and all agents access, learn from, and grow their knowledge using ONLY native OpenClaw features.

---

## Design Principles

1. **Native-only** — No third-party memory hacks. Only features built into OpenClaw.
2. **Grows daily** — Dreaming auto-processes every conversation into long-term knowledge.
3. **Structured knowledge** — Memory/Wiki builds evidence-backed claims.
4. **Codex integration** — Agents read/write the full One Mind Codex via MCP.
5. **Robotics-ready** — MCP is the universal protocol for any future machine.

---

## VPS Baseline (Pre-Migration)

| Field | Value |
|-------|-------|
| VPS | DigitalOcean 104.248.4.212 |
| OS | Ubuntu 24.04.3 LTS |
| Specs | 4 vCPU, 16GB RAM, 193G disk (149G free) |
| OpenClaw (before) | 2026.4.5 (3e72c03) |
| OpenClaw (after) | **2026.4.8** |
| Tailscale | onemind-main, 100.77.52.93 |
| Gateway token | `op://aaeyris5z44pdqspk277u2oppi/OpenClaw/gateway-token` |

---

## The Three Native Memory Layers

### Layer 1: Native Memory Plugin (replaces ClawVault)

**Why switch:** ClawVault is a third-party memory slot. The native `memory` plugin is the one that dreaming, memory-wiki, and all future OpenClaw features build on. Dreaming REQUIRES the native memory plugin — it does not work with ClawVault in the memory slot.

**What it does:**
- Fact extraction from every conversation
- Before-prompt recall of relevant context
- Session context injection
- Auto-checkpointing on compaction
- Weekly reflection summaries
- SQLite-backed storage (per-agent)

**Config change:**
```json
"plugins.slots.memory": "memory"
```

Previous: `"clawvault"` → Now: `"memory"`

### Layer 2: Dreaming (Light → Deep → REM)

Three cooperative phases that process daily conversations into long-term knowledge:

| Phase | Schedule | What It Does | Output |
|-------|----------|-------------|--------|
| **Light** | After each session | Scores short-term memories, processes daily note chunks | Scored memory entries |
| **Deep** | Daily | Promotes durable truths to MEMORY.md — replay-safe (reconciles, never duplicates) | `MEMORY.md` |
| **REM** | Weekly | Surfaces "possible lasting truths," contradiction detection | `dreams.md` (Dream Diary) |

**Session corpus ingestion (2026.4.7+):**
- Redacted session transcripts are ingested nightly
- Per-day corpus files: `memory/.dreams/session-corpus/YYYY-MM-DD.txt`
- Cursor checkpointing prevents re-processing
- Checkpoint file: `memory/.dreams/session-ingestion.json`

**Config:**
```json
"dreaming": {
  "enabled": true,
  "frequency": "daily"
}
```

**Optional tuning:**
- `recencyHalfLifeDays` — How fast memories decay in recall scoring
- `maxAgeDays` — Hard cutoff for old memories
- `dreaming.rem.frequency` — Override REM schedule (default: weekly)

**CLI commands:**
```bash
openclaw memory rem-harness    # Preview what REM would surface
openclaw memory promote-explain # Show why a memory was promoted
openclaw memory doctor          # Repair/status check
```

**Files created:**
```
~/.openclaw/workspace/
├── MEMORY.md              ← Promoted durable knowledge (gold)
├── dreams.md              ← Dream diary (REM truths + contradictions)
└── memory/
    └── .dreams/
        ├── session-corpus/
        │   ├── 2026-04-08.txt
        │   ├── 2026-04-09.txt
        │   └── ...
        └── session-ingestion.json  ← Cursor checkpoint
```

### Layer 3: Memory/Wiki (Structured Knowledge Base)

**Requires:** OpenClaw 2026.4.7+

A structured knowledge base that lives alongside conversation memory:

| Feature | What It Means |
|---------|--------------|
| Structured claims + evidence | Each "fact" has explicit evidence linking it back to source |
| Compiled digest retrieval | Summarized knowledge digests agents can retrieve |
| Claim-health linting | Auto-detects stale, contradictory, or weak claims |
| Contradiction clustering | Groups conflicting knowledge for human review |
| Staleness dashboards | Shows what knowledge is aging out |
| Freshness-weighted search | Recent knowledge ranks higher |
| Sync/query/apply tooling | CLI tools for wiki management |

**Config:**
```json
"plugins.allow": [..., "memory-wiki"]
```

**CLI commands:**
```bash
openclaw memory wiki sync     # Sync wiki from conversation memory
openclaw memory wiki query    # Search structured knowledge
openclaw memory wiki apply    # Apply wiki updates
```

---

## Codex Integration

### Git-Synced Codex on VPS

```
Mac (source of truth)
    ↓ git push
GitHub (OneMind-OS/onemind-codex)
    ↓ cron git pull (every 1 min)
/home/onemind/onemind-codex/  ← VPS copy (text only, ~10MB)
```

The VPS codex is text-only (all binaries excluded via .gitignore). Contains 8,500+ markdown files across all 4 domains.

### memorySearch.extraPaths

Points OpenClaw's memory search at the codex:
```json
"memorySearch": {
  "extraPaths": ["/home/onemind/onemind-codex"]
}
```

This makes all 8,500+ .md files searchable via `openclaw memory search`.

### MCP Servers (File + Git Access)

```json
"mcpServers": {
  "codex-fs": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-filesystem", "/home/onemind/onemind-codex"]
  },
  "codex-git": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-git", "--repository", "/home/onemind/onemind-codex"]
  }
}
```

This gives agents:
- **Read/write** any markdown file in the codex
- **Git history** — diffs, blame, commits
- **Version control** — agents can commit changes

---

## plugins.allow (Complete List)

```json
[
  "memory",
  "memory-wiki",
  "lossless-claw",
  "discord",
  "bluebubbles",
  "acpx",
  "openai",
  "amazon-bedrock",
  "anthropic",
  "browser",
  "lobster",
  "diffs",
  "llm-task",
  "open-prose"
]
```

**Changes from previous:**
- REMOVED: `clawvault` (replaced by native `memory`)
- ADDED: `memory` (native memory plugin)
- ADDED: `memory-wiki` (structured knowledge)

### plugins.slots

```json
{
  "memory": "memory",
  "contextEngine": "lossless-claw"
}
```

---

## Agent Architecture (9 Agents)

| Agent | Model | Domain | Channel |
|-------|-------|--------|---------|
| **Legacy** | GPT-5.2 | All | Default orchestrator |
| **Spartan** | GPT-4o | HP 27-28 | Body, fitness |
| **Oracle** | GPT-4o | HP 25-26, 29-34 | Mind, money, meaning |
| **Eden** | GPT-4o | LE | Sustainability |
| **Haven** | GPT-4o | LE | Family, home |
| **Guardian** | GPT-4o | SO | Security, infrastructure |
| **Grid** | GPT-4o | SO | Technology, systems |
| **Trinity** | GPT-4o | Cross-domain | Coordination |
| **Empire** | GPT-5.2 | GE | Business, ventures |

All agents share the native memory system. Each has its own SQLite memory database but shares the same dreaming and wiki infrastructure.

---

## Implementation Sequence

### Step 1: Upgrade OpenClaw
```bash
sudo -u onemind npm install -g openclaw@latest
# Verify: openclaw --version → 2026.4.8
```

### Step 2: Switch memory slot
```bash
openclaw config set "plugins.slots.memory" "memory"
```

### Step 3: Update plugins.allow
```bash
openclaw config set "plugins.allow" '["memory","memory-wiki","lossless-claw","discord","bluebubbles","acpx","openai","amazon-bedrock","anthropic","browser","lobster","diffs","llm-task","open-prose"]'
```

### Step 4: Enable dreaming
```bash
openclaw config set "dreaming.enabled" true
openclaw config set "dreaming.frequency" "daily"
```

### Step 5: Git clone codex on VPS
```bash
rm -rf /home/onemind/onemind-codex
git clone https://github.com/OneMind-OS/onemind-codex.git /home/onemind/onemind-codex
echo "* * * * * cd /home/onemind/onemind-codex && git pull --ff-only 2>&1 | logger -t onemind-pull" | crontab -
```

### Step 6: Set memorySearch.extraPaths
```bash
openclaw config set "memorySearch.extraPaths" '["/home/onemind/onemind-codex"]'
```

### Step 7: Configure MCP servers
```bash
openclaw config set "mcpServers" '{"codex-fs":{"command":"npx","args":["-y","@modelcontextprotocol/server-filesystem","/home/onemind/onemind-codex"]},"codex-git":{"command":"npx","args":["-y","@modelcontextprotocol/server-git","--repository","/home/onemind/onemind-codex"]}}'
```

### Step 8: Restart + verify
```bash
systemctl restart openclaw
openclaw memory status
openclaw memory wiki sync  # Test wiki
```

---

## Why This Architecture Is Bulletproof

| Concern | How It's Handled |
|---------|-----------------|
| **Long-term** | All 3 layers are NATIVE OpenClaw — survive every update |
| **Growing daily** | Dreaming auto-ingests session transcripts → promotes to MEMORY.md |
| **Structured** | Memory/Wiki builds evidence-backed claims from conversations |
| **Codex access** | memorySearch + MCP = agents read/write any .md file |
| **Version control** | Git-tracked codex with cron pull — infinite history |
| **Robotics-ready** | MCP is universal — any future device connects via MCP |
| **No vendor lock** | Everything is markdown + SQLite. No proprietary formats. |

---

## Vault Separation (April 8, 2026)

The local codex vault was separated from binary media:

| What | Where | Size |
|------|-------|------|
| **Text codex** (md, code, configs) | `/Users/zeuslegacy/Documents/onemind-codex/` | ~1-2 GB |
| **Binary media** (pdf, images, docs) | Staging → OM PRIMARY | ~77 GB |
| **Full backup** | `/Volumes/OM PRIMARY/onemind-codex/` | ~79 GB |

Binary extensions moved out:
- Images: jpg, jpeg, png, gif, webp, bmp, tiff, tif, psd, ai, eps, heic, ico, JPG, PNG, HEIC
- Documents: pdf, doc, docx, xls, xlsx, ppt, pptx
- Audio: mp3, wav, flac, aac, m4a
- Video: mp4, mov, avi, mkv, webm
- Archives: zip, rar, 7z, dmg, iso, pkg, gz, tar
- Fonts: ttf, woff, woff2, otf, eot
- Design: sketch, fig
- Other: sqlite, db

Text extensions kept:
- Markdown: md
- Code: ts, js, jsx, tsx, mjs, cjs, py, sh, bash, rb, go, rs, java, c, cpp, h, hpp
- Web: html, htm, php, css, scss, less, liquid, svg
- Data: json, yaml, yml, toml, xml, csv, sql, graphql, prisma
- Config: conf, cfg, ini, env, gitignore, gitattributes, editorconfig
- Text: txt, srt, log, tpl, njk, hbs, ejs, pug

---

*Part of the One Mind framework. See [ONEMIND-CODEX.md](../../../ONEMIND-CODEX.md) for the master document.*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
