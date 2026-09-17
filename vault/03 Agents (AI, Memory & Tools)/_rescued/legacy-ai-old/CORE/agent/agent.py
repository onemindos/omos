"""
Legacy AI - Phase 2 Agent
Real-time multimodal AI agent with hybrid cognitive architecture

This is the main entrypoint for the LiveKit agent worker.
It handles voice, vision, and tool execution with persistent memory via Letta.
"""

import os
import asyncio
from typing import Optional
from livekit import agents, rtc
from livekit.agents import AgentSession, llm
from livekit.plugins import assemblyai, cartesia, openai
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Import our custom modules
from models import ModelSelector, get_model_config
from tools import register_all_tools

# ============================================================
# Configuration
# ============================================================

LIVEKIT_URL = os.getenv("LIVEKIT_URL", "ws://livekit:7880")
LIVEKIT_API_KEY = os.getenv("LIVEKIT_API_KEY")
LIVEKIT_API_SECRET = os.getenv("LIVEKIT_API_SECRET")

LETTA_URL = os.getenv("LETTA_URL", "http://letta:8283")
LETTA_AGENT_ID = os.getenv("LETTA_AGENT_ID")

MODE = os.getenv("MODE", "online")  # online | offline | auto

# ============================================================
# Vision Support (Optional)
# ============================================================

async def get_latest_video_frame(room: rtc.Room) -> Optional[rtc.VideoFrame]:
    """
    Capture the latest frame from any video track in the room.
    Used for vision-enabled conversations.
    """
    for participant in room.remote_participants.values():
        for track_pub in participant.track_publications.values():
            if track_pub.track and isinstance(track_pub.track, rtc.RemoteVideoTrack):
                video_stream = rtc.VideoStream(track_pub.track)
                async for event in video_stream:
                    await video_stream.aclose()
                    return event.frame
    return None


async def before_llm_callback(agent_session: AgentSession, chat_ctx: llm.ChatContext):
    """
    Inject video frame into LLM context before processing.
    This enables vision-based queries like "What do you see?"
    """
    if not hasattr(agent_session, '_room'):
        return
    
    frame = await get_latest_video_frame(agent_session._room)
    if frame:
        # Add image to conversation context
        chat_ctx.messages.append(
            llm.ChatMessage(
                role="user",
                content=[llm.ChatImage(image=frame)]
            )
        )

# ============================================================
# Agent Class
# ============================================================

class LegacyAIAgent(agents.Agent):
    """
    Legacy AI - Cortana-inspired personal assistant
    
    Features:
    - Persistent memory across sessions (via Letta)
    - Real-time voice conversation (<100ms latency)
    - Vision understanding (camera/screen share)
    - Tool integration (MCP protocol)
    - Hybrid mode support (online ↔ offline)
    """
    
    def __init__(self):
        super().__init__(
            instructions=(
                "You are Legacy AI, a Cortana-inspired personal assistant. "
                "You have persistent memory and learn from every interaction. "
                "Be proactive, insightful, and conversational. "
                "Remember user preferences and adapt to their communication style. "
                "When using tools, explain what you're doing before executing. "
                "If you don't know something, say so honestly."
            )
        )
    
    async def on_chat_message(self, msg: rtc.ChatMessage):
        """Handle text chat messages (fallback for voice)"""
        print(f"💬 Text message from {msg.participant.identity}: {msg.message}")
        # Process via same LLM pipeline
        pass


# ============================================================
# Main Entrypoint
# ============================================================

async def entrypoint(ctx: agents.JobContext):
    """
    Main entrypoint for LiveKit agent worker.
    Called when a participant joins a room.
    """
    
    # Connect to LiveKit room
    await ctx.connect()
    print(f"✅ Connected to room: {ctx.room.name}")
    
    # Initialize model selector (handles hybrid mode)
    model_selector = ModelSelector(mode=MODE)
    
    # Get model configuration based on current mode
    config = get_model_config(model_selector.mode)
    print(f"🧠 Running in {model_selector.mode.upper()} mode")
    
    # Start connectivity monitor for auto mode
    if MODE == "auto":
        asyncio.create_task(model_selector.monitor_connectivity())
    
    # Configure AgentSession based on mode
    session = AgentSession(
        stt=assemblyai.STT(api_key=os.getenv("ASSEMBLYAI_API_KEY")),
        
        llm=openai.LLM(
            base_url=f"{LETTA_URL}/v1/chat/completions",
            model=f"agent-{LETTA_AGENT_ID}",
            api_key=os.getenv("LETTA_API_KEY", "not-needed-for-local")
        ),
        
        tts=cartesia.TTS(
            api_key=os.getenv("CARTESIA_API_KEY"),
            model="sonic-english",
            speed_alpha=1.0  # Normal speed
        )
    )
    
    # Start agent
    await session.start(
        agent=LegacyAIAgent(),
        room=ctx.room,
        before_llm_cb=before_llm_callback  # Enable vision
    )
    
    print(f"🤖 Agent started for participant: {ctx.participant.identity}")
    
    # Initial greeting
    await session.generate_reply(
        instructions=(
            "Greet the user warmly as Legacy AI. "
            "Mention that you have persistent memory and remember past conversations. "
            "Ask how you can help today."
        )
    )
    
    print("✅ Agent ready and listening...")


# ============================================================
# Worker Options
# ============================================================

if __name__ == "__main__":
    # Run the agent worker
    agents.cli.run_app(
        agents.WorkerOptions(
            entrypoint_fnc=entrypoint,
            api_key=LIVEKIT_API_KEY,
            api_secret=LIVEKIT_API_SECRET,
            ws_url=LIVEKIT_URL
        )
    )
