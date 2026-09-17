---
id: app-touchdesigner
title: TouchDesigner
category: Applications
status: documented
created: 2026-06-17
updated: 2026-06-17
tags: [real-time-graphics, node-based, multimedia, interactive, vj]
---

# TouchDesigner

## Application Overview

**TouchDesigner** is a professional node-based real-time graphics and multimedia creation platform developed by Derivative. It's used for creating interactive installations, live VJ performances, data visualizations, and multimedia experiences. The node graph-based workflow enables rapid prototyping and production-ready output without traditional coding (though Python integration is available).

| Property | Value |
|----------|-------|
| **Developer** | Derivative Inc. (Toronto, Canada) |
| **Bundle Identifier** | `ca.derivative.TouchDesigner` |
| **Version** | 2025.32820 |
| **Release Date** | Active (weekly builds, quarterly stable releases) |
| **Architecture** | arm64 (Apple Silicon native) |
| **Install Location** | `/Applications/TouchDesigner.app` |
| **App Size** | 1.7 GB |
| **License** | Freemium (Educational: free; Commercial: subscription-based) |
| **Website** | https://derivative.ca |
| **Documentation** | https://docs.derivative.ca |

## Copyright & Legal

```
Derivative Inc.
TouchDesigner is a registered trademark
Licensed under proprietary freemium model
```

## Key Features

### Node-Based Workflow
- **Operators (OPs)**: 200+ node types for different operations
  - **CHOPs**: Channel operators (animation, math, audio)
  - **TOPs**: Texture operators (video, image processing)
  - **SOPs**: Surface operators (3D geometry, mesh)
  - **POPs**: Particle operators (particle system nodes)
  - **DATs**: Data operators (text, tables, scripting)
  - **COMPs**: Component operators (groups, reusable modules)
- **Non-Destructive Workflow**: Edit nodes without re-rendering
- **Real-Time Feedback**: Instant visual preview of node changes
- **Custom TOE Files**: Save and reuse node graphs as components

### 3D Graphics & Rendering
- **GPU Rendering**: NVIDIA CUDA, AMD HIP, Apple Metal acceleration
- **Physically-Based Materials**: PBR shader support
- **Lighting System**: Multiple light types, global illumination simulation
- **Post-Effects**: Bloom, DOF, color grading, distortion
- **Instancing**: Render millions of objects efficiently
- **VR/360° Support**: Stereoscopic rendering, omni-directional output

### Real-Time Performance
- **GLSL Shaders**: Write custom fragment and vertex shaders
- **GPU Compute**: General-purpose GPU computation (CUDA/HIP)
- **Performance Optimization**: LOD, culling, instancing
- **Multi-GPU**: Distribute rendering across graphics processors
- **Resolution Independence**: Scale from HD to 8K output

### Animation & Timeline
- **CHOP Animation**: Channel-based animation system
- **Timeline COMP**: Frame-accurate sequencing and playback
- **Expression Language**: Python-based parameter expressions
- **Procedural Generation**: Script-based content creation
- **Keyframe Animation**: Traditional keyframe-based animation support

### Real-Time Audio
- **Audio Input**: Record and analyze live audio
- **FFT Analysis**: Frequency spectrum visualization
- **Audio Reactive**: Trigger effects based on audio features
- **MIDI Support**: Hardware controller input
- **Spatial Audio**: 3D audio positioning and HRTF

### Interactivity & Control
- **MIDI Controllers**: Map hardware sliders/buttons to parameters
- **OSC (Open Sound Control)**: Network-based parameter control
- **DMX512**: Lighting protocol integration
- **Mouse/Touch Input**: Respond to user interaction
- **Syphon/NDI Output**: Stream content to other apps
- **UDP/TCP Networking**: Send/receive data from external systems

### Python Integration
- **TScript API**: Python scripting within TouchDesigner
- **Custom Operators**: Write Python to process data
- **External Libraries**: NumPy, SciPy, sklearn integration
- **Real-Time Scripts**: Execute Python code per frame
- **Async Callbacks**: Background task execution

## User Data Locations

| Path | Size | Purpose |
|------|------|---------|
| `~/Library/Application Support/Derivative/` | Variable | Global preferences and config |
| `~/Library/Application Support/Derivative/TouchDesigner099/` | GB | Version-specific data |
| `~/Library/Application Support/Derivative/TouchDesigner099/FontCache/` | MB | Pre-rendered font glyphs |
| `~/Library/Application Support/Derivative/TouchDesigner099/plugins/` | Variable | Custom operators and extensions |
| `~/Library/Preferences/ca.derivative.TouchDesigner.plist` | KB | App preferences |
| `~/Library/Caches/ca.derivative.TouchDesigner/` | GB | GPU cache, shader compilation |
| Project Files: User-defined locations | Variable | `.toe` files (XML-based project format) |

## File Format Support

| Format | Type | Purpose |
|--------|------|---------|
| `.toe` | TouchDesigner Project | Native (XML text, versionable) |
| `.tox` | TouchDesigner Component | Reusable node graph package |
| `.toe.backup` | Auto-save backup | Crash recovery |
| `.avi/.mov/.mp4` | Video | Import/playback |
| `.png/.jpg/.exr/.tiff` | Image | Texture input |
| `.exr` | OpenEXR | 32-bit float image sequences |
| `.abc` | Alembic | 3D geometry import |
| `.obj/.fbx/.usdz` | 3D Model | Mesh import |
| `.wav/.aiff` | Audio | Sound input |
| `.midi` | MIDI Data | Sequence import |

## Performance Characteristics

- **Memory**: 500 MB – 4+ GB depending on scene complexity
- **GPU**: Leverages Metal heavily on Apple Silicon (very efficient)
- **CPU**: Multi-threaded; 4+ cores beneficial for Python and audio processing
- **Disk**: Project files typically 1-100 MB; media storage independent
- **Network**: Low-latency requirement for real-time control (< 10 ms optimal)
- **Resolution**: Supports 4K (3840×2160) and beyond; 8K feasible with GPU memory

## Licensing Model

### Educational (Free)
- Full feature set
- Unlimited output resolution
- No commercial use
- Personal/educational projects only

### Commercial (Subscription)
- Perpetual licensing
- Commercial use rights
- Priority support
- Includes plugins and extensions

## Supported GPUs

| GPU | Framework | Performance |
|-----|-----------|-------------|
| Apple Silicon (M1/M2/M3/M4) | Metal | Excellent (native) |
| NVIDIA RTX Series | CUDA | Excellent |
| AMD RDNA 2/3 | HIP | Good |
| Intel Arc | Limited | Fair |

## Integration with One Mind System

**Use Case**: Real-time interactive systems, live VJ, data visualization, installation art
- **Pillar**: GE (Generational Empire) — if producing branded experiences or products
- **Pillar**: SO (Sovereign Operations) — if managing interactive infrastructure
- **Relevant Domains**: Creative technology, real-time systems, multimedia installations

## System Requirements

- **macOS**: 10.13+ (Intel) / 11.0+ (Apple Silicon)
- **RAM**: 8 GB minimum; 16+ GB recommended
- **GPU**: Metal support; 2 GB VRAM minimum (4+ GB recommended)
- **Disk**: 2-5 GB for installation; project storage independent
- **Network**: Optional; required for OSC/Syphon/NDI features

## Common Use Cases

- **Live VJ**: Real-time visual mixing for music performances
- **Interactive Installations**: Sensor-driven art and immersive experiences
- **Corporate Events**: Real-time data visualization and stage graphics
- **Live Broadcast**: Dynamic graphics and real-time compositing
- **Architectural Projection**: Mapping and interactive building projections
- **Museum Exhibits**: Immersive, interactive museum installations
- **Scientific Visualization**: Real-time data and simulation visualization
- **Game Prototyping**: Rapid prototyping of interactive experiences

## Workflow Example: Interactive Installation

1. **Design Concept**: Sketch interaction model
2. **Build Node Graph**: Connect TOPs, CHOPs, SOPs for content
3. **Create Interactivity**: Add OSC/MIDI input handlers
4. **Add Shaders**: Write GLSL for custom effects
5. **Optimize**: Profile and optimize for target hardware
6. **Test**: Run in perform mode on target display
7. **Deploy**: Export to kiosk or installation environment

## Community & Ecosystem

- **TD Forum**: Active community with 100K+ members
- **Pre-Built Components**: Marketplace for ready-made modules
- **Custom Operators**: Community-contributed operators
- **Educational Resources**: Extensive tutorials and documentation
- **Conferences**: Annual Derivative conference with developer talks

## Removal Instructions

To completely uninstall TouchDesigner:

```bash
rm -rf /Applications/TouchDesigner.app
rm -rf ~/Library/Application\ Support/Derivative
rm -rf ~/Library/Preferences/ca.derivative.TouchDesigner.plist
rm -rf ~/Library/Caches/ca.derivative.TouchDesigner*
rm -rf ~/Library/Saved\ Application\ State/ca.derivative.TouchDesigner.savedState
```

## Alternatives

- **Nuke** (Compositing-focused, node-based)
- **Houdini** (Procedural 3D, can do similar tasks)
- **Processing/p5.js** (Code-based creative coding)
- **Resolume** (VJ-focused alternative)
- **vvvv** (Graph-based creative coding platform)
- **Notch** (GPU-accelerated real-time graphics)

## Notes

- Industry-standard in live performance, installations, and interactive design
- Freemium model makes it accessible; commercial licensing affordable
- Node-based workflow has low barrier to entry; scalable to complex projects
- Python integration enables advanced procedural workflows
- Active development; new features regularly added
- Strong vendor support and community
- Can be extended with custom operators via C++
- Projects are version-controllable (`.toe` files are XML text)
- Competitive with proprietary tools at a fraction of cost

## Performance Tips

- Use SOP instancing for rendering large numbers of objects
- Implement LOD (Level of Detail) for complex scenes
- Leverage GPU compute for CPU-intensive operations
- Profile performance with built-in Performance Monitor
- Use Texture Memory Pool to manage GPU memory efficiently

---

*Last verified: June 17, 2026 (v2025.32820)*
