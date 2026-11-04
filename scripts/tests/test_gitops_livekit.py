"""Offline contracts for the staged LiveKit SFU and private credential boundary."""
import importlib.util
from pathlib import Path
import shutil
import subprocess
import tempfile
import unittest
import yaml

ROOT = Path(__file__).resolve().parents[2]
APP = ROOT / 'platform/livekit'


class LiveKitPilot(unittest.TestCase):
    def test_api_private_bind_and_direct_rtc_are_separate(self):
        deployment = yaml.safe_load((APP / 'deployment.yaml').read_text())
        pod = deployment['spec']['template']['spec']
        self.assertTrue(pod['hostNetwork'])
        self.assertFalse(pod['automountServiceAccountToken'])
        self.assertEqual(deployment['spec']['replicas'], 1)
        self.assertEqual(deployment['spec']['strategy']['type'], 'Recreate')
        container = pod['containers'][0]
        self.assertEqual(container['args'][-2:], ['--bind', '$(NODE_INTERNAL_IP)'])
        self.assertEqual(container['env'][0]['valueFrom']['fieldRef']['fieldPath'], 'status.hostIP')
        self.assertIn('@sha256:', container['image'])
        self.assertTrue(container['securityContext']['readOnlyRootFilesystem'])
        self.assertFalse(container['securityContext']['allowPrivilegeEscalation'])
        self.assertEqual(container['securityContext']['capabilities']['drop'], ['ALL'])
        self.assertTrue(pod['securityContext']['runAsNonRoot'])
        config = yaml.safe_load((APP / 'livekit.yaml').read_text())
        self.assertEqual(config['rtc']['udp_port'], 7882)
        self.assertEqual(config['rtc']['tcp_port'], 7881)
        self.assertFalse(config['turn']['enabled'])
        self.assertNotIn('keys', config)
        self.assertNotIn('redis', config)
        self.assertEqual(config['key_file'], '/etc/livekit-keys/keys.yaml')

    def test_manual_activation_private_https_and_explicit_host_exception(self):
        argo = yaml.safe_load((ROOT / 'platform/argocd/platform-livekit.yaml').read_text())
        self.assertNotIn('automated', argo['spec']['syncPolicy'])
        service = yaml.safe_load((APP / 'service.yaml').read_text())
        self.assertEqual(service['spec']['type'], 'ClusterIP')
        self.assertEqual([p['port'] for p in service['spec']['ports']], [7880])
        ingress = yaml.safe_load((APP / 'ingress.yaml').read_text())
        self.assertEqual(ingress['spec']['ingressClassName'], 'tailscale')
        namespace = yaml.safe_load((APP / 'namespace.yaml').read_text())
        labels = namespace['metadata']['labels']
        self.assertEqual(labels['pod-security.kubernetes.io/enforce'], 'privileged')
        self.assertEqual(labels['pod-security.kubernetes.io/audit'], 'restricted')
        apps = yaml.safe_load((ROOT / 'platform/hermes/mcp/apps.yaml').read_text())['apps']
        self.assertFalse(next(a for a in apps if a['id'] == 'livekit')['enabled'])

    @unittest.skipUnless(shutil.which('kubectl'), 'Native render requires kubectl; source contracts always run')
    def test_render_has_no_development_credentials_or_persistent_state(self):
        rendered = subprocess.check_output(['kubectl', 'kustomize', str(APP)], text=True)
        resources = [r for r in yaml.safe_load_all(rendered) if r]
        self.assertEqual(len([r for r in resources if r['kind'] == 'Deployment']), 1)
        self.assertFalse(any(r['kind'] in {'Secret', 'PersistentVolumeClaim', 'ClusterRoleBinding'} for r in resources))
        config = next(r for r in resources if r['kind'] == 'ConfigMap')
        self.assertIn('livekit.yaml', config['data'])

    def test_sealer_payload_has_exact_namespace_and_no_embedded_default_keys(self):
        spec = importlib.util.spec_from_file_location('prepare_livekit', ROOT / 'scripts/prepare_livekit_secret.py')
        helper = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(helper)
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / 'private.yaml'
            path.write_text('fixture: ' + 'x' * 32)
            data = helper.payload(path)
            self.assertEqual(data['metadata'], {'name': 'livekit-api', 'namespace': 'livekit'})
            self.assertEqual(set(data['data']), {'keys.yaml'})
            for invalid in ['[]', '{}', 'fixture: short', 'fixture: 12']:
                path.write_text(invalid)
                with self.assertRaises(ValueError):
                    helper.payload(path)


if __name__ == '__main__':
    unittest.main()
