---
skill: entity-system
type: install
created: 2026-05-04
tags: [type/skill, zone/so, status/active]
---

# INSTALL: Entity System Skill

> Self-installation instructions for any agent loading this skill.

## What to Load

Read these files in order:

1. `SKILL.md` — core classification and routing instructions
2. `ONE-MIND-ENTITY-SYSTEM.md` (vault root) — full entity definitions and pillar map

## What to Remember

After reading, confirm you can answer:

- What are the 11 entity types?
- What are the 4 pillars and their domain ranges?
- What frontmatter fields are required for every entity?
- Where does an unrouted capture go?

## No Dependencies

This skill requires only:
- `read_file` — to read vault files
- `write_file` — to create new entity files
- `list_directory` — to navigate folder structure

No API keys. No database. No internet. Works on any filesystem.

## Activation Phrase

When a user says any of the following, activate this skill:
- "classify this"
- "where does this go"
- "route this to the right place"
- "what type is this"
- "create a [goal/project/task/etc]"
- "add this to the vault"


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
