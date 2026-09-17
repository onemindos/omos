"""
OneMind OS — MAVLink Bridge Skeleton (Drone)
Pattern: NATS cmd → MAVLink command to drone → NATS tel/evt back

Subject law:
  cmd.drone.<id>.<action>     ← Legacy publishes here
  tel.drone.<id>.<reading>    ← bridge publishes telemetry here (fast, no JetStream)
  evt.drone.<id>.<event>      ← bridge publishes state changes here (JetStream)

One bridge process per drone family (PX4/ArduPilot).
Adding a new drone = register its ID in the node registry, no new code.

Requirements:
  pip install pymavlink nats-py

IMPORTANT: NATS is NOT the flight controller.
  - PX4/ArduPilot still runs the 400Hz inner loop natively
  - NATS handles: launch, land, goto, mission, return, telemetry reporting
  - Never pipe raw IMU/servo data through NATS — that stays internal to the autopilot
"""

import asyncio
import json
import threading
import nats
from pymavlink import mavutil

NATS_URL   = "nats://localhost:4222"
DRONE_ID   = "skyhawk-1"                        # from node registry
DRONE_CONN = "udpin:0.0.0.0:14550"             # MAVLink connection string
                                                 # UDP = simulator or companion computer
                                                 # Serial = "/dev/ttyUSB0,57600"

# ── MAVLink connection (runs in background thread) ────────────────────────────

mav = None

def connect_drone():
    global mav
    print(f"[mavlink-bridge] connecting to drone at {DRONE_CONN}")
    mav = mavutil.mavlink_connection(DRONE_CONN)
    mav.wait_heartbeat()
    print(f"[mavlink-bridge] heartbeat received — drone online")


# ── command handlers (NATS → MAVLink) ────────────────────────────────────────

async def handle_takeoff(nc, msg):
    """
    NATS subject: cmd.drone.skyhawk-1.takeoff
    Payload: { "payload": { "altitude_m": 10 } }
    """
    body     = json.loads(msg.data)
    altitude = body.get("payload", {}).get("altitude_m", 10)

    # arm the drone
    mav.arducopter_arm()
    mav.motors_armed_wait()

    # send takeoff command
    mav.mav.command_long_send(
        mav.target_system,
        mav.target_component,
        mavutil.mavlink.MAV_CMD_NAV_TAKEOFF,
        0,          # confirmation
        0, 0, 0, 0, # unused params
        0, 0,       # lat/lon (0 = current position)
        altitude,
    )

    # publish event back to bus
    await nc.publish(
        f"evt.drone.{DRONE_ID}.takeoff_initiated",
        json.dumps({"ok": True, "altitude_m": altitude}).encode()
    )
    if msg.reply:
        await nc.publish(msg.reply,
            json.dumps({"ok": True, "altitude_m": altitude}).encode())


async def handle_land(nc, msg):
    """
    NATS subject: cmd.drone.skyhawk-1.land
    """
    mav.mav.command_long_send(
        mav.target_system, mav.target_component,
        mavutil.mavlink.MAV_CMD_NAV_LAND,
        0, 0, 0, 0, 0, 0, 0, 0
    )
    await nc.publish(
        f"evt.drone.{DRONE_ID}.landing",
        json.dumps({"ok": True}).encode()
    )
    if msg.reply:
        await nc.publish(msg.reply, json.dumps({"ok": True}).encode())


async def handle_return_home(nc, msg):
    """
    NATS subject: cmd.drone.skyhawk-1.return_home
    """
    mav.mav.command_long_send(
        mav.target_system, mav.target_component,
        mavutil.mavlink.MAV_CMD_NAV_RETURN_TO_LAUNCH,
        0, 0, 0, 0, 0, 0, 0, 0
    )
    await nc.publish(
        f"evt.drone.{DRONE_ID}.returning_home",
        json.dumps({"ok": True}).encode()
    )


async def handle_goto(nc, msg):
    """
    NATS subject: cmd.drone.skyhawk-1.goto
    Payload: { "payload": { "lat": 41.0, "lon": -73.5, "alt_m": 30 } }
    """
    body = json.loads(msg.data)
    p    = body.get("payload", {})
    lat  = p.get("lat", 0)
    lon  = p.get("lon", 0)
    alt  = p.get("alt_m", 20)

    mav.mav.mission_item_int_send(
        mav.target_system, mav.target_component,
        0,                                          # seq
        mavutil.mavlink.MAV_FRAME_GLOBAL_RELATIVE_ALT_INT,
        mavutil.mavlink.MAV_CMD_NAV_WAYPOINT,
        2,                                          # current = guided mode
        1,                                          # autocontinue
        0, 0, 0, 0,
        int(lat * 1e7), int(lon * 1e7), alt
    )
    await nc.publish(
        f"evt.drone.{DRONE_ID}.waypoint_set",
        json.dumps({"ok": True, "lat": lat, "lon": lon, "alt_m": alt}).encode()
    )


# ── telemetry publisher (MAVLink → NATS) ─────────────────────────────────────

async def publish_telemetry(nc):
    """
    Reads MAVLink telemetry from drone continuously.
    Publishes to tel.drone.<id>.* subjects.
    These are fire-and-forget (no JetStream) — high frequency.
    Sampling here: publish every message as it arrives.
    """
    loop = asyncio.get_event_loop()
    while True:
        msg = await loop.run_in_executor(
            None, lambda: mav.recv_match(blocking=True, timeout=1.0)
        )
        if msg is None:
            continue

        msg_type = msg.get_type()

        if msg_type == "GLOBAL_POSITION_INT":
            await nc.publish(
                f"tel.drone.{DRONE_ID}.gps",
                json.dumps({
                    "lat":  msg.lat  / 1e7,
                    "lon":  msg.lon  / 1e7,
                    "alt_m": msg.relative_alt / 1000.0,
                    "hdg":  msg.hdg / 100.0,
                }).encode()
            )

        elif msg_type == "BATTERY_STATUS":
            await nc.publish(
                f"tel.drone.{DRONE_ID}.battery",
                json.dumps({
                    "voltage_v":   msg.voltages[0] / 1000.0 if msg.voltages else 0,
                    "remaining_pct": msg.battery_remaining,
                }).encode()
            )

        elif msg_type == "HEARTBEAT":
            mode = mavutil.mavlink.enums["MAV_MODE_FLAG"].get(msg.base_mode, {})
            await nc.publish(
                f"tel.drone.{DRONE_ID}.heartbeat",
                json.dumps({
                    "armed":  bool(msg.base_mode & mavutil.mavlink.MAV_MODE_FLAG_SAFETY_ARMED),
                    "flight_mode": str(msg.custom_mode),
                    "system_status": str(msg.system_status),
                }).encode()
            )


# ── main ─────────────────────────────────────────────────────────────────────

async def main():
    # connect to drone first (blocking — must have MAVLink before subscribing)
    connect_drone()

    nc = await nats.connect(NATS_URL)
    print("[mavlink-bridge] connected to NATS")

    subject_prefix = f"cmd.drone.{DRONE_ID}"

    await nc.subscribe(f"{subject_prefix}.takeoff",
        cb=lambda m: asyncio.create_task(handle_takeoff(nc, m)))
    await nc.subscribe(f"{subject_prefix}.land",
        cb=lambda m: asyncio.create_task(handle_land(nc, m)))
    await nc.subscribe(f"{subject_prefix}.return_home",
        cb=lambda m: asyncio.create_task(handle_return_home(nc, m)))
    await nc.subscribe(f"{subject_prefix}.goto",
        cb=lambda m: asyncio.create_task(handle_goto(nc, m)))

    print(f"[mavlink-bridge] subscribed to {subject_prefix}.* — bridge live")

    # run telemetry publisher concurrently
    await asyncio.gather(
        publish_telemetry(nc),
        asyncio.Event().wait(),   # keep main loop alive
    )


if __name__ == "__main__":
    asyncio.run(main())


# ── end-to-end flow ───────────────────────────────────────────────────────────
#
#  Legacy publishes:
#    cmd.drone.skyhawk-1.takeoff
#    { "id": "ulid", "source": "agent.legacy", "payload": { "altitude_m": 15 } }
#
#  MAVLink bridge hears it
#  Bridge arms drone via MAVLink
#  Bridge sends MAV_CMD_NAV_TAKEOFF via MAVLink UDP
#  Drone takes off (PX4/ArduPilot handles the inner 400Hz loop)
#
#  Bridge continuously publishes back:
#    tel.drone.skyhawk-1.gps       { lat, lon, alt_m, hdg }
#    tel.drone.skyhawk-1.battery   { voltage_v, remaining_pct }
#    tel.drone.skyhawk-1.heartbeat { armed, flight_mode }
#
#  Any NATS subscriber (dashboard, safety monitor, legacy) reacts in real time.
#
#  The drone still speaks MAVLink internally.
#  NATS is only the high-level command + telemetry layer.
