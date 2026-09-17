---
title: "iNav"
type: reference
status: planned
tags: [platform, hardware, drones, autopilot, navigation, open-source]
---
# iNav

> Navigation-focused firmware. Simpler than ArduPilot for GPS waypoint missions on smaller/cheaper hardware.

## Use Case
- Simple GPS waypoint missions on budget hardware
- Fixed-wing long-range autonomous flights
- When ArduPilot is overkill for the hardware
- Good for lightweight mapping/survey drones

## Key Features
| Feature | Notes |
|---------|-------|
| Navigation | Full waypoint missions, RTH, position hold |
| Hardware | Runs on cheap F4/F7 boards (same as Betaflight) |
| Fixed-wing | Excellent fixed-wing support |
| Mixer | Custom motor/servo mixing |
| OSD | Navigation data on FPV display |
| Programming | Logic conditions (if battery < X, then RTH) |

## vs ArduPilot
| | iNav | ArduPilot |
|---|---|---|
| Hardware cost | $30 F4 board | $200+ Pixhawk |
| Setup | Simpler | Complex |
| Features | Basic nav | Full autonomy stack |
| Companion computer | No | Yes (Pi, Jetson) |
| Scripting | Logic conditions | Lua + MAVLink |

## Configurator
- iNav Configurator (Electron app)
- Mission Control (mobile planning app)

## Resources
- https://github.com/iNavFlight/inav
- https://github.com/iNavFlight/inav/wiki
