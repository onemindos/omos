"""Fail-closed behavior of the PgSTAC/AGE migration compatibility boundary."""
from hashlib import sha256
from pathlib import Path
import sys
import tempfile
from types import ModuleType
import unittest
from unittest.mock import patch

import yaml

ROOT = Path(__file__).resolve().parents[2]
SOURCE = ROOT / 'platform/stac/migrate-compat.py'


class MigrationCompatibility(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.root = Path(self.directory.name)
        self.sql_file = self.root / 'pgstac.0.9.12.sql'
        self.calls = []
        migration = ModuleType('pypgstac.migrate')
        migration.migrations_dir = str(self.root)
        migration.get_sql = lambda filename: self.sql_file.read_text()
        package = ModuleType('pypgstac')
        package.__version__ = '0.9.12'
        package.migrate = migration
        cli = ModuleType('pypgstac.pypgstac')
        cli.cli = lambda: self.calls.append('original CLI')
        self.upstream = migration
        self.module = ModuleType('stac_compat')
        with patch.dict(sys.modules, {'pypgstac': package, 'pypgstac.migrate': migration,
                                     'pypgstac.pypgstac': cli}):
            exec(compile(SOURCE.read_text(), str(SOURCE), 'exec'), self.module.__dict__)
        self.sql = 'BEGIN;\n' + self.module.ORIGINAL + ';\nSELECT set_version(\'0.9.12\');\nCOMMIT;\n'
        self.sql_file.write_text(self.sql)
        # Test a pinned synthetic input without vendoring the full upstream SQL.
        self.module.EXPECTED_SQL_SHA256 = sha256(self.sql_file.read_bytes()).hexdigest()

    def test_rewrite_preserves_the_whole_transaction_and_all_other_sql(self):
        rewritten = self.module.compatibility_sql(str(self.sql_file))
        self.assertEqual(rewritten, self.sql.replace(self.module.ORIGINAL, self.module.COMPATIBLE))
        self.assertTrue(rewritten.startswith('BEGIN;\n'))
        self.assertTrue(rewritten.endswith('COMMIT;\n'))

    def test_changed_sql_is_rejected_before_it_can_execute(self):
        self.sql_file.write_text(self.sql + '-- unexpected upstream change\n')
        with self.assertRaisesRegex(RuntimeError, 'checksum changed'):
            self.module.compatibility_sql(str(self.sql_file))

    def test_upgrades_missing_and_ambiguous_rewrites_are_rejected(self):
        with self.assertRaisesRegex(RuntimeError, 'bootstrap'):
            self.module.compatibility_sql('pgstac.0.9.11-0.9.12.sql')
        for sql in ['SELECT 1;', self.sql + self.module.ORIGINAL]:
            self.sql_file.write_text(sql)
            self.module.EXPECTED_SQL_SHA256 = sha256(self.sql_file.read_bytes()).hexdigest()
            with self.assertRaisesRegex(RuntimeError, 'exactly one'):
                self.module.compatibility_sql(str(self.sql_file))

    def test_original_cli_runs_only_for_the_reviewed_client(self):
        self.module.main()
        self.assertIs(self.upstream.get_sql, self.module.compatibility_sql)
        self.assertEqual(self.calls, ['original CLI'])
        self.module.__version__ = '0.9.13'
        with self.assertRaisesRegex(RuntimeError, 'client version'):
            self.module.main()
        self.assertEqual(self.calls, ['original CLI'])

    def test_production_and_synthetic_runs_use_the_identical_reviewed_wrapper(self):
        production = yaml.safe_load((ROOT / 'platform/stac/migration-compat.configmap.yaml').read_text())
        diagnostic = yaml.safe_load((ROOT / 'platform/stac-validation/scripts.yaml').read_text())
        self.assertEqual(production['data']['migrate-compat.py'], SOURCE.read_text())
        self.assertEqual(diagnostic['data']['migrate-compat.py'], SOURCE.read_text())


if __name__ == '__main__':
    unittest.main()
