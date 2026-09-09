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
