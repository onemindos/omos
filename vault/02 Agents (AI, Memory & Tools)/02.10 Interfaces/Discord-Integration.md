---
id: ui-02-20260111053925
title: "Discord Integration - OneMind Command Center"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Discord Integration - OneMind Command Center"
---

# Discord Integration - OneMind Command Center

> Primary interface for OneMind operations via Discord.
> Status: Active Development

---

## Credentials

| Item | Value |
|------|-------|
| Bot Token | `MTQ1MzQyNjkzNTQ0MjA0NzE5MQ.G2NkpD.6XSiJpxf7WXqgfpgBA1sh0MDU-MIAAx9PGcfic` |
| Client/App ID | `1453426935442047191` |
| Server ID | `1453411552181883027` |
| Owner User ID | `1024395579742228570` |
| Bot Name | OneMind OS#9821 |

**OAuth2 URL:**
```
https://discord.com/oauth2/authorize?client_id=1453426935442047191&permissions=2150746176&integration_type=0&scope=bot+applications.commands
```

---

## Server Structure

```
OneMind Server
├── COMMAND CENTER (category)
│   ├── #legacy-ai                        # Main chat with Legacy-AI
│   ├── #approvals                        # HITL approval requests
│   └── #alerts                           # System notifications
├── PILLARS (category)
│   ├── #holistic-performance-ai          # HP agent (green)
│   ├── #legacy-evolution-ai              # LE agent (blue)
│   ├── #generational-entrepreneurship-ai # GE agent (orange)
│   └── #innovative-technology-ai         # IT agent (red)
├── LOGS (category)
│   └── #agent-activity                   # Agent action logs
└── VOICE (category) [Future]
    └── voice-command                     # Voice channel
```

---

## Agent Configuration

| Agent | ID | Color | Channel |
|-------|-----|-------|---------|
| Legacy-AI | `agent-da585f6f-189f-472a-974c-29b9578f9b55` | Purple | #legacy-ai |
| Holistic-Performance-AI | `agent-1e1155fa-de36-475f-90c7-59e1761cb063` | Green | #holistic-performance-ai |
| Legacy-Evolution-AI | `agent-3ef12371-a11f-48c9-93f3-109dcf87e224` | Blue | #legacy-evolution-ai |
| Generational-Entrepreneurship-AI | `agent-052cf675-bcb6-468e-bbec-79e5be465c7a` | Orange | #generational-entrepreneurship-ai |
| Innovative-Technology-AI | `agent-2655ec04-d529-478e-b2b2-56b63b6f5f37` | Red | #innovative-technology-ai |

---

## Available Commands

| Command | Description |
|---------|-------------|
| `/agent` | Show current agent info |
| `/agents` | List all available agents |
| `/switch <name>` | Switch to different agent |
| `/status` | Check system connection status |
| `/memory` | View agent memory blocks |
| `/tools` | List agent's attached tools |
| `/reasoning on\|off` | Toggle reasoning display |
| `/help` | Show all commands |

---

## How to Use

### Chat with Agents

Simply type in the appropriate channel:
- Type in `#legacy-ai` → talks to Legacy-AI
- Type in `#holistic-performance-ai` → talks to HP agent
- Type in `#legacy-evolution-ai` → talks to LE agent
- Type in `#generational-entrepreneurship-ai` → talks to GE agent
- Type in `#innovative-technology-ai` → talks to IT agent

Responses are formatted as rich embeds with:
- Agent name and color
- Response content
- Reasoning (can be toggled with `/reasoning off`)
- Timestamp

### Switch Agents

Use `/switch <name>` where name is:
- `legacy` or `Legacy-AI`
- `hp` or `HP-Conductor`
- `le` or `LE-Conductor`
- `ge` or `GE-Conductor`
- `it` or `IT-Conductor`

### View Agent Memory

Use `/memory` to see the current agent's memory blocks.

### Check Status

Use `/status` to see:
- Letta API connection status
- Todoist API connection status
- Current active agent
- Server information

---

## Running the Bot

### Local Development

```bash
cd "00-09 Command (Meta)/03 Legacy-AI (Co-Pilot)/Discord-Bot"
pip install -r requirements.txt
python bot.py
```

### Production

Use PM2 or systemd to keep the bot running:

```bash
pm2 start bot.py --name onemind-bot --interpreter python3
pm2 save
```

---

## Development Roadmap

### Phase 1: Foundation (Complete)
- [x] Bot setup and connection
- [x] Channel auto-creation
- [x] Text chat with agents
- [x] Basic slash commands
- [x] Rich embed responses

### Phase 2: HITL & Notifications
- [ ] Approval buttons in #approvals
- [ ] Alert notifications
- [ ] Daily summaries
- [ ] Letta tools for Discord

### Phase 3: Advanced Features
- [ ] Todoist task display
- [ ] Memory editing
- [ ] Tool management
- [ ] Agent creation

### Phase 4: Voice (Future)
- [ ] Voice channel support
- [ ] Deepgram STT integration
- [ ] Cartesia TTS integration
- [ ] Conversational mode

---

## Code Location

```
00-09 Command (Meta)/
└── 03 Legacy-AI (Co-Pilot)/
    └── Discord-Bot/
        ├── bot.py              # Main entry point
        ├── config.py           # Configuration
        ├── .env                # Secrets
        ├── requirements.txt    # Dependencies
        ├── services/
        │   └── letta_client.py # Letta API client
        ├── utils/
        │   └── embeds.py       # Embed builders
        └── README.md           # Setup docs
```

---

## Related

- [[Telegram-Integration]] - Text-only backup interface
- [[LiveKit-Voice-Agent]] - Future voice interface
- [[00-09 Command (Meta)/03 Legacy-AI (Co-Pilot)/01.02 Letta Agents/]] - Agent configurations
- [[Task-Engine - Todoist Sync/06-Agent-Integration]] - Task integration
