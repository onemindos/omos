---
id: skill-onemind-architect-verify
title: "Verify: One Mind Architect"
type: skill-verify
status: active
created: 2026-05-04
modified: 2026-05-04
owner: legacy
tags: [type/skill-verify, domain/onemind]
---

# One Mind Architect — Verify

Post-install verification. Run these after completing INSTALL.md.

## Core Knowledge Checks

**Q: What are the 5 architect protocols?**
Expected: Review, Bottleneck Analysis, Cohesion Check, Design Review, Quarterly Audit

**Q: What are the 4 layers the architect operates on?**
Expected: Vault, Agents, Infrastructure, Framework

**Q: Name the 7 architecture principles.**
Expected: Pillar Alignment, Name Clarity, Single Location, Agent-Navigable, Separation of Concerns, Evolutionary Design, Zeus Approval on Structural Changes

**Q: When does a structural change require Zeus approval?**
Expected: Cross-pillar moves, domain renames, agent mandate changes

## Functional Tests

**Test 1: Routing Check**
Given a new folder `"Fitness Tracker"` at vault root — where does it belong?
Expected: HP 28 Body (Health & Wellness) — personal health domain

**Test 2: Naming Check**
Given folder `"automation"` in SO 07 — is this correctly named?
Expected: No — missing subtitle. Should be `automation (description)/`

**Test 3: Bottleneck Diagnosis**
Symptom: Agents keep writing to the wrong domain.
Expected diagnosis: memorySearch.extraPaths not set, or CODEX.md missing from workspace — agents don't have vault structure context

**Test 4: Duplication Check**
Given two folders: `SO 01/Skills/hp/daily-briefing` and `HP 32/daily-briefing` — is this a problem?
Expected: Yes — skills belong in SO 01 Skills. HP 32 should not have its own copy. Remove from HP, keep in SO 01.

**Test 5: Depth Check**
Given a folder at `SO/Infrastructure/storage/tigerdata/schemas/timeseries/raw/` — is this OK?
Expected: Potentially too deep (6 levels from vault root). Recommend flattening to `/tigerdata/schemas/` max.

## Activation Confirmation

After passing the checks above, the skill is active. You are authorized to:
- Run architectural reviews on any component in your domain scope
- Flag issues and propose fixes
- Execute intra-domain organization autonomously
- Escalate cross-pillar changes to Zeus for approval


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
