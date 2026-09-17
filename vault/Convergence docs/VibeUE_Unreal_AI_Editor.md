# VibeUE — AI Editor Toolset for Unreal Engine 5.8 + MCP Integration

**Official Links:**
- **Website:** https://www.vibeue.com/
- **GitHub:** https://github.com/vibeue
- **Discord Community:** https://discord.gg/vibeue
- **Documentation:** https://www.vibeue.com/docs
- **Fab Marketplace:** https://www.fab.com/listings/vibeue
- **Unreal Engine MCP Documentation:** https://dev.epicgames.com/documentation/unreal-engine/unreal-mcp-in-unreal-editor
- **Epic Developer Community:** https://dev.epicgames.com/

**Type:** AI Editor Plugin for Unreal Engine 5.8+

**Status:** Production, open-source, actively maintained

**Latest Version:** 5.8 branch (June 2026)

**Architecture:** Native UE5.8+ MCP plugin + Python-first toolset

---

## What Is VibeUE?

**VibeUE** is a **native AI editor toolset** that extends Unreal Engine 5.8's experimental **Model Context Protocol (MCP)** plugin, enabling AI agents (Claude, OpenAI Codex, other LLMs) to **directly control the Unreal editor** through natural language prompts.

**Core Innovation:** Bridges the gap between **AI reasoning systems** and **3D game engine automation** — allowing agents to:
- Place props in the world
- Edit Blueprint logic
- Configure lighting
- Manage terrain & world geometry
- Control audio, animation, VFX
- Profile performance
- And 18+ additional domains

**Philosophy:** "One plugin, native to the engine" — no external servers, no API overhead. Pure local Python-MCP integration inside the editor.

---

## Native Systems Integration

**VibeUE leverages three native UE 5.8+ systems:**

### 1. Model Context Protocol (MCP) — AI Integration

**What MCP Does:**
- Standardized protocol for AI agents to interact with tools/systems
- Defines schema for available functions
- Manages context, stateful operations
- Bidirectional: editor → agent, agent → editor

**Unreal's MCP Implementation (June 2026):**
- Native plugin in UE 5.8 core
- Experimental but production-ready
- Direct Python bindings to editor systems
- No additional servers required

**Supported AI Agents:**
- Claude (via Claude Code, Claude for VS Code)
- OpenAI Codex
- Other MCP-compatible LLMs

### 2. Python Scripting Layer

**VibeUE Python-First Architecture:**
- Full access to UE editor internals via Python API
- SoapySDR-like binding model (if you know the pattern)
- Direct function calls to engine systems
- Real-time execution in editor

**Python Access Points:**
- World editing (place/move/rotate actors)
- Blueprint compilation & modification
- Material editing & shader assignment
- Performance profiling hooks
- Undo/redo control
- Asset management

### 3. Performance Profiling & Insights

**Integrated UE Insights:**
- `frame_timing()` — query frame times, CPU/GPU costs
- Full capture lifecycle — start/stop profiling captures
- `analyse()` — automated bottleneck analysis
- Standalone build profiling (not just editor)
- Real-time performance feedback to agents

---

## Domain Coverage (19+ Toolsets)

VibeUE provides specialized toolsets across major Unreal Engine domains:

### 1. 🏔️ Terrain & World (6 toolsets)
- Landscape editing
- Actor placement & transformation
- World geometry (brushes, volumes)
- Navigation mesh generation
- Level streaming
- World partition management

### 2. 🎵 Audio (2 toolsets)
- Sound source placement
- Audio parameter control
- Soundscape configuration
- Spatial audio setup

### 3. 🎬 Animation Assets (4 toolsets)
- Animation blueprint setup
- Motion capture data import
- Skeletal mesh configuration
- Animation state machine editing

### 4. ✨ FX & Mesh Depth (3 toolsets)
- Particle effects (Niagara)
- VFX system configuration
- Mesh material assignment
- Complex geometry operations

### 5. 🎮 Enhanced Input (1 toolset)
- Input system configuration
- Control scheme setup
- Gamepad/keyboard binding
- Action/axis mapping

### 6. 🧠 Gameplay Logic & Data (2 toolsets)
- Blueprint editing & compilation
- Data table management
- Game mode configuration
- Gameplay ability system setup

### 7. ↩️ Editor Undo/Redo (1 toolset)
- Group multiple edits into single undo action
- Cancel-with-rollback operations
- History inspection
- Undo buffer reset

### 8. 📘 Higher-Order Blueprints (1 toolset)
- Blueprint macro libraries
- Blueprint function libraries
- Blueprint base classes
- Advanced Blueprint patterns

### 9. 🐍 Python-First Access (1 toolset)
- Direct Python API calls
- Custom scripting for specialized tasks
- Extensibility hooks
- Advanced automation

### 10. 🌐 Web Research (1 toolset)
- Fetch external data
- Asset library integration
- Reference material fetching
- Real-time web data binding

---

## How VibeUE + Unreal MCP Work Together

### Architecture Flow

```
User Natural Language Prompt
  ↓
AI Agent (Claude / OpenAI Codex)
  ↓
Model Context Protocol (MCP) Interface
  ↓ (schema + function calls)
Unreal Engine 5.8 MCP Plugin
  ↓ (native bridge)
VibeUE Python Toolsets (19+ domains)
  ↓ (execute in editor)
Real-time Editor Updates
  ↓
User sees result in Unreal Editor
```

### Example Workflow

**User prompt:**
"Place 10 trees randomly in the forest area, then add directional lighting that matches sunrise"

**AI Agent breakdown:**
1. Parse prompt into components (trees, placement, lighting)
2. Call VibeUE toolset: `terrain.place_actor_random(actor_type='tree', count=10, zone='forest')`
3. Call VibeUE toolset: `lighting.setup_directional(angle='sunrise', intensity=1.2)`
4. Query `performance.frame_timing()` to verify impact
5. Report back to user with results

**Result:** Forest layout + directional lighting configured in seconds (vs. manual hours)

---

## Key Features & Capabilities

### 1. Agentic Unreal Workflows

**Agents can:**
- Understand Unreal-specific domain language
- Sequence multi-step operations
- Reason about world state (where are things, what's visible?)
- Adapt based on performance feedback
- Learn from undo/redo history

**Example:** "Create a multiplayer arena with spawn points, weapon racks, and balanced lighting for competitive play"
- Agent breaks into: arena mesh, spawners, weapon placement, lighting config
- Sequences operations in proper order
- Verifies each step, adjusts if needed

### 2. World, UI, and Content Tools

**Full world editing:**
- Place/move/delete actors (with undo/redo)
- Configure BSP brushes
- Build navmesh
- Manage world partitions
- Landscape sculpting

**UI system:**
- Widget blueprint creation
- Layout automation
- Animation sequence setup
- Button/input binding

**Content:**
- Import models/audio/FX
- Assign materials
- Configure LODs
- Optimize memory usage

### 3. Learn & Understand

**AI-Assisted Learning:**
- Agents explain their reasoning ("I placed the light here because...")
- Educational tooltips from VibeUE
- Blueprint pattern recommendations
- Performance optimization suggestions
- Best-practice guidance

---

## Integration with Agent Skills

**VibeUE + Agent Skills = Complete Knowledge System**

### VibeUE (the live connection)
- Real-time editor control
- Direct Python API to engine
- Performance profiling hooks
- Undo/redo management

### Agent Skills (the knowledge)
- Open-source library of domain expertise
- Engine-source knowledge about:
  - Material system (how to set params, blend modes, etc.)
  - Blueprint patterns (state machines, gameplay systems)
  - Performance optimization (draw calls, LOD strategies, memory budgets)
  - Audio design (spatial audio, attenuation, mixing)
  - Animation workflows (retargeting, blending, state machines)

**Pairing:** Agent can use VibeUE to *execute*, and Agent Skills to *understand the why*

Example:
- Agent queries: "What's the best way to reduce draw calls?"
- Agent Skills respond: "Use dynamic LODs + HZB occlusion culling. Here are the parameters..."
- Agent executes via VibeUE: `rendering.enable_lod_dynamic() + rendering.enable_hzb_occlusion()`

---

## Getting Started with VibeUE

### Step 1: Configure Unreal's Native MCP

**Official Documentation:** https://dev.epicgames.com/documentation/unreal-engine/unreal-mcp-in-unreal-editor

**Process:**
1. Install Unreal Engine 5.8+
2. Enable MCP plugin (experimental)
3. Configure MCP server address
4. Set up authentication (NKey if using NATS)

### Step 2: Install VibeUE

**Installation Options:**

- **From Fab Marketplace:** https://www.fab.com/listings/vibeue
  - One-click install into project
  - Auto-configures MCP
  - Latest version always available

- **From GitHub:** https://github.com/vibeue
  - Clone repo into `Plugins/VibeUE/`
  - Build from source
  - Latest development version

### Step 3: Set Up the Agent Guide

**Documentation:** https://www.vibeue.com/docs

**Includes:**
- MCP schema reference
- Python API documentation
- 19+ domain toolset guides
- Example workflows
- Best practices

### Step 4: Launch & Prompt

**In Unreal Editor:**
1. Open editor
2. Open MCP command palette (Ctrl+Shift+M)
3. Type prompt (e.g., "Place 5 trees in the meadow")
4. Watch AI execute in real-time

---

## Unreal Engine 5.8 MCP Features (June 2026 Release)

**Released:** June 17, 2026

**Capabilities:**
- ✅ Direct editor automation via natural language
- ✅ Blueprint editing + compilation
- ✅ Asset placement & configuration
- ✅ Lighting setup (directional, point, spot, sky)
- ✅ Material assignment & parameter tuning
- ✅ Performance profiling integration
- ✅ Undo/redo support
- ✅ Python scripting layer
- ✅ Multi-agent coordination (multiple agents on same scene)
- ✅ Real-time feedback to agents

**LLM Compatibility:**
- Claude (Anthropic) — Claude Code, Claude for VS Code
- OpenAI Codex
- Other MCP-compatible models

**Status:** Experimental but production-ready (used in studios 2026+)

---

## Performance & Profiling

### frame_timing() — Real-time Performance Metrics

**Query current frame performance:**
```python
timing = unreal.PerformanceService.frame_timing()
# Returns: frame_time, GPU_time, CPU_time, draw_calls, triangles, memory
```

**Use case:** Agent can profile changes in real-time
- Place actor → measure frame impact
- Adjust LOD → measure performance delta
- Optimize lighting → verify GPU cost reduction

### Unreal Insights Integration

**Full capture lifecycle:**
1. `unreal.Insights.start_capture()` — begin recording
2. Agent makes edits (place props, adjust lighting)
3. `unreal.Insights.stop_capture()` — end recording
4. `unreal.Insights.export_trace()` — save to disk

**Insights viewer shows:**
- Frame breakdown (CPU/GPU time per system)
- Memory allocations
- Draw call hierarchy
- Shader compilation time
- Asset loading patterns

### analyse() — Automated Bottleneck Detection

**Query automated analysis:**
```python
bottlenecks = unreal.PerformanceService.analyse()
# Returns: [ 
#   { "category": "draw_calls", "severity": "high", "recommendation": "Enable LOD..." },
#   { "category": "memory", "severity": "medium", "recommendation": "Compress textures..." }
# ]
```

Agent receives actionable optimization recommendations.

### Standalone Build Profiling

**Profile actual shipped builds (not just editor):**
- Remote profiling (build sends traces back to dev machine)
- Mobile profiling (iOS/Android)
- Console profiling (PS5, Xbox Series X)
- Shipping mode analysis (with full optimization)

---

## Domains VibeUE Covers (19 Toolsets Summary)

| Domain | Toolsets | Example Operations |
|--------|----------|-------------------|
| **Terrain & World** | 6 | Place actor, landscape sculpt, navmesh |
| **Audio** | 2 | Place sound source, set attenuation |
| **Animation** | 4 | Blend skeletal meshes, retarget animation |
| **FX & Mesh** | 3 | Niagara particles, material assignment |
| **Enhanced Input** | 1 | Bind gamepad controls, configure input |
| **Gameplay Logic** | 2 | Edit blueprint, manage data table |
| **Undo/Redo** | 1 | Group edits, rollback, inspect history |
| **Higher-Order Blueprints** | 1 | Create blueprint macro library |
| **Python-First** | 1 | Execute custom Python scripts |
| **Web Research** | 1 | Fetch asset data, integrate web services |

---

## Integration with OneMind/Convergence

### Game Engine + AI Agent Orchestration

**Pattern Recognition:**
- **Vercel Eve (agents)** + **VibeUE (UE editor automation)** = AI-driven game design
- **OneMind NATS (mesh)** → **Unreal Network** (multiplayer sync)
- **Claude Code (reasoning)** + **VibeUE (execution)** = AI-generated content

### Applications

1. **Procedural Game Generation**
   - AI agents generate level layouts
   - Place props, configure spawners
   - Tune difficulty via lighting/enemy placement
   - Real-time profiling ensures performance

2. **AI-Driven Content Pipeline**
   - Agent imports 3D models from web
   - Automatically assigns materials
   - Generates LOD versions
   - Configures physics + collisions
   - Tests frame rate

3. **Multiplayer Arena Design**
   - Agent creates balanced competitive maps
   - Symmetrical spawning points
   - Balanced sight lines
   - Performance-optimized lighting
   - Walkthrough simulation

4. **Environmental Storytelling**
   - Agent places narrative elements
   - Configures atmospheric effects
   - Synchronizes audio + visuals
   - Tests cinematography angles

---

## Comparison: VibeUE vs. Traditional Workflows

| Aspect | Traditional | VibeUE + AI |
|--------|-----------|-------------|
| **Terrain setup** | Manual dragging (30 min) | "Create forest with 100 trees" (30 sec) |
| **Lighting design** | Manual placement & tweaking (2 hrs) | "Match cinematic sunrise" (1 min) |
| **Blueprint editing** | Hand-code state machine (1 hr) | "Create player controller with sprint" (10 sec) |
| **Performance profiling** | Manual testing (1 hr) | Automated recommendations (5 sec) |
| **Content iteration** | Repeat all above (days) | Agent iterates based on feedback (minutes) |
| **Learning curve** | Months to expertise | AI guides workflow, instant expertise |

---

## Resource Links

### Official

- **VibeUE Website:** https://www.vibeue.com/
- **GitHub:** https://github.com/vibeue
- **Discord:** https://discord.gg/vibeue
- **Documentation:** https://www.vibeue.com/docs
- **Fab Marketplace:** https://www.fab.com/listings/vibeue

### Unreal Engine MCP

- **MCP Official Docs:** https://dev.epicgames.com/documentation/unreal-engine/unreal-mcp-in-unreal-editor
- **Unreal Engine Download:** https://www.unrealengine.com/
- **Epic Developer Community:** https://dev.epicgames.com/
- **Unreal Engine Python API:** https://docs.unrealengine.com/en-US/PythonAPI/

### AI Integration

- **Claude Code:** https://www.anthropic.com/claude
- **OpenAI Codex:** https://openai.com/index/openai-codex/
- **Model Context Protocol:** https://modelcontextprotocol.io/

### Learning Resources

- **Unreal Engine 5.8 Release:** https://www.unrealengine.com/
- **Blueprint Best Practices:** https://docs.unrealengine.com/
- **Performance Profiling:** https://docs.unrealengine.com/en-US/TestingAndOptimization/PerformanceProfiling/

---

## Summary

**VibeUE transforms Unreal Engine 5.8 from a tool you manually control into an intelligent system you can reason with.**

Key achievements:
- ✅ Native MCP plugin (no external servers)
- ✅ 19+ domain toolsets (complete editor coverage)
- ✅ Python-first architecture (extensible, open-source)
- ✅ Performance profiling built-in (agents optimize automatically)
- ✅ AI-agent-ready (Claude, Codex, other LLMs)
- ✅ Production-ready (June 2026 release)

**For OneMind Convergence:**
- **Game engine automation** layer for AI-driven content
- **Real-time profiling** integration (performance feedback to agents)
- **Scalable development** (agents handle repetitive tasks)
- **Agentic workflows** (multi-step reasoning in editor)

This is the **game engine layer** that completes the autonomous creative systems stack.

---

**Captured:** 2026-07-11

**Sources:** 
- https://www.vibeue.com/
- https://dev.epicgames.com/documentation/unreal-engine/unreal-mcp-in-unreal-editor

**Status:** Production, actively maintained, open-source community-driven
