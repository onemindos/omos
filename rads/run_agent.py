"""Run the land agent."""
import asyncio
from rads.agents.land_agent import LandAgent

if __name__ == "__main__":
    asyncio.run(LandAgent().run())
