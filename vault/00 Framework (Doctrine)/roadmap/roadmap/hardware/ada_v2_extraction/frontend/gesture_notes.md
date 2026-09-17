# Gesture Control Notes

## Status in ADA V2

Gesture control appears to be implemented in the **frontend** using MediaPipe's JavaScript SDK, not in the Python backend.

From reviewing the repository:

- MediaPipe Face Landmarker is used in `authenticator.py` for face recognition
- Gesture detection likely uses MediaPipe Hands in the React/Electron frontend
- The backend receives gesture events via Socket.IO

## ADA V2 Gesture Workflow (Inferred)

```
┌─────────────────────────────────────────┐
│  Electron Frontend (React + Three.js)  │
│  ├── MediaPipe Hands (JavaScript)       │
│  │   └── @mediapipe/hands npm package  │
│  ├── Gesture Recognition Logic          │
│  │   └── detect pinch, swipe, grab, etc│
│  └── Socket.IO emit('gesture', {...})   │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│  Python Backend (FastAPI)               │
│  └── Socket.IO on('gesture', handler)   │
│      └── Trigger actions, update UI     │
└─────────────────────────────────────────┘
```

## For OneMind OS Integration

### Option A: Flutter + MediaPipe

Google offers MediaPipe packages for Flutter:

- `google_ml_kit` - Face detection, pose detection
- `camera` - Camera access

```yaml
# pubspec.yaml
dependencies:
  google_ml_kit: ^0.16.0
  camera: ^0.10.0
```

Example gesture detection flow:

```dart
class GestureController {
  final handDetector = GoogleMlKit.vision.handLandmarkDetector();

  void processFrame(CameraImage image) async {
    final hand = await handDetector.processImage(inputImage);
    if (detectPinchGesture(hand)) {
      // Emit to NATS
      natsClient.publish('gesture.pinch', {'confidence': 0.95});
    }
  }
}
```

### Option B: Backend MediaPipe (Python)

For server-side processing (useful for robotics/drones):

```python
import mediapipe as mp
from mediapipe.tasks.python import vision

mp_hands = mp.solutions.hands

class GestureDetector:
    def __init__(self):
        self.hands = mp_hands.Hands(
            static_image_mode=False,
            max_num_hands=2,
            min_detection_confidence=0.5
        )

    def detect(self, frame_rgb):
        results = self.hands.process(frame_rgb)
        if results.multi_hand_landmarks:
            for hand_landmarks in results.multi_hand_landmarks:
                # Analyze landmarks for gestures
                gesture = self._classify_gesture(hand_landmarks)
                yield gesture

    def _classify_gesture(self, landmarks):
        # Example: Detect pinch (thumb tip close to index tip)
        thumb_tip = landmarks.landmark[4]
        index_tip = landmarks.landmark[8]

        distance = ((thumb_tip.x - index_tip.x)**2 +
                   (thumb_tip.y - index_tip.y)**2)**0.5

        if distance < 0.05:
            return "pinch"
        return "none"
```

## Gesture Types to Support

| Gesture          | MediaPipe Landmarks            | Use Case         |
| ---------------- | ------------------------------ | ---------------- |
| Pinch            | Thumb tip + Index tip distance | Select, zoom     |
| Grab (fist)      | All finger tips curled         | Grab object      |
| Point            | Index extended, others curled  | Select target    |
| Swipe            | Hand movement velocity         | Navigate, scroll |
| Thumbs Up        | Thumb extended upward          | Confirm, approve |
| Open Hand        | All fingers extended           | Stop, pause      |
| Two-finger pinch | Thumb + Index on both hands    | Scale object     |

## Integration Priority

For OneMind OS, gesture control would be most valuable for:

1. **3D CAD Viewer** - Rotate, zoom, pan models with gestures
2. **Robotics Control** - Point to target, grab to move
3. **Drone Piloting** - Gesture-based navigation commands
4. **Presentation Mode** - Swipe through slides, point at elements

## Dependencies

```bash
# For Python backend gesture detection
pip install mediapipe opencv-python

# For Flutter frontend
# Add to pubspec.yaml:
# google_ml_kit: ^0.16.0
# camera: ^0.10.0
```

## Next Steps

1. Decide on frontend vs backend gesture processing
2. If frontend: Implement in Flutter using `google_ml_kit`
3. If backend: Create a gesture service that processes camera streams
4. Define NATS topics: `gesture.detected`, `gesture.{type}`
5. Connect to CAD viewer, drone controller, etc.

---

**Note:** Full gesture detection code was not extracted from ADA V2 as it appears
to be in the JavaScript frontend. The face authentication code in `authenticator.py`
demonstrates how to use MediaPipe in Python, which can be adapted for hand/gesture
detection using similar patterns.
