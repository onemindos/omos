"""Contracts for the Hermes Desktop fixes in platform/hermes/onemind-hermes/desktop and nats/gateway-deps.txt."""
import importlib.util
import os
import re
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / 'platform/hermes/onemind-hermes'
DEPS = BASE / 'nats/gateway-deps.txt'
APPLY_PY = BASE / 'desktop/apply_desktop.py'
APPLY_SH = BASE / 'desktop/apply.sh'
CHECK_PY = BASE / 'desktop/check_desktop.py'


def load(path, name):
    spec = importlib.util.spec_from_file_location(name, path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


class GatewayDepsTests(unittest.TestCase):
    def requirement_lines(self):
        return [l for l in DEPS.read_text().splitlines() if l.strip() and not l.startswith('#')]

    def test_every_requirement_is_pinned_and_hashed(self):
        for line in self.requirement_lines():
            self.assertRegex(line, r'^[A-Za-z0-9_.\-]+==[0-9][^ ]* --hash=sha256:[0-9a-f]{64}$', line)

    def test_edge_tts_for_desktop_speech_is_present(self):
        names = {l.split('==')[0] for l in self.requirement_lines()}
        self.assertIn('edge-tts', names)
        self.assertIn('tabulate', names)  # edge-tts needs it; the image's venv does not ship it


class ApplyDesktopTests(unittest.TestCase):
    def setUp(self):
        self.mod = load(APPLY_PY, 'apply_desktop')

    def test_junk_names(self):
        self.assertTrue(self.mod.is_junk('._SKILL.md'))
        self.assertTrue(self.mod.is_junk('.DS_Store'))
        self.assertFalse(self.mod.is_junk('SKILL.md'))
        self.assertFalse(self.mod.is_junk('.hidden-skill'))

    def test_scan_finds_appledouble_and_wrong_owner_without_following_links(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp) / 'skills'
            (root / 'devops/x').mkdir(parents=True)
            (root / 'devops/x/SKILL.md').write_text('# x\n')
            (root / 'devops/x/._SKILL.md').write_text('junk')
            (root / '._devops').write_text('junk')
            outside = Path(tmp) / 'outside'
            outside.mkdir()
            (outside / '._keep').write_text('not ours')
            os.symlink(outside, root / 'link')
            me, my_group = os.getuid(), os.getgid()
            junk, wrong = self.mod.scan(str(root), me, my_group, recurse=True)
            self.assertEqual(sorted(os.path.basename(p) for p in junk), ['._SKILL.md', '._devops'])
            self.assertEqual(wrong, [])  # everything is owned by the current user
            junk, wrong = self.mod.scan(str(root), me + 1, my_group, recurse=True)
            self.assertIn(str(root / 'devops/x/SKILL.md'), wrong)
            self.assertNotIn(str(outside / '._keep'), junk)

    def test_non_recursive_scan_stays_in_the_top_directory(self):
        with tempfile.TemporaryDirectory() as tmp:
            Path(tmp, '._SOUL.md').write_text('junk')
            Path(tmp, 'sub').mkdir()
            Path(tmp, 'sub/._deep').write_text('junk')
            junk, _ = self.mod.scan(tmp, os.getuid(), os.getgid(), recurse=False)
            self.assertEqual([os.path.basename(p) for p in junk], ['._SOUL.md'])


class ScriptSafetyTests(unittest.TestCase):
    def test_apply_script_has_dry_run_and_check_and_never_touches_secrets(self):
        sh = APPLY_SH.read_text()
        self.assertIn('--check', sh)
        self.assertIn('KUBE_CONTEXT', sh)
        for text in (sh, APPLY_PY.read_text()):
            self.assertNotRegex(text, r'open\([^)]*(auth\.json|\.env)', 'must not read credential files')
        self.assertIn('--dry-run', APPLY_PY.read_text())

    def test_check_script_prints_no_values(self):
        src = CHECK_PY.read_text()
        self.assertNotRegex(src, r'print\([^)]*\.text\b(?!\[:)', 'response bodies must be truncated and redacted')
        self.assertIn('TOKENISH.sub', src)
        self.assertNotIn('os.environ[', src)


if __name__ == '__main__':
    unittest.main()
