---
domain: 01
name: Sovereign Defense — PX4 Mission Planning Guide
status: active
owner: zeus
tags: [sovereign-defense, px4, ardupilot, missions, qgroundcontrol, autonomous]
created: 2026-05-13
updated: 2026-05-13
---

# PX4 Mission Planning Guide

**How to build, upload, and verify the 8 pre-loaded autonomous missions (M01-M08) in QGroundControl. Step-by-step from zero to a flight controller that protects Zeus without any connectivity.**

---

## Prerequisites

Before this guide: drone hardware assembled, PX4 firmware flashed to Pixhawk 6C, basic stabilization tested. If you haven't done that, follow the [Holybro Pixhawk 6C quickstart](https://docs.holybro.com/autopilot/pixhawk-6c/quickstart).

**Software needed:**

| Tool | Where | Cost |
| --- | --- | --- |
| QGroundControl (QGC) | qgroundcontrol.com | Free |
| PX4 firmware | firmware.ardupilot.org or via QGC | Free |
| Meshtastic app | meshtastic.org | Free |
| Mission Planner (optional, ArduPilot alt) | ardupilot.org/planner | Free |

---

## One-Time Setup

### 1. Flash PX4 Firmware

Connect Pixhawk 6C via USB to Mac. Open QGC → Vehicle Setup → Firmware → PX4 Pro Stable.

Let QGC flash and reboot. Verify: QGC shows vehicle connected, telemetry green.

### 2. Configure Failsafes (Critical — Do This First)

In QGC → Vehicle Setup → Safety:

| Failsafe | Setting | Why |
| --- | --- | --- |
| **RC Loss Action** | Return to Launch | If RC signal dies → drone comes home |
| **RC Loss Timeout** | 5 seconds | Wait 5s before RTL (prevents false triggers) |
| **Data Link Loss** | Return to Launch | If telemetry dies → drone comes home |
| **Data Link Timeout** | 60 seconds | Wait 60s (mesh may have brief gaps) |
| **Battery Failsafe — Warning** | 20% | Alert Zeus at 20% |
| **Battery Failsafe — Critical** | 15% → RTL | Auto-RTH at 15%, no exception |
| **Battery Failsafe — Emergency** | 10% → Land** | Land immediately if 10% (couldn't RTH) |
| **Geofence** | Return on breach | Define fence = property boundary + 100m buffer |
| **Geofence Action** | Return to Launch | Breach = go home, not land in-place |
| **Max Altitude** | 120m (400ft) | Hard ceiling, FAA compliance |

**Geofence setup:** QGC → Plan → Geofence → Add Inclusion Zone → draw polygon around property + buffer → Upload.

### 3. Set Home Position

Home position = where the drone will RTH. Set this to your charging/storage location on property. In QGC: right-click map → Set Home Position here. This gets saved to the flight controller.

### 4. Configure RC Channels for Mission Triggers

In QGC → Vehicle Setup → Parameters → RC_MAP:

| RC Channel | Assign To | Function |
| --- | --- | --- |
| CH5 | Flight mode switch | STABILIZE / POSITION / MISSION |
| CH6 | Aux 1 | Kill switch (disarm motors emergency) |
| CH7 | Aux 2 | Mission select (toggle between loaded missions) |
| CH8 | Aux 3 | RTH override (forces RTH regardless of mode) |

On Herelink: Map the 4 programmable buttons to CH6-CH9 via QGC.

| Herelink Button | Channel | Function |
| --- | --- | --- |
| **A** | CH7 position 1 | Start M01 (Guardian Orbit) |
| **B** | CH7 position 2 | Start M02 (Property Patrol) |
| **C** | CH8 | RTH ALL |
| **D** | CH6 | Emergency land (cut motors on ground) |

---

## Building Each Mission

### QGC Plan View Basics

QGC → Plan (top menu bar) → map loads, centered on your GPS location.

- **Waypoint:** Click map → green dot appears → drag to position → set altitude in right panel
- **Survey:** Plan → Survey → draw polygon → QGC auto-generates grid pattern
- **Upload:** Plan → Upload → saves to flight controller SD card
- **Download:** Plan → Download → saves to your Mac as .plan file

Always save `.plan` files to a folder on your Mac as backup. Losing the SD card means losing your missions.

---

### M01 — GUARDIAN ORBIT

**Trigger:** Panic button / dead-man switch / WARFIGHTER activation  
**Behavior:** Launches, flies to Zeus's GPS, orbits indefinitely at 30m altitude, 30m radius

**Build steps:**

1. QGC → Plan → Mission
2. Add Takeoff waypoint at home location, altitude 10m
3. Add LOITER_UNLIMITED waypoint at Zeus's typical position (adjust on the fly via parameter)
   - Altitude: 30m AGL
   - Radius: 30m
   - Direction: CCW (clockwise keeps camera pointing inward)
4. No RTH waypoint — M01 loops until battery failsafe triggers
5. Upload → rename `M01_GUARDIAN_ORBIT.plan`

**Key parameters to set:**
- `NAV_LOITER_RAD`: 30 (meters)
- `MIS_TAKEOFF_ALT`: 10 (meters)

**Dynamic position update:** The Meshtastic → MAVLink bridge updates the LOITER waypoint with Zeus's real-time GPS. Set `MPC_XY_TRAJ_P` higher for faster position tracking when Zeus moves.

---

### M02 — PROPERTY PATROL

**Trigger:** Scheduled cron / button press  
**Behavior:** Flies pre-loaded waypoint fence around property perimeter, repeats until battery or recall

**Build steps:**

1. QGC → Plan → Survey, OR manually place waypoints
2. Manual waypoints (recommended for property patrol):
   - Walk your perimeter with phone GPS on
   - Note lat/lon at each corner and any notable points (gate, barn, access road)
   - Enter each as a waypoint in QGC at 50m altitude
3. Set speed: 5-8 m/s (camera can capture stable imagery at this speed)
4. Add LOITER for 10s at each waypoint if you want it to pause and look
5. Add RTH at end OR set to LOITER_UNLIMITED at last waypoint to hold
6. Upload → rename `M02_PROPERTY_PATROL.plan`

---

### M03 — PERIMETER SWEEP

**Trigger:** Button press / mesh command  
**Behavior:** Visits each perimeter sensor location, hovers 10s, takes photo, moves to next

**Build steps:**

1. Get GPS coordinates of each Home Assistant sensor on perimeter
2. QGC → Plan → Mission
3. Add DO_DIGICAM_CONTROL command at each sensor waypoint (triggers camera shutter)
4. Add DELAY command (10 seconds) at each sensor point for hover
5. Connect all sensor waypoints in perimeter order
6. RTH at end
7. Upload → rename `M03_PERIMETER_SWEEP.plan`

---

### M04 — FOLLOW BEACON

**Trigger:** Button press / mesh command `D1:FOLLOW:ZEUS`  
**Behavior:** Follows Zeus's Meshtastic GPS beacon at 30m altitude, 15m behind

**This is not a static QGC mission — it's a dynamic mode.**

In PX4: enable `FOLLOW_TARGET` mode (requires telemetry link or MAVLink injection):

1. QGC → Parameters → `FLW_TGT_ALT`: 30 (meters)
2. `FLW_TGT_DST`: 15 (meters behind target)
3. `FLW_TGT_MAX_VEL`: 5 (m/s max follow speed)

The Meshtastic-to-MAVLink bridge (running on drone companion Jetson) reads Zeus's GPS from mesh and injects it as a FOLLOW_TARGET MAVLink message every 5 seconds.

**Bridge code (runs on Jetson):**

```python
# meshtastic_to_follow_target.py
import meshtastic
import meshtastic.serial_interface
from pymavlink import mavutil
import time

mesh = meshtastic.serial_interface.SerialInterface()
mav = mavutil.mavlink_connection('/dev/ttyACM0', baud=57600)

def on_receive(packet, interface):
    decoded = packet.get('decoded', {})
    pos = decoded.get('position', {})
    if pos and packet.get('fromId') == 'ZEUS':
        lat = pos.get('latitude', 0)
        lon = pos.get('longitude', 0)
        alt = pos.get('altitude', 0)
        # Inject as FOLLOW_TARGET MAVLink message
        mav.mav.follow_target_send(
            int(time.time() * 1000),  # timestamp ms
            1,  # est_capabilities
            int(lat * 1e7),  # lat (1e7 degrees)
            int(lon * 1e7),  # lon
            float(alt),      # alt MSL
            [0,0,0],         # vel (unknown)
            [0,0,0],         # acc (unknown)
            [1,0,0,0],       # attitude quaternion
            0, 0             # rates
        )

mesh.localNode.setOwner('DRONE1')
pub.subscribe(on_receive, "meshtastic.receive.position")
while True:
    time.sleep(0.1)
```

---

### M05 — EMERGENCY OVERWATCH

**Trigger:** WARFIGHTER auto  
**Behavior:** Climbs to maximum safe altitude (120m), wide-angle camera, slow 360° orbit, records

**Build steps:**

1. QGC → Plan → Mission
2. Takeoff waypoint, altitude 120m
3. LOITER_UNLIMITED waypoint at home position, altitude 120m, radius 50m
4. Camera set to widest FOV, video record mode
5. Upload → rename `M05_EMERGENCY_OVERWATCH.plan`

Companion Jetson starts YOLO detection and face recognition when M05 activates. All tagged contacts logged with timestamp and GPS.

---

### M06 — RTH SAFE

**Trigger:** Any RTH command (voice, button, failsafe, mesh)  
**Behavior:** Navigate to home coordinates, descend to 5m, land, disarm

This is PX4's built-in Return mode — no custom mission needed.

In QGC → Parameters:
- `RTL_RETURN_ALT`: 50m (climb to this altitude before RTH if currently lower)
- `RTL_DESCEND_ALT`: 10m (hover at this altitude above home before landing)
- `RTL_LAND_DELAY`: 0.5s (brief pause before landing)

Trigger: RC → flip to RTL mode, or Herelink Button C, or voice "RTH", or mesh `D1:RTH`.

---

### M07 — AREA DENIAL SIREN

**Trigger:** WARFIGHTER auto  
**Behavior:** Flies to threat bearing at 20m altitude, activates siren + strobe, orbits threat point aggressively

**Build steps:**

1. This mission requires the threat bearing calculated by YOLO/acoustic triangulation
2. The Jetson companion dynamically generates this mission at runtime:

```python
# generate_m07.py — called when WARFIGHTER activates
def generate_area_denial_mission(threat_lat, threat_lon, home_lat, home_lon):
    mission = [
        # Takeoff
        {'seq': 0, 'frame': 3, 'command': 22, 'param7': 20},  # TAKEOFF 20m
        # Fly to threat orbit point
        {'seq': 1, 'frame': 0, 'command': 17,  # LOITER_UNLIM
          'param1': 0, 'param2': 0, 'param3': 15,  # radius 15m
          'param4': 0,
          'x': int(threat_lat * 1e7),
          'y': int(threat_lon * 1e7),
          'z': 20},  # 20m altitude
    ]
    # Upload via MAVLink
    upload_mission(mission)
    # Activate siren + strobe GPIO pins
    activate_deterrence()
```

---

### M08 — SEARCH GRID

**Trigger:** Mesh command `D1:SEARCH:lat,lon,radius`  
**Behavior:** Expanding square search pattern from given GPS center point

**Build steps:**

1. Jetson companion generates expanding square dynamically from given center + radius
2. QGC Survey pattern at 30m altitude, 70% overlap for thorough coverage
3. Camera continuously recording, YOLO running on live feed
4. Any detection → Jetson sends alert on mesh and to OpenClaw (if online)
5. RTH when target found or search area exhausted

---

## Testing Protocol (Before Any Live Use)

### Ground Test (No Flight)

1. Power on drone, connect QGC via telemetry
2. Upload each mission (M01-M08)
3. Verify in QGC: all waypoints appear on map in correct positions
4. Verify failsafe parameters loaded: Vehicle Setup → Safety → check all values
5. Arm drone in QGC (not motors), verify no pre-arm warnings
6. Simulate RC loss: unplug RC receiver → verify QGC shows "RC Loss" alert → RTH should appear in QGC log

### Low-Altitude Test Flight (CRITICAL — do this before any autonomous use)

1. Open field, no obstructions, no people nearby
2. Test M02 (Property Patrol) at 20m altitude first — lowest-risk mission
3. Watch drone execute waypoints. Verify it follows plan exactly.
4. Test failsafe: kill RC during M02 → verify drone stops mission and RTH
5. Test battery: let it fly until battery alert → verify RTH triggers at configured %
6. Test geofence: manually fly to geofence boundary → verify it stops and turns back

### Panic Button Test

1. Drone on ground, armed and ready
2. Press panic button → verify Jetson receives trigger → verify M01 uploads and starts
3. Verify drone launches and orbits Zeus's GPS
4. Verify siren and strobe activate
5. Say "STAND DOWN WARFIGHTER" → verify drone transitions to LOITER → then RTH

---

## Parameter Cheat Sheet

Key PX4 parameters to set for this system. In QGC → Parameters → search by name.

| Parameter | Value | Purpose |
| --- | --- | --- |
| `RTL_RETURN_ALT` | 50 | RTH climb altitude (meters) |
| `RTL_DESCEND_ALT` | 10 | RTH hover altitude before landing |
| `NAV_LOITER_RAD` | 30 | Default orbit radius (meters) |
| `MIS_TAKEOFF_ALT` | 10 | Minimum takeoff altitude |
| `FLW_TGT_ALT` | 30 | Follow target altitude |
| `FLW_TGT_DST` | 15 | Follow target distance behind |
| `GF_ACTION` | 2 | Geofence breach action (2 = RTL) |
| `GF_MAX_HOR_DIST` | 500 | Max horizontal geofence distance |
| `CBRK_IO_SAFETY` | 22027 | Disable safety button (if no physical button) |
| `BAT1_CAPACITY` | Set to battery mAh | Accurate battery % |
| `BAT_CRIT_THR` | 0.15 | Critical battery = 15% |
| `BAT_EMERGEN_THR` | 0.10 | Emergency battery = 10% |
| `COM_RC_LOSS_T` | 5.0 | Seconds before RC loss failsafe |
| `COM_DL_LOSS_T` | 60.0 | Seconds before data link loss failsafe |

---

## Mission Files Backup Location

All `.plan` files saved to: `_system/drone-missions/`

| File | Mission |
| --- | --- |
| `M01_GUARDIAN_ORBIT.plan` | WARFIGHTER guardian orbit |
| `M02_PROPERTY_PATROL.plan` | Standard patrol route |
| `M03_PERIMETER_SWEEP.plan` | Sensor-by-sensor sweep |
| `M06_RTH_SAFE.plan` | RTH (PX4 built-in, no file needed) |
| `M05_EMERGENCY_OVERWATCH.plan` | Max altitude overwatch |
| `M04_FOLLOW_BEACON.py` | Dynamic follow mode (Jetson script) |
| `M07_AREA_DENIAL.py` | Dynamic area denial (Jetson script) |
| `M08_SEARCH_GRID.py` | Dynamic search grid (Jetson script) |

*Full suite: [[README]] | Operating procedures: [[wartime-weapon-protocol]]*
