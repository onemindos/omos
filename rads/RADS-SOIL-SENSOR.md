# RADS — Soil Sensor (Simulator)
**Service:** `com.onemind.rads-soil-sensor`
**Source:** `~/omos-rads/soil_sensor.py`
**Status at archive:** Disabled — hardware not connected (Exit 78)

## What It Does
Simulated soil moisture sensor. Publishes synthetic moisture readings over MQTT every 3 seconds, simulating a real field sensor draining and responding to irrigation events.

**Flow:**
1. Connects to local MQTT broker (`127.0.0.1:1883`)
2. Starts at 65% moisture, drains by 1.5–4.0% every 3 seconds (simulating evaporation/absorption)
3. Publishes readings to `rad/demo/land/soil/moisture`
4. Subscribes to `rad/demo/land/commands/irrigate` — when irrigation command received, moisture jumps +30% (simulating watering)

## This is a SIMULATOR — not real hardware
It was built to test the land_agent decision loop without physical sensors. The real hardware integration is pending.

## Config (via env vars)
| Var | Default | Description |
|-----|---------|-------------|
| `MQTT_HOST` | `127.0.0.1` | MQTT broker host |
| `MQTT_PORT` | `1883` | MQTT broker port |

## MQTT Topics
| Topic | Direction | Description |
|-------|-----------|-------------|
| `rad/demo/land/soil/moisture` | Publish | Sends simulated moisture readings every 3s |
| `rad/demo/land/commands/irrigate` | Subscribe | Receives irrigation trigger → boosts moisture |

## Payload Format
```json
{"field": 3, "moisture": 47.2, "ts": 1719432000.123}
```

## Migration Plan
- Port to real hardware (RPi + capacitive soil sensor)
- Replace MQTT with NATS native pub/sub
- Deploy to dedicated sim/edge node
- Field ID currently hardcoded as `3` — parameterize for multi-field support

## Source Code
```python
import json, os, random, time
import paho.mqtt.client as mqtt
BROKER = os.environ.get('MQTT_HOST', '127.0.0.1')
PORT = int(os.environ.get('MQTT_PORT', '1883'))
TOPIC = 'rad/demo/land/soil/moisture'
IRRIGATE_TOPIC = 'rad/demo/land/commands/irrigate'
moisture = 65.0

def on_connect(client, userdata, flags, rc, properties=None):
    print(f'[sensor] connected to {BROKER}:{PORT} rc={rc}', flush=True)
    client.subscribe(IRRIGATE_TOPIC)

def on_message(client, userdata, msg):
    global moisture
    print('[sensor] irrigation received -> wetting soil', flush=True)
    moisture = min(95.0, moisture + 30.0)

client = mqtt.Client(client_id='soil-sensor', protocol=mqtt.MQTTv311)
client.on_connect = on_connect
client.on_message = on_message
client.connect(BROKER, PORT, keepalive=30)
client.loop_start()
try:
    while True:
        moisture = max(5.0, moisture - random.uniform(1.5, 4.0))
        payload = json.dumps({'field': 3, 'moisture': round(moisture, 1), 'ts': time.time()})
        client.publish(TOPIC, payload)
        print(f'[sensor] -> {TOPIC} {payload}', flush=True)
        time.sleep(3)
except KeyboardInterrupt:
    client.loop_stop()
```
