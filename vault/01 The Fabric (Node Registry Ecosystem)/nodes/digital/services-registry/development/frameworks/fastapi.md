# FastAPI

> Modern, fast Python web framework for building APIs with automatic OpenAPI documentation.

---

## Overview

| Attribute         | Value                                                |
| ----------------- | ---------------------------------------------------- |
| **Type**          | Web Framework / API Server                           |
| **Language**      | Python 3.8+                                          |
| **Default Port**  | 8000                                                 |
| **License**       | MIT                                                  |
| **Documentation** | [fastapi.tiangolo.com](https://fastapi.tiangolo.com) |

---

## What is FastAPI?

FastAPI is a modern Python web framework that provides:

- **High performance** — On par with NodeJS and Go
- **Type hints** — Full Python type annotation support
- **Async/await** — Native async support
- **Auto documentation** — OpenAPI (Swagger) and ReDoc
- **Validation** — Automatic request validation via Pydantic
- **Dependency injection** — Clean, testable architecture

---

## OneMind Usage

### Primary Use Cases

```text
┌─────────────────────────────────────────────────────────────────┐
│                      OneMind FastAPI                             │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Agent API Gateway                      │   │
│  │  /agents • /runs • /sessions • /awareness                │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    WebSocket Streaming                    │   │
│  │  Real-time agent responses • Event streams               │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Webhook Handlers                       │   │
│  │  Nango callbacks • NATS bridge • External services       │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

| Use Case           | Description                           |
| ------------------ | ------------------------------------- |
| **Agent API**      | REST endpoints for Legacy agent       |
| **Streaming**      | SSE/WebSocket for real-time responses |
| **Webhooks**       | Handle callbacks from integrations    |
| **Health Checks**  | Kubernetes/Docker health endpoints    |
| **LiveKit Bridge** | Voice session management              |
| **NATS Gateway**   | HTTP to NATS message bridge           |

---

## Installation

```bash
# Core FastAPI
pip install fastapi

# With all extras
pip install "fastapi[all]"

# Production server
pip install uvicorn[standard]

# Common dependencies
pip install pydantic pydantic-settings python-multipart
```

---

## Basic Application

### Minimal App

```python
from fastapi import FastAPI

app = FastAPI(
    title="OneMind API",
    description="Legacy Agent API Gateway",
    version="1.0.0"
)

@app.get("/")
async def root():
    return {"message": "OneMind API"}

@app.get("/health")
async def health():
    return {"status": "healthy"}
```

### Run Server

```bash
# Development
uvicorn main:app --reload --host 0.0.0.0 --port 8000

# Production
uvicorn main:app --host 0.0.0.0 --port 8000 --workers 4
```

---

## OneMind Agent API

### Application Structure

```python
# main.py
from contextlib import asynccontextmanager
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from api.routes import agents, sessions, awareness, workflows
from core.config import settings
from core.db import init_db, close_db


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    await init_db()
    yield
    # Shutdown
    await close_db()


app = FastAPI(
    title="OneMind Legacy API",
    version="1.0.0",
    lifespan=lifespan
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Routes
app.include_router(agents.router, prefix="/agents", tags=["Agents"])
app.include_router(sessions.router, prefix="/sessions", tags=["Sessions"])
app.include_router(awareness.router, prefix="/awareness", tags=["Awareness"])
app.include_router(workflows.router, prefix="/workflows", tags=["Workflows"])


@app.get("/health")
async def health():
    return {"status": "healthy", "version": "1.0.0"}
```

### Agent Routes

```python
# api/routes/agents.py
from fastapi import APIRouter, HTTPException, Depends
from pydantic import BaseModel
from typing import Optional
from datetime import datetime

router = APIRouter()


class AgentRunRequest(BaseModel):
    message: str
    session_id: Optional[str] = None
    context: Optional[dict] = None


class AgentRunResponse(BaseModel):
    run_id: str
    session_id: str
    response: str
    created_at: datetime


@router.get("/")
async def list_agents():
    """List available agents."""
    return {
        "agents": [
            {"id": "legacy", "name": "Legacy", "status": "active"}
        ]
    }


@router.get("/{agent_id}")
async def get_agent(agent_id: str):
    """Get agent details."""
    if agent_id != "legacy":
        raise HTTPException(status_code=404, detail="Agent not found")

    return {
        "id": "legacy",
        "name": "Legacy",
        "model": "claude-3-5-sonnet",
        "status": "active",
        "tools": ["search", "calendar", "tasks", "home"]
    }


@router.post("/{agent_id}/runs", response_model=AgentRunResponse)
async def run_agent(agent_id: str, request: AgentRunRequest):
    """Execute agent with message."""
    from core.agent import legacy_agent

    result = await legacy_agent.run(
        message=request.message,
        session_id=request.session_id,
        context=request.context
    )

    return AgentRunResponse(
        run_id=result.run_id,
        session_id=result.session_id,
        response=result.response,
        created_at=datetime.utcnow()
    )
```

---

## Streaming Responses

### Server-Sent Events (SSE)

```python
from fastapi import APIRouter
from fastapi.responses import StreamingResponse
from typing import AsyncGenerator
import json

router = APIRouter()


async def generate_stream(message: str) -> AsyncGenerator[str, None]:
    """Stream agent response tokens."""
    from core.agent import legacy_agent

    async for chunk in legacy_agent.stream(message):
        data = json.dumps({"token": chunk, "done": False})
        yield f"data: {data}\n\n"

    yield f"data: {json.dumps({'done': True})}\n\n"


@router.post("/agents/{agent_id}/stream")
async def stream_agent(agent_id: str, request: AgentRunRequest):
    """Stream agent response via SSE."""
    return StreamingResponse(
        generate_stream(request.message),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
        }
    )
```

### WebSocket

```python
from fastapi import WebSocket, WebSocketDisconnect
from typing import Dict

# Connection manager
class ConnectionManager:
    def __init__(self):
        self.active_connections: Dict[str, WebSocket] = {}

    async def connect(self, session_id: str, websocket: WebSocket):
        await websocket.accept()
        self.active_connections[session_id] = websocket

    def disconnect(self, session_id: str):
        self.active_connections.pop(session_id, None)

    async def send(self, session_id: str, message: dict):
        if ws := self.active_connections.get(session_id):
            await ws.send_json(message)


manager = ConnectionManager()


@router.websocket("/ws/{session_id}")
async def websocket_endpoint(websocket: WebSocket, session_id: str):
    await manager.connect(session_id, websocket)

    try:
        while True:
            data = await websocket.receive_json()

            # Process message
            async for chunk in legacy_agent.stream(data["message"]):
                await manager.send(session_id, {
                    "type": "token",
                    "content": chunk
                })

            await manager.send(session_id, {"type": "done"})

    except WebSocketDisconnect:
        manager.disconnect(session_id)
```

---

## Dependency Injection

### Database Session

```python
from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession
from typing import AsyncGenerator

from core.db import async_session_maker


async def get_db() -> AsyncGenerator[AsyncSession, None]:
    async with async_session_maker() as session:
        try:
            yield session
            await session.commit()
        except Exception:
            await session.rollback()
            raise


@router.get("/sessions")
async def list_sessions(db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Session))
    return result.scalars().all()
```

### Current User

```python
from fastapi import Depends, HTTPException, Security
from fastapi.security import APIKeyHeader

api_key_header = APIKeyHeader(name="X-API-Key")


async def get_current_user(api_key: str = Security(api_key_header)):
    if api_key != settings.API_KEY:
        raise HTTPException(status_code=403, detail="Invalid API key")
    return {"user": "zeus", "role": "admin"}


@router.get("/me")
async def get_me(user: dict = Depends(get_current_user)):
    return user
```

### Settings

```python
# core/config.py
from pydantic_settings import BaseSettings
from functools import lru_cache


class Settings(BaseSettings):
    # API
    API_KEY: str
    CORS_ORIGINS: list[str] = ["http://localhost:3000"]

    # Database
    DATABASE_URL: str

    # NATS
    NATS_URL: str = "nats://localhost:4222"

    # Redis
    REDIS_URL: str = "redis://localhost:6379"

    class Config:
        env_file = ".env"


@lru_cache
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
```

---

## Request Validation

### Pydantic Models

```python
from pydantic import BaseModel, Field, validator
from typing import Optional, Literal
from datetime import datetime


class AgentMessage(BaseModel):
    content: str = Field(..., min_length=1, max_length=10000)
    role: Literal["user", "assistant", "system"] = "user"
    metadata: Optional[dict] = None

    @validator("content")
    def content_not_empty(cls, v):
        if not v.strip():
            raise ValueError("Content cannot be empty")
        return v.strip()


class AwarenessUpdate(BaseModel):
    level: Literal["dormant", "passive", "active", "focused", "immersed"]
    trigger: Optional[str] = None

    class Config:
        json_schema_extra = {
            "example": {
                "level": "active",
                "trigger": "user_message"
            }
        }


class WorkflowTrigger(BaseModel):
    name: str
    params: Optional[dict] = {}
    scheduled_at: Optional[datetime] = None
```

### Path and Query Parameters

```python
from fastapi import Query, Path
from typing import Optional


@router.get("/sessions")
async def list_sessions(
    limit: int = Query(default=10, ge=1, le=100),
    offset: int = Query(default=0, ge=0),
    status: Optional[str] = Query(default=None, regex="^(active|completed|archived)$")
):
    """List sessions with pagination."""
    pass


@router.get("/sessions/{session_id}")
async def get_session(
    session_id: str = Path(..., min_length=36, max_length=36)
):
    """Get session by ID (UUID format)."""
    pass
```

---

## Error Handling

### Exception Handlers

```python
from fastapi import FastAPI, Request, HTTPException
from fastapi.responses import JSONResponse
from fastapi.exceptions import RequestValidationError


app = FastAPI()


class AgentError(Exception):
    def __init__(self, message: str, code: str = "AGENT_ERROR"):
        self.message = message
        self.code = code


@app.exception_handler(AgentError)
async def agent_error_handler(request: Request, exc: AgentError):
    return JSONResponse(
        status_code=500,
        content={
            "error": exc.code,
            "message": exc.message,
            "path": str(request.url)
        }
    )


@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request: Request, exc: RequestValidationError):
    return JSONResponse(
        status_code=422,
        content={
            "error": "VALIDATION_ERROR",
            "details": exc.errors()
        }
    )


@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception):
    return JSONResponse(
        status_code=500,
        content={
            "error": "INTERNAL_ERROR",
            "message": "An unexpected error occurred"
        }
    )
```

---

## Middleware

### Request Logging

```python
import time
import logging
from fastapi import Request

logger = logging.getLogger(__name__)


@app.middleware("http")
async def log_requests(request: Request, call_next):
    start_time = time.time()

    response = await call_next(request)

    process_time = time.time() - start_time
    logger.info(
        f"{request.method} {request.url.path} - "
        f"{response.status_code} - {process_time:.3f}s"
    )

    response.headers["X-Process-Time"] = str(process_time)
    return response
```

### Rate Limiting

```python
from fastapi import Request, HTTPException
from collections import defaultdict
import time


class RateLimiter:
    def __init__(self, requests_per_minute: int = 60):
        self.requests_per_minute = requests_per_minute
        self.requests: dict = defaultdict(list)

    def is_allowed(self, client_id: str) -> bool:
        now = time.time()
        minute_ago = now - 60

        # Clean old requests
        self.requests[client_id] = [
            t for t in self.requests[client_id] if t > minute_ago
        ]

        if len(self.requests[client_id]) >= self.requests_per_minute:
            return False

        self.requests[client_id].append(now)
        return True


limiter = RateLimiter(requests_per_minute=60)


@app.middleware("http")
async def rate_limit_middleware(request: Request, call_next):
    client_ip = request.client.host

    if not limiter.is_allowed(client_ip):
        raise HTTPException(status_code=429, detail="Too many requests")

    return await call_next(request)
```

---

## Background Tasks

### Simple Background Task

```python
from fastapi import BackgroundTasks


async def log_agent_run(run_id: str, message: str):
    """Log agent run to database."""
    # Async database operation
    await db.execute(
        "INSERT INTO agent_logs (run_id, message) VALUES ($1, $2)",
        run_id, message
    )


@router.post("/agents/{agent_id}/runs")
async def run_agent(
    agent_id: str,
    request: AgentRunRequest,
    background_tasks: BackgroundTasks
):
    result = await legacy_agent.run(request.message)

    # Queue background task
    background_tasks.add_task(log_agent_run, result.run_id, request.message)

    return result
```

### With Task Queue (Redis)

```python
from arq import create_pool
from arq.connections import RedisSettings


async def get_redis_pool():
    return await create_pool(RedisSettings(host="redis"))


@router.post("/workflows/{workflow_name}")
async def trigger_workflow(workflow_name: str):
    redis = await get_redis_pool()

    # Queue job
    job = await redis.enqueue_job(
        "run_workflow",
        workflow_name
    )

    return {"job_id": job.job_id, "status": "queued"}
```

---

## Docker Compose

```yaml
services:
  api:
    build:
      context: ./agents/legacy
      dockerfile: Dockerfile
    container_name: onemind-api
    restart: unless-stopped
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql://onemind:${DB_PASSWORD}@postgres:5432/onemind
      - REDIS_URL=redis://redis:6379
      - NATS_URL=nats://nats:4222
      - API_KEY=${API_KEY}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - onemind-network
```

### Dockerfile

```dockerfile
FROM python:3.12-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application
COPY . .

# Run with uvicorn
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
```

---

## Testing

### Test Client

```python
from fastapi.testclient import TestClient
from main import app
import pytest

client = TestClient(app)


def test_health():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json()["status"] == "healthy"


def test_list_agents():
    response = client.get("/agents")
    assert response.status_code == 200
    assert len(response.json()["agents"]) > 0


@pytest.mark.asyncio
async def test_agent_run():
    response = client.post(
        "/agents/legacy/runs",
        json={"message": "Hello"}
    )
    assert response.status_code == 200
    assert "response" in response.json()
```

### Async Tests

```python
import pytest
from httpx import AsyncClient
from main import app


@pytest.mark.asyncio
async def test_async_endpoint():
    async with AsyncClient(app=app, base_url="http://test") as client:
        response = await client.get("/health")
        assert response.status_code == 200
```

---

## Production Deployment

### Gunicorn + Uvicorn Workers

```bash
# Production command
gunicorn main:app \
    --workers 4 \
    --worker-class uvicorn.workers.UvicornWorker \
    --bind 0.0.0.0:8000 \
    --access-logfile - \
    --error-logfile -
```

### Environment Variables

```bash
# API Configuration
API_KEY=your_secure_api_key
CORS_ORIGINS=["https://onemind.cloud"]

# Database
DATABASE_URL=postgresql://user:pass@host:5432/onemind

# Services
REDIS_URL=redis://redis:6379
NATS_URL=nats://nats:4222

# Logging
LOG_LEVEL=info
```

### Health Check Endpoint

```python
from fastapi import APIRouter
from datetime import datetime

router = APIRouter()


@router.get("/health")
async def health():
    return {
        "status": "healthy",
        "timestamp": datetime.utcnow().isoformat(),
        "version": "1.0.0"
    }


@router.get("/health/ready")
async def readiness():
    """Check all dependencies."""
    checks = {}

    # Database
    try:
        await db.execute("SELECT 1")
        checks["database"] = "ok"
    except Exception:
        checks["database"] = "error"

    # Redis
    try:
        await redis.ping()
        checks["redis"] = "ok"
    except Exception:
        checks["redis"] = "error"

    all_ok = all(v == "ok" for v in checks.values())

    return {
        "status": "ready" if all_ok else "not_ready",
        "checks": checks
    }
```

---

## OpenAPI Documentation

### Access Docs

- **Swagger UI**: `http://localhost:8000/docs`
- **ReDoc**: `http://localhost:8000/redoc`
- **OpenAPI JSON**: `http://localhost:8000/openapi.json`

### Custom Documentation

```python
from fastapi import FastAPI

app = FastAPI(
    title="OneMind Legacy API",
    description="""
    ## Legacy Agent API

    This API provides access to the Legacy AI agent.

    ### Features
    - Execute agent runs
    - Stream responses
    - Manage sessions
    - Control awareness levels
    """,
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_tags=[
        {"name": "Agents", "description": "Agent operations"},
        {"name": "Sessions", "description": "Session management"},
        {"name": "Awareness", "description": "Awareness state control"},
    ]
)
```

---

## Best Practices

### Project Structure

```text
agents/legacy/
├── main.py              # FastAPI app entry
├── api/
│   ├── __init__.py
│   ├── deps.py          # Dependencies
│   └── routes/
│       ├── agents.py
│       ├── sessions.py
│       └── awareness.py
├── core/
│   ├── config.py        # Settings
│   ├── db.py            # Database
│   └── agent.py         # Agent logic
├── models/
│   ├── schemas.py       # Pydantic models
│   └── db.py            # SQLAlchemy models
└── tests/
    └── test_api.py
```

### Performance Tips

| Tip                    | Description                        |
| ---------------------- | ---------------------------------- |
| **Use async**          | All I/O operations should be async |
| **Connection pooling** | Reuse DB/Redis connections         |
| **Response caching**   | Cache expensive computations       |
| **Pagination**         | Always paginate list endpoints     |
| **Background tasks**   | Offload non-critical work          |

---

## Troubleshooting

### Import Errors

```bash
# Verify installation
pip show fastapi uvicorn

# Reinstall if needed
pip install --force-reinstall fastapi[all]
```

### CORS Issues

```python
# Allow all origins (development only)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
```

### Slow Startup

```python
# Use lifespan for expensive initialization
@asynccontextmanager
async def lifespan(app: FastAPI):
    # Load models, connect to DB once
    await init_resources()
    yield
    await cleanup_resources()
```

### WebSocket Disconnect

```python
# Always handle disconnect gracefully
try:
    while True:
        data = await websocket.receive_json()
except WebSocketDisconnect:
    # Clean up
    manager.disconnect(session_id)
except Exception as e:
    logger.error(f"WebSocket error: {e}")
    await websocket.close(code=1011)
```

---

## Resources

- [FastAPI Documentation](https://fastapi.tiangolo.com)
- [FastAPI GitHub](https://github.com/tiangolo/fastapi)
- [Pydantic v2](https://docs.pydantic.dev)
- [Uvicorn](https://www.uvicorn.org)
- [Starlette](https://www.starlette.io)

---

_Last updated: January 2026_


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
