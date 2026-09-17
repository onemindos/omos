---
title: "OpenClaw — Custom Plugin Skills"
type: reference
status: active
created: 2026-05-10
modified: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO-02, openclaw/skills]
---

# OpenClaw — Custom Plugin Skills

Custom skills deployed via the `onemind-skills` local plugin on omos-core. All agents see these as `openclaw-extra` skills in the gateway.

**Last verified:** 2026-05-11 — 46 custom skills loaded, 0 errors, all visible to model.

---

## Architecture — How Skills Work

### Loading Model

Skills are **NOT** loaded into the agent's context window as full text. OpenClaw uses a two-stage system:

1. **Metadata menu:** On every message, the agent sees a compact `<available_skills>` XML block in its system prompt containing just the **name**, **description** (from SKILL.md frontmatter), and **file location** of each skill. All 46 skills fit in ~2-3K tokens.

2. **On-demand loading:** The system prompt instructs the agent:
   - Scan `<available_skills>` descriptions before replying
   - If one skill clearly matches → read its full SKILL.md via the read file tool
   - If multiple match → pick the most specific, read it
   - If none match → don't read any
   - Never read more than one skill upfront

This means the full SKILL.md content only enters the context window when the agent determines it's relevant to the current message.

### Skill Visibility

All 46 custom skills are visible to **every agent** — Legacy, Zenith, Unity, Eternal, Sovereignty, Sentinel. The workspace files (IDENTITY.md, TOOLS.md) guide each agent on which skills are most relevant to its domain, but any agent can invoke any skill.

### Plugin Architecture

Skills are registered through a **local plugin** called `onemind-skills`. No code runs — it's a directory structure that OpenClaw's plugin system reads.

```
~/.openclaw/onemind-skills-plugin/
├── openclaw.plugin.json   ← manifest: { skills: ["./skills"] }
├── package.json           ← metadata + openclaw.extensions
├── index.js               ← empty no-op (required by loader)
└── skills/
    ├── cfo/SKILL.md
    ├── spartan/SKILL.md
    ├── oracle/SKILL.md
    └── ... (46 skill folders)
```

Installed via: `openclaw plugins install --link ~/.openclaw/onemind-skills-plugin/`

The `--link` flag means OpenClaw reads the folder in-place. On gateway restart, OpenClaw auto-generates `~/.openclaw/plugin-skills/` as a **symlink cache** — this is normal runtime behavior, not stale data.

### Source of Truth

- **Vault (authoritative):** `01 Intelligence/Skills (Agent Skills)/` — organized by pillar subfolder
- **VPS (deployed):** `~/.openclaw/onemind-skills-plugin/skills/` — flat slug-named folders
- **Cut archive:** `01 Intelligence/Skills (Agent Skills)/_cut/` — 15 retired skills preserved

---

## Deployed Skills — 46 Custom (as of 2026-05-11)

### HP — Holistic Performance (12 skills)

| Skill | Domain | Purpose |
|-------|--------|---------|
| `spartan` | HP 28 | Body, fitness, training, nutrition, recovery |
| `oracle` | HP 25-34 | Mind, habits, personal finance, investments |
| `genesis` | HP 25-26 | Identity, life purpose, values, spiritual practice |
| `human-enhancement` | HP 27 | Capability development, skill stacking |
| `librarian` | HP 27 | External intelligence gathering, research, daily briefs |
| `pulse` | HP 27 | Real-time feed monitoring, news, signal detection |
| `morning-briefing` | HP 32 | Morning brief — priorities, calendar, energy |
| `evening-review` | HP 32 | Evening review — accomplishments, lessons, tomorrow prep |
| `weekly-scorecard` | HP 32 | Weekly review — task completions, domain health, scores |
| `life-score-intelligence` | HP 25 | Life domain scoring and tracking |
| `user-persona` | HP 25 | Zeus's detailed profile for context |
| `ikigai-purpose-discovery` | HP 25 | Everyday alignment and purpose check |

### SO — Sovereign Operations (8 skills)

| Skill | Domain | Purpose |
|-------|--------|---------|
| `grid` | SO 02 | VPS health, Docker, systemd, databases, backups, SSH |
| `guardian` | SO 04 | Cybersecurity, credentials, 1Password, SSL, legal deadlines |
| `overwatch` | SO 07 | n8n workflow management, Node-RED, automation health |
| `eagleeye` | SO 02 | Drone operations, computer vision, physical automation |
| `signals` | SO 03 | Interrupt and alert layer — watches for critical events |
| `transit` | SO 02 | Network, connectivity, VPN, Tailscale management |
| `problem-solving` | SO 05 | Define-gather-generate-evaluate-decide framework |
| `research` | SO 06 | Multi-source, cross-referenced research |

### LE — Legacy Evolution (6 skills)

| Skill | Domain | Purpose |
|-------|--------|---------|
| `eden` | LE | Sustainability, environment, home ecology |
| `homestead` | LE 50 | Farm management, crops, livestock, irrigation |
| `nomad` | LE 50 | Off-grid systems, emergency prep, sovereign living |
| `haven` | LE 50-74 | Family, relationships, home management |
| `helix` | LE 52-56 | Parenting, elder care, generational planning |
| `legacy-library-books` | LE | Book management, reading lists, library |

### GE — Generational Empire (8 skills)

| Skill | Domain | Purpose |
|-------|--------|---------|
| `empire` | GE 75-99 | Business strategy, ventures, brand building |
| `creed` | GE 75-81 | Brand narrative, content strategy, storytelling |
| `trinity` | SO 06 | Deep R&D, emerging tech, competitive intelligence |
| `cfo` | GE 79 | Capital, fundraising, financial modeling, compliance |
| `video-production` | GE 78 | Video content, scripting, CapCut, voice clips |
| `webflow` | GE 75 | Webflow sites, CMS, publishing, brand content |
| `circle` | GE 81 | Community platform management |
| `contact-crm-relationships` | GE 82 | Contact management, CRM, networking |

### OM — One Mind System (12 skills)

| Skill | Domain | Purpose |
|-------|--------|---------|
| `code-method` | Framework | CODE× method — Capture, Organize, Direct, Execute, × Multiply |
| `entity-system` | Framework | 8 entity types — notes, projects, areas, resources, etc. |
| `om-dashboard` | Operations | Dashboard reads, system overview |
| `onemind-architect` | Framework | Community-facing architecture, vault structure design |
| `onemind-defrag` | Maintenance | Restructure and reorganize vault content |
| `onemind-inbox` | Operations | Process and clear the inbox queue |
| `onemind-maintainer` | Operations | Read/write/maintain vault domain docs |
| `onemind-onboarding` | Community | Teach new students the One Mind system |
| `onemind-report` | Reporting | Generate system reports — task counts, domain health |
| `onemind-search` | Operations | Find information across the vault |
| `onemind-status` | Health | Full system meta-layer health check |
| `onemind-sync` | Sync | Keep vault and Todoist in sync |

---

## Bundled Skills (from OpenClaw — not custom)

These load automatically from OpenClaw's built-in extensions. Not managed by the plugin.

| Skill | Source | Status |
|-------|--------|--------|
| `1password` | bundled | ✓ ready |
| `browser-automation` | bundled (browser ext) | ✓ ready |
| `camsnap` | bundled | ✓ ready |
| `clawhub` | bundled | ✓ ready |
| `coding-agent` | bundled | ✓ ready |
| `gemini` | bundled | ✓ ready |
| `gh-issues` | bundled | ✓ ready |
| `github` | bundled | ✓ ready |
| `healthcheck` | bundled | ✓ ready |
| `mcporter` | bundled | ✓ ready |
| `nano-pdf` | bundled | ✓ ready |
| `node-connect` | bundled | ✓ ready |
| `obsidian` | bundled | ✓ ready |
| `obsidian-vault-maintainer` | bundled (memory-wiki ext) | ✓ ready |
| `openai-whisper` | bundled | ✓ ready |
| `openai-whisper-api` | bundled | ✓ ready |
| `oracle` (CLI) | bundled | ✓ ready |
| `session-logs` | bundled | ✓ ready |
| `skill-creator` | bundled | ✓ ready |
| `summarize` | bundled | ✓ ready |
| `taskflow` | bundled | ✓ ready |
| `taskflow-inbox-triage` | bundled | ✓ ready |
| `tavily` | bundled | ✓ ready |
| `tmux` | bundled | ✓ ready |
| `video-frames` | bundled | ✓ ready |
| `weather` | bundled | ✓ ready |
| `wiki-maintainer` | bundled | ✓ ready |
| `voice-call` | bundled | ⏸ disabled |

28 additional bundled skills show "missing requirements" — they need CLIs or tokens not installed on the VPS (e.g. `discord`, `slack`, `apple-notes`, `xurl`). These are not errors.

---

## Cut Skills Archive — 15 Retired (2026-05-11)

These were removed during the 67→46 consolidation. Source files preserved in `01 Intelligence/Skills (Agent Skills)/_cut/`.

| Cut Skill | Reason |
|-----------|--------|
| `apex` | Merged into `cfo` |
| `code-review` | Generic — covered by `coding-agent` (bundled) |
| `daily-briefing` | Replaced by `morning-briefing` + `evening-review` |
| `debugging` | Generic — covered by `coding-agent` (bundled) |
| `empire-ops` | Merged into `empire` |
| `farm-operations` | Merged into `homestead` |
| `guardian-ops` | Merged into `guardian` |
| `likekit-platform` | Product discontinued |
| `onemind-consolidate` | Merged into `onemind-defrag` |
| `onemind-init` | One-time setup — no longer needed |
| `onemind-maintain` | Merged into `onemind-maintainer` |
| `project-planning` | Generic — covered by `taskflow` (bundled) |
| `task-management` | Generic — covered by `taskflow` (bundled) |
| `three-paths` | Experimental — not used |
| `writing` | Generic — covered by base model capability |

6 agent pack skills (`agent-legacy`, `agent-zenith`, `agent-unity`, `agent-eternal`, `agent-sovereignty`, `agent-sentinel`) were also removed from OpenClaw deployment — redundant with workspace files (IDENTITY.md, SOUL.md, TOOLS.md, HEARTBEAT.md) that auto-load per agent. Source files kept in vault.

---

## Agent Workspace Files (Auto-Loaded)

Each agent's workspace at `~/.openclaw/workspace-{id}/` loads these files automatically at startup — separate from the skills system:

| File | Content |
|------|---------|
| `IDENTITY.md` | Name, ID, emoji, node, pillar, role description |
| `SOUL.md` | Mission, core values, rules, autonomy levels |
| `TOOLS.md` | Always-available + on-demand skill recommendations |
| `HEARTBEAT.md` | Cron schedule, escalation rules, health checks |
| `AGENTS.md` | Fleet roster — who the other agents are |
| `MEMORY.md` | Memory context and persistent state |
| `USER.md` | Zeus's preferences, communication style, details |

These define the agent's identity and behavior. Skills define domain expertise loaded on demand.

---

## Adding a New Skill — Exact Flow

Follow this every time. No exceptions.

### Step 1: Write the SKILL.md in the vault

Create the file at the correct pillar subfolder:

```
01 Intelligence/Skills (Agent Skills)/
├── hp/my-skill (Descriptive Name)/SKILL.md   ← HP skill
├── so/my-skill (Descriptive Name)/SKILL.md   ← SO skill
├── le/my-skill (Descriptive Name)/SKILL.md   ← LE skill
├── ge/my-skill (Descriptive Name)/SKILL.md   ← GE skill
└── onemind/my-skill (Descriptive Name)/SKILL.md  ← System skill
```

**Minimum required SKILL.md format:**

```yaml
---
name: my-skill
description: One-line description of when to use this skill. This is what appears in the agent's available_skills menu — make it clear and trigger-based.
---

# My Skill Name

Full instructions for the agent when this skill is loaded on demand.
Include: what to do, what tools to use, what files to read/write,
what format to output, what domain this maps to.
```

**Rules:**
- `name` = the slug (lowercase, hyphens, no spaces) — this becomes the folder name on VPS
- `description` = what the agent scans to decide whether to load this skill — include trigger phrases
- Body = full instructions only loaded when the skill is invoked

### Step 2: Deploy to VPS

```bash
# Create folder with the slug name and copy the SKILL.md
SKILL_NAME="my-skill"
VAULT_PATH="00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Skills (Agent Skills)/hp/${SKILL_NAME} (Descriptive Name)/SKILL.md"

ssh onemindos-core "mkdir -p ~/.openclaw/onemind-skills-plugin/skills/${SKILL_NAME}"
scp "${VAULT_PATH}" "onemindos-core:~/.openclaw/onemind-skills-plugin/skills/${SKILL_NAME}/SKILL.md"
```

### Step 3: Restart gateway

```bash
ssh onemindos-core "systemctl --user restart openclaw-gateway"
```

### Step 4: Verify

```bash
# Check it loaded
ssh onemindos-core "openclaw skills list 2>&1 | grep my-skill"
# Expected: ✓ ready | 📦 my-skill | description... | openclaw-extra

# Full health check
ssh onemindos-core "openclaw skills check 2>&1 | head -12"

# Count total custom skills
ssh onemindos-core "openclaw skills list 2>&1 | grep -c 'openclaw-extra'"
```

### Step 5: Update this doc

Add the skill to the appropriate pillar table in the "Deployed Skills" section above.

### Checklist (copy-paste for each new skill)

```
□ SKILL.md created in vault: 01 Intelligence/Skills/{pillar}/{slug} (Name)/SKILL.md
□ Frontmatter has name (slug) + description (trigger-based one-liner)
□ Body has full agent instructions
□ SCP'd to VPS: ~/.openclaw/onemind-skills-plugin/skills/{slug}/SKILL.md
□ Gateway restarted: systemctl --user restart openclaw-gateway
□ Verified: openclaw skills list | grep {slug} → ✓ ready
□ PLUGIN-SKILLS.md updated with new entry
□ Git committed
```

---

## Removing a Skill

```bash
# 1. Remove from VPS
ssh onemindos-core "rm -rf ~/.openclaw/onemind-skills-plugin/skills/my-skill"
ssh onemindos-core "systemctl --user restart openclaw-gateway"

# 2. Move vault source to _cut/ archive (don't delete)
mv "01 Intelligence/Skills/{pillar}/my-skill (Name)/" "01 Intelligence/Skills/_cut/my-skill/"

# 3. Update PLUGIN-SKILLS.md — remove from deployed table, add to cut archive table
```

---

## Updating a Skill

Just overwrite the SKILL.md and restart:

```bash
scp "path/to/SKILL.md" onemindos-core:~/.openclaw/onemind-skills-plugin/skills/my-skill/SKILL.md
ssh onemindos-core "systemctl --user restart openclaw-gateway"
```

---

## Bulk Deploy (All Skills at Once)

Use when redeploying the full set or after major changes:

```bash
# 1. Stage locally with clean slug names (no parenthetical suffixes)
mkdir -p /tmp/onemind-skills-deploy
SKILLS_BASE="00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Skills (Agent Skills)"

for pillar in hp so le ge onemind; do
  for skill_dir in "${SKILLS_BASE}/${pillar}/"*/; do
    name=$(basename "$skill_dir" | sed 's/ (.*//')  # strip parenthetical
    [ "$name" = "General" ] && continue
    [ "$name" = "_cut" ] && continue
    [ -f "${skill_dir}/SKILL.md" ] || continue
    mkdir -p "/tmp/onemind-skills-deploy/${name}"
    cp "${skill_dir}/SKILL.md" "/tmp/onemind-skills-deploy/${name}/SKILL.md"
  done
done

# 2. Verify count before deploying
echo "Skills staged: $(ls -1 /tmp/onemind-skills-deploy/ | wc -l)"

# 3. Rsync to VPS (--delete removes skills no longer in staging)
rsync -avz --delete /tmp/onemind-skills-deploy/ onemindos-core:~/.openclaw/onemind-skills-plugin/skills/

# 4. Restart and verify
ssh onemindos-core "systemctl --user restart openclaw-gateway && sleep 3 && openclaw skills list 2>&1 | grep -c 'openclaw-extra'"

# 5. Clean up
rm -rf /tmp/onemind-skills-deploy
```

---

## Health Check Commands

```bash
# Full readiness check (shows eligible, visible, missing requirements)
ssh onemindos-core "openclaw skills check"

# Count custom skills loaded
ssh onemindos-core "openclaw skills list 2>&1 | grep -c 'openclaw-extra'"

# Verify all custom skills are visible to model
ssh onemindos-core "openclaw skills check 2>&1 | sed -n '/Ready and visible/,/Missing requirements/p' | grep '📦'"

# Check for errors in gateway logs
ssh onemindos-core "journalctl --user -u openclaw-gateway --since '5 min ago' --no-pager 2>&1 | grep -iE 'error|warn|fail|skill'"

# Verify specific skill
ssh onemindos-core "openclaw skills list 2>&1 | grep my-skill"
```

---

## CLI Commands That DO NOT Exist

Despite what you might expect, OpenClaw has **no**:
- ~~`openclaw skills add`~~ — use the plugin folder method above
- ~~`openclaw skills remove`~~ — delete the folder and restart
- ~~`openclaw skills refresh`~~ — restart the gateway instead
- ~~`openclaw config set agents[X].skills`~~ — skills are universal, not per-agent

---

## Configuration

Skills can be enabled/disabled in `~/.openclaw/openclaw.json`:

```json
{
  "skills": {
    "entries": {
      "voice-call": { "enabled": false },
      "coding-agent": { "enabled": true }
    }
  }
}
```

Limits (defaults usually fine — only change if hitting caps):

| Setting | What |
|---------|------|
| `maxSkillsLoadedPerSource` | Max skills per source (bundled/extra/workspace) |
| `maxSkillsInPrompt` | Max skills in the `<available_skills>` block |
| `maxSkillsPromptChars` | Max characters for the skills metadata section |
| `maxSkillFileBytes` | Max size for a single SKILL.md to be considered |

---

## Plugin Recreation (If Ever Needed)

If the plugin gets corrupted or you need to set it up on a new node:

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

# 4. Create empty index.js (required by plugin loader)
echo 'export default function activate() {}' > ~/.openclaw/onemind-skills-plugin/index.js

# 5. Install as linked plugin
openclaw plugins install --link ~/.openclaw/onemind-skills-plugin/

# 6. Deploy skills via bulk rsync (see Bulk Deploy section above)

# 7. Restart gateway
systemctl --user restart openclaw-gateway
```

---

## Related

- `01 Intelligence/Skills (Agent Skills)/ADDING-SKILLS-TO-OPENCLAW.md` — Detailed runbook (same flow)
- `01 Intelligence/Skills (Agent Skills)/SKILL-LAYER-ARCHITECTURE.md` — How the skill system works
- [[../workspaces (Agent Workspaces)/README|Workspaces]] — Agent workspace files (IDENTITY/SOUL/TOOLS)
- [[../cron (Scheduled Jobs)/CRON|CRON]] — Scheduled jobs per agent
- [[../hooks (Event Hooks)/HOOKS|HOOKS]] — Event hooks
