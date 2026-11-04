"""Hosted Todoist recovery contracts; fixtures never connect or authorize."""
import ast
import contextlib
import importlib.util
import io
import tempfile
from types import SimpleNamespace
from unittest.mock import patch

import yaml
import json
import os
from pathlib import Path
import subprocess
import unittest

HERE = Path(__file__).parent
BASE = HERE.parents[1]
EXPECTED = {
    'url': 'https://ai.todoist.net/mcp', 'auth': 'oauth', 'enabled': True,
    'tools': {'exclude': ['search', 'fetch', 'export-project-template', 'import-project-template']},
    'sampling': {'enabled': False},
}


def load():
    spec = importlib.util.spec_from_file_location('todoist_reconcile', HERE / 'reconcile.py')
    assert spec is not None and spec.loader is not None
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


class ReconcileTests(unittest.TestCase):
    def test_only_todoist_supported_hermes_config_update(self):
        self.assertTrue((HERE / 'server.json').is_file(), 'Todoist public manifest missing')
        calls = []
        def runner(argv, **kwargs):
            calls.append((argv, kwargs))
            return subprocess.CompletedProcess(argv, 0)
        load().reconcile(HERE / 'server.json', '/isolated-home', '/hermes', runner)
        self.assertEqual(len(calls), 1)
        argv, kwargs = calls[0]
        self.assertEqual(argv[:4], ['/hermes', 'config', 'set', 'mcp_servers.todoist'])
        self.assertEqual(json.loads(argv[4]), EXPECTED)
        self.assertEqual(kwargs['env']['HERMES_HOME'], '/isolated-home')
        self.assertTrue(kwargs['capture_output'])
        self.assertNotIn('include', EXPECTED['tools'])


def stage_directory():
    scratch = Path(os.environ.get('TMPDIR') or os.environ.get('RUNNER_TEMP')
                   or Path.home() / '.cache/hermes/scratch')
    scratch.mkdir(parents=True, exist_ok=True)
    return Path(tempfile.mkdtemp(prefix='todoist-mcp-contract-', dir=scratch))


class PreservationTests(unittest.TestCase):
    def test_error_output_suppressed(self):
        def runner(argv, **kwargs):
            return subprocess.CompletedProcess(argv, 1, 'provider-output', 'provider-error')
        with contextlib.redirect_stdout(io.StringIO()) as out:
            with self.assertRaisesRegex(RuntimeError, '^Hermes config update failed for todoist$'):
                load().reconcile(HERE / 'server.json', '/isolated-home', '/hermes', runner)
        self.assertEqual(out.getvalue(), '')

    def test_private_or_unexpected_config_rejected_before_cli(self):
        for change in ({'headers': {'Authorization': 'opaque-fixture'}},
                       {'url': 'https://example.invalid/mcp'},
                       {'tools': {'include': ['user-info']}},
                       {'sampling': {'enabled': True}},
                       {'oauth': {'access_token': 'opaque-fixture'}}):
            with self.subTest(fields=list(change)):
                manifest = stage_directory() / 'invalid.json'
                manifest.write_text(json.dumps(dict(EXPECTED, **change)))
                def runner(*args, **kwargs):
                    self.fail('Invalid manifest reached CLI')
                with self.assertRaisesRegex(ValueError, 'public Todoist OAuth'):
                    load().reconcile(manifest, '/isolated-home', '/hermes', runner)

    def test_actual_cli_twice_preserves_other_settings_and_private_state(self):
        cli = Path(os.environ.get('HERMES_TEST_CLI', '/opt/hermes/bin/hermes'))
        if not cli.is_file():
            self.skipTest('Actual Hermes CLI absent; invocation contracts still run')
        stage = stage_directory()
        original = {
            'model': {'provider': 'openai-codex', 'default': 'unchanged'},
            'mcp_servers': {name: {'url': 'https://example.invalid/' + name}
                            for name in ('whop', 'zoho-ops', 'clickup', 'dropbox', 'ops', 'data', 'tak', 'other')},
            'custom': {'retained': True}, 'display': {'skin': 'default'},
        }
        config = stage / 'config.yaml'
        drifted = dict(original, mcp_servers=dict(original['mcp_servers'], todoist={'url': 'https://example.invalid/drift'}))
        config.write_text(yaml.safe_dump(drifted))
        config.chmod(0o600)
        paths = {'mcp-tokens/todoist-opaque-fixture.json': b'opaque-oauth-fixture',
                 'auth.json': b'opaque-provider-fixture', '.env': b'# opaque-environment-fixture\n',
                 'profiles/other/config.yaml': b'other-profile-opaque-fixture'}
        for relative, content in paths.items():
            target = stage / relative
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(content)
            target.chmod(0o600)
        (stage / 'mcp-tokens').chmod(0o700)
        before = {relative: ((stage / relative).read_bytes(), (stage / relative).stat().st_mode,
                             (stage / relative).stat().st_uid, (stage / relative).stat().st_mtime_ns)
                  for relative in paths}
        for _ in range(2):
            load().reconcile(HERE / 'server.json', stage, cli)
            updated = yaml.safe_load(config.read_text())
            self.assertEqual(updated['mcp_servers'].pop('todoist'), EXPECTED)
            self.assertEqual(updated, original)
            self.assertEqual(config.stat().st_mode & 0o777, 0o600)
            for relative, state in before.items():
                target = stage / relative
                self.assertEqual((target.read_bytes(), target.stat().st_mode,
                                  target.stat().st_uid, target.stat().st_mtime_ns), state)
            self.assertEqual((stage / 'mcp-tokens').stat().st_mode & 0o777, 0o700)
        print('EVIDENCE: actual Hermes CLI twice; unrelated config, private bytes/modes/owner/mtime and other profile unchanged')


class GitOpsTests(unittest.TestCase):
    def test_startup_defaults_mounts_and_backup_contract(self):
        spec = yaml.safe_load((BASE / 'deployment.yaml').read_text())['spec']['template']['spec']
        initial = next(c for c in spec['initContainers'] if c['name'] == 'initial-config')
        self.assertIn("'/todoist/reconcile.py'", initial['command'][-1])
        self.assertIn({'name': 'todoist-config', 'mountPath': '/todoist', 'readOnly': True}, initial['volumeMounts'])
        self.assertIn({'name': 'todoist-config', 'configMap': {'name': 'onemind-hermes-todoist'}}, spec['volumes'])
        self.assertEqual(initial['resources']['requests'], {'cpu': '10m', 'memory': '32Mi'})
        self.assertEqual(initial['resources']['limits'], {'cpu': '500m', 'memory': '192Mi'})
        defaults = yaml.safe_load((BASE / 'config.defaults.yaml').read_text())
        self.assertEqual(defaults['mcp_servers']['todoist'], EXPECTED)
        kustomization = yaml.safe_load((BASE / 'kustomization.yaml').read_text())
        cm = next(c for c in kustomization['configMapGenerator'] if c['name'] == 'onemind-hermes-todoist')
        self.assertEqual(set(cm['files']), {'reconcile.py=mcp/todoist/reconcile.py', 'server.json=mcp/todoist/server.json'})
        self.assertIn('backup.yaml', kustomization['resources'])
        backup = yaml.safe_load((BASE / 'backup.yaml').read_text())
        self.assertEqual(backup['spec']['sourcePVC'], 'onemind-hermes-data')
        self.assertEqual(backup['spec']['restic']['repository'], 'onemind-hermes-restic')
        sealed = yaml.safe_load((BASE / 'backup.sealedsecret.yaml').read_text())
        self.assertIn('RESTIC_PASSWORD', sealed['spec']['encryptedData'])
        self.assertFalse(any(e['name'].startswith('TODOIST_') for e in spec['containers'][0]['env']))

    def test_fresh_and_existing_startup_twice_ownership_before_cli(self):
        # Execute the real init code with fixture-only path rewriting. Mock chown/root
        # detection (no root required), and unrelated Zoho/Whop subprocesses only.
        spec = yaml.safe_load((BASE / 'deployment.yaml').read_text())['spec']['template']['spec']
        initial = next(c for c in spec['initContainers'] if c['name'] == 'initial-config')
        self.assertIn("'/todoist/reconcile.py'", initial['command'][-1])
        cli = Path(os.environ.get('HERMES_TEST_CLI', '/opt/hermes/bin/hermes'))
        actual_run = subprocess.run
        for fresh in (True, False):
            with self.subTest(fresh=fresh):
                stage = stage_directory()
                home = stage / 'data'
                home.mkdir()
                cron = stage / 'cron'
                cron.mkdir()
                defaults = stage / 'defaults'
                defaults.mkdir()
                (defaults / 'config.yaml').write_bytes((BASE / 'config.defaults.yaml').read_bytes())
                if not fresh:
                    (home / 'config.yaml').write_text(yaml.safe_dump({'custom': {'retained': True}}))
                    (home / 'config.yaml').chmod(0o600)
                mapping = {'/opt/data': str(home), '/cron': str(cron), '/defaults': str(defaults),
                           '/todoist': str(HERE), '/zoho': str(BASE / 'mcp/zoho'), '/whop': str(BASE / 'mcp/whop')}
                class Paths(ast.NodeTransformer):
                    def visit_Constant(self, node):
                        if isinstance(node.value, str):
                            for prefix, replacement in mapping.items():
                                if node.value == prefix or node.value.startswith(prefix + '/'):
                                    node.value = replacement + node.value[len(prefix):]
                                    break
                        return node
                code = compile(ast.fix_missing_locations(Paths().visit(ast.parse(initial['command'][-1]))), '<isolated-init>', 'exec')
                for attempt in range(2):
                    ownership = []
                    calls = []
                    def run(argv, **kwargs):
                        calls.append(argv)
                        self.assertIn(home, ownership)
                        self.assertIn(home / 'config.yaml', ownership)
                        self.assertEqual((home / 'config.yaml').stat().st_mode & 0o777, 0o600)
                        if str(argv[1]) == str(HERE / 'reconcile.py') and cli.is_file():
                            return actual_run(argv + ['--cli', str(cli)], **kwargs)
                        return subprocess.CompletedProcess(argv, 0)
                    with patch('os.geteuid', return_value=0), patch('os.chown', side_effect=lambda p, *a: ownership.append(p)), \
                            patch('pwd.getpwnam', return_value=SimpleNamespace(pw_uid=os.getuid(), pw_gid=os.getgid())), \
                            patch('subprocess.run', side_effect=run):
                        exec(code, {})
                    self.assertEqual([Path(a[1]).parent.name for a in calls], ['zoho', 'whop', 'todoist'])
                    current = yaml.safe_load((home / 'config.yaml').read_text())
                    if cli.is_file():
                        self.assertEqual(current['mcp_servers']['todoist'], EXPECTED)
                    if not fresh:
                        self.assertEqual(current['custom'], {'retained': True})
                if cli.is_file():
                    print('EVIDENCE: ' + ('fresh' if fresh else 'existing') + ' init twice; real Hermes config CLI, simulated root/chown ordering')


if __name__ == '__main__':
    unittest.main()
