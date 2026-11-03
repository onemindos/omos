"""Hosted Whop MCP recovery contracts: no OAuth, live state, or cleanup."""
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

import yaml

HERE = Path(__file__).parent
BASE = HERE.parents[1]


def stage_directory():
    scratch = Path(os.environ.get('TMPDIR') or os.environ.get('RUNNER_TEMP')
                   or Path.home() / '.cache/hermes/scratch')
    scratch.mkdir(parents=True, exist_ok=True)
    return Path(tempfile.mkdtemp(prefix='whop-mcp-contract-', dir=scratch))


def load():
    spec = importlib.util.spec_from_file_location('whop_reconcile', HERE / 'reconcile.py')
    assert spec is not None and spec.loader is not None
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


class ReconcileTests(unittest.TestCase):
    def test_only_whop_supported_hermes_config_update(self):
        calls = []
        def runner(argv, **kwargs):
            calls.append((argv, kwargs))
            return subprocess.CompletedProcess(argv, 0)
        load().reconcile(HERE / 'server.json', '/isolated-home', '/hermes', runner)
        self.assertEqual(len(calls), 1)
        argv, kwargs = calls[0]
        self.assertEqual(argv[:4], ['/hermes', 'config', 'set', 'mcp_servers.whop'])
        self.assertEqual(json.loads(argv[4]), {
            'url': 'https://mcp.whop.com/mcp', 'auth': 'oauth',
            'connect_timeout': 300, 'oauth': {'timeout': 300},
            'sampling': {'enabled': False}})
        self.assertEqual(kwargs['env']['HERMES_HOME'], '/isolated-home')
        self.assertTrue(kwargs['capture_output'])

    def test_error_output_suppressed(self):
        def runner(argv, **kwargs):
            return subprocess.CompletedProcess(argv, 1, 'provider-output', 'provider-error')
        with self.assertRaisesRegex(RuntimeError, '^Hermes config update failed for whop$'):
            load().reconcile(HERE / 'server.json', '/isolated-home', '/hermes', runner)

    def test_private_or_unexpected_config_rejected_before_cli(self):
        manifest = stage_directory() / 'invalid.json'
        manifest.write_text(json.dumps({'url': 'https://example.invalid/mcp'}))
        def runner(*args, **kwargs):
            self.fail('Invalid manifest reached CLI')
        with self.assertRaisesRegex(ValueError, 'public Whop OAuth'):
            load().reconcile(manifest, '/isolated-home', '/hermes', runner)

    def test_actual_hermes_cli_preserves_other_settings_and_private_state(self):
        cli = Path(os.environ.get('HERMES_TEST_CLI', '/opt/hermes/bin/hermes'))
        if not cli.is_file():
            self.skipTest('Actual Hermes CLI not present; invocation contracts still run')
        stage = stage_directory()
        original = {'model': {'provider': 'openai-codex', 'default': 'unchanged'},
                    'mcp_servers': {'zoho-ops': {'url': '${ZOHO_OPS_MCP_URL}'},
                                    'other': {'url': 'https://example.invalid/mcp'}},
                    'custom': {'retained': True}}
        config = stage / 'config.yaml'
        config.write_text(yaml.safe_dump(original))
        config.chmod(0o600)
        tokens = stage / 'mcp-tokens'
        tokens.mkdir(mode=0o700)
        private = tokens / 'whop-opaque-fixture.json'
        private.write_bytes(b'opaque-fixture-not-a-credential')
        private.chmod(0o600)
        for _ in range(2):
            load().reconcile(HERE / 'server.json', stage, cli)
            updated = yaml.safe_load(config.read_text())
            self.assertEqual(updated['mcp_servers'].pop('whop'),
                             json.loads((HERE / 'server.json').read_text()))
            self.assertEqual(updated, original)
            self.assertEqual(private.read_bytes(), b'opaque-fixture-not-a-credential')
            self.assertEqual(private.stat().st_mode & 0o777, 0o600)


class GitOpsTests(unittest.TestCase):
    def test_startup_and_encrypted_backup_contract(self):
        spec = yaml.safe_load((BASE / 'deployment.yaml').read_text())['spec']['template']['spec']
        initial = next(c for c in spec['initContainers'] if c['name'] == 'initial-config')
        self.assertIn("'/whop/reconcile.py'", initial['command'][-1])
        self.assertIn({'name': 'whop-config', 'mountPath': '/whop', 'readOnly': True},
                      initial['volumeMounts'])
        defaults = yaml.safe_load((BASE / 'config.defaults.yaml').read_text())
        self.assertEqual(defaults['mcp_servers']['whop'], json.loads((HERE / 'server.json').read_text()))
        backup = yaml.safe_load((BASE / 'backup.yaml').read_text())
        self.assertEqual(backup['spec']['sourcePVC'], 'onemind-hermes-data')
        self.assertEqual(backup['spec']['restic']['repository'], 'onemind-hermes-restic')
        kustomization = yaml.safe_load((BASE / 'kustomization.yaml').read_text())
        self.assertIn('backup.yaml', kustomization['resources'])
        self.assertIn('backup.sealedsecret.yaml', kustomization['resources'])
        cm = next(c for c in kustomization['configMapGenerator'] if c['name'] == 'onemind-hermes-whop')
        self.assertEqual(set(cm['files']), {'reconcile.py=mcp/whop/reconcile.py',
                                           'server.json=mcp/whop/server.json'})
        sealed = yaml.safe_load((BASE / 'backup.sealedsecret.yaml').read_text())
        self.assertEqual(sealed['kind'], 'SealedSecret')
        self.assertEqual(sealed['metadata']['name'], backup['spec']['restic']['repository'])
        self.assertIn('RESTIC_PASSWORD', sealed['spec']['encryptedData'])

    def test_no_whop_standalone_cli_or_auth_wiring(self):
        spec = yaml.safe_load((BASE / 'deployment.yaml').read_text())['spec']['template']['spec']
        self.assertNotIn('install-whop-cli', {c['name'] for c in spec['initContainers']})
        main = spec['containers'][0]
        names = {e['name'] for e in main['env']}
        self.assertFalse(any(name.startswith('WHOP_') for name in names))
        self.assertNotIn('PATH', names)
        cm = next(c for c in yaml.safe_load((BASE / 'kustomization.yaml').read_text())
                  ['configMapGenerator'] if c['name'] == 'onemind-hermes-whop')
        self.assertFalse(any('install_cli' in item for item in cm['files']))
        self.assertNotIn('login', (HERE / 'reconcile.py').read_text())


if __name__ == '__main__':
    unittest.main()
