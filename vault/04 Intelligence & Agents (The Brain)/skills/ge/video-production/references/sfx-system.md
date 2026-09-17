# SFX System Reference — Video Production

Sound effects live in `public/audio/sfx/` and are generated via ffmpeg synthesis (no external dependencies).

| File | Size | Purpose | When to Use |
| ---- | ---- | ------- | ----------- |
| `whoosh.mp3` | 4.5 KB | Transition swoosh | Between segments |
| `reveal.mp3` | 2.9 KB | Subtle high ding | Key reveals, new concept introductions |
| `impact.mp3` | 3.7 KB | Bass emphasis hit | Big numbers, dramatic moments |
| `ambient-bg.mp3` | 934 KB | 4-min lo-fi synth pad | Global background, volume 0.08 |

---

## ffmpeg Generation Commands

```bash
cd ~/Documents/onemind-video/public/audio/sfx

# Whoosh — transition swoosh
ffmpeg -f lavfi -i "anoisesrc=d=0.3:c=pink:r=44100" \
  -af "afade=in:0:0.05,afade=out:st=0.15:d=0.15,highpass=f=2000,lowpass=f=8000" \
  -y whoosh.mp3

# Reveal — subtle high ding
ffmpeg -f lavfi -i "sine=frequency=1200:duration=0.2" \
  -af "afade=in:0:0.02,afade=out:st=0.1:d=0.1,aecho=0.6:0.3:20:0.3" \
  -y reveal.mp3

# Impact — bass emphasis hit
ffmpeg -f lavfi -i "sine=frequency=60:duration=0.25" \
  -af "afade=in:0:0.01,afade=out:st=0.1:d=0.15,bass=g=10" \
  -y impact.mp3

# Ambient background — 4-min lo-fi synth pad
ffmpeg -f lavfi -i "sine=frequency=220:duration=240" \
  -f lavfi -i "sine=frequency=330:duration=240" \
  -f lavfi -i "sine=frequency=440:duration=240" \
  -filter_complex "[0:a][1:a][2:a]amix=inputs=3:duration=first,volume=0.15,lowpass=f=800,afade=in:0:2,afade=out:st=238:d=2" \
  -y ambient-bg.mp3
```

> **Future improvement:** Replace `ambient-bg.mp3` with a Suno AI-generated lo-fi track for production quality.

---

## Usage in Composition

```tsx
// Global ambient (entire video)
<Audio src={staticFile("audio/sfx/ambient-bg.mp3")} volume={0.08} />

// Per-segment SFX (at segment start)
<Sequence from={OFF.s03}>
  <Audio src={staticFile("audio/sfx/whoosh.mp3")} volume={0.15} />
</Sequence>
```

> **Note:** `ambient-bg.mp3` is a placeholder ffmpeg synth pad. For production quality, replace with a Suno AI lo-fi track or royalty-free music.
