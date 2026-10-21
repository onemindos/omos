#!/usr/bin/env python3
"""Validate local, redacted portal trust receipts for human review only.

No network, subprocesses, credential access, referenced-file reads or mutations.
Hashes bind supplied statements; they do not authenticate the observer or evidence.
"""

import argparse
from datetime import datetime, timedelta, timezone
import hashlib
import json
from pathlib import Path
import re
import sys


SCHEMA = "onemind.portal-acceptance.v1"
SCOPE = "portal-trust-recovery"
# Deliberately separate live coverage from regression/fixture coverage.
REQUIRED_CHECKS = (
    "release.exact-revision-replicas",
    "security.anonymous-denied",
    "security.machine-read-only",
    "security.backend-mutations-denied",
    "session.refresh-recovery-both-replicas",
    "session.expiry-fail-closed-both-replicas",
    "session.provider-failure-fail-closed",
    "roles.existing-human-isolation",
    "roles.actor-scope-isolation",
    "recovery.portal-isolated-restore",
    "recovery.grafana-durable-restore",
    "recovery.independent-backup-restore",
    "recovery.identity-key-availability",
    "recovery.reviewed-rollback-compatibility",
)
SHA = re.compile(r"[a-f0-9]{40}")
HASH = re.compile(r"sha256:[a-f0-9]{64}")
TIME = re.compile(
    r"[0-9]{4}-[0-9]{2}-[0-9]{2}T(?:[01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]"
    r"(?:\.[0-9]+)?(?:Z|[+-](?:[01][0-9]|2[0-3]):[0-5][0-9])"
)


class InvalidReceipt(ValueError):
    """Input is not a v1 receipt. Messages never echo receipt values."""


def require(condition, message):
    if not condition:
        raise InvalidReceipt(message)


def fields(value, expected, label):
    require(isinstance(value, dict) and set(value) == set(expected),
            f"{label}: required/unknown fields")


def timestamp(value):
    require(isinstance(value, str) and TIME.fullmatch(value) is not None,
            "timestamp must be RFC3339 with timezone")
    try:
        return datetime.fromisoformat(value.replace("Z", "+00:00")).astimezone(timezone.utc)
    except ValueError:
        raise InvalidReceipt("invalid calendar timestamp") from None


def unique_object(pairs):
    result = {}
    for key, value in pairs:
        require(key not in result, "duplicate JSON field")
        result[key] = value
    return result


def parse_receipt(raw):
    try:
        return json.loads(raw, object_pairs_hook=unique_object,
                          parse_constant=reject_constant)
    except (ValueError, RecursionError) as error:
        if isinstance(error, InvalidReceipt):
            raise
        raise InvalidReceipt("invalid JSON") from None


def reject_constant(_):
    raise InvalidReceipt("non-finite JSON number")


def validate(receipt, *, expected_commit, now=None, max_age_hours=24):
    """Return a deterministic review report; never certify an action or phase.

    Missing named checks are incomplete, malformed/unknown checks are invalid,
    any explicit failure dominates incomplete. No evidence is fetched or trusted.
    """
    require(isinstance(expected_commit, str) and SHA.fullmatch(expected_commit),
            "expected commit must be a full lowercase Git SHA")
    require(type(max_age_hours) is int and 1 <= max_age_hours <= 168,
            "max age must be 1..168 hours")
    now = now if now is not None else datetime.now(timezone.utc)
    require(isinstance(now, datetime) and now.tzinfo is not None
            and now.utcoffset() is not None, "review time must have a timezone")
    fields(receipt, ("schemaVersion", "commit", "observedAt", "scope", "checks"), "receipt")
    require(receipt["schemaVersion"] == SCHEMA, "unsupported schema")
    require(isinstance(receipt["commit"], str) and SHA.fullmatch(receipt["commit"]),
            "commit must be a full lowercase Git SHA")
    require(receipt["scope"] == SCOPE, "unsupported named scope")
    observed = timestamp(receipt["observedAt"])
    require(observed <= now, "receipt observation is in the future")
    require(isinstance(receipt["checks"], list), "checks must be an array")
    stale_before = now - timedelta(hours=max_age_hours)
    reasons = []
    if receipt["commit"] != expected_commit:
        reasons.append("commit-mismatch")
    if observed < stale_before:
        reasons.append("stale-receipt")
    seen, failed, incomplete = set(), [], []
    for check in receipt["checks"]:
        fields(check, ("name", "status", "observedAt", "evidence"), "check")
        name = check["name"]
        require(isinstance(name, str) and name in REQUIRED_CHECKS, "unknown check name")
        require(name not in seen, "duplicate named check")
        seen.add(name)
        require(check["status"] in ("passed", "failed", "not-tested"), "invalid check status")
        at = timestamp(check["observedAt"])
        require(at <= observed, "check observation is after receipt observation")
        evidence = check["evidence"]
        fields(evidence, ("kind", "ref", "sha256"), "evidence")
        require(evidence["kind"] in ("observation", "fixture"), "invalid evidence kind")
        require(isinstance(evidence["ref"], str) and 0 < len(evidence["ref"]) <= 1024
                and evidence["ref"] == evidence["ref"].strip()
                and not any(ord(c) < 32 for c in evidence["ref"]), "invalid evidence reference")
        require(isinstance(evidence["sha256"], str) and HASH.fullmatch(evidence["sha256"]),
                "evidence must have a SHA-256 content reference")
        if check["status"] == "failed":
            failed.append(name)
        if check["status"] == "not-tested" or evidence["kind"] == "fixture" or at < stale_before:
            incomplete.append(name)
    missing = sorted(set(REQUIRED_CHECKS) - seen)
    state = "failed" if failed else "incomplete" if reasons or missing or incomplete else "passed"
    canonical = json.dumps(receipt, sort_keys=True, separators=(",", ":"), ensure_ascii=True)
    return {
        "schemaVersion": "onemind.portal-acceptance-review.v1",
        "scope": SCOPE,
        "commit": receipt["commit"],
        "reviewedAt": now.astimezone(timezone.utc).isoformat(),
        "maxAgeHours": max_age_hours,
        "receiptSha256": "sha256:" + hashlib.sha256(canonical.encode("utf-8")).hexdigest(),
        "status": state,
        "reasons": reasons,
        "missingChecks": missing,
        "failedChecks": sorted(failed),
        "incompleteChecks": sorted(incomplete),
        "humanReviewRequired": True,
        "phase6Done": False,
        "migrationAuthorized": False,
        "actionAuthorized": False,
        "evidenceContentsVerified": False,
    }


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("receipt", type=Path, help="local redacted JSON receipt (read only)")
    parser.add_argument("--expected-commit", required=True, help="human-selected full Git SHA")
    parser.add_argument("--max-age-hours", type=int, default=24,
                        help="review freshness, 1..168 hours; default 24, not a health SLA")
    args = parser.parse_args(argv)
    try:
        receipt = parse_receipt(args.receipt.read_text(encoding="utf-8"))
        report = validate(receipt, expected_commit=args.expected_commit,
                          max_age_hours=args.max_age_hours)
    except (InvalidReceipt, OSError, UnicodeError, RecursionError):
        # No paths, evidence references, tokens or untrusted values echoed on error.
        print(json.dumps({"status": "invalid", "humanReviewRequired": True,
                          "phase6Done": False, "migrationAuthorized": False,
                          "actionAuthorized": False}), file=sys.stdout)
        return 2
    print(json.dumps(report, indent=2, sort_keys=True))
    return 0 if report["status"] == "passed" else 1


if __name__ == "__main__":
    sys.exit(main())
