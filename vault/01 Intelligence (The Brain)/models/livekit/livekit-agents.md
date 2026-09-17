# LiveKit Agents Framework

**Package:** `livekit-agents[openai,deepgram,cartesia,silero]`
**Python:** 3.11+
**Framework version:** 1.x (major rewrite from 0.x)

## Pipeline Components

| Component | Default | Alternatives |
|-----------|---------|-------------|
| STT | Deepgram Nova-3 | OpenAI Whisper, AssemblyAI |
| LLM | GPT-4o | Any OpenAI-compatible |
| TTS | Cartesia Sonic-2 | ElevenLabs, Deepgram Aura, OpenAI |
| VAD | Silero | Built-in |
| Turn detection | MultilingualModel | Manual |

## Minimal Agent Setup

```python
from livekit.agents import AutoSubscribe, JobContext, WorkerOptions, cli, llm
from livekit.agents.voice_assistant import VoiceAssistant
from livekit.plugins import cartesia, deepgram, openai, silero

async def entrypoint(ctx: JobContext):
    initial_ctx = llm.ChatContext().append(
        role="system",
        text="You are Legacy, Zeus's AI. Be concise, sharp, and direct."
    )
    await ctx.connect(auto_subscribe=AutoSubscribe.AUDIO_ONLY)

    assistant = VoiceAssistant(
        vad=silero.VAD.load(),
        stt=deepgram.STT(),
        llm=openai.LLM(model="gpt-4o"),
        tts=cartesia.TTS(),
        chat_ctx=initial_ctx,
    )
    assistant.start(ctx.room)
    await assistant.say("Hello Zeus, I'm ready.")

if __name__ == "__main__":
    cli.run_app(WorkerOptions(entrypoint_fnc=entrypoint))
```

## Key Concepts

- **Room:** WebRTC session where user and agent co-exist
- **Worker:** Process that hosts the agent — connects to LiveKit server
- **VAD (Voice Activity Detection):** Detects when user starts/stops speaking
- **Turn detection:** Knows when to interrupt vs wait

## LiveKit SIP

Connect real phone numbers to LiveKit rooms — agent answers phone calls.

```
User dials phone number → SIP trunk → LiveKit room → Agent picks up
```

Use case: Legacy agent accessible by phone, not just app.

## Deployment on VPS

Agent worker runs on the DO VPS, connects to LiveKit Cloud:

```bash
python agent.py start  # connects to LIVEKIT_URL, waits for room
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
