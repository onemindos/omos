"""Run dependency-free portal workstream contracts in existing GitOps gates."""
from pathlib import Path
import subprocess
import unittest

ROOT = Path(__file__).resolve().parents[2]


class PortalWorkstreamTests(unittest.TestCase):
    def test_offline_node_contracts(self):
        result = subprocess.run(
            ['node', '--test', 'services/onemind-portal/scripts/test-cost-allocation.cjs'],
            cwd=ROOT, capture_output=True, text=True, timeout=60)
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)


if __name__ == '__main__':
    unittest.main()
