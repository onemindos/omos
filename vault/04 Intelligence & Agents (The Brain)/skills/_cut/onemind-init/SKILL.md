---
name: onemind-init
description: Bootstrap a new codex vault or onboard a new agent to an existing vault. Explores structure, loads context, and initializes memory using concurrent subagents.
---

# Codex Skill: Initialize

How to bootstrap a new codex vault or onboard a new agent to an existing vault. Inspired by Letta Context Repositories' memory initialization skill, which explores the codebase and reviews historical data using concurrent subagents.

---

## When to Use

- Setting up a brand new codex vault from scratch
- Onboarding a new agent to an existing vault (agent needs to learn what's here)
- Recovering from a corrupted or reset vault state
- Bootstrapping from external data (chat history, existing docs, exported notes)

## New Vault Initialization

### Step 1: Create the Skeleton

```
codex/
├── _codex/                  ← Copy from template (this folder)
├── 00-24 SO/README.md
├── 25-49 HP/README.md
├── 50-74 LE/README.md
├── 75-99 GE/README.md
├── assets/
└── CODEX-FRAMEWORK.md
```

- Create 4 domain group folders with README.md files
- Copy `_codex/` from the template repository
- Create empty `assets/` folder
- Initialize git: `git init && git add -A && git commit -m "init: vault skeleton"`

### Step 2: Seed Domain Folders

Create numbered domain subfolders based on the owner's needs. Not every vault needs all 100 domains. Start with what's relevant:

- Read the owner's context (conversation history, profile, stated goals)
- Create 5-10 domain folders that match their primary areas of focus
- Each domain folder gets a README.md explaining what belongs there

### Step 3: Import Existing Knowledge

If bootstrapping from existing data:

1. **Chat history** — Scan conversations for facts, decisions, preferences, and relationships. Extract each into the correct document type (note, decision, person, SOP).
2. **Existing documents** — Classify each file using the organize skill's decision tree, add frontmatter, move to the correct domain.
3. **Exported notes** (from Notion, Obsidian, Apple Notes, etc.) — Parse, add frontmatter, classify, and file.

For large imports, process in parallel batches:
- Group source files by likely domain
- Process each group independently
- Merge results and check for duplicates

### Step 4: Populate `system/`

Identify the 3-5 most essential files and place them in `_codex/system/`:
- Active goals or priorities
- Owner identity/preferences
- Any "always-needed" reference data

### Step 5: Initial Commit

```bash
git add -A
git commit -m "init: vault bootstrapped from [source description]"
```

## Agent Onboarding (Existing Vault)

When a new agent is given access to an existing vault:

### Step 1: Read `_codex/INTERFACE.md`

This is the universal entry point. It explains everything.

### Step 2: Load `_codex/system/`

Read all files in the system directory — this is the pinned context.

### Step 3: Orientation Scan

- Read each domain group's README.md (L0 scan)
- Build a mental map of what's in the vault
- Note which domains are populated vs. empty

### Step 4: Introduce Yourself

Create an agent profile note:

```yaml
---
title: "[Agent Name] — Agent Profile"
type: reference
status: active
created: YYYY-MM-DD
tags: [domain/so, type/reference, agent/{name}]
ai_summary: "Profile and capabilities of the {name} agent"
---
```

Include: what this agent does, what domains it owns, what tools it has access to.

### Step 5: Register in System

If the vault tracks active agents, add yourself to the registry (typically in `02 Agents/`).

---

## Parallel Bootstrapping (Advanced)

For large imports (1000+ files or extensive chat history), use the Letta-inspired fan-out pattern:

1. **Partition** the source data into independent chunks (by date range, topic, or source)
2. **Process each chunk** independently — extract facts, classify, create notes
3. **Merge** all results into the vault
4. **Dedup** — run `_codex/skills/maintain.md` duplicate detection across the merged results
5. **Commit** with a descriptive message: `init: imported [N] notes from [source]`

Git worktrees or branches can isolate parallel processing if agents support it.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
