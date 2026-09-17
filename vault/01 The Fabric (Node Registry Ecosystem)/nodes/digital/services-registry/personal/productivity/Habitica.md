# Habitica Integration Strategy for OneMind OS

> **Status:** Research Complete | Future Consideration
> **Priority:** Low (Custom gamification preferred)
> **Decision:** Defer in favor of Todoist + Exist.io + Custom Legacy Gamification
> **Review Date:** Q3 2026

---

## Executive Summary

Habitica is an open-source gamified habit tracking platform with extensive API capabilities. While it offers compelling RPG mechanics, the decision was made to defer integration in favor of:
1. **Todoist** - Already in use with 116+ tasks, native Exist.io integration
2. **Exist.io** - Health correlations and insights
3. **Custom Legacy Gamification** - Build exactly what we need without another task system

This document preserves all research for future custom development.

---

## Table of Contents

1. [Platform Overview](#platform-overview)
2. [Why Habitica Was Considered](#why-habitica-was-considered)
3. [Why Deferred](#why-deferred)
4. [Complete API Reference](#complete-api-reference)
5. [Integration Patterns](#integration-patterns)
6. [Custom Gamification Alternative](#custom-gamification-alternative)
7. [Community Extensions](#community-extensions)
8. [Self-Hosting Option](#self-hosting-option)
9. [Future Considerations](#future-considerations)

---

## Platform Overview

### What is Habitica?

Habitica is an open-source habit-building program that treats your life like an RPG:
- **Level up** as you complete tasks
- **Take damage** when you miss dailies
- **Earn gold** to buy equipment, pets, and mounts
- **Join parties** for social accountability
- **Battle bosses** through group quests

### Key Statistics

| Metric | Value |
|--------|-------|
| Active Users | 4+ million |
| GitHub Stars | 13.6k |
| Contributors | 946 |
| Open Source | Yes (GPL-3.0) |
| API Version | v3 (current) |
| Mobile Apps | iOS + Android |
| Price | Free (Premium $9/month for extra features) |

### Repository

- **Main:** https://github.com/HabitRPG/habitica
- **Android:** https://github.com/HabitRPG/habitica-android
- **iOS:** https://github.com/HabitRPG/habitica-ios
- **API Docs:** https://habitica.com/apidoc

---

## Why Habitica Was Considered

### Compelling Features

1. **Gamification Psychology**
   - Loss aversion (HP damage for missed dailies)
   - Progress visualization (XP bars, levels)
   - Variable rewards (random drops)
   - Social pressure (party accountability)

2. **Comprehensive Task System**
   - **Habits:** +/- buttons for tracking behaviors
   - **Dailies:** Recurring tasks with custom schedules
   - **To-Dos:** One-time tasks with checklists
   - **Rewards:** Custom incentives purchased with gold

3. **Open Source Advantage**
   - Full source code access
   - Self-hosting capability
   - Customizable mechanics
   - No vendor lock-in

4. **Extensive API**
   - 100+ endpoints
   - Webhooks for real-time events
   - Complete CRUD for all entities
   - Data export capabilities

---

## Why Deferred

### Decision Matrix

| Factor | Habitica | Todoist + Custom | Winner |
|--------|----------|------------------|--------|
| Already using | No | Yes (116 tasks) | Todoist |
| Exist.io integration | Manual | Native | Todoist |
| Cognitive overhead | Two task systems | One system | Todoist |
| Customization | Fork required | Full control | Custom |
| Gamification | Built-in | Must build | Habitica |
| Mobile themes | No | N/A | - |
| Maintenance | Another app | Integrated | Custom |

### Key Reasons

1. **Todoist Investment**
   - Already have 116 tasks created
   - Familiar workflow established
   - Account cleanup tasks due TOMORROW

2. **Exist.io Native Integration**
   - Todoist → Exist.io automatic sync
   - Habitica → Exist.io requires custom bridge
   - Correlations work out of the box

3. **One System Philosophy**
   - Two task systems = cognitive overhead
   - Context switching costs
   - Sync complexity

4. **Custom Control**
   - Build gamification exactly as desired
   - Integrate with smart home
   - Zeus-specific mechanics

---

## Complete API Reference

### Authentication

```python
import requests

HABITICA_API = "https://habitica.com/api/v3"

headers = {
    "x-api-user": "YOUR_USER_ID",
    "x-api-key": "YOUR_API_TOKEN",
    "x-client": "onemind-os-UUID"  # Required for third-party apps
}
```

### Task Endpoints

#### Get All Tasks
```python
GET /api/v3/tasks/user

# Query parameters:
# type: habits, dailys, todos, rewards, completedTodos
# dueDate: for computing nextDue field
```

#### Create Task
```python
POST /api/v3/tasks/user

{
    "text": "Complete API integration",
    "type": "todo",  # habit, daily, todo, reward
    "tags": ["uuid-of-tag"],
    "alias": "api-integration",  # for easy reference
    "notes": "Documentation and testing",
    "date": "2026-01-25",  # due date (todos only)
    "priority": 2,  # 0.1=Trivial, 1=Easy, 1.5=Medium, 2=Hard
    "attribute": "int",  # str, int, per, con
    "checklist": [
        {"text": "Read docs", "completed": false},
        {"text": "Write code", "completed": false}
    ],
    "collapseChecklist": false,
    "reminders": [
        {
            "id": "uuid",
            "startDate": "2026-01-24T09:00:00.000Z",
            "time": "2026-01-24T09:00:00.000Z"
        }
    ]
}
```

#### Daily-Specific Options
```python
{
    "type": "daily",
    "frequency": "weekly",  # daily, weekly, monthly, yearly
    "repeat": {
        "su": false, "m": true, "t": true, "w": true,
        "th": true, "f": true, "s": false
    },
    "everyX": 1,  # every X days/weeks/months
    "startDate": "2026-01-21T00:00:00.000Z",
    "streak": 0,
    "daysOfMonth": [1, 15],  # for monthly
    "weeksOfMonth": [1, 3]   # for monthly
}
```

#### Habit-Specific Options
```python
{
    "type": "habit",
    "up": true,    # enable + button
    "down": true   # enable - button
}
```

#### Score Task
```python
POST /api/v3/tasks/{taskId}/score/{direction}
# direction: "up" or "down"

# Response includes:
{
    "delta": 0.9746,  # score change
    "hp": 49.06,
    "mp": 37.20,
    "exp": 101.93,
    "gp": 77.09,
    "lvl": 19,
    "class": "rogue",
    "_tmp": {
        "drop": {...}  # if item dropped
    }
}
```

### User Endpoints

#### Get User Profile
```python
GET /api/v3/user

# Selective fields:
GET /api/v3/user?userFields=stats,items.mounts
```

#### Update User
```python
PUT /api/v3/user

{
    "stats.hp": 50,
    "preferences.hair.bangs": 3,
    "profile.name": "Zeus"
}
```

### Webhook Configuration

```python
POST /api/v3/user/webhook

{
    "url": "https://onemind-os.zeuslegacy.com/webhooks/habitica",
    "label": "OneMind OS Integration",
    "enabled": true,
    "type": "taskActivity",
    "options": {
        "created": true,
        "updated": true,
        "deleted": true,
        "scored": true,
        "checklistScored": true
    }
}

# Webhook types:
# - taskActivity: task events
# - groupChatReceived: group messages
# - userActivity: level up, achievements
# - questActivity: quest progress
```

### Tags System

```python
# Get all tags
GET /api/v3/tags

# Create tag
POST /api/v3/tags
{"name": "OneMind"}

# Add tag to task
POST /api/v3/tasks/{taskId}/tags/{tagId}
```

### Groups & Social

```python
# Get party
GET /api/v3/groups/party

# Post chat message
POST /api/v3/groups/{groupId}/chat
{"message": "Completed all dailies!"}

# Accept quest
POST /api/v3/groups/{groupId}/quests/accept
```

### Data Export

```python
# Full JSON export
GET /export/userdata.json

# Task history CSV
GET /export/history.csv

# Avatar PNG
GET /export/avatar-{uuid}.png
```

---

## Integration Patterns

### Pattern 1: Headless Habitica (API Only)

Use Habitica purely for RPG mechanics, never open the app:

```python
# agents/legacy/tools/habitica_headless.py

class HabiticaHeadless:
    """Use Habitica as a gamification engine only"""
    
    def __init__(self, user_id: str, api_token: str):
        self.headers = {
            "x-api-user": user_id,
            "x-api-key": api_token,
            "x-client": "onemind-os-legacy"
        }
        self.base_url = "https://habitica.com/api/v3"
    
    async def score_achievement(self, description: str, difficulty: float = 1.0):
        """Score a custom achievement in Habitica"""
        # Create temporary habit and score it
        habit = await self.create_task({
            "type": "habit",
            "text": description,
            "up": True,
            "down": False,
            "priority": difficulty
        })
        
        result = await self.score_task(habit['id'], 'up')
        await self.delete_task(habit['id'])
        
        return {
            "xp_gained": result['exp'] - self.last_exp,
            "gold_gained": result['gp'] - self.last_gp,
            "level_up": result['lvl'] > self.last_level,
            "item_drop": result.get('_tmp', {}).get('drop')
        }
```

### Pattern 2: Webhook-Driven Reactions

Let Legacy react to Habitica events:

```python
# agents/legacy/webhooks/habitica_webhook.py

@router.post("/webhooks/habitica")
async def handle_habitica_webhook(event: HabiticaEvent):
    """React to Habitica events in real-time"""
    
    if event.type == "scored" and event.direction == "up":
        task = event.task
        
        # Streak milestones
        if task.streak in [7, 21, 30, 60, 90, 365]:
            await legacy.celebrate_streak(task.text, task.streak)
            await smart_home.celebration_lights()
        
        # Level up detection
        if event.user.lvl > await get_cached_level():
            await legacy.announce_level_up(event.user.lvl)
            await discord.post_achievement(
                f"🎮 Zeus reached Level {event.user.lvl}!"
            )
    
    elif event.type == "scored" and event.direction == "down":
        # Gentle accountability
        await legacy.send_encouragement(event.task.text)
    
    elif event.type == "cron":
        # Daily reset - check yesterday's performance
        await legacy.daily_summary()
```

### Pattern 3: Bidirectional Sync

Keep Habitica and another system in sync:

```python
# agents/legacy/tools/habitica_sync.py

async def sync_todoist_to_habitica():
    """Mirror Todoist tasks to Habitica for gamification"""
    
    # Get Todoist tasks
    todoist_tasks = await todoist.get_tasks()
    
    for task in todoist_tasks:
        habitica_id = task.metadata.get('habitica_id')
        
        if not habitica_id:
            # Create in Habitica
            h_task = await habitica.create_task({
                "type": "todo",
                "text": task.content,
                "notes": f"Synced from Todoist: {task.id}",
                "priority": map_priority(task.priority)
            })
            await todoist.update_metadata(task.id, 
                habitica_id=h_task['id'])
        
        # Sync completion status
        if task.is_completed and not await habitica.is_completed(habitica_id):
            await habitica.score_task(habitica_id, 'up')
```

---

## Custom Gamification Alternative

Instead of Habitica, build custom gamification into Legacy:

### Zeus RPG System

```python
# agents/legacy/gamification/zeus_rpg.py

from dataclasses import dataclass
from datetime import datetime

@dataclass
class ZeusStats:
    level: int = 1
    xp: int = 0
    gold: int = 0
    health: int = 50
    max_health: int = 50
    
    # Attributes
    strength: int = 1      # Physical tasks
    intelligence: int = 1  # Learning/coding
    perception: int = 1    # Attention to detail
    constitution: int = 1  # Consistency/streaks
    
    # Streaks
    streaks: dict = field(default_factory=dict)
    
    @property
    def xp_to_next_level(self):
        return self.level * 100 + 50  # 150, 250, 350...

class ZeusRPG:
    """Custom gamification engine for OneMind OS"""
    
    XP_MULTIPLIERS = {
        "trivial": 0.1,
        "easy": 1.0,
        "medium": 1.5,
        "hard": 2.0,
        "epic": 3.0
    }
    
    async def complete_task(self, task: dict, difficulty: str = "medium"):
        """Award XP and gold for task completion"""
        base_xp = 10
        base_gold = 5
        
        multiplier = self.XP_MULTIPLIERS.get(difficulty, 1.0)
        
        # Streak bonus
        streak = await self.get_streak(task['id'])
        streak_bonus = 1 + (streak * 0.1)  # +10% per day of streak
        
        xp_earned = int(base_xp * multiplier * streak_bonus)
        gold_earned = int(base_gold * multiplier)
        
        self.stats.xp += xp_earned
        self.stats.gold += gold_earned
        
        # Check level up
        level_up = await self.check_level_up()
        
        # Update streak
        await self.update_streak(task['id'])
        
        return {
            "xp_earned": xp_earned,
            "gold_earned": gold_earned,
            "level_up": level_up,
            "new_streak": streak + 1,
            "total_xp": self.stats.xp,
            "total_gold": self.stats.gold
        }
    
    async def miss_daily(self, task: dict):
        """Take damage for missing a daily"""
        damage = 10  # Base damage
        
        # Streak loss penalty
        streak = await self.get_streak(task['id'])
        if streak > 7:
            damage += streak  # Lose more for breaking long streaks
        
        self.stats.health -= damage
        await self.reset_streak(task['id'])
        
        if self.stats.health <= 0:
            await self.handle_death()
        
        return {
            "damage_taken": damage,
            "current_health": self.stats.health,
            "streak_lost": streak
        }
    
    async def check_level_up(self) -> bool:
        """Check and handle level up"""
        if self.stats.xp >= self.stats.xp_to_next_level:
            self.stats.xp -= self.stats.xp_to_next_level
            self.stats.level += 1
            self.stats.max_health += 5
            self.stats.health = self.stats.max_health
            
            # Stat points
            await self.award_stat_points(1)
            
            return True
        return False
    
    async def spend_gold(self, item: str, cost: int) -> bool:
        """Purchase rewards with gold"""
        if self.stats.gold >= cost:
            self.stats.gold -= cost
            await self.grant_reward(item)
            return True
        return False
```

### Legacy Integration

```python
# agents/legacy/tools/gamification_tools.py

@tool
async def complete_task_with_rewards(task_id: str):
    """Complete a task and show gamification results"""
    task = await todoist.get_task(task_id)
    
    # Complete in Todoist
    await todoist.complete_task(task_id)
    
    # Award in RPG system
    result = await zeus_rpg.complete_task(
        task=task,
        difficulty=estimate_difficulty(task)
    )
    
    # Build response
    response = f"✅ Completed: {task.content}\n"
    response += f"   +{result['xp_earned']} XP | +{result['gold_earned']} Gold\n"
    response += f"   🔥 {result['new_streak']}-day streak!"
    
    if result['level_up']:
        response += f"\n\n🎮 LEVEL UP! You're now Level {zeus_rpg.stats.level}!"
        await smart_home.celebration_mode()
    
    return response

@tool
async def show_stats():
    """Display Zeus's current RPG stats"""
    stats = zeus_rpg.stats
    
    return f"""
    ⚔️ **ZEUS - Level {stats.level}**
    
    ❤️ Health: {stats.health}/{stats.max_health}
    ✨ XP: {stats.xp}/{stats.xp_to_next_level}
    💰 Gold: {stats.gold}
    
    **Attributes:**
    💪 STR: {stats.strength}
    🧠 INT: {stats.intelligence}
    👁️ PER: {stats.perception}
    🛡️ CON: {stats.constitution}
    """

@tool
async def claim_reward(reward_name: str):
    """Spend gold on a custom reward"""
    rewards = {
        "coffee_break": {"cost": 10, "action": "15 minute break"},
        "game_time": {"cost": 50, "action": "30 min gaming"},
        "nice_dinner": {"cost": 100, "action": "Order nice food"},
        "day_off": {"cost": 500, "action": "Take tomorrow off"}
    }
    
    reward = rewards.get(reward_name)
    if not reward:
        return f"Unknown reward. Available: {list(rewards.keys())}"
    
    if await zeus_rpg.spend_gold(reward_name, reward['cost']):
        return f"🎁 Claimed: {reward['action']}! Enjoy!"
    else:
        return f"Not enough gold. Need {reward['cost']}, have {zeus_rpg.stats.gold}"
```

---

## Community Extensions

### Available for Habitica

| Extension | Description | Link |
|-----------|-------------|------|
| Habitica for Chrome | Browser extension | Chrome Web Store |
| Habitica for Obsidian | Sync with Obsidian | Community plugins |
| Habitica Chat Client | Desktop chat | GitHub |
| Habitica IFTTT | Automation triggers | IFTTT |
| Habitica Zapier | Workflow automation | Zapier |
| Stylus CSS | Custom themes (web) | Browser extension |

### Custom CSS Theming (Subscribers Only)

```css
/* Example: OneMind Purple Theme */
.header-wrapper {
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%) !important;
}

.task-wrapper {
    background: #0f0f23 !important;
    border-left: 3px solid #7c3aed !important;
}

.task-wrapper:hover {
    border-left-color: #06b6d4 !important;
}

/* Accent colors matching OneMind OS */
:root {
    --primary: #7c3aed;    /* Purple */
    --secondary: #06b6d4;  /* Cyan */
    --background: #0f0f23;
    --surface: #1a1a2e;
}
```

---

## Self-Hosting Option

### Requirements

- Node.js 20+
- MongoDB 3.6+
- Redis (optional, for sessions)

### Docker Setup

```yaml
# docker-compose.habitica.yml
version: '3.8'

services:
  habitica:
    build: 
      context: https://github.com/HabitRPG/habitica.git
      dockerfile: Dockerfile-Dev
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=development
      - ADMIN_EMAIL=zeus@onemind.com
      - BASE_URL=http://localhost:3000
    depends_on:
      - mongo
      - redis
    volumes:
      - ./config.json:/app/config.json

  mongo:
    image: mongo:4.4
    volumes:
      - habitica-mongo:/data/db

  redis:
    image: redis:alpine

volumes:
  habitica-mongo:
```

### Benefits of Self-Hosting

1. **Full customization** - Modify game mechanics
2. **Privacy** - Data stays on your server
3. **No subscription** - Premium features unlocked
4. **Custom themes** - Full CSS/JS control
5. **API modifications** - Add custom endpoints

### Drawbacks

1. **Maintenance burden** - Updates, security
2. **No community** - Solo experience
3. **No mobile apps** - Web only unless you fork apps
4. **Resource cost** - Server hosting

---

## Future Considerations

### When to Reconsider Habitica

1. **If Todoist becomes limiting** - Need more gamification
2. **If building multiplayer features** - Family/team accountability
3. **If custom RPG too complex** - Habitica already solved this
4. **If self-hosting desired** - Full control over mechanics

### Potential Hybrid Approach

```
                    ┌─────────────────────┐
                    │    TODOIST          │
                    │  (Task Management)  │
                    └──────────┬──────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────┐
│                      Legacy                          │
│            (Gamification Engine)                     │
│  ┌─────────────┐    ┌─────────────┐                 │
│  │ Zeus RPG    │    │ Habitica    │                 │
│  │ (Custom)    │ OR │ (Headless)  │                 │
│  └─────────────┘    └─────────────┘                 │
└──────────────────────────────────────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │     EXIST.IO        │
                    │   (Correlations)    │
                    └─────────────────────┘
```

---

## References

- [Habitica API Documentation](https://habitica.com/apidoc)
- [GitHub Repository](https://github.com/HabitRPG/habitica)
- [Guidance for Blacksmiths (Developers)](https://habitica.fandom.com/wiki/Guidance_for_Blacksmiths)
- [Extensions Wiki](https://habitica.fandom.com/wiki/Extensions,_Add-Ons,_and_Customizations)
- [API Usage Guidelines](https://github.com/HabitRPG/habitica/wiki/API-Usage-Guidelines)

---

## Document History

| Date | Author | Changes |
|------|--------|---------|
| 2026-01-21 | Zeus + Legacy | Initial research and documentation |
| - | - | Decision: Defer in favor of Todoist + Exist.io |

---

*This document is part of the OneMind OS integration planning series.*
*Linear Issue: OMOS-121 (Review Habitica Integration)*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
