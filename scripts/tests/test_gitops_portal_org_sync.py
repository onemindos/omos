"""Authentik org sync wiring contracts (docs/catalog/authentik-org-sync.md). Offline; no cluster or Authentik access."""
from pathlib import Path
import unittest
import yaml

ROOT = Path(__file__).resolve().parents[2]
PORTAL = ROOT / "services/onemind-portal"


def load(file):
    return [d for d in yaml.safe_load_all((ROOT / file).read_text()) if d]


class PortalOrgSyncTests(unittest.TestCase):
    def test_cluster_overlay_ships_the_sync_disabled(self):
        # Reviewed state: off until the owner seals backstage-authentik-catalog and binds the portal group
        # to the onemind-backstage application in Authentik. Enabling changes this assertion in the same PR.
        authentik = load("platform/backstage/app-config.cluster.yaml")[0]["authentik"]
        self.assertIs(authentik["orgSync"]["enabled"], False)
        self.assertEqual(authentik["orgSync"]["userGroups"], [])
        self.assertEqual(authentik["token"], "${AUTHENTIK_CATALOG_TOKEN}")
        self.assertEqual(authentik["baseUrl"], "http://platform-auth-helm-authentik-server.auth.svc.cluster.local")

    def test_portal_defaults_never_enable_it(self):
        for name in ("app-config.yaml", "app-config.production.yaml"):
            with self.subTest(file=name):
                config = load(f"services/onemind-portal/{name}")[0]
                self.assertIsNot(((config.get("authentik") or {}).get("orgSync") or {}).get("enabled"), True)

    def test_token_comes_only_from_the_optional_secret(self):
        # Optional, so a missing secret never blocks pod creation; the module then removes what it published.
        deployment = next(d for d in load("platform/backstage/deployment.yaml") if d["kind"] == "Deployment")
        containers = deployment["spec"]["template"]["spec"]["containers"]
        env = [e for c in containers for e in c.get("env", []) if e["name"] == "AUTHENTIK_CATALOG_TOKEN"]
        self.assertEqual(env, [{"name": "AUTHENTIK_CATALOG_TOKEN", "valueFrom": {"secretKeyRef": {
            "name": "backstage-authentik-catalog", "key": "AUTHENTIK_CATALOG_TOKEN", "optional": True}}}])

    def test_backend_registers_the_module_and_the_image_build_runs_its_tests(self):
        index = (PORTAL / "packages/backend/src/index.ts").read_text()
        self.assertIn("backend.add(import('@internal/plugin-onemind-org-backend'));", index)
        dockerfile = (PORTAL / "packages/backend/Dockerfile").read_text()
        self.assertRegex(dockerfile, r"(?m)^RUN yarn test --watch=false --runInBand .*\bplugins/plugin-onemind-org-backend/src\b")


if __name__ == "__main__":
    unittest.main()
