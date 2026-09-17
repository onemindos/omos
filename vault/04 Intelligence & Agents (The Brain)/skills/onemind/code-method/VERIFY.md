# Verifying CODE× Method Skill

## Required Checks

- [ ] **Agent confirms the deletion rule: nothing is deleted without Zeus's explicit approval — not duplicates, not empty folders, not stale files, not anything**
- [ ] Agent can describe all 5 CODE× stages in one sentence each
- [ ] Agent can name all 4 domain groups and their number ranges
- [ ] Agent can route a content item to the correct domain without looking it up
- [ ] Agent knows the folder naming convention (`Name (Subtitle)/`)
- [ ] Agent knows the priority system (P1-P4) and agent assignment map
- [ ] Agent knows what to log after every execution

## Capture Test

Given: Zeus shares `https://www.instagram.com/reel/XXXXX/`

Agent must:
- [ ] Detect as URL → Instagram reel
- [ ] Extract via oembed
- [ ] Classify domain based on content
- [ ] Present: route + summary + optional task
- [ ] Execute on approval: write note to vault, create task if needed

## Organize Test

Given these items, agent correctly identifies target domain + subfolder:

| Item | Expected Target |
|------|----------------|
| Zeus's personal credit score report | `HP / 30 Finance / Credit` |
| JV partnership agreement | `GE / 82 Partnerships / Strategic-Partners / Agreements` |
| New keynote speech script | `GE / 75 / Speaking / Keynotes` |
| n8n automation workflow | `SO / 07 Automations / n8n` |
| Zeus's morning routine doc | `HP / 32 Rhythm / Routines-Habits` |
| Legacy Restore LLC operating agreement | `GE / 78 Holdings / Legal / Operating-Agreements` |
| New business idea not yet validated | `GE / 80 Innovation / Ideas / Raw-Ideas` |
| Real estate property closing doc | `GE / 76 Real Estate / Legal / Deeds` |
| Travel itinerary | `HP / 34 Joy / Travel` |
| New VPS server setup guide | `SO / 02 Infrastructure / nodes` |

## Direct Test

Given a list of 10 unranked items and `active-goals.md`, agent must:
- [ ] Assign P1-P4 to each item
- [ ] Assign each item to the correct specialist agent
- [ ] Identify any cross-pillar items
- [ ] Build an execution queue capped at 5 P1 items

## Execute Test

Given a prioritized task, agent must:
- [ ] Load context before starting
- [ ] Identify task type (CREATE / UPDATE / RESEARCH / COMMUNICATE / AUTOMATE)
- [ ] Save output to the correct domain with proper frontmatter
- [ ] Create action log entry
- [ ] Capture any new tasks discovered during execution (loop back to C)

## Reflect Test

Given a completed work session, agent must:
- [ ] Produce a daily log with lessons learned
- [ ] Identify at least one Insight, Update, or Action from the session
- [ ] Know where to file each type of reflection output

## Result

✅ All checks pass = skill fully operational
⚠️ Routing test errors = re-read SKILL.md Organize section
⚠️ Priority test errors = re-read SKILL.md Direct section
❌ Required checks fail = not ready, re-run INSTALL.md


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
