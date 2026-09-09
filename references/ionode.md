# IOnode — ESP32 Native NATS Firmware

**Repo:** ionode.io / github.com/ionode  
**Language:** C++ (ESP32 firmware)  
**License:** MIT  
**Role in OneMind:** Flash-and-forget edge sensor node — any ESP32 becomes a NATS publisher

Flash any ESP32 and it immediately speaks NATS. Browser-based flasher at ionode.io/flash.html — no tools, configure from your phone.

## What it does

- Publishes sensor readings directly to NATS subjects
- Supports GPIO, I2C, SPI peripherals
- Auto-reconnect, TLS, NKey auth
- OTA firmware updates
- Zero-code sensor integration

## Why it matters for OneMind

Cheapest path to homestead telemetry. ESP32 boards cost $3-8. Flash IOnode, point at omos-core NATS, done:

```
ESP32 (IOnode) → tel.{node}.env.{metric} → NATS FABRIC → ClickHouse time machine
```

## Hardware targets

ESP32, ESP32-S3, ESP32-C3, ESP32-C6 (any 4MB+ flash variant)

## Integration point

```bash
# IOnode publishes to configurable subject prefix
tel.homestead.{sensor_id}.temperature
tel.homestead.{sensor_id}.humidity
tel.homestead.{sensor_id}.motion
```

## See Also

- [[WireClaw]] — AI agent variant of ESP32 NATS node (natural language hardware control)
- [[iron-scada]] — for Modbus/OPC-UA industrial hardware
