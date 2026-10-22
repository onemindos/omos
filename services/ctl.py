"""
OneMind OS — Subscriber Control CLI
=====================================
Usage:
    python3 ctl.py status
    python3 ctl.py test <subject> '<json_payload>'
    python3 ctl.py rules show [subscriber_name]
    python3 ctl.py rules set <key> '<json_rule>'
    python3 ctl.py rules delete <key>
    python3 ctl.py mode <incognito|high_alert|watch_me|default>
    python3 ctl.py agent <legacy|codex> '<prompt text>'
    python3 ctl.py context
"""

import asyncio
import json
import sys

import nats

NATS_URL = "nats://127.0.0.1:4222"
RULES_BUCKET = "rules"


async def get_status(nc, js):
    print(f"NATS connected: {nc.connected_url}\n")
    # Show agent state
    try:
        kv = await js.key_value("agent-state")
        keys = await kv.keys()
        print("Agent state:")
        for k in sorted(keys):
            e = await kv.get(k)
            print(f"  {k} = {e.value.decode()}")
    except Exception as e:
        print(f"  agent-state: {e}")
    print()
    # Show zeus context
    try:
        kv = await js.key_value("zeus-context")
        keys = await kv.keys()
        print("Zeus context:")
        for k in sorted(keys):
            e = await kv.get(k)
            print(f"  {k} = {e.value.decode()}")
    except Exception as e:
        print(f"  zeus-context: {e}")
    print()
    # Show running subscribers via launchctl
    import subprocess
    result = subprocess.run(["launchctl", "list"], capture_output=True, text=True)
    subs = [l for l in result.stdout.splitlines() if "onemind.subscriber" in l]
    print("Subscriber services:")
    if subs:
        for s in subs:
            print(f"  {s}")
    else:
        print("  (none loaded)")


async def test_event(nc, subject, payload_str):
    payload = json.loads(payload_str) if payload_str and payload_str != "{}" else {}
    await nc.publish(subject, json.dumps(payload).encode())
    print(f"Published: {subject}")
    print(f"Payload:   {json.dumps(payload)}")


async def show_rules(js, name_filter=None):
    kv = await js.key_value(RULES_BUCKET)
    try:
        keys = await kv.keys()
    except Exception:
        print("No rules found (bucket may be empty)")
        return
    found = 0
    for k in sorted(keys):
        if name_filter and not k.startswith(name_filter + "/"):
            continue
        e = await kv.get(k)
        rule = json.loads(e.value.decode())
        print(f"\n{k}:")
        print(f"  subject: {rule.get('subject')}")
        for a in rule.get("actions", []):
            atype = a.get("type")
            if atype == "ha_service":
                print(f"  action:  ha_service {a.get('domain')}.{a.get('service')} → {a.get('entity_id')}")
            elif atype == "telegram":
                print(f"  action:  telegram → {a.get('message','')[:60]}")
            elif atype == "prompt_agent":
                print(f"  action:  prompt_agent → {a.get('agent')} | {a.get('prompt','')[:60]}")
            elif atype == "kv_put":
                print(f"  action:  kv_put {a.get('bucket')}.{a.get('key')} = {a.get('value')}")
            else:
                print(f"  action:  {a}")
        found += 1
    if found == 0:
        print("No rules found" + (f" for '{name_filter}'" if name_filter else ""))
    else:
        print(f"\nTotal: {found} rules")


async def set_rule(js, key, rule_str):
    kv = await js.key_value(RULES_BUCKET)
    rule = json.loads(rule_str)
    await kv.put(key, json.dumps(rule).encode())
    print(f"Set: {key} → subject={rule.get('subject')}")


async def delete_rule(js, key):
    kv = await js.key_value(RULES_BUCKET)
    await kv.delete(key)
    print(f"Deleted: {key}")


async def activate_mode(nc, mode):
    subject_map = {
        "incognito":  "cmd.mode.incognito",
        "high_alert": "cmd.mode.security.alert",
        "watch_me":   "cmd.mode.watch_me",
        "default":    "cmd.mode.default",
    }
    subject = subject_map.get(mode)
    if not subject:
        print(f"Unknown mode: {mode}")
        print(f"Options: {list(subject_map)}")
        return
    await nc.publish(subject, b"{}")
    print(f"Mode activated: {mode} → {subject}")


async def prompt_agent(nc, agent, prompt_text):
    msg = {
        "prompt": prompt_text,
        "source": "ctl",
        "ts": __import__("datetime").datetime.utcnow().isoformat()
    }
    subject = f"prompt.agent.{agent}"
    await nc.publish(subject, json.dumps(msg).encode())
    print(f"Prompted agent '{agent}' via {subject}")
    print(f"Prompt: {prompt_text}")


async def show_context(js):
    for bucket in ["zeus-context", "agent-state", "config"]:
        try:
            kv = await js.key_value(bucket)
            keys = await kv.keys()
            print(f"\n[{bucket}]")
            for k in sorted(keys):
                e = await kv.get(k)
                print(f"  {k} = {e.value.decode()}")
        except Exception as e:
            print(f"\n[{bucket}] error: {e}")


async def main():
    if len(sys.argv) < 2:
        print(__doc__)
        return

    nc = await nats.connect(NATS_URL)
    js = nc.jetstream()

    cmd = sys.argv[1]

    try:
        if cmd == "status":
            await get_status(nc, js)

        elif cmd == "test":
            subject = sys.argv[2] if len(sys.argv) > 2 else "events.test"
            payload = sys.argv[3] if len(sys.argv) > 3 else "{}"
            await test_event(nc, subject, payload)

        elif cmd == "rules":
            sub = sys.argv[2] if len(sys.argv) > 2 else "show"
            if sub == "show":
                name_filter = sys.argv[3] if len(sys.argv) > 3 else None
                await show_rules(js, name_filter)
            elif sub == "set":
                await set_rule(js, sys.argv[3], sys.argv[4])
            elif sub == "delete":
                await delete_rule(js, sys.argv[3])

        elif cmd == "mode":
            mode = sys.argv[2] if len(sys.argv) > 2 else "default"
            await activate_mode(nc, mode)

        elif cmd == "agent":
            agent = sys.argv[2] if len(sys.argv) > 2 else "legacy"
            prompt_text = sys.argv[3] if len(sys.argv) > 3 else "ping"
            await prompt_agent(nc, agent, prompt_text)

        elif cmd == "context":
            await show_context(js)

        else:
            print(f"Unknown command: {cmd}")
            print(__doc__)

    finally:
        await nc.close()


if __name__ == "__main__":
    asyncio.run(main())
