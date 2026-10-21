"""Catalog-only topology contracts against the source-owned storage definitions."""
from pathlib import Path
import unittest
import yaml

ROOT = Path(__file__).resolve().parents[2]


class NativeTopologyContracts(unittest.TestCase):
    def test_existing_storage_identities_match_source_templates(self):
        catalog = list(yaml.safe_load_all((ROOT / 'platform/backstage/catalog/cluster.yaml').read_text()))
        claims = [e for e in catalog if e['kind'] == 'Resource' and e['spec']['type'] == 'kubernetes-persistent-volume-claim']
        self.assertEqual({(e['metadata']['annotations']['onemind.io/kubernetes-namespace'], e['metadata']['annotations']['onemind.io/kubernetes-pvc']) for e in claims}, {
            ('backstage', 'data-backstage-postgres-0'), ('monitoring', 'grafana-data'),
            ('nats', 'nats-js-nats-0'), ('nats', 'nats-js-nats-1'), ('nats', 'nats-js-nats-2')})
        pg = next(e for e in yaml.safe_load_all((ROOT / 'platform/backstage/postgres.yaml').read_text()) if e['kind'] == 'StatefulSet')
        self.assertEqual(pg['metadata']['name'], 'backstage-postgres')
        self.assertEqual(pg['spec']['volumeClaimTemplates'][0]['metadata']['name'], 'data')
        self.assertEqual(pg['spec']['replicas'], 1)
        nats = yaml.safe_load((ROOT / 'platform/nats/nats-statefulset.yaml').read_text())
        self.assertEqual(nats['spec']['replicas'], 3)
        self.assertEqual(nats['spec']['volumeClaimTemplates'][0]['metadata']['name'], 'nats-js')
        grafana = yaml.safe_load((ROOT / 'platform/values/kube-prometheus-stack.yaml').read_text())['grafana']
        self.assertEqual(grafana['persistence']['existingClaim'], 'grafana-data')

    def test_catalog_references_resolve_and_namespace_groups_are_not_clusters(self):
        docs = []
        for name in ['cluster', 'apps', 'org']:
            docs.extend(yaml.safe_load_all((ROOT / f'platform/backstage/catalog/{name}.yaml').read_text()))
        refs = {f"{e['kind'].lower()}:{e['metadata'].get('namespace', 'default')}/{e['metadata']['name']}" for e in docs}
        for e in docs:
            if e['kind'] == 'Resource' and e['spec']['type'] == 'kubernetes-persistent-volume-claim':
                for target in e['spec']['dependsOn']:
                    self.assertIn(target, refs)
                self.assertIn('system:default/' + e['spec']['system'], refs)
            if e['kind'] == 'System' and e['metadata']['name'].startswith('cluster-'):
                self.assertTrue(e['metadata']['title'].startswith('Namespace / '))
        self.assertEqual((ROOT / 'platform/backstage/catalog/cluster.yaml').read_bytes(), (ROOT / 'services/onemind-portal/catalog/cluster.yaml').read_bytes())
