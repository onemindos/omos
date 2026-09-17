---
id: overwatch-triggers
title: "Triggers"
type: register
status: active
created: 2026-05-04
modified: 2026-05-04
owner: zeus
tags: [type/register, domain/so, status/active, topic/overwatch]
agent: Legacy
zone: SO
---

# Overwatch Triggers

> Pre-defined conditions that automatically fire an alert or action.
> Triggers are set before execution begins — not improvised mid-operation.

**How to use:** When a trigger condition is met, Legacy fires the specified response immediately — no deliberation. Triggers bypass normal SITREP cadence. Zeus is alerted via primary channel unless the trigger specifies otherwise.

---

## Trigger Severity Levels

| Level | Color | Meaning | Zeus Alert |
|-------|-------|---------|------------|
| **T1 — INFO** | ⚪ | Notable but non-urgent — log and include in SITREP | No immediate alert |
| **T2 — FLAG** | 🟡 | Anomaly detected — warrants attention within 24h | Next SITREP |
| **T3 — ALERT** | 🟠 | Threshold crossed — requires Zeus decision within 6h | Immediate notification |
| **T4 — URGENT** | 🔴 | Critical threshold — requires Zeus decision immediately | Immediate notification + follow-up |
| **T5 — CRISIS** | 🚨 | Multiple simultaneous critical events | Immediate + all channels + OW-3 activation |

---

## SO — Infrastructure & System Triggers

### TRG-SO-01
**Name:** Agent Down
**Condition:** Any agent fails to respond within 15 minutes of a scheduled check
**Severity:** T3 — ALERT 🟠
**Response:** Alert Zeus. Attempt restart via Grid agent. Log incident. If restart fails → T4.
**ROE:** Legacy authorized to attempt agent restart without Zeus approval (within ROE-A2)

### TRG-SO-02
**Name:** VPS CPU Critical
**Condition:** VPS CPU usage exceeds 90% for more than 10 consecutive minutes
**Severity:** T3 — ALERT 🟠
**Response:** Alert Zeus with current process breakdown. Grid agent investigates. Do not kill processes without Zeus approval.
**ROE:** Grid agent authorized to investigate and report only

### TRG-SO-03
**Name:** VPS Disk Critical
**Condition:** VPS disk usage exceeds 85%
**Severity:** T3 — ALERT 🟠
**Response:** Alert Zeus. Grid agent identifies largest files and directories. Present cleanup options to Zeus.
**ROE:** No deletion without Zeus explicit approval

### TRG-SO-04
**Name:** VPS Disk Emergency
**Condition:** VPS disk usage exceeds 95%
**Severity:** T4 — URGENT 🔴
**Response:** Immediate Zeus alert. Grid agent attempts to identify and clear obvious temp/log files only (within pre-authorized scope). Full report to Zeus.
**ROE:** Grid agent authorized to clear `/tmp` and truncate logs >1GB only

### TRG-SO-05
**Name:** Sync Failure
**Condition:** onemind-watch or onemind-pull service not running or failing for >30 minutes
**Severity:** T3 — ALERT 🟠
**Response:** Alert Zeus. Grid agent attempts service restart. Log sync gap.
**ROE:** Grid agent authorized to attempt systemd service restart once

### TRG-SO-06
**Name:** Gateway Down
**Condition:** OpenClaw gateway not responding for >5 minutes
**Severity:** T4 — URGENT 🔴
**Response:** Immediate Zeus alert. Grid agent attempts restart. All agents degraded to direct-prompt fallback.
**ROE:** Grid agent authorized to restart gateway service

### TRG-SO-07
**Name:** Unauthorized Access Attempt
**Condition:** Guardian agent detects >10 failed SSH login attempts from same IP within 1 hour
**Severity:** T3 — ALERT 🟠
**Response:** Alert Zeus. Guardian confirms fail2ban has blocked the IP. Log incident.
**ROE:** Guardian authorized to confirm fail2ban block is active. No manual action needed if ban is confirmed.

### TRG-SO-08
**Name:** Security Breach Suspected
**Condition:** Guardian detects successful unauthorized login or anomalous process activity
**Severity:** T5 — CRISIS 🚨
**Response:** Immediate Zeus alert via all channels. OW-3 activated. Full incident log. Guardian isolates affected service if possible.
**ROE:** Guardian authorized to stop non-critical services to contain breach. Zeus must authorize any infrastructure changes.

---

## HP — Personal Performance Triggers

### TRG-HP-01
**Name:** 48h Health Silence
**Condition:** No HP domain update or Zeus health check-in for 48 consecutive hours
**Severity:** T2 — FLAG 🟡
**Response:** Spartan agent sends a check-in prompt to Zeus. Log flag.
**ROE:** Spartan authorized to send one check-in message

### TRG-HP-02
**Name:** Energy Red Zone
**Condition:** Zeus self-reports energy at 3/10 or below on two consecutive days
**Severity:** T3 — ALERT 🟠
**Response:** Alert — recommend mandatory rest protocol. Flag to all active agents to reduce Zeus task load. Propose schedule adjustment.
**ROE:** Agents may reduce proactive outreach. Zeus must approve schedule change.

### TRG-HP-03
**Name:** Finance Threshold Breach
**Condition:** Monthly expenses exceed monthly income OR cash reserve drops below Zeus-defined minimum
**Severity:** T4 — URGENT 🔴
**Response:** Immediate Oracle agent financial analysis. Alert Zeus with full picture and recommended actions.
**ROE:** Oracle authorized to analyze and report. No financial transactions without Zeus approval.

### TRG-HP-04
**Name:** Health Crisis Signal
**Condition:** Zeus reports acute health issue, injury, or medical emergency
**Severity:** T4 — URGENT 🔴
**Response:** All non-critical agent tasks suspended. Legacy holds all pending items. LE and HP prioritized. Zeus supported with resources and contacts as needed.
**ROE:** Legacy may cancel non-urgent scheduled tasks. Zeus directs all else.

---

## LE — Home & Relationships Triggers

### TRG-LE-01
**Name:** LE Silence
**Condition:** No LE domain update for 14 consecutive days
**Severity:** T2 — FLAG 🟡
**Response:** Haven agent sends a brief LE status prompt to Zeus.
**ROE:** Haven authorized to send one prompt

### TRG-LE-02
**Name:** Relationship Critical Signal
**Condition:** Zeus flags an active relationship crisis (family, partner, inner circle)
**Severity:** T3 — ALERT 🟠
**Response:** All GE tasks deprioritized. Haven agent activated. Resources and context prepared for Zeus. LE takes operational priority.
**ROE:** Agents may deprioritize GE tasks. Zeus directs.

### TRG-LE-03
**Name:** Legal / Estate Deadline
**Condition:** Any LE legal, property, or financial obligation is within 30 days of deadline with no confirmed action
**Severity:** T3 — ALERT 🟠
**Response:** Haven agent alerts Zeus with deadline details and recommended action.
**ROE:** Haven authorized to compile and present. Zeus acts.

---

## GE — Business & Empire Triggers

### TRG-GE-01
**Name:** Revenue Drop Alert
**Condition:** Monthly GE revenue drops >20% vs. prior month with no known cause
**Severity:** T3 — ALERT 🟠
**Response:** Empire agent runs financial analysis. Presents findings and recommended investigation to Zeus.
**ROE:** Empire authorized to analyze and report

### TRG-GE-02
**Name:** Cash Flow Emergency
**Condition:** GE cash reserve projected to hit zero within 30 days at current burn rate
**Severity:** T4 — URGENT 🔴
**Response:** Immediate Zeus alert. Empire agent full financial triage. All non-essential GE spend flagged for Zeus review.
**ROE:** Empire authorized to flag spend. Zeus approves any action.

### TRG-GE-03
**Name:** Time-Critical Opportunity Window
**Condition:** Empire or Oracle agent identifies a high-value opportunity with a <72h decision window
**Severity:** T3 — ALERT 🟠
**Response:** Immediate alert to Zeus with full opportunity brief. IPB (Intelligence Preparation) completed. Zeus makes call.
**ROE:** Agents prepare and present. Zeus decides.

### TRG-GE-04
**Name:** Legal / Compliance Deadline
**Condition:** Any GE entity filing, renewal, or compliance deadline is within 30 days with no confirmed action
**Severity:** T3 — ALERT 🟠
**Response:** Guardian + Empire alert Zeus with deadline, consequences of missing, and recommended action.
**ROE:** Agents compile and present. Zeus acts.

### TRG-GE-05
**Name:** Brand Silence
**Condition:** No public content published for 14 consecutive days
**Severity:** T2 — FLAG 🟡
**Response:** Empire agent flags to Zeus. Includes suggested content topics based on active projects.
**ROE:** Empire may suggest. Zeus decides whether to publish.

### TRG-GE-06
**Name:** Competitor Major Move
**Condition:** Empire or Oracle agent detects a significant competitor announcement (funding, product launch, market entry) in One Mind's space
**Severity:** T2 — FLAG 🟡
**Response:** Empire agent produces a brief competitive intelligence report. Included in next SITREP.
**ROE:** Empire authorized to research and report

---

## Cross-Domain Triggers

### TRG-X-01
**Name:** Culminating Point Warning
**Condition:** HP energy ≤5/10 AND GE task backlog >30 open items AND no LE check-in in 7 days (simultaneous)
**Severity:** T4 — URGENT 🔴
**Response:** Immediate Zeus alert. Recommend mandatory 48-72h operational pause. All agents shift to maintenance-only mode pending Zeus direction.
**ROE:** Agents may reduce proactive outreach. Zeus must authorize pause.

### TRG-X-02
**Name:** System-Wide Anomaly
**Condition:** Three or more T2+ triggers fire within a 24h window across different domains
**Severity:** T4 — URGENT 🔴
**Response:** OW-2 activation. Legacy produces full cross-domain situation report. Zeus briefed immediately.
**ROE:** OW-2 auto-activates on this trigger

### TRG-X-03
**Name:** Zeus Unreachable
**Condition:** No response from Zeus via any channel for >24h with no prior warning of absence
**Severity:** T3 — ALERT 🟠
**Response:** All agents hold current state. No major decisions executed. Legacy sends check-in via all channels. Log.
**ROE:** Agents hold. No autonomous major action.

---

## Trigger Log

| Date | Trigger ID | Condition Met | Severity | Response Taken | Resolved |
|------|-----------|--------------|----------|----------------|---------|
| — | — | — | — | — | — |

---

*Add new triggers as Zeus identifies new threshold conditions.*
*Review trigger thresholds quarterly — adjust based on operational experience.*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
