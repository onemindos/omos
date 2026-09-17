"""
Legacy Consciousness - Lightweight Event Processor

THE GATEKEEPER - Filters and enriches events before Legacy sees them.

Protects Legacy from noise while ensuring important events get through.

Features:
1. Noise filtering - Drop typing indicators, heartbeats, etc.
2. Awareness respect - Don't wake Legacy if dormant/resting
3. Event batching - Throttle high-volume events
4. Context enrichment - Add human-readable descriptions
5. Priority routing - Urgent events bypass filters
6. Correlation - Link related events across services
7. Time-aware filtering - Different awareness by time of day
8. Location-aware filtering - Per-site awareness levels
9. Configurable rules - UI can modify filter rules via API

Architecture:
- Subscribes to: events.> (all incoming events)
- Publishes to: legacy.inbox (filtered events for Legacy)
- Uses NATS KV Store for awareness state (no Redis needed!)

Integration Modes:
- Standalone: Creates own NATS connection and KV store
- Integrated: Uses existing KV store from Agno (preferred)
"""

import asyncio
import json
import logging
import os
import time
from datetime import datetime
from enum import IntEnum
from typing import Optional

from backend.consciousness.core.bus import Event, EventBus, create_event_bus
from backend.consciousness.core.kv_store import (
    AwarenessLevel,
    KVStore,
    FilterConfig,
    create_kv_store,
    get_kv_store,
)
from backend.consciousness.processing.correlator import (
    CorrelationEngine,
    create_correlator,
    get_correlator,
)

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("consciousness.processor")

# Default filter config (used as fallback, actual config comes from KV store)
DEFAULT_FILTER_CONFIG = FilterConfig()


class LightweightProcessor:
    """
    Lightweight event processor - the gatekeeper for Legacy.

    Filters noise, respects awareness levels, and enriches events
    before routing them to Legacy's inbox.

    Features:
    - Time-aware filtering (awareness varies by time of day)
    - Location-aware filtering (per-site awareness levels)
    - Configurable filter rules (UI can modify via API)
    - Uses NATS KV Store for state - no Redis dependency!

    Integration Modes:
    - Standalone: Creates own NATS/KV (for separate container)
    - Integrated: Uses existing KV store from Agno (preferred)
    """

    def __init__(
        self,
        nats_url: str = "nats://localhost:4222",
        kv_store: Optional[KVStore] = None,
        event_bus: Optional[EventBus] = None,
    ):
        self.nats_url = nats_url
        self._integrated_mode = kv_store is not None

        self.event_bus: Optional[EventBus] = event_bus
        self.kv_store: Optional[KVStore] = kv_store
        self.correlator: Optional[CorrelationEngine] = None

        # Cached filter config (updated via KV Watch - real-time, no polling)
        self._filter_config: FilterConfig = DEFAULT_FILTER_CONFIG
        self._filter_config_last_refresh: float = 0

        # In-memory batch buffers
        self._batch_buffers: dict[str, list[Event]] = {}
        self._batch_timers: dict[str, float] = {}

        # Stats
        self._stats = {
            "received": 0,
            "filtered_noise": 0,
            "filtered_awareness": 0,
            "filtered_time": 0,
            "filtered_location": 0,
            "batched": 0,
            "forwarded": 0,
            "correlated": 0,
        }

        # Running state
        self._running = False
        self._tasks: list[asyncio.Task] = []

    async def start(self, listen: bool = True):
        """
        Start the processor.

        Args:
            listen: If True, blocks and listens for events (standalone mode).
                    If False, just sets up subscriptions (integrated mode).
        """
        logger.info("Starting Consciousness Processor...")

        # Connect to NATS if not in integrated mode
        if not self._integrated_mode:
            logger.info(f"Connecting to NATS: {self.nats_url}")
            self.event_bus = await create_event_bus(self.nats_url)
            logger.info("Connected to NATS")

            logger.info("Initializing NATS KV Store...")
            self.kv_store = await create_kv_store(self.event_bus.js)
            logger.info("KV Store ready")
        else:
            logger.info("Running in integrated mode (using existing KV store)")

        # Initialize correlator
        logger.info("Initializing correlation engine...")
        self.correlator = create_correlator(self.kv_store)
        logger.info("Correlator ready")

        # Load initial filter config
        await self._refresh_filter_config()

        # Subscribe to all events
        await self.event_bus.subscribe("events.>", self.on_event)
        logger.info("Subscribed to events.> (all events)")

        self._running = True

        # Start batch flush task
        flush_task = asyncio.create_task(self._batch_flush_loop())
        self._tasks.append(flush_task)

        # Start real-time config watch (replaces polling loop)
        if self.kv_store:
            await self.kv_store.watch_config(self._on_config_change)
            logger.info("Config watch started (real-time updates via KV Watch)")

        logger.info("Processor ready - filtering events for Legacy...")

        # Block and listen (standalone mode) or return (integrated mode)
        if listen:
            await self.event_bus.listen()

    async def _refresh_filter_config(self):
        """Refresh filter config from KV store (called on startup)."""
        try:
            if self.kv_store:
                self._filter_config = await self.kv_store.get_filter_config()
                self._filter_config_last_refresh = time.time()
                logger.debug("Filter config refreshed")
        except Exception as e:
            logger.warning(f"Failed to refresh filter config: {e}")

    def _on_config_change(self, key: str, value: any) -> None:
        """
        Handle config changes from KV Watch (real-time updates).

        This replaces the 60-second polling loop with instant updates.
        When any config key changes, this callback is triggered immediately.
        """
        try:
            if key == "filter_config" and value:
                self._filter_config = FilterConfig(**value) if isinstance(value, dict) else self._filter_config
                self._filter_config_last_refresh = time.time()
                logger.info(f"Filter config updated in real-time: {key}")
            else:
                logger.debug(f"Config changed: {key}")
        except Exception as e:
            logger.warning(f"Failed to apply config change {key}: {e}")

    async def on_event(self, event: Event):
        """
        Process incoming event through the filter pipeline.

        Pipeline:
        1. Noise filter (drop typing, heartbeats)
        2. Critical check (bypass for emergencies)
        3. Time-aware filtering (check schedule)
        4. Location-aware filtering (per-site awareness)
        5. Awareness filter (respect Legacy's state)
        6. Batch check (group high-volume events)
        7. Enrich and forward
        """
        self._stats["received"] += 1
        event_type = str(event.type)

        # Extract location from event if present
        event_location = event.data.get("location") or event.data.get("location_id")

        # === STEP 1: Noise Filter ===
        if self._is_noise(event_type):
            self._stats["filtered_noise"] += 1
            logger.debug(f"Filtered noise: {event_type}")
            return

        # === STEP 2: Critical Events Bypass ===
        if self._is_critical(event_type):
            logger.warning(f"CRITICAL event - bypassing filters: {event_type}")
            await self._forward_event(event, bypass=True)
            return

        # === STEP 3: Get Effective Awareness (time + location aware) ===
        awareness = await self._get_effective_awareness(event_location)
        min_priority = self._filter_config.priority_thresholds.get(awareness.value, 1)

        if event.priority < min_priority:
            self._stats["filtered_awareness"] += 1
            logger.debug(f"Filtered by awareness ({awareness.name}): {event_type} (priority {event.priority})")
            return

        # === STEP 4: Batch Check ===
        batch_window = self._filter_config.batch_windows.get(event_type)
        if batch_window and awareness <= AwarenessLevel.AWARE:
            await self._add_to_batch(event, batch_window)
            return

        # === STEP 5: Forward Event ===
        await self._forward_event(event)

    def _is_noise(self, event_type: str) -> bool:
        """Check if event is noise that should always be filtered."""
        # Check exact matches
        if event_type in self._filter_config.noise_events:
            return True
        # Check patterns
        return any(
            pattern in event_type for pattern in self._filter_config.noise_patterns
        )

    def _is_critical(self, event_type: str) -> bool:
        """Check if event is critical and should bypass all filters."""
        # Check exact matches
        if event_type in self._filter_config.critical_events:
            return True
        # Check patterns
        return any(
            pattern in event_type for pattern in self._filter_config.critical_patterns
        )

    async def _get_effective_awareness(
        self,
        location_id: Optional[str] = None,
    ) -> AwarenessLevel:
        """
        Get effective awareness level considering time and location.

        Priority: Location override > Time schedule > Manual setting
        """
        if self.kv_store:
            return await self.kv_store.get_effective_awareness(
                location_id=location_id,
                at_time=datetime.utcnow(),
            )
        return AwarenessLevel.AWARE

    async def _add_to_batch(self, event: Event, window: int):
        """Add event to batch buffer for later processing."""
        event_type = str(event.type)

        if event_type not in self._batch_buffers:
            self._batch_buffers[event_type] = []
            self._batch_timers[event_type] = time.time()

        self._batch_buffers[event_type].append(event)
        self._stats["batched"] += 1
        logger.debug(f"Batched: {event_type} ({len(self._batch_buffers[event_type])} in buffer)")

    async def _batch_flush_loop(self):
        """Periodically flush batch buffers."""
        while True:
            await asyncio.sleep(10)  # Check every 10 seconds
            await self._flush_batches()

    async def _flush_batches(self):
        """Flush any batches that have exceeded their time window."""
        now = time.time()
        to_flush = []

        for event_type, timer in list(self._batch_timers.items()):
            window = BATCHABLE_EVENTS.get(event_type, 60)
            if now - timer >= window:
                to_flush.append(event_type)

        for event_type in to_flush:
            events = self._batch_buffers.pop(event_type, [])
            self._batch_timers.pop(event_type, None)

            if events:
                await self._forward_batched_events(event_type, events)

    async def _forward_batched_events(self, event_type: str, events: list[Event]):
        """Forward a batch of events as a summary."""
        if not events:
            return

        # Create a summary event
        count = len(events)
        latest = events[-1]

        # Build summary based on event type
        if "github.commit" in event_type:
            summary = f"{count} commits pushed"
            data = {
                "count": count,
                "repos": list(set(e.data.get("repo", "unknown") for e in events)),
                "latest": latest.data,
            }
        elif "todoist.task" in event_type:
            summary = f"{count} task updates"
            data = {
                "count": count,
                "tasks": [e.data.get("title", "unknown") for e in events[:5]],
            }
        elif "home.motion" in event_type:
            summary = f"Motion detected in {count} areas"
            data = {
                "count": count,
                "areas": list(set(e.data.get("friendly_name", "unknown") for e in events)),
            }
        else:
            summary = f"{count} {event_type} events"
            data = {"count": count}

        # Create summary event
        summary_event = Event(
            type=f"{event_type}.batch",
            source="processor",
            data={
                "summary": summary,
                "batch_count": count,
                **data,
            },
            priority=max(e.priority for e in events),  # Highest priority from batch
        )

        await self._forward_event(summary_event)
        logger.info(f"Flushed batch: {summary}")

    async def _forward_event(self, event: Event, bypass: bool = False):
        """Forward event to Legacy's inbox with enriched context."""
        # Enrich with human-readable context
        context = self._build_context(event)

        # Add context to event data
        enriched_data = {
            **event.data,
            "_context": context,
            "_processed_at": datetime.utcnow().isoformat(),
            "_bypass_filters": bypass,
        }

        # Create enriched event for Legacy
        enriched_event = Event(
            type=event.type,
            source=event.source,
            data=enriched_data,
            timestamp=event.timestamp,
            priority=event.priority,
        )

        # Publish to Legacy's inbox
        try:
            await self.event_bus._js.publish(
                subject="legacy.inbox",
                payload=enriched_event.to_bytes(),
                headers={
                    "Priority": str(event.priority),
                    "Source": event.source,
                    "Context": context[:100],  # First 100 chars
                },
            )
            self._stats["forwarded"] += 1
            logger.info(f"Forwarded to Legacy: {event.type} - {context[:50]}...")

        except Exception as e:
            logger.error(f"Failed to forward event: {e}")

    def _build_context(self, event: Event) -> str:
        """Build human-readable context for the event."""
        event_type = str(event.type)
        data = event.data

        # Task events
        if "todoist" in event_type:
            if "created" in event_type:
                due = f" due {data.get('due')}" if data.get("due") else ""
                return f"New task: '{data.get('title', 'Unknown')}'{due}"
            elif "completed" in event_type:
                return f"Completed: '{data.get('title', 'Unknown')}'"
            elif "updated" in event_type:
                return f"Updated: '{data.get('title', 'Unknown')}'"
            elif "batch" in event_type:
                return data.get("summary", f"{data.get('count', 0)} task events")

        # GitHub events
        elif "github" in event_type:
            if "pr" in event_type:
                action = data.get("action", "updated")
                return f"PR #{data.get('pr_number', '?')} {action}: '{data.get('title', 'Unknown')}'"
            elif "issue" in event_type:
                action = data.get("action", "updated")
                return f"Issue #{data.get('issue_number', '?')} {action}"
            elif "commit" in event_type or "push" in event_type:
                if "batch" in event_type:
                    return data.get("summary", f"{data.get('count', 0)} commits")
                return f"Push to {data.get('repo', 'unknown')}: {data.get('commits', 1)} commits"

        # Home events
        elif "home" in event_type:
            friendly = data.get("friendly_name", data.get("entity_id", "Unknown"))
            state = data.get("state", "unknown")

            if "door" in event_type:
                return f"{friendly} {'opened' if state == 'on' else 'closed'}"
            elif "motion" in event_type:
                if "batch" in event_type:
                    return data.get("summary", "Motion detected")
                return f"Motion: {friendly}"
            elif "alarm" in event_type:
                return f"ALARM: {friendly} - {state}"
            else:
                return f"{friendly}: {state}"

        # Calendar events
        elif "calendar" in event_type:
            return f"Calendar: {data.get('title', 'Event')} at {data.get('start', 'unknown time')}"

        # Security events
        elif "security" in event_type:
            return f"SECURITY: {data.get('description', event_type)}"

        # Batch events
        elif "batch" in event_type:
            return data.get("summary", f"Batch of {data.get('count', 0)} events")

        # Default
        return f"{event_type}: {json.dumps(data)[:80]}"

    async def get_stats(self) -> dict:
        """Get processor statistics."""
        effective_awareness = await self._get_effective_awareness()
        base_awareness = AwarenessLevel.AWARE
        if self.kv_store:
            base_awareness = await self.kv_store.get_awareness_level()

        return {
            **self._stats,
            "batch_buffers": {k: len(v) for k, v in self._batch_buffers.items()},
            "base_awareness": base_awareness.name,
            "effective_awareness": effective_awareness.name,
            "filter_config_age": time.time() - self._filter_config_last_refresh,
            "integrated_mode": self._integrated_mode,
        }

    async def stop(self):
        """Stop the processor."""
        self._running = False

        # Cancel background tasks
        for task in self._tasks:
            task.cancel()
            try:
                await task
            except asyncio.CancelledError:
                pass
        self._tasks.clear()

        # Flush any remaining batches
        await self._flush_batches()

        # Only close resources if not in integrated mode
        if not self._integrated_mode:
            if self.kv_store:
                await self.kv_store.close()
            if self.event_bus:
                await self.event_bus.disconnect()

        logger.info(f"Processor stopped. Stats: {self._stats}")


# =============================================================================
# Global Processor Instance
# =============================================================================

_processor: Optional[LightweightProcessor] = None


async def create_processor(
    nats_url: str = "nats://localhost:4222",
    kv_store: Optional[KVStore] = None,
    event_bus: Optional[EventBus] = None,
) -> LightweightProcessor:
    """
    Create and start the global processor instance.

    Args:
        nats_url: NATS URL (only used in standalone mode)
        kv_store: Existing KV store (integrated mode)
        event_bus: Existing event bus (integrated mode)

    Returns:
        The processor instance
    """
    global _processor
    _processor = LightweightProcessor(
        nats_url=nats_url,
        kv_store=kv_store,
        event_bus=event_bus,
    )
    # Start in non-blocking mode for integrated use
    await _processor.start(listen=False)
    return _processor


def get_processor() -> Optional[LightweightProcessor]:
    """Get the global processor instance."""
    return _processor


async def stop_processor() -> None:
    """Stop the global processor."""
    global _processor
    if _processor:
        await _processor.stop()
        _processor = None


async def main():
    """Main entry point for standalone mode."""
    processor = LightweightProcessor(
        nats_url=os.environ.get("NATS_URL", "nats://nats:4222"),
    )

    try:
        await processor.start(listen=True)
    except KeyboardInterrupt:
        await processor.stop()


if __name__ == "__main__":
    asyncio.run(main())


__all__ = [
    "LightweightProcessor",
    "create_processor",
    "get_processor",
    "stop_processor",
]
