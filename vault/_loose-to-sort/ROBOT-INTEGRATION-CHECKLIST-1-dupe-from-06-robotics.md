# Robot Integration: Complete Implementation Checklist

> **Status:** Ready to Build (Software Layer Prepared)
> **Timeline:** 4-6 weeks when ready to implement
> **Cost:** $997-$3,356 hardware + $135/mo services
> **Last Updated:** 2026-01-22

---

## Prerequisites (Hardware - When Ready)

- [ ] **OpenMind Edge Compute Device** - $997 (or DIY Jetson $500)
- [ ] **Unitree Robot** - You already have
- [ ] **Optional: DJI Mini 4 Pro Drone** - $759
- [ ] **Optional: Second Unitree Robot** - $1,600-$3,500
- [ ] **Network:** Tailscale VPN configured
- [ ] **Budget Approval:** $135/mo for services (AWS, Deepgram, Cartesia, FABRIC)

---

## Software Implementation Phases

### Phase 1: RemoteAgent Base (Week 1-2)

#### 1.1 Create Directory Structure
```bash
mkdir -p agents/edge/{tools,config}
touch agents/edge/__init__.py
touch agents/edge/remote_agent.py
touch agents/edge/config.py
touch agents/edge/requirements.txt
```

**Files to Create:**

**`agents/edge/requirements.txt`:**
```
agno>=0.1.0
nats-py>=2.7.0
asyncio-mqtt>=0.16.0
python-dotenv>=1.0.0
pydantic>=2.5.0
```

**`agents/edge/config.py`:**
```python
import os
from dataclasses import dataclass

@dataclass
class EdgeConfig:
    location: str = os.getenv("EDGE_LOCATION", "barn")
    nats_hub_url: str = os.getenv("NATS_HUB_URL", "nats://100.102.21.44:4222")
    nats_leaf_port: int = int(os.getenv("NATS_LEAF_PORT", "7422"))
    mqtt_port: int = int(os.getenv("MQTT_PORT", "1883"))

    # OM1 Integration
    om1_api_url: str = os.getenv("OM1_API_URL", "http://localhost:8090")
    om1_api_key: str = os.getenv("OM1_API_KEY", "")

    # Robot specific
    robot_type: str = os.getenv("ROBOT_TYPE", "unitree_go1")
    robot_id: str = os.getenv("ROBOT_ID", "go1_001")
```

**`agents/edge/remote_agent.py`:**
```python
#!/usr/bin/env python3
"""
RemoteAgent - Edge agent for OneMind OS robot control.
Connects to NATS hub and receives commands from Legacy.
"""

import asyncio
import json
import nats
from datetime import datetime
from typing import Dict, Any, Callable
from config import EdgeConfig

class RemoteAgent:
    """Edge agent that receives commands from Legacy via NATS"""

    def __init__(self, config: EdgeConfig):
        self.config = config
        self.nc = None
        self.tools: Dict[str, Callable] = {}

    async def connect(self):
        """Connect to NATS hub and subscribe to commands"""
        self.nc = await nats.connect(self.config.nats_hub_url)

        # Subscribe to commands for this location
        await self.nc.subscribe(
            f"commands.{self.config.location}.>",
            cb=self._handle_command
        )

        # Announce presence
        await self.nc.publish("agents.online", json.dumps({
            "location": self.config.location,
            "robot_id": self.config.robot_id,
            "robot_type": self.config.robot_type,
            "capabilities": list(self.tools.keys()),
            "timestamp": datetime.utcnow().isoformat()
        }).encode())

        print(f"✅ RemoteAgent online: {self.config.location}/{self.config.robot_id}")

    async def _handle_command(self, msg):
        """Handle incoming commands from Legacy"""
        try:
            subject = msg.subject  # e.g., "commands.barn.robot.move"
            parts = subject.split(".")
            tool_name = parts[-1] if len(parts) > 2 else "unknown"

            data = json.loads(msg.data.decode())
            print(f"📥 Command received: {tool_name} - {data}")

            if tool_name in self.tools:
                result = await self.tools[tool_name](data)
                await msg.respond(json.dumps({
                    "success": True,
                    "data": result,
                    "timestamp": datetime.utcnow().isoformat()
                }).encode())
            else:
                await msg.respond(json.dumps({
                    "success": False,
                    "error": f"Unknown tool: {tool_name}",
                    "available_tools": list(self.tools.keys())
                }).encode())

        except Exception as e:
            print(f"❌ Error handling command: {e}")
            await msg.respond(json.dumps({
                "success": False,
                "error": str(e)
            }).encode())

    def register_tool(self, name: str, func: Callable):
        """Register a tool/capability"""
        self.tools[name] = func
        print(f"🔧 Registered tool: {name}")

    async def publish_event(self, event_type: str, data: Dict[str, Any]):
        """Publish an event back to Legacy"""
        await self.nc.publish(
            f"events.{self.config.location}.{event_type}",
            json.dumps({
                "location": self.config.location,
                "robot_id": self.config.robot_id,
                "timestamp": datetime.utcnow().isoformat(),
                **data
            }).encode()
        )

    async def heartbeat(self, interval: int = 30):
        """Send periodic heartbeat to cloud"""
        while True:
            await asyncio.sleep(interval)
            await self.publish_event("heartbeat", {
                "status": "online",
                "tools": list(self.tools.keys())
            })

# Entry point
async def main():
    config = EdgeConfig()
    agent = RemoteAgent(config)

    # TODO: Register tools here
    # agent.register_tool("move", move_robot)
    # agent.register_tool("status", get_robot_status)

    await agent.connect()

    # Keep running
    await agent.heartbeat()

if __name__ == "__main__":
    asyncio.run(main())
```

**Tasks:**
- [ ] Create all files above
- [ ] Test NATS connection to cloud hub
- [ ] Deploy to test device (Pi/Jetson)
- [ ] Verify "agents.online" event received in cloud

---

#### 1.2 Deploy NATS Leaf Node

**`infra/compose/edge/docker-compose.yml`:**
```yaml
version: '3.8'

services:
  nats-leaf:
    image: nats:2.10-alpine
    container_name: onemind-nats-leaf
    ports:
      - "4222:4222"  # Client connections
      - "7422:7422"  # Leaf node port
      - "8222:8222"  # Monitoring
    volumes:
      - ./nats-leaf.conf:/etc/nats/nats-server.conf
    command: ["-c", "/etc/nats/nats-server.conf"]
    restart: unless-stopped
    networks:
      - onemind-edge

  remote-agent:
    build:
      context: ../../../
      dockerfile: infra/containers/Dockerfile.remote-agent
    container_name: onemind-remote-agent
    environment:
      - EDGE_LOCATION=barn
      - NATS_HUB_URL=nats://nats-leaf:4222
      - ROBOT_TYPE=unitree_go1
      - ROBOT_ID=go1_001
      - OM1_API_URL=http://om1:8090
    depends_on:
      - nats-leaf
    restart: unless-stopped
    networks:
      - onemind-edge
    privileged: true  # For hardware access

networks:
  onemind-edge:
    driver: bridge
```

**`infra/compose/edge/nats-leaf.conf`:**
```
port: 4222
http_port: 8222

leafnodes {
  remotes = [
    {
      url: "nats://100.102.21.44:7422"  # Cloud hub
    }
  ]
}
```

**`infra/containers/Dockerfile.remote-agent`:**
```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY agents/edge/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY agents/edge/ .

CMD ["python", "remote_agent.py"]
```

**Tasks:**
- [ ] Create edge docker-compose.yml
- [ ] Configure NATS leaf node
- [ ] Build RemoteAgent container
- [ ] Test leaf → hub connection
- [ ] Verify bidirectional messaging

---

### Phase 2: Unitree Control Tools (Week 2-3)

#### 2.1 Unitree SDK Integration

**`agents/edge/tools/unitree.py`:**
```python
"""
Unitree Robot Control Tools
Requires: unitree_legged_sdk (pip install unitree-legged-sdk)
"""

import asyncio
from typing import Dict, Any

# Placeholder - actual SDK integration needed
class UnitreeController:
    """Controller for Unitree Go1/Go2 robots"""

    def __init__(self, robot_type: str = "go1"):
        self.robot_type = robot_type
        self.connected = False

    async def connect(self):
        """Connect to robot via UDP"""
        # TODO: Initialize Unitree SDK connection
        # from unitree_legged_sdk import RobotInterface
        # self.robot = RobotInterface()
        self.connected = True
        print(f"✅ Connected to {self.robot_type}")

    async def move(self, direction: str, distance: float) -> Dict[str, Any]:
        """Move robot in specified direction"""
        if not self.connected:
            raise RuntimeError("Robot not connected")

        # TODO: Implement via Unitree SDK
        # if direction == "forward":
        #     await self.robot.walk_forward(distance)

        return {
            "action": "move",
            "direction": direction,
            "distance": distance,
            "status": "completed"
        }

    async def navigate_to(self, lat: float, lon: float) -> Dict[str, Any]:
        """Navigate to GPS coordinates"""
        # TODO: Implement waypoint navigation
        return {
            "action": "navigate",
            "target": {"lat": lat, "lon": lon},
            "status": "en_route"
        }

    async def get_status(self) -> Dict[str, Any]:
        """Get robot status"""
        # TODO: Get from Unitree SDK
        return {
            "battery": 78,
            "mode": "standing",
            "location": {"lat": 0, "lon": 0},
            "speed": 0,
            "connected": self.connected
        }

    async def return_to_dock(self) -> Dict[str, Any]:
        """Return to charging dock"""
        # TODO: Navigate to dock location
        return {
            "action": "return_to_dock",
            "status": "en_route"
        }

# Tool functions for RemoteAgent
async def move_robot(data: Dict[str, Any]) -> Dict[str, Any]:
    """Tool: Move robot"""
    controller = UnitreeController()
    await controller.connect()
    return await controller.move(
        direction=data.get("direction", "forward"),
        distance=data.get("distance", 1.0)
    )

async def get_robot_status(data: Dict[str, Any]) -> Dict[str, Any]:
    """Tool: Get robot status"""
    controller = UnitreeController()
    await controller.connect()
    return await controller.get_status()

async def navigate_robot(data: Dict[str, Any]) -> Dict[str, Any]:
    """Tool: Navigate to coordinates"""
    controller = UnitreeController()
    await controller.connect()
    return await controller.navigate_to(
        lat=data.get("lat"),
        lon=data.get("lon")
    )
```

**Tasks:**
- [ ] Install Unitree SDK on edge device
- [ ] Test UDP connection to robot
- [ ] Implement move commands
- [ ] Implement status queries
- [ ] Test waypoint navigation

---

#### 2.2 OpenMind OM1 Integration

**`agents/edge/tools/om1.py`:**
```python
"""
OpenMind OM1 API Integration
Provides high-level robot control via OM1 platform
"""

import aiohttp
from typing import Dict, Any

class OM1Client:
    """Client for OpenMind OM1 API"""

    def __init__(self, base_url: str, api_key: str):
        self.base_url = base_url
        self.api_key = api_key
        self.session = None

    async def __aenter__(self):
        self.session = aiohttp.ClientSession(
            headers={"Authorization": f"Bearer {self.api_key}"}
        )
        return self

    async def __aexit__(self, *args):
        await self.session.close()

    async def send_command(self, robot_id: str, command: Dict[str, Any]) -> Dict[str, Any]:
        """Send high-level command to robot via OM1"""
        async with self.session.post(
            f"{self.base_url}/api/robots/{robot_id}/commands",
            json=command
        ) as resp:
            return await resp.json()

    async def get_telemetry(self, robot_id: str) -> Dict[str, Any]:
        """Get robot telemetry from OM1"""
        async with self.session.get(
            f"{self.base_url}/api/robots/{robot_id}/telemetry"
        ) as resp:
            return await resp.json()

    async def stream_video(self, robot_id: str, camera: str = "front") -> str:
        """Get video stream URL from robot camera"""
        async with self.session.get(
            f"{self.base_url}/api/robots/{robot_id}/cameras/{camera}/stream"
        ) as resp:
            data = await resp.json()
            return data.get("stream_url")

# Tool functions
async def om1_command(data: Dict[str, Any]) -> Dict[str, Any]:
    """Send command via OM1"""
    async with OM1Client(
        base_url=data.get("om1_url"),
        api_key=data.get("om1_key")
    ) as client:
        return await client.send_command(
            robot_id=data.get("robot_id"),
            command=data.get("command")
        )
```

**Tasks:**
- [ ] Install OM1 on edge device ($997 device)
- [ ] Configure OM1 API credentials
- [ ] Test OM1 → Unitree connection
- [ ] Implement command forwarding
- [ ] Test video streaming to LiveKit

---

### Phase 3: Cloud Integration (Week 3-4)

#### 3.1 Legacy Agent Tools

**`agents/legacy/tools/robot_control.py`:**
```python
"""
Robot Control Tools for Legacy Agent
Sends commands to edge robots via NATS
"""

from agno.tools import tool
import nats
import json
from typing import Dict, Any

NATS_URL = "nats://localhost:4222"

@tool
async def move_robot(
    location: str,
    direction: str,
    distance: float = 1.0,
    robot_id: str = "go1_001"
) -> str:
    """
    Move robot at specified location.

    Args:
        location: Location name (barn, greenhouse, home)
        direction: Movement direction (forward, backward, left, right)
        distance: Distance in meters
        robot_id: Robot identifier

    Returns:
        Status message
    """
    nc = await nats.connect(NATS_URL)

    response = await nc.request(
        f"commands.{location}.robot.move",
        json.dumps({
            "robot_id": robot_id,
            "direction": direction,
            "distance": distance
        }).encode(),
        timeout=5.0
    )

    result = json.loads(response.data.decode())
    await nc.close()

    if result.get("success"):
        return f"Robot moved {direction} {distance}m at {location}"
    else:
        return f"Failed to move robot: {result.get('error')}"

@tool
async def get_robot_status(
    location: str,
    robot_id: str = "go1_001"
) -> Dict[str, Any]:
    """
    Get robot status at location.

    Args:
        location: Location name
        robot_id: Robot identifier

    Returns:
        Robot status (battery, location, mode)
    """
    nc = await nats.connect(NATS_URL)

    response = await nc.request(
        f"commands.{location}.robot.status",
        json.dumps({"robot_id": robot_id}).encode(),
        timeout=5.0
    )

    result = json.loads(response.data.decode())
    await nc.close()

    return result.get("data", {})

@tool
async def dispatch_robot_patrol(
    location: str,
    area: str,
    duration_minutes: int = 30,
    robot_id: str = "go1_001"
) -> str:
    """
    Dispatch robot for patrol mission.

    Args:
        location: Location name
        area: Area to patrol (perimeter, greenhouse, barn)
        duration_minutes: Patrol duration
        robot_id: Robot identifier

    Returns:
        Mission status
    """
    nc = await nats.connect(NATS_URL)

    response = await nc.request(
        f"commands.{location}.robot.patrol",
        json.dumps({
            "robot_id": robot_id,
            "area": area,
            "duration": duration_minutes
        }).encode(),
        timeout=5.0
    )

    result = json.loads(response.data.decode())
    await nc.close()

    if result.get("success"):
        return f"Robot dispatched for {area} patrol at {location}"
    else:
        return f"Failed to dispatch robot: {result.get('error')}"

@tool
async def get_robot_video_stream(
    location: str,
    robot_id: str = "go1_001",
    camera: str = "front"
) -> str:
    """
    Get live video stream URL from robot camera.

    Args:
        location: Location name
        robot_id: Robot identifier
        camera: Camera name (front, rear, left, right, bottom)

    Returns:
        LiveKit stream URL
    """
    nc = await nats.connect(NATS_URL)

    response = await nc.request(
        f"commands.{location}.robot.video",
        json.dumps({
            "robot_id": robot_id,
            "camera": camera
        }).encode(),
        timeout=5.0
    )

    result = json.loads(response.data.decode())
    await nc.close()

    if result.get("success"):
        stream_url = result.get("data", {}).get("stream_url")
        return f"Video stream: {stream_url}"
    else:
        return f"Failed to get video: {result.get('error')}"
```

**Tasks:**
- [ ] Add robot tools to Legacy agent
- [ ] Test NATS request/response
- [ ] Implement error handling
- [ ] Add tool descriptions
- [ ] Test voice commands

---

#### 3.2 Update Legacy Agent

**`agents/legacy/agent.py` (additions):**
```python
# Add to imports
from .tools.robot_control import (
    move_robot,
    get_robot_status,
    dispatch_robot_patrol,
    get_robot_video_stream
)

# Add to Legacy agent tools list
legacy = Agent(
    name="Legacy",
    # ... existing config ...
    tools=[
        # ... existing tools ...
        move_robot,
        get_robot_status,
        dispatch_robot_patrol,
        get_robot_video_stream,
    ],
)
```

**Tasks:**
- [ ] Import robot tools
- [ ] Add to agent tools list
- [ ] Update agent instructions (mention robot capabilities)
- [ ] Test agent initialization
- [ ] Test tool discovery

---

### Phase 4: Video Streaming (Week 4-5)

#### 4.1 Robot Camera → LiveKit

**`agents/edge/tools/video_stream.py`:**
```python
"""
Video Streaming: Robot cameras → LiveKit
"""

import asyncio
from livekit import rtc

async def stream_camera_to_livekit(
    robot_id: str,
    camera_name: str,
    livekit_url: str,
    livekit_token: str
) -> str:
    """
    Stream robot camera to LiveKit room.

    Returns:
        LiveKit room URL
    """
    room = rtc.Room()

    await room.connect(
        url=livekit_url,
        token=livekit_token
    )

    # TODO: Capture video from robot camera
    # video_source = get_robot_camera_feed(robot_id, camera_name)

    # Publish to LiveKit
    # track = rtc.LocalVideoTrack.create_video_track(
    #     name=f"{robot_id}_{camera_name}",
    #     source=video_source
    # )
    # await room.local_participant.publish_track(track)

    room_url = f"{livekit_url}/room/{room.name}"
    return room_url
```

**Tasks:**
- [ ] Capture robot camera feed (Unitree SDK)
- [ ] Create LiveKit video track
- [ ] Publish to LiveKit room
- [ ] Test video streaming
- [ ] Add to UI dashboard

---

#### 4.2 UI Integration

**`ui/src/components/robots/RobotDashboard.tsx`:**
```typescript
// Robot dashboard with video streams
import { LiveKitRoom, VideoTrack } from '@livekit/components-react';

export function RobotDashboard() {
  const [robots, setRobots] = useState([]);

  useEffect(() => {
    // Fetch robot list from API
    fetch('/api/robots')
      .then(res => res.json())
      .then(data => setRobots(data));
  }, []);

  return (
    <div className="robot-dashboard">
      {robots.map(robot => (
        <RobotCard key={robot.id} robot={robot} />
      ))}
    </div>
  );
}

function RobotCard({ robot }) {
  return (
    <div className="robot-card">
      <h3>{robot.id} - {robot.location}</h3>
      <div className="robot-status">
        <span>Battery: {robot.battery}%</span>
        <span>Status: {robot.status}</span>
      </div>

      {/* Live video stream */}
      <LiveKitRoom
        serverUrl={robot.livekit_url}
        token={robot.livekit_token}
      >
        <VideoTrack />
      </LiveKitRoom>

      <div className="robot-controls">
        <button onClick={() => moveRobot(robot.id, 'forward')}>
          Forward
        </button>
        <button onClick={() => getRobotStatus(robot.id)}>
          Status
        </button>
      </div>
    </div>
  );
}
```

**Tasks:**
- [ ] Create robot dashboard UI
- [ ] Add LiveKit video components
- [ ] Implement robot controls
- [ ] Test multi-robot view
- [ ] Add to Flutter app

---

### Phase 5: Automated Workflows (Week 5-6)

#### 5.1 Morning Patrol Workflow

**`agents/legacy/workflows/morning_patrol.py`:**
```python
"""
Morning Patrol: Automated daily robot inspection
"""

from agno.workflows import Workflow
from ..tools.robot_control import dispatch_robot_patrol, get_robot_status

morning_patrol = Workflow(
    name="Morning Patrol",
    description="Automated morning inspection with robot",
    schedule="0 7 * * *",  # 7:00 AM daily
)

@morning_patrol.step("check_battery")
async def check_battery():
    """Ensure robot has enough battery"""
    status = await get_robot_status(location="barn", robot_id="go1_001")

    if status.get("battery", 0) < 30:
        return {
            "status": "skip",
            "reason": f"Low battery: {status['battery']}%"
        }

    return {"status": "continue"}

@morning_patrol.step("dispatch_patrol")
async def dispatch_patrol():
    """Send robot on morning patrol"""
    result = await dispatch_robot_patrol(
        location="barn",
        area="greenhouse,coop,perimeter",
        duration_minutes=45,
        robot_id="go1_001"
    )

    return {"status": "patrolling", "message": result}

@morning_patrol.step("report")
async def generate_report():
    """Generate morning report for user"""
    # TODO: Analyze patrol data
    # TODO: Send voice summary to user

    return {
        "status": "completed",
        "report": "Morning patrol complete"
    }
```

**Tasks:**
- [ ] Create workflow definitions
- [ ] Test scheduled execution
- [ ] Add error handling
- [ ] Implement reporting
- [ ] Add to Legacy workflows

---

### Phase 6: Testing & Deployment (Week 6)

#### 6.1 End-to-End Testing

**Test Checklist:**
- [ ] **Voice Command Test**
  - Say: "Legacy, move robot forward 2 meters"
  - Verify: Robot moves, confirms via voice

- [ ] **Status Query Test**
  - Say: "Legacy, what's the robot battery level?"
  - Verify: Reports accurate battery percentage

- [ ] **Patrol Test**
  - Say: "Legacy, have the robot patrol the greenhouse"
  - Verify: Robot navigates, returns, reports findings

- [ ] **Video Stream Test**
  - Say: "Legacy, show me the robot camera"
  - Verify: Video appears in UI, <200ms latency

- [ ] **Multi-Robot Test** (if 2+ robots)
  - Say: "Legacy, have both robots patrol different areas"
  - Verify: Coordinated execution via FABRIC

- [ ] **Emergency Stop Test**
  - Say: "Legacy, stop the robot immediately"
  - Verify: Robot stops within 1 second

- [ ] **Return to Dock Test**
  - Say: "Legacy, send robot back to charging dock"
  - Verify: Navigates to dock, connects to charger

---

#### 6.2 Deployment Scripts

**`scripts/deploy_edge_agent.sh`:**
```bash
#!/bin/bash
# Deploy RemoteAgent to edge device

set -e

EDGE_HOST=${1:-"barn-edge"}
EDGE_USER=${2:-"ubuntu"}

echo "🚀 Deploying RemoteAgent to $EDGE_HOST..."

# Copy files to edge device
scp -r agents/edge/* ${EDGE_USER}@${EDGE_HOST}:/opt/onemind/edge/
scp infra/compose/edge/docker-compose.yml ${EDGE_USER}@${EDGE_HOST}:/opt/onemind/

# Deploy via Docker Compose
ssh ${EDGE_USER}@${EDGE_HOST} << 'EOF'
cd /opt/onemind
docker-compose down
docker-compose up -d --build
docker-compose logs -f
EOF

echo "✅ Deployment complete!"
```

**Tasks:**
- [ ] Create deployment scripts
- [ ] Test on staging edge device
- [ ] Document deployment process
- [ ] Create rollback procedure
- [ ] Deploy to production

---

## API Endpoints (To Add)

### Robot Control API

**`agents/legacy/main.py` (additions):**
```python
# GET /api/robots - List all robots
@app.get("/api/robots")
async def list_robots():
    """List all connected robots"""
    # Query NATS for online robots
    return {"robots": [...]}

# GET /api/robots/{robot_id}/status
@app.get("/api/robots/{robot_id}/status")
async def get_robot_status(robot_id: str):
    """Get robot status"""
    # Query via NATS
    return {"robot_id": robot_id, "status": {...}}

# POST /api/robots/{robot_id}/command
@app.post("/api/robots/{robot_id}/command")
async def send_robot_command(robot_id: str, command: dict):
    """Send command to robot"""
    # Send via NATS
    return {"success": True, "command": command}

# GET /api/robots/{robot_id}/video/{camera}
@app.get("/api/robots/{robot_id}/video/{camera}")
async def get_robot_video(robot_id: str, camera: str):
    """Get robot video stream URL"""
    # Get LiveKit stream URL
    return {"stream_url": "..."}
```

---

## Documentation Cross-References

### Main Documents
- **Integration Overview:** `docs/planning/UNITREE-OPENMIND-INTEGRATION.md`
- **Fleet Costs:** `docs/planning/MULTI-ROBOT-FLEET-COSTS.md`
- **RemoteAgent Architecture:** `docs/architecture/REMOTEAGENT-ARCHITECTURE.md`
- **NATS Multi-Location:** `docs/architecture/MULTI-LOCATION-NATS.md`
- **This Checklist:** `docs/planning/ROBOT-INTEGRATION-CHECKLIST.md`

### Related Docs
- OneMind Architecture: `docs/architecture/ONEMIND-ARCHITECTURE-PLAN.md`
- Omnipresent Mode: `docs/architecture/OMNIPRESENT-ARCHITECTURE.md`
- LiveKit Voice: `docs/legacy/VOICE_INTEGRATION.md`

---

## Linear Issue: OMOS-???

**Title:** Integrate Unitree Robot with OneMind OS via OpenMind OM1

**Status:** Ready to Implement (Hardware pending)

**Labels:** robotics, integration, high-priority, 4-6-weeks

**Description:** See `UNITREE-OPENMIND-INTEGRATION.md` for full details

**Assignee:** Claude Code (when ready)

**Estimate:** 4-6 weeks

---

## Quick Start (When Ready)

```bash
# 1. Create edge directory structure
cd ~/onemind-os
mkdir -p agents/edge/{tools,config}

# 2. Copy checklist templates
# (All code snippets from this document)

# 3. Install on edge device
ssh edge-device
cd /opt/onemind
git clone https://github.com/yourusername/onemind-os.git
cd onemind-os/infra/compose/edge
docker-compose up -d

# 4. Test connection
python3 -c "import nats; print('NATS OK')"

# 5. Deploy RemoteAgent
./scripts/deploy_edge_agent.sh barn-edge

# 6. Test voice command
# Say: "Legacy, what's the robot status?"
```

---

## Budget Summary (When Ready to Purchase)

### Hardware (One-Time)
- OpenMind Edge Compute: $997
- DJI Mini 4 Pro (optional): $759
- Unitree Go2 (optional): $1,600
- **TOTAL:** $997-$3,356

### Services (Monthly)
- AWS Bedrock: $50
- Deepgram STT: $20
- Cartesia TTS: $10
- OpenMind FABRIC: $50
- Tailscale Pro: $5
- **TOTAL:** $135/mo

### Development
- Claude Code: $20/mo (4-6 weeks)

**First Year Total:** $2,617-$5,976 (depending on hardware)

---

## Success Criteria

### MVP (Minimum Viable Product)
- [ ] Voice: "Move robot forward" works
- [ ] Robot status query works
- [ ] Video stream appears in UI
- [ ] Battery monitoring works
- [ ] Emergency stop works

### V1 (Production Ready)
- [ ] Autonomous patrol routes
- [ ] Return to dock
- [ ] Multi-robot coordination
- [ ] Automated workflows (morning patrol)
- [ ] Full UI integration

### V2 (Advanced)
- [ ] Multi-location fleet
- [ ] Drone integration
- [ ] Computer vision analysis
- [ ] Predictive maintenance
- [ ] Fleet analytics dashboard

---

**Status:** ✅ Software layer prepared and ready to build
**Next Action:** Purchase hardware when ready, then start Phase 1
**Timeline:** 4-6 weeks from start to production

*Last Updated: 2026-01-22*
*Prepared for: Implementation when hardware arrives*
