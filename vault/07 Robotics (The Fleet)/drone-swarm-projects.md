---
title: Drone Swarm & Multi-Vehicle Open Source Projects
tags: [robotics, drones, swarm, px4, ardupilot, ros2]
owner: zeus
status: reference
created: 2026-06-11
---

========================================
DRONE SWARM & MULTI-VEHICLE OPEN SOURCE PROJECTS
========================================

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. QGroundControl Multi-Vehicle Support
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

▸ QGroundControl (Official)
  URL: https://github.com/mavlink/qgroundcontrol
  Stars: 4634
  Docs: https://docs.qgroundcontrol.com/master/en/
  Multi-Vehicle Docs: https://docs.qgroundcontrol.com/master/en/FlyView/MultiVehicleView.html
  Capabilities:
  - Native multi-vehicle support (view and control multiple drones simultaneously)
  - Vehicle switching from UI
  - Multi-vehicle mission planning
  - Individual vehicle telemetry display
  - MAVLink-based communication
  Notes: QGC supports multiple vehicles out-of-the-box but limited swarm coordination features

▸ QGC User Guide
  URL: https://github.com/mavlink/qgc-user-guide
  Stars: 48
  Docs: https://docs.qgroundcontrol.com/master/en/
  Notes: Official documentation source

▸ QGC Custom Builds & Extensions
  CubePilot QGC Fork: https://github.com/CubePilot/qgroundcontrol-herelink
  Notes: Hardware-specific customizations

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
2. PX4 Multi-Vehicle Swarm Coordination
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

▸ PX4 Autopilot (Official)
  URL: https://github.com/PX4/PX4-Autopilot
  Stars: 11874
  Docs: https://docs.px4.io/main/en/
  Multi-Vehicle Docs: https://docs.px4.io/main/en/simulation/multi-vehicle-simulation.html
  Capabilities:
  - SITL multi-vehicle simulation support
  - MAVLink multi-vehicle messaging
  - Swarm coordination via external controllers
  - Gazebo multi-vehicle simulation
  Notes: Core autopilot; swarm logic typically implemented externally

▸ Prometheus (AMOV Lab)
  URL: https://github.com/amov-lab/Prometheus
  Stars: 3160
  Docs: https://github.com/amov-lab/Prometheus/wiki
  Capabilities:
  - Multi-UAV formation control
  - Swarm collision avoidance
  - PX4 integration
  - ROS-based architecture
  - Gazebo simulation
  Notes: Popular Chinese open-source UAV swarm platform

▸ XTDrone
  URL: https://github.com/robin-shaun/XTDrone
  Stars: 1656
  Docs: https://www.yuque.com/xtdrone/manual_en
  Capabilities:
  - Multi-vehicle PX4 simulation
  - ROS integration
  - Gazebo/jMAVSim support
  - Computer vision and AI demos
  - Swarm formation flying
  Notes: Educational platform for UAV research

▸ PX4 Multi-Drone Framework
  URL: https://github.com/fpratissoli/PX4-multi-drone
  Stars: 3
  Capabilities:
  - Python-based multi-drone control
  - MAVSDK integration
  - Single and swarm operations
  - Real-time state monitoring
  - Mission planning for choreography
  Notes: Lightweight Python framework

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
3. ArduPilot Swarm/Multi-Copter Management
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

▸ ArduPilot (Official)
  URL: https://github.com/ArduPilot/ardupilot
  Stars: 12000+
  Docs: https://ardupilot.org/copter/
  Multi-Vehicle Docs: https://ardupilot.org/dev/docs/sim-multi-vehicle.html
  Capabilities:
  - SITL multi-vehicle simulation
  - MAVLink protocol support
  - Swarm via MAVProxy
  - Follow-me mode for formations
  Notes: Requires external coordination layer for true swarm behavior

▸ MAVProxy
  URL: https://github.com/ArduPilot/MAVProxy
  Docs: https://ardupilot.org/mavproxy/
  Capabilities:
  - Command-line GCS
  - Multi-vehicle management
  - Scripting support for swarm control
  - Module system for extensions
  Notes: Lightweight tool for multi-vehicle ArduPilot control

▸ SkyMaster Drone Platform
  URL: https://github.com/deepNblue/skymaster-drone-platform
  Stars: <10
  Capabilities:
  - Enterprise fleet management
  - PX4/ArduPilot/DJI support
  - Real-time monitoring
  - Mission planning
  - Swarm control interface
  Notes: New project, actively developed

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
4. ROS2-Based Drone Swarm Frameworks
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

▸ Crazyswarm2
  URL: https://github.com/IMRCLab/crazyswarm2
  Stars: 150+
  Docs: https://imrclab.github.io/crazyswarm2/
  Capabilities:
  - ROS2-native swarm control
  - Crazyflie micro-drone swarms
  - Motion capture integration
  - Centralized and decentralized modes
  - Research-grade trajectory planning
  Notes: Leading micro-drone swarm platform, successor to Crazyswarm

▸ Aerostack2
  URL: https://github.com/aerostack2/aerostack2
  Stars: 200+
  Docs: https://aerostack2.github.io/
  Capabilities:
  - ROS2-based multi-drone framework
  - Behavior trees for mission control
  - PX4 and Crazyflie support
  - Gazebo simulation
  - Modular architecture
  - Multi-robot coordination
  Notes: Modern ROS2 framework for aerial robotics

▸ Aerial Autonomy Stack
  URL: https://github.com/JacopoPan/aerial-autonomy-stack
  Stars: 449
  Capabilities:
  - Perception-based swarm control
  - PX4/ArduPilot support
  - ROS2 integration
  - YOLO vision processing
  - LiDAR support
  - NVIDIA Jetson deployment
  Notes: End-to-end autonomy stack with perception

▸ PX4-ROS2 Bridge
  URL: https://github.com/PX4/px4_ros_com
  Docs: https://docs.px4.io/main/en/ros/ros2_comm.html
  Capabilities:
  - uXRCE-DDS middleware
  - ROS2-PX4 communication
  - Multi-vehicle ROS2 support
  Notes: Official PX4-ROS2 interface

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
5. MAVSDK Multi-Drone Python Examples
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

▸ MAVSDK-Python (Official)
  URL: https://github.com/mavlink/MAVSDK-Python
  Stars: 400+
  Docs: https://mavsdk.mavlink.io/main/en/python/
  API Docs: https://mavsdk-python-docs.s3.amazonaws.com/index.html
  Examples: https://github.com/mavlink/MAVSDK-Python/tree/main/examples
  Capabilities:
  - Async Python API for drones
  - Multi-vehicle support via multiple System objects
  - PX4/ArduPilot compatible
  - MAVLink-based
  - Mission, telemetry, action plugins
  Notes: Official Python bindings, examples include multi-vehicle patterns

▸ PX4 Multi-Drone Framework (mentioned above)
  URL: https://github.com/fpratissoli/PX4-multi-drone
  Capabilities:
  - MAVSDK-Python wrapper for swarms
  - Single drone and swarm APIs
  - State monitoring
  - Mission planning
  - Choreography support

▸ MAVSDK Multi-Vehicle Examples
  Example repos to search:
  - https://github.com/search?q=MAVSDK+multi+vehicle+python
  - https://github.com/search?q=MAVSDK-Python+swarm
  Pattern: Create multiple System() instances, one per vehicle
  
  Basic multi-drone pattern:
  ```python
  from mavsdk import System
  
  drone1 = System()
  await drone1.connect(system_address="udp://:14540")
  
  drone2 = System()
  await drone2.connect(system_address="udp://:14541")
  ```

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ADDITIONAL NOTABLE PROJECTS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

▸ DroneKit-Python (Legacy but still used)
  URL: https://github.com/dronekit/dronekit-python
  Docs: https://dronekit-python.readthedocs.io/
  Notes: MAVLink Python API, multi-vehicle capable, less active

▸ PyMavlink
  URL: https://github.com/ArduPilot/pymavlink
  Docs: https://mavlink.io/en/mavgen_python/
  Notes: Low-level MAVLink Python library, multi-vehicle support

▸ Swarm Robotics Research (Academic)
  - MIT ACL Swarm: https://github.com/mit-acl
  - ETH Zurich Drones: https://github.com/ethz-asl
  - USC RESL: https://github.com/uscresl

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
KEY DOCUMENTATION RESOURCES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

• PX4 Multi-Vehicle Sim: https://docs.px4.io/main/en/simulation/multi-vehicle-simulation.html
• ArduPilot Multi-Vehicle: https://ardupilot.org/dev/docs/sim-multi-vehicle.html
• QGC Multi-Vehicle: https://docs.qgroundcontrol.com/master/en/FlyView/MultiVehicleView.html
• MAVSDK Python: https://mavsdk.mavlink.io/main/en/python/
• ROS2 PX4: https://docs.px4.io/main/en/ros/ros2_comm.html
• MAVLink Protocol: https://mavlink.io/en/

========================================
