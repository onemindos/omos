---
name: basic-skill
description: >
  [REPLACE THIS] Describe exactly what this skill does and when Legacy should
  activate it. Be specific — include the phrases a user would say that should
  trigger this skill. Example: "Use when the user asks about X, Y, or Z."
metadata:
  author: zeus
  version: "1.0.0"
---

# Basic Skill

## When to Use

[Describe the specific scenarios when Legacy should activate this skill. Be
concrete — list the user intents, phrases, or task types that should trigger it.]

- User asks about [X]
- User wants to [Y]
- User says "[example phrase]"

## Available Tools

| Tool           | Description                              |
| -------------- | ---------------------------------------- |
| `basic_action` | Does the primary thing this skill is for |

## How to Use

### Perform the main action

```
basic_action({ input: "the value you want to process" })
```

Expected response:

```json
{
  "status": "ok",
  "result": "..."
}
```

## Notes

- [Add any edge cases, caveats, or things Legacy should know]
- [Environment variables required: `MY_SERVICE_API_KEY`]


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
