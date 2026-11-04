"""Compile embedded Python hooks without executing their privileged operations."""
from pathlib import Path
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]


class PythonHookSyntax(unittest.TestCase):
    def test_all_resource_job_python_commands_compile(self):
        checked = 0
        for path in (ROOT / 'platform/ecosystem/resources').glob('*.yaml'):
            for doc in yaml.safe_load_all(path.read_text()):
                if not isinstance(doc, dict) or doc.get('kind') != 'Job':
                    continue
                containers = doc['spec']['template']['spec'].get('containers', [])
                for container in containers:
                    command = container.get('command', [])
                    if command not in (['python', '-c'], ['python3', '-c']):
                        continue
                    with self.subTest(path=str(path.relative_to(ROOT)), container=container['name']):
                        compile(container['args'][0], str(path), 'exec')
                    checked += 1
        self.assertGreater(checked, 0)


if __name__ == '__main__':
    unittest.main()
