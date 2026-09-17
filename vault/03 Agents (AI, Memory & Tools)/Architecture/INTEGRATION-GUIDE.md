# Integration Guide — Step by Step

> How to build One Mind OS from the ground up. Start with the foundation, add layers incrementally.

## Prerequisites

- Python 3.11+
- PostgreSQL 15+ with pgvector extension
- NATS Server 2.10+
- Node.js 18+ (for Flutter tooling)
- Docker (optional, for containerized deployment)
- API keys: OpenAI (or Anthropic/Google), ElevenLabs (optional)

## Phase 1: Foundation (Week 1-2)

### 1.1 Set Up PostgreSQL + PgVector

```bash
# Docker
docker run -d --name onemind-db \
  -e POSTGRES_USER=ai \
  -e POSTGRES_PASSWORD=ai \
  -e POSTGRES_DB=onemind \
  -p 5532:5432 \
  pgvector/pgvector:pg16

# Enable pgvector extension
psql -h localhost -p 5532 -U ai -d onemind -c "CREATE EXTENSION IF NOT EXISTS vector;"
```

### 1.2 Set Up NATS Server

```bash
# Docker
docker run -d --name onemind-nats \
  -p 4222:4222 -p 8222:8222 \
  nats:latest --jetstream

# Or native
nats-server --jetstream --store_dir /data/nats
```

### 1.3 Install Agno

```bash
pip install agno[all]
# Or minimal:
pip install agno pgvector psycopg[binary] nats-py
```

### 1.4 Create Basic Agent (Verify It Works)

```python
# test_agent.py
from agno.agent import Agent
from agno.models.openai import OpenAIChat
from agno.db.postgres import PostgresDb

db = PostgresDb(
    db_url="postgresql+psycopg://ai:ai@localhost:5532/onemind",
    table_name="agent_sessions",
)

agent = Agent(
    name="Legacy",
    model=OpenAIChat(id="gpt-4.1"),
    db=db,
    enable_user_memories=True,
    instructions=[
        "You are Legacy, a direct and forward-thinking AI assistant.",
        "Get to the point. Don't sugar-coat responses.",
    ],
)

agent.print_response("Hello, who are you?", stream=True)
```

```bash
export OPENAI_API_KEY=your-key
python test_agent.py
```

## Phase 2: Agno Cognitive Layer (Week 2-3)

### 2.1 Set Up Knowledge Base

```python
# knowledge/setup.py
from agno.knowledge.combined import CombinedKnowledgeBase
from agno.knowledge.pdf_url import PDFUrlKnowledgeBase
from agno.knowledge.json import JSONKnowledgeBase
from agno.vectordb.pgvector import PgVector, SearchType

vector_db = PgVector(
    table_name="onemind_knowledge",
    db_url="postgresql+psycopg://ai:ai@localhost:5532/onemind",
    search_type=SearchType.hybrid,
)

knowledge = CombinedKnowledgeBase(
    sources=[
        # Add your knowledge sources here
        JSONKnowledgeBase(path="data/legacy_codex.json"),
    ],
    vector_db=vector_db,
)

# Load knowledge (run once, or when sources change)
# knowledge.load(recreate=False)
```

### 2.2 Build Agent Team

Create the full team structure as defined in [AGNO-LAYER.md](./AGNO-LAYER.md). Start with just the Digital agent and add others incrementally.

### 2.3 Deploy AgentOS

```python
# main.py
from agno.agent_os import AgentOS
from teams.onemind import one_mind

agent_os = AgentOS(
    description="One Mind OS — Legacy",
    teams=[one_mind],
)

# Run: uvicorn main:agent_os.app --host 0.0.0.0 --port 7777
```

### 2.4 Test via REST API

```bash
curl -X POST http://localhost:7777/teams/onemind/runs \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What can you do?",
    "session_id": "master-delacruz-main",
    "user_id": "master-delacruz"
  }'
```

## Phase 3: Interaction Gateway (Week 3-4)

### 3.1 Build Gateway Service

Create the FastAPI gateway as defined in [INTERACTION-GATEWAY.md](./INTERACTION-GATEWAY.md).

```python
# gateway/main.py
from fastapi import FastAPI, WebSocket
from gateway.session_router import send_to_agno, broadcast, is_duplicate
import nats

app = FastAPI(title="One Mind OS Gateway")

@app.websocket("/chat/{user_id}")
async def chat_channel(ws: WebSocket, user_id: str):
    await ws.accept()
    try:
        while True:
            data = await ws.receive_json()
            if not is_duplicate(data["text"]):
                response = await send_to_agno(user_id, data["text"])
                await broadcast(user_id, response)
    except Exception:
        pass
```

### 3.2 Add Voice Channel

Integrate Gemini Live API for real-time voice.

### 3.3 Add Channel Adapters

Start with Telegram (easiest), then add WhatsApp, Slack as needed.

### 3.4 Connect Flutter App

Point the Flutter app's WebSocket to the Gateway.

## Phase 4: Heartbeat System (Week 4-5)

### 4.1 Create HEARTBEAT.md

Create `workspace/HEARTBEAT.md` with your monitoring checklist as defined in [HEARTBEAT.md](./HEARTBEAT.md).

### 4.2 Build Heartbeat Service

Implement the `AgnoHeartbeat` class from [HEARTBEAT.md](./HEARTBEAT.md). Connect it to your data sources (Gmail, Calendar, ClickUp).

### 4.3 Test Heartbeat

```bash
# Run manually first
python heartbeat/run.py --once

# Then set up cron or the Python loop
python heartbeat/run.py --daemon --interval 15
```

## Phase 5: NATS Bridge (Week 5-6)

### 5.1 Build Robot Control Toolkit

Create the `RobotControlTools` Agno toolkit as defined in [AGNO-LAYER.md](./AGNO-LAYER.md).

### 5.2 Build Context Aggregator

Deploy the `ContextAggregator` service as defined in [NATS-BUS.md](./NATS-BUS.md).

### 5.3 Test Without a Robot

Use a NATS simulator to publish fake sensor data and verify Agno receives and reasons about it.

```python
# test/simulate_robot.py
import nats, json, asyncio

async def simulate():
    nc = await nats.connect("nats://localhost:4222")
    while True:
        await nc.publish(
            "onemind.vision.robot1",
            json.dumps({
                "robot_id": "robot1",
                "description": "Empty living room, lights on",
                "timestamp": time.time(),
            }).encode()
        )
        await asyncio.sleep(2)

asyncio.run(simulate())
```

## Phase 6: OM1 Embodiment (Week 6-8)

### 6.1 Install OM1

```bash
git clone https://github.com/OpenMind/OM1.git
cd OM1
pip install -r requirements.txt
```

### 6.2 Create Custom Plugins

Build the `AgnoNATSInput` and `NATSPublisherAction` plugins as defined in [OM1-LAYER.md](./OM1-LAYER.md).

### 6.3 Configure OM1

Create the JSON5 config with your custom plugins, robot hardware, and LLM settings.

### 6.4 Test with Simulator

OM1 includes a WebSocket simulator for testing without physical hardware. Use it to verify the full pipeline: Agno → NATS → OM1 → (simulated) actions → NATS → Agno.

### 6.5 Connect Physical Robot

Once the simulator works, connect to your actual Unitree G1/Go2 or TurtleBot.

## Phase 7: Personality and Polish (Week 8-9)

### 7.1 Create SOUL.md and USER.md

Define Legacy's personality as specified in [PERSONALITY.md](./PERSONALITY.md).

### 7.2 Load into Knowledge Base

Add SOUL.md and USER.md to Agno's knowledge base for personality-aware responses.

### 7.3 Sync Personality to OM1

Implement the personality projection via NATS as defined in [PERSONALITY.md](./PERSONALITY.md).

### 7.4 End-to-End Testing

Test the full flow: voice → Gateway → Agno → NATS → OM1 → physical action + digital response → Gateway → voice response.

## File Structure (Final)

```
onemindos-v2/
├── docs/
│   ├── ARCHITECTURE.md
│   ├── INTERACTION-GATEWAY.md
│   ├── AGNO-LAYER.md
│   ├── OM1-LAYER.md
│   ├── NATS-BUS.md
│   ├── HEARTBEAT.md
│   ├── UNIFIED-SESSION.md
│   ├── LLM-STRATEGY.md
│   ├── PERSONALITY.md
│   ├── SHARED-VISION.md
│   └── INTEGRATION-GUIDE.md
├── gateway/
│   ├── main.py
│   ├── session_router.py
│   ├── channels/
│   │   ├── voice.py
│   │   ├── chat.py
│   │   ├── telegram.py
│   │   └── nats_bridge.py
│   └── tts/
│       ├── gemini.py
│       └── elevenlabs.py
├── agno_layer/
│   ├── main.py
│   ├── agents/
│   │   ├── physical.py
│   │   ├── digital.py
│   │   ├── hp.py
│   │   ├── le.py
│   │   ├── ge.py
│   │   └── it.py
│   ├── teams/
│   │   └── onemind.py
│   ├── tools/
│   │   ├── robot_control.py
│   │   └── mcp_registry.py
│   ├── knowledge/
│   │   └── setup.py
│   └── memory/
│       └── setup.py
├── om1_plugins/
│   ├── inputs/
│   │   └── agno_nats_input.py
│   └── actions/
│       └── nats_publisher.py
├── heartbeat/
│   ├── run.py
│   ├── gatherers/
│   │   ├── gmail.py
│   │   ├── calendar.py
│   │   ├── clickup.py
│   │   └── nats_fleet.py
│   └── workspace/
│       └── HEARTBEAT.md
├── context_aggregator/
│   └── aggregator.py
├── personality/
│   ├── SOUL.md
│   └── USER.md
├── config/
│   ├── agno.env
│   ├── nats.conf
│   └── om1.json5
├── docker-compose.yml
└── README.md
```

## Docker Compose (Development)

```yaml
version: '3.8'

services:
  postgres:
    image: pgvector/pgvector:pg16
    environment:
      POSTGRES_USER: ai
      POSTGRES_PASSWORD: ai
      POSTGRES_DB: onemind
    ports:
      - "5532:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  nats:
    image: nats:latest
    command: --jetstream
    ports:
      - "4222:4222"
      - "8222:8222"

  agno:
    build: ./agno_layer
    ports:
      - "7777:7777"
    environment:
      - DATABASE_URL=postgresql+psycopg://ai:ai@postgres:5432/onemind
      - NATS_URL=nats://nats:4222
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    depends_on:
      - postgres
      - nats

  gateway:
    build: ./gateway
    ports:
      - "8000:8000"
    environment:
      - AGNO_URL=http://agno:7777
      - NATS_URL=nats://nats:4222
      - GEMINI_API_KEY=${GEMINI_API_KEY}
    depends_on:
      - agno
      - nats

  heartbeat:
    build: ./heartbeat
    environment:
      - AGNO_URL=http://agno:7777
      - NATS_URL=nats://nats:4222
    depends_on:
      - agno
      - nats

  aggregator:
    build: ./context_aggregator
    environment:
      - NATS_URL=nats://nats:4222
    depends_on:
      - nats

volumes:
  pgdata:
```

## Verification Checklist

After each phase, verify:

- [ ] **Phase 1**: Agent responds to text prompts with memory persistence
- [ ] **Phase 2**: Team delegates to correct specialist agent; RAG returns relevant knowledge
- [ ] **Phase 3**: Can chat via Flutter WebSocket; voice STT/TTS works; Telegram bot responds
- [ ] **Phase 4**: Heartbeat runs on schedule; notifications arrive; HEARTBEAT_OK when idle
- [ ] **Phase 5**: Fake robot events flow through NATS → Agno reasons about them → responds
- [ ] **Phase 6**: Physical robot receives commands and moves; sensor data flows back to Agno
- [ ] **Phase 7**: Same personality across voice, text, and robot speech; unified session works across channels
