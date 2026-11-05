"""Cheap source contracts; the dispatched image smoke is the runtime gate."""
import ast
from pathlib import Path
import re
import os
import subprocess
import tempfile
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

    def test_private_activation_patch_is_inert_and_narrow(self):
        relative = Path("platform/hermes/onemind-hermes/deployment.yaml")
        original = (ROOT / relative).read_text()
        upstream = "nousresearch/hermes-agent:v0.22.5@sha256:9c5561579f3ffdf333401797094ac3a743738d6c25d40f4331c7bf77b694238d"
        private = "ghcr.io/onemindos/hermes-agent@sha256:7e02bd636a1a6d66276bd006581eec72009aaa582af22a48605d8eae0bed9693"
        expected = original.replace(
            "      automountServiceAccountToken: false\n",
            "      automountServiceAccountToken: false\n"
            "      imagePullSecrets:\n        - name: onemind-hermes-registry\n", 1,
        ).replace("        - name: hermes\n          image: " + upstream,
                  "        - name: hermes\n          image: " + private, 1)
        with tempfile.TemporaryDirectory(dir=os.environ.get("TMPDIR")) as folder:
            target = Path(folder) / relative
            target.parent.mkdir(parents=True)
            target.write_text(original)
            subprocess.run(["git", "apply", str(HERE / "activation.patch")],
                           cwd=folder, check=True, capture_output=True)
            self.assertEqual(target.read_text(), expected)
        self.assertNotIn(private, original)
        self.assertNotIn("onemind-hermes-registry", original)
        kustomization = (HERE.parent / "kustomization.yaml").read_text()
        self.assertNotIn("activation.patch", kustomization)
        self.assertNotIn("bot-screen", kustomization)


if __name__ == "__main__":
    unittest.main()
