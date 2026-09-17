```skill
# SKILL.md — homestead-control

## What This Skill Does

Controls and monitors Zeus's 6-building homestead through native OpenClaw `nodes` tool commands and WireClaw NATS sensor data. This is your SOP for all homestead operations — building status, environmental monitoring, device control, security checks, and multi-agent routing.

**This skill does NOT replace existing tools.** It teaches you WHEN and HOW to use:
- `nodes` → status, describe, run, notify, camera_snap (native OpenClaw)
- `exec` with `host=node` → run commands on specific building nodes
- `nexus_query` → read WireClaw sensor data from NATS
- `homestead_status` → aggregated status across all 6 buildings (custom tool from this skill)
- `homestead_alert` → threshold-based alert publishing (custom tool from this skill)

## When To Use It

Use this skill when Zeus asks about:
- **Building status** → "How's the barn?", "Status report", "Are all buildings online?"
- **Environmental data** → "Temperature in the greenhouse?", "Air quality in the barn?"
- **Security** → "Is the guest house locked?", "Any motion in the office after hours?"
- **Agriculture** → "Soil moisture?", "FarmBot status?", "Greenhouse yield?"
- **Wellness** → "Is my bedroom sleep-ready?", "Dojo session stats?"
- **Maintenance** → "Which nodes need attention?", "Uptime across buildings?"
- **Control** → "Turn on barn fans", "Lock the guest house", "Set greenhouse irrigation"

## Building Registry

| Node Name | Building | Location | Sensors | Voice | Special |
| --- | --- | --- | --- | --- | --- |
| `main-house` | Main House | Central, 4BR/3BA | motion, temp, door x5 rooms | kitchen, bedroom | HA hub, Zigbee coordinator |
| `office` | Office | Adjacent to house | motion, temp, power, network | desk | Power monitoring, UPS |
| `dojo` | Dojo | Training building | motion, temp, humidity, sound | mat area | Session tracking |
| `greenhouse` | Greenhouse | Growing structure | temp, humidity, soil x6 zones | none | FarmBot, MimiClaw |
| `guest-house` | Guest House | Separate dwelling | motion, temp, door | bedroom | Smart lock, occupancy |
| `barn` | Barn | Livestock + workshop | temp x3 stalls, motion, air quality | workbench | Climate control, MimiClaw |

## Node Commands — How To Check Buildings

### Check if all buildings are online

```

nodes → status

```

This returns all 6 nodes with online/offline state, uptime, and last-seen.

### Get a specific building's capabilities

```

nodes → describe --node barn

```

Returns what the node can do: system.run, system.notify, and any device capabilities.

### Run a command on a building's node

```

nodes → run --node greenhouse -- sensors --json

````

Or for more complex commands:

```json
{
  "action": "run",
  "node": "barn",
  "command": ["sensors", "--json"],
  "cwd": "/home/pi",
  "commandTimeoutMs": 5000
}
````

**Common commands per building:**

| Building    | Commands to Run                                                                                |
| ----------- | ---------------------------------------------------------------------------------------------- |
| main-house  | `sensors --json`, `ha-cli state summary`, `ha-cli state lock.front_door`                       |
| office      | `sensors --json`, `uptime`, `df -h /`, `ha-cli state sensor.ups_battery`                       |
| dojo        | `sensors --json`, `uptime`                                                                     |
| greenhouse  | `sensors --json`, `ha-cli state farmbot`, `ha-cli state switch.irrigation_zone_*`              |
| guest-house | `sensors --json`, `ha-cli state lock.guest_front_door`, `ha-cli state binary_sensor.occupancy` |
| barn        | `sensors --json`, `ha-cli state climate.barn`, `ha-cli state switch.barn_fans`                 |

### Send a notification to a building

```json
{
  "action": "notify",
  "node": "main-house",
  "message": "Guest arriving in 30 minutes"
}
```

### Camera/screen capture (if device supports it)

```json
{
  "action": "camera_snap",
  "node": "barn"
}
```

## WireClaw NATS Subjects — Sensor Data

WireClaw ESP32 sensors publish to NATS at 30ms latency. Use `nexus_query` to read them.

### Subject Convention

```
wireclaw.{building}.{room}.{sensor_type}
```

### Read latest temperature from all buildings

```
nexus_query({ subject: "wireclaw.*.*.temp", last: 1 })
```

### Read greenhouse soil moisture

```
nexus_query({ subject: "wireclaw.greenhouse.zone-a.soil", last: 1 })
nexus_query({ subject: "wireclaw.greenhouse.zone-b.soil", last: 1 })
nexus_query({ subject: "wireclaw.greenhouse.zone-c.soil", last: 1 })
```

### Read barn stall temperatures

```
nexus_query({ subject: "wireclaw.barn.stall-*.temp", last: 1 })
```

### Check for recent motion events

```
nexus_query({ subject: "wireclaw.guest-house.*.motion", last: 5 })
```

### Full subject map

```
# Main House
wireclaw.main-house.kitchen.{motion|temp}
wireclaw.main-house.bedroom.{temp|humidity}
wireclaw.main-house.living.{motion|door}
wireclaw.main-house.front.door

# Office
wireclaw.office.desk.{motion|temp}
wireclaw.office.server.temp
wireclaw.office.main.power

# Dojo
wireclaw.dojo.main.{motion|temp|humidity|sound}

# Greenhouse (6 zones)
wireclaw.greenhouse.zone-{a|b|c}.{soil|temp|humidity}
wireclaw.greenhouse.ambient.{temp|humidity}

# Guest House
wireclaw.guest-house.bedroom.{motion|temp}
wireclaw.guest-house.{main|front}.door

# Barn
wireclaw.barn.stall-{1|2|3}.temp
wireclaw.barn.main.{motion|temp|air-quality}
wireclaw.barn.workshop.motion
```

## Status Report Format

When Zeus asks for a "homestead status report" or "how are the buildings", compile this format:

```
🏠 HOMESTEAD STATUS — {date} @ {time}

1. MAIN HOUSE          {✅|⚠️|❌} {Online|Offline} (uptime {Xd})
   Temp: {X}°F | Humidity: {X}% | Motion: {room} ({Xm ago})
   Doors: {status} | Energy: {X} kWh today

2. OFFICE               {✅|⚠️|❌} {Online|Offline} (uptime {Xd})
   Temp: {X}°F | Disk: {X}% | Power: {X}W
   {occupancy status}

3. DOJO                  {✅|⚠️|❌} {Online|Offline} (uptime {Xd})
   Temp: {X}°F | Humidity: {X}% | Sound: {X}dB
   Last session: {date/time}

4. GREENHOUSE            {✅|⚠️|❌} {Online|Offline} (uptime {Xd})
   Zone A: Soil {X}% {🟢|🟡|🔴} | Temp: {X}°F | Humidity: {X}%
   Zone B: Soil {X}% {🟢|🟡|🔴} | Temp: {X}°F
   Zone C: Soil {X}% {🟢|🟡|🔴}
   FarmBot: {status}

5. GUEST HOUSE           {✅|⚠️|❌} {Online|Offline} (uptime {Xd})
   Temp: {X}°F | Occupancy: {occupied|vacant} | Door: {locked|unlocked}
   Last checkout: {date}

6. BARN                  {✅|⚠️|❌} {Online|Offline} (uptime {Xd})
   Ambient: {X}°F | Stall 1: {X}°F | Stall 2: {X}°F | Stall 3: {X}°F
   Fans: {ON|OFF} | Air: {good|moderate|poor}

📊 Summary: {X}/6 online | {X} advisories | {active alerts}
```

## Alert Thresholds

Use `homestead_alert` to publish alerts, or raise them proactively during HEARTBEAT checks.

| Condition                                        | Severity   | Action                                   |
| ------------------------------------------------ | ---------- | ---------------------------------------- |
| Any node offline > 5 min                         | `high`     | Alert Zeus + check network               |
| barn stall temp > 90°F                           | `critical` | Turn on fans + alert Zeus                |
| barn stall temp > 80°F                           | `warning`  | Log trend, prepare advisory              |
| greenhouse soil < 25%                            | `high`     | Trigger FarmBot irrigation for that zone |
| greenhouse soil < 35%                            | `warning`  | Advisory — irrigation may be needed soon |
| greenhouse temp > 100°F                          | `critical` | Open vents + alert Zeus                  |
| guest-house door open + unoccupied + after 11 PM | `high`     | Security alert                           |
| main-house bedroom temp > 76°F at bedtime        | `warning`  | Suggest AC adjustment (→ HP agent)       |
| dojo sound > 70dB at 5-8 AM                      | `info`     | Log training session start               |
| office power > 500W for 2+ hours                 | `warning`  | Check for runaway process                |

## Multi-Agent Routing

Not every homestead question is yours to answer alone. Route to the right agent:

| Request Type                                          | Route To                                                          | Example                                                            |
| ----------------------------------------------------- | ----------------------------------------------------------------- | ------------------------------------------------------------------ |
| Infrastructure, status, control                       | **You (legacy)**                                                  | "Turn on barn fans", "Status report"                               |
| Maintenance planning, long-term improvements, estate  | **le** (sessions_send)                                            | "Maintenance schedule", "Guest house renovation plan"              |
| Zeus's health impact from environment                 | **hp** (sessions_send)                                            | "Is my sleep environment optimal?", "Air quality affecting focus?" |
| Greenhouse economics, barn productivity, energy costs | **ge** (sessions_send)                                            | "Greenhouse yield vs cost?", "ROI on solar panels?"                |
| Complex analysis involving multiple domains           | **Orchestrate** (query yourself, then send summaries to le/hp/ge) | "Full homestead quarterly review"                                  |

### Delegation examples

Route to LE for maintenance:

```
sessions_send({
  sessionKey: "le",
  message: "Zeus wants a maintenance schedule for all 6 buildings. Current status: [attach node status data]. Please create a quarterly maintenance plan."
})
```

Route to HP for wellness:

```
sessions_send({
  sessionKey: "hp",
  message: "Bedroom temp has been 78°F for the past 3 nights. Office air quality trending moderate. Recommend environment adjustments for sleep and focus?"
})
```

Route to GE for economics:

```
sessions_send({
  sessionKey: "ge",
  message: "Greenhouse Zone A produced 40 lbs of lettuce this month. Energy cost for all 6 buildings: $X. Analyze yield per dollar and suggest optimizations."
})
```

## Exec Node Shorthand

For simple one-off commands, you can also use the `exec` tool with `host=node`:

```json
{
  "tool": "exec",
  "host": "node",
  "node": "greenhouse",
  "command": "cat /var/log/sensors/last-24h.json"
}
```

This is equivalent to `nodes → run --node greenhouse -- cat /var/log/sensors/last-24h.json` but uses the exec tool pathway. Prefer `nodes → run` for homestead operations (more explicit).

## Home Assistant Commands via Nodes

Each Pi has `ha-cli` installed pointing to the central HA instance. Common HA operations:

```bash
# Toggle a device
ha-cli service call switch.toggle --entity_id switch.barn_fans

# Set climate
ha-cli service call climate.set_temperature --entity_id climate.barn --temperature 72

# Lock/unlock
ha-cli service call lock.lock --entity_id lock.guest_front_door
ha-cli service call lock.unlock --entity_id lock.guest_front_door

# Query state
ha-cli state get sensor.greenhouse_zone_a_soil_moisture
ha-cli state get binary_sensor.guest_house_occupancy
```

Run these on any building node:

```json
{
  "action": "run",
  "node": "barn",
  "command": [
    "ha-cli",
    "service",
    "call",
    "switch.toggle",
    "--entity_id",
    "switch.barn_fans"
  ]
}
```

## MimiClaw Smart Sensors

Two buildings have MimiClaw ESP32-S3 upgrades (local ReAct reasoning on the chip):

- **Greenhouse**: Pre-analyzes soil + temp + humidity trends. Publishes enriched events like "Zone A trending dry, recommend irrigation in ~2 hours" instead of raw numbers.
- **Barn**: Monitors livestock comfort thresholds. Publishes "Stall 2 heat stress likely by 3 PM based on trend" instead of just temperature.

MimiClaw data arrives on the same NATS subjects but with enriched payloads containing `analysis` and `recommendation` fields.

## HeyClaw Voice Pucks

Five buildings have voice pucks (not greenhouse — too noisy). All voice requests arrive as normal agent messages. Context clues for where Zeus is:

- If the voice channel identifies a specific node, you know which building Zeus is in.
- Adjust responses accordingly ("The barn fans are on" vs "Your barn fans are on — you're standing 50 feet away").

```

```
