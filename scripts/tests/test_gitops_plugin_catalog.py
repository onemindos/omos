"""The plugin-ecosystem catalog is generated, in sync in both copies, and wired into the cluster."""
from pathlib import Path
import subprocess
import sys
import unittest
import yaml

ROOT = Path(__file__).resolve().parents[2]
PLUGINS_TARGET = '/app/catalog-plugins/plugins.yaml'


class PluginCatalogContracts(unittest.TestCase):
    def test_generated_files_are_current_and_refs_resolve(self):
        r = subprocess.run([sys.executable, str(ROOT / 'scripts/generate_plugin_catalog.py'), '--check'],
                           capture_output=True, text=True)
        self.assertEqual(r.returncode, 0, r.stderr)

    def test_cluster_mounts_and_registers_the_file(self):
        kust = yaml.safe_load((ROOT / 'platform/backstage/kustomization.yaml').read_text())
        gen = {g['name']: g for g in kust['configMapGenerator']}
        self.assertEqual(gen['backstage-catalog-plugins']['files'], ['catalog/plugins.yaml'])
        dep = next(d for d in yaml.safe_load_all((ROOT / 'platform/backstage/deployment.yaml').read_text())
                   if d and d['kind'] == 'Deployment' and d['metadata']['name'] == 'backstage')
        pod = dep['spec']['template']['spec']
        self.assertIn({'name': 'catalog-plugins', 'configMap': {'name': 'backstage-catalog-plugins'}}, pod['volumes'])
        self.assertIn('/app/catalog-plugins', [m['mountPath'] for m in pod['containers'][0]['volumeMounts']])
        cfg = yaml.safe_load((ROOT / 'platform/backstage/app-config.cluster.yaml').read_text())
        targets = [loc['target'] for loc in cfg['catalog']['locations']]
        self.assertIn(PLUGINS_TARGET, targets)
        prod = yaml.safe_load((ROOT / 'services/onemind-portal/app-config.production.yaml').read_text())
        self.assertEqual([t for t in targets if t != PLUGINS_TARGET],
                         [loc['target'] for loc in prod['catalog']['locations']])


if __name__ == '__main__':
    unittest.main()
