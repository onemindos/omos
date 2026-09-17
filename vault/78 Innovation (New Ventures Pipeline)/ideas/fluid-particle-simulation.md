---
title: "Fluid Particle Simulation"
type: reference
domain: GE
status: active
created: 2026-05-27
updated: 2026-05-27
tags: [creative-tech, realtime, simulation, particles, generative, touchdesigner, visual, vfx]
---

# Fluid Particle Simulation

## What It Is
Fluid particle simulation is a computational technique that models fluid behavior (liquids, gases, smoke, fire) using discrete particles rather than a continuous mesh. Each particle carries properties like velocity, density, and pressure, and interacts with neighboring particles according to physics equations. Used across film VFX, games, interactive installations, and real-time generative art.

## Main Simulation Methods

| Method | Description | Speed | Best For |
|--------|-------------|-------|----------|
| SPH (Smoothed-Particle Hydrodynamics) | Mesh-free Lagrangian — fluid properties calculated from weighted neighbors | Slow (7-20+ substeps/frame) | Slow constrained fluids; physically accurate |
| FLIP (Fluid Implicit Particle) | Hybrid particle + volume — stores data in particles, pressure projection on temp volumes | Fast (few substeps/frame) | Fast-moving fluids (rivers, waterfalls), new particle injection |
| APIC (Affine Particle-In-Cell) | Variant of PIC — better vortex/rotation preservation | Medium | More energetic motion, less splashy than FLIP |
| Shallow Water Equations | 2D simplified method using render target textures — ground, velocity, height maps | Very fast | Real-time water surfaces in games/installations |
| GPU Particle Systems | Compute shader-based — each GPU texture pixel is a particle | Very fast | Massive particle counts, stylized/art-directed effects |
| Flow Fields (Vector Fields) | Particles follow a pre-computed or dynamic vector field | Very fast | Stylized smoke, atmospheric, crowd effects |

**SPH stability note:** Can destabilize if particles are placed too close together. FLIP handles overlapping particles without instability.

## Tools That Support It

| Tool | Method | Context |
|------|--------|---------|
| Houdini (SideFX) | FLIP, SPH, APIC | Film/TV VFX industry standard |
| TouchDesigner (Derivative) | GPU particles, Flex physics, ParticleSOP | Real-time installations, live AV |
| Unity VFX Graph | GPU particles, compute shaders | Games, interactive experiences |
| Unreal Engine Niagara | CPU + GPU particles, pre-built fluid templates | Games, virtual production |
| Blender (Mantaflow) | FLIP, APIC, SPH | Animation, VFX, free/open-source |
| Three.js (WebGPU) | GPU particles, MPM-style grid | Web-based interactive experiences |

## In TouchDesigner Specifically
- **ParticleSOP** — native particle creation and manipulation
- **Particle POP** — Particle Operator system
- **Flex physics simulator** — GPU-accelerated, particle-based, optimized for fluid behavior
- **GPU particle pattern**: each pixel in a TOP (texture operator) = one particle; instanced via Geometry COMP
- Python scripting for custom force/behavior logic
- Real-time capable at performance framerates (60fps+)

## Real-World Use Cases
- Film/TV: water, smoke, fire, destruction sequences (Houdini dominant)
- Games: rain, portals, magic effects, dynamic weather (Unity/Unreal)
- Interactive installations: audience-reactive particle clouds, live AV performances (TouchDesigner)
- Scientific visualization: CFD, blood flow, aerodynamics
- Web: interactive hero sections, data visualization (Three.js/WebGPU)

## Related Concepts
- [[touch-designer]] — primary environment for real-time fluid particle work in OneMind stack
- [[real-time-story-creation]] — production context (audio-reactive particles, face overlays, scene transitions)

## Notes
<!-- Add: specific TD networks, Flex setup notes, performance benchmarks on M2 Max GPU -->
