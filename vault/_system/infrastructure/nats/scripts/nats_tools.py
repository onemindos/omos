"""
NATS Tools — Give Legacy Eyes and Hands on the Event Bus
=========================================================

These tools let the Legacy agent:
1. Observe: See what's happening on NATS topics
2. Publish: Send commands/events to NATS
3. System status: Get system-wide health snapshot

Registered in backend/tools/registry.py as function-based tools.
"""

import json
import logging
from datetime import datetime
from typing import Optional

logger = logging.getLogger("onemind.tools.nats")


def nats_observe(topic: str = "system.>", limit: int = 10) -> str:
    """
    Observe recent events from the NATS event bus.

    Use this to see what's happening in the system right now.
    Topics follow the pattern: {category}.{identifier}

    Common topics:
    - system.> — System health and heartbeats
    - digital.> — Digital agent activity
    - physical.> — Physical entity activity (robots, sensors)
    - nexus.> — Nexus dispatch events
    - comms.> — Communication protocol events
    - events.> — Event gateway
    - human.> — Human-in-the-loop events

    Args:
        topic: NATS subject pattern to observe (supports wildcards: >, *)
        limit: Maximum number of recent events to return

    Returns:
        JSON string with recent events on the topic
    """
    import asyncio

    async def _observe():
        try:
            from backend.bus import get_bus
            nc = await get_bus()
            if not nc:
                return json.dumps({
                    "status": "disconnected",
                    "message": "NATS not available",
                    "events": [],
                })

            # Collect messages for a short window
            events = []
            collected = asyncio.Event()

            async def handler(msg):
                try:
                    data = json.loads(msg.data.decode()) if msg.data else {}
                except (json.JSONDecodeError, UnicodeDecodeError):
                    data = {"raw": msg.data.decode() if msg.data else ""}

                events.append({
                    "subject": msg.subject,
                    "data": data,
                    "timestamp": datetime.utcnow().isoformat(),
                })
                if len(events) >= limit:
                    collected.set()

            sub = await nc.subscribe(topic, cb=handler)

            # Wait up to 2 seconds for messages
            try:
                await asyncio.wait_for(collected.wait(), timeout=2.0)
            except asyncio.TimeoutError:
                pass

            await sub.unsubscribe()

            return json.dumps({
                "status": "ok",
                "topic": topic,
                "events_captured": len(events),
                "events": events,
            }, default=str)

        except Exception as e:
            return json.dumps({
                "status": "error",
                "message": str(e),
                "events": [],
            })

    try:
        loop = asyncio.get_event_loop()
        if loop.is_running():
            # We're already in an async context; create a future
            import concurrent.futures
            # For Agno tool calls that happen in async context
            future = asyncio.ensure_future(_observe())
            # Can't block here, return a note
            return json.dumps({
                "status": "note",
                "message": f"Observation request submitted for topic: {topic}",
                "tip": "Use system_overview for immediate status",
            })
    except RuntimeError:
        pass

    return asyncio.run(_observe())


def nats_publish(topic: str, message: str, event_type: str = "command") -> str:
    """
    Publish a message to the NATS event bus.

    Use this to send commands, events, or notifications to the system.

    Args:
        topic: NATS subject to publish to (e.g., "digital.agent.action", "system.alert")
        message: The message content to publish
        event_type: Type of event (command, alert, notification, status)

    Returns:
        Confirmation of publish
    """
    import asyncio

    async def _publish():
        try:
            from backend.bus import publish
            data = {
                "type": event_type,
                "content": message,
                "source": "legacy",
                "timestamp": datetime.utcnow().isoformat(),
            }
            success = await publish(topic, data)
            if success:
                return json.dumps({
                    "status": "published",
                    "topic": topic,
                    "event_type": event_type,
                    "message": message[:200],
                })
            else:
                return json.dumps({
                    "status": "failed",
                    "message": "NATS not available",
                })
        except Exception as e:
            return json.dumps({
                "status": "error",
                "message": str(e),
            })

    try:
        loop = asyncio.get_event_loop()
        if loop.is_running():
            future = asyncio.ensure_future(_publish())
            return json.dumps({
                "status": "submitted",
                "topic": topic,
                "message": f"Publishing to {topic}",
            })
    except RuntimeError:
        pass

    return asyncio.run(_publish())


def system_overview() -> str:
    """
    Get the current system overview showing all agents, teams, devices,
    NATS connectivity, and heartbeat status.

    Use this to understand what's available in the system right now.

    Returns:
        JSON string with complete system overview
    """
    overview = {
        "timestamp": datetime.utcnow().isoformat(),
        "agents": [],
        "teams": [],
        "nats": {},
        "heartbeat": {},
        "tools": {},
    }

    # Get agents info
    try:
        from backend.agents.presets import AGENT_PRESETS
        for agent_id, preset in AGENT_PRESETS.items():
            overview["agents"].append({
                "id": agent_id,
                "name": preset.get("name", agent_id),
                "model": preset.get("model_id", "unknown"),
                "tools": preset.get("tools", []),
                "memory_enabled": preset.get("enable_user_memories", False),
                "reasoning_enabled": preset.get("enable_reasoning", False),
            })
    except Exception as e:
        overview["agents_error"] = str(e)

    # Get teams info
    try:
        from backend.agents.teams.presets import TEAM_PRESETS
        for team_id, preset in TEAM_PRESETS.items():
            overview["teams"].append({
                "id": team_id,
                "name": preset.get("name", team_id),
                "members": preset.get("members", []),
                "mode": preset.get("mode", "coordinate"),
            })
    except Exception as e:
        overview["teams_error"] = str(e)

    # Get NATS status
    try:
        from backend.bus import get_bus_status
        overview["nats"] = get_bus_status()
    except Exception as e:
        overview["nats"] = {"status": "unavailable", "error": str(e)}

    # Get heartbeat status
    try:
        from backend.heartbeat import get_heartbeat
        hb = get_heartbeat()
        if hb:
            overview["heartbeat"] = hb.get_status()
        else:
            overview["heartbeat"] = {"status": "not_initialized"}
    except Exception as e:
        overview["heartbeat"] = {"status": "error", "error": str(e)}

    # Get tool availability
    try:
        from backend.tools.registry import get_tool_availability
        overview["tools"] = get_tool_availability()
    except Exception as e:
        overview["tools"] = {"error": str(e)}

    return json.dumps(overview, default=str, indent=2)
