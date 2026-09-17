---
id: so-01-skill-onemind-report
name: onemind-report
description: Generate a One Mind system report — task counts by pillar, sync status, inbox depth, and infrastructure health. Use when Zeus asks for a "report", "system report", "how's the system", "give me a dashboard", or wants a pillar breakdown with stall detection using "report deep".
type: skill
domain: SO/01
agent: Legacy
status: active
updated: 2026-05-10
tags: [type/skill, domain/SO, domain/01-intelligence, skill/reporting]
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Skill: One Mind System Report

Generate a One Mind system status report. Default mode produces a fast dashboard (task counts, sync, infrastructure). Pass `deep` to run a full pillar breakdown with stall detection.

---

## Trigger Phrases

- "report", "system report", "how's the system", "/report", "/report deep"
- "give me a dashboard", "what's the state of One Mind"

---

## Step 1: Get date

```bash
date '+%A, %B %-d %Y'
```

---

## Step 2: Pull Todoist task summary

Load the Todoist token:

```bash
source ~/.secrets/onemind.env && echo $TODOIST_API_TOKEN | cut -c1-8
```

Fetch all active tasks:

```bash
source ~/.secrets/onemind.env && curl -s -X GET "https://api.todoist.com/rest/v2/tasks" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" | python3 -c "
import json, sys
from datetime import date

tasks = json.load(sys.stdin)
today = date.today().isoformat()
pillars = {'SO': 0, 'HP': 0, 'LE': 0, 'GE': 0, 'Uncategorized': 0}
overdue = 0

SO_KW = ['so','sovereign','infrastructure','security','research','automation','intelligence']
HP_KW = ['hp','holistic','body','fitness','mind','money','health','spartan','oracle']
LE_KW = ['le','legacy evolution','family','home','haven','eden']
GE_KW = ['ge','empire','business','brand','ventures','career','generational']

for t in tasks:
    pname = (t.get('project_id') or '').lower()
    if any(k in pname for k in SO_KW): p = 'SO'
    elif any(k in pname for k in HP_KW): p = 'HP'
    elif any(k in pname for k in LE_KW): p = 'LE'
    elif any(k in pname for k in GE_KW): p = 'GE'
    else: p = 'Uncategorized'
    pillars[p] += 1
    due = (t.get('due') or {}).get('date','')
    if due and due < today: overdue += 1

total = sum(pillars.values())
print(f'TOTAL={total} SO={pillars[\"SO\"]} HP={pillars[\"HP\"]} LE={pillars[\"LE\"]} GE={pillars[\"GE\"]} OVERDUE={overdue}')
"
```

---

## Step 3: Check inbox

```bash
ls "/root/onemind-codex/00-24 SO (Sovereign Operations)/99 Inbox (Queue) (The Dock)/" 2>/dev/null | grep "\.md$" | wc -l | tr -d ' '
```

---

## Step 4: Check sync status

```bash
cd /root/onemind-codex && \
  AHEAD=$(git rev-list --count origin/main..HEAD 2>/dev/null || echo 0) && \
  BEHIND=$(git rev-list --count HEAD..origin/main 2>/dev/null || echo 0) && \
  HEAD=$(git rev-parse --short HEAD 2>/dev/null || echo unknown) && \
  echo "AHEAD=$AHEAD BEHIND=$BEHIND HEAD=$HEAD"
```

---

## Step 5: Check infrastructure services

```bash
echo "WATCH=$(systemctl is-active onemind-watch.service 2>/dev/null || echo unknown)"
echo "PULL=$(systemctl is-active onemind-pull.timer 2>/dev/null || echo unknown)"
echo "UPTIME=$(uptime -p 2>/dev/null || echo unknown)"
```

---

## Step 6: Render dashboard

Parse all output above and render:

```
## 📊 One Mind Report — {date}

### 🎯 Execution
Todoist: {total} active | ⚠️ {overdue} overdue | 📥 Inbox: {inbox_count}

### 🏛 Pillar Snapshot
SO ■■■░░░ {n} tasks    {⚠️ stalled? if 0}
HP ■■░░░░ {n} tasks    {⚠️ stalled? if 0}
LE ░░░░░░ {n} tasks    {⚠️ stalled? if 0}
GE ■■■■░░ {n} tasks    {⚠️ stalled? if 0}

### 🔄 Sync
Git: {✅ in sync | ⚠️ {n} ahead / {n} behind} | HEAD: {sha}
onemind-watch: {active ✅ / inactive ❌}
onemind-pull: {active ✅ / inactive ❌}

### 🖥 Infrastructure
Uptime: {uptime}
VPS: omos-core (/root/onemind-codex)
```

Bar calculation: 1 bar per task, max 6 bars. Round down. All 4 pillars always shown (show ⚠️ stalled? if 0 tasks).

---

## Step 7: Deep mode (only if triggered with "deep")

Warn first: "Running deep analysis across all pillars..."

### Deep — SO Analysis

```bash
find "/root/onemind-codex/00-24 SO (Sovereign Operations)" -name "README.md" -not -path "*/.obsidian/*" | head -15 | while read f; do
  echo "=== $f ==="
  grep -A 20 "## ⚡ Active Tasks" "$f" 2>/dev/null | head -15
done
```

```bash
find "/root/onemind-codex/00-24 SO (Sovereign Operations)" -name "README.md" -not -path "*/.obsidian/*" -mtime +7 | head -10
```

### Deep — HP Analysis

Same pattern for `/root/onemind-codex/25-49 HP (Holistic Performance)`.

### Deep — LE Analysis

Same pattern for `/root/onemind-codex/50-74 LE (Legacy Evolution)`.

### Deep — GE Analysis

Same pattern for `/root/onemind-codex/75-99 GE (Generational Empire)`.

### Deep render addition

Append after dashboard:

```
---
## 🔬 Deep Analysis

### ⚡ SO — Sovereign Operations
Active tasks in READMEs: {list or "none found"}
Stale READMEs (7+ days no update): {list or "none"}

### 💪 HP — Holistic Performance
...

### 🏡 LE — Legacy Evolution
...

### 💰 GE — Generational Empire
...

### 🚨 Flags
{Pillar with 0 Todoist tasks AND stale READMEs = likely stalled}
{Tasks in vault READMEs but not in Todoist = sync drift}
```

---

## Rules

- Always show the dashboard first, even in deep mode
- Deep mode is slow (~40 README reads) — always warn before starting
- Stale = `find -mtime +7` (file not modified in 7+ days)
- 0-task pillar always gets ⚠️ stalled? flag
- Keep the dashboard under 20 lines — tight and scannable
- Do NOT commit anything — report is read-only
