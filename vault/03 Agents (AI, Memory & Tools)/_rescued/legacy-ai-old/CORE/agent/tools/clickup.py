"""
ClickUp Integration
MCP tools for task and project management
"""

import os
import aiohttp
from typing import List, Dict, Any, Optional


CLICKUP_API_KEY = os.getenv("CLICKUP_API_KEY")
CLICKUP_TEAM_ID = os.getenv("CLICKUP_TEAM_ID")
CLICKUP_BASE_URL = "https://api.clickup.com/api/v2"


async def get_clickup_tasks(status: Optional[str] = None) -> List[Dict]:
    """
    Get user's tasks from ClickUp.
    
    Args:
        status: Filter by status (e.g., "to do", "in progress", "complete")
    
    Returns:
        List of task objects
    """
    headers = {"Authorization": CLICKUP_API_KEY}
    params = {}
    
    if status:
        params["statuses[]"] = status
    
    async with aiohttp.ClientSession() as session:
        async with session.get(
            f"{CLICKUP_BASE_URL}/team/{CLICKUP_TEAM_ID}/task",
            headers=headers,
            params=params
        ) as resp:
            if resp.status == 200:
                data = await resp.json()
                return data.get("tasks", [])
            else:
                print(f"❌ ClickUp API error: {resp.status}")
                return []


async def create_clickup_task(
    name: str,
    description: str = "",
    priority: int = 3,
    due_date: Optional[int] = None,
    list_id: Optional[str] = None
) -> Dict:
    """
    Create a new task in ClickUp.
    
    Args:
        name: Task name
        description: Task description
        priority: Priority level (1=Urgent, 2=High, 3=Normal, 4=Low)
        due_date: Unix timestamp (milliseconds)
        list_id: ClickUp list ID (defaults to first list in team)
    
    Returns:
        Created task object
    """
    headers = {
        "Authorization": CLICKUP_API_KEY,
        "Content-Type": "application/json"
    }
    
    # Get default list if not provided
    if not list_id:
        list_id = os.getenv("CLICKUP_DEFAULT_LIST_ID")
    
    payload = {
        "name": name,
        "description": description,
        "priority": priority
    }
    
    if due_date:
        payload["due_date"] = due_date
    
    async with aiohttp.ClientSession() as session:
        async with session.post(
            f"{CLICKUP_BASE_URL}/list/{list_id}/task",
            headers=headers,
            json=payload
        ) as resp:
            if resp.status == 200:
                return await resp.json()
            else:
                print(f"❌ ClickUp create task error: {resp.status}")
                return {}


async def update_task_status(task_id: str, status: str) -> bool:
    """
    Update task status in ClickUp.
    
    Args:
        task_id: ClickUp task ID
        status: New status (e.g., "in progress", "complete")
    
    Returns:
        True if successful
    """
    headers = {
        "Authorization": CLICKUP_API_KEY,
        "Content-Type": "application/json"
    }
    
    async with aiohttp.ClientSession() as session:
        async with session.put(
            f"{CLICKUP_BASE_URL}/task/{task_id}",
            headers=headers,
            json={"status": status}
        ) as resp:
            return resp.status == 200


# ============================================================
# MCP Tool Definitions
# ============================================================

def register_clickup_tools() -> List[Dict[str, Any]]:
    """Register ClickUp tools in MCP protocol format"""
    return [
        {
            "name": "get_tasks",
            "description": "Get user's tasks from ClickUp, optionally filtered by status",
            "input_schema": {
                "type": "object",
                "properties": {
                    "status": {
                        "type": "string",
                        "description": "Filter by status (to do, in progress, complete)",
                        "enum": ["to do", "in progress", "complete"]
                    }
                }
            },
            "function": get_clickup_tasks
        },
        {
            "name": "create_task",
            "description": "Create a new task in ClickUp",
            "input_schema": {
                "type": "object",
                "properties": {
                    "name": {
                        "type": "string",
                        "description": "Task name"
                    },
                    "description": {
                        "type": "string",
                        "description": "Task description"
                    },
                    "priority": {
                        "type": "integer",
                        "description": "Priority (1=Urgent, 2=High, 3=Normal, 4=Low)",
                        "default": 3
                    },
                    "due_date": {
                        "type": "integer",
                        "description": "Unix timestamp in milliseconds"
                    }
                },
                "required": ["name"]
            },
            "function": create_clickup_task
        },
        {
            "name": "update_task_status",
            "description": "Update the status of a task",
            "input_schema": {
                "type": "object",
                "properties": {
                    "task_id": {
                        "type": "string",
                        "description": "ClickUp task ID"
                    },
                    "status": {
                        "type": "string",
                        "description": "New status",
                        "enum": ["to do", "in progress", "complete"]
                    }
                },
                "required": ["task_id", "status"]
            },
            "function": update_task_status
        }
    ]
