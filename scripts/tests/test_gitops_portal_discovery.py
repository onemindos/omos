"""Offline GitHub catalog discovery contracts. No GitHub, cluster or portal credentials."""
from pathlib import Path
import json
import unittest
import yaml

ROOT = Path(__file__).resolve().parents[2]
PORTAL = ROOT / "services/onemind-portal"
MODULE = "@backstage/plugin-catalog-backend-module-github"


class GithubDiscoveryTests(unittest.TestCase):
    def test_module_is_pinned_locked_and_registered_once(self):
        # 0.14.0 is the version in the Backstage 1.55.3 release manifest (backstage.json).
        backend = json.loads((PORTAL / "packages/backend/package.json").read_text())
        self.assertEqual(backend["dependencies"][MODULE], "0.14.0")
        self.assertIn(f'"{MODULE}@npm:0.14.0":', (PORTAL / "yarn.lock").read_text())
        index = (PORTAL / "packages/backend/src/index.ts").read_text()
        self.assertEqual(index.count(f"backend.add(import('{MODULE}'));"), 1)

    def test_provider_registers_only_existing_root_descriptors_with_the_existing_app(self):
        cluster = yaml.safe_load((ROOT / "platform/backstage/app-config.cluster.yaml").read_text())
        providers = cluster["catalog"]["providers"]["github"]
        self.assertEqual(list(providers), ["onemindos"])
        provider = providers["onemindos"]
        self.assertEqual(provider["organization"], "onemindos")
        # Literal path: a wildcard would also ingest the Scaffolder skeleton descriptors.
        self.assertEqual(provider["catalogPath"], "/catalog-info.yaml")
        # omos is the only repository the App installation selects; @backstage/integration refuses
        # to read any other onemindos repository with it, so a wider filter only adds failing Locations.
        self.assertEqual(provider["filters"], {"branch": "main", "repository": "^omos$", "allowForks": False})
        self.assertIs(provider["validateLocationsExist"], True)
        # No credential or host of its own: it authenticates through integrations.github.
        self.assertNotIn("host", provider)
        self.assertNotIn("${", json.dumps(provider))
        app = cluster["integrations"]["github"][0]["apps"][0]
        self.assertEqual(app["allowedInstallationOwners"], ["onemindos"])
        allowed = {kind for rule in cluster["catalog"]["rules"] for kind in rule.get("allow", [])}
        self.assertFalse(allowed & {"User", "Group", "Template"})


if __name__ == "__main__":
    unittest.main()
