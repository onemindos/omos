# Interaction Gateway

> The topmost layer of One Mind OS. Accepts voice, text, and channel inputs from anywhere. Routes everything to Agno as text. Returns responses via TTS or text to the originating channel.

## Purpose

The Interaction Gateway sits ABOVE both Agno and OM1. It is the single entry point for all human interaction with Legacy. It replaces LiveKit rooms with a lightweight FastAPI service that manages unified sessions across all input channels.

## Position in the Stack

```
YOU (voice, text, chat, AR glasses, robot mic)
 │
 ▼
┌─────────────────────────────────────────────────┐
│         INTERACTION GATEWAY                      │
│  (this document)                                 │
│                                                  │
│  Accepts ALL inputs → converts to text →         │
│  sends to Agno with shared session_id →          │
│  receives response → converts back to            │
│  voice/text → sends to all active channels       │
└──────────────────┬──────────────────────────────┘
                   ▼
             AGNO (Brain)
                   ▼
             NATS (Nervous System)
                   ▼
             OM1 (Body)
```

## Supported Input Channels

### Voice — Gemini Live API (Primary)

Real-time conversational voice using Google's Gemini Live API over WebSocket.

- **STT**: Gemini Live handles speech-to-text with voice activity detection (VAD), interruption handling, and turn-taking
- **TTS**: Gemini Live or ElevenLabs for natural speech output
- **Transport**: WebSocket connection from Flutter app or web client
- **Latency**: ~200-500ms STT, ~300ms TTS

Why Gemini Live and not Agno's native audio: Agno supports audio input/output via OpenAI audio models, but real-time streaming voice (live conversation with interruption) is not yet natively supported — it's an open feature request. Gemini Live fills this gap.

```python
# Voice channel WebSocket handler
@app.websocket("/voice/{user_id}")
async def voice_channel(ws: WebSocket, user_id: str):
    await ws.accept()
    connected_channels.setdefault(user_id, []).append(ws)
    try:
        while True:
            audio_data = await ws.receive_bytes()
            text = await gemini_live_stt(audio_data)
            if is_duplicate(text):
                continue
            response = await send_to_agno(user_id, text)
            audio_out = await tts_convert(response)  # Gemini or ElevenLabs
            await ws.send_bytes(audio_out)
            await broadcast(user_id, response)
    finally:
        connected_channels[user_id].remove(ws)
```

### Text Chat — Flutter App / Web UI

WebSocket connection from the Flutter mobile/desktop app or a web-based chat interface.

```python
@app.websocket("/chat/{user_id}")
async def chat_channel(ws: WebSocket, user_id: str):
    await ws.accept()
    connected_channels.setdefault(user_id, []).append(ws)
    try:
        while True:
            data = await ws.receive_json()
            if is_duplicate(data["text"]):
                continue
            response = await send_to_agno(user_id, data["text"])
            await broadcast(user_id, response)
    finally:
        connected_channels[user_id].remove(ws)
```

### Channel Adapters — Telegram, WhatsApp, Slack, etc.

Each messaging platform gets its own adapter that converts platform-specific messages to the unified format. Inspired by the OpenClaw channel adapter pattern.

```python
# Telegram adapter example
from telegram import Update
from telegram.ext import Application, MessageHandler, filters

async def handle_telegram(update: Update, context):
    user_id = resolve_user_id(update.effective_user.id)
    text = update.message.text
    if is_duplicate(text):
        return
    response = await send_to_agno(user_id, text)
    await update.message.reply_text(response)
    await broadcast(user_id, response)
```

Supported platforms (add as needed):

- Telegram (bot API)
- WhatsApp (via Twilio or WhatsApp Business API)
- Slack (Socket Mode — no public URL needed)
- Discord (bot)
- iMessage (via BlueBubbles or similar bridge)
- Signal (via signal-cli)
- SMS (via Twilio)

### Robot Mic — OM1 ASR via NATS

When someone speaks directly to the robot, OM1's built-in Google ASR captures it and publishes to NATS. The Gateway subscribes and forwards to Agno.

```python
async def start_robot_voice_listener():
    nc = await nats.connect(NATS_URL)

    async def on_robot_speech(msg):
        data = json.loads(msg.data.decode())
        user_id = "master-delacruz"
        text = data.get("transcript", "")
        if text and not is_duplicate(text):
            response = await send_to_agno(user_id, text)
            # Send response back to robot for TTS
            await nc.publish(
                f"onemind.commands.{data['robot_id']}",
                json.dumps({"action": "speak", "text": response}).encode()
            )
            await broadcast(user_id, response)

    await nc.subscribe("onemind.audio.*", cb=on_robot_speech)
```

### AR Glasses — Brilliant Labs Frame

The Frame sends text/voice input via Bluetooth to the Flutter app, which routes through the same Gateway. No special handling needed — it's just another input channel feeding the same session.

## Session Routing

All channels share one `session_id` per user. This is the core mechanism that replaces LiveKit rooms.

```python
active_sessions = {}  # user_id → session_id

def get_session_id(user_id: str) -> str:
    if user_id not in active_sessions:
        active_sessions[user_id] = f"{user_id}-main"
    return active_sessions[user_id]

async def send_to_agno(user_id: str, message: str) -> str:
    session_id = get_session_id(user_id)
    async with httpx.AsyncClient() as client:
        resp = await client.post(
            f"{AGNO_URL}/teams/onemind/runs",
            json={
                "message": message,
                "session_id": session_id,
                "user_id": user_id,
                "stream": False,
            }
        )
        return resp.json()["content"]
```

See [UNIFIED-SESSION.md](./UNIFIED-SESSION.md) for full session management details.

## Deduplication

When you speak in the same room as the robot, both your phone mic and the robot's mic hear you. The Gateway prevents double-processing.

```python
import hashlib, time

dedup_cache = {}  # hash → timestamp

def is_duplicate(text: str, window: float = 3.0) -> bool:
    h = hashlib.md5(text.lower().strip().encode()).hexdigest()
    now = time.time()
    # Clean old entries
    dedup_cache = {k: v for k, v in dedup_cache.items() if now - v < window}
    if h in dedup_cache:
        return True
    dedup_cache[h] = now
    return False
```

For more sophisticated dedup (fuzzy matching when STT produces slightly different transcripts from different mics), use embedding similarity with a threshold of 0.9.

## Response Broadcasting

When Agno responds, the Gateway sends the response to ALL connected channels for that user.

```python
connected_channels = {}  # user_id → [websocket1, ws2, ...]

async def broadcast(user_id: str, response: str):
    for ws in connected_channels.get(user_id, []):
        try:
            await ws.send_json({"type": "response", "text": response})
        except Exception:
            pass  # Channel disconnected
```

Voice channels receive TTS audio. Text channels receive text. Telegram/WhatsApp receive formatted messages. Each adapter handles its own output format.

## NATS Event Injection

Robot events (sensor data, detections, alerts) flow through NATS and get injected into the user's active Agno session as system events.

```python
async def start_nats_event_listener():
    nc = await nats.connect(NATS_URL)

    async def on_robot_event(msg):
        event = json.loads(msg.data.decode())
        user_id = "master-delacruz"
        event_text = f"[ROBOT EVENT] {event['description']}"
        if not is_duplicate(event_text):
            response = await send_to_agno(user_id, event_text)
            # Only broadcast if response is actionable
            if "HEARTBEAT_OK" not in response:
                await broadcast(user_id, response)

    await nc.subscribe("onemind.vision.*", cb=on_robot_event)
    await nc.subscribe("onemind.alerts.*", cb=on_robot_event)
```

## Deployment

The Gateway runs as a standalone FastAPI service. It can be deployed alongside Agno on Railway or on any server with network access to both the NATS server and Agno's AgentOS API.

```bash
# Start the gateway
uvicorn gateway.main:app --host 0.0.0.0 --port 8000

# Required environment variables
AGNO_URL=http://localhost:7777
NATS_URL=nats://localhost:4222
GEMINI_API_KEY=your-key
ELEVENLABS_API_KEY=your-key  # optional, for TTS
```

## File Structure

```
gateway/
├── main.py              # FastAPI app, startup hooks
├── session_router.py    # Session ID management, dedup
├── channels/
│   ├── voice.py         # Gemini Live API WebSocket
│   ├── chat.py          # Flutter/Web WebSocket
│   ├── telegram.py      # Telegram bot adapter
│   ├── whatsapp.py      # WhatsApp adapter
│   ├── slack.py         # Slack Socket Mode adapter
│   └── nats_bridge.py   # Robot events → session injection
├── tts/
│   ├── gemini.py        # Gemini TTS
│   └── elevenlabs.py    # ElevenLabs TTS
└── config.py            # Environment, channel registry
```
