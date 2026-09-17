# Frontend Architecture

**Version:** 2.0
**Last Updated:** 2026-02-06
**Framework:** Flutter 3.10+ (Web)

---

## Table of Contents

1. [Overview](#overview)
2. [Project Structure](#project-structure)
3. [State Management with Riverpod](#state-management-with-riverpod)
4. [Service Layer Architecture](#service-layer-architecture)
5. [Model Classes](#model-classes)
6. [Screen Organization](#screen-organization)
7. [Widget Reusability](#widget-reusability)
8. [Routing with GoRouter](#routing-with-gorouter)
9. [Theme and Styling](#theme-and-styling)
10. [WebSocket Integration](#websocket-integration)
11. [Error Handling Patterns](#error-handling-patterns)
12. [Loading States](#loading-states)
13. [Form Validation](#form-validation)
14. [Best Practices](#best-practices)

---

## Overview

OneMind OS v2 frontend is a **Flutter web application** built with:

- **Language:** Dart 3.0+
- **Framework:** Flutter 3.10+ (Web target)
- **State Management:** Riverpod 2.x
- **Routing:** GoRouter 13.x
- **HTTP Client:** http package
- **WebSocket:** web_socket_channel package
- **Markdown Rendering:** flutter_markdown package

### Key Architectural Decisions

1. **Riverpod for state management** - Type-safe, testable, compile-time DI
2. **Service layer pattern** - Separation of concerns (API calls vs UI logic)
3. **Immutable models** - Data integrity and predictability
4. **Provider-per-resource** - Fine-grained reactivity
5. **WebSocket for real-time** - Efficient bi-directional communication
6. **Material Design 3** - Modern, accessible UI components

---

## Project Structure

```
frontend/
├── lib/
│   ├── main.dart                          # App entry point + routing
│   │
│   ├── screens/                           # Full-page views
│   │   ├── chat_screen.dart              # Main chat interface
│   │   ├── agents_screen.dart            # Agent list view
│   │   ├── agent_form_screen.dart        # Agent create/edit form
│   │   ├── teams_screen.dart             # Team list view
│   │   ├── team_form_screen.dart         # Team create/edit form
│   │   ├── memories_screen.dart          # Memory management
│   │   ├── sessions_screen.dart          # Session history
│   │   ├── knowledge_screen.dart         # Knowledge base UI
│   │   ├── workflows_screen.dart         # Workflow management
│   │   ├── evaluations_screen.dart       # Evaluation dashboard
│   │   └── analytics_screen.dart         # Analytics & metrics
│   │
│   ├── components/                        # Reusable widgets
│   │   ├── markdown_message.dart         # Message rendering
│   │   ├── typing_indicator.dart         # Chat typing animation
│   │   ├── agent_status_indicator.dart   # Agent state display
│   │   ├── settings_panel.dart           # Settings sidebar
│   │   ├── message_actions.dart          # Message action buttons
│   │   ├── quick_settings_drawer.dart    # Quick settings UI
│   │   ├── quick_info_drawer.dart        # Info drawer UI
│   │   ├── realtime_indicator.dart       # WebSocket status
│   │   └── realtime_toast.dart           # Toast notifications
│   │
│   ├── services/                          # Business logic layer
│   │   ├── api_service.dart              # REST API client
│   │   ├── websocket_service.dart        # WebSocket client
│   │   ├── agent_service.dart            # Agent operations
│   │   ├── enhanced_agent_service.dart   # Advanced agent features
│   │   └── export_service.dart           # Data export utilities
│   │
│   ├── providers/                         # Riverpod state providers
│   │   ├── api_providers.dart            # API data providers
│   │   ├── websocket_providers.dart      # WebSocket providers
│   │   └── agent_config_provider.dart    # Agent config state
│   │
│   └── models/                            # Data models
│       ├── agent_model.dart              # Agent entity
│       ├── team_model.dart               # Team entity
│       ├── tool_model.dart               # Tool metadata
│       ├── tool_config.dart              # Tool configuration
│       ├── model_info.dart               # LLM model info
│       ├── memory_model.dart             # Memory entity
│       ├── message.dart                  # Chat message
│       ├── agent_event.dart              # Agent events
│       ├── session_model.dart            # Session entity
│       ├── knowledge_model.dart          # Knowledge base
│       ├── workflow_model.dart           # Workflow entity
│       ├── approval_model.dart           # Approval request
│       └── analytics_model.dart          # Analytics data
│
├── pubspec.yaml                           # Dependencies
├── analysis_options.yaml                  # Linter rules
└── web/                                   # Web-specific assets
    ├── index.html
    └── favicon.png
```

---

## State Management with Riverpod

### Why Riverpod?

Riverpod is a **compile-time safe**, **testable**, and **scalable** state management solution.

**Advantages over Provider:**
- No `BuildContext` required
- Compile-time safety (no runtime errors)
- Better testing support
- Automatic disposal
- Dependency injection built-in

### Provider Types

#### 1. FutureProvider - Async data fetching

**Use case:** Load data from API once

```dart
/// Fetches list of agents from backend
final agentsProvider = FutureProvider<List<AgentModel>>((ref) async {
  return await ApiService.listAgents();
});
```

**Consumer usage:**
```dart
Consumer(
  builder: (context, ref, child) {
    final agentsAsync = ref.watch(agentsProvider);

    return agentsAsync.when(
      data: (agents) => ListView.builder(
        itemCount: agents.length,
        itemBuilder: (context, index) => AgentTile(agents[index]),
      ),
      loading: () => CircularProgressIndicator(),
      error: (err, stack) => Text('Error: $err'),
    );
  },
)
```

#### 2. FutureProvider.family - Parameterized async data

**Use case:** Load data based on an ID or parameter

```dart
/// Fetches specific agent by ID
final agentProvider = FutureProvider.family<AgentModel, String>((ref, agentId) async {
  return await ApiService.getAgent(agentId);
});
```

**Consumer usage:**
```dart
final agentAsync = ref.watch(agentProvider('agent-123'));
```

#### 3. Provider - Synchronous data or services

**Use case:** Provide access to services or mutations

```dart
/// Provides agent mutation operations
final agentMutationsProvider = Provider((ref) => AgentMutations(ref));

class AgentMutations {
  final Ref ref;

  AgentMutations(this.ref);

  Future<AgentModel> createAgent(AgentModel agent) async {
    final newAgent = await ApiService.createAgent(agent);
    // Invalidate cache to refresh UI
    ref.invalidate(agentsProvider);
    return newAgent;
  }

  Future<void> deleteAgent(String agentId) async {
    await ApiService.deleteAgent(agentId);
    ref.invalidate(agentsProvider);
  }
}
```

**Consumer usage:**
```dart
final mutations = ref.read(agentMutationsProvider);
await mutations.createAgent(newAgent);
```

#### 4. StreamProvider - Real-time data streams

**Use case:** WebSocket events, Firebase streams

```dart
/// Streams WebSocket events
final websocketEventsProvider = StreamProvider<WebSocketEvent>((ref) {
  final wsService = ref.watch(websocketServiceProvider);
  return wsService.eventStream;
});
```

**Consumer usage:**
```dart
final eventsAsync = ref.watch(websocketEventsProvider);

eventsAsync.when(
  data: (event) => Text('Event: ${event.type}'),
  loading: () => Text('Connecting...'),
  error: (err, stack) => Text('Disconnected'),
);
```

#### 5. StateNotifierProvider - Mutable state

**Use case:** Local UI state, forms, filters

```dart
class FilterState extends StateNotifier<String> {
  FilterState() : super('');

  void setFilter(String value) {
    state = value;
  }
}

final filterProvider = StateNotifierProvider<FilterState, String>((ref) {
  return FilterState();
});
```

**Consumer usage:**
```dart
// Read state
final filter = ref.watch(filterProvider);

// Update state
ref.read(filterProvider.notifier).setFilter('research');
```

### Provider Patterns in OneMind OS

#### Pattern 1: List + Detail + Mutations

```dart
// List all resources
final agentsProvider = FutureProvider<List<AgentModel>>((ref) async {
  return await ApiService.listAgents();
});

// Get single resource by ID
final agentProvider = FutureProvider.family<AgentModel, String>((ref, id) async {
  return await ApiService.getAgent(id);
});

// Mutation operations
final agentMutationsProvider = Provider((ref) => AgentMutations(ref));
```

#### Pattern 2: Filtered Lists

```dart
// Memories filtered by agent
final memoriesByAgentProvider = FutureProvider.family<List<MemoryModel>, String?>(
  (ref, agentId) async {
    return await ApiService.listMemories(agentId: agentId);
  }
);
```

#### Pattern 3: Cache Invalidation

```dart
class AgentMutations {
  Future<AgentModel> updateAgent(String id, AgentModel agent) async {
    final updated = await ApiService.updateAgent(id, agent);

    // Invalidate both list and detail
    ref.invalidate(agentsProvider);
    ref.invalidate(agentProvider(id));

    return updated;
  }
}
```

---

## Service Layer Architecture

### ApiService - HTTP Client

**File:** `lib/services/api_service.dart`

**Purpose:** Centralized HTTP client for all backend API calls

**Structure:**
```dart
class ApiService {
  static const String baseUrl = 'http://localhost:7777';

  // Agent APIs
  static Future<List<AgentModel>> listAgents() async { ... }
  static Future<AgentModel> getAgent(String agentId) async { ... }
  static Future<AgentModel> createAgent(AgentModel agent) async { ... }
  static Future<AgentModel> updateAgent(String id, AgentModel agent) async { ... }
  static Future<void> deleteAgent(String agentId) async { ... }

  // Team APIs
  static Future<List<TeamModel>> listTeams() async { ... }
  // ... (similar CRUD operations)

  // Tool APIs
  static Future<List<ToolModel>> listTools() async { ... }

  // Model APIs
  static Future<List<ModelInfo>> listModels({String? provider}) async { ... }

  // Memory APIs
  static Future<List<MemoryModel>> listMemories({String? agentId}) async { ... }

  // Knowledge APIs
  static Future<List<KnowledgeBaseModel>> listKnowledgeBases() async { ... }

  // Session APIs
  static Future<List<SessionModel>> listSessions({String? agentId}) async { ... }

  // Workflow APIs
  static Future<List<WorkflowModel>> listWorkflows() async { ... }

  // Approval APIs
  static Future<List<ApprovalRequest>> listApprovals() async { ... }
}
```

**Example Implementation:**
```dart
static Future<List<AgentModel>> listAgents() async {
  final response = await http.get(Uri.parse('$baseUrl/api/agent-configs/'));

  if (response.statusCode == 200) {
    final List<dynamic> data = json.decode(response.body);
    return data.map((json) => AgentModel.fromJson(json)).toList();
  } else {
    throw Exception('Failed to load agents: ${response.body}');
  }
}
```

**Error Handling:**
```dart
try {
  final agents = await ApiService.listAgents();
  // Success
} on SocketException {
  // Network error
  print('No internet connection');
} on HttpException {
  // HTTP error
  print('Server error');
} on FormatException {
  // JSON parsing error
  print('Invalid response format');
} catch (e) {
  // Unknown error
  print('Unexpected error: $e');
}
```

### WebSocketService - Real-time Updates

**File:** `lib/services/websocket_service.dart`

**Purpose:** Manage WebSocket connections for real-time events

**Features:**
- Auto-reconnect on disconnect
- Topic-based subscriptions
- Event stream handlers
- Connection status monitoring

**Example Usage:**
```dart
final wsService = WebSocketService();

// Connect to sessions topic
await wsService.connect(topic: WebSocketTopic.sessions);

// Listen to events
wsService.eventStream.listen((event) {
  print('Event: ${event.eventType}');
  print('Data: ${event.data}');
});

// Disconnect
await wsService.disconnect();
```

---

## Model Classes

### Immutable Data Models

All models are **immutable** (`const` constructors) for predictability and performance.

**Example: AgentModel**

```dart
class AgentModel {
  // Identity
  final String? id;
  final String name;
  final String? description;
  final String? role;

  // Model configuration
  final String modelId;
  final String modelName;

  // Instructions & Tools
  final List<String> instructions;
  final List<String> tools;

  // Flags
  final bool markdown;
  final bool addHistoryToContext;
  final int? numHistoryRuns;

  const AgentModel({
    this.id,
    required this.name,
    this.description,
    this.role,
    required this.modelId,
    required this.modelName,
    this.instructions = const [],
    this.tools = const [],
    this.markdown = true,
    this.addHistoryToContext = true,
    this.numHistoryRuns = 5,
  });

  // JSON serialization
  factory AgentModel.fromJson(Map<String, dynamic> json) {
    return AgentModel(
      id: json['id'],
      name: json['name'],
      description: json['description'],
      role: json['role'],
      modelId: json['model_id'],
      modelName: json['model_name'],
      instructions: List<String>.from(json['instructions'] ?? []),
      tools: List<String>.from(json['tools'] ?? []),
      markdown: json['markdown'] ?? true,
      addHistoryToContext: json['add_history_to_context'] ?? true,
      numHistoryRuns: json['num_history_runs'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      if (id != null) 'id': id,
      'name': name,
      'description': description,
      'role': role,
      'model_id': modelId,
      'model_name': modelName,
      'instructions': instructions,
      'tools': tools,
      'markdown': markdown,
      'add_history_to_context': addHistoryToContext,
      'num_history_runs': numHistoryRuns,
    };
  }

  // copyWith for immutable updates
  AgentModel copyWith({
    String? id,
    String? name,
    String? description,
    String? role,
    String? modelId,
    String? modelName,
    List<String>? instructions,
    List<String>? tools,
    bool? markdown,
    bool? addHistoryToContext,
    int? numHistoryRuns,
  }) {
    return AgentModel(
      id: id ?? this.id,
      name: name ?? this.name,
      description: description ?? this.description,
      role: role ?? this.role,
      modelId: modelId ?? this.modelId,
      modelName: modelName ?? this.modelName,
      instructions: instructions ?? this.instructions,
      tools: tools ?? this.tools,
      markdown: markdown ?? this.markdown,
      addHistoryToContext: addHistoryToContext ?? this.addHistoryToContext,
      numHistoryRuns: numHistoryRuns ?? this.numHistoryRuns,
    );
  }
}
```

### Model Conventions

1. **Immutability** - Use `const` constructors and `final` fields
2. **JSON serialization** - Implement `fromJson` and `toJson`
3. **copyWith method** - For creating modified copies
4. **Null safety** - Use `?` for optional fields
5. **Default values** - Provide sensible defaults in constructor
6. **Field naming** - Use camelCase (Dart convention)

---

## Screen Organization

### Screen Anatomy

Each screen follows a consistent structure:

```dart
class ExampleScreen extends ConsumerWidget {
  const ExampleScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // 1. Watch providers
    final dataAsync = ref.watch(dataProvider);

    // 2. Build scaffold
    return Scaffold(
      appBar: _buildAppBar(context, ref),
      body: _buildBody(context, ref, dataAsync),
      floatingActionButton: _buildFAB(context, ref),
    );
  }

  // 3. Helper methods
  PreferredSizeWidget _buildAppBar(BuildContext context, WidgetRef ref) { ... }
  Widget _buildBody(BuildContext context, WidgetRef ref, AsyncValue<Data> data) { ... }
  Widget _buildFAB(BuildContext context, WidgetRef ref) { ... }
}
```

### Example: AgentsScreen

```dart
class AgentsScreen extends ConsumerWidget {
  const AgentsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final agentsAsync = ref.watch(agentsProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Agents'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => ref.refresh(agentsProvider),
          ),
        ],
      ),
      body: agentsAsync.when(
        data: (agents) => _buildAgentsList(context, ref, agents),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, stack) => _buildError(err),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/agents/create'),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildAgentsList(BuildContext context, WidgetRef ref, List<AgentModel> agents) {
    if (agents.isEmpty) {
      return const Center(child: Text('No agents found'));
    }

    return ListView.builder(
      itemCount: agents.length,
      itemBuilder: (context, index) {
        final agent = agents[index];
        return ListTile(
          leading: CircleAvatar(child: Text(agent.name[0])),
          title: Text(agent.name),
          subtitle: Text(agent.description ?? 'No description'),
          trailing: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              IconButton(
                icon: const Icon(Icons.edit),
                onPressed: () => context.go('/agents/edit/${agent.id}'),
              ),
              IconButton(
                icon: const Icon(Icons.delete),
                onPressed: () => _confirmDelete(context, ref, agent),
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildError(Object error) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.error, size: 64, color: Colors.red),
          const SizedBox(height: 16),
          Text('Error: $error'),
        ],
      ),
    );
  }

  Future<void> _confirmDelete(BuildContext context, WidgetRef ref, AgentModel agent) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Agent'),
        content: Text('Are you sure you want to delete "${agent.name}"?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('Delete'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final mutations = ref.read(agentMutationsProvider);
      await mutations.deleteAgent(agent.id!);

      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Deleted ${agent.name}')),
        );
      }
    }
  }
}
```

---

## Widget Reusability

### Component Structure

Reusable widgets are placed in `lib/components/` and follow these principles:

1. **Single responsibility** - One widget, one purpose
2. **Parameterization** - Accept data via constructor
3. **Composition** - Build complex UIs from simple widgets
4. **Const constructors** - Enable performance optimizations

### Example: MarkdownMessage

**File:** `lib/components/markdown_message.dart`

```dart
/// Renders markdown-formatted messages with syntax highlighting
class MarkdownMessage extends StatelessWidget {
  final String content;
  final TextStyle? style;

  const MarkdownMessage({
    super.key,
    required this.content,
    this.style,
  });

  @override
  Widget build(BuildContext context) {
    return Markdown(
      data: content,
      styleSheet: MarkdownStyleSheet(
        p: style ?? Theme.of(context).textTheme.bodyMedium,
        code: TextStyle(
          fontFamily: 'monospace',
          backgroundColor: Colors.grey[800],
          fontSize: 13,
        ),
      ),
      syntaxHighlighter: _CodeHighlighter(),
    );
  }
}

class _CodeHighlighter extends SyntaxHighlighter {
  @override
  TextSpan format(String source) {
    // Syntax highlighting logic
  }
}
```

**Usage:**
```dart
MarkdownMessage(
  content: '**Bold** and `code`',
  style: TextStyle(fontSize: 16),
)
```

### Example: AgentStatusIndicator

```dart
class AgentStatusIndicator extends StatelessWidget {
  final AgentStatus status;
  final String? statusText;

  const AgentStatusIndicator({
    super.key,
    required this.status,
    this.statusText,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(
          _getStatusIcon(),
          color: _getStatusColor(),
          size: 16,
        ),
        if (statusText != null) ...[
          const SizedBox(width: 4),
          Text(
            statusText!,
            style: TextStyle(color: _getStatusColor()),
          ),
        ],
      ],
    );
  }

  IconData _getStatusIcon() {
    switch (status) {
      case AgentStatus.idle: return Icons.circle;
      case AgentStatus.thinking: return Icons.psychology;
      case AgentStatus.working: return Icons.loop;
      case AgentStatus.error: return Icons.error;
    }
  }

  Color _getStatusColor() {
    switch (status) {
      case AgentStatus.idle: return Colors.grey;
      case AgentStatus.thinking: return Colors.blue;
      case AgentStatus.working: return Colors.green;
      case AgentStatus.error: return Colors.red;
    }
  }
}
```

---

## Routing with GoRouter

### Router Configuration

**File:** `lib/main.dart`

```dart
final GoRouter _router = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(
      path: '/',
      builder: (context, state) => const ChatScreen(),
    ),
    GoRoute(
      path: '/agents',
      builder: (context, state) => const AgentsScreen(),
    ),
    GoRoute(
      path: '/agents/create',
      builder: (context, state) => const AgentFormScreen(),
    ),
    GoRoute(
      path: '/agents/edit/:id',
      builder: (context, state) {
        final agentId = state.pathParameters['id']!;
        return AgentFormScreen(agentId: agentId);
      },
    ),
    // ... more routes
  ],
);
```

### Navigation Patterns

**Navigate to route:**
```dart
context.go('/agents');
```

**Navigate with parameters:**
```dart
context.go('/agents/edit/agent-123');
```

**Navigate back:**
```dart
context.pop();
```

**Navigate with result:**
```dart
final result = await context.push<bool>('/agents/create');
if (result == true) {
  // Agent was created
}
```

---

## Theme and Styling

### Material Design 3

**Theme Configuration:**

```dart
MaterialApp.router(
  theme: ThemeData(
    colorScheme: ColorScheme.fromSeed(
      seedColor: Colors.cyan,  // OneMind cyan brand color
      brightness: Brightness.dark,
    ),
    useMaterial3: true,
  ),
  routerConfig: _router,
);
```

### Cyan on Black Theme

OneMind OS uses a **cyan accent on black background** for:
- Modern, professional look
- Reduced eye strain in dark environments
- High contrast for accessibility

**Color Palette:**
- Primary: `Colors.cyan[400]` (main actions)
- Background: `Colors.black` (app background)
- Surface: `Colors.grey[900]` (cards, panels)
- Text: `Colors.white` (primary text)
- Text secondary: `Colors.grey[400]` (secondary text)

**Component Styling:**
```dart
Card(
  color: Colors.grey[900],
  elevation: 4,
  child: Padding(
    padding: const EdgeInsets.all(16),
    child: Column(
      children: [
        Text(
          'Title',
          style: TextStyle(
            color: Colors.cyan[400],
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Description',
          style: TextStyle(color: Colors.grey[400]),
        ),
      ],
    ),
  ),
)
```

---

## WebSocket Integration

### Real-time Event Handling

**Pattern 1: StreamProvider**

```dart
final websocketEventsProvider = StreamProvider<WebSocketEvent>((ref) {
  final wsService = WebSocketService();
  wsService.connect(topic: WebSocketTopic.all);
  return wsService.eventStream;
});
```

**Pattern 2: Consumer with StreamProvider**

```dart
Consumer(
  builder: (context, ref, child) {
    final eventsAsync = ref.watch(websocketEventsProvider);

    eventsAsync.whenData((event) {
      // React to events
      if (event.eventType == WebSocketEventType.sessionCreated) {
        // Invalidate sessions provider to refresh UI
        ref.invalidate(sessionsProvider);
      }
    });

    return YourWidget();
  },
)
```

**Pattern 3: Side Effects**

```dart
ref.listen<AsyncValue<WebSocketEvent>>(websocketEventsProvider, (previous, next) {
  next.whenData((event) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Event: ${event.eventType}')),
    );
  });
});
```

---

## Error Handling Patterns

### Pattern 1: AsyncValue.when

```dart
asyncValue.when(
  data: (data) => SuccessWidget(data),
  loading: () => LoadingWidget(),
  error: (err, stack) => ErrorWidget(err),
);
```

### Pattern 2: Try-Catch in Mutations

```dart
try {
  await mutations.createAgent(agent);

  if (context.mounted) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Agent created')),
    );
    context.pop();
  }
} catch (e) {
  if (context.mounted) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
    );
  }
}
```

### Pattern 3: Error Boundaries

```dart
ErrorBoundary(
  child: YourWidget(),
  onError: (error, stackTrace) {
    // Log to error tracking service
    print('Error: $error');
  },
);
```

---

## Loading States

### Pattern 1: CircularProgressIndicator

```dart
asyncValue.when(
  data: (data) => DataWidget(data),
  loading: () => const Center(
    child: CircularProgressIndicator(),
  ),
  error: (err, stack) => ErrorWidget(err),
);
```

### Pattern 2: Skeleton Loaders

```dart
loading: () => ListView.builder(
  itemCount: 5,
  itemBuilder: (context, index) => ShimmerTile(),
),
```

### Pattern 3: Pull-to-Refresh

```dart
RefreshIndicator(
  onRefresh: () async {
    ref.refresh(agentsProvider);
  },
  child: ListView(...),
)
```

---

## Form Validation

### Pattern 1: GlobalKey with Form

```dart
class AgentFormScreen extends ConsumerStatefulWidget {
  @override
  ConsumerState<AgentFormScreen> createState() => _AgentFormScreenState();
}

class _AgentFormScreenState extends ConsumerState<AgentFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        children: [
          TextFormField(
            controller: _nameController,
            decoration: const InputDecoration(labelText: 'Name'),
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter a name';
              }
              if (value.length < 3) {
                return 'Name must be at least 3 characters';
              }
              return null;
            },
          ),
          ElevatedButton(
            onPressed: _submitForm,
            child: const Text('Save'),
          ),
        ],
      ),
    );
  }

  void _submitForm() {
    if (_formKey.currentState!.validate()) {
      // Form is valid, proceed
      final agent = AgentModel(
        name: _nameController.text,
        // ...
      );

      final mutations = ref.read(agentMutationsProvider);
      mutations.createAgent(agent);
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    super.dispose();
  }
}
```

### Pattern 2: Real-time Validation

```dart
TextFormField(
  autovalidateMode: AutovalidateMode.onUserInteraction,
  validator: (value) {
    if (value == null || value.isEmpty) {
      return 'Required';
    }
    return null;
  },
)
```

---

## Best Practices

### 1. Provider Best Practices

- Use `ref.watch` to rebuild on changes
- Use `ref.read` for one-time reads (e.g., in callbacks)
- Use `ref.listen` for side effects (toasts, navigation)
- Invalidate providers after mutations: `ref.invalidate(provider)`
- Use `.family` for parameterized providers

### 2. Widget Best Practices

- Use `const` constructors whenever possible
- Extract helper methods for complex build logic
- Keep `build()` method simple and readable
- Use `ConsumerWidget` instead of `StatelessWidget` when watching providers

### 3. Performance Best Practices

- Avoid rebuilding entire screens - use `Consumer` for granular updates
- Use `ListView.builder` for long lists (lazy loading)
- Implement pagination for large datasets
- Cache images and assets
- Use `RepaintBoundary` for expensive widgets

### 4. Code Organization

- One widget per file (exceptions: small private widgets)
- Group related screens in subdirectories
- Keep models, services, providers separate
- Use barrel files (`index.dart`) for convenient imports

### 5. Testing Recommendations

- Unit test models (JSON serialization)
- Unit test services (API calls with mocks)
- Widget test screens (UI logic)
- Integration test critical flows (create agent, chat)
- Use `ProviderContainer` for testing providers

---

## Quick Reference

### Common Patterns Cheat Sheet

**Fetch data once:**
```dart
final dataProvider = FutureProvider((ref) async => ApiService.getData());
```

**Fetch data with parameter:**
```dart
final dataProvider = FutureProvider.family((ref, id) async => ApiService.getData(id));
```

**Mutation with cache invalidation:**
```dart
await ApiService.mutateData();
ref.invalidate(dataProvider);
```

**Navigate programmatically:**
```dart
context.go('/path');
context.pop();
```

**Show snackbar:**
```dart
ScaffoldMessenger.of(context).showSnackBar(
  SnackBar(content: Text('Message')),
);
```

**Handle async data:**
```dart
asyncValue.when(
  data: (data) => Widget(data),
  loading: () => CircularProgressIndicator(),
  error: (err, stack) => Text('Error: $err'),
);
```

---

**Last Updated:** 2026-02-06
**Version:** 2.0
**Status:** Production-ready

For backend architecture, see `DATABASE_SCHEMA.md` and `API_DOCUMENTATION.md`.
For deployment, see `DEPLOYMENT_GUIDE.md`.


---

## Related

- [[00-24 SO (Sovereign Operations)/05 Development (Code & Testing) (The Forge)/05 Development|05 Development]]
