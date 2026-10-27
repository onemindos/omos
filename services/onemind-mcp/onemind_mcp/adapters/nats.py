"""NATS HTTP monitoring, not message publishing or cluster administration."""
from ..upstream import Upstream


def register(add_tool, upstream: Upstream) -> None:
    async def server_stats() -> dict:
        """Read NATS node uptime, connection counts, traffic and slow consumers. This is one node, not a cluster total."""
        data = await upstream.get("/varz")
        fields = ("server_name", "version", "uptime", "connections", "total_connections",
                  "in_msgs", "out_msgs", "in_bytes", "out_bytes", "slow_consumers", "mem", "cpu")
        return {key: data[key] for key in fields if key in data}

    async def jetstream_stats() -> dict:
        """Read this NATS node's JetStream capacity and usage summary without messages, account details or stream configuration."""
        data = await upstream.get("/jsz")
        fields = ("server_id", "now", "memory", "storage", "reserved_memory", "reserved_storage",
                  "accounts", "ha_assets", "streams", "consumers", "messages", "bytes")
        return {key: data[key] for key in fields if key in data}

    add_tool(server_stats)
    add_tool(jetstream_stats)
