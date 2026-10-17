"""Verify pinned foundation chart bytes and render the production values."""
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile
import urllib.request

ROOT = Path(__file__).resolve().parents[1]


def main():
    charts = json.loads((ROOT / "platform/ecosystem/foundation-lock.json").read_text())
    with tempfile.TemporaryDirectory(prefix="onemind-chart-check-") as tmp:
        for chart in charts:
            with urllib.request.urlopen(chart["url"], timeout=90) as response:
                data = response.read()
            if hashlib.sha256(data).hexdigest() != chart["sha256"]:
                raise RuntimeError("Chart checksum mismatch: " + chart["name"])
            archive = Path(tmp) / (chart["name"] + ".tgz")
            archive.write_bytes(data)
            subprocess.run([
                "helm", "template", chart["name"], str(archive),
                "--namespace", chart["namespace"], "--include-crds",
                "-f", str(ROOT / "platform/ecosystem/values" / (chart["name"] + ".yaml")),
            ], check=True, stdout=subprocess.DEVNULL)
            print("Verified and rendered " + chart["name"] + " " + chart["version"])
    core = ROOT / "platform/ecosystem/charts/onemind-core"
    if core.exists():
        subprocess.run(["helm", "lint", str(core)], check=True)
        subprocess.run(["helm", "template", "onemind-core", str(core), "--namespace", "onemind-core"],
                       check=True, stdout=subprocess.DEVNULL)


if __name__ == "__main__":
    main()
