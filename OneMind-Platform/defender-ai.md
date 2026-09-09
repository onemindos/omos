# Defender AI — Reference Platform

## What It Is
360° live video streaming platform with AI-powered object detection, geo-registration, and slew-to-cue tracking. Targets military, law enforcement, and first responders. Hardware form factors: body-worn cameras, vehicle-mounted systems, airborne pods.

## What It Does
- Ingests live video from multiple camera sources (RTSP/LTE/5G)
- Runs real-time AI inference on the edge (NVIDIA GPU) detecting people, vehicles, aircraft, threats
- Geo-registers every detected object — maps pixel coordinates to lat/lon on the ground
- Outputs a live stream of tracked objects with position, classification, and confidence
- Slew-to-cue: operators can click a detected object and a physical camera slews to track it
- Integrates with C2 platforms (their term — this is what OneMind is)

## Why We're Building Our Own
- Defender is a closed commercial product — no source access, no customization
- Their output is essentially a CoT-compatible detection stream — we can build this natively on NATS
- We own the inference pipeline end-to-end — model selection, edge deployment, subject naming
- Our version feeds directly into ORP for fusion with the full SA picture

---

## OneMind Vision Layer — `omos-vision`

### What We Build
A Python microservice that replicates Defender AI's core capability on our NATS bus.

### Architecture
```
RTSP / HLS camera streams (IP cameras, drones, body-worn)
  → omos-vision service
      → frame extraction (OpenCV / FFmpeg)
      → AI inference (YOLO v11 / ONNX Runtime / NVIDIA Triton)
      → geo-registration (camera intrinsics + extrinsics → WGS84)
      → NATS publish: det.visual.{type}.{id}
          → ORP fuses with TAK picture
          → CoT → TAK Server → ATAK / CloudTAK
```

### NATS Subject
```
det.visual.person.{track_id}
det.visual.vehicle.{track_id}
det.visual.aircraft.{track_id}
det.visual.unknown.{track_id}
```

### Payload (GeoJSON Feature)
```json
{
  "type": "Feature",
  "geometry": { "type": "Point", "coordinates": [lon, lat, alt] },
  "properties": {
    "uid": "det-visual-{track_id}",
    "type": "a-h-G",
    "callsign": "VID-TRACK-42",
    "classification": "vehicle",
    "confidence": 0.94,
    "source_camera": "cam-north-01",
    "sensor_footprint": [[lon,lat], ...],
    "timestamp": "2026-08-30T12:00:00Z"
  }
}
```

### Tech Stack
- Python 3.12
- Ultralytics YOLO v11 (or ONNX Runtime for custom models)
- OpenCV for frame extraction
- nats.py for bus publishing
- Optional: NVIDIA Triton Inference Server for GPU clusters

### Sensor Footprint on Map
Camera FOV projected onto terrain → published as `det.visual.footprint.{cam_id}` → rendered in CloudTAK as a polygon overlay via deck.gl or MapLibre fill layer.

### Hardware Support
- Any IP camera (RTSP stream)
- DJI drone video feed
- Body-worn camera (RTSP output)
- Edge deployment: NVIDIA Jetson Orin for on-device inference
