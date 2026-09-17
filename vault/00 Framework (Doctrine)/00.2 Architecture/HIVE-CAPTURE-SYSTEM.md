---
id: arch-00.2-20260304-capture
title: "Hive Capture System"
type: architecture
status: active
created: 2026-03-04
modified: 2026-03-04
owner: zeus
tags: [type/architecture, status/active, layer/capture, system/hive]
ai_summary: "Complete architecture of the OneMind Hive Capture System — 14 input sources, the 5-step HARVEST process, agent roles, and competitive analysis vs GTD/PARA/CODE/Readwise."
---

# Hive Capture System

> **"Other frameworks make YOU the inbox. OneMind makes the SYSTEM the inbox."**

The Hive Capture System is the ingestion layer of OneMind — the mechanism by which information from the world flows into the Codex and becomes intelligence, tasks, or memory. It is the difference between a knowledge system you maintain and a knowledge system that maintains itself.

---

## The Core Problem With Every Other Capture System

Every productivity framework ever built has the same hidden assumption:

> **You are the capture mechanism.**

| Framework | Capture Model | What You Must Do |
|-----------|--------------|-----------------|
| **GTD** | "Mind like water — clear your RAM" | Manually write down every thought, email, idea |
| **PARA** | "Second brain — store for later" | Manually drag files into 4 folders |
| **CODE / BASB** | "Capture → Organize → Distill → Express" | Manually highlight, save, tag, and summarize notes |
| **Readwise** | Reader + highlights sync | Manually read, highlight, then manually review |
| **Obsidian raw** | Daily notes + links | Manually write, tag, and backlink everything |

**The failure mode:** Your system is only as good as your willpower to feed it. You don't capture when you're busy, tired, or overwhelmed. The system goes stale exactly when you need it most.

### The OneMind Shift

```
LEGACY FRAMEWORKS          ONEMIND HIVE CAPTURE
─────────────────          ────────────────────
You → Inbox → Sort         World → Agents → Distilled Signal → You
Manual everywhere          Automated everywhere
You are the processor      You are the decision-maker
Input = willpower          Input = infrastructure
Breaks under pressure      Runs 24/7 regardless
```

**The Hive Capture System runs without you.** Agents harvest, agents distill, agents route. You receive processed intelligence — not raw noise. You only touch what requires your judgment.

---

## The 14 Input Sources

Everything that could carry signal for your life, organized by how it enters the system.

### Category A — Long-Form Knowledge
*Deep learning content. High signal density. Requires distillation.*

| # | Source | Format | Signal Type |
|---|--------|--------|-------------|
| 1 | **YouTube Videos** | Video/audio/transcript | Teaching, analysis, interviews |
| 2 | **Articles** | HTML/text | News, opinion, research |
| 3 | **Substack** | Email/web | Long-form newsletters |
| 4 | **Ebooks** | EPUB/MOBI | Books (digital format) |
| 5 | **Physical Books** | Paper → photo/scan | Books (offline) |
| 6 | **PDF/Docs** | PDF/DOCX | Reports, whitepapers, contracts |

### Category B — Live Stream
*High volume, low signal per item. Requires filtering before distillation.*

| # | Source | Format | Signal Type |
|---|--------|--------|-------------|
| 7 | **Social Media** | Text/image/video | Trends, ideas, connections |
| 8 | **RSS Feeds** | XML/JSON | Curated web updates |
| 9 | **Emails** | Text/HTML | Communications, newsletters, alerts |
| 10 | **Newsletters** | Email/web | Curated industry intel |
| 11 | **Bookmarks** | URL | Saved for later — varied types |

### Category C — Conversation Intelligence
*Already processed through AI. Highest distillation potential.*

| # | Source | Format | Signal Type |
|---|--------|--------|-------------|
| 12 | **ChatGPT Conversations** | Text/JSON | Research, ideas, drafts |
| 13 | **Claude Conversations** | Text/JSON | Strategy, analysis, code |

### Category D — Raw Input
*Unstructured. Requires the most transformation.*

| # | Source | Format | Signal Type |
|---|--------|--------|-------------|
| 14 | **Audio** | MP3/M4A/WAV | Voice memos, podcasts, meetings |
| — | **Web Scraping** | HTML/JSON | Structured data extraction |

---

## The 5-Step HARVEST Process

Every input — regardless of source — passes through the same 5 steps. This is the protocol.

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ HARVEST  │───▶│ DISTILL  │───▶│  ROUTE   │───▶│  STORE   │───▶│   ACT    │
│          │    │          │    │          │    │          │    │          │
│ Pull raw │    │ AI strips│    │ Classify │    │ Place in │    │ Task if  │
│ content  │    │ noise →  │    │ domain + │    │ correct  │    │ needed / │
│ into     │    │ extract  │    │ entity   │    │ Codex    │    │ surface  │
│ system   │    │ signal   │    │ type     │    │ location │    │ insight  │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
     │               │               │               │               │
  Automated       AI Agent        AI Agent        Automated       AI Notifies
  (tool/API)      (LEGACY)        (LEGACY)        (Codex write)   Zeus if needed
```

### Step 1 — HARVEST
*Getting the raw content into the system. No human intervention required.*

The agent or automation pulls the content in its native format and stores it temporarily in the processing queue.

**How each source enters:**

| Source | Harvest Mechanism | Tool |
|--------|------------------|------|
| YouTube Videos | URL → transcript via API; summary via AI | N8N + YouTube API + Whisper |
| Articles | Readwise Reader save → API → Codex | Readwise API / Jina.ai reader |
| Substack | Email forwarded → parsed | N8N email listener |
| Social Media | Saved post / monitored keyword | N8N + platform APIs |
| RSS Feeds | Polled on schedule | N8N RSS node |
| Emails | IMAP listener / forwarding rule | N8N email node |
| Newsletters | Email → parsed like Substack | N8N email node |
| Bookmarks | Browser extension → webhook | Raindrop.io API / custom |
| GPT Conversations | Export JSON → webhook | ChatGPT export + N8N |
| Claude Conversations | Project export → webhook | Claude export + N8N |
| Ebooks | Send to Kindle/Reader → highlights sync | Readwise API |
| Physical Books | Photo → OCR → text | iOS Shortcut + Vision API |
| PDF/Docs | File drop to watched folder | N8N file watcher |
| Audio | Upload / record → auto-transcribe | Whisper API / LiveKit |
| Web Scraping | Trigger URL → scraped content | N8N + Jina.ai / Playwright |

### Step 2 — DISTILL
*AI strips the noise and extracts the signal. This is where Hive earns its name.*

LEGACY agent processes raw content and produces a **Hive Extract** — a structured summary that contains only what matters.

**Hive Extract Format:**
```
SOURCE: [type] | [title/URL] | [date]
DOMAIN_GUESS: [HP / LE / GE / UI] | CONFIDENCE: [high/medium/low]
ENTITY_TYPE: [knowledge / task / project / insight / metric / tool]
─────────────────────────────────────────────────────────────────
CORE INSIGHT (1-3 sentences):
[What is the single most important thing here?]

KEY IDEAS (bullets):
- [idea 1]
- [idea 2]
- [idea 3]

ACTIONABLE: [yes/no]
IF YES:
  - TASK: [specific action]
  - PROJECT: [if bigger than 1 task]
  - URGENCY: [now / this week / someday]

QUOTES / VERBATIM (if valuable):
> [exact quote]

CONNECTIONS (links to existing Codex notes):
- [[related-note-1]]
- [[related-note-2]]

SIGNAL STRENGTH: [high / medium / low]
ZEUS NEEDS TO SEE: [yes/no] — [reason if yes]
```

**Distillation rules by source type:**

| Source Type | Distillation Focus | Typical Output |
|-------------|-------------------|----------------|
| Long-form (books, YouTube, articles) | Core ideas + key quotes + tasks generated | Knowledge note + possible project |
| Live stream (RSS, social, email) | Filter by relevance first → if relevant, extract insight | Insight note or task |
| Conversation (GPT/Claude exports) | Decisions made + knowledge generated + open loops | Knowledge note + task list |
| Audio (voice memo, podcast) | Verbatim transcript → then summarize | Note + tasks if present |
| PDF/Docs | Depends on type: contract → tasks; report → knowledge | Varies |

### Step 3 — ROUTE
*The Hive Extract is classified and assigned a destination.*

LEGACY routes every extract using two coordinates:

**Coordinate 1 — Domain:**
```
HP   → Does this improve Zeus as a person? (health, mind, skills, finance, identity)
LE   → Does this protect or grow the family/legacy?
GE   → Does this create business value or revenue?
UI   → Is this about systems, AI, infrastructure, or meta-level operations?
```

**Coordinate 2 — Entity Type:**
```
knowledge     → A piece of information worth keeping
insight       → A synthesized observation LEGACY surfaces proactively
task          → A discrete action to be done
project       → Coordinated effort with an outcome
tool          → A system, API, or integration
metric        → A number to track
```

**Routing matrix examples:**

| Input | Domain | Entity | Goes to |
|-------|--------|--------|---------|
| YouTube: "How to structure a 500-year family wealth plan" | LE | knowledge | `50-74 LE/` → Knowledge folder |
| Article: "New tax law affecting LLCs" | GE | task | Todoist → GE task: "Research LLC tax impact" |
| Podcast: Peter Attia on Zone 2 training | HP | knowledge | `25-49 HP/` → Knowledge folder |
| Claude convo about N8N automation | UI | knowledge + task | `00-24 UI/` + task in Todoist |
| Email: Invoice from vendor | GE | task | Todoist → GE task: "Pay invoice [vendor]" |
| Voice memo: "Idea for offer structure" | GE | project | GE → `Inbox` draft + maybe project |
| Newsletter: AI tool released | UI | insight | Surfaced to Zeus → `01 Command/` |

### Step 4 — STORE
*The routed extract is written into the Codex in the right location with full metadata.*

Every stored item gets the standard frontmatter:
```yaml
---
id: [domain]-[type]-[YYYYMMDD]-[slug]
title: "[Title]"
type: [knowledge | insight | task-log | project-seed]
source: "[original URL or filename]"
source_type: [youtube | article | podcast | email | conversation | book | pdf]
captured: [YYYY-MM-DD]
domain: [HP | LE | GE | UI]
tags: [domain/hp, source/youtube, topic/health, status/unreviewed]
signal_strength: [high | medium | low]
reviewed: false
ai_summary: "[1 sentence]"
---
```

**Storage locations:**
```
All domains have an Inbox subfolder as the first landing spot:
│
├── 06 Inbox (Queue)/06.1 Unprocessed/     ← before LEGACY processes it
│
After LEGACY routes it:
├── 25-49 HP/ → [relevant subfolder]       ← if HP
├── 50-74 LE/ → [relevant subfolder]       ← if LE
├── 75-99 GE/ → [relevant subfolder]       ← if GE
└── 00-24 UI/ → 04 Resources / 12 Research ← if UI
```

Items marked `reviewed: false` appear in the daily Zeus review queue in `01 Command/`.

### Step 5 — ACT
*If the extract has an action or insight worth Zeus's attention, LEGACY fires the appropriate action.*

**Three possible ACT outputs:**

```
1. SILENT STORE
   Signal strength low + not actionable
   → Stored, tagged, available for future RAG queries
   → Zeus never notified (doesn't need to be)

2. TASK CREATED
   Actionable item detected
   → Task created in Todoist with domain label + due date
   → Zeus sees it in normal task review, not interrupted

3. PROACTIVE INSIGHT
   High signal or time-sensitive
   → LEGACY messages Zeus via WhatsApp/Discord
   → Format: "📡 [Domain] Insight: [1 sentence] → [link to Codex note]"
   → Zeus can ignore, acknowledge, or act
```

---

## The Comparison That Matters

### How 2026 OneMind Beats Every Previous Approach

```
SCENARIO: You watch a 90-minute YouTube documentary on regenerative farming.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

GTD (2001):
  You'd have to remember to write down "watch farming doc" in your capture list
  Then manually decide if any tasks come from it
  90% chance you never get back to it

PARA / BASB (2017):
  You'd save the link to a "Resources" folder
  Maybe highlight key moments with Readwise Reader
  It sits there. You'd need to deliberately review it someday.
  Still manual. Still dependent on your energy.

CODE (2022):
  You Capture the link
  Organize the highlights
  Distill into your own words (manually)
  Express in a note or essay
  High quality — but requires 2-3 hours of your time per piece of content

READWISE (best existing tool):
  Auto-syncs highlights from Kindle, articles, etc.
  Daily review emails sent to you
  Still passive — you read highlights and hope they stick
  No routing, no tasks, no domain classification, no agent action

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

ONEMIND HIVE CAPTURE (2026):
  You paste the URL (or LEGACY catches it from your bookmark extension)
  LEGACY fetches the transcript via YouTube API
  LEGACY distills: 3 core ideas, 2 key quotes, 1 task ("Research soil prep methods for homestead")
  LEGACY routes it to LE (Legacy Domain — family homestead)
  Note stored in 50-74 LE/ with full metadata + backlinks to existing homestead notes
  Todoist task created: "Research soil prep methods for Zone 9" — tagged LE, due this week
  If signal is high: Zeus gets a WhatsApp message: "📡 LE Insight: Regenerative grazing cycles
  could apply to your homestead plan → [[farming-notes-20260304]]"

  Your time invested: ~15 seconds to paste the URL.
  System time: 3 minutes automated.
  Outcome: Knowledge stored, task created, homestead project advanced.
```

### The Adoption Curve

```
2001  GTD          → Capture everything manually. Pure human effort.
2007  Evernote     → Digital capture. Still manual tagging and filing.
2017  PARA         → Better structure, still manual, no AI.
2020  Roam/Obsidian → Bidirectional links. Still manual, just smarter storage.
2022  CODE / BASB  → Disciplined manual workflow. Readwise adds automation.
2023  AI notes apps → ChatGPT summaries added. Still siloed.
2024  Recall/Mem0  → Memory for AI sessions. No life OS integration.
2026  ONEMIND      → Fully automated capture pipeline feeding a unified life OS.
                     Agents run 24/7. You receive intelligence, not noise.
```

---

## Implementation Stack

The tools that make the Hive run:

```
CAPTURE LAYER          PROCESSING LAYER           OUTPUT LAYER
──────────────         ─────────────────          ─────────────
Raindrop.io            N8N (orchestration)        Obsidian (Codex storage)
Readwise Reader        LEGACY AI (Letta)          Todoist (task creation)
N8N email listener     Whisper API (transcription) Discord/WhatsApp (alerts)
iOS Shortcuts          Jina.ai (web reader)       01 Command (review queue)
Browser extension      OpenAI / Claude (distill)  Letta memory (RAG)
YouTube API            Vector DB (similarity)     Agent memory blocks
File watcher
```

---

## LEGACY's Role in Capture

LEGACY is not a passive storage system. In the capture layer, LEGACY acts as a **Hive Processor** — an always-on intelligence that sits between the world and Zeus.

**LEGACY's capture responsibilities:**

| Responsibility | Description |
|---------------|-------------|
| **Monitoring** | Watches the Inbox queue for new items 24/7 |
| **Distillation** | Runs extraction prompts against each new item |
| **Routing** | Applies domain + entity classification |
| **Deduplication** | Checks against existing Codex notes before storing |
| **Connection finding** | Identifies links to existing knowledge [[like this]] |
| **Task generation** | Creates Todoist tasks from actionable extracts |
| **Priority surfacing** | Flags high-signal items for Zeus's attention |
| **Memory indexing** | Adds to Letta archival memory for future RAG queries |

**Guardrails on capture (what LEGACY never does):**
- Never deletes original source content
- Never files something as reviewed without Zeus confirming
- Never creates a project (only tasks and project seeds — Zeus approves projects)
- Never shares captured content externally

---

## Metrics: How We Know the Hive is Working

The capture system reports these metrics to `01 Command/` weekly:

| Metric | What It Measures | Target |
|--------|-----------------|--------|
| **Items harvested** | Total inputs processed | Volume baseline |
| **Signal rate** | % marked high/medium signal | >40% |
| **Distillation lag** | Time from harvest to stored note | <30 min |
| **Task generation rate** | Tasks created per 10 items | 2-4 |
| **Review backlog** | Items in queue marked `reviewed: false` | <20 |
| **Domain distribution** | % routed to HP / LE / GE / UI | Reflects life balance |
| **Insight alerts** | Proactive messages sent to Zeus | 2-5/week |
| **RAG hits** | Future searches that surface capture items | Growing |

---

## Phase Rollout

Not everything gets automated on day one. This is the build order:

```
PHASE 1 — Manual Hive (Now)
  ✅ Zeus pastes URLs → LEGACY distills on demand
  ✅ Voice memo → LEGACY summarizes
  ✅ Email forwarding → LEGACY reads + routes
  → Output: Distillation speed 10x faster. Still somewhat manual entry.

PHASE 2 — Semi-Automated Hive
  ⬜ N8N RSS poller running
  ⬜ Readwise Reader API connected
  ⬜ Email inbox listener live
  ⬜ Bookmark extension → webhook → N8N → LEGACY
  → Output: 60-70% of captures automated.

PHASE 3 — Full Hive
  ⬜ YouTube URL auto-detected from browser history
  ⬜ GPT/Claude conversation exports automated
  ⬜ Physical book OCR shortcut built
  ⬜ Audio transcription pipeline live
  ⬜ Vector similarity deduplication running
  → Output: 90%+ capture automated. Zeus only touches high-signal items.

PHASE 4 — Predictive Hive
  ⬜ LEGACY recommends what to consume based on what's missing
  ⬜ "You haven't captured anything about [topic] and you have a goal that needs it"
  ⬜ Agent-initiated research fills knowledge gaps autonomously
  → Output: The system feeds itself.
```

---

*"The Hive doesn't sleep. It captures while you rest, distills while you work, routes while you think. By the time you open your command dashboard, the world has already been processed."*
