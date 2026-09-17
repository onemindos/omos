"""
Temporal Context Injection
==========================

Middleware and utilities for injecting temporal context into:
- Agent system prompts
- Tool contexts
- API responses

This is what makes time "permanent" in the system - every component
automatically has access to rich time awareness.
"""

import logging
from datetime import datetime
from typing import Optional, Dict, Any, Callable
from functools import wraps
from zoneinfo import ZoneInfo

from .models import TemporalContext, TimeOfDay, DayType
from .service import get_temporal_service

logger = logging.getLogger("consciousness.temporal.context")


# =============================================================================
# CONTEXT GENERATOR
# =============================================================================


class TemporalContextGenerator:
    """
    Generates temporal context for injection into various components.

    This provides a simplified interface to the TemporalService
    specifically for context generation.
    """

    def __init__(
        self,
        default_timezone: str = "America/Los_Angeles",
    ):
        """
        Initialize the context generator.

        Args:
            default_timezone: Default timezone for context
        """
        self.default_timezone = default_timezone
        self._timezone = ZoneInfo(default_timezone)

    async def get_context(
        self,
        user_id: str = "zeus",
        timezone: Optional[str] = None,
    ) -> TemporalContext:
        """
        Get the full temporal context.

        Args:
            user_id: User to generate context for
            timezone: Override timezone

        Returns:
            Complete TemporalContext
        """
        service = get_temporal_service()
        return await service.get_context(user_id, timezone)

    def get_basic_context(
        self,
        timezone: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Get basic temporal context (no async, no external data).

        Useful for quick context without database lookups.

        Args:
            timezone: Override timezone

        Returns:
            Dict with basic time info
        """
        tz = ZoneInfo(timezone) if timezone else self._timezone
        now = datetime.now(tz)

        return {
            "now": now.isoformat(),
            "timezone": timezone or self.default_timezone,
            "time_of_day": self._classify_time_of_day(now).value,
            "day_type": self._classify_day_type(now).value,
            "hour": now.hour,
            "day_of_week": now.weekday(),
            "is_weekend": now.weekday() >= 5,
            "is_work_hours": 9 <= now.hour < 17 and now.weekday() < 5,
        }

    def to_prompt_string(
        self,
        context: TemporalContext,
    ) -> str:
        """
        Convert context to a string for system prompts.

        Args:
            context: The temporal context

        Returns:
            Formatted string for prompt injection
        """
        return context.to_prompt_context()

    def to_brief_string(
        self,
        timezone: Optional[str] = None,
    ) -> str:
        """
        Get a brief time string for quick context.

        Args:
            timezone: Override timezone

        Returns:
            Brief time description
        """
        basic = self.get_basic_context(timezone)
        return (
            f"{basic['time_of_day'].replace('_', ' ').title()} "
            f"({basic['day_type']})"
        )

    def _classify_time_of_day(self, dt: datetime) -> TimeOfDay:
        """Classify time of day."""
        hour = dt.hour
        if 5 <= hour < 7:
            return TimeOfDay.EARLY_MORNING
        elif 7 <= hour < 12:
            return TimeOfDay.MORNING
        elif 12 <= hour < 17:
            return TimeOfDay.AFTERNOON
        elif 17 <= hour < 21:
            return TimeOfDay.EVENING
        elif 21 <= hour <= 23:
            return TimeOfDay.NIGHT
        else:
            return TimeOfDay.LATE_NIGHT

    def _classify_day_type(self, dt: datetime) -> DayType:
        """Classify day type."""
        if dt.weekday() >= 5:
            return DayType.WEEKEND
        return DayType.WORKDAY


# =============================================================================
# PROMPT TEMPLATES
# =============================================================================


TEMPORAL_SYSTEM_PROMPT_TEMPLATE = """
## Current Time Context

{context}

Use this temporal awareness to:
- Prioritize time-sensitive tasks
- Respect focus windows for deep work
- Consider deadline pressure when planning
- Adapt communication style to time of day
- Suggest optimal timing for actions
"""


MINIMAL_TIME_PROMPT = """
Current time: {now} ({timezone})
Time of day: {time_of_day}
Day type: {day_type}
"""


# =============================================================================
# INJECTION UTILITIES
# =============================================================================


def inject_temporal_context(
    system_prompt: str,
    context: TemporalContext,
    template: str = TEMPORAL_SYSTEM_PROMPT_TEMPLATE,
) -> str:
    """
    Inject temporal context into a system prompt.

    Args:
        system_prompt: Original system prompt
        context: Temporal context to inject
        template: Template to use for injection

    Returns:
        System prompt with temporal context appended
    """
    temporal_section = template.format(
        context=context.to_prompt_context()
    )

    return f"{system_prompt}\n\n{temporal_section}"


def inject_minimal_time(
    system_prompt: str,
    timezone: str = "America/Los_Angeles",
) -> str:
    """
    Inject minimal time context (no async).

    Args:
        system_prompt: Original system prompt
        timezone: Timezone to use

    Returns:
        System prompt with time info appended
    """
    generator = TemporalContextGenerator(timezone)
    basic = generator.get_basic_context()

    time_section = MINIMAL_TIME_PROMPT.format(
        now=basic["now"],
        timezone=basic["timezone"],
        time_of_day=basic["time_of_day"],
        day_type=basic["day_type"],
    )

    return f"{system_prompt}\n\n{time_section}"


# =============================================================================
# DECORATOR FOR AUTOMATIC CONTEXT
# =============================================================================


def with_temporal_context(
    user_id: str = "zeus",
    timezone: Optional[str] = None,
    context_key: str = "temporal_context",
):
    """
    Decorator to automatically inject temporal context.

    Adds temporal context to the function's kwargs.

    Usage:
        @with_temporal_context()
        async def my_agent_function(prompt: str, temporal_context: TemporalContext):
            # temporal_context is automatically provided
            pass

    Args:
        user_id: User to get context for
        timezone: Override timezone
        context_key: Key to use in kwargs
    """
    def decorator(func: Callable):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            # Generate context
            service = get_temporal_service()
            context = await service.get_context(user_id, timezone)

            # Inject into kwargs
            kwargs[context_key] = context

            return await func(*args, **kwargs)

        return wrapper
    return decorator


# =============================================================================
# CONTEXT ENRICHER
# =============================================================================


class ContextEnricher:
    """
    Enriches various data structures with temporal context.

    Use this to add time awareness to:
    - Agent run parameters
    - Tool call contexts
    - API responses
    """

    def __init__(
        self,
        default_timezone: str = "America/Los_Angeles",
    ):
        self.generator = TemporalContextGenerator(default_timezone)

    async def enrich_agent_params(
        self,
        params: Dict[str, Any],
        user_id: str = "zeus",
    ) -> Dict[str, Any]:
        """
        Enrich agent run parameters with temporal context.

        Args:
            params: Original parameters
            user_id: User ID

        Returns:
            Parameters with temporal context added
        """
        context = await self.generator.get_context(user_id)

        # Add to context injection if present
        if "context" in params and isinstance(params["context"], dict):
            params["context"]["temporal"] = {
                "now": context.now.isoformat(),
                "time_of_day": context.time_of_day.value,
                "day_type": context.day_type.value,
                "deadline_pressure": context.deadline_pressure,
                "is_focus_time": context.is_focus_time,
                "productivity_score": context.productivity_score,
            }
        else:
            params["temporal_context"] = context

        return params

    def enrich_tool_result(
        self,
        result: Dict[str, Any],
    ) -> Dict[str, Any]:
        """
        Enrich a tool result with time metadata.

        Args:
            result: Original tool result

        Returns:
            Result with time metadata
        """
        basic = self.generator.get_basic_context()

        result["_temporal"] = {
            "executed_at": basic["now"],
            "time_of_day": basic["time_of_day"],
        }

        return result

    async def enrich_api_response(
        self,
        response: Dict[str, Any],
        user_id: str = "zeus",
        include_full_context: bool = False,
    ) -> Dict[str, Any]:
        """
        Enrich an API response with temporal info.

        Args:
            response: Original response
            user_id: User ID
            include_full_context: Include full context or just basics

        Returns:
            Response with temporal info
        """
        if include_full_context:
            context = await self.generator.get_context(user_id)
            response["temporal"] = context.model_dump()
        else:
            response["temporal"] = self.generator.get_basic_context()

        return response


# =============================================================================
# SINGLETON
# =============================================================================

_context_generator: Optional[TemporalContextGenerator] = None
_context_enricher: Optional[ContextEnricher] = None


def get_context_generator() -> TemporalContextGenerator:
    """Get the global context generator."""
    global _context_generator
    if _context_generator is None:
        _context_generator = TemporalContextGenerator()
    return _context_generator


def get_context_enricher() -> ContextEnricher:
    """Get the global context enricher."""
    global _context_enricher
    if _context_enricher is None:
        _context_enricher = ContextEnricher()
    return _context_enricher
