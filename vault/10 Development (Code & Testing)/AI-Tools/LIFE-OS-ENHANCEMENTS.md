# Life OS Cherry-Pick Enhancements

This document describes the enhancements added to Life OS, inspired by open-source habit trackers (Loop Habits/uhabits, Habo) and recipe managers (Mealie, Tandoor).

---

## 1. Skip vs Miss Distinction (Loop Habits Pattern)

### Overview
Instead of binary completion tracking (done/not done), habits now support three states:

| Status | Effect | Use Case |
|--------|--------|----------|
| **Completed** | Counts toward streak and total | Normal completion |
| **Skipped** | Doesn't break streak, doesn't count | Rest day, intentional break |
| **Missed** | Breaks streak | Forgot/failed |

### API Endpoints

```bash
# Complete a habit
POST /api/habits/{habit_id}/complete
Body: {
  "status": "completed",  # or "skipped" or "missed"
  "notes": "Optional notes"
}

# Skip a habit (shortcut)
POST /api/habits/{habit_id}/skip?notes=Rest%20day
```

### Backend Changes
- `backend/lifeos/habits/models.py`: Added `CompletionStatus` enum
- `backend/lifeos/habits/service.py`: Updated `_update_streak()` to respect skip status
- `backend/lifeos/habits/api.py`: Added `/skip` endpoint

### Frontend Changes
- `frontend/lib/shared/models/habit.dart`: Added `CompletionStatus` and `HabitMood` enums
- `frontend/lib/lifeos/habits/providers/habits_provider.dart`: Added `skipHabit()` method

### Database Schema
```sql
ALTER TABLE habit_completions
ADD COLUMN status VARCHAR(20) DEFAULT 'completed'
  CHECK (status IN ('completed', 'skipped', 'missed'));
```

---

## 2. Mood/Energy Tracking (Habo Pattern)

### Overview
When completing habits, users can optionally track their emotional and physical state:

- **Mood**: great, good, okay, bad, terrible
- **Energy Level**: 1-10 scale
- **Intensity**: 1-5 (how hard was it)

### API Usage

```bash
POST /api/habits/{habit_id}/complete
Body: {
  "status": "completed",
  "mood": "good",
  "energy_level": 7,
  "intensity": 3,
  "notes": "Felt strong today"
}
```

### Use Cases
- Identify patterns (e.g., "I always skip workouts when energy < 5")
- Track mood over time to see habit impact
- Adjust habit difficulty based on intensity trends

### Frontend
The provider now accepts mood/energy when completing:
```dart
notifier.completeHabit(
  habitId,
  mood: HabitMood.good,
  energyLevel: 7,
  intensity: 3,
);
```

---

## 3. Habit Reminders (Notifications Integration)

### Overview
Habits with `reminder_enabled = true` and a `target_time` will trigger push notifications at the scheduled time minus `reminder_minutes_before`.

### How It Works
1. A background job checks habits periodically
2. If current time is within the reminder window, sends notification
3. User can **Complete**, **Skip**, or **Snooze** from the notification

### API Endpoints

```bash
# Manually trigger reminder check (for testing)
POST /api/habits/reminders/check
# Returns: {"reminders_sent": 3}

# Snooze a reminder
POST /api/habits/{habit_id}/snooze?minutes=15
```

### Notification Content
```
Title: "Habit Reminder: Morning Exercise"
Body: "Time to morning exercise! (15 min until scheduled time)"
Actions: [Complete] [Skip] [Snooze 15m]
```

### Setup
1. Ensure habits have `reminder_enabled: true` and `target_time` set
2. Configure notification channels in user preferences
3. Call `/api/habits/reminders/check` periodically (e.g., every minute via cron)

### Database Schema
```sql
CREATE TABLE habit_reminders_sent (
    id UUID PRIMARY KEY,
    habit_id UUID REFERENCES habits(id),
    sent_date DATE NOT NULL,
    sent_at TIMESTAMPTZ,
    snoozed_until TIMESTAMPTZ,
    UNIQUE(habit_id, sent_date)
);
```

---

## 4. Meal Plan Templates (Mealie/Tandoor Pattern)

### Overview
Save and reuse weekly meal patterns. Create once, apply to any week.

### Features
- **Create templates** manually or from existing meal plan
- **Apply templates** to any start date
- **Track usage** (times used, last used)
- **Duration**: 1-28 days

### API Endpoints

```bash
# List templates
GET /api/meals/templates

# Create template
POST /api/meals/templates
Body: {
  "name": "Healthy Week",
  "description": "Balanced meals for a healthy week",
  "duration_days": 7,
  "tags": ["healthy", "balanced"],
  "meals": [
    {"relative_day": 0, "meal_type": "breakfast", "recipe_id": "uuid..."},
    {"relative_day": 0, "meal_type": "lunch", "custom_meal": "Salad"},
    ...
  ]
}

# Apply template to meal plan
POST /api/meals/templates/{template_id}/apply
Body: {
  "start_date": "2025-02-03",
  "overwrite_existing": false
}
# Returns: {"meals_created": 21, "meals_skipped": 0}

# Create template from current meal plan
POST /api/meals/templates/from-plan?name=My%20Week&start_date=2025-01-27
```

### Template Meal Structure
```json
{
  "relative_day": 0,  // 0 = Day 1 (usually Monday)
  "meal_type": "dinner",
  "recipe_id": "uuid...",  // OR
  "custom_meal": "Takeout",
  "servings": 2,
  "notes": "Family favorite"
}
```

### Database Schema
```sql
CREATE TABLE meal_plan_templates (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    duration_days INT DEFAULT 7,
    tags TEXT[],
    meals JSONB NOT NULL DEFAULT '[]',
    is_default BOOLEAN DEFAULT false,
    times_used INT DEFAULT 0,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
```

---

## Summary of New Files

### Backend
| File | Description |
|------|-------------|
| `backend/lifeos/habits/reminders.py` | Habit reminder service |
| `backend/lifeos/meals/models.py` | Added template models |
| `backend/lifeos/meals/service.py` | Added template methods |
| `backend/lifeos/meals/api.py` | Added template endpoints |

### Database
| Change | Location |
|--------|----------|
| `habit_completions.status` | Skip/Miss tracking |
| `habit_completions.mood/energy` | Mood tracking |
| `habit_reminders_sent` | Reminder tracking |
| `meal_plan_templates` | Template storage |

### Frontend
| File | Changes |
|------|---------|
| `habit.dart` | Added `CompletionStatus`, `HabitMood` |
| `habits_provider.dart` | Added `skipHabit()`, mood/energy params |

---

## Testing

### Skip/Miss
```bash
# Complete normally
curl -X POST "http://localhost:8080/api/habits/{id}/complete" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'

# Skip
curl -X POST "http://localhost:8080/api/habits/{id}/skip"

# Verify streak preserved after skip
curl "http://localhost:8080/api/habits/{id}"
```

### Mood/Energy
```bash
curl -X POST "http://localhost:8080/api/habits/{id}/complete" \
  -H "Content-Type: application/json" \
  -d '{"mood": "good", "energy_level": 8}'
```

### Reminders
```bash
# Create habit with reminder
curl -X POST "http://localhost:8080/api/habits?user_id=zeus" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Morning Run",
    "target_time": "07:00:00",
    "reminder_enabled": true,
    "reminder_minutes_before": 15
  }'

# Trigger check
curl -X POST "http://localhost:8080/api/habits/reminders/check"
```

### Templates
```bash
# Create template
curl -X POST "http://localhost:8080/api/meals/templates?user_id=zeus" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Quick Weeknight",
    "duration_days": 5,
    "meals": [
      {"relative_day": 0, "meal_type": "dinner", "custom_meal": "Pasta"}
    ]
  }'

# Apply to next week
curl -X POST "http://localhost:8080/api/meals/templates/{id}/apply?user_id=zeus" \
  -H "Content-Type: application/json" \
  -d '{"start_date": "2025-02-03"}'
```

---

## Licensing Note

All source projects (Loop Habits/uhabits, Habo, Mealie, Tandoor) are AGPL/GPL licensed. For personal use with ported algorithms (not copied code), this is completely fine. No distribution = no license obligations.

---

## 5. Bridge Module - Unified Context & Knowledge Graph (NEW)

### Overview

The Bridge module connects all LifeOS features to the Learning Machine, enabling AI agents to:
- Get unified context for any entity (tasks, habits, goals, children, etc.)
- Create relationships between entities (task supports goal, task assigned to child)
- Record facts that accumulate over time for pattern detection
- Calculate Intelligence Scores (PI/LI/WI) from real data
- Detect patterns (best days, energy correlations, estimation accuracy)
- Predict task/habit completion likelihood

### Key Features

| Feature | Description |
|---------|-------------|
| **Unified Context** | Single API call returns LifeOS data + facts + relations + predictions |
| **Knowledge Graph** | Connect any entity to any other with typed relations |
| **Pattern Detection** | Automatically find completion patterns and correlations |
| **Predictions** | Predict completion likelihood based on patterns |
| **Intelligence Scores** | PI (Performance), LI (Legacy), WI (Wealth) from real metrics |
| **Family Support** | Multi-user with parent/child relationships |
| **Auto-Registration** | LifeOS items automatically become Learning entities |

### Documentation

Full documentation available at:
- **Complete Guide**: [BRIDGE-MODULE-DOCUMENTATION.md](./BRIDGE-MODULE-DOCUMENTATION.md)
- **Quick Reference**: [BRIDGE-QUICK-REFERENCE.md](./BRIDGE-QUICK-REFERENCE.md)

### API Endpoints

Base URL: `/api/bridge`

| Method | Path | Description |
|--------|------|-------------|
| GET | `/context/{type}/{id}` | Get unified entity context |
| POST | `/register` | Register entity |
| POST | `/link` | Create relationship |
| POST | `/fact/{type}/{id}` | Add fact |
| GET | `/scores` | Get PI/LI/WI scores |
| GET | `/patterns` | Get user patterns |
| GET | `/predictions/{date}` | Get day predictions |
| GET | `/family` | Get family overview |
| POST | `/backfill` | Backfill existing data |

### AI Tools

7 new tools for Legacy agent:
- `unified_context` - Get complete context
- `link_entities` - Create relationships
- `add_entity_context` - Add facts
- `get_intelligence_scores` - Get PI/LI/WI
- `get_day_predictions` - Get predictions
- `get_user_patterns` - Get patterns
- `get_family_overview` - Get family

### Files Created

```
backend/lifeos/bridge/
├── __init__.py      # Module exports (~30 lines)
├── models.py        # Pydantic schemas (~310 lines)
├── service.py       # Core logic (~770 lines)
├── api.py           # REST endpoints (~410 lines)
└── triggers.py      # Auto-registration (~450 lines)
```

### Integration

Uses existing tables (no new migrations required):
- `learning_entities` - Entity storage
- `learning_entity_facts` - Facts with confidence
- `learning_entity_relations` - Knowledge graph edges
- `learning_profiles` - User/family profiles
