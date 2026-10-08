"""Enforce YAML lint for changed files while legacy formatting debt is reported."""
import os
from pathlib import Path
import re
import subprocess

base = os.environ.get("YAML_BASE_SHA", "")
if not base or set(base) == {"0"}:
    base = subprocess.check_output(["git", "rev-parse", "HEAD^"], text=True).strip()
if not re.fullmatch(r"[0-9a-f]{40}", base):
    raise SystemExit("Expected a commit SHA for YAML validation")
paths = subprocess.check_output([
    "git", "diff", "--name-only", "--diff-filter=ACMR", "-z", base, "HEAD",
]).decode().split("\0")
files = [p for p in paths if Path(p).suffix in (".yaml", ".yml") and Path(p).is_file()]
if files:
    subprocess.run(["yamllint", "-c", ".yamllint.yml", *files], check=True)
print(f"Changed YAML gate passed ({len(files)} files checked)")
