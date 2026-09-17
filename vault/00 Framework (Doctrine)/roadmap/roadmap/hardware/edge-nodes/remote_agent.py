"""
Remote Agent - Legacy's Local Presence at Edge Locations
========================================================

This agent runs at each edge location and provides Legacy with:
- Local shell command execution
- Home Assistant device control
- Frigate camera feed access
- Offline LLM via Ollama when cloud unavailable
- NATS message bridging between edge and cloud

Environment Variables:
    EDGE_LOCATION       - Location identifier (home, farm, office)
    AGENT_NAME          - Agent name for NATS (remote-{location})
    NATS_URL            - Local NATS leaf node URL
    HOMEASSISTANT_URL   - Home Assistant API URL
    HOMEASSISTANT_TOKEN - Home Assistant long-lived access token
    FRIGATE_URL         - Frigate NVR API URL
    OLLAMA_URL          - Ollama API URL for offline LLM
    OFFLINE_LLM_MODEL   - Model to use when offline (phi3:mini)
"""

import os
import asyncio
import subprocess
from typing import Optional
from contextlib import asynccontextmanager

import httpx
import nats
from fastapi import FastAPI
from pydantic import BaseModel

# =============================================================================
# Configuration
# =============================================================================

EDGE_LOCATION = os.getenv("EDGE_LOCATION", "unknown")
AGENT_NAME = os.getenv("AGENT_NAME", f"remote-{EDGE_LOCATION}")
NATS_URL = os.getenv("NATS_URL", "nats://localhost:4222")
HOMEASSISTANT_URL = os.getenv("HOMEASSISTANT_URL", "http://localhost:8123")
HOMEASSISTANT_TOKEN = os.getenv("HOMEASSISTANT_TOKEN", "")
FRIGATE_URL = os.getenv("FRIGATE_URL", "http://localhost:5000")
OLLAMA_URL = os.getenv("OLLAMA_URL", "http://localhost:11434")
OFFLINE_LLM_MODEL = os.getenv("OFFLINE_LLM_MODEL", "phi3:mini")

# =============================================================================
# NATS Client
# =============================================================================

nc: Optional[nats.NATS] = None


async def connect_nats():
    """Connect to NATS and subscribe to command topics."""
    global nc
    nc = await nats.connect(NATS_URL)
    print(f"[{AGENT_NAME}] Connected to NATS at {NATS_URL}")

    # Subscribe to commands for this location
    await nc.subscribe(f"commands.{EDGE_LOCATION}.*", cb=handle_command)
    await nc.subscribe(f"commands.{EDGE_LOCATION}.shell", cb=handle_shell_command)
    await nc.subscribe(f"commands.{EDGE_LOCATION}.homeassistant", cb=handle_ha_command)

    print(f"[{AGENT_NAME}] Subscribed to commands.{EDGE_LOCATION}.*")


async def handle_command(msg):
    """Handle generic commands."""
    subject = msg.subject
    data = msg.data.decode()
    print(f"[{AGENT_NAME}] Received command on {subject}: {data}")
    # TODO: Implement command routing


async def handle_shell_command(msg):
    """Execute shell commands locally."""
    import json
    try:
        data = json.loads(msg.data.decode())
        command = data.get("command", "")
        timeout = data.get("timeout", 30)

        result = subprocess.run(
            command,
            shell=True,
            capture_output=True,
            text=True,
            timeout=timeout
        )

        response = {
            "success": result.returncode == 0,
            "stdout": result.stdout,
            "stderr": result.stderr,
            "returncode": result.returncode,
            "location": EDGE_LOCATION
        }

        if msg.reply:
            await nc.publish(msg.reply, json.dumps(response).encode())

    except Exception as e:
        if msg.reply:
            await nc.publish(msg.reply, json.dumps({"error": str(e)}).encode())


async def handle_ha_command(msg):
    """Execute Home Assistant commands."""
    import json
    try:
        data = json.loads(msg.data.decode())
        service = data.get("service", "")
        entity_id = data.get("entity_id", "")

        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{HOMEASSISTANT_URL}/api/services/{service}",
                headers={"Authorization": f"Bearer {HOMEASSISTANT_TOKEN}"},
                json={"entity_id": entity_id}
            )

            if msg.reply:
                await nc.publish(msg.reply, json.dumps({
                    "success": response.status_code == 200,
                    "location": EDGE_LOCATION
                }).encode())

    except Exception as e:
        if msg.reply:
            await nc.publish(msg.reply, json.dumps({"error": str(e)}).encode())


# =============================================================================
# FastAPI App (Health Check)
# =============================================================================

@asynccontextmanager
async def lifespan(app: FastAPI):
    """Startup and shutdown events."""
    await connect_nats()
    yield
    if nc:
        await nc.close()


app = FastAPI(
    title=f"Remote Agent - {EDGE_LOCATION}",
    lifespan=lifespan
)


class HealthResponse(BaseModel):
    status: str
    location: str
    agent: str
    nats_connected: bool


@app.get("/health")
async def health() -> HealthResponse:
    """Health check endpoint."""
    return HealthResponse(
        status="healthy",
        location=EDGE_LOCATION,
        agent=AGENT_NAME,
        nats_connected=nc is not None and nc.is_connected if nc else False
    )


@app.get("/status")
async def status():
    """Detailed status endpoint."""
    return {
        "location": EDGE_LOCATION,
        "agent": AGENT_NAME,
        "nats": {
            "url": NATS_URL,
            "connected": nc.is_connected if nc else False
        },
        "integrations": {
            "homeassistant": HOMEASSISTANT_URL,
            "frigate": FRIGATE_URL,
            "ollama": OLLAMA_URL
        }
    }


# =============================================================================
# Main
# =============================================================================

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8090)
