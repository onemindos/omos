"""Synthetic-only regressions against actual Hermes filters, not reimplementations.

Run in a NEW interpreter with Hermes dependencies on PYTHONPATH. --policy-root
may select a patched scratch copy; all other modules come from --runtime-root.
No dotenv loader, CLI authentication, session DB or network calls are used.
"""
import argparse
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

BOOT = 'OP_SERVICE_ACCOUNT_TOKEN'
APP = 'CLICKUP_API_TOKEN'
# Deliberately not shaped like real credentials.
MARKER = 'synthetic-bootstrap-marker'
APP_A = 'synthetic-app-profile-a'
APP_B = 'synthetic-app-profile-b'


class ChildPolicyTests(unittest.TestCase):
    def setUp(self):
        passthrough.clear_env_passthrough()
        passthrough._config_passthrough.clear()
        (HOME / 'config.yaml').write_text(json.dumps({
            'terminal': {'env_passthrough': [APP]}}))

    def assert_child(self, env, expected_app=None):
        self.assertFalse(any(k.upper() == BOOT for k in env))
        self.assertNotIn(MARKER, env.values())
        code = ('import os; assert "OP_SERVICE_ACCOUNT_TOKEN" not in os.environ; '
                'assert "synthetic-bootstrap-marker" not in os.environ.values()')
        if expected_app is not None:
            self.assertEqual(env.get(APP), expected_app)
            code += '; assert os.environ.get("CLICKUP_API_TOKEN") == ' + repr(expected_app)
        subprocess.run([sys.executable, '-c', code], env=env, check=True,
                       capture_output=True, timeout=10)

    def test_foreground_terminal_and_snapshot_extras(self):
        os.environ[BOOT] = MARKER
        try:
            for extra in ({}, {BOOT: MARKER}, {'_HERMES_FORCE_' + BOOT: MARKER}):
                self.assert_child(local._make_run_env(extra))
        finally:
            os.environ.pop(BOOT)

    def test_patch_rejects_unknown_source_without_writing(self):
        from apply_source_patch import apply, PATCHES
        with tempfile.TemporaryDirectory(dir=HOME) as directory:
            root = Path(directory)
            # First file is valid; second fails. Neither may be modified.
            before = {}
            for index, (relative, _, _, _) in enumerate(PATCHES):
                target = root / relative
                target.parent.mkdir(parents=True, exist_ok=True)
                content = ((args.runtime_root / relative).read_bytes() if index == 0
                           else b'unrecognized-source-marker\n')
                target.write_bytes(content)
                before[target] = content
            with self.assertRaises(ValueError):
                apply(root)
            for target, content in before.items():
                self.assertEqual(target.read_bytes(), content)

    def test_terminal_without_bootstrap_passthrough(self):
        self.assert_child(local._sanitize_subprocess_env({BOOT: MARKER, APP: APP_A}), APP_A)

    def test_actual_filter_casefold_and_force(self):
        for name in (BOOT, BOOT.lower(), '_HERMES_FORCE_' + BOOT):
            for unwrap in (False, True):
                with self.subTest(name=name, unwrap=unwrap):
                    out = {}
                    local._filter_secret_env({name: MARKER}, out, unwrap_force=unwrap)
                    self.assert_child(out)

    def test_config_and_skill_cannot_reallow_bootstrap(self):
        (HOME / 'config.yaml').write_text(json.dumps({
            'terminal': {'env_passthrough': [BOOT, BOOT.lower(), APP]}}))
        passthrough.register_env_passthrough([BOOT, BOOT.lower(), APP])
        self.assertFalse(passthrough.is_env_passthrough(BOOT))
        self.assertFalse(passthrough.is_env_passthrough(BOOT.lower()))
        self.assertNotIn(BOOT, passthrough.get_all_passthrough())
        source = {BOOT: MARKER, BOOT.lower(): MARKER, APP: APP_A}
        self.assert_child(local._sanitize_subprocess_env(source), APP_A)
        self.assert_child(code_env._scrub_child_env(source), APP_A)

    def test_execute_code_interpreter_probe_has_no_bootstrap(self):
        os.environ[BOOT] = MARKER
        try:
            probe = code_env._probe_python(sys.executable,
                'import os; print("OP_SERVICE_ACCOUNT_TOKEN" in os.environ)', text=True)
            assert probe is not None
            self.assertEqual(probe.returncode, 0)
            self.assertEqual(probe.stdout.strip(), 'False')
            self.assertTrue(code_env._is_usable_python(sys.executable))
            self.assertTrue(code_env._python_environment_prefix(sys.executable))
        finally:
            os.environ.pop(BOOT)

    def test_execute_code_default_and_app_opt_in(self):
        self.assert_child(code_env._scrub_child_env({BOOT: MARKER, APP: APP_A}), APP_A)
        (HOME / 'config.yaml').write_text('{}')
        passthrough._config_passthrough.clear()
        env = code_env._scrub_child_env({BOOT: MARKER, APP: APP_A})
        self.assertNotIn(APP, env)
        self.assert_child(env)

    def test_nonterminal_even_inherit_credentials(self):
        for inherit in (False, True):
            self.assert_child(local.hermes_subprocess_env(
                base_env={BOOT: MARKER, BOOT.lower(): MARKER},
                inherit_credentials=inherit))

    def test_scoped_app_never_falls_back_to_other_profile(self):
        multiplex = scope.set_multiplex_context(True)
        try:
            for values, expected in (({APP: APP_A, BOOT: MARKER}, APP_A),
                                     ({APP: APP_B, BOOT: MARKER}, APP_B),
                                     ({BOOT: MARKER}, None)):
                token = scope.set_secret_scope(values)
                try:
                    for env in (local._sanitize_subprocess_env({APP: APP_A, BOOT: MARKER}),
                                code_env._scrub_child_env({APP: APP_A, BOOT: MARKER}),
                                local._sanitize_subprocess_env({}),
                                code_env._scrub_child_env({})):
                        self.assertEqual(env.get(APP), expected)
                        self.assert_child(env, expected)
                finally:
                    scope.reset_secret_scope(token)
        finally:
            scope.reset_multiplex_context(multiplex)

    def test_native_source_retains_only_explicit_bootstrap(self):
        # The native source uses its source-environment context, not app passthrough.
        os.environ[BOOT] = MARKER
        os.environ[APP] = APP_A
        try:
            token_value = native.get_source_environment()[BOOT]
            native_env = native._op_child_env(token_value)
            self.assertEqual(native_env[BOOT], MARKER)
            self.assertNotIn(APP, native_env)
            self.assert_child(local._sanitize_subprocess_env(dict(os.environ)), APP_A)
        finally:
            os.environ.pop(BOOT)
            os.environ.pop(APP)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--runtime-root', type=Path, required=True)
    parser.add_argument('--policy-root', type=Path)
    args, remaining = parser.parse_known_args()
    # Clear inherited env BEFORE importing Hermes; only public source is imported.
    scratch = os.environ.get('TMPDIR')
    with tempfile.TemporaryDirectory(prefix='hermes-policy-', dir=scratch) as directory:
        HOME = Path(directory)
        os.environ.clear()
        sys.dont_write_bytecode = True
        os.environ.update(HOME=directory, HERMES_HOME=directory, TMPDIR=directory,
                          PATH=os.defpath, PYTHONDONTWRITEBYTECODE='1')
        sys.path.insert(0, str(args.runtime_root))
        if args.policy_root:
            for name, relative in (
                ('tools.environments.local_env_policy', 'tools/environments/local_env_policy.py'),
                ('tools.code_execution_env', 'tools/code_execution_env.py'),
            ):
                spec = importlib.util.spec_from_file_location(name, args.policy_root / relative)
                assert spec is not None and spec.loader is not None
                module = importlib.util.module_from_spec(spec)
                sys.modules[name] = module
                spec.loader.exec_module(module)
        from tools.environments import local
        from tools import code_execution_env as code_env, env_passthrough as passthrough
        from agent import secret_scope as scope
        from agent.secret_sources import onepassword as native
        unittest.main(argv=[sys.argv[0]] + remaining)
