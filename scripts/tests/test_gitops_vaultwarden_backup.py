"""Exercise the actual native-backup wrapper's recovery and deletion boundaries."""
from __future__ import annotations

import hashlib
import json
import os
from pathlib import Path
import sqlite3
import subprocess
import tempfile
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
APP = ROOT / "platform/vaultwarden"
DOCS = list(yaml.safe_load_all((APP / "backup.yaml").read_text()))
CRON = next(d for d in DOCS if d["kind"] == "CronJob")
JOB = CRON["spec"]["jobTemplate"]["spec"]
POD = JOB["template"]["spec"]
CONTAINER = POD["containers"][0]
SCRIPT = CONTAINER["args"][0]


class VaultwardenNativeBackup(unittest.TestCase):
    def setUp(self):
        self.scratch = tempfile.TemporaryDirectory()
        self.addCleanup(self.scratch.cleanup)
        self.data = Path(self.scratch.name) / "data"
        self.data.mkdir(mode=0o700)
        with sqlite3.connect(self.data / "db.sqlite3") as db:
            db.executescript("CREATE TABLE users(id TEXT); CREATE TABLE ciphers(id TEXT);"
                             "INSERT INTO users VALUES('fixture-only');")
        self.live_hash = self.sha(self.data / "db.sqlite3")
        self.binary = Path(self.scratch.name) / "vaultwarden-fixture"
        self.binary.write_text("""#!/usr/bin/env python3
import os
from pathlib import Path
import sqlite3
import sys
p = Path(os.environ['DATA_FOLDER'])
mode = os.environ.get('FIXTURE_MODE', 'success')
if mode == 'failure':
    print('Backup failed. fixture-only')
    sys.exit(1)
if mode == 'unsafe':
    print("Backup to '%s/db.sqlite3' was successful" % p)
    sys.exit(0)
backup = p / 'db_20990101_120000.sqlite3'
if mode == 'empty':
    backup.touch()
else:
    with sqlite3.connect('file:%s?mode=ro' % (p / 'db.sqlite3'), uri=True) as db:
        db.execute('VACUUM INTO ?', (str(backup),))
print("Backup to '%s' was successful" % backup)
""")
        self.binary.chmod(0o700)

    @staticmethod
    def sha(path):
        return hashlib.sha256(path.read_bytes()).hexdigest()

    def run_script(self, mode="success"):
        return subprocess.run(["/bin/sh", "-ec", SCRIPT], env={
            **os.environ, "DATA_FOLDER": str(self.data),
            "VAULTWARDEN_BIN": str(self.binary), "FIXTURE_MODE": mode,
        }, text=True, capture_output=True)

    def test_completed_marker_matches_a_restorable_native_database(self):
        result = self.run_script()
        self.assertEqual(result.returncode, 0, result.stderr)
        marker = json.loads((self.data / ".sqlite-backup-complete.json").read_text())
        backup = self.data / marker["database"]
        self.assertEqual(marker["sha256"], self.sha(backup))
        self.assertRegex(marker["completedAt"], r"^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\dZ$")
        with sqlite3.connect(backup) as db:
            self.assertEqual(db.execute("PRAGMA integrity_check").fetchone(), ("ok",))
            self.assertEqual(db.execute("SELECT count(*) FROM users").fetchone(), (1,))
        self.assertEqual(self.sha(self.data / "db.sqlite3"), self.live_hash)
        self.assertEqual((self.data / ".sqlite-backup-complete.json").stat().st_mode & 0o777, 0o600)

    def test_retention_removes_only_old_exact_native_names(self):
        for i in range(1, 10):
            (self.data / f"db_202601{i:02d}_120000.sqlite3").write_bytes(b"old-native-copy")
        protected = ["db_other.sqlite3", "db_20260101_120000.sqlite3-wal",
                     "db_202601011_120000.sqlite3", "config.json", "rsa_key.pem"]
        for name in protected:
            (self.data / name).write_bytes(b"protected-fixture")
        attachments = self.data / "attachments"
        attachments.mkdir()
        (attachments / "fixture").write_bytes(b"attachment-fixture")
        result = self.run_script()
        self.assertEqual(result.returncode, 0, result.stderr)
        copies = sorted(p.name for p in self.data.glob("db_202601??_120000.sqlite3"))
        self.assertEqual(copies, [f"db_202601{i:02d}_120000.sqlite3" for i in range(4, 10)])
        self.assertTrue((self.data / "db_20990101_120000.sqlite3").is_file())
        self.assertTrue(all((self.data / name).read_bytes() == b"protected-fixture" for name in protected))
        self.assertEqual((attachments / "fixture").read_bytes(), b"attachment-fixture")
        self.assertEqual(self.sha(self.data / "db.sqlite3"), self.live_hash)

    def test_failures_never_advance_last_success_marker_or_remove_old_backups(self):
        marker = self.data / ".sqlite-backup-complete.json"
        marker.write_text('{"previous":"successful"}\n')
        older = self.data / "db_20260101_120000.sqlite3"
        older.write_bytes(b"old-native-copy")
        for mode in ("failure", "unsafe", "empty"):
            with self.subTest(mode=mode):
                result = self.run_script(mode)
                self.assertNotEqual(result.returncode, 0)
                self.assertEqual(marker.read_text(), '{"previous":"successful"}\n')
                self.assertEqual(older.read_bytes(), b"old-native-copy")
                self.assertEqual(self.sha(self.data / "db.sqlite3"), self.live_hash)
                (self.data / "db_20990101_120000.sqlite3").unlink(missing_ok=True)

    def test_job_is_bounded_local_and_co_located_without_api_rights(self):
        self.assertEqual(CRON["spec"]["schedule"], "20 13 * * *")
        self.assertEqual(CRON["spec"]["timeZone"], "Etc/UTC")
        self.assertEqual(CRON["spec"]["concurrencyPolicy"], "Forbid")
        self.assertEqual(JOB["activeDeadlineSeconds"], 600)
        self.assertLessEqual(JOB["backoffLimit"], 1)
        self.assertFalse(POD["automountServiceAccountToken"])
        self.assertEqual(POD["securityContext"]["runAsUser"], 1000)
        self.assertEqual(POD["volumes"][0]["persistentVolumeClaim"]["claimName"], "data-pvc")
        affinity, = POD["affinity"]["podAffinity"]["requiredDuringSchedulingIgnoredDuringExecution"]
        self.assertEqual(affinity["labelSelector"]["matchLabels"],
                         {"app.kubernetes.io/name": "vaultwarden", "app.kubernetes.io/component": "web"})
        self.assertEqual(affinity["topologyKey"], "kubernetes.io/hostname")
        self.assertNotIn("app.kubernetes.io/name", JOB["template"]["metadata"]["labels"])
        app_image = yaml.safe_load((APP / "deployment.yaml").read_text())["spec"]["template"]["spec"]["containers"][0]["image"]
        self.assertEqual(CONTAINER["image"], app_image)
        self.assertTrue(CONTAINER["securityContext"]["readOnlyRootFilesystem"])
        self.assertNotIn("envFrom", CONTAINER)
        self.assertFalse(any(d["kind"] in {"Role", "RoleBinding", "ClusterRole", "ClusterRoleBinding", "Secret"} for d in DOCS))
        self.assertTrue(all("argocd.argoproj.io/hook" not in d["metadata"].get("annotations", {}) for d in DOCS))

    def test_encrypted_snapshot_is_distinct_and_preserves_daily_schedule(self):
        source, = [d for d in DOCS if d["kind"] == "ReplicationSource"]
        self.assertEqual(source["spec"]["trigger"], {"schedule": "40 13 * * *"})
        self.assertEqual(source["spec"]["sourcePVC"], "data-pvc")
        restic = source["spec"]["restic"]
        self.assertEqual(restic["copyMethod"], "Snapshot")
        self.assertEqual(restic["retain"], {"daily": 7, "weekly": 4, "monthly": 3, "yearly": 1})
        sealed = yaml.safe_load((APP / "backup-data.sealedsecret.yaml").read_text())
        self.assertEqual(sealed["metadata"]["namespace"], "vaultwarden")
        self.assertEqual(sealed["metadata"]["name"], restic["repository"])
        self.assertEqual(set(sealed["spec"]["encryptedData"]),
                         {"RESTIC_REPOSITORY", "RESTIC_PASSWORD", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_DEFAULT_REGION"})
        self.assertTrue(all(value.startswith("Ag") for value in sealed["spec"]["encryptedData"].values()))
        for key in ("sealedsecrets.bitnami.com/namespace-wide", "sealedsecrets.bitnami.com/cluster-wide"):
            self.assertNotIn(key, sealed["metadata"].get("annotations", {}))


if __name__ == "__main__":
    unittest.main()
