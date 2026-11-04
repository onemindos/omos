"""Idempotent installer for Legacy's NATS presence on onemind-hermes. Runs inside the pod as `hermes`.

    python3 apply_nats.py [--dry-run] [--status] [--disable] [--with-context] [--src /opt/data/nats-defs]

Installs, from --src (pushed by nats/apply.sh), the `nats-platform` plugin (and with --with-context the
`nats-context` plugin) into $HERMES_HOME/plugins and lists them in `plugins.enabled`. Nothing else in
config.yaml is touched. The connection settings (URL, credentials, CA, identity) are not written here: they come
from the Deployment environment, so Git stays the source of truth for them.

- Anything it would overwrite is first copied to $HERMES_HOME/backups/legacy-nats-<UTC stamp>/.
- A file or key that already matches is left alone and reported "unchanged"; a second run changes nothing.
- --disable removes both plugins from `plugins.enabled` (the files stay). That is the rollback switch.
- --status changes nothing: it reports the plugin files, the enabled list, and which NATS variables are present
  in this pod (names only, never values).
- No credential is read, copied, printed or written.
"""

import datetime
import hashlib
import os
import pathlib
import shutil
import sys

sys.path.insert(0, "/opt/hermes")
from hermes_cli.config import atomic_config_write, read_user_config_raw  # noqa: E402

ARGS = set(sys.argv[1:])
DRY = "--dry-run" in ARGS
STATUS = "--status" in ARGS
DISABLE = "--disable" in ARGS
WITH_CONTEXT = "--with-context" in ARGS
SRC = sys.argv[sys.argv.index("--src") + 1] if "--src" in sys.argv else "/opt/data/nats-defs"
HOME = os.environ.get("HERMES_HOME", "/opt/data")
STAMP = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%dT%H%M%SZ")
BACKUP = os.path.join(HOME, "backups", "legacy-nats-" + STAMP)
DEPS = os.path.join(HOME, ".local", "gateway-deps")
PLUGINS = ["nats-platform"] + (["nats-context"] if WITH_CONTEXT else [])
ALL_PLUGINS = ["nats-platform", "nats-context"]
# Names only. Their presence is what the plugin needs; the values are never read here.
ENV_NAMES = [
    "NATS_URL", "NATS_USER", "NATS_PASSWORD", "NATS_CA_FILE", "NATS_TLS_SERVER_NAME",
    "HERMES_NATS_AGENT", "HERMES_NATS_OWNER", "HERMES_NATS_SESSION_NAME",
    "HERMES_NATS_REFUSE_DUPLICATE", "NATS_ALLOW_ALL_USERS",
]
rc = 0


def say(action, what):
    print("%s %s" % (("WOULD " if DRY and action != "unchanged" else "") + action, what))


def digest(path):
    try:
        with open(path, "rb") as fh:
            return hashlib.sha256(fh.read()).hexdigest()
    except FileNotFoundError:
        return None


def backup(path):
    if DRY or not os.path.exists(path):
        return
    dest = os.path.join(BACKUP, os.path.relpath(path, HOME))
    os.makedirs(os.path.dirname(dest), exist_ok=True)
    shutil.copy2(path, dest)


def source_files(name):
    root = os.path.join(SRC, "plugins", name)
    for base, dirs, files in os.walk(root):
        dirs[:] = [d for d in dirs if d != "__pycache__"]
        for f in sorted(files):
            if f.startswith("."):
                continue
            full = os.path.join(base, f)
            yield os.path.relpath(full, root), full


def install_plugin(name):
    for rel, full in source_files(name):
        dest = os.path.join(HOME, "plugins", name, rel)
        if digest(dest) == digest(full):
            say("unchanged", "plugin %s/%s" % (name, rel))
            continue
        say("update" if os.path.exists(dest) else "create", "plugin %s/%s" % (name, rel))
        if DRY:
            continue
        backup(dest)
        os.makedirs(os.path.dirname(dest), exist_ok=True)
        tmp = dest + ".tmp"
        shutil.copyfile(full, tmp)
        os.replace(tmp, dest)


def set_enabled(wanted_plugins):
    """Make `plugins.enabled` contain exactly our plugins that are wanted, keeping every other entry."""
    path = os.path.join(HOME, "config.yaml")
    cfg = read_user_config_raw(pathlib.Path(path))
    if not isinstance(cfg.get("plugins"), dict):
        cfg["plugins"] = {}
    current = list(cfg["plugins"].get("enabled") or [])
    new = [p for p in current if p not in ALL_PLUGINS] + list(wanted_plugins)
    if new == current or sorted(new) == sorted(current):
        say("unchanged", "config.yaml plugins.enabled = %s" % current)
        return
    say("update", "config.yaml plugins.enabled: %s -> %s" % (current, new))
    if DRY:
        return
    backup(path)
    cfg["plugins"]["enabled"] = new
    atomic_config_write(pathlib.Path(path), cfg)


def report_status():
    print("pod environment (names only):")
    for name in ENV_NAMES:
        print("  %-30s %s" % (name, "set" if os.environ.get(name) else "MISSING"))
    ca = os.environ.get("NATS_CA_FILE", "")
    print("  CA file readable: %s" % (os.access(ca, os.R_OK) if ca else False))
    marker = os.path.join(DEPS, ".requirements.sha256")
    print("gateway deps: %s (%s)" % (DEPS, "installed" if os.path.exists(marker) else "NOT INSTALLED"))
    for name in ALL_PLUGINS:
        files = list(source_files(name)) if os.path.isdir(os.path.join(SRC, "plugins", name)) else []
        same = [digest(os.path.join(HOME, "plugins", name, rel)) == digest(full) for rel, full in files]
        state = "absent" if not os.path.isdir(os.path.join(HOME, "plugins", name)) else (
            "matches Git" if files and all(same) else "differs from the staged copy")
        print("plugin %-13s %s" % (name, state))
    cfg = read_user_config_raw(pathlib.Path(os.path.join(HOME, "config.yaml")))
    print("plugins.enabled: %s" % ((cfg.get("plugins") or {}).get("enabled") or []))


def main():
    global rc
    if STATUS:
        report_status()
        return 0
    if DISABLE:
        set_enabled([])
        print("Disabled. Restart the gateway (nats/apply.sh --restart) for it to take effect.")
        return rc
    for name in PLUGINS:
        install_plugin(name)
    if not os.path.exists(os.path.join(DEPS, ".requirements.sha256")):
        print("ERROR gateway deps are not installed in %s (the install-gateway-deps init container failed or has not run)" % DEPS)
        rc = 1
    for name in ("NATS_URL", "NATS_USER", "NATS_PASSWORD", "NATS_CA_FILE", "HERMES_NATS_SESSION_NAME"):
        if not os.environ.get(name):
            print("ERROR %s is not set in this pod: the Deployment change is not rolled out yet" % name)
            rc = 1
    if rc:
        print("Not enabling the plugin: fix the errors above first (nothing is switched on).")
        return rc
    set_enabled(PLUGINS)
    print("Backups of anything overwritten: %s" % (BACKUP if os.path.isdir(BACKUP) else "none needed"))
    print("Enabled in config. Restart the gateway (nats/apply.sh --restart) for it to load.")
    return rc


if __name__ == "__main__":
    sys.exit(main())
