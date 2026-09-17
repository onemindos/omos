---
title: "CodeProject.AI"
type: reference
status: planned
tags: [platform, iot, surveillance, ai, object-detection, self-hosted]
---
# CodeProject.AI

> Local AI object detection server. Drop-in replacement for DeepStack. Runs on CPU/GPU, self-hosted.

## Key Features
| Feature | Notes |
|---------|-------|
| Object detection | Person, vehicle, animal, face |
| Self-hosted | No cloud, no API keys |
| GPU support | NVIDIA CUDA, Intel OpenVINO, Apple M-series |
| Blue Iris native | Built-in integration |
| Frigate compatible | As detection backend |
| Modules | License plate, face recognition, scene classification |

## Deploy
```bash
docker run -d --name codeproject-ai \
  -p 32168:32168 \
  -v codeproject-data:/app/data \
  codeproject/ai-server:latest

# GPU (NVIDIA)
docker run -d --gpus all --name codeproject-ai \
  -p 32168:32168 \
  codeproject/ai-server:cuda12_2
```

## API
```bash
curl -X POST http://localhost:32168/v1/vision/detection \
  -F "image=@photo.jpg"
```

## Integration
- Frigate: Use as detector backend
- Blue Iris: Native AI module
- Home Assistant: Via Frigate or custom REST sensor
