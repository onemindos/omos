---
title: "File System (MCP)"
type: service-card
status: active
created: 2026-06-08
updated: 2026-06-08
owner: zeus
category: tool
tags: [tool, mcp, filesystem, local, storage]
---

# File System (MCP)

> Native MCP file system tool — gives agents direct read/write access to local directories.

## Role in OneMind OS
- Agents read and write vault files, configs, scripts
- Scoped to approved directories only (vault, _system, etc.)
- Foundation for any agent that needs to persist or retrieve local data

## Config
- Interface: MCP tool (native Hermes)
- Scope: configured paths only — not full disk

## Use Cases
- Reading/writing vault markdown files
- Accessing configs and dotfiles
- Script output persistence
