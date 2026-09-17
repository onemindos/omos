# Cross-Device Real-Time Sync (Flutter + FastAPI)

Real-time synchronization system enabling instant updates across all connected Flutter devices (web, iOS, Android).

## Overview

When a user creates, updates, or deletes data on one device, all other connected devices receive the update immediately via WebSocket.

```
┌─────────────────────────────────────────────────────────────────┐
│                    BACKEND (FastAPI)                            │
│                                                                 │
│  ┌─────────────┐      ┌─────────────────────┐                   │
│  │ Session API │ ───► │ SyncConnectionManager│                  │
│  │ (CRUD ops)  │      │ broadcasts to all   │                   │
│  └─────────────┘      │ connected devices   │                   │
│                       └──────────┬──────────┘                   │
│                                  │                               │
│  WebSocket: /api/sync/ws?user_id={userId}                       │
└──────────────────────────────────┼──────────────────────────────┘
                                   │
           ┌───────────────────────┼───────────────────────┐
           │                       │                       │
   ┌───────▼───────┐       ┌───────▼───────┐       ┌───────▼───────┐
   │ Flutter Phone │       │ Flutter Web   │       │ Flutter iPad  │
   │               │       │               │       │               │
   │ SyncService   │       │ SyncService   │       │ SyncService   │
   │ ↓             │       │ ↓             │       │ ↓             │
   │ SyncListener  │       │ SyncListener  │       │ SyncListener  │
   │ ↓             │       │ ↓             │       │ ↓             │
   │ Invalidates   │       │ Invalidates   │       │ Invalidates   │
   │ providers     │       │ providers     │       │ providers     │
   └───────────────┘       └───────────────┘       └───────────────┘
```

## Event Types

| Event | Trigger | Flutter Action |
|-------|---------|----------------|
| `session.created` | New chat started | Invalidates `sessionsProvider` |
| `session.updated` | Chat renamed | Invalidates `sessionsProvider` |
| `session.deleted` | Chat deleted | Invalidates `sessionsProvider` |
| `memory.created` | Memory saved | Invalidates `memoriesProvider` |
| `memory.updated` | Memory edited | Invalidates `memoriesProvider` |
| `memory.deleted` | Memory removed | Invalidates `memoriesProvider` |

## Backend Components

### WebSocket Endpoint

**File:** `backend/agno/routing/sync_ws.py`

**Endpoint:** `ws://api.onemindos.dev/api/sync/ws?user_id={userId}`

```python
@router.websocket("/api/sync/ws")
async def sync_websocket(
    websocket: WebSocket,
    user_id: str = Query(default="zeus"),
):
    await sync_manager.connect(websocket, user_id)
    try:
        while True:
            data = await websocket.receive_text()
            if data == "ping":
                await websocket.send_text("pong")
    except WebSocketDisconnect:
        sync_manager.disconnect(websocket, user_id)
```

### SyncConnectionManager

Tracks WebSocket connections per user and broadcasts events.

```python
from backend.agno.routing import get_sync_manager

# Get the global sync manager
sync_manager = get_sync_manager()

# Broadcast events
await sync_manager.broadcast_session_created(user_id, session_id)
await sync_manager.broadcast_session_updated(user_id, session_id)
await sync_manager.broadcast_session_deleted(user_id, session_id)
```

### Notify Endpoints

**File:** `backend/agno/sessions/api.py`

For operations where direct integration isn't possible (e.g., Agno framework sessions), Flutter calls notify endpoints after mutations:

| Endpoint | Purpose |
|----------|---------|
| `POST /api/sessions/notify/created` | Broadcast session.created |
| `POST /api/sessions/notify/updated` | Broadcast session.updated |
| `POST /api/sessions/notify/deleted` | Broadcast session.deleted |
| `POST /api/sessions/notify/bulk-deleted` | Broadcast multiple deletions |

**Request:**
```json
{
  "session_id": "abc123"
}
```

**Query Parameter:** `user_id=zeus`

### Status Endpoint

```
GET /api/sync/status
```

Returns:
```json
{
  "status": "healthy",
  "total_connections": 3,
  "users_connected": 1
}
```

## Flutter Implementation

### Dependencies

Add to `pubspec.yaml`:
```yaml
dependencies:
  web_socket_channel: ^2.4.0
  flutter_riverpod: ^2.4.0
```

### File Structure

```
frontend/lib/platform/
├── services/
│   └── sync_service.dart      # WebSocket client
├── providers/
│   ├── sync_listener_provider.dart  # Event → Provider bridge
│   └── app_providers.dart     # Data providers (sessions, memories)
├── api/
│   └── agents/
│       └── session_client.dart  # API client with sync notifications
└── config/
    └── flavors.dart           # Environment-specific URLs
```

### 1. SyncService (WebSocket Client)

**File:** `frontend/lib/platform/services/sync_service.dart`

```dart
import 'dart:async';
import 'dart:convert';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../config/flavors.dart';
import '../utils/app_logger.dart';

/// Event types from the sync WebSocket
enum SyncEventType {
  sessionCreated,
  sessionUpdated,
  sessionDeleted,
  memoryCreated,
  memoryUpdated,
  memoryDeleted,
  unknown,
}

/// Represents a sync event from the server
class SyncEvent {
  final SyncEventType type;
  final String? sessionId;
  final String? memoryId;
  final DateTime timestamp;

  SyncEvent({
    required this.type,
    this.sessionId,
    this.memoryId,
    DateTime? timestamp,
  }) : timestamp = timestamp ?? DateTime.now();

  factory SyncEvent.fromJson(Map<String, dynamic> json) {
    return SyncEvent(
      type: _parseType(json['type']),
      sessionId: json['session_id'],
      memoryId: json['memory_id'],
      timestamp: json['timestamp'] != null
          ? DateTime.tryParse(json['timestamp'].toString())
          : null,
    );
  }

  static SyncEventType _parseType(String? type) {
    switch (type) {
      case 'session.created': return SyncEventType.sessionCreated;
      case 'session.updated': return SyncEventType.sessionUpdated;
      case 'session.deleted': return SyncEventType.sessionDeleted;
      case 'memory.created': return SyncEventType.memoryCreated;
      case 'memory.updated': return SyncEventType.memoryUpdated;
      case 'memory.deleted': return SyncEventType.memoryDeleted;
      default: return SyncEventType.unknown;
    }
  }
}

/// Connection states
enum SyncConnectionState {
  disconnected,
  connecting,
  connected,
  reconnecting,
}

/// WebSocket service for real-time sync
class SyncService {
  WebSocketChannel? _channel;
  final _eventController = StreamController<SyncEvent>.broadcast();
  final _stateController = StreamController<SyncConnectionState>.broadcast();

  Timer? _pingTimer;
  Timer? _reconnectTimer;
  String? _userId;
  int _reconnectAttempts = 0;
  SyncConnectionState _connectionState = SyncConnectionState.disconnected;

  static const int _maxReconnectAttempts = 10;
  static const Duration _pingInterval = Duration(seconds: 30);
  static const Duration _baseReconnectDelay = Duration(seconds: 2);

  Stream<SyncEvent> get events => _eventController.stream;
  Stream<SyncConnectionState> get connectionState => _stateController.stream;
  SyncConnectionState get currentState => _connectionState;
  bool get isConnected => _connectionState == SyncConnectionState.connected;

  /// Connect to the sync WebSocket
  Future<void> connect(String userId) async {
    if (_connectionState == SyncConnectionState.connected ||
        _connectionState == SyncConnectionState.connecting) {
      return;
    }

    _userId = userId;
    _setConnectionState(SyncConnectionState.connecting);

    final wsUrl = _getWsUrl();
    AppLogger.i('[SyncService] Connecting to $wsUrl');

    try {
      _channel = WebSocketChannel.connect(
        Uri.parse('$wsUrl?user_id=$userId'),
      );
      await _channel!.ready;

      _setConnectionState(SyncConnectionState.connected);
      _reconnectAttempts = 0;
      AppLogger.i('[SyncService] Connected successfully');

      _startPingTimer();
      _channel!.stream.listen(
        _handleMessage,
        onError: _handleError,
        onDone: _handleDisconnect,
        cancelOnError: false,
      );
    } catch (e) {
      AppLogger.e('[SyncService] Connection failed', e);
      _setConnectionState(SyncConnectionState.disconnected);
      _scheduleReconnect();
    }
  }

  String _getWsUrl() {
    if (!FlavorConfig.isInitialized) {
      return 'wss://api.onemindos.dev/api/sync/ws';
    }
    final baseUrl = FlavorConfig.instance.apiUrl;
    final wsUrl = baseUrl
        .replaceFirst('http://', 'ws://')
        .replaceFirst('https://', 'wss://');
    return '$wsUrl/api/sync/ws';
  }

  void _handleMessage(dynamic message) {
    try {
      if (message == 'pong') return;
      final data = jsonDecode(message.toString());
      if (data is Map<String, dynamic>) {
        final event = SyncEvent.fromJson(data);
        AppLogger.d('[SyncService] Received event: $event');
        _eventController.add(event);
      }
    } catch (e) {
      AppLogger.w('[SyncService] Failed to parse message: $e');
    }
  }

  void _handleError(dynamic error) {
    AppLogger.e('[SyncService] WebSocket error', error);
    _scheduleReconnect();
  }

  void _handleDisconnect() {
    AppLogger.w('[SyncService] WebSocket disconnected');
    _setConnectionState(SyncConnectionState.disconnected);
    _pingTimer?.cancel();
    _scheduleReconnect();
  }

  void _startPingTimer() {
    _pingTimer?.cancel();
    _pingTimer = Timer.periodic(_pingInterval, (_) {
      if (_connectionState == SyncConnectionState.connected) {
        try {
          _channel?.sink.add('ping');
        } catch (e) {
          AppLogger.w('[SyncService] Failed to send ping');
        }
      }
    });
  }

  void _scheduleReconnect() {
    if (_reconnectAttempts >= _maxReconnectAttempts || _userId == null) return;

    _reconnectTimer?.cancel();
    _reconnectAttempts++;

    final delay = _baseReconnectDelay * _reconnectAttempts;
    AppLogger.i('[SyncService] Reconnecting in ${delay.inSeconds}s (attempt $_reconnectAttempts)');

    _setConnectionState(SyncConnectionState.reconnecting);
    _reconnectTimer = Timer(delay, () => connect(_userId!));
  }

  void _setConnectionState(SyncConnectionState state) {
    _connectionState = state;
    _stateController.add(state);
  }

  void disconnect() {
    _pingTimer?.cancel();
    _reconnectTimer?.cancel();
    _channel?.sink.close();
    _channel = null;
    _setConnectionState(SyncConnectionState.disconnected);
  }

  void dispose() {
    disconnect();
    _eventController.close();
    _stateController.close();
  }
}

// Riverpod providers
final syncServiceProvider = Provider<SyncService>((ref) {
  final service = SyncService();
  ref.onDispose(service.dispose);
  return service;
});

final syncEventsProvider = StreamProvider<SyncEvent>((ref) {
  return ref.watch(syncServiceProvider).events;
});

final syncConnectionStateProvider = StreamProvider<SyncConnectionState>((ref) {
  return ref.watch(syncServiceProvider).connectionState;
});
```

### 2. SyncListenerProvider (Event Handler)

**File:** `frontend/lib/platform/providers/sync_listener_provider.dart`

```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../services/sync_service.dart';
import '../utils/app_logger.dart';
import 'app_providers.dart';

/// Listens to sync events and invalidates data providers
final syncListenerProvider = Provider<void>((ref) {
  final userId = ref.watch(userIdProvider);
  final syncService = ref.watch(syncServiceProvider);

  // Auto-connect when user ID is available
  syncService.connect(userId);

  // Listen for sync events and invalidate relevant providers
  ref.listen<AsyncValue<SyncEvent>>(
    syncEventsProvider,
    (previous, next) {
      next.whenData((event) {
        AppLogger.i('[SyncListener] Received event: ${event.type}');

        switch (event.type) {
          case SyncEventType.sessionCreated:
          case SyncEventType.sessionUpdated:
          case SyncEventType.sessionDeleted:
            AppLogger.d('[SyncListener] Invalidating sessionsProvider');
            ref.invalidate(sessionsProvider);
            break;

          case SyncEventType.memoryCreated:
          case SyncEventType.memoryUpdated:
          case SyncEventType.memoryDeleted:
            AppLogger.d('[SyncListener] Invalidating memoriesProvider');
            ref.invalidate(memoriesProvider);
            break;

          case SyncEventType.unknown:
            AppLogger.w('[SyncListener] Unknown event type');
            break;
        }
      });
    },
  );
});

/// Exposes sync connection status for UI
final syncConnectionStatusProvider = Provider<bool>((ref) {
  final state = ref.watch(syncConnectionStateProvider);
  return state.whenOrNull(
        data: (s) => s == SyncConnectionState.connected,
      ) ??
      false;
});
```

### 3. SessionClient (API Client with Sync)

**File:** `frontend/lib/platform/api/agents/session_client.dart`

```dart
import '../base/api_client_base.dart';
import '../../../shared/models/session.dart';

class SessionClient extends ApiClientBase {
  SessionClient({required super.baseUrl});

  String userId = 'zeus';

  /// Notify sync service about a session event
  Future<void> _notifySync(String eventType, String sessionId) async {
    try {
      await dio.post(
        '/api/sessions/notify/$eventType',
        data: {'session_id': sessionId},
        queryParameters: {'user_id': userId},
      );
      logger.d('Sync notified: $eventType for session $sessionId');
    } catch (e) {
      // Don't fail the operation if sync notification fails
      logger.w('Failed to notify sync service: $e');
    }
  }

  Future<Session> createSession({
    String? sessionName,
    String agentId = 'legacy',
  }) async {
    try {
      final response = await dio.post(
        '/sessions',
        data: {
          'agent_id': agentId,
          if (sessionName != null) 'session_name': sessionName,
          'user_id': userId,
        },
        queryParameters: {'type': 'agent'},
      );
      final session = Session.fromJson(response.data);
      logger.i('Session created: ${session.sessionId}');

      // Notify sync service for cross-device update
      await _notifySync('created', session.sessionId);

      return session;
    } catch (e, stackTrace) {
      logger.e('Failed to create session: $e', stackTrace: stackTrace);
      rethrow;
    }
  }

  Future<void> deleteSession(String sessionId) async {
    try {
      await dio.delete('/sessions/$sessionId');
      logger.i('Session deleted: $sessionId');

      // Notify sync service for cross-device update
      await _notifySync('deleted', sessionId);
    } catch (e, stackTrace) {
      logger.e('Failed to delete session: $e', stackTrace: stackTrace);
      rethrow;
    }
  }

  Future<Session> renameSession(String sessionId, String newName) async {
    try {
      final response = await dio.put(
        '/sessions/$sessionId/rename',
        data: {'session_name': newName},
      );
      logger.i('Session renamed: $sessionId to $newName');

      // Notify sync service for cross-device update
      await _notifySync('updated', sessionId);

      return Session.fromJson(response.data);
    } catch (e, stackTrace) {
      logger.e('Failed to rename session: $e', stackTrace: stackTrace);
      rethrow;
    }
  }
}
```

### 4. App Initialization

**File:** `frontend/lib/main.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'platform/providers/sync_listener_provider.dart';

void main() {
  runApp(const ProviderScope(child: MyApp()));
}

class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Initialize sync listener (auto-connects WebSocket)
    ref.read(syncListenerProvider);

    return MaterialApp(
      // ... your app
    );
  }
}
```

### 5. UI Connection Indicator (Optional)

```dart
class SyncStatusIndicator extends ConsumerWidget {
  const SyncStatusIndicator({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isConnected = ref.watch(syncConnectionStatusProvider);

    return Container(
      width: 8,
      height: 8,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        color: isConnected ? Colors.green : Colors.red,
      ),
    );
  }
}
```

### Usage Examples

#### Listen for specific events in a widget:

```dart
class SessionsScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Sessions auto-refresh when sync events arrive
    final sessions = ref.watch(sessionsProvider);

    return sessions.when(
      data: (list) => ListView.builder(
        itemCount: list.length,
        itemBuilder: (_, i) => SessionTile(session: list[i]),
      ),
      loading: () => const CircularProgressIndicator(),
      error: (e, _) => Text('Error: $e'),
    );
  }
}
```

#### Manual sync status check:

```dart
final syncService = ref.read(syncServiceProvider);
if (!syncService.isConnected) {
  await syncService.connect(userId);
}
```

## Registration

### Backend (agentos_setup.py)

The sync routers are registered in `get_custom_routers()`:

```python
# Real-Time Sync WebSocket (cross-device)
# Endpoint: ws://api.onemindos.dev/api/sync/ws?user_id=zeus
# Also provides: GET /api/sync/status
from backend.agno.routing import sync_router
routers.append(("sync_ws", sync_router))

# Session sync notifications (cross-device HTTP endpoints)
from backend.agno.sessions.api import router as sessions_sync_router
routers.append(("sessions_sync", sessions_sync_router))
```

### Flutter (main.dart)

```dart
import 'platform/providers/sync_listener_provider.dart';

// In your app initialization
ref.read(syncListenerProvider);
```

## Testing

### Verify Connection

Check backend logs for:
```
[Sync] User zeus connected. Total connections: 1
```

Check Flutter logs for:
```
[SyncService] Connected successfully
```

### Test Cross-Device Sync

1. Open app on two devices
2. Verify `/api/sync/status` shows `total_connections: 2`
3. Create a chat on Device A
4. Verify chat appears on Device B within ~1 second
5. Delete the chat on Device B
6. Verify it disappears from Device A

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| WebSocket won't connect | Wrong URL scheme | Ensure `ws://` for HTTP, `wss://` for HTTPS |
| Events not received | User ID mismatch | Verify same `user_id` on all devices |
| Stale data after reconnect | Provider not invalidated | Check `syncListenerProvider` is initialized |
| Connection drops frequently | Network issues | Reconnect logic handles automatically |

## Adding New Event Types

### Backend

1. Add broadcast method to `SyncConnectionManager`:
```python
async def broadcast_knowledge_updated(self, user_id: str, knowledge_id: str):
    await self.broadcast_to_user(user_id, {
        "type": "knowledge.updated",
        "knowledge_id": knowledge_id,
    })
```

### Flutter

1. Add to `SyncEventType` enum:
```dart
enum SyncEventType {
  // ...
  knowledgeUpdated,
}
```

2. Add parsing in `SyncEvent._parseType()`:
```dart
case 'knowledge.updated':
  return SyncEventType.knowledgeUpdated;
```

3. Add handling in `syncListenerProvider`:
```dart
case SyncEventType.knowledgeUpdated:
  ref.invalidate(knowledgeItemsProvider);
  break;
```

## Related Files

| Component | Location |
|-----------|----------|
| WebSocket endpoint | `backend/agno/routing/sync_ws.py` |
| Notify endpoints | `backend/agno/sessions/api.py` |
| Router registration | `backend/agno/agentos_setup.py` (in `get_custom_routers()`) |
| Flutter sync service | `frontend/lib/platform/services/sync_service.dart` |
| Flutter sync listener | `frontend/lib/platform/providers/sync_listener_provider.dart` |
| Session client | `frontend/lib/platform/api/agents/session_client.dart` |
| App initialization | `frontend/lib/main.dart` |
