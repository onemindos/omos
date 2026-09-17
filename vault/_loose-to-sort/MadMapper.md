---
id: app-madmapper
title: MadMapper
category: Applications
status: documented
created: 2026-06-17
updated: 2026-06-17
tags: [projection-mapping, visual-effects, live-performance, vj]
---

# MadMapper

## Application Overview

**MadMapper** (officially "MadMapper Demo 6.1.2") is a professional projection mapping and visual effects software designed for real-time 3D projection content creation, live VJ performances, and interactive installations. It enables precise mapping of video content onto irregular surfaces, architectural elements, and 3D geometries.

| Property | Value |
|----------|-------|
| **Developer** | GarageCube (France) |
| **Bundle Identifier** | `com.garagecube.MadMapperDemo` |
| **Product Name** | MadMapper |
| **Version** | 6.1.2 (Demo Version) |
| **Release Date** | Active (regular updates) |
| **Architecture** | arm64 (Apple Silicon native) |
| **Install Location** | `/Applications/MadMapperDemo 6.1.2.app` |
| **App Size** | 400 MB |
| **License** | Proprietary (Freemium model; demo is feature-limited) |
| **Website** | https://www.madmapper.com |

## Copyright & Legal

```
GarageCube SAS
MadMapper is a registered trademark
```

## Key Features

### Projection Mapping
- **Surface Calibration**: Map content onto non-planar surfaces using 3D point cloud
- **Quad Mapping**: Simple quad-based surface mapping
- **3D Object Mapping**: Project onto imported 3D models (OBJ, DAE)
- **Keystone Correction**: Manual corner/edge pin adjustment
- **Mesh Deformation**: Fine-grained vertex control for precise alignment

### Real-Time Rendering
- **GPU Acceleration**: Metal (Apple Silicon optimized)
- **Shader Support**: GLSL shader programming for effects
- **Blend Modes**: 20+ layer blending modes
- **Particle Systems**: GPU-accelerated particle generators
- **Image Processing**: Color correction, blur, distortion effects

### Content Management
- **Multi-Layer Timeline**: Timeline-based sequencing
- **Media Support**: Video (H.264, ProRes), images, live camera feeds
- **Syphon/NDI**: Real-time video stream input/output
- **MIDI/OSC Control**: External device triggering and parameter control
- **DMX512**: Lighting protocol integration

### Interactive Features
- **Interactive Mapping**: Respond to user input (mouse, touch, sensors)
- **OSC Messaging**: Open Sound Control for network triggers
- **MIDI Controller Support**: Assign hardware controllers to effects
- **Sensor Integration**: Kinect, depth cameras for interactive installations

### Performance & Optimization
- **Multi-GPU Support**: Distribute rendering across graphics processors
- **Real-Time Preview**: Instant visual feedback
- **Output Scaling**: Support for 4K+ projection systems
- **Framerate Control**: Lock to 24/30/60 fps for synchronization

## User Data Locations

| Path | Size | Purpose |
|------|------|---------|
| `~/Library/Application Support/MadMapper/` | Variable | Projects, presets, custom effects |
| `~/Library/Application Support/MadMapper/MadMapper/` | MB-GB | User library, media cache |
| `~/Library/Preferences/com.madmapper.MadMapper.plist` | KB | General preferences |
| `~/Library/Preferences/com.garagecube.MadMapperDemo.plist` | KB | Demo-specific preferences |
| `~/Library/Logs/MadMapper/` | MB | Debug and performance logs |
| `~/Library/Caches/com.garagecube.MadMapperDemo/` | Variable | Texture cache, media thumbnails |

## File Format Support

| Format | Type | Support |
|--------|------|---------|
| `.madmap` | MadMapper Project | Native (Read/Write) |
| `.mp4/.mov` | Video | Import |
| `.png/.jpg/.tiff` | Image | Import |
| `.syphon` | Real-time video stream | Input/Output |
| `.ndi` | Network video stream | Input/Output |
| `.obj/.dae` | 3D Model | Import for mapping |
| `.glsl` | GLSL Shader | Import for effects |

## Connectivity Options

### Video I/O
- **Syphon**: macOS-native real-time video sharing with other apps
- **NDI (Network Device Interface)**: Network video streaming
- **Live Camera**: USB webcams and built-in cameras
- **Screen Capture**: Capture desktop or app windows

### Control Interfaces
- **MIDI**: Hardware music controllers, sliders, drum pads
- **OSC (Open Sound Control)**: Network protocol for parameter control
- **DMX512**: Stage lighting protocol via USB adapter
- **Network Sockets**: Custom TCP/UDP messaging

### Output
- **Projector Output**: Full-screen projection to connected display
- **Output Mapping**: Multiple projector support with edge blending
- **Recording**: Capture output to video file

## Performance Characteristics

- **Memory**: 400 MB – 2 GB depending on resolution and layer count
- **GPU**: Metal framework; leverages Apple Silicon GPU efficiently
- **CPU**: Multi-threaded; can utilize 4+ cores
- **Network**: Low latency for OSC/NDI streaming (1-10 ms)
- **Disk**: Project files typically 1-50 MB; media storage scales with content library

## Demo vs. Full Version

**Demo Limitations** (6.1.2):
- Output limited to 2048×2048 pixels
- Some advanced features restricted
- Watermark on output (may apply)
- Full licensing required for commercial use

**Full Version Advantages**:
- Unlimited output resolution (8K+)
- Professional plugins and effects
- Priority support
- Perpetual license or subscription model

## Supported Use Cases

- **Festival/Concert Production**: Real-time VJ mapping for stage visuals
- **Architectural Projection**: Building facades, sculptures, monuments
- **Installation Art**: Interactive and generative installations
- **Live Events**: Conference presentations, ceremonies, performances
- **Commercial Advertising**: Projection-based brand activations
- **Theater/Dance**: Set design and lighting integration

## Integration with One Mind System

**Use Case**: Real-time visual content creation and projection systems
- **Pillar**: GE (Generational Empire) — if producing branded/commercial content
- **Pillar**: SO (Sovereign Operations) — if managing live event infrastructure
- **Relevant Domains**: Media production, events, interactive installations

## System Requirements

- **macOS**: 10.13+ (Intel) / 11.0+ (Apple Silicon)
- **RAM**: 4 GB minimum; 8+ GB recommended
- **GPU**: Metal support (integrated or discrete); 1 GB VRAM minimum
- **Display**: One or more connected displays/projectors

## Common Workflows

1. **Create Project**: Start new MadMapper project
2. **Import Media**: Add video/image content to timeline
3. **Calibrate Surface**: Define mapping surface (quad, mesh, or 3D model)
4. **Apply Effects**: Layer shaders and color corrections
5. **Setup Control**: Configure MIDI/OSC triggers
6. **Live Performance**: Execute mapping with real-time parameter adjustments

## Removal Instructions

To completely uninstall MadMapper:

```bash
rm -rf "/Applications/MadMapperDemo 6.1.2.app"
rm -rf ~/Library/Application\ Support/MadMapper
rm -rf ~/Library/Preferences/com.madmapper.MadMapper.plist
rm -rf ~/Library/Preferences/com.garagecube.MadMapperDemo.plist
rm -rf ~/Library/Caches/com.garagecube.MadMapperDemo*
rm -rf ~/Library/Logs/MadMapper
```

## Alternatives

- **Resolume Arena** (Commercial, VJ-focused)
- **Derivative TouchDesigner** (Node-based, real-time design platform)
- **Processing** (Open-source, coding-based visualization)
- **Quartz Composer** (Discontinued; node-based)
- **Mapping Pro by Lightconverse** (Focused projector control)

## Notes

- Industry-standard in projection mapping and live visual effects
- Demo version suitable for learning; full version required for professional work
- Active community with tutorials, presets, and plugin ecosystem
- Supports advanced workflows: multi-projector edge blending, 3D calibration, AI content generation
- Regular updates and feature additions
- Strong integration with live performance tools (VJing software, audio tools, lighting control)

## Technical Architecture

- **Rendering**: OpenGL → Metal abstraction layer
- **Timeline Engine**: Frame-accurate playback with synchronization
- **Shader System**: GLSL 3.2 support with custom effects
- **Export**: ProRes, H.264 encoding
- **Scripting**: Lua scripting for custom effects and automation

---

*Last verified: June 17, 2026 (v6.1.2 Demo)*
