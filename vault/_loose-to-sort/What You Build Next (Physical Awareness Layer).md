# What You Build Next (Physical Awareness Layer)

## NATS Subjects — The Complete Namespace

```
┌─────────────────────────────────────────────────────────────────────────┐
│  NATS SUBJECTS — The Complete Namespace                                  │
│                                                                          │
│  ═══ AGENT PROTOCOL (Synadia Agent Protocol v0.3) ═══                    │
│  agents.prompt.{agent}.{owner}.{name}    → prompt an agent               │
│  agents.hb.{agent}.{owner}.{name}        → heartbeat (auto)             │
│  agents.status.{agent}.{owner}.{name}    → on-demand status             │
│                                                                          │
│  ═══ FLEET (Edge Devices — Physical World) ═══                           │
│  fleet.telemetry.{device_id}             → periodic sensor data          │
│  fleet.cmd.{device_id}                   → commands TO devices           │
│  fleet.status.{device_id}                → request/reply health          │
│  fleet.alert.{device_id}                 → urgent events (fence, smoke)  │
│  fleet.position.{device_id}              → GPS updates                   │
│  fleet.video.{device_id}                 → video frame stream            │
│                                                                          │
│  ═══ HOME (IoT via HomeAssistant bridge) ═══                             │
│  home.state.{entity_id}                  → entity state changes          │
│  home.cmd.{entity_id}                    → commands to devices           │
│  home.event.{trigger}                    → automations fired             │
│                                                                          │
│  ═══ CONTEXT (Zeus Awareness — KV Watches) ═══                           │
│  zeus.location                           → where you are physically      │
│  zeus.focus_mode                         → deep work / available / sleep │
│  zeus.activity                           → working / driving / walking   │
│  zeus.device                             → which device active           │
│                                                                          │
│  ═══ AGENT INTERNAL (JetStream Persistence) ═══                          │
│  agent.{agent}.tool_call                 → tool usage log                │
│  agent.{agent}.decision                  → reasoning trace               │
│  agent.{agent}.error                     → failures                      │
│  task.created.{id}                       → new task                      │
│  task.completed.{id}                     → task done                     │
│  comms.{platform}.{direction}            → message log                   │
│  system.health.{node}                    → node health                   │
│  system.backup.{timestamp}               → backup events                 │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Event Flow Diagrams

```
PHYSICAL WORLD                    NATS FABRIC                     LEGACY (Brain)
─────────────────                 ───────────                     ──────────────

Drone-1 GPS update ──────────→ fleet.position.drone-1 ──────→ nats-fleet plugin
                                                                  (buffers in memory)

Fence sensor trip ───────────→ fleet.alert.fence-north ─────→ nats-fleet plugin
                                                                  → Legacy REACTS:
                                                                    "Alert: north fence
                                                                     breach detected.
                                                                     Launching drone."
                                                                  → fleet.cmd.drone-1
                                                                    {action: "investigate",
                                                                     waypoint: [33.45,-112.07]}

Smoke detector ──────────────→ home.event.smoke_alarm ──────→ HA bridge
                                                                  → Legacy REACTS:
                                                                    "🚨 SMOKE: kitchen.
                                                                     Killed HVAC. Alerting."
                                                                  → Telegram to Zeus
                                                                  → fleet.cmd.irrigation
                                                                    {action: "zone3_on"}

Frame glasses (you look) ────→ frame-bridge plugin ─────────→ Legacy processes:
                                                                  "I see broken post.
                                                                   GPS: 33.451,-112.069
                                                                   Adding repair task."
                                                                  → task.created.fence-repair
                                                                  → frame display: "✅ Logged"

Zeus drives home ────────────→ zeus.location = "driving" ───→ KV watch fires
                                 zeus.activity = "commute"       → Legacy adjusts:
                                                                    "Zeus driving. Queue
                                                                     non-urgent until home."
```
