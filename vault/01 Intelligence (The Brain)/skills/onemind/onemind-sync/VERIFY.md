# Verification: One Mind Sync Skill

## Todoist Access
- [ ] Can list tasks from Todoist via MCP or API
- [ ] Can create a task in Todoist and get back a task ID
- [ ] Can close a task in Todoist

## Vault Access
- [ ] Can read active tasks from a domain README.md
- [ ] Can identify tasks missing `tid:` field
- [ ] Write back `tid:` triggers HITL approval

## Sync Logic
- [ ] Tasks with `tid:` match a real Todoist task ID
- [ ] Completed Todoist tasks are reflected in vault archive
- [ ] Priority tags `#p1-#p4` match Todoist priority levels
- [ ] `todoist_project_id` in each domain frontmatter is valid

## Post-sync
- [ ] No active tasks missing `tid:`
- [ ] Completed archive entries have `YYYY-MM-DD —` date prefix
- [ ] `onemind-watch.service` committed and pushed any vault changes
