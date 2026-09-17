# Dual Heartbeat System

> Two independent heartbeats: Agno (slow, digital awareness, every 15 min) and OM1 (fast, physical survival, every 0.5 sec). They share context through NATS but crash independently.

## Why Two Heartbeats

| Factor | OM1 Heartbeat | Agno Heartbeat |
|--------|---------------|----------------|
| **Purpose** | Physical survival | Digital awareness |
| **Frequency** | Every 0.5-2 seconds | Every 5-30 minutes |
| **LLM** | #2 (fast, local, cheap) | #1 (powerful, cloud, costs tokens) |
| **Cost per hour** | ~3,600 calls (tiny model) | ~2-4 calls (big model) |
| **If it stops** | Robot crashes into wall | You miss an email |
| **Example** | "Wall ahead — turn left" | "Meeting in 15 min, no prep doc" |

Combining them into one loop would either make the robot too slow (waiting for email checks) or make digital checks absurdly expensive (running GPT-4.1 twice per second).

## OM1 Heartbeat — Fast Loop (Physical)

This is OM1's native cognitive loop. It runs continuously and NEVER stops while the robot is powered on.

```
Every 0.5-2 seconds:
  1. Read ALL sensors (camera, mic, LiDAR, IMU)
  2. Convert sensor data to text (NLDB)
  3. Check NATS for any Agno commands
  4. Fuse everything into one prompt
  5. LLM #2 decides: move? speak? wait? dodge?
  6. Execute action via HAL
  7. Publish sensor summary to NATS
  8. Repeat immediately
```

This is built into OM1's core runtime — no custom implementation needed. You configure it via `cortex_rate_hz` in the OM1 config.

```json5
{
  "cortex_rate_hz": 1.0,  // 1 cycle per second
  "llm": {
    "model": "qwen3-30b",
    "provider": "ollama"
  }
}
```

## Agno Heartbeat — Slow Loop (Digital)

Inspired by the OpenClaw heartbeat pattern. A cron-like scheduler triggers an Agno agent run every N minutes with a checklist of things to monitor.

### Architecture

```
┌─────────────────────────────────────────────┐
│           AGNO HEARTBEAT SERVICE             │
│                                              │
│  Cron Scheduler (every 15 minutes)           │
│       │                                      │
│       ▼                                      │
│  Step 1: GATHER (Python, no LLM)             │
│    ├── Gmail API → new emails?               │
│    ├── Calendar API → upcoming events?       │
│    ├── ClickUp API → overdue tasks?          │
│    ├── NATS → fleet world state              │
│    ├── NATS → any robot alerts?              │
│    └── Flock Safety → camera alerts?         │
│       │                                      │
│       ▼                                      │
│  Step 2: BUILD PROMPT                        │
│    ├── Load HEARTBEAT.md checklist           │
│    ├── Inject gathered data                  │
│    └── Create heartbeat run payload          │
│       │                                      │
│       ▼                                      │
│  Step 3: REASON (LLM #1)                    │
│    ├── Agno agent processes heartbeat        │
│    ├── Decides what needs attention          │
│    └── Returns actions or HEARTBEAT_OK       │
│       │                                      │
│       ▼                                      │
│  Step 4: ACT                                 │
│    ├── Send notifications (if needed)        │
│    ├── Publish robot commands (if needed)     │
│    └── Log result                            │
│       │                                      │
│       ▼                                      │
│  Step 5: SLEEP until next cycle              │
└─────────────────────────────────────────────┘
```

### HEARTBEAT.md — The Checklist

Stored in the Agno workspace. The agent reads this on every heartbeat cycle to know what to monitor. Modifiable at runtime by telling Legacy "update my heartbeat checklist."

```markdown
# One Mind OS Heartbeat Checklist

## Digital Checks
- Scan inbox for emails from family or business contacts
- Check calendar for events in next 2 hours
- Review ClickUp for overdue or blocked tasks
- Check Slack/Discord for unread @mentions

## Physical Checks (via NATS)
- Robot fleet battery levels — alert if any below 20%
- Unusual sensor events (unknown person, loud noise)
- Any robot errors or connectivity issues

## Proactive Actions
- If meeting in < 30 min with no prep: remind me
- If VIP email unanswered > 2 hours: nudge me
- If all clear: return HEARTBEAT_OK (silent)

## Rules
- Never notify between 10 PM and 7 AM unless urgent
- Max 3 notifications per heartbeat cycle
- Always include actionable next step in notifications
```

### Implementation

```python
import asyncio, json, time
import httpx, nats
from datetime import datetime

class AgnoHeartbeat:
    def __init__(
        self,
        agno_url: str = "http://localhost:7777",
        nats_url: str = "nats://localhost:4222",
        interval_minutes: int = 15,
        user_id: str = "master-delacruz",
        active_hours: tuple = (7, 22),  # 7 AM - 10 PM
    ):
        self.agno_url = agno_url
        self.nats_url = nats_url
        self.interval = interval_minutes * 60
        self.user_id = user_id
        self.active_hours = active_hours
        self.session_id = f"{user_id}-heartbeat"

    async def start(self):
        """Run the heartbeat loop forever."""
        while True:
            hour = datetime.now().hour
            if self.active_hours[0] <= hour < self.active_hours[1]:
                await self._run_heartbeat()
            await asyncio.sleep(self.interval)

    async def _run_heartbeat(self):
        # Step 1: Gather context (no LLM needed)
        context = await self._gather_context()

        # Step 2: Load checklist
        checklist = self._load_checklist()

        # Step 3: Build prompt and send to Agno
        prompt = self._build_prompt(checklist, context)
        response = await self._send_to_agno(prompt)

        # Step 4: Act on response
        if "HEARTBEAT_OK" not in response:
            await self._notify(response)

    async def _gather_context(self) -> dict:
        context = {}

        # Gather from NATS — robot fleet status
        try:
            nc = await nats.connect(self.nats_url)
            resp = await nc.request("onemind.fleet.worldstate", b"query", timeout=3.0)
            context["fleet"] = json.loads(resp.data.decode())
            await nc.drain()
        except Exception:
            context["fleet"] = {"error": "NATS unavailable"}

        # Gather from Gmail API
        try:
            context["emails"] = await self._check_gmail()
        except Exception:
            context["emails"] = {"error": "Gmail check failed"}

        # Gather from Calendar API
        try:
            context["calendar"] = await self._check_calendar()
        except Exception:
            context["calendar"] = {"error": "Calendar check failed"}

        # Gather from ClickUp API
        try:
            context["tasks"] = await self._check_clickup()
        except Exception:
            context["tasks"] = {"error": "ClickUp check failed"}

        return context

    def _load_checklist(self) -> str:
        try:
            with open("workspace/HEARTBEAT.md", "r") as f:
                return f.read()
        except FileNotFoundError:
            return "Check for anything urgent. If nothing, return HEARTBEAT_OK."

    def _build_prompt(self, checklist: str, context: dict) -> str:
        return f"""[HEARTBEAT RUN]

Follow this checklist strictly. Do not infer or repeat old tasks.
If nothing needs attention, reply HEARTBEAT_OK.

## Checklist
{checklist}

## Current Context
{json.dumps(context, indent=2)}

## Rules
- Only report actionable items
- Include specific next steps
- Max 3 notifications
- If all clear, return HEARTBEAT_OK"""

    async def _send_to_agno(self, prompt: str) -> str:
        async with httpx.AsyncClient(timeout=30.0) as client:
            resp = await client.post(
                f"{self.agno_url}/teams/onemind/runs",
                json={
                    "message": prompt,
                    "session_id": self.session_id,
                    "user_id": self.user_id,
                    "stream": False,
                },
            )
            return resp.json().get("content", "HEARTBEAT_OK")

    async def _notify(self, message: str):
        """Send notification to user's preferred channel."""
        # Via Interaction Gateway broadcast
        async with httpx.AsyncClient() as client:
            await client.post(
                "http://localhost:8000/notify",
                json={"user_id": self.user_id, "message": message},
            )

    # Data source helpers (implement per service)
    async def _check_gmail(self) -> dict:
        # Use Gmail API or MCP to check for new emails
        pass

    async def _check_calendar(self) -> dict:
        # Use Google Calendar API to check upcoming events
        pass

    async def _check_clickup(self) -> dict:
        # Use ClickUp API to check overdue tasks
        pass
```

### Cron-based Alternative

For simpler deployments, use system cron instead of a Python loop:

```bash
# /etc/cron.d/onemind-heartbeat
*/15 7-22 * * * onemind python /opt/onemind/heartbeat/run.py
```

## How They Feed Each Other

```
OM1 Fast Loop                    NATS                     Agno Slow Loop
(every 0.5s)                                              (every 15min)
     │                                                         │
     ├── sees person ──────► onemind.vision.robot1 ──────► reads world state
     │                                                         │
     ├── battery 15% ──────► onemind.status.robot1 ──────► "robot1 low battery"
     │                                                         │
     │                      onemind.commands.robot1 ◄────── "go charge"
     │                                                         │
     ├── receives "charge" ◄────────────────────────────── published command
     │                                                         │
     └── walks to charger                                      └── logs heartbeat
```

OM1's fast heartbeat continuously publishes sensor data. Agno's slow heartbeat reads that data when it wakes up. If Agno decides to act (e.g., "send robot to charge"), it publishes a command that OM1 picks up on its next fast cycle (within 0.5 seconds).


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
