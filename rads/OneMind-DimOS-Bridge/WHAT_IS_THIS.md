# OneMind × DimOS — NATS Bridge

## What is this folder?

This is the **glue layer** between DimOS (the robot OS) and OneMind OS (your NATS fabric).

Without this, DimOS robots are isolated — they run their own agent loop but have
no idea your NATS bus, TAK server, or Legacy exist.

With this, every robot (quadruped, drone, arm, humanoid) becomes a **first-class
citizen on your fabric** — it shows up on your ATAK map, Legacy can command it
in natural language, and its telemetry flows into ClickHouse like any other CoT source.

---

## The Big Picture

```
┌──────────────────────────────────────────────────────────────┐
│                    OneMind Fabric (NATS)                      │
│                                                              │
│  fabric.robot.<id>.odom   ◄─ robot pose + velocity (JSON)   │
│  fabric.robot.<id>.cot    ◄─ TAK CoT XML  → ATAK map icon   │
│  fabric.robot.<id>.cmd    ─► velocity commands from Legacy   │
│  fabric.robot.<id>.status ◄─ heartbeat / online signal      │
│                                                              │
├──────────────────────┬───────────────────────────────────────┤
│  Legacy (cloud pod)  │  DimOS (on your Mac / robot hardware) │
│                      │                                       │
│  legacy_commander.py │  module.py  ← drop into any blueprint │
│  "move forward 2m"   │  blueprints.py  ← ready-to-run        │
└──────────────────────┴───────────────────────────────────────┘
```

---

## Files in this folder

| File | What it is |
|------|-----------|
| `module.py` | The DimOS Module. Drop this into any DimOS blueprint and the robot is on your fabric. Publishes odometry + CoT to NATS, receives velocity commands back. |
| `blueprints.py` | Two ready-made DimOS blueprints: one for real Go2 hardware, one for MuJoCo simulation. |
| `legacy_commander.py` | Python async client for Legacy (or any script) to discover robots, read their state, and send movement commands over NATS. |

---

## How to use it (step by step)

### Step 1 — Install DimOS on your Mac

```bash
uv venv --python 3.12
source .venv/bin/activate
uv pip install 'dimos[base,unitree,sim]' nats-py
```

### Step 2 — Copy this folder into your project

```
~/Desktop/OneMind-DimOS-Bridge/   ← you are here
```

### Step 3 — Run a simulated robot with the bridge (zero hardware)

```bash
ROBOT_ID=go2-sim \
NATS_URL=nats://localhost:4222 \
dimos --simulation run omos-go2-sim-nats
```

The robot runs in MuJoCo. The bridge publishes its state to your NATS bus.

### Step 4 — Watch it appear on NATS

```bash
nats sub 'fabric.robot.go2-sim.odom'   # live telemetry
nats sub 'fabric.robot.go2-sim.cot'    # TAK CoT XML events
```

### Step 5 — Command it from Legacy

```bash
# Move forward
nats pub fabric.robot.go2-sim.cmd '{"linear":{"x":0.5},"angular":{"z":0.0}}'

# Turn left
nats pub fabric.robot.go2-sim.cmd '{"linear":{"x":0.0},"angular":{"z":0.5}}'

# Stop
nats pub fabric.robot.go2-sim.cmd '{"linear":{"x":0.0},"angular":{"z":0.0}}'
```

### Step 6 — Real hardware (when you have it)

```bash
ROBOT_ID=go2-alpha \
ROBOT_IP=192.168.1.50 \
NATS_URL=nats://localhost:4222 \
dimos run omos-go2-nats
```

One env var change. Everything else stays the same.

---

## TAK Map integration

The bridge already publishes CoT XML to `fabric.robot.<id>.cot`.

To get robots on your ATAK/WebTAK map, add one rule to your TAK feed pod
(same pattern as the MTA train feeds):

```yaml
nats_subject: "fabric.robot.*.cot"
tak_inject_subject: "tak.cot.inject"
```

Every robot appears as a moving ground vehicle icon.
Position defaults to (0,0) until you wire in real GPS.

---

## What comes next

1. Wire GPS into the CoT (use ATAK GeoCam or a USB GPS puck → feed into odom)
2. Add a DimOS camera module → stream video to MediaMTX → HelmCam on ATAK
3. Teach Legacy natural language commands: "go2, patrol the perimeter"
4. Add more robot types (drone, arm) — same bridge module, different blueprint

This is the foundation. Every piece plugs into what you already have.