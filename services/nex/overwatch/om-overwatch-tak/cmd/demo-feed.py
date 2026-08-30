#!/usr/bin/env python3
"""
demo-feed — live outbound-test feed for om-overwatch-tak.

Drives REAL fusion (publishes to the spine's ingest subjects, not directly to
ent.air), so the full chain is exercised:

    this feed → spine fuses → ent.air.{id}.state → om-overwatch-tak → tak.cot.inject
              → NatsBridgePlugin → every connected ATAK client

Produces two moving tracks so an operator immediately sees the fused picture:
  • DEMO-PATROL  — friendly multirotor circling the center point   (blue)
  • DEMO-THREAT  — hostile air track inbound toward the center     (red)

Run it once a device is connected; tracks appear and move live. Ctrl-C to stop.
Cleanup:  nats kv del OVERWATCH_STATE DEMO-PATROL ; nats kv del OVERWATCH_STATE DEMO-THREAT

Env:
  NATS_URL     default nats://127.0.0.1:4222
  CENTER_LAT   default 25.7617   (Miami — matches the plugin probe example)
  CENTER_LON   default -80.1918
  TICK_SECS    default 2
"""
import json
import math
import os
import subprocess
import sys
import time

NATS_URL = os.getenv("NATS_URL", "nats://127.0.0.1:4222")
CLAT = float(os.getenv("CENTER_LAT", "25.7617"))
CLON = float(os.getenv("CENTER_LON", "-80.1918"))
TICK = float(os.getenv("TICK_SECS", "2"))


def pub(subject: str, payload: dict):
    """Publish JSON to a subject via the nats CLI."""
    subprocess.run(
        ["nats", "pub", subject, json.dumps(payload), "--server", NATS_URL],
        check=False,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )


def bearing(from_lat, from_lon, to_lat, to_lon) -> float:
    """Initial great-circle bearing in degrees (0=N, 90=E)."""
    dlon = math.radians(to_lon - from_lon)
    y = math.sin(dlon) * math.cos(math.radians(to_lat))
    x = math.cos(math.radians(from_lat)) * math.sin(math.radians(to_lat)) - math.sin(
        math.radians(from_lat)
    ) * math.cos(math.radians(to_lat)) * math.cos(dlon)
    return (math.degrees(math.atan2(y, x)) + 360) % 360


def main():
    print(f"demo-feed → {NATS_URL}  center=({CLAT},{CLON})  tick={TICK}s")
    print("  DEMO-PATROL (friendly, circling)  +  DEMO-THREAT (hostile, inbound)")
    print("  Ctrl-C to stop.  Cleanup: nats kv del OVERWATCH_STATE DEMO-PATROL / DEMO-THREAT")

    patrol_r = 0.012      # ~1.3 km orbit radius (deg)
    patrol_a = 0.0        # orbit angle (rad)
    threat_start = (0.09, 0.06)  # NE offset from center (deg lat, lon)
    threat_t = 0.0        # 0..1 progress inbound

    tick = 0
    while True:
        # ── friendly patrol: circle the center ──
        patrol_a += math.radians(12)
        plat = CLAT + patrol_r * math.sin(patrol_a)
        plon = CLON + patrol_r * math.cos(patrol_a)
        phdg = (math.degrees(patrol_a) + 90) % 360  # tangent to circle
        # Flat telemetry payload (spine.extractPos reads lat/lon in degrees).
        pub("tel.drone.DEMO-PATROL.mavlink", {
            "lat": round(plat, 6), "lon": round(plon, 6),
            "alt": 120.0, "heading": round(phdg, 1), "velocity": 14.0,
            "callsign": "DEMO-Patrol-1",
        })

        # ── hostile threat: inbound toward center, then loop ──
        threat_t += 0.06
        if threat_t > 1.0:
            threat_t = 0.0
        tlat = CLAT + threat_start[0] * (1 - threat_t)
        tlon = CLON + threat_start[1] * (1 - threat_t)
        thdg = bearing(tlat, tlon, CLAT, CLON)
        pub("det.air.contact.DEMO-THREAT", {
            "id": "DEMO-THREAT",
            "entity_type": "aircraft_unknown",
            "affiliation": "hostile",
            "name": "DEMO-Inbound",
            "latitude": round(tlat, 6), "longitude": round(tlon, 6),
            "altitude": 300.0, "heading": round(thdg, 1), "velocity": 45.0,
            "confidence": 0.82, "source": "demo",
        })

        tick += 1
        if tick % 5 == 1:
            print(f"  t={tick:>4}  patrol@{plat:.4f},{plon:.4f}  "
                  f"threat@{tlat:.4f},{tlon:.4f} (inbound {threat_t*100:.0f}%)")
        time.sleep(TICK)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\ndemo-feed stopped.")
        sys.exit(0)
