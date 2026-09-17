---
title: "Hook"
type: reference
domain: SO
status: active
created: 2026-05-27
updated: 2026-05-27
tags: [concept, ai-agents, architecture, execution, triggers, hermes]
---

# Hook

## Definition
An intercept point in an agent's execution loop. A hook fires at a specific moment — before a tool call, after a response, on error, on memory write — letting you inject logic without modifying the core agent loop.

## Hook Types (AI Agent Context)

| Hook | Fires When | Use Case |
|------|-----------|----------|
| **pre-tool** | Before a tool executes | Validate, log, block, modify args |
| **post-tool** | After tool returns | Audit output, transform, route |
| **pre-response** | Before final answer is sent | Filter, redact, format |
| **post-response** | After response delivered | Log, trigger downstream actions |
| **memory-read** | Agent reads from memory | Inject context, enrich |
| **memory-write** | Agent writes to memory | Validate, deduplicate |
| **on-error** | Tool or agent fails | Alert, fallback, recover |

## In Hermes
Hooks are configured in `~/.hermes/config.yaml` under the `hooks:` key.
`hooks_auto_accept: false` — hooks require approval by default.

## OneMind OS Angle
Hooks are how **OM Cortex** would intercept and influence agent behavior without owning it:
- Pre-tool hook → Cortex validates tool call intent against current operating mode
- Post-response hook → Cortex routes output to the right pillar (SO/HP/LE/GE)
- Memory hook → Cortex ensures vault writes land in the correct domain

Combined with [[chain-of-thought]]: hooks can observe AND intercept the reasoning step before action fires.

## Related Concepts
- [[chain-of-thought]] — CoT + hooks = observe reasoning, intercept before execution
- [[om-cortex]] — Cortex implemented as a hook layer over Legacy
- [[interdependent]] — hooks create explicit dependency points between components

## Notes
<!-- Add Hermes hook config examples, NATS-based hook patterns here -->
