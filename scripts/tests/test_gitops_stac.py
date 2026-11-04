"""STAC (platform/stac) contracts: the WebODM sync's unit tests, the embedded code, and the write-access boundary."""
from pathlib import Path
import importlib.util
import re
import subprocess
import sys
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
STAC = ROOT / "platform/stac"
SERVICE = ROOT / "services/stac-webodm-sync"


def docs(path):
    return [d for d in yaml.safe_load_all(path.read_text()) if d]


def all_stac_docs():
    out = []
    for f in sorted(STAC.glob("*.yaml")):
        out.extend((f.name, d) for d in docs(f))
    return out


def containers(doc):
    spec = doc["spec"]
    if doc["kind"] == "CronJob":
        spec = spec["jobTemplate"]["spec"]
    return spec["template"]["spec"]["containers"]


class SyncService(unittest.TestCase):
    def test_service_unit_tests_pass(self):
        run = subprocess.run([sys.executable, "-m", "unittest", "discover", "-s", str(SERVICE)],
                             cwd=ROOT, capture_output=True, text=True)
        self.assertEqual(run.returncode, 0, run.stdout + run.stderr)

    def test_configmap_holds_exactly_the_committed_sync_py(self):
        spec = importlib.util.spec_from_file_location("render_configmap", SERVICE / "render_configmap.py")
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        committed = (STAC / "webodm-sync-code.configmap.yaml").read_text()
        self.assertEqual(committed, module.render(), "run python services/stac-webodm-sync/render_configmap.py")
        self.assertEqual(yaml.safe_load(committed)["data"]["sync.py"], (SERVICE / "sync.py").read_text())


class WriteBoundary(unittest.TestCase):
    def test_only_the_writer_deployment_enables_transactions(self):
        enabled = {}
        for name, d in all_stac_docs():
            if d["kind"] != "Deployment":
                continue
            env = {e["name"]: e.get("value") for c in containers(d) for e in c.get("env", [])}
            enabled[d["metadata"]["name"]] = env.get("ENABLE_TRANSACTIONS_EXTENSIONS")
        self.assertEqual(enabled["stac-api"], "false")
        self.assertEqual(enabled["stac-api-writer"], "true")
        self.assertNotIn("true", [v for k, v in enabled.items() if k != "stac-api-writer"])

    def test_the_writer_has_no_ingress_and_only_the_sync_job_may_reach_it(self):
        for _name, d in all_stac_docs():
            if d["kind"] == "Ingress":
                self.assertNotEqual(d["spec"]["defaultBackend"]["service"]["name"], "stac-api-writer")
                self.assertEqual(d["spec"]["ingressClassName"], "tailscale")
        policy = next(d for _n, d in all_stac_docs()
                      if d["kind"] == "NetworkPolicy" and d["metadata"]["name"] == "api-writer")
        sources = [f for rule in policy["spec"]["ingress"] for f in rule["from"]]
        self.assertEqual(sources, [{"podSelector": {"matchLabels": {
            "app.kubernetes.io/name": "stac", "app.kubernetes.io/component": "webodm-sync"}}}])

    def test_hostnames_are_stac_and_stac_browser(self):
        hosts = {d["metadata"]["annotations"]["tailscale.com/hostname"] for _n, d in all_stac_docs()
                 if d["kind"] == "Ingress"}
        self.assertEqual(hosts, {"stac", "stac-browser"})


class Supply(unittest.TestCase):
    def test_every_image_is_pinned_by_tag_and_digest(self):
        images = [c["image"] for _n, d in all_stac_docs() if d["kind"] in ("Deployment", "Job", "CronJob")
                  for c in containers(d)]
        self.assertGreaterEqual(len(images), 5)
        for image in images:
            self.assertRegex(image, r":[A-Za-z0-9._-]+@sha256:[0-9a-f]{64}$", image)
            self.assertNotIn(":latest", image)

    def test_every_container_has_requests_and_a_memory_limit_and_drops_capabilities(self):
        for _n, d in all_stac_docs():
            if d["kind"] not in ("Deployment", "Job", "CronJob"):
                continue
            for c in containers(d):
                res = c["resources"]
                self.assertIn("cpu", res["requests"], c["name"])
                self.assertIn("memory", res["requests"], c["name"])
                self.assertIn("memory", res["limits"], c["name"])
                self.assertEqual(c["securityContext"]["capabilities"]["drop"], ["ALL"], c["name"])
                self.assertIs(c["securityContext"]["allowPrivilegeEscalation"], False, c["name"])

    def test_no_plaintext_secret_is_committed_and_example_files_are_not_synced(self):
        for name, d in all_stac_docs():
            self.assertNotEqual(d["kind"], "Secret", name)
        for example in STAC.glob("*.secret.example"):
            self.assertRegex(example.read_text(), r"<[^>]+>")  # placeholders only
        app = next(iter(docs(ROOT / "platform/argocd/platform-stac.yaml")))
        self.assertIn("*.secret*", app["spec"]["source"]["directory"]["exclude"])

    def test_namespace_is_labelled_for_its_layer(self):
        ns = next(d for _n, d in all_stac_docs() if d["kind"] == "Namespace")
        self.assertEqual(ns["metadata"]["labels"]["onemind.io/layer"], "make-sense")


class Database(unittest.TestCase):
    def test_stac_database_role_and_secret_exist_on_the_shared_cluster(self):
        res = ROOT / "platform/ecosystem/resources"
        database = next(d for d in docs(res / "databases.yaml") if d["spec"]["name"] == "stac")
        self.assertEqual(database["spec"]["owner"], "stac")
        self.assertEqual({e["name"] for e in database["spec"]["extensions"]}, {"postgis", "btree_gist", "unaccent"})
        cluster = next(d for d in docs(res / "postgres-cluster.yaml") if d["kind"] == "Cluster")
        role = next(r for r in cluster["spec"]["managed"]["roles"] if r["name"] == "stac")
        self.assertFalse(role["superuser"])
        self.assertEqual(role["passwordSecret"]["name"], "db-stac")
        self.assertTrue((res / "secrets/onemind-data--db-stac.json").exists())

    def test_migration_is_a_pinned_sync_hook_with_the_reviewed_compatibility_wrapper(self):
        job = next(d for n, d in all_stac_docs() if n == "migrate.yaml")
        ann = job["metadata"]["annotations"]
        self.assertEqual(ann["argocd.argoproj.io/hook"], "Sync")
        self.assertIn("BeforeHookCreation", ann["argocd.argoproj.io/hook-delete-policy"])
        c = job["spec"]["template"]["spec"]["containers"][0]
        self.assertIn("python3 /scripts/migrate-compat.py migrate", c["args"][0])
        self.assertIn({"name": "migration-compat", "mountPath": "/scripts", "readOnly": True}, c["volumeMounts"])
        tag = re.search(r"pgstac-pypgstac:(v[0-9.]+)@", c["image"]).group(1)
        self.assertRegex(tag, r"^v\d+\.\d+\.\d+$")


if __name__ == "__main__":
    unittest.main()
