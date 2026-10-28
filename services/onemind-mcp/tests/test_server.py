import asyncio
from concurrent.futures import ThreadPoolExecutor
import json
import logging
import secrets

import httpx
import pytest
from pydantic import ValidationError
from starlette.testclient import TestClient

from onemind_mcp import server
from onemind_mcp.config import Catalog, load_catalog
from onemind_mcp.server import create_app, load_client_tokens
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


def build(catalog, handler, clients=TOKEN):
    return create_app(catalog, clients, ["testserver"], upstream_factory=lambda url, token:
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


OTHER = secrets.token_urlsafe(32)
CLIENTS = {"hermes": TOKEN, "claude-code": OTHER}


def bearer(token):
    return {**HEADERS, "Authorization": f"Bearer {token}"}


def audit(caplog):
    return [json.loads(record.message) for record in caplog.records if record.name == "onemind.audit"]


def run_entrypoint(monkeypatch, env):
    """Run main() with the deployment's environment names; capture the app instead of serving it."""
    for name in ("ONEMIND_MCP_TOKEN", "ONEMIND_MCP_TOKENS_FILE"):
        monkeypatch.delenv(name, raising=False)
    for name, value in {"ONEMIND_APPS_FILE": "config/apps.yaml", "ONEMIND_ALLOWED_HOSTS": "testserver", **env}.items():
        monkeypatch.setenv(name, value)
    apps = []
    monkeypatch.setattr(server.uvicorn, "run", lambda app, **options: apps.append(app))
    server.main()
    return apps[0]


def test_each_client_token_is_audited_as_its_client(setup, caplog):
    caplog.set_level(logging.INFO, logger="onemind.audit")
    with TestClient(build(setup, lambda r: httpx.Response(200, json={"server_name": OTHER}), CLIENTS)) as client:
        for token, tool in [(OTHER, "nats_server_stats"), (TOKEN, "apps_list"), (OTHER, "apps_list")]:
            result = rpc(client, "tools/call", {"name": tool, "arguments": {}}, bearer(token))["result"]
            assert not result.get("isError")
            assert OTHER not in json.dumps(result)
        for wrong in [TOKEN + OTHER, OTHER[:-1], secrets.token_urlsafe(32)]:
            assert client.post("/mcp", json={}, headers=bearer(wrong)).status_code == 401
    assert [(entry["client"], entry["tool"]) for entry in audit(caplog)] == [
        ("claude-code", "nats_server_stats"), ("hermes", "apps_list"), ("claude-code", "apps_list")]
    assert OTHER not in caplog.text and TOKEN not in caplog.text


def test_concurrent_calls_keep_their_own_client(setup, caplog):
    caplog.set_level(logging.INFO, logger="onemind.audit")
    arrived, both = [], asyncio.Event()
    async def handler(request):
        # Hold each upstream call until both are in flight, so the two audits interleave.
        arrived.append(request.url.path)
        if len(arrived) == 2:
            both.set()
        await asyncio.wait_for(both.wait(), 5)
        return httpx.Response(200, json={})
    calls = [(OTHER, "nats_server_stats"), (TOKEN, "nats_jetstream_stats")]
    with TestClient(build(setup, handler, CLIENTS)) as client, ThreadPoolExecutor(2) as pool:
        results = list(pool.map(lambda call: rpc(client, "tools/call", {"name": call[1], "arguments": {}},
                                                 bearer(call[0]))["result"], calls))
    assert not any(result.get("isError") for result in results)
    assert sorted((entry["client"], entry["tool"]) for entry in audit(caplog)) == [
        ("claude-code", "nats_server_stats"), ("hermes", "nats_jetstream_stats")]


@pytest.mark.parametrize("env_token, file_clients", [
    pytest.param(True, {}, id="deployed-single-token"),
    pytest.param(True, {"claude-code": OTHER}, id="env-hermes-plus-file"),
    pytest.param(False, CLIENTS, id="file-only"),
])
def test_entrypoint_client_sources(setup, monkeypatch, tmp_path, caplog, env_token, file_clients):
    caplog.set_level(logging.INFO, logger="onemind.audit")
    env = {"ONEMIND_MCP_TOKEN": TOKEN} if env_token else {}
    if file_clients:
        path = tmp_path / "tokens"
        path.write_text(json.dumps(file_clients))
        env["ONEMIND_MCP_TOKENS_FILE"] = str(path)
    expected = {"hermes": TOKEN, **file_clients}
    with TestClient(run_entrypoint(monkeypatch, env)) as client:
        for token in expected.values():
            rpc(client, "tools/call", {"name": "apps_list", "arguments": {}}, bearer(token))
        if OTHER not in expected.values():
            assert client.post("/mcp", json={}, headers=bearer(OTHER)).status_code == 401
    assert [entry["client"] for entry in audit(caplog)] == list(expected)


def test_tokens_file_accepts_json_or_name_token_lines(tmp_path):
    path = tmp_path / "tokens"
    path.write_text(json.dumps(CLIENTS))
    assert load_client_tokens(path) == CLIENTS
    path.write_text(f"# one client per line\nhermes:{TOKEN}\r\n\n  claude-code : {OTHER}  \n")
    assert load_client_tokens(path) == CLIENTS


@pytest.mark.parametrize("content", [
    pytest.param('{"hermes": "' + TOKEN + '",}', id="invalid-json"),
    pytest.param(f'{{"hermes": "{TOKEN}", "hermes": "{OTHER}"}}', id="duplicate-json-name"),
    pytest.param(f"hermes:{TOKEN}\nhermes:{OTHER}\n", id="duplicate-line-name"),
    pytest.param(f"{TOKEN}:hermes\n", id="swapped-name-and-token"),
    pytest.param(f"{TOKEN}\n", id="missing-separator"),
    pytest.param(f"Hermes:{TOKEN}\n", id="invalid-name"),
    pytest.param("hermes:short\n", id="short-token"),
    pytest.param(f"hermes:{TOKEN[:20]} {TOKEN[20:]}\n", id="whitespace-in-token"),
    pytest.param(json.dumps({"hermes": {"token": TOKEN}}), id="non-string-token"),
    pytest.param("# no clients\n", id="no-lines"),
    pytest.param("{}", id="empty-json"),
])
def test_tokens_file_fails_closed_without_quoting_secrets(tmp_path, content):
    path = tmp_path / "tokens"
    path.write_text(content)
    with pytest.raises(ValueError, match="ONEMIND_MCP_TOKENS_FILE") as error:
        load_client_tokens(path)
    assert not any(secret[:16] in str(error.value) for secret in (TOKEN, OTHER))


def test_client_configuration_fails_closed(setup, monkeypatch, tmp_path):
    for clients, message in [({}, "ONEMIND_MCP_TOKENS_FILE"), ({"Claude Code": OTHER}, "names must match"),
                             ({"hermes": TOKEN, "claude-code": "short"}, "32"),
                             ({"hermes": TOKEN, "claude-code": TOKEN}, "own token")]:
        with pytest.raises(ValueError, match=message):
            create_app(setup, clients, ["testserver"])
    path = tmp_path / "tokens"
    path.write_text(f"hermes:{OTHER}\n")
    with pytest.raises(ValueError, match="not both"):
        run_entrypoint(monkeypatch, {"ONEMIND_MCP_TOKEN": TOKEN, "ONEMIND_MCP_TOKENS_FILE": str(path)})
    with pytest.raises(FileNotFoundError):
        run_entrypoint(monkeypatch, {"ONEMIND_MCP_TOKEN": TOKEN, "ONEMIND_MCP_TOKENS_FILE": str(tmp_path / "missing")})
    with pytest.raises(ValueError, match="ONEMIND_MCP_TOKENS_FILE"):
        run_entrypoint(monkeypatch, {})


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


NODES = [f"http://nats-{i}.nats-headless.nats.svc.cluster.local:8222" for i in range(3)]


def nats_call(client, tool, arguments=None):
    result = rpc(client, "tools/call", {"name": tool, "arguments": arguments or {}})["result"]
    return result, (None if result.get("isError") else json.loads(result["content"][0]["text"]))


def test_single_url_fallback_names_node(setup):
    with TestClient(build(setup, lambda r: httpx.Response(200, json={"server_name": "nats-0"}))) as client:
        _, data = nats_call(client, "nats_server_stats")
    assert data == {"node": "nats", "server_name": "nats-0"}


def test_multinode_schema_is_backward_compatible(setup, monkeypatch):
    monkeypatch.setenv("NATS_MONITOR_URLS", ",".join(NODES))
    with TestClient(build(setup, lambda r: httpx.Response(200, json={}))) as client:
        tools = {t["name"]: t for t in rpc(client, "tools/list")["result"]["tools"]}
    for name in ("nats_server_stats", "nats_jetstream_stats"):
        schema = tools[name]["inputSchema"]
        assert "node" in schema["properties"]
        assert "node" not in schema.get("required", [])
        assert "nats-0, nats-1, nats-2" in json.dumps(schema)
        assert "svc.cluster.local" not in json.dumps(tools[name])


@pytest.mark.parametrize("failure", ["http_500", "connect_error"])
def test_multinode_fails_over_in_order(setup, monkeypatch, failure):
    monkeypatch.setenv("NATS_MONITOR_URLS", " , ".join(NODES))
    requests = []
    def handler(request):
        requests.append(request.url.host.split(".")[0])
        if request.url.host.startswith("nats-0."):
            if failure == "connect_error":
                raise httpx.ConnectError("nats-0 secret detail", request=request)
            return httpx.Response(500)
        return httpx.Response(200, json={"server_id": "nats-1", "streams": 39, "config": {"x": 1}})
    with TestClient(build(setup, handler)) as client:
        _, data = nats_call(client, "nats_jetstream_stats")
    assert data == {"node": "nats-1", "server_id": "nats-1", "streams": 39, "unavailable_nodes": ["nats-0"]}
    assert requests == ["nats-0", "nats-1"]


def test_multinode_explicit_node_does_not_fail_over(setup, monkeypatch):
    monkeypatch.setenv("NATS_MONITOR_URLS", ",".join(NODES))
    requests = []
    def handler(request):
        requests.append(request.url.host.split(".")[0])
        if request.url.host.startswith("nats-1."):
            return httpx.Response(503)
        return httpx.Response(200, json={"server_name": request.url.host.split(".")[0]})
    with TestClient(build(setup, handler)) as client:
        _, data = nats_call(client, "nats_server_stats", {"node": "nats-2"})
        assert data == {"node": "nats-2", "server_name": "nats-2"}
        result, _ = nats_call(client, "nats_server_stats", {"node": "nats-1"})
        assert result["isError"]
        unknown, _ = nats_call(client, "nats_server_stats", {"node": "nats-9"})
        assert unknown["isError"]
        assert "nats-0, nats-1, nats-2" in json.dumps(unknown)
    assert requests == ["nats-2", "nats-1"]


def test_multinode_all_down_is_bounded_error(setup, monkeypatch):
    monkeypatch.setenv("NATS_MONITOR_URLS", ",".join(NODES))
    def handler(request):
        raise httpx.ConnectError(f"refused {TOKEN}", request=request)
    with TestClient(build(setup, handler)) as client:
        result, _ = nats_call(client, "nats_server_stats")
    assert result["isError"]
    text = json.dumps(result)
    assert "nats-0, nats-1, nats-2" in text
    assert TOKEN not in text and "svc.cluster.local" not in text


def test_multinode_url_list_validation(setup, monkeypatch):
    monkeypatch.delenv("NATS_MONITOR_URL")
    monkeypatch.setenv("NATS_MONITOR_URLS", f"{NODES[0]},{NODES[0]}/")
    assert setup.apps[0].connections() == [(NODES[0], None)]
    monkeypatch.setenv("NATS_MONITOR_URLS", f"{NODES[0]},http://user:pw@nats-1:8222")
    with pytest.raises(ValueError, match="NATS_MONITOR_URLS"):
        create_app(setup, TOKEN, ["testserver"])
    monkeypatch.setenv("NATS_MONITOR_URLS", ",".join(f"http://n{i}:8222" for i in range(9)))
    with pytest.raises(ValueError, match="more than 8"):
        create_app(setup, TOKEN, ["testserver"])


def test_single_node_adapter_rejects_url_list(setup, monkeypatch):
    monkeypatch.setenv("PROMETHEUS_URLS", "http://p0:9090,http://p1:9090")
    with pytest.raises(ValueError, match="accepts one URL"):
        create_app(setup, TOKEN, ["testserver"])
