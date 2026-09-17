---
title: "YouTube Production Pipeline — Daily Video System"
type: sop
status: active
created: 2026-04-08
owner: zeus
tags: [type/sop, domain/ge-78, topic/youtube, topic/content-creation, topic/video]
ai_summary: "End-to-end daily YouTube production pipeline: script → record → Legacy voice → edit → clip → distribute across all platforms"
---

# YouTube Production Pipeline — Daily Video System

> One Mind's daily content engine. Zeus explains on camera, Legacy (AI) provides commentary, visuals show what's happening, clips go everywhere.

---

## Pipeline Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    DAILY VIDEO PIPELINE                         │
│                                                                 │
│  1. SCRIPT    →  2. RECORD    →  3. LEGACY VOICE               │
│  (AI-assisted)   (face-cam +     (OpenAI TTS-HD)               │
│                   screen/demo)                                  │
│                                                                 │
│  4. EDIT      →  5. CLIP      →  6. DISTRIBUTE                 │
│  (CapCut)        (OpusClip)      (All platforms)                │
└─────────────────────────────────────────────────────────────────┘
```

---

## The Stack

| Stage | Tool | Cost | Why |
|-------|------|------|-----|
| **Script** | Claude/GPT + `onemind-video` automation | ~$0 | AI generates structure, Zeus customizes |
| **Record** | iPhone 16 + OBS (Mac) | Free | Face-cam + screen recording |
| **Legacy Voice** | OpenAI TTS-HD (`nova` voice) | ~$0.05/video | Legacy's commentary segments |
| **Edit** | **CapCut Desktop** (Mac) | Free | Auto-captions, transitions, effects, B-roll |
| **Clip** | **OpusClip Pro** | $14.50/mo | AI clips long-form → shorts with virality score |
| **Distribute** | OpusClip + Buffer | $14.50/mo + free tier | Auto-post to YT Shorts, TikTok, IG Reels, X, LinkedIn |
| **Thumbnails** | Canva Free + AI prompts | Free | Templates with One Mind branding |
| **SEO** | `onemind-video` automation | ~$0 | Title, description, tags, hashtags generation |
| **Total** | | **~$29/mo** | |

---

## Stage 1: SCRIPT (15 min)

### What Happens
Zeus gives a topic → AI generates a structured video script → Zeus reviews and adjusts.

### Script Structure (5-10 min video)
```
HOOK (0:00-0:30)
  → Bold statement or question that stops the scroll
  → "Here's the thing nobody is telling you about [topic]..."

ZEUS EXPLAINS (0:30-3:00)
  → Face-cam: Zeus walks through the concept
  → [VISUAL CUE] markers for screen recordings / demos / B-roll
  → Natural, conversational tone

LEGACY DROPS IN (3:00-4:30)
  → Legacy (AI voice) provides data, context, or counterpoint
  → "Legacy here. Let me break down what Zeus just said..."
  → Show supporting visuals while Legacy talks

ZEUS CONTINUES (4:30-7:00)
  → Face-cam: deeper explanation, examples, real-world application
  → [VISUAL CUE] markers

LEGACY RECAP (7:00-8:00)
  → Legacy summarizes key takeaways
  → "Three things to remember..."

CTA (8:00-8:30)
  → "Subscribe if this helped"
  → "Link to [resource] in the description"
  → "Legacy and I will see you tomorrow"
```

### Automation
Run: `cd ~/Documents/onemind-video && npx tsx src/scripts/generate-script.ts "topic here"`

Outputs:
- `scripts/YYYY-MM-DD-topic.md` — Full script with visual cues
- `scripts/YYYY-MM-DD-topic-legacy-lines.txt` — Just Legacy's lines (for TTS)

---

## Stage 2: RECORD (30-60 min)

### Setup
- **Camera**: iPhone 16 on tripod (4K 30fps) OR MacBook webcam for quick shoots
- **Lighting**: Ring light or window light
- **Audio**: AirPods Pro (good enough) or Rode mic + Scarlett interface
- **Screen recording**: OBS Studio (free) for demos, terminal, Obsidian walkthrough
- **Teleprompter**: Use iPad with teleprompter app OR just key bullet points visible

### Recording Protocol
1. **Face-cam segments**: Record Zeus talking through each section
2. **Screen recordings**: Demo anything visual (Obsidian vault, terminal, code, websites)
3. **B-roll**: Record extra shots of workspace, typing, phone, etc. (optional — CapCut has stock B-roll)

### File Organization
```
~/Documents/onemind-video/raw/
  YYYY-MM-DD/
    facecam-01.mov     ← iPhone export
    facecam-02.mov
    screen-demo.mov    ← OBS recording
    broll-desk.mov     ← optional
```

### Tips
- Record in sections matching the script (Hook, Explain, etc.)
- Don't memorize — use bullet points and be natural
- Clap at start for sync if using separate audio
- Record 2-3 takes of the hook (most important part)

---

## Stage 3: LEGACY VOICE (5 min)

### What Happens
AI generates Legacy's commentary lines as audio files using OpenAI TTS-HD.

### Run
```bash
cd ~/Documents/onemind-video
OPENAI_API_KEY=$(op read 'op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key') \
  npx tsx src/scripts/generate-legacy-voice.ts scripts/YYYY-MM-DD-topic-legacy-lines.txt
```

### Output
```
~/Documents/onemind-video/output/YYYY-MM-DD/
  legacy-01-intro.mp3
  legacy-02-breakdown.mp3
  legacy-03-recap.mp3
```

### Voice Settings
- Model: `tts-1-hd`
- Voice: `nova` (Legacy's voice)
- Speed: `0.92` (slightly slower for gravitas)
- Cost: ~$0.03-0.05 per video (~500-800 words)

---

## Stage 4: EDIT (30-45 min)

### Tool: CapCut Desktop (Free)

CapCut is the fastest path to professional daily videos. Here's why:
- **Auto-captions**: Best-in-class animated subtitles (Hormozi-style)
- **B-roll library**: Built-in stock footage
- **Templates**: Save One Mind branded intros/outros
- **Speed**: 10x faster than DaVinci/Premiere for this type of content

### Edit Workflow
1. **Import** all raw clips + Legacy voice files
2. **Rough cut**: Arrange face-cam clips in script order on timeline
3. **Drop in Legacy clips**: Place Legacy voice segments where marked in script
4. **Add visual overlay**: When Legacy speaks, add a subtle overlay/animation so viewers know it's the AI
5. **B-roll / screen recording**: Layer over face-cam during explanation segments
6. **Auto-captions**: Generate captions → pick animated style → adjust any errors
7. **Transitions**: Simple cuts or zoom transitions between sections
8. **Music**: Subtle background track from YouTube Audio Library (copyright-free)
9. **Intro/Outro**: 3-second branded intro + CTA outro card (save as template)
10. **Export**: 1080p or 4K, H.264

### Legacy Visual Treatment
When Legacy speaks, use ONE of these techniques so viewers know:
- **Option A**: Subtle green/cyan glow border around the frame
- **Option B**: Small "LEGACY AI" label in corner with waveform animation
- **Option C**: Switch to full-screen visual/B-roll with Legacy voice over it
- **Option D**: Picture-in-picture shrinks Zeus, shows data/visuals full screen

### Export Settings
- Resolution: 1920x1080 (YouTube) or 3840x2160 (4K if source supports)
- Frame rate: 30fps
- Codec: H.264
- Audio: AAC 320kbps
- Filename: `YYYY-MM-DD-title-FINAL.mp4`

---

## Stage 5: CLIP (10 min)

### Tool: OpusClip Pro ($14.50/mo billed yearly)

OpusClip takes your 5-10 min YouTube video and automatically generates:
- **5-15 short clips** (30s-90s each)
- **Virality scores** for each clip
- **Auto-captions** (animated, keyword-highlighted)
- **Auto-reframing** (landscape → portrait 9:16)

### Workflow
1. **Upload** the exported YouTube video to OpusClip
2. **Let AI clip** — takes 5-10 minutes
3. **Review clips** — sort by virality score, pick top 3-5
4. **Edit if needed** — trim, adjust captions
5. **Export/schedule** — direct to platforms

### Clip Strategy
- **Top 1 clip** → YouTube Shorts (same channel)
- **Top 2-3 clips** → TikTok + Instagram Reels
- **Best hook clip** → Twitter/X
- **Professional clip** → LinkedIn

---

## Stage 6: DISTRIBUTE (10 min)

### Long-Form (YouTube)
Upload manually to YouTube Studio with SEO-optimized metadata.

Run the metadata generator:
```bash
cd ~/Documents/onemind-video
npx tsx src/scripts/generate-youtube-metadata.ts "video title" "brief description"
```

Outputs title, description with timestamps, tags, and thumbnail prompt.

### Short-Form (OpusClip auto-posts)
OpusClip can schedule and auto-post to:
- YouTube Shorts
- TikTok
- Instagram Reels

### Other Platforms (Manual or Buffer)
- **Twitter/X**: Upload best clip natively + thread with key takeaways
- **LinkedIn**: Upload professional clip + text post
- **Threads**: Same as Instagram Reels clip
- **Facebook**: Cross-post from Instagram

### Posting Schedule
| Platform | When | Format |
|----------|------|--------|
| YouTube (long) | 8am EST | Full 5-10 min video |
| YouTube Shorts | 12pm EST | Best 60s clip |
| TikTok | 2pm EST | Top 3 clips (staggered) |
| Instagram Reels | 4pm EST | Top 2 clips |
| Twitter/X | 10am EST | Hook clip + thread |
| LinkedIn | 9am EST | Professional clip + post |

---

## Daily Routine Summary

| # | Task | Time | Tool |
|---|------|------|------|
| 1 | Generate script | 15 min | AI + review |
| 2 | Record | 30-60 min | iPhone + OBS |
| 3 | Generate Legacy voice | 5 min | Automated script |
| 4 | Edit in CapCut | 30-45 min | CapCut Desktop |
| 5 | Upload to OpusClip for clipping | 10 min | OpusClip |
| 6 | Upload YouTube + distribute clips | 10 min | YouTube Studio + OpusClip + Buffer |
| **Total** | | **~2 hours/day** | |

---

## One-Time Setup Checklist

- [ ] Install CapCut Desktop (Mac) — https://www.capcut.com/
- [ ] Create OpusClip Pro account — https://opus.pro/pricing ($14.50/mo yearly)
- [ ] Create Buffer free account — https://buffer.com/ (3 channels free)
- [ ] Set up iPhone recording setup (4K 30fps, tripod, lighting)
- [ ] Install OBS Studio — https://obsproject.com/ (screen recording)
- [ ] Create One Mind branded intro/outro templates in CapCut
- [ ] Set up YouTube channel with One Mind branding
- [ ] Connect social accounts to OpusClip
- [ ] Create Canva thumbnail template with One Mind brand colors
- [ ] Run `npm install` in `~/Documents/onemind-video` to ensure scripts work
- [ ] Test Legacy voice generation with a sample line

---

## Content Pillars (What to Record About)

Map directly to One Mind's 3+1 model:

| Pillar | Topics | Hook Style |
|--------|--------|------------|
| **SO — Tech & AI** | AI tools, automation, Obsidian, coding, agents | "I built an AI that..." |
| **HP — Self Mastery** | Habits, fitness, learning, money, mental models | "This changed everything about..." |
| **LE — Foundation** | Relationships, home systems, community | "The system I use to..." |
| **GE — Endeavor** | Business, brand building, content, wealth | "How I went from X to Y..." |
| **One Mind Method** | CODEX method, framework explainers, vault tours | "Here's how AI and I manage..." |

---

## File Structure

```
~/Documents/onemind-video/
├── src/scripts/
│   ├── generate-script.ts          ← AI script generator
│   ├── generate-legacy-voice.ts    ← OpenAI TTS for Legacy lines
│   ├── generate-youtube-metadata.ts ← SEO title/desc/tags
│   └── generate-social-captions.ts  ← Platform-specific captions
├── scripts/                         ← Generated scripts (daily)
│   └── YYYY-MM-DD-topic.md
├── output/                          ← Legacy voice files (daily)
│   └── YYYY-MM-DD/
│       └── legacy-*.mp3
├── raw/                             ← Raw recordings (daily)
│   └── YYYY-MM-DD/
│       ├── facecam-*.mov
│       └── screen-*.mov
└── templates/
    ├── script-template.md           ← Reusable script structure
    └── thumbnail-prompts.md         ← Thumbnail generation prompts
```

---

*Part of the One Mind framework. Domain: GE 78 Marketing (Content & Growth).*


---

## Related

- [[75-99 GE (Generational Empire)/75 Zeus Delacruz (Brand & Empire)/75 Zeus Delacruz|75 Zeus Delacruz (Brand & Empire)]]
