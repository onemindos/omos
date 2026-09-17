---
id: app-blender
title: Blender
category: Applications
status: documented
created: 2026-06-17
updated: 2026-06-17
tags: [3d-modeling, animation, rendering, cad, visual-effects]
---

# Blender

## Application Overview

**Blender** is a free and open-source 3D computer graphics software for creating animated films, visual effects, art, 3D models, motion graphics, interactive 3D applications, virtual reality, and video games.

| Property | Value |
|----------|-------|
| **Developer** | Blender Foundation & Community |
| **Bundle Identifier** | `org.blenderfoundation.blender` |
| **Version** | 5.1.1 |
| **Release Date** | Active (released quarterly) |
| **Architecture** | arm64 (Apple Silicon native) |
| **Install Location** | `/Applications/Blender.app` |
| **App Size** | 879 MB |
| **License** | GPL v3 (Open Source) |
| **Website** | https://www.blender.org |
| **Source Code** | https://github.com/blender/blender |

## Copyright & Legal

```
Blender Foundation & Contributors
Licensed under GNU General Public License v3
```

## Key Features

### Modeling
- **Mesh Modeling**: Vertex, edge, and face-level editing
- **Sculpting**: High-poly sculpting with dynamic topology
- **Modifiers**: Non-destructive modeling with 40+ procedural modifiers
- **Geometry Nodes**: Node-based procedural geometry generation
- **CAD Tools**: Boolean operations, exact modeling

### Rendering
- **Cycles**: GPU/CPU physically-based path tracer
- **Eevee**: Real-time viewport rendering engine
- **OptiX/HIP**: NVIDIA/AMD GPU acceleration
- **Volumetrics**: Smoke, fire, and fluid simulations

### Animation & Simulation
- **Rigging**: Skeletal animation with IK/FK
- **Motion Capture**: BVH file support
- **Particle Systems**: Procedural particles and hair
- **Fluid Dynamics**: Liquid and smoke simulation
- **Cloth Simulation**: Dynamic cloth and soft bodies
- **Rigid Body Physics**: Collision and gravity simulation

### Visual Effects
- **Compositor**: Node-based post-processing
- **Video Sequence Editor**: Multi-track video editing
- **Motion Tracking**: Camera tracking and stabilization
- **Grease Pencil**: 2D drawing within 3D space

### Texture & Shading
- **Shader Editor**: Node-based material system
- **Texture Painting**: Paint directly on 3D models
- **UV Mapping**: Automatic and manual UV unwrapping
- **Procedural Textures**: Built-in texture generation

## User Data Locations

| Path | Size | Purpose |
|------|------|---------|
| `~/Library/Application Support/Blender/` | Variable | Startup scripts, presets, add-ons |
| `~/Library/Application Support/Blender/[version]/` | Variable | Version-specific config (cache, preferences) |
| `~/Library/Application Support/Blender/[version]/cache/` | MB-GB | Rendered tile cache, UDIM textures |
| `~/Library/Application Support/Blender/[version]/startup.blend` | KB | Default startup file |
| `~/Library/Application Support/Blender/[version]/config/startup/` | KB | Startup scenes and preferences |
| `~/Library/Preferences/org.blenderfoundation.blender.plist` | Variable | Application preferences |
| `~/Library/Caches/org.blenderfoundation.blender/` | GB | Temporary cache files |
| `~/Library/Containers/org.blenderfoundation.blender.thumbnailer/` | Variable | Thumbnail renderer sandbox |

## Supported File Formats

| Format | Type | I/O |
|--------|------|-----|
| `.blend` | Blender Native | Read/Write |
| `.obj` | Wavefront OBJ | Read/Write |
| `.glb/.gltf` | Khronos glTF | Read/Write |
| `.fbx` | Autodesk FBX | Read/Write |
| `.usdz/.usd` | Pixar USD | Read/Write |
| `.stl` | Stereolithography | Read/Write |
| `.ply` | Point cloud | Read/Write |
| `.dae` | COLLADA | Read/Write |
| `.abc` | Alembic | Read/Write |
| `.exr` | OpenEXR (32-bit float) | Write |
| `.png/.jpg/.tiff` | Standard Image | Read/Write |

## Performance Characteristics

- **Memory**: 2-8 GB for typical projects; 16+ GB for heavy VFX
- **CPU**: Leverages all cores for rendering and simulation
- **GPU**: CUDA/OptiX (NVIDIA), HIP (AMD), Metal (Apple)
- **Disk**: Project files typically 50 MB – 1 GB; cache can grow large
- **Network**: Standalone; no network requirements (optional)

## Render Engine Capabilities

### Cycles (Physically-Based Path Tracer)
- Unbiased/biased rendering options
- GPU and CPU rendering
- Denoising (AI-powered via OptiX)
- Adaptive sampling

### Eevee (Real-Time Viewport)
- Screen-space reflections
- Global illumination
- Bloom and lens effects
- Subsurface scattering

## Add-On Ecosystem

Blender supports 300+ official and community add-ons:
- **Hard-Surface Modeling**: Boolean operations, retopo, modeling tools
- **Procedural Generation**: Parametric furniture, trees, cities
- **Rendering**: Advanced materials, light linking
- **Animation**: Motion capture, constraints, rigging helpers
- **VFX**: Fluid, smoke, volumetric effects
- **Import/Export**: Format converters, bridge to engines

## Integration with One Mind System

**Use Case**: 3D design, visualization, and content creation
- **Pillar**: GE (Generational Empire) — product visualization, brand assets
- **Pillar**: SO (Sovereign Operations) — if managing CAD/design infrastructure
- **Relevant Domains**: Design, content, visualization, prototyping

## System Requirements

- **macOS**: 10.13+ (Intel) / 11.0+ (Apple Silicon)
- **RAM**: 4 GB minimum; 16+ GB recommended
- **Disk**: 2-5 GB available for cache and projects
- **GPU**: Optional but strongly recommended for rendering

## Removal Instructions

To completely uninstall Blender:

```bash
rm -rf /Applications/Blender.app
rm -rf ~/Library/Application\ Support/Blender
rm -rf ~/Library/Caches/org.blenderfoundation.blender*
rm -rf ~/Library/Preferences/org.blenderfoundation.blender.plist
rm -rf ~/Library/Containers/org.blenderfoundation.blender*
rm -rf ~/Library/Application\ Scripts/org.blenderfoundation.blender*
```

## Alternatives

- **Autodesk Maya** (Commercial, VFX-focused)
- **3ds Max** (Commercial, game/VFX oriented)
- **Houdini** (Procedural, VFX powerhouse)
- **Cinema 4D** (Commercial, motion graphics)
- **Lightwave 3D** (Legacy)

## Notes

- 100% open-source; no vendor lock-in
- Highly extensible via Python scripting
- Active development; major releases every 3-4 months
- Mature enough for professional studio production
- Large community with tutorials, documentation, and support
- Industry adoption growing (used by companies like Tangent Animation, King's College London)

---

*Last verified: June 17, 2026 (v5.1.1)*
