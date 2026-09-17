---
title: "Symlink"
type: reference
domain: SO
status: active
created: 2026-05-27
updated: 2026-05-27
tags: [filesystem, os, infrastructure, dev, concept]
---

# Symlink

## Definition
A symbolic link (symlink) is a filesystem pointer that references another file or directory by path. The OS transparently resolves it — anything reading the symlink sees the target's content. Unlike a hard link, a symlink can cross filesystems and point to directories.

```bash
ln -s /actual/target/path /where/the/link/lives
```

## Key Properties

| Type | Description |
|------|-------------|
| Soft (symlink) | Points to a path — breaks if target moves/deleted |
| Hard link | Points to the inode — survives target rename, same FS only |
| Directory symlink | Link to a directory; `ls -la` shows `->` arrow |

## OneMind OS Angle
Used heavily for config management across the fleet — symlinks let Hermes profiles, NATS config files, and agent runtime configs live in the vault (`~/Documents/onemind/`) while being readable at the paths services expect. Key pattern: keep source of truth in vault, symlink into place.

Example uses:
- `~/.hermes/config.yaml` → vault copy under SO/03 Infrastructure
- Agent `.env` files symlinked from SO/04 Security references
- Multi-profile Hermes setups sharing a common skills/ directory

## Related Concepts
- [[infrastructure]] — symlinks are a core node management primitive
- [[nats-server]] — config files often managed this way

## Notes
<!-- Add specific symlink patterns used in OneMind OS here -->
