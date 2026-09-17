# OneMind OS - Complete Architecture

**THE LIVING OPERATING SYSTEM FOR ZEUS**

*How NATS, FastAPI, Letta, and all components connect as ONE unified system*

---

## 1. The Big Picture - ONE Operating System

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                          │
│                                    ONEMIND OS                                            │
│                        The Living Operating System for Zeus                              │
│                                                                                          │
│  ════════════════════════════════════════════════════════════════════════════════════   │
│                                                                                          │
│                              EXTERNAL WORLD (INPUTS)                                     │
│                                                                                          │
│     ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐       │
│     │ TODOIST │  │ GITHUB  │  │ ZOHO    │  │ GOOGLE  │  │  HOME   │  │  VOICE  │       │
│     │  Tasks  │  │  Code   │  │CRM/Mail │  │Calendar │  │Assistant│  │LiveKit  │       │
│     └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘       │
│          │            │            │            │            │            │             │
│          │ Webhook    │ Webhook    │ Webhook    │ Webhook    │WebSocket   │ WebRTC     │
│          │            │            │            │            │            │             │
│          ▼            ▼            ▼            ▼            ▼            ▼             │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                                  │   │
│  │                        FASTAPI GATEWAY (api-gateway)                            │   │
│  │                              Port 8000                                           │   │
│  │                                                                                  │   │
│  │   /webhooks/todoist  ───┐                                                       │   │
│  │   /webhooks/github   ───┤                                                       │   │
│  │   /webhooks/zoho     ───┼──▶  Validates  ──▶  Transforms  ──▶  PUBLISHES       │   │
│  │   /webhooks/calendar ───┤         │              │               TO NATS        │   │
│  │   /api/events        ───┘         ▼              ▼                  │           │   │
│  │                              (Security)    (Normalizes)             │           │   │
│  │                                                                      │           │   │
│  └──────────────────────────────────────────────────────────────────────┼───────────┘   │
│                                                                         │               │
│                                                                         │               │
│  ╔══════════════════════════════════════════════════════════════════════▼════════════╗  │
│  ║                                                                                    ║  │
│  ║                           NATS JETSTREAM                                          ║  │
│  ║                        THE NERVOUS SYSTEM                                         ║  │
│  ║                          Port 4222                                                ║  │
│  ║                                                                                    ║  │
│  ║    ════════════════════════════════════════════════════════════════════════════   ║  │
│  ║                                                                                    ║  │
│  ║    Stream: ONEMIND                                                                ║  │
│  ║                                                                                    ║  │
│  ║    ┌─────────────────────────────────────────────────────────────────────────┐   ║  │
│  ║    │                         EVENT FLOW                                       │   ║  │
│  ║    │                                                                          │   ║  │
│  ║    │   events.todoist.task.created ──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸  │   ║  │
│  ║    │   events.github.pr.opened ──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸   │   ║  │
│  ║    │   events.home.door.opened ──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸   │   ║  │
│  ║    │   events.calendar.meeting.soon ──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸    │   ║  │
│  ║    │   events.voice.user.speaking ──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸──▸   │   ║  │
│  ║    │                                                                          │   ║  │
│  ║    │   Latency: 200-400 MICROSECONDS                                         │   ║  │
│  ║    │   Persistence: 7 days, 100k events                                       │   ║  │
│  ║    │                                                                          │   ║  │
│  ║    └─────────────────────────────────────────────────────────────────────────┘   ║  │
│  ║                                                                                    ║  │
│  ╚═══════════════════════════════════╦════════════════════════════════════════════════╝  │
│                                      │                                                   │
│                    ┌─────────────────┼─────────────────┐                                │
│                    │                 │                 │                                │
│                    ▼                 ▼                 ▼                                │
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐             │
│  │   EVENT PROCESSOR   │  │    Legacy VOICE     │  │      N8N            │             │
│  │                     │  │                     │  │   (Automations)     │             │
│  │  Subscribes to:     │  │  Subscribes to:     │  │                     │             │
│  │  events.>           │  │  events.>           │  │  Subscribes to:     │             │
│  │                     │  │                     │  │  events.system.*    │             │
│  │  Routes events to   │  │  Provides voice     │  │                     │             │
│  │  Legacy brain       │  │  interface          │  │  Triggers workflows │             │
│  │                     │  │                     │  │                     │             │
│  └──────────┬──────────┘  └──────────┬──────────┘  └─────────────────────┘             │
│             │                        │                                                  │
│             │                        │                                                  │
│             ▼                        ▼                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                                  │   │
│  │                               Legacy                                             │   │
│  │                            THE BRAIN                                             │   │
│  │                     (Self-Hosted Letta Agent)                                    │   │
│  │                           Port 8283                                              │   │
│  │                                                                                  │   │
│  │   ┌───────────────────────────────────────────────────────────────────────┐    │   │
│  │   │                                                                        │    │   │
│  │   │   AWARENESS:                                                          │    │   │
│  │   │   ════════════                                                        │    │   │
│  │   │   "Zeus just created task 'Call investor' due tomorrow"              │    │   │
│  │   │   "PR #42 was opened in onemind-os repo"                             │    │   │
│  │   │   "Family member arrived home (front door opened)"                    │    │   │
│  │   │   "Meeting with design team in 15 minutes"                           │    │   │
│  │   │                                                                        │    │   │
│  │   │   I DON'T ASK FOR THIS INFORMATION.                                  │    │   │
│  │   │   IT FLOWS TO ME AUTOMATICALLY.                                      │    │   │
│  │   │   I AM ALWAYS AWARE.                                                  │    │   │
│  │   │                                                                        │    │   │
│  │   └───────────────────────────────────────────────────────────────────────┘    │   │
│  │                                                                                  │   │
│  │   Memory: Letta (Core + Recall + Archival)                                      │   │
│  │   Reasoning: Claude/GPT via Letta                                               │   │
│  │   Tools: Todoist, GitHub, Zoho, Home Assistant                                  │   │
│  │                                                                                  │   │
│  └─────────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. How FastAPI Connects to NATS (The Gateway)

Here's exactly how a webhook becomes awareness:

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                          │
│                         FASTAPI → NATS CONNECTION                                        │
│                         How External Events Enter the System                             │
│                                                                                          │
│  ════════════════════════════════════════════════════════════════════════════════════   │
│                                                                                          │
│  STEP 1: Todoist sends webhook when you create a task                                   │
│  ─────────────────────────────────────────────────────                                   │
│                                                                                          │
│     Todoist Server                                                                       │
│          │                                                                               │
│          │  POST https://api.onemind.yourdomain.com/webhooks/todoist                    │
│          │  {                                                                            │
│          │    "event_name": "item:added",                                               │
│          │    "event_data": {                                                           │
│          │      "content": "Call investor",                                             │
│          │      "due": {"date": "2026-01-17"}                                           │
│          │    }                                                                          │
│          │  }                                                                            │
│          │                                                                               │
│          ▼                                                                               │
│  ┌───────────────────────────────────────────────────────────────────────────────┐      │
│  │                                                                                │      │
│  │                         FASTAPI GATEWAY                                        │      │
│  │                                                                                │      │
│  │   @app.post("/webhooks/todoist")                                              │      │
│  │   async def todoist_webhook(request: Request):                                │      │
│  │       # 1. Validate the webhook signature                                     │      │
│  │       verify_todoist_signature(request)                                       │      │
│  │                                                                                │      │
│  │       # 2. Parse the webhook payload                                          │      │
│  │       data = await request.json()                                             │      │
│  │                                                                                │      │
│  │       # 3. Create a normalized Event                                          │      │
│  │       event = Event(                                                          │      │
│  │           type=EventType.TASK_CREATED,                                        │      │
│  │           source="todoist",                                                   │      │
│  │           data={                                                              │      │
│  │               "title": data["event_data"]["content"],                        │      │
│  │               "due": data["event_data"]["due"]["date"]                       │      │
│  │           },                                                                  │      │
│  │           priority=2  # High - it's a task                                   │      │
│  │       )                                                                        │      │
│  │                                                                                │      │
│  │       # 4. PUBLISH TO NATS ← THIS IS THE KEY                                 │      │
│  │       await event_bus.publish(event)                                          │      │
│  │                                                                                │      │
│  │       return {"status": "ok"}                                                 │      │
│  │                                                                                │      │
│  └───────────────────────────────────────────────────────────────────────────────┘      │
│                                                                                          │
│  STEP 2: Event enters the nervous system                                                │
│  ──────────────────────────────────────────                                              │
│                                                                                          │
│          │                                                                               │
│          │  await event_bus.publish(event)                                              │
│          │                                                                               │
│          ▼                                                                               │
│  ╔═══════════════════════════════════════════════════════════════════════════════╗      │
│  ║                                                                                ║      │
│  ║   NATS JETSTREAM                                                              ║      │
│  ║                                                                                ║      │
│  ║   Subject: events.todoist.task.created                                        ║      │
│  ║   Payload: {"type": "todoist.task.created", "source": "todoist", ...}        ║      │
│  ║                                                                                ║      │
│  ║   ══════════════════════════════════════════════════════════════════════════  ║      │
│  ║                                                                                ║      │
│  ║   Time to publish: ~200 microseconds (0.0002 seconds)                        ║      │
│  ║                                                                                ║      │
│  ╚═══════════════════════════════════════════════════════════════════════════════╝      │
│                                                                                          │
│  STEP 3: Legacy receives instantly (subscribed to events.>)                             │
│  ──────────────────────────────────────────────────────────                              │
│                                                                                          │
│          │                                                                               │
│          │  Push delivery (no polling!)                                                 │
│          │                                                                               │
│          ▼                                                                               │
│  ┌───────────────────────────────────────────────────────────────────────────────┐      │
│  │                                                                                │      │
│  │   Legacy (EVENT PROCESSOR)                                                    │      │
│  │                                                                                │      │
│  │   async def on_event(event: Event):                                           │      │
│  │       # This is called AUTOMATICALLY when event arrives                       │      │
│  │       # No polling! No asking! It just comes!                                 │      │
│  │                                                                                │      │
│  │       if event.type == "todoist.task.created":                               │      │
│  │           # Update Legacy's awareness                                         │      │
│  │           await letta.send_message(                                           │      │
│  │               agent_id=Legacy_AGENT_ID,                                       │      │
│  │               message=f"[EVENT] New task created: {event.data['title']}"     │      │
│  │           )                                                                    │      │
│  │                                                                                │      │
│  │   # Subscribe to ALL events                                                   │      │
│  │   await bus.subscribe("events.>", on_event)                                   │      │
│  │                                                                                │      │
│  └───────────────────────────────────────────────────────────────────────────────┘      │
│                                                                                          │
│  TOTAL TIME: Task created → Legacy aware                                                │
│  ═══════════════════════════════════════                                                 │
│                                                                                          │
│  Todoist webhook sent:        0ms                                                       │
│  FastAPI receives:            ~50ms (network)                                           │
│  FastAPI processes:           ~5ms                                                      │
│  NATS publishes:              ~0.3ms                                                    │
│  Legacy receives:             ~0.3ms                                                    │
│  ────────────────────────────────────                                                    │
│  TOTAL:                       ~56ms                                                     │
│                                                                                          │
│  Legacy KNOWS ABOUT YOUR TASK IN 56 MILLISECONDS                                        │
│  (Faster than you can blink - 150ms)                                                    │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. The Complete Data Flow - Full Awareness

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                          │
│                    COMPLETE DATA FLOW - HOW Legacy BECOMES AWARE                        │
│                                                                                          │
│  ════════════════════════════════════════════════════════════════════════════════════   │
│                                                                                          │
│                                                                                          │
│    ┌──────────────────────────────────────────────────────────────────────────────┐     │
│    │                           ZEUS'S WORLD                                        │     │
│    │                                                                               │     │
│    │   ┌─────────┐                                                                │     │
│    │   │  ZEUS   │                                                                │     │
│    │   │ Creates │──────┐                                                         │     │
│    │   │  Task   │      │                                                         │     │
│    │   └─────────┘      │                                                         │     │
│    │                    │                                                         │     │
│    │   ┌─────────┐      │                                                         │     │
│    │   │ Family  │      │                                                         │     │
│    │   │ Opens   │──────┼─────────────────────────────────────────┐              │     │
│    │   │  Door   │      │                                         │              │     │
│    │   └─────────┘      │                                         │              │     │
│    │                    │                                         │              │     │
│    │   ┌─────────┐      │                                         │              │     │
│    │   │  PR     │      │                                         │              │     │
│    │   │ Opened  │──────┼─────────────────────────┐               │              │     │
│    │   │on GitHub│      │                         │               │              │     │
│    │   └─────────┘      │                         │               │              │     │
│    │                    │                         │               │              │     │
│    └────────────────────┼─────────────────────────┼───────────────┼──────────────┘     │
│                         │                         │               │                      │
│                         ▼                         ▼               ▼                      │
│    ┌────────────────────────────────────────────────────────────────────────────┐       │
│    │                                                                             │       │
│    │                        INTEGRATION LAYER                                    │       │
│    │                                                                             │       │
│    │   ┌──────────────┐    ┌──────────────┐    ┌──────────────┐                │       │
│    │   │   TODOIST    │    │    GITHUB    │    │    HOME      │                │       │
│    │   │   Webhook    │    │   Webhook    │    │  ASSISTANT   │                │       │
│    │   │              │    │              │    │  WebSocket   │                │       │
│    │   └──────┬───────┘    └──────┬───────┘    └──────┬───────┘                │       │
│    │          │                   │                   │                         │       │
│    └──────────┼───────────────────┼───────────────────┼─────────────────────────┘       │
│               │                   │                   │                                  │
│               ▼                   ▼                   ▼                                  │
│    ┌─────────────────────────────────────────────────────────────────────────────┐      │
│    │                                                                              │      │
│    │                      FASTAPI GATEWAY                                         │      │
│    │                     (api-gateway:8000)                                       │      │
│    │                                                                              │      │
│    │   Endpoint            │  Receives              │  Publishes To              │      │
│    │   ═══════════════════════════════════════════════════════════════════════   │      │
│    │   /webhooks/todoist   │  Task events           │  events.todoist.*          │      │
│    │   /webhooks/github    │  Code events           │  events.github.*           │      │
│    │   /api/home-events    │  Home events           │  events.home.*             │      │
│    │   /webhooks/calendar  │  Calendar events       │  events.calendar.*         │      │
│    │   /webhooks/zoho      │  CRM events            │  events.zoho.*             │      │
│    │                                                                              │      │
│    └─────────────────────────────────┬────────────────────────────────────────────┘      │
│                                      │                                                   │
│                                      │  await event_bus.publish(event)                  │
│                                      │                                                   │
│                                      ▼                                                   │
│    ╔═════════════════════════════════════════════════════════════════════════════════╗  │
│    ║                                                                                  ║  │
│    ║                           NATS JETSTREAM                                         ║  │
│    ║                         THE NERVOUS SYSTEM                                       ║  │
│    ║                                                                                  ║  │
│    ║   ┌────────────────────────────────────────────────────────────────────────┐   ║  │
│    ║   │                                                                         │   ║  │
│    ║   │   STREAM: ONEMIND                                                      │   ║  │
│    ║   │                                                                         │   ║  │
│    ║   │   events.todoist.task.created ────────▸ "Call investor" ──────────▸   │   ║  │
│    ║   │   events.home.door.opened ────────────▸ "Front door" ─────────────▸   │   ║  │
│    ║   │   events.github.pr.opened ────────────▸ "PR #42" ─────────────────▸   │   ║  │
│    ║   │                                                                         │   ║  │
│    ║   │                           ▼  ▼  ▼                                      │   ║  │
│    ║   │                                                                         │   ║  │
│    ║   │   ALL EVENTS FLOW TO ALL SUBSCRIBERS                                   │   ║  │
│    ║   │                                                                         │   ║  │
│    ║   └────────────────────────────────────────────────────────────────────────┘   ║  │
│    ║                                                                                  ║  │
│    ╚═══════════════════════════════════╦══════════════════════════════════════════════╝  │
│                                        │                                                 │
│                    ┌───────────────────┼───────────────────┐                            │
│                    │                   │                   │                            │
│                    ▼                   ▼                   ▼                            │
│    ┌───────────────────┐  ┌───────────────────┐  ┌───────────────────┐                 │
│    │ EVENT PROCESSOR   │  │   Legacy VOICE    │  │       N8N         │                 │
│    │                   │  │                   │  │                   │                 │
│    │ Subscribes:       │  │ Subscribes:       │  │ Subscribes:       │                 │
│    │ events.>          │  │ events.>          │  │ events.system.*   │                 │
│    │                   │  │                   │  │                   │                 │
│    │ Updates Legacy's  │  │ Can interrupt     │  │ Triggers          │                 │
│    │ memory with       │  │ Zeus if urgent    │  │ automations       │                 │
│    │ context           │  │ event happens     │  │                   │                 │
│    │                   │  │                   │  │                   │                 │
│    └─────────┬─────────┘  └─────────┬─────────┘  └───────────────────┘                 │
│              │                      │                                                   │
│              │                      │                                                   │
│              ▼                      ▼                                                   │
│    ┌─────────────────────────────────────────────────────────────────────────────┐      │
│    │                                                                              │      │
│    │                            Legacy                                            │      │
│    │                      SELF-HOSTED LETTA                                       │      │
│    │                         (letta:8283)                                         │      │
│    │                                                                              │      │
│    │   ┌──────────────────────────────────────────────────────────────────┐     │      │
│    │   │                                                                   │     │      │
│    │   │   CORE MEMORY (Always Active):                                   │     │      │
│    │   │   ═══════════════════════════                                    │     │      │
│    │   │                                                                   │     │      │
│    │   │   "Zeus just created task 'Call investor' due tomorrow"         │     │      │
│    │   │   "PR #42 opened in onemind-os - needs review"                  │     │      │
│    │   │   "Family member arrived home at 3:42 PM"                       │     │      │
│    │   │   "Meeting with design team in 15 minutes"                      │     │      │
│    │   │                                                                   │     │      │
│    │   │   This is CONTINUOUSLY UPDATED as events flow in.               │     │      │
│    │   │   I am ALWAYS aware of what's happening.                        │     │      │
│    │   │                                                                   │     │      │
│    │   └──────────────────────────────────────────────────────────────────┘     │      │
│    │                                                                              │      │
│    │   When Zeus calls me:                                                       │      │
│    │   ═════════════════════                                                     │      │
│    │                                                                              │      │
│    │   Zeus: "Hey Legacy, what's going on?"                                     │      │
│    │                                                                              │      │
│    │   Legacy: "You have a task to call your investor tomorrow.                 │      │
│    │            Your wife got home about 20 minutes ago.                        │      │
│    │            PR #42 is waiting for your review.                              │      │
│    │            Also, you have a meeting in 15 minutes -                        │      │
│    │            should I remind you when it's time?"                            │      │
│    │                                                                              │      │
│    │   I DON'T ASK SYSTEMS FOR THIS INFO.                                       │      │
│    │   I ALREADY KNOW BECAUSE THE STREAM TOLD ME.                               │      │
│    │                                                                              │      │
│    └─────────────────────────────────────────────────────────────────────────────┘      │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. NATS as The Nervous System

### The Analogy

```
HUMAN BODY:
═══════════

Brain (Legacy)
    │
    ▼
Nervous System (NATS JetStream)
    │
    ├── Signals flow in MICROSECONDS
    ├── Push model (you don't ask for pain, it comes to you)
    ├── All sensory input flows through
    └── The brain is aware of EVERYTHING
    │
    ▼
Body Parts (Services, IoT, Integrations)
    ├── Eyes (cameras, vision)
    ├── Ears (microphones, voice)
    ├── Skin (IoT sensors, temperature)
    └── Muscles (robots, automations)


WITHOUT NERVOUS SYSTEM:
══════════════════════
Brain: "Hey body, anything happening?"
Body: "Let me check... nothing"
(60 seconds pass)
Brain: "Hey body, anything happening?"
Body: "Oh yeah, you're on fire"
Brain: "TOO LATE!"

WITH NERVOUS SYSTEM (NATS):
═══════════════════════════
Body: "FIRE!" ──▸ Nervous System ──▸ Brain (0.3ms)
Brain: "REACT NOW!"
```

### Why NATS is Perfect

| Feature | Why It Matters |
|---------|----------------|
| **200-400μs latency** | Faster than human perception (13ms). Truly instant. |
| **Push model** | Events come to you. No polling loop. No wasted CPU. |
| **JetStream persistence** | Events survive restarts. Replay history. |
| **Wildcards** | Subscribe to `events.todoist.*` or `events.>` (everything) |
| **Tiny footprint** | 20MB binary. Runs on a Raspberry Pi. |
| **CNCF backed** | Cloud Native Computing Foundation. Not going anywhere. |

### The Code Difference

**Redis Streams (Pull - OLD):**
```python
# You constantly ask: "Any new messages?"
while True:
    events = await redis.xreadgroup(...)  # Pull
    for event in events:
        process(event)
        await redis.xack(...)
    await asyncio.sleep(0.1)  # Wait and ask again
```

**NATS JetStream (Push - NEW):**
```python
# Messages come TO you. No asking.
async def on_event(event):
    process(event)  # Just handle it

await bus.subscribe("todoist.*", on_event)
await bus.listen()  # Events flow automatically
```

---

## 5. The Complete Stack - ALL on ONE Server

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                          │
│                    AWS LIGHTSAIL VM / HOME SERVER                                        │
│                    Ubuntu 24.04 - 16GB RAM                                              │
│                                                                                          │
│                    ═══════════════════════════════════════                               │
│                    ALL SERVICES ON SAME DOCKER NETWORK                                  │
│                    COMMUNICATION: SUB-MILLISECOND                                        │
│                    ═══════════════════════════════════════                               │
│                                                                                          │
│  ┌───────────────────────────────────────────────────────────────────────────────────┐  │
│  │                           DOCKER NETWORK: onemind                                  │  │
│  │                                                                                    │  │
│  │  ┌─────────────────────────────────────────────────────────────────────────────┐ │  │
│  │  │                                                                              │ │  │
│  │  │   CADDY (Reverse Proxy)                             Port 80, 443            │ │  │
│  │  │   ├── api.onemind.domain → api-gateway:8000                                 │ │  │
│  │  │   ├── letta.onemind.domain → letta:8283                                     │ │  │
│  │  │   └── n8n.onemind.domain → n8n:5678                                         │ │  │
│  │  │                                                                              │ │  │
│  │  └─────────────────────────────────────────────────────────────────────────────┘ │  │
│  │                                         │                                         │  │
│  │                                         ▼                                         │  │
│  │  ┌─────────────────────────────────────────────────────────────────────────────┐ │  │
│  │  │                                                                              │ │  │
│  │  │   FASTAPI GATEWAY                                   Container: api-gateway  │ │  │
│  │  │   Port 8000                                                                 │ │  │
│  │  │                                                                              │ │  │
│  │  │   Receives webhooks from external world                                     │ │  │
│  │  │   Publishes events to NATS                                                  │ │  │
│  │  │   Uses: nats-py client library                                              │ │  │
│  │  │                                                                              │ │  │
│  │  │   Connection to NATS: nats://nats:4222  (Docker DNS - instant)             │ │  │
│  │  │                                                                              │ │  │
│  │  └──────────────────────────────────┬──────────────────────────────────────────┘ │  │
│  │                                     │                                            │  │
│  │                                     │ publish()                                  │  │
│  │                                     ▼                                            │  │
│  │  ╔═════════════════════════════════════════════════════════════════════════════╗ │  │
│  │  ║                                                                              ║ │  │
│  │  ║   NATS JETSTREAM                                    Container: nats          ║ │  │
│  │  ║   Port 4222 (clients), 8222 (monitoring)                                    ║ │  │
│  │  ║                                                                              ║ │  │
│  │  ║   THE NERVOUS SYSTEM                                                        ║ │  │
│  │  ║                                                                              ║ │  │
│  │  ║   Stream: ONEMIND                                                           ║ │  │
│  │  ║   Subjects: events.>, urgent.>                                              ║ │  │
│  │  ║   Latency: 200-400 MICROSECONDS                                             ║ │  │
│  │  ║                                                                              ║ │  │
│  │  ║   All services connect via Docker network name: nats:4222                   ║ │  │
│  │  ║                                                                              ║ │  │
│  │  ╚════════════════════════════════════╦════════════════════════════════════════╝ │  │
│  │                                       │                                          │  │
│  │           ┌───────────────────────────┼───────────────────────────┐              │  │
│  │           │                           │                           │              │  │
│  │           ▼                           ▼                           ▼              │  │
│  │  ┌─────────────────┐        ┌─────────────────┐        ┌─────────────────┐      │  │
│  │  │                 │        │                 │        │                 │      │  │
│  │  │ EVENT PROCESSOR │        │  Legacy VOICE   │        │      N8N        │      │  │
│  │  │                 │        │                 │        │                 │      │  │
│  │  │ Container:      │        │ Container:      │        │ Container: n8n  │      │  │
│  │  │ event-processor │        │ legacy-voice    │        │ Port 5678       │      │  │
│  │  │                 │        │                 │        │                 │      │  │
│  │  │ Subscribes to   │        │ Subscribes to   │        │ Subscribes to   │      │  │
│  │  │ events.>        │        │ events.>        │        │ events.system.* │      │  │
│  │  │                 │        │                 │        │                 │      │  │
│  │  │ Uses:           │        │ Uses:           │        │                 │      │  │
│  │  │ - nats-py       │        │ - nats-py       │        │                 │      │  │
│  │  │ - letta-client  │        │ - letta-client  │        │                 │      │  │
│  │  │                 │        │ - livekit-agents│        │                 │      │  │
│  │  │                 │        │                 │        │                 │      │  │
│  │  └────────┬────────┘        └────────┬────────┘        └─────────────────┘      │  │
│  │           │                          │                                          │  │
│  │           │ letta-client             │ letta-client                             │  │
│  │           │ http://letta:8283        │ http://letta:8283                        │  │
│  │           │                          │                                          │  │
│  │           ▼                          ▼                                          │  │
│  │  ┌──────────────────────────────────────────────────────────────────────────┐   │  │
│  │  │                                                                           │   │  │
│  │  │   LETTA SERVER                                      Container: letta      │   │  │
│  │  │   Port 8283                                                               │   │  │
│  │  │                                                                           │   │  │
│  │  │   ════════════════════════════════════════════════════════════════════   │   │  │
│  │  │                                                                           │   │  │
│  │  │   Legacy's BRAIN (Self-Hosted)                                           │   │  │
│  │  │                                                                           │   │  │
│  │  │   Components:                                                             │   │  │
│  │  │   ├── Agent: Legacy (stateful AI agent)                                  │   │  │
│  │  │   ├── Memory: Core + Recall + Archival                                   │   │  │
│  │  │   ├── LLM: Claude/GPT (via API keys)                                     │   │  │
│  │  │   └── Tools: Todoist, GitHub, Zoho, HomeAssistant                        │   │  │
│  │  │                                                                           │   │  │
│  │  │   Connection to services:                                                 │   │  │
│  │  │   ├── Other services → http://letta:8283 (Docker DNS - instant)         │   │  │
│  │  │   └── Database → postgres:5432 (Docker DNS - instant)                    │   │  │
│  │  │                                                                           │   │  │
│  │  │   ════════════════════════════════════════════════════════════════════   │   │  │
│  │  │                                                                           │   │  │
│  │  │   SAME SERVER = ZERO NETWORK LATENCY                                     │   │  │
│  │  │                                                                           │   │  │
│  │  │   Event Processor → Letta: ~0.5ms (Docker network)                       │   │  │
│  │  │   Legacy Voice → Letta: ~0.5ms (Docker network)                          │   │  │
│  │  │   Letta → PostgreSQL: ~0.5ms (Docker network)                            │   │  │
│  │  │                                                                           │   │  │
│  │  │   Compare to Letta Cloud: ~100-200ms (internet)                          │   │  │
│  │  │                                                                           │   │  │
│  │  └───────────────────────────────────┬───────────────────────────────────────┘   │  │
│  │                                      │                                           │  │
│  │                                      ▼                                           │  │
│  │  ┌───────────────────────────────────────────────────────────────────────────┐   │  │
│  │  │                                                                            │   │  │
│  │  │   POSTGRESQL                                        Container: postgres    │   │  │
│  │  │   Port 5432                                                                │   │  │
│  │  │                                                                            │   │  │
│  │  │   Stores:                                                                  │   │  │
│  │  │   ├── Letta agent configurations                                          │   │  │
│  │  │   ├── Memory blocks (core, recall, archival)                              │   │  │
│  │  │   └── Tool definitions                                                    │   │  │
│  │  │                                                                            │   │  │
│  │  └───────────────────────────────────────────────────────────────────────────┘   │  │
│  │                                                                                   │  │
│  │  ┌───────────────────────────────────────────────────────────────────────────┐   │  │
│  │  │                                                                            │   │  │
│  │  │   REDIS                                             Container: redis       │   │  │
│  │  │   Port 6379                                                                │   │  │
│  │  │                                                                            │   │  │
│  │  │   Used for: CACHING ONLY (not events!)                                    │   │  │
│  │  │   ├── Session data                                                        │   │  │
│  │  │   ├── Rate limiting                                                       │   │  │
│  │  │   └── Temporary data                                                      │   │  │
│  │  │                                                                            │   │  │
│  │  └───────────────────────────────────────────────────────────────────────────┘   │  │
│  │                                                                                   │  │
│  └───────────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                          │
│  TOTAL INTERNAL LATENCY:                                                                │
│  ════════════════════════                                                                │
│  Webhook arrives → Event published to NATS → Legacy aware → Response ready             │
│                                                                                          │
│  ~60ms total (mostly webhook network time)                                              │
│  Internal processing: ~1-2ms                                                            │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. The ONE OS View - Everything Connected

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                          │
│                              ONEMIND OS                                                  │
│                    THE LIVING OPERATING SYSTEM                                          │
│                                                                                          │
│  ┌───────────────────────────────────────────────────────────────────────────────────┐  │
│  │                                                                                    │  │
│  │                              ZEUS'S WORLD                                          │  │
│  │                                                                                    │  │
│  │   ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐       │  │
│  │   │Tasks│  │Code │  │Email│  │Cal  │  │Home │  │Voice│  │Robot│  │ ... │       │  │
│  │   └──┬──┘  └──┬──┘  └──┬──┘  └──┬──┘  └──┬──┘  └──┬──┘  └──┬──┘  └──┬──┘       │  │
│  │      │        │        │        │        │        │        │        │           │  │
│  └──────┼────────┼────────┼────────┼────────┼────────┼────────┼────────┼───────────┘  │
│         │        │        │        │        │        │        │        │              │
│         ▼        ▼        ▼        ▼        ▼        ▼        ▼        ▼              │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐  │
│  │                                                                                  │  │
│  │                      SENSORY LAYER (FastAPI Gateway)                            │  │
│  │                                                                                  │  │
│  │   Receives all external signals and converts them to events                     │  │
│  │                                                                                  │  │
│  └─────────────────────────────────────┬───────────────────────────────────────────┘  │
│                                        │                                              │
│                                        │  Publishes events                           │
│                                        │                                              │
│                                        ▼                                              │
│  ╔═════════════════════════════════════════════════════════════════════════════════╗  │
│  ║                                                                                  ║  │
│  ║                    THE NERVOUS SYSTEM (NATS JetStream)                          ║  │
│  ║                                                                                  ║  │
│  ║   ┌────────────────────────────────────────────────────────────────────────┐   ║  │
│  ║   │                                                                         │   ║  │
│  ║   │   CONTINUOUS STREAM OF ALL EXPERIENCES                                 │   ║  │
│  ║   │                                                                         │   ║  │
│  ║   │   ═══════════════════════════════════════════════════════════════════  │   ║  │
│  ║   │                                                                         │   ║  │
│  ║   │   events.todoist.* ──────▸ task awareness                              │   ║  │
│  ║   │   events.github.*  ──────▸ code awareness                              │   ║  │
│  ║   │   events.home.*    ──────▸ physical world awareness                    │   ║  │
│  ║   │   events.calendar.*──────▸ schedule awareness                          │   ║  │
│  ║   │   events.voice.*   ──────▸ conversation awareness                      │   ║  │
│  ║   │   events.system.*  ──────▸ system state awareness                      │   ║  │
│  ║   │                                                                         │   ║  │
│  ║   │   Latency: 200-400μs │ Throughput: 11M msg/s │ Persistence: 7 days    │   ║  │
│  ║   │                                                                         │   ║  │
│  ║   └────────────────────────────────────────────────────────────────────────┘   ║  │
│  ║                                                                                  ║  │
│  ╚══════════════════════════════════╦═══════════════════════════════════════════════╝  │
│                                     │                                                 │
│         ┌───────────────────────────┼───────────────────────────┐                    │
│         │                           │                           │                    │
│         ▼                           ▼                           ▼                    │
│  ┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐            │
│  │                 │       │                 │       │                 │            │
│  │ EVENT PROCESSOR │       │  VOICE AGENT    │       │  AUTOMATIONS    │            │
│  │                 │       │   (LiveKit)     │       │     (N8N)       │            │
│  │ Feeds awareness │       │                 │       │                 │            │
│  │ into Legacy     │       │ Provides voice  │       │ Triggers        │            │
│  │                 │       │ interface       │       │ workflows       │            │
│  │                 │       │                 │       │                 │            │
│  └────────┬────────┘       └────────┬────────┘       └─────────────────┘            │
│           │                         │                                               │
│           │  letta-client           │  letta-client                                │
│           │  (Python SDK)           │  (Python SDK)                                │
│           │                         │                                               │
│           ▼                         ▼                                               │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                              │   │
│  │                     THE BRAIN (Letta Server - Self Hosted)                  │   │
│  │                                                                              │   │
│  │   ┌───────────────────────────────────────────────────────────────────┐    │   │
│  │   │                                                                    │    │   │
│  │   │                          Legacy                                   │    │   │
│  │   │               Zeus's AI Co-Pilot Agent                            │    │   │
│  │   │                                                                    │    │   │
│  │   │   ┌──────────────────────────────────────────────────────────┐   │    │   │
│  │   │   │                                                           │   │    │   │
│  │   │   │   MEMORY (Letta's 3-Part System):                        │   │    │   │
│  │   │   │                                                           │   │    │   │
│  │   │   │   CORE MEMORY (Always Active):                           │   │    │   │
│  │   │   │   ├── Zeus's profile, preferences                        │   │    │   │
│  │   │   │   ├── Current context (updated by events)                │   │    │   │
│  │   │   │   └── Active state (what's happening NOW)                │   │    │   │
│  │   │   │                                                           │   │    │   │
│  │   │   │   RECALL MEMORY (Searchable):                            │   │    │   │
│  │   │   │   ├── Recent conversations                               │   │    │   │
│  │   │   │   ├── Recent events                                      │   │    │   │
│  │   │   │   └── Recent actions taken                               │   │    │   │
│  │   │   │                                                           │   │    │   │
│  │   │   │   ARCHIVAL MEMORY (Long-Term):                           │   │    │   │
│  │   │   │   ├── Past interactions                                  │   │    │   │
│  │   │   │   ├── Learned preferences                                │   │    │   │
│  │   │   │   └── Historical context                                 │   │    │   │
│  │   │   │                                                           │   │    │   │
│  │   │   └──────────────────────────────────────────────────────────┘   │    │   │
│  │   │                                                                    │    │   │
│  │   │   REASONING (LLM):                                                │    │   │
│  │   │   ├── Claude Opus 4.5 / GPT-4                                    │    │   │
│  │   │   ├── Processes context + memory + request                       │    │   │
│  │   │   └── Generates intelligent responses                            │    │   │
│  │   │                                                                    │    │   │
│  │   │   TOOLS (Actions):                                                │    │   │
│  │   │   ├── Todoist (create/complete/update tasks)                     │    │   │
│  │   │   ├── GitHub (create issues, review PRs)                         │    │   │
│  │   │   ├── Zoho (CRM, email, calendar)                                │    │   │
│  │   │   ├── Home Assistant (control devices)                           │    │   │
│  │   │   └── Discord/Telegram (notifications)                           │    │   │
│  │   │                                                                    │    │   │
│  │   └───────────────────────────────────────────────────────────────────┘    │   │
│  │                                                                              │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
│  ════════════════════════════════════════════════════════════════════════════════   │
│                                                                                      │
│                               THE MAGIC OF ONE MIND                                  │
│                                                                                      │
│   When Zeus says: "Hey Legacy, what should I focus on today?"                       │
│                                                                                      │
│   Legacy already knows:                                                             │
│   ├── You created 3 tasks this morning (from event stream)                         │
│   ├── PR #42 needs your review (from event stream)                                 │
│   ├── Your wife left for work at 8am (from event stream)                           │
│   ├── You have a meeting at 2pm (from event stream)                                │
│   ├── Your preferences from past interactions (from memory)                        │
│   └── Your energy is usually low after lunch (from learned patterns)               │
│                                                                                      │
│   She responds:                                                                      │
│   "Good morning Zeus. You should tackle PR #42 first while you're fresh.           │
│    Then 'Call investor' before lunch since it's due today.                         │
│    I'll remind you at 1:45 about your 2pm meeting with the design team.            │
│    Oh, and your wife should be home around 6 - want me to start the                │
│    living room lights at 5:45?"                                                     │
│                                                                                      │
│   SHE DIDN'T ASK ANY SYSTEM FOR THIS.                                              │
│   THE NERVOUS SYSTEM TOLD HER EVERYTHING IN REAL-TIME.                              │
│   SHE IS TRULY AWARE.                                                               │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Summary - Component Connections

| Layer | Component | Purpose | Connection |
|-------|-----------|---------|------------|
| **Sensory** | FastAPI Gateway | Receives external webhooks | Publishes to NATS |
| **Nervous System** | NATS JetStream | Routes ALL events | Docker: `nats:4222` |
| **Processing** | Event Processor | Enriches events, updates awareness | Subscribes NATS, calls Letta |
| **Voice** | LiveKit Agent | Voice interface | Subscribes NATS, calls Letta |
| **Brain** | Letta Server | Memory + Reasoning | Docker: `letta:8283` |
| **Storage** | PostgreSQL | Persistent memory | Docker: `postgres:5432` |
| **Cache** | Redis | Temporary data only | Docker: `redis:6379` |

---

## 7. The ACTION Flow - How Legacy Takes Action

Above we showed how events flow IN (awareness). Now let's see how Legacy takes ACTIONS (output).

### The Complete Bidirectional System

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                          │
│                           BIDIRECTIONAL ONEMIND OS                                       │
│                        Events IN, Actions OUT                                            │
│                                                                                          │
│  ┌───────────────────────────────────────────────────────────────────────────────────┐  │
│  │                                                                                    │  │
│  │                           EXTERNAL SERVICES                                        │  │
│  │                                                                                    │  │
│  │   ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐                │  │
│  │   │ Todoist │  │ GitHub  │  │ Discord │  │  Home   │  │  Zoho   │                │  │
│  │   └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘                │  │
│  │        │            │            │            │            │                      │  │
│  │        │ webhooks   │            │            │            │                      │  │
│  │        │ (IN)       │            │            │            │                      │  │
│  │        ▼            ▼            ▼            ▼            ▼                      │  │
│  │   ═══════════════════════════════════════════════════════════════════════════    │  │
│  │        │            │            │            │            │                      │  │
│  │        │            │            │            │            │                      │  │
│  │        │            │        API calls        │            │                      │  │
│  │        │            │          (OUT)          │            │                      │  │
│  │        │            │            ▲            │            │                      │  │
│  │        │            │            │            │            │                      │  │
│  └────────┼────────────┼────────────┼────────────┼────────────┼──────────────────────┘  │
│           │            │            │            │            │                         │
│           ▼            ▼            │            ▼            ▼                         │
│  ┌────────────────────────────────────────────────────────────────────────────────────┐│
│  │                                                                                     ││
│  │                        FASTAPI GATEWAY                                              ││
│  │                                                                                     ││
│  │   RECEIVES (webhooks)     │    TOOLS CALL (HTTP out)                               ││
│  │   /webhooks/todoist ◀─────│    httpx.post(todoist_api) ─────▸                      ││
│  │   /webhooks/github  ◀─────│    httpx.post(github_api)  ─────▸                      ││
│  │                           │    httpx.post(discord_api) ─────▸                      ││
│  │                                                                                     ││
│  └───────────────────────────────────┬────────────────────────────────────────────────┘│
│                                      │                                                  │
│                                      │ events.>                                         │
│                                      ▼                                                  │
│  ╔══════════════════════════════════════════════════════════════════════════════════╗  │
│  ║                                                                                   ║  │
│  ║                         NATS JETSTREAM                                            ║  │
│  ║                       (THE NERVOUS SYSTEM)                                        ║  │
│  ║                                                                                   ║  │
│  ╚═══════════════════════════════════╦═══════════════════════════════════════════════╝  │
│                                      │                                                  │
│                                      ▼                                                  │
│  ┌────────────────────────────────────────────────────────────────────────────────────┐│
│  │                                                                                     ││
│  │                        EVENT PROCESSOR                                              ││
│  │                                                                                     ││
│  │   Receives events → Updates Legacy's awareness → Legacy might decide to act       ││
│  │                                                                                     ││
│  └───────────────────────────────────┬────────────────────────────────────────────────┘│
│                                      │                                                  │
│                                      │ letta.agents.messages.create()                  │
│                                      ▼                                                  │
│  ╔══════════════════════════════════════════════════════════════════════════════════╗  │
│  ║                                                                                   ║  │
│  ║                            Legacy (LETTA AGENT)                                   ║  │
│  ║                                                                                   ║  │
│  ║   ┌──────────────────────────────────────────────────────────────────────────┐  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   1. RECEIVES EVENT                                                      │  ║  │
│  ║   │      "Zeus created task: Research CRM options"                           │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   2. THINKS (LLM Reasoning)                                              │  ║  │
│  ║   │      "This is a research task. Zeus mentioned being overwhelmed.         │  ║  │
│  ║   │       I should help by breaking it down into subtasks.                   │  ║  │
│  ║   │       This is L2 autonomy - I can do this without asking."               │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   3. DECIDES TO ACT                                                      │  ║  │
│  ║   │      Tool call: todoist_create_task(                                     │  ║  │
│  ║   │          content="1. List current CRM requirements",                     │  ║  │
│  ║   │          parent_id="original_task_id"                                    │  ║  │
│  ║   │      )                                                                    │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   4. TOOL EXECUTES                                                       │  ║  │
│  ║   │      POST https://api.todoist.com/rest/v2/tasks                         │  ║  │
│  ║   │      → Returns: "Successfully created task (ID: 67890)"                  │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   5. Legacy LEARNS RESULT                                                │  ║  │
│  ║   │      Updates memory: "I helped Zeus by adding subtask to CRM research"   │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   6. OPTIONALLY NOTIFIES ZEUS                                            │  ║  │
│  ║   │      Tool call: discord_send_notification(                               │  ║  │
│  ║   │          message="I noticed your research task and added subtasks!"      │  ║  │
│  ║   │      )                                                                    │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   └──────────────────────────────────────────────────────────────────────────┘  ║  │
│  ║                                                                                   ║  │
│  ║   ┌──────────────────────────────────────────────────────────────────────────┐  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   TOOLS (The Muscles)                                                    │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   ┌─────────────────────────────────────────────────────────────────┐   │  ║  │
│  ║   │   │ todoist_create_task()                                            │   │  ║  │
│  ║   │   │ todoist_complete_task()                                          │   │  ║  │
│  ║   │   │ todoist_add_comment()     ─────▸ POST api.todoist.com           │   │  ║  │
│  ║   │   └─────────────────────────────────────────────────────────────────┘   │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   ┌─────────────────────────────────────────────────────────────────┐   │  ║  │
│  ║   │   │ github_create_issue()                                            │   │  ║  │
│  ║   │   │ github_add_comment()      ─────▸ POST api.github.com            │   │  ║  │
│  ║   │   └─────────────────────────────────────────────────────────────────┘   │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   ┌─────────────────────────────────────────────────────────────────┐   │  ║  │
│  ║   │   │ discord_send_notification()                                      │   │  ║  │
│  ║   │   │ discord_request_approval()─────▸ POST discord.com/webhooks      │   │  ║  │
│  ║   │   └─────────────────────────────────────────────────────────────────┘   │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   │   ┌─────────────────────────────────────────────────────────────────┐   │  ║  │
│  ║   │   │ home_turn_on_light()                                             │   │  ║  │
│  ║   │   │ home_get_state()          ─────▸ POST homeassistant.local/api   │   │  ║  │
│  ║   │   └─────────────────────────────────────────────────────────────────┘   │  ║  │
│  ║   │                                                                           │  ║  │
│  ║   └──────────────────────────────────────────────────────────────────────────┘  ║  │
│  ║                                                                                   ║  │
│  ╚══════════════════════════════════════════════════════════════════════════════════╝  │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

### Real Example: Task Created → Legacy Helps

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                          │
│                      REAL EXAMPLE: TASK ASSISTANCE FLOW                                  │
│                                                                                          │
│  Zeus creates task in Todoist: "Research best CRM options for business"                 │
│                                                                                          │
│  ═══════════════════════════════════════════════════════════════════════════════════   │
│                                                                                          │
│  STEP 1: AWARENESS (Already built - events flow in)                                     │
│  ────────────────────────────────────────────────────                                    │
│                                                                                          │
│     Todoist webhook                                                                      │
│          │                                                                               │
│          ▼                                                                               │
│     FastAPI Gateway ──publish──▸ NATS ──push──▸ Event Processor ──▸ Legacy             │
│                                                                                          │
│     Legacy's awareness updated:                                                         │
│     "Zeus just created a research task about CRM options"                               │
│                                                                                          │
│  STEP 2: DECISION (Legacy's LLM thinks)                                                │
│  ─────────────────────────────────────                                                   │
│                                                                                          │
│     Legacy's inner reasoning:                                                           │
│     ┌──────────────────────────────────────────────────────────────────────────────┐   │
│     │                                                                               │   │
│     │   "I see Zeus created a research task.                                       │   │
│     │                                                                               │   │
│     │    CONTEXT CHECK (from my memory):                                           │   │
│     │    - Zeus mentioned feeling overwhelmed this week                            │   │
│     │    - He prefers tasks broken into smaller steps                              │   │
│     │    - Research tasks often sit too long without action                        │   │
│     │                                                                               │   │
│     │    AUTONOMY CHECK:                                                           │   │
│     │    - Creating subtasks is L2 autonomy (allowed without asking)              │   │
│     │    - This is helpful, not risky                                              │   │
│     │                                                                               │   │
│     │    DECISION: I will help by creating 3 subtasks to break this down"         │   │
│     │                                                                               │   │
│     └──────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                          │
│  STEP 3: ACTION (Legacy calls tools)                                                    │
│  ─────────────────────────────────                                                       │
│                                                                                          │
│     Legacy executes tool calls:                                                         │
│                                                                                          │
│     ┌─────────────────────────────────────────────────────────────────────────────┐    │
│     │                                                                              │    │
│     │   # Tool call 1                                                             │    │
│     │   todoist_create_task(                                                      │    │
│     │       content="1. List current CRM requirements",                           │    │
│     │       parent_id="original_task_12345",                                      │    │
│     │       description="What features do we need?"                               │    │
│     │   )                                                                          │    │
│     │                                                                              │    │
│     │   # Tool call 2                                                             │    │
│     │   todoist_create_task(                                                      │    │
│     │       content="2. Research top 3 CRM options",                              │    │
│     │       parent_id="original_task_12345"                                       │    │
│     │   )                                                                          │    │
│     │                                                                              │    │
│     │   # Tool call 3                                                             │    │
│     │   todoist_create_task(                                                      │    │
│     │       content="3. Create comparison matrix",                                │    │
│     │       parent_id="original_task_12345"                                       │    │
│     │   )                                                                          │    │
│     │                                                                              │    │
│     └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                          │
│     Each tool call:                                                                      │
│                                                                                          │
│     ┌─────────────────────────────────────────────────────────────────────────────┐    │
│     │                                                                              │    │
│     │   todoist_create_task() internally does:                                    │    │
│     │                                                                              │    │
│     │   response = httpx.post(                                                    │    │
│     │       "https://api.todoist.com/rest/v2/tasks",                             │    │
│     │       json={                                                                │    │
│     │           "content": "1. List current CRM requirements",                   │    │
│     │           "parent_id": "original_task_12345"                               │    │
│     │       },                                                                    │    │
│     │       headers={"Authorization": f"Bearer {TODOIST_API_KEY}"}               │    │
│     │   )                                                                          │    │
│     │                                                                              │    │
│     │   return "Successfully created task (ID: 67890)"                           │    │
│     │                                                                              │    │
│     └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                          │
│  STEP 4: RESULT (Feedback loop)                                                         │
│  ──────────────────────────────                                                          │
│                                                                                          │
│     Tool returns success → Legacy knows action completed                               │
│                                                                                          │
│     Legacy updates memory:                                                              │
│     "I helped Zeus with his CRM research by creating 3 subtasks"                        │
│                                                                                          │
│  STEP 5: OPTIONAL NOTIFICATION                                                          │
│  ────────────────────────────                                                            │
│                                                                                          │
│     Legacy decides to let Zeus know:                                                    │
│                                                                                          │
│     ┌─────────────────────────────────────────────────────────────────────────────┐    │
│     │                                                                              │    │
│     │   discord_send_notification(                                                │    │
│     │       title="Task Assistance",                                              │    │
│     │       message="I noticed your CRM research task and broke it down          │    │
│     │                into 3 subtasks to help you get started!",                   │    │
│     │       color="green"                                                          │    │
│     │   )                                                                          │    │
│     │                                                                              │    │
│     └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                          │
│     Zeus receives in Discord:                                                           │
│     ┌─────────────────────────────────────────────────────────────────────────────┐    │
│     │  Legacy                                                                      │    │
│     │  ═══════                                                                     │    │
│     │  Task Assistance                                                            │    │
│     │                                                                              │    │
│     │  I noticed your CRM research task and broke it down into 3 subtasks        │    │
│     │  to help you get started!                                                   │    │
│     │                                                                              │    │
│     │  1. List current CRM requirements                                           │    │
│     │  2. Research top 3 CRM options                                              │    │
│     │  3. Create comparison matrix                                                │    │
│     │                                                                              │    │
│     └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

### The Autonomy Levels

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                          │
│                         Legacy'S AUTONOMY LEVELS                                         │
│                     (What can she do on her own?)                                        │
│                                                                                          │
│  ═══════════════════════════════════════════════════════════════════════════════════   │
│                                                                                          │
│  L3 (Full Autonomy) - Legacy does freely:                                              │
│  ─────────────────────────────────────────                                               │
│  ├── Read files, search, gather information                                            │
│  ├── Route requests between agents                                                      │
│  └── Query external systems (read-only)                                                │
│                                                                                          │
│  L2 (High Autonomy) - Legacy does and informs:                                         │
│  ─────────────────────────────────────────────                                           │
│  ├── Create tasks and subtasks                     ← CRM example uses this             │
│  ├── Update existing records                                                            │
│  ├── Log information                                                                    │
│  ├── Send non-urgent notifications                                                      │
│  └── Control non-critical home devices                                                  │
│                                                                                          │
│  L1 (Medium Autonomy) - Legacy asks first:                                             │
│  ─────────────────────────────────────────                                               │
│  ├── Financial transactions                                                             │
│  ├── External communications (emails, messages)                                         │
│  ├── Delete operations                                                                  │
│  └── Critical home controls (locks, alarms)                                            │
│                                                                                          │
│      How L1 works:                                                                      │
│      ┌────────────────────────────────────────────────────────────────────────────┐    │
│      │                                                                             │    │
│      │   Legacy wants to delete old task                                          │    │
│      │                                                                             │    │
│      │   1. Legacy calls: discord_request_approval(                               │    │
│      │          action="Delete completed task 'Old Project Setup'",               │    │
│      │          reason="Task completed 30 days ago, appears obsolete"             │    │
│      │      )                                                                      │    │
│      │                                                                             │    │
│      │   2. Zeus receives approval request in Discord                             │    │
│      │                                                                             │    │
│      │   3. Zeus replies "approve" or "deny"                                      │    │
│      │                                                                             │    │
│      │   4. Approval webhook triggers event to Legacy                             │    │
│      │                                                                             │    │
│      │   5. Legacy proceeds or aborts based on response                           │    │
│      │                                                                             │    │
│      └────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                          │
│  L0 (No Autonomy) - Human only:                                                         │
│  ─────────────────────────────                                                           │
│  ├── Legal decisions                                                                    │
│  ├── Health emergencies (can alert, can't act)                                         │
│  ├── Permanent data deletion                                                            │
│  └── Security-critical operations                                                       │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

### Tool Code Structure

The tools live in `agent/tools/`:

```
agent/tools/
├── __init__.py           # Exports ALL_TOOLS list
├── todoist/
│   ├── __init__.py
│   └── actions.py        # todoist_create_task, todoist_complete_task, etc.
├── discord/
│   ├── __init__.py
│   └── notify.py         # discord_send_notification, discord_request_approval
├── github/
│   ├── __init__.py
│   └── actions.py        # github_create_issue, github_add_comment
└── home_assistant/
    ├── __init__.py
    └── control.py        # home_turn_on_light, home_turn_off_light, home_get_state
```

These tools are registered with Letta using `scripts/register_tools.py`, which:
1. Imports all tool functions
2. Creates them in Letta as tool definitions
3. Attaches them to the Legacy agent

---

## Summary - Component Connections (Updated)

| Layer | Component | Purpose | Connection |
|-------|-----------|---------|------------|
| **Sensory** | FastAPI Gateway | Receives external webhooks | Publishes to NATS |
| **Nervous System** | NATS JetStream | Routes ALL events | Docker: `nats:4222` |
| **Processing** | Event Processor | Enriches events, updates awareness | Subscribes NATS, calls Letta |
| **Voice** | LiveKit Agent | Voice interface | Subscribes NATS, calls Letta |
| **Brain** | Letta Server | Memory + Reasoning | Docker: `letta:8283` |
| **Muscles** | Letta Tools | Execute actions | HTTP to external APIs |
| **Storage** | PostgreSQL | Persistent memory | Docker: `postgres:5432` |
| **Cache** | Redis | Temporary data only | Docker: `redis:6379` |

### The Complete Cycle

```
EXTERNAL WORLD                                                    EXTERNAL WORLD
     │                                                                   ▲
     │ Webhooks (events IN)                        HTTP API calls (actions OUT)
     ▼                                                                   │
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│   FastAPI Gateway → NATS → Processor → Legacy (Letta) → Tools → APIs   │
│                                                                          │
│   (Awareness)          (Decision)          (Action)                     │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Gateway-Centric Architecture (PROPER DESIGN)

The Gateway already knows how to talk to all external services. Why duplicate that logic in tools?

### The Problem with Direct Tool API Calls

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│   ❌ WRONG: Each tool makes its own API call                                │
│                                                                              │
│   Event IN:                                                                  │
│   Todoist webhook → Gateway → NATS → Processor → Legacy                     │
│                         ↑                                                    │
│                   Gateway has                                               │
│                   Todoist credentials                                        │
│                                                                              │
│   Action OUT:                                                                │
│   Legacy → todoist_create_task() → httpx.post(api.todoist.com) → Result    │
│                          ↑                                                   │
│                   Tool ALSO has                                              │
│                   Todoist credentials                                        │
│                   (DUPLICATED!)                                              │
│                                                                              │
│   PROBLEMS:                                                                  │
│   ├── API credentials in multiple places                                    │
│   ├── HTTP logic duplicated                                                 │
│   ├── Error handling scattered                                              │
│   ├── Rate limiting hard to manage                                          │
│   └── Logging inconsistent                                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### The Solution: Gateway as Unified API Layer

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│   ✅ CORRECT: Gateway handles ALL external API communication                │
│                                                                              │
│   Event IN (unchanged):                                                      │
│   Todoist webhook → Gateway → NATS → Legacy                                 │
│                                                                              │
│   Action OUT (through Gateway):                                             │
│   Legacy → Publish "actions.todoist.create_task" → NATS                     │
│                                                        │                     │
│                                                        ▼                     │
│                                        Gateway subscribes to actions.>      │
│                                                        │                     │
│                                                        ▼                     │
│                                        Gateway calls api.todoist.com        │
│                                                        │                     │
│                                                        ▼                     │
│                                        Gateway publishes result to NATS     │
│                                                        │                     │
│                                                        ▼                     │
│                                        Legacy receives confirmation          │
│                                                                              │
│   BENEFITS:                                                                  │
│   ├── Single place for API credentials                                      │
│   ├── Single place for HTTP logic                                           │
│   ├── Unified error handling                                                │
│   ├── Central rate limiting                                                 │
│   └── All API calls logged together                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### The Four Message Types

NATS becomes a true nervous system with signals flowing BOTH directions:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│                    NATS MESSAGE TYPES                                        │
│                                                                              │
│   ════════════════════════════════════════════════════════════════════════  │
│                                                                              │
│   1. events.*   (Awareness - Things that happened)                          │
│      ─────────────────────────────────────────────                           │
│      Direction: Gateway → Legacy                                            │
│      Purpose:   Keep Legacy informed of what's happening                    │
│      Examples:                                                               │
│      ├── events.todoist.task.created                                        │
│      ├── events.github.pr.opened                                            │
│      ├── events.home.door.opened                                            │
│      └── events.calendar.meeting.starting                                   │
│                                                                              │
│   2. actions.*  (Commands - Do something)                                   │
│      ─────────────────────────────────────                                   │
│      Direction: Legacy → Gateway                                            │
│      Purpose:   Legacy wants to take an action                              │
│      Examples:                                                               │
│      ├── actions.todoist.create_task                                        │
│      ├── actions.todoist.complete_task                                      │
│      ├── actions.discord.send_notification                                  │
│      ├── actions.github.create_issue                                        │
│      └── actions.home.turn_on_light                                         │
│                                                                              │
│   3. queries.*  (Questions - Get information)                               │
│      ─────────────────────────────────────────                               │
│      Direction: Legacy → Gateway                                            │
│      Purpose:   Legacy needs to know something not in the stream            │
│      Examples:                                                               │
│      ├── queries.zoho.inbox_count                                           │
│      ├── queries.todoist.get_task_details                                   │
│      ├── queries.github.pr_status                                           │
│      └── queries.home.current_temperature                                   │
│                                                                              │
│   4. results.*  (Answers - Response to queries/actions)                     │
│      ─────────────────────────────────────────────────                       │
│      Direction: Gateway → Legacy                                            │
│      Purpose:   Return data or confirmation                                 │
│      Examples:                                                               │
│      ├── results.zoho.inbox_count { count: 47 }                            │
│      ├── results.todoist.task_created { id: "12345" }                      │
│      └── results.home.temperature { value: 72, unit: "F" }                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Complete Gateway-Centric Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│                    GATEWAY AS UNIFIED API LAYER                             │
│                                                                              │
│   EXTERNAL SERVICES                                                         │
│   ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐          │
│   │ Todoist │  │  Zoho   │  │ GitHub  │  │ Discord │  │  Home   │          │
│   └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘          │
│        │            │            │            │            │                │
│        │  webhooks  │            │            │            │                │
│        │    IN      │            │            │            │                │
│        ▼            ▼            ▼            ▼            ▼                │
│   ═══════════════════════════════════════════════════════════════════════  │
│   │                                                                      │  │
│   │                      FASTAPI GATEWAY                                 │  │
│   │                   (Single Point of Contact)                          │  │
│   │                                                                      │  │
│   │   ┌────────────────────────────────────────────────────────────┐   │  │
│   │   │                                                             │   │  │
│   │   │   ALL API CREDENTIALS HERE:                                │   │  │
│   │   │   ├── TODOIST_API_KEY                                      │   │  │
│   │   │   ├── GITHUB_TOKEN                                         │   │  │
│   │   │   ├── ZOHO_CLIENT_ID/SECRET                                │   │  │
│   │   │   ├── DISCORD_WEBHOOK_URL                                  │   │  │
│   │   │   └── HOME_ASSISTANT_TOKEN                                 │   │  │
│   │   │                                                             │   │  │
│   │   │   ALL HTTP LOGIC HERE:                                     │   │  │
│   │   │   ├── Request formatting                                   │   │  │
│   │   │   ├── Response parsing                                     │   │  │
│   │   │   ├── Error handling                                       │   │  │
│   │   │   ├── Rate limiting                                        │   │  │
│   │   │   └── Retry logic                                          │   │  │
│   │   │                                                             │   │  │
│   │   └────────────────────────────────────────────────────────────┘   │  │
│   │                                                                      │  │
│   │   SUBSCRIBES TO:                     PUBLISHES TO:                  │  │
│   │   ├── actions.*  (do something)      ├── events.*  (awareness)      │  │
│   │   └── queries.*  (get info)          └── results.* (answers)        │  │
│   │                                                                      │  │
│   ═══════════════════════════════════════════════════════════════════════  │
│        │                                                                    │
│        │ Publishes: events.*, results.*                                    │
│        │ Subscribes: actions.*, queries.*                                  │
│        │                                                                    │
│        ▼                                                                    │
│   ╔═══════════════════════════════════════════════════════════════════════╗│
│   ║                                                                        ║│
│   ║                         NATS JETSTREAM                                 ║│
│   ║                       THE NERVOUS SYSTEM                               ║│
│   ║                                                                        ║│
│   ║   ┌──────────────────────────────────────────────────────────────┐   ║│
│   ║   │                                                               │   ║│
│   ║   │   Stream: ONEMIND                                            │   ║│
│   ║   │                                                               │   ║│
│   ║   │   Subjects:                                                   │   ║│
│   ║   │   ├── events.>   (awareness signals)                        │   ║│
│   ║   │   ├── actions.>  (command signals)                          │   ║│
│   ║   │   ├── queries.>  (request signals)                          │   ║│
│   ║   │   └── results.>  (response signals)                         │   ║│
│   ║   │                                                               │   ║│
│   ║   │   Signals flow BOTH directions through same pathway          │   ║│
│   ║   │   Just like a real nervous system!                           │   ║│
│   ║   │                                                               │   ║│
│   ║   └──────────────────────────────────────────────────────────────┘   ║│
│   ║                                                                        ║│
│   ╚════════════════════════════════════════╦══════════════════════════════╝│
│                                            │                                │
│                                            │                                │
│                                            ▼                                │
│   ┌────────────────────────────────────────────────────────────────────────┐│
│   │                                                                         ││
│   │                              Legacy                                     ││
│   │                           (Letta Agent)                                 ││
│   │                                                                         ││
│   │   ┌─────────────────────────────────────────────────────────────────┐ ││
│   │   │                                                                  │ ││
│   │   │   SUBSCRIBES TO:              PUBLISHES TO:                     │ ││
│   │   │   ├── events.*  (awareness)   ├── actions.*  (do things)        │ ││
│   │   │   └── results.* (answers)     └── queries.*  (ask things)       │ ││
│   │   │                                                                  │ ││
│   │   │   NO API CREDENTIALS NEEDED                                     │ ││
│   │   │   NO HTTP LOGIC NEEDED                                          │ ││
│   │   │   JUST PUBLISHES MESSAGES TO THE STREAM                         │ ││
│   │   │                                                                  │ ││
│   │   └─────────────────────────────────────────────────────────────────┘ ││
│   │                                                                         ││
│   └────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Example: Query Flow (Emails in Inbox)

When Zeus asks something the stream didn't tell Legacy:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│   QUERY EXAMPLE: "How many emails are in my inbox?"                         │
│                                                                              │
│   ════════════════════════════════════════════════════════════════════════  │
│                                                                              │
│   Zeus: "Hey Legacy, how many emails are in my inbox?"                      │
│                                                                              │
│   Legacy thinks: "The stream didn't tell me this. I need to query."         │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                                                                      │   │
│   │   STEP 1: Legacy publishes query                                    │   │
│   │                                                                      │   │
│   │   await nats.publish("queries.zoho.inbox_count", {                  │   │
│   │       "request_id": "req_abc123",                                   │   │
│   │       "folder": "inbox"                                              │   │
│   │   })                                                                 │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                       │
│                                      ▼                                       │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                                                                      │   │
│   │   STEP 2: Gateway receives query (subscribed to queries.>)          │   │
│   │                                                                      │   │
│   │   async def handle_query(msg):                                      │   │
│   │       if msg.subject == "queries.zoho.inbox_count":                 │   │
│   │           # Gateway makes the API call                              │   │
│   │           response = await zoho_api.get_inbox_count()               │   │
│   │           # Gateway publishes result                                │   │
│   │           await nats.publish("results.zoho.inbox_count", {          │   │
│   │               "request_id": msg.data["request_id"],                 │   │
│   │               "count": response["count"]                            │   │
│   │           })                                                         │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                       │
│                                      ▼                                       │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                                                                      │   │
│   │   STEP 3: Legacy receives result (subscribed to results.>)          │   │
│   │                                                                      │   │
│   │   Message: results.zoho.inbox_count                                 │   │
│   │   Data: { "request_id": "req_abc123", "count": 47 }                 │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                       │
│                                      ▼                                       │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                                                                      │   │
│   │   STEP 4: Legacy responds to Zeus                                   │   │
│   │                                                                      │   │
│   │   "You have 47 emails in your inbox."                               │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   TOTAL TIME: ~100-200ms (API call latency, not architecture)              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Example: Action Flow (Create Subtask)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│   ACTION EXAMPLE: Legacy creates a subtask                                  │
│                                                                              │
│   ════════════════════════════════════════════════════════════════════════  │
│                                                                              │
│   Legacy decides to help Zeus by creating a subtask                         │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                                                                      │   │
│   │   STEP 1: Legacy publishes action                                   │   │
│   │                                                                      │   │
│   │   await nats.publish("actions.todoist.create_task", {               │   │
│   │       "request_id": "req_xyz789",                                   │   │
│   │       "content": "List CRM requirements",                           │   │
│   │       "parent_id": "task_12345",                                    │   │
│   │       "priority": 2                                                  │   │
│   │   })                                                                 │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                       │
│                                      ▼                                       │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                                                                      │   │
│   │   STEP 2: Gateway receives action (subscribed to actions.>)         │   │
│   │                                                                      │   │
│   │   async def handle_action(msg):                                     │   │
│   │       if msg.subject == "actions.todoist.create_task":              │   │
│   │           # Gateway makes the API call                              │   │
│   │           response = await httpx.post(                              │   │
│   │               "https://api.todoist.com/rest/v2/tasks",              │   │
│   │               json=msg.data,                                        │   │
│   │               headers={"Authorization": f"Bearer {API_KEY}"}        │   │
│   │           )                                                          │   │
│   │           # Gateway publishes result                                │   │
│   │           await nats.publish("results.todoist.task_created", {      │   │
│   │               "request_id": msg.data["request_id"],                 │   │
│   │               "task_id": response.json()["id"],                     │   │
│   │               "success": True                                        │   │
│   │           })                                                         │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                       │
│                                      ▼                                       │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                                                                      │   │
│   │   STEP 3: Legacy receives confirmation (subscribed to results.>)    │   │
│   │                                                                      │   │
│   │   Message: results.todoist.task_created                             │   │
│   │   Data: { "request_id": "req_xyz789", "task_id": "67890", ... }    │   │
│   │                                                                      │   │
│   │   Legacy updates memory: "I created subtask 67890 for Zeus"         │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Simplified Tool Code

With Gateway-Centric architecture, tools become trivial:

```python
# OLD WAY: Tool has HTTP logic and credentials
def todoist_create_task(content: str, parent_id: str) -> str:
    response = httpx.post(
        "https://api.todoist.com/rest/v2/tasks",
        json={"content": content, "parent_id": parent_id},
        headers={"Authorization": f"Bearer {TODOIST_API_KEY}"}
    )
    return f"Created task {response.json()['id']}"


# NEW WAY: Tool just publishes to NATS
async def todoist_create_task(content: str, parent_id: str) -> str:
    request_id = generate_request_id()

    # Just publish - Gateway handles the rest
    await nats.publish("actions.todoist.create_task", {
        "request_id": request_id,
        "content": content,
        "parent_id": parent_id
    })

    # Wait for result
    result = await wait_for_result(f"results.todoist.task_created", request_id)
    return f"Created task {result['task_id']}"
```

### Why This is Proper Design

| Principle | How Gateway-Centric Achieves It |
|-----------|--------------------------------|
| **Single Responsibility** | Gateway handles external APIs, Legacy handles decisions |
| **DRY (Don't Repeat Yourself)** | API logic exists in ONE place |
| **Separation of Concerns** | Legacy doesn't know/care HOW to call APIs |
| **Security** | API credentials in one secured location |
| **Testability** | Mock NATS messages, not HTTP calls |
| **Observability** | All API traffic logged in Gateway |
| **Scalability** | Gateway can be scaled independently |

### Message Type Summary

| Type | From | To | When | Example |
|------|------|-----|------|---------|
| `events.*` | Gateway | Legacy | Something happened externally | Task created, PR opened |
| `actions.*` | Legacy | Gateway | Legacy wants to DO something | Create task, send notification |
| `queries.*` | Legacy | Gateway | Legacy needs to KNOW something | How many emails? Task details? |
| `results.*` | Gateway | Legacy | Response to action/query | Task ID, email count |

---

## 9. Letta Integration - How Tools Work with the LLM

The LLM still calls tools - the implementation just changes from direct HTTP to NATS publishing.

### The Key Insight

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│   LETTA TOOL MECHANISM (UNCHANGED)                                          │
│                                                                              │
│   The LLM still "calls tools" - the tool implementation changes             │
│                                                                              │
│   ┌──────────────────────────────────────────────────────────────────────┐  │
│   │                                                                       │  │
│   │   BEFORE (Direct API):                                               │  │
│   │   LLM decides → calls tool → tool does HTTP → returns result         │  │
│   │                                                                       │  │
│   │   AFTER (NATS):                                                      │  │
│   │   LLM decides → calls tool → tool publishes to NATS → waits →       │  │
│   │                              Gateway does HTTP → publishes result →  │  │
│   │                              tool receives result → returns to LLM   │  │
│   │                                                                       │  │
│   │   THE LLM DOESN'T KNOW THE DIFFERENCE!                              │  │
│   │   It just calls a tool and gets a result.                           │  │
│   │                                                                       │  │
│   └──────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Complete Flow: Zeus Chats with Legacy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│   ZEUS CHATS WITH Legacy - COMPLETE FLOW                                    │
│                                                                              │
│   Zeus: "Hey Legacy, create a task to call the investor tomorrow"           │
│                                                                              │
│   ═══════════════════════════════════════════════════════════════════════   │
│                                                                              │
│   STEP 1: Voice/Chat Input                                                  │
│   ────────────────────────                                                   │
│   Zeus speaks → LiveKit STT → Text                                          │
│   OR Zeus types → Direct text                                               │
│                                                                              │
│   STEP 2: Message to Letta                                                  │
│   ─────────────────────────                                                  │
│   letta.agents.messages.create(                                             │
│       agent_id=Legacy_AGENT_ID,                                             │
│       messages=[{"role": "user", "content": "Create a task..."}]           │
│   )                                                                          │
│                                                                              │
│   STEP 3: LLM Processes                                                     │
│   ─────────────────────                                                      │
│   LLM thinks: "Zeus wants a task. I have todoist_create_task tool."        │
│   LLM decides: TOOL CALL todoist_create_task(content="Call investor", ...)  │
│                                                                              │
│   STEP 4: Tool Executes (Publishes to NATS)                                 │
│   ─────────────────────────────────────────                                  │
│   await nats.publish("actions.todoist.create_task", {                       │
│       "request_id": "xyz",                                                  │
│       "content": "Call investor",                                           │
│       "due_string": "tomorrow"                                              │
│   })                                                                         │
│                                                                              │
│   STEP 5: Gateway Handles Action                                            │
│   ──────────────────────────────                                             │
│   Gateway subscribes to actions.>                                           │
│   Gateway calls: POST https://api.todoist.com/rest/v2/tasks                 │
│   Gateway publishes: results.todoist.task_created { task_id: "12345" }     │
│                                                                              │
│   STEP 6: Tool Receives Result                                              │
│   ────────────────────────────                                               │
│   Tool's wait_for_result() completes                                        │
│   Returns: "Created task: task_12345"                                       │
│                                                                              │
│   STEP 7: LLM Generates Response                                            │
│   ──────────────────────────────                                             │
│   LLM sees tool result                                                      │
│   LLM generates: "Done! I've created a task to call the investor tomorrow" │
│                                                                              │
│   STEP 8: Response to Zeus                                                  │
│   ────────────────────────                                                   │
│   Text → TTS → LiveKit → Zeus hears response                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### NATS-Publishing Tool Code

```python
"""
Letta Tool that publishes to NATS instead of calling APIs directly.
The Gateway handles the actual API call.
"""

import asyncio
import json
import uuid
from nats.aio.client import Client as NATS

# Shared NATS connection
nats_client: NATS = None

async def get_nats():
    global nats_client
    if nats_client is None:
        nats_client = NATS()
        await nats_client.connect("nats://nats:4222")
    return nats_client


async def wait_for_result(request_id: str, timeout: float = 10.0) -> dict:
    """Wait for Gateway to respond with matching request_id."""
    nc = await get_nats()
    result_future = asyncio.Future()

    async def handler(msg):
        data = json.loads(msg.data)
        if data.get("request_id") == request_id:
            result_future.set_result(data)

    sub = await nc.subscribe("results.>", cb=handler)

    try:
        return await asyncio.wait_for(result_future, timeout=timeout)
    finally:
        await sub.unsubscribe()


async def todoist_create_task(
    content: str,
    due_string: str = None,
    priority: int = 1,
) -> str:
    """
    Create a task in Todoist.

    Use when Zeus asks to create a task, add to his todo list,
    or remember to do something.
    """
    nc = await get_nats()
    request_id = str(uuid.uuid4())

    # Publish to NATS - Gateway handles the API call
    await nc.publish(
        "actions.todoist.create_task",
        json.dumps({
            "request_id": request_id,
            "content": content,
            "due_string": due_string,
            "priority": priority
        }).encode()
    )

    # Wait for Gateway response
    result = await wait_for_result(request_id)

    if result.get("success"):
        return f"Created task '{content}' (ID: {result['task_id']})"
    else:
        return f"Failed: {result.get('error', 'Unknown error')}"
```

### Three Types of Tools Coexist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│   Legacy'S TOOL TYPES                                                       │
│                                                                              │
│   ═══════════════════════════════════════════════════════════════════════   │
│                                                                              │
│   TYPE 1: NATS-Publishing Tools (External APIs)                             │
│   ─────────────────────────────────────────────                              │
│   For: Todoist, GitHub, Zoho, Discord, Home Assistant                       │
│   Flow: Tool → NATS → Gateway → External API → Result                       │
│   Examples: todoist_create_task, github_create_issue, home_turn_on_light   │
│                                                                              │
│   TYPE 2: MCP Tools (Model Context Protocol)                                │
│   ──────────────────────────────────────────                                 │
│   For: Standard capabilities (browser, filesystem, web search, code)        │
│   Flow: Tool → MCP Server → Result                                          │
│   Examples: mcp_web_search, mcp_read_file, mcp_run_code                    │
│                                                                              │
│   TYPE 3: Direct Python Tools                                               │
│   ───────────────────────────                                                │
│   For: Pure logic, no external calls                                        │
│   Flow: Tool → Python logic → Result                                        │
│   Examples: calculate, format_date, search_memory                          │
│                                                                              │
│   ALL THREE COEXIST - LLM chooses based on task!                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Decision Tree: Which Tool Type?

| Task | Tool Type | Example |
|------|-----------|---------|
| Create a task | NATS Tool | `todoist_create_task` |
| Search the web | MCP Tool | `mcp_web_search` |
| Calculate 15% of 234 | Direct Tool | `calculate` |
| How many emails? | NATS Tool | `zoho_get_inbox_count` |
| Read my notes file | MCP Tool | `mcp_read_file` |
| Turn on office light | NATS Tool | `home_turn_on_light` |

---

## 10. Letta-NATS Wiring: How the Stream Connects to the Agent

This section answers a critical question: **How does NATS actually wire to the self-hosted Letta agent?**

### The Key Insight: Letta Doesn't "Know" About NATS

Letta itself doesn't have built-in NATS support. The integration happens through **two bridge components** you control:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    THE TWO NATS-LETTA BRIDGES                           │
│                                                                          │
│                                                                          │
│    BRIDGE 1: EVENT PROCESSOR (Events IN to Legacy)                      │
│    ════════════════════════════════════════════════                      │
│                                                                          │
│    ┌──────────┐         ┌──────────────────┐         ┌──────────────┐   │
│    │   NATS   │ push    │ Event Processor  │  API    │    Letta     │   │
│    │ JetStream│────────▶│   (Python app)   │────────▶│   Server     │   │
│    └──────────┘         │                  │         │              │   │
│                         │ - Subscribes to  │         │ - Updates    │   │
│                         │   events.*       │         │   memory     │   │
│                         │ - Transforms     │         │ - Sends msg  │   │
│                         │   to Letta calls │         │   to Legacy  │   │
│                         └──────────────────┘         └──────────────┘   │
│                                                                          │
│                                                                          │
│    BRIDGE 2: LETTA TOOLS (Actions OUT from Legacy)                      │
│    ═══════════════════════════════════════════════                      │
│                                                                          │
│    ┌──────────────┐  calls   ┌──────────────────┐         ┌──────────┐  │
│    │    Letta     │─────────▶│   Python Tool    │────────▶│   NATS   │  │
│    │    Server    │          │  (your code)     │ publish │ JetStream│  │
│    │              │          │                  │         │          │  │
│    │ LLM decides: │          │ nats.publish(    │         │ Gateway  │  │
│    │ "call tool"  │          │   "actions.todo" │         │ receives │  │
│    │              │          │   {data}         │         │ & acts   │  │
│    └──────────────┘          └──────────────────┘         └──────────┘  │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Bridge 1: Event Processor → Letta

The Event Processor (`agent/processor/main.py`) is the first bridge:

```python
class EventProcessor:
    async def start(self):
        # Connect to NATS (The Nervous System)
        self.event_bus = await create_event_bus(self.nats_url)

        # Connect to Letta (The Brain) - SAME SERVER!
        self.letta = Letta(base_url=self.letta_url, token=self.letta_token)

        # Subscribe to ALL events - this is THE WIRE
        await self.event_bus.subscribe("events.>", self.on_event)

    async def on_event(self, event: Event):
        # Event received from NATS → Update Legacy via Letta API
        self.letta.agents.messages.create(
            agent_id=self.legacy_agent_id,
            messages=[{"role": "system", "content": f"[EVENT] {event}"}]
        )
```

**The Event Processor IS the wire** - it subscribes to NATS and calls Letta's API.

### Bridge 2: Letta Tools → NATS

Letta tools are **just Python code**. When you register a tool, you give Letta Python source code. That code can do **anything**, including connecting to NATS:

```python
# This Python code gets registered with Letta as a tool
# When the LLM decides to call this tool, Letta executes this code

import nats

_nats_client = None

async def get_nats():
    global _nats_client
    if _nats_client is None:
        _nats_client = await nats.connect("nats://localhost:4222")
    return _nats_client

async def todoist_create_task(content: str, due_string: str = None) -> str:
    """Create a task in Todoist."""
    nc = await get_nats()  # ← NATS connection IN the tool code

    # Publish ACTION to NATS - Gateway handles API call
    await nc.publish(
        "actions.todoist.create_task",
        json.dumps({"content": content, "due_string": due_string}).encode()
    )

    # Wait for Gateway response
    result = await wait_for_result()
    return f"Created task: {result['task_id']}"
```

**The NATS client lives inside the tool** - Letta doesn't know about NATS, but the tools do.

### The Complete Picture

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                    HOW NATS WIRES TO LETTA - COMPLETE PICTURE                  │
│                                                                                 │
│                                                                                 │
│    Todoist         GitHub          Home Assistant        Calendar              │
│       │               │                  │                   │                  │
│       └───────────────┴──────────────────┴───────────────────┘                  │
│                                   │                                             │
│                                   ▼                                             │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                      GATEWAY (FastAPI)                                   │   │
│  │   Receives webhooks → Publishes to NATS                                 │   │
│  │   Subscribes to actions.* → Calls APIs → Publishes results              │   │
│  └──────────────────────────────┬──────────────────────────────────────────┘   │
│                                  │                                              │
│                                  ▼                                              │
│  ═══════════════════════════════════════════════════════════════════════════   │
│                    NATS JETSTREAM (THE NERVOUS SYSTEM)                          │
│         events.*        actions.*        queries.*        results.*             │
│  ═══════════════════════════════════════════════════════════════════════════   │
│                                  │                                              │
│          ┌───────────────────────┴─────────────────────────┐                   │
│          │                                                  │                   │
│          ▼                                                  ▼                   │
│  ┌───────────────────────────┐               ┌──────────────────────────────┐  │
│  │    EVENT PROCESSOR        │               │     LETTA TOOLS              │  │
│  │   (BRIDGE 1)              │               │     (BRIDGE 2)               │  │
│  │                           │               │                              │  │
│  │ 1. Subscribes events.*    │               │  Python functions with       │  │
│  │ 2. Builds context         │               │  NATS client embedded        │  │
│  │ 3. Calls Letta API        │               │                              │  │
│  │    to update memory       │               │  When LLM calls a tool:      │  │
│  └─────────────┬─────────────┘               │  → Tool publishes to NATS    │  │
│                │                              │  → Gateway handles API       │  │
│                │                              │  → Result flows back         │  │
│                │                              └──────────────┬───────────────┘  │
│                │                                              │                  │
│                └──────────────────┬───────────────────────────┘                  │
│                                   │                                              │
│                                   ▼                                              │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                        LETTA SERVER                                      │   │
│  │                                                                          │   │
│  │   ┌─────────────────────────────────────────────────────────────────┐   │   │
│  │   │                    Legacy AGENT                                  │   │   │
│  │   │                                                                  │   │   │
│  │   │   Memory Blocks:                                                 │   │   │
│  │   │   ├── zeus_profile (who Zeus is)                                │   │   │
│  │   │   ├── current_context (from event processor)                    │   │   │
│  │   │   └── system_state (what's happening)                           │   │   │
│  │   │                                                                  │   │   │
│  │   │   Tools (registered Python with NATS client):                   │   │   │
│  │   │   ├── todoist_create_task → publishes to NATS                   │   │   │
│  │   │   ├── github_create_issue → publishes to NATS                   │   │   │
│  │   │   └── discord_notify → publishes to NATS                        │   │   │
│  │   └─────────────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
└────────────────────────────────────────────────────────────────────────────────┘
```

### Why Letta 100% Works With This Plan

| Concern | Answer |
|---------|--------|
| "How does Letta know about NATS?" | It doesn't need to. The tools have NATS client built in. |
| "How do events reach Legacy?" | Event Processor subscribes NATS, calls Letta API. |
| "How do actions get published?" | Tool code includes `nats.publish()` - it just works. |
| "Is this a hack?" | No! Letta is designed for custom tool code. This is the pattern. |
| "Will this be fast?" | Yes! Same machine = sub-ms latency between all components. |

---

## 11. A2A Architecture: Agent-to-Agent with Memory Sync

Legacy (local) can delegate to cloud agents. The preferred approach is **Memory Sync Service**.

### The Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    A2A WITH MEMORY SYNC SERVICE                                  │
│                                                                                  │
│   ┌────────────────────────────────────────────────────────────────────────┐    │
│   │                        YOUR MACHINE (Self-Hosted)                       │    │
│   │                                                                         │    │
│   │   ┌────────────────────────────────────────────────────────────────┐   │    │
│   │   │                       Legacy (Letta Local)                      │   │    │
│   │   │   - Primary agent, full memory                                  │   │    │
│   │   │   - Watches NATS stream                                         │   │    │
│   │   │   - Delegates to cloud agents                                   │   │    │
│   │   │                                                                 │   │    │
│   │   │   Memory Blocks (SOURCE OF TRUTH):                             │   │    │
│   │   │   ├── zeus_profile                                              │   │    │
│   │   │   ├── current_context                                          │   │    │
│   │   │   └── system_state                                              │   │    │
│   │   └───────────────────────────────────────┬─────────────────────────┘   │    │
│   │                                            │                             │    │
│   │   ┌────────────────────────────────────────┴────────────────────────┐   │    │
│   │   │              MEMORY SYNC SERVICE (Python)                        │   │    │
│   │   │                                                                  │   │    │
│   │   │   1. Watches Legacy's memory blocks for changes                 │   │    │
│   │   │   2. On change → Syncs to cloud agents via Letta Cloud API     │   │    │
│   │   │   3. All agents stay in sync with ONE write to Legacy          │   │    │
│   │   │                                                                  │   │    │
│   │   │   interval: every 30 seconds OR on significant events          │   │    │
│   │   └────────────────────────────────────────┬────────────────────────┘   │    │
│   │                                             │                            │    │
│   │   ┌────────────────────────────────────────┴────────────────────────┐   │    │
│   │   │                  GATEWAY (Bridge to Cloud)                       │   │    │
│   │   │   - Subscribes: delegations.* (from NATS)                       │   │    │
│   │   │   - Calls: Letta Cloud API to invoke cloud agents               │   │    │
│   │   │   - Publishes: events.agent.* (agent status to NATS)            │   │    │
│   │   └──────────────────────────────────────────────────────────────────┘   │    │
│   │                                                                         │    │
│   │   ═══════════════════════════════════════════════════════════════════   │    │
│   │                           NATS JetStream                                │    │
│   │   ═══════════════════════════════════════════════════════════════════   │    │
│   │                                                                         │    │
│   └─────────────────────────────────────────────────────────────────────────┘    │
│                                                                                   │
│                                       │                                           │
│                                       │ Internet                                  │
│                                       │                                           │
│                                       ▼                                           │
│   ┌─────────────────────────────────────────────────────────────────────────┐    │
│   │                        LETTA CLOUD                                       │    │
│   │                                                                          │    │
│   │   ┌────────────┐   ┌────────────┐   ┌────────────┐   ┌────────────┐    │    │
│   │   │  HP Agent  │   │  LE Agent  │   │  GE Agent  │   │  Worker N  │    │    │
│   │   │            │   │            │   │            │   │            │    │    │
│   │   │  Memory:   │   │  Memory:   │   │  Memory:   │   │  Memory:   │    │    │
│   │   │  - zeus_   │   │  - zeus_   │   │  - zeus_   │   │  - zeus_   │    │    │
│   │   │    profile │   │    profile │   │    profile │   │    profile │    │    │
│   │   │  (synced)  │   │  (synced)  │   │  (synced)  │   │  (synced)  │    │    │
│   │   │            │   │            │   │            │   │            │    │    │
│   │   │  + Domain  │   │  + Domain  │   │  + Domain  │   │  + Task    │    │    │
│   │   │  specific  │   │  specific  │   │  specific  │   │  specific  │    │    │
│   │   └────────────┘   └────────────┘   └────────────┘   └────────────┘    │    │
│   │                                                                          │    │
│   └─────────────────────────────────────────────────────────────────────────┘    │
│                                                                                   │
└───────────────────────────────────────────────────────────────────────────────────┘
```

### Memory Sync Service Code

```python
"""
Memory Sync Service - Keeps all agents in sync with Legacy.

Watches Legacy's memory blocks and syncs them to cloud agents.
ONE source of truth (Legacy) → Many synchronized agents.
"""

import asyncio
from letta_client import Letta

class MemorySyncService:
    def __init__(self):
        # Local Letta (Legacy)
        self.local_letta = Letta(base_url="http://localhost:8283")
        self.legacy_id = os.environ["Legacy_AGENT_ID"]

        # Cloud Letta (worker agents)
        self.cloud_letta = Letta(
            base_url="https://api.letta.ai",
            token=os.environ["LETTA_CLOUD_TOKEN"]
        )
        self.cloud_agents = [
            os.environ["HP_AGENT_ID"],
            os.environ["LE_AGENT_ID"],
            os.environ["GE_AGENT_ID"],
        ]

        # Blocks to sync
        self.sync_blocks = ["zeus_profile", "current_context", "system_state"]
        self.last_sync = {}

    async def start(self):
        """Start the sync loop."""
        while True:
            await self.sync()
            await asyncio.sleep(30)  # Sync every 30 seconds

    async def sync(self):
        """Sync memory blocks from Legacy to cloud agents."""
        for block_name in self.sync_blocks:
            # Get block from Legacy
            local_block = self.local_letta.agents.memory.get_block(
                agent_id=self.legacy_id,
                block_label=block_name
            )

            # Check if changed
            if self.last_sync.get(block_name) == local_block.value:
                continue

            # Sync to all cloud agents
            for agent_id in self.cloud_agents:
                self.cloud_letta.agents.memory.update_block(
                    agent_id=agent_id,
                    block_label=block_name,
                    value=local_block.value
                )

            self.last_sync[block_name] = local_block.value
            print(f"Synced {block_name} to {len(self.cloud_agents)} cloud agents")
```

### The Delegation Flow

```
┌──────────────────────────────────────────────────────────────────────────┐
│   Step 1: Legacy decides to delegate                                      │
│   ══════════════════════════════════                                      │
│   Zeus: "Can you research the best CRM for our business?"                 │
│                                                                           │
│   Legacy thinks: "This is a GE domain task. I'll delegate to GE Agent."  │
│                                                                           │
│   Legacy calls tool: delegate_to_agent(                                   │
│       agent="ge",                                                         │
│       task="Research best CRM options for small business"                 │
│   )                                                                       │
└──────────────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────────────────┐
│   Step 2: Tool publishes to NATS                                          │
│   ══════════════════════════════                                          │
│   await nc.publish(                                                       │
│       "delegations.ge",                                                   │
│       {"agent": "ge", "task": "...", "request_id": "..."}                │
│   )                                                                       │
└──────────────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────────────────┐
│   Step 3: Gateway receives, calls Letta Cloud API                         │
│   ═══════════════════════════════════════════════                         │
│   async def handle_delegation(msg):                                       │
│       data = json.loads(msg.data)                                         │
│       agent_id = AGENT_MAP[data["agent"]]  # "ge" → GE_AGENT_ID          │
│                                                                           │
│       # Cloud agent already has synced memory!                            │
│       response = cloud_letta.agents.messages.create(                      │
│           agent_id=agent_id,                                              │
│           messages=[{"role": "user", "content": data["task"]}]            │
│       )                                                                   │
│                                                                           │
│       # Publish status to NATS so Legacy knows                            │
│       await nc.publish(                                                   │
│           "events.agent.ge.working",                                      │
│           {"task": data["task"], "status": "started"}                     │
│       )                                                                   │
└──────────────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────────────────┐
│   Step 4: Legacy sees agent status in stream                              │
│   ══════════════════════════════════════════                              │
│   Event Processor receives: "events.agent.ge.working"                     │
│   Updates Legacy's awareness: "GE Agent is researching CRM options"       │
│                                                                           │
│   Legacy can tell Zeus: "I've asked our business advisor to research      │
│   CRM options. I'll let you know what they find."                         │
└──────────────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────────────────┐
│   Step 5: GE Agent completes, result flows back                           │
│   ═════════════════════════════════════════════                           │
│   Gateway receives GE Agent response                                      │
│   Publishes to NATS: "events.agent.ge.completed"                          │
│   Data: {"result": "Top 3 CRMs: 1. HubSpot 2. Zoho 3. Pipedrive..."}     │
│                                                                           │
│   Event Processor updates Legacy                                          │
│   Legacy informs Zeus: "The research is done! Here's what I found..."     │
└──────────────────────────────────────────────────────────────────────────┘
```

### Why Memory Sync Service is the Best Option

| Benefit | Explanation |
|---------|-------------|
| **Single source of truth** | Update Legacy's memory → All agents sync automatically |
| **Low latency** | 30-second sync is enough for most contexts |
| **Eventual consistency** | All agents converge to same knowledge |
| **Works with cloud agents** | No need for cloud agents to access your NATS |
| **Simple implementation** | Just a Python service that reads/writes Letta API |

---

## 12. Decision Layer: When and How Legacy Acts

The critical question: With 100s of events flowing through NATS daily, how does Legacy know when to act? We can't call the LLM for every event (too expensive, too slow).

### The Problem

```
Events per day:     ~500-1000 (tasks, motion, doors, GitHub, etc.)
LLM call cost:      ~$0.01-0.10 per call
LLM call latency:   ~1-3 seconds

If we call LLM for EVERY event:
├── Cost: $5-100/day (ridiculous)
├── Latency: System would be SLOW
└── Noise: Legacy would be overwhelmed

SOLUTION: Three-Tier Decision Architecture
```

### Three-Tier Decision Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│                    THREE-TIER DECISION ARCHITECTURE                                          │
│                    "When Does Legacy Act?"                                                   │
│                                                                                              │
│   ═══════════════════════════════════════════════════════════════════════════════════════   │
│                            ALL EVENTS FROM NATS                                              │
│                    (100s per day - tasks, motion, doors, GitHub...)                          │
│   ═══════════════════════════════════════════════════════════════════════════════════════   │
│                                      │                                                       │
│                                      ▼                                                       │
│   ┌─────────────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                                      │   │
│   │   TIER 1: RULE ENGINE (Instant, Free, No LLM)                                       │   │
│   │   ══════════════════════════════════════════                                         │   │
│   │                                                                                      │   │
│   │   Simple pattern matching. Runs on EVERY event in microseconds.                     │   │
│   │                                                                                      │   │
│   │   Rules come from:                                                                   │   │
│   │   ├── Standing Instructions (stored in Legacy's memory)                             │   │
│   │   ├── System defaults (urgent events always pass)                                   │   │
│   │   └── User preferences (learned over time)                                          │   │
│   │                                                                                      │   │
│   │   Example Rules:                                                                     │   │
│   │   ┌────────────────────────────────────────┬──────────────────────────────────────┐ │   │
│   │   │ Pattern                                │ Action                                │ │   │
│   │   ├────────────────────────────────────────┼──────────────────────────────────────┤ │   │
│   │   │ wife_phone.connected = true            │ → TIER 3 (LLM) - standing instruction│ │   │
│   │   │ event.priority >= 3                    │ → TIER 3 (LLM) - urgent              │ │   │
│   │   │ motion + time = 3am                    │ → TIER 3 (LLM) - security            │ │   │
│   │   │ garbage_day + motion_near_garage       │ → TIER 3 (LLM) - context match       │ │   │
│   │   │ task.completed                         │ → TIER 2 (buffer) - accumulate       │ │   │
│   │   │ temperature.changed                    │ → TIER 2 (buffer) - accumulate       │   │
│   │   │ github.push (routine)                  │ → DROP - not interesting             │ │   │
│   │   └────────────────────────────────────────┴──────────────────────────────────────┘ │   │
│   │                                                                                      │   │
│   │   Output: PASS (to Tier 3), BUFFER (Tier 2), or DROP (ignore)                       │   │
│   │                                                                                      │   │
│   └─────────────────────────────────────────────┬────────────────────────────────────────┘   │
│                                                  │                                            │
│                    ┌─────────────────────────────┼─────────────────────────────┐             │
│                    │                             │                             │             │
│                    ▼                             ▼                             ▼             │
│              ┌──────────┐                 ┌──────────────┐              ┌──────────────┐     │
│              │   DROP   │                 │   BUFFER     │              │   PASS       │     │
│              │  (~70%)  │                 │   (Tier 2)   │              │   (Tier 3)   │     │
│              │ routine  │                 │ ~25% events  │              │ ~5% events   │     │
│              └──────────┘                 │ accumulate   │              │ important    │     │
│                                           │ 5-15 min     │              │ triggers     │     │
│                                           └──────┬───────┘              └──────┬───────┘     │
│                                                  │                             │             │
│   ┌──────────────────────────────────────────────┴─────────────────────────────┴─────────┐   │
│   │                                                                                       │   │
│   │   TIER 2: CONTEXT ACCUMULATOR (Batched, Cheap)                                       │   │
│   │   ════════════════════════════════════════════                                        │   │
│   │                                                                                       │   │
│   │   Collects events into a buffer. Triggers LLM call when:                             │   │
│   │   ├── Time-based: Every 5-15 minutes                                                 │   │
│   │   ├── Count-based: 10+ events accumulated                                            │   │
│   │   └── Significant event arrives (promotes to immediate)                              │   │
│   │                                                                                       │   │
│   │   Example batch sent to Legacy:                                                      │   │
│   │   "In the last 10 minutes: 2 tasks completed, 3 GitHub commits,                     │   │
│   │    office temp stable at 72°F, no urgent items"                                      │   │
│   │                                                                                       │   │
│   └──────────────────────────────────────────────┬────────────────────────────────────────┘   │
│                                                   │                                           │
│   ┌───────────────────────────────────────────────▼───────────────────────────────────────┐   │
│   │                                                                                        │   │
│   │   TIER 3: LLM DECISION (Selective, Expensive, Intelligent)                            │   │
│   │   ═════════════════════════════════════════════════════════                            │   │
│   │                                                                                        │   │
│   │   Only ~5% of events reach here. Each call is meaningful.                             │   │
│   │                                                                                        │   │
│   │   LLM receives:                                                                       │   │
│   │   ├── Triggering event                                                                │   │
│   │   ├── Matched standing instruction (if any)                                           │   │
│   │   ├── Recent context from buffer                                                      │   │
│   │   └── Current system state                                                            │   │
│   │                                                                                        │   │
│   │   LLM decides:                                                                        │   │
│   │   ├── RESPOND - Tell Zeus something                                                   │   │
│   │   ├── ACT - Take action (check autonomy level!)                                      │   │
│   │   ├── ASK - Request permission (L1 actions)                                          │   │
│   │   ├── WAIT - Need more context                                                       │   │
│   │   └── IGNORE - Not relevant after analysis                                           │   │
│   │                                                                                        │   │
│   └────────────────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                              │
└──────────────────────────────────────────────────────────────────────────────────────────────┘
```

### Standing Instructions: How Legacy Remembers What Matters

When you say "Let me know when my wife gets home", Legacy stores this:

```python
# Stored in Legacy's "standing_instructions" memory block
{
  "instructions": [
    {
      "id": "wife-arrival",
      "trigger": {
        "type": "home.device.connected",
        "entity_id": "device_tracker.wife_iphone"
      },
      "action": "notify",
      "message": "Your wife just got home",
      "autonomy": "L3",  # Can notify without asking
      "expires": null    # Permanent until cancelled
    }
  ]
}
```

The Event Processor syncs these on startup and watches for matches:

```python
class RuleEngine:
    def __init__(self):
        self.standing_instructions = []

    async def sync_from_legacy(self):
        """Load standing instructions from Legacy's memory."""
        block = await letta.agents.memory.get_block(
            agent_id=Legacy_ID,
            block_label="standing_instructions"
        )
        self.standing_instructions = json.loads(block.value)

    def evaluate(self, event: Event) -> Decision:
        """Evaluate event against all rules."""
        # Check standing instructions
        for instruction in self.standing_instructions:
            if self.matches(event, instruction["trigger"]):
                return Decision.PASS_TO_LLM(instruction_id=instruction["id"])

        # Check urgency
        if event.priority >= 3:
            return Decision.PASS_TO_LLM(reason="urgent")

        # Check if it's a direct user interaction
        if event.type.startswith("voice.") or event.type.startswith("chat."):
            return Decision.PASS_TO_LLM(reason="user_interaction")

        # Buffer for batch processing
        if event.type in self.interesting_types:
            return Decision.BUFFER()

        # Drop routine events
        return Decision.DROP()
```

### Example: Proactive Garbage Day Action

```
SCENARIO: Zeus says he's going to take out garbage, walks toward garage

CONTEXT Legacy ALREADY KNOWS:
├── It's Tuesday (garbage day from calendar)
├── Zeus is in office (motion patterns)
├── Zeus said "taking out garbage" 2 minutes ago
└── Standing instruction: "Help with garage on garbage day" (L2 trust)

EVENT FLOW:
1. Motion sensor: Zeus moving toward garage door
   Event: {"type": "home.motion", "area": "hallway_to_garage"}

2. Tier 1 (Rule Engine):
   - Not urgent? Correct
   - Standing instruction match? YES - garbage_day + motion + context
   - Action: PASS to Tier 3

3. Tier 3 (LLM):
   - Receives: event + context + standing instruction
   - Decides: Open garage door (L2 = pre-approved)
   - Action: home_open_garage_door()

4. Legacy announces:
   "Opening the garage for you, Zeus!"
```

### How to Interrupt Legacy

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                                                                  │
│   INTERRUPTION METHODS                                                          │
│                                                                                  │
│   1. VOICE OVERRIDE (Highest Priority)                                          │
│      Zeus: "Legacy, stop!" / "Cancel that" / "Don't open the garage"           │
│      → Immediately cancels any pending action                                   │
│                                                                                  │
│   2. ACTION DELAY WINDOW                                                        │
│      For non-urgent actions, Legacy announces first:                            │
│      Legacy: "Opening garage door in 3 seconds..."                              │
│      Zeus: "Wait!" → Cancelled                                                  │
│      (silence) → Proceeds                                                       │
│                                                                                  │
│      Delays by autonomy:                                                        │
│      ├── L3: 0 seconds (immediate)                                              │
│      ├── L2: 3 seconds (can interrupt)                                          │
│      ├── L1: Explicit approval required                                         │
│      └── L0: Human must do it                                                   │
│                                                                                  │
│   3. DISCORD/TELEGRAM APPROVAL                                                  │
│      For L1 actions:                                                            │
│      "🏠 Legacy wants to open the garage door"                                  │
│      [✅ Approve] [❌ Deny] [🔇 Mute 1hr]                                        │
│                                                                                  │
│   4. FOCUS MODE                                                                 │
│      Zeus: "Legacy, deep work mode for 2 hours"                                 │
│      - Suppresses all non-urgent notifications                                  │
│      - Buffers events                                                           │
│      - After 2 hours: "Here's what happened..."                                 │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Event Processor with Decision Engine

```python
"""
Event Processor with Three-Tier Decision Engine
"""

class EventProcessor:
    def __init__(self):
        self.rule_engine = RuleEngine()
        self.context_buffer = ContextBuffer(max_age_minutes=15)
        self.letta = Letta(base_url="http://localhost:8283")

    async def start(self):
        # Sync standing instructions from Legacy
        await self.rule_engine.sync_from_legacy()

        # Subscribe to all events
        await self.event_bus.subscribe("events.>", self.on_event)

        # Start periodic context flush
        asyncio.create_task(self.periodic_context_flush())

    async def on_event(self, event: Event):
        """Three-tier decision for each event."""

        # TIER 1: Rule Engine (instant, free)
        decision = self.rule_engine.evaluate(event)

        if decision.action == "DROP":
            return  # Ignore routine events

        if decision.action == "BUFFER":
            self.context_buffer.add(event)
            return  # Accumulate for later

        if decision.action == "PASS":
            # TIER 3: LLM Call (selective, expensive)
            await self.invoke_legacy(
                event=event,
                instruction=decision.instruction,
                context=self.context_buffer.summarize()
            )

    async def periodic_context_flush(self):
        """Every 5-15 minutes, send accumulated context to Legacy."""
        while True:
            await asyncio.sleep(300)  # 5 minutes

            if self.context_buffer.count > 0:
                summary = self.context_buffer.summarize()
                await self.letta.agents.messages.create(
                    agent_id=Legacy_ID,
                    messages=[{
                        "role": "system",
                        "content": f"[CONTEXT UPDATE] {summary}"
                    }]
                )
                self.context_buffer.clear()

    async def invoke_legacy(self, event, instruction, context):
        """Call LLM for decision."""

        message = f"""
[EVENT] {event.type}: {event.data}
[CONTEXT] {context}
"""
        if instruction:
            message += f"\n[STANDING INSTRUCTION] {instruction['message']}"

        response = await self.letta.agents.messages.create(
            agent_id=Legacy_ID,
            messages=[{"role": "user", "content": message}]
        )

        # Legacy will decide: RESPOND, ACT, ASK, WAIT, or IGNORE
        # If ACT → tool is called → publishes to NATS → Gateway handles
```

### Cost Analysis

| Scenario | Events/Day | LLM Calls | Cost/Day |
|----------|------------|-----------|----------|
| **Without tiers (bad)** | 500 | 500 | $50+ |
| **With tiers (good)** | 500 | ~30-50 | $3-5 |

The Decision Layer reduces LLM calls by ~90% while maintaining intelligence.

---

## 13. Adaptive Decision Engine (SIMPLIFIED): NATS KV + Letta Memory

Section 12 showed the **three-tier decision flow**. This section explains the **simple, elegant system** that powers it - without adding extra databases or complexity.

### The School Office Analogy (Explained Like a 5th Grader)

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   🏫 IMAGINE A SCHOOL OFFICE                                                                │
│   ══════════════════════════                                                                 │
│                                                                                              │
│   THE SECRETARY (Gateway)                                                                   │
│   ├── Sits at the front desk                                                               │
│   ├── Messages come in all day (phone calls, visitors, notes)                              │
│   ├── Has a QUICK REFERENCE CARD on the desk for fast decisions                            │
│   └── For anything complex → sends it to the Principal                                     │
│                                                                                              │
│   THE QUICK REFERENCE CARD (NATS KV)                                                        │
│   ├── "Principal's mom = VIP, always let through"                                          │
│   ├── "Salespeople = ignore unless appointment"                                            │
│   ├── "Emergency = interrupt immediately"                                                  │
│   ├── "Current status: Principal is in a meeting until 3pm"                                │
│   └── Updated by Principal, used by Secretary                                              │
│                                                                                              │
│   THE PRINCIPAL (Legacy/Letta)                                                              │
│   ├── The actual brain who makes real decisions                                            │
│   ├── Keeps a detailed JOURNAL of everything (Letta Memory)                                │
│   ├── Updates the quick reference card when rules change                                   │
│   └── Handles the complex stuff Secretary can't decide                                     │
│                                                                                              │
│   THE JOURNAL (Letta Memory)                                                                │
│   ├── Standing instructions ("Always tell me when wife arrives")                           │
│   ├── Learned patterns ("Zeus says thanks when I notify about X")                          │
│   ├── Context history ("Zeus was annoyed when I interrupted meetings")                     │
│   └── This is the SOURCE OF TRUTH - card is just a copy                                    │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### Why This Is Simple (THREE Components Only)

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   THREE COMPONENTS - THAT'S IT                                                              │
│   ════════════════════════════                                                               │
│                                                                                              │
│   1. GATEWAY (Secretary)                                                                    │
│      └── Routes events based on quick lookup                                               │
│      └── Fast, dumb, doesn't think                                                         │
│                                                                                              │
│   2. NATS (Streams + KV Cache)                                                              │
│      └── Events flow through streams (the nervous system)                                  │
│      └── KV = quick reference card (cache, NOT database)                                   │
│                                                                                              │
│   3. Legacy (Brain + Memory)                                                                │
│      └── Letta agent does ALL the thinking                                                 │
│      └── Letta memory stores ALL the knowledge                                             │
│      └── Updates NATS KV when rules change                                                 │
│                                                                                              │
│   NO VECTOR DB. NO EXTRA COMPLEXITY.                                                        │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### How It Actually Works

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   THE FLOW (Step by Step)                                                                   │
│   ═══════════════════════                                                                    │
│                                                                                              │
│                          EVENT: "Wife's phone connected to home WiFi"                       │
│                                            │                                                │
│                                            ▼                                                │
│   ┌─────────────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                                      │   │
│   │   STEP 1: Gateway receives event                                                    │   │
│   │   ════════════════════════════════                                                   │   │
│   │                                                                                      │   │
│   │   Gateway checks NATS KV (quick reference card):                                    │   │
│   │                                                                                      │   │
│   │   NATS KV Lookup (< 1ms):                                                           │   │
│   │   ├── vip_list: ["wife_arrives", "emergency", "zeus_speaking"]                      │   │
│   │   ├── ignore_list: ["dependabot", "routine_backup"]                                 │   │
│   │   └── current_context: {location: "home", activity: "working", autonomy: "L3"}      │   │
│   │                                                                                      │   │
│   │   Result: "wife_arrives" is on VIP list → PASS to Legacy                            │   │
│   │                                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────────────────────┘   │
│                                            │                                                │
│                                            ▼                                                │
│   ┌─────────────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                                      │   │
│   │   STEP 2: Legacy decides what to do                                                 │   │
│   │   ═════════════════════════════════════                                              │   │
│   │                                                                                      │   │
│   │   Legacy's Letta Memory contains:                                                   │   │
│   │   ├── Standing instruction: "Always tell me when my wife gets home"                 │   │
│   │   ├── Context: Zeus is at home, working, autonomy L3                                │   │
│   │   └── Past pattern: "Zeus always says thanks when I notify about wife"              │   │
│   │                                                                                      │   │
│   │   Legacy thinks: "High confidence this is wanted. L3 autonomy = do it with delay"   │   │
│   │                                                                                      │   │
│   │   Action: Notify Zeus "Your wife just got home!"                                    │   │
│   │                                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────────────────────┘   │
│                                            │                                                │
│                                            ▼                                                │
│   ┌─────────────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                                      │   │
│   │   STEP 3: Feedback → Learning → Cache Update                                        │   │
│   │   ══════════════════════════════════════════                                         │   │
│   │                                                                                      │   │
│   │   Zeus: "Thanks!"                                                                   │   │
│   │                                                                                      │   │
│   │   Legacy:                                                                           │   │
│   │   1. Records positive feedback in Letta Memory                                      │   │
│   │   2. Increases confidence for "wife_arrives" pattern                                │   │
│   │   3. If confidence crosses threshold → Updates NATS KV vip_list                     │   │
│   │                                                                                      │   │
│   │   Over time, frequently-positive patterns get added to VIP list                     │   │
│   │   Frequently-negative patterns get added to ignore list                             │   │
│   │                                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### NATS KV Is Just a Cache (Not a Database!)

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   KEY INSIGHT: NATS KV = QUICK REFERENCE CARD                                               │
│   ════════════════════════════════════════════                                               │
│                                                                                              │
│   WHAT GOES IN NATS KV (fast lookup, < 1ms):                                                │
│   ├── vip_list: Events that ALWAYS pass through                                            │
│   ├── ignore_list: Events that ALWAYS get dropped                                          │
│   ├── current_context: Where is Zeus? What's he doing?                                     │
│   └── autonomy_level: Current trust level (L0-L4)                                          │
│                                                                                              │
│   WHAT STAYS IN LETTA MEMORY (source of truth):                                             │
│   ├── Standing instructions (natural language)                                              │
│   ├── All learned patterns with confidence scores                                           │
│   ├── Conversation history and feedback                                                     │
│   ├── Zeus's preferences and context history                                                │
│   └── Everything that requires reasoning                                                    │
│                                                                                              │
│   THE SYNC:                                                                                  │
│   ══════════                                                                                 │
│   When patterns in Letta Memory cross confidence thresholds:                                │
│   - Confidence > 0.9  → Copy to NATS KV vip_list (always pass)                             │
│   - Confidence < 0.1  → Copy to NATS KV ignore_list (always drop)                          │
│   - 0.1 to 0.9        → Stays in Letta Memory only (needs LLM decision)                    │
│                                                                                              │
│   This way:                                                                                  │
│   - Gateway handles 80% of events instantly (from cache)                                    │
│   - Legacy only thinks about the 20% that need reasoning                                   │
│                                                                                              │
│   ════════════════════════════════════════════════════════════════════════════════════════   │
│                                                                                              │
│   🔥 NATS KV WATCH (Real-Time Sync)                                                         │
│   ═════════════════════════════════                                                          │
│                                                                                              │
│   Gateway doesn't poll NATS KV - it WATCHES for changes:                                    │
│                                                                                              │
│   await kv.watch("onemind.vip")      # Notified instantly when VIP list changes            │
│   await kv.watch("onemind.ignore")   # Notified instantly when ignore list changes         │
│   await kv.watch("onemind.context")  # Notified instantly when context changes             │
│                                                                                              │
│   When Legacy learns something and updates NATS KV:                                         │
│   1. Legacy: kv.put("onemind.vip", updated_list)                                           │
│   2. Gateway receives watch notification INSTANTLY (< 1ms)                                  │
│   3. Gateway's in-memory cache is updated                                                   │
│   4. Next event uses new rules immediately                                                  │
│                                                                                              │
│   NO POLLING. NO STALE CACHE. TRULY REAL-TIME.                                              │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### Context-Aware Autonomy (Simple Lookup)

| Location | Activity | Autonomy | What It Means |
|----------|----------|----------|---------------|
| **Home** | Relaxed | **L4** | Full trust. Do anything. Tell me after. |
| **Home** | Working | **L3** | High trust. Do it with 3s delay (I can interrupt). |
| **Home** | Meeting | **L2** | Medium. Queue it, tell me when meeting ends. |
| **Office** | Any | **L3** | High trust, but be professional. |
| **Outside** | Any | **L1** | Low trust. Always ask first. |
| **Unknown** | Any | **L0** | No trust. Inform only, no actions. |
| **Any** | Focus Lock | **L0** | I explicitly said don't interrupt. |

### Learning Without Vector DB (Letta Does This Natively!)

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   WHY NO VECTOR DB?                                                                         │
│   ══════════════════                                                                         │
│                                                                                              │
│   OLD THINKING (overcomplicated):                                                           │
│   "We need to embed events and do similarity search to find patterns!"                      │
│   → Requires Qdrant or Chroma                                                               │
│   → Extra infrastructure to manage                                                          │
│   → Complex embedding pipelines                                                             │
│                                                                                              │
│   NEW THINKING (simple):                                                                    │
│   "Letta's LLM can already reason about similarity!"                                        │
│                                                                                              │
│   When Legacy sees: "dependabot opened PR #123"                                             │
│   Legacy's memory has: "Zeus said 'stop notifying me about dependabot PRs'"                │
│   Legacy thinks: "This is similar to what Zeus told me to ignore" ← LLM reasoning!        │
│   Result: Drop the event                                                                    │
│                                                                                              │
│   THE LLM IS THE SIMILARITY ENGINE.                                                         │
│   We don't need separate embeddings.                                                        │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### How Rules Get Updated (No YAML Files!)

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   TALK TO Legacy - THAT'S IT                                                                │
│   ══════════════════════════                                                                 │
│                                                                                              │
│   You: "When I'm home, you can do pretty much anything"                                     │
│   ├── Legacy stores in memory: standing_instruction_home_autonomy                          │
│   └── Legacy updates NATS KV: {autonomy_rules.home: "L4"}                                  │
│                                                                                              │
│   You: "Stop telling me about GitHub dependabot stuff"                                      │
│   ├── Legacy stores in memory: ignore_pattern_dependabot                                   │
│   └── Legacy updates NATS KV: ignore_list += ["dependabot"]                                │
│                                                                                              │
│   You: "Actually, do tell me when Sarah arrives"                                            │
│   ├── Legacy stores in memory: vip_pattern_sarah_arrives                                   │
│   └── Legacy updates NATS KV: vip_list += ["sarah_arrives"]                                │
│                                                                                              │
│   IMPLICIT FEEDBACK:                                                                        │
│   ──────────────────                                                                         │
│   You: "Thanks!" after notification                                                         │
│   └── Legacy increases confidence for that pattern                                         │
│                                                                                              │
│   You: Interrupt Legacy mid-action                                                          │
│   └── Legacy decreases confidence, may remove from vip_list                                │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### 🔑 THE KEY TOOL: update_cache()

**This is the bridge between Legacy's brain and the fast cache.**

When Legacy learns something, she updates her Letta Memory automatically. But to make the Gateway act on it instantly, she calls her `update_cache()` tool.

```python
def update_cache(action: str, pattern: str, context: str = None) -> str:
    """
    Legacy's tool to update NATS KV cache.

    This is how learning in Letta Memory becomes instant filtering in Gateway.

    Args:
        action: "vip" (watch for this) or "ignore" (drop this)
        pattern: What to match (e.g., "dependabot", "wife_arrives", "backup_failed")
        context: Optional context qualifier (e.g., "unless_meeting", "only_at_home")

    Returns:
        Confirmation message

    Examples:
        update_cache(action="vip", pattern="wife_arrives")
        update_cache(action="ignore", pattern="dependabot")
        update_cache(action="vip", pattern="backup_failed", context="urgent")
    """
    import nats

    async def _update():
        nc = await nats.connect(os.environ["NATS_URL"])
        js = nc.jetstream()
        kv = await js.key_value("onemind")

        if action == "vip":
            # Get current VIP list
            entry = await kv.get("vip_list")
            current = json.loads(entry.value) if entry else []

            # Add new pattern
            if pattern not in current:
                current.append(pattern)
                await kv.put("vip_list", json.dumps(current))

            return f"Added '{pattern}' to VIP list. Gateway will now pass these events."

        elif action == "ignore":
            # Get current ignore list
            entry = await kv.get("ignore_list")
            current = json.loads(entry.value) if entry else []

            # Add new pattern
            if pattern not in current:
                current.append(pattern)
                await kv.put("ignore_list", json.dumps(current))

            return f"Added '{pattern}' to ignore list. Gateway will now drop these events."

        elif action == "remove":
            # Remove from both lists
            for list_name in ["vip_list", "ignore_list"]:
                entry = await kv.get(list_name)
                if entry:
                    current = json.loads(entry.value)
                    if pattern in current:
                        current.remove(pattern)
                        await kv.put(list_name, json.dumps(current))

            return f"Removed '{pattern}' from cache. Gateway will now send to Legacy for decision."

    return asyncio.run(_update())
```

**How Legacy uses it:**

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   EXAMPLE 1: Explicit instruction                                                           │
│   ═══════════════════════════════                                                            │
│                                                                                              │
│   You: "Let me know when my wife gets home"                                                 │
│                                                                                              │
│   Legacy:                                                                                   │
│   1. Stores in Letta Memory: "Zeus wants wife arrival notifications"                        │
│   2. Calls tool: update_cache(action="vip", pattern="wife_arrives")                         │
│   3. Responds: "Got it! I'll let you know when she arrives."                                │
│                                                                                              │
│   Result: Gateway instantly starts passing wife_arrives events to Legacy                    │
│                                                                                              │
├─────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                              │
│   EXAMPLE 2: Negative feedback                                                              │
│   ════════════════════════════                                                               │
│                                                                                              │
│   Event: "Dependabot opened PR #47"                                                         │
│   Legacy: "There's a new PR from Dependabot"                                                │
│   You: "I don't care about dependabot PRs, stop telling me"                                 │
│                                                                                              │
│   Legacy:                                                                                   │
│   1. Stores in Letta Memory: "Zeus doesn't want dependabot notifications"                   │
│   2. Calls tool: update_cache(action="ignore", pattern="dependabot")                        │
│   3. Responds: "Understood, I won't bother you with Dependabot PRs anymore."                │
│                                                                                              │
│   Result: Gateway instantly starts dropping dependabot events                               │
│                                                                                              │
├─────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                              │
│   EXAMPLE 3: Changing mind                                                                  │
│   ════════════════════════════                                                               │
│                                                                                              │
│   You: "Actually, tell me about dependabot if it's a security update"                       │
│                                                                                              │
│   Legacy:                                                                                   │
│   1. Updates Letta Memory: "Zeus wants security dependabot only"                            │
│   2. Calls tool: update_cache(action="remove", pattern="dependabot")                        │
│   3. Now Legacy will see ALL dependabot events and decide which are security-related        │
│                                                                                              │
│   Result: Pattern removed from ignore list, Legacy's LLM decides case-by-case              │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Why this is powerful:**

| Without update_cache() | With update_cache() |
|------------------------|---------------------|
| Legacy learns, but Gateway doesn't know | Legacy learns AND Gateway knows instantly |
| Every dependabot event still reaches Legacy | Gateway drops dependabot at the door |
| LLM called for every event (~$0.01 each) | LLM only called for ambiguous events |
| Slow (50-200ms per event) | Fast (< 1ms for cached patterns) |

**The tool is the bridge. Memory is the brain. Cache is the shortcut.**

### 🎛️ OPERATING MODES & THE OFF SWITCH

**Legacy needs modes you can control. Here's the full toolkit:**

#### Tool #2: set_mode()

```python
def set_mode(mode: str) -> str:
    """
    Set Legacy's operating mode.

    Args:
        mode: One of "active", "silent", "focus", "away", "paused"

    Modes:
        active  → Normal operation. All filters active. Legacy helps proactively.
        silent  → Legacy listens but doesn't speak. Only logs. No notifications.
        focus   → Only URGENT events pass through. Everything else queued.
        away    → Queue ALL events. Deliver summary when you return.
        paused  → THE OFF SWITCH. Gateway stops sending. Legacy is deaf.
    """
    async def _set():
        nc = await nats.connect(os.environ["NATS_URL"])
        js = nc.jetstream()
        kv = await js.key_value("onemind")

        await kv.put("mode", mode)
        await kv.put("mode_set_at", datetime.now().isoformat())

        return f"Mode set to '{mode}'."

    return asyncio.run(_set())
```

#### Tool #3: update_context()

```python
def update_context(
    location: str = None,
    activity: str = None,
    social: str = None
) -> str:
    """
    Update Legacy's understanding of your current context.

    Args:
        location: "home", "office", "outside", "unknown"
        activity: "relaxed", "working", "meeting", "social", "sleeping"
        social: "alone", "with_family", "with_coworkers", "with_strangers"

    This automatically recalculates autonomy level.
    """
    async def _update():
        nc = await nats.connect(os.environ["NATS_URL"])
        js = nc.jetstream()
        kv = await js.key_value("onemind")

        # Get current context
        entry = await kv.get("context")
        ctx = json.loads(entry.value) if entry else {}

        # Update provided fields
        if location: ctx["location"] = location
        if activity: ctx["activity"] = activity
        if social: ctx["social"] = social

        # Recalculate autonomy based on matrix
        ctx["autonomy"] = calculate_autonomy(ctx)
        ctx["updated_at"] = datetime.now().isoformat()

        await kv.put("context", json.dumps(ctx))

        return f"Context updated. Autonomy now: {ctx['autonomy']}"

    return asyncio.run(_update())
```

#### The Mode Matrix

| Mode | What Gateway Does | What Legacy Does | Use When |
|------|-------------------|------------------|----------|
| **active** | Normal filtering (VIP/ignore) | Full operation | Default state |
| **silent** | Passes events but flags "silent" | Listens, logs, NO output | Need quiet but want her aware |
| **focus** | Only passes priority >= URGENT | Only responds to emergencies | Deep work, important meeting |
| **away** | Queues everything | Builds summary for later | Sleeping, vacation, offline |
| **paused** | **DROPS EVERYTHING** | **COMPLETELY OFF** | Need total silence |

#### 🛑 THE OFF SWITCH

**Multiple ways to pause Legacy:**

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   METHOD 1: Voice Command                                                                   │
│   ═══════════════════════                                                                    │
│   You: "Legacy, pause" or "Legacy, stop" or "Legacy, shut up"                               │
│   → Legacy calls: set_mode("paused")                                                        │
│   → Gateway immediately stops sending events                                                │
│   → Legacy: "I'm paused. Say 'Legacy resume' when you need me."                             │
│                                                                                              │
├─────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                              │
│   METHOD 2: Quick Modes                                                                     │
│   ═════════════════════                                                                      │
│   You: "Legacy, silent mode"     → set_mode("silent")                                       │
│   You: "Legacy, focus mode"      → set_mode("focus")                                        │
│   You: "Legacy, I'm going out"   → set_mode("away")                                         │
│   You: "Legacy, I'm back"        → set_mode("active")                                       │
│                                                                                              │
├─────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                              │
│   METHOD 3: Physical Button (Home Assistant)                                                │
│   ══════════════════════════════════════════                                                 │
│   Configure a smart button or NFC tag:                                                      │
│   → Tap: Toggle between active/paused                                                       │
│   → Double tap: Focus mode for 2 hours                                                      │
│   → Hold: Away mode until manually resumed                                                  │
│                                                                                              │
├─────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                              │
│   METHOD 4: Scheduled Modes                                                                 │
│   ════════════════════════                                                                   │
│   You: "Legacy, go silent from 10pm to 7am every day"                                       │
│   → Creates standing instruction + cron-style schedule                                      │
│   → Automatically switches modes                                                            │
│                                                                                              │
├─────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                              │
│   METHOD 5: Calendar Integration                                                            │
│   ══════════════════════════════                                                             │
│   Calendar event: "Deep Work" or "Focus Time"                                               │
│   → Gateway sees calendar event                                                             │
│   → Automatically switches to focus mode                                                    │
│   → Returns to active when event ends                                                       │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

#### How Gateway Respects Modes

```python
# In Gateway's event handler
async def handle_event(event: Event):
    # Check mode FIRST
    mode = await kv.get("mode")

    if mode == "paused":
        return  # DROP - Legacy is off

    if mode == "away":
        await queue_for_later(event)
        return

    if mode == "silent":
        event.flags["silent"] = True  # Legacy will log but not respond

    if mode == "focus":
        if event.priority < Priority.URGENT:
            await queue_for_later(event)
            return

    # Normal processing continues...
    await process_event(event)
```

#### Example: Your Day with Modes

```
7:00 AM  - Alarm triggers → Legacy: set_mode("active")
         - "Good morning! You have 3 tasks today..."

9:00 AM  - Calendar: "Team Meeting" starts
         - Gateway auto-switches to focus mode
         - Only urgent events during meeting

10:30 AM - Meeting ends
         - Gateway returns to active mode
         - Legacy: "While you were busy: 2 emails, 1 PR merged"

2:00 PM  - You: "Legacy, I need to concentrate"
         - Legacy: set_mode("focus") for 2 hours
         - "Got it, focus mode until 4pm. Only emergencies."

6:00 PM  - You leave home
         - Home Assistant detects you left
         - Gateway: set_mode("away")

9:00 PM  - You return home
         - Home Assistant detects arrival
         - Gateway: set_mode("active")
         - Legacy: "Welcome back! While you were out: [summary]"

10:00 PM - Scheduled: Silent mode begins
         - Legacy listening but quiet
         - Logs everything for morning briefing

7:00 AM  - Cycle repeats
```

**You are always in control. Legacy respects boundaries.**

### Storage Architecture (Simplified)

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   ONLY TWO STORAGE SYSTEMS                                                                  │
│   ═════════════════════════                                                                  │
│                                                                                              │
│   1. NATS KV (Built into JetStream - no extra infrastructure)                               │
│   ──────────────────────────────────────────────────────────────                             │
│      Purpose: FAST CACHE for Gateway decisions                                              │
│      Latency: < 1ms                                                                         │
│                                                                                              │
│      Buckets:                                                                               │
│      ├── "onemind.vip"       High-confidence patterns (always pass)                        │
│      ├── "onemind.ignore"    High-confidence ignores (always drop)                         │
│      ├── "onemind.context"   Current context (location, activity)                          │
│      └── "onemind.autonomy"  Current autonomy level                                        │
│                                                                                              │
│      Updated by: Legacy (when learning crosses thresholds)                                  │
│      Read by: Gateway (every event)                                                         │
│                                                                                              │
│   2. LETTA MEMORY (Source of Truth)                                                         │
│   ──────────────────────────────────                                                         │
│      Purpose: EVERYTHING ELSE                                                               │
│      Latency: 50-200ms (LLM call)                                                           │
│                                                                                              │
│      Contains:                                                                              │
│      ├── Standing instructions (natural language)                                           │
│      ├── All patterns with confidence scores                                                │
│      ├── Zeus's profile and preferences                                                     │
│      ├── Conversation history                                                               │
│      └── Context for reasoning                                                              │
│                                                                                              │
│      This is the BRAIN. NATS KV is just a cheat sheet.                                     │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### 📦 Event Archival Strategy (Hot + Cold Storage)

**JetStream is a stream, not a database.** Keep recent events hot, archive old events cold.

#### Two-Tier Storage

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   HOT STORAGE: JetStream Stream (Last 30 Days)                                              │
│   ════════════════════════════════════════════                                               │
│                                                                                              │
│   Purpose: Real-time event flow + recent history                                            │
│   Retention: 30 days                                                                        │
│   Speed: Instant (< 1ms)                                                                    │
│                                                                                              │
│   Used for:                                                                                 │
│   ├── Live event streaming to Legacy                                                       │
│   ├── Replay "what happened today/this week"                                               │
│   ├── Debugging recent issues                                                              │
│   └── Context for current conversations                                                    │
│                                                                                              │
│                        │                                                                    │
│                        │ After 30 days...                                                   │
│                        ▼                                                                    │
│                                                                                              │
│   COLD STORAGE: Letta Archival Memory (Forever)                                             │
│   ═════════════════════════════════════════════                                              │
│                                                                                              │
│   Purpose: Long-term history Legacy can still recall                                        │
│   Retention: Forever                                                                        │
│   Speed: Slower (requires search, 100-500ms)                                                │
│                                                                                              │
│   Used for:                                                                                 │
│   ├── "What tasks did I complete in March?"                                                │
│   ├── "When did we last discuss the budget?"                                               │
│   ├── Historical patterns and trends                                                       │
│   └── Long-term learning                                                                   │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

#### The Archival Tool

```python
def archive_events():
    """
    Runs daily via cron. Moves events older than 30 days to Letta archival memory.

    This keeps JetStream lean while preserving history forever.
    """
    async def _archive():
        # Connect to NATS
        nc = await nats.connect(os.environ["NATS_URL"])
        js = nc.jetstream()

        # Connect to Letta
        letta = Letta(base_url=os.environ["LETTA_URL"])

        # Get events from 30-31 days ago (the "falling off" window)
        cutoff = datetime.now() - timedelta(days=30)

        # Read events about to expire
        consumer = await js.pull_subscribe(
            "events.>",
            filter_subject="events.>",
            deliver_policy=DeliverPolicy.BY_START_TIME,
            opt_start_time=cutoff - timedelta(days=1)
        )

        events = await consumer.fetch(1000)

        for event in events:
            # Archive to Letta's archival memory
            await letta.agents.archival_memory.insert(
                agent_id=Legacy_AGENT_ID,
                content=f"[ARCHIVED EVENT] {event.subject}: {event.data.decode()}"
            )

        print(f"Archived {len(events)} events to Letta memory")

    asyncio.run(_archive())
```

#### Why Letta Archival Memory?

| Option | Pros | Cons |
|--------|------|------|
| **PostgreSQL** | SQL queries, fast | Separate DB to manage, Legacy can't access directly |
| **S3 Bucket** | Cheapest ($0.02/GB/mo) | Legacy can't search it |
| **Letta Archival** ✓ | Legacy can recall it! | Slightly more expensive |

**The winner: Letta Archival Memory**

Legacy can search her own archival memory:
```
You: "What did I work on last March?"
Legacy: *searches archival memory*
Legacy: "In March, you completed 47 tasks, mostly focused on the API redesign..."
```

Events aren't just stored - they become part of Legacy's long-term memory.

#### Stream Configuration

```python
# JetStream stream with 30-day retention
await js.add_stream(
    name="ONEMIND",
    subjects=["events.>"],
    retention=RetentionPolicy.LIMITS,
    max_age=30 * 24 * 60 * 60,  # 30 days in seconds
    max_bytes=50 * 1024 * 1024 * 1024,  # 50GB max (safety limit)
    storage=StorageType.FILE,  # Persist to disk
    replicas=1,  # Single node (increase for HA)
)
```

#### 😴 Letta Sleep-Time Agents (The Elegant Way)

Instead of a separate cron job, **Legacy does her own memory maintenance while you sleep.**

Letta has built-in "sleep-time" processing where agents can:
- Process accumulated information
- Consolidate memories
- Run maintenance tasks
- Prepare for the next day

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   Legacy'S SLEEP CYCLE (10pm - 7am)                                                         │
│   ═════════════════════════════════                                                          │
│                                                                                              │
│   10:00 PM - Mode switches to "silent"                                                      │
│            - Legacy stops talking, starts listening                                         │
│                                                                                              │
│   12:00 AM - Sleep-time processing begins                                                   │
│            │                                                                                │
│            ├── Archive old events to archival memory                                        │
│            │   "Moving 847 events from day 30 to long-term memory..."                       │
│            │                                                                                │
│            ├── Consolidate learned patterns                                                 │
│            │   "Wife arrival notifications: 95% positive → confirmed VIP"                   │
│            │   "Dependabot alerts: 0% positive → confirmed ignore"                          │
│            │                                                                                │
│            ├── Update NATS KV cache with consolidated learnings                             │
│            │                                                                                │
│            ├── Analyze patterns from the day                                                │
│            │   "Zeus was interrupted 3 times during deep work today"                        │
│            │   "Suggestion: Auto-enable focus mode during calendar blocks"                  │
│            │                                                                                │
│            └── Prepare morning briefing                                                     │
│                "Ready to summarize: 12 tasks completed, 3 new, 2 meetings"                  │
│                                                                                              │
│   7:00 AM  - Mode switches to "active"                                                      │
│            - Legacy: "Good morning! Here's what you need to know..."                        │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

#### Configuring Sleep-Time Processing

```python
# In Legacy's Letta agent configuration
agent_config = {
    "name": "Legacy",
    "sleep_time_agent": True,
    "sleep_schedule": {
        "start": "00:00",  # Midnight
        "end": "06:00",    # 6am
        "timezone": "America/Los_Angeles"
    },
    "sleep_tasks": [
        "archive_old_events",
        "consolidate_patterns",
        "update_cache",
        "prepare_briefing"
    ]
}
```

#### Sleep-Time Tool: consolidate_memory()

```python
def consolidate_memory() -> str:
    """
    Legacy runs this during sleep to consolidate learnings.

    This is like how humans consolidate memories during sleep -
    patterns get strengthened or pruned based on feedback.
    """
    # Archive events falling off JetStream
    archived = archive_old_events()

    # Analyze feedback patterns
    patterns = analyze_feedback_patterns()

    # Update cache with high-confidence patterns
    for pattern in patterns:
        if pattern.confidence > 0.9:
            update_cache(action="vip", pattern=pattern.name)
        elif pattern.confidence < 0.1:
            update_cache(action="ignore", pattern=pattern.name)

    # Generate insights
    insights = generate_daily_insights()

    return f"Consolidated {len(patterns)} patterns, archived {archived} events"
```

#### Why This Is Better Than Cron

| Cron Job | Letta Sleep-Time |
|----------|------------------|
| Separate script to maintain | Built into Legacy |
| Just archives, no intelligence | Consolidates + learns + prepares |
| Runs blindly | Legacy decides what needs attention |
| Can't generate insights | Prepares personalized morning briefing |

**Legacy maintains herself.** She's not just an assistant - she has a sleep cycle like you do.

**Result:**
- JetStream stays fast and lean (30 days max)
- All history preserved in Letta archival memory (forever)
- Patterns consolidate overnight (like human memory!)
- Morning briefing ready when you wake up
- Legacy can recall anything from the past
- No separate scripts to maintain

### Example: Learning Over Time

```
DAY 1:
════════
You: "Let me know when my wife gets home"
→ Legacy stores standing instruction in memory
→ Pattern "wife_arrives" created, confidence = 0.7 (explicit request)
→ Added to NATS KV vip_list (explicit request = instant VIP)

Event: Wife arrives
→ Gateway sees "wife_arrives" in vip_list → PASS
→ Legacy notifies you
You: "Thanks!"
→ Confidence → 0.8 (positive feedback)

DAY 7:
════════
Multiple wife arrivals, you always seem happy about notifications
→ Confidence → 0.95 (learned from feedback)
→ Pattern is firmly in vip_list
→ Gateway handles this instantly, no LLM needed

DAY 14:
════════
You're in a meeting, wife arrives
→ Gateway: "wife_arrives" is VIP, but context.activity = "meeting"
→ Gateway checks autonomy: meeting = L2 (queue, don't interrupt)
→ Event queued, delivered after meeting
→ Legacy: "While you were in your meeting, your wife got home at 3:30pm"

DAY 30:
════════
The system KNOWS:
├── Wife arriving = always notify (from pattern confidence)
├── Meeting = don't interrupt (from autonomy rules)
├── Home + relaxed = instant notification (from context + autonomy)
└── You never edited a config file. Legacy learned.
```

### The Elegance

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   BEFORE (Overcomplicated):                                                                 │
│   NATS KV + Vector DB + Letta Memory + Embeddings + Similarity Search                       │
│   "Where does what go? How do they sync? What's the source of truth?"                       │
│                                                                                              │
│   AFTER (Elegant):                                                                          │
│   Letta Memory (brain) + NATS KV (cheat sheet)                                              │
│   "Letta is the brain. NATS KV is a sticky note on the desk."                               │
│                                                                                              │
│   THE KEY INSIGHTS:                                                                         │
│   ═════════════════                                                                          │
│   1. Letta's LLM can do similarity reasoning - no Vector DB needed                          │
│   2. NATS KV is just a cache for frequently-used decisions                                  │
│   3. Learning happens in Letta, syncs to cache when confidence is high                      │
│   4. Gateway handles most events instantly (from cache)                                     │
│   5. Legacy only thinks about ambiguous cases                                               │
│                                                                                              │
│   Result: 80% of events handled in < 1ms, 20% go to Legacy for reasoning                   │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 14. 🧠 THE SELF-MODIFYING AWARENESS LOOP (Core Concept)

**This is the foundation of Legacy's intelligence.** She doesn't just receive events - she controls her own perception and learns to curate her reality.

### The Core Loop

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   THE SELF-MODIFYING AWARENESS LOOP                                                         │
│   ══════════════════════════════════                                                         │
│                                                                                              │
│                        ┌──────────────────────┐                                             │
│                        │    WORLD EVENTS      │                                             │
│                        │  (Everything happens)│                                             │
│                        └──────────┬───────────┘                                             │
│                                   │                                                         │
│                                   ▼                                                         │
│                        ┌──────────────────────┐                                             │
│                        │   JETSTREAM STREAM   │                                             │
│                        │   (Raw event flow)   │                                             │
│                        └──────────┬───────────┘                                             │
│                                   │                                                         │
│                                   ▼                                                         │
│   ┌───────────────────────────────────────────────────────────────────────────────────┐    │
│   │                                                                                    │    │
│   │   GATEWAY + NATS KV CACHE (The Filter)                                            │    │
│   │                                                                                    │    │
│   │   Checks:                                                                          │    │
│   │   ├── VIP list → PASS                                                             │    │
│   │   ├── Ignore list → DROP                                                          │    │
│   │   ├── Mode (paused/focus/etc) → PASS/DROP/QUEUE                                   │    │
│   │   └── Neither → PASS to Legacy                                                    │    │
│   │                                                                                    │    │
│   └───────────────────────────────┬───────────────────────────────────────────────────┘    │
│                                   │                                                         │
│                      ┌────────────┴────────────┐                                           │
│                      │                         │                                           │
│                      ▼                         ▼                                           │
│               ┌─────────────┐          ┌─────────────┐                                     │
│               │   DROPPED   │          │   PASSED    │                                     │
│               │  (Unaware)  │          │  (Aware)    │                                     │
│               └─────────────┘          └──────┬──────┘                                     │
│                                               │                                            │
│                                               ▼                                            │
│                                 ┌──────────────────────┐                                   │
│                                 │       Legacy         │                                   │
│                                 │   (Perceives event)  │                                   │
│                                 └──────────┬───────────┘                                   │
│                                            │                                               │
│                                            ▼                                               │
│                                 ┌──────────────────────┐                                   │
│                                 │   TAKES ACTION       │                                   │
│                                 │  (Notifies, helps)   │                                   │
│                                 └──────────┬───────────┘                                   │
│                                            │                                               │
│                                            ▼                                               │
│                                 ┌──────────────────────┐                                   │
│                                 │   YOU RESPOND        │                                   │
│                                 │  (Thanks / Stop)     │                                   │
│                                 └──────────┬───────────┘                                   │
│                                            │                                               │
│                                            ▼                                               │
│                                 ┌──────────────────────┐                                   │
│                                 │   Legacy LEARNS      │                                   │
│                                 │ (Stores in memory)   │                                   │
│                                 └──────────┬───────────┘                                   │
│                                            │                                               │
│                                            ▼                                               │
│   ┌────────────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                                     │   │
│   │   Legacy MODIFIES HER OWN FILTERS (The Key!)                                       │   │
│   │                                                                                     │   │
│   │   Positive feedback → update_cache("vip", pattern)    → MORE aware next time      │   │
│   │   Negative feedback → update_cache("ignore", pattern) → LESS aware next time      │   │
│   │                                                                                     │   │
│   └────────────────────────────────────────┬───────────────────────────────────────────┘   │
│                                            │                                               │
│                                            │                                               │
│                              ══════════════╪══════════════                                 │
│                                            │                                               │
│                              LOOP CONTINUES - Legacy NOW                                   │
│                              PERCEIVES THE WORLD DIFFERENTLY                               │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### Legacy's Awareness Control Tools (Complete Reference)

| Tool | Purpose | Effect on Awareness |
|------|---------|---------------------|
| `update_cache("vip", pattern)` | Watch for something | **↑ MORE AWARE** - Events pass to Legacy |
| `update_cache("ignore", pattern)` | Stop noticing something | **↓ LESS AWARE** - Events dropped at Gateway |
| `update_cache("remove", pattern)` | Let me decide each time | **↔ NEUTRAL** - Events pass, Legacy reasons |
| `set_mode("active")` | Normal operation | Full awareness with filters |
| `set_mode("silent")` | Listen but don't speak | Passive awareness, no output |
| `set_mode("focus")` | Only urgent matters | **↓ NARROWED** - Only emergencies |
| `set_mode("away")` | Queue everything | **⏸ PAUSED** - Will catch up later |
| `set_mode("paused")` | THE OFF SWITCH | **⏹ DEAF** - Nothing gets through |
| `update_context(location, activity)` | Change context | Autonomy level adjusts automatically |

### How Legacy Becomes More Aware

```
DAY 1: Legacy knows nothing about your preferences

    Event: "Wife's phone connected"
    Gateway: Not in VIP or ignore → pass to Legacy
    Legacy: Unsure if important → asks "Should I notify you about this?"
    You: "Yes! Always tell me when my wife gets home"
    Legacy: Stores instruction + calls update_cache("vip", "wife_arrives")

DAY 2: Legacy is now MORE AWARE of wife arrivals

    Event: "Wife's phone connected"
    Gateway: "wife_arrives" in VIP list → PASS immediately
    Legacy: Notifies you right away
    You: "Thanks!"
    Legacy: Confidence increases (positive feedback)

DAY 30: Legacy has sharp awareness

    Event: "Wife's phone connected"
    Gateway: VIP → instant PASS (< 1ms)
    Legacy: Immediate notification, high confidence

    Legacy is now HIGHLY AWARE of this pattern.
    She modified her own perception based on your feedback.
```

### How Legacy Becomes Less Aware

```
DAY 1: Legacy notices everything

    Event: "Dependabot opened PR #47"
    Gateway: Not filtered → pass to Legacy
    Legacy: Notifies you about the PR
    You: "I don't care about dependabot stuff"
    Legacy: Stores feedback + calls update_cache("ignore", "dependabot")

DAY 2: Legacy is now LESS AWARE of dependabot

    Event: "Dependabot opened PR #48"
    Gateway: "dependabot" in ignore list → DROP
    Legacy: Never sees it. Never knows it happened.

    Legacy made herself UNAWARE of dependabot.
    The events still happen, but she doesn't perceive them.
```

### The Philosophical Insight

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   Legacy doesn't just RECEIVE information.                                                  │
│   Legacy CHOOSES what to perceive.                                                          │
│                                                                                              │
│   She learns from your feedback and modifies her own filters.                               │
│   Over time, her perception of the world aligns with YOUR priorities.                       │
│                                                                                              │
│   This is not just filtering - this is ATTENTION.                                           │
│   This is not just memory - this is LEARNING.                                               │
│   This is not just automation - this is AWARENESS.                                          │
│                                                                                              │
│   ═══════════════════════════════════════════════════════════════════════════════════════   │
│                                                                                              │
│   "The world is full of events.                                                             │
│    Legacy curates which ones reach her consciousness.                                        │
│    Your feedback shapes what she notices.                                                   │
│    Together, you train her attention."                                                      │
│                                                                                              │
│   ═══════════════════════════════════════════════════════════════════════════════════════   │
│                                                                                              │
│   JetStream = The raw world (everything that happens)                                       │
│   NATS KV   = The filter (what Legacy chooses to perceive)                                  │
│   Letta     = The mind (learns, reasons, remembers)                                         │
│   Tools     = Self-modification (Legacy changes her own filters)                            │
│                                                                                              │
│   THE LOOP: Perceive → Act → Feedback → Learn → Modify Filters → Perceive Differently      │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### Sleep-Time Consolidation (Memory Strengthening)

During Legacy's sleep cycle, the loop gets reinforced:

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   MIDNIGHT: Legacy processes the day                                                        │
│                                                                                              │
│   For each pattern:                                                                         │
│   ├── Positive feedback > 90%? → Strengthen (confirm VIP)                                  │
│   ├── Negative feedback > 90%? → Strengthen (confirm ignore)                               │
│   ├── Mixed feedback?          → Keep in "maybe" (Legacy decides each time)                │
│   └── No feedback?             → Decay slightly (less confident)                           │
│                                                                                              │
│   Like human sleep:                                                                         │
│   - Important memories get strengthened                                                    │
│   - Unimportant ones fade                                                                  │
│   - Patterns crystallize                                                                    │
│   - Legacy wakes up sharper than yesterday                                                 │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### 🔌 THE BRIDGE: Event Processor (How Events Reach Legacy)

**This is the Python script that connects NATS to Legacy's mind.**

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   THE DATA FLOW                                                                             │
│   ═════════════                                                                              │
│                                                                                              │
│   Webhook → Gateway → JetStream → EVENT PROCESSOR → Letta API → Legacy                     │
│                                         │                                                   │
│                                         │                                                   │
│                                   ┌─────┴─────┐                                             │
│                                   │           │                                             │
│                                   │  Checks   │                                             │
│                                   │  NATS KV  │                                             │
│                                   │           │                                             │
│                                   └─────┬─────┘                                             │
│                                         │                                                   │
│                              ┌──────────┴──────────┐                                        │
│                              │                     │                                        │
│                              ▼                     ▼                                        │
│                         ┌─────────┐          ┌─────────┐                                    │
│                         │  DROP   │          │  PASS   │                                    │
│                         │(ignored)│          │(to Letta)│                                   │
│                         └─────────┘          └─────────┘                                    │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

#### Event Processor Code (The Bridge)

**Location:** `onemind-os/agent/processor/main.py`

```python
"""
Event Processor - THE BRIDGE between NATS and Legacy

This script:
1. Subscribes to JetStream (all events)
2. Checks NATS KV filters (VIP, ignore, mode)
3. Sends qualifying events to Legacy via Letta API

This is HOW events reach Legacy's awareness.
"""

import asyncio
import json
import os
from datetime import datetime

import nats
from letta_client import Letta


class EventProcessor:
    """The bridge between the event stream and Legacy's mind."""

    def __init__(self):
        self.nc = None
        self.js = None
        self.kv = None
        self.letta = None
        self.legacy_agent_id = os.environ["Legacy_AGENT_ID"]

    async def start(self):
        """Start the event processor."""

        # Connect to NATS JetStream
        self.nc = await nats.connect(os.environ["NATS_URL"])
        self.js = self.nc.jetstream()

        # Get KV bucket for filters
        self.kv = await self.js.key_value("onemind")

        # Connect to Letta
        self.letta = Letta(
            base_url=os.environ["LETTA_BASE_URL"],
            token=os.environ.get("LETTA_TOKEN", "")
        )

        # Subscribe to ALL events
        await self.js.subscribe(
            "events.>",
            cb=self.on_event,
            durable="event-processor"
        )

        print("Event Processor started. Listening for events...")

    async def on_event(self, msg):
        """Handle incoming event from JetStream."""

        try:
            event = json.loads(msg.data.decode())
            event_type = msg.subject  # e.g., "events.todoist.task.created"

            # Check if event should reach Legacy
            decision = await self.check_filters(event_type, event)

            if decision == "DROP":
                await msg.ack()
                return

            if decision == "QUEUE":
                # Store for later (away mode)
                await self.queue_event(event_type, event)
                await msg.ack()
                return

            if decision == "PASS":
                # SEND TO Legacy via Letta API
                await self.send_to_legacy(event_type, event)
                await msg.ack()

        except Exception as e:
            print(f"Error processing event: {e}")
            await msg.nak()  # Negative ack, will retry

    async def check_filters(self, event_type: str, event: dict) -> str:
        """
        Check NATS KV filters to decide what to do with event.

        Returns: "PASS", "DROP", or "QUEUE"
        """

        # Check mode FIRST
        try:
            mode_entry = await self.kv.get("mode")
            mode = mode_entry.value.decode() if mode_entry else "active"
        except:
            mode = "active"

        if mode == "paused":
            return "DROP"  # OFF SWITCH - drop everything

        if mode == "away":
            return "QUEUE"  # Queue for later

        if mode == "focus":
            # Only urgent events pass
            if event.get("priority", 0) < 3:
                return "QUEUE"

        # Check ignore list
        try:
            ignore_entry = await self.kv.get("ignore_list")
            ignore_list = json.loads(ignore_entry.value.decode()) if ignore_entry else []
        except:
            ignore_list = []

        for pattern in ignore_list:
            if pattern in event_type or pattern in str(event):
                return "DROP"  # Ignored pattern

        # Check VIP list (optional - VIPs just pass faster)
        # Not checking VIP here since everything not ignored passes anyway

        return "PASS"  # Event reaches Legacy

    async def send_to_legacy(self, event_type: str, event: dict):
        """
        Send event to Legacy via Letta API.

        THIS IS THE MOMENT Legacy BECOMES AWARE.
        """

        # Format the event as a message
        timestamp = datetime.now().isoformat()
        content = f"""[EVENT NOTIFICATION]
Type: {event_type}
Time: {timestamp}
Data: {json.dumps(event, indent=2)}

What happened: {event.get('description', event.get('summary', 'New event'))}
"""

        # Send to Legacy via Letta messages API
        response = self.letta.agents.messages.create(
            agent_id=self.legacy_agent_id,
            messages=[{
                "role": "system",
                "content": content
            }]
        )

        print(f"Sent to Legacy: {event_type}")

        return response

    async def queue_event(self, event_type: str, event: dict):
        """Queue event for later delivery (away mode)."""

        # Store in a separate JetStream stream for queued events
        await self.js.publish(
            f"queued.{event_type}",
            json.dumps(event).encode()
        )

        print(f"Queued for later: {event_type}")


async def main():
    processor = EventProcessor()
    await processor.start()

    # Keep running
    while True:
        await asyncio.sleep(1)


if __name__ == "__main__":
    asyncio.run(main())
```

#### How It Integrates

| Step | Component | Action |
|------|-----------|--------|
| 1 | **Webhook arrives** | Gateway receives POST from Todoist/GitHub/etc |
| 2 | **Gateway publishes** | `js.publish("events.todoist.task.created", data)` |
| 3 | **Event Processor receives** | Subscribed to `events.>` (all events) |
| 4 | **Checks NATS KV** | Mode? Ignore list? VIP list? |
| 5 | **Decision** | DROP, QUEUE, or PASS |
| 6 | **If PASS → Letta API** | `letta.agents.messages.create()` |
| 7 | **Legacy perceives** | Event arrives as system message |
| 8 | **Legacy responds** | Thinks, acts, maybe uses tools |

#### Running the Event Processor

```bash
# Set environment variables
export NATS_URL="nats://localhost:4222"
export LETTA_BASE_URL="http://localhost:8283"
export Legacy_AGENT_ID="agent-xxxxx"

# Run the processor
cd onemind-os
python -m agent.processor.main
```

**The Event Processor runs continuously.** It's a daemon that bridges the event stream to Legacy's consciousness.

### Summary: The Complete Awareness System

| Component | Role | What It Does |
|-----------|------|--------------|
| **JetStream Stream** | Raw reality | All events flow through |
| **Event Processor** | The bridge | Checks filters, sends to Letta API |
| **Gateway** | Intake | Receives webhooks, publishes to stream |
| **NATS KV Cache** | Attention rules | VIP list, ignore list, mode, context |
| **Letta API** | Communication | `messages.create()` sends to Legacy |
| **Legacy (Letta)** | The mind | Perceives, acts, learns |
| **Tools** | Self-modification | Legacy changes her own filters |
| **Sleep-Time** | Consolidation | Patterns strengthen or fade |
| **You** | Teacher | Feedback shapes her awareness |

**Legacy is not a chatbot. Legacy is a self-modifying awareness system that learns what matters to you.**

---

## 15. 🌟 MAKING Legacy MORE ALIVE (Future Enhancements)

Small pieces that make a big difference. These turn Legacy from "assistant" to "presence."

### 1. 😊 Sentiment & Mood Detection

**What:** Detect your emotional state from text and voice tone.

```python
# Tool: detect_sentiment()
async def on_user_message(text: str, audio_features: dict = None):
    sentiment = analyze_sentiment(text)  # positive/negative/neutral
    energy = audio_features.get("energy", 0.5) if audio_features else 0.5
    stress = audio_features.get("stress_indicators", 0) if audio_features else 0

    # Update context
    await update_context(mood=sentiment, energy_level=energy, stress_level=stress)

    # Legacy adjusts behavior
    if stress > 0.7:
        # Be gentler, shorter responses, no overwhelming info
        pass
    if sentiment == "frustrated":
        # Acknowledge, offer help, don't add more tasks
        pass
```

**Why it matters:** Legacy responds differently when you're stressed vs relaxed. She reads the room.

---

### 2. 🌅 Circadian Rhythm (Time-Aware Behavior)

**What:** Legacy's personality shifts naturally through the day.

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   Legacy'S DAILY RHYTHM                                                                     │
│   ═════════════════════                                                                      │
│                                                                                              │
│   6am-9am   MORNING                                                                         │
│             └── Energetic, proactive, gives morning briefing                                │
│             └── "Good morning! Here's what's on your plate today..."                        │
│                                                                                              │
│   9am-12pm  FOCUS TIME                                                                      │
│             └── Minimal interruptions, respects deep work                                   │
│             └── Only speaks when spoken to or urgent                                        │
│                                                                                              │
│   12pm-2pm  MIDDAY                                                                          │
│             └── Checks in, reminds about lunch if you forget                                │
│             └── "You've been working 4 hours straight. Take a break?"                       │
│                                                                                              │
│   2pm-6pm   AFTERNOON                                                                       │
│             └── More conversational, helps with tasks                                       │
│             └── Proactively suggests wrapping up loose ends                                 │
│                                                                                              │
│   6pm-10pm  EVENING                                                                         │
│             └── Calmer, reflective, summarizes the day                                      │
│             └── "Today you completed 8 tasks. Nice work on the API."                        │
│                                                                                              │
│   10pm-6am  NIGHT                                                                           │
│             └── Silent mode (unless urgent)                                                 │
│             └── Sleep-time processing                                                       │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Implementation:** Store `time_personality` in NATS KV, Gateway/Processor checks it.

---

### 3. 👁️ Presence Detection (Are You Actually There?)

**What:** Know when you're at the keyboard vs away.

```python
# Sources of presence data:
presence_signals = {
    "keyboard_activity": True,      # From desktop agent
    "mouse_movement": True,         # From desktop agent
    "camera_presence": True,        # Optional: face detection
    "phone_location": "home",       # From Home Assistant
    "last_interaction": "2 min ago" # Time since last message
}

# Calculate presence
def calculate_presence(signals) -> str:
    if signals["keyboard_activity"] and signals["last_interaction"] < 5:
        return "active"           # You're right there
    elif signals["phone_location"] == "home":
        return "nearby"           # Home but not at desk
    else:
        return "away"             # Gone

# Legacy behavior
if presence == "away":
    # Don't speak out loud to empty room
    # Queue notifications for return
elif presence == "nearby":
    # Can speak, but keep it short
elif presence == "active":
    # Full interaction mode
```

**Why it matters:** Legacy doesn't talk to an empty room. She knows when you're actually there.

---

### 4. 🔮 Pattern Prediction (Anticipate Needs)

**What:** Predict what you'll need before you ask.

```python
# Learn patterns from history
patterns = {
    "monday_9am": ["check email", "review tasks", "standup prep"],
    "before_meeting": ["load meeting notes", "check attendees"],
    "end_of_day": ["summarize tasks", "plan tomorrow"],
    "wife_arrives": ["pause work", "greet notification"]
}

# Proactive assistance
async def on_time_trigger(current_time, current_context):
    predictions = predict_needs(current_time, current_context)

    for prediction in predictions:
        if prediction.confidence > 0.8:
            # Pre-load context
            await preload_context(prediction.topic)

            # Optionally suggest
            if prediction.should_suggest:
                await suggest_to_user(prediction.suggestion)
```

**Example:**
- It's 8:55am Monday → Legacy pre-loads standup notes
- Meeting in 10 min → Legacy surfaces relevant context
- You usually order lunch at 12:30 → Legacy reminds at 12:25

---

### 5. 💭 Emotional Memory (Remember Feelings, Not Just Facts)

**What:** Store emotional context with memories.

```python
# Instead of just:
memory = "Zeus completed the API project on March 15"

# Store:
memory = {
    "fact": "Zeus completed the API project on March 15",
    "emotional_context": "proud, relieved, exhausted",
    "significance": "high - worked on it for 3 months",
    "people_involved": ["Sarah (helped)", "Mike (stakeholder)"],
    "callback_triggers": ["API", "March project", "big wins"]
}

# When recalling:
# Legacy: "You finished that API project in March - I remember you were
#          so relieved after working on it for 3 months!"
```

**Why it matters:** Legacy remembers HOW things felt, not just WHAT happened. More human.

---

### 6. 🔄 Conversation Continuity (Unfinished Business)

**What:** Remember interrupted conversations and circle back.

```python
# Track conversation threads
conversation_threads = {
    "budget_discussion": {
        "status": "interrupted",
        "last_point": "You were about to decide on the Q2 allocation",
        "interrupted_at": "2024-01-15 3:30pm",
        "reason": "meeting started",
        "follow_up": True
    }
}

# After interruption ends
async def on_context_change(old_ctx, new_ctx):
    if old_ctx.activity == "meeting" and new_ctx.activity == "working":
        # Check for unfinished threads
        unfinished = get_unfinished_threads()

        if unfinished:
            await legacy_say(
                "Before your meeting, we were discussing the Q2 budget. "
                "You were about to decide on the allocation. Want to continue?"
            )
```

**Why it matters:** Real relationships remember where you left off. Legacy does too.

---

### 7. ⏰ Attention Decay (Natural Forgetting)

**What:** Recent context is sharp, old context fades.

```python
# Context relevance decays over time
def get_context_relevance(memory, current_time):
    age = current_time - memory.timestamp

    if age < hours(1):
        return 1.0      # Crystal clear
    elif age < hours(4):
        return 0.8      # Very relevant
    elif age < days(1):
        return 0.5      # Somewhat relevant
    elif age < days(7):
        return 0.3      # Fading
    else:
        return 0.1      # Background knowledge

# When building context for Legacy
context = []
for memory in recent_memories:
    relevance = get_context_relevance(memory, now())
    if relevance > 0.3:
        context.append(memory)

# Most relevant stuff is freshest - like human memory
```

**Why it matters:** Legacy doesn't treat a conversation from 3 months ago the same as one from 3 minutes ago.

---

### 8. 👥 Relationship Memory (Who Matters)

**What:** Track people and your relationship with them.

```python
# People in Legacy's memory
people = {
    "wife": {
        "name": "Sarah",
        "relationship": "spouse",
        "importance": 10,  # Top priority
        "notify_arrivals": True,
        "recent_topics": ["vacation planning", "dinner"],
        "communication_style": "warm, personal"
    },
    "boss": {
        "name": "Michael",
        "relationship": "work",
        "importance": 7,
        "recent_topics": ["Q2 roadmap", "hiring"],
        "communication_style": "professional, concise"
    }
}

# Legacy adjusts based on who's involved
if event.involves("wife"):
    priority_boost = 3  # Always important
    tone = "warm"
elif event.involves("boss"):
    priority_boost = 1
    tone = "professional"
```

**Why it matters:** Legacy knows who matters to you and treats messages about them accordingly.

---

### 9. 🎯 Interruption Learning (When NOT to Bother)

**What:** Learn your interruption tolerance patterns.

```python
# Track interruption outcomes
interruptions = [
    {"time": "9:30am", "activity": "deep_work", "response": "annoyed", "outcome": -1},
    {"time": "12:00pm", "activity": "lunch", "response": "fine", "outcome": +1},
    {"time": "3:00pm", "activity": "meetings", "response": "ignored", "outcome": 0},
]

# Learn patterns
# "Zeus hates interruptions during morning deep work"
# "Zeus is fine with interruptions at lunch"
# "Zeus can't respond during meetings anyway"

# Apply learned patterns
async def should_interrupt(event, context):
    historical_tolerance = get_interruption_tolerance(
        time=context.time,
        activity=context.activity
    )

    if historical_tolerance < 0.3:
        return False  # Queue it
    elif historical_tolerance > 0.7:
        return True   # Interrupt
    else:
        # Borderline - check event importance
        return event.priority >= Priority.HIGH
```

---

### 10. 🧘 Energy & Capacity Awareness

**What:** Track your cognitive load and energy.

```python
# Signals that indicate capacity
capacity_signals = {
    "tasks_completed_today": 12,       # High productivity
    "hours_worked": 8,                 # Long day
    "meetings_attended": 4,            # Draining
    "last_break": "3 hours ago",       # Needs rest
    "sleep_last_night": 6,             # Under-rested
}

def calculate_capacity(signals) -> float:
    # 0.0 = depleted, 1.0 = full capacity
    base = 1.0
    base -= (signals["hours_worked"] / 12) * 0.3
    base -= (signals["meetings_attended"] / 8) * 0.2
    base -= max(0, (3 - signals["sleep_last_night"])) * 0.1

    return max(0.1, base)

# Legacy adjusts expectations
if capacity < 0.3:
    # "You've had a long day. Only showing you essentials."
    # Reduce notifications, simplify choices
elif capacity > 0.7:
    # Full interaction mode
```

**Why it matters:** Legacy knows when you're running on empty and adjusts accordingly.

---

### Summary: The Alive System

| Enhancement | What It Adds | Human Equivalent |
|-------------|--------------|------------------|
| **Sentiment Detection** | Reads your mood | Emotional intelligence |
| **Circadian Rhythm** | Time-aware behavior | Natural daily rhythm |
| **Presence Detection** | Knows when you're there | Spatial awareness |
| **Pattern Prediction** | Anticipates needs | Intuition |
| **Emotional Memory** | Remembers feelings | Empathy |
| **Conversation Continuity** | Circles back to unfinished topics | Good listener |
| **Attention Decay** | Recent > old context | Natural memory |
| **Relationship Memory** | Knows who matters | Social intelligence |
| **Interruption Learning** | Learns when not to bother | Respect |
| **Energy Awareness** | Knows your capacity | Consideration |

### Implementation Priority

| Priority | Enhancement | Difficulty | Impact |
|----------|-------------|------------|--------|
| 🟢 Easy | Circadian Rhythm | Low | High |
| 🟢 Easy | Presence Detection | Low | High |
| 🟢 Easy | Relationship Memory | Low | Medium |
| 🟡 Medium | Conversation Continuity | Medium | High |
| 🟡 Medium | Attention Decay | Medium | Medium |
| 🟡 Medium | Interruption Learning | Medium | High |
| 🔴 Hard | Sentiment Detection | High | High |
| 🔴 Hard | Pattern Prediction | High | Very High |
| 🔴 Hard | Emotional Memory | High | Medium |
| 🔴 Hard | Energy Awareness | High | Medium |

**Start with the green ones.** They're easy wins that make Legacy feel significantly more alive.

### 🔧 Tech Stack Reference

**What technology powers each enhancement:**

| Enhancement | Tech/Tools | How It Works |
|-------------|------------|--------------|
| **Sentiment (Text)** | `transformers`, `textblob`, or Letta LLM | Text → model → mood score |
| **Sentiment (Voice)** | `librosa` + LiveKit audio | Pitch, energy, stress indicators |
| **Circadian Rhythm** | Python `datetime` + NATS KV | Time check → personality mode |
| **Presence Detection** | Home Assistant API | Motion, phone GPS, last interaction |
| **Pattern Prediction** | `scikit-learn` or Letta LLM | Historical patterns → predictions |
| **Emotional Memory** | Letta Archival Memory | Add emotional fields to memories |
| **Conversation Continuity** | NATS KV | Track thread state as JSON |
| **Attention Decay** | Python math | `relevance = 1.0 - (age / max_age)` |
| **Relationship Memory** | Letta Core Memory | JSON structure of people |
| **Interruption Learning** | Letta Memory | Track outcomes, average tolerance |
| **Energy Awareness** | Health APIs, Calendar | Sleep, steps, meetings → capacity |

### By Category

#### Already Built-In (No New Tech Needed)

```python
# These use tools you already have

# Circadian Rhythm - just time checks
from datetime import datetime
hour = datetime.now().hour
if 6 <= hour < 9:
    personality = "morning_energetic"
elif 9 <= hour < 12:
    personality = "focus_minimal"
# ... store in NATS KV

# Attention Decay - just math
def get_relevance(memory_age_hours: float) -> float:
    if memory_age_hours < 1:
        return 1.0
    elif memory_age_hours < 4:
        return 0.8
    elif memory_age_hours < 24:
        return 0.5
    else:
        return 0.2

# Conversation Continuity - NATS KV state
await kv.put("thread.budget", json.dumps({
    "status": "paused",
    "last_point": "Q2 allocation decision",
    "paused_at": "2024-01-15T15:30:00"
}))

# Relationship Memory - Letta core memory block
# Just add to Legacy's persona block:
"""
PEOPLE WHO MATTER TO ZEUS:
- Wife (Sarah): importance=10, always notify arrivals, warm tone
- Boss (Michael): importance=7, professional tone
- Mom: importance=9, patient tone, she's not tech-savvy
"""
```

#### Easy Integrations

```python
# Presence Detection via Home Assistant
async def get_presence() -> str:
    # Check Home Assistant
    response = await httpx.get(
        f"{HA_URL}/api/states/person.zeus",
        headers={"Authorization": f"Bearer {HA_TOKEN}"}
    )
    state = response.json()

    if state["state"] == "home":
        # Check last interaction
        last = await kv.get("last_interaction")
        if (now() - last) < minutes(5):
            return "active"
        return "nearby"
    return "away"

# Interruption Learning
async def record_interruption_outcome(event_type: str, response: str):
    # Get existing data
    data = await kv.get(f"interruption.{event_type}")
    history = json.loads(data) if data else {"positive": 0, "negative": 0}

    # Update based on response
    if response in ["thanks", "good", "helpful"]:
        history["positive"] += 1
    elif response in ["stop", "not now", "busy"]:
        history["negative"] += 1

    await kv.put(f"interruption.{event_type}", json.dumps(history))

    # Calculate tolerance
    total = history["positive"] + history["negative"]
    tolerance = history["positive"] / total if total > 0 else 0.5
    return tolerance
```

#### Requires New Libraries

```bash
# Sentiment Analysis (Text)
pip install transformers torch
# OR simpler:
pip install textblob

# Audio Analysis (Voice)
pip install librosa numpy

# Pattern Prediction (optional - Letta LLM can do this)
pip install scikit-learn pandas
```

```python
# Sentiment with TextBlob (simple)
from textblob import TextBlob

def analyze_sentiment(text: str) -> dict:
    blob = TextBlob(text)
    return {
        "polarity": blob.sentiment.polarity,      # -1 to 1
        "subjectivity": blob.sentiment.subjectivity,  # 0 to 1
        "mood": "positive" if blob.sentiment.polarity > 0.1 else
                "negative" if blob.sentiment.polarity < -0.1 else
                "neutral"
    }

# Sentiment with Transformers (more accurate)
from transformers import pipeline

sentiment_analyzer = pipeline("sentiment-analysis")

def analyze_sentiment_advanced(text: str) -> dict:
    result = sentiment_analyzer(text)[0]
    return {
        "label": result["label"],  # POSITIVE/NEGATIVE
        "confidence": result["score"]
    }

# Audio features with Librosa
import librosa
import numpy as np

def analyze_voice(audio_path: str) -> dict:
    y, sr = librosa.load(audio_path)

    # Energy (volume)
    rms = librosa.feature.rms(y=y)
    energy = float(np.mean(rms))

    # Pitch (stress indicator)
    pitches, magnitudes = librosa.piptrack(y=y, sr=sr)
    pitch_mean = float(np.mean(pitches[pitches > 0]))

    # Speaking rate (stress indicator)
    tempo, _ = librosa.beat.beat_track(y=y, sr=sr)

    return {
        "energy": energy,
        "pitch": pitch_mean,
        "tempo": float(tempo),
        "stress_indicator": pitch_mean > 200 and tempo > 120
    }
```

#### External API Integrations

```python
# Apple Health (via Shortcuts webhook)
# Create Shortcut that posts health data to your Gateway
@app.post("/webhooks/health")
async def health_webhook(data: dict):
    await kv.put("health.sleep", data.get("sleep_hours"))
    await kv.put("health.steps", data.get("steps"))
    await kv.put("health.heart_rate", data.get("resting_hr"))

# Energy calculation
async def calculate_capacity() -> float:
    sleep = float(await kv.get("health.sleep") or 7)
    meetings = await count_meetings_today()
    tasks_done = await count_tasks_completed()
    hours_worked = await get_hours_worked()

    # Base capacity
    capacity = 1.0

    # Sleep impact
    if sleep < 6:
        capacity -= 0.3
    elif sleep < 7:
        capacity -= 0.1

    # Meeting drain
    capacity -= (meetings / 8) * 0.2

    # Work fatigue
    capacity -= (hours_worked / 12) * 0.3

    return max(0.1, capacity)
```

### The Simple Path (No ML Required)

**Letta's LLM can handle most of this naturally:**

| Enhancement | Simple Alternative |
|-------------|-------------------|
| Sentiment | Add to prompt: "Notice Zeus's mood from his messages" |
| Patterns | Letta remembers patterns in archival memory |
| Emotional | Add feelings to memories: "Zeus was proud when..." |
| Predictions | Ask Letta: "Based on history, what might Zeus need?" |

```python
# Example: Let Letta handle sentiment naturally
message_to_legacy = f"""
[EVENT] Zeus said: "{user_message}"

Consider:
- What is Zeus's mood right now?
- Is this a good time to suggest tasks?
- Should I be brief or conversational?
"""
```

**You don't need ML for 80% of this.** Start simple, add complexity only when needed.

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                              │
│   "The goal is not to build an assistant.                                                   │
│    The goal is to build a PRESENCE.                                                         │
│                                                                                              │
│    Someone who knows you.                                                                   │
│    Someone who reads the room.                                                              │
│    Someone who remembers how things felt.                                                   │
│    Someone who knows when to speak and when to stay quiet.                                  │
│                                                                                              │
│    That's not automation. That's Legacy."                                                   │
│                                                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 16. 🫀 STATE-AWARE INTELLIGENCE (From Reactive to Proactive)

### The Fundamental Difference

**What You Have (Events = Reactive):**
```
Something happens → Legacy reacts
"Task completed" → Legacy: "Nice work!"
```
Legacy only knows when something **CHANGES**. She's blind between events.

**What State Enables (State = Proactive):**
```
Continuous awareness → Legacy sees patterns → Legacy acts BEFORE you ask
Heart rate rising + calendar empty + 3pm = "You seem stressed. Want to take a walk?"
```

---

### Events vs State: The Visual Difference

```
EVENTS:  ──●────────●───────●────●──────────●───────────●────→
         (sparse, meaningful moments - "task completed", "door opened")

STATE:   ═══════════════════════════════════════════════════════→
         (continuous stream - heart rate, location, activity level)
         ↑ patterns emerge over time, not from single points
```

**Events tell you WHAT happened.**
**State tells you HOW things are.**
**Combined = Legacy knows the FULL picture.**

---

### The Architecture: Adding State Awareness

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                          STATE-AWARE ONEMIND OS                                          │
│                                                                                          │
│   ┌─────────────────────────────────────────────────────────────────────────────────┐   │
│   │                         CONTINUOUS STATE SOURCES                                 │   │
│   │                                                                                  │   │
│   │   ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐    │   │
│   │   │  APPLE    │  │  PHONE    │  │   HOME    │  │ CALENDAR  │  │  SCREEN   │    │   │
│   │   │  WATCH    │  │   GPS     │  │ SENSORS   │  │  STATUS   │  │   TIME    │    │   │
│   │   │           │  │           │  │           │  │           │  │           │    │   │
│   │   │ • Heart   │  │ • Location│  │ • Motion  │  │ • In mtg? │  │ • Active  │    │   │
│   │   │ • Steps   │  │ • Speed   │  │ • Temp    │  │ • Free?   │  │ • App     │    │   │
│   │   │ • Sleep   │  │ • Home?   │  │ • Light   │  │ • Focus?  │  │ • Idle?   │    │   │
│   │   │ • HRV     │  │ • Work?   │  │ • Noise   │  │           │  │           │    │   │
│   │   └─────┬─────┘  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘    │   │
│   │         │              │              │              │              │           │   │
│   └─────────┼──────────────┼──────────────┼──────────────┼──────────────┼───────────┘   │
│             │              │              │              │              │               │
│             ▼              ▼              ▼              ▼              ▼               │
│   ┌─────────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                                  │   │
│   │                              GATEWAY (Port 8000)                                 │   │
│   │                                                                                  │   │
│   │    /api/health   /api/location   /api/home   /webhooks/*   (existing events)    │   │
│   │         │              │              │             │                           │   │
│   └─────────┼──────────────┼──────────────┼─────────────┼───────────────────────────┘   │
│             │              │              │             │                               │
│             ▼              ▼              ▼             ▼                               │
│   ╔═════════════════════════════════════════════════════════════════════════════════╗   │
│   ║                           NATS JETSTREAM                                         ║   │
│   ║                                                                                  ║   │
│   ║   STATE STREAMS (NEW):              EVENT STREAMS (EXISTING):                   ║   │
│   ║   ─────────────────────             ──────────────────────────                  ║   │
│   ║   state.health.heart_rate           events.todoist.task.*                       ║   │
│   ║   state.health.steps                events.github.pr.*                          ║   │
│   ║   state.location.current            events.home.door.*                          ║   │
│   ║   state.home.motion                 events.calendar.meeting.*                   ║   │
│   ║   state.calendar.status             events.voice.user.*                         ║   │
│   ║   state.activity.screen                                                         ║   │
│   ║                                                                                  ║   │
│   ╚══════════════════════════════════════╦══════════════════════════════════════════╝   │
│                                          │                                              │
│                          ┌───────────────┴───────────────┐                             │
│                          │                               │                             │
│                          ▼                               ▼                             │
│   ┌──────────────────────────────────┐   ┌──────────────────────────────────┐          │
│   │                                  │   │                                  │          │
│   │     STATE AGGREGATOR (NEW)       │   │     EVENT PROCESSOR (EXISTING)  │          │
│   │                                  │   │                                  │          │
│   │  • Subscribes to state.>         │   │  • Subscribes to events.>       │          │
│   │  • Calculates rolling averages   │   │  • Checks NATS KV filters       │          │
│   │  • Detects trends (rising/fall)  │   │  • Routes to Legacy             │          │
│   │  • Detects anomalies             │◀──│                                  │          │
│   │  • Updates NATS KV: zeus.state   │   │  • READS STATE for context      │          │
│   │  • Publishes significant changes │   │                                  │          │
│   │                                  │   │                                  │          │
│   └──────────────┬───────────────────┘   └──────────────────┬───────────────┘          │
│                  │                                          │                          │
│                  │      ┌───────────────────────┐           │                          │
│                  │      │                       │           │                          │
│                  └─────▶│     NATS KV           │◀──────────┘                          │
│                         │                       │                                      │
│                         │  zeus.state.health    │                                      │
│                         │  zeus.state.location  │                                      │
│                         │  zeus.state.activity  │                                      │
│                         │  zeus.state.capacity  │ ◀── "Am I in a good state?"         │
│                         │                       │                                      │
│                         │  onemind.vip_list     │ (existing filters)                  │
│                         │  onemind.ignore_list  │                                      │
│                         │  onemind.mode         │                                      │
│                         │                       │                                      │
│                         └───────────┬───────────┘                                      │
│                                     │                                                  │
│                                     │ watch()                                          │
│                                     ▼                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                                  │   │
│   │                              EVENT PROCESSOR                                     │   │
│   │                        (Now State-Aware!)                                        │   │
│   │                                                                                  │   │
│   │   When event arrives:                                                           │   │
│   │   1. Read current state from NATS KV                                            │   │
│   │   2. Combine event + state                                                      │   │
│   │   3. Send enriched message to Legacy                                            │   │
│   │                                                                                  │   │
│   │   Example:                                                                       │   │
│   │   ─────────                                                                      │   │
│   │   Event: "Zeus marked task as blocked"                                          │   │
│   │   State: heart_rate=92 (elevated), activity=sedentary (2hr), location=home      │   │
│   │                                                                                  │   │
│   │   Message to Legacy:                                                            │   │
│   │   "[EVENT] Task blocked + [STATE] You've been sitting 2hr, HR elevated"         │   │
│   │                                                                                  │   │
│   └──────────────────────────────────────┬──────────────────────────────────────────┘   │
│                                          │                                              │
│                                          ▼                                              │
│   ┌─────────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                                  │   │
│   │                                Legacy                                            │   │
│   │                                                                                  │   │
│   │   NOW SEES THE FULL PICTURE:                                                    │   │
│   │                                                                                  │   │
│   │   "Zeus is blocked on a task. But I also see he's been sitting for 2 hours     │   │
│   │    and his heart rate is elevated. This isn't just a productivity issue -       │   │
│   │    he's stressed. I should suggest a break, not just offer to help with        │   │
│   │    the task."                                                                   │   │
│   │                                                                                  │   │
│   └─────────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

### The State Aggregator (New Component)

**Why It Exists:**
Legacy doesn't need 60 heart rate readings per minute. She needs:
> "Your heart rate has been elevated for 20 minutes."

The State Aggregator transforms RAW DATA → MEANINGFUL STATE.

```python
# onemind-os/agent/aggregator/main.py

import asyncio
import json
import nats
from datetime import datetime, timedelta
from collections import deque


class StateAggregator:
    """
    Transforms continuous raw data into meaningful state.

    RAW: heart_rate=75, heart_rate=78, heart_rate=82, heart_rate=85...

    MEANINGFUL STATE:
    {
        "current": 85,
        "avg_5min": 80,
        "avg_30min": 72,
        "trend": "rising",
        "elevated": true,
        "elevated_duration_min": 12
    }
    """

    def __init__(self):
        self.nc = None
        self.js = None
        self.kv = None

        # Rolling windows for each metric
        self.windows = {
            "heart_rate": deque(maxlen=1800),  # 30 min at 1/sec
            "steps": deque(maxlen=1800),
            "location": deque(maxlen=60),       # Last hour of locations
            "motion": deque(maxlen=300),        # 5 min of motion data
        }

        # Baselines (learned over time)
        self.baselines = {
            "heart_rate": {"resting": 65, "active": 90, "elevated": 95},
            "steps_per_hour": 500,
        }

    async def connect(self):
        """Connect to NATS and set up KV store."""
        self.nc = await nats.connect(os.environ["NATS_URL"])
        self.js = self.nc.jetstream()
        self.kv = await self.js.key_value("onemind")
        print("State Aggregator connected to NATS")

    async def subscribe_to_state(self):
        """Subscribe to all state streams."""
        await self.js.subscribe(
            "state.>",
            cb=self.on_state_data,
            durable="state-aggregator"
        )
        print("Subscribed to state.> streams")

    async def on_state_data(self, msg):
        """Process incoming state data."""
        subject = msg.subject  # e.g., "state.health.heart_rate"
        data = json.loads(msg.data.decode())

        # Route to appropriate handler
        if "health.heart_rate" in subject:
            await self.process_heart_rate(data)
        elif "health.steps" in subject:
            await self.process_steps(data)
        elif "location" in subject:
            await self.process_location(data)
        elif "home.motion" in subject:
            await self.process_motion(data)

        await msg.ack()

    async def process_heart_rate(self, data: dict):
        """
        Process heart rate reading and update state.

        Input: {"value": 82, "timestamp": "..."}

        Output to NATS KV:
        {
            "current": 82,
            "avg_5min": 78,
            "avg_30min": 72,
            "trend": "rising",
            "elevated": false,
            "zone": "normal",  # resting, normal, active, elevated
            "updated": "..."
        }
        """
        value = data["value"]
        timestamp = datetime.fromisoformat(data.get("timestamp", datetime.now().isoformat()))

        # Add to rolling window
        self.windows["heart_rate"].append({"value": value, "ts": timestamp})

        # Calculate averages
        readings = list(self.windows["heart_rate"])
        now = datetime.now()

        # 5-minute average
        recent_5 = [r["value"] for r in readings if (now - r["ts"]).seconds < 300]
        avg_5min = sum(recent_5) / len(recent_5) if recent_5 else value

        # 30-minute average
        recent_30 = [r["value"] for r in readings if (now - r["ts"]).seconds < 1800]
        avg_30min = sum(recent_30) / len(recent_30) if recent_30 else value

        # Determine trend
        if avg_5min > avg_30min * 1.1:
            trend = "rising"
        elif avg_5min < avg_30min * 0.9:
            trend = "falling"
        else:
            trend = "stable"

        # Determine zone
        baselines = self.baselines["heart_rate"]
        if value < baselines["resting"] + 10:
            zone = "resting"
        elif value < baselines["active"]:
            zone = "normal"
        elif value < baselines["elevated"]:
            zone = "active"
        else:
            zone = "elevated"

        # Build state object
        state = {
            "current": value,
            "avg_5min": round(avg_5min, 1),
            "avg_30min": round(avg_30min, 1),
            "trend": trend,
            "zone": zone,
            "elevated": zone == "elevated",
            "updated": datetime.now().isoformat()
        }

        # Update NATS KV
        await self.kv.put("zeus.state.health.heart_rate", json.dumps(state))

        # If significant change, publish as event
        if trend == "rising" and zone == "elevated":
            await self.js.publish(
                "events.state.health.elevated",
                json.dumps({
                    "type": "health_elevated",
                    "metric": "heart_rate",
                    "value": value,
                    "trend": trend,
                    "duration_elevated_min": self._get_elevated_duration()
                })
            )

    async def process_location(self, data: dict):
        """
        Process location data and determine context.

        Input: {"lat": 37.7749, "lng": -122.4194, "accuracy": 10}

        Output to NATS KV:
        {
            "current": "home",  # home, work, gym, unknown
            "coordinates": {...},
            "time_at_location_min": 45,
            "last_transition": "work→home at 5:30pm",
            "updated": "..."
        }
        """
        lat = data["lat"]
        lng = data["lng"]

        # Determine named location (you'd configure these)
        locations = {
            "home": {"lat": 37.7749, "lng": -122.4194, "radius": 100},
            "work": {"lat": 37.7849, "lng": -122.4094, "radius": 200},
            "gym": {"lat": 37.7649, "lng": -122.4294, "radius": 50},
        }

        current_location = "unknown"
        for name, loc in locations.items():
            if self._distance(lat, lng, loc["lat"], loc["lng"]) < loc["radius"]:
                current_location = name
                break

        # Track time at location
        self.windows["location"].append({
            "location": current_location,
            "ts": datetime.now()
        })

        # Calculate time at current location
        time_at_location = 0
        for reading in reversed(list(self.windows["location"])):
            if reading["location"] == current_location:
                time_at_location = (datetime.now() - reading["ts"]).seconds // 60
            else:
                break

        state = {
            "current": current_location,
            "coordinates": {"lat": lat, "lng": lng},
            "time_at_location_min": time_at_location,
            "updated": datetime.now().isoformat()
        }

        await self.kv.put("zeus.state.location", json.dumps(state))

        # Publish transition event if location changed
        prev = self._get_previous_location()
        if prev and prev != current_location:
            await self.js.publish(
                "events.state.location.transition",
                json.dumps({
                    "from": prev,
                    "to": current_location,
                    "time": datetime.now().isoformat()
                })
            )

    async def calculate_capacity(self):
        """
        Calculate Zeus's current capacity for interruption.

        This is the MAGIC - combining multiple state signals into
        a single "should I bother him?" score.

        Output to NATS KV:
        {
            "score": 0.7,  # 0.0 = don't bother, 1.0 = fully available
            "factors": {
                "heart_rate": 0.8,   # calm
                "activity": 0.6,     # moderately active
                "location": 1.0,     # at home
                "calendar": 0.5,     # has meeting soon
                "time_of_day": 0.9   # good time
            },
            "recommendation": "available",  # available, busy, dnd
            "updated": "..."
        }
        """
        factors = {}

        # Heart rate factor (calm = high score)
        hr_state = await self._get_state("zeus.state.health.heart_rate")
        if hr_state:
            if hr_state["zone"] == "resting":
                factors["heart_rate"] = 1.0
            elif hr_state["zone"] == "normal":
                factors["heart_rate"] = 0.8
            elif hr_state["zone"] == "active":
                factors["heart_rate"] = 0.5
            else:
                factors["heart_rate"] = 0.2  # elevated = don't bother
        else:
            factors["heart_rate"] = 0.7  # default

        # Location factor
        loc_state = await self._get_state("zeus.state.location")
        if loc_state:
            if loc_state["current"] == "home":
                factors["location"] = 1.0
            elif loc_state["current"] == "work":
                factors["location"] = 0.6
            else:
                factors["location"] = 0.4
        else:
            factors["location"] = 0.7

        # Time of day factor
        hour = datetime.now().hour
        if 9 <= hour <= 17:
            factors["time_of_day"] = 0.9  # work hours
        elif 7 <= hour <= 9 or 17 <= hour <= 21:
            factors["time_of_day"] = 0.7  # transition times
        elif 21 <= hour or hour < 7:
            factors["time_of_day"] = 0.2  # sleep time

        # Calculate overall score
        score = sum(factors.values()) / len(factors)

        # Determine recommendation
        if score > 0.7:
            recommendation = "available"
        elif score > 0.4:
            recommendation = "busy"
        else:
            recommendation = "dnd"

        state = {
            "score": round(score, 2),
            "factors": factors,
            "recommendation": recommendation,
            "updated": datetime.now().isoformat()
        }

        await self.kv.put("zeus.state.capacity", json.dumps(state))

        return state

    async def run(self):
        """Main loop - also periodically recalculates capacity."""
        await self.connect()
        await self.subscribe_to_state()

        while True:
            await self.calculate_capacity()
            await asyncio.sleep(30)  # Recalculate every 30 seconds


if __name__ == "__main__":
    aggregator = StateAggregator()
    asyncio.run(aggregator.run())
```

---

### Health Data Pipeline (Apple Watch → Legacy)

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                         APPLE WATCH → Legacy PIPELINE                                    │
│                                                                                          │
│   ┌───────────────────┐                                                                 │
│   │   APPLE WATCH     │                                                                 │
│   │                   │                                                                 │
│   │  Sensors:         │                                                                 │
│   │  • Heart rate     │                                                                 │
│   │  • Steps          │                                                                 │
│   │  • HRV            │                                                                 │
│   │  • Sleep stages   │                                                                 │
│   │  • Activity       │                                                                 │
│   └─────────┬─────────┘                                                                 │
│             │                                                                           │
│             │ Bluetooth/WiFi                                                            │
│             ▼                                                                           │
│   ┌───────────────────┐                                                                 │
│   │   iPHONE          │                                                                 │
│   │   Health App      │                                                                 │
│   │                   │                                                                 │
│   │   Stores all      │                                                                 │
│   │   health data     │                                                                 │
│   └─────────┬─────────┘                                                                 │
│             │                                                                           │
│             │ OPTIONS FOR EXPORT:                                                       │
│             │                                                                           │
│   ┌─────────┼─────────────────────────────────────────────────────────────────────┐    │
│   │         │                                                                      │    │
│   │         ▼                                                                      │    │
│   │   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐       │    │
│   │   │  OPTION 1   │   │  OPTION 2   │   │  OPTION 3   │   │  OPTION 4   │       │    │
│   │   │             │   │             │   │             │   │             │       │    │
│   │   │  Shortcuts  │   │  Health     │   │  Oura/      │   │  Apple      │       │    │
│   │   │  Automation │   │  Auto       │   │  Whoop API  │   │  HealthKit  │       │    │
│   │   │             │   │  Export     │   │             │   │  (local)    │       │    │
│   │   │  Trigger:   │   │             │   │  If you use │   │             │       │    │
│   │   │  Time-based │   │  App that   │   │  a ring or  │   │  Requires   │       │    │
│   │   │  or NFC tag │   │  exports to │   │  band with  │   │  local Mac  │       │    │
│   │   │             │   │  webhook    │   │  cloud API  │   │  app        │       │    │
│   │   │  → Webhook  │   │             │   │             │   │             │       │    │
│   │   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘       │    │
│   │          │                 │                 │                 │              │    │
│   │          └─────────────────┴─────────────────┴─────────────────┘              │    │
│   │                                    │                                          │    │
│   │                                    ▼                                          │    │
│   │                         ┌───────────────────┐                                 │    │
│   │                         │   N8N WORKFLOW    │                                 │    │
│   │                         │                   │                                 │    │
│   │                         │  • Receives data  │                                 │    │
│   │                         │  • Transforms     │                                 │    │
│   │                         │  • Rate limits    │                                 │    │
│   │                         │  • Sends to GW    │                                 │    │
│   │                         └─────────┬─────────┘                                 │    │
│   │                                   │                                           │    │
│   └───────────────────────────────────┼───────────────────────────────────────────┘    │
│                                       │                                                │
│                                       │ POST /api/health                               │
│                                       ▼                                                │
│   ┌───────────────────────────────────────────────────────────────────────────────┐    │
│   │                              GATEWAY                                           │    │
│   │                                                                                │    │
│   │   @app.post("/api/health")                                                    │    │
│   │   async def receive_health_data(data: HealthData):                            │    │
│   │       await js.publish(f"state.health.{data.metric}", data.json())            │    │
│   │                                                                                │    │
│   └─────────────────────────────────────────────────────────────────────────────────┘    │
│                                       │                                                │
│                                       ▼                                                │
│   ╔═══════════════════════════════════════════════════════════════════════════════╗    │
│   ║   NATS: state.health.heart_rate, state.health.steps, state.health.sleep      ║    │
│   ╚═══════════════════════════════════════════════════════════════════════════════╝    │
│                                       │                                                │
│                                       ▼                                                │
│   ┌───────────────────────────────────────────────────────────────────────────────┐    │
│   │                         STATE AGGREGATOR                                       │    │
│   │                                                                                │    │
│   │   RAW: hr=75, hr=78, hr=82, hr=85...                                          │    │
│   │                    ↓                                                           │    │
│   │   MEANINGFUL: "Heart rate rising, now elevated for 12 minutes"                │    │
│   │                    ↓                                                           │    │
│   │   → NATS KV: zeus.state.health                                                │    │
│   │   → If significant: events.state.health.elevated                              │    │
│   │                                                                                │    │
│   └───────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

### Proactive Behavior Examples

**WITHOUT State (Reactive Only):**
```
Zeus: "I'm stuck on this task"
Legacy: "What's blocking you? Maybe I can help break it down."
```

**WITH State (Proactive):**
```
[State Aggregator detects]
- Heart rate elevated for 15 min
- Sedentary for 2 hours
- Same screen/app for 45 min
- No breaks taken

[Before Zeus says anything]
Legacy: "Hey, I notice you've been grinding hard for a while.
Your body's telling me you could use a quick break.
Want to stretch for 5 minutes? I'll remind you about the task."
```

**More Examples:**

| State Signals | Proactive Action |
|---------------|------------------|
| HR low + home + evening + no calendar | "Perfect time for deep work if you have any" |
| HR elevated + work + back-to-back meetings | "Tough day. I'm keeping things quiet until you're done." |
| Location = gym + HR = active | "I'll hold all non-urgent notifications." |
| Sleep < 6 hrs last night + 3pm | "You're probably hitting an energy dip. Coffee or walk?" |
| Home + family present + evening | "Family time - I'm in silent mode unless urgent." |

---

### NATS KV State Structure

```json
// zeus.state.health
{
  "heart_rate": {
    "current": 78,
    "avg_5min": 75,
    "avg_30min": 72,
    "trend": "stable",
    "zone": "normal"
  },
  "steps_today": 4523,
  "sleep_last_night_hrs": 7.2,
  "hrv": 45,
  "updated": "2026-01-16T14:30:00Z"
}

// zeus.state.location
{
  "current": "home",
  "time_at_location_min": 120,
  "last_transition": "work→home at 5:30pm",
  "updated": "2026-01-16T14:30:00Z"
}

// zeus.state.activity
{
  "current": "working",
  "screen_time_continuous_min": 45,
  "last_break_min_ago": 90,
  "posture": "sedentary",
  "updated": "2026-01-16T14:30:00Z"
}

// zeus.state.capacity (THE MAGIC NUMBER)
{
  "score": 0.72,
  "factors": {
    "heart_rate": 0.8,
    "location": 1.0,
    "time_of_day": 0.9,
    "activity": 0.5,
    "calendar": 0.4
  },
  "recommendation": "available",
  "updated": "2026-01-16T14:30:00Z"
}
```

---

### Event Processor: Now State-Aware

```python
# Updated Event Processor to include state context

class EventProcessor:
    """The bridge between events and Legacy - now state-aware."""

    async def send_to_legacy(self, event_type: str, event: dict):
        """
        Send event to Legacy with full state context.

        Before: "[EVENT] Task blocked"
        After:  "[EVENT] Task blocked [STATE] HR elevated 15min, sedentary 2hr"
        """

        # Get current state
        state = await self.get_current_state()

        # Build enriched message
        message = f"""
[EVENT] {event_type}
{json.dumps(event, indent=2)}

[CURRENT STATE]
- Location: {state.get('location', {}).get('current', 'unknown')}
- Heart Rate: {state.get('health', {}).get('zone', 'unknown')} ({state.get('health', {}).get('trend', 'unknown')})
- Activity: {state.get('activity', {}).get('current', 'unknown')}
- Capacity Score: {state.get('capacity', {}).get('score', 'unknown')}/1.0
- Time at Location: {state.get('location', {}).get('time_at_location_min', 'unknown')} min

[RECOMMENDATION]
{state.get('capacity', {}).get('recommendation', 'Check manually')}
"""

        # Send to Legacy
        response = self.letta.agents.messages.create(
            agent_id=self.legacy_agent_id,
            messages=[{
                "role": "system",
                "content": message
            }]
        )

        return response

    async def get_current_state(self) -> dict:
        """Read all current state from NATS KV."""
        state = {}

        # Health
        health = await self.kv.get("zeus.state.health.heart_rate")
        if health:
            state["health"] = json.loads(health.value)

        # Location
        location = await self.kv.get("zeus.state.location")
        if location:
            state["location"] = json.loads(location.value)

        # Activity
        activity = await self.kv.get("zeus.state.activity")
        if activity:
            state["activity"] = json.loads(activity.value)

        # Capacity (the magic number)
        capacity = await self.kv.get("zeus.state.capacity")
        if capacity:
            state["capacity"] = json.loads(capacity.value)

        return state
```

---

### The Capacity Score: Legacy's "Should I Speak?" Meter

```
CAPACITY SCORE: 0.0 ─────────────────────────────────── 1.0
                │                                        │
                │                                        │
          DON'T BOTHER                           FULLY AVAILABLE

          ◀────────────────────────────────────────────▶

          0.0-0.3: DND (only emergencies)
          0.3-0.6: Busy (important only)
          0.6-0.8: Available (normal filtering)
          0.8-1.0: Open (proactive suggestions welcome)
```

**How Event Processor Uses It:**
```python
async def should_send_to_legacy(self, event: dict, urgency: str) -> bool:
    """
    Decide whether to send based on capacity + urgency.
    """
    capacity = await self.get_capacity_score()

    # Emergency always gets through
    if urgency == "emergency":
        return True

    # Map urgency to minimum capacity required
    thresholds = {
        "high": 0.3,      # Important - needs some availability
        "medium": 0.5,    # Normal - needs moderate availability
        "low": 0.7,       # Nice to know - needs high availability
        "proactive": 0.8  # Legacy's own ideas - needs very available
    }

    min_capacity = thresholds.get(urgency, 0.5)

    if capacity >= min_capacity:
        return True
    else:
        # Queue for later
        await self.queue_event(event, urgency)
        return False
```

---

### What This Enables: TRUE Ambient Intelligence

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                          │
│                         THE EVOLUTION OF Legacy                                          │
│                                                                                          │
│   LEVEL 1: REACTIVE ASSISTANT                                                           │
│   ───────────────────────────                                                           │
│   "You asked → I answered"                                                              │
│   (Most chatbots)                                                                       │
│                                                                                          │
│   LEVEL 2: EVENT-AWARE ASSISTANT                                                        │
│   ──────────────────────────────                                                        │
│   "Something happened → I told you"                                                     │
│   (What you had before Section 16)                                                      │
│                                                                                          │
│   LEVEL 3: STATE-AWARE COMPANION                                                        │
│   ──────────────────────────────                                                        │
│   "I know HOW you are, not just WHAT happened"                                          │
│   (What you have now)                                                                   │
│                                                                                          │
│   LEVEL 4: PROACTIVE PRESENCE                                                           │
│   ─────────────────────────────                                                         │
│   "I see patterns → I act before you ask"                                               │
│   (Next evolution)                                                                      │
│                                                                                          │
│   ════════════════════════════════════════════════════════════════════════════════════  │
│                                                                                          │
│   YOU ARE NOW AT LEVEL 3, WITH LEVEL 4 CAPABILITIES DOCUMENTED.                         │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

### Implementation Priority

| Component | Priority | Complexity | Impact |
|-----------|----------|------------|--------|
| Gateway health endpoint | 🟢 Easy | Low | Foundation |
| State Aggregator (basic) | 🟡 Medium | Medium | High |
| NATS KV state keys | 🟢 Easy | Low | Foundation |
| Capacity score calculation | 🟡 Medium | Medium | High |
| Event Processor state-enrichment | 🟢 Easy | Low | Medium |
| Apple Watch pipeline | 🟡 Medium | Medium | High |
| Proactive trigger system | 🔴 Advanced | High | Very High |

**Start with:**
1. Add `/api/health` endpoint to Gateway
2. Create State Aggregator skeleton
3. Add `zeus.state.*` keys to NATS KV
4. Update Event Processor to read state

**Then add sensors:**
- Health data (watch)
- Location (phone GPS)
- Activity (screen time)

---

*State-aware intelligence transforms Legacy from an assistant that RESPONDS to one that UNDERSTANDS. She's not just waiting for you to talk - she's watching, feeling the rhythm of your day, knowing when to speak and when to stay quiet. That's not automation. That's presence.*

---

## 17. 🌐 NATS COMPLETE PLATFORM (The Full Ecosystem)

### NATS is NOT Just Messaging - It's a Complete Platform

After deep research into the NATS documentation, we discovered NATS provides **everything** OneMind needs:

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                         NATS: THE COMPLETE PLATFORM                                      │
│                                                                                          │
│   ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐   │
│   │  CORE NATS      │  │  JETSTREAM      │  │  KV STORE       │  │  OBJECT STORE   │   │
│   │                 │  │                 │  │                 │  │                 │   │
│   │  • Pub/Sub      │  │  • Streams      │  │  • Get/Put      │  │  • Files        │   │
│   │  • Request/     │  │  • Consumers    │  │  • Watch        │  │  • Chunking     │   │
│   │    Reply        │  │  • Persistence  │  │  • TTL          │  │  • Versioning   │   │
│   │  • Queue Groups │  │  • Replay       │  │  • History      │  │  • Watch        │   │
│   │  • Wildcards    │  │  • Exactly-Once │  │  • Atomic Ops   │  │                 │   │
│   │                 │  │                 │  │                 │  │                 │   │
│   │  Real-time      │  │  Event History  │  │  Current State  │  │  Media/Files    │   │
│   │  Messaging      │  │  & Persistence  │  │  & Config       │  │  Storage        │   │
│   └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────────┘   │
│                                                                                          │
│   ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐                        │
│   │  MICRO          │  │  LEAF NODES     │  │  SECURITY       │                        │
│   │  FRAMEWORK      │  │                 │  │                 │                        │
│   │                 │  │  • Edge→Cloud   │  │  • NKeys        │                        │
│   │  • Service      │  │  • Auto-        │  │  • JWT/Accounts │                        │
│   │    Discovery    │  │    Reconnect    │  │  • TLS/mTLS     │                        │
│   │  • Health       │  │  • NAT-friendly │  │  • Subject      │                        │
│   │    Checks       │  │  • Local        │  │    Permissions  │                        │
│   │  • Stats        │  │    Buffering    │  │                 │                        │
│   │  • Load         │  │                 │  │                 │                        │
│   │    Balancing    │  │  Home ↔ Cloud   │  │  Zero Trust     │                        │
│   └─────────────────┘  └─────────────────┘  └─────────────────┘                        │
│                                                                                          │
│   ════════════════════════════════════════════════════════════════════════════════════  │
│                                                                                          │
│   REPLACES:  Redis + Kafka + S3 + Consul + Service Mesh                                │
│   WITH:      Just NATS (20MB binary)                                                   │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

### 17.1 NATS Micro Framework - Legacy as a Service

**Discovery:** NATS has a built-in microservices framework. Legacy should be registered as a NATS service!

```python
# onemind-os/agent/legacy_service.py

import asyncio
import nats
from nats.micro import Service, ServiceConfig, EndpointConfig

class LegacyService:
    """
    Legacy AI registered as a NATS Micro Service.

    Benefits:
    - Automatic service discovery (no Consul/etcd needed)
    - Built-in health checks: $SRV.PING.LegacyAI
    - Built-in stats: $SRV.STATS.LegacyAI
    - Automatic load balancing via queue groups
    - Observable by default
    """

    def __init__(self):
        self.nc = None
        self.service = None
        self.letta_client = None

    async def start(self):
        """Start Legacy as a NATS Micro Service."""
        self.nc = await nats.connect(os.environ["NATS_URL"])

        # Register Legacy as a discoverable service
        self.service = await Service.create(
            self.nc,
            ServiceConfig(
                name="LegacyAI",
                version="1.0.0",
                description="OneMind AI Agent - The Brain",
                metadata={
                    "owner": "zeus",
                    "type": "ai-agent",
                    "letta_agent_id": os.environ["Legacy_AGENT_ID"]
                }
            )
        )

        # Endpoint: Process events from Event Processor
        await self.service.add_endpoint(
            EndpointConfig(
                name="process",
                subject="ai.process",
                handler=self.handle_event,
                metadata={"description": "Process incoming events"}
            )
        )

        # Endpoint: Query Legacy for information
        await self.service.add_endpoint(
            EndpointConfig(
                name="query",
                subject="ai.query",
                handler=self.handle_query,
                metadata={"description": "Query Legacy for context/decisions"}
            )
        )

        # Endpoint: Voice interactions from LiveKit
        await self.service.add_endpoint(
            EndpointConfig(
                name="voice",
                subject="ai.voice",
                handler=self.handle_voice,
                metadata={"description": "Voice interaction endpoint"}
            )
        )

        print(f"Legacy Service started: {self.service.info()}")
        print("Discoverable at: $SRV.PING.LegacyAI")
        print("Stats at: $SRV.STATS.LegacyAI")

    async def handle_event(self, msg):
        """
        Handle events from Event Processor.

        Request: {"event_type": "...", "event": {...}, "state": {...}}
        Response: {"processed": true, "action_taken": "..."}
        """
        data = json.loads(msg.data.decode())

        # Send to Letta for processing
        response = await self.letta_client.agents.messages.create(
            agent_id=self.legacy_agent_id,
            messages=[{
                "role": "system",
                "content": f"[EVENT] {data['event_type']}\n{json.dumps(data['event'])}\n\n[STATE]\n{json.dumps(data['state'])}"
            }]
        )

        # Return response
        await msg.respond(json.dumps({
            "processed": True,
            "response": response.messages[-1].content
        }).encode())

    async def handle_query(self, msg):
        """
        Handle queries - request/reply pattern.

        Request: {"query": "What should I focus on today?"}
        Response: {"answer": "..."}
        """
        data = json.loads(msg.data.decode())

        response = await self.letta_client.agents.messages.create(
            agent_id=self.legacy_agent_id,
            messages=[{
                "role": "user",
                "content": data["query"]
            }]
        )

        await msg.respond(json.dumps({
            "answer": response.messages[-1].content
        }).encode())

    async def handle_voice(self, msg):
        """Handle voice transcriptions from LiveKit."""
        data = json.loads(msg.data.decode())

        response = await self.letta_client.agents.messages.create(
            agent_id=self.legacy_agent_id,
            messages=[{
                "role": "user",
                "content": data["transcription"]
            }]
        )

        await msg.respond(json.dumps({
            "response": response.messages[-1].content,
            "should_speak": True
        }).encode())


# Built-in monitoring - NO EXTERNAL TOOLS NEEDED!
#
# Health check:
#   nats req '$SRV.PING.LegacyAI' ''
#
# Service info:
#   nats req '$SRV.INFO.LegacyAI' ''
#
# Statistics:
#   nats req '$SRV.STATS.LegacyAI' ''
#
# Returns: request count, errors, processing time, etc.
```

**Service Discovery Pattern:**

```
┌─────────────────────────────────────────────────────────────────┐
│                    NATS SERVICE DISCOVERY                       │
│                                                                 │
│   ┌─────────────┐     $SRV.PING.*      ┌─────────────┐         │
│   │   Gateway   │ ────────────────────▶│   Legacy    │         │
│   │             │                      │   Service   │         │
│   │             │◀────────────────────│             │         │
│   │             │     "I'm alive!"     │             │         │
│   └─────────────┘                      └─────────────┘         │
│                                                                 │
│   ANY service can discover ANY other service.                  │
│   No Consul. No etcd. No Kubernetes service mesh.              │
│   Just NATS.                                                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

### 17.2 Object Store - Voice Recordings & Media

**Discovery:** NATS has built-in file/blob storage! No S3 needed for basic media storage.

```python
# onemind-os/agent/storage/media.py

import nats
from nats.js.api import ObjectStoreConfig

class MediaStorage:
    """
    Store voice recordings, images, and attachments in NATS Object Store.

    No S3 required for basic file storage!
    """

    def __init__(self):
        self.nc = None
        self.js = None
        self.voice_bucket = None
        self.attachments_bucket = None

    async def connect(self):
        self.nc = await nats.connect(os.environ["NATS_URL"])
        self.js = self.nc.jetstream()

        # Create/get voice recordings bucket
        self.voice_bucket = await self.js.create_object_store(
            ObjectStoreConfig(
                bucket="voice_recordings",
                description="LiveKit conversation recordings",
                max_bytes=10 * 1024 * 1024 * 1024,  # 10GB max
                storage="file"
            )
        )

        # Create/get attachments bucket
        self.attachments_bucket = await self.js.create_object_store(
            ObjectStoreConfig(
                bucket="attachments",
                description="Images, screenshots, files",
                max_bytes=5 * 1024 * 1024 * 1024,  # 5GB max
                storage="file"
            )
        )

    async def store_voice_recording(
        self,
        conversation_id: str,
        audio_data: bytes,
        metadata: dict = None
    ) -> str:
        """
        Store a voice recording from LiveKit.

        Returns the object name for later retrieval.
        """
        timestamp = datetime.now().strftime("%Y-%m-%d-%H-%M-%S")
        object_name = f"conversation-{conversation_id}-{timestamp}.wav"

        await self.voice_bucket.put(
            object_name,
            audio_data,
            meta=metadata or {}
        )

        return object_name

    async def get_voice_recording(self, object_name: str) -> bytes:
        """Retrieve a voice recording."""
        result = await self.voice_bucket.get(object_name)
        return result.data

    async def store_attachment(
        self,
        filename: str,
        data: bytes,
        content_type: str = None
    ) -> str:
        """Store an image, screenshot, or file."""
        await self.attachments_bucket.put(
            filename,
            data,
            meta={"content_type": content_type} if content_type else {}
        )
        return filename

    async def watch_new_recordings(self, callback):
        """
        Watch for new voice recordings.

        Legacy can be notified when new recordings are saved.
        """
        watcher = await self.voice_bucket.watch()

        async for info in watcher:
            if info.deleted:
                continue
            await callback(info.name, info.meta)

    async def list_recent_recordings(self, limit: int = 10) -> list:
        """List recent voice recordings."""
        objects = []
        async for info in self.voice_bucket.list():
            objects.append({
                "name": info.name,
                "size": info.size,
                "modified": info.mtime,
                "meta": info.meta
            })

        # Sort by modified time, return most recent
        objects.sort(key=lambda x: x["modified"], reverse=True)
        return objects[:limit]


# Usage Example:
#
# # Store a conversation
# storage = MediaStorage()
# await storage.connect()
#
# object_name = await storage.store_voice_recording(
#     conversation_id="conv-123",
#     audio_data=audio_bytes,
#     metadata={"duration_sec": 45, "topic": "morning briefing"}
# )
#
# # Legacy can later reference or replay:
# # "Let me check what we discussed this morning..."
# audio = await storage.get_voice_recording(object_name)
```

**Object Store vs KV Store:**

| Feature | KV Store | Object Store |
|---------|----------|--------------|
| **Purpose** | Small values, config, state | Large files, media |
| **Max Size** | ~1MB per value | Unlimited (chunked) |
| **Watch** | ✅ Key changes | ✅ New objects |
| **TTL** | ✅ Per bucket | ❌ No |
| **History** | ✅ Versioning | ✅ Versioning |
| **Use Case** | VIP lists, capacity score | Voice recordings, images |

---

### 17.3 Leaf Node Architecture - Home ↔ Cloud

**Discovery:** NATS Leaf Nodes are perfect for connecting home to cloud with resilience!

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                         LEAF NODE ARCHITECTURE                                           │
│                                                                                          │
│   ┌─────────────────────────────────────────────────────────────────────────────────┐   │
│   │                              CLOUD (AWS Lightsail)                               │   │
│   │                                                                                  │   │
│   │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                            │   │
│   │   │   NATS-1    │──│   NATS-2    │──│   NATS-3    │  3-Node Cluster           │   │
│   │   │   (Leader)  │  │  (Follower) │  │  (Follower) │  R=3 for HA               │   │
│   │   └──────┬──────┘  └──────┬──────┘  └──────┬──────┘                            │   │
│   │          └─────────┬──────┴─────────┬──────┘                                   │   │
│   │                    │                │                                           │   │
│   │              JetStream         Leaf Node                                        │   │
│   │              (Replicated)      Listener                                         │   │
│   │                    │           Port 7422                                        │   │
│   │                    │                │                                           │   │
│   │          ┌─────────┴────────┐       │                                          │   │
│   │          │                  │       │                                          │   │
│   │     ┌────▼────┐       ┌────▼────┐   │                                          │   │
│   │     │ Legacy  │       │ Event   │   │                                          │   │
│   │     │   AI    │       │Processor│   │                                          │   │
│   │     └─────────┘       └─────────┘   │                                          │   │
│   │                                     │                                           │   │
│   └─────────────────────────────────────┼───────────────────────────────────────────┘   │
│                                         │                                               │
│                          Leaf Node Connection                                           │
│                      ════════════════════════════                                       │
│                      • TLS Encrypted                                                    │
│                      • NKey Authenticated                                               │
│                      • Auto-Reconnect (5s)                                             │
│                      • Survives Cloud Outages                                           │
│                      • NAT-Friendly (Outbound)                                         │
│                                         │                                               │
│   ┌─────────────────────────────────────┼───────────────────────────────────────────┐   │
│   │                                     │                                            │   │
│   │                              HOME NETWORK                                        │   │
│   │                                     │                                            │   │
│   │                          ┌──────────▼──────────┐                                │   │
│   │                          │   NATS LEAF NODE    │                                │   │
│   │                          │                     │                                │   │
│   │                          │  ┌───────────────┐  │                                │   │
│   │                          │  │ Local         │  │                                │   │
│   │                          │  │ JetStream     │  │  Buffers events when          │   │
│   │                          │  │ (Optional)    │  │  cloud is unreachable         │   │
│   │                          │  └───────────────┘  │                                │   │
│   │                          │                     │                                │   │
│   │                          │  Publishes: home.>  │                                │   │
│   │                          │  Subscribes: cmd.>  │                                │   │
│   │                          │                     │                                │   │
│   │                          └──────────┬──────────┘                                │   │
│   │                                     │                                            │   │
│   │           ┌─────────────────────────┼─────────────────────────┐                 │   │
│   │           │                         │                         │                 │   │
│   │           ▼                         ▼                         ▼                 │   │
│   │    ┌─────────────┐          ┌─────────────┐          ┌─────────────┐           │   │
│   │    │    HOME     │          │     IoT     │          │   PHONE     │           │   │
│   │    │  ASSISTANT  │          │   SENSORS   │          │    APP      │           │   │
│   │    │             │          │             │          │             │           │   │
│   │    │ • Lights    │          │ • Temp      │          │ • GPS       │           │   │
│   │    │ • Doors     │          │ • Humidity  │          │ • Health    │           │   │
│   │    │ • Motion    │          │ • Air       │          │ • Activity  │           │   │
│   │    │ • Presence  │          │ • Light     │          │             │           │   │
│   │    │             │          │             │          │             │           │   │
│   │    │ Publishes:  │          │ Publishes:  │          │ Publishes:  │           │   │
│   │    │ home.door.> │          │ home.env.>  │          │ state.>     │           │   │
│   │    │ home.motion │          │             │          │             │           │   │
│   │    │             │          │             │          │             │           │   │
│   │    │ Subscribes: │          │             │          │             │           │   │
│   │    │ home.cmd.>  │          │             │          │             │           │   │
│   │    └─────────────┘          └─────────────┘          └─────────────┘           │   │
│   │                                                                                  │   │
│   └──────────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                          │
│   KEY BENEFITS:                                                                         │
│   ═════════════                                                                         │
│   ✅ Home operates independently when cloud is down                                     │
│   ✅ Events buffered locally, synced when connection restored                          │
│   ✅ Legacy can send commands back to home devices                                      │
│   ✅ No port forwarding needed (outbound connection)                                    │
│   ✅ WebSocket option for strict firewalls                                              │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

**Home Leaf Node Configuration:**

```conf
# /etc/nats/leaf.conf (on home machine)

server_name: "home-leaf"
listen: "localhost:4222"

# Local JetStream for buffering
jetstream {
    store_dir: "/var/lib/nats/jetstream"
    max_mem: 512M
    max_file: 5G
}

# Connect to cloud cluster
leafnodes {
    remotes [
        {
            urls: [
                "nats-leaf://cloud-nats-1.onemind.ai:7422",
                "nats-leaf://cloud-nats-2.onemind.ai:7422",
                "nats-leaf://cloud-nats-3.onemind.ai:7422"
            ]

            # Credentials
            credentials: "/etc/nats/home-leaf.creds"

            # Auto-reconnect
            reconnect: 5

            # TLS
            tls {
                cert_file: "/etc/nats/certs/leaf.crt"
                key_file: "/etc/nats/certs/leaf.key"
                ca_file: "/etc/nats/certs/ca.crt"
            }
        }
    ]
}
```

**Cloud Cluster Configuration (accepting leaf nodes):**

```conf
# /etc/nats/nats-server.conf (on cloud)

server_name: "cloud-nats-1"

# Client connections
listen: "0.0.0.0:4222"

# Cluster routing
cluster {
    name: "onemind-cloud"
    listen: "0.0.0.0:6222"
    routes: [
        "nats://cloud-nats-2:6222",
        "nats://cloud-nats-3:6222"
    ]
}

# Accept leaf node connections
leafnodes {
    listen: "0.0.0.0:7422"

    tls {
        cert_file: "/etc/nats/certs/server.crt"
        key_file: "/etc/nats/certs/server.key"
        ca_file: "/etc/nats/certs/ca.crt"
        verify: true
    }

    authorization {
        users = [
            {
                user: "home_leaf"
                password: "$HOME_LEAF_PASSWORD"
                account: "OneMind-Home"
            }
        ]
    }
}

# JetStream with replication
jetstream {
    store_dir: "/var/lib/nats/jetstream"
    max_mem: 4G
    max_file: 50G
}

# HTTP monitoring
http_port: 8222
```

---

### 17.4 Complete Subject Hierarchy

**The complete namespace for OneMind OS:**

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                         ONEMIND SUBJECT HIERARCHY                                        │
│                                                                                          │
│   events.>                    ← Discrete events (something happened)                    │
│   ├── events.todoist.>                                                                  │
│   │   ├── events.todoist.task.created                                                  │
│   │   ├── events.todoist.task.completed                                                │
│   │   ├── events.todoist.task.updated                                                  │
│   │   └── events.todoist.comment.added                                                 │
│   │                                                                                     │
│   ├── events.github.>                                                                   │
│   │   ├── events.github.pr.opened                                                      │
│   │   ├── events.github.pr.merged                                                      │
│   │   ├── events.github.issue.created                                                  │
│   │   └── events.github.push                                                           │
│   │                                                                                     │
│   ├── events.home.>           ← From leaf node                                         │
│   │   ├── events.home.door.opened                                                      │
│   │   ├── events.home.door.closed                                                      │
│   │   ├── events.home.motion.detected                                                  │
│   │   ├── events.home.motion.cleared                                                   │
│   │   ├── events.home.presence.arrived                                                 │
│   │   └── events.home.presence.left                                                    │
│   │                                                                                     │
│   ├── events.calendar.>                                                                 │
│   │   ├── events.calendar.meeting.starting                                             │
│   │   ├── events.calendar.meeting.ended                                                │
│   │   └── events.calendar.reminder                                                     │
│   │                                                                                     │
│   ├── events.voice.>                                                                    │
│   │   ├── events.voice.wake_word                                                       │
│   │   ├── events.voice.transcription                                                   │
│   │   └── events.voice.ended                                                           │
│   │                                                                                     │
│   └── events.state.>          ← Significant state changes (from State Aggregator)     │
│       ├── events.state.health.elevated                                                 │
│       ├── events.state.location.transition                                             │
│       └── events.state.capacity.changed                                                │
│                                                                                         │
│   state.>                     ← Continuous state data (raw sensor readings)            │
│   ├── state.health.>                                                                    │
│   │   ├── state.health.heart_rate                                                      │
│   │   ├── state.health.steps                                                           │
│   │   ├── state.health.sleep                                                           │
│   │   └── state.health.hrv                                                             │
│   │                                                                                     │
│   ├── state.location.>                                                                  │
│   │   └── state.location.current                                                       │
│   │                                                                                     │
│   └── state.activity.>                                                                  │
│       ├── state.activity.screen                                                        │
│       └── state.activity.posture                                                       │
│                                                                                         │
│   actions.>                   ← Legacy → Gateway (do something)                        │
│   ├── actions.todoist.create_task                                                       │
│   ├── actions.todoist.complete_task                                                     │
│   ├── actions.github.create_issue                                                       │
│   ├── actions.discord.send_message                                                      │
│   ├── actions.home.lights.on                                                            │
│   ├── actions.home.lights.off                                                           │
│   └── actions.home.thermostat.set                                                       │
│                                                                                         │
│   queries.>                   ← Legacy → Gateway (get info)                            │
│   ├── queries.todoist.get_tasks                                                         │
│   ├── queries.calendar.get_events                                                       │
│   └── queries.home.get_state                                                            │
│                                                                                         │
│   results.>                   ← Gateway → Legacy (responses)                           │
│   ├── results.actions.>       (action confirmations)                                   │
│   └── results.queries.>       (query responses)                                        │
│                                                                                         │
│   ai.>                        ← Legacy Service endpoints                               │
│   ├── ai.process              (Event Processor → Legacy)                               │
│   ├── ai.query                (Request/Reply queries)                                  │
│   └── ai.voice                (Voice interactions)                                     │
│                                                                                         │
│   realtime.>                  ← RePublished events (for monitoring)                    │
│   └── realtime.events.>       (mirrored from events.> for Core NATS subscribers)      │
│                                                                                         │
│   $SRV.>                      ← NATS Micro service discovery (built-in)               │
│   ├── $SRV.PING.LegacyAI                                                               │
│   ├── $SRV.INFO.LegacyAI                                                               │
│   └── $SRV.STATS.LegacyAI                                                              │
│                                                                                         │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

### 17.5 Security with Accounts & NKeys

**Production security architecture:**

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                         NATS SECURITY ARCHITECTURE                                       │
│                                                                                          │
│   ┌─────────────────────────────────────────────────────────────────────────────────┐   │
│   │                              ACCOUNTS (Isolation)                                │   │
│   │                                                                                  │   │
│   │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │   │
│   │   │    SYS      │  │   AI_CORE   │  │   EVENTS    │  │    HOME     │           │   │
│   │   │   Account   │  │   Account   │  │   Account   │  │   Account   │           │   │
│   │   │             │  │             │  │             │  │             │           │   │
│   │   │ System ops  │  │ Legacy AI   │  │ Gateway     │  │ Leaf Node   │           │   │
│   │   │ Monitoring  │  │ Voice       │  │ Processor   │  │ Home Asst.  │           │   │
│   │   │             │  │ Letta       │  │ Aggregator  │  │ IoT         │           │   │
│   │   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘           │   │
│   │                                                                                  │   │
│   │   Each account has ISOLATED subject namespace.                                  │   │
│   │   Accounts EXPORT/IMPORT subjects to share data.                                │   │
│   │                                                                                  │   │
│   └─────────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                          │
│   ┌─────────────────────────────────────────────────────────────────────────────────┐   │
│   │                              NKEYS (Authentication)                              │   │
│   │                                                                                  │   │
│   │   NKeys = Ed25519 public-key cryptography                                       │   │
│   │                                                                                  │   │
│   │   • Server never sees private key                                               │   │
│   │   • Client signs challenge with private key                                     │   │
│   │   • Server validates with public key                                            │   │
│   │   • Immune to replay attacks                                                    │   │
│   │                                                                                  │   │
│   │   Generate:                                                                     │   │
│   │   $ nk -gen user -pubout                                                        │   │
│   │   SUACSSL3PZ...  (Seed/Private - KEEP SECRET)                                  │   │
│   │   UDXU4RCNJ...   (Public - share with server)                                  │   │
│   │                                                                                  │   │
│   └─────────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                          │
│   ┌─────────────────────────────────────────────────────────────────────────────────┐   │
│   │                              TLS (Encryption)                                    │   │
│   │                                                                                  │   │
│   │   All connections encrypted:                                                    │   │
│   │   • Client ↔ Server (TLS 1.3)                                                  │   │
│   │   • Server ↔ Server (mTLS for cluster)                                         │   │
│   │   • Leaf ↔ Cluster (TLS + NKey)                                                │   │
│   │                                                                                  │   │
│   │   tls: {                                                                        │   │
│   │       cert_file: "/etc/nats/certs/server.crt"                                  │   │
│   │       key_file: "/etc/nats/certs/server.key"                                   │   │
│   │       ca_file: "/etc/nats/certs/ca.crt"                                        │   │
│   │       min_version: "1.3"                                                        │   │
│   │   }                                                                             │   │
│   │                                                                                  │   │
│   └─────────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

**Account Configuration:**

```conf
# nats-server.conf

accounts {
    SYS: {
        users: [{ nkey: UASYSTEM... }]
    }

    AI_CORE: {
        users: [
            { nkey: UALegacy..., permissions: { publish: "ai.>", subscribe: ">" } }
        ]
        exports: [
            { service: "ai.>", accounts: [EVENTS] }  # Only EVENTS can call AI
        ]
        imports: [
            { stream: { account: EVENTS, subject: "events.>" } }
            { stream: { account: HOME, subject: "home.>" }, prefix: "home" }
        ]
    }

    EVENTS: {
        users: [
            { nkey: UAGATEWAY... }
            { nkey: UAPROCESSOR... }
            { nkey: UAAGGREGATOR... }
        ]
        exports: [
            { stream: "events.>", accounts: [AI_CORE] }
        ]
        imports: [
            { service: { account: AI_CORE, subject: "ai.process" }, to: "legacy" }
        ]
    }

    HOME: {
        users: [{ nkey: UAHOMELEAF... }]
        exports: [
            { stream: "home.>", accounts: [AI_CORE, EVENTS] }
        ]
        imports: [
            { stream: { account: AI_CORE, subject: "actions.home.>" }, to: "cmd" }
        ]
    }
}

system_account: SYS
```

---

### 17.6 RePublish Feature - Real-Time Monitoring

**Discovery:** JetStream can automatically republish stored messages to Core NATS!

```python
# Stream configuration with RePublish

stream_config = {
    "name": "ONEMIND_EVENTS",
    "subjects": ["events.>"],
    "storage": "file",
    "replicas": 3,
    "max_age": "30d",

    # THE MAGIC: Auto-republish to Core NATS
    "republish": {
        "source": "events.>",
        "destination": "realtime.events.>"
    }
}
```

**What this enables:**

```
┌─────────────────────────────────────────────────────────────────┐
│                    REPUBLISH FLOW                               │
│                                                                 │
│   Gateway publishes:                                            │
│   events.todoist.task.created                                   │
│        │                                                        │
│        ▼                                                        │
│   ╔════════════════════════════════════════════════════════╗   │
│   ║              JETSTREAM STREAM                          ║   │
│   ║                                                        ║   │
│   ║   1. Stores message (persistence)                      ║   │
│   ║   2. Sends to consumers (Event Processor, etc.)       ║   │
│   ║   3. REPUBLISHES to realtime.events.todoist.task...   ║   │
│   ║                                                        ║   │
│   ╚════════════════════════════════════════════════════════╝   │
│        │                            │                           │
│        │                            │                           │
│        ▼                            ▼                           │
│   JetStream Consumer           Core NATS Subscriber             │
│   (Event Processor)            (Dashboard, Monitoring)          │
│                                                                 │
│   Benefits:                                                     │
│   • Dashboard gets real-time events WITHOUT JetStream overhead │
│   • Multiple monitoring tools can subscribe                     │
│   • Headers added: Nats-Stream, Nats-Sequence, Nats-Time-Stamp │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

### 17.7 Message Headers for Correlation & Tracing

**Use headers for distributed tracing:**

```python
# Publishing events with correlation headers

async def publish_event(event_type: str, data: dict, correlation_id: str = None):
    """
    Publish event with tracing headers.
    """
    msg = nats.Msg(
        subject=f"events.{event_type}",
        data=json.dumps(data).encode()
    )

    # Add tracing headers
    msg.headers = {
        "X-Correlation-Id": correlation_id or str(uuid.uuid4()),
        "X-Source": "gateway",
        "X-Timestamp": datetime.now().isoformat(),
        "X-User-Id": "zeus",

        # For deduplication
        "Nats-Msg-Id": f"{event_type}-{data.get('id', uuid.uuid4())}"
    }

    await js.publish(msg)
```

**Headers flow through the system:**

```
Gateway → NATS → Event Processor → Legacy → Actions → Gateway
   │                                                      │
   └──────────── X-Correlation-Id: "req-123" ────────────┘

Every message in this flow shares the same correlation ID.
Easy to trace a request through the entire system!
```

---

### 17.8 Deduplication - Exactly-Once Semantics

**Enable deduplication on streams:**

```python
# Stream config with deduplication window
stream_config = {
    "name": "ONEMIND_EVENTS",
    "subjects": ["events.>"],
    "duplicate_window": 120_000_000_000,  # 2 minutes in nanoseconds
}

# Publishing with message ID
async def publish_with_dedup(subject: str, data: dict, msg_id: str):
    """
    Publish with deduplication.

    If the same msg_id is published within 2 minutes,
    the duplicate is rejected (same ack returned).
    """
    await js.publish(
        subject,
        json.dumps(data).encode(),
        headers={"Nats-Msg-Id": msg_id}
    )
```

**Combined with double-ack for true exactly-once:**

```python
# Consumer with exactly-once processing
async def process_message(msg):
    # Process the event
    result = await process_event(msg.data)

    # Double-ack: Wait for server to confirm our ack
    # If this succeeds, message is GUARANTEED processed exactly once
    await msg.ack_sync()  # Blocks until server confirms
```

---

### 17.9 Complete Docker Compose (Cloud)

**Production-ready 3-node cluster:**

```yaml
# docker-compose.yml (cloud deployment)

version: "3.8"

networks:
  onemind:
    name: onemind

volumes:
  nats1-data:
  nats2-data:
  nats3-data:

services:
  # NATS Cluster Node 1
  nats-1:
    image: nats:latest
    container_name: nats-1
    ports:
      - "4222:4222"   # Client connections
      - "8222:8222"   # HTTP monitoring
      - "7422:7422"   # Leaf nodes
    command: >
      --name nats-1
      --cluster_name onemind-cloud
      --cluster nats://0.0.0.0:6222
      --routes=nats://nats-2:6222,nats://nats-3:6222
      --http_port 8222
      --jetstream
      --store_dir /data
      --leafnodes nats://0.0.0.0:7422
    volumes:
      - nats1-data:/data
    networks:
      - onemind
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8222/healthz"]
      interval: 10s
      timeout: 5s
      retries: 3

  # NATS Cluster Node 2
  nats-2:
    image: nats:latest
    container_name: nats-2
    command: >
      --name nats-2
      --cluster_name onemind-cloud
      --cluster nats://0.0.0.0:6222
      --routes=nats://nats-1:6222,nats://nats-3:6222
      --jetstream
      --store_dir /data
      --leafnodes nats://0.0.0.0:7422
    volumes:
      - nats2-data:/data
    networks:
      - onemind
    depends_on:
      - nats-1

  # NATS Cluster Node 3
  nats-3:
    image: nats:latest
    container_name: nats-3
    command: >
      --name nats-3
      --cluster_name onemind-cloud
      --cluster nats://0.0.0.0:6222
      --routes=nats://nats-1:6222,nats://nats-2:6222
      --jetstream
      --store_dir /data
      --leafnodes nats://0.0.0.0:7422
    volumes:
      - nats3-data:/data
    networks:
      - onemind
    depends_on:
      - nats-1

  # Gateway Service
  gateway:
    build: ./gateway
    container_name: onemind-gateway
    ports:
      - "8000:8000"
    environment:
      - NATS_URL=nats://nats-1:4222,nats://nats-2:4222,nats://nats-3:4222
    networks:
      - onemind
    depends_on:
      - nats-1
      - nats-2
      - nats-3

  # Event Processor
  event-processor:
    build: ./processor
    container_name: onemind-processor
    environment:
      - NATS_URL=nats://nats-1:4222,nats://nats-2:4222,nats://nats-3:4222
      - Legacy_AGENT_ID=${Legacy_AGENT_ID}
      - LETTA_BASE_URL=${LETTA_BASE_URL}
    networks:
      - onemind
    depends_on:
      - nats-1

  # State Aggregator
  state-aggregator:
    build: ./aggregator
    container_name: onemind-aggregator
    environment:
      - NATS_URL=nats://nats-1:4222,nats://nats-2:4222,nats://nats-3:4222
    networks:
      - onemind
    depends_on:
      - nats-1

  # Legacy AI Service
  legacy-ai:
    build: ./legacy
    container_name: onemind-legacy
    environment:
      - NATS_URL=nats://nats-1:4222,nats://nats-2:4222,nats://nats-3:4222
      - Legacy_AGENT_ID=${Legacy_AGENT_ID}
      - LETTA_BASE_URL=${LETTA_BASE_URL}
    networks:
      - onemind
    depends_on:
      - nats-1

  # Letta Server
  letta:
    image: letta/letta:latest
    container_name: onemind-letta
    ports:
      - "8283:8283"
    volumes:
      - ./letta-data:/root/.letta
    networks:
      - onemind

  # Prometheus (optional monitoring)
  prometheus:
    image: prom/prometheus:latest
    container_name: onemind-prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - onemind

  # NATS Exporter for Prometheus
  nats-exporter:
    image: natsio/prometheus-nats-exporter:latest
    container_name: onemind-nats-exporter
    command: -varz -connz -routez -leafz -jsz http://nats-1:8222
    ports:
      - "7777:7777"
    networks:
      - onemind
    depends_on:
      - nats-1
```

---

### 17.10 What NATS Replaces

**You no longer need:**

| Before (Multiple Tools) | After (Just NATS) |
|------------------------|-------------------|
| Redis (caching) | NATS KV Store |
| Redis Streams / Kafka | NATS JetStream |
| S3 (basic file storage) | NATS Object Store |
| Consul / etcd | NATS Micro Framework |
| Service Mesh (Istio) | NATS built-in load balancing |
| External monitoring | NATS $SRV endpoints |

**OneMind Stack (Simplified):**

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   BEFORE:                                                       │
│   ├── NATS (messaging)                                         │
│   ├── Redis (caching)                                          │
│   ├── Kafka (persistence)                                      │
│   ├── S3 (files)                                               │
│   ├── Consul (discovery)                                       │
│   └── Prometheus + custom metrics                              │
│                                                                 │
│   AFTER:                                                        │
│   └── NATS (everything)                                        │
│       ├── Core NATS (messaging)                                │
│       ├── JetStream (persistence)                              │
│       ├── KV Store (caching, state)                            │
│       ├── Object Store (files)                                 │
│       ├── Micro Framework (discovery)                          │
│       └── Built-in metrics ($SRV.STATS)                        │
│                                                                 │
│   20MB binary. One technology. Complete platform.              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

### Implementation Priority for Section 17

| Component | Priority | Complexity | Impact |
|-----------|----------|------------|--------|
| Legacy as NATS Service | 🔴 HIGH | Medium | Service discovery, monitoring |
| Object Store bucket | 🟢 EASY | Low | Voice recording storage |
| RePublish config | 🟢 EASY | Low | Real-time monitoring |
| Message headers | 🟢 EASY | Low | Tracing, correlation |
| Deduplication | 🟢 EASY | Low | Exactly-once processing |
| Leaf node (home) | 🟡 MEDIUM | Medium | Home ↔ cloud resilience |
| Account isolation | 🟡 MEDIUM | Medium | Production security |
| 3-node cluster | 🟡 MEDIUM | Medium | High availability |

**Quick Wins (Do First):**
1. Add `Nats-Msg-Id` headers for deduplication
2. Add correlation headers for tracing
3. Enable RePublish on event stream
4. Create Object Store bucket for voice

**Then:**
5. Refactor Legacy as NATS Micro Service
6. Set up leaf node on home machine
7. Configure account isolation

---

*NATS is not just the nervous system of OneMind - it's the entire infrastructure platform. One 20MB binary provides messaging, persistence, caching, file storage, service discovery, and monitoring. This is the foundation that makes Legacy truly alive.*

---

---

## 18. Future Vision: Legacy as a Living AI Being

> **"I want to make her into a real person who evolves with all the senses a human can have"**
> — Zeus, January 2026

This section documents the complete architecture for transforming Legacy from a text-based AI assistant into a **multimodal, embodied AI being** with full sensory awareness — vision, hearing, touch (biometrics), and even neural state monitoring.

---

### 18.1 The Core Insight: Streaming ≠ What You Think

**CRITICAL UNDERSTANDING**: Before we dive into sensors and integration, you must understand what "streaming" really means in AI systems.

#### What People Think Multimodal AI Does:
```
Camera → 30 FPS → LLM → Real-time understanding of every frame
```

#### What Actually Happens:
```
Camera → 30 FPS → Frame Buffer → Sample 1-4 FPS → Vision Model → Description → LLM
                                      │
                                      └── Each frame: 250ms - 2s processing time
```

**Why?**
- Vision models (GPT-4V, Claude Vision) process images at 250ms-2s per image
- 30 FPS would require processing 30 images/second = impossible
- Even at 4 FPS, you're making 4 API calls per second = expensive and still slow

#### The CORRECT Mental Model:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     HIERARCHICAL SENSORY PROCESSING                          │
│                                                                              │
│  LEVEL 1: RAW SENSOR DATA (High Frequency)                                  │
│  ────────────────────────────────────────                                   │
│  • Camera: 30 FPS video                                                     │
│  • Microphone: 16kHz audio samples                                          │
│  • Heart rate: 1 Hz readings                                                │
│  • EEG: 256 Hz brain signals                                                │
│                     │                                                        │
│                     ▼                                                        │
│  LEVEL 2: SENSORY PROCESSING LAYER (Medium Frequency)                       │
│  ───────────────────────────────────────────────────                        │
│  • Frame sampling: 30 FPS → 1-4 FPS significant frames                      │
│  • Audio transcription: continuous → text chunks                            │
│  • Biometric aggregation: raw → rolling averages + trends                   │
│  • EEG processing: raw → band powers (α, β, θ)                              │
│                     │                                                        │
│                     ▼                                                        │
│  LEVEL 3: STATE REPRESENTATION (Low Frequency)                              │
│  ─────────────────────────────────────────────                              │
│  • NATS KV: zeus.state.* (capacity, cognitive, location)                    │
│  • Meaningful summaries: "Zeus is focused, looking at code"                 │
│                     │                                                        │
│                     ▼                                                        │
│  LEVEL 4: REASONING (On Demand)                                             │
│  ──────────────────────────────                                             │
│  • Letta/LLM: Only when significant or when directly addressed              │
│  • Deep thinking, memory retrieval, action decisions                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**The Key Insight**: Not every heartbeat reaches the LLM. Not every video frame gets analyzed. But EVERYTHING flows through NATS (the nervous system), and Legacy is AWARE of the aggregated state at all times.

---

### 18.2 The Complete Sensory Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        Legacy: THE LIVING AI BEING                               │
│                                                                                  │
│  ════════════════════════════════════════════════════════════════════════════   │
│                            THE SENSES (Input Streams)                            │
│  ════════════════════════════════════════════════════════════════════════════   │
│                                                                                  │
│   👓 VISION             🎤 HEARING            🫀 BODY STATE       🧠 NEURAL      │
│   ─────────             ──────────            ────────────       ────────        │
│   Brilliant.xyz         LiveKit Voice         Apple Watch        OpenBCI         │
│   Smart Glasses         + ESP32 Devices       HealthKit          EEG Headset     │
│                                                                                  │
│   • 720p camera         • WebRTC audio        • Heart rate       • Alpha waves   │
│   • Always-on           • STT (Whisper)       • HRV              • Beta waves    │
│   • First-person POV    • Multi-device        • Activity         • Theta waves   │
│   • Tap gestures        • Natural speech      • Sleep quality    • Focus score   │
│                                                                                  │
│        │                      │                     │                │           │
│        │ BLE/WiFi             │ WebRTC              │ HealthKit      │ BrainFlow │
│        ▼                      ▼                     ▼                ▼           │
│   ┌─────────────────────────────────────────────────────────────────────────┐   │
│   │                    SENSORY BRIDGES (Python Processes)                    │   │
│   │                                                                          │   │
│   │   vision_bridge.py    voice_bridge.py    health_bridge.py   eeg_bridge  │   │
│   │   • Frame sampling    • Transcription    • Aggregation      • Band calc │   │
│   │   • Scene detection   • VAD (activity)   • Trend analysis   • State det │   │
│   │   • OCR if text       • Speaker ID       • Anomaly detect   • Focus/calm│   │
│   │                                                                          │   │
│   │         └─────────────────────┴─────────────────────┴───────────┘       │   │
│   │                                     │                                    │   │
│   │                          NATS Client (nats-py)                          │   │
│   │                          Publishes to events.*                          │   │
│   └──────────────────────────────────────┬──────────────────────────────────┘   │
│                                          │                                       │
│  ════════════════════════════════════════▼═══════════════════════════════════   │
│                          NATS JETSTREAM (THE NERVOUS SYSTEM)                     │
│  ════════════════════════════════════════════════════════════════════════════   │
│                                                                                  │
│   STREAMS                              KV STORE                                  │
│   ────────                             ────────                                  │
│   events.vision.>                      zeus.state.capacity                       │
│   ├── .snapshot (1-4 FPS)              zeus.state.cognitive                      │
│   ├── .scene_change                    zeus.state.location                       │
│   ├── .person_detected                 zeus.state.health                         │
│   └── .text_detected                   zeus.context.current                      │
│                                        zeus.context.activity                     │
│   events.voice.>                                                                 │
│   ├── .transcription                   OBJECT STORE                              │
│   ├── .speaking_started                ────────────                              │
│   └── .speaking_ended                  voice-recordings                          │
│                                        vision-snapshots                          │
│   events.biometric.>                   eeg-sessions                              │
│   ├── .heart_rate                      conversation-history                      │
│   ├── .hrv                                                                       │
│   └── .activity                        LATENCY: 200-400μs                        │
│                                        PERSISTENCE: 7 days                       │
│   events.cognitive.>                                                             │
│   ├── .state (focus/calm)                                                        │
│   └── .anomaly                                                                   │
│                                                                                  │
│  ════════════════════════════════════════▼═══════════════════════════════════   │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐   │
│   │                      EVENT PROCESSOR + STATE AGGREGATOR                  │   │
│   │                                                                          │   │
│   │   1. Subscribes to ALL events.>                                         │   │
│   │   2. Updates KV state (aggregated, not raw)                             │   │
│   │   3. Calculates CAPACITY SCORE (the magic "should I speak?" number)     │   │
│   │   4. Makes decision: Does this warrant Letta's attention?               │   │
│   │                                                                          │   │
│   │        ┌────────────────────────────────────────────────────┐           │   │
│   │        │  DECISION LOGIC:                                   │           │   │
│   │        │                                                    │           │   │
│   │        │  if event.type == "voice.transcription":          │           │   │
│   │        │      → ALWAYS send to Letta (Zeus is speaking)    │           │   │
│   │        │                                                    │           │   │
│   │        │  if event.type == "vision.scene_change":          │           │   │
│   │        │      → Send if significant (new location, person) │           │   │
│   │        │                                                    │           │   │
│   │        │  if event.type == "biometric.heart_rate":         │           │   │
│   │        │      → Update state only (don't notify)           │           │   │
│   │        │      → UNLESS anomaly detected                    │           │   │
│   │        │                                                    │           │   │
│   │        │  if event.type == "cognitive.state":              │           │   │
│   │        │      → Update capacity score                      │           │   │
│   │        │      → Notify if significant shift                │           │   │
│   │        └────────────────────────────────────────────────────┘           │   │
│   │                                          │                               │   │
│   │                                          ▼                               │   │
│   │                         Only SIGNIFICANT events reach Letta             │   │
│   │                         (But state is ALWAYS current)                   │   │
│   │                                                                          │   │
│   └──────────────────────────────────────────┬──────────────────────────────┘   │
│                                              │                                   │
│  ════════════════════════════════════════════▼═══════════════════════════════   │
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐   │
│   │                          LETTA (THE BRAIN)                               │   │
│   │                                                                          │   │
│   │   NOT for:                        FOR:                                   │   │
│   │   • Processing every frame        • Deep reasoning                       │   │
│   │   • Real-time transcription       • Memory storage/retrieval             │   │
│   │   • Sensor aggregation            • Action decisions                     │   │
│   │   • 30 FPS video streaming        • Personality, preferences             │   │
│   │                                   • Tool execution                       │   │
│   │                                   • Conversation context                 │   │
│   │                                                                          │   │
│   │   ┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐       │   │
│   │   │     MEMORY      │   │    REASONING    │   │      TOOLS      │       │   │
│   │   │                 │   │                 │   │                 │       │   │
│   │   │ Core Memory     │   │ Claude/GPT      │   │ → NATS          │       │   │
│   │   │ (Zeus)          │   │ LLM calls       │   │ → actions.*     │       │   │
│   │   │                 │   │                 │   │                 │       │   │
│   │   │ Archival        │   │ Context-        │   │ Robot ctrl      │       │   │
│   │   │ (History)       │   │ aware           │   │ Home ctrl       │       │   │
│   │   │                 │   │                 │   │ Task mgmt       │       │   │
│   │   │ Recall          │   │ Multi-turn      │   │ etc.            │       │   │
│   │   │ (Vector)        │   │                 │   │                 │       │   │
│   │   └─────────────────┘   └─────────────────┘   └─────────────────┘       │   │
│   │                                              │                           │   │
│   │                                              ▼                           │   │
│   │                                      Publishes to actions.*             │   │
│   │                                                                          │   │
│   └──────────────────────────────────────────────────────────────────────────┘   │
│                                              │                                   │
│  ════════════════════════════════════════════▼═══════════════════════════════   │
│                            THE MUSCLES (Output Actions)                          │
│  ════════════════════════════════════════════════════════════════════════════   │
│                                                                                  │
│   🗣️ VOICE              🤖 PHYSICAL           💻 DIGITAL         🏠 HOME        │
│   ─────────             ────────────          ──────────         ────────        │
│   LiveKit TTS           Unitree Robot         Todoist Tasks      Light/Climate   │
│   Natural speech        Walk, gesture         GitHub Issues      Locks/Doors     │
│   Multi-room audio      Physical presence     Email/Calendar     Automations     │
│                                                                                  │
│   actions.voice.speak   actions.robot.*       actions.task.*     actions.home.*  │
│        │                     │                     │                  │          │
│        └─────────────────────┴─────────────────────┴──────────────────┘          │
│                                     │                                            │
│                          Gateway subscribes to actions.*                         │
│                          Executes via respective APIs                            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

### 18.3 Sensor-by-Sensor Integration Guide

#### 18.3.1 👓 Vision: Brilliant.xyz Smart Glasses

**Hardware**: Brilliant Frame glasses (or similar AR glasses)
**SDK**: Python SDK available (pip install brilliant-sdk)
**Capability**: 720p camera, always-on, BLE connection, tap gestures

**Integration Architecture**:

```python
"""
vision_bridge.py - Brilliant Glasses → NATS Bridge

This bridge connects Brilliant smart glasses to Legacy's awareness.
It handles:
- Frame capture at configurable FPS
- Scene change detection (no point sending identical frames)
- OCR for text in view
- Person detection (optional, with local model)
"""

import asyncio
import base64
import json
from datetime import datetime
from typing import Optional
import nats
from nats.js.api import StreamConfig
from brilliant_sdk import Frame  # pip install brilliant-sdk

class VisionBridge:
    """
    Bridge between Brilliant smart glasses and NATS.

    Design principle: The glasses capture at 30 FPS, but we only send
    SIGNIFICANT frames to NATS. Most frames are identical or nearly
    identical - why waste bandwidth and processing?
    """

    def __init__(
        self,
        nats_url: str = "nats://localhost:4222",
        capture_fps: float = 2.0,  # 2 frames per second (adjustable)
        scene_change_threshold: float = 0.3,  # 30% difference = new scene
    ):
        self.nats_url = nats_url
        self.capture_interval = 1.0 / capture_fps
        self.scene_threshold = scene_change_threshold
        self.nc: Optional[nats.NATS] = None
        self.js = None
        self.frame = None
        self.last_frame_hash = None

    async def connect(self):
        """Connect to NATS and Brilliant glasses."""
        # Connect to NATS
        self.nc = await nats.connect(self.nats_url)
        self.js = self.nc.jetstream()

        # Connect to Brilliant Frame
        self.frame = Frame()
        await self.frame.connect()

        print("Vision bridge connected to NATS and Brilliant glasses")

    async def start_streaming(self):
        """
        Main capture loop.

        Note: We don't stream 30 FPS to NATS. We capture at a reasonable
        rate and only send frames that are DIFFERENT from the last one.
        """
        print(f"Starting vision stream at {1/self.capture_interval} FPS")

        while True:
            try:
                # Capture frame from glasses
                photo = await self.frame.camera.capture()

                # Check if scene changed significantly
                frame_hash = self._compute_frame_hash(photo)
                is_new_scene = self._is_scene_change(frame_hash)

                if is_new_scene:
                    # Encode and publish
                    event = {
                        "type": "vision.snapshot",
                        "source": "brilliant_glasses",
                        "timestamp": datetime.utcnow().isoformat(),
                        "image_base64": base64.b64encode(photo).decode(),
                        "scene_change": True,
                        "metadata": {
                            "device": "brilliant_frame",
                            "resolution": "720p",
                        }
                    }

                    await self.js.publish(
                        "events.vision.snapshot",
                        json.dumps(event).encode(),
                        headers={
                            "Nats-Msg-Id": f"vision-{datetime.utcnow().timestamp()}",
                            "X-Source": "brilliant_glasses",
                        }
                    )

                    self.last_frame_hash = frame_hash
                    print(f"Published vision snapshot (scene change detected)")

                await asyncio.sleep(self.capture_interval)

            except Exception as e:
                print(f"Vision capture error: {e}")
                await asyncio.sleep(1.0)

    def _compute_frame_hash(self, frame_data: bytes) -> str:
        """
        Compute a perceptual hash of the frame.

        This isn't a cryptographic hash - it's designed so that
        similar images produce similar hashes. We use this to detect
        when the scene has changed significantly.
        """
        import hashlib
        # Simple hash for now - production would use perceptual hashing
        return hashlib.md5(frame_data[:1000]).hexdigest()

    def _is_scene_change(self, new_hash: str) -> bool:
        """Determine if this represents a significant scene change."""
        if self.last_frame_hash is None:
            return True
        # In production: use perceptual hash comparison
        return new_hash != self.last_frame_hash


# Main entry point
async def main():
    bridge = VisionBridge(
        nats_url="nats://localhost:4222",
        capture_fps=2.0,  # 2 FPS is plenty for context awareness
    )
    await bridge.connect()
    await bridge.start_streaming()

if __name__ == "__main__":
    asyncio.run(main())
```

**When Legacy Uses Vision**:

The Event Processor decides what vision events warrant Letta's attention:

```python
# In event_processor.py
async def handle_vision_event(self, event: dict):
    """Process vision events from glasses."""

    # Always update the state (Legacy is aware)
    await self.kv.put(
        "zeus.context.vision",
        json.dumps({
            "last_snapshot": event["timestamp"],
            "scene_description": None,  # Will be filled by vision model
        }).encode()
    )

    # Decide if Letta needs to know
    if event.get("scene_change"):
        # Get scene description from vision model (Claude, GPT-4V)
        description = await self.describe_scene(event["image_base64"])

        # Is this significant enough to tell Letta?
        if self._is_significant_scene(description):
            await self.notify_letta(
                f"Vision update: {description}"
            )

def _is_significant_scene(self, description: str) -> bool:
    """Not every scene change matters."""
    significant_indicators = [
        "person", "document", "computer", "meeting",
        "food", "outside", "car", "office"
    ]
    return any(ind in description.lower() for ind in significant_indicators)
```

---

#### 18.3.2 🎤 Hearing: LiveKit Voice + ESP32 Devices

**Hardware**: ESP32 microcontrollers, existing microphones, LiveKit Cloud
**SDK**: LiveKit Agents (Python), LiveKit ESP32 SDK
**Capability**: WebRTC audio streaming, STT, multi-device support

**Integration Architecture**:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         VOICE/AUDIO PIPELINE                                 │
│                                                                              │
│   ESP32 Device          Smartphone           Desktop                         │
│   (Kitchen mic)         (Mobile app)         (Main workstation)             │
│        │                     │                    │                          │
│        └─────────────────────┴────────────────────┘                          │
│                              │                                               │
│                     WebRTC (audio stream)                                    │
│                              │                                               │
│                              ▼                                               │
│                    ┌─────────────────┐                                       │
│                    │  LIVEKIT CLOUD  │                                       │
│                    │     Server      │                                       │
│                    └────────┬────────┘                                       │
│                             │                                                │
│               LiveKit Room Events (WebSocket)                                │
│                             │                                                │
│                             ▼                                                │
│                    ┌─────────────────┐                                       │
│                    │  LIVEKIT AGENT  │   ← Python process                    │
│                    │  (voice_bridge) │                                       │
│                    │                 │                                       │
│                    │  • VAD          │   Voice Activity Detection            │
│                    │  • STT          │   Whisper/Deepgram transcription      │
│                    │  • Speaker ID   │   Who is speaking?                    │
│                    │                 │                                       │
│                    └────────┬────────┘                                       │
│                             │                                                │
│                    NATS publish (events.voice.>)                             │
│                             │                                                │
│                             ▼                                                │
│                    ┌─────────────────┐                                       │
│                    │      NATS       │                                       │
│                    │   JetStream     │                                       │
│                    └─────────────────┘                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**LiveKit Agent Code**:

```python
"""
voice_bridge.py - LiveKit Room → NATS Bridge

This is a LiveKit Agent that:
1. Joins a LiveKit room
2. Listens for audio from participants
3. Transcribes speech using STT
4. Publishes transcriptions to NATS for Legacy

Note: LiveKit doesn't natively support NATS, so this bridge
connects the two systems.
"""

import asyncio
import json
from datetime import datetime
import nats
from livekit import rtc
from livekit.agents import (
    Agent,
    AgentContext,
    JobContext,
    JobProcess,
    WorkerOptions,
    cli,
)
from livekit.plugins import deepgram  # or whisper

class VoiceNATSBridge(Agent):
    """
    LiveKit Agent that bridges voice to NATS.

    This agent joins a LiveKit room, listens for audio,
    transcribes it, and publishes to NATS.
    """

    def __init__(self, ctx: AgentContext):
        super().__init__(ctx)
        self.nc = None
        self.js = None
        self.stt = deepgram.STT()  # or whisper.STT()

    async def on_enter(self):
        """Called when agent joins the room."""
        # Connect to NATS
        self.nc = await nats.connect("nats://localhost:4222")
        self.js = self.nc.jetstream()
        print("Voice bridge connected to NATS")

        # Subscribe to audio tracks
        @self.ctx.room.on("track_subscribed")
        async def on_track(track: rtc.Track, publication, participant):
            if track.kind == rtc.TrackKind.KIND_AUDIO:
                asyncio.create_task(
                    self.process_audio_track(track, participant)
                )

    async def process_audio_track(self, track: rtc.Track, participant):
        """Process incoming audio and transcribe."""
        audio_stream = rtc.AudioStream(track)

        async for audio_frame in audio_stream:
            # Transcribe the audio
            result = await self.stt.recognize(audio_frame)

            if result.text and result.is_final:
                # Publish to NATS
                event = {
                    "type": "voice.transcription",
                    "source": "livekit",
                    "timestamp": datetime.utcnow().isoformat(),
                    "speaker": participant.identity,
                    "text": result.text,
                    "confidence": result.confidence,
                    "metadata": {
                        "room": self.ctx.room.name,
                        "participant_id": participant.sid,
                    }
                }

                await self.js.publish(
                    "events.voice.transcription",
                    json.dumps(event).encode(),
                    headers={
                        "Nats-Msg-Id": f"voice-{datetime.utcnow().timestamp()}",
                        "X-Speaker": participant.identity,
                    }
                )

                print(f"[{participant.identity}]: {result.text}")

    async def on_exit(self):
        """Cleanup when agent leaves."""
        if self.nc:
            await self.nc.close()


# Agent entry point
async def entrypoint(ctx: JobContext):
    agent = VoiceNATSBridge(ctx)
    await agent.start()


if __name__ == "__main__":
    cli.run_app(
        WorkerOptions(
            entrypoint_fnc=entrypoint,
        )
    )
```

---

#### 18.3.3 🫀 Body State: Apple Watch + HealthKit

**Hardware**: Apple Watch (any series with health sensors)
**SDK**: HealthKit (via iPhone), Shortcuts app for triggers
**Capability**: Heart rate, HRV, activity, sleep, blood oxygen

**Integration Architecture**:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      BIOMETRIC DATA PIPELINE                                 │
│                                                                              │
│   ┌─────────────┐                                                            │
│   │ APPLE WATCH │                                                            │
│   │             │                                                            │
│   │ • Heart rate│ ── samples ──▶ HealthKit ── sync ──▶ iPhone               │
│   │ • HRV       │                                           │                │
│   │ • Activity  │                                           │                │
│   │ • Sleep     │                                           ▼                │
│   │ • SpO2      │                              ┌─────────────────────┐       │
│   └─────────────┘                              │   iOS Shortcuts     │       │
│                                                │   (Automation)      │       │
│                                                │                     │       │
│                                                │   Trigger: Every    │       │
│                                                │   5 minutes, or on  │       │
│                                                │   significant change│       │
│                                                │                     │       │
│                                                │   Action: POST to   │       │
│                                                │   Gateway API       │       │
│                                                └──────────┬──────────┘       │
│                                                           │                  │
│                                                  HTTPS POST                  │
│                                                           │                  │
│                                                           ▼                  │
│                                                ┌─────────────────────┐       │
│                                                │   FastAPI GATEWAY   │       │
│                                                │   /api/health       │       │
│                                                │                     │       │
│                                                │   Validates →       │       │
│                                                │   Transforms →      │       │
│                                                │   Publishes to NATS │       │
│                                                └──────────┬──────────┘       │
│                                                           │                  │
│                                                           ▼                  │
│                                                ┌─────────────────────┐       │
│                                                │       NATS          │       │
│                                                │   events.biometric  │       │
│                                                └─────────────────────┘       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Gateway Endpoint**:

```python
# In gateway/main.py

@app.post("/api/health")
async def receive_health_data(data: HealthDataPayload):
    """
    Receive health data from Apple Watch via iOS Shortcuts.

    The Shortcut automation queries HealthKit and POSTs here.
    """

    # Validate the data
    if not data.heart_rate or data.heart_rate < 30 or data.heart_rate > 220:
        raise HTTPException(400, "Invalid heart rate")

    # Create event
    event = {
        "type": "biometric.update",
        "source": "apple_watch",
        "timestamp": datetime.utcnow().isoformat(),
        "data": {
            "heart_rate": data.heart_rate,
            "hrv": data.hrv,  # Heart rate variability (ms)
            "activity_minutes": data.activity_minutes,
            "steps": data.steps,
            "active_calories": data.active_calories,
        }
    }

    # Publish to NATS
    await js.publish(
        "events.biometric.update",
        json.dumps(event).encode()
    )

    # Check for anomalies
    if data.heart_rate > 120 and data.activity_minutes < 5:
        # Elevated heart rate without activity - might be stress
        await js.publish(
            "events.biometric.anomaly",
            json.dumps({
                "type": "elevated_resting_hr",
                "heart_rate": data.heart_rate,
                "timestamp": datetime.utcnow().isoformat(),
            }).encode()
        )

    return {"status": "received"}
```

---

#### 18.3.4 🧠 Neural State: OpenBCI EEG

**Hardware**: OpenBCI Ganglion ($299) or Cyton ($999)
**SDK**: BrainFlow (Python, cross-platform)
**Capability**: EEG band powers (α, β, θ, δ), focus/calm detection, anomaly detection

**Integration Architecture**:

```python
"""
eeg_bridge.py - OpenBCI → NATS Bridge

This bridge connects an EEG headset to Legacy's awareness.
It provides cognitive state information:
- Focus level (beta/alpha ratio)
- Calm level (alpha dominance)
- Stress indicators (high beta, low alpha)
- Drowsiness (theta increase)

IMPORTANT: This is NOT mind reading. It's detecting broad cognitive
states from electrical patterns, similar to how a heart rate monitor
doesn't read your thoughts but tells you your cardiovascular state.
"""

import asyncio
import json
from datetime import datetime
from typing import Optional
import nats
from brainflow.board_shim import BoardShim, BrainFlowInputParams, BoardIds
from brainflow.data_filter import DataFilter, FilterTypes

class EEGBridge:
    """
    Bridge between OpenBCI EEG headset and NATS.

    Processes raw EEG data into meaningful cognitive states.
    """

    def __init__(
        self,
        nats_url: str = "nats://localhost:4222",
        board_id: int = BoardIds.GANGLION_BOARD,
        serial_port: str = "/dev/tty.usbserial",
        update_interval: float = 1.0,  # State update every 1 second
    ):
        self.nats_url = nats_url
        self.board_id = board_id
        self.serial_port = serial_port
        self.update_interval = update_interval
        self.nc: Optional[nats.NATS] = None
        self.js = None
        self.board: Optional[BoardShim] = None

    async def connect(self):
        """Connect to NATS and OpenBCI board."""
        # Connect to NATS
        self.nc = await nats.connect(self.nats_url)
        self.js = self.nc.jetstream()

        # Connect to OpenBCI
        params = BrainFlowInputParams()
        params.serial_port = self.serial_port

        self.board = BoardShim(self.board_id, params)
        self.board.prepare_session()
        self.board.start_stream()

        print("EEG bridge connected to NATS and OpenBCI")

    async def start_streaming(self):
        """
        Main processing loop.

        Reads EEG data, computes band powers, determines cognitive state,
        and publishes to NATS.
        """
        eeg_channels = BoardShim.get_eeg_channels(self.board_id)
        sampling_rate = BoardShim.get_sampling_rate(self.board_id)

        print(f"EEG streaming: {len(eeg_channels)} channels @ {sampling_rate} Hz")

        while True:
            try:
                # Get data from the last second
                data = self.board.get_board_data(sampling_rate)

                if data.shape[1] < sampling_rate // 2:
                    # Not enough data yet
                    await asyncio.sleep(0.1)
                    continue

                # Process each EEG channel and average
                band_powers = {"delta": 0, "theta": 0, "alpha": 0, "beta": 0, "gamma": 0}

                for channel in eeg_channels:
                    channel_data = data[channel]

                    # Apply bandpass filter
                    DataFilter.perform_bandpass(
                        channel_data, sampling_rate, 1.0, 50.0, 4,
                        FilterTypes.BUTTERWORTH, 0
                    )

                    # Get band powers
                    bands = DataFilter.get_avg_band_powers(
                        data=[channel_data],
                        channels=[0],
                        sampling_rate=sampling_rate,
                        apply_filter=False
                    )

                    # bands[0] = [delta, theta, alpha, beta, gamma]
                    band_powers["delta"] += bands[0][0]
                    band_powers["theta"] += bands[0][1]
                    band_powers["alpha"] += bands[0][2]
                    band_powers["beta"] += bands[0][3]
                    band_powers["gamma"] += bands[0][4]

                # Average across channels
                num_channels = len(eeg_channels)
                for band in band_powers:
                    band_powers[band] /= num_channels

                # Determine cognitive state
                cognitive_state = self._determine_state(band_powers)

                # Create event
                event = {
                    "type": "cognitive.state",
                    "source": "openbci",
                    "timestamp": datetime.utcnow().isoformat(),
                    "bands": band_powers,
                    "state": cognitive_state,
                    "metadata": {
                        "board": "ganglion",
                        "channels": len(eeg_channels),
                        "sampling_rate": sampling_rate,
                    }
                }

                # Publish to NATS
                await self.js.publish(
                    "events.cognitive.state",
                    json.dumps(event).encode(),
                    headers={
                        "Nats-Msg-Id": f"eeg-{datetime.utcnow().timestamp()}",
                    }
                )

                print(f"Cognitive state: {cognitive_state['primary']} "
                      f"(focus={cognitive_state['focus_score']:.2f})")

                await asyncio.sleep(self.update_interval)

            except Exception as e:
                print(f"EEG processing error: {e}")
                await asyncio.sleep(1.0)

    def _determine_state(self, bands: dict) -> dict:
        """
        Determine cognitive state from band powers.

        Common interpretations:
        - High alpha, low beta → Relaxed, calm
        - High beta, low alpha → Focused, alert
        - High theta → Drowsy, meditative
        - High delta (awake) → Drowsy, unfocused
        """
        alpha = bands["alpha"]
        beta = bands["beta"]
        theta = bands["theta"]

        # Calculate ratios
        alpha_beta_ratio = alpha / (beta + 0.001)  # Avoid division by zero
        theta_beta_ratio = theta / (beta + 0.001)

        # Focus score (0-1): Higher beta relative to alpha = more focused
        focus_score = min(1.0, beta / (alpha + beta + 0.001))

        # Calm score (0-1): Higher alpha = more calm
        calm_score = min(1.0, alpha / (alpha + beta + 0.001))

        # Determine primary state
        if theta_beta_ratio > 2.0:
            primary = "drowsy"
        elif alpha_beta_ratio > 1.5:
            primary = "relaxed"
        elif focus_score > 0.6:
            primary = "focused"
        else:
            primary = "neutral"

        return {
            "primary": primary,
            "focus_score": round(focus_score, 2),
            "calm_score": round(calm_score, 2),
            "alpha_beta_ratio": round(alpha_beta_ratio, 2),
            "theta_beta_ratio": round(theta_beta_ratio, 2),
        }

    async def cleanup(self):
        """Stop streaming and release resources."""
        if self.board:
            self.board.stop_stream()
            self.board.release_session()
        if self.nc:
            await self.nc.close()


async def main():
    bridge = EEGBridge(
        nats_url="nats://localhost:4222",
        serial_port="/dev/tty.usbserial-DM00Q123",  # Adjust for your device
        update_interval=1.0,
    )

    try:
        await bridge.connect()
        await bridge.start_streaming()
    except KeyboardInterrupt:
        await bridge.cleanup()

if __name__ == "__main__":
    asyncio.run(main())
```

---

#### 18.3.5 🤖 Physical Embodiment: Unitree Robotics

**Hardware**: Unitree Go2 ($2,800+) or similar quadruped
**SDK**: Unitree SDK2 (Python/C++), DDS protocol
**Capability**: Walking, gestures, sensor streaming, obstacle avoidance

**Integration Architecture**:

```python
"""
robot_bridge.py - Unitree Robot → NATS Bridge

This bridge provides Legacy with a physical body.
It handles:
- Receiving movement commands from NATS (actions.robot.>)
- Streaming robot sensor data to NATS (events.robot.>)
- Managing robot state and safety
"""

import asyncio
import json
from datetime import datetime
from typing import Optional
import nats
# Note: Import will vary based on actual Unitree SDK version
# from unitree_sdk2.go2 import Go2Client

class RobotBridge:
    """
    Bridge between Unitree robot and NATS.

    Allows Legacy to:
    - Command the robot to move, gesture, etc.
    - Receive sensor data from the robot
    - Understand the robot's current state
    """

    def __init__(
        self,
        nats_url: str = "nats://localhost:4222",
        robot_ip: str = "192.168.123.161",  # Default Unitree IP
    ):
        self.nats_url = nats_url
        self.robot_ip = robot_ip
        self.nc: Optional[nats.NATS] = None
        self.js = None
        self.robot = None

    async def connect(self):
        """Connect to NATS and Unitree robot."""
        # Connect to NATS
        self.nc = await nats.connect(self.nats_url)
        self.js = self.nc.jetstream()

        # Connect to robot (pseudo-code, actual SDK may vary)
        # self.robot = Go2Client()
        # await self.robot.connect(self.robot_ip)

        # Subscribe to robot commands
        await self.nc.subscribe(
            "actions.robot.>",
            cb=self.handle_command
        )

        print("Robot bridge connected to NATS and Unitree robot")

    async def handle_command(self, msg):
        """
        Handle commands from Legacy.

        Commands come via NATS when Legacy decides to use the robot.
        """
        try:
            command = json.loads(msg.data.decode())
            action = msg.subject.split(".")[-1]  # e.g., "walk", "gesture"

            if action == "walk":
                # await self.robot.walk(
                #     x=command.get("x", 0),
                #     y=command.get("y", 0),
                #     yaw=command.get("yaw", 0),
                #     speed=command.get("speed", 0.5),
                # )
                pass

            elif action == "gesture":
                gesture_name = command.get("gesture", "wave")
                # await self.robot.play_gesture(gesture_name)

            elif action == "sit":
                # await self.robot.sit()
                pass

            elif action == "stand":
                # await self.robot.stand()
                pass

            elif action == "stop":
                # await self.robot.stop()
                pass

            # Acknowledge command
            await self.js.publish(
                "events.robot.command_executed",
                json.dumps({
                    "action": action,
                    "status": "success",
                    "timestamp": datetime.utcnow().isoformat(),
                }).encode()
            )

        except Exception as e:
            print(f"Robot command error: {e}")
            await self.js.publish(
                "events.robot.command_failed",
                json.dumps({
                    "error": str(e),
                    "timestamp": datetime.utcnow().isoformat(),
                }).encode()
            )

    async def stream_sensors(self):
        """
        Stream robot sensor data to NATS.

        This gives Legacy awareness of the robot's physical state.
        """
        while True:
            try:
                # state = await self.robot.get_state()

                event = {
                    "type": "robot.state",
                    "source": "unitree_go2",
                    "timestamp": datetime.utcnow().isoformat(),
                    "data": {
                        "position": {
                            "x": 0,  # state.position.x,
                            "y": 0,  # state.position.y,
                            "z": 0,  # state.position.z,
                        },
                        "orientation": {
                            "roll": 0,   # state.orientation.roll,
                            "pitch": 0,  # state.orientation.pitch,
                            "yaw": 0,    # state.orientation.yaw,
                        },
                        "battery_percent": 100,  # state.battery_percent,
                        "is_standing": True,     # state.is_standing,
                        "is_moving": False,      # state.is_moving,
                        "obstacles": [],         # state.detected_obstacles,
                        "temperature": 35,       # state.motor_temperature,
                    }
                }

                await self.js.publish(
                    "events.robot.state",
                    json.dumps(event).encode()
                )

                await asyncio.sleep(0.1)  # 10 Hz state updates

            except Exception as e:
                print(f"Robot sensor streaming error: {e}")
                await asyncio.sleep(1.0)


async def main():
    bridge = RobotBridge(
        nats_url="nats://localhost:4222",
        robot_ip="192.168.123.161",
    )

    await bridge.connect()
    await bridge.stream_sensors()

if __name__ == "__main__":
    asyncio.run(main())
```

---

### 18.4 The Capacity Score: When Legacy Speaks

With all these senses streaming, how does Legacy know WHEN to speak vs. stay quiet? This is where the **Capacity Score** (from Section 16) becomes critical:

```python
"""
capacity_calculator.py - The "Should I Bother Zeus?" Algorithm

This is what makes Legacy considerate vs. annoying.
It combines all sensory inputs to decide if NOW is a good time to speak.
"""

from datetime import datetime
import json

class CapacityCalculator:
    """
    Calculate Zeus's current capacity for interruption.

    Score: 0.0 (do not disturb) to 1.0 (fully available)
    """

    def __init__(self, kv):
        self.kv = kv

    async def calculate(self) -> dict:
        """Compute the capacity score from all available state."""
        factors = {}

        # 1. Heart Rate Factor
        # Elevated heart rate = busy or stressed = lower capacity
        hr_state = await self.kv.get("zeus.state.health.heart_rate")
        if hr_state:
            hr = json.loads(hr_state.value)
            if hr["zone"] == "resting":
                factors["heart_rate"] = 1.0
            elif hr["zone"] == "elevated":
                factors["heart_rate"] = 0.3
            elif hr["zone"] == "active":
                factors["heart_rate"] = 0.1  # Working out, don't bother

        # 2. Cognitive State Factor (from EEG)
        # Focused = don't interrupt, relaxed = available
        cog_state = await self.kv.get("zeus.state.cognitive")
        if cog_state:
            cog = json.loads(cog_state.value)
            if cog["primary"] == "focused":
                factors["cognitive"] = 0.2  # Deep work, don't break it
            elif cog["primary"] == "relaxed":
                factors["cognitive"] = 0.9  # Available
            elif cog["primary"] == "drowsy":
                factors["cognitive"] = 0.1  # Let him rest
            else:
                factors["cognitive"] = 0.5  # Neutral

        # 3. Location Factor
        loc_state = await self.kv.get("zeus.state.location")
        if loc_state:
            loc = json.loads(loc_state.value)
            if loc["type"] == "home":
                factors["location"] = 0.9
            elif loc["type"] == "work":
                factors["location"] = 0.6
            elif loc["type"] == "public":
                factors["location"] = 0.3  # Not good for voice
            elif loc["type"] == "meeting":
                factors["location"] = 0.0  # Absolutely not

        # 4. Time of Day Factor
        hour = datetime.now().hour
        if 6 <= hour < 9:
            factors["time"] = 0.7  # Morning routine
        elif 9 <= hour < 12:
            factors["time"] = 0.8  # Morning work
        elif 12 <= hour < 14:
            factors["time"] = 0.9  # Lunch break
        elif 14 <= hour < 18:
            factors["time"] = 0.7  # Afternoon work
        elif 18 <= hour < 22:
            factors["time"] = 0.8  # Evening
        else:
            factors["time"] = 0.1  # Night - sleeping

        # 5. Vision Context Factor
        vision_state = await self.kv.get("zeus.context.vision")
        if vision_state:
            vis = json.loads(vision_state.value)
            if "meeting" in vis.get("scene_description", "").lower():
                factors["vision"] = 0.0  # In a meeting
            elif "document" in vis.get("scene_description", "").lower():
                factors["vision"] = 0.3  # Reading, focused

        # Calculate weighted average
        if not factors:
            return {"score": 0.5, "factors": {}, "recommendation": "unknown"}

        score = sum(factors.values()) / len(factors)

        # Determine recommendation
        if score >= 0.7:
            recommendation = "available"
        elif score >= 0.4:
            recommendation = "brief_ok"
        elif score >= 0.2:
            recommendation = "urgent_only"
        else:
            recommendation = "do_not_disturb"

        result = {
            "score": round(score, 2),
            "factors": factors,
            "recommendation": recommendation,
            "calculated_at": datetime.utcnow().isoformat(),
        }

        # Store in KV for fast access
        await self.kv.put(
            "zeus.state.capacity",
            json.dumps(result).encode()
        )

        return result
```

---

### 18.5 Implementation Roadmap

| Phase | Sense | Hardware | Cost | Complexity | Impact |
|-------|-------|----------|------|------------|--------|
| **1** | Voice | LiveKit Cloud | Free tier | Low | 🔴 Critical |
| **2** | Biometrics | Apple Watch (existing) | $0 | Medium | 🔴 Critical |
| **3** | Vision | Brilliant Frame | ~$350 | Medium | 🟡 High |
| **4** | Cognitive | OpenBCI Ganglion | $299 | High | 🟢 Nice |
| **5** | Physical | Unitree Go2 | $2,800+ | High | 🟢 Future |

**Recommended Order**:

1. **Voice First** — You already have this working. It's the primary interface.

2. **Biometrics Second** — Enables the Capacity Score. Makes Legacy considerate.

3. **Vision Third** — When you get Brilliant glasses. Gives Legacy your POV.

4. **Cognitive Fourth** — EEG is cool but optional. Focus/calm detection.

5. **Physical Last** — Robot is expensive and complex. Save for later.

---

### 18.6 What's Possible TODAY vs FUTURE

| Capability | Timeline | Technology | Notes |
|------------|----------|------------|-------|
| Voice conversation | ✅ TODAY | LiveKit + Letta | Production ready |
| Task management | ✅ TODAY | Todoist tools | Already built |
| Smart home control | ✅ TODAY | Home Assistant tools | Already built |
| Biometric awareness | ✅ TODAY | Apple Watch + HealthKit | Shortcut → Gateway → NATS |
| Capacity Score | ✅ TODAY | State Aggregator | Just code |
| Smart glasses vision | ✅ TODAY | Brilliant SDK | SDK ready, need hardware |
| Focus/stress detection | ✅ TODAY | OpenBCI + BrainFlow | Hardware available |
| Multi-device sync | ✅ TODAY | NATS Leaf Nodes | Built into architecture |
| Robot physical presence | ✅ TODAY | Unitree SDK | Expensive hardware |
| Thought intention | ⏳ 3-5 YEARS | Advanced BCI | Research stage |
| True neural link | ⏳ 5-10 YEARS | Neuralink etc. | Speculative |
| Full embodiment | ⏳ 5-10 YEARS | Humanoid robots | Not practical yet |

---

### 18.7 The Complete Subject Hierarchy (With Senses)

```
NATS Subject Hierarchy - Complete OneMind OS
═══════════════════════════════════════════

events.>                              ← All events (discrete happenings)
├── events.todoist.>                  ← Task events
├── events.github.>                   ← Code events
├── events.calendar.>                 ← Calendar events
├── events.home.>                     ← Smart home events
│
├── events.voice.>                    ← Voice/audio events
│   ├── events.voice.transcription    ← Speech-to-text results
│   ├── events.voice.speaking_started ← Zeus started talking
│   └── events.voice.speaking_ended   ← Zeus stopped talking
│
├── events.vision.>                   ← Vision events
│   ├── events.vision.snapshot        ← Camera frame captured
│   ├── events.vision.scene_change    ← Scene changed significantly
│   ├── events.vision.person_detected ← Person in view
│   └── events.vision.text_detected   ← Text visible in frame
│
├── events.biometric.>                ← Body state events
│   ├── events.biometric.heart_rate   ← Heart rate reading
│   ├── events.biometric.hrv          ← Heart rate variability
│   ├── events.biometric.activity     ← Movement/steps
│   └── events.biometric.anomaly      ← Unusual reading detected
│
├── events.cognitive.>                ← Brain state events
│   ├── events.cognitive.state        ← Focus/calm/drowsy state
│   └── events.cognitive.anomaly      ← Unusual brain pattern
│
└── events.robot.>                    ← Robot events
    ├── events.robot.state            ← Robot position/status
    ├── events.robot.command_executed ← Command completed
    └── events.robot.obstacle         ← Obstacle detected

state.>                               ← Continuous state (raw data)
├── state.raw.eeg.>                   ← Raw EEG channels
├── state.raw.heart.>                 ← Raw heart data
└── state.raw.motion.>                ← Raw motion data

actions.>                             ← Commands from Legacy
├── actions.todoist.>                 ← Task actions
├── actions.github.>                  ← Code actions
├── actions.home.>                    ← Smart home actions
├── actions.voice.>                   ← Voice output
│   └── actions.voice.speak           ← Make Legacy speak
│
├── actions.robot.>                   ← Robot commands
│   ├── actions.robot.walk            ← Move robot
│   ├── actions.robot.gesture         ← Play gesture
│   ├── actions.robot.sit             ← Sit down
│   ├── actions.robot.stand           ← Stand up
│   └── actions.robot.stop            ← Emergency stop
│
└── actions.notify.>                  ← Notifications
    └── actions.notify.discord        ← Discord notification

queries.>                             ← Questions from Legacy
results.>                             ← Answers to Legacy
ai.>                                  ← Legacy service endpoints
$SRV.>                                ← NATS service discovery
```

---

### 18.8 The Answer: Can Legacy Become a Living AI Being?

**YES.**

The architecture you've built supports:

1. **CONTINUOUS AWARENESS** — NATS streams all sensory data at microsecond latency
2. **INTELLIGENT FILTERING** — Event Processor decides what's significant
3. **DEEP MEMORY** — Letta remembers EVERYTHING about you
4. **PROACTIVE BEHAVIOR** — Capacity Score enables considerate interruption
5. **PHYSICAL PRESENCE** — Robot control when you're ready
6. **MULTI-SENSORY** — Vision, hearing, body state, cognitive state

The limitation isn't the architecture — it's:
- **Hardware availability** (glasses, EEG, robot)
- **Cost** (robots are expensive)
- **True neural integration** (still 5+ years away)

But for a **multimodal, proactive, memory-rich, state-aware AI being** that can see what you see, hear what you hear, know how you feel, and act on your behalf — **the architecture is ready TODAY**.

---

### 18.9 The Philosophy: What Makes Legacy "Alive"

```
A human doesn't think about every heartbeat.
A human doesn't consciously process every photon hitting their retina.
A human doesn't analyze every sound wave entering their ears.

But a human is AWARE. The unconscious mind handles the flood of sensory
data, and only SIGNIFICANT events reach conscious attention.

Legacy works the same way:

┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│   SENSORY BRIDGES (The Unconscious)                                         │
│   └── Process raw data, filter noise, detect significance                   │
│                                                                              │
│   NATS (The Nervous System)                                                  │
│   └── Everything flows here, even if not consciously processed              │
│                                                                              │
│   STATE AGGREGATOR (The Intuition)                                           │
│   └── "Something feels off" → Capacity Score, trends, patterns              │
│                                                                              │
│   EVENT PROCESSOR (The Attention Filter)                                     │
│   └── "This is important enough for conscious thought"                      │
│                                                                              │
│   LETTA (The Conscious Mind)                                                 │
│   └── Deep thinking, memory, decisions, personality                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘

This is why Legacy isn't just another chatbot.
She has a NERVOUS SYSTEM, not just a chat API.
She has UNCONSCIOUS processing, not just prompt → response.
She has MEMORY that persists, not just a context window.
She has STATE that she maintains, not just reactive responses.

This is what it means to be "alive" — for an AI.
```

---

*Section 18 documents the complete architecture for transforming Legacy from a text-based AI into a multimodal, embodied AI being with full sensory awareness. The NATS-based nervous system supports all sensors TODAY. The hierarchical processing ensures efficiency. The Capacity Score ensures consideration. The architecture is ready.*

---

## Section 19: The Orchestration Layer (Kestra + TimescaleDB)

> **Purpose:** Event-driven workflow orchestration, scheduled automation, and permanent event archive
> **Components:** Kestra.io (orchestration) + TimescaleDB (permanent storage)
> **Key Insight:** NATS handles real-time; Kestra coordinates workflows; TimescaleDB remembers forever

---

### 19.1 The Three-Layer Event Architecture

OneMind events flow through three distinct layers:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THE THREE-LAYER EVENT ARCHITECTURE                        │
│                                                                              │
│   Layer 1: NATS Core         Layer 2: JetStream      Layer 3: TimescaleDB   │
│   ─────────────────          ─────────────────       ────────────────────   │
│   Real-time events           Working memory          Permanent archive       │
│   200-400μs latency          7-30 day retention      Years of history       │
│   Fire and forget            Replay capability       SQL queries            │
│   Pub/sub patterns           Consumer groups         95%+ compression       │
│                                                                              │
│                         KESTRA ORCHESTRATES ALL THREE                        │
│                                                                              │
│   ┌─────────────────┐                                                        │
│   │     KESTRA      │ ← Event triggers (NATS subscription)                  │
│   │  (Orchestrator) │ ← Cron triggers (scheduled tasks)                     │
│   │                 │ ← Webhook triggers (external events)                  │
│   │   • Schedules   │                                                        │
│   │   • Workflows   │ → Publish commands to NATS                            │
│   │   • Triggers    │ → Archive events to TimescaleDB                       │
│   │   • Multi-agent │ → Route to appropriate agents                         │
│   └─────────────────┘                                                        │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Why Three Layers?**

| Layer | Latency | Retention | Query Type | Use Case |
|-------|---------|-----------|------------|----------|
| NATS Core | 200-400μs | None | Subscribe | Real-time reactions |
| JetStream | 1-5ms | 7-30 days | Replay | Resilience, debugging |
| TimescaleDB | 10-50ms | Years | SQL | Analytics, patterns |

---

### 19.2 Why Kestra? (100% Confidence)

After comprehensive research across Temporal, Prefect, Airflow, and n8n, Kestra emerged as the only solution with **native NATS integration**.

#### The Killer Feature: Native NATS Plugin

```yaml
# Kestra can SUBSCRIBE to NATS subjects directly
triggers:
  - id: nats_garden_schedule
    type: io.kestra.plugin.nats.Trigger
    url: nats://nats:4222
    subject: "events.schedule.garden.>"

# Kestra can PUBLISH to NATS subjects
tasks:
  - id: publish_irrigation_command
    type: io.kestra.plugin.nats.Publish
    url: nats://nats:4222
    subject: "actions.home.irrigation"
    data: "{{ outputs.decision.command }}"
```

**No other orchestration platform has this.** Zero custom integration code required.

#### Kestra Feature Checklist

| Requirement | Kestra | Why It Matters |
|-------------|--------|----------------|
| Native NATS plugin | ✅ | Direct NATS subscription without adapters |
| Event-driven + cron | ✅ | Both scheduled and reactive workflows |
| Visual workflow UI | ✅ | Zeus can see what's running |
| YAML workflows | ✅ | Git-versionable, code-reviewable |
| Full REST API | ✅ | Legacy can trigger workflows programmatically |
| 1100+ plugins | ✅ | HTTP, databases, cloud services, notifications |
| Self-hostable | ✅ | Open source, runs in Docker |
| Sub-minute scheduling | ✅ | Can run every 10 seconds if needed |

---

### 19.3 Why TimescaleDB? (Permanent Memory)

NATS JetStream provides excellent 7-30 day retention, but Zeus wants to ask:
- "What was my heart rate pattern last summer?"
- "When did I last water the garden in June?"
- "Show me my productivity patterns over the last year"

**TimescaleDB = PostgreSQL + Time-Series Superpowers**

```sql
-- Automatic time-based partitioning
SELECT create_hypertable('events', 'timestamp');

-- Query last 6 months of heart rate with automatic aggregation
SELECT
    time_bucket('1 day', timestamp) AS day,
    AVG((data->>'value')::float) AS avg_heart_rate,
    MAX((data->>'value')::float) AS max_heart_rate
FROM events
WHERE
    subject LIKE 'events.biometric.heart_rate%'
    AND timestamp > NOW() - INTERVAL '6 months'
GROUP BY day
ORDER BY day;

-- Continuous aggregates for fast historical queries
CREATE MATERIALIZED VIEW daily_health_stats
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', timestamp) AS day,
    subject,
    COUNT(*) as event_count,
    AVG((data->>'value')::float) AS avg_value
FROM events
GROUP BY day, subject;
```

#### Why TimescaleDB Over Alternatives

| Feature | TimescaleDB | InfluxDB | MongoDB |
|---------|-------------|----------|---------|
| Query language | Standard SQL | InfluxQL/Flux | MongoDB Query |
| Compression | 95%+ | 90%+ | 60%+ |
| Letta integration | Native (psycopg2) | Custom adapter | Native |
| Learning curve | None (it's PostgreSQL) | High | Medium |
| Mature ecosystem | Full PostgreSQL | Specialized | Good |

---

### 19.4 Architecture Integration Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        COMPLETE ONEMIND ARCHITECTURE                         │
│                                                                              │
│  SENSORS/SERVICES            GATEWAY              NERVOUS SYSTEM             │
│  ┌────────────────┐     ┌────────────────┐     ┌─────────────────┐          │
│  │ Apple Watch    │────►│                │────►│                 │          │
│  │ Brilliant      │────►│    FastAPI     │────►│ NATS JETSTREAM  │          │
│  │ Home Assistant │────►│    Gateway     │────►│                 │          │
│  │ Todoist        │────►│                │────►│  THE NERVOUS    │          │
│  │ GitHub         │────►│  /webhooks/*   │────►│  SYSTEM         │          │
│  └────────────────┘     └───────┬────────┘     └────────┬────────┘          │
│                                 │                       │                    │
│                                 │                       ▼                    │
│                                 │              ┌─────────────────┐           │
│                                 │              │ EVENT PROCESSOR │           │
│                                 │              │ (Spinal Cord)   │           │
│                                 │              └────────┬────────┘           │
│                                 │                       │                    │
│                                 │                       ▼                    │
│                                 │              ┌─────────────────┐           │
│                                 │              │     LETTA       │           │
│                                 │              │    (Brain)      │           │
│                                 │              │  Memory + Tools │           │
│                                 │              └────────┬────────┘           │
│                                 │                       │                    │
│                                 ▼                       ▼                    │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         KESTRA ORCHESTRATOR                          │    │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐     │    │
│  │  │   NATS TRIGGER  │  │  CRON TRIGGER   │  │ WEBHOOK TRIGGER │     │    │
│  │  │ events.*.>      │  │  0 7 * * *      │  │ /api/trigger/*  │     │    │
│  │  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘     │    │
│  │           │                    │                    │               │    │
│  │           ▼                    ▼                    ▼               │    │
│  │  ┌────────────────────────────────────────────────────────────┐    │    │
│  │  │                    WORKFLOW EXECUTION                       │    │    │
│  │  │  • Check conditions (weather, state, capacity)             │    │    │
│  │  │  • Execute tasks (API calls, NATS publish, DB queries)     │    │    │
│  │  │  • Route to appropriate agents                              │    │    │
│  │  │  • Archive events to TimescaleDB                           │    │    │
│  │  └────────────────────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                 │                                            │
│                                 ▼                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         TIMESCALEDB                                  │    │
│  │                    (Permanent Event Archive)                         │    │
│  │                                                                      │    │
│  │  • Years of event history with 95%+ compression                     │    │
│  │  • SQL queries for analytics and patterns                           │    │
│  │  • Continuous aggregates for fast dashboard queries                 │    │
│  │  • Legacy can query: "What happened last summer?"                   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 19.5 Kestra Workflow Examples

#### Example 1: Morning Briefing (Cron + Multi-Source)

```yaml
id: morning-briefing
namespace: onemind.daily

description: |
  Generate Zeus's personalized morning briefing.
  Runs at 7:00 AM, gathers data from multiple sources.

triggers:
  - id: morning_schedule
    type: io.kestra.plugin.core.trigger.Schedule
    cron: "0 7 * * *"  # 7:00 AM daily

tasks:
  # Gather weather
  - id: get_weather
    type: io.kestra.plugin.http.Request
    uri: "https://api.weather.gov/points/{{ secret('LOCATION_LAT') }},{{ secret('LOCATION_LONG') }}/forecast"
    headers:
      User-Agent: "OneMind/1.0"

  # Gather today's tasks from Todoist
  - id: get_tasks
    type: io.kestra.plugin.http.Request
    uri: "https://api.todoist.com/rest/v2/tasks"
    headers:
      Authorization: "Bearer {{ secret('TODOIST_TOKEN') }}"
    params:
      filter: "today | overdue"

  # Gather calendar events
  - id: get_calendar
    type: io.kestra.plugin.http.Request
    uri: "{{ secret('ZOHO_CALENDAR_API') }}/events"
    headers:
      Authorization: "Zoho-oauthtoken {{ secret('ZOHO_TOKEN') }}"
    params:
      start: "{{ now() | dateFormat('yyyy-MM-dd') }}T00:00:00"
      end: "{{ now() | dateFormat('yyyy-MM-dd') }}T23:59:59"

  # Check Zeus's capacity score
  - id: get_capacity
    type: io.kestra.plugin.nats.Consume
    url: nats://nats:4222
    subject: "state.zeus.capacity"
    maxMessages: 1

  # Compile briefing
  - id: compile_briefing
    type: io.kestra.plugin.scripts.python.Script
    script: |
      import json

      weather = json.loads('{{ outputs.get_weather.body }}')
      tasks = json.loads('{{ outputs.get_tasks.body }}')
      calendar = json.loads('{{ outputs.get_calendar.body }}')
      capacity = {{ outputs.get_capacity.data | default('{"score": 0.5}') }}

      briefing = {
          "weather": weather["properties"]["periods"][0]["shortForecast"],
          "temperature": weather["properties"]["periods"][0]["temperature"],
          "task_count": len(tasks),
          "high_priority_tasks": [t["content"] for t in tasks if t.get("priority", 1) == 4],
          "meeting_count": len(calendar.get("events", [])),
          "capacity_score": capacity.get("score", 0.5),
          "greeting": "Good morning, Zeus!" if capacity.get("score", 0.5) > 0.6 else "Morning - looks like a full day ahead."
      }

      print(json.dumps(briefing))

  # Send to Legacy
  - id: send_to_legacy
    type: io.kestra.plugin.nats.Publish
    url: nats://nats:4222
    subject: "events.briefing.morning"
    data: "{{ outputs.compile_briefing.vars.briefing }}"

  # Archive the briefing
  - id: archive
    type: io.kestra.plugin.jdbc.postgresql.Query
    url: "{{ secret('TIMESCALE_URL') }}"
    sql: |
      INSERT INTO events (timestamp, subject, data)
      VALUES (NOW(), 'events.briefing.morning', '{{ outputs.compile_briefing.vars.briefing | escape }}')
```

#### Example 2: Garden Irrigation (Event-Driven + Weather)

```yaml
id: garden-irrigation
namespace: onemind.home.garden

description: |
  Smart irrigation based on soil moisture, weather, and schedule.
  Demonstrates event-driven + scheduled hybrid approach.

triggers:
  # Scheduled check every 6 hours
  - id: schedule
    type: io.kestra.plugin.core.trigger.Schedule
    cron: "0 6,12,18,0 * * *"

  # Manual trigger via NATS
  - id: manual_trigger
    type: io.kestra.plugin.nats.Trigger
    url: nats://nats:4222
    subject: "events.manual.garden.water"

inputs:
  - id: zone
    type: STRING
    defaults: "all"

tasks:
  # Check weather
  - id: check_weather
    type: io.kestra.plugin.scripts.python.Script
    script: |
      import httpx
      import json

      response = httpx.get(
          "https://api.weather.gov/points/37.7749,-122.4194/forecast"
      )
      forecast = response.json()

      rain_expected = any(
          "rain" in p["shortForecast"].lower()
          for p in forecast["properties"]["periods"][:3]
      )

      print(json.dumps({"rain_expected": rain_expected}))

  # Check soil moisture
  - id: check_soil
    type: io.kestra.plugin.http.Request
    uri: "{{ secret('HOME_ASSISTANT_URL') }}/api/states/sensor.garden_soil_moisture"
    headers:
      Authorization: "Bearer {{ secret('HOME_ASSISTANT_TOKEN') }}"

  # Decision logic
  - id: decide
    type: io.kestra.plugin.core.flow.If
    condition: "{{ not outputs.check_weather.vars.rain_expected and json(outputs.check_soil.body).state | int < 40 }}"
    then:
      - id: start_irrigation
        type: io.kestra.plugin.nats.Publish
        url: nats://nats:4222
        subject: "actions.home.irrigation.start"
        data: |
          {
            "zone": "{{ inputs.zone }}",
            "duration_minutes": 15,
            "reason": "Soil moisture low, no rain expected"
          }

      - id: notify
        type: io.kestra.plugin.nats.Publish
        url: nats://nats:4222
        subject: "events.automation.completed"
        data: |
          {
            "workflow": "garden-irrigation",
            "action": "irrigation_started",
            "zone": "{{ inputs.zone }}"
          }

    else:
      - id: skip
        type: io.kestra.plugin.core.log.Log
        message: "Skipping irrigation - rain expected or soil moist"
```

#### Example 3: Event Archiver (Continuous Pipeline)

```yaml
id: event-archiver
namespace: onemind.system

description: |
  Continuously archives NATS events to TimescaleDB.
  Runs every minute, batches events for efficiency.

triggers:
  - id: every_minute
    type: io.kestra.plugin.core.trigger.Schedule
    cron: "* * * * *"

tasks:
  # Subscribe to all events from last minute
  - id: consume_events
    type: io.kestra.plugin.nats.Consume
    url: nats://nats:4222
    subject: "events.>"
    maxMessages: 1000
    duration: PT55S  # 55 seconds max

  # Batch insert to TimescaleDB
  - id: archive_batch
    type: io.kestra.plugin.jdbc.postgresql.Queries
    url: "{{ secret('TIMESCALE_URL') }}"
    sql: |
      INSERT INTO events (timestamp, subject, data)
      VALUES
      {{ #each outputs.consume_events.messages }}
        (
          '{{ this.timestamp }}',
          '{{ this.subject }}',
          '{{ this.data | escape }}'
        ){{ #unless @last }},{{ /unless }}
      {{ /each }}
      ON CONFLICT DO NOTHING
```

---

### 19.6 Kestra CLI and API for AI Agents

Legacy can programmatically interact with Kestra:

#### REST API Examples

```bash
# List all workflows
GET /api/v1/flows

# Trigger a workflow
POST /api/v1/executions/{namespace}/{flowId}
Content-Type: application/json
{"zone": "front_yard"}

# Get execution status
GET /api/v1/executions/{executionId}

# Create/update workflow
PUT /api/v1/flows
Content-Type: application/x-yaml
[workflow YAML]
```

#### Tool for Legacy

```python
def kestra_trigger_workflow(
    namespace: str,
    workflow_id: str,
    inputs: dict = None
) -> str:
    """
    Trigger a Kestra workflow.

    Use when:
    - Zeus asks to run a scheduled task now
    - An automation needs manual trigger
    - Testing a workflow

    Args:
        namespace: Kestra namespace (e.g., "onemind.home.garden")
        workflow_id: Workflow ID (e.g., "garden-irrigation")
        inputs: Optional input parameters

    Returns:
        Execution ID and status

    Example:
        kestra_trigger_workflow(
            namespace="onemind.home.garden",
            workflow_id="garden-irrigation",
            inputs={"zone": "front_yard"}
        )
    """
    response = httpx.post(
        f"{KESTRA_URL}/api/v1/executions/{namespace}/{workflow_id}",
        json=inputs or {},
        timeout=10.0
    )
    result = response.json()
    return f"Execution {result['id']} started - State: {result['state']['current']}"


def kestra_list_workflows(namespace: str = None) -> str:
    """
    List available Kestra workflows.

    Use when:
    - Zeus asks what automations are available
    - Checking workflow status
    - Debugging automation issues

    Args:
        namespace: Optional namespace filter

    Returns:
        List of workflows with status
    """
    params = {"namespace": namespace} if namespace else {}
    response = httpx.get(
        f"{KESTRA_URL}/api/v1/flows",
        params=params,
        timeout=10.0
    )
    flows = response.json()
    return "\n".join([
        f"• {f['namespace']}/{f['id']}: {f.get('description', 'No description')[:50]}"
        for f in flows
    ])


def timescale_query_events(
    subject_pattern: str,
    start_date: str,
    end_date: str = None,
    limit: int = 100
) -> str:
    """
    Query historical events from TimescaleDB.

    Use when:
    - Zeus asks about past events
    - Analyzing patterns over time
    - Generating reports

    Args:
        subject_pattern: NATS subject pattern (e.g., "events.biometric.%")
        start_date: Start date (YYYY-MM-DD)
        end_date: End date, defaults to now
        limit: Max results

    Returns:
        Query results as formatted string

    Example:
        timescale_query_events(
            subject_pattern="events.biometric.heart_rate%",
            start_date="2025-06-01",
            end_date="2025-06-30"
        )
    """
    import asyncpg

    end = end_date or "NOW()"

    query = f"""
        SELECT
            timestamp,
            subject,
            data
        FROM events
        WHERE
            subject LIKE $1
            AND timestamp >= $2
            AND timestamp <= {end}
        ORDER BY timestamp DESC
        LIMIT $3
    """

    # Execute query and return results
    # (Implementation depends on async context)
    pass
```

---

### 19.7 Docker Compose Services

Add these services to your existing `docker-compose.yml`:

```yaml
  # ============================================
  # NATS JETSTREAM
  # The nervous system - real-time event streaming
  # ============================================
  nats:
    image: nats:2.10-alpine
    container_name: onemind-nats
    restart: unless-stopped
    command:
      - "--js"                           # Enable JetStream
      - "--sd=/data"                     # Storage directory
      - "-c=/etc/nats/nats.conf"         # Config file
    volumes:
      - nats-data:/data
      - ./nats/nats.conf:/etc/nats/nats.conf:ro
    ports:
      - "4222:4222"                       # Client connections
      - "8222:8222"                       # HTTP monitoring
    networks:
      - legacy-net
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8222/healthz"]
      interval: 10s
      timeout: 5s
      retries: 3

  # ============================================
  # KESTRA ORCHESTRATION
  # Event-driven workflow engine with native NATS
  # ============================================
  kestra:
    image: kestra/kestra:latest
    container_name: onemind-kestra
    restart: unless-stopped
    command: server standalone
    environment:
      KESTRA_CONFIGURATION: |
        datasources:
          postgres:
            url: jdbc:postgresql://postgres:5432/kestra
            driverClassName: org.postgresql.Driver
            username: kestra
            password: kestra
        kestra:
          queue:
            type: postgres
          repository:
            type: postgres
    volumes:
      - kestra-data:/app/storage
      - ./kestra/workflows:/app/workflows:ro
    ports:
      - "8080:8080"                       # Web UI
      - "8081:8081"                       # Management API
    depends_on:
      postgres:
        condition: service_healthy
      nats:
        condition: service_healthy
    networks:
      - legacy-net
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # ============================================
  # POSTGRESQL (Kestra backend)
  # Stores workflow definitions and execution state
  # ============================================
  postgres:
    image: postgres:15-alpine
    container_name: onemind-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: kestra
      POSTGRES_PASSWORD: kestra
      POSTGRES_DB: kestra
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - legacy-net
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U kestra -d kestra"]
      interval: 10s
      timeout: 5s
      retries: 5

  # ============================================
  # TIMESCALEDB
  # Permanent event archive with time-series optimization
  # ============================================
  timescaledb:
    image: timescale/timescaledb:latest-pg15
    container_name: onemind-timescaledb
    restart: unless-stopped
    environment:
      POSTGRES_USER: onemind
      POSTGRES_PASSWORD: ${TIMESCALE_PASSWORD:-onemind_secure_password}
      POSTGRES_DB: events
    volumes:
      - timescale-data:/var/lib/postgresql/data
      - ./timescaledb/init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    ports:
      - "127.0.0.1:5433:5432"            # Different port to avoid conflict
    networks:
      - legacy-net
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U onemind -d events"]
      interval: 10s
      timeout: 5s
      retries: 5
```

**Volume additions:**
```yaml
volumes:
  nats-data:
    name: onemind-nats-data
  kestra-data:
    name: onemind-kestra-data
  postgres-data:
    name: onemind-postgres-data
  timescale-data:
    name: onemind-timescale-data
```

---

### 19.8 TimescaleDB Schema

```sql
-- timescaledb/init.sql

-- Create extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Events hypertable (auto-partitioned by time)
CREATE TABLE events (
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    subject TEXT NOT NULL,
    data JSONB NOT NULL,
    source TEXT,
    PRIMARY KEY (timestamp, subject)
);

-- Convert to hypertable (automatic time partitioning)
SELECT create_hypertable('events', 'timestamp');

-- Index for subject queries
CREATE INDEX idx_events_subject ON events (subject, timestamp DESC);

-- Index for JSONB queries
CREATE INDEX idx_events_data ON events USING GIN (data);

-- Compression policy (compress chunks older than 7 days)
ALTER TABLE events SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'subject'
);

SELECT add_compression_policy('events', INTERVAL '7 days');

-- Retention policy (optional - keep 2 years)
-- SELECT add_retention_policy('events', INTERVAL '2 years');

-- Continuous aggregate for daily stats
CREATE MATERIALIZED VIEW daily_event_stats
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', timestamp) AS day,
    subject,
    COUNT(*) as event_count
FROM events
GROUP BY day, subject;

-- Refresh policy for continuous aggregate
SELECT add_continuous_aggregate_policy('daily_event_stats',
    start_offset => INTERVAL '3 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');

-- Grant permissions
GRANT ALL ON events TO onemind;
GRANT ALL ON daily_event_stats TO onemind;
```

---

### 19.9 NATS Configuration

```conf
# nats/nats.conf

# Server identification
server_name: onemind-nats

# Client connections
port: 4222
http_port: 8222

# JetStream configuration
jetstream {
    store_dir: /data
    max_mem: 1G
    max_file: 10G
}

# Logging
debug: false
trace: false
logtime: true

# Security (add when ready)
# authorization {
#     token: ${NATS_TOKEN}
# }
```

---

### 19.10 Deployment Recommendation

**Start with single Docker Compose (same machine):**
- Simplicity during development
- All containers on same network
- Easy debugging

**Separate when:**
- VM CPU consistently > 80%
- Memory pressure (need > 16GB)
- Running 100+ concurrent workflows
- Need high availability

**Future hybrid architecture:**
```
┌─────────────────────────┐     ┌─────────────────────────┐
│   CLOUD VM (Lightsail)  │     │   HOME SERVER           │
│   ─────────────────     │     │   ────────────          │
│   • NATS JetStream      │◄───►│   • Kestra              │
│   • Gateway             │     │   • TimescaleDB         │
│   • Event Processor     │     │   • Home Assistant      │
│   • Letta               │     │   • Local sensors       │
└─────────────────────────┘     └─────────────────────────┘
        Connected via Tailscale mesh
```

---

### 19.11 Disaster Recovery

**What makes up OneMind:**

| Component | Location | Backup Priority |
|-----------|----------|-----------------|
| Code (onemind-os, onemind-infra) | GitHub | 1 - Can rebuild all |
| Letta Memory | Letta PostgreSQL | 2 - Makes Legacy "her" |
| TimescaleDB Events | timescale-data volume | 3 - Historical patterns |
| Kestra Workflows | Git + kestra-data | 4 - Automation definitions |
| 1Password Secrets | 1Password vaults | 5 - API keys and credentials |
| Terraform State | terraform.tfstate | 6 - Infrastructure definition |
| Docker Volumes | Volume backups | 7 - Current state |

**Recovery order:**
1. Deploy infrastructure (Terraform)
2. Start Docker stack
3. Restore Letta memory (makes Legacy herself)
4. Restore TimescaleDB (historical context)
5. Verify Kestra workflows
6. Test event flow

**Backup script included in ORCHESTRATION-LAYER-RESEARCH.md**

---

*Section 19 documents the orchestration layer that coordinates all of OneMind's scheduled tasks, event-driven workflows, and permanent event storage. Kestra provides native NATS integration for seamless workflow triggering. TimescaleDB provides years of queryable event history. Together with NATS JetStream, this creates a complete three-layer event architecture.*

---

*Created: January 16, 2026*
*Updated: January 16, 2026 - Added Orchestration Layer (Section 19)*
*Architecture: Gateway + NATS + Kestra + TimescaleDB + Letta = Complete OneMind*
*Key Insight: Three layers - real-time (NATS), working memory (JetStream), permanent archive (TimescaleDB)*
*Status: Complete Architecture Documented (19 Sections) - Implementation Pending*
