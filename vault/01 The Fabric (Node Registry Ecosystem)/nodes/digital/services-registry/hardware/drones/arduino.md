---
title: "Arduino"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, robotics, drones, microcontroller, iot, edge, embedded, ros2]
---
# Arduino

> Open-source microcontroller platform. Low-cost, accessible embedded computing. Entry point for hardware prototyping — sensors, actuators, motor controllers. Integrates with ROS2 via micro-ROS.

## What It Is

Arduino = hardware (AVR/ARM MCU boards) + IDE + C++ framework + ecosystem of libraries. The standard starting point for embedded/IoT/robotics:
- Simple C++ (no OS, no scheduler beyond `setup()`/`loop()`)
- 14,000+ libraries (sensors, motors, displays, comms)
- Huge community, long-tail hardware support

## Key Boards

| Board | MCU | Notes |
|-------|-----|-------|
| **Arduino Uno R4** | RA4M1 (32-bit ARM) | General purpose, WiFi variant |
| **Arduino Nano 33 BLE** | nRF52840 | BLE + IMU, good for wearables/drones |
| **Arduino Nano 33 IoT** | SAMD21 + ESP32 | WiFi/BT, IoT sensors |
| **Arduino Mega 2560** | ATmega2560 | Many I/O pins, robotics use |
| **Arduino Due** | SAM3X8E (ARM Cortex-M3) | 32-bit, higher compute |
| **Portenta H7** | STM32H747 (dual-core) | Industrial, ML inference |

Note: ESP32 (Espressif) is NOT Arduino hardware but runs in Arduino IDE — extremely common in IoT builds.

## micro-ROS (Arduino → ROS2)

micro-ROS runs a cut-down ROS2 stack on MCUs — Arduino becomes a ROS2 node:
```cpp
#include <micro_ros_arduino.h>
#include <std_msgs/msg/int32.h>

rcl_publisher_t publisher;
std_msgs__msg__Int32 msg;

void setup() {
  set_microros_transports();  // serial/WiFi/UDP
  rcl_init(); // simplified
  // create node, executor, publisher
}

void loop() {
  msg.data = analogRead(A0);
  rcl_publish(&publisher, &msg, NULL);
  rclc_executor_spin_some(&executor, RCL_MS_TO_NS(1));
}
```

micro-ROS Agent runs on the companion computer (Pi/Jetson) as the DDS bridge:
```bash
# On companion computer
ros2 run micro_ros_agent micro_ros_agent serial --dev /dev/ttyUSB0
```

## Common Use Cases in OneMind Fleet

| Use Case | Board | Protocol |
|----------|-------|---------|
| Sensor node (soil, temp, humidity) | Nano 33 IoT | MQTT/NATS |
| Motor controller (greenhouse actuators) | Mega 2560 | Serial/CAN |
| Drone ESC signal converter | Nano | PWM |
| CNC/tractor aux controller | Mega | G-code/Serial |
| ROS2 peripheral node | Any + micro-ROS | micro-ROS/DDS |
| Barn/IoT sensor array | ESP32 (Arduino IDE) | MQTT → NATS |

## Arduino IDE vs PlatformIO

| Tool | Best For |
|------|---------|
| **Arduino IDE 2** | Beginners, quick sketches |
| **PlatformIO** | Production, VS Code, dependency management, CI |

PlatformIO recommended for OneMind fleet builds — proper project structure, library pinning, OTA.

## Arduino vs Raspberry Pi vs ESP32

| Platform | OS | Best For |
|----------|-----|---------|
| **Arduino** | Bare metal | Real-time sensing, motor control, deterministic I/O |
| **ESP32** | FreeRTOS | WiFi/BT IoT nodes, runs Arduino framework |
| **Raspberry Pi** | Linux | Vision, AI, ROS2, companion computer |

Rule: Arduino for hard real-time peripherals. Pi for processing. NATS ties them together.

## NATS Integration

Direct (ESP32/nano with WiFi):
```cpp
// ESP32 NATS publish via WiFi
#include <WiFi.h>
#include <NATS.hpp>

NATS nats(client, "nats.omos.local", 4222);
nats.publish("eden.soil.moisture.zone1", "42.5");
```

Or bridge via micro-ROS → ROS2 → Zenoh → NATS.

## Related

- [[ros2-drones]] — micro-ROS bridges Arduino into ROS2
- [[zenoh]] — transport layer (micro-ROS can use Zenoh)
- [[nats]] — NATS direct publish from ESP32/WiFi Arduinos
- [[iot-fleet]] — Arduino nodes registered as edge devices in The Fabric
