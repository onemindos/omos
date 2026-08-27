"""
OneMind OS — Create/seed all KV buckets
Run after NATS server start or domain change.
"""
import asyncio
import nats
from nats.js.api import KeyValueConfig

NATS_URL = "nats://127.0.0.1:4222"

BUCKETS = {
    "zeus-context": {
        "history": 20,
        "initial": {
            "location": "home-office",
            "focus_mode": "active",
            "status": "online",
        }
    },
    "agent-state": {
        "history": 5,
        "initial": {
            "legacy.status": "idle",
            "hermes.status": "idle",
        }
    },
    "config": {
        "history": 5,
        "initial": {
            "nats.version": "2.14.1",
            "setup.date": "2026-05-22",
        }
    },
    "rules": {
        "history": 5,
        "initial": {}  # seeded by seed_rules.py
    },
    "connections": {"history": 5, "initial": {}},
    "schedules":   {"history": 5, "initial": {}},
    "triggers":    {"history": 5, "initial": {}},
}


async def main():
    nc = await nats.connect(NATS_URL)
    js = nc.jetstream()

    for bucket_name, spec in BUCKETS.items():
        try:
            kv = await js.key_value(bucket_name)
            print(f"  EXISTS:  {bucket_name}")
        except Exception:
            kv = await js.create_key_value(
                KeyValueConfig(bucket=bucket_name, history=spec["history"])
            )
            print(f"  CREATED: {bucket_name}")

        for key, value in spec.get("initial", {}).items():
            try:
                existing = await kv.get(key)
                print(f"    {key} = {existing.value.decode()} (kept)")
            except Exception:
                await kv.put(key, value.encode())
                print(f"    {key} = {value} (seeded)")

    print("\nDone.")
    await nc.close()


if __name__ == "__main__":
    asyncio.run(main())
