import hashlib
import json
import pytest
from fastapi.testclient import TestClient
from onemind.api import app
from onemind import store
from onemind.workflows import execute_capability


@pytest.fixture
def client(monkeypatch):
    monkeypatch.setenv("UI_GATEWAY_KEY", "test-ui-gateway")
    monkeypatch.setenv("READ_API_KEY", "test-reader")
    return TestClient(app)


UI = {"Authorization": "Bearer test-ui-gateway", "X-OneMind-Actor": "user:default/operator"}
READ = {"Authorization": "Bearer test-reader"}


def test_anonymous_cannot_inspect_capabilities(client):
    assert client.get("/capabilities").status_code == 401


def test_read_credential_cannot_mutate_world(client):
    response = client.post("/objects", headers=READ,
                           json={"id": "asset:pump", "kind": "asset", "name": "Pump"})
    assert response.status_code == 403


def test_gateway_requires_asserted_user_identity(client):
    assert client.get("/capabilities", headers={"Authorization": "Bearer test-ui-gateway"}).status_code == 403


def test_unknown_action_never_reaches_execution(client):
    response = client.post("/missions", headers=UI,
                           json={"capability": "arbitrary-shell", "idempotency_key": "test"})
    assert response.status_code == 422


def test_restart_cannot_target_existing_identity_service(client):
    response = client.post("/missions", headers=UI, json={"capability": "restart-worker",
                           "target": "auth/platform-auth-helm-authentik-server", "idempotency_key": "test"})
    assert response.status_code == 403


def test_idempotency_key_cannot_be_reused_for_other_work(client, monkeypatch):
    monkeypatch.setattr(store, "create_mission", lambda *args: ({"request_hash": "different-request"}, False))
    response = client.post("/missions", headers=UI,
                           json={"capability": "inspect-system", "idempotency_key": "same-key"})
    assert response.status_code == 409


def test_identical_request_returns_existing_mission(client, monkeypatch):
    spec = {"capability": "inspect-system", "prompt": "Inspect current OneMind conditions", "target": ""}
    digest = hashlib.sha256(json.dumps(spec, sort_keys=True).encode()).hexdigest()
    monkeypatch.setattr(store, "create_mission", lambda *args: ({"id": "existing", "request_hash": digest}, False))
    monkeypatch.setattr(store, "audit", lambda *args: pytest.fail("Duplicate request must not record another action"))
    response = client.post("/missions", headers=UI,
                           json={"capability": "inspect-system", "idempotency_key": "same-key"})
    assert response.status_code == 202
    assert response.json()["id"] == "existing"


def test_completed_external_work_is_not_reexecuted(monkeypatch):
    monkeypatch.setattr(store, "mission", lambda mid: {"state": "completed", "result": {"answer": "done"}})
    assert execute_capability.__wrapped__(None, {"id": "mission"}) == {"answer": "done"}


def test_ambiguous_agent_dispatch_is_not_repeated(monkeypatch):
    monkeypatch.setattr(store, "mission", lambda mid: {"state": "running"})
    monkeypatch.setattr(store, "dispatched", lambda mid: True)
    with pytest.raises(RuntimeError, match="ambiguous"):
        execute_capability.__wrapped__(None, {"id": "mission", "spec": {"capability": "delegate-agent"}})


def test_cancelled_mission_cannot_start_an_activity(monkeypatch):
    monkeypatch.setattr(store, "mission", lambda mid: {"state": "cancelled"})
    with pytest.raises(RuntimeError, match="cancelled"):
        execute_capability.__wrapped__(None, {"id": "mission"})
