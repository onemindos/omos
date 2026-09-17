# OpenMind.org Database Requirements

## Overview

OpenMind (OM1) is an AI-native operating system for autonomous machines (robots). While OpenMind.org doesn't prescribe a specific database, the architecture requires several types of data persistence.

---

## Core Storage Requirements

### 1. **Memory Module** (Short-term & Long-term Memory)

**What it stores:**
- Recent sensor observations (vision, audio, proprioception)
- Episodic memory (what happened and when)
- Semantic memory (learned facts and concepts)
- Vector embeddings for similarity search

**Database recommendation:**
```
PostgreSQL + pgvector

Why:
✅ Vector similarity search (pgvector)
✅ JSONB for flexible sensor data
✅ Transactional integrity
✅ Full-text search
✅ Proven reliability for robotics

Example schema:
CREATE TABLE robot_memories (
    id UUID PRIMARY KEY,
    timestamp TIMESTAMPTZ,
    type TEXT,  -- 'episodic', 'semantic', 'procedural'
    content TEXT,
    embedding vector(1536),  -- For similarity search
    metadata JSONB,
    created_at TIMESTAMPTZ
);

-- Find similar memories
SELECT * FROM robot_memories
ORDER BY embedding <=> query_vector
LIMIT 10;
```

---

### 2. **Planning Module** (Goals, Tasks, State)

**What it stores:**
- Current goals and sub-goals
- Task queue (pending actions)
- World state model
- Action outcomes (success/failure)

**Database recommendation:**
```
PostgreSQL (same as memory)

Why:
✅ ACID transactions (critical for planning)
✅ Complex queries for goal hierarchies
✅ Foreign keys for task dependencies

Example schema:
CREATE TABLE goals (
    id UUID PRIMARY KEY,
    parent_goal_id UUID REFERENCES goals(id),
    description TEXT,
    status TEXT,  -- 'pending', 'active', 'completed', 'failed'
    priority INT,
    created_at TIMESTAMPTZ
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    goal_id UUID REFERENCES goals(id),
    action_type TEXT,
    parameters JSONB,
    status TEXT,
    result JSONB,
    executed_at TIMESTAMPTZ
);
```

---

### 3. **Sensor Data Streams** (Real-time & Historical)

**What it stores:**
- Camera frames (images/video)
- LiDAR point clouds
- IMU readings (acceleration, gyroscope)
- Motor encoders (position, velocity)
- Environmental sensors (temperature, etc.)

**Database recommendation:**
```
TimescaleDB (time-series extension for PostgreSQL)

Why:
✅ Optimized for high-frequency sensor data
✅ Automatic partitioning by time
✅ Compression (90%+ savings)
✅ time_bucket() for aggregations
✅ Retention policies (auto-delete old data)

Example schema:
CREATE TABLE sensor_data (
    time TIMESTAMPTZ NOT NULL,
    sensor_id TEXT NOT NULL,
    sensor_type TEXT,  -- 'camera', 'lidar', 'imu', etc.
    value JSONB,
    metadata JSONB
);

SELECT create_hypertable('sensor_data', 'time');

-- Aggregate sensor readings
SELECT
    time_bucket('1 second', time) as bucket,
    sensor_id,
    avg((value->>'x')::float) as avg_x
FROM sensor_data
WHERE sensor_type = 'imu'
  AND time > now() - INTERVAL '1 hour'
GROUP BY bucket, sensor_id;
```

---

### 4. **Event Bus** (Real-time Communication)

**What it stores:**
- Inter-module messages
- Sensor → Perception → Intelligence → Action pipeline
- Robot-to-robot communication (FABRIC)

**Database recommendation:**
```
NATS JetStream (NOT a database, but critical for OpenMind)

Why:
✅ Real-time pub/sub (microsecond latency)
✅ Persistent streams (replay events)
✅ MQTT bridge (IoT sensors)
✅ Leaf nodes (distributed robots)
✅ Exactly-once delivery

Your implementation already uses this! ✅
```

---

### 5. **SLAM & Mapping** (Spatial Data)

**What it stores:**
- 3D maps (occupancy grids, point clouds)
- Localization data (where am I?)
- Landmarks (recognized features)

**Database recommendation:**
```
PostgreSQL + PostGIS (spatial extension)

Why:
✅ Geospatial queries (distance, containment)
✅ 3D geometry support
✅ Path planning queries
✅ Integration with ROS

Example schema:
CREATE EXTENSION postgis;

CREATE TABLE maps (
    id UUID PRIMARY KEY,
    name TEXT,
    created_at TIMESTAMPTZ,
    bounds GEOMETRY(POLYGON, 4326)
);

CREATE TABLE landmarks (
    id UUID PRIMARY KEY,
    map_id UUID REFERENCES maps(id),
    position GEOMETRY(POINT, 4326),
    type TEXT,
    confidence FLOAT
);

-- Find nearby landmarks
SELECT * FROM landmarks
WHERE ST_DWithin(position, robot_position, 10.0);  -- Within 10 meters
```

---

## OpenMind.org Official Stack (Robotics)

Based on the OM1 architecture, the typical robotics stack is:

```
┌──────────────────────────────────────────────────┐
│         OPENMIND ROBOTICS STACK                   │
├──────────────────────────────────────────────────┤
│                                                   │
│  Real-time Layer (Milliseconds)                   │
│  ─────────────────────────────                    │
│  • NATS JetStream (event bus)                    │
│  • Redis (cache, coordination)                   │
│  • ROS topics (sensor streams)                   │
│                                                   │
│  Persistence Layer (Seconds)                      │
│  ──────────────────────────                       │
│  • PostgreSQL + pgvector (memory, planning)      │
│  • TimescaleDB (sensor history)                  │
│  • PostGIS (SLAM, maps)                          │
│                                                   │
│  Archival Layer (Minutes/Hours)                   │
│  ────────────────────────────                     │
│  • S3/MinIO (video recordings)                   │
│  • Parquet (compressed sensor logs)              │
│                                                   │
└──────────────────────────────────────────────────┘
```

---

## Your "Digital-First" Implementation

You're doing a **digital adaptation** of OpenMind (calendar/email instead of cameras/sensors).

### Your Current Database Needs:

| OpenMind Component | Your Digital Equivalent | Database Need |
|--------------------|------------------------|---------------|
| **Memory Module** | Agent memories, embeddings | ✅ PostgreSQL + pgvector |
| **Planning Module** | Agent goals, tasks | ✅ PostgreSQL |
| **Sensor Streams** | Webhook events, API data | ⚠️ Events table (decide: TimescaleDB or Parquet) |
| **Event Bus** | NATS JetStream | ✅ Already using NATS |
| **SLAM/Mapping** | N/A (digital, no spatial) | ❌ Not needed |

---

## Recommended Database Stack for YOU

Based on your digital-first OpenMind implementation:

### **Minimal Stack (Getting Started)**

```yaml
PostgreSQL 16 + pgvector:
  - Agent configs, teams, users
  - Knowledge bases (RAG)
  - Memories (vector search)
  - Sessions, approvals
  - Events (last 7 days, partitioned)

NATS JetStream:
  - Real-time event bus
  - Heartbeat loop communication
  - Webhook → Agent → Action pipeline

Cost: $25-33/month (App Platform)
```

### **Production Stack (Scaling)**

```yaml
PostgreSQL 16 + pgvector:
  - Transactional data (same as above)

TimescaleDB:
  - Event archival (90 days, compressed)
  - Sensor data (if adding IoT later)
  - Analytics (continuous aggregates)

NATS JetStream:
  - Real-time event bus

Redis (optional):
  - Rate limiting, caching
  - Real-time counters

Cost: $24/month (Droplet) + $5/month (Redis Cloud)
```

---

## Comparison: Robotics vs Digital

| Aspect | OpenMind Robotics | Your Digital Adaptation |
|--------|-------------------|------------------------|
| **Sensor Data** | Camera (30 FPS) = 1.8M frames/day | Webhooks = 100 events/day |
| **Data Volume** | 100 GB+/day (video, LiDAR) | 20 MB/day (JSON events) |
| **Latency** | Milliseconds (real-time control) | Seconds (API calls OK) |
| **Database** | TimescaleDB required (high volume) | PostgreSQL sufficient (low volume) |
| **SLAM** | PostGIS required (spatial queries) | Not needed (no physical space) |
| **Memory** | pgvector (embeddings) | ✅ Same (pgvector) |

**Key Insight:** Robotics OpenMind needs TimescaleDB due to massive sensor data. Your digital version can start with regular PostgreSQL.

---

## Decision Matrix for YOU

### Use PostgreSQL (Regular) if:
- ✅ Event volume < 500K/day
- ✅ Want App Platform managed hosting
- ✅ Avoid version lock-in
- ✅ Budget-conscious ($25-33/mo)

### Add TimescaleDB if:
- ✅ Event volume > 500K/day
- ✅ Planning to add IoT sensors (ESP32, etc.)
- ✅ Need automatic compression/retention
- ✅ Willing to self-host ($24/mo Droplet)

### Add PostGIS if:
- ✅ Adding physical robots (location tracking)
- ✅ Need spatial queries (distance, path planning)
- ❌ NOT needed for digital-only implementation

---

## Summary

**OpenMind.org (Robotics) Needs:**
1. PostgreSQL + pgvector (memory, planning)
2. TimescaleDB (sensor data at scale)
3. PostGIS (SLAM, spatial queries)
4. NATS JetStream (real-time event bus)
5. S3/MinIO (video archival)

**Your Digital OpenMind Needs:**
1. ✅ PostgreSQL + pgvector (memory, planning, events)
2. ✅ NATS JetStream (event bus)
3. ⚠️ TimescaleDB (optional, only if scaling to 500K+ events/day)
4. ❌ PostGIS (not needed for digital)

**Start simple, scale when needed.**

---

## Sources

- [OpenMind Architecture Overview](https://docs.openmind.org/full_autonomy_guidelines/architecture_overview)
- [OpenMind Explained (Medium)](https://medium.com/@takerwoods/openmind-explained-how-an-open-os-and-trust-network-could-reshape-robotics-eae0a9609cf8)
- [OM1: OpenMind's Open-Source OS](https://humanoidroboticstechnology.com/industry-news/openmind-launches-om1-beta-an-open-source-os-for-intelligent-robots/)
- [OpenMind wants to be the Android of humanoid robots (TechCrunch)](https://techcrunch.com/2025/08/04/openmind-wants-to-be-the-android-operating-system-of-humanoid-robots/)

---

_Last updated: 2026-02-11_
