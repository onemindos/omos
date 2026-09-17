# OneMind OS: Evolution Summary

> **From 77 Screens to Conversation-First Mobile**
> **Last Updated:** 2026-01-29

---

## The Big Picture

### Before → After

```
BEFORE (Desktop-First)               AFTER (Mobile-First)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
77 screens                    →      12 primary screens
6 drawer sections             →      4 bottom nav + capture FAB
Separate task/habit/goal      →      Universal Assets
Agents hidden in menus        →      Agents always visible
Manual everything             →      AI-powered automation
Desktop layouts               →      Mobile-optimized
Navigate to create            →      Capture anywhere
Module-based thinking         →      Conversation-first
Complex tactical UI           →      Simple, focused UI
```

---

## Linear Issue Summary

### Total Issues: 70

```
EXISTING (50 issues)
├── Context Engine (6)
├── Intelligence Scores (5)
├── Smart Chat (4)
├── Intelligent Scheduling (4)
├── Workflow Automation (4)
├── Inbox Intelligence (4)
├── Operating Rhythm (4)
├── Frontend - Original (4)
├── Life OS Completion (6)
└── Advanced Features (9)

NEW MOBILE-FIRST (20 issues)
├── Mobile Navigation (1)
│   └── OMOS-466: Mobile Navigation Redesign (8pts)
├── Universal Assets (2)
│   ├── OMOS-467: Universal Asset Schema (5pts)
│   └── OMOS-468: Asset Unified API (5pts)
├── Core Mobile Screens (4)
│   ├── OMOS-469: Universal Capture Component (5pts)
│   ├── OMOS-470: Home Dashboard Mobile (5pts)
│   ├── OMOS-471: Execute View (Today) (5pts)
│   └── OMOS-472: Command Chat Mobile (5pts)
├── Agent & Routing (3)
│   ├── OMOS-473: Agent Visibility Layer (3pts)
│   ├── OMOS-474: @Mention Agent Routing Mobile (3pts)
│   └── OMOS-475: /Command Quick Actions Mobile (3pts)
├── Secondary Screens (4)
│   ├── OMOS-476: Insights Mobile View (3pts)
│   ├── OMOS-477: Profile & Settings Mobile (3pts)
│   ├── OMOS-478: Voice Capture Integration (3pts)
│   └── OMOS-479: Photo Capture Processing (3pts)
├── Support Screens (5)
│   ├── OMOS-480: Smart Suggestions Widget (3pts)
│   ├── OMOS-481: Agent Status Dashboard (3pts)
│   ├── OMOS-482: Memory UI Component (3pts)
│   ├── OMOS-483: Knowledge Mobile UI (3pts)
│   └── OMOS-484: Mobile Approval UI (3pts)
└── Integration (1)
    └── OMOS-485: Integration & Polish Pass (5pts)
```

---

## How It All Connects

### Architecture Stack

```
┌─────────────────────────────────────────────────────────────┐
│                       MOBILE UI (Flutter)                    │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │  Home   │ │ Command │ │ Execute │ │ Profile │ + Capture │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
├─────────────────────────────────────────────────────────────┤
│                    AGNO AGENT LAYER                          │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ Legacy (Mother) → Routes to Specialists via @mentions   │ │
│  │ /Commands → Quick Actions                               │ │
│  │ Workflows → Automated Operations                        │ │
│  │ Approvals → Human-in-the-Loop                          │ │
│  └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    CONTEXT ENGINE                            │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ State + Assets + Scores + Calendar + History + Health   │ │
│  └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    DATA LAYER (Life OS)                      │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ Universal Assets (tasks, habits, goals, events, etc.)   │ │
│  │ Memory (user preferences, agent learnings)              │ │
│  │ Knowledge (documents, RAG)                              │ │
│  └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    INTEGRATIONS                              │
│  [Todoist] [Calendar] [Oura] [Home Assistant] [Finance]     │
└─────────────────────────────────────────────────────────────┘
```

### Agno Features → Mobile Mapping

| Agno Feature | Mobile Location | Component |
|--------------|-----------------|-----------|
| Agent Reasoning | Command | ThinkingLoader |
| Memory System | Knowledge | MemoryList |
| Knowledge Base | Knowledge | DocumentBrowser |
| Tools | Execute | ActionButtons |
| Skills | Command | AgentBadge |
| Teams | Home | ActiveTeamsCard |
| Workflows | Home | WorkflowStatusCard |
| Approvals | Modal | ApprovalSheet |
| Context | Everywhere | ContextIndicator |
| Scores | Home | ScoreCard |

---

## The New Mobile Experience

### Primary Flow

```
1. User opens app → Home Dashboard
   - See scores (PI: 87, LI: 82, WI: 83)
   - See today's focus
   - See active agents working

2. User taps Capture FAB → Capture Modal
   - Type, speak, or photograph
   - AI classifies automatically
   - One tap to save

3. User taps Command → Chat with Legacy
   - Natural language requests
   - @research for deep research
   - /tasks for quick list
   - Agent responses with identity

4. User taps Execute → Today's Actions
   - Focus blocks
   - Tasks to complete
   - Habits to check off
   - Schedule timeline

5. Agents work autonomously
   - Morning briefing at 6am
   - Task triage continuous
   - Evening review at 9pm
   - Approval requests push notifications
```

### Screen Count: 77 → 12

```
PRIMARY (Bottom Nav)
├── Home         ← Dashboard, scores, agents
├── Command      ← Chat with Legacy
├── Execute      ← Today's tasks/habits
├── Profile      ← Settings, config
└── [+ Capture]  ← Universal input (center FAB)

SECONDARY (From Primary)
├── Tasks        ← Full asset list
├── Calendar     ← Time-based view
├── Insights     ← Analytics deep dive
├── Goals        ← Three Paths progress
├── Agents       ← Agent management
├── Knowledge    ← Memory and docs
└── Settings     ← Deep configuration
```

---

## 8-Agent Swarm

### Week 1: Foundation

| Agent | Model | Focus | Key Deliverables |
|-------|-------|-------|------------------|
| Agent 1 | Opus 4.5 | Backend Architecture | Universal Asset schema, API |
| Agent 2 | Opus 4.5 | Context Engine | Aggregators, score calculation |
| Agent 3 | Opus 4.5 | Agent Routing | @mentions, /commands |
| Agent 4 | Sonnet | Mobile Navigation | New nav, routes |

### Week 2: UI Build

| Agent | Model | Focus | Key Deliverables |
|-------|-------|-------|------------------|
| Agent 5 | Sonnet | Home & Execute | Dashboard, today view |
| Agent 6 | Sonnet | Command & Capture | Chat, capture modal |
| Agent 7 | Sonnet | Insights & Profile | Analytics, settings |
| Agent 8 | Sonnet | Integration | Wire everything |

### Week 3: Polish

All agents on integration, testing, and polish.

---

## Realistic Assessment

### What's Working Now
- ✅ Full Agno framework with agents, tools, teams, workflows
- ✅ Life OS backend (90% complete)
- ✅ 77-screen Flutter app (working but complex)
- ✅ Memory and knowledge systems
- ✅ Context provider from Home Assistant
- ✅ Score calculation (basic)

### What Needs Building
- Universal Asset schema (unifies all data)
- Mobile-first navigation
- Capture modal with voice/photo
- Home dashboard with scores
- Execute view for today
- Command chat with routing
- Agent visibility layer
- Proper mobile UX polish

### Risk Factors
1. **Backend Changes**: Universal Asset requires data migration
2. **Navigation Refactor**: 77 → 12 screens is significant
3. **Voice/Photo**: External API dependencies
4. **Integration**: Many moving parts to connect

### Mitigation
- Agent 1 starts immediately (no dependencies)
- Agent 4 starts immediately (no dependencies)
- Backend changes have compatibility layer
- Week 3 dedicated to integration

---

## Success Metrics

### Week 1 Milestones
- [ ] Asset API returning data
- [ ] Context Engine building <500ms
- [ ] @mentions routing correctly
- [ ] New nav working on device

### Week 2 Milestones
- [ ] Home/Execute screens functional
- [ ] Chat/Capture working end-to-end
- [ ] All secondary screens built
- [ ] First integration pass complete

### Week 3 Milestones
- [ ] All screens connected
- [ ] No critical bugs
- [ ] Performance targets met
- [ ] Ready for production

### Final State
- 12 primary screens (down from 77)
- Conversation-first interface
- Visible autonomous agents
- Universal capture (text/voice/photo)
- Mobile-optimized throughout
- Full Agno feature utilization

---

## Key Documents

| Document | Purpose |
|----------|---------|
| [ONEMIND-MOBILE-FIRST-VISION.md](../../architecture/ONEMIND-MOBILE-FIRST-VISION.md) | Full vision document |
| [AGNO-MOBILE-INTEGRATION.md](../../architecture/AGNO-MOBILE-INTEGRATION.md) | Agno → Mobile mapping |
| [AGENT-BUILD-ASSIGNMENTS.md](AGENT-BUILD-ASSIGNMENTS.md) | 8-agent work assignments |
| [LINEAR_PROJECT_CONSCIOUSNESS_V2.md](LINEAR_PROJECT_CONSCIOUSNESS_V2.md) | Original 50 issues |

---

## Next Steps

1. **Deploy Agent 1 & 4** (no dependencies)
   - Agent 1: Start Universal Asset schema
   - Agent 4: Start Mobile Navigation

2. **Deploy Agent 2 & 3** (after Agent 1)
   - Agent 2: Context Engine
   - Agent 3: Routing

3. **Deploy Agents 5-7** (after Agent 4)
   - Agent 5: Home & Execute
   - Agent 6: Command & Capture
   - Agent 7: Insights & Profile

4. **Deploy Agent 8** (after all others)
   - Integration & Polish

---

**The transformation from 77 desktop screens to 12 mobile-first screens, powered by visible autonomous agents and universal capture, in 3 weeks.**
