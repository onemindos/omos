"""Builds asset command envelopes the way the policy service does (asset-command.md
section 6.3), for tests only. Python counterpart of #550's internal/commandtest; keys are
generated per run and none is stored in the repository."""

import json
import uuid
from datetime import datetime, timedelta, timezone

from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey

from ros2bridge import command as C


def rfc(t: datetime) -> str:
    u = t.astimezone(timezone.utc)
    return u.strftime("%Y-%m-%dT%H:%M:%S.") + f"{u.microsecond // 1000:03d}Z"


def new_key():
    k = Ed25519PrivateKey.generate()
    return k.public_key(), k


_UNSET = object()


def build(key, verb="rtl", asset="sim-tb1", data=None, now=None, ttl=timedelta(seconds=30), level=None,
          approvers=None, status="approved", simulated=True, source=C.POLICY_SOURCE, hash_data=_UNSET,
          cid=None, task_id="", no_sig=False, kid="policy-test", mutate=None, cls="robot", subject_verb=None):
    now = now or datetime.now(timezone.utc)
    data = {} if data is None else data
    hd = data if hash_data is _UNSET else hash_data
    v = C.VERBS.get(verb)
    level = level or (v.level if v else "A3")
    if approvers is None:
        approvers = ["user:zeus"] + (["user:alice"] if v and v.approvals == 2 else [])
    subject = f"cmd.fleet.{cls}.{asset}.{subject_verb or verb}"
    env = {
        "v": 1, "id": cid or str(uuid.uuid4()), "ts": rfc(now), "cmd": subject, "source": source,
        "subject": subject, "data": data, "action_id": "act-" + uuid.uuid4().hex[:8],
        "expires_at": rfc(now + ttl), "simulated": simulated,
        "approval": {"required": True, "status": status, "approver": approvers[-1] if approvers else "",
                     "reason": "", "level": level,
                     "params_hash": C.params_hash("asset." + verb, 1, cls, asset, hd),
                     "approvals": [{"by": a, "at": rfc(now - timedelta(seconds=1)), "step": i + 1}
                                   for i, a in enumerate(approvers)]},
        "policy": {"decision_id": "dec-" + uuid.uuid4().hex[:8], "allowed": True},
        "corr": str(uuid.uuid4()),
    }
    if task_id:
        env["task_id"] = task_id
    if mutate:
        mutate(env)
    if not no_sig:
        env = C.sign(env, kid, key)
    return subject, json.dumps(env).encode()
