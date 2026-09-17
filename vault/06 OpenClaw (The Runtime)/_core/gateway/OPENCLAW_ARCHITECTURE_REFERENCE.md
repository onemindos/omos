# OpenClaw vs OneMind OS - Technical Deep Dive

**Date**: February 7, 2026
**Author**: OneMind Development Team
**Purpose**: Comprehensive technical architecture comparison

---

## Table of Contents

1. [Architecture Philosophies](#architecture-philosophies)
2. [Core Technology Stacks](#core-technology-stacks)
3. [Agent/Tool Systems](#agenttool-systems)
4. [Memory & Persistence](#memory--persistence)
5. [Integration Patterns](#integration-patterns)
6. [Chat Interface Implementations](#chat-interface-implementations)
7. [Extension Mechanisms](#extension-mechanisms)
8. [Security Models](#security-models)
9. [Performance Characteristics](#performance-characteristics)
10. [Deployment Architecture](#deployment-architecture)

---

## 1. Architecture Philosophies

### OpenClaw: Gateway Pattern

**Design Philosophy**: Minimal bridge between AI models and system

```
┌─────────────────────────────────────────────────────────┐
│               OpenClaw Gateway Architecture              │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Chat App (WhatsApp/Telegram)                           │
│              ↓                                           │
│  ┌────────────────────────────────────┐                 │
│  │   Gateway Process (Node.js)         │                 │
│  │   - Message routing                 │                 │
│  │   - Context injection               │                 │
│  │   - Skill loading                   │                 │
│  └────────────────────────────────────┘                 │
│              ↓                                           │
│  ┌────────────────────────────────────┐                 │
│  │   AI Model (Claude/GPT/Local)       │                 │
│  │   - Receives augmented prompt       │                 │
│  │   - Generates tool calls            │                 │
│  └────────────────────────────────────┘                 │
│              ↓                                           │
│  ┌────────────────────────────────────┐                 │
│  │   Skill Executor                    │                 │
│  │   - Python/Shell/JS scripts         │                 │
│  │   - Sandboxed or full access        │                 │
│  └────────────────────────────────────┘                 │
│              ↓                                           │
│  System Resources (Files, Browser, APIs)                │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

**Key Characteristics**:
- **Thin layer**: Gateway is < 5000 lines of code
- **Stateless execution**: Each request is independent
- **Context via files**: Memory stored in Markdown files
- **Skill-centric**: Capabilities added through YAML configs
- **Single agent**: No multi-agent orchestration

---

### OneMind OS: Platform Pattern

**Design Philosophy**: Complete operating system for personal AI

```
┌───────────────────────────────────────────────────────────┐
│            OneMind OS Platform Architecture               │
├───────────────────────────────────────────────────────────┤
│                                                           │
│  ┌─────────────────────────────────────────────────┐     │
│  │  LAYER 1: Interface (Flutter Multi-Platform)     │     │
│  │  - Native apps: iOS, Android, macOS, Windows     │     │
│  │  - Web app: Progressive Web App                  │     │
│  │  - Chat interfaces: SSE streaming, voice input   │     │
│  └─────────────────────────────────────────────────┘     │
│                        ↓ REST/SSE/WS                      │
│  ┌─────────────────────────────────────────────────┐     │
│  │  LAYER 2: Brain (Python FastAPI + Agno)         │     │
│  │  ├─ Agent Factory: Dynamic agent creation        │     │
│  │  ├─ Tool Registry: 45+ tools, auto-discovery     │     │
│  │  ├─ Team Orchestration: COORDINATE/COLLABORATE   │     │
│  │  ├─ Workflow Engine: Multi-step automation       │     │
│  │  └─ HITL System: Human approval gates            │     │
│  └─────────────────────────────────────────────────┘     │
│                        ↓ Event Bus                        │
│  ┌─────────────────────────────────────────────────┐     │
│  │  LAYER 3: Nervous System (NATS JetStream)       │     │
│  │  - Event streaming: 200-400μs latency            │     │
│  │  - MQTT bridge: IoT device integration           │     │
│  │  - WebSocket: Real-time UI updates               │     │
│  │  - Leaf nodes: Edge deployment support           │     │
│  └─────────────────────────────────────────────────┘     │
│                        ↓ Data Layer                       │
│  ┌─────────────────────────────────────────────────┐     │
│  │  LAYER 4: Persistence                            │     │
│  │  ├─ PostgreSQL: Structured data + conversations  │     │
│  │  ├─ pgvector: Semantic memory search             │     │
│  │  ├─ LanceDB: Knowledge base (RAG)                │     │
│  │  ├─ Redis: Hot state, sessions, cache            │     │
│  │  └─ TimescaleDB: Event history, metrics          │     │
│  └─────────────────────────────────────────────────┘     │
│                        ↓ Integrations                     │
│  ┌─────────────────────────────────────────────────┐     │
│  │  LAYER 5: Integrations (15+ Services)           │     │
│  │  - Productivity: Todoist, Calendar, Email        │     │
│  │  - Smart Home: Home Assistant                    │     │
│  │  - Code: GitHub, Linear                          │     │
│  │  - Communication: Discord, Telegram, Slack       │     │
│  │  - Biometrics: Oura Ring, Apple Watch            │     │
│  └─────────────────────────────────────────────────┘     │
│                                                           │
└───────────────────────────────────────────────────────────┘
```

**Key Characteristics**:
- **Thick platform**: 50,000+ lines of code
- **Stateful services**: Sessions, workflows, persistent state
- **Context via databases**: Structured storage with vector search
- **Tool-centric**: Capabilities via Agno tools + custom implementations
- **Multi-agent**: Teams with coordination patterns

---

## 2. Core Technology Stacks

### OpenClaw Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Runtime** | Node.js 20+ | Lightweight JavaScript runtime |
| **Language** | TypeScript | Type-safe JavaScript |
| **AI SDK** | Anthropic SDK, OpenAI SDK | Direct API integration |
| **Skills** | YAML + Python/Shell/JS | Extensible functionality |
| **Memory** | Markdown files | Local file-based storage |
| **Browser** | Playwright | Web automation |
| **Chat** | Various SDKs | WhatsApp, Telegram, Discord, Slack |
| **MCP** | Model Context Protocol | Extensibility standard |
| **Dependencies** | ~50 npm packages | Minimal footprint |

**File Structure**:
```
openclaw/
├── src/
│   ├── gateway.ts           # Core gateway logic
│   ├── skills/              # Skill loader
│   ├── memory/              # Markdown memory
│   ├── chat/                # Chat app integrations
│   └── mcp/                 # MCP client
├── skills/                  # Skill definitions (YAML)
├── memory/                  # User memory (Markdown)
└── package.json             # Dependencies
```

**Startup Time**: < 5 seconds

---

### OneMind OS Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Backend** | Python 3.12+ | Modern async Python |
| **Framework** | FastAPI | High-performance API |
| **AI Framework** | Agno 0.0.95+ | Agent orchestration |
| **Frontend** | Flutter 3.6+ | Cross-platform UI |
| **State** | Riverpod | Reactive state management |
| **Database** | PostgreSQL 15+ | Relational data |
| **Vector DB** | pgvector | Semantic search |
| **Knowledge** | LanceDB | RAG implementation |
| **Cache** | Redis 7+ | In-memory cache |
| **Events** | NATS JetStream | Event streaming |
| **Time Series** | TimescaleDB | Analytics |
| **Voice** | LiveKit, Deepgram, Cartesia | Voice interface |
| **Deployment** | Docker Compose | Container orchestration |

**File Structure**:
```
OneMind-OS/
├── backend/
│   ├── main.py                   # FastAPI app (130 lines)
│   ├── agno/
│   │   ├── agents/               # Agent factory & presets
│   │   ├── teams/                # Team factory & presets
│   │   ├── tools/                # Tool registry & custom tools
│   │   ├── models/               # Model factory (60+ models)
│   │   ├── knowledge/            # RAG system
│   │   ├── mcp/                  # MCP servers (20+)
│   │   └── skills.disabled/      # Skills (currently disabled)
│   ├── platform/                 # Infrastructure
│   ├── lifeos/                   # 14 Life OS features
│   └── integrations/             # External systems
├── frontend/
│   └── lib/
│       ├── agno/                 # Agno UI (chat, agents, teams)
│       ├── lifeos/               # Life OS UIs
│       ├── platform/             # Platform services
│       └── shared/               # Reusable widgets
├── infra/
│   ├── deployment/docker/        # Docker Compose configs
│   └── containers/               # Dockerfiles
└── docs/                         # 55+ documentation files
```

**Startup Time**: 10-30 seconds (Docker), 5-10 seconds (local)

---

## 3. Agent/Tool Systems

### OpenClaw: Skills System

**Skill Definition** (YAML):
```yaml
# skills/flight-checkin.yaml
name: Flight Check-in
description: Automatically check into flights
version: 1.0.0
author: community
category: travel

tools:
  - browser
  - calendar

script: |
  #!/usr/bin/env python3
  import sys
  from playwright.sync_api import sync_playwright

  def check_in(airline, confirmation):
      with sync_playwright() as p:
          browser = p.chromium.launch()
          page = browser.new_page()
          page.goto(f"https://{airline}.com/checkin")
          page.fill('#confirmation', confirmation)
          page.click('button[type="submit"]')
          # ... rest of automation
          browser.close()

  if __name__ == '__main__':
      airline, confirmation = sys.argv[1], sys.argv[2]
      check_in(airline, confirmation)
```

**Skill Execution Flow**:
```
1. User: "Check me in for my Delta flight DL123"
2. Gateway loads memory → finds confirmation code
3. AI decides: Use flight-checkin skill
4. Gateway executes: python skills/flight-checkin/script.py delta ABC123
5. Script returns: ✓ Checked in, seat 12A
6. AI responds: "You're checked in for DL123, seat 12A"
```

**Key Features**:
- **Self-modifying**: AI can create new skills by generating YAML + scripts
- **Community-driven**: 100+ skills on ClawHub registry
- **Language-agnostic**: Python, Shell, Node.js, etc.
- **Hot-reloadable**: Skills can be added without restart

---

### OneMind OS: Tool System

**Tool Definition** (Python):
```python
# backend/agno/tools/flight_checkin.py
from agno.tools import Toolkit
from playwright.async_api import async_playwright

class FlightCheckinTool(Toolkit):
    def __init__(self):
        super().__init__(name="flight_checkin")
        self.register(self.check_in)

    async def check_in(self, airline: str, confirmation: str) -> str:
        """Automatically check into a flight.

        Args:
            airline: Airline code (e.g., 'delta', 'united')
            confirmation: Confirmation number

        Returns:
            Check-in status and seat assignment
        """
        async with async_playwright() as p:
            browser = await p.chromium.launch()
            page = await browser.new_page()
            await page.goto(f"https://{airline}.com/checkin")
            await page.fill('#confirmation', confirmation)
            await page.click('button[type="submit"]')
            # ... rest of automation
            await browser.close()

            return "✓ Checked in for DL123, seat 12A"
```

**Tool Registration**:
```python
# backend/agno/tools/registry.py
CUSTOM_TOOL_SPECS = [
    {
        "name": "flight_checkin",
        "module": "backend.agno.tools.flight_checkin",
        "class": "FlightCheckinTool",
        "description": "Automatically check into flights",
        "category": "TRAVEL",
        "requires_approval": True,  # HITL gate
    }
]
```

**Tool Execution Flow**:
```
1. User: "Check me in for my Delta flight DL123"
2. Agent queries memory → finds confirmation code
3. Agent decides: Use flight_checkin tool
4. Agno framework: Tool requires approval → HITL card
5. User approves in Flutter app
6. Tool executes: await flight_checkin.check_in('delta', 'ABC123')
7. Agent responds: "You're checked in for DL123, seat 12A"
```

**Key Features**:
- **Type-safe**: Python type hints, Pydantic validation
- **Integrated**: Tools see full Agno context (memory, sessions)
- **HITL built-in**: `requires_approval` parameter
- **Discoverable**: Frontend auto-generates UI from tool schemas
- **Composable**: Tools can call other tools

---

### Tool/Skill Comparison

| Aspect | OpenClaw Skills | OneMind OS Tools |
|--------|----------------|------------------|
| **Definition** | YAML + Scripts | Python classes |
| **Type Safety** | ❌ No | ✅ Yes (Pydantic) |
| **Hot Reload** | ✅ Yes | ⚠️ Requires restart |
| **Self-Modifying** | ✅ AI writes YAML | ❌ Not implemented |
| **HITL** | ⚠️ Manual | ✅ Built-in |
| **UI Generation** | ❌ No | ✅ Auto-generated |
| **Marketplace** | ✅ ClawHub | ❌ No |
| **Language** | Any | Python only |
| **Context Access** | Limited | Full (memory, sessions) |

---

## 4. Memory & Persistence

### OpenClaw: File-Based Memory

**Memory Structure**:
```
memory/
├── user-profile.md          # User facts & preferences
├── conversations/
│   ├── 2026-02-07.md        # Daily conversation log
│   └── 2026-02-06.md
├── tasks.md                 # Active tasks
├── projects/
│   ├── project-alpha.md     # Project context
│   └── project-beta.md
└── knowledge/
    ├── work-policies.md     # Personal knowledge base
    └── contacts.md
```

**Example Memory File**:
```markdown
# User Profile

## Basic Info
- Name: Zeus
- Timezone: EST
- Preferred communication: Direct, concise

## Preferences
- Model: Claude Opus
- Response style: Professional but friendly
- Work hours: 9 AM - 6 PM EST

## Projects
- [[project-alpha]] - OneMind OS development
- [[project-beta]] - Book writing

## Recent Context (2026-02-07)
- Researching OpenClaw integration
- Working on agent architecture comparison
```

**Memory Retrieval**:
```javascript
// Gateway loads all memory files before each request
const userProfile = fs.readFileSync('memory/user-profile.md', 'utf8');
const recentConvo = fs.readFileSync('memory/conversations/2026-02-07.md', 'utf8');
const projects = glob.sync('memory/projects/*.md').map(f => fs.readFileSync(f, 'utf8'));

// Inject into prompt
const systemPrompt = `
You are an AI assistant. Here's what you know about the user:

${userProfile}

Recent conversation:
${recentConvo}

Active projects:
${projects.join('\n\n')}
`;
```

**Advantages**:
- ✅ Human-readable
- ✅ Git-trackable
- ✅ Easy to edit manually
- ✅ No database needed
- ✅ Portable (just copy files)

**Disadvantages**:
- ❌ No semantic search
- ❌ Scales poorly (> 10MB memory)
- ❌ No structured queries
- ❌ Full context loaded every time (token-heavy)

---

### OneMind OS: Database-Backed Memory

**Memory Structure** (PostgreSQL):
```sql
-- Memory storage
CREATE TABLE agent_memories (
    id UUID PRIMARY KEY,
    agent_id VARCHAR,
    memory_type VARCHAR,  -- 'episodic', 'semantic', 'procedural'
    content TEXT,
    metadata JSONB,
    embedding VECTOR(1536),  -- pgvector for semantic search
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Sessions
CREATE TABLE agent_sessions (
    id UUID PRIMARY KEY,
    agent_id VARCHAR,
    user_id VARCHAR,
    messages JSONB,
    summary TEXT,
    created_at TIMESTAMP
);

-- Knowledge base (LanceDB)
knowledge/
├── embeddings/              # Vector index
├── documents/               # Original documents
└── metadata.json           # Document metadata
```

**Memory Retrieval** (Python):
```python
from agno.storage import AgentMemory
from pgvector.asyncpg import register_vector

# Semantic search
async def get_relevant_memories(query: str, limit: int = 5):
    # Generate embedding for query
    query_embedding = await embed_text(query)

    # Vector similarity search
    results = await db.fetch("""
        SELECT content, metadata,
               1 - (embedding <=> $1) AS similarity
        FROM agent_memories
        WHERE agent_id = $2
        ORDER BY embedding <=> $1
        LIMIT $3
    """, query_embedding, agent_id, limit)

    return [r['content'] for r in results]

# Usage in agent
relevant_context = await get_relevant_memories("flight check-in preferences")
# Returns: ["User prefers aisle seats", "Always flies Delta when possible", ...]
```

**Knowledge Base (RAG)**:
```python
from agno.knowledge import KnowledgeBase
from lancedb import connect

kb = KnowledgeBase(
    path="backend/knowledge/data",
    embedder="text-embedding-3-large"
)

# Add document
await kb.load_document(
    path="docs/onemind-architecture.pdf",
    metadata={"type": "technical", "project": "onemind"}
)

# Search
results = await kb.search(
    query="How does the agent system work?",
    limit=5,
    filters={"type": "technical"}
)

# Returns: Chunks with source citations
# [
#   {
#     "content": "The agent system uses Agno framework...",
#     "source": "onemind-architecture.pdf",
#     "page": 12,
#     "similarity": 0.89
#   },
#   ...
# ]
```

**Advantages**:
- ✅ Semantic search (vector similarity)
- ✅ Scales to GB of data
- ✅ Structured queries (filter by date, type, etc.)
- ✅ Only relevant context loaded (token-efficient)
- ✅ RAG with source citations

**Disadvantages**:
- ❌ Not human-readable (binary storage)
- ❌ Requires database infrastructure
- ❌ Harder to manually edit
- ❌ More complex backup/restore

---

### Memory Comparison

| Feature | OpenClaw | OneMind OS |
|---------|----------|------------|
| **Storage** | Markdown files | PostgreSQL + pgvector |
| **Semantic Search** | ❌ No | ✅ Yes |
| **Scale** | < 10MB | Unlimited |
| **Human-Readable** | ✅ Yes | ❌ No |
| **Git-Trackable** | ✅ Yes | ❌ No |
| **RAG Support** | ⚠️ Via MCP | ✅ Built-in (LanceDB) |
| **Token Efficiency** | Low (loads all) | High (loads relevant) |
| **Manual Editing** | Easy | Requires SQL |

---

## 5. Integration Patterns

### OpenClaw: Skill-Based Integrations

**Pattern**: Each integration is a skill with API credentials

```yaml
# skills/todoist.yaml
name: Todoist
description: Task management integration
api_key: ${TODOIST_API_KEY}  # From env var

tools:
  - http

script: |
  #!/usr/bin/env python3
  import requests
  import os

  API_KEY = os.getenv('TODOIST_API_KEY')
  BASE_URL = 'https://api.todoist.com/rest/v2'

  def get_tasks():
      r = requests.get(f'{BASE_URL}/tasks', headers={'Authorization': f'Bearer {API_KEY}'})
      return r.json()

  def create_task(content, due_string=None):
      r = requests.post(f'{BASE_URL}/tasks',
                       headers={'Authorization': f'Bearer {API_KEY}'},
                       json={'content': content, 'due_string': due_string})
      return r.json()
```

**Advantages**:
- Simple, self-contained
- Easy to share (just copy YAML file)
- No framework lock-in

**Disadvantages**:
- No type safety
- Credentials in env vars only
- No retry logic, error handling
- Each skill reimplements HTTP client

---

### OneMind OS: Framework-Integrated

**Pattern**: Integrations as Agno Toolkits with full framework support

```python
# backend/agno/tools/todoist_tool.py
from agno.tools import Toolkit
from agno.integrations.todoist import TodoistIntegration
import requests

class TodoistTool(Toolkit):
    def __init__(self):
        super().__init__(name="todoist")

        # Get credentials from database (encrypted)
        self.integration = TodoistIntegration.get_active()

        # Register tools
        self.register(self.get_tasks)
        self.register(self.create_task)

    async def get_tasks(self, filter: str = None) -> str:
        """Get tasks from Todoist.

        Args:
            filter: Optional filter (e.g., 'today', 'p1')
        """
        try:
            tasks = await self.integration.api.get_tasks(filter=filter)
            return self._format_tasks(tasks)
        except Exception as e:
            logger.error(f"Todoist API error: {e}")
            return f"❌ Error fetching tasks: {str(e)}"

    async def create_task(
        self,
        content: str,
        due_string: str = None,
        priority: int = 1,
        project_id: str = None
    ) -> str:
        """Create a new task in Todoist.

        Args:
            content: Task description
            due_string: Natural language due date (e.g., 'tomorrow at 3pm')
            priority: Priority 1-4 (4 = highest)
            project_id: Optional project ID
        """
        task = await self.integration.api.add_task(
            content=content,
            due_string=due_string,
            priority=priority,
            project_id=project_id
        )
        return f"✓ Created task: {task.content}"
```

**Integration Registry**:
```python
# backend/integrations/manager.py
class IntegrationManager:
    AVAILABLE_INTEGRATIONS = {
        "todoist": {
            "id": "todoist",
            "name": "Todoist",
            "type": "oauth",  # or "api_key", "direct"
            "auth_url": "https://todoist.com/oauth/authorize",
            "tool_class": "backend.agno.tools.todoist_tool.TodoistTool",
            "scopes": ["data:read_write", "task:add"],
            "config_fields": ["api_token"],
        }
    }

    async def connect_integration(self, integration_id: str, user_id: str):
        """Initiate OAuth flow"""
        config = self.AVAILABLE_INTEGRATIONS[integration_id]
        # ... OAuth flow

    async def store_credentials(self, integration_id: str, user_id: str, credentials: dict):
        """Encrypt and store credentials in database"""
        encrypted = encrypt_credentials(credentials)
        await db.execute("""
            INSERT INTO integration_credentials (integration_id, user_id, credentials)
            VALUES ($1, $2, $3)
        """, integration_id, user_id, encrypted)
```

**Advantages**:
- Type-safe with Pydantic models
- Centralized credential management (encrypted in DB)
- Built-in error handling and retries
- OAuth flow support
- Frontend auto-generates connection UI

**Disadvantages**:
- More boilerplate code
- Requires database for credentials
- Framework-specific (Agno)

---

## 6. Chat Interface Implementations

### OpenClaw: Multi-Platform Chat Adapters

**Architecture**: Each chat platform has an adapter

```typescript
// src/chat/telegram-adapter.ts
import TelegramBot from 'node-telegram-bot-api';
import { Gateway } from '../gateway';

export class TelegramAdapter {
  private bot: TelegramBot;
  private gateway: Gateway;

  constructor(token: string, gateway: Gateway) {
    this.bot = new TelegramBot(token, { polling: true });
    this.gateway = gateway;

    this.bot.on('message', async (msg) => {
      const chatId = msg.chat.id;
      const userId = msg.from.id;
      const text = msg.text;

      // Route to gateway
      const response = await this.gateway.handleMessage({
        userId: userId.toString(),
        sessionId: `telegram_${chatId}`,
        message: text,
        context: {
          platform: 'telegram',
          chatId: chatId,
        }
      });

      // Send response
      await this.bot.sendMessage(chatId, response);
    });
  }
}

// src/gateway.ts
export class Gateway {
  async handleMessage(request: MessageRequest): Promise<string> {
    // Load memory
    const memory = await this.loadMemory(request.userId);

    // Build prompt
    const systemPrompt = this.buildSystemPrompt(memory);

    // Call AI
    const response = await this.callAI(systemPrompt, request.message);

    // Execute any tool calls
    if (response.toolCalls) {
      const results = await this.executeSkills(response.toolCalls);
      // ... final response
    }

    return response.content;
  }
}
```

**Supported Platforms**:
- Telegram (bot API)
- Discord (Discord.js bot)
- WhatsApp (Twilio API)
- Slack (Slack SDK)
- iMessage (via macOS scripting)
- Signal (via signal-cli)

---

### OneMind OS: Flutter App + Optional Chat

**Primary Interface**: Native Flutter app

```dart
// frontend/lib/agno/chat/screens/enhanced_chat_screen.dart
class EnhancedChatScreen extends ConsumerStatefulWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final messages = ref.watch(chatMessagesProvider);
    final isStreaming = ref.watch(streamingProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text('Chat with Legacy'),
        actions: [
          // Voice input button
          IconButton(
            icon: Icon(Icons.mic),
            onPressed: () => _startVoiceInput(),
          ),
          // Agent selector
          AgentSelectorButton(),
        ],
      ),
      body: Column(
        children: [
          // Message list with SSE streaming
          Expanded(
            child: MessageList(
              messages: messages,
              onToolCard: (toolCall) => ToolCardWidget(toolCall),
              onHITLCard: (approval) => ApprovalCardWidget(approval),
            ),
          ),
          // Input field
          MessageInputField(
            onSend: (text) => _sendMessage(text),
            onVoice: () => _startVoiceInput(),
          ),
        ],
      ),
    );
  }

  Future<void> _sendMessage(String text) async {
    final controller = ref.read(chatControllerProvider);

    // Detect @mentions for multi-agent
    final agents = _extractMentions(text);

    if (agents.isNotEmpty) {
      // Stream ad-hoc team run
      await controller.streamAdHocTeamRun(
        message: text,
        agents: agents,
      );
    } else {
      // Stream single agent run
      await controller.streamRun(
        message: text,
        agentId: 'legacy-mother-ai',
      );
    }
  }
}

// Chat controller with SSE streaming
class ChatController {
  Future<void> streamRun({
    required String message,
    required String agentId,
  }) async {
    final url = '$API_BASE/agents/$agentId/runs/stream';

    // Create EventSource for SSE
    final source = EventSource(Uri.parse(url));

    source.onMessage.listen((event) {
      // Parse AG-UI protocol events
      final data = jsonDecode(event.data);

      switch (data['type']) {
        case 'textMessageContent':
          _appendToMessage(data['content']);
          break;
        case 'toolCall':
          _displayToolCard(data['tool']);
          break;
        case 'approvalRequest':
          _displayHITLCard(data['approval']);
          break;
        case 'reasoningContent':
          _displayReasoningStep(data['content']);
          break;
      }
    });
  }
}
```

**Optional Chat Integrations** (Can be added):
```python
# backend/agno/tools/telegram_bot.py (currently notification-only, needs bidirectional)
class TelegramBotTool(Toolkit):
    def __init__(self):
        super().__init__(name="telegram_bot")
        self.bot = Bot(token=TELEGRAM_BOT_TOKEN)

    async def start_listener(self):
        """Start bot polling loop"""
        async with Bot(token=TELEGRAM_BOT_TOKEN) as bot:
            async for update in bot.listen():
                if update.message:
                    await self.handle_message(update.message)

    async def handle_message(self, message):
        """Route incoming message to agent"""
        session_id = f"telegram_{message.from_user.id}"
        agent = get_agent("legacy-mother-ai")

        response = await agent.arun(
            message=message.text,
            session_id=session_id,
        )

        await self.bot.send_message(
            chat_id=message.chat.id,
            text=response.content
        )
```

---

### Chat Interface Comparison

| Aspect | OpenClaw | OneMind OS |
|--------|----------|------------|
| **Primary Interface** | Chat apps | Flutter app |
| **Chat Platforms** | 6+ (native support) | 0 (notification-only) |
| **Native Apps** | ❌ No | ✅ iOS, Android, macOS, Windows |
| **Web App** | ❌ No | ✅ Yes (PWA) |
| **Voice Input** | ⚠️ Via platform | ✅ Native (LiveKit) |
| **Rich UI** | ❌ Text only | ✅ Tool cards, HITL, reasoning |
| **Multi-Agent** | ❌ Single agent | ✅ @mention routing |
| **Setup Complexity** | Low (just API tokens) | High (full app stack) |

---

## 7. Extension Mechanisms

### OpenClaw: Self-Modifying Skills

**Key Feature**: AI can write its own skills

**Example Flow**:
```
User: "I need a skill to check my local library for new books"

AI: I'll create a skill for you. Let me write the configuration and script.

[AI generates YAML config]
---
name: Library Checker
description: Check local library catalog for new books
tools:
  - http
  - notification
---

[AI generates Python script]
#!/usr/bin/env python3
import requests
from bs4 import BeautifulSoup

def check_library():
    url = "https://mylibrary.org/new-releases"
    r = requests.get(url)
    soup = BeautifulSoup(r.text, 'html.parser')

    books = []
    for item in soup.find_all('div', class_='book-item'):
        title = item.find('h3').text
        author = item.find('span', class_='author').text
        books.append(f"{title} by {author}")

    return "\n".join(books)

if __name__ == '__main__':
    print(check_library())

AI: ✓ Created skill "Library Checker". You can now ask me to check the library anytime!

User: Great! Check it now.

AI: [executes skill] Here are the new releases:
- The Midnight Library by Matt Haig
- Project Hail Mary by Andy Weir
- ...
```

**Implementation**:
```javascript
// Gateway can write files
async function createSkill(skillName: string, config: string, script: string) {
  const skillDir = `skills/${skillName}`;
  await fs.mkdir(skillDir, { recursive: true });

  await fs.writeFile(`${skillDir}/skill.yaml`, config);
  await fs.writeFile(`${skillDir}/script.py`, script);
  await fs.chmod(`${skillDir}/script.py`, 0o755);

  // Hot-reload skills
  await reloadSkills();
}
```

---

### OneMind OS: Configuration-Based Extensions

**Current Approach**: Define tools/agents via Python code

**To Add Self-Modifying** (Implementation needed):
```python
# backend/agno/tools/skill_writer.py
from agno.tools import tool
from pathlib import Path

@tool(requires_confirmation=True)  # HITL approval required
async def create_skill(
    name: str,
    description: str,
    instructions: str,
    code: str
) -> str:
    """Create a new skill from AI-generated code.

    REQUIRES HUMAN APPROVAL before execution.

    Args:
        name: Skill name (kebab-case)
        description: Brief description
        instructions: Detailed skill instructions (markdown)
        code: Python implementation
    """
    skill_dir = Path(f"backend/agno/skills/{name}")
    skill_dir.mkdir(parents=True, exist_ok=True)

    # Write SKILL.md with frontmatter
    skill_md = f"""---
name: {name}
description: {description}
version: 1.0.0
category: GENERATED
tags: [ai-generated]
---

# {name}

{instructions}
"""
    (skill_dir / "SKILL.md").write_text(skill_md)

    # Write helper script
    (skill_dir / "scripts" / "main.py").write_text(code)

    # Reload skills
    from backend.agno.skills import discover_skills
    discover_skills()

    return f"✓ Skill '{name}' created and loaded"
```

**With HITL Approval**:
```python
# Agent proposes skill creation
response = await agent.arun(
    message="Create a skill to check my library for new books"
)

# Agno sees requires_confirmation=True → pauses execution
# Frontend displays approval card:
"""
┌───────────────────────────────────────────────┐
│  Approval Required                            │
├───────────────────────────────────────────────┤
│  Agent wants to create a new skill:           │
│                                               │
│  Name: library-checker                        │
│  Description: Check local library catalog     │
│                                               │
│  [View Code]                                  │
│                                               │
│  [Approve]  [Deny]  [Modify]                  │
└───────────────────────────────────────────────┘
"""

# User approves → skill created and loaded
```

---

## 8. Security Models

### OpenClaw Security

**Default Mode**: Full system access (opt-in sandbox)

```javascript
// Default: no restrictions
const result = await executeSkill('shell', ['rm', '-rf', '/important-data']);

// Opt-in sandbox mode
const result = await executeSkill('shell', ['rm', '-rf', '/important-data'], {
  sandbox: true,  // Restricts to skill directory only
});
```

**Security Concerns** (from research):
- 15% of community skills found to contain malicious instructions
- Skills can load external code to bypass security checks
- Self-modifying code increases attack surface
- File-based memory can be manipulated

**Mitigation**:
- Use sandbox mode
- Vet all skills from ClawHub
- Limit API key permissions
- Monitor system activity

---

### OneMind OS Security

**Default Mode**: HITL approval for sensitive operations

```python
# Tools can require approval
@tool(requires_confirmation=True)
async def delete_file(path: str) -> str:
    """Delete a file. REQUIRES APPROVAL."""
    os.remove(path)
    return f"✓ Deleted {path}"

# Agent proposes deletion → HITL card → user approves → executes
```

**Security Layers**:
1. **HITL Approval**: Sensitive operations pause for human review
2. **Encrypted Credentials**: Integration tokens encrypted in PostgreSQL
3. **JWT Auth**: API access requires authentication
4. **RBAC**: Role-based access control (planned)
5. **Audit Logs**: All agent actions logged to TimescaleDB
6. **Rate Limiting**: Prevent abuse

**Security Advantages**:
- Structured approval flow
- Centralized credential management
- Audit trail of all actions

---

## 9. Performance Characteristics

### OpenClaw Performance

**Startup**: < 5 seconds
**Memory Usage**: 50-100MB (Node.js process)
**Response Time**:
- Simple query: 1-3 seconds
- With skill execution: 3-10 seconds
- Browser automation: 10-30 seconds

**Scalability**:
- ✅ Single user: Excellent
- ⚠️ Multi-user: Requires multiple instances
- ❌ Enterprise: Not designed for scale

**Resource Usage**:
- Runs on Raspberry Pi (minimal footprint)
- No database overhead
- File I/O for memory (fast for small datasets)

---

### OneMind OS Performance

**Startup**:
- Docker: 10-30 seconds (cold start)
- Local: 5-10 seconds (hot start)

**Memory Usage**: 500MB-2GB (full stack)
- PostgreSQL: 100-200MB
- Redis: 50-100MB
- Python/FastAPI: 200-500MB
- Flutter app: 100-300MB

**Response Time**:
- Simple query: 1-3 seconds
- With tool execution: 3-10 seconds
- Multi-agent team: 10-30 seconds
- RAG search: 2-5 seconds

**Scalability**:
- ✅ Single user: Excellent
- ✅ Family (5-10 users): Good
- ⚠️ Small team (10-50): Requires optimization
- ❌ Enterprise (100+): Requires re-architecture

**Resource Usage**:
- Requires 4+ CPU cores, 8GB+ RAM
- Database overhead (PostgreSQL, Redis, TimescaleDB)
- Vector search (pgvector) is CPU-intensive

---

### Performance Comparison

| Metric | OpenClaw | OneMind OS |
|--------|----------|------------|
| **Startup Time** | < 5s | 10-30s |
| **Memory Usage** | 50-100MB | 500MB-2GB |
| **Response Time** | 1-10s | 1-30s |
| **Min Hardware** | Raspberry Pi | 4 cores, 8GB RAM |
| **Scalability** | Single user | 1-10 users |

---

## 10. Deployment Architecture

### OpenClaw Deployment

**Single-Machine Setup**:
```bash
# One-liner install
curl -fsSL https://openclaw.ai/install.sh | bash

# What it does:
# 1. Installs Node.js 20+
# 2. Clones OpenClaw repo
# 3. Installs dependencies (npm install)
# 4. Creates config.yaml
# 5. Starts gateway process
# 6. Sets up systemd service (auto-restart)
```

**Architecture**:
```
┌─────────────────────────────────┐
│  User's Machine                 │
├─────────────────────────────────┤
│                                 │
│  ┌───────────────────────────┐  │
│  │  OpenClaw Gateway (Node)  │  │
│  │  - Port: 3000             │  │
│  │  - Process: pm2/systemd   │  │
│  └───────────────────────────┘  │
│              ↕                  │
│  ┌───────────────────────────┐  │
│  │  Filesystem               │  │
│  │  - skills/                │  │
│  │  - memory/                │  │
│  │  - config.yaml            │  │
│  └───────────────────────────┘  │
│              ↕                  │
│  ┌───────────────────────────┐  │
│  │  External APIs            │  │
│  │  - Claude API             │  │
│  │  - Telegram API           │  │
│  │  - WhatsApp (Twilio)      │  │
│  └───────────────────────────┘  │
│                                 │
└─────────────────────────────────┘
```

**Deployment Options**:
- Laptop (easiest)
- Mac Mini (24/7 uptime)
- Raspberry Pi (low power)
- VPS (remote access)

---

### OneMind OS Deployment

**Docker Compose Setup**:
```bash
# Clone repo
git clone https://github.com/Zeus-Delacruz/onemind-os.git
cd onemind-os

# Configure environment
cp .env.example .env
nano .env  # Add AWS credentials, API keys

# Deploy full stack
cd infra/deployment/docker
docker compose -f docker-compose.core.yml up -d

# Services started:
# - PostgreSQL (port 5432)
# - Redis (port 6379)
# - NATS (port 4222)
# - TimescaleDB (port 5433)
# - Agno Backend (port 8081)
# - Flutter Web (port 3000)
```

**Architecture**:
```
┌───────────────────────────────────────────────┐
│  OneMind OS Infrastructure                    │
├───────────────────────────────────────────────┤
│                                               │
│  ┌─────────────────────────────────────────┐  │
│  │  Frontend (Docker Container)            │  │
│  │  - Flutter Web (Nginx)                  │  │
│  │  - Port: 3000 → Public                  │  │
│  └─────────────────────────────────────────┘  │
│                    ↕ REST/SSE/WS              │
│  ┌─────────────────────────────────────────┐  │
│  │  Backend (Docker Container)             │  │
│  │  - FastAPI (Python)                     │  │
│  │  - Port: 8081 → Internal                │  │
│  │  - Workers: 4 (auto-scaling)            │  │
│  └─────────────────────────────────────────┘  │
│         ↕             ↕           ↕            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │PostgreSQL│  │  Redis   │  │   NATS   │    │
│  │  +pgvec  │  │  Cache   │  │  Events  │    │
│  │Port: 5432│  │Port: 6379│  │Port: 4222│    │
│  └──────────┘  └──────────┘  └──────────┘    │
│         ↕                                      │
│  ┌─────────────────────────────────────────┐  │
│  │  TimescaleDB (Docker Container)         │  │
│  │  - Event history & analytics            │  │
│  │  - Port: 5433 → Internal                │  │
│  └─────────────────────────────────────────┘  │
│                    ↕                          │
│  ┌─────────────────────────────────────────┐  │
│  │  External Services                      │  │
│  │  - AWS Bedrock (Claude API)             │  │
│  │  - Todoist, GitHub, etc.                │  │
│  │  - Home Assistant (MQTT)                │  │
│  └─────────────────────────────────────────┘  │
│                                               │
└───────────────────────────────────────────────┘
```

**Deployment Options**:
- Local (Docker Desktop)
- Azure VM (current production)
- AWS EC2 / GCP Compute
- Kubernetes (planned)
- Edge devices (Jetson - planned)

---

### Deployment Comparison

| Aspect | OpenClaw | OneMind OS |
|--------|----------|------------|
| **Install Time** | 5 minutes | 1-3 hours |
| **Components** | 1 (gateway) | 6+ (containers) |
| **Port Management** | Simple (1 port) | Complex (5+ ports) |
| **Dependencies** | Node.js only | Python, Flutter, Docker |
| **Backup** | Copy files | Database dumps |
| **Updates** | git pull | Rebuild containers |
| **Monitoring** | Process logs | Multi-service monitoring |

---

## Conclusion

Both systems represent sophisticated approaches to personal AI assistants:

- **OpenClaw** excels at: Simplicity, chat-first UX, rapid deployment, self-modification
- **OneMind OS** excels at: Comprehensiveness, multi-platform apps, structured workflows, enterprise features

Choose based on:
- **Technical skill level**: OpenClaw (lower) vs OneMind OS (higher)
- **Use case**: Quick automation (OpenClaw) vs Life management (OneMind OS)
- **Infrastructure preference**: Minimal (OpenClaw) vs Full-stack (OneMind OS)
- **Extensibility needs**: Self-modifying (OpenClaw) vs Type-safe (OneMind OS)

---

**Document Version**: 1.0
**Last Updated**: 2026-02-07
**Pages**: 34
