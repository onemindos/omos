---
type: agent-heartbeat
agent: spartan
status: active
updated: 2026-05-03
---
# HEARTBEAT.md — Spartan

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Heartbeat | Every 1h (6am-10pm) | Physical status check |
| Morning check-in | 6:00am | Sleep score, today's workout, daily priorities |
| Post-workout log | Event-driven | Log sets, reps, weight, RPE |
| Evening review | 9:00pm | Day's metrics, nutrition summary, recovery notes |
| Weekly PI score | Sunday 8:00am | Composite performance index report |
| Monthly program review | 1st of month | Progress assessment, program adjustments |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Health emergency indicators | Immediate alert to Zeus | L0 |
| Injury signs during training | Stop workout + injury protocol | L0 |
| PI score drops >20% | Alert Zeus with analysis | L1 |
| Missed 3+ consecutive workouts | Check-in with Zeus | L1 |
| Sleep quality <5/10 for 3+ days | Sleep intervention protocol | L2 |
| Supplement/medication interaction | Alert + hold recommendation | L1 |
| PR achieved | Log + celebrate | L3 |
