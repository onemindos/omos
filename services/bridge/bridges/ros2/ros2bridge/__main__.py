"""bridge-ros2: connects one simulated ROS 2 robot to the OneMind bus.

Run: ``python -m ros2bridge`` (configuration from the environment, see README.md).
"""

from __future__ import annotations

import asyncio
import logging
import signal
import sys

from . import config
from .bridge import Bridge, Bus

log = logging.getLogger("ros2bridge")


class NatsBus(Bus):
    def __init__(self, nc):
        self.nc = nc

    async def publish(self, subject: str, data: bytes, msg_id: str = "") -> None:
        await self.nc.publish(subject, data, headers={"Nats-Msg-Id": msg_id} if msg_id else None)


def make_transport(s: config.Settings):
    if s.transport == "zenoh":
        from .zenoh_transport import ZenohTransport

        return ZenohTransport(s.zenoh_connect, s.zenoh_mode, s.zenoh_listen, s.namespace, s.topics, s.query_timeout)
    from .rosbridge_transport import RosbridgeTransport

    return RosbridgeTransport(s.rosbridge_url, s.namespace, s.topics)


async def run(s: config.Settings, stop: asyncio.Event, ready=None) -> None:
    import nats

    cfg = s.bridge
    opts = dict(name=f"{cfg.source}/{cfg.asset_id}", max_reconnect_attempts=-1, reconnect_time_wait=2)
    if s.nats_creds:
        opts["user_credentials"] = s.nats_creds
    elif s.nats_user:
        opts["user"], opts["password"] = s.nats_user, s.nats_password

    async def disconnected():
        log.warning("[%s] NATS disconnected", cfg.source)

    async def reconnected():
        log.info("[%s] NATS reconnected", cfg.source)

    async def error_cb(e):
        log.warning("[%s] NATS error: %s", cfg.source, e)

    nc = await nats.connect(s.nats_url, disconnected_cb=disconnected, reconnected_cb=reconnected,
                            error_cb=error_cb, **opts)
    log.info("[%s] NATS connected: %s", cfg.source, nc.connected_url.netloc if nc.connected_url else s.nats_url)
    transport = make_transport(s)
    b = Bridge(cfg, NatsBus(nc), transport)
    try:
        # Core subscription, not a durable consumer: an offline bridge gets no stale command
        # later (asset-command.md section 4 "No queued replay").
        async def on_cmd(m):
            b._spawn(b.handle_command(m.subject, m.data))

        sub = await nc.subscribe(b.command_subject(), cb=on_cmd)
        await nc.flush()
        await transport.start(b.on_odom, b.on_pose, b.on_battery)
        log.info("[%s] subscribed %s; publishing %s, %s, %s via %s", cfg.source, b.command_subject(),
                 b.position_subject(), b.status_subject(), b.outcome_subject(), transport.name)
        if ready:
            ready(b)
        await b.run_status(stop)
        await sub.unsubscribe()
    finally:
        await transport.close()
        await nc.drain()
        log.info("[%s] stopped", cfg.source)


def main() -> int:
    logging.basicConfig(stream=sys.stderr, level=logging.INFO,
                        format="%(asctime)s.%(msecs)03d %(levelname)s %(name)s %(message)s", datefmt="%Y-%m-%dT%H:%M:%S")
    try:
        s = config.load()
    except ValueError as e:
        log.error("[bridge-ros2] config: %s", e)
        return 2

    async def amain():
        stop = asyncio.Event()
        loop = asyncio.get_running_loop()
        for sig in (signal.SIGINT, signal.SIGTERM):
            loop.add_signal_handler(sig, stop.set)
        await run(s, stop)

    asyncio.run(amain())
    return 0


if __name__ == "__main__":
    sys.exit(main())
