"""Rendered safety contracts for the private personal vault deployment."""
import ipaddress
import subprocess
import shutil
import unittest
from pathlib import Path
import yaml

ROOT = Path(__file__).resolve().parents[2]


class VaultwardenSafety(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        app = ROOT / 'platform/vaultwarden'
        k = yaml.safe_load((app / 'kustomization.yaml').read_text())
        cls.docs = [d for resource in k['resources']
                    for d in yaml.safe_load_all((app / resource).read_text())]
        generator, = k['configMapGenerator']
        cls.docs.append({'kind': 'ConfigMap', 'metadata': {'name': generator['name']},
                         'data': dict(v.split('=', 1) for v in generator['literals'])})

    def one(self, kind):
        matches = [d for d in self.docs if d['kind'] == kind]
        self.assertEqual(len(matches), 1)
        return matches[0]

    def test_render_keeps_single_writer_and_hashed_configuration_reference(self):
        d = self.one('Deployment')
        self.assertEqual(d['spec']['replicas'], 1)
        self.assertEqual(d['spec']['strategy']['type'], 'Recreate')
        pod = d['spec']['template']['spec']
        self.assertFalse(pod['automountServiceAccountToken'])
        self.assertTrue(pod['securityContext']['runAsNonRoot'])
        container, = pod['containers']
        self.assertRegex(container['image'], r'^vaultwarden/server:1\.37\.4@sha256:[0-9a-f]{64}$')
        cm = self.one('ConfigMap')
        self.assertIn({'configMapRef': {'name': cm['metadata']['name']}}, container['envFrom'])
        for key in ['SIGNUPS_ALLOWED', 'INVITATIONS_ALLOWED', 'SSO_ENABLED', 'SSO_SIGNUPS_ALLOWED']:
            self.assertEqual(cm['data'][key], 'false')
        self.assertEqual(cm['data']['SIGNUPS_DOMAINS_WHITELIST'], '')
        self.assertNotIn('ADMIN_TOKEN', cm['data'])
        self.assertTrue(container['securityContext']['readOnlyRootFilesystem'])
        self.assertIn('ALL', container['securityContext']['capabilities']['drop'])
        self.assertTrue(all(container[p]['httpGet']['path'] == '/alive'
                            for p in ['startupProbe', 'readinessProbe', 'livenessProbe']))
        for d in self.docs:
            self.assertNotIn(d['kind'], ['Role', 'RoleBinding', 'ClusterRole', 'ClusterRoleBinding'])

    @unittest.skipUnless(shutil.which('kubectl'), 'Kustomize render verified when kubectl is available')
    def test_kustomize_rewrites_configuration_reference_to_content_hash(self):
        docs = list(yaml.safe_load_all(subprocess.check_output(
            ['kubectl', 'kustomize', str(ROOT / 'platform/vaultwarden')], text=True)))
        cm, = [d for d in docs if d['kind'] == 'ConfigMap']
        deployment, = [d for d in docs if d['kind'] == 'Deployment']
        container, = deployment['spec']['template']['spec']['containers']
        self.assertIn({'configMapRef': {'name': cm['metadata']['name']}}, container['envFrom'])
        self.assertNotEqual(cm['metadata']['name'], 'vaultwarden-configuration')
        self.assertEqual(cm['data'], self.one('ConfigMap')['data'])

    def test_storage_survives_gitops_deletion_and_no_public_listener(self):
        for kind in ['Namespace', 'PersistentVolumeClaim']:
            annotations = self.one(kind)['metadata']['annotations']['argocd.argoproj.io/sync-options']
            self.assertIn('Prune=false', annotations)
            self.assertIn('Delete=false', annotations)
        self.assertEqual(self.one('Ingress')['spec']['ingressClassName'], 'tailscale')
        service = self.one('Service')
        self.assertEqual(service['spec'].get('type', 'ClusterIP'), 'ClusterIP')
        policies = [d for d in self.docs if d['kind'] == 'NetworkPolicy']
        web = next(d for d in policies if d['spec']['podSelector'].get('matchLabels', {}).get('app.kubernetes.io/component') == 'web')
        peers = web['spec']['ingress'][0]['from']
        self.assertEqual(len(peers), 1)
        self.assertEqual(peers[0]['namespaceSelector']['matchLabels']['kubernetes.io/metadata.name'], 'tailscale')
        self.assertEqual(peers[0]['podSelector']['matchLabels']['tailscale.com/parent-resource'], 'vaultwarden')
        blocks = [peer['ipBlock'] for rule in web['spec']['egress'] for peer in rule['to'] if 'ipBlock' in peer]
        for private in ['10.109.3.1', '169.254.169.254', '100.100.100.100', '192.168.1.1', '127.0.0.1']:
            self.assertFalse(any(ipaddress.ip_address(private) in ipaddress.ip_network(b['cidr'])
                                 and not any(ipaddress.ip_address(private) in ipaddress.ip_network(n)
                                             for n in b.get('except', [])) for b in blocks))

    def test_admin_ciphertext_is_name_bound_and_no_plaintext_secret_is_rendered(self):
        self.assertFalse(any(d['kind'] == 'Secret' for d in self.docs))
        admin = next(d for d in self.docs if d['kind'] == 'SealedSecret' and d['metadata']['name'] == 'vaultwarden-admin')
        self.assertEqual(admin['metadata']['namespace'], 'vaultwarden')
        self.assertEqual(set(admin['spec']['encryptedData']), {'ADMIN_TOKEN'})
        annotations = admin['metadata'].get('annotations', {})
        self.assertNotIn('sealedsecrets.bitnami.com/cluster-wide', annotations)
        self.assertNotIn('sealedsecrets.bitnami.com/namespace-wide', annotations)


if __name__ == '__main__':
    unittest.main()
