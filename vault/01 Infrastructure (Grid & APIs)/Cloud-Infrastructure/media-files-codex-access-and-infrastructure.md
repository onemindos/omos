---
id: media-files-codex-access-and-infrastructure
title: "Media Files, Codex Access & Infrastructure Decisions"
type: note
source: apple-notes
imported: 2026-03-07
status: staged
tags:
  - apple-notes-import
  - infrastructure
  - codex
  - media
  - nats
  - openclaw
  - vps
  - architecture
---

# Media Files, Codex Access & Infrastructure Decisions

## Media Files: The Right Answer

**Don't try to put media on the VPS or in Obsidian.** Here's the clean architecture:

| Layer | Tool | What Lives There | Why |
| --- | --- | --- | --- |
| Storage | Cloudflare R2 (or AWS S3) | All raw media — photos, videos, PDFs, audio | $0.015/GB/month, stable URLs, CDN-fast |
| Index | Codex `54 Media Index/` | Markdown files with frontmatter metadata | Searchable by Legacy, backed up in git |
| Thumbnails | Syncthing to VPS | Compressed previews only | Legacy can visually reference without full file |
| Access | R2 public URL in frontmatter | Legacy fetches when needed | One URL, always accessible, no disk bloat |

Each media item in the Codex index looks like:

```yaml
---
title: Zeus Family Reunion 2025
type: photo-album
url: https://r2.yourdomain.com/media/family/reunion-2025/
local_path: ~/Dropbox/OneMind-Media/family/reunion-2025/
tags: [family, le, photos, 2025]
people: [Zeus, kids, extended-family]
domain: LE
---
Summer family reunion at the farm. 47 photos.
```

---

## The Correct Architecture for Codex Access

```
            MULTI-DEVICE CODEX SYNC

Mac (Obsidian + edit)
iPhone (Obsidian app)  ── Syncthing ──→ VPS /home/openclaw/codex/
Dell (read + backup) ─┘        │
                          GitHub
                     (version backup)

                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
obsidian-cli search   read_file (native)   OpenClaw vector
(full-text + names)   (any .md file)       index (BM25 + semantic)
        │                   │                   │
        └───────────────────┴───────────────────┘
     Legacy answers from Codex directly — NO API round-trip
```

### Why Syncthing over GitHub API

- GitHub API = 1MB/min rate limit, no real-time, requires auth tokens
- Syncthing = real-time sync, LAN-speed (Mac↔VPS over Tailscale hits ~50MB/s), zero API calls, works offline
- GitHub stays as the version control backup, not the sync layer

### The Three Access Tiers for Legacy

| Tier | Tool | What It Does | When Used |
| --- | --- | --- | --- |
| Native read | `read_file` (built-in) | Opens any `.md` in Codex by path | "Read my morning protocol" — exact file |
| Search | `obsidian-cli search-content` (skill loaded) | Full-text search across all 400+ files | "Find anything about nutrition" |
| Semantic | OpenClaw's vector store (extend it) | Nearest-neighbor match on query | "What have I decided about the farm?" |

The only missing piece: **tell Legacy where the Codex lives on disk.** One line in `AGENTS.md`:

```markdown
## Codex Location
The OneMind Codex vault lives at: /home/openclaw/codex/
Use obsidian-cli with this vault or read_file for direct access.
Domain map: 00-24 SO | 25-49 HP | 50-74 LE | 75-99 GE
```

### Codex Sync Paths

```
Codex (GitHub) → git pull cron every 5min → /home/openclaw/codex
→ memorySearch.extraPaths (native, already live in config)
```

```
OneMind Codex (vault .md files)
  → git pull cron → VPS extraPaths → Legacy's native memorySearch

OneMind Media (content library)
  → DO Spaces bucket → Gradient Knowledge Base (auto-indexed)
  → Queryable via kbaas.do-ai.run API or ADK agents
```

---

## GitHub Sync (Changes Flow Both Ways)

### Flow 1: Zeus edits on Mac → pushes to GitHub → VPS syncs

```
Obsidian (Mac) → git commit (Obsidian Git plugin)
  → GitHub push → VPS cron (git pull every 5 min)
  → Syncthing sync (immediate if on same Tailscale network)
```

### Flow 2: Legacy writes to Codex → shows up on Mac

```
Legacy write_file /home/openclaw/codex/06 Inbox/note.md
  → Syncthing pushes to Mac immediately
  → Obsidian shows new note in vault
  → Git plugin auto-commits on schedule
  → GitHub updated
```

This means Legacy can **write back to your Codex** — save research, create family notes, log events, build out LE entries for your kids — and you see it instantly in Obsidian on your Mac.

---

## The Right Stack

```
Mac (Obsidian) ←→ GitHub ←→ VPS (git clone, write access)
                                ↑
                        obsidian-cli (npm)
                        gives Legacy vault-aware writes:
                        proper frontmatter, links, note structure
```

**obsidian-cli on VPS** is the right call. It's not Obsidian the app (no GUI needed), just a CLI tool that understands vault structure. Already packaged as the `obsidian-1.0.0` OpenClaw skill you have installed.

### How It Works End-to-End

**Legacy reads:**

```
memory_search "anything" → finds files instantly via search index
read_file /home/openclaw/codex/path/to/note.md → reads full file
```

**Legacy writes (new note):**

```
obsidian-cli create "Note Title" --vault /home/openclaw/codex --content "..."
→ creates properly formatted .md in 06 Inbox/
→ Legacy runs: git add . && git commit -m "..." && git push
→ GitHub updated
→ Your Mac Obsidian pulls it in next sync cycle (1-5 min)
→ You see it in Obsidian, triage it to the right folder
```

**You write:**

```
Edit anything in Obsidian on Mac
→ Obsidian Git auto-pushes to GitHub
→ VPS auto-pulls (cron, 5 min)
→ Legacy sees your new content
```

---

## What Needs to Change

Right now the deploy key is **read-only**. We need to upgrade it to write access, plus set up `obsidian-cli` on the VPS.

**Step 1 — You do this (30 seconds):**

Go to: **github.com/Zeus-Delacruz/onemindcodex-v0.4/settings/keys**

Add the key with "Allow write access" checked this time:

```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBilpHxjjyPbX3U6XrToqyj0HvrDoxzx6H34h83QEGE8 openclaw-vps-codex-
```

(Rename it `openclaw-vps-readwrite` if you want to be accurate.)

**Steps 2-5 — I run once you confirm:**

- Clone repo to VPS with write access
- Configure git user on VPS (`Legacy AI <legacy@onemind.os>`)
- Install `obsidian-cli` on VPS
- Set vault default to `/home/openclaw/codex`
- Update `AGENTS.md` so Legacy knows the write protocol (pull → write → commit → push)
- Test: tell Legacy to create a note and confirm it appears in your Obsidian

---

## Revised Stack: What's Actually Hosted

```
        DO Droplet (OpenClaw Hub)

  • openclaw gateway (systemd, --force)
  • 4 agents now → 10-15 agents over time
  • Codex on disk (/home/openclaw/codex)
  • codex-autopush.service (inotifywait)
  • codex cron pull (*/1 * * * *)
  • Discord / BlueBubbles / Telegram channels
  • hooks endpoint (webhook receiver)
  • sherpa-onnx-tts (local voice, CPU)

  NOTHING ELSE — no NATS container, no n8n,
  no TimescaleDB, no Letta, no OM1 runtime

              │ API calls
      ┌───────┼───────┐
      ▼       ▼       ▼
  Synadia   OM1      Composio
  Cloud     Cloud    (future,
  (managed  (perception enterprise
  NATS,     APIs,    tools)
  free tier) $2k
             credits)
```

### Recommendation: General Purpose — 4 vCPU / 16GB RAM / NVMe SSD

| Spec | Value | Why |
| --- | --- | --- |
| vCPU | 4 dedicated | Multi-agent concurrent spawning, no noisy neighbor |
| RAM | 16GB | 10-15 agents + Codex index + OS headroom |
| Disk | 160GB NVMe | Fast git ops, Codex, concurrent writes |
| Price | ~$96/mo | Leaves half your DO credits for volumes, snapshots |
| Region | NYC3 or SFO3 | Lowest latency to you |

**Optional:** +100GB Block Volume ($10/mo) for Codex storage, independently expandable.

---

## What EXISTS as Code but is NOT Deployed

### `onemindos-v0.4`:

| Code | Location | State |
| --- | --- | --- |
| NATS server config | hub.conf, leaf.conf | Config written, server not running anywhere |
| Nexus layer | registry.py, comms.py, awareness.py | Python code, not deployed |
| OpenMind bridge | bridge.py | OM1-pattern digital cortex loop, not running |
| Backend (Agno AgentOS) | docker-compose.yml | In compose, never started |
| Flutter frontend | docker-compose.yml | In compose, never built |
| TimescaleDB init | init.sql | Schema written, DB not running |
| NATS event bus | bus.py | Bus implementation, unused |
| Event flow system | EVENT_FLOW.md + nats_bridge.py | Fully designed, not deployed |

### `onemindcodex-v0.4`:

- Complete Obsidian vault with 4-domain structure — fully populated and working as knowledge base
- Letta plugin (`.letta`, `plugins`) — installed but superseded by OpenClaw

### Architecture Drift Problem (Critical)

Your docs have multiple conflicting layers written at different times.

---

## Does NATS Still Matter? Yes — But Only for One Thing

**The NATS purpose narrows to exactly this:**

```
Solar-punk IoT devices (MQTT)
├── Home Assistant sensors
├── Victron Cerbo GX (solar/battery)
├── Environmental sensors
├── Future Meshtastic LoRa nodes
└── Any device that speaks MQTT natively

        ↓ NATS MQTT port 1883

NATS JetStream (hub)
  • Fan-out to subscribers
  • 30-day event buffer
  • Multi-location leaf nodes (house ↔ office ↔ van)

        ↓ NATS subject subscription

OpenClaw skill (nats-listener)
  • Subscribes to events.>
  • Injects context into Legacy's awareness
```

This is irreplaceable with Redis. This is also the solar-punk community use case — your homestead/office/off-grid van all sharing a single nervous system via NATS leaf nodes with MQTT on the edge.

### NATS → OpenClaw: How It Already Works

You already have it in docker-compose — it's the `legacy-bridge` service. Here's the actual flow:

```
IoT device (MQTT)  ───→ NATS port 1883
Home Assistant     ───→ NATS (native MQTT)
Victron Cerbo GX   ───→ NATS (native MQTT)

    NATS JetStream (subscribes to events.>)
              │
    legacy-bridge service
    (backend/legacy/bridge.py)
              │ POST webhook
    OpenClaw gateway :18789
              │
          Legacy
    reads the event, decides action
    (turn off lights, alert Zeus, etc.)
```

The `bridge.py` is the glue — it subscribes to NATS subjects and forwards them as webhook events to OpenClaw. Legacy sees them as incoming hooks and reacts via skills. That's the complete loop with **zero custom code needed**.

---

## Codex Viewer Options

> Can we view the Codex in the Control UI?

**Option 1: Read-only Markdown Viewer Tab** — Add a tab to the Control UI that shows a file tree on the left, clicking a file renders the markdown on the right. Would need custom UI code added to the OpenClaw frontend bundle.

**Option 2: Embedded Obsidian Publish / Static Site** — Generate a static website from the Codex using Obsidian Publish, mkdocs, or Quartz → serve it on a subdomain → embed it as an iframe tab in the Control UI.

**Option 3: File Browser via OpenClaw Canvas** — OpenClaw has a "Canvas" feature for custom HTML panels. You could create a canvas that lists files and renders markdown using JavaScript — the quickest option since it doesn't require modifying the Control UI source.

**Option 4: VS Code Server / Code-Server** — Run code-server on the VPS pointed at `/home/openclaw/codex` → full VS Code in browser with markdown preview. Add it as a link/iframe in the Control UI.

Canvas option (#3) is fastest to implement. Option #4 gives the most power (full editing + preview + git integration).

---

## DO Knowledge Base vs Codex

> DO Knowledge Base → Media agent only

Codex stays git + extraPaths. That's native, fast, already working.

DO Knowledge (via Gradient RAG) is for: uploading video, PDFs, images from Dropbox/DO Spaces → searchable knowledge base via Gradient API.

**Use case:** OneMind media agent that researches your content library, past recordings, documents for community posts, scripts, briefings. Not for Codex. Add this when you're actively producing content.

---

## Obsidian CLI → obsidian-local-rest-api on Mac

Install the `obsidian-local-rest-api` plugin in Obsidian on Mac → it exposes a local REST API on port 27123. Legacy can call it via a webhook skill to:

- Create notes
- Search vault
- Open specific files
- Trigger Obsidian commands

The Codex git sync already covers 95% of this (Legacy can write files directly to `/home/openclaw/codex` → autopush service commits them). The REST API adds the UI-trigger layer (open a note, run a plugin). Low priority — add after migration settled.

---

## Why n8n Is Not Needed

For your specific workflows (webhooks, scheduling, calling APIs, agent routing), OpenClaw has every primitive:

| n8n Concept | OpenClaw Equivalent |
| --- | --- |
| Trigger nodes | `hooks` block → webhook receiver |
| Schedule trigger | `proactive-agent` cron → scheduled triggers |
| Action nodes | Skills (JS modules) → HTTP calls, API integrations |
| Router nodes | Agent-to-agent routing via `openclaw.json` |
| Event chain | OpenClaw hooks emit → other services subscribe |

The only thing n8n adds is a visual UI for non-technical users. You don't need that. **Drop n8n.**
