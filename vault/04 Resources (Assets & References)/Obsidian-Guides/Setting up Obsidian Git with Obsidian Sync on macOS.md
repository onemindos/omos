---
id: ui-06-20260111053926
title: "Setting up Obsidian Git with Obsidian Sync on macOS"
type: reference
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/06]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Setting up Obsidian Git with Obsidian Sync on macOS"
---

# Setting up Obsidian Git with Obsidian Sync on macOS

Your setup strategy is sound: **Git handles version control on desktop** (for VS Code/Claude Code integration), while **Obsidian Sync handles cross-device synchronization** to iPhone. These systems can coexist without conflicts when configured properly.

## Understanding the architecture

Git and Obsidian Sync serve complementary roles. Obsidian Sync provides real-time synchronization across all your devices, including mobile. Git provides version control, backup, and integration with development tools on your desktop machines. The key is configuring Git to ignore device-specific files that Obsidian Sync manages, preventing conflicts between the two systems.

## Prerequisites

**Install Git on macOS:**
```bash
# Check if Git is already installed
git --version

# Install via Homebrew (recommended)
brew install git

# Or install Xcode Command Line Tools
xcode-select --install
```

**Configure Git globally (first time only):**
```bash
git config --global user.name "Your Name"
git config --global user.email "[email protected]"
git config --global credential.helper osxkeychain
```

The credential helper stores your GitHub token securely in macOS Keychain, so you won't need to enter it repeatedly.

## Creating the GitHub repository

**Via GitHub web interface:**
1. Log into GitHub at https://github.com
2. Click the "+" button (top right) → "New repository"
3. Repository name: Choose something like `obsidian-vault` or `notes-vault`
4. **Privacy: Select "Private"** (essential for personal notes)
5. **Do NOT check** "Add README", "Add .gitignore", or "Choose a license" - these create initial commits that will conflict with your existing vault
6. Click "Create repository"
7. Copy the HTTPS URL shown (looks like `https://github.com/username/vault-name.git`)

## Generating a GitHub Personal Access Token

GitHub deprecated password authentication, so you need a Personal Access Token:

1. GitHub → Settings (click your profile picture) → Developer settings (bottom of left sidebar)
2. Personal access tokens → Tokens (classic) → Generate new token (classic)
3. Note: "Obsidian Git Access"
4. Expiration: Choose "No expiration" or custom period
5. Scopes: Check **"repo"** (grants full repository access)
6. Click "Generate token" at bottom
7. **Copy the token immediately** - you won't see it again
8. Store it securely (password manager recommended)

## Setting up .gitignore first

Before initializing Git, create a proper `.gitignore` file to prevent conflicts between Git and Obsidian Sync. This is **critical** for avoiding issues.

```bash
# Navigate to your vault directory
cd /path/to/your/obsidian-vault

# Create .gitignore file
cat > .gitignore << 'EOF'
# === Obsidian workspace files (always exclude) ===
.obsidian/workspace.json
.obsidian/workspace-mobile.json
.obsidian/workspace
.obsidian/workspace.json.bak

# === Cache and temporary files ===
.obsidian/cache/

# === Obsidian Sync metadata ===
.obsidian/sync.json

# === Device-specific settings ===
.obsidian/app.json
.obsidian/appearance.json
.obsidian/hotkeys.json

# === Obsidian Git plugin data ===
.obsidian-git-data

# === macOS system files ===
.DS_Store
.AppleDouble
.LSOverride
._*

# === Obsidian trash ===
.trash/

# === Plugin code (keep config only) ===
.obsidian/plugins/*/main.js
.obsidian/plugins/*/styles.css
.obsidian/plugins/*/manifest.json

# === Large plugin data ===
.smart-connections/
.pro/copilot-index-*.json

# === Keep these important files ===
# These lines are comments - don't ignore:
# !.obsidian/community-plugins.json
# !.obsidian/core-plugins.json
# !.obsidian/plugins/*/data.json
# !.obsidian/graph.json
# !.obsidian/templates.json
EOF
```

**What this .gitignore does:**

- **Excludes workspace files** that change every time you open/close tabs - these cause constant conflicts
- **Excludes Obsidian Sync metadata** so Git doesn't interfere with Sync's operation
- **Excludes device-specific settings** (appearance, hotkeys) that should differ per machine
- **Excludes plugin code** but keeps plugin configuration (data.json files)
- **Excludes macOS system files** like .DS_Store
- **Keeps essential configuration** like plugin lists and settings that should be versioned

## Initializing the Git repository

```bash
# Ensure you're in your vault directory
cd /path/to/your/obsidian-vault

# Initialize Git repository
git init

# Rename default branch to main (GitHub standard)
git branch -M main

# Add the remote repository
git remote add origin https://github.com/YOUR-USERNAME/YOUR-REPO.git

# Verify remote is set correctly
git remote -v
```

## Initial commit and push for 1000+ files

For large vaults, this process takes 5-15 minutes. Ensure stable internet connection.

```bash
# Stage all files (respects .gitignore)
git add .

# Check what's being added
git status

# Create initial commit
git commit -m "Initial commit: Import Obsidian vault"

# Configure larger buffer for the large initial push
git config http.postBuffer 524288000  # 500MB
git config http.timeout 300  # 5 minutes

# Push to GitHub (first push sets upstream)
git push -u origin main
```

**When prompted for credentials:**
- Username: Your GitHub username
- Password: **Paste the Personal Access Token** you generated earlier

The macOS Keychain will store these credentials securely, so you won't need to enter them again.

## Installing the Obsidian Git plugin

**In Obsidian:**

1. Open Settings (gear icon, bottom left)
2. Navigate to **Community plugins**
3. If first time: Toggle OFF "Safe mode" → Click "Turn on community plugins"
4. Click **"Browse"** button
5. Search for **"Git"**
6. Find "Obsidian Git" by Vinzent03 → Click **"Install"**
7. After installation completes, click **"Enable"**
8. You should see "Git" listed under "Installed plugins"

## Configuring auto-commit every 30 minutes

Navigate to **Settings → Obsidian Git** (under Community plugins section).

**Essential settings for 30-minute auto-backup:**

| Setting | Value | Why |
|---------|-------|-----|
| **Vault backup interval (minutes)** | **30** | Main setting - commits AND pushes every 30 min |
| **Auto pull interval (minutes)** | **30** | Pulls changes from GitHub every 30 min |
| **Disable push** | **OFF** | Must be OFF for auto-push to work |
| **Push on backup** | **ON** | Ensures backup interval includes push operation |
| **Pull updates on startup** | **ON** | Syncs latest when opening Obsidian |
| **Pull changes before push** | **ON** | Prevents conflicts by pulling first |

**Commit message settings:**

| Setting | Recommended Value |
|---------|------------------|
| **Commit message on auto backup** | `vault backup: {{date}}` |
| **{{date}} placeholder format** | `YYYY-MM-DD HH:mm:ss` |
| **List filenames in commit body** | **ON** (optional but helpful) |

For multi-device tracking, you can use: `{{hostname}} backup: {{date}}` to identify which Mac made the commit.

**Branch and merge settings:**

| Setting | Value |
|---------|-------|
| **Current branch** | `main` |
| **Pull/Sync method** | Rebase (cleaner history) or Merge (safer) |

**Notification settings:**

- **Disable notifications**: OFF (so you see backup confirmations)
- **Show status bar**: ON (displays sync status at bottom)

## Verifying it's working correctly

**Test manual operations first:**

1. Make a small change to a note
2. Open Command Palette: **Cmd+P**
3. Type and run: **"Git: Commit and sync"**
4. You should see notification: "Pushing X files..." then "Pushed successfully"
5. Visit your GitHub repository page - verify the commit appears
6. Check the commit message and files changed

**Monitor automatic backup:**

1. Make another small change
2. Wait 30 minutes
3. Watch for automatic "Pushing files..." notification
4. Verify commit appears on GitHub with timestamp

**Check status bar:**

Look at bottom-right of Obsidian window - should show Git icon and current branch (main). This updates during operations.

**Command line verification:**

```bash
cd /path/to/vault

# Check repository status
git status

# View recent commits
git log --oneline -10

# Verify remote connection
git remote -v

# Check last push succeeded
git log origin/main -1
```

## How Git and Obsidian Sync work together

**The division of labor:**

**Git handles (desktop only):**
- Version control and commit history
- Backup to GitHub
- Integration with VS Code and Claude Code
- Manual commits for important milestones
- Automatic commits every 30 minutes on desktop

**Obsidian Sync handles (all devices):**
- Real-time synchronization of markdown files
- Mobile device synchronization (iPhone)
- Immediate sync of changes across devices
- Workspace and settings sync (if enabled in Sync settings)
- Conflict resolution for simultaneous edits

**Why they don't conflict:**

1. **.gitignore excludes** Obsidian Sync's metadata files (like `.obsidian/sync.json`)
2. **.gitignore excludes** workspace files that change constantly
3. **Git operates on desktop only** - don't install Obsidian Git plugin on iPhone
4. **Obsidian Sync handles mobile** - seamless, automatic
5. **Both systems operate independently** on your actual notes (markdown files)
6. **Time separation** - Obsidian Sync works in real-time; Git commits periodically

**Recommended workflow:**

- **On desktop Macs**: Both Git (for version control) and Obsidian Sync (for device sync) run simultaneously
- **On iPhone**: Only Obsidian Sync runs - no Git plugin needed or recommended
- **Content syncing**: Obsidian Sync keeps all devices current
- **Version history**: Git provides deep history and rollback capability on desktop
- **Development integration**: Git enables VS Code and Claude Code integration

## Using separate Git directory (advanced, optional)

For maximum separation between Git and Obsidian Sync, you can store Git metadata outside the vault:

```bash
cd /path/to/vault

# Move .git directory outside vault
mv .git ~/Documents/vault-git-metadata

# Create pointer to external Git directory
echo "gitdir: $HOME/Documents/vault-git-metadata" > .git

# Verify it works
git status
```

This prevents Obsidian Sync from ever touching Git files. However, the `.gitignore` approach above is usually sufficient and simpler.

## Common troubleshooting issues

**Plugin shows "Git is not ready"**

Check Git is properly installed and repository initialized:
```bash
which git  # Should show path like /usr/local/bin/git
cd /path/to/vault
ls -la .git  # Should see .git directory
git status  # Should show branch info
```

**Authentication failures when pushing**

Your token may have expired or wasn't stored correctly:
```bash
# Remove stored credentials
git credential-osxkeychain erase
host=github.com
protocol=https

# Next push will prompt for credentials again
# Enter username and paste token as password
```

**"Failed to push" errors**

Usually means remote has changes you don't have locally:
```bash
# Pull first, then push
git pull origin main --rebase
git push origin main
```

Or use plugin: Command Palette → "Git: Pull" then "Git: Push"

**Constant workspace.json conflicts**

If you didn't exclude workspace files initially:
```bash
# Remove from tracking
git rm --cached .obsidian/workspace.json
git rm --cached .obsidian/workspace-mobile.json
git commit -m "Stop tracking workspace files"
git push
```

**Plugin stuck at "Checking repo status"**

Usually a remote/branch name issue:
```bash
# Verify remote is named "origin"
git remote -v

# Rename if needed
git remote rename oldname origin

# Verify branch is "main"
git branch
```

**index.lock error after crash**

```bash
cd /path/to/vault
rm -f .git/index.lock
```

Then restart Obsidian.

**Large vault performance issues**

For very large vaults (10,000+ files):
- Increase backup interval to 60 minutes
- Or disable auto-backup (set to 0) and commit manually
- Use Command Palette → "Git: Commit and sync" when desired
- Consider archiving old notes to separate repository

**SSH authentication hanging**

If you set up SSH instead of HTTPS and it hangs:
```bash
# Test SSH connection
ssh -T [email protected]

# If prompts for passphrase repeatedly, add to keychain:
ssh-add --apple-use-keychain ~/.ssh/id_ed25519

# Configure SSH to use keychain
cat >> ~/.ssh/config << 'EOF'
Host github.com
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_ed25519
EOF
```

However, HTTPS with tokens is generally more reliable for Obsidian Git on macOS.

## Best practices for commit messages

The plugin supports dynamic placeholders:

**Variables available:**
- `{{date}}` - Current timestamp (format customizable)
- `{{hostname}}` - Computer name (useful for multi-device setups)
- `{{numFiles}}` - Number of changed files
- `{{files}}` - List of changed files (in commit body)

**Recommended templates:**

For single device:
```
vault backup: {{date}}
```

For multiple desktop devices:
```
{{hostname}}: {{date}}
```

For detailed tracking:
```
auto: {{numFiles}} changed - {{date}}
```

Enable "List filenames affected by commit in commit body" for transparency about what changed.

## VS Code and Claude Code integration

With Git set up, you can now integrate with development tools:

**VS Code:**
1. Open your vault folder in VS Code: `code /path/to/vault`
2. Git integration works automatically
3. Use Source Control panel (Cmd+Shift+G) to view changes
4. Make commits directly from VS Code if preferred
5. VS Code respects your .gitignore

**Claude Code / other AI tools:**
- They can now access your vault's Git history
- View version history and changes over time
- Analyze patterns in your note-taking
- Reference specific commits or changes

**Important**: When committing from external tools (VS Code, command line), the changes still sync via Obsidian Sync to your iPhone. The systems work independently.

## Configuration values quick reference

**File locations:**
- Vault: `~/Documents/YourVaultName/` (or wherever you chose)
- Git metadata: `~/Documents/YourVaultName/.git/`
- Plugin settings: `~/Documents/YourVaultName/.obsidian/plugins/obsidian-git/data.json`

**Key settings recap:**
```
Vault backup interval: 30 minutes
Auto pull interval: 30 minutes
Disable push: OFF
Pull updates on startup: ON
Pull before push: ON
Commit message: vault backup: {{date}}
Date format: YYYY-MM-DD HH:mm:ss
```

**Essential commands:**
```bash
git status              # Check current state
git log --oneline -10   # View recent commits
git pull origin main    # Pull latest changes
git push origin main    # Push your changes
git remote -v           # Check remote URL
```

**Obsidian commands (Cmd+P):**
- `Git: Commit and sync` - Complete backup operation
- `Git: Open source control view` - Visual file staging
- `Git: Open history view` - Browse commits
- `Git: Pull` - Get latest from GitHub
- `Git: Push` - Send commits to GitHub

## Final setup checklist

Before considering setup complete:

- [ ] Git installed and configured on macOS
- [ ] GitHub repository created (private)
- [ ] Personal Access Token generated and stored securely
- [ ] .gitignore created with proper exclusions
- [ ] Repository initialized and initial commit pushed
- [ ] Obsidian Git plugin installed and enabled
- [ ] Auto-commit configured for 30-minute intervals
- [ ] Manual commit tested successfully
- [ ] Commit appears on GitHub.com
- [ ] Auto-commit tested (wait 30 minutes after a change)
- [ ] Obsidian Sync enabled and working on all devices
- [ ] iPhone syncing properly (via Obsidian Sync only)
- [ ] VS Code can open vault and see Git history
- [ ] Status bar shows Git status in Obsidian

Your setup is now complete. You have version control with Git on desktop (enabling VS Code and Claude Code integration), automatic commits every 30 minutes, backup to GitHub, and seamless synchronization to iPhone via Obsidian Sync. The two systems operate independently without conflicts, giving you the best of both worlds.