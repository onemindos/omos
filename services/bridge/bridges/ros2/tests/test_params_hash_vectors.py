"""Cross-implementation params_hash vectors shared with the policy service (omos #554) and
the MAVLink bridge (omos #550): same input, same hash (asset-command.md section 5).
The file lives at services/schemas/command/testdata/params-hash-vectors.json."""

import json
import pathlib

import pytest

from cmdtest import build, new_key
from ros2bridge import command as C
from ros2bridge import jcs

VECTOR_FILE = pathlib.Path(__file__).resolve().parents[4] / "schemas" / "command" / "testdata" / "params-hash-vectors.json"
VECTORS = {v["name"]: v for v in json.loads(VECTOR_FILE.read_text())["vectors"]}
PUB, PRIV = new_key()
TSK = "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8"


def test_vector_count():
    assert len(VECTORS) >= 10


@pytest.mark.parametrize("name", sorted(VECTORS))
def test_params_hash_vector(name):
    v = VECTORS[name]
    params = json.loads(v["params_json"]) if v["params_json"] else None
    t = v["target"]
    assert C.params_hash(v["action_type"], v["action_version"], t["entity_type"], t["entity_id"], params) == v["params_hash"]
    doc = {"action_type": v["action_type"], "action_version": v["action_version"], "target": t,
           "params": {} if params is None else params}
    assert jcs.canonicalize(doc).decode() == v["canonical"]


def verifier():
    return C.Verifier(C.Config(asset_class="robot", asset_id="sim-tb1", keys={"policy-test": PUB}))


GOTO = {"lat": 30.0001, "lon": -97.0001, "hae_m": 150.0}


@pytest.mark.parametrize("kw,code", [
    # D2-P1: a task-bearing command carries the task.assign / task.abort hash; accepted when signed.
    (dict(verb="goto", data=GOTO, task_id=TSK, params_hash=VECTORS["task-assign-goto"]["params_hash"]), None),
    (dict(verb="hold", data={"reason": "abort"}, task_id=TSK, params_hash=VECTORS["task-abort"]["params_hash"]), None),
    # Without task_id the asset.{verb} hash is recomputed, so a task hash is refused.
    (dict(verb="goto", data=GOTO, params_hash=VECTORS["task-assign-goto"]["params_hash"]), C.NOT_AUTHORIZED),
    (dict(verb="goto", data=GOTO, task_id=TSK, params_hash="sha256:XYZ"), C.SCHEMA_INVALID),
    (dict(verb="goto", data=GOTO, task_id="task 1", params_hash=VECTORS["task-assign-goto"]["params_hash"]), C.SCHEMA_INVALID),
    (dict(verb="goto", data=GOTO, task_id=TSK, params_hash=VECTORS["task-assign-goto"]["params_hash"], no_sig=True), C.BAD_SIGNATURE),
])
def test_task_bound_params_hash(kw, code):
    subj, body = build(PRIV, **kw)
    env = json.loads(body)
    if code is None:
        res = verifier().verify(subj, env, C.utcnow())
        assert res.task_bound and res.params_hash == kw["params_hash"]
    else:
        with pytest.raises(C.Refusal) as e:
            verifier().verify(subj, env, C.utcnow())
        assert e.value.code == code
