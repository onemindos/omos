---
name: mac-repo-hygiene
description: "Audit Zeus's Mac code folders for duplicate repos and bloat, verify repos are pushed, and rewire remotes to the onemindos GitHub org."
version: 2.0.0
author: Legacy
tags: [github, mac, repos, cleanup, onemindos, dedup]
---

# Mac Repo Hygiene

Inventory repos on Zeus's Mac, prove they are safe (pushed, clean), report bloat,
and move work into the `onemindos` GitHub org through pull requests.

## Reach and authority (read first)

- The cloud agent reaches the Mac only over Tailscale/SSH while the Mac is online.
  Resolve it from `tailscale status`; do not hard-code addresses. Offline means
  report and stop. The Mac is a leaf; Git on GitHub is the source of truth.
- Inventory, `git status`, `git fetch`, `du` and hashing are read-only and need no approval.
- Deleting anything (repos, `node_modules`, binaries, logs, launchd entries, vault
  folders, empty-folder cleanup, `rm`, `remove_tree`, emptying Trash) needs Zeus to
  write "execute" or "override" in chat for that specific deletion (SOUL rule 4).
  A "safe to delete" verdict is a proposal. Deduplicate by byte-verifying and
  quarantining (move to a clearly named quarantine folder Zeus can inspect), never
  by blind deletion.
- Pushing or creating repos and rewriting remotes change shared state: do them on a
  branch with a PR, and force pushes only on Zeus's explicit word.
- Stopping or unloading services on the Mac (launchd, `nats`, `nex`) is an action
  with effect outside software: propose it and wait for "execute" or "override".

## 1. Inventory a folder

```bash
for d in "$BASE"/*/; do
  name="$(basename "$d"); size=$(du -sh "$d" 2>/dev/null | cut -f1)
  if [ -d "$d/.git" ]; then
    remote=$(git -C "$d" remote get-url origin 2>/dev/null || echo NO-REMOTE)
    dirty=$(git -C "$d" status --short 2>/dev/null | wc -l | tr -d ' ')
    echo "GIT | $size | $name | $remote | dirty: $dirty"
  else
    echo "NO GIT | $size | $name"
  fi
done
```

Run it over SSH from a script file on the Mac with a quoted heredoc
(`<< 'EOF'`), so the pod does not expand `$`. Quote paths: several folders have a
trailing space (for example the media folder and the Desktop `atak ` folder).

## 2. Batch sync-check before anyone proposes deleting a local repo

For each repo: `git fetch origin`, compare `HEAD` with `origin/<branch>`, list
`git status --porcelain` (ignoring `.DS_Store`), `git stash list`, and
`git log @{u}..HEAD`. A repo is re-clonable only if HEAD equals remote, the tree is
clean, no stash, no unpushed commits. Present a table: repo, remote, sync,
uncommitted, size, verdict. The verdict is a proposal.

## 3. Orphan and non-git folders

- No remote and only scaffold commits ("Initial commit from Create Next App") with
  only generated files: likely nothing unique. Still propose, do not delete.
- A non-git folder is not automatically safe. If it holds markdown, architecture
  notes or planning docs, read it first; it may be the only copy. Offer to push it
  to a repo.
- If an unrecognized repo has a `.env`, report the key names (not values) and
  advise Zeus to rotate any keys that are his. Never print or store the values.

## 4. Remotes, org repos and identity

- Remote form: `git@github.com:onemindos/<repo>.git`. Repo creation in the org
  goes through the approved Composio GitHub tool or `gh`, with Zeus's approval;
  see the `composio` skill. Look up the connected account at run time.
- GH007 "would publish a private email" means the commit author is a personal
  address. Set the commit email to the GitHub noreply address for Zeus's account
  (ask Zeus for it; do not guess or store it here), and re-author unpushed commits.
- Merges are never squashed. Zeus approves and merges PRs.
- Add `.DS_Store` to `.gitignore` and `git rm --cached` it.

## 5. Bloat report (no cleanup without approval)

`du -sh "$BASE"/*/ | sort -rh`. Typical offenders: `node_modules/` (50 MB to 1 GB;
on deep pnpm trees `rm -rf` silently fails, so use perl `remove_tree`), compiled
Go binaries, `logs/`. A folder that keeps coming back is usually a running
process's working directory (check `lsof -p <pid>`, `launchctl list`).

## 6. NATS on the Mac

The Mac runs a leaf node that forwards to the cluster. It is a forwarding pipe, not
a storage node, and not the source of truth. JetStream must not run there: past
runs accumulated gigabytes of duplicate stream data the cluster already holds.
If Mac JetStream data reappears, report it and propose removing the `jetstream`
block from the leaf config and purging the directory; stopping the service and
the purge each need "execute" or "override". Verify afterwards that the leaf is
still connected (`/leafz` on the local monitor port).

## 7. Vault and external-drive audits

See `references/obsidian-vault-dedup.md` for the byte-by-byte comparison script
and folder-by-folder workflow. Rules: go one folder at a time and advance only
when Zeus names the next folder; print only the summary and the list of unique
files, never every duplicate line; identity documents are always unique and are
moved to the asset drive only on Zeus's word; before adding new binary or
sensitive patterns to a vault repo, update its `.gitignore` first (once in Git
history, removal needs history rewriting). Keep vault changes from triggering
cluster deploys (exclude them from any Argo path).

## Pitfalls

- macOS moves "deleted" folders to `~/.Trash`; check there if something vanished.
- SSH heredocs without quotes expand `$` pod-side.
- `~/go/bin` can hold stale binaries shadowed by Homebrew; check `ps aux` for the
  real path before assuming what is live.
- Go module cache (`~/go/pkg/mod`) is regenerable, but deleting still needs the word.
