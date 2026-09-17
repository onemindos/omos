---
id: ui-08-openclaw-practical-guide
title: "OpenClaw Practical Guide — Cut Token Spend 3x & Build Agents That Work"
type: reference
status: active
created: 2026-03-05
modified: 2026-03-05
owner: zeus
tags: [type/reference, status/active, domain/so, platform/openclaw]
ai_summary: "Practical guide to OpenClaw internals — 6 components, dual memory system, gateway message flow, tools + cron automation, multi-agent setups, and 5 common mistakes that waste tokens or break security."
---

# OpenClaw Practical Guide

> Cut token spend by 3x and build agents that actually work — by understanding what's happening under the hood.

---

## Why This Matters

Most people use OpenClaw blind. They install skills, write prompts, tweak settings — but have no idea what's happening under the hood. A black box. Even though it's completely transparent.

Key things most users don't know:

- OpenClaw has **two types of memory** — bootstrap and semantic search. If you're only using one, your agent is running at half capacity.
- You can control your agent from a **plain text editor** and get results 95% better than most people.
- If you're running your agent "out of the box", you're spending **3x more tokens** (= money) than you need to.

---

## The Six Parts of OpenClaw

OpenClaw is built from six parts. Each one works independently, but together they form a full-fledged agent.

### 1. Gateway — The Core

Everything goes through it. It maintains persistent connections to Telegram, WhatsApp, Discord, Slack. When a message comes in — Gateway decides which agent gets it, pulls up the history, assembles context and sends it to the LLM.

The response flies back the same way. On top of that it spins up a WebSocket API so you can plug in your own interface or external integration.

### 2. Agent — The Brain

Receives assembled context from Gateway: chat history, memory files, available tools. Thinks, decides which tool to call, builds a response. If needed — kicks off a chain: called a tool, got a result, thought some more, called another one. Keeps going until the final answer is ready.

### 3. Tools — The Hands

- **exec** — runs shell commands on your server
- **browser** — opens pages, clicks, takes screenshots
- **file** — reads and writes files
- **message** — sends messages to channels
- **memory** — searches through long-term notes

Each tool is a separate capability you can turn on or off.

### 4. Workspace — Long-Term Memory

A folder with files where everything the agent needs to know between sessions is written down. Who it is, what tone to use, what you like, what decisions you've made together, what happened yesterday. **Without workspace the agent wakes up with a blank slate every single time.**

### 5. Sessions — Conversation Memory

Full history of a specific dialogue: what you wrote, what it replied, which tools it called. Each session lives on its own and doesn't bleed into others (unless you messed up the settings — more on that at the end).

### 6. Nodes — Physical Devices

Your Mac, phone, remote server. They connect to Gateway and expand what the agent can do: snap a photo, take a screenshot, grab geolocation. Gateway on the server is the brain. Node on your Mac is the eyes and hands.

**All of this is text files.** Not a database, not binaries. Plain `.md` and `.json` you open in an editor and change by hand.

---

## Workspace — The Superpower Nobody Configures

Workspace is a folder with files. Sounds boring. But it's exactly what separates a dumb bot from an agent you can actually work with.

Without workspace the agent wakes up with a blank head every time. Doesn't remember who you are. Doesn't remember what you discussed last week. Doesn't remember decisions you made together. Every conversation starts from zero. And every time you're spending tokens just to re-explain context.

### Workspace Files and Their Roles

| File | Role |
|------|------|
| `AGENTS.md` | **Operating manual.** How the agent should think. When to use which tool, safety rules, order of operations. Its playbook. |
| `SOUL.md` | **Personality.** Tone, boundaries, priorities. Want brief answers with no unsolicited advice? Put it here. |
| `USER.md` | **Your profile.** How to address you, what you do, what you prefer. Read before every single response. |
| `MEMORY.md` | **Long-term memory.** Facts that must not get lost. "We only trade on DEX, no CEX." "Primary RPC is Alchemy." |
| `YYYY-MM-DD.md` | **Daily logs.** What happened today, tasks in progress, what you discussed. Tomorrow it picks up context. |
| `IDENTITY.md` | **Name and vibe.** "My name is Legacy, I'm your systems assistant." Sets the tone for everything. |
| `HEARTBEAT.md` | **Periodic check checklist.** "Check email", "see if monitoring is running." |
| `TOOLS.md` | **Local tool hints.** Where scripts live, which commands are available, shortcuts. |

---

## The Two Levels of Memory (And Why Most People Only Use One)

### Level 1: Bootstrap — Always In Context

Every time it runs, Gateway takes `AGENTS.md`, `SOUL.md`, `USER.md`, `IDENTITY.md`, today's daily log and injects them into context **before the LLM even sees your message**. This is bootstrap — the first level of memory. The agent "remembers" the contents of these files every single time.

### Level 2: Semantic Search — On-Demand Recall

If the memory plugin is enabled, the agent can search through `MEMORY.md` and other notes via a vector index. You ask "which DEX do we trade on?" — it finds the right chunk even if you wrote it two months ago. Not by keywords — **by meaning**. You wrote "we only trade on DEX" — the agent finds it when you ask "where do we open positions."

### The Key Difference

| | Bootstrap | Semantic Search |
|---|---|---|
| **When used** | Every single message | Only when relevant |
| **Token cost** | Constant (eats tokens every request) | Minimal (only pulls what's needed) |
| **Reliability** | Guaranteed in context | No guarantee the right fact is found |
| **Best for** | Critical rules, identity, tone | Decisions, facts, task context |

### The Strategy

Put the **critical stuff in bootstrap** — what the agent needs in every message. Tone, rules, who you are. **Everything else goes into `MEMORY.md` and daily logs** — decisions, facts, task context. Semantic search will pull them when needed.

- If only bootstrap is configured → agent runs at **half power**
- If neither is configured → you're just **burning money** on tokens every day

### Example

You tell the agent: "Remember that I only trade on DEX and I watch pools on Curve." It writes this to `MEMORY.md`. A week later you ask: "Where do we usually check liquidity?" The agent searches its memory and answers without extra questions. You didn't spend tokens re-explaining. It just knows.

---

## Gateway — How a Message Becomes a Response

Gateway is the thing everything goes through. Every message, every response, every tool call. If the agent is the brain, Gateway is the circulatory system. It stops — everything dies.

It's a long-running daemon process. You start it once — it sits there and works. Maintains persistent connections to channels, waits for messages, processes them, returns responses.

### The 5-Step Message Pipeline

**Step 1 — Receive:** Gateway maintains a persistent connection to the Telegram API. An event comes in — a new message from you.

**Step 2 — Route:** Gateway checks the config: which agent should handle this message? Bindings are defined there — channel → agent mappings. It also determines the `SessionId`: is this a continuation or a new session.

**Step 3 — Build Context:** Gateway reads the session history from the `.jsonl` file. Pulls in bootstrap files from workspace: `AGENTS.md`, `SOUL.md`, `USER.md`, everything else. Adds available skills. Packs it all together and sends it to the LLM.

**Step 4 — Agent Loop:** The LLM returns either text or a tool call. If it's a tool call — Gateway executes it: runs a command, opens a browser, writes a file — and feeds the result back into context. The LLM looks at the result, thinks further. Might call another tool. And another. The loop keeps spinning until a final answer appears.

**Step 5 — Persist:** The response streams back to Telegram. You see "typing..." in real time. Meanwhile the entire exchange — your message, the agent's response, all tool calls — gets written to `.jsonl`. `sessions.json` gets updated.

### API Access

Gateway spins up a WebSocket API on port `18789`. Through it you can plug in your own UI — a web interface, a mobile app. Or integrate with external systems via HTTP. There's even an OpenAI-compatible endpoint.

By default Gateway only listens on **localhost**. For remote access:

- **Tailscale** — bridge Gateway into your private network
- **SSH tunnel** — `ssh -L 18789:127.0.0.1:18789 user@server`
- **Never expose 18789 to the open internet** — that's full access to all your data

---

## Tools + Cron — The Agent That Works Without You

Tools are how the agent interacts with the outside world. Without them it can only talk. With them — it acts.

### exec — The Most Powerful (and Dangerous) Tool

Runs shell commands. Through it the agent can do practically anything: run a script, install a package, process a file, deploy code.

**Three modes — from safe to risky:**

| Mode | Description | Risk |
|------|-------------|------|
| `sandbox` | Runs inside a Docker container. Separate sandbox, no connection to your server. | Safe |
| `gateway` | Runs on your server with restrictions. Whitelist controls which commands are allowed. | Medium |
| `full` | No restrictions. The agent does whatever it wants. | Dangerous |

### browser — Browser Automation

The agent can open a page, click an element, type text, take a screenshot, save a PDF. Two profiles:

- **openclaw** — fully isolated, managed
- **chrome** — controls your regular Chrome via an extension

### Cron — The Game Changer

This is what turns an agent from a chatbot into a worker. You set up a schedule once and forget:

```bash
openclaw cron add --schedule "0 9 * * *" \
  --agent personal \
  --prompt "Check new emails, send summary to Telegram" \
  --announce
```

Every morning at 9:00 the agent wakes up, does what you asked, and sends the result to the channel. Without you lifting a finger.

- `--announce` — result gets delivered to the channel
- `--no-deliver` — quiet execution, agent does everything internally

### Heartbeat — Periodic Health Checks

A periodic check based on the checklist in `HEARTBEAT.md`. Gateway asks the agent to go through the items: is monitoring running, is there disk space left, are there errors in the logs. If something's wrong — the agent messages you.

### Real Example: Automated Email Summary

You want the agent to check Gmail every morning and send a summary to Telegram:

1. Enable browser tool
2. Set up a cron for 9:00
3. Write in `AGENTS.md`: "Open gmail.com, use `snapshot()` for the email list, filter by relevant senders, compose a summary and send via message tool."

Every morning the agent opens a browser with a saved session, reads the inbox, filters, sends you: *"3 emails from clients today: X wants an update, Y asks about payment, Z confirmed the meeting."*

You haven't finished your coffee — and it's already sorted everything out.

---

## Multi-Agent — One Gateway, an Army of Agents

Gateway can serve as many agents as you want. Each agent is a separate folder in `~/.openclaw/agents/`. Its own workspace, its own sessions, its own memory files.

### How It Works

- **Work agent** — knows about your stack, project, team
- **Personal agent** — knows about habits, schedule, preferences
- **Monitoring agent** — checks servers via heartbeat, alerts on issues
- **Research agent** — parses sources with browser tool, saves summaries
- **Trading agent** — watches pools, pings you on opportunities

Channel mapping is done through **bindings** in `config.json`. Write to one Telegram chat — goes to the work agent. Write to another — to the personal one. Same Gateway, routing by rules.

### Session Isolation via dmScope

| Mode | Behavior |
|------|----------|
| `"main"` | All DMs in one session (default — **unsafe for multi-user**) |
| `"per-channel-peer"` | Each user gets isolated session (**safe**) |
| `"per-agent"` | Isolation by agent (for multi-agent setups) |

**Main rule:** If more than one person has access to an agent, you **must** set `dmScope: "per-channel-peer"`. Without it, sessions from different users collapse into one.

---

## 5 Mistakes That Break Your Agent

### Mistake 1: dmScope set to "main" with multiple users

OpenClaw's default — all direct messages in one channel get dumped into a single session. Two different people message your Telegram bot → the agent sees both as one conversation. It can respond to one person with context from another's conversation.

**Fix:** Set `dmScope: "per-channel-peer"` in config. Each user gets an isolated session.

### Mistake 2: exec tool set to full mode

You've given the LLM full shell access to your server. No whitelist, no restrictions. The agent can delete files, install anything, read any data. On a live server — a disaster.

**Fix:** Switch to `sandbox` or `gateway` with a proper `exec-approvals.json`.

### Mistake 3: No workspace or it's empty

An agent without workspace is an LLM without context. Every conversation starts from a blank slate. You explain who you are, what you do, how to communicate — every single time. Every time you spend tokens on the same thing.

**Fix:** At minimum set up `AGENTS.md` + `SOUL.md` + `USER.md`. That's 15 minutes of work that pays for itself from the first conversation.

### Mistake 4: No compaction strategy

Long dialogues grow into thousands of tokens. OpenClaw can compress old messages into a summary, but if the agent didn't write the important stuff to `MEMORY.md` before compression — it's gone. The agent forgets decisions you made three days ago.

**Fix:** Enable memory flush before compaction — the agent automatically dumps important facts into long-term memory.

### Mistake 5: Port 18789 exposed to the internet

Gateway API listens on localhost by default. But if you forwarded the port for remote access and forgot to lock it down — that's full access to all agents, sessions, data. Anyone who finds that port can read conversations, send commands, modify workspace.

**Fix:** Tailscale or SSH tunnel. Never expose 18789 directly.

---

## Summary

OpenClaw is not a black box. Every component is a text file you can open and edit. Every session is JSONL you can read and parse. Every config is JSON you control.

| Action | Impact |
|--------|--------|
| Configure workspace (bootstrap files) | Agent remembers who you are — every time |
| Enable semantic search (MEMORY.md) | Agent recalls decisions from weeks ago |
| Use both memory levels | 3x token savings vs re-explaining context |
| Set up cron + heartbeat | Agent works while you sleep |
| Fix dmScope + exec + port | No security disasters |

Set up your workspace, use both levels of memory, close the holes from the checklist — and your agent will operate on a different level.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
