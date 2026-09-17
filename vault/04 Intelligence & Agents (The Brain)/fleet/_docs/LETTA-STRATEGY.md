# Letta Strategy — OpenClaw Integration

**Status**: OpenClaw is the primary agent interface. Letta is a backend service + dev tool.
**Credits**: $1,000 USD in Letta token credits across all supported models.

---

## The Shift

Letta is no longer the primary AI interface — OpenClaw is. Legacy lives in OpenClaw, talks to Zeus through it, and uses it to maintain context, run skills, and handle day-to-day. Letta chat UI is irrelevant to this workflow.

**But the credits are still extremely valuable — just used differently.**

Two distinct plays:

1. **Letta API** — backend memory/compute service that OpenClaw calls into
2. **Letta Code** — separate coding agent tool for *building* OpenClaw itself

---

## Play 1: Letta API as OpenClaw's Backend Brain

### What to Use It For

#### A. Archival Memory / Vector Knowledge Base
Embed the entire OneMind Codex vault into a Letta agent's archival memory. OpenClaw gains an OpenClaw skill that does semantic search against it — no need to build or host your own vector DB.

```
OneMind Codex → Letta archival (embedded) → OpenClaw `query_codex` skill → Legacy answers from real memory
```

**What to embed:**
- All `00-24 SO/` docs (frameworks, protocols, SOPs, infra configs)
- `25-49 HP/`, `50-74 LE/`, `75-99 GE/` domain docs
- Past decisions, architectural notes, session logs
- `backend/legacy/` workspace files (SOUL, AGENTS, HEARTBEAT, TOOLS)

**Credit cost**: Cheap. Embeddings use small models. Archival search queries are lightweight.

#### B. Background / Long-Running Agents
Delegate tasks to Letta agents that run for minutes/hours without OpenClaw blocking:

| Letta Agent | Task | When OpenClaw Delegates |
|-------------|------|------------------------|
| `researcher` | Deep-dive research, multi-hour synthesis | "Research X and send me a summary" |
| `analyst` | Data analysis, report generation | Dashboard queries, metrics |
| `codex-memory` | Vault knowledge search | Any knowledge retrieval task |
| `content-gen` | Draft docs, emails, plans | Content creation workflows |

OpenClaw dispatches via its `letta-dispatch` skill → Letta runs → result comes back via webhook or polling.

#### C. Free Multi-Model Experimentation
Letta gives access to all major providers on one bill. Use the credits to test which models work best for which tasks before committing to direct API contracts:

- `anthropic/claude-opus-4` — complex reasoning
- `openai/gpt-4o` — general tasks
- `google/gemini-2.0-flash` — fast/cheap
- `openai/text-embedding-3-small` — embeddings (very cheap)

Test → decide → wire best models directly into OpenClaw's model routing.

#### D. Sleep-Time Compute
Letta supports agents that "think" during downtime — background reflection processes that run when OpenClaw is idle. Useful for:

- Nightly synthesis of daily memory logs
- Periodic vault organization passes
- Background research on queued topics

### OpenClaw Implementation

Create a new skill: `backend/legacy/skills/letta-memory/`

```javascript
// index.js
const LETTA_BASE_URL = process.env.LETTA_BASE_URL; // https://api.letta.com/v1
const LETTA_API_KEY = process.env.LETTA_API_KEY;
const CODEX_AGENT_ID = process.env.LETTA_CODEX_AGENT_ID;

export const skills = [
  {
    name: 'query_codex',
    description: 'Semantic search across the OneMind Codex knowledge base',
    parameters: {
      query: { type: 'string', description: 'What to search for' },
      limit: { type: 'number', default: 5 }
    },
    handler: async ({ query, limit = 5 }) => {
      const res = await fetch(
        `${LETTA_BASE_URL}/agents/${CODEX_AGENT_ID}/archival?query=${encodeURIComponent(query)}&limit=${limit}`,
        { headers: { Authorization: `Bearer ${LETTA_API_KEY}` } }
      );
      const results = await res.json();
      return results.map(r => r.text).join('\n\n---\n\n');
    }
  },
  {
    name: 'delegate_to_letta',
    description: 'Delegate a long-running task to a background Letta agent',
    parameters: {
      agent_id: { type: 'string', description: 'Letta agent ID to delegate to' },
      task: { type: 'string', description: 'Task description' }
    },
    handler: async ({ agent_id, task }) => {
      const res = await fetch(`${LETTA_BASE_URL}/agents/${agent_id}/messages`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${LETTA_API_KEY}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ messages: [{ role: 'user', content: task }] })
      });
      const data = await res.json();
      return data?.messages?.at(-1)?.content ?? 'Task dispatched.';
    }
  }
];
```

Add to `.env.example`:
```bash
LETTA_API_KEY=sk-let-xxx
LETTA_BASE_URL=https://api.letta.com/v1
LETTA_CODEX_AGENT_ID=agent-xxx   # The codex memory agent
LETTA_RESEARCHER_AGENT_ID=agent-xxx
```

### Codex Upload Script

Run once to populate the knowledge base:

```bash
#!/bin/bash
# scripts/upload_codex_to_letta.sh
AGENT_ID="your-codex-agent-id"
VAULT="/Users/zeuslegacy/Documents/onemindcodex-v0.4"

find "$VAULT" -name "*.md" -not -path "*/.git/*" -not -path "*/.obsidian/*" | while read file; do
  content=$(cat "$file")
  relative="${file#$VAULT/}"
  echo "Uploading: $relative"
  curl -s -X POST "$LETTA_BASE_URL/agents/$AGENT_ID/archival" \
    -H "Authorization: Bearer $LETTA_API_KEY" \
    -H "Content-Type: application/json" \
    -d "{\"text\": \"FILE: $relative\n\n$content\"}" > /dev/null
done

echo "Done. Codex embedded in Letta archival memory."
```

---

## Play 2: Letta Code as the Dev Tool for Building OpenClaw

**Letta Code** (`npm install -g @letta-ai/letta-code`) is a memory-first coding agent — think Claude Code, but with persistent memory across sessions and model-agnostic.

### Why It Matters for This Project

Every other coding agent (Claude Code, Gemini CLI, Codex) forgets your context when the session ends. Letta Code persists a memory repository (git-backed) so the agent gets *smarter about your codebase over time* — not just in one session.

**Key features:**

| Feature | What It Does | Why It Matters for OpenClaw |
|---------|-------------|----------------------------|
| `/init` | Agent deep-reads your codebase and forms structured memory | Run once on the OpenClaw repo — agent knows all patterns, conventions, file relationships |
| `/remember` | Explicit trigger to reflect and learn from recent session | Run after major changes to entrench patterns |
| Skill Learning | Learns reusable `.md` skills from past experience | Teach it OpenClaw skill patterns once, it applies them forever |
| Context Repos | Git-backed memory with progressive disclosure | Agent memory is versioned, sharable, inspectable |
| Sleep-time compute | Background reflection during downtime | Agent processes your last session and forms better context overnight |
| Model-agnostic | Swap any model (Claude, GPT-4o, Gemini) | Use whichever model is best for the task, costs billed to Letta credits |
| `/search` | Semantic search over all past conversations | Find where you discussed a particular pattern or decision |

### Getting Started

```bash
# Install
npm install -g @letta-ai/letta-code

# Set up with your Letta credentials
export LETTA_API_KEY=sk-let-xxx

# Start in the OpenClaw repo
cd ~/Documents/onemindos-v0.4

# Launch
letta-code

# Initialize memory for the project (first time)
> /init

# Agent deeply reads the codebase, forms structured memory repository
# Future sessions: agent already knows OpenClaw architecture, patterns, file layout
```

### Teach It OpenClaw-Specific Skills

After working through a complex task, trigger skill learning:

```
> /learn  # or the skill learning command
```

Skills to develop with it over time:
- **openclaw-skill-pattern** — how to create new OpenClaw skills (structure, exports, env vars)
- **nats-subject-convention** — subject naming rules for inter-runtime comms
- **agno-agent-preset** — how to add new Agno worker agents
- **letta-archival-upload** — batch uploading documents to Letta archival

These get stored as `.md` files in the git-backed memory repo — versionable and reusable across future sessions.

### Context Repositories (Feb 2026 — Latest Feature)

Memory is stored as a git repo locally, synced to the Letta server:
- Every memory update is a git commit with an informative message
- Concurrent subagents can each work in isolated worktrees and merge
- `/memfs` command enables it for existing agents
- Memory files have YAML frontmatter for progressive disclosure

This means: after `/init` on the OpenClaw repo, the agent has a structured, versioned understanding of the entire codebase that carries across every future session.

---

## Credit Allocation Strategy

| Use Case | Model | Est. Cost | Priority |
|----------|-------|-----------|----------|
| Codex vault embeddings (batch upload) | `text-embedding-3-small` | ~$5-10 | First |
| Archival search queries (daily use) | `gpt-4o-mini` | ~$1-5/mo | High |
| Letta Code dev sessions | `claude-sonnet-4-5` or `gpt-4o` | ~$20-50/session | High |
| Background researcher agent | `claude-sonnet-4-5` | $10-30/task | Medium |
| Model comparison testing | Various | $50-100 | One-time |
| Sleep-time / background reflection | `gpt-4o-mini` | ~$5-10/mo | Low |

**Avoid**: Using expensive models (Opus, GPT-4 Turbo) for anything that gpt-4o-mini can handle (embeddings, archival search, simple delegation).

---

## Architecture Summary

```
┌─────────────────────────────────────────────────────────────┐
│                      ZEUS                                    │
│          WhatsApp / iMessage / Voice / Flutter               │
└─────────────────────────┬───────────────────────────────────┘
                          │
              ┌───────────▼──────────────┐
              │       OPENCLAW           │
              │   (Legacy's Runtime)     │
              │   Primary Interface      │
              └─────┬────────────────────┘
                    │
       ┌────────────┴────────────┐
       │                         │
       ▼                         ▼
Anthropic/OpenAI          LETTA API (Backend)
(real-time chat)          ├── Archival Memory (Codex KB)
                          ├── Background Agents (Research)
                          └── Sleep-time Compute

   SEPARATE TOOL
   ┌────────────────┐
   │  LETTA CODE    │  ← Terminal coding agent
   │  (dev tool)    │     Used to BUILD OpenClaw
   │  /init on repo │     Memory persists across sessions
   └────────────────┘     Billed to Letta credits
```

---

## Immediate Next Steps

1. **Run `/init` with Letta Code** on the onemindos-v0.4 repo — one-time setup
2. **Create `codex-memory` agent** in Letta (via API or ADE)
3. **Run the codex upload script** to embed the vault
4. **Create `backend/legacy/skills/letta-memory/`** skill in OpenClaw
5. **Add Letta env vars** to docker-compose and .env.example
6. **Test archival search** from OpenClaw via the new skill

---

*Last updated: 2026-02-28*
*Related: [CLAUDE.md](CLAUDE.md), [backend/legacy/skills/](backend/legacy/skills/)*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
