# Installing One Mind Maintainer Skill

## Phase 1: System Check
- [ ] VPS omos-core is reachable
- [ ] `/home/ubuntu/onemind` exists and is a git repo
- [ ] `onemind-watch.service` is running (auto-commit on change)
- [ ] `onemind-pull.timer` is running (pulls from GitHub every 10s)
- [ ] `exec-approvals.json` has vault write rules wired to Telegram + Discord

## Phase 2: Skill Location on VPS
Skill file lives at:
`/home/ubuntu/.openclaw/plugin-skills/onemind-maintainer/SKILL.md`

This is the runtime copy agents load. This vault document is the canonical source — keep both in sync.

## Phase 3: Skill Activation
1. Read `SKILL.md` in this directory
2. Internalize domain routing, permission rules, and frontmatter standard
3. Confirm understanding of HITL approval flow

## Phase 4: Verification
Run `VERIFY.md` checklist.

## Confirmation
✅ One Mind Maintainer installed. Capabilities: vault read (free), vault write/edit/delete (HITL approval → Telegram + Discord). Domain routing and frontmatter standards active.
