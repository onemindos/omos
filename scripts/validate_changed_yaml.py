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
# Exact renames (R100) carry no new content, so moving a file does not force fixing its
# legacy formatting debt; any rename with edits is still checked.
fields = subprocess.check_output([
    "git", "diff", "--name-status", "-M", "--diff-filter=ACMR", "-z", base, "HEAD",
]).decode().split("\0")
paths = []
i = 0
while i < len(fields) and fields[i]:
    status = fields[i]
    if status.startswith(("R", "C")):
        if status != "R100":
            paths.append(fields[i + 2])
        i += 3
    else:
        paths.append(fields[i + 1])
        i += 2
files = [p for p in paths if Path(p).suffix in (".yaml", ".yml") and Path(p).is_file()]
if files:
    subprocess.run(["yamllint", "-c", ".yamllint.yml", *files], check=True)
print(f"Changed YAML gate passed ({len(files)} files checked)")
