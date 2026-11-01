"""Verifier tests: the table from #550's command_test.go adapted to the robot catalogue,
plus envelopes signed by the policy service code (fixtures/policy_vectors.json)."""

import base64
import json
import os
from datetime import timedelta

import pytest
from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PublicKey

from cmdtest import build, new_key
from ros2bridge import command as C

PUB, PRIV = new_key()
_, OTHER = new_key()
NOW = C.utcnow()
GOOD_GOTO = {"lat": 30.0001, "lon": -97.0001, "hae_m": 150.0, "speed_mps": 0.2}


def verifier(**kw):
    cfg = dict(asset_class="robot", asset_id="sim-tb1", keys={"policy-test": PUB})
    cfg.update(kw)
    return C.Verifier(C.Config(**cfg))


def test_refuses_non_sim_asset():
    with pytest.raises(ValueError):
        verifier(asset_id="tb1")
    with pytest.raises(ValueError):
        C.Verifier(C.Config("robot", "sim-tb1", {}, sig_required=True))


def _tamper(env):
    env["data"] = {"lat": 1}


CASES = [
    ("goto approved", dict(verb="goto", data=GOOD_GOTO), None, None),
    ("goto with yaw", dict(verb="goto", data={"lat": 0.0, "lon": 0.0, "hae_m": 0, "yaw_deg": 359.9}), None, None),
    ("hold", dict(verb="hold", data={"duration_s": 30, "reason": "wait"}), None, None),
    ("hold empty", dict(verb="hold"), None, None),
    ("rtl", dict(verb="rtl"), None, None),
    ("estop", dict(verb="estop"), None, None),
    ("estop_release two humans", dict(verb="estop_release"), None, None),
    ("estop_release one human", dict(verb="estop_release", approvers=["user:zeus"]), C.NOT_AUTHORIZED, None),
    ("estop_release at A3", dict(verb="estop_release", level="A3"), C.NOT_AUTHORIZED, None),
    ("A4 level needs two humans even for rtl", dict(verb="rtl", level="A4"), C.NOT_AUTHORIZED, None),
    ("pending approval", dict(verb="rtl", status="pending"), C.NOT_AUTHORIZED, None),
    ("no approvals", dict(verb="rtl", approvers=[]), C.NOT_AUTHORIZED, None),
    ("agent approver", dict(verb="rtl", approvers=["agent:hermes/zeus/legacy"]), C.AGENT_CANNOT_CONFIRM, None),
    ("same human twice", dict(verb="estop_release", approvers=["user:zeus", "user:zeus"]), C.NOT_AUTHORIZED, None),
    ("not simulated", dict(verb="rtl", simulated=False), C.NOT_SIMULATED, None),
    ("simulated missing", dict(verb="rtl", mutate=lambda e: e.pop("simulated")), C.NOT_SIMULATED, None),
    ("non-sim subject", dict(verb="rtl"), C.NOT_SIMULATED, "cmd.fleet.robot.tb1.rtl"),
    ("other sim asset", dict(verb="rtl"), C.SCHEMA_INVALID, "cmd.fleet.robot.sim-tb2.rtl"),
    ("drone class subject", dict(verb="rtl"), C.SCHEMA_INVALID, "cmd.fleet.drone.sim-tb1.rtl"),
    ("unknown verb", dict(verb="strike"), C.UNKNOWN_ACTION, None),
    ("drone-only verb", dict(verb="takeoff"), C.UNKNOWN_ACTION, None),
    ("wrong source", dict(verb="rtl", source="agent.legacy"), C.NOT_AUTHORIZED, None),
    ("unsigned", dict(verb="rtl", no_sig=True), C.BAD_SIGNATURE, None),
    ("wrong key", dict(verb="rtl", key=OTHER), C.BAD_SIGNATURE, None),
    ("unknown kid", dict(verb="rtl", kid="nope"), C.BAD_SIGNATURE, None),
    ("expired", dict(verb="rtl", now=NOW - timedelta(minutes=1)), C.EXPIRED, None),
    ("ttl too long", dict(verb="rtl", ttl=timedelta(minutes=10)), C.EXPIRED, None),
    ("future ts", dict(verb="rtl", now=NOW + timedelta(minutes=1)), C.EXPIRED, None),
    ("hash mismatch", dict(verb="goto", data=GOOD_GOTO, hash_data={"lat": 1.0, "lon": 1.0, "hae_m": 80.0}),
     C.NOT_AUTHORIZED, None),
    ("bad params", dict(verb="goto", data={"lat": 91.0, "lon": 0.0, "hae_m": 10.0}), C.SCHEMA_INVALID, None),
    ("goto without hae", dict(verb="goto", data={"lat": 1.0, "lon": 0.0}), C.SCHEMA_INVALID, None),
    ("speed too high", dict(verb="goto", data={"lat": 1.0, "lon": 0.0, "hae_m": 0, "speed_mps": 51}), C.SCHEMA_INVALID, None),
    ("hold bad duration", dict(verb="hold", data={"duration_s": 1.5}), C.SCHEMA_INVALID, None),
    ("extra param", dict(verb="rtl", data={"payload": "drop"}), C.SCHEMA_INVALID, None),
    ("data not object", dict(verb="rtl", data=[1]), C.SCHEMA_INVALID, None),
    ("policy denied", dict(verb="rtl", mutate=lambda e: e.update(policy={"decision_id": "d", "allowed": False})),
     C.NOT_AUTHORIZED, None),
    ("bad id", dict(verb="rtl", cid="NOT-A-UUID"), C.SCHEMA_INVALID, None),
    ("bad ts", dict(verb="rtl", mutate=lambda e: e.update(ts="yesterday")), C.SCHEMA_INVALID, None),
    ("tampered after signing", dict(verb="rtl"), C.BAD_SIGNATURE, "tamper"),
]


@pytest.mark.parametrize("name,opts,want,subject", CASES, ids=[c[0] for c in CASES])
def test_verify(name, opts, want, subject):
    opts = dict(opts)
    key = opts.pop("key", PRIV)
    opts.setdefault("now", NOW)
    subj, body = build(key, **opts)
    if subject == "tamper":
        env = json.loads(body)
        _tamper(env)
        body = json.dumps(env).encode()
    elif subject:
        subj = subject
    v = verifier()
    if want is None:
        res = v.verify(subj, json.loads(body), NOW)
        assert res.verb.name == opts["verb"] and res.signed
    else:
        with pytest.raises(C.Refusal) as ei:
            v.verify(subj, json.loads(body), NOW)
        assert ei.value.code == want, ei.value.message


def test_replay_refused():
    v = verifier()
    subj, body = build(PRIV, verb="rtl", now=NOW)
    v.verify(subj, json.loads(body), NOW)
    with pytest.raises(C.Refusal) as ei:
        v.verify(subj, json.loads(body), NOW + timedelta(seconds=1))
    assert ei.value.code == C.REPLAYED
    # after expiry the id is forgotten, but then the command itself is expired
    with pytest.raises(C.Refusal) as ei:
        v.verify(subj, json.loads(body), NOW + timedelta(seconds=40))
    assert ei.value.code == C.EXPIRED


def test_sig_optional_when_not_required():
    v = verifier(keys={}, sig_required=False)
    subj, body = build(PRIV, verb="rtl", now=NOW, no_sig=True)
    assert v.verify(subj, json.loads(body), NOW).signed is False


def test_params_hash_is_bound_to_target():
    a = C.params_hash("asset.goto", 1, "robot", "sim-tb1", GOOD_GOTO)
    assert a != C.params_hash("asset.goto", 1, "drone", "sim-tb1", GOOD_GOTO)
    assert a != C.params_hash("asset.goto", 1, "robot", "sim-tb2", GOOD_GOTO)
    assert C.params_hash("asset.rtl", 1, "robot", "sim-tb1", None) == C.params_hash("asset.rtl", 1, "robot", "sim-tb1", {})


def test_parse_keys():
    raw = PUB.public_bytes_raw()
    keys = C.parse_keys(f"policy-1={base64.b64encode(raw).decode()}, ")
    assert list(keys) == ["policy-1"]
    for bad in ("nokid", "k=notbase64!", "k=" + base64.b64encode(b"short").decode()):
        with pytest.raises(ValueError):
            C.parse_keys(bad)


def test_parse_time_accepts_go_nanoseconds():
    assert C.parse_time("2026-10-29T12:00:00.123456789Z").microsecond == 123456
    assert C.parse_time("2026-10-29T12:00:00+02:00").hour == 12
    for bad in ("2026-10-29", "2026-10-29 12:00:00Z", 5, None):
        with pytest.raises(ValueError):
            C.parse_time(bad)


# ---- envelopes signed by the policy service code (PR #554) ----
VEC = json.load(open(os.path.join(os.path.dirname(__file__), "fixtures", "policy_vectors.json")))
POLICY_PUB = Ed25519PublicKey.from_public_bytes(base64.b64decode(VEC["public_key_b64"]))
VNOW = C.parse_time(VEC["now"])


@pytest.mark.parametrize("vec", VEC["vectors"], ids=lambda v: v["name"])
def test_policy_signed_vectors_verify(vec):
    v = C.Verifier(C.Config("robot", "sim-tb1", {VEC["kid"]: POLICY_PUB}))
    res = v.verify(vec["subject"], vec["envelope"], VNOW + timedelta(seconds=1))
    assert res.signed and res.params_hash == vec["params_hash"]


def test_policy_vector_tampered_fails():
    vec = VEC["vectors"][0]
    env = json.loads(json.dumps(vec["envelope"]))
    env["data"]["lat"] = 30.3
    v = C.Verifier(C.Config("robot", "sim-tb1", {VEC["kid"]: POLICY_PUB}))
    with pytest.raises(C.Refusal) as ei:
        v.verify(vec["subject"], env, VNOW)
    assert ei.value.code == C.BAD_SIGNATURE
