# Complete Skill Building Guide

Everything you need to know about Agent Skills — the open standard, the theory, the spec, OpenClaw's implementation, and proven best practices for writing skills that actually work.

---

## Table of Contents

1. [What Are Agent Skills?](#1-what-are-agent-skills)
2. [Why Skills Exist — The Context Window Problem](#2-why-skills-exist-the-context-window-problem)
3. [How Skills Work — Progressive Disclosure](#3-how-skills-work-progressive-disclosure)
4. [The SKILL.md Specification](#4-the-skillmd-specification)
5. [Optional Directories](#5-optional-directories)
6. [OpenClaw's Node.js Extension Model](#6-openclaws-nodejs-extension-model)
7. [The Power of Skills](#7-the-power-of-skills)
8. [Writing Effective Instructions](#8-writing-effective-instructions)
9. [Skill Design Patterns](#9-skill-design-patterns)
10. [Secrets and Environment](#10-secrets-and-environment)
11. [Testing and Debugging](#11-testing-and-debugging)
12. [Checklist Before Shipping](#12-checklist-before-shipping)

---

## 1. What Are Agent Skills?

Agent Skills are an **open standard** (originally developed by Anthropic, now maintained at agentskills.io) for giving AI agents portable, on-demand capabilities.

**Core definition**: A skill is a directory containing, at minimum, a `SKILL.md` file. That file holds YAML metadata and Markdown instructions. Everything else is optional.

```
my-skill/
├── SKILL.md          # Required — metadata + instructions
├── scripts/          # Optional — executable code
├── references/       # Optional — deep documentation
└── assets/           # Optional — templates, schemas, static files
```

**Why it's an open standard**: The same skill folder can run on Claude Code, OpenHands, Databricks, Snowflake Cortex, Spring AI, and OpenClaw — any agent runtime that implements the spec. Build once, deploy anywhere compatible.

**The atomic unit**: A skill encodes one domain of expertise or one integration. It is not a monolith. Keep skills focused.

---

## 2. Why Skills Exist — The Context Window Problem

Even with 200K–1M token context windows, agents face a fundamental tradeoff:

> **Everything loaded into context costs tokens, slows inference, and increases noise.**

Imagine Legacy has 35 skills. If you dump all 35 skill guides into the system prompt:
- That's ~175,000 tokens of mostly irrelevant context every single conversation
- The model spends attention on "how to control smart home lights" when the user is asking about GitHub PRs
- Latency increases, cost increases, accuracy decreases

**Without skills**, you're forced into one of three bad choices:

| Approach | Problem |
|---|---|
| Cram everything into system prompt | Bloated, slow, expensive, noisy |
| Build one giant monolithic agent | Brittle, unmaintainable, hard to update |
| Make agents that do nothing domain-specific | Useless for real work |

**Skills solve this** by applying the principle of progressive disclosure — agents only load what they need, when they need it.

---

## 3. How Skills Work — Progressive Disclosure

Skills load in three stages, each costing progressively more context:

### Stage 1: Discovery (Startup)

When the OpenClaw gateway boots, it loads **only the `name` and `description`** from every skill's frontmatter.

- Cost: ~100 tokens per skill
- 35 skills = ~3,500 tokens total at startup
- Legacy now "knows" what skills exist and what they do

### Stage 2: Activation (Task Matched)

When a user request matches a skill's description, the agent reads the **full `SKILL.md` body** into context.

- Cost: up to ~5,000 tokens (recommended limit)
- This is the procedure, examples, and reference for that specific skill
- Legacy now knows exactly how to execute the task

### Stage 3: Execution (As Needed)

The agent follows the instructions, loading **referenced files or running bundled scripts** only when the procedure requires it.

- Cost: only what's needed
- A deep `REFERENCE.md` is never loaded unless the skill's instructions specifically call for it
- Scripts are executed, not just read

**The result**: An agent with 35 skills pays a ~3,500-token overhead at startup, then loads only 1-3 skill docs per conversation, rather than paying 175,000 tokens every time.

---

## 4. The SKILL.md Specification

Every skill MUST have a `SKILL.md` file. This is the complete spec.

### Frontmatter Fields

The YAML block at the top of the file (between `---` markers):

```yaml
---
name: my-skill-name
description: >
  What this skill does and exactly when to use it. Include specific
  keywords that help the agent recognize relevant tasks.
license: MIT
compatibility: Requires Node.js 18+. Network access to api.example.com.
metadata:
  author: zeus
  version: "1.0.0"
  domain: infrastructure
allowed-tools: Bash(git:*) Read Write
---
```

#### `name` (Required)

- Must match the parent directory name exactly
- 1–64 characters
- Only lowercase letters (`a-z`), numbers (`0-9`), and hyphens (`-`)
- Cannot start or end with a hyphen
- No consecutive hyphens (`--`)

```
✅ github-api
✅ home-assistant-1.0.0
✅ nexus-bridge
❌ GitHub-API      (uppercase)
❌ -my-skill       (starts with hyphen)
❌ my--skill       (consecutive hyphens)
```

#### `description` (Required)

- 1–1,024 characters
- Must describe **what** the skill does AND **when** to activate it
- Include specific keywords agents use when matching tasks

```yaml
# Good: specific, includes trigger keywords
description: >
  Controls Home Assistant smart home devices, scenes, and automations via REST API.
  Use when the user asks about lights, locks, climate, sensors, cameras, or home status.
  Supports turn on/off, brightness, temperature, and scene activation.

# Bad: vague, no trigger keywords
description: Helps control the home.
```

#### `license` (Optional)

Short license name or reference to a bundled `LICENSE.txt`:

```yaml
license: MIT
license: Apache-2.0
license: Proprietary — see LICENSE.txt
```

#### `compatibility` (Optional, ≤500 chars)

Use only if the skill has specific requirements:

```yaml
compatibility: Requires Node.js 18+, nats npm package, and a running NATS server.
compatibility: Designed for OpenClaw. Requires HOME_ASSISTANT_URL and HOME_ASSISTANT_TOKEN env vars.
```

#### `metadata` (Optional)

Arbitrary key-value pairs — anything useful for tooling, auditing, or future Agent Skills extensions:

```yaml
metadata:
  author: zeus
  version: "1.0.0"
  domain: infrastructure
  forge: true           # built by skill-forge
  emoji: "⚡"
```

#### `allowed-tools` (Optional, Experimental)

Space-delimited list of tools the skill is pre-approved to use without user confirmation:

```yaml
allowed-tools: Bash(git:*) Bash(npm:*) Read Write
```

### Markdown Body

Everything after the frontmatter. No format restrictions — write what helps Legacy execute the task reliably.

**Recommended structure:**

```markdown
# Skill Name

## When to Use
Clear, specific trigger conditions. When should Legacy activate this?

## Architecture / How It Works
Brief explanation of the system being controlled. What are the key concepts?

## Available Tools
| Tool | Description |
|------|-------------|
| tool_name | What it does |

## Quick Reference / Examples

### Common Task 1
\`\`\`
tool_call({ param: "value" })
\`\`\`

### Common Task 2
\`\`\`
tool_call({ param: "value", options: { flag: true } })
\`\`\`

## Edge Cases / Notes
Things the agent should know but would likely get wrong otherwise.
```

**Keep `SKILL.md` under 500 lines.** Move heavy reference material to `references/REFERENCE.md` and link to it.

---

## 5. Optional Directories

### `scripts/`

Executable code the agent can run directly. Scripts should:

- Be self-contained, or clearly document their dependencies
- Handle errors gracefully with useful messages
- Accept arguments via stdin or CLI flags

```
scripts/
├── extract.py        # Python extraction script
├── validate.sh       # Bash validation
└── transform.js      # Node.js transformation
```

In `SKILL.md`, reference scripts with relative paths:

```markdown
Run the extraction script:

\`\`\`bash
python scripts/extract.py --input file.pdf --output out.txt
\`\`\`
```

### `references/`

Deep documentation loaded on demand. Keep individual files focused and small — the agent only loads these when the skill's instructions specifically redirect to them.

```
references/
├── REFERENCE.md      # Full API reference
├── FORMS.md          # Form templates
├── AUTH.md           # Auth flow details
└── ERRORS.md         # Error codes and remediation
```

### `assets/`

Static files the agent uses during execution:

```
assets/
├── schema.json           # JSON schema for validation
├── template.md           # Document template to fill in
├── config.example.json   # Config file template
└── diagrams/             # Architecture diagrams (for context)
```

---

## 6. OpenClaw's Node.js Extension Model

OpenClaw extends the open Agent Skills standard with **live executable tools**. Every exported async function in `index.js` becomes a registered OpenClaw tool that Legacy can call directly.

### The index.js Contract

```javascript
"use strict";

// Exported functions = registered tools
async function my_tool({ param1, param2 = "default" }) {
  // Real work happens here
  const result = await doSomething(param1);
  return { status: "ok", data: result };
}

async function another_tool({ required_param }) {
  if (!required_param) throw new Error("required_param is required");
  return { result: "..." };
}

// Only exported functions become tools
module.exports = { my_tool, another_tool };
```

**Rules for exported tool functions:**

- Must be `async` (or return a Promise)
- Accept a single destructured object parameter
- Return a plain object (serialized back to Legacy)
- Throw an `Error` with a clear message on failure (OpenClaw catches and surfaces it)
- Name matches what you reference in `SKILL.md`

### package.json

```json
{
  "name": "my-skill",
  "version": "1.0.0",
  "description": "Short description matching SKILL.md",
  "main": "index.js",
  "dependencies": {
    "node-fetch": "^3.3.2"
  }
}
```

### Tool Parameter Patterns

**Required parameter:**
```javascript
async function create_item({ title, project_id }) {
  if (!title) throw new Error("title is required");
  if (!project_id) throw new Error("project_id is required");
  // ...
}
```

**Optional parameter with default:**
```javascript
async function list_items({ project_id, limit = 50, offset = 0 }) {
  // ...
}
```

**Flexible input (one of several modes):**
```javascript
async function control_device({ entity_id, action, brightness, temperature }) {
  if (!entity_id) throw new Error("entity_id is required");
  if (!action) throw new Error("action is required (turn_on, turn_off, toggle)");
  // brightness and temperature are optional, only used when action supports them
}
```

### Environment Variable Pattern

```javascript
// At the top of index.js — load once, fail fast
const BASE_URL = process.env.MY_SERVICE_URL;
const API_KEY = process.env.MY_SERVICE_API_KEY;

if (!BASE_URL) throw new Error("MY_SERVICE_URL env var is required");
if (!API_KEY) throw new Error("MY_SERVICE_API_KEY env var is required");
```

Or defer validation to tool execution (better for optional skills):

```javascript
async function my_tool({ param }) {
  const apiKey = process.env.MY_SERVICE_API_KEY;
  if (!apiKey) throw new Error("MY_SERVICE_API_KEY is not configured");
  // ...
}
```

### HTTP Request Pattern (fetch)

```javascript
async function api_call({ endpoint, method = "GET", body = null }) {
  const BASE_URL = process.env.MY_SERVICE_URL || "https://api.example.com";
  const API_KEY = process.env.MY_SERVICE_API_KEY;

  const res = await fetch(`${BASE_URL}${endpoint}`, {
    method,
    headers: {
      "Authorization": `Bearer ${API_KEY}`,
      "Content-Type": "application/json",
    },
    body: body ? JSON.stringify(body) : undefined,
    signal: AbortSignal.timeout(30000), // 30s timeout
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`API error ${res.status}: ${text}`);
  }

  return res.json();
}
```

### NATS Publish Pattern

```javascript
const { connect, JSONCodec } = require("nats");

const NATS_URL = process.env.NATS_URL || "nats://localhost:4222";
const jc = JSONCodec();
let _nc = null;

async function getNats() {
  if (_nc) return _nc;
  _nc = await connect({ servers: NATS_URL, timeout: 3000 });
  return _nc;
}

async function publish_event({ subject, data }) {
  const nc = await getNats();
  nc.publish(subject, jc.encode(data));
  return { published: true, subject };
}

module.exports = { publish_event };
```

---

## 7. The Power of Skills

### Procedural Reliability

Natural language instructions ("figure out how to do X") produce variable results. A skill encodes the **exact procedure** — the specific API calls, the correct sequence of operations, the edge cases. Legacy stops improvising on high-stakes tasks and starts executing a known-good procedure.

### Organizational Knowledge, Versioned

Skills can encode knowledge that exists nowhere else:
- Your homestead's specific camera naming convention
- Your exact PR review checklist
- Your deployment procedure with all the non-obvious steps
- Domain-specific heuristics built up over months

That knowledge is now **git-versioned**, auditable, shareable, and reproducible.

### Domain Expertise On Demand

A `legal-review` skill can encode contract review procedures. A `credit-repair` skill can encode dispute letter strategy. A `recipe-scraper` skill knows exactly how to extract structured data from cooking sites. Each is a pocket expert, loaded only when needed.

### Autonomous Self-Extension

Legacy has a `skill-forge` skill that can:
1. Research a new API's documentation
2. Identify which endpoints are needed
3. Generate a complete `index.js` + `SKILL.md`
4. Install it into the workspace

**Legacy can extend its own capabilities at runtime** without human coding.

### Separation of Concerns

- **Skill author**: Knows the domain, writes the skill once
- **Legacy**: Consumes the skill, executes tasks
- **You**: Give Legacy access to the skill

They don't need to understand each other's systems.

### Cross-Runtime Portability

Because the base `SKILL.md` format is an open standard, the same skill instructions can run on:
- OpenClaw (with `index.js` for execution)
- Claude Code (with `scripts/` for execution)
- Any other Agent Skills-compatible runtime

---

## 8. Writing Effective Instructions

### The Description is Everything

The description field in frontmatter is what the agent reads to decide "should I use this skill?" Write it to match exactly how users phrase requests.

```yaml
# Weak — too vague
description: Manages tasks.

# Strong — specific, includes trigger phrases
description: >
  Manages Todoist tasks, projects, and sections for the OneMind infrastructure.
  Use when the user says tasks, todos, Todoist, to-do list, project management,
  or when creating, listing, completing, or prioritizing tasks. Supports
  creating tasks with due dates, priorities, and labels.
```

### Instructions Should Be Procedures, Not Descriptions

Don't describe what the tools do — tell Legacy how to use them step by step.

```markdown
## Bad: Description
The `ha_control` tool controls Home Assistant entities. It accepts an entity_id
and an action parameter.

## Good: Procedure
### Turning Something On
1. Get the entity_id (e.g., "light.office_desk", "switch.garage_fan")
2. Call: `ha_control({ action: "turn_on", entity_id: "entity_id_here" })`
3. If the user asked for a specific brightness, add: `brightness: 0-255`
4. Verify with `ha_state({ entity_id: "..." })` if the user wants confirmation
```

### Always Include Examples

Concrete examples are worth more than paragraphs of explanation. Show the exact function call pattern:

```markdown
## Examples

### Turn on office lights at 50%
\`\`\`
ha_control({ action: "turn_on", entity_id: "light.office_desk", brightness: 128 })
\`\`\`

### Get home status
\`\`\`
ha_status()
\`\`\`

### Activate movie scene
\`\`\`
ha_scene({ scene: "movie_time" })
\`\`\`
```

### Document Edge Cases Explicitly

What would the model likely get wrong without guidance?

```markdown
## Edge Cases

- **Entity not found**: If `ha_state` returns null, the entity_id may be wrong.
  Call `ha_list({ domain: "light" })` to see all valid entity IDs first.
- **Brightness range**: 0–255, not 0–100. 50% = 128, not 50.
- **Scenes vs automations**: Use `ha_scene` for named scenes, `ha_automation`
  for triggering automation. They are different.
```

### Progressive Disclosure Within Your Skill

Keep the main `SKILL.md` body fast and scannable. Move deep reference to `references/`:

```markdown
## Quick Reference
See [REFERENCE.md](references/REFERENCE.md) for the full entity ID list.
```

---

## 9. Skill Design Patterns

### Pattern 1: Single-Service Wrapper

The most common pattern. One skill = one external service.

```
github-api-1.0.3/
├── SKILL.md          ← When/how to use GitHub via this skill
├── index.js          ← Exports: create_pr, list_issues, merge_pr, etc.
└── package.json      ← @octokit/rest dependency
```

### Pattern 2: Multi-Tool Domain Skill

One skill covering a coherent domain with multiple interaction modes.

```
home-assistant-1.0.0/
├── SKILL.md          ← When to use each tool
└── index.js          ← Exports: ha_control, ha_status, ha_state, ha_scene, ha_list
```

### Pattern 3: Infrastructure Bridge

Skills that connect Legacy to other runtimes or message buses.

```
nexus-bridge/
├── SKILL.md          ← NATS subject conventions, dispatch patterns
└── index.js          ← Exports: nexus_dispatch, nexus_query, nexus_publish
```

### Pattern 4: Autonomous Tool Builder (Meta-Skill)

A skill that helps build other skills.

```
skill-forge-1.0.0/
├── SKILL.md          ← How to research APIs and generate skills
└── index.js          ← Exports: skill_forge_create, skill_forge_validate
```

### Pattern 5: Sensor / Perception Skill

Skills that pull data in (versus pushing commands out).

```
legacy-pulse/
├── SKILL.md          ← When to check system health metrics
└── index.js          ← Exports: get_system_vitals, get_agent_status
```

### Pattern 6: Script-Heavy Skill

Skills where the core logic is in executable scripts, and `index.js` just shells out.

```
data-processor/
├── SKILL.md          ← Workflow overview
├── index.js          ← Exports: process_file (shells to scripts/process.py)
├── scripts/
│   └── process.py    ← Heavy processing logic in Python
└── package.json
```

---

## 10. Secrets and Environment

### The Pattern

OpenClaw's gateway is started with `op run` on the VPS, which injects 1Password secrets as environment variables before the process starts. Skills receive them via `process.env`.

```javascript
// In index.js — read from env
const API_KEY = process.env.MY_SERVICE_API_KEY;
const BASE_URL = process.env.MY_SERVICE_URL || "https://api.example.com";
```

### Adding a New Secret

```bash
# 1. Create the item in 1Password
op item create \
  --category=login \
  --vault="00-24 SO (Sovereign Operations)" \
  --title="My Service" \
  "api-key[password]=the-actual-value"

# 2. Add to /opt/onemind/.env.tpl on VPS
# MY_SERVICE_API_KEY=op://aaeyris5z44pdqspk277u2oppi/My Service/api-key

# 3. Restart openclaw to pick up the new secret
systemctl restart openclaw
```

### Never Do This

```javascript
// ❌ Never hardcode
const API_KEY = "sk-abc123-real-secret";

// ❌ Never commit real values
// MY_SERVICE_API_KEY=sk-abc123-real-secret
```

### .env for Local Dev

When running locally (`openclaw gateway`), create a local `.env` file (gitignored) or export vars in your shell:

```bash
export MY_SERVICE_API_KEY="your-local-test-key"
openclaw gateway
```

---

## 11. Testing and Debugging

### Test a Tool Function Directly

```javascript
// test.js — run with: node test.js
require("dotenv").config();
const { my_tool } = require("./index.js");

(async () => {
  try {
    const result = await my_tool({ param1: "test", param2: "value" });
    console.log("Result:", JSON.stringify(result, null, 2));
  } catch (err) {
    console.error("Error:", err.message);
  }
})();
```

### Check Skill Loaded by Gateway

After restarting, check journal:

```bash
ssh root@104.248.4.212 "journalctl -u openclaw -f --no-pager | grep -i 'skill\|error'"
```

### Common Failures

| Symptom | Likely Cause | Fix |
|---|---|---|
| Tool not available after restart | `index.js` syntax error | `node index.js` locally to check |
| Tool throws "not configured" | Env var missing | Check env.tpl and restart |
| Tool hangs forever | No timeout on fetch | Add `AbortSignal.timeout(30000)` |
| Wrong tool called | Weak `description` field | Strengthen trigger keywords in frontmatter |
| Tool not activated at all | Description doesn't match phrases | Add synonyms to description |

### Validate Before Deploy

```bash
# Check syntax
node -c index.js

# Check exports
node -e "const s = require('./index.js'); console.log(Object.keys(s))"

# Check SKILL.md frontmatter (if skills-ref installed)
npx skills-ref validate .
```

---

## 12. Checklist Before Shipping

**SKILL.md**
- [ ] `name` matches directory name exactly
- [ ] `description` is specific, includes trigger keywords, ≤1024 chars
- [ ] Body has "When to Use" section
- [ ] Body has tool usage examples with exact call signatures
- [ ] Edge cases documented
- [ ] Body is under 500 lines

**index.js**
- [ ] All exported functions are async
- [ ] All required parameters validated with clear error messages
- [ ] All `fetch` calls have timeout (`AbortSignal.timeout`)
- [ ] No hardcoded secrets
- [ ] `module.exports` matches tools referenced in SKILL.md

**package.json**
- [ ] `name` matches directory/skill name
- [ ] All dependencies listed
- [ ] `npm install` runs without errors

**Secrets**
- [ ] All env vars documented in SKILL.md frontmatter or body
- [ ] Added to `/opt/onemind/.env.tpl` on VPS
- [ ] 1Password item created

**Deploy**
- [ ] `npm install` run in skill directory
- [ ] Gateway restarted
- [ ] Journal checked for errors
- [ ] Tool called once manually to verify it works
