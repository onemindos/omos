# Component Library Reference — Video Production

## Explainer.tsx — Core Visuals (18 components)

**Backgrounds:**
- `DarkGrid` — Animated dark grid with radial glow pulse
- `OneMindGradient` — Brand gradient background (teal → purple)

**Content Slides:**
- `TitleSlide` — Big title + subtitle + animated accent line
- `BulletList` — Staggered-reveal bullet points. **Props:** `items: string[]` (NOT `bullets`)
- `BigCallout` — Giant number with label. **Does NOT accept `emoji` prop**
- `SplitCompare` — Side-by-side comparison. **Props:** nested `left: { title, items }` and `right: { title, items }`
- `StatCounter` — Animated counting number with label

**Diagrams:**
- `FlowDiagram` — Horizontal step-by-step flow with arrows. **Does NOT accept `accentColor`**
- `ArchDiagram` — Layered architecture boxes. **Uses `width` NOT `isWide`**

**Code:**
- `CodeBlock` — Terminal-style typewriter. **Props:** `lines: string[]` (NOT `code: string`)

**Branding:**
- `LegacyBadge` — Animated waveform badge showing "LEGACY AI"
- `ZeusBadge` — Shows "ZEUS" badge
- `OneMindLogo` — Logo component
- `OneMindWatermark` — Logo watermark in corner
- `LowerThird` — Sliding name/topic bar

**Episode-Specific:**
- `PhaseCard` — 3-column phase cards (OpenClaw Dreaming)
- `ScoringBar` — Animated horizontal bar with percentage

> **CRITICAL prop signatures** — These are easy to get wrong:
> - `BulletList` takes `items: string[]` not `bullets`
> - `FlowDiagram` does NOT accept `accentColor`
> - `BigCallout` does NOT accept `emoji`
> - `CodeBlock` uses `lines: string[]` NOT `code: string`
> - `SplitCompare` uses nested `left/right` objects
> - `ArchDiagram` uses `width` not `isWide`

---

## MotionFX.tsx — Enhanced Motion Graphics (9 components)

Added in V2 for visual energy and dialogue support:

| Component | Purpose | Key Props |
| --------- | ------- | --------- |
| `ParticleField` | Floating glowing dots — ambient energy | `count`, `color`, `speed`, `maxSize` |
| `SceneTransition` | Segment transition effects | `type: "zoom" \| "wipe" \| "glitch" \| "fade"`, `duration`, `color` |
| `SpeakerBubble` | Who's talking indicator with waveform | `speaker: "zeus" \| "legacy"`, `text`, `position` |
| `FloatingEmoji` | Emojis that float up and fade | `emojis: string[]`, `density`, `speed` |
| `PulseRing` | Expanding rings for emphasis | `color`, `count`, `maxSize`, `speed` |
| `GlowText` | Text with animated glow | `text`, `color`, `fontSize`, `glowIntensity` |
| `ConversationLayout` | Split speaker indicators for dialogue | `activeSpeaker: "zeus" \| "legacy" \| "both"`, `children` |
| `KineticText` | Word-by-word reveal animation | `text`, `highlightWords`, `highlightColor`, `wordsPerBeat` |
| `EnergyBurst` | Radial burst for emphasis moments | `color`, `delay` |

**All MotionFX components use deterministic pseudo-random** (golden angle seed) for Remotion compatibility — no `Math.random()`.

**`ConversationLayout`** is the key V2 component — it shows Zeus (top-left, ⚡, amber) and Legacy (top-right, 🧠, teal) indicators with animated waveform bars for the active speaker. Use frame-based switching:

```tsx
const activeSpeaker = frame < 90 ? "zeus" : frame < 200 ? "legacy" : "zeus";
<ConversationLayout activeSpeaker={activeSpeaker}>
  {/* segment content */}
</ConversationLayout>
```

---

## ParticleBackground.tsx — Legacy Particles (from v1-v3 promos)

Older particle effect from promo videos. For new episodes, prefer `ParticleField` from MotionFX.tsx.
