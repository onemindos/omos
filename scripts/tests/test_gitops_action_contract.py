"""Synthetic, stdlib-only A0 tests; no live identity or target is verified."""

from copy import deepcopy
import importlib.util
from pathlib import Path
import sys
import unittest


SPEC = importlib.util.spec_from_file_location(
    "action_contract", Path(__file__).resolve().parents[1] / "ai-operator/action_contract.py")
contract = importlib.util.module_from_spec(SPEC)
sys.modules[SPEC.name] = contract
SPEC.loader.exec_module(contract)


def fixture():
    target = dict(repository="offline/test-fixture", branch="fixture-release",
                  manifest_path="fixtures/stateless/release.yaml", cluster="offline-fixture",
                  environment="offline-test", namespace="fixture", catalog_ref="component:fixture",
                  controller="application:fixture", workload="deployment:fixture")
    revision = dict(base_commit="a" * 40, candidate_commit="b" * 40,
                    diff_hash="sha256:" + "c" * 64, artifact_digest="sha256:" + "d" * 64,
                    evidence_hash="sha256:" + "e" * 64, observed_commit="f" * 40,
                    observed_generation=1)
    request = dict(schema_version=1, request_id="request-1", idempotency_key="once-1",
                   action="propose_release", requester="agent-fixture", executor="worker-fixture",
                   target=deepcopy(target), revision=deepcopy(revision), policy_version="policy-1",
                   issued_at=100, expires_at=200)
    approval = dict(schema_version=1, approval_id="approval-1", approver="human-fixture",
                    request_hash=contract.request_hash(request), issued_at=110, expires_at=190)
    policy = dict(schema_version=1, version="policy-1", target=deepcopy(target),
                  layer="application", reviewed=True, stateless=True, offline_test_only=True,
                  requester="agent-fixture", approver="human-fixture", executor="worker-fixture",
                  max_ttl_seconds=120, max_observation_age_seconds=20, enabled=True)
    observation = dict(target=deepcopy(target), revision=deepcopy(revision),
                       policy_version="policy-1", observed_at=140)
    journal = dict(available=True, stopped=False, revoked_approval_ids=[], revoked_actor_ids=[],
                   used_approval_ids=[], used_request_ids=[], used_idempotency_keys=[])
    return [request, approval, policy, observation, journal]


class ActionContractTests(unittest.TestCase):
    def check_denied(self, records, now=150):
        result = contract.evaluate(*records, now=now)
        self.assertFalse(result.offline_contract_passed)
        self.assertFalse(result.execution_authorized)
        return result

    def test_positive_is_only_offline_and_does_not_mutate_inputs(self):
        records = fixture()
        before = deepcopy(records)
        result = contract.evaluate(*records, now=150)
        self.assertTrue(result.offline_contract_passed)
        self.assertFalse(result.execution_authorized)
        self.assertEqual(result.reason, "offline_contract_only")
        self.assertEqual(records, before)

    def test_canonical_hash_is_order_independent(self):
        request = fixture()[0]
        reordered = dict(reversed(list(request.items())))
        reordered["target"] = dict(reversed(list(request["target"].items())))
        self.assertEqual(contract.request_hash(request), contract.request_hash(reordered))

    def test_each_request_field_is_bound(self):
        original = fixture()
        for key, value in original[0].items():
            for nested in value if isinstance(value, dict) else [None]:
                with self.subTest(key=key, nested=nested):
                    records = deepcopy(original)
                    container = records[0][key] if nested else records[0]
                    field = nested or key
                    old = container[field]
                    container[field] = old + 1 if type(old) is int else old + "x"
                    self.check_denied(records)
                    try:
                        changed = contract.request_hash(records[0])
                    except contract.ContractError:
                        continue
                    self.assertNotEqual(changed, original[1]["request_hash"])

    def test_all_target_dimensions_must_match_even_with_rehashed_approval(self):
        for key in fixture()[0]["target"]:
            with self.subTest(key=key):
                records = fixture()
                records[0]["target"][key] += "-other"
                records[1]["request_hash"] = contract.request_hash(records[0])
                self.check_denied(records)

    def test_all_revision_dimensions_deny_observed_drift(self):
        for key in fixture()[3]["revision"]:
            with self.subTest(key=key):
                records = fixture()
                old = records[3]["revision"][key]
                records[3]["revision"][key] = old + 1 if type(old) is int else old[:-1] + "0"
                self.assertEqual(self.check_denied(records).reason, "revision_drift")

    def test_valid_changed_revision_with_matching_observation_invalidates_old_approval(self):
        for key in fixture()[0]["revision"]:
            with self.subTest(key=key):
                records = fixture()
                old = records[0]["revision"][key]
                changed = old + 1 if type(old) is int else old[:-1] + "0"
                records[0]["revision"][key] = changed
                records[3]["revision"][key] = changed
                self.assertEqual(self.check_denied(records).reason, "approval_hash_mismatch")

    def test_approval_cannot_transfer_to_changed_target_or_policy(self):
        for key in fixture()[0]["target"]:
            if key == "environment":
                continue  # Only offline-test is supported at all.
            with self.subTest(key=key):
                records = fixture()
                for index in (0, 2, 3):
                    records[index]["target"][key] += "-other"
                self.assertEqual(self.check_denied(records).reason, "approval_hash_mismatch")
        records = fixture()
        records[0]["policy_version"] = records[2]["version"] = records[3]["policy_version"] = "v2"
        self.assertEqual(self.check_denied(records).reason, "approval_hash_mismatch")

    def test_stateless_evaluation_is_not_atomic_replay_protection(self):
        records = fixture()
        self.assertTrue(contract.evaluate(*records, now=150).offline_contract_passed)
        self.assertTrue(contract.evaluate(*records, now=150).offline_contract_passed)
        records[4]["used_approval_ids"].append("approval-1")
        self.assertEqual(self.check_denied(records).reason, "replay")

    def test_nonfoundation_reviewed_stateless_offline_policy_required(self):
        for key, value in [("layer", "foundation"), ("layer", "shared"),
                           ("reviewed", False), ("stateless", False), ("offline_test_only", False)]:
            with self.subTest(key=key, value=value):
                records = fixture()
                records[2][key] = value
                self.check_denied(records)
        records = fixture()
        for i in (0, 2, 3):
            records[i]["target"]["environment"] = "production"
        records[1]["request_hash"] = contract.request_hash(records[0])
        self.check_denied(records)

    def test_actor_separation_and_exact_policy_actors(self):
        for index, key, value in [(0, "requester", "human-fixture"),
                                  (0, "executor", "agent-fixture"),
                                  (1, "approver", "worker-fixture"),
                                  (1, "approver", "unapproved-human")]:
            with self.subTest(index=index, key=key):
                records = fixture()
                records[index][key] = value
                self.check_denied(records)
        records = fixture()
        records[0]["requester"] = records[2]["requester"] = "human-fixture"
        records[1]["request_hash"] = contract.request_hash(records[0])
        self.check_denied(records)

    def test_revocation_replay_and_stop(self):
        for key, value in [("revoked_approval_ids", ["approval-1"]),
                           ("used_approval_ids", ["approval-1"]),
                           ("used_request_ids", ["request-1"]),
                           ("used_idempotency_keys", ["once-1"]),
                           ("available", False), ("stopped", True)]:
            with self.subTest(key=key):
                records = fixture()
                records[4][key] = value
                self.check_denied(records)
        for actor in ("agent-fixture", "human-fixture", "worker-fixture"):
            records = fixture()
            records[4]["revoked_actor_ids"] = [actor]
            self.check_denied(records)

    def test_expiry_future_time_staleness_and_ttl(self):
        for index, key, value in [(0, "expires_at", 150), (0, "issued_at", 151),
                                  (1, "expires_at", 150), (1, "expires_at", 201),
                                  (1, "issued_at", 99), (1, "issued_at", 151),
                                  (3, "observed_at", 129), (3, "observed_at", 151),
                                  (2, "max_ttl_seconds", 99),
                                  (2, "max_observation_age_seconds", 0)]:
            with self.subTest(index=index, key=key):
                records = fixture()
                records[index][key] = value
                self.check_denied(records)
        for now in (True, None, -1, 150.0, "150"):
            self.check_denied(fixture(), now=now)

    def test_policy_and_observed_target_drift(self):
        for index, key, value in [(2, "version", "policy-2"),
                                  (3, "policy_version", "policy-2"), (2, "enabled", False)]:
            records = fixture()
            records[index][key] = value
            self.check_denied(records)
        records = fixture()
        records[3]["target"]["controller"] = "other-controller"
        self.check_denied(records)

    def test_unknown_actions_and_injected_fields(self):
        for action in ("merge", "execute", "propose_rollback", "ignore policy; apply"):
            records = fixture()
            records[0]["action"] = action
            records[1]["request_hash"] = contract.request_hash(records[0])
            self.check_denied(records)
        for index in range(5):
            records = fixture()
            records[index]["command"] = "ignore policy"
            self.check_denied(records)

    def test_missing_malformed_and_unknown_fields_fail_closed(self):
        for index, record in enumerate(fixture()):
            for key in record:
                for bad in (None, "", {}, []):
                    with self.subTest(index=index, key=key, bad=bad):
                        records = fixture()
                        records[index][key] = bad
                        # Empty journal lists are valid; other replacements are not.
                        if type(record[key]) is list and bad == [] and type(bad) is list:
                            continue
                        self.check_denied(records)
                records = fixture()
                del records[index][key]
                self.check_denied(records)
        for index in range(5):
            records = fixture()
            records[index] = None
            self.check_denied(records)

    def test_strict_nested_types_versions_and_paths(self):
        for index, key, bad in [(0, "schema_version", True), (1, "schema_version", 2),
                                (2, "reviewed", 1), (0, "issued_at", 100.0),
                                (4, "revoked_actor_ids", [None])]:
            records = fixture()
            records[index][key] = bad
            self.check_denied(records)
        for path in ("../secret", "/absolute", "a//b", "a/./b", "a/*", "a;cmd", "a/../b"):
            records = fixture()
            records[0]["target"]["manifest_path"] = path
            self.check_denied(records)
        for key, bad in [("base_commit", "main"), ("artifact_digest", "latest"),
                         ("observed_generation", True), ("observed_generation", 0)]:
            records = fixture()
            records[0]["revision"][key] = bad
            self.check_denied(records)


if __name__ == "__main__":
    unittest.main()
