# Whisper

> **OpenAI's speech-to-text for voice transcription**

## Overview

Whisper is OpenAI's automatic speech recognition (ASR) system. Provides highly accurate transcription across multiple languages, enabling Legacy AI to understand voice input.

## Status

| Attribute   | Value             |
| ----------- | ----------------- |
| Status      | ✅ Active         |
| Category    | AI Agents         |
| Type        | Speech-to-Text    |
| Integration | OpenAI API, Local |

## Use Cases in OneMind

1. **Voice Commands** - Transcribe voice input for Legacy AI
2. **Meeting Notes** - Transcribe recordings
3. **Voice Memos** - Convert audio notes to text
4. **Accessibility** - Enable voice-first interactions

## Model Variants

| Model     | Size | Speed   | Accuracy | Use Case          |
| --------- | ---- | ------- | -------- | ----------------- |
| whisper-1 | API  | Fast    | High     | Production        |
| tiny      | 39M  | Fastest | Fair     | Real-time         |
| base      | 74M  | Fast    | Good     | Quick transcripts |
| small     | 244M | Medium  | Better   | Balanced          |
| medium    | 769M | Slow    | Great    | Quality focus     |
| large-v3  | 1.5B | Slowest | Best     | Maximum accuracy  |

## Installation

### OpenAI API (Recommended)

```bash
pip install openai
```

### Local (faster-whisper)

```bash
pip install faster-whisper
```

### Local (Original)

```bash
pip install openai-whisper
```

## API Integration

### OpenAI API

```python
from openai import OpenAI
from pathlib import Path

client = OpenAI()

async def transcribe_audio(audio_path: str, language: str = None) -> str:
    """Transcribe audio file using OpenAI Whisper API."""
    with open(audio_path, "rb") as audio_file:
        transcript = client.audio.transcriptions.create(
            model="whisper-1",
            file=audio_file,
            language=language,  # Optional: "en", "es", etc.
            response_format="text"
        )
    return transcript

async def transcribe_with_timestamps(audio_path: str) -> dict:
    """Transcribe with word-level timestamps."""
    with open(audio_path, "rb") as audio_file:
        transcript = client.audio.transcriptions.create(
            model="whisper-1",
            file=audio_file,
            response_format="verbose_json",
            timestamp_granularities=["word"]
        )
    return transcript
```

### Local with faster-whisper

```python
from faster_whisper import WhisperModel

# Load model (runs on GPU if available)
model = WhisperModel("large-v3", device="cuda", compute_type="float16")
# Or CPU: WhisperModel("large-v3", device="cpu", compute_type="int8")

def transcribe_local(audio_path: str) -> str:
    """Transcribe using local Whisper model."""
    segments, info = model.transcribe(
        audio_path,
        beam_size=5,
        language="en",
        vad_filter=True  # Filter silence
    )

    text = " ".join([segment.text for segment in segments])
    return text.strip()

def transcribe_with_segments(audio_path: str) -> list:
    """Get transcription with time segments."""
    segments, info = model.transcribe(audio_path)

    result = []
    for segment in segments:
        result.append({
            "start": segment.start,
            "end": segment.end,
            "text": segment.text.strip()
        })

    return result
```

### Real-time Streaming

```python
import pyaudio
import numpy as np
from faster_whisper import WhisperModel

model = WhisperModel("base", device="cpu")

def stream_transcribe():
    """Real-time microphone transcription."""
    CHUNK = 1024
    FORMAT = pyaudio.paFloat32
    CHANNELS = 1
    RATE = 16000

    p = pyaudio.PyAudio()
    stream = p.open(
        format=FORMAT,
        channels=CHANNELS,
        rate=RATE,
        input=True,
        frames_per_buffer=CHUNK
    )

    print("Listening... (Ctrl+C to stop)")

    audio_buffer = []

    try:
        while True:
            data = stream.read(CHUNK)
            audio_buffer.append(np.frombuffer(data, dtype=np.float32))

            # Process every 3 seconds
            if len(audio_buffer) * CHUNK / RATE >= 3:
                audio_data = np.concatenate(audio_buffer)
                segments, _ = model.transcribe(audio_data)

                for segment in segments:
                    print(f">> {segment.text}")

                audio_buffer = []

    except KeyboardInterrupt:
        stream.close()
        p.terminate()
```

## Legacy AI Integration

```python
from agno import Tool
import tempfile

@Tool
async def transcribe_voice(audio_path: str) -> str:
    """Transcribe audio file to text.

    Args:
        audio_path: Path to audio file (mp3, wav, m4a, etc.)
    """
    return await transcribe_audio(audio_path)

@Tool
async def transcribe_meeting(
    audio_path: str,
    include_timestamps: bool = True
) -> dict:
    """Transcribe a meeting recording with timestamps.

    Args:
        audio_path: Path to meeting audio
        include_timestamps: Whether to include time markers
    """
    if include_timestamps:
        result = await transcribe_with_timestamps(audio_path)
        return {
            "text": result.text,
            "segments": result.segments,
            "duration": result.duration
        }
    else:
        text = await transcribe_audio(audio_path)
        return {"text": text}

@Tool
async def voice_to_command(audio_data: bytes) -> str:
    """Convert voice input to text command.

    Args:
        audio_data: Raw audio bytes from microphone
    """
    # Save to temp file
    with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as f:
        f.write(audio_data)
        temp_path = f.name

    text = await transcribe_audio(temp_path)
    return text.strip()
```

## LiveKit Integration

```python
from livekit import rtc

async def transcribe_livekit_audio(track: rtc.AudioTrack):
    """Transcribe audio from LiveKit track."""
    audio_buffer = []

    async for frame in track:
        audio_buffer.append(frame.data)

        # Process every 5 seconds
        if len(audio_buffer) * frame.duration >= 5.0:
            audio_data = b"".join(audio_buffer)

            # Save and transcribe
            with tempfile.NamedTemporaryFile(suffix=".wav") as f:
                # Write WAV header + data
                f.write(audio_data)
                f.flush()

                text = await transcribe_audio(f.name)
                yield text

            audio_buffer = []
```

## Docker Deployment

```yaml
version: "3.8"

services:
  whisper:
    image: onerahmet/openai-whisper-asr-webservice:latest
    ports:
      - "9000:9000"
    environment:
      - ASR_MODEL=large-v3
      - ASR_ENGINE=faster_whisper
    volumes:
      - whisper_cache:/root/.cache/huggingface
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]

volumes:
  whisper_cache:
```

### Whisper API Service

```python
import httpx

WHISPER_SERVICE_URL = "http://localhost:9000"

async def transcribe_via_service(audio_path: str) -> str:
    """Transcribe using self-hosted Whisper service."""
    async with httpx.AsyncClient() as client:
        with open(audio_path, "rb") as f:
            response = await client.post(
                f"{WHISPER_SERVICE_URL}/asr",
                files={"audio_file": f},
                params={
                    "task": "transcribe",
                    "language": "en",
                    "output": "txt"
                }
            )
        return response.text
```

## Environment Variables

```bash
OPENAI_API_KEY=your_key

# For local deployment
WHISPER_MODEL=large-v3
WHISPER_DEVICE=cuda  # or cpu
WHISPER_SERVICE_URL=http://localhost:9000
```

## Performance Tips

1. **Use faster-whisper** - 4x faster than original
2. **GPU acceleration** - CUDA for speed
3. **VAD filtering** - Skip silence
4. **Batch processing** - Process multiple files together
5. **Model selection** - Balance speed vs accuracy

## Supported Audio Formats

- MP3, MP4, MPEG
- WAV, WEBM
- M4A, OGG, FLAC

## Best Practices

1. **Clean audio** - Better quality = better transcription
2. **Appropriate model** - Don't use large-v3 for real-time
3. **Language hints** - Specify language when known
4. **Error handling** - Audio files may be corrupted
5. **Privacy** - Consider local deployment for sensitive audio

## Resources

- [OpenAI Whisper](https://openai.com/research/whisper)
- [faster-whisper](https://github.com/guillaumekln/faster-whisper)
- [Whisper GitHub](https://github.com/openai/whisper)
- [Whisper API Docs](https://platform.openai.com/docs/guides/speech-to-text)

---

Last Updated: January 25, 2026


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
