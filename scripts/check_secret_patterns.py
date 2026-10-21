"""Report possible hardcoded password locations without printing their values."""
from pathlib import Path
import re
import sys

found = False
for path in Path("platform").rglob("*"):
    if path.suffix not in (".yaml", ".yml") or not path.is_file():
        continue
    for number, line in enumerate(path.read_text().splitlines(), 1):
        if "#" not in line and re.search(r"\bpassword\s*:", line):
            if not any(s in line for s in ("secretKeyRef", "external-secret", "valueFrom")):
                print(f"{path}:{number}: review password configuration")
                found = True
sys.exit(0 if found else 1)
