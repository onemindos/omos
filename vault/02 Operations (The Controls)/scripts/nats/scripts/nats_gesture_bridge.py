"""
NATS Gesture Bridge - Generic bridge between gesture sources and NATS

This bridge:
1. Receives gesture data from various sources (OSC, direct Python, etc.)
2. Normalizes the data format
3. Publishes to NATS topics
4. Optionally stores gesture history for pattern recognition

Topics Published:
    gesture.raw.{source}     - Raw data from each source
    gesture.processed.hands  - Normalized hand data
    gesture.classified       - High-level classified gestures
    gesture.intent           - Interpreted user intent

Dependencies:
    pip install nats-py python-osc

Usage:
    # Start the bridge
    python nats_gesture_bridge.py

    # Or import and start programmatically
    from nats_gesture_bridge import GestureBridge
    bridge = GestureBridge()
    await bridge.start()
"""

import asyncio
import json
import time
from dataclasses import dataclass, asdict, field
from typing import Dict, List, Optional, Callable, Any
from collections import deque
from enum import Enum


class GestureSource(Enum):
    MEDIAPIPE = "mediapipe"
    KINECT = "kinect"
    REALSENSE = "realsense"
    TOUCHDESIGNER = "touchdesigner"
    LEAP = "leap"


@dataclass
class NormalizedGesture:
    """Normalized gesture format for NATS publishing."""

    source: str
    gesture_type: str
    gesture_id: int
    hand: str
    confidence: float
    position: Dict[str, float]  # x, y, z (normalized 0-1)
    direction: Optional[Dict[str, float]] = None
    velocity: float = 0.0
    timestamp: float = field(default_factory=time.time)
    metadata: Dict[str, Any] = field(default_factory=dict)

    def to_json(self) -> str:
        return json.dumps(asdict(self))


class GestureBridge:
    """
    Bridge gestures from various sources to NATS.

    Handles:
    - OSC input (from TouchDesigner, etc.)
    - Direct Python callbacks
    - Gesture normalization
    - History tracking for pattern recognition
    - NATS publishing
    """

    def __init__(
        self,
        nats_url: str = "nats://localhost:4222",
        osc_port: int = 5005,
        history_size: int = 100,
    ):
        self.nats_url = nats_url
        self.osc_port = osc_port
        self.history_size = history_size

        self.nc = None  # NATS connection
        self.osc_server = None

        # Gesture history for pattern recognition
        self.history: Dict[str, deque] = {
            "left": deque(maxlen=history_size),
            "right": deque(maxlen=history_size),
        }

        # Intent handlers
        self.intent_handlers: List[Callable] = []

        self._running = False

    async def connect(self):
        """Connect to NATS."""
        try:
            import nats

            self.nc = await nats.connect(self.nats_url)
            print(f"[GestureBridge] Connected to NATS at {self.nats_url}")
        except Exception as e:
            print(f"[GestureBridge] NATS connection failed: {e}")
            self.nc = None

    async def publish_raw(self, source: GestureSource, data: dict):
        """Publish raw gesture data."""
        if self.nc:
            topic = f"gesture.raw.{source.value}"
            await self.nc.publish(topic, json.dumps(data).encode())

    async def publish_gesture(self, gesture: NormalizedGesture):
        """Publish normalized gesture."""
        if not self.nc:
            return

        # Add to history
        self.history[gesture.hand].append(gesture)

        # Publish to classified topic
        await self.nc.publish("gesture.classified", gesture.to_json().encode())

        # Check for intent patterns
        intent = self._recognize_intent(gesture)
        if intent:
            await self.nc.publish("gesture.intent", json.dumps(intent).encode())
            print(f"[GestureBridge] Intent: {intent['action']}")

    def _recognize_intent(self, gesture: NormalizedGesture) -> Optional[dict]:
        """
        Recognize higher-level intent from gesture patterns.

        Example intents:
        - "select": Point + hold
        - "zoom_in": Two-hand pinch moving apart
        - "zoom_out": Two-hand pinch moving together
        - "dismiss": Swipe away gesture
        - "confirm": Thumbs up
        - "cancel": Open palm (stop)
        """
        history = self.history[gesture.hand]

        if len(history) < 2:
            return None

        # Simple intent recognition

        # Confirm: Thumbs up
        if gesture.gesture_type == "thumbs_up" and gesture.confidence > 0.8:
            return {
                "action": "confirm",
                "confidence": gesture.confidence,
                "timestamp": gesture.timestamp,
            }

        # Cancel/Stop: Open palm
        if gesture.gesture_type == "open_palm" and gesture.confidence > 0.8:
            return {
                "action": "stop",
                "confidence": gesture.confidence,
                "timestamp": gesture.timestamp,
            }

        # Select: Point and hold
        if gesture.gesture_type == "point":
            # Check if pointing in same direction for multiple frames
            recent = list(history)[-5:]
            if all(g.gesture_type == "point" for g in recent):
                avg_direction = {
                    "x": sum(g.direction["x"] for g in recent if g.direction)
                    / len(recent),
                    "y": sum(g.direction["y"] for g in recent if g.direction)
                    / len(recent),
                    "z": sum(g.direction.get("z", 0) for g in recent if g.direction)
                    / len(recent),
                }
                return {
                    "action": "select",
                    "direction": avg_direction,
                    "confidence": gesture.confidence,
                    "timestamp": gesture.timestamp,
                }

        # Swipe navigation
        if gesture.gesture_type.startswith("swipe_"):
            direction = gesture.gesture_type.split("_")[1]
            return {
                "action": f"navigate_{direction}",
                "velocity": gesture.velocity,
                "confidence": gesture.confidence,
                "timestamp": gesture.timestamp,
            }

        return None

    def register_osc_handlers(self, dispatcher):
        """Register OSC message handlers."""

        def handle_gesture_osc(address, *args):
            """Handle /gesture/* OSC messages."""
            try:
                # Expected format: type, confidence, hand, x, y, z, dir_x, dir_y, dir_z, velocity
                gesture = NormalizedGesture(
                    source=GestureSource.TOUCHDESIGNER.value,
                    gesture_type=args[0],
                    gesture_id=args[1] if len(args) > 1 else 0,
                    hand=args[2] if len(args) > 2 else "right",
                    confidence=args[3] if len(args) > 3 else 0.5,
                    position={
                        "x": args[4] if len(args) > 4 else 0,
                        "y": args[5] if len(args) > 5 else 0,
                        "z": args[6] if len(args) > 6 else 0,
                    },
                    direction=(
                        {
                            "x": args[7] if len(args) > 7 else 0,
                            "y": args[8] if len(args) > 8 else 0,
                            "z": args[9] if len(args) > 9 else 0,
                        }
                        if len(args) > 7
                        else None
                    ),
                    velocity=args[10] if len(args) > 10 else 0,
                )

                asyncio.create_task(self.publish_gesture(gesture))

            except Exception as e:
                print(f"[OSC] Error parsing gesture: {e}")

        def handle_skeleton_osc(address, *args):
            """Handle /skeleton/* OSC messages."""
            # Full skeleton data from Kinect via TouchDesigner
            asyncio.create_task(
                self.publish_raw(
                    GestureSource.TOUCHDESIGNER, {"address": address, "args": args}
                )
            )

        dispatcher.map("/gesture/*", handle_gesture_osc)
        dispatcher.map("/skeleton/*", handle_skeleton_osc)

    async def start_osc_server(self):
        """Start OSC server for receiving gestures."""
        try:
            from pythonosc import dispatcher as osc_dispatcher
            from pythonosc import osc_server

            dispatcher = osc_dispatcher.Dispatcher()
            self.register_osc_handlers(dispatcher)

            self.osc_server = osc_server.AsyncIOOSCUDPServer(
                ("0.0.0.0", self.osc_port), dispatcher, asyncio.get_event_loop()
            )

            transport, protocol = await self.osc_server.create_serve_endpoint()
            print(f"[GestureBridge] OSC server listening on port {self.osc_port}")

        except ImportError:
            print("[GestureBridge] python-osc not installed, OSC disabled")
        except Exception as e:
            print(f"[GestureBridge] OSC server error: {e}")

    async def start(self):
        """Start the gesture bridge."""
        self._running = True

        await self.connect()
        await self.start_osc_server()

        print("[GestureBridge] Running. Press Ctrl+C to stop.")

        try:
            while self._running:
                await asyncio.sleep(1)
        except asyncio.CancelledError:
            pass
        finally:
            await self.stop()

    async def stop(self):
        """Stop the gesture bridge."""
        self._running = False
        if self.nc:
            await self.nc.close()
        print("[GestureBridge] Stopped")

    # Direct Python API
    async def process_gesture(
        self,
        gesture_type: str,
        hand: str = "right",
        position: tuple = (0.5, 0.5, 0.5),
        direction: tuple = None,
        confidence: float = 0.8,
        source: GestureSource = GestureSource.MEDIAPIPE,
    ):
        """
        Process a gesture from Python code.

        Example:
            await bridge.process_gesture(
                gesture_type="point",
                hand="right",
                position=(0.6, 0.4, 0.3),
                direction=(0.1, -0.2, 0.0),
                confidence=0.9
            )
        """
        gesture = NormalizedGesture(
            source=source.value,
            gesture_type=gesture_type,
            gesture_id=0,
            hand=hand,
            confidence=confidence,
            position={"x": position[0], "y": position[1], "z": position[2]},
            direction=(
                {"x": direction[0], "y": direction[1], "z": direction[2]}
                if direction
                else None
            ),
        )

        await self.publish_gesture(gesture)


# =============================================================================
# Main
# =============================================================================


async def main():
    """Run the gesture bridge."""
    bridge = GestureBridge(nats_url="nats://localhost:4222", osc_port=5005)

    await bridge.start()


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\nShutting down...")
