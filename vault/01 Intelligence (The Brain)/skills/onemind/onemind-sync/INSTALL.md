# Installing One Mind Sync Skill

## Phase 1: System Check
- [ ] `onemind-maintainer` skill installed (vault writes need HITL)
- [ ] Todoist MCP or API access available
- [ ] Each domain README.md has `todoist_project_id` in frontmatter

## Phase 2: Skill Location on VPS
`/home/ubuntu/.openclaw/plugin-skills/onemind-sync/SKILL.md`

## Phase 3: Skill Activation
1. Read `SKILL.md` in this directory
2. Internalize the `tid:` bridge concept and sync rules table
3. Confirm Todoist access before running any sync

## Phase 4: Verification
Run `VERIFY.md` checklist.

## Confirmation
✅ One Mind Sync installed. Workflow: scan vault tasks → check Todoist state → create missing `tid:` links → archive completed tasks. All vault writes require HITL approval.
