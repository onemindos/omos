# Two-Tier LLM Strategy

> LLM #1 (Agno) thinks deeply. LLM #2 (OM1) reacts fast. Neither replaces the other.

## Why Two LLMs

A single LLM cannot serve both physical reflexes and deep reasoning. When the robot is walking toward a wall, it needs a response in 200ms. When you ask it to research competitors and draft an email, it needs 3-5 seconds of deep reasoning. These are fundamentally different workloads.

Think of it like the human nervous system: your spinal cord pulls your hand off a hot stove (200ms, no conscious thought) while your prefrontal cortex decides what to have for dinner (seconds of deliberation).

## LLM #1 — The Thinker (Agno)

| Attribute | Value |
|-----------|-------|
| **Location** | Cloud (API calls) |
| **Models** | GPT-4.1, Claude Sonnet 4.5, Gemini 2.5 Pro |
| **Response time** | 1-5 seconds |
| **Cost** | ~$0.002-0.01 per request |
| **Invocation** | On-demand (chat, webhooks, heartbeat) |
| **Calls per hour** | 2-50 (varies by usage) |

**What it handles:**
- Complex reasoning and planning
- Knowledge base retrieval (RAG)
- Long-term memory operations
- Multi-agent team coordination
- Digital tool execution (email, calendar, APIs)
- Interpreting robot sensor events
- Generating responses for chat/voice
- Heartbeat analysis (every 15 min)

**Model selection per agent:** Agno supports different models per agent within a team. Use this strategically:

```python
# Research agent — needs deep reasoning
research_agent = Agent(model=OpenAIChat(id="gpt-4.1"))

# Quick response agent — needs speed
triage_agent = Agent(model=OpenAIChat(id="gpt-4.1-mini"))

# Code agent — needs strong coding
code_agent = Agent(model=Claude(id="claude-sonnet-4-5-20250929"))
```

## LLM #2 — The Reflex (OM1)

| Attribute | Value |
|-----------|-------|
| **Location** | Local (on-device) or low-latency cloud |
| **Models** | Qwen3-30B (local), GPT-4.1-nano (cloud fallback) |
| **Response time** | 200-500ms |
| **Cost** | $0 (local) or ~$0.0001 per request |
| **Invocation** | Continuous loop (every 0.5-2 sec) |
| **Calls per hour** | ~1,800-7,200 |

**What it handles:**
- Real-time obstacle avoidance
- Interpreting sensor data into physical actions
- Translating high-level commands ("patrol") into motor sequences
- Speech generation for the robot's speaker
- Safety-critical decisions (stop, dodge, alert)

**What it does NOT handle:**
- Email, calendar, digital tools
- Knowledge base searches
- Long-term memory
- Multi-agent coordination
- Complex reasoning or planning

**OM1's dual-model mechanism:** OM1 can run local and cloud models in parallel, taking the first good response. This maximizes both speed and quality:

```json5
{
  "llm": {
    "model": "qwen3-30b",
    "provider": "ollama",
    "fallback": {
      "model": "gpt-4.1-nano",
      "provider": "openai"
    },
    "parallel": true  // run both, take first good response
  }
}
```

## How They Interact

LLM #1 and LLM #2 never call each other directly. They communicate through NATS:

```
You: "Patrol the house and email me if you see anyone"

LLM #1 (Agno) reasons:
  → Task 1: Send "patrol house" command → NATS → OM1
  → Task 2: Subscribe to person-detection events from OM1

LLM #2 (OM1) receives "patrol house":
  → Plans walking path
  → Avoids obstacles (every 0.5 sec)
  → Camera sees person → publishes event → NATS

LLM #1 (Agno) receives "person detected":
  → Drafts email
  → Sends via Gmail MCP tool
```

## Cost Projection

| Component | Calls/Hour | Cost/Call | Monthly Cost |
|-----------|-----------|-----------|-------------|
| LLM #2 (local) | ~3,600 | $0 | $0 |
| LLM #1 (chat) | ~10 | ~$0.005 | ~$36 |
| LLM #1 (heartbeat) | ~4 | ~$0.008 | ~$23 |
| LLM #1 (events) | ~5 | ~$0.003 | ~$11 |
| **Total** | | | **~$70/month** |

Running LLM #2 locally eliminates the biggest cost driver. If you ran GPT-4.1 at 3,600 calls/hour, it would cost ~$500/day.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
