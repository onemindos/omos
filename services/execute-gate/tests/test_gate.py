import io
import json
import re

import pytest

from execute_gate.audit import Audit
from execute_gate.canonical import call_hash, canonical_json, summarise
from execute_gate.gate import Gate
from execute_gate.policy import Policy
from execute_gate.server import decide
from execute_gate.store import MemoryStore

OWNER = "@zeus:onemindos.dev"
BOT = "@gate:onemindos.dev"
ROOM = "!approvals:onemindos.dev"
DEL = "kubernetes/resources_delete"
DOOR = "homeassistant/lock_unlock"
ARGS = {"kind": "ConfigMap", "name": "gateway-selftest", "namespace": "tool-gateway"}


class Clock:
    def __init__(self):
        self.t = 1000.0

    def __call__(self):
        return self.t


class Notifier:
    def __init__(self):
        self.sent = []
        self.fail = False

    def send(self, text):
        if self.fail:
            raise OSError("down")
        self.sent.append(text)


class Bus:
    def __init__(self):
        self.events = []

    def __call__(self, subject, payload):
        self.events.append((subject, json.loads(payload)))


@pytest.fixture
def env():
    clock, notifier, bus, out = Clock(), Notifier(), Bus(), io.StringIO()
    gate = Gate(MemoryStore(), Policy(physical=("homeassistant/lock_*",)), OWNER, ROOM, BOT,
                notifier, Audit(bus, out, clock), clock=clock)
    return gate, clock, notifier, bus, out


def request_id(notifier):
    return re.search(r"\[([0-9A-F]{4})\]", notifier.sent[-1]).group(1)


def approve(gate, rid, word="execute", sender=OWNER, room=ROOM, ev="$e1"):
    return gate.handle_message(room, sender, ev, f"{word} {rid}")


# ---------------------------------------------------------------- canonicalisation

def test_canonical_is_order_and_whitespace_independent():
    a = {"b": 1, "a": {"y": [1, 2], "x": "é"}}
    b = json.loads('{ "a": {"x": "é", "y": [1, 2]}, "b": 1 }')
    assert canonical_json(a) == canonical_json(b) == '{"a":{"x":"é","y":[1,2]},"b":1}'
    assert call_hash(DEL, a) == call_hash(DEL, b)


def test_canonical_integral_float_equals_int_but_other_values_differ():
    assert call_hash(DEL, {"n": 1.0}) == call_hash(DEL, {"n": 1})
    assert call_hash(DEL, {"n": 1.5}) != call_hash(DEL, {"n": 1})
    assert call_hash(DEL, {"n": True}) != call_hash(DEL, {"n": 1})


def test_hash_binds_tool_and_arguments():
    assert call_hash(DEL, ARGS) != call_hash("kubernetes/resources_get", ARGS)
    assert call_hash(DEL, ARGS) != call_hash(DEL, {**ARGS, "name": "other"})


def test_canonical_rejects_nan_and_nonstring_keys():
    with pytest.raises(ValueError):
        canonical_json({"a": float("nan")})
    with pytest.raises(ValueError):
        canonical_json({1: "x"})


def test_summary_masks_secrets_and_truncates():
    s = summarise({"api_token": "hunter2", "name": "x" * 500})
    assert "hunter2" not in s and "[redacted]" in s and len(s) < 800


# ---------------------------------------------------------------- flow

def test_default_deny_creates_request_and_announces(env):
    gate, _, notifier, bus, out = env
    ok, reason, rid = gate.evaluate("hermes", DEL, ARGS)
    assert (ok, reason) == (False, "approval_required")
    assert rid and rid in notifier.sent[0] and "execute " + rid in notifier.sent[0]
    assert bus.events[0][0] == "evt.toolgw.approval.requested"
    assert "gateway-selftest" not in out.getvalue()  # argument values never audited


def test_pending_retry_does_not_reannounce(env):
    gate, _, notifier, *_ = env
    gate.evaluate("hermes", DEL, ARGS)
    ok, reason, _ = gate.evaluate("hermes", DEL, ARGS)
    assert (ok, reason) == (False, "pending_approval") and len(notifier.sent) == 1


def test_full_approval_then_single_use(env):
    gate, _, notifier, bus, _ = env
    gate.evaluate("hermes", DEL, ARGS)
    rid = request_id(notifier)
    assert approve(gate, rid) == "granted"
    assert gate.evaluate("hermes", DEL, dict(reversed(list(ARGS.items()))))[0] is True
    ok, reason, _ = gate.evaluate("hermes", DEL, ARGS)  # replay of the same call
    assert ok is False and reason == "approval_required"
    subjects = [s for s, _ in bus.events]
    assert subjects.count("evt.toolgw.approval.granted") == 1
    assert subjects.count("evt.toolgw.approval.consumed") == 1


def test_changed_argument_needs_new_approval(env):
    gate, _, notifier, *_ = env
    gate.evaluate("hermes", DEL, ARGS)
    approve(gate, request_id(notifier))
    assert gate.evaluate("hermes", DEL, {**ARGS, "name": "other"})[0] is False


def test_approved_id_cannot_be_reused_after_consumption(env):
    gate, _, notifier, *_ = env
    gate.evaluate("hermes", DEL, ARGS)
    rid = request_id(notifier)
    approve(gate, rid)
    gate.evaluate("hermes", DEL, ARGS)
    assert approve(gate, rid, ev="$e2") == "unknown_id"
    assert gate.evaluate("hermes", DEL, ARGS)[0] is False


def test_ttl_expiry_of_pending_and_approved(env):
    gate, clock, notifier, *_ = env
    gate.evaluate("hermes", DEL, ARGS)
    rid = request_id(notifier)
    clock.t += 601
    assert approve(gate, rid) == "expired"
    # approved but not consumed in time
    gate.evaluate("hermes", DEL, ARGS)
    rid2 = request_id(notifier)
    assert rid2 and approve(gate, rid2, ev="$e2") == "granted"
    clock.t += 601
    ok, reason, _ = gate.evaluate("hermes", DEL, ARGS)
    assert ok is False and reason == "approval_required"


def test_just_inside_ttl_still_works(env):
    gate, clock, notifier, *_ = env
    gate.evaluate("hermes", DEL, ARGS)
    clock.t += 599
    assert approve(gate, request_id(notifier)) == "granted"


def test_wrong_user_cannot_approve(env):
    gate, _, notifier, bus, _ = env
    gate.evaluate("hermes", DEL, ARGS)
    rid = request_id(notifier)
    for who in ("@legacy:onemindos.dev", BOT, "@zeus:evil.example", "@ZEUS:onemindos.dev"):
        assert approve(gate, rid, sender=who, ev="$" + who) == ("ignored" if who == BOT else "not_owner")
    assert gate.evaluate("hermes", DEL, ARGS)[0] is False
    assert any(s == "evt.toolgw.approval.denied" for s, _ in bus.events)


def test_wrong_room_ignored(env):
    gate, _, notifier, *_ = env
    gate.evaluate("hermes", DEL, ARGS)
    assert approve(gate, request_id(notifier), room="!other:onemindos.dev") == "ignored"
    assert gate.evaluate("hermes", DEL, ARGS)[0] is False


def test_wrong_or_unknown_id(env):
    gate, _, notifier, *_ = env
    gate.evaluate("hermes", DEL, ARGS)
    rid = request_id(notifier)
    other = "0000" if rid != "0000" else "FFFF"
    assert approve(gate, other) == "unknown_id"
    assert gate.evaluate("hermes", DEL, ARGS)[0] is False
    assert approve(gate, rid, ev="$ok") == "granted"


def test_physical_needs_override_destructive_needs_execute(env):
    gate, _, notifier, *_ = env
    gate.evaluate("hermes", DEL, ARGS)
    d = request_id(notifier)
    assert approve(gate, d, word="override") == "wrong_word"
    assert gate.evaluate("hermes", DEL, ARGS)[0] is False

    gate.evaluate("hermes", DOOR, {"entity": "lock.front"})
    p = request_id(notifier)
    assert "override " + p in notifier.sent[-1]
    assert approve(gate, p, word="execute", ev="$2") == "wrong_word"
    assert gate.evaluate("hermes", DOOR, {"entity": "lock.front"})[0] is False
    assert approve(gate, p, word="override", ev="$3") == "granted"
    assert gate.evaluate("hermes", DOOR, {"entity": "lock.front"})[0] is True


def test_unclassified_tool_is_destructive(env):
    gate, *_ = env
    assert gate.policy.classify("anything/new_tool") == "destructive"
    assert gate.policy.classify(DOOR) == "physical"


def test_event_replay_is_ignored(env):
    gate, _, notifier, *_ = env
    gate.evaluate("hermes", DEL, ARGS)
    rid = request_id(notifier)
    assert approve(gate, rid, ev="$same") == "granted"
    assert approve(gate, rid, ev="$same") == "replay"


def test_edited_message_and_command_shape(env):
    gate, _, notifier, *_ = env
    gate.evaluate("hermes", DEL, ARGS)
    rid = request_id(notifier)
    assert gate.handle_message(ROOM, OWNER, "$1", f"* execute {rid}", edited=True) == "ignored"
    assert gate.handle_message(ROOM, OWNER, "$2", f"execute {rid}", edited=True) == "not_plain_message"
    assert gate.handle_message(ROOM, OWNER, "$3", f"please execute {rid}") == "ignored"
    assert gate.handle_message(ROOM, OWNER, "$4", f"execute {rid} now") == "ignored"
    assert gate.handle_message(ROOM, OWNER, "$5", f"  EXECUTE {rid.lower()} ") == "granted"


def test_disallowed_subject_and_bad_input_denied(env):
    gate, *_ = env
    assert gate.evaluate("legacy", DEL, ARGS)[:2] == (False, "subject_not_allowed")
    assert gate.evaluate("hermes", DEL, "{not json")[:2] == (False, "bad_arguments")
    assert gate.evaluate("hermes", "", ARGS)[:2] == (False, "bad_request")


def test_arguments_as_json_text_hash_same_as_object(env):
    gate, _, notifier, *_ = env
    gate.evaluate("hermes", DEL, json.dumps(ARGS))
    approve(gate, request_id(notifier))
    assert gate.evaluate("hermes", DEL, ARGS)[0] is True


def test_announce_failure_retries_on_next_call_and_store_failure_denies(env):
    gate, _, notifier, *_ = env
    notifier.fail = True
    assert gate.evaluate("hermes", DEL, ARGS)[:2] == (False, "approval_required")
    notifier.fail = False
    assert gate.evaluate("hermes", DEL, ARGS)[:2] == (False, "pending_approval")
    assert len(notifier.sent) == 1

    class Broken(MemoryStore):
        def get(self, key):
            raise RuntimeError("kv down")

    gate.store = Broken()
    assert gate.evaluate("hermes", DEL, ARGS)[:2] == (False, "error")


# ---------------------------------------------------------------- AuthZEN wrapper

def test_decide_authzen_shape(env):
    gate, *_ = env
    body = {"subject": {"type": "client", "id": "hermes"}, "action": {"name": "tools/call"},
            "resource": {"type": "tool", "id": DEL}, "context": {"arguments": ARGS}}
    r = decide(gate, body)
    assert r["decision"] is False and r["context"]["request_id"]
    assert decide(gate, {**body, "action": {"name": "tools/list"}})["decision"] is False
    assert decide(gate, {})["decision"] is False
