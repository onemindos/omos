# LiveKit

> **Real-Time Voice & Video Communication**

## Overview

LiveKit is an open-source WebRTC infrastructure for building real-time voice and video applications. It powers Legacy AI's voice interface.

## Status

- **Active**: ✅ Yes
- **Deployment**: Self-hosted
- **Last Updated**: January 25, 2026

## Use Cases

1. **Legacy AI Voice** - Talk to Legacy via voice
2. **Real-Time Transcription** - Speech-to-text
3. **Voice Commands** - Hands-free control
4. **Audio Processing** - Sound analysis

## Key Features

- WebRTC-based (low latency)
- Open source & self-hostable
- Supports audio, video, and data channels
- Built-in recording
- Multi-participant rooms
- End-to-end encryption

## Architecture

```
User (Microphone)
    ↓
LiveKit Client SDK
    ↓
LiveKit Server (self-hosted)
    ↓
Agent Worker
    ↓
Legacy AI Agent
    ↓
Response (TTS)
    ↓
User (Speaker)
```

## Installation

### Server (Docker)

```bash
docker run -d \
  --name livekit \
  -p 7880:7880 \
  -p 7881:7881 \
  -p 7882:7882/udp \
  -e LIVEKIT_KEYS="devkey: secret" \
  livekit/livekit-server
```

### Python SDK

```bash
pip install livekit livekit-agents
```

## Configuration

### Environment Variables

```env
LIVEKIT_URL=ws://localhost:7880
LIVEKIT_API_KEY=devkey
LIVEKIT_API_SECRET=secret
```

### Server Config (livekit.yaml)

```yaml
port: 7880
rtc:
  port_range_start: 50000
  port_range_end: 60000
  use_external_ip: true
keys:
  devkey: secret
logging:
  level: info
```

## Agent Integration

### Voice Agent Worker

```python
from livekit import agents
from livekit.agents import llm, voice

async def entrypoint(ctx: agents.JobContext):
    # Connect to room
    await ctx.connect()

    # Initialize voice assistant
    assistant = voice.VoiceAssistant(
        vad=voice.VAD(),  # Voice Activity Detection
        stt=voice.STT(),   # Speech-to-Text
        llm=get_legacy_agent(),  # Your LLM
        tts=voice.TTS(),   # Text-to-Speech
    )

    await assistant.start(ctx.room)
```

### Room Management

```python
from livekit import api

async def create_room(room_name: str):
    client = api.LiveKitAPI(
        url=os.getenv("LIVEKIT_URL"),
        api_key=os.getenv("LIVEKIT_API_KEY"),
        api_secret=os.getenv("LIVEKIT_API_SECRET"),
    )

    room = await client.room.create_room(
        api.CreateRoomRequest(name=room_name)
    )
    return room
```

### Access Tokens

```python
from livekit import api

def create_token(room_name: str, participant_name: str) -> str:
    token = api.AccessToken(
        api_key=os.getenv("LIVEKIT_API_KEY"),
        api_secret=os.getenv("LIVEKIT_API_SECRET"),
    )

    token.with_identity(participant_name)
    token.with_name(participant_name)
    token.with_grants(api.VideoGrants(
        room_join=True,
        room=room_name,
    ))

    return token.to_jwt()
```

## Flutter Integration

```dart
import 'package:livekit_client/livekit_client.dart';

class VoiceService {
  late Room room;

  Future<void> connect(String token) async {
    room = await LiveKitClient.connect(
      'wss://livekit.onemind.zeusdelacruz.com',
      token,
      roomOptions: RoomOptions(
        adaptiveStream: true,
        dynacast: true,
      ),
    );
  }

  Future<void> enableMicrophone() async {
    await room.localParticipant?.setMicrophoneEnabled(true);
  }
}
```

## Voice Flow

1. **User speaks** → Captured by client microphone
2. **Audio sent** → WebRTC to LiveKit server
3. **VAD detects** → Speech activity detected
4. **STT processes** → Audio → Text
5. **Legacy AI responds** → Process query, generate response
6. **TTS converts** → Text → Audio
7. **Audio sent** → WebRTC to client
8. **User hears** → Response played

## Speech-to-Text Options

| Provider           | Quality   | Speed  | Cost         |
| ------------------ | --------- | ------ | ------------ |
| **Deepgram**       | Excellent | Fast   | $0.0043/min  |
| **OpenAI Whisper** | Excellent | Medium | $0.006/min   |
| **AssemblyAI**     | Good      | Fast   | $0.00025/sec |
| **Local Whisper**  | Good      | Slow   | Free         |

## Text-to-Speech Options

| Provider          | Quality | Speed     | Cost            |
| ----------------- | ------- | --------- | --------------- |
| **ElevenLabs**    | Best    | Fast      | $0.30/1K chars  |
| **OpenAI TTS**    | Great   | Fast      | $0.015/1K chars |
| **Cartesia**      | Great   | Very Fast | Variable        |
| **Coqui (Local)** | Good    | Medium    | Free            |

## Recording

```python
# Start recording
await client.room.start_recording(
    api.StartRecordingRequest(
        room=room_name,
        output=api.RecordingOutput(
            s3=api.S3Upload(
                bucket="recordings",
                region="us-east-1",
            )
        ),
    )
)
```

## Deployment

### Docker Compose

```yaml
services:
  livekit:
    image: livekit/livekit-server
    ports:
      - "7880:7880"
      - "7881:7881"
      - "7882:7882/udp"
    environment:
      - LIVEKIT_KEYS=devkey:secret
    volumes:
      - ./livekit.yaml:/etc/livekit.yaml
    command: --config /etc/livekit.yaml

  livekit-agent:
    build: ./agents/voice
    environment:
      - LIVEKIT_URL=ws://livekit:7880
      - LIVEKIT_API_KEY=devkey
      - LIVEKIT_API_SECRET=secret
```

## Best Practices

- ✅ Use VAD to reduce processing (only when speaking)
- ✅ Implement wake word detection for always-listening
- ✅ Use streaming STT/TTS for low latency
- ✅ Handle network interruptions gracefully
- ✅ Implement echo cancellation
- ✅ Test with various microphone types

## Resources

- [LiveKit Documentation](https://docs.livekit.io)
- [LiveKit Agents](https://docs.livekit.io/agents)
- [GitHub Repository](https://github.com/livekit/livekit)

---

**Related**: [Agno](./Agno.md) | [Twilio](./Twilio.md)
