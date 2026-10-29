---
name: hermes-plugin-pattern
description: Build Hermes plugins. Plugin types, hooks, register pattern.
version: 1.0.0
tags: [hermes, plugin, architecture, onemind]
platforms: [linux, macos]
---

# Hermes Plugin Development — OneMind Pattern

The universal recipe for making services first-class inside Hermes. Every plugin: **connect → inject → tools**.

## The Three-Layer Pattern

Every plugin: **connect → inject → tools**

1. **Connect** — `register(ctx)`: set up data source, start subscriptions, open connections
2. **Inject** — `pre_llm_call(context)` hook: deep-copy state into prompt before LLM sees it
3. **Tools** — `register_tool(name, schema, handler)`: callable functions the LLM invokes

## Five Plugin Types (from Hermes docs)

| Type | Plugin kind enum | What it does | OneMind examples |
|------|-----------------|--------------|-----------------|
| **General (standalone)** | `standalone` | Tools + hooks | MAVEN (5 tools), nats-context (9 tools) |
| **General (backend)** | `backend` | Provides services to other plugins | nats-platform (bus transport) |
| **Gateway (platform)** | `platform` | Interactive channel (TAK, Telegram, Discord) | tak-platform (GeoChat bridge) |
| **LLM provider (exclusive)** | `model-provider` | Replace LLM with different backend | (choose own syntax) |
| **Memory (exclusive)** | `memory` | Replace built-in memory provider | (none yet) |

## The Standard Plugin Structure

```
plugin-name/
├── plugin.yaml          # manifest — name, version, description, requires_env
├── __init__.py           # register(ctx) — wires everything together
├── state.py              # optional — thread-safe state manager
├── subscriptions.py      # optional — NATS/subscription manager
├── formatter.py          # optional — format context for injection
└── tools.py              # optional — tool handler functions
```

**plugin.yaml template:**
```yaml
name: plugin-name
version: "0.1.0"
description: "What this plugin does"
author: OneMind OS / Zeus
requires_env: []          # needed envvars, checked at registration
```

## The register() Pattern

```python
def register(ctx):
    # 1. CONNECT
    state = StateManager()
    manager = SubscriptionManager(state, url)
    manager.start()

    # 2. INJECT — pre_llm_call hook
    def _inject(payload: dict) -> dict | str:
        if state is None:
            return ""
        snapshot = state.snapshot()
        return {"context": format_context(snapshot)}
    ctx.register_hook("pre_llm_call", inject)

    # 3. TOOLS — register callable tools
    ctx.register_tool(
        name="tool_name",
        schema={
            "name": "tool_name",
            "description": "What it does",
            "parameters": {
                "type": "object",
                "properties": {
                    "param1": {"type": "string", "description": "..."},
                    "param2": {"type": "number", "default": 30},
                },
                "required": ["param1"],
            },
        },
        handler=handle_tool,
    )
```

## The Hook Catalog (from VALID_HOOKS)

From `hermes_cli/plugins.py`:

| Category | Hooks |
|----------|-------|
| **Directive/control** | `pre_tool_call`, `pre_llm_call`, `pre_verify`, `pre_gateway_dispatch` |
| **Transform** | `transform_tool_result`, `transform_terminal_output`, `transform_llm_output`, `pre_transcription` |
| **Observer** | `post_tool_call`, `post_llm_call`, `pre_api_request`, `post_api_request`, `api_request_error`, `on_stream_start`, `on_stream_delta`, `on_stream_end`, `on_interim_message`, `on_session_start`, `on_session_end`, `on_session_finalize`, `on_session_reset`, `agent_loop_stopped`, `on_skill_lifecycle`, `subagent_start`, `subagent_stop`, `pre_approval_request`, `post_approval_response`, `pre_command`, `kanban_task_claimed`, `kanban_task_completed`, `kanban_task_blocked` |

## Available Hooks for Context Injection

```python
def _inject(payload: dict) -> dict | str:
    # Your state snapshot → format → return dict
    snapshot = state.snapshot()
    return {"context": format_context(snapshot)}

ctx.register_hook("pre_llm_call", _inject)
```

Takes: `context` (formatted string/dict). Filter for LLM-visible context.

## Tool Handler Pattern

```python
def handle_tool(params, **kwargs):
    del kwargs
    # params = LLM-provided arguments
    result = do_something(params)
    return json.dumps({"success": True, "result": result})
```

## Conventions

- **Registration once** — `register(ctx)` calls exactly once per plugin
- **Deep-copy context** — `snapshot.copy()` before LLM sees it
- **Graceful check_fn** — keep graceful under service failures
- **Thread safety** — background threads on shared state
- **Priority-based injection** — weight context relevance (TAK ops > system health)

## CLI Diagnostic

```bash
# Check plugins surfaces after registration
hermes plugins capabilities plugin-name
# List registered capabilities and their friendly names
```

## Hermes Discovery Hierarchy

Later sources override earlier ones:

1. **bundled** — `<repo>/plugins/` built-in
2. **project** — `.hermes/plugins/` project-specific
3. **pip** — `hermes_agent.plugins` entry points
4. **nix** — NixOS declarative installs

For OneMind: all plugins in `~/.hermes/plugins/`, enabled via `plugins.enabled` in config.yaml.

## OneMind Plugin Examples

### Context-Only (nats-context)
```python
# register once → start subscriptions → hook → tools
_state = LiveState()
manager = SubscriptionManager(_state, url)
manager.start()

ctx.register_hook("pre_llm_call", context_inject)
for tool in TOOLS:
    ctx.register_tool(name=tool.name, schema=tool.schema, handler=tool.handler)
```

### Tools-Only (MAVEN)
```python
# register once → all tools
ctx.register_tool(
    name="tool_name",
    schema={...},
    handler=handle_tool,

    name="tool_name2",
    schema={...},
    handler=handle_tool2,
)
```

### Platform Gateway (tak-platform)
```python
# Gateway channels have special kind:
# ctx.register_hook("pre_gateway_dispatch", ...) fires on
```

## Distribution and Installation

```bash
# Install Git-based plugin
hermes plugins install <github-url> --ref 0123456789abcdef
# Install local file
hermes plugins install file:///path/to/plugin
# Enable/disable
hermes plugins enable <name>
hermes plugins disable <name>
```

## Build Checklist For OneMind Plugin

```
[ ] plugin.yaml — name, version, description, requires_env + optional_env
[ ] __init__.py — register(ctx) calls register_tool/hook + subscriptions
[ ] pre_llm_call — context injection hook (if context plugin)
[ ] tools — each with schema + handler + description
[ ] descriptions — LLM-friendly descriptions
[ ] requires_env — env vars validated at startup
[ ] graceful check_fn — keeps graceful under service failures
[ ] threading — no shared mutable state without locks
[ ] priority — context weighting — TAK > system health > integrations
[ ] test — hermes plugins capabilities plugin-name
```