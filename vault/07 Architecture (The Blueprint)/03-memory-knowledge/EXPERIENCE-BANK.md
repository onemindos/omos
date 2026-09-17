# Experience Bank - Behavioral Pattern Learning System

> **Created:** 2026-01-25
> **Status:** PLANNING
> **Linear Issue:** OMOS-XXX (to be created)

---

## Overview

The **Experience Bank** is a system that enables Legacy to learn your habits, detect behavioral patterns, and make predictions and recommendations based on your history. Unlike simple memory, the Experience Bank understands *when*, *where*, and *how* you do things.

**Key Insight:** The LLM (Claude) does ALL the reasoning - we just need to store patterns for it to analyze.

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         EXPERIENCE BANK ARCHITECTURE                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         DATA SOURCES                                   │  │
│  │                                                                        │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │  Activity   │  │  Session    │  │  Entity     │  │  External   │  │  │
│  │  │  Events     │  │  Summaries  │  │  Facts      │  │  Data       │  │  │
│  │  │             │  │             │  │             │  │             │  │  │
│  │  │ • Chat msgs │  │ • Topics    │  │ • People    │  │ • Oura      │  │  │
│  │  │ • Tool calls│  │ • Decisions │  │ • Projects  │  │ • Calendar  │  │  │
│  │  │ • Workflows │  │ • Outcomes  │  │ • Locations │  │ • Todoist   │  │  │
│  │  │ • Commands  │  │ • Duration  │  │ • Events    │  │ • GitHub    │  │  │
│  │  │             │  │             │  │             │  │             │  │  │
│  │  │ TimescaleDB │  │ PostgreSQL  │  │ PostgreSQL  │  │ Nango OAuth │  │  │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  │  │
│  │         │                │                │                │         │  │
│  └─────────┼────────────────┼────────────────┼────────────────┼─────────┘  │
│            │                │                │                │            │
│            └────────────────┴────────────────┴────────────────┘            │
│                                      │                                      │
│                                      ▼                                      │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                       PATTERN DETECTOR                                 │  │
│  │                       (Cron Job - Every Hour)                          │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    PATTERN TYPES                                 │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐   │  │  │
│  │  │  │   TIME    │  │ SEQUENCE  │  │CORRELATION│  │ LOCATION  │   │  │  │
│  │  │  │           │  │           │  │           │  │           │   │  │  │
│  │  │  │ "Codes at │  │ "Coffee → │  │ "7h sleep │  │ "Deep work│   │  │  │
│  │  │  │  6am on   │  │  Standup  │  │  = +30%   │  │  at home  │   │  │  │
│  │  │  │  weekdays"│  │  → Code"  │  │  output"  │  │  office"  │   │  │  │
│  │  │  │           │  │           │  │           │  │           │   │  │  │
│  │  │  └───────────┘  └───────────┘  └───────────┘  └───────────┘   │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Detection Methods:                                                    │  │
│  │  • Time patterns: GROUP BY hour, day_of_week, count occurrences       │  │
│  │  • Sequences: Track event order within 2-hour windows                 │  │
│  │  • Correlations: Compare metrics across days (sleep vs productivity)  │  │
│  │  • Location: Join activity with presence data                         │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                      │                                      │
│                                      ▼                                      │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        PATTERN STORE                                   │  │
│  │                        (TimescaleDB)                                   │  │
│  │                                                                        │  │
│  │  experience_patterns (hypertable)                                      │  │
│  │  ├── time TIMESTAMPTZ                                                 │  │
│  │  ├── user_id VARCHAR(255)                                             │  │
│  │  ├── pattern_type VARCHAR(50)  -- time, sequence, correlation, loc   │  │
│  │  ├── pattern_name VARCHAR(255) -- human readable name                 │  │
│  │  ├── pattern_data JSONB        -- full pattern details                │  │
│  │  ├── confidence FLOAT          -- 0.0 to 1.0                          │  │
│  │  └── occurrences INT           -- how many times seen                 │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                      │                                      │
│                                      ▼                                      │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        LLM REASONING                                   │  │
│  │                        (Claude - No Custom Algorithm)                  │  │
│  │                                                                        │  │
│  │  Input: Detected patterns from last 30 days                           │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Prompt:                                                         │  │  │
│  │  │  "Based on Zeus's patterns:                                      │  │  │
│  │  │   - Codes 6-9am weekdays (85% confidence)                        │  │  │
│  │  │   - Coffee → standup → code review sequence                      │  │  │
│  │  │   - 7+ hours sleep correlates with 30% more commits              │  │  │
│  │  │                                                                  │  │  │
│  │  │  1. What predictions for tomorrow (Monday)?                      │  │  │
│  │  │  2. Any recommendations right now (Sunday 9pm)?                  │  │  │
│  │  │  3. What correlations are most actionable?"                      │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Output:                                                               │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────────┐   │  │
│  │  │ PREDICTIONS │  │RECOMMENDATIONS│  │ CORRELATIONS              │   │  │
│  │  │             │  │             │  │                             │   │  │
│  │  │ "Tomorrow   │  │ "Get to bed │  │ "Your sleep directly        │   │  │
│  │  │  you'll     │  │  by 10pm to │  │  impacts coding output.     │   │  │
│  │  │  likely     │  │  maximize   │  │  Prioritize 7+ hours."      │   │  │
│  │  │  code 6-9am"│  │  tomorrow"  │  │                             │   │  │
│  │  └─────────────┘  └─────────────┘  └─────────────────────────────┘   │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Components Needed

| Component | What It Does | Already Have? | Build Effort |
|-----------|--------------|---------------|--------------|
| **TimescaleDB** | Time-series storage for patterns | ✅ YES (port 5433) | - |
| **Pattern Detector** | Cron job to analyze data | ❌ NO | 2-3 days |
| **Pattern Store** | Tables for detected patterns | ❌ NO | 1 day |
| **LLM Reasoning** | Claude analyzes patterns | ✅ YES | - |
| **Activity Events** | Track what you do when | ⚠️ PARTIAL | 1-2 days |
| **External Data** | Oura, calendar, etc. | ⚠️ PARTIAL (Nango ready) | 2-3 days |

**Total Build Effort:** ~8-10 days

---

## Database Schema

### TimescaleDB Tables (port 5433)

```sql
-- =============================================================================
-- EXPERIENCE BANK: Activity Events (raw data)
-- =============================================================================
CREATE TABLE IF NOT EXISTS activity_events (
    time TIMESTAMPTZ NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(50) NOT NULL,     -- chat, tool_call, workflow, command
    event_name VARCHAR(255) NOT NULL,    -- specific action
    event_data JSONB DEFAULT '{}',       -- full event details
    location_id VARCHAR(100),            -- where it happened
    device_id VARCHAR(100),              -- which device
    session_id VARCHAR(100),             -- which session
    pillar VARCHAR(10)                   -- hp, le, ge, it
);

SELECT create_hypertable('activity_events', 'time');

CREATE INDEX idx_activity_events_user ON activity_events(user_id, time DESC);
CREATE INDEX idx_activity_events_type ON activity_events(event_type);

-- =============================================================================
-- EXPERIENCE BANK: Detected Patterns
-- =============================================================================
CREATE TABLE IF NOT EXISTS experience_patterns (
    time TIMESTAMPTZ NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    pattern_type VARCHAR(50) NOT NULL,   -- time, sequence, correlation, location
    pattern_name VARCHAR(255) NOT NULL,  -- human readable name
    pattern_data JSONB NOT NULL,         -- full pattern details
    confidence FLOAT DEFAULT 0.5,        -- 0.0 to 1.0
    occurrences INT DEFAULT 1,           -- how many times seen
    last_seen TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE
);

SELECT create_hypertable('experience_patterns', 'time');

CREATE INDEX idx_patterns_user ON experience_patterns(user_id, time DESC);
CREATE INDEX idx_patterns_type ON experience_patterns(pattern_type);
CREATE INDEX idx_patterns_active ON experience_patterns(is_active) WHERE is_active = TRUE;

-- =============================================================================
-- EXPERIENCE BANK: External Metrics (from Oura, etc.)
-- =============================================================================
CREATE TABLE IF NOT EXISTS external_metrics (
    time TIMESTAMPTZ NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    source VARCHAR(50) NOT NULL,         -- oura, garmin, calendar, todoist
    metric_name VARCHAR(100) NOT NULL,   -- sleep_hours, hrv, tasks_completed
    metric_value FLOAT NOT NULL,
    metric_data JSONB DEFAULT '{}'       -- additional details
);

SELECT create_hypertable('external_metrics', 'time');

CREATE INDEX idx_metrics_user ON external_metrics(user_id, time DESC);
CREATE INDEX idx_metrics_source ON external_metrics(source);
```

---

## Pattern Types

### 1. Time Patterns

Detects when you typically do things.

```json
{
  "pattern_type": "time",
  "pattern_name": "morning_coding",
  "pattern_data": {
    "activity": "code_review",
    "hour_start": 6,
    "hour_end": 9,
    "days": ["monday", "tuesday", "wednesday", "thursday", "friday"],
    "frequency": 0.85
  },
  "confidence": 0.85,
  "occurrences": 42
}
```

**Detection Query:**
```sql
SELECT
    EXTRACT(HOUR FROM time) as hour,
    EXTRACT(DOW FROM time) as day_of_week,
    event_name,
    COUNT(*) as occurrences
FROM activity_events
WHERE user_id = 'zeus' AND time > NOW() - INTERVAL '30 days'
GROUP BY hour, day_of_week, event_name
HAVING COUNT(*) > 5
ORDER BY occurrences DESC;
```

### 2. Sequence Patterns

Detects what typically follows what.

```json
{
  "pattern_type": "sequence",
  "pattern_name": "morning_routine",
  "pattern_data": {
    "sequence": ["coffee", "standup", "code_review"],
    "avg_gap_minutes": [15, 30],
    "total_duration_minutes": 45
  },
  "confidence": 0.72,
  "occurrences": 28
}
```

**Detection Query:**
```sql
WITH event_pairs AS (
    SELECT
        e1.event_name as first_event,
        e2.event_name as second_event,
        EXTRACT(EPOCH FROM (e2.time - e1.time))/60 as gap_minutes
    FROM activity_events e1
    JOIN activity_events e2
        ON e1.user_id = e2.user_id
        AND e2.time > e1.time
        AND e2.time < e1.time + INTERVAL '2 hours'
    WHERE e1.user_id = 'zeus'
)
SELECT first_event, second_event, AVG(gap_minutes), COUNT(*)
FROM event_pairs
GROUP BY first_event, second_event
HAVING COUNT(*) > 5;
```

### 3. Correlation Patterns

Detects what affects what.

```json
{
  "pattern_type": "correlation",
  "pattern_name": "sleep_productivity",
  "pattern_data": {
    "when": "sleep_hours >= 7",
    "then": "commits_count +30%",
    "baseline_metric": 5.2,
    "correlated_metric": 6.8,
    "correlation_coefficient": 0.68
  },
  "confidence": 0.68,
  "occurrences": 21
}
```

**Detection Query:**
```sql
WITH daily_metrics AS (
    SELECT
        DATE(time) as day,
        MAX(CASE WHEN metric_name = 'sleep_hours' THEN metric_value END) as sleep,
        MAX(CASE WHEN metric_name = 'commits' THEN metric_value END) as commits
    FROM external_metrics
    WHERE user_id = 'zeus' AND time > NOW() - INTERVAL '30 days'
    GROUP BY DATE(time)
)
SELECT
    CORR(sleep, commits) as correlation,
    AVG(CASE WHEN sleep >= 7 THEN commits END) as high_sleep_commits,
    AVG(CASE WHEN sleep < 7 THEN commits END) as low_sleep_commits
FROM daily_metrics;
```

### 4. Location Patterns

Detects where you do what.

```json
{
  "pattern_type": "location",
  "pattern_name": "deep_work_location",
  "pattern_data": {
    "activity": "deep_coding",
    "location": "home_office",
    "frequency": 0.92,
    "avg_duration_minutes": 120
  },
  "confidence": 0.92,
  "occurrences": 35
}
```

---

## Pattern Detector Service

```python
# core/experience/detector.py
"""
Pattern Detector - Runs hourly to find patterns in activity data.
"""

import asyncio
from datetime import datetime, timedelta
from typing import List, Dict
import asyncpg

class PatternDetector:
    def __init__(self, timescale_url: str):
        self.timescale_url = timescale_url

    async def detect_all_patterns(self, user_id: str, days: int = 30) -> List[Dict]:
        """Run all pattern detection algorithms."""
        patterns = []

        # Detect each pattern type
        patterns.extend(await self.detect_time_patterns(user_id, days))
        patterns.extend(await self.detect_sequence_patterns(user_id, days))
        patterns.extend(await self.detect_correlation_patterns(user_id, days))
        patterns.extend(await self.detect_location_patterns(user_id, days))

        # Store patterns
        await self.store_patterns(patterns)

        return patterns

    async def detect_time_patterns(self, user_id: str, days: int) -> List[Dict]:
        """Detect when user typically does things."""
        # SQL query to find time-based patterns
        # Return list of pattern dicts
        pass

    async def detect_sequence_patterns(self, user_id: str, days: int) -> List[Dict]:
        """Detect what typically follows what."""
        pass

    async def detect_correlation_patterns(self, user_id: str, days: int) -> List[Dict]:
        """Detect what affects what."""
        pass

    async def detect_location_patterns(self, user_id: str, days: int) -> List[Dict]:
        """Detect where user does what."""
        pass

    async def store_patterns(self, patterns: List[Dict]):
        """Store detected patterns in TimescaleDB."""
        pass

# Cron job entry point
async def run_pattern_detection():
    detector = PatternDetector(TIMESCALE_URL)
    await detector.detect_all_patterns(user_id="zeus", days=30)

# Schedule: Run every hour
# 0 * * * * python -m core.experience.detector
```

---

## LLM Integration

The LLM does all reasoning - no custom algorithm needed:

```python
# core/experience/reasoning.py

async def get_experience_insights(user_id: str) -> str:
    """Use LLM to reason about detected patterns."""

    # Get recent patterns
    patterns = await get_recent_patterns(user_id, days=30)

    # Format for LLM
    patterns_text = format_patterns_for_llm(patterns)

    # Ask Claude to reason
    prompt = f"""
    Based on {user_id}'s behavioral patterns from the last 30 days:

    {patterns_text}

    Please provide:
    1. **Predictions** - What will {user_id} likely do tomorrow?
    2. **Recommendations** - What should {user_id} do right now to optimize?
    3. **Insights** - What interesting correlations or patterns stand out?

    Be specific and actionable. Reference the actual patterns.
    """

    response = await legacy.run(prompt, user_id=user_id)
    return response
```

---

## API Endpoints

```python
# core/experience/api.py

@router.get("/experience/patterns")
async def get_patterns(
    user_id: str = "zeus",
    pattern_type: Optional[str] = None,
    days: int = 30
):
    """Get detected patterns for a user."""
    pass

@router.get("/experience/insights")
async def get_insights(user_id: str = "zeus"):
    """Get LLM-generated insights from patterns."""
    pass

@router.post("/experience/detect")
async def trigger_detection(user_id: str = "zeus"):
    """Manually trigger pattern detection."""
    pass

@router.get("/experience/predictions")
async def get_predictions(user_id: str = "zeus"):
    """Get predictions for upcoming activities."""
    pass
```

---

## Implementation Phases

### Phase 1: Foundation (3 days)
- [ ] Create TimescaleDB schema (activity_events, experience_patterns, external_metrics)
- [ ] Add activity event logging to existing endpoints
- [ ] Basic pattern store API

### Phase 2: Pattern Detection (4 days)
- [ ] Time pattern detector
- [ ] Sequence pattern detector
- [ ] Location pattern detector
- [ ] Correlation pattern detector
- [ ] Cron job setup

### Phase 3: Integration (3 days)
- [ ] LLM reasoning integration
- [ ] API endpoints
- [ ] External data ingestion (Oura via Nango)
- [ ] UI display (1INTEL integration)

---

## Success Criteria

| Metric | Target |
|--------|--------|
| Time patterns detected | 5+ per user |
| Sequence patterns detected | 3+ per user |
| Correlation patterns detected | 2+ per user |
| Prediction accuracy | 70%+ |
| User finds insights valuable | Yes |

---

## Related Documentation

- `docs/architecture/LEARNING-MACHINE.md` - Memory system
- `docs/ONEINTEL.md` - Analytics UI
- `docs/architecture/FAMILY-PROFILES.md` - User profiles
- `infra/compose/docker-compose.core.yml` - TimescaleDB config

---

*Experience Bank - Teaching Legacy to understand your habits*
