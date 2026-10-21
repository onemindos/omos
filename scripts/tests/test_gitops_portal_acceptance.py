"""Offline receipt contracts. All generated observations below are synthetic tests."""

import contextlib
import copy
from datetime import datetime, timedelta, timezone
import hashlib
import importlib.util
import io
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch


SOURCE = Path(__file__).resolve().parents[1] / "gitops/portal_acceptance.py"
SPEC = importlib.util.spec_from_file_location("portal_acceptance", SOURCE)
acceptance = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(acceptance)
NOW = datetime(2026, 10, 21, 1, tzinfo=timezone.utc)
COMMIT = "a" * 40


def receipt():
    return {
        "schemaVersion": acceptance.SCHEMA,
        "scope": acceptance.SCOPE,
        "commit": COMMIT,
        "observedAt": NOW.isoformat(),
        "checks": [{
            "name": name, "status": "passed", "observedAt": NOW.isoformat(),
            "evidence": {"kind": "observation", "ref": "redacted-receipt.json#" + name,
                         "sha256": "sha256:" + hashlib.sha256(name.encode()).hexdigest()},
        } for name in acceptance.REQUIRED_CHECKS],
    }


class PortalAcceptanceTests(unittest.TestCase):
    def review(self, value, **kwargs):
        return acceptance.validate(value, expected_commit=COMMIT, now=NOW, **kwargs)

    def test_complete_observation_pass_is_scoped_and_never_authorizes(self):
        result = self.review(receipt())
        self.assertEqual(result["status"], "passed")
        self.assertTrue(result["humanReviewRequired"])
        for field in ("phase6Done", "migrationAuthorized", "actionAuthorized", "evidenceContentsVerified"):
            self.assertFalse(result[field])

    def test_each_required_gate_individually_blocks_when_missing(self):
        for name in acceptance.REQUIRED_CHECKS:
            with self.subTest(name=name):
                value = receipt()
                value["checks"] = [c for c in value["checks"] if c["name"] != name]
                result = self.review(value)
                self.assertEqual(result["status"], "incomplete")
                self.assertEqual(result["missingChecks"], [name])

    def test_empty_receipt_is_incomplete_not_vacuous_pass(self):
        value = receipt()
        value["checks"] = []
        self.assertEqual(self.review(value)["status"], "incomplete")

    def test_not_tested_stays_incomplete(self):
        value = receipt()
        value["checks"][0]["status"] = "not-tested"
        self.assertEqual(self.review(value)["status"], "incomplete")

    def test_explicit_failure_dominates_missing_and_fixture(self):
        value = receipt()
        failed = value["checks"].pop()
        value["checks"][0]["status"] = "failed"
        value["checks"][0]["evidence"]["kind"] = "fixture"
        result = self.review(value)
        self.assertEqual(result["status"], "failed")
        self.assertIn(failed["name"], result["missingChecks"])
        self.assertEqual(result["failedChecks"], [value["checks"][0]["name"]])

    def test_all_fixture_and_mixed_fixture_passes_are_incomplete(self):
        for count in (1, len(acceptance.REQUIRED_CHECKS)):
            value = receipt()
            for check in value["checks"][:count]:
                check["evidence"]["kind"] = "fixture"
            result = self.review(value)
            self.assertEqual(result["status"], "incomplete")
            self.assertEqual(len(result["incompleteChecks"]), count)

    def test_stale_receipt_and_old_check_under_fresh_envelope(self):
        old = (NOW - timedelta(hours=25)).isoformat()
        value = receipt()
        value["checks"][0]["observedAt"] = old
        self.assertEqual(self.review(value)["status"], "incomplete")
        for check in value["checks"]:
            check["observedAt"] = old
        value["observedAt"] = old
        self.assertIn("stale-receipt", self.review(value)["reasons"])

    def test_freshness_boundary_and_bounded_override(self):
        value = receipt()
        for check in value["checks"]:
            check["observedAt"] = (NOW - timedelta(hours=24)).isoformat()
        self.assertEqual(self.review(value)["status"], "passed")
        self.assertEqual(self.review(value, max_age_hours=23)["status"], "incomplete")
        for hours in (0, -1, 169, True, float("inf")):
            with self.assertRaises(acceptance.InvalidReceipt):
                self.review(value, max_age_hours=hours)

    def test_future_envelope_and_future_check_rejected(self):
        for target in ("envelope", "check"):
            value = receipt()
            obj = value if target == "envelope" else value["checks"][0]
            obj["observedAt"] = (NOW + timedelta(seconds=1)).isoformat()
            with self.assertRaises(acceptance.InvalidReceipt):
                self.review(value)

    def test_timezone_offsets_are_compared_as_instants(self):
        value = receipt()
        value["observedAt"] = "2026-10-20T21:00:00-04:00"
        self.assertEqual(self.review(value)["status"], "passed")

    def test_invalid_times_and_naive_review_time(self):
        for bad in (None, "2026-10-21", "2026-10-21T01:00:00", "2026-02-30T01:00:00Z",
                    "2026-10-21T01:00:00+25:00", "2026-10-21T01:00:00+01:99",
                    "2026-10-21T01:00:00-00:99", "2026-10-21T01:00:60Z"):
            value = receipt()
            value["observedAt"] = bad
            with self.assertRaises(acceptance.InvalidReceipt):
                self.review(value)
        with self.assertRaises(acceptance.InvalidReceipt):
            acceptance.validate(receipt(), expected_commit=COMMIT, now=NOW.replace(tzinfo=None))

    def test_commit_mismatch_is_incomplete_and_short_sha_invalid(self):
        value = receipt()
        value["commit"] = "b" * 40
        self.assertIn("commit-mismatch", self.review(value)["reasons"])
        for bad in ("main", "abcd1234", "A" * 40, None):
            value["commit"] = bad
            with self.assertRaises(acceptance.InvalidReceipt):
                self.review(value)

    def test_strict_schema_and_no_blanket_done_input(self):
        variants = []
        for key in receipt():
            value = receipt()
            del value[key]
            variants.append(value)
        for key, val in (("phase6Done", True), ("schemaVersion", "v0"), ("scope", "phase6"),
                         ("checks", {}), ("tokens", "never-echo-this")):
            value = receipt()
            value[key] = val
            variants.append(value)
        for value in variants + [None, [], 1]:
            with self.subTest(value=value), self.assertRaises(acceptance.InvalidReceipt):
                self.review(value)

    def test_duplicate_unknown_and_malformed_checks(self):
        for change in (lambda v: v["checks"].append(copy.deepcopy(v["checks"][0])),
                       lambda v: v["checks"][0].update(name="all-security"),
                       lambda v: v["checks"][0].update(status=True),
                       lambda v: v["checks"][0].update(status="done"),
                       lambda v: v["checks"].append(None)):
            value = receipt()
            change(value)
            with self.assertRaises(acceptance.InvalidReceipt):
                self.review(value)

    def test_evidence_requires_type_reference_and_content_digest(self):
        for key, bad in (("kind", "live-ish"), ("ref", ""), ("ref", "line\nbreak"),
                         ("sha256", "sha256:abc"), ("sha256", None)):
            value = receipt()
            value["checks"][0]["evidence"][key] = bad
            with self.assertRaises(acceptance.InvalidReceipt):
                self.review(value)
        value = receipt()
        del value["checks"][0]["evidence"]["ref"]
        with self.assertRaises(acceptance.InvalidReceipt):
            self.review(value)

    def test_strict_json_rejects_duplicates_nonfinite_and_malformed(self):
        for raw in ('{"commit":1,"commit":2}', '{"x":NaN}', '{"x":Infinity}', '{', ''):
            with self.assertRaises(acceptance.InvalidReceipt):
                acceptance.parse_receipt(raw)

    def test_hash_is_order_stable_changes_with_evidence_and_input_not_modified(self):
        value = receipt()
        original = copy.deepcopy(value)
        first = self.review(value)["receiptSha256"]
        self.assertEqual(value, original)
        self.assertEqual(first, self.review(dict(reversed(list(value.items()))))["receiptSha256"])
        value["checks"][0]["evidence"]["sha256"] = "sha256:" + "f" * 64
        self.assertNotEqual(first, self.review(value)["receiptSha256"])

    def test_cli_read_only_local_input_exit_codes_and_no_reference_reads(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "receipt.json"
            for state, code in (("passed", 0), ("failed", 1), ("not-tested", 1)):
                value = receipt()
                value["checks"][0]["status"] = state
                path.write_text(json.dumps(value), encoding="utf-8")
                before = path.read_bytes()
                output = io.StringIO()
                real_validate = acceptance.validate
                with contextlib.redirect_stdout(output), patch.object(
                        acceptance, "validate", side_effect=lambda r, **kw: real_validate(r, now=NOW, **kw)):
                    self.assertEqual(acceptance.main([str(path), "--expected-commit", COMMIT]), code)
                self.assertFalse(json.loads(output.getvalue())["actionAuthorized"])
                self.assertEqual(path.read_bytes(), before)
                self.assertEqual(list(Path(directory).iterdir()), [path])
            path.write_text('{"private":"never-echo-this"}', encoding="utf-8")
            for input_path in (path, path.parent / "missing.json"):
                output = io.StringIO()
                with contextlib.redirect_stdout(output):
                    self.assertEqual(acceptance.main([str(input_path), "--expected-commit", COMMIT]), 2)
                self.assertEqual(json.loads(output.getvalue())["status"], "invalid")
                self.assertNotIn("never-echo-this", output.getvalue())


if __name__ == "__main__":
    unittest.main()
