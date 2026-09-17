# RADS — Virtual World Simulation Setup

> Goal: prove the entire RADS loop — sensor → MQTT → NATS → Hermes decides →
> MCP/A2A → Zenoh/ROS 2 → robot/drone acts → shows on the One Mind map —
> with ZERO physical hardware and NO local GPU, running on a cloud VM,
> watched live from the Mac. Focus: testing + learning a virtual homestead/farm.
> Created: 2026-05-31 · Project: One Mind · Companion to RADS-Technical-Architecture.md

---

## 0. The mental-model correction (read this first)

**NATS is NOT the middle-man for everything.** That is the key misunderstanding.

- NATS is the **coordination / brain / state layer** — "WHAT should happen" and "WHAT is true right now."
- A robot's **split-second internal control** (motors, balance, flight) is NOT NATS — it is **Zenoh/DDS/ROS 2**, the robot's own fast reflexes.
- A cheap **sensor** is NOT NATS — it is **MQTT**.
- They all connect to NATS through small **bridge** programs.

**How to structure your AGENTS (the important rule):**
Agents (Hermes) only ever speak **NATS + MCP/A2A**. They never speak MAVLink, ROS 2, or Zenoh directly.
A **bridge per device-type** translates NATS commands into the device's native language.

```
Agent thinks/decides  ->  NATS  ->  bridge translates  ->  device's own protocol
   (Hermes, MCP/A2A)      (spine)     (per device type)     (MQTT / ROS2 / Zenoh / MAVLink)
```

This keeps agents simple and protocol-agnostic. The bridges handle the "last mile." Swap a
device and you only change its bridge — the agent never changes.

---

## 1. The full virtual stack (what runs where)

```
CLOUD CPU VM (Linux)                              YOUR MAC (just a browser)
────────────────────                              ────────────────────────
Virtual world:                                    Foxglove Studio (web)  ← watch live
  Gazebo / PX4 SITL  (virtual drone)              QGroundControl (map)   ← drone map
  FarmBot web app    (virtual garden + 3D)        Browser → One Mind UI  ← the cockpit
  Sensor simulators  (soil/water/weather scripts)
        │                                          (SSH tunnel / Tailscale to reach VM)
        ▼  native protocols
  Bridges:
    mqtt-bridge   (sensors  ↔ NATS)
    ros2-bridge   (robots   ↔ NATS via DDS/Zenoh)
    mavlink-bridge(drone    ↔ NATS via MAVSDK)
        │
        ▼
  NATS + JetStream      ← the spine (state, presence, commands)
        │
        ▼
  Hermes agent + Cortex ← the brain (decides, MCP/A2A actions)
        │
        ▼
  One Mind UI / map     ← served to the Mac
```

Everything except the Mac browser runs on ONE cloud CPU VM. No GPU needed (see §5).

---

## 2. Software needed for autonomous drones

| Piece | Role | Tool |
|---|---|---|
| Autopilot / flight controller | Drone reflexes, stabilize, fly waypoints | **PX4** (or ArduPilot) |
| Drone messaging protocol | How parts talk to the autopilot | **MAVLink** |
| Robot brain OS | Smart logic, navigation, sensors | **ROS 2** |
| ROS ↔ autopilot bridge | Connect brain to flight controller | **MAVSDK** (or MAVROS) |
| Ground control / map | "Go here", monitor | **QGroundControl** |
| Real-time machine fabric | Fast robot/drone comms under ROS 2 | **Zenoh** or **DDS** |
| Run it with NO hardware | Flight code in pure software | **PX4 SITL** (Software-In-The-Loop) |

Core trio to remember: **ROS 2 + PX4 + MAVLink.**

---

## 3. 100% virtual test coverage matrix (don't miss anything)

Every RADS component → how to test it virtually → the tool. Aim: nothing untested.

| RADS component | Virtual test | Tool (free, mostly CPU) |
|---|---|---|
| Drone flight + autonomy | Fly virtual drone, waypoints, failsafes | **PX4 SITL** + **Gazebo** (or **SIH**, runs on Mac) |
| Ground robot / rover | Drive virtual rover, obstacle avoidance | **Gazebo** or **Webots** |
| Robot control loops / RL | Test AI movement policies fast | **MuJoCo** or **gym-pybullet-drones** (PyBullet, CPU) |
| Farm/garden robot (homestead) | Plant/water/sequence a virtual garden + 3D view | **FarmBot web app** (self-host, MIT, ships MQTT + 3D) |
| Soil / water / weather sensors | Fake sensor scripts publishing readings | Python script → **Mosquitto** (MQTT broker) |
| Home / energy / IoT devices | Simulate home actuation + state | **Home Assistant** (already running as ha-bridge) |
| MQTT layer | Sensor telemetry in/out | **Mosquitto** or EMQX |
| NATS spine | Pub/sub, streams, KV, presence | **NATS + JetStream** (already running) |
| Zenoh/DDS robot fabric | Real-time robot messages + ROS 2 | **ROS 2** + **Zenoh** bridge |
| AI agent decisions | Hermes reacts to events, takes actions | **Hermes** + **MCP/A2A** (already running) |
| Agent ↔ agent delegation | One agent hands task to another | **A2A** over NATS |
| Live visualization | Watch all data/robots in real time | **Foxglove Studio** (web, on Mac) |
| The cockpit / map | Unified human view | **One Mind UI** |
| Multi-community federation | Two RADS leaf nodes sync up | **NATS leaf nodes** (run 2 VMs/containers) |

If every row above runs green in sim, the whole architecture is proven before buying anything.

---

## 4. Source opportunities (real projects to leverage)

| Project | What it gives you | Note |
|---|---|---|
| **PX4 SITL** (docs.px4.io) | Real drone autopilot, no hardware | **SIH** mode runs natively on Apple-Silicon Mac; Gazebo best on Linux VM |
| **Gazebo** (gazebosim.org) | 3D physics world, sensors, robots | CPU ok; GPU only for camera realism |
| **Webots** | Friendly robot simulator | Lightweight alternative to Gazebo |
| **MuJoCo** (DeepMind, free) | Fast physics for control/RL | CPU, great for learning agent control |
| **gym-pybullet-drones** | Python drone sim for RL/agents | CPU-only, simple to start |
| **FarmBot Web App** (MIT, github.com/FarmBot/Farmbot-Web-App) | Virtual garden UI + 3D + **Dockerized MQTT** | Perfect homestead analog; speaks MQTT into your stack |
| **Home Assistant** | IoT/home/energy simulation + dashboards | Already running locally (ha-bridge) |
| **Mosquitto / EMQX** | MQTT broker for sensors | Tiny, runs anywhere |
| **Foxglove Studio** | Live web dashboard for robots/data | View the whole sim from the Mac browser |
| **QGroundControl** | Drone map + manual control | Connects to PX4 SITL over UDP 14550 |
| **NVIDIA Isaac Sim** | Photorealistic vision sim | LATER — needs rented cloud GPU |

---

## 5. The GPU truth

- **CPU is enough** for: flight dynamics, navigation, all protocols (MQTT/NATS/Zenoh/ROS 2),
  agent decisions, FarmBot, Home Assistant, Foxglove. This is ~95% of what you want to learn.
- **GPU only needed** for photorealistic **camera vision** (an AI "seeing" through a drone cam):
  Isaac Sim / AirSim. Defer it. When needed, **rent a cloud GPU by the hour and turn it off after.**
- Don't let "no GPU" block you — it blocks almost nothing at the testing/learning stage.

---

## 6. Mac + cloud topology

1. **One Linux CPU VM** (your existing cloud), ideally Docker Compose, running:
   Gazebo + PX4 SITL, Mosquitto, NATS + JetStream, Hermes, the bridges, FarmBot, One Mind UI.
2. **Reach it from the Mac** via **Tailscale** (easiest) or SSH tunnel. (PX4 docs cover SSH/mavlink-router for remote sim.)
3. **On the Mac:** browser only → Foxglove (live data), QGroundControl (drone map), One Mind UI (cockpit).

---

## 7. Phased plan (testing + learning first)

**Phase 1 — One sensor, one decision (the smallest real proof):**
Fake soil-moisture script → Mosquitto (MQTT) → mqtt-bridge → NATS → Hermes sees threshold →
publishes a command on NATS → logged → appears in One Mind UI. No drone yet. Proves the spine + brain.

**Phase 2 — Add the virtual drone:**
PX4 SITL + Gazebo. mavlink-bridge exposes the drone on NATS. Hermes commands "scan field 3";
watch it fly in Gazebo/QGroundControl, telemetry flowing to Foxglove.

**Phase 3 — Add the homestead:**
Self-host FarmBot web app (virtual garden + MQTT). Hermes runs a watering sequence. Add Home
Assistant for energy/home. Now you have a living virtual farm.

**Phase 4 — Make it a RAD + agents talking:**
Two NATS leaf nodes (two "communities") federating. A2A delegation between a "land agent" and a
"drone agent." This proves the community → network → map model end to end.

---

## 8. The homestead scenario you can watch live

"Virtual Homestead, morning cycle":
1. Weather + soil sensor scripts publish over **MQTT** → NATS.
2. **Hermes** (land agent) sees Field 3 dry + sunny → decides "water it, then inspect from air."
3. Tier check (autonomy): watering = auto; drone flight = notify-and-go.
4. **A2A**: land agent delegates "inspect Field 3" to the **drone agent**.
5. Drone agent commands **PX4 SITL** via mavlink-bridge → drone flies in **Gazebo**.
6. **FarmBot** runs the watering sequence in its virtual garden.
7. All telemetry streams to **Foxglove**; status + heartbeats show on the **One Mind map**.
8. You watch the entire RADS loop happen — on your Mac, no hardware, no GPU.

That single run exercises every protocol and every layer. When it's green, you've learned the
system and proven it works — ready to swap virtual devices for real ones later.
