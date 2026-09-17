# Full FABRIC Fleet: Ground Robots + Drones

> **Created:** 2026-01-22
> **Vision:** Complete autonomous fleet with OpenMind OM1 structure
> **Integration:** Ground robots + aerial drones + specialty robots

---

## Overview

Combine **ground robots** (Unitree, Boston Dynamics) with **aerial drones** (Zeus, DJI, Skydio) into a **unified FABRIC fleet** using OpenMind's architecture.

**Key Concept:** ONE edge device controls BOTH robots AND drones at each location.

```
┌───────────────────────────────────────────────────────────────┐
│         SINGLE EDGE DEVICE ($997) Per Location                │
│                                                                │
│  ┌──────────────────────────────────────────────────────┐    │
│  │    OneMind Unified Fleet Manager                     │    │
│  │    (Combines Robot + Drone managers)                 │    │
│  └────────────────┬─────────────────────────────────────┘    │
│                   │                                            │
│         ┌─────────┴─────────┐                                 │
│         │                   │                                 │
│         ▼                   ▼                                 │
│  ┌─────────────┐     ┌─────────────┐                         │
│  │   Robot     │     │   Drone     │                         │
│  │  Manager    │     │  Manager    │                         │
│  │  (OM1)      │     │  (OneMind)  │                         │
│  └──────┬──────┘     └──────┬──────┘                         │
│         │                   │                                 │
│    ┌────┴────┐         ┌────┴────┬────┐                      │
│    ▼         ▼         ▼         ▼    ▼                      │
│ ┌──────┐ ┌──────┐  ┌─────┐  ┌─────┐ ┌─────┐                 │
│ │Unit- │ │Unit- │  │Zeus │  │Farm │ │Sec- │                 │
│ │ree 1 │ │ree 2 │  │Drone│  │Drone│ │urity│                 │
│ └──────┘ └──────┘  └─────┘  └─────┘ └─────┘                 │
│   Ground   Ground     Air       Air     Air                   │
│                                                                │
│  All controlled via NATS + Legacy voice commands              │
└────────────────────────────────────────────────────────────────┘
```

---

## Architecture: Unified Fleet Manager

### Core Concept

Instead of separate managers, create **Unified Fleet Manager** that coordinates BOTH robots and drones:

```python
# agents/edge/unified_fleet/manager.py

from drone_manager.core import DroneManager
from openmind_om1 import RobotManager  # OpenMind OM1
import asyncio

class UnifiedFleetManager:
    """
    Unified management for ground robots + aerial drones

    Coordinates multi-modal operations:
    - Ground robots for detailed inspection
    - Drones for aerial overview
    - Specialty robots for specific tasks
    """

    def __init__(self, nats_url: str):
        self.nats_url = nats_url
        self.nc = None

        # Initialize sub-managers
        self.robot_manager = RobotManager(nats_url)
        self.drone_manager = DroneManager(nats_url)

        # Fleet state
        self.robots = {}  # {robot_id: robot_state}
        self.drones = {}  # {drone_id: drone_state}

    async def connect(self):
        """Connect to NATS and initialize managers"""
        self.nc = await nats.connect(self.nats_url)

        # Connect sub-managers
        await self.robot_manager.connect()
        await self.drone_manager.connect()

        # Subscribe to unified fleet commands
        await self.nc.subscribe(
            "commands.farm.fleet.>",
            cb=self.handle_fleet_command
        )

    async def handle_fleet_command(self, msg):
        """
        Handle unified fleet commands

        Examples:
          commands.farm.fleet.patrol             # Coordinate robots + drones
          commands.farm.fleet.search.{target}    # Multi-modal search
          commands.farm.fleet.inspect.{location} # Ground + air inspection
        """
        subject_parts = msg.subject.split(".")
        # commands.farm.fleet.{command}.{params}
        command = subject_parts[3]

        data = json.loads(msg.data.decode()) if msg.data else {}

        if command == "patrol":
            await self.coordinated_patrol(data)

        elif command == "search":
            target = subject_parts[4] if len(subject_parts) > 4 else "unknown"
            await self.coordinated_search(target, data)

        elif command == "inspect":
            location = subject_parts[4] if len(subject_parts) > 4 else None
            await self.coordinated_inspection(location, data)

        elif command == "emergency":
            emergency_type = subject_parts[4]
            await self.emergency_response(emergency_type, data)

    # ============================================
    # Coordinated Operations
    # ============================================

    async def coordinated_patrol(self, params: dict):
        """
        Coordinate ground robots and drones for comprehensive patrol

        Strategy:
        - Drones provide aerial overview (fast, wide coverage)
        - Ground robots investigate points of interest (detailed)
        """

        # Start aerial patrol (fast sweep)
        drone_patrol_task = asyncio.create_task(
            self.drone_manager.command_drone(
                "security",
                "patrol",
                {
                    "pattern": "perimeter",
                    "altitude": 20,
                    "speed": 7,  # m/s
                }
            )
        )

        # Start ground patrol (detailed)
        robot_patrol_task = asyncio.create_task(
            self.robot_manager.command_robot(
                "unitree1",
                "patrol",
                {
                    "waypoints": params.get("ground_waypoints", []),
                    "speed": 1.5,  # m/s
                }
            )
        )

        # Monitor both patrols
        while True:
            # Check drone detections
            drone_status = await self.drone_manager.get_drone_status("security")
            if drone_status.threat_detected:
                # Dispatch ground robot to investigate
                threat_location = drone_status.threat_location

                await self.robot_manager.command_robot(
                    "unitree1",
                    "goto",
                    {
                        "latitude": threat_location.lat,
                        "longitude": threat_location.lon,
                    }
                )

                # Drone hovers and monitors
                await self.drone_manager.command_drone(
                    "security",
                    "orbit",
                    {
                        "latitude": threat_location.lat,
                        "longitude": threat_location.lon,
                        "radius": 15,
                        "altitude": 20,
                    }
                )

            # Check if patrols complete
            if drone_patrol_task.done() and robot_patrol_task.done():
                break

            await asyncio.sleep(5)

    async def coordinated_search(self, target: str, params: dict):
        """
        Multi-modal search (air + ground)

        Use Case: Lost item, missing animal, person search
        """

        search_area = params.get("area")

        # Phase 1: Aerial sweep (fast coverage)
        drone_results = await self.drone_manager.command_drone(
            "zeus",
            "grid_search",
            {
                "area": search_area,
                "altitude": 30,
                "overlap": 70,
            }
        )

        # Analyze aerial photos
        potential_locations = await analyze_aerial_search(
            drone_results.images,
            target=target
        )

        if not potential_locations:
            return {"status": "not_found", "coverage": "complete"}

        # Phase 2: Ground investigation (detailed)
        for location in potential_locations:
            await self.robot_manager.command_robot(
                "unitree1",
                "goto",
                location
            )

            # Close-up inspection
            video = await self.robot_manager.get_robot_video("unitree1")
            analysis = await legacy.analyze_video(
                video,
                objectives=[f"confirm_{target}_present"]
            )

            if analysis.confirmed:
                return {
                    "status": "found",
                    "location": location,
                    "confidence": analysis.confidence
                }

        return {"status": "not_found", "locations_checked": len(potential_locations)}

    async def coordinated_inspection(self, location_name: str, params: dict):
        """
        Multi-angle inspection (air + ground)

        Use Case: Inspect building, equipment, fence line
        """

        location = await db.get_location(location_name)

        # Aerial overview (top + sides)
        aerial_task = asyncio.create_task(
            self.drone_manager.command_drone(
                "zeus",
                "orbit",
                {
                    "latitude": location.lat,
                    "longitude": location.lon,
                    "radius": 10,
                    "altitude": 15,
                    "orbits": 2,  # Full 720° coverage
                }
            )
        )

        # Ground-level inspection (front + details)
        ground_task = asyncio.create_task(
            self.robot_manager.command_robot(
                "unitree1",
                "inspect",
                {
                    "target": location_name,
                    "angles": ["front", "sides", "back"],
                    "close_ups": True,
                }
            )
        )

        # Wait for both to complete
        await asyncio.gather(aerial_task, ground_task)

        # Combine results
        aerial_video = await self.drone_manager.get_latest_video("zeus")
        ground_video = await self.robot_manager.get_latest_video("unitree1")

        analysis = await legacy.analyze_videos(
            [aerial_video, ground_video],
            objectives=[
                "structural_damage",
                "equipment_condition",
                "maintenance_needs",
            ]
        )

        return {
            "location": location_name,
            "inspection_time": datetime.now(),
            "aerial_coverage": "complete",
            "ground_coverage": "complete",
            "issues": analysis.issues,
            "recommendations": analysis.recommendations,
        }

    async def emergency_response(self, emergency_type: str, params: dict):
        """
        Coordinate emergency response with all available units

        Examples:
        - Fire: Aerial assessment + ground evacuation
        - Intruder: Aerial tracking + ground interception
        - Animal attack: Both units converge
        """

        if emergency_type == "intruder":
            intruder_location = params["location"]

            # Drone: Track from above
            await self.drone_manager.command_drone(
                "security",
                "track",
                {
                    "subject": "person",
                    "location": intruder_location,
                }
            )

            # Robot: Approach on ground (loudspeaker warning)
            await self.robot_manager.command_robot(
                "unitree1",
                "goto",
                intruder_location
            )

            await self.robot_manager.command_robot(
                "unitree1",
                "speak",
                {
                    "message": "You are trespassing. Leave immediately. "
                              "You are being recorded."
                }
            )

        elif emergency_type == "fire":
            fire_location = params["location"]

            # Drone: Assess from safe distance
            await self.drone_manager.command_drone(
                "security",
                "orbit",
                {
                    "latitude": fire_location.lat,
                    "longitude": fire_location.lon,
                    "radius": 50,  # Safe distance
                    "altitude": 30,
                }
            )

            # Robot: Do NOT approach (heat damage)
            # Instead, evacuate area / check for people
            await self.robot_manager.command_robot(
                "unitree1",
                "patrol",
                {
                    "area": "safe_zone",
                    "objective": "evacuate_check"
                }
            )
```

---

## Use Case 1: Morning Property Check

**Scenario:** Comprehensive morning inspection using ALL units

```python
@workflow
async def morning_property_check():
    """
    Full property inspection with entire fleet

    Robots: Check buildings, equipment, feed/water
    Drones: Aerial overview, crop health, perimeter
    """

    legacy.speak("Starting morning property check with full fleet.")

    # ============================================
    # Phase 1: Aerial Survey (Fast Overview)
    # ============================================

    # Zeus: General property scan
    zeus_task = asyncio.create_task(
        drone_tool.command_drone(
            "zeus",
            "mission",
            generate_property_scan_pattern()
        )
    )

    # Farm Drone: Crop health check
    farm_task = asyncio.create_task(
        weekly_crop_health_scan(["east_field", "west_field"])
    )

    # Security Drone: Perimeter patrol
    security_task = asyncio.create_task(
        perimeter_patrol()
    )

    # ============================================
    # Phase 2: Ground Inspection (Detailed)
    # ============================================

    # Unitree 1: Barn and livestock check
    unitree1_task = asyncio.create_task(
        robot_tool.command_robot(
            "unitree1",
            "mission",
            {
                "waypoints": [
                    {"location": "barn_entrance", "action": "inspect"},
                    {"location": "feed_storage", "action": "check_levels"},
                    {"location": "water_troughs", "action": "check_water"},
                    {"location": "chicken_coop", "action": "count_chickens"},
                ]
            }
        )
    )

    # Unitree 2: Equipment and greenhouse check
    unitree2_task = asyncio.create_task(
        robot_tool.command_robot(
            "unitree2",
            "mission",
            {
                "waypoints": [
                    {"location": "equipment_yard", "action": "inspect_equipment"},
                    {"location": "greenhouse", "action": "check_plants"},
                    {"location": "irrigation_pump", "action": "check_status"},
                ]
            }
        )
    )

    # ============================================
    # Phase 3: Wait for All Units
    # ============================================

    results = await asyncio.gather(
        zeus_task,
        farm_task,
        security_task,
        unitree1_task,
        unitree2_task,
        return_exceptions=True
    )

    # ============================================
    # Phase 4: Compile Report
    # ============================================

    report = """
    Morning Property Check Complete

    Aerial Survey:
    - Zeus: {zeus_status}
    - Farm Drone: {farm_status}
    - Security: {security_status}

    Ground Inspection:
    - Unitree 1 (Barn): {unitree1_status}
    - Unitree 2 (Equipment): {unitree2_status}

    Issues Detected: {total_issues}
    Immediate Attention: {critical_issues}

    All systems operational. Property secure.
    """

    legacy.speak(
        "Morning property check complete. "
        f"{total_issues} issues detected. "
        "Detailed report sent to your email."
    )

    return report
```

---

## Use Case 2: Livestock Monitoring

**Scenario:** Ground robot counts animals, drone provides aerial view

```python
@workflow
async def coordinated_livestock_check(pasture_name: str):
    """
    Multi-angle livestock monitoring
    """

    pasture = await db.get_pasture(pasture_name)

    legacy.speak(f"Checking livestock in {pasture_name}")

    # Drone: Aerial count
    aerial_count_task = asyncio.create_task(
        drone_tool.command_drone(
            "security",
            "grid_search",
            {
                "area": pasture.boundary,
                "altitude": 40,
                "objective": "count_animals"
            }
        )
    )

    # Robot: Ground inspection
    ground_check_task = asyncio.create_task(
        robot_tool.command_robot(
            "unitree1",
            "patrol",
            {
                "area": pasture.name,
                "objective": "check_animal_health"
            }
        )
    )

    # Wait for both
    aerial_result, ground_result = await asyncio.gather(
        aerial_count_task,
        ground_check_task
    )

    # Compare counts
    aerial_count = aerial_result.animal_count
    expected_count = pasture.expected_count

    if aerial_count < expected_count:
        legacy.speak(
            f"WARNING: Only {aerial_count} animals counted. "
            f"Expected {expected_count}. Investigating."
        )

        # Search for missing animals
        await coordinated_search(
            target="missing_animals",
            area=pasture.extended_boundary
        )

    # Check health issues from ground robot
    if ground_result.health_issues:
        legacy.speak(
            f"Health issues detected: {len(ground_result.health_issues)} animals "
            "need attention."
        )

    legacy.speak(f"Livestock check complete. {aerial_count} animals, all accounted for.")
```

---

## Use Case 3: Security Response

**Scenario:** Intruder detected, coordinate response

```python
@workflow
async def coordinated_security_response(threat_location: dict):
    """
    Full fleet security response

    - Drone: Track from above, provide overview
    - Robot: Approach on ground, issue warning, record evidence
    """

    legacy.speak("SECURITY THREAT. Deploying full response.")

    # ============================================
    # Phase 1: Aerial Tracking
    # ============================================

    # Skydio: Track intruder
    await drone_tool.command_drone(
        "security",
        "track",
        {
            "subject": "person",
            "location": threat_location,
            "mode": "spotlight",
            "altitude_offset": 15,
        }
    )

    # Turn on strobe lights (deterrent)
    await drone_tool.command_drone("security", "strobe_lights", {"enabled": True})

    # ============================================
    # Phase 2: Ground Interception
    # ============================================

    # Unitree robot: Approach and warn
    await robot_tool.command_robot(
        "unitree1",
        "goto",
        threat_location
    )

    # Play warning message (loudspeaker)
    await robot_tool.command_robot(
        "unitree1",
        "speak",
        {
            "message": "STOP. You are trespassing on private property. "
                      "You are being recorded from the air and ground. "
                      "Leave immediately or police will be called.",
            "volume": "max",
            "repeat": 3,
        }
    )

    # ============================================
    # Phase 3: Monitor and Record
    # ============================================

    # Both units record
    drone_video_task = asyncio.create_task(
        drone_tool.record_video("security", duration=300)  # 5 min
    )

    robot_video_task = asyncio.create_task(
        robot_tool.record_video("unitree1", duration=300)
    )

    # Alert owner
    await send_critical_alert({
        "type": "intruder",
        "severity": "critical",
        "description": "Intruder on property. Full fleet deployed.",
        "drone_feed": "livekit://drone-security-feed",
        "robot_feed": "livekit://robot-unitree1-feed",
    })

    # Wait for user command or intruder leaves
    legacy.speak("Monitoring intruder. Say 'call 911' to alert authorities.")

    response = await wait_for_voice_command(timeout=180)

    if "911" in response:
        await call_authorities({
            "type": "active_intrusion",
            "units_deployed": ["security_drone", "unitree1_robot"],
            "evidence": [drone_video_task, robot_video_task],
        })

    legacy.speak("Security response complete. Evidence recorded.")
```

---

## Use Case 4: Emergency Evacuation

**Scenario:** Fire detected, evacuate livestock with robots, assess with drones

```python
@workflow
async def emergency_evacuation(emergency_type: str, location: dict):
    """
    Coordinate evacuation with full fleet
    """

    if emergency_type == "fire":
        legacy.speak("FIRE DETECTED. Initiating emergency evacuation.")

        # Call 911 immediately
        await call_authorities({
            "type": "fire",
            "location": location,
            "severity": "critical"
        })

        # Drones: Assess fire spread (safe distance)
        zeus_task = asyncio.create_task(
            drone_tool.command_drone(
                "zeus",
                "orbit",
                {
                    "location": location,
                    "radius": 50,  # Safe distance
                    "altitude": 40,
                    "duration": 300,  # Until told to stop
                }
            )
        )

        # Provide real-time fire spread info to firefighters
        await share_video_with_emergency_services("zeus")

        # Robots: Evacuate livestock (if safe)
        if is_safe_for_ground_approach(location):
            # Open gates, herd animals away
            await robot_tool.command_robot(
                "unitree1",
                "herd_animals",
                {
                    "from_area": "near_fire",
                    "to_area": "safe_zone",
                    "urgency": "high"
                }
            )

        legacy.speak("Emergency services notified. Live video feed active.")
```

---

## Fleet Composition Examples

### Starter Fleet ($6K)

```
1x OpenMind Edge ($997)
1x Unitree Go1 (owned) ($0)
1x Zeus DIY Drone ($765)
1x DJI Mini 4 Pro ($759)
1x Skydio 2+ ($1,099)
──────────────────────────
Total: $3,620

Capabilities:
- 1 ground robot
- 3 drones (personal, farm, security)
- Coordinated operations
- Voice control
```

### Farm Fleet ($25K)

```
2x OpenMind Edge ($1,994)  - Barn + greenhouse
2x Unitree Go1 ($7,000)
1x Zeus DIY ($765)
1x DJI Agras T40 ($16,000)
1x DJI Mavic 3 Multispectral ($5,500)
1x Skydio 2+ ($1,099)
──────────────────────────
Total: $32,358

Capabilities:
- 2 locations
- 2 ground robots
- 4 drones (personal, spraying, monitoring, security)
- Full agricultural automation
- Multi-location coordination
```

### Enterprise Fleet ($100K+)

```
3x OpenMind Edge ($2,991)
2x Boston Dynamics Spot ($150,000)
2x Unitree Go1 ($7,000)
1x Unitree H1 Humanoid ($90,000)
2x Zeus DIY ($1,530)
2x DJI Agras T40 ($32,000)
2x DJI Mavic 3 Multispectral ($11,000)
2x Skydio 2+ ($2,198)
──────────────────────────────
Total: $296,719

Capabilities:
- 3 locations
- 5 ground robots (including humanoid)
- 8 drones
- Complete automation
- Research platform
```

---

## Voice Command Examples

### Fleet Coordination

```
"Legacy, do a full property check"
→ Deploys all units (robots + drones) for comprehensive inspection

"Search for the missing calf"
→ Aerial sweep + ground search

"Inspect the barn from all angles"
→ Drone orbits, robot walks around

"Emergency at the greenhouse"
→ All available units converge

"What's the fleet status?"
→ "2 robots patrolling, 1 drone charging, 2 drones standby"
```

### Individual Units

```
"Unitree 1, check the chicken coop"
"Zeus, give me an aerial view of the property"
"Farm drone, start spraying"
"Security drone, patrol the perimeter"
"Have both robots meet at the barn"
```

---

## Integration with Existing Robot Documentation

This fleet integration builds on:
1. **Robot Integration** docs (Unitree + OpenMind OM1)
2. **Multi-Drone Fleet** docs (Zeus + DJI + Skydio)

**Now combined into unified FABRIC fleet.**

---

## ROI Analysis

### Cost Per Capability

**Aerial Coverage:**
- Drones: $765-$16K one-time
- Human (ATV): $15K/year salary + $5K ATV

**Ground Inspection:**
- Robots: $3,500-$75K one-time
- Human: $15K/year salary + $3K equipment

**24/7 Security:**
- Robot + Drone: $4,600 one-time
- Human guard: $175K/year ($20/hr * 24/7)

**Payback:**
- Security: 2-3 weeks
- Farm operations: 6-12 months
- Overall: <1 year for full fleet

---

## Conclusion

Full FABRIC fleet provides:
- ✅ **Multi-modal coverage** (air + ground)
- ✅ **Coordinated operations** (robots + drones work together)
- ✅ **24/7 autonomous monitoring**
- ✅ **Emergency response** (all units converge)
- ✅ **Voice control** through Legacy AI
- ✅ **Scales efficiently** (one edge per location)

**The OpenMind OM1 architecture, applied to BOTH robots and drones, creates a truly unified autonomous fleet.**

---

*All documentation complete! Ready for Linear issue creation.*
