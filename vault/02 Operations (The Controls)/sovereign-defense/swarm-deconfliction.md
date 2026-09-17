---
domain: 01
name: Sovereign Defense — Swarm Deconfliction
status: active
owner: zeus
tags: [sovereign-defense, swarm, deconfliction, collision-avoidance, multi-drone]
created: 2026-05-13
updated: 2026-05-13
---

# Swarm Deconfliction

**Preventing Zeus's drones from hitting each other.** With 3+ drones orbiting the same GPS point, friendly collisions are the most likely failure mode. This document covers the approaches, their tradeoffs, and what to implement at each phase.

---

## The Problem

When WARFIGHTER activates, 3 drones simultaneously:
- Launch from nearby positions
- Fly to the same GPS coordinate (Zeus)
- Orbit the same point at the same altitude
- Respond to the same voice commands simultaneously

Without deconfliction, they will collide. This is the hardest unsolved engineering problem in the system.

---

## Deconfliction Approaches — Overview

| Approach | Complexity | Reliability | Latency | Phase |
| --- | --- | --- | --- | --- |
| **Altitude separation** | Low | High | Zero | Phase 1 (now) |
| **Role assignment** | Low-Medium | High | Near-zero | Phase 1 (now) |
| **Geographic partitioning** | Medium | High | Near-zero | Phase 2 |
| **Token ring (take turns)** | Medium | Medium | Low | Phase 2 |
| **Leader-follower formation** | High | Medium | Low | Phase 3 |
| **Decentralized swarm AI** | Very High | Medium | Low | Phase 4+ |

---

## Approach 1: Altitude Separation (Implement Now)

**The simplest and most reliable method. Each drone is permanently assigned a different altitude band.**

Each drone is pre-configured with a fixed mission altitude. They fly at different heights — physically impossible to collide at their assigned altitudes.

| Drone | Callsign | Assigned Altitude | Role |
| --- | --- | --- | --- |
| **D1** | OVERWATCH | 50m | High overwatch, wide-angle camera, maximum coverage |
| **D2** | GUARDIAN | 25m | Mid-level orbit, close surveillance |
| **D3** | SCOUT | 12m | Low-level close protection, detail camera |
| **D4** (if present) | RELAY | 35m | Fills the gap between D1 and D2, relay |

**PX4 configuration:**

In each drone's QGC → Parameters:
- `MIS_TAKEOFF_ALT` for D1: 50
- `MIS_TAKEOFF_ALT` for D2: 25
- `MIS_TAKEOFF_ALT` for D3: 12
- Each M01 GUARDIAN_ORBIT waypoint altitude set to the drone's assigned altitude

**Separation buffer:** With 12m, 25m, 50m — minimum separation is 13m. At 5-10m/s horizontal speed, this is adequate for consumer drone maneuvers. Not adequate for FPV-speed drones (20+ m/s), but your drones won't be flying that fast in formation.

**Result:** Even if all three drones are orbiting Zeus at the same GPS lat/lon, they physically cannot collide because they're at different altitudes. Dead simple. No software needed. Zero latency. Implement this first.

---

## Approach 2: Role Assignment (Implement Now)

**Each drone has a fixed role and behavior. They don't all do the same thing at the same time.**

In the WARFIGHTER protocol, the 3 drones are not interchangeable — they each have a specific mission:

| Drone | Mission | Orbit Radius | Direction |
| --- | --- | --- | --- |
| **D1 (OVERWATCH)** | M05 Emergency Overwatch — max alt, wide camera, full scene | 50m | Clockwise |
| **D2 (GUARDIAN)** | M01 Guardian Orbit — mid-level protective circle | 30m | Counter-clockwise |
| **D3 (SCOUT)** | M07 Area Denial — flies toward detected threat bearing, siren | Threat-directed | N/A |

With different altitudes, radii, AND directions:
- D1 and D2 never cross the same point at the same altitude
- D3 flies away from Zeus toward the threat — physically separating from D1/D2

**Implementation:** Pre-program each drone's SD card with its specific mission set. D1's mission files have altitude 50m. D2's have altitude 25m. D3's have altitude 12m and directed toward threat. When WARFIGHTER activates, each drone executes its own role automatically.

---

## Approach 3: Geographic Partitioning (Phase 2)

**Divide the operational area into sectors. Each drone owns a sector — cannot enter another drone's sector.**

For property patrol and search operations with multiple drones:

```
PROPERTY AERIAL VIEW (top-down)

  ┌─────────────────────────────────────┐
  │          NORTH                       │
  │    D1 SECTOR (north half)            │
  │                                      │
  │              ─ ─ ─ ─ ─              │
  │                                      │
  │    D2 SECTOR (south half)            │
  │          SOUTH                       │
  └─────────────────────────────────────┘
```

Or for circular orbit around Zeus:

```
                  N
                  ↑
           D1 (0°-120°)
          /              \
   D3 (240°-360°)   D2 (120°-240°)
          \              /
                  S
```

Each drone's M01 orbit has a `YAW_TO_WAYPOINT` parameter that constrains it to its sector arc. Even at the same radius, they never occupy the same angular position simultaneously.

**PX4 implementation:** Set `NAV_LOITER_RAD` to same radius for all drones but offset their initial position around the orbit by 120° each. Their continuous orbit keeps them 120° apart indefinitely.

**Orbit synchronization via Meshtastic:** Each drone broadcasts its current angular position on mesh (0-360°) every 5 seconds. If two drones are within 30° of each other, the lower-priority drone slows (hover) until separation increases. Implemented in the companion Jetson.

---

## Approach 4: Token Ring (Phase 2)

**One token circulates among drones. Only the drone holding the token can execute a maneuver near other drones.**

Useful for coordinated actions — like all drones converging on Zeus's position after a WARFIGHTER trigger.

**Protocol:**

1. On WARFIGHTER activation, token holder = D1 (always, by rule)
2. D1 launches immediately, flies to Zeus GPS, begins orbit
3. D1 broadcasts "TOKEN: D2" on mesh when it's safely orbiting
4. D2 launches on receiving token, flies to Zeus GPS at its altitude, begins orbit
5. D2 broadcasts "TOKEN: D3" when safely orbiting
6. D3 launches on receiving token

**Staggered launch prevents simultaneous collision risk:**

| Time | Event |
| --- | --- |
| T+0s | D1 launches |
| T+3s | D1 reaches 20m altitude, broadcasts TOKEN:D2 |
| T+3s | D2 launches |
| T+6s | D2 reaches 20m altitude, broadcasts TOKEN:D3 |
| T+6s | D3 launches |
| T+15s | All 3 in formation, deconflicted by altitude + stagger |

**Fallback:** If token is not received within 8 seconds (drone that has token crashed or lost comms), the next drone launches anyway. Safety > perfect coordination.

**Implemented in Meshtastic mesh bridge on each Jetson.**

---

## Approach 5: Leader-Follower Formation (Phase 3)

**One drone is the leader (D1). D2 and D3 track the leader's position and maintain fixed offsets.**

More sophisticated — allows the formation to move together coherently when Zeus moves.

```
          D1 (LEADER)
         /            \
        /  offset       \
      D2 (-5m, +15m alt)  D3 (+5m, +5m alt)
```

D2 always stays 5m to the left and 15m above D1.  
D3 always stays 5m to the right and 5m above D1.

When D1 moves, D2 and D3 move to maintain their offsets — no collision possible because offsets guarantee physical separation.

**PX4 implementation:** [PX4 Formation Flight](https://docs.px4.io/main/en/advanced_features/formation_controller/) — experimental feature available in PX4 main branch. Requires OFFBOARD mode and companion computer sending formation positions.

**Companion Jetson code (simplified):**

```python
# follower.py — runs on D2 and D3 companion Jetsons
# D2 offset: x=-5m, y=0, z=+15m relative to leader

import meshtastic
from pymavlink import mavutil

OFFSET_X = -5.0  # meters (left of leader)
OFFSET_Y = 0.0   # meters (same forward position)
OFFSET_Z = 15.0  # meters (above leader)

def on_leader_position(packet):
    if packet.get('fromId') != 'D1': return
    pos = packet['decoded']['position']
    # Convert leader lat/lon to local NED frame
    # Apply offset in NED
    # Convert back to lat/lon
    # Send POSITION_TARGET_GLOBAL_INT to autopilot
    target_lat, target_lon, target_alt = apply_offset(
        pos['latitude'], pos['longitude'], pos['altitude'],
        OFFSET_X, OFFSET_Y, OFFSET_Z
    )
    send_position_target(target_lat, target_lon, target_alt)

mesh.subscribe(on_leader_position, "meshtastic.receive.position")
```

D1 broadcasts its GPS position on mesh every 500ms. D2 and D3 receive it, apply their offset, and command their autopilot to the offset position. The swarm moves as a rigid formation.

**Collision guarantee:** By design — if D1 is at position P, D2 is always at P + offset, D3 is always at P + different offset. They cannot converge to the same point unless the offset math breaks. Verify the offset math before deployment.

---

## Approach 6: Sense-and-Avoid (Phase 3-4)

**Each drone detects other drones using onboard sensors and avoids them autonomously.**

Each drone already has a depth camera (Intel RealSense D435) for obstacle avoidance. The same camera can detect other drones:

1. YOLO fine-tuned on drone images (your own drones specifically, known visual signatures)
2. When a drone is detected within 10m → veer away from it
3. Resume mission after clearing separation

**Limitations:**
- Works only within depth camera range (~3-5m effective for small fast objects)
- High latency if using vision-based detection (~100ms per frame)
- False positive risk if flying near trees, birds, other aircraft
- Not a primary deconfliction method — a last-resort failsafe

**PX4 integration:** Obstacle avoidance with `COM_OBS_AVOID: 1` in PX4 parameters. The companion computer runs obstacle avoidance locally and commands the autopilot to maneuver around detected obstacles via OFFBOARD mode.

---

## Recommended Implementation Stack

Build these in order — each phase adds a layer on top of the previous:

### Phase 1 (Now — Before First Multi-Drone Flight)

**Altitude Separation + Role Assignment — implement in mission planning**

- [ ] Assign permanent altitude to each drone: D1=50m, D2=25m, D3=12m
- [ ] Pre-program each drone's SD card with its role-specific missions
- [ ] Test: arm all 3 drones on ground simultaneously, verify mission parameters per drone
- [ ] Never allow two drones at the same altitude unless you've verified horizontal separation >50m

### Phase 2 (Before WARFIGHTER Drills)

**Token Ring + Geographic Partitioning**

- [ ] Implement Meshtastic token broadcast on each companion Jetson
- [ ] Test staggered launch: D1 launches → waits for D1 stable → D2 launches → waits → D3 launches
- [ ] Implement 120° angular partitioning for orbit separation
- [ ] Test: 3 drones orbiting same GPS point — verify they maintain 120° separation

### Phase 3 (Before Sustained Fleet Ops)

**Leader-Follower Formation**

- [ ] Enable PX4 OFFBOARD mode on D2 and D3
- [ ] Implement follower.py on D2 and D3 companion Jetsons
- [ ] Test: D1 moves, D2 and D3 follow in formation
- [ ] Verify offset math with GPS ground truth
- [ ] Enable depth camera obstacle avoidance as fallback layer

---

## Failure Mode Analysis

| Failure | Effect Without Deconfliction | Effect With Altitude Separation |
| --- | --- | --- |
| Two drones orbit same point | Collision risk at orbit crossings | No collision — different altitudes |
| Token ring message lost | Both D2 and D3 launch simultaneously | Collision risk during ascent — add timeout |
| D1 crashes during formation | D2/D3 may fly toward crash site | RTH triggers — they go home, not to crash site |
| GPS drift (both drones drift to same point) | Collision possible | Altitude separation still prevents collision |
| One drone runs out of battery, descends | Descending D1 may hit ascending D3 | Monitor: D1 descends from 50m, D3 is at 12m — safe |
| Companion Jetson freezes | Drone flies last commanded path | Altitude separation still protects |

**The one failure altitude separation doesn't protect against:** A drone at 50m descending through 25m and 12m to land, while D2 and D3 are actively orbiting at those altitudes. 

**Solution:** Landing corridor protocol. Before any drone descends to land:
1. Drone broadcasts "LANDING: D1" on mesh with landing waypoint GPS
2. Other drones within 30m horizontal radius move 40m away horizontally
3. D1 descends through clear corridor
4. D1 broadcasts "LANDED: D1" — other drones return to orbit

---

## Quick Reference — Deconfliction Rules

**Non-negotiable minimums:**

1. No two drones at the same altitude unless >50m horizontal separation
2. Staggered launch: minimum 3-second gap between drone launches
3. D1 always launches first, orbits at 50m
4. D2 always launches second, orbits at 25m
5. D3 always launches third, orbits at 12m or flies toward threat
6. Landing corridor broadcast before any descent
7. If comms lost between drones: all drones RTH (not hover)

*Full suite: [[README]] | Operating procedures: [[wartime-weapon-protocol]] | Hardware: [[hardware-bom-acquisition]]*
