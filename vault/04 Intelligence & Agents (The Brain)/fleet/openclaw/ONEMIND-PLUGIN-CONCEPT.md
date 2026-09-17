# OneMind OpenClaw Plugin — Concept

> What a `@onemind/openclaw-plugin` would look like using `@openclaw/plugin-sdk`
> Date: 2026-05-17

---

## Yes — `@openclaw/plugin-sdk` Is for Creating Your Own Plugins

The SDK gives you:

| Capability | What It Does |
|------------|--------------|
| `definePlugin()` | Main entry point — registers your plugin with OpenClaw |
| Hook registration | Tap into `pre_tool_call`, `post_tool_call`, `llm_output`, `model_call_*` |
| Tool contracts | Declare custom tools your plugin provides to the agent |
| Channel adapters | Build your own messaging platform connector |
| MCP tool handlers | `createPluginToolsMcpHandlers()` — expose plugin tools via MCP |
| Sidecars | Start background services that run alongside the gateway |
| ClawHub distribution | Publish as npm package, list on ClawHub marketplace |

---

## What Would a OneMind Plugin Be?

### Purpose

A unified plugin that turns OpenClaw into the OneMind orchestration layer:

```
@onemind/openclaw-plugin
├── NATS bridge (pub/sub to OneMind event bus)
├── DimOS connector (robotics control tools)
├── OM1 connector (embodied agent tools)
├── OneMind identity (system prompt injection, personality)
├── Cross-agent routing (forward tasks to Hermes when appropriate)
└── OneMind dashboard tab (custom Control UI panel)
```

---

## Plugin Structure

```
@onemind/openclaw-plugin/
├── package.json
├── src/
│   ├── index.ts              # definePlugin() entry
│   ├── hooks/
│   │   ├── nats-bridge.ts    # NATS subscribe/publish on lifecycle
│   │   ├── identity.ts       # Inject OneMind system context
│   │   └── routing.ts        # Route tasks to other agents
│   ├── tools/
│   │   ├── nats-publish.ts   # Agent tool: publish to NATS
│   │   ├── nats-request.ts   # Agent tool: request/reply on NATS
│   │   ├── dimos-command.ts  # Agent tool: send command to DimOS
│   │   ├── om1-action.ts     # Agent tool: trigger OM1 action
│   │   └── hermes-delegate.ts# Agent tool: delegate task to Hermes
│   ├── channels/
│   │   └── nats-channel.ts   # NATS as a "messaging platform" adapter
│   └── ui/
│       └── dashboard.tsx     # Custom Control UI panel (optional)
├── contracts.json            # Tool schemas for ClawHub
└── README.md
```

---

## Core Implementation

```typescript
// src/index.ts
import { definePlugin } from '@openclaw/plugin-sdk';
import { createNatsBridge } from './hooks/nats-bridge';
import { injectIdentity } from './hooks/identity';
import { routeToHermes } from './hooks/routing';
import { natsPublishTool, natsRequestTool } from './tools/nats-publish';
import { dimosCommandTool } from './tools/dimos-command';
import { om1ActionTool } from './tools/om1-action';
import { hermesDelegateTool } from './tools/hermes-delegate';

export default definePlugin({
  name: '@onemind/openclaw-plugin',
  version: '0.1.0',
  description: 'OneMind orchestration layer for OpenClaw',

  hooks: {
    // Start NATS connection when gateway boots
    async on_gateway_ready(ctx) {
      await createNatsBridge(ctx);
    },

    // Inject OneMind identity into every agent session
    async on_session_start(ctx) {
      await injectIdentity(ctx);
    },

    // Before each tool call — check if it should route to Hermes
    async pre_tool_call(ctx, toolCall) {
      return routeToHermes(ctx, toolCall);
    },

    // After each tool call — publish result to NATS
    async post_tool_call(ctx, result) {
      await ctx.nats?.publish(
        `onemind.tools.${result.toolName}.result`,
        JSON.stringify(result.output)
      );
    },

    // Transform LLM output — add OneMind formatting/routing
    async llm_output(ctx, output) {
      // Could route certain outputs to NATS subjects
      if (output.includes('[broadcast]')) {
        await ctx.nats?.publish('onemind.broadcast', output);
      }
      return output;
    }
  },

  // Tools the agent can invoke
  contracts: {
    tools: [
      natsPublishTool,
      natsRequestTool,
      dimosCommandTool,
      om1ActionTool,
      hermesDelegateTool,
    ]
  },

  // Plugin config schema (user sets in openclaw config)
  config: {
    nats_url: { type: 'string', default: 'nats://localhost:4222' },
    nats_subjects: { type: 'array', default: ['onemind.events.>'] },
    dimos_endpoint: { type: 'string', default: 'http://localhost:8080' },
    om1_endpoint: { type: 'string', default: 'http://localhost:9090' },
    hermes_endpoint: { type: 'string', default: 'http://localhost:4200' },
    identity_prompt: { type: 'string', default: 'You are OneMind...' },
  }
});
```

---

## Example Tool: NATS Publish

```typescript
// src/tools/nats-publish.ts
export const natsPublishTool = {
  name: 'nats_publish',
  description: 'Publish a message to a NATS subject on the OneMind event bus',
  parameters: {
    type: 'object',
    properties: {
      subject: { type: 'string', description: 'NATS subject (e.g. onemind.robotics.go2.walk)' },
      payload: { type: 'string', description: 'JSON payload to publish' },
    },
    required: ['subject', 'payload']
  },
  async execute(ctx, { subject, payload }) {
    await ctx.nats.publish(subject, payload);
    return { success: true, subject, timestamp: Date.now() };
  }
};
```

---

## Example Tool: Delegate to Hermes

```typescript
// src/tools/hermes-delegate.ts
export const hermesDelegateTool = {
  name: 'hermes_delegate',
  description: 'Delegate a task to Hermes Agent (Python-native tasks, long-running ops)',
  parameters: {
    type: 'object',
    properties: {
      task: { type: 'string', description: 'Task description for Hermes' },
      wait: { type: 'boolean', description: 'Wait for completion?', default: false },
    },
    required: ['task']
  },
  async execute(ctx, { task, wait }) {
    const res = await fetch(`${ctx.config.hermes_endpoint}/api/task`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ task, wait })
    });
    return await res.json();
  }
};
```

---

## Installation & Usage

```bash
# Development
cd @onemind/openclaw-plugin
npm install
npm run build

# Install into OpenClaw
openclaw plugins install ./dist
# or from npm registry:
openclaw plugins install @onemind/openclaw-plugin

# Configure in OpenClaw config
# ~/.openclaw/config.yaml (or equivalent)
plugins:
  @onemind/openclaw-plugin:
    nats_url: "nats://localhost:4222"
    nats_subjects: ["onemind.>"]
    dimos_endpoint: "http://192.168.1.100:8080"
    hermes_endpoint: "http://localhost:4200"
```

---

## Plugin vs Hook vs MCP Server — When to Use What

| Need | Solution |
|------|----------|
| Persistent NATS subscription | **Plugin** (long-lived, starts with gateway) |
| React to a tool call once | **Hook** (stateless, per-event) |
| Share tools with multiple agents | **MCP Server** (platform-agnostic) |
| Custom agent tools | **Plugin** (tool contracts) |
| Cross-agent communication | **Plugin** (NATS bridge) OR **MCP Server** |
| Shell script integration | **Hook** (Hermes shell hooks) |
| Distribution to community | **Plugin** on ClawHub (npm package) |

---

## Next Steps

1. Scaffold with `npm init` + `@openclaw/plugin-sdk` as dev dependency
2. Start with NATS bridge hook only (MVP)
3. Add tools incrementally (nats_publish → dimos_command → hermes_delegate)
4. Test locally with `openclaw plugins install ./`
5. Publish to ClawHub when stable
