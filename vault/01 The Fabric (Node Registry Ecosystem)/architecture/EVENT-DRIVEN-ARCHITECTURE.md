# Event-Driven Architecture — OpenClaw & Hermes

> Comparison of event-driven patterns in both platforms for OneMind integration.
> Date: 2026-05-17

---

## OpenClaw Event-Driven Gateway (TypeScript/Node.js)

### Internal Event Flow

```
Inbound events (Telegram poll, Discord WS, Webhook, CLI)
    → Channel adapter (plugin)
        → Gateway event bus (internal)
            → Agent turn lifecycle:
                model_call_start → tool_call → tool_result → llm_output → turn_complete
            → Plugin hooks fire at each lifecycle stage
        → Outbound delivery (streaming SSE, WS to Control UI)
```

### Key Event Mechanisms

| Mechanism | Purpose |
|-----------|---------|
| WebSocket protocol v4 | Real-time Control UI + SDK clients |
| SSE (Server-Sent Events) | Session history live streaming |
| Plugin hooks | `pre_tool_call`, `post_tool_call`, `llm_output`, `model_call_*` |
| Heartbeat agent events | Distinguish scheduled runs from chat-triggered runs |
| StatusReactionController | Emoji lifecycle across platforms (🧠→🛠️→💻→✅) |
| Restart trace logs | signal → drain → close → next-start → ready spans |

### Plugin Hook Points

```typescript
// Available plugin lifecycle hooks in OpenClaw
pre_tool_call    // Before any tool executes
post_tool_call   // After tool returns result
model_call_start // Before LLM API call
model_call_end   // After LLM response received
llm_output       // When streaming LLM output arrives
turn_complete    // Agent turn finished
```

---

## Hermes Event-Driven Agent Loop (Python)

### Internal Event Flow

```
Inbound (CLI/TUI, Gateway platforms, API server, ACP, Webhooks)
    → Agent loop (async Python)
        → Plugin hooks: pre_llm_call, post_llm_call, pre_tool_call, post_tool_call
        → transform_llm_output, transform_terminal_output
        → on_session_start, on_session_end
    → Background processes:
        → watch_patterns (pattern-match on stdout → notify agent)
        → notify_on_complete (auto-notify when bg process ends)
        → Self-improvement fork (background memory/skill review)
        → Curator (7-day cron, grades/prunes/consolidates skills)
    → Webhook direct-delivery (zero-LLM push, no agent needed)
```

### Key Event Mechanisms

| Mechanism | Purpose |
|-----------|---------|
| Plugin hooks (decorators) | pre/post_llm_call, pre/post_tool_call, on_session_start/end |
| transform_llm_output | Reshape LLM output before conversation delivery |
| transform_terminal_output | Filter/modify terminal output before agent sees it |
| watch_patterns | Pattern-match stdout → trigger agent notification |
| notify_on_complete | Fire event when background process exits |
| Webhook direct-delivery | Zero-LLM push — forward payloads straight to chat |
| Shell hooks | Wire any shell script as a lifecycle hook (no Python needed) |

### Plugin Hook Points

```python
# Available plugin lifecycle hooks in Hermes
pre_llm_call              # Before LLM API call
post_llm_call             # After LLM response
pre_tool_call             # Before tool execution (can block)
post_tool_call            # After tool returns
transform_llm_output      # Mutate LLM output before delivery
transform_terminal_output # Mutate terminal output before agent
on_session_start          # Session initialization
on_session_end            # Session teardown
```

---

## NATS Integration Strategy (Neither Has It Natively)

**Neither OpenClaw nor Hermes uses NATS out of the box.**

### Option A: OpenClaw Plugin (Recommended for OneMind)

Build a custom OpenClaw plugin using `@openclaw/plugin-sdk` that:
- Subscribes to NATS subjects on gateway startup
- Publishes agent events to NATS subjects
- Exposes NATS publish/subscribe as agent tools

```typescript
// @onemind/openclaw-nats-plugin
import { definePlugin } from '@openclaw/plugin-sdk';
import { connect } from 'nats';

export default definePlugin({
  hooks: {
    async on_gateway_ready(ctx) {
      const nc = await connect({ servers: 'nats://localhost:4222' });
      const sub = nc.subscribe('onemind.events.>');
      for await (const msg of sub) {
        // Inject NATS event into agent session
        ctx.injectEvent({ type: 'nats_message', subject: msg.subject, data: msg.data });
      }
    },
    async post_tool_call(ctx, result) {
      // Publish tool results to NATS for other agents/services
      await ctx.nats.publish(`onemind.tools.${result.toolName}`, result.output);
    }
  },
  contracts: {
    tools: [
      { name: 'nats_publish', description: 'Publish message to NATS subject', parameters: { subject: 'string', payload: 'string' } },
      { name: 'nats_request', description: 'Request/reply on NATS subject', parameters: { subject: 'string', payload: 'string', timeout: 'number' } }
    ]
  }
});
```

### Option B: Hermes Shell Hook + Webhook

```bash
# ~/.hermes/hooks/nats-bridge.sh — shell hook for Hermes
#!/bin/bash
# Subscribe to NATS and forward to Hermes webhook
nats sub "onemind.events.>" --each "curl -X POST http://localhost:4200/webhook -d '{\"subject\": \"{{Subject}}\", \"data\": \"{{Body}}\"}'"
```

### Option C: NATS MCP Server (Platform-Agnostic)

Build a standalone MCP server that both OpenClaw and Hermes connect to:

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  OpenClaw   │────▶│  NATS MCP Server │◀────│  Hermes Agent   │
│  (MCP client)│    │  (bridge)        │     │  (MCP client)   │
└─────────────┘     └────────┬─────────┘     └─────────────────┘
                             │
                    ┌────────▼─────────┐
                    │   NATS Cluster   │
                    │  (pub/sub bus)   │
                    └──────────────────┘
```

---

## Why Plugin > Hook for NATS

| Aspect | Plugin | Hook |
|--------|--------|------|
| Lifecycle | Starts with gateway, persistent connection | Fires per-event, stateless |
| NATS connection | Single long-lived connection | Would reconnect every invocation |
| Bidirectional | Can subscribe AND publish | Hooks only react, don't subscribe |
| Tools | Can register `nats_publish` / `nats_request` as agent tools | Cannot expose new tools |
| State | Plugin maintains subscription state | Hooks are stateless |
| Distribution | npm package, installable by others | Local config only |

**Verdict**: NATS needs a persistent subscriber — that's a **plugin**, not a hook. Hooks fire reactively on events that already happened. A plugin runs continuously and can *create* events from external sources.

---

## Dual-Run Architecture (Mac)

```
┌──────────────────────────────────────────────────────────────┐
│  Mac (M2 Max / M4 Pro)                                       │
│                                                              │
│  ┌─────────────────────────┐  ┌────────────────────────────┐ │
│  │  OpenClaw Gateway       │  │  Hermes Agent Gateway      │ │
│  │  Node.js / TypeScript   │  │  Python 3.11+ async        │ │
│  │  Port: 3210             │  │  Port: 4200                │ │
│  │  ~/Library/App Support/ │  │  ~/.hermes/                │ │
│  │  LaunchAgent managed    │  │  hermes --gateway          │ │
│  │                         │  │                            │ │
│  │  Plugins:               │  │  Plugins:                  │ │
│  │  · @onemind/nats-bridge │  │  · nats shell hook         │ │
│  │  · @onemind/dimos       │  │  · custom providers        │ │
│  └─────────────────────────┘  └────────────────────────────┘ │
│                                                              │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  NATS Server (nats-server)  Port: 4222                   │ │
│  │  Subjects: onemind.events.>, onemind.tools.>,            │ │
│  │            onemind.agents.>, onemind.robotics.>          │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                              │
│  Shared: API keys, MCP servers (DimOS, OM1, NATS MCP)       │
│  Constraint: Separate bot tokens per platform per gateway    │
└──────────────────────────────────────────────────────────────┘
```
