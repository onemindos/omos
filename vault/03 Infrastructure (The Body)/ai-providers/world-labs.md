---
type: ai-provider
category: 3d-generation
vendor: World Labs
founded: 2024
status: active
sovereignty: ❌ cloud-only
created: 2026-05-24
tags:
  - spatial-intelligence
  - 3d-generation
  - gaussian-splatting
  - marble
---

# World Labs — Spatial Intelligence Platform

**URL:** https://www.worldlabs.ai/  
**Platform:** https://marble.worldlabs.ai/  
**API:** https://platform.worldlabs.ai/  
**Founder:** Fei-Fei Li (ex-Google, Stanford AI Lab)

## What is World Labs?

**Spatial intelligence company** building frontier models for 3D world generation, perception, and interaction.

### Mission Statement (from homepage)
> "World Labs is building the next frontier of generative AI — one where models can understand and interact with the world to empower use cases from storytelling to simulation."

> "Spatial intelligence transforms seeing into doing, understanding into reasoning, and imagining into creating."

## Key Product: **Marble**

**Text/Image-to-3D world generation platform**

### Capabilities
- **Multimodal inputs:** Text prompts, image prompts, sketches
- **3D Gaussian Splatting (3DGS):** High-fidelity rendering
- **Streamable LoD:** Level-of-detail system (Spark 2.0 tech)
- **Persistent worlds:** Explore, edit, modify generated scenes
- **Web-based:** Real-time 3D in browser (no local GPU needed)

### Example Prompts (from homepage demo)
- "Autumn forest with a wooden cabin"
- "Futuristic amphitheater in a neon-lit city"
- "Japanese garden with a pond and stone bridge"
- "Anime-style train station at sunset"
- "Modern bathroom with marble tiles"

### Tech Stack (from blog)
- **Gaussian Splatting:** `.spz` format (compressed 3DGS)
- **Streaming:** 100k Gaussian splats per scene
- **WebGL/WebGPU:** Browser-native rendering
- **Spark 2.0:** LoD system for scalable streaming

## OneMind Integration Potential

### Use Cases

#### 1. Virtual World Building
- Generate 3D environments for planning/visualization
- Prototype home renovation designs
- Create virtual sets for video content

**Example:**
- **Prompt:** "Modern home office with glass walls, standing desk, and mountain view"
- **Output:** 3D scene you can walk through

#### 2. Drone Flight Simulation
- Generate test environments for drone navigation
- Obstacle course generation
- Physics-based flight testing

**Example:**
- **Prompt:** "Forest with trees, clearings, and a landing pad"
- **Use:** Test drone pathfinding algorithms

#### 3. Architecture Visualization
- Visualize building designs before construction
- Client presentations
- Interior design mockups

**Example:**
- **Prompt:** "Living room with minimalist furniture, floor-to-ceiling windows, and city view"
- **Use:** Show Zeus design options for future home

#### 4. Digital Twin (Home)
- 3D model of your home for IoT planning
- Security camera coverage visualization
- Robot/drone navigation maps

**Example:**
- **Prompt:** Upload photos of each room → Generate 3D model
- **Use:** Plan smart home device placement

#### 5. Storytelling & Content Creation
- Generate scenes for YouTube videos
- Virtual backgrounds for presentations
- Game environment prototyping

### API Integration (Research Needed)
- **Platform URL:** https://platform.worldlabs.ai/
- **Pricing:** Unknown (research needed)
- **Rate limits:** Unknown
- **API format:** REST? gRPC? GraphQL?

## Sovereignty Analysis ⚠️

| Factor | Status |
|--------|--------|
| **Cloud dependency** | ❌ Cloud-only (no local deployment) |
| **Data privacy** | ⚠️ Prompts/images sent to World Labs servers |
| **Availability** | ⚠️ Requires internet, subject to rate limits |
| **Cost** | ⚠️ Unknown pricing (could be expensive) |
| **Lock-in** | ⚠️ Proprietary format (.spz) |

### Alternatives (Sovereign)
- **Blender + Procedural Generation:** Free, local, scriptable
- **Unity/Unreal Engine:** Free (personal use), local
- **ComfyUI + ControlNet:** Text-to-3D via Stable Diffusion (local)
- **Gaussian Splatting (open-source):** GitHub repos for DIY 3DGS

## Recommendation

### Decision Tree
```
Do you NEED 3D world generation for critical operations?
├─ YES → Evaluate alternatives first (Blender, Unity)
│   └─ Still need AI generation? → Test Marble (exploratory)
└─ NO → Skip for now (nice-to-have, not core OneMind function)
```

### Priority: **Low** (Exploratory)
- **Not critical path** for OneMind OS
- **Nice-to-have** for content creation, visualization
- **Wait for API pricing** before committing
- **Test free tier first** (if available)

### When to Revisit
1. **Drone simulation needs** become urgent
2. **Architecture/home design** projects start
3. **API pricing** is revealed and affordable
4. **Local deployment** becomes available (unlikely)

## Research Blog Posts (from homepage)

### 1. "Streaming 3DGS worlds on the web" (April 14, 2026)
- **Tech:** Spark 2.0's Level-of-Detail system
- **Format:** `.spz` (compressed 3D Gaussian Splatting)
- **Streaming:** Progressively load 100k+ splats
- **Relevance:** Technical deep dive for understanding pipeline

### 2. "3D as code" (March 3, 2026)
- **Concept:** Text = universal interface for software, 3D = universal interface for space
- **Vision:** Humans + AI systems generate/edit/share worlds together
- **Relevance:** Philosophy behind World Labs' mission

### 3. "World Labs Announces New Funding" (Feb 18, 2026)
- **Funding:** Likely Series A/B (amount not disclosed)
- **Investors:** Unknown
- **Relevance:** Company has runway for development

## Next Steps (Low Priority)

### Phase 1: Exploration (Optional)
- [ ] Create World Labs account (https://marble.worldlabs.ai/)
- [ ] Test Marble platform with sample prompts
- [ ] Evaluate output quality vs. Blender
- [ ] Research API pricing (https://platform.worldlabs.ai/)

### Phase 2: Integration (If Valuable)
- [ ] Test API for OneMind use cases
- [ ] Build Hermes skill: `world-labs-generation`
- [ ] Design NATS bridge (if API is async)
- [ ] Evaluate cost per generation

### Phase 3: Alternatives (If Expensive)
- [ ] Research local 3DGS tools (GitHub)
- [ ] Test Blender procedural generation
- [ ] Evaluate ComfyUI text-to-3D workflows

## Related Files
- YOLOv8 (vision AI): `ai-providers/yolov8-ultralytics.md`
- Pi AI HAT 2 (edge vision): `hardware/pi-ai-hat-2.md`
- Comfy (alternative 3D/image gen): `ai-providers/comfy.md`
