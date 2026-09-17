---
title: "TouchDesigner"
type: reference
domain: GE
status: active
created: 2026-05-27
updated: 2026-05-27
tags: [creative-tech, realtime, visual-programming, generative, media, tool, derivative]
---

# TouchDesigner

## What It Is
TouchDesigner is a visual development platform by Derivative used to build real-time interactive 3D and 2D applications. It is a procedural, node-based visual programming language — no compilation required. The system is always-alive: it continuously generates output even while you are modifying it. Used by artists, designers, and developers for interactive media, installations, live performances, and data-driven visuals.

Official site: https://derivative.ca
Docs: https://docs.derivative.ca

## How It Works
Everything in TouchDesigner is an Operator (OP). Operators are wired together into networks:

| Operator Family | Role |
|-----------------|------|
| TOP (Texture Operator) | 2D image and video processing |
| CHOP (Channel Operator) | Audio, motion, data streams |
| SOP (Surface Operator) | 3D geometry |
| DAT (Data Operator) | Text, tables, scripts |
| COMP (Component) | Containers, UI, 3D scene objects |
| MAT (Material Operator) | Shaders and surface rendering |

## Key Features
- Real-time GPU-accelerated rendering
- Live editing — modify while running, no restart
- Multi-GPU support with per-GPU binding
- Protocol support: Art-Net, Dante, NDI, OpenCV, OpenGL, OSC, MIDI, serial, WebSocket
- Python scripting inside the network
- Projection mapping integration (Vioso, Scalable Display)
- Performer tracking (BlackTrax integration)
- Instancing via Geometry COMP for massive particle counts
- GPU particle simulation: each pixel in a TOP represents a particle
- Flex physics simulator (GPU-accelerated, particle-based, excels with fluids)

## Use Cases
- Live performances and VJing
- Projection mapping
- Interactive installations and immersive experiences
- Generative art and data visualization
- Museum and theater productions
- Motion-controlled visuals and lighting
- Audiovisual performances

## Licensing & Pricing

| License | Price | Use |
|---------|-------|-----|
| Non-Commercial | Free | Personal, educational, evaluation |
| Educational | $300 USD | Students/accredited institutions only |
| Commercial | $600 USD | Any compensated work |
| Pro | $2,200 USD | Advanced features + direct support |

Annual update costs: $150 (Educational), $300 (Commercial), $400 (Pro).
Floating Cloud licenses available.

## Platform Support
- **Windows**: Primary platform. Minimum 4GB GPU RAM, 8GB+ recommended.
- **macOS**: Supported. Minimum macOS 12 Monterey (13 Ventura recommended). Apple Silicon highly recommended. Intel Mac requires discrete AMD GPU.

## Related Concepts
- [[fluid-particle-simulation]] — runs natively inside TD via Flex + ParticleSOP
- [[autodesk-flow]] — production pipeline for managing asset versions TD renders
- [[real-time-story-creation]] — primary OneMind use case

## Notes
<!-- Add: specific project networks, license key location in 1Password, hardware node assignment -->
