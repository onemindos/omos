# OneMind Insights: Personal Analytics Engine

> **Status:** Planning Phase | Use Exist.io Now, Build This Q2-Q4 2026  
> **Created:** 2026-01-21  
> **Goal:** Replace Exist.io with a self-hosted, AI-powered correlation engine  
> **Philosophy:** Build insights INTO Legacy, not as a separate service

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Comparison: OneMind Insights vs Exist.io](#comparison-onemind-insights-vs-existio)
3. [Open Source Alternatives Analysis](#open-source-alternatives-analysis)
4. [Exist.io Feature Reference](#existio-feature-reference)
5. [Data Types & Platforms Reference](#data-types--platforms-reference)
6. [Architecture Overview](#architecture-overview)
7. [Components to Build](#components-to-build)
8. [Database Schema](#database-schema)
9. [Files to Create](#files-to-create)
10. [Implementation Roadmap](#implementation-roadmap)
11. [Dependencies](#dependencies)
12. [Health Scores & Behavioral Archetypes](#health-scores--behavioral-archetypes)
13. [Health Data Aggregator Evaluation](#health-data-aggregator-evaluation)

---

## Executive Summary

**OneMind Insights** is a custom-built personal analytics and correlation engine that will replace Exist.io. Unlike Exist.io's static dashboards, OneMind Insights leverages Claude Opus 4.5 to provide conversational, actionable insights from your personal data.

### The Core Equation

```
EXIST.IO:    Data Sources → Correlations → Static Dashboard
ONEMIND:     Data Sources → Correlations → AI Reasoning → Conversational Interface
```

### Why Build This?

| Factor | Exist.io | OneMind Insights |
| ------ | -------- | ---------------- |
| **Control** | Their servers | 100% yours |
| **AI Model** | Unknown/Basic | Claude Opus 4.5 |
| **Custom Metrics** | ~80 templates | Unlimited |
| **Conversation** | None | Ask Legacy anything |
| **Actions** | View only | Legacy can act |
| **Cost** | $6.99/month forever | Free after build |
| **Privacy** | 3rd party Australia | Fully owned |
| **Semantic Search** | No | pgvector embeddings |
| **Causal Reasoning** | No | AI-powered |

---

## Comparison: OneMind Insights vs Exist.io

### Feature Comparison

| Feature | Exist.io | OneMind Insights | Winner |
| ------- | -------- | ---------------- | ------ |
| **Automatic data sync** | ✅ 20+ integrations | 🔨 Build needed | Exist |
| **Correlation engine** | ✅ Weekly generation | ✅ Daily + on-demand | Tie |
| **AI insight generation** | ⚠️ Template sentences | ✅ Claude reasoning | OneMind |
| **Natural language queries** | ❌ No | ✅ "What affects my sleep?" | OneMind |
| **Take action on insights** | ❌ No | ✅ Legacy can act | OneMind |
| **Custom metrics** | ⚠️ ~80 templates | ✅ Unlimited | OneMind |
| **Semantic search** | ❌ No | ✅ "Find productive days" | OneMind |
| **Mobile app** | ✅ iOS/Android | ⚠️ PWA (future) | Exist |
| **Weekly email reports** | ✅ Yes | ✅ + Discord/Telegram | Tie |
| **Self-hosted** | ❌ No | ✅ Yes | OneMind |
| **API for writing data** | ✅ OAuth2 | ✅ Direct + API | Tie |
| **Historical backfill** | ✅ Yes | ✅ Yes | Tie |
| **Trend detection** | ✅ Rolling averages | ✅ + ML patterns | OneMind |
| **Anomaly detection** | ⚠️ Basic | ✅ Z-score, IQR | OneMind |
| **Cost** | $6.99/mo ($84/yr) | $0 after build | OneMind |

### Correlation Intelligence Comparison

**Exist.io generates template sentences:**
```
"You are more productive when you listen to music more."
"Your weight is higher when you wake more during sleep."
```

**OneMind Insights with Claude:**
```
"Based on 90 days of data, I found that you complete 40% more tasks 
on days after sleeping 7+ hours. Your most productive days also 
correlate with morning exercise before 10am. 

Recommendation: Protect your mornings—prioritize sleep and exercise 
before deep work blocks for maximum productivity."
```

---

## Open Source Alternatives Analysis

### Why NOT Use Existing Open Source?

| Project | Stars | Pros | Cons | Verdict |
| ------- | ----- | ---- | ---- | ------- |
| **[HPI](https://github.com/karlicoss/HPI)** | 1.6k ⭐ | 50+ integrations, Python, self-hosted | No UI, no correlations, requires coding | Aggregation only |
| **[Heedy](https://github.com/heedy/heedy)** | 423 ⭐ | Web UI, Docker, plugin system | Abandoned (4 years), small community | Dead project |
| **[QS Ledger](https://github.com/markwk/qs_ledger)** | 1.1k ⭐ | Jupyter notebooks, correlation explorer | Not maintained (5 years), manual only | Research tool |
| **[Open Wearables](https://github.com/the-momentum/open-wearables)** | 400 ⭐ | New (Oct 2025), AI-ready, wearables focus | Early stage, limited scope | Watch this |
| **[ActivityWatch](https://github.com/ActivityWatch/activitywatch)** | 16k ⭐ | Great time tracking, privacy-focused | Computer usage only, no correlations | Niche tool |

### The Truth

**There is NO open-source Exist.io equivalent with automatic AI correlations.**

All existing projects fall into two categories:
1. **Aggregators** (HPI, QS Ledger) - Collect data, but no intelligence
2. **Niche trackers** (ActivityWatch) - Track one thing well

**Conclusion:** Build your own using existing infrastructure.

---

## Exist.io Feature Reference

### Core Features to Replicate

#### 1. Automatic Correlations
- **What:** Find relationships between metrics (Pearson correlation)
- **When:** Generated weekly on Sundays
- **How:** Compare all pairs of attributes over 60+ days
- **Output:** Percentage strength, P-value, star rating (1-5)

#### 2. Insights
- **What:** Interesting events in user data
- **Types:**
  - Day-level: "Yesterday was your highest steps in 30 days"
  - Week-level: "You completed 45 tasks this week, 20% more than usual"
  - Month-level: "Total steps for May: 250,000"
- **Priority:** 1 (real-time), 2 (daily), 3 (weekly), 4 (monthly)

#### 3. Averages & Trends
- **Rolling averages:** Last 60 days
- **Day-of-week breakdowns:** See patterns by day
- **Progress bars:** Daily progress vs your average

#### 4. Mood Tracking
- **Scale:** 1-9 (bad to great)
- **Notes:** Free-text daily notes
- **Correlations:** What makes you happiest?

#### 5. Custom Tags
- **Boolean attributes:** "meditation", "headache", "worked from home"
- **Quick entry:** Tap to add, automatic correlation

### Correlation Star System

| Stars | P-Value | Confidence |
| ----- | ------- | ---------- |
| ⭐⭐⭐⭐⭐ | < 0.001 | Certain to be related |
| ⭐⭐⭐⭐ | < 0.01 | Very likely related |
| ⭐⭐⭐ | < 0.05 | Likely related |
| ⭐⭐ | < 0.10 | Possibly related |
| ⭐ | < 0.20 | Weak evidence |

### Insight Types

| Type | Example | Trigger |
| ---- | ------- | ------- |
| **High/Low** | "Highest steps in 30 days" | New max/min in window |
| **Streak** | "5-day streak of 10k steps" | Consecutive threshold met |
| **Comparison** | "20% more productive than last week" | Week-over-week change |
| **Pattern** | "You listen to more music on Fridays" | Day-of-week correlation |
| **Milestone** | "Total 1 million steps this year" | Round number reached |

---

## Data Types & Platforms Reference

### Exist.io Connected Platforms

| Platform | Category | Data Synced |
| -------- | -------- | ----------- |
| **Apple Health / Watch** | Health | Steps, sleep, workouts, heart rate, HRV |
| **Health Connect (Android)** | Health | Steps, sleep, workouts, heart rate |
| **Fitbit** | Wearable | Steps, sleep, heart rate, weight |
| **Garmin** | Wearable | Steps, sleep, workouts, heart rate |
| **Oura Ring** | Wearable | Sleep, readiness, activity, HRV |
| **Withings** | Wearable | Weight, body fat, blood pressure, sleep |
| **Strava** | Fitness | Workouts, distance, time |
| **RescueTime** | Productivity | Productive/distracting time, apps used |
| **Toggl** | Productivity | Time tracked per project |
| **Todoist** | Productivity | Tasks completed, by priority |
| **GitHub** | Productivity | Commits per day |
| **Foreground** | Productivity | App usage, screen time |
| **Google Calendar** | Events | Events count, time in events |
| **Apple Calendar** | Events | Events count, time in events |
| **Instapaper** | Media | Articles read, word count |
| **Last.fm** | Media | Tracks played |
| **Spotify** | Media | Tracks played (via Last.fm) |
| **Trakt** | Media | TV/movies watched |
| **Swarm/Foursquare** | Location | Check-ins |
| **Mastodon** | Social | Posts, replies, boosts, favorites |
| **Apple Weather** | Weather | Temp, precipitation, cloud cover, wind |

### All Attribute Types (80+ Templates)

#### Activity (8 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `steps` | Integer | count | Apple Health, Fitbit, Garmin |
| `steps_active_min` | Duration | minutes | Apple Health, Fitbit |
| `steps_distance` | Float | km | Apple Health, Fitbit |
| `steps_elevation` | Float | km | Fitbit |
| `floors` | Integer | count | Fitbit |
| `stand_hours` | Integer | hours | Apple Watch |
| `cycle_min` | Duration | minutes | Apple Health |
| `cycle_distance` | Float | km | Apple Health, Strava |
| `active_energy` | Float | kJ | Apple Health, Fitbit |

#### Workouts (3 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `workouts` | Integer | count | Apple Health, Strava |
| `workouts_min` | Duration | minutes | Apple Health, Strava |
| `workouts_distance` | Float | km | Strava |

#### Productivity (10 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `productive_min` | Duration | minutes | RescueTime |
| `neutral_min` | Duration | minutes | RescueTime |
| `distracting_min` | Duration | minutes | RescueTime |
| `commits` | Integer | count | GitHub |
| `tasks_completed` | Integer | count | Todoist |
| `words_written` | Integer | count | Manual |
| `emails_sent` | Integer | count | Manual |
| `emails_received` | Integer | count | Manual |
| `pomodoros_min` | Duration | minutes | Manual |
| `keystrokes` | Integer | count | RescueTime |

#### Mood (4 attributes)
| Attribute | Type | Scale | Source |
| --------- | ---- | ----- | ------ |
| `mood` | Scale | 1-9 | Manual |
| `mood_note` | String | 1000 chars | Manual |
| `energy_level` | Scale | 1-9 | Manual |
| `stress_level` | Scale | 1-9 | Manual |

#### Sleep (9 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `sleep` | Duration | minutes | Apple Health, Oura, Fitbit |
| `time_in_bed` | Duration | minutes | Apple Health, Oura |
| `sleep_light` | Duration | minutes | Oura, Fitbit |
| `sleep_deep` | Duration | minutes | Oura, Fitbit |
| `sleep_rem` | Duration | minutes | Oura, Fitbit |
| `sleep_start` | Time | minutes from midnight | Oura |
| `sleep_end` | Time | minutes from midnight | Oura |
| `sleep_awakenings` | Integer | count | Oura, Fitbit |

#### Events (2 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `events` | Integer | count | Google/Apple Calendar |
| `events_duration` | Duration | minutes | Google/Apple Calendar |

#### Food & Drink (12 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `coffees` | Integer | count | Manual |
| `alcoholic_drinks` | Integer | count | Manual |
| `water` | Integer | ml | Manual |
| `energy` | Float | kJ | Manual |
| `carbohydrates` | Float | g | Manual |
| `fat` | Float | g | Manual |
| `fibre` | Float | g | Manual |
| `protein` | Float | g | Manual |
| `sugar` | Float | g | Manual |
| `sodium` | Float | mg | Manual |
| `cholesterol` | Float | mg | Manual |
| `caffeine` | Float | mg | Manual |

#### Health & Body (11 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `weight` | Float | kg | Withings, Manual |
| `body_fat` | Percentage | 0.0-1.0 | Withings |
| `lean_mass` | Float | kg | Withings |
| `heartrate` | Integer | bpm | Apple Health, Fitbit |
| `heartrate_max` | Integer | bpm | Apple Health |
| `heartrate_resting` | Integer | bpm | Apple Health, Oura |
| `heartrate_variability` | Integer | ms | Oura, Apple Health |
| `meditation_min` | Duration | minutes | Manual |
| `menstrual_flow` | Integer | 0-4 | Manual |
| `sexual_activity` | Integer | count | Manual |

#### Location (2 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `checkins` | Integer | count | Swarm |
| `location` | String | lat,lng | GPS |

#### Media (7 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `tracks` | Integer | count | Last.fm, Spotify |
| `articles_read` | Integer | count | Instapaper |
| `pages_read` | Integer | count | Instapaper |
| `mobile_screen_min` | Duration | minutes | Foreground |
| `gaming_min` | Duration | minutes | Trakt, Manual |
| `tv_min` | Duration | minutes | Trakt |
| `movies_min` | Duration | minutes | Trakt |

#### Social (7 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `mastodon_posts` | Integer | count | Mastodon |
| `mastodon_replies` | Integer | count | Mastodon |
| `mastodon_boosts` | Integer | count | Mastodon |
| `mastodon_favourites` | Integer | count | Mastodon |
| `tweets` | Integer | count | Deprecated |
| `twitter_mentions` | Integer | count | Deprecated |

#### Weather (8 attributes)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `weather_temp_max` | Float | °C | Apple Weather |
| `weather_temp_min` | Float | °C | Apple Weather |
| `weather_precipitation` | Float | inches/hr | Apple Weather |
| `weather_cloud_cover` | Percentage | 0.0-1.0 | Apple Weather |
| `weather_wind_speed` | Float | km/hr | Apple Weather |
| `weather_summary` | String | text | Apple Weather |
| `weather_icon` | String | icon name | Apple Weather |
| `day_length` | Duration | minutes | Calculated |

#### Finance (1 attribute)
| Attribute | Type | Unit | Source |
| --------- | ---- | ---- | ------ |
| `money_spent` | Float | local currency | Manual |

#### Custom Tags (Boolean)
Examples: `meditation`, `headache`, `worked_from_home`, `alcohol`, `exercise`, `travel`, `sick`

### Value Types

| Type | Code | Description | Example |
| ---- | ---- | ----------- | ------- |
| Integer | 0 | Whole numbers | 10,542 steps |
| Float | 1 | Decimal numbers | 73.5 kg |
| String | 2 | Text (max 250 chars) | "Partly cloudy" |
| Duration | 3 | Minutes as integer | 480 (8 hours) |
| Time of Day | 4 | Minutes from midnight | 420 (7:00 AM) |
| Percentage | 5 | Float 0.0 to 1.0 | 0.75 (75%) |
| Boolean | 7 | True/False | 1 or 0 |
| Scale | 8 | Integer 1-9 | 7 (good mood) |

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            ONEMIND INSIGHTS                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        DATA INGESTION LAYER                           │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                       │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐         │  │
│  │  │ Todoist │ │ Apple   │ │ GitHub  │ │ Oura    │ │ Weather │         │  │
│  │  │         │ │ Health  │ │         │ │ Ring    │ │ API     │         │  │
│  │  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘         │  │
│  │       │           │           │           │           │               │  │
│  │       └───────────┴───────────┴───────────┴───────────┘               │  │
│  │                               │                                       │  │
│  │                               ▼                                       │  │
│  │                    ┌─────────────────────┐                            │  │
│  │                    │  Normalizer Agent   │                            │  │
│  │                    │  (Legacy)           │                            │  │
│  │                    └──────────┬──────────┘                            │  │
│  │                               │                                       │  │
│  └───────────────────────────────┼───────────────────────────────────────┘  │
│                                  │                                          │
│  ┌───────────────────────────────┼───────────────────────────────────────┐  │
│  │                        STORAGE LAYER                                  │  │
│  ├───────────────────────────────┼───────────────────────────────────────┤  │
│  │                               ▼                                       │  │
│  │            ┌─────────────────────────────────────┐                    │  │
│  │            │         PostgreSQL + pgvector       │                    │  │
│  │            ├─────────────────────────────────────┤                    │  │
│  │            │  • daily_metrics (all data points)  │                    │  │
│  │            │  • correlations (cached results)    │                    │  │
│  │            │  • insights (generated insights)    │                    │  │
│  │            │  • daily_embeddings (semantic)      │                    │  │
│  │            │  • averages (rolling stats)         │                    │  │
│  │            └─────────────────────────────────────┘                    │  │
│  │                                                                       │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                  │                                          │
│  ┌───────────────────────────────┼───────────────────────────────────────┐  │
│  │                     ANALYSIS LAYER                                    │  │
│  ├───────────────────────────────┼───────────────────────────────────────┤  │
│  │                               ▼                                       │  │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐       │  │
│  │  │  Correlation    │  │  Trend          │  │  Anomaly        │       │  │
│  │  │  Engine         │  │  Detection      │  │  Detection      │       │  │
│  │  │  (Pearson,      │  │  (Linear reg,   │  │  (Z-score,      │       │  │
│  │  │  Spearman)      │  │  Seasonal)      │  │  IQR)           │       │  │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────┘       │  │
│  │                               │                                       │  │
│  └───────────────────────────────┼───────────────────────────────────────┘  │
│                                  │                                          │
│  ┌───────────────────────────────┼───────────────────────────────────────┐  │
│  │                      AI LAYER                                         │  │
│  ├───────────────────────────────┼───────────────────────────────────────┤  │
│  │                               ▼                                       │  │
│  │            ┌─────────────────────────────────────┐                    │  │
│  │            │       Claude Opus 4.5               │                    │  │
│  │            │       (AWS Bedrock)                 │                    │  │
│  │            ├─────────────────────────────────────┤                    │  │
│  │            │  • Natural language insights        │                    │  │
│  │            │  • Multi-metric reasoning           │                    │  │
│  │            │  • Actionable recommendations       │                    │  │
│  │            │  • Question answering               │                    │  │
│  │            └─────────────────────────────────────┘                    │  │
│  │                               │                                       │  │
│  │  ┌─────────────────┐          │                                       │  │
│  │  │  Semantic       │◄─────────┘                                       │  │
│  │  │  Search         │  "Find days like today"                          │  │
│  │  │  (pgvector)     │                                                  │  │
│  │  └─────────────────┘                                                  │  │
│  │                                                                       │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                  │                                          │
│  ┌───────────────────────────────┼───────────────────────────────────────┐  │
│  │                      OUTPUT LAYER                                     │  │
│  ├───────────────────────────────┼───────────────────────────────────────┤  │
│  │                               ▼                                       │  │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐       │  │
│  │  │  Legacy Agent   │  │  Dashboard UI   │  │  Notifications  │       │  │
│  │  │  (Conversation) │  │  (React)        │  │  (Discord/      │       │  │
│  │  │                 │  │                 │  │   Telegram)     │       │  │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────┘       │  │
│  │                                                                       │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Components to Build

### 1. Data Layer (PostgreSQL + pgvector)

**Purpose:** Store all personal metrics in a unified schema

**Key Tables:**
- `daily_metrics` - All data points by date
- `correlations` - Cached correlation results
- `insights` - Generated insights
- `daily_embeddings` - Semantic vectors for days
- `averages` - Rolling averages by metric
- `metric_templates` - Attribute definitions

### 2. Data Ingestion Layer

**Purpose:** Pull data from various sources into unified format

**Priority Data Sources:**

| Priority | Source | Metrics | Complexity |
| -------- | ------ | ------- | ---------- |
| P0 | **Todoist** | tasks_completed, by priority | Low (API) |
| P0 | **ZeusRPG** | xp, gold, streak, level | Low (internal) |
| P1 | **Apple Health** | steps, sleep, workouts, HR | Medium (HealthKit) |
| P1 | **Calendar** | events, event_duration | Low (API) |
| P1 | **GitHub** | commits | Low (API) |
| P2 | **Weather** | temp, conditions | Low (API) |
| P2 | **Oura** | sleep, HRV, readiness | Medium (API) |
| P2 | **RescueTime** | productive_min, distracting_min | Low (API) |
| P3 | **Last.fm** | tracks | Low (API) |
| P3 | **Location** | checkins, home_time | High (GPS) |

### 3. Correlation Engine

**Purpose:** Calculate statistical relationships between metrics

**Methods to Implement:**
- Pearson correlation (linear relationships)
- Spearman correlation (ranked relationships)
- Rolling window correlations (30, 60, 90 days)
- Lagged correlations (yesterday's X → today's Y)
- Day-of-week patterns
- Seasonal decomposition

### 4. Trend Detection

**Purpose:** Identify patterns over time

**Methods:**
- Linear regression (increasing/decreasing/stable)
- Rolling averages (7, 30, 60 days)
- Day-of-week breakdown
- Month-over-month comparison
- Year-over-year comparison

### 5. Anomaly Detection

**Purpose:** Alert on unusual days

**Methods:**
- Z-score (how many standard deviations from mean)
- IQR (interquartile range outliers)
- Streak detection (consecutive days above/below threshold)
- New high/low detection

### 6. AI Insights Generator

**Purpose:** Turn statistics into human-readable, actionable insights

**Capabilities:**
- Generate natural language from correlations
- Combine multiple insights into recommendations
- Answer natural language questions
- Generate daily/weekly summaries
- Create embeddings for semantic search

### 7. Semantic Search (pgvector)

**Purpose:** Find similar days, patterns, and context

**Queries:**
- "Find days when I was highly productive"
- "What happened on days like today?"
- "When was my sleep best?"

### 8. Legacy Agent Tools

**Purpose:** Expose insights via conversational interface

**Tools:**
- `insights_get_correlations` - Get significant correlations
- `insights_query` - Ask questions about your data
- `insights_get_trends` - Get trend analysis
- `insights_find_similar_days` - Semantic search
- `insights_weekly_report` - Generate weekly summary
- `insights_record_metric` - Manual data entry

### 9. Scheduler

**Purpose:** Automate data ingestion and analysis

**Jobs:**
- Hourly: Sync data from sources
- Daily 3 AM: Calculate correlations
- Daily 10 PM: Generate daily summary
- Sunday 9 AM: Generate weekly report
- Monthly: Generate monthly insights

### 10. UI Components (Optional)

**Purpose:** Visual dashboard in OneMind UI

**Components:**
- Correlation cards
- Trend charts
- Insight feed
- Weekly report view
- Metric entry forms

---

## Database Schema

```sql
-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- ============================================
-- CORE METRICS TABLE
-- ============================================
CREATE TABLE daily_metrics (
    id SERIAL PRIMARY KEY,
    user_id UUID DEFAULT 'zeus',
    date DATE NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    value FLOAT,
    value_type INTEGER DEFAULT 0,  -- 0=int, 1=float, 3=duration, 5=pct, 7=bool, 8=scale
    source VARCHAR(50) NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, date, metric_name)
);

CREATE INDEX idx_metrics_date ON daily_metrics(date);
CREATE INDEX idx_metrics_name ON daily_metrics(metric_name);
CREATE INDEX idx_metrics_source ON daily_metrics(source);

-- ============================================
-- METRIC TEMPLATES (Like Exist's 80+ attributes)
-- ============================================
CREATE TABLE metric_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    label VARCHAR(200) NOT NULL,
    group_name VARCHAR(50) NOT NULL,
    value_type INTEGER NOT NULL,
    unit VARCHAR(50),
    description TEXT,
    is_manual BOOLEAN DEFAULT FALSE,
    priority INTEGER DEFAULT 100
);

-- ============================================
-- CORRELATIONS CACHE
-- ============================================
CREATE TABLE correlations (
    id SERIAL PRIMARY KEY,
    metric_a VARCHAR(100) NOT NULL,
    metric_b VARCHAR(100) NOT NULL,
    correlation_type VARCHAR(20) DEFAULT 'pearson',
    coefficient FLOAT NOT NULL,
    p_value FLOAT,
    sample_size INTEGER,
    stars INTEGER,  -- 1-5 confidence rating
    is_significant BOOLEAN,
    insight TEXT,  -- AI-generated
    second_person TEXT,  -- "You are more X when Y"
    period_days INTEGER DEFAULT 60,
    lag_days INTEGER DEFAULT 0,
    calculated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(metric_a, metric_b, correlation_type, period_days, lag_days)
);

CREATE INDEX idx_correlations_significant ON correlations(is_significant);
CREATE INDEX idx_correlations_strength ON correlations(coefficient);

-- ============================================
-- INSIGHTS
-- ============================================
CREATE TABLE insights (
    id SERIAL PRIMARY KEY,
    insight_type VARCHAR(50) NOT NULL,  -- correlation, high, low, streak, trend, summary
    priority INTEGER DEFAULT 2,  -- 1=realtime, 2=daily, 3=weekly, 4=monthly
    target_date DATE,
    title VARCHAR(200),
    body TEXT,
    html TEXT,
    metrics_involved TEXT[],
    confidence FLOAT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_insights_date ON insights(target_date);
CREATE INDEX idx_insights_type ON insights(insight_type);
CREATE INDEX idx_insights_unread ON insights(is_read) WHERE is_read = FALSE;

-- ============================================
-- DAILY EMBEDDINGS (Semantic Search)
-- ============================================
CREATE TABLE daily_embeddings (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    embedding VECTOR(1536),  -- Claude/OpenAI embedding dimension
    summary TEXT,  -- "Productive day, 8hr sleep, exercised"
    tags TEXT[],  -- ['high_productivity', 'good_sleep']
    metrics_snapshot JSONB,  -- All metrics for that day
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_embeddings_vector ON daily_embeddings 
    USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- ============================================
-- AVERAGES (Rolling Stats)
-- ============================================
CREATE TABLE averages (
    id SERIAL PRIMARY KEY,
    metric_name VARCHAR(100) NOT NULL,
    calculated_date DATE NOT NULL,
    period_days INTEGER DEFAULT 60,
    overall FLOAT,
    monday FLOAT,
    tuesday FLOAT,
    wednesday FLOAT,
    thursday FLOAT,
    friday FLOAT,
    saturday FLOAT,
    sunday FLOAT,
    UNIQUE(metric_name, calculated_date, period_days)
);

-- ============================================
-- CUSTOM TAGS
-- ============================================
CREATE TABLE custom_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    label VARCHAR(200) NOT NULL,
    color VARCHAR(7) DEFAULT '#6366f1',
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- DATA SOURCE SYNC STATUS
-- ============================================
CREATE TABLE sync_status (
    id SERIAL PRIMARY KEY,
    source VARCHAR(50) UNIQUE NOT NULL,
    last_sync TIMESTAMP,
    last_success TIMESTAMP,
    error_message TEXT,
    is_active BOOLEAN DEFAULT TRUE
);
```

---

## Files to Create

### Directory Structure

```
agents/legacy/
├── insights/
│   ├── __init__.py
│   ├── config.py                    # Configuration and constants
│   │
│   ├── ingestion/
│   │   ├── __init__.py
│   │   ├── base.py                  # Abstract DataSource class
│   │   ├── todoist.py               # Todoist data source
│   │   ├── zeusrpg.py               # ZeusRPG data source
│   │   ├── apple_health.py          # Apple Health data source
│   │   ├── calendar.py              # Calendar data source
│   │   ├── github.py                # GitHub commits source
│   │   ├── weather.py               # Weather API source
│   │   ├── oura.py                  # Oura Ring source
│   │   ├── rescuetime.py            # RescueTime source
│   │   └── manual.py                # Manual entry handler
│   │
│   ├── analysis/
│   │   ├── __init__.py
│   │   ├── correlation.py           # Correlation engine
│   │   ├── trends.py                # Trend detection
│   │   ├── anomaly.py               # Anomaly detection
│   │   ├── patterns.py              # Day-of-week patterns
│   │   └── averages.py              # Rolling averages
│   │
│   ├── ai/
│   │   ├── __init__.py
│   │   ├── generator.py             # Claude insight generator
│   │   ├── embeddings.py            # Embedding creation
│   │   ├── prompts.py               # Prompt templates
│   │   └── reasoning.py             # Multi-metric reasoning
│   │
│   ├── search/
│   │   ├── __init__.py
│   │   └── semantic.py              # pgvector semantic search
│   │
│   ├── scheduler.py                 # APScheduler jobs
│   └── reports.py                   # Weekly/monthly reports
│
├── tools/
│   └── insights.py                  # Legacy agent tools
│
└── db/
    └── models/
        ├── metrics.py               # Metric models
        ├── correlations.py          # Correlation models
        └── insights.py              # Insight models

infra/
└── migrations/
    ├── 001_insights_schema.sql      # Core tables
    ├── 002_metric_templates.sql     # Pre-populate templates
    └── 003_indexes.sql              # Performance indexes

ui/src/components/insights/
├── CorrelationCard.tsx              # Single correlation display
├── CorrelationList.tsx              # List of correlations
├── TrendChart.tsx                   # Trend visualization
├── InsightFeed.tsx                  # Insight stream
├── WeeklyReport.tsx                 # Weekly summary
├── MetricEntry.tsx                  # Manual entry form
├── DayView.tsx                      # Single day metrics
└── InsightsDashboard.tsx            # Main dashboard
```

### File Count Summary

| Category | Files | Priority |
| -------- | ----- | -------- |
| **Ingestion** | 10 | P0-P2 |
| **Analysis** | 5 | P0 |
| **AI** | 4 | P0 |
| **Search** | 1 | P1 |
| **Tools** | 1 | P1 |
| **Models** | 3 | P0 |
| **Migrations** | 3 | P0 |
| **UI Components** | 8 | P2 |
| **Total** | 35 | - |

---

## Implementation Roadmap

### Phase 1: Foundation (Week 1-2)
- [ ] Create PostgreSQL schema (migrations)
- [ ] Build base DataSource class
- [ ] Implement Todoist ingestion
- [ ] Implement ZeusRPG ingestion
- [ ] Create metric templates table
- [ ] Basic CRUD for daily_metrics

### Phase 2: Correlation Engine (Week 3-4)
- [ ] Implement Pearson correlation
- [ ] Implement Spearman correlation
- [ ] Add rolling window analysis
- [ ] Calculate P-values and star ratings
- [ ] Build correlation cache
- [ ] Generate "second_person" sentences

### Phase 3: Analysis Features (Week 5)
- [ ] Trend detection (linear regression)
- [ ] Rolling averages calculation
- [ ] Day-of-week pattern detection
- [ ] Anomaly detection (Z-score)
- [ ] Streak detection

### Phase 4: AI Integration (Week 6-7)
- [ ] Create insight generator with Claude
- [ ] Build embedding pipeline
- [ ] Implement semantic search
- [ ] Create question-answering system
- [ ] Build multi-metric reasoning

### Phase 5: Legacy Tools (Week 8)
- [ ] `insights_get_correlations` tool
- [ ] `insights_query` tool
- [ ] `insights_get_trends` tool
- [ ] `insights_find_similar_days` tool
- [ ] `insights_weekly_report` tool

### Phase 6: Automation (Week 9)
- [ ] Set up APScheduler
- [ ] Hourly data sync jobs
- [ ] Daily correlation calculation
- [ ] Weekly report generation
- [ ] Notification system (Discord)

### Phase 7: Additional Sources (Week 10-12)
- [ ] Apple Health integration
- [ ] Calendar integration
- [ ] GitHub integration
- [ ] Weather API integration
- [ ] Oura Ring integration (if applicable)

### Phase 8: UI Dashboard (Week 13-14)
- [ ] Correlation cards component
- [ ] Trend charts
- [ ] Insight feed
- [ ] Weekly report view
- [ ] Manual entry forms

---

## Dependencies

### Python Dependencies

```toml
# pyproject.toml additions for OneMind Insights

[project.dependencies]
# Data analysis
numpy = "^1.26"
scipy = "^1.12"
pandas = "^2.2"

# Database
asyncpg = "^0.29"
pgvector = "^0.2"
sqlalchemy = "^2.0"

# Scheduling
apscheduler = "^3.10"

# HTTP clients
httpx = "^0.27"
aiohttp = "^3.9"

# Date handling
python-dateutil = "^2.8"

# Statistics
statsmodels = "^0.14"  # For seasonal decomposition

# Embeddings (optional, if not using Claude directly)
openai = "^1.0"  # For embeddings API
```

### External Services Required

| Service | Purpose | Required? |
| ------- | ------- | --------- |
| **PostgreSQL 15+** | Database with pgvector | ✅ Yes |
| **Claude via Bedrock** | AI insights | ✅ Yes |
| **Todoist API** | Task data | ✅ Yes |
| **Apple HealthKit** | Health data | ⚠️ Optional |
| **Google Calendar API** | Events | ⚠️ Optional |
| **GitHub API** | Commits | ⚠️ Optional |
| **OpenWeather API** | Weather context | ⚠️ Optional |
| **Oura API** | Sleep/HRV | ⚠️ Optional |

---

## Timeline Summary

| Phase | Duration | Deliverable |
| ----- | -------- | ----------- |
| **Now** | 2026 Q1-Q2 | Use Exist.io, collect patterns |
| **Phase 1-2** | 4 weeks | Database + Correlations |
| **Phase 3-5** | 5 weeks | Analysis + AI + Tools |
| **Phase 6-7** | 4 weeks | Automation + Sources |
| **Phase 8** | 2 weeks | UI Dashboard |
| **Total** | ~15 weeks | Full replacement |
| **Drop Exist.io** | Q4 2026 | Self-sufficient |

---

## The Killer Feature: Conversational Insights

```
Zeus: "Legacy, what affects my productivity the most?"

Legacy: "Based on 90 days of data, I found 3 significant correlations:

1. **Sleep duration** (r=0.72, ⭐⭐⭐⭐⭐): You complete 40% more tasks 
   on days after sleeping 7+ hours.

2. **Morning exercise** (r=0.58, ⭐⭐⭐⭐): Days you exercise before 10am, 
   you earn 35% more XP.

3. **Meeting count** (r=-0.45, ⭐⭐⭐⭐): More than 3 meetings correlates 
   with 25% fewer completed tasks.

**Recommendation:** Protect your mornings for exercise and deep work. 
Consider batching meetings to specific days."
```

**Exist.io can show you correlations. OneMind Insights can reason about them.**

---

## Philosophy

> "The goal isn't to collect data—it's to understand yourself."

OneMind Insights embodies the core equation:

```
HUMAN (Zeus) + AI (Legacy) = ONE MIND
```

Not: `Zeus + Habitica + HPI + Exist.io = Fragmented Mind`

---

---

## Health Scores & Behavioral Archetypes

> **Inspiration:** Sahha.ai offers pre-built health scores and "archetypes" (behavioral profiles). We'll build our own with Claude for full control and explainability.

### Why Build Custom Health Scores?

| Factor | Sahha (Black Box) | OneMind Insights (Custom) |
| ------ | ----------------- | ------------------------- |
| **Explainability** | "Your score is 72" | "Your score is 72 because sleep was 6hrs (-15), but exercise (+10) and low stress (+7)" |
| **Customization** | Fixed algorithms | Tune weights per Zeus's goals |
| **AI Reasoning** | Pre-baked | Claude explains WHY |
| **Transparency** | Proprietary | Open formula, debuggable |
| **Action-taking** | View only | Legacy can act on scores |

### Health Scores to Build

| Score | Metrics Used | Range | Update |
| ----- | ------------ | ----- | ------ |
| **Readiness Score** | Sleep, HRV, resting HR, previous day strain | 0-100 | Daily AM |
| **Recovery Score** | Sleep quality, sleep duration, HRV trend | 0-100 | Daily AM |
| **Strain Score** | Active calories, workout intensity, steps | 0-100 | Daily PM |
| **Productivity Score** | Tasks completed, focus time, distraction ratio | 0-100 | Daily PM |
| **Balance Score** | Work hours, leisure time, social events | 0-100 | Weekly |
| **Consistency Score** | Streak adherence, routine stability | 0-100 | Weekly |
| **Wellness Score** | Composite of all above | 0-100 | Daily |

### Score Calculation Example: Readiness Score

```python
def calculate_readiness_score(metrics: dict) -> tuple[int, str]:
    """
    Calculate readiness score with explainable breakdown.
    Returns (score, explanation)
    """
    score = 50  # Base score
    factors = []
    
    # Sleep (40% weight)
    sleep_hours = metrics.get('sleep_duration', 0) / 60
    if sleep_hours >= 7.5:
        sleep_bonus = 20
        factors.append(f"Sleep {sleep_hours:.1f}hrs (+{sleep_bonus})")
    elif sleep_hours >= 6:
        sleep_bonus = 10
        factors.append(f"Sleep {sleep_hours:.1f}hrs (+{sleep_bonus})")
    else:
        sleep_bonus = -15
        factors.append(f"Sleep {sleep_hours:.1f}hrs ({sleep_bonus})")
    score += sleep_bonus
    
    # HRV (25% weight)
    hrv = metrics.get('hrv', 0)
    hrv_baseline = metrics.get('hrv_baseline', 50)
    hrv_diff = ((hrv - hrv_baseline) / hrv_baseline) * 100
    if hrv_diff > 10:
        hrv_bonus = 15
        factors.append(f"HRV {hrv}ms, {hrv_diff:.0f}% above baseline (+{hrv_bonus})")
    elif hrv_diff > -10:
        hrv_bonus = 5
        factors.append(f"HRV {hrv}ms, normal range (+{hrv_bonus})")
    else:
        hrv_bonus = -10
        factors.append(f"HRV {hrv}ms, {abs(hrv_diff):.0f}% below baseline ({hrv_bonus})")
    score += hrv_bonus
    
    # Resting HR (20% weight)
    rhr = metrics.get('resting_hr', 60)
    rhr_baseline = metrics.get('rhr_baseline', 55)
    if rhr <= rhr_baseline:
        rhr_bonus = 10
        factors.append(f"Resting HR {rhr}bpm, at/below baseline (+{rhr_bonus})")
    elif rhr <= rhr_baseline + 5:
        rhr_bonus = 0
        factors.append(f"Resting HR {rhr}bpm, slightly elevated (+0)")
    else:
        rhr_bonus = -10
        factors.append(f"Resting HR {rhr}bpm, elevated ({rhr_bonus})")
    score += rhr_bonus
    
    # Previous day strain (15% weight)
    yesterday_strain = metrics.get('yesterday_strain', 50)
    if yesterday_strain < 30:
        strain_bonus = 5
        factors.append(f"Yesterday low strain (+{strain_bonus})")
    elif yesterday_strain > 70:
        strain_bonus = -5
        factors.append(f"Yesterday high strain ({strain_bonus})")
    else:
        strain_bonus = 0
    score += strain_bonus
    
    score = max(0, min(100, score))  # Clamp 0-100
    explanation = f"Readiness: {score}/100 — " + ", ".join(factors)
    
    return score, explanation
```

### Behavioral Archetypes

> **Concept:** Sahha assigns labels like "highly social" or "morning person" based on patterns. We'll build dynamic archetypes that Claude can reference and explain.

| Archetype | Detection Logic | Example Insight |
| --------- | --------------- | --------------- |
| **Morning Person** | 80%+ days wake before 7am, best productivity before noon | "You're a morning person—your task completion rate is 45% higher before noon." |
| **Night Owl** | 80%+ days sleep after midnight, peak productivity 6-10pm | "Your night owl pattern shows 2x productivity after 8pm vs morning." |
| **Streak Builder** | Maintains 5+ day streaks regularly | "You thrive on streaks. Current: 12 days exercise, 8 days meditation." |
| **Weekend Warrior** | 3x+ activity on weekends vs weekdays | "Your weekend warrior pattern: 15k avg weekend steps vs 5k weekdays." |
| **Social Energizer** | Mood correlates positively with social events | "Social events boost your mood by avg +1.5 points the next day." |
| **Solo Recharger** | Productivity higher after low-social days | "Your best work days follow solo recharge days—protect your alone time." |
| **Consistent Sleeper** | <30min std deviation in sleep/wake times | "Your consistent 11pm-6:30am sleep window correlates with 20% higher HRV." |
| **Variable Sleeper** | >90min std deviation in sleep/wake times | "Irregular sleep pattern detected. Days with consistent sleep: +15% productivity." |
| **Recovery Responder** | Clear performance boost after rest days | "You respond well to recovery—day after rest day shows +25% output." |
| **High-Volume Achiever** | Completes 15+ tasks daily consistently | "High-volume achiever: averaging 18 tasks/day, 40% above your 90-day norm." |

### Archetype Detection Example

```python
def detect_archetypes(metrics_30_days: list[dict]) -> list[dict]:
    """
    Analyze 30 days of data to detect behavioral archetypes.
    Returns list of {archetype, confidence, evidence, insight}
    """
    archetypes = []
    
    # Morning Person Detection
    wake_times = [m['sleep_end'] for m in metrics_30_days if m.get('sleep_end')]
    early_wakes = sum(1 for t in wake_times if t < 420)  # Before 7am (420 mins)
    if len(wake_times) > 0:
        morning_ratio = early_wakes / len(wake_times)
        if morning_ratio >= 0.8:
            archetypes.append({
                'archetype': 'Morning Person',
                'confidence': morning_ratio,
                'evidence': f"Woke before 7am on {early_wakes}/{len(wake_times)} days",
                'insight': "Your best work aligns with early hours. Protect morning focus time."
            })
    
    # Streak Builder Detection
    current_streaks = calculate_active_streaks(metrics_30_days)
    long_streaks = [s for s in current_streaks if s['length'] >= 5]
    if len(long_streaks) >= 2:
        archetypes.append({
            'archetype': 'Streak Builder',
            'confidence': min(1.0, len(long_streaks) / 3),
            'evidence': f"Maintaining {len(long_streaks)} active streaks of 5+ days",
            'insight': "Streaks motivate you. Consider adding new streak goals."
        })
    
    # Social Energizer vs Solo Recharger
    social_mood_correlation = calculate_lagged_correlation(
        metrics_30_days, 'social_events', 'mood', lag=1
    )
    if social_mood_correlation > 0.4:
        archetypes.append({
            'archetype': 'Social Energizer',
            'confidence': social_mood_correlation,
            'evidence': f"Social events → next-day mood correlation: {social_mood_correlation:.2f}",
            'insight': "Social connection boosts your energy. Schedule regular social time."
        })
    elif social_mood_correlation < -0.3:
        archetypes.append({
            'archetype': 'Solo Recharger',
            'confidence': abs(social_mood_correlation),
            'evidence': f"Low-social days → better mood/productivity next day",
            'insight': "You recharge alone. Protect recovery time after social events."
        })
    
    return archetypes
```

### Legacy Conversation Examples

**Morning Check-in with Scores:**
```
Legacy: "Good morning, Zeus. Your Readiness Score is 85/100 today:
• Sleep 7.8hrs (+20) ✅
• HRV 58ms, 12% above baseline (+15) ✅
• Resting HR 52bpm, at baseline (+10) ✅
• Yesterday moderate strain (+0)

You're primed for a high-output day. Your Morning Person archetype suggests 
front-loading deep work before noon. Shall I block your morning for focus time?"
```

**Weekly Archetype Update:**
```
Zeus: "Legacy, how am I doing this week?"

Legacy: "This week's behavioral patterns show:

**Active Archetypes:**
• 🌅 Morning Person (92% confidence) — Avg wake: 6:15am
• 🔥 Streak Builder (3 active streaks) — Exercise: 14 days, Meditation: 7 days, Todoist: 21 days
• 💪 Recovery Responder — Your Monday rest day boosted Tuesday output by 35%

**Emerging Pattern:**
You're shifting toward a more consistent sleep schedule (std dev dropped from 
45min to 28min over 2 weeks). If maintained, expect +10-15% HRV improvement.

**Recommendation:** Continue protecting Sunday/Monday as recovery, then push 
hard Tuesday-Friday. Your data shows this rhythm works best for you."
```

**Score-Based Decision Support:**
```
Zeus: "Should I do a hard workout today?"

Legacy: "Based on your current scores:
• Readiness: 62/100 (below your 75 threshold)
• Recovery: 55/100 (HRV trending down 3 days)
• Yesterday's Strain: 78/100 (high)

**Recommendation:** Light activity only today. Your Recovery Responder 
archetype shows you'll get 30% more from tomorrow's workout if you rest today.

Alternative: 30min walk or yoga instead of HIIT. Want me to adjust your 
workout plan?"
```

---

## Health Data Aggregator Evaluation

> **TL;DR:** Terra is the best aggregator, but $499/mo is overkill for 1 user. Build direct integrations for personal use.

### Terra API Deep Dive

| Aspect | Details |
| ------ | ------- |
| **Pricing** | $499/mo (Quick Start) or $399/mo (annual) |
| **Credits** | 100,000/mo included, $0.005/credit overage |
| **Free Tier** | ❌ None |
| **Startup Program** | ❌ Not publicly advertised |
| **Solo/Indie Tier** | ❌ None found |
| **30-Day Refund** | ✅ Yes, money-back guarantee |
| **Usage Spike Protection** | ✅ First 6 months, one-time refund eligible |

### Cost Analysis for Personal Use

| Approach | Monthly Cost | Annual Cost | Notes |
| -------- | ------------ | ----------- | ----- |
| **Terra Quick Start** | $499 | $5,988 | Overkill for 1 user |
| **Terra Annual** | $399 | $4,788 | Still expensive |
| **Direct APIs** | $0 | $0 | All personal use APIs are free |
| **ROOK** | ~$100-200? | ~$1,200-2,400 | Contact for pricing |
| **Sahha** | $0 (Free tier) | $0 | But black-box AI |

### Recommendation: Direct Integration Strategy

For a **single user** (Zeus), direct API integrations are **100% free**:

| Data Source | API Cost | Complexity | Priority |
| ----------- | -------- | ---------- | -------- |
| **Oura Ring** | Free (personal) | Low | P0 |
| **Apple Health** | Free (HealthKit) | Medium | P1 |
| **Todoist** | Free | Low | P0 |
| **GitHub** | Free | Low | P1 |
| **Google Calendar** | Free | Low | P1 |
| **Weather API** | Free (OpenWeather) | Low | P2 |
| **RescueTime** | Free (limited) | Low | P2 |

### When to Reconsider Terra

| Scenario | Recommendation |
| -------- | -------------- |
| Adding 10+ wearable types | Consider ROOK first ($100-200/mo) |
| Need real-time BLE streaming | Terra is the only option |
| Building for multiple users | Terra makes sense at scale |
| Want pre-built health scores | Build custom (better explainability) |

### The Verdict

**For OneMind Insights (personal use):**
1. ✅ Direct integrations (Oura, Apple Health, Todoist, GitHub) — **FREE**
2. ✅ Build custom health scores with Claude — **Better than Sahha's black box**
3. ❌ Skip Terra until scaling to multiple users or needing real-time streaming
4. 📝 Re-evaluate in Q3 2026 if adding 5+ more wearable types

---

## Related Documentation

- [ZEUSRPG-GAMIFICATION.md](./docs/ZEUSRPG-GAMIFICATION.md) - Gamification system
- [TODOIST-EXISTIO-INTEGRATION.md](./docs/integrations/TODOIST-EXISTIO-INTEGRATION.md) - Current integration plan
- [HABITICA-INTEGRATION.md](./docs/integrations/HABITICA-INTEGRATION.md) - Deferred for Q3 2026
- [ONEMIND-INSIGHTS-ARCHITECTURE.md](../OneMind-Codex/ONEMIND-INSIGHTS-ARCHITECTURE.md) - Architecture deep dive

---

*"Know thyself—with data."*
