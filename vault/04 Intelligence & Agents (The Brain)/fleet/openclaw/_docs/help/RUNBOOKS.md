---
name: OpenClaw — Operational Runbooks
type: runbook
status: active
updated: 2026-05-11
owner: zeus
tags: [type/runbook, domain/SO-02, openclaw/help]
---

# Operational Runbooks

Six core procedures for maintaining the OpenClaw fleet on omos-core.

---

## 1. Fix Broken Workspace Backup Push

**When:** `git status` on VM shows diverged repo — `watch.sh` can't push.

```bash
ssh ubuntu@100.123.126.34
cd ~/.openclaw

# Inspect the divergence
git fetch origin
git log HEAD..origin/main --oneline    # what remote has
git log origin/main..HEAD --oneline    # what local has

# Rebase local commits on top of remote
git stash push -u -m "pre-rebase save"
git pull --rebase origin main
git stash pop
git push origin main

# Verify clean
git status
```

Expected: `Your branch is up to date with 'origin/main'`

---

## 2. Restore Agent Workspace from Backup

**When:** Agent workspace files are corrupted, missing, or need rollback.

```bash
ssh ubuntu@100.123.126.34

# Clone the backup repo
git clone https://github.com/OneMind-OS/openclaw-core-workspaces.git /tmp/ws-restore

# Restore — workspace-{id} agents (legacy, zenith, unity, eternal, sovereignty,
# spartan, oracle, eden, haven, empire, guardian, grid, trinity, overwatch, librarian):
cp /tmp/ws-restore/workspace-zenith/* ~/.openclaw/workspace-zenith/

# Restore — agents/{id}/workspace agents (sentinel, pulse, genesis, helix, nomad, apex, creed, eagleeye):
cp /tmp/ws-restore/workspace-apex/* ~/.openclaw/agents/apex/workspace/

# Restart gateway to reload
sudo systemctl restart openclaw.service

# Cleanup
rm -rf /tmp/ws-restore
```

---

## 3. Add a New Agent

```bash
ssh ubuntu@100.123.126.34
```

1. Edit `~/.openclaw/openclaw.json` → `agents.list` — add entry:
```json
{
  "id": "new-agent-id",
  "name": "Agent Name",
  "workspace": "~/.openclaw/workspace-new-agent-id",
  "skills": [
    "domain-skill",
    "onemind-maintainer", "taskflow", "tavily", "summarize",
    "telegram", "onemind-inbox", "onemind-sync"
  ]
}
```
2. Create workspace directory and populate 7 files (copy from existing agent as template)
3. Add Discord binding in `openclaw.json` → `bindings`
4. Run `openclaw gateway reload`
5. Create agent profile at `agents/{id}.md` in this vault
6. Create workspace folder at `workspaces (Agent Workspaces)/{Name} ({Name} Workspace)/`
7. Update `agents/README.md` tables and `help/README.md` What's Here table

---

## 4. Add or Modify a Cron Job

```bash
# Create a new cron job
ssh ubuntu@100.123.126.34 "openclaw cron create \
  --agent zenith \
  --name 'zenith-morning-brief' \
  --schedule '0 7 * * *' \
  --tz 'America/New_York' \
  --prompt 'Run morning HP brief'"

# List all jobs + IDs
ssh ubuntu@100.123.126.34 "openclaw cron list"

# Disable / enable
ssh ubuntu@100.123.126.34 "openclaw cron disable <job-id>"
ssh ubuntu@100.123.126.34 "openclaw cron enable <job-id>"

# Trigger manually to test
ssh ubuntu@100.123.126.34 "openclaw cron trigger <job-id>"

# View run history
ssh ubuntu@100.123.126.34 "openclaw cron runs <job-id>"
```

After adding, document in `cron (Scheduled Jobs)/CRON.md`.

---

## 5. Add a New MCP Server

1. SSH to omos-core
2. Edit `~/.openclaw/openclaw.json` → `mcp.servers`:
```json
"new-server": {
  "command": "npx",
  "args": ["-y", "@org/new-mcp-server"],
  "env": { "API_KEY": "op://aaeyris5z44pdqspk277u2oppi/Item Name/field" }
}
```
3. Reload: `ssh ubuntu@100.123.126.34 "openclaw gateway reload"`
4. Document in `mcp/README.md`

---

## 6. Rotate Gateway Token

The gateway auth token lives in `openclaw.json` → `gateway.auth.token`.

```bash
ssh ubuntu@100.123.126.34

# Generate new token
NEW_TOKEN=$(openssl rand -hex 24)
echo "New token: $NEW_TOKEN"

# Edit ~/.openclaw/openclaw.json — update gateway.auth.token
# Also update gateway.remote.token if used externally

openclaw gateway reload
```

Update any clients using the token (Claude Code MCP config, Tailscale proxy, etc.).

---

## Diagnostics

```bash
# Gateway health
ssh ubuntu@100.123.126.34 "openclaw gateway status"

# All active agents
ssh ubuntu@100.123.126.34 "openclaw agents list"

# Is watcher alive?
ssh ubuntu@100.123.126.34 "ps aux | grep watch.sh | grep -v grep"

# Backup repo sync status
ssh ubuntu@100.123.126.34 "git -C ~/.openclaw fetch origin && git -C ~/.openclaw status"

# Gateway logs (live)
ssh ubuntu@100.123.126.34 "journalctl -u openclaw.service -f"

# Watcher log
ssh ubuntu@100.123.126.34 "tail -30 /tmp/openclaw-backup.log"
```

## Related

- [[../workspaces (Agent Workspaces)/backup-repo|Backup Repo]] — Workspace backup repo docs
- [[../agents/README|Agent Registry]] — All 23 agents
- [[../cron (Scheduled Jobs)/CRON|Cron Jobs]] — Full cron schedule
- [[../mcp/README|MCP Servers]] — MCP server management
- [[../gateway/VPS-GATEWAY|VPS Gateway]] — Gateway config reference
