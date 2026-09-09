"""Run the MAVLink→NATS bridge (Phase 2)."""
import asyncio
from rads.bridges.mavlink_bridge import run_bridge

if __name__ == "__main__":
    asyncio.run(run_bridge())
