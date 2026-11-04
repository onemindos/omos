"""Cheap source contracts; the dispatched image smoke is the runtime gate."""
import ast
from pathlib import Path
import re
import unittest

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]


class ImageContract(unittest.TestCase):
    def test_exact_upstream_desktop_packages(self):
        # Keep this explicit: no xfce4/session/locker metapackage side effects.
        packages = {
            "tigervnc-standalone-server", "xfce4-panel", "xfwm4", "xfdesktop4",
            "xfce4-settings", "xfce4-terminal", "dbus-x11", "x11-xserver-utils",
            "x11-utils", "x11-xkb-utils", "xauth", "fonts-dejavu-core",
        }
        source = (HERE / "Dockerfile").read_text()
        apt = source.split("install -y --no-install-recommends", 1)[1].split("&&", 1)[0]
        self.assertEqual(set(apt.replace("\\", "").split()), packages)
        self.assertRegex(source, r"FROM nousresearch/hermes-agent:v0\.22\.5@sha256:[0-9a-f]{64}")
        self.assertIn('ensure("cua-driver", explicit=True)', source)
        self.assertNotIn("ENTRYPOINT", source)

    def test_manual_only_credentials_and_gate_before_publish(self):
        source = (ROOT / ".github/workflows/build-hermes-runtime.yml").read_text()
        self.assertIn("workflow_dispatch:", source)
        self.assertNotIn("pull_request:", source)
        self.assertNotIn("pull_request_target:", source)
        self.assertIn("default: runtime-fix", source)
        self.assertLess(source.index("/opt/bot-screen/smoke_test.py"), source.index("uses: docker/login-action"))
        for flag in ("--network none", "--memory 3g", "--user 10000:10000", "--security-opt no-new-privileges"):
            self.assertIn(flag, source)

    def test_smoke_is_valid_python_and_has_real_boundaries(self):
        source = (HERE / "smoke_test.py").read_text()
        ast.parse(source)
        for boundary in ('runtime.start(', 'runtime.stop()', 'client.connect(', 'client.recv(',
                         '"xdpyinfo"', 'browser.dock_argv(', 'DevToolsActivePort', 'os.getuid() == 10000'):
            self.assertIn(boundary, source)
        self.assertNotIn("min_free_memory_mb", source)
        self.assertNotIn("mock", source)


if __name__ == "__main__":
    unittest.main()
