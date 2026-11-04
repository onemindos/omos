"""Data-group gateway backends: trimmed OpenAPI specs, their ConfigMap mounts, and the ClickHouse read-only user."""
from pathlib import Path
import importlib.util
import json
import re
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / "platform/tool-gateway"
spec = importlib.util.spec_from_file_location("trim_openapi", ROOT / "scripts/gateway/trim_openapi_get.py")
trim = importlib.util.module_from_spec(spec)
spec.loader.exec_module(trim)


def read(path):
    return yaml.safe_load((ROOT / path).read_text())


class OpenApiTargets(unittest.TestCase):
    def test_configmaps_are_current(self):
        for path, text in trim.build().items():
            self.assertEqual(path.read_text(), text, f"{path.name} is stale: run scripts/gateway/trim_openapi_get.py")

    def test_every_exposed_operation_is_a_get_with_a_generic_response(self):
        for path in trim.build():
            doc = json.loads(yaml.safe_load(path.read_text())["data"]["openapi.yaml"])
            self.assertTrue(doc["openapi"].startswith("3.0"), path.name)  # the gateway parser rejects 3.1 constructs
            self.assertNotIn("components", doc)
            for p, methods in doc["paths"].items():
                self.assertEqual(set(methods), {"get"}, f"{path.name} {p}")

    def test_tools_match_classification(self):
        for name in ("stac", "core"):
            doc = json.loads(yaml.safe_load((BASE / f"manifests/openapi-{name}.configmap.yaml").read_text())["data"]["openapi.yaml"])
            tools = {m["get"]["operationId"] for m in doc["paths"].values()}
            self.assertEqual(tools, set(read(f"platform/tool-gateway/policy/classification/{name}.yaml")["read"]), name)

    def test_core_spec_only_names_get_routes_of_the_service(self):
        api = (ROOT / "services/onemind-core/onemind/api.py").read_text()
        routes = set(re.findall(r'@app\.get\("([^"]+)"', api))
        doc = json.loads(yaml.safe_load((BASE / "manifests/openapi-core.configmap.yaml").read_text())["data"]["openapi.yaml"])
        for p in doc["paths"]:
            self.assertIn(p, routes)

    def test_each_spec_is_mounted_where_the_target_reads_it(self):
        values = read("platform/tool-gateway/values.yaml")
        mounts = {m["mountPath"]: m["name"] for m in values["extraVolumeMounts"]}
        volumes = {v["name"]: v["configMap"]["name"] for v in values["extraVolumes"]}
        for name in ("stac", "core"):
            target = read(f"platform/tool-gateway/targets/{name}.yaml")["target"]["openapi"]
            self.assertEqual(Path(target["schema"]["file"]).parent.as_posix(), f"/specs/{name}")
            self.assertEqual(volumes[mounts[f"/specs/{name}"]], f"gateway-openapi-{name}")


class ClickHouse(unittest.TestCase):
    def test_legacy_ro_is_readonly_and_has_no_plaintext_password(self):
        chi = read("platform/data/clickhouse/clickhouse.yaml")["spec"]["configuration"]
        users, profiles = chi["users"], chi["profiles"]
        self.assertEqual(users["legacy_ro/profile"], "legacy_ro")
        self.assertEqual(profiles["legacy_ro/readonly"], 1)
        self.assertNotIn("legacy_ro/password", users)  # never a literal
        # A hash, not valueFrom.secretKeyRef: a secret ref adds a container env var and the operator restarts ClickHouse.
        self.assertRegex(users["legacy_ro/password_sha256_hex"], r"^[0-9a-f]{64}$")

    def test_backend_is_pinned_read_only_and_authenticated(self):
        dep = list(yaml.safe_load_all((BASE / "manifests/backends/clickhouse/deployment.yaml").read_text()))[0]
        c = dep["spec"]["template"]["spec"]["containers"][0]
        self.assertRegex(c["image"], r"@sha256:[0-9a-f]{64}$")
        env = {e["name"]: e for e in c["env"]}
        self.assertNotIn("CLICKHOUSE_ALLOW_WRITE_ACCESS", env)
        self.assertNotIn("CLICKHOUSE_MCP_AUTH_DISABLED", env)
        self.assertIn("valueFrom", env["CLICKHOUSE_MCP_AUTH_TOKEN"])
        self.assertTrue(c["securityContext"]["readOnlyRootFilesystem"])
        self.assertIn("limits", c["resources"])


if __name__ == "__main__":
    unittest.main()
