# Rename onemind → zeus — Execution Plan
Created: June 3, 2026 — Paused, resume when machine is free

## Why
- macOS username is `onemind` — conflicts with the vault folder also called `onemind`
- zeus = the operator (you), onemind = the system you run
- Clean path after rename: /Users/zeus/onemind/

## Final Target Structure
  /Users/zeus/               ← you (the operator)
    onemind/                 ← THE vault (Obsidian + git)
    dev/                     ← all code repos
    Screenshots/             ← screenshots
    + macOS defaults

## Machine Naming Convention (already set, no change needed)
  Username:  zeus          → /Users/zeus/
  Vault:     onemind       → /Users/zeus/onemind/
  Hostname:  omos-mbp      → Tailscale/Headscale fleet identity

---

## Execution Steps

### STEP 1 — Create temp admin (while logged in as onemind)
  System Settings > Users & Groups > Add Account
  Type: Administrator
  Name: migrate
  Password: anything temporary

### STEP 2 — Log out of onemind
  Apple menu > Log Out

### STEP 3 — Log in as migrate

### STEP 4 — Run rename script
  sudo bash /Users/onemind/rename-to-zeus.sh

### STEP 5 — Log out of migrate, log in as zeus

### STEP 6 — Run post-login cleanup
  bash ~/post-login-cleanup.sh

### STEP 7 — Open Hermes and tell Legacy AI to resume
  Say: "resume the zeus rename — run post-rename verification"
  Legacy will handle:
  - Verify all paths updated /Users/onemind → /Users/zeus
  - Update Hermes config.yaml
  - Update all docker-compose volume paths (nextcloud, ocis, uptime-kuma, homepage)
  - Update Headscale + Headplane configs
  - Update OpenClaw configs
  - Update LaunchAgents/LaunchDaemons
  - Update Filesystem MCP scope
  - Update memory with new paths
  - Delete migrate account

---

## Scripts Location
  ~/rename-to-zeus.sh        ← run as migrate in Step 4
  ~/post-login-cleanup.sh    ← run as zeus in Step 6

## DO NOT proceed without a free machine window — takes ~10 min, requires full logout.
