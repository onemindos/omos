"""Offline checks for the Phase 7a parallel Maputnik migration.

Run: python3 -m unittest discover -s tests -p test_openchoreo_maputnik.py -v
These tests never mutate or connect to Kubernetes.
"""

from pathlib import Path
import unittest

import yaml


ROOT = Path(__file__).resolve().parents[1]
APP = ROOT / "platform/openchoreo/projects/onetak/maputnik"
IMAGE = (
    "docker.io/maputnik/editor@sha256:"
    "35d3074330b01a55b14cd29c7a0168f09728afabc673b7d262eb70546eed7d75"
)


def load(path):
    return list(yaml.safe_load_all(path.read_text()))


class MaputnikMigrationTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.component, cls.workload = load(APP / "component.yaml")
        cls.release = load(APP / "release-v1.yaml")[0]
        cls.binding = load(APP / "development.yaml")[0]
        cls.redirect = load(APP / "canonical-path.yaml")[0]

    def test_same_project_and_explicit_promotion(self):
        self.assertEqual(self.component["spec"]["owner"]["projectName"], "onetak")
        self.assertIs(self.component["spec"]["autoDeploy"], False)
        self.assertEqual(self.binding["spec"]["owner"], {
            "projectName": "onetak", "componentName": "maputnik",
        })
        self.assertEqual(self.binding["spec"]["releaseName"], "maputnik-v1")
        self.assertEqual(self.binding["spec"]["environment"], "development")

    def test_existing_image_is_immutable(self):
        self.assertEqual(self.workload["spec"]["container"], {"image": IMAGE})
        self.assertEqual(self.release["spec"]["workload"]["container"]["image"], IMAGE)

    def test_generated_type_snapshot_matches_source(self):
        source = load(ROOT / "platform/openchoreo/platform/onemind-service.yaml")[0]
        self.assertEqual(self.release["spec"]["componentType"]["spec"], source["spec"])

    def test_generated_workload_snapshot_matches_source(self):
        expected = {k: v for k, v in self.workload["spec"].items() if k != "owner"}
        self.assertEqual(self.release["spec"]["workload"], expected)
        self.assertEqual(
            self.release["spec"]["componentProfile"]["traits"],
            self.component["spec"]["traits"],
        )

    def test_health_probes_are_functional_http_checks(self):
        probes = next(t for t in self.component["spec"]["traits"]
                      if t["name"] == "health-probes")["parameters"]
        for name in ("startup", "readiness", "liveness"):
            self.assertIs(probes[name]["enabled"], True)
            self.assertEqual(probes[name]["handler"], {
                "httpGet": {"path": "/", "port": 8888},
            })
        self.assertEqual({t["name"] for t in self.release["spec"]["traits"]}, {
            "cilium-gateway-ingress", "health-probes",
        })

    def test_resource_budget_does_not_expand_the_legacy_app(self):
        old = load(ROOT / "platform/onetak/maputnik.yaml")[0]
        expected = old["spec"]["template"]["spec"]["containers"][0]["resources"]
        self.assertEqual(
            self.binding["spec"]["componentTypeEnvironmentConfigs"]["resources"],
            expected,
        )

    def test_relative_assets_get_canonical_prefix(self):
        spec = self.redirect["spec"]
        self.assertEqual(spec["parentRefs"], [{
            "name": "gateway-default", "namespace": "openchoreo-data-plane",
        }])
        self.assertEqual(spec["hostnames"], [
            "development-default.apps.openchoreo.100-110-185-30.nip.io",
        ])
        rule = spec["rules"][0]
        self.assertEqual(rule["matches"], [{
            "path": {"type": "Exact", "value": "/maputnik-http"},
        }])
        self.assertEqual(rule["filters"], [{
            "type": "RequestRedirect",
            "requestRedirect": {
                "path": {"type": "ReplaceFullPath", "replaceFullPath": "/maputnik-http/"},
                "statusCode": 302,
            },
        }])
        self.assertNotIn("backendRefs", rule)

    def test_parallel_candidate_preserves_legacy_and_storage(self):
        legacy = load(ROOT / "platform/onetak/maputnik.yaml")
        self.assertEqual({d["kind"] for d in legacy}, {"Deployment", "Service"})
        self.assertEqual(legacy[0]["spec"]["replicas"], 1)
        kinds = {d["kind"] for f in APP.glob("*.yaml") for d in load(f) if d}
        self.assertEqual(kinds, {
            "Component", "Workload", "ComponentRelease", "ReleaseBinding", "HTTPRoute",
        })
        self.assertNotIn("env", self.workload["spec"]["container"])
        self.assertNotIn("files", self.workload["spec"]["container"])


if __name__ == "__main__":
    unittest.main()
