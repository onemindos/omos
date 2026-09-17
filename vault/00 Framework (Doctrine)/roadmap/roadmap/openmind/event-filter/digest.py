"""
Digest Engine
=============

Batches multiple notifications into single digest messages.
Uses Redis for time-windowed aggregation.

Features:
- Regular digest: Always batches for specified window
- Look-back digest: Checks if recent message sent, sends immediately if not
- Scheduled digest: Sends at specific times (daily briefing style)
- AI-aware: Considers awareness mode in batching decisions
"""

import asyncio
import json
import logging
import time as time_module
from datetime import datetime, timedelta
from typing import Optional, List, Dict, Any, Callable, Awaitable
from uuid import UUID, uuid4

from backend.notifications.models import (
    Notification,
    NotificationPriority,
    DigestConfig,
    DigestEvent,
    DigestSummary,
)

logger = logging.getLogger("notifications.digest")


class DigestEngine:
    """
    Batches notifications into digests using Redis.

    Strategies:
    1. Regular: Always digest for window duration
    2. Look-back: Send immediately if no recent notification, else digest
    3. Scheduled: Send at specific times (e.g., 9 AM daily)

    AI-Aware Enhancement:
    - Omnipresent mode: Bypass digest, send all immediately
    - Dormant mode: Aggressive digest, only critical bypasses
    """

    def __init__(
        self,
        redis_client,
        nats_client=None,
        on_flush: Optional[Callable[[DigestSummary], Awaitable[None]]] = None,
    ):
        self.redis = redis_client
        self.nats = nats_client
        self.on_flush = on_flush
        self._scheduler_task: Optional[asyncio.Task] = None
        self._running = False

    async def start(self):
        """Start the digest scheduler."""
        self._running = True
        self._scheduler_task = asyncio.create_task(self._run_scheduler())
        logger.info("Digest engine started")

    async def stop(self):
        """Stop the digest scheduler."""
        self._running = False
        if self._scheduler_task:
            self._scheduler_task.cancel()
            try:
                await self._scheduler_task
            except asyncio.CancelledError:
                pass
        logger.info("Digest engine stopped")

    # -------------------------------------------------------------------------
    # ADD EVENT TO DIGEST
    # -------------------------------------------------------------------------

    async def add_event(
        self,
        notification: Notification,
        config: DigestConfig,
        awareness_mode: str = "present",
    ) -> bool:
        """
        Add notification to digest window.

        Returns True if notification was digested, False if should send immediately.
        """
        # Critical always bypasses digest
        if notification.priority == NotificationPriority.CRITICAL:
            logger.debug(f"Critical notification bypasses digest: {notification.id}")
            return False

        # Omnipresent mode: send everything immediately
        if awareness_mode == "omnipresent":
            logger.debug(f"Omnipresent mode bypasses digest: {notification.id}")
            return False

        # Digest disabled
        if not config.enabled:
            return False

        # Build digest key
        digest_key = self._build_key(notification)

        # Add to sorted set with timestamp score
        event_data = DigestEvent(notification=notification).model_dump_json()
        score = time_module.time()

        await self.redis.zadd(digest_key, {event_data: score})
        await self.redis.expire(digest_key, config.window_seconds + 60)

        # Schedule flush if not already scheduled
        schedule_key = f"{digest_key}:scheduled"
        if not await self.redis.exists(schedule_key):
            await self.redis.setex(schedule_key, config.window_seconds, "1")

            # Publish flush job to NATS
            if self.nats:
                await self._schedule_flush(
                    digest_key,
                    notification.user_id,
                    config.window_seconds,
                )

            logger.debug(f"Digest scheduled for {config.window_seconds}s: {digest_key}")

        current_count = await self.redis.zcard(digest_key)
        logger.debug(f"Event added to digest {digest_key}, count: {current_count}")

        return True

    async def add_event_lookback(
        self,
        notification: Notification,
        config: DigestConfig,
        lookback_seconds: int = 60,
    ) -> bool:
        """
        Look-back digest strategy.

        If a notification was sent recently, add to digest.
        If not, send immediately (first notification goes through).
        """
        # Check if recent notification was sent
        lookback_key = f"digest:lookback:{notification.user_id}"
        recent = await self.redis.get(lookback_key)

        if recent:
            # Recent notification sent, digest this one
            return await self.add_event(notification, config)
        else:
            # No recent notification, send immediately
            # Mark that we're sending
            await self.redis.setex(lookback_key, lookback_seconds, "1")
            return False

    # -------------------------------------------------------------------------
    # FLUSH DIGEST
    # -------------------------------------------------------------------------

    async def flush(
        self,
        digest_key: str,
        user_id: str,
    ) -> Optional[DigestSummary]:
        """
        Flush a digest and return summary.

        Called by scheduler when window expires.
        """
        # Get all events in window
        events_raw = await self.redis.zrange(digest_key, 0, -1)
        if not events_raw:
            return None

        # Parse events
        events = []
        for raw in events_raw:
            try:
                event = DigestEvent.model_validate_json(raw)
                events.append(event)
            except Exception as e:
                logger.warning(f"Failed to parse digest event: {e}")

        if not events:
            return None

        # Get window timestamps
        scores = await self.redis.zrange(digest_key, 0, -1, withscores=True)
        timestamps = [s for _, s in scores] if scores else []

        window_start = datetime.fromtimestamp(min(timestamps)) if timestamps else datetime.utcnow()
        window_end = datetime.utcnow()

        # Build summary
        summary = DigestSummary(
            user_id=user_id,
            digest_key=digest_key,
            events=events,
            event_count=len(events),
            window_start=window_start,
            window_end=window_end,
        )

        # Clear the digest
        await self.redis.delete(digest_key)
        await self.redis.delete(f"{digest_key}:scheduled")

        logger.info(f"Flushed digest {digest_key}: {len(events)} events")

        # Call flush callback
        if self.on_flush:
            try:
                await self.on_flush(summary)
            except Exception as e:
                logger.error(f"Digest flush callback failed: {e}")

        return summary

    async def flush_all_for_user(self, user_id: str) -> List[DigestSummary]:
        """Flush all pending digests for a user."""
        pattern = f"digest:{user_id}:*"
        keys = []

        async for key in self.redis.scan_iter(pattern):
            if not key.endswith(":scheduled"):
                keys.append(key)

        summaries = []
        for key in keys:
            summary = await self.flush(key, user_id)
            if summary:
                summaries.append(summary)

        return summaries

    # -------------------------------------------------------------------------
    # STATS
    # -------------------------------------------------------------------------

    async def get_pending_count(self, user_id: str) -> int:
        """Get count of pending digest events for user."""
        pattern = f"digest:{user_id}:*"
        total = 0

        async for key in self.redis.scan_iter(pattern):
            if not key.endswith(":scheduled"):
                count = await self.redis.zcard(key)
                total += count

        return total

    async def get_pending_digests(self, user_id: str) -> Dict[str, int]:
        """Get pending digest keys and event counts."""
        pattern = f"digest:{user_id}:*"
        digests = {}

        async for key in self.redis.scan_iter(pattern):
            if not key.endswith(":scheduled"):
                count = await self.redis.zcard(key)
                digests[key] = count

        return digests

    # -------------------------------------------------------------------------
    # INTERNAL HELPERS
    # -------------------------------------------------------------------------

    def _build_key(self, notification: Notification) -> str:
        """Build Redis key for digest grouping."""
        base = f"digest:{notification.user_id}"

        # Use custom digest key if provided
        if notification.digest_key:
            return f"{base}:{notification.digest_key}"

        # Default grouping by category
        category = notification.category or "general"
        return f"{base}:{category}"

    async def _schedule_flush(
        self,
        digest_key: str,
        user_id: str,
        delay_seconds: int,
    ):
        """Schedule a flush job via NATS."""
        if self.nats:
            job = {
                "type": "digest.flush",
                "digest_key": digest_key,
                "user_id": user_id,
                "flush_at": (datetime.utcnow() + timedelta(seconds=delay_seconds)).isoformat(),
            }
            await self.nats.publish("notifications.jobs", json.dumps(job).encode())

    async def _run_scheduler(self):
        """
        Background task that checks for expired digests.

        This is a fallback - primary scheduling is via NATS jobs.
        """
        while self._running:
            try:
                await asyncio.sleep(30)  # Check every 30 seconds

                # Find expired digests
                pattern = "digest:*:scheduled"
                async for key in self.redis.scan_iter(pattern):
                    # If scheduled key expired but digest still exists
                    digest_key = key.replace(":scheduled", "")
                    if await self.redis.exists(digest_key):
                        exists = await self.redis.exists(key)
                        if not exists:
                            # Scheduled key expired, flush digest
                            user_id = digest_key.split(":")[1]
                            await self.flush(digest_key, user_id)

            except asyncio.CancelledError:
                break
            except Exception as e:
                logger.error(f"Digest scheduler error: {e}")


# =============================================================================
# FACTORY
# =============================================================================

_digest_engine: Optional[DigestEngine] = None


def get_digest_engine() -> Optional[DigestEngine]:
    """Get the global digest engine."""
    return _digest_engine


async def create_digest_engine(
    redis_client,
    nats_client=None,
    on_flush: Optional[Callable[[DigestSummary], Awaitable[None]]] = None,
) -> DigestEngine:
    """Create and start the digest engine."""
    global _digest_engine
    _digest_engine = DigestEngine(redis_client, nats_client, on_flush)
    await _digest_engine.start()
    return _digest_engine
