---
type: communication-protocol
updated: 2026-05-26
---

---

## Universal Communications Protocol (NexusMessage)

Every message in OneMind OS — text, voice, camera, sensor, agent — flows through a single canonical format: **NexusMessage**.

### Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    HUMAN INPUTS                         │
│  Text Chat │ Voice │ Camera │ Files │ Sensor Dashboard  │
└─────────┬───────┬────────┬──────┬────────┬──────────────┘
          │       │        │      │        │
          ▼       ▼        ▼      ▼        ▼
    ┌─────────────────────────────────────────────┐
    │           NexusMessage Envelope             │
    │  content + modality + attachments + routing  │
    └───────────────────┬─────────────────────────┘
                        │
              ┌─────────▼─────────┐
              │   NATS Bus        │
              │   comms.{mod}.    │
              │   {src_type}.     │
              │   {src_id}        │
              └─────────┬─────────┘
                        │
          ┌─────────────┼─────────────┐
          │             │             │
    ┌─────▼────┐  ┌─────▼────┐  ┌────▼─────┐
    │  Legacy  │  │ Socket.IO│  │  Agents  │
    │  Sees    │  │ Bridge   │  │  Aware   │
    │  ALL     │  │ Forwards │  │  Listen  │
    └──────────┘  └──────────┘  └──────────┘
```

### NexusMessage Format

```python
NexusMessage(
    id="uuid",
    content="What's in this image?",
    modality="vision",                     # text|voice|vision|audio|sensor|command|multimodal
    source_id="human:zeus",
    source_type="human",                   # human|agent|team|robot|sensor|system
    attachments=[
        Attachment(type="image", url="http://...", mime_type="image/jpeg"),
        Attachment(type="audio", data_base64="...", mime_type="audio/wav"),
    ],
    data={},
    routing=Routing(target_ids=["researcher"], mode="direct"),  # auto|direct|broadcast
    session_id="conv-123",
    correlation_id="msg-456",
    priority="normal",                     # critical|high|normal|low|background
    timestamp="2025-01-15T...",
)
```

### NATS Subject Convention

```
comms.{modality}.{source_type}.{source_id}

Examples:
  comms.text.human.zeus           # Human text chat
  comms.voice.human.zeus          # Voice command (transcribed)
  comms.vision.sensor.garage-cam  # Camera snapshot
  comms.sensor.robot.arm-01       # Robot telemetry
  comms.text.agent.researcher     # Agent response
```

### JetStream Streams

| Stream | Subjects | Purpose |
|--------|----------|---------|
| DIGITAL | digital.> | Agent activity |
| PHYSICAL | physical.> | Robot/sensor events |
| HUMAN | human.> | Human-in-the-loop |
| SYSTEM | system.> | Health/heartbeats |
| NEXUS | nexus.> | Orchestration events |
| **COMMS** | **comms.>** | **Universal multimodal messages** |

### NATS Switches

```python
"comms.text": True,
"comms.voice": True,
"comms.vision": True,
"comms.audio": True,
"comms.sensor": True,
"comms.command": True,
"comms.chat": True,
```

### Modality-Aware Entities

```python
EntityDescriptor(
    id="researcher",
    input_modalities=["text", "image"],
    output_modalities=["text"],
    capabilities=["chat", "reasoning", "tool_use", "vision"],
)
```

Auto-detected: Voice agents → `audio`, Vision models → `image`, Robots → `sensor/command/video`

---
# OneMind OS — Communications Protocol

## Channels

| Platform | Purpose | Access |
|----------|---------|--------|
| Telegram | Primary Zeus ↔ Agent comms | Bot + Userbot |
| Discord | Community / dev collab | Server |
| SMS/iMessage | Personal, family, urgent | imsg CLI |
| Email | Formal, business, receipts | Fastmail (himalaya CLI) |
| Signal | Secure 1:1 | — |

## Routing Rules

1. **Urgent / safety** → SMS or Telegram DM (immediate)
2. **Task results** → Telegram (home channel)
3. **Daily briefs** → Telegram (morning/evening)
4. **External / public** → Never direct from Legacy — route to Unity profile
5. **Business comms** → Email via Fastmail
6. **Community** → Discord or Skool (via GoHighLevel)

## Agent-to-Agent

| Path | Transport | Subject Pattern |
|------|-----------|-----------------|
| Legacy → OpenClaw | NATS | agents.prompt.oc.* |
| OpenClaw → Legacy | NATS / Hermes API | agents.prompt.hermes.* |
| Legacy → Claude Code | NATS | agents.prompt.cc.* |
| Cron notifications | Hermes deliver | origin / telegram / all |

## Notification Tiers

| Tier | Criteria | Delivery |
|------|----------|----------|
| 🔴 Critical | Security, health, money >$500 | Telegram + SMS immediately |
| 🟡 Important | Task complete, errors, decisions needed | Telegram |
| 🟢 Info | Status updates, FYI | Telegram (silent) or queue |
| ⚪ Log | Routine confirmations | No delivery (local only) |

## Standing Comms

- **07:00 ET** — Morning brief (priorities, calendar, pending)
- **22:00 ET** — Evening wrap (summary, tomorrow's plan)
- **On-demand** — Zeus messages in any channel → immediate response

## Escalation

If Zeus is offline (23:00–06:00 ET):
- 🔴 Critical → deliver immediately anyway
- 🟡 Important → queue, deliver at 07:00
- 🟢/⚪ → queue silently

---

## Voice Infrastructure (Legacy Reference)

> **Status:** Historical — Letta/Vapi era (pre-Hermes). Preserved for reference.
> Original system: Legacy-Comms (Vapi/Twilio) → Legacy-AI (Letta). Now unified in Hermes.

### Phone Numbers

| Purpose | Number | Provider |
|---------|--------|----------|
| Legacy-Comms (Primary) | +1 (203) 583-3480 | Twilio via Vapi |
| Legacy-Comms (Vapi Native) | +1 (203) 760-5895 | Vapi |
| Zeus Personal | +1 (203) 303-4844 | — |

### Security Protocol

**Call Sign:** `Zulu Echo Uniform Sierra November Yankee Charlie` (ZEUSNYC)
**Address:** Mr. Delacruz

**Auth Flow:**
1. Call Legacy → *"Legacy comms. Authenticate."*
2. Say call sign
3. Legacy confirms → *"Call sign confirmed. Channel secure."*

### Comms Modes

| Mode | Trigger | Behavior |
|------|---------|----------|
| Full Comms | "Full comms" | Responds to all speech |
| Standby | "Activate standby mode" | Only responds to "Legacy" prefix |
| Silent | "Go silent" | No proactive communication |
| Emergency Only | "Emergency only mode" | Critical alerts only |

### Legacy Agent IDs (Letta Era)

| Agent | ID | Platform |
|-------|----|----------|
| Legacy-AI (Brain) | `agent-da585f6f-189f-472a-974c-29b9578f9b55` | Letta |
| Legacy-Comms (Voice) | `agent-028b397a-efe8-45f8-804c-8fed8605f8bd` | Letta/Vapi |
| Vapi Assistant | `cb290108-e8b7-4e47-b4e3-e003ca822308` | Vapi |
| Vapi Phone ID | `d2a1ddd0-46df-4baf-9639-7ef529a768ae` | Vapi |
| Studio Flow | `FWb0ca50ad9d7370016079068b6822af57` | Twilio |

### Intelligent Communication Activation (ICA) — Vision

> Goal: Legacy proactively communicates via the right channel based on presence + urgency.

**Decision Tree:**
```
Legacy needs to communicate
         ↓
    Standby mode? → YES → Wait for wake word
         ↓ NO
    Zeus at home? → YES → Speak via Home Assistant speaker
         ↓ NO
    Urgency CRITICAL → Phone call
         ↓ NORMAL
    Send via preferred channel (SMS / Email / Telegram)
```

**Planned channels:**
- Phone: Twilio/Vapi ✅ Active
- SMS: Twilio ⚠️ Needs upgrade
- Email: SendGrid (legacy@domain) 📋 Planned
- Local Voice: Home Assistant + Wyoming + Piper 📋 Planned
- Presence Detection: LLM Vision + ESP32-CAM 📋 Planned

### Vapi Settings Reference

| Setting | Value |
|---------|-------|
| silenceTimeoutSeconds | 300 |
| maxDurationSeconds | 3600 |
| backgroundDenoisingEnabled | true |
| recordingEnabled | true |
| firstMessage | "Legacy comms. Authenticate." |
| endCallMessage | "Comms offline. Legacy out." |
| endCallPhrases | "End call", "Goodbye Legacy", "Comms out" |

---

## Command Grammar

Structured shorthand for issuing commands across agents, robots, and services. Each token narrows scope — left to right — until Legacy has everything needed to act without asking.

### Syntax

```
{contact} - {agent} - {intent} - {command}
```

| Token | Role | Examples |
|-------|------|---------|
| `contact` | Urgency / priority of the contact | `urgent`, `now`, `when ready`, `scheduled` |
| `agent` | Which agent or system receives it | `legacy`, `sentinel`, `zenith`, `drone`, `robot`, `home` |
| `intent` | Domain / what category of action | `drone`, `camera`, `lights`, `nav`, `report`, `launch` |
| `command` | Specific action to execute | `up`, `down`, `sweep`, `lock`, `status`, `follow` |

### Example

```
contact (urgency) - legacy (agent) - drone (intent) - up (command)
→ "Legacy, launch the drone closest to me now."
```

### More Examples

| Shorthand | Full Execution |
|-----------|---------------|
| `urgent - legacy - drone - up` | Launch nearest drone immediately |
| `urgent - legacy - drone - sweep barn` | Launch drone, sweep Barn perimeter |
| `now - sentinel - camera - lock` | Sentinel locks all cameras to current zone |
| `now - legacy - findzeus - sweep` | Run full Find Zeus sweep immediately |
| `urgent - legacy - killswitch - L3` | Shut down NATS bus |
| `when ready - zenith - deploy - api` | Zenith deploys the API when current task is done |
| `now - home - lights - off` | All lights off via Home Assistant |
| `urgent - legacy - robot - recall` | Recall all active robots to home position |
| `now - legacy - status - all` | Full system status report |

### Rules

- Tokens can be dropped if implied by context — `drone up` = `urgent - legacy - drone - up` when Zeus is clearly talking to Legacy
- Urgency defaults to `now` if omitted
- Agent defaults to `legacy` if omitted
- Legacy maps `intent + command` to the correct NATS subject and fires it — no clarification needed unless the command is ambiguous
- NATS subject derived: `cmd.{agent}.{intent}.{command}` (e.g. `cmd.drone.launch.nearest`)

---

## Zeus → AI Intent Keywords

Short phrases with defined execution behavior. Legacy recognizes these instantly and adjusts response mode accordingly.

| Keyword / Phrase | Intent | Legacy Behavior |
|------------------|--------|----------------|
| `now execute` | Fire immediately | No confirmation, no options — act and report after |
| `in parallel` | Concurrent execution | All independent calls fire simultaneously |
| `copy comms` / `comms copy` | Radio check | Reply: "Copy that." — nothing else |
| `go` / `execute` | Start the task | Same as now execute — immediate action |
| `brief` | Summary mode | Short, structured, no filler |
| `browse` | Session summary | Run session_search(), report recent sessions + status |
| `pause` | Hold | Stop current task, await next instruction |
| `status` | Current state | Report what's running / pending / blocked |
| `d` | Delete (vault triage) | Delete the current file immediately, confirm + next |
| `k` | Keep (vault triage) | Keep + route/move the file, confirm + next |
| `leave` | Skip (vault triage) | Skip this file, move to next — no action |
| `next` | Advance | Move to next item in current workflow |
| `L1` | Approval required | Present plan, wait for Zeus to confirm before acting |

### Rules
- **Zeus never repeats himself** — if a keyword was used once, it's standing until changed
- **No filler responses** — keyword triggers → action, not explanation
- **Tactical tone always** — direct, structured, tables preferred


---

## Related Protocol Docs (Vault Cross-References)

Intent keywords and comms channels are the Zeus-facing layer. Full protocol depth lives in these vault docs:

| Doc | Path | What It Covers |
|-----|------|----------------|
| NexusMessage Spec | `SO/00 Framework/blueprint/system-design/general/system-primitives.md` | Canonical message envelope, NATS subject schema |
| Legacy Agent Profile | `SO/01 Intelligence/agents/_roster/_COMMAND/legacy/00 - Legacy Agent.md` | Legacy's comms channels, escalation, morning/evening routines |
| NATS Fabric Plugin | `SO/01 Intelligence/agents/_runtimes/hermes/NATS-FABRIC-PLUGIN.md` | NATS pub/sub integration with Hermes |
| Proactive Architecture | `SO/07 Architecture/05-execution-proactive/PROACTIVE_ARCHITECTURE.md` | How Legacy initiates comms proactively |
| Voice-First Architecture | `SO/04 Operations/voice/VOICE-FIRST-SWARM-ARCHITECTURE.md` | Voice channel design, LiveKit, wake words |
| Hive Protocol | `SO/01 Intelligence/memory/general/HIVE-PROTOCOL.md` | Agent-to-agent comms protocol |

> **Note:** `comms.md` is the Zeus-facing interface layer — intent keywords, channel routing, notification tiers. The docs above are the implementation layer. Don't duplicate them here.
