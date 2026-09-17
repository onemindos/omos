---
name: onemind-architect
description: "Community-facing architecture skill — helps One Mind community members design their own vault, set up their pillars, configure agents, and customize domains for their life. Also serves as the internal architecture auditor for Zeus's own system. Use when helping someone set up One Mind, designing vault architecture, answering framework questions, or running system architecture audits. Activates on: architecture, setup, vault design, new member, onboarding, customize domains, pillar setup, system design, community setup, template."
id: skill-onemind-architect
title: "Skill: One Mind Architect"
type: skill
status: active
created: 2026-05-04
modified: 2026-06-06
owner: legacy
tags: [type/skill, domain/onemind, domain/so, area/architecture, area/community]
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# One Mind Architect — Community & System Design Skill

> Help anyone build their One Mind — from first vault to full system.

Two modes: **Community Architect** (help others design their One Mind) and **System Auditor** (evaluate Zeus's own architecture).

---

## Mode 1: Community Architect (Primary)

Guide community members through setting up their own One Mind vault.

### Onboarding Flow

1. **Assess** — What does this person need? Starter (3+1 pillars, basic domains) or Advanced (custom domains, agents)?
2. **Pillars** — Help them define their 3 life pillars + 1 operating layer. Default: Personal, Home, Business + Operations.
3. **Domains** — Start with 5-10 domains total (not 100). Help them pick the ones that match their life right now.
4. **Templates** — Point them to `onemind-codex-template` for the starter vault.
5. **First Captures** — Help them capture 3-5 real items from their life into the right domains.
6. **Review Cadence** — Set up their first weekly review rhythm.

### Design Principles for Community Vaults
- Start simple — 3+1 pillars, 5-10 domains max
- Domains should reflect THEIR life, not Zeus's
- No agents needed at first — the vault structure is enough
- CODE× method is universal — teach the loop first
- They can always add depth later

### Common Questions
| Question | Answer |
|----------|--------|
| "How many domains do I need?" | Start with 5-10. Add more only when you feel a gap. |
| "Do I need agents?" | No. The vault works without agents. Agents are Layer 5. |
| "Can I rename pillars?" | Yes. The names are yours. The 3+1 structure matters, not the labels. |
| "What goes in Operations?" | Whatever makes your other 3 pillars work — tools, tech, processes. |

---

## Mode 2: System Auditor (Internal)

> Analyze, evaluate, and evolve the One Mind system architecture across vault, agents, infrastructure, and domain structure.

---

## What the System Auditor Does

The Architect operates at the system level. While other skills handle execution within the system, the Architect examines the system itself — how pieces fit together, where friction lives, what's missing, and what needs to evolve. It applies to all 4 layers of One Mind:

| Layer | What Gets Architected |
|-------|-----------------------|
| **Vault** | Domain structure, folder taxonomy, routing logic, naming conventions |
| **Agents** | Agent roles, skill assignments, capability gaps, coordination patterns |
| **Infrastructure** | VPS topology, services, observability, storage, networking |
| **Framework** | CODE× method, pillar definitions, entity model, doctrine |

---

## When to Activate

| Trigger | What to Do |
|---------|-----------|
| Significant new feature or system added | Post-implementation review |
| System feels slow, cluttered, or fragmented | Bottleneck analysis |
| Major development phase completed | Cohesion check across all components |
| Planning structural changes | Pre-change architectural design |
| Quarterly review cycle | Full system health audit |
| Agent asks "does this fit the architecture?" | Design alignment check |

---

## Protocol 1: Architectural Review

Use after any significant change to vault structure, agent setup, or infrastructure.

```
STEP 1 — SCOPE
  Identify: what changed? What domain/layer does it touch?
  Load: relevant _README.md files, agent configs, infra docs

STEP 2 — ALIGNMENT CHECK
  Ask:
  - Does this follow the naming convention? (Name (Subtitle)/)
  - Does it belong in the right domain group (SO/HP/LE/GE)?
  - Does it duplicate anything that already exists?
  - Does the depth feel right (2-3 levels, no more)?
  - If an agent: does it have a clear domain scope and no overlap with existing agents?
  - If infrastructure: does it fit the existing service topology?

STEP 3 — ISSUE IDENTIFICATION
  Flag any of:
  - Naming violations (missing subtitle, wrong prefix)
  - Domain misrouting (content in wrong pillar)
  - Duplication (same concept in two places)
  - Depth problems (too flat = hard to find, too deep = too granular)
  - Orphaned skills/tools with no home
  - Agents with overlapping mandates

STEP 4 — RECOMMENDATION
  For each issue:
  - Describe the problem clearly
  - Propose the fix
  - Estimate impact (high/medium/low)
  - Note if Zeus approval is required before acting

STEP 5 — EXECUTE (if approved)
  Apply fixes. Update _README.md. Log changes.
```

---

## Protocol 2: Bottleneck Analysis

Use when something feels slow, cluttered, or inefficient.

```
SYMPTOMS → DIAGNOSIS → FIX

Symptom: Agents can't find the right skill
  Diagnosis: Skills folder flat or poorly categorized
  Fix: Apply pillar-aligned skill taxonomy (onemind/so/hp/le/ge/platforms)

Symptom: Inbox always full, nothing gets processed
  Diagnosis: No regular Organize cadence, or routing rules unclear
  Fix: Schedule weekly inbox triage; update code-method routing map

Symptom: Vault feels cluttered at domain root
  Diagnosis: Loose files not routed, subfolders missing
  Fix: Run full domain organization pass (see code-method Organize section)

Symptom: Two agents doing the same thing
  Diagnosis: Mandate overlap, unclear domain boundaries
  Fix: Redefine scope in IDENTITY.md files; route one agent's work to the other

Symptom: VPS services restarting or failing
  Diagnosis: Resource contention, missing health checks, no alerting
  Fix: Add to observability stack; set memory limits; document in omos-core.md

Symptom: Knowledge siloed — agents can't find relevant context
  Diagnosis: memorySearch.extraPaths not set, or vault not indexed
  Fix: Confirm CODEX.md in all workspaces; verify extraPaths config
```

---

## Protocol 3: Cohesion Check

Use after completing a major phase or set of changes. Ensures all new components integrate with the whole.

**Vault Cohesion:**
```bash
# Check for naming violations (folders without subtitles)
find "DOMAIN" -maxdepth 3 -type d | grep -v "(" | grep -v "\.git" | grep -v "_"

# Check for loose files at domain roots
find "DOMAIN" -maxdepth 1 -type f -name "*.md" ! -name "_README.md"

# Check for empty folders (might be stubs that need content)
find "DOMAIN" -type d -empty
```

**Agent Cohesion:**
- Every agent has IDENTITY.md, CODEX.md, USER.md in its workspace
- Every agent's domain scope is non-overlapping
- Every agent has the skills it needs in the correct pillar bucket
- No agent is missing a specialist for its primary domain

**Infrastructure Cohesion:**
- Every active service is documented in `omos-core.md` or the relevant node doc
- Every service has an entry in the observability stack
- No orphaned services running without documentation

---

## Protocol 4: Design Review (Pre-Change)

Before making structural changes, use this to evaluate the design.

**Questions to answer before any structural change:**

1. **Single Responsibility** — Does this new folder/agent/service do one thing clearly?
2. **Right Home** — Is this in the correct pillar/domain/layer?
3. **Name Clarity** — Does the name + subtitle tell you exactly what's inside without opening it?
4. **Depth** — Are we going 2-3 levels deep? Not more, not less?
5. **Duplication** — Does this duplicate anything that already exists?
6. **Future Proof** — Will this still make sense in 6 months when there are 2x as many items?
7. **Agent Findable** — Can an agent navigate to this without a human explaining the structure?

If the answer to any of these is "no" or "unsure" — redesign before building.

---

## Protocol 5: Quarterly Architecture Audit

Full system health check. Run once per quarter.

**Vault Audit:**
- [ ] Run loose file check across all 4 domain groups
- [ ] Check all domain roots have `_README.md`
- [ ] Check no folder violates naming convention
- [ ] Identify any domains that have grown too large (>50 subfolders = split candidate)
- [ ] Check inbox age — no items older than 2 weeks unprocessed

**Agent Audit:**
- [ ] All 9 agents have current IDENTITY.md, CODEX.md, USER.md
- [ ] No two agents have mandate overlap
- [ ] Each agent has skills assigned in the correct pillar bucket
- [ ] Memory indexing progress — note % complete
- [ ] Check for skill gaps — domains with no dedicated skill

**Infrastructure Audit:**
- [ ] All active VPS services documented in node docs
- [ ] All services appear in observability stack
- [ ] Backup verified working
- [ ] Security hardening still current (fail2ban, SSH key-only, firewall rules)
- [ ] Codex sync working (onemind-watch.service active, pull timer running)

**Framework Audit:**
- [ ] AGENTS.md and CLAUDE.md are current
- [ ] Domain registry (DOMAIN-REGISTRY.md) reflects actual structure
- [ ] No stale naming (check for old terms: "quadrant", "UI", "codex-" prefixes)
- [ ] ONEMIND-CODEX.md still accurate as canonical document

---

## One Mind Architecture Principles

These govern all design decisions. When in doubt, apply these:

**1. Pillar Alignment**
Everything belongs to one of the 4 pillars: SO (Operations), HP (Personal), LE (Home), GE (Business). If it doesn't fit cleanly, it goes to SO.

**2. Name Clarity**
Every folder, file, and agent name must communicate its purpose without explanation. If you have to open it to understand it, the name is wrong.

**3. Single Location**
One concept, one home. No duplicates. No cross-posting. Use links/references if something is relevant in multiple domains.

**4. Agent-Navigable**
The structure must be navigable by an AI agent without a human tour guide. Depth ≤ 3 levels. Consistent naming. README at every domain root.

**5. Separation of Concerns**
Skills ≠ Tools. Knowledge ≠ Code. Vault ≠ Infrastructure. Keep these layers clean.

**6. Evolutionary Design**
The system grows. Design for change. Placeholder folders are good. Don't over-engineer upfront — build for the next 6 months, not the next 6 years.

**7. Zeus Approval on Structural Changes**
Cross-pillar moves, domain renames, and agent mandate changes always get Zeus sign-off. Intra-domain organization can be autonomous.

---

## Current Architecture State (Reference)

### Vault Structure
```
SO (00-24) ← Operations — the platform
HP (25-49) ← Personal
LE (50-74) ← Home
GE (75-99) ← Business
99 Inbox   ← All captures land here first
```

### Agent Network (9 agents on omos-core)
| Agent | Domain | Model |
|-------|--------|-------|
| Legacy ⚡ | All — orchestrator | claude-sonnet-4.6 |
| Spartan 💪 | HP 27-28 — body, fitness | claude-sonnet-4.6 |
| Oracle 🔮 | HP 25-26, 29-34 — mind, money | claude-sonnet-4.6 |
| Eden 🌱 | LE — sustainability | claude-sonnet-4.6 |
| Haven 🏡 | LE 50-74 — family, home | claude-sonnet-4.6 |
| Empire 💰 | GE 75-99 — ventures, brand, wealth | claude-sonnet-4.6 |
| Guardian 🛡️ | SO 04 — security | claude-sonnet-4.6 |
| Grid ⚙️ | SO 02, 07 — infrastructure, automation | claude-sonnet-4.6 |
| Trinity 🔺 | SO 01 — agent coordination | claude-sonnet-4.6 |

### Skill Architecture
```
Skills (Agent Skills)/
  onemind/   ← vault + system skills (11)
  so/        ← operations skills (5)
  hp/        ← personal skills (9)
  le/        ← home skills (4)
  ge/        ← business skills (7)

Tools (Agent Tools)/
  platforms/ ← service connectors (12)
  _meta/     ← templates + specialist (2)
```

### Infrastructure (omos-core)
- VPS: AWS, Tailscale `100.98.25.10`, public `64.227.8.23`
- Gateway: OpenClaw v2026.4.14, port 18789
- Sync: onemind-watch.service (write → GitHub), onemind-pull.timer (GitHub → VPS)
- Observability: Datadog + OTEL (in progress)


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
