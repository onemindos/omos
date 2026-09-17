# OneMind OS Product Comparison

> **An objective comparison of AI-powered enhancement framework and second brain software products**
>
> Comparing OneMind OS against market competitors

---

## Overview Matrix

### Software / Enhancement Framework Products

| Product           | Company       | Focus                 | AI Integration         | Open Source | Self-Hosted | Price             |
| ----------------- | ------------- | --------------------- | ---------------------- | ----------- | ----------- | ----------------- |
| **OneMind OS**    | Zeus Delacruz | Enhancement Framework + AI Agents   | Native (Claude, Letta) | Yes         | Yes         | Free              |
| **Notion AI**     | Notion Labs   | Workspace + AI        | Add-on                 | No          | No          | $10+/mo           |
| **Obsidian**      | Obsidian      | PKM + Plugins         | Via plugins            | No (core)   | Yes (local) | Free / $50/yr     |
| **Mem.ai**        | Mem           | AI-first notes        | Native                 | No          | No          | $15+/mo           |
| **Reflect**       | Reflect Notes | Networked notes + AI  | Native                 | No          | No          | $10/mo            |
| **Capacities**    | Capacities    | Object-based notes    | Add-on                 | No          | No          | Free / $13/mo     |
| **Tana**          | Tana          | Supertags + AI        | Native                 | No          | No          | Waitlist / $10/mo |
| **Roam Research** | Roam Research | Networked thought     | Limited                | No          | No          | $15/mo            |
| **Logseq**        | Logseq        | Outliner + PKM        | Via plugins            | Yes         | Yes         | Free              |
| **Anytype**       | Anytype       | Local-first workspace | Limited                | Yes         | Yes         | Free              |
| **Personal AI**   | Personal.ai   | AI memory companion   | Native                 | No          | No          | $40+/mo           |
| **Motion**        | Motion        | AI Calendar + Tasks   | Native                 | No          | No          | $19-34/mo         |
| **Benji**         | Benji         | Enhancement Framework / Habits      | Limited                | No          | No          | $8-10/mo          |

### Smart Home / Physical World

| Product             | Company         | Focus                  | AI Integration | Open Source | Price     |
| ------------------- | --------------- | ---------------------- | -------------- | ----------- | --------- |
| **OneMind OS**      | Zeus Delacruz   | Enhancement Framework + IoT/Robotics | Native agents  | Yes         | Free      |
| **Josh.ai**         | Josh.ai Inc     | Luxury smart home      | JoshGPT        | No          | $10,000+  |
| **1X Technologies** | 1X Technologies | Humanoid robots (NEO)  | Proprietary    | No          | TBD ($$$) |

### Agent Frameworks / Infrastructure

| Product        | Company       | Focus                 | Letta Compatible | Open Source |
| -------------- | ------------- | --------------------- | ---------------- | ----------- |
| **OneMind OS** | Zeus Delacruz | Enhancement Framework + Agents      | Yes              | Yes         |
| **Animus**     | Animus        | Embodied AI framework | Yes              | Yes         |

### Philosophy / Education (Non-Software)

| Resource            | Company         | Focus                   | Relationship to OneMind   |
| ------------------- | --------------- | ----------------------- | ------------------------- |
| **Acorn Land Labs** | Acorn Land Labs | Solar-punk homesteading | Philosophical inspiration |

---

## Notion AI

**Company:** Notion Labs (San Francisco)
**Founded:** 2013 | **AI Launch:** 2023
**Users:** 30M+

### What It Does

All-in-one workspace combining notes, databases, wikis, and project management with AI assistance.

### AI Capabilities

| Feature     | Description                                |
| ----------- | ------------------------------------------ |
| AI Writing  | Drafting, summarizing, editing text        |
| Q&A         | Ask questions about your workspace content |
| Autofill    | AI-generated database properties           |
| Translation | Multi-language support                     |

### Strengths ✅

- **Massive ecosystem** — templates, integrations, community
- **Familiar interface** — easy adoption for teams
- **Database power** — relational data with views
- **Collaboration** — real-time multi-user editing
- **Cross-platform** — web, desktop, mobile apps
- **Enterprise-ready** — SOC2, SSO, permissions

### Weaknesses ❌

- **Cloud-dependent** — no offline-first, no self-hosting
- **AI is shallow** — no persistent memory, no agents
- **Vendor lock-in** — data export is limited
- **Performance** — can be slow with large workspaces
- **No local files** — everything in Notion's cloud
- **Privacy concerns** — data on Notion's servers

### vs OneMind OS

| Dimension      | Notion AI                | OneMind OS                         |
| -------------- | ------------------------ | ---------------------------------- |
| AI depth       | Surface-level assistant  | Deep agent integration with memory |
| Data ownership | Notion's cloud           | Self-hosted, local-first option    |
| Automation     | Zapier/Make integrations | Native n8n, NATS, workflows        |
| Life scope     | Work-focused             | Full life (HP, LE, GE, UI)         |
| Customization  | Templates only           | Full code access                   |
| Price          | $10+/mo per user         | Free (self-hosted)                 |

---

## Obsidian

**Company:** Obsidian (Canada)
**Founded:** 2020
**Users:** 1M+

### What It Does

Local-first markdown editor with bidirectional linking, plugins, and graph visualization.

### AI Capabilities (via plugins)

| Plugin            | Function                  |
| ----------------- | ------------------------- |
| Copilot           | ChatGPT/Claude in sidebar |
| Smart Connections | AI-powered note linking   |
| Text Generator    | GPT-powered writing       |
| Khoj              | Self-hosted AI search     |

### Strengths ✅

- **Local-first** — your files, your computer, plain markdown
- **Plugin ecosystem** — 1000+ community plugins
- **Graph view** — visualize note connections
- **Privacy** — no cloud required (Sync is optional)
- **Customizable** — CSS themes, hotkeys, workflows
- **One-time cost** — no subscription for core features
- **Zettelkasten-friendly** — built for linked notes

### Weaknesses ❌

- **No native AI** — relies on plugins (varying quality)
- **Plugin maintenance** — community-dependent, can break
- **Steep learning curve** — power requires configuration
- **No native tasks** — task management via plugins only
- **No database views** — markdown tables, not Notion-style DBs
- **Mobile experience** — functional but not optimal
- **No real-time collab** — single-user focused

### vs OneMind OS

| Dimension       | Obsidian                      | OneMind OS                     |
| --------------- | ----------------------------- | ------------------------------ |
| AI integration  | Plugin-based, fragmented      | Native agent architecture      |
| Automation      | Limited (Templater, Dataview) | Full workflow automation (n8n) |
| Task management | Plugins (Tasks, Kanban)       | Integrated (Todoist, native)   |
| Agent memory    | None                          | Persistent (Letta)             |
| IoT/Robotics    | None                          | Native support                 |
| Philosophy      | Tool for thought              | Human enhancement framework          |

**Note:** OneMind OS actually uses Obsidian as the PKM layer — they're complementary, not competitors.

---

## Mem.ai

**Company:** Mem (San Francisco)
**Founded:** 2020
**Funding:** $29M

### What It Does

AI-first note-taking with automatic organization, smart search, and writing assistance.

### AI Capabilities

| Feature           | Description                                |
| ----------------- | ------------------------------------------ |
| Mem Chat          | Ask questions, get answers from your notes |
| Smart Write       | AI drafting and editing                    |
| Auto-organization | AI suggests collections and tags           |
| Related memories  | Surface relevant past notes                |

### Strengths ✅

- **AI-native design** — built for AI from day one
- **Frictionless capture** — quick notes without structure
- **Temporal context** — timeline-based organization
- **Smart search** — semantic, not just keyword
- **Clean UX** — minimalist, focused interface
- **Meeting integration** — calendar and transcription

### Weaknesses ❌

- **Cloud-only** — no local storage, no self-hosting
- **No agent autonomy** — AI assists but doesn't act
- **Limited structure** — weak on projects/areas/databases
- **Expensive** — $15-25/month
- **Vendor lock-in** — proprietary format
- **No automation** — no workflow capabilities
- **Individual focus** — weak team features

### vs OneMind OS

| Dimension      | Mem.ai              | OneMind OS                     |
| -------------- | ------------------- | ------------------------------ |
| AI autonomy    | Assistant only      | Agents with autonomy levels    |
| Data ownership | Mem's cloud         | Self-hosted, local files       |
| Structure      | Flat / AI-organized | Explicit 4-domain architecture |
| Automation     | None                | Full workflow orchestration    |
| Life scope     | Notes/knowledge     | Full enhancement framework                   |
| Price          | $15-25/mo           | Free                           |

---

## Reflect

**Company:** Reflect Notes
**Founded:** 2021

### What It Does

Networked note-taking with AI assistant, backlinks, and daily notes focus.

### AI Capabilities

| Feature       | Description                         |
| ------------- | ----------------------------------- |
| AI Assistant  | GPT-4 powered chat about your notes |
| Transcription | Voice notes to text                 |
| Smart linking | AI-suggested connections            |

### Strengths ✅

- **Beautiful design** — polished, minimal interface
- **End-to-end encryption** — privacy-focused
- **Daily notes** — journaling-first workflow
- **Fast sync** — real-time across devices
- **Voice notes** — native transcription
- **iOS excellence** — best-in-class mobile app

### Weaknesses ❌

- **Limited AI depth** — basic assistant, no agents
- **No self-hosting** — cloud-only
- **Weak structure** — no areas/projects/databases
- **No automation** — no workflows
- **Individual only** — no team features
- **Limited integrations** — small ecosystem

### vs OneMind OS

| Dimension     | Reflect            | OneMind OS               |
| ------------- | ------------------ | ------------------------ |
| AI depth      | Chat assistant     | Multi-agent architecture |
| Data location | Cloud (encrypted)  | Self-hosted              |
| Structure     | Flat + daily notes | 4-domain enhancement framework         |
| Automation    | None               | Native workflows         |
| Price         | $10/mo             | Free                     |

---

## Capacities

**Company:** Capacities (Germany)
**Founded:** 2021

### What It Does

Object-based note-taking where everything is a typed object (person, book, project, etc.).

### AI Capabilities

| Feature      | Description                 |
| ------------ | --------------------------- |
| AI Assistant | Chat about your content     |
| Smart blocks | AI-generated content blocks |

### Strengths ✅

- **Object-based** — structured data model
- **Daily notes** — strong journaling support
- **Media handling** — images, PDFs, web clips
- **Type system** — custom object types
- **Graph view** — relationship visualization
- **Free tier** — generous free plan

### Weaknesses ❌

- **No self-hosting** — cloud-only
- **Limited AI** — basic assistant
- **No automation** — manual workflows only
- **Young product** — still evolving
- **No agents** — no autonomous AI

### vs OneMind OS

| Dimension    | Capacities      | OneMind OS                 |
| ------------ | --------------- | -------------------------- |
| Data model   | Objects/types   | Domains/folders + entities |
| AI           | Basic assistant | Agent architecture         |
| Automation   | None            | Full workflow support      |
| Self-hosting | No              | Yes                        |

---

## Tana

**Company:** Tana (Norway)
**Founded:** 2020

### What It Does

Outliner with "supertags" that add structured fields to any node, plus AI integration.

### AI Capabilities

| Feature      | Description                              |
| ------------ | ---------------------------------------- |
| Tana AI      | GPT-4 integration for content generation |
| AI Commands  | Custom AI workflows on nodes             |
| Smart fields | AI-populated supertag fields             |

### Strengths ✅

- **Supertags** — powerful structured data in an outliner
- **AI commands** — customizable AI workflows
- **Outliner power** — infinitely nested, flexible
- **Queries** — Roam-like live queries
- **Innovative** — pushing PKM boundaries

### Weaknesses ❌

- **Steep learning curve** — complex concepts
- **Still invite-only** — limited access
- **No offline** — requires internet
- **No self-hosting** — cloud-only
- **Performance** — can lag with large workspaces
- **No mobile app** — web only currently

### vs OneMind OS

| Dimension    | Tana                  | OneMind OS                |
| ------------ | --------------------- | ------------------------- |
| Structure    | Supertags/outliner    | Domain folders + entities |
| AI           | Commands + generation | Full agent architecture   |
| Offline      | No                    | Yes                       |
| Self-hosting | No                    | Yes                       |
| Mobile       | No                    | Planned                   |

---

## Logseq

**Company:** Logseq (Open source)
**Founded:** 2020

### What It Does

Open-source outliner with bidirectional links, local-first storage, and plugin system.

### AI Capabilities

| Plugin         | Function              |
| -------------- | --------------------- |
| GPT-3 Logseq   | OpenAI integration    |
| Logseq Copilot | AI writing assistance |

### Strengths ✅

- **Open source** — fully transparent, community-driven
- **Local-first** — your files, markdown/org-mode
- **Free** — no cost
- **Plugin system** — extensible
- **Privacy** — no cloud required
- **Outliner + pages** — hybrid approach

### Weaknesses ❌

- **Performance issues** — slow with large graphs
- **Mobile lag** — behind Obsidian
- **AI is plugin-only** — not native
- **Smaller ecosystem** — fewer plugins than Obsidian
- **No automation** — limited workflow support
- **Development pace** — slower than commercial tools

### vs OneMind OS

| Dimension       | Logseq            | OneMind OS           |
| --------------- | ----------------- | -------------------- |
| AI              | Plugin-based      | Native agents        |
| Automation      | Limited           | Full n8n integration |
| Task management | Built-in (basic)  | Todoist + native     |
| Structure       | Outliner/journals | Domain architecture  |
| Open source     | Yes               | Yes                  |

---

## Anytype

**Company:** Anytype (Berlin)
**Founded:** 2019

### What It Does

Local-first, encrypted workspace with objects, types, and peer-to-peer sync.

### AI Capabilities

Limited — focus is on local-first and privacy, not AI.

### Strengths ✅

- **Local-first** — data on your device
- **End-to-end encrypted** — privacy by design
- **P2P sync** — no central server required
- **Open source** — transparent codebase
- **Object model** — structured data
- **Free** — no subscription

### Weaknesses ❌

- **Minimal AI** — not an AI-first product
- **Early stage** — still developing features
- **Sync issues** — P2P can be unreliable
- **Limited ecosystem** — few integrations
- **Learning curve** — unique concepts

### vs OneMind OS

| Dimension      | Anytype       | OneMind OS     |
| -------------- | ------------- | -------------- |
| AI integration | Minimal       | Native agents  |
| Privacy        | Encrypted P2P | Self-hosted    |
| Automation     | None          | Full workflows |
| Maturity       | Early         | Early          |

---

## Personal AI

**Company:** Personal.ai (US)
**Founded:** 2021
**Funding:** $10.5M

### What It Does

AI that learns from your data to become a personalized assistant with persistent memory.

### AI Capabilities

| Feature        | Description                        |
| -------------- | ---------------------------------- |
| Memory stacks  | Train AI on your content           |
| Personal model | AI learns your style and knowledge |
| Lounge         | Chat with your AI twin             |
| Integrations   | Slack, email, calendar sync        |

### Strengths ✅

- **True personalization** — AI learns YOU specifically
- **Memory persistence** — remembers across sessions
- **Multi-source** — learns from various inputs
- **AI twin concept** — digital you
- **Business use case** — clone yourself for customers

### Weaknesses ❌

- **Expensive** — $40+/month
- **Cloud-only** — no self-hosting
- **Not a PKM** — AI companion, not knowledge system
- **Limited structure** — memory, not organization
- **No automation** — can't take actions
- **Proprietary** — no data portability

### vs OneMind OS

| Dimension         | Personal AI       | OneMind OS          |
| ----------------- | ----------------- | ------------------- |
| AI memory         | Yes, personalized | Yes (Letta)         |
| Action capability | Talk only         | Full agent autonomy |
| PKM               | No                | Yes (Obsidian)      |
| Automation        | No                | Yes (n8n)           |
| Self-hosting      | No                | Yes                 |
| Price             | $40+/mo           | Free                |

---

## Josh.ai

**Company:** Josh.ai Inc. (Denver, CO)
**Founded:** 2015
**Focus:** Luxury smart home voice control

### What It Does

Privacy-focused AI voice assistant for high-end smart home control. Works with Crestron, Lutron, Sonos, and other luxury home automation systems.

### AI Capabilities

| Feature           | Description                           |
| ----------------- | ------------------------------------- |
| JoshGPT           | GPT-powered conversational interface  |
| Natural language  | Complex multi-device commands         |
| Context awareness | Understands room, time, preferences   |
| Proactive alerts  | Security and automation notifications |

### Strengths ✅

- **Privacy-first** — data never sold for ads or marketing
- **Luxury integration** — works with high-end home systems
- **Natural voice** — "OK Josh, dim the lights and play jazz"
- **Beautiful hardware** — Josh Micro, Nano, Edge devices
- **Professional install** — certified integrator network
- **Continuous updates** — OTA improvements

### Weaknesses ❌

- **Expensive** — requires professional installation ($10K+)
- **Luxury only** — not for average consumers
- **Closed ecosystem** — proprietary system
- **No PKM** — smart home only, not life management
- **No task management** — voice control, not productivity

### vs OneMind OS

| Dimension   | Josh.ai          | OneMind OS             |
| ----------- | ---------------- | ---------------------- |
| Smart home  | ✅ Best-in-class | ✅ Via Home Assistant  |
| PKM         | ❌ None          | ✅ Full system         |
| AI agents   | ⚠️ Voice only    | ✅ Multi-modal agents  |
| Self-hosted | ❌ Proprietary   | ✅ Yes                 |
| Price       | $10,000+         | Free                   |
| Target      | Luxury homes     | Tech-savvy individuals |

---

## Motion (usemotion.com)

**Company:** Motion (Y Combinator)
**Founded:** 2019
**Users:** 1M+
**Funding:** Y Combinator backed

### What It Does

AI-powered productivity platform combining calendar, tasks, projects, docs, and meeting notes with automatic scheduling and prioritization.

### AI Capabilities

| Feature            | Description                                   |
| ------------------ | --------------------------------------------- |
| AI Task Planner    | Auto-prioritizes based on deadlines, duration |
| AI Project Manager | Generates projects from descriptions          |
| AI Calendar        | Schedules and reschedules automatically       |
| AI Notetaker       | Meeting transcription and follow-ups          |
| AI Docs            | Drafting, editing, organizing                 |
| AI Workflows       | Turns SOPs into action items                  |

### Strengths ✅

- **All-in-one** — replaces 5-10 separate tools
- **AI scheduling** — auto-reschedules when things change
- **Project generation** — create entire projects from prompts
- **Mobile + desktop** — cross-platform apps
- **Team features** — collaboration, enterprise ready
- **Polished UX** — consumer-grade design
- **SOC2 compliant** — enterprise security

### Weaknesses ❌

- **Cloud-only** — no self-hosting option
- **No PKM** — tasks/calendar only, not knowledge
- **Subscription required** — $19-34/month
- **No AI agents** — AI assists but doesn't act autonomously
- **No memory persistence** — context resets
- **Work-focused** — not an enhancement framework (no personal domains)
- **No IoT/robotics** — digital productivity only

### vs OneMind OS

| Dimension        | Motion             | OneMind OS           |
| ---------------- | ------------------ | -------------------- |
| Task/Calendar AI | ✅ Best-in-class   | ⚠️ Via Todoist       |
| PKM              | ❌ None            | ✅ Full system       |
| AI agents        | ❌ Assistants only | ✅ Autonomous agents |
| Life scope       | Work only          | Full life (HP/LE/GE) |
| Self-hosted      | ❌ No              | ✅ Yes               |
| Mobile app       | ✅ Polished        | 🚧 In progress       |
| Price            | $19-34/mo          | Free                 |
| UX polish        | ✅ Consumer-ready  | ⚠️ Developer tool    |

---

## Benji

**Company:** Benji (benji.so)
**Focus:** "The Life OS" — habits, planner, metrics

### What It Does

Life management app combining habits, todos, planner, timeline, and personal metrics tracking across all platforms.

### AI Capabilities

Limited — focus is on tracking and planning, not AI assistance.

### Strengths ✅

- **Life OS branding** — directly comparable positioning
- **Cross-platform** — iOS, Android, macOS, Windows, Linux, Web
- **Timeline view** — see how you've changed over time
- **Habit tracking** — streaks, analytics, reminders
- **Sharing** — collaborate with family on lists and plans
- **Community** — accountability features
- **Affordable** — $8-10/month

### Weaknesses ❌

- **Minimal AI** — not AI-native
- **No PKM** — habits/tasks only, not knowledge
- **Cloud-only** — no self-hosting
- **No agents** — passive tracking tool
- **No automation** — manual workflows
- **Limited integrations** — small ecosystem
- **No IoT** — digital only

### vs OneMind OS

| Dimension      | Benji            | OneMind OS          |
| -------------- | ---------------- | ------------------- |
| Enhancement claim  | ✅ Yes           | ✅ Yes              |
| AI agents      | ❌ None          | ✅ Native           |
| PKM            | ❌ None          | ✅ Full system      |
| Habit tracking | ✅ Core feature  | ⚠️ Via integrations |
| Mobile app     | ✅ All platforms | 🚧 In progress      |
| Self-hosted    | ❌ No            | ✅ Yes              |
| Price          | $8-10/mo         | Free                |

---

## Animus

**Company:** Animus (animus.uno)
**Focus:** Open-source agentic framework for embodied AI

### What It Does

Open-source framework for creating AI agents with persistent identity ("souls") that can transfer across different physical forms (robots, avatars, devices).

### Core Architecture

| Component    | Description                         |
| ------------ | ----------------------------------- |
| Thalamus     | Sensory processing hub              |
| Cochlea      | Audio input processing              |
| Letta Kernel | Core infrastructure (uses Letta!)   |
| SMCP         | Internal MCP server                 |
| Soul State   | Persistent identity across hardware |

### Supported LLMs

GPT-4o, Claude 3.5 Sonnet, Gemini Pro, Llama 3.1, Mistral Large, Qwen 2.5, DeepSeek Coder, Phi 3.5, Codestral, o1-preview, Grok

### Strengths ✅

- **Open source** — fully transparent, community-driven
- **Embodied AI** — bridging digital and physical
- **Uses Letta** — same memory infrastructure as OneMind
- **Hardware agnostic** — same "soul" across devices
- **Modular** — brain-inspired architecture
- **Self-sovereign** — agents own their identity

### Weaknesses ❌

- **Robotics-focused** — not a life management system
- **No PKM** — agent framework, not knowledge system
- **No task management** — infrastructure, not application
- **Early stage** — limited documentation
- **Developer-only** — requires coding skills
- **Crypto-adjacent** — token/contract elements

### vs OneMind OS

| Dimension       | Animus          | OneMind OS              |
| --------------- | --------------- | ----------------------- |
| Agent framework | ✅ Core focus   | ✅ Uses Agno/Letta      |
| Embodied AI     | ✅ Primary goal | ✅ IoT/robotics planned |
| PKM             | ❌ None         | ✅ Full system          |
| Task management | ❌ None         | ✅ Todoist + native     |
| Uses Letta      | ✅ Yes          | ✅ Yes                  |
| Open source     | ✅ Yes          | ✅ Yes                  |

**Note:** Animus and OneMind share similar visions (embodied AI, Letta, open source). They're more complementary than competitive — Animus focuses on agent infrastructure, OneMind focuses on enhancement framework.

---

## 1X Technologies

**Company:** 1X Technologies (Norway)
**Founded:** 2014
**Funding:** $125M+
**Product:** NEO humanoid robot

### What It Does

Builds humanoid robots (NEO, EVE) designed for homes and workplaces, with advanced AI and physical capabilities.

### Products

| Robot   | Description                                       |
| ------- | ------------------------------------------------- |
| **NEO** | Bipedal humanoid for home use ($200 deposit open) |
| **EVE** | Wheeled android for commercial/industrial         |

### Strengths ✅

- **World-class robotics** — one of the leaders in humanoids
- **Home-focused** — NEO designed for household tasks
- **Significant funding** — $125M+ including OpenAI investment
- **Nordic engineering** — high quality hardware
- **Consumer timeline** — targeting home deployment

### Weaknesses ❌

- **Hardware company** — robots, not software platform
- **No enhancement framework** — no PKM, tasks, or knowledge system
- **Not available yet** — NEO still in development
- **Expensive** — consumer robots will cost thousands
- **No AI agent framework** — proprietary system
- **Single product** — robots only, not productivity suite

### vs OneMind OS

| Dimension         | 1X Technologies  | OneMind OS                |
| ----------------- | ---------------- | ------------------------- |
| Robotics          | ✅ Core focus    | ✅ Integration layer      |
| Enhancement Framework | ❌ None          | ✅ Full system            |
| Hardware          | ✅ Builds robots | ❌ Integrates with robots |
| Software platform | ❌ Proprietary   | ✅ Open source            |
| Available now     | ⚠️ Coming soon   | ✅ Yes                    |
| Price             | $10,000+?        | Free                      |

**Note:** 1X builds the robots; OneMind could potentially orchestrate them. They're complementary — 1X for hardware, OneMind for the "brain" layer.

---

## Acorn Land Labs

**Company:** Acorn Land Labs (US)
**Focus:** Solar-punk, off-grid living, sustainable homesteading

### What It Does

Educational platform teaching self-sufficiency in the 5 basic needs: Food, Water, Energy, Shelter, Sanitation. Offers courses, simulator app, and community.

### Products

| Offering        | Description                            |
| --------------- | -------------------------------------- |
| Land Lab SIM    | Off-grid simulation app                |
| Off-Grid Course | Comprehensive homesteading training    |
| Community       | Support network for sustainable living |
| Book            | "The Acorn Home" guide                 |

### Strengths ✅

- **Solar-punk philosophy** — same foundational values as OneMind
- **Practical skills** — food, water, energy, shelter, sanitation
- **Education-focused** — courses and community
- **Proven system** — working land lab demonstration
- **Sustainable tech** — methane, solar, rainwater, composting

### Weaknesses ❌

- **Not a software platform** — education, not OS
- **No AI** — human knowledge transfer only
- **No digital tools** — physical homesteading focus
- **Niche audience** — off-grid enthusiasts
- **No automation** — manual sustainable living

### vs OneMind OS

| Dimension         | Acorn Land Labs    | OneMind OS        |
| ----------------- | ------------------ | ----------------- |
| Solar-punk values | ✅ Core philosophy | ✅ Foundational   |
| Digital platform  | ❌ Education only  | ✅ Full OS        |
| AI agents         | ❌ None            | ✅ Native         |
| Physical world    | ✅ Homesteading    | ✅ IoT/robotics   |
| Self-sufficiency  | ✅ Core goal       | ✅ Part of vision |

**Note:** Acorn Land Labs is a philosophical sibling, not a software competitor. OneMind acknowledges them as inspiration in its solar-punk values. They teach the "why" and physical skills; OneMind provides the digital coordination layer.

---

## Feature Comparison Matrix

| Feature                 | OneMind | Notion | Obsidian | Mem | Reflect | Tana | Logseq | Anytype | Personal AI |
| ----------------------- | ------- | ------ | -------- | --- | ------- | ---- | ------ | ------- | ----------- |
| **AI Agents**           | ✅      | ❌     | ❌       | ❌  | ❌      | ⚠️   | ❌     | ❌      | ❌          |
| **Persistent Memory**   | ✅      | ❌     | ❌       | ⚠️  | ❌      | ❌   | ❌     | ❌      | ✅          |
| **Autonomy Levels**     | ✅      | ❌     | ❌       | ❌  | ❌      | ❌   | ❌     | ❌      | ❌          |
| **Workflow Automation** | ✅      | ⚠️     | ⚠️       | ❌  | ❌      | ❌   | ❌     | ❌      | ❌          |
| **Self-Hosted**         | ✅      | ❌     | ✅       | ❌  | ❌      | ❌   | ✅     | ✅      | ❌          |
| **Open Source**         | ✅      | ❌     | ❌       | ❌  | ❌      | ❌   | ✅     | ✅      | ❌          |
| **Local-First**         | ✅      | ❌     | ✅       | ❌  | ❌      | ❌   | ✅     | ✅      | ❌          |
| **Mobile App**          | 🚧      | ✅     | ✅       | ✅  | ✅      | ❌   | ✅     | ✅      | ✅          |
| **Team Collab**         | 🚧      | ✅     | ⚠️       | ❌  | ❌      | ⚠️   | ❌     | ⚠️      | ❌          |
| **IoT/Robotics**        | ✅      | ❌     | ❌       | ❌  | ❌      | ❌   | ❌     | ❌      | ❌          |
| **Free Tier**           | ✅      | ⚠️     | ✅       | ❌  | ❌      | ❌   | ✅     | ✅      | ❌          |

Legend: ✅ Yes | ⚠️ Partial | ❌ No | 🚧 In Progress

---

## Market Positioning

### The Landscape

```
                    HIGH AI DEPTH
                         │
         Personal AI     │     OneMind OS
              •          │          •
                         │
    ─────────────────────┼─────────────────────
    CLOUD-ONLY           │           SELF-HOSTED
                         │
        Mem.ai  •        │        • Logseq
        Notion  •        │        • Obsidian
        Reflect •        │        • Anytype
        Tana    •        │
                         │
                    LOW AI DEPTH
```

### Competitive Moats

| Competitor      | Their Moat                  | OneMind Counter               |
| --------------- | --------------------------- | ----------------------------- |
| **Notion**      | Network effects, enterprise | Open source, no lock-in       |
| **Obsidian**    | Plugin ecosystem, community | Complementary (uses Obsidian) |
| **Mem.ai**      | AI-first UX, simplicity     | Deeper AI with agents         |
| **Personal AI** | True personalization        | Similar + action capability   |
| **Logseq**      | Open source community       | Similar + better AI           |

---

## Honest Assessment of OneMind OS

### What OneMind Does Better Than All Competitors

1. **Agent Architecture** — Only system with true AI agents that have autonomy levels
2. **Memory Persistence** — Letta integration provides genuine long-term AI memory
3. **Action Capability** — Agents can execute, not just advise
4. **Life Scope** — Covers all domains (self, family, business, tech)
5. **IoT/Robotics** — No competitor touches physical world integration
6. **Self-Hosted + Open** — Full data sovereignty

### What OneMind Does Worse Than Competitors

1. **Maturity** — < 1 year old vs. Notion's 10+ years
2. **Mobile** — No native mobile app (competitors have polished apps)
3. **Team/Collab** — Individual-focused (Notion excels here)
4. **UX Polish** — Developer tool, not consumer-ready
5. **Community** — Tiny user base vs. millions for competitors
6. **Documentation** — Complex, requires technical understanding
7. **Onboarding** — Steep learning curve, no guided setup

### Where OneMind Ties

1. **Knowledge management** — Obsidian layer is comparable
2. **Automation** — N8N matches Notion integrations differently
3. **Privacy** — Matches Anytype, better than cloud-only options

---

## Recommendations

### Choose Notion If...

- You need team collaboration
- You want polished, consumer-ready UX
- You're okay with cloud dependency
- You need enterprise features (SSO, permissions)

### Choose Obsidian If...

- You want local-first with plugins
- You're building a Zettelkasten
- You prefer markdown portability
- You'll add AI via plugins

### Choose Mem.ai If...

- You want frictionless AI note-taking
- You prefer AI-organized over manual structure
- You're okay paying for convenience
- You don't need automation

### Choose OneMind OS If...

- You want AI agents, not just assistants
- You require self-hosting and data ownership
- You're technical and can configure systems
- You want unified enhancement framework (not just notes)
- You think in generational terms
- You accept early-adopter tradeoffs

---

## Conclusion

OneMind OS occupies a unique position: **the only self-hosted, AI-agent-native human enhancement framework.**

However, this comes with honest tradeoffs:

- It's new and unproven at scale
- It requires technical sophistication
- Mobile and UX are behind commercial competitors
- The community is tiny

**For most users:** Notion or Obsidian remain safer choices.
**For AI-forward, technical users:** OneMind offers capabilities no competitor matches.

The question is not "which is best" but "which fits your needs, skills, and risk tolerance."

---

_Document created: January 2026_
_Product comparison for strategic positioning_


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
