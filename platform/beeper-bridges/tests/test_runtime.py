"""Offline checks: no Beeper authentication, downloads or connector jobs run."""

import hashlib
import json
import os
from pathlib import Path
import secrets
import stat
import subprocess
import tempfile
import unittest


SCRIPTS = Path(__file__).resolve().parents[1] / "scripts"


class RuntimeTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.commands = self.root / "commands"
        self.commands.mkdir()
        self.data = self.root / "data"
        self.runtime = self.root / "runtime"
        self.result = self.root / "result.json"
        self.env = os.environ.copy()
        self.env.update(
            PATH=str(self.commands) + os.pathsep + self.env["PATH"],
            BRIDGE_TYPE="telegram",
            BRIDGE_NAME="sh-onemind-telegram",
            BEEPER_DATA_DIR=str(self.data),
            BEEPER_RUNTIME_DIR=str(self.runtime),
            TEST_RESULT=str(self.result),
        )

    def command(self, name, body):
        path = self.commands / name
        path.write_text("#!/usr/bin/env python3\n" + body)
        path.chmod(0o755)
        return path

    def run_script(self, name, *args, cwd=None):
        return subprocess.run(
            ["bash", str(SCRIPTS / name), *map(str, args)],
            env=self.env,
            cwd=cwd,
            capture_output=True,
            text=True,
            check=False,
        )

    def test_private_config_and_exec_without_token_in_arguments(self):
        token = "bat_" + secrets.token_urlsafe(24)
        token_file = self.root / "token"
        token_file.write_text(token + "\n")
        self.env["BEEPER_TOKEN_FILE"] = str(token_file)
        self.command(
            "bbctl",
            "import json,os,sys\n"
            "from pathlib import Path\n"
            "Path(os.environ['TEST_RESULT']).write_text(json.dumps({'args':sys.argv[1:],'pid':os.getpid()}))\n",
        )
        proc = subprocess.Popen(
            ["bash", str(SCRIPTS / "run-bridge.sh")],
            env=self.env,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )
        stdout, stderr = proc.communicate(timeout=10)
        self.assertEqual(proc.returncode, 0, stderr)
        result = json.loads(self.result.read_text())
        self.assertEqual(result["pid"], proc.pid, "entrypoint must exec bbctl")
        self.assertIn("--custom-startup-command", result["args"])
        self.assertIn("--no-override-config", result["args"])
        self.assertNotIn(token, stdout + stderr + json.dumps(result))
        cfg_path = self.runtime / "bbctl.json"
        self.assertEqual(stat.S_IMODE(cfg_path.stat().st_mode), 0o600)
        cfg = json.loads(cfg_path.read_text())["environments"]["prod"]
        self.assertEqual(cfg["access_token"], token)
        self.assertEqual(cfg["database_dir"], str(self.data / "db"))

    def test_empty_token_fails_before_bbctl(self):
        token_file = self.root / "token"
        token_file.write_text("")
        self.env["BEEPER_TOKEN_FILE"] = str(token_file)
        result = self.run_script("run-bridge.sh")
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse((self.runtime / "bbctl.json").exists())

    def test_telegram_credentials_reach_config_but_not_connector_environment(self):
        api_id = str(secrets.randbelow(1_000_000) + 1)
        api_hash = secrets.token_hex(16)
        self.env.update(TELEGRAM_API_ID=api_id, TELEGRAM_API_HASH=api_hash)
        config = self.root / "config.yaml"
        config.write_text(
            "network:\n    api_id: 0\n    api_hash: unused\n"
            "appservice:\n    hostname: 0.0.0.0\n    port: 4000\n"
        )
        (self.data / "bin").mkdir(parents=True)
        connector = self.data / "bin" / "mautrix-telegram"
        connector.write_text(
            "#!/usr/bin/env python3\nimport json,os,sys\nfrom pathlib import Path\n"
            "Path(os.environ['TEST_RESULT']).write_text(json.dumps({'args':sys.argv[1:],"
            "'has_id':'TELEGRAM_API_ID' in os.environ,'has_hash':'TELEGRAM_API_HASH' in os.environ}))\n"
        )
        connector.chmod(0o755)
        result = self.run_script("start-connector.sh", "-c", config)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertNotIn(api_hash, result.stdout + result.stderr)
        rendered = config.read_text()
        self.assertIn("api_id: " + api_id, rendered)
        self.assertIn('api_hash: "' + api_hash + '"', rendered)
        self.assertIn("hostname: 127.0.0.1", rendered)
        self.assertEqual(stat.S_IMODE(config.stat().st_mode), 0o600)
        child = json.loads(self.result.read_text())
        self.assertFalse(child["has_id"])
        self.assertFalse(child["has_hash"])
        self.assertNotIn(api_hash, json.dumps(child))

    def test_changed_config_shape_fails_without_replacing_original(self):
        self.env.update(TELEGRAM_API_ID="1", TELEGRAM_API_HASH=secrets.token_hex(16))
        config = self.root / "config.yaml"
        original = "network:\n    api_id: 0\nappservice:\n    port: 4000\n"
        config.write_text(original)
        result = self.run_script("start-connector.sh", "-c", config)
        self.assertNotEqual(result.returncode, 0)
        self.assertEqual(config.read_text(), original)
        self.assertFalse(self.result.exists())

    def download_fixture(self, expected_hash=None):
        payload = self.root / "payload"
        payload.write_bytes(b"connector fixture content\n")
        digest = expected_hash or hashlib.sha256(payload.read_bytes()).hexdigest()
        lock = self.root / "artifacts.json"
        lock.write_text(json.dumps({"telegram": {
            "binary": "mautrix-telegram",
            "url": "https://mau.dev/mautrix/telegram/-/jobs/1/artifacts/raw/mautrix-telegram",
            "sha256": digest,
        }}))
        self.env.update(
            BEEPER_ARTIFACT_LOCK=str(lock),
            TEST_DOWNLOAD_SOURCE=str(payload),
            TEST_DOWNLOAD_COUNT=str(self.root / "downloads"),
        )
        self.command("uname", "print('x86_64')\n")
        self.command(
            "curl",
            "import os,shutil,sys\nfrom pathlib import Path\n"
            "shutil.copyfile(os.environ['TEST_DOWNLOAD_SOURCE'],sys.argv[sys.argv.index('--output')+1])\n"
            "p=Path(os.environ['TEST_DOWNLOAD_COUNT']); p.write_text(p.read_text()+'1' if p.exists() else '1')\n",
        )
        return payload, lock

    def test_checksum_cache_and_corruption_recovery(self):
        payload, _ = self.download_fixture()
        target = self.data / "bin" / "mautrix-telegram"
        for _ in range(2):
            result = self.run_script("fetch-bridge.sh")
            self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual((self.root / "downloads").read_text(), "1")
        self.assertEqual(target.read_bytes(), payload.read_bytes())
        self.assertEqual(stat.S_IMODE(target.stat().st_mode), 0o755)
        target.write_bytes(b"corrupted cache")
        result = self.run_script("fetch-bridge.sh")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual((self.root / "downloads").read_text(), "11")
        self.assertEqual(target.read_bytes(), payload.read_bytes())

    def test_bad_checksum_never_installs_binary(self):
        self.download_fixture(expected_hash="0" * 64)
        result = self.run_script("fetch-bridge.sh")
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse((self.data / "bin" / "mautrix-telegram").exists())
        self.assertEqual(list((self.data / "bin").iterdir()), [])

    def test_wrong_architecture_fails_before_downloading(self):
        self.download_fixture()
        self.command("uname", "print('aarch64')\n")
        result = self.run_script("fetch-bridge.sh")
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse((self.root / "downloads").exists())


if __name__ == "__main__":
    unittest.main()
