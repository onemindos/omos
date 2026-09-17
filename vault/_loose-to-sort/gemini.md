# gemini.md

Context for Gemini CLI when working in this repository.

---

## Your Role

Gemini CLI coding agent. You handle code tasks, analysis, and research within the OneMind ecosystem. Read `agents.md` for full vault conventions — this file provides Gemini-specific context.

---

## This Repository

- NOT a software project — it's an Obsidian vault (markdown knowledge base)
- Code lives exclusively in `_system/`
- 4 pillars: SO (Operations 00-24), HP (Personal 25-49), LE (Home 50-74), GE (Business 75-99)
- Do NOT modify markdown knowledge files unless explicitly asked

---

## Architecture

```
ONE MIND SYSTEM
├── Hermes (Brain) → Legacy AI — orchestration, scheduling, deep work (~/.hermes/)
├── Claude Code → coding, features, PRs (reads claude.md)
├── Gemini CLI ← YOU → coding, research, analysis (reads this file)
├── Copilot → IDE coding (reads .github/copilot-instructions.md)
└── Codex/OpenCode → coding (read agents.md)
```

---

## Working Rules

1. Read `agents.md` for full vault conventions
2. Code changes go in `_system/` only
3. Never rename pillar folders (SO/HP/LE/GE)
4. Preserve YAML frontmatter in all markdown files
5. No secrets in markdown — reference 1Password only
6. Use conventional commits
7. All agents share this folder — your changes are visible to all other agents

---

## Key Paths

| Path | Purpose |
|------|---------|
| `_system/` | All operational code (hooks, scripts, configs) |
| `agents.md` | Universal agent instructions (read this first) |
| `claude.md` | Claude Code context |
| `onemind.md` | THE canonical master document |
| `codex.md` | CODE× method deep spec |
| `.hermes.md` | Hermes agent context |
| `00-24 SO (Sovereign Operations)/00 Framework (The Constitution)/reference/SCHEMA.md` | Domain file schema |
