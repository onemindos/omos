"""Exercise history suppression with the exact deployed Kyverno 1.19.1 engine.

Set KYVERNO_CLI to a checksum-verified 1.19.1 CLI. No cluster access is used.
Native-number tests also exercise the exact pinned JMESPath dependency because
Kubernetes dynamic-client int64 values differ from CLI JSON float64 values.
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
# Admission enforces these; the other policies stay Audit (see the policy file header).
ENFORCED_POLICIES = {IMAGE_POLICY, "require-namespace-layer-label"}
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
    # Dynamic-client numbers must normalize without treating unknown state as history.
    for section, key in (("spec", "replicas"), ("status", "replicas"),
                         ("status", "observedGeneration"), ("metadata", "generation")):
        for suffix, value in (("missing", None), ("unknown", "unknown"),
                              ("boolean", False), ("object", {})):
            r = copy.deepcopy(retired_rs)
            r["metadata"]["name"] = f"rs-{section}-{key.lower()}-{suffix}"
            if suffix == "missing":
                r[section].pop(key, None)
            else:
                r[section][key] = value
            add(r)
    add(workload("Deployment", "deployment-scaled-zero",
                 status={"replicas": 0, "observedGeneration": 3}))
    add(workload("ReplicationController", "replicationcontroller-zero",
                 references=[owner("Deployment")], status={"replicas": 0}))

    completed_job = workload("Job", "job-complete", references=[owner("CronJob", "batch/v1")],
                             status={"active": 0, "conditions": [{"type": "Complete", "status": "True"}]})
    add(completed_job, True)
    for key in ("active", "terminating"):
        for suffix, value in (("unknown", "unknown"), ("boolean", False), ("object", {})):
            r = copy.deepcopy(completed_job)
            r["metadata"]["name"] = f"job-{key}-{suffix}"
            r["status"][key] = value
            add(r)
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

    def test_only_violation_free_policies_enforce(self):
        names = set()
        for policy in yaml.safe_load_all(POLICY.read_text()):
            name = policy["metadata"]["name"]
            names.add(name)
            with self.subTest(policy=name):
                expected = "Enforce" if name in ENFORCED_POLICIES else "Audit"
                self.assertEqual(policy["spec"]["validationFailureAction"], expected)
                for rule in policy["spec"]["rules"]:
                    # Kyverno 1.19.1 lets a rule-level failureAction override the spec field.
                    self.assertNotIn("failureAction", rule["validate"])
        self.assertLessEqual(ENFORCED_POLICIES, names)


def _wave(doc):
    return int(((doc.get("metadata") or {}).get("annotations") or {}).get("argocd.argoproj.io/sync-wave", "0"))


def _labelled_namespaces(paths):
    found = {}
    for path in paths:
        try:
            docs = list(yaml.safe_load_all(path.read_text()))
        except yaml.YAMLError:
            continue
        for doc in docs:
            if isinstance(doc, dict) and doc.get("kind") == "Namespace":
                labels = doc["metadata"].get("labels") or {}
                if labels.get("onemind.io/layer"):
                    found[doc["metadata"]["name"]] = doc
    return found


class NamespaceBootstrap(unittest.TestCase):
    """Enforce must never deny a namespace that Argo CD auto-creates on a rebuild from zero.

    Argo CD's CreateNamespace=true creates an unlabelled namespace unless the Application ships the
    Namespace itself or sets managedNamespaceMetadata. A namespace labelled by a *different*
    Application races require-namespace-layer-label, unless platform-security labels it: its
    Namespaces apply in the same sync as, and before, the Enforce ClusterPolicies.
    """

    SECURITY = ROOT / "platform/security"

    def test_security_namespaces_precede_policies(self):
        namespaces = _labelled_namespaces(sorted(self.SECURITY.glob("*.yaml")))
        policy_waves = {_wave(p) for p in yaml.safe_load_all(POLICY.read_text())}
        self.assertTrue(namespaces)
        for name, doc in namespaces.items():
            with self.subTest(namespace=name):
                self.assertLess(_wave(doc), min(policy_waves))

    def test_auto_created_namespaces_are_labelled_before_enforce(self):
        security = _labelled_namespaces(sorted(self.SECURITY.glob("*.yaml")))
        checked = 0
        for path in sorted((ROOT / "platform").rglob("*.yaml")):
            if "parked" in path.parts:
                continue
            try:
                docs = list(yaml.safe_load_all(path.read_text()))
            except yaml.YAMLError:
                continue
            for app in docs:
                if not (isinstance(app, dict) and app.get("kind") == "Application"
                        and str(app.get("apiVersion", "")).startswith("argoproj.io/")):
                    continue
                spec = app["spec"]
                policy = spec.get("syncPolicy") or {}
                namespace = spec.get("destination", {}).get("namespace")
                if "CreateNamespace=true" not in (policy.get("syncOptions") or []) or not namespace:
                    continue
                # argocd exists before Kyverno (Helm bootstrap) and must not use managedNamespaceMetadata.
                if namespace == "argocd":
                    continue
                checked += 1
                managed = ((policy.get("managedNamespaceMetadata") or {}).get("labels") or {})
                own = set()
                for source in spec.get("sources") or [spec.get("source")]:
                    if not source or source.get("chart") or not source.get("path"):
                        continue
                    base = ROOT / source["path"]
                    deep = (base / "kustomization.yaml").exists() or (source.get("directory") or {}).get("recurse")
                    own |= set(_labelled_namespaces(sorted(base.rglob("*.yaml") if deep else base.glob("*.yaml"))))
                with self.subTest(application=app["metadata"]["name"], namespace=namespace):
                    self.assertTrue(
                        namespace in own or namespace in security or managed.get("onemind.io/layer"),
                        "label it in its own source, in platform/security/shared-namespaces.yaml, "
                        "or with managedNamespaceMetadata",
                    )
        self.assertGreater(checked, 0)


class KyvernoNativeNumberContext(unittest.TestCase):
    """Reproduce the reports controller's native-int64 context without JSON laundering."""

    @classmethod
    def setUpClass(cls):
        go = shutil.which("go")
        if not go:
            raise RuntimeError("Go is required for the Kyverno native-number regression")
        cls.temp = tempfile.TemporaryDirectory(prefix="kyverno-native-context-")
        cls.addClassCleanup(cls.temp.cleanup)
        cls.binary = Path(cls.temp.name) / "native-context"
        subprocess.run([go, "build", "-mod=readonly", "-buildvcs=false", "-trimpath", "-o", str(cls.binary), "."],
                       cwd=ROOT / "scripts/tests/jmespath_native", check=True,
                       capture_output=True, text=True)

    def evaluate(self, queries, contexts):
        completed = subprocess.run([str(self.binary)], check=True, capture_output=True, text=True,
                                   input=json.dumps({"queries": queries, "contexts": contexts}))
        return json.loads(completed.stdout)

    def test_history_guard_has_identical_native_and_json_classification(self):
        queries = [entry["variable"]["jmesPath"]
                   for policy in yaml.safe_load_all(POLICY.read_text())
                   if policy["metadata"]["name"] in (IMAGE_POLICY, RESOURCE_POLICY)
                   for rule in policy["spec"]["rules"] for entry in rule.get("context", [])
                   if entry["name"] == "retiredWorkloadHistory"]
        self.assertEqual(len(queries), 3)
        cases = fixtures()
        results = self.evaluate(queries, [{"request": {"object": r}} for r, _ in cases])
        for mode, rows in results.items():
            for (resource, expected), values in zip(cases, rows):
                for index, result in enumerate(values):
                    # The policy variable's default:false keeps null/errors evaluated.
                    with self.subTest(mode=mode, resource=resource["metadata"]["name"], rule=index):
                        self.assertEqual(result.get("value") is True and not result.get("error"), expected)
        self.assertEqual(results["json-float64"], results["kubernetes-int64"])

    def test_original_numeric_comparisons_reproduce_scanner_failure(self):
        queries = ["request.object.spec.replicas == `0`",
                   "request.object.status.observedGeneration >= request.object.metadata.generation",
                   "to_number(to_string(request.object.spec.replicas)) == `0`",
                   "to_number(to_string(request.object.status.observedGeneration)) >= "
                   "to_number(to_string(request.object.metadata.generation))"]
        resource = workload("ReplicaSet", "native-type-proof", references=[owner("Deployment")],
                            status={"replicas": 0, "observedGeneration": 3})
        results = self.evaluate(queries, [{"request": {"object": resource}}])
        self.assertEqual([x["value"] for x in results["json-float64"][0]], [True] * 4)
        self.assertEqual([x["value"] for x in results["kubernetes-int64"][0]],
                         [False, None, True, True])


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
        cls.messages = {}
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
            cls.messages[mode] = {(reference["name"], result["policy"], result["rule"]): result.get("message", "")
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
                key = (name, policy, prefix + rule)
                with self.subTest(mode=mode, resource=name, rule=rule):
                    if mode == "admission-update" and policy == IMAGE_POLICY:
                        # The Enforce image rules match CREATE only (see the policy header): Kyverno
                        # 1.19.1 never applies allowExistingViolations to a rule with preconditions, so
                        # matching UPDATE would deny every update of an already-violating object.
                        self.assertIsNone(self.results[mode].get(key))
                    else:
                        self.assertEqual(self.results[mode].get(key), result)

    def test_background_skips_only_controller_history(self):
        self.assert_guarded_results("background", True)

    def test_empty_uid_is_not_an_admission_bypass(self):
        self.assert_guarded_results("empty-uid", True)

    def test_admission_create_evaluates_even_historical_status(self):
        self.assert_guarded_results("admission-create", False)

    def test_admission_update_skips_image_rules(self):
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
