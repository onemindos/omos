"""Idempotent installer for Legacy's doctrine and worker roles. Runs inside the onemind-hermes pod as `hermes`.

    python3 apply_roles.py [--dry-run] [--src /opt/data/roles-defs]

Installs, from --src (pushed by roles/apply.sh): SOUL.md, the `legacy-roles` delegation plugin, one Hermes
profile per role in roles.yaml, and the default profile's `delegation.backend` / `plugins.enabled` keys.

- Anything it would overwrite is first copied to $HERMES_HOME/backups/legacy-roles-<UTC stamp>/.
- A file or key that already matches is left alone and reported "unchanged"; a second run changes nothing.
- No credential is read, copied or written. Role profiles get no auth.json: they inherit the default profile's
  pooled openai-codex credentials (hermes_cli.auth.read_credential_pool falls back to the root store).
- Config goes through Hermes' own atomic_config_write (comments kept, existing keys never dropped).
"""

import datetime
import pathlib
import os
import shutil
import subprocess
import sys

sys.path.insert(0, "/opt/hermes")
from ruamel.yaml import YAML  # noqa: E402  (PyYAML is not in Hermes' venv)
from hermes_cli.config import atomic_config_write, read_user_config_raw  # noqa: E402

DRY = "--dry-run" in sys.argv
SRC = sys.argv[sys.argv.index("--src") + 1] if "--src" in sys.argv else "/opt/data/roles-defs"
HOME = os.environ.get("HERMES_HOME", "/opt/data")
PLUGIN = "legacy-roles"
STAMP = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%dT%H%M%SZ")
BACKUP = os.path.join(HOME, "backups", "legacy-roles-" + STAMP)
rc = 0


def say(action, what):
    print("%s %s" % (("WOULD " if DRY and action != "unchanged" else "") + action, what))


def backup(path):
    """Copy an existing file into this run's backup dir, keeping its path relative to HERMES_HOME."""
    if DRY or not os.path.exists(path):
        return
    dest = os.path.join(BACKUP, os.path.relpath(path, HOME))
    os.makedirs(os.path.dirname(dest), exist_ok=True)
    shutil.copy2(path, dest)


def read(path):
    try:
        with open(path, encoding="utf-8") as fh:
            return fh.read()
    except FileNotFoundError:
        return None


def put_file(dest, text, label):
    """Write text to dest when it differs (backing up the old file)."""
    if read(dest) == text:
        say("unchanged", label)
        return
    say("update" if os.path.exists(dest) else "create", label)
    if DRY:
        return
    backup(dest)
    os.makedirs(os.path.dirname(dest), exist_ok=True)
    tmp = dest + ".tmp"
    with open(tmp, "w", encoding="utf-8") as fh:
        fh.write(text)
    os.replace(tmp, dest)


def put_config(path, wanted, label):
    """Set the given dotted keys in a config.yaml when they differ. wanted: {"a.b": value}."""
    cfg = read_user_config_raw(pathlib.Path(path))
    changed = []
    for dotted, value in wanted.items():
        node = cfg
        parts = dotted.split(".")
        for part in parts[:-1]:
            if not isinstance(node.get(part), dict):
                node[part] = {}
            node = node[part]
        if node.get(parts[-1]) != value:
            node[parts[-1]] = value
            changed.append(dotted)
    if not changed:
        say("unchanged", label)
        return
    say("update", "%s (%s)" % (label, ", ".join(changed)))
    if not DRY:
        backup(path)
        atomic_config_write(pathlib.Path(path), cfg)


def hermes(*args):
    exe = shutil.which("hermes") or "/opt/hermes/.venv/bin/hermes"
    return subprocess.run([exe, *args], capture_output=True, text=True, env=dict(os.environ, HERMES_HOME=HOME))


def main():
    global rc
    with open(os.path.join(SRC, "roles.yaml"), encoding="utf-8") as fh:
        spec = YAML(typ="safe").load(fh)
    defaults, roles = spec["defaults"], spec["roles"]
    default_cfg_path = os.path.join(HOME, "config.yaml")
    default_cfg = read_user_config_raw(pathlib.Path(default_cfg_path))

    # 1. Doctrine: SOUL.md is Git-owned (platform/hermes/seed/onemind-hermes/SOUL.md).
    put_file(os.path.join(HOME, "SOUL.md"), read(os.path.join(SRC, "SOUL.md")), "SOUL.md")

    # 2. The delegation backend plugin.
    plugin_src = os.path.join(SRC, "plugin", PLUGIN)
    for name in sorted(os.listdir(plugin_src)):
        if name.startswith((".", "__pycache__")):
            continue
        put_file(os.path.join(HOME, "plugins", PLUGIN, name), read(os.path.join(plugin_src, name)), "plugin %s/%s" % (PLUGIN, name))

    # 3. One profile per role. Created without bundled skills, an alias or any credential.
    for name, role in roles.items():
        from hermes_cli.profiles import get_profile_dir, profile_exists
        if not profile_exists(name):
            say("create", "profile %s" % name)
            if not DRY:
                done = hermes("profile", "create", name, "--no-alias", "--no-skills", "--description", role["description"].strip())
                if done.returncode != 0:
                    print("ERROR profile %s: %s" % (name, (done.stderr or done.stdout).strip()[-400:]))
                    rc = 1
                    continue
        pdir = str(get_profile_dir(name))
        if DRY and not os.path.isdir(pdir):
            continue
        wanted = {
            "model.provider": defaults["provider"],
            "model.default": role["model"],
            "credential_pool_strategies": {defaults["provider"]: defaults["pool_strategy"]},
            "platform_toolsets.cli": list(role["toolsets"]),
            "skills.external_dirs": list(defaults["skills_external_dirs"]),
        }
        # Same safety posture as the default profile: copied, never invented here.
        for key in ("approvals", "security"):
            if isinstance(default_cfg.get(key), dict):
                wanted[key] = default_cfg[key]
        if isinstance(default_cfg.get("agent"), dict) and default_cfg["agent"].get("disabled_toolsets"):
            wanted["agent.disabled_toolsets"] = list(default_cfg["agent"]["disabled_toolsets"])
        # Tool gateway: whatever the default profile has. ${VAR} references stay unexpanded (raw read).
        if isinstance(default_cfg.get("mcp_servers"), dict) and default_cfg["mcp_servers"]:
            wanted["mcp_servers"] = default_cfg["mcp_servers"]
        put_config(os.path.join(pdir, "config.yaml"), wanted, "profile %s config.yaml" % name)
        put_file(os.path.join(pdir, "SOUL.md"), role["brief"].strip() + "\n", "profile %s SOUL.md" % name)
        from hermes_cli.profiles import read_profile_meta
        want_desc = " ".join(role["description"].split())
        if " ".join(str(read_profile_meta(pathlib.Path(pdir)).get("description") or "").split()) != want_desc:
            say("update", "profile %s description" % name)
            if not DRY:
                hermes("profile", "describe", name, "--text", want_desc)
        else:
            say("unchanged", "profile %s description" % name)

    # 4. Default profile: select the backend and enable the plugin (existing entries kept).
    enabled = list((default_cfg.get("plugins") or {}).get("enabled") or [])
    if PLUGIN not in enabled:
        enabled.append(PLUGIN)
    put_config(default_cfg_path, {"delegation.backend": PLUGIN, "plugins.enabled": enabled}, "config.yaml delegation.backend, plugins.enabled")

    if not DRY:
        for root, dirs, files in os.walk(os.path.join(HOME, "plugins", PLUGIN)):
            for item in [root, *[os.path.join(root, f) for f in files]]:
                try:
                    shutil.chown(item, "hermes", "hermes")
                except (LookupError, PermissionError):
                    pass
    print("Backups of anything overwritten: %s" % (BACKUP if os.path.isdir(BACKUP) else "none needed"))
    return rc


if __name__ == "__main__":
    sys.exit(main())
