---
type: sop
status: in-progress
tags: [git, infrastructure, backup, forgejo, codeberg, github, version-control]
owner: zeus
created: 2026-06-03
---

# Git Multi-Remote Safety Net — Setup SOP

> STATUS: IN PROGRESS — Resume from "NEXT STEPS" section below.

One push from MBP hits three destinations simultaneously:
GitHub (primary) → Codeberg (cloud backup) → Forgejo (local self-hosted)

---

## Architecture

```
MBP (git push)
  └── origin (multi-remote)
       ├── github.com/onemindos/<repo>       ← cloud primary (Microsoft infra)
       ├── codeberg.org/onemindos/<repo>      ← cloud backup (EU non-profit)
       └── localhost:3000/onemindos/<repo>    ← local Forgejo (self-hosted, offline-capable)

Fleet Nodes (DIMOS, OM1, House, Barn, Office)
  └── git pull from Forgejo over Headscale mesh
       OR git-sync daemon for auto-pull on new commits
```

---

## Tools Installed (DONE)

| Tool     | Version  | Purpose                              | Status    |
|----------|----------|--------------------------------------|-----------|
| lazygit  | 0.62.1   | Terminal TUI git client              | INSTALLED |
| Gitify   | 6.20.0   | macOS menu bar — GitHub + Codeberg   | INSTALLED |
| Forgejo  | binary   | Self-hosted git platform             | NOT YET   |

---

## GitHub (DONE)

- Org renamed: OneMind-OS → onemindos
- All local remotes updated to github.com/onemindos/...
- Verified live with git ls-remote
- Note: One-Mind vault repo uses SSH remote (git@github.com:onemindos/One-Mind.git) — keep as-is

Affected repos and their current GitHub remote:

| Repo                   | GitHub Remote URL                                          | Notes                        |
|------------------------|------------------------------------------------------------|------------------------------|
| onemind-command        | https://github.com/onemindos/onemind-app.git               | repo name is onemind-app     |
| onemind-codex-live     | https://github.com/onemindos/onemind-codex-live.git        |                              |
| onemind-codex-template | https://github.com/onemindos/onemind-codex-template.git    |                              |
| onemind-dashboard      | none — needs GitHub repo created first                     | local dir doesn't exist yet  |
| One-Mind (vault)       | git@github.com:onemindos/One-Mind.git                      | SSH remote, include in backup|

Repos in ~/dev/ that are third-party forks (DO NOT add multi-remote):
- homepage, langfuse, paperclip, hermes-webui (upstream forks, not ours)
- Crucix, wwv-src (not yet — future, pending Zeus decision)

---

## Codeberg (PENDING — waiting on Zeus)

Platform: codeberg.org (free, non-profit, EU-hosted, Forgejo-based)

### Steps to complete:
1. Create account at codeberg.org
2. Create org named: onemindos
3. Go to codeberg.org/user/settings/applications
4. Generate API token — name it "Legacy CLI"
   - Scopes: repository (read+write), organization (read+write)
5. Paste token to Legacy — I will run the rest automatically

### What Legacy will do with the token:
- Create all 4 repos under onemindos org via Codeberg API
- Add codeberg.org/onemindos/<repo> as second push remote on all 4 repos
- Test push to verify

---

## Forgejo — Local Self-Hosted (NOT STARTED)

Platform: Forgejo native binary (no Docker, no container)
Install method: Homebrew (v15.0.3 available — no macOS binary from Forgejo directly, brew is the clean path)
Location: Run directly on MBP as a LaunchAgent background service
Data dir: ~/forgejo-data
Web UI: http://localhost:3000
Access from fleet: http://forgejo.mesh.onemind.dev (via Headscale)

### Install steps (Legacy will run these):
```bash
# 1. Install via Homebrew (macOS arm64 native)
brew install forgejo

# 2. Create data directory
mkdir -p ~/forgejo-data

# 3. First launch — completes setup via web UI at http://localhost:3000
#    Choose: SQLite (no extra DB needed), set admin user, org name onemindos
forgejo web --work-path ~/forgejo-data

# 4. Install as LaunchAgent (auto-start on login)
# Legacy will write the plist to ~/Library/LaunchAgents/dev.forgejo.plist
# with: forgejo web --work-path /Users/onemind/forgejo-data
```

### What Legacy will do after install:
- Write and load LaunchAgent plist for auto-start
- Create onemindos org on local Forgejo via API
- Create all 5 repos under that org (onemind-command, onemind-codex-live, onemind-codex-template, onemind-dashboard, One-Mind)
- Add http://localhost:3000/onemindos/<repo>.git as third push remote on all repos
- Verify push works

---

## Multi-Remote Final State (per repo)

Once complete, each repo will have:
```bash
git remote -v
# origin  https://github.com/onemindos/<repo>.git (fetch)
# origin  https://github.com/onemindos/<repo>.git (push)
# origin  https://codeberg.org/onemindos/<repo>.git (push)
# origin  http://localhost:3000/onemindos/<repo>.git (push)
```

One `git push` hits all three simultaneously.

---

## Fleet Node Sync Strategy

Nodes do NOT push — they pull only.

Option A — cron pull (simple):
```bash
*/15 * * * * cd ~/dev/<repo> && git pull --ff-only
```

Option B — git-sync daemon (cleaner for robots):
```bash
# github.com/kubernetes/git-sync — single binary
# Watches repo, auto-pulls on new commits
# Nodes hit Forgejo over Headscale mesh (internal, no internet required)
```

---

## Gitify Setup (after Codeberg token is added)

1. Open /Applications/Gitify.app
2. Sign in with GitHub account
3. Add Codeberg account (Settings > Accounts > Add Gitea/Forgejo)
   - Host: https://codeberg.org
   - Token: (same token used for CLI setup)
4. Notifications from both platforms in one menu bar icon

---

## AGENTS.md Update Needed (PENDING)

File: ~/dev/onemind-command/AGENTS.md
Issue: Still references "OneMind-OS" org name in several places
Action: Legacy to run find/replace after Codeberg + Forgejo are live

---

## NEXT STEPS (resume here)

  [ ] 1. Zeus creates Codeberg account + onemindos org + API token
           → codeberg.org → sign up → create org "onemindos"
           → codeberg.org/user/settings/applications → new token "Legacy CLI"
           → scopes: repository (read+write) + organization (read+write)
  [ ] 2. Zeus pastes Codeberg API token to Legacy
  [ ] 3. Legacy creates 5 Codeberg repos + adds as second push remote on all repos
  [ ] 4. Legacy installs Forgejo via brew install forgejo
  [ ] 5. Legacy runs first-launch setup (Zeus completes web UI at localhost:3000)
           → choose SQLite, set admin username/password, org: onemindos
  [ ] 6. Legacy writes LaunchAgent plist + loads it (auto-start on login)
  [ ] 7. Legacy creates 5 Forgejo repos + adds as third push remote on all repos
  [ ] 8. Test push all three remotes on all repos
  [ ] 9. Configure Gitify — sign in GitHub + add Codeberg account
           → Settings > Accounts > Add Gitea/Forgejo > host: https://codeberg.org
  [ ] 10. Create onemind-dashboard GitHub repo (currently no remote, local dir missing)
  [ ] 11. Update AGENTS.md in onemind-command (still references OneMind-OS org name)
  [ ] 12. Document final remote state per repo in this file (update table below)

---

## Notes & Decisions

- No Docker — Forgejo runs as native binary + LaunchAgent
- Codeberg chosen over GitLab: EU-hosted, non-profit, no tracking, Forgejo-based
- Multi-remote push chosen over Forgejo push mirrors for simplicity
  (push mirrors = extra config per repo in Forgejo UI; multi-remote = one git config)
- All internal fleet traffic over Headscale — Forgejo never public-facing
- lazygit for daily terminal git operations
- Gitify for notification aggregation across all platforms
