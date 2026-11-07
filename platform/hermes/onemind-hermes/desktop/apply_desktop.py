"""Idempotent fixes for the live onemind-hermes volume that Hermes Desktop trips over.

    python3 apply_desktop.py [--step files|webapp] [--dry-run] [--webapp-off] [--home /opt/data]

Two steps, run by desktop/apply.sh in this order:

- `--step files` (the default) runs in the pod as root (it changes file ownership) and uses only the standard
  library. It never reads or writes `config.yaml`, `auth.json` or any `.env`.
- `--step webapp` runs as the `hermes` user inside Hermes' own Python environment (it edits `config.yaml`
  through Hermes' atomic writer) and sets only the `webapp:` block, see below.

Step `files` fixes (found in the 2026-11-04 Desktop audit, see ../README.md "Hermes Desktop: feature status"):

- `/opt/data/onemind-skills` was copied from a Mac with tar: every entry belongs to uid 501 and the tree is full of
  AppleDouble `._*` files. The agent runs as `hermes` (uid 10000), so `skill_manage` failed with
  `PermissionError` on those external skills (and Desktop's Skills page cannot edit them). The applier gives the
  tree to `hermes` and removes the AppleDouble and `.DS_Store` files. Removed files are first copied into
  `/opt/data/backups/desktop-<UTC stamp>/appledouble.tar`.
- The same `._*` files sitting directly in `/opt/data` (for example `._SOUL.md`) are removed the same way.

Step `webapp` (Hermes v0.23.0, 2026-11-07) turns on the browser Desktop, the Webapp, on its own port. It sets
`webapp.enabled`, `host`, `port` and `public_url` in `config.yaml` to WEBAPP below and leaves every other key,
including any other key inside `webapp:`, alone. `--webapp-off` sets only `webapp.enabled: false` (the rollback
switch). The gateway reads the block at start, so a change needs a gateway restart (`apply.sh --restart`).
`config.defaults.yaml` carries the same block for a fresh volume; a test keeps the two equal. Anything it
overwrites is first copied to `/opt/data/backups/desktop-<UTC stamp>/config.yaml`.

A second run prints `unchanged`.
"""
import argparse
import datetime
import os
import pathlib
import pwd
import shutil
import sys
import tarfile

AGENT_USER = "hermes"
JUNK_NAMES = {".DS_Store"}
# The Webapp listens on all interfaces of the pod (a wildcard bind turns on dashboard sign-in) and is served to
# the tailnet only, by the `onemind-hermes-web` Tailscale Ingress (service.yaml). public_url is that address.
WEBAPP = {
    "enabled": True,
    "host": "0.0.0.0",
    "port": 9120,
    "public_url": "https://onemind-hermes-web.tail717ba0.ts.net",
}


def is_junk(name):
    return name.startswith("._") or name in JUNK_NAMES


def scan(root, uid, gid, recurse):
    """Return (junk files, entries with the wrong owner) under root. Symlinks are never followed."""
    junk, wrong = [], []
    if not os.path.isdir(root):
        return junk, wrong
    if recurse:
        walker = os.walk(root)
    else:
        top = next(os.walk(root))
        walker = [(top[0], [], top[2])]
    for base, dirs, files in walker:
        for name in dirs + files:
            path = os.path.join(base, name)
            st = os.lstat(path)
            if name in files and is_junk(name):
                junk.append(path)
            elif st.st_uid != uid or st.st_gid != gid:
                wrong.append(path)
    return junk, wrong


def plan_webapp(config, off=False):
    """Pure: the `webapp:` block that `config` should end up with, and the (key, old, new) changes to get there.

    Keys that are not ours are kept. A `webapp` that is not a mapping is replaced. `off` sets only `enabled`.
    """
    current = config.get("webapp") if isinstance(config, dict) else None
    block = dict(current) if isinstance(current, dict) else {}
    wanted = {"enabled": False} if off else dict(WEBAPP)
    changes = []
    for key, value in wanted.items():
        if key not in block or block[key] != value or type(block[key]) is not type(value):
            changes.append((key, block.get(key), value))
            block[key] = value
    if current is not None and not isinstance(current, dict):
        changes.append(("webapp", current, "mapping"))
    return block, changes


def webapp_step(args):
    """Set `webapp:` in config.yaml with Hermes' own reader and atomic writer (needs /opt/hermes on the path)."""
    sys.path.insert(0, "/opt/hermes")
    from hermes_cli.config import atomic_config_write, read_user_config_raw

    path = pathlib.Path(args.home) / "config.yaml"
    config = read_user_config_raw(path)
    block, changes = plan_webapp(config, off=args.webapp_off)
    if not changes:
        print("unchanged: config.yaml webapp = %s" % {k: block[k] for k in sorted(block)})
        return 0
    verb = "would update" if args.dry_run else "update"
    for key, old, new in changes:
        print("%s config.yaml webapp.%s: %r -> %r" % (verb, key, old, new))
    if args.dry_run:
        return 0
    stamp = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    bdir = os.path.join(args.home, "backups", "desktop-" + stamp)
    os.makedirs(bdir, exist_ok=True)
    shutil.copy2(str(path), os.path.join(bdir, "config.yaml"))
    print("backup: %s/config.yaml" % bdir)
    config["webapp"] = block
    atomic_config_write(path, config)
    print("done. Restart the gateway for it to take effect (desktop/apply.sh --restart).")
    return 0


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--dry-run", action="store_true")
    ap.add_argument("--step", choices=("files", "webapp"), default="files")
    ap.add_argument("--webapp-off", action="store_true", help="webapp step: set webapp.enabled false")
    ap.add_argument("--home", default=os.environ.get("HERMES_HOME", "/opt/data"))
    args = ap.parse_args()
    if args.step == "webapp":
        return webapp_step(args)

    pw = pwd.getpwnam(AGENT_USER)
    uid, gid = pw.pw_uid, pw.pw_gid
    skills = os.path.join(args.home, "onemind-skills")

    s_junk, s_wrong = scan(skills, uid, gid, recurse=True)
    h_junk, _ = scan(args.home, uid, gid, recurse=False)
    junk = s_junk + h_junk
    wrong = s_wrong
    # The skills root itself, if it is not the agent's.
    if os.path.isdir(skills):
        st = os.lstat(skills)
        if st.st_uid != uid or st.st_gid != gid:
            wrong.append(skills)

    if not junk and not wrong:
        print("unchanged: onemind-skills is owned by %s and holds no AppleDouble files" % AGENT_USER)
        return 0

    verb = "would " if args.dry_run else ""
    print("%sremove %d AppleDouble/.DS_Store files (%d under onemind-skills, %d in %s)"
          % (verb, len(junk), len(s_junk), len(h_junk), args.home))
    print("%schown %d entries under onemind-skills to %s (uid %d)" % (verb, len(wrong), AGENT_USER, uid))
    if args.dry_run:
        return 0

    if junk:
        stamp = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%dT%H%M%SZ")
        bdir = os.path.join(args.home, "backups", "desktop-" + stamp)
        os.makedirs(bdir, exist_ok=True)
        with tarfile.open(os.path.join(bdir, "appledouble.tar"), "w") as tar:
            for path in junk:
                tar.add(path, arcname=os.path.relpath(path, args.home), recursive=False)
        os.chown(bdir, uid, gid)
        os.chown(os.path.join(bdir, "appledouble.tar"), uid, gid)
        print("backup: %s/appledouble.tar" % bdir)
        for path in junk:
            os.remove(path)
    for path in wrong:
        if os.path.lexists(path):
            os.lchown(path, uid, gid)
    print("done")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
