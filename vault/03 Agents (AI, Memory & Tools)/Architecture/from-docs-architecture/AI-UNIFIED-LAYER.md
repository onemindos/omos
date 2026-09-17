# OneMind AI Unified Layer Architecture
## How the Platform Becomes AI-Aware

*The technical blueprint for creating a system where AI knows everything*

---

## The Core Question

> "How does the whole platform become AI-aware? Is there an AI unified layer we have to create so it's aware of the whole platform and its content for the context?"

**Answer: Yes.** And you've already built most of it.

---

## What You Already Have

```
EXISTING CONSCIOUSNESS INFRASTRUCTURE
═══════════════════════════════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────────────────────────────┐
│                            CONSCIOUSNESS MODULE                                 │
│                         backend/consciousness/                                  │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                              CORE                                        │   │
│  │                                                                         │   │
│  │  EventBus (NATS)          KVStore (NATS)         ObjectStore (NATS)    │   │
│  │  ═════════════════        ══════════════         ═════════════════      │   │
│  │  • Real-time events       • State storage        • Memory/files         │   │
│  │  • Pub/sub                • Session data         • Long-term store      │   │
│  │  • Event routing          • Hot state            • Embeddings           │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                            │
│                                    ▼                                            │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                           AWARENESS                                      │   │
│  │                                                                         │   │
│  │  StateManager             PatternEngine          TemporalAwareness     │   │
│  │  ══════════════           ═════════════          ══════════════════     │   │
│  │  • Location               • Pattern detect       • Time context         │   │
│  │  • Activity               • Biometrics           • Deadlines            │   │
│  │  • Devices                • Insights             • Calendar             │   │
│  │  • Biometrics             • Anomalies            • Time slots           │   │
│  │  • Context summary        • Correlations                                │   │
│  │                                                                         │   │
│  │  LegacyState.get_context_summary() → "📍 Location: Home                │   │
│  │                                        🎯 Activity: Working (45min)     │   │
│  │                                        💻 Task: Coding OneMind          │   │
│  │                                        ⏰ Next: Team call in 30min"     │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                            │
│                                    ▼                                            │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                           PROCESSING                                     │   │
│  │                                                                         │   │
│  │  Processor                Correlator             Archiver               │   │
│  │  ═════════════            ════════════           ═══════════            │   │
│  │  • Event filter           • Cross-service        • TimescaleDB          │   │
│  │  • Awareness gate         • Entity linking       • Historical           │   │
│  │  • Priority routing       • Pattern detect       • Replay               │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                            │
│                                    ▼                                            │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                     EXECUTION ENGINE                                     │   │
│  │                   backend/consciousness/engine/                          │   │
│  │                                                                         │   │
│  │  • Receives tasks from inbox                                            │   │
│  │  • Classifies (category, urgency, complexity)                          │   │
│  │  • Selects execution target (agent/team/workflow)                      │   │
│  │  • Injects context from inbox item                                     │   │
│  │  • Executes and tracks progress                                        │   │
│  │  • HITL approval integration                                           │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘

ALSO EXISTING:
══════════════

Agno Module (backend/agno/):
├── Agents (create_legacy_agent, specialists)
├── Teams (research, dev, finance)
├── Workflows (steps: agent, parallel, approval, function)
├── Knowledge (LanceDB embeddings, RAG)
├── Memory (session, conversation history)
└── Tools (Composio, MCP, browser, GitHub, etc.)

LifeOS Module (backend/lifeos/):
├── Tasks
├── Habits
├── Goals
├── Calendar
├── Inbox
├── Journal
└── Meals
```

---

## What's Missing: The Unified Context Layer

The **consciousness module** tracks real-time state (location, activity, devices). But it doesn't know about:

1. **Your Items** - Tasks, projects, habits, goals, skills
2. **Your Paths** - HP, LE, GE progress and scores
3. **Your Knowledge Base** - Documents, notes, recipes
4. **Your Intelligence Scores** - HPI, LEI, GEI, UI
5. **Your History** - Past conversations, patterns, learnings

We need a **Unified Context Layer** that aggregates ALL of this into a single context that any AI interface can use.

```
THE MISSING PIECE: UNIFIED CONTEXT LAYER
═══════════════════════════════════════════════════════════════════════════════════

                           DATA SOURCES (What exists)
                           ═════════════════════════
          ┌──────────────────────────────────────────────────────────┐
          │                                                          │
          │  PostgreSQL        LanceDB           Redis              │
          │  ══════════════    ═══════════       ═══════════        │
          │  • Items           • Knowledge       • Hot state        │
          │  • Users           • Embeddings      • Sessions         │
          │  • Scores          • Search          • Cache            │
          │  • History         • Semantic        │                  │
          │                                                          │
          └──────────────────────────────────────────────────────────┘
                                    │
                                    │  CURRENTLY DISCONNECTED
                                    │
                                    ▼
          ┌──────────────────────────────────────────────────────────┐
          │                                                          │
          │              UNIFIED CONTEXT LAYER  ← MISSING            │
          │                                                          │
          │  This layer needs to:                                    │
          │  • Aggregate from ALL sources                           │
          │  • Build semantic understanding                         │
          │  • Calculate intelligence scores                        │
          │  • Generate context for AI prompts                      │
          │  • Update in real-time as data changes                  │
          │                                                          │
          └──────────────────────────────────────────────────────────┘
                                    │
                                    │  PROVIDES CONTEXT TO
                                    │
                                    ▼
          ┌──────────────────────────────────────────────────────────┐
          │                                                          │
          │                    AI INTERFACES                         │
          │                                                          │
          │  Chat       Voice       Agents      Workflows            │
          │  ═══════    ═══════     ════════    ═════════            │
          │  UIO Chat   LiveKit     Specialists  Briefings           │
          │                                                          │
          └──────────────────────────────────────────────────────────┘
```

---

## The Architecture: How AI Becomes Aware

```
COMPLETE AI-AWARE ARCHITECTURE
═══════════════════════════════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────────────────────────────┐
│                                                                                 │
│                             TIER 1: DATA LAYER                                  │
│                         (Sources of Truth)                                      │
│                                                                                 │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐│
│  │  POSTGRESQL    │  │    LANCEDB     │  │     REDIS      │  │     NATS       ││
│  │                │  │                │  │                │  │                ││
│  │ • Items table  │  │ • Embeddings   │  │ • Hot state    │  │ • Events       ││
│  │ • Users        │  │ • Knowledge    │  │ • Session      │  │ • Pub/sub      ││
│  │ • Scores       │  │ • Semantic     │  │ • Cache        │  │ • Real-time    ││
│  │ • History      │  │   search       │  │ • Context      │  │                ││
│  │ • Integrations │  │                │  │   cache        │  │                ││
│  └────────────────┘  └────────────────┘  └────────────────┘  └────────────────┘│
│          │                  │                   │                   │           │
│          └──────────────────┴───────────────────┴───────────────────┘           │
│                                      │                                          │
│                                      ▼                                          │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│                         TIER 2: UNIFIED CONTEXT LAYER                           │
│                    backend/context/ (NEW - TO BE CREATED)                       │
│                                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │                                                                          │  │
│  │                         CONTEXT ENGINE                                   │  │
│  │                       context/engine.py                                  │  │
│  │                                                                          │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐   │  │
│  │  │                      AGGREGATORS                                │   │  │
│  │  │                                                                 │   │  │
│  │  │  StateAggregator     ItemsAggregator     KnowledgeAggregator   │   │  │
│  │  │  ═══════════════     ════════════════    ═══════════════════   │   │  │
│  │  │  Consciousness       Universal items    LanceDB + RAG          │   │  │
│  │  │  state manager       from PostgreSQL     retrieval             │   │  │
│  │  │                                                                 │   │  │
│  │  │  ScoresAggregator    CalendarAggregator  HistoryAggregator     │   │  │
│  │  │  ═══════════════     ══════════════════  ══════════════════    │   │  │
│  │  │  Calculate HPI/      Google Calendar     Conversation          │   │  │
│  │  │  LEI/GEI/UI          integration         history               │   │  │
│  │  └─────────────────────────────────────────────────────────────────┘   │  │
│  │                                │                                        │  │
│  │                                ▼                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐   │  │
│  │  │                    CONTEXT BUILDER                              │   │  │
│  │  │                   context/builder.py                            │   │  │
│  │  │                                                                 │   │  │
│  │  │  build_full_context()     → Complete context for chat          │   │  │
│  │  │  build_agent_context()    → Context for specific agent         │   │  │
│  │  │  build_workflow_context() → Context for workflow step          │   │  │
│  │  │  build_briefing_context() → Context for morning/evening brief  │   │  │
│  │  │  build_scheduling_context()→ Context for DO DATE calculation   │   │  │
│  │  └─────────────────────────────────────────────────────────────────┘   │  │
│  │                                │                                        │  │
│  │                                ▼                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐   │  │
│  │  │                    CONTEXT CACHE                                │   │  │
│  │  │                   context/cache.py                              │   │  │
│  │  │                                                                 │   │  │
│  │  │  • Redis-backed for fast access                                │   │  │
│  │  │  • TTL-based invalidation                                      │   │  │
│  │  │  • Event-driven updates (NATS subscription)                    │   │  │
│  │  │  • Warm cache on startup                                       │   │  │
│  │  └─────────────────────────────────────────────────────────────────┘   │  │
│  │                                                                          │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                      │                                          │
│                                      ▼                                          │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│                           TIER 3: AI INTERFACES                                 │
│                    (All access context through unified layer)                   │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                         │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │   │
│  │  │                         UIO CHAT                                │  │   │
│  │  │                   agno/chat/enhanced.py                         │  │   │
│  │  │                                                                 │  │   │
│  │  │  • Context injection from Context Engine                       │  │   │
│  │  │  • @agent syntax for specialist invocation                     │  │   │
│  │  │  • /commands for quick actions                                 │  │   │
│  │  │  • Multi-step workflow detection and execution                 │  │   │
│  │  │  • Citations from knowledge base                               │  │   │
│  │  │                                                                 │  │   │
│  │  │  FLOW:                                                         │  │   │
│  │  │  User message → Context Engine → Build prompt → LLM → Response │  │   │
│  │  └─────────────────────────────────────────────────────────────────┘  │   │
│  │                                                                         │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐       │   │
│  │  │   VOICE    │  │   AGENTS   │  │ WORKFLOWS  │  │ BRIEFINGS  │       │   │
│  │  │            │  │            │  │            │  │            │       │   │
│  │  │ LiveKit    │  │ Specialists│  │ Multi-step │  │ Morning/   │       │   │
│  │  │ → Context  │  │ → Context  │  │ → Context  │  │ Evening    │       │   │
│  │  │   Engine   │  │   Engine   │  │   Engine   │  │ → Context  │       │   │
│  │  │            │  │            │  │            │  │   Engine   │       │   │
│  │  └────────────┘  └────────────┘  └────────────┘  └────────────┘       │   │
│  │                                                                         │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## The Context Engine: Implementation

### File Structure

```
backend/context/
├── __init__.py              # Exports
├── engine.py                # Main context engine
├── builder.py               # Context builders
├── cache.py                 # Redis-backed cache
├── aggregators/
│   ├── __init__.py
│   ├── state.py             # From consciousness StateManager
│   ├── items.py             # From PostgreSQL items table
│   ├── knowledge.py         # From LanceDB
│   ├── scores.py            # Intelligence scores calculator
│   ├── calendar.py          # From Google Calendar
│   └── history.py           # Conversation history
└── models.py                # Pydantic models
```

### Core Implementation

```python
# backend/context/engine.py

class ContextEngine:
    """
    The unified context engine that makes AI aware of EVERYTHING.

    This is the brain that aggregates all data sources and provides
    context to any AI interface (chat, voice, agents, workflows).
    """

    def __init__(
        self,
        db_pool,              # PostgreSQL
        lancedb,              # Vector search
        redis_client,         # Cache
        state_manager,        # From consciousness
        nats_client,          # For real-time updates
    ):
        self.db = db_pool
        self.lancedb = lancedb
        self.redis = redis_client
        self.state_manager = state_manager
        self.nats = nats_client

        # Aggregators
        self.state = StateAggregator(state_manager)
        self.items = ItemsAggregator(db_pool)
        self.knowledge = KnowledgeAggregator(lancedb)
        self.scores = ScoresAggregator(db_pool)
        self.calendar = CalendarAggregator()
        self.history = HistoryAggregator(db_pool)

        # Cache
        self.cache = ContextCache(redis_client)

    async def get_full_context(self, user_id: str) -> FullContext:
        """
        Get complete context for a user.

        This is the main entry point - returns everything the AI needs.
        """
        # Check cache first
        cached = await self.cache.get(f"context:{user_id}")
        if cached:
            return FullContext.model_validate_json(cached)

        # Aggregate from all sources in parallel
        state, items, scores, calendar, history = await asyncio.gather(
            self.state.get(user_id),
            self.items.get_active(user_id),
            self.scores.calculate_all(user_id),
            self.calendar.get_today(user_id),
            self.history.get_recent(user_id, limit=10),
        )

        context = FullContext(
            user_id=user_id,
            timestamp=datetime.utcnow(),

            # Real-time state (from consciousness)
            state=state,

            # Intelligence scores
            scores=scores,

            # Active items by path
            hp_items=self._filter_by_path(items, "HP"),
            le_items=self._filter_by_path(items, "LE"),
            ge_items=self._filter_by_path(items, "GE"),

            # Today's schedule
            calendar=calendar,

            # Conversation context
            history=history,
        )

        # Cache for quick access
        await self.cache.set(f"context:{user_id}", context.model_dump_json(), ttl=60)

        return context

    async def get_agent_context(
        self,
        user_id: str,
        agent_id: str,
        query: Optional[str] = None,
    ) -> AgentContext:
        """
        Get context for a specific agent.

        Includes relevant knowledge from RAG if query provided.
        """
        full = await self.get_full_context(user_id)

        # Get agent config to know what context it needs
        agent_config = await self._get_agent_config(agent_id)

        # Get relevant knowledge if query provided
        knowledge = []
        if query:
            knowledge = await self.knowledge.search(
                query=query,
                filters=agent_config.knowledge_filters,
                limit=5,
            )

        return AgentContext(
            user_id=user_id,
            agent_id=agent_id,
            state_summary=full.state.get_summary(),
            scores=full.scores,
            relevant_items=self._filter_for_agent(full, agent_config),
            relevant_knowledge=knowledge,
            calendar_summary=full.calendar.get_summary(),
        )

    async def build_prompt_context(
        self,
        user_id: str,
        query: Optional[str] = None,
    ) -> str:
        """
        Build a text context block for injection into prompts.

        This is what gets injected into the system prompt.
        """
        ctx = await self.get_full_context(user_id)

        sections = []

        # User info
        sections.append(f"USER: {ctx.state.user_name}")
        sections.append(f"TIME: {datetime.now().strftime('%A, %B %d, %Y at %I:%M %p')}")

        # Intelligence scores
        sections.append(f"""
INTELLIGENCE SCORES:
- HPI (Holistic Performance): {ctx.scores.hpi} {self._trend(ctx.scores.hpi_trend)}
- LEI (Legacy Evolution): {ctx.scores.lei} {self._trend(ctx.scores.lei_trend)}
- GEI (Generational Endeavor): {ctx.scores.gei} {self._trend(ctx.scores.gei_trend)}
- SO (Sovereign): {ctx.scores.ui}""")

        # Real-time state
        if ctx.state.context_summary:
            sections.append(f"CURRENT STATE:\n{ctx.state.context_summary}")

        # Today's schedule
        if ctx.calendar.events:
            schedule_lines = [f"TODAY'S SCHEDULE ({len(ctx.calendar.events)} events):"]
            for event in ctx.calendar.events[:5]:
                schedule_lines.append(f"- {event.time}: {event.title}")
            sections.append("\n".join(schedule_lines))

        # Active tasks by priority
        high_priority = [i for i in ctx.hp_items + ctx.le_items + ctx.ge_items
                        if i.priority == "high"]
        if high_priority:
            tasks_lines = ["HIGH PRIORITY TASKS:"]
            for item in high_priority[:5]:
                tasks_lines.append(f"- [{item.path}] {item.title}")
            sections.append("\n".join(tasks_lines))

        # Recent context
        if ctx.history:
            recent = ctx.history[0]
            sections.append(f"LAST CONVERSATION: {recent.topic or 'General'} ({recent.time_ago})")

        return "\n\n".join(sections)

    async def subscribe_to_updates(self, user_id: str):
        """
        Subscribe to real-time context updates via NATS.

        When data changes, invalidate cache and notify listeners.
        """
        async def handle_event(msg):
            # Invalidate cache on any relevant event
            await self.cache.invalidate(f"context:{user_id}")

            # Publish context update event
            await self.nats.publish(
                f"context.updated.{user_id}",
                json.dumps({"user_id": user_id, "reason": msg.subject}).encode()
            )

        # Subscribe to relevant event streams
        await self.nats.subscribe(f"events.*.{user_id}", cb=handle_event)
        await self.nats.subscribe(f"items.*.{user_id}", cb=handle_event)
        await self.nats.subscribe(f"state.*.{user_id}", cb=handle_event)
```

### Intelligence Scores Calculator

```python
# backend/context/aggregators/scores.py

class ScoresAggregator:
    """
    Calculates HPI, LEI, GEI, and UI scores.

    Based on the CODEX framework definitions.
    """

    # Score component weights
    HPI_WEIGHTS = {
        "sleep_quality": 0.20,
        "recovery": 0.20,
        "habit_completion": 0.25,
        "training_load": 0.15,
        "mood_average": 0.10,
        "goal_progress": 0.10,
    }

    LEI_WEIGHTS = {
        "family_rituals": 0.25,
        "home_maintenance": 0.20,
        "quality_time": 0.20,
        "wealth_posture": 0.15,
        "communication": 0.10,
        "legacy_projects": 0.10,
    }

    GEI_WEIGHTS = {
        "runway": 0.25,
        "free_cash_flow": 0.25,
        "revenue_growth": 0.20,
        "client_health": 0.15,
        "legal_hygiene": 0.10,
        "capital_efficiency": 0.05,
    }

    async def calculate_all(self, user_id: str) -> IntelligenceScores:
        """Calculate all intelligence scores for a user."""

        # Get component data
        hpi_data = await self._get_hp_components(user_id)
        lei_data = await self._get_le_components(user_id)
        gei_data = await self._get_ge_components(user_id)

        # Calculate weighted scores
        hpi = self._calculate_weighted(hpi_data, self.HPI_WEIGHTS)
        lei = self._calculate_weighted(lei_data, self.LEI_WEIGHTS)
        gei = self._calculate_weighted(gei_data, self.GEI_WEIGHTS)

        # Unified score is average
        ui = (hpi + lei + gei) / 3

        # Get trends (compare to previous period)
        hpi_trend, lei_trend, gei_trend = await self._get_trends(user_id)

        return IntelligenceScores(
            hpi=round(hpi, 1),
            lei=round(lei, 1),
            gei=round(gei, 1),
            ui=round(ui, 1),
            hpi_trend=hpi_trend,
            lei_trend=lei_trend,
            gei_trend=gei_trend,
            components={
                "hpi": hpi_data,
                "lei": lei_data,
                "gei": gei_data,
            },
            calculated_at=datetime.utcnow(),
        )

    async def _get_hp_components(self, user_id: str) -> Dict[str, float]:
        """Get HP score components."""
        components = {}

        # Sleep quality (from Oura or manual)
        sleep = await self._get_sleep_score(user_id)
        components["sleep_quality"] = sleep

        # Recovery (from Oura or manual)
        recovery = await self._get_recovery_score(user_id)
        components["recovery"] = recovery

        # Habit completion (from habits module)
        habits = await self._get_habit_completion(user_id, path="HP")
        components["habit_completion"] = habits

        # Training load (from activity/fitness)
        training = await self._get_training_score(user_id)
        components["training_load"] = training

        # Mood average (from journal/mood logs)
        mood = await self._get_mood_score(user_id)
        components["mood_average"] = mood

        # Goal progress (HP goals)
        goals = await self._get_goal_progress(user_id, path="HP")
        components["goal_progress"] = goals

        return components
```

---

## How the Chat Uses Context

```python
# backend/agno/chat/enhanced.py

class EnhancedChatEngine:
    """
    Chat engine that uses the unified context layer.
    """

    def __init__(self, context_engine: ContextEngine):
        self.context = context_engine

    async def process_message(
        self,
        user_id: str,
        message: str,
        session_id: str,
    ) -> ChatResponse:
        """
        Process a chat message with full context.
        """

        # 1. Build prompt context
        prompt_context = await self.context.build_prompt_context(user_id, query=message)

        # 2. Check for @agent mentions
        if "@" in message:
            return await self._handle_agent_mention(user_id, message, prompt_context)

        # 3. Check for /commands
        if message.startswith("/"):
            return await self._handle_command(user_id, message, prompt_context)

        # 4. Normal chat with full context
        full_prompt = f"""
You are Legacy, the AI partner in the OneMind human optimization system.

CURRENT CONTEXT:
================
{prompt_context}

CONVERSATION:
================
User: {message}

Respond with full awareness of the user's context. Be concise but helpful.
If suggesting actions, consider their schedule and energy level.
Always reference specific data points when relevant.
"""

        # Get response from LLM
        response = await self._get_llm_response(full_prompt, session_id)

        return ChatResponse(
            content=response,
            context_used=prompt_context,
            session_id=session_id,
        )
```

---

## Technology Stack Summary

```
TECHNOLOGY STACK
═══════════════════════════════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────────────────────────────┐
│                                                                                 │
│  LAYER              │ TECHNOLOGY        │ PURPOSE                              │
│  ══════════════════════════════════════════════════════════════════════════    │
│                     │                   │                                      │
│  Data Storage       │ PostgreSQL        │ Items, users, scores, history        │
│                     │ LanceDB           │ Knowledge base, embeddings, RAG      │
│                     │ Redis             │ Hot state, sessions, context cache   │
│                     │ NATS              │ Events, real-time pub/sub            │
│                     │                   │                                      │
│  Consciousness      │ StateManager      │ Location, activity, devices, bio     │
│                     │ PatternEngine     │ Pattern detection, insights          │
│                     │ TemporalAwareness │ Time, deadlines, calendar            │
│                     │                   │                                      │
│  Context Layer      │ ContextEngine     │ Aggregate all sources (NEW)          │
│  (TO BUILD)         │ Aggregators       │ State, items, knowledge, scores      │
│                     │ ContextCache      │ Redis-backed, event-invalidated      │
│                     │                   │                                      │
│  AI/Agents          │ Agno              │ Agent framework                      │
│                     │ LangChain/Phidata │ LLM orchestration                    │
│                     │ Composio          │ 2000+ app integrations               │
│                     │ MCP               │ Local tool servers                   │
│                     │                   │                                      │
│  Interfaces         │ FastAPI           │ REST API                             │
│                     │ WebSocket         │ Real-time chat                       │
│                     │ LiveKit           │ Voice interface                      │
│                     │ Flutter           │ Mobile/web frontend                  │
│                     │                   │                                      │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Implementation Plan for 24-Hour Venture

### Phase 1: Context Engine Foundation (4 hours)

```
FILES TO CREATE:
----------------
backend/context/__init__.py
backend/context/engine.py          # Main context engine
backend/context/builder.py         # Context builders
backend/context/cache.py           # Redis cache
backend/context/models.py          # Pydantic models

backend/context/aggregators/__init__.py
backend/context/aggregators/state.py      # Wrap consciousness StateManager
backend/context/aggregators/items.py      # PostgreSQL items queries
backend/context/aggregators/knowledge.py  # LanceDB RAG
backend/context/aggregators/scores.py     # Calculate HPI/LEI/GEI/UI
backend/context/aggregators/calendar.py   # Google Calendar
backend/context/aggregators/history.py    # Conversation history
```

### Phase 2: Enhanced Chat (4 hours)

```
FILES TO CREATE:
----------------
backend/agno/chat/__init__.py
backend/agno/chat/enhanced.py      # Main chat engine
backend/agno/chat/commands.py      # /command handlers
backend/agno/chat/agents.py        # @agent router
backend/agno/chat/workflows.py     # Multi-step workflow detection
backend/agno/chat/citations.py     # Citation tracking
```

### Phase 3: Universal Items (4 hours)

```
DATABASE MIGRATION:
-------------------
migrations/20260129_universal_items.sql
- Create item_type, path_type, workforce_type enums
- Create unified items table
- Create intelligence_scores table
- Create items_search index

FILES TO CREATE:
----------------
backend/lifeos/items/__init__.py
backend/lifeos/items/service.py    # CRUD operations
backend/lifeos/items/models.py     # Pydantic models
backend/lifeos/items/api.py        # REST endpoints
```

### Phase 4: Scheduling Engine (4 hours)

```
FILES TO CREATE:
----------------
backend/platform/scheduling/__init__.py
backend/platform/scheduling/intelligent.py  # DO DATE algorithm
backend/platform/scheduling/calendar.py     # Google Calendar sync
backend/platform/scheduling/energy.py       # Energy pattern analysis
backend/platform/scheduling/briefings.py    # Morning/evening briefs
```

### Phase 5: Agent Enhancements (4 hours)

```
FILES TO MODIFY:
----------------
backend/agno/agents/models.py      # Add persona fields
backend/agno/workflows/workflow/steps/approval.py  # Smart approval modes

FILES TO CREATE:
----------------
backend/agno/agents/personas.py    # Persona management
backend/agno/agents/invent.py      # Natural language agent creation
backend/agno/agents/queue.py       # Task queue management
```

### Phase 6: Frontend Integration (4 hours)

```
FILES TO CREATE:
----------------
frontend/lib/context/context_service.dart      # Context API client
frontend/lib/context/context_provider.dart     # Riverpod provider

frontend/lib/shared/widgets/chat_enhanced.dart # Enhanced chat widget
frontend/lib/shared/widgets/scores_card.dart   # Intelligence scores display
frontend/lib/shared/widgets/workforce_card.dart # Agent status cards
```

---

## GitHub Issues for Project

```
EPIC: AI Unified Context Layer
------------------------------

Issue #1: [Backend] Create Context Engine
- Create backend/context/ module
- Implement aggregators for all data sources
- Add Redis caching with NATS invalidation
- Expose context through API

Issue #2: [Backend] Implement Intelligence Scores
- Create ScoresAggregator
- Implement HPI, LEI, GEI, UI calculations
- Add trend analysis
- Store historical scores in PostgreSQL

Issue #3: [Backend] Enhanced Chat Engine
- Create backend/agno/chat/ module
- Implement @agent syntax parsing
- Implement /command handling
- Add multi-step workflow detection
- Integrate with Context Engine

Issue #4: [Database] Universal Items Migration
- Create item_type, path_type enums
- Create unified items table
- Create indexes for efficient queries
- Migrate existing data

Issue #5: [Backend] Intelligent Scheduling
- Implement DO DATE algorithm
- Add Google Calendar integration
- Add energy pattern analysis
- Create morning/evening briefing workflows

Issue #6: [Frontend] Context Integration
- Create context service/provider
- Update chat UI for @agents and /commands
- Add scores display to dashboard
- Add workforce status cards

Issue #7: [Backend] Agent Enhancements
- Add persona fields to agent config
- Implement "Invent" feature
- Add smart approval modes
- Create task queue management
```

---

## Conclusion

You've already built **90% of the infrastructure**. The consciousness module, state manager, pattern engine, execution engine - it's all there.

What's missing is the **glue layer** that:
1. **Aggregates** all data sources into one context
2. **Provides** that context to all AI interfaces
3. **Calculates** the intelligence scores
4. **Caches** for performance

The 24-hour venture should focus on:
1. **Context Engine** (the glue)
2. **Enhanced Chat** (the primary interface)
3. **Intelligence Scores** (the differentiation)
4. **Universal Items** (the data model)

Everything else builds on top of this foundation.

**Human + AI = ONE MIND** - now let's make it aware.
