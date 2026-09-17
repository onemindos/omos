# Fresh Install Guide — MacBook Pro → Local OpenClaw Host

> **Created:** 2026-05-15
> **Purpose:** Complete instructions to factory reset the MacBook Pro and rebuild it as the always-on OpenClaw host, replacing the Lightsail VPS (omos-core).
> **Decision rationale:** M2 Max (12-core, 32GB) outperforms the 8-vCPU VPS in every dimension. Eliminates ~$200/mo cloud cost. Latency drops to zero. Full browser visibility for automation.

---

## Hardware

| Field | Value |
|-------|-------|
| Model | MacBook Pro Mac14,6 |
| Chip | Apple M2 Max (12-core: 8P + 4E) |
| RAM | 32 GB unified |
| Storage | 1 TB NVMe (~434 GB free pre-reset) |
| macOS | 26.x (latest at time of reset) |
| Battery | 82% / 233 cycles |
| Setup | Docked, clamshell mode, two external displays, always on power |

---

## Phase 0 — Before You Reset

### Back Up These (they don't live in the vault)

- [ ] **1Password** — Confirm vault syncs to cloud (all secrets are here, nothing local-only)
- [ ] **SSH keys** — `~/.ssh/` — back up to 1Password or external drive
  - Key files: `config`, `codex_deploy`, any other private keys
- [ ] **Git config** — `~/.gitconfig` (name, email, signing key, aliases)
- [ ] **Tailscale** — Note your device name (`zeus-macbook-pro`) and IP (`100.114.142.55`). You'll re-auth after reset
- [ ] **Time Machine** — Ensure latest backup completes. You won't restore FROM it (fresh install), but it's your safety net
- [ ] **Browser bookmarks** — Export from Brave/Chrome/Arc (or confirm synced to account)
- [ ] **VS Code settings** — Confirm Settings Sync is ON (GitHub account). Extensions, keybindings, snippets all sync
- [ ] **Claude Code / Copilot** — Session history is cloud-synced, nothing to save
- [ ] **OpenClaw vault** — The onemind vault is in git. Confirm `git push` is current:
  ```bash
  cd ~/Documents/onemind && git status && git push
  ```
- [ ] **Obsidian settings** — `.obsidian/` is in the vault repo. Confirm committed
- [ ] **Fonts** — Note any custom-installed fonts (JetBrains Mono, etc.)
- [ ] **Dock / app list** — Screenshot your dock for reference

### What You Do NOT Need to Back Up

- OpenClaw agent sessions (207 MB on VPS — start fresh, that's the point)
- Brave browser data (29.5 GB of RAM hog — good riddance)
- Any VPS-specific paths or configs
- Node modules, caches, build artifacts

### Reference: VPS Config to Extract (Not Copy)

These are the VALUES to reference when building fresh config. Do NOT copy the VPS config files directly (they contain Linux paths, dead Discord config, plaintext secrets, and cruft).

**From VPS `~/.openclaw/openclaw.json`:**
- Telegram bot token → 1Password: `op://aaeyris5z44pdqspk277u2oppi/c2hsruwzaum6z4d7qjjyduadpq/credential`
- Telegram group ID: `-1003905828404` (forum mode)
- Telegram authorized user (Zeus): `577573558`
- Telegram thread IDs:
  ```
  General: 1          ⚡ Legacy: 292      Zenith: 318
  Unity: 319          Eternal: 320        Sovereignty: 321
  Sentinel: 256       Approvals: 322      📥 Inbox: 323
  📋 Directives: 324  📊 Daily: 325       📆 Weekly: 326
  📣 Briefings: 327   🖥️ Status: 328      🚨 Errors: 329
  ⏱️ Cron: 330        🔄 Active: 331      Completed: 332
  ```
- Model config: Primary `github-copilot/claude-sonnet-4.6`, fallback `amazon-bedrock/us.anthropic.claude-sonnet-4-6`
- MCP servers to reconfigure: `github`, `todoist`, `filesystem` (codex-vault), `context7`, `home-assistant`
- Agent IDs: `main` (Legacy), `zenith`, `unity`, `eternal`, `sovereignty`, `sentinel`

**From VPS `gateway.systemd.env` (secrets — pull from 1Password, not the file):**
- AWS credentials (for Bedrock)
- GitHub PAT
- Todoist API key
- Home Assistant JWT
- Any other API keys

---

## Phase 1 — Factory Reset macOS

1. Sign out of iCloud, App Store, Find My Mac
2. Restart → Hold Power button → Options → Recovery Mode
3. Disk Utility → Erase "Macintosh HD" (APFS)
4. Reinstall macOS (latest)
5. Create user account: `zeuslegacy`
6. Skip migration assistant — this is a FRESH start

---

## Phase 2 — Core System Setup

Run these in order after first boot.

### 2.1 — Xcode Command Line Tools
```bash
xcode-select --install
```

### 2.2 — Homebrew
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/opt/homebrew/bin/brew shellenv)"
```

### 2.3 — Essential CLI Tools
```bash
brew install git gh jq ripgrep fd tree wget curl
```

### 2.4 — Node.js via nvm
```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash
source ~/.zshrc
nvm install --lts
nvm alias default lts/*
node -v && npm -v
```

### 2.5 — Tailscale
```bash
brew install --cask tailscale
```
Open Tailscale → Sign in → Confirm device name is `zeus-macbook-pro` and IP is `100.114.142.55` (or note new IP).

### 2.6 — 1Password + CLI
```bash
brew install --cask 1password
brew install --cask 1password-cli
```
Sign in to 1Password. Enable CLI integration in 1Password settings.

### 2.7 — Git Config
```bash
git config --global user.name "Zeus Delacruz"
git config --global user.email "YOUR_EMAIL"
git config --global init.defaultBranch main
git config --global pull.rebase true
```

### 2.8 — SSH Keys
Restore from 1Password backup, or generate new:
```bash
mkdir -p ~/.ssh && chmod 700 ~/.ssh
# Restore keys from 1Password, then:
chmod 600 ~/.ssh/id_* ~/.ssh/codex_deploy
chmod 644 ~/.ssh/*.pub
# Restore ~/.ssh/config
```

### 2.9 — Clone the Vault
```bash
mkdir -p ~/Documents
cd ~/Documents
git clone git@github.com:OneMind-OS/One-Mind.git onemind
```

---

## Phase 3 — Applications

### 3.1 — Productivity
```bash
brew install --cask obsidian
brew install --cask visual-studio-code
brew install --cask arc                # Personal browser
brew install --cask google-chrome      # OpenClaw automation browser (dedicated)
brew install --cask iterm2             # Optional — better terminal
```

### 3.2 — VS Code Extensions
Settings Sync should restore everything automatically on sign-in. If not, key extensions:
- GitHub Copilot + Copilot Chat
- Claude Code (if extension exists)
- Markdown All in One
- YAML
- GitLens

### 3.3 — Obsidian
Open `~/Documents/onemind` as vault. Settings are in `.obsidian/` (committed to git).

---

## Phase 4 — OpenClaw Installation

### 4.1 — Install OpenClaw
```bash
npm install -g openclaw@latest
openclaw --version
# Should be >= 2026.5.12
```

### 4.2 — Initialize Gateway
```bash
openclaw init
```

### 4.3 — Configure Models
```bash
# Primary model — GitHub Copilot (free with Copilot subscription)
openclaw config set models.primary "github-copilot/claude-sonnet-4.6"

# Fallback — Amazon Bedrock
openclaw config set models.fallback "amazon-bedrock/us.anthropic.claude-sonnet-4-6"
```

Set API keys via environment variables (NOT in openclaw.json):
```bash
# Create ~/.openclaw/gateway.env (or equivalent)
# Pull values from 1Password — NEVER hardcode in config files

# AWS for Bedrock
export AWS_ACCESS_KEY_ID="op://aaeyris5z44pdqspk277u2oppi/AWS-Bedrock/access-key-id"
export AWS_SECRET_ACCESS_KEY="op://aaeyris5z44pdqspk277u2oppi/AWS-Bedrock/secret-access-key"
export AWS_DEFAULT_REGION="us-east-1"

# GitHub Copilot — handled by gh auth
```

### 4.4 — Configure Agents

6 agents. Create each:

```bash
# Legacy — Orchestrator
openclaw agent create main --name "Legacy" --model "github-copilot/claude-sonnet-4.6"

# Zenith — HP (Personal)
openclaw agent create zenith --name "Zenith" --model "github-copilot/claude-sonnet-4.6"

# Unity — LE (Home)  
openclaw agent create unity --name "Unity" --model "github-copilot/claude-sonnet-4.6"

# Eternal — GE (Business)
openclaw agent create eternal --name "Eternal" --model "github-copilot/claude-sonnet-4.6"

# Sovereignty — SO (Operations)
openclaw agent create sovereignty --name "Sovereignty" --model "github-copilot/claude-sonnet-4.6"

# Sentinel — Watchdog (runs on omos-sentinel, but register locally too)
openclaw agent create sentinel --name "Sentinel" --model "github-copilot/claude-sonnet-4.6"
```

Copy workspace identity files from the vault into each agent workspace:
```bash
VAULT="$HOME/Documents/onemind"
AGENTS_DIR="$VAULT/00-24 SO (Sovereign Operations)/01 Intelligence (The Brain)/Agents (Digital arms)/_COMMAND"

# For each agent, copy their identity files to ~/.openclaw/workspace-{id}/
# Structure: IDENTITY.md, SOUL.md, TOOLS.md, HEARTBEAT.md, AGENTS.md, USER.md, MEMORY.md
```

### 4.5 — Configure Telegram
```bash
# Set bot token from 1Password
TELEGRAM_TOKEN=$(op read "op://aaeyris5z44pdqspk277u2oppi/c2hsruwzaum6z4d7qjjyduadpq/credential")
openclaw config set channels.telegram.accounts.default.token "$TELEGRAM_TOKEN"

# Set authorized user
openclaw config set channels.telegram.dmPolicy allowlist
openclaw config set channels.telegram.accounts.default.allowlist "[577573558]"
```

Forum group and thread routing will need to be configured in `openclaw.json` — reference the thread IDs from Phase 0.

### 4.6 — Configure MCP Servers

```json
{
  "mcp": {
    "servers": {
      "codex-vault": {
        "command": "npx",
        "args": ["-y", "@modelcontextprotocol/server-filesystem", "/Users/zeuslegacy/Documents/onemind"]
      },
      "github": {
        "command": "npx",
        "args": ["-y", "@modelcontextprotocol/server-github"],
        "env": {
          "GITHUB_PERSONAL_ACCESS_TOKEN": "${GITHUB_PAT}"
        }
      },
      "context7": {
        "command": "npx",
        "args": ["-y", "@upstash/context7-mcp@latest"]
      },
      "todoist": {
        "command": "npx",
        "args": ["-y", "@modelcontextprotocol/server-todoist"],
        "env": {
          "TODOIST_API_KEY": "${TODOIST_API_KEY}"
        }
      },
      "home-assistant": {
        "command": "npx",
        "args": ["-y", "@modelcontextprotocol/server-home-assistant"],
        "env": {
          "HA_TOKEN": "${HA_JWT_TOKEN}"
        }
      }
    }
  }
}
```

> **CRITICAL:** Use environment variable references or `op://` references for ALL secrets. Never put plaintext API keys in `openclaw.json`.

### 4.7 — Configure Browser

```json
{
  "browser": {
    "enabled": true,
    "defaultProfile": "openclaw",
    "headless": false,
    "executablePath": "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
    "profiles": {
      "openclaw": {
        "cdpPort": 18800,
        "color": "#FF4500"
      },
      "user": {
        "driver": "existing-session",
        "attachOnly": true,
        "color": "#00AA00"
      }
    }
  }
}
```

- `openclaw` profile: Dedicated Chrome instance, visible on second monitor
- `user` profile: Attach to your real Chrome session for logged-in automation
- `headless: false`: You SEE every action the agent takes

### 4.8 — Gateway Auth

```bash
# Set gateway token — use the ENVIRONMENT VARIABLE method, not the config file
# This prevents OpenClaw from rotating it on config writes
export OPENCLAW_GATEWAY_TOKEN="onemind-zeus-gateway-$(openssl rand -hex 16)"
# Save this token to 1Password immediately
```

> **WARNING:** Do NOT put `gateway.auth.token` in `openclaw.json`. The env var `OPENCLAW_GATEWAY_TOKEN` is the sole source of truth. If the field exists in the config file, OpenClaw rotates it on every config write, causing token mismatch loops.

---

## Phase 5 — Always-On Configuration (Clamshell Mode)

### 5.1 — Prevent Sleep
```bash
# Disable all sleep/standby while on power
sudo pmset -a sleep 0
sudo pmset -a standby 0
sudo pmset -a autopoweroff 0
sudo pmset -a hibernatemode 0
sudo pmset -a disablesleep 1

# Verify
pmset -g
```

### 5.2 — Clamshell Mode Requirements
- MacBook connected to power adapter (mandatory)
- At least one external display connected
- External keyboard + mouse (Bluetooth or USB)
- Close the lid → Mac continues running on external displays

### 5.3 — Auto-Login (Optional)
System Settings → Users & Groups → Login Options → Automatic Login → `zeuslegacy`

### 5.4 — Launch OpenClaw on Boot

Create a LaunchAgent:
```bash
cat > ~/Library/LaunchAgents/com.openclaw.gateway.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.openclaw.gateway</string>
    <key>ProgramArguments</key>
    <array>
        <string>/Users/zeuslegacy/.nvm/versions/node/vXX.X.X/bin/node</string>
        <string>/Users/zeuslegacy/.nvm/versions/node/vXX.X.X/bin/openclaw</string>
        <string>gateway</string>
        <string>start</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/Users/zeuslegacy/.openclaw/logs/gateway.stdout.log</string>
    <key>StandardErrorPath</key>
    <string>/Users/zeuslegacy/.openclaw/logs/gateway.stderr.log</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin</string>
    </dict>
</dict>
</plist>
EOF

mkdir -p ~/.openclaw/logs
# Update the node path above to match your actual nvm node version
# Then load it:
launchctl load ~/Library/LaunchAgents/com.openclaw.gateway.plist
```

### 5.5 — Vault Auto-Sync (Replaces VPS inotifywait)

Create a LaunchAgent for git auto-commit + push:
```bash
cat > ~/Library/LaunchAgents/com.onemind.vault-sync.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.onemind.vault-sync</string>
    <key>ProgramArguments</key>
    <array>
        <string>/bin/bash</string>
        <string>/Users/zeuslegacy/.openclaw/scripts/vault-sync.sh</string>
    </array>
    <key>StartInterval</key>
    <integer>30</integer>
    <key>RunAtLoad</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/Users/zeuslegacy/.openclaw/logs/vault-sync.stdout.log</string>
    <key>StandardErrorPath</key>
    <string>/Users/zeuslegacy/.openclaw/logs/vault-sync.stderr.log</string>
</dict>
</plist>
EOF
```

Create the sync script:
```bash
mkdir -p ~/.openclaw/scripts
cat > ~/.openclaw/scripts/vault-sync.sh << 'SCRIPT'
#!/bin/bash
cd /Users/zeuslegacy/Documents/onemind || exit 1

# Pull first
git pull --rebase --quiet 2>/dev/null

# Check for changes
if [[ -n $(git status --porcelain) ]]; then
    git add -A
    git commit -m "auto: vault sync $(date +%Y-%m-%d_%H:%M:%S)" --quiet
    git push --quiet
fi
SCRIPT
chmod +x ~/.openclaw/scripts/vault-sync.sh
```

---

## Phase 6 — Post-Install Verification

Run through this checklist after everything is set up:

- [ ] `openclaw --version` shows >= 2026.5.12
- [ ] `openclaw gateway status` shows running
- [ ] `openclaw agent list` shows all 6 agents
- [ ] Send a test message in Telegram → Legacy responds
- [ ] `tailscale status` shows connected, other nodes visible
- [ ] `git -C ~/Documents/onemind status` shows clean working tree
- [ ] Open Obsidian → vault loads, dataview queries work
- [ ] Open VS Code → Settings Sync restored, Copilot active
- [ ] `pmset -g` confirms sleep=0, standby=0, autopoweroff=0
- [ ] Close MacBook lid → external displays stay on, SSH still works
- [ ] Browser automation: `openclaw browser test` or equivalent
- [ ] MCP servers: test `codex-vault` read, `github` API call, `todoist` list
- [ ] Check RAM: `top -l 1 | grep PhysMem` — should show 18-21 GB free (no Brave bloat)

---

## Phase 7 — VPS Decommission (After 1 Week Burn-In)

Only after the Mac has run stable for 7 days:

1. **Stop VPS agents**: `ssh omos-core "openclaw agent stop --all"`
2. **Stop VPS services**: `ssh omos-core "systemctl --user stop onemind-watch onemind-pull openclaw-gateway"`
3. **Final backup**: `ssh omos-core "tar czf ~/omos-core-final-backup.tar.gz ~/.openclaw/"`
4. **Download backup**: `scp omos-core:~/omos-core-final-backup.tar.gz ~/Desktop/`
5. **Update Tailscale**: Remove `omos-core` from Tailscale admin if decommissioning entirely
6. **Terminate Lightsail instance** — saves ~$200/mo
7. **Update DNS/references**: Any docs pointing to `44.193.145.121` or `100.123.126.34`
8. **Keep `omos-sentinel`** as external watchdog (separate node, low cost)

---

## Architecture: Before vs After

```
BEFORE (Cloud):
  Mac (Obsidian/VS Code) → git push → GitHub → git pull → VPS (OpenClaw)
  Telegram → VPS → agents → VPS → Telegram
  Browser automation: headless on VPS (can't see anything)

AFTER (Local):
  Mac = EVERYTHING
  ├── Obsidian (vault UI)
  ├── VS Code (coding)
  ├── OpenClaw gateway + 6 agents (local process)
  ├── Chrome (visible browser automation on second monitor)
  ├── Telegram polling (Mac → Telegram servers, no public IP needed)
  └── Git sync (auto-commit every 30s → GitHub → omos-sentinel pulls)
  
  omos-sentinel (VPS) = watchdog only (low cost, monitors Mac health)
```

---

## RAM Budget (Post-Reset, Disciplined Use)

| Process | Expected RAM |
|---------|-------------|
| macOS + WindowServer | ~4 GB |
| OpenClaw gateway + 6 agents | ~1.5 GB |
| Obsidian | ~0.5 GB |
| VS Code + extensions | ~1.5 GB |
| Chrome (OpenClaw, ~5 tabs) | ~1 GB |
| Arc (personal, ~10 tabs) | ~2 GB |
| Misc (Tailscale, 1Password, etc.) | ~0.5 GB |
| **Total estimated** | **~11 GB** |
| **Free for burst / file ops** | **~21 GB** |

> Previous state: Brave alone was eating 29.5 GB. The fresh install with Chrome + Arc discipline should leave 20+ GB free at all times.

---

## Quick Reference — Key Paths

| What | Path |
|------|------|
| Vault | `~/Documents/onemind/` |
| OpenClaw config | `~/.openclaw/openclaw.json` |
| OpenClaw logs | `~/.openclaw/logs/` |
| Agent workspaces | `~/.openclaw/workspace-{id}/` |
| Gateway env | `~/.openclaw/gateway.env` |
| Vault sync script | `~/.openclaw/scripts/vault-sync.sh` |
| LaunchAgents | `~/Library/LaunchAgents/com.openclaw.*.plist` |
| SSH keys | `~/.ssh/` |
| Node (nvm) | `~/.nvm/versions/node/` |
| Homebrew | `/opt/homebrew/` |

---

## Tailscale Mesh (Expected Post-Reset)

| Node | IP | Role |
|------|-----|------|
| `zeus-macbook-pro` | `100.114.142.55` | Primary host (OpenClaw + vault + everything) |
| `omos-sentinel` | `100.66.235.38` | External watchdog |
| `zeus-iphone-17` | `100.122.203.68` | Mobile access |
| `omos-core` | `100.123.126.34` | DECOMMISSIONED after burn-in |
| `omos-dell` | `100.66.84.71` | Edge node |
| `omos-ha` | `100.64.25.12` | Home Assistant |
| `omos-hermes` | `100.111.255.72` | Hermes node |

---

## Appendix A — Current App & Service Audit (Pre-Reset)

> Inventoried 2026-05-15. Mark each item ✅ KEEP or ❌ SKIP before resetting.

### A.1 — GUI Applications (`/Applications/`)

| App | Category | Verdict | Notes |
|-----|----------|---------|-------|
| **1Password** | Security | ✅ KEEP | Secrets vault — install via `brew install --cask 1password` |
| **balenaEtcher** | Utility | ❌ SKIP | USB flasher — install only if needed later |
| **BambuStudio** | Hardware | ❓ DECIDE | 3D printer slicer — only if actively printing |
| **Blender** | Creative | ❓ DECIDE | 3D modeling — large app, install only if needed |
| **Blueprint IDE** | Dev | ❓ DECIDE | What is this used for? |
| **Brave Browser** | Browser | ❌ SKIP | RAM monster (29.5 GB). Replace with Arc (personal) + Chrome (OpenClaw) |
| **ChatGPT** | AI | ❓ DECIDE | Desktop app — or just use web? |
| **Circle** | Community | ✅ KEEP | Community management — `brew install --cask circle` or direct download |
| **Claude** | AI | ✅ KEEP | Desktop app for Claude |
| **cmux** | Utility | ❓ DECIDE | What is this? |
| **Discord** | Comms | ❌ SKIP | Deactivated 2026-05-12, Telegram is sole messaging layer |
| **Docker** | Dev | ✅ KEEP | Container runtime — `brew install --cask docker` |
| **Dropbox** | Storage | ❓ DECIDE | Backup layer? Or switch fully to Proton Drive? |
| **IDrive** | Backup | ❓ DECIDE | Cloud backup — still needed with Time Machine + git? |
| **Letta Code** | AI | ❓ DECIDE | Letta AI agent IDE — still using? |
| **OBS** | Creative | ❓ DECIDE | Screen recording / streaming — only if actively using |
| **Obsidian** | Core | ✅ KEEP | Vault UI — `brew install --cask obsidian` |
| **Ollama** | AI | ❌ SKIP | Local LLMs — using cloud models (Copilot/Bedrock), not needed on fresh install |
| **OpenClaw** | Core | ✅ KEEP | Agent gateway — `npm install -g openclaw@latest` |
| **Perplexity** | AI | ❓ DECIDE | AI search — or just use web? |
| **Proton Drive** | Storage | ✅ KEEP | Encrypted cloud storage |
| **Proton Mail** | Email | ✅ KEEP | Encrypted email client |
| **ProtonVPN** | Security | ✅ KEEP | VPN — privacy layer |
| **Raindrop.io** | Productivity | ❓ DECIDE | Bookmark manager — still using? |
| **Raycast** | Productivity | ✅ KEEP | Launcher — replaces Spotlight. Login item |
| **Rectangle** | Utility | ✅ KEEP | Window management — `brew install --cask rectangle` |
| **Safari** | Browser | ✅ KEEP | Ships with macOS — no action needed |
| **Stats** | Utility | ✅ KEEP | Menu bar system monitor — `brew install --cask stats` |
| **Tailscale** | Network | ✅ KEEP | Mesh VPN — `brew install --cask tailscale` |
| **Telegram** | Comms | ✅ KEEP | Primary messaging layer for OpenClaw |
| **Todoist** | Productivity | ✅ KEEP | Task execution layer |
| **Trident** | Unknown | ❓ DECIDE | What is this? |
| **Vaydeer keyboard** | Hardware | ❓ DECIDE | Keyboard software — still using this keyboard? |
| **Visual Studio Code** | Core | ✅ KEEP | Primary editor — `brew install --cask visual-studio-code` |
| **WhatsApp** | Comms | ❓ DECIDE | Personal messaging — or phone only? |
| **Xcode** | Dev | ✅ KEEP | Required for iOS dev (onemind-app) + command line tools |
| **Yubico Authenticator** | Security | ✅ KEEP | YubiKey manager — hardware 2FA |
| **Zoho WorkDrive TrueSync** | Storage | ❓ DECIDE | Zoho sync — still needed? |

**Summary: 15 ✅ KEEP, 4 ❌ SKIP, 17 ❓ DECIDE**

### A.2 — Homebrew CLI Tools (Formulae)

| Tool | Verdict | Notes |
|------|---------|-------|
| **awscli** | ✅ KEEP | AWS/Bedrock access |
| **azure-cli** | ❓ DECIDE | Still using Azure? |
| **bfg** | ❌ SKIP | Git history cleaner — one-time use, install if needed |
| **cloudflared** | ❓ DECIDE | Cloudflare tunnel — still using? |
| **cmake** | ❌ SKIP | Build tool — install only if compiling C/C++ |
| **cocoapods** | ✅ KEEP | iOS dependency manager (for onemind-app) |
| **docker / docker-completion** | ✅ KEEP | Container runtime CLI |
| **doctl** | ❓ DECIDE | DigitalOcean CLI — still on DO? |
| **exiftool** | ❌ SKIP | Photo metadata — install if needed |
| **fdupes** | ❌ SKIP | Duplicate finder — one-time use |
| **ffmpeg** | ✅ KEEP | Media processing — used by agents (hyperframes) |
| **flutter** | ✅ KEEP | onemind-app framework |
| **gcloud-cli** | ❓ DECIDE | Google Cloud — still using? |
| **gh** | ✅ KEEP | GitHub CLI |
| **git + git-lfs** | ✅ KEEP | Core |
| **git-crypt** | ❓ DECIDE | Encrypted git files — still using? |
| **git-filter-repo** | ❌ SKIP | One-time repo cleanup — install if needed |
| **gnupg** | ✅ KEEP | GPG encryption |
| **imagemagick** | ❓ DECIDE | Image processing — agents use it? |
| **jq** | ✅ KEEP | JSON processing — essential |
| **livekit / livekit-cli** | ❓ DECIDE | Real-time audio/video — still using? |
| **mole** | ❓ DECIDE | SSH tunnel tool — still using? |
| **nats** | ❓ DECIDE | NATS messaging — still in architecture? |
| **ngrok** | ❓ DECIDE | Tunneling — still using? |
| **nmap** | ❓ DECIDE | Network scanner — security audits? |
| **node** | ✅ KEEP | JavaScript runtime (use nvm instead of brew for fresh install) |
| **python@3.13 / python@3.14** | ✅ KEEP | Python runtime |
| **rclone** | ❓ DECIDE | Cloud storage sync — still using? |
| **ripgrep (rg)** | ✅ KEEP | Fast search — used everywhere |
| **ruby** | ✅ KEEP | Required by CocoaPods |
| **sqlite** | ✅ KEEP | Database — OpenClaw SQLite migration coming |
| **sshpass** | ❓ DECIDE | SSH password auth — should be using keys instead |
| **swiftformat / swiftlint** | ✅ KEEP | Swift formatting (onemind-app) |
| **terraform** | ❓ DECIDE | Infrastructure as code — still using? |
| **unar** | ❌ SKIP | Archive extractor — install if needed |
| **uv** | ✅ KEEP | Fast Python package manager |
| **xcodegen** | ✅ KEEP | Xcode project generation |
| **zig** | ❓ DECIDE | Zig compiler — what for? |

> **Note:** Many formulae (aom, cairo, harfbuzz, libheif, etc.) are **dependencies** of the above — they'll install automatically. No need to track them individually.

### A.3 — npm Global Packages

| Package | Verdict | Notes |
|---------|---------|-------|
| **openclaw** (2026.4.21) | ✅ KEEP | Upgrade to latest on fresh install |
| **openclaw-mcp** | ✅ KEEP | MCP bridge for OpenClaw |
| **@openai/codex** | ❓ DECIDE | OpenAI Codex CLI — still using? |
| **@letta-ai/letta-code** | ❓ DECIDE | Letta agent — still using? |
| **@chrusic/todoist-mcp-server-extended** | ✅ KEEP | Todoist MCP server |
| **pnpm** | ✅ KEEP | Package manager |
| **twilio-cli** | ❓ DECIDE | Twilio — still using? |
| **agent-browser** | ❓ DECIDE | Browser agent — replaced by OpenClaw browser? |

### A.4 — Background Services (LaunchAgents / LaunchDaemons)

| Service | Type | Verdict | Notes |
|---------|------|---------|-------|
| `ai.openclaw.gateway` | User Agent | ✅ KEEP | OpenClaw gateway — rebuild fresh |
| `ai.openclaw.node` | User Agent | ✅ KEEP | OpenClaw node — rebuild fresh |
| `com.onemind.env` | User Agent | ✅ KEEP | OneMind environment — rebuild fresh |
| `com.onemind.legacy-tunnel` | User Agent | ❓ DECIDE | Tunnel to VPS? Won't need after migration |
| `com.dropbox.dropboxmacupdate.*` | User Agent | ❓ DECIDE | Dropbox updater — depends on Dropbox decision |
| `com.google.keystone.*` | System Agent | ❌ SKIP | Chrome auto-updater — reinstalls with Chrome |
| `com.iDrive.IDFinderPluginApp` | System Agent | ❓ DECIDE | IDrive Finder plugin — depends on IDrive decision |
| `com.docker.socket` | System Daemon | ✅ KEEP | Docker socket — reinstalls with Docker |
| `com.docker.vmnetd` | System Daemon | ✅ KEEP | Docker networking — reinstalls with Docker |
| `com.prosoftnet.idrivedaemon` | System Daemon | ❓ DECIDE | IDrive background service |
| `com.prosoftnet.idsyncdaemon` | System Daemon | ❓ DECIDE | IDrive sync daemon |
| `com.zoho.truesync.upgrader` | System Daemon | ❓ DECIDE | Zoho TrueSync — depends on Zoho decision |

### A.5 — Login Items

| Item | Verdict | Notes |
|------|---------|-------|
| **Raycast** | ✅ KEEP | Launcher — should start on login |
| **Dropbox** | ❓ DECIDE | Depends on Dropbox decision |

### A.6 — Currently Running Processes (Notable)

Already captured above. Key background processes to be aware of:

- **Proton Drive + Proton Mail** — Running with helper processes
- **Zoho TrueSync** — Syncing in background  
- **IDrive** — Backup daemon + Finder plugin + monitor
- **Dropbox** — Sync + updater + file provider
- **Brave Browser** — 29.5 GB RAM hog (multiple renderer processes)
- **Circle** — Community app with helpers
- **Raindrop.io** — Bookmark manager

---

## Appendix B — Decision Summary

**To finalize before reset, answer these:**

1. **Dropbox** — Keep or replace with Proton Drive only?
2. **IDrive** — Keep or rely on Time Machine + git?
3. **Zoho WorkDrive** — Still needed?
4. **Cloud CLIs** — Which clouds are you actively using? (AWS ✅, Azure ?, GCP ?, DigitalOcean ?)
5. **Letta Code / Agent Browser** — Still experimenting or done?
6. **BambuStudio / Blender / OBS** — Active creative tools or dormant?
7. **ChatGPT / Perplexity desktop apps** — Or just use web versions?
8. **WhatsApp** — Desktop app or phone only?
9. **LiveKit / NATS / Terraform** — Active infrastructure tools or legacy?
10. **Raindrop.io** — Active bookmark workflow?
