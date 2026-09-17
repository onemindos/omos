# OneMind OS Unified Architecture Plan
## Legacy AI - Solarpunk Tactical Operating System

### Vision Summary
A unified operating system where everything is a connectable node/asset - humans, agents, robots, drones, locations, machines, vehicles. The goal: manage off-grid homestead infrastructure with AI + robotics for food production, maintenance, security, etc.

---

## Key Clarifications (User Feedback)

### 1. Operations vs Assets
- **Operations** (NOT assets): Inbox, Tasks, Calendar - these are command center screens
- **Assets**: Skills and Tools belong to assets, not as standalone features

### 2. Universal Chat Interface
- NOT agent-default - should be a universal interface
- Can be directed to: specific agent, team, OR use tools directly
- Default model from model factory (configurable)
- Can choose Gemini for multimodal when needed
- Can chat with teams but not make teams the default

### 3. Additional Requirements
- **Voice**: Need voice input/output integration
- **Multi-user**: Different users with their own context
- **OpenMind**: Robots show up via OpenMind heartbeat loop (`backend/openmind/bridge.py`)
- **Agno**: Designed for agents - universal chat needs tool access too

---

## Current State Analysis

### What Already Exists (Solid Foundation)

| System | Location | Status |
|--------|----------|--------|
| **Asset Tracker** | `backend/assets/tracker.py` | 95% complete |
| **Soul System** | `backend/agents/souls/__init__.py` | Ready |
| **Capability Registry** | `backend/capabilities/__init__.py` | Ready |
| **Workflow Engine** | `backend/workflows/engine.py` | Ready |
| **Tools Registry** | `backend/tools/registry.py` | Ready |
| **Model Factory** | `backend/models.py` | Ready |
| **OpenMind Bridge** | `backend/openmind/bridge.py` | Ready |
| **DB Schema** | `infra/db/003_expansion_schema.sql` | Ready |
| **Chat Screen** | `frontend/lib/screens/chat_screen.dart` | Agent-focused (needs update) |

### Asset Types Already Supported
- `LOCATION` - buildings, zones, outdoor areas
- `DEVICE` - phones, sensors, cameras, beacons
- `HUMAN` - with biometrics (heart_rate, blood_oxygen, steps, stress)
- `MACHINE` - robots, drones, vehicles, 3D printers, CNC

### Model Factory (`backend/models.py`)
- OpenAI: gpt-4o, gpt-4o-mini, o1, o1-mini
- AWS Bedrock: Claude Sonnet 4.5
- Google Gemini: (add for multimodal)
- Azure OpenAI: (configurable)
- Groq: (fast inference)

---

## Implementation Plan

### Phase 1: Universal Chat Interface
**Critical: Replace agent-focused chat with universal interface**

#### 1.1 Backend: Universal Chat Endpoint
**File: `backend/chat/universal.py`**
```python
# New universal chat handler
class UniversalChat:
    """
    Universal chat interface that can:
    1. Use default model directly (no agent)
    2. Route to specific agent
    3. Route to team
    4. Access tools directly
    """

    async def chat(
        self,
        message: str,
        mode: str = "direct",  # direct | agent | team
        target_id: Optional[str] = None,  # agent_id or team_id
        model_id: Optional[str] = None,  # override default model
        user_id: str = "default",
        tools: List[str] = [],  # tools to enable
    ):
        if mode == "direct":
            # Use model factory directly with optional tools
            model = get_model(model_id or "gpt-4o")
            # ...
        elif mode == "agent":
            # Route to specific agent
            agent = get_agent(target_id)
            # ...
        elif mode == "team":
            # Route to team
            team = get_team(target_id)
            # ...
```

#### 1.2 API Endpoint
```
POST /api/chat/universal
{
    "message": "...",
    "mode": "direct" | "agent" | "team",
    "target_id": "optional-agent-or-team-id",
    "model_id": "gpt-4o" | "gemini-2.0-flash" | etc,
    "user_id": "user-123",
    "tools": ["web_search", "home_assistant"]
}
```

#### 1.3 Frontend: Update Chat Screen
**File: `frontend/lib/screens/chat_screen.dart`**

Changes:
- Remove agent-focused dropdown
- Add mode selector: Direct / Agent / Team
- Add model selector (from model factory)
- Add tool picker
- Add voice input button
- Add user context

```dart
// Mode selector
enum ChatMode { direct, agent, team }

// Chat screen state
ChatMode _chatMode = ChatMode.direct;
String? _selectedModel = 'gpt-4o';
String? _selectedAgentId;
String? _selectedTeamId;
List<String> _enabledTools = [];
```

#### 1.4 Model Selector Widget
**New: `frontend/lib/widgets/model_selector.dart`**
- Fetch available models from `/api/models`
- Show provider grouping (OpenAI, Gemini, Claude)
- Indicate multimodal support (Gemini, GPT-4o)

---

### Phase 2: Voice Integration
**Add voice input/output to chat**

#### 2.1 Backend: Voice Endpoints
**File: `backend/voice/service.py`**
```python
# Speech-to-Text (STT)
POST /api/voice/transcribe
# Input: audio file
# Output: text

# Text-to-Speech (TTS)
POST /api/voice/synthesize
# Input: text
# Output: audio URL
```

#### 2.2 Frontend: Voice Button
- Add microphone button to chat input
- Use Web Speech API (browser) or backend STT
- Add speaker button to play responses

---

### Phase 3: Multi-User Support
**Add user context to all operations**

#### 3.1 Database: Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    preferences JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Add user_id to existing tables
ALTER TABLE tasks ADD COLUMN user_id UUID REFERENCES users(id);
ALTER TABLE memories ADD COLUMN user_id UUID;
ALTER TABLE sessions ADD COLUMN user_id UUID;
```

#### 3.2 Backend: User Context Middleware
- Extract user from token/session
- Pass user_id to all services
- Filter queries by user

#### 3.3 Frontend: User Provider
- Store current user in Riverpod
- Pass user_id to all API calls

---

### Phase 4: Sidebar Reorganization
**Proper separation: Operations vs Assets**

#### 4.1 New Sidebar Structure
```
COMMAND CENTER (Operations)
├── Inbox          → /
├── Tasks          → /tasks
├── Calendar       → /calendar
└── Chat           → /chat (NEW: Universal Chat)

ASSETS & ENTITIES
├── All Assets     → /assets (NEW: unified browser)
├── Agents         → /agents
├── Teams          → /teams
├── Locations      → /locations (filter of assets)
└── Machines       → /machines (filter of assets)

WORKSPACE
├── Projects       → /projects
├── Documents      → /documents
├── Sheets         → /sheets
├── Knowledge      → /knowledge
└── Memories       → /memories

AUTOMATION
├── Workflows      → /workflows
├── Builder        → /workflow-builder
├── Approvals      → /approvals
└── Sessions       → /sessions

MONITORING
├── System Pulse   → /system
├── Topology       → /topology
├── Events         → /events
├── Activity       → /activity
└── Analytics      → /analytics

SETTINGS
├── Profile        → /profile
├── Settings       → /settings
├── NATS Control   → /nats-control
└── Field Assets   → /field-assets
```

#### 4.2 Files to Modify
- `frontend/lib/widgets/sidebar.dart` - New structure
- `frontend/lib/main.dart` - Add routes

---

### Phase 5: Skills & Tools (Asset Attachments)
**Skills and Tools belong to Assets**

#### 5.1 Database Schema
**File: `infra/db/004_skills_schema.sql`**
```sql
CREATE TABLE skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    category TEXT,  -- physical, digital, hybrid
    difficulty INT DEFAULT 1,
    prerequisites UUID[],
    tools_required TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE asset_skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) ON DELETE CASCADE,
    skill_id UUID REFERENCES skills(id) ON DELETE CASCADE,
    progress FLOAT DEFAULT 0,  -- 0-100
    status TEXT DEFAULT 'not_started',
    training_logs JSONB DEFAULT '[]',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

CREATE TABLE asset_tools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) ON DELETE CASCADE,
    tool_name TEXT NOT NULL,  -- matches tools registry
    enabled BOOLEAN DEFAULT TRUE,
    config JSONB DEFAULT '{}',
    permissions JSONB DEFAULT '{}'
);
```

#### 5.2 Asset Detail Screen Enhancement
When viewing any asset, show:
- **Overview**: Type, status, location, telemetry
- **Skills**: Assigned skills with progress
- **Tools**: Available tools for this asset
- **Tasks**: Assigned tasks
- **Docs**: Attached documents

---

### Phase 6: OpenMind Robot Integration
**Robots appear via OpenMind heartbeat loop**

#### 6.1 How OpenMind Works
```
backend/openmind/bridge.py:
├── OpenMindLoop (heartbeat at 0.1Hz default)
├── DigitalSensorySystem (gather digital data)
├── Fuser (create situation report)
├── CortexLLM (Agno agent decides actions)
└── ActionOrchestrator (execute)

NATS Topics:
├── openmind.sense.digital (sensory data)
├── openmind.cortex.plan (AI decisions)
└── asset.{type}.{id}.telemetry (machine updates)
```

#### 6.2 Robot Asset Registration
When OpenMind detects a robot:
1. Create/update asset in `assets` table
2. Publish to `asset.machine.{id}.telemetry`
3. WebSocket broadcasts to frontend
4. Asset browser shows real-time status

#### 6.3 Integration Points
- `backend/openmind/sensors/` - Add robot sensors
- `backend/assets/tracker.py` - Link to OpenMind data
- `frontend/lib/screens/asset_detail_screen.dart` - Show live telemetry

---

## File Manifest

### New Files
| File | Purpose |
|------|---------|
| `backend/chat/universal.py` | Universal chat handler |
| `backend/chat/__init__.py` | Module init |
| `backend/voice/service.py` | Voice STT/TTS |
| `infra/db/004_skills_schema.sql` | Skills + asset_tools |
| `frontend/lib/widgets/model_selector.dart` | Model picker |
| `frontend/lib/widgets/voice_input.dart` | Voice button |
| `frontend/lib/screens/assets_browser_screen.dart` | Unified asset view |
| `frontend/lib/screens/asset_detail_screen.dart` | Asset drill-down |

### Files To Modify
| File | Changes |
|------|---------|
| `backend/main.py` | Register universal chat routes |
| `backend/models.py` | Add Gemini for multimodal |
| `frontend/lib/screens/chat_screen.dart` | Universal chat UI |
| `frontend/lib/widgets/sidebar.dart` | New nav structure |
| `frontend/lib/main.dart` | Add new routes |

---

## API Design Summary

### Universal Chat API
```
POST /api/chat/universal
GET  /api/models                    # List available models
GET  /api/tools                     # List available tools
```

### Assets API (Enhanced)
```
GET    /api/assets/                 # List all with filters
POST   /api/assets/                 # Create asset
GET    /api/assets/{id}             # Get with relationships
PATCH  /api/assets/{id}             # Update
DELETE /api/assets/{id}             # Delete
GET    /api/assets/{id}/skills      # Get skills
POST   /api/assets/{id}/skills      # Assign skill
GET    /api/assets/{id}/tools       # Get tools
POST   /api/assets/{id}/tools       # Enable tool
```

### Voice API
```
POST /api/voice/transcribe          # Audio → Text
POST /api/voice/synthesize          # Text → Audio
```

### Users API
```
GET  /api/users/me                  # Current user
POST /api/users/                    # Create user
PATCH /api/users/{id}               # Update preferences
```

---

## Key Design Decisions

### 1. Universal Chat (Not Agent-Default)
Chat interface is model-first, not agent-first. Users can:
- Chat directly with models (default)
- Optionally route to agents
- Optionally route to teams
- Enable specific tools

### 2. Skills & Tools Belong to Assets
Skills and tools are attachments to assets (humans, machines, agents), not standalone features.

### 3. Operations vs Assets
- **Operations**: Inbox, Tasks, Calendar, Chat (command center)
- **Assets**: Humans, Agents, Machines, Locations (entities)

### 4. OpenMind for Robots
Robots register via OpenMind heartbeat loop, not manual creation. Live telemetry flows through NATS.

### 5. Multi-Model Support
Model factory supports multiple providers:
- GPT-4o (default)
- Gemini (multimodal)
- Claude (reasoning)
- User can switch per conversation

---

## Implementation Order

1. **Phase 1**: Universal Chat Interface (backend + frontend)
2. **Phase 2**: Voice Integration (STT/TTS)
3. **Phase 3**: Multi-User Support
4. **Phase 4**: Sidebar Reorganization
5. **Phase 5**: Skills & Tools for Assets
6. **Phase 6**: OpenMind Robot Integration

---

## Critical Files to Read Before Implementation

1. `backend/models.py` - Model factory (add Gemini)
2. `backend/tools/registry.py` - Available tools
3. `backend/openmind/bridge.py` - OpenMind heartbeat
4. `frontend/lib/screens/chat_screen.dart` - Current chat (to update)
5. `frontend/lib/widgets/sidebar.dart` - Current nav (to reorganize)
