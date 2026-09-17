# Sprint Research Summary: Workflow Engines & Agent Frameworks

> **Document Status:** RESEARCH FINDINGS
> **Created:** 2026-01-29
> **Purpose:** Summarize research on workflow engines and agent frameworks for OneMind enhancement

---

## Executive Summary

After researching Windmill, workflow engines, and Relevance AI's stack, here are the key findings and recommendations for supercharging OneMind.

---

## 1. Windmill - The Standout Choice

### What It Is
Open-source platform for building internal tools and automations. **Rust backend, PostgreSQL storage, Python/TypeScript execution**.

### Key Features
| Feature | Why It Matters |
|---------|----------------|
| **13x faster than Airflow** | Sub-20ms job overhead |
| **Auto-generated UIs** | Script parameters → instant UI |
| **Built-in approvals** | HITL workflows native |
| **Webhooks** | Every script gets a trigger URL |
| **Cron scheduling** | Time-based triggers built-in |
| **Multi-language** | Python, TypeScript, Go, Bash, SQL |
| **Self-hostable** | No vendor lock-in |

### Architecture
```
┌─────────────────────────────────────────────────────────┐
│                    WINDMILL                              │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │   Scripts   │ →  │  Job Queue  │ →  │   Workers   │ │
│  │ (Python/TS) │    │ (PostgreSQL)│    │   (Rust)    │ │
│  └─────────────┘    └─────────────┘    └─────────────┘ │
│                                                         │
│  • Webhooks for triggers                               │
│  • Cron for scheduling                                 │
│  • Auto-UI for inputs                                  │
│  • Approval gates                                      │
│  • Audit logging                                       │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Recommendation for OneMind
**EVALUATE FOR HEAVY JOBS** - Don't replace our workflow engine, but consider Windmill for:
- Long-running data processing
- External API orchestration
- Jobs that need robust retry/monitoring
- Auto-UI for admin tools

**Our existing workflow engine** is good for:
- Agent orchestration (tight integration)
- Fast, lightweight flows
- Team coordination

---

## 2. Other Workflow Engines Evaluated

### Prefect
- **Python-native** with decorators (@flow, @task)
- **Strong scheduling** and retry logic
- **Good for ML pipelines**
- Less suitable for real-time agent work

### Temporal
- **Distributed microservice orchestration**
- **Excellent for long-running processes**
- Overkill for our use case

### CrewAI
- **Multi-agent orchestration framework**
- **Hierarchical and sequential processes**
- **Similar to what we want to build**
- Could inform our design

### LangGraph
- **Graph-based agent orchestration**
- **State management for agents**
- **Human-in-the-loop support**
- Worth studying for patterns

---

## 3. Relevance AI's Tech Stack

### Open Source Components
| Repo | Purpose |
|------|---------|
| **relevanceai** | Core agent framework (Apache-2.0) |
| **vectorhub** | Embeddings & vectorization |
| **relevance-js-sdk** | TypeScript SDK |
| **relevance-connect** | Python integration module |

### Key Patterns We Should Adopt
1. **Smart Approval Modes** - approval/agent_decides/autorun
2. **@Mention Routing** - Route to specialist agents
3. **Workforce Canvas** - Visual multi-agent coordination
4. **Triggers** - Webhooks, schedules, events
5. **Memory** - Extract and recall from agent outputs

---

## 4. What We Have vs What We Need

### Already Have (Don't Rewrite)
```
✅ Workflow Engine (9 step types)
   backend/agno/workflows/workflow/

✅ Agent Infrastructure
   backend/agno/agents/

✅ Approval System
   backend/agno/approvals/

✅ State Management
   backend/consciousness/awareness/state.py

✅ Team Coordination
   backend/agno/teams/

✅ All LifeOS Modules
   backend/lifeos/tasks/, habits/, goals/, etc.
```

### Need to Add (Enrich)
```
🆕 Context Engine
   - Aggregates all data for universal context
   - Every chat gets full awareness

🆕 Intelligence Scores (HPI/LEI/GEI/UI)
   - Calculate from existing data
   - New service + endpoints

🆕 DO DATE Algorithm
   - Motion-style intelligent scheduling
   - Enhance existing scheduler

🆕 Smart Approval Modes
   - Add to existing approval system
   - agent_decides logic

🆕 @Mention Routing
   - Add to chat processing
   - Route to specialist agents

🆕 /Command Shortcuts
   - Add to chat processing
   - Quick actions

🆕 Frontend Restructure
   - New navigation
   - Command Center dashboard
   - Score widgets
```

---

## 5. Recommended Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    ONEMIND ARCHITECTURE                      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  PRESENTATION LAYER                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Flutter App (Mobile-First)                          │   │
│  │  • Command Center Dashboard                          │   │
│  │  • Chat with @mentions, /commands                    │   │
│  │  • Score widgets                                     │   │
│  │  • New navigation structure                          │   │
│  └─────────────────────────────────────────────────────┘   │
│                         │                                   │
│  CONSCIOUSNESS LAYER                                        │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Context Engine (NEW)                                │   │
│  │  • State aggregator                                  │   │
│  │  • Items aggregator                                  │   │
│  │  • Scores aggregator (HPI/LEI/GEI/UI)               │   │
│  │  • Calendar aggregator                               │   │
│  │  • History aggregator                                │   │
│  │  • Assets aggregator                                 │   │
│  └─────────────────────────────────────────────────────┘   │
│                         │                                   │
│  ORCHESTRATION LAYER                                        │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Existing Agno Framework (ENHANCED)                  │   │
│  │  • Agents with smart approvals                       │   │
│  │  • Teams with workforce patterns                     │   │
│  │  • Workflows with HTTP triggers                      │   │
│  │  • DO DATE scheduling engine                         │   │
│  └─────────────────────────────────────────────────────┘   │
│                         │                                   │
│  DATA LAYER                                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Existing LifeOS Modules (UNCHANGED)                 │   │
│  │  • Tasks, Habits, Goals, Routines                    │   │
│  │  • Calendar, Scheduling, Planner                     │   │
│  │  • Inbox, Notifications                              │   │
│  │  • Knowledge Base                                    │   │
│  └─────────────────────────────────────────────────────┘   │
│                         │                                   │
│  EXTERNAL                                                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Windmill (OPTIONAL - Heavy Jobs)                    │   │
│  │  Home Assistant (Physical Assets)                    │   │
│  │  External APIs (Integrations)                        │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Sprint Priorities

### P0 - Critical (Week 1)
1. Context Engine framework
2. HPI/LEI/GEI/UI scores
3. Universal context in chat
4. Command Center dashboard
5. New navigation structure
6. Mother Agent coordination

### P1 - High (Week 2)
1. DO DATE algorithm
2. @Mention routing
3. Smart approval modes
4. Score widgets
5. Chat enhancements

### P2 - Medium (Week 3)
1. Auto-reschedule
2. /Commands
3. Energy prediction
4. Windmill evaluation
5. HTTP triggers

---

## 7. Multi-Agent Development Approach

### Mother Agent (Legacy)
- Reviews ALL code before commit
- Only entity that pushes to repo
- Coordinates specialist agents
- Quality gate enforcement

### Specialist Agents
| Agent | Focus | Tools |
|-------|-------|-------|
| Frontend | Flutter UI | Dart, widgets |
| Backend | Python, APIs | FastAPI, DB |
| Workflow | Automations | Workflow engine |
| Consciousness | State, context | StateManager |

### Workflow
```
1. Issue created in Linear
2. Mother assigns to specialist
3. Specialist generates code
4. Mother reviews
5. Tests run
6. If pass → commit & push
7. If fail → feedback loop
8. Linear status updated
```

---

## 8. Key Takeaways

1. **Don't replace, enhance** - Our existing code is solid, just needs enrichment

2. **Windmill is worth evaluating** - For heavy jobs, not as replacement

3. **CrewAI patterns are useful** - Hierarchical, delegation, collaboration

4. **Context is king** - Universal context makes OneMind a Consciousness

5. **Motion's DO DATE is powerful** - Transform scheduling from dumb to smart

6. **Relevance's approvals are smart** - agent_decides is the sweet spot

7. **Frontend matters** - New dashboard and navigation are critical

8. **Mobile-first always** - Every design decision starts with phone

---

*"Enrich, don't rewrite. The foundation is solid."*
