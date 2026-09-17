"""
Inbound Webhook Receiver
========================

Receives webhooks from external services and publishes them to NATS.
This is the entry point for external events into the OneMind OS nervous system.

Usage:
    POST /webhooks/{source}     - Receive webhook from any source
    POST /webhooks/{source}/{event_type}  - With explicit event type

Headers:
    X-Webhook-Secret: Optional shared secret for verification
    X-Event-Type: Optional event type override

The webhook data is published to NATS as:
    events.{source}.{event_type}

Examples:
    POST /webhooks/github       -> events.github.push (inferred from payload)
    POST /webhooks/todoist      -> events.todoist.item_added
    POST /webhooks/custom/alert -> events.custom.alert
"""

import json
import hmac
import hashlib
import logging
import os
from datetime import datetime, timezone
from typing import Optional, Dict, Any
from uuid import uuid4

from fastapi import APIRouter, Request, HTTPException, Header, BackgroundTasks
from pydantic import BaseModel

logger = logging.getLogger("webhooks.receiver")

router = APIRouter(prefix="/webhooks", tags=["webhooks"])


# =============================================================================
# WEBHOOK SECRET VERIFICATION
# =============================================================================

# Source-specific secrets (loaded from environment)
WEBHOOK_SECRETS: Dict[str, str] = {
    "github": os.getenv("GITHUB_WEBHOOK_SECRET", ""),
    "todoist": os.getenv("TODOIST_WEBHOOK_SECRET", ""),
    "linear": os.getenv("LINEAR_WEBHOOK_SECRET", ""),
    "slack": os.getenv("SLACK_WEBHOOK_SECRET", ""),
    "stripe": os.getenv("STRIPE_WEBHOOK_SECRET", ""),
    "default": os.getenv("WEBHOOK_SECRET", ""),  # Fallback
}


def verify_signature(
    source: str,
    payload: bytes,
    signature: Optional[str],
    signature_header: str = "X-Webhook-Signature",
) -> bool:
    """Verify webhook signature if secret is configured."""
    secret = WEBHOOK_SECRETS.get(source) or WEBHOOK_SECRETS.get("default")

    if not secret:
        # No secret configured, allow all
        return True

    if not signature:
        logger.warning(f"No signature provided for {source} webhook")
        return False

    # Different services use different signature formats
    if source == "github":
        # GitHub: sha256=<signature>
        expected = "sha256=" + hmac.new(
            secret.encode(), payload, hashlib.sha256
        ).hexdigest()
        return hmac.compare_digest(expected, signature)

    elif source == "stripe":
        # Stripe uses timestamp + signature
        # Simplified - real implementation would parse Stripe's format
        return True

    else:
        # Default: simple HMAC comparison
        expected = hmac.new(secret.encode(), payload, hashlib.sha256).hexdigest()
        return hmac.compare_digest(expected, signature)


# =============================================================================
# EVENT TYPE INFERENCE
# =============================================================================

def infer_event_type(source: str, data: Dict[str, Any]) -> str:
    """Infer the event type from the webhook payload."""

    # GitHub
    if source == "github":
        action = data.get("action", "")
        if "pull_request" in data:
            return f"pull_request.{action}" if action else "pull_request"
        elif "issue" in data:
            return f"issues.{action}" if action else "issues"
        elif "commits" in data or "head_commit" in data:
            return "push"
        elif "ref_type" in data:
            return "create"
        return data.get("action", "event")

    # Todoist
    elif source == "todoist":
        return data.get("event_name", "event")

    # Linear
    elif source == "linear":
        return data.get("type", "event")

    # Home Assistant
    elif source == "homeassistant" or source == "ha":
        event_type = data.get("event_type", "event")
        # Normalize HA events
        if event_type.startswith("homeassistant_"):
            event_type = event_type.replace("homeassistant_", "")
        return event_type

    # Slack
    elif source == "slack":
        return data.get("type", data.get("event", {}).get("type", "event"))

    # Stripe
    elif source == "stripe":
        return data.get("type", "event").replace(".", "_")

    # Calendar (Google, Apple, etc.)
    elif source in ["calendar", "gcal", "apple_calendar"]:
        return data.get("eventType", data.get("type", "event"))

    # Oura
    elif source == "oura":
        return data.get("event_type", "event")

    # Default: look for common fields
    else:
        for field in ["event_type", "eventType", "type", "action", "event"]:
            if field in data:
                return str(data[field])
        return "event"


# =============================================================================
# NATS PUBLISHER
# =============================================================================

async def publish_to_nats(subject: str, data: Dict[str, Any]):
    """Publish event to NATS JetStream."""
    try:
        from backend.consciousness.core.bus import get_event_bus

        bus = await get_event_bus()
        if bus:
            await bus.publish(subject, data)
            logger.info(f"Published to NATS: {subject}")
            return True
        else:
            logger.warning("NATS bus not available, event not published")
            return False

    except ImportError:
        logger.warning("NATS bus module not available")
        return False
    except Exception as e:
        logger.error(f"Failed to publish to NATS: {e}")
        return False


# =============================================================================
# WEBHOOK ENDPOINTS
# =============================================================================


class WebhookResponse(BaseModel):
    """Response model for webhook receipt."""
    status: str
    event_id: str
    subject: str
    message: str


@router.post("/{source}", response_model=WebhookResponse)
@router.post("/{source}/{event_type}", response_model=WebhookResponse)
async def receive_webhook(
    request: Request,
    source: str,
    event_type: Optional[str] = None,
    background_tasks: BackgroundTasks = None,
    x_webhook_signature: Optional[str] = Header(None, alias="X-Webhook-Signature"),
    x_hub_signature_256: Optional[str] = Header(None, alias="X-Hub-Signature-256"),
    x_event_type: Optional[str] = Header(None, alias="X-Event-Type"),
    x_github_event: Optional[str] = Header(None, alias="X-GitHub-Event"),
):
    """
    Receive a webhook from an external service.

    The webhook payload is published to NATS as: events.{source}.{event_type}

    Args:
        source: The source service (github, todoist, linear, etc.)
        event_type: Optional explicit event type (inferred from payload if not provided)

    Headers:
        X-Webhook-Signature: Signature for verification (source-specific)
        X-Hub-Signature-256: GitHub's signature header
        X-Event-Type: Override event type
        X-GitHub-Event: GitHub's event type header
    """
    # Get raw body for signature verification
    body = await request.body()

    # Verify signature
    signature = x_hub_signature_256 or x_webhook_signature
    if not verify_signature(source, body, signature):
        raise HTTPException(status_code=401, detail="Invalid webhook signature")

    # Parse payload
    try:
        data = json.loads(body) if body else {}
    except json.JSONDecodeError:
        # Some webhooks use form data
        form = await request.form()
        data = dict(form)

    # Determine event type
    final_event_type = (
        event_type or
        x_event_type or
        x_github_event or
        infer_event_type(source, data)
    )

    # Normalize source name
    source = source.lower().replace("-", "_")

    # Build NATS subject
    subject = f"events.{source}.{final_event_type}"

    # Generate event ID
    event_id = str(uuid4())

    # Enrich payload
    enriched_data = {
        "event_id": event_id,
        "source": source,
        "event_type": final_event_type,
        "timestamp": datetime.now(timezone.utc).isoformat(),
        "data": data,
        # Preserve original fields at top level for compatibility
        **data,
    }

    # Publish to NATS (in background to not block response)
    if background_tasks:
        background_tasks.add_task(publish_to_nats, subject, enriched_data)
        message = "Webhook received, processing in background"
    else:
        success = await publish_to_nats(subject, enriched_data)
        message = "Webhook received and published" if success else "Webhook received, NATS unavailable"

    logger.info(f"Webhook received: {subject} (event_id={event_id})")

    return WebhookResponse(
        status="received",
        event_id=event_id,
        subject=subject,
        message=message,
    )


@router.get("/sources")
async def list_webhook_sources():
    """List available webhook sources and their event types."""
    return {
        "sources": {
            "github": {
                "events": ["push", "pull_request.opened", "pull_request.merged",
                          "pull_request.review_requested", "issues.opened", "issues.assigned"],
                "signature_header": "X-Hub-Signature-256",
            },
            "todoist": {
                "events": ["item_added", "item_completed", "item_updated", "item_deleted"],
                "signature_header": "X-Todoist-Hmac-SHA256",
            },
            "linear": {
                "events": ["Issue.created", "Issue.updated", "Comment.created"],
                "signature_header": "Linear-Signature",
            },
            "homeassistant": {
                "events": ["state_changed", "automation_triggered", "script_started"],
                "signature_header": "X-Webhook-Signature",
            },
            "oura": {
                "events": ["daily_readiness", "daily_sleep", "daily_activity"],
                "signature_header": "X-Oura-Signature",
            },
            "calendar": {
                "events": ["event.created", "event.updated", "event.starting"],
                "signature_header": "X-Webhook-Signature",
            },
            "custom": {
                "events": ["any"],
                "description": "Generic webhook endpoint for custom integrations",
            },
        },
        "usage": {
            "endpoint": "POST /webhooks/{source}",
            "example": "POST /webhooks/github",
            "with_type": "POST /webhooks/custom/my_event_type",
        },
    }


@router.post("/test")
async def test_webhook(request: Request):
    """Test endpoint that echoes webhook data and publishes to test stream."""
    body = await request.body()
    try:
        data = json.loads(body) if body else {}
    except json.JSONDecodeError:
        data = {"raw": body.decode()}

    # Publish to test subject
    subject = "events.test.webhook"
    event_id = str(uuid4())

    enriched = {
        "event_id": event_id,
        "source": "test",
        "event_type": "webhook",
        "timestamp": datetime.now(timezone.utc).isoformat(),
        "data": data,
    }

    await publish_to_nats(subject, enriched)

    return {
        "status": "received",
        "event_id": event_id,
        "subject": subject,
        "echo": data,
    }
