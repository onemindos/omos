"""Offline SQLite migration, archive safety and GitOps wiring contracts."""
from pathlib import Path
import importlib.util
import io
import sqlite3
import tarfile
import tempfile
import unittest
import yaml

ROOT = Path(__file__).resolve().parents[2]
SOURCE = ROOT / 'platform/monitoring/grafana-state.py'
SPEC = importlib.util.spec_from_file_location('grafana_state', SOURCE)
state = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(state)


def fixture(folder):
    database = folder / 'grafana.db'
    with sqlite3.connect(database) as db:
        for table in ['user', 'dashboard', 'data_source']:
            db.execute(f'CREATE TABLE {table} (id INTEGER PRIMARY KEY, value TEXT)')
            db.execute(f'INSERT INTO {table} VALUES (1, "fixture")')
    out = io.BytesIO()
    with tarfile.open(fileobj=out, mode='w:') as tar:
        tar.add(database, arcname='grafana.db')
    return out.getvalue()


class GrafanaStateTests(unittest.TestCase):
    def test_seed_preserves_all_table_counts_and_requires_remote_restore(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp); data = root / 'data'; data.mkdir()
            payload = fixture(root)
            seen = []
            def verified(path):
                seen.append(state.fingerprint(path))
                return {'restoreVerified': True, 'backup': 'fixture'}
            receipt = state.seed(payload, data, verified)
            self.assertTrue(receipt['restoreVerified'])
            self.assertEqual(state.fingerprint(data / 'grafana.db'), seen[0])
            self.assertEqual(len(seen), 1)
            self.assertTrue((data / state.MARKER).exists())
            with self.assertRaises(ValueError): state.seed(payload, data, verified)

    def test_failed_backup_never_seeds_or_marks_a_volume(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp); data = root / 'data'; data.mkdir()
            with self.assertRaises(ValueError):
                state.seed(fixture(root), data, lambda _: {'restoreVerified': False})
            self.assertFalse((data / 'grafana.db').exists())
            self.assertFalse((data / state.MARKER).exists())

    def test_unknown_duplicate_or_link_members_are_rejected(self):
        for name, kind in [('../grafana.db', tarfile.REGTYPE), ('grafana.db', tarfile.SYMTYPE),
                           ('plugin', tarfile.REGTYPE)]:
            with self.subTest(name=name, kind=kind), tempfile.TemporaryDirectory() as tmp:
                payload = io.BytesIO()
                with tarfile.open(fileobj=payload, mode='w:') as tar:
                    entry = tarfile.TarInfo(name); entry.type = kind; entry.linkname = '/etc/passwd'
                    tar.addfile(entry)
                with self.assertRaises(ValueError): state.seed_snapshot(payload.getvalue(), Path(tmp))
        with tempfile.TemporaryDirectory() as tmp:
            payload = io.BytesIO()
            with tarfile.open(fileobj=payload, mode='w:') as tar:
                for _ in range(2): tar.addfile(tarfile.TarInfo('grafana.db'))
            with self.assertRaises(ValueError): state.seed_snapshot(payload.getvalue(), Path(tmp))

    def test_empty_oversized_and_non_grafana_databases_are_rejected(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            with self.assertRaises(ValueError): state.seed_snapshot(b'', root)
            with self.assertRaises(ValueError): state.seed_snapshot(b'x' * (state.MAX_SEED + 1), root)
            db = root / 'other.db'
            with sqlite3.connect(db) as conn: conn.execute('CREATE TABLE other (id INTEGER)')
            with self.assertRaises(ValueError): state.fingerprint(db)

    def test_mounted_script_and_source_match(self):
        docs = list(yaml.safe_load_all((ROOT / 'platform/monitoring/grafana-state.yaml').read_text()))
        config = next(d for d in docs if d['kind'] == 'ConfigMap')
        self.assertEqual(config['data']['state.py'], SOURCE.read_text())
        pvc = next(d for d in docs if d['kind'] == 'PersistentVolumeClaim')
        self.assertEqual(pvc['spec']['resources']['requests']['storage'], '5Gi')
        self.assertEqual(pvc['spec']['storageClassName'], 'do-block-storage-retain')
        self.assertEqual(pvc['metadata']['annotations']['argocd.argoproj.io/sync-options'], 'Prune=false')
        cron = next(d for d in docs if d['kind'] == 'CronJob')
        self.assertEqual(cron['spec']['concurrencyPolicy'], 'Forbid')
        self.assertTrue(cron['spec']['suspend'])  # enabled only after durable cutover
        for template in [next(d for d in docs if d['kind'] == 'Job')['spec']['template'], cron['spec']['jobTemplate']['spec']['template']]:
            self.assertFalse(template['spec']['automountServiceAccountToken'])
            self.assertEqual(template['spec']['securityContext']['runAsUser'], 472)
            self.assertIn('@sha256:', template['spec']['containers'][0]['image'])


if __name__ == '__main__':
    unittest.main()
