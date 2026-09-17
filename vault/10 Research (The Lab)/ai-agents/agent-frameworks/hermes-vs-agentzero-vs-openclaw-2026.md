---
title: "Hermes Agent vs Agent Zero vs OpenClaw — Unbiased Comparison"
type: reference
status: active
created: 2026-04-06
owner: zeus
tags: [type/reference, domain/so, topic/agents, topic/research]
ai_summary: "Research-based comparison of three personal AI agent platforms (Hermes Agent, Agent Zero, OpenClaw) for the OneMind life OS use case. Covers ecosystem size, features, protocols, extensibility, robotics/hardware support, and forkability. All data sourced from live research on April 6, 2026."
---

# Hermes Agent vs Agent Zero vs OpenClaw — Unbiased Comparison

> All data in this document is sourced from live research (GitHub repos, official websites, documentation) on April 6, 2026. No opinions — just numbers, features, and facts.

---

## Critical Finding: Hermes Agent ↔ OpenClaw Relationship

Before comparing, a key fact: **Hermes Agent has a built-in migration tool from OpenClaw** (`hermes claw migrate`). The Hermes GitHub repo tags include "openclaw" and "clawdbot." The migration tool imports SOUL.md, memories, skills, command allowlists, messaging settings, API keys, TTS assets, and workspace instructions. Hermes can also install OpenClaw skills from ClawHub.

This means Hermes is either a fork of or a spiritual successor to OpenClaw — built by Nous Research (an AI research lab) with added capabilities around model training, RL environments, and research pipelines. **Your existing OpenClaw investment is fully portable to Hermes if you ever choose to migrate.**

---

## The Numbers (April 6, 2026)

| Metric | OpenClaw | Hermes Agent | Agent Zero |
|--------|----------|--------------|------------|
| **GitHub Stars** | 349,379 | 26,900 | 16,800 |
| **Forks** | 70,021 | 3,500 | 3,400 |
| **Contributors** | 1,539 | 272 | 40 |
| **Releases** | 82 | 6 (first: Feb 2026) | 48 |
| **Latest Release** | v2026.4.5 | v0.7.0 (2026.4.3) | v1.7 (3 days ago) |
| **Primary Language** | TypeScript | Python (93%) | Python (52%), JS/HTML (44%) |
| **License** | MIT | MIT | MIT |
| **Backed By** | Community (massive) | Nous Research (AI lab) | Agent Zero s.r.o. (Czech company) |
| **Open Issues** | — | 659 | 205 |
| **Open PRs** | — | 1,100 | 161 |
| **Commit Activity** | Active | Extremely active (commits hours ago) | Active (commits 3 days ago) |

**Ecosystem scale is not even close.** OpenClaw has 13x Hermes's stars and 20x Agent Zero's contributors. But Hermes is the fastest-growing — going from 0 to 26.9K stars in ~2 months (launched Feb 2026).

---

## Feature Comparison

### Messaging Platforms

| Platform | OpenClaw | Hermes Agent | Agent Zero |
|----------|----------|--------------|------------|
| Telegram | ✅ | ✅ | ❌ (via API only) |
| Discord | ✅ | ✅ | ❌ |
| Slack | ✅ | ✅ | ❌ |
| WhatsApp | ✅ | ✅ | ❌ |
| Signal | ✅ | ✅ | ❌ |
| Email | ✅ | ✅ | ❌ |
| SMS | ✅ | ✅ | ❌ |
| Matrix | ✅ | ✅ | ❌ |
| Mattermost | ✅ | ✅ | ❌ |
| Home Assistant | ✅ | ✅ | ❌ |
| DingTalk | ✅ | ✅ | ❌ |
| Feishu/WeCom | ✅ | ✅ | ❌ |
| Web UI | ✅ | ❌ (CLI-first) | ✅ (primary) |
| CLI/Terminal | ✅ | ✅ (primary) | ❌ |
| **Total Platforms** | **28+** | **14+** | **1 (Web UI) + REST API** |

**Verdict:** OpenClaw leads with 28+ platforms. Hermes is strong at 14+. Agent Zero is Web UI only — no native messaging gateway. For a life OS where you message your agent from your phone, this is a major gap for Agent Zero.

### Skills / Extensions Ecosystem

| Metric | OpenClaw | Hermes Agent | Agent Zero |
|--------|----------|--------------|------------|
| **Skill Registry** | ClawHub (47,341 skills) | agentskills.io + ClawHub + skills.sh + LobeHub | Plugin Hub (a0-plugins, ~50 plugins) |
| **Curated Collection** | Awesome Skills (19,500 stars, 2,868 curated) | 40+ bundled skills | ~8 built-in plugin-skills |
| **Skill Format** | OpenClaw native | agentskills.io standard (Anthropic open standard) | SKILL.md (compatible with Claude Code, Cursor) + plugin.yaml |
| **Cross-Compatibility** | ClawHub | Can install from ClawHub, LobeHub, GitHub, skills.sh | Independent plugin system |
| **Plugin Hub Stars** | 7,539 (ClawHub) | N/A (uses existing hubs) | 46 (a0-plugins) |

**Verdict:** OpenClaw's 47,341 skills on ClawHub is an enormous moat. Hermes is smart — instead of building its own registry, it can install from ClawHub, LobeHub, skills.sh, and GitHub URLs. Agent Zero's plugin ecosystem is the smallest (46 stars, ~50 plugins, 25 contributors), but its SKILL.md compatibility with other agents is a plus.

### Protocols & Interoperability

| Protocol | OpenClaw | Hermes Agent | Agent Zero |
|----------|----------|--------------|------------|
| **MCP** (tool connectivity) | ✅ Client | ✅ Client + Server mode | ✅ Client + Server |
| **ACP** (Agent Communication) | ✅ | ✅ (editor integration: VS Code, Zed, JetBrains) | ❌ |
| **A2A** (Agent-to-Agent) | ❌ | ❌ | ✅ (FastA2A) |
| **AG-UI** (Agent-User) | ❌ | ❌ | ❌ |

**Verdict:** No platform supports the full Protocol Triangle yet. Agent Zero is the only one with A2A support. Hermes has ACP for editor integration. All three support MCP.

### Memory & Learning

| Capability | OpenClaw | Hermes Agent | Agent Zero |
|------------|----------|--------------|------------|
| **Persistent Memory** | ✅ | ✅ | ✅ (FAISS vector search) |
| **Cross-Session Recall** | ✅ | ✅ (FTS5 + LLM summarization) | ✅ (categorized: main, fragments, solutions) |
| **Auto Skill Creation** | ✅ (via `self-improving-agent` skill — 3K ⭐, 355K downloads, auto-extracts recurring learnings into new skills) | ✅ (built-in, creates skills from experience) | ❌ |
| **Skill Self-Improvement** | ✅ (via `Self-Improving + Proactive Agent` — 889 ⭐, 148K downloads, tiered memory with auto-promotion/demotion) | ✅ (built-in, skills improve during use) | ❌ |
| **User Modeling** | ❌ | ✅ (Honcho dialectic modeling) | ❌ |
| **Memory Dashboard** | ❌ | ❌ | ✅ (browse, search, edit, delete) |
| **Knowledge Base** | ✅ (Memory-Wiki — actively developed, committed hours ago) | ❌ | ✅ (custom document directory) |
| **Nudge System** | ❌ | ✅ (agent-curated periodic nudges) | ❌ |
| **Dreaming / Memory Consolidation** | ✅ (built-in — background sweep scores short-term signals, promotes to long-term MEMORY.md with threshold gates) | ❌ | ❌ |

**Verdict:** Both OpenClaw and Hermes have strong learning loops. Hermes has it built-in natively; OpenClaw achieves equivalent capability via ClawHub skills (the #1 most popular skill on the platform) plus built-in Dreaming for memory consolidation. Agent Zero has the best memory management UI. Key difference: Hermes has unique user modeling via Honcho.

### Execution Environments

| Environment | OpenClaw | Hermes Agent | Agent Zero |
|-------------|----------|--------------|------------|
| **Local** | ✅ | ✅ | ✅ |
| **Docker** | ✅ | ✅ | ✅ (primary, with Kali Linux) |
| **SSH Remote** | ❌ | ✅ | ❌ |
| **Daytona** (serverless) | ❌ | ✅ | ❌ |
| **Modal** (serverless GPU) | ❌ | ✅ | ❌ |
| **Singularity** (HPC) | ❌ | ✅ | ❌ |
| **Total Backends** | **2** | **6** | **2** |

**Verdict:** Hermes has the most flexible execution — 6 backends including serverless (Daytona/Modal where your environment hibernates when idle). This matters for running on a $5 VPS vs a GPU cluster.

### Model Support

| Provider | OpenClaw | Hermes Agent | Agent Zero |
|----------|----------|--------------|------------|
| **OpenAI** | ✅ | ✅ | ✅ |
| **Anthropic** | ✅ | ✅ (Claude Code auth) | ✅ |
| **OpenRouter** (200+ models) | ✅ | ✅ | ✅ |
| **Ollama** (local) | ✅ | ✅ | ✅ |
| **Azure OpenAI** | ✅ | ✅ | ✅ |
| **Nous Portal** | ❌ | ✅ (native OAuth) | ❌ |
| **DeepSeek** | ✅ | ✅ | ❌ |
| **GitHub Copilot** | ❌ | ✅ (OAuth) | ✅ |
| **Z.AI / GLM** | ❌ | ✅ | ❌ |
| **Kimi / Moonshot** | ❌ | ✅ | ❌ |
| **Alibaba / DashScope** | ❌ | ✅ | ❌ |
| **Custom Endpoint** | ✅ | ✅ | ✅ |
| **Total Providers** | **~10** | **20+** | **~8** |

**Verdict:** Hermes supports the most providers natively (20+), including many Chinese/Asian providers. All three support custom endpoints, so any OpenAI-compatible API works.

### Voice & Media

| Feature | OpenClaw | Hermes Agent | Agent Zero |
|---------|----------|--------------|------------|
| **Text-to-Speech** | ✅ (ElevenLabs) | ✅ (multiple engines) | ✅ |
| **Speech-to-Text** | ❌ | ✅ (Whisper, local) | ✅ (Whisper) |
| **Voice Mode** | ❌ | ✅ (CLI + Telegram + Discord VC) | ✅ (Web UI) |
| **Image Generation** | ✅ | ✅ | ❌ |
| **Vision/Screenshots** | ✅ | ✅ | ✅ (Browser Agent) |
| **Browser Automation** | ❌ | ✅ | ✅ (Playwright) |

### Scheduling & Automation

| Feature | OpenClaw | Hermes Agent | Agent Zero |
|---------|----------|--------------|------------|
| **Cron Scheduling** | ❌ | ✅ (natural language, platform delivery) | ❌ |
| **Subagent Spawning** | ✅ | ✅ (isolated, parallel) | ✅ (hierarchical superior/subordinate) |
| **Python RPC** | ❌ | ✅ (zero-context-cost tool pipelines) | ❌ |
| **Workflow Engine** | Lobster (1,071 stars) | ❌ | ❌ |

### Research & Training

| Feature | OpenClaw | Hermes Agent | Agent Zero |
|---------|----------|--------------|------------|
| **Batch Trajectory Generation** | ❌ | ✅ | ❌ |
| **RL Training** (Atropos) | ❌ | ✅ | ❌ |
| **Trajectory Compression** | ❌ | ✅ | ❌ |
| **ShareGPT Export** | ❌ | ✅ | ❌ |

**Verdict:** Only Hermes has research/training capabilities, which makes sense — it's built by Nous Research, a model training lab. This is irrelevant for a life OS but could matter if you want to fine-tune models on your own data.

---

## Hardware / Robotics / Physical World

This is where the OneMind use case gets specific.

| Capability | OpenClaw | Hermes Agent | Agent Zero |
|------------|----------|--------------|------------|
| **ROS 2 Bridge** | ✅ RosClaw (459 stars, won SF hackathon) | ❌ | ❌ |
| **Smart Glasses** | ✅ VisionClaw (2,000 stars, Meta Ray-Ban) | ❌ | ❌ |
| **Home Assistant** | ✅ (native platform) | ✅ (messaging gateway) | ❌ |
| **Drone Control** | Via RosClaw + ROS 2 | ❌ | ❌ |
| **Robotics** | Via RosClaw + ROS 2 | ❌ | ❌ |
| **IoT/Hardware** | Via platform plugins | Via MCP servers | Via MCP servers |

**Verdict:** OpenClaw is the ONLY platform with a robotics bridge (RosClaw) and smart glasses integration (VisionClaw). These are production projects with community traction. For the OneMind use case involving homesteading, robotics, drones, and glasses — this is a significant differentiator. Hermes and Agent Zero would require building these integrations from scratch.

---

## Forkability & Branding

| Factor | OpenClaw | Hermes Agent | Agent Zero |
|--------|----------|--------------|------------|
| **License** | MIT | MIT | MIT |
| **Can Fork & Rebrand** | ✅ | ✅ | ✅ |
| **Language** | TypeScript | Python | Python + JS |
| **Codebase Size** | Large (mature) | Large (active) | Medium |
| **Packaging** | npm / Docker | pip / Homebrew / Nix / Docker | pip / Docker |
| **White-label Ready** | ✅ (SwarmClaw precedent — your OneMind Web) | ✅ (SOUL.md + profiles) | ✅ (fully customizable prompts) |

**Verdict:** All three are MIT and forkable. You've already forked SwarmClaw (OpenClaw's multi-agent variant) into OneMind Web. That's a proven path. Hermes has a SOUL.md personality system and profiles for brand customization. Agent Zero's entire behavior is prompt-defined, making it highly customizable.

---

## OneMind Use Case Scoring

Scoring each platform 1-5 against your specific requirements:

| Requirement | OpenClaw | Hermes | Agent Zero | Notes |
|-------------|----------|--------|------------|-------|
| **Life/family/business management** | 5 | 4 | 3 | OpenClaw: you already run 9 agents + 30 skills for this. Hermes: can migrate it. A0: no native messaging. |
| **Software/apps/hardware dev** | 4 | 5 | 4 | Hermes: best code execution (6 backends, Python RPC). All can write code. |
| **Fork and brand** | 5 | 4 | 4 | You've already forked SwarmClaw → OneMind Web. Proven path. |
| **Robotics/drones** | 5 | 1 | 1 | Only OpenClaw has RosClaw (ROS 2 bridge). |
| **Smart glasses** | 5 | 1 | 1 | Only OpenClaw has VisionClaw (Meta Ray-Ban). You already have this fork. |
| **Home automation** | 5 | 4 | 2 | OpenClaw + Hermes both have Home Assistant. A0 requires custom API work. |
| **Open source** | 5 | 5 | 5 | All MIT. |
| **Multiple interfaces** | 5 | 4 | 2 | OpenClaw: 28+ platforms. Hermes: 14+. A0: Web UI only. |
| **Ecosystem / skills** | 5 | 4 | 2 | OpenClaw: 47K skills. Hermes: can use ClawHub + others. A0: ~50 plugins. |
| **Self-improvement / learning** | 4 | 5 | 3 | Both have learning loops. Hermes: built-in + user modeling (Honcho). OpenClaw: via top ClawHub skills + built-in Dreaming. Hermes edge: native integration. |
| **AI research / model training** | 1 | 5 | 1 | Only Hermes (Nous Research lab). |
| **Privacy / local models** | 4 | 4 | 5 | A0: Docker sandbox + SearXNG private search. All support Ollama. |
| **TOTAL** | **53/60** | **46/60** | **33/60** | — |

---

## Summary: What Each Platform Is Best At

### OpenClaw — The Ecosystem Giant
- **Best for:** Life operating system with maximum reach (28+ platforms, 47K skills, robotics, glasses)
- **Strength:** Largest community, most integrations, proven at your scale (already running)
- **Weakness:** Self-improvement via skills rather than built-in, no cron scheduling, no research pipeline
- **Your investment:** 9 agents, 30 skills, VPS deployment, SwarmClaw fork → OneMind Web, VisionClaw fork → OneMind Mobile

### Hermes Agent — The Learning Machine
- **Best for:** Developers/researchers who want an agent that gets smarter over time
- **Strength:** Self-improving skills, RL training, 20+ providers, 6 execution backends, Nous Research backing
- **Weakness:** No robotics/glasses, smaller ecosystem (but uses OpenClaw's), CLI-first (no Web UI)
- **Migration path:** `hermes claw migrate` imports everything from OpenClaw

### Agent Zero — The Tinkerer's Playground
- **Best for:** Power users who want to build agents from scratch with total control
- **Strength:** "Computer as a tool" philosophy, A2A protocol, privacy focus, crypto governance
- **Weakness:** No messaging gateway (Web UI only), tiny ecosystem (~50 plugins), no robotics/glasses
- **Unique:** A0T governance token, Skool community, Czech company structure

---

## Recommendation (Data-Based)

**For the OneMind use case specifically:**

1. **OpenClaw remains the right primary platform.** The robotics bridge (RosClaw), smart glasses (VisionClaw), 28+ messaging platforms, and 47K skills are capabilities neither Hermes nor Agent Zero offer. You've already built on it (OneMind Web, OneMind Mobile). The ecosystem is 13x larger than the next competitor.

2. **Hermes Agent is worth watching closely.** Built by Nous Research (serious AI lab), growing explosively (26.9K stars in 2 months), and has a direct migration path from OpenClaw. The self-improving skill loop is genuinely novel. If it develops robotics/hardware integrations, it could become a compelling alternative. The fact that it ALREADY reads from ClawHub means it's designed to coexist with or succeed OpenClaw.

3. **Agent Zero is the wrong fit for OneMind.** No messaging gateway (critical for a life OS), tiny ecosystem, no robotics/hardware, and a crypto token governance model adds complexity without value for a personal framework. The A2A protocol support is notable but insufficient to overcome the other gaps.

**Bottom line:** OpenClaw now. Keep Hermes on the radar. Skip Agent Zero.

---

*Research conducted April 6, 2026. All numbers verified against live GitHub repos and official websites.*


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]
