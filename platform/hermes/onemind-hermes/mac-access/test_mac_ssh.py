"""Offline contracts; all key material is generated in memory or scratch."""
import base64
import importlib.util
import os
from pathlib import Path
import stat
import struct
import subprocess
import tempfile
import unittest

HERE = Path(__file__).resolve().parent


class AccessTests(unittest.TestCase):
    def load(self):
        self.assertTrue((HERE / 'mac_ssh.py').exists(), 'secure SSH support missing')
        spec = importlib.util.spec_from_file_location('mac_ssh', HERE / 'mac_ssh.py')
        assert spec is not None and spec.loader is not None
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        return module

    def test_actual_openssh_configuration(self):
        module = self.load()
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / 'ssh-privatekey').touch(mode=0o400)
            (root / 'known_hosts').touch()
            (root / 'ssh-publickey').touch()
            command = module.ssh_command(root, 'zeus', ['true'])
            # Exercise OpenSSH's real option parser, not only string assertions.
            result = subprocess.run(command[:1] + ['-G'] + command[1:],
                                    capture_output=True, text=True, check=True)
            config = result.stdout
            for option in ['stricthostkeychecking true', 'batchmode yes',
                           'identitiesonly yes', 'identityagent none',
                           'forwardagent no', 'forwardx11 no', 'requesttty false',
                           'passwordauthentication no', 'kbdinteractiveauthentication no',
                           'globalknownhostsfile /dev/null', 'updatehostkeys false',
                           'hostkeyalgorithms ssh-ed25519']:
                self.assertIn(option, config)
            self.assertIn('userknownhostsfile ' + str(root / 'known_hosts'), config)
            self.assertEqual(command[-2:], ['zeus@omos-mbp.hermes.svc.cluster.local', 'true'])
            self.assertEqual(str(root / 'ssh-privatekey'), command[command.index('-i') + 1])
            self.assertEqual(sorted(p.name for p in root.iterdir()),
                             ['known_hosts', 'ssh-privatekey', 'ssh-publickey'])


    def test_cli_rejects_missing_mount_without_key_output(self):
        result = subprocess.run(['python', str(HERE / 'mac_ssh.py'), '--credentials',
                                 str(HERE / 'missing'), '--', 'true'], capture_output=True, text=True)
        self.assertNotEqual(result.returncode, 0)
        self.assertEqual(result.stdout, '')
        self.assertEqual(result.stderr, 'Mac SSH refused: check mount, permissions and account.\n')

    def test_active_mount_only_hermes_and_resources_unchanged(self):
        from ruamel.yaml import YAML
        import json
        yaml = YAML(typ='safe')
        path = HERE.parent
        current = yaml.load((path / 'deployment.yaml').read_text())
        base = yaml.load(subprocess.check_output(['git', 'show',
            '2a93f0cf381b1f8d32c31847430adb63b19cd72e:platform/hermes/onemind-hermes/deployment.yaml'], text=True))
        spec = current['spec']['template']['spec']
        old_spec = base['spec']['template']['spec']
        for group in ('containers', 'initContainers'):
            for container, old in zip(spec[group], old_spec[group]):
                # The hermes container's memory request was raised on 2026-11-07 (observed use ~1.5 GiB);
                # everything else in resources must stay as at the base commit.
                expected = json.loads(json.dumps(old['resources']))
                if container['name'] == 'hermes':
                    expected['requests']['memory'] = container['resources']['requests']['memory']
                self.assertEqual(container['resources'], expected)
                # The image may move with Hermes upgrades; it must stay digest-pinned and the same in every container.
                self.assertIn('@sha256:', container['image'])
                self.assertEqual(container['image'], spec['containers'][0]['image'])
                mounts = {m['name']: m for m in container.get('volumeMounts', [])}
                self.assertEqual('mac-ssh' in mounts, container['name'] in ('hermes', 'initial-config'))
                self.assertEqual('mac-ssh-input' in mounts, container['name'] == 'initial-config')
                if 'mac-ssh' in mounts:
                    self.assertEqual(mounts['mac-ssh']['mountPath'], '/var/run/legacy-mac-ssh')
                    self.assertEqual(mounts['mac-ssh'].get('readOnly', False), container['name'] == 'hermes')
                    self.assertIs(mounts['mac-ssh-client']['readOnly'], True)
                if 'mac-ssh-input' in mounts:
                    self.assertEqual(mounts['mac-ssh-input']['mountPath'], '/var/run/legacy-mac-ssh-input')
                    self.assertIs(mounts['mac-ssh-input']['readOnly'], True)
            # Init containers may be appended (install-langfuse-plugin); the originals keep their order and resources.
            if group == 'containers':
                self.assertEqual(len(spec[group]), len(old_spec[group]))
            else:
                self.assertGreaterEqual(len(spec[group]), len(old_spec[group]))
        init = spec['initContainers'][0]
        self.assertEqual(init['command'][0], '/opt/hermes/.venv/bin/python3')
        self.assertIn("sys.executable, '/etc/legacy-mac-ssh/prepare_mac_ssh.py'", init['command'][2])
        ram = next(v for v in spec['volumes'] if v['name'] == 'mac-ssh')
        self.assertEqual(ram, {'name': 'mac-ssh', 'emptyDir': {'medium': 'Memory', 'sizeLimit': '1Mi'}})
        from prepare_mac_ssh import SOURCE, DESTINATION
        self.assertEqual(str(SOURCE), '/var/run/legacy-mac-ssh-input')
        self.assertEqual(str(DESTINATION), '/var/run/legacy-mac-ssh')
        kustomization = yaml.load((path / 'kustomization.yaml').read_text())
        client = next(c for c in kustomization['configMapGenerator'] if c['name'] == 'onemind-hermes-mac-ssh-client')
        self.assertIn('prepare_mac_ssh.py=mac-access/prepare_mac_ssh.py', client['files'])
        volume = next(v for v in spec['volumes'] if v['name'] == 'mac-ssh-input')
        self.assertEqual(volume['secret']['defaultMode'], 0o400)
        self.assertIs(volume['secret']['optional'], True)
        sealed = json.loads((HERE / 'onemind-hermes-mac-ssh.sealedsecret.json').read_text())
        self.assertEqual(sealed['kind'], 'SealedSecret')
        self.assertEqual(sealed['metadata']['namespace'], 'hermes')
        self.assertEqual(sealed['metadata']['name'], volume['secret']['secretName'])
        self.assertEqual(set(sealed['spec']['encryptedData']), {'ssh-privatekey', 'ssh-publickey', 'known_hosts'})
        self.assertNotIn('data', sealed['spec']['template'])
        self.assertNotIn('stringData', sealed['spec']['template'])

    def test_fail_closed_permissions_symlinks_and_existing_unrestricted_key(self):
        module = self.load()
        public = (HERE / 'legacy-mac-access.pub').read_text()
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            host = root / 'host.pub'
            host.write_text(public)
            expected = module.fingerprint(public)
            ssh_dir = root / '.ssh'
            ssh_dir.mkdir(mode=0o700)
            target = ssh_dir / 'authorized_keys'
            target.write_text(public)
            target.chmod(0o600)
            before = target.read_bytes()
            with self.assertRaises(ValueError):
                module.enroll(public, expected, host, root)
            self.assertEqual(target.read_bytes(), before)
            target.chmod(0o644)
            with self.assertRaises(ValueError):
                module.enroll(public, expected, host, root)
            other_home = root / 'other'
            other_home.mkdir()
            (other_home / '.ssh').symlink_to(ssh_dir)
            with self.assertRaises(ValueError):
                module.enroll(public, expected, host, other_home)
            for name in ('ssh-privatekey', 'known_hosts', 'ssh-publickey'):
                (root / name).touch()
            (root / 'ssh-privatekey').chmod(0o600)
            with self.assertRaises(ValueError):
                module.ssh_command(root, 'zeus', ['true'])

    def test_observed_host_fingerprint_matches_real_openssh(self):
        module = self.load()
        observed = 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICDjSXg1sV7uPXb/nEhA1Lpf6mNtdG9omHyh2HwfVuZV'
        expected = 'SHA256:q8Al9q/LiMEPgwtN+VMoekF+8oB+TNYqjF3xDHzE7d8'
        self.assertEqual(module.fingerprint(observed), expected)
        result = subprocess.run(['ssh-keygen', '-lf', '-'], input=observed,
                                capture_output=True, text=True, check=True)
        self.assertIn(expected, result.stdout)

    def test_owner_enrollment_fingerprint_and_idempotence(self):
        module = self.load()
        self.assertTrue(hasattr(module, 'enroll'), 'owner enrollment missing')
        raw = struct.pack('>I', 11) + b'ssh-ed25519' + struct.pack('>I', 32) + os.urandom(32)
        public = 'ssh-ed25519 ' + base64.b64encode(raw).decode()
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            host = root / 'host.pub'
            host.write_text(public + '\n')
            authorized = root / '.ssh' / 'authorized_keys'
            authorized.parent.mkdir(mode=0o700)
            original = b'# preserve owner entries\nssh-rsa unrelated existing-owner-key'
            authorized.write_bytes(original)
            authorized.chmod(0o600)
            with self.assertRaises(ValueError):
                module.enroll(public, 'SHA256:WRONG', host, root)
            self.assertEqual(authorized.read_bytes(), original)
            expected = module.fingerprint(public)
            self.assertEqual(module.enroll(public, expected, host, root), 'enrolled')
            first = authorized.read_bytes()
            self.assertTrue(first.startswith(original + b'\n'))
            self.assertIn(b'restrict,no-agent-forwarding,no-port-forwarding,no-pty,no-user-rc,no-X11-forwarding ', first)
            self.assertEqual(module.enroll(public, expected, host, root), 'already enrolled')
            self.assertEqual(authorized.read_bytes(), first)
            self.assertEqual(stat.S_IMODE(authorized.stat().st_mode), 0o600)


if __name__ == '__main__':
    unittest.main()
