"""Guards for the Alertmanager Telegram receiver (cluster review A17).

The receiver is in the kube-prometheus-stack base config, not in an AlertmanagerConfig, because that
CRD only accepts the chat ID as a literal integer. Alertmanager reads the bot token and the chat ID as
files from Secret monitoring/alertmanager-telegram, mounted as an optional volume: a missing Secret only
fails Telegram sends and never stops Alertmanager from starting. The owner seals that Secret; once the
SealedSecret file exists, these tests also check that its keys are the files the receiver reads.
"""
import unittest
from pathlib import Path, PurePosixPath

import yaml

ROOT = Path(__file__).resolve().parents[2]
VALUES = ROOT / "platform/values/kube-prometheus-stack.yaml"
SEALED = ROOT / "platform/monitoring/alertmanager-telegram.sealedsecret.yaml"
SECRET_NAME = "alertmanager-telegram"


def alertmanager():
    return yaml.safe_load(VALUES.read_text())["alertmanager"]


def telegram_config(am):
    receiver = next(r for r in am["config"]["receivers"] if r["name"] == "telegram")
    (config,) = receiver["telegram_configs"]
    return config


def secret_mount(am):
    spec = am["alertmanagerSpec"]
    volume = next(v for v in spec["volumes"] if (v.get("secret") or {}).get("secretName") == SECRET_NAME)
    mount = next(m for m in spec["volumeMounts"] if m["name"] == volume["name"])
    return volume, mount


def credential_keys(am):
    """The Secret keys the receiver reads: the file names of bot_token_file and chat_id_file."""
    config = telegram_config(am)
    return {PurePosixPath(config[key]).name for key in ("bot_token_file", "chat_id_file")}


class TelegramReceiver(unittest.TestCase):
    def setUp(self):
        self.am = alertmanager()

    def test_token_and_chat_id_are_files_from_the_mounted_secret(self):
        config = telegram_config(self.am)
        # Inline values would put the token or the chat ID in Git.
        self.assertNotIn("bot_token", config)
        self.assertNotIn("chat_id", config)
        _, mount = secret_mount(self.am)
        base = PurePosixPath(mount["mountPath"])
        self.assertEqual(base / "bot-token", PurePosixPath(config["bot_token_file"]))
        self.assertEqual(base / "chat-id", PurePosixPath(config["chat_id_file"]))
        self.assertTrue(mount.get("readOnly"))
        self.assertTrue(config["send_resolved"])

    def test_a_missing_secret_cannot_stop_alertmanager(self):
        volume, _ = secret_mount(self.am)
        self.assertTrue(volume["secret"].get("optional"))
        # `secrets:` mounts are never optional: without the Secret the pod could not start.
        self.assertNotIn(SECRET_NAME, self.am["alertmanagerSpec"].get("secrets") or [])

    def test_config_survives_the_config_reloader(self):
        # prometheus-config-reloader expands $(VAR) references before Alertmanager reads the file.
        self.assertNotIn("$(", yaml.safe_dump(self.am["config"]))

    def test_watchdog_is_never_notified_and_critical_and_warning_reach_telegram(self):
        route = self.am["config"]["route"]
        # The operator rejects a root route with matchers or continue.
        self.assertEqual("null", route["receiver"])
        self.assertNotIn("matchers", route)
        self.assertFalse(route.get("continue", False))
        heartbeat, telegram = route["routes"][:2]
        self.assertEqual("null", heartbeat["receiver"])
        self.assertEqual(['alertname =~ "Watchdog|InfoInhibitor"'], heartbeat["matchers"])
        self.assertFalse(heartbeat.get("continue", False))
        self.assertEqual("telegram", telegram["receiver"])
        self.assertEqual(['severity =~ "critical|warning"'], telegram["matchers"])


class OwnerSealedSecret(unittest.TestCase):
    """Skipped until the owner commits the SealedSecret; the receiver fails safe without it."""

    def setUp(self):
        if not SEALED.is_file():
            self.skipTest(f"{SEALED.relative_to(ROOT)} is not committed yet (the owner seals it)")
        docs = [d for d in yaml.safe_load_all(SEALED.read_text()) if d]
        self.assertEqual(1, len(docs), f"{SEALED} must hold exactly one document")
        self.doc = docs[0]

    def test_is_the_sealed_secret_the_receiver_mounts(self):
        self.assertEqual("bitnami.com/v1alpha1", self.doc.get("apiVersion"))
        self.assertEqual("SealedSecret", self.doc.get("kind"))
        meta = self.doc.get("metadata") or {}
        self.assertEqual((SECRET_NAME, "monitoring"), (meta.get("name"), meta.get("namespace")))

    def test_holds_exactly_the_keys_the_receiver_reads_and_nothing_in_plaintext(self):
        spec = self.doc.get("spec") or {}
        encrypted = spec.get("encryptedData") or {}
        self.assertEqual(credential_keys(alertmanager()), set(encrypted))
        for key, value in encrypted.items():
            looks_sealed = isinstance(value, str) and value.startswith("Ag") and len(value) >= 100
            self.assertTrue(looks_sealed, f"{key} does not look sealed")
        for plain in ("data", "stringData"):
            self.assertNotIn(plain, self.doc)
        self.assertFalse((spec.get("template") or {}).get("data"), "template.data would be plaintext")


if __name__ == "__main__":
    unittest.main()
