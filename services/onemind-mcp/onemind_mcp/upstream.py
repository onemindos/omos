"""Bounded HTTP transport shared by adapters. No agent-controlled destinations."""
import json

import httpx
from mcp.server.fastmcp.exceptions import ToolError


class Upstream:
    def __init__(self, base_url: str, token: str | None = None,
                 transport: httpx.AsyncBaseTransport | None = None):
        self.base_url = base_url
        self.token = token
        self.auth_header = "Authorization"
        self.verify: str | bool = True
        self.transport = transport

    async def get(self, path: str, params: dict | None = None):
        return await self.request("GET", path, params=params)

    async def post(self, path: str, body: dict):
        return await self.request("POST", path, body=body)

    async def request(self, method: str, path: str, *, params: dict | None = None, body: dict | None = None):
        # Paths are constants in reviewed adapters, never caller input.
        if not path.startswith("/") or path.startswith("//"):
            raise ToolError("Invalid adapter endpoint")
        headers = ({self.auth_header: (f"Bearer {self.token}" if self.auth_header == "Authorization" else self.token)}
                   if self.token else {})
        try:
            async with httpx.AsyncClient(
                timeout=httpx.Timeout(10, connect=3), follow_redirects=False,
                trust_env=False, transport=self.transport, headers=headers, verify=self.verify,
            ) as client:
                async with client.stream(method, self.base_url + path, params=params, json=body) as response:
                    if response.status_code not in {200, 201}:
                        raise ToolError(f"Upstream returned HTTP {response.status_code}")
                    body = bytearray()
                    async for chunk in response.aiter_bytes():
                        body.extend(chunk)
                        if len(body) > 1_048_576:
                            raise ToolError("Result exceeds 1 MiB; narrow the query")
                    data = json.loads(body)
                    if not isinstance(data, (dict, list)):
                        raise ToolError("Upstream returned an unexpected response")
                    return data
        except (httpx.HTTPError, ValueError):
            # Exception strings can include upstream URLs and credentials.
            raise ToolError("Upstream unavailable or returned invalid JSON") from None
