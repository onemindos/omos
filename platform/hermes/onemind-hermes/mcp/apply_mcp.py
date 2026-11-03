"""Idempotent applier for the `mcp_servers` section of the live Hermes config.yaml.
Runs inside the onemind-hermes pod, as the `hermes` user (mcp/apply.sh starts it).

    python3 apply_mcp.py [--dry-run] [--defaults config.defaults.yaml] [--config /opt/data/config.yaml]

config.defaults.yaml is only copied into an EMPTY volume, so the live config must be changed in
place. The desired `mcp_servers` entries are read from the defaults file (one source of truth).

- Every server in the defaults is set in the live config exactly as there (replaced if it differs).
- Servers in REMOVE are deleted from the live config.
- Other servers and every other key are left alone; comments and layout are kept (round-trip YAML).
- Before the first change a copy is written next to the config: config.yaml.bak-mcp-<UTC timestamp>.
- A second run with no drift writes nothing and prints only "unchanged" lines.
"""

import datetime
import os
import shutil
import sys

sys.path.insert(0, "/opt/hermes")
from ruamel.yaml import YAML  # noqa: E402  (PyYAML is not in Hermes' venv)

# Retired: the LiteLLM MCP hub entry (replaced by the tool gateway, Phase 5).
REMOVE = ["onemind"]


def arg(flag, default):
    return sys.argv[sys.argv.index(flag) + 1] if flag in sys.argv else default


DRY = "--dry-run" in sys.argv
DEFAULTS = arg("--defaults", "/tmp/mcp-apply/config.defaults.yaml")
CONFIG = arg("--config", "/opt/data/config.yaml")

yaml = YAML()  # round-trip: keeps comments and key order
yaml.preserve_quotes = True
yaml.width = 4096


def plain(node):
    """Comparable plain python value of a (possibly round-trip) node."""
    if hasattr(node, "items"):
        return {str(k): plain(v) for k, v in node.items()}
    if isinstance(node, (list, tuple)):
        return [plain(v) for v in node]
    return node


def main():
    with open(DEFAULTS, encoding="utf-8") as fh:
        desired = (YAML(typ="safe").load(fh) or {}).get("mcp_servers") or {}
    if not desired:
        print("ERROR: no mcp_servers in %s" % DEFAULTS)
        return 1
    with open(CONFIG, encoding="utf-8") as fh:
        conf = yaml.load(fh)
    live = conf.get("mcp_servers")
    if live is None:
        live = conf["mcp_servers"] = {}
    changed = False
    for name, spec in desired.items():
        if name in live and plain(live[name]) == plain(spec):
            print("unchanged %s" % name)
            continue
        print("%s %s: %s" % ("WOULD" if DRY else "DO", name, "update" if name in live else "add"))
        live[name] = spec
        changed = True
    for name in REMOVE:
        if name in live and name not in desired:
            print("%s %s: remove" % ("WOULD" if DRY else "DO", name))
            del live[name]
            changed = True
        else:
            print("unchanged %s (absent)" % name)
    unmanaged = sorted(set(live) - set(desired) - set(REMOVE))
    if unmanaged:
        print("unmanaged, left alone: %s" % ", ".join(unmanaged))
    if not changed or DRY:
        return 0
    stamp = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    backup = "%s.bak-mcp-%s" % (CONFIG, stamp)
    shutil.copy2(CONFIG, backup)
    print("backup %s" % backup)
    staged = CONFIG + ".new"
    with open(staged, "w", encoding="utf-8") as fh:
        yaml.dump(conf, fh)
    shutil.copymode(CONFIG, staged)
    os.replace(staged, CONFIG)
    print("written %s (restart Hermes so it reloads mcp_servers)" % CONFIG)
    return 0


sys.exit(main())
