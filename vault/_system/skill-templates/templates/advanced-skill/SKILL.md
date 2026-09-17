---
name: advanced-skill
description: >
  [REPLACE THIS] Full-featured skill with scripts, deep reference docs, and
  static assets. Use as a starting point for complex domain integrations that
  require executable scripts, detailed API references, or template assets.
  Activate when the user [describe specific triggers].
license: MIT
compatibility: Node.js 18+. Requires [list any system tools, e.g., python3, jq].
metadata:
  author: zeus
  version: "1.0.0"
  domain: "[domain name]"
---

# Advanced Skill

## When to Use

[Describe the exact conditions that should activate this skill.]

- User asks to [X]
- User wants to [Y]
- User mentions "[keyword]" in the context of [Z]

## Architecture

[Explain how this skill works — what systems it connects to, what the data flow
looks like. A few sentences or a diagram if helpful.]

```
User request
  → Legacy activates skill
    → index.js tool called
      → scripts/process.py or API call
        → result returned to Legacy
```

## Available Tools

| Tool               | Description                               |
| ------------------ | ----------------------------------------- |
| `advanced_process` | Primary processing tool                   |
| `advanced_analyze` | Analyzes data with optional script output |
| `advanced_report`  | Generates a structured report             |

## Quick Reference

### Process something

```
advanced_process({ input: "value", mode: "standard" })
```

### Run analysis

```
advanced_analyze({ target: "item-id", depth: "full" })
```

### Generate report

```
advanced_report({ format: "markdown", timeframe: "7d" })
```

## Detailed Reference

For the complete API reference, see [references/REFERENCE.md](references/REFERENCE.md).

For data schemas and validation rules, see [assets/schema.json](assets/schema.json).

## Scripts

The `scripts/` folder contains standalone scripts that can be run directly:

```bash
# Process a file
python3 scripts/process.py --input data.json --output result.json

# Validate a payload
python3 scripts/validate.py --schema assets/schema.json --data payload.json
```

## Common Patterns

### [Pattern 1 Name]

1. Call `advanced_process` with [specific params]
2. If result contains [X], then [do Y]
3. Optionally confirm with `advanced_report`

### [Pattern 2 Name]

1. Use `advanced_analyze` first to understand the target
2. Based on the analysis, call `advanced_process` with the correct `mode`

## Edge Cases

- **[Edge case 1]**: If [condition], do [action] instead
- **[Edge case 2]**: [Explain the non-obvious behavior]
- **Empty results**: Not an error — means no matching data

## Required Environment Variables

| Variable             | Description              |
| -------------------- | ------------------------ |
| `ADVANCED_SKILL_URL` | Base URL for the service |
| `ADVANCED_SKILL_KEY` | API key                  |
