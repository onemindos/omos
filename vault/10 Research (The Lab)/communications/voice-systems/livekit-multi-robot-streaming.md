# LiveKit Multi-Robot Streaming Integration
## Real-Time Video Coordination with Brilliant Frame + 5 Robots

**Date**: February 7, 2026
**Use Case**: Stream Brilliant Frame glasses camera to multiple robots while viewing all robot cameras simultaneously
**Voice Control**: "Legacy, come pick this up and bring it to the barn"

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Use Case Scenario](#use-case-scenario)
3. [Technical Implementation](#technical-implementation)
4. [Code Examples](#code-examples)
5. [Data Flow Diagrams](#data-flow-diagrams)
6. [Deployment](#deployment)
7. [Cost Analysis](#cost-analysis)
8. [Performance Metrics](#performance-metrics)

---

## Architecture Overview

### System Components

```
┌──────────────────────────────────────────────────────────┐
│ Brilliant Frame (AR Glasses)                             │
│  - Camera: 640x480 @ 30fps                               │
│  - Microphone: Voice commands                            │
│  - Display: AR overlay for robot feeds                   │
│  - Connection: BLE → Flutter App                         │
└────────────────┬─────────────────────────────────────────┘
                 │ WebRTC (via Flutter)
┌────────────────▼─────────────────────────────────────────┐
│ LiveKit Cloud (Real-Time Media Server)                   │
│  Room: "property_inspection"                             │
│                                                           │
│  Participants (6 total):                                 │
│  ┌─────────────────────────────────────────────────┐    │
│  │ 1. brilliant_frame (You)                        │    │
│  │    - Publishes: camera + microphone             │    │
│  │    - Subscribes: 5 robot cameras                │    │
│  ├─────────────────────────────────────────────────┤    │
│  │ 2. unitree_go2                                  │    │
│  │    - Publishes: camera                          │    │
│  │    - Subscribes: brilliant_frame camera         │    │
│  ├─────────────────────────────────────────────────┤    │
│  │ 3. unitree_g1                                   │    │
│  │    - Publishes: camera                          │    │
│  │    - Subscribes: brilliant_frame camera         │    │
│  ├─────────────────────────────────────────────────┤    │
│  │ 4. openmind_bot_1                               │    │
│  │    - Publishes: camera                          │    │
│  │    - Subscribes: brilliant_frame camera         │    │
│  ├─────────────────────────────────────────────────┤    │
│  │ 5. drone_1                                      │    │
│  │    - Publishes: camera                          │    │
│  │    - Subscribes: brilliant_frame camera         │    │
│  ├─────────────────────────────────────────────────┤    │
│  │ 6. workshop_bot                                 │    │
│  │    - Publishes: camera                          │    │
│  │    - Subscribes: brilliant_frame camera         │    │
│  └─────────────────────────────────────────────────┘    │
│                                                           │
│  Tracks:                                                 │
│  - 6 video tracks (1 from each participant)              │
│  - 1 audio track (from brilliant_frame)                  │
└────────────────┬─────────────────────────────────────────┘
                 │ Audio stream
┌────────────────▼─────────────────────────────────────────┐
│ Agno Backend (Legacy Agent)                              │
│  - Subscribes to brilliant_frame audio track             │
│  - Transcribes voice commands (Gemini)                   │
│  - Analyzes Frame camera with vision (Gemini)            │
│  - Selects closest robot for task                        │
│  - Sends commands to robot controllers                   │
└────────────────┬─────────────────────────────────────────┘
                 │ REST API commands
┌────────────────▼─────────────────────────────────────────┐
│ Robot Controllers (FastAPI microservices)                │
│  - Unitree Controller (unitree_sdk2)                     │
│  - OpenMind Controller (OM1 runtime)                     │
│  - Drone Controller (DJI SDK)                            │
└───────────────────────────────────────────────────────────┘
```

---

## Use Case Scenario

### Detailed Walkthrough

**Initial State**:
- You're wearing Brilliant Frame glasses in your house
- 5 robots are distributed across property:
  - Unitree Go2: Near barn (charging)
  - Unitree G1: Workshop
  - OpenMind Bot: Solar panel area
  - Drone: Landing pad
  - Workshop Bot: Indoor shop

**Step-by-Step Flow**:

```
┌─ PHASE 1: Setup Streaming Session ────────────────────┐
│ Time: 0:00                                             │
│                                                        │
│ 1. You activate voice mode in OneMind OS app          │
│    "Legacy, start collaboration mode"                 │
│                                                        │
│ 2. Flutter app:                                       │
│    - Connects Frame via BLE                           │
│    - Requests LiveKit token from Agno                 │
│    - Joins LiveKit room "property_inspection"         │
│    - Publishes Frame camera + mic                     │
│                                                        │
│ 3. All 5 robots automatically join LiveKit room       │
│    - Each publishes their camera feed                 │
│    - Each subscribes to your Frame camera             │
│                                                        │
│ 4. Frame AR overlay shows:                            │
│    - 5 robot camera feeds (thumbnail grid)            │
│    - Robot status indicators                          │
│    - "Collaboration mode active" message              │
│                                                        │
│ Result: Everyone can see everyone                     │
└────────────────────────────────────────────────────────┘

┌─ PHASE 2: Object Identification ──────────────────────┐
│ Time: 0:30                                             │
│                                                        │
│ 1. You look at red toolbox on workbench               │
│    Frame camera captures view → streams to LiveKit    │
│                                                        │
│ 2. All 5 robots receive your Frame view               │
│    Each robot's AI processes visual data:             │
│    - Detects: Red toolbox                             │
│    - Location: Indoor workbench                       │
│    - Distance from each robot calculated              │
│                                                        │
│ 3. Agno also receives Frame video stream              │
│    Continuously analyzes with Gemini Vision:          │
│    - Object: "Red metal toolbox"                      │
│    - Location: "Indoor workbench, north wall"         │
│    - Context: "User's gaze focused on toolbox"        │
│                                                        │
│ Result: Shared visual context established             │
└────────────────────────────────────────────────────────┘

┌─ PHASE 3: Voice Command ───────────────────────────────┐
│ Time: 0:45                                             │
│                                                        │
│ You: "Legacy, come pick this up and bring it to barn" │
│                                                        │
│ 1. Frame microphone captures audio                    │
│    → Streams to LiveKit audio track                   │
│                                                        │
│ 2. Agno subscribes to audio track                     │
│    - Gemini transcribes: "Legacy, come pick this up   │
│      and bring it to barn"                            │
│    - Detects wake word: "Legacy"                      │
│                                                        │
│ 3. Legacy agent processes command + visual context:   │
│    Input:                                             │
│    - Voice: "come pick this up"                       │
│    - Visual: Red toolbox on workbench                 │
│    - Destination: "barn"                              │
│                                                        │
│    Processing:                                        │
│    - Understands "this" = red toolbox (from vision)   │
│    - Task: fetch + deliver                            │
│    - Queries robot positions (via status APIs)        │
│    - Calculates distances to workbench                │
│                                                        │
│ Result: Command understood with full context          │
└────────────────────────────────────────────────────────┘

┌─ PHASE 4: Robot Selection ─────────────────────────────┐
│ Time: 0:48                                             │
│                                                        │
│ Legacy agent decision logic:                          │
│                                                        │
│ Robot Distances to Workbench:                         │
│ ├─ Unitree Go2: 150m (outdoor near barn)             │
│ ├─ Unitree G1: 5m (already in workshop) ✓ SELECTED   │
│ ├─ OpenMind Bot: 200m (solar panel area)             │
│ ├─ Drone: 100m (landing pad)                         │
│ └─ Workshop Bot: 2m (indoor shop)                    │
│                                                        │
│ Selection criteria:                                   │
│ 1. Can manipulate objects? (needs gripper/arm)        │
│    - Go2: No ❌                                        │
│    - G1: Yes ✓                                        │
│    - OpenMind: No ❌                                   │
│    - Drone: No ❌                                      │
│    - Workshop Bot: Yes ✓                              │
│                                                        │
│ 2. Closest with manipulation capability?              │
│    - Workshop Bot: 2m                                 │
│    - G1: 5m                                           │
│                                                        │
│ 3. Can navigate to barn? (indoor → outdoor)           │
│    - Workshop Bot: No (indoor only) ❌                 │
│    - G1: Yes (mobile, indoor/outdoor) ✓               │
│                                                        │
│ **Decision: Unitree G1 selected**                     │
│                                                        │
│ Legacy agent sends command:                           │
│ POST /api/unitree/g1/fetch                            │
│ {                                                     │
│   "object": "red_toolbox",                            │
│   "pickup_location": {"x": 10.5, "y": 3.2, "z": 0},  │
│   "delivery_location": "barn_entrance",               │
│   "guidance": "frame_camera_stream"                   │
│ }                                                     │
│                                                        │
│ Result: G1 receives task, starts execution            │
└────────────────────────────────────────────────────────┘

┌─ PHASE 5: Robot Execution with Real-Time Streaming ───┐
│ Time: 0:50 - 3:00                                      │
│                                                        │
│ [0:50] G1 acknowledges task                           │
│    - Subscribes to Frame camera (enhanced priority)   │
│    - Uses your view for visual navigation             │
│    - Starts moving toward workbench                   │
│                                                        │
│ [0:55] You see G1's camera feed in Frame AR overlay   │
│    - G1 approaching workbench (POV from G1)           │
│    - Other 4 robot feeds still visible (smaller)      │
│    - Frame displays: "G1: Navigating to target"       │
│                                                        │
│ [1:10] G1 reaches workbench                           │
│    - G1's camera shows toolbox from its perspective   │
│    - G1's AI compares:                                │
│      • Your Frame view (what you see)                 │
│      • Its own camera view                            │
│      → Confirms: "Same red toolbox"                   │
│                                                        │
│ [1:15] G1 extends manipulator arm                     │
│    - You watch in real-time via G1's camera feed      │
│    - Frame AR overlay shows:                          │
│      "G1: Picking up toolbox"                         │
│                                                        │
│ [1:20] Toolbox secured                                │
│    - G1 confirms grip via force sensors               │
│    - Sends status to Agno:                            │
│      {"status": "object_secured", "grip_force": 85}   │
│    - Legacy notifies you via Frame:                   │
│      "Toolbox secured. Heading to barn."              │
│                                                        │
│ [1:25 - 2:45] G1 navigates to barn                    │
│    - Autonomous SLAM navigation                       │
│    - You see G1's camera feed (journey POV)           │
│    - Frame displays GPS trail on AR overlay           │
│    - Other robots continue their work (visible)       │
│                                                        │
│ [2:50] G1 reaches barn entrance                       │
│    - Gently places toolbox on designated shelf        │
│    - Confirms delivery via camera + weight sensors    │
│                                                        │
│ [3:00] Task complete                                  │
│    - G1 sends completion status to Agno               │
│    - Legacy speaks via Frame:                         │
│      "Task complete. Toolbox delivered to barn.       │
│       G1 battery: 67%. Returning to standby."         │
│    - Frame AR overlay shows:                          │
│      ✓ "Mission Complete"                             │
│                                                        │
│ Result: Task executed successfully with full visibility│
└─────────────────────────────────────────────────────────┘
```

---

## Technical Implementation

### 1. LiveKit Room Setup

**File**: `backend/services/livekit_coordinator.py`

```python
from livekit import api, rtc
import asyncio
from typing import Dict, List

class LiveKitMultiRobotCoordinator:
    """Manages LiveKit room for multi-robot collaboration"""

    def __init__(self):
        self.livekit_url = "wss://livekit.onemind.example.com"
        self.api_key = os.getenv("LIVEKIT_API_KEY")
        self.api_secret = os.getenv("LIVEKIT_API_SECRET")
        self.api_client = api.LiveKitAPI(
            url=self.livekit_url,
            api_key=self.api_key,
            api_secret=self.api_secret
        )
        self.active_rooms: Dict[str, rtc.Room] = {}

    async def create_collaboration_room(self, room_name: str = "property_inspection"):
        """Create LiveKit room for property-wide collaboration"""

        # Create room with settings
        room = await self.api_client.room.create_room(
            api.CreateRoomRequest(
                name=room_name,
                empty_timeout=300,  # Room stays alive 5 min after last participant
                max_participants=10,  # You + up to 9 robots
                metadata=json.dumps({
                    "type": "multi_robot_collaboration",
                    "created_at": datetime.now().isoformat(),
                    "owner": "brilliant_frame"
                })
            )
        )

        return room

    async def generate_token(
        self,
        identity: str,
        name: str,
        room_name: str,
        can_publish: bool = True,
        can_subscribe: bool = True
    ) -> str:
        """Generate access token for participant"""

        token = api.AccessToken(self.api_key, self.api_secret)
        token.with_identity(identity)
        token.with_name(name)
        token.with_grants(api.VideoGrants(
            room_join=True,
            room=room_name,
            can_publish=can_publish,
            can_subscribe=can_subscribe,
            can_publish_data=True  # For metadata/commands
        ))

        # Token valid for 6 hours
        token.with_ttl(timedelta(hours=6))

        return token.to_jwt()

    async def join_as_brilliant_frame(self, room_name: str):
        """Join room from Brilliant Frame perspective (via Flutter app)"""

        room = rtc.Room()

        # Generate token
        token = await self.generate_token(
            identity="brilliant_frame",
            name="Your Glasses",
            room_name=room_name
        )

        # Connect
        await room.connect(self.livekit_url, token)

        # Store room reference
        self.active_rooms["brilliant_frame"] = room

        # Setup event handlers
        room.on("track_subscribed", self.on_track_subscribed)
        room.on("track_unsubscribed", self.on_track_unsubscribed)
        room.on("participant_connected", self.on_participant_connected)

        return room

    async def join_as_robot(
        self,
        robot_id: str,
        robot_name: str,
        room_name: str
    ):
        """Join room from robot perspective"""

        room = rtc.Room()

        # Generate token for robot
        token = await self.generate_token(
            identity=robot_id,
            name=robot_name,
            room_name=room_name
        )

        # Connect
        await room.connect(self.livekit_url, token)

        # Store room reference
        self.active_rooms[robot_id] = room

        return room

    async def on_track_subscribed(
        self,
        track: rtc.Track,
        publication: rtc.TrackPublication,
        participant: rtc.RemoteParticipant
    ):
        """Handle new track subscription"""

        print(f"Subscribed to {participant.identity}'s {track.kind} track")

        if track.kind == rtc.TrackKind.KIND_VIDEO:
            # Start processing video frames
            asyncio.create_task(self.process_video_track(track, participant))
        elif track.kind == rtc.TrackKind.KIND_AUDIO:
            # Start processing audio for voice commands
            asyncio.create_task(self.process_audio_track(track, participant))

    async def process_video_track(
        self,
        track: rtc.VideoTrack,
        participant: rtc.RemoteParticipant
    ):
        """Process incoming video frames"""

        async for frame in rtc.VideoStream(track):
            # Convert frame to bytes for Gemini vision
            frame_bytes = await frame.to_argb_bytes()

            if participant.identity == "brilliant_frame":
                # This is user's Frame camera
                # Analyze for object detection
                await self.analyze_frame_view(frame_bytes)

    async def process_audio_track(
        self,
        track: rtc.AudioTrack,
        participant: rtc.RemoteParticipant
    ):
        """Process incoming audio for voice commands"""

        if participant.identity == "brilliant_frame":
            # This is user's microphone
            async for frame in rtc.AudioStream(track):
                audio_bytes = frame.data

                # Send to Gemini for transcription
                transcript = await self.transcribe_audio(audio_bytes)

                if "legacy" in transcript.lower():
                    # Wake word detected
                    await self.process_voice_command(transcript)

    async def on_participant_connected(self, participant: rtc.RemoteParticipant):
        """Handle new participant joining"""

        print(f"Participant joined: {participant.identity} ({participant.name})")

        # Notify all participants
        await self.broadcast_data_message({
            "type": "participant_joined",
            "identity": participant.identity,
            "name": participant.name,
            "timestamp": datetime.now().isoformat()
        })
```

---

### 2. Brilliant Frame Integration

**File**: `frontend/lib/services/brilliant_frame_livekit.dart`

```dart
import 'package:livekit_client/livekit_client.dart';
import 'package:flutter_blue_plus/flutter_blue_plus.dart';

class BrilliantFrameLiveKitService {
  late Room room;
  late LocalVideoTrack cameraTrack;
  late LocalAudioTrack micTrack;
  BluetoothDevice? frameDevice;

  // Track subscriptions
  final Map<String, RemoteVideoTrack> robotCameras = {};

  Future<void> startCollaborationMode() async {
    // 1. Connect to Brilliant Frame via BLE
    await _connectToFrame();

    // 2. Get LiveKit token from Agno backend
    final token = await _getLiveKitToken(
      identity: 'brilliant_frame',
      roomName: 'property_inspection',
    );

    // 3. Connect to LiveKit room
    room = Room();

    // Setup event listeners
    room.on<TrackSubscribedEvent>(_onTrackSubscribed);
    room.on<TrackUnsubscribedEvent>(_onTrackUnsubscribed);
    room.on<ParticipantConnectedEvent>(_onParticipantConnected);

    await room.connect(
      'wss://livekit.onemind.example.com',
      token,
      roomOptions: RoomOptions(
        adaptiveStream: true,
        dynacast: true,
      ),
    );

    // 4. Publish Frame camera
    await _publishFrameCamera();

    // 5. Publish Frame microphone
    await _publishFrameMicrophone();

    print('Collaboration mode active');
  }

  Future<void> _connectToFrame() async {
    // Scan for Brilliant Frame
    final devices = await FlutterBluePlus.scanResults.first;
    frameDevice = devices
        .firstWhere((r) => r.device.name == 'Frame')
        .device;

    await frameDevice!.connect();
    print('Connected to Brilliant Frame');
  }

  Future<String> _getLiveKitToken({
    required String identity,
    required String roomName,
  }) async {
    final response = await agnoClient.post('/api/livekit/token', {
      'identity': identity,
      'room_name': roomName,
    });

    return response.data['token'];
  }

  Future<void> _publishFrameCamera() async {
    // Get camera stream from Frame via BLE
    final frameCameraStream = await _getFrameCameraStream();

    // Create LiveKit video track
    cameraTrack = await LocalVideoTrack.createCameraTrack(
      CameraCaptureOptions(
        cameraPosition: CameraPosition.front,
        params: VideoParametersPresets.h720,
      ),
    );

    // Publish to room
    await room.localParticipant?.publishVideoTrack(cameraTrack);

    print('Frame camera published to LiveKit');
  }

  Future<void> _publishFrameMicrophone() async {
    // Get microphone stream from Frame via BLE
    micTrack = await LocalAudioTrack.create();

    // Publish to room
    await room.localParticipant?.publishAudioTrack(micTrack);

    print('Frame microphone published to LiveKit');
  }

  void _onTrackSubscribed(TrackSubscribedEvent event) {
    print('Subscribed to ${event.participant.identity} ${event.track.kind}');

    if (event.track is RemoteVideoTrack) {
      // This is a robot camera feed
      final robotId = event.participant.identity;
      robotCameras[robotId] = event.track as RemoteVideoTrack;

      // Display in AR overlay
      _displayRobotCameraInAR(robotId, event.track as RemoteVideoTrack);
    }
  }

  void _displayRobotCameraInAR(String robotId, RemoteVideoTrack track) {
    // Send video frames to Frame display via BLE
    track.onVideoFrameReceived = (frame) {
      _sendFrameToAROverlay(robotId, frame);
    };
  }

  Future<void> _sendFrameToAROverlay(String robotId, VideoFrame frame) async {
    // Convert frame to format Frame can display
    final frameData = await frame.toBytes();

    // Send via BLE to Frame
    await frameDevice?.writeCharacteristic(
      displayCharacteristic,
      frameData,
    );
  }

  void _onParticipantConnected(ParticipantConnectedEvent event) {
    print('Robot joined: ${event.participant.name}');

    // Show notification in Frame AR overlay
    _showARNotification('${event.participant.name} joined');
  }

  Future<void> _showARNotification(String message) async {
    await frameDevice?.writeCharacteristic(
      textCharacteristic,
      utf8.encode(message),
    );
  }

  Future<Stream<Uint8List>> _getFrameCameraStream() async {
    // Read camera data from Frame via BLE
    // Returns stream of camera frames
    return frameDevice!
        .characteristicStream(cameraCharacteristic)
        .map((data) => Uint8List.fromList(data));
  }
}
```

---

## Data Flow Diagrams

### Real-Time Streaming Flow

```
┌──────────────────────────────────────────────────────────┐
│ PHASE 1: Stream Setup                                    │
└──────────────────────────────────────────────────────────┘

You activate "Legacy, start collaboration"
    ↓
Flutter App:
    ├─→ Connect Frame (BLE)
    ├─→ Request LiveKit token (Agno API)
    └─→ Join LiveKit room
         ↓
    Publish Frame camera (640x480, 30fps, 500kbps)
         ↓
    Publish Frame mic (48kHz, mono, 64kbps)

Robots (all 5) automatically join:
    ├─→ unitree_go2 publishes camera
    ├─→ unitree_g1 publishes camera
    ├─→ openmind_bot publishes camera
    ├─→ drone publishes camera
    └─→ workshop_bot publishes camera

LiveKit distributes streams:
    ├─→ Your Frame receives 5 robot cameras
    ├─→ All robots receive your Frame camera
    └─→ Agno receives Frame audio + video

Result: Bidirectional streaming active

┌──────────────────────────────────────────────────────────┐
│ PHASE 2: Voice Command Processing                        │
└──────────────────────────────────────────────────────────┘

You: "Legacy, come pick this up"
    ↓ (Frame mic captures audio)
LiveKit audio track
    ↓ (Agno subscribes)
Gemini transcribes
    ↓
Legacy agent receives:
    • Voice: "Legacy, come pick this up"
    • Frame video frame (current view)
    ↓
Gemini Vision analyzes Frame video:
    • Object detected: "Red toolbox"
    • Location: "Indoor workbench"
    ↓
Legacy combines:
    • Command: "pick this up" → red toolbox
    • Destination: "bring to barn"
    ↓
Query robot positions (status APIs)
    ↓
Calculate distances
    ↓
Select optimal robot: unitree_g1
    ↓
Send command: POST /api/unitree/g1/fetch

┌──────────────────────────────────────────────────────────┐
│ PHASE 3: Robot Execution with Visual Feedback            │
└──────────────────────────────────────────────────────────┘

Unitree G1 receives task
    ↓
Subscribes to Frame camera (priority++)
    ↓
Uses your view for visual guidance:
    ├─→ Sees red toolbox from your perspective
    ├─→ Compares with its own camera
    └─→ Confirms target match
         ↓
Navigates to workbench
    ↓ (Publishes camera feed)
Your Frame receives G1 camera
    ↓
Frame AR overlay displays:
    • G1 POV (primary)
    • Other 4 robots (thumbnails)
    • Status: "G1: Approaching target"
         ↓
G1 picks up toolbox
    ↓
Confirms grip (force sensors)
    ↓
Sends status to Agno: "object_secured"
    ↓
Legacy notifies via Frame:
    "Toolbox secured. Heading to barn."
         ↓
G1 navigates to barn (autonomous)
    • SLAM navigation
    • Obstacle avoidance
    • You watch live via G1 camera
         ↓
Delivers toolbox
    ↓
Sends status: "task_complete"
    ↓
Legacy notifies via Frame:
    "Task complete. G1 returning to standby."

Result: Task executed with full transparency
```

---

## Code Examples

### 3. Robot LiveKit Publisher

**File**: `robot_controllers/livekit_publisher.py`

```python
from livekit import rtc
import asyncio
from typing import Optional

class RobotLiveKitPublisher:
    """Publishes robot camera to LiveKit and subscribes to Frame camera"""

    def __init__(self, robot_id: str, robot_name: str, robot_type: str):
        self.robot_id = robot_id
        self.robot_name = robot_name
        self.robot_type = robot_type  # 'unitree', 'openmind', 'drone', etc.
        self.room: Optional[rtc.Room] = None
        self.frame_camera_track: Optional[rtc.RemoteVideoTrack] = None

    async def start(self, room_name: str = "property_inspection"):
        """Start publishing to LiveKit"""

        # 1. Get token from Agno
        token = await self._get_token(room_name)

        # 2. Create and connect room
        self.room = rtc.Room()
        self.room.on("track_subscribed", self._on_track_subscribed)

        await self.room.connect(
            "wss://livekit.onemind.example.com",
            token
        )

        # 3. Publish robot camera
        await self._publish_robot_camera()

        print(f"{self.robot_name} joined LiveKit room")

    async def _get_token(self, room_name: str) -> str:
        """Request LiveKit token from Agno backend"""
        async with aiohttp.ClientSession() as session:
            async with session.post(
                "http://agno-backend/api/livekit/token",
                json={
                    "identity": self.robot_id,
                    "name": self.robot_name,
                    "room_name": room_name
                }
            ) as resp:
                data = await resp.json()
                return data['token']

    async def _publish_robot_camera(self):
        """Publish robot's camera feed to LiveKit"""

        # Get camera feed based on robot type
        if self.robot_type == "unitree":
            camera_source = await self._get_unitree_camera()
        elif self.robot_type == "openmind":
            camera_source = await self._get_openmind_camera()
        elif self.robot_type == "drone":
            camera_source = await self._get_drone_camera()
        else:
            raise ValueError(f"Unknown robot type: {self.robot_type}")

        # Create video track
        video_track = rtc.LocalVideoTrack.create_video_track(
            f"{self.robot_id}_camera",
            camera_source
        )

        # Publish to LiveKit
        options = rtc.TrackPublishOptions()
        options.video_codec = rtc.VideoCodec.VP8

        await self.room.local_participant.publish_track(
            video_track,
            options
        )

        print(f"{self.robot_name} camera published")

    async def _get_unitree_camera(self) -> rtc.VideoSource:
        """Get camera feed from Unitree robot"""
        from unitree_sdk2py import RobotState

        robot = RobotState()
        video_source = rtc.VideoSource(640, 480)

        async def capture_frames():
            while True:
                # Get frame from Unitree SDK
                frame = robot.GetCameraFrame()

                # Convert to LiveKit video frame
                lk_frame = rtc.VideoFrame(
                    width=640,
                    height=480,
                    type=rtc.VideoBufferType.RGBA,
                    data=frame
                )

                video_source.capture_frame(lk_frame)
                await asyncio.sleep(1/30)  # 30 FPS

        asyncio.create_task(capture_frames())
        return video_source

    async def _get_openmind_camera(self) -> rtc.VideoSource:
        """Get camera feed from OpenMind robot via OM1"""
        # Similar to Unitree, but using OM1 SDK
        pass

    async def _get_drone_camera(self) -> rtc.VideoSource:
        """Get camera feed from drone"""
        # Using drone SDK (e.g., DJI)
        pass

    async def _on_track_subscribed(
        self,
        track: rtc.Track,
        publication: rtc.TrackPublication,
        participant: rtc.RemoteParticipant
    ):
        """Handle subscribed tracks"""

        if participant.identity == "brilliant_frame":
            if track.kind == rtc.TrackKind.KIND_VIDEO:
                # This is the user's Frame camera
                self.frame_camera_track = track
                print(f"{self.robot_name} now sees Frame camera")

                # Start using Frame view for navigation
                asyncio.create_task(self._process_frame_view())

    async def _process_frame_view(self):
        """Process Frame camera view for visual guidance"""

        if not self.frame_camera_track:
            return

        async for frame in rtc.VideoStream(self.frame_camera_track):
            # Convert frame to bytes
            frame_bytes = await frame.to_argb_bytes()

            # Send to robot's vision system
            await self._analyze_frame_view(frame_bytes)

    async def _analyze_frame_view(self, frame_bytes: bytes):
        """Analyze what user is looking at"""

        # Robot can use this to:
        # 1. Identify target objects
        # 2. Understand user's intent
        # 3. Navigate using shared visual context

        # Example: Detect objects in Frame view
        # Then compare with robot's own camera view
        pass
```

---

## Deployment

### Infrastructure Setup

```yaml
# docker-compose.yml

version: '3.8'

services:
  # Self-hosted LiveKit server
  livekit:
    image: livekit/livekit-server:latest
    command: --config /etc/livekit.yaml
    volumes:
      - ./livekit-config.yaml:/etc/livekit.yaml
    ports:
      - "7880:7880"   # HTTP
      - "7881:7881"   # WebRTC TCP
      - "7882:7882/udp"  # WebRTC UDP
      - "50000-60000:50000-60000/udp"  # RTP
    environment:
      - LIVEKIT_KEYS=APIXXX:secretXXX
    restart: unless-stopped

  # Agno backend
  agno:
    build: ./backend
    ports:
      - "8000:8000"
    environment:
      - LIVEKIT_URL=ws://livekit:7880
      - LIVEKIT_API_KEY=APIXXX
      - LIVEKIT_API_SECRET=secretXXX
    depends_on:
      - livekit
    restart: unless-stopped

  # Robot controllers (one per robot type)
  unitree-controller:
    build: ./robot_controllers/unitree
    environment:
      - LIVEKIT_URL=ws://livekit:7880
      - AGNO_API_URL=http://agno:8000
    network_mode: host  # For local robot communication
    restart: unless-stopped

  openmind-controller:
    build: ./robot_controllers/openmind
    environment:
      - LIVEKIT_URL=ws://livekit:7880
      - AGNO_API_URL=http://agno:8000
    network_mode: host
    restart: unless-stopped
```

**LiveKit Configuration**:

```yaml
# livekit-config.yaml

port: 7880
rtc:
  port_range_start: 50000
  port_range_end: 60000
  use_external_ip: true

keys:
  APIXXX: secretXXX

room:
  auto_create: true
  empty_timeout: 300
  max_participants: 10

video:
  codec: vp8

audio:
  codec: opus
```

---

## Cost Analysis

### Self-Hosted LiveKit (Recommended)

**Server Requirements**:
- CPU: 4 cores
- RAM: 8GB
- Bandwidth: 10 Mbps upload (for 6 participants)
- Storage: 50GB

**Monthly Cost**:

| Provider | Instance | Cost |
|----------|----------|------|
| **Hetzner** | CX31 | €9.50 ($10) |
| **DigitalOcean** | 4GB Droplet | $24 |
| **AWS Lightsail** | 8GB instance | $40 |

**Recommendation**: Hetzner CX31 ✅ ($10/month)

### LiveKit Cloud (Alternative)

**Pricing**: $0.04/minute per participant

**Your scenario** (6 participants):
- Cost per minute: $0.24
- 1 hour session: $14.40
- 10 hours/month: $144

**Verdict**: Self-hosted is 14x cheaper ✅

---

## Performance Metrics

### Bandwidth Requirements

| Stream | Resolution | FPS | Bitrate | Bandwidth |
|--------|-----------|-----|---------|-----------|
| Frame camera (out) | 640x480 | 30 | 500 kbps | 0.5 Mbps |
| Robot 1 camera (in) | 640x480 | 30 | 500 kbps | 0.5 Mbps |
| Robot 2 camera (in) | 640x480 | 30 | 500 kbps | 0.5 Mbps |
| Robot 3 camera (in) | 640x480 | 30 | 500 kbps | 0.5 Mbps |
| Robot 4 camera (in) | 640x480 | 30 | 500 kbps | 0.5 Mbps |
| Robot 5 camera (in) | 640x480 | 30 | 500 kbps | 0.5 Mbps |
| Frame mic (out) | 48kHz | - | 64 kbps | 0.064 Mbps |
| **Total upstream** | - | - | - | **0.564 Mbps** |
| **Total downstream** | - | - | - | **2.5 Mbps** |
| **Total bandwidth** | - | - | - | **3.064 Mbps** |

✅ **Works on 4G** (typical: 5-50 Mbps upload, 10-100 Mbps download)
✅ **Works on 5G** (typical: 50-300 Mbps upload)

### Latency

| Component | Latency |
|-----------|---------|
| Frame camera → LiveKit | 50-100ms |
| LiveKit → Robot camera | 50-100ms |
| **Total glass-to-glass** | **100-200ms** |
| Voice command processing | 300-500ms |
| Robot response time | 500-1000ms |

**User experience**: Near real-time ✅

---

## Summary

### What This Enables

✅ **Real-time multi-robot coordination**
✅ **Shared visual context** (all see what you see)
✅ **Voice control with visual understanding**
✅ **Live monitoring** (watch robots in AR overlay)
✅ **Bidirectional streaming** (you see robots, robots see you)

### Tech Stack Integration

```
Flutter (Frame UI)
  ↔ LiveKit (WebRTC streaming)
  ↔ Agno (intelligence + coordination)
  ↔ Robot APIs (control)
```

### Cost: $10-40/month

**Bandwidth**: 3 Mbps (works on 4G/5G)

**Latency**: 100-200ms (near real-time)

---

**File**: `/Users/zeuslegacy/Documents/onemind-workspace/LIVEKIT_MULTI_ROBOT_STREAMING.md`
**Status**: Complete multi-robot streaming integration guide
**Use**: Reference for implementing LiveKit-based robot collaboration


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]
