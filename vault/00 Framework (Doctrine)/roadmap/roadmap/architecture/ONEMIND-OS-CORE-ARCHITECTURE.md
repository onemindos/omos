# OneMind OS - Core Architecture

**The True Heart of the System**

---

## The Hierarchy (Corrected Understanding)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│                            ONEMIND OS                                    │
│                    (The Living Operating System)                         │
│                                                                          │
│   ════════════════════════════════════════════════════════════════════  │
│                                                                          │
│                     ★ THE HEART: EVENT STREAM BUS ★                     │
│                                                                          │
│       The living river of ALL experiences, actions, and context         │
│                                                                          │
│       ─────────────────────────────────────────────────────────────     │
│       │                                                           │     │
│       │   YOUR actions ──────────────────────────────────────▶   │     │
│       │   OTHER human actions ───────────────────────────────▶   │     │
│       │   System events ─────────────────────────────────────▶   │     │
│       │   IoT/sensors ───────────────────────────────────────▶   │     │
│       │   Digital activity ──────────────────────────────────▶   │     │
│       │                                                           │     │
│       │              ══════════════════════════                   │     │
│       │              ║   CONTINUOUS STREAM   ║                   │     │
│       │              ║   (Redis Streams)     ║                   │     │
│       │              ══════════════════════════                   │     │
│       │                         │                                 │     │
│       │                         ▼                                 │     │
│       │              FULL CONTEXT OF ALL THINGS                   │     │
│       │              HAPPENING IN REAL-TIME                       │     │
│       │                                                           │     │
│       ─────────────────────────────────────────────────────────────     │
│                                 │                                        │
│                                 ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                         Legacy                                   │   │
│   │                    (The AI Agent)                                │   │
│   │                                                                  │   │
│   │   One piece of OneMind OS that:                                 │   │
│   │   • WATCHES the event stream (awareness)                        │   │
│   │   • TALKS to you (voice interface)                              │   │
│   │   • ACTS on your behalf (tool execution)                        │   │
│   │   • DELEGATES to worker agents                                  │   │
│   │                                                                  │   │
│   │   Legacy is aware of ALL experiences around you                 │   │
│   │   because she's connected to the living stream                  │   │
│   │                                                                  │   │
│   └─────────────────────────────────────────────────────────────────┘   │
│                                 │                                        │
│                                 ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                         LETTA                                    │   │
│   │               (Intelligence + Memory Layer)                      │   │
│   │                                                                  │   │
│   │   Just the brain that Legacy uses:                              │   │
│   │   • Memory storage (what Legacy remembers)                      │   │
│   │   • Reasoning (how Legacy thinks)                               │   │
│   │   • Tool calling (how Legacy executes)                          │   │
│   │                                                                  │   │
│   │   Letta is NOT the heart - it's the intelligence                │   │
│   │   that processes what comes from the stream                     │   │
│   │                                                                  │   │
│   └─────────────────────────────────────────────────────────────────┘   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## The Key Insight

```
WRONG understanding:
  "Legacy-AI is the main thing, event bus is a feature"

RIGHT understanding:
  "The EVENT STREAM is the heart, Legacy is the agent that uses it"
```

---

## What This Means

### The Event Stream Bus IS OneMind OS

The stream is the **living awareness** - it's what makes the system "alive":

```
Without the stream:
  Legacy is just a chatbot you call when you need something
  ❌ Not aware
  ❌ Not watching
  ❌ Not alive

With the stream:
  Legacy knows everything happening around you in real-time
  ✅ Aware of your actions (task created, email sent)
  ✅ Aware of others' actions (family member came home)
  ✅ Aware of systems (server alert, calendar reminder)
  ✅ ALIVE - always watching, always ready
```

### Legacy is ONE Component

Legacy is the **agent** that:
- Connects to the stream
- Processes the context
- Talks to you
- Takes action

But Legacy is not the whole OS. The OS is the entire living system.

### Letta is Just Intelligence

Letta provides:
- Memory (so Legacy remembers)
- Reasoning (so Legacy thinks)
- Execution (so Legacy acts)

But Letta doesn't make Legacy "alive" - the **stream** does.

---

## Repository Naming (Recommended Change)

Current:
```
legacy-ai/     ← Name focuses on the agent
```

Should be:
```
onemind-os/    ← Name reflects the whole system
├── stream/    ← THE HEART (event bus)
├── agents/
│   └── legacy/  ← The main agent (one piece)
├── voice/
└── tools/
```

---

## The Analogy

```
HUMAN BODY:

  Bloodstream = Event Stream Bus
  ├── Carries oxygen, nutrients, signals to EVERYWHERE
  ├── Always flowing
  └── The body dies without it

  Brain = Legacy (the agent)
  ├── Processes information from the bloodstream
  ├── Makes decisions
  └── Controls actions

  Neural Networks = Letta (intelligence)
  ├── How the brain thinks
  ├── Memory storage
  └── Pattern recognition

THE HEART IS THE BLOODSTREAM (EVENT STREAM)
NOT THE BRAIN (Legacy)
```

---

## Summary

| Component | What It Is | Role |
|-----------|------------|------|
| **Event Stream Bus** | THE HEART | Living awareness of all experiences |
| **Legacy** | The Agent | One piece - watches stream, talks, acts |
| **Letta** | Intelligence | Memory + reasoning Legacy uses |
| **OneMind OS** | The Whole System | Stream + Agent + Intelligence + Tools |

---

**The stream is what makes Legacy alive.**
**Without the stream, Legacy is just another chatbot.**

---

*Created: January 16, 2026*
