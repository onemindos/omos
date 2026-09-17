# Edge Compute Requirements

> OneMind OS Edge Deployment Hardware Guide

## Overview

Each edge location runs 7 containerized services for local IoT, AI processing, and robotics. This document outlines the compute requirements and hardware recommendations.

---

## Edge Services (7 containers)

| Service | RAM | CPU | GPU/TPU | Purpose |
|---------|-----|-----|---------|---------|
| **nats** (leaf) | 30MB | Low | - | Local NATS + MQTT, bridges to cloud hub |
| **livekit** | 100-200MB | Medium | - | Voice/video for robotics (WebRTC) |
| **homeassistant** | 200-400MB | Low | - | Smart home automation |
| **frigate** | 500MB-2GB | **High** | Coral TPU recommended | Camera AI (NVR) |
| **ollama** | **2-4GB** | High | **GPU recommended** | Offline LLM |
| **remote-agent** | 50-100MB | Low | - | Legacy's local presence |
| **tailscale** | 30MB | Low | - | VPN mesh network |

### Service Details

#### NATS Leaf Node
- Bridges local MQTT devices to cloud NATS hub
- Provides local pub/sub when cloud is unavailable
- Port 4222 (NATS), 1883 (MQTT)

#### LiveKit (Local)
- Self-hosted WebRTC server for robotics
- Robot voice control and feedback
- Video streaming from robots/drones
- No cloud dependency or per-minute costs
- Ports: 7880 (API), 7881 (RTC), 7882 (TURN), 50000-50100 (media)

#### Home Assistant
- Smart home hub with 2000+ integrations
- Zigbee/Z-Wave via USB dongles
- Port 8123

#### Frigate NVR
- Real-time camera AI (person/car/animal detection)
- RTSP stream processing
- **Coral TPU strongly recommended** for multiple cameras
- Ports: 5000 (UI), 8554 (RTSP), 8555 (WebRTC)

#### Ollama
- Offline LLM for when cloud is unavailable
- Recommended models:
  - `phi3:mini` (2.7GB) - Best for Jetson/RPi
  - `llama3.2:3b` (2GB) - Lightweight alternative
  - `qwen2.5:3b` (2GB) - Good for coding tasks
- Port 11434

#### Remote Agent
- Legacy's local presence at each location
- Executes commands from cloud
- Reports sensor data back
- Port 8090

---

## Minimum vs Recommended Specs

| Spec | Minimum | Recommended | Heavy Workload |
|------|---------|-------------|----------------|
| **RAM** | 4GB | **8GB** | 16GB+ |
| **CPU** | 4-core ARM | 6-core ARM/x86 | 8-core |
| **Storage** | 32GB SD | **64GB SSD** | 128GB+ NVMe |
| **AI Accel** | None | Coral USB TPU | GPU (CUDA) |

### RAM Breakdown

| Config | Estimated Total |
|--------|-----------------|
| **Minimal** (no Ollama, 2 cameras) | ~2GB |
| **Standard** (phi3:mini, 4 cameras) | ~5GB |
| **Full** (llama3:8b, 8+ cameras) | ~10GB+ |

---

## Hardware Recommendations

### By Budget

| Device | Price | RAM | AI Capability | Best For |
|--------|-------|-----|---------------|----------|
| **Raspberry Pi 5** | $80 | 8GB | Coral USB ($60) | Home automation, light AI |
| **Orange Pi 5** | $100 | 16GB | Coral USB | Budget multi-camera |
| **Jetson Orin Nano** | $250 | 8GB | 1024 CUDA cores | **Robotics + AI** ⭐ |
| **Intel N100 Mini PC** | $150 | 16GB | Coral USB | High camera count |
| **Jetson AGX Orin** | $1,000+ | 32-64GB | 2048 CUDA cores | Heavy robotics/drones |

### By Use Case

#### 🏠 Home Automation Only
**Raspberry Pi 5 (4GB)** - $60
- Home Assistant + Frigate (2-3 cameras)
- No Ollama needed (use cloud LLM)
- Add Coral USB if more cameras

#### 🏠 Home + Basic AI
**Raspberry Pi 5 (8GB) + Coral USB** - $140
- Full stack with phi3:mini
- 4-6 cameras with Coral acceleration
- Good for most homes

#### 🤖 Robotics (Recommended)
**Jetson Orin Nano (8GB)** - $250
- GPU acceleration for Ollama
- Real-time robot vision processing
- LiveKit video streaming
- Can run llama3:8b

#### 📹 High Camera Count (10+)
**Intel N100 Mini PC + Coral** - $200
- x86 compatibility
- 16GB RAM for many camera streams
- Add M.2 Coral for best performance

#### 🚁 Heavy Robotics/Drones
**Jetson AGX Orin** - $1,000+
- Multiple robots/drones
- Real-time computer vision
- Large language models (llama3:70b)

---

## Jetson Orin Nano Setup (Recommended)

The Jetson Orin Nano is the sweet spot for robotics deployments.

### Specs
- **GPU:** 1024 CUDA cores (Ampere)
- **CPU:** 6-core Arm Cortex-A78AE
- **RAM:** 8GB LPDDR5
- **Storage:** NVMe slot (add your own SSD)
- **Power:** 7-15W

### Why Jetson for Robotics?
1. **GPU acceleration** for both Ollama and robot vision
2. **Low power** - can run on battery
3. **NVIDIA ecosystem** - Isaac ROS, DeepStream, TensorRT
4. **GPIO/I2C/SPI** - Direct hardware control
5. **Native ARM Docker** - All services run natively

### Initial Setup

```bash
# Flash JetPack 6.x via SDK Manager
# Boot and configure

# Install Docker (if not included)
sudo apt-get update
sudo apt-get install docker.io docker-compose-plugin

# Add user to docker group
sudo usermod -aG docker $USER

# Enable GPU in Docker
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker

# Clone OneMind OS
git clone https://github.com/your-org/onemind-os.git
cd onemind-os

# Configure edge environment
cp infra/compose/env/edge.env.example infra/compose/env/.env
# Edit .env with your settings

# Start edge services
cd infra/compose
docker compose -f docker-compose.edge.yml up -d
```

### Enable GPU for Ollama

Update `docker-compose.edge.yml`:

```yaml
ollama:
  image: ollama/ollama:latest
  container_name: edge-ollama
  restart: unless-stopped
  volumes:
    - ollama_data:/root/.ollama
  ports:
    - "11434:11434"
  networks:
    - onemind-edge
  deploy:
    resources:
      reservations:
        devices:
          - driver: nvidia
            count: all
            capabilities: [gpu]
```

---

## Raspberry Pi 5 Setup

For home automation without heavy AI.

### Specs
- **CPU:** 4-core Arm Cortex-A76 @ 2.4GHz
- **RAM:** 4GB or 8GB LPDDR4X
- **Storage:** microSD or NVMe HAT
- **Power:** 5V/5A USB-C

### Setup

```bash
# Flash Raspberry Pi OS (64-bit) via Imager
# Boot and configure

# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Clone and configure
git clone https://github.com/your-org/onemind-os.git
cd onemind-os
cp infra/compose/env/edge.env.example infra/compose/env/.env

# Start (without GPU flags)
cd infra/compose
docker compose -f docker-compose.edge.yml up -d
```

### Add Coral USB TPU

```bash
# Install Coral runtime
echo "deb https://packages.cloud.google.com/apt coral-edgetpu-stable main" | \
  sudo tee /etc/apt/sources.list.d/coral-edgetpu.list
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
sudo apt-get update
sudo apt-get install libedgetpu1-std

# Pass to Frigate container (uncomment in docker-compose.edge.yml)
# devices:
#   - /dev/bus/usb:/dev/bus/usb
```

---

## Storage Considerations

| Use Case | Recommended | Notes |
|----------|-------------|-------|
| Home Assistant only | 32GB SD | Minimal writes |
| + Frigate (clips) | 128GB SSD | Video storage |
| + Ollama models | 64GB+ SSD | Model files |
| Full stack | **256GB NVMe** | Best performance |

### Frigate Storage Math

| Resolution | FPS | Retain Days | Storage/Camera |
|------------|-----|-------------|----------------|
| 1080p | 5 | 7 | ~20GB |
| 1080p | 5 | 30 | ~85GB |
| 4K | 5 | 7 | ~80GB |

---

## Network Requirements

| Service | Bandwidth | Latency | Notes |
|---------|-----------|---------|-------|
| NATS leaf | Low | <100ms | Event messages |
| LiveKit | **High** | <50ms | Video streams |
| Frigate | Medium | N/A | Local only |
| Ollama | Low | N/A | Local only |
| Tailscale | Low | <100ms | VPN tunnel |

### Recommended Network
- **Wired Ethernet** for edge device
- **WiFi 6** for mobile robots
- **Tailscale** for secure cloud connection

---

## Power Considerations

| Device | Idle | Load | Notes |
|--------|------|------|-------|
| RPi 5 | 3W | 8W | USB-C 5V/5A |
| Jetson Orin Nano | 7W | 15W | DC barrel |
| Intel N100 | 10W | 25W | 12V adapter |
| Jetson AGX Orin | 15W | 60W | DC barrel |

### UPS Recommendations
- **CyberPower CP1500** - Desktop UPS, good for single edge
- **APC Back-UPS** - Budget option
- **Battery bank** - For mobile/outdoor deployments

---

## Quick Decision Matrix

| Your Needs | Recommended Hardware | Total Cost |
|------------|---------------------|------------|
| Home Assistant only | RPi 5 (4GB) | ~$70 |
| Home + 4 cameras | RPi 5 (8GB) + Coral | ~$150 |
| Home + basic AI | RPi 5 (8GB) + Coral | ~$150 |
| **Robotics** | **Jetson Orin Nano** | **~$300** |
| 10+ cameras | N100 Mini PC + Coral | ~$250 |
| Heavy AI/Robotics | Jetson AGX Orin | ~$1,200 |

---

## Related Documentation

- [Edge Compose File](../../infra/compose/docker-compose.edge.yml)
- [Edge Environment Variables](../../infra/compose/env/edge.env.example)
- [LiveKit Configuration](../../infra/configs/livekit/livekit.yaml)
- [NATS Hub-Leaf Architecture](./MULTI-LOCATION-NATS.md)
- [Remote Agent Architecture](./REMOTEAGENT-ARCHITECTURE.md)
