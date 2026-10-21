"""Offline proposal-contract regression; not app/runtime acceptance."""
from pathlib import Path
import unittest
import yaml

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / 'services/onemind-portal/templates/app-onboarding'


class AppLifecycleTemplateTests(unittest.TestCase):
    def test_existing_proposal_remains_draft_only(self):
        template = yaml.safe_load((BASE / 'template.yaml').read_text())
        steps = template['spec']['steps']
        self.assertEqual([s['action'] for s in steps],
                         ['fetch:template', 'publish:github:pull-request'])
        publish = steps[1]['input']
        self.assertTrue(publish['draft'])
        self.assertEqual(publish['targetPath'], 'proposals/apps/${{ parameters.name }}')
        self.assertEqual(publish['token'], '${{ secrets.githubToken }}')

    def test_review_contract_and_no_old_live_directory_instruction(self):
        text = (BASE / 'skeleton/README.md').read_text()
        for section in ('Identity', 'Ownership', 'Artifact', 'Access', 'Runtime',
                        'Signals', 'Relationships', 'AI access', 'Data', 'Recovery',
                        'Lifecycle', 'Cost', 'Release', 'Acceptance', 'Retirement'):
            self.assertIn(f'**{section}:**', text)
        self.assertIn('not-tested / not-applicable', text)
        self.assertIn('not the live', text)
        self.assertNotIn('Add the matching card to', text)
        self.assertIn('No automatic namespace/PVC/Secret/backup deletion', text)
        self.assertIn('Image rollback may not undo database changes', text)


if __name__ == '__main__':
    unittest.main()
