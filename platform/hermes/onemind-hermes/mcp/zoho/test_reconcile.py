"""Focused tests; use stdlib unittest, with PyYAML for manifest checks."""
import copy
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

import yaml

HERE = Path(__file__).parent
SPEC = importlib.util.spec_from_file_location('zoho_reconcile', HERE / 'reconcile.py')
assert SPEC is not None and SPEC.loader is not None
MODULE = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(MODULE)
BASE = HERE.parents[1]
MANIFEST = json.loads((HERE / 'servers.json').read_text())


class ReconcileTests(unittest.TestCase):
    def test_all_nine_exact_names(self):
        self.assertEqual(set(MANIFEST), {'zoho-ops', 'zoho-data-analytics', 'zoho-desk',
            'zoho-projects', 'zoho-crm', 'zoho-finance', 'zoho-salesiq',
            'zoho-workdrive', 'zoho-mail-calendar'})

    def test_cli_scopes_to_zoho_and_preserves_environment_references(self):
        calls = []
        def runner(argv, **kwargs):
            calls.append((argv, kwargs))
            return subprocess.CompletedProcess(argv, 0)
        MODULE.reconcile(HERE / 'servers.json', '/test-home', '/test-cli', runner)
        self.assertEqual(len(calls), 9)
        for argv, kwargs in calls:
            self.assertTrue(argv[3].startswith('mcp_servers.zoho-'))
            self.assertEqual(json.loads(argv[4]), MANIFEST[argv[3].split('.', 1)[1]])
            self.assertEqual(kwargs['env']['HERMES_HOME'], '/test-home')
            self.assertTrue(kwargs['capture_output'])

    def test_cli_failure_fails_closed_without_echoing_output(self):
        def runner(argv, **kwargs):
            return subprocess.CompletedProcess(argv, 1, '', 'private credential')
        with self.assertRaisesRegex(RuntimeError, '^Hermes config update failed for zoho-ops$'):
            MODULE.reconcile(HERE / 'servers.json', '/test-home', '/test-cli', runner)

    def test_seed_matches_reconciler(self):
        defaults = yaml.safe_load((BASE / 'config.defaults.yaml').read_text())
        for name, settings in MANIFEST.items():
            self.assertEqual(defaults['mcp_servers'][name], settings)
        self.assertTrue({'ops', 'data', 'tak'} <= set(defaults['mcp_servers']))

    def test_secrets_resources_and_environment_bindings(self):
        kustomization = yaml.safe_load((BASE / 'kustomization.yaml').read_text())
        deployment = yaml.safe_load((BASE / 'deployment.yaml').read_text())
        spec = deployment['spec']['template']['spec']
        env_names = {r['secretRef']['name'] for r in spec['containers'][0]['envFrom']}
        for name, settings in MANIFEST.items():
            relative = 'mcp/zoho/' + name + '.sealedsecret.yaml'
            self.assertIn(relative, kustomization['resources'])
            secret = yaml.safe_load((BASE / relative).read_text())
            key = name.upper().replace('-', '_') + '_MCP_URL'
            self.assertEqual(secret['metadata']['namespace'], 'hermes')
            self.assertEqual(secret['metadata']['name'], name + '-mcp')
            self.assertIn(name + '-mcp', env_names)
            self.assertEqual(set(secret['spec']['encryptedData']), {key})
            self.assertEqual(secret['metadata']['annotations']['argocd.argoproj.io/sync-wave'], '-2')
            self.assertNotIn('stringData', secret['spec'])
        self.assertIn('onemind-hermes-runtime', env_names)
        self.assertIn('hermes-dashboard-auth', env_names)

    def test_startup_reconciles_existing_volume(self):
        spec = yaml.safe_load((BASE / 'deployment.yaml').read_text())['spec']['template']['spec']
        init = next(c for c in spec['initContainers'] if c['name'] == 'initial-config')
        script = init['command'][-1]
        self.assertIn("subprocess.run([", script)
        self.assertIn("'/zoho/reconcile.py'", script)
        # Reconciliation must be outside the initial empty-volume seed conditional.
        namespace = {'__name__': '__main__'}
        compile(script, '<initial-config>', 'exec')
        mount = next(m for m in init['volumeMounts'] if m['name'] == 'zoho-config')
        self.assertEqual(mount['mountPath'], '/zoho')
        self.assertTrue(mount['readOnly'])

    def test_read_only_explicit_policy(self):
        for settings in MANIFEST.values():
            self.assertFalse(settings['sampling']['enabled'])
            self.assertFalse(settings['tools']['resources'])
            self.assertFalse(settings['tools']['prompts'])
            for tool in settings['tools']['include']:
                self.assertNotIn('*', tool)
                self.assertRegex(tool, r'^Zoho\w+_(?:get|Get|list|List|Search|Fetch|search|fetch)')
                self.assertNotRegex(tool.lower(), r'delete|create|update|send|execute|payment_link|card')


if __name__ == '__main__':
    unittest.main()
