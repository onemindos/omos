---
id: flutter-app-dev
title: "OneMind Flutter App — Development Guide"
type: reference
status: active
created: 2026-04-17
owner: zeus
tags: [type/reference, domain/so, project/flutter-app]
ai_summary: "Complete development guide for the OneMind private Flutter app — gateway protocol spec, API surface, voice/TTS architecture, auth, and project setup."
---

# OneMind Flutter App — Development Guide

> Private app for Zeus. Connects to all OneMind infrastructure via OpenClaw gateway.
> Covers iOS, Android, Linux desktop (future: VisionClaw robotics).

---

## 1. Architecture Overview

```
┌─────────────────────────────────────────────────┐
│             OneMind Flutter App                  │
│  iOS · Android · Linux Desktop · (Embedded)     │
├─────────────────────────────────────────────────┤
│  Flutter UI (Dart)  · Widgets · Theme · Nav     │
├ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─┤
│  [Future] Rust Core (libonemind) via dart:ffi   │
│  Gateway auth · Session state · Offline queue   │
├─────────────────────────────────────────────────┤
│  WebSocket (JSON-RPC)  │  REST/SSE  │  Audio    │
└────────────┬───────────┴─────┬──────┴────┬──────┘
             │                 │           │
             ▼                 ▼           ▼
┌─────────────────────────────────────────────────┐
│        OpenClaw Gateway (v2026.4.15+)            │
│  VPS: 64.227.8.23 · Tailscale: 100.98.25.10     │
│  Port: 18789 · Tailscale-served HTTPS            │
│  URL: wss://omos-core.tail717ba0.ts.net  │
├─────────────────────────────────────────────────┤
│  9 Agents · MCP Servers · Codex Vault R/W        │
│  TTS (4 providers) · Voice Call · Memory · Cron  │
└─────────────────────────────────────────────────┘
```

**Key constraint:** Private app, internal only. Connects exclusively to Zeus's infrastructure via Tailscale mesh.

---

## 1a. Why Flutter (Decision Rationale)

This decision was evaluated against four alternatives. Flutter wins on the combination of **one codebase, real Linux path, strong native plugin ecosystem, and Tailscale mesh compatibility**.

| Option | Pros | Cons | Verdict |
|--------|------|------|---------|
| **Flutter (Dart)** | One codebase for iOS + Android + Linux desktop. Mature native plugins for camera, mic, BLE, push. Fast iteration. Render engine is Skia — no web bridge tax. Linux desktop = future VisionClaw robotics UI. | Not native iOS/Android look-and-feel (Cupertino widgets exist but aren't pixel-perfect). Dart ecosystem smaller than JS/Swift. | **CHOSEN** — best balance of reach, performance, and future Linux path |
| **React Native** | JS ecosystem, huge community, Expo for fast iteration. | Weak Linux support (no official desktop target). Voice/vision requires heavy native bridging. JS bridge adds latency for real-time audio. No embedded/robotics path. | Rejected — no Linux, voice bridge tax |
| **Native (Swift + Kotlin)** | Best platform UX. Proven — Chowder + VisionClaw already exist in Swift/Kotlin. | Two codebases forever. No Linux desktop. Double the maintenance for every feature. | Rejected — can't maintain two codebases long-term |
| **KMP (Kotlin Multiplatform)** | Shared business logic, native UI per platform. Kotlin is expressive. | No Linux desktop path. iOS UI still needs SwiftUI layer. Young ecosystem. Worse than Flutter for cross-platform consistency. | Rejected — only if native UX > Linux/embedded flexibility |

### What About Rust?

Rust is **not a replacement for Flutter** — it's a **complement**. See Section 1c below. The plan:
- **Flutter** = UI shell + platform plugins (camera, mic, push, haptics)
- **Rust** (future) = Core runtime beneath Flutter (gateway auth, session state, WebSocket streaming, device identity, local cache, offline queueing)
- Access via `dart:ffi` + `ffigen`

### Native Plugins Still Needed

Even with Flutter, some features require native platform plugins:
- **Camera** — `camera` package (uses AVCaptureSession / Camera2 under the hood)
- **Microphone** — `record` or `flutter_sound` for PCM capture
- **Live Activities** (iOS) — Swift native code via platform channel
- **BLE** (wearables) — `flutter_blue_plus`
- **Push Notifications** — `firebase_messaging` or APNs via platform channel
- **Location** — `geolocator` for Find Zeus mode
- **Secure Storage** — `flutter_secure_storage` (Keychain / KeyStore)

---

## 1b. Existing Mobile Assets (Port from onemind-mobile)

The `onemind-mobile/` repo contains **three working native apps** that have already solved the gateway integration problems. These are the source of truth for protocol behavior — port the patterns, not just the theory.

### What Exists

| App | Language | What It Proves | Key Files |
|-----|----------|---------------|-----------|
| **Chowder iOS** | Swift | WebSocket JSON-RPC chat with reconnect, history, tool events, Live Activities | `ChatService.swift` |
| **VisionClaw iOS** | Swift | REST bridge, audio pipeline (PCM 16kHz→24kHz), video pipeline (JPEG 1fps), tool calling | `OpenClawBridge.swift` |
| **VisionClaw Android** | Kotlin | Same REST bridge on Android, reactive state via StateFlow | `OpenClawBridge.kt` |

### Chowder `ChatService.swift` → Port to Dart `GatewayClient`

This is the **most valuable file to port**. It contains the complete WebSocket protocol implementation:

| Pattern | Swift Implementation | Dart Equivalent |
|---------|---------------------|-----------------|
| **WebSocket connection** | `URLSessionWebSocketTask` with `?client=chowder-ios` query param | `WebSocketChannel.connect(uri.replace(queryParameters: {'client': 'onemind-flutter'}))` |
| **JSON-RPC framing** | Outbound: `type: "req"`, inbound: `type: "event"` / `type: "res"` | Same JSON structure, use `jsonEncode`/`jsonDecode` |
| **Request IDs** | Counter: `req-1`, `req-2`, etc. | `int _requestCounter = 0; String get _nextId => 'req-${++_requestCounter}';` |
| **Device pairing** | `UIDevice.current.identifierForVendor` sent on connect | `DeviceInfoPlugin().androidInfo.id` / `iosInfo.identifierForVendor` |
| **Session keys** | `"agent:main:main"` for chat, `"agent:main:glass"` for vision | Same format: `"agent:{agentId}:{channel}"` |
| **Auto-reconnect** | State tracking (connected/disconnected), reconnect on failure | `Timer.periodic` + exponential backoff |
| **History polling** | Deduplication via `seenSequenceNumbers`, `seenToolCallIds`, `seenTimestamps` | Same triple-dedup with `Set<int>`, `Set<String>`, `Set<String>` |
| **Server events** | `event` type with sub-events: `tick`, `health`, `agent/lifecycle`, `chat` | Dart `StreamController` per event type, or single stream with type filtering |
| **Delegate pattern** | `ChatServiceDelegate` protocol: connect, disconnect, delta, finish, error, thinking, tool events, bot identity, user profile, history | Dart: callback functions or `Stream<ChatEvent>` with sealed class |

### VisionClaw `OpenClawBridge.swift` → Port REST Patterns

| Pattern | Implementation | Notes |
|---------|---------------|-------|
| **REST endpoint** | POST `/v1/chat/completions` with Bearer token | Already in our `GatewayClient` starter code |
| **Custom headers** | `x-openclaw-message-channel: glass`, `x-openclaw-session-key` | Add to Dart `Dio` interceptor |
| **Connection state machine** | `notConfigured → checking → connected / unreachable` | Port as Dart `enum` + `ValueNotifier` |
| **Conversation history** | Last 10 turns kept in-memory, sent with each request | Port with `Queue<Message>` bounded to 10 |
| **Timeouts** | 120s request, 5s ping | Configure in Dio: `connectTimeout`, `receiveTimeout` |

### Audio Pipeline (from VisionClaw) → Reference for Voice

VisionClaw's proven audio pipeline:
- **Input**: PCM Int16, 16kHz, mono → sent via Gemini WebSocket
- **Output**: PCM 24kHz → played via `AVAudioEngine`
- **Flutter equivalent**: `record` package for PCM capture, `just_audio` for playback
- **Note**: The gateway's TTS providers (OpenAI, Google, ElevenLabs, Microsoft) accept text and return audio — simpler than Gemini's bidirectional audio stream

### Key Headers to Preserve

```dart
// These headers are proven in the native apps — use exactly
final headers = {
  'Authorization': 'Bearer $gatewayToken',
  'x-openclaw-message-channel': channel,  // 'main', 'glass', 'flutter'
  'x-openclaw-session-key': sessionKey,   // 'agent:main:main'
  'Content-Type': 'application/json',
};
```

---

## 1c. Rust Core Layer (Future — Not for MVP)

GPT's analysis recommended a Rust core layer beneath Flutter. This is a **future optimization**, not an MVP requirement. Document it now so the Flutter architecture doesn't paint us into a corner.

### What Rust Would Own

| Responsibility | Why Rust | Current Location |
|---------------|----------|-----------------|
| **Gateway auth + session state** | Zero-copy token handling, no GC pauses during WebSocket streaming | Dart `GatewayClient` class |
| **WebSocket streaming** | `tokio-tungstenite` — battle-tested, lower memory than Dart isolate | Dart `web_socket_channel` |
| **Device identity** | Stable device fingerprinting with OS keychain bindings | Platform channels |
| **Local cache** | `sled` or `redb` embedded DB — faster than SQLite for key-value patterns | SharedPreferences / Hive |
| **Offline queueing** | Guaranteed delivery with WAL — messages persist across crashes | Not yet implemented |
| **Headless mode** | Runs on Linux/ARM without a display server — for VisionClaw embedded + robotics | Not possible with Flutter alone |

### Architecture With Rust

```
┌─────────────────────────────────────────────┐
│          Flutter UI Shell (Dart)             │
│  Widgets · Navigation · Theme · Animations  │
├─────────────────────────────────────────────┤
│  dart:ffi  ←→  ffigen generated bindings    │
├─────────────────────────────────────────────┤
│          Rust Core Runtime (libonemind)      │
│  Gateway Client · Session State · Cache     │
│  Offline Queue · Device Identity · Crypto   │
├─────────────────────────────────────────────┤
│  tokio (async) · tungstenite (WS) · serde   │
│  sled (cache) · ring (crypto) · tracing     │
└─────────────────────────────────────────────┘
```

### Why Not Now

1. **Flutter's Dart WebSocket and HTTP are sufficient** for a private app with one user
2. **Adds build complexity** — cross-compile Rust for iOS (aarch64-apple-ios), Android (aarch64-linux-android), Linux (x86_64/aarch64)
3. **MVP needs speed** — Dart-only ships faster
4. **Revisit when**: Linux desktop build is real, or offline-first becomes a requirement, or robotics nodes need headless gateway access

### Keeping the Door Open

To ensure the Flutter app can adopt Rust later without rewriting:
- **Isolate gateway logic** in a single `GatewayClient` class (already in our starter code)
- **Don't scatter WebSocket handling** across widgets — keep it in a service layer
- **Use streams/callbacks** for data flow — same interface whether backed by Dart or FFI
- **Session state in a repository pattern** — swap storage backend without changing consumers

---

## 2. Gateway Connection Details

| Field | Value |
|-------|-------|
| **WebSocket URL** | `wss://omos-core.tail717ba0.ts.net` |
| **REST Base URL** | `https://omos-core.tail717ba0.ts.net` |
| **Public IP** | `64.227.8.23` |
| **Tailscale IP** | `100.98.25.10` |
| **Gateway Port** | `18789` |
| **Auth Mode** | Token-based |
| **Gateway Token** | Stored in 1Password: `op://SO Vault/OpenClaw Gateway Token/credential` |
| **OpenClaw Version** | `v2026.4.15` |
| **Tailscale Peers** | Mac `100.114.142.55` · iPhone `100.122.203.68` · VPS `100.98.25.10` |

### Authentication

All requests require the gateway token:

```dart
// WebSocket — token in URL query param
final ws = WebSocketChannel.connect(
  Uri.parse('wss://omos-core.tail717ba0.ts.net?token=$gatewayToken'),
);

// REST — token in Authorization header
final response = await dio.get(
  'https://omos-core.tail717ba0.ts.net/api/v1/ping',
  options: Options(headers: {'Authorization': 'Bearer $gatewayToken'}),
);
```

---

## 3. Protocols & Transport

The gateway speaks **three standard protocols**. No proprietary SDK required.

### 3.1 WebSocket + JSON-RPC (Primary — Real-Time)

All real-time communication uses WebSocket with JSON-RPC 2.0 messages.

**Evidence:** 413 WebSocket references, 764 `method:` references, 16 explicit `jsonrpc` references in gateway dist/.

```dart
// Send a JSON-RPC request
final request = jsonEncode({
  'jsonrpc': '2.0',
  'id': uuid.v4(),
  'method': 'chat.send',
  'params': {
    'sessionKey': sessionKey,
    'message': 'Hello Legacy',
  },
});
ws.sink.add(request);

// Listen for responses
ws.stream.listen((data) {
  final response = jsonDecode(data);
  if (response['method'] == 'chat.delta') {
    // Streaming token
  }
});
```

**Flutter package:** `web_socket_channel` (stable, maintained)

### 3.2 REST + OpenAI-Compatible (Stateless Calls)

Standard HTTP endpoints for one-shot operations.

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/v1/ping` | GET | Health check |
| `/api/v1/server/info` | GET | Server info, version |
| `/api/v1/chat/query` | POST | One-shot chat query |
| `/api/v1/message/text` | POST | Send text message |
| `/api/v1/search` | POST | Search |
| `/api/v1/download` | GET | File download |
| `/upload/image` | POST | Image upload |
| `/v1/chat/completions` | POST | OpenAI-compatible chat |
| `/v1/embeddings` | POST | OpenAI-compatible embeddings |
| `/v1/files` | POST | File management |
| `/mcp` | POST | MCP protocol endpoint |
| `/tools/invoke` | POST | Direct tool invocation |
| `/api/channels` | GET | Channel status |

**Flutter package:** `dio` (preferred) or `http`

### 3.3 SSE — Server-Sent Events (Streaming)

For streaming chat completions without full WebSocket.

**Evidence:** 6 `text/event-stream` references, 7 `EventSource` references in dist/.

```dart
// SSE streaming via dio
final response = await dio.post(
  '$baseUrl/v1/chat/completions',
  data: {'messages': [...], 'stream': true},
  options: Options(
    responseType: ResponseType.stream,
    headers: {'Authorization': 'Bearer $gatewayToken'},
  ),
);

response.data.stream
  .transform(utf8.decoder)
  .transform(const LineSplitter())
  .where((line) => line.startsWith('data: '))
  .map((line) => line.substring(6))
  .listen((data) {
    if (data == '[DONE]') return;
    final chunk = jsonDecode(data);
    // Handle streaming token
  });
```

**Flutter package:** `eventsource` or manual SSE via `dio` stream

---

## 4. Full RPC Method Surface

Every method callable via WebSocket JSON-RPC. This IS the API.

### 4.1 Chat

| Method | Purpose |
|--------|---------|
| `chat.send` | Send a message to an agent |
| `chat.subscribe` | Subscribe to chat events for a session |
| `chat.unsubscribe` | Unsubscribe from chat events |
| `chat.history` | Get chat history for a session |
| `chat.abort` | Abort current agent response |
| `chat.completion` | Get completion suggestions |
| `chat.inject` | Inject system message into context |
| `chat.voice` | Send voice message |
| `chat.photo` | Send photo |
| `chat.gif` | Send GIF |
| `chat.sticker` | Send sticker |
| `chat.doodle` | Send doodle |
| `chat.todo` | Create todo from chat |

### 4.2 Agents

| Method | Purpose |
|--------|---------|
| `agents.list` | List all configured agents |
| `agents.create` | Create a new agent |
| `agents.update` | Update agent config |
| `agents.delete` | Delete an agent |
| `agents.defaults` | Get/set agent defaults |

### 4.3 Sessions

| Method | Purpose |
|--------|---------|
| `sessions.create` | Create a new chat session |
| `sessions.get` | Get session details |
| `sessions.delete` | Delete a session |
| `sessions.abort` | Abort session execution |
| `sessions.compact` | Compact session context |
| `sessions.changed` | Check for session changes |

### 4.4 Models

| Method | Purpose |
|--------|---------|
| `models.list` | List available models |
| `models.authStatus` | Check model auth/availability |
| `models.providers` | List model providers |
| `models.mode` | Get/set model mode |
| `models.bedrockDiscovery` | AWS Bedrock model discovery |
| `models.legacy` | Legacy model compat |

### 4.5 Memory

| Method | Purpose |
|--------|---------|
| `memory.backend` | Get memory backend info |
| `memory.citations` | Get memory citations |
| `memory.grow` | Add to agent memory |
| `memory.size` | Get memory size |

### 4.6 Config

| Method | Purpose |
|--------|---------|
| `config.get` | Get config value |
| `config.set` | Set config value |
| `config.patch` | Patch config |
| `config.write` | Write config to disk |
| `config.apply` | Apply config changes |
| `config.observe` | Watch config changes (reactive) |
| `config.schema` | Get config JSON schema |
| `config.open` | Open config file |

### 4.7 Cron / Scheduling

| Method | Purpose |
|--------|---------|
| `cron.add` | Schedule a recurring job |
| `cron.list` | List all cron jobs |
| `cron.remove` | Remove a cron job |
| `cron.run` | Manually trigger a job |
| `cron.runs` | Get job run history |
| `cron.status` | Get cron system status |
| `cron.update` | Update a cron job |
| `cron.webhook` | Cron webhook endpoint |

### 4.8 Channels (Messaging Platforms)

| Method | Purpose |
|--------|---------|
| `channels.status` | Get all channel status |
| `channels.logout` | Logout from a channel |

**Supported channels:** Telegram, Discord, WhatsApp, Slack, Matrix, Signal, SMS, email, and more.

### 4.9 Device

| Method | Purpose |
|--------|---------|
| `device.health` | Device health check |
| `device.id` | Get device ID |
| `device.info` | Get device info |
| `device.permissions` | Get/set device permissions |
| `device.status` | Get device status |

### 4.10 Plugins

| Method | Purpose |
|--------|---------|
| `plugins.allow` | Allow a plugin |
| `plugins.deny` | Deny a plugin |
| `plugins.enabled` | List enabled plugins |
| `plugins.entries` | List all plugin entries |
| `plugins.installs` | List installed plugins |
| `plugins.load` | Load a plugin |

### 4.11 Tools

| Method | Purpose |
|--------|---------|
| `tools/invoke` | Invoke a tool directly |

### 4.12 Speech / TTS

| Method | Purpose |
|--------|---------|
| Speech provider selection | Via `config.set` for `messages.tts.providers.*` |
| TTS directives | Inline `[[tts:speed=1.2]]` in messages |
| TTS auto-mode | Automatic TTS for agent replies |

---

## 5. Voice & TTS Architecture

OpenClaw 4.15 ships a full voice stack. This is what powers voice mode in the app.

### 5.1 Text-to-Speech (TTS) Providers

| Provider | Model | Output Formats | Status |
|----------|-------|----------------|--------|
| **OpenAI** | `gpt-4o-mini-tts` | WAV, MP3, PCM | Available |
| **Google/Gemini** | Gemini TTS (new in 4.15) | WAV, PCM | Available |
| **ElevenLabs** | Various voices | WAV, MP3 | Auto-enabled in 4.15 |
| **Microsoft** | Azure TTS | WAV, MP3 | Auto-enabled in 4.15 |

**Config path:** `messages.tts.providers.<provider>.apiKey`

**TTS directives** (inline in messages):
```
[[tts:speed=1.2]]     — Adjust speech speed
[[tts:voice=alloy]]   — Select voice
```

### 5.2 Real-Time Voice

| Capability | Provider Type | Purpose |
|------------|---------------|---------|
| `realtimeVoice` | Voice output provider | Real-time voice synthesis |
| `realtimeTranscription` | STT provider | Speech-to-text input |
| `voice-call` extension | Telephony | Make/receive phone calls |

**Voice call CLI** (exists in gateway):
```
voicecall call --to <phone> --message <text> --mode <mode>
```

### 5.3 Audio Formats

| Format | Use Case |
|--------|----------|
| **PCM** | Raw audio for real-time streaming |
| **WAV** | Standard playback |
| **Opus** | Compressed, low-latency voice |

### 5.4 Flutter Voice Implementation

```dart
// Request TTS for a message
final ttsRequest = jsonEncode({
  'jsonrpc': '2.0',
  'id': uuid.v4(),
  'method': 'chat.voice',
  'params': {
    'sessionKey': sessionKey,
    'messageId': messageId,
    'provider': 'openai', // or 'google', 'elevenlabs', 'microsoft'
  },
});
ws.sink.add(ttsRequest);

// Play audio response
// Audio data comes back as binary frames on WebSocket
// or as a URL to stream via REST
```

**Flutter packages for voice:**
- `just_audio` — Audio playback (WAV, MP3, streaming)
- `flutter_sound` — Recording + playback (PCM, WAV, Opus)
- `record` — Microphone capture
- `speech_to_text` — On-device STT (or use gateway's realtimeTranscription)

---

## 6. Agent Roster

These are the 9 agents available via the gateway. All use `github-copilot/claude-sonnet-4.6`.

| ID | Agent | Emoji | Domain | Scope |
|----|-------|-------|--------|-------|
| `main` | **Legacy** | ⚡ | All | Orchestrator — routes to specialists |
| `spartan` | **Spartan** | 💪 | HP 27-28 | Body, fitness, nutrition |
| `oracle` | **Oracle** | 🔮 | HP 25-26, 29-34 | Mind, money, skills |
| `eden` | **Eden** | 🌱 | LE | Sustainability, environment |
| `haven` | **Haven** | 🏡 | LE 50-74 | Family, home, community |
| `empire` | **Empire** | 💰 | GE 75-99 | Ventures, brand, wealth |
| `guardian` | **Guardian** | 🛡️ | SO 11 | Security, access control |
| `grid` | **Grid** | ⚙️ | SO 01, 09 | Infrastructure, automation |
| `trinity` | **Trinity** | 🔺 | SO 03 | Agent coordination |

**Important:** The agent ID for Legacy is `main`, not `legacy`.

```dart
// List agents
ws.sink.add(jsonEncode({
  'jsonrpc': '2.0',
  'id': uuid.v4(),
  'method': 'agents.list',
  'params': {},
}));

// Send message to specific agent
ws.sink.add(jsonEncode({
  'jsonrpc': '2.0',
  'id': uuid.v4(),
  'method': 'chat.send',
  'params': {
    'agentId': 'main', // or 'spartan', 'oracle', etc.
    'message': 'What\'s on today?',
  },
}));
```

---

## 7. MCP Servers (Available via Gateway)

Agents have access to these MCP servers. The Flutter app benefits indirectly — ask an agent and it uses MCP tools.

| Server | What It Provides |
|--------|-----------------|
| `context7` | Up-to-date library/framework docs |
| `github` | GitHub repos, issues, PRs, file contents |
| `codex-vault` | Read/write access to `/root/onemind-codex/` |

---

## 8. Flutter Project Setup

### 8.1 Create Project

```bash
# Ensure Flutter is installed
flutter --version  # Need 3.22+ for Dart 3.4+

# Create new project
flutter create onemind_app \
  --org com.onemindos \
  --description "OneMind private command center" \
  --platforms ios,android,linux

cd onemind_app
```

### 8.2 Dependencies (pubspec.yaml)

```yaml
name: onemind_app
description: OneMind private command center
version: 1.0.0+1

environment:
  sdk: '>=3.4.0 <4.0.0'
  flutter: '>=3.22.0'

dependencies:
  flutter:
    sdk: flutter

  # Networking
  web_socket_channel: ^3.0.0    # WebSocket (JSON-RPC transport)
  dio: ^5.4.0                   # REST + SSE streaming
  connectivity_plus: ^6.0.0     # Network state monitoring

  # State Management
  riverpod: ^2.5.0              # State management (or flutter_riverpod)
  flutter_riverpod: ^2.5.0

  # Audio / Voice
  just_audio: ^0.9.36           # Audio playback
  record: ^5.1.0                # Microphone recording
  flutter_sound: ^9.6.0         # Full audio toolkit (PCM, WAV, Opus)

  # Storage
  hive: ^2.2.3                  # Local key-value storage
  hive_flutter: ^1.1.0
  flutter_secure_storage: ^9.0.0 # Secure token storage

  # UI
  flutter_markdown: ^0.7.1       # Render markdown responses
  cached_network_image: ^3.3.0   # Image caching
  shimmer: ^3.0.0                # Loading states

  # Utilities
  uuid: ^4.3.0                  # Request IDs
  json_annotation: ^4.8.0       # JSON serialization
  intl: ^0.19.0                 # Date/time formatting
  url_launcher: ^6.2.0          # Open URLs
  share_plus: ^7.2.0            # Share content

dev_dependencies:
  flutter_test:
    sdk: flutter
  build_runner: ^2.4.0
  json_serializable: ^6.7.0
  flutter_lints: ^3.0.0
```

### 8.3 Project Structure

```
onemind_app/
├── lib/
│   ├── main.dart                      # App entry point
│   ├── app.dart                       # MaterialApp / theme
│   │
│   ├── config/
│   │   ├── constants.dart             # Gateway URLs, timeouts
│   │   ├── theme.dart                 # Solarpunk theme tokens
│   │   └── routes.dart                # Named routes
│   │
│   ├── core/
│   │   ├── gateway/
│   │   │   ├── gateway_client.dart    # WebSocket + JSON-RPC client
│   │   │   ├── rpc_methods.dart       # All RPC method constants
│   │   │   ├── rpc_request.dart       # JSON-RPC request builder
│   │   │   └── rpc_response.dart      # JSON-RPC response parser
│   │   ├── rest/
│   │   │   ├── rest_client.dart       # Dio-based REST client
│   │   │   └── sse_client.dart        # SSE streaming client
│   │   ├── auth/
│   │   │   ├── auth_service.dart      # Token management
│   │   │   └── secure_storage.dart    # Encrypted token storage
│   │   └── voice/
│   │       ├── tts_service.dart       # Text-to-speech
│   │       ├── stt_service.dart       # Speech-to-text
│   │       └── audio_player.dart      # Audio playback manager
│   │
│   ├── models/
│   │   ├── agent.dart                 # Agent model
│   │   ├── session.dart               # Chat session model
│   │   ├── message.dart               # Chat message model
│   │   ├── device.dart                # Device info model
│   │   └── cron_job.dart              # Cron job model
│   │
│   ├── providers/                     # Riverpod providers
│   │   ├── gateway_provider.dart      # Gateway connection state
│   │   ├── agents_provider.dart       # Agent list
│   │   ├── sessions_provider.dart     # Active sessions
│   │   ├── chat_provider.dart         # Chat messages
│   │   └── voice_provider.dart        # Voice mode state
│   │
│   ├── screens/
│   │   ├── home/
│   │   │   └── home_screen.dart       # Dashboard / command center
│   │   ├── chat/
│   │   │   ├── chat_screen.dart       # Chat with agent
│   │   │   └── chat_input.dart        # Text + voice input
│   │   ├── agents/
│   │   │   └── agents_screen.dart     # Agent roster / status
│   │   ├── voice/
│   │   │   └── voice_screen.dart      # Full voice mode
│   │   ├── memory/
│   │   │   └── memory_screen.dart     # Memory search / browse
│   │   ├── cron/
│   │   │   └── cron_screen.dart       # Scheduled jobs
│   │   └── settings/
│   │       └── settings_screen.dart   # Config, connection, TTS
│   │
│   └── widgets/
│       ├── agent_avatar.dart          # Agent emoji + status
│       ├── message_bubble.dart        # Chat bubble (markdown)
│       ├── voice_button.dart          # Push-to-talk / voice toggle
│       ├── connection_indicator.dart   # Gateway connection status
│       └── domain_badge.dart          # SO/HP/LE/GE color badge
│
├── assets/
│   ├── fonts/                         # Custom fonts
│   └── images/                        # Static assets
│
├── ios/
├── android/
├── linux/
└── test/
```

### 8.4 Core Gateway Client (Starter Code)

```dart
// lib/core/gateway/gateway_client.dart

import 'dart:async';
import 'dart:convert';
import 'package:uuid/uuid.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class GatewayClient {
  final String url;
  final String token;
  final _uuid = const Uuid();

  WebSocketChannel? _channel;
  final _responseControllers = <String, Completer<Map<String, dynamic>>>{};
  final _eventController = StreamController<Map<String, dynamic>>.broadcast();

  Stream<Map<String, dynamic>> get events => _eventController.stream;
  bool get isConnected => _channel != null;

  GatewayClient({required this.url, required this.token});

  Future<void> connect() async {
    final uri = Uri.parse('$url?token=$token');
    _channel = WebSocketChannel.connect(uri);

    _channel!.stream.listen(
      (data) => _handleMessage(jsonDecode(data)),
      onError: (error) => _handleError(error),
      onDone: () => _handleDisconnect(),
    );
  }

  void _handleMessage(Map<String, dynamic> message) {
    final id = message['id'] as String?;

    if (id != null && _responseControllers.containsKey(id)) {
      // Response to a request we sent
      _responseControllers[id]!.complete(message);
      _responseControllers.remove(id);
    } else {
      // Server-pushed event (chat.delta, session updates, etc.)
      _eventController.add(message);
    }
  }

  void _handleError(dynamic error) {
    // Reconnect logic here
  }

  void _handleDisconnect() {
    _channel = null;
    // Auto-reconnect with exponential backoff
  }

  /// Send a JSON-RPC request and await the response
  Future<Map<String, dynamic>> call(
    String method, [
    Map<String, dynamic>? params,
  ]) async {
    final id = _uuid.v4();
    final completer = Completer<Map<String, dynamic>>();
    _responseControllers[id] = completer;

    final request = jsonEncode({
      'jsonrpc': '2.0',
      'id': id,
      'method': method,
      if (params != null) 'params': params,
    });

    _channel!.sink.add(request);
    return completer.future.timeout(const Duration(seconds: 30));
  }

  /// Fire-and-forget (no response expected)
  void notify(String method, [Map<String, dynamic>? params]) {
    final request = jsonEncode({
      'jsonrpc': '2.0',
      'method': method,
      if (params != null) 'params': params,
    });
    _channel!.sink.add(request);
  }

  void disconnect() {
    _channel?.sink.close();
    _channel = null;
  }
}
```

### 8.5 Usage Example

```dart
// Connect to gateway
final gateway = GatewayClient(
  url: 'wss://omos-core.tail717ba0.ts.net',
  token: secureStorage.read('gateway_token'),
);
await gateway.connect();

// List agents
final agents = await gateway.call('agents.list');

// Create a chat session
final session = await gateway.call('sessions.create', {
  'agentId': 'main',
});

// Subscribe to chat events
gateway.notify('chat.subscribe', {
  'sessionKey': session['result']['sessionKey'],
});

// Send a message
await gateway.call('chat.send', {
  'sessionKey': session['result']['sessionKey'],
  'message': 'Good morning Legacy. What\'s on my plate today?',
});

// Listen for streaming responses
gateway.events
  .where((e) => e['method'] == 'chat.delta')
  .listen((delta) {
    // Update UI with streaming text
    print(delta['params']['content']);
  });
```

---

## 9. Solarpunk Theme Tokens

Match the OneMind web UI. CSS tokens from `onemind-web/src/app/globals.css`, mapped to Flutter.

```dart
// lib/config/theme.dart

import 'package:flutter/material.dart';

class OneMindColors {
  // Solarpunk Primary
  static const moss = Color(0xFF2D5016);
  static const amber = Color(0xFFF59E0B);
  static const sky = Color(0xFF0EA5E9);
  static const coral = Color(0xFFF97316);

  // Domain Colors
  static const so = Color(0xFF6366F1);   // Sovereign Operations — indigo
  static const hp = Color(0xFF22C55E);   // Holistic Performance — green
  static const le = Color(0xFF3B82F6);   // Legacy Evolution — blue
  static const ge = Color(0xFFF59E0B);   // Generational Empire — amber

  // Tactical Dark (Zeus Red + Cyan)
  static const zeusRed = Color(0xFFDC2626);
  static const zeusCyan = Color(0xFF06B6D4);

  // Backgrounds
  static const bgDark = Color(0xFF0A0A0A);
  static const bgCard = Color(0xFF1A1A1A);
  static const bgSurface = Color(0xFF2A2A2A);

  // Text
  static const textPrimary = Color(0xFFF5F5F5);
  static const textSecondary = Color(0xFFA3A3A3);
  static const textMuted = Color(0xFF737373);
}

// Agent avatar colors
const agentColors = {
  'main': OneMindColors.amber,      // ⚡ Legacy
  'spartan': OneMindColors.hp,      // 💪 Spartan
  'oracle': Color(0xFF8B5CF6),      // 🔮 Oracle — purple
  'eden': Color(0xFF10B981),        // 🌱 Eden — emerald
  'haven': OneMindColors.le,        // 🏡 Haven
  'empire': OneMindColors.ge,       // 💰 Empire
  'guardian': OneMindColors.zeusRed, // 🛡️ Guardian
  'grid': OneMindColors.zeusCyan,   // ⚙️ Grid
  'trinity': OneMindColors.so,      // 🔺 Trinity
};
```

---

## 10. App Features (MVP → v1.0)

### MVP (Phase 1)

- [ ] Gateway connection + auth (token from secure storage)
- [ ] Connection status indicator (connected / reconnecting / offline)
- [ ] Agent list with status
- [ ] Chat with Legacy (main agent) — text only
- [ ] Streaming response rendering (markdown)
- [ ] Chat history per session
- [ ] Session management (create, switch, delete)
- [ ] Basic settings (gateway URL, token)

### Phase 2 — Voice

- [ ] TTS playback on agent responses
- [ ] Push-to-talk voice input
- [ ] Voice mode screen (full-screen, minimal UI)
- [ ] TTS provider selection (OpenAI / Google / ElevenLabs / Microsoft)
- [ ] Audio streaming (PCM/WAV)

### Phase 3 — Full Command Center

- [ ] Multi-agent chat (switch between all 9 agents)
- [ ] Domain-colored UI (SO/HP/LE/GE badges)
- [ ] Memory search
- [ ] Cron job viewer + manual trigger
- [ ] Photo/image send in chat
- [ ] Device health monitoring
- [ ] Config viewer (read-only)

### Phase 4 — Integration

- [ ] Todoist widget / quick-add
- [ ] Home Assistant quick controls (lights, sensors)
- [ ] Notification routing (agent alerts → push notification)
- [ ] Vault file browser (via codex-vault MCP)
- [ ] Location sharing for Find Zeus mode
- [ ] Linux desktop build for VisionClaw robotics

---

## 10a. Page Architecture (Full Screen Map)

The Flutter app mirrors and extends two existing dashboards — the OpenClaw native dashboard (19 tabs) and onemind-web (25 views) — while adding mobile-native screens for the Codex vault.

### OpenClaw Native Dashboard (Already Exists — Built Into Gateway)

| Section | Tabs |
|---------|------|
| **Chat** | chat |
| **Control** | overview, channels, instances, sessions, usage, cron |
| **Agent** | agents, skills, nodes, dreams |
| **Settings** | config, communications, appearance, automation, infrastructure, aiAgents, debug, logs |

### onemind-web (Already Exists — Next.js 16)

| Group | Pages |
|-------|-------|
| **Workspace** | Home, Agents, Org Chart, Office, Inbox, Chatrooms, Sessions, Projects |
| **Execution** | Tasks, Missions, Schedules, Memory, Runs |
| **Knowledge** | Knowledge, Skills, Connectors, Webhooks, MCP Servers, Extensions |
| **System** | Secrets, Providers, Usage, Activity, Autonomy, Wallets, Logs, Settings |

### Flutter App — Bottom Tab Bar (5 Tabs)

```
┌──────┬──────────┬──────────┬──────────┬──────────┐
│  ⚡  │    💬    │    🧠    │    📊    │    ⚙️    │
│ Home │   Chat   │  Agents  │  Codex   │ Command  │
└──────┴──────────┴──────────┴──────────┴──────────┘
```

### Tab 1: ⚡ Home (Dashboard)

| Page | Source | What It Shows |
|------|--------|---------------|
| Dashboard | New | Cards: agent status, active chats, today's tasks, quick actions |
| Activity Feed | onemind-web `activity` | Real-time event stream from all agents |
| Morning Briefing | New (skill) | Pull from morning-briefing skill |
| Notifications | New | Push notification history + routing |

### Tab 2: 💬 Chat (Agent Conversations)

| Page | Source | What It Shows |
|------|--------|---------------|
| Chat List | OpenClaw `chat` + onemind-web `chatrooms` | All active sessions, grouped by agent |
| Chat Thread | OpenClaw `chat` | Full conversation with streaming, tool calls, thinking |
| Voice Chat | New (VisionClaw port) | Real-time voice conversation (TTS + mic) |
| Quick Capture | New (skill) | "Hey Legacy, save this..." — one-tap capture to vault |
| Group Chat | onemind-web `chatrooms` | Multi-agent conversation |

### Tab 3: 🧠 Agents (AI Workforce)

| Page | Source | What It Shows |
|------|--------|---------------|
| Agent Roster | OpenClaw `agents` + onemind-web `agents` | All 9 agents with status, domain badge, model |
| Agent Detail | onemind-web `agents/[id]` | Config, memory, tools, active sessions |
| Org Chart | onemind-web `org-chart` | Visual agent hierarchy |
| Skills | OpenClaw `skills` | Installed skills, ClewHub search, enable/disable |
| Memory | onemind-web `memory` + OpenClaw `dreams` | Memory entries, dream diary, dreaming status |
| MCP Servers | onemind-web `mcp-servers` | Connected MCP servers |
| Runs | onemind-web `runs` | Execution traces and logs |
| Autonomy | onemind-web `autonomy` | Per-agent autonomy levels (L0-L3) |

### Tab 4: 📊 Codex (Your Life OS) — The Killer Differentiator

This is what makes the Flutter app more than a dashboard. It reads the vault via MCP and presents the entire life OS as navigable, domain-colored screens.

| Page | Source | What It Shows |
|------|--------|---------------|
| Domain Map | New | Visual 3+1 grid — SO/HP/LE/GE with health scores |
| Domain Detail | New | Drill into any domain — goals, projects, metrics |
| Vault Browser | New (via codex-vault MCP) | Browse vault files, read notes, search |
| Tasks | onemind-web `tasks` + Todoist skill | Today's tasks, grouped by domain |
| Projects | onemind-web `projects` | Active projects with status |
| Goals | New | Goals from vault frontmatter, progress tracking |
| Reviews | New (skill) | Weekly scorecard, evening reviews, life scores |
| Metrics | New | PI/WI/LI intelligence scores, trends |
| Calendar | New | Schedules from onemind-web + Todoist |
| Contacts | New (CRM skill) | Relationship tracker from LE domain |

### Tab 5: ⚙️ Command (System Control)

| Page | Source | What It Shows |
|------|--------|---------------|
| System Overview | OpenClaw `overview` | Gateway status, uptime, connected devices |
| Sessions | OpenClaw `sessions` | Active sessions with context weight |
| Usage | OpenClaw `usage` + onemind-web `usage` | Token usage, cost tracking, time series |
| Channels | OpenClaw `channels` | Telegram, Discord, connected channels |
| Instances | OpenClaw `instances` | Running instances/nodes |
| Cron | OpenClaw `cron` | Scheduled jobs (heartbeat, reviews) |
| Providers | onemind-web `providers` | LLM provider config |
| Connectors | onemind-web `connectors` | Third-party integrations |
| Secrets | onemind-web `secrets` | Credential management |
| Config | OpenClaw `config` | Raw gateway config editor |
| Logs | OpenClaw `logs` + `debug` | System logs, debug output |
| Devices | New | Paired devices, Tailscale peers |
| Infrastructure | OpenClaw `infrastructure` | VPS status, services |

### Slide-Over Panels (Not Full Pages)

| Panel | Trigger | Content |
|-------|---------|---------|
| Agent Picker | Long-press chat | Switch agent mid-conversation |
| Domain Picker | Capture flow | Route content to SO/HP/LE/GE |
| Quick Actions | 3D Touch / Force Touch | Capture, voice note, task, search |
| Search | Pull-down or swipe | Universal search across vault + agents + sessions |
| Find Zeus | Codex tab | Location from Bouncie/HA/phone |

### Page Count Summary

| Tab | Pages | New | From OpenClaw/onemind-web |
|-----|-------|-----|--------------------------|
| Home | 4 | 2 | 2 |
| Chat | 5 | 2 | 3 |
| Agents | 8 | 1 | 7 |
| Codex | 10 | 8 | 2 |
| Command | 13 | 1 | 12 |
| **Total** | **40** | **14 new** | **26 mirrored** |

---

## 11. Security Considerations

| Concern | Mitigation |
|---------|------------|
| Token storage | `flutter_secure_storage` — encrypted keychain (iOS) / keystore (Android) |
| Network | Tailscale mesh — traffic never leaves the mesh. WireGuard encrypted. |
| No public API | Gateway binds to loopback, served via Tailscale only |
| No secrets in code | Token entered at first launch, stored in secure storage |
| Certificate pinning | Optional — Tailscale handles TLS termination |

---

## 12. Development Environment

### Prerequisites

```bash
# Flutter SDK
brew install flutter   # macOS
flutter doctor         # Verify setup

# Dart SDK (bundled with Flutter)
dart --version         # Need 3.4+

# Xcode (iOS builds)
xcode-select --install
sudo xcodebuild -license accept

# Android Studio (Android builds)
# Install via JetBrains Toolbox or brew install --cask android-studio

# Tailscale (required to reach gateway)
brew install --cask tailscale
# Must be logged in to the OneMind tailnet
```

### Run

```bash
cd onemind_app

# iOS Simulator
flutter run -d ios

# Android Emulator
flutter run -d android

# Linux Desktop
flutter run -d linux

# All platforms
flutter devices  # List available devices
```

### Test

```bash
flutter test                    # Unit tests
flutter test integration_test/  # Integration tests
```

---

## 13. ACP Protocol (Agent Client Protocol) — Reference

The gateway also speaks ACP v0.18.2 (962 references in dist/). This is the protocol OpenClaw uses for editor integrations (VS Code, etc.). For the Flutter app, **WebSocket JSON-RPC is the primary transport** — ACP is documented here for completeness.

ACP standardizes communication between code editors and coding agents. The wire format is JSON-RPC over WebSocket — same as the primary protocol. If future needs arise (e.g., integrating the Flutter app as a "coding client"), ACP compatibility is already there.

---

## 14. Related Repos & Resources

| Resource | Path / URL |
|----------|-----------|
| OneMind Web (UI reference) | `/Users/zeuslegacy/Documents/onemind-web/` |
| OneMind Mobile (native apps) | `/Users/zeuslegacy/Documents/onemind-mobile/` |
| Chowder iOS (WebSocket chat — **port this**) | `onemind-mobile/samples/Chowder-iOS/Chowder/` |
| VisionClaw iOS (REST bridge + audio) | `onemind-mobile/samples/VisionClaw-iOS/` |
| VisionClaw Android (Kotlin REST bridge) | `onemind-mobile/samples/VisionClaw-Android/` |
| Solarpunk theme tokens | `onemind-web/src/app/globals.css` |
| Agent workspace files | VPS: `~/.openclaw/workspace*/` |
| OpenClaw dist (protocol source) | VPS: `/usr/lib/node_modules/openclaw/dist/` |
| Gateway config | VPS: `/root/.openclaw/openclaw.json` |
| This doc | `/Users/zeuslegacy/Documents/onemind-codex/FLUTTER-APP-DEV.md` |

---

*Created: 2026-04-17. Part of the One Mind framework — SO (Sovereign Operations) / 10 Development.*
