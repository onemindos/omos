---
title: "DroneKit-Python"
type: reference
status: planned
tags: [platform, hardware, drones, sdk, python, autonomy]
---
# DroneKit-Python

> Python SDK for autonomous drone control. Script missions, custom behaviors, AI integration.

## Install
```bash
pip install dronekit dronekit-sitl
```

## Connect
```python
from dronekit import connect, VehicleMode

# USB
vehicle = connect('/dev/ttyACM0', baud=115200, wait_ready=True)

# Network (MAVProxy forwarding)
vehicle = connect('udp:127.0.0.1:14550', wait_ready=True)

# SITL (simulation)
vehicle = connect('tcp:127.0.0.1:5760', wait_ready=True)
```

## Common Operations
```python
# Arm and takeoff
vehicle.mode = VehicleMode("GUIDED")
vehicle.armed = True
vehicle.simple_takeoff(10)  # meters

# Goto
from dronekit import LocationGlobalRelative
point = LocationGlobalRelative(-34.364114, 149.166022, 20)
vehicle.simple_goto(point, groundspeed=5)

# Return to launch
vehicle.mode = VehicleMode("RTL")

# Read telemetry
print(f"GPS: {vehicle.location.global_frame}")
print(f"Battery: {vehicle.battery}")
print(f"Mode: {vehicle.mode.name}")
print(f"Armed: {vehicle.armed}")
```

## NATS Integration
```python
import asyncio, nats
from dronekit import connect

vehicle = connect('/dev/ttyACM0', baud=115200)

async def telemetry_loop():
    nc = await nats.connect()
    while True:
        await nc.publish("telemetry.drone.recon.gps", 
            f"{vehicle.location.global_frame.lat},{vehicle.location.global_frame.lon}".encode())
        await nc.publish("telemetry.drone.recon.battery",
            f"{vehicle.battery.voltage},{vehicle.battery.level}".encode())
        await asyncio.sleep(1)
```

## Resources
- https://dronekit-python.readthedocs.io
- https://github.com/dronekit/dronekit-python
