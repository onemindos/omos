"""Authenticated read-only acceptance from inside the actual Legacy container."""
import asyncio
import json
import os
import re
import traceback
from pathlib import Path

import httpx
from mcp import ClientSession
from mcp.client.streamable_http import streamable_http_client
import yaml

URL = "http://onemind-mcp.hermes.svc.cluster.local:8000/mcp"


async def observe_rollout():
    """Only named workload status; no configuration, Secret or application data."""
    identity = Path("/var/run/secrets/kubernetes.io/serviceaccount")
    report = {}
    async with httpx.AsyncClient(
        verify=str(identity / "ca.crt"), trust_env=False,
        headers={"Authorization": "Bearer " + (identity / "token").read_text()}, timeout=10,
    ) as kube:
        for name in ("onemind-mcp", "hermes-legacy"):
            response = await kube.get(f"https://kubernetes.default.svc/apis/apps/v1/namespaces/hermes/deployments/{name}")
            if response.status_code != 200:
                report[name] = {"http_status": response.status_code}
                continue
            data = response.json()
            report[name] = {"available_replicas": data["status"].get("availableReplicas", 0),
                            "generation": data["metadata"]["generation"],
                            "observed_generation": data["status"].get("observedGeneration", 0)}
        response = await kube.get("https://kubernetes.default.svc/apis/argoproj.io/v1alpha1/namespaces/argocd/applications/platform-onemind-mcp")
        if response.status_code == 200:
            data = response.json()
            status = data["status"]
            report["argo"] = {"sync": status.get("sync", {}).get("status"),
                              "health": status.get("health", {}).get("status"),
                              "phase": status.get("operationState", {}).get("phase"),
                              "revisions": status.get("sync", {}).get("revisions", []),
                              "source_count": len(data["spec"].get("sources", []))}
    return report


async def main():
    token = os.environ.get("ONEMIND_MCP_TOKEN")
    if not token or os.environ.get("ONEMIND_MCP_TOOLSET_REVISION") != "app-coverage-v2":
        print(json.dumps({"phase": "waiting_for_legacy_rollout", **await observe_rollout()}, indent=2))
        return 75
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
        async with streamable_http_client(URL, http_client=client) as streams:
            read, write = streams[:2]
            async with ClientSession(read, write) as session:
                await session.initialize()
                tools = {tool.name for tool in (await session.list_tools()).tools}
                expected = {"apps_list", "nats_server_stats", "nats_jetstream_stats", "prometheus_query", "apps_describe", "stac_collections", "stac_search",
                            "ollama_models", "ollama_running_models", "alertmanager_alerts",
                            "loki_query", "opencost_allocation"}
                assert expected == tools
                report["tool_count"] = len(tools)
                for name, args in [("apps_list", {}), ("nats_server_stats", {}),
                                   ("nats_jetstream_stats", {}), ("prometheus_query", {"query": "up"}),
                                   ("apps_describe", {"app_id": "nextcloud"}),
                                   ("stac_collections", {}), ("stac_search", {"bbox": [-75, 40, -74, 41], "limit": 1}),
                                   ("ollama_models", {}), ("ollama_running_models", {}),
                                   ("alertmanager_alerts", {}), ("loki_query", {"query": '{namespace="hermes"}', "limit": 1}),
                                   ("opencost_allocation", {"window": "1h"})]:
                    result = await session.call_tool(name, args)
                    payload = result.model_dump(by_alias=True)
                    error_flag = payload.get("isError", payload.get("is_error"))
                    if not isinstance(error_flag, bool):
                        raise RuntimeError("Unrecognized MCP result error flag")
                    if error_flag:
                        error_text = " ".join(item.get("text", "") for item in payload.get("content", []))
                        report["failure"] = {"tool": name,
                            "http_status_codes": re.findall(r"Upstream returned HTTP ([0-9]{3})", error_text),
                            "flags": [flag for flag in ("invalid JSON", "unexpected response", "exceeds 1 MiB", "unavailable", "timeout") if flag in error_text]}
                        if name == "loki_query":
                            async with httpx.AsyncClient(trust_env=False, timeout=10) as diagnostic:
                                response = await diagnostic.get("http://loki-gateway.monitoring.svc.cluster.local/loki/api/v1/query_range",
                                    params={"query": '{namespace="hermes"}', "limit": 1, "since": "1h", "direction": "backward"})
                                report["failure"]["direct_status"] = response.status_code
                                if response.status_code != 200:
                                    # Classify known errors without exporting response content or logs.
                                    text = response.text.lower()
                                    report["failure"]["direct_flags"] = [term for term in (
                                        "parse error", "invalid", "timeout", "max", "limit", "too many", "not ready", "no org id", "ring", "replica", "unhealthy") if term in text]
                        if name == "loki_query":
                            identity = Path("/var/run/secrets/kubernetes.io/serviceaccount")
                            async with httpx.AsyncClient(verify=str(identity / "ca.crt"), trust_env=False,
                                headers={"Authorization": "Bearer " + (identity / "token").read_text()}, timeout=10) as kube:
                                response = await kube.get("https://kubernetes.default.svc/api/v1/namespaces/monitoring/services/loki-gateway")
                                if response.status_code == 200:
                                    report["failure"]["service_ports"] = response.json()["spec"]["ports"]
                        print(json.dumps(report, indent=2))
                        raise RuntimeError(f"{name} failed; no response body exported")
                    if name == "apps_list":
                        catalog = payload.get("structuredContent", payload.get("structured_content"))
                        if not catalog:
                            catalog = json.loads(next(item["text"] for item in payload["content"] if item.get("type") == "text"))
                        entries = catalog["apps"]
                        report["catalog"] = {
                            "domains": len(entries),
                            "configured": sorted(item["id"] for item in entries if item["access"] == "configured"),
                            "blocked": sorted(item["id"] for item in entries if item["access"] == "blocked"),
                        }
                        assert report["catalog"]["domains"] == 32
                        assert len(report["catalog"]["configured"]) == 7
                        assert all(not item["tools"] for item in entries if item["access"] == "blocked")
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
            if (data["status"].get("observedGeneration", 0) < data["metadata"]["generation"]
                    or data["status"].get("availableReplicas", 0) != data["spec"]["replicas"]):
                print(json.dumps({"phase": "waiting_for_rollout", **await observe_rollout()}, indent=2))
                return 75
            assert data["spec"]["replicas"] == 1
            report[name] = {"available_replicas": 1, "generation": data["metadata"]["generation"],
                            "images": [c["image"] for c in data["spec"]["template"]["spec"]["containers"]]}
        response = await kube.get("https://kubernetes.default.svc/apis/argoproj.io/v1alpha1/namespaces/argocd/applications/platform-onemind-mcp")
        response.raise_for_status()
        status = response.json()["status"]
        report["argo"] = {"sync": status["sync"]["status"], "health": status["health"]["status"],
                          "revisions": status["sync"].get("revisions", [])}
        if report["argo"]["sync"] != "Synced" or report["argo"]["health"] != "Healthy":
            print(json.dumps({"phase": "waiting_for_argo", **report}, indent=2))
            return 75
    print(json.dumps(report, indent=2))


if __name__ == "__main__":
    try:
        result = asyncio.run(main())
    except Exception as error:
        traceback.print_exception(error)
        raise SystemExit(78)
    raise SystemExit(result)
