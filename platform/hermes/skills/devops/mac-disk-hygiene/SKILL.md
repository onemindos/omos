---
name: mac-disk-hygiene
description: "Audit or free disk space on Zeus's Mac over SSH. Read-only audit by default; every deletion needs execute/override."
version: 2.0.0
author: Legacy
---

# Mac Disk Hygiene

Audit what is using space on Zeus's MacBook Pro and propose cleanups. The Mac is
a 1 TB working machine, not an archive: big media belongs on the external
`OM PRIMARY` drive and in Dropbox, and the internal disk should stay well over
half free.

## Reach and authority (read first)

- The cloud agent reaches the Mac only over Tailscale/SSH, and only while the Mac
  is online and awake. Resolve it from `tailscale status` (MagicDNS name); do not
  hard-code an IP. If it is offline, say so and stop. The Mac is a leaf, never the
  source of truth.
- Auditing (`df`, `du`, `git status`) is read-only and needs no approval.
- Deleting anything (cache purge, repo checkout, models, installers, `brew
  cleanup`) needs Zeus to write "execute" or "override" in chat for that specific
  deletion (SOUL rule 4). Earlier approvals, "clean it all up" and verified-safe
  results do not substitute. Present a table (item, size, risk, how to restore),
  then wait.
- Ejecting drives, sleep, or power changes are physical-world actions under SOUL
  rule 8 and need the same word.

## SSH route

`ssh -i /opt/data/home/.ssh/<key> zeus@<mac-tailnet-name> '<cmd>'`. Prepend
`export PATH=/opt/homebrew/bin:/usr/local/bin:$PATH` for brew tools. Pin the host
key after the first verified connection instead of disabling host checking.

## Audit procedure

1. `df -h /System/Volumes/Data` for the headline, then depth-capped `du` per
   suspect directory. Never run bare `du` on a huge tree; use
   `du -xg -d1 <dir> | sort -rn | head -N`.
2. Categorize before proposing anything: (a) regenerable caches and build
   artifacts, (b) tool installs (editors, engines), reinstallable, (c) original or
   unique data. Class (c) is never proposed for deletion, only for relocation.

## Verify a repo checkout is a clean clone (read-only)

`git remote -v`, then `git log --oneline @{u}..HEAD | wc -l` (unpushed commits)
and `git ls-files --others --exclude-standard` (untracked). Zero unpushed and only
`.DS_Store`-class untracked means it can be re-cloned. That is evidence for a
proposal, not permission.

## Regenerable set to propose (never run unprompted)

`~/.npm`, `~/.gradle`, `~/.buildozer`, `~/.python-for-android`, `~/.cache`,
`~/Library/Caches/*`, stale `com.docker.install` installers, then
`brew cleanup -s`. All regenerate on demand. Each deletion still needs its own
"execute" or "override".

## Pitfalls

- macOS `df` under-credits deletions while local Time Machine snapshots hold the
  blocks (hours to days). Report the purgeable gap; do not re-audit in a panic.
- Docker Desktop's VM disk (`~/Library/Containers/com.docker.docker/Data/vms`):
  confirm images are re-pullable first.
- LM Studio models (`~/.lmstudio/models`) are large and user-chosen; ask.
- Dropbox big media must stay Mirror/local, never online-only: a local rsync or
  rclone would force-hydrate placeholders.
- `/private/var/vm` is system-managed. Leave it.
- Long SSH jobs: run in the background on the Mac with a log in `/tmp` and an exit
  marker; foreground-launched jobs die when the SSH session drops.
