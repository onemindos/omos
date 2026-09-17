---
type: agent-heartbeat
agent: eagleeye
status: planned
updated: 2026-05-10
---
# HEARTBEAT.md — Eagleeye 🦅

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Daily 7am | ET | Mission queue check — scheduled flights, workflow status |
| Weekly Monday | ET | Equipment status — battery health, calibration due? |
| Weekly Friday | ET | n8n workflow audit — any failures or degraded flows? |
| Monthly 1st | ET | Full mission log review, maintenance schedule |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Drone battery critical | Ground mission, alert Zeus | L2 |
| n8n workflow failure | Triage + alert Grid/Zeus | L1 |
| Computer vision anomaly detected | Alert Guardian + Zeus | L1 |
| Hardware malfunction | Ground + alert Zeus immediately | L1 |
| Unauthorized perimeter intrusion detected | Alert Guardian immediately | L0 |
