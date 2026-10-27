"""Offline GitHub catalog discovery contracts. No GitHub, cluster or portal credentials."""
from pathlib import Path
import json
import unittest
import yaml

ROOT = Path(__file__).resolve().parents[2]
PORTAL = ROOT / "services/onemind-portal"
MODULE = "@backstage/plugin-catalog-backend-module-github"


def entity_ref(entity):
    meta = entity["metadata"]
    return f'{entity["kind"]}:{meta.get("namespace", "default")}/{meta["name"]}'.lower()


def parse_ref(value, default_kind):
    kind, sep, rest = value.partition(":")
    if not sep:
        kind, rest = default_kind, value
    namespace, sep, name = rest.partition("/")
    if not sep:
        namespace, name = "default", rest
    return f"{kind}:{namespace}/{name}".lower()


def static_refs():
    refs = set()
    for name in ("apps.yaml", "cluster.yaml", "org.yaml"):
        text = (ROOT / "platform/backstage/catalog" / name).read_text()
        refs |= {entity_ref(e) for e in yaml.safe_load_all(text) if e}
    return refs


def cluster_config():
    return yaml.safe_load((ROOT / "platform/backstage/app-config.cluster.yaml").read_text())


class GithubDiscoveryTests(unittest.TestCase):
    def test_module_is_pinned_locked_and_registered_once(self):
        # 0.14.0 is the version in the Backstage 1.55.3 release manifest (backstage.json).
        backend = json.loads((PORTAL / "packages/backend/package.json").read_text())
        self.assertEqual(backend["dependencies"][MODULE], "0.14.0")
        self.assertIn(f'"{MODULE}@npm:0.14.0":', (PORTAL / "yarn.lock").read_text())
        index = (PORTAL / "packages/backend/src/index.ts").read_text()
        self.assertEqual(index.count(f"backend.add(import('{MODULE}'));"), 1)

    def test_provider_registers_only_existing_root_descriptors_with_the_existing_app(self):
        cluster = cluster_config()
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

    def test_discovered_descriptors_never_redefine_static_entities(self):
        root = yaml.safe_load((ROOT / "catalog-info.yaml").read_text())
        self.assertEqual(entity_ref(root), "location:default/omos-service-descriptors")
        discovered = [entity_ref(root)]
        for target in root["spec"]["targets"]:
            self.assertRegex(target, r"^\./services/[a-z0-9-]+/catalog-info\.yaml$")
            discovered += [entity_ref(e) for e in yaml.safe_load_all((ROOT / target).read_text()) if e]
        self.assertGreater(len(discovered), 1)
        self.assertEqual(len(discovered), len(set(discovered)))
        # Backstage rejects a second definition of an entity ref; it never merges the two.
        self.assertEqual(sorted(set(discovered) & static_refs()), [])

    def test_service_descriptors_use_catalog_owners_systems_and_relations(self):
        static = static_refs()
        allowed = {kind for rule in cluster_config()["catalog"]["rules"] for kind in rule.get("allow", [])}
        descriptors = {
            path.relative_to(ROOT).as_posix(): [e for e in yaml.safe_load_all(path.read_text()) if e]
            for path in sorted(ROOT.glob("services/*/catalog-info.yaml"))
        }
        known = static | {entity_ref(e) for entities in descriptors.values() for e in entities}
        self.assertIn("services/onemind-portal/catalog-info.yaml", descriptors)
        for path, entities in descriptors.items():
            for entity in entities:
                with self.subTest(path=path, entity=entity_ref(entity)):
                    self.assertIn(entity["kind"], allowed)
                    spec = entity["spec"]
                    self.assertIn(parse_ref(spec["owner"], "group"), static)
                    self.assertIn(parse_ref(spec["system"], "system"), static)
                    for relation in spec.get("dependsOn", []) + spec.get("dependencyOf", []):
                        self.assertIn(parse_ref(relation, "component"), known)


if __name__ == "__main__":
    unittest.main()
