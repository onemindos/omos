#!/usr/bin/env python3
"""Isolated contract tests; fake Docker never contacts a daemon or registry."""
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile
import unittest


MOCK_DOCKER = r'''#!/usr/bin/env python3
import json, os, pathlib, sys
args = sys.argv[1:]
with open(os.environ['RELEASE_TEST_LOG'], 'a') as f:
    f.write(json.dumps(args)+'\n')
if args[:2] == ['buildx', 'build']:
    root = pathlib.Path(args[-1])
    assert (root/'package.json').exists()
    assert not (root/'.env').exists(), 'untracked credentials reached build context'
    assert not (root/'node_modules').exists(), 'host binaries reached build context'
    if os.getenv('FAIL_BUILD'): sys.exit(12)
    dest = pathlib.Path(args[args.index('--metadata-file')+1])
    dest.write_text(json.dumps({'containerimage.digest': 'sha256:'+'a'*64}))
if args[:3] == ['buildx', 'imagetools', 'inspect'] and os.getenv('FAIL_INSPECT'):
    sys.exit(13)
'''


class ReleaseTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.project = self.root / 'services/onemind-portal'
        (self.project / 'scripts').mkdir(parents=True)
        self.script = self.project / 'scripts/release-image.sh'
        shutil.copy2(Path(__file__).with_name('release-image.sh'), self.script)
        (self.project / 'package.json').write_text('{}\n')
        self.git('init', '-q')
        self.git('config', 'user.email', 'release-test@example.invalid')
        self.git('config', 'user.name', 'Release Test')
        self.git('add', '.')
        self.git('commit', '-qm', 'fixture')
        # Ignore files as a normal working copy would; archive must still omit them.
        (self.root / '.git/info/exclude').write_text('.env\nnode_modules/\n')
        (self.project / '.env').write_text('TEST_FIXTURE_ONLY=not-a-credential\n')
        (self.project / 'node_modules').mkdir()
        self.bin = self.root / 'bin'
        self.bin.mkdir()
        fake = self.bin / 'docker'
        fake.write_text(MOCK_DOCKER)
        fake.chmod(0o755)
        self.log = self.root / 'docker.log'
        self.output = self.root / 'output'
        self.env = dict(os.environ, PATH=str(self.bin)+os.pathsep+os.environ['PATH'],
                        RELEASE_TEST_LOG=str(self.log))

    def git(self, *args):
        return subprocess.check_output(['git', '-C', str(self.root), *args], text=True).strip()

    def run_release(self, *args, **env):
        return subprocess.run(['bash', str(self.script), '--output-dir', str(self.output), *args],
                              env=dict(self.env, **env), capture_output=True, text=True)

    def calls(self):
        return [json.loads(line) for line in self.log.read_text().splitlines()]

    def test_default_loads_without_publishing(self):
        result = self.run_release()
        self.assertEqual(result.returncode, 0, result.stderr)
        record = json.loads((self.output/'release.json').read_text())
        self.assertFalse(record['published'])
        self.assertEqual(record['commit'], self.git('rev-parse', 'HEAD'))
        self.assertFalse((self.output/'image-digest.txt').exists())
        build = next(a for a in self.calls() if a[:2] == ['buildx', 'build'])
        self.assertIn('--load', build)
        self.assertEqual(build[build.index('--platform')+1], 'linux/amd64')
        self.assertFalse(Path(build[-1]).exists(), 'temporary source context leaked')

    def test_push_verifies_remote_before_emitting_digest(self):
        result = self.run_release('--push')
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertTrue(json.loads((self.output/'release.json').read_text())['published'])
        self.assertEqual((self.output/'image-digest.txt').read_text().strip(), 'sha256:'+'a'*64)
        self.assertTrue(any(a[:3] == ['buildx', 'imagetools', 'inspect'] for a in self.calls()))

    def test_dirty_source_fails_closed(self):
        (self.project/'package.json').write_text('{"dirty":true}')
        self.assertNotEqual(self.run_release('--push').returncode, 0)
        self.assertFalse(self.log.exists())

    def test_build_failure_cannot_emit_release(self):
        self.assertNotEqual(self.run_release('--push', FAIL_BUILD='1').returncode, 0)
        self.assertFalse((self.output/'image-digest.txt').exists())
        self.assertFalse((self.output/'release.json').exists())

    def test_remote_failure_cannot_emit_release(self):
        self.assertNotEqual(self.run_release('--push', FAIL_INSPECT='1').returncode, 0)
        self.assertFalse((self.output/'image-digest.txt').exists())
        self.assertFalse((self.output/'release.json').exists())

    def test_previous_release_is_not_overwritten(self):
        self.output.mkdir()
        (self.output/'image-digest.txt').write_text('previous-result\n')
        self.assertNotEqual(self.run_release('--push').returncode, 0)
        self.assertEqual((self.output/'image-digest.txt').read_text(), 'previous-result\n')
        self.assertFalse(any(a[:2] == ['buildx', 'build'] for a in self.calls()))

    def test_unknown_option_refused(self):
        self.assertEqual(self.run_release('--deploy').returncode, 2)
        self.assertFalse(self.log.exists())


if __name__ == '__main__':
    unittest.main()
