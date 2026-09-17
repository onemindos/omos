---
domain: 01
name: Sovereign Defense — Tactical Drone & Robot Operations Suite
status: active
owner: zeus
tags: [sovereign-defense, drones, tactical, security, operations, robotics]
created: 2026-05-13
updated: 2026-05-13
---

# Sovereign Defense

**Tactical drone and robot operating system for OMOS. Covers everything from peacetime surveillance to WARFIGHTER autonomous deployment.**

This folder is the operational doctrine layer for Zeus's physical defense stack — drones, robots, sensors, mesh comms, edge AI, wearable HUD, and the protocols that bind them together across all 5 threat levels.

---

## Documents in This Suite

| File | Purpose | Status |
| --- | --- | --- |
| [[wartime-weapon-protocol]] | Full operating procedures — 5 threat levels, command chains, hardware loadouts, voice commands, PX4 missions, training drills | ✅ Complete |
| [[field-reference-card]] | Single-page laminated quick reference — memorize this, carry this | ✅ Complete |
| [[hardware-bom-acquisition]] | Phased hardware purchase plan — what to buy, in what order, from where, at what cost | ✅ Complete |
| [[legal-framework]] | US drone law by threat level — what's legal, where the line is, liability exposure | ✅ Complete |
| [[px4-mission-guide]] | How to build, upload, and test the 8 pre-loaded PX4 missions in QGroundControl | ✅ Complete |
| [[counter-drone-protocol]] | Detecting and responding to hostile drones targeting Zeus | ✅ Complete |
| [[swarm-deconfliction]] | Multi-drone collision avoidance — preventing friendly fire in the air | ✅ Complete |

---

## Relationship to Other Documents

```
RESEARCH (what we studied):
  04 Research/Robotics/drone-stack-architecture-openclaw-dimos.md
  04 Research/Robotics/drone-combat-readiness-assessment.md
  04 Research/Robotics/openmind-om1-deep-analysis.md
  04 Research/Robotics/om1-vs-dimos-vs-openclaw-comparison.md
  04 Research/Vision/OMOS-EagleEye.md
  04 Research/Robotics/Mesh-Comms/meshtastic-mesh-integration.md

OPERATIONS (what we do — THIS FOLDER):
  sovereign-defense/ ← YOU ARE HERE
    wartime-weapon-protocol.md  ← master doctrine
    field-reference-card.md     ← carry this
    hardware-bom-acquisition.md ← buy this
    legal-framework.md          ← know this
    px4-mission-guide.md        ← build this
    counter-drone-protocol.md   ← watch for this
    swarm-deconfliction.md      ← solve this
```

---

## 5-Minute Summary

**What this system does:**
Zeus's drones and robots operate intelligently at peacetime and autonomously at wartime — without internet, without a phone, without typing.

**The 5 modes:**

| Level | Name | Trigger | Command Method |
| --- | --- | --- | --- |
| 5 | SOVEREIGN 🟢 | Always | Telegram → OpenClaw → DimOS |
| 4 | VIGILANT 🔵 | Anomaly detected | Telegram + auto-dispatch |
| 3 | GUARDED 🟡 | Comms degraded | Physical RC + mesh radio |
| 2 | TACTICAL 🟠 | Internet down | Voice + mesh + edge AI |
| 1 | WARFIGHTER 🔴 | Under fire | Panic button → full autonomy |

**The key insight:**
OpenClaw runs the fleet at LEVEL 5-4. When internet dies, Meshtastic mesh takes over. When mesh dies, PX4 autonomous missions run with zero connectivity. The system is designed to degrade gracefully down to the flight controller itself — which needs nothing external to protect Zeus.

**Start here:**
Read [[field-reference-card]] first. That's what you carry and memorize. Everything else is reference.
