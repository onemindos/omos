# Unified Session Management — Replacing LiveKit

> One `session_id` replaces LiveKit rooms. All channels share it. One conversation, one memory, regardless of input source.

## The Problem

Without session binding, each channel creates isolated conversations:

```
You (voice):   "Check my email"     → Session A → checks email
Robot (mic):   "Check my email"     → Session B → checks email AGAIN
You (Flutter): "What did you find?" → Session C → "Find what?"
```

LiveKit solved this with "rooms" — a shared context ID. One Mind OS solves it with Agno's native `session_id`, which is simpler and requires zero media infrastructure.

## The Solution

Every channel sends the same `session_id` to Agno's REST API. Agno maintains one conversation, one memory, one context per session.

```python
# ALL channels use this same function
async def send_to_agno(user_id: str, message: str) -> str:
    session_id = f"{user_id}-main"  # same ID always
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

## What LiveKit Did vs. What Replaces It

| LiveKit Feature | One Mind OS Replacement |
|----------------|----------------------|
| Room (shared context) | Agno `session_id` — all channels send the same ID |
| Audio transport | Gemini Live API WebSocket (STT/TTS) |
| Video transport | OM1 VLM → text via NATS (no video stream needed) |
| Participant tracking | NATS subjects per robot + Gateway channel registry |
| Data channel (text) | Agno AgentOS REST API + WebSocket SSE |
| Recording / playback | Agno session storage in PostgreSQL |

## What Agno's Session Gives You Automatically

When all channels use the same `session_id`:

- **Chat history** — All messages from all channels, in order
- **Memory** — "User prefers morning emails" persists forever
- **Knowledge context** — RAG searches use full session context
- **Team state** — If a task is mid-flight, it continues
- **User memories** — Extracted facts carry across sessions

## Deduplication

When you speak in the same room as the robot, both mics hear you. The Gateway deduplicates within a 3-second window.

```python
import hashlib, time

dedup_cache = {}

def is_duplicate(text: str, window: float = 3.0) -> bool:
    h = hashlib.md5(text.lower().strip().encode()).hexdigest()
    now = time.time()
    if h in dedup_cache and now - dedup_cache[h] < window:
        return True
    dedup_cache[h] = now
    return False
```

## Response Broadcasting

The Gateway sends Agno's response to ALL connected channels for the user. Voice channels get TTS audio. Text channels get text.

```python
connected_channels = {}  # user_id → [ws1, ws2, ...]

async def broadcast(user_id: str, response: str):
    for ws in connected_channels.get(user_id, []):
        try:
            await ws.send_json({"type": "response", "text": response})
        except Exception:
            pass
```

## Robot Event Injection

OM1 events (sensor data, detections) get injected into the user's active session as system messages, maintaining context continuity.

```python
async def on_robot_event(msg):
    event = json.loads(msg.data.decode())
    user_id = "master-delacruz"
    event_text = f"[ROBOT EVENT] {event['description']}"
    if not is_duplicate(event_text):
        response = await send_to_agno(user_id, event_text)
        if "HEARTBEAT_OK" not in response:
            await broadcast(user_id, response)
```

## When Would You Need LiveKit?

Only if you wanted multi-party video conferencing where multiple humans and robots share live video/audio in a call. For a personal AI assistant where ONE user talks to ONE entity across multiple channels, `session_id` + WebSocket gateway is all you need.
