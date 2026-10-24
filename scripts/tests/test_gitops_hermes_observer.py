"""Negative privilege and evidence-redaction checks for Hermes observation."""
import importlib.util
import io
import json
from pathlib import Path
import time
import unittest
from unittest.mock import patch
import urllib.error

import yaml

ROOT = Path(__file__).resolve().parents[2]
AGENT = ROOT / 'platform/hermes/agent'
SPEC = importlib.util.spec_from_file_location('cluster_observer', AGENT / 'cluster_observer.py')
observer = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(observer)


class Observer(unittest.TestCase):
    def test_exact_readonly_rbac_and_identity(self):
        docs = list(yaml.safe_load_all((AGENT / 'cluster-observer.yaml').read_text()))
        role = next(d for d in docs if d['kind'] == 'ClusterRole')
        expected = {
            '': {'nodes', 'namespaces', 'pods', 'services', 'endpoints', 'events',
                 'persistentvolumeclaims', 'persistentvolumes'},
            'apps': {'deployments', 'statefulsets', 'daemonsets', 'replicasets'},
            'batch': {'jobs', 'cronjobs'}, 'networking.k8s.io': {'ingresses', 'networkpolicies'},
            'discovery.k8s.io': {'endpointslices'}, 'metrics.k8s.io': {'nodes', 'pods'},
            'argoproj.io': {'applications'}, 'postgresql.cnpg.io': {'clusters'},
            'volsync.backube': {'replicationsources'},
        }
        actual = {}
        for rule in role['rules']:
            self.assertEqual(set(rule['verbs']), {'get', 'list', 'watch'})
            self.assertNotIn('nonResourceURLs', rule)
            for group in rule['apiGroups']:
                actual[group] = set(rule['resources'])
        self.assertEqual(actual, expected)
        binding = next(d for d in docs if d['kind'] == 'ClusterRoleBinding')
        self.assertEqual(binding['subjects'], [{'kind': 'ServiceAccount', 'name': 'hermes-observer',
                                                'namespace': 'hermes'}])
    def test_raw_sensitive_object_fields_never_leave_summarizer(self):
        sensitive = 'sentinel-value-that-must-not-be-reported'
        item = {'metadata': {'name': 'app', 'namespace': 'demo', 'annotations': {'token': sensitive}},
                'spec': {'containers': [{'env': [{'name': 'PASSWORD', 'value': sensitive}]}]},
                'message': sensitive,
                'status': {'phase': 'Running', 'conditions': [{'type': 'Ready', 'status': 'True',
                                                               'message': sensitive}],
                           'containerStatuses': [{'name': 'app', 'state': {'waiting': {
                               'reason': 'CrashLoopBackOff', 'message': sensitive}}}]}}
        for kind in observer.PATHS:
            result = json.dumps(observer.summarize(kind, item))
            self.assertNotIn(sensitive, result, kind)
            self.assertNotIn('PASSWORD', result, kind)

    def test_missing_api_is_incomplete_not_healthy_and_body_not_printed(self):
        class Unavailable:
            def open(self, request, timeout):
                raise urllib.error.HTTPError(request.full_url, 403, 'sensitive error message', {}, None)
        with patch.object(observer, 'PATHS', {'pods': '/api/v1/pods'}), \
                patch.object(Path, 'read_text', return_value='fake-token'):
            result = observer.collect(Unavailable(), time.monotonic() + 5)
        self.assertFalse(result['pods']['complete'])
        self.assertEqual(result['pods']['httpStatus'], 403)
        self.assertNotIn('sensitive', json.dumps(result))
        self.assertNotIn('fake-token', json.dumps(result))

    def test_pagination_stays_fixed_origin_and_get_only(self):
        class Pages:
            def __init__(self):
                self.requests = []
            def open(self, request, timeout):
                self.requests.append(request)
                body = {'items': [{'metadata': {'name': str(len(self.requests))}}],
                        'metadata': {'continue': 'opaque/page?token' if len(self.requests) == 1 else ''}}
                return io.BytesIO(json.dumps(body).encode())
        opener = Pages()
        with patch.object(observer, 'PATHS', {'pods': '/api/v1/pods'}), \
                patch.object(Path, 'read_text', return_value='fake-token'):
            result = observer.collect(opener, time.monotonic() + 5)
        self.assertTrue(result['pods']['complete'])
        self.assertEqual(result['pods']['count'], 2)
        for request in opener.requests:
            self.assertEqual(request.get_method(), 'GET')
            self.assertTrue(request.full_url.startswith('https://kubernetes.default.svc/api/v1/pods?'))
        self.assertIn('continue=opaque%2Fpage%3Ftoken', opener.requests[1].full_url)

    def test_unbounded_pages_and_deadlines_fail_explicitly(self):
        class Endless:
            def open(self, request, timeout):
                return io.BytesIO(b'{"items": [], "metadata": {"continue": "again"}}')
        with patch.object(observer, 'PATHS', {'pods': '/api/v1/pods'}), \
                patch.object(Path, 'read_text', return_value='fake-token'), \
                patch.object(observer, 'MAX_PAGES', 2):
            result = observer.collect(Endless(), time.monotonic() + 5)
            expired = observer.collect(Endless(), time.monotonic() - 1)
        self.assertFalse(result['pods']['complete'])
        self.assertEqual(result['pods']['error'], 'ValueError')
        self.assertEqual(expired['pods']['error'], 'TimeoutError')

    def test_no_redirect_of_credentials(self):
        self.assertIsNone(observer.NoRedirect().redirect_request(None, None, 302, '', {},
                                                               'https://unrelated.invalid/'))


if __name__ == '__main__':
    unittest.main()
