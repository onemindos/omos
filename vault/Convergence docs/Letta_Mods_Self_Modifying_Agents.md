# Letta Mods: Self-Modifying Agent Runtime Extensions

**Link:** https://docs.letta.com/letta-agent/mods

**Type:** Agent runtime extension system (part of Letta Agent platform)

**Status:** Active, production-ready feature in Letta Code

**Language:** TypeScript/JavaScript (JS, MJS, TS, TSX)

**Architecture:** Trusted local code that modifies agent runtime

---

## What It Is

**Letta Mods are self-modifying runtime extensions** that allow AI agents to alter their own execution environment.

**Core Concept:**
```
Traditional AI:
  Agent operates within fixed constraints
  └─ Cannot change its own environment

Letta with Mods:
  Agent runs as self-modifying system
  └─ Can add/change tools, commands, hooks, policies
  └─ Can adapt environment on-the-fly
```

**Key Distinction (Per Letta Docs):**
```
Skills:    Teach agent HOW to do a task
Mods:      Change the HARNESS itself (constraints, environment)

Example:
  Skill: "How to write a poem"
  Mod:   "Add a /poetry command that stores outputs"
```

---

## What Mods Can Do

### 1. Add Commands (Slash Commands)

**Capability:**
```typescript
/review      — Show git diff
/whereami    — Display current agent and working directory
/custom      — Any custom command
```

**How It Works:**
```
User types: /review
  ↓
Mod intercepts command
  ↓
Executes registered handler
  ↓
Returns output to agent
```

**Example (From Docs):**
```typescript
// ~/.letta/mods/whereami.ts
export default function activate(letta) {
  return letta.commands.register({
    id: "whereami",
    description: "Show the active agent and working directory",
    run(ctx) {
      return {
        type: "output",
        output: `Agent: ${ctx.agent.name ?? ctx.agent.id}\nCWD: ${ctx.cwd}`,
      }
    }
  })
}
```

### 2. Register Tools/Functions

**Capability:**
Agent can call registered tools directly in its reasoning.

**Example:**
```
Agent: "I need to analyze this code. Let me use the code_analyzer tool..."
  ↓
Mod provides: code_analyzer(path, options)
  ↓
Agent invokes tool
  ↓
Tool returns analysis
```

### 3. Event Hooks

**Capability:**
Intercept and modify agent behavior at key points.

**Available Hooks:**
```
on_turn_start    — Before agent processes input
on_tool_call     — Before tool execution
on_tool_result   — After tool returns
on_message       — Message sent/received
on_error         — Error occurred
```

**Use Case:**
```
Hook: on_turn_start
  └─ Agent receives input
  └─ Mod runs: "Check if user is authorized"
  └─ If not authorized: Block turn, return permission error
  └─ If authorized: Continue normally
```

### 3. Permission Policies

**Capability:**
Define what agent can access/do.

**Example:**
```
Policy: "Agent can only access /home/user/projects/"
  ↓
Agent tries: access /etc/passwd
  ↓
Mod policy blocks: Permission denied
```

### 4. Provider Adapters

**Capability:**
Customize how agent integrates with external systems.

**Example:**
```
Mod adapter for custom LLM provider:
  Input:   Agent's reasoning request
  Process: Translate to provider's format
  Call:    Custom LLM API
  Output:  Convert response back to Letta format
```

### 5. State Management

**Capability:**
Persist and access mod-specific state.

**Use Case:**
```
Mod tracks: User preferences, session state, custom config
  ↓
Agent benefits: "Remember I like concise responses"
  ↓
Mod implements: Hook on_turn_start filters to short responses
```

### 6. UI/Terminal Modifications

**Capability:**
Extend command-line interface or embedded UI.

**Example:**
```
Mod adds: Custom status display
  └─ Shows real-time memory usage
  └─ Shows active tools
  └─ Shows permission violations
```

---

## How Mods Work: Lifecycle

### 1. Discovery (Startup)

```
Letta starts
  ↓
Scans: ~/.letta/mods/
  ↓
Finds: All .js, .mjs, .ts, .tsx files
  ↓
Loads: Each file's activate() function
```

### 2. Activation

```
Mod file exports: activate(letta) function
  ↓
Letta calls: activate(letta)
  ↓
Mod registers: Commands, tools, hooks, policies
  ↓
Returns: cleanup function (optional)
```

### 3. Runtime

```
Agent operates
  ↓
Events occur (turn start, tool call, message, etc.)
  ↓
Hooks execute: Mod code runs
  ↓
Commands available: /command accessible
  ↓
Tools available: Agent can invoke
```

### 4. Reload

```
User edits mod file
  ↓
Runs: /reload (slash command)
  ↓
Letta reloads all mods WITHOUT restarting
  ↓
Changes take effect immediately
```

### 5. Cleanup

```
If mod returns cleanup function:
  ├─ On reload: Previous cleanup runs first
  ├─ On shutdown: Cleanup runs
  └─ Cleanup cancels: Timers, subscriptions, resources
```

---

## Architecture: Self-Modifying Agent

### The Loop

```
Agent operates
  ↓
Agent recognizes problem:
  "I don't have a tool to do X"
  "I can't access resource Y"
  "Permission policy is wrong for my task"
  ↓
Agent writes a mod (agent-assisted):
  "Write a Letta mod that adds the X tool"
  ↓
Mod saved to: ~/.letta/mods/new-tool.ts
  ↓
Agent requests: /reload
  ↓
Mods reloaded (no restart needed)
  ↓
Agent now has new capability
  ↓
Agent continues task with new capability
```

### Key Feature: Agent-Driven Mods

**Design Principle:**
Mods are meant to be written by agents, not humans.

**Why:**
- Agents understand code better than humans understand agent needs
- Agents have deep access to codebase (can inspect mod docs)
- Mods are rapid iteration (reload without restart)
- Agents can test and refine mods in real-time

**Example Prompt:**
```
"Write a Letta Code mod that adds a /review command 
for reviewing my current git diff."

Agent can:
  1. Inspect existing mods in ~/.letta/mods/
  2. Read mod documentation in codebase
  3. Write new mod
  4. Tell you: "Run /reload when ready"
  5. Test new command
  6. Refine if needed
```

---

## Tactical Applications

### Scenario 1: Adaptive Agent in Contested Environment

**Situation:**
```
Agent deployed to forward base.
Environment hostile (network attacks, permission changes, etc.).
Agent needs to adapt capabilities on-the-fly.
```

**How Mods Enable This:**
```
Agent recognizes threat:
  "My usual reconnaissance tool is blocked"
  ↓
Agent writes mod:
  "Add alternative_recon tool using different API"
  ↓
Mod activated: /reload
  ↓
Agent continues operation with new capability
  ↓
No restart, no downtime, adaptive on-the-fly
```

### Scenario 2: Self-Improving Command Suite

**Situation:**
```
Agent starts with basic commands (/help, /status).
As agent operates, it learns what commands would help.
Agent enhances itself continuously.
```

**How Mods Enable This:**
```
Hour 1: Agent has /help, /status
  ↓
Hour 2: Agent adds /analyze (learned need)
  ↓
Hour 3: Agent adds /predict (learned need)
  ↓
Hour 4: Agent adds /synthesize (learned need)
  ↓
Result: Agent capabilities grow over time
```

### Scenario 3: Operator-Requested Mod

**Situation:**
```
Commander needs agent to track specific KPI.
Wants real-time dashboard display.
Doesn't want to restart agent.
```

**How Mods Enable This:**
```
Commander: "Add /kpi command showing current status"
  ↓
Agent writes mod: Hooks into on_turn_start
  ↓
Calculates KPI from agent state
  ↓
Displays: /kpi status in real-time
  ↓
Commander sees live dashboard without restart
```

### Scenario 4: Permission Policy Adaptation

**Situation:**
```
Agent begins in restricted environment.
As trust increases, should get more access.
Manually changing policies = downtime.
```

**How Mods Enable This:**
```
Day 1: Restricted to /tmp/
  ↓
Day 2: Agent proves reliability
  ↓
Policy mod updated: Now access ~/projects/
  ↓
/reload (no downtime)
  ↓
Agent continues with expanded permissions
```

---

## Integration with OneMind OS / Tactical C4ISR

### Potential Integration Points

**Letta Mods with OneMind:**

```
OneMind Tactical Agent
├─ TAK integration mod
│   └─ Reads CoT events from TAK Server
│   └─ Converts to agent reasoning
│   └─ Agent processes tactical data
│
├─ NATS integration mod
│   └─ Subscribes to NATS subjects
│   └─ Routes real-time data to agent
│   └─ Agent responds via NATS
│
├─ PostGIS integration mod
│   └─ Queries spatial data
│   └─ Interprets geography for tactics
│   └─ Agent plans with spatial awareness
│
├─ DDF integration mod
│   └─ Reads intelligence catalog
│   └─ Correlates with current ops
│   └─ Agent provides intel-driven decisions
│
└─ Adaptive hooks
    └─ on_threat_detected: Escalate to commander
    └─ on_resource_denied: Switch to backup tool
    └─ on_permission_change: Adjust behavior
```

### Self-Modifying Tactical Agent

```
Agent deploys to tactical network
  ↓
Agent discovers:
  "I can't reach TAK Server through normal port"
  ↓
Agent writes mod:
  "Use alternate TAK protocol through backdoor port"
  ↓
Mod: /reload
  ↓
Agent regains TAK connectivity
  ↓
Operations continue without commander intervention
```

### Agent-Driven Capability Enhancement

```
Commander: "I need better threat analysis"
  ↓
Agent: "I can add a threat_score mod that correlates DDF intelligence with positions"
  ↓
Agent writes mod that:
  └─ Queries DDF for threat level
  └─ Compares with TAK positions
  └─ Calculates threat score
  └─ Displays in /threat command
  ↓
Mod: /reload
  ↓
Commander now has /threat command
  ↓
No new coding, agent self-enhanced
```

---

## Technical Details

### File Structure

```
~/.letta/mods/
├── whereami.ts          (single file = single mod)
├── threat-analyzer.ts   (another mod)
└── custom-tools.js      (JS also supported)
```

### Mod File Template

```typescript
// ~/.letta/mods/my-mod.ts

export default function activate(letta) {
  // Check what capabilities are available
  if (!letta.capabilities.commands) return;
  
  // Register command
  const cleanup = letta.commands.register({
    id: "my-command",
    description: "What this command does",
    run(ctx) {
      // ctx has access to:
      // - ctx.agent (current agent)
      // - ctx.cwd (working directory)
      // - ctx.input (user input)
      // - ctx.state (mod state)
      
      return {
        type: "output",
        output: "Result here"
      }
    }
  })
  
  // Return cleanup function
  return () => {
    cleanup()  // Unregister command on reload
  }
}
```

### Supported File Types

- `.js` (JavaScript)
- `.mjs` (ES Module JavaScript)
- `.ts` (TypeScript)
- `.tsx` (TypeScript + JSX)

### Reload Command

```
User: /reload
  ↓
Letta: Re-discovers all mods in ~/.letta/mods/
  ↓
For each loaded mod: Run cleanup() if defined
  ↓
For each mod file: Load and activate(letta)
  ↓
Result: All mods updated without restart
```

### Trusted Code Requirement

**Critical Security Model:**

```
Mods run as fully trusted code with:
  ├─ Full filesystem access
  ├─ Full network access
  ├─ Same process as Letta
  └─ No sandboxing

Therefore: Only install/write mods you trust
```

---

## Comparison: Letta Mods vs. Other Extensibility

| System | Skills | Mods | Custom Agents |
|--------|--------|------|---------------|
| **What** | Teach agent task | Change harness | Write agent from scratch |
| **Scope** | Agent capabilities | System constraints | Everything |
| **Reload** | Hot (immediate) | Hot (immediate) | Cold (restart) |
| **Complexity** | Medium | Advanced | High |
| **Use Case** | "Do X better" | "Let me do X at all" | "Complete redesign" |
| **Security** | Sandboxed | Trusted | Trusted |
| **Agent-written** | Yes | Yes (designed for it) | Sometimes |

---

## Why Mods Matter for Autonomous Agents

### Problem

Traditional agents:
- Fixed capabilities at deployment
- Cannot adapt to new environments
- Cannot extend themselves
- Bottleneck: require engineer restart/recompile

### Letta Solution: Mods

```
Agents can now:
✓ Adapt on-the-fly (no restart)
✓ Extend themselves (write mods)
✓ Operate independently (don't need engineer)
✓ Improve over time (learn → write better mods)
```

### Strategic Implication

**Self-modifying agents can:**
- Operate in contested environments (adapt capabilities)
- Improve without interruption (hot reload)
- Survive permission changes (adapt policies)
- Learn by enhancing themselves (operational improvement)

---

## Example Mods

### Mod 1: Git Diff Reviewer

```typescript
export default function activate(letta) {
  return letta.commands.register({
    id: "review",
    description: "Show git diff and request review",
    async run(ctx) {
      const { exec } = await import('child_process');
      const diff = exec('git diff').toString();
      return {
        type: "output",
        output: `Current changes:\n${diff}`
      }
    }
  })
}
```

### Mod 2: Real-Time KPI Display

```typescript
export default function activate(letta) {
  if (!letta.capabilities.hooks) return;
  
  return letta.hooks.register({
    id: "kpi-monitor",
    on_turn_start: async (ctx) => {
      // Calculate KPI on every turn
      const kpi = calculateKPI(ctx.agent.state);
      console.log(`KPI: ${kpi.value} (${kpi.status})`);
    }
  })
}
```

### Mod 3: Conditional Permission Policy

```typescript
export default function activate(letta) {
  if (!letta.capabilities.permissions) return;
  
  return letta.permissions.register({
    id: "trust-based-access",
    can_access: (resource, agent) => {
      if (agent.trust_level >= 0.9) return true;  // High trust = full access
      if (resource.startsWith('/tmp')) return true;  // Always allow temp
      return false;  // Deny by default
    }
  })
}
```

---

## Related Resources in Convergence

- **convergence-resource-capture (skill)** — Capturing resources like this Letta Mods doc
- **PX4_ArduPilot_Autopilots.md** — Self-modifying control systems for drones
- **TAKpilot.md** — AI agent for tactical decision-making
- **Operational_Genome_Definition.md** — Procedures that mods can implement dynamically

---

**Captured:** 2026-07-11

**Status:** Production-ready feature in Letta Agent

**Relevance:** Self-modifying agents for autonomous operations, especially tactical deployments where restart = downtime/failure

**Key Innovation:** Agents can write mods themselves (agent-driven enhancement)

**Architecture Pattern:** Trusted local code + hot reload = adaptive systems

**Tactical Value:** Agents adapt to new environments, escalate permissions, extend capabilities WITHOUT human intervention

**Strategic Implication:** Forward-deployed agents can improve themselves in real-time, adapt to failures, and operate independently

**Comparison to OneMind:** Similar to NATS subscriptions + dynamic routing. Agents can rewire their own connections on-the-fly.

**Security Model:** Mods are trusted code (no sandboxing). Only install mods you trust.

**Reload Pattern:** `/reload` = hot reload all mods without restarting agent. Critical for continuity in tactical ops.

**Self-Modification Loop:** Agent recognizes problem → writes mod → requests reload → has new capability → continues operation.
