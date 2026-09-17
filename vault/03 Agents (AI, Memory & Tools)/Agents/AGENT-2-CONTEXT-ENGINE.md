# Agent 2: Context Engine (Opus 4.5)

> **Model:** Opus 4.5
> **Focus:** Context Engine - Aggregating ALL user data for AI
> **Issues:** 9 | **Points:** 22
> **Dependencies:** Agent 1 (Universal Assets)

---

## Your Mission

You are Agent 2 of the 8-agent swarm. Your job is building the **Context Engine** - the system that aggregates ALL user data into a unified context object that gets injected into every AI interaction. This is what makes OneMind a "consciousness" rather than a chatbot.

---

## Your Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-416** | Context Engine Framework | 3 | P0-Critical |
| **OMOS-417** | State Aggregator | 2 | - |
| **OMOS-418** | Items Aggregator | 2 | - |
| **OMOS-419** | Scores Aggregator | 2 | - |
| **OMOS-420** | Calendar Aggregator | 2 | - |
| **OMOS-421** | History Aggregator | 2 | - |
| **OMOS-451** | Assets Aggregator | 3 | P1 |
| **OMOS-456** | Biometrics Aggregator (Oura) | 3 | P1 |
| **OMOS-479** | Photo Capture Processing | 3 | P1 |

---

## Dependencies

**Wait for Agent 1** to complete:
- `backend/lifeos/assets/models.py` - Asset schema
- `backend/lifeos/assets/service.py` - AssetService

Once Agent 1 posts handoff, you can begin.

---

## Existing Code Context

### Consciousness Module (Build On This)
```
backend/consciousness/
├── bus.py              # NATS JetStream event bus
├── kv_store.py         # NATS KV Store for state
├── awareness/
│   ├── state.py        # StateManager - existing state tracking
│   └── device_presence.py
└── api.py              # Awareness endpoints
```

### Integrations (Use These)
```
backend/integrations/
├── homeassistant/      # Home Assistant client
│   └── client.py       # HAClient for device states
└── oura/               # Oura integration (if exists)
```

---

## Implementation Requirements

### Files to Create

```
backend/consciousness/context/
├── __init__.py
├── engine.py           # Main ContextEngine class
├── models.py           # UnifiedContext, MinimalContext
├── api.py              # /api/v1/context endpoints
├── cache.py            # Caching layer
└── aggregators/
    ├── __init__.py
    ├── state.py        # Location, energy, focus, device
    ├── items.py        # Tasks, habits, goals from assets
    ├── scores.py       # PI/LI/WI (wraps score service)
    ├── calendar.py     # Calendar events
    ├── history.py      # Recent activity
    ├── assets.py       # Physical/digital assets
    └── biometrics.py   # Oura sleep, readiness, HRV

backend/agno/capture/
└── photo_processor.py  # GPT-4 Vision analysis
```

### OMOS-416: Context Engine Framework

**Create:** `backend/consciousness/context/engine.py`

```python
import asyncio
from datetime import datetime
from typing import Optional
from .models import UnifiedContext, MinimalContext
from .cache import ContextCache
from .aggregators import (
    StateAggregator,
    ItemsAggregator,
    ScoresAggregator,
    CalendarAggregator,
    HistoryAggregator,
    AssetsAggregator,
    BiometricsAggregator,
)

class ContextEngine:
    """
    Central engine that aggregates ALL user data into unified context.
    This is injected into every AI interaction.

    Performance target: build() < 500ms
    """

    def __init__(self, db, oura_client=None, ha_client=None):
        self.db = db
        self.cache = ContextCache(ttl_seconds=60)

        # Initialize aggregators
        self.aggregators = {
            'state': StateAggregator(db, ha_client),
            'items': ItemsAggregator(db),
            'scores': ScoresAggregator(db),
            'calendar': CalendarAggregator(db),
            'history': HistoryAggregator(db),
            'assets': AssetsAggregator(db, ha_client),
            'biometrics': BiometricsAggregator(oura_client),
        }

    async def build(self, user_id: str) -> UnifiedContext:
        """
        Build full context for AI operations.

        Aggregates data in parallel for speed.
        Target: < 500ms total.
        """
        # Check cache first
        cached = self.cache.get(user_id)
        if cached:
            return cached

        # Aggregate all data in parallel
        start = datetime.utcnow()

        results = await asyncio.gather(
            self.aggregators['state'].aggregate(user_id),
            self.aggregators['items'].aggregate(user_id),
            self.aggregators['scores'].aggregate(user_id),
            self.aggregators['calendar'].aggregate(user_id),
            self.aggregators['history'].aggregate(user_id),
            self.aggregators['assets'].aggregate(user_id),
            self.aggregators['biometrics'].aggregate(user_id),
            return_exceptions=True  # Don't fail if one aggregator fails
        )

        # Handle any failures gracefully
        state, items, scores, calendar, history, assets, biometrics = [
            r if not isinstance(r, Exception) else None
            for r in results
        ]

        context = UnifiedContext(
            user_id=user_id,
            timestamp=datetime.utcnow(),
            build_time_ms=(datetime.utcnow() - start).total_seconds() * 1000,
            state=state,
            items=items,
            scores=scores,
            calendar=calendar,
            history=history,
            assets=assets,
            biometrics=biometrics,
        )

        # Cache the result
        self.cache.set(user_id, context)

        return context

    async def build_minimal(self, user_id: str) -> MinimalContext:
        """
        Build minimal context for quick operations.
        Only includes state and today's items.

        Target: < 100ms
        """
        state, items = await asyncio.gather(
            self.aggregators['state'].aggregate(user_id),
            self.aggregators['items'].aggregate_today(user_id),
        )

        return MinimalContext(
            user_id=user_id,
            timestamp=datetime.utcnow(),
            state=state,
            today_items=items,
        )

    async def invalidate_cache(self, user_id: str):
        """Invalidate cached context when data changes."""
        self.cache.delete(user_id)

    def to_prompt(self, context: UnifiedContext) -> str:
        """
        Convert context to prompt format for injection into AI.
        """
        sections = []

        if context.state:
            sections.append(f"""## Current State
- Location: {context.state.location}
- Energy Level: {context.state.energy_level}
- Active Focus: {context.state.focus_mode}
- Current Path: {context.state.active_pillar}
""")

        if context.scores:
            sections.append(f"""## Intelligence Scores
- PI (Performance): {context.scores.pi}/100 ({context.scores.pi_trend})
- LI (Legacy): {context.scores.li}/100 ({context.scores.li_trend})
- WI (Wealth): {context.scores.wi}/100 ({context.scores.wi_trend})
- Unified: {context.scores.unified}/100
""")

        if context.items:
            tasks = [f"  - {t.title} ({t.priority})" for t in context.items.tasks[:5]]
            habits = [f"  - {h.title} ({'✓' if h.completed_today else '○'})" for h in context.items.habits[:5]]
            sections.append(f"""## Today's Items
Tasks ({len(context.items.tasks)}):
{chr(10).join(tasks)}

Habits ({context.items.habits_completed}/{len(context.items.habits)}):
{chr(10).join(habits)}
""")

        if context.calendar:
            events = [f"  - {e.time}: {e.title}" for e in context.calendar.today[:5]]
            sections.append(f"""## Calendar
{chr(10).join(events)}
""")

        if context.biometrics:
            sections.append(f"""## Biometrics (Oura)
- Sleep Score: {context.biometrics.sleep_score}/100
- Readiness: {context.biometrics.readiness_score}/100
- HRV: {context.biometrics.hrv} ms
""")

        return "\n".join(sections)
```

**Create:** `backend/consciousness/context/models.py`

```python
from datetime import datetime
from typing import Optional, List
from pydantic import BaseModel

class StateContext(BaseModel):
    location: str = "unknown"  # home, office, mobile
    energy_level: str = "medium"  # high, medium, low
    focus_mode: Optional[str] = None  # deep_work, meetings, etc
    active_pillar: Optional[str] = None  # hp, le, ge
    device: Optional[str] = None  # mobile, desktop

class ItemsContext(BaseModel):
    tasks: List[dict] = []
    habits: List[dict] = []
    goals: List[dict] = []
    habits_completed: int = 0
    tasks_completed_today: int = 0

class ScoresContext(BaseModel):
    pi: float = 0
    li: float = 0
    wi: float = 0
    unified: float = 0
    pi_trend: str = "stable"
    li_trend: str = "stable"
    wi_trend: str = "stable"

class CalendarContext(BaseModel):
    today: List[dict] = []
    tomorrow: List[dict] = []
    available_slots: List[dict] = []

class HistoryContext(BaseModel):
    recent_completions: List[dict] = []
    recent_conversations: List[str] = []
    productivity_patterns: dict = {}

class AssetsContext(BaseModel):
    home_status: dict = {}  # Home Assistant devices
    wearables_status: dict = {}
    financial_summary: dict = {}
    vehicle_status: dict = {}

class BiometricsContext(BaseModel):
    sleep_score: Optional[int] = None
    sleep_duration: Optional[float] = None
    readiness_score: Optional[int] = None
    hrv: Optional[int] = None
    resting_hr: Optional[int] = None
    energy_level: str = "medium"
    recommendations: List[str] = []

class UnifiedContext(BaseModel):
    """Complete context for AI operations."""
    user_id: str
    timestamp: datetime
    build_time_ms: float = 0

    state: Optional[StateContext] = None
    items: Optional[ItemsContext] = None
    scores: Optional[ScoresContext] = None
    calendar: Optional[CalendarContext] = None
    history: Optional[HistoryContext] = None
    assets: Optional[AssetsContext] = None
    biometrics: Optional[BiometricsContext] = None

class MinimalContext(BaseModel):
    """Minimal context for quick operations."""
    user_id: str
    timestamp: datetime
    state: Optional[StateContext] = None
    today_items: Optional[ItemsContext] = None
```

### Aggregator Implementations

**Create:** `backend/consciousness/context/aggregators/state.py`

```python
from ..models import StateContext

class StateAggregator:
    def __init__(self, db, ha_client=None):
        self.db = db
        self.ha = ha_client

    async def aggregate(self, user_id: str) -> StateContext:
        # Get location from Home Assistant
        location = "unknown"
        if self.ha:
            try:
                presence = await self.ha.get_entity_state("person.zeus")
                location = presence.get("state", "unknown")
            except:
                pass

        # Get energy from recent biometrics or user input
        energy = await self._get_energy_level(user_id)

        # Get focus mode if active
        focus = await self._get_active_focus(user_id)

        return StateContext(
            location=location,
            energy_level=energy,
            focus_mode=focus,
            active_pillar=await self._get_active_pillar(user_id),
        )

    async def _get_energy_level(self, user_id: str) -> str:
        # Check Oura readiness if available, else default
        query = "SELECT energy_level FROM user_state WHERE user_id = $1"
        row = await self.db.fetchrow(query, user_id)
        return row["energy_level"] if row else "medium"

    async def _get_active_focus(self, user_id: str) -> Optional[str]:
        query = """
            SELECT mode FROM focus_sessions
            WHERE user_id = $1 AND ended_at IS NULL
            ORDER BY started_at DESC LIMIT 1
        """
        row = await self.db.fetchrow(query, user_id)
        return row["mode"] if row else None

    async def _get_active_pillar(self, user_id: str) -> Optional[str]:
        # Infer from current tasks or explicit setting
        return None
```

**Create:** `backend/consciousness/context/aggregators/items.py`

```python
from datetime import datetime
from ..models import ItemsContext
# Import from Agent 1's work
from backend.lifeos.assets.service import AssetService
from backend.lifeos.assets.models import AssetType, AssetFilters

class ItemsAggregator:
    def __init__(self, db):
        self.db = db
        self.assets = AssetService(db)

    async def aggregate(self, user_id: str) -> ItemsContext:
        today = await self.aggregate_today(user_id)

        # Get active goals
        goals_filter = AssetFilters(asset_type=AssetType.GOAL, status="in_progress")
        goals = await self.assets.list(goals_filter)

        today.goals = [g.model_dump() for g in goals[:10]]
        return today

    async def aggregate_today(self, user_id: str) -> ItemsContext:
        """Get just today's items for minimal context."""
        today_assets = await self.assets.get_today()

        tasks = [a for a in today_assets if a.asset_type == AssetType.TASK]
        habits = [a for a in today_assets if a.asset_type == AssetType.HABIT]

        return ItemsContext(
            tasks=[t.model_dump() for t in tasks],
            habits=[h.model_dump() for h in habits],
            habits_completed=len([h for h in habits if h.status == "completed"]),
            tasks_completed_today=len([t for t in tasks if t.status == "completed"]),
        )
```

**Create:** `backend/consciousness/context/aggregators/biometrics.py`

```python
from typing import Optional
from ..models import BiometricsContext

class BiometricsAggregator:
    """Aggregates biometric data from Oura Ring."""

    def __init__(self, oura_client=None):
        self.oura = oura_client

    async def aggregate(self, user_id: str) -> Optional[BiometricsContext]:
        if not self.oura:
            return None

        try:
            # Get today's sleep and readiness
            sleep = await self.oura.get_sleep_data()
            readiness = await self.oura.get_readiness_data()

            return BiometricsContext(
                sleep_score=sleep.get("score"),
                sleep_duration=sleep.get("total_sleep_duration", 0) / 3600,  # hours
                readiness_score=readiness.get("score"),
                hrv=readiness.get("hrv_balance"),
                resting_hr=readiness.get("resting_heart_rate"),
                energy_level=self._calculate_energy(readiness.get("score", 50)),
                recommendations=self._get_recommendations(sleep, readiness),
            )
        except Exception as e:
            print(f"Oura aggregation failed: {e}")
            return None

    def _calculate_energy(self, readiness_score: int) -> str:
        if readiness_score >= 85:
            return "high"
        elif readiness_score >= 65:
            return "medium"
        return "low"

    def _get_recommendations(self, sleep: dict, readiness: dict) -> list:
        recs = []
        if sleep.get("score", 100) < 70:
            recs.append("Sleep score low - consider earlier bedtime")
        if readiness.get("score", 100) < 70:
            recs.append("Recovery low - recommend lighter day")
        return recs
```

### OMOS-479: Photo Capture Processing

**Create:** `backend/agno/capture/photo_processor.py`

```python
import base64
from openai import AsyncOpenAI
from typing import Optional

class PhotoProcessor:
    """Process captured photos with GPT-4 Vision for asset creation."""

    def __init__(self, openai_api_key: str):
        self.client = AsyncOpenAI(api_key=openai_api_key)

    async def analyze(self, image_bytes: bytes, user_context: str = "") -> dict:
        """
        Analyze image and extract:
        - Text (OCR)
        - Context (what's in the image)
        - Suggested asset type
        - Suggested title
        """
        base64_image = base64.b64encode(image_bytes).decode('utf-8')

        response = await self.client.chat.completions.create(
            model="gpt-4-vision-preview",
            messages=[
                {
                    "role": "system",
                    "content": """You are analyzing images for a personal productivity system.
Extract any text visible (OCR), describe the context, and suggest how to capture this
as an actionable item.

Return JSON:
{
    "extracted_text": "any text visible in image",
    "context": "description of what's in the image",
    "suggested_type": "task|note|event|reminder",
    "suggested_title": "short title for the item",
    "suggested_pillar": "hp|le|ge|null",
    "confidence": 0.0-1.0
}"""
                },
                {
                    "role": "user",
                    "content": [
                        {
                            "type": "text",
                            "text": f"Analyze this image. User context: {user_context}"
                        },
                        {
                            "type": "image_url",
                            "image_url": {
                                "url": f"data:image/jpeg;base64,{base64_image}"
                            }
                        }
                    ]
                }
            ],
            max_tokens=500,
        )

        # Parse JSON response
        import json
        try:
            return json.loads(response.choices[0].message.content)
        except:
            return {
                "extracted_text": "",
                "context": response.choices[0].message.content,
                "suggested_type": "note",
                "suggested_title": "Captured image",
                "confidence": 0.5
            }
```

### API Endpoints

**Create:** `backend/consciousness/context/api.py`

```python
from fastapi import APIRouter, Depends
from database import get_db, Database
from .engine import ContextEngine
from .models import UnifiedContext, MinimalContext

router = APIRouter(prefix="/api/v1/context", tags=["context"])

@router.get("/", response_model=UnifiedContext)
async def get_full_context(user_id: str = "default", db: Database = Depends(get_db)):
    """Get full unified context for AI operations."""
    engine = ContextEngine(db)
    return await engine.build(user_id)

@router.get("/minimal", response_model=MinimalContext)
async def get_minimal_context(user_id: str = "default", db: Database = Depends(get_db)):
    """Get minimal context for quick operations."""
    engine = ContextEngine(db)
    return await engine.build_minimal(user_id)

@router.get("/prompt")
async def get_context_as_prompt(user_id: str = "default", db: Database = Depends(get_db)):
    """Get context formatted as prompt text for AI injection."""
    engine = ContextEngine(db)
    context = await engine.build(user_id)
    return {"prompt": engine.to_prompt(context)}

@router.post("/invalidate")
async def invalidate_context(user_id: str = "default", db: Database = Depends(get_db)):
    """Invalidate cached context."""
    engine = ContextEngine(db)
    await engine.invalidate_cache(user_id)
    return {"status": "invalidated"}
```

---

## Success Criteria

- [ ] `ContextEngine.build()` returns complete context in <500ms
- [ ] All 7 aggregators implemented and working
- [ ] Oura biometrics flowing (when configured)
- [ ] Photo processing with GPT-4 Vision working
- [ ] Caching reduces repeated calls
- [ ] Context converts to prompt format correctly
- [ ] API endpoints functional

---

## When Done

### 1. Update Linear Issues to "In Review"

Update all 9 issues: OMOS-416, 417, 418, 419, 420, 421, 451, 456, 479

### 2. Post Handoff

```markdown
## Handoff: Agent 2 → Agents 3, 4

### Completed
- ✅ Context Engine in `backend/consciousness/context/engine.py`
- ✅ All 7 aggregators in `backend/consciousness/context/aggregators/`
- ✅ Photo processor in `backend/agno/capture/photo_processor.py`
- ✅ API endpoints in `backend/consciousness/context/api.py`

### For Dependent Agents
- Import: `from backend.consciousness.context.engine import ContextEngine`
- Usage: `context = await engine.build(user_id)`
- Prompt: `prompt_text = engine.to_prompt(context)`

### API Endpoints
- `GET /api/v1/context` - Full context
- `GET /api/v1/context/minimal` - Quick context
- `GET /api/v1/context/prompt` - As prompt text

### Performance
- Full build: ~XXms (target <500ms)
- Cached: ~Xms
```

---

**You are building the "consciousness" of OneMind. Every AI interaction will use your context. Make it comprehensive but fast.**
