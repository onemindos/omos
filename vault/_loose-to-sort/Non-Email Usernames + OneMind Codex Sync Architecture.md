# Non-Email Usernames + OneMind Codex Sync Architecture

---

## Non-Email Usernames
*(not email addresses — just username strings)*

`zeusdelacruz` · `creditexecutive` · `creditexecutive1` · `legacywarrior` · `legacywarrior1` · `legacywarriorx` · `zeuslegacy` · `zeuslegacydelacruz` · `zeusd` · `wealthwarrior` · `admin` · `adelacruz418` · `legacy_ai` · `OneMind` · `OneMindOS` · `ANTDELACRUZ` · `ZEUSDELACRUZ` · `Creditexecutive18`

---

## Other People's Credentials (Stored in Vault)

| Email | Count |
|-------|-------|
| Ajhailey@theclassof700.com | 2 |
| ahailey | 2 |
| Jeremylewis357@gmail.com | 1 |
| Josh@getmormedia.com | 1 |
| eflores1524@gmail.com | 1 |
| emkunderwood.ig@gmail.com | 1 |
| Luzmdelacruz@icloud.com | 1 |

---

## 1Password Vault Audit (263 items, all vaults)

### Active System Emails

| Email | Count | Domain | Purpose |
|-------|-------|--------|---------|
| zeus.ops@legacyhq.ai | 19 | legacyhq.ai | Business / AI ops |
| ui@zeusdelacruz.me | 19 | zeusdelacruz.me | SO/tech alias |
| finance@zeusdelacruz.me | 18 | zeusdelacruz.me | Personal finance |
| shopping@zeusdelacruz.me | 15 | zeusdelacruz.me | Retail |
| zeus.vaults@pm.me | 12 | ProtonMail | Root of trust / 1P login |
| zeus@zeusdelacruz.com | 8 | zeusdelacruz.com | Public brand |
| partnernetwork@legacyhq.ai | 7 | legacyhq.ai | Affiliates |
| performance@zeusdelacruz.me | 6 | zeusdelacruz.me | Fitness/wellness |
| finance@legacyhq.ai | 6 | legacyhq.ai | Business finance |
| travel@zeusdelacruz.me | 5 | zeusdelacruz.me | Travel |
| social@zeusdelacruz.me | 4 | zeusdelacruz.me | Brand/social |
| legal@zeusdelacruz.me | 4 | zeusdelacruz.me | Legal |
| education@zeusdelacruz.me | 4 | zeusdelacruz.me | Courses |
| recreation@zeusdelacruz.me | 2 | zeusdelacruz.me | Gaming |
| legacy@zeusdelacruz.me | 1 | zeusdelacruz.me | LE |
| le@zeusdelacruz.me | 1 | zeusdelacruz.me | LE |
| legal@legacyhq.ai | 1 | legacyhq.ai | Business legal |
| legacywarriorx@proton.me | 1 | ProtonMail | Secondary secure |

### ⚠️ Problem Emails (Need Action)

| Email | Count | Issue |
|-------|-------|-------|
| me@zeusdelacruz.com | 18 | NOT a registered Zoho alias — biggest gap |
| executivebycredit@gmail.com | 6 | Old Gmail — should migrate |
| legacyplatformholdings@gmail.com | 4 | Old Gmail — should migrate |
| Finance@zeusdelacruz.me | 2 | Capitalization bug (same as finance@) |
| Shopping@zeusdelacruz.me | 2 | Capitalization bug (same as shopping@) |
| zeus@onemindcodex.com | 3 | Domain may retire |

---

## OneMind Codex — Sync Architecture

### USB Drive Feasibility

You CAN do it now for the markdown + code portion. You can't today for the full vault due to size:

| What | Size | USB Feasible? |
|------|------|---------------|
| Markdown + code (git-tracked) | ~1 GB | **Yes, today** |
| `.git` directory | 46.2 GB (bloated) | No — needs cleanup first |
| `_docs` (PDFs) | 46.7 GB | Needs 256GB+ drive |
| `_images` | 29.4 GB | Same |
| **Total vault** | **125.7 GB** | Needs 256GB USB-C SSD ($50–75) |

The blocker is `.git` at **46.2 GB**. After:
1. Clean `.git` with `git filter-repo` → shrinks to ~1 GB
2. Fix `.gitignore` to exclude binaries
3. Convert most-used PDFs to markdown over time

...vault drops to ~80 GB. A **256GB Samsung T7 SSD ($75)** carries everything.

Git works from USB — plug into any Mac/Linux:
```bash
cd /Volumes/USB/onemind-codex
git status          # works
git commit -am "..."  # works
git push            # works (if internet available)
obsidian .          # open as vault (if Obsidian installed)
```

---

### Obsidian Plugin → OpenClaw Chatbot (Offgrid)

**Concept: `obsidian-openclaw` Plugin**

```
┌─────────────────────────────┐
│  Obsidian (Mac or Mobile)   │
│                             │
│  OpenClaw Chat Panel        │
│  [You]: How do I...?        │
│  [Legacy]: Based on your    │
│  codex note in HP/27...     │
└─────────────────────────────┘
            │
     HTTP/WebSocket
            ↓
  OpenClaw Gateway (VPS or local)
```

- **Online mode:** Chat panel → OpenClaw gateway on VPS via WebSocket (Tailscale `100.77.52.93`)
- **Offgrid mode Option 1:** Queue messages locally → sync when back online (stores in `.openclaw-queue.json`)
- **Offgrid mode Option 2:** Local Ollama fallback — detects no VPS connectivity, routes to local Ollama with vault as context

---

### Read-Only Phone — Conflict Risk = Zero

Since phone is read-only:
- **No conflict files ever** — conflicts only happen when two devices write simultaneously
- Obsidian Sync Standard ($4/mo) is perfect — browse, search, read
- If edits needed: prompt OpenClaw via Telegram → Legacy edits on VPS → syncs everywhere

> **Workflow: Phone reads. OpenClaw writes. Mac does both.**

---

### Obsidian CLI Tools

**Obsidian CLI (Desktop)**
```bash
obsidian daily                           # Open today's daily note
obsidian search query="credit"           # Search your vault
obsidian daily:append content="- [ ] Review docs"  # Add task
obsidian read                            # Read current file
obsidian tags counts                     # View all tags
obsidian create name="Note" template=Travel  # From template
```
Requires Obsidian desktop running. Enable in `Settings → General → CLI`.

**Obsidian Headless Sync (VPS/Server — THIS IS HUGE)**
```bash
npm install -g obsidian-headless

ob login
ob sync-setup --vault "OneMind Codex"
ob sync --continuous           # Real-time sync, no GUI needed
ob sync-config --mode pull-only  # VPS only receives, never sends
```

Instead of `codex-pull.timer` (git pull every 1 min), run:
```bash
ob sync --continuous --path /home/onemind/onemind-os/codex/
```
This gives the VPS **real-time vault access including binaries** — synced in seconds instead of minutes.

---

### Repo + Vault + Filesystem in One

The codex IS all three simultaneously:

| Layer | How It Works Today |
|-------|-------------------|
| **Filesystem** | Plain folders + files on disk. Any OS can read it. |
| **Obsidian Vault** | Open the folder in Obsidian → wiki-links, graph view, search. `.obsidian` stores config. |
| **Git Repo** | `.git` tracks changes. GitHub is the remote. VPS pulls it. |

The three layers are **additive, not conflicting**:
- `.obsidian` is ignored by git (in `.gitignore`)
- `.git` is ignored by Obsidian Sync and Dropbox
- The actual files are shared by all three

```
/onemind-codex/
├── .git/               ← Git sees this (version control)
├── .obsidian/          ← Obsidian sees this (vault config)
├── 00–24 SO/           ← ALL THREE see this (content)
├── 25–49 HP/           ← ALL THREE see this
├── _codex/             ← ALL THREE see this
├── _docs/              ← Filesystem + Obsidian (git ignores binaries)
└── _images/            ← Filesystem + Obsidian (git ignores binaries)
```

> **Key insight: Git tracks the text. Obsidian Sync / Dropbox handles the binaries. The filesystem is the universal truth.**

---

### Sync Method Coverage

| Sync Method | What Gets There | Coverage |
|-------------|----------------|----------|
| Git | Markdown + code only | ~1 GB (the brain) |
| Obsidian Sync | MD + selected PDFs/images | Up to 100 GB (configurable) |
| USB-C SSD | Everything | Full vault (~80 GB after cleanup) |
| Dropbox (optional) | Full folder mirror | 2 TB ceiling + API access |

Each device gets exactly the subset it needs. The Mac is the only machine that always has everything.

---

### Do We Still Need Dropbox?

**Without Dropbox (Obsidian Sync Only + Git)**
```
Mac (Obsidian Desktop)
  ⇕ Obsidian Sync (real-time, E2E encrypted)
VPS (Obsidian Headless — ob sync --continuous)
  ⇕ Obsidian Sync
iPhone (Obsidian Mobile — read-only)

Mac --git push--> GitHub --git pull--> VPS (text-only backup path)
```
Cost: $8–16/mo (Obsidian Sync Plus)

**With Dropbox (Belt + Suspenders)**
Cost: $20–28/mo (Obsidian Sync + Dropbox Plus)

**What Dropbox Adds That Obsidian Sync Doesn't:**

| Capability | Obsidian Sync | Dropbox |
|-----------|--------------|---------|
| Sync vault to VPS | ✅ (headless) | ✅ (dbxcli) |
| Sync to phone | ✅ (mobile app) | ❌ |
| Sync to Mac | ✅ | ✅ |
| Access from any browser | ❌ | ✅ (dropbox.com) |
| Access from non-Obsidian devices | ❌ | ✅ |
| API for agents | ❌ | ✅ (full REST API, longpoll) |
| Storage ceiling | 100 GB max | 2 TB |
| File versioning / restore | 60-day history | 180-day history |
| Sharing files with others | ❌ | ✅ |

**Bottom line:** If you only need Mac ↔ VPS ↔ Phone, **Obsidian Sync alone is sufficient.** Add Dropbox if you need browser access, agent API, storage >100 GB, or file sharing.

---

### Do We Still Need RAG?

**Yes — RAG and sync solve completely different problems:**

```
SYNC  = "Get the file to the right machine"
RAG   = "Find the right answer inside 16,000+ files"
```

Obsidian Sync **unlocks** better RAG — it doesn't replace it:

| Before (git only on VPS) | After (Obsidian Headless on VPS) |
|--------------------------|----------------------------------|
| VPS has ~3,601 markdown files | VPS has **ALL files** — markdown, PDFs, images |
| OpenClaw searches markdown only | OpenClaw can now **process PDFs/binaries** |
| RAG: markdown → embed → pgvector | RAG: **PDF → extract text → embed → pgvector** |

Flow:
```
Obsidian Sync delivers files to VPS
            ↓
PDF extraction pipeline (pdfplumber/unstructured)
            ↓
Text chunks → OpenAI embeddings → pgvector
            ↓
OpenClaw memory_search queries pgvector
            ↓
Agent answers with PDF knowledge
```

---

### Dropbox as Backup Layer

**Backup Depth Chart:**

| Layer | What | Recovery Scenario |
|-------|------|------------------|
| **1: Obsidian Sync** | 60-day version history, E2E encrypted | "I need yesterday's version of this note" |
| **2: Git + GitHub** | Full text history, infinite (markdown) | "I need the version from 6 months ago" |
| **3: Dropbox** | 180-day history, ALL file types including binaries | "I accidentally deleted a PDF" |
| **4: USB-C SSD** | Physical snapshot, offline | "Everything online is gone" (nuclear) |

**Key value:** Dropbox is a **separate provider** from Obsidian — independent failure domains. If one goes down, your data is still in the other.

**Revised Architecture:**
```
Mac (source of truth)
├── Obsidian Sync ──→ VPS (headless, real-time)
│                    └─→ Phone (selective, read/light-edit)
├── Dropbox ──→ Dropbox Cloud (full mirror, 180-day history, API)
└── Git ──→ GitHub (text version control, infinite history)
```

---

### Updated Recommendation

| Layer | Tool | Cost |
|-------|------|------|
| **Primary Sync** | Obsidian Sync Plus | $8–16/mo — Mac ↔ VPS ↔ Phone (real-time) |
| **Version Control** | Git + GitHub | Free — markdown + code, infinite history |
| **Backup** | Dropbox Plus *or* iCloud | $12/mo *or* $3/mo — full vault mirror, independent provider |
| **Offline** | USB-C SSD (256GB Samsung T7) | $75 one-time — physical snapshot, nuclear recovery |

**Total: $20–28/mo** with Dropbox, or **$11–19/mo** with iCloud as backup instead.

---

### How Git Fits In

Git is **not part of the real-time sync loop** — it's a separate layer for version control and deployment.

**Two Sync Planes:**

| | Plane 1: Real-Time Sync (Obsidian Sync) | Plane 2: Version Control (Git) |
|-|----------------------------------------|-------------------------------|
| **Devices** | Mac ↔ Obsidian Cloud ↔ VPS ↔ Phone | Mac → GitHub ← VPS |
| **Timing** | Automatic, seconds | On commit — manual or scripted |
| **Files** | ALL files (md, pdf, images) | TEXT files only (md, code, config) |
| **Purpose** | Every device has the latest | History, rollback, collaboration |

**When an Agent Edits a File:**
```
VPS: Legacy edits workout-plan.md
├── INSTANT: File saved to VPS disk
├── ~1-5 sec: Obsidian Headless Sync detects change
│           → pushes to Obsidian Cloud
│           → Phone sees it
│           → Mac sees it (if awake)
└── LATER: Git commit happens (see options below)
```

**Three Ways Git Commits Happen:**

**Option A: Automated Timer**
```bash
# codex-push.timer (runs every 5 min on VPS)
cd /home/onemind/onemind-os/codex
git add -A
git diff --cached --quiet || git commit -m "auto: agent updates $(date +%Y-%m-%d-%H%M)"
git push origin main
```

**Option B: Event-Driven (Better)**
```bash
# In OpenClaw's post-edit hook
after_file_edit() {
  cd /home/onemind/onemind-os/codex
  git add "$EDITED_FILE"
  git commit -m "agent($AGENT_NAME): update $EDITED_FILE"
  git push origin main
}
```

**What Git Gives You That Obsidian Sync Doesn't:**

| Capability | Obsidian Sync | Git |
|-----------|--------------|-----|
| Real-time sync | ✅ seconds | ❌ on commit only |
| Binary files (PDFs, images) | ✅ | ❌ |
| Full history forever | ❌ 60 days max | ✅ infinite |
| Diff / blame | ❌ | ✅ |
| Branching | ❌ | ✅ |
| Rollback to any point | ❌ limited | ✅ any commit |
| CI/CD triggers | ❌ | ✅ GitHub Actions |
| Public site deploy | ❌ | ✅ push → Cloudflare Pages |

---

### codex-watch.service — Repair Plan

`codex-watch.service` already exists on VPS — disabled since 2026-03-17 due to detached HEAD errors.

**Current status:**

| Service | Status |
|---------|--------|
| codex-watch.service | Disabled (paused 2026-03-17 — detached HEAD errors) |
| codex-pull.timer | Active (pulls GitHub → VPS every 1 min) |

**What `codex-watch.sh` already does:**
```bash
# Already on VPS at /home/onemind/codex-watch.sh
inotifywait -m -r \
  -e close_write -e create -e delete -e move \
  --exclude '\.git' \
  "$REPO/codex/" "$REPO/openclaw/workspaces/" |
while read -r dir event file; do
  sleep 3                    # 3-second debounce
  git add codex/ openclaw/workspaces/
  git commit -m "vps: auto-sync $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  git push
done
```

**Fix 1: Detached HEAD Guard**
```bash
# Add to top of the commit block in codex-watch.sh
BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null)
if [[ -z "$BRANCH" ]]; then
  echo "[$(date)] ERROR: Detached HEAD — running git checkout main"
  git checkout main
fi
```

**Fix 2: Better Commit Messages**
```bash
# Replace the commit line with:
CHANGED=$(git diff --cached --name-only | head -5 | tr '\n' ', ')
git commit -m "agent: update ${CHANGED%,}" -q
```

Produces commits like:
```
agent: update codex/25-49 HP/27 Body/workout-plan.md
agent: update codex/00-24 SO/24 Inbox/meeting-notes.md, codex/00-24 SO/03 Agents/memory.json
```

**With Both Running on VPS:**
```
Agent edits file on VPS disk
├── inotifywait detects change (instant)
│   ├── Obsidian Headless Sync → Obsidian Cloud → Phone + Mac (1-5 sec)
│   └── codex-watch.sh → git commit + push → GitHub (3-5 sec)
└── Both fire from the SAME filesystem event — no polling, no timers
```

**New VPS Service Stack:**

| Service | Replaces | Purpose |
|---------|----------|---------|
| `ob sync --continuous` | `codex-pull.timer` | Real-time bidirectional vault sync (all files) |
| `codex-watch.service` (re-enabled) | — | Event-driven git commit + push (text files) |
| `openclaw.service` | — | Agent gateway (unchanged) |
| `openclaw-reload.service` | — | Auto-restart on config changes (unchanged) |

`codex-pull.timer` becomes **redundant** once Obsidian Headless Sync is running — disable it.

---

### Option A: Dropbox + Obsidian Sync (Most Complete)

```
┌────────────────────────────────────┐
│  Mac (Dropbox folder)              │
│  ~/Dropbox/onemind-codex/          │
│  Obsidian + VS Code + Git          │
│  Obsidian Sync enabled (md only)   │
│  .git excluded from Dropbox via xattr │
└────────────────────────────────────┘
         │              │              │
   Dropbox sync   Obsidian Sync    git push
         ↓              ↓              ↓
  ┌──────────┐  ┌──────────────┐  ┌──────────┐
  │ Dropbox  │  │  Obsidian    │  │  GitHub  │
  │  Cloud   │  │   Cloud      │  │ (md+code)│
  │  (2TB)   │  │  (md only)   │  └──────────┘
  └──────────┘  └──────────────┘
         │              │           git pull
     API/RAG      Obsidian Sync         ↓
         ↓              ↓         ┌──────────┐
  ┌──────────┐  ┌──────────────┐  │   VPS    │
  │  RAG/    │  │    iPhone    │  │  Agents  │
  │ DO Know. │  │   Obsidian   │  │   (md)   │
  │  Agents  │  │ (read+edit)  │  └──────────┘
  └──────────┘  └──────────────┘
```

| Item | Details |
|------|---------|
| **Cost** | $12/mo (Dropbox Plus) + $4/mo (Obsidian Sync Standard) = **$16/mo** |
| **Mobile** | Full read + edit via Obsidian Sync. Toggle what syncs (md only = fast, or include images/PDFs up to 1 GB) |
| **Risk** | Obsidian warns against dual-sync. Low for single user — possible conflict files occasionally |
| **RAG** | Dropbox API connectors work |
| **VPS** | git pull for md. Dropbox API for binaries on-demand |

---

### Dropbox CLI Reference

Two tools available:

**dbxcli (Dropbox's official CLI — Go binary)**
```bash
# Install on Mac
brew install dbxcli
# Or download from github.com/dropbox/dbxcli

# Commands
dbxcli ls /onemind-codex/          # List files
dbxcli get /onemind-codex/file.md  # Download file
dbxcli put local.md /onemind-codex/ # Upload file
dbxcli search "credit repair"       # Search filenames
dbxcli du                           # Show space usage
```

**Dropbox HTTP API v2** (api.dropboxapi.com)

| Endpoint | Use |
|----------|-----|
| `files/list_folder` | Browse directory contents |
| `files/download` | Download a file |
| `files/search` | Search by filename |
| `files/upload` | Upload <150 MB files |
| `files/list_folder/longpoll` | Watch for changes (event-driven) |

> The `longpoll` endpoint is gold — OpenClaw agent could subscribe to Dropbox change events and react to new files.

---

### PDF / Binary Indexing Options

**Option A: OpenClaw Native Search (Already Working)**
- Hybrid BM25 + vector search over the codex
- `memorySearch.extraPaths` → `/home/onemind/onemind-os/codex` — already indexing vault
- OpenAI embeddings, MMR diversity reranking, temporal decay
- ⚠️ Only searches **git-tracked files** (~3,601 markdown files) — does NOT search PDFs/binaries

**Option B: pgvector + PDF Processing (Best for your stack)**
Pipeline:
1. `dbxcli` or Dropbox API pulls PDFs from Dropbox to VPS
2. Python script extracts text (`PyPDF2`, `pdfplumber`, or `unstructured`)
3. Embed with OpenAI or local model (Nomic Embed)
4. Store chunks + vectors in existing pgvector database
5. OpenClaw agents query via a custom skill

**Option C: ChromaDB / Milvus Lite (Simpler)**
- MemSearch (Zilliz) — ONNX local embeddings, no API key, BM25+vector
- Mount Dropbox folder → MemSearch watches it → auto-indexes

**File Type Support:**

| Type | Extraction Tool | Quality |
|------|----------------|---------|
| PDF | `pdfplumber`, `PyMuPDF`, `unstructured` | Good for text PDFs; OCR needed for scanned |
| DOCX/XLSX/PPTX | `python-docx`, `openpyxl`, `python-pptx` | Excellent |
| Markdown | Direct read | Perfect |
| Images | OCR (tesseract) or VLM (Gemini, GPT-4o) | Varies |
| HTML | `beautifulsoup4`, `defuddle` | Good |
