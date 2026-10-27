import json
import logging
import secrets

import httpx
import pytest
from pydantic import ValidationError
from starlette.testclient import TestClient

from onemind_mcp.config import Catalog, load_catalog
from onemind_mcp.server import create_app
from onemind_mcp.upstream import Upstream

TOKEN = secrets.token_urlsafe(32)
HEADERS = {"Authorization": f"Bearer {TOKEN}", "Accept": "application/json, text/event-stream"}


def rpc(client, method, params=None, headers=None):
    response = client.post("/mcp", json={"jsonrpc": "2.0", "id": 1, "method": method,
                                         "params": params or {}}, headers=headers or HEADERS)
    assert response.status_code == 200, response.text
    return response.json()


@pytest.fixture
def setup(monkeypatch):
    monkeypatch.setenv("NATS_MONITOR_URL", "http://nats.nats.svc.cluster.local:8222")
    monkeypatch.setenv("PROMETHEUS_URL", "http://prometheus.monitoring.svc.cluster.local:9090")
    return load_catalog("config/apps.yaml")


def build(catalog, handler):
    return create_app(catalog, TOKEN, ["testserver"], upstream_factory=lambda url, token:
                      Upstream(url, token, httpx.MockTransport(handler)))


def test_protocol_discovery_resources_call_and_audit(setup, caplog):
    caplog.set_level(logging.INFO, logger="onemind.audit")
    def handler(request):
        assert str(request.url) == "http://nats.nats.svc.cluster.local:8222/varz"
        return httpx.Response(200, json={"connections": 4, "uptime": "2d", "authorization": TOKEN})
    with TestClient(build(setup, handler)) as client:
        assert client.get("/healthz").status_code == 200
        assert client.get("/readyz").status_code == 200
        init = rpc(client, "initialize", {"protocolVersion": "2025-11-25", "capabilities": {},
                                         "clientInfo": {"name": "test", "version": "1"}})
        assert "serverInfo" in init["result"]
        tools = rpc(client, "tools/list")["result"]["tools"]
        assert {t["name"] for t in tools} == {"apps_list", "apps_describe", "nats_server_stats", "nats_jetstream_stats", "prometheus_query"}
        assert all(t["annotations"]["readOnlyHint"] for t in tools)
        result = rpc(client, "tools/call", {"name": "nats_server_stats", "arguments": {}})["result"]
        assert not result.get("isError")
        assert "connections" in json.dumps(result)
        assert TOKEN not in json.dumps(result)
        resources = rpc(client, "resources/list")["result"]["resources"]
        assert resources[0]["uri"] == "onemind://apps"
        catalog = rpc(client, "resources/read", {"uri": "onemind://apps"})
        assert "prometheus_query" in json.dumps(catalog)
        assert "svc.cluster.local" not in json.dumps(catalog)
    logs = [json.loads(record.message) for record in caplog.records if record.name == "onemind.audit"]
    assert logs[-1]["client"] == "hermes"
    assert logs[-1]["tool"] == "nats_server_stats"
    assert TOKEN not in caplog.text


def test_unauthenticated_and_rebinding_rejected(setup):
    with TestClient(build(setup, lambda r: httpx.Response(200, json={}))) as client:
        assert client.post("/mcp", json={}).status_code == 401
        assert client.get("/mcp", headers={"Authorization": "Bearer wrong"}).status_code == 401
        response = client.post("/mcp", headers={**HEADERS, "Host": "attacker.example"}, json={})
        assert response.status_code == 421
        assert client.post("/mcp", headers={**HEADERS, "Origin": "https://attacker.example"}, json={}).status_code == 403


@pytest.mark.parametrize("status", [302, 401, 500])
def test_upstream_errors_and_redirects_do_not_leak(setup, status):
    requests = []
    def handler(request):
        requests.append(request)
        return httpx.Response(status, headers={"Location": "http://attacker.example"}, text=TOKEN)
    with TestClient(build(setup, handler)) as client:
        result = rpc(client, "tools/call", {"name": "nats_server_stats", "arguments": {}})["result"]
        assert result["isError"]
        assert TOKEN not in json.dumps(result)
    assert len(requests) == 1


def test_promql_query_encoded_validation_limit_and_secret_redaction(setup):
    requests = []
    def handler(request):
        requests.append(request)
        assert request.url.path == "/api/v1/query"
        assert request.url.params["timeout"] == "5s"
        assert request.url.params["limit"] == "100"
        return httpx.Response(200, json={"status": "success", "data": {
            "resultType": "vector", "result": [{"metric": {"label": TOKEN}, "value": [0, "1"]}]}})
    with TestClient(build(setup, handler)) as client:
        invalid = rpc(client, "tools/call", {"name": "prometheus_query", "arguments": {"query": ""}})["result"]
        assert invalid["isError"]
        assert not requests
        result = rpc(client, "tools/call", {"name": "prometheus_query", "arguments": {"query": 'up{job="app"}'}})["result"]
        assert not result.get("isError")
        assert TOKEN not in json.dumps(result)
    assert requests[0].url.params["query"] == 'up{job="app"}'


@pytest.mark.parametrize("kind", ["oversized", "timeout", "invalid_json", "too_many_series"])
def test_bounded_failure(setup, kind):
    def handler(request):
        if kind == "timeout":
            raise httpx.ReadTimeout("sensitive error", request=request)
        if kind == "oversized":
            return httpx.Response(200, content=b"x" * 1_048_577)
        if kind == "invalid_json":
            return httpx.Response(200, text="not json")
        return httpx.Response(200, json={"status": "success", "data": {"resultType": "vector", "result": [{}] * 101}})
    with TestClient(build(setup, handler)) as client:
        result = rpc(client, "tools/call", {"name": "prometheus_query", "arguments": {"query": "up"}})["result"]
        assert result["isError"]
        assert "sensitive error" not in json.dumps(result)


def test_fail_closed_configuration(setup, monkeypatch):
    with pytest.raises(ValueError, match="32"):
        create_app(setup, "short", ["testserver"])
    with pytest.raises(ValueError, match="Host"):
        create_app(setup, TOKEN, ["*"])
    monkeypatch.delenv("PROMETHEUS_URL")
    with pytest.raises(ValueError, match="PROMETHEUS_URL"):
        create_app(setup, TOKEN, ["testserver"])
    with pytest.raises(ValidationError):
        Catalog(apps=[setup.apps[0], setup.apps[0]])


def test_add_app_instance_without_core_changes(setup, monkeypatch):
    monkeypatch.setenv("EDGE_NATS_URL", "http://edge.nats.svc.cluster.local:8222")
    setup.apps.append(setup.apps[0].model_copy(update={"id": "edge", "url_env": "EDGE_NATS_URL"}))
    with TestClient(build(setup, lambda r: httpx.Response(200, json={}))) as client:
        names = {t["name"] for t in rpc(client, "tools/list")["result"]["tools"]}
        assert "edge_server_stats" in names
        assert "edge_jetstream_stats" in names
