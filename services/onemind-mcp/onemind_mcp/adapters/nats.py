"""NATS HTTP monitoring, not message publishing or cluster administration."""
from typing import Annotated
from urllib.parse import urlsplit

from mcp.server.fastmcp.exceptions import ToolError
from pydantic import Field

from ..upstream import Upstream


def node_labels(upstreams: list[Upstream]) -> dict[str, Upstream]:
    """Short node names (first DNS label, e.g. nats-1); never full hostnames."""
    labels = [(urlsplit(upstream.base_url).hostname or "").split(".")[0] for upstream in upstreams]
    if len(set(labels)) != len(labels) or not all(labels):
        labels = [f"node-{index}" for index in range(len(upstreams))]
    return dict(zip(labels, upstreams))


def register(add_tool, upstreams: list[Upstream]) -> None:
    nodes = node_labels(upstreams)
    Node = Annotated[str | None, Field(max_length=63, description=(
        f"NATS node to read, one of: {', '.join(nodes)}. "
        "Omit to read the first reachable node, failing over in that order."))]

    async def read(path: str, node: str | None, fields: tuple[str, ...]) -> dict:
        unavailable = []
        if node is not None:
            if node not in nodes:
                raise ToolError(f"Unknown NATS node; configured nodes: {', '.join(nodes)}")
            data = await nodes[node].get(path)
        else:
            for node, upstream in nodes.items():
                try:
                    data = await upstream.get(path)
                    break
                except ToolError:
                    unavailable.append(node)
            else:
                raise ToolError(f"No NATS monitor node reachable (tried {', '.join(unavailable)})")
        result = {"node": node, **{key: data[key] for key in fields if key in data}}
        return {**result, "unavailable_nodes": unavailable} if unavailable else result

    async def server_stats(node: Node = None) -> dict:
        """Read one NATS node's uptime, connection counts, traffic and slow consumers; not a cluster total. Call once per node to compare nodes."""
        return await read("/varz", node, (
            "server_name", "version", "uptime", "connections", "total_connections",
            "in_msgs", "out_msgs", "in_bytes", "out_bytes", "slow_consumers", "mem", "cpu"))

    async def jetstream_stats(node: Node = None) -> dict:
        """Read one NATS node's JetStream capacity and usage summary without messages, account details or stream configuration."""
        return await read("/jsz", node, (
            "server_id", "now", "memory", "storage", "reserved_memory", "reserved_storage",
            "accounts", "ha_assets", "streams", "consumers", "messages", "bytes"))

    add_tool(server_stats)
    add_tool(jetstream_stats)
