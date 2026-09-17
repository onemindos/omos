---
provider: YOLOv8 / Ultralytics
category: Computer Vision AI
website: https://yolov8.com/
docs: https://docs.ultralytics.com/
github: https://github.com/ultralytics/ultralytics
status: active
added: 2026-05-24
tags: [computer-vision, object-detection, image-segmentation, classification, yolo, ultralytics, real-time, open-source]
---

# YOLOv8 / Ultralytics — State-of-the-Art Computer Vision

## Overview

YOLOv8 (You Only Look Once v8) by Ultralytics is the latest generation of real-time computer vision models. It's an open-source, state-of-the-art solution for object detection, image classification, instance segmentation, pose estimation, and more.

**Latest:** YOLO26 is the newest model (announced 2026) — next-generation vision AI.

## Key Features

- **Real-Time Object Detection** — Detect and classify objects in images/video instantly
- **Image Segmentation** — Instance segmentation for precise object boundaries
- **Classification** — Image classification tasks
- **Pose Estimation** — Human pose and keypoint detection
- **Tracking** — Multi-object tracking across frames
- **Python API** — Easy integration via pip install
- **CLI Tools** — Command-line interface for quick tasks
- **Model Training** — Train custom models on your datasets
- **Export Formats** — Deploy to ONNX, TensorRT, CoreML, TFLite, Edge TPU
- **Ultralytics Platform** — Cloud training, deployment, and dataset management

## Capabilities

| Task | Description |
|---|---|
| Object Detection | Bounding boxes around detected objects |
| Instance Segmentation | Pixel-perfect object masks |
| Classification | Image category prediction |
| Pose Estimation | Human skeleton/keypoint detection |
| Oriented Bounding Boxes (OBB) | Rotated object detection |
| Tracking | Multi-object tracking (ByteTrack, BoT-SORT) |

## Use Cases

1. **Drone/IoT Fleet Computer Vision** — Real-time object detection on Cognitum Seed devices
2. **Security/Surveillance** — Person detection, intrusion alerts
3. **Automation** — Quality control, defect detection
4. **AR/Robotics** — Spatial understanding, object interaction
5. **Agriculture** — Crop monitoring, pest detection
6. **Traffic/Transport** — Vehicle detection, license plate recognition
7. **Retail** — Customer analytics, inventory tracking

## Installation

```bash
pip install ultralytics
```

## Basic Usage

```python
from ultralytics import YOLO

# Load pretrained model
model = YOLO('yolov8n.pt')  # nano (fastest), s, m, l, x (most accurate)

# Run inference
results = model('path/to/image.jpg')

# Results
results[0].boxes  # bounding boxes
results[0].masks  # segmentation masks
results[0].keypoints  # pose keypoints
```

## Python API

```python
# Train custom model
model.train(data='coco128.yaml', epochs=100, imgsz=640)

# Validate
metrics = model.val()

# Export
model.export(format='onnx')  # or 'torchscript', 'tflite', 'coreml', 'engine'
```

## CLI Usage

```bash
# Predict
yolo predict model=yolov8n.pt source='image.jpg'

# Train
yolo train data=coco128.yaml model=yolov8n.pt epochs=100

# Validate
yolo val model=yolov8n.pt data=coco128.yaml

# Export
yolo export model=yolov8n.pt format=onnx
```

## Models Available

| Model | Size | Speed | Accuracy |
|---|---|---|---|
| YOLOv8n | Nano | Fastest | Good |
| YOLOv8s | Small | Fast | Better |
| YOLOv8m | Medium | Balanced | Great |
| YOLOv8l | Large | Slower | Excellent |
| YOLOv8x | Extra Large | Slowest | Best |
| YOLO26 | Latest | — | State-of-the-art |

## Datasets Supported

- COCO (Common Objects in Context)
- Pascal VOC
- ImageNet
- Custom datasets (YAML config)

## Integration with OneMind OS

### Priority Use Cases

1. **Cognitum Seed Fleet** — Deploy YOLOv8 on IoT devices for real-time vision
   - Person detection (home security)
   - Object tracking (robotics navigation)
   - Pose estimation (gesture control)

2. **Drone Vision** — Aerial object detection
   - Infrastructure inspection
   - Search and rescue
   - Agriculture monitoring

3. **Home Assistant Integration** — Smart home vision
   - Person detection triggers
   - Package delivery detection
   - Pet monitoring

4. **OpenClaw Canvas** — Visual feedback
   - Real-time detection overlays on dashboard
   - Camera feed analysis

### Deployment Strategy

```yaml
# OneMind YOLO deployment plan
edge_devices:
  - Cognitum Seed (ONNX export for efficiency)
  - Raspberry Pi (TFLite)
  - Jetson Nano (TensorRT)

cloud_training:
  - Ultralytics Platform (managed)
  - Modal (serverless GPU)
  - Local MBP (small datasets only)

models_to_train:
  - Person detection (Zeus at home/office)
  - Device recognition (identify IoT devices)
  - Gesture control (hand poses for UI)
```

## Integration Status

- [x] Open-source library evaluated
- [ ] Install locally for testing
- [ ] Train custom model (OneMind-specific objects)
- [ ] Deploy to Cognitum Seed device
- [ ] Integrate with Home Assistant
- [ ] Add to OpenClaw Canvas dashboard
- [ ] API integration for automated workflows

## Pricing & Deployment

### Open Source (Free)
- YOLOv8 is **completely free** — Apache-2.0 license
- Install via: `pip install ultralytics`
- Run locally on CPU, GPU, or edge devices
- No API keys, no usage limits, full source access

### Ultralytics Platform (Paid)
- Optional cloud platform for teams
- Dataset management, model training, deployment
- Web UI + REST API
- **Check pricing:** https://www.ultralytics.com/pricing

### Deployment Options
- **Local:** Python, CLI, Docker
- **Edge:** NVIDIA Jetson, Raspberry Pi, Coral TPU
- **Mobile:** iOS (CoreML), Android (TFLite), ONNX
- **Web:** ONNX.js, TensorFlow.js
- **Cloud:** Any cloud provider (AWS, GCP, Azure)
- **Embedded:** TensorRT, OpenVINO, NCNN

No vendor lock-in — sovereign-friendly ✅
**Ultralytics Platform:** Cloud service — research pricing at https://platform.ultralytics.com/

## Performance

- **Speed:** Real-time (30+ FPS on modern GPUs)
- **Accuracy:** mAP 50-95 up to 53.9% (YOLO8x on COCO)
- **Size:** Models range from 3MB (nano) to 136MB (x-large)

## Community & Support

- **GitHub:** 130.0k+ stars
- **Documentation:** Comprehensive (docs.ultralytics.com)
- **Discord/Forum:** Active community
- **Integrations:** 100+ third-party tools
- **Languages:** Python, CLI, REST API

## Related Technologies

- OpenCV (image processing)
- PyTorch (model backend)
- ONNX Runtime (deployment)
- TensorRT (NVIDIA acceleration)
- CoreML (Apple Silicon optimization)

## Competitors

- **Detectron2** (Facebook/Meta)
- **MMDetection** (OpenMMLab)
- **TensorFlow Object Detection API**
- **Roboflow** (platform/tooling)

## Notes

- Latest YOLO version as of 2026 is **YOLO26**
- Ultralytics maintains the official implementation
- Best choice for production real-time vision tasks
- Excellent for edge deployment (phones, IoT, drones)
- Can train custom models with as little as 100 images

## TODO

- [ ] Benchmark YOLOv8 vs YOLO26 on MBP M2 Max
- [ ] Create OneMind custom dataset (Zeus, devices, spaces)
- [ ] Train person detection model (Zeus-specific)
- [ ] Export to ONNX for Cognitum Seed
- [ ] Integrate with Home Assistant automation
- [ ] Build OpenClaw Canvas real-time detection panel
- [ ] Document deployment workflow as skill

## References

- Official Site: https://yolov8.com/
- Docs: https://docs.ultralytics.com/
- GitHub: https://github.com/ultralytics/ultralytics
- Platform: https://platform.ultralytics.com/
- Paper: https://arxiv.org/abs/2305.09972 (YOLOv8 tech report)
