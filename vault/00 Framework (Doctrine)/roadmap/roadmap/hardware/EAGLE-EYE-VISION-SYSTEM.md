# Eagle Eye - Legacy Vision & Awareness System

**Last Updated:** January 15, 2026  
**Status:** Planning  
**Classification:** Level 3 - Autonomous Operations Center  
**Codename:** EAGLE EYE

---

## Executive Summary

Eagle Eye is Legacy's comprehensive awareness system - the ability to SEE, TRACK, and UNDERSTAND everything in Zeus's digital and physical world.

**Core Capabilities:**
- 👁️ **Screen Vision** - See and control any screen
- 📍 **Location Tracking** - Know where Zeus is at all times
- 📷 **Camera Vision** - Eyes in every room and entrance
- 🧠 **Situational Awareness** - Understand context and act proactively

---

## Section 1-2: Concept & Screen Vision

- Eagle Eye = AI that sees, tracks, understands
- Screen Vision: Screenshot → Claude Vision → Click/Type → Repeat

## Section 3: Location Tracking

- Phone GPS → Home Assistant → Zones → Legacy queries
- Triggers: enter/leave zones, proximity alerts, stationary checks

## Section 4: Camera Vision

- Yi cameras → Home Assistant → Legacy snapshots → Claude Vision analysis
- Package detection, security alerts, kids monitoring

---

## Section 5: Situational Awareness

### Overview

Situational Awareness is the BRAIN layer - combining all vision inputs to understand CONTEXT and act proactively.

```
┌─────────────────────────────────────────────────────────────────┐
│                    SITUATIONAL AWARENESS                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   INPUTS                           CONTEXT ENGINE                │
│   ┌─────────────┐                 ┌─────────────────────────┐   │
│   │ Screen      │────────────────►│                         │   │
│   │ Vision      │                 │   "Zeus is working in   │   │
│   └─────────────┘                 │    Monarch Money on     │   │
│                                   │    his Mac. He's at     │   │
│   ┌─────────────┐                 │    home. Kids are in   │   │
│   │ Location    │────────────────►│    the living room.    │   │
│   │ Tracking    │                 │    Package at door.     │   │
│   └─────────────┘                 │    No security alerts." │   │
│                                   │                         │   │
│   ┌─────────────┐                 └───────────┬─────────────┘   │
│   │ Camera      │────────────────►            │                 │
│   │ Vision      │                             │                 │
│   └─────────────┘                             ▼                 │
│                                   ┌─────────────────────────┐   │
│   ┌─────────────┐                 │    PROACTIVE ACTIONS    │   │
│   │ Calendar    │────────────────►│                         │   │
│   │ Tasks       │                 │ • Don't disturb Zeus    │   │
│   └─────────────┘                 │ • Monitor kids quietly  │   │
│                                   │ • Queue package reminder│   │
│   ┌─────────────┐                 │ • Prep for next meeting │   │
│   │ Time of Day │────────────────►│                         │   │
│   │ Patterns    │                 └─────────────────────────┘   │
│   └─────────────┘                                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Context Building

Legacy constantly builds a picture of "what's happening right now":

```python
class SituationalAwareness:
    def __init__(self):
        self.context = {}
    
    async def update_context(self):
        # Location
        self.context["zeus_location"] = await get_location("zeus_iphone")
        
        # Time context
        self.context["time"] = datetime.now()
        self.context["is_work_hours"] = 9 <= self.context["time"].hour <= 17
        self.context["is_night"] = self.context["time"].hour >= 22 or self.context["time"].hour <= 6
        
        # Home status
        self.context["family_home"] = await check_who_is_home()
        self.context["security_armed"] = await get_security_status()
        
        # Active tasks
        self.context["todays_tasks"] = await get_todoist_tasks("today")
        self.context["next_meeting"] = await get_next_calendar_event()
        
        # Camera status (only if motion detected)
        if await motion_detected_recently():
            self.context["last_motion"] = await get_last_motion_event()
        
        # Screen activity (if HITL active)
        if self.hitl_active:
            self.context["current_screen"] = await get_current_activity()
        
        return self.context
```

### Proactive Actions Based on Context

| Context | Legacy Action |
|---------|---------------|
| **Zeus working + kids home** | Monitor kids quietly, don't interrupt unless urgent |
| **Zeus left office + meeting in 30 min** | "Reminder: You have a meeting at 3pm" |
| **Zeus at gym + package delivered** | Queue notification for when he returns |
| **Night time + motion outside** | Heightened security monitoring |
| **Zeus stationary 4 hours + no lunch** | "Mr. Delacruz, reminder to eat" |
| **Kids room empty + bedtime passed** | Check other cameras, gentle reminder |
| **Zeus driving + incoming call** | Let call through, prep summary |

### Daily Context Patterns

```
┌─────────────────────────────────────────────────────────────────┐
│                      DAILY RHYTHM                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  6:00 AM - MORNING MODE                                         │
│  ├── Check Zeus location: Home                                  │
│  ├── Prep morning brief                                         │
│  ├── Check weather, traffic                                     │
│  └── Queue first announcements                                  │
│                                                                  │
│  8:00 AM - WORK MODE                                            │
│  ├── Zeus typically leaves for office                           │
│  ├── Arm security when family leaves                            │
│  ├── Reduce home monitoring frequency                           │
│  └── Focus on calendar, tasks, business                         │
│                                                                  │
│  12:00 PM - MIDDAY CHECK                                        │
│  ├── Check Zeus location                                        │
│  ├── Afternoon task reminders                                   │
│  └── School pickup logistics (if applicable)                    │
│                                                                  │
│  5:00 PM - EVENING TRANSITION                                   │
│  ├── Zeus typically leaving office                              │
│  ├── Calculate ETA, prep home                                   │
│  ├── Evening brief preparation                                  │
│  └── Family coordination mode                                   │
│                                                                  │
│  8:00 PM - FAMILY TIME                                          │
│  ├── Reduce work notifications                                  │
│  ├── Kids bedtime reminders                                     │
│  └── Light home automation                                      │
│                                                                  │
│  10:00 PM - NIGHT MODE                                          │
│  ├── Security patrol every 15 min                               │
│  ├── All exterior cameras active                                │
│  ├── Motion = immediate alert                                   │
│  └── Tomorrow prep (calendar, tasks)                            │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Section 6: Implementation Checklist

### Phase 1: Screen Vision
- [ ] Deploy computer-use Docker container
- [ ] Test screenshot capture
- [ ] Test Claude Vision analysis
- [ ] Test click/type automation
- [ ] Create first UI automation task

### Phase 2: Location Tracking
- [ ] Install Home Assistant Companion on iPhone
- [ ] Configure zones (Home, Office, Gym, School, etc.)
- [ ] Test zone enter/leave triggers
- [ ] Create location-based automations
- [ ] Test proximity alerts

### Phase 3: Camera Vision
- [ ] Integrate Yi cameras with Home Assistant
- [ ] Test snapshot capture via HA API
- [ ] Test Claude Vision analysis on snapshots
- [ ] Configure motion detection triggers
- [ ] Create camera-triggered automations

### Phase 4: Situational Awareness
- [ ] Build context engine
- [ ] Integrate all data sources
- [ ] Create daily rhythm patterns
- [ ] Test proactive actions
- [ ] Fine-tune timing and triggers

---

## Quick Reference

### Legacy Commands

| Command | Action |
|---------|--------|
| "Legacy, where am I?" | Query current location |
| "Legacy, who's at the door?" | Snapshot + analyze front door |
| "Legacy, check on the kids" | Snapshot + analyze kids room |
| "Legacy, what am I looking at?" | Screenshot + analyze current screen |
| "Legacy, security status" | Report all cameras, motion, alerts |
| "Legacy, what's the context?" | Full situational report |

### Automation Triggers

```yaml
# Home Assistant automations for Eagle Eye

automation:
  - alias: "Front door motion - notify Legacy"
    trigger:
      - platform: state
        entity_id: binary_sensor.yi_front_door_motion
        to: "on"
    action:
      - service: rest_command.notify_legacy
        data:
          event: "motion_front_door"
          
  - alias: "Zeus arrives home"
    trigger:
      - platform: zone
        entity_id: device_tracker.zeus_iphone
        zone: zone.home
        event: enter
    action:
      - service: rest_command.notify_legacy
        data:
          event: "zeus_home"
          
  - alias: "Zeus leaves home"
    trigger:
      - platform: zone
        entity_id: device_tracker.zeus_iphone
        zone: zone.home
        event: leave
    action:
      - service: rest_command.notify_legacy
        data:
          event: "zeus_left"
```

---

## The Vision

**Eagle Eye transforms Legacy from assistant to guardian.**

- Always watching (cameras, screens)
- Always knowing (location, schedule)
- Always understanding (context, patterns)
- Always acting (proactive, protective)

**This is Level 3: Autonomous Operations Center.**

**This is your personal Eagle Eye.**

---

**Document Complete**

*Related: Legacy-OPERATIONS-CENTER.md (infrastructure)*
