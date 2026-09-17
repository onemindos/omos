# 5 Essential Cron Automations

> *The 5 cron jobs that make your agent stack proactive.*
> *Add these to `openclaw.json` under the `"cron"` array.*

---

## Automation 1 — Morning Brief (Daily 6:00 AM)

```json
{
  "id": "morning-briefing",
  "schedule": "0 6 * * *",
  "agent": "main",
  "channel": "discord",
  "account": "sentinel",
  "prompt": "You are Legacy ⚡. Generate the daily morning brief for Zeus. Structure: 🌅 TODAY'S FOCUS (top priority from current goals), ⚡ SYSTEM STATUS (agent stack health), 📋 TOP 5 TASKS (from active task list), 🧠 ONE INSIGHT (specific to Zeus's current context), 🔔 FLAGS (anything unresolved or time-sensitive). Be concise — Zeus reads this in 3 minutes.",
  "enabled": true
}
```

**Why it matters:** Zeus starts every day oriented. No more "what was I working on?"

---

## Automation 2 — EOD Review (Daily 8:00 PM)

```json
{
  "id": "eod-review",
  "schedule": "0 20 * * *",
  "agent": "main",
  "channel": "discord",
  "account": "sentinel",
  "prompt": "You are Legacy ⚡. Generate the end-of-day review. Structure: ✅ WINS (what moved today), 🔄 CARRY FORWARD (what to reschedule vs. drop), 📊 DAY QUALITY (1-10 + one honest sentence), 🔔 TOMORROW'S SETUP (single most important MIT). Under 200 words total.",
  "enabled": true
}
```

**Why it matters:** The day closes deliberately. Nothing carries forward accidentally.

---

## Automation 3 — Weekly Report (Sundays 5:00 PM)

```json
{
  "id": "weekly-report",
  "schedule": "0 17 * * 0",
  "agent": "main",
  "channel": "discord",
  "account": "sentinel",
  "prompt": "You are Legacy ⚡. Generate the full weekly intelligence report. This feeds Zeus's Sunday review. Cover: 🏆 TOP WINS, 📋 TASKS COMPLETED by pillar, 🎯 GOAL PROGRESS (Green/Yellow/Red per goal), ⚠️ OPEN LOOPS, 📈 PATTERNS, 🔭 NEXT WEEK SETUP (top 3 priorities), 🤖 AGENT ACTIVITY. Be comprehensive — this replaces scrolling a week of messages.",
  "enabled": true
}
```

**Why it matters:** Zeus walks into his weekly review with full context already prepared.

---

## Automation 4 — Monthly Goal Audit (1st of Month, 7:00 AM)

```json
{
  "id": "monthly-goal-audit",
  "schedule": "0 7 1 * *",
  "agent": "main",
  "channel": "discord",
  "account": "sentinel",
  "prompt": "You are Legacy ⚡. It's the 1st of the month. Run the monthly goal audit. For each active goal in Zeus's context: (1) what was the target for last month? (2) what happened? (3) On track / Behind / Complete? (4) What's the single most important action this month to move it forward? Present as a table then 3 key decisions Zeus needs to make this month.",
  "enabled": true
}
```

**Why it matters:** Goals don't slip silently. The first of every month triggers a mandatory audit.

---

## Automation 5 — Empire Weekly Update (Mondays 9:00 AM)

```json
{
  "id": "empire-weekly",
  "schedule": "0 9 * * 1",
  "agent": "empire",
  "channel": "discord",
  "account": "sentinel",
  "prompt": "You are Empire 💰. Generate the weekly empire update. Cover: 💰 REVENUE (what came in last week, trend), 📈 ACTIVE VENTURES (one sentence status per active venture), 🎯 BRAND ACTIVITY (content created, audience growth, notable interactions), 🏗️ THIS WEEK'S EMPIRE PRIORITIES (top 3 GE actions), 🚨 FLAGS (any GE item at risk or needing Zeus's decision). Direct and numbers-first.",
  "enabled": true
}
```

**Why it matters:** Empire gets dedicated attention every Monday. Business doesn't slide into the background.

---

## Schedule Overview

| Automation | Schedule | Agent | Channel |
|------------|----------|-------|---------|
| Morning Brief | Daily 6:00 AM | Legacy | Discord |
| EOD Review | Daily 8:00 PM | Legacy | Discord |
| Weekly Report | Sunday 5:00 PM | Legacy | Discord |
| Monthly Goal Audit | 1st of month 7:00 AM | Legacy | Discord |
| Empire Weekly | Monday 9:00 AM | Empire | Discord |

---

## Implementation

1. Open `~/.openclaw/openclaw.json`
2. Find the `"cron"` array (or add one if not present)
3. Add the JSON objects above inside the array
4. Restart OpenClaw: `systemctl --user restart openclaw-gateway`
5. Verify: `openclaw cron list` — should show all 5 jobs

To test any cron immediately (without waiting for schedule): `openclaw cron run [id]`

---

*Part of the One Mind OS OpenClaw Deep Dive — Course 3*
