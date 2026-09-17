---
title: "FreeMoCap"
type: reference
domain: SO
status: active
created: 2026-05-26
updated: 2026-05-26
tags: [motion-capture, computer-vision, blender, unity, unreal, 3d, open-source, research]
---

# FreeMoCap — Free Motion Capture

> Open-source markerless motion capture system. No suits, no markers — just cameras.

## Source
- GitHub: https://github.com/freemocap/freemocap

## Overview
FreeMoCap is a free, open-source markerless motion capture tool that uses computer vision (MediaPipe/OpenCV under the hood) to track full-body skeletal movement from standard cameras. No expensive hardware required.

## Engine Integrations
- **Blender** ✅ — exports rigs and animation data directly
- **Unity** ✅ — compatible output for game/sim pipelines
- **Unreal Engine** ✅ — compatible for real-time and cinematic pipelines

## OneMind Use Cases
- Drone/robot body tracking and movement modeling
- Avatar animation for OneMind virtual presence
- Physical performance analysis (HP/28 Body)
- Real-time motion data into NATS event stream
- Training data generation for AI movement models

## Related Domains
- `00-24 SO/07 Development` — integration projects
- `25-49 HP/28 Body (Health & Wellness)` — biomechanics, movement analysis
- `00-24 SO/06 OpenClaw` — avatar/canvas applications
