# Course 5: EMPIRE — AI AGENT ARCHITECTURE

> **Design, evaluate, and govern AI systems that run your life at scale.**

**Tier:** Empire ($197/mo)
**Circle Space:** 🤖 AI & Agents (ID: 2617142)
**Prerequisite:** Courses 1, 2, and 3 — full stack running, skills built, Claude Code in workflow
**Target Student:** Someone who wants to go beyond using the system to designing it — building production-grade agent architectures with evals, guardrails, HITL, and observability baked in
**Outcome:** A professionally designed multi-agent system with evaluation frameworks, guardrails at every layer, full observability, human-in-the-loop checkpoints, and reasoning engines that think before they act

---

## What You'll Build

By the end of this course:
- An evaluation framework that measures whether your agents are actually performing
- Guardrail layers at the input, output, and behavior levels
- Human-in-the-loop checkpoints for high-stakes agent decisions
- Full observability: every agent action traced, logged, and reviewable
- Session management that handles scale — concurrent users, multi-turn flows
- A reasoning engine pattern: agents that plan before they act
- A knowledge management architecture your agents never get lost in

---

## Course Structure — 7 Modules, 40 Lessons

---

### MODULE 1: REASONING ENGINES
*"The difference between a chatbot and an agent is that an agent thinks before it acts."*

**Lesson 1.1 — What a Reasoning Engine Is**
A reasoning engine is the cognitive layer between input and action. It plans, decomposes, and decides — before touching tools or files.
- Reactive agents: input → immediate action (chatbot behavior)
- Reasoning agents: input → think → plan → act → verify → report
- The ReAct pattern: Reason + Act interleaved
- The planning problem: why agents that skip planning hallucinate more
- One Mind OS reasoning: why Legacy always plans before executing multi-step tasks

**Lesson 1.2 — The Planning Loop**
How to structure agent prompts so they reason before acting.
- Step 1: Understand — "What is actually being asked?"
- Step 2: Decompose — "What are the sub-tasks?"
- Step 3: Sequence — "What order do they go in? What depends on what?"
- Step 4: Execute — "Do one step at a time, verify each"
- Step 5: Synthesize — "Combine results into a coherent output"
- The planning prompt pattern: forcing explicit reasoning before tool calls
- Hands-on: rewrite one of your agent's system prompts to include explicit planning

**Lesson 1.3 — Extended Thinking Models**
Claude 3.7 Sonnet Extended and o3 — when to route to reasoning models.
- What extended thinking does: internal chain-of-thought before responding
- Use cases: complex multi-step tasks, ambiguous instructions, high-stakes decisions
- Cost tradeoff: reasoning models cost 5–10x more per token
- Routing strategy: standard tasks → Haiku/Sonnet, complex reasoning → extended thinking
- OpenClaw config: `model.reasoning` override for specific agents or tasks
- Hands-on: enable extended thinking for Legacy on planning tasks

**Lesson 1.4 — Multi-Step Task Orchestration**
Breaking complex goals into agent-executable steps.
- The task decomposition pattern: one big goal → numbered steps → one step per agent call
- Dependency graphs: when step 3 can't happen until step 2 completes
- Agent-to-agent delegation: Legacy assigns sub-tasks to domain agents
- The subagent pool: worker agents that execute, report back, get dismissed
- Error handling: what happens when a step fails mid-sequence
- Hands-on: design a multi-step task flow for a real goal in your GE domain

**Lesson 1.5 — Agentic Loops and Stopping Conditions**
Agents that iterate until done — with guardrails that prevent infinite loops.
- The basic loop: execute → evaluate → decide (continue / stop / escalate)
- Stopping conditions: task complete, error threshold hit, human required, max iterations
- The max_iterations parameter: how OpenClaw prevents runaway loops
- Success criteria: how does the agent know it's done?
- The reflection step: agents that grade their own output before reporting
- Hands-on: build a loop with explicit stopping conditions for one automation

---

### MODULE 2: EVALS
*"If you can't measure it, you can't improve it. Evals are how you know your agents are actually working."*

**Lesson 2.1 — Why Evals Matter (And Why Most People Skip Them)**
Evals are the highest-leverage thing you're not doing. This lesson makes the case.
- The problem: you've been shipping agents without knowing if they work
- "It seemed right" is not an eval
- What evals catch: hallucinations, instruction drift, tool misuse, wrong outputs
- The business case: one bad agent action without evals vs systematic quality control
- Types of evals: automated, human review, ground truth comparison

**Lesson 2.2 — The Eval Framework**
How to structure evaluation for One Mind OS agents.
- **Correctness**: did the agent do the right thing?
- **Completeness**: did it cover everything required?
- **Faithfulness**: did it stay true to source material (no hallucination)?
- **Format**: did it output in the expected structure?
- **Tone**: did it match the expected personality / communication style?
- Building an eval rubric per agent type
- The eval dataset: 20 real inputs + expected outputs per agent

**Lesson 2.3 — Automated Evals**
Run evals programmatically — not just by feel.
- LLM-as-judge: using Claude to grade Claude's outputs
- The eval prompt pattern: "Here is the task, here is the response, score it 1–5 on these criteria"
- Building an eval script: input → agent call → LLM grade → log result
- Running evals on schedule: weekly automated eval runs via n8n
- Regression testing: "Did the latest config change break anything?"
- Hands-on: write an eval for your Legacy agent on 5 real tasks

**Lesson 2.4 — Human Eval and the Review Loop**
Some things require human judgment. Build a review loop into your system.
- The human eval queue: outputs that need Zeus review before being trusted
- Flagging criteria: high-stakes decisions, low-confidence outputs, new task types
- The review workflow: agent flags → n8n → Discord message → Zeus approves/rejects → result logged
- Building a feedback signal: good/bad ratings that improve agent behavior over time
- The learning loop: eval results → memory update → agent improves

**Lesson 2.5 — Tracking Improvement Over Time**
Evals are only useful if you track trends.
- The eval log: CSV or database of scores over time
- Metrics to track: average score by agent, score by task type, regression rate
- The improvement cycle: eval → identify weakness → update config → re-eval → confirm fix
- The One Mind OS agent scorecard: monthly eval summary per agent
- Hands-on: set up your eval logging and run your first monthly agent review

---

### MODULE 3: GUARDRAILS
*"Agents are powerful. Guardrails are what make them trustworthy."*

**Lesson 3.1 — The Guardrail Philosophy**
Guardrails are not a limitation on what your agents can do — they're what allows you to trust them with more.
- The spectrum: no guardrails (chaos) → perfect guardrails (useless, can't do anything)
- The goal: maximum capability within defined trust boundaries
- The One Mind OS approach: guardrails are domain-scoped
- Why domain separation IS a guardrail: Spartan doesn't touch your finances
- The autonomy level system (L0–L3) as the primary guardrail framework

**Lesson 3.2 — Input Guardrails**
Validate what goes INTO your agents before they process it.
- Input sanitization: prevent prompt injection attacks
- Topic scoping: agents reject requests outside their domain
- Length limits: prevent context overflow attacks
- Source verification: only trust inputs from known channels
- The input guard layer: a pre-processing step before any agent sees a message
- Hands-on: add input validation to your primary agent's channel binding

**Lesson 3.3 — Output Guardrails**
Validate what comes OUT before it reaches you or takes action.
- Output schema validation: enforce JSON structure on tool outputs
- Confidence thresholds: agents express uncertainty and escalate when unsure
- Action confirmation: agents describe planned action, pause for approval on high-risk ops
- The output review pattern: draft → review → approve → execute
- Sensitive data filtering: strip PII, credentials, private data from outputs
- Hands-on: add output schema validation to one agent's critical tool

**Lesson 3.4 — Behavioral Guardrails**
Constraints on what agents are allowed to DO — not just say.
- The permission matrix: agent × action × domain = allowed/denied
- Write restrictions: which agents can write to which vault paths
- Tool allow-lists per agent: Guardian can access security tools, Spartan cannot
- The `allowedTools` config in openclaw.json
- Rate limiting: prevent agents from spamming APIs or taking too many actions
- The nuclear option: `autonomy: L0` — human approval required for everything

**Lesson 3.5 — The Guardrail Stack in One Mind OS**
All the layers working together as a coherent system.
- Layer 1: Domain scoping (structural)
- Layer 2: Input validation (pre-processing)
- Layer 3: Autonomy levels (decision-time)
- Layer 4: Tool allow-lists (action-time)
- Layer 5: Output validation (post-processing)
- Layer 6: Human review queue (escalation)
- Mapping your current system: where are the gaps?
- Hands-on: audit your agent stack against all 6 layers

---

### MODULE 4: HUMAN IN THE LOOP (HITL)
*"The best AI systems keep humans in control of what matters. HITL is how you design that."*

**Lesson 4.1 — The HITL Philosophy**
Human-in-the-loop is not a weakness — it's a deliberate design decision about where human judgment is irreplaceable.
- The autonomy spectrum: fully manual → HITL → fully autonomous
- Why full autonomy is wrong for high-stakes decisions
- The One Mind OS HITL principle: agents execute, humans decide
- The three HITL patterns: approval-before-action, review-after-action, exception-escalation

**Lesson 4.2 — Approval-Before-Action**
Agent plans action → pauses → you approve → agent executes.
- Use cases: financial transactions, sending external communications, file deletions, calendar changes
- The approval flow: agent → formats proposed action → sends to Discord/Telegram → waits → you respond yes/no → agent acts
- Timeout handling: what happens if you don't respond in 10 minutes?
- Batch approvals: reviewing multiple pending actions at once
- Hands-on: implement approval-before-action for one high-stakes agent behavior

**Lesson 4.3 — Review-After-Action**
Agent acts → logs what it did → you review on your own schedule.
- Use cases: routine tasks, low-stakes automations, information gathering
- The action log: structured record of every agent action with timestamp, context, result
- The daily digest: "Here's everything your agents did today"
- Spot checks: random sampling of agent actions for quality review
- Building the review habit: 5 minutes every morning reviewing the action log

**Lesson 4.4 — Exception Escalation**
Agent tries → hits ambiguity or high-risk → escalates to you instead of guessing.
- The escalation trigger: agent detects it's about to do something uncertain
- Escalation message format: "I was asked to X, but I'm unsure because Y — should I: A or B?"
- Response routing: escalation goes to Discord, you reply, agent continues
- The escalation log: track what agents escalate to identify training opportunities
- Reducing escalations over time: each escalation is a training signal
- Hands-on: add escalation logic to your most autonomous agent

**Lesson 4.5 — Building the HITL Interface**
The practical tools that make HITL smooth and not annoying.
- Discord as your HITL command center: channels for approvals, reviews, escalations
- Telegram bot for mobile approvals: quick yes/no on your phone
- n8n HITL workflows: timeout handling, reminder escalation, logging
- The approval UI: structured messages that make decisions easy
- Keyboard shortcuts: quick approve/deny from anywhere
- The goal: HITL should add 60 seconds to your day, not 60 minutes

---

### MODULE 5: SESSIONS AND MEMORY MANAGEMENT
*"Production agents handle many conversations. Session architecture determines whether they stay coherent."*

**Lesson 5.1 — Session Architecture**
How OpenClaw manages concurrent conversations without agents losing track of who they're talking to.
- What a session is: one isolated conversation context between one agent and one peer
- Session lanes: main (4), cron (3), subagent (8), nested
- Per-channel-peer isolation: Zeus's Discord thread vs your assistant's thread vs a subagent's work
- Concurrent sessions: how many can run simultaneously, how to raise the limit
- Session lifecycle: created → active → idle → expired → cleaned up

**Lesson 5.2 — Session State and Continuity**
How agents remember what happened in a conversation — and across conversations.
- In-session state: the conversation history in the current context window
- Cross-session memory: what persists after a session ends (memory files, vault writes)
- The handoff problem: new session with a returning user — how does the agent remember them?
- Session metadata: storing context that survives session boundaries
- The "catch me up" pattern: agent reads recent memory on session start

**Lesson 5.3 — Session Management at Scale**
When you have community members, clients, or automated services all talking to your agents.
- Multi-user sessions: each user gets isolated context
- Agent identity per user: same agent, different memory scope per user
- Rate limiting per session: prevent any one user from monopolizing the agent
- Session quotas: Sovereign tier vs Legacy tier vs Empire tier
- The queue: what happens when all session lanes are full?
- Hands-on: design session architecture for your One Mind community (75 members)

**Lesson 5.4 — Memory Systems Deep Dive**
The full memory architecture — all layers, all mechanisms, all tradeoffs.
- Working memory: in-context (fastest, limited by window size)
- Short-term memory: recent session files (fast, survives session end)
- Long-term memory: vault files (persistent, requires search to retrieve)
- Episodic memory: event logs, what happened when (journal-style)
- Semantic memory: knowledge base, facts, concepts (vault knowledge notes)
- Procedural memory: skills, workflows, how-to (vault skill files)
- The retrieval decision: when does the agent use which memory layer?

**Lesson 5.5 — Compaction: The Long Game**
Agents that run for months accumulate context. Compaction is how they stay intelligent.
- The compaction problem: 200K context fills in ~3 weeks of active use
- What compaction does: summarizes old context into compressed memories
- Lossless compaction: preserve all facts, compress the narrative
- The Lossless Claw plugin: DAG-based compaction that doesn't lose information
- Compaction triggers: token threshold, time-based, manual
- Testing compaction: run it, verify your agent still knows what it should
- The memory aging policy: what fades, what stays permanently
- Hands-on: trigger a compaction on your longest-running agent, verify nothing important was lost

---

### MODULE 6: OBSERVABILITY AND TRACING
*"You can't fix what you can't see. Observability is the nervous system of your AI stack."*

**Lesson 6.1 — The Observability Mindset**
Production systems are observable by design, not as an afterthought.
- The 3 pillars: Logs (what happened), Metrics (how it's performing), Traces (how it got there)
- Why observability matters more for AI than traditional software: non-deterministic behavior
- The One Mind OS observability goal: every agent action visible, searchable, and alertable
- The difference between debugging (reactive) and observability (proactive)

**Lesson 6.2 — Logging Agent Behavior**
Capture a complete record of every agent action.
- OpenClaw gateway logs: what's included, format, location
- Structured logging: JSON logs that are machine-parseable
- Log levels: DEBUG, INFO, WARN, ERROR — what belongs at each level
- What to log: tool calls, memory reads/writes, channel messages, errors, latency
- Log rotation: don't let logs eat your disk
- Log aggregation: centralizing logs from all nodes (journald, Docker, n8n)

**Lesson 6.3 — Tracing: Follow One Request End to End**
A trace shows you the complete journey of one agent interaction.
- What a trace contains: timestamps, agent, input, reasoning, tool calls, output, duration
- The trace ID: unique identifier that links all steps of one interaction
- Distributed tracing: tracing across agents when Legacy delegates to domain agents
- OpenTelemetry: the open standard for distributed tracing
- Jaeger or Honeycomb: visualization tools for traces
- The trace-to-debug workflow: something went wrong → find the trace → identify the failure point

**Lesson 6.4 — Debugging Agent Failures**
Systematic debugging methodology for when agents do the wrong thing.
- The 5 failure modes: wrong tool, wrong parameters, wrong interpretation, context overflow, guardrail block
- Debug mode: enabling verbose output in OpenClaw
- The reproduction case: capture the exact input that caused the failure
- Isolation testing: test the agent in isolation vs full stack
- Memory debugging: is the agent acting on stale or wrong memories?
- Tool debugging: is the tool returning what the agent expects?
- The fix → verify → monitor cycle

**Lesson 6.5 — Metrics and Alerting**
Numbers that tell you how your system is performing — before it fails.
- Key metrics: response latency, error rate, tool call success rate, session count, memory size
- Baseline: what does "normal" look like for your stack?
- Alert thresholds: latency > 5s, error rate > 5%, disk > 80%
- DigitalOcean monitoring: built-in CPU/RAM/disk alerts
- Custom metrics: n8n workflow to log custom KPIs to a dashboard
- The weekly metrics review: 5-minute health check on your system

**Lesson 6.6 — The Observability Stack for One Mind OS**
Putting it all together — a practical observability setup for your personal AI stack.
- journald: system service logs (already running)
- Docker logging: container output captured
- n8n execution history: workflow audit trail
- OpenClaw logs: agent conversation and tool traces
- A simple dashboard: Grafana (if you want it) or just structured log files + grep
- The "something broke" workflow: alert fires → check logs → find trace → fix

---

### MODULE 7: KNOWLEDGE MANAGEMENT
*"Your agents are only as smart as what they can find. Knowledge architecture is what makes them brilliant."*

**Lesson 7.1 — The Knowledge Architecture Problem**
Most vaults grow organically and become unsearchable. This lesson teaches intentional knowledge design.
- The two failure modes: too little structure (nothing is findable) vs too much structure (nothing gets added)
- The One Mind OS balance: domain structure + entity types + semantic search
- Knowledge decay: information that was true last year may be wrong today
- The evergreen vs temporal distinction: permanent facts vs dated notes

**Lesson 7.2 — Why Git-Based PKM is Superior**
A comprehensive argument for why storing your knowledge in Git (vs Notion, Obsidian Sync, Roam) is the right architectural choice.
- **Portability**: plain markdown files work with any tool, now or in 20 years
- **Version history**: every change tracked, every version recoverable
- **Branching**: experiment with vault restructuring without risking your live system
- **Diff and audit**: see exactly what changed and when
- **AI-native**: any LLM can read plain text files — no API required
- **Sovereignty**: your data lives on your hardware, not in a SaaS database
- **Collaboration**: git workflows enable multi-agent write coordination without conflicts
- **Offline-first**: your vault works without internet
- The practical argument: your VPS agents sync via git in 3 seconds — try that with Notion
- Comparison: Notion (vendor lock-in, API-only for AI), Obsidian Sync (proprietary), Roam (limited export)

**Lesson 7.3 — Knowledge Entity Design**
The 11 entity types of One Mind OS — and how they form a knowledge graph.
- Goals → Projects → Tasks (the execution hierarchy)
- Knowledge → Skills → Tools (the intelligence hierarchy)
- Metrics → Assets → Agents → Humans → Calendar (the context hierarchy)
- Cross-entity linking: how a Goal links to Projects links to Tasks links to Metrics
- The entity relationship graph: visualizing your vault as a knowledge graph
- Frontmatter as graph metadata: `related`, `blocks`, `supports`, `agent`

**Lesson 7.4 — Search Architecture: How Agents Find What They Need**
Getting retrieval right is what separates a useful vault from a digital landfill.
- Keyword search (BM25): fast, exact matches — good for IDs and specific terms
- Semantic search (vector): meaning-based — finds relevant content even with different words
- Hybrid search: both combined — the best of both worlds
- Dataview: structured queries on frontmatter — the SQL of your vault
- Smart Connections plugin: Obsidian-native semantic search
- Agent memory search: `memorySearch.extraPaths` and how it uses your vault
- Hands-on: test all 4 search methods on the same question — compare results

**Lesson 7.5 — Knowledge Maintenance and Decay**
Your vault is only useful if it's accurate. This lesson builds the maintenance habit.
- The review cadence: which notes need quarterly review, which are evergreen
- Staleness markers: frontmatter `review: YYYY-QQ` tag that agents can query
- The stale content queue: "These notes haven't been reviewed in 6 months"
- Archival vs deletion: when to move to `_archive/` vs permanently remove
- The `× Multiply` integration: knowledge maintenance is part of the review loop
- Agent-assisted maintenance: asking Trinity to surface outdated research notes

**Lesson 7.6 — The Agent Knowledge Handoff**
When one agent's output becomes another agent's input. Designing for continuity.
- The structured output contract: standardized formats between agents
- The knowledge pipeline: research (Trinity) → strategy (Legacy) → execution (Empire/Grid)
- Shared memory zones: vault paths multiple agents can read but only one can write
- The knowledge quality standard: what makes a note useful for an agent (vs just for a human)
- The living document pattern: notes that agents update automatically over time

---

## Resources Included

| Resource | Format | Access |
|----------|--------|--------|
| Eval rubric template (5 agents) | Spreadsheet | Empire |
| Eval script (automated LLM-as-judge) | Python | Empire |
| Guardrail audit checklist | PDF | Empire |
| HITL Discord workflow (n8n JSON) | Import | Empire |
| Session architecture diagram | SVG | Empire |
| Observability stack setup guide | Written | Empire |
| Knowledge graph visualization (Obsidian config) | Download | Empire |
| "Why Git-Based PKM is Superior" | Full guide | Empire |
| Agent performance scorecard template | Spreadsheet | Empire |

---

## What Comes Next

You've completed the full One Mind OS curriculum. You can now:
- Teach this system to others (Certification path — coming)
- Build custom implementations for clients (Setup-as-a-Service)
- Publish skills and MCP servers to the community
- Apply for the One Mind Architect certification

*The door is open. The system is yours.*
