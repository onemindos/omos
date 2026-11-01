"""nats-context — OneMind OS omnipresence plugin for Hermes.

Injects a live [LIVE OPERATIONAL PICTURE] block into every LLM call via the
pre_llm_call hook, and exposes NATS tools (publish, request, peek, JetStream,
KV, server info, stream listing) so Hermes can actively interact with the bus.
"""
from __future__ import annotations

import logging
import os
from typing import Any

from .state import LiveState
from .subscriptions import SubscriptionManager
from .formatter import format_context
from .tools import (
    nats_publish,
    nats_request,
    nats_peek,
    nats_jetstream_fetch,
    nats_kv_get,
    nats_kv_put,
    nats_kv_list,
    nats_server_info,
    nats_list_streams,
)

logger = logging.getLogger("nats-context")

# Module-level singletons — one per Hermes process
_state: LiveState | None = None
_manager: SubscriptionManager | None = None


def register(ctx: Any) -> None:
    """Entry point called by Hermes when the plugin is loaded."""
    global _state, _manager

    nats_url = os.environ.get("NATS_URL", "nats://localhost:4222")
    logger.info("nats-context: registering plugin (NATS_URL=%s)", nats_url)

    _state = LiveState()
    _manager = SubscriptionManager(_state, nats_url=nats_url)
    _manager.start()

    # ── pre_llm_call hook ──────────────────────────────────────────────────
    def _inject_operational_picture(payload: dict) -> str | dict:
        """Append the current operational picture to every user message."""
        if _state is None:
            return ""
        snapshot = _state.snapshot()

        # Skip injection if nothing has been observed yet
        has_data = any([
            snapshot["tak"],
            snapshot["entities"],
            snapshot["agents"],
            snapshot["bridges"],
            snapshot["threats"],
            snapshot["intel"],
            snapshot["detections"],
            snapshot["telemetry"],
            snapshot["nats_stats"],
        ])
        if not has_data:
            return ""

        context = format_context(snapshot)
        return {"context": context}

    ctx.register_hook("pre_llm_call", _inject_operational_picture)

    # ── NATS tools ─────────────────────────────────────────────────────────
    ctx.register_tool(
        toolset="nats-context",
        name="nats_publish",
        description=(
            "Publish a message to any NATS subject. "
            "Use to send commands, trigger workflows, or push state updates. "
            "Example subjects: cmd.agent.hermes.zeus.run, evt.task.42.complete"
        ),
        schema={
            "type": "object",
            "properties": {
                "subject": {"type": "string", "description": "NATS subject to publish to"},
                "payload": {"type": "string", "description": "JSON-encoded message body", "default": "{}"},
            },
            "required": ["subject"],
        },
        handler=nats_publish,
        is_async=False,
    )

    ctx.register_tool(
        toolset="nats-context",
        name="nats_request",
        description=(
            "Send a request to a NATS subject and wait for the first reply (request/reply pattern). "
            "Useful for querying services synchronously."
        ),
        schema={
            "type": "object",
            "properties": {
                "subject": {"type": "string"},
                "payload": {"type": "string", "default": "{}"},
                "timeout": {"type": "number", "default": 5.0, "description": "Seconds to wait for reply"},
            },
            "required": ["subject"],
        },
        handler=nats_request,
        is_async=False,
    )

    ctx.register_tool(
        toolset="nats-context",
        name="nats_peek",
        description=(
            "Subscribe to a NATS subject wildcard for 2 seconds and collect up to n messages. "
            "Use to inspect live traffic on a subject, e.g. 'tak.cot.>' or 'ent.drone.>'."
        ),
        schema={
            "type": "object",
            "properties": {
                "subject": {"type": "string", "description": "Subject pattern (wildcards * and > supported)"},
                "n": {"type": "integer", "default": 10, "description": "Maximum messages to collect"},
            },
            "required": ["subject"],
        },
        handler=nats_peek,
        is_async=False,
    )

    ctx.register_tool(
        toolset="nats-context",
        name="nats_jetstream_fetch",
        description=(
            "Fetch recent messages from a JetStream stream. "
            "Streams available: TAK, ENTITY, TELEMETRY, INTEL, EVENTS, AGENTS, BRIDGES."
        ),
        schema={
            "type": "object",
            "properties": {
                "stream": {"type": "string", "description": "JetStream stream name"},
                "subject_filter": {"type": "string", "description": "Optional subject filter within the stream"},
                "n": {"type": "integer", "default": 10, "description": "Max messages to fetch"},
            },
            "required": ["stream"],
        },
        handler=nats_jetstream_fetch,
        is_async=False,
    )

    ctx.register_tool(
        toolset="nats-context",
        name="nats_kv_get",
        description="Read a value from a NATS KV bucket by key.",
        schema={
            "type": "object",
            "properties": {
                "bucket": {"type": "string"},
                "key": {"type": "string"},
            },
            "required": ["bucket", "key"],
        },
        handler=nats_kv_get,
        is_async=False,
    )

    ctx.register_tool(
        toolset="nats-context",
        name="nats_kv_put",
        description="Write a value to a NATS KV bucket key.",
        schema={
            "type": "object",
            "properties": {
                "bucket": {"type": "string"},
                "key": {"type": "string"},
                "value": {"description": "Any JSON-serialisable value"},
            },
            "required": ["bucket", "key"],
        },
        handler=nats_kv_put,
        is_async=False,
    )

    ctx.register_tool(
        toolset="nats-context",
        name="nats_kv_list",
        description="List all keys in a NATS KV bucket.",
        schema={
            "type": "object",
            "properties": {
                "bucket": {"type": "string"},
            },
            "required": ["bucket"],
        },
        handler=nats_kv_list,
        is_async=False,
    )

    ctx.register_tool(
        toolset="nats-context",
        name="nats_server_info",
        description="Return NATS server metadata: version, cluster, connections, JetStream status.",
        schema={
            "type": "object",
            "properties": {},
        },
        handler=lambda: nats_server_info(),
        is_async=False,
    )

    ctx.register_tool(
        toolset="nats-context",
        name="nats_list_streams",
        description="List all JetStream streams: names, subject filters, message counts, byte sizes.",
        schema={
            "type": "object",
            "properties": {},
        },
        handler=lambda: nats_list_streams(),
        is_async=False,
    )

    logger.info("nats-context: registered pre_llm_call hook + 9 NATS tools")
