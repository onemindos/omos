# Founder Education: OneMind OS Architecture

> Deep dive into how the frontend and backend communicate

---

## The Big Picture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           YOUR DEVICE                                    │
│  ┌─────────────────────┐              ┌─────────────────────────────┐   │
│  │   Flutter App       │    HTTP/SSE  │   Python Backend            │   │
│  │   (agno_client)     │◄────────────►│   (FastAPI + Agno)          │   │
│  │                     │              │                             │   │
│  │   • UI/UX           │              │   • AI Agents               │   │
│  │   • State (Riverpod)│              │   • Memory/Knowledge        │   │
│  │   • Local Storage   │              │   • Tool Execution          │   │
│  └─────────────────────┘              └─────────────────────────────┘   │
│         ▲                                        ▲                      │
│         │                                        │                      │
│         ▼                                        ▼                      │
│  ┌─────────────────────┐              ┌─────────────────────────────┐   │
│  │   SharedPreferences │              │   PostgreSQL / SQLite       │   │
│  │   (settings, cache) │              │   (memories, sessions)      │   │
│  └─────────────────────┘              └─────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
                                   │
                                   ▼
                    ┌─────────────────────────────┐
                    │   External Services         │
                    │   • Anthropic/OpenAI APIs   │
                    │   • MCP Servers             │
                    │   • Home Assistant          │
                    │   • CalDAV                  │
                    └─────────────────────────────┘
```

---

## Part 1: The Frontend (Flutter)

### What is `agno_client.dart`?

It's the **messenger** between your UI and the backend. When you tap a button or send a message, `agno_client` packages that request and sends it to the Python server.

**Location**: `frontend/lib/platform/api/agno_client.dart`

### The Client Architecture

```
AgnoClient (the facade - one interface for everything)
    │
    ├── AgentClient       → Runs AI agents, streams responses
    ├── SessionClient     → Manages chat sessions
    ├── WorkflowClient    → Multi-agent workflows
    ├── ConsciousnessClient → Memory, awareness, presence
    ├── HealthClient      → System health checks
    ├── ModelsClient      → AI model selection (Claude, GPT, etc.)
    ├── MediaClient       → File uploads, media management
    ├── SettingsClient    → User preferences
    ├── IntegrationsClient → OAuth, external services
    └── NotificationsClient → Push notifications
```

### How to Use It (In Your Code)

```dart
// 1. Get the client via Riverpod provider
final client = ref.read(agnoClientProvider);

// 2. Make API calls
final agents = await client.listAgents();           // GET /api/v1/agents
final response = await client.runAgent('zeus');     // POST /api/v1/agents/zeus/run
final memories = await client.getMemories();        // GET /api/v1/memories
```

### The Provider Pattern

```dart
// In app_providers.dart
final agnoClientProvider = Provider<AgnoClient>((ref) {
  final baseUrl = ref.watch(apiBaseUrlProvider);  // Usually http://localhost:8000
  return AgnoClient(baseUrl: baseUrl);
});

// In any widget/screen
class MyScreen extends ConsumerWidget {
  Widget build(BuildContext context, WidgetRef ref) {
    final client = ref.read(agnoClientProvider);
    // Now you can call client.anything()
  }
}
```

---

## Part 2: The Backend (Python)

### What is `main.py` / the Agno Framework?

The Python backend is a **FastAPI** server that:
1. Receives HTTP requests from the Flutter app
2. Runs AI agents using the **Agno** framework
3. Streams responses back using **SSE** (Server-Sent Events)

**Location**: `backend/main.py`

### Key Backend Components

```
backend/
├── main.py                 # FastAPI app entry point
├── agno/
│   ├── agents/             # AI agent definitions
│   │   ├── zeus.py         # Your personal agent
│   │   ├── researcher.py   # Research agent
│   │   └── ...
│   ├── api.py              # API routes (/api/v1/agents, etc.)
│   └── mcp/                # MCP server integrations
├── platform/
│   ├── lifeos/             # Life management features
│   ├── integrations/       # External service connectors
│   └── streams/            # Real-time event streams
```

### How Agents Work

```python
# backend/agno/agents/zeus.py
from agno import Agent, tool

class ZeusAgent(Agent):
    name = "zeus"
    model = "claude-sonnet-4-20250514"
    instructions = "You are Zeus, a personal AI assistant..."

    @tool
    def search_memories(self, query: str) -> list:
        """Search through stored memories"""
        # Tool implementation
        pass

    @tool
    def create_task(self, title: str, due_date: str) -> dict:
        """Create a new task"""
        # Tool implementation
        pass
```

---

## Part 3: The Communication Protocol

### HTTP Requests (Standard API Calls)

```
Flutter                          Python
───────                          ──────
client.listAgents()    ────►     GET /api/v1/agents
                       ◄────     [{"id": "zeus", "name": "Zeus"}, ...]

client.getMemories()   ────►     GET /api/v1/memories
                       ◄────     [{"id": "1", "content": "..."}, ...]
```

### SSE (Server-Sent Events) for Streaming

When you chat with an agent, responses stream back in real-time:

```
Flutter                                  Python
───────                                  ──────
client.streamAgentRun(
  agentId: 'zeus',         ────►        POST /api/v1/agents/zeus/run
  message: 'Hello'
)
                                        ┌─ Agent thinks...
                           ◄─SSE────    │  event: text_delta
                                        │  data: {"content": "Hello"}
                           ◄─SSE────    │
                                        │  event: text_delta
                                        │  data: {"content": "! How"}
                           ◄─SSE────    │
                                        │  event: tool_call_start
                                        │  data: {"tool": "search_memories"}
                           ◄─SSE────    │
                                        │  event: tool_call_end
                                        │  data: {"result": [...]}
                           ◄─SSE────    │
                                        │  event: run_finished
                                        └─ data: {"status": "complete"}
```

### SSE Event Types (AG-UI Protocol)

| Event | When It Fires | What It Contains |
|-------|---------------|------------------|
| `text_delta` | Agent is typing | Partial text content |
| `tool_call_start` | Agent calls a tool | Tool name, arguments |
| `tool_call_end` | Tool finished | Tool result |
| `step_started` | Reasoning step | Step description |
| `run_finished` | Agent done | Final status |
| `run_error` | Something broke | Error message |

---

## Part 4: Data Flow Example

### "Send a message to Zeus"

```
1. USER TAPS SEND
   └─► CommandScreen._sendMessage()
       └─► commandProvider.sendMessage("What's on my calendar?")

2. FLUTTER PREPARES REQUEST
   └─► agnoClient.streamAgentRun(
         agentId: 'zeus',
         sessionId: 'abc123',
         message: "What's on my calendar?"
       )

3. HTTP REQUEST SENT
   └─► POST http://localhost:8000/api/v1/agents/zeus/run
       Body: {"message": "What's on my calendar?", "session_id": "abc123"}

4. PYTHON RECEIVES REQUEST
   └─► FastAPI routes to agno/api.py
       └─► Loads ZeusAgent
       └─► Agent.run(message)

5. AGENT PROCESSES
   └─► Claude API called with:
       - System prompt (agent instructions)
       - User message
       - Available tools

6. AGENT DECIDES TO USE TOOL
   └─► Calls @tool get_calendar_events()
       └─► Fetches from CalDAV
       └─► Returns events

7. AGENT GENERATES RESPONSE
   └─► "You have 3 meetings today: ..."

8. SSE STREAMS BACK
   └─► Flutter receives events
       └─► UI updates in real-time
       └─► Message bubble shows typing
       └─► Tool card shows calendar

9. RUN COMPLETES
   └─► Session saved to database
   └─► Memory extracted (optional)
```

---

## Part 5: The Stub Methods

### Why Stubs Exist

Many `agno_client.dart` methods are **stubs** - they return fake data:

```dart
// This is a stub - backend endpoint doesn't exist yet
Future<Map<String, dynamic>?> getUsageMetrics({int days = 7}) async => null;
```

**Purpose**: Let the Flutter app compile and run while backend features are being built.

### How to Implement a Real Endpoint

**Step 1: Add Python Endpoint**
```python
# backend/agno/api.py
@router.get("/api/v1/metrics/usage")
async def get_usage_metrics(days: int = 7):
    # Real implementation
    return {"usage_by_day": [...], "total_requests": 1234}
```

**Step 2: Update Flutter Client**
```dart
// frontend/lib/platform/api/agno_client.dart
Future<Map<String, dynamic>?> getUsageMetrics({int days = 7}) async {
  final response = await get('/api/v1/metrics/usage', params: {'days': days});
  return response;
}
```

---

## Part 6: Key Files Reference

### Frontend (Flutter)

| File | Purpose |
|------|---------|
| `lib/platform/api/agno_client.dart` | API client facade |
| `lib/platform/providers/app_providers.dart` | Riverpod providers |
| `lib/command/screens/command_screen.dart` | Main chat UI |
| `lib/command/providers/command_provider.dart` | Chat state management |
| `lib/shared/theme/os.dart` | Design system tokens |

### Backend (Python)

| File | Purpose |
|------|---------|
| `backend/main.py` | FastAPI app entry |
| `backend/agno/api.py` | API routes |
| `backend/agno/agents/*.py` | Agent definitions |
| `backend/platform/lifeos/*.py` | LifeOS features |

---

## Part 7: Running the System

### Start Backend
```bash
cd backend
uv run uvicorn main:app --reload --port 8000
```

### Start Frontend
```bash
cd frontend
flutter run
```

### Verify Connection
```bash
# Backend should respond
curl http://localhost:8000/api/v1/health

# Should return: {"status": "healthy", ...}
```

---

## Part 8: Common Patterns

### Adding a New Feature

1. **Define the API contract** (what data goes back and forth)
2. **Implement backend endpoint** in `backend/agno/api.py`
3. **Add client method** in `agno_client.dart`
4. **Create provider** in `app_providers.dart`
5. **Build UI** using the provider

### Example: Adding "Get Tasks"

```python
# 1. Backend (Python)
@router.get("/api/v1/tasks")
async def get_tasks(user_id: str = "zeus"):
    tasks = await db.fetch_tasks(user_id)
    return {"tasks": tasks}
```

```dart
// 2. Client (Dart)
Future<List<Task>> getTasks() async {
  final response = await get('/api/v1/tasks');
  return (response['tasks'] as List)
      .map((t) => Task.fromJson(t))
      .toList();
}

// 3. Provider
final tasksProvider = FutureProvider<List<Task>>((ref) async {
  final client = ref.read(agnoClientProvider);
  return client.getTasks();
});

// 4. UI
class TasksScreen extends ConsumerWidget {
  Widget build(context, ref) {
    final tasksAsync = ref.watch(tasksProvider);
    return tasksAsync.when(
      data: (tasks) => ListView(children: tasks.map(TaskCard.new).toList()),
      loading: () => CircularProgressIndicator(),
      error: (e, _) => Text('Error: $e'),
    );
  }
}
```

---

## Summary

| Layer | Technology | Responsibility |
|-------|------------|----------------|
| **UI** | Flutter + Riverpod | Display, user interaction |
| **API Client** | agno_client.dart | HTTP requests, SSE parsing |
| **Backend** | FastAPI + Agno | Routes, agent execution |
| **AI** | Anthropic/OpenAI | Language model inference |
| **Storage** | PostgreSQL/SQLite | Persistence |

**The flow**: User action → Provider → AgnoClient → HTTP → FastAPI → Agent → AI → Response → SSE → UI Update

---

*Created for Zeus - Founder of OneMind OS*
