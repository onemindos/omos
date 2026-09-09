# RADS — Land Agent
**Service:** `com.onemind.rads-land-agent`
**Source:** `~/omos-rads/land_agent.py`
**Status at archive:** Disabled — hardware not connected (Exit 78)

## What It Does
Autonomous irrigation decision agent. Subscribes to soil moisture readings over MQTT and triggers irrigation commands when moisture drops below a threshold.

**Flow:**
1. Connects to local MQTT broker (`127.0.0.1:1883`)
2. Subscribes to `rad/demo/land/soil/moisture`
3. If moisture < 30% → publishes `irrigate` command to `rad/demo/land/commands/irrigate`
4. Enters cooldown until moisture recovers above 40%

## Config (via env vars)
| Var | Default | Description |
|-----|---------|-------------|
| `MQTT_HOST` | `127.0.0.1` | MQTT broker host |
| `MQTT_PORT` | `1883` | MQTT broker port |
| `MOISTURE_THRESHOLD` | `30` | Moisture % below which irrigation triggers |

## MQTT Topics
| Topic | Direction | Description |
|-------|-----------|-------------|
| `rad/demo/land/soil/moisture` | Subscribe | Receives moisture readings from soil sensor |
| `rad/demo/land/commands/irrigate` | Publish | Sends irrigation commands |

## Migration Plan
- Move to sim/edge node (RPi or Jetson)
- Update MQTT → NATS natively (remove paho dependency)
- Hardware: physical soil moisture sensor (not yet connected)

## Source Code
```python
import json, os
import paho.mqtt.client as mqtt
BROKER = os.environ.get('MQTT_HOST', '127.0.0.1')
PORT = int(os.environ.get('MQTT_PORT', '1883'))
MOISTURE_TOPIC = 'rad/demo/land/soil/moisture'
IRRIGATE_TOPIC = 'rad/demo/land/commands/irrigate'
THRESHOLD = float(os.environ.get('MOISTURE_THRESHOLD', '30'))
cooling_down = False

def on_connect(client, userdata, flags, rc, properties=None):
    print(f'[agent] connected to {BROKER}:{PORT} rc={rc}', flush=True)
    client.subscribe(MOISTURE_TOPIC)

def on_message(client, userdata, msg):
    global cooling_down
    data = json.loads(msg.payload)
    moisture = data.get('moisture')
    print(f'[agent] field={data.get("field")} moisture={moisture}', flush=True)
    if moisture is not None and moisture < THRESHOLD and not cooling_down:
        cmd = json.dumps({'action': 'irrigate', 'field': data.get('field'), 'reason': f'moisture {moisture} < {THRESHOLD}'})
        client.publish(IRRIGATE_TOPIC, cmd)
        print(f'[agent] DECISION -> {IRRIGATE_TOPIC} {cmd}', flush=True)
        cooling_down = True
    elif moisture is not None and moisture >= THRESHOLD + 10:
        cooling_down = False

client = mqtt.Client(client_id='land-agent', protocol=mqtt.MQTTv311)
client.on_connect = on_connect
client.on_message = on_message
client.connect(BROKER, PORT, keepalive=30)
client.loop_forever()
```
