# Hermes Agent — The Complete Guide

> **Title:** HERMES AGENT: THE COMPLETE GUIDE. From Zero to Self-Improving AI Employee
> **Author:** YanXbt ([@IBuzovskyi](https://x.com/ibuzovskyi))
> **Source:** [x.com/ibuzovskyi/status/2059675518966894767](https://x.com/ibuzovskyi/status/2059675518966894767?s=12)
> **Type:** X Article | **Posted:** 2026-05-27 | **Ingested:** 2026-06-07

---

## Table of Contents

1. [What Hermes Agent Actually Is](#1-what-hermes-agent-actually-is)
2. [Hermes vs OpenClaw vs Claude Code](#2-hermes-vs-openclaw-vs-claude-code)
3. [Installation](#3-installation)
4. [Model Selection](#4-model-selection)
5. [Messaging Platform](#5-messaging-platform)
6. [First Things to Do](#6-first-things-to-do)
7. [The Dashboard](#7-the-dashboard)
8. [Use Cases](#8-use-cases)
9. [Self-Improvement — The Actual Edge](#9-self-improvement--the-actual-edge)
10. [Security](#10-security)
11. [Appendix: Complete Command Reference](#appendix-complete-command-reference)

---

## 1. What Hermes Agent Actually Is

Hermes Agent is a 24/7 autonomous AI employee built by Nous Research. An agent that works while you sleep, proactively figures out tasks aligned with your goals, and gets better every single session.

**Three things separate it from everything else:**

- **Memory.** Everything lives in markdown files on your computer — not in the cloud, not in a black box. On your machine. You can read it, edit it, delete it. Full transparency.
- **Self-improvement.** Every task it completes, it reviews. What worked. What didn't. How to do it better next time. It edits its own skills after every session.
- **Session recall.** Every conversation you've ever had is logged with FTS5 full-text search and LLM summarization. Ask it what you talked about 3 months ago. It knows.

---

## 2. Hermes vs OpenClaw vs Claude Code

### Hermes vs OpenClaw

OpenClaw has gotten bloated. Feels slow. Updates break it — every single one. For someone who lives with AI agents, that's painful. For someone newer to agents, it's a dealbreaker.

Hermes is lighter, snappier, and updates don't destroy your setup. That reliability alone is the reason to switch.

**Additional Hermes advantages:**
- Built-in multi-agent via Kanban (v0.12.0+) — agents claim tasks from a board, work in parallel, hand off when blocked
- Nous Portal with curated models built in
- 166 tracked skills (87 bundled + 79 optional) across 26+ categories
- 20+ messaging platforms: Telegram, Discord, Slack, WhatsApp, Signal, Matrix, Teams, and more

### Hermes vs Claude Code vs Codex

Different tools. Different jobs. Use both.

| Tool | Best For |
|------|----------|
| **Hermes** | General purpose employee. Day-to-day tasks. Research. Documents. Spreadsheets. Computer administration. Business advisor. Prototypes. Everything that needs to continuously improve over time. Think: Chief of Staff. |
| **Claude Code / Codex** | Deep focused coding sessions. Large-scale complex applications. End-to-end testing. Two windows open side by side, locked in. |

Hermes handles everything else so you can go deep when you need to.

---

## 3. Installation

**Linux / macOS / WSL2:**
```bash
curl -fsSL https://raw.githubusercontent.com/NousResearch/hermes-agent/main/scripts/install.sh
```

**Windows (native, PowerShell — early beta):**
```powershell
iex (irm https://raw.githubusercontent.com/NousResearch/hermes-agent/main/scripts/install.ps1)
```

**Android (Termux):** same curl one-liner as Linux — installer auto-detects Termux.

**Start Hermes:**
```bash
hermes
```

Quick setup walks you through model selection and messaging platform.

> If you have OpenClaw installed: you'll get an option to import memories. **Recommendation: start clean.** Two separate agents with separate memories and skills is better than merging everything.

---

## 4. Model Selection

### Expensive — `claude-opus-4` / `claude-sonnet-4`

Best for: complex reasoning, long `/goal` tasks, nuanced writing, business advisor role. Claude handles ambiguity better than any other model.

- **opus-4** → when Hermes needs to make judgment calls mid-task
- **sonnet-4** → speed with 90% of opus quality

Catch: Anthropic disabled OAuth for agents. API key required — you pay per token. Budget: a few hundred dollars a month with heavy use.

### Moderate — `GPT-5.5`

Best for: coding tasks, prototyping, budget-conscious daily driver. Works with your existing $20/month ChatGPT subscription. Way more usage than Anthropic at a fraction of the cost.

> Start here if you're new to Hermes. Switch to Claude once you know what tasks need the upgrade.

### Affordable — Qwen 3.7 Max / Grok / Nous Portal

| Model | Best For |
|-------|----------|
| **Qwen 3.7 Max** | Long-horizon autonomous tasks. 35 hours continuous execution, 1,000+ tool calls without losing context. Top model in OpenRouter's curated picker (PR #32809). Free tier via OpenRouter. |
| **XAI Grok** | If you already pay for SuperGrok. OAuth works. Strong on X-related tasks and research. Pairs well with xurl skill. |
| **Nous Portal** | $20/month flat, curated models, built by the Hermes team. One bill, no API key management. |

**Switch models anytime:**
```bash
hermes model
```

No code changes, no reinstall. Different profiles can run different models simultaneously.

---

## 5. Messaging Platform

**Telegram. Set it up. Don't overthink this.**

Telegram is the only messaging platform actively building for AI agents. Topics. Agent-to-agent communication. Constantly adding features. Free.

Setup: Hermes walks you through it during install. Copy a token from Telegram's BotFather into Hermes. Done. Five minutes total.

Other supported platforms: Discord, Slack, WhatsApp, Signal, Matrix, Mattermost, Email, SMS, DingTalk, Feishu, Microsoft Teams, Google Chat, and more. 20+ total from one gateway process.

---

## 6. First Things to Do

### Step 1: Tell Hermes about yourself

Send this as your first message (fill in your details):

```
Here's who I am and what I'm working on:

Name: [Your name]
What I do: [your role / business]
What I'm building: [current projects]
My goals: [what you want to achieve in the next 3-6 months]
How I work: [your working style, preferences, timezone]
```

This goes into memory. Every proactive task Hermes does from now on is filtered through these goals. Skip this step and it's working blind.

### Step 2: Set up your first cron job

```
I want you to schedule a task for yourself to do every single 
night at 2am. It should be a micro app, UI, or system that helps 
us get closer to my personal goals and ambitions. It should save 
us time, make us more productive, or just be generally useful. 
Make sure it surprises me and makes me happy.
```

Wake up every morning with something new built for you.

### Step 3: Learn `/goal`

`/goal` turns the agent from reactive chatbot into a background worker. You set an objective, Hermes breaks it into tasks, executes autonomously, and continues until done.

| Command | Action |
|---------|--------|
| `/goal [description]` | Start autonomous execution |
| `/goal status` | Check what's running |
| `/goal pause` | Pause without losing context |
| `/goal resume` | Continue after pause |
| `/goal clear` | End the current goal |
| `/subgoal [text]` | Add conditions mid-execution |

---

## 7. The Dashboard

```bash
hermes dashboard
```

Opens at `localhost:9119`. **Open the Skills tab first.**

### Models tab
Swap models instantly. No terminal needed. Set different models per profile.

### Cron tab
See all scheduled tasks. Create complex ones manually with more control than typing in Telegram.

### Skills tab
Every skill Hermes has learned is here. Browse, toggle on/off, read the markdown.

A well-used agent: 150+ skills. It knows your specific workflows by name.

**Built-in skills to turn on immediately:**
- **Browser automation** — crawl the web, extract data, interact with pages
- **Computer use** — control your screen, move files, open apps
- **Image generation** — connects to OpenAI or Grok image models, generates assets mid-task
- **Video generation** — generates videos directly from conversations

### Plugins tab
Extra capabilities via API keys. Turn on browser-use, fire-crawl, computer-use.

### Profiles tab
One profile = one Hermes agent with its own memories, skills, and model. Create multiple profiles for different roles: Chief of Staff, Head of Content, Research Agent, DevOps Engineer. Run them simultaneously.

### Kanban board — the most powerful screen

**How to use it correctly:**
1. Every morning, write down your to-do list
2. Take every task an AI agent can handle → drop into Triage → walk away
3. Hermes splits tasks into subtasks → moves to To-Do → assigns to sub-agents → executes in parallel

**Statuses:** Triage → To-Do → Ready → In Progress → Blocked → Done

The daemon runs continuously (v0.16+), checking for new tasks every 60 seconds. No wasted tokens between tasks.

---

## 8. Use Cases

### 1. Daily Tutor
```
Check this video out. Read the transcript. 
Then every morning at 8am remind me of a concept or two from 
the video, then quiz me on it.
[paste YouTube link]
```
Hermes pulls the transcript, extracts key concepts, schedules a cron job. Every morning at 8am: lesson + quiz delivered to Telegram.

### 2. Computer Administrator
Install Tailscale on all your devices. Now Hermes can move across every machine you own.

```
Get me the agent.md file from my Mac Studio 
and drop it on my MacBook Pro.
```

Any file. Any device. No SSH, no Google Drive, no manual transfer.

### 3. Session Recall
```
Remind me of every business idea we discussed last month.
```
Every conversation logged. Every idea recalled on demand. FTS5 full-text search with LLM summarization across your entire history.

### 4. X Content Workflow with xurl

`xurl` is a Hermes skill that gives your agent direct access to X — searching posts, reading discussions, publishing content.

```
hermes goal "Every morning and evening, research the latest 
AI agent discussions, analyze against my content style, 
check previously published topics, generate a thread, 
evaluate virality score, and publish via xurl if it passes the bar."
```

Structure it in stages: Data Collection → Style Check → Repetition Check → Draft → Quality Score → Publish decision.

> Key mistake: don't let it publish automatically on day one. Review outputs manually for 5–7 runs first.

### 5. Mission Control
```
Build out a mission control for me where we can 
add tools to improve our workflows. Start with:
- A content pipeline where I can add ideas and you move them forward
- A memory wiki showing everything we've worked on
- A docs page showing every artifact we've created
```

### 6. Prototype Builder
```
Build me a landing page for [idea]. 
One page. Hero section, 3 features, pricing, CTA. 
Deploy to localhost. Send me the link.
```
Hermes knows your stack from memory. It doesn't ask what framework you use. It uses what you always use.

### 7. Business Advisor
```
I'm thinking about [decision or problem]. 
Here's the context: [paste context]. 
What are the 3 strongest arguments for and against? 
What would you do?
```
Unlike ChatGPT, Hermes knows your business, your goals, your constraints. Advice filtered through your actual situation.

### 8. Overnight `/goal` Runs
```
/goal build a competitor research report on [competitor 1], [competitor 2], [competitor 3].
Cover: pricing, feature set, content strategy, top performing posts last 30 days, 
estimated team size, recent funding or news.
Format as a structured markdown doc. Save to file.
```

**Set `max_turns` higher for complex overnight tasks:**
```bash
hermes config set goals.max_turns 20    # research, reports, content drafts
hermes config set goals.max_turns 50    # code refactoring, multi-step builds
```

Don't go higher without a reason. Every turn costs tokens.

### 9. Multi-Agent Org Chart

```
Profile: Chief of Staff
soul.md: You coordinate all other agents. Morning brief at 7am.
Delegate research to Head of Research. Delegate content to Head of Content.
Report back with merged summary.

Profile: Head of Research  
soul.md: You find information. Monitor X for top posts in [niche].
Track competitor activity. Send weekly intelligence brief.

Profile: Head of Content
soul.md: You produce content. Draft posts based on research briefs.
Repurpose articles into 5 formats. Maintain content calendar.
```

One morning brief in Telegram. Multiple agents working in parallel. You read the output.

---

## 9. Self-Improvement — The Actual Edge

**The loop:**
1. You give Hermes a task
2. It executes
3. After completion, it reviews: what worked, what didn't, optimal path
4. It saves that as a skill in `~/.hermes/skills/`
5. Next time you give it the same task, it uses that skill directly

You correct it once. It never makes that mistake again.

Six months in, it knows:
- How you like things formatted
- Which tools work for your specific workflows
- What tone you prefer in different contexts
- Which shortcuts to take and which to avoid

Skills are transparent. They live as markdown files on your computer. No black box.

```bash
hermes tools
```

Turn on immediately: `browser-automation`, `computer-use`

---

## 10. Security

**The honest take:** Security concerns with Hermes are vastly overrated for basic use. It only does what you tell it to do.

**What you don't need for personal use:**
- A separate iCloud account
- A dedicated Google account
- A VPS (VPS recommendations from AI YouTubers are sponsored content)

**What you do need:**
- Common sense
- Reviewing prompts before hitting enter on anything destructive
- Ground rule in `soul.md`: *"Never send money to anyone without explicit confirmation"*

> If something breaks: open the Hermes folder in Claude Code, say "I have this problem, please fix it." Solves 100% of issues.

### Production Security Stack

**Layer 1 — Bitwarden Secrets Manager:**
```bash
hermes secrets bitwarden setup   # wizard: installs bws, prompts for token
hermes secrets bitwarden status  # verify connection
hermes secrets bitwarden sync    # dry-run: see what gets applied
```

One bootstrap token in `.env`. All real credentials live in Bitwarden. Every Hermes instance pulls secrets at startup. Rotate a key once in the web app — every instance picks it up on next restart.

**Layer 2 — iron-proxy egress firewall:**
```bash
hermes egress install   # downloads iron-proxy binary, SHA-256 verified
hermes egress setup     # interactive wizard
hermes egress start     # spawn managed proxy daemon
hermes egress status    # binary + config + active mappings
hermes egress setup --from-bitwarden  # pull real credentials from BSM at proxy startup
```

Instead of injecting real credentials into the sandbox, Hermes gives the agent opaque proxy tokens. iron-proxy intercepts at the network boundary, swaps for the real credential, forwards the request. The sandbox never holds the actual key.

The two layers compose: rotate in Bitwarden, propagates to all sandboxes automatically.

---

## Appendix: Complete Command Reference

```bash
# Start
hermes

# Model selection
hermes model

# Dashboard
hermes dashboard                        # opens localhost:9119

# Goals
/goal [description]                     # start autonomous execution
/goal status                            # check what's running
/goal pause                             # pause without losing context
/goal resume                            # continue after pause
/goal clear                             # end the current goal
/subgoal [text]                         # add conditions mid-execution

# Config
hermes config set goals.max_turns 20   # research/reports
hermes config set goals.max_turns 50   # complex builds

# Tools
hermes tools                            # see all available tools

# Secrets (production)
hermes secrets bitwarden setup
hermes secrets bitwarden status
hermes secrets bitwarden sync

# Egress firewall (production)
hermes egress install
hermes egress setup
hermes egress start
hermes egress status
hermes egress setup --from-bitwarden
```

---

## Resources

**Official:**
- [Hermes Agent Docs](https://github.com/NousResearch/hermes-agent) — installation, configuration, full CLI reference
- Skills Hub — community skills, browse and install
- GitHub — source, issues, PRs

**Related articles by @ibuzovskyi:**
- The Complete Hermes /goal Playbook — 21 Workflows
- Hermes /goal — The Full Guide
- How to Make Hermes + xurl Actually Work as a System
- Hermes x Bitwarden — The Security Stack AI Agents Actually Need
