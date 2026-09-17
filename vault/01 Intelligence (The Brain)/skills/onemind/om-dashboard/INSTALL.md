# Installing One Mind Dashboard Skill

## Phase 1: System Check
- [ ] VPS omos-core is reachable (`ssh -i ~/.ssh/id_ed25519 ubuntu@100.123.126.34`)
- [ ] OpenClaw gateway is running (`sudo openclaw channels status --probe`)
- [ ] `/home/ubuntu/.openclaw/plugin-skills/` directory exists
- [ ] Agent has `om-dashboard` in its skills list in `openclaw.json`

## Phase 2: Skill Location on VPS
Runtime copy lives at:
`/home/ubuntu/.openclaw/plugin-skills/om-dashboard/SKILL.md`

This vault document is the canonical source — keep both in sync when editing.

To push an update from vault to VPS:
```bash
scp -i ~/.ssh/id_ed25519 \
  "SKILL.md" \
  ubuntu@100.123.126.34:/home/ubuntu/.openclaw/plugin-skills/om-dashboard/SKILL.md
```
Then reload: `kill -HUP $(pgrep -f 'openclaw.*index.js')`

## Phase 3: Skill Activation
All 19 agents have `om-dashboard` in their skills list as of 2026-05-10.

To add to a new agent:
```bash
sudo -u ubuntu openclaw config get agents.list | python3 -c "..."
# Or edit openclaw.json directly — add "om-dashboard" to agent's skills array
```

## Phase 4: Cron Setup (Overwatch daily refresh)
```bash
sudo -u ubuntu openclaw cron add \
  --agent overwatch \
  --name "om-dashboard-daily" \
  --cron "0 9 * * *" \
  --tz "America/New_York" \
  --session isolated \
  --announce \
  --channel discord \
  --to 1502954950764793870 \
  --message "Run /om-dashboard gateway and save the output to the vault. Post a one-line summary here when done." \
  --model amazon-bedrock/us.anthropic.claude-haiku-4-5-20251001-v1:0
```

## Phase 5: Verification
Run `VERIFY.md` checklist.

## Confirmation
✅ One Mind Dashboard installed. Capabilities: generate live HTML dashboards for gateway, agents, cron, channels, vault, or pillar status. Saves to vault, announces to Discord. All 19 agents can invoke on demand.
