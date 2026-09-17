# System — Pinned Context

Files in this directory are **always loaded** into the agent's context window at the start of every interaction. This is the "hot zone" — the most essential files the agent needs regardless of task.

> Inspired by Letta Context Repositories: the `system/` directory designates which files are always fully loaded into the system prompt.

---

## Rules

1. **Keep it small** — Total size should stay under 2,000 tokens. Every token here is consumed on every interaction.
2. **Only facts, not procedures** — Put "how to" instructions in `_codex/skills/`. System holds "what is" context.
3. **Agent-managed** — Agents can move files in and out of `system/` as they learn what's most relevant. This is the primary mechanism for agent self-organization.
4. **Human override** — Owner can lock files in `system/` by adding `pinned: true` to frontmatter.

## What Belongs Here

- Active goals and priorities (the 3-5 things that matter right now)
- Identity context (who the owner is, core preferences)
- Critical reference data that's needed in >50% of interactions

## What Does NOT Belong Here

- Procedures and SOPs (use `_codex/skills/`)
- Historical logs (use domain folders)
- Large reference documents (use L1/L2 tiered loading)
- Anything over 500 tokens (summarize it instead)

## How Agents Use This

At the start of every interaction:
1. Load ALL files in `_codex/system/` into context
2. Use this context to inform responses without needing to search
3. During reflection, evaluate: "Should anything be promoted to or demoted from `system/`?"

## Self-Organization

Agents are encouraged to manage `system/` contents:
- **Promote**: If a file is accessed in >50% of recent sessions → consider moving/symlinking to `system/`
- **Demote**: If a `system/` file hasn't been relevant in 7+ days → move it back to its domain folder
- **Summarize**: If a needed file is too large for `system/`, create a summary version here and keep the full version in its domain
