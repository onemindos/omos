"""Prometheus read APIs; no lifecycle or configuration endpoints."""
from typing import Annotated

from mcp.server.fastmcp.exceptions import ToolError
from pydantic import Field

from ..upstream import Upstream


def register(add_tool, upstream: Upstream) -> None:
    async def query(query: Annotated[str, Field(min_length=1, max_length=4096)]) -> dict:
        """Run an instant PromQL metrics query, e.g. up or sum(rate(http_requests_total[5m])). Returns at most 100 series; use selective queries."""
        data = await upstream.get("/api/v1/query", {"query": query, "timeout": "5s", "limit": 100})
        if data.get("status") != "success":
            raise ToolError("Prometheus rejected the query")
        # Bound output even on older Prometheus versions that ignore limit.
        result = data.get("data", {})
        if not isinstance(result, dict):
            raise ToolError("Unexpected Prometheus result")
        series = result.get("result")
        if isinstance(series, list) and result.get("resultType") in {"vector", "matrix"} and len(series) > 100:
            raise ToolError("More than 100 series; narrow the query")
        return {key: data[key] for key in ("status", "data", "warnings", "infos") if key in data}

    add_tool(query)
