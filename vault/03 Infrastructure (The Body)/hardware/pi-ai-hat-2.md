---
type: hardware
category: edge-ai
vendor: Raspberry Pi Foundation
model: AI HAT 2
status: coming-q2-2026
price: ~$70-100
created: 2026-05-24
tags:
  - raspberry-pi
  - edge-ai
  - neural-accelerator
  - yolov8
  - hailo
---

# Raspberry Pi AI HAT 2 — Edge Vision Accelerator

**Status:** Coming Q2 2026 (check availability)  
**OneMind Role:** Real-time edge vision (cameras, drones, security)

## Specs

| Feature | Value |
|---------|-------|
| **NPU** | Hailo-8L AI accelerator (13 TOPS) |
| **Inference Speed** | 30+ FPS (YOLOv8n at 640×640) |
| **Supported Models** | YOLOv8, MobileNet, ResNet, EfficientDet, custom ONNX |
| **Power** | 3W typical |
| **Compatibility** | **Raspberry Pi 5 ONLY** (requires PCIe interface) |
| **Interface** | PCIe Gen 2.0 (via Pi 5's M.2 slot) |
| **Software** | Hailo SDK, TensorFlow Lite, PyTorch, ONNX Runtime |
| **Price** | ~$70–100 (estimated) |

## Why This Matters for OneMind

### Problem: Pi Alone = Too Slow for Vision AI
| Device | YOLOv8n Speed | Usable? |
|--------|---------------|---------|
| Pi Zero 2 W | ~0.1 FPS | ❌ Unusable |
| Pi 4 (4GB) | ~2 FPS | ❌ Too slow |
| Pi 5 (8GB, CPU only) | ~5 FPS | ⚠️ Barely usable |
| **Pi 5 + AI HAT 2** | **30 FPS** | ✅ Real-time! |
| Mac M2 Max | 100 FPS | ✅ Fast (but not portable) |

### Solution: Edge Vision Nodes
Deploy **Pi 5 + AI HAT 2** at strategic locations for on-device inference with zero latency.

## OneMind Integration Use Cases

### 1. Security Camera + AI
**Deployment:**
```
Pi 5 + AI HAT 2 + Camera Module
→ YOLOv8 (on-device, 30 FPS)
→ Person detection
→ NATS alert to Legacy/Home Assistant
→ Zero cloud, zero Mac dependency
```

**Benefits:**
- **Privacy:** No video sent to Mac (only NATS alerts)
- **Latency:** <100ms detection-to-alert
- **Bandwidth:** No video streaming (only events)
- **Resilience:** Works if Mac is offline

**Architecture:**
```
Camera → Pi 5 (YOLOv8 @ 30 FPS) → Person detected
  ↓
NATS pub: events.vision.person_detected {cam_id, timestamp, bbox, confidence}
  ↓
NATS → HA Bridge → Home Assistant → Telegram notification
```

### 2. Drone Vision (Onboard)
**Deployment:**
```
Drone + Pi 5 + AI HAT 2 + Camera
→ Object avoidance (YOLOv8)
→ Landing pad recognition
→ Gesture recognition
```

**Benefits:**
- **Autonomous:** No Mac/network needed mid-flight
- **Low power:** 3W NPU (battery-friendly)
- **Real-time:** 30 FPS for obstacle detection

### 3. Smart Doorbell
**Deployment:**
```
Pi 5 + AI HAT 2 + Camera + PIR Sensor
→ Face recognition (local dataset)
→ Package detection
→ Visitor classification (known/unknown/delivery)
```

**Architecture:**
```
PIR trigger → Camera → YOLOv8 person + face detection
→ Known face? → "Welcome home, Zeus"
→ Unknown? → "Visitor at door" + snapshot to Telegram
→ Package? → "Package delivered" + log to Home Assistant
```

### 4. Home Assistant AI Node (Room Occupancy)
**Deployment:**
```
Pi 5 + AI HAT 2 + Ceiling-mounted camera
→ Person detection per room
→ NATS → Home Assistant
→ Auto-adjust lights, HVAC
```

**Architecture:**
```
Living room camera → YOLOv8 → 2 people detected
  ↓
NATS: sensors.occupancy.living_room = 2
  ↓
Home Assistant: Turn on lights, set HVAC to "home" mode
```

### 5. Gesture Control
**Deployment:**
```
Pi 5 + AI HAT 2 + Camera (wall-mounted)
→ Hand gesture recognition (custom model)
→ Control lights, music, scenes
```

**Example:**
- Wave hand → Lights on
- Peace sign → Play music
- Thumbs up → Increase brightness

## vs. Mac M2 Max

| Feature | Pi 5 + AI HAT 2 | Mac M2 Max |
|---------|-----------------|------------|
| **Speed** | 30 FPS (YOLOv8n) | 100 FPS (YOLOv8n) |
| **Power** | 3W (NPU only) | 15–30W (whole system) |
| **Cost** | $150 ($80 Pi 5 + $70 HAT) | $3000+ (already owned) |
| **Deployment** | Edge (cameras, drones, remote) | Central hub (high throughput) |
| **Portability** | ✅ Battery-powered | ❌ Desktop only |
| **Sovereignty** | ✅ On-device | ✅ On-device |
| **Best for** | Distributed vision (10+ cameras) | Central processing (aggregate) |

### Hybrid Strategy (Recommended)
- **Mac M2 Max:** Central AI hub for high-throughput tasks
- **Pi 5 + AI HAT 2:** Edge vision nodes (cameras, drones, doorbell)
- **Pi Zero 2 W:** Dumb sensors (temp, motion, GPIO)

**Example Fleet:**
- 1x Mac M2 Max (central hub, NATS server, Home Assistant)
- 5x Pi 5 + AI HAT 2 (cameras: driveway, front door, backyard, living room, drone)
- 10x Pi Zero 2 W (sensors: rooms, garage, gates)

**Total cost:** $150 × 5 + $15 × 10 = **$900** for edge fleet

## Installation & Setup (When Available)

### 1. Order Hardware
```bash
# Components needed:
- Raspberry Pi 5 (8GB recommended): $80
- AI HAT 2: $70–100
- Pi Camera Module v3: $25
- Power supply (5V/5A): $10
- microSD card (64GB): $10

# Total per node: ~$195
```

### 2. Install Raspberry Pi OS (64-bit)
```bash
# Use Raspberry Pi Imager
# Select: Raspberry Pi OS (64-bit)
# Enable SSH, configure Wi-Fi, set hostname: pi5-cam-01
```

### 3. Install Hailo SDK
```bash
ssh onemind@pi5-cam-01.local

# Install Hailo runtime
sudo apt update
sudo apt install hailo-all

# Test NPU
hailortcli fw-control identify
# Should show: Hailo-8L, 13 TOPS
```

### 4. Install YOLOv8 + Hailo Backend
```bash
# Install Ultralytics
pip3 install ultralytics

# Export YOLOv8 for Hailo
from ultralytics import YOLO
model = YOLO('yolov8n.pt')
model.export(format='hailo', imgsz=640)

# This creates: yolov8n.hef (Hailo Executable Format)
```

### 5. Run Inference
```python
# vision_node.py
from hailo_platform import HEF, VDevice, ConfigureParams
import cv2
import json
import subprocess

# Load model
hef = HEF('yolov8n.hef')
target = VDevice()
infer_model = target.create_infer_model(hef)

# Open camera
cap = cv2.VideoCapture(0)

while True:
    ret, frame = cap.read()
    if not ret:
        break
    
    # Preprocess
    input_data = preprocess(frame)  # Resize to 640×640, normalize
    
    # Inference on Hailo NPU (30 FPS!)
    results = infer_model.run(input_data)
    
    # Postprocess
    detections = postprocess(results)  # Parse bounding boxes
    
    # Publish to NATS if person detected
    for det in detections:
        if det['class'] == 'person' and det['confidence'] > 0.5:
            event = {
                "cam_id": "pi5-cam-01",
                "timestamp": time.time(),
                "class": "person",
                "confidence": det['confidence'],
                "bbox": det['bbox']
            }
            subprocess.run([
                "nats", "pub", "events.vision.person_detected",
                json.dumps(event)
            ])
            print(f"Person detected: {event}")
    
    # Optional: Display (for debugging)
    for det in detections:
        draw_box(frame, det)
    cv2.imshow('Vision Node', frame)
    
    if cv2.waitKey(1) & 0xFF == ord('q'):
        break

cap.release()
cv2.destroyAllWindows()
```

### 6. Run as systemd Service
```bash
sudo nano /etc/systemd/system/vision-node.service

# Paste:
[Unit]
Description=OneMind Vision Node
After=network.target

[Service]
ExecStart=/usr/bin/python3 /home/onemind/vision_node.py
Restart=always
User=onemind

[Install]
WantedBy=multi-user.target

# Enable
sudo systemctl enable vision-node
sudo systemctl start vision-node
```

## Benchmark (Expected)

| Model | Resolution | Speed | Power | Latency |
|-------|------------|-------|-------|---------|
| YOLOv8n | 640×640 | 30 FPS | 3W | 33ms |
| YOLOv8s | 640×640 | 20 FPS | 3W | 50ms |
| YOLOv8m | 640×640 | 10 FPS | 3W | 100ms |

**Comparison:**
- **Pi 5 (CPU only):** YOLOv8n @ 5 FPS (6x slower)
- **Mac M2 Max:** YOLOv8n @ 100 FPS (3.3x faster, but not portable)

## Next Steps (When Available — Q2 2026)

- [ ] Check Pi AI HAT 2 availability (watch official store)
- [ ] Order 1x test unit (Pi 5 + AI HAT 2)
- [ ] Benchmark YOLOv8n vs. Mac M2 Max
- [ ] Test NATS pub/sub latency
- [ ] Deploy first camera node (driveway)
- [ ] Build Home Assistant integration
- [ ] Scale to 5 nodes (front door, backyard, living room, drone)

## Related Files
- Pi Zero 2 W (sensors): `hardware/pi-zero-2w.md`
- YOLOv8 Mac deployment: `ai-providers/yolov8-mac-deployment-guide.md`
- Home Assistant: `50-74 LE/50 Home Automation/`
- NATS fabric: `00-24 SO/03 Infrastructure/nats-server.md` (if exists)
