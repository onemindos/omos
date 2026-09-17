# OneMind OS - Edge Deployment Guide

## Overview

Deploy Legacy presence to edge locations using Jetson Nano (or similar ARM devices). Each edge operates independently with offline capability.

## Hardware Requirements

### Recommended: Jetson Nano 4GB
- **CPU**: Quad-core ARM A57
- **GPU**: 128 CUDA cores (for local AI)
- **RAM**: 4GB
- **Storage**: 64GB+ microSD or NVMe
- **Cost**: ~$150

### Alternatives

| Device | RAM | GPU | Best For |
|--------|-----|-----|----------|
| Jetson Nano 4GB | 4GB | 128 CUDA | Full edge (recommended) |
| Jetson Orin Nano | 8GB | 1024 CUDA | Heavy AI workloads |
| Raspberry Pi 5 | 8GB | None | Light edge (no local LLM) |
| Orange Pi 5 | 8GB | NPU | Budget option |

## Software Stack

```
┌─────────────────────────────────────────────────┐
│                 JETSON EDGE                      │
├─────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────┐ │
│  │   Ollama    │  │   Frigate   │  │  NATS   │ │
│  │ (Local LLM) │  │  (Cameras)  │  │  Leaf   │ │
│  └─────────────┘  └─────────────┘  └─────────┘ │
│         │                │              │       │
│         └────────────────┴──────────────┘       │
│                         │                       │
│              ┌──────────▼──────────┐            │
│              │    RemoteAgent      │            │
│              │  (Edge Controller)  │            │
│              └─────────────────────┘            │
└─────────────────────────────────────────────────┘
                         │
                    NATS Leaf
                         │
                         ▼
                   Cloud NATS Hub
```

## Offline LLM Models

### For Jetson Nano 4GB (2-3GB available for model)

| Model | Size | Speed | Quality | Use Case |
|-------|------|-------|---------|----------|
| **Phi-3 Mini 4K** | 2.3GB | Fast | Good | General tasks |
| **TinyLlama 1.1B** | 0.6GB | Very fast | Basic | Simple commands |
| **Llama 3.2 3B Q4** | 2.0GB | Medium | Good | Complex reasoning |
| **Gemma 2B** | 1.4GB | Fast | Good | Balanced |

### For Jetson Orin Nano 8GB

| Model | Size | Speed | Quality | Use Case |
|-------|------|-------|---------|----------|
| **Llama 3.1 8B Q4** | 4.7GB | Medium | Excellent | Full capability |
| **Mistral 7B Q4** | 4.1GB | Medium | Excellent | Best quality |
| **CodeLlama 7B** | 4.0GB | Medium | Excellent | Code tasks |

### Recommended: Phi-3 Mini for Jetson Nano

```bash
# Install Ollama on Jetson
curl -fsSL https://ollama.com/install.sh | sh

# Pull Phi-3 Mini (best balance for 4GB)
ollama pull phi3:mini

# Test
ollama run phi3:mini "What time is it in Tokyo?"
```

## Docker Compose for Edge

```yaml
# infra/compose/edge/docker-compose.yml
version: "3.8"

services:
  # Local LLM
  ollama:
    image: ollama/ollama
    container_name: ollama
    restart: unless-stopped
    volumes:
      - ollama_data:/root/.ollama
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: all
              capabilities: [gpu]

  # Camera AI (if cameras connected)
  frigate:
    image: ghcr.io/blakeblackshear/frigate:stable
    container_name: frigate
    restart: unless-stopped
    privileged: true
    volumes:
      - ./frigate.yml:/config/config.yml
      - /media/frigate:/media/frigate
      - /dev/bus/usb:/dev/bus/usb  # For Coral TPU
    ports:
      - "5000:5000"  # Web UI
      - "8554:8554"  # RTSP
    environment:
      - FRIGATE_RTSP_PASSWORD=${RTSP_PASSWORD}

  # MQTT for Frigate events
  mosquitto:
    image: eclipse-mosquitto
    container_name: mosquitto
    restart: unless-stopped
    ports:
      - "1883:1883"
    volumes:
      - ./mosquitto.conf:/mosquitto/config/mosquitto.conf

  # NATS Leaf Node (bridges to cloud)
  nats-leaf:
    image: nats:alpine
    container_name: nats-leaf
    restart: unless-stopped
    command: -c /etc/nats/nats-leaf.conf
    volumes:
      - ./nats-leaf.conf:/etc/nats/nats-leaf.conf
      - ./edge.creds:/etc/nats/edge.creds
    ports:
      - "4222:4222"

  # Edge Agent
  edge-agent:
    build:
      context: ../../..
      dockerfile: infra/containers/Dockerfile.edge
    container_name: edge-agent
    restart: unless-stopped
    environment:
      - EDGE_LOCATION=${EDGE_LOCATION:-home}
      - NATS_URL=nats://nats-leaf:4222
      - OLLAMA_URL=http://ollama:11434
      - FRIGATE_URL=http://frigate:5000
    depends_on:
      - ollama
      - nats-leaf

volumes:
  ollama_data:
  frigate_data:
```

## NATS Leaf Configuration

```conf
# nats-leaf.conf
port: 4222

leafnodes {
  remotes [
    {
      url: "nats://cloud.onemindos.dev:7422"
      credentials: "/etc/nats/edge.creds"
    }
  ]
}

# Local JetStream for offline queuing
jetstream {
  store_dir: "/data/jetstream"
  max_mem: 256MB
  max_file: 1GB
}
```

## Edge Agent Code

```python
# agents/edge/main.py
import os
import asyncio
from agno import Agent, RemoteAgent
from agno.models import Ollama
from agno.tools import Tool

# Configuration
LOCATION = os.getenv("EDGE_LOCATION", "home")
NATS_URL = os.getenv("NATS_URL", "nats://localhost:4222")
OLLAMA_URL = os.getenv("OLLAMA_URL", "http://localhost:11434")
FRIGATE_URL = os.getenv("FRIGATE_URL", "http://localhost:5000")

# Local model for offline operation
local_model = Ollama(
    model="phi3:mini",
    host=OLLAMA_URL
)

# Define edge-specific tools
class FrigateTool(Tool):
    """Query Frigate camera system."""

    async def get_events(self, hours: int = 1) -> list:
        """Get recent detection events."""
        async with aiohttp.ClientSession() as session:
            async with session.get(
                f"{FRIGATE_URL}/api/events",
                params={"after": time.time() - hours * 3600}
            ) as resp:
                return await resp.json()

    async def get_snapshot(self, camera: str) -> bytes:
        """Get current camera snapshot."""
        async with aiohttp.ClientSession() as session:
            async with session.get(
                f"{FRIGATE_URL}/api/{camera}/latest.jpg"
            ) as resp:
                return await resp.read()

class AudioTool(Tool):
    """Control audio I/O."""

    async def speak(self, text: str):
        """Text-to-speech through local speakers."""
        # Use piper or espeak for local TTS
        os.system(f'espeak "{text}"')

    async def listen(self, duration: int = 5) -> str:
        """Listen for speech and transcribe."""
        # Use whisper.cpp for local transcription
        pass

# Create edge agent
edge_agent = Agent(
    name=f"legacy-{LOCATION}",
    model=local_model,
    tools=[FrigateTool(), AudioTool()],
    instructions=f"""You are Legacy's edge presence at {LOCATION}.

    You can:
    - Check cameras for people/objects
    - Speak through local speakers
    - Listen for voice commands
    - Report sensor readings

    When offline, handle commands locally.
    When online, sync with cloud Legacy."""
)

# Connect to cloud via NATS
async def main():
    remote = RemoteAgent(
        agent=edge_agent,
        nats_url=NATS_URL,
        location=LOCATION
    )

    # Register capabilities
    await remote.register_capabilities([
        "camera",
        "audio",
        "presence",
        "offline_llm"
    ])

    await remote.connect()
    print(f"Edge agent '{LOCATION}' connected")

    # Handle local wake word
    while True:
        # Listen for "Hey Legacy"
        if await detect_wake_word():
            response = await edge_agent.run(
                await transcribe_audio()
            )
            await AudioTool().speak(response)

        await asyncio.sleep(0.1)

if __name__ == "__main__":
    asyncio.run(main())
```

## Frigate Configuration

```yaml
# frigate.yml
mqtt:
  host: mosquitto
  port: 1883

cameras:
  front_door:
    ffmpeg:
      inputs:
        - path: rtsp://user:pass@camera-ip:554/stream
          roles:
            - detect
            - record
    detect:
      width: 1280
      height: 720
      fps: 5

  living_room:
    ffmpeg:
      inputs:
        - path: rtsp://user:pass@camera2-ip:554/stream
          roles:
            - detect
    detect:
      width: 1280
      height: 720
      fps: 5

detectors:
  coral:
    type: edgetpu
    device: usb  # Coral USB Accelerator

objects:
  track:
    - person
    - dog
    - cat
    - car
```

## Deployment Steps

### 1. Flash Jetson

```bash
# Download JetPack 5.x from NVIDIA
# Flash to SD card using Etcher
```

### 2. Initial Setup

```bash
# SSH into Jetson
ssh jetson@192.168.1.x

# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Install NVIDIA Container Toolkit
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list
sudo apt-get update && sudo apt-get install -y nvidia-docker2
sudo systemctl restart docker
```

### 3. Clone and Configure

```bash
git clone https://github.com/yourusername/onemind-os.git
cd onemind-os/infra/compose/edge

# Configure environment
cp .env.example .env
nano .env
# Set EDGE_LOCATION=home (or office, farm, etc.)
```

### 4. Generate NATS Credentials

On cloud server:
```bash
# Generate credentials for this edge
nsc add user -a ONEMIND -n home-edge
nsc describe user -a ONEMIND -n home-edge -o /tmp/home-edge.creds
```

Copy to Jetson:
```bash
scp /tmp/home-edge.creds jetson@192.168.1.x:~/onemind-os/infra/compose/edge/edge.creds
```

### 5. Deploy

```bash
# Pull and start
docker compose pull
docker compose up -d

# Check logs
docker compose logs -f

# Test Ollama
docker exec -it ollama ollama run phi3:mini "Hello"

# Test Frigate
curl http://localhost:5000/api/stats
```

## Offline Operation

When cloud is unreachable:

1. **NATS Leaf** queues messages locally (up to 1GB)
2. **Ollama** handles LLM requests locally
3. **Frigate** continues recording and detecting
4. **Edge Agent** responds to wake words using local model

When connection restores:
1. NATS syncs queued messages to cloud
2. Edge agent reports offline events
3. Cloud Legacy gets full picture

## Resource Usage

| Service | RAM | CPU | GPU |
|---------|-----|-----|-----|
| Ollama (idle) | 200MB | 1% | 0% |
| Ollama (inference) | 2.5GB | 50% | 80% |
| Frigate | 500MB | 30% | 0% (uses Coral) |
| NATS Leaf | 50MB | 1% | 0% |
| Edge Agent | 100MB | 5% | 0% |
| **Total (idle)** | **~1GB** | **~40%** | **0%** |
| **Total (active)** | **~3.5GB** | **~90%** | **80%** |

## Locations

| Location | Hardware | Features |
|----------|----------|----------|
| Home | Jetson Nano + Coral | 6 cameras, audio, presence |
| Office | Jetson Nano | 2 cameras, audio |
| Farm | Jetson Orin Nano | Sensors, FarmBot, weather |
| Car | Raspberry Pi | GPS, OBD-II |

## Troubleshooting

### Ollama out of memory
```bash
# Use smaller quantization
ollama pull phi3:mini-q4_0
```

### Frigate high CPU
```bash
# Ensure Coral TPU is detected
lsusb | grep Google
# Should show: Google Inc. Coral USB Accelerator
```

### NATS not connecting
```bash
# Check leaf connection
docker logs nats-leaf | grep "leafnode"
# Should show: "Leafnode connection created"
```

### Edge agent offline
```bash
# Test local NATS
docker exec -it nats-leaf nats pub test "hello"
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
