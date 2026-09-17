"""
System Bus - NATS Client Singleton
====================================

The nervous system of OneMind OS. Every subsystem publishes and
subscribes through this single NATS connection.

Usage:
    from backend.bus import get_bus, publish, subscribe

    # Publish an event
    await publish("digital.agent.action", {"action": "notify", "text": "..."})

    # Subscribe to a topic
    async def handler(msg):
        data = json.loads(msg.data)
        ...
    await subscribe("physical.>", handler)

    # Get the raw NATS client
    nc = await get_bus()
"""

import os
import json
import asyncio
import logging
from typing import Optional, Any, Dict, Callable, Awaitable
from datetime import datetime

logger = logging.getLogger("onemind.bus")

# Singleton NATS connection
_nc = None
_js = None
_connected = False
_subscriptions: list = []

# Connection stats
_stats = {
    "connected_at": None,
    "messages_published": 0,
    "messages_received": 0,
    "reconnects": 0,
    "last_error": None,
}


async def get_bus():
    """Get the NATS connection, connecting if needed."""
    global _nc, _js, _connected, _stats
    if _nc and _connected:
        return _nc

    try:
        import nats as nats_lib
        from nats.js.api import StreamConfig

        nats_url = os.getenv("NATS_URL", "nats://localhost:4222")
        logger.info(f"Connecting to NATS at {nats_url}")

        async def disconnected_cb():
            global _connected
            _connected = False
            logger.warning("NATS disconnected")

        async def reconnected_cb():
            global _connected, _stats
            _connected = True
            _stats["reconnects"] += 1
            logger.info("NATS reconnected")

        async def error_cb(e):
            global _stats
            _stats["last_error"] = str(e)
            logger.error(f"NATS error: {e}")

        _nc = await nats_lib.connect(
            nats_url,
            disconnected_cb=disconnected_cb,
            reconnected_cb=reconnected_cb,
            error_cb=error_cb,
            max_reconnect_attempts=-1,  # infinite reconnect
            reconnect_time_wait=2,
        )
        _connected = True
        _stats["connected_at"] = datetime.utcnow().isoformat()

        # Set up JetStream
        _js = _nc.jetstream()

        # Create core streams
        streams = [
            ("DIGITAL", ["digital.>"]),
            ("PHYSICAL", ["physical.>"]),
            ("HUMAN", ["human.>"]),
            ("SYSTEM", ["system.>"]),
            ("EVENTS", ["events.>"]),
            ("NEXUS", ["nexus.>"]),
            ("COMMS", ["comms.>"]),
        ]
        for name, subjects in streams:
            try:
                await _js.add_stream(name=name, subjects=subjects)
                logger.debug(f"Stream '{name}' ready")
            except Exception:
                # Stream may already exist
                pass

        logger.info(f"NATS connected with JetStream ({len(streams)} streams)")
        return _nc

    except Exception as e:
        logger.warning(f"NATS not available: {e} — running without event bus")
        _stats["last_error"] = str(e)
        return None


async def publish(subject: str, data: Any, headers: Optional[Dict[str, str]] = None):
    """
    Publish a message to the bus.

    Args:
        subject: NATS subject (e.g., "digital.agent.action", "system.heartbeat")
        data: Dict or string payload (auto-serialized to JSON)
        headers: Optional NATS headers
    """
    global _stats
    nc = await get_bus()
    if not nc:
        return False

    try:
        if isinstance(data, dict):
            payload = json.dumps(data, default=str).encode()
        elif isinstance(data, str):
            payload = data.encode()
        elif isinstance(data, bytes):
            payload = data
        else:
            payload = json.dumps(data, default=str).encode()

        await nc.publish(subject, payload, headers=headers)
        _stats["messages_published"] += 1
        return True
    except Exception as e:
        logger.error(f"Failed to publish to {subject}: {e}")
        _stats["last_error"] = str(e)
        return False


async def subscribe(subject: str, handler: Callable, queue: str = None):
    """
    Subscribe to a subject pattern.

    Args:
        subject: NATS subject pattern (supports wildcards: >, *)
        handler: Async callback(msg) 
        queue: Optional queue group for load balancing
    """
    global _stats
    nc = await get_bus()
    if not nc:
        return None

    try:
        async def wrapped_handler(msg):
            _stats["messages_received"] += 1
            await handler(msg)

        if queue:
            sub = await nc.subscribe(subject, cb=wrapped_handler, queue=queue)
        else:
            sub = await nc.subscribe(subject, cb=wrapped_handler)

        _subscriptions.append(sub)
        logger.info(f"Subscribed to: {subject}")
        return sub
    except Exception as e:
        logger.error(f"Failed to subscribe to {subject}: {e}")
        return None


def get_bus_status() -> Dict[str, Any]:
    """Get bus connection status and stats."""
    return {
        "connected": _connected,
        "url": os.getenv("NATS_URL", "nats://localhost:4222"),
        "stats": _stats.copy(),
        "subscriptions": len(_subscriptions),
    }


async def shutdown_bus():
    """Gracefully close the NATS connection."""
    global _nc, _connected
    if _nc:
        try:
            await _nc.drain()
            _connected = False
            logger.info("NATS connection closed")
        except Exception as e:
            logger.error(f"Error closing NATS: {e}")
