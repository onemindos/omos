---
title: "PyBullet"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, robotics, drones, simulation, physics, ml, reinforcement-learning, python]
---
# PyBullet

> Python bindings for the Bullet physics engine. Fast, lightweight robotics/drone sim used primarily for RL training and rapid prototyping. Lower fidelity than MuJoCo but zero cost and easy to set up.

## What It Is

PyBullet wraps the Bullet real-time physics SDK (used in games/VFX) for robotics simulation:
- Rigid body dynamics, constraints, articulated bodies
- URDF/SDF model loading (ROS standard robot format)
- OpenGL visualization (optional, can run headless)
- Direct Python control — no ROS required

## Install

```bash
pip install pybullet
```

## Basic Usage

```python
import pybullet as p
import pybullet_data

# Connect to GUI or headless
client = p.connect(p.GUI)          # visual
# client = p.connect(p.DIRECT)    # headless (faster, for training)

p.setAdditionalSearchPath(pybullet_data.getDataPath())
p.setGravity(0, 0, -9.81)

# Load ground plane + robot
plane = p.loadURDF("plane.urdf")
robot = p.loadURDF("r2d2.urdf", basePosition=[0, 0, 1])

# Step simulation
for _ in range(1000):
    p.stepSimulation()
    time.sleep(1./240.)
```

## Loading Drones / Robots

```python
# Load a quadrotor URDF
drone = p.loadURDF("cf2x.urdf", basePosition=[0, 0, 1])

# Apply forces to rotors
p.applyExternalForce(
    objectUniqueId=drone,
    linkIndex=-1,
    forceObj=[0, 0, thrust],
    posObj=[0, 0, 0],
    flags=p.LINK_FRAME
)
```

## PyBullet Aerial Gym / Drone Envs

Popular RL environments built on PyBullet for drones:
- **gym-pybullet-drones** (UToronto) — CF2x Crazyflie RL env, multi-agent
- **Aerial Gym Simulator** (Isaac/Bullet hybrid)

```bash
git clone https://github.com/utiasDSL/gym-pybullet-drones.git
pip install -e .
python3 examples/fly.py
```

## PyBullet vs Gazebo vs MuJoCo

| Sim | Fidelity | Speed | ROS2 | RL Training | Setup |
|-----|----------|-------|------|-------------|-------|
| **PyBullet** | Medium | ✅ Fast | Manual | ✅ Good | ✅ Easy |
| **Gazebo** | High | Slow | ✅ Native | ❌ Slow | Medium |
| **MuJoCo** | Very High | ✅ Fast | Manual | ✅ Best | Medium |

## SITL Note

PyBullet doesn't natively integrate with PX4/ArduPilot SITL. It's used for:
- RL policy training (train in PyBullet, deploy to real hardware or Gazebo for validation)
- Custom drone dynamics research
- Rapid algorithm testing before Gazebo/hardware

For full SITL (autopilot in the loop), use Gazebo.

## OneMind OS Angle

PyBullet = the fast RL training gym. Train navigation policies, obstacle avoidance, multi-drone coordination on CPU. Export trained policies → validate in Gazebo with PX4 SITL → deploy.

## Related

- [[gazebo]] — higher-fidelity SITL sim
- [[mujoco]] — higher-fidelity RL sim
- [[ros2-drones]] — deploy trained policies here
- [[lerobot]] — RL training framework (uses its own sim envs)
