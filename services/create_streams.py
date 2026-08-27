"""
OneMind OS — Create/verify all JetStream event streams
Run once after NATS server start. Safe to re-run.
"""
import asyncio
import nats
from nats.js.api import StreamConfig, RetentionPolicy, StorageType

NATS_URL = "nats://127.0.0.1:4222"

THIRTY_DAYS = 30 * 24 * 3600
SEVEN_DAYS  = 7  * 24 * 3600

STREAMS = [
    StreamConfig(name="AGENT_EVENTS",      subjects=["agent.>"],                      max_age=THIRTY_DAYS, storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
    StreamConfig(name="ZEUS_EVENTS",       subjects=["events.zeus.>"],                max_age=THIRTY_DAYS, storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
    StreamConfig(name="SECURITY_EVENTS",   subjects=["events.security.>"],            max_age=THIRTY_DAYS, storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
    StreamConfig(name="HEALTH_EVENTS",     subjects=["events.health.>"],              max_age=THIRTY_DAYS, storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
    StreamConfig(name="COMMS",             subjects=["events.comms.>", "comms.>"],    max_age=THIRTY_DAYS, storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
    StreamConfig(name="COMMANDS",          subjects=["cmd.mode.>"],                   max_age=SEVEN_DAYS,  storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
    StreamConfig(name="AGENT_PROMPTS",     subjects=["prompt.agent.>"],               max_age=SEVEN_DAYS,  storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
    StreamConfig(name="SUBSCRIBER_EVENTS", subjects=["events.subscriber.>"],          max_age=SEVEN_DAYS,  storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
    StreamConfig(name="SYSTEM",            subjects=["system.>"],                     max_age=THIRTY_DAYS, storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
    StreamConfig(name="TASKS",             subjects=["task.>"],                       max_age=THIRTY_DAYS, storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
    StreamConfig(name="SAAS_EVENTS",       subjects=["saas.>"],                       max_age=THIRTY_DAYS, storage=StorageType.FILE, retention=RetentionPolicy.LIMITS),
]


async def main():
    nc = await nats.connect(NATS_URL)
    js = nc.jetstream()

    for cfg in STREAMS:
        try:
            info = await js.stream_info(cfg.name)
            print(f"  EXISTS:  {cfg.name} ({info.state.messages} msgs)")
        except Exception:
            try:
                await js.add_stream(cfg)
                print(f"  CREATED: {cfg.name} — {cfg.subjects}")
            except Exception as e:
                print(f"  ERROR:   {cfg.name} → {e}")

    print("\nDone.")
    await nc.close()


if __name__ == "__main__":
    asyncio.run(main())
