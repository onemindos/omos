# Verification: One Mind Maintainer Skill

Run these checks after install to confirm the skill is operating correctly.

## Infrastructure
- [ ] `systemctl status onemind-watch.service` → active (running) on omos-core
- [ ] `systemctl status onemind-pull.timer` → active (running) on omos-core
- [ ] `ls /home/ubuntu/onemind/` → vault folders present (00-24 SO, 25-49 HP, etc.)
- [ ] `git -C /home/ubuntu/onemind remote -v` → points to OneMind-OS/One-Mind

## Permissions
- [ ] Reading a vault file requires no approval
- [ ] Attempting to write a file triggers exec-approval pause
- [ ] Approval request appears in Telegram AND Discord
- [ ] Approving → file is written
- [ ] Rejecting → agent stops and surfaces reason

## Auto-commit
- [ ] After an approved write, `onemind-watch.service` commits within 30s
- [ ] Commit message format: `agent: <filename>`
- [ ] Push to GitHub succeeds

## Skill Integrity
- [ ] Agent can read `/home/ubuntu/onemind/AGENTS.md`
- [ ] Agent routes new content to correct pillar (SO/HP/LE/GE)
- [ ] Frontmatter matches SCHEMA.md standard
- [ ] Terminology rules followed (no "quadrant", "Codex Method", etc.)
