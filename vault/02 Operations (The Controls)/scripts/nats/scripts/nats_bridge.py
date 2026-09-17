"""
Socket.IO ↔ NATS Bridge
========================

Bridges NATS JetStream events to Socket.IO namespaces/rooms.

NATS subjects are mapped to Socket.IO rooms:
    NATS: system.>       → /system namespace, room: events:system
    NATS: digital.>      → /events namespace, room: events:digital.{subject}
    NATS: physical.>     → /events namespace, room: events:physical.{subject}
    NATS: human.>        → /approvals namespace (for approval requests)
                         → /events namespace (for notifications)

This runs as a background task inside the backend process.
"""

import asyncio
import json
import logging
from typing import Optional, TYPE_CHECKING

if TYPE_CHECKING:
    import socketio

logger = logging.getLogger("realtime.nats_bridge")


class NATSSocketIOBridge:
    """
    Bridge between NATS event bus and Socket.IO server.
    
    Subscribes to NATS subjects and pushes events to Socket.IO rooms.
    """

    def __init__(self, sio: "socketio.AsyncServer", nats_client=None):
        self.sio = sio
        self.nats = nats_client
        self._subscriptions = []
        self._running = False

    async def start(self, nats_client=None):
        """Start listening to NATS and forwarding to Socket.IO."""
        if nats_client:
            self.nats = nats_client

        if not self.nats:
            logger.warning("No NATS client provided — bridge not started")
            return

        self._running = True
        logger.info("Starting NATS → Socket.IO bridge")

        # Subscribe to all major stream subjects
        subjects = [
            ("system.>", self._handle_system_event),
            ("digital.>", self._handle_digital_event),
            ("physical.>", self._handle_physical_event),
            ("human.>", self._handle_human_event),
            ("nexus.>", self._handle_nexus_event),
            ("comms.>", self._handle_comms_event),
        ]

        for subject, handler in subjects:
            try:
                sub = await self.nats.subscribe(subject, cb=handler)
                self._subscriptions.append(sub)
                logger.info(f"Subscribed to NATS: {subject}")
            except Exception as e:
                logger.error(f"Failed to subscribe to {subject}: {e}")

    async def stop(self):
        """Stop all NATS subscriptions."""
        self._running = False
        for sub in self._subscriptions:
            try:
                await sub.unsubscribe()
            except Exception:
                pass
        self._subscriptions.clear()
        logger.info("NATS → Socket.IO bridge stopped")

    async def _handle_system_event(self, msg):
        """Forward system events to /system namespace."""
        try:
            data = json.loads(msg.data.decode())
            subject = msg.subject

            await self.sio.emit(
                "nats_event",
                {"subject": subject, "data": data},
                namespace="/system",
                room="events:system",
            )
        except Exception as e:
            logger.error(f"Error forwarding system event: {e}")

    async def _handle_digital_event(self, msg):
        """Forward digital events to /events namespace."""
        try:
            data = json.loads(msg.data.decode())
            subject = msg.subject

            # Emit to the specific subject room AND the wildcard room
            await self.sio.emit(
                "nats_event",
                {"subject": subject, "data": data},
                namespace="/events",
                room=f"events:{subject}",
            )

            # Also emit to the domain-level room
            await self.sio.emit(
                "nats_event",
                {"subject": subject, "data": data},
                namespace="/events",
                room="events:digital.>",
            )
        except Exception as e:
            logger.error(f"Error forwarding digital event: {e}")

    async def _handle_physical_event(self, msg):
        """Forward physical events to /events namespace."""
        try:
            data = json.loads(msg.data.decode())
            subject = msg.subject

            await self.sio.emit(
                "nats_event",
                {"subject": subject, "data": data},
                namespace="/events",
                room=f"events:{subject}",
            )

            await self.sio.emit(
                "nats_event",
                {"subject": subject, "data": data},
                namespace="/events",
                room="events:physical.>",
            )
        except Exception as e:
            logger.error(f"Error forwarding physical event: {e}")

    async def _handle_human_event(self, msg):
        """Forward human events to /approvals namespace."""
        try:
            data = json.loads(msg.data.decode())
            subject = msg.subject

            if "approval" in subject:
                # Approval requests go to /approvals namespace
                await self.sio.emit(
                    "new_approval",
                    {"subject": subject, "data": data},
                    namespace="/approvals",
                    room="approvals:all",
                )
            else:
                # Notifications go to /events
                await self.sio.emit(
                    "nats_event",
                    {"subject": subject, "data": data},
                    namespace="/events",
                    room="events:human.>",
                )
        except Exception as e:
            logger.error(f"Error forwarding human event: {e}")

    async def _handle_nexus_event(self, msg):
        """Forward Nexus orchestration events to /events namespace."""
        try:
            data = json.loads(msg.data.decode())
            subject = msg.subject

            # Nexus events go to /events namespace under nexus rooms
            await self.sio.emit(
                "nexus_event",
                {"subject": subject, "data": data},
                namespace="/events",
                room=f"events:{subject}",
            )

            # Also emit to the wildcard nexus room
            await self.sio.emit(
                "nexus_event",
                {"subject": subject, "data": data},
                namespace="/events",
                room="events:nexus.>",
            )
        except Exception as e:
            logger.error(f"Error forwarding nexus event: {e}")

    async def _handle_comms_event(self, msg):
        """Forward universal comms events to /chat namespace."""
        try:
            data = json.loads(msg.data.decode())
            subject = msg.subject

            # Chat/voice/vision messages go to /chat namespace
            await self.sio.emit(
                "comms_event",
                {"subject": subject, "data": data},
                namespace="/chat",
                room=f"comms:{subject}",
            )

            # Also emit to the modality-level room (e.g. comms:voice, comms:vision)
            parts = subject.split(".")
            if len(parts) > 1:
                modality = parts[1]
                await self.sio.emit(
                    "comms_event",
                    {"subject": subject, "data": data},
                    namespace="/chat",
                    room=f"comms:{modality}",
                )

            # And a catch-all room for the Nexus screen
            await self.sio.emit(
                "comms_event",
                {"subject": subject, "data": data},
                namespace="/events",
                room="events:comms.>",
            )
        except Exception as e:
            logger.error(f"Error forwarding comms event: {e}")
