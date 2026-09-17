---
id: upgrade-03-diagrams
title: "Upgrade 03: Visual Architecture Diagrams"
type: blueprint
status: deployed
created: 2026-03-27
deployed: 2026-03-27
owner: zeus
tags: [type/blueprint, status/deployed, domain/00-framework]
destination: "_codex/assets/diagrams/"
---

# Upgrade 03: Visual Architecture Diagrams

> ✅ **DEPLOYED** to `_codex/assets/diagrams/` — commit `eb2cfb78` (2026-03-27)

> **Source**: PAI has 23 professional chalk-style infographics explaining every component.
> **Destination**: `_codex/assets/diagrams/`
> **Why**: OneMind has ZERO visual documentation. For a framework meant to be understood, taught, and sold — this is a critical gap. Diagrams are used in: README, book, course, community onboarding, pitch decks.

---

## What's Included

### Mermaid Diagrams (Render in Obsidian)
These work immediately — paste into any `.md` file and Obsidian renders them.

| Diagram | File | Purpose |
|---------|------|---------|
| System Architecture | `mermaid/architecture.md` | Full system overview |
| CODEX Method Flow | `mermaid/codex-flow.md` | C → O → D → E → × loop |
| Domain Map | `mermaid/domain-map.md` | 100 domains across 3+1 groups |
| Memory Architecture | `mermaid/memory-tiers.md` | 3-tier memory flow |
| Agent Architecture | `mermaid/agent-architecture.md` | How 4 agents coordinate |
| Data Flow | `mermaid/data-flow.md` | How information moves through the vault |

### Image Generation Briefs
Detailed prompts for generating PAI-quality infographics using Midjourney, DALL-E, or a designer.

| Brief | File | Style |
|-------|------|-------|
| System Overview | `image-briefs/system-overview.md` | Chalk-on-dark infographic |
| CODEX Method | `image-briefs/codex-method.md` | Circular flow diagram |
| Domain Wheel | `image-briefs/domain-wheel.md` | 100-spoke radial chart |
| Before/After | `image-briefs/before-after.md` | Split comparison |
| Agent Coordination | `image-briefs/agent-coordination.md` | Network diagram |
| Memory Flow | `image-briefs/memory-flow.md` | Layered architecture |
| OneMind vs Second Brain | `image-briefs/vs-second-brain.md` | Feature comparison |
| The Equation | `image-briefs/the-equation.md` | Hero graphic |

---

## How to Deploy

1. Copy Mermaid `.md` files to wherever they're referenced
2. Use image briefs with Midjourney/DALL-E to generate PNGs
3. Save generated images to `_codex/assets/diagrams/`
4. Reference in README, ONEMIND-CODEX.md, and docs

## PAI Equivalent

PAI commissioned 23 custom chalk-style infographics (dark background, hand-drawn aesthetic). They appear in the README and throughout docs. Our Mermaid diagrams give us instant functional diagrams; the image briefs let us generate matching quality when ready.
