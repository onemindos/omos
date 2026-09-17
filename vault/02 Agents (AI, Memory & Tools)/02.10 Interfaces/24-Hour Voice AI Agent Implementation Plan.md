---
id: ui-02-20260111053925
title: "24-Hour Voice AI Agent Implementation Plan"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "24-Hour Voice AI Agent Implementation Plan"
---

# 24-Hour Voice AI Agent Implementation Plan
## Letta Cloud + LiveKit Cloud with Azure OpenAI

**This complete guide takes you from zero to a working voice-enabled AI agent system in under 24 hours, optimized for Azure credits and designed for non-technical users.**

---

## Hour-by-Hour Implementation Timeline

### Hours 1-3: Foundation Setup

**Hour 1: Account Creation & API Keys**
- Create LiveKit Cloud account at https://cloud.livekit.io (10 min)
  - Navigate to Settings → API Keys
  - Copy: WebSocket URL, API Key, API Secret
- Create Letta Cloud account at https://app.letta.com (10 min)
  - Go to Settings → API Keys → Create API Key
  - Save your `LETTA_API_KEY` securely
- Sign up for Deepgram at https://console.deepgram.com (15 min)
  - Get free $200 credit for speech-to-text
  - Copy API key from dashboard
- Sign up for Cartesia at https://play.cartesia.ai (15 min)
  - Get API key for text-to-speech
- Set up Azure OpenAI resource (10 min)
  - Go to portal.azure.com → Create "Azure OpenAI"
  - Choose East US 2 or Sweden Central region
  - Deploy `gpt-4o-mini-realtime-preview` model

**Hour 2: Development Environment Setup**
- Install Python 3.9+ if not already installed
- Install Git and create GitHub account (if new to Git)
- Clone the official Letta-LiveKit template:
```bash
git clone https://github.com/letta-ai/letta-voice.git
cd letta-voice
```
- Create Python virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate
pip install -r requirements.txt
```
- Create `.env` file and add all API keys:
```env
LIVEKIT_URL=wss://your-project.livekit.cloud
LIVEKIT_API_KEY=your_key
LIVEKIT_API_SECRET=your_secret
LETTA_API_KEY=letta_sk_...
AZURE_OPENAI_ENDPOINT=https://[name].openai.azure.com/
AZURE_OPENAI_API_KEY=your_azure_key
DEEPGRAM_API_KEY=your_deepgram_key
CARTESIA_API_KEY=your_cartesia_key
```

**Hour 3: Configure Your First Letta Agent**
- Open Letta Agent Development Environment at https://app.letta.com
- Click "Create Agent"
- Select architecture: `voice_convo_agent` (optimized for voice)
- Configure memory blocks:
  - **Persona**: "You are a helpful AI assistant with perfect memory of our conversations"
  - **Human**: "User is exploring voice AI technology"
- Enable sleep-time agent for better memory management
- Set `max_message_buffer_length: 20` for good latency
- Choose model: **gpt-4o-mini** (fast and cheap for testing)
- Attach built-in tools: `send_message`, `archival_memory_search`
- Copy your agent ID (format: `agent-xxxxx`)

### Hours 4-8: Basic Voice Agent Deployment

**Hour 4: Local Testing (Console Mode)**
Test your setup without LiveKit first:
```bash
# Still in letta-voice directory
python main.py console
```
This launches a terminal-based voice interface. Speak to test:
- STT working? (Deepgram transcribing your voice)
- LLM responding? (Letta agent answering)
- TTS working? (Cartesia speaking responses)
- Memory persisting? (Ask follow-up questions)

**Hour 5: Configure Azure OpenAI in Code**
Update `main.py` to use Azure OpenAI Realtime API:
```python
from livekit.plugins import azure

session = AgentSession(
    llm=azure.RealtimeModel(
        azure_endpoint=os.getenv("AZURE_OPENAI_ENDPOINT"),
        api_key=os.getenv("AZURE_OPENAI_API_KEY"),
        api_version="2024-10-01-preview",
        model="gpt-4o-mini-realtime",
        voice="alloy",
    )
)
```
This ensures your Azure credits are used for the voice model.

**Hour 6: Test with LiveKit Playground**
Run in dev mode to connect to LiveKit Cloud:
```bash
python main.py dev
```
- Agent starts and registers with LiveKit
- Open https://agents-playground.livekit.io
- Select your project → Click "Connect"
- **Speak to your agent through the web interface**
- Test interruptions, multi-turn conversations
- Verify memory: "What did I just ask you?"

**Hour 7: Add Custom Keys to Letta Cloud**
To ensure all costs bill to your Azure account:
- Go to https://app.letta.com/models
- Click "Add Custom Provider" → "Azure OpenAI"
- Enter your Azure credentials:
  - API Key, Endpoint URL, Deployment Name, API Version
- In your agent settings, select model from "Your models" section
- **Result**: Letta uses your Azure OpenAI, bypassing Letta's usage limits

**Hour 8: First Deployment to LiveKit Cloud**
Deploy to production using LiveKit CLI:
```bash
# Install LiveKit CLI (if not done)
# Mac: brew install livekit-cli
# Windows: winget install LiveKit.LiveKitCLI

# Authenticate
lk cloud auth

# Deploy agent
lk agent create
```
LiveKit auto-generates a Dockerfile and deploys. Your agent is now running 24/7 in the cloud!

### Hours 9-16: Letta Integration & Memory Setup

**Hour 9: Advanced Memory Configuration**
Return to Letta ADE and enhance your agent:
- Add custom memory blocks:
  - **Project Context**: "This is a voice AI system for [your use case]"
  - **User Preferences**: Store user-specific settings
- Configure memory management:
  - Enable `enable_sleeptime: true` for background memory processing
  - Set sleep-time agent to `claude-sonnet-4` (better quality for summarization)
  - Adjust buffer: `min_message_buffer_length: 5`

**Hour 10: Tool Calling Setup**
Add functional capabilities to your agent:
- In Letta ADE, navigate to Tools section
- Attach built-in tools:
  - `web_search`: Real-time web searching
  - `run_code`: Execute Python code
  - `send_message`: Agent-to-agent communication
  - `archival_memory_insert`: Long-term storage
- Test each tool in the ADE simulator:
  - "Search the web for latest AI news"
  - "Calculate 15% of $250"

**Hour 11: Implement Cached Audio Tokens**
Optimize Azure OpenAI costs by 90%:
```python
# In main.py, configure session with cached context
session = AgentSession(
    llm=azure.RealtimeModel(
        # ... existing config
    ),
    chat_ctx=ChatContext(
        messages=[
            {
                "role": "system",
                "content": "You are a helpful assistant...",  # This gets cached
                "cache_control": {"type": "ephemeral"}
            }
        ]
    )
)
```
Repeated context now costs $20/1M tokens instead of $100/1M.

**Hour 12: Multi-Turn Conversation Testing**
Create a test script to validate memory:
```python
# test_memory.py
async def test_conversation():
    # Conversation 1
    response1 = await agent.send_message("My favorite color is blue")
    print(response1)
    
    # Conversation 2 (5 min later)
    response2 = await agent.send_message("What's my favorite color?")
    assert "blue" in response2.lower()
    print("✅ Memory test passed!")
```

**Hours 13-16: Context Window Management**
Implement smart context compression to reduce costs:

**Hour 13**: Set up automatic summarization
- In Letta agent config, enable conversation summarization
- Trigger: Every 20 messages or when context reaches 80% of limit
- Sleeptime agent creates summary, main agent continues with compressed history

**Hour 14**: Configure archival memory
- Important facts get stored in `archival_memory` (unlimited storage)
- Configure retrieval: Top-5 most relevant memories injected per query
- Test: Share 10 facts, verify agent retrieves correct ones later

**Hour 15**: Optimize for latency
- Test current latency in playground (measure time to first word)
- Target: Under 1.5 seconds customer-perceived latency
- If too slow:
  - Switch to `gpt-4o-mini-realtime` (faster, 70% cheaper)
  - Reduce `max_message_buffer_length` to 15
  - Enable LiveKit's server-side VAD for instant turn detection

**Hour 16**: Create test conversations dataset
Document 10 test scenarios:
1. Simple greeting
2. Complex multi-turn reasoning
3. Tool use (web search)
4. Memory recall from previous session
5. Interruption handling
6. Long monologue response
7. Clarification questions
8. Error recovery
9. Context switching
10. Emotional understanding

### Hours 17-20: Testing & Refinement

**Hour 17: Automated Testing**
Create evaluation suite:
```bash
# In letta-voice directory
mkdir tests
```

Create `tests/test_agent.py`:
```python
import pytest
from livekit_client import Letta

@pytest.mark.asyncio
async def test_basic_response():
    client = Letta(token=os.getenv("LETTA_API_KEY"))
    response = await client.agents.messages.create(
        agent_id="agent-your-id",
        messages=[{"role": "user", "content": "Hello"}]
    )
    assert response.messages[0].content
    assert len(response.messages[0].content) > 0

@pytest.mark.asyncio
async def test_memory_persistence():
    # Test memory across sessions
    # ... (implement based on test_memory.py from Hour 12)
```

Run tests: `pytest tests/`

**Hour 18: Cost Monitoring Setup**
Configure Azure cost tracking:
- Go to portal.azure.com → Cost Management
- Create budget: $100/month with alerts at 50%, 75%, 90%
- Add cost tags to OpenAI resource: `project:voice-agent`, `env:production`
- Set up daily cost emails
- Check current spend:
  - Azure OpenAI usage in tokens/minute
  - LiveKit agent session minutes
  - Total daily burn rate

Create monitoring dashboard:
```python
# monitor_costs.py
from azure.monitor.query import MetricsQueryClient

def check_daily_spend():
    # Query Azure OpenAI token usage
    # Calculate: tokens × price per token
    # Alert if exceeding budget
```

**Hour 19: Load Testing**
Simulate concurrent users:
```bash
# Install testing tool
pip install locust

# Create locustfile.py
from locust import User, task, between
import asyncio

class VoiceAgentUser(User):
    wait_time = between(5, 15)
    
    @task
    def have_conversation(self):
        # Simulate 3-minute conversation
        # Measure: latency, error rate, memory accuracy
```

Run: `locust -f locustfile.py --users 10 --spawn-rate 2`

**Hour 20: Quality Assurance Checklist**
Test all critical paths:
- ✅ Agent joins room within 3 seconds
- ✅ Responds to simple queries under 1.5s
- ✅ Handles interruptions gracefully
- ✅ Remembers facts from 10+ messages ago
- ✅ Tools execute successfully (web search, calculations)
- ✅ No crashes after 30-minute conversation
- ✅ Audio quality is clear and natural
- ✅ Agent doesn't hallucinate API keys or sensitive data
- ✅ Costs stay within budget ($0.15/minute average)
- ✅ Logs are readable and helpful for debugging

### Hours 21-24: Multi-Agent Foundation

**Hour 21: Multi-Agent Architecture Design**
Plan your multi-agent system:
- **Supervisor Agent**: Routes queries to specialist agents
- **Research Agent**: Handles web search and data gathering
- **Memory Agent**: Manages long-term knowledge base
- **Action Agent**: Executes tasks (API calls, calculations)

Create agents in Letta ADE:
```bash
# Research Agent
- Persona: "You are a research specialist focused on gathering accurate information"
- Tools: web_search, archival_memory_insert
- Tags: ["worker", "research"]

# Memory Agent  
- Persona: "You manage long-term knowledge storage"
- Tools: archival_memory_search, archival_memory_insert
- Tags: ["worker", "memory"]

# Supervisor Agent
- Persona: "You coordinate between specialist agents"
- Tools: send_message_to_agent_async, send_message_to_agents_matching_all_tags
- Tags: ["supervisor"]
```

**Hour 22: Implement Agent-to-Agent Communication**
Enable inter-agent messaging:

In Letta ADE, attach communication tools to supervisor:
```python
# Supervisor can now route tasks
"Search for latest AI news" 
→ Supervisor: send_message_to_agent_async(agent_id="research-agent-id", ...)
→ Research Agent performs web search
→ Supervisor receives results, responds to user
```

Test flow:
1. User asks complex question
2. Supervisor analyzes and routes to research agent
3. Research agent gathers data
4. Memory agent stores important facts
5. Supervisor synthesizes and responds
6. User receives comprehensive answer with sources

**Hour 23: Shared Memory Configuration**
Create shared knowledge base:
- In Letta, create memory block: `label: "shared_knowledge"`
- Attach to all agents (research, memory, supervisor)
- Any agent can read/write to this shared context
- Test: Research agent adds fact → Supervisor can reference it immediately

**Hour 24: Multi-Agent Testing & Deployment**
Deploy multi-agent system:
```bash
# Deploy each agent separately to LiveKit
cd supervisor-agent/
lk agent create

cd ../research-agent/
lk agent create

cd ../memory-agent/
lk agent create
```

Configure routing in LiveKit:
- Single room receives user input
- Supervisor agent active by default
- Supervisor spawns worker agents as needed
- Parallel processing for independent tasks

Final validation:
- Test complex query requiring multiple agents
- Measure end-to-end latency (should be under 5s)
- Verify all agents communicate successfully
- Check costs: Multi-agent adds ~$0.05/minute overhead

---

## Architecture Decision Matrix

### Decision 1: GitHub vs Cloud-Only Management

**Recommended: Hybrid Model**

| What to Manage | Where | Rationale |
|----------------|-------|-----------|
| **Agent code** | GitHub | Version control, collaboration, rollback capability |
| **Dockerfile, requirements.txt** | GitHub | Infrastructure as code, reproducible deployments |
| **Environment variable template** | GitHub (`.env.example`) | Onboarding new developers |
| **Actual API keys** | Local `.env` (never commit) | Security |
| **Agent memory/personality** | Letta Cloud Dashboard (ADE) | Real-time editing without redeployment |
| **Tool configurations** | Letta Cloud Dashboard | Non-developer can adjust without coding |
| **Monitoring dashboards** | LiveKit + Letta Cloud | Real-time observability |
| **Deployment configs** | `livekit.toml` in GitHub | Track deployment settings |

**Pros:**
- ✅ Non-developers can edit agent behavior in ADE without touching code
- ✅ Developers get full version control for infrastructure
- ✅ Clear separation: code vs configuration
- ✅ Easy rollback: code (Git), config (Letta versions)

**Workflow:**
1. Developer commits code to GitHub
2. CI/CD deploys to LiveKit Cloud via `lk agent deploy`
3. Non-developer edits personality/memory in Letta ADE
4. Changes take effect immediately without redeployment

### Decision 2: Python vs Node.js

**Recommended: Python**

| Factor | Python | Node.js |
|--------|--------|---------|
| **Letta integration** | ✅ Official, well-documented | ⚠️ Limited examples |
| **LiveKit examples** | ✅ 104+ examples | ⚠️ Fewer resources |
| **Learning curve** | ✅ Easier for non-devs | ⚠️ Steeper |
| **Code length** | ✅ ~50 lines for working agent | ⚠️ More verbose |
| **Plugin ecosystem** | ✅ More STT/TTS/LLM plugins | ⚠️ Growing |

**Use Node.js only if:**
- Building full-stack web app with Next.js frontend
- Team is exclusively JavaScript-experienced
- Need tight integration with existing Node.js services

**Python wins for:**
- Fastest path to working voice agent
- Best documentation and community support
- Official Letta-LiveKit template is Python-only

### Decision 3: Which Realtime Model to Use

**For Azure Credits: Azure OpenAI Realtime API (gpt-4o-mini-realtime)**

| Model | Cost/Min | Latency | Azure Credits? | Recommendation |
|-------|----------|---------|----------------|----------------|
| **Azure OpenAI GPT-4o Realtime** | $0.30 | 1.18s | ✅ YES | Development/complex queries |
| **Azure OpenAI GPT-4o-Mini Realtime** | $0.08 | ~1.2s | ✅ YES | **BEST for Azure credits** |
| Amazon Nova Sonic | $0.06 | 1.09s | ❌ NO | Production alternative after credits exhausted |
| Gemini 2.5 Flash Live | $0.18 | 1.41s | ❌ NO | If need 1M context window |

**Strategy:**
1. **Development (Weeks 1-4)**: Use **gpt-4o-mini-realtime** to maximize Azure credit usage
2. **Production (Months 2+)**: Evaluate cost vs performance
   - If Azure credits remain: Continue with Mini
   - If budget-constrained: Migrate to Nova Sonic (80% cheaper)
   - If need advanced reasoning: Upgrade to full GPT-4o Realtime

**Cost optimization:**
- Use cached audio tokens (90% savings on repeated context)
- Implement context summarization (reduce token usage)
- Smart routing: Mini for simple queries, full GPT-4o for complex

### Decision 4: Virtual Avatar Plugin Recommendation

**Recommended: Tavus (for plug-and-play digital twin)**

| Plugin | Setup | Cost | Digital Twin Quality | Best For |
|--------|-------|------|----------------------|----------|
| **Tavus** | ⭐⭐⭐⭐⭐ Easiest | $$ Medium | ⭐⭐⭐⭐⭐ Excellent | **Professional digital twins** |
| Hedra | ⭐⭐⭐⭐ Easy | $ Low | ⭐⭐⭐⭐ Good | Creative/animated characters |
| Beyond Presence | ⭐⭐⭐⭐ Easy | $$$ High | ⭐⭐⭐⭐⭐ Hyper-realistic | Enterprise, ultra-realistic needs |
| Simli | ⭐⭐⭐⭐ Easy | $$ Medium | ⭐⭐⭐⭐ Good | Low-latency requirements |

**Tavus Setup (< 30 minutes):**
1. Sign up at https://tavus.io
2. Create replica in Tavus Studio:
   - Upload 2-minute video of yourself speaking
   - AI generates your digital twin
   - Get `replica_id`
3. Create persona (personality/voice characteristics)
4. Install plugin:
```bash
pip install livekit-plugins-tavus
```
5. Add to agent code:
```python
from livekit.plugins import tavus

avatar = tavus.AvatarSession(
    replica_id="r4c41453d2",
    persona_id="p2fbd605"
)
await avatar.start(session, room=ctx.room)
```

**Result**: Your digital twin appears in the LiveKit room, lip-synced to agent's voice.

**When to integrate avatar:**
- ✅ After basic voice agent works (Hour 24+)
- ✅ When frontend is ready (web/mobile interface)
- ✅ User needs visual engagement (education, healthcare, sales)

---

## Step-by-Step Setup Guide

### Pre-Flight Checklist

**System Requirements:**
- Operating System: Windows 10+, macOS 10.15+, or Linux
- Python: 3.9 or higher (`python --version`)
- RAM: 8GB minimum, 16GB recommended
- Internet: Stable connection (voice streaming requires low latency)
- Terminal/Command Prompt: Basic familiarity helpful

**Skills Needed:**
- ⚠️ Minimal: Copy-paste commands, edit text files
- ✅ No coding required for initial setup
- ✅ No server management needed (fully cloud-based)

### Part 1: LiveKit Cloud Setup

**Step 1.1: Create Account**
1. Visit https://cloud.livekit.io
2. Click "Sign Up" → Use Google/GitHub or email
3. Verify email
4. You're in! Free tier includes:
   - 500 agent session minutes/month
   - 5GB bandwidth
   - Access to playground

**Step 1.2: Create Project**
1. Click "Projects" in bottom left
2. "Create New Project" button
3. Name: `voice-agent-production`
4. Click "Create"

**Step 1.3: Get API Credentials**
1. Select your project
2. Navigate to "Settings" → "Keys"
3. Copy and save (use password manager):
   - **WebSocket URL**: `wss://voice-agent-production-xxxxx.livekit.cloud`
   - **API Key**: `APIxxxxxxxxxx`
   - **API Secret**: `secretxxxxxxxxxx`

**Step 1.4: Install LiveKit CLI** (for deployment)
```bash
# macOS
brew install livekit-cli

# Windows
winget install --id=LiveKit.LiveKitCLI -e

# Linux
curl -sSL https://get.livekit.io | bash
```

Verify: `lk version` (should show version number)

**Step 1.5: Authenticate CLI**
```bash
lk cloud auth
```
This opens browser → log in → CLI is now connected to your account.

### Part 2: Letta Cloud Setup

**Step 2.1: Create Account**
1. Visit https://app.letta.com
2. Sign up with Google/GitHub
3. Complete onboarding survey (select "Building voice agents")

**Step 2.2: Get API Key**
1. Click your profile → "Settings"
2. Navigate to "API Keys"
3. Click "Create API Key"
4. Name: `voice-agent-dev`
5. Copy key: `letta_sk_xxxxxxxx` (save securely)

**Step 2.3: Create Your First Agent**
1. Click "Create Agent" (big blue button)
2. **Agent Name**: "VoiceAssistant-v1"
3. **Architecture**: Select `voice_convo_agent` (optimized for voice)
4. **Memory Configuration**:
   - **Core Memory - Persona**:
     ```
     You are a helpful AI assistant named Alex.
     You have perfect memory of all our conversations.
     You are friendly, concise, and proactive.
     You can search the web and help with tasks.
     ```
   - **Core Memory - Human**:
     ```
     User is exploring voice AI technology.
     User prefers direct answers without excessive politeness.
     ```
5. **Model Selection**:
   - **Main Agent**: `gpt-4o-mini` (fast, cheap)
   - **Sleep-time Agent**: Enable, select `claude-sonnet-4` (better for memory processing)
6. **Message Buffer**:
   - `max_message_buffer_length`: 20
   - `min_message_buffer_length`: 5
7. **Tools**: Attach these built-in tools:
   - ✅ `send_message` (for multi-agent later)
   - ✅ `web_search`
   - ✅ `archival_memory_search`
   - ✅ `archival_memory_insert`
8. Click "Create Agent"
9. **Copy Agent ID** from URL: `https://app.letta.com/agents/agent-xxxxxxxx`

**Step 2.4: Test in ADE Simulator**
1. In the chat interface on right side, type: "Hello, what can you help me with?"
2. Agent responds (may take 5-10 seconds first time)
3. Test memory: "Remember my favorite color is blue"
4. Ask: "What's my favorite color?" → Should recall "blue"
5. Test tool: "Search the web for latest AI breakthroughs"
6. Verify web search tool executes and returns results

### Part 3: Azure OpenAI Setup

**Step 3.1: Create Azure OpenAI Resource**
1. Go to https://portal.azure.com
2. Search bar: Type "Azure OpenAI"
3. Click "Create"
4. Fill in:
   - **Subscription**: Select your subscription with credits
   - **Resource Group**: Create new → Name: `voice-ai-rg`
   - **Region**: **East US 2** (has Realtime API)
   - **Name**: `voice-agent-openai` (must be globally unique)
   - **Pricing Tier**: Standard S0
5. **Networking**: Select "All networks" (simplest for development)
6. Click "Review + Create" → "Create"
7. Wait 2-3 minutes for deployment

**Step 3.2: Deploy Realtime Model**
1. After deployment, click "Go to resource"
2. Click "Go to Azure AI Foundry portal" button
3. In AI Foundry:
   - Left sidebar → "Models + endpoints"
   - Click "+ Deploy model" → "Deploy base model"
   - Search: `gpt-4o-mini-realtime-preview`
   - Click "Confirm"
4. Configuration:
   - **Deployment name**: `voice-realtime-mini`
   - **Model version**: Latest (2024-12-17 or newer)
   - **Deployment type**: **Global Standard** (best availability)
   - **Tokens per minute rate limit**: 50K (default)
5. Click "Deploy" (takes 1-2 minutes)

**Step 3.3: Get Credentials**
1. Return to Azure Portal → Your OpenAI resource
2. Left sidebar → "Keys and Endpoint"
3. Copy:
   - **KEY 1**: (this is your API key)
   - **Endpoint**: `https://voice-agent-openai.openai.azure.com/`
   - **Location**: (e.g., eastus2)

**Step 3.4: Add Custom Keys to Letta Cloud**
1. Go to https://app.letta.com/models
2. Click "Add Custom Provider"
3. Select "Azure OpenAI"
4. Fill in:
   - **Display Name**: My Azure OpenAI
   - **API Key**: [paste from Step 3.3]
   - **Endpoint**: [paste from Step 3.3]
   - **Deployment Name**: `voice-realtime-mini`
   - **API Version**: `2024-10-01-preview`
5. Click "Test Connection" (should show ✅ Success)
6. Click "Save"
7. Return to your agent → Edit model → Under "Your models" dropdown, select "Azure OpenAI gpt-4o-mini-realtime"

**Result**: All Letta Cloud usage now bills directly to your Azure account (using your credits!).

### Part 4: Speech Provider Setup

**Step 4.1: Deepgram (Speech-to-Text)**
1. Visit https://console.deepgram.com
2. Sign up (free $200 credit)
3. After login, click "Create API Key"
4. Name: `livekit-voice-agent`
5. Copy API key: `xxxxxxxxxxxxxxxx`

**Step 4.2: Cartesia (Text-to-Speech)**
1. Visit https://play.cartesia.ai
2. Sign up with Google/email
3. Dashboard → "API Keys"
4. Click "Create API Key"
5. Copy key: `xxxxxxxxxxxxxxxx`

### Part 5: Project Setup

**Step 5.1: Install Python (if needed)**
```bash
# Check if installed
python --version
# Should show Python 3.9+ 

# If not installed:
# Windows: Download from python.org
# macOS: brew install python@3.11
# Linux: sudo apt install python3.11
```

**Step 5.2: Install Git (if needed)**
```bash
# Check if installed
git --version

# If not installed:
# Windows: Download from git-scm.com
# macOS: xcode-select --install
# Linux: sudo apt install git
```

**Step 5.3: Clone Template**
```bash
# Open terminal/command prompt
# Navigate to where you want the project:
cd Documents  # or wherever you keep projects

# Clone official Letta-LiveKit template
git clone https://github.com/letta-ai/letta-voice.git
cd letta-voice
```

**Step 5.4: Create Virtual Environment**
```bash
# Create isolated Python environment
python -m venv venv

# Activate it:
# macOS/Linux:
source venv/bin/activate

# Windows (Command Prompt):
venv\Scripts\activate.bat

# Windows (PowerShell):
venv\Scripts\Activate.ps1

# You should see (venv) prefix in terminal
```

**Step 5.5: Install Dependencies**
```bash
# Install all required Python packages
pip install -r requirements.txt

# This installs:
# - livekit-agents
# - livekit-plugins (OpenAI, Azure, Deepgram, Cartesia)
# - letta-client
# - python-dotenv
# Takes 2-3 minutes
```

**Step 5.6: Create Environment Variables File**
```bash
# Create .env file from template
cp .env.example .env

# Open .env in text editor:
# Windows: notepad .env
# macOS: open -e .env
# Linux: nano .env
```

**Paste this content** (replace with your actual keys from previous steps):
```env
# LiveKit Cloud (from Part 1, Step 1.3)
LIVEKIT_URL=wss://voice-agent-production-xxxxx.livekit.cloud
LIVEKIT_API_KEY=APIxxxxxxxxxx
LIVEKIT_API_SECRET=secretxxxxxxxxxx

# Letta Cloud (from Part 2, Step 2.2)
LETTA_API_KEY=letta_sk_xxxxxxxx

# Azure OpenAI (from Part 3, Step 3.3)
AZURE_OPENAI_ENDPOINT=https://voice-agent-openai.openai.azure.com/
AZURE_OPENAI_API_KEY=xxxxxxxxxxxxxxxx
AZURE_OPENAI_DEPLOYMENT=voice-realtime-mini
AZURE_OPENAI_API_VERSION=2024-10-01-preview

# Speech Providers (from Part 4)
DEEPGRAM_API_KEY=xxxxxxxxxxxxxxxx
CARTESIA_API_KEY=xxxxxxxxxxxxxxxx
```

Save and close the file.

**Step 5.7: Update Agent Code for Azure**
Open `main.py` in text editor and find the section where `AgentSession` is created (around line 20-30).

Replace with:
```python
from livekit.plugins import azure
from livekit.agents import AgentSession
import os

session = AgentSession(
    llm=azure.RealtimeModel(
        azure_endpoint=os.getenv("AZURE_OPENAI_ENDPOINT"),
        api_key=os.getenv("AZURE_OPENAI_API_KEY"),
        api_version=os.getenv("AZURE_OPENAI_API_VERSION"),
        model=os.getenv("AZURE_OPENAI_DEPLOYMENT"),
        voice="alloy",  # or "shimmer", "echo"
    )
)
```

### Part 6: Testing & Deployment

**Step 6.1: Local Console Test**
First, test without LiveKit (fastest debugging):
```bash
python main.py console
```

What happens:
1. Agent loads (5-10 seconds)
2. Prompt appears: "Ready to chat! Press Enter to start..."
3. Press Enter, then **speak into your microphone**
4. Agent transcribes, thinks, responds with voice
5. Test memory: Say something, then ask about it later
6. Press Ctrl+C to exit

**Troubleshooting:**
- ❌ "KeyError: LIVEKIT_URL" → Check `.env` file exists and is named correctly (not `.env.txt`)
- ❌ "No module named livekit" → Activate venv: `source venv/bin/activate`
- ❌ No audio input → Check microphone permissions in system settings
- ❌ "Letta API error" → Verify `LETTA_API_KEY` is correct

**Step 6.2: Test with LiveKit Playground**
Now test the full cloud integration:
```bash
python main.py dev
```

Terminal shows:
```
Connecting to LiveKit...
Agent registered successfully
Waiting for room connection...
```

**In your browser:**
1. Go to https://agents-playground.livekit.io
2. You should see your project listed
3. Click on your project
4. Click "Connect" button
5. **Allow microphone access** when prompted
6. You see audio visualizer
7. **Start speaking** - agent responds!

**Test checklist:**
- ✅ Agent joins within 5 seconds
- ✅ Audio visualizer shows your voice
- ✅ Agent responds in under 2 seconds
- ✅ You can interrupt agent mid-sentence
- ✅ Multi-turn conversation works
- ✅ Agent remembers previous messages
- ✅ Live transcription appears (if enabled)

**Step 6.3: Deploy to Production**
Deploy to LiveKit Cloud for 24/7 availability:
```bash
# Stop dev mode (Ctrl+C)

# Link your project
lk project add
# Browser opens → Select your project → Confirm

# Deploy agent
lk agent create
```

What happens:
1. CLI detects Python app
2. Auto-generates `Dockerfile`
3. Builds container image (2-5 minutes)
4. Uploads to LiveKit Cloud
5. Agent starts running 24/7

You'll see:
```
✓ Agent created successfully
✓ Name: letta-voice
✓ Status: Running
✓ URL: https://cloud.livekit.io/projects/your-project/agents/letta-voice
```

**Step 6.4: Test Production Agent**
1. Return to https://agents-playground.livekit.io
2. Now your production agent is available
3. Click "Connect" → Agent joins automatically
4. This agent runs 24/7, accessible from anywhere

**Step 6.5: Invite Others to Test**
Share playground access:
1. Go to LiveKit Cloud Dashboard
2. Project Settings → "Invite Members"
3. Enter email addresses
4. They can test your agent in playground

---

## Cost Optimization Strategy

### Maximizing Azure Credits

**Azure for Startups Credits:**
- Initial credit: $200 (30-day validity)
- Extended: Up to $5,000 total ($1,000 immediate + $4,000 after BizSpark verification)
- Valid for 12 months

**What Azure Credits Cover:**
- ✅ Azure OpenAI Realtime API (voice model)
- ✅ Azure Virtual Machines (if self-hosting LiveKit)
- ✅ Azure Kubernetes Service (for agent deployment)
- ✅ Azure Storage (for recordings, logs)
- ✅ Azure Monitor (observability)
- ❌ LiveKit Cloud subscription (separate billing)
- ❌ Deepgram/Cartesia (separate billing, but Deepgram has $200 free credit)

### Cost Breakdown (Real Numbers)

**Per-Conversation Costs:**

| Component | Provider | Cost per Minute | Azure Credits? |
|-----------|----------|-----------------|----------------|
| Voice Model | Azure OpenAI Mini Realtime | $0.08 | ✅ YES |
| Speech-to-Text | Deepgram | $0.0048 | ❌ No (but $200 free) |
| Text-to-Speech | Cartesia | $0.05 | ❌ No |
| Agent Hosting | LiveKit Cloud | $0.02-$0.04 | ❌ No |
| **Total per Minute** | | **~$0.15** | ~50% covered |

**Monthly Cost Projections:**

| Usage Level | Minutes/Month | Cost (Unoptimized) | Cost (Optimized) | Azure Credits Used |
|-------------|---------------|-------------------|------------------|---------------------|
| **Light** (1 hour/day) | 1,800 | $270 | $144 | $144 |
| **Medium** (5 hours/day) | 9,000 | $1,350 | $720 | $720 |
| **Heavy** (24/7, 10% utilization) | 4,320 | $648 | $346 | $346 |

**Optimization Strategies Breakdown:**

1. **Use Cached Audio Tokens** (90% savings on context)
   - Unoptimized: $0.30/min → Optimized: $0.08/min
   - Implementation: Cache system prompts, conversation history
   - Savings: ~$200/month at medium usage

2. **Context Compression** (50% token reduction)
   - Summarize every 20 messages
   - Remove filler content ("OK", "Sure")
   - Savings: ~$100/month at medium usage

3. **Smart Model Routing**
   - Simple queries → gpt-4o-mini ($0.08/min)
   - Complex reasoning → gpt-4o ($0.30/min)
   - Ratio: 80% mini, 20% full model
   - Blended cost: ~$0.12/min (vs $0.30/min always-full)

4. **Self-Host LiveKit on Azure**
   - LiveKit Cloud: $199/month + $0.02/min
   - Azure VM B2s: ~$30/month (covered by credits)
   - Azure bandwidth: $0.05/GB (covered by credits)
   - **Savings: $200+/month, fully covered by Azure credits**

### Hybrid Cloud Cost Optimization

**Recommended Stack for Maximum Azure Credit Utilization:**

```
┌─────────────────────────────────────────────────┐
│             USER (Web/Mobile/Phone)             │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│          LiveKit Self-Hosted on Azure           │
│         (Azure Kubernetes Service - AKS)        │
│         Cost: ~$75/month → Azure Credits ✅     │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│         Azure OpenAI Realtime API (Mini)        │
│         Cost: $0.08/min → Azure Credits ✅      │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│                Letta Cloud                      │
│       (Custom Azure Keys → Azure billing)       │
│         Cost: $0 (using your Azure API)         │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│           Speech Providers (Cached)             │
│   Deepgram: $200 free → Cartesia: $0.05/min    │
└─────────────────────────────────────────────────┘
```

**Total Monthly Cost (Medium Usage, 9,000 min):**
- Unoptimized (all cloud services): $1,350
- Optimized (hybrid + Azure credits): $450 out-of-pocket
  - Azure OpenAI: $720 (covered by credits)
  - LiveKit self-hosted: $75 (covered by credits)
  - Deepgram: $0 (free credit)
  - Cartesia: $450 (only real expense for months 1-12)

**After Azure Credits Exhausted (Month 13+):**
- Consider migrating to Amazon Nova Sonic: $0.06/min (80% cheaper)
- Total monthly cost drops to ~$540 for same usage

### Cost Monitoring Setup

**Azure Cost Management Dashboard:**
```bash
# Install Azure CLI
# Windows: winget install Microsoft.AzureCLI
# macOS: brew install azure-cli

# Login
az login

# Set up budget alert
az consumption budget create \
  --resource-group voice-ai-rg \
  --budget-name voice-agent-monthly \
  --amount 500 \
  --time-grain Monthly \
  --start-date 2025-01-01 \
  --notifications \
    actual-percentage=50 \
    actual-percentage=75 \
    actual-percentage=90 \
    forecasted-percentage=100
```

**Daily Cost Check Script:**
Create `check_costs.py`:
```python
import os
from azure.monitor.query import MetricsQueryClient
from azure.identity import DefaultAzureCredential

credential = DefaultAzureCredential()
client = MetricsQueryClient(credential)

def check_daily_openai_cost():
    # Query OpenAI token usage
    resource_uri = f"/subscriptions/{subscription_id}/resourceGroups/voice-ai-rg/providers/Microsoft.CognitiveServices/accounts/voice-agent-openai"
    
    metrics = client.query_resource(
        resource_uri=resource_uri,
        metric_names=["TokenTransaction"],
        timespan=timedelta(days=1)
    )
    
    total_tokens = sum(t.total for t in metrics.metrics[0].timeseries[0].data)
    
    # Calculate cost (approximate)
    # Mini Realtime: $0.08 input + $0.32 output per 1M tokens
    estimated_cost = (total_tokens / 1_000_000) * 0.40
    
    print(f"Today's OpenAI spend: ${estimated_cost:.2f}")
    print(f"Monthly projection: ${estimated_cost * 30:.2f}")
    
    if estimated_cost > 20:  # $20/day = $600/month
        print("⚠️ WARNING: High spend detected!")

check_daily_openai_cost()
```

Run daily: `python check_costs.py`

---

## Scaling Roadmap (Post-24 Hours)

### Week 2: Multi-Agent Coordination

**Supervisor-Worker Pattern Implementation:**

1. **Create Specialist Agents in Letta ADE** (1 hour)
   - Research Agent: Web search, data gathering
   - Analysis Agent: Data analysis, reasoning
   - Memory Manager: Long-term knowledge storage
   - Action Agent: API calls, task execution

2. **Configure Agent Tags** (30 min)
   ```python
   # Via Letta API or ADE
   research_agent.tags = ["worker", "research", "web"]
   analysis_agent.tags = ["worker", "analysis", "reasoning"]
   memory_agent.tags = ["worker", "memory"]
   action_agent.tags = ["worker", "actions", "api"]
   ```

3. **Update Supervisor Logic** (2 hours)
   ```python
   # In supervisor agent's tools
   async def route_query(user_query: str):
       # Analyze query type
       if "search" in query or "latest" in query:
           # Route to research agent
           return await send_to_agent_async(
               agent_id="research-agent-id",
               message=user_query
           )
       elif "analyze" in query or "why" in query:
           return await send_to_agent_async(
               agent_id="analysis-agent-id",
               message=user_query
           )
       # ... etc
   ```

4. **Test Multi-Agent Flow** (1 hour)
   - User: "Search for and analyze the latest AI market trends"
   - Supervisor → Research Agent (searches web)
   - Research Agent → Supervisor (returns data)
   - Supervisor → Analysis Agent (analyzes data)
   - Analysis Agent → Supervisor (returns insights)
   - Supervisor → User (synthesized response)

**Expected Timeline:** 4-5 hours total

### Week 3: Virtual Avatar Integration

**Tavus Digital Twin Setup:**

1. **Create Tavus Account \& Replica** (1 hour)
   - Sign up at https://tavus.io
   - Record 2-minute video:
     - Look at camera, well-lit
     - Speak naturally, various expressions
     - Read provided script
   - Upload video → AI processes (takes 2-4 hours)

2. **Get Replica ID** (5 min)
   - Once processed, copy `replica_id` from dashboard
   - Create persona: Define personality, voice tone
   - Get `persona_id`

3. **Install Tavus Plugin** (10 min)
   ```bash
   pip install livekit-plugins-tavus
   ```

4. **Update Agent Code** (30 min)
   ```python
   from livekit.plugins import tavus
   
   async def entrypoint(ctx: JobContext):
       await ctx.connect()
       
       # Create agent session with Letta
       session = AgentSession(
           llm=letta.LLM.with_letta(agent_id="your-agent-id")
       )
       
       # Create avatar
       avatar = tavus.AvatarSession(
           replica_id=os.getenv("TAVUS_REPLICA_ID"),
           persona_id=os.getenv("TAVUS_PERSONA_ID")
       )
       
       # Start both
       await session.start(room=ctx.room)
       await avatar.start(session, room=ctx.room)
   ```

5. **Deploy \& Test** (30 min)
   ```bash
   lk agent deploy
   ```
   - Open playground
   - Your digital twin appears in video
   - Speaks with synchronized lip movements
   - Maintains eye contact with camera

**Expected Timeline:** 2-3 hours (+ 2-4 hour wait for replica generation)

### Month 2: Telephony Expansion

**SIP Integration for Phone Calls:**

1. **Choose SIP Provider** (30 min)
   - Recommended: Twilio (easiest setup)
   - Alternative: Telnyx (cheaper for high volume)
   - Sign up, get phone number

2. **Configure LiveKit SIP Trunk** (1 hour)
   ```bash
   # Via CLI
   lk sip inbound create \
     --name "voice-agent-inbound" \
     --numbers "+15551234567" \
     --trunk-id "your-trunk-id"
   
   # Create dispatch rule
   lk sip dispatch create \
     --name "route-to-voice-agent" \
     --trunk-id "your-trunk-id" \
     --rule-type "direct" \
     --room-prefix "call-"
   ```

3. **Configure Twilio** (30 min)
   - Point phone number webhook to LiveKit SIP URL
   - Test: Call your number → LiveKit receives call
   - Agent answers automatically

4. **Update Agent for Phone Context** (1 hour)
   ```python
   async def entrypoint(ctx: JobContext):
       await ctx.connect()
       
       # Detect if call is from phone (SIP)
       is_phone_call = ctx.room.metadata.get("sip_trunk_id") is not None
       
       if is_phone_call:
           # Adjust behavior for phone
           instructions = """
           You are speaking on a phone call.
           Keep responses concise (under 30 seconds).
           Ask if user wants info repeated.
           Speak numbers clearly (e.g., "one, two, three").
           """
       else:
           instructions = "Standard web conversation instructions..."
       
       session = AgentSession(llm=..., instructions=instructions)
       await session.start(room=ctx.room)
       
       # Greet caller
       await session.generate_reply(
           instructions="Greet the caller warmly and ask how you can help"
       )
   ```

5. **Test Phone Integration** (30 min)
   - Call your number
   - Agent answers within 2 rings
   - Have full conversation
   - Test DTMF (press numbers) if needed

**Expected Timeline:** 3-4 hours total

**Cost Implications:**
- Twilio phone number: $1/month
- Incoming calls: $0.0085/minute
- Outgoing calls: $0.013/minute
- LiveKit SIP processing: Included in agent session minutes
- Total added cost: ~$8-15/month for moderate phone usage

### Month 3: Frontend/Web Interface

**React Frontend with Letta-LiveKit:**

1. **Clone Frontend Starter** (10 min)
   ```bash
   git clone https://github.com/livekit-examples/agent-starter-react
   cd agent-starter-react
   npm install
   ```

2. **Configure for Your Agent** (1 hour)
   - Update `.env.local`:
     ```
     NEXT_PUBLIC_LIVEKIT_URL=wss://your-project.livekit.cloud
     LIVEKIT_API_KEY=your_key
     LIVEKIT_API_SECRET=your_secret
     ```
   - Edit `app/page.tsx` to customize UI:
     - Your branding/logo
     - Custom colors
     - Agent name/avatar

3. **Add Letta-Specific Features** (2 hours)
   - Display agent memory in sidebar:
     ```typescript
     // Fetch agent memory from Letta API
     const [agentMemory, setAgentMemory] = useState(null);
     
     useEffect(() => {
       fetch(`https://api.letta.com/v1/agents/${agentId}/memory`, {
         headers: { Authorization: `Bearer ${LETTA_API_KEY}` }
       })
       .then(res => res.json())
       .then(data => setAgentMemory(data));
     }, []);
     ```
   - Show conversation history
   - Allow memory editing (for admin users)

4. **Deploy Frontend** (1 hour)
   ```bash
   # Deploy to Vercel (easiest)
   npm install -g vercel
   vercel login
   vercel deploy
   ```
   - Your web app is now live at https://your-app.vercel.app
   - Users can talk to agent via browser

5. **Custom Domain** (30 min)
   - Buy domain (Namecheap, Google Domains)
   - Add to Vercel: Settings → Domains
   - Update DNS records
   - Your voice agent now at https://voice.yourdomain.com

**Expected Timeline:** 4-5 hours total

**Result:** Professional web interface where users can:
- Click to start conversation
- See live transcription
- View agent memory
- Talk via voice or type
- See avatar (if integrated)

---

## Monitoring & Maintenance

### Daily Monitoring (5 minutes/day)

**Quick Health Check:**
```bash
# Check agent status
lk agent status

# View recent logs
lk agent logs --tail 50

# Look for errors
lk agent logs --grep "ERROR" --tail 100
```

**What to Monitor:**
1. **Agent Status**: Should always show "Running"
2. **Error Rate**: Under 1% of requests
3. **Average Latency**: Under 1.5 seconds
4. **Active Sessions**: Matches expected user activity

### Weekly Reviews (30 minutes/week)

**Performance Metrics (LiveKit Dashboard):**
- Sessions count: Trending up?
- TTFT (Time to First Token): Under 1s?
- Session duration: Are conversations engaging?
- Error categories: What's failing most?

**Cost Analysis (Azure Cost Management):**
- Daily spend vs budget
- Token usage trends
- Anomaly detection (unusual spikes?)
- Forecast for month-end

**Conversation Quality Review:**
- Listen to 5-10 random recordings
- Check transcription accuracy
- Verify memory persistence
- Note areas for improvement

### Monthly Maintenance (2 hours/month)

**1. Agent Behavior Tuning** (1 hour)
- Review conversation logs in Letta ADE
- Update system prompts based on common issues
- Adjust memory summarization triggers
- Add/remove tools as needed

**2. Model Evaluation** (30 min)
- Test alternative models for cost/quality
- Compare: gpt-4o-mini vs full gpt-4o vs Nova Sonic
- Run benchmark conversations
- Update model if significant improvement

**3. Cost Optimization Review** (30 min)
- Analyze spending patterns
- Identify expensive operations
- Implement additional caching
- Renegotiate if approaching credit limits

**4. Security Updates** (30 min)
- Rotate API keys
- Review access logs
- Update dependencies: `pip install --upgrade -r requirements.txt`
- Check for CVEs in used packages

### How to Edit Agent Behavior (No Redeploy Needed)

**Via Letta Cloud Dashboard:**
1. Go to https://app.letta.com → Your Agent
2. **Edit Persona**:
   - Click on "Core Memory" section
   - Edit "Persona" block directly
   - Changes apply to next conversation immediately
3. **Add/Remove Tools**:
   - Tools section → Toggle checkboxes
   - New tools available instantly
4. **Adjust Memory Settings**:
   - Settings → Message Buffer → Change values
   - Summarization frequency updates live

**When You MUST Redeploy:**
- Changing voice model (Azure → Nova Sonic)
- Updating Python code logic
- Adding new plugins (e.g., avatar integration)
- Modifying LiveKit connection settings

**How to Redeploy:**
```bash
# Make code changes
git add .
git commit -m "Update: improved conversation flow"
git push origin main

# Deploy to LiveKit
lk agent deploy

# Rollback if issues
lk agent rollback
```

### Debugging Common Issues

**Issue: High Latency (> 2 seconds)**
- **Check**: Azure OpenAI latency in Azure Monitor
- **Fix**: 
  - Use `gpt-4o-mini-realtime` instead of full model
  - Reduce `max_message_buffer_length` to 15
  - Enable cached audio tokens
- **Verify**: Test in playground, measure TTFT

**Issue: Memory Not Persisting**
- **Check**: Letta agent has write permissions to memory blocks
- **Fix**:
  - In Letta ADE, verify memory blocks are not read-only
  - Check agent has `archival_memory_insert` tool
  - Enable sleep-time agent for background processing
- **Verify**: Ask agent to remember fact → Query it later

**Issue: High Costs**
- **Check**: Azure Cost Management → Breakdown by service
- **Fix**:
  - Implement cached audio tokens (90% savings)
  - Add context summarization every 15 messages
  - Switch to gpt-4o-mini for simple queries
- **Verify**: Monitor daily spend → Should drop 50%+

**Issue: Agent Not Joining Rooms**
- **Check**: `lk agent status` → Should show "Running"
- **Fix**:
  - Verify `LIVEKIT_URL` in environment variables
  - Check API key validity: `lk token create --test`
  - Restart agent: `lk agent deploy`
- **Verify**: Playground should show agent available

### What to Version Control

**DO Commit to GitHub:**
```
✅ main.py (agent code)
✅ requirements.txt
✅ Dockerfile
✅ livekit.toml
✅ .env.example (template only)
✅ README.md
✅ tests/ (automated tests)
```

**DO NOT Commit (add to .gitignore):**
```
❌ .env (contains secrets)
❌ venv/ (Python environment)
❌ __pycache__/ (build artifacts)
❌ *.pyc (compiled Python)
❌ .letta/ (local data)
```

**Cloud-Only Configuration (Not in Git):**
- Letta agent memory blocks (managed in ADE)
- LiveKit room settings (managed in dashboard)
- API keys (environment variables)
- Monitoring dashboards (LiveKit/Azure)

### Backup Strategy

**Weekly Backups:**
1. **Export Letta Agent**:
   - ADE → Your Agent → Settings → "Export Agent"
   - Downloads `.af` file (Agent File format)
   - Store in secure location (not Git if contains sensitive info)

2. **Backup Code to GitHub**:
   ```bash
   git push origin main
   git tag -a v1.0-backup-$(date +%Y%m%d)
   git push origin --tags
   ```

3. **Document Configuration**:
   - Create `DEPLOYMENT.md` with:
     - Current model settings
     - API provider accounts
     - Environment variable list (no values)
     - Known issues and workarounds

**Recovery Process:**
If agent breaks:
1. Check `lk agent logs` for errors
2. Rollback to previous version: `lk agent rollback --version 1`
3. If Letta agent corrupted: Import `.af` backup in ADE
4. If code broken: `git revert` to last working commit

---

## Frequently Asked Questions

**Q: Can I use this commercially?**
A: Yes, all components (Letta, LiveKit, Azure OpenAI) have commercial licenses. Check each provider's terms for specific limits.

**Q: What happens when Azure credits run out?**
A: You'll start paying standard Azure rates. Consider migrating to Amazon Nova Sonic (80% cheaper) or optimizing further with cached tokens.

**Q: Can users interrupt the agent mid-sentence?**
A: Yes! Azure OpenAI Realtime API has built-in turn detection and interruption handling. Agent stops when user speaks.

**Q: How many concurrent users can this handle?**
A: Depends on your LiveKit plan:
- Developer: ~10 concurrent
- Starter: ~50 concurrent
- Pro: ~200 concurrent
- Enterprise: Custom (thousands)

**Q: Can the agent make phone calls proactively?**
A: Yes, with LiveKit SIP integration. Use `lk sip call` command or API to initiate outbound calls.

**Q: Is the conversation data secure?**
A: Yes. Azure OpenAI is SOC2/HIPAA compliant. LiveKit uses end-to-end encryption. Letta stores data encrypted at rest. Ensure you configure proper access controls.

**Q: Can I use a different LLM provider?**
A: Yes! Replace `azure.RealtimeModel` with:
- `openai.realtime.RealtimeModel` (OpenAI directly)
- `google.gemini.LiveModel` (Gemini Live API)
- Custom HTTP endpoint for any provider

**Q: How long does memory persist?**
A: Letta memory persists indefinitely. Conversation history is kept until you manually delete it. Archival memory has unlimited storage.

**Q: Can I test before deploying?**
A: Yes! Use `python main.py console` for local terminal testing. Then `python main.py dev` for cloud testing without production deployment.

**Q: What if I'm not technical at all?**
A: Use Letta Cloud's Agent Development Environment exclusively:
- No-code agent creation
- Visual tool management
- Test directly in web interface
- For voice integration, hire a developer for initial setup (~4 hours), then you manage behavior in ADE yourself

---

## Quick Reference Card

**Essential Commands:**
```bash
# Local testing
python main.py console          # Terminal-based voice chat
python main.py dev             # Connect to LiveKit Cloud

# Deployment
lk agent create                # Initial deployment
lk agent deploy                # Update existing agent
lk agent rollback              # Revert to previous version

# Monitoring
lk agent status                # Check if running
lk agent logs --follow         # Live log streaming
lk agent logs --grep ERROR     # Find errors

# Cost checking
az costmanagement query        # Azure cost data
python check_costs.py          # Custom cost script
```

**Essential URLs:**
- Letta ADE: https://app.letta.com
- LiveKit Dashboard: https://cloud.livekit.io
- Agents Playground: https://agents-playground.livekit.io
- Azure Portal: https://portal.azure.com
- Letta Docs: https://docs.letta.com
- LiveKit Docs: https://docs.livekit.io

**Emergency Contacts:**
- Letta Discord: https://discord.gg/letta
- LiveKit Support: https://livekit.io/support
- Azure Support: portal.azure.com → Support

---

## Success Metrics

**By Hour 24, you should have:**
- ✅ Working voice agent responding in under 1.5 seconds
- ✅ Agent deployed 24/7 to LiveKit Cloud
- ✅ Memory persisting across conversations
- ✅ Tool calling functional (web search at minimum)
- ✅ Azure credits being used for all LLM costs
- ✅ Cost under $0.15/minute of conversation
- ✅ Basic multi-agent foundation in place

**By Week 2:**
- ✅ Multi-agent coordination working (supervisor + 2-3 workers)
- ✅ Optimized costs to under $0.10/minute
- ✅ Avatar integration (if needed for use case)

**By Month 1:**
- ✅ Telephony integration (if needed)
- ✅ Custom web frontend deployed
- ✅ Monitoring dashboards configured
- ✅ 50+ successful test conversations
- ✅ Ready for beta users

**This completes your 24-hour implementation plan. You now have everything needed to build, deploy, and scale a production-ready voice AI agent system optimized for Azure credits with low-code approaches throughout. Good luck!**