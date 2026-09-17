---
title: "LPR/ALPR Research — Sovereign License Plate Recognition"
type: research
domain: SO/03 Infrastructure
topic: computer-vision, security, iot
status: reference
created: 2026-05-24
updated: 2026-05-26
tags: [type/research, domain/so, topic/vision, topic/lpr, topic/security]
---

# LPR / ALPR Research — Sovereign License Plate Recognition

**Context:** Flock Safety cameras are law enforcement only — no civilian API access. This doc covers sovereign alternatives for OneMind OS.

---

## Options Compared

| Solution | Accuracy | Speed | Privacy | Cost | Sovereign | Fit |
|----------|----------|-------|---------|------|-----------|-----|
| **OpenALPR** | ⭐⭐⭐⭐⭐ | Fast | ✅ | Free (AGPL) | ✅ | ⭐⭐⭐⭐ |
| **Ultimat ANPR** | ⭐⭐⭐⭐⭐ | Very Fast | ✅ | Free | ✅ | ⭐⭐⭐⭐⭐ |
| **Plate Recognizer** | ⭐⭐⭐⭐⭐ | Fast | ⚠️ Cloud | $50/mo | ❌ | ⭐⭐⭐ |
| **ANPR Pipeline** | ⭐⭐⭐⭐ | Fast | ✅✅ HMAC | Free | ✅ | ⭐⭐⭐⭐ |
| **YOLOv8 + EasyOCR** | ⭐⭐⭐ | Fast (GPU) | ✅ | Free | ✅ | ⭐⭐⭐⭐ |

---

## Option 1 — OpenALPR (Recommended Start)

**GitHub:** https://github.com/openalpr/openalpr
**License:** AGPL-3.0 (personal use fine)
**Status:** Active (2024)

Features: Multi-country (US, EU, Brazil, AU, Korea), 95%+ accuracy, Python SDK, macOS via Homebrew.

```bash
# Install
brew install openalpr

# Test
curl -O https://raw.githubusercontent.com/openalpr/openalpr/master/samples/us-1.jpg
alpr us-1.jpg
```

```python
from openalpr import Alpr
import cv2

alpr = Alpr("us", "/etc/openalpr/openalpr.conf", "/usr/share/openalpr/runtime_data")
alpr.set_top_n(3)

cap = cv2.VideoCapture("rtsp://camera-ip/stream")
ret, frame = cap.read()

results = alpr.recognize_ndarray(frame)
for plate in results['results']:
    print(f"Plate: {plate['plate']}, Confidence: {plate['confidence']}")
```

---

## Option 2 — Ultimat ANPR (Best Performance)

**Tech:** TensorFlow / TensorRT / OpenVINO — NPU/GPU accelerated
**Platforms:** Mac M2 Max (Metal), Pi 5 + AI HAT 2 (Hailo NPU), Jetson
**Charset:** Latin, Korean, Chinese

Best fit for production deployment on edge nodes. Fastest inference of all options.

---

## Option 3 — YOLOv8 + EasyOCR (Most Sovereign)

100% local, no licenses, Metal GPU on M2 Max.

```python
from ultralytics import YOLO
import easyocr
import cv2

vehicle_model = YOLO('yolov8n.pt')
reader = easyocr.Reader(['en'], gpu=True)

frame = cv2.imread('car.jpg')
results = vehicle_model(frame)

for r in results:
    for box in r.boxes:
        if box.cls == 2:  # Car class
            x1, y1, x2, y2 = box.xyxy[0]
            plate_region = frame[int(y1):int(y2), int(x1):int(x2)]
            plate_text = reader.readtext(plate_region)
            if plate_text:
                print(f"Plate: {plate_text[0][1]}")
```

Tradeoff: requires tuning for accuracy. Good starting point since YOLOv8 is already in the stack.

---

## Recommendation

- **Prototype:** OpenALPR — easiest setup, battle-tested
- **Production:** Ultimat ANPR — fastest, NPU-ready for edge nodes
- **Fallback:** YOLOv8 + EasyOCR — if staying fully in-stack

---

## Full Architecture (OneMind LPR System)

```
Camera (RTSP) → Frame capture (Pi Zero or Mac)
  ↓
Vehicle detection (YOLOv8 — Mac M2 Max or Pi 5 + AI HAT 2)
  ↓
Plate detection + OCR (OpenALPR or Ultimat ANPR)
  ↓
NATS pub: events.vision.vehicle_detected
  {
    "plate": "ABC1234",
    "confidence": 0.95,
    "vehicle_type": "car",
    "color": "silver",
    "timestamp": "2026-05-26T14:23:45Z",
    "camera_id": "driveway-cam-01"
  }
  ↓
NATS → HA Bridge → Home Assistant automation
  ↓
Actions:
  - Known plate → Open garage door
  - Unknown plate → Telegram alert
  - Delivery truck → Package arriving notification
  - All vehicles → Log to local SQLite
```

---

## Privacy

### Plate Hashing (GDPR compliant)

```python
import hashlib, hmac

SECRET_KEY = "..."  # Store in 1Password

def hash_plate(plate_text):
    return hmac.new(
        SECRET_KEY.encode(),
        plate_text.encode(),
        hashlib.sha256
    ).hexdigest()

# Store hash, not raw plate text
hashed = hash_plate("ABC1234")
```

Rules:
- Never upload plates to cloud
- Store hashes + timestamps in local SQLite only
- Delete records after 30 days (configurable)
- Known plates: whitelist in Home Assistant

---

## Action Items

- [ ] Install OpenALPR, test with sample US plates
- [ ] Deploy YOLOv8 vehicle detection, chain with OpenALPR
- [ ] Build NATS publisher script (`events.vision.vehicle_detected`)
- [ ] Install IP camera at driveway, configure RTSP stream
- [ ] Build Home Assistant automations (known/unknown plate, delivery)
- [ ] Evaluate Ultimat ANPR for production performance vs OpenALPR

---

## Related

- `SO/03 Infrastructure/hardware/iot-hardware-roadmap.md` — Pi Zero + Pi AI HAT 2 specs
- `SO/03 Infrastructure/vision/obs-automation-research.md` — OBS camera feeds
- Home Assistant: `LE/50 Home Automation/`
