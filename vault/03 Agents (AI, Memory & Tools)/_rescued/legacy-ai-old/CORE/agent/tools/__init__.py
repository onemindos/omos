"""
Tool Registration
Registers all MCP tools with the Letta agent
"""

from typing import List, Dict, Any
from .clickup import register_clickup_tools
from .home_assistant import register_home_assistant_tools
from .supabase_logger import register_supabase_tools


def register_all_tools() -> List[Dict[str, Any]]:
    """
    Register all MCP tools for the Letta agent.
    
    Returns:
        List of tool definitions in MCP protocol format
    """
    tools = []
    
    # Task management (ClickUp)
    tools.extend(register_clickup_tools())
    
    # Smart home (Home Assistant)
    tools.extend(register_home_assistant_tools())
    
    # Data logging (Supabase)
    tools.extend(register_supabase_tools())
    
    print(f"✅ Registered {len(tools)} MCP tools")
    return tools
