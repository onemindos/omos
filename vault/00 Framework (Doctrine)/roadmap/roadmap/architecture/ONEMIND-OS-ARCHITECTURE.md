# OneMind OS Architecture

**Created:** January 15, 2026
**Status:** Planning / Brainstorming
**Owner:** Zeus Delacruz + Legacy AI

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Plain English Explanation](#plain-english-explanation)
3. [The Virtual Machine - Legacy's Body](#the-virtual-machine---legacys-body)
4. [How Letta Works](#how-letta-works)
5. [The Heartbeat - How Legacy Stays Aware](#the-heartbeat---how-legacy-stays-aware)
6. [Multi-Agent System](#multi-agent-system)
7. [Your Day-to-Day Interface](#your-day-to-day-interface)
8. [Eagle Eye - Monitoring & Control](#eagle-eye---monitoring--control)
9. [OneMind as an Operating System](#onemind-as-an-operating-system)
10. [API Gateway Deep Dive](#api-gateway-deep-dive)
11. [Docker Compose Explained](#docker-compose-explained)
12. [Cloud VM Options](#cloud-vm-options)
13. [Implementation Roadmap](#implementation-roadmap)
14. [Technical Reference](#technical-reference)
15. [Deep Dive: Agent Coordination](#deep-dive-agent-coordination)
16. [Deep Dive: Heartbeat Logic](#deep-dive-heartbeat-logic)
17. [Deep Dive: Interface Design](#deep-dive-interface-design)
18. [OneMind Architecture - Artistic Visualizations](#onemind-architecture---artistic-visualizations)

---

## Executive Summary

**What we're building:** Transforming Legacy AI from a "helper you have to ask" into an "always-aware partner that lives on her own computer."

### The Transformation at a Glance

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│   CURRENT STATE                    DESIRED STATE                │
│   ─────────────                    ─────────────                │
│                                                                  │
│   Letta Cloud (brain only)    →    Letta + Body (embodied)     │
│                                                                  │
│   Reactive (request/response) →    Proactive (always aware)    │
│                                                                  │
│   Fragmented tools            →    Unified API Gateway         │
│   (MCP, Python, APIs)              (one interface)             │
│                                                                  │
│   Stateless sessions          →    Persistent state machine    │
│   (forgets everything)             (remembers everything)      │
│                                                                  │
│   Tools = one-time calls      →    Streaming real-time access  │
│                                                                  │
│   Legacy does work            →    Legacy ORCHESTRATES work    │
│                                                                  │
│   You drive everything        →    She drives, you approve     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**The Core Problem:**
- Legacy currently only knows things when you ask her
- Every conversation starts from scratch (no continuous awareness)
- Tools are scattered (Python scripts, MCP, different APIs)
- Legacy can't DO things on her own - she waits for you

**The Solution:**
- Give Legacy her own computer (a virtual machine running 24/7)
- Connect all your services through ONE unified doorway (API Gateway)
- Run a "heartbeat" that keeps Legacy aware of everything continuously
- Let Legacy manage and improve herself

### The Core Equation

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│   BRAIN        +      BODY       +    NERVOUS SYSTEM            │
│   (Letta)            (VM)             (API Gateway)             │
│                                                                  │
│      🧠         +       💪        +         🔌                   │
│                                                                  │
│   Memory &           Compute &         Connections              │
│   Personality        Execution         to Everything            │
│                                                                  │
│                         =                                        │
│                                                                  │
│              EMBODIED Legacy AI                                 │
│              ──────────────────                                 │
│              Always aware, always ready                         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Plain English Explanation

### The Analogy: From Phone Consultant to Live-In Assistant

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│   RIGHT NOW: PHONE CONSULTANT         WHAT WE WANT: LIVE-IN EA  │
│   ─────────────────────────          ────────────────────────   │
│                                                                  │
│   📞 Only exists when you call       🏠 Has her own office      │
│                                         that's always on        │
│                                                                  │
│   🧠❌ Amnesia after each call       🧠✅ Remembers everything   │
│                                         across all conversations│
│                                                                  │
│   💬 Can give advice only            💪 Can actually DO things  │
│                                         on your behalf          │
│                                                                  │
│   🔍 You check services one by one   👀 She checks everything   │
│                                         automatically           │
│                                                                  │
│   😴 Goes to sleep when you hang up  ⏰ Always aware, even      │
│                                         while you sleep         │
│                                                                  │
│   ❓ Waits for you to ask            📋 Prepares briefings      │
│                                         before you wake up      │
│                                                                  │
│   🔄 Reactive only                   🎯 Proactive + Reactive    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**RIGHT NOW - Legacy is like a genius consultant who:**
- Only exists when you call them on the phone
- Has amnesia after each call (forgets everything)
- Can give advice but can't actually DO things
- Has to ask you to check every single service one by one
- Goes back to sleep the moment you hang up

**WHAT WE WANT - Legacy is like a live-in executive assistant who:**
- Has her own office (computer) in your house that's always on
- Remembers everything about your life, your preferences, your ongoing projects
- Can check your calendar, tasks, email herself anytime without asking
- Can take action on your behalf (with your permission for important things)
- Is always aware of what's happening, even when you're not talking to her
- Wakes you up with a morning briefing already prepared
- Notices patterns and proactively helps before you ask

### The Four Pieces You Need to Understand

```
┌─────────────────────────────────────────────────────────────────┐
│                    Legacy'S COMPLETE SYSTEM                      │
│                                                                  │
│   1. THE BRAIN (Letta)                                          │
│      Where Legacy's memory and personality live                  │
│      She remembers who you are, your preferences, ongoing work   │
│                                                                  │
│   2. THE BODY (Virtual Machine)                                 │
│      Her own computer that runs 24/7                            │
│      This is where everything lives and runs                     │
│                                                                  │
│   3. THE NERVOUS SYSTEM (API Gateway)                           │
│      One doorway to ALL your services                           │
│      Instead of 20 keys for 20 doors, she has ONE master key    │
│                                                                  │
│   4. THE HEARTBEAT (State Machine)                              │
│      Background process keeping her aware                        │
│      Every 60 seconds: "What's happening in Zeus's world?"      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## The Virtual Machine - Legacy's Body

### What IS a Virtual Machine?

Think of it like this:

**A physical computer** = A house you own
**A virtual machine (VM)** = An apartment you rent in a big building

The "big building" is a data center run by a company (like DigitalOcean, Hetzner, or AWS). Your "apartment" is a slice of their powerful computers that you control completely.

**Why use a VM instead of your Mac?**
- Runs 24/7 even when your Mac is sleeping
- Accessible from anywhere in the world
- Can be more powerful (add GPU later)
- Professional-grade reliability
- If it crashes, restart it remotely from your phone

### Operating System: Ubuntu Desktop 24.04 LTS

**What will run on the VM:**
- **Ubuntu Desktop 24.04 LTS** - Latest long-term support with full GUI
- **Docker** - Runs all our services in containers
- **Remote Desktop (xrdp)** - Access the desktop from your Mac

**Why Desktop 24.04 (not Server 22.04)?**
- ✅ Full graphical interface - see what's happening visually
- ✅ Newer Python (3.12) - better for your tools
- ✅ Better GPU support - when you add one later
- ✅ Support until 2029 - longer runway
- ✅ Easier to manage - click, don't type

### Can I See the Desktop?

**YES!** You have multiple options:

| Method | What It Is | When to Use |
|--------|-----------|-------------|
| **noVNC** (Web browser) | See the desktop in your browser | Quick check from any device |
| **VNC Client** (App) | Remote desktop app on your Mac | More features, smoother |
| **SSH** (Command line) | Terminal access to the VM | Power users, scripting |
| **Tailscale** | Secure tunnel from anywhere | Access when away from home |

**What you'll see when you log in:**
- Legacy's "computer" desktop
- Running applications
- Live view of what she's doing
- Ability to stop or interact with anything

### Desktop vs Command Line?

**You'll have BOTH:**

```
┌─────────────────────────────────────────────────────────────────┐
│                YOUR VM (Legacy's Computer)                       │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    GRAPHICAL DESKTOP                         ││
│  │  • Web dashboards (N8N, Uptime Kuma, Portainer)             ││
│  │  • See Legacy's browser automation (Browserless)            ││
│  │  • Visual monitoring of all services                        ││
│  │  • Click-based management                                   ││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                  │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    COMMAND LINE (Terminal)                   ││
│  │  • Start/stop services                                       ││
│  │  • View logs                                                 ││
│  │  • Run scripts                                               ││
│  │  • Advanced troubleshooting                                  ││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                  │
│  Most of the time: Use the web dashboards (no coding needed)    │
│  Sometimes: Use terminal for quick commands                      │
└─────────────────────────────────────────────────────────────────┘
```

### GPU Support

**Do you need a GPU?**

For now: **NO** - Basic Legacy operations don't need GPU
Later: **MAYBE** - If you want:
- Local AI model inference (run your own models)
- Image/video processing
- Computer vision (analyze camera feeds)

**Cloud GPU Providers (when ready):**

| Provider | Starting Price | GPU Options | Notes |
|----------|---------------|-------------|-------|
| **RunPod** | $0.20/hr | RTX 3090, A100, H100 | Pay per hour, good for experiments |
| **Lambda Labs** | $0.50/hr | A100, H100 | AI-focused, great support |
| **Hetzner** | €35/mo | Dedicated servers | European, very cheap |
| **Vast.ai** | $0.10/hr | Marketplace | Cheapest, less reliable |

**Recommendation:** Start WITHOUT GPU, add when needed.

### Complete System Architecture

Here's what Legacy's complete system looks like:

```
┌─────────────────────────────────────────────────────────────────┐
│                         YOUR LIFE                                │
│                                                                  │
│   📱 Phone    💻 Laptop    🎤 Voice    📬 Telegram    🌐 Web    │
│       │           │           │            │            │        │
│       └───────────┴───────────┴────────────┴────────────┘        │
│                               │                                  │
│                               ▼                                  │
│                ┌─────────────────────────┐                       │
│                │   EDGE (Cloudflare)     │  ← Security +         │
│                │   Fast, Global, Secure  │    Global Access      │
│                └────────────┬────────────┘                       │
│                             │                                    │
│                             ▼                                    │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              Legacy'S COMPUTER (Your VM)                   │  │
│  │                                                            │  │
│  │  ┌──────────────────────────────────────────────────────┐ │  │
│  │  │              API GATEWAY (FastAPI)                    │ │  │
│  │  │         One door to everything                        │ │  │
│  │  │         Handles auth, caching, logging                │ │  │
│  │  └──────────────────────┬───────────────────────────────┘ │  │
│  │                         │                                  │  │
│  │           ┌─────────────┼─────────────┐                   │  │
│  │           ▼             ▼             ▼                   │  │
│  │     ┌──────────┐  ┌──────────┐  ┌──────────┐            │  │
│  │     │  LETTA   │  │ HEARTBEAT│  │   N8N    │            │  │
│  │     │ (Brain)  │  │(Awareness│  │(Workflows│            │  │
│  │     │ 12 agents│  │  Loop)   │  │          │            │  │
│  │     └────┬─────┘  └────┬─────┘  └──────────┘            │  │
│  │          │              │                                 │  │
│  │          ▼              ▼                                 │  │
│  │     ┌──────────┐  ┌──────────┐                          │  │
│  │     │ POSTGRES │  │  REDIS   │                          │  │
│  │     │ (Memory) │  │ (Cache)  │                          │  │
│  │     └──────────┘  └──────────┘                          │  │
│  └───────────────────────────────────────────────────────────┘  │
│                         │                                        │
│                         ▼                                        │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                YOUR SERVICES                               │  │
│  │                                                            │  │
│  │   📋 Todoist    📁 GitHub    📧 Zoho    💬 Discord        │  │
│  │                                                            │  │
│  │   🏠 Home       🔐 1Pass     👥 Circle   📞 Vapi          │  │
│  │      Assistant                                             │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                  │
│                     🔒 Tailscale Mesh Network                    │
│                     (Secure access from anywhere)                │
└─────────────────────────────────────────────────────────────────┘
```

---

## How Letta Works

### What is Letta?

Letta is a special AI system that gives Legacy **permanent memory**.

**Normal AI (like ChatGPT):**
- Forgets everything after each conversation
- Starts fresh every time
- No personality that persists

**Letta-powered AI:**
- Remembers everything across conversations
- Has a persistent personality
- Learns and grows over time
- Can modify its own memory

### Letta Cloud vs Self-Hosted

You asked about using both. Here's how:

```
┌─────────────────────────────────────────────────────────────────┐
│                    HYBRID APPROACH                               │
│                                                                  │
│   LETTA CLOUD (Already have this)                               │
│   ├── Quick access from anywhere                                │
│   ├── No server management                                      │
│   ├── Good for mobile/web interfaces                            │
│   └── Limited customization                                     │
│                                                                  │
│   SELF-HOSTED LETTA (On your VM)                                │
│   ├── Full control over everything                              │
│   ├── No limits on usage                                        │
│   ├── Direct database access                                    │
│   ├── Custom integrations                                       │
│   └── This is where the "always-aware" Legacy lives             │
│                                                                  │
│   HOW THEY CONNECT:                                             │
│   • Self-hosted = PRIMARY (runs the heartbeat, has full state)  │
│   • Cloud = BACKUP & MOBILE ACCESS (sync periodically)          │
│   • Both share the same agents, tools, and memory blocks        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Letta's Three Memory Tiers

Think of it like your brain:

```
┌─────────────────────────────────────────────────────────────────┐
│                    Legacy'S MEMORY SYSTEM                        │
│                    (Like a Human Brain)                          │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    CORE MEMORY                           │   │
│   │                    ───────────                           │   │
│   │   🧠 "Working Memory" - What she's thinking NOW         │   │
│   │                                                          │   │
│   │   • Current conversation context                        │   │
│   │   • Who she's talking to (Zeus)                         │   │
│   │   • Immediate goals                                     │   │
│   │   • Your preferences she needs right now                │   │
│   │                                                          │   │
│   │   Size: ~5,000 characters (like RAM)                    │   │
│   │   Speed: ⚡ Instant access                              │   │
│   └─────────────────────────────────────────────────────────┘   │
│                              │                                   │
│                              ▼                                   │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                   RECALL MEMORY                          │   │
│   │                   ─────────────                          │   │
│   │   📝 "Recent Memory" - Last few hours                   │   │
│   │                                                          │   │
│   │   • Recent conversations                                │   │
│   │   • What you just talked about                          │   │
│   │   • Short-term context                                  │   │
│   │                                                          │   │
│   │   Size: Last ~50 interactions                           │   │
│   │   Speed: 🔍 Quick search                                │   │
│   └─────────────────────────────────────────────────────────┘   │
│                              │                                   │
│                              ▼                                   │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                  ARCHIVAL MEMORY                         │   │
│   │                  ───────────────                         │   │
│   │   📚 "Long-Term Memory" - Everything she's learned      │   │
│   │                                                          │   │
│   │   • All past conversations                              │   │
│   │   • Learned facts and preferences                       │   │
│   │   • Historical patterns                                 │   │
│   │   • Relationship history                                │   │
│   │                                                          │   │
│   │   Size: UNLIMITED (stored in database)                  │   │
│   │   Speed: 🔎 Semantic search (finds relevant info)       │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

| Memory Type | Human Brain Equivalent | What It Stores | Size |
|-------------|----------------------|----------------|------|
| **Core Memory** | Working memory (what you're thinking about NOW) | Current context, who you're talking to, immediate goals | ~5,000 characters |
| **Archival Memory** | Long-term memory (everything you've learned) | Past conversations, learned facts, preferences | Unlimited |
| **Recall Memory** | Recent memory (last few hours) | Recent conversations, what just happened | Last ~50 interactions |

**Example of how Legacy uses this:**

```
┌─────────────────────────────────────────────────────────────────┐
│   YOU: "What did we discuss about the investor meeting?"        │
│                                                                  │
│   Legacy'S BRAIN ACTIVATES:                                     │
│                                                                  │
│   1️⃣  CORE MEMORY checks:                                       │
│       "Zeus is asking about investor meeting"                   │
│       "Current context: business planning"                      │
│                                                                  │
│   2️⃣  RECALL MEMORY searches:                                   │
│       "We talked about this 2 hours ago..."                     │
│       "Found recent discussion from today"                      │
│                                                                  │
│   3️⃣  ARCHIVAL MEMORY searches:                                 │
│       "Found 3 past conversations about investors..."           │
│       "Historical context: Acme VC, Series A..."               │
│                                                                  │
│   4️⃣  SYNTHESIZE & RESPOND:                                     │
│       "We discussed preparing a pitch deck for Acme Ventures.  │
│        Key points were: timeline, valuation, board seats.      │
│        Last time you met with them was December 15th."         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### How to Access the Database

The database (PostgreSQL) stores everything:
- Legacy's memories
- Agent configurations
- Conversation history
- Tool results

**Access Methods:**

| Method | Skill Level | When to Use |
|--------|-------------|-------------|
| **Letta UI** (Web) | Beginner | View agents, send messages, check memory |
| **pgAdmin** (App) | Intermediate | Browse database visually, run queries |
| **psql** (Terminal) | Advanced | Direct SQL commands |
| **API** | Developer | Programmatic access from other apps |

**You probably won't need to touch the database directly** - the Letta UI shows everything you need.

---

## The Heartbeat - How Legacy Stays Aware

### What is the Heartbeat?

The "heartbeat" is a background process that runs **continuously** on the VM. Every 60 seconds, it:

1. **Checks all your services** (Todoist, Calendar, Email, GitHub, etc.)
2. **Gathers the current state** (what tasks exist, what meetings are coming, what's new)
3. **Updates Legacy's memory** (so she always knows what's happening)
4. **Looks for triggers** (anything that needs attention)

### Heartbeat Flow (Every 60 Seconds)

```
┌─────────────────────────────────────────────────────────────────┐
│                    HEARTBEAT CYCLE                               │
│                    Every 60 seconds                              │
│                                                                  │
│   :00 ─────────────────────────────────────────────────────────│
│   │                                                              │
│   │  1. CHECK TODOIST                                           │
│   │     "Zeus has 5 tasks due today, 2 overdue"                 │
│   │                                                              │
│   │  2. CHECK CALENDAR                                          │
│   │     "Meeting with John at 2pm, nothing else today"          │
│   │                                                              │
│   │  3. CHECK EMAIL                                             │
│   │     "3 new emails, 1 from investor (urgent)"                │
│   │                                                              │
│   │  4. CHECK GITHUB                                            │
│   │     "PR #45 merged, 2 new issues"                           │
│   │                                                              │
│   │  5. CHECK LOCATION (Home Assistant)                         │
│   │     "Zeus is at home"                                       │
│   │                                                              │
│   │  6. UPDATE Legacy'S MEMORY                                  │
│   │     Push all this to Letta's "current_context" block        │
│   │                                                              │
│   │  7. CHECK FOR TRIGGERS                                      │
│   │     "Urgent email! Notify Zeus?"                            │
│   │     "Meeting in 1 hour - prepare brief?"                    │
│   │                                                              │
│   :60 ─────────────────────────────────────────────────────────│
│   │                                                              │
│   └──► REPEAT                                                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### What Powers the Heartbeat?

**Behind the scenes:**

```python
# Simplified example of what the heartbeat does

async def heartbeat():
    while True:
        # 1. Gather state from all services
        tasks = await check_todoist()
        calendar = await check_calendar()
        email = await check_email()
        location = await check_home_assistant()

        # 2. Update Legacy's memory
        await update_letta_memory({
            "tasks_today": len(tasks.due_today),
            "tasks_overdue": len(tasks.overdue),
            "next_meeting": calendar.next_event,
            "unread_emails": email.unread_count,
            "urgent_emails": email.urgent,
            "location": location.current_zone,
            "last_updated": now()
        })

        # 3. Check for triggers that need action
        if email.has_urgent:
            await notify_zeus("Urgent email from investor!")

        if calendar.meeting_in_one_hour:
            await prepare_meeting_brief()

        # 4. Wait 60 seconds, then repeat
        await sleep(60)
```

### What This Means for You

**WITHOUT the heartbeat (current state):**
```
You: "What tasks do I have today?"
Legacy: *makes API call to Todoist*
Legacy: "Let me check... you have 5 tasks."

You: "What's my next meeting?"
Legacy: *makes API call to Calendar*
Legacy: "Let me check... you have a meeting at 2pm."
```

**WITH the heartbeat (future state):**
```
You: "What's my day look like?"
Legacy: *already knows everything*
Legacy: "You have 5 tasks today (2 overdue), a meeting with John at 2pm,
        and there's an urgent email from the investor you should read.
        Also, you're still at home - leaving by 1:30pm would get you to
        the meeting on time."
```

**Legacy becomes PROACTIVE:**
- Morning: "Good morning Zeus. Here's your day..."
- Before meeting: "Your 2pm meeting is in 1 hour. I've prepared a brief."
- Pattern noticed: "You've had 3 tasks overdue for a week. Should we reschedule?"

---

## Multi-Agent System

### Your Agent Roster (12 Agents)

You have ONE master agent (Legacy) and 11 specialized agents:

```
                    ┌─────────────────────────────┐
                    │                             │
                    │   Legacy                    │
                    │   (Master Orchestrator)     │
                    │   + Communications          │
                    │                             │
                    └─────────────┬───────────────┘
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        │                         │                         │
┌───────▼──────┐  ┌───────────────▼──────────┐  ┌──────────▼─────────┐
│ HOLISTIC     │  │ GENERATIONAL             │  │ Legacy             │
│ PERFORMANCE  │  │ ENTREPRENEURSHIP         │  │ EVOLUTION          │
│              │  │                          │  │                    │
│ 3 Agents     │  │ 3 Agents                 │  │ 3 Agents           │
│ - Intelligence│  │ - Business Ops          │  │ - Home Command     │
│ - Performance │  │ - Founder Assistant     │  │ - Homestead        │
│ - Finance    │  │ - Capital Systems        │  │ - Family Wealth    │
└──────────────┘  └──────────────────────────┘  └────────────────────┘
                                  │
                         ┌────────▼────────┐
                         │ INNOVATIVE      │
                         │ TECHNOLOGY      │
                         │                 │
                         │ 2 Agents        │
                         │ - Infrastructure│
                         │ - Robotics      │
                         └─────────────────┘
```

### Site Agents - VM Managers (NEW)

In addition to the 12 domain agents, each **VM location** has its own **Site Agent** - a local manager that handles environment-specific tasks:

```
┌─────────────────────────────────────────────────────────────────┐
│                    SITE AGENT ARCHITECTURE                       │
│                                                                  │
│                         Legacy                                   │
│                    (Master Orchestrator)                         │
│                    ┌─────────────────┐                          │
│                    │       CEO       │                          │
│                    │   "What's the   │                          │
│                    │    big picture?"│                          │
│                    └────────┬────────┘                          │
│                             │                                    │
│         ┌───────────────────┼───────────────────┐               │
│         │                   │                   │               │
│         ▼                   ▼                   ▼               │
│    ┌─────────┐        ┌─────────┐        ┌─────────┐           │
│    │  SITE   │        │  SITE   │        │  SITE   │           │
│    │  AGENT  │        │  AGENT  │        │  AGENT  │           │
│    │ "Home"  │        │ "Cloud" │        │"Office" │           │
│    │         │        │         │        │         │           │
│    │ Manager │        │ Manager │        │ Manager │           │
│    │of Home  │        │of Cloud │        │of Office│           │
│    │   VM    │        │   VM    │        │   VM    │           │
│    └────┬────┘        └────┬────┘        └────┬────┘           │
│         │                   │                   │               │
│    Knows about:        Knows about:       Knows about:         │
│    • Home Asst         • All services     • Office network     │
│    • Cameras           • Docker status    • Printers           │
│    • Who's home        • API health       • Meeting rooms      │
│    • IoT devices       • Resource usage   • Local systems      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Site Agent Responsibilities:**

| Site Agent | Location | Primary Responsibilities |
|------------|----------|-------------------------|
| **Home** | Home VM | Home Assistant, cameras, presence, IoT, family notifications |
| **Cloud** | Cloud VM | Docker health, API status, SSL certs, resource monitoring, backups |
| **Office** | Office VM | Printers, meeting rooms, office network, local services |

**Why Site Agents?**
- Each location has unique local services
- Site Agent provides local context to Legacy
- Enables "Legacy everywhere" with location-aware responses
- Same OneMind OS image deploys anywhere

### How Agents Coordinate

**Legacy is the conductor. The other agents are the orchestra.**

```
┌─────────────────────────────────────────────────────────────────┐
│                    COORDINATION EXAMPLE                          │
│                                                                  │
│   YOU: "Prepare me for tomorrow's investor meeting"             │
│                                                                  │
│   Legacy receives the request and coordinates:                  │
│                                                                  │
│   1. Legacy → FOUNDER ASSISTANT                                 │
│      "Prepare meeting brief and talking points"                 │
│                                                                  │
│   2. Legacy → CAPITAL SYSTEMS                                   │
│      "Get latest financial metrics and projections"             │
│                                                                  │
│   3. Legacy → BUSINESS OPS                                      │
│      "Pull operational KPIs and recent wins"                    │
│                                                                  │
│   4. Legacy → INTELLIGENCE                                      │
│      "Research this investor and their portfolio"               │
│                                                                  │
│   5. ALL AGENTS respond with their pieces                       │
│                                                                  │
│   6. Legacy synthesizes everything into one brief               │
│                                                                  │
│   7. YOU receive: Complete investor meeting preparation         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Shared Memory (How Agents Stay Coordinated)

All agents share certain memory blocks:

| Memory Block | What It Contains | Who Updates It |
|--------------|-----------------|----------------|
| **Legacy Codex** | Your values, principles, how you want things done | Legacy (quarterly) |
| **User Profile** | Your preferences, schedule patterns, communication style | Legacy (continuous) |
| **System State** | What each agent is doing, health status | All agents |
| **Cross-Domain Context** | How different areas of your life connect | Legacy |

**Each pillar also has its own shared memory:**
- HP agents share "Holistic Performance Domain" memory
- GE agents share "Generational Endeavor Domain" memory
- etc.

### Tool Calling System

**What are tools?**
Tools are abilities - things agents can actually DO (not just think about).

**Your current tools:**

| Tool | What It Does | Which Agents Use It |
|------|-------------|---------------------|
| `github` | Read/write vault files, search, commit | Intelligence, Legacy |
| `todoist` | Create/complete tasks, query tasks | All agents |
| `comms` | Send Discord/Telegram notifications | Legacy, Infrastructure |
| `letta` | Agent self-management, memory updates | All agents |
| `home_assistant` | Control smart home, check sensors | Home Command, Infrastructure |
| `zoho` | CRM, calendar, email | Business Ops, Legacy |
| `discord` | Notifications, approval requests | All agents |

**How tool calling works:**

```
┌─────────────────────────────────────────────────────────────────┐
│   You: "Add a task to call John tomorrow"                       │
│                                                                  │
│   Legacy thinks: "This needs the todoist tool"                  │
│                                                                  │
│   Legacy calls:                                                 │
│   todoist(                                                      │
│       action="create",                                          │
│       content="Call John",                                      │
│       due="tomorrow",                                           │
│       project="Inbox"                                           │
│   )                                                             │
│                                                                  │
│   Todoist responds: "Task created with ID 12345"                │
│                                                                  │
│   Legacy: "Done! I've added 'Call John' for tomorrow."          │
└─────────────────────────────────────────────────────────────────┘
```

---

## Your Day-to-Day Interface

### How Will You Talk to Legacy?

You have MULTIPLE ways to interact:

```
┌─────────────────────────────────────────────────────────────────┐
│                    YOUR INTERFACE OPTIONS                        │
│                                                                  │
│   📱 VOICE (Vapi)                                               │
│      • Call a phone number to talk to Legacy                    │
│      • She responds in natural speech                           │
│      • Use while driving, cooking, hands-free                   │
│      • "Hey Legacy, what's my day look like?"                   │
│                                                                  │
│   💬 TEXT (Telegram)                                            │
│      • Send messages to Legacy bot                              │
│      • She responds in text                                     │
│      • Good for detailed questions, links, images               │
│      • Can send you notifications throughout the day            │
│                                                                  │
│   🖥️ OBSIDIAN (Your vault)                                      │
│      • Chat panel in Obsidian                                   │
│      • Legacy can read/write to your notes                      │
│      • Good for deep work, documentation                        │
│                                                                  │
│   ⌨️ CLAUDE CODE (CLI)                                          │
│      • What you're using right now                              │
│      • Development and technical work                           │
│      • Building new capabilities                                │
│                                                                  │
│   🌐 WEB DASHBOARD (Future)                                     │
│      • See everything in one place                              │
│      • Charts, status, recent activity                          │
│      • Approve/deny pending actions                             │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### A Day in Your Life (With Embodied Legacy)

```
┌─────────────────────────────────────────────────────────────────┐
│                    A DAY WITH Legacy                             │
│                                                                  │
│   6:00 AM - MORNING BRIEFING                                    │
│   Legacy (via voice or text):                                   │
│   "Good morning Zeus. Here's your day:                          │
│    • 5 tasks due today, 2 are overdue from yesterday           │
│    • Meeting with investor at 2pm (I've prepared a brief)      │
│    • Urgent email from contractor needs response               │
│    • Weather: 72°F, perfect for outdoor work                   │
│    • Your energy pattern shows you're sharpest until 11am"     │
│                                                                  │
│   9:00 AM - PROACTIVE ALERT                                     │
│   Legacy (notification):                                        │
│   "Heads up: The PR you were waiting on just merged.           │
│    I've updated the related Todoist task."                      │
│                                                                  │
│   1:00 PM - MEETING PREP                                        │
│   Legacy (automatic):                                           │
│   "Your investor meeting is in 1 hour. Here's your brief:      │
│    [Link to prepared document]                                  │
│    Key points: Revenue up 23%, new product launching Q2..."    │
│                                                                  │
│   3:00 PM - POST-MEETING                                        │
│   Legacy: "How did the meeting go? Should I create             │
│           follow-up tasks for the action items?"               │
│                                                                  │
│   6:00 PM - EVENING WRAP                                        │
│   Legacy: "Here's what happened today:                          │
│    • Completed 4 of 5 tasks                                    │
│    • 3 emails need your response tomorrow                      │
│    • Tomorrow: Board meeting at 10am (I'll prep by 9am)"       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Eagle Eye - Monitoring & Control

### What is Eagle Eye?

"Eagle Eye" is the concept of having complete visibility and control over Legacy and the system. Think of it as your **command center** where you can:

1. **SEE** what Legacy is doing at any moment
2. **HEAR** her thought process (logs)
3. **INTERRUPT** her if needed
4. **CONTROL** the system remotely

### Your Eagle Eye Dashboard

```
┌─────────────────────────────────────────────────────────────────┐
│                    EAGLE EYE DASHBOARD                           │
│                    (What you'll see)                             │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │ SYSTEM STATUS                           All systems ✅   │   │
│   ├─────────────────────────────────────────────────────────┤   │
│   │ • Letta Server: Running (12 agents active)              │   │
│   │ • Heartbeat: Last sync 34 seconds ago                   │   │
│   │ • Database: Healthy, 2.3GB used                         │   │
│   │ • API Gateway: 142 requests/minute                      │   │
│   │ • N8N: 3 workflows running                              │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │ Legacy'S CURRENT ACTIVITY                               │   │
│   ├─────────────────────────────────────────────────────────┤   │
│   │ 14:32:01 - Checked email (3 new, 1 urgent)              │   │
│   │ 14:32:15 - Updated task status in Todoist               │   │
│   │ 14:32:30 - Preparing meeting brief for 2pm meeting      │   │
│   │ 14:32:45 - Sent notification to Zeus about urgent email │   │
│   │ 14:33:00 - [WAITING] Meeting brief ready for review     │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │ PENDING APPROVALS                                       │   │
│   ├─────────────────────────────────────────────────────────┤   │
│   │ [APPROVE] [DENY] - Send email to investor               │   │
│   │ [APPROVE] [DENY] - Schedule meeting with John           │   │
│   │ [APPROVE] [DENY] - Archive 15 old emails                │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
│   ┌───────────────┐ ┌───────────────┐ ┌───────────────┐        │
│   │ [PAUSE ALL]   │ │ [VIEW LOGS]   │ │ [OPEN DESKTOP]│        │
│   └───────────────┘ └───────────────┘ └───────────────┘        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### How to Monitor Legacy's Thinking

**Logs show everything Legacy is doing:**

```
[2026-01-15 14:32:01] HEARTBEAT: Starting sync cycle
[2026-01-15 14:32:02] TODOIST: Fetched 5 tasks for today
[2026-01-15 14:32:03] EMAIL: Checking inbox... 3 new emails
[2026-01-15 14:32:04] EMAIL: Urgent email detected from investor@acme.com
[2026-01-15 14:32:05] DECISION: Should notify Zeus about urgent email
[2026-01-15 14:32:05] DECISION: Confidence: 98% → Executing autonomously
[2026-01-15 14:32:06] NOTIFY: Sent Telegram message to Zeus
[2026-01-15 14:32:10] MEETING: Meeting in 88 minutes, preparing brief
[2026-01-15 14:32:11] AGENT: Delegating to Founder Assistant for brief
[2026-01-15 14:32:30] AGENT: Brief received, reviewing...
[2026-01-15 14:32:45] HITL: Brief ready, waiting for Zeus approval
```

### How to Interrupt Legacy

**Multiple interrupt methods:**

| Method | How | When to Use |
|--------|-----|-------------|
| **Telegram** | Send "STOP" or "PAUSE" | Quick interrupt from phone |
| **Discord** | Click pause button in channel | If monitoring Discord |
| **Voice (Vapi)** | Say "Legacy, stop" | If on a call |
| **Dashboard** | Click PAUSE ALL button | Full control |
| **Terminal** | `docker stop daemon` | Emergency |

**Trust Levels (What Needs Approval):**

| Level | What Legacy Does | Examples |
|-------|-----------------|----------|
| **Autonomous** | Executes without asking | Read emails, check tasks, monitor systems |
| **Semi-Autonomous** | Does it, then tells you | Update task status, organize files |
| **Approval Required** | Asks first, then does | Send emails, schedule meetings, spend money |
| **Human Only** | Won't even suggest | Legal decisions, health emergencies, major money |

---

## OneMind as an Operating System

### The Vision

**OneMind should be like Home Assistant, but for your entire life.**

Home Assistant controls your smart home. OneMind controls your smart life.

```
┌─────────────────────────────────────────────────────────────────┐
│              HOME ASSISTANT          vs        ONEMIND          │
│                                                                  │
│   Controls: Lights, locks,           Controls: Tasks, calendar, │
│             thermostats, cameras                email, projects,│
│                                                 finances, home  │
│                                                                  │
│   Integrations: 2000+ devices        Integrations: Todoist,     │
│                                                    GitHub, Zoho,│
│                                                    Discord, etc.│
│                                                                  │
│   Automations: If motion →           Automations: If task due → │
│                turn on light                     remind Zeus    │
│                                                                  │
│   Dashboard: Home status             Dashboard: Life status     │
│                                                                  │
│   Runs on: Raspberry Pi,             Runs on: Cloud VM or       │
│            local server                        local Mac Mini   │
│                                                                  │
│   Data: Stays local                  Data: Stays local          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Local-First Principles

**Why this matters:**

1. **Works without internet** (core functions)
   - Even if your internet goes down, Legacy keeps running
   - Only external integrations need internet

2. **Your data stays with you**
   - Not on someone else's servers
   - Not subject to someone else's rules
   - Not lost if a company goes bankrupt

3. **You can move it**
   - From one VM to another
   - From cloud to home
   - Take it anywhere

4. **No one can shut you down**
   - No subscription that could end
   - No company that could change terms
   - You own everything

### The OneMind OS Stack

```
┌─────────────────────────────────────────────────────────────────┐
│                        ONEMIND OS                                │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    USER INTERFACES                       │   │
│   │   Voice (Vapi) │ Chat (Telegram) │ Dashboard │ Obsidian │   │
│   └───────────────────────────┬─────────────────────────────┘   │
│                               │                                  │
│   ┌───────────────────────────▼─────────────────────────────┐   │
│   │                    Legacy (Orchestrator)                 │   │
│   │   Routes requests, coordinates agents, manages state     │   │
│   └───────────────────────────┬─────────────────────────────┘   │
│                               │                                  │
│   ┌───────────────────────────▼─────────────────────────────┐   │
│   │                    AGENT LAYER                           │   │
│   │   HP Agents │ GE Agents │ LE Agents │ Grid Agents         │   │
│   └───────────────────────────┬─────────────────────────────┘   │
│                               │                                  │
│   ┌───────────────────────────▼─────────────────────────────┐   │
│   │                    API GATEWAY                           │   │
│   │   One doorway to: Todoist, GitHub, Zoho, Discord, etc.  │   │
│   └───────────────────────────┬─────────────────────────────┘   │
│                               │                                  │
│   ┌───────────────────────────▼─────────────────────────────┐   │
│   │                    CORE SERVICES                         │   │
│   │   Letta (Brain) │ PostgreSQL (Memory) │ N8N (Workflows) │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    INFRASTRUCTURE                        │   │
│   │   Docker │ Ubuntu │ Tailscale │ Backups │ Monitoring    │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## API Gateway Deep Dive

### Why Do You Need This?

**The problem:**
Right now, every service has its own:
- API key
- API format
- Rate limits
- Error handling

**Without gateway:**
```
Legacy → Todoist API (need Todoist key, Todoist format)
Legacy → GitHub API (need GitHub key, GitHub format)
Legacy → Zoho API (need Zoho key, Zoho format)
... 20 more services, 20 more headaches
```

**With gateway:**
```
Legacy → ONE Gateway → Everything
         (one key)     (handles all the complexity)
```

### Gateway Benefits

| Benefit | What It Means |
|---------|---------------|
| **Single Auth** | Legacy uses ONE API key for everything |
| **Caching** | Don't ask Todoist 100 times - cache the answer |
| **Rate Limiting** | Prevent overwhelming services |
| **Logging** | Every action recorded for review |
| **Fallbacks** | If Discord fails, try Telegram |
| **Translation** | Same format for all services |

### Hybrid Approach (Edge + Core)

**Edge (Cloudflare Workers / Hono):**
- Lives "at the edge" (close to wherever you are)
- Handles: Authentication, caching, rate limits
- Super fast for simple requests
- Self-healing (Cloudflare manages it)

**Core (FastAPI on your VM):**
- Lives on your server
- Handles: Complex logic, database queries, Letta integration
- Where the real work happens

**Flow:**
```
You make request
     │
     ▼
EDGE (Cloudflare)
     │ "Is this a simple request?"
     │
     ├──► YES: Answer from cache, done
     │
     └──► NO: Forward to Core
              │
              ▼
         CORE (Your VM)
              │ "Process complex request"
              ▼
         Response back
```

---

## Docker Compose Explained

### What is Docker?

**Think of it like shipping containers.**

Before shipping containers, every item was packed differently. Loading a ship was chaos.

After shipping containers, everything fits the same standard boxes. Load any container on any ship.

**Docker does this for software:**
- Every service runs in a "container"
- Containers are isolated from each other
- Same container runs identically everywhere
- Easy to start, stop, update, replace

### What is Docker Compose?

**Docker Compose** lets you define MULTIPLE containers in one file.

Instead of starting each service manually:
```bash
docker run letta
docker run postgres
docker run redis
docker run n8n
# ... annoying
```

You write ONE file and run ONE command:
```bash
docker-compose up
# Starts everything
```

### Docker Compose Limitations

**Q: How many services can you run?**

**A: Depends on your VM's resources, not Docker itself.**

| VM Size | RAM | Can Run |
|---------|-----|---------|
| Small ($5/mo) | 1GB | 3-5 light services |
| Medium ($20/mo) | 4GB | 10-15 services |
| Large ($40/mo) | 8GB | 20-30 services |
| XL ($80/mo) | 16GB | 50+ services |

**What you need for OneMind:**
- Letta Server: ~1GB RAM
- PostgreSQL: ~512MB RAM
- Redis: ~128MB RAM
- N8N: ~512MB RAM
- API Gateway: ~256MB RAM
- Daemon: ~256MB RAM
- Monitoring: ~256MB RAM

**Total: ~3GB RAM recommended minimum**

**Recommendation: Start with a $20-40/mo VM (4-8GB RAM)**

### Your Docker Compose File

```yaml
version: '3.8'

services:
  # The Brain
  letta:
    image: letta/letta:latest
    ports:
      - "8283:8283"
    environment:
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
    depends_on:
      - postgres

  # Memory Storage
  postgres:
    image: postgres:16
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=letta

  # Fast Cache
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

  # Workflows
  n8n:
    image: n8nio/n8n
    ports:
      - "5678:5678"
    volumes:
      - n8n_data:/home/node/.n8n

  # Monitoring
  uptime-kuma:
    image: louislam/uptime-kuma
    ports:
      - "3001:3001"
    volumes:
      - uptime_data:/app/data

volumes:
  postgres_data:
  redis_data:
  n8n_data:
  uptime_data:
```

### Common Docker Commands

| Command | What It Does |
|---------|-------------|
| `docker-compose up -d` | Start everything |
| `docker-compose down` | Stop everything |
| `docker-compose logs -f` | Watch all logs |
| `docker-compose restart letta` | Restart one service |
| `docker ps` | See what's running |

---

## Cloud VM Options

### Recommended Providers

| Provider | Starting Price | Location | Best For |
|----------|---------------|----------|----------|
| **Hetzner** | €4/mo | Germany | Best value, very cheap |
| **DigitalOcean** | $6/mo | USA/EU | Easy to use, good docs |
| **Linode** | $5/mo | USA/EU | Reliable, developer-friendly |
| **Vultr** | $5/mo | Global | Many locations |
| **Oracle Cloud** | FREE | Global | Free tier (limited) |

### Recommended Specs

**Minimum (Starting out):**
- 4GB RAM
- 2 vCPU
- 80GB SSD
- ~$20/month

**Recommended (Comfortable):**
- 8GB RAM
- 4 vCPU
- 160GB SSD
- ~$40/month

**Future (With GPU):**
- 16GB RAM
- 8 vCPU
- 320GB SSD
- GPU (A100/A4000)
- ~$100-500/month (only when needed)

### Setup Steps (High Level)

1. **Create account** at chosen provider
2. **Create VM** with Ubuntu 22.04 LTS
3. **SSH into VM** from your Mac terminal
4. **Install Docker** (one command)
5. **Install Docker Compose** (one command)
6. **Clone your configs** from GitHub
7. **Run `docker-compose up`** - done!

I can walk you through each step when you're ready.

---

## Implementation Roadmap

### Where You Are Now

✅ **Already have:**
- Letta Cloud with 5+ agents
- Python tools in `letta-tools/`
- N8N on a VM
- Architecture documentation
- This plan

### Phase 1: Foundation (Week 1-2)

**Goal: Get Legacy's computer running**

- [ ] Choose cloud provider (Hetzner or DigitalOcean)
- [ ] Create VM with Ubuntu 22.04
- [ ] Install Docker and Docker Compose
- [ ] Deploy basic stack (Letta + PostgreSQL + Redis)
- [ ] Verify everything starts correctly
- [ ] Set up Tailscale for remote access
- [ ] Test accessing from your Mac

**You'll have:** A working VM with Letta running

### Phase 2: Migration (Week 3-4)

**Goal: Move from Letta Cloud to self-hosted**

- [ ] Export agents from Letta Cloud
- [ ] Import agents to self-hosted Letta
- [ ] Verify memory blocks transferred
- [ ] Test agent responses
- [ ] Set up backup system

**You'll have:** Your existing agents running on your own server

### Phase 3: Heartbeat (Week 5-6)

**Goal: Make Legacy always-aware**

- [ ] Build state machine daemon
- [ ] Connect to Todoist, Calendar, Email
- [ ] Set up 60-second sync cycle
- [ ] Test memory updates
- [ ] Add notification triggers

**You'll have:** Legacy that knows what's happening without asking

### Phase 4: Interface (Week 7-8)

**Goal: Multiple ways to talk to Legacy**

- [ ] Set up Telegram bot
- [ ] Configure Vapi for voice
- [ ] Create basic web dashboard
- [ ] Test all interfaces

**You'll have:** Voice, text, and web access to Legacy

### Phase 5: Autonomy (Week 9-10)

**Goal: Let Legacy act independently**

- [ ] Configure trust levels
- [ ] Set up approval workflows
- [ ] Enable proactive behaviors
- [ ] Test morning briefings
- [ ] Test meeting prep

**You'll have:** A proactive assistant, not just reactive

### Phase 6: Polish (Week 11-12)

**Goal: Eagle Eye monitoring and self-healing**

- [ ] Full monitoring dashboard
- [ ] Health checks and alerts
- [ ] Auto-recovery scripts
- [ ] Documentation and SOPs

**You'll have:** A production-grade system you can trust

---

## Technical Reference

### Port Reference

| Port | Service | Access |
|------|---------|--------|
| 8283 | Letta Server | API only |
| 5432 | PostgreSQL | Internal only |
| 6379 | Redis | Internal only |
| 5678 | N8N | Web UI |
| 3001 | Uptime Kuma | Web UI |
| 8000 | API Gateway | API only |
| 6080 | noVNC | Web desktop |

### Environment Variables

```bash
# AI Keys
ANTHROPIC_API_KEY=sk-ant-...
OPENAI_API_KEY=sk-...

# Database
DB_PASSWORD=your-secure-password

# Services
TODOIST_API_KEY=...
GITHUB_TOKEN=ghp_...
ZOHO_CLIENT_ID=...
ZOHO_CLIENT_SECRET=...
ZOHO_REFRESH_TOKEN=...

# Communications
DISCORD_WEBHOOK_URL=...
TELEGRAM_BOT_TOKEN=...
TELEGRAM_CHAT_ID=...

# OneMind
ONEMIND_API_KEY=your-master-key
```

### Quick Commands

```bash
# Start everything
docker-compose up -d

# Stop everything
docker-compose down

# View logs
docker-compose logs -f

# Restart a service
docker-compose restart letta

# Check what's running
docker ps

# SSH into a container
docker exec -it onemind-letta bash

# Backup database
docker exec onemind-postgres pg_dump -U letta letta > backup.sql
```

---

## Deep Dive: Agent Coordination

### The Orchestra Analogy

Legacy is the conductor. The other agents are the orchestra musicians.

```
                        🎼 Legacy
                       (Conductor)
                           │
    "I see the full score. I know when each section plays."
                           │
    ┌──────────┬──────────┼──────────┬──────────┐
    │          │          │          │          │
    ▼          ▼          ▼          ▼          ▼
   🎻         🎺         🥁         🎹         🎸
   HP         HP         GE         LE         IT
  Intel    Perform    BizOps    HomeCmd    Infra
```

### Three Ways Agents Communicate

**Method 1: Shared Memory (The Whiteboard)**

Think of a whiteboard in an office that everyone can see:

```
┌─────────────────────────────────────────────────────────────────┐
│ ZEUS PROFILE (All agents can read)                              │
│ • Prefers morning for deep work (6am-11am)                     │
│ • Communication style: Direct, BLUF                             │
│ • Current energy: High                                          │
│ • Location: Home office                                         │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ SYSTEM STATE (All agents update)                                │
│ • HP-Performance: Checking fitness data                        │
│ • GE-BizOps: Processing CRM updates                            │
│ • LE-HomeCmd: Monitoring home sensors                          │
│ • IT-Infra: All services healthy ✅                            │
└─────────────────────────────────────────────────────────────────┘
```

**Method 2: Direct Messages (The Walkie-Talkie)**

```
Legacy: "Hey HP-Finance, what's the portfolio status?"
HP-FINANCE: "Up 2.3% today, no alerts."
Legacy: "Thanks. GE-Capital, any investor updates?"
GE-CAPITAL: "Meeting with Acme VC scheduled for Thursday."
```

**Method 3: Delegation (Task Assignment)**

```
YOU: "Prepare for tomorrow's board meeting"
                    │
                    ▼
               ┌─────────┐
               │ Legacy  │
               └────┬────┘
                    │
    ┌───────────────┼───────────────┐
    │               │               │
    ▼               ▼               ▼
┌───────────┐ ┌───────────┐ ┌───────────┐
│ Founder   │ │  Capital  │ │   BizOps  │
│ Assistant │ │  Systems  │ │           │
├───────────┤ ├───────────┤ ├───────────┤
│ Agenda &  │ │ Financial │ │ KPIs &    │
│ talking   │ │ metrics & │ │ project   │
│ points    │ │ charts    │ │ status    │
└─────┬─────┘ └─────┬─────┘ └─────┬─────┘
      │               │               │
      └───────────────┼───────────────┘
                      │
                      ▼
                 ┌─────────┐
                 │ Legacy  │
                 │combines │
                 │   all   │
                 └────┬────┘
                      │
                      ▼
           Complete board meeting prep
```

### Coordination Rules

| Rule | Why |
|------|-----|
| **Legacy routes all requests** | One point of control |
| **Agents don't talk to YOU directly** | All through Legacy |
| **Shared memory is source of truth** | Prevents conflicts |
| **Agents can work in parallel** | Faster results |
| **If stuck, escalate to Legacy** | Clear hierarchy |

---

## Deep Dive: Heartbeat Logic

### What Happens Every 60 Seconds

**SECONDS 0-5: GATHER DATA**

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   TODOIST   │  │  CALENDAR   │  │    EMAIL    │
│ • Today: 5  │  │ • Next mtg: │  │ • Unread: 3 │
│ • Overdue:2 │  │   2pm John  │  │ • Urgent: 1 │
│ • @legacy:1 │  │ • Tomorrow: │  │ • From:     │
│             │  │   10am board│  │   investor  │
└─────────────┘  └─────────────┘  └─────────────┘

┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   GITHUB    │  │    HOME     │  │   WEATHER   │
│ • New PR: 1 │  │  ASSISTANT  │  │ • Temp: 72°F│
│ • Issues: 3 │  │ • Location: │  │ • Rain: 20% │
│             │  │   Home      │  │             │
└─────────────┘  └─────────────┘  └─────────────┘
```

**SECONDS 5-10: UPDATE MEMORY**

All gathered data → Letta "current_context" block:

```
CURRENT_CONTEXT (updated every 60 sec)

last_sync: 2026-01-15T14:32:00Z
tasks_today: 5
tasks_overdue: 2
tasks_legacy: 1 (tagged @legacy)
next_meeting: "2pm - John (investor pitch)"
meeting_in_minutes: 88
unread_emails: 3
urgent_emails: ["investor@acme.com - RE: Term sheet"]
location: "home"
weather: "72°F, 20% rain"
energy_estimate: "high"
```

**SECONDS 10-15: CHECK TRIGGERS**

| Trigger Check | Result | Action |
|--------------|--------|--------|
| Urgent email exists? | YES | → NOTIFY |
| Meeting in < 60 min? | NO (88 min) | - |
| Meeting in < 120 min? | YES | → PREP |
| Task tagged @legacy? | YES | → PROCESS |
| Overdue > 3 days? | NO | - |
| System health issues? | NO | - |

**SECONDS 15-30: EXECUTE ACTIONS**

```
ACTION 1: Notify Zeus about urgent email
📱 Telegram: "🚨 Urgent email from investor@acme.com
              Subject: RE: Term sheet"

ACTION 2: Start meeting prep (2 hours before)
Legacy → Founder Assistant: "Prepare brief for 2pm"
Legacy → Intelligence: "Research John's investments"

ACTION 3: Process @legacy task
Task: "Review Q4 financials @legacy"
→ Delegate to HP-Finance
→ Notify Zeus when ready
```

**SECONDS 30-60: SLEEP**

Wait until next cycle...

### Trigger Types

| Trigger | Condition | Action |
|---------|-----------|--------|
| **Urgent Email** | From VIP or marked urgent | Notify immediately |
| **Meeting Soon** | < 60 min away | Send reminder |
| **Meeting Prep** | < 120 min away | Start preparing |
| **@legacy Task** | Tagged @legacy | Process and act |
| **Overdue Alert** | > 3 days overdue | Suggest reschedule |
| **Location Change** | Zeus moved | Context actions |
| **System Alert** | Service down | Notify + fix |

### Learning From Patterns

```
PATTERN: "Zeus checks email first thing in morning"
ADAPTATION: Pre-fetch and categorize at 5:45am

PATTERN: "Zeus ignores notifications 6am-11am"
ADAPTATION: Queue non-urgent until 11am

PATTERN: "Zeus needs 30 min prep before investor meetings"
ADAPTATION: Start prep 2 hours before, not 1
```

---

## Deep Dive: Interface Design

### Voice Interface (Vapi)

**How a voice conversation works:**

```
YOU: *call Legacy's number*

Legacy: "Hey Zeus, what's up?"
        (Recognizes your number, casual greeting)

YOU: "What's my day look like?"

[Behind the scenes]
1. Vapi converts speech → text
2. Text sent to Legacy (Letta)
3. Legacy checks current_context (already knows!)
4. Legacy composes response
5. Response → Vapi → speech (ElevenLabs)

Legacy: "You've got five tasks today, two overdue.
         Meeting with John at 2pm - brief is ready.
         Urgent email from investor about term sheet.
         Want me to summarize it?"

YOU: "Yeah, summarize"

Legacy: "Investor asking about timeline for closing.
         Wants to meet Tuesday to finalize.
         Should I draft a response confirming?"

YOU: "Yes, Tuesday 3pm works"

Legacy: "Got it. I'll draft the email and send
         to Telegram for your approval."
```

### Text Interface (Telegram)

**Why Telegram?**
- Works everywhere (phone, tablet, desktop)
- Rich formatting (bold, links, buttons)
- Bots are first-class citizens
- Push notifications
- Media support

**Example conversation:**

```
YOU: Add task: Call John tomorrow at 10am

Legacy: ✅ Task created:
        📋 Call John
        📅 Tomorrow, 10:00 AM
        📁 Inbox
        [View in Todoist]
```

**Notification example:**

```
Legacy: 🔔 NOTIFICATION
        Meeting with John starts in 1 hour.

        📄 **Meeting Brief Ready**
        • John is partner at Acme VC
        • Last meeting: Series A terms
        • Topics: timeline, valuation, board

        [View Full Brief] [Reschedule]
```

**Approval example:**

```
Legacy: 📧 APPROVAL NEEDED

        Draft email to investor@acme.com:
        ─────────────────────────────────
        Subject: RE: Term Sheet - Tuesday

        Hi John,
        Tuesday at 3pm works perfectly.
        Looking forward to finalizing.

        Best,
        Zeus
        ─────────────────────────────────

        [✅ Send] [✏️ Edit] [❌ Cancel]
```

### Dashboard Interface (Web)

```
┌─────────────────────────────────────────────────────────────────┐
│                    ONEMIND DASHBOARD                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  TODAY'S FOCUS              │  Legacy ACTIVITY                  │
│  ─────────────              │  ───────────────                  │
│  ⏰ 6am Deep work           │  14:32 Checked email (3 new)     │
│  📋 5 tasks (2 overdue)     │  14:33 Preparing meeting brief   │
│  📅 2pm John meeting        │  14:35 Brief ready for review    │
│                             │  [View Full Logs]                 │
│                                                                  │
│  PENDING APPROVALS                                              │
│  ─────────────────                                              │
│  📧 Email to investor  [✅ Approve] [✏️ Edit] [❌ Deny]         │
│  📅 Calendar invite    [✅ Approve] [✏️ Edit] [❌ Deny]         │
│                                                                  │
│  SYSTEM        AGENTS        QUICK ACTIONS                      │
│  Letta ✅      Legacy ✅     [📝 New Task]                      │
│  Postgres ✅   HP-Intel ✅   [📧 Check Email]                   │
│  N8N ✅        GE-Ops ✅     [🔊 Call Legacy]                   │
│  Gateway ✅    IT-Infra ✅   [⏸️ Pause System]                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### The Key Principle: Meet Zeus Where He Is

| Situation | Best Interface |
|-----------|---------------|
| Hands busy (driving, cooking) | Voice (Vapi) |
| Quick check or action | Telegram |
| Deep planning session | Dashboard |
| Writing code | Claude Code |
| In Obsidian | Chat panel |

**All channels connect to the SAME Legacy with the SAME memory.**

---

---

## OneMind Architecture - Artistic Visualizations

### The Mountain of OneMind

```
                                      ·
                               *            ·
                         ·          ✦            *
                    *                                  ·
                              ╔═══════╗
                              ║ ZEUS  ║
                              ║  ◉◉   ║
                              ╚═══╤═══╝
                                  │
                                 ╱ ╲
                                ╱   ╲
                               ╱     ╲
                              ╱ Legacy╲
                             ╱   🧠    ╲
                            ╱           ╲
                           ╱─────────────╲
                          ╱       │       ╲
                         ╱        │        ╲
                        ╱    ┌────┴────┐    ╲
                       ╱     │HEARTBEAT│     ╲
                      ╱      │ ❤ 60s ❤│      ╲
                     ╱       └────┬────┘       ╲
                    ╱             │             ╲
                   ╱──────────────┼──────────────╲
                  ╱    │          │          │    ╲
                 ╱   ┌─┴─┐      ┌─┴─┐      ┌─┴─┐   ╲
                ╱    │ HP│      │ GE│      │ LE│    ╲
               ╱     │ 🧠│      │ 💰│      │ 🏠│     ╲
              ╱      └───┘      └───┘      └───┘      ╲
             ╱                   │                     ╲
            ╱────────────────────┼────────────────────╲
           ╱                     │                     ╲
          ╱      AGENTS          │       TOOLS          ╲
         ╱   ┌──┬──┬──┬──┐       │    ┌──┬──┬──┬──┐     ╲
        ╱    │01│02│03│04│       │    │📋│🐙│📧│💬│      ╲
       ╱     └──┴──┴──┴──┘       │    └──┴──┴──┴──┘       ╲
      ╱──────────────────────────┼────────────────────────╲
     ╱                           │                         ╲
    ╱         ╔══════════════════╧══════════════════╗       ╲
   ╱          ║         VIRTUAL MACHINE             ║        ╲
  ╱           ║  Ubuntu 24.04 │ Docker │ Tailscale  ║         ╲
 ╱            ╚══════════════════════════════════════╝          ╲
╱~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~╲
╲________________________________________________________________╱
                        YOUR LIFE FOUNDATION

            "One Mind. One System. Infinite Possibility."
```

### The Four Pillars Rising

```
                              ╱╲
                             ╱  ╲
                            ╱ UI ╲
                           ╱ TECH ╲
                          ╱   ⚡   ╲
                         ╱_________╲
                              │
     ╱╲                       │                       ╱╲
    ╱  ╲                      │                      ╱  ╲
   ╱ HP ╲                     │                     ╱ GE ╲
  ╱ SELF ╲════════════════════╪════════════════════╱WEALTH╲
 ╱   🧠   ╲                   │                   ╱   💰   ╲
╱__________╲                  │                  ╱__________╲
      │                       │                       │
      │                  ╱╲   │   ╱╲                  │
      │                 ╱  ╲  │  ╱  ╲                 │
      │                ╱ LE ╲ │ ╱    ╲                │
      │               ╱Legacy╲│╱      ╲               │
      │              ╱   🏠   ╲        ╲              │
      │             ╱__________╲        ╲             │
      │                   │                           │
      └───────────────────┴───────────────────────────┘
                          │
              ════════════╧════════════
                    THE FOUNDATION
              ══════════════════════════
```

### The Hive Mind (12 Agents)

```
                              ·  ✧  ·
                         ·    Legacy    ·
                        ·    ╱     ╲    ·
                       · ╱           ╲ ·
                  ·   ╱       🧠       ╲   ·
                 ╱   ╱   ORCHESTRATOR   ╲   ╲
               ╱    ╱                     ╲    ╲
             ╱    ╱     MASTER  AGENT      ╲    ╲
           ╱    ╱___________________________╲    ╲
         ╱   ╱               │               ╲   ╲
       ╱   ╱    ┌────────────┼────────────┐   ╲   ╲
      ·   │     │            │            │    │   ·
     ·    │     ▼            ▼            ▼    │    ·
    ·     │  ╭─────╮     ╭─────╮     ╭─────╮   │     ·
   ·      │  │ 01  │     │ 02  │     │ 03  │   │      ·
  ·   HP  │  │Intel│     │Perf │     │ Fin │   │  HP   ·
  ·  ◯────┼──│ 📚  │     │ 💪  │     │ 📈  │───┼────◯  ·
  ·       │  ╰─────╯     ╰─────╯     ╰─────╯   │       ·
   ·      │     │            │            │    │      ·
    ·     │     ▼            ▼            ▼    │     ·
     ·    │  ╭─────╮     ╭─────╮     ╭─────╮   │    ·
      ·   │  │ 04  │     │ 05  │     │ 06  │   │   ·
       ╲  │  │BzOps│     │Found│     │ Cap │   │  ╱
    GE  ◯─┼──│ 🏢  │     │ 🤝  │     │ 💵  │───┼─◯  GE
         ╲│  ╰─────╯     ╰─────╯     ╰─────╯   │╱
          │     │            │            │    │
          │     ▼            ▼            ▼    │
          │  ╭─────╮     ╭─────╮     ╭─────╮   │
          │  │ 07  │     │ 08  │     │ 09  │   │
     LE ◯─┼──│Home │     │Hmstd│     │Wlth │───┼─◯  LE
          │  │ 🏠  │     │ 🌾  │     │ 👨‍👩‍👧  │   │
          │  ╰─────╯     ╰─────╯     ╰─────╯   │
          │       ╲          │          ╱      │
          │        ╲    ╭────┴────╮    ╱       │
          │         ╲   │   IT    │   ╱        │
     IT ◯─┼──────────╲──│ 10 │ 11 │──╱─────────┼─◯  IT
          │           ╲ │ ⚡ │ 🤖 │ ╱          │
          ╲            ╲╰─────────╯╱           ╱
           ╲            ╲    │    ╱           ╱
            ╲            ╲   │   ╱           ╱
             ╲___________╲__│__╱___________╱
                          ╲ │ ╱
                           ╲│╱
                      ══════╧══════
                       SHARED MEMORY
                      ══════════════
                    │ Codex │ Profile │
                    │ State │ Context │
                      ══════════════

                 "12 Minds. 1 Purpose. YOUR Life."
```

### The Flow of Awareness (Day/Night Cycle)

```
    ☀️ MORNING                                           🌙 NIGHT
    ══════════                                           ══════════

         6:00 AM                                         6:00 PM
            │                                               │
            ▼                                               ▼
    ┌───────────────┐                               ┌───────────────┐
    │ Legacy WAKES  │                               │ Legacy WRAPS  │
    │    ◉    ◉     │                               │    ◉    ◉     │
    │      ‿       │                               │      ‿        │
    └───────┬───────┘                               └───────┬───────┘
            │                                               │
            ▼                                               ▼
    ╔═══════════════╗                               ╔═══════════════╗
    ║ CHECK:        ║                               ║ SUMMARIZE:    ║
    ║ • Email 📧    ║                               ║ • Completions ║
    ║ • Tasks ✓     ║                               ║ • Pending     ║
    ║ • Calendar 📅 ║                               ║ • Tomorrow    ║
    ╚═══════╤═══════╝                               ╚═══════╤═══════╝
            │                                               │
            ▼                                               ▼
    ╭───────────────╮                               ╭───────────────╮
    │"Good morning  │                               │"Here's what   │
    │ Zeus! Here's  │                               │ happened      │
    │ your day..."  │                               │ today..."     │
    ╰───────────────╯                               ╰───────────────╯
            │                                               │
            │         ┌─────────────────────┐               │
            │         │                     │               │
            └────────►│   YOUR DAY FLOWS    │◄──────────────┘
                      │                     │
                      │  ❤ ❤ ❤ ❤ ❤ ❤ ❤ ❤   │
                      │   (Heartbeat runs)  │
                      │                     │
                      │  ▸ Monitors         │
                      │  ▸ Coordinates      │
                      │  ▸ Alerts           │
                      │  ▸ Prepares         │
                      │                     │
                      └─────────────────────┘
```

### The Complete Landscape

```
                                    ·  ✦
                         *    ·          ·    *
                    ·              ·              ·
        ✦                    *                         ✦

                         ┌─────────┐
                         │  ZEUS   │
                         │    ◉    │
                         └────┬────┘
                              │
    ════════════════════════════════════════════════════════
                         ∿ CLOUD ∿
    ════════════════════════════════════════════════════════
                              │
                    ╔═════════╧═════════╗
                    ║                   ║
                    ║   L E G A C Y     ║
                    ║   ┌───────────┐   ║
                    ║   │ 🧠 BRAIN  │   ║
                    ║   │  (Letta)  │   ║
                    ║   └───────────┘   ║
                    ║                   ║
                    ╚═════════╤═════════╝
                              │
             ┌────────────────┼────────────────┐
             │                │                │
             ▼                ▼                ▼
      ┌────────────┐   ┌────────────┐   ┌────────────┐
      │INTELLIGENCE│   │ PERFORMANCE│   │  FINANCE   │
      │     📚     │   │     💪     │   │     📈     │
      └─────┬──────┘   └─────┬──────┘   └─────┬──────┘
            │                │                │
            └────────────────┴────────────────┘
                             │
    ════════════════════════════════════════════════════════
                    ∿ THE VIRTUAL MACHINE ∿
    ════════════════════════════════════════════════════════
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
        ▼                    ▼                    ▼
   ╔═════════╗          ╔═════════╗          ╔═════════╗
   ║POSTGRES ║          ║  REDIS  ║          ║   N8N   ║
   ║   📦    ║          ║   ⚡    ║          ║   🔄    ║
   ║ Memory  ║          ║  Cache  ║          ║Workflows║
   ╚════╤════╝          ╚════╤════╝          ╚════╤════╝
        │                    │                    │
        └────────────────────┴────────────────────┘
                             │
    ════════════════════════════════════════════════════════
                    ∿ API GATEWAY ∿
              ┌──────────────┴──────────────┐
    ════════════════════════════════════════════════════════
              │              │              │
              ▼              ▼              ▼
         ┌────────┐    ┌────────┐    ┌────────┐
         │TODOIST │    │ GITHUB │    │  ZOHO  │
         │   ✓    │    │   🐙   │    │   📧   │
         └────────┘    └────────┘    └────────┘
              │              │              │
              ▼              ▼              ▼
         ┌────────┐    ┌────────┐    ┌────────┐
         │DISCORD │    │TELEGRAM│    │  HOME  │
         │   💬   │    │   ✈    │    │ASST 🏠 │
         └────────┘    └────────┘    └────────┘

    ════════════════════════════════════════════════════════
                       Y O U R    L I F E
    ════════════════════════════════════════════════════════
```

### The OneMind Cosmos

```
                                        ✧
                              ✦                    ✧
                    ✧                                        ✦
          ✦                         ✧
                          ╔═══════════════════╗
                   ✧      ║    Z E U S        ║      ✧
                          ║   ◉  Human  ◉     ║
                          ╚════════╤══════════╝
                                   │
                    ═══════════════╪═══════════════
                         ∿∿∿ VOICE ∿∿∿ TEXT ∿∿∿
                    ═══════════════╪═══════════════
                                   │
                                   ▼
                    ╭─────────────────────────────╮
                   ╱    ██╗     ███████╗ ██████╗   ╲
                  ╱     ██║     ██╔════╝██╔════╝    ╲
                 ║      ██║     █████╗  ██║  ███╗    ║
                 ║      ██║     ██╔══╝  ██║   ██║    ║
                  ╲     ███████╗███████╗╚██████╔╝   ╱
                   ╲    ╚══════╝╚══════╝ ╚═════╝   ╱
                    ╰──────────────┬──────────────╯
                       THE MASTER ORCHESTRATOR
                                   │
                    ╔══════════════╧══════════════╗
                    ║    ♥ HEARTBEAT (60 sec) ♥   ║
                    ║  ·  ·  ·  ❤  ·  ·  ·  ❤  · ║
                    ╚══════════════╤══════════════╝
                                   │
            ┌──────────────────────┼──────────────────────┐
            │                      │                      │
            ▼                      ▼                      ▼

                    "HUMAN + AI = ONE MIND"
```

---

## Related Documents

- [Legacy-OPERATIONS-CENTER.md](./Legacy-OPERATIONS-CENTER.md) - Original ops center plan
- [ONEMIND-CODEX.md](./ONEMIND-CODEX.md) - Master OneMind blueprint
- [00 - Agent Roster Overview.md](./00-24%20UI%20(Unified%20Intelligence)/02%20Agents%20(AI%2C%20Memory%20%26%20Tools)/02.2%20AI-Agent-Frameworks/00%20-%20Agent%20Roster%20Overview.md) - 12-agent architecture
- [00 - Legacy - Master Orchestrator.md](./00-24%20UI%20(Unified%20Intelligence)/02%20Agents%20(AI%2C%20Memory%20%26%20Tools)/02.5%20Legacy-AI-Development/00%20-%20Legacy%20-%20Master%20Orchestrator.md) - Legacy's full specification

---

**Document Status:** Living document - will be updated as architecture evolves

*Last Updated: January 15, 2026*
