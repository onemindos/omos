import importlib.util
import json
import os
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch, MagicMock

spec = importlib.util.spec_from_file_location("backup", Path(__file__).with_name("backup.py"))
backup = importlib.util.module_from_spec(spec)
spec.loader.exec_module(backup)


class IntegrityTests(unittest.TestCase):
    def test_owned_database_discovery_refuses_empty_backup(self):
        with patch.object(backup, "command", return_value=b""):
            with self.assertRaisesRegex(RuntimeError, "empty owned-database"):
                backup.postgres_backup(Path("/tmp/fake"), owned=True)

    def test_owned_restore_accepts_empty_plugins_but_not_empty_platform(self):
        metadata = {"databases": [{"file": "catalog.dump"}, {"file": "unused-plugin.dump"}]}
        with patch.object(backup, "postgres_restore_one", side_effect=[{"tables": 12}, {"tables": 0}]):
            self.assertEqual(backup.postgres_restore(Path("/tmp/fake"), metadata), {"databases": 2, "tables": 12})
        with patch.object(backup, "postgres_restore_one", return_value={"tables": 0}):
            with self.assertRaisesRegex(RuntimeError, "no application tables"):
                backup.postgres_restore(Path("/tmp/fake"), metadata)

    def test_large_archive_round_trip(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            source, encrypted, restored = [root / name for name in ("source", "encrypted", "restored")]
            source.write_bytes(os.urandom(backup.CHUNK * 2 + 17))
            backup.encrypt(source, encrypted, b"k" * 32)
            backup.decrypt(encrypted, restored, b"k" * 32)
            self.assertEqual(backup.digest(source), backup.digest(restored))

    def test_tampering_and_wrong_key_leave_no_cleartext(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            source, encrypted, restored = [root / name for name in ("source", "encrypted", "restored")]
            source.write_bytes(b"private archive")
            backup.encrypt(source, encrypted, b"k" * 32)
            with self.assertRaises(Exception):
                backup.decrypt(encrypted, restored, b"x" * 32)
            self.assertFalse(restored.exists())
            payload = bytearray(encrypted.read_bytes())
            payload[-1] ^= 1
            encrypted.write_bytes(payload)
            with self.assertRaises(Exception):
                backup.decrypt(encrypted, restored, b"k" * 32)
            self.assertFalse(restored.exists())

    def test_failed_command_does_not_expose_credentials(self):
        response = MagicMock(returncode=1, stderr=b"secret password", stdout=b"")
        with patch.object(backup.subprocess, "run", return_value=response):
            with self.assertRaisesRegex(RuntimeError, r"nats failed \(exit 1\)") as caught:
                backup.command(["nats", "--password", "secret password"])
        self.assertNotIn("secret password", str(caught.exception))

    def test_restore_never_uses_production_nats_credentials(self):
        server = MagicMock()
        server.poll.return_value = None
        commands = []

        def command(args, env=None):
            commands.append((args, env))
            return json.dumps({"config": {"name": "EVENTS"}}).encode()

        with patch.dict(os.environ, {"NATS_PASSWORD": "production-secret", "NATS_URL": "tls://production",
                                     "NATS_CA": "/production-ca", "NATS_CONTEXT": "production"}), \
                patch.object(backup.subprocess, "Popen", return_value=server) as start, \
                patch.object(backup.socket, "create_connection"), \
                patch.object(backup, "command", side_effect=command):
            backup.nats_restore(Path("/tmp/fake"), {"streams": ["EVENTS"]})
        self.assertIn("127.0.0.1", start.call_args.args[0])
        for _, env in commands:
            self.assertEqual(env["NATS_URL"], "nats://127.0.0.1:14222")
            self.assertNotIn("NATS_PASSWORD", env)
            self.assertNotIn("NATS_CA", env)
            self.assertNotIn("NATS_CONTEXT", env)
        server.terminate.assert_called_once()

    def test_failed_postgres_restore_drops_only_temporary_database(self):
        calls = []

        def command(args, env=None):
            calls.append(args)
            self.assertEqual(env["PGHOST"], "127.0.0.1")
            self.assertNotIn("PGPASSWORD", env)
            self.assertNotIn("PGDATABASE", env)
            if args[0] == "pg_restore":
                raise RuntimeError("restore failed")
            return b""

        with patch.object(backup, "command", side_effect=command):
            with self.assertRaisesRegex(RuntimeError, "restore failed"):
                backup.postgres_restore(Path("/tmp/fake"), {"database": "cot"})
        self.assertEqual(calls[-1][0], "dropdb")
        name = calls[-1][-1]
        self.assertTrue(name.startswith("onemind_restore_"))
        self.assertEqual(calls[0][-1], name)
        self.assertNotEqual(name, "cot")


if __name__ == "__main__":
    unittest.main()
