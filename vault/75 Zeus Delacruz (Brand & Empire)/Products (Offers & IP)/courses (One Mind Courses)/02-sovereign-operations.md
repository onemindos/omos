# Course 2: SOVEREIGN OPERATIONS

> **Deploy the full stack. Your AI runs 24/7.**

**Tier:** Legacy ($97/mo)
**Circle Space:** ⚙️ SO — Sovereign Operations (ID: 2617120)
**Prerequisite:** Course 1 (ONE MIND OS) — vault live, Legacy agent running locally
**Target Student:** Course 1 graduate ready to go from local to always-on production
**Outcome:** Full One Mind OS stack deployed on VPS — all agents live, automations running, system operating 24/7 without you babysitting it

---

## What You'll Build

By the end of this course:
- A production VPS running OpenClaw with 3–9 agents active 24/7
- Tailscale private network connecting your Mac, VPS, and iPhone
- All 9 domain agents configured with correct identities and vault access
- n8n automation stack running: morning brief, task sync, event pipelines
- Vault on GitHub with auto-commit from VPS agents
- systemd services keeping everything alive after reboots
- A security-hardened server you actually own

---

## Course Structure — 8 Modules, 45 Lessons

---

### MODULE 1: THE PRODUCTION MINDSET
*"Running locally is a toy. Running in production is a life."*

**Lesson 1.1 — Why Local Isn't Enough**
Your Mac sleeps. Your agents don't have to. This lesson makes the case for always-on infrastructure.
- The problem: agent only works when your laptop is open
- The goal: Legacy runs while you sleep, travels, and works
- What "production" means for a personal AI system
- The One Mind OS stack: what lives where

**Lesson 1.2 — The Stack Overview**
Everything you're about to build — in one diagram.
- Mac: VS Code, Claude Code, Obsidian vault (source of truth)
- GitHub: vault repo, version control, sync bridge
- VPS (onemindos-core): OpenClaw gateway, all 9 agents, MCP servers
- VPS (omos-titan): n8n automations, databases, services
- Tailscale: private network connecting everything
- Cloudflare: domains, DNS, optional tunnel

**Lesson 1.3 — Choosing Your VPS**
Where to host your agents. Honest breakdown of options.
- DigitalOcean Lightsail: cheapest start ($5–$10/mo), what Zeus uses
- DigitalOcean Droplet: more control, same provider
- Hetzner: best price/performance for EU
- AWS EC2: overkill unless you know AWS
- What specs you actually need: 2 vCPU, 2GB RAM minimum for OpenClaw
- Region choice: closest to you = lowest latency

---

### MODULE 2: VPS SETUP FROM ZERO
*"Your first server. You own this. No one can take it from you."*

**Lesson 2.1 — Provision Your VPS**
Create the server. Step-by-step, no experience required.
- Create account on DigitalOcean / Lightsail
- Create Ubuntu 22.04 LTS instance
- SSH key setup: generate on Mac, add to server
- First login: `ssh root@your.server.ip`
- The feeling of owning a server

**Lesson 2.2 — Server Hardening (Non-Negotiable)**
Security before anything else. A 30-minute process that prevents 99% of attacks.
- Create a non-root user with sudo
- Disable root SSH login
- SSH key-only authentication (disable passwords)
- UFW firewall: allow only what you need (22, 80, 443, Tailscale)
- fail2ban: auto-ban brute force attempts
- Unattended upgrades: auto-patch security vulnerabilities
- The hardening checklist (downloadable PDF)

**Lesson 2.3 — Install the Runtime Stack**
Everything OpenClaw needs to run.
- Node.js 20+ via nvm (the right way)
- npm, git
- `npm install -g openclaw`
- Clone your vault: `git clone git@github.com:yourusername/your-vault.git`
- SSH deploy key: passwordless vault push from VPS to GitHub

**Lesson 2.4 — systemd Services: Keeping Agents Alive**
Your agents need to restart automatically after crashes and reboots.
- What systemd is and why it matters
- Writing a service file for OpenClaw gateway
- `systemctl enable`, `start`, `status`, `logs`
- Watching logs in real-time: `journalctl -u openclaw-gateway -f`
- Testing: reboot the server, watch agents come back up

---

### MODULE 3: TAILSCALE — YOUR PRIVATE NETWORK
*"Zero-config networking. Your devices talk to each other like they're in the same room."*

**Lesson 3.1 — What Tailscale Is (And Why It Changes Everything)**
The simplest way to understand Tailscale: it's a VPN that doesn't feel like a VPN.
- The problem Tailscale solves: public IPs, port forwarding, firewall hell
- How it works: WireGuard tunnels + Tailscale's coordination server
- Your Tailnet: a private network only your devices can join
- Zero open ports required for internal services

**Lesson 3.2 — Install Tailscale on Everything**
Mac, VPS, iPhone — all on the same private network.
- Install on VPS: `curl -fsSL https://tailscale.com/install.sh | sh`
- Install on Mac: download from tailscale.com
- Install on iPhone: App Store
- Login to your Tailscale account on all devices
- Verify: ping your VPS from your Mac using its Tailscale IP
- Name your devices: onemindos-core, zeus-macbook-pro, zeus-iphone

**Lesson 3.3 — Tailscale Serve: Expose Services Without Public Ports**
Access your n8n, dashboards, and OpenClaw from any device — privately.
- `tailscale serve https:443 proxy http://127.0.0.1:18789`
- Your OpenClaw dashboard now at `https://onemindos-core.tailnet-name.ts.net`
- Access from iPhone, Mac, anywhere on Tailnet
- No public exposure — only your devices can reach it

**Lesson 3.4 — Connecting OpenClaw Mac Client to VPS Gateway**
Your Mac talks to the VPS agents as a remote client.
- `openclaw config set gateway.url wss://onemindos-core.tailnet-name.ts.net`
- Device pairing: `openclaw device pair`
- Test: send a message from Mac, agent responds from VPS

---

### MODULE 4: DEPLOYING ALL 9 AGENTS
*"One agent is a tool. Nine agents are a team."*

**Lesson 4.1 — The Agent Roster and Why Each One Exists**
The 9 agents of One Mind OS — their domains, personalities, and purpose.

| Agent | Emoji | Domain | What It Handles |
|-------|-------|--------|-----------------|
| Legacy | ⚡ | All (orchestrator) | Coordination, planning, high-level decisions |
| Spartan | 💪 | HP 27-28 | Body, fitness, health, physical performance |
| Oracle | 🔮 | HP 25-26, 29-34 | Mind, money, learning, identity |
| Eden | 🌱 | LE (sustainability) | Food growing, environment, homestead |
| Haven | 🏡 | LE 50-74 | Family, home, relationships, estate |
| Empire | 💰 | GE 75-99 | Ventures, brand, wealth, business |
| Guardian | 🛡️ | SO 04 | Security, access control, digital protection |
| Grid | ⚙️ | SO 01-02 | Infrastructure, automation, technical ops |
| Trinity | 🔺 | SO 06 | Research, R&D, emerging intelligence |

**Lesson 4.2 — Agent Workspace Architecture**
How agent workspaces are structured on the VPS — and why.
- Each agent gets: `~/.openclaw/workspace-{id}/`
- Shared files: CODEX.md (universal vault context), USER.md (Zeus's profile)
- Agent-specific files: IDENTITY.md, SOUL.md (personality per agent)
- Symlinks for shared files — change once, updates all agents
- The workspace init workflow: new agent in 5 minutes

**Lesson 4.3 — Writing Agent Identities (IDENTITY.md + SOUL.md)**
What makes each agent feel like a distinct intelligence.
- IDENTITY.md: name, emoji, tagline, domain, model, workspace path
- SOUL.md: personality traits, communication style, expertise, values, guardrails
- The difference between Legacy (philosophical, big-picture) vs Spartan (direct, data-driven)
- Hands-on: write IDENTITY.md and SOUL.md for your 3 priority agents

**Lesson 4.4 — The openclaw.json Config (Full Breakdown)**
The master configuration file — every setting that matters.
- `agents.list[]`: defining each agent (id, name, model, workspace, channels)
- `agents.defaults`: shared settings (model, memory, tools)
- `channels`: Discord, Telegram, iMessage bindings
- `memory`: extraPaths pointing agents at your vault
- `gateway`: port, trustedProxies, WebSocket settings
- The #1 gotcha: `agents.defaults.model` overrides per-agent model settings
- Hands-on: configure all 9 agents in openclaw.json

**Lesson 4.5 — Vault Access: Giving Agents Eyes**
Your agents need to read (and write) your Codex. This lesson wires it up.
- MCP filesystem server: `@modelcontextprotocol/server-filesystem`
- Configure per-agent in openclaw.json: `mcpServers.codex-vault`
- `memorySearch.extraPaths`: vault path for semantic search
- Test: ask Empire agent to read your GE domain goals
- Test: ask Spartan to read your HP body goals
- The git auto-commit loop: agents write → `onemind-watch` commits → GitHub updated

**Lesson 4.6 — Channel Bindings: One Agent Per Channel**
Route conversations to the right agent automatically.
- Discord: create server, create channels per agent, bind in config
- Telegram: one bot per agent (or one bot, agent-switching by command)
- The routing logic: `channels.discord.channelBindings`
- Testing all 9 agents: send a message to each channel
- The daily workflow: which channel do you use for what?

**Lesson 4.7 — Model Configuration: Right Model for the Job**
Not every agent needs the same model. Cost optimization + performance tuning.
- Legacy (orchestrator): Claude Sonnet — needs reasoning + context
- Domain agents (Spartan, Oracle, etc): Claude Sonnet or Haiku
- Worker agents (Grid, Guardian): Haiku for cost efficiency
- The GitHub Copilot API path: `github-copilot/claude-sonnet-4.6`
- Reasoning models (Claude 3.7 Sonnet Extended): when to use
- Fallback chains: `model.primary` → `model.fallback`

---

### MODULE 5: AUTO-SYNC — YOUR VAULT IS ALWAYS CURRENT
*"What you write on Mac appears on VPS in 3 seconds."*

**Lesson 5.1 — The Sync Architecture**
How the bidirectional vault sync works — the full picture.
- Mac (Obsidian) → git push → GitHub (source of truth)
- VPS: `onemind-pull.timer` (every 10s → `git pull --rebase`)
- VPS: `onemind-watch.service` (inotifywait → agent writes → auto-commit → push)
- The result: Mac and VPS always in sync within ~15 seconds

**Lesson 5.2 — SSH Deploy Key Setup**
Passwordless VPS→GitHub push. Required for `onemind-watch` to work.
- Generate deploy key on VPS: `ssh-keygen -t ed25519 -f ~/.ssh/codex_deploy`
- Add public key to GitHub repo as deploy key (write access)
- Configure git on VPS: `GIT_SSH_COMMAND=ssh -i ~/.ssh/codex_deploy`
- Test: make a change on VPS, watch it appear on GitHub

**Lesson 5.3 — onemind-watch.service: Auto-Commit from VPS**
Every file an agent writes gets committed and pushed automatically.
- Write the `onemind-watch.service` systemd unit
- inotifywait: watching the vault folder for changes
- 3-second debounce: avoid commit storms
- Auto-commit message: `[agent] auto-sync YYYY-MM-DD HH:MM`
- Enable and start: `systemctl --user enable onemind-watch`

**Lesson 5.4 — onemind-pull.timer: Always Pulling from GitHub**
Keep the VPS vault current when you edit on Mac.
- Write the `onemind-pull.service` and `onemind-pull.timer` systemd units
- `git pull --rebase` every 10 seconds
- Handle merge conflicts gracefully
- Monitor: `journalctl -u onemind-pull -f`

---

### MODULE 6: N8N — THE AUTOMATION LAYER
*"Agents react to your life. n8n makes your life trigger your agents."*

**Lesson 6.1 — What n8n Is and Why It's Your Automation Brain**
n8n vs Zapier vs Make — the honest comparison for a sovereignty-first stack.
- n8n is open-source, self-hosted, no per-task pricing
- 400+ integrations: every SaaS you use has a node
- Visual workflow builder: no code required for most automations
- The role in One Mind OS: SaaS events → normalize → NATS → agent

**Lesson 6.2 — Deploy n8n on Your VPS**
n8n lives on `omos-titan` (or your single VPS if you have one).
- Docker Compose setup (the right way)
- Environment variables: encryption key, timezone, webhook URL
- systemd Docker service: auto-start on reboot
- Tailscale Serve: expose n8n dashboard privately
- First login and basic orientation

**Lesson 6.3 — The 5 Automations Every One Mind Setup Needs**
The starter automation pack — download and import directly.

1. **Morning Brief** — 6:00 AM → pull Todoist tasks + calendar + weather → format → send to Legacy via Discord
2. **Todoist Task Completed** → NATS event → Legacy notified → vault updated
3. **Circle.so New Member** → welcome DM via Discord → tag in member registry
4. **Evening Digest** — 9:00 PM → day summary → wins captured → tomorrow preview
5. **Weekly Review Trigger** — Sunday 7 PM → create weekly review note in vault → notify via Telegram

**Lesson 6.4 — Webhook Ingestion: Any SaaS → Your Agents**
The universal webhook pattern — works with any service that sends webhooks.
- n8n Webhook node: your inbound URL
- Parse and normalize the payload
- Format as a message for your agent
- Send via Discord webhook or OpenClaw HTTP endpoint
- Example: Stripe payment → Empire agent notified → vault updated

**Lesson 6.5 — Cron Jobs: Your Agents' Internal Clock**
Some automations aren't triggered by events — they run on schedule.
- The OpenClaw cron configuration: `heartbeat.schedule`
- Cron syntax explained simply
- Beyond morning brief: weekly domain audits, monthly reviews, quarterly resets
- The proactive agent: Legacy checks in on you without you asking

**Lesson 6.6 — Build Your First Custom Workflow**
Pick something from your life and automate it end-to-end.
- Student picks their own use case
- Walk through: identify trigger → map flow → build nodes → test → activate
- Examples: fitness tracker → Spartan agent, email from client → Empire pipeline, weather alert → Eden agent

---

### MODULE 7: MEMORY AND INTELLIGENCE
*"An agent without memory is just a chatbot. Memory is what makes it yours."*

**Lesson 7.1 — The 3-Tier Memory Architecture**
How One Mind OS stores and retrieves intelligence across time.
- Work memory: active state, current tasks, what's happening now (`_codex/memory/work/`)
- Learning memory: signals, patterns, feedback, what worked (`_codex/memory/learning/`)
- Relationship memory: interaction history, preferences, Zeus-specific context (`_codex/memory/relationship/`)

**Lesson 7.2 — How Memory Search Works**
The hybrid retrieval system — semantic + keyword.
- BM25: keyword search (fast, exact matches)
- Vector embeddings: semantic search (finds meaning, not just words)
- Hybrid: combines both for best recall
- `memorySearch.extraPaths`: what each agent can search
- Scoping memory: why Empire shouldn't search HP files

**Lesson 7.3 — Compaction and the Long Game**
Context windows fill up. This lesson explains how agents handle it gracefully.
- The compaction problem: 200K tokens fills up faster than you think
- Memory flush: saving key facts before context resets
- The Lossless Claw plugin: DAG-based summary that preserves everything
- Session memory (`/memories/session/`): cleared after conversation
- The accumulation effect: agents get smarter over months, not days

**Lesson 7.4 — Teaching Your Agents More About You**
The USER.md file is never done. This lesson builds the habit of feeding your agents.
- When to update USER.md: after life changes, new goals, preference shifts
- What to include: communication preferences, decision-making style, what you hate
- Domain-specific USER.md: Spartan knows your training history, Oracle knows your money mindset
- The "agent briefing" habit: monthly refresh of key agent files

---

### MODULE 8: MONITORING AND MAINTENANCE
*"You built it. Now keep it running."*

**Lesson 8.1 — Health Checks: Is Everything Alive?**
Daily 60-second check that your stack is healthy.
- `systemctl status openclaw-gateway`
- `journalctl -u openclaw-gateway -n 20`
- `docker ps` (n8n and databases)
- The health check script: run it, get a green/red status in 5 seconds
- Automating the health check: ping yourself every morning

**Lesson 8.2 — Observability: What's Happening in Your Stack**
See inside your running system.
- Gateway logs: who's talking to what
- n8n execution history: which workflows ran, which failed
- Memory search logs: what agents are looking up
- Disk and RAM monitoring: before you run out
- Setting up DigitalOcean alerts: CPU/RAM/disk notifications

**Lesson 8.3 — Tracing and Debugging**
When something goes wrong, you need to find the exact failure point. This lesson builds the debugging muscle.
- Enabling debug mode: `openclaw gateway --debug` for verbose output
- Reading gateway logs: understanding tool calls, memory lookups, channel messages
- The trace of one interaction: input → memory search → tool call → output → channel send
- The 5 most common failure modes: wrong tool called, stale memory, context overflow, bad parameters, guardrail block
- Reproduction: capture the exact input that caused the failure and replay it
- Isolation: test agent alone vs full stack to narrow the problem
- Fix → verify → monitor: don't just fix, confirm it stays fixed
- Hands-on: introduce a deliberate misconfiguration, debug it using logs

**Lesson 8.4 — Sessions Management**
How to monitor and manage active sessions in your gateway.
- Listing active sessions: `openclaw sessions list`
- Session lanes at a glance: how many main/cron/subagent lanes are in use
- Killing a stuck session: `openclaw sessions kill <session-id>`
- Session idle timeout: configuring when sessions auto-expire
- The session overflow condition: all lanes full — what happens to new messages
- Monitoring session health: metrics to watch, alert thresholds
- Hands-on: simulate a session backlog and observe how the gateway handles it

**Lesson 8.5 — Backup Strategy**
Your vault is on GitHub. Your agent configs need protection too.
- What to back up: `~/.openclaw/`, vault (already on GitHub), n8n workflows
- `openclaw export`: backup gateway config
- DigitalOcean automated snapshots: $2/mo for full VPS backups
- Offsite: n8n workflow JSON exports to vault `_system/` folder

**Lesson 8.4 — Course 2 Graduation Challenge**
Your stack is live. Prove it.
1. All agents responding on Discord
2. Morning brief arriving automatically at 6 AM
3. Vault syncing bidirectionally (Mac edit → VPS → GitHub within 15 seconds)
4. n8n dashboard accessible via Tailscale
5. Post your `openclaw status` screenshot + morning brief screenshot to Inner Circle

---

## Resources Included

| Resource | Format | Access |
|----------|--------|--------|
| VPS setup + hardening script | Shell script | Legacy+ |
| `openclaw.json` full production config (9 agents) | Download | Legacy+ |
| All 9 `IDENTITY.md` templates | Download pack | Legacy+ |
| All 9 `SOUL.md` templates | Download pack | Legacy+ |
| `onemind-watch.service` + `onemind-pull.timer` unit files | Download | Legacy+ |
| n8n 5-automation starter pack | JSON import | Legacy+ |
| VPS hardening checklist | PDF | Legacy+ |
| Health check script | Shell script | Legacy+ |
| docker-compose.yml for n8n stack | Download | Legacy+ |

---

## Up Next → Course 3: CLAUDE CODE & AI TOOLS

You have a running stack. Course 3 teaches you to build on top of it — writing your own skills, extending your agents, and using Claude Code to develop with AI as your pair programmer.
