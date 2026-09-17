---
title: "Blender"
type: reference
status: active
tags: [platform, creative, 3d, modeling, animation, vfx, open-source]
---
# Blender

> Open-source 3D creation suite. Modeling, animation, rendering, VFX, compositing, video editing.

## Key Features
| Feature | Notes |
|---------|-------|
| 3D Modeling | Sculpt, poly, NURBS |
| Animation | Keyframe, rigging, motion capture |
| Rendering | Cycles (raytracing), EEVEE (realtime) |
| VFX | Compositing, motion tracking |
| Video editing | Built-in NLE |
| Scripting | Full Python API |
| Geometry Nodes | Procedural modeling |
| Grease Pencil | 2D animation in 3D space |

## Python Scripting
```python
import bpy

# Create a cube
bpy.ops.mesh.primitive_cube_add(location=(0, 0, 0))

# Render
bpy.context.scene.render.filepath = '/tmp/render.png'
bpy.ops.render.render(write_still=True)
```

## Integration
- Python API for automated rendering
- CLI rendering for batch jobs
- AI-assisted (Stable Diffusion texture generation)
