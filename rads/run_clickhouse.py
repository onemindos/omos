"""Run the ClickHouse Time Machine writer."""
import asyncio
from rads.storage.clickhouse_writer import ClickHouseWriter

if __name__ == "__main__":
    asyncio.run(ClickHouseWriter().run())
