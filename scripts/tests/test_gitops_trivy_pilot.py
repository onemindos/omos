"""Render the pinned scanner and verify its effective pilot security boundary."""
import copy
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest
import urllib.request

import yaml

ROOT = Path(__file__).resolve().parents[2]
SCANNER = {'kind': 'ServiceAccount', 'name': 'trivy-operator', 'namespace': 'trivy-system'}


def permits(documents, namespace, group, resource, verb, name=None):
    """Resolve explicit scanner bindings, including a ClusterRole bound by RoleBinding."""
    roles = {(d['kind'], d['metadata'].get('namespace', ''), d['metadata']['name']): d
             for d in documents if d['kind'] in ('Role', 'ClusterRole')}
    for binding in documents:
        kind = binding['kind']
        if kind not in ('RoleBinding', 'ClusterRoleBinding'):
            continue
        if SCANNER not in binding.get('subjects', []):
            continue
        if kind == 'RoleBinding' and binding['metadata']['namespace'] != namespace:
            continue
        ref = binding['roleRef']
        role_ns = binding['metadata']['namespace'] if ref['kind'] == 'Role' else ''
        role = roles[(ref['kind'], role_ns, ref['name'])]
        for rule in role['rules']:
            if (group in rule['apiGroups'] or '*' in rule['apiGroups']) and (
                resource in rule['resources'] or '*' in rule['resources']
            ) and (verb in rule['verbs'] or '*' in rule['verbs']):
                if 'resourceNames' not in rule or name in rule['resourceNames']:
                    return True
    return False


class TrivyPilotTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.values = yaml.safe_load((ROOT / 'platform/ecosystem/values/trivy-operator.yaml').read_text())
        cls.app = yaml.safe_load((ROOT / 'platform/security/ecosystem-trivy-operator.yaml').read_text())
        cls.lock = next(c for c in json.loads((ROOT / 'platform/ecosystem/foundation-lock.json').read_text())
                        if c['name'] == 'trivy-operator')
        cls.tmp = tempfile.TemporaryDirectory(prefix='trivy-pilot-test-')
        cls.addClassCleanup(cls.tmp.cleanup)
        local = os.environ.get('TRIVY_CHART_ARCHIVE')
        if local:
            archive = Path(local).read_bytes()
        else:
            with urllib.request.urlopen(cls.lock['url'], timeout=60) as response:
                archive = response.read()
        if hashlib.sha256(archive).hexdigest() != cls.lock['sha256']:
            raise ValueError('Pinned Trivy chart checksum mismatch')
        chart = Path(cls.tmp.name) / 'chart.tgz'
        chart.write_bytes(archive)
        result = subprocess.run([
            'helm', 'template', 'trivy-operator', str(chart), '--namespace', 'trivy-system',
            '--include-crds', '--kube-version', '1.36.0', '-f',
            str(ROOT / 'platform/ecosystem/values/trivy-operator.yaml'),
        ], capture_output=True, text=True, check=True, timeout=60)
        cls.rendered = [d for d in yaml.safe_load_all(result.stdout) if d]
        cls.rbac = cls.rendered + list(yaml.safe_load_all(
            (ROOT / 'platform/security/trivy-pilot-rbac.yaml').read_text()))

    def resource(self, kind, name):
        return next(d for d in self.rendered if d['kind'] == kind and d['metadata']['name'] == name)

    def test_locked_chart_is_discovered_by_standard_render_gate(self):
        source = self.app['spec']['sources'][0]
        self.assertEqual(source['targetRevision'], self.lock['version'])
        self.assertEqual(source['repoURL'], self.lock['repo'])
        self.assertEqual(self.app['spec']['destination']['namespace'], self.lock['namespace'])
        spec = importlib.util.spec_from_file_location('helm_gate', ROOT / 'scripts/validate_helm_renders.py')
        gate = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(gate)
        self.assertIn('platform/security/ecosystem-trivy-operator.yaml', [p for p, _ in gate.load_apps()])

    def test_runtime_configuration_scans_only_two_namespaces_without_optional_scanners(self):
        dep = self.resource('Deployment', 'trivy-operator')
        env = {e['name']: e['value'] for e in dep['spec']['template']['spec']['containers'][0]['env']}
        self.assertEqual(set(env['OPERATOR_TARGET_NAMESPACES'].split(',')), {'nextcloud', 'n8n'})
        self.assertEqual(env['OPERATOR_TARGET_WORKLOADS'], 'replicaset,statefulset')
        config = self.resource('ConfigMap', 'trivy-operator-config')['data']
        self.assertEqual(config['OPERATOR_VULNERABILITY_SCANNER_ENABLED'], 'true')
        self.assertEqual(config['OPERATOR_CONCURRENT_SCAN_JOBS_LIMIT'], '1')
        for key in ['SBOM_GENERATION', 'CLUSTER_SBOM_CACHE', 'CONFIG_AUDIT_SCANNER',
                    'RBAC_ASSESSMENT_SCANNER', 'INFRA_ASSESSMENT_SCANNER',
                    'CLUSTER_COMPLIANCE', 'EXPOSED_SECRET_SCANNER']:
            self.assertEqual(config['OPERATOR_' + key + '_ENABLED'], 'false', key)
        self.assertEqual(config['OPERATOR_ACCESS_GLOBAL_SECRETS_SERVICE_ACCOUNTS'], 'false')
        self.assertEqual(config['OPERATOR_METRICS_VULN_ID_ENABLED'], 'false')
        self.assertFalse(any(d['kind'] in ('ClusterComplianceReport', 'StatefulSet', 'PersistentVolumeClaim',
                                         'Ingress', 'ValidatingWebhookConfiguration',
                                         'MutatingWebhookConfiguration') for d in self.rendered))

    def test_scanner_can_observe_workloads_and_publish_findings_but_cannot_change_apps(self):
        for namespace in ['nextcloud', 'n8n', 'trivy-system']:
            self.assertTrue(permits(self.rbac, namespace, 'apps', 'replicasets', 'watch'))
            self.assertTrue(permits(self.rbac, namespace, 'batch', 'jobs', 'list'))
            self.assertTrue(permits(self.rbac, namespace, 'aquasecurity.github.io', 'vulnerabilityreports', 'create'))
            for report in ['sbomreports', 'exposedsecretreports']:
                self.assertTrue(permits(self.rbac, namespace, 'aquasecurity.github.io', report, 'watch'))
                self.assertFalse(permits(self.rbac, namespace, 'aquasecurity.github.io', report, 'create'))
            self.assertFalse(permits(self.rbac, namespace, 'apps', 'deployments', 'patch'))
        for namespace in ['nextcloud', 'n8n', 'auth', 'kube-system', 'backstage']:
            for verb in ['get', 'list', 'create', 'update', 'delete']:
                self.assertFalse(permits(self.rbac, namespace, '', 'secrets', verb), (namespace, verb))
            self.assertFalse(permits(self.rbac, namespace, 'batch', 'jobs', 'create'))
            self.assertFalse(permits(self.rbac, namespace, '', 'pods/log', 'get'))
        self.assertFalse(permits(self.rbac, 'auth', 'aquasecurity.github.io', 'vulnerabilityreports', 'create'))
        self.assertFalse(permits(self.rbac, '', '', 'nodes/proxy', 'get'))
        self.assertFalse(permits(self.rbac, '', '', 'nodes', 'list'))
        self.assertTrue(permits(self.rbac, 'trivy-system', 'batch', 'jobs', 'create'))
        self.assertTrue(permits(self.rbac, 'trivy-system', '', 'pods/log', 'get'))
        self.assertTrue(permits(self.rbac, 'trivy-system', '', 'secrets', 'create'))

    def test_permission_resolver_detects_accidental_clusterwide_binding(self):
        documents = copy.deepcopy(self.rbac)
        binding = next(d for d in documents if d['kind'] == 'RoleBinding'
                       and d['metadata']['name'] == 'trivy-pilot-jobs')
        binding['kind'] = 'ClusterRoleBinding'
        binding['roleRef']['kind'] = 'ClusterRole'
        role = next(d for d in documents if d['kind'] == 'Role' and d['metadata']['name'] == 'trivy-pilot-jobs')
        role['kind'] = 'ClusterRole'
        role['metadata'].pop('namespace')
        self.assertTrue(permits(documents, 'auth', '', 'secrets', 'get'))
        self.assertTrue(permits(documents, 'auth', 'batch', 'jobs', 'create'))

    def test_scan_configuration_has_pinned_images_resources_and_no_token_or_host_access(self):
        pod = self.resource('Deployment', 'trivy-operator')['spec']['template']['spec']
        self.assertIn('@sha256:', pod['containers'][0]['image'])
        self.assertTrue(pod['securityContext']['runAsNonRoot'])
        self.assertEqual(pod['securityContext']['seccompProfile']['type'], 'RuntimeDefault')
        scan = self.resource('ConfigMap', 'trivy-operator')['data']
        self.assertEqual(scan['scanJob.automountServiceAccountToken'], 'false')
        self.assertEqual(scan['vulnerabilityReports.scanJobsInSameNamespace'], 'false')
        context = json.loads(scan['scanJob.podTemplatePodSecurityContext'])
        self.assertTrue(context['runAsNonRoot'])
        self.assertEqual(context['seccompProfile']['type'], 'RuntimeDefault')
        image = self.resource('ConfigMap', 'trivy-operator-trivy-config')['data']
        self.assertIn('@sha256:', image['trivy.tag'])
        self.assertEqual(image['trivy.command'], 'image')
        for key in ['requests.cpu', 'requests.memory', 'requests.ephemeral-storage',
                    'limits.cpu', 'limits.memory', 'limits.ephemeral-storage']:
            self.assertTrue(image['trivy.resources.' + key])
        self.assertFalse(any('hostPath' in v for v in pod.get('volumes', [])))

    def test_monitor_and_dashboard_preserve_namespace_labels_and_no_data(self):
        monitor = self.resource('ServiceMonitor', 'trivy-operator')
        self.assertEqual(monitor['metadata']['labels']['release'], 'monitoring')
        self.assertTrue(monitor['spec']['endpoints'][0]['honorLabels'])
        service = self.resource('Service', 'trivy-operator')
        for key, value in monitor['spec']['selector']['matchLabels'].items():
            self.assertEqual(service['metadata']['labels'][key], value)
        cm = yaml.safe_load((ROOT / 'platform/monitoring/grafana-dashboard-trivy-pilot.yaml').read_text())
        board = json.loads(cm['data']['trivy-image-pilot.json'])
        self.assertEqual(board['uid'], 'onemind-trivy-pilot')
        for panel in board['panels']:
            for target in panel.get('targets', []):
                query = target['expr']
                self.assertNotIn('or vector(0)', query)
                if 'trivy_image_vulnerabilities' in query:
                    self.assertIn('namespace=~"nextcloud|n8n"', query)
        roles = [d for d in self.rendered if d['kind'] == 'ClusterRole']
        self.assertEqual(len(roles), 3)  # upstream report-view aggregation, not operator rights
        for role in roles:
            for rule in role['rules']:
                self.assertEqual(set(rule['verbs']), {'get', 'list', 'watch'})
                self.assertEqual(rule['apiGroups'], ['aquasecurity.github.io'])


if __name__ == '__main__':
    unittest.main()
