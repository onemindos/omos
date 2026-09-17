# Gesture Systems

## 🔗 Sources

| Source               | Link                                                                                                    | Description                                           |
| -------------------- | ------------------------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| **Instagram Reel**   | [instagram.com/reel/DTF-F-giHSp](https://www.instagram.com/reel/DTF-F-giHSp/?igsh=MWJibGptMHdpM2hwNw==) | Kinect + TouchDesigner + Unreal particle control demo |
| **ADA V2 Gestures**  | [github.com/nazirlouis/ada_v2](https://github.com/nazirlouis/ada_v2)                                    | MediaPipe hand gesture reference code                 |
| **MediaPipe**        | [mediapipe.dev](https://mediapipe.dev/)                                                                 | Google's ML framework for hand/body tracking          |
| **TouchDesigner**    | [derivative.ca](https://derivative.ca/)                                                                 | Visual programming for real-time content              |
| **Azure Kinect SDK** | [Microsoft Docs](https://docs.microsoft.com/en-us/azure/kinect-dk/)                                     | Depth camera body tracking                            |

**Extraction Date:** 2026-01-30

---

## 📋 Overview

This folder contains resources for implementing gesture-based control in OneMind OS, ranging from simple webcam-based approaches to professional depth camera setups.

### The Instagram Reel Concept

> _"I developed a real-time interactive system using Unreal Engine, Kinect, and TouchDesigner, where hand tracking data directly controls and moves particles. By capturing body and hand movement through Kinect and processing it in TouchDesigner, the system drives live particle behavior inside Unreal. This project explores the connection between physical movement and digital space, turning gestures into responsive, living visuals."_

This is the inspiration for the advanced gesture pipeline documented in `kinect_pipeline.md`.

---

## 📁 Contents

| File                     | Purpose                                            | OneMind Use Case                            |
| ------------------------ | -------------------------------------------------- | ------------------------------------------- |
| `kinect_pipeline.md`     | Full Kinect + TouchDesigner + Unreal architecture  | Immersive gesture control for drones/robots |
| `mediapipe_gestures.py`  | Webcam-based gesture detection with classification | Quick-start gesture recognition             |
| `nats_gesture_bridge.py` | Bridge between gesture sources and NATS            | Connect any gesture source to OneMind       |

---

## 🎮 Approach Comparison

```
                    COMPLEXITY
                        ↑
                        │
    ┌───────────────────┼───────────────────┐
    │                   │                   │
    │  Kinect/Azure     │  Professional     │
    │  + TouchDesigner  │  Installation     │
    │  + Unreal Engine  │                   │
    │                   │                   │
    ├───────────────────┼───────────────────┤
    │                   │                   │
    │  RealSense D435   │  Prosumer         │
    │  + Python SDK     │  Development      │
    │                   │                   │
    ├───────────────────┼───────────────────┤
    │                   │                   │
    │  MediaPipe        │  Consumer         │
    │  + Webcam         │  Quick Start      │
    │                   │                   │
    └───────────────────┴───────────────────┘
        LOW ────────────────────────── HIGH
                    FIDELITY
```

---

## 🚀 Quick Start: MediaPipe (Webcam)

The simplest approach - works with any webcam:

```bash
pip install mediapipe opencv-python
python mediapipe_gestures.py
```

---

## 🔧 Integration with OneMind

All gesture systems publish to NATS topics:

```
gesture.raw.{source}     # Raw tracking data from Kinect/MediaPipe/etc
gesture.processed.hands  # Normalized hand position and state
gesture.classified       # High-level gestures (point, grab, swipe, pinch)
gesture.intent           # Interpreted user intent (select, navigate, confirm)
```

Subscribers (drones, robots, Unreal visualizer, Flutter UI) react to these topics.

---

## 🎯 OneMind OS Use Cases

| Use Case          | Gesture         | Action                     |
| ----------------- | --------------- | -------------------------- |
| **Drone Control** | Point           | Fly in pointed direction   |
| **Drone Control** | Open Palm       | Emergency stop / hover     |
| **Robot Arm**     | Pinch + Move    | Grab and position object   |
| **Smart Home**    | Swipe           | Navigate dashboard panels  |
| **CAD Review**    | Two-hand rotate | Rotate 3D model            |
| **Confirmation**  | Thumbs Up       | Confirm action             |
| **Cancel**        | Open Palm       | Cancel/Stop current action |
