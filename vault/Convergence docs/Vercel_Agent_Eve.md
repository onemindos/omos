## Vercel Agent Eve — The Framework for Building Agents

**Link:** https://vercel.com/eve

**Type:** Agent framework (production-ready)

**Status:** Live, production framework

**Related Resources:**
- **Vercel:** https://vercel.com/
- **Vercel Agents Documentation:** https://vercel.com/docs/agents
- **Vercel Workflows:** https://vercel.com/docs/workflows (underlying execution engine)
- **GitHub:** https://github.com/vercel/ai (AI SDK)

**Core Philosophy:** "Like Next.js for web apps, but for agents" — Directory-based agent framework

---

## What Eve Is

**Eve is a directory-based agent framework** that lets you define an agent using three file types:

1. **instructions.md** — Markdown file defining agent role/behavior
2. **agent.ts** — TypeScript configuration (model selection, runtime config)
3. **tools/** — TypeScript files that become callable tools (no registration needed)
4. **skills/** — Markdown playbooks (conditionally loaded)
5. **channels/** — Integration with external platforms (Slack, Discord, Teams, Web)
6. **connections/** — Authentication for services (GitHub, Stripe, Linear)
7. **subagents/** — Delegate to specialized agents
8. **schedules/** — Durable background jobs
9. **sandbox/** — Customize isolated execution environment

**Installation:** `npx eve@latest init my-agent`

---

## Core Architecture

### Foundation: Directory Structure

```
agent/
├── instructions.md          # Agent role & behavior (Markdown)
├── agent.ts                 # Model selection & config (TypeScript)
├── tools/                   # Tools (TypeScript files → callable functions)
│   ├── weather.ts
│   ├── search.ts
│   └── calculator.ts
├── skills/                  # Reusable playbooks (Markdown)
│   ├── research.md
│   ├── analysis.md
│   └── reporting.md
├── channels/                # Where agent appears
│   ├── slack.ts
│   ├── discord.ts
│   └── web.ts
├── connections/             # Auth for external services
│   ├── github.ts
│   └── stripe.ts
├── subagents/               # Specialized child agents
│   ├── researcher.ts
│   └── analyst.ts
├── schedules/               # Background jobs
│   ├── daily-report.ts
│   └── weekly-digest.ts
└── sandbox/                 # Execution environment
    └── sandbox.ts
```

---

## Component Breakdown

### 1. instructions.md (Agent Identity)

**What it is:** Complete agent definition in Markdown

**Example:**
```markdown
# Identity

You are an expert weather assistant.
You can fetch the weather for any city in the world.

## Responsibilities

- Get current weather
- Provide forecasts
- Answer weather-related questions

## Tools Available

- fetch_weather (city: string) → weather data
- get_forecast (city: string, days: number) → forecast data

## Behavior Guidelines

- Always cite your source
- Provide context for temperature
- Warn about extreme weather
```

**Key:** An instructions.md file alone is a complete agent.

---

### 2. agent.ts (Runtime Configuration)

**What it is:** TypeScript configuration for model selection and runtime

**Example:**
```typescript
import { createAgent } from '@vercel/ai';

export default createAgent({
  model: 'claude-3.5-sonnet',  // Model selection
  maxTokens: 4096,
  temperature: 0.7,
  apiKey: process.env.ANTHROPIC_API_KEY,
  // Custom runtime config
});
```

**When needed:** Optional. Only needed when you want custom model/config.

---

### 3. tools/ — Callable Functions

**What it is:** TypeScript files that become agent-callable tools

**Mechanism:**
- Add a .ts file to tools/
- Filename becomes tool name
- Function exports become tool functions
- No registration required

**Example (tools/weather.ts):**
```typescript
export async function fetchWeather(city: string) {
  const response = await fetch(
    `https://api.weather.com/current?city=${city}`
  );
  return response.json();
}

export async function getForecast(city: string, days: number = 7) {
  const response = await fetch(
    `https://api.weather.com/forecast?city=${city}&days=${days}`
  );
  return response.json();
}
```

**Result:** Agent can now call `fetchWeather()` and `getForecast()` directly

---

### 4. skills/ — Conditional Playbooks

**What it is:** Markdown playbooks loaded when relevant

**Mechanism:**
- Skills are Markdown files
- Loaded only when agent needs them
- Agent gets focused guidance without bloat
- Conditional loading prevents context explosion

**Example (skills/research.md):**
```markdown
# Research Skill

When the user asks you to research something:

1. Use the search tool to find sources
2. Summarize findings into key points
3. Provide citations
4. Ask clarifying questions if needed

## Sources to Prefer

- Academic databases
- Official documentation
- Recent news articles
```

**Benefit:** Agent doesn't carry all skills in every prompt. Only loaded when relevant.

---

### 5. channels/ — Multi-Platform Presence

**What it is:** Connect same agent to different platforms

**Supported Channels:**
- **Slack**
- **Discord**
- **Microsoft Teams**
- **Google Chat**
- **Web Chat** (embedded)
- **WhatsApp** (via Twilio)
- **API** (direct HTTP)
- **Linear** (tool integration)
- **Cron** (scheduled execution)

**Example (channels/slack.ts):**
```typescript
import { createSlackChannel } from '@vercel/eve';

export default createSlackChannel({
  signingSecret: process.env.SLACK_SIGNING_SECRET,
  botToken: process.env.SLACK_BOT_TOKEN,
  // Agent automatically available in configured Slack workspace
});
```

**Key:** Same agent code, multiple surfaces. No duplication.

---

### 6. connections/ — Service Authentication

**What it is:** Handle OAuth/API key management for external services

**Purpose:**
- Tools call services without managing tokens
- Connections handle all auth
- Token rotation automatic

**Example (connections/github.ts):**
```typescript
import { createConnection } from '@vercel/eve';

export default createConnection({
  provider: 'github',
  clientId: process.env.GITHUB_CLIENT_ID,
  clientSecret: process.env.GITHUB_CLIENT_SECRET,
});
```

**Supported Services:**
- GitHub
- Stripe
- Linear
- Slack
- Discord
- And more

---

### 7. subagents/ — Delegation

**What it is:** Create specialized child agents for specific tasks

**How It Works:**
```
Main Agent
├─ Analyze user request
├─ Delegate to Researcher subagent
│   └─ Research returns findings
├─ Delegate to Writer subagent
│   └─ Writer returns draft
└─ Combine results and respond
```

**Example (subagents/researcher.ts):**
```typescript
import { createSubagent } from '@vercel/eve';

export default createSubagent({
  name: 'Researcher',
  instructions: `You are a research specialist. 
    Find information, verify sources, and summarize findings.`,
  tools: ['search', 'fetch_page', 'extract_data'],
});
```

**Benefit:** Parallel delegation, specialized focus, better reasoning

---

### 8. schedules/ — Durable Background Jobs

**What it is:** Run agents automatically on schedule

**Mechanism:**
- Cron-based scheduling
- Durable execution (survives restarts)
- State persisted
- Ideal for reports, digests, monitoring

**Example (schedules/daily-report.ts):**
```typescript
import { createSchedule } from '@vercel/eve';

export default createSchedule({
  cron: '0 9 * * *',  // 9 AM daily
  run: async (agent) => {
    await agent.execute({
      message: 'Generate and send daily report to team@company.com',
    });
  },
});
```

**Use Cases:**
- Daily reports
- Weekly digests
- Monitoring jobs
- Data collection
- Scheduled backups

---

### 9. sandbox/ — Isolated Execution

**What it is:** Customizable sandbox environment

**Default:** Isolated JavaScript/TypeScript runtime

**Customization:** Choose backend, add permissions

**Example (sandbox/sandbox.ts):**
```typescript
import { createSandbox } from '@vercel/eve';

export default createSandbox({
  backend: 'docker',  // or 'node', 'wasmtime'
  permissions: ['fs_read', 'fs_write', 'network'],
  timeout: 30000,     // 30 seconds
  memory: 512,        // MB
});
```

---

## Eve Runtime: Durable Execution

### Core Features

#### 1. Checkpointed Steps

```
Agent receives message
  ↓ CHECKPOINT
Agent processes step 1
  ↓ CHECKPOINT
Agent processes step 2
  ↓ CHECKPOINT
Agent responds
```

**Benefit:** If interrupted, resume from last checkpoint (not from start)

#### 2. Park and Resume

```
Agent waiting for user input
  ↓ PARK (state saved)
User sends reply (potentially hours/days later)
  ↓ RESUME from saved state
Agent continues
```

**Benefit:** Agent survives process crashes, deployments, timeouts

#### 3. Durable Execution Layer

```
Agent Runtime (Eve)
├─ Vercel Workflows (orchestration + checkpointing)
├─ AI Gateway (model calls + streaming)
├─ Vercel Sandbox (isolated execution)
└─ Vercel Connect (auth + MCP endpoints)
```

**Result:** Agent continues even if infrastructure fails

---

## Eve's Vercel Integration Stack

### Layer 1: Runtime
- **Vercel Workflows** — Checkpointed steps, state persistence, event streaming
- **AI Gateway** — Model calls, token counting, streaming, rate limiting
- **Vercel Sandbox** — Isolated execution, file tools, resource limits
- **Vercel Connect** — MCP endpoint support, HTTP endpoints

### Layer 2: Tools & Agents
- Tools (TypeScript functions)
- Subagents (specialized agents)
- Direct function calling

### Layer 3: Channels
- Slack integration
- Discord integration
- Microsoft Teams
- Google Chat
- Web chat (embedded)
- WhatsApp (Twilio)
- API (HTTP)
- Linear
- Cron (scheduled)

### Layer 4: Infrastructure
- Durable execution
- State persistence
- Enterprise observability
- Governance controls

---

## Production Features

### 1. Durable Execution

Every step is checkpointed. Agents park when waiting, resume on next message.

```
Agent Runtime
├─ Checkpoint every step
├─ Persist state to durable storage
├─ Survive crashes + deployments
└─ Resume from last checkpoint
```

### 2. Enterprise Governance

Built-in:
- Access control (RBAC)
- Audit logging
- Resource limits
- Compliance tracking

### 3. Observability

Monitor:
- Agent execution status (Running/Stopped)
- CPU/memory allocation
- Cost per agent
- Error tracking
- Performance metrics

### 4. Sandboxed Compute

Every agent gets:
- Isolated execution environment
- 2-4 CPU cores (configurable)
- 4-8GB RAM (configurable)
- File system (sandboxed)
- Network access (controlled)

---

## Comparison to Other Agent Frameworks

| Feature | Eve | LangChain | AutoGen | Letta |
|---------|-----|----------|---------|-------|
| **Directory-based** | ✅ | ❌ | ❌ | ✅ |
| **Durable execution** | ✅ (Workflows) | ❌ | ❌ | ✅ |
| **Multi-channel** | ✅ | ❌ | ❌ | ✅ |
| **Skill system** | ✅ (Markdown) | ✅ (Tools) | ❌ | ✅ (Skills) |
| **Subagents** | ✅ | ❌ | ✅ | ✅ |
| **Scheduled jobs** | ✅ (Workflows) | ❌ | ❌ | ✅ (Crons) |
| **Enterprise sandbox** | ✅ | ❌ | ❌ | ✅ |
| **Production-ready** | ✅ | Mostly | ❌ | ✅ |
| **Vercel-native** | ✅ | ❌ | ❌ | ✅ |

---

## Getting Started

### Quick Start

```bash
npx eve@latest init my-agent
cd my-agent
npx eve run
```

### Minimal Agent (3 files)

**instructions.md:**
```markdown
# Identity

You are a helpful assistant.
```

**agent.ts:**
```typescript
// Optional - uses default model if omitted
```

**Run:**
```bash
npx eve run
```

---

## Why Eve Matters

### Problem It Solves

Building production agents requires:
- Model integration ✅ (AI Gateway)
- Durable state ✅ (Workflows)
- Multi-channel presence ✅ (Chat SDK)
- Tool execution ✅ (Sandbox)
- Auth management ✅ (Connect)

**Before Eve:** Wire these together manually (months of work)

**With Eve:** Directory structure compiles it all automatically (minutes)

### Key Innovation

```
Traditional: Boilerplate glue code
Eve: Declarative directory → compiled agent
```

### Philosophy

"Like Next.js for web apps, but for agents."

- **Next.js:** app/ directory → web app
- **Eve:** agent/ directory → production agent

---

## Integration with OneMind

### Potential OneMind + Eve Pattern

```
OneMind Fabric (NATS bus)
  ↓
Eve Agent running on edge node
  ├─ instructions.md: "You are a tactical agent"
  ├─ tools/: 
  │   ├─ query_tts_data.ts (queries NATS)
  │   ├─ command_drone.ts (sends commands)
  │   └─ get_cot_position.ts (TAK integration)
  ├─ channels/:
  │   ├─ slack.ts (human interface)
  │   └─ nats.ts (machine-to-machine)
  └─ subagents/:
      ├─ analyzer (analyze intelligence)
      └─ planner (plan response)
```

### Why This Works

1. **Eve framework** provides production agent infrastructure
2. **NATS tools** connect agent to OneMind fabric
3. **Channels** expose agent to humans (Slack) + machines (NATS)
4. **Durable execution** means agent survives edge failures
5. **Subagents** specialize: one analyzes, one plans

---

## Status & Community

**Status:** Production-ready (Vercel official product)

**Model Support:**
- Claude (Anthropic)
- GPT-4 (OpenAI)
- Others via AI Gateway

**Documentation:** Full docs at https://vercel.com/eve

**Community:** Active (Vercel-backed)

---

**Captured:** 2026-07-11

**Type:** Production agent framework

**Status:** Live and actively maintained

**Relevance:** Complete agent framework for building production-grade autonomous agents with durable execution, multi-channel presence, and enterprise features

**Key Innovation:** Directory-based framework (like Next.js for web) + durable workflows + built-in integrations

**Comparison:** Similar to Letta (self-hosting oriented), but Vercel-native and production-ready

**Tactical Value:** Can run Eve agents on edge nodes with NATS integration for OneMind fabric connectivity

**Philosophy:** Framework, not just library. Handles all the infrastructure automatically.

**Best For:** Production agents needing durability, multi-channel, enterprise governance, and minimal boilerplate
