# OpenClaw Plugin Templates

A reference kit for building OpenClaw gateway plugins. Plugins extend the **gateway itself** — adding channels, endpoints, context engines, prompt hooks, slash commands, and RPC methods. This is different from skills, which give agents tools.

> **Skills** = agent-level tools (Node.js functions Legacy can call)
> **Plugins** = gateway-level extensions (modify how the gateway works)

---

## What This Folder Contains

```
plugin-templates/
├── README.md                       ← You are here
├── templates/
│   ├── basic-plugin/               ← Minimal starting point
│   ├── channel-plugin/             ← Custom messaging channel
│   ├── context-engine-plugin/      ← Custom memory/context engine
│   └── webhook-plugin/             ← HTTP endpoint + webhook receiver
└── examples/
    └── prompt-injector/            ← Working example: inject context into every prompt
```

---

## Quick Start — Build a Plugin in 4 Steps

### Step 1: Copy a template

```bash
cp -r plugin-templates/templates/basic-plugin ./my-plugin
cd my-plugin
```

### Step 2: Edit `package.json`

Update the name and entry point:

```json
{
  "name": "my-plugin",
  "version": "1.0.0",
  "openclaw": {
    "extensions": ["./src/index.ts"]
  }
}
```

### Step 3: Write your plugin logic in `src/index.ts`

```typescript
export default function register(api) {
  // Your gateway extensions here
}
```

### Step 4: Install and test

```bash
# Link locally for development (symlink, edits apply instantly)
openclaw plugins install ./my-plugin --link

# Or install from npm
openclaw plugins install @my-org/my-plugin

# Restart gateway to load
systemctl restart openclaw
```

---

## Plugin API Reference

Every plugin exports a `register(api)` function. The `api` object provides:

### Core Methods

| Method | Purpose | Example Use |
|--------|---------|-------------|
| `api.on(event, handler, opts)` | Hook into gateway lifecycle events | Inject context before every prompt |
| `api.registerGatewayMethod(name, handler)` | Add custom RPC methods | `myplugin.status` callable via gateway |
| `api.registerHttpRoute(opts)` | Expose HTTP endpoints on gateway | Receive webhooks from external services |
| `api.registerHook(event, handler, opts)` | Register event hooks | Run logic on `command:new`, `gateway:startup` |
| `api.registerCommand(opts)` | Custom slash commands (no AI turn) | `/mystatus` for quick checks |
| `api.registerChannel(opts)` | Add messaging channel | Custom chat platform integration |
| `api.registerContextEngine(id, factory)` | Custom context/memory engine | Like Lossless Claw |

### Lifecycle Events

| Event | When It Fires |
|-------|---------------|
| `gateway:startup` | After all channels started and hooks loaded |
| `before_prompt_build` | Before every prompt is assembled — inject/modify system context |
| `command:new` | When `/new` command is invoked |
| `command` | Any slash command |
| `agent:bootstrap` | Agent session initialization |

### Prompt Hook Return Values

```typescript
api.on("before_prompt_build", (event, ctx) => {
  return {
    prependSystemContext: "Added before system prompt.",
    appendSystemContext: "Added after system prompt.",
  };
}, { priority: 10 });
```

Priority: lower number = runs first.

---

## Plugin vs Skill Decision Guide

| Question | If Yes → |
|----------|----------|
| Does it add a tool the agent calls? | **Skill** |
| Does it modify how the gateway processes messages? | **Plugin** |
| Does it add a new communication channel? | **Plugin** |
| Does it expose an HTTP endpoint? | **Plugin** |
| Does it change how context/memory works? | **Plugin** |
| Does it add a slash command? | **Plugin** |
| Does an agent need to decide when to use it? | **Skill** |
| Should it run automatically on every message? | **Plugin** |

---

## Configuration

Enable plugins in `openclaw.json`:

```json5
{
  "plugins": {
    "enabled": true,
    "allow": ["my-plugin"],          // allowlist (optional)
    "deny": ["untrusted-plugin"],    // denylist (optional)
    "load": {
      "paths": ["~/my-plugin"]       // load from local path
    },
    "entries": {
      "my-plugin": {
        "enabled": true,
        "config": {
          "apiKey": "${MY_PLUGIN_API_KEY}"
        }
      }
    }
  }
}
```

---

## Currently Installed Plugins (VPS)

| Plugin | Purpose |
|--------|---------|
| `lobster` | Named automation routines |
| `bluebubbles` | iMessage via BlueBubbles |
| `diffs` | File diff tracking |
| `llm-task` | Background LLM task execution |
| `open-prose` | Prose/writing enhancements |
| `diagnostics-otel` | OpenTelemetry diagnostics |
| `acpx` | Agent Communication Protocol extensions |

---

## Security Notes

- Treat plugin installs like running arbitrary code — they execute in the gateway process
- Use `--link` for local development, `--pin` for npm installs to lock versions
- Secrets should use `${ENV_VAR}` references in config, never hardcoded values
- HTTP routes should use `auth: "plugin"` to require gateway authentication
