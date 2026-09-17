"""
Supabase Logging
MCP tools for conversation logging and analytics
"""

import os
from datetime import datetime
from typing import Dict, Any, Optional
from supabase import create_client, Client


SUPABASE_URL = os.getenv("SUPABASE_URL")
SUPABASE_KEY = os.getenv("SUPABASE_SERVICE_KEY")

# Initialize Supabase client
supabase: Client = create_client(SUPABASE_URL, SUPABASE_KEY)


async def log_conversation_message(
    session_id: str,
    role: str,
    content: str,
    metadata: Optional[Dict[str, Any]] = None
) -> bool:
    """
    Log a conversation message to Supabase.
    
    Args:
        session_id: LiveKit session/room ID
        role: Message role (user, assistant, system)
        content: Message content
        metadata: Optional metadata (latency, model, tokens, etc.)
    
    Returns:
        True if successful
    """
    try:
        data = {
            "session_id": session_id,
            "role": role,
            "content": content,
            "metadata": metadata or {},
            "created_at": datetime.utcnow().isoformat()
        }
        
        result = supabase.table("conversation_messages").insert(data).execute()
        return bool(result.data)
    
    except Exception as e:
        print(f"❌ Supabase logging error: {e}")
        return False


async def log_tool_execution(
    session_id: str,
    tool_name: str,
    inputs: Dict[str, Any],
    outputs: Any,
    success: bool,
    duration_ms: Optional[float] = None
) -> bool:
    """
    Log tool execution to Supabase for analytics.
    
    Args:
        session_id: LiveKit session/room ID
        tool_name: Name of the executed tool
        inputs: Tool input parameters
        outputs: Tool outputs/results
        success: Whether execution succeeded
        duration_ms: Execution duration in milliseconds
    
    Returns:
        True if successful
    """
    try:
        data = {
            "session_id": session_id,
            "tool_name": tool_name,
            "inputs": inputs,
            "outputs": outputs,
            "success": success,
            "duration_ms": duration_ms,
            "executed_at": datetime.utcnow().isoformat()
        }
        
        result = supabase.table("tool_executions").insert(data).execute()
        return bool(result.data)
    
    except Exception as e:
        print(f"❌ Supabase tool logging error: {e}")
        return False


async def get_conversation_history(
    session_id: str,
    limit: int = 50
) -> list:
    """
    Get conversation history from Supabase.
    
    Args:
        session_id: LiveKit session/room ID
        limit: Maximum messages to retrieve
    
    Returns:
        List of message objects
    """
    try:
        result = supabase.table("conversation_messages") \
            .select("*") \
            .eq("session_id", session_id) \
            .order("created_at", desc=True) \
            .limit(limit) \
            .execute()
        
        return result.data or []
    
    except Exception as e:
        print(f"❌ Supabase query error: {e}")
        return []


# ============================================================
# MCP Tool Definitions
# ============================================================

def register_supabase_tools() -> list:
    """Register Supabase logging tools in MCP protocol format"""
    return [
        {
            "name": "log_message",
            "description": "Log a conversation message for analytics and history",
            "input_schema": {
                "type": "object",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "Session identifier"
                    },
                    "role": {
                        "type": "string",
                        "description": "Message role",
                        "enum": ["user", "assistant", "system"]
                    },
                    "content": {
                        "type": "string",
                        "description": "Message content"
                    },
                    "metadata": {
                        "type": "object",
                        "description": "Optional metadata (latency, model, etc.)"
                    }
                },
                "required": ["session_id", "role", "content"]
            },
            "function": log_conversation_message
        },
        {
            "name": "log_tool_execution",
            "description": "Log tool execution for analytics",
            "input_schema": {
                "type": "object",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "Session identifier"
                    },
                    "tool_name": {
                        "type": "string",
                        "description": "Name of executed tool"
                    },
                    "inputs": {
                        "type": "object",
                        "description": "Tool input parameters"
                    },
                    "outputs": {
                        "description": "Tool outputs/results"
                    },
                    "success": {
                        "type": "boolean",
                        "description": "Whether execution succeeded"
                    },
                    "duration_ms": {
                        "type": "number",
                        "description": "Execution duration in milliseconds"
                    }
                },
                "required": ["session_id", "tool_name", "inputs", "success"]
            },
            "function": log_tool_execution
        },
        {
            "name": "get_conversation_history",
            "description": "Retrieve conversation history from database",
            "input_schema": {
                "type": "object",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "Session identifier"
                    },
                    "limit": {
                        "type": "integer",
                        "description": "Maximum messages to retrieve",
                        "default": 50
                    }
                },
                "required": ["session_id"]
            },
            "function": get_conversation_history
        }
    ]
