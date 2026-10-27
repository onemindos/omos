"""Authenticated read-only acceptance from inside the actual Legacy container."""
import asyncio
import json
import os
from pathlib import Path

import httpx
from mcp import ClientSession
from mcp.client.streamable_http import streamable_http_client
import yaml

URL = "http://onemind-mcp.hermes.svc.cluster.local:8000/mcp"


async def main():
    token = os.environ["ONEMIND_MCP_TOKEN"]
    # Check only the owned non-secret keys; never export the volume-managed config.
    config = yaml.safe_load(Path("/opt/data/config.yaml").read_text())
    entry = config["mcp_servers"]["onemind"]
    assert entry["url"] == URL
    assert entry["headers"]["Authorization"] == "Bearer ${ONEMIND_MCP_TOKEN}"
    assert config["tools"]["tool_search"]["enabled"] == "on"
    report = {"config_verified": True, "progressive_disclosure_enabled": True, "tools": {}}
    async with httpx.AsyncClient(trust_env=False, timeout=30) as anonymous:
        report["anonymous_status"] = (await anonymous.post(URL, json={})).status_code
        assert report["anonymous_status"] == 401
    async with httpx.AsyncClient(headers={"Authorization": f"Bearer {token}"},
                                 trust_env=False, timeout=30) as client:
        async with streamable_http_client(URL, http_client=client) as (read, write, _):
            async with ClientSession(read, write) as session:
                await session.initialize()
                tools = {tool.name for tool in (await session.list_tools()).tools}
                expected = {"apps_list", "nats_server_stats", "nats_jetstream_stats", "prometheus_query"}
                assert expected <= tools
                for name, args in [("apps_list", {}), ("nats_server_stats", {}),
                                   ("nats_jetstream_stats", {}), ("prometheus_query", {"query": "up"})]:
                    result = await session.call_tool(name, args)
                    if result.isError:
                        raise RuntimeError(f"{name} failed; no response body exported")
                    report["tools"][name] = "passed"
                resources = await session.list_resources()
                assert "onemind://apps" in {str(resource.uri) for resource in resources.resources}
                await session.read_resource("onemind://apps")
                report["resource_read"] = "passed"
    identity = Path("/var/run/secrets/kubernetes.io/serviceaccount")
    async with httpx.AsyncClient(
        verify=str(identity / "ca.crt"), trust_env=False,
        headers={"Authorization": "Bearer " + (identity / "token").read_text()}, timeout=10,
    ) as kube:
        for name in ("onemind-mcp", "hermes-legacy"):
            response = await kube.get(f"https://kubernetes.default.svc/apis/apps/v1/namespaces/hermes/deployments/{name}")
            response.raise_for_status()
            data = response.json()
            assert data["status"].get("observedGeneration", 0) >= data["metadata"]["generation"]
            assert data["status"].get("availableReplicas", 0) == data["spec"]["replicas"] == 1
            report[name] = {"available_replicas": 1, "generation": data["metadata"]["generation"],
                            "images": [c["image"] for c in data["spec"]["template"]["spec"]["containers"]]}
        response = await kube.get("https://kubernetes.default.svc/apis/argoproj.io/v1alpha1/namespaces/argocd/applications/platform-legacy")
        response.raise_for_status()
        status = response.json()["status"]
        report["argo"] = {"sync": status["sync"]["status"], "health": status["health"]["status"],
                          "revisions": status["sync"].get("revisions", [])}
        assert report["argo"]["sync"] == "Synced"
        assert report["argo"]["health"] == "Healthy"
    print(json.dumps(report, indent=2))


if __name__ == "__main__":
    asyncio.run(main())
