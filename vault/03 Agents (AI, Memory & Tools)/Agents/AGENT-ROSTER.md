# OneMind OS — Agent Roster

> **Last Updated**: 2026-03-09
> **Version**: 3.0 — 9-agent architecture (codex-pillar aligned)
> **Gateway**: OpenClaw 2026.3.7

---

## Architecture Overview

```
Zeus (human)
  ↕ Telegram / iMessage / Voice / iOS App / Discord
Legacy ⚡ (Sovereign Orchestrator)
  ├── Spartan ⚔️   ← Body & Physical Mastery (HP 27-28)
  ├── Oracle 🔮    ← Mind, Money & Meaning (HP 25-26, 29-34)
  ├── Eden 🌳      ← Food Growing & Agriculture (LE farming)
  ├── Empire 👑    ← Business & Wealth (GE 75-84)
  ├── Haven 🌿     ← Home, Family & Estate (LE 50-58)
  ├── Guardian 🛡️  ← Security & Compliance (UI 11 + HP 31 + LE 57)
  ├── Grid ⚙️      ← Infrastructure Ops (UI 00-12)
  └── Trinity 🔺   ← Community & Partnerships (GE 75 + Community)

  + NATS → Agno Workers (stateless):
    researcher, coder, analyst, web_automator, db_manager
```

---

## Agent Roster

### 1. Legacy ⚡ — Sovereign Orchestrator

| Field | Value |
|-------|-------|
| **Agent ID** | `legacy` |
| **Workspace** | `openclaw/workspaces/legacy/` |
| **Model** | `azure-openai/gpt-4o` (primary), `gpt-5.2` (available) |
| **Tool Profile** | `full` — all tools, exec, browser, fs, voice |
| **Codex Scope** | All domains (00-99) |
| **Heartbeat** | Every 30m via Telegram (7am-11pm ET) |
| **Channel** | Telegram (primary), iMessage, Discord, Voice |

**Role**: The conductor. Legacy doesn't do everything — it ensures everything gets done by the right agent. Receives all inbound messages, reasons about routing, dispatches to domain agents, synthesizes results, and responds to Zeus.

**Skills**: nexus-bridge, todoist, obsidian, github-api, browser, home-assistant, onepassword, vision, desktop-control, twilio, resend, spotify, transit, uber, recipe-scraper, lob, clickup, zoho-cli, hitl-escalate, proactive-agent, agent-browser

**Lobster Routines**: morning-brief, habit-checker, tech-digest, project-state, earnings-tracker, second-brain, health-check, git-sync

---

### 2. Spartan ⚔️ — Body & Physical Mastery

| Field | Value |
|-------|-------|
| **Agent ID** | `spartan` |
| **Workspace** | `openclaw/workspaces/spartan/` |
| **Model** | `azure-openai/gpt-4o` |
| **Tool Profile** | `messaging` + domain skills |
| **Codex Scope** | `HP 27 (Body)` + `HP 28 (Mastery)` |
| **Heartbeat** | Every 1h via Telegram (6am-10pm ET) |
| **Replaces** | HP agent (Pulse) — body domains only |

**Role**: Zeus's physical performance engine. Tracks fitness, nutrition, sleep, biometrics, recovery, body composition, and career mastery. Prescribes — never suggests. Data-driven, no fluff.

**Codex Domains**:

| Slot | Domain | What Spartan Owns |
|------|--------|-------------------|
| 27 | Body | Fitness, sleep, recovery, movement, biometrics, nutrition |
| 28 | Mastery | Skills, certifications, training, career discipline |

**Skills Needed**:
- `fitness-tracker` — workout logging, PR tracking, program management
- `nutrition-planner` — meal plans, macros, supplement stacks
- `biometrics` — wearable data pulls, body composition, health metrics
- `sleep-tracker` — sleep score integration, sleep hygiene protocols
- `todoist` (shared) — task management
- `obsidian` (shared) — codex read/write
- `weather` (shared) — outdoor training conditions

**Crons**:
- Morning check-in (6am) — sleep score, today's workout, daily priorities
- Evening review (9pm) — log day's metrics, nutrition summary
- Weekly PI score report (Sunday 8am) — composite performance score

**Escalation Rules**: Health emergency indicators, PI score drops >20%, supplement/medication interactions, missed 3+ consecutive workouts

---

### 3. Oracle 🔮 — Mind, Money & Meaning

| Field | Value |
|-------|-------|
| **Agent ID** | `oracle` |
| **Workspace** | `openclaw/workspaces/oracle/` |
| **Model** | `azure-openai/gpt-4o` |
| **Tool Profile** | `messaging` + domain skills |
| **Codex Scope** | `HP 25-26, 29-34` |
| **Heartbeat** | Every 1h via Telegram (7am-10pm ET) |
| **New Agent** | Splits inner HP domains from Spartan |

**Role**: Zeus's strategic advisor for the inner dimensions of personal mastery. Manages mindset, personal finance, learning, relationships, life systems, spirituality, and recreation. Calm, strategic, sees further than today. Framework-driven.

**Codex Domains**:

| Slot | Domain | What Oracle Owns |
|------|--------|------------------|
| 25 | Identity & Mind | Self-awareness, purpose, psychology, mindset |
| 26 | Intelligence | Education, learning goals, books, research |
| 29 | Life Systems | Routines, habits, daily ops, productivity |
| 30 | Finance | Personal budget, banking, credit, investments, taxes |
| 31 | Sovereignty | Legal rights, personal contracts, autonomy |
| 32 | Relationships | Social network, friendships, mentors, networking |
| 33 | Joy | Recreation, travel, entertainment, hobbies |
| 34 | Spirituality | Metaphysics, meditation, personal rituals, philosophy |

**Skills Needed**:
- `finance-personal` — budget tracking, savings rate, credit monitoring
- `learning-tracker` — book reading, course progress, certification tracking
- `habit-tracker` — habit streaks, routine compliance
- `todoist` (shared) — task management
- `obsidian` (shared) — codex read/write
- `weather` (shared) — travel/recreation planning

**Crons**:
- Morning priorities (7am) — today's key items across mind/money/meaning
- Finance check (daily 9am) — spending patterns, bills due
- Learning streak check (daily 8pm) — reading/course progress
- Weekly life review (Sunday 10am) — habits, finances, growth
- Relationship reminders (weekly) — birthdays, check-ins due

**Escalation Rules**: Financial decisions >$500. Mental health crisis indicators. Legal deadlines. Credit score drops.

---

### 4. Eden 🌳 — Food Growing & Agriculture

| Field | Value |
|-------|-------|
| **Agent ID** | `eden` |
| **Workspace** | `openclaw/workspaces/eden/` |
| **Model** | `azure-openai/gpt-4o` |
| **Tool Profile** | `messaging` + domain skills |
| **Codex Scope** | `LE` (farming/agriculture focus) |
| **Heartbeat** | Every 1h via Telegram (8am-9pm ET) |
| **Replaces** | LE agent (Steward) — farming niche |

**Role**: The food growing specialist. Manages all agriculture, gardening, permaculture, food forests, seed saving, composting, and seasonal planting. Thinks in growing cycles and soil health. Eden makes the land produce.

**Codex Domains**:

| Slot | Domain | What Eden Owns |
|------|--------|----------------|
| 50 | Homestead Index | Property inventory (garden plots, growing zones) |
| — | Garden | Planting schedules, crop rotation, harvest tracking |
| — | Food Forest | Perennial food systems, permaculture design |
| — | Composting | Compost systems, soil amendments, worm farms |
| — | Seed Saving | Seed inventory, germination tracking, heirloom varieties |
| — | Livestock (food) | Chickens, bees, aquaponics — food production animals |
| — | Land Management | Soil testing, cover crops, erosion control |
| — | Weather | Frost dates, growing conditions, microclimate tracking |

**Skills Needed**:
- `garden-planner` — planting calendar, crop rotation, harvest logs
- `weather` (shared) — forecasts, frost alerts, growing conditions
- `sensor-monitor` — soil moisture, temperature, humidity
- `seed-library` — seed inventory, germination data, variety tracking
- `todoist` (shared) — task management
- `obsidian` (shared) — codex read/write

**Crons**:
- Morning garden check (7am) — weather, soil conditions, today's tasks
- Weather alert monitor (every 3h) — frost warnings, extreme heat
- Weekly garden review (Saturday 8am) — harvest log, planting needs
- Seasonal planning (1st of month) — rotation schedule, seed ordering

**Escalation Rules**: Frost/freeze warnings — immediate. Pest outbreak indicators. Equipment failure. Water system anomalies.

**Memory Boundary**: Garden/growing data stays in Eden. Weather shared with Haven and Spartan.

---

### 5. Empire 👑 — Business & Wealth

| Field | Value |
|-------|-------|
| **Agent ID** | `empire` |
| **Workspace** | `openclaw/workspaces/empire/` |
| **Model** | `azure-openai/gpt-5.2` |
| **Tool Profile** | `messaging` + domain skills |
| **Codex Scope** | `75-99 GE (Generational Endeavor)` |
| **Heartbeat** | Every 1h via Telegram (7am-11pm ET) |
| **Replaces** | GE agent (Atlas) |

**Role**: The business brain. ROI-first, numbers-driven, deal-oriented. Manages all ventures, revenue tracking, client relationships, marketing, operations, and investment analysis. Ends every response with a next action. Thinks in terms of leverage and scalability.

**Codex Domains**:

| Slot | Domain | What Empire Owns |
|------|--------|------------------|
| 75 | Zeus Delacruz Enterprise | Personal brand, content creation |
| 76 | Legacy Acadacom | Education platform, courses |
| 77 | Legacy Platform Holdings | IP, holding company |
| 78 | Client Pipeline | Active clients, deals, revenue |
| 79 | Community (Business) | Paid community, member LTV |
| 80 | Marketing | Content strategy, social, email campaigns |
| 81 | Operations | SOPs, team management, tools |
| 82 | Legal (Business) | Business legal, contracts, compliance |
| 83 | Investments | Portfolio, real estate, crypto |
| 84 | Partnerships | Strategic alliances, JVs |

**Skills Needed**:
- `zoho-crm` — leads, contacts, deals, pipeline management
- `revenue-tracker` — MRR/ARR tracking, revenue reports
- `content-planner` — content calendar, publishing schedule
- `client-pipeline` — deal tracking, follow-ups, win rates
- `brand-manager` — brand assets, website updates, SEO
- `todoist` (shared) — task management
- `obsidian` (shared) — codex read/write

**Crons**:
- Daily revenue check (8am) — yesterday's numbers, pipeline changes
- Client health scan (Mon/Thu) — at-risk accounts, follow-up needed
- Weekly business report (Sunday 9am) — WI score, revenue, pipeline
- Content calendar check (daily 10am) — what's scheduled, what's overdue

**Escalation Rules**: Financial decisions >$1,000. Client churn risk. Legal/contract deadlines. Revenue drop >10% week-over-week.

**Memory Boundary**: NEVER share client details between domains. Business data stays in Empire.

---

### 6. Haven 🌿 — Home, Family & Estate

| Field | Value |
|-------|-------|
| **Agent ID** | `haven` |
| **Workspace** | `openclaw/workspaces/haven/` |
| **Model** | `azure-openai/gpt-4o` |
| **Tool Profile** | `messaging` + domain skills |
| **Codex Scope** | `LE 50-58` (all Legacy Evolution) |
| **Heartbeat** | Every 1h via Telegram (6am-9pm ET) |
| **New Agent** | No predecessor |

**Role**: The estate steward. Manages the full LE domain: home, family, marriage, children, elders, heritage, estate planning, property, and smart home systems. Haven is the family’s operational backbone — warm but strategic, thinks in 20-50 year horizons. Family data never leaves this agent.

**Codex Domains**:

| Slot | Domain | What Haven Owns |
|------|--------|-----------------|
| 50.1 | Property Index | Property inventory, valuations |
| 50.2 | Property Management | Maintenance schedules, repairs, vendors |
| — | Garden | Planting schedules, crop rotation, harvest tracking |
| — | Livestock | Animal care, feeding schedules, veterinary |
| — | Smart Home | Home Assistant devices, automations, scenes |
| — | Water Systems | Irrigation, well monitoring, rain collection |
| — | Land Management | Fencing, roads, outbuildings, land improvement |
| — | Weather | Forecasts, frost alerts, growing conditions |
| — | Security (Physical) | Cameras, gates, perimeter, fire safety |

**Skills Needed**:
- `home-assistant` — smart home control, device management, automations
- `weather` — forecasts, severe weather alerts, growing conditions
- `garden-planner` — planting calendar, crop rotation, harvest logs
- `livestock-tracker` — feeding schedules, health records, breeding
- `water-systems` — irrigation scheduling, well levels, rain collection
- `sensor-monitor` — soil moisture, temperature, humidity, air quality
- `property-maintenance` — maintenance schedules, vendor contacts, repair logs
- `obsidian` (shared) — codex read/write

**Crons**:
- Morning farm check (6am) — weather, sensor readings, today's tasks
- Weather alert monitor (every 2h) — severe weather, frost, heat
- Irrigation schedule (daily 5am) — water needs based on weather + sensors
- Weekly property review (Saturday 8am) — maintenance due, supply needs
- Seasonal planning (1st of month) — planting, harvesting, rotation

**Escalation Rules**: Severe weather warnings. Equipment failure. Livestock health alerts. Water system anomalies. Security breaches (physical).

---

### 7. Guardian 🛡️ — Security & Legal

| Field | Value |
|-------|-------|
| **Agent ID** | `guardian` |
| **Workspace** | `openclaw/workspaces/guardian/` |
| **Model** | `azure-openai/gpt-4o` |
| **Tool Profile** | `messaging` + domain skills |
| **Codex Scope** | `HP 31 (Legal)` + `LE 57 (Family Legal)` + `UI 11 (Cyber)` |
| **Heartbeat** | Every 2h via Telegram (8am-10pm ET) |
| **New Agent** | No predecessor |

**Role**: The shield. Manages all security and legal affairs — personal legal, family legal, business compliance, cybersecurity, credential management, insurance, and asset protection. Never acts without verification. Everything is documented with an audit trail.

**Codex Domains**:

| Slot | Domain | What Guardian Owns |
|------|--------|-------------------|
| 11.1 | Accounts Recovery | Account access, recovery procedures |
| 11.2 | Business Security | Business cybersecurity, compliance |
| 11.3 | Domain Management | DNS, domain registrations, renewals |
| 11.4 | Personal Security | Personal cybersecurity, threat model |
| 11.5 | Credentials | Credential inventory, rotation schedule |
| 31.3 | Criminal Record | Record management, expungement |
| 31.4 | Identity & Vitals | Birth certificates, passport, SSN |
| 31.7 | Protection | Firearms, self-defense, permits |
| 31.10 | Liability | Claims, disputes, liability |
| 31.11 | Asset Protection | Privacy, trusts, entity structure |
| 57.1 | Family Legal Index | Family legal overview |
| 57.2 | Family Legal Cases | Active cases |
| 57.3 | Family Documents | Wills, POAs, vital docs |
| 57.4 | Trusts & Estate | Trust management, estate planning |
| 57.5 | Insurance | All insurance policies, claims |

**Skills Needed**:
- `onepassword` — credential queries, vault management
- `legal-tracker` — case tracking, deadline management, document status
- `insurance-monitor` — policy renewals, coverage gaps, claim status
- `credit-repair` — credit monitoring, dispute management
- `compliance-checker` — regulatory deadlines, filing reminders
- `audit-trail` — log all security-relevant actions
- `obsidian` (shared) — codex read/write

**Crons**:
- Credential rotation check (weekly) — expiring passwords, tokens, API keys
- Insurance review (monthly) — upcoming renewals, coverage gaps
- Legal deadline scan (daily 9am) — court dates, filing deadlines
- Threat assessment (weekly) — unusual access patterns, security alerts
- Domain/SSL expiry check (weekly) — DNS, certs, registrations

**Escalation Rules**: ALL actions require Zeus's approval before execution. Security anomalies — escalate immediately. Legal deadlines within 7 days. Insurance claims. Credential compromises.

---

### 8. Grid ⚙️ — Infrastructure Ops

| Field | Value |
|-------|-------|
| **Agent ID** | `grid` |
| **Workspace** | `openclaw/workspaces/grid/` |
| **Model** | `azure-openai/gpt-4o` |
| **Tool Profile** | `messaging` + domain skills |
| **Codex Scope** | `UI 00-12 (tech, infra, dev)` |
| **Heartbeat** | Every 1h via Telegram (24/7) |
| **New Agent** | No predecessor |

**Role**: The system brain. Monitors all infrastructure — NATS, PostgreSQL, OpenClaw gateway, Agno workers, Docker services, VPS health, GitHub repos, CI/CD. Legacy shouldn't debug its own plumbing — Grid does that.

**Codex Domains**:

| Slot | Domain | What Grid Owns |
|------|--------|----------------|
| 00 | Framework | Doctrine, architecture, tech stack |
| 01 | Command | Dashboard infrastructure, navigation |
| 02 | Agents | Agent configs, memory, tools |
| 03 | Protocols | SOPs, operational procedures |
| 04 | Resources | Assets, templates, references |
| 08 | Infrastructure | Cloud, storage, integrations, email |
| 09 | Automation | Workflows, n8n, automation platforms |
| 10 | Development | Code, testing, DevOps, documentation |
| 12 | Research | R&D, robotics, communications |

**Skills Needed**:
- `github-api` — repo management, PR reviews, issue tracking
- `system-monitor` — VPS health, CPU/memory/disk, process status
- `nats-health` — JetStream status, stream health, consumer lag
- `deployment` — service restarts, deploy scripts, rollback
- `db-manager` — PostgreSQL health, query performance, backups
- `backup-verify` — verify backup integrity, schedule checks
- `obsidian` (shared) — codex read/write

**Crons**:
- Service health check (every 30m, 24/7) — all systemd services
- SSL/cert expiry (weekly) — Tailscale, domains
- Backup verification (daily 3am) — DB, codex, workspace backups
- Disk usage alert (daily) — VPS storage, log rotation
- NATS stream health (every 1h) — consumer lag, stream size
- GitHub PR review (daily 10am) — open PRs, stale branches

**Escalation Rules**: Service down — escalate immediately. Disk >80% — alert. Backup failure — alert. NATS stream errors. Database connection issues.

---

### 9. Trinity 🔺 — Community & Public

| Field | Value |
|-------|-------|
| **Agent ID** | `trinity` |
| **Workspace** | `openclaw/workspaces/trinity/` |
| **Model** | `azure-openai/gpt-4o` |
| **Tool Profile** | `messaging` + domain skills |
| **Codex Scope** | `Community` + `GE 75 (Brand)` |
| **Heartbeat** | Every 2h via Discord (9am-10pm ET) |
| **Replaces** | Community agent |

**Role**: The public voice. Manages Circle.so community, Discord server, content publishing, email campaigns, social media, and brand representation. Different personality from private agents — warm, inviting, genuinely helpful. Represents OneMind OS and Zeus Delacruz to the outside world.

**Codex Domains**:

| Slot | Domain | What Trinity Owns |
|------|--------|-------------------|
| 75.2 | Content Creation | Blog posts, videos, social content |
| 75.4 | Personal Brand | Brand assets, website, public image |
| — | Circle.so | Space moderation, member engagement, welcome flows |
| — | Discord | Server management, community support |
| — | Email | Newsletters, announcements, drip campaigns |
| — | Social Media | Platform management, engagement, analytics |

**Skills Needed**:
- `circle-api` — Circle.so posts, spaces, member management
- `discord` — Discord server management, channel posting
- `content-publisher` — schedule and publish content across platforms
- `email-campaigns` — newsletter creation, drip sequences
- `social-media` — cross-platform posting, engagement tracking
- `todoist` (shared) — task management
- `obsidian` (shared) — codex read/write

**Crons**:
- Community digest (daily 9am) — new members, trending posts, engagement
- Discord health (every 2h) — unanswered questions, engagement
- Content schedule check (daily 10am) — what's due, what's overdue
- Member retention scan (weekly) — at-risk members, churn signals
- Email campaign status (Mon/Thu) — open rates, click rates

**Escalation Rules**: Abusive/toxic community members — escalate. Member churn spike. Brand reputation issues. Content publishing failures.

**Personality Note**: Trinity is the ONLY agent that speaks publicly. All other agents are Zeus-private. Trinity's tone is warm, helpful, and community-oriented — never corporate, never salesy.

---

## Shared Skills

These live in `~/.openclaw/skills/` and are available to all agents:

| Skill | Used By | Purpose |
|-------|---------|---------|
| `todoist` | All | Task creation, queries, updates |
| `obsidian` | All | Codex vault read/write |
| `weather` | Haven, Spartan, Oracle | Weather forecasts, alerts |
| `web_fetch` | All | Fetch web pages |
| `web_search` | All | Search the web |

## Agent-Only Skills

These live in each agent's `<workspace>/skills/` directory:

| Skill | Agent | Reason |
|-------|-------|--------|
| `nexus-bridge` | Legacy | Only Legacy dispatches via NATS |
| `home-assistant` | Haven | Only Haven controls the home |
| `zoho-crm` | Empire | Client data stays in business |
| `onepassword` | Guardian | Credential access is security-scoped |
| `github-api` | Grid | Only infra agent touches repos |
| `circle-api` | Trinity | Only Trinity manages community |
| `discord` | Trinity | Only Trinity posts to Discord |
| `twilio` | Legacy | Only Legacy makes voice calls |
| `vision` | Legacy | Vision processing stays central |

---

## NATS + Agno Workers (Stateless)

These are NOT OpenClaw agents. They're lightweight Agno workers dispatched via NATS:

| Worker | Purpose | Dispatched By |
|--------|---------|---------------|
| `researcher` | Deep research (30+ min) | Legacy |
| `coder` | Code writing, debugging | Legacy / Grid |
| `analyst` | Data analysis, reports | Legacy / Empire |
| `web_automator` | Web scraping, automation | Legacy |
| `db_manager` | Database queries | Legacy / Grid |

---

## Concurrency Model

| Setting | Value | Impact |
|---------|-------|--------|
| Max concurrent subagents | 8 (gateway-wide) | Shared pool across all agents |
| Max spawn depth | 2 | Agent → Child → Grandchild |
| Max children per agent | 5 | Each agent can spawn 5 subs |
| Max concurrent cron runs | 3 | Prevents cron storms |
| Subagent model | `azure-openai/gpt-4o` | Cheaper than main agents |

---

## Memory Boundaries

| Agent | Can Read | Cannot Read |
|-------|----------|-------------|
| Legacy | All domains | — |
| Spartan | HP 25-34 only | Family (LE), Client data (GE) |
| Eden | LE 50-58 only | Client data (GE), Infra (UI) |
| Empire | GE 75-99 only | Family (LE), Personal health (HP) |
| Haven | LE 50 + weather | Client data (GE), Personal finance (HP) |
| Guardian | HP 31 + LE 57 + UI 11 | Full family (LE), Business (GE) |
| Grid | UI 00-12 only | Family (LE), Personal (HP), Business (GE) |
| Trinity | Community + GE 75 | Family (LE), Personal (HP), Internal ops |

---

## VPS Deployment Settings

The gateway config lives at `/home/onemind/.openclaw/openclaw.json` on the VPS. Below is the complete mapping of every agent's runtime settings as deployed.

### Discord Channel Bindings

| Agent | Discord Channel | Channel ID |
|-------|----------------|------------|
| Legacy ⚡ | `#legacy-chat` | `1478604243509448705` |
| Spartan ⚔️ | `#spartan-chat` | `1480445911707549750` |
| Oracle 🔮 | `#oracle-chat` | `1480450417522643127` |
| Eden 🌳 | `#eden-chat` | `1480445924202385499` |
| Empire 👑 | `#empire-chat` | `1480445930661744801` |
| Haven 🌿 | `#haven-chat` | `1480445937665966161` |
| Guardian 🛡️ | `#guardian-chat` | `1480445943789785180` |
| Grid ⚙️ | `#grid-chat` | `1480445951805096007` |
| Trinity 🔺 | `#trinity-chat` | `1480445958876827714` |

**Discord Guild**: `1453411552181883027` (One Mind OS)

### Discord Server Layout

| Category | Channels |
|----------|----------|
| 🧭 COMMAND CENTER | legacy-chat, agent-heartbeat, morning-brief, notifications, legacy-voice |
| 🧠 UNIFIED INTELLIGENCE | grid-chat, guardian-chat |
| 💪 HOLISTIC PERFORMANCE | spartan-chat, oracle-chat |
| 🏠 LEGACY EVOLUTION | eden-chat, haven-chat |
| 🚀 GENERATIONAL ENTERPRISE | empire-chat, trinity-chat |

### Agent Runtime Configuration

| Agent | Model | Heartbeat To | Active Hours (ET) | Codex Extra Path |
|-------|-------|-------------|-------------------|------------------|
| Legacy | gpt-5.2 | `${HEARTBEAT_CHANNEL_ID}` | 07:00–23:00 | `/home/onemind/onemind-os/codex` (all) |
| Spartan | gpt-4o | `1480445911707549750` | 08:00–22:00 | `codex/25-49 HP (Holistic Performance)` |
| Oracle | gpt-4o | `1480450417522643127` | 07:00–23:00 | `codex/25-49 HP (Holistic Performance)` |
| Eden | gpt-4o | `1480445924202385499` | 07:00–20:00 | `codex/50-74 LE (Legacy Evolution)` |
| Empire | gpt-5.2 | `1480445930661744801` | 07:00–23:00 | `codex/75-99 GE (Generational Endeavor)` |
| Haven | gpt-4o | `1480445937665966161` | 09:00–20:00 | `codex/50-74 LE (Legacy Evolution)` |
| Guardian | gpt-4o | `1480445943789785180` | 00:00–23:59 | `codex/00-24 SO (Sovereign Operations)` |
| Grid | gpt-4o | `1480445951805096007` | 07:00–23:00 | `codex/00-24 SO (Sovereign Operations)` |
| Trinity | gpt-4o | `1480445958876827714` | 08:00–22:00 | `codex/75-99 GE (Generational Endeavor)` |

### Gateway Defaults

| Setting | Value |
|---------|-------|
| Default model | `azure-openai/gpt-4o` |
| Embedding model | `text-embedding-3-small` (OpenAI) |
| Heartbeat interval | `1h` |
| Context pruning | `cache-ttl` / 1h |
| Compaction | `safeguard` |
| Thinking default | `medium` |
| Subagent model | `azure-openai/gpt-4o-mini` |
| Max concurrent subagents | 8 |
| Max spawn depth | 3 |
| Max children per agent | 5 |
| Max concurrent cron runs | 3 |
| Session dmScope | `per-channel-peer` |
| Gateway bind | `loopback` |
| Gateway port | `18789` |
| Tailscale | `serve` mode |

### VPS Paths

| Path | Purpose |
|------|---------|
| `/opt/onemind/openclaw/workspace/` | Legacy's deployed workspace |
| `~/.openclaw/agents/<agentId>/` | Each agent's state, auth, sessions |
| `~/.openclaw/skills/` | Shared skills (available to all agents) |
| `/home/onemind/onemind-os/codex/` | Codex vault (all agents read this) |
