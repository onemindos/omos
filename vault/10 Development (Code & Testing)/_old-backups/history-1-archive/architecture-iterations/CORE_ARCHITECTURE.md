# OneMind OS Core

The brain of OneMind OS - everything the AI system needs.

## Structure Overview

```
core/
├── server/           # FastAPI application (HTTP layer)
├── agents/           # All Agents
│
├── memory/           # Agent memory system → /memory
├── voice/            # LiveKit voice integration → /voice
├── multimodal/       # Camera/video/image/realtime analysis
│
├── consciousness/    # Awareness layer → /awareness
│   ├── core/         # Event infrastructure (bus, kv_store)
│   ├── processing/   # Event pipeline (processor, correlator)
│   ├── awareness/    # State tracking (state, device_presence)
│   ├── bridges/      # External integrations (composio, frigate)
│   └── utils/        # Utilities (notifier, deep_links)
│
├── tools/            # What agents CAN DO → /tools
├── skills/           # What agents KNOW → /skills
├── specialists/      # Domain experts → /agents
├── teams/            # Multi-agent teams → /teams
├── workflow/         # Automated processes → /workflows
│
├── knowledge/        # RAG system → /knowledge
├── approvals/        # HITL handler → /approvals
├── mcp/              # MCP servers → /mcp
├── models/           # Model factory → /models
├── integrations/     # External services → /settings
├── sync/             # Database sync engine
│
├── db.py             # Database singleton
├── registry.py       # Central discovery
└── prompts/          # Shared prompts
```

## Usage

```python
from core import get_shared_db, get_model, get_shared_skills
from core.agent import create_legacy_agent
from core.tools import todoist, github
from core.registry import get_agent, get_team
```

---

## Full Module Breakdown

```
core/
├── agno/                   # ALL Agno features (141 files)
│   ├── agents/             # ✅ Agents + specialists
│   ├── teams/              # ✅ Teams
│   ├── workflows/          # ✅ Workflows
│   ├── tools/              # ✅ Tools
│   ├── skills/             # ✅ Skills
│   ├── knowledge/          # ✅ RAG
│   ├── memory/             # ✅ Memory
│   ├── learning/           # ✅ Learning
│   ├── models/             # ✅ Model factory (moved from platform)
│   ├── multimodal/         # ✅ Multimodal
│   └── voice/              # ✅ Voice
│
├── consciousness/          # Awareness system (40 files)
│
├── platform/               # Pure infrastructure (10 files)
│   ├── os/                 # ✅ OneMindOS app
│   ├── db/                 # ✅ Database
│   ├── middleware/         # ✅ JWT, RBAC
│   └── observability/      # ✅ Metrics
│
├── lifeos/                 # Life OS features (68 files)
│   ├── inbox/
│   ├── tasks/
│   ├── calendar/
│   ├── goals/
│   └── ... (13 life management features)
│
├── integrations/           # External systems (33 files)
│
├── orchestration/          # Coordination (29 files)
│
├── safety/                 # Safety (9 files)
│
└── shared/                 # Shared utilities (3 files)
```
