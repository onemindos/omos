---
title: "ATAK — Android Team Awareness Kit"
type: reference
domain: SO
status: active
created: 2026-05-27
updated: 2026-05-27
tags: [military, situational-awareness, C2, mapping, real-time, ATAK, TAK]
---

# ATAK — Android Team Awareness Kit

## Definition
Military-grade situational awareness and command & control platform. Developed by the Air Force Research Laboratory (AFRL). Provides real-time common operating picture (COP) across distributed teams in the field.

## Platform Variants
| Version | Platform | Audience |
|---------|----------|----------|
| **ATAK** | Android | Military / government (CAC-gated) |
| **iTAK** | iOS | Civilian / first responders |
| **WinTAK** | Windows | Fixed command posts |
| **WebTAK** | Browser | TOC / rear command display |
| **CIVTAK** | Android | Public safety / civilian ops |

## Core Capabilities
- Real-time blue force tracking (friendly positions)
- Shared common operating picture — everyone sees the same map
- Encrypted comms over mesh (MANET) or internet (TAK Server)
- Data layers: terrain, threats, assets, waypoints, geo-fences
- Cursor on Target (CoT) protocol — standard message format for position/status reporting
- Plugin ecosystem: HAMMER, SNAP, ISR feeds, sensor integration

## TAK Ecosystem
- **TAK Server** — self-hosted backend, federates multiple TAK clients
- **FreeTAKServer** — open-source TAK server implementation
- **Cursor on Target (CoT)** — XML-based message standard, the data backbone
- **ATAK plugins** — extend capability: drones, radios, sensors, biometrics

## Relevant Links
- TAK Product Center: https://tak.gov
- FreeTAKServer: https://github.com/FreeTAKTeam/FreeTAKServer
- ATAK civilian download: https://play.google.com/store/apps/details?id=com.atakmap.app.civ

## Notes
<!-- Add TAK server setup notes, plugin list, CoT schema references here -->
