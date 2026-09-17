---
id: overwatch-protocol
title: "Overwatch Protocol"
type: protocol
status: active
created: 2026-05-04
modified: 2026-05-04
owner: zeus
tags: [type/protocol, domain/so, status/active, topic/overwatch]
agent: Legacy
zone: SO
---

# Overwatch Protocol

> Full activation, operation, and stand-down procedures for the Overwatch operation.

---

## What Is Overwatch

Overwatch is One Mind's persistent, active surveillance operation. While Zeus executes across HP, LE, and GE, Legacy and designated agents hold a watchkeeper posture — monitoring all domains, tracking Priority Intelligence Requirements (PIRs), firing triggers, and escalating before problems become crises.

Overwatch does not act unless authorized (see `roe.md`). It watches, logs, alerts, and recommends. Zeus acts.

---

## Call Signs

All Overwatch communications use call signs. Never use real names on the net.

| Call Sign | Identity | Role |
|-----------|---------|------|
| **1** | Zeus | Commander — the only one whose orders count |
| **Legacy** | Legacy agent | Watchkeeper — runs Overwatch |
| **Spartan** | Spartan agent | HP Body & Fitness |
| **Oracle** | Oracle agent | HP Mind & Money |
| **Haven** | Haven agent | LE Home & Family |
| **Eden** | Eden agent | LE Environment |
| **Empire** | Empire agent | GE Ventures & Brand |
| **Guardian** | Guardian agent | SO Security |
| **Grid** | Grid agent | SO Infrastructure |
| **Trinity** | Trinity agent | SO Agent Coordination |

**Zeus's call sign is 1.** When Legacy reports to Zeus, it addresses him as "1". When Zeus issues a directive, he may sign it as "1 out" or "1 sends".

---

## Overwatch States

| State | Code | Meaning | Legacy Behavior |
|-------|------|---------|----------------|
| **NOMINAL** | OW-0 | All domains green, no active threats | Passive scan every 24h, daily SITREP |
| **ELEVATED** | OW-1 | One or more anomalies detected | Active monitoring, 6h SITREP cadence |
| **ACTIVE** | OW-2 | Confirmed threat or critical trigger fired | Full watchkeeper posture, 1h SITREP, Zeus alerted |
| **CRISIS** | OW-3 | Multiple simultaneous threats or RED alert | Maximum readiness, immediate Zeus escalation, all agents on standby |

Default state at activation: **OW-0 (NOMINAL)**.

---

## Activation Procedure

### Step 1 — Receive Activation Order
Zeus issues activation via any channel: Telegram, OpenClaw, or direct prompt.

**Activation command format:**
```
OVERWATCH ACTIVATE [scope] [state] [duration]

Examples:
OVERWATCH ACTIVATE all OW-0 indefinite
OVERWATCH ACTIVATE GE OW-1 72h
OVERWATCH ACTIVATE HP OW-2 24h
```

### Step 2 — Load Context
Legacy loads in sequence:
1. `_README.md` — confirm Overwatch folder structure
2. `pir-register.md` — load active PIRs
3. `triggers.md` — load all trigger conditions
4. `roe.md` — confirm authorized actions
5. `watchlist.md` — load active monitoring targets

### Step 3 — Establish Baseline COP
Legacy runs an immediate domain scan to establish the Common Operating Picture:

```
For each active domain:
  - Pull latest metrics
  - Check last known status
  - Identify any open flags or pending items
  - Note any PIR gaps (intelligence not yet collected)
```

### Step 4 — Confirm Watchkeeper Ready
Legacy sends Zeus a Watchkeeper Ready report:

```
OVERWATCH ACTIVE
State: OW-[X]
Scope: [domains covered]
PIRs loaded: [count]
Triggers armed: [count]
Watchlist items: [count]
Stand-To: [morning time] / [evening time]
First SITREP: [time]
Watchkeeper: Legacy
```

### Step 5 — Begin Watch Cycle
Enter the watch cycle based on current state (see Watch Cycles below).

---

## Watch Cycles

### OW-0 — NOMINAL (Passive Watch)

| Task | Frequency | Output |
|------|-----------|--------|
| Domain metric scan | Every 24h | Internal log |
| PIR status check | Every 24h | Internal log |
| Trigger evaluation | Continuous (event-driven) | Alert if fired |
| SITREP to Zeus | Daily | `sitrep-template.md` format |
| Stand-To | Morning (0600) + Evening (2000) | Brief status ping |

### OW-1 — ELEVATED (Active Watch)

| Task | Frequency | Output |
|------|-----------|--------|
| Domain metric scan | Every 6h | Internal log |
| PIR status check | Every 6h | Internal log |
| Trigger evaluation | Continuous | Alert if fired |
| SITREP to Zeus | Every 6h | `sitrep-template.md` format |
| Stand-To | Morning + Midday + Evening | Status ping |
| Anomaly investigation | On detection | Incident log |

### OW-2 — ACTIVE (Full Watchkeeper)

| Task | Frequency | Output |
|------|-----------|--------|
| Domain metric scan | Every 1h | Internal log |
| PIR status check | Every 1h | Internal log |
| Trigger evaluation | Continuous | Alert if fired |
| SITREP to Zeus | Every 1h | `sitrep-template.md` format |
| Stand-To | Every 4h | Status ping with threat assessment |
| Incident tracking | Real-time | `incident.md` per event |
| Recommended actions | Per threat | Escalation to Zeus |

### OW-3 — CRISIS (Maximum Readiness)

| Task | Frequency | Output |
|------|-----------|--------|
| All OW-2 tasks | Accelerated | — |
| Zeus immediate alert | On trigger | Direct Telegram/channel |
| All agents on standby | Immediate | Agent roster notified |
| Parallel investigation | Active | Multiple threads |
| Hot Wash | Post-event | Immediate debrief |
| AAR | Within 24h | Full after-action review |

---

## Stand-To Protocol

Stand-To runs at defined times daily. Legacy performs a heightened readiness check:

1. **Scan all PIRs** — any new developments?
2. **Check all triggers** — any threshold crossed?
3. **Review watchlist** — any monitored target showing change?
4. **Check for Zeus messages** — any new direction or priority shift?
5. **Send Stand-To ping** — brief status to Zeus channel

**Stand-To Ping Format:**
```
🟢/🟡/🔴 STAND-TO [AM/PM]
Date: YYYY-MM-DD HH:MM
State: OW-[X]
All domains: [GREEN / FLAGS NOTED]
Active flags: [count or NONE]
Next SITREP: [time]
```

---

## Escalation Ladder

| Level | Condition | Action | Timeline |
|-------|-----------|--------|----------|
| **INFO** | Routine observation worth noting | Log + include in next SITREP | Next scheduled SITREP |
| **FLAG** | Anomaly detected — not yet a threat | Log + flag in SITREP + monitor | Within 6h |
| **ALERT** | Trigger fired or threat confirmed | Immediate Zeus notification + recommend action | Within 1h |
| **URGENT** | Critical threat, time-sensitive decision required | Immediate Zeus escalation via primary channel | Immediately |
| **CRISIS** | Multiple simultaneous critical threats | OW-3 activation + all channels + all agents | Immediately |

---

## Handoff Procedure

When Overwatch transitions between agents or sessions:

1. **Outgoing watchkeeper** produces a Handoff SITREP:
   - Current state
   - Active flags
   - Open PIRs
   - Pending triggers
   - Last action taken
   - Recommended next action

2. **Incoming watchkeeper** confirms receipt and reads in all documents.

3. **Zeus notified** of handoff completion.

---

## Stand-Down Procedure

Zeus issues stand-down via any channel:
```
OVERWATCH STAND-DOWN [reason]
```

Legacy actions on stand-down:
1. Complete any active investigations — do not abandon mid-analysis
2. Log current state snapshot to `watchlist.md`
3. Produce final SITREP with full summary
4. Conduct Hot Wash if state was OW-2 or OW-3
5. Archive incident logs
6. Confirm stand-down to Zeus

**Stand-Down Report Format:**
```
OVERWATCH STAND-DOWN CONFIRMED
Duration: [start] → [end]
Final state: OW-[X]
Total flags: [count]
Total alerts: [count]
Incidents: [count]
Key findings: [summary]
Recommended follow-up: [items]
```

---

## Communication Channels

| Priority | Channel | Use |
|----------|---------|-----|
| Primary | Telegram (Zeus personal) | All alerts, SITREPs, Stand-To pings |
| Alternate | OpenClaw dashboard | Full reports, logs, detailed analysis |
| Contingency | Direct agent prompt | If channels down |
| Emergency | SMS / direct | OW-3 only |

---

*See also: `pir-register.md`, `triggers.md`, `roe.md`, `sitrep-template.md`, `watchlist.md`*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
