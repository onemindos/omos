#!/usr/bin/env python3
"""Read-only MCP transport and component-query smoke check. Never prints tokens or log bodies."""
import argparse
from datetime import datetime, timedelta, timezone
import json
import os
from pathlib import Path
import subprocess
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen


class Observer:
    def __init__(self, url, token=None):
        self.url, self.token, self.session, self.counter = url, token, None, 0

    def call(self, method, params=None, notification=False):
        self.counter += 1
        message = {"jsonrpc": "2.0", "method": method}
        if not notification:
            message["id"] = self.counter
        if params is not None:
            message["params"] = params
        headers = {"Content-Type": "application/json", "Accept": "application/json, text/event-stream"}
        if self.token:
            headers["Authorization"] = "Bearer " + self.token
        if self.session:
            headers["Mcp-Session-Id"] = self.session
            headers["MCP-Protocol-Version"] = "2025-03-26"
        request = Request(self.url, data=json.dumps(message).encode(), headers=headers)
        with urlopen(request, timeout=30) as response:
            self.session = response.headers.get("Mcp-Session-Id", self.session)
            if notification:
                return {}
            if "text/event-stream" in response.headers.get("Content-Type", ""):
                result = None
                for line in response:
                    if line.startswith(b"data:"):
                        event = json.loads(line[5:].strip())
                        if event.get("id") == self.counter:
                            result = event
                            break
                if result is None:
                    raise RuntimeError("Observer stream ended without a matching result")
            else:
                result = json.load(response)
        if "error" in result:
            raise RuntimeError("Observer returned a JSON-RPC error")
        value = result["result"]
        if value.get("isError"):
            raise RuntimeError("Observer tool returned an error")
        return value


def verify(url, component=None, namespace="default", project="onetak", environment="development"):
    try:
        Observer(url).call("initialize", {
            "protocolVersion": "2025-03-26", "capabilities": {},
            "clientInfo": {"name": "onemind-observer-check", "version": "1.0"}})
    except HTTPError as error:
        if error.code not in (401, 403):
            raise RuntimeError("Unexpected unauthenticated HTTP status " + str(error.code)) from None
    else:
        raise RuntimeError("Observer accepted an unauthenticated MCP request")
    helper = Path(__file__).with_name("openchoreo-token.sh")
    token = subprocess.check_output([str(helper), "--raw"], text=True).strip()
    client = Observer(url, token)
    client.call("initialize", {"protocolVersion": "2025-03-26", "capabilities": {},
                               "clientInfo": {"name": "onemind-observer-check", "version": "1.0"}})
    client.call("notifications/initialized", notification=True)
    tools, cursor = set(), None
    while True:
        result = client.call("tools/list", {"cursor": cursor} if cursor else {})
        tools.update(tool["name"] for tool in result["tools"])
        cursor = result.get("nextCursor")
        if not cursor:
            break
    required = {"query_component_logs", "query_resource_metrics", "query_traces",
                "query_alerts", "query_incidents", "query_costs"}
    if not required.issubset(tools):
        raise RuntimeError("Observer is missing required observability tools")
    print("PASS: unauthenticated requests refused; authenticated MCP exposes " + str(len(tools)) + " tools")
    if component:
        now = datetime.now(timezone.utc)
        args = {"namespace": namespace, "project": project, "component": component, "environment": environment,
                "start_time": (now - timedelta(hours=1)).isoformat(), "end_time": now.isoformat()}
        for tool in ("query_component_logs", "query_resource_metrics"):
            client.call("tools/call", {"name": tool, "arguments": args})
            print("PASS: " + tool + " completed (inspect console data separately)")
    return tools


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--url", default=os.environ.get("OBSERVER_MCP_URL",
                                                       "https://onemind-observer.tail717ba0.ts.net/mcp"))
    parser.add_argument("--component")
    parser.add_argument("--namespace", default="default")
    parser.add_argument("--project", default="onetak")
    parser.add_argument("--environment", default="development")
    args = parser.parse_args()
    try:
        verify(args.url, args.component, args.namespace, args.project, args.environment)
    except (RuntimeError, HTTPError, URLError, subprocess.CalledProcessError, KeyError, ValueError) as error:
        # HTTP bodies can contain private data. Report only a safe category/status.
        detail = str(error) if isinstance(error, RuntimeError) else type(error).__name__
        raise SystemExit("Observer verification failed: " + detail) from None


if __name__ == "__main__":
    main()
