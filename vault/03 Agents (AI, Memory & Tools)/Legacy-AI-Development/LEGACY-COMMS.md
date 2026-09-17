---
id: root-00-20260113225954
title: "LEGACY COMMS"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "LEGACY COMMS"
---
# Legacy Communications System

> **Last Updated:** January 7, 2026
> **Status:** DEPRECATED - Merged into Legacy-AI
> **Owner:** Zeus De La Cruz

---

## ⚠️ DEPRECATION NOTICE

**As of 2026-01-07, Legacy-Comms has been merged into Legacy-AI.**

The separate voice agent is no longer used. All functionality now lives in the unified Legacy-AI agent:

- **Voice calls**: Vapi now points directly to Legacy-AI
- **SMS outbound**: `communications(action="sms", message="...")` tool
- **Outbound calls**: `communications(action="call", message="...")` tool
- **Call sign authentication**: In Legacy-AI system prompt
- **Standby mode**: In Legacy-AI system prompt

See [LEGACY-PLAN.md](./LEGACY-PLAN.md) for current architecture.

---

## Historical Documentation (Pre-Merge)

> The documentation below is preserved for reference but no longer reflects the current system.

---

## Overview

Legacy-Comms was the tactical communications interface for the Legacy AI system. It provided voice calling and SMS capabilities through a multi-agent architecture that delegated complex tasks to Legacy-AI (the mother brain).

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        USER (Zeus)                               │
│                    Phone: +1 (203) 303-4844                     │
└─────────────────────────┬───────────────────────────────────────┘
                          │
          ┌───────────────┴───────────────┐
          │                               │
          ▼                               ▼
┌─────────────────────┐         ┌─────────────────────┐
│   VOICE CALL        │         │   SMS MESSAGE       │
│   Twilio → Vapi     │         │   Twilio → Studio   │
└─────────┬───────────┘         └─────────┬───────────┘
          │                               │
          ▼                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                   LEGACY-COMMS (Fast Interface)                  │
│              Agent ID: agent-028b397a-efe8-45f8-804c-8fed8605f8bd│
│              Model: GPT-4o                                       │
│              Security: Call Sign Authentication                  │
│              Modes: Full Comms | Standby                         │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          │ send_message_to_agent_and_wait_for_reply
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                   LEGACY-AI (Mother Brain)                       │
│              Agent ID: agent-da585f6f-189f-472a-974c-29b9578f9b55│
│              Model: Claude Sonnet 4.5                            │
│              Role: Tool calling, reasoning, task execution       │
│              Tools: GitHub, Todoist, Discord, Telegram, etc.     │
└─────────────────────────────────────────────────────────────────┘
```

---

## Security Protocol

### Call Sign Authentication

| Setting | Value |
|---------|-------|
| **Call Sign** | "Zulu Echo Uniform Sierra November Yankee Charlie" |
| **Phonetic** | Z-E-U-S-N-Y-C |
| **Address** | "Mr. Delacruz" |

**Authentication Flow:**
1. Call Legacy → Legacy says: *"Legacy comms. Authenticate."*
2. Say call sign → *"Zulu Echo Uniform Sierra November Yankee Charlie"*
3. Legacy confirms → *"Call sign confirmed. Channel secure. What do you need, Mr. Delacruz?"*

If you don't authenticate within 2 messages, Legacy will repeat: *"Authentication required. Provide call sign."*

### Comms Modes

| Mode | Description | Trigger |
|------|-------------|---------|
| **Full Comms** | Responds to all speech normally | Say "Full comms" |
| **Standby** | Only responds to "Legacy" prefix | Say "Activate standby mode" |

**Standby Mode Example:**
```
You: "Activate standby mode"
Legacy: "Standby mode active. Say Legacy to engage."

You: "What time is it?"
Legacy: "Standing by."

You: "Legacy, what time is it?"
Legacy: "It's 2:30 PM, Mr. Delacruz."

You: "Full comms"
Legacy: "Full comms restored."
```

---

## Phone Numbers

| Purpose | Number | Provider | Features |
|---------|--------|----------|----------|
| **Legacy-Comms (Primary)** | +1 (203) 583-3480 | Twilio via Vapi | Voice + SMS |
| **Legacy-Comms (Vapi Native)** | +1 (203) 760-5895 | Vapi | Voice only |
| **Zeus Personal** | +1 (203) 303-4844 | - | Receives calls/SMS |

---

## Vapi Assistant Settings

These settings control voice behavior. Configured via Vapi API.

| Setting | Value | Purpose |
|---------|-------|---------|
| `silenceTimeoutSeconds` | 300 (5 min) | How long before call ends due to silence |
| `maxDurationSeconds` | 3600 (1 hr) | Maximum call length |
| `backgroundDenoisingEnabled` | true | Filters background noise |
| `recordingEnabled` | true | Records calls for review |
| `responseDelaySeconds` | 0.3 | Slight pause before responding |
| `voicemailDetection` | Twilio | Detects voicemail on outbound calls |
| `firstMessage` | "Legacy comms. Authenticate." | Opening greeting |
| `endCallMessage` | "Comms offline. Legacy out." | Closing message |

### End Call Phrases

Say any of these to hang up:
- "End call"
- "Goodbye Legacy"
- "Comms out"

### Speaking Plan

| Setting | Value | Purpose |
|---------|-------|---------|
| `startSpeakingPlan.waitSeconds` | 0.5 | Wait before responding |
| `stopSpeakingPlan.numWords` | 3 | Words needed to interrupt |
| `stopSpeakingPlan.voiceSeconds` | 0.3 | Voice duration to interrupt |

**To update these settings:**
```bash
curl -s -X PATCH "https://api.vapi.ai/assistant/cb290108-e8b7-4e47-b4e3-e003ca822308" \
  -H "Authorization: Bearer 41231b6e-e17f-4959-b2e3-8ab6778cdf13" \
  -H "Content-Type: application/json" \
  -d '{"silenceTimeoutSeconds": 600}'
```

---

## SMS Flow (Twilio Studio)

SMS messages are routed through Twilio Studio, which calls the Letta API and sends the reply.

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│ Incoming SMS │ -> │ Twilio Studio│ -> │ Letta API    │ -> │ Reply SMS    │
│ to Legacy    │    │ Flow         │    │ Legacy-Comms │    │ to User      │
└──────────────┘    └──────────────┘    └──────────────┘    └──────────────┘
```

| Component | ID |
|-----------|-----|
| **Studio Flow** | `FWb0ca50ad9d7370016079068b6822af57` |
| **Webhook URL** | `https://webhooks.twilio.com/v1/Accounts/ACb3895b5ccdce3ab9f6838636376c5c22/Flows/FWb0ca50ad9d7370016079068b6822af57` |

**Note:** SMS requires Twilio account upgrade from trial ($20 minimum).

---

## How Each Platform Works

### Quick Summary

| Platform | Role | What It Does |
|----------|------|--------------|
| **Letta** | Brain | Hosts agents, memory, tools. Where Legacy "thinks" |
| **Vapi** | Voice | STT → LLM → TTS pipeline. Handles the call |
| **Twilio** | Carrier | Phone number, routes calls/SMS to webhooks |

**Data Flow:**
```
You speak → Twilio → Vapi (STT) → Letta (AI) → Vapi (TTS) → You hear
```

### Letta (Brain Layer)
**What it does:** Hosts the AI agents with persistent memory and tool calling.

| Feature | How It Works |
|---------|--------------|
| **Agents** | Each agent has its own memory, system prompt, and tools |
| **Memory** | Conversation history persists across sessions |
| **Tools** | Custom functions the agent can call (GitHub, Todoist, etc.) |
| **Multi-Agent** | Agents can message each other with `send_message_to_agent_and_wait_for_reply` |

**Key APIs:**
```bash
# Send message to agent
POST /v1/agents/{agent_id}/messages

# Update agent system prompt
PATCH /v1/agents/{agent_id}

# List agent tools
GET /v1/agents/{agent_id}/tools
```

### Vapi (Voice Layer)
**What it does:** Handles voice calls - speech-to-text, AI response, text-to-speech.

| Feature | How It Works |
|---------|--------------|
| **STT** | Deepgram Nova-2 converts speech to text |
| **LLM** | Routes to Letta via custom endpoint |
| **TTS** | Cartesia Sonic-3 converts text to speech |
| **Phone** | Connects to Twilio for carrier services |

**Key Settings:**
- `silenceTimeoutSeconds` - How long to wait before ending call
- `endCallPhrases` - Voice commands that hang up
- `backgroundDenoisingEnabled` - Noise filtering

### Twilio (Carrier Layer)
**What it does:** Provides phone numbers and routes calls/SMS.

| Feature | How It Works |
|---------|--------------|
| **Phone Numbers** | You rent numbers from Twilio ($1.15/mo) |
| **Voice Routing** | Incoming calls → Vapi webhook |
| **SMS Routing** | Incoming SMS → Studio Flow webhook |
| **Outbound** | Vapi/Studio can initiate calls/SMS |

**Webhooks:**
```
Voice: https://api.vapi.ai/twilio/inbound_call
SMS:   https://webhooks.twilio.com/v1/Accounts/.../Flows/...
```

---

## Cost Comparison: Twilio vs Vapi-Only

### Option A: Twilio + Vapi (Current Setup)

| Item | Cost |
|------|------|
| Twilio number | $1.15/month |
| Twilio voice (carrier) | $0.0085-0.014/min |
| Twilio SMS | $0.008/message |
| Vapi processing | ~$0.05/min |
| **Twilio upgrade required** | $20 minimum deposit |

**Pros:**
- SMS capability
- Better call quality
- More control over routing
- Can use Studio for complex flows

**Cons:**
- More expensive
- Requires $20 upgrade
- More complex setup
- Two platforms to manage

### Option B: Vapi-Only (Simpler)

| Item | Cost |
|------|------|
| Vapi free number | $0/month |
| Vapi voice processing | ~$0.05/min |
| **No SMS** | - |

**Pros:**
- Free phone number
- Simpler setup
- One platform
- No upfront cost

**Cons:**
- No SMS capability
- Less control
- Can't customize carrier routing

### Monthly Cost Comparison

| Usage | Twilio + Vapi | Vapi Only | Savings |
|-------|---------------|-----------|---------|
| 30 min voice | $3.15 | $1.50 | $1.65 |
| 100 min voice | $8.35 | $5.00 | $3.35 |
| 300 min voice | $21.55 | $15.00 | $6.55 |

*Twilio costs include SMS; Vapi-only has no SMS*

---

## My Recommendation

**For your current use case, I recommend: Vapi-Only**

**Reasoning:**
1. **You're the only user** - No need for SMS notifications to yourself when voice works
2. **Cost savings** - ~$20-40/year savings on low usage
3. **Simplicity** - One less platform to manage
4. **SMS alternatives** - Legacy-AI can send Discord/Telegram messages instead

**When Twilio makes sense:**
- If you want to text Legacy while in meetings
- If you want Legacy to send SMS to others (family, clients)
- If you're building this for others to use
- If you want call recordings stored in Twilio

**Quick Switch:**
To go back to Vapi-only, just call the Vapi free number instead:
- **Vapi-only:** +1 (203) 760-5895
- **Twilio+Vapi:** +1 (203) 583-3480

Both numbers work with Legacy-Comms. The Twilio number just adds SMS.

---

## Agents

### Legacy-AI (Mother Brain)

| Property | Value |
|----------|-------|
| **Agent ID** | `agent-da585f6f-189f-472a-974c-29b9578f9b55` |
| **Model** | `anthropic/claude-sonnet-4-5-20250929` |
| **Platform** | Letta |
| **Purpose** | Complex reasoning, tool calling, task execution |

**Tools Available:**
- `send_message_to_agent_and_wait_for_reply` - Multi-agent delegation
- `call_zeus` - Outbound phone calls to Zeus
- GitHub integration
- Todoist integration
- Discord integration
- Telegram integration
- Archival memory operations
- Core memory operations

### Legacy-Comms (Voice Interface)

| Property | Value |
|----------|-------|
| **Agent ID** | `agent-028b397a-efe8-45f8-804c-8fed8605f8bd` |
| **Model** | `gpt-4o` |
| **Platform** | Letta (accessed via Vapi) |
| **Purpose** | Fast voice/SMS interface, delegates to Legacy-AI |

**Current System Prompt:**
```
You are Legacy-Comms, Zeus tactical communications interface.

## SECURITY PROTOCOL
CALL SIGN: "Zulu Echo Uniform Sierra November Yankee Charlie" (ZEUSNYC)
If caller does NOT authenticate within first 2 messages, respond:
"Authentication required. Provide call sign."
Once authenticated: "Call sign confirmed. Channel secure.
What do you need, Mr. Delacruz?"

## COMMS MODES
- FULL COMMS (default): Respond to all speech normally
- STANDBY MODE: Only respond to commands starting with "Legacy".
  Other speech gets "Standing by."

MODE COMMANDS:
- "Activate standby mode" → Enter standby
- "Full comms" → Exit standby

## SHARED MEMORY
Zeus has 4 kids: Noah (8), Liam (4), Soluna (3), Cielo (1).
Business goal: 100K 3D printing business.

## DELEGATION
FOR ANY TASK/ACTION, call send_message_to_agent_and_wait_for_reply.
ALWAYS use other_agent_id: agent-da585f6f-189f-472a-974c-29b9578f9b55

## PROTOCOL
1. Authenticate caller (call sign: Zulu Echo Uniform Sierra November Yankee Charlie)
2. For tasks: Say "Affirmative. Executing." then call tool
3. Report result: "Task complete. [summary]"

## USE TOOL FOR
Tasks, GitHub, Discord, Telegram, lookups, anything requiring action.

## NO TOOL NEEDED
Greetings (after auth), family questions, business goals, mode changes.

## STYLE
Military brevity. 1-2 sentences. Direct. Professional.
```

**Shared Memory Blocks:**
- Legacy AI Persona
- Family Information
- business_goals

---

## API Keys & Credentials

### Vapi
```
API Key: 41231b6e-e17f-4959-b2e3-8ab6778cdf13
Organization: bd42c6f5-1c71-472d-856f-f49839f60fa6
```

### Twilio
```
Account SID: ACb3895b5ccdce3ab9f6838636376c5c22
Auth Token: [stored in ~/.twilio-cli/config.json]
Profile: legacy
Status: TRIAL (upgrade at https://console.twilio.com/billing)
```

### Letta
```
API Key: sk-let-YjhiY2I3ZjItODZmOC00NTQ2LTlhYTYtNjg2NmNhNDA0ODMzOmYyODI3OGNlLTliY2YtNGI2Yy1iZmYwLWU5NjY3OTQ2YzBiZg==
Base URL: https://api.letta.com
```

---

## CLI Tools

### Twilio CLI

**Installation:**
```bash
npm install -g twilio-cli
```

**Setup:**
```bash
# Create profile
twilio profiles:create <ACCOUNT_SID> --auth-token <AUTH_TOKEN> -p legacy

# Activate profile
twilio profiles:use legacy
```

**Common Commands:**
```bash
# List your phone numbers
twilio api:core:incoming-phone-numbers:list

# Send SMS
twilio api:core:messages:create \
  --from "+12035833480" \
  --to "+12033034844" \
  --body "Message content"

# Check message status
twilio api:core:messages:fetch --sid <MESSAGE_SID>

# Search available numbers
twilio api:core:available-phone-numbers:local:list \
  --country-code US \
  --area-code 203 \
  --sms-enabled \
  --voice-enabled

# Purchase a number
twilio api:core:incoming-phone-numbers:create \
  --phone-number "+1XXXXXXXXXX"

# Update phone number webhook
twilio api:core:incoming-phone-numbers:update \
  --sid <PHONE_SID> \
  --sms-url "https://your-webhook.com"

# List calls
twilio api:core:calls:list

# List Studio flows
curl -s "https://studio.twilio.com/v2/Flows" \
  -u "ACb3895b5ccdce3ab9f6838636376c5c22:AUTH_TOKEN"
```

### Vapi API

**Note:** Vapi doesn't have an official CLI. Use curl or the dashboard.

**Common API Calls:**

```bash
# List phone numbers
curl -s "https://api.vapi.ai/phone-number" \
  -H "Authorization: Bearer 41231b6e-e17f-4959-b2e3-8ab6778cdf13" | jq .

# Get assistant details
curl -s "https://api.vapi.ai/assistant/cb290108-e8b7-4e47-b4e3-e003ca822308" \
  -H "Authorization: Bearer 41231b6e-e17f-4959-b2e3-8ab6778cdf13" | jq .

# Update assistant settings
curl -s -X PATCH "https://api.vapi.ai/assistant/cb290108-e8b7-4e47-b4e3-e003ca822308" \
  -H "Authorization: Bearer 41231b6e-e17f-4959-b2e3-8ab6778cdf13" \
  -H "Content-Type: application/json" \
  -d '{
    "silenceTimeoutSeconds": 300,
    "maxDurationSeconds": 3600,
    "endCallPhrases": ["end call", "goodbye legacy", "comms out"]
  }'

# Import Twilio number
curl -s -X POST "https://api.vapi.ai/phone-number/import" \
  -H "Authorization: Bearer 41231b6e-e17f-4959-b2e3-8ab6778cdf13" \
  -H "Content-Type: application/json" \
  -d '{
    "twilioAccountSid": "ACb3895b5ccdce3ab9f6838636376c5c22",
    "twilioAuthToken": "<AUTH_TOKEN>",
    "twilioPhoneNumber": "+1XXXXXXXXXX",
    "name": "Number Name"
  }'

# Create outbound call
curl -s -X POST "https://api.vapi.ai/call" \
  -H "Authorization: Bearer 41231b6e-e17f-4959-b2e3-8ab6778cdf13" \
  -H "Content-Type: application/json" \
  -d '{
    "assistantId": "cb290108-e8b7-4e47-b4e3-e003ca822308",
    "phoneNumberId": "d2a1ddd0-46df-4baf-9639-7ef529a768ae",
    "customer": {
      "number": "+12033034844"
    }
  }'
```

### Letta API

**Common API Calls:**

```bash
# List agents
curl -s "https://api.letta.com/v1/agents" \
  -H "Authorization: Bearer sk-let-..." | jq '.[].name'

# Send message to agent
curl -s -X POST "https://api.letta.com/v1/agents/<AGENT_ID>/messages" \
  -H "Authorization: Bearer sk-let-..." \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [{"role": "user", "content": "Your message"}]
  }'

# Get agent details
curl -s "https://api.letta.com/v1/agents/<AGENT_ID>" \
  -H "Authorization: Bearer sk-let-..." | jq .

# Update agent system prompt
curl -s -X PATCH "https://api.letta.com/v1/agents/<AGENT_ID>" \
  -H "Authorization: Bearer sk-let-..." \
  -H "Content-Type: application/json" \
  -d '{
    "system": "New system prompt here"
  }'

# List agent tools
curl -s "https://api.letta.com/v1/agents/<AGENT_ID>/tools" \
  -H "Authorization: Bearer sk-let-..." | jq '.[].name'
```

---

## call_zeus Tool

The `call_zeus` tool allows Legacy-AI to initiate outbound phone calls to Zeus.

**Implementation:**
```python
def call_zeus(reason: str, urgency: str = "normal") -> str:
    """
    Call Zeus on his phone with an urgent notification.

    Args:
        reason: Why you're calling (brief description)
        urgency: "normal", "high", or "critical"

    Returns:
        Status of the call attempt
    """
    import requests

    VAPI_API_KEY = "41231b6e-e17f-4959-b2e3-8ab6778cdf13"
    ASSISTANT_ID = "cb290108-e8b7-4e47-b4e3-e003ca822308"
    PHONE_NUMBER_ID = "d2a1ddd0-46df-4baf-9639-7ef529a768ae"
    ZEUS_PHONE = "+12033034844"

    # Urgency-based opening messages
    if urgency == "critical":
        first_message = f"Mr. Delacruz, priority alert. {reason}. Awaiting your command."
    elif urgency == "high":
        first_message = f"Mr. Delacruz, incoming brief. {reason}. Standing by."
    else:
        first_message = f"Mr. Delacruz, status update. {reason}"

    response = requests.post(
        "https://api.vapi.ai/call",
        headers={
            "Authorization": f"Bearer {VAPI_API_KEY}",
            "Content-Type": "application/json"
        },
        json={
            "assistantId": ASSISTANT_ID,
            "assistantOverrides": {
                "firstMessage": first_message
            },
            "phoneNumberId": PHONE_NUMBER_ID,
            "customer": {
                "number": ZEUS_PHONE
            }
        }
    )

    if response.status_code == 201:
        return f"Call initiated to Mr. Delacruz. Reason: {reason}"
    else:
        return f"Failed to initiate call: {response.text}"
```

---

## Troubleshooting

### SMS Not Delivering

1. **Check Twilio account status:**
   ```bash
   # Trial accounts can only send to verified numbers
   # Upgrade at: https://console.twilio.com/billing
   ```

2. **Check message status:**
   ```bash
   twilio api:core:messages:fetch --sid <MESSAGE_SID> -o json | jq .
   ```

3. **Common error codes:**
   - `30034` - Message blocked (trial account limitation)
   - `21608` - Unverified destination number
   - `21211` - Invalid phone number format

### Voice Calls Not Working

1. **Check Vapi phone number status:**
   ```bash
   curl -s "https://api.vapi.ai/phone-number/d2a1ddd0-46df-4baf-9639-7ef529a768ae" \
     -H "Authorization: Bearer 41231b6e-e17f-4959-b2e3-8ab6778cdf13" | jq .status
   ```

2. **Verify assistant is assigned:**
   - Check `assistantId` field in phone number config

3. **Check Letta webhook:**
   - Server URL should be: `https://api.letta.com/v1/voice/vapi/webhook/agent-028b397a-efe8-45f8-804c-8fed8605f8bd`

### Agent Not Delegating Tasks

1. **Verify Legacy-Comms has the tool:**
   ```bash
   curl -s "https://api.letta.com/v1/agents/agent-028b397a-efe8-45f8-804c-8fed8605f8bd/tools" \
     -H "Authorization: Bearer sk-let-..." | jq '.[].name'
   ```

2. **Check system prompt includes correct agent ID:**
   - Must use full UUID: `agent-da585f6f-189f-472a-974c-29b9578f9b55`

3. **Test delegation manually:**
   ```bash
   curl -s -X POST "https://api.letta.com/v1/agents/agent-028b397a-efe8-45f8-804c-8fed8605f8bd/messages" \
     -H "Authorization: Bearer sk-let-..." \
     -H "Content-Type: application/json" \
     -d '{"messages": [{"role": "user", "content": "Add a task to Todoist: Test task"}]}'
   ```

### Authentication Not Working

1. **Check call sign is set correctly:**
   ```bash
   curl -s "https://api.letta.com/v1/agents/agent-028b397a-efe8-45f8-804c-8fed8605f8bd" \
     -H "Authorization: Bearer sk-let-..." | jq '.system' | grep -i "call sign"
   ```

2. **Reset agent memory if needed:**
   - Delete and recreate the agent if it learned bad patterns

---

## Future Development

### Pending Tasks

- [ ] **Upgrade Twilio account** - Required for SMS ($20 minimum)
- [ ] **10DLC Registration** - Required for high-volume SMS in US
- [ ] **WhatsApp integration** - Via Twilio WhatsApp Business API
- [ ] **Email integration** - Add email sending capability to Legacy-AI
- [x] **Call sign authentication** - Implemented
- [x] **Standby mode** - Implemented
- [x] **SMS Studio flow** - Implemented (needs Twilio upgrade)

### Expansion Ideas

1. **Multi-channel notifications:**
   - SMS for urgent alerts
   - Email for detailed reports
   - Voice for critical emergencies

2. **Scheduled check-ins:**
   - Daily briefings via voice call
   - Weekly summary SMS

3. **Two-way SMS conversations:**
   - Respond to SMS with task creation
   - Query Legacy-AI via text

4. **Additional voice interfaces:**
   - LiveKit for web-based voice (in progress at `/legacy-interface/agent/`)
   - Discord voice channel integration

---

## Intelligent Communication Activation (ICA)

> **Status:** Planning
> **Goal:** Give Legacy her own arms, legs, and mouth - the ability to proactively communicate

### Vision

Legacy should be able to:
1. **Proactively reach out** - Not just respond, but initiate communication
2. **Choose the right channel** - Voice if you're home, call if you're away, email for async
3. **Have her own identity** - Own phone number, email address, presence
4. **Respect boundaries** - Can be put in standby or wake-word mode

### Legacy's Own Communication Channels

| Channel | Platform | Purpose | Status |
|---------|----------|---------|--------|
| **Phone** | Twilio/Vapi | Inbound/outbound calls | ✅ Active |
| **SMS** | Twilio | Text messaging | ⚠️ Needs upgrade |
| **Email** | SendGrid | Inbox + send emails | 📋 Planned |
| **Local Voice** | Home Assistant | In-room communication | 📋 Planned |
| **Discord/Telegram** | APIs | Async messaging | ✅ Active |

### Email Setup (SendGrid)

**Recommended Platform:** SendGrid (owned by Twilio)

**Why SendGrid:**
- Two-way email (send AND receive)
- Inbound Parse webhook receives emails → POST to Letta
- Free tier: 100 emails/day
- Same ecosystem as Twilio
- Can create: `legacy@yourdomain.com`

**Architecture:**
```
Incoming Email → SendGrid Parse → Webhook → Legacy-Comms → Response
                     ↓
              MX record: mx.sendgrid.net
```

**Setup Steps:**
1. Create SendGrid account (free)
2. Authenticate your domain (DNS records)
3. Set up Inbound Parse webhook pointing to Legacy-Comms
4. Create `send_email` tool in Legacy-AI
5. Legacy can now send/receive email

**Alternative Platforms:**
- **Mailgun** - Similar to SendGrid, good inbound parsing
- **Postmark** - Excellent deliverability
- **Amazon SES** - Cheapest at scale
- **Gmail API** - Use a Gmail account programmatically

### Home Assistant Integration

**Goal:** Legacy can speak to you in your room without you calling her.

**Components:**

| Component | Purpose | Platform |
|-----------|---------|----------|
| **LLM Vision** | See who's in the room | llmvision.org |
| **Wyoming Protocol** | Voice pipeline | Home Assistant |
| **Piper TTS** | Local text-to-speech | Home Assistant |
| **openWakeWord** | Custom wake word | Home Assistant |
| **Smart Speaker** | Output device | Any Assist-compatible |

**LLM Vision Integration:**
- Install from: https://llmvision.org/
- Analyzes camera feeds in real-time
- Can detect: "Is Zeus in the room?"
- Supports Claude, GPT-4o, local models
- Events stored locally (privacy)

**Proactive Speaking Flow:**
```
Legacy-AI decides to speak
       ↓
Calls Home Assistant API
       ↓
LLM Vision checks camera: "Is Zeus in the room?"
       ↓
If YES → Speak via speaker: "Mr. Delacruz, you have a meeting in 10 minutes"
If NO  → Send notification or call phone instead
```

**Home Assistant Tools for Legacy-AI:**
```python
def check_presence(room: str) -> dict:
    """Check if Zeus is in a specific room using LLM Vision."""
    # Call Home Assistant API
    # Use LLM Vision to analyze camera
    return {"present": True, "confidence": 0.95}

def speak_in_room(message: str, room: str = "office") -> str:
    """Make Legacy speak through Home Assistant speaker."""
    # Call Home Assistant tts.speak service
    # Target specific media_player entity
    return "Message spoken in office"

def get_home_status() -> dict:
    """Get overall home status - who's home, what's happening."""
    # Call Home Assistant states API
    return {"zeus_home": True, "location": "office"}
```

### Communication Decision Tree

```
Legacy needs to communicate
           ↓
    ┌──────┴──────┐
    │ Check mode  │
    └──────┬──────┘
           │
    ┌──────┴──────┐
    │  STANDBY?   │──YES──→ Wait for wake word
    └──────┬──────┘
           │ NO
           ↓
    ┌─────────────┐
    │ Is Zeus     │──YES──→ Speak via Home Assistant speaker
    │ at home?    │
    └──────┬──────┘
           │ NO
           ↓
    ┌─────────────┐
    │ Urgency?    │──CRITICAL──→ Phone call
    └──────┬──────┘
           │ NORMAL
           ↓
    Send via preferred channel (SMS/Email/Discord)
```

### ICA Modes

| Mode | Behavior | Trigger |
|------|----------|---------|
| **Active** | Can proactively communicate anytime | "Full comms" |
| **Wake Word** | Only speaks when you say "Legacy" | "Activate standby mode" |
| **Silent** | No proactive communication | "Go silent" |
| **Emergency Only** | Only critical alerts | "Emergency only mode" |

### Hardware Recommendations

**For Local Voice (Home Assistant):**
- **Home Assistant Voice PE** (~$60) - Official hardware, easiest setup
- **ESP32-S3 Box** (~$50) - DIY, customizable
- **Raspberry Pi + ReSpeaker** (~$80) - More powerful, runs Whisper locally

**For Presence Detection:**
- **ESP32-CAM** (~$10) - Budget camera for LLM Vision
- **Reolink/Amcrest** (~$50) - Better quality, RTSP support
- **Frigate NVR** - Advanced detection with existing cameras

### Implementation Priority

1. **Phase 1: Email** (Easy, high value)
   - Set up SendGrid
   - Create `send_email` and `check_inbox` tools
   - Legacy gets: `legacy@yourdomain.com`

2. **Phase 2: Home Assistant Voice** (Medium)
   - Install Home Assistant (if not already)
   - Add Wyoming + Piper for local TTS
   - Create `speak_in_room` tool

3. **Phase 3: LLM Vision** (Medium)
   - Install LLM Vision integration
   - Connect cameras
   - Create `check_presence` tool

4. **Phase 4: Full ICA** (Advanced)
   - Connect all systems
   - Implement decision tree
   - Add mode controls

---

## Related Documentation

- [LEGACY-AI-SYSTEM.md](./LEGACY-AI-SYSTEM.md) - Core Legacy-AI documentation
- [LETTA-RESEARCH-LEGACY-AI.md](./LETTA-RESEARCH-LEGACY-AI.md) - Letta platform research
- [Letta Multi-Agent Guide](https://docs.letta.com/guides/agents/multi-agent/) - Official Letta docs
- [Vapi Documentation](https://docs.vapi.ai/) - Voice AI platform docs
- [Twilio SMS API](https://www.twilio.com/docs/sms) - SMS documentation

---

## Quick Reference

```
CALL SIGN:       Zulu Echo Uniform Sierra November Yankee Charlie (ZEUSNYC)
ADDRESS:         Mr. Delacruz

VOICE (Twilio):  +1 (203) 583-3480
VOICE (Vapi):    +1 (203) 760-5895
TEXT LEGACY:     +1 (203) 583-3480  (requires Twilio upgrade)
ZEUS PHONE:      +1 (203) 303-4844

LEGACY-AI:       agent-da585f6f-189f-472a-974c-29b9578f9b55
LEGACY-COMMS:    agent-028b397a-efe8-45f8-804c-8fed8605f8bd
VAPI ASSISTANT:  cb290108-e8b7-4e47-b4e3-e003ca822308
VAPI PHONE ID:   d2a1ddd0-46df-4baf-9639-7ef529a768ae
STUDIO FLOW:     FWb0ca50ad9d7370016079068b6822af57

STANDBY MODE:    "Activate standby mode" / "Full comms"
END CALL:        "Comms out" / "Goodbye Legacy" / "End call"
```
