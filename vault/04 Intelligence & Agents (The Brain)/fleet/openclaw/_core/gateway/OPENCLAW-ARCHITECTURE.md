# OPENCLAW-ARCHITECTURE.md

> **OneMind OS — OpenClaw Layer: How everything connects.**
> No Agno. DO-native. Two repos. Three data flows.

---

## 1. The Two Repos

```
onemindos-v0.4/          ← CODE (this repo)
├── openclaw/             ← Legacy's runtime files & config
├── nats/                 ← NATS JetStream event bus
├── nexus/                ← cross-runtime comms protocol
├── openmind/             ← physical layer (robots, sensors, HA)
└── docker-compose.yml    ← NATS + openmind-bridge + legacy-bridge

onemindcodex-v0.4/        ← KNOWLEDGE (separate repo)
├── 00-24 SO/             ← Sovereign Operations
├── 25-49 HP/             ← Holistic Performance
├── 50-74 LE/             ← Legacy Evolution
└── 75-99 GE/             ← Generational Empire
```

**Why separate:** The Codex has a completely different lifecycle. You edit it
in Obsidian on Mac. Code changes go through git normally. Mixing them creates
noise in both directions.

**Why monorepo for OS:** `nexus/comms.py` defines NATS subjects used by both
OpenMind bridge and OpenClaw's nexus-bridge skill. They must version together.

---

## 2. The Full Stack (DO-Native, No Agno)

```
┌─────────────────────────────────────────────────────────────────────┐
│  OPENCLAW — Legacy's Runtime                                        │
│  DO App Platform (managed, git-push deploys, Tailscale built-in)   │
│  GRADIENT_API_KEY → all LLM inference bills against DO $5k credit  │
│  Models: Claude Sonnet/Opus 4.6, Qwen3-32B, DeepSeek (serverless)  │
│  Agents: legacy | hp | le | ge                                      │
│  Channels: iMessage (BlueBubbles) | Discord | Voice (ElevenLabs)   │
│  State backup: DO Spaces (auto-synced)                              │
└───────────────────────┬─────────────────────────────────────────────┘
                        │ webhook relay
┌───────────────────────▼─────────────────────────────────────────────┐
│  LEGACY-BRIDGE — NATS ↔ OpenClaw Relay                              │
│  Docker container (4GB Droplet)                                     │
│  Forwards NATS events → OpenClaw webhooks and vice versa            │
└───────────────────────┬─────────────────────────────────────────────┘
                        │ NATS subjects
┌───────────────────────▼─────────────────────────────────────────────┐
│  NATS JETSTREAM — Nervous System                                     │
│  Docker container (same 4GB Droplet as legacy-bridge)               │
│  Subjects: digital.*, physical.*, nexus.*, comms.*                  │
└───────────────────────┬─────────────────────────────────────────────┘
                        │
┌───────────────────────▼─────────────────────────────────────────────┐
│  OPENMIND-BRIDGE — Physical World                                    │
│  Docker container, 1Hz autonomous heartbeat                         │
│  Robots, sensors, cameras, Home Assistant → NATS                    │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│  POSTGRES (TimescaleDB + pgvector) — Managed DO Database            │
│  Sensor timeseries from OpenMind, event logs, embeddings            │
│  NOT used by OpenClaw (OpenClaw has its own SQLite)                 │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 3. Three Data Flows

### Flow A — Codex (Long-term Knowledge)

```
Mac (Obsidian editor)
  ↓ Obsidian Git plugin auto-commits on save
GitHub: Zeus-Delacruz/onemindcodex-v0.4 (private)
  ↓ deploy key (ssh-ed25519 on VPS)
  ↓ cron: */5 * * * * (git pull --ff-only)
/home/onemind/codex/           ← read-only on VPS, never committed from here
  ↓ memorySearch.extraPaths (native OpenClaw config)
Legacy's SQLite index (~/.openclaw/memory/<agentId>.sqlite)
  ↓ hybrid BM25+vector on every memory_search call
Legacy answers Zeus's question with actual Codex facts
```

**VPS is read-only.** Zero conflicts. Obsidian Git is the only writer.

---

### Flow B — Media (OneMind Media Library)

```
Media files (photos, videos, PDFs, transcripts, recordings)
  ↓ upload / sync
DO Spaces bucket: onemind-media
  ↓ data source in Gradient Knowledge Base
DO Gradient Knowledge Base (OpenSearch backend)
  ↓ auto-indexed on schedule (or on push via GitHub Action)
  ↓ hybrid retrieval: kbaas.do-ai.run API
OpenClaw skill (future: gradient-kb skill)
  ↓ Legacy calls skill to query media knowledge base
Legacy answers: "Here's the transcript from that call..." / "Here's the R2 URL for that photo album..."
```

**Why Gradient KB for media (not Codex extraPaths):** Media content is large,
unstructured, and benefits from managed OpenSearch with semantic chunking.
The Codex (`extraPaths`) is for structured vault .md files. These are different
retrieval problems.

**Dropbox (optional):** DO Gradient Knowledge Bases support Dropbox as a data
source. If media is primarily in Dropbox, point the KB data source there
directly — no intermediate Spaces bucket needed.

---

### Flow C — Physical Events (OpenMind → Legacy)

```
Home Assistant / sensors / robots
  ↓ NATS: physical.location.home.state, physical.heartbeat.*
NATS JetStream (Docker)
  ↓ NATS subject forwarding
legacy-bridge (Docker)
  ↓ webhook POST
OpenClaw gateway (App Platform)
  ↓ Legacy processes event
Legacy takes action: sets reminders, sends alerts, updates context
```

---

## 4. OpenClaw on DO App Platform

### Deployment

Based on the official template: `github.com/digitalocean-labs/openclaw-appplatform`

```bash
# One-click deploy via Deploy-to-DO button OR:
doctl apps create --spec .do/app.yaml
```

### Key Environment Variables

| Variable | Value | Purpose |
|---|---|---|
| `GRADIENT_API_KEY` | DO API key with GenAI:read | LLM inference → bills DO $5k credit |
| `ANTHROPIC_API_KEY` | Optional override | Direct Anthropic billing (skip if using Gradient) |
| `OPENCLAW_GATEWAY_TOKEN` | Generated on first deploy | Authenticated access only |
| `OPENCLAW_HOOKS_TOKEN` | For legacy-bridge webhooks | NATS→OpenClaw relay auth |
| `TAILSCALE_AUTH_KEY` | Ephemeral auth key | Private Tailscale access to UI |
| `ZEUS_WHATSAPP_NUMBER` | Zeus's number | BlueBubbles channel |
| `ELEVENLABS_API_KEY` | ElevenLabs key | Legacy's voice |
| `ELEVENLABS_VOICE_ID` | Voice ID | Legacy's voice profile |

### Access Modes

**With Tailscale (recommended):**
- OpenClaw UI accessible at `openclaw.your-tailnet.ts.net`
- No public ports. Private by default.

**Headless:**
- No Tailscale daemon. No web UI.
- Manage via: `doctl apps console <app-id> openclaw`
- Still gets all message channels (WhatsApp, Discord, etc.)

### Updates

Zero-downtime by design:
```bash
git push origin main  # triggers App Platform redeploy automatically
```

---

## 5. LLM Inference via DO Gradient (GRADIENT_API_KEY)

Setting `GRADIENT_API_KEY` routes all OpenClaw LLM calls through DigitalOcean's
Agentic Inference Cloud instead of directly to Anthropic/OpenAI. Billing comes
from your DO $5k credit.

**Available models (as of Feb 2026):**

| Model | ID | Notes |
|---|---|---|
| Claude Sonnet 4.6 | `anthropic-claude-sonnet-4-6` | Default for Legacy |
| Claude Opus 4.6 | `anthropic-claude-opus-4-6` | Heavy reasoning |
| Qwen3-32B | `alibaba-qwen3-32b` | Open-source, serverless only |
| DeepSeek | available | |
| GPT-5.3 Codex | `openai-gpt-5-3-codex` | |
| Llama (Meta) | various | |
| Mistral | various | |

**Per-agent model config in `openclaw.json`:**

```json
{
  "id": "legacy",
  "model": {
    "provider": "gradient",
    "name": "anthropic-claude-sonnet-4-6"
  }
}
```

---

## 6. Codex → VPS: memorySearch Config

The live config in `/home/onemind/.openclaw/openclaw.json`:

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

**CRITICAL format note:** `extraPaths` must be **plain strings**, NOT objects.
`["/path"]` works. `[{"path": "/path", "pattern": "**/*.md"}]` is REJECTED by OpenClaw.

**Domain-scoped agents** (hp/le/ge) have their own extraPaths scoped to their domain:

```json
// hp agent
"extraPaths": ["/home/onemind/codex/25-49 HP (Holistic Performance)"]

// le agent
"extraPaths": ["/home/onemind/codex/50-74 LE (Legacy Evolution)"]

// ge agent
"extraPaths": ["/home/onemind/codex/75-99 GE (Generational Empire)"]

// legacy (main orchestrator) — searches all
"extraPaths": ["/home/onemind/codex"]
```

---

## 7. Codex Access: How Legacy Uses It

### memory_search (90% of queries)

```
Zeus: "What time does Ash get out of school?"
Legacy: memory_search "Ash school schedule"
→ LE/51 Family/ash-profile.md (score: 0.93)
→ "Ash gets out at 3:15pm. Bus home unless you pick up."

Zeus: "What are my protein targets?"
Legacy: memory_search "nutrition macros protein targets"
→ HP/27 Body/nutrition.md (score: 0.91)
→ "Your target: 220g protein, 2800 kcal total."

Zeus: "What's the Circle onboarding flow?"
Legacy: memory_search "Circle onboarding protocol"
→ GE/79 Community/onboarding-protocol.md (score: 0.88)
→ walks through the steps
```

### read_file (when you have the path)

```
read_file "/home/onemind/codex/75-99 GE (Generational Empire)/75 Ventures/community-playbook.md"
```

### obsidian skill (create/move notes — optional)

```
obsidian create "note title"    → creates in 06 Inbox for triage
```

---

## 8. Deploy Key Status

| Item | Status |
|---|---|
| Key generated on VPS | ✅ `/home/onemind/.ssh/github_codex` |
| SSH config written on VPS | ✅ `~/.ssh/config` → `github-codex` host |
| Cron configured | ✅ `*/5 * * * *` git pull |
| Key added to GitHub | ❌ **USER ACTION REQUIRED** |
| Codex cloned | ❌ Blocked on above |

**GitHub URL:** `https://github.com/Zeus-Delacruz/onemindcodex-v0.4/settings/keys`

**Public key to add (with write access checked):**
```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBilpHxjjyPbX3U6XrToqyj0HvrDoxzx6H34h83QEGE8 openclaw-vps-codex-readonly
```

**After key is added, run:**
```bash
ssh onemind-gateway "sudo -u openclaw git clone git@github.com:Zeus-Delacruz/onemindcodex-v0.4.git /home/onemind/codex"
ssh onemind-gateway "sudo -u openclaw git -C /home/onemind/codex config user.name 'Legacy AI' && sudo -u openclaw git -C /home/onemind/codex config user.email 'legacy@onemind.os'"
ssh onemind-gateway "sudo systemctl restart openclaw"
```

---

## 9. DO Gradient Knowledge Base — OneMind Media

For media content (not Codex vault files):

```
Data sources supported:
  ✅ DO Spaces bucket or folder  ← preferred
  ✅ Dropbox folder              ← if media lives in Dropbox
  ✅ Amazon S3
  ✅ Direct file upload
  ✅ Web/sitemap URLs

File formats: .md, .pdf, .txt, .csv, .docx, .json, .html, .epub, .mp3 transcripts...

Embedding models:
  - Qwen3 Embedding 0.6B (8000 token window, multilingual)  ← best for media
  - GTE Large v1.5 (8192 tokens)
  - All-MiniLM-L6-v2 (256 tokens, fast, small)

Auto-index: schedule daily/weekly re-crawl of Spaces bucket
Retrieval: POST https://kbaas.do-ai.run/v1/<kb-uuid>/retrieve
  { "query": "...", "num_results": 5, "alpha": 0.5 }
```

**Future OpenClaw skill: `gradient-kb`**

```javascript
// skills/gradient-kb/index.js
// Queries a Gradient Knowledge Base from inside OpenClaw
// Called by Legacy: "Search media library for Zeus birthday 2025"
```

---

## 10. Rollout Tiers

### Tier 1 — Now (one action away)
- [x] OpenClaw running on Lightsail VPS
- [x] memorySearch.extraPaths configured
- [x] Deploy key generated
- [ ] **Add deploy key to GitHub → clone Codex → Codex live**

### Tier 2 — Migrate to DO (after Codex clone working)
- Snapshot current Lightsail → freeze as backup
- Deploy OpenClaw on DO App Platform (App Platform template)
- Set `GRADIENT_API_KEY` → inference costs from DO $5k credit
- Provision 4GB Droplet for NATS + legacy-bridge + openmind-bridge
- `docker-compose up -d nats legacy-bridge`
- Point `/home/onemind/codex` git remote to same URL (just re-clone)

### Tier 3 — Full DO Stack
- DO Managed Postgres (TimescaleDB for OpenMind sensor data)
- DO Spaces: `onemind-media` bucket + Gradient Knowledge Base
- GitHub Action: auto-sync Codex .md files to Spaces on push (optional redundancy)
- `gradient-kb` skill: Legacy queries media knowledge base

---

## 11. Current VPS State (Feb 26, 2026)

| Item | Status |
|---|---|
| Host | 134.122.8.135 (onemind-gateway alias) |
| Tailscale IP | 100.104.61.29 |
| OpenClaw | ✅ Active (systemd, PIDs confirmed) |
| Config | /home/onemind/.openclaw/openclaw.json (patched) |
| memorySearch | ✅ enabled, extraPaths configured |
| Docker | Installed, zero containers running |
| RAM | 2GB total, ~1.1GB used by OpenClaw, 217MB free |
| Disk | 58GB, 14% used |
| Codex | Not yet cloned (deploy key pending) |
| NATS | Not deployed (needs 4GB plan first) |


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
