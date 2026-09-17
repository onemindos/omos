---
domain: 01
name: Sovereign Defense — Field Reference Card
status: active
owner: zeus
tags: [sovereign-defense, tactical, quick-reference, field, carry]
created: 2026-05-13
updated: 2026-05-13
---

# Field Reference Card

**PRINT. LAMINATE. CARRY.**
This is the only document you need in the field. Memorize the voice commands and button mappings. Everything else lives in [[wartime-weapon-protocol]].

---

## THREAT LEVELS — Quick Escalation

| Level | Color | State | Auto-Trigger |
| --- | --- | --- | --- |
| **5 — SOVEREIGN** | 🟢 | Normal ops, internet up | Default |
| **4 — VIGILANT** | 🔵 | Anomaly detected | Perimeter trip / unknown face |
| **3 — GUARDED** | 🟡 | Confirmed threat / comms degrading | Multiple triggers / internet flicker |
| **2 — TACTICAL** | 🟠 | Internet dead / active intruder | No internet / GPS denied |
| **1 — WARFIGHTER** | 🔴 | Under fire / life threat | Panic button / dead-man switch |

**Escalate:** Say `"WARFIGHTER"` or press panic button  
**De-escalate:** Say `"STAND DOWN [LEVEL NAME]"` (never automatic)

---

## VOICE COMMANDS (Local — No Internet Needed)

| Say This | What Happens |
| --- | --- |
| `"DEPLOY"` | Nearest grounded drone launches, orbits Zeus at 30m |
| `"DEPLOY ALL"` | All available drones launch in formation |
| `"OVERWATCH"` | Nearest drone climbs to max alt, wide camera, records |
| `"FOLLOW ME"` | Nearest drone follows Zeus GPS at 30m offset |
| `"ORBIT HERE"` | Nearest drone orbits Zeus's current position |
| `"EYES NORTH"` / `"EYES SOUTH"` / `"EYES EAST"` / `"EYES WEST"` | Drone points camera in that direction |
| `"LIGHTS ON"` | Drone spotlight + strobe activates |
| `"LIGHTS OFF"` | Drone lights off |
| `"RTH"` or `"COME HOME"` | All drones return to home |
| `"LAND"` | Nearest drone lands |
| `"LAND ALL"` | All drones emergency land |
| `"STATUS"` | Audio: battery levels, positions, threats detected |
| `"WARFIGHTER"` | Escalate to LEVEL 1, full autonomous protection |
| `"STAND DOWN [LEVEL]"` | De-escalate |

---

## PHYSICAL BUTTONS (RC Controller)

| Button | Action |
| --- | --- |
| **BTN 1** | DEPLOY (nearest drone orbits Zeus) |
| **BTN 2** | RTH ALL |
| **BTN 3** | FOLLOW ME |
| **BTN 4** | EMERGENCY LAND ALL |
| **Joystick** | Manual override on selected drone |
| **Panic (chest/belt)** | WARFIGHTER — full auto immediately |

---

## MESH COMMANDS (When Internet Is Down)

Send short strings on channel `DRONE-CMD` (ch1) via Meshtastic:

| Command | Effect |
| --- | --- |
| `D1:FLY:lat,lon,alt` | Drone 1 flies to GPS coord |
| `D1:ORBIT:lat,lon,30,100` | Drone 1 orbits point (30m alt, 100m radius) |
| `D1:FOLLOW:ZEUS` | Drone 1 follows Zeus GPS beacon |
| `D1:RTH` | Drone 1 returns home |
| `D*:RTH` | ALL drones return home |
| `D*:LAND` | ALL drones emergency land |
| `PANIC` | All drones execute WARFIGHTER protocol |

---

## PRE-LOADED PX4 MISSIONS (Run Without Any Connectivity)

| Mission | Trigger | What It Does |
| --- | --- | --- |
| **M01 GUARDIAN ORBIT** | Panic / dead-man | Launch → orbit Zeus GPS at 30m alt, indefinitely |
| **M02 PROPERTY PATROL** | Schedule / button | Fly pre-loaded fence around property |
| **M03 PERIMETER SWEEP** | Button / mesh | Visit each sensor location, hover, photo |
| **M04 FOLLOW BEACON** | Button | Follow Meshtastic GPS beacon (Zeus) |
| **M05 EMERGENCY OVERWATCH** | WARFIGHTER auto | Max altitude, wide camera, 360° orbit, record |
| **M06 RTH SAFE** | Any RTH | Return, descend, land, disarm |
| **M07 AREA DENIAL** | WARFIGHTER auto | Orbit threat bearing, siren + strobe |
| **M08 SEARCH GRID** | Mesh command | Expanding square search from GPS center |

---

## WARFIGHTER SEQUENCE (What Happens Automatically)

```
T+0s  → All grounded drones launch. Strobes + siren ON.
T+2s  → Formation: 1 drone at 50m overwatch, 2-3 at 20m orbit around Zeus
T+5s  → YOLO tags all humans: BLUE (known) / RED (unknown)
T+10s → All cameras recording to SD card
T+30s → Zeus GPS broadcasting on mesh ALERT channel (ch2)
        Rotating battery coverage (RTH at 20%, fresh drone launches)
```

Dead-man switch: Zeus's mesh node silent >30s → auto-WARFIGHTER

---

## COMMS FAILOVER CASCADE

```
Internet + Tailscale  → OpenClaw (cloud, full intel)
  ↓ if down
LTE / WiFi            → DimOS local (edge AI, no cloud)
  ↓ if down
Meshtastic LoRa       → Mesh commands (short strings, encrypted)
  ↓ if down
RC Controller         → Direct manual override (Herelink)
  ↓ if down
Dead-man switch       → PX4 runs M01 autonomously (no input needed)
```

---

## EagleEye HUD — What You're Seeing

| Symbol | Meaning |
| --- | --- |
| ● | Zeus position |
| ▲ | Drone (airborne) |
| ◼ | Drone (grounded/standby) |
| ✖ | Threat contact |
| ◎ | Home base |
| 🔵 | Known face (blue tag) |
| 🔴 | Unknown face (red tag) |
| ⚠️ | Weapon detected |

Status bar reads: `D1: ■ AIRBORNE 47m | 73% | ORBIT`  
Format: `[Drone]: [status] [altitude] | [battery] | [behavior]`

---

## HARDWARE AT A GLANCE

| Item | Where | When Active |
| --- | --- | --- |
| Meshtastic node (Heltec V3) | Body (EDC) | Always |
| Panic button | Chest / belt (EDC) | Always |
| iPhone + ATAK | Pocket | LEVEL 4+ |
| EagleEye HUD (Xreal + Jetson) | On face/vest | LEVEL 4+ |
| Herelink RC controller | Hand | LEVEL 3+ |
| ATAK tablet (ruggedized) | Chest rig | LEVEL 2+ |
| Throat mic + bone conduction | Neck | LEVEL 2+ |

---

## RULES — NEVER FORGET

**NEVER:**
- Fly above 400ft AGL (FAA limit)
- Ram or physically contact any person
- Carry weapons or weapon-like payloads
- Override geofence for any reason

**ALWAYS:**
- Record to SD card when airborne
- Broadcast GPS on mesh when airborne
- RTH when battery <15%
- RTH on signal loss >60s

**STAND DOWN requires your explicit voice or button command. It never happens automatically.**

---

*Full doctrine: [[wartime-weapon-protocol]] | Full suite: [[README]]*
