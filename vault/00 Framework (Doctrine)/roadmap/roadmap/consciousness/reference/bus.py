"""
Legacy Consciousness - NATS JetStream Event Bus

THE NERVOUS SYSTEM - The living awareness stream.

This is what makes Legacy "alive" - the continuous flow of all experiences,
actions, and context in real-time. Without this stream, Legacy is just
another chatbot. With it, she's aware of everything happening around you.

NATS JetStream provides:
- Sub-millisecond latency (200-400us vs 1-5ms for Redis Streams)
- Push model (events come to you, no polling)
- Persistence (events survive restarts)
- Perfect for future IoT/robots (OpenMind OM1)

Event Types:
- todoist.* - Task events (created, completed, updated)
- github.* - Code events (PR, issue, commit)
- calendar.* - Schedule events (meeting soon, conflict)
- home.* - Home Assistant events (motion, door, etc.)
- system.* - Internal events (agent started, error, etc.)
- biometric.* - Neural interface data (HR, motion, sleep, stress)
- environment.* - IoT sensors (temp, humidity, air quality)
- presence.* - Location awareness (motion, occupancy, devices)
"""

import asyncio
import json
import logging
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from typing import Any, Callable, Optional, Union

import nats
from nats.js.api import (
    AckPolicy,
    ConsumerConfig,
    DeliverPolicy,
    RetentionPolicy,
    StorageType,
    StreamConfig,
)

logger = logging.getLogger("consciousness.bus")

# Type hints for new components (avoid circular imports)
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from backend.consciousness.core.kv_store import KVStore
    from backend.consciousness.core.object_store import MemoryStore
    from backend.consciousness.bridges.service import NATSService


class EventType(str, Enum):
    """Event types that flow through the bus."""

    # Todoist
    TASK_CREATED = "todoist.task.created"
    TASK_COMPLETED = "todoist.task.completed"
    TASK_UPDATED = "todoist.task.updated"
    TASK_DUE_SOON = "todoist.task.due_soon"

    # GitHub
    PR_OPENED = "github.pr.opened"
    PR_MERGED = "github.pr.merged"
    PR_REVIEW_REQUESTED = "github.pr.review_requested"
    ISSUE_CREATED = "github.issue.created"
    COMMIT_PUSHED = "github.commit.pushed"

    # Calendar
    MEETING_SOON = "calendar.meeting.soon"
    MEETING_STARTED = "calendar.meeting.started"
    SCHEDULE_CONFLICT = "calendar.schedule.conflict"

    # Home Assistant
    MOTION_DETECTED = "home.motion.detected"
    DOOR_OPENED = "home.door.opened"
    DEVICE_STATE_CHANGED = "home.device.changed"

    # System
    AGENT_STARTED = "system.agent.started"
    AGENT_ERROR = "system.agent.error"
    USER_CONNECTED = "system.user.connected"
    USER_DISCONNECTED = "system.user.disconnected"

    # =========================================================================
    # Biometric Events - Neural Interface Data
    # =========================================================================

    # Heart Rate / Cardiovascular
    BIOMETRIC_HEART = "biometric.heart"  # HR, HRV, SpO2 reading
    BIOMETRIC_HEART_ELEVATED = "biometric.heart.elevated"  # HR above normal
    BIOMETRIC_HEART_LOW = "biometric.heart.low"  # HR below normal
    BIOMETRIC_HRV_LOW = "biometric.hrv.low"  # Low heart rate variability
    BIOMETRIC_SPO2_LOW = "biometric.spo2.low"  # Blood oxygen below 90%

    # Motion / Activity
    BIOMETRIC_MOTION = "biometric.motion"  # Accelerometer reading
    BIOMETRIC_STEPS = "biometric.steps"  # Step count update
    BIOMETRIC_ACTIVITY_CHANGED = "biometric.activity.changed"  # Activity type changed

    # Sleep
    BIOMETRIC_SLEEP_STARTED = "biometric.sleep.started"  # User fell asleep
    BIOMETRIC_SLEEP_ENDED = "biometric.sleep.ended"  # User woke up
    BIOMETRIC_SLEEP_STAGE = "biometric.sleep.stage"  # Sleep stage changed

    # Temperature
    BIOMETRIC_TEMPERATURE = "biometric.temperature"  # Body temperature reading
    BIOMETRIC_TEMP_ELEVATED = "biometric.temp.elevated"  # Fever detection

    # Stress / Recovery
    BIOMETRIC_STRESS_HIGH = "biometric.stress.high"  # High stress detected
    BIOMETRIC_STRESS_LOW = "biometric.stress.low"  # Low stress / relaxed
    BIOMETRIC_RECOVERY_LOW = "biometric.recovery.low"  # Poor recovery score

    # =========================================================================
    # Biometric Alerts - Immediate Response Required
    # =========================================================================

    BIOMETRIC_ALERT_FALL = "biometric.alert.fall"  # Fall detected
    BIOMETRIC_ALERT_HR_CRITICAL = "biometric.alert.hr_critical"  # HR > 180 or < 40
    BIOMETRIC_ALERT_SPO2_CRITICAL = "biometric.alert.spo2_critical"  # SpO2 < 85%
    BIOMETRIC_ALERT_INACTIVE = "biometric.alert.inactive"  # No movement for too long

    # =========================================================================
    # Biometric Patterns - Derived Insights
    # =========================================================================

    BIOMETRIC_PATTERN_STRESS = "biometric.pattern.stress"  # Stress pattern detected
    BIOMETRIC_PATTERN_FATIGUE = "biometric.pattern.fatigue"  # Fatigue pattern detected
    BIOMETRIC_PATTERN_RECOVERY = "biometric.pattern.recovery"  # Recovery state
    BIOMETRIC_PATTERN_CIRCADIAN = "biometric.pattern.circadian"  # Sleep rhythm pattern

    # =========================================================================
    # Environment Events - IoT Sensors
    # =========================================================================

    ENVIRONMENT_TEMPERATURE = "environment.temperature"
    ENVIRONMENT_HUMIDITY = "environment.humidity"
    ENVIRONMENT_AIR_QUALITY = "environment.air_quality"
    ENVIRONMENT_LIGHT = "environment.light"
    ENVIRONMENT_NOISE = "environment.noise"

    # =========================================================================
    # Presence Events - Location Awareness
    # =========================================================================

    PRESENCE_MOTION = "presence.motion"  # Motion in a zone
    PRESENCE_DOOR = "presence.door"  # Door open/close
    PRESENCE_OCCUPANCY = "presence.occupancy"  # Room occupancy
    PRESENCE_DEVICE_CONNECTED = "presence.device.connected"
    PRESENCE_DEVICE_DISCONNECTED = "presence.device.disconnected"
    PRESENCE_LOCATION_CHANGED = "presence.location.changed"


@dataclass
class Event:
    """
    A single event flowing through the system.

    Events are immutable once created and contain:
    - type: What happened (EventType)
    - source: Where it came from (todoist, github, etc.)
    - data: The event payload
    - timestamp: When it happened (auto-set)
    - priority: Urgency level (0=low, 5=critical)
    """

    type: Union[EventType, str]
    source: str
    data: dict = field(default_factory=dict)
    timestamp: str = field(default_factory=lambda: datetime.utcnow().isoformat())
    priority: int = 1  # 0=low, 1=normal, 2=high, 3=urgent, 4=critical
    sequence: int = 0  # Set by NATS when persisted

    def to_bytes(self) -> bytes:
        """Convert to bytes for NATS publishing."""
        return json.dumps(
            {
                "type": str(self.type),
                "source": self.source,
                "data": self.data,
                "timestamp": self.timestamp,
                "priority": self.priority,
            }
        ).encode()

    @classmethod
    def from_bytes(cls, data: bytes, sequence: int = 0) -> "Event":
        """Create from NATS message payload."""
        payload = json.loads(data.decode())
        return cls(
            type=payload.get("type", "unknown"),
            source=payload.get("source", "unknown"),
            data=payload.get("data", {}),
            timestamp=payload.get("timestamp", datetime.utcnow().isoformat()),
            priority=int(payload.get("priority", 1)),
            sequence=sequence,
        )

    @property
    def is_urgent(self) -> bool:
        """Check if this is an urgent event."""
        return self.priority >= 3

    @property
    def subject(self) -> str:
        """Get the NATS subject for this event type."""
        return f"events.{self.type}"

    def __str__(self) -> str:
        return f"Event({self.type}, priority={self.priority}, source={self.source})"


# Type alias for event handlers
EventHandler = Callable[[Event], Any]  # Can be sync or async


class EventBus:
    """
    NATS JetStream-based event bus for real-time events.

    THE NERVOUS SYSTEM OF ONEMIND OS.

    Features:
    - Sub-millisecond event delivery (200-400us)
    - Push model (events come to your handler automatically)
    - Event persistence (JetStream stores to disk)
    - Consumer groups (multiple listeners)
    - Event replay (get history from any point)

    Usage:
        bus = EventBus("nats://localhost:4222")
        await bus.connect()

        # Publish events
        await bus.publish(Event(
            type=EventType.TASK_CREATED,
            source="todoist",
            data={"title": "Call investor", "due": "tomorrow"}
        ))

        # Subscribe to events (push model - they come to you!)
        async def on_task_event(event: Event):
            print(f"Got event: {event}")

        await bus.subscribe("todoist.*", on_task_event)

        # Start listening (this processes events as they arrive)
        await bus.listen()
    """

    # Stream configuration
    STREAM_NAME = "ONEMIND"  # The living stream of experiences
    STREAM_SUBJECTS = [
        "events.>",
        "action.>",
        "legacy.>",
        "urgent.>",
    ]  # All events + actions + inbox
    CONSUMER_NAME = "legacy-agent"

    # Additional streams
    ACTIONS_STREAM = "ONEMIND_ACTIONS"  # Outbound actions for audit

    def __init__(self, nats_url: str = "nats://localhost:4222"):
        """
        Initialize the event bus.

        Args:
            nats_url: NATS server URL
        """
        self.nats_url = nats_url
        self._nc: Optional[nats.NATS] = None
        self._js: Optional[nats.js.JetStreamContext] = None
        self._handlers: dict[str, list[EventHandler]] = {}
        self._subscriptions: list = []
        self._running = False

    @property
    def nc(self) -> Optional[nats.NATS]:
        """Get the NATS connection (for advanced operations)."""
        return self._nc

    @property
    def js(self) -> Optional[nats.js.JetStreamContext]:
        """Get the JetStream context (for KV, Object Store, etc.)."""
        return self._js

    @property
    def is_connected(self) -> bool:
        """Check if connected to NATS."""
        return self._nc is not None and self._nc.is_connected

    async def connect(self) -> None:
        """Connect to NATS and set up JetStream."""
        logger.info(f"Connecting to NATS: {self.nats_url}")

        # Error callback for slow consumer and other errors (NATS best practice)
        async def error_cb(e):
            if isinstance(e, nats.errors.SlowConsumerError):
                logger.warning(f"Slow consumer detected: {e}")
                # Don't unsubscribe - just log and continue
            else:
                logger.error(f"NATS error: {e}")

        # Disconnected callback
        async def disconnected_cb():
            logger.warning("NATS disconnected")

        # Reconnected callback
        async def reconnected_cb():
            logger.info(f"NATS reconnected to {self._nc.connected_url.netloc}")

        # Connect to NATS with error handling (NATS best practice)
        self._nc = await nats.connect(
            self.nats_url,
            error_cb=error_cb,
            disconnected_cb=disconnected_cb,
            reconnected_cb=reconnected_cb,
            reconnect_time_wait=1,
            max_reconnect_attempts=-1,  # Infinite reconnects
        )
        logger.info(f"NATS connected to {self._nc.connected_url.netloc}")

        # Get JetStream context
        self._js = self._nc.jetstream()

        # Create or update the stream
        await self._setup_stream()

    async def _setup_stream(self) -> None:
        """Create or update the JetStream stream."""
        stream_config = StreamConfig(
            name=self.STREAM_NAME,
            description="OneMind OS Event Stream - Legacy's consciousness",
            subjects=self.STREAM_SUBJECTS,
            retention=RetentionPolicy.LIMITS,  # Keep until limits hit
            storage=StorageType.FILE,  # Persist to disk
            max_msgs=500000,  # Keep last 500k events (increased for 90 days)
            max_bytes=1024 * 1024 * 2048,  # 2GB max (increased for 90 days)
            max_age=60 * 60 * 24 * 90,  # 90 days retention (in seconds)
            duplicate_window=60 * 2,  # 2 minute dedup window (in seconds)
        )

        try:
            # Try to update existing stream
            await self._js.update_stream(config=stream_config)
            logger.info(f"Updated stream: {self.STREAM_NAME}")
        except nats.js.errors.NotFoundError:
            # Stream doesn't exist, create it
            await self._js.add_stream(config=stream_config)
            logger.info(f"Created stream: {self.STREAM_NAME}")
        except Exception as e:
            # Stream exists with different config, just log it
            logger.warning(f"Stream setup note: {e}")

    async def disconnect(self) -> None:
        """Disconnect from NATS."""
        self._running = False

        # Unsubscribe from all subscriptions
        for sub in self._subscriptions:
            try:
                await sub.unsubscribe()
            except Exception:
                pass

        if self._nc:
            await self._nc.drain()
            logger.info("NATS disconnected")

    async def publish(self, event: Event) -> int:
        """
        Publish an event to the stream.

        Args:
            event: The event to publish

        Returns:
            Sequence number assigned by JetStream
        """
        if not self._js:
            raise RuntimeError("Not connected to NATS")

        # Publish to JetStream
        ack = await self._js.publish(
            subject=event.subject,
            payload=event.to_bytes(),
            headers={
                "Priority": str(event.priority),
                "Source": event.source,
            },
        )

        logger.debug(f"Published: {event} -> seq={ack.seq}")

        # For urgent events, also publish to a separate subject for priority handling
        if event.is_urgent:
            await self._js.publish(
                subject=f"urgent.{event.type}",
                payload=event.to_bytes(),
            )
            logger.info(f"URGENT event published: {event}")

        return ack.seq

    async def subscribe(
        self,
        pattern: str = "*",
        handler: EventHandler = None,
    ) -> None:
        """
        Subscribe to events matching a pattern.

        NATS wildcards:
        - * matches a single token: events.todoist.* matches events.todoist.task
        - > matches multiple tokens: events.> matches all events

        Args:
            pattern: Event type pattern (e.g., "todoist.*", "github.pr.*")
            handler: Async or sync function to call for each event
        """
        if handler:
            # Convert pattern to NATS subject format
            subject = (
                f"events.{pattern}" if not pattern.startswith("events.") else pattern
            )

            if subject not in self._handlers:
                self._handlers[subject] = []
            self._handlers[subject].append(handler)
            logger.info(f"Registered handler for pattern: {subject}")

    async def listen(self, queue_group: str = "consciousness-processors") -> None:
        """
        Start listening for events.

        This sets up push subscriptions for all registered handlers.
        Events are delivered to your handlers automatically - no polling!

        Args:
            queue_group: Queue group name for load balancing across instances.
                        Multiple instances with the same queue_group share the load.
                        NATS automatically distributes messages among group members.
        """
        if not self._js:
            raise RuntimeError("Not connected to NATS")

        self._running = True
        logger.info(f"Starting event listeners (queue_group={queue_group})...")

        # Create a durable consumer for each handler pattern
        for subject, handlers in self._handlers.items():
            consumer_name = f"{self.CONSUMER_NAME}-{subject.replace('.', '-').replace('*', 'any').replace('>', 'all')}"

            try:
                # Create push subscription with native backoff for reliability
                # backoff is in nanoseconds: 1s, 5s, 30s, 2min, 10min
                backoff_ns = [
                    1_000_000_000,      # 1 second
                    5_000_000_000,      # 5 seconds
                    30_000_000_000,     # 30 seconds
                    120_000_000_000,    # 2 minutes
                    600_000_000_000,    # 10 minutes
                ]
                sub = await self._js.subscribe(
                    subject=subject,
                    durable=consumer_name,
                    queue=queue_group,  # Queue group for horizontal scaling
                    config=ConsumerConfig(
                        deliver_policy=DeliverPolicy.NEW,  # Only new messages
                        ack_wait=30,  # 30 second ack timeout
                        max_deliver=5,  # Max 5 delivery attempts before DLQ advisory
                        backoff=backoff_ns,  # Native exponential backoff
                        deliver_group=queue_group,  # JetStream deliver group
                    ),
                )
                self._subscriptions.append(sub)
                logger.info(f"Subscribed to: {subject} as {consumer_name} (queue={queue_group})")

                # Start a task to process messages for this subscription
                asyncio.create_task(self._process_subscription(sub, handlers))

            except Exception as e:
                logger.error(f"Failed to subscribe to {subject}: {e}")

        # Keep running until stopped
        while self._running:
            await asyncio.sleep(1)

    async def _process_subscription(
        self,
        sub,
        handlers: list[EventHandler],
    ) -> None:
        """Process messages from a subscription."""
        async for msg in sub.messages:
            if not self._running:
                break

            try:
                # Parse the event
                event = Event.from_bytes(msg.data, sequence=msg.seq)

                # Call all handlers
                for handler in handlers:
                    try:
                        result = handler(event)
                        # Handle both sync and async handlers
                        if asyncio.iscoroutine(result):
                            await result
                    except Exception as e:
                        logger.error(f"Handler error for {event}: {e}")

                # Acknowledge the message
                await msg.ack()

            except Exception as e:
                logger.error(f"Error processing message: {e}")
                # NAK to trigger native backoff retry (configured in ConsumerConfig)
                await msg.nak()

    async def get_recent(
        self,
        count: int = 100,
        event_type: Optional[str] = None,
    ) -> list[Event]:
        """
        Get recent events from the stream.

        Args:
            count: Number of events to retrieve
            event_type: Optional filter by event type

        Returns:
            List of recent events (newest first)
        """
        if not self._js:
            raise RuntimeError("Not connected to NATS")

        events = []
        subject = f"events.{event_type}" if event_type else "events.>"

        try:
            # Get stream info to find the last sequence
            stream_info = await self._js.stream_info(self.STREAM_NAME)
            last_seq = stream_info.state.last_seq

            if last_seq == 0:
                return events

            # Calculate starting sequence
            start_seq = max(1, last_seq - count + 1)

            # Create ephemeral consumer to read history
            sub = await self._js.subscribe(
                subject=subject,
                config=ConsumerConfig(
                    deliver_policy=DeliverPolicy.BY_START_SEQUENCE,
                    opt_start_seq=start_seq,
                ),
            )

            # Read messages
            try:
                for _ in range(count):
                    try:
                        msg = await asyncio.wait_for(
                            sub.next_msg(),
                            timeout=1.0,
                        )
                        event = Event.from_bytes(msg.data, sequence=msg.seq)
                        events.append(event)
                        await msg.ack()
                    except asyncio.TimeoutError:
                        break
            finally:
                await sub.unsubscribe()

        except Exception as e:
            logger.error(f"Error getting recent events: {e}")

        # Return in reverse order (newest first)
        return list(reversed(events))

    async def get_event_by_sequence(self, sequence: int) -> Optional[Event]:
        """Get a specific event by its sequence number."""
        if not self._js:
            raise RuntimeError("Not connected to NATS")

        try:
            # Direct get by sequence
            msg = await self._js.get_msg(
                stream=self.STREAM_NAME,
                seq=sequence,
            )
            return Event.from_bytes(msg.data, sequence=sequence)
        except Exception:
            return None

    async def replay_from(
        self,
        start_time: datetime,
        handler: EventHandler,
        event_type: Optional[str] = None,
    ) -> None:
        """
        Replay events from a specific time.

        Useful for catching up after downtime or debugging.

        Args:
            start_time: Start replaying from this time
            handler: Handler to call for each event
            event_type: Optional filter by event type
        """
        if not self._js:
            raise RuntimeError("Not connected to NATS")

        subject = f"events.{event_type}" if event_type else "events.>"

        sub = await self._js.subscribe(
            subject=subject,
            config=ConsumerConfig(
                deliver_policy=DeliverPolicy.BY_START_TIME,
                opt_start_time=start_time,
            ),
        )

        try:
            async for msg in sub.messages:
                event = Event.from_bytes(msg.data, sequence=msg.seq)
                result = handler(event)
                if asyncio.iscoroutine(result):
                    await result
                await msg.ack()

                # Check if we've caught up to "now"
                if datetime.fromisoformat(event.timestamp) >= datetime.utcnow():
                    break
        finally:
            await sub.unsubscribe()


# Convenience function to create and connect bus
async def create_event_bus(nats_url: str = "nats://localhost:4222") -> EventBus:
    """Create and connect an event bus."""
    bus = EventBus(nats_url)
    await bus.connect()
    return bus
