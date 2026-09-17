---
title: "FFmpeg — Video Processing CLI"
type: reference
domain: SO
status: active
created: 2026-05-26
updated: 2026-05-26
tags: [video, tools, cli, content, open-source, ffmpeg]
---

# FFmpeg

**URL:** https://ffmpeg.org
**GitHub:** https://github.com/FFmpeg/FFmpeg (46k+ stars)
**License:** LGPL / GPL (open source, free)
**Status on Zeus's Mac:** INSTALLED — v8.1.1 at `/opt/homebrew/bin/ffmpeg`

---

## What It Is

The universal video/audio processing engine. Open source since 2000. Runs under the hood of VLC, YouTube, OBS, Handbrake, and nearly every video tool in existence. Command-line only — no GUI. Runs on Mac via Homebrew.

Zeus's M2 Max build includes hardware acceleration (VideoToolbox + AudioToolbox + NEON) — video processing is fast, uses the GPU/Neural Engine, not just CPU.

---

## What It Can Do

| Task | Command Pattern |
|------|----------------|
| Trim a clip | `ffmpeg -i input.mp4 -ss 00:01:00 -to 00:02:30 -c copy clip.mp4` |
| Convert format | `ffmpeg -i input.mov output.mp4` |
| Resize to vertical (TikTok/Reels 9:16) | `ffmpeg -i input.mp4 -vf "scale=1080:1920,setsar=1" output.mp4` |
| Extract audio | `ffmpeg -i input.mp4 -vn -acodec mp3 audio.mp3` |
| Burn in subtitles | `ffmpeg -i input.mp4 -vf subtitles=subs.srt output.mp4` |
| Merge clips | `ffmpeg -f concat -safe 0 -i list.txt -c copy merged.mp4` |
| Speed up video | `ffmpeg -i input.mp4 -vf "setpts=0.5*PTS" -af "atempo=2.0" fast.mp4` |
| Add watermark/logo | `ffmpeg -i input.mp4 -i logo.png -filter_complex "overlay=10:10" output.mp4` |
| Batch resize folder | `for f in *.mp4; do ffmpeg -i "$f" -vf scale=1080:1920 "out_$f"; done` |
| Extract thumbnail | `ffmpeg -i input.mp4 -ss 00:00:05 -vframes 1 thumb.jpg` |

---

## Content Pipeline Use Cases (Legacy Platform Holdings)

FFmpeg is the backbone for any automated content pipeline:

1. **Long → Short clips** — trim timestamps from transcripts, output clips
2. **Reformat for platform** — landscape 16:9 → vertical 9:16 for TikTok/Reels
3. **Burn captions** — combine Whisper `.srt` output + FFmpeg to make captioned clips
4. **Batch process** — one shell loop processes 100 clips overnight
5. **Audio extraction** — pull audio from recordings for podcast publishing

---

## Integration with OneMind OS Stack

```
Video file
  → Whisper (transcribe + generate .srt)
  → auto-editor (remove silence/dead air)
  → FFmpeg (resize, burn captions, format for platform)
  → GeeLark API / quso.ai (schedule + distribute)
```

This pipeline runs 100% on Zeus's Mac, zero cloud cost.

---

## Install (already installed)

```bash
# Already installed at /opt/homebrew/bin/ffmpeg (v8.1.1)
# To reinstall or upgrade:
brew install ffmpeg
brew upgrade ffmpeg
```

---

## Key Flags Reference

| Flag | Meaning |
|------|---------|
| `-i` | Input file |
| `-ss` | Start time (seek to) |
| `-to` | End time |
| `-t` | Duration |
| `-c copy` | Copy streams without re-encoding (fast) |
| `-vf` | Video filter |
| `-af` | Audio filter |
| `-vn` | No video (audio only) |
| `-an` | No audio (video only) |
| `-y` | Overwrite output without asking |

---

## Related Tools in Stack

- **auto-editor** — sits before FFmpeg, removes silence/dead air automatically → https://github.com/WyattBlue/auto-editor
- **Whisper** — generates `.srt` subtitle files that FFmpeg burns into video
- **quso.ai** — SaaS layer on top if manual pipeline is too much work
