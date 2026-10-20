"""Offline lifecycle safety contracts: no cluster, registry or paid inference."""

from copy import deepcopy
import importlib.util
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import Mock, patch

ROOT = Path(__file__).resolve().parents[2]
SPEC = importlib.util.spec_from_file_location("gitops_pipeline", ROOT / "scripts/gitops/pipeline.py")
pipeline = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(pipeline)
P = json.loads((ROOT / "scripts/gitops/profiles.json").read_text())["portal"]
COMMIT = "a" * 40
DIGEST = "sha256:" + "b" * 64
IMAGE = P["imageRepository"] + "@" + DIGEST


def release():
    return {"commit": COMMIT, "digest": DIGEST, "published": True, "platform": "linux/amd64",
            "image": P["imageRepository"] + ":sha-" + COMMIT, "gates": sorted(pipeline.RELEASE_GATES)}


def fixtures():
    app = {"spec": {"source": {"repoURL": P["repository"] + ".git", "path": P["gitopsPath"], "targetRevision": "main"},
                     "destination": {"namespace": P["namespace"], "server": "https://kubernetes.default.svc"}},
           "status": {"sync": {"status": "Synced", "revision": COMMIT}, "health": {"status": "Healthy"},
                      "operationState": {"phase": "Succeeded"}, "resources": [{"kind": "Deployment", "name": "backstage"}]}}
    deploy = {"metadata": {"generation": 3}, "spec": {"replicas": 2, "template": {"spec": {"containers": [
        {"name": "backstage", "image": IMAGE}]}}}, "status": {"observedGeneration": 3, "replicas": 2,
        "readyReplicas": 2, "updatedReplicas": 2, "availableReplicas": 2}}
    pod = {"metadata": {}, "spec": {"containers": [{"name": "backstage", "image": IMAGE}]}, "status": {
        "conditions": [{"type": "Ready", "status": "True"}],
        "containerStatuses": [{"name": "backstage", "ready": True, "imageID": IMAGE}]}}
    return app, deploy, {"items": [pod, deepcopy(pod)]}


class ArtifactTests(unittest.TestCase):
    def test_exact_published_registry_commit_platform_and_gates(self):
        self.assertEqual(pipeline.release_contract(release(), P), IMAGE)

    def test_invalid_or_stale_release_metadata_rejected(self):
        for key, value in [("commit", "main"), ("digest", "latest"), ("published", False),
                           ("platform", "linux/arm64"), ("image", "unapproved/image:latest"), ("gates", [])]:
            with self.subTest(key=key), self.assertRaises(pipeline.GateError):
                record = release()
                record[key] = value
                pipeline.release_contract(record, P)

    def test_registry_labels_and_source_equivalence_are_mandatory(self):
        with tempfile.TemporaryDirectory() as tmp:
            folder = Path(tmp)
            record = folder / "release.json"
            record.write_text(json.dumps(release()))
            config = folder / "config.json"
            raw = folder / "manifest.json"
            config.write_text(json.dumps({"os": "linux", "architecture": "amd64", "config": {"Labels": {
                "org.opencontainers.image.revision": COMMIT, "org.opencontainers.image.source": P["repository"]}}}))
            raw.write_text(json.dumps({"schemaVersion": 2}))
            run = Mock()
            run.command.side_effect = [config, raw]
            run.git.side_effect = [COMMIT, ""]
            self.assertEqual(pipeline.artifact(run, record, P)[1], IMAGE)
            run.command.side_effect = [config, raw]
            run.git.side_effect = [COMMIT, "services/onemind-portal/changed.ts"]
            with self.assertRaisesRegex(pipeline.GateError, "Source changed"):
                pipeline.artifact(run, record, P)
            config.write_text(json.dumps({"os": "linux", "architecture": "amd64", "config": {"Labels": {}}}))
            run.command.side_effect = [config]
            with self.assertRaisesRegex(pipeline.GateError, "revision mismatch"):
                pipeline.artifact(run, record, P)


class RolloutTests(unittest.TestCase):
    def test_matching_healthy_synced_revision_and_all_replicas(self):
        self.assertTrue(pipeline.rollout_state(*fixtures(), P, IMAGE, COMMIT))

    def test_wrong_revision_not_accepted(self):
        self.assertFalse(pipeline.rollout_state(*fixtures(), P, IMAGE, "c" * 40))

    def test_wrong_owner_or_target_refused(self):
        for source_change in ({"path": "other/path"}, {"targetRevision": "feature"}, {"repoURL": "https://untrusted.invalid/repo"}):
            a, d, pods = fixtures()
            a["spec"]["source"].update(source_change)
            with self.subTest(source=source_change), self.assertRaises(pipeline.GateError):
                pipeline.rollout_state(a, d, pods, P, IMAGE, COMMIT)

    def test_ready_old_generation_or_partial_rollout_cannot_pass(self):
        for key in ("readyReplicas", "availableReplicas", "updatedReplicas", "replicas", "observedGeneration"):
            a, d, pods = fixtures()
            d["status"][key] = 1
            with self.subTest(key=key):
                self.assertFalse(pipeline.rollout_state(a, d, pods, P, IMAGE, COMMIT))

    def test_zero_or_changed_desired_replicas_cannot_pass(self):
        for n in (0, 1, 3):
            a, d, pods = fixtures()
            d["spec"]["replicas"] = n
            self.assertFalse(pipeline.rollout_state(a, d, pods, P, IMAGE, COMMIT))

    def test_wrong_spec_image_or_actual_runtime_digest_cannot_pass(self):
        a, d, pods = fixtures()
        pods["items"][0]["status"]["containerStatuses"][0]["imageID"] = "sha256:" + "c" * 64
        self.assertFalse(pipeline.rollout_state(a, d, pods, P, IMAGE, COMMIT))
        a, d, pods = fixtures()
        d["spec"]["template"]["spec"]["containers"][0]["image"] = "old/image"
        self.assertFalse(pipeline.rollout_state(a, d, pods, P, IMAGE, COMMIT))

    def test_index_and_verified_platform_manifest_digest_supported(self):
        a, d, pods = fixtures()
        platform_digest = "sha256:" + "c" * 64
        pods["items"][0]["status"]["containerStatuses"][0]["imageID"] = "repo@" + platform_digest
        self.assertTrue(pipeline.rollout_state(a, d, pods, P, IMAGE, COMMIT, {DIGEST, platform_digest}))

    def test_running_or_failed_argo_operation_cannot_pass(self):
        for phase in ("Running", "Failed", "Error"):
            a, d, pods = fixtures()
            a["status"]["operationState"]["phase"] = phase
            self.assertFalse(pipeline.rollout_state(a, d, pods, P, IMAGE, COMMIT))

    def test_deleted_unready_or_missing_pods_cannot_pass(self):
        a, d, pods = fixtures()
        pods["items"][0]["metadata"]["deletionTimestamp"] = "observed"
        self.assertFalse(pipeline.rollout_state(a, d, pods, P, IMAGE, COMMIT))
        a, d, pods = fixtures()
        pods["items"][0]["status"]["conditions"][0]["status"] = "False"
        self.assertFalse(pipeline.rollout_state(a, d, pods, P, IMAGE, COMMIT))


class IntentSafetyTests(unittest.TestCase):
    def test_pin_changes_only_the_unique_reviewed_container_image(self):
        text = (ROOT / P["manifest"]).read_text()
        result, previous = pipeline.changed_image(text, P, IMAGE)
        self.assertEqual(result.replace(IMAGE, previous), text)
        self.assertNotIn(previous, result)
        with self.assertRaises(pipeline.GateError):
            pipeline.changed_image(result, P, IMAGE)

    def test_ambiguous_manifest_is_not_rewritten(self):
        text = (ROOT / P["manifest"]).read_text()
        with self.assertRaises(pipeline.GateError):
            pipeline.changed_image(text + "\n---\n" + text, P, IMAGE)

    def test_profile_does_not_silently_onboard_an_unknown_app(self):
        with self.assertRaises(pipeline.GateError):
            pipeline.profile(ROOT, "unknown-application")

    def test_dirty_checkout_is_not_overwritten(self):
        run = object.__new__(pipeline.Run)
        run.git = Mock(return_value=" M unrelated-user-file")
        with self.assertRaises(pipeline.GateError):
            run.clean()

    def test_retirement_never_executes_a_deletion_or_claims_backup_certified(self):
        app, _, _ = fixtures()
        with patch.object(pipeline, "kube_get", side_effect=[app, {"items": [{"metadata": {"name": "retained-data"}}]}]) as get:
            result = pipeline.retirement_plan(Mock(), P)
        self.assertFalse(result["executionPermitted"])
        self.assertEqual(result["namespacePVCsToRetain"], ["retained-data"])
        self.assertEqual(get.call_count, 2)
        self.assertTrue(P["protected"])

    def test_empty_retirement_inventory_is_unknown_not_permission_to_prune(self):
        app, _, _ = fixtures()
        app["status"]["resources"] = []
        with patch.object(pipeline, "kube_get", return_value=app), self.assertRaises(pipeline.GateError):
            pipeline.retirement_plan(Mock(), P)

    def test_logs_are_private_and_required_failures_stop(self):
        with tempfile.TemporaryDirectory() as tmp:
            run = pipeline.Run(ROOT, Path(tmp) / "run")
            with self.assertRaises(pipeline.GateError):
                run.command("intentional-failure", ["false"])
            log = next(run.output.glob("*.log"))
            self.assertEqual(log.stat().st_mode & 0o777, 0o600)

    def test_techdocs_no_longer_uses_floating_npx_resolution(self):
        package = json.loads((ROOT / "services/onemind-portal/package.json").read_text())
        self.assertEqual(package["devDependencies"]["@techdocs/cli"], "1.12.0")
        script = (ROOT / "services/onemind-portal/scripts/build-techdocs.sh").read_text()
        self.assertNotIn("npx", script)
        self.assertIn("yarn-4.13.0.cjs exec techdocs-cli", script)
        self.assertIn("techdocs_metadata.json", script)
        self.assertIn("search_index.json", script)


if __name__ == "__main__":
    unittest.main()
