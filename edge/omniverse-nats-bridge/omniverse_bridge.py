"""
OneMind Omniverse NATS Bridge
Subscribes to entity state updates on NATS and mirrors them as USD prims in
an Omniverse/Isaac Sim scene.  Runs as an Omniverse Extension inside Isaac Sim.

NATS subjects consumed:
  ent.*.*.state   — all entity state updates (robot, drone, person, …)
  tel.robot.>.    — robot telemetry (position, status)
  tel.drone.>.    — drone telemetry

USD prim layout:
  /World/Entities/{type}/{id}   — XformPrim, updated each message

Requirements:
  - NVIDIA Isaac Sim 2023.1+ with Omniverse Kit
  - nats-py installed in the Kit Python env
  - GPU node (RTX 3080+) — see docs/gpu-node-requirement.md

Environment variables:
  NATS_URL      — default nats://localhost:4222 (leaf sidecar)
  NATS_USER / NATS_PASS — optional credentials
  STAGE_PATH    — USD stage to update (default: current open stage)
"""

import asyncio
import json
import math
import os
import threading

import nats
import omni.usd
from pxr import Gf, UsdGeom, Sdf

NATS_URL   = os.environ.get("NATS_URL",  "nats://localhost:4222")
NATS_USER  = os.environ.get("NATS_USER", "")
NATS_PASS  = os.environ.get("NATS_PASS", "")

_stage     = None
_nc        = None
_loop      = None


def _get_or_create_prim(entity_type: str, entity_id: str):
    prim_path = f"/World/Entities/{entity_type}/{entity_id}"
    stage = omni.usd.get_context().get_stage()
    prim = stage.GetPrimAtPath(prim_path)
    if not prim.IsValid():
        prim = UsdGeom.Xform.Define(stage, prim_path).GetPrim()
    return prim


def _lat_lon_alt_to_xyz(lat: float, lon: float, alt: float = 0.0):
    # Simple flat-earth approximation centred at (0,0).
    # Replace with CesiumGeoreference math for large-area scenes.
    R = 6_371_000.0
    x = math.radians(lon) * R
    y = math.radians(lat) * R
    z = alt
    return Gf.Vec3d(x, y, z)


def _update_prim_from_state(data: dict):
    entity_type = data.get("type", "unknown")
    entity_id   = data.get("id",   "unknown")
    lat = data.get("lat", 0.0)
    lon = data.get("lon", 0.0)
    alt = data.get("alt", 0.0)

    prim = _get_or_create_prim(entity_type, entity_id)
    xform = UsdGeom.Xformable(prim)

    ops = {op.GetOpName(): op for op in xform.GetOrderedXformOps()}
    if "xformOp:translate" not in ops:
        xform.AddTranslateOp()

    xform.GetOrderedXformOps()[0].Set(_lat_lon_alt_to_xyz(lat, lon, alt))

    meta_attr = prim.GetAttribute("omni:meta:callsign")
    if not meta_attr:
        meta_attr = prim.CreateAttribute("omni:meta:callsign", Sdf.ValueTypeNames.String)
    meta_attr.Set(data.get("callsign", entity_id))


async def _msg_handler(msg):
    try:
        data = json.loads(msg.data.decode())
        _update_prim_from_state(data)
    except Exception as e:
        print(f"[omniverse-bridge] handler error: {e}")


async def _nats_loop():
    global _nc
    opts = {"servers": [NATS_URL]}
    if NATS_USER:
        opts["user"] = NATS_USER
        opts["password"] = NATS_PASS

    _nc = await nats.connect(**opts)
    print(f"[omniverse-bridge] connected to {NATS_URL}")

    await _nc.subscribe("ent.*.*.state", cb=_msg_handler)
    await _nc.subscribe("tel.robot.>",   cb=_msg_handler)
    await _nc.subscribe("tel.drone.>",   cb=_msg_handler)

    print("[omniverse-bridge] subscribed — streaming entities to USD stage")
    await asyncio.Event().wait()


def start():
    global _loop
    _loop = asyncio.new_event_loop()
    t = threading.Thread(target=_loop.run_until_complete, args=(_nats_loop(),), daemon=True)
    t.start()
    print("[omniverse-bridge] started")


def stop():
    if _nc:
        asyncio.run_coroutine_threadsafe(_nc.close(), _loop)


if __name__ == "__main__":
    start()
    threading.Event().wait()
