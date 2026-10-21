"""Pure offline A0 contracts. A passing decision is NEVER permission to execute.

All identities, policy, observations and journal entries are supplied test data,
not authenticated or independently verified evidence. No I/O or executor exists.
"""

from dataclasses import asdict, dataclass, fields
import hashlib
import json
import re
from typing import get_type_hints


@dataclass(frozen=True)
class Target:
    repository: str
    branch: str
    manifest_path: str
    cluster: str
    environment: str
    namespace: str
    catalog_ref: str
    controller: str
    workload: str


@dataclass(frozen=True)
class Revision:
    base_commit: str
    candidate_commit: str
    diff_hash: str
    artifact_digest: str
    evidence_hash: str
    observed_commit: str
    observed_generation: int


@dataclass(frozen=True)
class Request:
    schema_version: int
    request_id: str
    idempotency_key: str
    action: str
    requester: str
    executor: str
    target: Target
    revision: Revision
    policy_version: str
    issued_at: int
    expires_at: int


@dataclass(frozen=True)
class Approval:
    schema_version: int
    approval_id: str
    approver: str
    request_hash: str
    issued_at: int
    expires_at: int


@dataclass(frozen=True)
class Policy:
    schema_version: int
    version: str
    target: Target
    layer: str
    reviewed: bool
    stateless: bool
    offline_test_only: bool
    requester: str
    approver: str
    executor: str
    max_ttl_seconds: int
    max_observation_age_seconds: int
    enabled: bool


@dataclass(frozen=True)
class Observation:
    target: Target
    revision: Revision
    policy_version: str
    observed_at: int


@dataclass(frozen=True)
class Journal:
    available: bool
    stopped: bool
    revoked_approval_ids: list[str]
    revoked_actor_ids: list[str]
    used_approval_ids: list[str]
    used_request_ids: list[str]
    used_idempotency_keys: list[str]


@dataclass(frozen=True)
class Decision:
    offline_contract_passed: bool
    reason: str
    request_hash: str = ""
    execution_authorized: bool = False


class ContractError(ValueError):
    """Malformed or unsupported offline contract."""


def _decode(cls, raw):
    """Reject unknown/missing fields and coercions, including bool-as-int."""
    if type(raw) is not dict or set(raw) != {f.name for f in fields(cls)}:
        raise ContractError("fields")
    values = {}
    for name, kind in get_type_hints(cls).items():
        value = raw[name]
        if kind in (Target, Revision):
            value = _decode(kind, value)
        elif kind == list[str]:
            if type(value) is not list or any(not _text(x) for x in value):
                raise ContractError("list")
            value = list(value)
        elif type(value) is not kind:
            raise ContractError("type")
        elif kind is str and not _text(value):
            raise ContractError("text")
        elif kind is int and value < 0:
            raise ContractError("integer")
        values[name] = value
    result = cls(**values)
    if hasattr(result, "schema_version") and result.schema_version != 1:
        raise ContractError("version")
    if cls is Target:
        # Literal relative paths only; never commands, glob patterns or URLs.
        path = result.manifest_path
        if (not re.fullmatch(r"[A-Za-z0-9_./-]+", path)
                or any(p in ("", ".", "..") for p in path.split("/"))):
            raise ContractError("path")
        if not re.fullmatch(r"[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+", result.repository):
            raise ContractError("repository")
    if cls is Revision:
        for key in ("base_commit", "candidate_commit", "observed_commit"):
            if not re.fullmatch(r"[0-9a-f]{40}", getattr(result, key)):
                raise ContractError("commit")
        for key in ("diff_hash", "artifact_digest", "evidence_hash"):
            if not re.fullmatch(r"sha256:[0-9a-f]{64}", getattr(result, key)):
                raise ContractError("digest")
        if result.observed_generation < 1:
            raise ContractError("generation")
    if cls is Approval and not re.fullmatch(r"sha256:[0-9a-f]{64}", result.request_hash):
        raise ContractError("request hash")
    return result


def _text(value):
    return (type(value) is str and 0 < len(value) <= 1024
            and value == value.strip() and value.isascii()
            and all(32 <= ord(c) < 127 for c in value))


def request_hash(raw: dict) -> str:
    """Hash the COMPLETE validated v1 request; not a signature or identity proof.

    Canonical form is sorted ASCII JSON, compact separators, no floats/nulls.
    Invalid inputs raise ContractError. Approvals must hash this exact form.
    """
    request = _decode(Request, raw)
    canonical = json.dumps(asdict(request), sort_keys=True, separators=(",", ":"),
                           ensure_ascii=True, allow_nan=False).encode("ascii")
    return "sha256:" + hashlib.sha256(canonical).hexdigest()


def evaluate(request, approval, policy, observation, journal, *, now) -> Decision:
    """Evaluate a synthetic snapshot without mutating or consuming any state.

    Replay rejection depends on the supplied journal. This pure function cannot
    atomically reserve an approval or prevent concurrent reuse. A runtime would
    need authenticated inputs, fresh reads and a durable transactional journal.
    """
    try:
        if type(now) is not int or now < 0:
            raise ContractError("clock")
        r = _decode(Request, request)
        a = _decode(Approval, approval)
        p = _decode(Policy, policy)
        o = _decode(Observation, observation)
        j = _decode(Journal, journal)
        digest = request_hash(request)
    except (ContractError, TypeError, ValueError, RecursionError):
        return Decision(False, "malformed_contract")

    def deny(reason):
        return Decision(False, reason, digest)

    if not j.available or j.stopped or not p.enabled:
        return deny("unavailable_or_stopped")
    if (p.layer != "application" or not p.reviewed or not p.stateless
            or not p.offline_test_only or p.target.environment != "offline-test"):
        return deny("not_reviewed_nonfoundation_test_target")
    if r.action != "propose_release":
        return deny("unsupported_action")
    if r.target != p.target or o.target != p.target:
        return deny("target_mismatch")
    if r.policy_version != p.version or o.policy_version != p.version:
        return deny("policy_drift")
    if ((r.requester, a.approver, r.executor) != (p.requester, p.approver, p.executor)
            or len({r.requester, a.approver, r.executor}) != 3):
        return deny("actor_separation")
    if {r.requester, a.approver, r.executor}.intersection(j.revoked_actor_ids):
        return deny("revoked_actor")
    if a.approval_id in j.revoked_approval_ids:
        return deny("revoked_approval")
    if (a.approval_id in j.used_approval_ids or r.request_id in j.used_request_ids
            or r.idempotency_key in j.used_idempotency_keys):
        return deny("replay")
    if (p.max_ttl_seconds <= 0 or p.max_observation_age_seconds <= 0
            or not r.issued_at <= a.issued_at <= now < a.expires_at <= r.expires_at
            or r.expires_at - r.issued_at > p.max_ttl_seconds
            or not 0 <= now - o.observed_at <= p.max_observation_age_seconds):
        return deny("expired_or_stale")
    if a.request_hash != digest:
        return deny("approval_hash_mismatch")
    if o.revision != r.revision:
        return deny("revision_drift")
    return Decision(True, "offline_contract_only", digest)
