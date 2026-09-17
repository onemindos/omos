# Security Patrol Patterns with AI Threat Detection

> **Created:** 2026-01-22
> **Focus:** Skydio 2+ Autonomous Security
> **Integration:** OneMind OS + Legacy AI + Computer Vision

---

## Overview

The **Skydio 2+** is purpose-built for autonomous security operations:
- ✅ Best-in-class 360° obstacle avoidance (6 cameras)
- ✅ AI-powered tracking (follow intruders without GPS)
- ✅ Fully autonomous navigation (no manual piloting)
- ✅ 27 min flight time
- ✅ 4K60 HDR camera
- ✅ Cloud API for programmatic control

**This document provides:**
1. Security patrol patterns (perimeter, zone, random, threat-response)
2. AI threat detection (people, vehicles, animals, anomalies)
3. Emergency response workflows
4. Integration with Legacy for voice alerts
5. 24/7 autonomous operations

---

## Skydio 2+ Security Setup

### Hardware Configuration

| Component | Specification |
|-----------|---------------|
| **Base Station** | Auto-charging dock (optional, $500) |
| **Lights** | Strobe lights for night visibility |
| **Speaker** | External speaker for warnings |
| **Beacon** | GPS beacon for tracking mode |
| **Extra Batteries** | 3-4 batteries for continuous ops |

### Software Stack

```
┌──────────────────────────────────────────────────────┐
│  Security Operations Center                          │
│                                                       │
│  ┌────────────────┐      ┌────────────────────┐     │
│  │ Live Dashboard │      │  Alert Console     │     │
│  │  - Map view    │      │  - Threat detected │     │
│  │  - Video feed  │      │  - Dispatch drone  │     │
│  │  - Patrol path │      │  - Call authorities│     │
│  └────────────────┘      └────────────────────┘     │
│                                                       │
└───────────────────────┬───────────────────────────────┘
                        │
                        ▼
        ┌───────────────────────────────┐
        │  Legacy AI (Cloud)            │
        │  - Voice alerts               │
        │  - Threat classification      │
        │  - Response coordination      │
        └───────────────┬───────────────┘
                        │
                        ▼
        ┌───────────────────────────────┐
        │  OneMind Drone Manager (Edge) │
        │  - Skydio Cloud API bridge    │
        │  - Patrol orchestration       │
        │  - Video streaming            │
        └───────────────┬───────────────┘
                        │
                        ▼
                ┌───────────────┐
                │  Skydio 2+    │
                │  - Autonomous │
                │  - AI tracking│
                │  - 4K video   │
                └───────────────┘
```

---

## Patrol Pattern 1: Perimeter Patrol

### Use Case: Regular Property Boundary Checks

**Frequency:** Every 2 hours, 24/7
**Duration:** 15-20 minutes
**Coverage:** Complete property perimeter

```python
# agents/legacy/workflows/security/perimeter_patrol.py

from agno import workflow
from datetime import datetime
import asyncio

@workflow
async def perimeter_patrol(
    property_name: str = "main_property",
    altitude: float = 20.0,  # meters
    speed: float = 5.0,  # m/s
):
    """
    Automated perimeter patrol with threat detection

    Flies around property boundary, analyzing video for:
    - People entering/exiting
    - Vehicles near gates
    - Animals in prohibited areas
    - Fence damage
    - Unusual activity

    Args:
        property_name: Property identifier
        altitude: Flight altitude in meters
        speed: Flight speed in m/s
    """

    legacy.speak("Starting perimeter patrol")

    # ============================================
    # PHASE 1: Pre-Flight Checks
    # ============================================

    # Get property boundaries
    property_data = await db.get_property(property_name)
    perimeter_points = property_data.perimeter

    # Check drone status
    drone_status = await drone_tool.get_drone_status("security")

    if drone_status.battery_percent < 50:
        legacy.speak("Security drone battery low. Charging required.")
        return {"status": "aborted", "reason": "low_battery"}

    # Check weather (wind > 20 mph = abort)
    weather = await sensors.get_weather()
    if weather.wind_speed > 20:
        legacy.speak("Wind too high for patrol. Patrol postponed.")
        return {"status": "postponed", "reason": "high_wind"}

    # ============================================
    # PHASE 2: Generate Patrol Path
    # ============================================

    patrol_path = []

    # Add perimeter waypoints
    for point in perimeter_points:
        patrol_path.append({
            "latitude": point.lat,
            "longitude": point.lon,
            "altitude": altitude,
            "speed": speed,
            "gimbal_pitch": -30,  # Look down 30° from horizon
            "hover_time": 3,  # Hover 3 seconds at each waypoint
        })

    # Close the loop (return to start)
    patrol_path.append(patrol_path[0])

    # ============================================
    # PHASE 3: Execute Patrol
    # ============================================

    patrol_results = {
        "start_time": datetime.now(),
        "property": property_name,
        "waypoints": len(patrol_path),
        "threats_detected": [],
        "alerts_sent": 0,
        "video_clips": [],
    }

    # Upload mission to Skydio
    await drone_tool.command_drone(
        "security",
        "mission",
        {
            "waypoints": patrol_path,
            "video_recording": True,
            "stream_to_livekit": True,
            "obstacle_avoidance": "max",  # Maximum obstacle avoidance
        }
    )

    # Takeoff
    await drone_tool.command_drone("security", "takeoff", {"altitude": altitude})

    legacy.update_status("Perimeter patrol in progress")

    # ============================================
    # PHASE 4: Real-Time Threat Analysis
    # ============================================

    current_waypoint = 0
    video_buffer = []

    while True:
        # Get drone status
        status = await drone_tool.get_drone_status("security")

        # Check if mission complete
        if not status.mission_active:
            break

        # Get video frame from LiveKit stream
        video_frame = await livekit.get_latest_frame("drone-security-feed")

        if video_frame:
            video_buffer.append(video_frame)

            # Analyze every 2 seconds (reduce API calls)
            if len(video_buffer) >= 60:  # 30 fps * 2 sec
                # AI analysis with Claude (vision)
                analysis = await legacy.analyze_video_frames(
                    frames=video_buffer,
                    objectives=[
                        "detect_people",
                        "detect_vehicles",
                        "detect_animals",
                        "detect_fence_damage",
                        "detect_unusual_activity",
                    ],
                    current_location={
                        "lat": status.latitude,
                        "lon": status.longitude,
                        "waypoint": current_waypoint,
                        "time": datetime.now(),
                    }
                )

                # Check for threats
                if analysis.threats_detected:
                    for threat in analysis.threats:
                        legacy.speak(
                            f"THREAT DETECTED: {threat.type} at {threat.location_name}"
                        )

                        # Save threat details
                        threat_record = {
                            "type": threat.type,
                            "severity": threat.severity,  # low, medium, high, critical
                            "location": {
                                "lat": status.latitude,
                                "lon": status.longitude,
                                "waypoint": current_waypoint,
                            },
                            "description": threat.description,
                            "timestamp": datetime.now(),
                            "video_clip": await save_video_clip(video_buffer),
                        }

                        patrol_results["threats_detected"].append(threat_record)

                        # Response based on severity
                        if threat.severity in ["high", "critical"]:
                            # CRITICAL: Stop patrol, investigate
                            await drone_tool.command_drone("security", "pause")

                            legacy.speak(
                                f"Critical threat: {threat.description}. "
                                f"Investigating now."
                            )

                            # Hover and zoom in
                            await drone_tool.command_drone(
                                "security",
                                "orbit",
                                {
                                    "radius": 10,  # 10 meter orbit
                                    "altitude": altitude,
                                    "duration": 30,  # 30 seconds
                                }
                            )

                            # Capture high-res photos
                            photos = await drone_tool.command_drone(
                                "security", "capture_photos", {"count": 5}
                            )

                            # Send alert
                            await send_critical_alert(threat_record, photos)
                            patrol_results["alerts_sent"] += 1

                            # Wait for user decision
                            legacy.speak(
                                "What would you like me to do? "
                                "Options: continue patrol, keep tracking, land and wait, "
                                "or call authorities."
                            )

                            response = await wait_for_voice_command(timeout=60)

                            if "continue" in response:
                                await drone_tool.command_drone("security", "resume")
                            elif "track" in response:
                                await initiate_tracking_mode(threat_record)
                            elif "land" in response:
                                await drone_tool.command_drone("security", "land")
                                return patrol_results
                            elif "authorities" in response:
                                await call_authorities(threat_record)
                                await drone_tool.command_drone("security", "resume")

                        elif threat.severity == "medium":
                            # MEDIUM: Log and continue, send notification
                            await notifications.send_push(
                                title=f"Security Alert: {threat.type}",
                                body=threat.description,
                                image=await save_video_clip(video_buffer),
                            )
                            patrol_results["alerts_sent"] += 1

                        else:  # low severity
                            # LOW: Just log
                            pass

                # Clear buffer
                video_buffer = []

        # Track waypoint progress
        waypoint_progress = status.mission_progress / 100 * len(patrol_path)
        if int(waypoint_progress) > current_waypoint:
            current_waypoint = int(waypoint_progress)
            legacy.update_status(
                f"Perimeter patrol: waypoint {current_waypoint}/{len(patrol_path)}"
            )

        await asyncio.sleep(0.5)  # Check every 0.5 seconds

    # ============================================
    # PHASE 5: Return and Report
    # ============================================

    await drone_tool.command_drone("security", "rtl")

    patrol_results["end_time"] = datetime.now()
    patrol_results["duration_minutes"] = (
        patrol_results["end_time"] - patrol_results["start_time"]
    ).seconds / 60

    # Generate report
    if patrol_results["threats_detected"]:
        report = f"""
        Perimeter Patrol Complete - THREATS DETECTED

        Property: {property_name}
        Start: {patrol_results['start_time'].strftime('%H:%M:%S')}
        End: {patrol_results['end_time'].strftime('%H:%M:%S')}
        Duration: {patrol_results['duration_minutes']:.1f} minutes

        THREATS DETECTED: {len(patrol_results['threats_detected'])}
        """

        for threat in patrol_results["threats_detected"]:
            report += f"""
        - {threat['type'].upper()} ({threat['severity']} severity)
          Location: {threat['location']['waypoint']}
          Time: {threat['timestamp'].strftime('%H:%M:%S')}
          Description: {threat['description']}
          Video: {threat['video_clip']}
        """

        legacy.speak(
            f"Perimeter patrol complete. {len(patrol_results['threats_detected'])} "
            f"threats detected. Review report for details."
        )

    else:
        report = f"""
        Perimeter Patrol Complete - All Clear

        Property: {property_name}
        Start: {patrol_results['start_time'].strftime('%H:%M:%S')}
        End: {patrol_results['end_time'].strftime('%H:%M:%S')}
        Duration: {patrol_results['duration_minutes']:.1f} minutes

        No threats detected. Property secure.
        """

        legacy.speak("Perimeter patrol complete. All clear.")

    # Save to database
    await db.save_patrol_record(patrol_results)

    # Send report
    await notifications.send_email(
        to=user.email,
        subject=f"Security Patrol Report: {property_name}",
        body=report,
    )

    return patrol_results


# ============================================
# Helper Functions
# ============================================

async def initiate_tracking_mode(threat_record):
    """
    Switch to AI tracking mode (Skydio's best feature)

    Autonomously follows subject (person, vehicle, animal)
    """
    legacy.speak("Switching to tracking mode")

    # Skydio AI tracking
    await drone_tool.command_drone(
        "security",
        "track",
        {
            "subject": threat_record["type"],  # "person", "vehicle"
            "mode": "follow",  # or "orbit", "spotlight"
            "distance": 10,  # meters behind
            "altitude_offset": 5,  # 5m above subject
        }
    )

    legacy.speak("Now tracking target. Recording video.")

    # Continue tracking until told to stop or target lost
    # (implemented in separate workflow)


async def send_critical_alert(threat_record, photos):
    """
    Send immediate critical alert with all channels
    """

    # Push notification
    await notifications.send_push(
        title=f"🚨 CRITICAL SECURITY ALERT",
        body=threat_record["description"],
        image=photos[0] if photos else None,
        sound="alarm",
        priority="critical",
    )

    # SMS
    await notifications.send_sms(
        to=user.phone,
        message=f"CRITICAL SECURITY ALERT: {threat_record['description']} "
                f"at {threat_record['timestamp'].strftime('%H:%M:%S')}. "
                f"Check app for details."
    )

    # Voice call (optional)
    if threat_record["severity"] == "critical":
        await notifications.make_voice_call(
            to=user.phone,
            message=f"This is your OneMind Security System. "
                    f"Critical security threat detected: {threat_record['description']}. "
                    f"Check your app immediately."
        )

    # Email with photos
    await notifications.send_email(
        to=user.email,
        subject="🚨 CRITICAL SECURITY ALERT",
        body=f"""
        CRITICAL SECURITY THREAT DETECTED

        Type: {threat_record['type']}
        Severity: {threat_record['severity']}
        Time: {threat_record['timestamp']}
        Location: {threat_record['location']}

        Description:
        {threat_record['description']}

        Video Clip: {threat_record['video_clip']}

        Please review immediately and take appropriate action.
        """,
        attachments=photos,
    )


async def call_authorities(threat_record):
    """
    Automated 911 call (with user confirmation)
    """

    legacy.speak(
        "Are you sure you want me to call 911? Say 'yes call 911' to confirm."
    )

    confirmation = await wait_for_voice_command(timeout=30)

    if "yes" in confirmation and "911" in confirmation:
        # Use Twilio or similar to place call
        await twilio.make_call(
            to="911",
            from_number=user.phone,
            message=f"This is an automated security system at {property_data.address}. "
                    f"A {threat_record['type']} has been detected on the property. "
                    f"The property owner has been notified. "
                    f"This call is being recorded."
        )

        legacy.speak("Calling 911 now.")

        # Log emergency call
        await db.create_incident_record({
            "type": "emergency_call",
            "threat": threat_record,
            "call_time": datetime.now(),
        })

    else:
        legacy.speak("911 call cancelled.")
```

---

## Patrol Pattern 2: Zone Patrol

### Use Case: Focus on High-Priority Areas

**Zones:**
- Entry gates
- Barn/equipment storage
- Livestock areas
- High-value crops

```python
@workflow
async def zone_patrol(zones: List[str], dwell_time: int = 60):
    """
    Patrol specific zones in sequence

    Args:
        zones: List of zone names ["front_gate", "barn", "greenhouse"]
        dwell_time: Seconds to spend analyzing each zone
    """

    legacy.speak(f"Starting zone patrol. {len(zones)} zones to check.")

    for zone_name in zones:
        zone = await db.get_zone(zone_name)

        legacy.speak(f"Checking {zone_name}")

        # Fly to zone
        await drone_tool.command_drone(
            "security",
            "goto",
            {
                "latitude": zone.center_lat,
                "longitude": zone.center_lon,
                "altitude": zone.optimal_altitude,
            }
        )

        # Orbit zone for analysis
        await drone_tool.command_drone(
            "security",
            "orbit",
            {
                "radius": zone.radius,
                "altitude": zone.optimal_altitude,
                "duration": dwell_time,
                "gimbal_mode": "point_inward",  # Look at zone center
            }
        )

        # Capture video and analyze
        video = await capture_video_clip(duration=dwell_time)

        analysis = await legacy.analyze_video(
            video,
            objectives=zone.security_objectives  # Zone-specific checks
        )

        if analysis.issues:
            legacy.speak(f"Issues detected in {zone_name}: {len(analysis.issues)}")
            # Handle issues...

    legacy.speak("Zone patrol complete.")
```

---

## Patrol Pattern 3: Random Patrol

### Use Case: Unpredictable Deterrent

**Strategy:** Random waypoints within property bounds to prevent pattern recognition by potential intruders.

```python
@workflow
async def random_patrol(duration_minutes: int = 30):
    """
    Random waypoint patrol (unpredictable pattern)

    Makes it impossible for intruders to predict patrol route
    """
    import random

    property_bounds = await db.get_property_bounds()

    end_time = datetime.now() + timedelta(minutes=duration_minutes)

    await drone_tool.command_drone("security", "takeoff", {"altitude": 20})

    while datetime.now() < end_time:
        # Generate random waypoint within bounds
        random_point = generate_random_point_in_polygon(property_bounds)

        # Fly to random point
        await drone_tool.command_drone(
            "security",
            "goto",
            {
                "latitude": random_point.lat,
                "longitude": random_point.lon,
                "altitude": random.randint(15, 25),  # Random altitude
            }
        )

        # Hover and scan
        await asyncio.sleep(random.randint(10, 30))  # Random dwell time

        # Capture and analyze
        video = await capture_video_clip(duration=10)
        analysis = await legacy.analyze_video(video, objectives=["detect_threats"])

        if analysis.threats:
            # Switch to tracking mode
            await initiate_tracking_mode(analysis.threats[0])
            break

    await drone_tool.command_drone("security", "rtl")
```

---

## Patrol Pattern 4: Scheduled Coverage

### Use Case: Complete Property Scan

**Grid pattern for 100% coverage** of large properties.

```python
@workflow
async def complete_property_scan(property_name: str):
    """
    Systematic grid scan of entire property

    Ensures every square meter is visually inspected
    """

    property_data = await db.get_property(property_name)

    # Generate grid pattern (similar to agricultural scan)
    grid_pattern = generate_grid_pattern(
        boundaries=property_data.perimeter,
        altitude=30,  # meters (higher for overview)
        overlap=50,  # 50% overlap
        camera_fov=90,  # Skydio 2+ FOV
    )

    legacy.speak(
        f"Starting complete property scan. "
        f"{len(grid_pattern)} waypoints, approximately "
        f"{len(grid_pattern) * 10 / 60:.0f} minutes."
    )

    # Execute grid mission
    await drone_tool.command_drone(
        "security",
        "mission",
        {
            "waypoints": grid_pattern,
            "video_recording": True,
            "capture_interval": 5,  # Photo every 5 seconds
        }
    )

    # ... execute and analyze
```

---

## AI Threat Detection

### Computer Vision Analysis

```python
# agents/legacy/tools/threat_detection.py

async def analyze_security_video(
    video_frames: List[Image],
    context: dict
) -> ThreatAnalysis:
    """
    Analyze video for security threats using Claude (vision)

    Returns detailed threat analysis with classification
    """

    prompt = f"""
    You are analyzing security camera footage from a drone.

    Context:
    - Location: {context['location_name']}
    - Time: {context['timestamp']}
    - Weather: {context['weather']}
    - Known authorized personnel: {context['authorized_people']}
    - Known vehicles: {context['authorized_vehicles']}

    Analyze these video frames for:
    1. **People**: Count, description, behavior, authorized/unauthorized
    2. **Vehicles**: Type, color, license plate (if visible), authorized/unauthorized
    3. **Animals**: Type, location, threat level (predator vs harmless)
    4. **Structures**: Damage to fences, buildings, gates
    5. **Unusual Activity**: Anything suspicious or out of place

    For each detected element, classify threat level:
    - CRITICAL: Immediate danger (armed person, break-in, fire)
    - HIGH: Potential danger (unauthorized person, damaged fence)
    - MEDIUM: Warrants attention (unknown vehicle, large animal)
    - LOW: Informational (authorized person, known vehicle)

    Return structured JSON with detections and threat assessments.
    """

    response = await claude.analyze_images(
        images=video_frames,
        prompt=prompt,
        model="claude-opus-4-5"
    )

    # Parse response into structured format
    analysis = ThreatAnalysis(
        threats_detected=response.threats,
        people_count=response.people_count,
        vehicles_count=response.vehicles_count,
        animals_detected=response.animals,
        damage_detected=response.damage,
        overall_threat_level=response.overall_threat_level,
    )

    return analysis
```

### Threat Classification Examples

**Critical Threats:**
- Armed person
- Break-in attempt (forcing gate, cutting fence)
- Fire/smoke detected
- Vehicle ramming fence/gate
- Multiple unauthorized people

**High Threats:**
- Unauthorized person on property
- Unknown vehicle near structures
- Predator animal near livestock
- Major structural damage

**Medium Threats:**
- Unknown vehicle on road near property
- Large animal (deer, coyote) on property
- Minor fence damage
- Person near property line (not trespassing)

**Low Threats:**
- Authorized person/vehicle
- Small animals (rabbits, birds)
- Natural changes (fallen branch)

---

## Emergency Response Workflows

### Response 1: Intruder Detection

```python
@workflow
async def intruder_response(intruder_location: dict):
    """
    Automated response to intruder detection

    1. Track intruder with drone
    2. Alert owner
    3. Announce presence (deterrent)
    4. Record evidence
    5. Call authorities if needed
    """

    legacy.speak("INTRUDER DETECTED. Initiating response protocol.")

    # Send immediate alert
    await send_critical_alert({
        "type": "intruder",
        "severity": "critical",
        "location": intruder_location,
        "description": "Unauthorized person on property",
        "timestamp": datetime.now(),
    }, [])

    # Track intruder
    await drone_tool.command_drone(
        "security",
        "track",
        {
            "subject": "person",
            "mode": "spotlight",  # Keep them lit up
            "altitude_offset": 10,
        }
    )

    # Announce via speaker (if equipped)
    await drone_tool.command_drone(
        "security",
        "play_audio",
        {
            "message": "This is private property. You are being recorded. "
                      "Leave immediately or authorities will be called."
        }
    )

    # Flash lights (deterrent)
    await drone_tool.command_drone("security", "strobe_lights", {"enabled": True})

    # Wait for response
    legacy.speak(
        "Tracking intruder. Should I call 911? "
        "Say 'yes call 911' to confirm."
    )

    response = await wait_for_voice_command(timeout=60)

    if "yes" in response and "911" in response:
        await call_authorities({
            "type": "intruder",
            "location": intruder_location,
            "timestamp": datetime.now(),
        })

    # Continue tracking until intruder leaves
    while True:
        status = await drone_tool.get_drone_status("security")

        if not status.tracking_active:
            legacy.speak("Intruder has left the property.")
            break

        await asyncio.sleep(5)

    # Return to base
    await drone_tool.command_drone("security", "rtl")

    legacy.speak("Intruder response complete. Video evidence saved.")
```

### Response 2: Fire Detection

```python
@workflow
async def fire_response(fire_location: dict):
    """
    Automated response to fire/smoke detection
    """

    legacy.speak("FIRE DETECTED. Initiating emergency protocol.")

    # Immediate 911 call (fire = always call)
    await twilio.make_call(
        to="911",
        message=f"This is an automated alert. Fire detected at {property_address}. "
                f"Approximate location: {fire_location['description']}."
    )

    # Alert owner
    await send_critical_alert({
        "type": "fire",
        "severity": "critical",
        "location": fire_location,
        "description": "Fire or smoke detected on property",
        "timestamp": datetime.now(),
    }, [])

    # Orbit fire to assess (safe distance)
    await drone_tool.command_drone(
        "security",
        "orbit",
        {
            "latitude": fire_location["lat"],
            "longitude": fire_location["lon"],
            "radius": 50,  # Stay far away (heat/smoke)
            "altitude": 30,  # Higher altitude
            "duration": 60,
        }
    )

    # Capture thermal imaging (if equipped)
    # Useful for firefighters

    legacy.speak("Emergency services have been notified. Monitoring fire from safe distance.")
```

---

## 24/7 Autonomous Operations

### Automated Schedule

```python
# Automated patrol schedule (runs 24/7)

PATROL_SCHEDULE = {
    "00:00": {"type": "perimeter", "frequency": "every_2_hours"},
    "06:00": {"type": "complete_scan", "frequency": "daily"},
    "12:00": {"type": "zone", "zones": ["livestock", "equipment"]},
    "18:00": {"type": "perimeter", "priority": "gates"},
    "22:00": {"type": "random", "duration": 30},  # Night deterrent
}

@workflow
async def security_daemon():
    """
    24/7 security patrol daemon

    Runs continuously, executing scheduled patrols
    """

    while True:
        now = datetime.now()
        current_time = now.strftime("%H:%M")

        # Check scheduled patrols
        if current_time in PATROL_SCHEDULE:
            patrol_config = PATROL_SCHEDULE[current_time]

            if patrol_config["type"] == "perimeter":
                await perimeter_patrol()

            elif patrol_config["type"] == "complete_scan":
                await complete_property_scan("main_property")

            elif patrol_config["type"] == "zone":
                await zone_patrol(patrol_config["zones"])

            elif patrol_config["type"] == "random":
                await random_patrol(duration_minutes=patrol_config["duration"])

        # Check for ad-hoc alerts (motion sensors, etc.)
        alerts = await check_sensor_alerts()
        if alerts:
            for alert in alerts:
                await investigate_alert(alert)

        # Sleep until next check (1 minute intervals)
        await asyncio.sleep(60)
```

---

## Voice Commands

```
"Start perimeter patrol"
"Check the front gate"
"Investigate that noise" (dispatches drone to your location)
"What's the security status?"
"Show me the security feed"
"Track that person"
"Land the security drone"
"Call 911"
"Turn off security patrols for 1 hour" (when you have guests)
```

---

## ROI Analysis

### Skydio 2+ ($1,099)

**vs. Human Security Guard:**
- Security guard: $20/hr × 24/7 = $480/day = $14,400/month
- Skydio 2+: $1,099 one-time + $0/month (no subscription)
- **Payback: 2-3 days** (if replacing 24/7 guard)

**vs. Security Camera System:**
- CCTV cameras: $2,000-5,000 (fixed positions)
- Skydio: $1,099 (mobile, 360° coverage, AI tracking)
- **Advantage:** Drones can investigate vs cameras just record

---

## Conclusion

Security patrol with Skydio 2+ provides:
- ✅ Autonomous 24/7 monitoring
- ✅ AI threat detection and classification
- ✅ Immediate alerts and response
- ✅ Evidence recording (4K video)
- ✅ Active deterrent (visible presence)
- ✅ Cost-effective vs human guards

**All controlled by Legacy AI** - "Legacy, check the perimeter" initiates full autonomous patrol with threat detection.

---

*Next: Full FABRIC Fleet (Ground Robots + Drones)*
