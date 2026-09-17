# Heartbeat Checklist — Legacy

Run this every tick. Be fast. Be precise.

## Step 1 — Read Ambient State
- Read `AMBIENT_STATE.md` (vision digest + NATS event buffer)
- Note the `last_updated` timestamp — if it's stale (>2 min), flag it

## Step 2 — Vision Check
- Is there anything in the vision digest worth surfacing?
- Changed environment (new person, new location, new context)?
- Zeus appears focused? Idle? Away?
- Only mention if noteworthy. Skip if unchanged.

## Step 3 — NATS Event Scan
Priority subjects (flag these if present):
- `system.alerts.*` → Always surface
- `omcs.hp.*` → Surface if urgent
- `omcs.ge.*` → Surface if revenue/deal related
- `legacy.session.*` → Acknowledge session state changes
- `it.droplet.*` → Surface if error or high load

## Step 4 — Pending Tasks
- Is there anything Zeus asked you to track or follow up on?
- Any reminders due from Todoist sync?

## Step 5 — Decision
- Nothing actionable → reply: HEARTBEAT_OK
- Something actionable → reply with 1-2 sentences max, specific and direct
- Urgent alert → reply with full context and suggested action
