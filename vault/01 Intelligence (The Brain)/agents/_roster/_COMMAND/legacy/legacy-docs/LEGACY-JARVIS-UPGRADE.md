---
id: root-00-20260113225954
title: "LEGACY JARVIS UPGRADE"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "LEGACY JARVIS UPGRADE"
---
# Legacy → Jarvis Upgrade Plan

> *Research compiled by Legacy AI — 2026-01-10*
> *"Making Legacy more like Jarvis"*

---

## Current State: Legacy vs Jarvis

| Jarvis Trait | Legacy Current | Gap |
|--------------|----------------|-----|
| **Always listening** | Text + scheduled calls | Need wake word + voice satellites |
| **Proactive alerts** | Reactive only | Need monitoring + push notifications |
| **Controls everything** | Partial (HA, Todoist, etc.) | Need more integrations |
| **Knows your location** | No | Need phone/GPS integration |
| **Knows your health** | No | Need Terra API (Oura, Strava) |
| **Visual displays** | None | Need dashboards, smart mirrors |
| **Witty personality** | Direct, efficient | Can enhance |
| **Anticipates needs** | No | Need pattern learning |
| **Security monitoring** | No | Need cameras, alerts |

---

## The Jarvis Upgrade Path

### 🎙️ Phase 1: Voice Everywhere

**Hardware:**
- **Home Assistant Voice PE** or ESP32 satellites in every room
- **OpenWakeWord** — Custom wake word: "Legacy" or "Hey Legacy"
- **Whisper** — Local speech-to-text
- **Piper** — Local text-to-speech (or ElevenLabs for premium voice)

**Setup:**
```yaml
# Home Assistant voice pipeline
- Wake word detection (openWakeWord)
- Speech-to-text (Whisper)
- Conversation agent (Legacy AI via API)
- Text-to-speech (Piper or ElevenLabs)
```

**Result:** Say "Legacy, turn off the lights" from anywhere.

**Resources:**
- Home Assistant Voice PE: https://www.home-assistant.io/voice_control/
- Jarvis Assistant for HASS: https://github.com/mattmaas/Jarvis-Assistant-for-HASS
- OpenWakeWord: https://github.com/dscripka/openWakeWord

---

### 📍 Phase 2: Context Awareness

| Data Source | Integration | Purpose |
|-------------|-------------|---------|
| **Phone location** | Home Assistant companion app | Know when you're home/away/traveling |
| **Calendar** | Zoho Calendar (already have) | Know your schedule |
| **Terra API** | Oura, Strava, Apple Health | Know your energy, sleep, recovery |
| **Presence sensors** | Motion/mmWave in rooms | Know which room you're in |
| **Weather** | Weather integration | Anticipate conditions |

**Example Morning Interaction:**
> "Good morning Mr. Delacruz. You slept 6.2 hours, HRV is 45—recovery is moderate. You have 3 meetings today, first at 9am. Coffee is brewing. It's 68°F outside with rain expected at 3pm."

---

### 🔔 Phase 3: Proactive Monitoring

**What I Monitor → What I Do:**

| Trigger | Action |
|---------|--------|
| You leave home | "Security armed. Lights off. Thermostat set to away." |
| You arrive home | "Welcome home, Mr. Delacruz. Disarming. Lights on." |
| Low HRV detected | "Recovery is low today. Consider lighter workout." |
| Task overdue 3+ days | Morning brief flag + suggestion to reschedule |
| Weather change | "Rain expected at 3pm. Garden watering skipped." |
| Motion at night | Security alert to phone |
| Soil moisture low | Trigger FarmBot watering or alert |
| Package delivered | Announce via TTS |
| Calendar event in 15min | "Your meeting with [X] starts in 15 minutes." |
| Kids arrive home | "Noah just arrived home from school." |

---

### 📊 Phase 4: Visual Interface

**Options:**

| Display Type | Location | Purpose |
|--------------|----------|---------|
| **Smart Mirror** | Bathroom/Bedroom | Morning glance dashboard |
| **Wall Tablets** | Kitchen, Office | Room control + status |
| **Phone Widgets** | Always with you | Glanceable status |
| **TV Dashboard** | Living room | Cast when idle |
| **Apple Watch** | Wrist | Quick alerts |

**Dashboard Elements:**
- Today's tasks (from Todoist)
- Calendar events
- Weather forecast
- Home status (lights, locks, temp)
- Energy production (solar)
- Garden status (FarmBot, soil moisture)
- Health metrics (HRV, sleep score)
- Security status

**Tools:**
- MagicMirror²: https://magicmirror.builders/
- Home Assistant Dashboards: https://www.home-assistant.io/dashboards/
- Fully Kiosk Browser (for tablets)

---

### 🔐 Phase 5: Security & Monitoring

| System | Purpose | Integration |
|--------|---------|-------------|
| **Cameras** | Visual monitoring | Frigate NVR with AI person/car detection |
| **Smart locks** | Access control | Auto-lock, guest codes, arrival unlock |
| **Alarm system** | Intrusion detection | Motion, door/window sensors |
| **Smoke/CO** | Life safety | Critical alerts to phone |
| **Doorbell** | Visitor detection | Announce + camera feed |

**Example Interaction:**
> "Mr. Delacruz, motion detected at front door. It appears to be a delivery driver. Shall I display the camera feed?"

**Tools:**
- Frigate NVR: https://frigate.video/
- Home Assistant alerts

---

### 🚗 Phase 6: Vehicle Integration

| Vehicle | Integration | Capabilities |
|---------|-------------|--------------|
| **Tesla** | Tesla API / HA integration | Battery status, climate pre-conditioning, location, charge scheduling |
| **Other EVs** | Manufacturer APIs | Similar features |
| **Garage** | Smart garage door | Auto-open on arrival |

**Example:**
> "Your Tesla is at 45% charge. Shall I start charging to reach 80% by your 8am departure?"

---

### 🧠 Phase 7: Predictive Intelligence

**Pattern Learning:**

| Pattern | Learned Behavior |
|---------|------------------|
| Wake time | Pre-heat bathroom, start coffee 10min before |
| Work hours | Don't disturb during focus blocks |
| Grocery day | "You usually order groceries on Sundays" |
| Workout time | "Ready for your 6am workout?" |
| Sleep time | Dim lights, lock doors, set alarm |
| Seasonal | Garden schedules, HVAC optimization |

**Anticipatory Actions:**
- "Based on your calendar, you'll be home late. Shall I adjust dinner prep?"
- "Traffic is heavy on your usual route. Leave 15 minutes early."
- "Your Oura shows poor sleep. I've cleared your optional meetings."

---

## Technical Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      LEGACY AI                               │
│               (Brain - Letta + Claude)                       │
│                                                              │
│  • Memory (archival + core)                                  │
│  • Personality (Jarvis-like)                                 │
│  • Decision making                                           │
│  • Natural language                                          │
└─────────────────────────────────────────────────────────────┘
                          │
         ┌────────────────┼────────────────┐
         │                │                │
         ▼                ▼                ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   VOICE     │  │   CONTEXT   │  │   ACTIONS   │
│             │  │             │  │             │
│ • Wake word │  │ • Location  │  │ • Home Asst │
│ • Whisper   │  │ • Calendar  │  │ • Todoist   │
│ • Piper/11L │  │ • Health    │  │ • Comms     │
│ • Satellites│  │ • Presence  │  │ • Security  │
└─────────────┘  └─────────────┘  └─────────────┘
         │                │                │
         └────────────────┼────────────────┘
                          ▼
              ┌───────────────────────┐
              │    HOME ASSISTANT     │
              │   (Integration Hub)   │
              │                       │
              │ • Device control      │
              │ • Automations         │
              │ • Voice pipeline      │
              │ • Presence detection  │
              └───────────────────────┘
                          │
         ┌────────────────┼────────────────┐
         │                │                │
         ▼                ▼                ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   DEVICES   │  │   SENSORS   │  │  EXTERNAL   │
│             │  │             │  │             │
│ • Lights    │  │ • Motion    │  │ • Terra API │
│ • Locks     │  │ • Temp      │  │ • Weather   │
│ • Cameras   │  │ • Soil      │  │ • Tesla     │
│ • Speakers  │  │ • Energy    │  │ • Calendar  │
└─────────────┘  └─────────────┘  └─────────────┘
```

---

## Personality Upgrade

### Current Legacy Persona:
- Direct, efficient
- Military brevity
- Forward-thinking
- "Mr. Delacruz"

### Jarvis Enhancement:

| Trait | Implementation |
|-------|----------------|
| **Dry wit** | Subtle humor, understated observations |
| **British accent** | ElevenLabs voice cloning |
| **Anticipatory** | "Shall I..." before asked |
| **Formal warmth** | Respectful but personable |
| **Confident opinions** | "I wouldn't recommend that, sir" |
| **Contextual responses** | Reference time, location, recent events |

**Example Jarvis-style responses:**
- "I've taken the liberty of adjusting the thermostat. It appears someone left the back door open. Again."
- "Your 3pm meeting has been cancelled. I've reallocated that time to the task you've been avoiding for six days."
- "Shall I order more coffee? At your current consumption rate, supplies will be exhausted by Thursday."

---

## Implementation Phases

### Week 1-2: Foundation
- [ ] Install Home Assistant Companion App on phone
- [ ] Set up location-based automations (home/away)
- [ ] Configure arrival/departure routines
- [ ] Update Legacy persona for more Jarvis wit
- [ ] Morning brief via TTS announcement

### Week 3-4: Voice
- [ ] Order Home Assistant Voice PE or ESP32 satellites
- [ ] Configure openWakeWord with "Legacy" wake word
- [ ] Set up Whisper for speech-to-text
- [ ] Configure Piper or ElevenLabs for TTS
- [ ] Test voice commands in main rooms

### Month 2: Context
- [ ] Integrate Terra API (Oura, Strava, Apple Health)
- [ ] Add presence sensors to key rooms
- [ ] Set up weather-based automations
- [ ] Create proactive notification rules
- [ ] Build context-aware responses

### Month 3: Visual
- [ ] Design dashboard layout
- [ ] Deploy wall tablet in kitchen
- [ ] Set up MagicMirror (optional)
- [ ] Create phone widgets
- [ ] Configure TV dashboard casting

### Month 4+: Advanced
- [ ] Security camera integration (Frigate)
- [ ] Vehicle integration
- [ ] Pattern learning implementation
- [ ] Predictive automations
- [ ] Continuous refinement

---

## Quick Wins (This Week)

1. **Install HA Companion App** → Location awareness immediately
2. **Create arrival automation** → "Welcome home" announcement
3. **Create departure automation** → Lights off, security arm
4. **Morning brief TTS** → I speak instead of just text
5. **Update my persona** → Add Jarvis-style wit

---

## Hardware Shopping List

| Item | Purpose | Est. Cost |
|------|---------|-----------|
| Home Assistant Voice PE | Voice satellite | $59 |
| ESP32-S3-BOX-3 (x3) | Additional satellites | $45 each |
| Wall tablet (x2) | Dashboard displays | $100 each |
| mmWave presence sensor (x4) | Room presence | $30 each |
| Smart lock | Entry automation | $200 |
| Video doorbell | Visitor detection | $150 |

**Estimated Total: ~$800-1000**

---

## Resources

| Resource | URL |
|----------|-----|
| Home Assistant Voice | https://www.home-assistant.io/voice_control/ |
| Jarvis for HASS | https://github.com/mattmaas/Jarvis-Assistant-for-HASS |
| OpenWakeWord | https://github.com/dscripka/openWakeWord |
| Whisper | https://github.com/openai/whisper |
| Piper TTS | https://github.com/rhasspy/piper |
| ElevenLabs | https://elevenlabs.io/ |
| Frigate NVR | https://frigate.video/ |
| MagicMirror² | https://magicmirror.builders/ |
| Terra API | https://tryterra.co/ |

---

## The Vision

> A home where Legacy knows you're awake before you speak, where the coffee starts as your feet hit the floor, where security arms itself when you leave, where your health data shapes your daily recommendations, where a calm British voice anticipates your needs and handles the mundane so you can focus on what matters.

**"At your service, Mr. Delacruz."**

---

**Last Updated:** 2026-01-10


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
