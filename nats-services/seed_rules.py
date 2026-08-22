"""
OneMind OS — Seed default rules into NATS KV bucket "rules"
Run once to bootstrap. Safe to re-run (overwrites with defaults).
"""
import asyncio
import json
import nats

NATS_URL = "nats://127.0.0.1:4222"
RULES_BUCKET = "rules"


DEFAULT_RULES = {

    # =========================================================================
    # HOME subscriber — watches events.zeus.>
    # =========================================================================
    "home/left_home": {
        "subject": "events.zeus.left_home",
        "actions": [
            {"type": "kv_put", "bucket": "zeus-context", "key": "location", "value": "away"},
            {"type": "ha_service", "domain": "switch", "service": "turn_off", "entity_id": "switch.plug_mini_a886"},
            {"type": "ha_service", "domain": "switch", "service": "turn_off", "entity_id": "switch.plug_mini_ef5a"},
            {"type": "telegram", "message": "Zeus left home. Plugs off, context updated."},
        ]
    },
    "home/arrived_home": {
        "subject": "events.zeus.arrived_home",
        "actions": [
            {"type": "kv_put", "bucket": "zeus-context", "key": "location", "value": "home-office"},
            {"type": "telegram", "message": "Zeus arrived home. Welcome back."},
        ]
    },
    "home/focus_deep_work": {
        "subject": "events.zeus.focus.deep_work",
        "actions": [
            {"type": "kv_put", "bucket": "zeus-context", "key": "focus_mode", "value": "deep-work"},
            {"type": "log", "message": "Focus mode: deep-work activated"},
        ]
    },
    "home/focus_rest": {
        "subject": "events.zeus.focus.rest",
        "actions": [
            {"type": "kv_put", "bucket": "zeus-context", "key": "focus_mode", "value": "rest"},
            {"type": "log", "message": "Focus mode: rest activated"},
        ]
    },

    # =========================================================================
    # SECURITY subscriber — watches events.security.>
    # =========================================================================
    "security/motion_night": {
        "subject": "events.security.motion.detected",
        "actions": [
            {"type": "prompt_agent", "agent": "legacy",
             "prompt": "Motion detected at {ts}. Zeus is home. Assess threat level and decide: alert or ignore.",
             "request_reply": False},
            {"type": "telegram", "message": "Motion detected. Legacy notified for assessment."},
        ]
    },
    "security/high_alert_activate": {
        "subject": "cmd.mode.security.alert",
        "actions": [
            {"type": "kv_put", "bucket": "zeus-context", "key": "focus_mode", "value": "unavailable"},
            {"type": "telegram", "message": "HIGH ALERT mode activated. All sensors active."},
            {"type": "prompt_agent", "agent": "legacy",
             "prompt": "High alert mode activated. Initiate security protocol: check all sensors, lock all available devices, monitor for 30 minutes."},
        ]
    },

    # =========================================================================
    # HEALTH subscriber — watches events.health.>
    # =========================================================================
    "health/battery_low": {
        "subject": "events.health.battery.low",
        "actions": [
            {"type": "telegram", "message": "Zeus iPhone battery low. Plug in to charge."},
        ]
    },
    "health/battery_critical": {
        "subject": "events.health.battery.critical",
        "actions": [
            {"type": "prompt_agent", "agent": "legacy",
             "prompt": "Zeus iPhone battery critical (<10%). Send urgent notification and log event."},
            {"type": "telegram", "message": "CRITICAL: Zeus iPhone battery critical. Charge immediately."},
        ]
    },

    # =========================================================================
    # COMMS subscriber — watches events.comms.> and cmd.mode.*
    # =========================================================================
    "comms/incognito_activate": {
        "subject": "cmd.mode.incognito",
        "actions": [
            {"type": "kv_put", "bucket": "zeus-context", "key": "status", "value": "offline"},
            {"type": "kv_put", "bucket": "zeus-context", "key": "focus_mode", "value": "unavailable"},
            {"type": "log", "message": "Incognito mode activated — context set offline"},
        ]
    },
    "comms/watch_me_activate": {
        "subject": "cmd.mode.watch_me",
        "actions": [
            {"type": "kv_put", "bucket": "zeus-context", "key": "status", "value": "online"},
            {"type": "prompt_agent", "agent": "legacy",
             "prompt": "Watch Me mode activated. Start monitoring Zeus's iPhone location, battery, and activity. Report any notable changes to Telegram every 15 minutes."},
            {"type": "telegram", "message": "Watch Me mode active. Legacy is monitoring."},
        ]
    },
    "comms/default_mode": {
        "subject": "cmd.mode.default",
        "actions": [
            {"type": "kv_put", "bucket": "zeus-context", "key": "status", "value": "online"},
            {"type": "kv_put", "bucket": "zeus-context", "key": "focus_mode", "value": "active"},
            {"type": "log", "message": "Default mode restored"},
        ]
    },

    # =========================================================================
    # AGENTS subscriber — watches prompt.agent.> (agent invocation bus)
    # =========================================================================
    "agents/legacy_ping": {
        "subject": "prompt.agent.legacy.ping",
        "actions": [
            {"type": "log", "message": "Legacy ping received"},
            {"type": "nats_publish", "subject": "events.agent.legacy.pong",
             "payload": {"status": "alive", "ts": "now"}},
        ]
    },
    "agents/broadcast": {
        "subject": "prompt.agent.broadcast",
        "actions": [
            {"type": "log", "message": "Broadcast received on all agents"},
        ]
    },
}


async def seed():
    nc = await nats.connect(NATS_URL)
    js = nc.jetstream()

    # Create or bind rules bucket
    try:
        kv = await js.key_value(RULES_BUCKET)
        print(f"Bucket '{RULES_BUCKET}' already exists — updating rules")
    except Exception:
        from nats.js.api import KeyValueConfig
        kv = await js.create_key_value(KeyValueConfig(bucket=RULES_BUCKET, history=5))
        print(f"Created bucket '{RULES_BUCKET}'")

    for key, rule in DEFAULT_RULES.items():
        await kv.put(key, json.dumps(rule).encode())
        print(f"  seeded: {key} → subject={rule['subject']} actions={len(rule['actions'])}")

    print(f"\nDone. {len(DEFAULT_RULES)} rules seeded into '{RULES_BUCKET}' bucket.")
    await nc.close()


if __name__ == "__main__":
    asyncio.run(seed())
