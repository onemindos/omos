---
id: root-00-20260113225954
title: "LEGACY PLAN"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "LEGACY PLAN"
---
# Legacy AI System - Master Plan

> **Last Updated:** 2026-01-07
> **Purpose:** Comprehensive roadmap for Legacy AI development. Read this first in new sessions.

---

## Quick Status Overview

| Component | Status | Notes |
|-----------|--------|-------|
| **Legacy-AI (Unified)** | Active | Single agent: voice + brain + home control |
| **Legacy-Comms** | MERGED | Functionality merged into Legacy-AI |
| **Voice Calls** | Active | +1 (203) 583-3480 (Twilio via Vapi) |
| **SMS Outbound** | Removed | A2P 10DLC compliance required - using Telegram instead |
| **Outbound Calls** | Active | Legacy can call Zeus via Vapi |
| **Telegram** | Planned | Notifications and alerts |
| **Home Assistant** | Active | v2025.12.5, HA Cloud connected, API working |
| **Email (SendGrid)** | Planned | Gives outbound + inbound, no visual inbox |

---

## Critical IDs Reference

```
ADDRESS:         Mr. Delacruz
SECURITY:        Caller ID restriction (no verbal auth needed)

VOICE (Vapi):    +1 (203) 760-5895
VOICE (Twilio):  +1 (203) 583-3480
ZEUS PHONE:      +1 (203) 303-4844

LEGACY-AI:       agent-da585f6f-189f-472a-974c-29b9578f9b55 (unified agent)
LEGACY-COMMS:    agent-028b397a-efe8-45f8-804c-8fed8605f8bd (deprecated - merged)
VAPI ASSISTANT:  cb290108-e8b7-4e47-b4e3-e003ca822308 (points to Legacy-AI)
VAPI PHONE ID:   d2a1ddd0-46df-4baf-9639-7ef529a768ae

TWILIO SID:      ACb3895b5ccdce3ab9f6838636376c5c22
TWILIO VOICE:    legacy-voice-filter-2815-dev.twil.io (caller ID filter - ACTIVE)
LETTA API KEY:   sk-let-YjhiY2I3ZjItODZmOC00NTQ2LTlhYTYtNjg2NmNhNDA0ODMzOmYyODI3OGNlLTliY2YtNGI2Yy1iZmYwLWU5NjY3OTQ2YzBiZg==
VAPI API KEY:    41231b6e-e17f-4959-b2e3-8ab6778cdf13

HOME ASSISTANT:
  NABU CASA URL:  https://1mwywqa9zvz1hhh9lofv5wytww4nswti.ui.nabu.casa
  LOCAL URL:      http://homeassistant.local:8123
  API TOKEN:      eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJiMzY3Zjc1MzQ4MWI0YWZlODBlOGFlODk0ZjE5YmQ2NCIsImlhdCI6MTc2NzczODk4MiwiZXhwIjoyMDgzMDk4OTgyfQ.C__KGAQFBHRUxB4j-Mg7uu4NnbQv2WntovyZLNzKXeE
  VERSION:        2025.12.5
```

---

## Your Hardware Inventory

| Device | Purpose | Status |
|--------|---------|--------|
| **Beelink Mini PC** | Home Assistant host | i3, 8GB RAM, 256GB SSD - Running HA |
| **ATOM Echo** | Voice satellite | ✅ Integrated - TTS working |
| **Yi Cameras** | Vision/presence detection | ⏸️ Requires yi-hack firmware (see below) |
| **SwitchBot Plugs** | Smart outlets | ⏸️ Requires SwitchBot Hub for cloud API |
| **SwitchBot Humidity Sensor** | Environment monitoring | ⏸️ Requires SwitchBot Hub for cloud API |

### Device Integration Notes

**Yi Cameras:**

- Require custom firmware ([yi-hack-Allwinner](https://github.com/roleoroleo/yi-hack-Allwinner))
- Supported models: Yi 1080p Home (9FUS, BFUS, SFUS, BFCN, IFUS), Yi Dome X
- Enables RTSP streaming, MQTT, FTP access
- No cloud-based HA integration available

**SwitchBot:**

- **Option 1: Bluetooth** - Local, real-time, requires HA within Bluetooth range
- **Option 2: Cloud API** (recommended for multi-location) - Requires SwitchBot Hub
  - Get API token/secret from SwitchBot app → Developer Options
  - 10-minute sensor update interval
  - 10,000 API calls/day limit

---

## Key Questions Answered

### 1. Can Letta be the LLM for Home Assistant?

**Short answer:** YES - 90% confirmed, with two integration methods.

**Letta Endpoint:** `https://api.letta.com/v1/chat/completions` (OpenAI-compatible)

#### Method A: MCP (Recommended)

Home Assistant has native MCP server support. Letta supports MCP clients.

```text
Legacy-AI (Letta) ←──MCP──→ Home Assistant MCP Server
                              ↓
                    All entities auto-exposed
                    All services available
                    Legacy knows everything she can control
```

**Why MCP is best:**
- HA automatically exposes all entities to Legacy
- No manual tool creation needed
- Standardized protocol (future-proof)
- Legacy instantly knows what she has access to

#### Method B: Extended OpenAI Conversation (Fallback)

Install from HACS, configure:

```yaml
# Integration config:
Base URL: https://api.letta.com/v1
API Key: sk-let-YjhiY2I3ZjIt...  (your Letta API key)
Model: agent-da585f6f-189f-472a-974c-29b9578f9b55  (Legacy-AI agent ID)
```

**Note:** Function calling (for device control) needs testing with Letta's endpoint.

#### Method C: Letta Tools + HA REST API (Most Reliable)

Create tools in Legacy-AI that call Home Assistant's REST API directly:

```python
def control_ha_device(entity_id: str, action: str) -> str:
    """Control any Home Assistant device."""
    # POST https://ha.local/api/services/{domain}/{service}

def get_ha_states() -> dict:
    """Get all device states."""
    # GET https://ha.local/api/states
```

**Same Agent:** Use Legacy-AI directly for HA - no separate agent needed (unlike Vapi).

**Node-RED/n8n:** Not required. Letta + HA is sufficient for most use cases.

---

### 2. Can Legacy get vision through Vapi?

**Short answer:** No. Vapi is voice-only (STT → LLM → TTS).

#### At Home (via Home Assistant)

| Method | How It Works | Best For |
|--------|--------------|----------|
| **LLM Vision (HA)** | Camera → LLM analyzes → Text to Letta | Presence detection, room awareness |
| **Frigate NVR** | Object detection → Events → Letta | "Person detected in office" |
| **Direct API** | Periodic camera snapshot → Claude Vision → Letta | On-demand "what do you see?" |

**Architecture:**

```text
Yi Camera → Home Assistant → LLM Vision Integration
                                    ↓
                        "Zeus is in the office"
                                    ↓
                            Legacy-AI (Letta)
                                    ↓
                    "Mr. Delacruz, you have a meeting in 10 minutes"
                                    ↓
                         ATOM Echo speaks (via HA TTS)
```

#### On the Move (Mobile Vision)

| Method | How It Works | Cost | Latency |
|--------|--------------|------|---------|
| **Telegram Bot** | Send photo → Bot uses Claude Vision → Responds | Free (API costs) | 2-5 sec |
| **iOS Shortcut** | Take photo → Run shortcut → Send to Letta | Free | 3-5 sec |
| **Custom App** | Phone camera → Vision API → Letta | Dev time | 1-3 sec |
| **Ray-Ban Meta** | "Hey Meta, send to Legacy" → Vision | $299 glasses | Real-time |

**Easiest: Telegram Bot with Vision**

Since you already have Discord/Telegram, add an `analyze_image` tool to Legacy-AI:

```text
You: [Send photo to Legacy on Telegram]
Legacy Bot: Receives image → Sends to Claude Vision API → Gets description
Legacy: "I see you're at a coffee shop. The menu shows..."
```

**Most Seamless: iOS Shortcut**

1. "Hey Siri, show Legacy"
2. Takes photo with back camera
3. Sends to Claude Vision API
4. Sends description + question to Letta
5. Speaks response via Siri

**Implementation:**

1. Install LLM Vision integration in Home Assistant
2. Connect Yi cameras via RTSP
3. Create `check_presence` tool in Letta that calls HA API
4. Create `analyze_image` tool for mobile vision via Telegram
5. Legacy can now "see" both at home and on the go

#### Vision Solutions Comparison (Researched 2026-01-07)

| Solution | What It Does | Pros | Cons | Cost |
|----------|--------------|------|------|------|
| **Frigate NVR** | Local AI object detection | 100% local/private, fast, low latency | Requires Coral TPU (~$25-60) | Free + hardware |
| **LLM Vision** | Camera snapshots → Claude/GPT | Context-aware ("Is Zeus at desk?") | Per-image cost, latency | $0.003-0.01/image |
| **Double Take** | Facial recognition | Identifies specific people | Requires training, complex | Free |
| **Hybrid** | Frigate + LLM Vision | Best of both worlds | Most complex setup | Hardware + API |

**Recommended: Frigate + LLM Vision (Hybrid)**

```
Camera → Frigate (local, fast)
              ↓
        "Person detected"
              ↓
        LLM Vision (only when needed)
              ↓
        "Zeus is at his desk working"
              ↓
        Legacy-AI decides action
```

**Why hybrid:**

- Frigate handles 99% of detections locally (privacy, speed)
- LLM Vision triggers only on important events (saves cost)
- Legacy can ask "What's happening in the living room?"

**Hardware needed:**

- Google Coral USB TPU (~$25-60) for Frigate AI
- Yi cameras with yi-hack firmware

---

### 3. Will SendGrid give Legacy an inbox like Gmail?

**Short answer:** No visual inbox, but YES for receiving emails.

**What SendGrid provides:**
- **Outbound:** Send emails via API (Legacy can email you)
- **Inbound Parse:** Receive emails via webhook (emails → POST to your endpoint)
- **NO:** Visual inbox, web interface to read emails, Gmail-like experience

**How it works:**
```
Someone emails legacy@yourdomain.com
              ↓
    SendGrid receives it
              ↓
    POST webhook to your server
              ↓
    Your server sends to Letta
              ↓
    Legacy processes and can respond
```

**Alternatives for visual inbox:**

| Option | Visual Inbox | Cost | Complexity |
|--------|--------------|------|------------|
| **Gmail API** | Yes (Gmail) | Free (15GB) | Medium - OAuth setup |
| **Zoho Mail** | Yes | Free (5GB) | Easy - standard setup |
| **SendGrid + Custom UI** | Build your own | SendGrid cost | High |
| **Microsoft 365** | Yes (Outlook) | $6/month | Easy |

**Recommendation:**
- For MVP: Use SendGrid (no inbox needed - Legacy stores in memory)
- For visual inbox: Create a Gmail account (`legacy.ai.assistant@gmail.com`) and use Gmail API

---

## Cost Analysis

### Monthly Service Costs

| Service | Free Tier | Paid Tier | Your Expected Usage |
|---------|-----------|-----------|---------------------|
| **Letta Cloud** | Generous free tier | $20/month (Pro) | Free tier likely sufficient |
| **Vapi** | $10 free credits | Pay-as-you-go | ~$5-15/month typical |
| **Home Assistant Cloud** | None | $6.50/month | Remote access + TTS |
| **SendGrid** | 100 emails/day | $19.95/month | Free tier sufficient |
| **Twilio** | Trial | Pay-as-you-go | ~$1-5/month (if used) |
| **Yi Cameras** | One-time purchase | - | Already owned |

### Per-Use Costs

| Action | Cost | Breakdown |
|--------|------|-----------|
| **1 min voice call (Vapi)** | ~$0.10 | STT + LLM + TTS |
| **10 min voice call** | ~$0.70-1.00 | Depends on talk time |
| **SMS (Twilio)** | $0.0079/msg | Outbound text |
| **Email (SendGrid)** | Free | Under 100/day |
| **HA Cloud TTS** | Included | In $6.50/month |
| **LLM Vision (per image)** | $0.003-0.01 | See below |

### LLM Vision Costs

| Provider | Cost/Image | Notes |
|----------|------------|-------|
| **Claude (Anthropic)** | ~$0.003-0.01 | Best quality analysis |
| **GPT-4o (OpenAI)** | ~$0.005-0.01 | Good balance |
| **Gemini (Google)** | ~$0.001-0.005 | Cheapest cloud option |
| **Ollama (local)** | Free | Requires GPU, runs on Beelink |

**Monthly estimate (100 analyses/day):** ~$15/month with Claude, $0 with Ollama

### Total Monthly Estimate

| Usage Level | Monthly Cost |
|-------------|--------------|
| **Light** (few calls, HA basics) | $6.50 (just HA Cloud) |
| **Moderate** (daily calls, smart home) | $15-25 |
| **Heavy** (frequent calls, full automation) | $30-50 |

---

## Implementation Phases

### Phase 1: Home Assistant Setup (Priority: HIGH)

**Goal:** Get HA running on Beelink with voice control via ATOM Echo

**Steps:**
1. [ ] Install Home Assistant OS on Beelink
   - Download HA OS image for x86-64
   - Flash to USB, boot from USB, install to SSD
   - Initial setup at `http://homeassistant.local:8123`

2. [ ] Subscribe to Home Assistant Cloud ($6.50/month)
   - Enables remote access
   - Enables cloud TTS (no local setup needed)
   - Auto-backup to cloud
   - Mobile app push notifications

3. [ ] Set up ATOM Echo
   - Flash with ESPHome firmware
   - Connect to HA via ESPHome integration
   - Configure as voice satellite

4. [ ] Add Yi Cameras
   - Get RTSP stream URLs from Yi cameras
   - Add as generic camera integration
   - Test live view in HA dashboard

5. [ ] Add SwitchBot devices
   - Install SwitchBot integration (HACS)
   - Pair plugs and humidity sensor
   - Create basic automations

**Test:** Say "Hey Nabu, turn on the plug" through ATOM Echo

---

### Phase 2: LLM Vision Integration (Priority: MEDIUM)

**Goal:** Legacy can "see" who's in the room via cameras

**Steps:**
1. [ ] Install LLM Vision from HACS
   - Repository: https://github.com/valentinfrlch/ha-llmvision
   - Requires HACS installed first

2. [ ] Configure with Claude API
   - Add Anthropic API key
   - Point to Yi camera entities

3. [ ] Create presence detection automation
   ```yaml
   automation:
     - alias: "Check if Zeus is home"
       trigger:
         - platform: time_pattern
           minutes: "/15"  # Every 15 minutes
       action:
         - service: llmvision.analyze
           data:
             camera_entity: camera.yi_office
             prompt: "Is there a person in this image? If yes, describe them briefly."
   ```

4. [ ] Create Letta tool: `check_presence(room)`
   - Calls HA REST API
   - Uses LLM Vision to analyze camera
   - Returns presence status

---

### Phase 3: Legacy ↔ Home Assistant Bridge (Priority: HIGH)

**Goal:** Legacy can control and query Home Assistant

**Steps:**
1. [ ] Create Long-Lived Access Token in HA
   - Profile → Security → Create Token
   - Store securely in Letta environment

2. [ ] Create Letta tools:
   ```python
   # Tools to add to Legacy-AI

   def get_home_status() -> dict:
       """Get current state of all devices and sensors."""
       # GET /api/states

   def control_device(entity_id: str, action: str) -> str:
       """Turn devices on/off, set values."""
       # POST /api/services/{domain}/{service}

   def speak_in_room(message: str, room: str = "office") -> str:
       """Make Legacy speak through Home Assistant."""
       # POST /api/services/tts/cloud_say

   def check_presence(room: str) -> dict:
       """Check if Zeus is in a specific room using camera."""
       # Uses LLM Vision integration
   ```

3. [ ] Test bidirectional communication
   - Legacy → HA: "Turn off the office plug"
   - HA → Legacy: Voice command triggers Letta query

---

### Phase 4: Email Setup - SendGrid (Priority: MEDIUM)

**Goal:** Legacy can send and receive emails

**Steps:**
1. [ ] Create SendGrid account (free tier)
2. [ ] Authenticate your domain (DNS records)
3. [ ] Set up Inbound Parse webhook
   - Create endpoint in Legacy-Comms system
   - Point SendGrid to webhook URL
4. [ ] Create Letta tools:
   ```python
   def send_email(to: str, subject: str, body: str) -> str:
       """Send an email as Legacy."""

   def get_recent_emails() -> list:
       """Get emails received in last 24 hours."""
   ```

**Alternative:** Gmail API for visual inbox
- Create: `legacy.ai.assistant@gmail.com`
- Set up OAuth credentials
- Legacy can read/send through familiar Gmail interface

---

### Phase 5: ICA - Intelligent Communication Activation (Priority: FUTURE)

**Goal:** Legacy proactively communicates through the right channel

**Prerequisites:**
- Phase 1-4 complete
- All communication channels working

**Implementation:**
1. [ ] Create communication decision engine
2. [ ] Implement mode controls (Active, Standby, Silent)
3. [ ] Connect all channels:
   - Phone/SMS via Vapi/Twilio
   - Email via SendGrid
   - Voice via Home Assistant + ATOM Echo
   - Discord/Telegram for async

**Decision Tree:**
```
Legacy needs to communicate
         ↓
    Check mode → STANDBY? → Wait for wake word
         ↓ NO
    Zeus at home? → YES → Speak via ATOM Echo
         ↓ NO
    Urgency? → CRITICAL → Phone call
         ↓ NORMAL
    Send via preferred channel (SMS/Email/Discord)
```

---

## Quick Start Commands

### Check Legacy Status
```bash
# Test Letta connection
curl -s https://api.letta.com/v1/agents/agent-da585f6f-189f-472a-974c-29b9578f9b55 \
  -H "Authorization: Bearer $LETTA_API_KEY" | jq .name

# Test Vapi
curl -s https://api.vapi.ai/assistant/cb290108-e8b7-4e47-b4e3-e003ca822308 \
  -H "Authorization: Bearer 41231b6e-e17f-4959-b2e3-8ab6778cdf13" | jq .name
```

### Talk to Legacy
```bash
# Via Letta API
curl -X POST https://api.letta.com/v1/agents/agent-da585f6f-189f-472a-974c-29b9578f9b55/messages \
  -H "Authorization: Bearer $LETTA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"messages": [{"role": "user", "content": "Hello Legacy"}]}'

# Via phone
# Call: +1 (203) 760-5895
# Say call sign: "Zulu Echo Uniform Sierra November Yankee Charlie"
```

---

## Immediate Next Steps

1. **Today:** Install Home Assistant OS on Beelink
2. **Today:** Subscribe to Home Assistant Cloud
3. **This week:** Set up ATOM Echo and Yi cameras
4. **This week:** Upgrade Twilio ($20) for SMS capability
5. **Next:** Create Letta tools for HA integration
6. **Future:** SendGrid email setup
7. **Future:** Full ICA implementation

---

## Related Documentation

- [LEGACY-COMMS.md](./LEGACY-COMMS.md) - Communications system details
- [LEGACY-AI-SYSTEM.md](./LEGACY-AI-SYSTEM.md) - Core AI documentation
- [LETTA-RESEARCH-LEGACY-AI.md](./LETTA-RESEARCH-LEGACY-AI.md) - Letta platform research

---

## Session Resume Checklist

When starting a new Claude session, share this:

> "Read LEGACY-PLAN.md at root for full context. I'm working on [current phase].
> Key IDs are in that file. Last session we [what you did]."

This gives Claude instant context without re-explaining everything.

---

## Session Log

### 2026-01-06: Home Assistant Setup Complete

**Accomplishments:**
- HA REST API access confirmed via Nabu Casa URL
- SSH access configured with Claude's public key
- Initial backup created: `Legacy-Initial-Setup-2026-01-06` (slug: `e98fe189`)
- Comprehensive documentation created: `LEGACY-HA-REFERENCE.md`

**System Info:**
- HA OS: 16.3
- HA Core: 2025.12.5
- Machine: generic-x86-64 (Beelink)
- Disk: 234 GB total, 7.9 GB used
- Add-ons: Music Assistant, Let's Encrypt, Assist Microphone, Git pull, Terminal & SSH

**SSH Public Key (Claude):**
```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOtEsVg+0RB9Vf5kGrsnwW3iw2MFG0Jbg1PNC+dFrqT6 claude@anthropic
```

**Next Steps:**
- Set up ATOM Echo voice satellite
- Add Yi cameras to HA
- Add SwitchBot devices

### 2026-01-07: Legacy-Comms Merged into Legacy-AI

**Major Restructure Complete:**

- Merged Legacy-Comms into Legacy-AI (single unified agent)
- Vapi now points to Legacy-AI (agent-da585f6f-189f-472a-974c-29b9578f9b55)
- Added `communications` tool for SMS and outbound calls
- Enabled Twilio caller ID restriction (only Zeus's phone can call)
- Removed call sign authentication (caller ID is sufficient)
- Standby mode retained in Legacy-AI system prompt
- Deleted Twilio Studio Flow (was for SMS routing)
- Deleted old SMS serverless handler

**Security Model:**

- **Caller ID Filter**: Twilio function checks incoming number
- **Allowed**: +1 (203) 303-4844 (Zeus's phone)
- **Blocked**: All other numbers silently rejected
- **No verbal authentication needed**

**New Capabilities:**

- Legacy can send SMS: `communications(action="sms", message="...")`
- Legacy can call Zeus: `communications(action="call", message="...")`
- Unified voice + brain + home control in one agent

**Architecture:**

```
Phone Call → Twilio → Caller ID Filter → Vapi → Legacy-AI (Letta)
              │              │                        ↓
              │         (Zeus only)              Tools:
              │              │                   • home_assistant()
              │              ▼                   • communications()
              │         Other callers
              │         → Rejected
              │
SMS Outbound ← Twilio ← communications() tool
```

**Cloud Configurations (not in git):**

- Twilio phone voice URL: `legacy-voice-filter-2815-dev.twil.io/voice-filter`
- Vapi assistant model: `agent-da585f6f-189f-472a-974c-29b9578f9b55`
- Vapi first message: "Legacy online. How can I help, Mr. Delacruz?"
- Legacy-AI system prompt: Updated with communications instructions

**Files Changed:**

- `tools/communications_tool.py` - New SMS/call tool (gitignored - contains secrets)
- `LEGACY-PLAN.md` - Updated architecture docs
- `LEGACY-COMMS.md` - Marked as deprecated

**Test Command:**

- Call +1 (203) 583-3480 from Zeus's phone
- Legacy answers: "Legacy online. How can I help, Mr. Delacruz?"
- Try: "Turn on the lights"

### 2026-01-07 (Evening): SMS Removed, Cleanup Complete

**Decision:** Remove SMS capability due to A2P 10DLC compliance requirements.

US carriers now require A2P 10DLC registration for SMS from local numbers. Rather than deal with registration overhead, decided to use Telegram for notifications instead (simpler, free, richer content).

**Cleanup Actions:**

| Resource | Action |
|----------|--------|
| Letta: `communications` tool | Deleted |
| Letta: Legacy system prompt | Updated (SMS references removed) |
| Local: `tools/communications_tool.py` | Deleted |
| Twilio: Studio flows | Already cleaned (previous session) |
| Twilio: SMS webhook | Already empty |

**Kept (still needed):**

| Resource | Purpose |
|----------|---------|
| Twilio: voice-filter function | Caller ID restriction |
| Twilio: phone number | Voice calls to Legacy |

**Current Legacy Capabilities:**

- ✅ Voice calls (Twilio → Caller ID Filter → Vapi → Letta)
- ✅ Caller ID restriction (only Zeus's phone: +1 203-303-4844)
- ✅ Home Assistant control (lights, TTS, etc.)
- ✅ Agent delegation (business/personal specialists)
- ❌ SMS (removed - A2P compliance)
- 🔜 Telegram notifications (planned)

**Legacy-AI Tools:**

```
tool-2449d4bc-26fb-4677-89b4-43d7dde83a01 - home_assistant
```

**Legacy-AI System Prompt Summary:**

- Security: Caller ID restriction at phone level
- Address: "Mr. Delacruz"
- Modes: Full Comms (default), Standby Mode
- Personality: Direct, military brevity, forward-thinking
- Capabilities: Home Assistant control, agent delegation

**Next Steps:**

- Add Telegram bot for notifications
- Test voice call with updated system
