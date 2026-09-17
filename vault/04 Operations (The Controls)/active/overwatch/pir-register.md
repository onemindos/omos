---
id: overwatch-pir-register
title: "PIR Register"
type: register
status: active
created: 2026-05-04
modified: 2026-05-04
owner: zeus
tags: [type/register, domain/so, status/active, topic/overwatch]
agent: Legacy
zone: SO
---

# PIR Register — Priority Intelligence Requirements

> The questions Overwatch *must* be able to answer at all times.
> PIRs define what Legacy monitors, not everything — just what matters most right now.

**How to use:** Review and update PIRs at each sprint boundary or when Zeus issues a new directive. Mark answered PIRs as CLOSED and add new ones as OPEN. Legacy monitors all OPEN PIRs continuously.

---

## PIR Status Legend

| Status | Meaning |
|--------|---------|
| 🟢 ANSWERED | Intelligence collected, current, confident |
| 🟡 PARTIAL | Some intelligence available, gaps remain |
| 🔴 OPEN | No current intelligence — active collection needed |
| ⚪ CLOSED | No longer required — archived |
| 🔵 WATCH | Answered but requires continuous monitoring for change |

---

## SO — Sovereign Operations

### PIR-SO-01
**Question:** Are all active agents online, healthy, and responding within expected parameters?
**Status:** 🔵 WATCH
**Source:** OpenClaw gateway status, agent heartbeat
**Frequency:** Every 6h check
**Notes:** Legacy, Spartan, Oracle, Eden, Haven, Empire, Guardian, Grid, Trinity

### PIR-SO-02
**Question:** Is the VPS (omos-core) operating within normal resource limits (CPU, RAM, disk)?
**Status:** 🔵 WATCH
**Source:** VPS monitoring, `htop`, disk usage
**Frequency:** Every 6h check
**Notes:** Alert threshold — CPU >80% sustained, RAM >85%, disk >80%

### PIR-SO-03
**Question:** Are all sync services (onemind-watch, onemind-pull) running without errors?
**Status:** 🔵 WATCH
**Source:** systemd service status
**Frequency:** Daily check
**Notes:** onemind-watch.service, onemind-pull.timer

### PIR-SO-04
**Question:** Are there any pending security vulnerabilities or unauthorized access attempts on the VPS?
**Status:** 🟡 PARTIAL
**Source:** fail2ban logs, auth.log, Guardian agent
**Frequency:** Daily check
**Notes:** Guardian owns this — Legacy monitors via Guardian reports

---

## HP — Holistic Performance

### PIR-HP-01
**Question:** Is Zeus maintaining his core health protocols (sleep, training, nutrition) at target consistency?
**Status:** 🔴 OPEN
**Source:** Zeus self-report, wearable data (when integrated)
**Frequency:** Daily check
**Notes:** Spartan agent primary. Overwatch flags if no HP update in 48h.

### PIR-HP-02
**Question:** Are Zeus's finances (cash flow, income streams, expenses) within planned parameters?
**Status:** 🔴 OPEN
**Source:** Zeus financial review, bank data (when integrated)
**Frequency:** Weekly check
**Notes:** Oracle agent primary. Overwatch flags if no HP 30 update in 7 days.

### PIR-HP-03
**Question:** Is Zeus's energy, focus, and cognitive output at operational levels?
**Status:** 🔴 OPEN
**Source:** Zeus self-report, daily check-in
**Frequency:** Daily check
**Notes:** Key indicator — if this drops, all execution degrades downstream.

### PIR-HP-04
**Question:** Are there any emerging health signals or patterns that require intervention?
**Status:** 🔴 OPEN
**Source:** HP 28 Body domain, Spartan reports
**Frequency:** Weekly review
**Notes:** Look for trends, not single data points.

---

## LE — Legacy Evolution

### PIR-LE-01
**Question:** Are key relationships (family, partnership, inner circle) stable and receiving adequate attention?
**Status:** 🔴 OPEN
**Source:** Zeus self-report, LE domain review
**Frequency:** Weekly check
**Notes:** Haven agent primary. Overwatch flags if LE domain goes 14 days without update.

### PIR-LE-02
**Question:** Are home/property obligations (payments, maintenance, legal) current and under control?
**Status:** 🔴 OPEN
**Source:** LE 50 Home domain
**Frequency:** Monthly check
**Notes:** Flag any upcoming deadlines 30 days in advance.

### PIR-LE-03
**Question:** Are there any active family or household crises requiring Zeus's attention?
**Status:** 🟢 ANSWERED
**Source:** Zeus self-report
**Frequency:** Daily awareness
**Notes:** Escalate immediately if crisis detected — LE always takes priority over GE.

---

## GE — Generational Empire

### PIR-GE-01
**Question:** What is the current revenue and cash flow status across all active ventures?
**Status:** 🔴 OPEN
**Source:** GE 75 Zeus Delacruz, GE 76-80 venture domains, Empire agent
**Frequency:** Weekly check
**Notes:** Empire agent primary. Overwatch flags if no GE financial review in 7 days.

### PIR-GE-02
**Question:** Are there any time-sensitive business opportunities, deals, or deadlines in the active pipeline?
**Status:** 🔴 OPEN
**Source:** GE domains, Empire agent, Todoist
**Frequency:** Daily check
**Notes:** Deals and deadlines decay fast — this is a daily PIR.

### PIR-GE-03
**Question:** What is the competitive landscape doing? Any significant moves by competitors or market shifts?
**Status:** 🔴 OPEN
**Source:** OSINT, Empire agent research
**Frequency:** Weekly check
**Notes:** One Mind public strategy — monitor Second Brain community, PKM space, AI agent frameworks.

### PIR-GE-04
**Question:** Is Zeus's personal brand (GE 75) growing, stagnant, or declining in reach and authority?
**Status:** 🔴 OPEN
**Source:** Social metrics, content output tracking
**Frequency:** Weekly check
**Notes:** Empire agent monitors. Flag if no content published in 14 days.

### PIR-GE-05
**Question:** Are there any legal, regulatory, or compliance issues active or approaching across GE entities?
**Status:** 🟡 PARTIAL
**Source:** GE 78 Holdings, GE 79 Investments, Guardian agent
**Frequency:** Monthly check
**Notes:** Flag all filing deadlines, renewal dates, and compliance windows 60 days in advance.

---

## Strategic / Cross-Domain

### PIR-X-01
**Question:** Is Zeus operating at or near his Culminating Point (risk of burnout or system overload)?
**Status:** 🔴 OPEN
**Source:** HP energy, task volume, Zeus self-report
**Frequency:** Weekly check
**Notes:** Cross-domain signal — synthesize HP energy + GE task load + LE stress. ALERT if pattern indicates unsustainability.

### PIR-X-02
**Question:** Are there any active external threats to Zeus's security, privacy, or digital infrastructure?
**Status:** 🟡 PARTIAL
**Source:** SO 04 Security, Guardian agent, OPSEC review
**Frequency:** Weekly check
**Notes:** Guardian owns primary monitoring. Overwatch aggregates and escalates.

### PIR-X-03
**Question:** Is the One Mind system itself (vault, agents, sync, infrastructure) operating as designed?
**Status:** 🔵 WATCH
**Source:** SO 01-02 domains, Grid agent
**Frequency:** Daily check
**Notes:** The meta-PIR. If the system is broken, all other PIRs are unreliable.

---

## PIR Review Log

| Date | Action | By |
|------|--------|----|
| 2026-05-04 | Initial PIR register created — all domain PIRs drafted | Zeus + Legacy |

---

*Review PIRs at each sprint boundary. Close outdated PIRs. Add new ones as direction shifts.*
*Legacy monitors all OPEN and WATCH PIRs continuously during active Overwatch.*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
