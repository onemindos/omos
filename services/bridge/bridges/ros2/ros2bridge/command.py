"""Verification of asset commands on ``cmd.fleet.robot.{asset_id}.{verb}``.

A Python port of services/bridge/bridges/mavlink/internal/command (omos PR #550): the
same checks, in the same order, with the same refusal codes (onemind-xr
docs/contracts/asset-command.md section 4 step 6, section 6.3 and section 7). Only the
verb catalogue and parameter rules differ, because a ground robot has no arm, takeoff
or land. tests/test_command.py ports #550's table test and also verifies envelopes
signed by the policy service code (PR #554), so the two implementations cannot drift
silently. A command is accepted only if every check passes; anything else is refused
and nothing is sent to the robot.
"""

from __future__ import annotations

import base64
import hashlib
import json
import math
import re
import threading
from dataclasses import dataclass, field
from datetime import datetime, timedelta, timezone

from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey, Ed25519PublicKey

from . import jcs

# Refusal codes (asset-command.md section 6.4), identical to the MAVLink bridge.
SCHEMA_INVALID = "schema_invalid"
UNKNOWN_ACTION = "unknown_action"
NOT_SIMULATED = "not_simulated"
NOT_AUTHORIZED = "not_authorized"
AGENT_CANNOT_CONFIRM = "agent_cannot_confirm"
BAD_SIGNATURE = "bad_signature"
EXPIRED = "expired"
REPLAYED = "replayed"
PRECONDITION = "precondition_failed"
VEHICLE_REJECTED = "vehicle_rejected"
NO_ACK = "no_ack"
CONFLICT = "conflict"

POLICY_SOURCE = "service:policy"


@dataclass(frozen=True)
class Verb:
    name: str
    level: str
    approvals: int


# The verbs this bridge executes (asset-command.md section 3). Any other verb is refused.
VERBS = {
    "goto": Verb("goto", "A3", 1),
    "hold": Verb("hold", "A3", 1),
    "rtl": Verb("rtl", "A3", 1),
    "estop": Verb("estop", "A3", 1),
    "estop_release": Verb("estop_release", "A4", 2),
}

# Allowed parameter members per verb (services/schemas/actions/asset.{verb}.v1.schema.json).
_ALLOWED = {
    "goto": ("lat", "lon", "hae_m", "speed_mps", "yaw_deg"),
    "hold": ("duration_s", "reason"),
}


class Refusal(Exception):
    def __init__(self, code: str, message: str, detail: str = ""):
        super().__init__(f"{code}: {message}")
        self.code, self.message, self.detail = code, message, detail


def _level_rank(level) -> int:
    if isinstance(level, str) and len(level) == 2 and level[0] == "A" and level[1] in "01234":
        return int(level[1])
    return -1


_UUID = re.compile(r"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")
_RFC3339 = re.compile(r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,9})?(Z|[+-]\d{2}:\d{2})$")


def parse_time(s) -> datetime:
    if not isinstance(s, str) or not _RFC3339.match(s):
        raise ValueError(f"not an RFC 3339 timestamp: {s!r}")
    head, _, rest = s.partition(".")
    if rest:  # Python accepts at most 6 fractional digits; Go accepts 9.
        frac = re.match(r"\d+", rest).group(0)
        s = head + "." + frac[:6].ljust(6, "0") + rest[len(frac):]
    return datetime.fromisoformat(s.replace("Z", "+00:00"))


@dataclass
class Params:
    lat: float | None = None
    lon: float | None = None
    hae_m: float | None = None
    speed_mps: float | None = None
    yaw_deg: float | None = None
    duration_s: int | None = None
    reason: str | None = None


def _finite(v) -> bool:
    return isinstance(v, (int, float)) and not isinstance(v, bool) and math.isfinite(v)


def parse_params(verb: str, data) -> Params:
    if not isinstance(data, dict):
        raise Refusal(SCHEMA_INVALID, "data must be an object")
    allowed = _ALLOWED.get(verb, ())
    for k in data:
        if k not in allowed:
            raise Refusal(SCHEMA_INVALID, f"unexpected parameter {k!r} for {verb}")
    p = Params(**{k: data[k] for k in data})
    for k in ("lat", "lon", "hae_m", "speed_mps", "yaw_deg"):
        v = getattr(p, k)
        if v is not None and not _finite(v):
            raise Refusal(SCHEMA_INVALID, f"{k} must be a finite number")
    if verb == "goto":
        if p.lat is None or p.lon is None or not (-90 <= p.lat <= 90) or not (-180 <= p.lon <= 180):
            raise Refusal(SCHEMA_INVALID, "lat/lon out of range")
        if p.hae_m is None:
            raise Refusal(SCHEMA_INVALID, "hae_m required")
        if p.speed_mps is not None and not (0 < p.speed_mps <= 50):
            raise Refusal(SCHEMA_INVALID, "speed_mps must be in (0, 50]")
        if p.yaw_deg is not None and not (0 <= p.yaw_deg < 360):
            raise Refusal(SCHEMA_INVALID, "yaw_deg must be in [0, 360)")
    elif verb == "hold":
        d = p.duration_s
        if d is not None and (isinstance(d, bool) or not isinstance(d, int) or not 1 <= d <= 3600):
            raise Refusal(SCHEMA_INVALID, "duration_s must be an integer in [1, 3600]")
        if p.reason is not None and (not isinstance(p.reason, str) or len(p.reason) > 100):
            raise Refusal(SCHEMA_INVALID, "reason must be a string of at most 100 characters")
    return p


def params_hash(action_type: str, action_version: int, cls: str, asset_id: str, params) -> str:
    """"sha256:" + hex SHA-256 of the JCS form of {action_type, action_version, target, params}."""
    doc = {"action_type": action_type, "action_version": action_version,
           "target": {"entity_type": cls, "entity_id": asset_id},
           "params": {} if params is None else params}
    return "sha256:" + hashlib.sha256(jcs.canonicalize(doc)).hexdigest()


def signing_bytes(env: dict) -> bytes:
    return jcs.canonicalize({k: v for k, v in env.items() if k != "sig"})


def sign(env: dict, kid: str, key: Ed25519PrivateKey) -> dict:
    """Return a copy of env with an ed25519 ``sig``. Tests only: the policy owns the key."""
    out = {k: v for k, v in env.items() if k != "sig"}
    out["sig"] = {"alg": "ed25519", "kid": kid,
                  "value": base64.b64encode(key.sign(signing_bytes(out))).decode()}
    return out


@dataclass
class Config:
    asset_class: str
    asset_id: str
    keys: dict = field(default_factory=dict)  # kid -> Ed25519PublicKey
    sig_required: bool = True
    max_ttl: timedelta = timedelta(seconds=60)
    clock_skew: timedelta = timedelta(seconds=5)


@dataclass
class Result:
    verb: Verb
    params: Params
    signed: bool
    params_hash: str
    expires_at: datetime


class Verifier:
    """Checks commands and remembers command ids (nonces) until they expire."""

    def __init__(self, cfg: Config):
        if not cfg.asset_id.startswith("sim-"):
            raise ValueError(f"asset id {cfg.asset_id!r} does not start with sim- (decision 0004 #13)")
        if cfg.sig_required and not cfg.keys:
            raise ValueError("signatures required but no policy public keys configured")
        self.cfg = cfg
        self._seen: dict[str, datetime] = {}
        self._mu = threading.Lock()

    def verify(self, subject: str, env, now: datetime) -> Result:
        cfg = self.cfg
        toks = subject.split(".")
        if len(toks) != 5 or toks[0] != "cmd" or toks[1] != "fleet" or toks[2] != cfg.asset_class:
            raise Refusal(SCHEMA_INVALID, f"subject {subject!r} is not cmd.fleet.{cfg.asset_class}.{{id}}.{{verb}}")
        asset, verb_name = toks[3], toks[4]
        if not isinstance(env, dict):
            raise Refusal(SCHEMA_INVALID, "envelope must be an object")

        # Rule 1: simulated assets only.
        if not asset.startswith("sim-"):
            raise Refusal(NOT_SIMULATED, f"target {asset!r} is not a sim- asset")
        if asset != cfg.asset_id:
            raise Refusal(SCHEMA_INVALID, f"target {asset!r} is not this bridge's asset")
        if env.get("simulated") is not True:
            raise Refusal(NOT_SIMULATED, "simulated must be true")
        verb = VERBS.get(verb_name)
        if verb is None:
            raise Refusal(UNKNOWN_ACTION, f"verb {verb_name!r} is not in the bridge catalogue")

        # Envelope shape (omos command envelope v1 required fields + section 6.3).
        v = env.get("v")
        if not isinstance(v, int) or isinstance(v, bool) or v < 1:
            raise Refusal(SCHEMA_INVALID, "v must be >= 1")
        if not isinstance(env.get("id"), str) or not _UUID.match(env["id"]):
            raise Refusal(SCHEMA_INVALID, "id must be a lowercase UUID")
        if env.get("cmd") != subject or env.get("subject") != subject:
            raise Refusal(SCHEMA_INVALID, "cmd and subject must equal the delivery subject")
        if env.get("source") != POLICY_SOURCE:
            raise Refusal(NOT_AUTHORIZED, f"source {env.get('source')!r} is not {POLICY_SOURCE}")
        if not isinstance(env.get("action_id"), str) or not env["action_id"].startswith("act-"):
            raise Refusal(SCHEMA_INVALID, "action_id missing")
        try:
            ts = parse_time(env.get("ts"))
        except ValueError as e:
            raise Refusal(SCHEMA_INVALID, f"ts: {e}") from None
        try:
            exp = parse_time(env.get("expires_at"))
        except ValueError as e:
            raise Refusal(SCHEMA_INVALID, f"expires_at: {e}") from None

        # Rule 2: signature.
        signed = False
        if env.get("sig") is not None:
            self._check_sig(env)
            signed = True
        elif cfg.sig_required:
            raise Refusal(BAD_SIGNATURE, "unsigned command")

        # Rule 2: expiry and freshness.
        if not now < exp:
            raise Refusal(EXPIRED, f"expired at {env['expires_at']}")
        if ts > now + cfg.clock_skew:
            raise Refusal(EXPIRED, f"ts {env['ts']} is in the future")
        if exp - ts > cfg.max_ttl or not exp > ts:
            raise Refusal(EXPIRED, f"expires_at - ts must be in (0, {cfg.max_ttl}]")

        # Policy decision and human approval.
        pol = env.get("policy")
        if not isinstance(pol, dict) or pol.get("allowed") is not True or not pol.get("decision_id"):
            raise Refusal(NOT_AUTHORIZED, "no allowing policy decision")
        a = env.get("approval")
        if not isinstance(a, dict) or a.get("required") is not True or a.get("status") != "approved":
            raise Refusal(NOT_AUTHORIZED, "command is not approved")
        if _level_rank(a.get("level")) < _level_rank(verb.level):
            raise Refusal(NOT_AUTHORIZED, f"level {a.get('level')!r} below {verb.level} required for {verb.name}")
        need = verb.approvals
        if _level_rank(a.get("level")) >= 4 and need < 2:
            need = 2
        humans: set = set()
        steps: set = set()
        approver_listed = False
        for s in a.get("approvals") or []:
            by = s.get("by") if isinstance(s, dict) else None
            if not isinstance(by, str) or not by.startswith("user:"):
                raise Refusal(AGENT_CANNOT_CONFIRM, f"approval by non-human principal {by!r}")
            try:
                at = parse_time(s.get("at"))
            except ValueError:
                raise Refusal(NOT_AUTHORIZED, f"approval by {by} has a bad time") from None
            if at > ts + cfg.clock_skew:
                raise Refusal(NOT_AUTHORIZED, f"approval by {by} has a bad time")
            step = s.get("step")
            if by in humans or step in steps:
                raise Refusal(NOT_AUTHORIZED, f"duplicate approval step or approver {by}")
            humans.add(by)
            steps.add(step)
            if by == a.get("approver"):
                approver_listed = True
        if len(humans) < need:
            raise Refusal(NOT_AUTHORIZED, f"{verb.name} needs {need} distinct human approvals, has {len(humans)}")
        if not approver_listed:
            raise Refusal(NOT_AUTHORIZED, f"approver {a.get('approver')!r} not in approvals")

        # Parameters and the approved hash (approvals are bound to params_hash, section 4).
        data = env.get("data")
        params = parse_params(verb.name, data)
        av = env.get("action_version", 1)
        if not isinstance(av, int) or isinstance(av, bool):
            raise Refusal(SCHEMA_INVALID, "action_version must be an integer")
        try:
            h = params_hash("asset." + verb.name, av, cfg.asset_class, asset, data)
        except (TypeError, ValueError) as e:
            raise Refusal(SCHEMA_INVALID, f"params_hash: {e}") from None
        if a.get("params_hash") != h:
            raise Refusal(NOT_AUTHORIZED, "params_hash does not match the approved parameters")

        # Nonce: refuse a command id seen before; remember it until it expires.
        self._remember(env["id"], exp, now)
        return Result(verb, params, signed, h, exp)

    def _check_sig(self, env: dict) -> None:
        sig = env["sig"]
        if not isinstance(sig, dict) or sig.get("alg") != "ed25519":
            raise Refusal(BAD_SIGNATURE, f"alg {sig.get('alg') if isinstance(sig, dict) else sig!r} not supported")
        key = self.cfg.keys.get(sig.get("kid"))
        if key is None:
            raise Refusal(BAD_SIGNATURE, f"unknown kid {sig.get('kid')!r}")
        try:
            raw = base64.b64decode(sig.get("value") or "", validate=True)
        except (ValueError, TypeError):
            raw = b""
        if len(raw) != 64:
            raise Refusal(BAD_SIGNATURE, "malformed signature")
        try:
            msg = signing_bytes(env)
        except (TypeError, ValueError) as e:
            raise Refusal(BAD_SIGNATURE, f"canonicalize: {e}") from None
        try:
            key.verify(raw, msg)
        except InvalidSignature:
            raise Refusal(BAD_SIGNATURE, "signature does not verify") from None

    def _remember(self, cid: str, exp: datetime, now: datetime) -> None:
        with self._mu:
            for k in [k for k, t in self._seen.items() if now > t]:
                del self._seen[k]
            if cid in self._seen:
                raise Refusal(REPLAYED, f"command id {cid} already received")
            self._seen[cid] = exp


def parse_keys(s: str) -> dict:
    """POLICY_PUBKEYS: comma-separated kid=base64(raw 32-byte ed25519 public key)."""
    keys = {}
    for part in (s or "").split(","):
        part = part.strip()
        if not part:
            continue
        kid, sep, b64 = part.partition("=")
        if not sep or not kid:
            raise ValueError(f"POLICY_PUBKEYS: want kid=base64, got {part!r}")
        try:
            raw = base64.b64decode(b64, validate=True)
        except ValueError:
            raw = b""
        if len(raw) != 32:
            raise ValueError(f"POLICY_PUBKEYS: key {kid!r} is not a base64 ed25519 public key")
        keys[kid] = Ed25519PublicKey.from_public_bytes(raw)
    return keys


def utcnow() -> datetime:
    return datetime.now(timezone.utc)


def decode(data: bytes):
    try:
        return json.loads(data)
    except (json.JSONDecodeError, UnicodeDecodeError) as e:
        raise Refusal(SCHEMA_INVALID, f"envelope: {e}") from None
