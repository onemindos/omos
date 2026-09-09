# Octave — Reference Platform

## What It Is
Enterprise software platform for "mission-critical intelligence" across fixed sites and infrastructure. 14,000+ customers, 180 countries. Heavy in defense, government, critical infrastructure. Covers design, build, operate, and protect use cases for physical sites.

## What It Does
- **Fixed Asset Management** — tracks every sensor, camera, access point, and device at a facility
- **Physical Security** — integrates access control systems, alarm panels, video surveillance
- **Situational Awareness** — displays real-time status of all facility systems on a map
- **Incident Management** — logs events, dispatches response, tracks resolution
- **Digital Twin** — 3D model of facility with live sensor state overlaid
- **Cybersecurity + Physical** — converged security operations (IT + OT + physical)

## Data Protocols
- IoT: MQTT, OPC-UA, Modbus (industrial sensor protocols)
- Video: RTSP, ONVIF (IP camera standard)
- Access control: REST API integration with major access control vendors
- No published public API — enterprise integrations via professional services

## Why We're Building Our Own
- Octave is facility-first, static infrastructure — OneMind adds the mobile/field picture
- Their strength is monitoring fixed sites; our strength is the moving, fused intelligence picture
- Building our own means the facility layer feeds our NATS bus directly — no data silos
- A MQTT bridge → NATS is a small service; we don't need Octave's entire stack

---

## OneMind Site Layer — `omos-site`

### What We Build
An IoT/facility sensor ingestion pipeline that brings fixed-site awareness into the NATS bus and TAK picture.

### Architecture
```
Fixed site sensors:
  IP cameras (RTSP/ONVIF)
  Access control panels (REST/Wiegand)
  Environmental sensors (temperature, motion, intrusion)
  MQTT-enabled IoT devices (door sensors, power monitors)
  Alarm systems

  → MQTT broker (Mosquitto or NATS JetStream as MQTT bridge)
      → omos-site service
          → normalizes events to NATS subject grammar
              evt.facility.{site}.{sensor}.{event}
          → position-aware assets → Tile38 (geofence + spatial queries)
          → time-series data → ClickHouse
          → alert events → NATS → ORP → CoT → TAK Server
```

### NATS Subjects
```
evt.facility.{site}.access.{door_id}     → badge in/out event
evt.facility.{site}.camera.{cam_id}      → camera status / motion detect
evt.facility.{site}.alarm.{zone_id}      → alarm triggered
evt.facility.{site}.power.{panel_id}     → power status
evt.facility.{site}.sensor.{id}          → generic sensor reading
```

### CoT Integration
Facility alerts that require field response become CoT events:
```
Alarm triggered at north gate
  → NATS evt.facility.site-01.alarm.north-gate
      → ORP: correlates with nearby personnel CoT
          → Legacy AI: "Alarm at north gate, unit ALPHA-2 is 200m away"
              → CoT cmd → TAK Server → ATAK device of ALPHA-2
```

### Tech Stack
- Mosquitto MQTT broker (or NATS JetStream MQTT bridge)
- Python 3.12 omos-site service
- Tile38 for fixed asset positions + perimeter geofences
- ClickHouse for sensor time-series history
- NATS for event stream

### CloudTAK Display
- Fixed assets as persistent CoT entities (buildings, cameras, gates)
- Alert events as CoT events with callsign and position
- Sensor status overlaid as a MapLibre GeoJSON layer
- Historical sensor data via ClickHouse → Grafana dashboard

### Hardware Support
- Any ONVIF-compatible IP camera
- MQTT-enabled sensors (off-the-shelf IoT)
- Home Assistant integration (already in your Mac leaf node)
- Modbus/OPC-UA industrial sensors via bridge
