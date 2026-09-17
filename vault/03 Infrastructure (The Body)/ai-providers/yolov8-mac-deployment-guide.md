---
type: deployment-guide
provider: yolov8-ultralytics
platform: macOS-m2-max
status: planned
created: 2026-05-24
updated: 2026-05-24
tags:
  - yolov8
  - computer-vision
  - edge-ai
  - mac-m2
  - sovereign
  - onemind
---

# YOLOv8 on Mac M2 Max — OneMind Deployment Guide

## Overview

YOLOv8 running on Mac M2 Max with Neural Engine acceleration for real-time object detection. All on-device, zero cloud, fully sovereign.

## Current Status

- ❌ Not installed yet
- ✅ Mac M2 Max hardware ready (Neural Engine + Metal GPU)
- ✅ Python 3.9.6 available (`/usr/bin/python3`)
- 📝 Use cases documented below

## Hardware Advantages

### M2 Max Neural Engine
- Hardware-accelerated inference (faster than CPU)
- Low power consumption
- Real-time performance (30-100 FPS depending on model)

### Metal GPU
- Native Apple Silicon optimization
- Parallel processing for batch inference
- Energy efficient

### Privacy & Sovereignty
- **All processing on-device** ✅
- No cloud APIs required
- No data leaves the Mac
- Zero vendor lock-in

## Performance Estimates (M2 Max)

| Model | Speed | Use Case |
|-------|-------|----------|
| YOLOv8n (nano) | ~100 FPS | Edge devices, drones, IoT |
| YOLOv8s (small) | ~60 FPS | Balanced speed/accuracy |
| YOLOv8m (medium) | ~30 FPS | High accuracy needed |
| YOLOv8x (extra) | ~15 FPS | Maximum accuracy |

## OneMind OS Use Cases

### 1. Home Security & Monitoring 🏠
- Real-time person/animal detection from security cameras
- Alert when unknown person enters property
- Track package deliveries (detect when packages appear/disappear)
- Pet activity monitoring (where's the dog/cat?)
- **Integration:** Home Assistant + Frigate NVR + YOLOv8

### 2. Drone Fleet Vision 🚁
- Object detection for autonomous navigation
- Landing pad recognition
- Obstacle avoidance
- Aerial inspection (roof damage, crop health)
- Wildlife counting/tracking
- **Integration:** NATS → YOLOv8 → Flight controller

### 3. IoT Seed Device Fleet 🤖
- Edge vision processing on Cognitum Seed devices
- Smart doorbell (face recognition, package detection)
- Room occupancy detection for automation
- Gesture recognition for touchless control
- **Integration:** Seed → MQTT → NATS → YOLOv8 pipeline

### 4. Desktop Automation 💻
- Screen element detection (find buttons, icons, text in screenshots)
- Document analysis (detect tables, figures, diagrams in PDFs)
- Auto-sort photos by content (beach, mountains, people, food)
- Meme detection for social media workflows
- **Integration:** Hermes skill → YOLOv8 → file operations

### 5. OpenClaw Canvas Integration 🖥️
- Live camera feed overlays with detection boxes
- Real-time object tracking dashboard
- Visual alerts in command center
- **Integration:** Canvas WebSocket → YOLOv8 → live overlay

### 6. Telegram/Communication 📱
- Send photo → get detection results
- Auto-tag media in group chats
- Alert bot when specific objects detected
- **Integration:** Telegram MCP → YOLOv8 → reply with annotations

## Installation (When Ready)

```bash
# Install ultralytics package
pip3 install ultralytics

# Quick test with sample image
yolo detect predict model=yolov8n.pt source='https://ultralytics.com/images/bus.jpg'

# Or via Python
python3 -c "from ultralytics import YOLO; model = YOLO('yolov8n.pt'); model.predict('bus.jpg')"
```

### First-Time Setup
```bash
# Create venv (recommended)
python3 -m venv ~/.venvs/yolo
source ~/.venvs/yolo/bin/activate

# Install with dependencies
pip install ultralytics opencv-python pillow

# Download pre-trained model (happens automatically on first run)
yolo detect predict model=yolov8n.pt source='test.jpg'
```

## Integration Ideas for OneMind

### 1. Hermes Skill: `yolo-detect`
```bash
# CLI wrapper for common detection tasks
hermes skill create yolo-detect

# Usage examples:
hermes run yolo-detect --image ~/Desktop/photo.jpg
hermes run yolo-detect --camera 0 --live
hermes run yolo-detect --video ~/Videos/security.mp4 --save-annotations
```

### 2. NATS Event Bridge
```
Camera → YOLOv8 → NATS events → OpenClaw alerts
```

**Flow:**
1. Camera captures frame (Home Assistant, Frigate, or direct)
2. YOLOv8 analyzes frame
3. If "person" detected → publish to `events.vision.person_detected`
4. OpenClaw subscribes → triggers lights/alarms/notifications

**NATS subjects:**
- `vision.detect.request` — send image for detection
- `vision.detect.response.<request_id>` — results
- `events.vision.person_detected` — person alert
- `events.vision.package_detected` — package alert
- `events.vision.pet_detected` — pet alert

### 3. Home Assistant Integration
```yaml
# Example automation
automation:
  - alias: "YOLOv8 Person Detection Alert"
    trigger:
      platform: mqtt
      topic: "vision/person_detected"
    action:
      - service: notify.telegram
        data:
          message: "Person detected at front door"
      - service: light.turn_on
        entity_id: light.porch
```

### 4. Cron Job: Nightly Security Analysis
```bash
# Analyze overnight security footage
# ~/.hermes/scripts/nightly-vision-analysis.sh

#!/bin/bash
# Process all recordings from last 24h
find /security/recordings -name "*.mp4" -mtime -1 | while read video; do
    yolo detect predict model=yolov8m.pt source="$video" save=True
done

# Send summary to Telegram
hermes send telegram "Nightly vision analysis complete. Check ~/yolo/runs/detect/"
```

**Cron schedule:**
```bash
hermes cron create \
  --schedule "0 6 * * *" \
  --script ~/.hermes/scripts/nightly-vision-analysis.sh \
  --name "nightly-yolo-security"
```

### 5. Telegram Bot Command
```python
# /detect command in Telegram
# User sends photo → bot replies with annotated image + detection list

from ultralytics import YOLO

model = YOLO('yolov8n.pt')

# On photo received:
results = model.predict(photo_path)
annotated = results[0].plot()  # Draw boxes
summary = results[0].summary()  # Get detections

# Reply with annotated image + text
telegram.send_photo(annotated)
telegram.send_message(f"Detected: {summary}")
```

## Python API Examples

### Basic Detection
```python
from ultralytics import YOLO

# Load model
model = YOLO('yolov8n.pt')

# Predict
results = model.predict('image.jpg')

# Get detections
for r in results:
    boxes = r.boxes  # Bounding boxes
    for box in boxes:
        cls = int(box.cls[0])  # Class ID
        conf = float(box.conf[0])  # Confidence
        coords = box.xyxy[0].tolist()  # [x1, y1, x2, y2]
        label = model.names[cls]  # Class name
        print(f"{label}: {conf:.2f} at {coords}")
```

### Real-Time Webcam
```python
from ultralytics import YOLO
import cv2

model = YOLO('yolov8n.pt')
cap = cv2.VideoCapture(0)  # Mac FaceTime camera

while True:
    ret, frame = cap.read()
    if not ret:
        break
    
    results = model.predict(frame, stream=True)
    
    for r in results:
        annotated_frame = r.plot()
        cv2.imshow('YOLOv8', annotated_frame)
    
    if cv2.waitKey(1) & 0xFF == ord('q'):
        break

cap.release()
cv2.destroyAllWindows()
```

### Filter Specific Classes
```python
# Only detect people (class 0)
results = model.predict('image.jpg', classes=[0])

# Detect people and pets (dogs=16, cats=15, person=0)
results = model.predict('image.jpg', classes=[0, 15, 16])

# All 80 COCO classes available
print(model.names)  # {0: 'person', 1: 'bicycle', ...}
```

### Save Annotated Images
```python
results = model.predict('image.jpg', save=True, project='runs/detect', name='experiment')
# Saves to: runs/detect/experiment/image.jpg
```

## COCO Classes (80 total)

Most useful for OneMind:
- **0:** person
- **2:** car
- **15:** cat
- **16:** dog
- **24:** backpack
- **26:** handbag
- **28:** suitcase
- **39:** bottle
- **41:** cup
- **56:** chair
- **57:** couch
- **62:** tv
- **63:** laptop
- **64:** mouse
- **65:** remote
- **66:** keyboard
- **67:** cell phone

Full list: https://github.com/ultralytics/ultralytics/blob/main/ultralytics/cfg/datasets/coco.yaml

## Model Export Options

YOLOv8 can export to multiple formats for deployment:

```bash
# Export to CoreML (iOS)
yolo export model=yolov8n.pt format=coreml

# Export to TensorFlow Lite (Android, edge)
yolo export model=yolov8n.pt format=tflite

# Export to ONNX (universal)
yolo export model=yolov8n.pt format=onnx

# Export to TensorRT (NVIDIA)
yolo export model=yolov8n.pt format=engine
```

**For Seed devices:** Use TFLite or ONNX depending on hardware

## Resource Requirements

### Mac M2 Max
- **RAM:** 2-4 GB for YOLOv8n, 4-8 GB for YOLOv8x
- **Storage:** ~50 MB per model
- **Power:** Minimal (Neural Engine is efficient)

### Recommended Model per Use Case
- **Real-time webcam:** YOLOv8n (nano)
- **Security analysis:** YOLOv8s or YOLOv8m
- **Drone navigation:** YOLOv8n (edge deployment)
- **Document analysis:** YOLOv8m or YOLOv8x
- **Photo sorting:** YOLOv8s (batch processing)

## Next Steps (When Ready to Deploy)

### Phase 1: Testing
- [ ] Install ultralytics package
- [ ] Test with sample images
- [ ] Benchmark performance on M2 Max
- [ ] Test real-time webcam detection

### Phase 2: Hermes Integration
- [ ] Create `yolo-detect` skill
- [ ] Add CLI wrappers for common tasks
- [ ] Test integration with file operations

### Phase 3: NATS Bridge
- [ ] Design NATS subject schema
- [ ] Build vision detection service
- [ ] Subscribe OpenClaw to vision events

### Phase 4: Home Automation
- [ ] Integrate with Home Assistant
- [ ] Connect to Frigate NVR (if using)
- [ ] Set up person/package detection alerts

### Phase 5: Edge Deployment
- [ ] Export models for Seed devices
- [ ] Test on IoT fleet
- [ ] Build MQTT → NATS → YOLOv8 pipeline

## Related Files

- Provider doc: `ai-providers/yolov8-ultralytics.md`
- Home Assistant: `~/Documents/onemind/50-74 LE (Legacy Estate)/50 Home Automation (The Domain)/`
- IoT fleet: `~/Documents/onemind/00-24 SO (Sovereign Operations)/07 IoT Fleet (Sentinels)/`
- NATS config: `~/.nats/`

## Resources

- **Ultralytics Docs:** https://docs.ultralytics.com/
- **GitHub:** https://github.com/ultralytics/ultralytics
- **Model Zoo:** https://github.com/ultralytics/assets/releases
- **COCO Dataset:** https://cocodataset.org/
- **Community:** https://community.ultralytics.com/

---

**Status:** Documented and ready for deployment when needed ✅  
**Sovereignty:** 100% on-device, zero cloud ⭐  
**Next Action:** Install when vision use case becomes active priority
