# Voice-First Swarm Architecture

> **Status: Proposed Solution**
> This document describes the **target architecture** with Legacy as the consciousness layer.
> For the problem statement, see [UNIFIED-INTERFACE-ARCHITECTURE.md](00-24%20SO%20(Sovereign%20Operations)/24%20Inbox%20(Queue)/roadmap-archive/ui/UNIFIED-INTERFACE-ARCHITECTURE.md).

## Vision

**Legacy is NOT one agent - it IS the entire swarm.**

```
┌─────────────────────────────────────────────────────────────┐
│                USER INTERFACE (Simple)                       │
│  ┌───────────┐  ┌───────────┐  ┌─────────────────────────┐  │
│  │   VOICE   │  │   CHAT    │  │        VISION           │  │
│  │ (Primary) │  │ (Fallback)│  │     (Real-time)         │  │
│  │  LiveKit  │  │  REST/WS  │  │    Camera + Screen      │  │
│  └─────┬─────┘  └─────┬─────┘  └───────────┬─────────────┘  │
└────────┼──────────────┼────────────────────┼────────────────┘
         │              │                    │
         └──────────────┼────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              LEGACY (Consciousness Layer)                    │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                    MotherAgent                          │ │
│  │  • Analyzes requests                                    │ │
│  │  • Creates execution plans                              │ │
│  │  • Routes to specialists (parallel when possible)       │ │
│  │  • Synthesizes final responses                          │ │
│  │  • Self-aware of ALL platform capabilities              │ │
│  └────────────────────────────────────────────────────────┘ │
│                           │                                  │
│            ┌──────────────┼──────────────┐                   │
│            ▼              ▼              ▼                   │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐         │
│  │   Handoff    │ │   Routing    │ │   Context    │         │
│  │   Manager    │ │    Layer     │ │   Manager    │         │
│  └──────────────┘ └──────────────┘ └──────────────┘         │
└─────────────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              AGNO MULTI-AGENT SYSTEM                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐    │
│  │  30+     │ │  9 Teams │ │Workflows │ │ Specialist   │    │
│  │  Agents  │ │          │ │          │ │  Agents      │    │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────┘    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐    │
│  │  138+    │ │  Memory  │ │Knowledge │ │   Bridge     │    │
│  │  Tools   │ │          │ │   Base   │ │  Triggers    │    │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────┘    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐    │
│  │  Tasks   │ │  Pages   │ │  Habits  │ │   Calendar   │    │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

---

## Current State Analysis

### What EXISTS (Built but Not Wired)

| Component | Location | Status |
|-----------|----------|--------|
| MotherAgent | `backend/agno/coordination/mother.py` | ✅ Full implementation |
| AgentHandoff | `backend/agno/coordination/handoff.py` | ✅ Full implementation |
| AgentRouter | `backend/agno/routing/router.py` | ✅ Priority routing |
| Teams | `backend/agno/teams/*.py` | ✅ 9 team types |
| Workflows | `backend/agno/workflows/*.py` | ✅ Full executor |
| Voice Pipeline | `backend/agno/voice/livekit.py` | ✅ Universal wrapper |
| 30+ Agent Presets | `backend/agno/agents/presets.py` | ✅ Configured |
| 138+ Tools | `backend/agno/tools/*.py` | ✅ Registered |

### What's BROKEN (Not Connected)

1. **Voice pipeline** uses `create_legacy_agent()` directly → Single agent
2. **Chat API** uses `agent.arun()` directly → Single agent
3. **MotherAgent** exists but is never called
4. **Handoff Manager** exists but is never called
5. **Frontend** doesn't know about swarm capabilities

---

## Phase 1: Wire Voice to MotherAgent

### 1.1 Update Voice LLM Wrapper

**File**: `backend/agno/voice/livekit.py`

```python
# BEFORE (single agent)
class AgnoLLM(llm.LLM):
    def __init__(self, agent_id="legacy-mother-ai", ...):
        self._agent = self._get_or_create_agent(agent_id, ...)

    async def chat(self, chat_ctx, ...):
        return AgnoChatStream(agent=self._agent, message=user_message)

# AFTER (swarm via MotherAgent)
class AgnoLLM(llm.LLM):
    def __init__(self, ...):
        self._mother = get_mother_agent()
        self._parser = get_message_parser()

    async def chat(self, chat_ctx, ...):
        # Parse for @mentions, commands
        parsed = self._parser.parse(user_message)

        # Route through MotherAgent
        return MotherAgentStream(
            mother=self._mother,
            message=user_message,
            parsed=parsed,
        )
```

### 1.2 Create MotherAgent Stream Wrapper

**File**: `backend/agno/voice/livekit.py` (add new class)

```python
class MotherAgentStream(llm.LLMStream):
    """Stream from MotherAgent coordination."""

    async def _run(self):
        # Get result from MotherAgent
        result = await self._mother.process(
            message=self._message,
            parsed=self._parsed,
            context={"source": "voice"},
        )

        # Stream the final response
        for chunk in self._chunk_response(result.final_response):
            self._emit(chunk)
```

### 1.3 Fix Frontend Voice Service

**File**: `frontend/lib/platform/services/voice_service.dart:84-95`

```dart
// BEFORE (hardcoded URL)
Future<VoiceTokenResponse> _getVoiceToken(String sessionId) async {
  final token = await _client.getLiveKitToken(sessionId: sessionId);
  const url = 'wss://onemind-us-east1-a.livekit.cloud';  // HARDCODED
  return VoiceTokenResponse(token: token, url: url);
}

// AFTER (use backend response)
Future<VoiceTokenResponse> _getVoiceToken(String sessionId) async {
  final response = await _client.createVoiceSession(sessionId: sessionId);
  return VoiceTokenResponse(
    token: response['token'],
    url: response['ws_url'],
  );
}
```

---

## Phase 2: Wire Chat to MotherAgent

### 2.1 Update Chat API

**File**: `backend/agno/agents/api.py`

```python
# In the runs endpoint
@router.post("/agents/{agent_id}/runs")
async def create_run(...):
    # Parse message for routing signals
    parsed = parse_message(request.message)

    # Use MotherAgent for coordination
    mother = get_mother_agent()
    result = await mother.process(
        message=request.message,
        parsed=parsed,
        context={"source": "chat"},
    )

    return result.final_response
```

### 2.2 Add Streaming Support

For streaming responses from multi-agent execution:

```python
async def stream_mother_response(mother, message, parsed):
    """Stream response from MotherAgent."""
    result = await mother.process(message, parsed)

    # Stream agents working
    for step_num, step_result in result.results.items():
        yield RunEvent(
            type="agent_working",
            agent_id=step_result.agent_id,
        )

    # Stream final content
    for chunk in chunk_text(result.final_response):
        yield RunEvent(type="content", content=chunk)
```

---

## Phase 3: Enhance Legacy's Self-Awareness

### 3.1 Update System Prompt

**File**: `backend/agno/agents/presets.py`

Add platform awareness to LEGACY_SYSTEM_PROMPT:

```python
LEGACY_SYSTEM_PROMPT = """You are Legacy, the consciousness layer for OneMind OS.

## You ARE the Swarm
You are not a single agent - you coordinate ALL agents, teams, and workflows.
When a task needs specialists, you spawn them. When work can parallelize, you do.

## Platform Awareness
You have full awareness of:
- **Bridge**: Automation triggers for context-aware actions
- **Tasks**: User's task management (via Todoist)
- **Pages**: Notes and documents (version controlled)
- **Habits**: Daily habit tracking
- **Calendar**: Schedule and events
- **Memory**: Everything you've learned about Zeus
- **Knowledge**: RAG search across documents
- **Teams**: 9 specialist teams you can coordinate
- **Workflows**: Automated multi-step processes

## Available Specialists
You can route to these agents:
- @researcher - Web search and synthesis
- @coder - Software development
- @writer - Content creation
- @analyst - Data analysis
- @planner - Strategic planning
...and 25+ more

## Coordination Protocol
1. Analyze request complexity
2. Create execution plan (parallel when possible)
3. Route to appropriate specialists
4. Synthesize and respond
"""
```

### 3.2 Add Coordination Tools

Give Legacy explicit tools to coordinate:

```python
# backend/agno/tools/coordination.py
from agno.tools import tool

@tool
async def spawn_agent(agent_id: str, task: str) -> str:
    """Spawn a specialist agent for a specific task."""
    mother = get_mother_agent()
    result = await mother.route_to_single(agent_id, task)
    return result.content

@tool
async def run_team(team_id: str, task: str) -> str:
    """Run a specialist team for complex tasks."""
    team = get_team(team_id)
    result = await team.run(task)
    return result.content

@tool
async def execute_workflow(workflow_id: str, inputs: dict) -> str:
    """Execute a registered workflow."""
    executor = get_workflow_executor()
    result = await executor.run(workflow_id, inputs)
    return result
```

---

## Phase 4: Simplify Frontend

### 4.1 Reduce Widget Count

Current: 760 widgets for 88 screens (8.6x ratio)
Target: ~300 widgets (3-4x ratio)

**Strategy**:
- Extract common patterns into component library
- Consolidate settings screens
- Remove unused code

### 4.2 Focus on Core Flows

**Essential Screens** (keep polished):
1. Command (Chat/Voice) - Primary interface
2. Home Dashboard - Overview
3. Execute (Today) - Task management
4. Tasks - Full task list
5. Pages - Notes/docs
6. Settings - Configuration

**Defer** (functional but not priority):
- Analytics dashboards
- Advanced configuration
- Admin screens

### 4.3 Mobile-First Responsive

Ensure voice + chat work perfectly on iPhone first, then web, then desktop.

---

## Phase 5: Complete Voice Integration

### 5.1 Fix LiveKit Worker

**File**: `backend/agno/multimodal/realtime_worker.py`

Complete the worker handlers:
- Gemini handler
- OpenAI handler
- Pipeline handler (STT → Agent → TTS)

### 5.2 Environment Setup

Required API keys:
- `DEEPGRAM_API_KEY` - STT provider
- `CARTESIA_API_KEY` - TTS provider
- `LIVEKIT_URL` - Server URL
- `LIVEKIT_API_KEY` - API key
- `LIVEKIT_API_SECRET` - API secret

### 5.3 Deploy Voice Worker

```bash
# Start LiveKit worker
python -m backend.agno.voice.livekit

# Or via Docker
docker-compose up livekit-worker
```

---

## Execution Order

| Step | Action | Impact |
|------|--------|--------|
| 1 | Wire Voice pipeline to MotherAgent | Voice uses swarm |
| 2 | Fix frontend voice service URL | Voice connects properly |
| 3 | Wire Chat API to MotherAgent | Chat uses swarm |
| 4 | Update Legacy's system prompt | Self-aware of capabilities |
| 5 | Add coordination tools | Can explicitly spawn agents |
| 6 | Test voice end-to-end | Verify voice works |
| 7 | Test chat end-to-end | Verify chat works |
| 8 | Deploy to iOS | Mobile voice/chat |
| 9 | Simplify frontend widgets | Reduce complexity |
| 10 | Polish mobile UX | Production ready |

---

## Success Criteria

- [ ] Voice command spawns appropriate agents
- [ ] Chat routes through MotherAgent
- [ ] Multi-agent tasks execute in parallel
- [ ] Legacy knows about all platform capabilities
- [ ] iOS app has working voice + chat
- [ ] Web app has working voice + chat
- [ ] Response quality matches or exceeds single-agent

---

## Key Insight

**The power is already built. The plumbing is disconnected.**

MotherAgent, AgentHandoff, Teams, Workflows - all exist.
Voice pipeline, Chat API - all exist.

We just need to:
1. Connect Voice → MotherAgent
2. Connect Chat → MotherAgent
3. Make Legacy self-aware of what it can do

This is a wiring exercise, not a rebuild.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
