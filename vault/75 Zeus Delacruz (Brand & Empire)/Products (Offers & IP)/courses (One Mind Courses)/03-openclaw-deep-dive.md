# Course 3: OPENCLAW DEEP DIVE

> **Master the engine that runs One Mind OS.**

**Tier:** Legacy ($97/mo)
**Circle Space:** 🤖 AI & Agents (ID: 2617142)
**Prerequisite:** Course 2 — VPS live, OpenClaw installed, at least one agent running
**Target Student:** Someone with OpenClaw running but not fully configured — wants to understand every dial, build production-grade agent identities, use clawhub, run scheduled automations, and wire up all their channels
**Outcome:** All 9 agents fully configured with deep identities, custom skills installed, channels firing, cron jobs running, multi-node if applicable — OpenClaw running like a professional AI operation

---

## What You'll Build

By the end of this course:
- Every agent has a complete IDENTITY, SOUL, USER, and CODEX file
- `openclaw.json` fully understood and optimized — no mystery fields
- Clawhub explored — 3+ custom skills installed
- Discord and Telegram channels fully wired
- 5 cron automations running on schedule
- Agent routing working: the right task goes to the right agent automatically
- (Optional) Multi-node: core + second node running in parallel

---

## Course Structure — 8 Modules, 44 Lessons

---

### MODULE 1: HOW OPENCLAW ACTUALLY WORKS

**Lesson 1.1 — The Architecture From the Inside**
Before you configure it, understand what OpenClaw is and how it runs.
- The gateway: the central process that routes everything
- Agents: isolated processes with their own context, tools, memory
- Channels: how the outside world (Discord, Telegram, webhooks) reaches agents
- Sessions: isolated conversation contexts within an agent
- The message flow: channel receives → gateway routes → agent processes → channel responds
- How this differs from ChatGPT, Claude.ai, or raw API calls
- Why the gateway architecture enables things a stateless API cannot

**Lesson 1.2 — The Gateway Process**
What actually runs on your VPS and how to manage it.
- `openclaw gateway start` — what launches, what ports open
- The gateway as a systemd service: how to check status, restart, view logs
- `journalctl -u openclaw -f` — reading live logs
- The OpenClaw dashboard: what the UI shows, how to read it
- Gateway health: CPU/memory footprint of a 9-agent stack
- The startup sequence: gateway boots → loads config → starts agents → connects channels
- Graceful restart vs hard restart — when to use each

**Lesson 1.3 — Session Lanes and Concurrency**
How OpenClaw handles multiple conversations simultaneously.
- Main lanes (4): your primary interactive sessions
- Cron lanes (3): background scheduled tasks
- Subagent lanes (8): spawned when one agent calls another
- Nested lanes: subagent spawning further subagents
- What happens when all lanes are full: the queue
- Configuring lane counts: when to increase and the tradeoffs
- Per-agent session limits: preventing any one agent from hogging resources

**Lesson 1.4 — The Config System Deep Dive**
`openclaw.json` is the single source of truth. This lesson reads every field.
- The top-level structure: `gateway`, `agents`, `channels`, `models`, `memory`, `plugins`
- Defaults vs overrides: what inherits from the top level vs what's per-agent
- Config validation: how to check for errors before restarting
- `openclaw config get` and `openclaw config set` — CLI vs file editing
- Config reload: hot reload vs restart required changes
- Version-controlled config: why your `openclaw.json` should be in git
- Hands-on: complete annotated walkthrough of a real production config

---

### MODULE 2: AGENT IDENTITY ARCHITECTURE

**Lesson 2.1 — The Four Identity Files**
Every agent is defined by four documents. These are not optional extras — they ARE the agent.
- `IDENTITY.md` — who the agent is: name, emoji, role, domain, personality, tone
- `SOUL.md` — the agent's values, operating principles, how it handles edge cases
- `USER.md` — who the agent is talking to: Zeus's preferences, communication style, what matters
- `CODEX.md` — the vault: where it is, how to navigate it, domain routing rules
- Why four files instead of one long system prompt
- How these files map to the agent's system prompt at runtime
- The load order: IDENTITY → CODEX → USER → SOUL (and why that order matters)

**Lesson 2.2 — Writing IDENTITY.md**
The identity file sets the personality baseline for every interaction.
- The required fields: name, emoji, role, domain, tagline
- Writing the character brief: 2-3 sentences that define the agent's vibe
- Tone calibration: formal vs casual, terse vs expansive, direct vs exploratory
- Domain scope: what this agent handles and explicitly what it does not
- The handoff line: "For X, delegate to Y"
- Anti-patterns: vague identities, conflicting instructions, trying to make one agent do everything
- Hands-on: write or rewrite IDENTITY.md for two of your agents

**Lesson 2.3 — Writing SOUL.md**
The soul file governs behavior in ambiguous and high-stakes situations.
- What SOUL covers that IDENTITY doesn't: values, judgment calls, escalation rules
- The 5 SOUL sections: Mission, Values, Operating Principles, Escalation Triggers, Boundaries
- Writing operating principles: "When in doubt, X" statements that are actually actionable
- Escalation triggers: exactly when does this agent stop and ask for human input?
- Boundary statements: explicit things this agent will not do
- The difference between a soul-less agent (does whatever it's asked) vs a principled one
- Hands-on: write SOUL.md for Legacy, the orchestrator — your most critical agent

**Lesson 2.4 — Writing USER.md**
USER.md makes agents feel like they've known you for years.
- What belongs in USER.md: communication style, preferences, priorities, context
- Zeus's profile: background, what matters to him, how he likes to receive information
- The "don't make me repeat myself" principle: everything Zeus always says goes in USER.md
- Preference sections: response format preferences, depth preferences, frequency preferences
- Context that matters: time zone, language, abbreviations you use, things to always remember
- Updating USER.md: it's a living document — add to it when agents get things wrong
- Hands-on: build or audit your own USER.md — use the last 10 things you had to clarify

**Lesson 2.5 — Writing CODEX.md**
CODEX.md is the agent's map of the vault. Without it, agents wander.
- What CODEX.md contains: vault path, domain map, file conventions, where to write
- The domain routing table: "for HP content, write to 25-49/ for GE content, write to 75-99/"
- The frontmatter spec: what every file the agent creates should include
- Inbox-first rule: new captures always land in 99 Inbox before routing
- The write permissions model: which domains each agent is allowed to write to
- How agents use CODEX.md: it loads at session start and informs every file operation
- Hands-on: audit your current CODEX.md — are the paths correct? Is the routing complete?

**Lesson 2.6 — The 9-Agent Identity Matrix**
Building out all 9 agents with consistent, differentiated identities.
- Why differentiation matters: Legacy and Empire must feel like different minds
- The One Mind OS agent roster and their distinct personalities
- Cross-agent consistency: what all agents share (USER.md, CODEX.md base)
- Domain-specific identity elements: what makes Spartan feel like a fitness coach vs Oracle like a strategist
- The identity audit: reading all 9 IDENTITY.md files side by side for coherence
- Common mistakes: agents that sound the same, agents that scope-creep, agents with no clear handoff
- Hands-on: complete the identity matrix for all 9 — fill any gaps

---

### MODULE 3: THE OPENCLAW.JSON COMPLETE REFERENCE

**Lesson 3.1 — `gateway` Block**
The top-level gateway configuration.
- `port`: default 18789 — when and why to change it
- `trustedProxies`: required for Tailscale reverse proxy setups
- `maxSessions`: total concurrent session limit across all agents
- `timeout`: request timeout — what happens when an agent takes too long
- `cors`: origin whitelist for the dashboard
- `debug`: enable verbose logging for troubleshooting
- `logLevel`: INFO vs DEBUG vs WARN — production settings

**Lesson 3.2 — `agents` Block**
Configuring individual agents — every field.
- `id`: the unique identifier used in routing and logs
- `name`: display name (shown in dashboard and channel messages)
- `model`: primary model — can override the global default
- `workspace`: path to the agent's workspace folder with identity files
- `autonomy`: L0 (human approval all) through L3 (fully autonomous)
- `allowedTools`: explicit tool allowlist — security-critical for domain isolation
- `maxTokens`: context window limit per session
- `temperature`: creativity vs consistency tradeoff per agent
- `memory`: per-agent memory config overrides

**Lesson 3.3 — `models` Block**
Model routing: which AI powers which agent for which task.
- `default`: the baseline model for all agents
- `routing`: conditions that trigger a different model
- `providers`: configuring Anthropic, OpenAI, GitHub Copilot, Ollama
- `fallback`: what to use if the primary model is down
- The GitHub Copilot provider: `github-copilot/claude-sonnet-4.6` config
- Cost routing: cheap model for simple tasks, expensive for complex
- `openclaw models auth login-github-copilot` — re-authentication flow

**Lesson 3.4 — `channels` Block**
Every channel fully configured.
- `discord`: server ID, bot token, channel bindings, DM policy
- `telegram`: bot token, DM allowlist, group behavior
- `webhook`: inbound URLs for n8n and external triggers
- `cli`: the local terminal channel (always on)
- Channel-to-agent bindings: which channel routes to which agent
- The dmPolicy: `allowlist` vs `open` — why allowlist is the production default
- Multiple Discord servers: binding different agents to different servers

**Lesson 3.5 — `memory` Block**
Configuring how agents remember.
- `extraPaths`: directories added to semantic memory search (include your vault)
- `workingMemoryLimit`: max tokens of active context
- `shortTermExpiry`: when recent session files expire
- `longTermPath`: where persistent memory files are written
- `compaction`: trigger thresholds and settings (see Module 7)
- Per-agent memory overrides: some agents remember more than others
- The memory isolation rule: agents can't accidentally read each other's private memory

**Lesson 3.6 — `plugins` Block**
OpenClaw's plugin system for extending core behavior.
- What plugins are: registered modules that hook into the agent lifecycle
- `lossless-claw`: the compaction plugin — install, configure, enable
- Plugin installation: `openclaw plugins install <name>`
- Plugin lifecycle hooks: pre-message, post-message, pre-tool, post-tool
- Building a custom plugin: structure, entry point, registration
- Plugin security: what plugins can and cannot access
- Hands-on: install lossless-claw, run a test compaction

---

### MODULE 4: CLAWHUB — SKILLS ECOSYSTEM

**Lesson 4.1 — What Clawhub Is**
The package manager for OpenClaw agent capabilities.
- Skills vs MCP servers vs built-in tools — what's the difference?
- Clawhub as npm for agent skills: publish, install, version, update
- The official skill index: how to browse available skills
- Community skills vs verified skills
- `clawhub search <query>` — finding the right skill
- `clawhub info <skill>` — reading the skill spec before installing
- The skill trust model: what you should audit before installing a skill

**Lesson 4.2 — Installing and Managing Skills**
`clawhub install` to `clawhub remove` — the full lifecycle.
- `clawhub install <skill>` — where it installs, what it creates
- Skills go to `~/.openclaw/workspace/skills/` — shared across all agents
- Per-agent skills vs shared skills: when to scope a skill to one agent
- `clawhub list` — see everything installed
- `clawhub update <skill>` — keeping skills current
- `clawhub remove <skill>` — clean uninstall
- Skill versioning: pinning to a specific version for stability
- Hands-on: install the `todoist-api-rest` skill, verify it's accessible

**Lesson 4.3 — Using Installed Skills**
Wiring skills into agent conversations.
- How agents discover available skills: the skills/ directory scan
- Calling a skill from a conversation: "use the todoist skill to..."
- The skill invocation format: how the agent knows which skill to call
- Skill parameters: required vs optional, types, validation
- The skill response: what comes back and how to use it
- Skill chaining: using the output of one skill as input to another
- Debugging skill calls: what to check when a skill doesn't work

**Lesson 4.4 — The Anatomy of a Skill**
Understanding skill internals so you can build and debug them.
- The skill folder structure: `README.md`, `INSTALL.md`, `SKILL.md`, `VERIFY.md`
- `SKILL.md` is what the agent reads — the instruction set
- `INSTALL.md`: setup steps (API keys, dependencies, config)
- `VERIFY.md`: how to confirm the skill is working
- The skill instruction format: how to write clear, agent-executable steps
- Tool calls within skills: how skills invoke external APIs via curl/jq
- Error handling in skill instructions: what the agent should do when a step fails

**Lesson 4.5 — Building Your First Custom Skill**
Create a skill for something specific to your One Mind OS workflow.
- Skill ideation: what tasks do you repeat often that could be a skill?
- The skill scaffold: start from the template
- Writing the SKILL.md: step-by-step instructions the agent will follow
- Testing the skill: invoke it in a conversation, watch what happens
- Iterating: the most common fix is adding more precision to the instructions
- Publishing to clawhub (optional): share with the community
- Hands-on: build a skill that captures a note to the vault Inbox and tags it

---

### MODULE 5: CHANNELS — FULL CONFIGURATION

**Lesson 5.1 — Discord: Your AI Command Center**
Discord is the primary interface for your One Mind OS agents.
- Creating the Discord bot: developer portal, scopes, permissions
- Getting the bot token: where to store it (1Password, never in config files)
- Adding the bot to your server: the OAuth URL and required permissions
- Setting the token in OpenClaw: `openclaw config set channels.discord.accounts.main.token`
- Channel bindings: `#legacy` routes to Legacy, `#spartan` routes to Spartan
- DM policy: enabling or restricting who can DM the bots directly
- Role-based access: only certain roles can trigger certain agents
- Hands-on: wire all 9 agents to dedicated Discord channels

**Lesson 5.2 — Telegram: Mobile-First Command**
Telegram for fast, mobile-native agent access.
- Creating the Telegram bot: BotFather walkthrough
- Bot token storage: same as Discord — 1Password
- Setting the token: `openclaw config set channels.telegram.accounts.default.token`
- DM allowlist: `openclaw config set channels.telegram.dmPolicy allowlist`
- Adding yourself to the allowlist
- Use cases: quick captures on the go, approval responses, morning briefing
- Telegram vs Discord: when to use which
- Hands-on: send a message to Legacy via Telegram, verify response

**Lesson 5.3 — Webhooks: Connecting Everything Else**
Webhooks let external tools trigger your agents programmatically.
- The inbound webhook URL: how to get it from the gateway
- Securing webhooks: HMAC signatures, IP allowlisting
- n8n → OpenClaw: triggering agents from n8n workflow nodes
- Todoist → OpenClaw: agent responds when a task status changes
- Calendar → OpenClaw: agent briefings triggered by events
- Custom webhooks: any HTTP POST can trigger an agent with the right payload
- Hands-on: build one n8n → agent webhook flow

**Lesson 5.4 — The CLI Channel**
The terminal is always available — use it well.
- `openclaw chat` — interactive CLI conversation with your default agent
- `openclaw chat --agent spartan` — talk to a specific agent
- `openclaw run "task description"` — one-shot non-interactive execution
- Piping input: `echo "summarize this" | openclaw run --agent oracle -`
- Scripting with the CLI: bash scripts that invoke agent actions
- The CLI in VS Code terminal: using the agent while you code
- Use case: running agents from Claude Code via `run_in_terminal`

**Lesson 5.5 — Channel Routing Logic**
How the gateway decides which agent handles an incoming message.
- The routing priority: explicit binding > keyword match > default agent
- Channel-to-agent bindings: one channel = one agent (default)
- Keyword routing: messages containing "finance" route to Oracle
- The `@mention` override: `@spartan do 50 pushups` goes directly to Spartan
- Fallback routing: what handles messages that match nothing
- Cross-agent delegation: Legacy receives, routes to specialist, returns unified response
- The routing config: `agents[*].routing` in openclaw.json

---

### MODULE 6: CRON — SCHEDULED INTELLIGENCE

**Lesson 6.1 — Why Cron Changes Everything**
Cron turns reactive agents into proactive ones. This is where One Mind OS becomes autonomous.
- The difference between on-demand and scheduled: from tool to system
- What cron enables: morning briefings, end-of-day reviews, weekly reports, health pings
- OpenClaw native cron vs n8n cron: when to use each
- The cron lane: 3 dedicated lanes that never compete with interactive sessions
- Cron as the "× Multiply" step: the system that runs without you

**Lesson 6.2 — Configuring Cron Jobs**
Writing cron schedules in `openclaw.json`.
- Cron syntax: `0 7 * * *` (daily at 7am) — the 5 fields explained
- The `cron` block in agent config: schedule + prompt + channel
- Multiple cron jobs per agent: Legacy runs 3 different schedules
- Time zones: setting the cron TZ on the VPS
- The cron prompt: writing the instruction the agent executes on schedule
- Output routing: where does the cron result go? (Discord channel, vault file, both)
- Enabling/disabling cron without deleting it

**Lesson 6.3 — The 5 Essential Cron Automations**
Build these five — they run One Mind OS on autopilot.
- **Morning Briefing** (7:00 AM): Legacy reads active goals, today's Todoist tasks, weather, sends to Discord
- **EOD Review** (9:00 PM): Oracle reviews what was completed today, what carries over, sends summary
- **Weekly Report** (Sunday 8:00 AM): Legacy reads the week's vault changes, summarizes progress, sends to Discord
- **Vault Health Check** (Monday 6:00 AM): Trinity checks for stale notes, orphaned files, missing frontmatter
- **Memory Compaction** (Sunday 3:00 AM): Legacy triggers compaction on all agents with full sessions
- Hands-on: implement all 5 with actual cron prompts and Discord channel output

**Lesson 6.4 — Building Custom Cron Workflows**
Design your own schedules for your specific domains.
- GE: Empire runs weekly portfolio review every Friday
- HP: Spartan sends Monday workout plan, Spartan sends Friday progress report
- LE: Haven sends monthly family check-in prompt
- SO: Guardian runs security scan every 72 hours
- Cron as a publishing tool: agents that post to Circle.so on schedule
- Cron failure handling: what happens if the cron job errors?
- Hands-on: design one domain-specific cron workflow for GE or HP

**Lesson 6.5 — n8n + OpenClaw: Hybrid Scheduling**
Some automations are better in n8n. Know which to put where.
- When n8n wins: multi-step workflows with conditionals, external API calls, data transformation
- When OpenClaw cron wins: simple agent execution, vault writes, Discord sends
- The webhook bridge: n8n triggers OpenClaw via webhook
- The pattern: n8n handles data prep → fires webhook → OpenClaw agent processes → writes to vault
- Example: Todoist completed task → n8n → Legacy adds to vault archive
- Example: Calendar event ends → n8n → Oracle writes a reflection prompt to Inbox
- Hands-on: build one hybrid n8n + OpenClaw workflow

---

### MODULE 7: MEMORY AND COMPACTION

**Lesson 7.1 — Memory Architecture Recap**
What each memory layer does and where it lives on disk.
- Working memory (in-context): the active conversation, rebuilt each session
- Short-term memory: recent session files in the workspace
- Long-term memory: vault files, semantically indexed
- The `memorySearch.extraPaths` config: how the vault becomes searchable
- Semantic search under the hood: vector embeddings, similarity threshold
- The indexing process: how your vault content gets indexed (and when it re-indexes)
- Memory isolation: why Oracle's memories don't leak into Spartan

**Lesson 7.2 — Teaching Agents What They Need to Know**
Your agents are only as smart as their memory. This lesson is about feeding them well.
- The startup sequence: what agents read before the first message
- Pinned context vs memory search: what's always loaded vs what's retrieved on demand
- Writing vault notes that agents can actually use (vs notes written for humans)
- The knowledge note format: title, frontmatter, concise facts, no narrative padding
- Domain briefings: one `README.md` per domain that orients the agent
- Updating agent memory: when to write a new memory file vs update an existing one
- The feedback loop: agent gets something wrong → update the relevant file → re-test

**Lesson 7.3 — Compaction Deep Dive**
How to run compaction and verify nothing was lost.
- What compaction does: takes long conversation history, compresses to summary memories
- The lossless-claw plugin: DAG-based compaction that preserves factual accuracy
- Compaction trigger options: token threshold (auto), manual, scheduled via cron
- Running manual compaction: `openclaw compaction run --agent legacy`
- The compaction report: what was kept, what was summarized, what was dropped
- Post-compaction verification: test 5 things the agent should still know
- The compaction log: stored in the workspace, reviewable
- Hands-on: trigger a compaction on your most active agent, verify with a 5-question quiz

**Lesson 7.4 — Memory Hygiene**
Keeping memory clean over months of use.
- Memory rot: outdated facts that are worse than no facts
- The staleness audit: files with old dates that need review or deletion
- Agent-assisted cleanup: asking Trinity to surface memory files older than 90 days
- The archive pattern: don't delete old memory — move it to `_archive/memory/`
- Rebuilding from scratch: when to wipe an agent's memory and start clean
- The backup discipline: always backup before any major memory operation

---

### MODULE 8: MULTI-NODE AND ADVANCED PATTERNS

**Lesson 8.1 — Why Multi-Node?**
One VPS is enough to start. Two nodes unlocks specialization.
- The single-node limit: all 9 agents competing for one box's CPU/RAM
- The multi-node model: core handles orchestration, titan handles heavy compute
- Use cases for a second node: n8n, Postiz, heavy search, video processing
- Tailscale as the mesh: how nodes communicate securely without public exposure
- The cost math: 2× $12/mo = $24/mo for a professional AI operation
- When NOT to go multi-node: premature optimization

**Lesson 8.2 — Configuring Core + Titan**
The current One Mind OS architecture: `onemindos-core` + `omos-titan`.
- Core's role: Legacy orchestrator, all 9 agents, vault sync, primary gateway
- Titan's role: n8n workflows, Postiz social publishing, heavy background tasks
- Tailscale between nodes: `omos-titan` reachable at `omos-titan.tail717ba0.ts.net`
- Cross-node webhook calls: Core n8n → Titan workflow, Titan result → Core vault
- SSH between nodes: `ssh omos-titan "command"` from Core
- Shared secrets: how both nodes access the same 1Password vault
- The deployment script: provisioning a second node from scratch

**Lesson 8.3 — Agent-to-Agent Communication**
How Legacy delegates to specialist agents — and how specialists report back.
- The delegation pattern: Legacy receives task → decomposes → assigns subagent
- Subagent invocation: how OpenClaw spawns a subagent session
- The subagent contract: input format → expected output format
- Parallel execution: multiple subagents working simultaneously on independent steps
- The merge step: Legacy combines subagent outputs into a unified response
- Error handling: what Legacy does when a subagent fails or times out
- Hands-on: design a multi-agent workflow for a real GE task (e.g., competitive research)

**Lesson 8.4 — Production Hardening**
Your stack is live. Make it bulletproof.
- Process supervision: why systemd is better than `nohup` or `screen`
- Automatic restart on failure: `Restart=always` in the service unit
- Gateway health endpoint: `/health` — how to ping it from monitoring
- The watchdog pattern: a cron job that pings the health endpoint and alerts on failure
- Log rotation: `logrotate` for gateway logs so disk doesn't fill
- Dependency order: ensure vault sync is up before OpenClaw starts
- The disaster recovery test: shutdown the VPS, restore from backup, verify everything comes back
- Hands-on: write the systemd service unit for your gateway if not already done

**Lesson 8.5 — Graduation: Your Production AI Stack**
You've built a professional AI operation. This lesson documents it.
- The One Mind OS stack diagram: draw your architecture (all services, all connections)
- The runbook: a single document that lets you rebuild everything from scratch
- The config audit: one final pass through openclaw.json — everything intentional
- The identity audit: all 9 agents have complete IDENTITY, SOUL, USER, CODEX
- The channel audit: Discord, Telegram, and webhooks all firing correctly
- The cron audit: 5 core automations confirmed running
- The memory audit: vault indexed, compaction tested, memory healthy
- Graduation checkpoint: 30 minutes of normal use — everything working as expected

---

## Resources Included

| Resource | Format | Access |
|----------|--------|--------|
| Annotated `openclaw.json` (9-agent production config) | JSON | Legacy+ |
| All 9 IDENTITY.md templates | Markdown | Legacy+ |
| All 9 SOUL.md templates | Markdown | Legacy+ |
| USER.md template | Markdown | Legacy+ |
| CODEX.md template | Markdown | Legacy+ |
| 5 cron automation prompts | Markdown | Legacy+ |
| Custom skill scaffold | Folder template | Legacy+ |
| Discord bot setup guide | PDF | Legacy+ |
| Multi-node provisioning script | Bash | Legacy+ |
| Stack architecture diagram (editable) | SVG | Legacy+ |
| Production runbook template | Markdown | Legacy+ |
