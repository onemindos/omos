---
name: adding-skills-to-openclaw
description: "Runbook for adding custom skills to OpenClaw via the onemind-skills local plugin"
type: reference
status: active
created: 2026-05-11
modified: 2026-05-11
owner: zeus
tags: [type/runbook, domain/so-01, stage/active, infrastructure/openclaw]
---

# Adding Skills to OpenClaw — Complete Process

> Skills are SKILL.md files. OpenClaw loads them via a local plugin we created called `onemind-skills`. Drop a folder, restart gateway, done.

---

## How It Works

OpenClaw has **no native `openclaw skills add` command**. Skills can only be registered through:

1. **Bundled** — shipped inside OpenClaw (56 built-in, can't modify)
2. **ClawHub** — `openclaw skills install <slug>` from their cloud registry (public)
3. **Plugin skills** — any installed plugin that declares `"skills": ["./skills"]` in its manifest

We use option 3: a **local plugin** called `onemind-skills` that acts as a container for all our custom SKILL.md files. No code runs — it's just a directory wrapper that OpenClaw's plugin system reads.

```
Vault Skill (source of truth — 01 Intelligence/Skills/)
    ↓  scp to VPS
Plugin skills dir (~/.openclaw/onemind-skills-plugin/skills/{name}/SKILL.md)
    ↓  gateway restart
OpenClaw loads as "openclaw-extra" skill → visible to all agents + dashboard
```

---

## Architecture

### VPS Plugin Location

```
~/.openclaw/onemind-skills-plugin/
├── openclaw.plugin.json   ← manifest (tells OpenClaw to scan ./skills)
├── package.json           ← metadata + openclaw.extensions field
├── index.js               ← empty no-op (required by plugin loader)
└── skills/
    ├── morning-briefing/
    │   └── SKILL.md
    ├── spartan/
    │   └── SKILL.md
    └── {any-new-skill}/
        └── SKILL.md
```

### Vault Skill Source Location

```
00-24 SO/01 Intelligence/Skills (Agent Skills)/
├── hp/    ← HP pillar skills (spartan, oracle, identity, etc.)
├── so/    ← SO pillar skills (grid, guardian, overwatch, etc.)
├── le/    ← LE pillar skills (haven, helix, eden, nomad, etc.)
├── ge/    ← GE pillar skills (empire, apex, creed, finance, etc.)
└── onemind/  ← System-wide skills (dashboard, inbox, sync, etc.)
```

---

## Add a New Skill (2 Steps)

### Step 1: Create SKILL.md in the vault

**Location:** `01 Intelligence/Skills (Agent Skills)/{pillar}/{skill-name}/SKILL.md`

**Minimal format:**

```yaml
---
name: my-skill
description: One-line description of when to use this skill.
---

# My Skill

Instructions for the agent when this skill is loaded.
```

That's it. Only `name` and `description` are required in frontmatter. The body is freeform markdown.

### Step 2: Deploy to VPS

```bash
# Create the skill folder on VPS and copy SKILL.md
ssh onemindos-core "mkdir -p ~/.openclaw/onemind-skills-plugin/skills/my-skill"
scp "path/to/my-skill/SKILL.md" onemindos-core:~/.openclaw/onemind-skills-plugin/skills/my-skill/SKILL.md

# Restart gateway to pick it up
ssh onemindos-core "systemctl --user restart openclaw-gateway"
```

### Verify

```bash
ssh onemindos-core "openclaw skills list | grep my-skill"
ssh onemindos-core "openclaw skills info my-skill"
```

Should show: `✓ ready | openclaw-extra | Visible to model: yes`

---

## Batch Deploy (All Skills at Once)

```bash
# From Mac — deploy entire pillar
SKILLS_DIR="00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Skills (Agent Skills)"

# Copy all skills from a pillar (e.g., hp)
for skill_dir in "$SKILLS_DIR/hp/"*/; do
  name=$(basename "$skill_dir")
  [ "$name" = "General (Overview)" ] && continue  # skip overviews
  ssh onemindos-core "mkdir -p ~/.openclaw/onemind-skills-plugin/skills/$name"
  scp "$skill_dir/SKILL.md" "onemindos-core:~/.openclaw/onemind-skills-plugin/skills/$name/SKILL.md"
done

# One restart at the end
ssh onemindos-core "systemctl --user restart openclaw-gateway"

# Verify count
ssh onemindos-core "openclaw skills list | grep openclaw-extra | wc -l"
```

---

## Remove a Skill

Delete the folder from the plugin and restart:

```bash
ssh onemindos-core "rm -rf ~/.openclaw/onemind-skills-plugin/skills/my-skill"
ssh onemindos-core "systemctl --user restart openclaw-gateway"
```

---

## Update a Skill

Just overwrite the SKILL.md and restart:

```bash
scp "path/to/my-skill/SKILL.md" onemindos-core:~/.openclaw/onemind-skills-plugin/skills/my-skill/SKILL.md
ssh onemindos-core "systemctl --user restart openclaw-gateway"
```

---

## CLI Reference (Actual Working Commands)

```bash
# List all skills (built-in + extras)
openclaw skills list

# List skills for specific agent
openclaw skills list --agent main

# Get skill details
openclaw skills info morning-briefing

# Check skill readiness
openclaw skills check --agent main

# Install from ClawHub (cloud registry — not our custom skills)
openclaw skills install <slug>

# Search ClawHub
openclaw skills search <query>
```

**These commands do NOT exist** (despite old docs claiming they do):
- ~~`openclaw skills add`~~ — no such command
- ~~`openclaw skills remove`~~ — no such command
- ~~`openclaw skills refresh`~~ — no such command
- ~~`openclaw config set agents[X].skills`~~ — skills are universal, not per-agent

---

## How the Plugin Was Created (One-Time Setup)

If you ever need to recreate the plugin from scratch:

```bash
ssh onemindos-core

# 1. Create plugin directory
mkdir -p ~/.openclaw/onemind-skills-plugin/skills

# 2. Create plugin manifest
cat > ~/.openclaw/onemind-skills-plugin/openclaw.plugin.json << 'EOF'
{
  "id": "onemind-skills",
  "enabledByDefault": true,
  "activation": { "onStartup": true },
  "skills": ["./skills"],
  "contracts": {},
  "configSchema": {
    "type": "object",
    "additionalProperties": false,
    "properties": {}
  }
}
EOF

# 3. Create package.json
cat > ~/.openclaw/onemind-skills-plugin/package.json << 'EOF'
{
  "name": "onemind-skills",
  "version": "1.0.0",
  "private": true,
  "description": "One Mind custom skills plugin for OpenClaw",
  "type": "module",
  "openclaw": { "extensions": ["./index.js"] }
}
EOF

# 4. Create empty index.js (required by loader, does nothing)
echo 'export default function activate() {}' > ~/.openclaw/onemind-skills-plugin/index.js

# 5. Install as linked plugin
openclaw plugins install --link ~/.openclaw/onemind-skills-plugin/

# 6. Restart gateway
systemctl --user restart openclaw-gateway
```

---

## Troubleshooting

**Skill not showing after restart:**
- Check file exists: `ls ~/.openclaw/onemind-skills-plugin/skills/my-skill/SKILL.md`
- Check plugin is enabled: `openclaw plugins list | grep onemind`
- Check gateway logs: `journalctl --user -u openclaw-gateway --since "1 min ago" | grep -i skill`

**EINVAL error in logs:**
- Means a real directory exists in `~/.openclaw/plugin-skills/` with the same name (old manual placement). Delete it: `rm -rf ~/.openclaw/plugin-skills/my-skill`

**Plugin not loading:**
- Verify manifest: `cat ~/.openclaw/onemind-skills-plugin/openclaw.plugin.json | python3 -m json.tool`
- Must have `configSchema`, `contracts`, and `skills` fields
- `package.json` must have `openclaw.extensions` pointing to `index.js`

---

## Why a Plugin?

OpenClaw has **no native command to add local skills**. The only 3 paths:

1. **Bundled** — shipped inside OpenClaw (can't modify)
2. **ClawHub** — cloud registry (`openclaw skills install <slug>`) — public, hosted externally
3. **Plugin skills** — any enabled plugin with `"skills": ["./skills"]` in its manifest

There is no `openclaw skills add`, no `skills.customDirs` config, no local install path. The plugin wrapper is the correct and only way to register local custom skills that appear in the dashboard, CLI, and agent sessions.

---

## Process Checklist — New Skill

```
□ Create skill folder in vault: 01 Intelligence/Skills/{pillar}/my-skill/SKILL.md
□ Frontmatter: name + description (minimum required)
□ SCP to VPS: scp SKILL.md onemindos-core:~/.openclaw/onemind-skills-plugin/skills/my-skill/
□ Restart gateway: ssh onemindos-core "systemctl --user restart openclaw-gateway"
□ Verify: ssh onemindos-core "openclaw skills info my-skill"
□ Done — skill is live for all agents
```

---

## Agent Skill Visibility (Current — 2026-05-11)

All 46 custom skills are visible to **every agent** via the `<available_skills>` menu. Skills are NOT assigned per-agent — OpenClaw has no per-agent skill allowlist. Each agent's TOOLS.md guides it on which skills are most relevant to its domain, but any agent can invoke any skill.

**46 custom skills deployed** — see `02 Infrastructure/openclaw/plugin-skills/PLUGIN-SKILLS.md` for the full categorized list.

| Pillar | Count | Skills |
|--------|-------|--------|
| HP | 12 | spartan, oracle, genesis, human-enhancement, librarian, pulse, morning-briefing, evening-review, weekly-scorecard, life-score-intelligence, user-persona, ikigai-purpose-discovery |
| SO | 8 | grid, guardian, overwatch, eagleeye, signals, transit, problem-solving, research |
| LE | 6 | eden, homestead, nomad, haven, helix, legacy-library-books |
| GE | 8 | empire, creed, trinity, cfo, video-production, webflow, circle, contact-crm-relationships |
| OM | 12 | code-method, entity-system, om-dashboard, onemind-architect, onemind-defrag, onemind-inbox, onemind-maintainer, onemind-onboarding, onemind-report, onemind-search, onemind-status, onemind-sync |

---

## Related

- `SKILL-LAYER-ARCHITECTURE.md` — Universal model explanation
- `MULTI-SYSTEM-ROUTING.md` — Full skill assignment matrix
- `../../02 Infrastructure/openclaw/` — OpenClaw docs
- [OpenClaw CLI Docs](https://docs.openclaw.ai/cli)
