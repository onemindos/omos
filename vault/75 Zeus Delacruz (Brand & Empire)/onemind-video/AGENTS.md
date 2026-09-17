# AGENTS.md — One Mind Video Production

> Instructions for AI agents working in this repository.

---

## What This Is

**onemind-video** is a Remotion project that produces educational video content for the One Mind framework. Zeus (human) and Legacy (AI) create visual explainer episodes together using React components, OpenAI TTS voice generation, and ffmpeg for audio processing.

| Field | Value |
|-------|-------|
| Framework | Remotion 4.0.446 |
| Language | TypeScript, React 19 |
| Resolution | 1920×1080 @ 30fps |
| Output | H.264 MP4 |
| Location | `~/Documents/onemind-video/` |
| Skill file | `~/Documents/onemind-codex/.claude/skills/video-production/SKILL.md` |

**Read the skill file first** — it has the complete pipeline, component docs, voice settings, and troubleshooting.

---

## Quick Reference

### Build & Run

```bash
npm run studio             # Preview at localhost:3000
npx tsc --noEmit           # Type-check (must be zero errors before render)
npm run render:intro-v2    # Render current episode → out/onemind-intro-v2.mp4
```

### Generate Voice (V2 — Two Voices)

```bash
npm run voice:intro-v2
# Or with explicit key:
OPENAI_API_KEY=$(op read 'op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key') \
  npx tsx src/scripts/generate-intro-v2-voice.ts
```

### Generate SFX

SFX are pre-generated in `public/audio/sfx/`. Only regenerate if needed — see skill file for ffmpeg commands.

---

## Current State (April 2026)

**Active episode:** `OneMindIntroV2` — Zeus + Legacy two-voice dialogue, 5:04, 22.1 MB. Rendered and ready for review.

**Episode registry:**

| Composition | Duration | Status |
|-------------|----------|--------|
| `OneMindIntroV2` | 5:04 | Rendered (current) |
| `OneMindIntro` | 3:39 | Rendered (superseded by V2) |
| `OpenClawDreaming` | 3:20 | Rendered |
| `OneMindPromoV4` | 1:00 | Rendered (rejected) |
| `OneMindPromoV3` | 1:27 | Rendered |
| `OneMindPromo` | 2:11 | Rendered |
| `ComposioCLI` | 1:00 | Rendered |

---

## Two Episode Formats

### V2 — Zeus + Legacy Dialogue (current standard)
- Two speakers: Zeus (`onyx` voice @ 1.0x) and Legacy (`nova` voice @ 0.92x)
- Script format: `ZEUS: "text"` and `LEGACY: "text"` lines in `scripts/<name>-script.md`
- Generator: `src/scripts/generate-intro-v2-voice.ts` → outputs `timing.json` with frame counts
- Visual layers: Explainer.tsx + MotionFX.tsx + SFX + ambient background music
- Composition uses `ConversationLayout` with frame-based `activeSpeaker` switching

### V1 — Legacy Monologue (deprecated for new episodes)
- Single speaker: Legacy only (`nova` @ 0.92x)
- Script format: plain text lines under `## Segment N` headers
- Generator: `src/scripts/generate-intro-voice.ts` or `generate-dreaming-voice.ts`
- Visual layer: Explainer.tsx components only, no SFX system

---

## Component Libraries

### Explainer.tsx (18 components) — Core visuals
Backgrounds, content slides, diagrams, code blocks, branding elements.

**CRITICAL prop gotchas:**
- `BulletList` → `items: string[]` (NOT `bullets`)
- `FlowDiagram` → does NOT accept `accentColor`
- `BigCallout` → does NOT accept `emoji`
- `CodeBlock` → `lines: string[]` (NOT `code: string`)
- `SplitCompare` → nested `left: { title, items }` and `right: { title, items }`
- `ArchDiagram` → `width` (NOT `isWide`)

### MotionFX.tsx (9 components) — Enhanced motion graphics
`ParticleField`, `SceneTransition`, `SpeakerBubble`, `FloatingEmoji`, `PulseRing`, `GlowText`, `ConversationLayout`, `KineticText`, `EnergyBurst`.

All use deterministic pseudo-random (golden angle seed) — no `Math.random()`.

---

## Audio Architecture (V2)

Three layers:
1. **Voice** — Per-segment MP3s in `public/audio/<episode>/`
2. **SFX** — Per-segment triggers from `public/audio/sfx/` (whoosh, reveal, impact)
3. **Ambient** — Global `ambient-bg.mp3` at volume 0.08

**CRITICAL:** When concatenating voice clips, use ffmpeg `filter_complex` concat method, NOT the concat demuxer (fails on sample rate mismatch between TTS 24kHz output and generated silence).

---

## Voice Settings

| Role | Voice | Speed | Character |
|------|-------|-------|-----------|
| Zeus | `onyx` | 1.0× | Deep, warm, energetic |
| Legacy | `nova` | 0.92× | Warm, expressive |

Provider: OpenAI TTS-HD (`tts-1-hd`)
API Key: `op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key` (1Password CLI)

---

## Brand Tokens (from Explainer.tsx)

| Token | Value | Usage |
|-------|-------|-------|
| `bg` | `#0a0a0f` | Dark background |
| `teal` | `#64ffda` | Legacy's color |
| `amber` | `#fbbf24` | Zeus's color |
| `onemindRed` | `#e63946` | One Mind brand red |
| `cyan` | `#22d3ee` | Secondary accent |
| `purple` | `#a78bfa` | Thoughtful concepts |
| `coral` | `#f87171` | Errors, "without" states |
| `green` | `#34d399` | Success, "with" states |

---

## Known Issues

- `ambient-bg.mp3` is a synthetic ffmpeg pad — replace with Suno/royalty-free lo-fi for production
- `KineticText` component built but unused — good for key phrase emphasis
- No screenshot/photo slide component yet (ImageSlide planned)
- `activeSpeaker` frame timings are estimates — may need adjustment after watching
- No face-cam integration yet — future: intercut Zeus footage in CapCut

---

## What NOT to Do

- Do not use `Math.random()` in any component — Remotion requires deterministic rendering
- Do not use ffmpeg concat demuxer for joining TTS clips — use filter_complex
- Do not change voice assignments (Zeus=onyx, Legacy=nova) without updating all scripts
- Do not render without running `npx tsc --noEmit` first — type errors = broken video

---

_Part of the One Mind ecosystem. Full skill documentation: `~/Documents/onemind-codex/.claude/skills/video-production/SKILL.md`_
