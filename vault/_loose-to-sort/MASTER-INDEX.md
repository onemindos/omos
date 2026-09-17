---
type: index
domain: SO-00
status: active
created: 2026-05-14
updated: 2026-05-14
owner: zeus
depth: 1
agent: Legacy
tags: [type/index, domain/so, status/active]
ai_summary: "Master vault index — flat table of every domain file with path, ai_summary, agent, and tags. Single-file navigation for agents. Maintained by Legacy on cron."
---

# MASTER-INDEX

> Single-file navigation layer. Agent scans this one file to locate anything in the vault — no folder traversal required.
> **Maintained by Legacy.** Regenerate: run the `onemind-maintainer` skill → `regenerate-index` command.
> **Last generated:** 2026-05-14

---

## How to Use This File

```
AGENT WORKFLOW:
1. Read MASTER-INDEX.md
2. Scan the table for matching domain / tags / ai_summary keywords
3. Navigate directly to the path — no folder walk needed
4. Fall back to list_directory + _README.md if path not found here
```

**Depth key:**
- `1` = Pillar root (SO / HP / LE / GE README)
- `2` = Domain index (e.g. `HP 28 Body/28 Body.md`)
- `3` = Subfolder (e.g. `HP 28 Body/Sleep/`)
- `4` = Deep content (specific notes, guides, reference docs)

---

## Pillar Roots (depth: 1)

| Path | Agent | Summary |
| ---- | ----- | ------- |
| `00-24 SO (Sovereign Operations)/README.md` | Sovereignty | SO pillar root — operating layer index |
| `25-49 HP (Holistic Performance)/README.md` | Zenith | HP pillar root — personal layer index |
| `50-74 LE (Legacy Evolution)/README.md` | Unity | LE pillar root — home layer index |
| `75-99 GE (Generational Empire)/README.md` | Eternal | GE pillar root — business layer index |

---

## SO Domains (depth: 2)

| Path | Agent | Summary |
| ---- | ----- | ------- |
| `00-24 SO/00 Framework (The Constitution)/00 Framework.md` | Sovereignty | Doctrine, philosophy, architecture — the constitution |
| `00-24 SO/01 Intelligence (The Brain)/01 Intelligence.md` | Sovereignty | AI, agents, models, skills, memory, tools |
| `00-24 SO/02 Infrastructure (The Body)/` | Sovereignty | VPS, servers, networking, storage, services |
| `00-24 SO/03 Security (The Shield)/` | Sovereignty | Cybersecurity, access control, hardening |
| `00-24 SO/04 Operations (The Controls)/` | Sovereignty | Runbooks, SOPs, comms, commands, hooks |
| `00-24 SO/05 Platforms (The Connectors)/05 Platforms.md` | Sovereignty | External tools, APIs, integrations catalog |
| `00-24 SO/06 OpenClaw (The Runtime)/` | Sovereignty | Agent gateway, workspaces, CLI, plugins |
| `00-24 SO/07 Architecture (The Blueprint)/` | Sovereignty | System design, 9 architectural layers |
| `00-24 SO/08 OM Assets (Entity Registry)/` | Sovereignty | System registry, nodes, templates |
| `00-24 SO/09 Development (The Forge)/` | Sovereignty | DevOps, deploy, scripts, guides |
| `00-24 SO/10 Research (The Lab)/` | Sovereignty | R&D, robotics, AI agents, emerging tech |

---

## HP Domains (depth: 2)

| Path | Agent | Summary |
| ---- | ----- | ------- |
| `25-49 HP/25 Identity (Self, Purpose & Spirit)/25 Identity.md` | Zenith | Core identity — bio, values, purpose, personal history |
| `25-49 HP/26 Spirituality (Metaphysics)/` | Zenith | Spiritual practice, metaphysics, beliefs |
| `25-49 HP/27 Mind (Education & Learning)/` | Zenith | Education, learning systems, knowledge acquisition |
| `25-49 HP/28 Body (Health & Wellness)/28 Body.md` | Zenith | Health, fitness, biometrics, sleep, nutrition, recovery |
| `25-49 HP/29 Relationships (Social)/` | Zenith | Social network, inner circle, friendships, mentors |
| `25-49 HP/30 Finance (Money)/` | Zenith | Personal finance, budget, income, investments, credit |
| `25-49 HP/31 Mastery (Skills & Training)/` | Zenith | Skills, certifications, technical + creative mastery |
| `25-49 HP/32 Rhythm (Routines & Productivity)/` | Zenith | Daily/weekly rhythms, productivity systems, habits |
| `25-49 HP/33 Sovereignty (Legal)/` | Zenith | Legal docs, contracts, personal sovereignty |
| `25-49 HP/34 Joy (Recreation & Travel)/` | Zenith | Recreation, hobbies, travel, leisure |

---

## LE Domains (depth: 2)

| Path | Agent | Summary |
| ---- | ----- | ------- |
| `50-74 LE/50 Home (Sanctuary & Property)/` | Unity | Home management, property, living environment |
| `50-74 LE/51 Covenant (Relationship & Union)/` | Unity | Partnership, marriage, relationship systems |
| `50-74 LE/52 Offspring (Parenting & Education)/` | Unity | Parenting, child development, education |
| `50-74 LE/53 Parents (Elder Care & Support)/` | Unity | Elder care planning, parent support |
| `50-74 LE/54 Heritage (History & Traditions)/` | Unity | Genealogy, oral history, family traditions, annual recaps |
| `50-74 LE/55 Health (Family Wellness)/` | Unity | Family health systems, wellness tracking |
| `50-74 LE/56 Estate (Family Wealth & Legal)/` | Unity | Estate planning, family legal, generational wealth |
| `50-74 LE/57 Community (Extended Network)/` | Unity | Community, extended network, local ties |

---

## GE Domains (depth: 2)

| Path | Agent | Summary |
| ---- | ----- | ------- |
| `75-99 GE/75 Zeus Delacruz (Brand & Empire)/` | Eternal | Personal brand, identity, content engine, speaking |
| `75-99 GE/76 Real Estate (Landlord & Property)/` | Eternal | Real estate portfolio, properties, landlord ops |
| `75-99 GE/77 Acquisitions (Bought Businesses)/` | Eternal | Acquired businesses, operations, performance |
| `75-99 GE/78 Holdings (Legal & Entity Structure)/` | Eternal | Legal entities, LLC structure, corporate governance |
| `75-99 GE/79 Investments (Wealth & Assets)/` | Eternal | Investment portfolio, assets, wealth tracking |
| `75-99 GE/80 Innovation (New Ventures Pipeline)/` | Eternal | New venture pipeline, R&D, innovation projects |
| `75-99 GE/81 Strategy (Portfolio Intelligence)/` | Eternal | Portfolio strategy, competitive intelligence |
| `75-99 GE/82 Partnerships & Capital (Network & Deals)/` | Eternal | Partnerships, deal flow, capital network |

---

## Intelligence Layer (depth: 2-3)

| Path | Agent | Summary |
| ---- | ----- | ------- |
| `00-24 SO/01 Intelligence/agents/_docs/AGENT-NAME-VAULT.md` | Sovereignty | Private agent name candidates — internal only |
| `00-24 SO/01 Intelligence/agents/_docs/AGENT-ASSETS.md` | Sovereignty | Agent asset inventory |
| `00-24 SO/01 Intelligence/skills/` | Any | All agent skills — organized by pillar (hp/le/ge/so/onemind) |
| `00-24 SO/01 Intelligence/memory/` | Legacy | Memory architecture, memory files |
| `00-24 SO/01 Intelligence/models/` | Sovereignty | Model configs, routing, provider docs |

---

## Key Reference Files (depth: 3-4)

| Path | Agent | Summary |
| ---- | ----- | ------- |
| `SCHEMA.md` | Any | Domain file format — frontmatter spec, field reference, priority tags, sync rules |
| `ONEMIND-CODEX.md` | Any | Canonical master document — do not casually edit |
| `AGENTS.md` | Any | Universal agent instructions |
| `CLAUDE.md` | Any | Claude Code agent context |
| `00-24 SO/04 Operations/command/Core/FIRE-AND-FIELD-COMMANDS.md` | Any | Tactical command language — 24 commands (v2) |
| `00-24 SO/06 OpenClaw/_docs/reference/MEMORY-STACK.md` | Sovereignty | OpenClaw 3-plugin memory stack reference |
| `25-49 HP/25 Identity/Core-Identity/Foundation/Personal-Philosophy.md` | Zenith | Zeus's personal philosophy, mantras, beliefs |
| `25-49 HP/28 Body/Body-Care/lip-scrub.md` | Zenith | DIY lip scrub recipes |
| `25-49 HP/31 Mastery/Technical-Engineering/Mechanical Engineering/3d-printing.md` | Zenith | 3D printing reference — materials, printers, slicers, repositories |
| `00-24 SO/05 Platforms/tools/circle/overview/circle-so.md` | Sovereignty | Circle MCP — full capability list |
| `00-24 SO/05 Platforms/tools/ai/anthropic/` | Sovereignty | Anthropic platform catalog — Claude AI/Code/API/CoWork |

---

## Maintenance

**To regenerate this index:** Legacy runs `onemind-maintainer` skill → `regenerate-index`.

**Manual update rule:** When creating a new depth:2 domain file or a high-value depth:3-4 reference doc, add a row to the relevant section above. Agent should do this automatically on any vault write.

**Note:** This is a curated index, not an exhaustive file list. For exhaustive search, use `list_directory` + `_README.md` traversal or the vector layer (Layer 2).
