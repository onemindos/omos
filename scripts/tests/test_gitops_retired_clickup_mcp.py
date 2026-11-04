"""ClickUp stays retired on fresh seeds, PVC recovery and manual MCP reconciliation."""
import ast
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch

import yaml

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / 'platform/hermes/onemind-hermes'
CLI = Path(os.environ.get('HERMES_TEST_CLI', '/opt/hermes/bin/hermes'))
SCRATCH = Path(os.environ.get('TMPDIR') or os.environ.get('RUNNER_TEMP')
               or Path.home() / '.cache/hermes/scratch')
SCRATCH.mkdir(parents=True, exist_ok=True)


class ClickUpRetirementTests(unittest.TestCase):
    def test_defaults_do_not_restore_clickup(self):
        servers = yaml.safe_load((BASE / 'config.defaults.yaml').read_text())['mcp_servers']
        self.assertNotIn('clickup', servers)
        for name in ('todoist', 'dropbox', 'whop', 'ops', 'data', 'tak'):
            self.assertIn(name, servers)

    def test_startup_and_manual_removal_are_declared(self):
        deployment = yaml.safe_load((BASE / 'deployment.yaml').read_text())
        code = deployment['spec']['template']['spec']['initContainers'][0]['command'][-1]
        tree = ast.parse(code)
        calls = [node for node in ast.walk(tree) if isinstance(node, ast.List)]
        self.assertTrue(any([getattr(item, 'value', None) for item in node.elts] ==
                            ['/opt/hermes/bin/hermes', 'config', 'unset', 'mcp_servers.clickup']
                            for node in calls))
        applier = ast.parse((BASE / 'mcp/apply_mcp.py').read_text())
        removal = next(node.value for node in applier.body if isinstance(node, ast.Assign)
                       and any(isinstance(t, ast.Name) and t.id == 'REMOVE' for t in node.targets))
        self.assertIn('clickup', ast.literal_eval(removal))

    def test_actual_startup_removes_restored_clickup_twice(self):
        if not CLI.is_file():
            self.skipTest('Actual Hermes CLI unavailable')
        deployment = yaml.safe_load((BASE / 'deployment.yaml').read_text())
        initial = deployment['spec']['template']['spec']['initContainers'][0]
        original = {'model': {'default': 'unchanged'}, 'custom': {'retained': True},
                    'mcp_servers': {name: {'url': 'https://example.invalid/' + name}
                                    for name in ('todoist', 'dropbox', 'whop', 'other')}}
        for fresh in (True, False):
            with self.subTest(fresh=fresh), tempfile.TemporaryDirectory(dir=SCRATCH) as tmp:
                home = Path(tmp) / 'data'
                home.mkdir()
                defaults = Path(tmp) / 'defaults'
                defaults.mkdir()
                (defaults / 'config.yaml').write_bytes((BASE / 'config.defaults.yaml').read_bytes())
                cron = Path(tmp) / 'cron'
                cron.mkdir()
                if not fresh:
                    restored = dict(original, mcp_servers=dict(original['mcp_servers'], clickup={'url': 'https://mcp.clickup.com/mcp'}))
                    (home / 'config.yaml').write_text(yaml.safe_dump(restored))
                    (home / 'config.yaml').chmod(0o600)
                private = {'auth.json': b'opaque-provider-fixture', '.env': b'# unchanged\n',
                           'mcp-tokens/dropbox.json': b'opaque-oauth-fixture',
                           'profiles/other/config.yaml': b'untouched-profile'}
                for name, content in private.items():
                    path = home / name
                    path.parent.mkdir(parents=True, exist_ok=True)
                    path.write_bytes(content)
                    path.chmod(0o600)
                before = {name: ((home / name).read_bytes(), (home / name).stat()) for name in private}
                mapping = {'/opt/data': str(home), '/defaults': str(defaults), '/cron': str(cron),
                           '/opt/hermes/bin/hermes': str(CLI)}
                class Paths(ast.NodeTransformer):
                    def visit_Constant(self, node):
                        if isinstance(node.value, str):
                            for prefix, replacement in mapping.items():
                                if node.value == prefix or node.value.startswith(prefix + '/'):
                                    node.value = replacement + node.value[len(prefix):]
                                    break
                        return node
                code = compile(ast.fix_missing_locations(Paths().visit(ast.parse(initial['command'][-1]))), '<fixture-init>', 'exec')
                actual_run = subprocess.run
                for attempt in range(2):
                    calls = []
                    def run(argv, **kwargs):
                        if argv[:4] == [str(CLI), 'config', 'unset', 'mcp_servers.clickup']:
                            calls.append(argv)
                            return actual_run(argv, **kwargs)
                        # Existing Zoho/Whop/Todoist reconciliation is covered by its own suites.
                        return subprocess.CompletedProcess(argv, 0)
                    with patch('os.chown'), patch('subprocess.run', side_effect=run):
                        exec(code, {})
                    self.assertEqual(len(calls), int(not fresh and attempt == 0),
                                     'Startup must retire ClickUp outside the seed guard')
                    updated = yaml.safe_load((home / 'config.yaml').read_text())
                    self.assertNotIn('clickup', updated.get('mcp_servers', {}))
                    if not fresh:
                        self.assertEqual(updated, original)
                    self.assertEqual((home / 'config.yaml').stat().st_mode & 0o777, 0o600)
                    for name, (content, stat) in before.items():
                        path = home / name
                        self.assertEqual(path.read_bytes(), content)
                        current = path.stat()
                        self.assertEqual((current.st_mode, current.st_uid, current.st_mtime_ns),
                                         (stat.st_mode, stat.st_uid, stat.st_mtime_ns))

    def test_manual_applier_removes_clickup_and_is_idempotent(self):
        try:
            import ruamel.yaml  # noqa: F401 - runtime dependency of the existing applier
        except ImportError:
            self.skipTest('Hermes runtime ruamel.yaml unavailable')
        with tempfile.TemporaryDirectory(dir=SCRATCH) as tmp:
            config = Path(tmp) / 'config.yaml'
            config.write_text(yaml.safe_dump({'custom': True, 'mcp_servers': {
                'clickup': {'url': 'https://mcp.clickup.com/mcp'},
                'other': {'url': 'https://example.invalid/mcp'}}}))
            command = [sys.executable, str(BASE / 'mcp/apply_mcp.py'),
                       '--defaults', str(BASE / 'config.defaults.yaml'), '--config', str(config)]
            before = config.read_bytes()
            subprocess.run(command + ['--dry-run'], capture_output=True, check=True)
            self.assertEqual(config.read_bytes(), before)
            subprocess.run(command, capture_output=True, check=True)
            updated = yaml.safe_load(config.read_text())
            self.assertNotIn('clickup', updated['mcp_servers'])
            self.assertEqual(updated['mcp_servers']['other'], {'url': 'https://example.invalid/mcp'})
            self.assertTrue(updated['custom'])
            before = config.read_bytes(), config.stat().st_mtime_ns
            subprocess.run(command, capture_output=True, check=True)
            self.assertEqual((config.read_bytes(), config.stat().st_mtime_ns), before)


if __name__ == '__main__':
    unittest.main()
