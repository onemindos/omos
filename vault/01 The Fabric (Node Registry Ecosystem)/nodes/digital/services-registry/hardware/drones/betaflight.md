---
title: "Betaflight"
type: reference
status: planned
tags: [platform, hardware, drones, fpv, racing, firmware]
---
# Betaflight

> FPV racing/freestyle drone firmware. Manual pilot control — not for autonomous missions.

## Use Case
- FPV racing builds
- Freestyle acrobatics
- Manual recon flights (pilot in the loop)
- NOT for autonomous/waypoint missions (use ArduPilot/PX4 for that)

## Key Features
| Feature | Notes |
|---------|-------|
| PID tuning | Precision flight controller tuning |
| OSD | On-screen display overlay in goggles |
| Rates | Customizable stick response curves |
| GPS Rescue | Emergency RTH (not full nav) |
| Blackbox | Flight data logging for tuning |

## Configurator
```bash
# Betaflight Configurator (Electron app)
# Download from GitHub releases or Chrome Web Store
```

## Hardware
- Flight controllers: F4/F7/H7 processors
- ESCs: BLHeli_32, AM32
- Receivers: ELRS, Crossfire, FrSky
- VTX: DJI O3, HDZero, Walksnail, analog

## Resources
- https://betaflight.com
- https://github.com/betaflight/betaflight
