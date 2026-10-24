"""Safety contracts for the disposable native-crash diagnostic fixture."""
from copy import deepcopy
from pathlib import Path
import shutil
import subprocess
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
FIXTURE = ROOT / 'platform/stac-validation'


def load(path):
    return yaml.safe_load(path.read_text())


class Isolation(unittest.TestCase):
    def test_only_disposable_resources_and_manual_application_are_declared(self):
        docs = [load(p) for p in FIXTURE.glob('*.yaml')]
        self.assertCountEqual([d['kind'] for d in docs],
                              ['Namespace', 'NetworkPolicy', 'ConfigMap', 'SealedSecret', 'Job', 'Job'])
        for doc in docs:
            self.assertEqual(doc['metadata'].get('namespace', 'stac-validation'), 'stac-validation')
        app = load(ROOT / 'platform/argocd/platform-stac-validation.yaml')
        self.assertNotIn('finalizers', app['metadata'])
        self.assertEqual(app['spec']['syncPolicy'], {'retry': {'limit': 0}})
        self.assertEqual(app['spec']['source']['path'], 'platform/stac-validation')
        mapping = load(ROOT / 'services/onemind-portal/catalog/workload-map.yaml')
        self.assertEqual(mapping['argoApplications']['platform-stac-validation'], 'component:default/stac')
        policy = next(d for d in docs if d['kind'] == 'NetworkPolicy')
        self.assertEqual(policy['spec'], {'podSelector': {}, 'policyTypes': ['Ingress', 'Egress'],
                                          'ingress': [], 'egress': []})
        registry = next(d for d in docs if d['kind'] == 'SealedSecret')
        self.assertEqual(registry['metadata']['name'], 'registry-pull')
        self.assertEqual(set(registry['spec']['encryptedData']), {'.dockerconfigjson'})
        self.assertEqual(registry['spec']['template']['type'], 'kubernetes.io/dockerconfigjson')
        self.assertNotIn('bitnami.com/sealed-secrets/namespace-wide', registry['metadata'].get('annotations', {}))
        self.assertNotIn('bitnami.com/sealed-secrets/cluster-wide', registry['metadata'].get('annotations', {}))

    def test_jobs_have_no_database_credentials_external_access_or_persistent_storage(self):
        image = load(ROOT / 'platform/ecosystem/resources/postgres-image-catalog.yaml')['spec']['images'][0]['image']
        client = load(ROOT / 'platform/stac/migrate.yaml')['spec']['template']['spec']['containers'][0]['image']
        for case in ['control', 'age']:
            job = load(FIXTURE / f'{case}.yaml')
            self.assertEqual(job['spec']['backoffLimit'], 0)
            self.assertEqual(job['spec']['activeDeadlineSeconds'], 600)
            pod = job['spec']['template']['spec']
            self.assertEqual(pod['restartPolicy'], 'Never')
            self.assertIs(pod['automountServiceAccountToken'], False)
            for forbidden in ['hostNetwork', 'hostPID', 'hostIPC']:
                self.assertNotIn(forbidden, pod)
            self.assertEqual(pod['imagePullSecrets'], [{'name': 'registry-pull'}])
            for volume in pod['volumes']:
                self.assertTrue(set(volume) <= {'name', 'emptyDir', 'configMap'})
            server = pod['initContainers'][0]
            self.assertEqual(server['image'], 'ghcr.io/onemindos/onemind-postgres@' + image.split('@', 1)[1])
            self.assertEqual(server['imagePullPolicy'], 'IfNotPresent')
            self.assertEqual(server['restartPolicy'], 'Always')
            self.assertIn('startupProbe', server)
            self.assertEqual(server['resources']['limits']['memory'], '1Gi')
            migration = pod['containers'][0]
            self.assertEqual(migration['image'], client)
            self.assertEqual(migration['resources']['limits']['memory'], '512Mi')
            self.assertEqual(next(e['value'] for e in migration['env'] if e['name'] == 'PGHOST'), '127.0.0.1')
            for container in [server, migration]:
                self.assertNotIn('envFrom', container)
                self.assertTrue(all('valueFrom' not in e for e in container['env']))
                self.assertTrue(container['securityContext']['readOnlyRootFilesystem'])
                self.assertEqual(container['securityContext']['capabilities']['drop'], ['ALL'])

    def test_cases_differ_only_in_identity_preload_and_wave(self):
        control = load(FIXTURE / 'control.yaml')
        age = load(FIXTURE / 'age.yaml')
        self.assertEqual(control['metadata']['annotations']['argocd.argoproj.io/sync-wave'], '1')
        self.assertEqual(age['metadata']['annotations']['argocd.argoproj.io/sync-wave'], '2')
        normalized = []
        for doc in [control, age]:
            doc = deepcopy(doc)
            doc['metadata'] = {}
            doc['spec']['template']['metadata'] = {}
            pod = doc['spec']['template']['spec']
            pod['initContainers'][0]['env'] = []
            pod['containers'][0]['env'] = [e for e in pod['containers'][0]['env'] if e['name'] != 'VALIDATION_CASE']
            normalized.append(doc)
        self.assertEqual(*normalized)
        self.assertEqual(control['spec']['template']['spec']['initContainers'][0]['env'], [{'name': 'AGE_PRELOAD', 'value': ''}])
        self.assertEqual(age['spec']['template']['spec']['initContainers'][0]['env'], [{'name': 'AGE_PRELOAD', 'value': 'age'}])

    def test_embedded_shell_is_valid_and_preserves_loopback_and_attempt_barrier(self):
        scripts = load(FIXTURE / 'scripts.yaml')['data']
        for name, script in scripts.items():
            result = subprocess.run(['bash', '-n'], input=script, text=True, capture_output=True)
            self.assertEqual(result.returncode, 0, name + result.stderr)
        postgres = scripts['postgres.sh']
        self.assertIn('listen_addresses=127.0.0.1', postgres)
        self.assertIn('host all all 127.0.0.1/32 trust', postgres)
        self.assertIn('restart_after_crash=off', postgres)
        self.assertLess(postgres.index('test -e /scratch/attempted'), postgres.index('initdb -D'))
        self.assertIn('log_min_error_statement=error', postgres)
        migration = scripts['migrate.sh']
        self.assertIn('NOSUPERUSER NOCREATEDB', migration)
        for extension in ['postgis', 'btree_gist', 'unaccent']:
            self.assertIn(f'CREATE EXTENSION {extension};', migration)
        self.assertNotIn('CREATE EXTENSION age', migration)
        self.assertEqual(migration.count('\npypgstac migrate\n'), 2)

    @unittest.skipUnless(shutil.which('shellcheck'), 'ShellCheck is optional; Bash syntax is always checked')
    def test_embedded_shellcheck_when_available(self):
        for name, script in load(FIXTURE / 'scripts.yaml')['data'].items():
            result = subprocess.run(['shellcheck', '-s', 'bash', '-'], input=script, text=True, capture_output=True)
            self.assertEqual(result.returncode, 0, name + result.stdout + result.stderr)


if __name__ == '__main__':
    unittest.main()
