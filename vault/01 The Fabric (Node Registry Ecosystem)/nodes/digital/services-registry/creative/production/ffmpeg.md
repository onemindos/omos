---
title: "FFmpeg"
type: reference
status: active
tags: [platform, creative, production, media, video, audio, cli]
---
# FFmpeg

> Universal media processing CLI. Convert, transcode, stream, record, filter any audio/video format.

## Common Commands
```bash
# Convert format
ffmpeg -i input.mov output.mp4

# Extract audio
ffmpeg -i video.mp4 -vn -acodec copy audio.aac

# Resize video
ffmpeg -i input.mp4 -vf scale=1920:1080 output.mp4

# Trim video
ffmpeg -i input.mp4 -ss 00:01:00 -t 00:00:30 output.mp4

# GIF from video
ffmpeg -i input.mp4 -vf "fps=10,scale=480:-1" output.gif

# Concatenate files
ffmpeg -f concat -i files.txt -c copy output.mp4

# Stream RTSP to file
ffmpeg -i rtsp://camera/stream -c copy recording.mp4

# Extract frames
ffmpeg -i video.mp4 -vf "fps=1" frame_%04d.png
```

## Key Features
| Feature | Notes |
|---------|-------|
| Formats | Every format (MP4, MKV, AVI, MOV, WebM, FLAC, etc.) |
| Codecs | H.264, H.265, VP9, AV1, ProRes, DNxHD |
| Filters | Scale, crop, overlay, color, denoise, stabilize |
| Streaming | RTMP, RTSP, HLS, DASH |
| Hardware accel | VideoToolbox (Mac), NVENC (NVIDIA), QSV (Intel) |

## Integration
- Surveillance: RTSP → recording/clips
- OBS: output processing
- Automation: batch transcode via scripts
- Drone: video processing pipeline
