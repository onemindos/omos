"""Offline build/activation preparation contracts. No image or vault access."""
import importlib.util
from pathlib import Path
import subprocess
import unittest
from unittest.mock import patch

import yaml

HERE = Path(__file__).parent
REPO = HERE.parents[3]
WORKFLOW = REPO / '.github/workflows/build-hermes-runtime.yml'
spec = importlib.util.spec_from_file_location('prepare_image', HERE / 'proposal/prepare_image.py')
assert spec is not None and spec.loader is not None
prepare = importlib.util.module_from_spec(spec)
spec.loader.exec_module(prepare)


class BuildContracts(unittest.TestCase):
    def setUp(self):
        # BaseLoader preserves Actions' "on" key instead of YAML 1.1 boolean coercion.
        self.workflow = yaml.load(WORKFLOW.read_text(), Loader=yaml.BaseLoader)
        self.steps = self.workflow['jobs']['build']['steps']

    def test_manual_only_and_minimum_permissions(self):
        self.assertEqual(self.workflow['on'], {'workflow_dispatch': ''})
        self.assertEqual(self.workflow['permissions'], {'contents': 'read', 'packages': 'write'})
        self.assertEqual(self.steps[0]['with']['persist-credentials'], 'false')
        for step in self.steps:
            if 'uses' in step:
                self.assertRegex(step['uses'], r'^[\w/-]+@[0-9a-f]{40}$')
        self.assertEqual(WORKFLOW.read_text().count('${{ secrets.'), 1)
        self.assertIn('${{ secrets.GITHUB_TOKEN }}', WORKFLOW.read_text())

    def test_build_test_then_publish_same_image(self):
        build = next(s for s in self.steps if s.get('uses', '').startswith('docker/build-push-action@'))
        self.assertEqual(build['with']['platforms'], 'linux/amd64')
        self.assertEqual(build['with']['load'], 'true')
        self.assertEqual(build['with']['push'], 'false')
        self.assertEqual(build['with']['context'], str((HERE / 'runtime-fix').relative_to(REPO)))
        test = next(s for s in self.steps if '--user 10000:10000' in s.get('run', ''))
        for flag in ('--network none', '--read-only', 'test_child_policy.py', '--cap-drop ALL'):
            self.assertIn(flag, test['run'])
        login = next(s for s in self.steps if s.get('uses', '').startswith('docker/login-action@'))
        publish = next(s for s in self.steps if s.get('id') == 'digest')
        self.assertLess(self.steps.index(test), self.steps.index(login))
        self.assertLess(self.steps.index(login), self.steps.index(publish))
        self.assertIn('docker push "$IMAGE"', publish['run'])
        self.assertIn('docker buildx imagetools inspect "$pinned"', publish['run'])
        self.assertIn('hermes-runtime-image.txt', publish['run'])
        tag = self.workflow['jobs']['build']['env']['IMAGE']
        for key in ('github.sha', 'github.run_id', 'github.run_attempt'):
            self.assertIn('${{ ' + key + ' }}', tag)
        for step in self.steps:
            if 'run' in step:
                subprocess.run(['bash', '-n'], input=step['run'], text=True, check=True)

    def test_digest_substitution_only_changes_proposal_gateway(self):
        # Synthetic digest is a format fixture only, never a deployable image claim.
        ref = 'ghcr.io/onemindos/hermes-agent@sha256:' + 'a' * 64
        original = (HERE / 'proposal/deployment.yaml').read_text()
        rendered = prepare.render(ref)
        before = yaml.safe_load(original)['spec']['template']['spec']
        after = yaml.safe_load(rendered)['spec']['template']['spec']
        self.assertEqual(after['containers'][0]['image'], ref)
        self.assertEqual(before['initContainers'], after['initContainers'])
        after['containers'][0]['image'] = before['containers'][0]['image']
        self.assertEqual(before, after)
        self.assertEqual(original, (HERE / 'proposal/deployment.yaml').read_text())

    def test_rejects_unsafe_references_and_snapshot_drift(self):
        for ref in ('', prepare.UPSTREAM, 'ghcr.io/onemindos/hermes-agent:latest',
                    'ghcr.io/other/hermes-agent@sha256:' + 'a' * 64,
                    'ghcr.io/onemindos/hermes-agent@' + prepare.UPSTREAM.split('@')[1],
                    'ghcr.io/onemindos/hermes-agent@sha256:' + 'A' * 64):
            with self.subTest(ref=ref), self.assertRaises(ValueError):
                prepare.render(ref)
        with patch.object(Path, 'read_text', return_value='unexpected snapshot'):
            with self.assertRaisesRegex(ValueError, 'drift'):
                prepare.render('ghcr.io/onemindos/hermes-agent@sha256:' + 'a' * 64)


if __name__ == '__main__':
    unittest.main()
