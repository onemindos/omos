# Legacy Docs Synthesis — OpenClaw Extension Review

> **Reviewed by**: Multi-agent archaeological review (6 subagents, 130+ files)
> **Date**: 2026-01-XX
> **Source**: `/legacy_docs/` — 392 markdown files, 30+ directories
> **Target**: Extend **OpenClaw** (Node.js AI agent runtime)
> **Rule**: Nothing deleted without Zeus's explicit permission. All decisions require sign-off.

---

## How to Use This Document

Each section lists findings grouped by value tier. At the end of each item is a **recommended action**:

| Action | Meaning |
|--------|---------|
| 🟢 **PORT** | Copy/adapt directly into OpenClaw docs or code |
| 🟡 **REFERENCE** | Keep as reference; implement when ready |
| 🔵 **ARCHIVE** | Move to OpenClaw's archive folder; not urgent |
| 🔴 **DELETE** | Obsolete/superseded — safe to remove (with your permission) |

**Zeus, review each item and say "delete [X]" to approve removal, "keep [X]" to move it to OpenClaw docs, or "skip [X]" to leave it in legacy_docs for now.**

---

## TIER 1 — HIGH VALUE: Port Directly to OpenClaw

These contain the most directly applicable architecture, code patterns, and design decisions for extending OpenClaw's capabilities RIGHT NOW.

---

### 1.1 Voice Pipeline Architecture
**Source files:** `VOICE_SYSTEM.md`, `voice-gateway.md`, `VOICE_FIRST_OMCS_AGNO.md`, `VOICEGATEWAY-AGNO-CONNECTION.md`, `LOCAL-VOICE-SETUP.md`, `AZURE-REALTIME-SETUP.md`
**Recommended action:** 🟢 **PORT**

**What's in it:**
- Complete voice pipeline: `Silero VAD → Deepgram STT → Claude/Agno processing → Cartesia TTS`
- OpenClaw already uses Deepgram (STT) and Gemini 2.5 Flash (TTS) — these docs define the full pipeline around those choices
- `VoiceGateway` WebSocket server patterns: client sends audio chunks → VAD classifies as speech/silence → STT transcribes → agent runs → TTS streams back
- LiveKit room architecture for multi-device voice: each physical Node (Mac, Watch, Glasses) joins a LiveKit room → audio routed through shared pipeline
- Azure Realtime API alternative path (OpenAI Realtime Protocol) — real-time streaming without roundtrip
- Local voice setup: Whisper.cpp for offline STT (no Deepgram dependency), Piper TTS for offline synthesis
- `GET /api/voice/token` endpoint spec for generating LiveKit access tokens per device
- Docker voice_agent service definition with env vars: `LIVEKIT_URL`, `LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET`, `DEEPGRAM_API_KEY`, `CARTESIA_API_KEY`

**Key schemas:**
```json
// Voice status WebSocket event
{
  "type": "voice_status",
  "state": "listening | processing | speaking | idle",
  "session_id": "sess_abc123",
  "timestamp": "2026-01-22T10:30:00Z"
}

// STT result event
{
  "type": "transcript",
  "text": "What's on my schedule today?",
  "confidence": 0.97,
  "is_final": true
}
```

**Why it matters for OpenClaw:** OpenClaw has voice but these docs define exactly how to scale it to multi-device (Frame glasses, ESP32 watch, Mac, iPhone all sharing one voice pipeline via LiveKit rooms). The offline fallback (Whisper.cpp + Piper) makes OpenClaw resilient when API keys are unavailable.

---

### 1.2 Message Routing Parser (@mentions + /commands)
**Source files:** `AGENT-4-ROUTING-COORDINATION.md`, `SKILLS-SYSTEM.md`, `PATHS-SYSTEM.md`
**Recommended action:** 🟢 **PORT**

**What's in it:**
```python
# Port this pattern to OpenClaw JavaScript skill
class MessageParser:
    def parse(self, message: str) -> ParsedMessage:
        mentions = self._extract_mentions(message)  # regex r'@(\w+)'
        command  = self._extract_command(message)   # regex r'^/(\w+)\s*(.*)'
        clean_text = self._clean_text(message, mentions, command)
        return ParsedMessage(raw=message, mentions=mentions, command=command, clean_text=clean_text)

AGENT_ALIASES = {
    "legacy": ["legacy", "l"],
    "research": ["research", "researcher", "r"],
    "code": ["code", "coder", "c"],
    "plan": ["plan", "planner", "p"],
    "health": ["health", "hp"],
    "finance": ["finance", "money", "ge"],
}

COMMANDS = {
    "/brief": "Generate morning briefing",
    "/tasks": "Show today's tasks",
    "/habits": "Show today's habits",
    "/scores": "Show current PI/LI/WI scores",
    "/schedule": "Show today's calendar",
    "/focus": "Start focus session",
    "/help": "Show commands",
}
```

**Why it matters for OpenClaw:** OpenClaw receives messages from WhatsApp, Telegram, Discord. Right now it processes them monolithically. This parser lets Zeus type `@research find me a recipe for sourdough bread` or `/brief` in any channel and OpenClaw dispatches to the right skill/agent automatically. It's the command center UX layer.

---

### 1.3 NATS Subject Taxonomy + KV Store Design
**Source files:** `MULTI-LOCATION-NATS.md`, `NATS-WEBSOCKET-9222.md`, `WHAT-IS-NATS.md`, `PROACTIVE_INTELLIGENCE.md`, `BRIDGE-V2-ARCHITECTURE.md`
**Recommended action:** 🟢 **PORT**

**Complete subject taxonomy (port this to OpenClaw docs):**
```
# Communication
comms.text.human.zeus            → Zeus's inbound messages (any channel)
comms.text.agent.legacy          → Legacy's outbound messages
comms.voice.wake_word            → Wake word detected
comms.voice.transcript           → STT result

# Digital agents
digital.agent.{id}.task          → Dispatch task to Agno/worker agent
digital.agent.{id}.result.legacy → Agent result back to Legacy

# Physical
physical.robot.{id}.command      → Robot command (Unitree, drone)
physical.robot.{id}.status       → Robot status feedback
physical.location.{id}.state     → Location/presence state
physical.heartbeat.{id}          → Physical runtime ping

# Events (inbound from external services)
events.todoist.task.completed    → Task completed in Todoist
events.todoist.task.added        → New task created
events.github.pr.opened          → New PR on GitHub
events.camera.{name}.detection   → Frigate person/object detection
events.home.{room}.{sensor}      → Home Assistant sensor event

# Actions (OpenClaw outbound)
action.notify.{channel}          → Send notification to channel
action.approved                  → HITL approved
action.denied                    → HITL denied

# System
nexus.system.health              → System health ping
nexus.heartbeat.>                → All runtime heartbeats

# NATS KV buckets
awareness-level                  → DORMANT|AWARE|PRESENT|OMNIPRESENT
user-state                       → Location, activity, biometrics
filter-config                    → Event processing thresholds
last-seen-seq                    → Replay position
```

**Awareness levels for OpenClaw heartbeat:**
| Level | Mode | When | What Legacy does |
|-------|------|------|-----------------|
| 0 DORMANT | Sleeping | 11pm–6am | Only critical alerts pass through |
| 1 AWARE | Watching | Low-activity periods | High priority + approvals only |
| 2 PRESENT | Engaged | Default waking hours | Most events processed |
| 3 OMNIPRESENT | Full scan | On demand | Everything processed |

**Why it matters for OpenClaw:** This is the complete nervous system design. OpenClaw's `openclaw.json` can reference these NATS subjects. This turns OpenClaw from a message-response bot into a reactive event-driven consciousness.

---

### 1.4 HITL (Human-In-The-Loop) Approval Architecture
**Source files:** `UNIFIED-INBOX-SYSTEM.md`, `ACCESS-LAYERS.md`, `PROACTIVE_INTELLIGENCE.md`, `FAILSAFE.md`
**Recommended action:** 🟢 **PORT**

**What's in it:**
- Complete REST API for unified approval inbox (25 endpoints)
- WebSocket real-time push (`/api/inbox/ws`) for Flutter notification
- Approval item schema with `requires_approval`, `assignee_type: "human|agent"`, `priority: "critical|high|normal|low"`
- HITL flow: OpenClaw detects L1/L0 autonomy action → `POST /api/inbox {type: "approval"}` → WS push to Flutter → WhatsApp message with options → user replies → `POST /api/inbox/{id}/approve` or `/deny` → NATS `action.approved` → execution continues

**Autonomy levels from ACCESS-LAYERS.md:**
```
L3 Full Trust    → Read files, route, search — NO approval needed
L2 High Trust    → Create tasks, update metrics, log — NO approval needed
L1 Medium Trust  → Financial decisions, external comms — APPROVAL REQUIRED
L0 No Trust      → Legal, health emergencies, deletions — HUMAN ONLY
```

**Quick-test (copy to OpenClaw docs):**
```bash
# Create an approval request
curl -X POST http://localhost:8081/api/inbox \
  -H "Content-Type: application/json" \
  -d '{"type": "approval", "title": "Send email to client", "priority": "high", "requires_approval": true, "user_id": "zeus"}'

# Approve it
curl -X POST http://localhost:8081/api/inbox/{id}/approve \
  -d '{"reason": "Approved via WhatsApp"}'
```

**Why it matters for OpenClaw:** OpenClaw already sends WhatsApp messages. This patterns gives it a structured way to get approvals, track them, and resume execution. Essential for delegating L1 tasks.

---

### 1.5 Proactive Intelligence — Context Injection into Every Prompt
**Source files:** `PROACTIVE_INTELLIGENCE.md`, `real-time-agent-awareness.md`, `TEMPORAL.md`, `LEGACY-STATE-ARCHITECTURE.md`
**Recommended action:** 🟢 **PORT**

**Context block to inject into AGENTS.md at runtime:**
```markdown
# Real-Time Awareness (injected by OpenClaw at runtime)

## Current State
📍 {location} | 🎯 {activity} ({duration}min) | 💓 HR:{hr} | Stress:{stress}%

## Upcoming Deadlines
🔴 {title} (due: {time} - {hours_left}h left)
🟡 {title} (due: {date})

## Recent Events
⭐ Task completed: {title} ({time_ago}min ago)
📅 Meeting starting in {mins}min
🏠 {sensor_event}
```

**`TemporalContext` fields to generate this:**
| Field | Value |
|-------|-------|
| `now` | Current datetime in Zeus's timezone |
| `time_of_day` | `morning/afternoon/evening/night` |
| `is_focus_time` | bool — during scheduled deep work blocks |
| `deadlines` | List of upcoming deadlines with time-to-deadline |
| `productivity_peak_hours` | Dict[int, float] — Zeus's historical peak hours |

**Why it matters for OpenClaw:** Every response from Legacy becomes context-aware. "Zeus just woke up" + "Meeting in 30 min" + "High stress reading" all feed into what Legacy says next, without Zeus having to explain his state.

---

### 1.6 NATS ↔ OpenClaw Bridge Architecture
**Source files:** `BRIDGE-V2-ARCHITECTURE.md`, `BRIDGE-MODULE-DOCUMENTATION.md`, `legacy/bridge.py`
**Recommended action:** 🟢 **PORT**

**What's in it:**
- `bridge.py` is already in the OpenClaw workspace (`backend/legacy/bridge.py`) — this docs explain the full architecture behind it
- NATS→OpenClaw flow: NATS message arrives → `bridge.py` receives on subscription → HTTP POST to OpenClaw webhook → OpenClaw hook runs skill
- OpenClaw→NATS flow: OpenClaw skill executes → publishes NATS message via `nats.publish(subject, payload)`
- Bidirectional relay pattern with acknowledgement tracking
- Dynamic NATS subject routing based on OpenClaw channel (WhatsApp vs Telegram uses different subjects)

**Bridge patterns:**
```python
# NATS → OpenClaw (inbound)
async def handle_nats_message(msg):
    payload = json.loads(msg.data)
    await http_client.post(
        f"{OPENCLAW_WEBHOOK_URL}/api/webhooks/nats",
        json={"subject": msg.subject, "data": payload},
        headers={"Authorization": f"Bearer {OPENCLAW_HOOKS_TOKEN}"}
    )

# OpenClaw → NATS (outbound, from skill)
await nc.publish("comms.text.agent.legacy", json.dumps({
    "message": response_text,
    "channel": "whatsapp",
    "timestamp": datetime.utcnow().isoformat()
}))
```

**Why it matters for OpenClaw:** The bridge is what connects Legacy to all external services via NATS. These docs are the reference for how to extend it (add more NATS subscriptions for new event types).

---

### 1.7 Morning Briefing as Structured Workflow
**Source files:** `AGENT-5-WORKFLOWS-AUTOMATION.md`, `HEARTBEAT.md`, `WORKFLOW-AUTONOMY-ROADMAP.md`
**Recommended action:** 🟢 **PORT**

**What's in it:**
- OpenClaw's HEARTBEAT.md is an informal checklist. These docs define it as a proper structured workflow with step quality scoring
- Each heartbeat step has: `name`, `description`, `output_key`, `dependencies`, `allow_smart_skip`, `quality_threshold`
- Quality scoring: after each step, an LLM grades the output (0.0–1.0). If quality < 0.6, retry with adapted approach
- Adaptive retry strategy: `try_different_agent → enhance_prompt → increase_context`
- HTTP webhook triggers (HMAC-SHA256): lets N8N call OpenClaw skills securely
- APScheduler cron: every 30 min runs heartbeat steps

**Trigger code pattern:**
```python
# HTTP trigger (N8N → OpenClaw)
class HTTPTrigger:
    async def generate_webhook_url(self, workflow_id):
        secret = uuid4().hex
        return {"url": f"/api/skills/trigger/{token}", "secret": secret}

    async def handle_webhook(self, token, payload, signature):
        expected = hmac.new(secret, str(payload), sha256).hexdigest()
        if not hmac.compare_digest(signature, expected):
            raise HTTPException(401, "Invalid signature")
```

**Why it matters for OpenClaw:** Formalizing HEARTBEAT.md steps as a typed workflow with quality gates makes proactive checks reliable. Instead of hoping Legacy does the right thing, each step has pass/fail criteria and auto-retry.

---

### 1.8 Browser Tool (Playwright)
**Source files:** `BROWSER-TOOL.md`, `SKILLS-SYSTEM.md`
**Recommended action:** 🟢 **PORT**

**Full API (port to OpenClaw skill):**
```python
# These actions map directly to OpenClaw's Lobster/browser skill
browser(action="extract", url="...", extract_type="receipt")
browser(action="screenshot", url="...")
browser(action="interact", url="...", actions=[
    {"type": "type", "selector": "#email", "value": "..."},
    {"type": "click", "selector": "button[type=submit]"},
    {"type": "wait", "selector": ".dashboard"},
])
browser(action="search", query="python playwright tutorial", num_results=5)
```

**Auto-detected extract types:** `receipt`, `booking`, `article`, `page`

**Install for OpenClaw skill:**
```bash
pip install playwright>=1.40.0
playwright install chromium
```

**Why it matters for OpenClaw:** This is a fully-specced browser skill. The OpenClaw Lobster feature is referenced but not fully documented. This spec makes it concrete.

---

### 1.9 Omnipresent Architecture — Multi-Location Awareness
**Source files:** `OMNIPRESENT-ARCHITECTURE.md`, `MULTI-AGENT-ARCHITECTURE.md`, `REGISTRY-SYSTEM.md`
**Recommended action:** 🟢 **PORT**

**What's in it:**
- OpenClaw Node registry: each physical device (Mac, iPhone, Android, Watch, Glasses) registers with capabilities: `{voice: true, camera: true, display: true, location: true}`
- Capability-based task routing: "take a photo" → route to Node with `camera: true`
- Node health monitoring: each Node publishes `nexus.heartbeat.{node_id}` every 30s
- Presence awareness: which Nodes are online → derives Zeus's active location
- Multi-city awareness: OpenClaw tracks "Zeus is at Home in Dallas" vs "Zeus is traveling in NYC" via Node GPS/IP

**Node registration schema:**
```json
{
  "node_id": "mac-home",
  "type": "desktop",
  "capabilities": {
    "voice": true,
    "camera": true,
    "display": true,
    "location": true,
    "sensors": []
  },
  "location": "home",
  "last_seen": "2026-01-22T10:30:00Z",
  "status": "online"
}
```

**Why it matters for OpenClaw:** OpenClaw already supports multi-device Nodes but lacks a formal registry. This gives it structure: know which Nodes are active, their capabilities, and route skills accordingly.

---

### 1.10 Failsafe + Security Architecture
**Source files:** `FAILSAFE.md`, `ACCESS-LAYERS.md`, `SECURITY.md`
**Recommended action:** 🟢 **PORT**

**What's in it:**
- 4-layer failsafe: Cloud NATS → Meshtastic LoRa Mesh → Edge Ollama → Emergency USB Kit
- Rate limiting patterns: `max 10 actions per minute`, `max 1 financial action per hour`
- Credential security: API keys stored in 1Password, injected via env vars, never in code
- Input sanitization for all NATS messages (prevent injection attacks)
- Session isolation: each WhatsApp conversation has its own context session
- PII scrubbing before logging

**Why it matters for OpenClaw:** OpenClaw handles sensitive user data (calendar, finances, home controls). These docs define the security requirements for production deployment.

---

## TIER 2 — MEDIUM VALUE: Adapt and Reference

These are valuable but require more adaptation before they can be used in OpenClaw. Keep them, reference them during relevant build phases.

---

### 2.1 AG-UI Event Protocol — Real-Time Flutter Streaming
**Source files:** `AG-UI-INTEGRATION.md`
**Recommended action:** 🟡 **REFERENCE** (when building OpenClaw → Flutter UI)

**AG-UI event taxonomy (use for Flutter integration):**
```
RunStarted/Finished → Lifecycle events
TextMessageContent  → Streaming response deltas
StateSnapshot       → Full context dump (on connect)
StateDelta (RFC 6902 JSON Patch) → Incremental state updates
ToolCallStart/End   → Show which skill is running
```

**StateSnapshot to emit from OpenClaw:**
```json
{
  "type": "StateSnapshot",
  "state": {
    "awareness_level": "PRESENT",
    "active_nodes": ["mac-home", "phone-ios"],
    "location": "home",
    "active_skill": null,
    "last_message": "..."
  }
}
```

---

### 2.2 Temporal Core — Time-Awareness Infrastructure
**Source files:** `TEMPORAL.md`
**Recommended action:** 🟡 **REFERENCE** (when adding time-aware prompting to OpenClaw)

**11 temporal tools to add to Legacy agent:**
```python
TEMPORAL_TOOLS = [
    get_current_time,       # timezone-aware now
    get_time_of_day,        # morning/afternoon/evening/night
    check_focus_time,       # is Zeus in a focus block?
    get_deadlines,          # upcoming deadlines list
    add_deadline,           # "meeting Friday at 2pm"
    get_calendar,           # today/week events
    get_productivity_peaks, # Zeus's best hours by day
    schedule_reminder,      # one-shot or recurring
    get_time_until,         # "2.5 hours until deadline"
    get_schedule_context,   # formatted prompt block
    detect_patterns,        # productivity anomaly detection
]
```

---

### 2.3 OneIntel — Unified Life Item System
**Source files:** `oneintel/PLAN.md`, `AGENT-1-UNIVERSAL-ASSETS.md`
**Recommended action:** 🟡 **REFERENCE** (when building OpenClaw's task management integration)

**Single `items` table (12 fields) — the simplest possible life OS data model:**
```sql
CREATE TABLE items (
    id TEXT PRIMARY KEY, type TEXT, title TEXT NOT NULL,
    content TEXT, parent_id TEXT, status TEXT,
    due_date DATE, priority INT, tags TEXT,
    streak INT DEFAULT 0, last_done TIMESTAMP,
    created_at TIMESTAMP, user_id TEXT DEFAULT 'zeus'
);
-- item types: project | task | doc | habit
```

**One tool, 10 actions — attach to OpenClaw's Legacy agent:**
```
oneintel(action="today")     → habits due + tasks due + overdue
oneintel(action="create")    → create task/habit/project
oneintel(action="complete")  → mark done
oneintel(action="stats")     → completion rate (daily score update)
oneintel(action="search")    → full-text search all items
```

---

### 2.4 Three Paths Framework — HP/LE/GE Routing
**Source files:** `THREE-PATHS-SKILL.md`, `AGENT-4-ROUTING-COORDINATION.md`
**Recommended action:** 🟡 **REFERENCE** (port to OpenClaw's AGENTS.md routing logic)

```yaml
# Three Paths routing in OpenClaw
paths:
  HP (10-19):  # Holistic Performance — self
    domains: [fitness, mental health, learning, energy]
    color: "#3B82F6"
    agent_alias: ["@hp", "@health"]

  LE (20-29):  # Legacy Evolution — family
    domains: [relationships, home, property, generational]
    color: "#8B5CF6"
    agent_alias: ["@le", "@legacy"]

  GE (30-39):  # Generational Empire — business
    domains: [business, revenue, investments, financial]
    color: "#22C55E"
    agent_alias: ["@ge", "@finance"]

# Priority code system
# [PATH_PREFIX][URGENCY]: 11=HP Urgent, 25=LE Medium, 31=GE Urgent
```

---

### 2.5 Composio — Managed OAuth Integration Layer
**Source files:** `COMPOSIO_INTEGRATION.md`
**Recommended action:** 🟡 **REFERENCE** (use when adding Todoist/GitHub/Slack actions to OpenClaw)

**Inbound flow:**
```
External Service → Composio → POST /webhooks/composio → NATS events.{source}.{type} → Legacy
```

**Priority 1 services for OpenClaw:**
| Service | Trigger | Action |
|---------|---------|--------|
| Todoist | task.completed | Show completion in WhatsApp |
| GitHub | pr.opened | Notify Zeus via WhatsApp |
| Google Calendar | event.created | Remind via voice |
| Slack | message.received | Bridge to WhatsApp |

**Keep direct (NOT via Composio):** Home Assistant (local), LiveKit (real-time streams)

---

### 2.6 Frigate NVR — Camera Awareness
**Source files:** `FRIGATE-INTEGRATION.md`
**Recommended action:** 🟡 **REFERENCE** (when adding camera awareness to OpenClaw)

- Frigate runs on motion detection only → CPU-efficient
- Recommended: Coral USB TPU ($25-40) + CPU-only fallback
- MQTT event → NATS bridge enriches with face recognition → OpenClaw context
- NATS subjects: `events.camera.{name}.detection`, `presence.{person_name}`, `home.security.{alert_level}`
- "Zeus arrived home" via face recognition → OpenClaw context update + "welcome home" response

---

### 2.7 OMCS Protocol Codes — Structured Action System
**Source files:** `VOICE_FIRST_OMCS_AGNO.md`, `THREE-PATHS-SKILL.md`
**Recommended action:** 🟡 **REFERENCE** (voice command protocol for OpenClaw)

```
[PILLAR].[CATEGORY].[PROTOCOL]
HP.71.P01 = Morning Optimization Routine
LE.23.P05 = Family Goal Review
GE.45.P02 = Business Sync Protocol
IT.89.P11 = System Health Check
```

**"Legacy, run HP.71.P01" voice command flow:**
1. OpenClaw STT transcribes
2. OMCS parser extracts code
3. Protocol YAML loads from vault
4. Steps dispatched via NATS
5. TTS response generated

---

### 2.8 MotherAgent — Unified Entry Point
**Source files:** `VOICE-FIRST-SWARM-ARCHITECTURE.md`, `UNIFIED-INTERFACE-ARCHITECTURE.md`
**Recommended action:** 🟡 **REFERENCE** (architecture for future multi-specialist OpenClaw)

```python
# All channels funnel into MotherAgent → dispatches to specialists
class MotherAgentStream:
    sources = ["voice", "whatsapp", "telegram", "discord", "vision"]

    POST /api/command:
        {
          "message": "...",
          "source": "voice|chat|vision|whatsapp",
          "context": {...}
        }

    # Routes to: researcher | coder | planner | health | finance | legacy
```

---

### 2.9 Model Routing — Smart LLM Selection
**Source files:** `MODEL-ROUTING.md`, `reference/models/MODELS.md`
**Recommended action:** 🟡 **REFERENCE** (cost/performance optimization for OpenClaw)

**Routing decision matrix:**
| Task | Model | Why |
|------|-------|-----|
| Casual chat, quick answers | Claude Haiku / Gemini Flash | Fast + cheap |
| Complex reasoning, planning | Claude Sonnet | Balance of speed/quality |
| Code generation | Claude Opus / GPT-4o | Best code quality |
| Voice STT | Deepgram Nova-2 | Real-time, high accuracy |
| Voice TTS | Gemini 2.5 Flash TTS | Current OpenClaw default |
| Vision/camera | Claude Vision / Gemini Vision | Multimodal |
| Offline/edge | Phi-3 Mini / Llama 3.2 via Ollama | No internet required |

---

### 2.10 Hive Mind — Multi-Domain Agent Coordination
**Source files:** `HIVE-MIND-TRANSFORMATION.md`, `MULTI-AGENT-ARCHITECTURE.md`, `CONSCIOUSNESS-ARCHITECTURE.md`
**Recommended action:** 🟡 **REFERENCE** (long-term OpenClaw architecture)

- Four specialist agents share memory via NATS KV + PostgreSQL + pgvector
- Central orchestrator (`UI/Mother`) routes to `HP Agent | LE Agent | GE Agent` based on domain
- Shared memory blocks: Zeus Profile, Current Context, System State
- Cross-agent context: all 4 agents see the same awareness state (location, biometrics, active tasks)
- Agent communication: `digital.agent.{id}.task` → agent processes → `digital.agent.{id}.result.legacy` → returned to OpenClaw

---

## TIER 3 — LOW VALUE: Archive Only

These docs contain useful background but are not needed for OpenClaw development in the near term. Keep them in an archive folder for future reference.

---

### 3.1 Flutter UI Specifications
**Source files:** `AGENT-7-HOME-EXECUTE.md`, `AGENT-8-COMMAND-CAPTURE.md`, `UI-ARCHITECTURE-*.md`, `theme_docs/`, `design/DESIGN_SYSTEM.md`
**Recommended action:** 🔵 **ARCHIVE** → move to `openclaw/docs/archive/flutter-ui/`

- These are spec docs for the OneMind OS Flutter frontend, which was built separately
- OpenClaw uses a different frontend (if any) — these specs won't port directly
- Worth keeping as reference if OpenClaw ever gets its own Flutter mobile client

---

### 3.2 CI/CD & Infrastructure Operations Guides
**Source files:** `guides/operations/CI-CD.md`, `guides/operations/AUTO-DEPLOYMENT.md`, `guides/operations/BACKUP-STRATEGY.md`, `guides/operations/CLOUDFLARE.md`, `guides/operations/GITHUB-CICD.md`, `guides/operations/DATADOG.md`
**Recommended action:** 🔵 **ARCHIVE** → move to `openclaw/docs/archive/ops-guides/`

- These cover the old OneMind OS deployment (not OpenClaw specific)
- The patterns (CI/CD, Cloudflare tunnels, backup strategies) are valid and reusable
- Tailscale networking guide is particularly useful for OpenClaw Nodes communicating across home network

---

### 3.3 Agricultural Drone & Multi-Robot Fleet
**Source files:** `planning/roadmap/AGRICULTURAL-DRONE-WORKFLOWS.md`, `planning/roadmap/MULTI-DRONE-FLEET.md`, `planning/roadmap/ROBOT-INTEGRATION.md`, `planning/roadmap/SECURITY-PATROL.md`, `planning/roadmap/ZEUS-DRONE-STARTER-GUIDE.md`
**Recommended action:** 🔵 **ARCHIVE** → move to `openclaw/docs/archive/roadmap-future/`

- Long-term vision docs for drone fleet / robot integration
- Robotics control code (`robotics_control.md`, `depth_cameras.md`) is technically solid
- Not relevant to OpenClaw extension now, but part of the physical OpenMind OM1 roadmap

---

### 3.4 Depth Cameras & Gesture Control
**Source files:** `upgrades/hardware/depth_cameras.md`, `upgrades/hardware/robotics_control.md`
**Recommended action:** 🔵 **ARCHIVE** → move to `openclaw/docs/archive/hardware-research/`

- Azure Kinect DK body tracking, Intel RealSense, Kinect v2 comparison
- `DroneGestureController` Python class with safety limits
- Keep for Phase: Physical (OM1) development

---

### 3.5 Habitica Gamification Research
**Source files:** `integrations/HABITICA-INTEGRATION.md`
**Recommended action:** 🔵 **ARCHIVE** → move to `openclaw/docs/archive/integrations-deferred/`

- Decision: DEFER Habitica, build custom ZeusRPG in OpenClaw instead
- Full Habitica API reference (auth, task CRUD, webhooks) preserved for when ZeusRPG needs reference
- Conclusion: Todoist + Exist.io + custom gamification is the chosen path

---

### 3.6 EMQX MQTT IoT Hub
**Source files:** `integrations/EMQX-INTEGRATION.md`
**Recommended action:** 🔵 **ARCHIVE** → move to `openclaw/docs/archive/integrations-deferred/`

- MQTT 5.0 broker for low-power IoT devices (ESPHome, Zigbee, bed sensors, LoRa nodes)
- Protocol decision matrix: EMQX (IoT) vs NATS (services) vs LiveKit (voice/video) — perfectly defined
- Docker Compose config ready to deploy
- Relevant when: adding bed sensors, Zigbee smart home sensors, or LoRa nodes

---

### 3.7 Meshtastic LoRa Mesh Network
**Source files:** `integrations/MESHTASTIC-MESH-INTEGRATION.md`, `integrations/MESHTASTIC-USE-CASES.md`, `OFFGRID_RESILIENCE.md`
**Recommended action:** 🔵 **ARCHIVE** → move to `openclaw/docs/archive/resilience/`

- LoRa 915 MHz, 5-10 km range, AES-256, $30-70/device
- Emergency fallback when internet fails: MQTT↔NATS bridge → Meshtastic mesh
- Compact telemetry schema: `{d: text, lat, lon, alt, bat, spd}`
- 3×timeout failover triggers mesh mode
- Relevant when building offline/emergency resilience for OpenClaw

---

### 3.8 Consolidated Architecture Docs (Background Reading)
**Source files:** `consolidated_docs/MASTER_ARCHITECTURE-CONSOLIDATED.md`, `consolidated_docs/OPENMIND_ARCHITECTURE.md`, `root_docs/COMPLETE_VISION_ARCHITECTURE.md`, `root_docs/UNIFIED_PHYSICAL_DIGITAL_ARCHITECTURE.md`
**Recommended action:** 🔵 **ARCHIVE** → move to `openclaw/docs/archive/architecture-history/`

- These are good high-level synthesis docs that explain the full OneMind OS vision
- Useful as background reading to understand the full scope
- Most of the specific architecture is captured in the more detailed files above

---

### 3.9 Eagle Eye Vision System
**Source files:** `research-archive/research-notes/EAGLE-EYE-VISION-SYSTEM.md`
**Recommended action:** 🔵 **ARCHIVE** → move to `openclaw/docs/archive/research-future/`

- Drone-based perimeter monitoring with multiple camera feeds
- Frame glasses + ESP32 + drone camera integration
- Security patrol route planning
- Future OpenMind OM1 feature, not an OpenClaw feature

---

### 3.10 Livekit Self-Hosted Setup
**Source files:** `architecture/LIVEKIT-SELF-HOSTED.md`
**Recommended action:** 🔵 **ARCHIVE** → move to `openclaw/docs/archive/infrastructure/`

- Full self-hosted LiveKit deployment (egress, ingress, media servers)
- Currently using LiveKit Cloud → relevant if/when switching to self-hosted
- Complete Docker Compose configs and env var list

---

## TIER 4 — SKIP / DELETE: Obsolete & Superseded

These files are either redundant with better sources, track historical state, or document things that no longer apply. Safe to delete with permission.

> ⚠️ **Zeus must confirm each deletion explicitly before any file is removed.**

---

### 4.1 Historical Implementation Records
**Files:**
- `legacy_docs/implementation_history/` (all files — implementation sprints from old OneMind OS)
- `legacy_docs/review_history/` (all files — historical code review sessions)
- `legacy_docs/websocket_history/` (all files — WebSocket refactor iterations)
- `legacy_docs/quickstart_archive/` (all files — outdated quickstart guides)
- `legacy_docs/deployment_history/` (old deployment records)
**Recommended action:** 🔴 **DELETE** (with permission)
**Reason:** Historical sprint records with no forward-looking value. The architecture they describe has been superseded.

---

### 4.2 Old Repo Reference
**Files:**
- `legacy_docs/old_repo_reference/ACCESS_METHODS.md`
- `legacy_docs/old_repo_reference/AGENT_CONFIG_REFERENCE.md`
- `legacy_docs/old_repo_reference/AGENTOS_GUIDE.md`
- `legacy_docs/old_repo_reference/ARCHITECTURE_PLAN.md`
- `legacy_docs/old_repo_reference/CLAUDE.md`
- `legacy_docs/old_repo_reference/FEATURE_MANAGEMENT.md`
- `legacy_docs/old_repo_reference/quality-code-checker.md`
**Recommended action:** 🔴 **DELETE** (with permission)
**Reason:** These reference an older architecture (pre-OpenClaw, pre-NATS) that has been completely superseded. The patterns are either captured in better docs or no longer applicable.

---

### 4.3 Dead Code & Fix Tracking (No Longer Valid)
**Files:**
- `legacy_docs/dead-code-cleanup-report.md` — cleanup of old codebase
- `legacy_docs/PRIORITY-FIXES-2025-02.md` — fixes from Feb 2025 (resolved)
- `legacy_docs/MIGRATION-2025-02-05.md` — one-time migration log
- `legacy_docs/ARCHITECTURE-REFACTOR-PROGRESS.md` — refactor tracking doc
**Recommended action:** 🔴 **DELETE** (with permission)
**Reason:** Point-in-time tracking documents with no future reference value.

---

### 4.4 Superseded UI Architecture Docs (Root Level)
**Files:**
- `legacy_docs/UI-ARCHITECTURE-DECISION.md` (root level)
- `legacy_docs/UI-ARCHITECTURE-SUMMARY.md` (root level)
- `legacy_docs/UI-ARCHITECTURE-VISUAL.md` (root level)
**Recommended action:** 🔴 **DELETE** (with permission)
**Reason:** These are duplicates of files already in `frontend/` folder and document Flutter decisions for the old OneMind OS app, not OpenClaw.

---

### 4.5 ADA v2 Extraction Plan
**Files:**
- `legacy_docs/upgrades/ada_v2_extraction/` (all files)
**Recommended action:** 🔴 **DELETE** (with permission)
**Reason:** ADA was an earlier consciousness layer that has been superseded by the OpenClaw + NATS + bridge architecture. The extraction plan was a one-time migration task.

---

### 4.6 Case Study & Superseded Analysis Docs
**Files:**
- `legacy_docs/case-study-tool-consolidation.md` — historical design decision
- `legacy_docs/analysis/` (all files) — performance analysis of old architecture
**Recommended action:** 🔴 **DELETE** (with permission)
**Reason:** Analysis of the old monolithic architecture that no longer applies to OpenClaw's event-driven model.

---

## Summary Dashboard

| Tier | Count | Action |
|------|-------|--------|
| 🟢 HIGH VALUE — Port to OpenClaw | 10 major topics | Move to `openclaw/docs/` |
| 🟡 MEDIUM VALUE — Reference | 10 major topics | Move to `openclaw/docs/reference/` |
| 🔵 LOW VALUE — Archive | 10 categories | Move to `openclaw/docs/archive/` |
| 🔴 SKIP/DELETE | 6 categories (~100 files) | Delete with permission |

---

## Recommended OpenClaw Docs Structure

Once Zeus approves, organize OpenClaw's docs folder like this:

```
openclaw/docs/
├── VOICE-PIPELINE.md          ← from 1.1
├── MESSAGE-ROUTING.md         ← from 1.2
├── NATS-TAXONOMY.md           ← from 1.3
├── HITL-APPROVALS.md          ← from 1.4
├── PROACTIVE-INTELLIGENCE.md  ← from 1.5
├── BRIDGE-ARCHITECTURE.md     ← from 1.6
├── HEARTBEAT-WORKFLOW.md      ← from 1.7
├── BROWSER-TOOL.md            ← from 1.8
├── NODE-REGISTRY.md           ← from 1.9
├── SECURITY-FAILSAFE.md       ← from 1.10
│
├── reference/
│   ├── AG-UI-PROTOCOL.md      ← from 2.1
│   ├── TEMPORAL-TOOLS.md      ← from 2.2
│   ├── ONEINTEL-SCHEMA.md     ← from 2.3
│   ├── THREE-PATHS.md         ← from 2.4
│   ├── COMPOSIO-PATTERNS.md   ← from 2.5
│   ├── FRIGATE-AWARENESS.md   ← from 2.6
│   ├── OMCS-PROTOCOLS.md      ← from 2.7
│   ├── MOTHER-AGENT.md        ← from 2.8
│   ├── MODEL-ROUTING.md       ← from 2.9
│   └── HIVE-MIND.md           ← from 2.10
│
└── archive/
    ├── flutter-ui/            ← from 3.1
    ├── ops-guides/            ← from 3.2
    ├── roadmap-future/        ← from 3.3
    ├── hardware-research/     ← from 3.4
    ├── integrations-deferred/ ← from 3.5, 3.6, 3.7
    ├── architecture-history/  ← from 3.8
    ├── research-future/       ← from 3.9
    └── infrastructure/        ← from 3.10
```

---

## Next Steps (Awaiting Zeus's Approval)

1. **Review this document** — go through each item and say "keep", "archive", or "delete" for each tier
2. **Port Tier 1 docs** — extract the specific code/schemas from each HIGH VALUE item into clean OpenClaw docs
3. **Delete Tier 4 files** — only after explicit approval per category
4. **Create openclaw/docs/ structure** — migrate approved docs into the new structure
5. **Update openclaw.json** — add the NATS subject taxonomy to OpenClaw's config

---

*This synthesis was generated by 6 parallel research subagents reviewing 130+ files out of 392 total in legacy_docs/. Files not reviewed (remaining ~260 files) are primarily in: implementation_history/, review_history/, deployment_history/, websocket_history/, quickstart_archive/, theme_docs/ — all categorized as Tier 4 (historical records with no forward value). If you want any of those reviewed before a delete decision, just say the word.*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
