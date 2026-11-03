"""OpenAPI targets of the tool gateway: the mounted spec, the classification and the fragment must agree."""
from pathlib import Path
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / "platform/tool-gateway"
CLASSES = ("read", "write", "destructive", "physical")
METHODS = ("get", "put", "post", "delete", "patch")


def openapi_targets():
    for p in sorted((BASE / "targets").glob("*.yaml")):
        frag = yaml.safe_load(p.read_text())
        if "openapi" in frag["target"]:
            yield frag


class OpenApiTargets(unittest.TestCase):
    def test_there_is_at_least_one(self):
        self.assertTrue(list(openapi_targets()))

    def test_spec_matches_classification_and_fragment(self):
        values = yaml.safe_load((BASE / "values.yaml").read_text())
        volumes = {v["name"]: v["configMap"]["name"] for v in values["extraVolumes"]}
        mounts = {m["name"]: m for m in values["extraVolumeMounts"]}
        for frag in openapi_targets():
            name = frag["name"]
            with self.subTest(target=name):
                cm = yaml.safe_load((BASE / f"manifests/openapi-{name}.configmap.yaml").read_text())
                self.assertEqual(cm["metadata"]["name"], f"gateway-openapi-{name}")
                self.assertEqual(cm["metadata"]["namespace"], "tool-gateway")
                spec = yaml.safe_load(cm["data"]["openapi.yaml"])
                ops = {o["operationId"]: (m.upper(), p)
                       for p, item in spec["paths"].items() for m, o in item.items() if m in METHODS}
                cls = yaml.safe_load((BASE / f"policy/classification/{name}.yaml").read_text())
                listed = {n: c for c in CLASSES for n in cls.get(c) or []}
                self.assertEqual(set(ops), set(listed), "spec operations and classification must be identical")
                # HTTP method sanity: DELETE and kick-style POSTs are never read or write.
                for op, (method, path) in ops.items():
                    if method == "DELETE" or "/kick" in path:
                        self.assertEqual(listed[op], "destructive", op)
                    if method == "GET":
                        self.assertEqual(listed[op], "read", op)
                # The server in the spec is the host the fragment calls.
                self.assertEqual(spec["servers"][0]["url"], "http://" + frag["target"]["openapi"]["host"])
                # Mounted where the fragment reads it.
                path = frag["target"]["openapi"]["schema"]["file"]
                vol = next(v for v, c in volumes.items() if c == cm["metadata"]["name"])
                self.assertEqual(path, f"{mounts[vol]['mountPath']}/openapi.yaml")
                self.assertTrue(mounts[vol]["readOnly"])
                # Backing secret and egress exist.
                for ref in (frag.get("env") or {}).values():
                    sealed = yaml.safe_load((BASE / f"manifests/gateway-target-{name}.sealedsecret.yaml").read_text())
                    self.assertEqual(sealed["metadata"]["name"], ref["secret"])
                    self.assertIn(ref["key"], sealed["spec"]["encryptedData"])
                self.assertTrue((BASE / f"manifests/egress-{name}.yaml").exists())

    def test_mediamtx_never_exposes_config_endpoints(self):
        cm = yaml.safe_load((BASE / "manifests/openapi-mediamtx.configmap.yaml").read_text())
        spec = yaml.safe_load(cm["data"]["openapi.yaml"])
        for path in spec["paths"]:
            self.assertNotIn("/config/", path)
            self.assertNotIn("/auth/", path)


if __name__ == "__main__":
    unittest.main()
