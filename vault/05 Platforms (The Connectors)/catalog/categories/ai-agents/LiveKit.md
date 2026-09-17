# LiveKit

> Real-time audio, video, and data communication for voice interfaces and video calls.

---

## Overview

| Attribute         | Value                                      |
| ----------------- | ------------------------------------------ |
| **Type**          | WebRTC SFU (Selective Forwarding Unit)     |
| **Protocol**      | WebRTC, WebSocket                          |
| **Default Ports** | 7880 (HTTP), 7881 (RTC), 7882 (TURN)       |
| **License**       | Apache 2.0                                 |
| **Documentation** | [docs.livekit.io](https://docs.livekit.io) |

---

## What is LiveKit?

LiveKit is an open-source real-time communication platform that provides:

- **Audio/Video rooms** — Multi-participant calls
- **Voice AI** — Real-time voice agents
- **Screen sharing** — Share screens and applications
- **Data channels** — Low-latency data messaging
- **Recording** — Capture sessions
- **Egress/Ingress** — Stream to/from external sources

---

## OneMind Usage

### Primary Use Cases

```
┌─────────────────────────────────────────────────────────────────┐
│                       OneMind LiveKit                            │
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │  Voice AI    │    │   Watch UI   │    │  Video Calls │       │
│  │  Interface   │    │   Voice      │    │  w/ Agent    │       │
│  └──────────────┘    └──────────────┘    └──────────────┘       │
│         │                   │                   │                │
│         └───────────────────┼───────────────────┘                │
│                             ▼                                    │
│                    ┌──────────────┐                              │
│                    │   LiveKit    │                              │
│                    │    Server    │                              │
│                    └──────────────┘                              │
│                             │                                    │
│         ┌───────────────────┼───────────────────┐                │
│         ▼                   ▼                   ▼                │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │   Legacy     │    │  Transcribe  │    │    TTS       │       │
│  │   Agent      │    │   (Deepgram) │    │   (11Labs)   │       │
│  └──────────────┘    └──────────────┘    └──────────────┘       │
└─────────────────────────────────────────────────────────────────┘
```

| Use Case            | Description                       |
| ------------------- | --------------------------------- |
| **Voice Interface** | Talk to Legacy agent hands-free   |
| **ESP32 Watch**     | Voice commands from wearable      |
| **Mobile App**      | Voice chat on the go              |
| **Video Presence**  | See Legacy's avatar while talking |
| **Screen Share**    | Share context with agent          |
| **Recording**       | Capture conversations for memory  |

---

## Docker Compose

```yaml
services:
  livekit:
    image: livekit/livekit-server:latest
    container_name: onemind-livekit
    restart: unless-stopped
    ports:
      - "7880:7880" # HTTP API
      - "7881:7881" # RTC (WebRTC)
      - "7882:7882" # TURN/TLS
    environment:
      - LIVEKIT_KEYS=devkey:secret
    command: --config /etc/livekit.yaml --dev
    volumes:
      - ./livekit.yaml:/etc/livekit.yaml:ro
    networks:
      - onemind-network

  livekit-egress:
    image: livekit/egress:latest
    container_name: onemind-livekit-egress
    restart: unless-stopped
    environment:
      - EGRESS_CONFIG_FILE=/etc/egress.yaml
    volumes:
      - ./egress.yaml:/etc/egress.yaml:ro
      - egress_data:/data
    networks:
      - onemind-network

volumes:
  egress_data:
```

---

## Configuration

### livekit.yaml

```yaml
# livekit.yaml
port: 7880
rtc:
  port_range_start: 50000
  port_range_end: 60000
  tcp_port: 7881
  use_external_ip: true

keys:
  devkey: secret # Development
  # Production: use secure generated keys

logging:
  level: info
  json: false

room:
  empty_timeout: 300 # 5 minutes
  max_participants: 50

turn:
  enabled: true
  domain: turn.onemind.cloud
  tls_port: 5349
  udp_port: 3478
```

### Production Keys

```bash
# Generate secure API key
livekit-server generate-keys

# Output:
# API Key: APIsomeRandomKey
# API Secret: someVeryLongRandomSecret
```

---

## Voice Agent Architecture

### Real-time Voice Pipeline

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│   User   │────▶│ LiveKit  │────▶│ Deepgram │────▶│  Legacy  │
│   Mic    │     │  Server  │     │   STT    │     │  Agent   │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                       │                                 │
                       │                                 │
                       ▼                                 ▼
                ┌──────────┐     ┌──────────┐     ┌──────────┐
                │   User   │◀────│ LiveKit  │◀────│  11Labs  │
                │ Speaker  │     │  Server  │     │   TTS    │
                └──────────┘     └──────────┘     └──────────┘
```

### Voice Agent (Python)

```python
from livekit import agents
from livekit.agents import JobContext, WorkerOptions, cli
from livekit.agents.pipeline import VoicePipelineAgent
from livekit.plugins import deepgram, silero, openai

async def entrypoint(ctx: JobContext):
    # Connect to room
    await ctx.connect()

    # Create voice agent
    agent = VoicePipelineAgent(
        vad=silero.VAD.load(),
        stt=deepgram.STT(),
        llm=openai.LLM(model="gpt-4"),
        tts=openai.TTS(),
    )

    # Start the agent
    agent.start(ctx.room)

    # Initial greeting
    await agent.say("Hello Zeus, I'm Legacy. How can I help?")

if __name__ == "__main__":
    cli.run_app(WorkerOptions(entrypoint_fnc=entrypoint))
```

### With Bedrock/Claude

```python
from livekit.plugins import anthropic

async def entrypoint(ctx: JobContext):
    await ctx.connect()

    agent = VoicePipelineAgent(
        vad=silero.VAD.load(),
        stt=deepgram.STT(),
        llm=anthropic.LLM(model="claude-3-5-sonnet-20241022"),
        tts=elevenlabs.TTS(voice_id="your_voice_id"),
    )

    agent.start(ctx.room)
```

---

## Client Integration

### JavaScript/React

```typescript
import { Room, RoomEvent } from "livekit-client";

const room = new Room();

// Connect to room
await room.connect("wss://livekit.onemind.cloud", token);

// Handle events
room.on(RoomEvent.TrackSubscribed, (track, publication, participant) => {
  if (track.kind === "audio") {
    const audioElement = track.attach();
    document.body.appendChild(audioElement);
  }
});

// Publish microphone
await room.localParticipant.setMicrophoneEnabled(true);
```

### ESP32 Watch (Simplified)

```cpp
// ESP32 sends audio to gateway, which bridges to LiveKit
#include <WiFi.h>
#include <WebSocketsClient.h>
#include <driver/i2s.h>

WebSocketsClient ws;

void setup() {
  // Initialize I2S for microphone
  i2s_config_t i2s_config = {
    .mode = I2S_MODE_MASTER | I2S_MODE_RX,
    .sample_rate = 16000,
    .bits_per_sample = I2S_BITS_PER_SAMPLE_16BIT,
    // ...
  };
  i2s_driver_install(I2S_NUM_0, &i2s_config, 0, NULL);

  // Connect to WebSocket bridge
  ws.begin("gateway.onemind.local", 8000, "/ws/voice");
}

void loop() {
  // Read audio and send via WebSocket
  int16_t samples[160];
  size_t bytes_read;
  i2s_read(I2S_NUM_0, samples, sizeof(samples), &bytes_read, portMAX_DELAY);
  ws.sendBIN((uint8_t*)samples, bytes_read);
}
```

---

## Room Management

### Create Room (Python)

```python
from livekit import api

# Create client
lk = api.LiveKitAPI(
    url='http://localhost:7880',
    api_key='devkey',
    api_secret='secret'
)

# Create room
room = await lk.room.create_room(
    api.CreateRoomRequest(
        name='zeus-voice-session',
        empty_timeout=300,
        max_participants=2
    )
)
```

### Generate Access Token

```python
from livekit import api
import time

def create_token(room_name: str, participant_name: str) -> str:
    token = api.AccessToken(
        api_key='devkey',
        api_secret='secret'
    )

    token.with_identity(participant_name)
    token.with_name(participant_name)
    token.with_grants(api.VideoGrants(
        room_join=True,
        room=room_name,
        can_publish=True,
        can_subscribe=True
    ))

    return token.to_jwt()

# Create token for Zeus
token = create_token('zeus-voice-session', 'zeus')
```

### FastAPI Token Endpoint

```python
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI()

class TokenRequest(BaseModel):
    room: str
    identity: str

@app.post("/api/livekit/token")
async def get_token(request: TokenRequest):
    token = create_token(request.room, request.identity)
    return {"token": token}
```

---

## Agent Workers

### Running Voice Agent

```bash
# Install dependencies
pip install livekit-agents livekit-plugins-deepgram livekit-plugins-silero livekit-plugins-openai

# Run worker
python voice_agent.py dev
```

### Agent Configuration

```python
# agent_config.py
from livekit.agents import WorkerOptions

worker_options = WorkerOptions(
    entrypoint_fnc=entrypoint,
    api_key="devkey",
    api_secret="secret",
    ws_url="ws://localhost:7880",
)

if __name__ == "__main__":
    from livekit.agents import cli
    cli.run_app(worker_options)
```

---

## Recording (Egress)

### Record Room

```python
from livekit import api

async def start_recording(room_name: str):
    lk = api.LiveKitAPI()

    egress = await lk.egress.start_room_composite_egress(
        api.RoomCompositeEgressRequest(
            room_name=room_name,
            file_outputs=[
                api.EncodedFileOutput(
                    file_type=api.EncodedFileType.MP4,
                    filepath="/recordings/{room_name}/{time}.mp4"
                )
            ]
        )
    )

    return egress.egress_id
```

### egress.yaml

```yaml
# egress.yaml
api_key: devkey
api_secret: secret
ws_url: ws://livekit:7880

# S3 storage
s3:
  access_key: ${AWS_ACCESS_KEY_ID}
  secret: ${AWS_SECRET_ACCESS_KEY}
  region: us-east-1
  bucket: onemind-recordings
```

---

## Monitoring

### Server Stats

```bash
# Check server status
curl http://localhost:7880/

# List rooms
curl -H "Authorization: Bearer $(livekit-cli create-token --api-key devkey --api-secret secret --list --valid-for 1m)" \
  http://localhost:7880/twirp/livekit.RoomService/ListRooms
```

### livekit-cli

```bash
# Install CLI
brew install livekit-cli

# Configure
livekit-cli project add \
  --url http://localhost:7880 \
  --api-key devkey \
  --api-secret secret

# List rooms
livekit-cli room list

# Join room (for testing)
livekit-cli room join --room test-room --identity tester
```

### Health Check

```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "http://localhost:7880"]
  interval: 30s
  timeout: 10s
  retries: 3
```

---

## Production Deployment

### TURN Server (Required for NAT traversal)

```yaml
# For production, configure TURN
turn:
  enabled: true
  domain: turn.onemind.cloud
  cert_file: /etc/livekit/turn.crt
  key_file: /etc/livekit/turn.key
  tls_port: 5349
  udp_port: 3478
  external_tls: true
```

### SSL/TLS

```yaml
# Use with nginx or load balancer
rtc:
  use_external_ip: true

# Or configure TLS directly
tls:
  cert_file: /etc/livekit/cert.pem
  key_file: /etc/livekit/key.pem
```

### Environment Variables

```bash
# LiveKit Server
LIVEKIT_KEYS=APIkey:secretkey
LIVEKIT_LOG_LEVEL=info
LIVEKIT_REDIS_ADDRESS=redis:6379

# Voice Agent
LIVEKIT_URL=wss://livekit.onemind.cloud
LIVEKIT_API_KEY=APIkey
LIVEKIT_API_SECRET=secretkey
DEEPGRAM_API_KEY=your_deepgram_key
OPENAI_API_KEY=your_openai_key
```

---

## Troubleshooting

### Connection Issues

```bash
# Check server is running
curl http://localhost:7880

# Check WebRTC connectivity
livekit-cli room join --room test --identity test

# Check firewall ports
sudo ufw status
```

### Audio Issues

```bash
# Check audio tracks
livekit-cli room list-participants --room your-room

# Verify codec support
# LiveKit uses Opus for audio by default
```

### Agent Not Responding

```bash
# Check agent worker logs
python voice_agent.py dev --log-level debug

# Verify STT/TTS API keys
echo $DEEPGRAM_API_KEY
echo $OPENAI_API_KEY
```

---

## Resources

- [LiveKit Documentation](https://docs.livekit.io)
- [LiveKit GitHub](https://github.com/livekit/livekit)
- [LiveKit Agents](https://github.com/livekit/agents)
- [Python SDK](https://github.com/livekit/python-sdks)
- [JavaScript SDK](https://github.com/livekit/client-sdk-js)

---

_Last updated: January 2026_


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
