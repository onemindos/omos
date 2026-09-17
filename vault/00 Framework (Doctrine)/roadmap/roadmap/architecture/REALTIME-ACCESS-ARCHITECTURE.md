# OneMind OS: Real-Time Access Architecture

> **Created:** January 2026
> **Purpose:** Real-time voice, HITL, 24/7 access, and native app architecture
> **Key Concern:** Legacy must be accessible everywhere, always, with instant response

---

## Section 1: The Real-Time Layer

### 1.1 Cloud + Local Hybrid Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                        REAL-TIME HYBRID ARCHITECTURE                                 │
│                                                                                      │
│   ┌─────────────────────────────────────┐     ┌─────────────────────────────────┐  │
│   │         CLOUD (AWS Lightsail)       │     │        LOCAL (Home Server)      │  │
│   │         ──────────────────          │     │        ─────────────────        │  │
│   │                                     │     │                                 │  │
│   │  ┌─────────────────────────────┐   │     │  ┌─────────────────────────┐   │  │
│   │  │     LIVEKIT CLOUD           │   │     │  │    HOME ASSISTANT       │   │  │
│   │  │     (Voice Rooms)           │   │     │  │    (Smart Home)         │   │  │
│   │  │     < 100ms latency         │   │     │  │    Local network        │   │  │
│   │  └─────────────────────────────┘   │     │  └─────────────────────────┘   │  │
│   │                                     │     │                                 │  │
│   │  ┌─────────────────────────────┐   │     │  ┌─────────────────────────┐   │  │
│   │  │     NATS JETSTREAM          │   │     │  │    NATS LEAF NODE       │   │  │
│   │  │     (Central Hub)           │◄──┼─────┼──►   (Local Events)        │   │  │
│   │  │     200-400μs               │   │     │  │    Sub-ms local         │   │  │
│   │  └─────────────────────────────┘   │     │  └─────────────────────────┘   │  │
│   │                                     │     │                                 │  │
│   │  ┌─────────────────────────────┐   │     │  ┌─────────────────────────┐   │  │
│   │  │     LETTA                   │   │     │  │    KESTRA               │   │  │
│   │  │     (The Brain)             │   │     │  │    (Orchestrator)       │   │  │
│   │  │     Memory + Reasoning      │   │     │  │    Workflows            │   │  │
│   │  └─────────────────────────────┘   │     │  └─────────────────────────┘   │  │
│   │                                     │     │                                 │  │
│   │  ┌─────────────────────────────┐   │     │  ┌─────────────────────────┐   │  │
│   │  │     GATEWAY                 │   │     │  │    TIMESCALEDB          │   │  │
│   │  │     (Webhooks)              │   │     │  │    (Event Archive)      │   │  │
│   │  └─────────────────────────────┘   │     │  └─────────────────────────┘   │  │
│   │                                     │     │                                 │  │
│   └──────────────────┬──────────────────┘     └──────────────┬──────────────────┘  │
│                      │                                       │                      │
│                      │         TAILSCALE MESH                │                      │
│                      │         (Encrypted, NAT-traversal)    │                      │
│                      └───────────────────────────────────────┘                      │
│                                                                                      │
│   LATENCY TARGETS:                                                                  │
│   • Voice (LiveKit): < 200ms end-to-end                                            │
│   • NATS events: < 1ms local, < 50ms cloud                                         │
│   • LLM response: 500ms-2s (streaming helps)                                       │
│   • Action execution: < 100ms                                                       │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Latency Optimization Strategy

| Component | Target | How |
|-----------|--------|-----|
| Voice capture | < 50ms | LiveKit's WebRTC, local VAD |
| Voice → NATS | < 10ms | Direct publish from LiveKit agent |
| NATS → Letta | < 5ms | In-process call or gRPC |
| Letta thinking | 500ms-2s | Streaming responses |
| Response → TTS | < 100ms | ElevenLabs streaming |
| Total voice loop | < 2.5s | Acceptable for conversation |

---

## Section 2: LiveKit Multimodal Integration

### 2.1 Where LiveKit Fits

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         LIVEKIT IN THE ONEMIND STACK                                 │
│                                                                                      │
│   USER DEVICES                     LIVEKIT CLOUD                  ONEMIND           │
│   ────────────                     ────────────                   ───────           │
│                                                                                      │
│   ┌──────────┐                   ┌──────────────────┐                               │
│   │  Phone   │───────────────────►                  │                               │
│   │  (App)   │  WebRTC           │   LIVEKIT        │                               │
│   └──────────┘                   │   CLOUD          │                               │
│                                   │                  │                               │
│   ┌──────────┐                   │  ┌────────────┐  │      ┌─────────────────┐     │
│   │  Laptop  │───────────────────►  │   ROOM     │  │      │                 │     │
│   │  (Web)   │  WebRTC           │  │            │  │      │  Legacy AGENT   │     │
│   └──────────┘                   │  │  Zeus +    │◄─┼──────►  (LiveKit SDK)  │     │
│                                   │  │  Legacy    │  │      │                 │     │
│   ┌──────────┐                   │  │            │  │      │  Processes:     │     │
│   │ Desktop  │───────────────────►  └────────────┘  │      │  • STT          │     │
│   │  (App)   │  WebRTC           │                  │      │  • Letta call   │     │
│   └──────────┘                   │                  │      │  • TTS          │     │
│                                   │                  │      │  • NATS pub     │     │
│   ┌──────────┐                   │                  │      │                 │     │
│   │ Smart    │───────────────────►                  │      └────────┬────────┘     │
│   │ Glasses  │  WebRTC           │                  │               │              │
│   └──────────┘                   └──────────────────┘               │              │
│                                                                      │              │
│                                                                      ▼              │
│                                                          ┌─────────────────┐       │
│                                                          │  NATS JETSTREAM │       │
│                                                          │                 │       │
│                                                          │  events.voice.> │       │
│                                                          │  actions.voice.>│       │
│                                                          └─────────────────┘       │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 LiveKit Multimodal Agent Flow

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                      LIVEKIT MULTIMODAL AGENT FLOW                                   │
│                                                                                      │
│   ZEUS SPEAKS                                                                       │
│        │                                                                            │
│        ▼                                                                            │
│   ┌─────────────────┐                                                               │
│   │  LIVEKIT VAD    │  Voice Activity Detection                                     │
│   │  (< 50ms)       │  Detects speech start/end                                     │
│   └────────┬────────┘                                                               │
│            │                                                                         │
│            ▼                                                                         │
│   ┌─────────────────┐                                                               │
│   │  STT (Deepgram) │  Speech-to-Text                                               │
│   │  (streaming)    │  Real-time transcription                                      │
│   └────────┬────────┘                                                               │
│            │                                                                         │
│            ▼                                                                         │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                             │   │
│   │   Legacy AGENT (Python)                                                     │   │
│   │                                                                             │   │
│   │   1. Receive transcript                                                     │   │
│   │   2. Publish to NATS: events.voice.transcript                              │   │
│   │   3. Call Letta API with transcript                                        │   │
│   │   4. Stream Letta response                                                 │   │
│   │   5. Publish to NATS: events.legacy.response                               │   │
│   │   6. Send to TTS                                                           │   │
│   │                                                                             │   │
│   │   INTERRUPT DETECTION:                                                      │   │
│   │   • If Zeus speaks during response → immediate stop                        │   │
│   │   • Publish: events.voice.interrupt                                        │   │
│   │   • Cancel TTS, process new input                                          │   │
│   │                                                                             │   │
│   └─────────────────────────────────────────────────────────────────────────────┘   │
│            │                                                                         │
│            ▼                                                                         │
│   ┌─────────────────┐                                                               │
│   │  TTS (Eleven)   │  Text-to-Speech                                               │
│   │  (streaming)    │  Streaming audio chunks                                       │
│   └────────┬────────┘                                                               │
│            │                                                                         │
│            ▼                                                                         │
│   ZEUS HEARS Legacy                                                                 │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 2.3 Voice Commands for Execution Control

```python
# agent/voice/commands.py

INTERRUPT_PHRASES = [
    "stop", "cancel", "wait", "hold on", "never mind",
    "legacy stop", "legacy cancel", "abort", "halt"
]

CONFIRMATION_PHRASES = [
    "yes", "do it", "go ahead", "confirm", "approved",
    "legacy yes", "proceed", "execute"
]

REJECTION_PHRASES = [
    "no", "don't", "reject", "deny", "cancel that",
    "legacy no", "negative", "abort that"
]

class VoiceCommandProcessor:
    """Process voice commands for real-time control."""

    def __init__(self, nats_client):
        self.nats = nats_client
        self.pending_actions = {}  # execution_id -> action details

    async def process_transcript(self, transcript: str) -> str:
        """Check for command phrases in transcript."""

        lower = transcript.lower().strip()

        # Check for interrupts first (highest priority)
        for phrase in INTERRUPT_PHRASES:
            if phrase in lower:
                await self.handle_interrupt()
                return "interrupt"

        # Check for confirmations (HITL responses)
        for phrase in CONFIRMATION_PHRASES:
            if phrase in lower:
                await self.handle_confirmation()
                return "confirm"

        # Check for rejections (HITL responses)
        for phrase in REJECTION_PHRASES:
            if phrase in lower:
                await self.handle_rejection()
                return "reject"

        return "normal"

    async def handle_interrupt(self):
        """Cancel all running executions."""
        # Publish interrupt to NATS
        await self.nats.publish(
            "commands.interrupt",
            json.dumps({
                "type": "interrupt",
                "scope": "all",
                "timestamp": datetime.utcnow().isoformat()
            })
        )

        # Cancel any pending Kestra workflows
        for exec_id in list(self.pending_actions.keys()):
            await self.cancel_kestra_execution(exec_id)

    async def handle_confirmation(self):
        """Approve pending HITL action."""
        if self.pending_actions:
            # Get most recent pending action
            exec_id = list(self.pending_actions.keys())[-1]
            action = self.pending_actions[exec_id]

            await self.nats.publish(
                "commands.hitl.approve",
                json.dumps({
                    "execution_id": exec_id,
                    "action": action,
                    "approved_by": "voice",
                    "timestamp": datetime.utcnow().isoformat()
                })
            )

            del self.pending_actions[exec_id]
```

---

## Section 3: Voice-Based Cancellation System

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         VOICE CANCELLATION ARCHITECTURE                              │
│                                                                                      │
│   "Legacy, CANCEL THAT"                                                             │
│          │                                                                          │
│          ▼                                                                          │
│   ┌─────────────────┐                                                               │
│   │  STT Detection  │  < 100ms                                                      │
│   │  "cancel"       │                                                               │
│   └────────┬────────┘                                                               │
│            │                                                                         │
│            ▼                                                                         │
│   ┌─────────────────┐                                                               │
│   │ Command Processor│                                                              │
│   │ Recognizes      │                                                               │
│   │ interrupt phrase│                                                               │
│   └────────┬────────┘                                                               │
│            │                                                                         │
│            ▼                                                                         │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                             │   │
│   │                        NATS: commands.interrupt                             │   │
│   │                                                                             │   │
│   └────────┬────────────────────┬────────────────────┬──────────────────────────┘   │
│            │                    │                    │                              │
│            ▼                    ▼                    ▼                              │
│   ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐                    │
│   │                 │  │                 │  │                 │                    │
│   │  KESTRA         │  │  LETTA          │  │  TTS/VOICE      │                    │
│   │                 │  │                 │  │                 │                    │
│   │  Cancel running │  │  Stop current   │  │  Stop speaking  │                    │
│   │  workflow via   │  │  tool execution │  │  immediately    │                    │
│   │  API call       │  │                 │  │                 │                    │
│   │                 │  │                 │  │                 │                    │
│   └─────────────────┘  └─────────────────┘  └─────────────────┘                    │
│                                                                                      │
│   CANCELLATION TARGETS:                                                             │
│   • TTS output (immediate)                                                          │
│   • Letta tool execution (if running)                                              │
│   • Kestra workflow (via API)                                                       │
│   • Pending HITL approval (clears queue)                                           │
│   • Any action.* in progress                                                        │
│                                                                                      │
│   LATENCY: Voice → Cancel = < 500ms                                                │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Section 4: HITL (Human-in-the-Loop) System

### 4.1 HITL Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                          HITL (HUMAN-IN-THE-LOOP) SYSTEM                             │
│                                                                                      │
│   ACTION CLASSIFICATION                                                             │
│   ─────────────────────                                                             │
│                                                                                      │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │  L3: AUTONOMOUS          L2: NOTIFY           L1: APPROVAL      L0: NEVER   │   │
│   │  ────────────────        ─────────            ───────────       ─────────   │   │
│   │                                                                             │   │
│   │  • Read files            • Create tasks       • Send money      • Delete    │   │
│   │  • Search                • Update metrics     • External comms  • Legal     │   │
│   │  • Internal queries      • Log events         • Modify config   • Medical   │   │
│   │  • Status checks         • Routine notifs     • API keys        • Emergency │   │
│   │                                                                             │   │
│   │  Just do it              Do it, tell Zeus     Ask first         Human only  │   │
│   │                                                                             │   │
│   └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
│   L1 APPROVAL FLOW                                                                  │
│   ────────────────                                                                  │
│                                                                                      │
│   Legacy decides to send email                                                      │
│          │                                                                          │
│          ▼                                                                          │
│   ┌─────────────────┐                                                               │
│   │ Check autonomy  │                                                               │
│   │ level → L1      │                                                               │
│   └────────┬────────┘                                                               │
│            │                                                                         │
│            ▼                                                                         │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                    NATS: hitl.approval.request                              │   │
│   │                                                                             │   │
│   │  {                                                                          │   │
│   │    "id": "hitl_12345",                                                      │   │
│   │    "action": "send_email",                                                  │   │
│   │    "target": "client@example.com",                                          │   │
│   │    "content_preview": "Re: Project update...",                              │   │
│   │    "urgency": "medium",                                                     │   │
│   │    "timeout_seconds": 300,                                                  │   │
│   │    "channels": ["voice", "app", "discord"]                                  │   │
│   │  }                                                                          │   │
│   │                                                                             │   │
│   └────────┬───────────────────────┬───────────────────────┬────────────────────┘   │
│            │                       │                       │                        │
│            ▼                       ▼                       ▼                        │
│   ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐              │
│   │     VOICE       │     │      APP        │     │    DISCORD      │              │
│   │                 │     │    (Dioxus)     │     │                 │              │
│   │  "Zeus, may I   │     │  ┌───────────┐  │     │  Legacy wants   │              │
│   │   send this     │     │  │ APPROVE   │  │     │  to send email  │              │
│   │   email to..."  │     │  │  DENY     │  │     │  [✓] [✗]        │              │
│   │                 │     │  └───────────┘  │     │                 │              │
│   └────────┬────────┘     └────────┬────────┘     └────────┬────────┘              │
│            │                       │                       │                        │
│            │                       │                       │                        │
│            ▼                       ▼                       ▼                        │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                    NATS: hitl.approval.response                             │   │
│   │                                                                             │   │
│   │  { "id": "hitl_12345", "approved": true, "channel": "voice" }              │   │
│   │                                                                             │   │
│   └────────┬────────────────────────────────────────────────────────────────────┘   │
│            │                                                                         │
│            ▼                                                                         │
│   ┌─────────────────┐                                                               │
│   │ Execute action  │                                                               │
│   │ OR cancel       │                                                               │
│   └─────────────────┘                                                               │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 HITL Implementation

```python
# agent/hitl/manager.py

from enum import Enum
from dataclasses import dataclass
import asyncio

class AutonomyLevel(Enum):
    L3_AUTONOMOUS = 3  # Just do it
    L2_NOTIFY = 2      # Do it, notify after
    L1_APPROVAL = 1    # Ask first
    L0_FORBIDDEN = 0   # Never (human only)

# Action classification
ACTION_LEVELS = {
    # L3: Autonomous
    "read_file": AutonomyLevel.L3_AUTONOMOUS,
    "search_code": AutonomyLevel.L3_AUTONOMOUS,
    "get_status": AutonomyLevel.L3_AUTONOMOUS,
    "query_memory": AutonomyLevel.L3_AUTONOMOUS,

    # L2: Notify
    "create_task": AutonomyLevel.L2_NOTIFY,
    "update_task": AutonomyLevel.L2_NOTIFY,
    "log_event": AutonomyLevel.L2_NOTIFY,
    "send_notification": AutonomyLevel.L2_NOTIFY,

    # L1: Approval required
    "send_email": AutonomyLevel.L1_APPROVAL,
    "send_message": AutonomyLevel.L1_APPROVAL,
    "transfer_money": AutonomyLevel.L1_APPROVAL,
    "modify_settings": AutonomyLevel.L1_APPROVAL,
    "create_pr": AutonomyLevel.L1_APPROVAL,
    "deploy": AutonomyLevel.L1_APPROVAL,

    # L0: Forbidden
    "delete_data": AutonomyLevel.L0_FORBIDDEN,
    "medical_decision": AutonomyLevel.L0_FORBIDDEN,
    "legal_action": AutonomyLevel.L0_FORBIDDEN,
}

@dataclass
class HITLRequest:
    id: str
    action: str
    details: dict
    urgency: str  # low, medium, high, critical
    timeout_seconds: int
    channels: list  # ["voice", "app", "discord"]

class HITLManager:
    def __init__(self, nats_client):
        self.nats = nats_client
        self.pending_requests = {}

    async def request_approval(self, action: str, details: dict) -> bool:
        """Request human approval for an action."""

        level = ACTION_LEVELS.get(action, AutonomyLevel.L1_APPROVAL)

        if level == AutonomyLevel.L0_FORBIDDEN:
            return False  # Never allowed

        if level == AutonomyLevel.L3_AUTONOMOUS:
            return True  # Just do it

        if level == AutonomyLevel.L2_NOTIFY:
            # Do it, notify after
            asyncio.create_task(self.notify_action(action, details))
            return True

        # L1: Need approval
        request = HITLRequest(
            id=f"hitl_{uuid.uuid4().hex[:8]}",
            action=action,
            details=details,
            urgency="medium",
            timeout_seconds=300,
            channels=["voice", "app", "discord"]
        )

        # Publish request to all channels
        await self.nats.publish(
            "hitl.approval.request",
            json.dumps(asdict(request))
        )

        # Store pending request
        self.pending_requests[request.id] = request

        # Wait for response (with timeout)
        try:
            response = await asyncio.wait_for(
                self.wait_for_response(request.id),
                timeout=request.timeout_seconds
            )
            return response.get("approved", False)
        except asyncio.TimeoutError:
            # No response = denied
            return False

    async def wait_for_response(self, request_id: str) -> dict:
        """Wait for approval response from any channel."""

        # Subscribe to response subject
        sub = await self.nats.subscribe(f"hitl.approval.response.{request_id}")

        async for msg in sub.messages:
            response = json.loads(msg.data)
            del self.pending_requests[request_id]
            return response
```

---

## Section 5: Dioxus Native App

### 5.1 Why Dioxus?

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              DIOXUS FOR ONEMIND APP                                  │
│                                                                                      │
│   WHY DIOXUS?                                                                       │
│   ───────────                                                                       │
│                                                                                      │
│   ✅ Rust-based (performance, safety)                                               │
│   ✅ Cross-platform (iOS, Android, Desktop, Web)                                    │
│   ✅ Single codebase                                                                │
│   ✅ Native performance (not Electron)                                              │
│   ✅ Hot reload for development                                                     │
│   ✅ React-like syntax (familiar)                                                   │
│   ✅ Small binary size                                                              │
│   ✅ WebSocket support (for NATS bridge)                                            │
│                                                                                      │
│   ALTERNATIVES CONSIDERED:                                                          │
│   • Flutter - Good but Dart, larger binaries                                        │
│   • React Native - JS runtime overhead                                              │
│   • Tauri - Desktop only (no mobile)                                                │
│   • SwiftUI + Kotlin - Two codebases                                                │
│                                                                                      │
│   DIOXUS WINS: Native performance + single codebase + Rust safety                  │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 App Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           DIOXUS APP ARCHITECTURE                                    │
│                                                                                      │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                           DIOXUS APP                                        │   │
│   │                                                                             │   │
│   │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │   │
│   │   │    CHAT     │  │    HITL     │  │   STATUS    │  │  SETTINGS   │      │   │
│   │   │    VIEW     │  │   QUEUE     │  │  DASHBOARD  │  │             │      │   │
│   │   │             │  │             │  │             │  │             │      │   │
│   │   │  Messages   │  │  Pending    │  │  Capacity   │  │  Servers    │      │   │
│   │   │  Voice btn  │  │  approvals  │  │  Events     │  │  API keys   │      │   │
│   │   │  History    │  │  [✓] [✗]    │  │  Workflows  │  │  Profile    │      │   │
│   │   │             │  │             │  │             │  │             │      │   │
│   │   └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └─────────────┘      │   │
│   │          │                │                │                               │   │
│   │          └────────────────┼────────────────┘                               │   │
│   │                           │                                                 │   │
│   │                           ▼                                                 │   │
│   │   ┌─────────────────────────────────────────────────────────────────────┐  │   │
│   │   │                      STATE MANAGEMENT                               │  │   │
│   │   │                                                                     │  │   │
│   │   │   messages: Vec<Message>                                            │  │   │
│   │   │   pending_hitl: Vec<HITLRequest>                                    │  │   │
│   │   │   connection_status: ConnectionStatus                               │  │   │
│   │   │   capacity_score: f32                                               │  │   │
│   │   │   recent_events: Vec<Event>                                         │  │   │
│   │   │                                                                     │  │   │
│   │   └────────────────────────────┬────────────────────────────────────────┘  │   │
│   │                                │                                           │   │
│   └────────────────────────────────┼───────────────────────────────────────────┘   │
│                                    │                                               │
│                                    ▼                                               │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                         CONNECTION LAYER                                    │   │
│   │                                                                             │   │
│   │   ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐      │   │
│   │   │   NATS Bridge   │     │  LiveKit SDK    │     │   REST API      │      │   │
│   │   │   (WebSocket)   │     │  (Voice)        │     │   (Fallback)    │      │   │
│   │   │                 │     │                 │     │                 │      │   │
│   │   │  events.>       │     │  Join room      │     │  /api/message   │      │   │
│   │   │  hitl.>         │     │  Push-to-talk   │     │  /api/status    │      │   │
│   │   │  state.>        │     │  Audio stream   │     │                 │      │   │
│   │   │                 │     │                 │     │                 │      │   │
│   │   └────────┬────────┘     └────────┬────────┘     └────────┬────────┘      │   │
│   │            │                       │                       │               │   │
│   └────────────┼───────────────────────┼───────────────────────┼───────────────┘   │
│                │                       │                       │                   │
│                ▼                       ▼                       ▼                   │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                             │   │
│   │   CLOUD                        LIVEKIT CLOUD               GATEWAY          │   │
│   │   (NATS WebSocket)             (Voice)                     (REST)           │   │
│   │                                                                             │   │
│   └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 Dioxus Code Example

```rust
// src/main.rs

use dioxus::prelude::*;

fn main() {
    dioxus::launch(App);
}

fn App() -> Element {
    // Connection state
    let mut connection_status = use_signal(|| ConnectionStatus::Connecting);
    let mut messages = use_signal(|| Vec::<Message>::new());
    let mut pending_hitl = use_signal(|| Vec::<HITLRequest>::new());

    // Connect to NATS via WebSocket bridge
    use_effect(move || {
        spawn(async move {
            let ws = connect_nats_bridge().await;
            connection_status.set(ConnectionStatus::Connected);

            // Subscribe to events
            ws.subscribe("events.legacy.>").await;
            ws.subscribe("hitl.approval.request").await;

            // Handle incoming messages
            while let Some(msg) = ws.next().await {
                match msg.subject.as_str() {
                    s if s.starts_with("events.legacy.response") => {
                        messages.write().push(Message::from_nats(msg));
                    }
                    "hitl.approval.request" => {
                        pending_hitl.write().push(HITLRequest::from_nats(msg));
                    }
                    _ => {}
                }
            }
        });
    });

    rsx! {
        div { class: "app",
            // Header with connection status
            Header { status: connection_status() }

            // Main content
            div { class: "content",
                // Chat view
                ChatView { messages: messages() }

                // HITL approval queue
                HITLQueue {
                    requests: pending_hitl(),
                    on_approve: move |id| approve_hitl(id),
                    on_deny: move |id| deny_hitl(id),
                }
            }

            // Voice button
            VoiceButton {
                on_start: || start_voice_session(),
                on_stop: || stop_voice_session(),
            }
        }
    }
}

#[component]
fn HITLQueue(
    requests: Vec<HITLRequest>,
    on_approve: EventHandler<String>,
    on_deny: EventHandler<String>,
) -> Element {
    rsx! {
        div { class: "hitl-queue",
            h3 { "Pending Approvals" }

            for request in requests {
                div { class: "hitl-item",
                    p { "{request.action}: {request.details}" }

                    div { class: "hitl-buttons",
                        button {
                            class: "approve",
                            onclick: move |_| on_approve.call(request.id.clone()),
                            "✓ Approve"
                        }
                        button {
                            class: "deny",
                            onclick: move |_| on_deny.call(request.id.clone()),
                            "✗ Deny"
                        }
                    }
                }
            }
        }
    }
}
```

---

## Section 6: Letta → TimescaleDB Direct Connection

### 6.1 How Letta Queries TimescaleDB

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                      LETTA → TIMESCALEDB CONNECTION                                  │
│                                                                                      │
│   LETTA doesn't connect directly - it uses a TOOL                                   │
│                                                                                      │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                             │   │
│   │   LETTA (Brain)                                                             │   │
│   │        │                                                                    │   │
│   │        │  "What was Zeus's heart rate last week?"                          │   │
│   │        │                                                                    │   │
│   │        ▼                                                                    │   │
│   │   ┌─────────────────────────────────────────────────────────────────────┐  │   │
│   │   │                                                                     │  │   │
│   │   │   TOOL: query_event_history                                         │  │   │
│   │   │                                                                     │  │   │
│   │   │   def query_event_history(                                          │  │   │
│   │   │       subject_pattern: str,                                         │  │   │
│   │   │       start_date: str,                                              │  │   │
│   │   │       end_date: str = None,                                         │  │   │
│   │   │       aggregation: str = "daily"                                    │  │   │
│   │   │   ) -> str:                                                         │  │   │
│   │   │       """Query historical events from TimescaleDB."""               │  │   │
│   │   │       ...                                                           │  │   │
│   │   │                                                                     │  │   │
│   │   └──────────────────────────────┬──────────────────────────────────────┘  │   │
│   │                                  │                                         │   │
│   │                                  ▼                                         │   │
│   │                          PostgreSQL Connection                             │   │
│   │                          (psycopg2 / asyncpg)                              │   │
│   │                                  │                                         │   │
│   │                                  ▼                                         │   │
│   │                           TIMESCALEDB                                      │   │
│   │                                                                             │   │
│   └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 TimescaleDB Query Tool

```python
# agent/tools/timescale/query.py

import asyncpg
from typing import Optional
import os

TIMESCALE_URL = os.environ.get(
    "TIMESCALE_URL",
    "postgresql://onemind:password@timescaledb:5432/events"
)

def query_event_history(
    subject_pattern: str,
    start_date: str,
    end_date: Optional[str] = None,
    aggregation: str = "daily",
    limit: int = 100
) -> str:
    """
    Query historical events from TimescaleDB.

    Use this when:
    - Zeus asks about past events
    - Analyzing patterns over time
    - Generating reports or summaries
    - Looking up what happened on a specific date

    Args:
        subject_pattern: NATS subject pattern (e.g., "events.biometric.heart_rate%")
        start_date: Start date (YYYY-MM-DD)
        end_date: End date (YYYY-MM-DD), defaults to now
        aggregation: "raw", "hourly", or "daily"
        limit: Maximum results to return

    Returns:
        Formatted query results

    Examples:
        query_event_history(
            subject_pattern="events.biometric.heart_rate%",
            start_date="2026-01-01",
            end_date="2026-01-07",
            aggregation="daily"
        )

        query_event_history(
            subject_pattern="events.github.%",
            start_date="2026-01-15",
            aggregation="raw"
        )
    """
    import asyncio

    async def _query():
        conn = await asyncpg.connect(TIMESCALE_URL)

        try:
            end = end_date or "NOW()"

            if aggregation == "raw":
                query = """
                    SELECT timestamp, subject, data
                    FROM events
                    WHERE subject LIKE $1
                      AND timestamp >= $2::date
                      AND timestamp <= $3::date
                    ORDER BY timestamp DESC
                    LIMIT $4
                """
                rows = await conn.fetch(query, subject_pattern, start_date, end, limit)

                results = []
                for row in rows:
                    results.append({
                        "time": row["timestamp"].isoformat(),
                        "subject": row["subject"],
                        "data": row["data"]
                    })

            elif aggregation == "hourly":
                query = """
                    SELECT
                        time_bucket('1 hour', timestamp) AS hour,
                        COUNT(*) as event_count,
                        AVG((data->>'value')::float) as avg_value
                    FROM events
                    WHERE subject LIKE $1
                      AND timestamp >= $2::date
                      AND timestamp <= $3::date
                    GROUP BY hour
                    ORDER BY hour DESC
                    LIMIT $4
                """
                rows = await conn.fetch(query, subject_pattern, start_date, end, limit)

                results = []
                for row in rows:
                    results.append({
                        "hour": row["hour"].isoformat(),
                        "count": row["event_count"],
                        "avg_value": float(row["avg_value"]) if row["avg_value"] else None
                    })

            else:  # daily
                query = """
                    SELECT
                        time_bucket('1 day', timestamp) AS day,
                        COUNT(*) as event_count,
                        AVG((data->>'value')::float) as avg_value,
                        MIN((data->>'value')::float) as min_value,
                        MAX((data->>'value')::float) as max_value
                    FROM events
                    WHERE subject LIKE $1
                      AND timestamp >= $2::date
                      AND timestamp <= $3::date
                    GROUP BY day
                    ORDER BY day DESC
                    LIMIT $4
                """
                rows = await conn.fetch(query, subject_pattern, start_date, end, limit)

                results = []
                for row in rows:
                    results.append({
                        "day": row["day"].strftime("%Y-%m-%d"),
                        "count": row["event_count"],
                        "avg": round(float(row["avg_value"]), 2) if row["avg_value"] else None,
                        "min": round(float(row["min_value"]), 2) if row["min_value"] else None,
                        "max": round(float(row["max_value"]), 2) if row["max_value"] else None
                    })

            return json.dumps(results, indent=2)

        finally:
            await conn.close()

    return asyncio.run(_query())
```

---

## Section 7: 24/7 Access Strategy

### 7.1 Access Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           24/7 ACCESS ARCHITECTURE                                   │
│                                                                                      │
│   GOAL: Legacy is accessible EVERYWHERE, ALWAYS                                     │
│                                                                                      │
│  ════════════════════════════════════════════════════════════════════════════════   │
│                                                                                      │
│   ACCESS MODES                                                                      │
│   ────────────                                                                      │
│                                                                                      │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                             │   │
│   │   MODE 1: FULL ONLINE                                                       │   │
│   │   ──────────────────                                                        │   │
│   │   • Internet available                                                      │   │
│   │   • Cloud Letta accessible                                                  │   │
│   │   • Full capabilities                                                       │   │
│   │                                                                             │   │
│   │   Phone ──► LiveKit Cloud ──► Cloud NATS ──► Letta ──► Full response       │   │
│   │                                                                             │   │
│   └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                             │   │
│   │   MODE 2: LOCAL NETWORK ONLY                                                │   │
│   │   ──────────────────────────                                                │   │
│   │   • Home network available                                                  │   │
│   │   • No internet                                                             │   │
│   │   • Local Letta fallback                                                    │   │
│   │                                                                             │   │
│   │   Phone ──► Local NATS (Tailscale) ──► Local Letta ──► Response            │   │
│   │                                                                             │   │
│   │   Capabilities:                                                             │   │
│   │   ✅ Voice (local STT/TTS)                                                  │   │
│   │   ✅ Smart home control                                                     │   │
│   │   ✅ Local file access                                                      │   │
│   │   ⚠️ Limited memory (local cache)                                          │   │
│   │   ❌ External APIs (Todoist, GitHub, etc.)                                  │   │
│   │                                                                             │   │
│   └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                             │   │
│   │   MODE 3: FULLY OFFLINE                                                     │   │
│   │   ────────────────────                                                      │   │
│   │   • No network at all                                                       │   │
│   │   • Local-only model                                                        │   │
│   │   • Queue actions for later                                                 │   │
│   │                                                                             │   │
│   │   Phone (local) ──► On-device model ──► Queue for sync                     │   │
│   │                                                                             │   │
│   │   Capabilities:                                                             │   │
│   │   ✅ Basic voice commands                                                   │   │
│   │   ✅ Note taking (queued)                                                   │   │
│   │   ✅ Reminders (local)                                                      │   │
│   │   ⚠️ Simple responses (small model)                                        │   │
│   │   ❌ Full reasoning                                                         │   │
│   │   ❌ Memory access                                                          │   │
│   │                                                                             │   │
│   │   SYNC: When online, queue syncs automatically                             │   │
│   │                                                                             │   │
│   └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Multi-Access Points

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         Legacy ACCESS POINTS                                         │
│                                                                                      │
│   EVERYWHERE Legacy LIVES                                                           │
│   ───────────────────────                                                           │
│                                                                                      │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│   │   PHONE     │  │   LAPTOP    │  │   DESKTOP   │  │   GLASSES   │              │
│   │   (Dioxus)  │  │   (Web)     │  │   (Dioxus)  │  │   (Future)  │              │
│   │             │  │             │  │             │  │             │              │
│   │  • Voice    │  │  • Voice    │  │  • Voice    │  │  • Voice    │              │
│   │  • Chat     │  │  • Chat     │  │  • Chat     │  │  • Vision   │              │
│   │  • HITL     │  │  • HITL     │  │  • HITL     │  │  • Always-on│              │
│   │  • Offline  │  │             │  │             │  │             │              │
│   │    queue    │  │             │  │             │  │             │              │
│   └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │
│          │                │                │                │                      │
│          └────────────────┴────────────────┴────────────────┘                      │
│                                    │                                                │
│                                    ▼                                                │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                             │   │
│   │   UNIFIED IDENTITY                                                          │   │
│   │   ────────────────                                                          │   │
│   │                                                                             │   │
│   │   • Same Letta agent ID across all access points                           │   │
│   │   • Conversation continues across devices                                   │   │
│   │   • Memory is unified in cloud                                             │   │
│   │   • Context syncs automatically                                            │   │
│   │                                                                             │   │
│   │   Example:                                                                  │   │
│   │   Morning (phone): "Legacy, remind me to call John"                        │   │
│   │   Afternoon (laptop): "Legacy, what did I ask you this morning?"           │   │
│   │   Legacy: "You asked me to remind you to call John."                       │   │
│   │                                                                             │   │
│   └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
│   ADDITIONAL ACCESS POINTS                                                          │
│   ────────────────────────                                                          │
│                                                                                      │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│   │   DISCORD   │  │  TELEGRAM   │  │   WATCH     │  │   CAR       │              │
│   │             │  │             │  │   (Future)  │  │   (Future)  │              │
│   │  Bot in     │  │  Bot for    │  │             │  │             │              │
│   │  your       │  │  mobile     │  │  Quick      │  │  CarPlay/   │              │
│   │  server     │  │  messaging  │  │  commands   │  │  Android    │              │
│   │             │  │             │  │             │  │  Auto       │              │
│   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘              │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 7.3 Offline Queue System

```python
# agent/offline/queue.py

import sqlite3
from dataclasses import dataclass
from datetime import datetime
from typing import Optional

@dataclass
class QueuedAction:
    id: str
    action_type: str  # "message", "task", "note", "reminder"
    content: str
    created_at: datetime
    synced_at: Optional[datetime] = None

class OfflineQueue:
    """Local queue for offline actions that sync when online."""

    def __init__(self, db_path: str = "offline_queue.db"):
        self.conn = sqlite3.connect(db_path)
        self._init_db()

    def _init_db(self):
        self.conn.execute("""
            CREATE TABLE IF NOT EXISTS queue (
                id TEXT PRIMARY KEY,
                action_type TEXT NOT NULL,
                content TEXT NOT NULL,
                created_at TIMESTAMP NOT NULL,
                synced_at TIMESTAMP
            )
        """)
        self.conn.commit()

    def add(self, action: QueuedAction):
        """Add action to offline queue."""
        self.conn.execute(
            "INSERT INTO queue VALUES (?, ?, ?, ?, ?)",
            (action.id, action.action_type, action.content,
             action.created_at, action.synced_at)
        )
        self.conn.commit()

    def get_pending(self) -> list[QueuedAction]:
        """Get all unsynced actions."""
        cursor = self.conn.execute(
            "SELECT * FROM queue WHERE synced_at IS NULL ORDER BY created_at"
        )
        return [QueuedAction(*row) for row in cursor.fetchall()]

    def mark_synced(self, action_id: str):
        """Mark action as synced."""
        self.conn.execute(
            "UPDATE queue SET synced_at = ? WHERE id = ?",
            (datetime.utcnow(), action_id)
        )
        self.conn.commit()

    async def sync_when_online(self, nats_client):
        """Sync all pending actions when connection restored."""
        pending = self.get_pending()

        for action in pending:
            try:
                # Publish to NATS
                await nats_client.publish(
                    f"events.offline.{action.action_type}",
                    json.dumps({
                        "id": action.id,
                        "content": action.content,
                        "created_at": action.created_at.isoformat(),
                        "source": "offline_queue"
                    })
                )

                self.mark_synced(action.id)

            except Exception as e:
                # Will retry next sync
                print(f"Failed to sync {action.id}: {e}")
```

---

## Section 8: Complete Real-Time Data Flow

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                       COMPLETE REAL-TIME DATA FLOW                                   │
│                                                                                      │
│   ZEUS SAYS: "Legacy, cancel the garden watering and create a task to fix it"      │
│                                                                                      │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │  T+0ms     Voice captured by LiveKit                                        │   │
│   │            │                                                                │   │
│   │  T+50ms    VAD detects speech end                                          │   │
│   │            │                                                                │   │
│   │  T+200ms   STT transcription complete                                      │   │
│   │            │                                                                │   │
│   │  T+210ms   Command processor detects "cancel" → INTERRUPT                  │   │
│   │            │                                                                │   │
│   │  T+220ms   ├── NATS: commands.interrupt                                    │   │
│   │            │   ├── Kestra receives → cancels garden workflow               │   │
│   │            │   ├── TTS stops (if speaking)                                 │   │
│   │            │   └── Legacy acknowledges                                     │   │
│   │            │                                                                │   │
│   │  T+250ms   Full transcript sent to Letta                                   │   │
│   │            │                                                                │   │
│   │  T+300ms   Letta processes request                                         │   │
│   │            │                                                                │   │
│   │  T+800ms   Letta decides: create_task (L2: NOTIFY)                        │   │
│   │            │                                                                │   │
│   │  T+850ms   Task created in Todoist                                         │   │
│   │            │                                                                │   │
│   │  T+900ms   NATS: events.task.created (notification)                        │   │
│   │            │                                                                │   │
│   │  T+1000ms  Letta response: "Cancelled. Created task: Fix garden..."       │   │
│   │            │                                                                │   │
│   │  T+1100ms  TTS begins streaming response                                   │   │
│   │            │                                                                │   │
│   │  T+1500ms  Zeus hears response                                             │   │
│   │            │                                                                │   │
│   │  T+2000ms  Complete. Archive to TimescaleDB                                │   │
│   │                                                                             │   │
│   └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
│   TOTAL: ~2 seconds voice-to-voice with cancellation + task creation               │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

*Document created: January 2026*
*Purpose: Real-time access architecture for OneMind OS*
*Key insight: Legacy must be EVERYWHERE, ALWAYS, with multiple access modes*
