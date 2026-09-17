# OneMind OS — The Vision

> Master vision document extracted from early planning docs (Feb 2026).
> Stack-agnostic — these are the *ideas*, not the tech. Current stack: OpenClaw + NATS.

---

## The Core Equation

**HUMAN (Zeus) + AI (LEGACY) = ONE MIND**

One unified intelligence that manages the full spectrum of life:
digital (business, research, communication) + physical (home, garden, robots, family).

Legacy sits **above** all systems — orchestrating, not executing. Every tool, robot, and automation is a skill Legacy can call.

---

## The 5-Domain Architecture

```
UI (00-24)  — The Brain        → Legacy + tech infrastructure
HP (25-49)  — Self Path        → Health, fitness, cognition, finance
LE (50-74)  — Legacy Path      → Family, home, estate, generational wealth
GE (75-99)  — Wealth Path      → Business operations, ventures, revenue
OS          — Physical Layer   → Robots, sensors, cameras, garden, home
```

---

## The Heartbeat Loop (Core Intelligence Pattern)

Legacy runs a continuous **SENSE → FUSE → CORTEX → ACTION → REMEMBER** cycle:

```
1. SENSE    — Gather all inputs (calendar, email, sensors, cameras, NATS events)
2. FUSE     — Build unified context (digital state + physical state + goals + memory)
3. CORTEX   — Route to right agent or human (digital agent / physical agent / HITL)
4. ACTION   — Execute (API calls, robot commands, IoT, voice response)
5. REMEMBER — Store learnings back to shared memory, evolve
```

Proactive cadence: check every 30 minutes minimum. React instantly to NATS events.

---

## Intelligent Event Routing

Every event has a routing decision:
- **Digital domain** → Legacy handles directly (email, calendar, research, code)
- **Physical domain** → Dispatch to specialist (robot command, Home Assistant, FarmBot)
- **Cross-domain** → Multi-step coordination (e.g., "prepare house for arrival" = route + lock + climate + grocery order)
- **High stakes** → Human-in-the-loop (financial decisions, medical, legal, deletions)

Routing keywords for physical: `robot, garden, water, plant, lights, door, camera, temperature, move, navigate, pick up`

---

## Agent Roster (Future Expansion)

| Agent | Domain | Specialty |
|-------|--------|-----------|
| **Legacy (UI)** | All | Central orchestrator, voice, communication |
| **HP Agent** | 25-49 | Health coach, fitness, nutrition, sleep, cognition |
| **LE Agent** | 50-74 | Family steward, homeschool, household, legacy docs |
| **GE Agent** | 75-99 | Business advisor, revenue, ventures |
| **Eden** | Physical | FarmBot garden automation |
| **HomeOS** | Physical | Home Assistant orchestration |
| **Scout** | Physical | Mobile patrol/inspection robot |

All agents share unified context via NATS and a shared memory store.

---

## HP Agent Vision (Self Path)

**Goal:** Full personal optimization loop — proactive, not reactive.

**Capabilities:**
- Morning biometric analysis → adaptive daily plan (sleep score → adjust workout intensity, push meetings if recovery low)
- Adaptive workout generation (considers recovery data, recent training load, goals)
- Nutrition optimizer (macros, meal suggestions, shopping list generation)
- Sleep analyzer (score quality, suggest protocol adjustments — caffeine timing, nap windows)
- Learning curriculum planner (study schedules, skill tracking, knowledge gaps)
- Habit tracker (streaks, reminders, pattern analysis)
- Cognitive performance tracking (identify decline patterns → intervene)

**Example:** "You had 6.2 hrs sleep, 40% deep — below baseline. Workout intensity reduced 20%, first meeting pushed to 9:30 AM."

---

## LE Agent Vision (Family + Legacy)

**Goal:** Coordinate a large, multi-household family operation autonomously.

**Scale:** 10-12 children, 3 separate households, homeschool curriculum, estate planning.

**Capabilities:**
- **Multi-household calendar** — sync, conflict detection, optimal scheduling across all 3 households
- **Per-child homeschool tracking** — grade level, curriculum, daily schedule, completed assignments, progress milestones, strengths/challenges, parent observations
- **Values checker** — before major decisions, check alignment with documented family values framework. "This property purchase conflicts with your 'minimize debt' principle."
- **Legacy documentation** — capture stories, lessons, wisdom from Zeus and ancestors. Searchable. "What did my grandfather teach me about business?"
- **Household operations** — chore rotation, maintenance schedules (HVAC filters, appliances), shopping lists from inventory, meal planning
- **Relationship maintenance** — check-in reminders, birthday alerts, anniversary preparation

**Daily briefing schedule:**
- 7 AM: Today's schedule, weather, urgent tasks, birthdays
- 12 PM: Progress updates, schedule changes
- 6 PM: What's left today, tomorrow preview, family time suggestions
- Sunday PM: Week ahead, goals check-in, family meeting agenda

---

## Home Automation Vision

**Platform:** Home Assistant (local-first, 2,000+ integrations)

**Domains:**
- **Climate:** Smart thermostats, HVAC zones, temp/humidity sensors, sleep/wake/away modes
- **Lighting:** Circadian rhythm automation (warm at night, bright morning), motion-activated, scene modes
- **Security:** Smart locks, door/window sensors, motion sensors, glass break detection
- **Access control:** Garage/gate controllers, facial recognition auto-unlock for family
- **Energy:** Solar system monitoring, battery storage tracking, load balancing, off-grid capability
- **Appliances:** Washer/dryer status, coffee maker start, irrigation control

**Voice automation chains:**
- `"work mode"` → office lights bright, distractions minimized, DND on devices
- `"family movie time"` → living room dims, TV on, snacks reminder
- `"emergency lockdown"` → lock all, all exterior lights max, alert contacts
- `"sleep mode"` → lights off, thermostat adjusts, all doors locked, security armed, random light patterns (presence simulation)

**Legacy orchestrates HA** — Legacy gives intent, HomeOS agent executes at device level.

---

## Camera & Security Vision

**Platform:** Frigate NVR (local AI, no cloud)

**Infrastructure:**
- 8-12 exterior cameras (all entrances, perimeter, driveways)
- 4-6 interior cameras (common areas only — never bedrooms/bathrooms)
- Google Coral TPU for local AI processing (50-200ms detection)
- 4-8 TB NVR storage (30+ days 24/7 recording)

**AI capabilities:**
- Person, vehicle, animal, package detection
- **Facial recognition** — trained on family photos, identifies each child/adult across all cameras
- Loitering detection, zone intrusion alerts
- Privacy mode (guest mode disables interior cameras; all data stays local)

**Use cases:**
- `"Where is Noah?"` → Frigate face search across cams → "Living room, 2 minutes ago"
- `"Has my wife left?"` → Vehicle detection → "Her car left at 8:47 AM"
- Autonomous: Package delivered → notify. Unknown person at 2 AM → emergency alert + exterior lights
- `"Show me the backyard right now"` → Live stream on device

---

## Robotics Vision (Phased)

### Phase 1: Mobile Robots (Year 1-2)
- **Patrol robot** — autonomous nightly security patrols, anomaly detection
- **Delivery robot** — fetch/deliver items between rooms or buildings
- **Inspection robot** — check livestock, plants, solar panels, fences
- **Telepresence robot** — video call on wheels for remote family interaction

**Platform:** ROS 2 (Nav2 navigation, SLAM mapping, collision avoidance)
**Hardware:** TurtleBot 4 (learning, ~$1,500), Clearpath Jackal (outdoor, ~$20K), or custom ESP32 (~$500)

### Phase 2: Manipulator Robots (Year 2-3)
- Pick up packages, deliver to rooms
- Organize items (toys, tools, inventory)
- Feed livestock (automated food dispensing)
- Harvest greenhouse produce
- **Hardware:** UR5e (collaborative arm, ~$30K), Stretch Robot (~$25K), or custom gripper on mobile base

### Phase 3: Humanoid Robots (Year 3-5)
- General-purpose household assistance
- Open doors, navigate stairs, carry up to 50 lbs
- Household chores (laundry, dishes)
- Child supervision (remote monitoring + physical intervention if authorized)
- Assist elderly family members
- Outdoor work (gardening, maintenance, repairs)
- **Platforms to watch:** Figure 02, Tesla Optimus, 1X NEO (~$30K target), Unitree H1

**Safety protocols:**
- Emergency stop via voice + physical button
- Force-limiting motors (can't harm)
- Child detection → automatic slow-down
- Geofenced operation zones
- Manual override always available

### Phase 4: Agricultural Automation (Year 3+)
**Platform:** FarmBot + Custom

- Automated planting by season/calendar
- Precision watering (soil moisture sensors)
- Weed detection and organic treatment
- Pest monitoring
- Harvest assistance
- Greenhouse climate control
- Livestock monitoring (health, feed levels, water)

---

## NATS Event Taxonomy (Physical + Digital)

```
events.digital.*       → Calendar, email, webhooks, task updates
events.physical.*      → Sensors, robots, IoT, HA state changes
cortex.decisions.*     → Agent plans and routing decisions
action.digital.*       → API calls, tool executions
action.physical.*      → Robot commands, actuators, HA calls
human.approvals.*      → HITL approval workflows
comms.text.human.zeus  → Zeus's inbound messages
comms.text.agent.legacy → Legacy's outbound messages
nexus.system.health    → System health pings
```

---

## Offline / Decentralized Mode

The system must function when internet is down:
- All core functionality runs on local hardware (no cloud for critical paths)
- NATS leaf nodes enable local-first with cloud sync when available
- Legacy can operate fully offline (local LLM fallback or cached responses)
- Robots navigate and execute pre-programmed tasks without internet
- Home automation operates independently via Home Assistant

---

## Future Canvas Ideas

- **NATS Pulse Visualization** — OpenClaw canvas screen that renders live NATS message flow as animated particle effects, showing agent activity and system health as a living visual (think: each message = a spark traveling between nodes)
- **Tactical Map** — Live OpenStreetMap canvas overlaying all assets (robots, family members by location, deliveries) with real-time position updates
- **System Heartbeat Dashboard** — Cortex screen showing all 5 domain agents, their current state, active tasks, and event queue

---

## The North Star

> A sovereign, solar-punk AI operating system for the Delacruz family —
> running 24/7 on owned hardware, managing digital life and physical home,
> raising the next generation with AI-augmented education, keeping the family
> coordinated across 3 households, building generational wealth through autonomous
> business operations, and embodying a future where technology serves human
> flourishing rather than extracting from it.
