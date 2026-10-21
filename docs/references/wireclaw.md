# WireClaw — AI Agent on ESP32

**Repo:** wireclaw.io  
**Language:** C++ (ESP32 firmware)  
**Chips:** ESP32-C6, ESP32-S3, ESP32-C3 (4MB flash required)  
**Role in OneMind:** Natural language hardware control at the edge — $5 AI agent node

AI agent that lives on a microcontroller and controls real hardware. Tell it what you want in plain language — over Telegram, serial, or NATS — and it wires up GPIO pins, reads sensors, switches relays, and sets up automation rules that keep running without the AI. Flash from browser at wireclaw.io/flash.html.

## Key behaviors

- Natural language → hardware actions (GPIO, sensors, relays, UART bridge)
- Automation rules persist and run in the main loop — no cloud, no LLM after rule creation
- Memory: writes preferences/notes to flash, reloads on reboot
- NATS transport: receives commands and publishes state over NATS subjects
- Bridges serial devices (Arduinos, GPS, CO2 sensors, RFID) over UART
- Talks to other WireClaw devices on the network

## Why it matters for OneMind

Cheapest possible AI-actuated edge node. OneMind agents can send natural language commands to WireClaw nodes:

```
cmd.homestead.wireclaw.{id}.execute → "When temperature > 28, activate fan relay"
tel.homestead.wireclaw.{id}.sensor  ← sensor readings
evt.homestead.wireclaw.{id}.rule    ← rule created/triggered
```

The rule then runs locally, no further LLM involvement. Offline-capable.

## See Also

- [[ionode]] — IOnode is the sensor-only variant (no AI, just NATS publisher)
- [[iron-scada]] — for industrial-grade Modbus/PLC hardware
