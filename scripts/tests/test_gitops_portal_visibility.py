"""Offline GitOps/catalog wiring contracts. No browser or cluster credentials."""
from pathlib import Path
import unittest
import yaml

ROOT = Path(__file__).resolve().parents[2]


def load(file):
    return list(yaml.safe_load_all((ROOT / file).read_text()))


class PortalVisibilityTests(unittest.TestCase):
    def test_paired_catalogs(self):
        self.assertEqual(
            (ROOT / "services/onemind-portal/catalog/cluster.yaml").read_bytes(),
            (ROOT / "platform/backstage/catalog/cluster.yaml").read_bytes())

    def test_native_workload_and_application_identity(self):
        definitions = [
            ("k8s-backstage-backstage-21b99354", "platform/backstage/deployment.yaml", "Deployment", "backstage", "platform-backstage"),
            ("k8s-backstage-backstage-postgres-17d2f351", "platform/backstage/postgres.yaml", "StatefulSet", "backstage-postgres", "platform-backstage"),
            ("k8s-nats-nats-0dee54ac", "platform/nats/nats-statefulset.yaml", "StatefulSet", "nats", "platform-nats"),
        ]
        entities = load("platform/backstage/catalog/cluster.yaml")
        for identity, file, kind, name, application in definitions:
            with self.subTest(identity=identity):
                entity = next(e for e in entities if e["kind"] == "Component" and e["metadata"]["name"] == identity)
                annotation = entity["metadata"]["annotations"]
                resource = next(r for r in load(file) if r["kind"] == kind and r["metadata"]["name"] == name)
                self.assertEqual(resource["metadata"]["namespace"], annotation["backstage.io/kubernetes-namespace"])
                selector = dict(pair.split("=", 1) for pair in annotation["backstage.io/kubernetes-label-selector"].split(","))
                for key, value in selector.items():
                    self.assertEqual(resource["metadata"]["labels"][key], value)
                    self.assertEqual(resource["spec"]["selector"]["matchLabels"][key], value)
                    self.assertEqual(resource["spec"]["template"]["metadata"]["labels"][key], value)
                app = load(f"platform/argocd/{application}.yaml")[0]
                self.assertEqual(annotation["argocd/app-name"], application)
                self.assertEqual(annotation["argocd/instance-name"], "onemind-cloud")
                self.assertEqual(app["spec"]["destination"]["namespace"], resource["metadata"]["namespace"])
                self.assertEqual(app["spec"]["source"]["path"], str(Path(file).parent))
                self.assertEqual(app["spec"]["source"]["targetRevision"], "main")


if __name__ == "__main__":
    unittest.main()
