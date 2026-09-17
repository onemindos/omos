# Technology Radar

> *Track technologies on your radar: adopt, trial, assess, or hold.*
> *Inspired by ThoughtWorks Technology Radar. Updated quarterly.*
> *Lives in: `SO/06 Research/technology-radar.md`*

---

```yaml
---
type: technology-radar
domain: 06
agent: Trinity
owner: Zeus
status: active
updated: YYYY-MM-DD
---
```

---

## What This Is

The technology radar maps everything in the tech landscape that matters to One Mind OS. Each technology has a quadrant (what kind of thing it is) and a ring (what you should do with it).

**Quadrants:**
- **Tools** — Software tools, IDEs, CLIs, utilities
- **Platforms** — Clouds, infrastructure platforms, operating environments
- **Languages & Frameworks** — Programming languages, frameworks, libraries
- **Techniques** — Patterns, methods, approaches

**Rings:**
- **Adopt** — Use in production. High confidence. Default choice.
- **Trial** — Try in projects. Promising. Building confidence.
- **Assess** — Research and evaluate. On the horizon.
- **Hold** — Don't adopt. Use caution with existing usage.

---

## Adopt — Use These Now

These are proven, in active use, high confidence.

### Tools

| Technology | Version/Status | Why | Since |
|-----------|----------------|-----|-------|
| Obsidian | Latest | PKM + vault interface | [DATE] |
| VS Code | Latest | Code + vault editing | [DATE] |
| 1Password | Latest | Credential management | [DATE] |
| Tailscale | Latest | Private network overlay | [DATE] |
| OpenClaw | Latest | AI agent runtime | [DATE] |

### Platforms

| Technology | Why | Since |
|-----------|-----|-------|
| DigitalOcean | VPS hosting — simple, reliable | [DATE] |
| GitHub | Source of truth for vault + code | [DATE] |
| Discord | Agent communication channel | [DATE] |

### Languages & Frameworks

| Technology | Use Case | Since |
|-----------|---------|-------|
| Markdown | Everything — vault, docs, templates | Always |
| Bash | VPS automation and scripting | [DATE] |
| Python | Agent skills, data processing | [DATE] |

### Techniques

| Technique | Domain | Since |
|-----------|--------|-------|
| Git-based PKM | Vault architecture | [DATE] |
| Agent-per-domain architecture | SO 01 | [DATE] |
| Cron-driven intelligence | SO 07 | [DATE] |

---

## Trial — Actively Evaluating

Promising — using in specific contexts to build confidence.

| Technology | Quadrant | Testing For | Started | Notes |
|-----------|---------|-------------|---------|-------|
| [FILL] | | | | |

---

## Assess — On the Radar

Worth watching. Not ready to trial yet.

| Technology | Quadrant | Why Interesting | Added | Signals |
|-----------|---------|----------------|-------|---------|
| [FILL] | | | | |

---

## Hold — Avoid or Use Cautiously

Technologies to avoid adopting, or to reduce dependency on.

| Technology | Reason | Alternative | Since |
|-----------|--------|------------|-------|
| Notion | Data lock-in, no git, SaaS dependency | Obsidian + git vault | [DATE] |
| [FILL] | | | |

---

## Quarterly Review Process

Ask Trinity to run radar review:

> *"Trinity, run the technology radar quarterly review. For each item: (1) should it move rings based on new information? (2) any new technologies to add? (3) any items to remove? Present changes with reasoning."*

Changes from last quarter: [FILL after review]

---

## Notes on Evaluation Criteria

**To move to Adopt:** Used in production for 3+ months. Confident it solves the problem better than alternatives. Team/system understanding is solid.

**To move to Trial:** Clear use case identified. Low risk to experiment. Willing to invest evaluation time.

**To move to Assess:** Something interesting happened — new version, new use case, new evidence. Worth tracking.

**To move to Hold:** Reliability concerns, security issues, vendor lock-in risk, or better alternative exists.

---

*Agent Trinity: Maintain this radar. Propose quarterly updates based on your research activity and new evidence. Flag any 'Adopt' item that should move to 'Hold' immediately.*
