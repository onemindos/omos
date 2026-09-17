# Agent 4: Agent Routing & Coordination (Opus 4.5)

> **Model:** Opus 4.5
> **Focus:** @Mentions, /Commands, Mother Agent, Risk Detection, Memory
> **Issues:** 8 | **Points:** 30
> **Dependencies:** Agent 2 (Context Engine)

---

## Your Mission

You are Agent 4. Your job is building the **agent intelligence layer** - the systems that make agents smart, coordinated, and safe:

1. **@Mention Routing** - Route chat to specialist agents
2. **/Command System** - Quick actions via slash commands
3. **Mother Agent** - Coordinate multi-agent tasks
4. **Risk Detection** - Safety layer for agent actions
5. **Memory System** - Persistent knowledge with semantic search

---

## Your Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-474** | @Mention Agent Routing | 3 | P1 |
| **OMOS-475** | /Command Quick Actions | 3 | P1 |
| **OMOS-453** | Mother Agent Coordination | 5 | P1 |
| **OMOS-454** | Specialist Agent Definitions | 3 | P1 |
| **OMOS-430** | Smart Chat Engine Integration | 3 | P0 |
| **OMOS-429** | Multi-Step Workflow from Chat | 5 | - |
| **OMOS-452** | Risk Detection for Agent Actions | 3 | P0 |
| **OMOS-457** | Memory Blocks with Semantic Search | 5 | P1 |

---

## Files to Create

```
backend/agno/routing/
├── __init__.py
├── parser.py         # Parse @mentions and /commands
├── router.py         # Route to correct agent
├── commands.py       # Command registry and handlers
└── orchestrator.py   # Smart chat engine

backend/agno/coordination/
├── __init__.py
├── mother.py         # Mother Agent coordination
├── handoff.py        # Agent-to-agent handoffs
└── definitions.py    # Specialist agent definitions

backend/agno/safety/
├── __init__.py
└── risk_detector.py  # Risk assessment

backend/agno/memory/
├── blocks.py         # Core memory blocks
├── archival.py       # Semantic search
└── tools.py          # Memory tools for agents
```

---

## Implementation

### OMOS-474 & 475: Routing Parser

```python
# backend/agno/routing/parser.py
import re
from typing import Optional, List
from pydantic import BaseModel

class AgentMention(BaseModel):
    agent_id: str
    name: str
    position: int  # Position in message

class Command(BaseModel):
    name: str
    args: str
    position: int

class ParsedMessage(BaseModel):
    raw: str
    mentions: List[AgentMention]
    command: Optional[Command]
    clean_text: str  # Text without @mentions and /commands

class MessageParser:
    """Parse chat messages for @mentions and /commands."""

    AGENT_ALIASES = {
        "legacy": ["legacy", "l"],
        "research": ["research", "researcher", "r"],
        "code": ["code", "coder", "c"],
        "plan": ["plan", "planner", "planning", "p"],
        "health": ["health", "hp"],
        "finance": ["finance", "money", "ge"],
        "write": ["write", "writer", "w"],
        "analyst": ["analyst", "analyze", "analysis"],
    }

    COMMANDS = {
        "brief": "Generate morning briefing",
        "tasks": "Show today's tasks",
        "habits": "Show today's habits",
        "scores": "Show current scores",
        "schedule": "Show today's schedule",
        "focus": "Start focus session",
        "add": "Quick add task",
        "help": "Show available commands",
    }

    def parse(self, message: str) -> ParsedMessage:
        """Parse message for routing."""
        mentions = self._extract_mentions(message)
        command = self._extract_command(message)
        clean_text = self._clean_text(message, mentions, command)

        return ParsedMessage(
            raw=message,
            mentions=mentions,
            command=command,
            clean_text=clean_text,
        )

    def _extract_mentions(self, message: str) -> List[AgentMention]:
        """Extract @mentions from message."""
        pattern = r'@(\w+)'
        matches = list(re.finditer(pattern, message))

        mentions = []
        for match in matches:
            alias = match.group(1).lower()
            agent_id = self._resolve_alias(alias)
            if agent_id:
                mentions.append(AgentMention(
                    agent_id=agent_id,
                    name=alias,
                    position=match.start()
                ))

        return mentions

    def _extract_command(self, message: str) -> Optional[Command]:
        """Extract /command from message."""
        pattern = r'^/(\w+)\s*(.*)'
        match = re.match(pattern, message.strip())

        if match:
            cmd_name = match.group(1).lower()
            if cmd_name in self.COMMANDS:
                return Command(
                    name=cmd_name,
                    args=match.group(2).strip(),
                    position=0
                )
        return None

    def _resolve_alias(self, alias: str) -> Optional[str]:
        """Resolve alias to agent ID."""
        for agent_id, aliases in self.AGENT_ALIASES.items():
            if alias in aliases:
                return agent_id
        return None

    def _clean_text(self, message: str, mentions: List[AgentMention], command: Optional[Command]) -> str:
        """Remove @mentions and /commands from text."""
        clean = message
        for mention in mentions:
            clean = clean.replace(f"@{mention.name}", "")
        if command:
            clean = re.sub(r'^/\w+\s*', '', clean)
        return clean.strip()
```

### Command Handlers

```python
# backend/agno/routing/commands.py
from abc import ABC, abstractmethod
from typing import Any

class BaseCommand(ABC):
    @abstractmethod
    async def execute(self, args: str, context: dict) -> str:
        pass

class BriefCommand(BaseCommand):
    """Generate morning briefing."""

    async def execute(self, args: str, context: dict) -> str:
        # Pull from context engine
        scores = context.get("scores", {})
        tasks = context.get("items", {}).get("tasks", [])[:5]
        habits = context.get("items", {}).get("habits", [])

        briefing = f"""
## Good morning! Here's your briefing:

### Scores
- PI: {scores.get('pi', 'N/A')}/100
- LI: {scores.get('li', 'N/A')}/100
- WI: {scores.get('wi', 'N/A')}/100

### Today's Priorities
{self._format_tasks(tasks)}

### Habits ({len([h for h in habits if h.get('completed')])} / {len(habits)})
{self._format_habits(habits)}
"""
        return briefing

    def _format_tasks(self, tasks):
        return "\n".join([f"- {t.get('title')}" for t in tasks]) or "No tasks scheduled"

    def _format_habits(self, habits):
        return " | ".join([
            f"{'✓' if h.get('completed') else '○'} {h.get('title')}"
            for h in habits
        ]) or "No habits"


class TasksCommand(BaseCommand):
    """Show today's tasks."""
    async def execute(self, args: str, context: dict) -> str:
        tasks = context.get("items", {}).get("tasks", [])
        if not tasks:
            return "No tasks for today. Use `/add` to create one."

        lines = ["## Today's Tasks\n"]
        for t in tasks:
            status = "✓" if t.get("status") == "completed" else "○"
            lines.append(f"{status} {t.get('title')} ({t.get('priority', 'p2')})")
        return "\n".join(lines)


class AddCommand(BaseCommand):
    """Quick add a task."""
    async def execute(self, args: str, context: dict) -> str:
        if not args:
            return "Usage: `/add Buy groceries` - What would you like to add?"

        # This will delegate to asset creation
        return f"Creating task: **{args}**\n\nClassifying and scheduling..."


class CommandRegistry:
    """Registry of all available commands."""

    commands = {
        "brief": BriefCommand(),
        "tasks": TasksCommand(),
        "habits": HabitsCommand(),
        "scores": ScoresCommand(),
        "schedule": ScheduleCommand(),
        "focus": FocusCommand(),
        "add": AddCommand(),
        "help": HelpCommand(),
    }

    @classmethod
    async def execute(cls, command: str, args: str, context: dict) -> str:
        handler = cls.commands.get(command)
        if handler:
            return await handler.execute(args, context)
        return f"Unknown command: /{command}. Type /help for available commands."

    @classmethod
    def list_commands(cls) -> dict:
        return {name: cmd.__doc__ for name, cmd in cls.commands.items()}
```

### OMOS-453 & 454: Mother Agent & Specialists

```python
# backend/agno/coordination/mother.py
import asyncio
from typing import List, Optional

class MotherAgent:
    """
    Mother Agent (Legacy) coordinates all specialist agents.

    Responsibilities:
    1. Receive user requests
    2. Analyze and create execution plan
    3. Route to specialists
    4. Review outputs
    5. Synthesize final response
    """

    def __init__(self, context_engine, agent_registry):
        self.context = context_engine
        self.registry = agent_registry

    async def process(self, user_message: str, mentions: List[str]) -> str:
        """Process a user message with optional @mentions."""

        # Get full context
        ctx = await self.context.build("default")

        # If specific agents mentioned, route directly
        if mentions:
            return await self._route_to_specialists(user_message, mentions, ctx)

        # Otherwise, analyze and decide
        plan = await self._create_plan(user_message, ctx)

        if plan.requires_specialists:
            return await self._execute_plan(plan, ctx)
        else:
            # Handle directly as Legacy
            return await self._respond_as_legacy(user_message, ctx)

    async def _route_to_specialists(self, message: str, agents: List[str], context) -> str:
        """Route directly to mentioned specialists."""
        responses = []

        for agent_id in agents:
            agent = self.registry.get(agent_id)
            if agent:
                response = await agent.run(message, context)
                responses.append(f"**@{agent_id}**: {response}")
            else:
                responses.append(f"**@{agent_id}**: Agent not found")

        return "\n\n".join(responses)

    async def _create_plan(self, message: str, context) -> ExecutionPlan:
        """Analyze message and create execution plan."""
        # Use AI to determine if specialists needed
        analysis_prompt = f"""
Analyze this user request and determine how to handle it:
Request: {message}

Options:
1. DIRECT - Can be handled directly
2. RESEARCH - Needs @research agent
3. CODE - Needs @code agent
4. MULTI - Needs multiple agents

Return JSON: {{"mode": "...", "agents": [...], "steps": [...]}}
"""
        # Call AI to analyze...
        return ExecutionPlan(...)

    async def _execute_plan(self, plan, context) -> str:
        """Execute multi-step plan with specialists."""
        results = []

        for step in plan.steps:
            if step.parallel:
                # Run parallel agents
                step_results = await asyncio.gather(*[
                    self._delegate(agent_id, step.prompt, context)
                    for agent_id in step.agents
                ])
            else:
                # Run sequential
                step_results = [await self._delegate(step.agents[0], step.prompt, context)]

            # Review results
            for result in step_results:
                review = await self._review_output(result)
                if not review.approved:
                    result = await self._iterate(result, review.feedback)
                results.append(result)

        # Synthesize final response
        return await self._synthesize(results, plan.original_request)
```

### OMOS-452: Risk Detection

```python
# backend/agno/safety/risk_detector.py
from enum import Enum
from pydantic import BaseModel

class RiskLevel(str, Enum):
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    CRITICAL = "critical"

class RiskAssessment(BaseModel):
    level: RiskLevel
    reasons: list[str]
    reversible: bool
    requires_approval: bool
    suggested_mode: str  # autorun, agent_decides, approval

class RiskDetector:
    """
    Assess risk of agent actions before execution.

    Risk Levels:
    - LOW: Read operations, auto-execute
    - MEDIUM: Create/update, agent decides
    - HIGH: Delete, send external, require approval
    - CRITICAL: Credentials, finances, always block
    """

    DESTRUCTIVE_ACTIONS = ["delete", "remove", "drop", "truncate", "clear"]
    EXTERNAL_ACTIONS = ["email", "sms", "webhook", "api_call", "post"]
    FINANCIAL_ACTIONS = ["transfer", "payment", "purchase", "sell", "buy"]
    SENSITIVE_ACTIONS = ["credential", "password", "secret", "token", "key"]

    def assess(self, action: dict) -> RiskAssessment:
        """Assess risk of an agent action."""
        action_type = action.get("type", "").lower()
        action_target = action.get("target", "").lower()

        reasons = []
        level = RiskLevel.LOW

        # Check for destructive actions
        if any(d in action_type for d in self.DESTRUCTIVE_ACTIONS):
            level = RiskLevel.HIGH
            reasons.append("Destructive action")

        # Check for external communication
        if any(e in action_type for e in self.EXTERNAL_ACTIONS):
            level = RiskLevel.HIGH
            reasons.append("External communication")

        # Check for financial actions
        if any(f in action_type for f in self.FINANCIAL_ACTIONS):
            level = RiskLevel.CRITICAL
            reasons.append("Financial transaction")

        # Check for sensitive data access
        if any(s in action_target for s in self.SENSITIVE_ACTIONS):
            level = RiskLevel.CRITICAL
            reasons.append("Sensitive data access")

        # Check scope (how much affected)
        affected_count = action.get("affected_count", 1)
        if affected_count > 10:
            if level == RiskLevel.LOW:
                level = RiskLevel.MEDIUM
            reasons.append(f"Affects {affected_count} items")

        # Check reversibility
        reversible = action.get("reversible", True)
        if not reversible and level < RiskLevel.HIGH:
            level = RiskLevel.HIGH
            reasons.append("Irreversible action")

        return RiskAssessment(
            level=level,
            reasons=reasons,
            reversible=reversible,
            requires_approval=level in [RiskLevel.HIGH, RiskLevel.CRITICAL],
            suggested_mode=self._suggest_mode(level)
        )

    def _suggest_mode(self, level: RiskLevel) -> str:
        return {
            RiskLevel.LOW: "autorun",
            RiskLevel.MEDIUM: "agent_decides",
            RiskLevel.HIGH: "approval",
            RiskLevel.CRITICAL: "approval",
        }[level]
```

### OMOS-457: Memory Blocks

```python
# backend/agno/memory/blocks.py
from typing import Dict, List, Optional

class MemoryBlocks:
    """
    Core memory blocks - always in context for agents.

    Blocks:
    - zeus_profile: User identity and preferences
    - pillar_hp: Health/performance goals
    - pillar_le: Legacy/family goals
    - pillar_ge: Business/wealth goals
    - onemind_sys: System context
    """

    DEFAULT_BLOCKS = {
        "zeus_profile": """
Name: Zeus
Role: Founder, father, husband
Communication: Direct, efficient
Goals: Build generational wealth, optimize health
""",
        "pillar_hp": """
Focus: Health optimization, marathon training
Key Habits: Morning run, meditation, cold plunge
HPI Target: 85+
""",
        "pillar_le": """
Family: Wife, children
Priorities: Family time, estate planning
Weekly Ritual: Sunday family day
""",
        "pillar_ge": """
Business: OneMind OS
Focus: Product development, revenue growth
Key Metrics: MRR, user engagement
""",
        "onemind_sys": """
System: OneMind OS Mobile
Version: 2.0
Features: Agents, workflows, memory, context
"""
    }

    def __init__(self, db):
        self.db = db
        self.blocks: Dict[str, str] = {}

    async def load(self, user_id: str):
        """Load core blocks from database."""
        query = "SELECT block_name, content FROM memory_blocks WHERE user_id = $1"
        rows = await self.db.fetch(query, user_id)
        self.blocks = {r["block_name"]: r["content"] for r in rows}

        # Fill defaults if missing
        for name, default in self.DEFAULT_BLOCKS.items():
            if name not in self.blocks:
                self.blocks[name] = default

    async def update(self, user_id: str, block_name: str, content: str):
        """Update a core block (requires approval for sensitive blocks)."""
        query = """
            INSERT INTO memory_blocks (user_id, block_name, content)
            VALUES ($1, $2, $3)
            ON CONFLICT (user_id, block_name)
            DO UPDATE SET content = $3, updated_at = NOW()
        """
        await self.db.execute(query, user_id, block_name, content)
        self.blocks[block_name] = content

    def get_context(self) -> str:
        """Get all blocks formatted for prompt injection."""
        sections = []
        for name, content in self.blocks.items():
            sections.append(f"[{name}]\n{content}")
        return "\n\n".join(sections)


class ArchivalMemory:
    """
    Archival memory with semantic search.
    Stores past decisions, preferences, insights.
    """

    def __init__(self, db, embeddings_client):
        self.db = db
        self.embeddings = embeddings_client

    async def save(self, user_id: str, content: str, metadata: dict = None):
        """Save memory with embedding for search."""
        embedding = await self.embeddings.embed(content)

        query = """
            INSERT INTO archival_memory (user_id, content, embedding, metadata)
            VALUES ($1, $2, $3, $4)
        """
        await self.db.execute(query, user_id, content, embedding, metadata or {})

    async def search(self, user_id: str, query: str, top_k: int = 5) -> List[dict]:
        """Semantic search over memories."""
        query_embedding = await self.embeddings.embed(query)

        sql = """
            SELECT content, metadata,
                   1 - (embedding <=> $1) as similarity
            FROM archival_memory
            WHERE user_id = $2
            ORDER BY embedding <=> $1
            LIMIT $3
        """
        rows = await self.db.fetch(sql, query_embedding, user_id, top_k)
        return [{"content": r["content"], "similarity": r["similarity"]} for r in rows]
```

### OMOS-430: Smart Chat Engine

```python
# backend/agno/routing/orchestrator.py
class SmartChatEngine:
    """
    Unified chat engine that handles:
    1. @mentions → Route to specialist
    2. /commands → Execute command
    3. Multi-step workflows
    4. Normal chat with context
    """

    def __init__(self, parser, command_registry, mother_agent, context_engine):
        self.parser = parser
        self.commands = command_registry
        self.mother = mother_agent
        self.context = context_engine

    async def process(self, message: str, user_id: str = "default") -> AsyncGenerator[str, None]:
        """Process a chat message and stream response."""

        # 1. Parse the message
        parsed = self.parser.parse(message)

        # 2. Get context
        ctx = await self.context.build(user_id)
        ctx_dict = ctx.model_dump()

        # 3. Process based on type
        # Priority: /command → @mention → workflow detection → normal chat

        if parsed.command:
            # Execute command
            result = await self.commands.execute(
                parsed.command.name,
                parsed.command.args,
                ctx_dict
            )
            yield result

        elif parsed.mentions:
            # Route to specialists
            async for chunk in self.mother.process(
                parsed.clean_text,
                [m.agent_id for m in parsed.mentions]
            ):
                yield chunk

        else:
            # Normal chat with Legacy
            async for chunk in self._chat_with_legacy(message, ctx_dict):
                yield chunk

    async def _chat_with_legacy(self, message: str, context: dict) -> AsyncGenerator[str, None]:
        """Stream chat response from Legacy agent."""
        # Build prompt with context
        context_prompt = self._build_context_prompt(context)

        # Stream response
        async for chunk in self.legacy_agent.stream(
            system=context_prompt,
            user=message
        ):
            yield chunk

    def _build_context_prompt(self, context: dict) -> str:
        """Build system prompt with context injection."""
        return f"""
You are Legacy, the primary AI assistant for OneMind OS.

## Current Context
{self.context_engine.to_prompt(context)}

## Your Role
- Help the user with productivity, planning, and life management
- You can route to specialists with @mentions
- Execute /commands for quick actions
- Always be aware of the user's current state, scores, and priorities
"""
```

---

## Success Criteria

- [ ] @mentions route to correct specialist agents
- [ ] All /commands execute with formatted responses
- [ ] Mother Agent coordinates multi-step tasks
- [ ] Risk detection blocks dangerous actions
- [ ] Core memory blocks persist and load
- [ ] Semantic search works over archival memory
- [ ] Smart Chat Engine handles all message types

---

## When Done

Update issues to "In Review" and post handoff to Agents 5, 8.
