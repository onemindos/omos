"""Run Mac SSH contracts on either nonroot developers or root CI runners."""
import importlib.util
import os
from pathlib import Path
import pwd
import subprocess
import sys
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
DIRECTORY = ROOT / 'platform/hermes/onemind-hermes/mac-access'
CI_UID = 10000
CI_GID = 10000


def stage_directory():
    parent = Path(os.environ.get('RUNNER_TEMP') or os.environ.get('TMPDIR')
                  or Path.home() / '.cache/hermes/scratch')
    parent.mkdir(parents=True, exist_ok=True)
    return Path(tempfile.mkdtemp(prefix='mac-ssh-ci-contract-', dir=parent))


class RootRunnerContracts(unittest.TestCase):
    """Do not fake getuid or weaken the production nonroot-only requirement."""

    def test_existing_contracts_execute_as_actual_nonroot_identity(self):
        # OpenSSH also requires a passwd entry for its invoking UID. Use a real
        # unprivileged runner account; the separate root-transfer test below
        # exercises the production target UID/GID 10000 even if absent from NSS.
        owner = next((p for p in pwd.getpwall() if p.pw_uid > 0 and p.pw_gid > 0
                      and p.pw_name in ('hermes', 'runner', 'nobody')), None)
        self.assertIsNotNone(owner, 'a registered nonroot account is required')
        assert owner is not None
        stage = stage_directory()
        os.chown(stage, owner.pw_uid, owner.pw_gid)
        os.chmod(stage, 0o700)
        env = dict(os.environ, TMPDIR=str(stage), HOME=str(stage),
                   GIT_CONFIG_COUNT='1', GIT_CONFIG_KEY_0='safe.directory',
                   GIT_CONFIG_VALUE_0=str(ROOT))
        result = subprocess.run(
            [sys.executable, '-m', 'unittest', 'discover', '-s', str(DIRECTORY), '-v'],
            cwd=ROOT, env=env, user=owner.pw_uid, group=owner.pw_gid, extra_groups=[],
            capture_output=True, text=True, timeout=90,
        )
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn('Ran 15 tests', result.stderr)
        self.assertIn('OK', result.stderr)
        print('EVIDENCE: all 15 Mac SSH contracts ran as nonroot UID/GID '
              + str(owner.pw_uid) + '/' + str(owner.pw_gid))

    def test_real_root_preparation_transfers_ownership_for_nonroot_readback(self):
        stage = stage_directory()
        os.chmod(stage, 0o755)
        source, target = stage / 'input', stage / 'ram'
        source.mkdir(mode=0o700)
        target.mkdir()
        target.chmod(0o1777)
        spec = importlib.util.spec_from_file_location('mac_ssh_root_prepare', DIRECTORY / 'prepare_mac_ssh.py')
        assert spec is not None and spec.loader is not None
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        for name in module.FILES:
            (source / name).write_bytes(b'synthetic-contract-fixture')
            (source / name).chmod(0o400)
        self.assertTrue(module.prepare(source, target, CI_UID, CI_GID))
        self.assertTrue(module.prepare(source, target, CI_UID, CI_GID))
        code = (
            'import os,stat,sys;from pathlib import Path;'
            'assert os.getuid()==10000 and os.getgid()==10000;'
            'root=Path(sys.argv[1]);assert stat.S_IMODE(root.stat().st_mode)==0o700;'
            'files=[root/n for n in ("ssh-privatekey","ssh-publickey","known_hosts")];'
            'assert all(p.stat().st_uid==10000 and p.stat().st_gid==10000 '
            'and stat.S_IMODE(p.stat().st_mode)==0o400 '
            'and p.read_bytes()==b"synthetic-contract-fixture" for p in files);'
            'print("root-init to nonroot readback verified")'
        )
        result = subprocess.run(
            [sys.executable, '-c', code, str(target)], user=CI_UID, group=CI_GID,
            extra_groups=[], capture_output=True, text=True, timeout=15,
        )
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertEqual(result.stdout.strip(), 'root-init to nonroot readback verified')


def load_tests(loader, tests, pattern):
    if os.geteuid() == 0:
        return loader.loadTestsFromTestCase(RootRunnerContracts)
    location = str(DIRECTORY)
    if location not in sys.path:
        sys.path.insert(0, location)
    suite = unittest.TestSuite()
    for name in ('test_mac_ssh', 'test_prepare_mac_ssh'):
        spec = importlib.util.spec_from_file_location('mac_ssh_contracts_' + name, DIRECTORY / (name + '.py'))
        assert spec is not None and spec.loader is not None
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        suite.addTests(loader.loadTestsFromModule(module))
    return suite
