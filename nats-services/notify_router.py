"""
OneMind OS — Notification Router
=================================
Subscribes to cmd.notify on NATS and routes messages to configured platforms.
Reads routing config from NATS KV bucket "notifications".

Subject: cmd.notify
Payload:
{
    "message": "Text to send",
    "channels": ["telegram", "discord", "push"],  // optional, defaults to all
    "priority": "low|normal|high|critical",       // optional, default "normal"
    "title": "Optional title",                    // optional
    "context": "home|security|agents|health"      // optional, for rule matching
}

Also subscribes to legacy subjects for backward compat:
  cmd.telegram.send  → routes to telegram
  cmd.discord.send   → routes to discord

Config KV keys (bucket: "notifications"):
  config/telegram   → {"bot_token": "...", "chat_id": "...", "enabled": true}
  config/discord    → {"webhook_url": "...", "enabled": true}
  config/push       → {"fcm_key": "...", "enabled": true}
  config/sms        → {"twilio_sid": "...", "twilio_token": "...", "to": "...", "from": "...", "enabled": true}
  config/email      → {"smtp_host": "...", "smtp_port": 587, "user": "...", "pass": "...", "to": "...", "enabled": true}
  rules/dnd         → {"start": "23:00", "end": "06:00", "allow_critical": true}
  rules/routing     → {"critical": ["telegram", "sms", "push"], "high": ["telegram", "push"], "normal": ["telegram"], "low": ["push"]}

Usage:
    python3 notify_router.py

Env:
    NATS_URL (default nats://127.0.0.1:4222)
    NOTIFY_KV_BUCKET (default "notifications")
"""

import asyncio
import json
import logging
import os
import sys
from datetime import datetime, time as dtime
from typing import Optional

import nats
from nats.js.kv import KeyValue

try:
    import httpx
    HTTP_CLIENT = True
except ImportError:
    import urllib.request
    HTTP_CLIENT = False

NATS_URL = os.environ.get("NATS_URL", "nats://127.0.0.1:4222")
KV_BUCKET = os.environ.get("NOTIFY_KV_BUCKET", "notifications")
LOG_DIR = os.path.join(os.path.dirname(__file__), "logs")

# ─── Logging ───────────────────────────────────────────────────────────────

def setup_logging() -> logging.Logger:
    os.makedirs(LOG_DIR, exist_ok=True)
    logger = logging.getLogger("notify")
    logger.setLevel(logging.INFO)
    fh = logging.FileHandler(os.path.join(LOG_DIR, "notify.log"))
    fh.setFormatter(logging.Formatter("%(asctime)s %(levelname)s %(message)s"))
    ch = logging.StreamHandler()
    ch.setFormatter(logging.Formatter("%(asctime)s %(levelname)s %(message)s"))
    logger.addHandler(fh)
    logger.addHandler(ch)
    return logger

log = setup_logging()

# ─── Platform Senders ──────────────────────────────────────────────────────

async def send_telegram(message: str, title: Optional[str], config: dict) -> bool:
    """Send via Telegram Bot API."""
    token = config.get("bot_token")
    chat_id = config.get("chat_id")
    if not token or not chat_id:
        log.warning("Telegram not configured (missing bot_token or chat_id)")
        return False

    text = f"*{title}*\n{message}" if title else message
    url = f"https://api.telegram.org/bot{token}/sendMessage"
    payload = {"chat_id": chat_id, "text": text, "parse_mode": "Markdown"}

    return await _http_post(url, payload)


async def send_discord(message: str, title: Optional[str], config: dict) -> bool:
    """Send via Discord webhook."""
    webhook_url = config.get("webhook_url")
    if not webhook_url:
        log.warning("Discord not configured (missing webhook_url)")
        return False

    content = f"**{title}**\n{message}" if title else message
    payload = {"content": content}

    return await _http_post(webhook_url, payload)


async def send_push(message: str, title: Optional[str], config: dict) -> bool:
    """Send push notification via ntfy.sh (self-hostable) or FCM."""
    ntfy_url = config.get("ntfy_url", "https://ntfy.sh")
    topic = config.get("topic", "onemind")

    url = f"{ntfy_url}/{topic}"
    headers = {"Title": title or "OneMind", "Priority": "default"}
    if config.get("token"):
        headers["Authorization"] = f"Bearer {config['token']}"

    return await _http_post(url, message, headers=headers, is_text=True)


async def send_sms(message: str, title: Optional[str], config: dict) -> bool:
    """Send SMS via Twilio."""
    sid = config.get("twilio_sid")
    token = config.get("twilio_token")
    from_num = config.get("from")
    to_num = config.get("to")
    if not all([sid, token, from_num, to_num]):
        log.warning("SMS not configured (missing Twilio credentials)")
        return False

    url = f"https://api.twilio.com/2010-04-01/Accounts/{sid}/Messages.json"
    text = f"{title}: {message}" if title else message
    payload = {"To": to_num, "From": from_num, "Body": text}

    return await _http_post(url, payload, auth=(sid, token), is_form=True)


async def send_email(message: str, title: Optional[str], config: dict) -> bool:
    """Send email via SMTP."""
    import smtplib
    from email.mime.text import MIMEText

    host = config.get("smtp_host")
    port = config.get("smtp_port", 587)
    user = config.get("user")
    passwd = config.get("pass")
    to_addr = config.get("to")
    from_addr = config.get("from", user)

    if not all([host, user, passwd, to_addr]):
        log.warning("Email not configured")
        return False

    try:
        msg = MIMEText(message)
        msg["Subject"] = title or "OneMind Notification"
        msg["From"] = from_addr
        msg["To"] = to_addr

        with smtplib.SMTP(host, port) as server:
            server.starttls()
            server.login(user, passwd)
            server.send_message(msg)
        return True
    except Exception as e:
        log.error(f"Email failed: {e}")
        return False


SENDERS = {
    "telegram": send_telegram,
    "discord": send_discord,
    "push": send_push,
    "sms": send_sms,
    "email": send_email,
}

# ─── HTTP Helper ───────────────────────────────────────────────────────────

async def _http_post(url: str, data, headers=None, auth=None, is_text=False, is_form=False) -> bool:
    """Send HTTP POST. Uses httpx if available, else urllib."""
    try:
        if HTTP_CLIENT:
            async with httpx.AsyncClient(timeout=10) as client:
                kwargs = {"headers": headers or {}}
                if auth:
                    kwargs["auth"] = auth
                if is_text:
                    kwargs["content"] = data if isinstance(data, str) else str(data)
                elif is_form:
                    kwargs["data"] = data
                else:
                    kwargs["json"] = data
                resp = await client.post(url, **kwargs)
                if resp.status_code >= 400:
                    log.error(f"HTTP {resp.status_code} from {url}: {resp.text[:200]}")
                    return False
                return True
        else:
            # Fallback to urllib (sync, but works)
            import urllib.request
            req_data = json.dumps(data).encode() if not is_text else data.encode()
            req = urllib.request.Request(url, data=req_data, method="POST")
            req.add_header("Content-Type", "application/json" if not is_text else "text/plain")
            if headers:
                for k, v in headers.items():
                    req.add_header(k, v)
            urllib.request.urlopen(req, timeout=10)
            return True
    except Exception as e:
        log.error(f"HTTP POST to {url} failed: {e}")
        return False

# ─── DND / Routing Logic ──────────────────────────────────────────────────

def is_dnd(dnd_config: dict) -> bool:
    """Check if current time is in Do Not Disturb window."""
    if not dnd_config:
        return False
    try:
        start = dtime.fromisoformat(dnd_config["start"])
        end = dtime.fromisoformat(dnd_config["end"])
        now = datetime.now().time()
        if start <= end:
            return start <= now <= end
        else:  # overnight (e.g., 23:00 → 06:00)
            return now >= start or now <= end
    except (KeyError, ValueError):
        return False


def get_channels_for_priority(routing: dict, priority: str, requested: list = None) -> list:
    """Determine which channels to use based on priority and request."""
    default_routing = {
        "critical": ["telegram", "sms", "push"],
        "high": ["telegram", "push"],
        "normal": ["telegram"],
        "low": ["push"],
    }
    route_map = routing if routing else default_routing
    priority_channels = route_map.get(priority, route_map.get("normal", ["telegram"]))

    if requested:
        # Intersect requested with allowed for this priority (but critical overrides)
        if priority == "critical":
            return requested  # critical can go anywhere requested
        return [ch for ch in requested if ch in priority_channels] or priority_channels

    return priority_channels

# ─── Main Loop ─────────────────────────────────────────────────────────────

async def run():
    nc = await nats.connect(NATS_URL)
    js = nc.jetstream()
    log.info(f"Connected to NATS at {NATS_URL}")

    # Create/get KV bucket for config
    try:
        kv = await js.key_value(KV_BUCKET)
    except Exception:
        kv = await js.create_key_value(bucket=KV_BUCKET)
    log.info(f"KV bucket '{KV_BUCKET}' ready")

    # Load config helper
    async def get_config(key: str) -> dict:
        try:
            entry = await kv.get(key)
            return json.loads(entry.value.decode())
        except Exception:
            return {}

    # ─── Handler: cmd.notify ───────────────────────────────────────────
    async def handle_notify(msg):
        try:
            data = json.loads(msg.data.decode())
        except json.JSONDecodeError:
            log.warning(f"Invalid JSON on {msg.subject}: {msg.data[:100]}")
            return

        message = data.get("message", "")
        title = data.get("title")
        priority = data.get("priority", "normal")
        requested_channels = data.get("channels")
        context = data.get("context", "")

        if not message:
            log.warning("Empty message received, ignoring")
            return

        log.info(f"Notify [{priority}] {context}: {message[:80]}")

        # Check DND
        dnd_config = await get_config("rules/dnd")
        if is_dnd(dnd_config) and priority != "critical":
            if not dnd_config.get("allow_critical"):
                log.info(f"DND active, suppressing {priority} notification")
                return
            log.info(f"DND active but priority={priority}, allowing critical only")
            if priority != "critical":
                return

        # Get routing rules
        routing = await get_config("rules/routing")
        channels = get_channels_for_priority(routing, priority, requested_channels)

        # Deliver to each channel
        results = {}
        for channel in channels:
            config = await get_config(f"config/{channel}")
            if not config.get("enabled", False):
                log.debug(f"Channel '{channel}' disabled, skipping")
                continue

            sender = SENDERS.get(channel)
            if not sender:
                log.warning(f"Unknown channel: {channel}")
                continue

            success = await sender(message, title, config)
            results[channel] = success
            status = "✓" if success else "✗"
            log.info(f"  {status} {channel}")

        # Publish delivery confirmation
        await nc.publish("events.notify.delivered", json.dumps({
            "message": message[:100],
            "priority": priority,
            "channels": results,
            "timestamp": datetime.now().isoformat(),
        }).encode())

    # ─── Handler: Legacy compat (cmd.telegram.send) ────────────────────
    async def handle_telegram_legacy(msg):
        """Backward compat: cmd.telegram.send → route to telegram."""
        try:
            data = json.loads(msg.data.decode())
        except json.JSONDecodeError:
            return
        text = data.get("text", data.get("message", ""))
        if text:
            await handle_notify(type("Msg", (), {
                "data": json.dumps({
                    "message": text,
                    "channels": ["telegram"],
                    "priority": "normal",
                }).encode(),
                "subject": "cmd.telegram.send",
            })())

    # Subscribe
    await nc.subscribe("cmd.notify", cb=handle_notify)
    await nc.subscribe("cmd.telegram.send", cb=handle_telegram_legacy)
    await nc.subscribe("cmd.discord.send", cb=handle_telegram_legacy)  # same pattern

    log.info("Subscribed to cmd.notify, cmd.telegram.send, cmd.discord.send")
    log.info("Notification Router ready.")

    # Keep alive
    try:
        while True:
            await asyncio.sleep(1)
    except (KeyboardInterrupt, asyncio.CancelledError):
        pass
    finally:
        await nc.drain()


if __name__ == "__main__":
    asyncio.run(run())
