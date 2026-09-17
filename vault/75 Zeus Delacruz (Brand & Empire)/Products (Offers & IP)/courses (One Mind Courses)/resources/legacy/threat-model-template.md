# Threat Model Template

> *Use this to define the threat model for any system: your agent stack, a specific server, your credentials, or your full digital infrastructure.*

---

```yaml
---
type: threat-model
domain: 04
system: "[System Name]"
owner: Zeus
agent: Guardian
status: active
reviewed: YYYY-MM-DD
next-review: YYYY-MM-DD
---
```

---

## System: [System Name]

**Description:** [What is this system? What does it do? What data does it hold?]

**Criticality:** [Critical / High / Medium / Low]
— Critical = if compromised, major harm to Zeus, his family, or the empire
— High = significant disruption or data loss
— Medium = inconvenient but recoverable
— Low = minimal impact

---

## Assets to Protect

| Asset | Sensitivity | Location |
|-------|-------------|----------|
| [e.g., API keys] | Critical | 1Password |
| [e.g., vault data] | High | GitHub + VPS |
| [e.g., agent configs] | Medium | VPS ~/.openclaw/ |

---

## Threat Actors

Who might attack this system, and why?

| Actor | Motivation | Capability | Likelihood |
|-------|-----------|-----------|-----------|
| Script kiddie | Opportunistic | Low | High |
| Credential stuffing bot | Monetization | Medium | High |
| Targeted attacker | Financial/reputational | High | Low |
| Malicious insider | N/A (solo operation) | N/A | N/A |

---

## Attack Surface

List every way an attacker could reach this system:

- **Network:** [Open ports, exposed services, public endpoints]
- **Credentials:** [What credentials grant access? Where are they stored?]
- **Supply chain:** [Third-party dependencies, npm packages, Docker images]
- **Physical:** [Who has physical access to the hardware?]
- **Social engineering:** [Phishing, SIM swap, account recovery attacks]

---

## Threat Scenarios

For each realistic attack scenario, assess impact and likelihood:

### Scenario 1: [Scenario Name]

**Attack path:** [How would attacker execute this?]
**Impact if successful:** [What breaks or gets stolen?]
**Likelihood:** [1=rare, 5=likely]
**Impact severity:** [1=minimal, 5=catastrophic]
**Current mitigations:** [What's already protecting against this?]
**Gaps:** [What's not protected yet?]

---

### Scenario 2: [Scenario Name]

**Attack path:** [How?]
**Impact:** [What breaks?]
**Likelihood:** [1-5]
**Impact severity:** [1-5]
**Current mitigations:** [What exists?]
**Gaps:** [What's missing?]

---

### Scenario 3: [Scenario Name]

**Attack path:** [How?]
**Impact:** [What breaks?]
**Likelihood:** [1-5]
**Impact severity:** [1-5]
**Current mitigations:** [What exists?]
**Gaps:** [What's missing?]

---

## Risk Matrix

| Scenario | Likelihood | Impact | Risk Score | Priority |
|----------|-----------|--------|------------|---------|
| [Scenario 1] | 3 | 5 | 15 | P1 |
| [Scenario 2] | 4 | 2 | 8 | P2 |
| [Scenario 3] | 2 | 3 | 6 | P3 |

Risk Score = Likelihood × Impact. Prioritize highest scores first.

---

## Mitigations — Action Plan

| Risk | Mitigation | Owner | Due | Status |
|------|-----------|-------|-----|--------|
| [Scenario 1] | [What to implement] | Guardian/Zeus | YYYY-MM-DD | Pending |
| [Scenario 2] | [What to implement] | Grid | YYYY-MM-DD | Pending |

---

## Out of Scope

[Explicitly list what this threat model does NOT cover — set boundaries to keep the analysis focused]

---

## Review Schedule

Review this threat model:
- After any major system change
- Quarterly at minimum
- Immediately after any security incident

**Last reviewed:** YYYY-MM-DD
**Next review:** YYYY-MM-DD

---

*Agent Guardian: Review this threat model on the schedule above. Flag any new threat vectors that emerge from system changes.*
