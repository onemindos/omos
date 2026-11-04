"""Protocol and adapter smoke test; run in an authorized pod/network with env credentials."""
import asyncio
import os

from mcp import ClientSession
from mcp.client.streamable_http import streamablehttp_client


async def main():
    token = os.environ["ONEMIND_MCP_TOKEN"]
    url = os.environ.get("ONEMIND_MCP_URL", "http://onemind-mcp.hermes.svc.cluster.local:8000/mcp")
    async with streamablehttp_client(url, headers={"Authorization": f"Bearer {token}"}) as (read, write, _):
        async with ClientSession(read, write) as session:
            await session.initialize()
            tools = await session.list_tools()
            names = {tool.name for tool in tools.tools}
            assert {"apps_list", "nats_server_stats", "nats_jetstream_stats", "prometheus_query"} <= names
            for name, arguments in [("apps_list", {}), ("nats_server_stats", {}),
                                    ("nats_jetstream_stats", {}), ("prometheus_query", {"query": "up"})]:
                result = await session.call_tool(name, arguments)
                if result.isError:
                    raise RuntimeError(f"{name} failed (inspect upstream reachability and permissions)")
                print(f"PASS {name}")


if __name__ == "__main__":
    asyncio.run(main())
