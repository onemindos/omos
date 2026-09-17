# Skill Scaffold

> *Use this scaffold to create a new clawhub-compatible skill for OpenClaw.*
> *Skills extend what your agents can do — API integrations, CLI tools, data sources.*

---

## Skill Directory Structure

```
skills/
└── your-skill-name/
    ├── README.md       ← Human-readable description
    ├── INSTALL.md      ← Setup instructions
    ├── SKILL.md        ← The actual skill definition (what the agent reads)
    └── VERIFY.md       ← How to confirm the skill is working
```

---

## README.md Template

```markdown
# Skill Name

Brief description: what this skill does, what API or tool it wraps.

## What It Does

- Capability 1
- Capability 2
- Capability 3

## Prerequisites

- Tool X installed
- API key in 1Password at: op://[vault]/[item]/[field]

## Version

1.0.0
```

---

## INSTALL.md Template

```markdown
# Installation

## Prerequisites

Install required CLI tools:

```bash
# Example
brew install jq curl
```

## Configuration

Set credentials:

```bash
export SKILL_API_KEY=$(op read "op://vault-id/item-name/field")
```

Or store in `.env`:

```bash
echo "SKILL_API_KEY=your-key" >> ~/.openclaw/.env
```

## Verification

After installing, run the verify step to confirm setup.
```

---

## SKILL.md Template

This is what your agent reads. Write it as instructions to the agent.

```markdown
# Skill: [Skill Name]

## Purpose

[One sentence: what capability this skill adds to the agent]

## When to Use This Skill

Use this skill when:
- [Trigger condition 1]
- [Trigger condition 2]

Do NOT use this skill when:
- [Anti-pattern or out-of-scope situation]

## How to Use

### [Action 1 Name]

To [perform action 1]:

```bash
[exact command]
```

Expected output: [what success looks like]

### [Action 2 Name]

To [perform action 2]:

```bash
[exact command]
```

## Error Handling

| Error | Cause | Fix |
|-------|-------|-----|
| [Error message] | [Why it happens] | [What to do] |

## Examples

**Example: [Use case]**
Input: [what you give it]
Command: `[command]`
Output: [what comes back]

## Notes

[Any gotchas, rate limits, or important context]
```

---

## VERIFY.md Template

```markdown
# Verification

Run this to confirm the skill is working:

```bash
[test command]
```

Expected output:

```
[expected response]
```

If you see an error instead:
1. Check that prerequisites are installed: `which [tool]`
2. Check that credentials are set: `echo $SKILL_API_KEY`
3. Check network access: `curl -I [endpoint]`
```

---

## Installing Your Skill

Once the 4 files are created:

```bash
# If publishing to clawhub (public)
clawhub publish ./skills/your-skill-name

# If using locally only — just place in workspace
cp -r ./skills/your-skill-name ~/.openclaw/workspace/skills/
```

Tell your agent to use the skill by referencing SKILL.md in its workspace or via CODEX.md:

```
Skills available in this workspace:
- your-skill-name: see skills/your-skill-name/SKILL.md
```

---

*Part of the One Mind OS OpenClaw Deep Dive — Course 3*
