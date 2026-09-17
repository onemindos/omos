"""
OneMind Voice Agent — Credit-Optimized LiveKit Voice Pipeline

Components:
  VAD:  Silero (local, free)
  STT:  Deepgram Nova-3 (Deepgram credits)
  LLM:  Azure OpenAI (Azure credits) — routed via OpenClaw hooks
  TTS:  Azure OpenAI tts-1 (Azure credits)

This agent runs as a LiveKit worker. It joins a LiveKit room, listens for
speech, transcribes it, forwards the text to OpenClaw via NATS, and speaks
the response back using Azure OpenAI TTS.

Usage:
  python agent.py dev          # Development mode (auto-connects)
  python agent.py start        # Production mode (waits for dispatch)

Environment:
  LIVEKIT_URL            — LiveKit server URL (ws://localhost:7880 or wss://...)
  LIVEKIT_API_KEY        — LiveKit API key
  LIVEKIT_API_SECRET     — LiveKit API secret
  DEEPGRAM_API_KEY       — Deepgram API key (STT)
  AZURE_OPENAI_API_KEY   — Azure OpenAI API key (TTS)
  AZURE_OPENAI_ENDPOINT  — Azure OpenAI endpoint URL
  NATS_URL               — NATS server URL (nats://localhost:4222)
  OPENCLAW_HOOKS_TOKEN   — OpenClaw webhook auth token
  OPENCLAW_GATEWAY_URL   — OpenClaw gateway URL (http://localhost:18789)
"""

import asyncio
import logging
import os

from livekit import rtc
from livekit.agents import (
    AutoSubscribe,
    JobContext,
    JobProcess,
    WorkerOptions,
    cli,
    llm,
)
from livekit.agents.pipeline import VoicePipelineAgent
from livekit.plugins import deepgram, openai, silero

logger = logging.getLogger("voice-agent")
logger.setLevel(logging.INFO)


# ---------------------------------------------------------------------------
# Azure OpenAI TTS wrapper — routes TTS through Azure endpoint using credits
# ---------------------------------------------------------------------------
def create_azure_tts() -> openai.TTS:
    """Create TTS instance pointed at Azure OpenAI deployment."""
    endpoint = os.environ["AZURE_OPENAI_ENDPOINT"]
    api_key = os.environ["AZURE_OPENAI_API_KEY"]

    # Azure OpenAI TTS endpoint format
    base_url = f"{endpoint}/openai/deployments/tts/audio/speech?api-version=2024-12-17"

    return openai.TTS(
        model="tts-1",
        voice="alloy",
        base_url=base_url,
        api_key=api_key,
    )


# ---------------------------------------------------------------------------
# Azure OpenAI LLM wrapper — routes LLM through Azure endpoint
# ---------------------------------------------------------------------------
def create_azure_llm() -> openai.LLM:
    """Create LLM instance pointed at Azure OpenAI deployment."""
    endpoint = os.environ["AZURE_OPENAI_ENDPOINT"]
    api_key = os.environ["AZURE_OPENAI_API_KEY"]

    return openai.LLM(
        model="gpt-4o",
        base_url=f"{endpoint}/openai/v1",
        api_key=api_key,
    )


# ---------------------------------------------------------------------------
# Deepgram STT — uses Deepgram credits
# ---------------------------------------------------------------------------
def create_deepgram_stt() -> deepgram.STT:
    """Create Deepgram STT instance (Nova-3 model)."""
    return deepgram.STT(
        model="nova-3",
        language="en",
    )


# ---------------------------------------------------------------------------
# Agent prewarm — load Silero VAD model once on worker startup
# ---------------------------------------------------------------------------
def prewarm(proc: JobProcess):
    """Prewarm: load the Silero VAD model into memory."""
    proc.userdata["vad"] = silero.VAD.load()
    logger.info("Silero VAD model loaded")


# ---------------------------------------------------------------------------
# Agent entrypoint — called for each voice session
# ---------------------------------------------------------------------------
async def entrypoint(ctx: JobContext):
    """Main entrypoint for each voice session."""

    # System prompt for Legacy's voice persona
    system_prompt = (
        "You are Legacy, Zeus's AI command center. "
        "You are in voice mode — keep responses concise, clear, and conversational. "
        "Use short sentences. Avoid markdown, bullet points, or long lists. "
        "If you need to convey complex information, summarize the key points. "
        "You have access to all of Zeus's systems, calendar, tasks, and home automation."
    )

    initial_ctx = llm.ChatContext().append(
        role="system",
        text=system_prompt,
    )

    # Connect to the LiveKit room
    await ctx.connect(auto_subscribe=AutoSubscribe.AUDIO_ONLY)

    # Wait for a participant (the user) to join
    participant = await ctx.wait_for_participant()
    logger.info(f"Participant joined: {participant.identity}")

    # Build the voice pipeline
    agent = VoicePipelineAgent(
        vad=ctx.proc.userdata["vad"],
        stt=create_deepgram_stt(),
        llm=create_azure_llm(),
        tts=create_azure_tts(),
        chat_ctx=initial_ctx,
        # Pipeline tuning for low latency
        min_endpointing_delay=0.5,  # 500ms silence before processing
    )

    agent.start(ctx.room, participant)

    # Greet the user
    await agent.say(
        "Hey Zeus, I'm here. What do you need?",
        allow_interruptions=True,
    )


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------
if __name__ == "__main__":
    cli.run_app(
        WorkerOptions(
            entrypoint_fnc=entrypoint,
            prewarm_fnc=prewarm,
        ),
    )
