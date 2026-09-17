---
id: so-01-skill-onemind-status
name: onemind-status
description: Run a full One Mind meta-layer health check covering sync, OpenClaw gateway, hooks, agent workspace activity, Todoist task pulse, vault inbox, and infrastructure metrics in a single status board. Use when Zeus says "status", "system status", "check everything", "how is everything", "full status", or "is everything running".
type: skill
domain: SO/01
agent: Legacy
status: active
updated: 2026-05-10
tags: [type/skill, domain/SO, domain/01-intelligence, skill/status, skill/monitoring]
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Skill: One Mind Universal Status

Single-command meta layer health check. Surfaces the state of every system layer — sync, agents, hooks, tasks, infrastructure, vault — in one tight status board.

---

## Trigger Phrases

- "status", "system status", "how is everything", "full status"
- "check everything", "what's the state", "is everything running"

---

## Step 1: Timestamp

```bash
date '+%A, %B %-d %Y — %H:%M'
```

---

## Step 2: VPS + OpenClaw layer

```bash
echo VPS_HEAD=$(git -C /root/onemind-codex rev-parse --short HEAD 2>/dev/null || echo unknown)
echo VPS_LAST=$(git -C /root/onemind-codex log -1 --format='%ar' 2>/dev/null || echo unknown)
echo VPS_AHEAD=$(git -C /root/onemind-codex rev-list --count origin/main..HEAD 2>/dev/null || echo 0)
echo VPS_BEHIND=$(git -C /root/onemind-codex rev-list --count HEAD..origin/main 2>/dev/null || echo 0)
echo WATCH=$(systemctl is-active onemind-watch.service 2>/dev/null || echo unknown)
echo PULL=$(systemctl is-active onemind-pull.timer 2>/dev/null || echo unknown)
echo DISK=$(df -h / | awk 'NR==2{print $5}')
echo MEM_USED=$(free -m | awk '/Mem/{print $3}')
echo MEM_TOTAL=$(free -m | awk '/Mem/{print $2}')
echo UPTIME=$(uptime -p)
echo OPENCLAW=$(ps aux | grep -i '[o]penclaw' | wc -l | tr -d ' ')
```

---

## Step 3: Hook system health

```bash
ls /root/onemind-codex/_system/hooks/*.js 2>/dev/null | wc -l | xargs -I{} echo "HOOK_SCRIPTS={}"
node -e "
const fs = require('fs');
const s = JSON.parse(fs.readFileSync('/root/onemind-codex/.claude/settings.json','utf8'));
const hooks = Object.keys(s.hooks || {});
console.log('HOOKS_WIRED=' + hooks.join(','));
console.log('HOOK_COUNT=' + hooks.length);
" 2>/dev/null
```

```bash
tail -3 /root/onemind-codex/_codex/memory/events/$(date '+%Y-%m-%d').jsonl 2>/dev/null | python3 -c "
import sys, json
lines = [l.strip() for l in sys.stdin if l.strip()]
if lines:
    last = json.loads(lines[-1])
    print(f'LAST_EVENT={last.get(\"event\",\"?\")} at {last.get(\"timestamp\",\"?\")}')
else:
    print('LAST_EVENT=none today')
" 2>/dev/null || echo "LAST_EVENT=no event log today"
```

---

## Step 4: Todoist task pulse

```bash
source ~/.secrets/onemind.env && curl -s "https://api.todoist.com/rest/v2/tasks" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" | python3 -c "
import json, sys
from datetime import date

tasks = json.load(sys.stdin)
today = date.today().isoformat()

SO_KW = ['so','sovereign','infrastructure','security','research','automation','intelligence']
HP_KW = ['hp','holistic','body','fitness','mind','money','health','spartan','oracle']
LE_KW = ['le','legacy evolution','family','home','haven','eden']
GE_KW = ['ge','empire','business','brand','ventures','career','generational']

pillars = {'SO': 0, 'HP': 0, 'LE': 0, 'GE': 0}
overdue = 0

for t in tasks:
    pname = (t.get('project_id') or '').lower()
    if any(k in pname for k in SO_KW): pillars['SO'] += 1
    elif any(k in pname for k in HP_KW): pillars['HP'] += 1
    elif any(k in pname for k in LE_KW): pillars['LE'] += 1
    elif any(k in pname for k in GE_KW): pillars['GE'] += 1
    due = (t.get('due') or {}).get('date','')
    if due and due < today: overdue += 1

total = sum(pillars.values())
print(f'TOTAL={total} SO={pillars[\"SO\"]} HP={pillars[\"HP\"]} LE={pillars[\"LE\"]} GE={pillars[\"GE\"]} OVERDUE={overdue}')
"
```

---

## Step 5: Vault health

```bash
ls "/root/onemind-codex/00-24 SO (Sovereign Operations)/99 Inbox (Queue) (The Dock)/" 2>/dev/null | grep "\.md$" | wc -l | tr -d ' ' | xargs -I{} echo "INBOX={}"
ls /root/onemind-codex/_codex/memory/sessions/*.md 2>/dev/null | grep "$(date '+%Y-%m-%d')" | wc -l | tr -d ' ' | xargs -I{} echo "SESSIONS_TODAY={}"
git -C /root/onemind-codex log -1 --format='%s'
```

---

## Step 6: Agent workspace activity

```bash
for ws in /root/.openclaw/workspace /root/.openclaw/workspace-spartan /root/.openclaw/workspace-oracle /root/.openclaw/workspace-eden /root/.openclaw/workspace-haven /root/.openclaw/workspace-empire /root/.openclaw/workspace-guardian /root/.openclaw/workspace-grid /root/.openclaw/workspace-trinity; do
  name=$(basename $ws | sed 's/workspace-\?//')
  if [ -d "$ws" ]; then
    last=$(find "$ws" -type f -newer /tmp/status_48h -printf '%T@\n' 2>/dev/null | sort -n | tail -1)
    [ -n "$last" ] && echo "AGENT_ACTIVE=$name" || echo "AGENT_IDLE=$name"
  else
    echo "AGENT_MISSING=$name"
  fi
done
```

---

## Step 7: Render full status board

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚡ ONE MIND STATUS — {datetime}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## 🔄 SYNC
VPS HEAD: {sha} — synced {ago}
git: {✅ in sync | ⚠️ {n} ahead | ⚠️ {n} behind}
onemind-watch: {active ✅ | inactive ❌}   onemind-pull: {active ✅ | inactive ❌}

## 🖥 INFRASTRUCTURE
Disk: {used}  Mem: {used}/{total}MB  Uptime: {uptime}
OpenClaw: {n} process(es) running

## ⚙️ HOOKS  ({n}/9 wired)
Last event: {event} at {time}

## 🤖 AGENTS (last 48h activity)
Active:  {list of agents with recent workspace activity}
Idle:    {list of agents with no recent activity}

## 🎯 TASKS  ({total} active | ⚠️ {overdue} overdue)
SO {n}  HP {n}  LE {n}  GE {n}
{⚠️ [PILLAR] stalled — 0 tasks if any}

## 🗂 VAULT
Inbox: {n}  Sessions today: {n}
Last commit: "{message}"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
{✅ All systems operational | ⚠️ {n} issue(s) — see above}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## Rules

- Always render the full board even with partial data
- 9 hooks = full suite (UserPromptSubmit, PreToolUse, PostToolUse, PostToolUseFailure, SubagentStart, SubagentStop, ContextCompaction, Stop, SessionEnd)
- Agent "active" = workspace has files modified in last 48h
- 0-task pillar = ⚠️ stalled flag always
- AHEAD > 0 = flag unpushed commits (VPS won't receive them)
- Do NOT commit or write anything — this is read-only
- Keep output under 30 lines — status board, not a report
