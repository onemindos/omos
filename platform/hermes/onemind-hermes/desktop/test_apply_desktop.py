"""Tests for the Webapp step of apply_desktop.py (pure planning, no Hermes needed)."""
import importlib.util
import pathlib
import unittest

HERE = pathlib.Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location("apply_desktop", HERE / "apply_desktop.py")
apply_desktop = importlib.util.module_from_spec(spec)
spec.loader.exec_module(apply_desktop)


class PlanWebapp(unittest.TestCase):
    def test_empty_config_gets_the_block(self):
        block, changes = apply_desktop.plan_webapp({})
        self.assertEqual(block, apply_desktop.WEBAPP)
        self.assertEqual(sorted(k for k, _, _ in changes), sorted(apply_desktop.WEBAPP))

    def test_second_run_is_unchanged(self):
        block, _ = apply_desktop.plan_webapp({})
        again, changes = apply_desktop.plan_webapp({"webapp": block})
        self.assertEqual(changes, [])
        self.assertEqual(again, block)

    def test_other_keys_are_kept(self):
        block, changes = apply_desktop.plan_webapp({"webapp": {"extra": 1, "port": 9999}})
        self.assertEqual(block["extra"], 1)
        self.assertEqual(block["port"], 9120)
        self.assertIn(("port", 9999, 9120), changes)

    def test_off_sets_only_enabled(self):
        on, _ = apply_desktop.plan_webapp({})
        block, changes = apply_desktop.plan_webapp({"webapp": on}, off=True)
        self.assertEqual(changes, [("enabled", True, False)])
        self.assertEqual(block["host"], "0.0.0.0")

    def test_string_true_is_not_true(self):
        _, changes = apply_desktop.plan_webapp({"webapp": dict(apply_desktop.WEBAPP, enabled="true")})
        self.assertEqual([k for k, _, _ in changes], ["enabled"])

    def test_non_mapping_is_replaced(self):
        block, changes = apply_desktop.plan_webapp({"webapp": "oops"})
        self.assertEqual(block, apply_desktop.WEBAPP)
        self.assertIn("webapp", [k for k, _, _ in changes])

    def test_defaults_file_matches(self):
        try:
            import yaml
        except ImportError:
            self.skipTest("PyYAML not installed")
        seed = yaml.safe_load((HERE.parent / "config.defaults.yaml").read_text())
        self.assertEqual(seed["webapp"], apply_desktop.WEBAPP)

    def test_webapp_is_wildcard_bound_with_tailnet_url(self):
        self.assertEqual(apply_desktop.WEBAPP["host"], "0.0.0.0")
        self.assertTrue(apply_desktop.WEBAPP["public_url"].endswith(".tail717ba0.ts.net"))
        self.assertTrue(apply_desktop.WEBAPP["public_url"].startswith("https://"))


if __name__ == "__main__":
    unittest.main()
