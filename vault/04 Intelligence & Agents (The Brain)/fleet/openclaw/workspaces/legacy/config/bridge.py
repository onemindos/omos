"""
NATS → OpenClaw Webhook Bridge
================================

Subscribes to key NATS subjects and forwards relevant events to the
OpenClaw Gateway webhook endpoint. This gives Legacy full awareness of
everything happening in the OneMind OS ecosystem without OpenClaw
needing to speak NATS natively.

Flow:
  NATS (Agno results, OM1 status, system events)
      ↓ this bridge
  OpenClaw /hooks/wake or /hooks/agent
      ↓
  Legacy sees the event, reasons, responds

Usage:
    python -m backend.legacy.bridge

Environment:
    NATS_URL                   NATS server URL (default: nats://localhost:4222)
    OPENCLAW_GATEWAY_URL       OpenClaw gateway URL (default: http://localhost:18789)
    OPENCLAW_HOOKS_TOKEN       Token for OpenClaw webhook auth (required)
    BRIDGE_AGENT_ID            OpenClaw agent to route events to (default: legacy)
"""

import asyncio
import json
import logging
import os
import sys
from datetime import datetime
from typing import Any, Dict, Optional

import httpx

logger = logging.getLogger("onemind.legacy.bridge")
logging.basicConfig(
    level=os.getenv("LOG_LEVEL", "INFO"),
    format="%(asctime)s [%(name)s] %(levelname)s %(message)s",
)

# Configuration
NATS_URL = os.getenv("NATS_URL", "nats://localhost:4222")
OPENCLAW_GATEWAY_URL = os.getenv("OPENCLAW_GATEWAY_URL", "http://localhost:18789")
OPENCLAW_HOOKS_TOKEN = os.getenv("OPENCLAW_HOOKS_TOKEN", "")
BRIDGE_AGENT_ID = os.getenv("BRIDGE_AGENT_ID", "legacy")

if not OPENCLAW_HOOKS_TOKEN:
    logger.error("OPENCLAW_HOOKS_TOKEN is not set — bridge cannot authenticate with OpenClaw")
    sys.exit(1)


# ==============================================================================
# Event routing rules
# Map NATS subjects to OpenClaw webhook actions + formatting
# ==============================================================================

def build_wake_payload(subject: str, data: Dict[str, Any]) -> Optional[Dict]:
    """Build a /hooks/wake payload for informational events."""
    return {
        "text": f"[{subject}] {_summarize(data)}",
        "mode": "now",
    }


def build_agent_payload(subject: str, data: Dict[str, Any], deliver: bool = False) -> Optional[Dict]:
    """Build a /hooks/agent payload for events requiring Legacy's attention."""
    return {
        "message": f"System event on `{subject}`:\n\n{json.dumps(data, indent=2, default=str)[:2000]}",
        "name": _source_name(subject),
        "agentId": BRIDGE_AGENT_ID,
        "wakeMode": "now",
        "deliver": deliver,
        "channel": "whatsapp" if deliver else "last",
    }


def _summarize(data: Dict) -> str:
    """Short text summary of event data."""
    if isinstance(data, dict):
        for key in ("message", "text", "summary", "status", "result", "task"):
            if key in data:
                val = str(data[key])
                return val[:200] + ("..." if len(val) > 200 else "")
    return str(data)[:200]


def _source_name(subject: str) -> str:
    """Human-readable source name from NATS subject."""
    parts = subject.split(".")
    if len(parts) >= 3:
        return f"{parts[0].capitalize()} — {parts[2].replace('_', ' ').title()}"
    return subject


# Subject → routing rules
# Format: (pattern, endpoint, should_deliver_to_user)
ROUTING_RULES = [
    # Agent results from Agno — wake Legacy (she dispatched these)
    ("digital.agent.*.result.legacy", "wake", False),

    # Critical system alerts — wake immediately
    ("system.alert.>", "wake", False),
    ("system.error.>", "wake", False),

    # Physical robot status changes — wake Legacy
    ("physical.robot.*.status", "wake", False),
    ("physical.robot.*.alert", "agent", True),   # Robot alerts → deliver to Zeus

    # Home Assistant critical events — deliver to Zeus
    ("physical.ha.alert.>", "agent", True),

    # Nexus heartbeat failures — wake Legacy
    ("nexus.heartbeat.failure.>", "wake", False),

    # Agno backend health changes
    ("system.health.agno", "wake", False),

    # OM1 heartbeat loss
    ("physical.heartbeat.lost.>", "agent", True),

    # Community events (circle.so webhooks forwarded by backend)
    ("events.community.member.joined", "agent", False),

    # Task completions (ClickUp webhooks)
    ("events.task.completed.>", "wake", False),
]


def match_subject(subject: str, pattern: str) -> bool:
    """Match NATS subject against a pattern (supports * and >)."""
    parts = subject.split(".")
    pat_parts = pattern.split(".")
    for i, p in enumerate(pat_parts):
        if p == ">":
            return True
        if i >= len(parts):
            return False
        if p != "*" and p != parts[i]:
            return False
    return len(parts) == len(pat_parts)


def get_routing(subject: str):
    """Return (endpoint, deliver) for a NATS subject, or None if not routed."""
    for pattern, endpoint, deliver in ROUTING_RULES:
        if match_subject(subject, pattern):
            return endpoint, deliver
    return None, None


# ==============================================================================
# OpenClaw webhook client
# ==============================================================================

async def post_to_openclaw(endpoint: str, payload: Dict, client: httpx.AsyncClient) -> bool:
    """POST to OpenClaw webhook. Returns True on success."""
    url = f"{OPENCLAW_GATEWAY_URL}/hooks/{endpoint}"
    try:
        response = await client.post(
            url,
            json=payload,
            headers={
                "Authorization": f"Bearer {OPENCLAW_HOOKS_TOKEN}",
                "Content-Type": "application/json",
            },
            timeout=10.0,
        )
        if response.status_code in (200, 202):
            logger.debug(f"✓ Forwarded to OpenClaw /{endpoint}: {response.status_code}")
            return True
        else:
            logger.warning(f"OpenClaw /{endpoint} returned {response.status_code}: {response.text[:200]}")
            return False
    except httpx.RequestError as e:
        logger.error(f"Failed to reach OpenClaw at {url}: {e}")
        return False


# ==============================================================================
# NATS subscriber
# ==============================================================================

async def run_bridge():
    """Main bridge loop — subscribe to NATS, forward to OpenClaw."""
    try:
        import nats
    except ImportError:
        logger.error("nats-py not installed. Run: pip install nats-py")
        sys.exit(1)

    logger.info(f"NATS→OpenClaw bridge starting")
    logger.info(f"  NATS:     {NATS_URL}")
    logger.info(f"  OpenClaw: {OPENCLAW_GATEWAY_URL}")
    logger.info(f"  Agent:    {BRIDGE_AGENT_ID}")

    async with httpx.AsyncClient() as http_client:
        # Connect to NATS
        nc = None
        while True:
            try:
                nc = await nats.connect(
                    NATS_URL,
                    max_reconnect_attempts=-1,
                    reconnect_time_wait=2,
                )
                logger.info("Connected to NATS")
                break
            except Exception as e:
                logger.warning(f"NATS not ready ({e}), retrying in 5s...")
                await asyncio.sleep(5)

        # Subscribe to all monitored subjects
        monitored_patterns = set(pattern for pattern, _, _ in ROUTING_RULES)

        async def message_handler(msg):
            subject = msg.subject
            try:
                data = json.loads(msg.data.decode())
            except (json.JSONDecodeError, UnicodeDecodeError):
                data = {"raw": msg.data.decode(errors="replace")}

            endpoint, deliver = get_routing(subject)
            if not endpoint:
                return

            logger.info(f"Event: {subject} → /{endpoint} (deliver={deliver})")

            if endpoint == "wake":
                payload = build_wake_payload(subject, data)
            else:
                payload = build_agent_payload(subject, data, deliver=deliver)

            if payload:
                await post_to_openclaw(endpoint, payload, http_client)

        # Subscribe to all relevant subjects
        subscriptions = []
        for pattern in monitored_patterns:
            sub = await nc.subscribe(pattern, cb=message_handler)
            subscriptions.append(sub)
            logger.debug(f"Subscribed: {pattern}")

        logger.info(f"Bridge active — monitoring {len(subscriptions)} subject patterns")

        # Keep alive
        try:
            while True:
                await asyncio.sleep(30)
                if nc.is_closed:
                    logger.error("NATS connection closed — exiting")
                    break
        except (KeyboardInterrupt, asyncio.CancelledError):
            logger.info("Bridge shutting down...")
        finally:
            for sub in subscriptions:
                try:
                    await sub.unsubscribe()
                except Exception:
                    pass
            await nc.drain()
            logger.info("Bridge stopped")


if __name__ == "__main__":
    asyncio.run(run_bridge())
