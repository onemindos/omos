"""RAM credential preparation contracts; synthetic bytes only."""
import importlib.util
import os
from pathlib import Path
import stat
import tempfile
import unittest

HERE = Path(__file__).resolve().parent


class PreparationTests(unittest.TestCase):
    def load(self):
        self.assertTrue((HERE / 'prepare_mac_ssh.py').exists(), 'RAM preparation missing')
        spec = importlib.util.spec_from_file_location('prepare', HERE / 'prepare_mac_ssh.py')
        assert spec is not None and spec.loader is not None
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        return module

    def test_projected_secret_becomes_owner_only_readable_nonroot_files(self):
        module = self.load()
        self.assertNotEqual(os.getuid(), 0, 'exercise actual nonroot reader')
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            source = root / 'input'
            source.mkdir()
            version = source / '..version'
            version.mkdir()
            (source / '..data').symlink_to(version.name)
            for name in module.FILES:
                (version / name).write_bytes(b'synthetic fixture')
                (version / name).chmod(0o400)
                (source / name).symlink_to('..data/' + name)
            destination = root / 'ram'
            destination.mkdir()
            self.assertTrue(module.prepare(source, destination, os.getuid(), os.getgid(), trusted_uid=os.getuid()))
            self.assertEqual(stat.S_IMODE(destination.stat().st_mode), 0o700)
            self.assertEqual(set(p.name for p in destination.iterdir()), set(module.FILES))
            for name in module.FILES:
                item = destination / name
                self.assertEqual(item.read_bytes(), b'synthetic fixture')
                self.assertEqual((item.stat().st_uid, item.stat().st_gid), (os.getuid(), os.getgid()))
                self.assertEqual(stat.S_IMODE(item.stat().st_mode), 0o400)
            from mac_ssh import ssh_command
            command = ssh_command(destination, 'zeus', ['true'])
            import subprocess
            result = subprocess.run(command[:1] + ['-G'] + command[1:], capture_output=True, text=True)
            self.assertEqual(result.returncode, 0)
            self.assertIn(str(destination / 'ssh-privatekey'), result.stdout)
            before = {name: (destination / name).read_bytes() for name in module.FILES}
            self.assertTrue(module.prepare(source, destination, os.getuid(), os.getgid(), trusted_uid=os.getuid()))
            self.assertEqual(before, {name: (destination / name).read_bytes() for name in module.FILES})

    def test_empty_optional_secret_is_unavailable_without_writes(self):
        module = self.load()
        with tempfile.TemporaryDirectory() as directory:
            source, target = Path(directory) / 'input', Path(directory) / 'ram'
            source.mkdir()
            target.mkdir()
            self.assertFalse(module.prepare(source, target, os.getuid(), os.getgid(), trusted_uid=os.getuid()))
            self.assertEqual(list(target.iterdir()), [])
        from unittest.mock import patch
        from types import SimpleNamespace
        import contextlib
        import io
        output = io.StringIO()
        with patch.object(module.os, 'geteuid', return_value=0), patch.object(module.pwd, 'getpwnam', return_value=SimpleNamespace(pw_uid=1234, pw_gid=5678)), patch.object(module, 'prepare', return_value=False) as prepare, contextlib.redirect_stdout(output):
            self.assertEqual(module.main(), 0)
            prepare.assert_called_once_with(module.SOURCE, module.DESTINATION, 1234, 5678)
        self.assertEqual(output.getvalue(), 'Mac SSH unavailable: optional Secret absent.\n')

    def test_rejects_partial_escape_bad_mode_owner_and_root_target(self):
        module = self.load()
        for defect in ('partial', 'escape', 'mode', 'owner', 'root', 'directory-symlink', 'hardlink'):
            with self.subTest(defect=defect), tempfile.TemporaryDirectory() as directory:
                root = Path(directory)
                source, target = root / 'input', root / 'ram'
                source.mkdir()
                target.mkdir()
                for name in module.FILES:
                    if defect == 'partial' and name == 'known_hosts':
                        continue
                    item = source / name
                    if name == 'ssh-privatekey' and defect in ('escape', 'hardlink'):
                        outside = root / 'outside'
                        outside.write_bytes(b'synthetic')
                        outside.chmod(0o400)
                        if defect == 'escape':
                            item.symlink_to(outside)
                        else:
                            os.link(outside, item)
                    else:
                        item.write_bytes(b'synthetic')
                        item.chmod(0o444 if defect == 'mode' else 0o400)
                if defect == 'directory-symlink':
                    link = root / 'link'
                    link.symlink_to(source)
                    source = link
                with self.assertRaises((ValueError, OSError)):
                    module.prepare(source, target, 0 if defect == 'root' else os.getuid(), os.getgid(),
                                   trusted_uid=os.getuid() + 1 if defect == 'owner' else os.getuid())
                self.assertEqual(list(target.iterdir()), [])

    def test_absent_optional_secret_with_empty_kubernetes_projection(self):
        module = self.load()
        with tempfile.TemporaryDirectory() as directory:
            source, target = Path(directory) / 'input', Path(directory) / 'ram'
            source.mkdir()
            target.mkdir()
            version = source / '..2026_11_05_01_00_00'
            version.mkdir()
            (source / '..data').symlink_to(version.name)
            self.assertFalse(module.prepare(source, target, os.getuid(), os.getgid(), trusted_uid=os.getuid()))
            self.assertEqual(list(target.iterdir()), [])
            # A real partial projection must never be mistaken for absence.
            (version / 'ssh-privatekey').write_bytes(b'synthetic fixture')
            (version / 'ssh-privatekey').chmod(0o400)
            (source / 'ssh-privatekey').symlink_to('..data/ssh-privatekey')
            with self.assertRaises(ValueError):
                module.prepare(source, target, os.getuid(), os.getgid(), trusted_uid=os.getuid())
            self.assertEqual(list(target.iterdir()), [])

    def test_fresh_world_writable_tmpfs_is_tightened_before_payload(self):
        module = self.load()
        with tempfile.TemporaryDirectory() as directory:
            source, target = Path(directory) / 'input', Path(directory) / 'ram'
            source.mkdir()
            target.mkdir()
            target.chmod(0o1777)
            for name in module.FILES:
                (source / name).write_bytes(b'synthetic fixture')
                (source / name).chmod(0o400)
            self.assertTrue(module.prepare(source, target, os.getuid(), os.getgid(), trusted_uid=os.getuid()))
            self.assertEqual(stat.S_IMODE(target.stat().st_mode), 0o700)
            self.assertTrue(all((target / name).read_bytes() == b'synthetic fixture' for name in module.FILES))

    def test_readonly_projected_volume_root_can_have_writable_mode_bits(self):
        module = self.load()
        from unittest.mock import patch
        from types import SimpleNamespace
        with tempfile.TemporaryDirectory() as directory:
            source, target = Path(directory) / 'input', Path(directory) / 'ram'
            source.mkdir()
            source.chmod(0o1777)
            target.mkdir()
            for name in module.FILES:
                (source / name).write_bytes(b'synthetic fixture')
                (source / name).chmod(0o400)
            with patch.object(module.os, 'statvfs', return_value=SimpleNamespace(f_flag=os.ST_RDONLY)):
                self.assertTrue(module.prepare(source, target, os.getuid(), os.getgid(), trusted_uid=os.getuid()))

    def test_restart_resumes_approved_staging_without_truncation_or_deletion(self):
        module = self.load()
        for prefix in (b'', b'synthetic', b'synthetic fixture'):
            with self.subTest(prefix_length=len(prefix)), tempfile.TemporaryDirectory() as directory:
                source, target = Path(directory) / 'input', Path(directory) / 'ram'
                source.mkdir()
                target.mkdir(mode=0o700)
                for name in module.FILES:
                    (source / name).write_bytes(b'synthetic fixture')
                    (source / name).chmod(0o400)
                staged = target / '.ssh-privatekey.new'
                staged.write_bytes(prefix)
                staged.chmod(0o600)
                self.assertTrue(module.prepare(source, target, os.getuid(), os.getgid(), trusted_uid=os.getuid()))
                self.assertEqual((target / 'ssh-privatekey').read_bytes(), b'synthetic fixture')
                before = {name: (target / name).read_bytes() for name in module.FILES}
                self.assertTrue(module.prepare(source, target, os.getuid(), os.getgid(), trusted_uid=os.getuid()))
                self.assertEqual(before, {name: (target / name).read_bytes() for name in module.FILES})

    def test_restart_rejects_different_credentials_without_replacing_existing_bytes(self):
        module = self.load()
        with tempfile.TemporaryDirectory() as directory:
            source, target = Path(directory) / 'input', Path(directory) / 'ram'
            source.mkdir()
            target.mkdir(mode=0o700)
            for name in module.FILES:
                (source / name).write_bytes(b'synthetic fixture')
                (source / name).chmod(0o400)
            unknown = target / '.ssh-privatekey.new'
            unknown.write_bytes(b'unrelated bytes')
            unknown.chmod(0o600)
            with self.assertRaises(ValueError):
                module.prepare(source, target, os.getuid(), os.getgid(), trusted_uid=os.getuid())
            self.assertEqual(unknown.read_bytes(), b'unrelated bytes')

    def test_main_refuses_nonroot_or_root_account_with_sanitized_error(self):
        module = self.load()
        from unittest.mock import patch
        from types import SimpleNamespace
        import contextlib
        import io
        for effective, owner in ((10000, 10000), (0, 0)):
            error = io.StringIO()
            with patch.object(module.os, 'geteuid', return_value=effective), patch.object(module.pwd, 'getpwnam', return_value=SimpleNamespace(pw_uid=owner, pw_gid=owner)), contextlib.redirect_stderr(error):
                self.assertEqual(module.main(), 1)
            self.assertEqual(error.getvalue(), 'Mac SSH preparation refused: check account and volume safety.\n')


if __name__ == '__main__':
    unittest.main()
