---
domain: 01
name: Sovereign Defense — Hardware BOM & Acquisition Plan
status: active
owner: zeus
tags: [sovereign-defense, hardware, bom, acquisition, drones, robotics]
created: 2026-05-13
updated: 2026-05-13
---

# Hardware BOM & Acquisition Plan

**What to buy, in what order, from where, at what budget tier. Phased across 4 stages so the system becomes operational before it's complete.**

---

## Phase 0 — Foundation (Under $200 | Buy Now)

Get comms and awareness up immediately. No drone hardware needed yet. This works alongside your existing DJI Mavic and Mac/DimOS setup.

| # | Item | Est. Cost | Source | Purpose |
| --- | --- | --- | --- | --- |
| 1 | Heltec V3 LoRa32 ×5 | $125 ($25 ea) | AliExpress / Amazon | Meshtastic mesh nodes (body + drones + base + vehicle) |
| 2 | 18650 battery holders ×5 | $15 | Amazon | Power for mesh nodes |
| 3 | 3D-printed enclosures ×5 | $10 filament | Print yourself | Weatherproof node housing |
| 4 | Chest-mount panic button (BLE) | $15 | Amazon (BLE button) | WARFIGHTER trigger EDC |
| 5 | ATAK install on Android tablet | Free | F-Droid / atak.app | Tactical map, blue force tracking |
| **Phase 0 Total** | | **~$165** | | |

**Milestone:** Meshtastic mesh operational between Zeus + 4 nodes. ATAK showing Zeus position. Manual panic button wired. Mesh command flow testable.

---

## Phase 1 — Tactical Visibility (Under $1,200 | 1-3 Months)

Get the EagleEye HUD and edge compute online. DimOS continues handling drone runtime on existing DJI hardware.

| # | Item | Est. Cost | Source | Purpose |
| --- | --- | --- | --- | --- |
| 6 | Xreal Air 2 Pro | $449 | xreal.com / Amazon | AR display for HUD |
| 7 | NVIDIA Jetson Orin Nano 8GB Dev Kit | $249 | NVIDIA / Arrow Electronics | Body-worn edge AI (voice + vision) |
| 8 | MicroSD card 256GB (×2) | $30 | Amazon | Jetson OS + model storage |
| 9 | USB-C hub (5-port) | $25 | Amazon | Peripheral connections to Jetson |
| 10 | 20,000mAh slim battery pack | $60 | Anker / Amazon | Body-worn power for Jetson + HUD |
| 11 | Ruggedized Android tablet (used) | $150 | eBay / Swappa | Dedicated ATAK device |
| 12 | Throat mic + bone conduction earpiece | $50 | Amazon / Aliexpress | Hands-free voice commands |
| **Phase 1 Total** | | **~$1,013** | | |

**Milestone:** EagleEye HUD Phase 1 operational (see OMOS-EagleEye build phases). Jetson running Whisper.cpp locally. Voice commands processed off-device without internet. Drone camera feed in HUD picture-in-picture.

---

## Phase 2 — Custom PX4 Fleet (Under $3,000 | 3-6 Months)

First non-DJI autonomous drone. This is the core of the wartime weapon — no geofencing, no kill switch, pre-loaded PX4 missions.

### Drone 1 — Custom PX4 Quad Build

| # | Item | Est. Cost | Source | Purpose |
| --- | --- | --- | --- | --- |
| 13 | Frame: iFlight Nazgul5 or DJI F450 clone | $80 | GetFPV / RaceDayQuads | Main structural frame |
| 14 | Flight controller: Holybro Pixhawk 6C | $250 | Holybro / Amazon | PX4 autopilot (industry standard) |
| 15 | 4× BLDC motors (T-Motor F40 or similar) | $120 | GetFPV | Propulsion |
| 16 | 4× 30A ESC or 4-in-1 ESC | $80 | GetFPV | Motor controllers |
| 17 | GPS module: Here3+ GPS/Compass | $120 | Holybro | Position hold, autonomous nav |
| 18 | Telemetry: RFD900x ×2 (ground + air) | $200 | RFDesign | Long-range encrypted telemetry |
| 19 | RC receiver: ExpressLRS (ELRS) 2.4GHz | $30 | Aliexpress | Manual RC link |
| 20 | Companion computer: Jetson Orin Nano (drone-mounted) | $249 | NVIDIA | Edge AI on drone (YOLO, face rec) |
| 21 | Camera: ELP 1080p USB | $35 | Amazon | Main camera |
| 22 | Camera: FLIR Lepton 3.5 (thermal module) | $300 | FLIR / GroupGets | Thermal vision |
| 23 | Depth camera: Intel RealSense D435 | $250 | Intel / Amazon | Obstacle avoidance |
| 24 | Siren module: 120dB piezo | $15 | Amazon | Audio deterrence |
| 25 | LED strobe array: 12V high-intensity | $25 | Amazon | Visual deterrence |
| 26 | LiPo battery: 4S 5000mAh ×4 | $160 ($40 ea) | Hobbyking | Flight endurance (hot-swap) |
| 27 | LiPo charger: ISDT Q6 Plus | $50 | Amazon | Fast charging |
| 28 | Heltec V3 (drone mesh node) | $25 | AliExpress | Mesh command receiver on drone |
| **Drone 1 Total** | | **~$1,989** | | |

### Controller Setup

| # | Item | Est. Cost | Source | Purpose |
| --- | --- | --- | --- | --- |
| 29 | Herelink RC controller | $499 | Holybro | Professional RC + programmable buttons |
| 30 | QGroundControl license | Free | qgroundcontrol.com | PX4 mission planning software |
| **Controller Total** | | **~$499** | | |

**Phase 2 Total: ~$2,488**

**Milestone:** First fully autonomous PX4 drone airborne. M01-M08 missions loaded and tested. Meshtastic → MAVLink command bridge working. Voice commands control drone without internet. Herelink physical buttons mapped to missions.

---

## Phase 3 — Full Fleet (Under $5,000 Additional | 6-12 Months)

Scale to 2-4 drones for formation flying, rotation during sustained ops, and redundancy.

| # | Item | Est. Cost | Source | Purpose |
| --- | --- | --- | --- | --- |
| 31 | Drone 2 (same build as Drone 1) | ~$1,800 | As above | Formation pair |
| 32 | Drone 3 (same build) | ~$1,800 | As above | Third formation slot |
| 33 | EagleEye HUD Phase 2-3 additions | ~$1,060 | See BOM below | Thermal + UWB for HUD |
| 34 | Solar Meshtastic repeater (property-mounted) | $40 | Heltec + solar panel | Extended mesh coverage |
| 35 | Ruggedized vest/chest rig with mounts | $150 | Amazon / Condor | Organized carry system |
| 36 | Deadman switch module (UART + GPS) | $35 | Custom / Adafruit | Auto-WARFIGHTER trigger |
| **Phase 3 Total** | | **~$4,885** | | |

**Milestone:** 3-drone fleet operating in formation. WARFIGHTER drill completed — panic button → all 3 launch in <5 seconds → guardian orbit. Dead-man switch tested. Rotation behavior verified (RTH at 20% battery, fresh drone launches).

---

## Phase 4 — EagleEye Thermal + Radar (Parallel Track | 6-12 Months)

These can be built in parallel with the drone fleet. See [[OMOS-EagleEye]] for full architecture.

| # | Item | Est. Cost | Source | Purpose |
| --- | --- | --- | --- | --- |
| 37 | FLIR Boson+ 320 thermal | $800 | FLIR.com | Through-wall heat signatures |
| 38 | Acconeer XM125 UWB radar ×2 | $100 | acconeer.com | Through-wall movement detection |
| 39 | Sony STARVIS IMX327 (night vision cam) | $80 | AliExpress | Low-light full scene |
| 40 | Novelda X4 (UWB heartbeat) | $80 | Novelda | Breathing/heartbeat through walls |
| **EagleEye Phase 2-3 Total** | | **~$1,060** | | |

---

## Budget Summary

| Phase | When | Cost | Cumulative | What You Get |
| --- | --- | --- | --- | --- |
| **Phase 0** | Now | $165 | $165 | Mesh network + panic button + ATAK |
| **Phase 1** | 1-3 months | $1,013 | $1,178 | EagleEye HUD + edge AI + voice commands |
| **Phase 2** | 3-6 months | $2,488 | $3,666 | First autonomous PX4 drone + Herelink |
| **Phase 3** | 6-12 months | $4,885 | $8,551 | 3-drone fleet + WARFIGHTER drill ready |
| **EagleEye 2-3** | 6-12 months | $1,060 | $9,611 | Thermal + radar HUD |

**Minimum to be combat-capable:** Phase 0 + Phase 2 = ~$2,653  
**Full sovereign defense stack:** ~$9,600

---

## Sourcing Notes

| Source | Best For | Notes |
| --- | --- | --- |
| **GetFPV** | FPV drone parts | US-based, fast shipping, good selection |
| **Holybro** | Pixhawk FC, GPS, Herelink | Official PX4 hardware manufacturer |
| **RFDesign** | RFD900x telemetry | Australian, long-range encrypted |
| **NVIDIA** | Jetson dev kits | Arrow Electronics for reseller pricing |
| **AliExpress** | Heltec modules, motors, basic parts | 3-4 week shipping, cheap |
| **GroupGets** | FLIR Lepton, thermal modules | Group buys for thermal sensors |
| **Amazon** | Cables, power, accessories | Prime for fast turnaround |
| **eBay / Swappa** | Used Android tablets | $100-150 for dedicated ATAK device |

---

## What NOT to Buy

| Item | Why Not |
| --- | --- |
| **DJI flight controller** | Proprietary, geofenced, no PX4 support, Chinese kill switch |
| **DJI goggles for HUD** | Closed ecosystem, no AR overlay, no sensor fusion |
| **goTenna Pro X2** | $300/node — Meshtastic does 90% at 8% of the cost |
| **Skydio X10** | $10,000+ — overkill for Phase 1-3 |
| **Shield AI Hivemind** | Defense contractor only, not available to civilians |
| **Any DJI drone for tactical** | Same reasons as DJI FC above |
| **Consumer drones (Parrot, Autel consumer)** | No PX4 support, closed autopilot |

**Autel EVO Max 4T** is the one exception — $6,500, no CCP kill switch, AES-256 encrypted link, thermal built in. Worth considering if budget allows and you want turnkey vs custom build.

---

## Phase 0 First Purchase Checklist

- [ ] Order 5× Heltec V3 LoRa32 from AliExpress (3-4 week wait) or Amazon (faster, $5/unit more)
- [ ] Download Meshtastic app (iOS/Android) and firmware flasher
- [ ] Download ATAK from atak.app (free for civilians)
- [ ] Order BLE panic button — any $15 Bluetooth button works with Meshtastic
- [ ] Flash Meshtastic firmware to all 5 nodes when they arrive
- [ ] Set up DRONE-CMD (ch1) and ALERT (ch2) channels with AES-256 keys
- [ ] Test mesh range on property — verify 5-10km in open terrain

*Full suite: [[README]] | Operating procedures: [[wartime-weapon-protocol]]*
