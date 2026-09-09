"""
MAVLink → NATS Bridge — Phase 2 keystone.
Streams PX4 SITL telemetry into NATS JetStream.

Subject: rad.<homestead_id>.robots.<drone_id>.telemetry
Command: rad.<homestead_id>.robots.<drone_id>.cmd

Usage:
    python -m rads.bridges.mavlink_bridge
    MAVLINK_ADDRESS=udp://:14540 python -m rads.bridges.mavlink_bridge

PX4 SITL setup:
    docker run -d --name px4-drone \
      jonasvautherin/px4-gazebo-headless:1.14.0 \
      <YOUR_MAC_TAILSCALE_IP> <YOUR_MAC_TAILSCALE_IP>
"""
import asyncio, json, time
import nats
from rads.config import (NATS_URL, HOMESTEAD_ID, MAVLINK_ADDRESS,
                         drone_telemetry_subject, drone_cmd_subject)
from rads.envelope import wrap, unwrap

DRONE_ID = "drone1"

async def run_bridge():
    try:
        from mavsdk import System
    except ImportError:
        print("[mavlink-bridge] mavsdk not installed — run: pip install mavsdk", flush=True)
        return

    nc = await nats.connect(NATS_URL)
    js = nc.jetstream()

    # ensure stream
    try:
        await nc.jsm().add_stream(
            name="RADS_ROBOTS",
            subjects=[f"rad.{HOMESTEAD_ID}.robots.>"],
            retention="limits",
            max_age=10 * 365 * 24 * 3600,
        )
    except Exception:
        pass

    drone = System()
    print(f"[mavlink-bridge] connecting to {MAVLINK_ADDRESS} ...", flush=True)
    await drone.connect(system_address=MAVLINK_ADDRESS)

    async for state in drone.core.connection_state():
        if state.is_connected:
            print(f"[mavlink-bridge] drone connected", flush=True)
            break

    # incoming commands from NATS → drone
    async def handle_cmd(msg):
        env = unwrap(msg.data)
        cmd = env["payload"]
        print(f"[mavlink-bridge] CMD → {cmd}", flush=True)
        # extend here: arm, takeoff, goto, land, etc.

    await nc.subscribe(drone_cmd_subject(DRONE_ID), cb=handle_cmd)

    # telemetry → NATS
    telemetry_subj = drone_telemetry_subject(DRONE_ID)
    async for pos in drone.telemetry.position():
        payload = {
            "drone_id":    DRONE_ID,
            "lat":         pos.latitude_deg,
            "lon":         pos.longitude_deg,
            "alt_m":       pos.absolute_altitude_m,
            "alt_rel_m":   pos.relative_altitude_m,
            "ts":          time.time(),
        }
        await js.publish(telemetry_subj, wrap("drone.telemetry", payload, source=DRONE_ID))
        print(f"[mavlink-bridge] {telemetry_subj} → lat={pos.latitude_deg:.5f} lon={pos.longitude_deg:.5f} alt={pos.relative_altitude_m:.1f}m", flush=True)

if __name__ == "__main__":
    asyncio.run(run_bridge())
