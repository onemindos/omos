#!/usr/bin/env python3
"""
OneMind bus heartbeat — emits live evt.> and tel.> traffic so the live viewer
always has something to show while we build. Publishes envelope-shaped events.

Run:  python3 heartbeat.py            (uses GLOBAL no-auth user 'onemind')
Stop: Ctrl-C  (emits evt.sys.node.offline first)

This is a DEV TOOL, not production. It publishes to:
  evt.sys.heartbeat.tick     every 3s  (envelope-shaped)
  tel.omos-mbp.cpu.load      every 3s  (live macOS load avg)
  tel.omos-mbp.mem.used_pct  every 3s  (live memory %)
"""
import json, time, uuid, subprocess, signal, sys, os
from datetime import datetime, timezone

try:
    import nats  # nats-py
except ImportError:
    print("nats-py not installed. Install: pip3 install nats-py", file=sys.stderr)
    sys.exit(1)

import asyncio

NODE = os.uname().nodename.split('.')[0]
URL = os.environ.get("OMD_NATS_URL", "nats://127.0.0.1:4222")

def now_iso():
    return datetime.now(timezone.utc).isoformat(timespec='milliseconds')

def envelope(typ, subject, data):
    """OneMind event envelope (P0-E shape)."""
    return json.dumps({
        "v": 1,
        "id": str(uuid.uuid4()),
        "ts": now_iso(),
        "type": typ,
        "source": f"heartbeat@{NODE}",
        "subject": subject,
        "data": data,
    })

def cpu_load():
    try:
        return os.getloadavg()[0]
    except Exception:
        return 0.0

def mem_used_pct():
    try:
        # macOS: vm_stat parse
        out = subprocess.run(["vm_stat"], capture_output=True, text=True).stdout
        pages = {}
        for line in out.splitlines():
            if ':' in line:
                k, v = line.split(':', 1)
                pages[k.strip()] = int(v.strip().rstrip('.') or 0)
        free = pages.get("Pages free", 0)
        active = pages.get("Pages active", 0)
        inactive = pages.get("Pages inactive", 0)
        wired = pages.get("Pages wired down", 0)
        total = free + active + inactive + wired
        used = active + wired
        return round(100 * used / total, 1) if total else 0.0
    except Exception:
        return 0.0

async def main():
    nc = await nats.connect(URL, name="onemind-heartbeat")
    print(f"[heartbeat] connected to {URL} as node={NODE}")
    # announce online
    await nc.publish("evt.sys.node.online",
                     envelope("evt.sys.node.online", "evt.sys.node.online",
                              {"node": NODE, "tool": "heartbeat"}).encode())

    stop = asyncio.Event()
    def _sig(*_): stop.set()
    for s in (signal.SIGINT, signal.SIGTERM):
        signal.signal(s, _sig)

    tick = 0
    while not stop.is_set():
        tick += 1
        ts = now_iso()
        # heartbeat event (envelope)
        await nc.publish("evt.sys.heartbeat.tick",
                         envelope("evt.sys.heartbeat.tick", "evt.sys.heartbeat.tick",
                                  {"node": NODE, "tick": tick}).encode())
        # telemetry (flat — telemetry is high-volume, kept lean)
        await nc.publish(f"tel.{NODE}.cpu.load",
                         json.dumps({"ts": ts, "node": NODE, "metric": "cpu.load", "value": cpu_load()}).encode())
        await nc.publish(f"tel.{NODE}.mem.used_pct",
                         json.dumps({"ts": ts, "node": NODE, "metric": "mem.used_pct", "value": mem_used_pct()}).encode())
        await nc.flush()
        try:
            await asyncio.wait_for(stop.wait(), timeout=3.0)
        except asyncio.TimeoutError:
            pass

    # announce offline
    await nc.publish("evt.sys.node.offline",
                     envelope("evt.sys.node.offline", "evt.sys.node.offline",
                              {"node": NODE, "tool": "heartbeat"}).encode())
    await nc.flush()
    await nc.drain()
    print("\n[heartbeat] stopped cleanly")

if __name__ == "__main__":
    asyncio.run(main())
