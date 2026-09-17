"""
Home Assistant Integration
MCP tools for smart home control
"""

import os
import aiohttp
from typing import List, Dict, Any, Optional


HOME_ASSISTANT_URL = os.getenv("HOME_ASSISTANT_URL", "http://homeassistant.local:8123")
HOME_ASSISTANT_TOKEN = os.getenv("HOME_ASSISTANT_TOKEN")


async def get_home_assistant_states() -> List[Dict]:
    """
    Get all entity states from Home Assistant.
    
    Returns:
        List of entity state objects
    """
    headers = {
        "Authorization": f"Bearer {HOME_ASSISTANT_TOKEN}",
        "Content-Type": "application/json"
    }
    
    async with aiohttp.ClientSession() as session:
        async with session.get(
            f"{HOME_ASSISTANT_URL}/api/states",
            headers=headers
        ) as resp:
            if resp.status == 200:
                return await resp.json()
            else:
                print(f"❌ Home Assistant API error: {resp.status}")
                return []


async def call_home_assistant_service(
    domain: str,
    service: str,
    entity_id: Optional[str] = None,
    data: Optional[Dict] = None
) -> bool:
    """
    Call a Home Assistant service.
    
    Args:
        domain: Service domain (e.g., "light", "switch", "climate")
        service: Service name (e.g., "turn_on", "turn_off", "toggle")
        entity_id: Target entity ID (e.g., "light.bedroom")
        data: Additional service data
    
    Returns:
        True if successful
    """
    headers = {
        "Authorization": f"Bearer {HOME_ASSISTANT_TOKEN}",
        "Content-Type": "application/json"
    }
    
    payload = data or {}
    if entity_id:
        payload["entity_id"] = entity_id
    
    async with aiohttp.ClientSession() as session:
        async with session.post(
            f"{HOME_ASSISTANT_URL}/api/services/{domain}/{service}",
            headers=headers,
            json=payload
        ) as resp:
            return resp.status == 200


async def get_entity_state(entity_id: str) -> Optional[Dict]:
    """
    Get the state of a specific entity.
    
    Args:
        entity_id: Entity ID (e.g., "sensor.temperature")
    
    Returns:
        Entity state object or None
    """
    headers = {
        "Authorization": f"Bearer {HOME_ASSISTANT_TOKEN}",
        "Content-Type": "application/json"
    }
    
    async with aiohttp.ClientSession() as session:
        async with session.get(
            f"{HOME_ASSISTANT_URL}/api/states/{entity_id}",
            headers=headers
        ) as resp:
            if resp.status == 200:
                return await resp.json()
            else:
                return None


# ============================================================
# MCP Tool Definitions
# ============================================================

def register_home_assistant_tools() -> List[Dict[str, Any]]:
    """Register Home Assistant tools in MCP protocol format"""
    return [
        {
            "name": "get_smart_home_states",
            "description": "Get all smart home device states from Home Assistant",
            "input_schema": {
                "type": "object",
                "properties": {}
            },
            "function": get_home_assistant_states
        },
        {
            "name": "control_device",
            "description": "Control a smart home device (turn on/off, adjust, etc.)",
            "input_schema": {
                "type": "object",
                "properties": {
                    "domain": {
                        "type": "string",
                        "description": "Device domain (light, switch, climate, etc.)",
                        "enum": ["light", "switch", "climate", "cover", "fan", "lock"]
                    },
                    "service": {
                        "type": "string",
                        "description": "Service to call (turn_on, turn_off, toggle, etc.)"
                    },
                    "entity_id": {
                        "type": "string",
                        "description": "Entity ID (e.g., light.bedroom)"
                    },
                    "data": {
                        "type": "object",
                        "description": "Additional service data (brightness, color, temperature, etc.)"
                    }
                },
                "required": ["domain", "service"]
            },
            "function": call_home_assistant_service
        },
        {
            "name": "get_device_state",
            "description": "Get the current state of a specific smart home device",
            "input_schema": {
                "type": "object",
                "properties": {
                    "entity_id": {
                        "type": "string",
                        "description": "Entity ID (e.g., sensor.temperature)"
                    }
                },
                "required": ["entity_id"]
            },
            "function": get_entity_state
        }
    ]
