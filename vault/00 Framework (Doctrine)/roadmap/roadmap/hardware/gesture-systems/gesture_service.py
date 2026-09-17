"""
gesture_service.py — MediaPipe Gesture Detection → NATS Publisher

Runs on any node with a webcam. Detects hand gestures with MediaPipe Hands,
publishes each detected gesture to NATS so Legacy and other runtimes can react.

Adapted from:
  - reference/concept-bytes/hand_tracking/handtracking.py
  - reference/gesture_systems/mediapipe_gestures.py

NATS Subject Published:   physical.gesture.local.detected
NATS Subject Control:     physical.gesture.local.command  (start/stop/config)

Payload:
  {
    "gesture": "THUMBS_UP",      # gesture name
    "hand": "right",             # "left" or "right"
    "confidence": 0.95,
    "position": [0.45, 0.32],    # normalized x,y (0-1) of palm center
    "timestamp": "2026-03-06T..."
  }

Usage:
    python3 gesture_service.py                       # with video window
    python3 gesture_service.py --no-video            # headless (saves CPU)
    python3 gesture_service.py --nats nats://localhost:4222

Setup:
    pip install mediapipe opencv-python nats-py

Dependencies:
    - mediapipe
    - opencv-python
    - nats-py
"""

import asyncio
import argparse
import json
import math
import time
from datetime import datetime, timezone

try:
    import cv2
    import mediapipe as mp
    CV_AVAILABLE = True
except ImportError:
    CV_AVAILABLE = False

try:
    import nats
    NATS_AVAILABLE = True
except ImportError:
    NATS_AVAILABLE = False

# ─── Gesture constants ────────────────────────────────────────────────────────

GESTURE_NONE        = "NONE"
GESTURE_POINT       = "POINT"          # index finger up only
GESTURE_OPEN_PALM   = "OPEN_PALM"      # all 5 fingers extended
GESTURE_CLOSED_FIST = "CLOSED_FIST"    # all fingers curled
GESTURE_THUMBS_UP   = "THUMBS_UP"      # only thumb extended
GESTURE_PEACE       = "PEACE"          # index + middle extended
GESTURE_PINCH       = "PINCH"          # thumb + index close together
GESTURE_THREE       = "THREE"          # index + middle + ring extended

# Pinch threshold (normalized landmark distance)
PINCH_THRESHOLD = 0.08

# ─── Landmark helpers ─────────────────────────────────────────────────────────

def _dist(a, b):
    return math.sqrt((a.x - b.x) ** 2 + (a.y - b.y) ** 2)

def _extended(tip, pip):
    """True if finger is extended (tip above PIP joint in image-y)."""
    return tip.y < pip.y

def classify_gesture(lm, handedness: str) -> dict:
    """
    Given MediaPipe hand landmarks, return gesture name and confidence.
    lm = hand_landmarks.landmark (list of 21 NormalizedLandmark)
    """
    # Thumb: different left/right logic
    if handedness == "Right":
        thumb_up = lm[4].x < lm[3].x
    else:
        thumb_up = lm[4].x > lm[3].x

    index_up = _extended(lm[8],  lm[6])
    middle_up = _extended(lm[12], lm[10])
    ring_up   = _extended(lm[16], lm[14])
    pinky_up  = _extended(lm[20], lm[18])

    pinch_dist = _dist(lm[4], lm[8])

    # Palm center (average of wrist + MCP joints)
    palm_x = sum(lm[i].x for i in [0, 5, 9, 13, 17]) / 5
    palm_y = sum(lm[i].y for i in [0, 5, 9, 13, 17]) / 5

    fingers = [thumb_up, index_up, middle_up, ring_up, pinky_up]
    count = sum(fingers)

    if pinch_dist < PINCH_THRESHOLD:
        gesture = GESTURE_PINCH
    elif count == 5:
        gesture = GESTURE_OPEN_PALM
    elif count == 0:
        gesture = GESTURE_CLOSED_FIST
    elif thumb_up and not index_up and not middle_up and not ring_up and not pinky_up:
        gesture = GESTURE_THUMBS_UP
    elif index_up and middle_up and not ring_up and not pinky_up and not thumb_up:
        gesture = GESTURE_PEACE
    elif index_up and middle_up and ring_up and not pinky_up and not thumb_up:
        gesture = GESTURE_THREE
    elif index_up and not middle_up and not ring_up and not pinky_up:
        gesture = GESTURE_POINT
    else:
        gesture = GESTURE_NONE

    return {
        "gesture": gesture,
        "hand": handedness.lower(),
        "confidence": round(1.0 - min(pinch_dist, 1.0), 2),
        "position": [round(palm_x, 3), round(palm_y, 3)],
        "timestamp": datetime.now(timezone.utc).isoformat(),
    }

# ─── Main service ─────────────────────────────────────────────────────────────

async def run(nats_url: str, show_video: bool, camera_id: int, cooldown_ms: int):
    if not CV_AVAILABLE:
        raise RuntimeError("opencv-python and mediapipe required: pip install mediapipe opencv-python")

    # Connect to NATS
    nc = None
    if NATS_AVAILABLE:
        try:
            nc = await nats.connect(nats_url)
            print(f"[gesture_service] Connected to NATS at {nats_url}")
        except Exception as e:
            print(f"[gesture_service] NATS unavailable ({e}) — running in print-only mode")
    else:
        print("[gesture_service] nats-py not installed — running in print-only mode")

    # MediaPipe setup
    mp_hands = mp.solutions.hands
    mp_draw  = mp.solutions.drawing_utils
    hands = mp_hands.Hands(
        static_image_mode=False,
        max_num_hands=1,
        min_detection_confidence=0.75,
        min_tracking_confidence=0.6,
    )

    cap = cv2.VideoCapture(camera_id)
    if not cap.isOpened():
        raise RuntimeError(f"Could not open camera {camera_id}")

    print(f"[gesture_service] Camera {camera_id} open. Press Q to quit.")

    last_gesture = None
    last_publish_time = 0.0

    while True:
        success, frame = cap.read()
        if not success:
            await asyncio.sleep(0.01)
            continue

        frame = cv2.flip(frame, 1)
        rgb   = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        results = hands.process(rgb)

        if results.multi_hand_landmarks and results.multi_handedness:
            for hand_lm, hand_info in zip(results.multi_hand_landmarks, results.multi_handedness):
                label = hand_info.classification[0].label  # "Left" or "Right"
                payload = classify_gesture(hand_lm.landmark, label)

                now = time.monotonic() * 1000  # ms
                changed = payload["gesture"] != last_gesture
                cooled  = (now - last_publish_time) > cooldown_ms

                if (changed or cooled) and payload["gesture"] != GESTURE_NONE:
                    last_gesture = payload["gesture"]
                    last_publish_time = now

                    msg = json.dumps(payload)
                    print(f"[gesture] {payload['gesture']} ({payload['hand']}) @ {payload['position']}")

                    if nc:
                        await nc.publish("physical.gesture.local.detected", msg.encode())

                if show_video:
                    mp_draw.draw_landmarks(frame, hand_lm, mp_hands.HAND_CONNECTIONS)
                    cv2.putText(
                        frame, payload["gesture"], (10, 40),
                        cv2.FONT_HERSHEY_SIMPLEX, 1.2, (0, 255, 0), 2
                    )

        if show_video:
            cv2.imshow("Gesture Service", frame)
            if cv2.waitKey(1) & 0xFF == ord("q"):
                break

        await asyncio.sleep(0)  # yield to event loop

    cap.release()
    if show_video:
        cv2.destroyAllWindows()
    if nc:
        await nc.close()


def main():
    parser = argparse.ArgumentParser(description="MediaPipe gesture → NATS publisher")
    parser.add_argument("--nats",      default="nats://localhost:4222", help="NATS server URL")
    parser.add_argument("--no-video",  action="store_true",             help="Headless mode (no video window)")
    parser.add_argument("--camera",    type=int, default=0,             help="Camera device ID (default: 0)")
    parser.add_argument("--cooldown",  type=int, default=500,           help="Min ms between repeated publishes of same gesture")
    args = parser.parse_args()

    asyncio.run(run(
        nats_url=args.nats,
        show_video=not args.no_video,
        camera_id=args.camera,
        cooldown_ms=args.cooldown,
    ))


if __name__ == "__main__":
    main()
