"""Run all RADS sensors."""
import asyncio
from rads.sensors.soil import SoilSensorArray

if __name__ == "__main__":
    asyncio.run(SoilSensorArray().run())
