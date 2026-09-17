# Unified Interface Architecture

> **Status: Problem Statement**
> This document describes the **current disconnected architecture** and the problem to solve.
> For the proposed solution, see [VOICE-FIRST-SWARM-ARCHITECTURE.md](VOICE-FIRST-SWARM-ARCHITECTURE.md).

## The Question

How do Voice, Chat, and Vision all work together using Agno's native patterns?

---

## Current Structure (What You Have)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          CURRENT FLOW                                    │
│                                                                         │
│   VOICE                    CHAT                      VISION             │
│   (LiveKit)               (REST API)                (Multimodal)        │
│      │                       │                          │               │
│      ▼                       ▼                          ▼               │
│  ┌─────────┐           ┌─────────────┐           ┌─────────────┐        │
│  │ livekit │           │ /agents/    │           │ /multimodal │        │
│  │  .py    │           │ {id}/runs   │           │   /analyze  │        │
│  └────┬────┘           └──────┬──────┘           └──────┬──────┘        │
│       │                       │                         │               │
│       │                       │                         │               │
│       ▼                       ▼                         ▼               │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │              DISCONNECTED PATHS                                  │   │
│  │                                                                  │   │
│  │  Voice: AgnoLLM → create_legacy_agent() → Single Agent          │   │
│  │  Chat:  API → agent.arun() → Single Agent                       │   │
│  │  Vision: analyze_image() → Separate processing                  │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│   PROBLEM: Three separate paths, no unified routing, no swarm          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Option A: Legacy as ROUTE-Mode Team

Make Legacy a Team that routes to specialists:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    LEGACY AS AGNO TEAM (ROUTE MODE)                     │
│                                                                         │
│                          USER INPUT                                     │
│                              │                                          │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                    UNIFIED API ENDPOINT                            │  │
│  │                                                                    │  │
│  │   POST /api/command                                                │  │
│  │   {                                                                │  │
│  │     "message": "Research Flutter animations",                      │  │
│  │     "media": [{"type": "image", "data": "..."}],  // Optional     │  │
│  │     "source": "voice" | "chat" | "vision",                         │  │
│  │     "session_id": "abc123"                                         │  │
│  │   }                                                                │  │
│  │                                                                    │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                          │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                         LEGACY TEAM                                │  │
│  │                       (mode=ROUTE)                                 │  │
│  │                                                                    │  │
│  │   Team(                                                            │  │
│  │     name="Legacy",                                                 │  │
│  │     mode=TeamMode.ROUTE,                                           │  │
│  │     members=[                                                      │  │
│  │       researcher_agent,                                            │  │
│  │       coder_agent,                                                 │  │
│  │       writer_agent,                                                │  │
│  │       analyst_agent,                                               │  │
│  │       vision_agent,                                                │  │
│  │       planner_agent,                                               │  │
│  │       ... (all specialists)                                        │  │
│  │     ],                                                             │  │
│  │     instructions=LEGACY_SYSTEM_PROMPT,                             │  │
│  │   )                                                                │  │
│  │                                                                    │  │
│  │   How ROUTE mode works:                                            │  │
│  │   1. Team leader (coordinator) receives request                    │  │
│  │   2. Leader analyzes and picks best member(s)                      │  │
│  │   3. Routes task to selected member(s)                             │  │
│  │   4. Returns response (synthesized if multiple)                    │  │
│  │                                                                    │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                          │
│           ┌──────────────────┼──────────────────┐                       │
│           ▼                  ▼                  ▼                       │
│     ┌──────────┐      ┌──────────┐      ┌──────────┐                    │
│     │Researcher│      │  Coder   │      │  Vision  │                    │
│     │  Agent   │      │  Agent   │      │  Agent   │                    │
│     └──────────┘      └──────────┘      └──────────┘                    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Pros**: Uses Agno's native Team routing. Clean architecture.
**Cons**: All specialists must be loaded as members. Less flexible.

---

## Option B: Legacy as Agent with Coordination Tools

Keep Legacy as a single Agent, but give it tools to spawn agents/teams:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                LEGACY AS AGENT WITH COORDINATION TOOLS                  │
│                                                                         │
│                          USER INPUT                                     │
│                              │                                          │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                    UNIFIED API ENDPOINT                            │  │
│  │                      POST /api/command                             │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                          │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                      LEGACY AGENT                                  │  │
│  │                                                                    │  │
│  │   Agent(                                                           │  │
│  │     name="Legacy",                                                 │  │
│  │     tools=[                                                        │  │
│  │       # Standard tools                                             │  │
│  │       todoist, github, browser, comms,                             │  │
│  │                                                                    │  │
│  │       # COORDINATION TOOLS (new)                                   │  │
│  │       ask_specialist,     # Route to single agent                  │  │
│  │       run_team,           # Spawn a team (dev, research, etc)      │  │
│  │       execute_workflow,   # Run a workflow                         │  │
│  │       analyze_image,      # Vision analysis                        │  │
│  │                                                                    │  │
│  │       # PLATFORM TOOLS                                             │  │
│  │       create_task, get_tasks,                                      │  │
│  │       create_page, get_pages,                                      │  │
│  │       query_memory, store_memory,                                  │  │
│  │       search_knowledge,                                            │  │
│  │     ],                                                             │  │
│  │     instructions=LEGACY_SYSTEM_PROMPT,                             │  │
│  │   )                                                                │  │
│  │                                                                    │  │
│  │   How it works:                                                    │  │
│  │   1. Legacy receives request                                       │  │
│  │   2. Legacy decides: handle directly OR use tool                   │  │
│  │   3. Tools spawn agents/teams/workflows as needed                  │  │
│  │   4. Legacy synthesizes and responds                               │  │
│  │                                                                    │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                          │
│                         (via tools)                                     │
│           ┌──────────────────┼──────────────────┐                       │
│           ▼                  ▼                  ▼                       │
│     ┌──────────┐      ┌──────────┐      ┌──────────┐                    │
│     │ask_      │      │run_team  │      │execute_  │                    │
│     │specialist│      │("dev")   │      │workflow  │                    │
│     └──────────┘      └──────────┘      └──────────┘                    │
│           │                  │                  │                       │
│           ▼                  ▼                  ▼                       │
│     ┌──────────┐      ┌──────────┐      ┌──────────┐                    │
│     │ Single   │      │  Team    │      │ Workflow │                    │
│     │ Agent    │      │ (9 types)│      │ Executor │                    │
│     └──────────┘      └──────────┘      └──────────┘                    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Pros**: Flexible. Legacy decides when to spawn. Uses existing teams/workflows.
**Cons**: Legacy must be smart enough to know when to use tools.

---

## Option C: Hybrid - Agent with Fallback to Team

Legacy is an Agent for simple requests, but escalates to Team for complex ones:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    HYBRID: AGENT + TEAM ESCALATION                      │
│                                                                         │
│                          USER INPUT                                     │
│                              │                                          │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                    UNIFIED API ENDPOINT                            │  │
│  │                      POST /api/command                             │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                          │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                      ROUTER                                        │  │
│  │                                                                    │  │
│  │   def route(message, media):                                       │  │
│  │       # Check for explicit @mentions                               │  │
│  │       if has_mentions(message):                                    │  │
│  │           return TEAM_MODE  # Multiple agents needed               │  │
│  │                                                                    │  │
│  │       # Check for complexity signals                               │  │
│  │       if is_complex(message):  # research, multi-step, etc        │  │
│  │           return TEAM_MODE                                         │  │
│  │                                                                    │  │
│  │       # Check for media                                            │  │
│  │       if has_media(media):                                         │  │
│  │           return TEAM_MODE  # Need vision agent                    │  │
│  │                                                                    │  │
│  │       # Simple request - direct to Legacy agent                    │  │
│  │       return AGENT_MODE                                            │  │
│  │                                                                    │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                    │                      │                             │
│           AGENT_MODE                TEAM_MODE                           │
│                    │                      │                             │
│                    ▼                      ▼                             │
│  ┌──────────────────────┐   ┌──────────────────────────────────┐       │
│  │    LEGACY AGENT      │   │        LEGACY TEAM               │       │
│  │                      │   │       (mode=ROUTE)               │       │
│  │  Fast, direct        │   │                                  │       │
│  │  Single model call   │   │  Coordinates specialists         │       │
│  │  For simple tasks    │   │  Parallel execution              │       │
│  │                      │   │  For complex tasks               │       │
│  └──────────────────────┘   └──────────────────────────────────┘       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Pros**: Fast for simple requests. Full power for complex ones.
**Cons**: More routing logic to maintain.

---

## The Unified Endpoint

Regardless of which option, here's the unified endpoint:

```python
# backend/agno/api/command.py

@router.post("/api/command")
async def unified_command(request: CommandRequest) -> StreamingResponse:
    """
    Unified endpoint for all user interactions.

    Accepts voice transcripts, chat messages, and images.
    Routes to appropriate handler based on content.
    """
    # Build context
    context = CommandContext(
        message=request.message,
        media=request.media,          # Images, audio frames
        source=request.source,        # "voice" | "chat" | "vision"
        session_id=request.session_id,
        user_id=request.user_id,
    )

    # Route to handler
    if request.source == "voice":
        # Voice: needs streaming for TTS
        return await stream_voice_response(context)
    else:
        # Chat/Vision: standard streaming
        return await stream_chat_response(context)


async def stream_chat_response(context: CommandContext):
    """Stream response from Legacy (Agent or Team)."""

    # Parse for routing signals
    parsed = parse_message(context.message)

    # Decide: Agent or Team?
    if needs_team(parsed, context.media):
        # Complex request - use Team
        legacy_team = get_legacy_team()
        async for chunk in legacy_team.arun_stream(
            context.message,
            images=context.media,
        ):
            yield chunk
    else:
        # Simple request - use Agent
        legacy_agent = get_legacy_agent()
        async for chunk in legacy_agent.arun_stream(context.message):
            yield chunk
```

---

## Voice Pipeline Integration

Voice connects to the same unified endpoint:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      VOICE FLOW                                         │
│                                                                         │
│   ┌─────────┐     ┌─────────┐     ┌───────────────┐     ┌─────────┐    │
│   │ Phone   │     │ LiveKit │     │ Voice Worker  │     │  TTS    │    │
│   │  Mic    │────▶│ Server  │────▶│               │────▶│(Cartesia│    │
│   └─────────┘     └─────────┘     │ 1. STT        │     └─────────┘    │
│                                   │ 2. Call API:  │                     │
│                                   │    POST /api/ │                     │
│                                   │    command    │                     │
│                                   │ 3. Get text   │                     │
│                                   └───────────────┘                     │
│                                          │                              │
│                                          ▼                              │
│                    ┌─────────────────────────────────────────┐          │
│                    │  SAME UNIFIED ENDPOINT                  │          │
│                    │  (Legacy Agent or Team handles it)      │          │
│                    └─────────────────────────────────────────┘          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

The voice worker just needs to:
1. STT the audio
2. Call the unified endpoint
3. Stream text back to TTS

---

## My Recommendation: Option B (Agent with Coordination Tools)

Why:
1. **Uses existing infrastructure** - Your teams, workflows already work
2. **Flexible** - Legacy decides dynamically what to spawn
3. **Simpler** - One agent, not a massive team with all specialists
4. **Matches your mental model** - "Legacy IS the swarm" = Legacy can USE the swarm

```python
# What needs to happen:

# 1. Add coordination tools to Legacy's preset
LEGACY_TOOLS = [
    # Existing
    "todoist", "github", "browser", "comms", "inbox",

    # NEW: Coordination tools
    "ask_specialist",    # @tool: route to specific agent
    "run_team",          # @tool: spawn team (dev, research, etc)
    "execute_workflow",  # @tool: run workflow

    # NEW: Platform awareness tools
    "create_task", "get_tasks",
    "create_page", "get_pages",
    "analyze_image",     # Vision
]

# 2. Create unified endpoint
POST /api/command

# 3. Wire voice worker to call unified endpoint
```

---

## Next Steps

1. **Create coordination tools** (`ask_specialist`, `run_team`, `execute_workflow`)
2. **Add to Legacy's preset**
3. **Create unified `/api/command` endpoint**
4. **Wire voice worker to use it**

This keeps Agno's architecture intact while giving Legacy swarm capabilities.
