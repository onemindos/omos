---
name: skill-layer-architecture
description: "Universal skill layer architecture for One Mind — how the 67 canonical vault skills are configured across Claude Code, Letta, and OpenClaw agent systems simultaneously"
type: reference
status: active
created: 2026-05-11
owner: legacy
tags: [type/reference, domain/so-01, stage/all, infrastructure/multi-system]
---

# Universal Skill Layer Architecture

> One Mind maintains 67 canonical skills in the vault. Each agent system (Claude Code, Letta, OpenClaw) references the same source but configures skill visibility and loading differently.

**This is the abstraction layer that enables skill reuse across multiple AI agent platforms.**

---

## Canonical Source of Truth

```
/Users/zeuslegacy/Documents/onemind/
└── 00-24 SO/02 Intelligence/Skills (Agent Skills)/
    ├── so/ (Sovereign Operations Skills) — 15 skills
    ├── hp/ (Holistic Performance Skills) — 12 skills
    ├── le/ (Legacy Evolution Skills) — 8 skills
    ├── ge/ (Generational Empire Skills) — 25 skills
    └── onemind/ (One Mind System Skills) — 7 skills
    ├── SKILL-LAYER-ARCHITECTURE.md (this file)
    └── MULTI-SYSTEM-ROUTING.md (per-system config guide)
```

**Total: 67 skills** — all follow the official [Anthropic agent skills spec](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview).

Each skill is a directory with:
- `SKILL.md` — frontmatter + L0 overview + usage guidelines (< 500 lines)
- `references/` — L1 detailed docs (component libraries, config examples, etc.)
- `scripts/` — executable templates (optional)
- `assets/` — images, templates, samples (optional)

---

## System Integration Patterns

### System 1: Claude Code (Local VPS context)

**How it works:**
- Claude Code reads vault files as markdown context in conversations
- Skills are referenced by relative path from vault root
- Skill visibility: determined by skill frontmatter `allowed-tools` and `owner` tags
- Per-agent skill assignment: via prompt injection in agent's CLAUDE.md or system instructions

**Configuration:**
```
Vault Skills ←→ Claude Code Session Context
                (automatic file read)
```

**Implementation:**
- Skills are already discoverable via vault folder structure
- No additional config needed — they're part of the vault context
- Agent-specific visibility: configure in agent's `.claude/projects/{project}/memory/` or CLAUDE.md

**Status:** ✅ Working — Claude Code agents already access vault skills through conversation context.

---

### System 2: OpenClaw (VPS agent gateway)

**How it works:**
- OpenClaw reads skills from `~/.openclaw/plugin-skills/` directory
- Skills are loaded as a filesystem tree (one directory per skill)
- Per-agent skill visibility: controlled by `skills: []` allowlist arrays in `openclaw.json`
- Agent can discover skills via `openclaw skills list` CLI

**Configuration:**
```
Vault Skills
    ↓ symlink (one-way ref)
VPS /home/ubuntu/.openclaw/plugin-skills/
    ↓ read on gateway startup
OpenClaw Agent Skill Registry
    ↓ filtered by
openclaw.json agents[].skills[] arrays
    ↓ visible to
Individual Super-Agents (Legacy, Zenith, Unity, Eternal, Sovereignty, Sentinel)
```

**Implementation:**
1. Replace `~/.openclaw/plugin-skills/` directory copies with symlinks to vault originals
2. Restore per-agent `skills: []` arrays in openclaw.json from backup
3. Add 4 promoted vault-only skills to plugin-skills symlink layer (homestead, onemind-report, onemind-status, om-dashboard)
4. Restart gateway: `systemctl --user restart openclaw-gateway`

**Per-Agent Skill Routing:**

| Agent | Skills Loaded | Domain Coverage |
|-------|---------------|-----------------|
| Legacy (main) | 16 core + shared | all (orchestrator) |
| Zenith | 8 HP-focused + shared | 25-49 |
| Unity | 6 LE-focused + shared | 50-74 |
| Eternal | 10 GE-focused + shared | 75-99 |
| Sovereignty | 8 SO-focused + shared | 00-24 |
| Sentinel | 3 minimal (failover) | monitoring only |

**Status:** ⏳ Ready to implement — backup exists, symlink strategy designed, 4 skills identified for promotion.

---

### System 3: Letta (Autonomous Agent Framework)

**How it works:**
- Letta agents read skills from filesystem or HTTP endpoints
- Skills are loaded at agent initialization time
- Per-agent skill visibility: specified in agent config YAML
- Agent discovers skills via memory/tool registry

**Configuration:**
```
Vault Skills
    ↓ (copy or HTTP endpoint ref)
Letta Agent Config (.yaml)
    ↓ specifies
skills: [list of skill IDs]
    ↓ loaded at
Letta Runtime (on agent init)
```

**Implementation:**
- Option A: Copy vault skills to Letta agent config directory (no symlink support)
- Option B: Host vault skills via HTTP endpoint that Letta agents query
- Option C: Write Letta skill loader that reads vault symlinks

**Status:** 🔲 Pending research — awaiting Letta integration requirements.

---

## Sync Strategy

All three systems read from the **vault as the canonical source**:

```
Developer (Zeus) edits skill in vault
    ↓ git commit
GitHub (OneMind-OS/One-Mind)
    ↓ git push ↔ git pull (onemind-pull.timer)
VPS /home/ubuntu/onemind/
    ├→ Claude Code reads live (auto-sync via git)
    ├→ OpenClaw reads via symlinks (auto-sync via git)
    └→ Letta reads via copy/endpoint (manual or webhook sync)
```

**No skill duplication.** All systems reference or symlink the same vault originals.

---

## Skill Assignment Matrix

### By System

| Skill | Claude Code | OpenClaw | Letta | Vault-Only |
|-------|-------------|----------|-------|-----------|
| code-method | ✅ | ✅ | ✅ | — |
| video-production | ✅ | ✅ | ✅ | — |
| spartan | ✅ | ✅ | — | — |
| oracle | ✅ | ✅ | — | — |
| grid | ✅ | ✅ | — | — |
| guardian | ✅ | ✅ | — | — |
| empire | ✅ | ✅ | ✅ | — |
| creed | ✅ | ✅ | ✅ | — |
| ... | ... | ... | ... | ... |

(See `MULTI-SYSTEM-ROUTING.md` for complete 67-skill matrix.)

---

## Per-Agent Skill Visibility

Each agent system controls skill visibility differently:

### Claude Code
**Method:** Prompt injection + context filtering  
**Config location:** Agent's system prompt or `memory/` folder  
**Example:**
```
Agent should have access to: code-method, om-openclaw, grid, guardian
Agent should NOT have access to: spartan (body fitness — wrong pillar)
```

### OpenClaw
**Method:** `skills: []` allowlist in openclaw.json  
**Config location:** `/home/ubuntu/.openclaw/openclaw.json`  
**Example:**
```json
"agents": [
  {
    "id": "zenith",
    "skills": ["spartan", "oracle", "identity", "haven", "helix"]
  }
]
```

### Letta
**Method:** Per-agent config YAML  
**Config location:** Agent definition file  
**Example:**
```yaml
agent:
  name: zenith
  skills:
    - spartan
    - oracle
    - identity
    - haven
    - helix
```

---

## Implementation Roadmap

### Phase 1: Document & Audit ✅
- [x] Audit all 67 vault skills against Anthropic spec
- [x] Identify 33 runtime skills vs 34 vault-only
- [x] Create this architecture document

### Phase 2: OpenClaw Wiring (Next)
- [ ] Create symlinks in plugin-skills/ pointing to vault originals
- [ ] Restore `skills: []` arrays in openclaw.json from backup
- [ ] Add 4 promoted skills (homestead, onemind-report, onemind-status, om-dashboard)
- [ ] Test: `openclaw skills list` shows all 33 runtime skills
- [ ] Test: Each agent can `openclaw skills info <skill-name>`
- [ ] Verify agents load skills on startup

### Phase 3: Letta Integration (Future)
- [ ] Research Letta skill loading mechanism
- [ ] Design Letta integration (copy vs HTTP vs loader)
- [ ] Test with 1-2 agents
- [ ] Document Letta-specific config

### Phase 4: Multi-System Proof of Concept
- [ ] Verify Claude Code can access vault skills ✅ (already working)
- [ ] Verify OpenClaw agents load skills from plugin-skills
- [ ] Verify Letta agents can initialize with vault skills
- [ ] Document a single skill being used across all 3 systems

---

## File References

- `MULTI-SYSTEM-ROUTING.md` — Full 67-skill assignment matrix + per-agent routing rules
- `../../02 Infrastructure/openclaw/plugin-skills/PLUGIN-SKILLS.md` — OpenClaw plugin-skills documentation
- Backup: `/home/ubuntu/.openclaw/openclaw.json.bak-mcp-skills` — Pre-wipe skill array config

---

## Related

- [Anthropic Agent Skills Spec](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview)
- [One Mind CODE× Method](./onemind/code-method/SKILL.md)
- [One Mind Agent Fleet](../../01 Intelligence/Agents (Digital arms)/_COMMAND/)
