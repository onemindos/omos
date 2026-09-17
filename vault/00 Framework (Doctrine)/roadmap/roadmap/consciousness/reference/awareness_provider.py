"""
Agent Awareness Provider
========================

Provides real-time context from the Consciousness layer to Agno agents.

This is the bridge between:
- NATS Event Bus (recent events)
- KV Store (current state)
- State Manager (user context)
- Temporal Service (deadlines, time context)

Usage:
    from backend.consciousness.awareness_provider import get_agent_context

    # In your agent run
    context = await get_agent_context()
    response = await agent.arun(message, run_context={"awareness": context})
"""

import logging
from datetime import datetime, timedelta
from typing import Any, Dict, List, Optional

logger = logging.getLogger("agno.awareness")


async def get_recent_events(
    count: int = 10,
    event_type: Optional[str] = None,
) -> List[Dict[str, Any]]:
    """
    Get recent events from NATS.

    Args:
        count: Number of recent events
        event_type: Filter by type (e.g., "todoist", "calendar")

    Returns:
        List of event dicts with type, timestamp, summary
    """
    try:
        from backend.agno.agentos_setup import get_event_bus

        bus = get_event_bus()
        if not bus:
            return []

        events = await bus.get_recent(count=count, event_type=event_type)

        return [
            {
                "type": e.type,
                "source": e.source,
                "timestamp": e.timestamp.isoformat() if e.timestamp else None,
                "summary": e.description or f"{e.source}.{e.type}",
                "priority": e.priority,
            }
            for e in events
        ]
    except Exception as e:
        logger.warning(f"Failed to get recent events: {e}")
        return []


async def get_current_state() -> Dict[str, Any]:
    """
    Get current awareness state from KV Store.

    Returns:
        Dict with awareness_level, focus, location, etc.
    """
    try:
        from backend.agno.agentos_setup import get_consciousness_kv_store

        kv = get_consciousness_kv_store()
        if not kv:
            return {}

        snapshot = await kv.get_state_snapshot()

        return {
            "awareness_level": snapshot.awareness_level.name,
            "focus_task": snapshot.focus_task,
            "focus_type": snapshot.focus_type,
            "last_interaction": snapshot.last_interaction.isoformat() if snapshot.last_interaction else None,
        }
    except Exception as e:
        logger.warning(f"Failed to get current state: {e}")
        return {}


async def get_state_summary() -> str:
    """
    Get human-readable state summary.

    Returns:
        String like "📍 Location: Home | 🎯 Activity: Working (45min) | 💻 Task: Code review"
    """
    try:
        from backend.consciousness.awareness.state import get_state_manager

        manager = get_state_manager()
        if not manager:
            return ""

        return manager.state.get_context_summary()
    except Exception as e:
        logger.debug(f"State manager not available: {e}")
        return ""


async def get_upcoming_deadlines(hours: int = 24) -> List[Dict[str, Any]]:
    """
    Get upcoming deadlines.

    Args:
        hours: Look-ahead window

    Returns:
        List of deadline dicts with title, due, urgency
    """
    try:
        from backend.consciousness.temporal import get_temporal_service, DeadlineUrgency

        temporal = get_temporal_service()
        if not temporal:
            return []

        # Get deadlines with urgency SOON or higher (< 24 hours)
        # If hours > 24, also include UPCOMING (< 1 week)
        urgencies = [DeadlineUrgency.OVERDUE, DeadlineUrgency.CRITICAL,
                     DeadlineUrgency.URGENT, DeadlineUrgency.SOON]
        if hours > 24:
            urgencies.append(DeadlineUrgency.UPCOMING)

        all_deadlines = []
        for urgency in urgencies:
            deadlines = await temporal.get_deadlines(urgency=urgency)
            all_deadlines.extend(deadlines)

        return [
            {
                "title": d.title,
                "due": d.due.isoformat() if d.due else None,
                "urgency": d.urgency.value if d.urgency else 5,
                "source": d.source,
            }
            for d in all_deadlines
        ]
    except Exception as e:
        logger.debug(f"Temporal service not available: {e}")
        return []


async def get_agent_context(
    include_events: bool = True,
    include_state: bool = True,
    include_deadlines: bool = True,
    event_count: int = 10,
    deadline_hours: int = 24,
) -> str:
    """
    Build comprehensive context string for agent injection.

    This is the main function to call when running an agent.

    Args:
        include_events: Include recent NATS events
        include_state: Include current awareness state
        include_deadlines: Include upcoming deadlines
        event_count: Number of recent events to include
        deadline_hours: Deadline look-ahead window

    Returns:
        Formatted context string for agent prompt

    Example:
        >>> context = await get_agent_context()
        >>> response = await agent.arun(
        ...     message,
        ...     additional_context=context,
        ... )
    """
    parts = []
    now = datetime.utcnow()
    parts.append(f"# Real-Time Awareness (as of {now.strftime('%Y-%m-%d %H:%M UTC')})\n")

    # State summary (human-readable)
    if include_state:
        summary = await get_state_summary()
        if summary:
            parts.append(f"## Current State\n{summary}\n")
        else:
            # Fall back to KV state
            state = await get_current_state()
            if state:
                state_lines = [
                    f"- Awareness: {state.get('awareness_level', 'PRESENT')}",
                ]
                if state.get("focus_task"):
                    state_lines.append(f"- Focus: {state['focus_task']} ({state.get('focus_type', 'work')})")
                parts.append(f"## Current State\n" + "\n".join(state_lines) + "\n")

    # Upcoming deadlines
    if include_deadlines:
        deadlines = await get_upcoming_deadlines(hours=deadline_hours)
        if deadlines:
            deadline_lines = []
            for d in deadlines[:5]:  # Top 5 deadlines
                urgency_emoji = {
                    "critical": "🔴",
                    "high": "🟠",
                    "medium": "🟡",
                    "low": "🟢",
                }.get(d["urgency"], "⚪")
                deadline_lines.append(f"- {urgency_emoji} {d['title']} (due: {d['due']})")
            parts.append(f"## Upcoming Deadlines\n" + "\n".join(deadline_lines) + "\n")

    # Recent events
    if include_events:
        events = await get_recent_events(count=event_count)
        if events:
            event_lines = []
            for e in events[:event_count]:
                priority_star = "⭐" if e["priority"] >= 3 else ""
                event_lines.append(f"- {priority_star}{e['summary']} ({e['type']})")
            parts.append(f"## Recent Events\n" + "\n".join(event_lines) + "\n")

    return "\n".join(parts) if len(parts) > 1 else ""


async def get_agent_context_dict(
    include_events: bool = True,
    include_state: bool = True,
    include_deadlines: bool = True,
) -> Dict[str, Any]:
    """
    Get context as a structured dict (for tool use or metadata).

    Returns:
        Dict with events, state, deadlines keys
    """
    result = {}

    if include_state:
        result["state"] = await get_current_state()
        result["state_summary"] = await get_state_summary()

    if include_deadlines:
        result["deadlines"] = await get_upcoming_deadlines()

    if include_events:
        result["recent_events"] = await get_recent_events()

    return result


# Convenience for backwards compatibility
build_agent_context = get_agent_context


__all__ = [
    "get_agent_context",
    "get_agent_context_dict",
    "get_recent_events",
    "get_current_state",
    "get_state_summary",
    "get_upcoming_deadlines",
    "build_agent_context",
]
