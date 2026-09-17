---
id: ui-08-openclaw-deep-dive
title: "OpenClaw Deep Dive — Architecture, File Structure & Practical Scenarios"
type: reference
status: active
created: 2026-03-05
modified: 2026-03-05
owner: zeus
tags: [type/reference, status/active, domain/so, platform/openclaw]
ai_summary: "Complete technical reference for OpenClaw internals — gateway, agents, workspace, sessions, tools, nodes, memory, cron, skills. Covers file structure, data flow, and practical examples."
---

# OpenClaw: Deep Dive

## Architecture, File Structure, and Practical Scenarios

---

## Introduction: OpenClaw as an Organism

OpenClaw isn't just software — it's a distributed system organized like a living organism. Think of its anatomy:

- **Gateway** — the heart and nervous system, pumping data and coordinating all processes
- **Agent** — the brain, thinking and making decisions
- **Tools** — the hands, executing actions
- **Workspace** — long-term memory and personal space
- **Sessions** — short-term conversation memory
- **Nodes** — additional sensory organs (camera, screen, microphone)

This document breaks down each component to the file and config level, shows how they interact, and illustrates everything with practical examples.

---

## Part 1: Gateway — The Heart of the System

### 1.1 What is Gateway and Why Does It Exist

Gateway is a long-running daemon process that:

- Maintains persistent connections to channels (Telegram, WhatsApp, Discord, Slack)
- Routes incoming messages to the correct agents
- Stores session state (conversation history)
- Exposes an API (HTTP/WebSocket) for UI and external integrations
- Runs periodic tasks (cron, heartbeat)
- Manages connected nodes (devices)

**Human analogy:** Gateway is the cardiovascular + nervous system. It pumps events (like blood) between channels and agents, transmitting signals (like nerve impulses) from sensory organs (channels) to the brain (agent) and back.

---

### 1.2 Gateway File Structure

All Gateway data lives in `~/.openclaw/`:

**`~/.openclaw/config.json`** — the main Gateway config:

- Auth settings (`gateway.auth.token` / `password`)
- WebSocket API port (default: `18789`)
- Channel configs (`telegram.token`, `whatsapp.credentials`, etc.)
- Exec security settings (sandbox, approvals)
- Browser tool configs (profiles, executable path)

**`~/.openclaw/agents/<agentId>/`** — data for each agent:

- `sessions/sessions.json` — metadata for all sessions (who, when, channel, status)
- `sessions/<SessionId>.jsonl` — full history of each conversation (one line = one message)
- `workspace/` — the agent's personal memory (covered in Part 2)

**`~/.openclaw/channels/`** — persistent channel data:

- `telegram/session.json` — Telegram bot session
- `whatsapp/.wwebjs_auth/` — WhatsApp session (QR auth)

**`~/.openclaw/exec-approvals.json`** — whitelist of approved commands for the exec tool

---

### 1.3 How Gateway Processes an Incoming Message

**Step 1: Message arrives in a channel (e.g., Telegram)**

- Gateway holds a persistent connection to the Telegram API
- Receives a webhook or long-polling event: new message from a user

**Step 2: Route to agent**

- Gateway checks `bindings` in config: which agent should handle this message?
- Determines `SessionId`: for DMs, usually one main context (`dmScope: "main"`) or isolated (`dmScope: "per-channel-peer"`)
- If it's a new session, creates a record in `sessions.json` and a new `.jsonl` file

**Step 3: Launch agent loop**

- Gateway reads session history from `.jsonl`
- Injects bootstrap files from workspace (`AGENTS.md`, `SOUL.md`, `USER.md`, etc.)
- Appends available skills
- Assembles the full context and sends it to the LLM

**Step 4: Process the response**

- LLM returns either text or tool calls
- If tool call, Gateway executes it (exec, browser, file write, etc.) and returns the result to context
- The loop repeats until the LLM returns a final response
- Response is streamed back to the channel (user sees typing in real time)

**Step 5: Persist to memory**

- The entire exchange (user message + assistant response + tool calls) is written to `.jsonl`
- `sessions.json` is updated (timestamp, message count)

---

### 1.4 Practical Example: Multi-Agent System

**Use case:** You have a work agent (for company tasks) and a personal agent (for life admin). You want `@work_bot` on Telegram to be the work agent, and your personal agent to respond in DMs.

**Setup:**

1. Create two workspaces:
   - `~/.openclaw/agents/work/workspace/` — with `AGENTS.md` for work tasks
   - `~/.openclaw/agents/personal/workspace/` — with `AGENTS.md` for personal matters
2. Define bindings in `config.json`:

   ```json
   { "channel": "telegram", "chatId": "@work_bot", "agentId": "work" }
   { "channel": "telegram", "chatId": "YOUR_USER_ID", "agentId": "personal" }
   ```

3. Set session isolation: `session.dmScope: "per-agent"`

**Result:** Messaging `@work_bot` routes to the "work" agent, which only sees its own workspace and sessions. Messaging in your DMs switches to "personal" with its own workspace and sessions.

---

### 1.5 API and Remote Access

Gateway exposes a WebSocket API on port `18789` (default), enabling:

- UI clients (web interface, mobile app)
- External system integrations via HTTP API
- OpenAI-compatible endpoint for tool compatibility

**Security:** By default Gateway only listens on `127.0.0.1` (localhost). For remote access:

- **VPN/Tailscale** — expose Gateway on your private network
- **SSH tunnel** — `ssh -L 18789:127.0.0.1:18789 user@server`
- **Never expose port `18789` directly to the internet** — it gives full access to all your data

---

## Part 2: Workspace — Long-Term Memory

### 2.1 What is Workspace

Workspace is a folder of Markdown files that define:

- Agent behavior (how to think, how to work with the user)
- Personality (tone, boundaries, name)
- User knowledge (preferences, context)
- Long-term memory (important facts, decisions, notes)
- Daily journal (operational notes by day)

**Human analogy:** Workspace is a "memory home" — your notebooks, notes, journals, instructions to yourself. The things you write down so you never forget them.

---

### 2.2 Workspace File Structure

Typical structure at `~/.openclaw/agents/<agentId>/workspace/`:

**`AGENTS.md`** — the agent's operational manual:

- How to make decisions
- When to use which tools
- Security rules (e.g., "always confirm exec commands")
- Docs-first approach ("read docs first, then write code")

**`SOUL.md`** — persona and tone:

- How to communicate (formal, friendly, concise?)
- Boundaries (what the agent won't do)
- Values and priorities

**`USER.md`** — user profile:

- How to address them (name, formal/informal)
- Preferences (response format, work style)
- Context (what the user does, their goals)

**`IDENTITY.md`** — the agent's name and vibe (e.g., "My name is Alex, I'm your automation assistant")

**`TOOLS.md`** — local tool hints:

- Where scripts live (e.g., `/home/user/scripts/process.py`)
- Which commands are available in this environment
- Shortcuts and aliases

**`HEARTBEAT.md`** — checklist for periodic checks:

- "Check for new emails from clients"
- "Verify monitoring is running"
- "Update status in Notion"

**`YYYY-MM-DD.md`** — daily journal:

- Brief notes about the day's events
- What was done, what decisions were made
- Current work context ("currently debugging module X")

**`MEMORY.md`** — long-term memory (summary of important things):

- Important facts that must never be forgotten
- Decisions made and their reasoning
- Lessons from past experience
- Usually only available in main/private context (not in public chats)

---

### 2.3 How the Agent Uses Workspace

**Bootstrap files:** Each time the agent loop starts, Gateway reads specific files from workspace and injects them into context before calling the LLM:

- `AGENTS.md` (always)
- `SOUL.md` (always)
- `USER.md` (always)
- `IDENTITY.md` (always)
- `TOOLS.md` (if present)
- Today's journal (`YYYY-MM-DD.md`)
- `MEMORY.md` (only for main/private context)

**Memory search:** If the memory plugin is enabled (default: `memory-core`), the agent can:

- Semantically search `MEMORY.md` and other `.md` files via a vector index
- Find relevant facts even if not explicitly mentioned in the bootstrap

**Writing to memory:** The agent can write to workspace if `workspace writable` is enabled (default: yes). Examples:

- User: "Remember that I prefer React with TypeScript" → Agent writes this to `MEMORY.md`
- User: "Log in the journal that we fixed the auth bug today" → Agent appends to `2025-02-15.md`

---

### 2.4 Practical Example: Personal Assistant with Memory

**Use case:** You want an agent that remembers your preferences, keeps a daily journal, and can search past notes.

**Setup:**

1. Create workspace at `~/.openclaw/agents/personal/workspace/`
2. Write in `USER.md`: "User: [Name], informal address, prefers short answers"
3. Write in `AGENTS.md`: "Always keep a journal. Log important events in `YYYY-MM-DD.md`, long-term decisions in `MEMORY.md`"
4. Tell the agent: "Remember that I work in Web3 and I'm interested in automation"
5. Agent writes this to `MEMORY.md`

**Result:**

- A week later you ask: "What did we do on Monday?" → Agent opens `2025-02-10.md` and tells you
- You ask: "What are my stack preferences?" → Agent searches `MEMORY.md` and responds: "You prefer React with TypeScript, you work in Web3"

---

## Part 3: Sessions — Short-Term Memory

### 3.1 What is a Session

A session is "conversation memory" — the context of one dialogue between a user and agent. Each session:

- Stores message history (user, assistant, tool calls)
- Is isolated from other sessions (unless configured otherwise)
- Can be private (DM) or group (chat)
- Has its own `SessionId` (unique identifier)

**Human analogy:** A session is "working memory" for a conversation. While you're talking, you remember the context of that specific discussion. End the conversation, switch topics — new session, new context.

---

### 3.2 Session File Structure

**`~/.openclaw/agents/<agentId>/sessions/sessions.json`** — metadata for all sessions:

- `SessionId`, channel, user
- Creation and last-message timestamps
- Message count, status (active/archived)

**`~/.openclaw/agents/<agentId>/sessions/<SessionId>.jsonl`** — full conversation history:

- Each line = one event (user message, assistant response, tool call, tool result)
- JSONL format: JSON Lines (each line is valid JSON)
- Easy to read and parse incrementally

---

### 3.3 dmScope: DM Session Isolation

**Problem:** By default, all DMs in one channel collapse into a single main context (`dmScope: "main"`). This means:

- If two different people DM your Telegram bot, the agent sees their messages in one session
- Context leak risk: the agent might accidentally reply to one user with information from another

**Solution:** DM session isolation via `session.dmScope`:

- `"main"` — all DMs in one session (unsafe for multi-user)
- `"per-channel-peer"` — each user gets their own isolated session (safe)
- `"per-agent"` — isolation by agent (for multi-agent setups)

**Recommendation:** If more than one person DMs your bot, always set `session.dmScope: "per-channel-peer"`.

---

### 3.4 Compaction: History Compression

Long conversations can grow to tens of thousands of tokens. OpenClaw automatically compacts history:

- Older messages are rolled into a summary
- Recent messages remain full
- Before compaction, Gateway can do a "memory flush" — asking the agent to write important things to `MEMORY.md`

---

## Part 4: Tools — The Agent's Hands

### 4.1 Tool Types

Tools are how the agent interacts with the outside world:

- **exec** — run shell commands
- **browser** — control a browser (open pages, click, take screenshots)
- **file** — read/write files
- **message** — send messages to channels
- **memory** — search long-term memory
- **node.*** — commands on connected devices (`camera.*`, `screen.*`, `location.*`, etc.)

---

### 4.2 exec Tool: Running Commands

`exec` lets the agent run shell commands. It's the most powerful — and most dangerous — tool.

**Three execution modes:**

- `host=sandbox` — runs in an isolated Docker container (default if sandboxing is on)
- `host=gateway` — runs on the Gateway machine (with policies and approvals)
- `host=node` — runs on a connected node (macOS app / headless)

**exec security levels:**

- `deny` — exec completely disabled
- `allowlist` — only commands in `exec-approvals.json` are permitted
- `full` — everything allowed (dangerous!)

**`exec-approvals.json`** — command whitelist:

```json
{ "pattern": "/usr/bin/git", "approved": true }
{ "pattern": "/usr/bin/python3", "args": ["script.py"], "approved": true }
{ "pattern": "rm", "approved": false }
```

---

### 4.3 browser Tool: Managed Browser

Browser tool is a separate isolated browser profile the agent can control programmatically.

**Two modes:**

- `openclaw` — managed profile, full isolation, no extension needed
- `chrome` — control your regular Chrome via an extension relay (requires open tab with extension)

**Commands:**

- `tabs.open(url)` — open a page
- `tabs.close(tabId)` — close a tab
- `click(selector)` — click an element
- `type(selector, text)` — enter text
- `snapshot()` — get page HTML/text
- `screenshot()` — take a screenshot
- `pdf()` — save page as PDF

**Config** (`openclaw.json`):

```json
"browser": {
  "enabled": true,
  "defaultProfile": "openclaw",
  "executablePath": "/usr/bin/chromium"
}
```

---

### 4.4 Practical Example: Automated Email Check

**Use case:** Every morning the agent checks your Gmail, finds emails from clients, and sends you a summary in Telegram.

**Setup:**

1. Enable browser tool (`browser.enabled: true`)
2. Add a cron job in Gateway:

   ```bash
   openclaw cron add --schedule "0 9 * * *" --agent personal --prompt "Go to Gmail, check emails from clients, send summary to Telegram"
   ```

3. Add instructions in `AGENTS.md`:
   > "To check Gmail: open gmail.com, click Inbox, use snapshot() to get the email list, filter emails from clients (domain @client.com), compose a summary and send via message tool to Telegram"

**Result:** Every morning at 9:00 the agent opens the browser, logs into Gmail, reads inbox via `snapshot()`, filters client emails, and sends you a Telegram summary.

---

## Part 5: Nodes — Additional Sensory Organs

### 5.1 What is a Node

A node is a physical device (computer, smartphone, tablet) connected to Gateway that provides additional capabilities:

- Camera (`camera.capture`)
- Screen (`screen.record`, `screen.capture`)
- Location (`location.get`)
- System commands (`system.run`)
- Canvas (drawing on screen)

**Human analogy:** Nodes are additional limbs and sensory organs. Gateway on a server is the brain; the node on your Mac is eyes, hands, and ears.

---

### 5.2 Connecting a Node

**macOS app:**

1. Download OpenClaw
2. Launch it, enter Gateway address (e.g., via Tailscale)
3. Node registers with Gateway and announces its capabilities

**Headless node:**

1. Run `openclaw node start` on a remote machine
2. Set Gateway address in config
3. Node connects and waits for commands

---

### 5.3 Practical Example: Home Monitoring

**Use case:** You want the agent to take screenshots of your home computer on demand.

**Setup:**

1. Gateway running on a VPS
2. OpenClaw running on your home Mac, connected to Gateway via Tailscale
3. Mac announces capability: `screen.capture`
4. You type in Telegram: "Show me my screen"
5. Agent calls `screen.capture` on the node → Node returns screenshot → Agent sends it to Telegram

**Result:** Remote screen monitoring without VNC or TeamViewer.

---

## Part 6: Memory — The Agent's Hippocampus

### 6.1 Two Levels of Memory

**Bootstrap (direct injection):** Files from workspace are always read in full and added to context before each run.

- Guaranteed to be in context
- Token-heavy if files are large

**Semantic search (vector index):** The memory plugin indexes `MEMORY.md` and other `.md` files, searching for relevant chunks by meaning.

- Efficient search over large note sets
- Doesn't guarantee the right fact will be found

---

### 6.2 Right Memory Strategy

**`MEMORY.md`** — long-lived facts that should ALWAYS be remembered:

- User preferences
- Decisions made and their reasoning
- Lessons from past experience ("don't use library X, it's buggy")

**`YYYY-MM-DD.md`** — events and context for a specific day:

- What was done today
- Current tasks and debugging context
- Temporary notes that can be archived later

**How to write:**

- Direct instruction: "Log in memory that I work in Web3"
- Agent decides autonomously: if workspace is writable and memory flush is enabled, the agent can automatically write important things before session compaction

---

## Part 7: Cron and Heartbeat — Daily Rhythm

### 7.1 Cron: Scheduled Tasks

Cron is a task scheduler on the Gateway side. Configure tasks like:

- "Every day at 9:00 check email and send summary"
- "Every 30 minutes check monitoring status"
- "Once a day back up workspace to git"

**Command:**

```bash
openclaw cron add --schedule "0 9 * * *" --agent work --prompt "Check for new tasks in Notion" --announce
```

**Flags:**

- `--announce` — result is delivered to a channel (e.g., sent to Telegram)
- `--no-deliver` — runs internally, result not sent

---

### 7.2 Heartbeat: The Agent's Pulse

Heartbeat is a short periodic check that Gateway asks the agent to perform:

1. Agent reads `HEARTBEAT.md`
2. Executes the checklist
3. If something's wrong, writes to channel (e.g., "Monitoring is down!")

**Example `HEARTBEAT.md`:**

```markdown
- Check if monitoring.py process is running
- Check if >10GB disk space is free
- Check for new errors in logs
```

---

## Part 8: Skills — Training the Agent

### 8.1 What is a Skill

A skill is a folder with a `SKILL.md` that teaches the agent how to use a specific tool or procedure.

**Example:** skill for working with PDFs at `~/.openclaw/skills/pdf/SKILL.md` — contains instructions on reading PDFs, extracting text, merging, converting to images.

---

### 8.2 Skill Sources

- **Bundled** — built-in skills, shipped with OpenClaw
- **`~/.openclaw/skills`** — global skills for all agents
- **`<workspace>/skills`** — skills for a specific agent (override global)

**Priority:** workspace > global > bundled

---

### 8.3 Gating: Conditional Skill Activation

A skill can be gated on specific conditions:

- Presence of a binary (e.g., only show git skill if `/usr/bin/git` exists)
- Config flag (e.g., only show browser skill if `browser.enabled=true`)
- Environment variable (e.g., `API_KEY` is set)

---

## Conclusion: Putting It All Together

OpenClaw is not just "an LLM with tools." It's a distributed system with:

| Component     | Role                                                                         |
| ------------- | ---------------------------------------------------------------------------- |
| **Gateway**   | Central coordinator — holds connections, routes messages, stores sessions     |
| **Workspace** | Long-term memory — defines behavior, personality, and agent knowledge        |
| **Sessions**  | Short-term memory — stores conversation context                              |
| **Tools**     | Actuators — let the agent act in the real world                              |
| **Nodes**     | Sensory organs — extend agent capabilities to physical devices               |
| **Skills**    | Training materials — teach the agent to use tools effectively                |

Understanding each component — where it lives on the filesystem, how it interacts with the rest — lets you:

- Configure complex multi-agent scenarios
- Safely grant agents access to critical systems
- Build durable memory that persists across sessions
- Automate routine tasks via cron and heartbeat
- Extend agent capabilities via nodes and custom skills

**Core principle:** OpenClaw is not a black box. Every component is transparent, every config is editable, every session is stored in plain JSONL. **You are always in full control of the system.**


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
