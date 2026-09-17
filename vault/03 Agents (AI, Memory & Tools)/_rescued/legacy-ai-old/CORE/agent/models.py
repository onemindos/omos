"""
Model Selection and Configuration
Handles hybrid mode switching between online and offline models
"""

import asyncio
import aiohttp
from typing import Dict, Any
from enum import Enum


class AgentMode(str, Enum):
    """Operating modes for the agent"""
    ONLINE = "online"
    OFFLINE = "offline"
    AUTO = "auto"


class ModelSelector:
    """
    Manages model selection and mode switching.
    
    Modes:
    - online: Use cloud APIs (AssemblyAI, Cartesia, Azure)
    - offline: Use local models (Whisper.cpp, Ollama, Piper)
    - auto: Automatically detect and switch
    """
    
    def __init__(self, mode: str = "online"):
        self.mode = AgentMode(mode)
        self._is_online = True
        self._monitor_task = None
    
    async def check_connectivity(self) -> bool:
        """
        Check if internet connection is available.
        Tests connectivity to Azure OpenAI (or any reliable endpoint).
        """
        try:
            async with aiohttp.ClientSession() as session:
                async with session.get(
                    "https://api.openai.com/v1/models",
                    timeout=aiohttp.ClientTimeout(total=2)
                ) as resp:
                    return resp.status == 200
        except Exception as e:
            print(f"⚠️  Connectivity check failed: {e}")
            return False
    
    async def monitor_connectivity(self):
        """
        Continuously monitor internet connection.
        Automatically switches modes when connectivity changes.
        """
        print("🔄 Starting connectivity monitor (auto mode)...")
        
        while True:
            is_online = await self.check_connectivity()
            
            # Detect mode change
            if is_online and not self._is_online:
                print("🌐 Internet restored - switching to ONLINE mode")
                self._is_online = True
                self.mode = AgentMode.ONLINE
                await self._notify_mode_change("online")
            
            elif not is_online and self._is_online:
                print("📴 Internet lost - switching to OFFLINE mode")
                self._is_online = False
                self.mode = AgentMode.OFFLINE
                await self._notify_mode_change("offline")
            
            # Check every 10 seconds
            await asyncio.sleep(10)
    
    async def _notify_mode_change(self, new_mode: str):
        """
        Notify systems of mode change.
        In Phase 2B, this will trigger:
        - Model swapping (Whisper.cpp ↔ AssemblyAI)
        - Database switching (SQLite ↔ Supabase)
        - Tool availability updates
        """
        # TODO: Implement in Phase 2B
        print(f"🔔 Mode changed to: {new_mode}")
        pass


def get_model_config(mode: AgentMode) -> Dict[str, Any]:
    """
    Get model configuration based on current mode.
    
    Returns dict with model providers for STT, LLM, TTS, Vision.
    """
    
    if mode == AgentMode.ONLINE:
        return {
            "stt": {
                "provider": "assemblyai",
                "model": "best",  # Highest accuracy
                "realtime": True
            },
            "llm": {
                "provider": "letta",
                "base_url": "http://letta:8283",
                "mode": "cloud"  # Letta uses Supabase Cloud
            },
            "tts": {
                "provider": "cartesia",
                "model": "sonic-english",
                "voice": "default"  # Can be cloned voice ID
            },
            "vision": {
                "provider": "azure",
                "model": "gpt-4o",
                "enabled": True
            }
        }
    
    elif mode == AgentMode.OFFLINE:
        # Phase 2B configuration (Jetson Nano edge deployment)
        return {
            "stt": {
                "provider": "whisper-cpp",
                "model": "large-v3-int8",
                "path": "/models/whisper-large-int8"
            },
            "llm": {
                "provider": "letta",
                "base_url": "http://letta:8283",
                "mode": "local"  # Letta uses SQLite
            },
            "tts": {
                "provider": "piper",
                "voice": "en_US-lessac-medium",
                "path": "/models/piper"
            },
            "vision": {
                "provider": "llama",
                "model": "llama-3.2-1b-vision",
                "tensorrt": True,
                "enabled": True
            }
        }
    
    else:  # AUTO mode defaults to online config
        return get_model_config(AgentMode.ONLINE)


# ============================================================
# Helper Functions
# ============================================================

def get_latency_target(mode: AgentMode) -> int:
    """
    Get expected latency target in milliseconds.
    
    Returns:
        Target latency for voice response
    """
    if mode == AgentMode.ONLINE:
        return 100  # <100ms for online mode
    else:
        return 200  # <200ms for offline mode (still very responsive)
