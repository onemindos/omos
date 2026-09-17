# Family Profiles - Multi-User Learning System

> **Created:** 2026-01-25
> **Status:** DEPLOYED TO PRODUCTION

---

## Overview

OneMind OS now supports **family profiles** - allowing Legacy to learn about and remember information about each family member. This uses Agno's native multi-user support via the `user_id` parameter.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    FAMILY PROFILE SYSTEM                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   learning_profiles                      │    │
│  │  ┌─────────┬─────────┬─────────┬─────────┬─────────┐   │    │
│  │  │  zeus   │  noah   │  liam   │ soluna  │  cielo  │   │    │
│  │  │ (father)│ (child) │ (child) │ (child) │ (child) │   │    │
│  │  └─────────┴─────────┴─────────┴─────────┴─────────┘   │    │
│  └─────────────────────────────────────────────────────────┘    │
│                              │                                   │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              learning_entities (Knowledge Graph)         │    │
│  │                                                          │    │
│  │  Zeus's entities:                                        │    │
│  │  ├── Noah Delacruz (person, child)                      │    │
│  │  ├── Liam Delacruz (person, child)                      │    │
│  │  ├── Soluna Delacruz (person, child)                    │    │
│  │  └── Cielo Delacruz (person, child)                     │    │
│  │                                                          │    │
│  │  Facts about each (learning_entity_facts):              │    │
│  │  ├── "Noah likes dinosaurs"                             │    │
│  │  ├── "Liam's favorite color is blue"                    │    │
│  │  └── etc. (added through conversations)                 │    │
│  └─────────────────────────────────────────────────────────┘    │
│                              │                                   │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │           learning_culture_principles                    │    │
│  │                                                          │    │
│  │  Legacy's principle:                                     │    │
│  │  "Remember and care about the Delacruz family:          │    │
│  │   Noah, Liam, Soluna, Cielo" (priority: 9)              │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Database Schema

### learning_profiles

Stores profile for each family member:

```sql
CREATE TABLE learning_profiles (
    user_id VARCHAR(255) PRIMARY KEY,    -- 'zeus', 'noah', 'liam', etc.
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    timezone VARCHAR(100) DEFAULT 'America/New_York',
    four_pillars JSONB DEFAULT '["HP", "LE", "GE", "IT"]',
    preferences JSONB,                    -- role, communication style, etc.
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### learning_entities

Knowledge graph entries for people Legacy learns about:

```sql
CREATE TABLE learning_entities (
    id UUID PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,       -- Owner (zeus)
    type VARCHAR(50) NOT NULL,           -- 'person', 'project', etc.
    name VARCHAR(255) NOT NULL,          -- 'Noah Delacruz'
    description TEXT,                     -- 'Child of Zeus Delacruz'
    pillar VARCHAR(10),                  -- 'le' (Life & Experience)
    metadata JSONB,                      -- {"relationship": "child", "family": "delacruz"}
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### learning_entity_facts

Facts Legacy learns about family members:

```sql
CREATE TABLE learning_entity_facts (
    id UUID PRIMARY KEY,
    entity_id UUID REFERENCES learning_entities(id),
    fact TEXT NOT NULL,                  -- "Noah likes dinosaurs"
    source VARCHAR(50) DEFAULT 'conversation',
    confidence FLOAT DEFAULT 1.0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Data in Production

### Family Profiles

| user_id | name | role |
|---------|------|------|
| zeus | Zeus Delacruz | father |
| noah | Noah Delacruz | child |
| liam | Liam Delacruz | child |
| soluna | Soluna Delacruz | child |
| cielo | Cielo Delacruz | child |

### Zeus's Profile Fields

| key | value | pillar |
|-----|-------|--------|
| role | "One Mind Operator" | - |
| focus_time | "5:30-9:00 AM" | hp |
| family_day | "Sunday" | le |
| decision_style | "data-driven" | - |

### Family Entities (Knowledge Graph)

| name | type | pillar | relationship |
|------|------|--------|--------------|
| Noah Delacruz | person | le | child |
| Liam Delacruz | person | le | child |
| Soluna Delacruz | person | le | child |
| Cielo Delacruz | person | le | child |

### Legacy's Culture Principles

| key | value | priority |
|-----|-------|----------|
| family_aware | Remember and care about the Delacruz family: Noah, Liam, Soluna, Cielo | 9 |

---

## How Agno Uses This

### Multi-User Memory Isolation

```python
# Each user_id gets isolated memories
agent.print_response("I love dinosaurs", user_id="noah")
agent.print_response("I like coding", user_id="zeus")

# Memories are isolated - noah's memories don't mix with zeus's
noah_memories = agent.get_user_memories(user_id="noah")
zeus_memories = agent.get_user_memories(user_id="zeus")
```

### Learning Facts About Family

When you talk about your kids, Legacy stores facts:

```
Zeus: "Noah had his first soccer game today, he scored 2 goals!"

Legacy stores:
- Entity: Noah Delacruz
- Fact: "Had first soccer game, scored 2 goals"
- Source: conversation
- Date: 2026-01-25
```

### Querying Family Info

```python
# Get all facts about Noah
facts = await get_entity_facts(entity_name="Noah Delacruz")

# Search across family
results = await search_entities(
    user_id="zeus",
    query="soccer",
    type="person"
)
```

---

## API Endpoints

### Get Family Profiles

```bash
curl http://100.102.21.44:8081/api/learning/profile?user_id=zeus
curl http://100.102.21.44:8081/api/learning/profile?user_id=noah
```

### Get Family Entities

```bash
# All entities for Zeus
curl http://100.102.21.44:8081/api/learning/entities?user_id=zeus

# Filter by type
curl "http://100.102.21.44:8081/api/learning/entities?user_id=zeus&type=person"
```

### Add Facts About Family Members

```bash
# Add a fact about Noah
curl -X POST http://100.102.21.44:8081/api/learning/entities/{noah_entity_id}/facts \
  -H "Content-Type: application/json" \
  -d '{"fact": "Loves building with Legos", "source": "conversation"}'
```

---

## Future Enhancements

| Feature | Description |
|---------|-------------|
| **Entity Relations** | Link family members (noah → sibling_of → liam) |
| **Automatic Extraction** | Legacy auto-creates facts from conversations |
| **Family Calendar** | Shared events and reminders |
| **Growth Tracking** | Track milestones, achievements |
| **Photo Memories** | Link photos to entity facts |

---

## Files Modified

| File | Changes |
|------|---------|
| `infra/configs/postgres/init-db.sql` | Added family profiles, entities, culture principle |
| `docs/architecture/FAMILY-PROFILES.md` | This documentation |

---

## Verification Queries

```sql
-- View all family profiles
SELECT user_id, name, preferences->>'role' as role
FROM learning_profiles ORDER BY user_id;

-- View family in knowledge graph
SELECT name, description, pillar
FROM learning_entities
WHERE metadata->>'family' = 'delacruz';

-- View Legacy's family awareness
SELECT key, value FROM learning_culture_principles
WHERE key = 'family_aware';
```

---

*Created for the Delacruz family - Zeus, Noah, Liam, Soluna, Cielo*
