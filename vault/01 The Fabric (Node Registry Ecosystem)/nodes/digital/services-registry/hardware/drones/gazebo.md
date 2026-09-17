---
title: "Gazebo"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, drones, robotics, simulation, sitl, ros2, physics]
---
# Gazebo

> Industry-standard robotics simulator. Native ROS2 integration. Sensor plugins for IMU, GPS, LIDAR, cameras. Primary SITL backend for PX4 and ArduPilot.

## Versions

| Version | Name | Status |
|---------|------|--------|
| Gazebo Classic (11) | "Gazebo" | Legacy, EOL 2025 |
| **Gazebo Harmonic** | "gz sim" | Current, recommended |
| Gazebo Fortress/Garden | Transitional | Skip |

PX4 SITL targets Gazebo Harmonic (`gz_x500` models). ArduPilot supports both.

## What It Simulates

- **Physics** — rigid body dynamics, collisions, joint constraints
- **Sensors** — IMU, GPS, magnetometer, barometer, LIDAR, depth camera, RGB camera, optical flow
- **Environments** — outdoor maps, indoor spaces, custom worlds
- **Multi-vehicle** — multiple drones/robots in same world

## ROS2 Integration

```bash
# Bridge Gazebo topics ↔ ROS2
sudo apt install ros-humble-ros-gz-bridge

# Launch bridge
ros2 run ros_gz_bridge parameter_bridge \
  /model/x500/pose@geometry_msgs/msg/PoseStamped@gz.msgs.Pose

# Or use ros_gz_sim launch files (preferred)
```

## PX4 SITL with Gazebo

```bash
cd PX4-Autopilot
make px4_sitl gz_x500              # X500 quadcopter
make px4_sitl gz_x500_depth        # X500 + depth camera
make px4_sitl gz_standard_vtol     # VTOL

# Headless (no GUI, faster, CI/CD)
HEADLESS=1 make px4_sitl gz_x500
```

## ArduPilot SITL with Gazebo

```bash
# Install ArduPilot Gazebo plugin
git clone https://github.com/ArduPilot/ardupilot_gazebo
cd ardupilot_gazebo && mkdir build && cd build
cmake .. && make -j4 && sudo make install

# Launch world
gz sim -v4 -r iris_runway.sdf

# In separate terminal
sim_vehicle.py -v ArduCopter -f gazebo-iris --model JSON
```

## Key Sensor Plugins (Gazebo Harmonic)

| Plugin | Sensor |
|--------|--------|
| `gz-sim-imu-system` | IMU (accel + gyro) |
| `gz-sim-navsat-system` | GPS (NavSat) |
| `gz-sim-magnetometer-system` | Magnetometer |
| `gz-sim-sensors-system` | Generic sensor framework |
| `gz-sim-air-pressure-system` | Barometer |
| `gz-sim-camera-system` | RGB camera |
| `gpu_lidar` | LIDAR point cloud |

## SITL (Software In The Loop)

SITL runs the autopilot firmware **as a process on your computer** — no hardware. Gazebo provides the physics and sensor data the firmware would get from real hardware.

```
Gazebo (physics + sensors)
  ↕ UDP/socket (sim data)
PX4/ArduPilot SITL process
  ↕ UDP MAVLink (port 14550/14540)
MAVSDK / MAVROS / QGroundControl
```

## Worlds / Models

- PX4 models: `PX4-Autopilot/Tools/simulation/gz/models/`
- Custom worlds: `.sdf` format (XML)
- Models: Iris, X500, standard VTOL, rover, boat

## Gazebo vs PyBullet vs MuJoCo

| Sim | ROS2 | SITL | RL/ML | Physics Fidelity | Resource Use |
|-----|------|------|-------|-----------------|-------------|
| **Gazebo** | ✅ Native | ✅ PX4/APM | ⚠️ Slow | High | High |
| **PyBullet** | Manual | Manual | ✅ Fast | Medium | Low |
| **MuJoCo** | Manual | Manual | ✅ Best | Very High | Medium |

## OneMind OS Angle

Gazebo is the test environment before deploying to real hardware:
```
Develop mission → Test in Gazebo SITL → Validate behavior → Deploy to real vehicle
```
NATS-driven commands go through the same code path in sim and on hardware — sim validates the full stack.

## Related

- [[px4]] — primary SITL autopilot
- [[ardupilot]] — also supports Gazebo
- [[ros2-drones]] — ROS2 nodes tested against Gazebo
- [[pybullet]] — lighter-weight RL sim
- [[mujoco]] — high-fidelity physics sim
- [[foxglove-studio]] — visualize ROS2 data from Gazebo runs
