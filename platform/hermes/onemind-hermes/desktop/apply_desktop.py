"""Idempotent fixes for the live onemind-hermes volume that Hermes Desktop trips over.

    python3 apply_desktop.py [--dry-run] [--home /opt/data]

Runs in the pod as root (it changes file ownership) and uses only the standard library.

What it fixes (found in the 2026-11-04 Desktop audit, see ../README.md "Hermes Desktop: feature status"):

- `/opt/data/onemind-skills` was copied from a Mac with tar: every entry belongs to uid 501 and the tree is full of
  AppleDouble `._*` files. The agent runs as `hermes` (uid 10000), so `skill_manage` failed with
  `PermissionError` on those external skills (and Desktop's Skills page cannot edit them). The applier gives the
  tree to `hermes` and removes the AppleDouble and `.DS_Store` files. Removed files are first copied into
  `/opt/data/backups/desktop-<UTC stamp>/appledouble.tar`.
- The same `._*` files sitting directly in `/opt/data` (for example `._SOUL.md`) are removed the same way.

A second run prints `unchanged`. It never reads or writes `config.yaml`, `auth.json` or any `.env`.
"""
import argparse
import datetime
import os
import pwd
import tarfile

AGENT_USER = "hermes"
JUNK_NAMES = {".DS_Store"}


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


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--dry-run", action="store_true")
    ap.add_argument("--home", default=os.environ.get("HERMES_HOME", "/opt/data"))
    args = ap.parse_args()

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
