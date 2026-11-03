"""Offline contracts plus real isolated Hermes CLI configuration; no credentials."""
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

import yaml

HERE = Path(__file__).parent
BASE = HERE.parent


def load(name):
    spec = importlib.util.spec_from_file_location(name, HERE / (name + '.py'))
    assert spec is not None and spec.loader is not None
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def stage():
    root = Path(os.environ.get('TMPDIR', '/opt/data/cache/scratch'))
    root.mkdir(parents=True, exist_ok=True)
    return Path(tempfile.mkdtemp(prefix='native-op-test-', dir=root))


class NativeTests(unittest.TestCase):
    def test_actual_cli_empty_and_existing_volume_idempotent(self):
        cli = Path('/opt/hermes/bin/hermes')
        if not cli.is_file():
            self.skipTest('Actual Hermes CLI unavailable')
        for existing in (False, True):
            home = stage()
            original = {'model': {'provider': 'openai-codex', 'default': 'unchanged'},
                        'mcp_servers': {'other': {'url': 'https://example.invalid/mcp'}},
                        'secrets': {'preserve_existing': ['OTHER_TOKEN']},
                        'terminal': {'backend': 'local', 'env_passthrough':
                                     ['EXISTING_APP_TOKEN', 'OP_SERVICE_ACCOUNT_TOKEN']}}
            config = home / 'config.yaml'
            config.write_text(yaml.safe_dump(original) if existing else
                              (BASE / 'config.defaults.yaml').read_text())
            config.chmod(0o600)
            # Empty volume seeding is the existing initial-config contract.
            for _ in range(2):
                load('reconcile').reconcile(HERE / 'settings.json', home, cli)
                updated = yaml.safe_load(config.read_text())
                self.assertEqual(updated['secrets']['onepassword'],
                                 json.loads((HERE / 'settings.json').read_text()))
                self.assertNotIn('OP_SERVICE_ACCOUNT_TOKEN', updated['terminal']['env_passthrough'])
                self.assertIn('CLICKUP_API_TOKEN', updated['terminal']['env_passthrough'])
                if existing:
                    self.assertEqual(updated['model'], original['model'])
                    self.assertEqual(updated['mcp_servers'], original['mcp_servers'])
                    self.assertEqual(updated['secrets']['preserve_existing'], ['OTHER_TOKEN'])
                    self.assertEqual(updated['terminal']['backend'], 'local')
                    self.assertIn('EXISTING_APP_TOKEN', updated['terminal']['env_passthrough'])
                self.assertFalse((home / 'cache/op_cache.json').exists())

    def test_cli_errors_do_not_echo_output(self):
        def runner(argv, **kwargs):
            return subprocess.CompletedProcess(argv, 1, 'untrusted-output', 'untrusted-error')
        with self.assertRaisesRegex(RuntimeError, '^Hermes terminal configuration read failed$'):
            load('reconcile').reconcile(HERE / 'settings.json', '/unused', '/hermes', runner)

    def test_bad_arch_fails_without_download(self):
        def fetch():
            self.fail('Unsupported architecture attempted download')
        with self.assertRaisesRegex(RuntimeError, 'linux/amd64 only'):
            load('install_cli').install(stage(), arch='aarch64', fetch=fetch)

    def test_bad_hash_never_installs(self):
        destination = stage()
        with self.assertRaisesRegex(RuntimeError, 'checksum mismatch'):
            load('install_cli').install(destination, arch='x86_64', fetch=lambda: b'not-an-archive')
        self.assertFalse((destination / 'op').exists())

    def test_active_gitops_is_unchanged_source_only(self):
        # Exact bytes from original base 85a52b599b1b18b8f4b6a38cee9040c195ee3c9b.
        hashes = {
            'deployment.yaml': '2e4212f88053d899d039eb0a496c15f9ee096f460c22c7127b701e78578d284d',
            'config.defaults.yaml': '271af0fa4a9a207045f4e0cfdebe52d9c5f9d21c10b3e47a1cb7f17c7f4d73f6',
            'kustomization.yaml': 'ab9740b777cea039d8b08c44a8a12cdeb896285d7ac41fadec18c392e49a1630',
        }
        for name, digest in hashes.items():
            content = (BASE / name).read_bytes()
            self.assertEqual(hashlib.sha256(content).hexdigest(), digest, name)
            self.assertNotIn(b'onepassword', content.lower(), name)
            self.assertNotIn(b'OP_SERVICE_ACCOUNT_TOKEN', content, name)
            self.assertNotIn(b'install_cli.py', content, name)
        spec = yaml.safe_load((BASE / 'deployment.yaml').read_text())['spec']['template']['spec']
        expected = ('nousresearch/hermes-agent:v0.22.5@sha256:'
                    '9c5561579f3ffdf333401797094ac3a743738d6c25d40f4331c7bf77b694238d')
        self.assertEqual(spec['containers'][0]['image'], expected)
        self.assertNotIn('install-onepassword-cli', [c['name'] for c in spec['initContainers']])
        self.assertNotIn('onepassword', yaml.safe_load(
            (BASE / 'config.defaults.yaml').read_text()).get('secrets', {}))
        for name in ('reconcile.py', 'install_cli.py', 'settings.json',
                     'runtime-fix/Dockerfile', 'runtime-fix/apply_source_patch.py',
                     'runtime-fix/test_child_policy.py'):
            self.assertTrue((HERE / name).is_file(), name)

    def test_non_active_proposal_wiring_and_scope(self):
        proposal = HERE / 'proposal'
        spec = yaml.safe_load((proposal / 'deployment.yaml').read_text())['spec']['template']['spec']
        self.assertEqual(spec['nodeSelector']['kubernetes.io/arch'], 'amd64')
        main = spec['containers'][0]
        env = {entry['name']: entry for entry in main['env']}
        self.assertEqual(env['OP_SERVICE_ACCOUNT_TOKEN']['valueFrom']['secretKeyRef'], {
            'name': 'hermes-onepassword-bootstrap', 'key': 'OP_SERVICE_ACCOUNT_TOKEN'})
        self.assertEqual(env['OP_CONFIG_DIR']['value'], '/var/run/onepassword')
        for init in spec['initContainers']:
            self.assertNotIn('envFrom', init)
            self.assertNotIn('OP_SERVICE_ACCOUNT_TOKEN', {e['name'] for e in init.get('env', [])})
            self.assertIn('@sha256:', init['image'])
            self.assertIn('requests', init['resources'])
            self.assertIn('limits', init['resources'])
        initial = next(c for c in spec['initContainers'] if c['name'] == 'initial-config')
        self.assertIn("'/onepassword/reconcile.py'", initial['command'][-1])
        self.assertIn({'name': 'onepassword-bin', 'mountPath': '/opt/onepassword',
                       'readOnly': True}, main['volumeMounts'])
        defaults = yaml.safe_load((proposal / 'config.defaults.yaml').read_text())
        self.assertEqual(defaults['secrets']['onepassword'], json.loads((HERE / 'settings.json').read_text()))
        self.assertEqual(defaults['terminal']['env_passthrough'], ['CLICKUP_API_TOKEN'])
        cm = next(c for c in yaml.safe_load((proposal / 'kustomization.yaml').read_text())
                  ['configMapGenerator'] if c['name'] == 'onemind-hermes-onepassword')
        self.assertEqual(len(cm['files']), 3)
        for entry in cm['files']:
            self.assertTrue((BASE / entry.split('=', 1)[1]).is_file())


if __name__ == '__main__':
    unittest.main()
