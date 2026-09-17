---
title: "OneMind OS — Dotfiles Management"
type: doctrine
status: planned
created: 2026-06-05
owner: zeus
tags: [dotfiles, bootstrap, portability, environment, setup]
---

# OneMind OS — Dotfiles Management

> Goal: The vault is your entire environment. Clone it on any machine and be back to 100% in minutes.

---

## The Problem

Every tool on macOS creates a `~/.toolname/` folder on first run. These folders live on the machine — not in the vault. When you move to a new Mac, they're gone. You lose config, skills, flows, and settings.

---

## The Solution

Move the config folders **into the vault**, then create **symlinks** from `~/` back to them. Tools still find everything at `~/.toolname` — nothing breaks. But now the config lives in git, travels with the vault, and can be restored on any machine instantly.

**Pattern (same as `.obsidian/`):**
```
~/onemind/.hermes/          ← actual data, in vault, versioned
~/.hermes → ~/onemind/.hermes   ← symlink, tool finds it here
```

---

## Dotfiles Philosophy

This is a standard practice called **dotfiles management**. Tools like `chezmoi`, `stow`, and `dotbot` all do this — but OneMind OS does it natively via the vault + a bootstrap script. No extra dependencies. Fully sovereign.

**What to version:**
- ✅ Config files, settings, skills, flows, key configs
- ❌ Session data, caches, binaries, databases, credentials

---

## Folders to Move into Vault

| Folder | What to keep | Gitignore |
|--------|-------------|-----------|
| `.hermes/` | `config.yaml`, skills, plugins, memories | sessions, audio cache, state.db |
| `.openclaw/` | `openclaw.json`, agent configs | sessions, trajectory files |
| `.nats/` | `nats.conf` | credentials |
| `.node-red/` | `settings.js`, flows | `node_modules/` |
| `.aws/` | `config` | `credentials` |
| `.ssh/` | `config`, `known_hosts` | private keys (`id_*`) |
| `.vscode/` | `settings.json`, `extensions.json` | — |

---

## Folders That Stay at Root (Never Move)

| Folder | Reason |
|--------|--------|
| `.cache/`, `.thumbnails/` | Generated cache — machine-specific |
| `.docker/` | Container runtime data |
| `.npm/`, `.nvm/`, `.bun/` | Package/version managers — binaries |
| `.dspy_cache/` | ML cache |
| `.agentmemory/`, `.iii/` | Live databases |
| `.codex/`, `.claude/`, `.copilot/`, `.gemini/` | Agent session data |
| `.zsh_sessions/` | Shell history |
| `.ssh/` keys | Private keys — never in git |
| `.azure/`, `.aws/credentials` | Credentials — never in git |
| `.secrets/` | Secrets — never in git |

---

## Bootstrap Script (Planned)

A single script at `_system/bootstrap.sh` that:

1. Checks for Homebrew, installs if missing
2. `git clone git@github.com:onemindos/One-Mind.git ~/onemind`
3. Creates symlinks for every dotfolder
4. Installs tools (nvm, bun, ollama, etc.)
5. Restores full OneMind OS environment

```bash
# Future usage — new Mac setup:
curl -fsSL https://raw.githubusercontent.com/onemindos/One-Mind/main/_system/bootstrap.sh | bash
```

One command. Full environment. Any machine.

---

## Status

- [ ] Move `.hermes/` into vault + symlink
- [ ] Move `.openclaw/` into vault + symlink
- [ ] Move `.nats/` into vault + symlink
- [ ] Move `.node-red/` into vault + symlink
- [ ] Move `.aws/config` into vault + symlink
- [ ] Move `.ssh/config` into vault + symlink
- [ ] Move `.vscode/` into vault + symlink
- [ ] Write `_system/bootstrap.sh`
- [ ] Test full restore on a clean machine

---

**Last updated:** 2026-06-05
**Owner:** Zeus
**Agent:** Legacy AI
