---
type: hardware
category: edge-compute
vendor: Raspberry Pi Foundation
model: Pi Zero 2 W
status: available
price: ~$15
created: 2026-05-24
tags:
  - raspberry-pi
  - iot
  - edge-device
  - sensors
---

# Raspberry Pi Zero 2 W — Edge IoT Node

**Status:** Available now — $15  
**OneMind Role:** Lightweight IoT sensor nodes, MQTT/NATS bridges, GPIO control

## Specs

| Feature | Value |
|---------|-------|
| **CPU** | 1 GHz quad-core ARM Cortex-A53 (64-bit) |
| **RAM** | 512 MB LPDDR2 |
| **Connectivity** | 2.4 GHz 802.11n Wi-Fi, Bluetooth 4.2 BLE |
| **GPIO** | 40-pin header (same as Pi 4) |
| **Storage** | microSD card (no eMMC) |
| **Power** | 5V micro-USB (low power: ~1–2W idle, 3–4W under load) |
| **Ports** | Mini HDMI, 2x micro-USB (1 power, 1 OTG) |
| **Camera** | CSI connector for Pi Camera Module |
| **Size** | 65mm × 30mm (tiny!) |

## OneMind Integration Use Cases

### 1. IoT Seed Lite (Sensors)
- **Temperature/humidity:** DHT22, BME280
- **Motion detection:** PIR sensor (HC-SR501)
- **Door/window sensors:** Reed switches, magnetic sensors
- **Light sensors:** TSL2561, BH1750
- **Air quality:** MQ-135, PMS5003

**Architecture:**
```
Pi Zero → Sensor data → MQTT/NATS → Mac M2 → Home Assistant
```

### 2. Camera Node (Not for AI)
- **Pi Camera Module v2/v3:** 1080p video streaming
- **Motion-activated recording**
- **Stream to Mac for YOLOv8 processing** (Pi Zero too weak to run inference)

**Architecture:**
```
Pi Zero + Camera → RTSP stream → Mac M2 (YOLOv8) → NATS alert → HA Bridge → Home Assistant
```

### 3. Smart Home GPIO Control
- **Zigbee/Z-Wave USB dongle:** Control smart home devices
- **IR blaster:** Control TV, AC, legacy devices
- **Relay control:** Switch 110V devices (lights, fans)

### 4. Drone Telemetry Node
- **GPS module:** NEO-6M, NEO-M8N
- **Altimeter/IMU:** BMP280, MPU6050
- **Stream data to NATS** for real-time tracking

### 5. Edge NATS Client
- **Publish sensor data to NATS**
- **Subscribe to commands from Legacy/NATS**
- **Lightweight message routing**
- **Run Python/Node.js scripts**

## Limitations ⚠️

| Limitation | Impact |
|------------|--------|
| **512 MB RAM** | Cannot run YOLOv8 (needs 2GB+ for even nano model) |
| **Slow CPU** | Python inference ~0.1 FPS (unusable for real-time) |
| **No Neural Engine** | No hardware acceleration for ML |
| **microSD** | Slow I/O, wear issues for heavy logging |

### Conclusion
- ✅ **Good for:** Sensors, GPIO, MQTT/NATS, lightweight scripts
- ❌ **Bad for:** Computer vision, AI inference, heavy compute
- ✅ **Best role:** Dumb sensor nodes that forward data to Mac M2 for processing

## vs. Other Pi Models

| Model | CPU | RAM | AI Capable? | Price | OneMind Role |
|-------|-----|-----|-------------|-------|--------------|
| **Pi Zero 2 W** | 1 GHz quad-core | 512 MB | ❌ | $15 | Sensors only |
| **Pi 4 (4GB)** | 1.8 GHz quad-core | 4 GB | ⚠️ (slow) | $55 | Edge scripts |
| **Pi 5 (8GB)** | 2.4 GHz quad-core | 8 GB | ⚠️ (slow) | $80 | Edge scripts |
| **Pi 5 + AI HAT 2** | 2.4 GHz + 13 TOPS | 8 GB | ✅ (30 FPS) | $150 | **Edge vision** |
| **Mac M2 Max** | Neural Engine + GPU | 32 GB | ✅ (100 FPS) | $3000 | **Central AI hub** |

### Recommendation
- **Pi Zero 2 W:** Deploy as sensor nodes (10+ across home/farm)
- **Pi 5 + AI HAT 2:** Deploy for edge vision (cameras, drones)
- **Mac M2 Max:** Central AI hub (aggregate all data, run heavy inference)

## Installation & Setup (When Ready)

### 1. Flash Raspberry Pi OS Lite (64-bit)
```bash
# Use Raspberry Pi Imager
# Select: Raspberry Pi OS Lite (64-bit)
# Configure Wi-Fi, SSH, hostname in Imager
```

### 2. Install NATS Client
```bash
# SSH into Pi Zero
ssh onemind@pi-zero-01.local

# Install NATS CLI
curl -sf https://binaries.nats.dev/nats-io/natscli/nats@latest | sh

# Test NATS connection to Mac
nats --server nats://omos-mbp.local:4222 pub test.hello "Hello from Pi Zero"
```

### 3. Install Sensor Libraries (Example: DHT22)
```bash
# Install Python + GPIO libraries
sudo apt update
sudo apt install python3-pip python3-gpiozero

# Install sensor library
pip3 install adafruit-circuitpython-dht
```

### 4. Create Sensor Script
```python
# sensor_node.py
import time
import board
import adafruit_dht
import json
import subprocess

dht = adafruit_dht.DHT22(board.D4)  # GPIO 4

while True:
    try:
        temp = dht.temperature
        humidity = dht.humidity
        data = {"temp": temp, "humidity": humidity, "node": "pi-zero-01"}
        
        # Publish to NATS
        subprocess.run([
            "nats", "pub", "sensors.temp.living_room",
            json.dumps(data)
        ])
        print(f"Published: {data}")
    except RuntimeError as e:
        print(f"Error: {e}")
    
    time.sleep(60)  # Every minute
```

### 5. Run as systemd Service
```bash
# Create service file
sudo nano /etc/systemd/system/sensor-node.service

# Paste:
[Unit]
Description=OneMind Sensor Node
After=network.target

[Service]
ExecStart=/usr/bin/python3 /home/onemind/sensor_node.py
Restart=always
User=onemind

[Install]
WantedBy=multi-user.target

# Enable & start
sudo systemctl enable sensor-node
sudo systemctl start sensor-node
```

## Next Steps (Phase 3 — Edge Hardware)
- [ ] Order 5x Pi Zero 2 W ($75 total)
- [ ] Order sensors (DHT22, PIR, reed switches)
- [ ] Test NATS pub/sub to Mac
- [ ] Deploy first sensor node (living room temp/humidity)
- [ ] Build Home Assistant integration
- [ ] Scale to 10+ nodes across property

## Related Files
- Pi AI HAT 2 (edge vision): `hardware/pi-ai-hat-2.md`
- YOLOv8 deployment guide: `ai-providers/yolov8-mac-deployment-guide.md`
- Home Assistant: `50-74 LE/50 Home Automation/`
