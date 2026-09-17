# Agent 3: Intelligence & Scheduling (Opus 4.5)

> **Model:** Opus 4.5
> **Focus:** Intelligence Scores (PI/LI/WI) + DO DATE Scheduling
> **Issues:** 9 | **Points:** 24
> **Dependencies:** Agent 2 (Context Engine)

---

## Your Mission

You are Agent 3. Your job is building the **Intelligence Scores** (PI, LI, WI, Unified) and **Intelligent Scheduling** (DO DATE algorithm, energy matching, auto-reschedule).

These scores are what make OneMind "intelligent" - they track performance across all three pillars and help schedule tasks at optimal times.

---

## Your Assigned Issues

### Intelligence Scores
| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-422** | HPI Calculator (Holistic Performance Index) | 3 | P0-Critical |
| **OMOS-423** | LEI Calculator (Legacy Evolution Index) | 2 | - |
| **OMOS-424** | GEI Calculator (Generational Endeavor Index) | 2 | - |
| **OMOS-425** | UI Calculator (Unified Index) | 2 | - |
| **OMOS-426** | Score History & Trending | 2 | - |

### Intelligent Scheduling
| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-431** | DO DATE Algorithm | 5 | P1-High |
| **OMOS-432** | Auto-Reschedule | 3 | - |
| **OMOS-433** | Energy Matching | 3 | - |
| **OMOS-434** | Path Balance Scheduling | 2 | - |

---

## Dependencies

**Wait for Agent 2** to complete:
- Context Engine with BiometricsAggregator
- Items Aggregator for task/habit data

Once Agent 2 posts handoff, you can begin.

---

## Implementation

### Files to Create

```
backend/lifeos/scores/
├── __init__.py
├── models.py        # Score models
├── hpi.py          # Holistic Performance Index
├── lei.py          # Legacy Evolution Index
├── gei.py          # Generational Endeavor Index
├── unified.py      # Unified Index
├── history.py      # Score persistence & trending
├── service.py      # ScoreService
└── api.py          # Score endpoints

backend/lifeos/scheduling/
├── __init__.py
├── models.py        # Scheduling models
├── do_date.py      # DO DATE algorithm
├── auto_reschedule.py
├── energy_matching.py
├── path_balance.py
├── service.py       # SchedulingService
└── api.py          # Scheduling endpoints
```

---

## Intelligence Scores

### Score Definitions

| Score | Name | Path | Components |
|-------|------|------|------------|
| **PI** | Performance Index | HP (Health) | Sleep, Recovery, Habits, Tasks, Energy |
| **LI** | Legacy Index | LE (Family) | Family Time, Legacy Tasks, Relationships |
| **WI** | Wealth Index | GE (Business) | Revenue, Business Tasks, Investments, Network |
| **UI** | Unified Index | All | Average of PI + LI + WI |

### OMOS-422: HPI (Performance Index)

**Create:** `backend/lifeos/scores/hpi.py`

```python
from datetime import datetime, timedelta
from typing import Optional
from .models import ScoreResult, ScoreComponent

class HPICalculator:
    """
    Holistic Performance Index - measures personal health and performance.

    Components & Weights:
    - Sleep Quality: 20% (from Oura)
    - Recovery Score: 20% (from Oura)
    - Habit Completion: 25% (from habits)
    - Task Completion: 20% (from tasks)
    - Energy Levels: 15% (from Oura/state)
    """

    WEIGHTS = {
        "sleep": 0.20,
        "recovery": 0.20,
        "habits": 0.25,
        "tasks": 0.20,
        "energy": 0.15,
    }

    def __init__(self, db, context_engine=None):
        self.db = db
        self.context = context_engine

    async def calculate(self, user_id: str) -> ScoreResult:
        """Calculate current HPI score."""
        components = []

        # Sleep Quality (from Oura via context)
        sleep_score = await self._get_sleep_score(user_id)
        components.append(ScoreComponent(
            name="sleep",
            value=sleep_score,
            weight=self.WEIGHTS["sleep"],
            source="oura"
        ))

        # Recovery Score
        recovery_score = await self._get_recovery_score(user_id)
        components.append(ScoreComponent(
            name="recovery",
            value=recovery_score,
            weight=self.WEIGHTS["recovery"],
            source="oura"
        ))

        # Habit Completion (7-day rolling average)
        habit_score = await self._get_habit_score(user_id)
        components.append(ScoreComponent(
            name="habits",
            value=habit_score,
            weight=self.WEIGHTS["habits"],
            source="lifeos"
        ))

        # Task Completion (7-day rolling average)
        task_score = await self._get_task_score(user_id)
        components.append(ScoreComponent(
            name="tasks",
            value=task_score,
            weight=self.WEIGHTS["tasks"],
            source="lifeos"
        ))

        # Energy Level
        energy_score = await self._get_energy_score(user_id)
        components.append(ScoreComponent(
            name="energy",
            value=energy_score,
            weight=self.WEIGHTS["energy"],
            source="state"
        ))

        # Calculate weighted score
        total = sum(c.value * c.weight for c in components)

        return ScoreResult(
            score_type="hpi",
            value=round(total, 1),
            components=components,
            calculated_at=datetime.utcnow(),
        )

    async def _get_sleep_score(self, user_id: str) -> float:
        """Get sleep score from Oura or default."""
        if self.context:
            ctx = await self.context.build_minimal(user_id)
            if ctx.biometrics and ctx.biometrics.sleep_score:
                return ctx.biometrics.sleep_score
        return 70  # Default if no data

    async def _get_recovery_score(self, user_id: str) -> float:
        """Get recovery/readiness score from Oura."""
        if self.context:
            ctx = await self.context.build_minimal(user_id)
            if ctx.biometrics and ctx.biometrics.readiness_score:
                return ctx.biometrics.readiness_score
        return 70

    async def _get_habit_score(self, user_id: str) -> float:
        """Calculate habit completion rate (7-day rolling)."""
        query = """
            SELECT
                COUNT(CASE WHEN status = 'completed' THEN 1 END)::float /
                NULLIF(COUNT(*), 0) * 100 as completion_rate
            FROM assets
            WHERE asset_type = 'habit'
            AND created_at >= NOW() - INTERVAL '7 days'
        """
        row = await self.db.fetchrow(query)
        return row["completion_rate"] or 70

    async def _get_task_score(self, user_id: str) -> float:
        """Calculate task completion rate (7-day rolling)."""
        query = """
            SELECT
                COUNT(CASE WHEN status = 'completed' THEN 1 END)::float /
                NULLIF(COUNT(*), 0) * 100 as completion_rate
            FROM assets
            WHERE asset_type = 'task'
            AND pillar = 'hp'
            AND created_at >= NOW() - INTERVAL '7 days'
        """
        row = await self.db.fetchrow(query)
        return row["completion_rate"] or 70

    async def _get_energy_score(self, user_id: str) -> float:
        """Get current energy level as score."""
        if self.context:
            ctx = await self.context.build_minimal(user_id)
            if ctx.state:
                energy_map = {"high": 90, "medium": 70, "low": 50}
                return energy_map.get(ctx.state.energy_level, 70)
        return 70
```

### OMOS-423-424: LEI & GEI

Similar pattern, different components:

**LEI (Legacy Evolution Index):**
- Family Time: 25%
- Legacy Tasks: 25%
- Relationship Health: 20%
- Teaching/Mentoring: 15%
- Estate Planning: 15%

**GEI (Generational Endeavor Index):**
- Revenue Progress: 25%
- Business Tasks: 25%
- Investment Growth: 20%
- Network Building: 15%
- Skill Development: 15%

### OMOS-425: Unified Index

```python
class UnifiedIndexCalculator:
    """Calculate unified score from HPI, LEI, GEI."""

    def __init__(self, hpi: HPICalculator, lei: LEICalculator, gei: GEICalculator):
        self.hpi = hpi
        self.lei = lei
        self.gei = gei

    async def calculate(self, user_id: str) -> ScoreResult:
        # Calculate all three in parallel
        hpi_result, lei_result, gei_result = await asyncio.gather(
            self.hpi.calculate(user_id),
            self.lei.calculate(user_id),
            self.gei.calculate(user_id),
        )

        unified = (hpi_result.value + lei_result.value + gei_result.value) / 3

        # Determine balance
        scores = {"hpi": hpi_result.value, "lei": lei_result.value, "gei": gei_result.value}
        lowest = min(scores, key=scores.get)
        balance_warning = f"{lowest.upper()} needs attention" if scores[lowest] < 60 else None

        return ScoreResult(
            score_type="unified",
            value=round(unified, 1),
            components=[
                ScoreComponent(name="hpi", value=hpi_result.value, weight=0.33),
                ScoreComponent(name="lei", value=lei_result.value, weight=0.33),
                ScoreComponent(name="gei", value=gei_result.value, weight=0.34),
            ],
            balance_warning=balance_warning,
            calculated_at=datetime.utcnow(),
        )
```

### OMOS-426: Score History

```python
class ScoreHistory:
    """Persist and trend score data."""

    def __init__(self, db):
        self.db = db

    async def save_snapshot(self, user_id: str, scores: dict):
        """Save daily score snapshot."""
        query = """
            INSERT INTO score_history (user_id, date, hpi, lei, gei, unified)
            VALUES ($1, CURRENT_DATE, $2, $3, $4, $5)
            ON CONFLICT (user_id, date) DO UPDATE
            SET hpi = $2, lei = $3, gei = $4, unified = $5
        """
        await self.db.execute(query, user_id, scores["hpi"], scores["lei"], scores["gei"], scores["unified"])

    async def get_trend(self, user_id: str, score_type: str, days: int = 7) -> str:
        """Calculate trend direction: up, down, stable."""
        query = """
            SELECT {score_type} as score, date
            FROM score_history
            WHERE user_id = $1
            ORDER BY date DESC
            LIMIT $2
        """.format(score_type=score_type)

        rows = await self.db.fetch(query, user_id, days)
        if len(rows) < 2:
            return "stable"

        recent = rows[0]["score"]
        older = rows[-1]["score"]
        diff = recent - older

        if diff > 5:
            return "up"
        elif diff < -5:
            return "down"
        return "stable"
```

---

## Intelligent Scheduling

### OMOS-431: DO DATE Algorithm

**The Key Insight:** Tasks have a **DUE DATE** (deadline) but should be scheduled for a **DO DATE** (optimal time to work on it).

```python
class DODateAlgorithm:
    """
    Calculate optimal DO DATE for tasks (Motion-style scheduling).

    Factors:
    1. Task priority and deadline
    2. Estimated duration
    3. Available time slots
    4. Energy prediction for slot
    5. Path balance (HP/LE/GE)
    6. User preferences
    """

    def __init__(self, db, context_engine):
        self.db = db
        self.context = context_engine

    async def calculate(self, task_id: str) -> list[TimeSlotScore]:
        """Calculate top 3 recommended time slots for a task."""
        task = await self._get_task(task_id)
        if not task:
            return []

        # Get available slots for next 7 days
        available_slots = await self._get_available_slots(task.user_id, days=7)

        # Score each slot
        scored_slots = []
        for slot in available_slots:
            score = await self._score_slot(task, slot)
            scored_slots.append(TimeSlotScore(
                start=slot.start,
                end=slot.end,
                score=score,
                reasons=self._get_score_reasons(task, slot, score)
            ))

        # Return top 3
        scored_slots.sort(key=lambda x: x.score, reverse=True)
        return scored_slots[:3]

    async def _score_slot(self, task, slot) -> float:
        """Score a time slot for a task (0-100)."""
        score = 50  # Base score

        # Factor 1: Deadline proximity (max +20)
        if task.due_date:
            days_until_due = (task.due_date - slot.start).days
            if days_until_due <= 1:
                score += 20  # Urgent
            elif days_until_due <= 3:
                score += 10
            elif days_until_due > 7:
                score -= 5  # Can wait

        # Factor 2: Energy matching (max +20)
        predicted_energy = await self._predict_energy(slot.start)
        task_energy = task.energy_level or "medium"

        if task_energy == "high" and predicted_energy == "high":
            score += 20  # Perfect match
        elif task_energy == "low" and predicted_energy == "low":
            score += 15  # Good for low-energy tasks
        elif task_energy == "high" and predicted_energy == "low":
            score -= 15  # Mismatch

        # Factor 3: Path balance (max +15)
        path_balance = await self._get_path_balance(slot.start.date())
        if task.pillar and path_balance.get(task.pillar, 0) < 0.3:
            score += 15  # This path is underrepresented

        # Factor 4: Duration fit (max +15)
        slot_duration = (slot.end - slot.start).total_seconds() / 60
        if task.duration_minutes:
            if slot_duration >= task.duration_minutes:
                score += 15  # Fits perfectly
            elif slot_duration >= task.duration_minutes * 0.8:
                score += 10  # Close enough
            else:
                score -= 10  # Too short

        # Factor 5: User preference (max +10)
        hour = slot.start.hour
        if 9 <= hour <= 11:  # Morning focus
            score += 10
        elif 14 <= hour <= 16:  # Afternoon focus
            score += 5

        return min(max(score, 0), 100)

    async def _predict_energy(self, time: datetime) -> str:
        """Predict energy level for a time based on patterns."""
        hour = time.hour

        # Simple heuristic (can be enhanced with ML)
        if 6 <= hour <= 10:
            return "high"  # Morning peak
        elif 10 <= hour <= 14:
            return "medium"
        elif 14 <= hour <= 16:
            return "low"  # Post-lunch dip
        elif 16 <= hour <= 19:
            return "medium"
        else:
            return "low"

    async def _get_path_balance(self, date) -> dict:
        """Get path distribution for a date."""
        query = """
            SELECT pillar, COUNT(*) as count
            FROM assets
            WHERE do_date::date = $1 AND status != 'completed'
            GROUP BY pillar
        """
        rows = await self.db.fetch(query, date)
        total = sum(r["count"] for r in rows) or 1
        return {r["pillar"]: r["count"] / total for r in rows}
```

### OMOS-432-434: Supporting Features

```python
class AutoReschedule:
    """Automatically reschedule when conflicts arise."""

    async def check_and_reschedule(self, user_id: str):
        """Check for conflicts and auto-reschedule."""
        conflicts = await self._find_conflicts(user_id)
        for conflict in conflicts:
            new_slot = await self.do_date.calculate(conflict.task_id)
            if new_slot:
                await self._reschedule(conflict.task_id, new_slot[0])
                await self._notify_user(conflict, new_slot[0])


class EnergyMatching:
    """Match task difficulty to energy levels."""

    TASK_ENERGY_MAP = {
        "high": ["deep_work", "creative", "planning"],
        "medium": ["meetings", "communication", "review"],
        "low": ["admin", "routine", "reading"],
    }

    async def match(self, task, available_slots) -> list:
        """Return slots that match task energy requirements."""
        matched = []
        for slot in available_slots:
            energy = await self._predict_energy(slot.start)
            if self._is_match(task.energy_level, energy):
                matched.append(slot)
        return matched


class PathBalance:
    """Ensure HP/LE/GE tasks are balanced throughout the week."""

    IDEAL_DISTRIBUTION = {"hp": 0.40, "le": 0.30, "ge": 0.30}

    async def get_balance(self, user_id: str, period: str = "week") -> dict:
        """Get current path balance."""
        query = """
            SELECT pillar, COUNT(*) as count
            FROM assets
            WHERE do_date >= NOW() - INTERVAL '7 days'
            GROUP BY pillar
        """
        rows = await self.db.fetch(query)
        total = sum(r["count"] for r in rows) or 1

        balance = {p: 0.0 for p in ["hp", "le", "ge"]}
        for row in rows:
            if row["pillar"]:
                balance[row["pillar"]] = row["count"] / total

        return balance

    async def get_recommendations(self, user_id: str) -> list[str]:
        """Suggest which path needs attention."""
        balance = await self.get_balance(user_id)
        recs = []
        for pillar, ideal in self.IDEAL_DISTRIBUTION.items():
            if balance.get(pillar, 0) < ideal * 0.5:
                recs.append(f"{pillar.upper()} is neglected - schedule some {pillar} tasks")
        return recs
```

---

## API Endpoints

**Create:** `backend/lifeos/scores/api.py`

```python
from fastapi import APIRouter, Depends

router = APIRouter(prefix="/api/v1/scores", tags=["scores"])

@router.get("/hpi")
async def get_hpi(db = Depends(get_db)):
    """Get current HPI score."""
    return await HPICalculator(db).calculate("default")

@router.get("/lei")
async def get_lei(db = Depends(get_db)):
    """Get current LEI score."""
    return await LEICalculator(db).calculate("default")

@router.get("/gei")
async def get_gei(db = Depends(get_db)):
    """Get current GEI score."""
    return await GEICalculator(db).calculate("default")

@router.get("/unified")
async def get_unified(db = Depends(get_db)):
    """Get unified score with all components."""
    return await UnifiedIndexCalculator(...).calculate("default")

@router.get("/all")
async def get_all_scores(db = Depends(get_db)):
    """Get all scores at once."""
    ...

@router.get("/history")
async def get_score_history(days: int = 30, db = Depends(get_db)):
    """Get score history for trending."""
    return await ScoreHistory(db).get_history("default", days)
```

**Create:** `backend/lifeos/scheduling/api.py`

```python
router = APIRouter(prefix="/api/v1/scheduling", tags=["scheduling"])

@router.post("/tasks/{task_id}/schedule")
async def schedule_task(task_id: str, db = Depends(get_db)):
    """Get optimal DO DATE for a task."""
    return await DODateAlgorithm(db).calculate(task_id)

@router.get("/balance")
async def get_path_balance(db = Depends(get_db)):
    """Get current path balance."""
    return await PathBalance(db).get_balance("default")

@router.post("/auto-reschedule")
async def trigger_auto_reschedule(db = Depends(get_db)):
    """Manually trigger auto-reschedule check."""
    return await AutoReschedule(db).check_and_reschedule("default")
```

---

## Success Criteria

- [ ] All 4 scores (HPI, LEI, GEI, Unified) calculating correctly
- [ ] Score history storing daily snapshots
- [ ] Trends (up/down/stable) calculated correctly
- [ ] DO DATE algorithm returning scored slots
- [ ] Energy matching working
- [ ] Path balance tracking and recommendations
- [ ] API endpoints functional

---

## When Done

Update issues OMOS-422, 423, 424, 425, 426, 431, 432, 433, 434 to "In Review"

Post handoff to Agents 4, 5, 7 (who need scores).
