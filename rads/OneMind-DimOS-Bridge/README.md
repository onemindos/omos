# OneMind × DimOS NATS Bridge

Connects any DimOS robot to the OneMind NATS fabric in **one module**.

```
┌─────────────────────────────────────────────────────┐
│               OneMind Fabric (NATS)                  │
│                                                      │
│  fabric.robot.<id>.odom   ← robot pose/velocity JSON │
│  fabric.robot.<id>.cot    ← TAK CoT XML (→ ATAK map) │
│  fabric.robot.<id>.cmd    → velocity command JSON     │
│  fabric.robot.<id>.status ← heartbeat                │
├────────────────────────┬────────────────────────────┤
│  Legacy (cloud pod)    │  DimOS (robot / Mac / edge) │
│  RobotCommander        │  NATSBridgeModule           │
│  legacy_commander.py   │  module.py                  │
└────────────────────────┴────────────────────────────┘
```

## Files

| File | Purpose |
|------|---------|
| `module.py` | DimOS Module — subscribes to odometry, publishes to NATS, receives commands |
| `blueprints.py` | Ready-made DimOS blueprints: `omos-go2-nats`, `omos-go2-sim-nats` |
| `legacy_commander.py` | Legacy/Python async client to command robots over NATS |

## Quickstart (no hardware)

```bash
# 1. Install DimOS + nats-py
uv venv --python 3.12 && source .venv/bin/activate
uv pip install 'dimos[base,unitree,sim]' nats-py

# 2. Start NATS (if not already running on the fabric)
docker run -p 4222:4222 nats:latest

# 3. Run Go2 in MuJoCo simulation with NATS bridge
ROBOT_ID=go2-sim NATS_URL=nats://localhost:4222 \
  dimos --simulation run omos-go2-sim-nats
```

## Watch the robot on the fabric (from Legacy pod)

```bash
# Live telemetry
nats sub 'fabric.robot.go2-sim.odom'

# Live CoT (paste into ATAK or TAK Server)
nats sub 'fabric.robot.go2-sim.cot'

# Command the robot to move forward
nats pub fabric.robot.go2-sim.cmd '{"linear":{"x":0.5},"angular":{"z":0.0}}'

# Stop
nats pub fabric.robot.go2-sim.cmd '{"linear":{"x":0.0},"angular":{"z":0.0}}'
```

## From Legacy (Python)

```python
from omos.dimos.nats_bridge.legacy_commander import RobotCommander
import asyncio

async def main():
    # Discover what robots are online
    robots = await RobotCommander.discover_robots()
    print("Online robots:", robots)

    # Command the first one
    robot = RobotCommander("go2-sim")
    await robot.forward(speed=0.5, duration_s=3.0)   # move 3 seconds then stop

    state = await robot.get_state()
    print(f"Position: x={state['x']:.2f} y={state['y']:.2f}")

asyncio.run(main())
```

## Real hardware

```bash
ROBOT_ID=go2-alpha ROBOT_IP=192.168.1.50 NATS_URL=nats://localhost:4222 \
  dimos run omos-go2-nats
```

## TAK Integration

The bridge publishes CoT XML to `fabric.robot.<id>.cot`.
The OneMind TAK feed pod can pick this up and forward to TAK Server with one rule:

```yaml
# In your feed pod config
nats_subject: "fabric.robot.*.cot"
tak_multicast: true
mission: "Robots"
```

Every robot then appears as a moving icon on ATAK/WebTAK automatically.
No GPS = position shows at (0,0) until you wire in a real GPS source.