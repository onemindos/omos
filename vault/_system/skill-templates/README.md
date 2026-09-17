# OpenClaw Skill Templates

A complete reference kit for building OpenClaw skills — from minimal stubs to production-ready integrations. Every template here follows both the open **Agent Skills** standard (agentskills.io) and OpenClaw's native Node.js execution model.

---

## What This Folder Contains

```
skill-templates/
├── README.md                       ← You are here — start here
├── SKILL-GUIDE.md                  ← Full deep-dive: spec, theory, best practices
├── templates/
│   ├── basic-skill/                ← Minimal starting point (use this first)
│   ├── api-integration/            ← REST API skill (most common pattern)
│   ├── nats-skill/                 ← NATS-connected skill (OneMind OS pattern)
│   └── advanced-skill/             ← Full-featured with scripts, references, assets
└── examples/
    ├── simple-weather/             ← Complete working weather skill example
    └── webhook-notifier/           ← Complete working webhook notifier example
```

---

## Quick Start — Build a New Skill in 5 Steps

### Step 1: Copy the right template

```bash
# For a REST API integration (most common):
cp -r skill-templates/templates/api-integration openclaw/workspaces/legacy/skills/my-skill-1.0.0

# For a basic utility skill:
cp -r skill-templates/templates/basic-skill openclaw/workspaces/legacy/skills/my-skill-1.0.0

# For an advanced skill with scripts/references:
cp -r skill-templates/templates/advanced-skill openclaw/workspaces/legacy/skills/my-skill-1.0.0
```

### Step 2: Edit `SKILL.md`

Update the frontmatter and write the instructions body:

```markdown
---
name: my-skill
description: What this skill does and exactly when Legacy should activate it.
---

# My Skill

## When to Use
...

## How It Works
...
```

### Step 3: Implement `index.js`

Export tool functions. Each exported async function becomes a live OpenClaw tool:

```javascript
async function my_tool({ param1, param2 }) {
  // do real work
  return { result: "..." };
}

module.exports = { my_tool };
```

### Step 4: Configure `package.json`

Update the `name` field to match, add any npm dependencies, then install:

```bash
cd openclaw/workspaces/legacy/skills/my-skill-1.0.0
npm install
```

### Step 5: Restart the gateway

```bash
# On VPS:
ssh root@104.248.4.212 "systemctl restart openclaw"

# Local dev:
openclaw gateway
```

---

## Choosing the Right Template

| Template | Use When |
|---|---|
| `basic-skill` | Simple utility, no external API, minimal deps |
| `api-integration` | Wrapping any REST API (auth, CRUD, pagination) |
| `nats-skill` | Need to publish/subscribe to NATS subjects across runtimes |
| `advanced-skill` | Complex domain with scripts, reference docs, static assets |

---

## OpenClaw vs. Agent Skills Standard

OpenClaw extends the open Agent Skills spec with **executable Node.js tools**.

| Feature | Agent Skills Spec | OpenClaw Extension |
|---|---|---|
| `SKILL.md` | ✅ Required | ✅ Same format |
| `scripts/` | Optional Bash/Python | Optional (Node.js via index.js) |
| `references/` | Optional docs | ✅ Supported |
| `assets/` | Optional static files | ✅ Supported |
| `index.js` | Not in spec | ✅ Exports live callable tools |
| `package.json` | Not in spec | ✅ Manages npm dependencies |
| Tool registration | Agent-defined | Automatic — exported functions = tools |
| Secret injection | Implementation-defined | `process.env.*` via `op run` |

---

## Skill Naming Conventions

| Pattern | Example | When to Use |
|---|---|---|
| `name-version` | `github-api-1.0.3` | External API integrations |
| `domain-version` | `home-assistant-1.0.0` | Platform integrations |
| `name` (no version) | `nexus-bridge` | Core infrastructure skills |
| `name` (no version) | `todoist` | Stable internal skills |

The `name` field in `SKILL.md` frontmatter must:
- Be lowercase, hyphens only
- Match the parent directory name
- Be ≤ 64 characters
- Not start or end with a hyphen

---

## Environment Variables

All secrets are injected by `op run` at gateway startup. Access them via `process.env`:

```javascript
const API_KEY = process.env.MY_SERVICE_API_KEY;
const BASE_URL = process.env.MY_SERVICE_URL || "https://api.example.com";
```

**Never hardcode secrets in skill files.** Add new secrets to:
1. 1Password vault `00-24 SO (Sovereign Operations)`
2. `/opt/onemind/.env.tpl` on VPS

---

## See Also

- [SKILL-GUIDE.md](./SKILL-GUIDE.md) — Complete theory, spec, and best practices
- [Agent Skills Spec](https://agentskills.io/specification) — Open standard reference
- [Legacy's live skills](../openclaw/workspaces/legacy/skills/) — 35 production skills
- [skill-forge skill](../openclaw/workspaces/legacy/skills/skill-forge-1.0.0/) — Autonomously build skills at runtime
