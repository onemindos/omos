## NATIVE HERMES INTEGRATION AUDIT
**2026-06-02 04:15 AM**

### ✓ FULLY NATIVE (Hermes MCP):
1. **Todoist** — MCP server @ ~/.openclaw/npm/node_modules/.bin/todoist-mcp
   - 36 tools (tasks, projects, sections, labels, comments)
   - Status: Token configured, needs gateway restart test
   
2. **Fastmail** — MCP server @ https://api.fastmail.com/mcp  
   - 19 tools (email, calendar, contacts)
   - Status: Configured in config.yaml
   - Account: zeuslegacy@fastmail.com
   
3. **Obsidian** — Filesystem MCP server
   - Path: ~/Documents/onemind/
   - Status: Operational

### ✗ NOT NATIVE (External Dependencies):
4. **Proton Mail** — Uses Himalaya CLI (external tool)
   - Bridge running (127.0.0.1:1143/1025)
   - Account: zeus.vaults@pm.me (trust anchor)
   - **ISSUE**: No Proton MCP server exists
   - **OPTIONS**:
     a) Remove Himalaya, use Proton web UI manually (recovery-only)
     b) Keep Himalaya for CLI access (external dependency)
     c) Build custom Proton MCP wrapper (engineering effort)

---

### RECOMMENDATION:

**EMAIL STRATEGY:**
- **Fastmail (zeuslegacy@fastmail.com)** → Primary operational email via Hermes MCP
  - Read/send/search email
  - Calendar management
  - Contact management
  
- **Proton Mail (zeus.vaults@pm.me)** → Trust anchor, recovery-only
  - Manual access via web UI or Proton Mail.app
  - DNS-independent fallback
  - Never exposed publicly
  - Remove Himalaya dependency

**RATIONALE:**
- Fastmail MCP is native, supports full CRUD
- Proton's value is as a protected recovery channel, not operational email
- Removing Himalaya eliminates external dependency
- Proton Bridge stays installed for manual access if needed

### CURRENT MCP SERVERS (all native):
- nats
- cua-driver (macOS computer use)
- peekaboo
- todoist ✓
- filesystem ✓
- context7
- home-assistant
- fastmail ✓
- zoho-mail
- wwv
- telegram

### NEXT ACTIONS:
1. Test Todoist MCP after gateway restart
2. Test Fastmail MCP (email read/send)
3. Decision: Keep or remove Himalaya for Proton
