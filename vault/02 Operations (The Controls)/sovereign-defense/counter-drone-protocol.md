---
domain: 01
name: Sovereign Defense — Counter-Drone Protocol
status: active
owner: zeus
tags: [sovereign-defense, counter-drone, detection, hostile-drone, security]
created: 2026-05-13
updated: 2026-05-13
---

# Counter-Drone Protocol

**Detecting, identifying, and responding to hostile drones targeting Zeus or his properties.**

The wartime protocol assumes Zeus deploys the drone. This document covers the inverse: what happens when someone deploys a drone against Zeus.

---

## The Threat Landscape

### Who Uses Drones Aggressively

| Actor | Likely Platform | Purpose | Indicators |
| --- | --- | --- | --- |
| Paparazzi / stalkers | Consumer DJI (Mavic, Mini) | Surveillance, photos | Persistent hovering, following |
| Corporate espionage | DJI Enterprise or custom | Property/IP theft, surveillance | Night ops, thermal, long hover |
| Criminal surveillance | Consumer DJI | Casing a property before burglary | Perimeter mapping patterns |
| Law enforcement (warranted) | DJI Enterprise, Skydio | Lawful surveillance | Announce before / during if warranted |
| Cartel / organized crime | Consumer or custom | High-value target surveillance | Professional ops, multiple drones |
| Military/state actor | Custom (non-consumer) | High threat | Unlikely for civilian, but escalate immediately |

### Consumer Drone Signatures

Most threat actors use consumer drones. These are identifiable:

| Platform | RF Frequency | Visual Signature | Sound | Notes |
| --- | --- | --- | --- | --- |
| DJI Mavic 3 | 2.4GHz + 5.8GHz OcuSync 3 | Foldable, gray, 249g | 64dB @ 4ft | Most common |
| DJI Mini 4 Pro | 2.4GHz + 5.8GHz | Small, white, 249g | Quieter than Mavic | Cheap, widespread |
| DJI Matrice 350 | 2.4GHz + 5.8GHz | Large, enterprise, >1kg | Louder | More capable, less common |
| Custom FPV | 5.8GHz video + ELRS control | Small, fast, any color | Very loud | Harder to detect |
| Fixed wing | Varies | Slow, gliding, long endurance | Very quiet | Rare, surveillance-focused |

---

## Detection Methods

### Layer 1 — RF Detection (Best for Consumer Drones)

DJI drones broadcast RF on specific frequencies. RF detectors can identify DJI control signals and video feeds without active jamming.

**Recommended hardware (civilian legal):**

| Device | Cost | What It Detects | Range |
| --- | --- | --- | --- |
| **DJI AeroScope G16** | ~$5,000 | DJI drones only, ID, GPS position, pilot position | 5km |
| **DroneWatcher app** | Free / $199/yr | DJI drones via RF signature | ~1.5km |
| **Autel Evo Detection Module** | ~$2,000 | Autel + some DJI | ~3km |
| **DIY: RTL-SDR + DragonOS** | $30 (SDR dongle) | Raw RF spectrum, manual ID | ~500m |
| **OpenDroneID sniffer** | Free (software) | Remote ID (FAA-mandated for new drones) | ~3km |

**Lowest cost starting point:** RTL-SDR v3 dongle ($30) + DragonOS (free Linux image) + GNU Radio → run on spare Raspberry Pi → alerts when drone RF signatures detected on 2.4GHz/5.8GHz.

**OpenDroneID:** As of 2023, all new drones sold in the US must broadcast Remote ID (like an aircraft transponder). This broadcasts drone serial number, operator GPS, altitude, and velocity. Any OpenDroneID-compatible receiver can see this. The Meshtastic + ATAK system can be extended to display Remote ID contacts.

### Layer 2 — Acoustic Detection

Drone motors and propellers produce distinct acoustic signatures. Acoustic sensors can detect drones before they're visible.

| System | Cost | Range | Notes |
| --- | --- | --- | --- |
| **Dedrone DroneTracker** | ~$10,000 | 500m | Enterprise, full suite |
| **DIY acoustic array** | ~$200 | 200-400m | 4× microphones + ML classifier |
| **Shotgun mic + Jetson** | ~$150 | 100-200m | Single directional mic + YOLO audio |

DIY approach: 4 directional microphones at property corners + Raspberry Pi running [DeepMind's DroneDetect model](https://github.com/AlexanderNeuwirth/DiscovAir) → triangulate bearing when drone acoustic signature detected → alert via Home Assistant → send mesh alert.

### Layer 3 — Visual Detection (Cameras + YOLO)

Existing cameras on property feed into YOLO running on Jetson. YOLO v8 can be trained on drone images.

**Models:**
- `yolov8n` pretrained on COCO (includes 'aircraft' class — imprecise but works)
- Fine-tuned on drone images: [Drone Detection Dataset (Roboflow)](https://universe.roboflow.com/yolo-drone-detection) — 6,000+ labeled images, fine-tune in 2 hours on Jetson

**Camera placement for maximum sky coverage:**

Deploy fisheye cameras (180° FOV) pointed upward at property corners. Cheap options:
- ELP 180° fisheye USB camera (~$40) → USB to Jetson → YOLO running at 10fps

### Layer 4 — EagleEye HUD Integration

When wearing the EagleEye HUD, YOLO + acoustic fusion can tag detected drones in the AR overlay:
- Detected drone → red triangle overlay at detected bearing
- RF signal detected → orange dot overlay at signal direction
- Remote ID received → ID number + operator GPS displayed

---

## Detection Pipeline (Integrated)

```
┌──────────────────────────────────────────────────────────┐
│  SENSOR LAYER                                             │
│  RF SDR (2.4/5.8GHz) → drone RF signature               │
│  Acoustic array → drone sound signature                  │
│  Fisheye cameras → YOLO visual detection                 │
│  OpenDroneID receiver → Remote ID broadcast              │
└──────────────────┬───────────────────────────────────────┘
                   │ raw detections
                   ▼
┌──────────────────────────────────────────────────────────┐
│  FUSION LAYER (Jetson Orin)                              │
│  Correlate RF + acoustic + visual into single contact   │
│  Calculate bearing and estimated range                   │
│  Classify: DJI / FPV / unknown / fixed wing             │
│  Confidence score (low / medium / high)                  │
└──────────────────┬───────────────────────────────────────┘
                   │ fused contact
                   ▼
┌──────────────────────────────────────────────────────────┐
│  ALERT LAYER                                             │
│  Home Assistant → property alert (sensor trigger)       │
│  Meshtastic → mesh broadcast (bearing + type)           │
│  OpenClaw (if online) → Telegram alert to Zeus          │
│  EagleEye HUD → visual tag in AR overlay                │
│  ATAK → contact plotted on tactical map                  │
│  Auto-escalate OMOS DEFCON based on confidence          │
└──────────────────────────────────────────────────────────┘
```

---

## Response Protocol by Threat Level

### Low Confidence Detection (single sensor, brief)

- Log the detection
- Dispatch Zeus's drone to investigate bearing
- Continue monitoring
- No DEFCON escalation yet

### Medium Confidence (two sensors, sustained)

- **DEFCON escalates to LEVEL 4 (VIGILANT)**
- Zeus's drone dispatched to investigate at 50m altitude
- EagleEye camera points to detected bearing
- OpenClaw sends Telegram alert: "⚠️ Possible drone detected — bearing [X], type [DJI/unknown], confidence [%]"
- Begin recording all sensors

### High Confidence (three sensors, or visual lock)

- **DEFCON escalates to LEVEL 3 (GUARDED)**
- Zeus's drone attempts visual intercept (approach and photograph the hostile drone)
- Try to capture Remote ID broadcast → identify operator
- Record everything to SD card
- Notify law enforcement if threatening behavior observed
- Do NOT jam, net, or physically intercept hostile drone (see Legal below)

### Active Surveillance Confirmed (drone circling, hovering, following)

- **DEFCON escalates to LEVEL 2 (TACTICAL)**
- Full documentation: video, RF logs, Remote ID (if available), timestamps
- Call law enforcement with documentation
- Zeus moves to covered position
- Own drones provide overwatch on Zeus (not to chase hostile drone)
- Activate property perimeter lights and cameras

---

## What You CANNOT Do (Legal Limits)

| Action | Legal Status | Law |
| --- | --- | --- |
| Visually tracking and photographing hostile drone | LEGAL | Surveillance in public airspace |
| RF detection and logging | LEGAL | Passive monitoring |
| Actively jamming drone control signals | **ILLEGAL** | 47 U.S.C. § 333 — federal felony. No exceptions. |
| Shooting down a drone (gun, net, laser) | **ILLEGAL** | FAA — shooting an aircraft is a federal felony. Also state laws. |
| Hacking/taking control of hostile drone | **ILLEGAL** | Computer fraud laws |
| Intercepting with own drone physically | **ILLEGAL** | Destruction of aircraft |

**The frustrating truth:** Counter-drone intervention is currently reserved for federal agencies (FAA, DHS, DOD). No civilian in the US can legally jam, shoot down, or capture a hostile drone — even one actively surveilling you on your own property.

Your legal remedies:

1. **Document and report to FAA** (faa.gov/uas/resources/law_enforcement) — FAA can identify the operator from Remote ID or registration
2. **Report to law enforcement** — stalking/harassment laws may apply if drone follows you
3. **Civil action** — trespass to chattels, invasion of privacy, harassment
4. **Contact DJI** — report the serial number (from Remote ID) to DJI, they can geofence the operator or revoke access
5. **Hire a private investigator** — they can trace the operator from RF patterns and video analysis

---

## Countermeasures (Legal Only)

While you can't shoot it down or jam it, you can make your property harder to surveil:

| Countermeasure | Effect | Cost |
| --- | --- | --- |
| **High-lumen spotlights** aimed upward | Overwhelm drone camera with light | $100-500 |
| **Roof-mounted mirrors/reflective panels** | Reflect and distort aerial imagery | $50 |
| **IR floodlights** (IR-sensitive cameras get blinded) | Overwhelm thermal cameras | $200 |
| **Roof/canopy structures** | Block line of sight for aerial surveillance | $500+ |
| **Privacy film on windows** (blocks interior thermal) | Prevents through-window thermal | $200 |
| **Vegetation screening** (mature trees, tall hedges) | Natural barrier to aerial view | Long term |
| **Move activities indoors** | Deny surveillance opportunity | Free |
| **Anti-drone mesh canopy** (over sensitive outdoor areas) | Physical barrier | $500+ |

---

## Remote ID — Your Best Legal Tool

The FAA Remote ID rule (effective September 2023) requires all drones sold in the US to broadcast:

- Drone serial number (UAS ID)
- Operator's control station GPS position
- Drone GPS position, altitude, speed, heading
- Timestamp

You can receive this with a $30 OpenDroneID USB dongle or with compatible ATAK plugin.

**If you receive Remote ID from a suspicious drone:**

1. Log the UAS ID (serial number format: FAA-registered number)
2. Log the operator GPS position (where they're standing)
3. Share with law enforcement — this is location data for the person controlling the drone
4. FAA can cross-reference the serial number to the registered owner

Not all illegal drones will broadcast Remote ID (older drones, DIY, FPV don't need to). But most consumer DJI drones sold after 2023 do.

---

## Build: Minimum Counter-Drone Detection Stack

Cheapest working detection system:

| Component | Cost | Purpose |
| --- | --- | --- |
| RTL-SDR v3 dongle | $30 | RF detection (2.4/5.8GHz) |
| Raspberry Pi 4 (2GB) | $45 | Processing host |
| 2× ELP 180° fisheye USB cameras | $80 | Sky-facing visual |
| DragonOS FlightAware Pi image | Free | Pre-built SDR software stack |
| OpenDroneID software | Free | Remote ID reception |
| DroneWatcher (app tier) | Free | RF signature matching |
| ATAK with DroneKit plugin | Free | Contact display on map |
| **Total** | **~$155** | Full detection pipeline |

Mount cameras at highest point on property, antennas in open area away from metal obstructions, Pi in weatherproof box.

---

## Future Capability: Drone-on-Drone Intercept (Not Legal Yet)

For reference — what legal counter-drone systems look like if regulations change:

| System | Approach | Status |
| --- | --- | --- |
| **Net launchers** (DJI DJI AeroScope + net drone) | Mechanical net capture | Currently illegal for civilians |
| **RF spoofing** (GPS denial) | Confuse drone navigation | Illegal under FCC + FAA |
| **Trained eagles** (Guard From Above) | Biological interception | Legal in some countries, impractical |
| **Laser dazzlers** | Blind sensors | Illegal (eye hazard, aircraft targeting law) |
| **Pursuit and shadow** (drone follows hostile drone) | Documentation, not capture | LEGAL — no contact, stay in your airspace |

**Shadow protocol (currently legal):** Launch your drone to follow the hostile drone at safe distance (no contact), record everything from multiple angles. When it returns to its operator, you've documented the operator location. Feed to law enforcement.

---

*Full suite: [[README]] | Operating procedures: [[wartime-weapon-protocol]] | Legal limits: [[legal-framework]]*
