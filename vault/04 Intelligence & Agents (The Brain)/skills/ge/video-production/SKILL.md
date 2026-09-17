---
name: video-production
description: "Manage the One Mind video production pipeline using Remotion, OpenAI TTS, and CapCut — write scripts, generate voices, build Remotion compositions, and render episodes. Use when creating, rendering, or editing educational video content, scripting a new episode, generating voice clips, or assembling in CapCut."
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Video Production Skill — One Mind Visual Explainers

> How Zeus + Legacy create educational video content using Remotion, OpenAI TTS, and CapCut.
> Updated: 2026-04-09 — includes V2 two-voice dialogue pipeline, MotionFX components, SFX system.

For component library details see `references/component-library.md`. For SFX generation see `references/sfx-system.md`. For episode registry and npm scripts see `references/episode-registry.md`. For production costs and real-life examples see `references/production-notes.md`.

---

## Philosophy

**This is a collaboration, not a delegation.** Zeus brings the vision, the real-life experience, and the face-cam presence. Legacy brings the visual explanation engine, voice narration, and production automation. Together = One Mind content.

**V2 Principle:** Videos should feel like two people having a real conversation — podcast energy, not lecture. Zeus opens casually ("Alright everybody"), Legacy adds context and humor, they interrupt each other, joke around, and keep it human.

---

## Project Location

```
~/Documents/onemind-video/
├── src/
│   ├── Root.tsx                      # Composition registry — add new episodes here
│   ├── components/
│   │   ├── Explainer.tsx             # Core visual explainer library (18 components)
│   │   ├── MotionFX.tsx              # Enhanced motion graphics library (9 components)
│   │   └── ParticleBackground.tsx    # Legacy particle effect (from v1-v3 promos)
│   ├── OneMindIntroV2.tsx            # Episode: Two-voice Zeus+Legacy intro (CURRENT)
│   └── scripts/
│       ├── generate-intro-v2-voice.ts  # Two-voice TTS (Zeus onyx + Legacy nova)
│       ├── generate-legacy-voice.ts    # Generic Legacy voice from text lines
│       ├── generate-script.ts          # AI script generator (GPT-4o)
│       ├── generate-youtube-metadata.ts
│       └── generate-social-captions.ts
├── public/
│   ├── audio/
│   │   ├── onemind-intro-v2/         # Two-voice clips + timing.json
│   │   └── sfx/                      # whoosh.mp3, reveal.mp3, impact.mp3, ambient-bg.mp3
│   └── images/onemind-logo.png
├── scripts/                          # Episode dialogue scripts (.md)
├── out/                              # Rendered MP4 output
└── templates/script-template.md
```

---

## Video Formats

**V2 (current standard) — Zeus + Legacy Dialogue**
- Two speakers: Zeus and Legacy have a real conversation
- Script format: `ZEUS: "text"` and `LEGACY: "text"` lines
- Voices: Zeus = `onyx` @ 1.0x, Legacy = `nova` @ 0.92x
- Generator: `generate-intro-v2-voice.ts` (outputs timing.json)
- Visual layer: Explainer.tsx + MotionFX.tsx

**V1 (deprecated) — Legacy Monologue**
- Single speaker: Legacy narrates everything; `nova` only; V1 generators

---

## Design System (BRAND tokens)

| Token | Value | Usage |
|-------|-------|-------|
| `bg` | `#0a0a0f` | Dark background |
| `bgCard` | `#12121a` | Card/panel background |
| `teal` | `#64ffda` | Primary accent, Legacy's color |
| `cyan` | `#22d3ee` | Secondary accent |
| `purple` | `#a78bfa` | Dreaming / thoughtful concepts |
| `amber` | `#fbbf24` | Warnings, highlights, Zeus references |
| `coral` | `#f87171` | Errors, "without" states |
| `green` | `#34d399` | Success, "with" states |
| `onemindRed` | `#e63946` | One Mind brand red (from logo) |

---

## Voice Settings

| Role | Voice | Speed | Character |
|------|-------|-------|-----------|
| Zeus | `onyx` | 1.0× | Deep, warm, energetic — the human host |
| Legacy | `nova` | 0.92× | Warm, expressive — the AI partner |

Provider: OpenAI TTS-HD (`tts-1-hd`). API Key: `op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key`

Timing constants: `PAUSE_BETWEEN_SPEAKERS_MS = 350ms` | `PAUSE_SAME_SPEAKER_MS = 150ms`

**CRITICAL:** Do NOT use ffmpeg's concat demuxer for joining voice clips — use `filter_complex` concat (silence clips must be 24kHz mono to match TTS output format).

---

## How to Create a New Episode (V2)

### Step 1: Write the Dialogue Script

Create `scripts/<episode-name>-script.md`:

```markdown
# Episode Name — Zeus + Legacy Dialogue Script

> Zeus Voice: OpenAI TTS-HD, model "onyx", speed 1.0
> Legacy Voice: OpenAI TTS-HD, model "nova", speed 0.92
> Format: Two hosts chatting — podcast energy, not lecture

---

## Segment 01 — Cold Open
ZEUS: "Alright everybody — let me show you something cool."
LEGACY: "And here's the context he's leaving out."
```

Writing guidelines: podcast energy, Zeus opens casually, Legacy adds depth and humor, beginner-friendly, real Zeus examples, one topic per segment.

### Step 2: Generate Two-Voice Clips

```bash
cd ~/Documents/onemind-video
OPENAI_API_KEY=$(op read 'op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key') \
  npx tsx src/scripts/generate-intro-v2-voice.ts
```

Outputs: `public/audio/<episode>/segment-NN.mp3` + `timing.json` with per-segment frame counts.

### Step 3: Build the Remotion Composition

1. Create `src/<EpisodeName>.tsx`
2. Import from `Explainer.tsx` and `MotionFX.tsx`
3. Build `SEG` timing object from `timing.json` (frames from JSON + 30 padding per segment)
4. Calculate cumulative `OFF` offsets and export `TOTAL_FRAMES`
5. Wrap segments in `<ConversationLayout activeSpeaker="zeus"|"legacy"|"both">` with frame-based switching
6. Add per-segment `<Audio>` for voice, per-segment SFX, and global ambient at volume 0.08

```tsx
const SEG = {
  s01: { audio: 199, total: 229 },  // frames from timing.json + 30 padding
  s02: { audio: 718, total: 748 },
};
const OFF = (() => {
  const keys = Object.keys(SEG) as (keyof typeof SEG)[];
  const result: Record<string, number> = {};
  let pos = 0;
  for (const k of keys) { result[k] = pos; pos += SEG[k].total; }
  result.total = pos;
  return result;
})();
export const TOTAL_FRAMES = OFF.total;
```

### Step 4: Register in Root.tsx

```tsx
import { EpisodeName, TOTAL_FRAMES } from "./EpisodeName";
<Composition id="EpisodeName" component={EpisodeName}
  durationInFrames={TOTAL_FRAMES} fps={30} width={1920} height={1080} />
```

### Step 5: Type-Check, Preview, Render

```bash
npx tsc --noEmit              # Must be zero errors
npm run studio                 # Preview at localhost:3000
npm run render:<episode>       # Render to out/
```

### Step 6: Verify Output

```bash
ffprobe -v quiet -show_entries format=duration,size -of default=noprint_wrappers=1 out/<episode>.mp4
```

### Step 7 (Optional): Final Assembly in CapCut

Import rendered Remotion MP4 + Zeus face-cam, intercut, add music, export 1080p @ 30fps.

---

## Audio Architecture (V2)

```
Layer 1: Per-segment voice clips (segment-01.mp3 ... segment-NN.mp3)
Layer 2: Per-segment SFX (whoosh for transitions, reveal/impact for emphasis)
Layer 3: Global ambient music (ambient-bg.mp3 at volume 0.08)
```

---

## Video Format Standards

- Resolution: 1920×1080 (16:9) — Frame rate: 30fps — Codec: H.264
- Style: Dark background (`#0a0a0f`), animated motion graphics, One Mind brand colors
- Logo: Always present as watermark

---

## Troubleshooting

**ffmpeg concat fails with "Could not find codec parameters"**
Generate silence at 24kHz mono: `ffmpeg -f lavfi -i anullsrc=r=24000:cl=mono -t 0.35 -y silence.mp3`
Then use filter_complex: `ffmpeg -i clip1.mp3 -i silence.mp3 -i clip2.mp3 -filter_complex "[0:a][1:a][2:a]concat=n=3:v=0:a=1[out]" -map "[out]" output.mp3`

**1Password authorization timeout**
Run `op signin` first, then retry. Or: `export OPENAI_API_KEY=$(op read '...')`

**Type errors after adding MotionFX components**
Import from `./components/MotionFX` (not `./components/Explainer`).

---

## Related

- `references/component-library.md` — full Explainer.tsx + MotionFX.tsx component API
- `references/sfx-system.md` — SFX generation commands and usage
- `references/episode-registry.md` — episode list and npm scripts
- `references/production-notes.md` — costs, Zeus's real numbers, known issues
