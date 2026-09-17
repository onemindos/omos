# Multimodal Real-Time Voice + Vision

Press-to-talk multimodal interface for any Agno agent.

## Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Flutter App                               │
│  ┌─────────────┐                                                │
│  │   [🎤]      │──► POST /api/realtime/session                  │
│  │Voice Button │    { agent_id, enable_vision }                 │
│  └─────────────┘                                                │
│         │                                                        │
│         ▼                                                        │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │           LiveKit Room Connection                        │    │
│  │  • Publish microphone audio                              │    │
│  │  • Publish camera video (if vision enabled)              │    │
│  │  • Subscribe to agent audio track                        │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│              Real-Time Agent Worker (Server)                     │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │           Provider (Gemini/OpenAI/Pipeline)              │    │
│  │  • Receives audio from user via LiveKit                  │    │
│  │  • Processes with AI model                               │    │
│  │  • Sends audio response back                             │    │
│  │  • Uses agent instructions and tools                     │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

## Providers

| Provider | Latency | How It Works | Best For |
|----------|---------|--------------|----------|
| `gemini` | ~200ms | Gemini 2.0 Flash Live (speech-to-speech) | Production voice |
| `openai` | ~300ms | GPT-4o Realtime API (speech-to-speech) | OpenAI ecosystem |
| `pipeline` | ~1-2s | STT → Any Agent → TTS | Custom models |

## API Endpoints

### Create Session

```http
POST /api/realtime/session
Content-Type: application/json

{
  "agent_id": "legacy-mother-ai",
  "provider": "gemini",
  "enable_vision": true,
  "voice": "Puck",
  "user_id": "zeus"
}
```

**Response:**
```json
{
  "room_name": "rt-legacy-abc12345",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "ws_url": "wss://livekit.example.com",
  "agent_id": "legacy-mother-ai",
  "provider": "gemini",
  "vision_enabled": true,
  "session_id": "abc12345-...",
  "voice": "Puck"
}
```

### Quick Start (Auto-Select Provider)

```http
POST /api/realtime/quick-start?agent_id=legacy-mother-ai&vision=true
```

### List Active Sessions

```http
GET /api/realtime/sessions
```

### End Session

```http
DELETE /api/realtime/session/{room_name}
```

### List Providers

```http
GET /api/realtime/providers
```

## Flutter Integration

### 1. Add Dependencies

```yaml
# pubspec.yaml
dependencies:
  livekit_client: ^2.1.0
  permission_handler: ^11.0.0
```

### 2. Create Realtime Service

```dart
// lib/services/realtime_service.dart
import 'package:livekit_client/livekit_client.dart';

class RealtimeService {
  Room? _room;
  LocalAudioTrack? _audioTrack;
  LocalVideoTrack? _videoTrack;

  /// Start a real-time voice session
  Future<void> startSession({
    required String agentId,
    bool enableVision = false,
    String provider = 'gemini',
  }) async {
    // 1. Create session via API
    final response = await http.post(
      Uri.parse('$baseUrl/api/realtime/session'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'agent_id': agentId,
        'provider': provider,
        'enable_vision': enableVision,
      }),
    );

    final session = jsonDecode(response.body);

    // 2. Connect to LiveKit room
    _room = Room();
    await _room!.connect(
      session['ws_url'],
      session['token'],
      roomOptions: RoomOptions(
        adaptiveStream: true,
        dynacast: true,
      ),
    );

    // 3. Publish microphone
    _audioTrack = await LocalAudioTrack.create();
    await _room!.localParticipant?.publishAudioTrack(_audioTrack!);

    // 4. Publish camera if vision enabled
    if (enableVision) {
      _videoTrack = await LocalVideoTrack.createCameraTrack();
      await _room!.localParticipant?.publishVideoTrack(_videoTrack!);
    }

    // 5. Listen for agent audio
    _room!.addListener(_onRoomEvent);
  }

  void _onRoomEvent(RoomEvent event) {
    if (event is TrackSubscribedEvent) {
      if (event.track.kind == TrackType.AUDIO) {
        // Agent is speaking - audio plays automatically
        print('Agent speaking...');
      }
    }
  }

  /// End the session
  Future<void> endSession() async {
    await _audioTrack?.dispose();
    await _videoTrack?.dispose();
    await _room?.disconnect();
    _room = null;
  }
}
```

### 3. Create Voice Button Widget

```dart
// lib/widgets/voice_button.dart
import 'package:flutter/material.dart';

class VoiceButton extends StatefulWidget {
  final String agentId;
  final bool enableVision;

  const VoiceButton({
    super.key,
    this.agentId = 'legacy-mother-ai',
    this.enableVision = false,
  });

  @override
  State<VoiceButton> createState() => _VoiceButtonState();
}

class _VoiceButtonState extends State<VoiceButton> {
  final _realtimeService = RealtimeService();
  bool _isActive = false;

  Future<void> _toggleVoice() async {
    if (_isActive) {
      await _realtimeService.endSession();
    } else {
      await _realtimeService.startSession(
        agentId: widget.agentId,
        enableVision: widget.enableVision,
      );
    }
    setState(() => _isActive = !_isActive);
  }

  @override
  Widget build(BuildContext context) {
    return FloatingActionButton(
      onPressed: _toggleVoice,
      backgroundColor: _isActive ? Colors.red : Colors.blue,
      child: Icon(
        _isActive ? Icons.stop : Icons.mic,
        color: Colors.white,
      ),
    );
  }
}
```

### 4. Full-Screen Voice Mode

```dart
// lib/screens/voice_mode_screen.dart
class VoiceModeScreen extends StatefulWidget {
  final String agentId;
  final bool enableVision;

  const VoiceModeScreen({
    super.key,
    required this.agentId,
    this.enableVision = false,
  });

  @override
  State<VoiceModeScreen> createState() => _VoiceModeScreenState();
}

class _VoiceModeScreenState extends State<VoiceModeScreen> {
  final _realtimeService = RealtimeService();
  bool _isConnected = false;
  String _status = 'Connecting...';

  @override
  void initState() {
    super.initState();
    _startSession();
  }

  Future<void> _startSession() async {
    try {
      await _realtimeService.startSession(
        agentId: widget.agentId,
        enableVision: widget.enableVision,
      );
      setState(() {
        _isConnected = true;
        _status = 'Listening...';
      });
    } catch (e) {
      setState(() => _status = 'Error: $e');
    }
  }

  @override
  void dispose() {
    _realtimeService.endSession();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            // Animated orb
            AnimatedContainer(
              duration: const Duration(milliseconds: 300),
              width: _isConnected ? 200 : 150,
              height: _isConnected ? 200 : 150,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: _isConnected ? Colors.blue : Colors.grey,
                boxShadow: _isConnected
                    ? [
                        BoxShadow(
                          color: Colors.blue.withOpacity(0.5),
                          blurRadius: 50,
                          spreadRadius: 20,
                        ),
                      ]
                    : null,
              ),
              child: Icon(
                _isConnected ? Icons.mic : Icons.mic_off,
                size: 80,
                color: Colors.white,
              ),
            ),
            const SizedBox(height: 40),
            Text(
              _status,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 18,
              ),
            ),
            const SizedBox(height: 20),
            if (widget.enableVision)
              const Text(
                '📹 Camera active',
                style: TextStyle(color: Colors.green),
              ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => Navigator.pop(context),
        backgroundColor: Colors.red,
        child: const Icon(Icons.close),
      ),
    );
  }
}
```

## Running the Worker

The real-time worker must be running to handle voice sessions:

```bash
# Start the worker
python -m core.multimodal.realtime_worker dev

# Or in production
python -m core.multimodal.realtime_worker start --production
```

## Environment Variables

```bash
# LiveKit (required)
LIVEKIT_URL=wss://your-livekit-server.livekit.cloud
LIVEKIT_API_KEY=your_api_key
LIVEKIT_API_SECRET=your_api_secret

# Gemini (for gemini provider)
GOOGLE_API_KEY=your_google_api_key

# OpenAI (for openai provider)
OPENAI_API_KEY=your_openai_api_key

# Pipeline mode (for pipeline provider)
DEEPGRAM_API_KEY=your_deepgram_key
CARTESIA_API_KEY=your_cartesia_key
```

## Voices

### Gemini Voices
| Voice | Description |
|-------|-------------|
| Puck | Playful, energetic |
| Charon | Deep, mysterious |
| Kore | Warm, nurturing |
| Fenrir | Bold, commanding |
| Aoede | Musical, expressive |

### OpenAI Voices
| Voice | Description |
|-------|-------------|
| alloy | Neutral, balanced |
| echo | Male, warm |
| shimmer | Female, soft |
| ash | Male, measured |
| ballad | Storytelling |
| coral | Female, friendly |
| sage | Authoritative |
| verse | Expressive |

## Architecture Files

```
core/multimodal/
├── realtime.py          # Session manager
├── realtime_api.py      # FastAPI endpoints
├── realtime_worker.py   # LiveKit worker
└── __init__.py          # Exports
```

## Example: Security Camera Agent

```dart
// Create a voice session with security agent + camera
await realtimeService.startSession(
  agentId: 'security-agent',
  enableVision: true,
  provider: 'gemini',
);

// User says: "What do you see?"
// Agent responds with camera analysis
```

## Example: Research Assistant

```dart
// Voice-only session for quick research
await realtimeService.startSession(
  agentId: 'research-agent',
  enableVision: false,
  provider: 'pipeline',  // Uses Claude via pipeline
);
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
