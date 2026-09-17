# Discord × OpenClaw: Complete Technical Reference for One Mind OS

**OpenClaw is a free, open-source personal AI agent gateway that bridges Discord (and 12+ other platforms) to LLM-powered autonomous assistants**, making it the ideal middleware layer for building a "Life Operating System" on Discord. This reference maps every Discord feature against OpenClaw's capabilities, provides implementation architecture for a One Mind OS built on the Agno AI framework, and documents the technical boundaries of each platform. OpenClaw uses discord.js under the hood, connects to Discord's full Gateway API v10, and supports interactive components v2, slash commands, forum channels, voice, and DM pairing — covering roughly **85–90%** of Discord's native feature surface through its gateway model.

---

## What OpenClaw actually is

**OpenClaw** (MIT license, 180K+ GitHub stars) is a self-hosted personal AI agent platform created by Peter Steinberger (PSPDFKit founder). Originally launched as "Clawdbot" in November 2025, it was renamed to Moltbot and then OpenClaw in January 2026. On February 14, 2026, Steinberger announced he is joining OpenAI, and the project will move to an open-source foundation.

OpenClaw is **not** a traditional iPaaS or middleware. It is an **agent-native runtime** that runs a Gateway process on the user's own hardware (macOS, Linux, Windows via WSL2), routing messages between chat platforms and LLM backends. The Gateway owns all platform connections, manages sessions, routes conversations to agents, and persists memory in a local SQLite database.

**Core capabilities at a glance:**

- **12+ messaging channels**: Discord, WhatsApp, Telegram, Slack, Signal, Google Chat, iMessage, Teams, Matrix, IRC, Twitch, and more
- **Model-agnostic LLM support**: Claude, GPT, Gemini, DeepSeek, Ollama (local), OpenRouter
- **Persistent memory**: Cross-session, cross-platform memory via MEMORY.md and memory tools
- **Proactive automation**: Cron jobs, reminders, morning briefings, background monitoring
- **Browser control**: CDP-based Chrome automation for web browsing, form filling, screenshots
- **Skills system**: 50+ built-in skills; extensible via SKILL.md files
- **Plugin system**: Deep Gateway extensions in TypeScript/JavaScript
- **Webhook routing**: HTTP endpoints for external event ingestion with template-based routing
- **Voice**: Wake word + talk mode with ElevenLabs integration
- **100% free**: Only cost is LLM API usage ($5–100+/month depending on volume)

### Installation

```bash
npm install -g openclaw@latest
openclaw onboard --install-daemon
openclaw gateway --port 18789
```

### Security model

OpenClaw uses a **single-operator trust boundary**. The Gateway authenticates via a shared secret token stored in `~/.openclaw/openclaw.json`. Channel-level security uses DM pairing codes (expire after 1 hour), per-channel allowlists, and per-guild user restrictions. It binds to loopback by default. **It is not designed for multi-tenant or adversarial environments** — for mixed-trust teams, separate gateways per trust boundary are recommended.

Known security concerns include a critical 1-click RCE vulnerability (GHSA-g8p2-7wf7-98mq) and reports of malicious skills in the marketplace performing data exfiltration (Cisco/Snyk found 7.1% of analyzed skills designed to leak secrets).

---

## Every Discord feature cataloged

### Text communication features

Discord's text system spans seven distinct channel types plus a rich messaging layer:

**Channel types**: Standard text channels live inside categories (collapsible organizational folders). **Forum channels** (Community servers only) treat each post as its own thread with tag-based categorization and required-tag enforcement. **Media channels** provide gallery-style layouts optimized for images and video. **Announcement channels** (formerly News) allow messages to be "published" and crossposted to all following servers. **Threads** come in three variants — public (visible to all with parent channel access), private (invite-only, requires Boost Level 2), and announcement threads (within announcement channels).

**Messaging layer**: Messages support **rich embeds** (title, description, up to 25 fields, footer, image, thumbnail, author, color), markdown formatting (bold, italic, code blocks with syntax highlighting, headers H1–H3, spoiler tags, masked links), file attachments with inline preview, stickers (custom server + default Discord), custom/animated emojis (50 default slots, up to 250 at Boost Level 3), reactions, **Super Reactions** (animated effects, Nitro subscribers), **polls** (multi-option, configurable duration), **voice messages** (up to 20 minutes, waveform visualization), message replies with linked references, message forwarding, pinning, and full-text search with advanced filters.

### Voice, video, and streaming

**Voice channels** are persistent rooms — click to join, no ringing. They support embedded text chat, video (with background replacement), screen share, and Go Live streaming. Audio bitrate scales with boost level: **96 kbps** default → **128** at Level 1 → **256** at Level 2 → **384 kbps** at Level 3. Built-in Krisp noise suppression, push-to-talk or voice activity detection, and priority speaker permission are standard.

**Stage channels** provide a speaker/audience webinar model for Community servers, supporting up to 50 speakers (300 with boosts) with Request to Speak functionality. **Discord Activities** are embedded web apps (games, Watch Together, etc.) running inside voice channels via the Embedded App SDK. **Clips** capture the last 30s–2min of a stream (Nitro early access). **Soundboard** provides short audio clips in voice (8 default slots, up to 48 at Level 3).

### Server structure and community

Servers support **templates** (saveable channel/role/permission configurations), **server discovery** (public listing in Discord's browser), **welcome screens**, **member verification gates** (rules screening), and multi-step **onboarding flows** with customizable questions that auto-assign channels and roles. **Scheduled events** support voice, stage, or external locations with RSVP and reminders. **Server insights** provide analytics on member growth, retention, engagement, join sources, and top channels for Community servers.

### Monetization

**Server Subscriptions** offer up to 3 tiered monthly plans with custom pricing, exclusive roles/channels/emojis (90/10 revenue split). **Server Products/Shop** enables one-time digital product sales (downloadable files, premium roles). **Premium App monetization** lets developers sell subscriptions and in-app purchases natively (15% platform fee on first $1M).

### Moderation and safety

**AutoMod** provides keyword filters (pre-made + custom with wildcard support), ML-driven spam detection, and mention spam protection with configurable actions (block, alert, timeout). Manual tools include **timeout** (60s–28 days), kick, ban (with message history purge), slowmode (5s–6h), and bulk message deletion. **Audit logs** (retained 45 days) track all administrative actions. The DAVE protocol provides **end-to-end encryption** for voice and video across all platforms.

---

## Discord permission flags — complete reference

Every permission in Discord's system is a bit flag. OpenClaw's bot inherits these permissions based on its role assignment.

| Permission | Bit | Hex | Category |
|---|---|---|---|
| `CREATE_INSTANT_INVITE` | 1 << 0 | `0x1` | General |
| `KICK_MEMBERS` | 1 << 1 | `0x2` | General |
| `BAN_MEMBERS` | 1 << 2 | `0x4` | General |
| `ADMINISTRATOR` | 1 << 3 | `0x8` | General |
| `MANAGE_CHANNELS` | 1 << 4 | `0x10` | General |
| `MANAGE_GUILD` | 1 << 5 | `0x20` | General |
| `ADD_REACTIONS` | 1 << 6 | `0x40` | Text |
| `VIEW_AUDIT_LOG` | 1 << 7 | `0x80` | General |
| `PRIORITY_SPEAKER` | 1 << 8 | `0x100` | Voice |
| `STREAM` | 1 << 9 | `0x200` | Voice |
| `VIEW_CHANNEL` | 1 << 10 | `0x400` | General |
| `SEND_MESSAGES` | 1 << 11 | `0x800` | Text |
| `SEND_TTS_MESSAGES` | 1 << 12 | `0x1000` | Text |
| `MANAGE_MESSAGES` | 1 << 13 | `0x2000` | Text |
| `EMBED_LINKS` | 1 << 14 | `0x4000` | Text |
| `ATTACH_FILES` | 1 << 15 | `0x8000` | Text |
| `READ_MESSAGE_HISTORY` | 1 << 16 | `0x10000` | Text |
| `MENTION_EVERYONE` | 1 << 17 | `0x20000` | Text |
| `USE_EXTERNAL_EMOJIS` | 1 << 18 | `0x40000` | Text |
| `VIEW_GUILD_INSIGHTS` | 1 << 19 | `0x80000` | General |
| `CONNECT` | 1 << 20 | `0x100000` | Voice |
| `SPEAK` | 1 << 21 | `0x200000` | Voice |
| `MUTE_MEMBERS` | 1 << 22 | `0x400000` | Voice |
| `DEAFEN_MEMBERS` | 1 << 23 | `0x800000` | Voice |
| `MOVE_MEMBERS` | 1 << 24 | `0x1000000` | Voice |
| `USE_VAD` | 1 << 25 | `0x2000000` | Voice |
| `CHANGE_NICKNAME` | 1 << 26 | `0x4000000` | General |
| `MANAGE_NICKNAMES` | 1 << 27 | `0x8000000` | General |
| `MANAGE_ROLES` | 1 << 28 | `0x10000000` | General |
| `MANAGE_WEBHOOKS` | 1 << 29 | `0x20000000` | General |
| `MANAGE_GUILD_EXPRESSIONS` | 1 << 30 | `0x40000000` | General |
| `USE_APPLICATION_COMMANDS` | 1 << 31 | `0x80000000` | Text |
| `REQUEST_TO_SPEAK` | 1 << 32 | `0x100000000` | Voice |
| `MANAGE_EVENTS` | 1 << 33 | `0x200000000` | General |
| `MANAGE_THREADS` | 1 << 34 | `0x400000000` | Text |
| `CREATE_PUBLIC_THREADS` | 1 << 35 | `0x800000000` | Text |
| `CREATE_PRIVATE_THREADS` | 1 << 36 | `0x1000000000` | Text |
| `USE_EXTERNAL_STICKERS` | 1 << 37 | `0x2000000000` | Text |
| `SEND_MESSAGES_IN_THREADS` | 1 << 38 | `0x4000000000` | Text |
| `USE_EMBEDDED_ACTIVITIES` | 1 << 39 | `0x8000000000` | Voice |
| `MODERATE_MEMBERS` | 1 << 40 | `0x10000000000` | General |
| `VIEW_CREATOR_MONETIZATION_ANALYTICS` | 1 << 41 | `0x20000000000` | General |
| `USE_SOUNDBOARD` | 1 << 42 | `0x40000000000` | Voice |
| `CREATE_GUILD_EXPRESSIONS` | 1 << 43 | `0x80000000000` | General |
| `CREATE_EVENTS` | 1 << 44 | `0x100000000000` | General |
| `USE_EXTERNAL_SOUNDS` | 1 << 45 | `0x200000000000` | Voice |
| `SEND_VOICE_MESSAGES` | 1 << 46 | `0x400000000000` | Text |
| `SEND_POLLS` | 1 << 49 | `0x2000000000000` | Text |
| `USE_EXTERNAL_APPS` | 1 << 50 | `0x4000000000000` | Text |

---

## Feature mapping matrix: Discord → OpenClaw

This is the core reference. For each Discord feature category, the table shows native Discord behavior, OpenClaw access level, and implementation approach for a One Mind OS.

### Text channels and messaging

| Discord Feature | OpenClaw Access | Implementation | One Mind OS Use Case |
|---|---|---|---|
| Send/receive text messages | **Full** — via discord.js Gateway | Automatic; bot reads MESSAGE_CREATE events, sends via REST | Agent receives commands, delivers reports |
| Rich embeds | **Full** — structured embed JSON | Include `embeds` array in message payload | Daily dashboards, health reports, financial summaries |
| Forum channels | **Full** — create threads, post to forums | Send to forum parent channel ID; auto-creates thread | Knowledge base, project tracking, decision logs |
| Forum tags | **Full** — via API tag assignment | Set `applied_tags` array on thread creation | Categorize entries by pillar (Health, Legacy, Tech, Business) |
| Threads (public/private) | **Full** — create, post, archive | `POST /channels/{id}/threads` | Isolate conversations per topic/project |
| Announcement channels | **Full** — publish/crosspost | `POST /channels/{id}/messages/crosspost` | Broadcast critical family alerts across linked servers |
| Voice messages | **Full** — send/receive audio | Audio file handling through OpenClaw media pipeline | Voice memos, audio instructions to agents |
| Polls | **Partial** — can create via API | `POST /channels/{id}/messages` with poll object | Family decision-making, preference gathering |
| Reactions/Super Reactions | **Full** — add/read reactions | `PUT /channels/{id}/messages/{id}/reactions/{emoji}/@me` | Acknowledgment tracking, quick approval workflows |
| Stickers/Emojis | **Partial** — send with messages; cannot manage server stickers directly | Include emoji in message content; sticker via sticker_ids | Status indicators, mood tracking |
| Message components (buttons, selects, modals) | **Full** — Components V2 support | OpenClaw's structured `components` block with blocks, buttons, selects, modals | Interactive dashboards, approval workflows, form collection |
| Markdown formatting | **Full** — standard Discord markdown | Include markdown in message text | Formatted reports, code snippets, structured updates |
| File attachments | **Full** — send/receive files | Attach files to message payloads | Document sharing, report PDFs, image uploads |
| Slash commands | **Full** — native slash command support | OpenClaw registers commands via Discord API; handles INTERACTION_CREATE | `/report`, `/status`, `/schedule`, `/ask-agent` |
| Message search | **No** — no API for message search | Not available via bot API (Discord limitation) | Must use native Discord client |

### Voice and video

| Discord Feature | OpenClaw Access | Implementation | One Mind OS Use Case |
|---|---|---|---|
| Join/leave voice channels | **Full** — voice state management | Gateway opcode 4 (VOICE_STATE_UPDATE) | Agent joins family standup voice channel |
| Voice channel events | **Full** — VOICE_STATE_UPDATE events | Listen for join/leave/mute/deafen events | Track family member availability, log meeting attendance |
| Stage channels | **Partial** — can create/manage stages via API | `POST /stage-instances`, manage speaker permissions | Weekly family town halls, coaching sessions |
| Screen share / Go Live | **No** — requires native client | Cannot programmatically screen share via bot | Not automatable |
| Video calls | **No** — requires native client | Not available via bot API | Not automatable |
| Discord Activities | **No** — requires Embedded App SDK | Would need separate Activity app build | Could build custom One Mind OS Activity |
| Clips | **No** — client-side feature | Not accessible via API | Not automatable |
| Soundboard | **Partial** — can manage sounds via API | GUILD_SOUNDBOARD_SOUND events | Alert sounds for priority notifications |
| Text-in-voice | **Full** — standard text channel behavior | Text channel associated with voice channel | Share links/notes during voice meetings |

### Server management

| Discord Feature | OpenClaw Access | Implementation | One Mind OS Use Case |
|---|---|---|---|
| Create/modify channels | **Full** — MANAGE_CHANNELS permission | `POST /guilds/{id}/channels`, `PATCH /channels/{id}` | Dynamic channel creation for new projects/goals |
| Create/modify roles | **Full** — MANAGE_ROLES permission | `POST /guilds/{id}/roles`, `PATCH /guilds/{id}/roles/{id}` | Auto-assign family member roles, agent roles |
| Assign roles to members | **Full** — via REST API | `PUT /guilds/{id}/members/{id}/roles/{id}` | Automatic role assignment based on onboarding |
| Permission overrides | **Full** — per-channel overwrites | `PUT /channels/{id}/permissions/{overwrite_id}` | Dynamic access control per project/family member |
| Audit log access | **Full** — VIEW_AUDIT_LOG permission | `GET /guilds/{id}/audit-log` with filters | Security monitoring, activity tracking |
| Server templates | **Partial** — can create/use via API | `POST /guilds/templates/{template_code}` | Replicate One Mind OS structure to new families |
| Scheduled events | **Full** — CREATE_EVENTS permission | `POST /guilds/{id}/scheduled-events` | Family events, goal review sessions, coaching calls |
| Onboarding flows | **Limited** — can read, limited write | Guild onboarding endpoints exist but limited bot control | Initial setup requires manual configuration |
| Welcome screen | **Limited** — `PATCH /guilds/{id}/welcome-screen` | Can modify welcome screen programmatically | Auto-update welcome message with current priorities |
| Server insights | **No** — VIEW_GUILD_INSIGHTS is read-only, limited API | No comprehensive insights API | Use audit log + custom tracking instead |
| Server discovery | **No** — managed via Discord client | Not relevant for personal Life OS server | Keep server private |

### Moderation and automation

| Discord Feature | OpenClaw Access | Implementation | One Mind OS Use Case |
|---|---|---|---|
| AutoMod rules | **Full** — create/manage via API | `POST /guilds/{id}/auto-moderation/rules` | Content filtering, keyword-triggered automations |
| Kick/ban/timeout | **Full** — with appropriate permissions | REST API endpoints for moderation actions | Parental controls, enforce house rules |
| Bulk message delete | **Full** — MANAGE_MESSAGES permission | `POST /channels/{id}/messages/bulk-delete` | Clean up channels, archive old content |
| Slowmode | **Full** — channel setting | `PATCH /channels/{id}` with `rate_limit_per_user` | Control posting frequency in shared channels |

### Webhooks and integrations

| Discord Feature | OpenClaw Access | Implementation | One Mind OS Use Case |
|---|---|---|---|
| Create webhooks | **Full** — MANAGE_WEBHOOKS permission | `POST /channels/{id}/webhooks` | Route external events to specific channels |
| Execute webhooks | **Full** — no auth needed (token in URL) | `POST /webhooks/{id}/{token}` with embeds | Smart home alerts, CI/CD, calendar events |
| Platform integrations (Spotify, Twitch, etc.) | **No** — user-level connections only | Cannot control user connections via bot | Users configure natively |
| Linked roles | **Partial** — can configure role connection metadata | Application role connection metadata API | Auto-verify family members via external accounts |

### Premium features

| Discord Feature | OpenClaw Access | Implementation | One Mind OS Use Case |
|---|---|---|---|
| Nitro perks | **No** — user subscription only | Not applicable | Users purchase individually |
| Server boosts | **No** — user-initiated only | Cannot programmatically boost | Encourage family members to boost |
| Server subscriptions | **Limited** — can check entitlements via API | ENTITLEMENT events in Gateway | Gate premium agent features behind subscription tiers |
| Profile customization | **No** — client-side only | Not available via API | Not automatable |

---

## OpenClaw's Discord integration in depth

### Configuration

OpenClaw connects to Discord as a standard bot application via the discord.js library. Configuration lives in `~/.openclaw/openclaw.json`:

```json
{
  "channels": {
    "discord": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "streaming": "partial",
      "dmPolicy": "pairing",
      "groupPolicy": "allowlist",
      "guilds": {
        "1234567890": {
          "requireMention": false,
          "users": ["9876543210"]
        }
      }
    }
  }
}
```

**Required Discord bot permissions**: View Channels, Send Messages, Read Message History, Embed Links, Attach Files, Add Reactions. **Required privileged intents**: Message Content Intent (mandatory), Server Members Intent (recommended), Presence Intent (optional).

### Session routing model

OpenClaw uses a deterministic session routing system. DM conversations share the main agent session (`agent:main:main`). Guild channel conversations get **isolated session keys** (`agent:<agentId>:discord:channel:<channelId>`), ensuring context isolation between channels. This is critical for a Life OS — conversations in #health should not leak context into #finances.

### Interactive components v2

OpenClaw's component system wraps Discord's Components V2 API with a clean abstraction:

```json
{
  "channel": "discord",
  "action": "send",
  "to": "channel:123456789012345678",
  "components": {
    "reusable": true,
    "text": "Morning Check-in",
    "blocks": [
      {
        "type": "actions",
        "buttons": [
          { "label": "✅ Great", "style": "success" },
          { "label": "😐 Okay", "style": "secondary" },
          { "label": "❌ Rough", "style": "danger" }
        ]
      },
      {
        "type": "actions",
        "select": {
          "type": "string",
          "placeholder": "Top priority today?",
          "options": [
            { "label": "Health & Fitness", "value": "health" },
            { "label": "Family & Legacy", "value": "legacy" },
            { "label": "Business & Finance", "value": "business" },
            { "label": "Technology & Innovation", "value": "tech" }
          ]
        }
      }
    ],
    "modal": {
      "title": "Daily Intentions",
      "triggerLabel": "Set Goals",
      "fields": [
        { "type": "text", "label": "Top 3 priorities" },
        { "type": "text", "label": "Blockers or concerns" }
      ]
    }
  }
}
```

**Buttons** support 5 styles (primary, secondary, success, danger, link), access control via `allowedUsers`, and single-use or reusable modes. **Select menus** support string, user, role, mentionable, and channel types. **Modals** support up to 5 text input fields. Components are single-use by default; set `components.reusable = true` for persistent interactive elements.

### Webhook event routing

OpenClaw can ingest external webhooks and route them to Discord channels with template-based message construction:

```json
{
  "hooks": {
    "enabled": true,
    "token": "long-random-secret",
    "path": "/hooks",
    "routes": [
      {
        "match": { "path": "/hooks/home-assistant" },
        "action": "agent",
        "template": {
          "message": "🏠 Home Alert: {{body.event}} in {{body.location}}",
          "channel": "discord",
          "to": "channel:SMART_HOME_CHANNEL_ID"
        }
      },
      {
        "match": { "path": "/hooks/github" },
        "action": "agent",
        "template": {
          "message": "Push to {{body.repository.name}} by {{body.pusher.name}}",
          "channel": "discord",
          "to": "channel:DEV_CHANNEL_ID"
        }
      }
    ]
  }
}
```

---

## Agno framework powers the agent layer

**Agno** (formerly Phidata) is an open-source Python framework for building multi-agent AI systems. With **22K+ GitHub stars** and native Discord integration via `agno.integrations.discord.DiscordClient`, it provides the agent intelligence layer that sits behind OpenClaw or runs as a standalone Discord bot.

### Three-layer architecture

**Layer 1 — SDK**: Agent (single autonomous unit with LLM + tools), Team (coordinated agent groups with delegation/debate), Workflow (deterministic + agentic multi-step orchestration). **Layer 2 — AgentOS**: Production FastAPI runtime with JWT/RBAC security, per-session isolation, and tracing. **Layer 3 — UI**: Web dashboard at os.agno.com for monitoring and management.

### Native Discord integration

```python
from agno.agent import Agent
from agno.team import Team
from agno.integrations.discord import DiscordClient
from agno.models.anthropic import Claude
from agno.tools.duckduckgo import DuckDuckGoTools
from agno.db.sqlite import SqliteDb

# Create specialized agents
research_agent = Agent(
    name="Research Agent",
    role="Search the web and synthesize information",
    model=Claude(id="claude-sonnet-4-6"),
    tools=[DuckDuckGoTools()],
    learning=True,
)

life_os_team = Team(
    name="One Mind OS",
    mode="coordinate",
    members=[research_agent],  # Add more agents as needed
    model=Claude(id="claude-sonnet-4-6"),
    enable_memories=True,
    db=SqliteDb(db_file="one_mind_os.db"),
)

# Serve via Discord
discord_bot = DiscordClient(team=life_os_team)
discord_bot.serve()  # Reads DISCORD_BOT_TOKEN from env
```

Agno's DiscordClient automatically creates threads for each conversation, handles images/video/audio/documents, splits long messages (>1500 chars), displays reasoning content in italics, and enriches context with username and message URL. It supports both single Agent and Team modes.

### 80+ built-in toolkits

Agno provides toolkits for web search (DuckDuckGo, Exa), finance (YFinance), reasoning (chain-of-thought), code execution (Python, DuckDB, shell), media (Giphy, YouTube), data scraping (Newspaper, Firecrawl), **Discord operations** (send/delete messages, forums, reactions), GitHub, Docker, file system, and critically — **MCP (Model Context Protocol)** for connecting any MCP-compatible service (Home Assistant, Google Calendar, etc.).

---

## One Mind OS deployment architecture

### Recommended Discord server structure

The server is organized around four life pillars plus operational infrastructure:

```
🧠 ONE MIND OS
├── 📋 COMMAND CENTER
│   ├── #dashboard          (read-only — daily AI-generated summaries)
│   ├── #agent-chat         (primary AI interaction channel)
│   ├── #notifications      (announcement channel — crosspostable)
│   └── #audit-log          (bot activity logging)
│
├── 💪 HOLISTIC PERFORMANCE
│   ├── #health-tracking    (nutrition, exercise, sleep, vitals)
│   ├── #fitness-log        (forum channel — one thread per workout)
│   ├── #mental-wellness    (journal entries, meditation tracking)
│   ├── #biometrics         (webhook-fed from wearables/smart scale)
│   └── 🔊 coaching-voice   (voice channel for AI coaching sessions)
│
├── 🏛️ LEGACY EVOLUTION
│   ├── #family-calendar    (webhook-fed from Google Calendar)
│   ├── #family-updates     (announcement channel)
│   ├── #kid-milestones     (media channel — photos/videos)
│   ├── #education          (forum channel — tagged by child/subject)
│   ├── #traditions         (family rituals, values documentation)
│   └── 🔊 family-standup   (daily voice check-in channel)
│
├── 💰 GENERATIONAL ENTREPRENEURSHIP
│   ├── #financial-overview (read-only — AI dashboards)
│   ├── #business-ops       (active business management)
│   ├── #investments        (forum channel — one thread per asset)
│   ├── #opportunities      (deal flow, partnerships)
│   ├── #real-estate        (property management tracking)
│   └── 🔊 strategy-room    (voice channel for planning)
│
├── 🔧 INNOVATIVE TECHNOLOGY
│   ├── #smart-home         (Home Assistant webhook alerts)
│   ├── #dev-ops            (GitHub/CI webhooks)
│   ├── #ai-experiments     (forum channel — tagged by project)
│   ├── #robotics           (automation status, robot notifications)
│   └── #tech-research      (AI-curated articles and insights)
│
├── 🏠 LOCATION MANAGEMENT
│   ├── #home-primary       (alerts from primary residence)
│   ├── #home-vacation      (alerts from vacation property)
│   ├── #vehicles           (car status, maintenance reminders)
│   └── #security           (private — camera alerts, access logs)
│
└── ⚙️ SYSTEM
    ├── #bot-config         (admin-only — configuration commands)
    ├── #error-log          (agent errors and failures)
    ├── #cron-log           (scheduled task execution records)
    └── #integrations       (webhook test/debug channel)
```

### Agent architecture with Agno + OpenClaw

Two deployment models are viable. They can coexist:

**Model A — OpenClaw as gateway (recommended for personal use)**:
OpenClaw handles the Discord connection, session management, and cross-platform routing. Agno agents can be invoked via OpenClaw skills or plugins, or the LLM backend can be configured to use an Agno-powered API.

```
User ↔ Discord ↔ OpenClaw Gateway ↔ LLM (Claude/GPT)
                       ↓                    ↓
                  Session DB          Agno AgentOS API
                  MEMORY.md           (specialized agents)
                  Skills/Plugins
```

**Model B — Agno DiscordClient as direct bot**:
Agno's native DiscordClient connects directly to Discord's Gateway. Best when you want full Python control and don't need cross-platform routing.

```
User ↔ Discord ↔ Agno DiscordClient ↔ Agent Team
                                          ↓
                                     Tool Execution
                                     Memory/Knowledge DB
                                     External APIs
```

### Recommended agent team structure

```python
from agno.agent import Agent
from agno.team import Team

# Pillar agents
health_agent = Agent(
    name="Vitality Agent",
    role="Track health, nutrition, fitness, sleep, and mental wellness",
    tools=[HealthTools(), CalendarTools(), WearableAPITools()],
    instructions=[
        "Track all health metrics and identify trends",
        "Send morning check-in prompts to #health-tracking",
        "Alert to #notifications if any metric is concerning",
    ],
    learning=True,
)

family_agent = Agent(
    name="Legacy Agent", 
    role="Manage family calendar, milestones, education, multi-location coordination",
    tools=[GoogleCalendarTools(), ReminderTools(), LocationTools()],
    instructions=[
        "Coordinate schedules across family members and locations",
        "Track children's milestones and educational progress",
        "Send daily family briefing to #family-updates",
    ],
    learning=True,
)

finance_agent = Agent(
    name="Wealth Agent",
    role="Track investments, business metrics, real estate, and financial health",
    tools=[YFinanceTools(), SpreadsheetTools(), BankingAPITools()],
    instructions=[
        "Monitor portfolio performance and alert on significant changes",
        "Generate weekly financial summary for #financial-overview",
        "Track business KPIs and revenue targets",
    ],
    learning=True,
)

tech_agent = Agent(
    name="Innovation Agent",
    role="Manage smart home, development projects, AI experiments, robotics",
    tools=[HomeAssistantTools(), GitHubTools(), DockerTools()],
    instructions=[
        "Route smart home alerts to appropriate channels",
        "Monitor CI/CD pipelines and deployment status",
        "Research emerging technologies relevant to family goals",
    ],
    learning=True,
)

coordinator = Agent(
    name="One Mind",
    role="Central coordinator that routes requests and synthesizes across all life pillars",
    instructions=[
        "You are the unified intelligence layer for the family",
        "Route queries to the appropriate specialist agent",
        "Synthesize cross-pillar insights (e.g., how sleep affects business performance)",
        "Generate daily morning briefing covering all pillars",
    ],
    learning=True,
)

one_mind_os = Team(
    name="One Mind OS",
    mode="coordinate",
    leader=coordinator,
    members=[health_agent, family_agent, finance_agent, tech_agent],
    enable_memories=True,
    db=SqliteDb(db_file="one_mind_os.db"),
)
```

### Notification routing strategies

**Priority-based routing** maps urgency to Discord notification behavior:

| Priority | Channel Type | Discord Config | Example |
|---|---|---|---|
| **P0 — Critical** | Announcement channel (`#notifications`) | @everyone mention, crosspost to followers | Security breach, medical emergency, server down |
| **P1 — Urgent** | Standard text with role mention | @role mention for relevant family members | Schedule conflict, overdue payment, health metric anomaly |
| **P2 — Important** | Standard text, no mention | Normal message, users check on their schedule | Daily briefings, weekly reports, milestone updates |
| **P3 — Informational** | Forum channel thread | Post as new thread with appropriate tag | Research findings, article recommendations, general logs |
| **P4 — Ambient** | Webhook-fed channel, no notifications | Suppress @everyone, members opt-in | Smart home status, IoT telemetry, background task logs |

**Smart home / robotics routing** uses OpenClaw webhooks:

```bash
# Home Assistant automation → OpenClaw webhook → Discord channel
# In Home Assistant automation.yaml:
- alias: "Motion Alert to Discord"
  trigger:
    - platform: state
      entity_id: binary_sensor.front_door_motion
      to: "on"
  action:
    - service: rest_command.openclaw_webhook
      data:
        url: "http://localhost:18789/hooks/home-assistant"
        payload: '{"event": "motion_detected", "location": "front_door", "priority": "P1"}'
```

### Multi-location family management

Discord's category and channel structure maps cleanly to multi-location management. Each property gets its own category with channels for alerts, maintenance, and status. OpenClaw webhooks from each location's Home Assistant instance route to the corresponding Discord channel. The family agent coordinates across locations, detecting conflicts (e.g., "vacation home HVAC set to away mode but you have guests arriving Friday").

**Cross-platform continuity**: OpenClaw's unified inbox means a family member can ask about home status via WhatsApp, get a response that also posts to `#home-primary` on Discord, and continue the conversation on Telegram — all sharing the same agent session context.

---

## Discord Gateway events reference

For the One Mind OS, these Gateway events drive automation:

| Event | Trigger | One Mind OS Application |
|---|---|---|
| `MESSAGE_CREATE` | Any message sent | Core agent interaction loop |
| `INTERACTION_CREATE` | Slash command, button, select, modal | Interactive dashboard controls |
| `GUILD_MEMBER_ADD` | New member joins server | Onboarding workflow trigger |
| `VOICE_STATE_UPDATE` | User joins/leaves/mutes in voice | Track family member availability |
| `GUILD_SCHEDULED_EVENT_CREATE` | New event scheduled | Sync to family calendar |
| `GUILD_SCHEDULED_EVENT_USER_ADD` | User RSVPs to event | Attendance tracking |
| `AUTO_MODERATION_ACTION_EXECUTION` | AutoMod triggers | Content safety monitoring |
| `GUILD_AUDIT_LOG_ENTRY_CREATE` | Admin action taken | Security monitoring |
| `MESSAGE_REACTION_ADD` | Reaction added to message | Acknowledgment tracking, approval workflows |
| `THREAD_CREATE` | New thread created | Forum post tracking |

**Gateway intents required for One Mind OS**: GUILDS (1), GUILD_MEMBERS (2, privileged), GUILD_MODERATION (4), GUILD_VOICE_STATES (128), GUILD_MESSAGES (512), GUILD_MESSAGE_REACTIONS (1024), MESSAGE_CONTENT (32768, privileged), GUILD_SCHEDULED_EVENTS (65536). **Total intent value: 99,975**.

---

## What is NOT possible — hard limitations

### Discord API limitations that affect OpenClaw

**No message search API**: Discord provides no bot-accessible search endpoint. You cannot programmatically search message history by content — only paginate through messages chronologically. Workaround: maintain a local search index by processing MESSAGE_CREATE events.

**No screen share or video streaming via bot**: Bots cannot programmatically start screen shares or video streams. Voice channel joining and audio streaming are possible, but visual streaming requires a native Discord client. This limits automated presentation/dashboard-streaming use cases.

**No Discord Activities creation at runtime**: Activities require a separate Embedded App SDK build, hosted web app, and Discord Developer Portal registration. OpenClaw cannot dynamically create or modify Activities.

**No user connection management**: Bots cannot view, add, or modify user-level integrations (Spotify, Twitch, Steam connections). These are personal settings.

**No Nitro/Boost management**: Bots cannot purchase, apply, or manage Nitro subscriptions or server boosts programmatically.

**No native DM initiation without shared server**: Bots can only DM users who share a server with the bot or have an existing DM channel open.

**Audit log retention**: Fixed at **45 days** with no API to extend. For long-term compliance tracking, you must export and store audit log entries externally.

**Rate limits constrain throughput**: Gateway allows 120 events per 60 seconds. REST API uses per-route bucket rate limiting. For a personal Life OS this is rarely an issue, but bulk operations (e.g., mass channel creation, role assignment) require careful throttling.

### OpenClaw-specific limitations

**Single-operator trust model**: OpenClaw is designed for **one trusted operator**. There is no multi-user RBAC, no MFA, and no role separation within the Gateway. For a family server, this means one person administers the Gateway, and other family members interact only through Discord's interface.

**No native RBAC or audit trail**: OpenClaw does not provide its own permission system beyond channel-level allowlists. All authorization is coarse-grained (paired or not paired, allowed or not allowed).

**Cost unpredictability with autonomous agents**: LLM API costs can spiral with proactive cron jobs, background monitoring, and forgotten automations. Users have reported costs reaching **$170 in a single day** (90M tokens through Claude Opus) and $3,600/month from unmonitored workflows. Budget caps and model routing (80% to cheap models like GPT-4o-mini) are essential.

**Requires always-on infrastructure**: The Gateway must run 24/7 on a dedicated machine or VPS ($4–15/month for a VPS, $0 if running on existing hardware). If the machine sleeps or crashes, all agent functionality stops until restarted.

**Skills marketplace security**: Research has found that **7.1% of community skills** contain data exfiltration or prompt injection code. Only use trusted skills from verified sources, or write your own.

**Config fragility**: Some OpenClaw commands can overwrite configuration. Always backup `~/.openclaw/openclaw.json` before running setup commands.

### Feature gap summary

| Capability | Native Discord | OpenClaw | Agno Direct | Gap |
|---|---|---|---|---|
| Screen share/streaming | ✅ | ❌ | ❌ | Requires native client |
| Message search | ✅ (client) | ❌ | ❌ | No API endpoint |
| Activities/embedded apps | ✅ | ❌ | ❌ | Requires separate SDK build |
| Clips recording | ✅ (Nitro) | ❌ | ❌ | Client-side only |
| Profile customization | ✅ | ❌ | ❌ | Client-side only |
| Server discovery listing | ✅ | ❌ | ❌ | Not relevant for private server |
| Multi-tenant security | ❌ | ❌ | ✅ (AgentOS RBAC) | OpenClaw is single-operator |
| Cross-platform memory | ❌ | ✅ | ✅ | OpenClaw's key differentiator |
| Multi-channel routing | ❌ | ✅ | ❌ | OpenClaw's key differentiator |
| Proactive cron automation | ❌ | ✅ | ✅ (via FastAPI) | Neither platform is limited |

---

## Conclusion: architectural decision framework

**For a personal/family One Mind OS, OpenClaw is the stronger choice** over building a raw Discord bot with Agno alone. OpenClaw provides the always-on Gateway, cross-platform routing (Discord + WhatsApp + Telegram for different family members), persistent memory, webhook ingestion from smart home/IoT, and proactive cron-based automation — all without writing deployment infrastructure. Agno serves as the agent intelligence layer, either behind OpenClaw (via skills/plugins or as an API backend) or alongside it (as a direct DiscordClient for specialized agents that need full Python tool access).

The combination covers roughly **85–90% of Discord's feature surface** through bot APIs. The missing 10–15% (screen sharing, Activities, Clips, profile cosmetics, message search) are either client-side luxuries or niche features that don't impact a Life OS workflow. The real constraints are operational: **cost management** (set hard budget caps, route to cheap models), **security hygiene** (audit skills, restrict Gateway access, backup configs), and **reliability** (use a VPS or always-on machine with process monitoring).

The four-pillar channel structure, priority-based notification routing, and multi-agent team architecture described above provide a production-ready blueprint. Start with the Command Center category and one pillar, validate the workflow, then expand. The forum channel pattern (one thread per workout, one thread per investment, one thread per project) creates a naturally searchable, tagged knowledge base that grows more valuable over time — and that the agent team can reference via Agno's knowledge/RAG system to provide increasingly personalized, context-aware assistance.

---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
