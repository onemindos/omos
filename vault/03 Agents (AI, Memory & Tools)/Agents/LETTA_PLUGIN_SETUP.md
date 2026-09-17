---
id: ui-02-20260111053925
title: "Letta Plugin Setup Reference"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Letta Plugin Setup Reference"
---

# Letta Plugin Setup Reference

> **Quick reference for the Letta plugin symlink setup**
> **Date:** 2025-12-29 | **Version:** 2.0.0

---

## Location

**Plugin Development:**
```
~/Development/obsidian-plugins/letta-obsidian
```

**Obsidian Sees It At:**
```
.obsidian/plugins/letta-ai-agent → (symlink to dev folder)
```

**Git Repos:**
- Plugin: `https://github.com/Zeus-Delacruz/letta-obsidian.git`
- Vault: `https://github.com/Zeus-Delacruz/One-Mind-OS.git`

---

## Quick Commands

### Update Plugin from GitHub
```bash
cd ~/Development/obsidian-plugins/letta-obsidian
git pull origin master
npm run build
# Reload Obsidian (Cmd+R)
```

### Make Changes to Plugin
```bash
cd ~/Development/obsidian-plugins/letta-obsidian
# Edit files
npm run build
# Reload Obsidian (Cmd+R)
git add .
git commit -m "Changes"
git push origin master
```

### Fix Broken Symlink
```bash
cd .obsidian/plugins
rm -f letta-ai-agent
ln -s ~/Development/obsidian-plugins/letta-obsidian letta-ai-agent
```

---

## How It Works

1. **Plugin code lives** in `~/Development/obsidian-plugins/letta-obsidian`
2. **Symlink points** from `.obsidian/plugins/letta-ai-agent` to that location
3. **Obsidian loads** from the symlink (sees it as normal plugin)
4. **Vault git ignores** the symlink (no conflicts)
5. **Plugin git tracks** the code independently

---

## Why This Setup?

- ✅ **No git conflicts** - Vault and plugin repos are separate
- ✅ **Easy updates** - Just `git pull` in plugin folder
- ✅ **Instant sync** - Changes appear immediately
- ✅ **Permanent location** - Won't break if things move

---

## Features (v2.0.0)

- 9 vault operations (view, create, edit, delete, search)
- HITL (Human-in-the-Loop) approval system
- Full vault control
- Persistent memory
- Memory block management
- Project-aware agent switching

---

## Full Documentation

See comprehensive setup guide:
```
~/Development/obsidian-plugins/letta-obsidian/VAULT_SETUP.md
```

---

## Troubleshooting

**Plugin not loading?**
1. Check symlink: `ls -l .obsidian/plugins/letta-ai-agent`
2. Rebuild: `cd ~/Development/obsidian-plugins/letta-obsidian && npm run build`
3. Hard restart Obsidian

**Git errors about plugin?**
- Plugin should be in `.gitignore`
- Run: `git rm --cached .obsidian/plugins/letta-ai-agent`

---

**Last Updated:** 2025-12-29
