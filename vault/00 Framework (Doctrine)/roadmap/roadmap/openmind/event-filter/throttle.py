"""
Throttle Engine
===============

Rate limits notifications per subscriber within time windows.
Prevents notification overload while ensuring critical alerts get through.

Features:
- Fixed window: Predefined time frame
- Dynamic window: Time frame from notification data
- Priority-aware: Different limits per priority level
- Group throttling: Rate limit by custom keys (e.g., per project)
"""

import logging
from datetime import datetime, timedelta
from typing import Optional, Dict

from backend.notifications.models import (
    Notification,
    NotificationPriority,
    ThrottleConfig,
)

logger = logging.getLogger("notifications.throttle")


class ThrottleEngine:
    """
    Rate limits notifications using Redis.

    Uses sliding window algorithm for smooth rate limiting.
    Different limits per priority level - critical is rarely throttled.
    """

    def __init__(self, redis_client):
        self.redis = redis_client

    # -------------------------------------------------------------------------
    # CHECK THROTTLE
    # -------------------------------------------------------------------------

    async def check(
        self,
        notification: Notification,
        config: ThrottleConfig,
    ) -> bool:
        """
        Check if notification should be throttled.

        Returns True if allowed, False if throttled.
        """
        if not config.enabled:
            return True  # Throttling disabled, allow

        # Critical notifications are rarely throttled
        if notification.priority == NotificationPriority.CRITICAL:
            max_count = config.critical_max
        elif notification.priority == NotificationPriority.HIGH:
            max_count = config.high_max
        elif notification.priority == NotificationPriority.NORMAL:
            max_count = config.normal_max
        else:
            max_count = config.low_max

        # Build throttle key
        key = self._build_key(notification)

        # Increment and check
        current = await self.redis.incr(key)

        if current == 1:
            # First request in window, set expiry
            await self.redis.expire(key, config.window_seconds)

        if current > max_count:
            logger.debug(
                f"Notification throttled: {notification.id} "
                f"(key={key}, count={current}, max={max_count})"
            )
            return False

        return True

    async def check_custom(
        self,
        user_id: str,
        throttle_key: str,
        window_seconds: int = 3600,
        max_count: int = 1,
    ) -> bool:
        """
        Check throttle with custom key and limits.

        Useful for one-off throttling like "max 1 per hour per project".
        """
        key = f"throttle:{user_id}:{throttle_key}"

        current = await self.redis.incr(key)
        if current == 1:
            await self.redis.expire(key, window_seconds)

        return current <= max_count

    # -------------------------------------------------------------------------
    # STATS
    # -------------------------------------------------------------------------

    async def get_count(
        self,
        user_id: str,
        throttle_key: Optional[str] = None,
    ) -> int:
        """Get current count in throttle window."""
        if throttle_key:
            key = f"throttle:{user_id}:{throttle_key}"
        else:
            key = f"throttle:{user_id}"

        count = await self.redis.get(key)
        return int(count) if count else 0

    async def get_remaining(
        self,
        user_id: str,
        config: ThrottleConfig,
        priority: NotificationPriority = NotificationPriority.NORMAL,
    ) -> int:
        """Get remaining notifications allowed in current window."""
        key = f"throttle:{user_id}"
        current = await self.redis.get(key)
        current = int(current) if current else 0

        if priority == NotificationPriority.CRITICAL:
            max_count = config.critical_max
        elif priority == NotificationPriority.HIGH:
            max_count = config.high_max
        elif priority == NotificationPriority.NORMAL:
            max_count = config.normal_max
        else:
            max_count = config.low_max

        return max(0, max_count - current)

    async def get_ttl(self, user_id: str) -> int:
        """Get seconds until throttle window resets."""
        key = f"throttle:{user_id}"
        ttl = await self.redis.ttl(key)
        return max(0, ttl)

    async def reset(self, user_id: str, throttle_key: Optional[str] = None):
        """Reset throttle for a user (admin action)."""
        if throttle_key:
            key = f"throttle:{user_id}:{throttle_key}"
        else:
            # Reset all throttle keys for user
            pattern = f"throttle:{user_id}*"
            async for key in self.redis.scan_iter(pattern):
                await self.redis.delete(key)
            return

        await self.redis.delete(key)
        logger.info(f"Throttle reset for {user_id}")

    # -------------------------------------------------------------------------
    # INTERNAL HELPERS
    # -------------------------------------------------------------------------

    def _build_key(self, notification: Notification) -> str:
        """Build Redis key for throttling."""
        base = f"throttle:{notification.user_id}"

        # Use custom throttle key if provided
        if notification.throttle_key:
            return f"{base}:{notification.throttle_key}"

        return base


# =============================================================================
# FACTORY
# =============================================================================

_throttle_engine: Optional[ThrottleEngine] = None


def get_throttle_engine() -> Optional[ThrottleEngine]:
    """Get the global throttle engine."""
    return _throttle_engine


def create_throttle_engine(redis_client) -> ThrottleEngine:
    """Create the throttle engine."""
    global _throttle_engine
    _throttle_engine = ThrottleEngine(redis_client)
    return _throttle_engine
