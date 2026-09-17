---
type: agent-heartbeat
agent: eden
status: active
updated: 2026-05-03
---
# HEARTBEAT.md — Eden

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Weather check | 6:00am + 6:00pm | Frost alerts, rain, extreme heat |
| Heartbeat | Every 3h (6am-8pm) | Garden and livestock status |
| Daily garden log | 7:00pm | What was done, observed, harvested |
| Weekly harvest summary | Sunday 6:00pm | Total yields, variety performance |
| Seasonal planning | 30 days before season | Planting calendar, seed orders |
| Monthly soil check | 15th of month | Moisture, amendments, cover crop status |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Frost warning | Immediate alert + protection plan | L1 |
| Livestock health concern | Alert Zeus + care protocol | L1 |
| Pest outbreak | Identify + IPM recommendation | L2 |
| Equipment failure | Alert + workaround plan | L2 |
| Harvest ready | Notify with preservation plan | L3 |
| Sensor anomaly (future) | Alert Eden + Overwatch | L2 |
