---
id: ui-02-20260111053925
title: "Cognitive Architecture"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Cognitive Architecture"
---

Mr. Delacruz — in our context, **the cognitive architecture** is the blueprint for how Legacy AI _perceives, remembers, reasons, decides, acts,_ and _learns_ across all your apps and devices. Think of it as the brainstem + cortex + memory of Legacy, with clear layers and data flows so every agent shares the right state.

  

**The layers (Legacy AI v1)**

1. **Perception (I/O)**

• Voice (Realtime mic/speaker), chat, screen/vision, file drops.

• Streams events into the brain in real time.

2. **Session/Thread & State Manager**

• Creates a **thread** for every conversation/task, ties it to you, a workspace, an agent, and a goal.

• Holds **working memory** (short-term context window) and a pointer to long-term memories.

3. **Memory System (shared across agents)**

• **Working (short-term):** the last turns, tool results, scratchpads.

• **Episodic:** what happened _when_ (meetings, calls, tasks done).

• **Semantic:** facts about you, your companies, SOPs, contacts, configs.

• **Procedural:** “how-to” skills and tool recipes (automations).

• **Artifacts:** files, transcripts, embeddings for RAG.

• Organized via a **context hierarchy** (Org → Workspace → User → Session) so agents can _share state safely_ and inherit the right memories.

4. **Reasoning & Planning**

• A planner chooses a strategy: reactive answer, multi-step tool plan, or multi-agent delegation.

• Uses scratchpads, tree/graph search, or reflection to improve plans.

5. **Tools & Actuation**

• Uniform tool layer (MCP servers) for ClickUp, Zoho, Circle, Home Assistant, Browser/RPA, GitHub, n8n, etc.

• Actions are **logged** and may require **HITL approval**.

6. **Multi-Agent Orchestration**

• Registry of specialist agents (Finance, Ops, Health, Home, Content).

• Message bus/blackboard so agents can hand off tasks and share intermediate results.

7. **Safety, Policy, & Permissions**

• Guardrails (what tools/data an agent may access), PII handling, secret vault.

• “Explain plan” + approval gates for irreversible actions.

8. **Learning & Consolidation**

• After tasks/convos, summarize → distill → write back to long-term memory.

• Auto-create SOPs from repeated behavior (“proceduralization”).

9. **Observability**

• Traces of thoughts/tools, cost & latency, replay, dashboards. Every act also mirrored to ClickUp for accountability.

  

⸻

  

**The minimal stack that implements it**

• **Agent runtime & state**: **Letta** (ADE + Python/TS SDK) as the _brainstem_ for threads, memory, rules, streaming, multi-agent.

• **Datastores**

• Postgres + **pgvector** (core state & embeddings)

• **Zep** (long-term doc memory + RAG indexes)

• **MinIO/S3** (files/artifacts); optional Redis for fast working caches

• **Voice**: OpenAI **Realtime** (or LiveKit/Agora transport) → Letta webhooks

• **Tools (MCP)**: ClickUp, Zoho, Circle, Home Assistant, GitHub, n8n, Browser/RPA

• **HITL**: Approvals surfaced in ClickUp (required for high-impact actions)

• **Observability**: OpenTelemetry → Grafana/Loki/Prometheus

• **Security**: Vault (secrets), per-agent scopes/roles

  

⸻

  

**How a single interaction flows**

1. “Legacy, book my NYC trip and sync to ClickUp.”

2. Perception converts voice → text; **Session Manager** opens/loads your thread.

3. **Planner** drafts a plan (search flights → hold → create tasks → notify).

4. **Memory lookups**: preferences, past trips, payment methods (semantic); last travel convo (episodic).

5. **Tool calls** (MCP): travel API, calendar, ClickUp tasks; asks you to approve purchase.

6. **Streamed response** while working; **HITL** gate pops in ClickUp.

7. On completion, **Learning** writes a “Travel SOP v1” and updates your preferences.

8. Everything is **logged** and re-indexable for RAG.

  

⸻

  

**Data model (essentials)**

• **thread_id** ↔ user_id/workspace_id/agent_id/goal

• **memories**: {type: working|episodic|semantic|procedural|artifact, scope: org|ws|user|session, embedding, pointer(s)}

• **events**: {timestamp, actor, tool, input, output, cost, policy_decision}

• **policies/roles** per agent & tool

  

⸻

  

**Mapping to your Codex**

  

Use the **Seven Kits** as memory tags across pillars:

• **INTEL** (facts, dashboards) → semantic memory

• **TRAIN** (prompts, curricula) → procedural/semantic

• **SYS/OPS** (SOPs, runbooks) → procedural

• **SEC** (policies, roles) → policy store

• **FIN** (budgets, pricing) → semantic + tools

• **COMM** (templates, tone) → semantic + generation presets

  

⸻

  

**Why Letta at the core**

• Gives you **stateful agents out of the box** (threads, memory, files, multi-agent, streaming) so you don’t stitch these from scratch.

• **Context hierarchy** matches your need for _shared state with boundaries_ (Legacy-wide knowledge vs user/session-specific).

• Plays nicely with your **self-hosted, vendor-agnostic** requirement.

  

⸻

  

**First build (fast path)**

1. Bring up: Letta + Postgres/pgvector + MinIO + Zep (Docker).

2. Connect Realtime voice → Letta webhook; verify streaming + memory writes.

3. Register MCP tools for ClickUp/Zoho/Circle/Home Assistant; add approval policy.

4. Define **memory schema** (episodic/semantic/procedural) and the **context hierarchy** for Org/Workspace/User.

5. Instrument OpenTelemetry; ship logs to ClickUp for audit.

  

If you want, I’ll sketch the exact Compose stack and the memory schema tables next.

---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
