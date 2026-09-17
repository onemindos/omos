## Direct Integration Setup Progress

**STATUS: 2026-06-02 03:52 AM**

### 1. TODOIST ✓ TOKEN UPDATED, ⚠️ NEEDS SESSION RESTART
- New token: b559882a309bae2a7a4b853e5ad919ae5dbd5386
- Saved to ~/.hermes/.env
- **ACTION REQUIRED:** Restart this Hermes session for MCP to pick up new token
  ```bash
  # Exit this session, then:
  hermes chat
  ```
- After restart, test: `mcp_todoist_get_projects_list`

### 2. FASTMAIL ✓ TOKEN SAVED, ⚠️ NEEDS MCP CONFIG
- Token: fmu1-5b7643a6-1060e13d013bbafc1c65b92ff145692b-0-8026b0e41a4e2a9208424757ee9473ae
- Saved to ~/.hermes/.env as FASTMAIL_MCP_TOKEN
- **ACTION REQUIRED:** Add to ~/.hermes/config.yaml under mcp_servers:
  ```yaml
  fastmail:
    transport: http
    url: https://api.fastmail.com/mcp
    headers:
      Authorization: Bearer ${FASTMAIL_MCP_TOKEN}
  ```
- Docs: https://www.fastmail.help/hc/en-us/articles/15869557281295

### 3. PROTON MAIL ⚠️ BRIDGE SETUP IN PROGRESS
- Proton Mail Bridge v3.24.2 installed
- Bridge GUI launched (PID 8405, 8377)
- **ACTION REQUIRED (ZEUS — INTERACTIVE):**
  
  **Step 1:** Bridge should show a window asking you to sign in
  - Click "Sign in" or "Add account"
  - Use your Proton Mail credentials (from 1Password SO vault)
  - Complete 2FA if prompted
  
  **Step 2:** Once logged in, Bridge will show IMAP/SMTP settings:
  - IMAP host: 127.0.0.1
  - IMAP port: (Bridge will display — usually 1143)
  - SMTP host: 127.0.0.1
  - SMTP port: (Bridge will display — usually 1025)
  - Username: (your Proton email)
  - Password: (Bridge-generated password — COPY THIS)
  
  **Step 3:** Copy the Bridge password, then run:
  ```bash
  himalaya account configure
  ```
  
  When prompted:
  - Account name: `proton`
  - Email: (your Proton address)
  - IMAP host: 127.0.0.1
  - IMAP port: (from Bridge, usually 1143)
  - SMTP host: 127.0.0.1
  - SMTP port: (from Bridge, usually 1025)
  - Auth: password
  - Password command: `echo "BRIDGE_PASSWORD_HERE"` (paste Bridge password)
  
  **Step 4:** Test:
  ```bash
  himalaya envelope list
  ```

### 4. OBSIDIAN ✓ FULLY OPERATIONAL
- Path: ~/Documents/onemind/
- Filesystem MCP active
- No action needed

---

**NEXT STEPS:**
1. Zeus: Sign into Proton Bridge GUI → copy IMAP/SMTP credentials
2. Zeus: Run `himalaya account configure` with Bridge settings
3. Zeus: Manually add Fastmail MCP block to ~/.hermes/config.yaml
4. Zeus: Restart Hermes session (`exit` then `hermes chat`)
5. Legacy: Test all integrations after restart
