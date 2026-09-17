# Verification: One Mind Inbox Skill

## Inbox State
- [ ] Can list inbox: `ls /home/ubuntu/onemind/00-24\ SO\ \(Sovereign\ Operations\)/99\ Inbox\ \(Queue\)/`
- [ ] Item count matches hook report (UserPromptSubmit shows "Inbox items waiting: N")

## Routing
- [ ] Each item correctly classified to SO / HP / LE / GE
- [ ] No items left unrouted without Zeus flag

## HITL Gate
- [ ] Move triggers exec-approval before executing
- [ ] Delete triggers exec-approval before executing
- [ ] Read of inbox items requires no approval

## Post-triage
- [ ] Routed items have correct frontmatter at destination
- [ ] Tasks have `tid:` populated after Todoist creation
- [ ] Inbox item deleted after successful routing
- [ ] `onemind-watch.service` committed and pushed changes
