"""Run the deployed offsite shell against local remotes using the real rclone engine.

Set RCLONE_BIN to the installed, verified rclone executable in CI. No cloud accounts,
credentials, network services or production backup data are used by these fixtures.
"""
from __future__ import annotations

import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import time
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
DOCS = list(yaml.safe_load_all((ROOT / "platform/backups/offsite-copy.yaml").read_text()))
CRON = next(d for d in DOCS if d["kind"] == "CronJob")
SCRIPT = CRON["spec"]["jobTemplate"]["spec"]["template"]["spec"]["containers"][0]["command"][-1]
RCLONE = os.environ.get("RCLONE_BIN") or shutil.which("rclone")


@unittest.skipUnless(RCLONE, "Install verified rclone or set RCLONE_BIN for engine tests")
class OffsiteCopyEngine(unittest.TestCase):
    def setUp(self):
        self.scratch = tempfile.TemporaryDirectory()
        self.addCleanup(self.scratch.cleanup)
        self.root = Path(self.scratch.name)
        self.source = self.root / "source"
        self.destination = self.root / "destination"
        self.binary_dir = self.root / "bin"
        self.tmp = self.root / "tmp"
        for directory in (self.source, self.destination, self.binary_dir, self.tmp):
            directory.mkdir()
        self.mirror = self.destination / "mirror"
        self.receipt = self.destination / "_offsite/last-success.json"
        wrapper = self.binary_dir / "rclone"
        wrapper.write_text(f"#!{sys.executable}\n" + '''import os
from pathlib import Path
import subprocess
import sys
import time

args = sys.argv[1:]
mode = os.environ.get("FIXTURE_MODE", "success")
if args[0] == "lsf" and mode == "list-fails":
    print("archive.enc")
    sys.exit(23)
if args[0] == "copy" and mode == "copy-fails":
    sys.exit(23)
if args[0] == "check" and mode == "corrupt-before-check":
    (Path(os.environ["FIXTURE_DESTINATION"]) / "mirror/archive.enc").write_bytes(b"bad-size")
result = subprocess.run([os.environ["FIXTURE_RCLONE"], *args])
if args[0] == "copy" and result.returncode == 0 and mode == "cross-age-threshold":
    time.sleep(1.3)
if args[0] == "lsf" and result.returncode == 0:
    if mode == "disappear-after-list":
        (Path(os.environ["FIXTURE_SOURCE"]) / "archive.enc").unlink()
    elif mode == "arrive-after-list":
        path = Path(os.environ["FIXTURE_SOURCE"]) / "late-arrival.enc"
        path.write_bytes(b"arrived after selection")
        os.utime(path, (time.time() - 3600, time.time() - 3600))
sys.exit(result.returncode)
''')
        wrapper.chmod(0o700)

    def fixture(self, name="archive.enc", data=b"encrypted fixture", age=3600):
        path = self.source / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(data)
        os.utime(path, (time.time() - age, time.time() - age))
        return path

    def run_script(self, mode="success", min_age="0s", max_age=""):
        env = {
            **os.environ,
            "PATH": str(self.binary_dir) + os.pathsep + os.environ["PATH"],
            "RCLONE_CONFIG": os.devnull,
            "RCLONE_CONFIG_SRC_TYPE": "local", "RCLONE_CONFIG_DST_TYPE": "local",
            "SRC_BUCKET": str(self.source), "DST_BUCKET": str(self.destination),
            "DST_PATH": "mirror", "COPY_MIN_AGE": min_age, "COPY_MAX_AGE": max_age,
            "TMPDIR": str(self.tmp), "FIXTURE_MODE": mode,
            "FIXTURE_RCLONE": str(Path(RCLONE).resolve()),
            "FIXTURE_SOURCE": str(self.source), "FIXTURE_DESTINATION": str(self.destination),
        }
        result = subprocess.run(["/bin/sh", "-eu", "-c", SCRIPT], env=env,
                                text=True, capture_output=True, timeout=30)
        self.assertEqual(list(self.tmp.iterdir()), [], "Selection manifest must be removed")
        return result

    def assert_success(self, result):
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        receipt = json.loads(self.receipt.read_text())
        self.assertEqual(receipt["source"], str(self.source))
        self.assertRegex(receipt["completed"], r"^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\dZ$")

    def test_age_threshold_crossing_uses_the_original_selection(self):
        self.fixture("archive.enc")
        self.fixture("too-young.enc", age=0)
        self.assert_success(self.run_script("cross-age-threshold", min_age="1s"))
        self.assertTrue((self.mirror / "archive.enc").is_file())
        self.assertFalse((self.mirror / "too-young.enc").exists())

    def test_late_arrival_does_not_change_the_verification_set(self):
        self.fixture()
        self.assert_success(self.run_script("arrive-after-list"))
        self.assertTrue((self.source / "late-arrival.enc").is_file())
        self.assertFalse((self.mirror / "late-arrival.enc").exists())

    def test_source_disappearing_after_selection_cannot_produce_false_success(self):
        self.fixture()
        result = self.run_script("disappear-after-list")
        self.assertNotEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertFalse(self.receipt.exists())

    def test_restic_locks_are_excluded_and_literal_paths_are_preserved(self):
        self.fixture("backups/restic/example/claim/locks/lock-id")
        self.fixture("backups/restic/example/claim/data/pack-id")
        self.fixture("# archive with spaces.enc")
        self.assert_success(self.run_script())
        self.assertFalse((self.mirror / "backups/restic/example/claim/locks/lock-id").exists())
        self.assertTrue((self.mirror / "backups/restic/example/claim/data/pack-id").is_file())
        self.assertTrue((self.mirror / "# archive with spaces.enc").is_file())

    def test_source_deletions_do_not_delete_existing_mirror_objects(self):
        deleted = self.fixture("previous.enc", b"keep independently")
        self.assert_success(self.run_script())
        deleted.unlink()
        self.fixture("new.enc")
        self.assert_success(self.run_script())
        self.assertEqual((self.mirror / "previous.enc").read_bytes(), b"keep independently")

    def test_overwrites_retain_the_previous_content_in_backup_directory(self):
        self.fixture(data=b"version one", age=7200)
        self.assert_success(self.run_script())
        self.fixture(data=b"version two", age=3600)
        self.assert_success(self.run_script())
        self.assertEqual((self.mirror / "archive.enc").read_bytes(), b"version two")
        versions = list((self.destination / "_overwritten").glob("*/archive.enc"))
        self.assertEqual(len(versions), 1)
        self.assertEqual(versions[0].read_bytes(), b"version one")

    def test_optional_max_age_is_applied_when_the_selection_is_created(self):
        self.fixture("recent.enc", age=3600)
        self.fixture("expired.enc", age=172800)
        self.assert_success(self.run_script(max_age="1d"))
        self.assertTrue((self.mirror / "recent.enc").is_file())
        self.assertFalse((self.mirror / "expired.enc").exists())

    def test_empty_selection_never_records_success_or_copies_unselected_objects(self):
        self.fixture(age=0)
        result = self.run_script(min_age="1h")
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse(self.receipt.exists())
        self.assertFalse(self.mirror.exists())

    def test_listing_copy_and_real_verification_failures_preserve_previous_receipt(self):
        self.fixture()
        self.receipt.parent.mkdir()
        self.receipt.write_bytes(b"previous successful receipt")
        for mode in ("list-fails", "copy-fails", "corrupt-before-check"):
            with self.subTest(mode=mode):
                result = self.run_script(mode)
                self.assertNotEqual(result.returncode, 0, result.stdout + result.stderr)
                self.assertEqual(self.receipt.read_bytes(), b"previous successful receipt")


if __name__ == "__main__":
    unittest.main()
