"""
Awareness-Enhanced Agent Runner
===============================

Wraps Agno agent runs with real-time context from the Consciousness layer.

This ensures agents are "aware" of:
- Recent NATS events
- Current user state
- Upcoming deadlines
- Device presence

Usage:
    from backend.consciousness.aware_runner import run_with_awareness

    # Instead of:
    #   response = await agent.arun(message)
    # Use:
    response = await run_with_awareness(agent, message, session_id=session_id)
"""

import logging
from typing import Any, Dict, Optional

from agno.agent import Agent

from backend.consciousness.awareness_provider import get_agent_context

logger = logging.getLogger("agno.aware_runner")


async def run_with_awareness(
    agent: Agent,
    message: str,
    *,
    session_id: Optional[str] = None,
    user_id: Optional[str] = None,
    include_events: bool = True,
    include_state: bool = True,
    include_deadlines: bool = True,
    event_count: int = 10,
    additional_context: Optional[str] = None,
    **kwargs,
) -> Any:
    """
    Run an agent with real-time awareness context injected.

    This wraps the standard agent.arun() to automatically inject
    context from the Consciousness layer.

    Args:
        agent: The Agno agent to run
        message: User message
        session_id: Session ID for context
        user_id: User ID for personalization
        include_events: Include recent NATS events
        include_state: Include current state
        include_deadlines: Include upcoming deadlines
        event_count: Number of events to include
        additional_context: Extra context to append
        **kwargs: Additional args passed to agent.arun()

    Returns:
        Agent run response

    Example:
        >>> from backend.agno.agents.factory import create_agent_from_preset
        >>> agent = await create_agent_from_preset("legacy")
        >>> response = await run_with_awareness(agent, "What's on my calendar today?")
    """
    # Build awareness context
    awareness_context = await get_agent_context(
        include_events=include_events,
        include_state=include_state,
        include_deadlines=include_deadlines,
        event_count=event_count,
    )

    # Combine with any additional context
    full_context = awareness_context
    if additional_context:
        full_context = f"{awareness_context}\n\n{additional_context}"

    # Log for debugging
    if awareness_context:
        logger.debug(f"Injecting awareness context ({len(awareness_context)} chars)")

    # Run the agent with combined context
    # Note: We use additional_context at run time, which OVERRIDES the agent's default
    # The agent's preset additional_context is for static profile info
    # This run-time context is for real-time awareness
    response = await agent.arun(
        message=message,
        session_id=session_id,
        user_id=user_id,
        additional_context=full_context if full_context else None,
        **kwargs,
    )

    return response


async def stream_with_awareness(
    agent: Agent,
    message: str,
    *,
    session_id: Optional[str] = None,
    user_id: Optional[str] = None,
    include_events: bool = True,
    include_state: bool = True,
    include_deadlines: bool = True,
    event_count: int = 10,
    additional_context: Optional[str] = None,
    **kwargs,
):
    """
    Stream agent response with real-time awareness context.

    Same as run_with_awareness but returns an async generator for streaming.

    Yields:
        RunResponse chunks for streaming to client
    """
    # Build awareness context
    awareness_context = await get_agent_context(
        include_events=include_events,
        include_state=include_state,
        include_deadlines=include_deadlines,
        event_count=event_count,
    )

    # Combine with any additional context
    full_context = awareness_context
    if additional_context:
        full_context = f"{awareness_context}\n\n{additional_context}"

    # Stream the agent response
    async for chunk in agent.arun_stream(
        message=message,
        session_id=session_id,
        user_id=user_id,
        additional_context=full_context if full_context else None,
        **kwargs,
    ):
        yield chunk


class AwareAgent:
    """
    Wrapper class that makes any agent awareness-enabled.

    Usage:
        aware_legacy = AwareAgent(get_agent("legacy"))
        response = await aware_legacy.run("What's happening?")
    """

    def __init__(
        self,
        agent: Agent,
        include_events: bool = True,
        include_state: bool = True,
        include_deadlines: bool = True,
        event_count: int = 10,
    ):
        self.agent = agent
        self.include_events = include_events
        self.include_state = include_state
        self.include_deadlines = include_deadlines
        self.event_count = event_count

    async def run(
        self,
        message: str,
        session_id: Optional[str] = None,
        **kwargs,
    ) -> Any:
        """Run with awareness."""
        return await run_with_awareness(
            self.agent,
            message,
            session_id=session_id,
            include_events=self.include_events,
            include_state=self.include_state,
            include_deadlines=self.include_deadlines,
            event_count=self.event_count,
            **kwargs,
        )

    async def stream(
        self,
        message: str,
        session_id: Optional[str] = None,
        **kwargs,
    ):
        """Stream with awareness."""
        async for chunk in stream_with_awareness(
            self.agent,
            message,
            session_id=session_id,
            include_events=self.include_events,
            include_state=self.include_state,
            include_deadlines=self.include_deadlines,
            event_count=self.event_count,
            **kwargs,
        ):
            yield chunk


__all__ = [
    "run_with_awareness",
    "stream_with_awareness",
    "AwareAgent",
]
