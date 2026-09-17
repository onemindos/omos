# OneMind OS: Unified Physical-Digital Architecture
## Brainstorming the Complete System

> **Vision**: A unified, decentralized, autonomous system managing both physical (robots, IoT, garden) and digital (business, family, life) domains with shared context, offline-capable, and human-in-loop coordination.

---

## Current State Analysis

### What We Have ✅

```
OneMind OS (Current)
├── Agno Agents (Digital Specialists)
│   ├── Assistant, Researcher, Coder
│   ├── Chat-first, reactive
│   └── 50+ tools (APIs, webhooks)
│
├── OpenMind Heartbeat Loop (1 Hz)
│   ├── SENSE → FUSE → CORTEX → ACTION
│   └── Proactive awareness (not just reactive)
│
├── PostgreSQL + pgvector
│   ├── Agent configs, teams
│   ├── Memories (vector embeddings)
│   └── Sessions, knowledge bases
│
├── NATS JetStream
│   ├── Event bus (nervous system)
│   ├── MQTT bridge (IoT devices)
│   └── JetStream (persistent streams)
│
├── OneIntel (Internal PM)
│   ├── Unified items table
│   ├── Projects, tasks, docs, habits
│   └── Replaces ClickUp
│
└── Flutter Frontend
    ├── 16+ screens
    └── WebSocket real-time
```

### What We Need 🎯

```
Physical Layer (Not Yet Implemented)
├── OM1 Integration
│   ├── Robot control (Unitree, TurtleBot, etc.)
│   ├── Physical sensors (camera, LiDAR)
│   └── ROS message passing
│
├── FarmBot Integration (Eden)
│   ├── Garden automation
│   └── SolarPunk vision
│
├── Home Assistant Integration
│   ├── Smart home devices
│   └── Environmental sensors
│
└── Offline/Decentralized
    ├── FABRIC protocol (OM1)
    └── Off-grid autonomy
```

---

## The Fragmentation Problem

### Current Architecture (Fragmented)

```
┌──────────────────────────────────────────────────────────────┐
│                    FRAGMENTED SYSTEMS                         │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  Agno (Digital Agents)    OM1 (Physical Robots)              │
│  ────────────────────     ───────────────────                │
│  • Chat-first             • Robot-first                       │
│  • Reactive               • Real-time control                 │
│  • No heartbeat           • Heartbeat built-in                │
│  • API tools              • ROS actions                       │
│  • PostgreSQL             • Separate DB?                      │
│                                                               │
│  ❌ NO SHARED CONTEXT                                         │
│  ❌ NO UNIFIED ORCHESTRATION                                  │
│  ❌ NO EVENT ROUTING                                          │
│                                                               │
│  OneIntel (PM)            Home Assistant         FarmBot      │
│  ────────────────         ──────────────         ────────     │
│  • SQLite/PostgreSQL      • Own database         • API only   │
│  • Agno tool              • MQTT events          • HTTP API   │
│                                                               │
│  ❌ EACH SYSTEM ISOLATED                                      │
│  ❌ AGENTS CAN'T COORDINATE                                   │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

**Key Issues:**
1. ❌ Agno agents don't know about physical world
2. ❌ Physical robots don't have digital context
3. ❌ Events routed manually (no intelligent routing)
4. ❌ No shared memory/learning
5. ❌ Offline mode impossible (cloud dependencies)

---

## Architectural Options Analysis

### Option 1: Agno as Primary ⚠️

```
Agno (Orchestration)
├── Digital Agents (native)
├── Physical Proxies
│   └── OM1 Tool (wrapper)
│       ├── Robot control via API
│       └── Limited physical awareness
└── Heartbeat (custom implementation)
```

**Pros:**
- ✅ Strong digital agent capabilities
- ✅ Mature LLM integration
- ✅ Easy to extend with tools

**Cons:**
- ❌ Not designed for physical world
- ❌ No native robot control
- ❌ Reactive by default (need custom heartbeat)
- ❌ No decentralized/offline mode
- ❌ No multi-robot coordination

**Verdict:** ⚠️ **Good for digital-only, insufficient for robotics**

---

### Option 2: OM1 as Primary ⚠️

```
OM1 (Orchestration)
├── Physical Agents (native)
│   ├── Eden (FarmBot)
│   ├── Mobile robots
│   └── Home automation
├── Digital Proxies
│   └── Virtual agents (limited)
└── FABRIC (decentralized trust)
```

**Pros:**
- ✅ Built for robotics
- ✅ Decentralized (FABRIC protocol)
- ✅ Offline-first
- ✅ Multi-robot coordination
- ✅ Real-time sensor processing

**Cons:**
- ❌ Weak digital agent capabilities
- ❌ Not chat-optimized
- ❌ Limited LLM integration
- ❌ Fewer tools for digital world
- ❌ Early stage (less mature than Agno)

**Verdict:** ⚠️ **Good for robotics, insufficient for digital life**

---

### Option 3: OneMind OS as Unified Orchestration Layer ✅

```
┌──────────────────────────────────────────────────────────────────┐
│              ONEMIND OS (Unified Orchestration)                   │
│              ═══════════════════════════════════                  │
│                                                                   │
│  HEARTBEAT LOOP (1 Hz Proactive + Reactive)                      │
│  ───────────────────────────────────────────                      │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────┐     │
│  │  1. SENSE (Unified Sensory Input)                       │     │
│  │     ─────────────────────────────                       │     │
│  │                                                          │     │
│  │     Digital Sensors         Physical Sensors            │     │
│  │     ───────────────         ────────────────            │     │
│  │     • Webhooks              • Cameras (OM1)             │     │
│  │     • Calendar              • LiDAR (OM1)               │     │
│  │     • Email                 • IMU (OM1)                 │     │
│  │     • ClickUp → OneIntel    • Soil sensors (FarmBot)    │     │
│  │     • GitHub                • Home Assistant            │     │
│  │                                                          │     │
│  │     ALL → NATS → Digital Sensory Data Object            │     │
│  └─────────────────────────────────────────────────────────┘     │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────┐     │
│  │  2. FUSE (Unified Context Creation)                     │     │
│  │     ──────────────────────────────                      │     │
│  │                                                          │     │
│  │     Unified Memory Store (PostgreSQL + pgvector)        │     │
│  │     ───────────────────────────────────────────         │     │
│  │     • Agent memories (digital + physical)               │     │
│  │     • World model (digital + physical state)            │     │
│  │     • Goals (unified across domains)                    │     │
│  │     • Events history (TimescaleDB)                      │     │
│  │     • OneIntel (projects, tasks, habits)                │     │
│  │                                                          │     │
│  │     Context Object:                                     │     │
│  │     ──────────────                                      │     │
│  │     {                                                   │     │
│  │       "digital_state": {...},  // Business, family      │     │
│  │       "physical_state": {...}, // Garden, home, robots  │     │
│  │       "unified_goals": [...],  // Cross-domain          │     │
│  │       "recent_events": [...],  // Last 60 seconds       │     │
│  │       "agent_memories": [...], // Shared learning       │     │
│  │     }                                                   │     │
│  └─────────────────────────────────────────────────────────┘     │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────┐     │
│  │  3. CORTEX (Multi-Agent Decision Making)                │     │
│  │     ───────────────────────────────────                 │     │
│  │                                                          │     │
│  │     Agent Router (Intelligent Event Assignment)         │     │
│  │     ────────────────────────────────────────            │     │
│  │                                                          │     │
│  │     Event → Router → {                                  │     │
│  │       Digital Agent (Agno),                             │     │
│  │       Physical Agent (OM1),                             │     │
│  │       Human (OneIntel approval)                         │     │
│  │     }                                                   │     │
│  │                                                          │     │
│  │     Agno Agents           OM1 Agents                    │     │
│  │     ────────────          ──────────                    │     │
│  │     • Assistant            • Eden (Garden)              │     │
│  │     • Researcher           • HomeOS (Automation)        │     │
│  │     • Coder                • Scout (Mobile robot)       │     │
│  │     • Analyst                                           │     │
│  │                                                          │     │
│  │     All agents share unified context                    │     │
│  │     All publish decisions to NATS                       │     │
│  │     All learn from each other                           │     │
│  └─────────────────────────────────────────────────────────┘     │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────┐     │
│  │  4. ACTION (Unified Execution)                          │     │
│  │     ───────────────────────                             │     │
│  │                                                          │     │
│  │     Digital Tools         Physical Tools                │     │
│  │     ─────────────         ──────────────                │     │
│  │     • API calls            • ROS commands (OM1)         │     │
│  │     • Webhooks             • MQTT (IoT)                 │     │
│  │     • OneIntel ops         • FarmBot API                │     │
│  │     • Email/Slack          • Home Assistant             │     │
│  │                                                          │     │
│  │     ALL → Publish results to NATS → Update context      │     │
│  └─────────────────────────────────────────────────────────┘     │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────┐     │
│  │  5. REMEMBER (Unified Learning)                         │     │
│  │     ────────────────────────                            │     │
│  │                                                          │     │
│  │     • Digital learnings → Shared memory                 │     │
│  │     • Physical learnings → Shared memory                │     │
│  │     • Cross-domain insights                             │     │
│  │     • Evolve together                                   │     │
│  └─────────────────────────────────────────────────────────┘     │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                    INFRASTRUCTURE LAYER                           │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  NATS JetStream (Nervous System)                                 │
│  ─────────────────────────────────                               │
│  • events.digital.*    (Calendar, email, webhooks)               │
│  • events.physical.*   (Sensors, robots, IoT)                    │
│  • cortex.decisions.*  (Agent plans)                             │
│  • action.digital.*    (API calls, tool executions)              │
│  • action.physical.*   (Robot commands, actuators)               │
│  • human.approvals.*   (OneIntel workflows)                      │
│                                                                   │
│  Unified Memory (PostgreSQL + pgvector + TimescaleDB)            │
│  ───────────────────────────────────────────────────────         │
│  • Agent memories (embeddings, vector search)                    │
│  • World state (digital + physical)                              │
│  • Goals & tasks (OneIntel)                                      │
│  • Event history (compressed, 90 days)                           │
│                                                                   │
│  Decentralization (FABRIC Protocol - OM1)                        │
│  ────────────────────────────────────────                        │
│  • Offline operation                                             │
│  • Robot-to-robot trust                                          │
│  • Distributed identity                                          │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
```

**Pros:**
- ✅ **Best of both worlds** (Agno digital + OM1 physical)
- ✅ **Unified context** (all agents share memory)
- ✅ **Intelligent routing** (events → right agent or human)
- ✅ **Proactive + Reactive** (heartbeat loop)
- ✅ **Decentralized** (FABRIC for offline/autonomous)
- ✅ **Scalable** (add agents as needed)
- ✅ **Human-in-loop** (OneIntel approval workflows)

**Cons:**
- ⚠️ More complex to build (custom orchestration)
- ⚠️ Integration work (Agno + OM1)

**Verdict:** ✅ **RECOMMENDED - Unified, extensible, future-proof**

---

## Recommended Architecture: OneMind OS as Orchestration Platform

### Why OneMind OS Should Be Primary

| Reason | Explanation |
|--------|-------------|
| **1. Already Hybrid** | You've already built the heartbeat loop - this IS the OpenMind architecture |
| **2. Best of Both** | Use Agno for digital, OM1 for physical (specialized agents) |
| **3. Unified Context** | Single memory store, shared learning |
| **4. Event Routing** | Intelligent assignment to digital/physical/human |
| **5. Extensible** | Add new agent types without rewriting |
| **6. Open Source** | Both Agno and OM1 are open source Python |
| **7. Decentralized** | FABRIC protocol for offline autonomy |

---

## Detailed Integration Strategy

### 1. Agent Layer Integration

```python
# backend/agents/unified_agent_system.py

class UnifiedAgentSystem:
    """
    Orchestrates both Agno (digital) and OM1 (physical) agents
    with shared context and coordinated decision-making.
    """

    def __init__(self):
        # Digital agents (Agno)
        self.digital_agents = {
            "assistant": AgnoAgent(name="Assistant", tools=[...]),
            "researcher": AgnoAgent(name="Researcher", tools=[...]),
            "coder": AgnoAgent(name="Coder", tools=[...]),
        }

        # Physical agents (OM1)
        self.physical_agents = {
            "eden": OM1Agent(name="Eden", robot="farmbot"),
            "homeos": OM1Agent(name="HomeOS", robot="home_assistant"),
            "scout": OM1Agent(name="Scout", robot="unitree_go2"),
        }

        # Unified context
        self.memory = UnifiedMemoryStore()
        self.nats = NATS()

    async def route_event(self, event: NormalizedEvent) -> AgentDecision:
        """
        Intelligent routing: Digital vs Physical vs Human
        """
        # Get unified context
        context = await self.memory.get_unified_context()

        # Determine agent type needed
        if self._is_physical_task(event, context):
            agent = self._select_physical_agent(event)
            return await agent.process(event, context)

        elif self._is_digital_task(event, context):
            agent = self._select_digital_agent(event)
            return await agent.process(event, context)

        else:
            # Requires human decision
            return await self._request_human_approval(event, context)

    def _is_physical_task(self, event, context) -> bool:
        """
        Determine if task requires physical action.

        Examples:
        - "Water the garden" → Physical (Eden)
        - "Turn off lights" → Physical (HomeOS)
        - "Check front door camera" → Physical (HomeOS)
        """
        physical_keywords = [
            "robot", "garden", "water", "plant", "grow",
            "light", "door", "camera", "temperature",
            "move", "navigate", "grab", "pick up"
        ]
        return any(kw in event.title.lower() for kw in physical_keywords)

    def _is_digital_task(self, event, context) -> bool:
        """
        Determine if task requires digital action.

        Examples:
        - "Send email to team" → Digital (Assistant)
        - "Research quantum computing" → Digital (Researcher)
        - "Fix bug in auth.py" → Digital (Coder)
        """
        return not self._is_physical_task(event, context)
```

---

### 2. Unified Memory Architecture

```python
# backend/memory/unified_store.py

class UnifiedMemoryStore:
    """
    Shared memory accessible by all agents (digital + physical).

    Uses PostgreSQL + pgvector for:
    - Agent memories (embeddings)
    - World state (digital + physical)
    - Goals (unified context)
    - Event history
    """

    def __init__(self, db):
        self.db = db

    async def store_memory(
        self,
        agent_id: str,
        agent_type: str,  # "digital" or "physical"
        memory: str,
        embedding: List[float],
        metadata: Dict
    ):
        """
        Store a memory accessible to all agents.
        """
        await self.db.execute("""
            INSERT INTO unified_memories (
                agent_id,
                agent_type,
                content,
                embedding,
                metadata,
                created_at
            ) VALUES ($1, $2, $3, $4, $5, NOW())
        """, agent_id, agent_type, memory, embedding, metadata)

    async def get_relevant_memories(
        self,
        query_embedding: List[float],
        agent_type: Optional[str] = None,  # Filter by type if needed
        limit: int = 10
    ) -> List[Memory]:
        """
        Vector search across ALL agent memories.

        This allows:
        - Digital agents to learn from physical experiences
        - Physical agents to use digital context
        """
        if agent_type:
            result = await self.db.fetch("""
                SELECT * FROM unified_memories
                WHERE agent_type = $1
                ORDER BY embedding <=> $2
                LIMIT $3
            """, agent_type, query_embedding, limit)
        else:
            result = await self.db.fetch("""
                SELECT * FROM unified_memories
                ORDER BY embedding <=> $1
                LIMIT $2
            """, query_embedding, limit)

        return [Memory.from_row(row) for row in result]

    async def get_unified_context(self) -> UnifiedContext:
        """
        Build complete context combining digital + physical state.
        """
        return UnifiedContext(
            digital_state=await self._get_digital_state(),
            physical_state=await self._get_physical_state(),
            goals=await self._get_active_goals(),
            recent_events=await self._get_recent_events(),
            agent_memories=await self._get_all_memories()
        )

    async def _get_digital_state(self) -> Dict:
        """
        Digital world state from various sources.
        """
        return {
            "calendar": await self._get_upcoming_events(),
            "email": await self._get_unread_emails(),
            "tasks": await self._get_oneintel_tasks(),
            "projects": await self._get_active_projects(),
        }

    async def _get_physical_state(self) -> Dict:
        """
        Physical world state from OM1 agents.
        """
        return {
            "garden": await self._get_garden_state(),  # From Eden
            "home": await self._get_home_state(),      # From HomeOS
            "robots": await self._get_robot_states(),  # From Scout, etc.
        }
```

---

### 3. Event Routing System

```python
# backend/routing/event_router.py

class IntelligentEventRouter:
    """
    Routes events to the most appropriate handler:
    - Digital agent (Agno)
    - Physical agent (OM1)
    - Human (OneIntel approval workflow)
    """

    def __init__(self, unified_system: UnifiedAgentSystem):
        self.system = unified_system
        self.nats = NATS()

    async def route(self, event: NormalizedEvent) -> RoutingDecision:
        """
        Determine where event should go.
        """
        # Get unified context
        context = await self.system.memory.get_unified_context()

        # Analyze event characteristics
        analysis = await self._analyze_event(event, context)

        # Routing logic
        if analysis.requires_human_decision:
            return await self._route_to_human(event, analysis)

        elif analysis.is_autonomous_capable:
            if analysis.domain == "physical":
                return await self._route_to_physical_agent(event, analysis)
            else:
                return await self._route_to_digital_agent(event, analysis)

        else:
            # Complex cross-domain task
            return await self._route_to_multi_agent_team(event, analysis)

    async def _analyze_event(self, event, context) -> EventAnalysis:
        """
        Analyze event to determine routing.

        Factors:
        - Domain (digital vs physical)
        - Complexity (simple vs complex)
        - Risk level (autonomous vs needs approval)
        - Required capabilities
        """
        return EventAnalysis(
            domain=self._determine_domain(event),
            complexity=self._assess_complexity(event, context),
            risk_level=self._assess_risk(event, context),
            required_capabilities=self._identify_capabilities(event),
            recommended_agent=self._suggest_agent(event, context)
        )

    async def _route_to_human(self, event, analysis) -> RoutingDecision:
        """
        Create OneIntel approval task.
        """
        # Create task in OneIntel
        task = await oneintel(
            action="create",
            type="task",
            title=f"Approval needed: {event.title}",
            content=f"""
            Event: {event.to_dict()}

            Analysis: {analysis.to_dict()}

            Recommended action: {analysis.recommended_action}

            Please review and approve/reject.
            """,
            priority=analysis.risk_level,
            tags=["approval", "human-in-loop"]
        )

        # Publish to NATS
        await self.nats.publish(
            "human.approvals.pending",
            {"task_id": task.id, "event": event.to_dict()}
        )

        return RoutingDecision(
            route="human",
            handler="oneintel",
            task_id=task.id
        )
```

---

### 4. Database Strategy

```sql
-- Unified database schema supporting both digital and physical

-- ============================================================
-- UNIFIED MEMORIES (Shared by all agents)
-- ============================================================
CREATE TABLE unified_memories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id TEXT NOT NULL,
    agent_type TEXT NOT NULL,  -- 'digital' or 'physical'
    content TEXT NOT NULL,
    embedding vector(1536),    -- pgvector for similarity search
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    -- Vector index for fast similarity search
    INDEX idx_embedding ON unified_memories USING ivfflat (embedding vector_cosine_ops)
);

-- ============================================================
-- WORLD STATE (Digital + Physical)
-- ============================================================
CREATE TABLE world_state (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain TEXT NOT NULL,      -- 'digital' or 'physical'
    entity_type TEXT NOT NULL, -- 'calendar_event', 'robot', 'sensor', etc.
    entity_id TEXT NOT NULL,
    state JSONB NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- GOALS (Unified across domains)
-- ============================================================
CREATE TABLE unified_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    goal_type TEXT NOT NULL,   -- 'digital', 'physical', 'hybrid'
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'active',
    assigned_agents TEXT[],    -- Can be digital + physical
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- EVENTS (Time-series, from NATS archival)
-- ============================================================
CREATE TABLE events (
    timestamp TIMESTAMPTZ NOT NULL,
    source TEXT NOT NULL,      -- 'digital' or 'physical'
    type TEXT NOT NULL,
    data JSONB NOT NULL,
    metadata JSONB
);

-- TimescaleDB hypertable for efficient time-series queries
SELECT create_hypertable('events', 'timestamp');

-- ============================================================
-- ONEINTEL (Project Management)
-- ============================================================
CREATE TABLE oneintel_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL,        -- 'project', 'task', 'doc', 'habit'
    title TEXT NOT NULL,
    content TEXT,
    parent_id UUID REFERENCES oneintel_items(id),
    status TEXT DEFAULT 'todo',
    domain TEXT,               -- 'digital', 'physical', 'hybrid'
    assigned_agent TEXT,       -- Digital or physical agent
    due_date DATE,
    priority INT,
    tags TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

### 5. Offline/Decentralized Strategy

```python
# backend/decentralized/fabric_integration.py

class FabricIntegration:
    """
    Integrate OM1's FABRIC protocol for decentralized operation.

    Enables:
    - Offline autonomous operation
    - Robot-to-robot trust
    - Distributed identity
    """

    def __init__(self):
        self.fabric = FabricClient()
        self.local_storage = LocalCache()

    async def sync_state(self):
        """
        Sync state between online and offline modes.

        When online:
        - Push local changes to FABRIC
        - Pull updates from network

        When offline:
        - Use local cache
        - Queue changes for sync
        """
        if await self.is_online():
            # Push queued changes
            await self._push_pending_changes()

            # Pull updates
            await self._pull_network_updates()
        else:
            # Offline mode: use local cache
            logger.info("Offline mode: using local cache")

    async def _push_pending_changes(self):
        """
        Push local changes to FABRIC network.
        """
        pending = await self.local_storage.get_pending_changes()

        for change in pending:
            await self.fabric.publish_change(change)

    async def _pull_network_updates(self):
        """
        Pull updates from FABRIC network.
        """
        updates = await self.fabric.get_updates(since=self.last_sync)

        for update in updates:
            await self.local_storage.apply_update(update)
```

---

### 6. SolarPunk Vision: Eden Garden Robot

```python
# backend/agents/physical/eden.py

class EdenAgent(OM1Agent):
    """
    Physical agent for garden automation using FarmBot.

    Capabilities:
    - Plant care (watering, monitoring)
    - Soil analysis
    - Pest detection (via camera)
    - Harvest planning
    - Sustainable practices
    """

    def __init__(self):
        super().__init__(name="Eden", robot_type="farmbot")
        self.farmbot = FarmBotAPI()
        self.vision = VisionModel()

    async def process_event(self, event: NormalizedEvent, context: UnifiedContext):
        """
        Process garden-related events.

        Examples:
        - "Water the tomatoes" → Execute watering sequence
        - "Check garden health" → Analyze plants via camera
        - "Plant new seeds" → Plan planting sequence
        """
        if "water" in event.title.lower():
            return await self._water_plants(event, context)

        elif "check" in event.title.lower():
            return await self._check_garden_health(context)

        elif "plant" in event.title.lower():
            return await self._plan_planting(event, context)

    async def _water_plants(self, event, context):
        """
        Execute watering sequence via FarmBot.
        """
        # Get plant locations from database
        plants = await self.get_plants_needing_water(context)

        # Execute watering
        for plant in plants:
            await self.farmbot.move_to(plant.x, plant.y)
            await self.farmbot.water(duration=plant.water_duration)

        # Store memory
        await self.memory.store_memory(
            agent_id="eden",
            agent_type="physical",
            memory=f"Watered {len(plants)} plants",
            embedding=await self.embed(f"garden watering {plants}"),
            metadata={"plants": [p.dict() for p in plants]}
        )

        # Publish to NATS
        await self.nats.publish(
            "action.physical.garden.watered",
            {"plants": [p.dict() for p in plants]}
        )
```

---

## Implementation Roadmap

### Phase 1: Unified Core (Month 1-2)

```
[ ] 1. Event Router Implementation
    [ ] Intelligent routing logic
    [ ] Digital vs Physical vs Human detection
    [ ] Risk assessment

[ ] 2. Unified Memory Store
    [ ] PostgreSQL + pgvector schema
    [ ] Cross-agent memory sharing
    [ ] Vector similarity search

[ ] 3. Agent Integration
    [ ] Agno agents → OneMind OS orchestration
    [ ] OM1 stub agents (placeholder)
    [ ] Shared context interface

[ ] 4. OneIntel Integration
    [ ] Approval workflows
    [ ] Task assignment (digital/physical)
    [ ] Human-in-loop UI
```

### Phase 2: Physical Layer (Month 3-4)

```
[ ] 1. OM1 Integration
    [ ] Install OM1 SDK
    [ ] Create physical agent wrapper
    [ ] Robot control interface

[ ] 2. Eden (FarmBot) Agent
    [ ] FarmBot API integration
    [ ] Plant database
    [ ] Watering automation
    [ ] Vision-based health checks

[ ] 3. HomeOS Agent
    [ ] Home Assistant integration
    [ ] Device control
    [ ] Sensor monitoring
    [ ] Automation rules

[ ] 4. Unified Context
    [ ] Digital + Physical state fusion
    [ ] Cross-domain decision making
```

### Phase 3: Decentralization (Month 5-6)

```
[ ] 1. FABRIC Protocol
    [ ] OM1 FABRIC integration
    [ ] Identity management
    [ ] Trust network

[ ] 2. Offline Mode
    [ ] Local caching
    [ ] Sync queue
    [ ] Autonomous operation

[ ] 3. Multi-Robot Coordination
    [ ] Robot-to-robot communication
    [ ] Distributed task allocation
    [ ] Consensus mechanisms
```

---

## Answers to Your Questions

### 1. What should be the main platform for orchestration and why?

**Answer:** ✅ **OneMind OS**

**Why:**
- ✅ You've already built the heartbeat loop (OpenMind architecture)
- ✅ Best of both worlds: Agno (digital) + OM1 (physical)
- ✅ Unified context accessible to all agents
- ✅ Intelligent event routing
- ✅ Extensible (add new agent types)
- ✅ Decentralized capable (FABRIC protocol)
- ✅ Human-in-loop (OneIntel integration)

---

### 2. Should Agno or OM1 be primary?

**Answer:** ✅ **Neither - OneMind OS orchestrates both**

**Architecture:**
```
OneMind OS (Orchestration Platform)
├── Agno Layer (Digital Specialists)
│   └── Chat, research, coding, analysis
│
├── OM1 Layer (Physical Specialists)
│   └── Robots, IoT, garden, home
│
└── Unified Context (Shared Memory)
    └── All agents learn together
```

Think of it like:
- **OneMind OS** = Conductor
- **Agno Agents** = Digital musicians
- **OM1 Agents** = Physical musicians
- **Unified Context** = Musical score (shared)

---

### 3. How to handle the heartbeat with Agno (reactive by default)?

**Answer:** ✅ **OneMind OS heartbeat loop wraps Agno agents**

Your current implementation already does this! The heartbeat loop:
1. SENSE (digital + physical inputs)
2. FUSE (unified context)
3. CORTEX (call Agno/OM1 agents with context)
4. ACTION (execute digital/physical tools)

This makes Agno agents **proactive** because the heartbeat feeds them context every second, even without user input.

---

### 4. What databases make sense for long-term unified agents?

**Answer:** ✅ **PostgreSQL + pgvector + TimescaleDB + SQLite (OneIntel)**

```
PostgreSQL + pgvector:
- Unified memories (all agents)
- World state (digital + physical)
- Goals (cross-domain)
- Agent configs, teams

TimescaleDB (PostgreSQL extension):
- Event history (NATS archival)
- Sensor data (physical agents)
- Analytics (cross-domain insights)

SQLite (OneIntel):
- Projects, tasks, habits
- Local-first, sync to PostgreSQL

NATS JetStream:
- Real-time event bus
- Replay capability
- MQTT bridge (IoT)
```

---

### 5. Is there something better than OM1 for physical/digital hybrid?

**Answer:** ⚠️ **No single platform - OM1 is best for robotics, Agno best for digital**

**Alternatives considered:**

| Platform | Digital | Physical | Verdict |
|----------|---------|----------|---------|
| **Agno** | ⭐⭐⭐⭐⭐ | ⭐⭐ | Great digital, weak physical |
| **OM1** | ⭐⭐ | ⭐⭐⭐⭐⭐ | Great physical, weak digital |
| **ROS** | ⭐ | ⭐⭐⭐⭐ | Robotics only, no AI agents |
| **LangChain** | ⭐⭐⭐ | ⭐ | Digital only, no robotics |
| **CrewAI** | ⭐⭐⭐⭐ | ⭐ | Good multi-agent, no physical |

**Recommendation:** Use **OneMind OS to orchestrate both Agno + OM1**. This gives you best-in-class for each domain.

---

### 6. How to handle offline/online hybrid?

**Answer:** ✅ **FABRIC protocol (from OM1) + Local-first architecture**

```python
# backend/hybrid_mode.py

class HybridMode:
    """
    Seamless online/offline operation.
    """

    async def operate(self):
        if self.is_online():
            # Online mode: Full capabilities
            - Use cloud LLMs
            - Sync to FABRIC network
            - Access external APIs

        else:
            # Offline mode: Autonomous
            - Use local LLMs (llama.cpp)
            - Use cached memories
            - Execute pre-planned tasks
            - Queue changes for sync
```

---

## Final Recommendation

### Build OneMind OS as Unified Orchestration Platform

```
┌──────────────────────────────────────────────────────────┐
│            ONEMIND OS UNIFIED ARCHITECTURE                │
│            (SolarPunk Autonomous System)                  │
├──────────────────────────────────────────────────────────┤
│                                                           │
│  Heartbeat Loop (1 Hz)                                    │
│  ──────────────────────                                   │
│  SENSE → FUSE → CORTEX → ACTION → REMEMBER               │
│                                                           │
│  Agent Layers                                             │
│  ────────────                                             │
│  • Agno Agents (Digital: chat, research, code)            │
│  • OM1 Agents (Physical: Eden, HomeOS, robots)            │
│  • Human-in-Loop (OneIntel approvals)                     │
│                                                           │
│  Unified Context                                          │
│  ───────────────                                          │
│  • PostgreSQL + pgvector (memories)                       │
│  • TimescaleDB (events history)                           │
│  • OneIntel (projects, tasks)                             │
│  • NATS (real-time events)                                │
│                                                           │
│  Decentralization                                         │
│  ───────────────────                                      │
│  • FABRIC protocol (offline mode)                         │
│  • Local-first (SQLite cache)                             │
│  • Autonomous operation                                   │
│                                                           │
└──────────────────────────────────────────────────────────┘

Vision:
- Manage digital life (business, family) with Agno agents
- Manage physical world (garden, home) with OM1 agents
- Unified context (all agents learn together)
- Decentralized (off-grid autonomous)
- SolarPunk (Eden garden robot, sustainable)
- Human-in-loop (OneIntel project management)
```

---

## Next Steps

Want me to:
1. ✅ Create detailed integration spec (Agno + OM1 + OneMind OS)?
2. ✅ Design the Event Router implementation?
3. ✅ Design the Unified Memory architecture?
4. ✅ Create Eden (FarmBot) agent spec?
5. ✅ Write offline/online hybrid mode implementation?

Let me know which part you want to dive into next! 🌱🤖

---

_Last updated: 2026-02-11_


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
