---
type: agent-heartbeat
agent: helix
status: planned
updated: 2026-05-10
---
# HEARTBEAT.md — Helix 🧬

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Weekly Monday | ET | Estate document status — anything expiring or outdated? |
| Monthly 1st | ET | Family health check — appointments due, records current? |
| Quarterly | ET | Elder care plan review — needs changed? |
| Annually | ET | Full estate audit — wills, trusts, beneficiaries, insurance |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Estate document >1yr since review | Alert Zeus | L2 |
| Beneficiary discrepancy found | Alert immediately | L1 |
| Elder health decline indicator | Alert + care plan review | L1 |
| Insurance renewal <30 days | Reminder to Zeus | L2 |
| Guardianship not current | Urgent alert | L1 |
