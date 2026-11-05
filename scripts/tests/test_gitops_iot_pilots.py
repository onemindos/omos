"""Boundaries for the separately activated private IoT evaluation apps."""
import base64
import importlib.util
from pathlib import Path
import secrets
import shutil
import subprocess
import unittest
from unittest.mock import patch
import yaml

ROOT = Path(__file__).resolve().parents[2]
APPS = ('openremote', 'thingsboard')


def resources(app):
    result = []
    for path in (ROOT / 'platform' / app).glob('*.yaml'):
        if path.name != 'kustomization.yaml':
            result.extend(x for x in yaml.safe_load_all(path.read_text()) if x)
    return result


class PrivateIotPilots(unittest.TestCase):
    def test_apps_require_separate_manual_sync_and_retained_data(self):
        for app in APPS:
            argo = yaml.safe_load((ROOT / f'platform/argocd/platform-{app}.yaml').read_text())
            self.assertNotIn('automated', argo['spec']['syncPolicy'])
            rs = resources(app)
            self.assertFalse(any(x['kind'] == 'Secret' for x in rs))
            for claim in (x for x in rs if x['kind'] == 'PersistentVolumeClaim'):
                self.assertEqual(claim['spec']['storageClassName'], 'do-block-storage-retain')
                self.assertIn('Prune=false', claim['metadata']['annotations']['argocd.argoproj.io/sync-options'])
            backup = [x for x in rs if x['kind'] == 'ReplicationSource']
            self.assertEqual({x['spec']['sourcePVC'] for x in backup},
                             {x['metadata']['name'] for x in rs if x['kind'] == 'PersistentVolumeClaim'})
            for item in backup:
                self.assertEqual(item['spec']['restic']['storageClassName'], 'do-block-storage')
                self.assertEqual(item['spec']['restic']['cacheStorageClassName'], 'do-block-storage')

    def test_web_is_tailnet_only_and_no_mqtt_loadbalancer(self):
        for app in APPS:
            for r in resources(app):
                if r['kind'] == 'Ingress':
                    self.assertEqual(r['spec']['ingressClassName'], 'tailscale')
                if r['kind'] == 'Service':
                    self.assertNotIn(r['spec'].get('type'), ('LoadBalancer', 'NodePort'))
            self.assertTrue(any(r['kind'] == 'NetworkPolicy' and
                                r['spec']['podSelector'] == {} and
                                not r['spec'].get('ingress') and not r['spec'].get('egress')
                                for r in resources(app)))

    def test_pods_have_pinned_images_resources_and_no_cluster_authority(self):
        for app in APPS:
            for r in resources(app):
                if r['kind'] not in ('Deployment', 'Job'):
                    continue
                pod = r['spec']['template']['spec']
                self.assertFalse(pod['automountServiceAccountToken'])
                self.assertTrue(pod['securityContext']['runAsNonRoot'])
                for c in pod['containers'] + pod.get('initContainers', []):
                    self.assertIn('@sha256:', c['image'])
                    self.assertIn('requests', c['resources'])
                    self.assertIn('limits', c['resources'])
                    self.assertFalse(c['securityContext']['allowPrivilegeEscalation'])
                    self.assertIn('ALL', c['securityContext']['capabilities']['drop'])

    def test_thingsboard_installer_and_password_replacement_precede_exposure(self):
        rs = resources('thingsboard')
        job = next(x for x in rs if x['kind'] == 'Job')
        app = next(x for x in rs if x['kind'] == 'Deployment' and x['metadata']['name'] == 'thingsboard')
        self.assertEqual(job['metadata']['annotations']['argocd.argoproj.io/hook'], 'Sync')
        self.assertLess(int(job['metadata']['annotations']['argocd.argoproj.io/sync-wave']),
                        int(app['metadata']['annotations']['argocd.argoproj.io/sync-wave']))
        c = app['spec']['template']['spec']['containers'][0]
        self.assertIn('tb-node:4.3.1.8@sha256:', c['image'])
        self.assertFalse(any(x['name'] in ('INSTALL_TB', 'UPGRADE_TB', 'LOAD_DEMO') for x in c['env']))
        data = next(x for x in rs if x['kind'] == 'ConfigMap')['data']
        self.assertIn('refuse installation or credential reset', data['guard.sh'])
        self.assertIn('public.onemind_bootstrap', data['guard.sh'])
        self.assertIn('BEGIN;', data['secure-admin.sql'])
        self.assertIn('COMMIT;', data['secure-admin.sql'])
        self.assertNotIn('DROP ', data['secure-admin.sql'])
        self.assertIn('SYSADMIN_PASSWORD_HASH', data['secure-admin.sql'])

    def test_openremote_uses_production_identity_and_one_https_origin(self):
        rs = resources('openremote')
        manager = next(x for x in rs if x['kind'] == 'Deployment' and x['metadata']['name'] == 'manager')
        env = {x['name']: x.get('value') for x in manager['spec']['template']['spec']['containers'][0]['env']}
        self.assertEqual(env['OR_DEV_MODE'], 'false')
        self.assertEqual(env['OR_SETUP_RUN_ON_RESTART'], 'false')
        self.assertEqual(env['OR_KEYCLOAK_PUBLIC_URI'], 'https://openremote.tail717ba0.ts.net/auth')
        self.assertNotIn('OR_SETUP_TYPE', env)
        keycloak = next(x for x in rs if x['kind'] == 'Deployment' and x['metadata']['name'] == 'keycloak')
        c = keycloak['spec']['template']['spec']['containers'][0]
        self.assertNotIn('start-dev', str(c))
        self.assertEqual(c['readinessProbe']['httpGet'], {'path': '/auth/health/ready', 'port': 'health'})

    @unittest.skipUnless(shutil.which('kubectl'), 'Render needs kubectl; source contracts always run')
    def test_kustomize_render_works_for_both_independent_apps(self):
        for app in APPS:
            out = subprocess.check_output(['kubectl', 'kustomize', str(ROOT / 'platform' / app)], text=True)
            rendered = list(yaml.safe_load_all(out))
            self.assertFalse(any(r['kind'] == 'Secret' for r in rendered))
            self.assertTrue(all(r['metadata'].get('namespace') == app for r in rendered if r['kind'] != 'Namespace'))


class PrivateCredentialPreparation(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        spec = importlib.util.spec_from_file_location('iot_prepare', ROOT / 'scripts/prepare_iot_secrets.py')
        cls.helper = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(cls.helper)

    def test_reject_missing_values_and_shared_backup_repository(self):
        with self.assertRaises(ValueError):
            self.helper.validate('openremote', {})
        values = {name: {key: secrets.token_urlsafe(32) for key in keys}
                  for name, keys in self.helper.REQUIRED['openremote'].items()}
        for name in values:
            if 'RESTIC_REPOSITORY' in values[name]:
                values[name]['RESTIC_REPOSITORY'] = 's3:https://storage.invalid/same'
        with self.assertRaises(ValueError):
            self.helper.validate('openremote', values)

    def test_thingsboard_plaintext_owner_password_is_replaced_with_hash(self):
        values = {name: {key: secrets.token_urlsafe(32) for key in keys}
                  for name, keys in self.helper.REQUIRED['thingsboard'].items()}
        values['thingsboard-env']['SYSADMIN_EMAIL'] = 'owner-' + secrets.token_hex(4) + '@example.invalid'
        values['thingsboard-env']['JWT_TOKEN_SIGNING_KEY'] = base64.b64encode(secrets.token_bytes(32)).decode()
        values['thingsboard-database-restic']['RESTIC_REPOSITORY'] = 's3:https://storage.invalid/isolated'
        result = self.helper.validate('thingsboard', values)
        self.assertNotIn('SYSADMIN_PASSWORD', result['thingsboard-env'])
        self.assertTrue(result['thingsboard-env']['SYSADMIN_PASSWORD_HASH'].startswith('$2'))

    def test_sealing_failure_does_not_echo_subprocess_diagnostics(self):
        with patch.object(self.helper.subprocess, 'run', return_value=subprocess.CompletedProcess([], 1, '', secrets.token_hex(32))):
            with self.assertRaisesRegex(RuntimeError, 'private diagnostics intentionally suppressed'):
                self.helper.seal('openremote', 'openremote-env', {}, Path('/public/cert'))

    def test_unexpected_sealing_output_is_rejected_without_echoing_payload(self):
        private_value = secrets.token_hex(32)
        invalid = yaml.safe_dump({'kind': 'Secret', 'metadata': {'name': 'openremote-env',
                                 'namespace': 'openremote'}, 'stringData': {'DB_PASSWORD': private_value}})
        with patch.object(self.helper.subprocess, 'run', return_value=subprocess.CompletedProcess([], 0, invalid, '')):
            with self.assertRaises(RuntimeError) as caught:
                self.helper.seal('openremote', 'openremote-env', {'DB_PASSWORD': private_value}, Path('/public/cert'))
        self.assertNotIn(private_value, str(caught.exception))


if __name__ == '__main__':
    unittest.main()
