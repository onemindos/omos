---
name: onemind-inbox
description: Process and clear the One Mind inbox at 00-24 SO/99 Inbox (Queue)/. Routes every new capture to the correct domain using pillar routing rules. Use during daily or weekly inbox triage.
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Codex Skill: One Mind Inbox

Process and clear the One Mind inbox. Every new capture lands in `00-24 SO/99 Inbox (Queue)/` first — this skill routes it to the correct domain.

---

## Trigger

Run when inbox has items waiting, or when Zeus asks for inbox triage.

---

## Inbox Location

`/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/99 Inbox (Queue)/`

---

## Triage Process

For each item in the inbox:

### Step 1 — Read and classify
- What pillar does this belong to? SO / HP / LE / GE
- What domain number? (00-24, 25-49, 50-74, 75-99)
- Is it a task, a note, a reference, or a project?

### Step 2 — Route
Move (with HITL approval) to the correct domain folder:

| Content | Route to |
|---------|---------|
| AI, agents, tech, automation, security | `00-24 SO/` |
| Personal growth, fitness, habits, money | `25-49 HP/` |
| Relationships, community, home | `50-74 LE/` |
| Career, ventures, brand, wealth | `75-99 GE/` |
| Unclear / needs Zeus decision | Leave in inbox, flag to Zeus |

### Step 3 — Format on arrival
- Add correct frontmatter (domain, agent, status, updated)
- Add to the correct domain's `## ⚡ Active Tasks` or `## 📎 Context & Notes`
- If it's a task, create in Todoist and add `tid:` reference

### Step 4 — Clear
Delete the original inbox item after successful routing (HITL approval required).

---

## Rules
- Never route blindly — if the destination is unclear, flag to Zeus
- Never delete from inbox without routing first
- All writes require HITL approval (see onemind-maintainer skill)
- Keep inbox at zero — it is a queue, not storage
