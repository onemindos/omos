from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import yaml

import importlib.util
import sys
from types import SimpleNamespace
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location("litellm_provider_config", ROOT / "platform/hermes/agent/litellm_provider_config.py")
provider = importlib.util.module_from_spec(spec)
spec.loader.exec_module(provider)
main, merge = provider.main, provider.merge


class LitellmProviderTests(unittest.TestCase):
    def test_merge_selects_gateway_without_persisting_key(self):
        config = merge({"model": {"provider": "bedrock", "default": "old"}}, "unit-test-marker")
        self.assertEqual(config["model"], {"provider": "litellm", "default": "hermes-agent"})
        self.assertEqual(
            config["providers"]["litellm"]["base_url"],
            "http://litellm.onemind-ai.svc.cluster.local:4000/v1",
        )
        self.assertEqual(config["providers"]["litellm"]["key_env"], "HERMES_LITELLM_API_KEY")
        self.assertEqual(config["providers"]["litellm"]["default_model"], "hermes-agent")

        self.assertNotIn("unit-test-marker", yaml.safe_dump(config))

    def test_empty_key_reports_inactive(self):
        config = {"model": {"provider": "bedrock"}, "providers": {}}
        self.assertEqual(merge(config.copy(), ""), config)

    def test_preserves_other_settings_and_file_permissions(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "config.yaml"
            path.write_text(yaml.safe_dump({"model": {"provider": "custom", "base_url": "http://old.invalid/v1", "api_key": "old-unit-marker"}, "platforms": {"matrix": {"enabled": True}}, "tools": ["shell"]}))
            path.chmod(0o600)
            main(str(path), "unit-test-marker")
            result = yaml.safe_load(path.read_text())
            self.assertEqual(result["platforms"], {"matrix": {"enabled": True}})
            self.assertEqual(result["tools"], ["shell"])
            self.assertNotIn("base_url", result["model"])
            self.assertNotIn("api_key", result["model"])
            self.assertNotIn("unit-test-marker", path.read_text())
            self.assertEqual(path.stat().st_mode & 0o777, 0o600)

    def test_empty_key_does_not_write_config(self):
        with tempfile.TemporaryDirectory() as directory:
            config_path = Path(directory) / "config.yaml"
            config_path.write_text("existing: true\n", encoding="utf-8")
            self.assertEqual(main(str(config_path), ""), 0)
            with config_path.open(encoding="utf-8") as stream:
                self.assertEqual(yaml.safe_load(stream), {"existing": True})



class BootstrapTests(unittest.TestCase):
    def provision(self, hermes_present):
        calls = []
        class Client:
            def __init__(self, **kwargs): pass
            def __enter__(self): return self
            def __exit__(self, *args): pass
            def post(self, url, headers, json):
                calls.append(json)
                return SimpleNamespace(status_code=200, raise_for_status=lambda: None)
        doc = yaml.safe_load((ROOT / "platform/ecosystem/resources/model-client-bootstrap.yaml").read_text())
        code = doc["spec"]["template"]["spec"]["containers"][0]["args"][0]
        env = {"MASTER_KEY": "test-marker", "core": "test-core", "openchoreo": "test-openchoreo"}
        if hermes_present: env["hermes"] = "test-hermes"
        with patch.dict("os.environ", env, clear=True), patch.dict(sys.modules, {"httpx": SimpleNamespace(Client=Client)}):
            exec(compile(code, "bootstrap", "exec"), {})
        return calls

    def test_absent_hermes_key_keeps_existing_clients_working(self):
        calls = self.provision(False)
        self.assertEqual([c["metadata"]["service"] for c in calls], ["core", "openchoreo"])
        self.assertTrue(all(c["models"] == ["onemind-chat", "onemind-fast"] for c in calls))

    def test_hermes_key_has_separate_scoped_model_allowlist(self):
        calls = self.provision(True)
        hermes = next(c for c in calls if c["metadata"]["service"] == "hermes")
        self.assertEqual(hermes["models"], ["hermes-agent", "local-tools", "local-embed"])
        self.assertEqual(hermes["key_alias"], "onemind-hermes")

if __name__ == "__main__":
    unittest.main()
