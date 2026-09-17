# OneMind Insights: Personal Analytics Engine

> **Status:** Planning Phase  
> **Created:** 2026-01-21  
> **Goal:** Replace Exist.io with a self-hosted, AI-powered correlation engine  
> **Philosophy:** Build insights INTO LEGACY, not as a separate service

---

## Executive Summary

**OneMind Insights** is a custom-built personal analytics and correlation engine that will replace Exist.io. Unlike Exist.io's static dashboards, OneMind Insights leverages Claude Opus 4.5 to provide conversational, actionable insights from your personal data.

### Why Build This?

| Factor | Exist.io | OneMind Insights |
|--------|----------|------------------|
| **Control** | Their servers | 100% yours |
| **AI Model** | Unknown/Basic | Claude Opus 4.5 |
| **Custom Metrics** | ~80 templates | Unlimited |
| **Conversation** | None | Ask LEGACY anything |
| **Actions** | View only | LEGACY can act on insights |
| **Cost** | $6/month forever | Free after build |
| **Privacy** | 3rd party | Fully owned |

---

## Comparison: Open Source Alternatives

### Tier 1: Best Open Source Options

| Project | Stars | Description | Correlation? | Self-Hosted |
|---------|-------|-------------|--------------|-------------|
| **[HPI](https://github.com/karlicoss/HPI)** | 1.6k ⭐ | Human Programming Interface - Python library to unify ALL your data | DIY | ✅ Yes |
| **[Heedy](https://github.com/heedy/heedy)** | 423 ⭐ | Personal metrics aggregator with analysis engine | Built-in | ✅ Yes |
| **[QS Ledger](https://github.com/markwk/qs_ledger)** | 1.1k ⭐ | Jupyter notebooks for data aggregation | DIY | ✅ Yes |
| **[Open Wearables](https://github.com/the-momentum/open-wearables)** | 400 ⭐ | Unify wearable health data (NEW - Oct 2025) | AI-ready | ✅ Yes |
| **[ActivityWatch](https://github.com/ActivityWatch/activitywatch)** | 16k ⭐ | Time tracking (focused on computer usage) | Limited | ✅ Yes |

### Deep Comparison

#### HPI (Human Programming Interface)
- **Best for:** Developers who want FULL control
- **Pros:** 100% open source, 50+ integrations, Python
- **Cons:** NO UI, no automatic correlations, requires coding
- **Verdict:** Great for aggregation, but YOU build the intelligence

#### Heedy
- **Best for:** Self-hosted with some analysis
- **Pros:** Web UI, plugin system, Docker support
- **Cons:** Last updated 4 years ago, small community
- **Verdict:** Good foundation but abandoned

#### QS Ledger
- **Best for:** Data scientists using Jupyter
- **Pros:** Correlation explorer, Todoist/Habitica integrations
- **Cons:** Not maintained (5 years), no automation
- **Verdict:** Good for manual analysis only

### The Truth

**There is NO open-source Exist.io equivalent with automatic AI correlations.**

The closest is building it yourself - which is exactly what OneMind Insights does.

---

## Why NOT Habitica + HPI?

| Factor | Habitica + HPI | OneMind Insights |
|--------|----------------|------------------|
| **Systems** | 3 (Todoist + Habitica + HPI) | 1 (LEGACY) |
| **Data Location** | Habitica = their servers | All local |
| **Gamification** | Limited to their mechanics | Custom ZeusRPG |
| **AI** | None | Claude Opus 4.5 |
| **Maintenance** | 3 projects to track | You own it |

**Key Insight:** HPI is just aggregation. It doesn't DO anything intelligent with data. You'd still need to build the correlation engine yourself.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         ONEMIND INSIGHTS                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌──────────────────┐     ┌──────────────────┐     ┌────────────────┐  │
│  │   DATA SOURCES   │     │   PROCESSING     │     │   STORAGE      │  │
│  ├──────────────────┤     ├──────────────────┤     ├────────────────┤  │
│  │ • Todoist        │     │                  │     │                │  │
│  │ • Apple Health   │────►│  LEGACY Agent    │────►│  PostgreSQL    │  │
│  │ • ZeusRPG        │     │  (Normalizer)    │     │  + pgvector    │  │
│  │ • Calendar       │     │                  │     │                │  │
│  │ • Location       │     └──────────────────┘     └───────┬────────┘  │
│  │ • Sleep (Oura?)  │                                      │           │
│  │ • Exercise       │                                      │           │
│  └──────────────────┘                                      │           │
│                                                            ▼           │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │                     CORRELATION ENGINE                           │  │
│  ├──────────────────────────────────────────────────────────────────┤  │
│  │  1. Statistical Analysis (Pearson, Spearman, Rolling Windows)    │  │
│  │  2. Pattern Detection (Weekly, Monthly, Seasonal)                │  │
│  │  3. Anomaly Detection (Z-score, IQR)                             │  │
│  │  4. Causal Inference (Granger causality tests)                   │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                      │                                 │
│                                      ▼                                 │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │                      AI INSIGHTS ENGINE                          │  │
│  ├──────────────────────────────────────────────────────────────────┤  │
│  │  Claude Opus 4.5 via AWS Bedrock                                 │  │
│  │  • Natural language insight generation                           │  │
│  │  • Multi-metric reasoning                                        │  │
│  │  • Actionable recommendations                                    │  │
│  │  • Semantic search via pgvector embeddings                       │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                      │                                 │
│                                      ▼                                 │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │                         OUTPUTS                                  │  │
│  ├──────────────────────────────────────────────────────────────────┤  │
│  │  • Conversational queries via LEGACY                             │  │
│  │  • Daily/Weekly insight reports                                  │  │
│  │  • Proactive notifications                                       │  │
│  │  • Dashboard widgets (OneMind UI)                                │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Components Required

### 1. Data Layer (PostgreSQL)

**Purpose:** Store all personal metrics in a unified schema

```sql
-- Core metrics table
CREATE TABLE daily_metrics (
    id SERIAL PRIMARY KEY,
    user_id UUID DEFAULT 'zeus',
    date DATE NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    value FLOAT NOT NULL,
    value_type VARCHAR(20) DEFAULT 'number',  -- number, duration, boolean, scale
    source VARCHAR(50) NOT NULL,              -- todoist, apple_health, zeusrpg
    metadata JSONB,                           -- extra context
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, date, metric_name)
);

-- Index for fast queries
CREATE INDEX idx_metrics_date ON daily_metrics(date);
CREATE INDEX idx_metrics_name ON daily_metrics(metric_name);
CREATE INDEX idx_metrics_source ON daily_metrics(source);

-- Correlations cache
CREATE TABLE correlations (
    id SERIAL PRIMARY KEY,
    metric_a VARCHAR(100) NOT NULL,
    metric_b VARCHAR(100) NOT NULL,
    correlation_type VARCHAR(20) DEFAULT 'pearson',
    coefficient FLOAT NOT NULL,
    p_value FLOAT,
    sample_size INT,
    confidence_interval JSONB,
    insight TEXT,                             -- AI-generated
    is_significant BOOLEAN,
    last_calculated TIMESTAMP DEFAULT NOW(),
    UNIQUE(metric_a, metric_b, correlation_type)
);

-- Daily context embeddings for semantic search
CREATE TABLE daily_embeddings (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    embedding VECTOR(1536),                   -- Claude/OpenAI embedding
    summary TEXT,                             -- "Productive day, 8hr sleep, exercised"
    tags TEXT[],                              -- ['high_productivity', 'good_sleep']
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insight history
CREATE TABLE insights (
    id SERIAL PRIMARY KEY,
    insight_type VARCHAR(50),                 -- correlation, anomaly, trend, recommendation
    title VARCHAR(200),
    body TEXT,
    metrics_involved TEXT[],
    confidence FLOAT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Files to Create:**
- `agents/legacy/db/models/metrics.py`
- `agents/legacy/db/models/correlations.py`
- `agents/legacy/db/models/insights.py`
- `infra/migrations/001_insights_schema.sql`

---

### 2. Data Ingestion Layer

**Purpose:** Pull data from various sources into unified format

```python
# agents/legacy/insights/ingestion/base.py

from abc import ABC, abstractmethod
from datetime import date
from typing import List
from pydantic import BaseModel

class MetricValue(BaseModel):
    date: date
    metric_name: str
    value: float
    source: str
    metadata: dict = {}

class DataSource(ABC):
    """Base class for all data sources."""
    
    @property
    @abstractmethod
    def source_name(self) -> str:
        pass
    
    @abstractmethod
    async def fetch_metrics(self, start_date: date, end_date: date) -> List[MetricValue]:
        pass
    
    @abstractmethod
    async def get_available_metrics(self) -> List[str]:
        pass
```

**Data Sources to Implement:**

| Source | Metrics | Priority |
|--------|---------|----------|
| **Todoist** | tasks_completed, p1_completed, p2_completed, overdue_tasks | P0 |
| **ZeusRPG** | xp_earned, gold_earned, streak_days, level | P0 |
| **Apple Health** | steps, sleep_hours, active_calories, heart_rate_avg | P1 |
| **Calendar** | meetings_count, meeting_hours, focus_time | P1 |
| **Location** | places_visited, home_time, commute_time | P2 |
| **Weather** | temperature, conditions (for correlation) | P2 |

**Files to Create:**
- `agents/legacy/insights/ingestion/todoist.py`
- `agents/legacy/insights/ingestion/zeusrpg.py`
- `agents/legacy/insights/ingestion/apple_health.py`
- `agents/legacy/insights/ingestion/calendar.py`
- `agents/legacy/insights/ingestion/weather.py`

---

### 3. Correlation Engine

**Purpose:** Calculate statistical relationships between metrics

```python
# agents/legacy/insights/correlation/engine.py

import numpy as np
from scipy import stats
from typing import List, Tuple, Optional
from dataclasses import dataclass

@dataclass
class CorrelationResult:
    metric_a: str
    metric_b: str
    coefficient: float
    p_value: float
    sample_size: int
    is_significant: bool
    correlation_type: str

class CorrelationEngine:
    """Statistical correlation analysis engine."""
    
    def __init__(self, db, significance_threshold: float = 0.05):
        self.db = db
        self.significance_threshold = significance_threshold
    
    async def calculate_all_correlations(self) -> List[CorrelationResult]:
        """Calculate pairwise correlations for all metrics."""
        metrics = await self.db.get_distinct_metrics()
        results = []
        
        for i, metric_a in enumerate(metrics):
            for metric_b in metrics[i+1:]:
                result = await self.calculate_correlation(metric_a, metric_b)
                if result:
                    results.append(result)
        
        return results
    
    async def calculate_correlation(
        self, 
        metric_a: str, 
        metric_b: str,
        method: str = 'pearson'
    ) -> Optional[CorrelationResult]:
        """Calculate correlation between two metrics."""
        
        # Get aligned data
        data = await self.db.get_aligned_metrics(metric_a, metric_b)
        
        if len(data) < 14:  # Minimum 2 weeks
            return None
        
        values_a = [d['value_a'] for d in data]
        values_b = [d['value_b'] for d in data]
        
        if method == 'pearson':
            coef, p_value = stats.pearsonr(values_a, values_b)
        elif method == 'spearman':
            coef, p_value = stats.spearmanr(values_a, values_b)
        else:
            raise ValueError(f"Unknown method: {method}")
        
        return CorrelationResult(
            metric_a=metric_a,
            metric_b=metric_b,
            coefficient=coef,
            p_value=p_value,
            sample_size=len(data),
            is_significant=p_value < self.significance_threshold and abs(coef) > 0.3,
            correlation_type=method
        )
    
    async def calculate_lagged_correlation(
        self,
        metric_a: str,
        metric_b: str,
        lag_days: int = 1
    ) -> Optional[CorrelationResult]:
        """Calculate correlation with time lag (e.g., sleep today → productivity tomorrow)."""
        # Implementation for lagged analysis
        pass
    
    async def detect_trends(
        self,
        metric_name: str,
        window_days: int = 30
    ) -> dict:
        """Detect trends in a metric over time."""
        data = await self.db.get_metric_history(metric_name, window_days)
        
        if len(data) < 7:
            return {"trend": "insufficient_data"}
        
        values = [d['value'] for d in data]
        x = np.arange(len(values))
        
        slope, intercept, r_value, p_value, std_err = stats.linregress(x, values)
        
        if p_value > 0.05:
            trend = "stable"
        elif slope > 0:
            trend = "increasing"
        else:
            trend = "decreasing"
        
        return {
            "trend": trend,
            "slope": slope,
            "r_squared": r_value ** 2,
            "p_value": p_value,
            "change_per_day": slope
        }
```

**Additional Analysis Methods:**
- Rolling window correlations
- Seasonal decomposition
- Anomaly detection (Z-score, IQR)
- Granger causality tests

**Files to Create:**
- `agents/legacy/insights/correlation/engine.py`
- `agents/legacy/insights/correlation/trends.py`
- `agents/legacy/insights/correlation/anomaly.py`
- `agents/legacy/insights/correlation/causality.py`

---

### 4. AI Insights Generator

**Purpose:** Turn statistics into human-readable, actionable insights

```python
# agents/legacy/insights/ai/generator.py

from typing import List
from agents.legacy.llm import get_claude_client

class InsightGenerator:
    """Generate natural language insights using Claude."""
    
    def __init__(self):
        self.claude = get_claude_client()
    
    async def generate_correlation_insight(
        self,
        metric_a: str,
        metric_b: str,
        coefficient: float,
        sample_size: int
    ) -> str:
        """Generate insight from correlation."""
        
        direction = "positively" if coefficient > 0 else "negatively"
        strength = self._get_strength_label(abs(coefficient))
        
        prompt = f"""
        Based on {sample_size} days of personal data analysis:
        
        Finding: {metric_a} is {strength} {direction} correlated with {metric_b}
        Correlation coefficient: {coefficient:.2f}
        
        Generate a brief, actionable insight for the user. Requirements:
        1. One sentence maximum
        2. Use plain language (no statistical jargon)
        3. Include a specific, actionable recommendation if possible
        4. Be encouraging, not preachy
        
        Example format: "You complete 40% more tasks on days after sleeping 7+ hours—prioritizing sleep could boost your productivity."
        """
        
        response = await self.claude.generate(prompt)
        return response.text.strip()
    
    async def generate_daily_summary(
        self,
        date: str,
        metrics: dict
    ) -> str:
        """Generate a summary of the day's metrics."""
        
        prompt = f"""
        Summarize this person's day based on their tracked metrics:
        Date: {date}
        Metrics: {metrics}
        
        Create a 2-3 sentence summary that:
        1. Highlights what went well
        2. Notes any areas for improvement
        3. Is encouraging and supportive
        """
        
        response = await self.claude.generate(prompt)
        return response.text.strip()
    
    async def generate_weekly_report(
        self,
        week_data: dict,
        correlations: List[dict],
        trends: List[dict]
    ) -> str:
        """Generate comprehensive weekly insights report."""
        
        prompt = f"""
        Generate a weekly insights report for Zeus based on this data:
        
        ## Weekly Metrics Summary
        {week_data}
        
        ## Significant Correlations Found
        {correlations}
        
        ## Trends Detected
        {trends}
        
        Create a report with:
        1. Key wins from the week
        2. Most important insight discovered
        3. One specific recommendation for next week
        4. Encouragement and motivation
        
        Keep it concise but meaningful. Use bullet points.
        """
        
        response = await self.claude.generate(prompt)
        return response.text
    
    def _get_strength_label(self, abs_coefficient: float) -> str:
        if abs_coefficient >= 0.7:
            return "strongly"
        elif abs_coefficient >= 0.5:
            return "moderately"
        else:
            return "weakly"
```

**Files to Create:**
- `agents/legacy/insights/ai/generator.py`
- `agents/legacy/insights/ai/embeddings.py`
- `agents/legacy/insights/ai/prompts.py`

---

### 5. Semantic Search (pgvector)

**Purpose:** Find similar days, patterns, and context

```python
# agents/legacy/insights/search/semantic.py

from typing import List
import numpy as np

class SemanticSearch:
    """Semantic search over daily embeddings using pgvector."""
    
    def __init__(self, db, embedding_client):
        self.db = db
        self.embeddings = embedding_client
    
    async def embed_day(self, date: str, summary: str) -> None:
        """Create and store embedding for a day."""
        
        embedding = await self.embeddings.create(summary)
        
        await self.db.execute("""
            INSERT INTO daily_embeddings (date, embedding, summary)
            VALUES ($1, $2, $3)
            ON CONFLICT (date) DO UPDATE SET embedding = $2, summary = $3
        """, date, embedding, summary)
    
    async def find_similar_days(
        self,
        query: str,
        limit: int = 5
    ) -> List[dict]:
        """Find days similar to a query."""
        
        query_embedding = await self.embeddings.create(query)
        
        results = await self.db.fetch("""
            SELECT date, summary, 1 - (embedding <=> $1) as similarity
            FROM daily_embeddings
            ORDER BY embedding <=> $1
            LIMIT $2
        """, query_embedding, limit)
        
        return results
    
    async def find_days_like(
        self,
        reference_date: str,
        limit: int = 5
    ) -> List[dict]:
        """Find days similar to a reference day."""
        
        ref = await self.db.fetchrow(
            "SELECT embedding FROM daily_embeddings WHERE date = $1",
            reference_date
        )
        
        if not ref:
            return []
        
        results = await self.db.fetch("""
            SELECT date, summary, 1 - (embedding <=> $1) as similarity
            FROM daily_embeddings
            WHERE date != $2
            ORDER BY embedding <=> $1
            LIMIT $3
        """, ref['embedding'], reference_date, limit)
        
        return results
```

**Files to Create:**
- `agents/legacy/insights/search/semantic.py`
- `agents/legacy/insights/search/patterns.py`

---

### 6. LEGACY Agent Tools

**Purpose:** Expose insights to LEGACY for conversational queries

```python
# agents/legacy/tools/insights.py

from agno import Tool

class InsightsTools:
    """LEGACY tools for OneMind Insights."""
    
    @Tool(name="insights_get_correlations")
    async def get_correlations(
        self,
        min_strength: float = 0.3,
        metric_filter: str = None
    ) -> list:
        """
        Get significant correlations from your personal data.
        
        Args:
            min_strength: Minimum correlation coefficient (0-1)
            metric_filter: Only show correlations involving this metric
        
        Returns:
            List of correlations with AI-generated insights
        """
        pass
    
    @Tool(name="insights_query")
    async def query_insights(
        self,
        question: str
    ) -> str:
        """
        Ask a natural language question about your personal data.
        
        Args:
            question: E.g., "What affects my productivity the most?"
        
        Returns:
            AI-generated answer based on your data
        """
        pass
    
    @Tool(name="insights_get_trends")
    async def get_trends(
        self,
        metric_name: str = None,
        days: int = 30
    ) -> dict:
        """
        Get trends for a specific metric or all metrics.
        
        Returns:
            Trend direction, change rate, and insight
        """
        pass
    
    @Tool(name="insights_find_similar_days")
    async def find_similar_days(
        self,
        description: str
    ) -> list:
        """
        Find days similar to a description.
        
        Args:
            description: E.g., "highly productive days with good sleep"
        
        Returns:
            List of matching dates with metrics
        """
        pass
    
    @Tool(name="insights_weekly_report")
    async def generate_weekly_report(self) -> str:
        """
        Generate a comprehensive weekly insights report.
        
        Returns:
            Formatted report with wins, insights, and recommendations
        """
        pass
    
    @Tool(name="insights_record_metric")
    async def record_metric(
        self,
        metric_name: str,
        value: float,
        date: str = None
    ) -> dict:
        """
        Manually record a metric value.
        
        Args:
            metric_name: Name of the metric
            value: Numeric value
            date: Date (defaults to today)
        """
        pass
```

**Files to Create:**
- `agents/legacy/tools/insights.py`

---

### 7. Scheduler / Cron Jobs

**Purpose:** Automate data ingestion and analysis

```python
# agents/legacy/insights/scheduler.py

from apscheduler.schedulers.asyncio import AsyncIOScheduler

class InsightsScheduler:
    """Schedule automated insights tasks."""
    
    def __init__(self, engine, generator, ingestion):
        self.engine = engine
        self.generator = generator
        self.ingestion = ingestion
        self.scheduler = AsyncIOScheduler()
    
    def start(self):
        # Sync data every hour
        self.scheduler.add_job(
            self.sync_all_sources,
            'interval',
            hours=1,
            id='sync_sources'
        )
        
        # Calculate correlations daily at 3 AM
        self.scheduler.add_job(
            self.calculate_correlations,
            'cron',
            hour=3,
            id='daily_correlations'
        )
        
        # Generate daily summary at 10 PM
        self.scheduler.add_job(
            self.generate_daily_summary,
            'cron',
            hour=22,
            id='daily_summary'
        )
        
        # Generate weekly report on Sundays at 9 AM
        self.scheduler.add_job(
            self.generate_weekly_report,
            'cron',
            day_of_week='sun',
            hour=9,
            id='weekly_report'
        )
        
        self.scheduler.start()
```

**Files to Create:**
- `agents/legacy/insights/scheduler.py`

---

### 8. UI Components (Optional)

**Purpose:** Visual dashboard in OneMind UI

```tsx
// ui/src/components/insights/CorrelationCard.tsx

interface CorrelationCardProps {
  metricA: string;
  metricB: string;
  coefficient: number;
  insight: string;
}

export function CorrelationCard({ metricA, metricB, coefficient, insight }: CorrelationCardProps) {
  const isPositive = coefficient > 0;
  const strength = Math.abs(coefficient);
  
  return (
    <div className="correlation-card">
      <div className="metrics">
        <span>{metricA}</span>
        <span className={isPositive ? "positive" : "negative"}>
          {isPositive ? "↑↑" : "↓↑"}
        </span>
        <span>{metricB}</span>
      </div>
      <div className="strength-bar">
        <div style={{ width: `${strength * 100}%` }} />
      </div>
      <p className="insight">{insight}</p>
    </div>
  );
}
```

**Files to Create:**
- `ui/src/components/insights/CorrelationCard.tsx`
- `ui/src/components/insights/TrendChart.tsx`
- `ui/src/components/insights/WeeklyReport.tsx`
- `ui/src/components/insights/InsightsDashboard.tsx`

---

## Component Summary

| Component | Purpose | Tech Stack | Priority |
|-----------|---------|------------|----------|
| **PostgreSQL Schema** | Store metrics & correlations | PostgreSQL + pgvector | P0 |
| **Data Ingestion** | Pull from sources | Python async | P0 |
| **Correlation Engine** | Statistical analysis | NumPy, SciPy | P0 |
| **AI Generator** | Natural language insights | Claude Opus 4.5 | P0 |
| **Semantic Search** | Find similar patterns | pgvector embeddings | P1 |
| **LEGACY Tools** | Conversational interface | Agno framework | P1 |
| **Scheduler** | Automate tasks | APScheduler | P1 |
| **UI Dashboard** | Visual insights | React, TailwindCSS | P2 |

---

## Dependencies

```toml
# pyproject.toml additions

[project.dependencies]
numpy = "^1.26"
scipy = "^1.12"
pandas = "^2.2"
asyncpg = "^0.29"
pgvector = "^0.2"
apscheduler = "^3.10"
httpx = "^0.27"
```

---

## Implementation Roadmap

### Phase 1: Foundation (Week 1-2)
- [ ] Create PostgreSQL schema
- [ ] Build base DataSource class
- [ ] Implement Todoist ingestion
- [ ] Implement ZeusRPG ingestion

### Phase 2: Correlation Engine (Week 3-4)
- [ ] Build correlation calculator
- [ ] Add trend detection
- [ ] Implement anomaly detection
- [ ] Store results in database

### Phase 3: AI Integration (Week 5-6)
- [ ] Create insight generator
- [ ] Build embedding pipeline
- [ ] Implement semantic search
- [ ] Create LEGACY tools

### Phase 4: Automation (Week 7)
- [ ] Set up scheduler
- [ ] Automate daily sync
- [ ] Automate weekly reports
- [ ] Add notification system

### Phase 5: UI (Week 8)
- [ ] Build correlation cards
- [ ] Create trend charts
- [ ] Build insights dashboard
- [ ] Mobile responsive design

---

## The Killer Feature: Conversational Insights

```
Zeus: "LEGACY, what affects my productivity the most?"

LEGACY: "Based on 90 days of data, I found 3 significant correlations:

1. **Sleep duration** (r=0.72): You complete 40% more tasks on days 
   after sleeping 7+ hours.

2. **Morning exercise** (r=0.58): Days you exercise before 10am, 
   you earn 35% more XP.

3. **Meeting count** (r=-0.45): More than 3 meetings correlates 
   with 25% fewer completed tasks.

**Recommendation:** Protect your mornings for exercise and deep work."
```

Exist.io shows correlations. OneMind Insights **reasons** about them.

---

## Related Documentation

- [ZEUSRPG-GAMIFICATION.md](../onemind-os/docs/ZEUSRPG-GAMIFICATION.md) - Gamification system
- [TODOIST-EXISTIO-INTEGRATION.md](../onemind-os/docs/integrations/TODOIST-EXISTIO-INTEGRATION.md) - Current integration plan
- [HABITICA-INTEGRATION.md](../onemind-os/docs/integrations/HABITICA-INTEGRATION.md) - Deferred for Q3 2026

---

## Philosophy

> "The goal isn't to collect data—it's to understand yourself."

OneMind Insights embodies the core equation: `HUMAN (Zeus) + AI (LEGACY) = ONE MIND`

Not: `Zeus + Habitica + HPI + Exist.io = Fragmented Mind`

---

*"Know thyself—with data."*
