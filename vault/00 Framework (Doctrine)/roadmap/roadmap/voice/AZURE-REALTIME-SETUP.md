# Azure OpenAI GPT-4o Realtime Setup Guide

Complete guide for using Azure OpenAI GPT-4o Realtime API with your Agno agents.

## Features

- **Real-time voice conversations** with ~300ms latency
- **Vision support** - agents can see camera feed
- **Tool calling** - agents can use tools during conversation
- **Any Agno agent** - works with any agent in your registry
- **Local LiveKit** - full control, runs on your infrastructure

## Prerequisites

1. **Azure OpenAI Access** with GPT-4o Realtime deployment
2. **5k Azure credits** (you have this ✓)
3. **Local LiveKit server** running
4. **LiveKit agents SDK** installed

## Step 1: Azure OpenAI Setup

### Create GPT-4o Realtime Deployment

1. Go to [Azure Portal](https://portal.azure.com)
2. Navigate to your Azure OpenAI resource
3. Click **Model deployments** → **Deploy model**
4. Select **gpt-4o-realtime-preview**
5. Name it: `gpt-4o-realtime` (or your preference)
6. Deploy the model

### Get Your Credentials

```bash
# Your Azure OpenAI endpoint (find in Azure Portal)
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com

# Your API key (Keys and Endpoint section)
AZURE_OPENAI_API_KEY=your-azure-api-key

# Your deployment name (from step above)
AZURE_OPENAI_DEPLOYMENT_NAME=gpt-4o-realtime

# API version (use this for Realtime API)
AZURE_OPENAI_API_VERSION=2024-10-01-preview
```

## Step 2: Environment Variables

Add these to your `.env` file:

```bash
# ============================================================================
# Azure OpenAI Configuration
# ============================================================================
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
AZURE_OPENAI_API_KEY=your-azure-api-key
AZURE_OPENAI_DEPLOYMENT_NAME=gpt-4o-realtime
AZURE_OPENAI_API_VERSION=2024-10-01-preview

# ============================================================================
# LiveKit Configuration (Local)
# ============================================================================
LIVEKIT_URL=ws://localhost:7880
LIVEKIT_API_KEY=devkey
LIVEKIT_API_SECRET=secret

# ============================================================================
# Optional: Pipeline Mode Providers
# ============================================================================
# For STT → Agent → TTS mode (if you want fallback option)
DEEPGRAM_API_KEY=your-deepgram-key
CARTESIA_API_KEY=your-cartesia-key
```

## Step 3: Install LiveKit Server (Local)

### Option A: Docker (Recommended)

```bash
# Run LiveKit server locally
docker run -d \
  --name livekit \
  -p 7880:7880 \
  -p 7881:7881 \
  -p 7882:7882/udp \
  -e LIVEKIT_KEYS="devkey: secret" \
  livekit/livekit-server:latest \
  --dev
```

### Option B: Binary

```bash
# Download from https://github.com/livekit/livekit/releases
# Extract and run
./livekit-server --dev
```

Verify it's running:
```bash
curl http://localhost:7880
# Should return: LiveKit server is running
```

## Step 4: Install Python Dependencies

```bash
cd backend

# Install LiveKit agents SDK with OpenAI support
pip install livekit-agents livekit-plugins-openai livekit-plugins-silero

# Optional: For pipeline mode
pip install livekit-plugins-deepgram livekit-plugins-cartesia
```

## Step 5: Start the Realtime Worker

The worker listens for voice sessions and handles the Azure OpenAI integration:

```bash
cd backend

# Run the realtime worker
python -m agno.multimodal.realtime_worker dev

# You should see:
# INFO:multimodal.realtime_worker:Starting OneMind Real-Time Agent Worker
# INFO:multimodal.realtime_worker:Prewarming realtime worker...
# INFO:livekit:Worker connected to LiveKit
```

Keep this running in a terminal!

## Step 6: Test Voice with Any Agent

### Option A: API Test (Quick Start)

```bash
# Create a voice session with your legacy agent
curl -X POST http://localhost:8000/api/realtime/quick-start \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "legacy-mother-ai",
    "vision": true
  }'

# Response:
{
  "room_name": "rt-legacy-abc123",
  "token": "eyJ...",
  "ws_url": "ws://localhost:7880",
  "agent_id": "legacy-mother-ai",
  "provider": "openai",
  "vision_enabled": true
}
```

### Option B: Full Configuration

```bash
# Create session with specific settings
curl -X POST http://localhost:8000/api/realtime/session \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "your-agent-id",
    "provider": "openai",
    "enable_vision": true,
    "voice": "alloy",
    "user_id": "zeus"
  }'
```

### Option C: Python Client

```python
import asyncio
from backend.agno.multimodal.realtime import create_realtime_session, RealtimeProvider

async def test_voice():
    # Create session
    session = await create_realtime_session(
        agent_id="legacy-mother-ai",
        provider=RealtimeProvider.OPENAI,
        enable_vision=True,
        user_id="zeus"
    )

    print(f"Room: {session.room_name}")
    print(f"Token: {session.token}")
    print(f"Connect to: {session.ws_url}")

    # Your Flutter app would connect using these credentials

asyncio.run(test_voice())
```

## Step 7: Connect from Flutter

Your Flutter app connects to LiveKit using the credentials from Step 6:

```dart
import 'package:livekit_client/livekit_client.dart';

// Create voice session
final response = await api.post('/api/realtime/quick-start', {
  'agent_id': 'legacy-mother-ai',
  'vision': true,
});

// Connect to LiveKit room
final room = Room();
await room.connect(
  response['ws_url'],
  response['token'],
);

// Enable microphone
await room.localParticipant.setMicrophoneEnabled(true);

// Enable camera (for vision)
await room.localParticipant.setCameraEnabled(true);

// Listen for agent responses
room.on<TrackSubscribed>((event) {
  if (event.track.kind == TrackKind.AUDIO) {
    // Agent is speaking!
  }
});
```

## Available Voices

Azure OpenAI supports these voices:

- **alloy** - Neutral, balanced
- **echo** - Male, warm
- **shimmer** - Female, soft
- **ash** - Male, measured
- **ballad** - Storytelling
- **coral** - Female, friendly
- **sage** - Authoritative
- **verse** - Expressive

Specify voice when creating session:
```json
{
  "voice": "alloy"
}
```

## Vision Support

To enable vision (camera feed):

1. Set `enable_vision: true` when creating session
2. Enable camera in your Flutter app
3. Agent automatically receives camera frames
4. Agent can describe what it sees in real-time

The agent's instructions automatically include:
> "Vision is enabled - you can see the camera feed. If you see something, comment on it naturally."

## Using with Different Agents

Works with **any agent** in your registry:

```python
# Research agent with voice
session = await create_realtime_session(
    agent_id="research-agent",
    provider=RealtimeProvider.OPENAI,
)

# CEO agent with voice
session = await create_realtime_session(
    agent_id="ceo-agent",
    provider=RealtimeProvider.OPENAI,
)

# Custom agent with voice + vision
session = await create_realtime_session(
    agent_id="your-custom-agent",
    provider=RealtimeProvider.OPENAI,
    enable_vision=True,
)
```

The agent's system prompt and tools are automatically loaded from the registry!

## Health Check

Verify everything is configured:

```bash
curl http://localhost:8000/api/realtime/health

# Response:
{
  "status": "healthy",
  "active_sessions": 0,
  "providers_available": ["openai (azure)", "pipeline"]
}
```

You should see `"openai (azure)"` in the providers list.

## Troubleshooting

### Worker not connecting to Azure

**Check logs:**
```bash
# In worker terminal, you should see:
INFO:multimodal.realtime_worker:Starting Azure OpenAI Realtime session
INFO:multimodal.realtime_worker:Azure endpoint: https://your-resource.openai.azure.com, deployment: gpt-4o-realtime
```

**If not:**
- Verify `AZURE_OPENAI_ENDPOINT` is set
- Verify `AZURE_OPENAI_API_KEY` is correct
- Check deployment name matches Azure Portal

### Audio not working

1. Check LiveKit server is running: `curl http://localhost:7880`
2. Check worker is connected: Look for "Worker connected" in logs
3. Verify room was created: `curl http://localhost:8000/api/realtime/sessions`
4. Check Flutter app connected to room

### Vision not working

1. Verify `enable_vision: true` in session creation
2. Check camera enabled in Flutter app
3. Look for "Subscribing to video" in worker logs
4. Grant camera permissions in browser/app

### Azure credits usage

Monitor your usage in Azure Portal:
- Realtime API uses **audio input tokens** + **audio output tokens**
- Vision adds **image tokens** per frame sent
- Typical 10-minute conversation: ~$0.30-0.50

## Cost Optimization

Realtime API pricing (as of 2024):
- Audio input: $100/million tokens (~$0.06/min)
- Audio output: $200/million tokens (~$0.12/min)
- Text tokens: Standard GPT-4o rates

Tips to save credits:
1. **Disable vision** when not needed (`enable_vision: false`)
2. **Lower FPS** for vision: `vision_fps: 1` (default)
3. **Use pipeline mode** for non-critical conversations (uses standard API)
4. **Set session timeouts** - rooms auto-close after 5 min empty

## Next Steps

1. **Test with your agents** - Try different agent_ids
2. **Build UI** - Add voice button to Flutter app
3. **Add tools** - Agents can call tools during conversation
4. **Monitor usage** - Track Azure costs
5. **Scale up** - Move to cloud LiveKit when needed

## Architecture Flow

```
User speaks → LiveKit (local) → Azure OpenAI GPT-4o Realtime
                ↓
         Agent instructions from registry
                ↓
User hears ← LiveKit (local) ← Azure OpenAI GPT-4o Realtime
```

With vision:
```
User's camera → LiveKit → Azure OpenAI (sees video)
                  ↓
            Describes what it sees
                  ↓
User hears   ← LiveKit ← Azure OpenAI
```

## Support

If you hit issues:
1. Check worker logs for errors
2. Verify Azure Portal shows API calls
3. Test with `curl` before Flutter integration
4. Check [LiveKit docs](https://docs.livekit.io) for troubleshooting

Your system is now set up for real-time voice + vision with any Agno agent using your Azure credits!
