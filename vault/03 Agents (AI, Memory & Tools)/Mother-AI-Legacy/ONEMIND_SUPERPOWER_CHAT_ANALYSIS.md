# OneMind Superpower Chat System - Deep Analysis

> Comprehensive research on LobeChat, LobeHub, modern AI chat design, and unified agent workforce architecture.

---

## Executive Summary

After deep analysis of LobeChat/LobeHub, your existing OneMind architecture, and modern AI chat paradigms, I've identified a path to create a **superpower unified structure** that combines:

1. **LobeHub's Agent Marketplace Model** - Curated, shareable, installable agents
2. **MCP Protocol** - Universal tool/plugin architecture
3. **AG-UI Protocol** - Standardized agent-frontend communication (you already have this!)
4. **Workforce Orchestration** - Multi-agent teams with specialization
5. **Flutter Native Implementation** - Cross-platform mobile-first design

**Key Insight**: OneMind is uniquely positioned because you already have:
- AG-UI protocol integration
- Tool card registry (Flutter MCP-UI)
- Multi-backend support (Agno, LangGraph, CrewAI)
- Conversation branching
- HITL approval flows

LobeChat is React/Web only. **You can build the Flutter-native superpower version**.

---

## Part 1: LobeChat/LobeHub Deep Dive

### 1.1 LobeHub Philosophy: "Agents as Units of Work"

LobeChat's core insight is treating **agents as the fundamental unit of work**, not just chat sessions:

```
Traditional Chat App:
  User → Chat Session → Messages → AI Response

LobeChat Model:
  User → Agent (with skills, tools, memory) → Conversations → Artifacts
         ↓
    Agent Marketplace (discover, install, share)
         ↓
    Agent Groups (multi-agent collaboration)
```

### 1.2 LobeHub Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              LOBEHUB ECOSYSTEM                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────┐    ┌───────────────────┐    ┌───────────────────┐   │
│  │   AGENT HUB       │    │   PLUGIN HUB      │    │   MCP MARKETPLACE │   │
│  │   (lobehub.com/   │    │   (lobehub.com/   │    │   (lobehub.com/   │   │
│  │    discover)      │    │    plugins)       │    │    mcp)           │   │
│  │                   │    │                   │    │                   │   │
│  │   500+ Agents     │    │   100+ Plugins    │    │   MCP Servers     │   │
│  │   Categories:     │    │   Function Call   │    │   - filesystem    │   │
│  │   - Programming   │    │   - Web Search    │    │   - github        │   │
│  │   - Translation   │    │   - Image Gen     │    │   - notion        │   │
│  │   - Writing       │    │   - Weather       │    │   - slack         │   │
│  │   - Research      │    │   - Calculator    │    │   - database      │   │
│  │   - Creative      │    │   - File Ops      │    │   - home-asst     │   │
│  └───────────────────┘    └───────────────────┘    └───────────────────┘   │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                           STATE MANAGEMENT (Zustand)                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────┐    ┌───────────────────┐    ┌───────────────────┐   │
│  │   SessionStore    │    │   AgentStore      │    │   ChatStore       │   │
│  │                   │    │                   │    │                   │   │
│  │   - sessions[]    │    │   - agents[]      │    │   - messages[]    │   │
│  │   - activeId      │    │   - activeAgent   │    │   - streaming     │   │
│  │   - groupId       │    │   - plugins[]     │    │   - toolCalls[]   │   │
│  │   - create()      │    │   - tools[]       │    │   - send()        │   │
│  │   - update()      │    │   - memory        │    │   - regenerate()  │   │
│  │   - remove()      │    │   - systemRole    │    │   - delete()      │   │
│  └───────────────────┘    └───────────────────┘    └───────────────────┘   │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                              UI LAYER (Lobe-UI)                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐       │
│  │  ChatItem    │ │ MessageInput │ │  Markdown    │ │  Highlighter │       │
│  │  (bubble)    │ │  (composer)  │ │  (renderer)  │ │  (code)      │       │
│  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.3 LobeChat Plugin Manifest Schema

The plugin system uses a declarative manifest that defines:

```json
{
  "identifier": "weather-plugin",           // Unique ID
  "meta": {
    "title": "Weather Plugin",
    "description": "Get weather forecasts",
    "avatar": "🌤️",
    "tags": ["weather", "utility"]
  },
  "api": [                                  // Function calling definitions
    {
      "name": "getWeather",
      "description": "Get weather for a location",
      "url": "https://api.example.com/weather",
      "parameters": {
        "type": "object",
        "properties": {
          "location": { "type": "string", "description": "City name" },
          "units": { "type": "string", "enum": ["celsius", "fahrenheit"] }
        },
        "required": ["location"]
      }
    }
  ],
  "ui": {                                   // Optional custom UI
    "url": "https://plugin.example.com/ui",
    "height": 200
  },
  "gateway": "https://plugin.example.com/gateway",  // API gateway
  "systemRole": "You can check weather using the getWeather function..."
}
```

### 1.4 LobeChat Agent Definition

Agents in LobeChat are defined with:

```json
{
  "identifier": "code-assistant",
  "meta": {
    "title": "Code Assistant",
    "description": "Expert programmer and code reviewer",
    "avatar": "👨‍💻",
    "tags": ["programming", "code-review", "debugging"]
  },
  "config": {
    "model": "gpt-4-turbo",
    "systemRole": "You are an expert programmer...",
    "temperature": 0.3,
    "plugins": ["web-search", "code-execution"],
    "tools": ["read_file", "write_file", "run_code"]
  }
}
```

---

## Part 2: OneMind Current Architecture Analysis

### 2.1 What You Already Have (Strengths)

| Feature | Status | Notes |
|---------|--------|-------|
| AG-UI Protocol | ✅ Complete | `agui_protocol.dart` - full event mapping |
| Tool Card Registry | ✅ Complete | Flutter MCP-UI with 12+ built-in cards |
| Enhanced Messages | ✅ Complete | Branching, HITL, tool results, reasoning |
| Multi-Backend | ✅ Complete | Agno, LangGraph, CrewAI support |
| HITL Approval | ✅ Complete | `hitl_approval_card.dart` |
| Conversation Branching | ✅ Complete | `branch_navigator.dart` |
| Voice Integration | ✅ Partial | Models ready, needs full implementation |
| Riverpod State | ✅ Complete | `enhanced_chat_provider.dart` |

### 2.2 What's Missing (Gaps to Fill)

| Feature | Status | Priority |
|---------|--------|----------|
| Agent Marketplace/Hub | ❌ Missing | HIGH |
| Agent Builder UI | ❌ Missing | HIGH |
| Plugin Marketplace | ❌ Missing | HIGH |
| Agent Groups (Multi-Agent) | ❌ Missing | HIGH |
| Personal Memory System | ❌ Missing | MEDIUM |
| Scheduling/Automation | ❌ Missing | MEDIUM |
| Pages (Collaborative Docs) | ❌ Missing | LOW |
| Workspaces (Team) | ❌ Missing | LOW |

---

## Part 3: The Superpower Unified Structure

### 3.1 OneMind Agent Workforce Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         ONEMIND AGENT WORKFORCE SYSTEM                           │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐    │
│  │                           ONEMIND HUB                                    │    │
│  │                    (Central Discovery & Management)                      │    │
│  ├─────────────────────────────────────────────────────────────────────────┤    │
│  │                                                                          │    │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐          │    │
│  │  │  AGENT STORE    │  │  SKILL STORE    │  │  WORKFLOW STORE │          │    │
│  │  │  (Marketplace)  │  │  (MCP Tools)    │  │  (Teams)        │          │    │
│  │  │                 │  │                 │  │                 │          │    │
│  │  │  • Browse       │  │  • Filesystem   │  │  • Predefined   │          │    │
│  │  │  • Install      │  │  • Web Search   │  │  • Custom       │          │    │
│  │  │  • Rate/Review  │  │  • Home Asst    │  │  • Templates    │          │    │
│  │  │  • Create       │  │  • Calendar     │  │  • Sharing      │          │    │
│  │  │  • Share        │  │  • Database     │  │                 │          │    │
│  │  │  • Fork         │  │  • Custom MCP   │  │                 │          │    │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────┘          │    │
│  │                                                                          │    │
│  └─────────────────────────────────────────────────────────────────────────┘    │
│                                                                                  │
├─────────────────────────────────────────────────────────────────────────────────┤
│                              AGENT LAYER                                         │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐     │
│  │                         AGENT DEFINITION                                │     │
│  ├────────────────────────────────────────────────────────────────────────┤     │
│  │                                                                         │     │
│  │  Agent {                                                                │     │
│  │    id: string                    // Unique identifier                   │     │
│  │    name: string                  // Display name                        │     │
│  │    avatar: string                // Icon/image                          │     │
│  │    category: AgentCategory       // coding, research, creative, etc.   │     │
│  │    systemPrompt: string          // Core instructions                   │     │
│  │    model: ModelConfig            // Default model + params              │     │
│  │    skills: List<Skill>           // MCP tools enabled                   │     │
│  │    memory: MemoryConfig          // What to remember                    │     │
│  │    triggers: List<Trigger>       // Auto-activation rules               │     │
│  │    ui: AgentUIConfig             // Custom UI hints                     │     │
│  │    permissions: Permissions      // What it can access                  │     │
│  │  }                                                                      │     │
│  │                                                                         │     │
│  │  Categories:                                                            │     │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐      │     │
│  │  │ Coding   │ │ Research │ │ Creative │ │ LifeOS   │ │ Home     │      │     │
│  │  │          │ │          │ │          │ │          │ │          │      │     │
│  │  │ • Debug  │ │ • Web    │ │ • Write  │ │ • Tasks  │ │ • Lights │      │     │
│  │  │ • Review │ │ • Papers │ │ • Design │ │ • Goals  │ │ • Climate│      │     │
│  │  │ • Refact │ │ • Data   │ │ • Music  │ │ • Habits │ │ • Security│     │     │
│  │  │ • Docs   │ │ • Fact   │ │ • Art    │ │ • Journal│ │ • Cameras│      │     │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘      │     │
│  │                                                                         │     │
│  └────────────────────────────────────────────────────────────────────────┘     │
│                                                                                  │
├─────────────────────────────────────────────────────────────────────────────────┤
│                              WORKFORCE LAYER                                     │
│                         (Multi-Agent Orchestration)                              │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐     │
│  │                           WORKFORCE TEAMS                               │     │
│  ├────────────────────────────────────────────────────────────────────────┤     │
│  │                                                                         │     │
│  │  Team {                                                                 │     │
│  │    id: string                                                           │     │
│  │    name: string                  // "Development Team"                  │     │
│  │    agents: List<Agent>           // Team members                        │     │
│  │    coordinator: Agent?           // Optional orchestrator               │     │
│  │    executionMode: ExecutionMode  // sequential | parallel | adaptive   │     │
│  │    handoffRules: List<Rule>      // When to pass to another agent       │     │
│  │    sharedMemory: SharedMemory    // Team-level context                  │     │
│  │  }                                                                      │     │
│  │                                                                         │     │
│  │  Execution Patterns:                                                    │     │
│  │                                                                         │     │
│  │  SEQUENTIAL:    User → Agent A → Agent B → Agent C → Response          │     │
│  │                                                                         │     │
│  │  PARALLEL:      User → ┬→ Agent A ─┐                                   │     │
│  │                        ├→ Agent B ─┼→ Merge → Response                 │     │
│  │                        └→ Agent C ─┘                                   │     │
│  │                                                                         │     │
│  │  ADAPTIVE:      User → Coordinator → Routes to best agent              │     │
│  │                        ↓                                                │     │
│  │                 Monitors → Escalates/Handoffs as needed                │     │
│  │                                                                         │     │
│  │  Example Teams:                                                         │     │
│  │  ┌─────────────────────┐  ┌─────────────────────┐                      │     │
│  │  │ Development Team    │  │ Research Team       │                      │     │
│  │  │                     │  │                     │                      │     │
│  │  │ • Code Writer       │  │ • Web Researcher    │                      │     │
│  │  │ • Code Reviewer     │  │ • Paper Analyst     │                      │     │
│  │  │ • Test Engineer     │  │ • Fact Checker      │                      │     │
│  │  │ • Doc Writer        │  │ • Synthesizer       │                      │     │
│  │  └─────────────────────┘  └─────────────────────┘                      │     │
│  │                                                                         │     │
│  └────────────────────────────────────────────────────────────────────────┘     │
│                                                                                  │
├─────────────────────────────────────────────────────────────────────────────────┤
│                              SKILL/TOOL LAYER                                    │
│                              (MCP Integration)                                   │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐     │
│  │                          SKILL DEFINITION                               │     │
│  ├────────────────────────────────────────────────────────────────────────┤     │
│  │                                                                         │     │
│  │  Skill {                                                                │     │
│  │    id: string                    // "home-assistant"                    │     │
│  │    name: string                  // "Home Assistant Control"            │     │
│  │    type: SkillType               // mcp_server | function | api        │     │
│  │    mcpServer?: MCPServerConfig   // If MCP-based                        │     │
│  │    tools: List<Tool>             // Available functions                 │     │
│  │    permissions: List<Permission> // Required access                     │     │
│  │    uiCards: Map<String, CardDef> // Tool → Flutter widget mapping      │     │
│  │  }                                                                      │     │
│  │                                                                         │     │
│  │  Tool {                                                                 │     │
│  │    name: string                  // "turn_on_lights"                    │     │
│  │    description: string           // For LLM function calling            │     │
│  │    parameters: JSONSchema        // Input schema                        │     │
│  │    resultSchema: JSONSchema?     // Output schema                       │     │
│  │    uiHint: String?               // "LightControlCard"                  │     │
│  │  }                                                                      │     │
│  │                                                                         │     │
│  │  Built-in Skills:                                                       │     │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐      │     │
│  │  │filesystem│ │web_search│ │home_asst │ │ calendar │ │ memory   │      │     │
│  │  │          │ │          │ │          │ │          │ │          │      │     │
│  │  │read_file │ │search    │ │lights    │ │get_events│ │remember  │      │     │
│  │  │write_file│ │extract   │ │climate   │ │create    │ │recall    │      │     │
│  │  │list_dir  │ │crawl     │ │sensors   │ │update    │ │forget    │      │     │
│  │  │search    │ │summarize │ │automats  │ │delete    │ │search    │      │     │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘      │     │
│  │                                                                         │     │
│  └────────────────────────────────────────────────────────────────────────┘     │
│                                                                                  │
├─────────────────────────────────────────────────────────────────────────────────┤
│                              MEMORY LAYER                                        │
│                          (Personal & Shared Memory)                              │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐     │
│  │                          MEMORY SYSTEM                                  │     │
│  ├────────────────────────────────────────────────────────────────────────┤     │
│  │                                                                         │     │
│  │  ┌─────────────────────────────┐  ┌─────────────────────────────┐      │     │
│  │  │     PERSONAL MEMORY         │  │     AGENT MEMORY            │      │     │
│  │  │                             │  │                             │      │     │
│  │  │  User Preferences:          │  │  Per-Agent Context:         │      │     │
│  │  │  • Communication style      │  │  • Learned patterns         │      │     │
│  │  │  • Technical level          │  │  • Successful approaches    │      │     │
│  │  │  • Time zone                │  │  • User corrections         │      │     │
│  │  │  • Name, pronouns           │  │  • Tool usage stats         │      │     │
│  │  │                             │  │                             │      │     │
│  │  │  Facts & Knowledge:         │  │  Session Memory:            │      │     │
│  │  │  • Work projects            │  │  • Current task context     │      │     │
│  │  │  • Family members           │  │  • Open files/docs          │      │     │
│  │  │  • Tech stack               │  │  • Recent decisions         │      │     │
│  │  │  • Interests                │  │                             │      │     │
│  │  └─────────────────────────────┘  └─────────────────────────────┘      │     │
│  │                                                                         │     │
│  │  Memory Operations:                                                     │     │
│  │  • remember(fact, category, confidence)                                │     │
│  │  • recall(query) → List<Fact>                                          │     │
│  │  • forget(factId)                                                       │     │
│  │  • updateConfidence(factId, delta)                                      │     │
│  │                                                                         │     │
│  │  Storage: Vector DB (embeddings) + Structured DB (facts)               │     │
│  │                                                                         │     │
│  └────────────────────────────────────────────────────────────────────────┘     │
│                                                                                  │
├─────────────────────────────────────────────────────────────────────────────────┤
│                           CONVERSATION LAYER                                     │
│                          (Enhanced Chat Engine)                                  │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐     │
│  │                      CONVERSATION FEATURES                              │     │
│  ├────────────────────────────────────────────────────────────────────────┤     │
│  │                                                                         │     │
│  │  ✅ Message Types:              ✅ Advanced Features:                   │     │
│  │  • User messages               • Conversation branching                │     │
│  │  • Assistant messages          • Chain of thought visualization        │     │
│  │  • System messages             • HITL approval flows                   │     │
│  │  • Tool calls + results        • Voice input/output                    │     │
│  │                                • File attachments                       │     │
│  │  ✅ Tool Result Rendering:     • @mentions for agents                  │     │
│  │  • Calendar cards              • /slash commands                       │     │
│  │  • Task cards                  • Real-time streaming                   │     │
│  │  • Weather cards               • Multi-model switching                 │     │
│  │  • Code execution cards                                                 │     │
│  │  • Search result cards         ✅ Actions:                              │     │
│  │  • Generic data cards          • Regenerate                            │     │
│  │  • Custom Generative UI        • Edit & resend                         │     │
│  │                                • Branch from message                   │     │
│  │                                • Copy, share, delete                   │     │
│  │                                • Speak (TTS)                           │     │
│  │                                                                         │     │
│  └────────────────────────────────────────────────────────────────────────┘     │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Unified Data Models (Flutter/Dart)

```dart
// ============================================================================
// AGENT DEFINITION
// ============================================================================

@freezed
sealed class OneMindAgent with _$OneMindAgent {
  const factory OneMindAgent({
    required String id,
    required String name,
    required String description,
    String? avatar,
    @Default([]) List<String> tags,
    required AgentCategory category,
    required AgentConfig config,
    @Default([]) List<String> skillIds,      // MCP skills/tools enabled
    MemoryConfig? memoryConfig,
    AgentUIConfig? uiConfig,
    AgentMetrics? metrics,                    // Usage stats, ratings
    String? authorId,
    @Default(false) bool isPublic,
    @Default(false) bool isInstalled,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _OneMindAgent;
}

@freezed
sealed class AgentConfig with _$AgentConfig {
  const factory AgentConfig({
    required String systemPrompt,
    @Default('claude-sonnet') String defaultModel,
    @Default(0.7) double temperature,
    @Default(4096) int maxTokens,
    @Default([]) List<String> enabledTools,
    @Default({}) Map<String, dynamic> modelParams,
    String? openingMessage,
    @Default([]) List<String> suggestedPrompts,
  }) = _AgentConfig;
}

enum AgentCategory {
  coding,
  research,
  creative,
  productivity,
  lifeos,
  home,
  business,
  education,
  health,
  custom,
}

// ============================================================================
// WORKFORCE / TEAM DEFINITION
// ============================================================================

@freezed
sealed class Workforce with _$Workforce {
  const factory Workforce({
    required String id,
    required String name,
    required String description,
    @Default([]) List<String> agentIds,       // Team members
    String? coordinatorAgentId,               // Optional orchestrator
    required ExecutionMode executionMode,
    @Default([]) List<HandoffRule> handoffRules,
    SharedMemory? sharedMemory,
    @Default({}) Map<String, dynamic> metadata,
  }) = _Workforce;
}

enum ExecutionMode {
  sequential,    // A → B → C
  parallel,      // A, B, C simultaneously
  adaptive,      // Coordinator routes dynamically
  hierarchical,  // Manager → Workers
}

@freezed
sealed class HandoffRule with _$HandoffRule {
  const factory HandoffRule({
    required String fromAgentId,
    required String toAgentId,
    required String condition,              // Natural language or expression
    @Default(false) bool requiresApproval,  // HITL for handoff
  }) = _HandoffRule;
}

// ============================================================================
// SKILL / TOOL DEFINITION
// ============================================================================

@freezed
sealed class Skill with _$Skill {
  const factory Skill({
    required String id,
    required String name,
    required String description,
    String? icon,
    @Default([]) List<String> tags,
    required SkillType type,
    MCPServerConfig? mcpConfig,              // If MCP-based
    @Default([]) List<Tool> tools,
    @Default([]) List<String> permissions,
    @Default({}) Map<String, String> uiCardMappings,  // tool → widget
    @Default(false) bool isInstalled,
    @Default(false) bool isEnabled,
  }) = _Skill;
}

enum SkillType {
  mcpServer,      // Full MCP server
  builtIn,        // Native implementation
  apiPlugin,      // REST API integration
  customFunction, // Custom Dart function
}

@freezed
sealed class MCPServerConfig with _$MCPServerConfig {
  const factory MCPServerConfig({
    required String serverName,
    required String command,                  // e.g., "npx"
    @Default([]) List<String> args,           // e.g., ["-y", "@mcp/filesystem"]
    @Default({}) Map<String, String> env,
    String? workingDirectory,
    @Default('stdio') String transport,       // stdio | sse | websocket
  }) = _MCPServerConfig;
}

@freezed
sealed class Tool with _$Tool {
  const factory Tool({
    required String name,
    required String description,
    required Map<String, dynamic> parameters,  // JSON Schema
    Map<String, dynamic>? resultSchema,
    String? uiHint,                            // Widget type hint
    @Default([]) List<String> exampleInputs,
  }) = _Tool;
}

// ============================================================================
// MEMORY DEFINITION
// ============================================================================

@freezed
sealed class MemoryFact with _$MemoryFact {
  const factory MemoryFact({
    required String id,
    required String content,
    required MemoryCategory category,
    @Default(1.0) double confidence,
    String? source,                           // Where learned from
    @Default([]) List<String> relatedFactIds,
    DateTime? createdAt,
    DateTime? lastAccessedAt,
    @Default(0) int accessCount,
  }) = _MemoryFact;
}

enum MemoryCategory {
  userPreference,
  userFact,
  projectContext,
  learnedPattern,
  correction,
  temporary,
}

@freezed
sealed class MemoryConfig with _$MemoryConfig {
  const factory MemoryConfig({
    @Default(true) bool enableAutoMemory,
    @Default([]) List<MemoryCategory> allowedCategories,
    @Default(100) int maxFacts,
    @Default(0.5) double minConfidenceThreshold,
  }) = _MemoryConfig;
}
```

### 3.3 State Management Architecture (Riverpod)

```dart
// ============================================================================
// ONEMIND HUB PROVIDERS
// ============================================================================

// Agent Store
final agentStoreProvider = StateNotifierProvider<AgentStoreNotifier, AgentStoreState>((ref) {
  return AgentStoreNotifier(ref);
});

class AgentStoreState {
  final List<OneMindAgent> installedAgents;
  final List<OneMindAgent> marketplaceAgents;
  final String? activeAgentId;
  final bool isLoading;
  final String? error;
}

// Skill Store
final skillStoreProvider = StateNotifierProvider<SkillStoreNotifier, SkillStoreState>((ref) {
  return SkillStoreNotifier(ref);
});

class SkillStoreState {
  final List<Skill> installedSkills;
  final List<Skill> availableSkills;
  final Map<String, MCPSession> activeSessions;  // Running MCP servers
  final bool isLoading;
}

// Workforce Store
final workforceStoreProvider = StateNotifierProvider<WorkforceStoreNotifier, WorkforceStoreState>((ref) {
  return WorkforceStoreNotifier(ref);
});

class WorkforceStoreState {
  final List<Workforce> teams;
  final String? activeTeamId;
  final WorkforceExecution? currentExecution;
}

// Memory Store
final memoryStoreProvider = StateNotifierProvider<MemoryStoreNotifier, MemoryStoreState>((ref) {
  return MemoryStoreNotifier(ref);
});

class MemoryStoreState {
  final List<MemoryFact> facts;
  final Map<String, List<MemoryFact>> agentMemories;
  final bool isLoading;
}

// ============================================================================
// CROSS-CUTTING PROVIDERS
// ============================================================================

// Active agent's tools
final activeAgentToolsProvider = Provider<List<Tool>>((ref) {
  final agentState = ref.watch(agentStoreProvider);
  final skillState = ref.watch(skillStoreProvider);

  if (agentState.activeAgentId == null) return [];

  final agent = agentState.installedAgents
      .firstWhere((a) => a.id == agentState.activeAgentId);

  return agent.skillIds
      .expand((skillId) => skillState.installedSkills
          .firstWhere((s) => s.id == skillId, orElse: () => const Skill.empty())
          .tools)
      .toList();
});

// Context for current conversation
final conversationContextProvider = Provider<ConversationContext>((ref) {
  final agent = ref.watch(activeAgentProvider);
  final memory = ref.watch(relevantMemoryProvider);
  final tools = ref.watch(activeAgentToolsProvider);

  return ConversationContext(
    systemPrompt: agent?.config.systemPrompt ?? '',
    memory: memory,
    availableTools: tools,
  );
});
```

---

## Part 4: Modern Chat Design Principles

### 4.1 The Modern AI Chat Era Paradigm

Based on analyzing Claude, ChatGPT, Grok, LobeChat, and others, here are the design principles:

#### 1. **Agent-First, Not Chat-First**
```
Old Model: Chat → AI responds
New Model: Agent (with capabilities) → Conversations → Artifacts
```

#### 2. **Tools as First-Class Citizens**
```
Old Model: AI generates text response
New Model: AI uses tools → Results rendered as rich UI cards
```

#### 3. **Memory & Continuity**
```
Old Model: Each conversation is isolated
New Model: Personal memory spans all conversations
           Agent memory provides specialized context
```

#### 4. **Collaboration & Teams**
```
Old Model: One AI, one conversation
New Model: Multiple specialized agents work together
           Handoffs, parallel execution, coordination
```

#### 5. **Transparency & Control**
```
Old Model: Black box responses
New Model: Chain of thought visible
           HITL approval for important actions
           Branching to explore alternatives
```

### 4.2 UI/UX Patterns from Research

| Pattern | Source | Implementation |
|---------|--------|----------------|
| Model Router Chips | Claude, ChatGPT | Horizontal chip row for quick model switching |
| Quick Action Pills | Claude | "Code", "Research", "Write" buttons above input |
| Reasoning Drawer | Grok, LobeChat | Collapsible "thinking" section |
| Tool Cards | MCP-UI, LobeChat | Rich widgets for tool results |
| Branch Navigator | LobeChat | Tree visualization for conversation forks |
| Agent Mentions | Slack, LobeChat | @agent syntax in messages |
| Slash Commands | Discord, Claude | /command syntax for actions |
| Context Sidebar | Motion | Calendar, tasks, home state awareness |

### 4.3 The OneMind Differentiator

What makes OneMind's approach unique:

1. **Flutter Native** - Not React/Web, true mobile-first
2. **LifeOS Integration** - Tasks, habits, goals, calendar, home automation
3. **AG-UI Protocol** - Multi-backend support already built
4. **Real-World Awareness** - Home state, location, time, wellness data
5. **Unified Workforce** - Agents collaborate on your life, not just chat

---

## Part 5: Implementation Roadmap

### Phase 1: Foundation (You're Here)
- [x] AG-UI Protocol integration
- [x] Tool Card Registry (Flutter MCP-UI)
- [x] Enhanced Message model
- [x] Branching conversations
- [x] HITL approval UI
- [ ] **Agent data model** ← Next
- [ ] **Skill data model** ← Next

### Phase 2: OneMind Hub
- [ ] Agent Store (local + cloud)
- [ ] Skill Store (MCP management)
- [ ] Hub UI (marketplace browsing)
- [ ] Agent installation flow
- [ ] Skill installation flow

### Phase 3: Agent Builder
- [ ] No-code agent builder UI
- [ ] System prompt editor
- [ ] Skill/tool selector
- [ ] Model configuration
- [ ] Test playground

### Phase 4: Workforce System
- [ ] Team definition model
- [ ] Execution orchestrator
- [ ] Handoff rules engine
- [ ] Shared memory
- [ ] Team templates

### Phase 5: Memory System
- [ ] Personal memory store
- [ ] Agent memory store
- [ ] Vector embeddings
- [ ] Memory UI (view/edit/delete)

### Phase 6: Advanced Features
- [ ] Scheduled agents (cron)
- [ ] Pages (multi-agent docs)
- [ ] Workspaces (team collaboration)
- [ ] Agent analytics

---

## Part 6: Key Files to Create

```
frontend/lib/
├── onemind_hub/
│   ├── models/
│   │   ├── agent.dart              # OneMindAgent, AgentConfig
│   │   ├── skill.dart              # Skill, Tool, MCPServerConfig
│   │   ├── workforce.dart          # Workforce, HandoffRule
│   │   └── memory.dart             # MemoryFact, MemoryConfig
│   │
│   ├── providers/
│   │   ├── agent_store_provider.dart
│   │   ├── skill_store_provider.dart
│   │   ├── workforce_provider.dart
│   │   └── memory_provider.dart
│   │
│   ├── services/
│   │   ├── agent_service.dart      # CRUD, sync with backend
│   │   ├── skill_service.dart      # MCP server management
│   │   ├── workforce_service.dart  # Orchestration logic
│   │   └── memory_service.dart     # Vector search, storage
│   │
│   ├── screens/
│   │   ├── hub_screen.dart         # Main marketplace
│   │   ├── agent_store_screen.dart # Browse agents
│   │   ├── skill_store_screen.dart # Browse skills
│   │   ├── agent_detail_screen.dart
│   │   ├── agent_builder_screen.dart
│   │   └── workforce_screen.dart   # Team management
│   │
│   └── widgets/
│       ├── agent_card.dart
│       ├── skill_card.dart
│       ├── agent_builder/
│       │   ├── prompt_editor.dart
│       │   ├── skill_picker.dart
│       │   └── model_config.dart
│       └── workforce/
│           ├── team_visualizer.dart
│           └── execution_monitor.dart
```

---

## Conclusion

OneMind has a unique opportunity to be the **Flutter-native superpower chat** that combines:

1. **LobeHub's Agent Marketplace** model for discovery and sharing
2. **MCP Protocol** for universal tool integration
3. **AG-UI Protocol** for multi-backend support (already built!)
4. **Workforce Orchestration** for multi-agent collaboration
5. **LifeOS Integration** for real-world awareness and action

The key insight is that you're not building "another chat app" - you're building a **unified agent workforce platform** that happens to have chat as its primary interface.

Next steps:
1. Implement Agent and Skill data models
2. Build the OneMind Hub UI
3. Create the Agent Builder
4. Add Workforce orchestration

Let me know when you're ready to start implementation!
