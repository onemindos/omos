# OneMind OS Architecture

> A comprehensive guide to understanding how OneMind OS works from backend to frontend.

## Table of Contents
1. [System Overview](#system-overview)
2. [Backend Architecture](#backend-architecture)
3. [Frontend Architecture](#frontend-architecture)
4. [Data Flow](#data-flow)
5. [Key Files Reference](#key-files-reference)
6. [API Endpoints](#api-endpoints)

---

## System Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              OneMind OS                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐     HTTP/SSE      ┌──────────────────────────────────┐   │
│  │   Flutter    │ ◄───────────────► │         FastAPI Backend          │   │
│  │   Frontend   │                   │         (AgentOS/Agno)           │   │
│  └──────────────┘                   └──────────────────────────────────┘   │
│         │                                         │                         │
│         │                                         │                         │
│    Riverpod                              ┌────────┴────────┐                │
│    State Mgmt                            │                 │                │
│                                     PostgreSQL         NATS                 │
│                                    (Memories,      (Events,                 │
│                                     Sessions)      Consciousness)           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Core Technologies
- **Backend**: Python 3.11+, FastAPI, Agno (AgentOS framework)
- **Frontend**: Flutter/Dart, Riverpod state management
- **Database**: PostgreSQL (via asyncpg/psycopg)
- **Events**: NATS JetStream (consciousness/awareness layer)
- **AI Models**: OpenRouter, OpenAI, Anthropic (configurable)

---

## Backend Architecture

### Entry Point: `app.py`

The main FastAPI application that bootstraps everything:

```python
# app.py - What happens on startup

1. Import AgentOS from Agno framework
2. Discover agents from backend/agents/
3. Mount custom routes (Memory, Sessions, Approvals, etc.)
4. Initialize StateManager for awareness system
5. Start listening on port 8080
```

### Directory Structure

```
backend/
├── agents/                    # AI Agent definitions
│   ├── legacy/               # Main assistant agent
│   │   ├── agent.py          # Agent config (model, tools, instructions)
│   │   └── tools.py          # Custom tools (calendar, email, etc.)
│   └── ...
│
├── agno/                      # Agno core components (custom implementations)
│   ├── memory/               # Semantic memory system
│   │   ├── api.py            # REST endpoints (/memories, /api/memory/*)
│   │   ├── db.py             # Database operations
│   │   └── search.py         # Vector similarity search
│   │
│   ├── sessions/             # Conversation session management
│   │   ├── api.py            # REST endpoints (/sessions/*)
│   │   └── service.py        # Session CRUD operations
│   │
│   ├── approvals/            # Human-in-the-loop approvals
│   │   ├── api.py            # REST endpoints (/approvals/*)
│   │   └── store.py          # Approval queue management
│   │
│   └── learning/             # User preferences & learning
│       ├── api.py            # REST endpoints (/api/learning/*)
│       └── profile.py        # User profile management
│
├── consciousness/            # Awareness & context system
│   ├── awareness/
│   │   ├── state.py          # StateManager - tracks user context
│   │   ├── device_presence.py # Device tracking
│   │   ├── patterns.py       # Pattern detection
│   │   └── temporal.py       # Time awareness
│   │
│   ├── core/
│   │   ├── bus.py            # NATS event bus
│   │   ├── kv_store.py       # Key-value state storage
│   │   └── object_store.py   # Object/file storage
│   │
│   └── processing/
│       ├── processor.py      # Event filter/router
│       └── correlator.py     # Event correlation
│
├── lifeos/                   # Life management features
│   ├── inbox/                # Notification inbox
│   ├── productivity/         # Tasks, habits, goals
│   ├── planning/             # Calendar, planner
│   └── lifestyle/            # Meals, recipes
│
└── platform/
    └── db/
        └── __init__.py       # Database connection management
```

### How Agents Work

```python
# backend/agents/legacy/agent.py (simplified)

from agno.agent import Agent
from agno.models import OpenRouter  # or OpenAI, Anthropic
from backend.platform.db import get_shared_db

legacy_agent = Agent(
    agent_id="legacy",
    name="Legacy",
    model=OpenRouter(id="anthropic/claude-sonnet-4"),

    # Shared database = shared memory across all agents
    db=get_shared_db(),
    enable_user_memories=True,

    # Tools the agent can use
    tools=[calendar_tool, email_tool, task_tool, ...],

    # System instructions
    instructions="You are Legacy, a personal AI assistant...",
)
```

### Database Connection Pattern

```python
# backend/platform/db/__init__.py

# Singleton pattern - ONE database for ALL agents
_shared_db = None

def get_shared_db() -> PostgresDb:
    """All agents share this connection = shared memories."""
    global _shared_db
    if _shared_db is None:
        _shared_db = PostgresDb(db_url=get_db_url())
    return _shared_db

# Async pool for services needing direct asyncpg
async def get_db_pool():
    """For StateManager and other async services."""
    ...
```

### API Router Mounting

```python
# app.py - mount_custom_routes()

# Each feature has its own router with endpoints
# Routers are mounted to the main FastAPI app

routers_to_mount = [
    ("", memory_router, "memory"),       # /memories, /api/memory/*
    ("", sessions_router, "sessions"),   # /sessions/*
    ("", approvals_router, "approvals"), # /approvals/*
    ("", inbox_router, "inbox"),         # /api/inbox/*
    ("", learning_router, "learning"),   # /api/learning/*
]

for prefix, router, name in routers_to_mount:
    app.include_router(router, prefix=prefix)
```

---

## Frontend Architecture

### Entry Point: `frontend/lib/main.dart`

```dart
void main() async {
  // Initialize services
  await dotenv.load();

  runApp(
    ProviderScope(  // Riverpod state management
      child: OneMindApp(),
    ),
  );
}
```

### Directory Structure

```
frontend/lib/
├── main.dart                  # App entry point
│
├── platform/                  # Core platform layer
│   ├── api/
│   │   └── agno_client.dart   # Backend API client (4,362 lines!)
│   │
│   ├── config/
│   │   └── environment.dart   # API URL configuration
│   │
│   ├── models/                # Data models (Freezed)
│   │   ├── agent.dart
│   │   ├── session.dart
│   │   ├── memory.dart
│   │   └── ...
│   │
│   └── providers/
│       └── app_providers.dart # Riverpod state providers
│
├── features/                  # Feature modules
│   ├── chat/                  # Main chat interface
│   │   ├── screens/
│   │   ├── widgets/
│   │   └── providers/
│   │
│   ├── productivity/          # Tasks, habits, goals
│   ├── planning/              # Calendar, planner
│   └── ...
│
└── shared/                    # Shared widgets & utilities
```

### Environment Configuration

```dart
// frontend/lib/platform/config/environment.dart

class Environment {
  static String get apiBaseUrl {
    if (kDebugMode) {
      // Development - via Tailscale to Azure VM
      if (Platform.isIOS || Platform.isAndroid) {
        return 'http://100.102.21.44:8081';  // Tailscale IP (Azure VM)
      }
      return 'http://localhost:8080';
    }
    // Production
    return 'https://api.onemindos.dev';
  }
}
```

### API Client Pattern

```dart
// frontend/lib/platform/api/agno_client.dart

class AgnoClient {
  final String baseUrl;
  final http.Client _client;

  // Example: Stream agent responses via SSE
  Stream<AgentResponse> streamAgentRun({
    required String agentId,
    required String message,
    String? sessionId,
  }) async* {
    final request = http.Request(
      'POST',
      Uri.parse('$baseUrl/agents/$agentId/runs'),
    );
    request.headers['Accept'] = 'text/event-stream';
    request.body = jsonEncode({
      'message': message,
      'session_id': sessionId,
      'stream': true,
    });

    final response = await _client.send(request);

    // Parse Server-Sent Events
    await for (final chunk in response.stream.transform(utf8.decoder)) {
      // Parse "data: {...}" lines
      yield AgentResponse.fromJson(jsonDecode(data));
    }
  }

  // Example: Get memories
  Future<List<Memory>> getMemories({String userId = 'zeus'}) async {
    final response = await _client.get(
      Uri.parse('$baseUrl/api/memory?user_id=$userId'),
    );
    return (jsonDecode(response.body)['memories'] as List)
        .map((m) => Memory.fromJson(m))
        .toList();
  }
}
```

### Riverpod State Management

```dart
// frontend/lib/platform/providers/app_providers.dart

// API Client provider
final agnoClientProvider = Provider<AgnoClient>((ref) {
  return AgnoClient(baseUrl: Environment.apiBaseUrl);
});

// Agents list (fetched from backend)
final agentsProvider = FutureProvider<List<Agent>>((ref) async {
  final client = ref.watch(agnoClientProvider);
  return client.getAgents();
});

// Sessions for current user
final sessionsProvider = FutureProvider<List<Session>>((ref) async {
  final client = ref.watch(agnoClientProvider);
  return client.getSessions();
});

// Memories with search capability
final memoriesProvider = FutureProvider.family<List<Memory>, String?>(
  (ref, query) async {
    final client = ref.watch(agnoClientProvider);
    if (query != null) {
      return client.searchMemories(query);
    }
    return client.getMemories();
  },
);
```

### Using Providers in Widgets

```dart
// Example: Chat screen

class ChatScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Watch agent list
    final agentsAsync = ref.watch(agentsProvider);

    return agentsAsync.when(
      data: (agents) => AgentSelector(agents: agents),
      loading: () => CircularProgressIndicator(),
      error: (err, stack) => Text('Error: $err'),
    );
  }

  void sendMessage(WidgetRef ref, String message) {
    final client = ref.read(agnoClientProvider);

    // Stream response from agent
    client.streamAgentRun(
      agentId: 'legacy',
      message: message,
      sessionId: currentSessionId,
    ).listen((response) {
      // Update UI with streaming response
      setState(() {
        currentResponse += response.content;
      });
    });
  }
}
```

---

## Data Flow

### Chat Message Flow

```
┌──────────────────────────────────────────────────────────────────────────┐
│                         CHAT MESSAGE FLOW                                 │
└──────────────────────────────────────────────────────────────────────────┘

User types message in Flutter
        │
        ▼
┌─────────────────┐
│  ChatProvider   │  (Riverpod)
│  sendMessage()  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     POST /agents/legacy/runs
│   AgnoClient    │ ──────────────────────────────►
│ streamAgentRun  │     {message, session_id, stream: true}
└────────┬────────┘
         │                                          ┌─────────────────┐
         │                                          │    FastAPI      │
         │                                          │   POST handler  │
         │                                          └────────┬────────┘
         │                                                   │
         │                                                   ▼
         │                                          ┌─────────────────┐
         │                                          │  Legacy Agent   │
         │                                          │   agent.run()   │
         │                                          └────────┬────────┘
         │                                                   │
         │                                          ┌────────┴────────┐
         │                                          │                 │
         │                                          ▼                 ▼
         │                                    ┌──────────┐     ┌──────────┐
         │                                    │  Tools   │     │  Memory  │
         │                                    │ (email,  │     │  Recall  │
         │                                    │ calendar)│     │          │
         │                                    └──────────┘     └──────────┘
         │                                          │                 │
         │                                          └────────┬────────┘
         │                                                   │
         │                                                   ▼
         │                                          ┌─────────────────┐
         │         SSE: data: {"content": "..."}    │    LLM API      │
         │ ◄─────────────────────────────────────── │  (OpenRouter)   │
         │         (streaming chunks)               └─────────────────┘
         │
         ▼
┌─────────────────┐
│   StreamBuilder │  (Flutter widget)
│   Updates UI    │
└─────────────────┘
```

### Memory Storage Flow

```
┌──────────────────────────────────────────────────────────────────────────┐
│                        MEMORY STORAGE FLOW                                │
└──────────────────────────────────────────────────────────────────────────┘

During conversation, agent extracts facts about user
        │
        ▼
┌─────────────────────────────────────────┐
│  Agent: "User prefers morning meetings" │
│  → Triggers memory storage              │
└────────────────────┬────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────┐
│           PostgresDb.upsert()           │
│  Table: ai.user_memories                │
│  - memory_id: uuid                      │
│  - user_id: "zeus"                      │
│  - memory: "Prefers morning meetings"   │
│  - embedding: [0.1, 0.2, ...]          │  ← Vector for semantic search
└────────────────────┬────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────┐
│        Future conversations...          │
│                                         │
│  User: "Schedule a meeting with John"   │
│                    │                    │
│                    ▼                    │
│  Agent recalls: user prefers mornings   │
│  Agent: "I'll schedule it at 10 AM"     │
└─────────────────────────────────────────┘
```

---

## Key Files Reference

### Backend

| File | Purpose |
|------|---------|
| `app.py` | Main entry point, mounts all routes |
| `backend/platform/db/__init__.py` | Database connection singleton |
| `backend/agents/legacy/agent.py` | Main assistant agent definition |
| `backend/agno/memory/api.py` | Memory REST endpoints |
| `backend/agno/sessions/api.py` | Session REST endpoints |
| `backend/agno/approvals/api.py` | HITL approval endpoints |
| `backend/consciousness/awareness/state.py` | StateManager for context |
| `backend/consciousness/__init__.py` | Consciousness layer exports |

### Frontend

| File | Purpose |
|------|---------|
| `lib/main.dart` | App entry point |
| `lib/platform/api/agno_client.dart` | All backend API calls |
| `lib/platform/config/environment.dart` | API URL configuration |
| `lib/platform/providers/app_providers.dart` | Riverpod state providers |
| `lib/platform/models/*.dart` | Data models (Freezed) |
| `lib/features/chat/` | Main chat interface |

---

## API Endpoints

### Core (Agno Built-in)
```
GET  /health              # Health check
GET  /agents              # List all agents
POST /agents/{id}/runs    # Run agent (chat)
GET  /teams               # List teams
POST /teams/{id}/runs     # Run team
GET  /workflows           # List workflows
```

### Memory (Custom)
```
GET  /api/memory              # List memories
POST /memories                # Create memory
PATCH /memories/{id}          # Update memory
DELETE /memories/{id}         # Delete memory
POST /api/memory/semantic-search  # Search memories
GET  /api/memory/stats        # Memory statistics
```

### Sessions (Custom)
```
GET  /sessions               # List sessions
POST /sessions               # Create session
GET  /sessions/{id}          # Get session
PUT  /sessions/{id}          # Update session
DELETE /sessions/{id}        # Delete session
PUT  /sessions/{id}/rename   # Rename session
GET  /sessions/{id}/messages # Get messages
POST /sessions/{id}/messages # Add message
```

### Approvals (HITL)
```
GET  /approvals              # List pending approvals
POST /approvals/{id}/approve # Approve action
POST /approvals/{id}/reject  # Reject action
GET  /approvals/stats        # Approval statistics
```

### Awareness
```
GET  /api/awareness          # Get current awareness state
```

### Learning
```
GET  /api/learning/profile   # User preferences
GET  /api/learning/entities  # Known entities
GET  /api/learning/culture   # Cultural preferences
GET  /api/learning/stats     # Learning statistics
```

### Productivity
```
GET/POST /api/tasks          # Task management
GET/POST /api/habits         # Habit tracking
GET/POST /api/goals          # Goal management
GET/POST /api/routines       # Daily routines
GET/POST /api/journal        # Journal entries
```

### Planning
```
GET/POST /api/calendar       # Calendar events
GET/POST /api/planner        # Planning items
GET  /api/activity           # Activity feed
```

---

## Running the System

### Backend
```bash
# From project root
cd /Users/zeuslegacy/Documents/onemind-os

# Set environment variables
export DATABASE_URL="postgresql+psycopg://onemind:password@100.102.21.44:5432/onemind"
export OPENROUTER_API_KEY="..."

# Run with uvicorn
python backend/onemindos.py
# or
uvicorn app:app --host 0.0.0.0 --port 8080 --reload
```

### Frontend
```bash
cd frontend

# Get dependencies
flutter pub get

# Run on different platforms
flutter run -d chrome      # Web browser
flutter run -d macos       # macOS desktop
flutter run -d "Zeus"      # iPhone (your device)
```

### Environment Variables

| Variable | Purpose | Example |
|----------|---------|---------|
| `DATABASE_URL` | PostgreSQL connection | `postgresql+psycopg://user:pass@host:5432/db` |
| `OPENROUTER_API_KEY` | LLM API key | `sk-or-...` |
| `NATS_URL` | NATS server | `nats://localhost:4222` |

---

## Debugging Tips

### Backend Logs
```bash
# Watch backend logs
tail -f logs/app.log

# Or run with verbose logging
LOG_LEVEL=DEBUG python backend/onemindos.py
```

### Test Endpoints
```bash
# Health check
curl http://localhost:8080/health

# List agents
curl http://localhost:8080/agents

# Test awareness
curl http://localhost:8080/api/awareness
```

### Flutter DevTools
```bash
# Open Flutter DevTools
flutter pub global activate devtools
flutter pub global run devtools
```

---

## Architecture Decisions

### Why Singleton Database?
All agents share one `PostgresDb` instance so they share memories. When Legacy learns something about you, other agents can recall it too.

### Why Riverpod over Provider/Bloc?
- Compile-time safety
- Better testing support
- Cleaner async handling
- No BuildContext needed

### Why SSE over WebSockets?
- Simpler implementation
- Better HTTP/2 support
- Sufficient for one-way streaming (agent responses)
- WebSockets used only where bidirectional needed

### Why NATS for Events?
- JetStream provides persistence
- Built-in key-value store
- Lightweight compared to Kafka
- Perfect for real-time awareness

---

*Last updated: 2025-01-28*
