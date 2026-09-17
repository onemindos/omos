---
title: "MuJoCo"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, robotics, drones, simulation, physics, ml, reinforcement-learning, deepmind]
---
# MuJoCo

> Multi-Joint dynamics with Contact. High-fidelity physics engine for robotics RL research. DeepMind acquired and open-sourced in 2022. Industry benchmark for legged robots, manipulation, and drone control.

## What It Is

MuJoCo (Mu-Jo-Co) is a physics engine purpose-built for model-based RL and optimal control:
- **Contacts** — soft contact model with friction, slip, rolling
- **Tendons** — biomechanical joint coupling
- **Actuators** — motors, servos, pneumatics, general force
- **Speed** — 10-100x faster than Gazebo for comparable scenes
- **Python API** — `dm_control`, `gymnasium`, native `mujoco` package

## Install

```bash
pip install mujoco
# Optional high-level wrappers
pip install dm_control
pip install gymnasium[mujoco]
```

## Basic Usage

```python
import mujoco
import numpy as np

# Load model from XML
model = mujoco.MjModel.from_xml_path("quadrotor.xml")
data = mujoco.MjData(model)

# Step simulation
for _ in range(1000):
    mujoco.mj_step(model, data)

# Read state
print(data.qpos)   # joint positions
print(data.qvel)   # joint velocities
print(data.sensordata)  # sensor readings
```

## Model Format

MuJoCo uses XML (`.xml`) model files:
- `worldbody` — kinematic tree
- `actuator` — motor definitions
- `sensor` — virtual sensors (IMU, force, position)
- `contact` — collision geometry

Pre-built models:
- `ant.xml`, `humanoid.xml` — locomotion benchmarks
- `quadrotor.xml` — drone in dm_control
- Unitree Go1/Go2 models available from Unitree + community

## Gymnasium Integration (RL Training)

```python
import gymnasium as gym

env = gym.make("HalfCheetah-v4")   # MuJoCo backend
obs, info = env.reset()

for _ in range(1000):
    action = env.action_space.sample()
    obs, reward, done, truncated, info = env.step(action)
    if done or truncated:
        obs, info = env.reset()
```

Standard benchmark envs: `Ant-v4`, `Humanoid-v4`, `HalfCheetah-v4`, `Hopper-v4`, `Walker2d-v4`

## dm_control (DeepMind)

Higher-level suite with tasks for manipulation and locomotion:
```python
from dm_control import suite
env = suite.load('cheetah', 'run')
timestep = env.reset()
```

## MuJoCo vs PyBullet vs Gazebo

| Sim | Physics Fidelity | Speed | RL Standard | ROS2 |
|-----|-----------------|-------|-------------|------|
| **MuJoCo** | ✅ Best | ✅ Fast | ✅ de facto | Manual |
| **PyBullet** | Medium | Fast | Good | Manual |
| **Gazebo** | High | Slow | Niche | ✅ Native |
| **Isaac Gym** | High | GPU-accelerated | Growing | Limited |

## Legged Robot Focus

MuJoCo is THE standard for legged robot RL (Unitree Go2, Boston Dynamics research, humanoids):
- DIMOS (OM1) training → MuJoCo → transfer to hardware
- Unitree Go2 SDK includes MuJoCo simulation config
- sim-to-real transfer: domain randomization in MuJoCo → deploy

## SITL Note

MuJoCo does not integrate with PX4/ArduPilot SITL directly. Use it for:
- Policy training (PPO, SAC, TD3) for drone/robot control
- Optimal control research
- Transfer learning (train in MuJoCo, fine-tune in Gazebo)

## OneMind OS Angle

MuJoCo is the training ground for DIMOS and legged fleet:
```
RL Policy training (MuJoCo) → Validation (Gazebo) → Hardware (DIMOS / OM1)
```
GPU-parallelized training via Isaac Gym integration possible for large-scale fleet RL.

## Related

- [[gazebo]] — SITL sim with autopilot integration
- [[pybullet]] — lighter alternative
- [[ros2-drones]] — policy deployment target
- [[lerobot]] — RL training framework (uses Gym envs including MuJoCo)
- [[unitree]] — primary hardware target for MuJoCo-trained policies
