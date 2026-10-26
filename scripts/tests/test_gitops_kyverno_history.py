"""Exercise history suppression with the exact deployed Kyverno 1.19.1 engine.

Set KYVERNO_CLI to a checksum-verified 1.19.1 CLI. No cluster access is used.
The CLI's default CREATE context has no request.uid, like the deployed scanner;
--set request.uid supplies the discriminator present in real AdmissionReviews.
Sources for that distinction and auto-generation behavior:
https://github.com/kyverno/kyverno/blob/v1.19.1/pkg/controllers/report/utils/scanner.go#L330-L343
https://github.com/kyverno/kyverno/blob/v1.19.1/pkg/engine/policycontext/policy_context.go#L185-L256
https://kyverno.io/docs/policy-types/cluster-policy/autogen/#exclusion-by-metadata
"""

import copy
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile
import unittest

import yaml


ROOT = Path(__file__).resolve().parents[2]
POLICY = ROOT / "platform/security/kyverno-policies-audit.yaml"
CLI = os.environ.get("KYVERNO_CLI") or shutil.which("kyverno")
IMAGE_POLICY = "disallow-latest-image-tag"
RESOURCE_POLICY = "require-requests-limits"
GUARDED_RULES = ("require-image-tag", "validate-image-tag", "validate-resources")
BAD_CONTAINER = {"name": "workload", "image": "busybox:latest", "resources": {}}
PROVIDER_BUDGET_SKIP = {
    "provider-system-pod", "provider-system-deployment", "provider-system-daemonset",
    "provider-system-active-replicaset", "provider-system-pod-untagged",
}


def owner(kind, api_version="apps/v1", **changes):
    reference = {"apiVersion": api_version, "kind": kind, "name": "parent",
                 "uid": "e7c18871-3a51-4448-a29a-5cd10c3f92fd", "controller": True}
    reference.update(changes)
    return reference


def workload(kind, name, *, status=None, references=None):
    metadata = {"name": name, "namespace": "history-tests", "generation": 3}
    if references is not None:
        metadata["ownerReferences"] = references
    pod = {"containers": [copy.deepcopy(BAD_CONTAINER)]}
    api_version = "v1" if kind in ("Pod", "ReplicationController") else (
        "batch/v1" if kind in ("Job", "CronJob") else "apps/v1")
    if kind == "Pod":
        spec = pod
    elif kind == "CronJob":
        spec = {"schedule": "0 0 * * *", "suspend": True,
                "jobTemplate": {"spec": {"template": {"spec": pod}}}}
    else:
        spec = {"template": {"metadata": {"labels": {"app": name}}, "spec": pod}}
        if kind not in ("Job",):
            spec.update({"replicas": 0, "selector": {"matchLabels": {"app": name}}})
    resource = {"apiVersion": api_version, "kind": kind, "metadata": metadata, "spec": spec}
    if status is not None:
        resource["status"] = status
    return resource


def fixtures():
    """Each expected skip is explicit; evaluation expectations do not inspect policy logic."""
    cases = []

    def add(resource, skip=False):
        cases.append((resource, skip))
        return resource

    retired_rs = workload("ReplicaSet", "rs-retired", references=[owner("Deployment")],
                          status={"replicas": 0, "observedGeneration": 3})
    add(retired_rs, True)
    for name, path, value in [
        ("rs-active", ("spec", "replicas"), 1),
        ("rs-draining", ("status", "replicas"), 1),
        ("rs-stale-status", ("status", "observedGeneration"), 2),
    ]:
        r = copy.deepcopy(retired_rs)
        r["metadata"]["name"] = name
        r[path[0]][path[1]] = value
        add(r)
    for name, section, key in [
        ("rs-default-replicas", "spec", "replicas"),
        ("rs-no-observation", "status", "observedGeneration"),
        ("rs-no-count", "status", "replicas"),
    ]:
        r = copy.deepcopy(retired_rs)
        r["metadata"]["name"] = name
        del r[section][key]
        add(r)
    add(workload("ReplicaSet", "rs-new-no-status", references=[owner("Deployment")]))
    for name, references in [
        ("rs-unowned", None), ("rs-orphan", []),
        ("rs-not-controller", [owner("Deployment", controller=False)]),
        ("rs-custom-owner", [owner("Deployment", "example.test/v1")]),
        ("rs-owner-no-uid", [owner("Deployment", uid="")]),
    ]:
        add(workload("ReplicaSet", name, references=references,
                     status={"replicas": 0, "observedGeneration": 3}))
    add(workload("Deployment", "deployment-scaled-zero",
                 status={"replicas": 0, "observedGeneration": 3}))
    add(workload("ReplicationController", "replicationcontroller-zero",
                 references=[owner("Deployment")], status={"replicas": 0}))

    completed_job = workload("Job", "job-complete", references=[owner("CronJob", "batch/v1")],
                             status={"active": 0, "conditions": [{"type": "Complete", "status": "True"}]})
    add(completed_job, True)
    failed = copy.deepcopy(completed_job)
    failed["metadata"]["name"] = "job-failed-terminal"
    failed["status"]["conditions"][0]["type"] = "Failed"
    add(failed, True)
    for name, changes in [
        ("job-still-active", {"active": 1}),
        ("job-still-terminating", {"terminating": 1}),
        ("job-failure-target", {"conditions": [{"type": "FailureTarget", "status": "True"}]}),
        ("job-not-complete", {"conditions": [{"type": "Complete", "status": "False"}]}),
        ("job-success-count-only", {"conditions": [], "succeeded": 1}),
        ("job-new-inactive", {"conditions": []}),
    ]:
        r = copy.deepcopy(completed_job)
        r["metadata"]["name"] = name
        r["status"].update(changes)
        add(r)
    add(workload("Job", "job-no-status", references=[owner("CronJob", "batch/v1")]))
    for name, references in [
        ("job-unowned-complete", None), ("job-orphan-complete", []),
        ("job-owner-no-uid", [owner("CronJob", "batch/v1", uid="")]),
        ("job-custom-owner", [owner("CronJob", "example.test/v1")]),
    ]:
        add(workload("Job", name, references=references, status=copy.deepcopy(completed_job["status"])))
    add(workload("CronJob", "cronjob-suspended"))

    for phase in ("Succeeded", "Failed", "Running", "Pending", "Unknown"):
        add(workload("Pod", "pod-" + phase.lower(), references=[owner("Job", "batch/v1")],
                     status={"phase": phase}), phase in ("Succeeded", "Failed"))
    add(workload("Pod", "pod-no-status", references=[owner("Job", "batch/v1")]))
    for name, references in [
        ("pod-unowned-terminal", None), ("pod-orphan-terminal", []),
        ("pod-not-controller", [owner("Job", "batch/v1", controller=False)]),
        ("pod-owner-no-uid", [owner("Job", "batch/v1", uid="")]),
        ("pod-custom-controller", [owner("Job", "example.test/v1")]),
    ]:
        add(workload("Pod", name, references=references, status={"phase": "Succeeded"}))
    for kind in ("ReplicaSet", "StatefulSet", "DaemonSet"):
        add(workload("Pod", "pod-terminal-" + kind.lower(), references=[owner(kind)],
                     status={"phase": "Succeeded"}), True)

    untagged = workload("Pod", "pod-untagged", status={"phase": "Running"})
    untagged["spec"]["containers"][0]["image"] = "busybox"
    add(untagged)
    good = workload("Pod", "pod-compliant", status={"phase": "Running"})
    good["spec"]["containers"][0].update({
        "image": "busybox:1.37.0",
        "resources": {"requests": {"cpu": "10m", "memory": "16Mi"}, "limits": {"memory": "32Mi"}},
    })
    add(good)
    security = workload("Pod", "pod-terminal-security-violations", references=[owner("Job", "batch/v1")],
                        status={"phase": "Succeeded"})
    security["spec"]["containers"][0]["securityContext"] = {"privileged": True}
    security["spec"]["volumes"] = [{"name": "host", "hostPath": {"path": "/"}}]
    add(security, True)

    # Provider ownership suppresses only background budgets, never image checks.
    for kind, name in [
        ("Pod", "provider-system-pod"),
        ("Deployment", "provider-system-deployment"),
        ("DaemonSet", "provider-system-daemonset"),
        ("ReplicaSet", "provider-system-active-replicaset"),
        ("Pod", "provider-system-pod-untagged"),
    ]:
        r = workload(kind, name, status={"phase": "Running"} if kind == "Pod" else {"replicas": 1})
        r["metadata"].update({"namespace": "kube-system", "labels": {"doks.digitalocean.com/managed": "true"}})
        if kind in ("Deployment", "ReplicaSet"):
            r["spec"]["replicas"] = 1
        if name == "provider-system-pod-untagged":
            r["spec"]["containers"][0]["image"] = "busybox"
        add(r)
    for name, namespace, labels in [
        ("provider-label-app-namespace", "history-tests", {"doks.digitalocean.com/managed": "true"}),
        ("provider-unlabelled-system-pod", "kube-system", {}),
        ("provider-false-system-pod", "kube-system", {"doks.digitalocean.com/managed": "false"}),
        ("provider-case-sensitive-system-pod", "kube-system", {"doks.digitalocean.com/managed": "True"}),
        ("provider-wrong-label-system-pod", "kube-system", {"doks.digitalocean.com/managed-by": "true"}),
    ]:
        r = workload("Pod", name, status={"phase": "Running"})
        r["metadata"].update({"namespace": namespace, "labels": labels})
        add(r)
    for name, top_labels in [
        ("provider-template-only-system-deployment", {}),
        ("provider-top-false-template-true-deployment", {"doks.digitalocean.com/managed": "false"}),
    ]:
        r = workload("Deployment", name, status={"replicas": 1})
        r["spec"]["replicas"] = 1
        r["metadata"].update({"namespace": "kube-system", "labels": top_labels})
        r["spec"]["template"]["metadata"]["labels"]["doks.digitalocean.com/managed"] = "true"
        add(r)
    return cases


class PolicyDelivery(unittest.TestCase):
    def test_server_defaults_are_explicit_without_rule_list_prepatching(self):
        app = yaml.safe_load((ROOT / "platform/argocd/platform-security.yaml").read_text())
        self.assertIn("ServerSideApply=true", app["spec"]["syncPolicy"]["syncOptions"])
        self.assertNotIn("RespectIgnoreDifferences=true", app["spec"]["syncPolicy"]["syncOptions"])
        self.assertFalse(app["spec"].get("ignoreDifferences"))
        for policy in yaml.safe_load_all(POLICY.read_text()):
            with self.subTest(policy=policy["metadata"]["name"]):
                spec = policy["spec"]
                self.assertTrue(spec["admission"])
                self.assertFalse(spec["emitWarning"])
                for rule in spec["rules"]:
                    self.assertTrue(rule["skipBackgroundRequests"])
                    self.assertTrue(rule["validate"]["allowExistingViolations"])


@unittest.skipUnless(CLI, "Set KYVERNO_CLI to the verified Kyverno 1.19.1 binary")
class KyvernoHistoryEngine(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        version = subprocess.check_output([CLI, "version"], text=True)
        if "Version: 1.19.1\n" not in version:
            raise RuntimeError("History guards must be tested with deployed Kyverno 1.19.1: " + version)
        cls.cases = fixtures()
        cls.temp = tempfile.TemporaryDirectory(prefix="kyverno-history-")
        cls.addClassCleanup(cls.temp.cleanup)
        resources = Path(cls.temp.name) / "resources.yaml"
        resources.write_text(yaml.safe_dump_all([r for r, _ in cls.cases], sort_keys=False))
        cls.results = {}
        for mode, variables in [
            ("background", []),
            ("empty-uid", ["--set", "request.uid="]),
            ("admission-create", ["--set", "request.uid=admission-proof,request.operation=CREATE"]),
            ("admission-update", ["--set", "request.uid=admission-proof,request.operation=UPDATE"]),
        ]:
            completed = subprocess.run(
                [CLI, "apply", str(POLICY), "--resource", str(resources), "--policy-report",
                 "--output-format", "json", "--remove-color", *variables],
                capture_output=True, text=True, check=False)
            if completed.returncode not in (0, 1):
                raise RuntimeError(completed.stdout + completed.stderr)
            try:
                # This CLI writes its deprecation warnings to stdout before JSON.
                payload, = [line for line in completed.stdout.splitlines() if line.startswith("{")]
                report = json.loads(payload)
            except (json.JSONDecodeError, ValueError) as error:
                raise AssertionError((completed.stdout + completed.stderr)[:4000]) from error
            if report["summary"]["error"]:
                raise AssertionError(completed.stdout + completed.stderr)
            cls.results[mode] = {(reference["name"], result["policy"], result["rule"]): result["result"]
                                 for result in report["results"] for reference in result["resources"]}

    def assert_guarded_results(self, mode, background):
        for resource, skip in self.cases:
            name = resource["metadata"]["name"]
            kind = resource["kind"]
            prefix = "" if kind == "Pod" else ("autogen-cronjob-" if kind == "CronJob" else "autogen-")
            expected = ["pass", "fail", "fail"]
            if name in ("pod-untagged", "provider-system-pod-untagged"):
                expected = ["fail", "pass", "fail"]
            elif name == "pod-compliant":
                expected = ["pass", "pass", "pass"]
            if background and skip:
                expected = ["skip"] * 3
            elif background and name in PROVIDER_BUDGET_SKIP:
                expected[2] = "skip"
            for rule, result in zip(GUARDED_RULES, expected):
                policy = RESOURCE_POLICY if rule == "validate-resources" else IMAGE_POLICY
                with self.subTest(mode=mode, resource=name, rule=rule):
                    self.assertEqual(self.results[mode].get((name, policy, prefix + rule)), result)

    def test_background_skips_only_controller_history(self):
        self.assert_guarded_results("background", True)

    def test_empty_uid_is_not_an_admission_bypass(self):
        self.assert_guarded_results("empty-uid", True)

    def test_admission_create_evaluates_even_historical_status(self):
        self.assert_guarded_results("admission-create", False)

    def test_admission_update_evaluates_even_historical_status(self):
        self.assert_guarded_results("admission-update", False)

    def test_other_security_rules_still_assess_terminal_pods(self):
        for mode, results in self.results.items():
            for rule in ("disallow-privileged", "disallow-hostpath"):
                with self.subTest(mode=mode, rule=rule):
                    self.assertEqual(results.get(("pod-terminal-security-violations",
                                                  "disallow-privileged-and-hostpath", rule)), "fail")

    def test_provider_boundary_never_skips_image_rules(self):
        for mode in ("background", "empty-uid"):
            for resource in PROVIDER_BUDGET_SKIP:
                prefix = "" if "pod" in resource else "autogen-"
                with self.subTest(mode=mode, resource=resource):
                    self.assertEqual(self.results[mode].get((resource, RESOURCE_POLICY,
                                                            prefix + "validate-resources")), "skip")
                    expected = ["fail", "pass"] if resource.endswith("untagged") else ["pass", "fail"]
                    for rule, result in zip(GUARDED_RULES[:2], expected):
                        self.assertEqual(self.results[mode].get((resource, IMAGE_POLICY, prefix + rule)), result)


if __name__ == "__main__":
    unittest.main()
