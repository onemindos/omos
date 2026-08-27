"""
OneMind OS — Config-Driven NATS Subscriber
==========================================
Reads rules from NATS KV bucket "rules", subscribes to subjects,
executes actions without LLM involvement.

Supports agent invocation: publish to prompt.agent.<name> to wake agents.

Usage:
    python3 subscriber.py <name> [--subject-prefix events.zeus.>]

Rules KV key format:  <name>/<rule_id>
Rules KV value:       JSON  { "subject": "events.zeus.*", "actions": [...] }

Action types:
  ha_service   — call Home Assistant service via NATS
  nats_publish — fire-and-forget to any subject
  nats_request — request/reply to any subject
  telegram     — send Telegram message via cmd.telegram.send
  kv_put       — write to any KV bucket
  prompt_agent — wake an agent (Legacy, Codex, OpenClaw) with a prompt
  log          — log only, no action
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

NATS_URL = os.environ.get("NATS_URL", "nats://127.0.0.1:4222")
RULES_BUCKET = "rules"
KV_CONTEXT_BUCKET = "zeus-context"
LOG_DIR = os.path.join(os.path.dirname(__file__), "logs")


def setup_logging(name: str) -> logging.Logger:
    os.makedirs(LOG_DIR, exist_ok=True)
    logger = logging.getLogger(name)
    logger.setLevel(logging.INFO)
    fh = logging.FileHandler(os.path.join(LOG_DIR, f"{name}.log"))
    fh.setFormatter(logging.Formatter("%(asctime)s %(levelname)s %(message)s"))
    ch = logging.StreamHandler()
    ch.setFormatter(logging.Formatter("%(asctime)s %(levelname)s %(message)s"))
    logger.addHandler(fh)
    logger.addHandler(ch)
    return logger


async def load_rules(kv: KeyValue, name: str) -> list:
    """Load all rules for this subscriber from the KV bucket (live reload every event)."""
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
            "domain": action["domain"],
            "service": action["service"],
            "entity_id": action.get("entity_id", ""),
        }
        resp = await nc.request("cmd.ha.service.call", json.dumps(req).encode(), timeout=5)
        logger.info(f"ha_service {action['domain']}.{action['service']} → {resp.data.decode()}")

    elif atype == "nats_publish":
        pub_payload = action.get("payload", {})
        # Template substitution from incoming payload
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
        logger.info(f"kv_put {action.get('bucket', KV_CONTEXT_BUCKET)}.{action['key']} = {value}")

    elif atype == "telegram":
        msg = action["message"]
        try:
            msg = msg.format(subject=subject, **payload)
        except Exception:
            pass
        await nc.publish("cmd.telegram.send", json.dumps({"text": msg}).encode())
        logger.info(f"telegram → {msg[:80]}")

    elif atype == "prompt_agent":
        # Wake an agent with a prompt — publish to prompt.agent.<agent_name>
        # Agent responds on reply subject if request/reply used
        agent = action.get("agent", "legacy")
        prompt_text = action.get("prompt", "")
        try:
            prompt_text = prompt_text.format(subject=subject, **payload)
        except Exception:
            pass
        msg = {
            "prompt": prompt_text,
            "source_subject": subject,
            "payload": payload,
            "ts": datetime.utcnow().isoformat()
        }
        target_subject = f"prompt.agent.{agent}"
        if action.get("request_reply", False):
            try:
                resp = await nc.request(target_subject, json.dumps(msg).encode(), timeout=30)
                logger.info(f"prompt_agent {agent} reply: {resp.data.decode()[:200]}")
            except Exception as e:
                logger.error(f"prompt_agent {agent} timeout/error: {e}")
        else:
            await nc.publish(target_subject, json.dumps(msg).encode())
            logger.info(f"prompt_agent → {target_subject}: {prompt_text[:80]}")

    elif atype == "log":
        logger.info(f"[rule-log] {action.get('message', subject)}")

    else:
        logger.warning(f"Unknown action type: {atype}")


async def run(name: str, subject_prefix: str):
    logger = setup_logging(name)
    logger.info(f"Subscriber '{name}' starting — prefix={subject_prefix}, nats={NATS_URL}")

    async def error_cb(e):
        logger.error(f"NATS error: {e}")

    async def disconnected_cb():
        logger.warning("NATS disconnected")

    async def reconnected_cb():
        logger.info("NATS reconnected")

    nc = await nats.connect(
        NATS_URL,
        error_cb=error_cb,
        disconnected_cb=disconnected_cb,
        reconnected_cb=reconnected_cb,
        max_reconnect_attempts=-1,  # reconnect forever
    )
    js = nc.jetstream()

    # Get or create rules KV bucket
    try:
        rules_kv = await js.key_value(RULES_BUCKET)
        logger.info(f"Rules bucket '{RULES_BUCKET}' connected")
    except Exception:
        logger.warning(f"KV bucket '{RULES_BUCKET}' not found — run seed_rules.py first. Running in passthrough mode.")
        rules_kv = None

    async def message_handler(msg):
        await msg.ack()
        subject = msg.subject
        try:
            payload = json.loads(msg.data.decode()) if msg.data else {}
        except Exception:
            payload = {}

        logger.info(f"Event: {subject} {str(payload)[:120]}")

        # Reload rules on every event — live config, no restarts needed
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

        # Publish execution event for audit trail
        await nc.publish(
            f"events.subscriber.{name}.executed",
            json.dumps({
                "subject": subject,
                "rules_matched": len(matched),
                "ts": datetime.utcnow().isoformat()
            }).encode(),
        )

    # Subscribe with queue group — queue name IS the durable name in nats-py 2.x
    # Queue group also prevents double-fire when MBP + Atlas both run this subscriber
    await js.subscribe(
        subject_prefix,
        queue=f"sub-{name}",
        cb=message_handler,
        manual_ack=True,
    )
    logger.info(f"Subscribed to '{subject_prefix}' (queue=sub-{name})")
    logger.info(f"Subscriber '{name}' ready.")

    try:
        while True:
            await asyncio.sleep(1)
    except (asyncio.CancelledError, KeyboardInterrupt):
        pass
    finally:
        await nc.close()
        logger.info(f"Subscriber '{name}' stopped")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 subscriber.py <name> [subject_prefix]")
        print("  name:           home | security | health | comms | agents")
        print("  subject_prefix: e.g. events.zeus.>")
        sys.exit(1)
    name = sys.argv[1]
    prefix = sys.argv[2] if len(sys.argv) > 2 else f"events.{name}.>"
    asyncio.run(run(name, prefix))
