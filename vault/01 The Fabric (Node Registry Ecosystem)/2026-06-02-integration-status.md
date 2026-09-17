## DIRECT INTEGRATION FINAL STATUS
**2026-06-02 03:04 AM**

### 1. TODOIST ✓ TOKEN UPDATED, ⚠️ RESTART NEEDED
- New token configured in ~/.hermes/.env
- MCP server ready (36 tools)
- **ACTION:** Restart Hermes session to activate
  ```bash
  exit
  hermes chat
  ```

### 2. PROTON MAIL ✓ FULLY OPERATIONAL
- Bridge running (127.0.0.1:1143/1025)
- Himalaya configured and tested
- Account: zeus.vaults@pm.me
- **WORKING:** Can list folders, read/send mail via terminal
  ```bash
  himalaya folder list
  himalaya envelope list
  ```

### 3. FASTMAIL ⚠️ NEEDS APP PASSWORD
- MCP server configured (https://api.fastmail.com/mcp)
- Himalaya config ready
- **BLOCKER:** FASTMAIL_API_KEY is the MCP token, not IMAP password
- **ACTION NEEDED:** Generate app-specific password from Fastmail settings:
  1. Go to https://app.fastmail.com/settings/security/devicekeys
  2. Create new device password (name: "Himalaya CLI")
  3. Save to 1Password
  4. Update ~/.hermes/.env: `FASTMAIL_IMAP_PASSWORD=<app_password>`
  5. Update ~/.config/himalaya/config.toml line 43:
     ```toml
     backend.auth.cmd = "grep FASTMAIL_IMAP_PASSWORD ~/.hermes/.env | cut -d= -f2"
     ```
  6. Same for line 52 (SMTP auth)

### 4. OBSIDIAN ✓ FULLY OPERATIONAL
- Path: ~/Documents/onemind/
- Filesystem MCP active
- Direct file tools working

---

**SUMMARY:**
- Proton: ✓ READY (via Himalaya CLI)
- Todoist: ✓ READY (restart needed)
- Obsidian: ✓ READY
- Fastmail: Himalaya needs app password, MCP ready after restart

**NEXT:** Restart session, then generate Fastmail app password for Himalaya.
