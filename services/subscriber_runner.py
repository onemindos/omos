"""
OneMind OS — Multi-Channel Subscriber Runner
=============================================
Runs ALL subscriber channels in a single process / single NATS connection.
Add new channels to CHANNELS dict — no new LaunchAgent or plist needed.

Channels are fully independent async tasks; one channel crashing restarts
only that channel (3 s back-off), not the whole process.
"""

import asyncio
import json
import logging
import os
import sys
from fnmatch import fnmatch
from datetime import datetime

import nats
from nats.js.kv import KeyValue

NATS_URL       = os.environ.get("NATS_URL", "nats://127.0.0.1:4222")
RULES_BUCKET   = "rules"
KV_CONTEXT_BUCKET = "zeus-context"
LOG_DIR        = os.path.join(os.path.dirname(__file__), "logs")
RESTART_DELAY  = 3   # seconds before restarting a crashed channel

# ── Channel registry ──────────────────────────────────────────────────────────
# To add a new subscriber: add one line here, restart the runner.
# subject_prefix supports NATS wildcards (> and *)
CHANNELS = {
    "home":     "events.zeus.>",
    "security": "events.security.>",
    "health":   "events.health.>",
    "comms":    "events.comms.>",
    "agents":   "prompt.agent.>",
}
# ─────────────────────────────────────────────────────────────────────────────


def setup_logging(name: str) -> logging.Logger:
    os.makedirs(LOG_DIR, exist_ok=True)
    logger = logging.getLogger(name)
    if logger.handlers:
        return logger
    logger.setLevel(logging.INFO)
    fh = logging.FileHandler(os.path.join(LOG_DIR, f"{name}.log"))
    fh.setFormatter(logging.Formatter("%(asctime)s %(levelname)s %(message)s"))
    ch = logging.StreamHandler()
    ch.setFormatter(logging.Formatter("%(asctime)s [%(name)s] %(levelname)s %(message)s"))
    logger.addHandler(fh)
    logger.addHandler(ch)
    return logger


root_log = setup_logging("runner")


async def load_rules(kv: KeyValue, name: str) -> list:
    rules = []
    try:
        keys = await kv.keys()
        for key in keys:
            if key.startswith(f"{name}/"):
                entry = await kv.get(key)
                if entry and entry.value:
                    rule = json.loads(entry.value.decode())
                    rule["_key"] = key
                    rules.append(rule)
    except Exception:
        pass
    return rules


async def execute_action(nc, js, action: dict, subject: str, payload: dict, logger):
    atype = action.get("type")

    if atype == "ha_service":
        req = {
            "domain":    action["domain"],
            "service":   action["service"],
            "entity_id": action.get("entity_id", ""),
        }
        resp = await nc.request("cmd.ha.service.call", json.dumps(req).encode(), timeout=5)
        logger.info(f"ha_service {action['domain']}.{action['service']} → {resp.data.decode()}")

    elif atype == "nats_publish":
        pub_payload = action.get("payload", {})
        if isinstance(pub_payload, dict):
            pub_payload = {k: v.format(**payload) if isinstance(v, str) else v
                           for k, v in pub_payload.items()}
        await nc.publish(action["subject"], json.dumps(pub_payload).encode())
        logger.info(f"nats_publish → {action['subject']}")

    elif atype == "nats_request":
        req_payload = action.get("payload", {})
        resp = await nc.request(action["subject"], json.dumps(req_payload).encode(), timeout=5)
        logger.info(f"nats_request → {action['subject']} = {resp.data.decode()[:200]}")

    elif atype == "kv_put":
        kv_bucket = await js.key_value(action.get("bucket", KV_CONTEXT_BUCKET))
        value = action["value"]
        if isinstance(value, str) and "{" in value:
            value = value.format(**payload)
        await kv_bucket.put(action["key"], value.encode())
        logger.info(f"kv_put {action.get('bucket', KV_CONTEXT_BUCKET)}.{action['key']}")

    elif atype == "telegram":
        msg = action["message"]
        try:
            msg = msg.format(subject=subject, **payload)
        except Exception:
            pass
        await nc.publish("cmd.telegram.send", json.dumps({"text": msg}).encode())
        logger.info(f"telegram → {msg[:80]}")

    elif atype == "prompt_agent":
        agent = action.get("agent", "legacy")
        prompt_text = action.get("prompt", "")
        try:
            prompt_text = prompt_text.format(subject=subject, **payload)
        except Exception:
            pass
        msg = {
            "prompt":         prompt_text,
            "source_subject": subject,
            "payload":        payload,
            "ts":             datetime.utcnow().isoformat(),
        }
        target = f"prompt.agent.{agent}"
        if action.get("request_reply", False):
            try:
                resp = await nc.request(target, json.dumps(msg).encode(), timeout=30)
                logger.info(f"prompt_agent {agent} reply: {resp.data.decode()[:200]}")
            except Exception as e:
                logger.error(f"prompt_agent {agent} error: {e}")
        else:
            await nc.publish(target, json.dumps(msg).encode())
            logger.info(f"prompt_agent → {target}: {prompt_text[:80]}")

    elif atype == "log":
        logger.info(f"[rule-log] {action.get('message', subject)}")

    else:
        logger.warning(f"Unknown action type: {atype}")


async def run_channel(nc, js, rules_kv, name: str, subject_prefix: str):
    """Single channel task — subscribes and processes messages."""
    logger = setup_logging(name)
    logger.info(f"Channel '{name}' subscribing to '{subject_prefix}'")

    async def message_handler(msg):
        await msg.ack()
        subject = msg.subject
        try:
            payload = json.loads(msg.data.decode()) if msg.data else {}
        except Exception:
            payload = {}

        logger.info(f"Event: {subject} {str(payload)[:120]}")

        rules = await load_rules(rules_kv, name) if rules_kv else []
        matched = [r for r in rules if fnmatch(subject, r.get("subject", ""))]

        if not matched:
            logger.debug(f"No matching rules for {subject}")
            return

        for rule in matched:
            logger.info(f"Rule match: {rule['_key']}")
            for action in rule.get("actions", []):
                try:
                    await execute_action(nc, js, action, subject, payload, logger)
                except Exception as e:
                    logger.error(f"Action failed: {action.get('type')} → {e}")

        await nc.publish(
            f"events.subscriber.{name}.executed",
            json.dumps({
                "subject":       subject,
                "rules_matched": len(matched),
                "ts":            datetime.utcnow().isoformat(),
            }).encode(),
        )

    await js.subscribe(
        subject_prefix,
        queue=f"sub-{name}",
        cb=message_handler,
        manual_ack=True,
    )
    logger.info(f"Channel '{name}' ready.")

    while True:
        await asyncio.sleep(1)


async def channel_supervisor(nc, js, rules_kv, name: str, subject_prefix: str):
    """Wraps run_channel with auto-restart on crash."""
    logger = setup_logging(name)
    while True:
        try:
            await run_channel(nc, js, rules_kv, name, subject_prefix)
        except asyncio.CancelledError:
            break
        except Exception as e:
            logger.error(f"Channel '{name}' crashed: {e} — restarting in {RESTART_DELAY}s")
            await asyncio.sleep(RESTART_DELAY)


async def main():
    root_log.info(f"Subscriber runner starting — {len(CHANNELS)} channels, nats={NATS_URL}")

    async def error_cb(e):        root_log.error(f"NATS error: {e}")
    async def disconnected_cb():  root_log.warning("NATS disconnected")
    async def reconnected_cb():   root_log.info("NATS reconnected")

    nc = await nats.connect(
        NATS_URL,
        error_cb=error_cb,
        disconnected_cb=disconnected_cb,
        reconnected_cb=reconnected_cb,
        max_reconnect_attempts=-1,
    )
    js = nc.jetstream()

    try:
        rules_kv = await js.key_value(RULES_BUCKET)
        root_log.info(f"Rules bucket '{RULES_BUCKET}' connected")
    except Exception:
        root_log.warning(f"KV bucket '{RULES_BUCKET}' not found — running in passthrough mode")
        rules_kv = None

    tasks = [
        asyncio.create_task(
            channel_supervisor(nc, js, rules_kv, name, prefix),
            name=f"sub-{name}",
        )
        for name, prefix in CHANNELS.items()
    ]

    root_log.info(f"All {len(tasks)} channels started.")

    try:
        await asyncio.gather(*tasks)
    except (asyncio.CancelledError, KeyboardInterrupt):
        root_log.info("Shutting down subscriber runner")
        for t in tasks:
            t.cancel()
        await asyncio.gather(*tasks, return_exceptions=True)
    finally:
        await nc.close()
        root_log.info("Subscriber runner stopped")


if __name__ == "__main__":
    asyncio.run(main())
