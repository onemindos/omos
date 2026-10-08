"""Real PostgreSQL checks for replay and transactional event persistence."""
import os
import uuid
import json
from pathlib import Path
from concurrent.futures import ThreadPoolExecutor
import pytest
from jsonschema import Draft202012Validator, FormatChecker
from psycopg_pool import ConnectionPool
from psycopg.rows import dict_row
from onemind import store
from onemind.events import envelope

pytestmark = pytest.mark.skipif(not os.environ.get("INTEGRATION_DATABASE_URL"), reason="Integration PostgreSQL not configured")


@pytest.fixture
def database(monkeypatch):
    pool = ConnectionPool(os.environ["INTEGRATION_DATABASE_URL"], open=False,
                          kwargs={"row_factory": dict_row})
    monkeypatch.setattr(store, "pool", pool)
    store.initialize()
    with pool.connection() as conn:
        conn.execute("TRUNCATE action_audit, event_outbox, missions, links, objects RESTART IDENTITY")
    yield pool
    pool.close()


def test_concurrent_request_has_one_mission_and_one_committed_event(database):
    def request(_):
        return store.create_mission(str(uuid.uuid4()), "user:default/operator", "same-request", "digest", {"capability": "inspect-system"})
    with ThreadPoolExecutor(max_workers=4) as workers:
        rows = list(workers.map(request, range(8)))
    assert len({row[0]["id"] for row in rows}) == 1
    assert sum(row[1] for row in rows) == 1
    assert len(store.pending_events()) == 1


def test_cancel_cannot_be_overwritten_by_a_finishing_activity(database):
    mid = str(uuid.uuid4())
    store.create_mission(mid, "user:default/operator", "cancel", "digest", {})
    store.set_state(mid, "running")
    store.set_state(mid, "cancelled")
    store.set_state(mid, "completed", result={"answer": "late completion"})
    assert store.mission(mid)["state"] == "cancelled"
    assert [x["payload"]["state"] for x in store.pending_events()] == ["requested", "running", "cancelled"]


def test_external_dispatch_can_only_be_claimed_once(database):
    mid = str(uuid.uuid4())
    with ThreadPoolExecutor(max_workers=4) as workers:
        claims = list(workers.map(lambda _: store.claim_dispatch("user:default/operator", "discovered/ref", mid), range(8)))
    assert sum(claims) == 1


def test_outbox_matches_existing_envelope_and_new_payload_contract(database):
    mid = str(uuid.uuid4())
    store.create_mission(mid, "user:default/operator", "contract", "digest", {})
    row = store.pending_events()[0]
    event = envelope(row)
    root = Path(__file__).resolve().parents[3]
    for file, value in [("schemas/envelope/envelope.schema.json", event),
                        ("schemas/missions/mission-state.schema.json", event["data"])]:
        Draft202012Validator(json.loads((root / file).read_text()), format_checker=FormatChecker()).validate(value)
    assert event["corr"] == mid
    assert event["subject"] == "evt.so.mission.requested"
    store.event_sent(row["id"])
    assert not store.pending_events()
