"""Exercise the deployed observer script without a cluster or credentials."""
import contextlib
import io
import json
import os
from pathlib import Path
import unittest
from unittest.mock import Mock, mock_open, patch

import yaml


class HealthKeeperTest(unittest.TestCase):
    def observe(self, pods):
        documents = list(yaml.safe_load_all(Path(__file__).with_name('health-keeper.yaml').read_text()))
        script = next(d['data']['run.py'] for d in documents if d['kind'] == 'ConfigMap')
        role = next(d for d in documents if d['kind'] == 'ClusterRole')
        self.assertTrue(all(set(rule['verbs']) <= {'get', 'list'} for rule in role['rules']))
        requests = []

        def respond(request, **kwargs):
            requests.append((request.get_method(), request.full_url))
            self.assertEqual(request.get_method(), 'GET')
            self.assertEqual(kwargs['timeout'], 10)
            payload = {'items': pods if request.full_url.endswith('/pods') else []}
            response = Mock()
            response.__enter__ = Mock(return_value=io.BytesIO(json.dumps(payload).encode()))
            response.__exit__ = Mock(return_value=False)
            return response

        stdout = io.StringIO()
        with patch.dict(os.environ, {'KUBERNETES_URL': 'https://cluster.invalid'}), \
                patch('builtins.open', mock_open(read_data='test-token')), \
                patch('ssl.create_default_context', return_value=Mock()), \
                patch('urllib.request.urlopen', side_effect=respond), \
                contextlib.redirect_stdout(stdout):
            exec(compile(script, 'deployed-health-keeper', 'exec'), {})
        self.assertEqual(len(requests), 3)
        return json.loads(stdout.getvalue())['findings']

    def test_pending_and_failed_job_pods_are_preserved(self):
        pods = [{'metadata': {'name': name, 'namespace': 'nats'},
                 'status': {'phase': phase, 'reason': reason}}
                for name, phase, reason in [('pending', 'Pending', 'Unschedulable'),
                                            ('backup-failed', 'Failed', 'Error')]]
        findings = self.observe(pods)
        self.assertEqual({f['pod'] for f in findings}, {'pending', 'backup-failed'})
        self.assertTrue(all(f['action'] == 'pod-unhealthy' for f in findings))

    def test_running_crash_loop_is_reported(self):
        findings = self.observe([{'metadata': {'name': 'app', 'namespace': 'omos'},
                                 'status': {'phase': 'Running',
                                            'conditions': [{'type': 'Ready', 'status': 'False'}],
                                            'containerStatuses': [{'state': {'waiting': {'reason': 'CrashLoopBackOff'}}}]}}])
        self.assertEqual(findings[0]['reasons'], ['CrashLoopBackOff'])

    def test_ready_and_completed_pods_are_not_incidents(self):
        self.assertEqual(self.observe([
            {'metadata': {'name': 'app', 'namespace': 'omos'},
             'status': {'phase': 'Running', 'conditions': [{'type': 'Ready', 'status': 'True'}]}},
            {'metadata': {'name': 'backup', 'namespace': 'nats'}, 'status': {'phase': 'Succeeded'}},
        ]), [])


if __name__ == '__main__':
    unittest.main()
