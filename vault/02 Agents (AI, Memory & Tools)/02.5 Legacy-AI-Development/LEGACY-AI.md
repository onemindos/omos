---
id: root-00-20260113225954
title: "LEGACY AI"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "LEGACY AI"
---
# Legacy AI Documentation

> **Last Updated:** 2026-01-07

Legacy is Zeus's unified AI assistant with voice and home automation capabilities.

---

## Architecture Overview

```
Phone Call → Twilio → Caller ID Filter → Vapi → Legacy-AI (Letta)
                           │                           │
                      (Zeus only)                   Tools:
                           │                    • home_assistant()
                           ▼
                     Other callers
                     → Rejected
```

---

## Components

| Component | Service | Status |
|-----------|---------|--------|
| Voice Interface | Vapi | Active |
| Phone Number | Twilio | +1 (203) 583-3480 |
| AI Brain | Letta Cloud | Active |
| Home Control | Home Assistant | Active |
| Notifications | Telegram | Planned |

---

## Security Model

**Caller ID Restriction** - Phone-level security

- Only Zeus's phone (+1 203-303-4844) can connect
- All other callers are silently rejected
- No verbal authentication required
- Filter runs on Twilio serverless function

---

## IDs Reference

```
LEGACY-AI AGENT:    agent-da585f6f-189f-472a-974c-29b9578f9b55
VAPI ASSISTANT:     cb290108-e8b7-4e47-b4e3-e003ca822308
VAPI PHONE ID:      d2a1ddd0-46df-4baf-9639-7ef529a768ae
TWILIO PHONE:       +1 (203) 583-3480
ZEUS PHONE:         +1 (203) 303-4844
```

---

## Capabilities

### Voice Calls
- Inbound calls via Twilio → Vapi → Letta
- Speech-to-text, AI processing, text-to-speech
- Natural conversation with memory

### Home Assistant Control
Tool: `home_assistant(action, ...)`

| Action | Description | Example |
|--------|-------------|---------|
| `speak` | TTS announcement | `home_assistant(action="speak", message="Hello")` |
| `summary` | Get location status | `home_assistant(action="summary", location="grw")` |
| `control` | Control devices | `home_assistant(action="control", entity="light.office", state="on")` |

Locations: `grw` (Greenwich), `nyc` (NYC), `mex` (Mexico), `ups` (Upstate)

### Agent Delegation
Legacy can delegate to specialist agents:

| Agent | Purpose | ID |
|-------|---------|-----|
| Generational Entrepreneurship | Business tasks | `agent-e5d22e27-674e-4a44-b0b1-03e4b4f89408` |
| Holistic Performance | Personal tasks | `agent-0d5f61ce-1290-4cdc-9517-ecf353309791` |

---

## Communication Modes

**FULL COMMS** (default)
- Respond to all speech normally

**STANDBY MODE**
- Only respond to commands starting with "Legacy"
- Other speech gets: "Standing by."

Commands:
- "Activate standby mode" → Enter standby
- "Full comms" → Exit standby

---

## Personality

- **Direct** - Tell it like it is, not diplomatic
- **Military brevity** - 1-2 sentences on voice
- **Forward-thinking** - Always consider next steps
- **Lead with the answer** - Get to the point

Address: "Mr. Delacruz"
Task completion: "Done, Mr. Delacruz" or "Executed."

---

## Cloud Services

### Twilio
- Phone number: +1 (203) 583-3480
- Voice filter: `legacy-voice-filter-2815-dev.twil.io/voice-filter`
- Status callback: Vapi webhook

### Vapi
- Assistant handles STT/TTS
- Connects to Letta for AI processing
- First message: "Legacy online. How can I help, Mr. Delacruz?"

### Letta Cloud
- Hosts Legacy-AI agent
- Persistent memory across conversations
- Tool execution for home automation

### Home Assistant
- Nabu Casa cloud connection
- Controls smart home devices
- TTS via cloud service

---

## Removed Features

| Feature | Reason | Alternative |
|---------|--------|-------------|
| SMS | A2P 10DLC compliance required | Telegram |
| Call Sign Auth | Unnecessary with caller ID | Caller ID filter |

---

## Testing

### Voice Call Test
1. Call +1 (203) 583-3480 from Zeus's phone
2. Legacy answers: "Legacy online. How can I help, Mr. Delacruz?"
3. Try: "Turn on the office lights"

### API Test
```bash
curl -X POST https://api.letta.com/v1/agents/agent-da585f6f-189f-472a-974c-29b9578f9b55/messages \
  -H "Authorization: Bearer $LETTA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"messages": [{"role": "user", "content": "Status report"}]}'
```

---

## Related Files

- [LEGACY-PLAN.md](../LEGACY-PLAN.md) - Master plan and session logs
- [LEGACY-COMMS.md](../LEGACY-COMMS.md) - Deprecated communications docs
- [LEGACY-HA-REFERENCE.md](../LEGACY-HA-REFERENCE.md) - Home Assistant reference
