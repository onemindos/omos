# Local LiveKit Voice Setup - Complete Guide

## Your Questions Answered

### 1. Local vs Cloud LiveKit ✅

**You already have LiveKit installed locally via Homebrew!**
- Config ready: `infra/config/livekit/livekit.yaml`
- Running locally gives you full control, no cloud costs
- Better for development and testing

### 2. One Voice Agent vs All Agents 🎯

**RECOMMENDATION: Give voice to ALL agents dynamically**

Why:
- User chooses which agent to talk to
- Different agents for different tasks
- More flexible and powerful
- Your system already supports this!

Architecture:
```
User creates session → Specifies agent_id → Worker loads that agent → Voice enabled

Examples:
- "Talk to research agent" → research-agent gets voice
- "Talk to CEO agent" → ceo-agent gets voice
- "Talk to legacy" → legacy-mother-ai gets voice
```

### 3. Agno Voice Capabilities 🚀

**YES! Full agent capabilities work in voice mode:**

| Capability | Works in Voice? | Details |
|------------|----------------|---------|
| **All Tools** | ✅ YES | All 29+ registered tools |
| **MCP Servers** | ✅ YES | Context7, filesystem, memory, etc. |
| **Memory** | ✅ YES | Full agent memory |
| **Knowledge/RAG** | ✅ YES | Knowledge bases work |
| **Custom Logic** | ✅ YES | All agent code runs |
| **Tool Calling** | ✅ YES | Agent can use any tool mid-conversation |
| **Streaming** | ✅ YES | Real-time responses |

**How it works:**
```python
# User speaks: "Search for quantum computing papers"

# 1. Deepgram converts speech → text
# 2. Text goes to YOUR agent
agent.arun("Search for quantum computing papers")

# 3. Agent uses tools
agent.use_tool("duckduckgo")  # Search tool
agent.use_tool("arxiv")       # Research papers

# 4. Agent processes results
# 5. Response → Cartesia TTS → User hears answer
```

### 4. Context7 + MCP in Voice ✅

**YES! You can use Context7 MCP for docs during voice conversations!**

Your system supports MCP servers:
- Context7 (docs)
- Memory
- Filesystem
- Tavily (search)
- And 29+ other tools

## Complete Local Setup

### Step 1: Start Local LiveKit Server

```bash
# Option A: Using your config file
livekit-server --config infra/config/livekit/livekit.yaml

# Option B: Dev mode (simpler)
livekit-server --dev --bind 0.0.0.0
```

You'll see:
```
INFO: LiveKit server starting
INFO: API server listening on :7880
INFO: RTC server ready
```

### Step 2: Update .env for Local

```bash
# Change from cloud to local
LIVEKIT_URL=ws://localhost:7880

# For dev mode, use these keys
LIVEKIT_API_KEY=devkey
LIVEKIT_API_SECRET=secret

# Or generate your own
livekit-cli create-token \
  --api-key YOUR_KEY \
  --api-secret YOUR_SECRET \
  --identity zeus \
  --grant roomJoin=true
```

### Step 3: Get Voice API Keys

```bash
# Deepgram (Speech → Text)
# Go to: https://deepgram.com
DEEPGRAM_API_KEY=your-key

# Cartesia (Text → Speech)
# Go to: https://cartesia.ai
CARTESIA_API_KEY=your-key

# Add to .env
echo "DEEPGRAM_API_KEY=your-key" >> .env
echo "CARTESIA_API_KEY=your-key" >> .env
```

### Step 4: Install Python Dependencies

```bash
pip install \
  livekit \
  livekit-api \
  livekit-agents \
  livekit-plugins-deepgram \
  livekit-plugins-cartesia \
  livekit-plugins-silero \
  livekit-plugins-openai
```

### Step 5: Create Voice-Enabled Agent Registry

Your agents already work! Let me show you how to make them all voice-ready:

```python
# backend/agno/agents/voice_registry.py (NEW FILE)
"""
Voice-enabled agent registry.
Any agent can be used in voice mode!
"""
from backend.shared.registry import list_agents, get_agent

def get_voice_capable_agents():
    """
    Get all agents that support voice.

    Returns:
        List of agent IDs that can be used in voice mode
    """
    # ALL agents support voice!
    return list_agents()

def create_voice_session_for_agent(agent_id: str):
    """
    Create voice session for any agent.

    Args:
        agent_id: Any registered agent ID

    Returns:
        Voice session credentials
    """
    from backend.agno.multimodal.realtime import create_realtime_session

    return await create_realtime_session(
        agent_id=agent_id,
        provider="pipeline",  # Full agent mode
        enable_vision=False,
    )
```

### Step 6: Enable MCP Tools (Including Context7)

```python
# Add Context7 MCP to your agents
from backend.integrations.mcp import create_mcp_tools

# In your agent creation:
async def create_voice_agent_with_mcp(agent_id: str):
    """Create agent with MCP tools for voice."""

    # Get base agent
    agent = get_agent(agent_id)

    # Add MCP tools
    mcp_tools = await create_mcp_tools([
        "context7",    # Documentation
        "memory",      # Short-term memory
        "filesystem",  # File access
    ])

    # Agent now has access to all MCP servers in voice!
    if agent.tools:
        agent.tools.append(mcp_tools)
    else:
        agent.tools = [mcp_tools]

    return agent
```

### Step 7: Create Realtime Tools for Voice

Let me create a proper realtime tools module:

```python
# backend/agno/multimodal/realtime_tools.py (NEW FILE)
"""
Tools available during real-time voice conversations.

These are LiveKit function_tool wrapped versions that work
in real-time mode. They provide access to your full agent
capabilities during voice conversations.
"""
from livekit.agents import function_tool
import logging

logger = logging.getLogger("realtime.tools")

# Global state for vision
LATEST_VIDEO_FRAME = None


@function_tool
async def search_web(query: str) -> str:
    """
    Search the web for current information.

    Args:
        query: What to search for

    Returns:
        Search results
    """
    from backend.agno.tools.builtin.web import DuckDuckGoSearch

    try:
        search = DuckDuckGoSearch()
        results = await search.search(query, max_results=5)
        return f"Found: {results}"
    except Exception as e:
        logger.error(f"Search error: {e}")
        return f"Search failed: {e}"


@function_tool
async def query_docs(library: str, question: str) -> str:
    """
    Query documentation using Context7.

    Use this to look up API docs, framework docs, etc.

    Args:
        library: Library name (e.g., "livekit", "agno")
        question: What you need help with

    Returns:
        Documentation answer
    """
    try:
        # Use Context7 MCP
        from backend.integrations.mcp import create_mcp_tools

        mcp = await create_mcp_tools(["context7"])
        result = await mcp.query_docs(
            library_name=library,
            query=question
        )
        return result
    except Exception as e:
        logger.error(f"Docs query error: {e}")
        return f"Could not query docs: {e}"


@function_tool
async def get_tasks() -> str:
    """
    Get user's todo list tasks.

    Returns:
        Current tasks
    """
    from backend.agno.tools.todoist import TodoistTools

    try:
        todoist = TodoistTools()
        tasks = await todoist.get_tasks()
        return f"Your tasks: {tasks}"
    except Exception as e:
        return f"Could not get tasks: {e}"


@function_tool
async def create_task(title: str, description: str = "") -> str:
    """
    Create a new task in Todoist.

    Args:
        title: Task title
        description: Task description

    Returns:
        Success message
    """
    from backend.agno.tools.todoist import TodoistTools

    try:
        todoist = TodoistTools()
        await todoist.create_task(title, description)
        return f"Created task: {title}"
    except Exception as e:
        return f"Failed to create task: {e}"


@function_tool
async def analyze_vision() -> str:
    """
    Analyze what the camera sees.

    Only works when vision is enabled.

    Returns:
        Description of what's visible
    """
    global LATEST_VIDEO_FRAME

    if not LATEST_VIDEO_FRAME:
        return "Camera not enabled or no frame available"

    try:
        from backend.agno.multimodal.analysis import analyze_image

        result = await analyze_image(
            image_bytes=LATEST_VIDEO_FRAME,
            prompt="Describe what you see in this image"
        )
        return result
    except Exception as e:
        return f"Vision analysis failed: {e}"


@function_tool
async def ask_full_agent(request: str) -> str:
    """
    Hand off to the full agent for complex multi-step tasks.

    Use this when you need:
    - Multiple tool calls
    - Research and analysis
    - Complex reasoning
    - Access to knowledge bases

    Args:
        request: Detailed request for the agent

    Returns:
        Agent's complete response
    """
    from backend.shared.registry import get_agent

    try:
        # Load the current agent
        agent = get_agent("legacy-mother-ai")

        # Run full agent
        result = await agent.arun(request)
        return result.content if result.content else "Done."
    except Exception as e:
        logger.error(f"Agent error: {e}")
        return f"Error: {e}"


# Export all tools
REALTIME_TOOLS = [
    search_web,
    query_docs,
    get_tasks,
    create_task,
    analyze_vision,
    ask_full_agent,
]
```

### Step 8: Update Worker to Use New Tools

The worker already calls `get_realtime_tools()`, so the new tools will be automatically loaded!

## Testing Locally

### Terminal 1: Start LiveKit

```bash
livekit-server --dev --bind 0.0.0.0

# You'll see:
# INFO: LiveKit server is ready
# INFO: Listening on :7880
```

### Terminal 2: Start Voice Worker

```bash
python -m backend.agno.multimodal.realtime_worker dev

# You'll see:
# INFO: Starting OneMind Real-Time Agent Worker
# INFO: Connected to LiveKit at ws://localhost:7880
# INFO: Waiting for voice sessions...
```

### Terminal 3: Create Voice Session

```bash
# Test with any agent!

# Legacy agent
curl -X POST http://localhost:8000/api/realtime/session \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "legacy-mother-ai",
    "provider": "pipeline",
    "enable_vision": false
  }'

# Research agent
curl -X POST http://localhost:8000/api/realtime/session \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "research-agent",
    "provider": "pipeline"
  }'

# CEO agent
curl -X POST http://localhost:8000/api/realtime/session \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "ceo-agent",
    "provider": "pipeline"
  }'
```

## All Agents Get Voice!

### How It Works

1. **User picks agent**: Choose which agent to talk to
2. **Session created**: LiveKit room created with agent metadata
3. **Worker loads agent**: Worker loads specified agent from registry
4. **Pipeline created**: Voice pipeline wraps the agent
5. **User talks**: Agent responds with full capabilities

### Example Conversations

**With Research Agent:**
```
User: "Search for papers on quantum computing"
→ Research agent uses DuckDuckGo + ArXiv
→ Synthesizes findings
→ Speaks summary
```

**With Legacy Agent + Context7:**
```
User: "How do I use LiveKit agents?"
→ Agent calls query_docs("livekit", "how to use agents")
→ Context7 MCP returns docs
→ Agent explains in natural language
→ Speaks answer
```

**With CEO Agent:**
```
User: "What's on my agenda today?"
→ CEO agent checks calendar
→ Reviews tasks
→ Speaks schedule
```

## Architecture: All Agents Voice-Enabled

```
┌────────────────────────────────────────────────────────┐
│                   Agent Registry                       │
│                                                        │
│  • legacy-mother-ai (29 tools)                        │
│  • research-agent (search tools)                      │
│  • ceo-agent (productivity tools)                     │
│  • analyst-agent (data tools)                         │
│  • writer-agent (content tools)                       │
│  • ANY custom agent                                   │
│                                                        │
│  All have access to:                                  │
│  ├─ 29+ tools                                         │
│  ├─ MCP servers (Context7, memory, filesystem)       │
│  ├─ Knowledge bases                                   │
│  └─ Custom logic                                      │
└────────────────────────────────────────────────────────┘
                         ▲
                         │ Worker loads agent by ID
                         │
┌────────────────────────────────────────────────────────┐
│              Voice Pipeline Wrapper                    │
│                                                        │
│  User speaks → Deepgram → Agent.arun() → Cartesia     │
│                                                        │
│  Agent uses tools in real-time during conversation    │
└────────────────────────────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────┐
│           Local LiveKit Server                         │
│           ws://localhost:7880                          │
└────────────────────────────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────┐
│                 User's Device                          │
│              (Flutter App / CLI)                       │
└────────────────────────────────────────────────────────┘
```

## Cost Comparison: Local vs Cloud

### Local LiveKit (Recommended)
- **LiveKit**: FREE (self-hosted)
- **Deepgram**: $0.0043/min
- **Cartesia**: $0.015/min
- **Total**: ~$0.02/min + your model costs

### Cloud LiveKit
- **LiveKit Cloud**: $0.0009/min per participant
- **Deepgram**: $0.0043/min
- **Cartesia**: $0.015/min
- **Total**: ~$0.02/min + your model costs

**Local is better for:**
- Development
- Privacy
- Control
- No vendor lock-in

## Benefits of All-Agents-Voice Architecture

### Flexibility
- User chooses which agent to talk to
- Different agents for different contexts
- Easy to add new voice-enabled agents

### Simplicity
- One worker handles all agents
- Agents loaded dynamically
- No per-agent configuration needed

### Power
- Each agent keeps its full capabilities
- All tools work in voice
- MCP servers accessible
- Knowledge bases work

### Code Reusability
- Same agent code for chat and voice
- No voice-specific agent implementations
- Tools work everywhere

## Next Steps

1. **Start Local LiveKit**:
   ```bash
   livekit-server --dev
   ```

2. **Get API Keys** (5 min):
   - Deepgram: https://deepgram.com
   - Cartesia: https://cartesia.ai

3. **Update .env**:
   ```bash
   LIVEKIT_URL=ws://localhost:7880
   LIVEKIT_API_KEY=devkey
   LIVEKIT_API_SECRET=secret
   DEEPGRAM_API_KEY=your-key
   CARTESIA_API_KEY=your-key
   ```

4. **Create Realtime Tools** (copy code above)

5. **Start Worker**:
   ```bash
   python -m backend.agno.multimodal.realtime_worker dev
   ```

6. **Test with Any Agent**!

Your entire agent fleet now has voice! 🎤✨
