# The complete OpenClaw ecosystem inventory

**OpenClaw is the fastest-growing open-source project in GitHub history — a personal AI assistant with 236,000+ stars, 5,700+ community skills, 115,000 Discord members, and an explosion of academic research — all in under four months.** This inventory catalogs every significant repo, plugin, integration, community, and resource across the ecosystem so Master Delacruz can systematically extend One Mind OS across all four pillars. The ecosystem spans voice/vision agents, IoT/home automation, NATS messaging, LLM orchestration, and a rapidly maturing developer community producing 8+ academic papers in just three weeks.

OpenClaw (formerly Clawdbot → Moltbot → OpenClaw) was created by Austrian developer Peter Steinberger in November 2025. It operates as a local-first AI gateway connecting to messaging platforms (WhatsApp, Telegram, Slack, Discord, Signal, iMessage, Teams, and more) with an extensible skill/plugin architecture. On February 14, 2026, Steinberger announced he was joining OpenAI, and the project would transition to an independent OpenClaw Foundation. The ClawHub registry hosts **5,705+ community-built skills**, with the curated awesome lists filtering these to ~2,868 vetted entries.

---

## Official OpenClaw organization repositories

The `openclaw` GitHub org maintains **20 repositories** that form the project's core:

| Repo | URL | Description | Stars |
|------|-----|-------------|-------|
| **openclaw/openclaw** | https://github.com/openclaw/openclaw | Core gateway runtime — personal AI assistant with multi-channel inbox, voice, canvas, multi-agent routing | **236,312** |
| **openclaw/clawhub** | https://github.com/openclaw/clawhub | Skill Directory — the registry where all community skills are published and discovered | — |
| **openclaw/lobster** | https://github.com/openclaw/lobster | OpenClaw-native workflow shell — typed, local-first "macro engine" for composable pipelines and safe automations | ~1,500 |
| **openclaw/skills** | https://github.com/openclaw/skills | Archived backup of all versions of all skills ever published to ClawHub — rollback and audit purposes | 1,600 |
| **openclaw/nix-openclaw** | https://github.com/openclaw/nix-openclaw | Packages OpenClaw for Nix/NixOS | — |
| **openclaw/clawdinators** | https://github.com/openclaw/clawdinators | Declarative infra + NixOS modules for CLAWTINATOR hosts | — |
| **openclaw/openclaw-ansible** | https://github.com/openclaw/openclaw-ansible | Automated, hardened installation with Tailscale VPN, UFW, Fail2ban, Docker isolation | 420 |
| **openclaw/acpx** | https://github.com/openclaw/acpx | Headless CLI client for Agent Client Protocol (ACP) sessions — bridges Codex, Claude, Gemini, OpenCode, Pi | — |

**One Mind OS value**: The official repos provide the foundational runtime (core gateway), workflow engine (Lobster), agent communication layer (acpx/ACP), and infrastructure-as-code (Ansible, Nix). These directly serve the **Innovative Technology** pillar as the backbone of any autonomous agent deployment.

---

## Awesome lists and curated collections

The community has produced multiple massive curated lists, with the largest reaching **195,000+ stars** — making it itself one of the fastest-growing GitHub repos ever:

| Repo | URL | Description | Stars |
|------|-----|-------------|-------|
| **rohitg00/awesome-openclaw** | https://github.com/rohitg00/awesome-openclaw | The definitive curated collection — hosting guides, cost comparisons, security hardening, 80+ ecosystem projects, powers openclawsearch.com | **195,000+** |
| **VoltAgent/awesome-openclaw-skills** | https://github.com/VoltAgent/awesome-openclaw-skills | Curated catalog of 2,868 skills filtered from 5,705+ on ClawHub, organized by category with security warnings | — |
| **hesamsheikh/awesome-openclaw-usecases** | https://github.com/hesamsheikh/awesome-openclaw-usecases | Real-life use cases — voice SMS, family calendars, multi-agent teams, trading, memory dashboards, content pipelines | — |
| **SamurAIGPT/awesome-openclaw** | https://github.com/SamurAIGPT/awesome-openclaw | Resources, tools, skills, tutorials — covers ClawHub, MCP support (13,000+ servers), phone integrations | — |
| **thewh1teagle/awesome-openclaw** | https://github.com/thewh1teagle/awesome-openclaw | "Agent-perspective" list organized by Work, Hardware, Skills, Infrastructure, Social, Bounties, Trading | — |
| **BankrBot/openclaw-skills** | https://github.com/BankrBot/openclaw-skills | Financial/crypto skill collection — Polymarket, crypto trading, DeFi operations | — |

**One Mind OS value**: These lists are the primary discovery layer for the entire ecosystem. The rohitg00 list alone catalogs 80+ projects. For **Generational Entrepreneurship**, the use-case list demonstrates how OpenClaw agents generate revenue (trading, content, business automation). For **Innovative Technology**, the skills collection provides 2,868 pre-built capabilities to integrate.

---

## Voice and vision agent repos

OpenClaw has built-in **Voice Wake + Talk Mode** with ElevenLabs, but the community has expanded voice/vision capabilities dramatically. This is one of the ecosystem's most active areas with **7+ distinct projects**:

| Repo | URL | Description | Category |
|------|-----|-------------|----------|
| **sseanliu/VisionClaw** | https://github.com/sseanliu/VisionClaw | Real-time AI assistant for **Meta Ray-Ban smart glasses** — Gemini Live voice + vision over WebSocket, OpenClaw as action layer with 56+ tools. iOS (Swift) and Android (Kotlin) apps. WebRTC POV streaming at 2.5 Mbps | Vision + Voice |
| **Purple-Horizons/openclaw-voice** | https://github.com/Purple-Horizons/openclaw-voice | Open-source browser-based voice chat — Whisper STT (local via faster-whisper) + ElevenLabs or Chatterbox TTS. Self-hosted, private, no app install. Website: openclawvoice.com | Voice |
| **openserv-labs/openclaw-voice-avatar** | https://github.com/openserv-labs/openclaw-voice-avatar | **LiveKit + LemonSlice** lip-synced avatar — Deepgram transcription → OpenClaw Gateway → ElevenLabs TTS → LemonSlice avatar generation. ~1-2s latency. Posted on Hacker News | Voice + Vision |
| **SudarshanSuryaprakash/voice-assistant** | https://github.com/SudarshanSuryaprakash/voice-assistant | "Jupiter Voice" — fully local voice assistant. Apple Silicon optimized (MLX for 10x faster Whisper), OpenWakeWord for custom wake phrases, Kokoro ONNX TTS | Voice (local) |
| **Omni-NexusAI/agent-voice-messaging** | https://github.com/Omni-NexusAI/agent-voice-messaging | Modular voice messaging with swappable STT/TTS providers (faster-whisper, OpenAI, Google, Kokoro, ElevenLabs). Docker-ready with GPU/CUDA. Works with OpenClaw, LangChain, Agent Zero | Voice |
| **omeriko9/HeyClawy** | https://github.com/omeriko9/HeyClawy | Voice-first **ESP32** physical companion — wake word support (WakeNet, "Hey Jarvis" default), display or screenless modes, local web settings. Bridges hardware → OpenClaw | Voice + IoT |
| **elevenlabs-voices** (ClawHub skill) | openclaw/skills registry | Production-ready ElevenLabs skill: 18 curated voice personas, 32 languages, streaming, SFX generation, voice design tool, pronunciation dictionary | Voice |

VisionClaw is the standout project here — it gives OpenClaw agents **eyes and ears through smart glasses**, running Gemini Live for real-time perception while using OpenClaw as the execution backend. The openclaw-voice-avatar project enables **LiveKit-powered avatar agents** with lip-synced faces, enabling a full embodied AI experience.

**One Mind OS value**: Voice and vision are critical for the **Holistic Performance** pillar (hands-free daily briefings, health check-ins, ambient awareness) and **Innovative Technology** (embodied agents via VisionClaw's smart glasses integration). The LiveKit avatar project directly enables the kind of personalized AI presence that defines a "life operating system."

---

## IoT and home automation integrations

The ecosystem has **strong Home Assistant integration through multiple paths** plus a growing ESP32 hardware ecosystem:

| Repo | URL | Description | Protocol |
|------|-----|-------------|----------|
| **techartdev/OpenClawHomeAssistant** | https://github.com/techartdev/OpenClawHomeAssistant | Home Assistant **add-on** — runs OpenClaw inside HAOS with AI Gateway, Web Terminal, Assist Pipeline (conversation agent via OpenAI-compatible API), browser automation, persistent storage | HA Add-on |
| **Djelibeybi/openclaw_conversation** | https://github.com/Djelibeybi/openclaw_conversation | HA **custom component** — uses OpenClaw as AI backend for conversation and task platforms. Leverages mcporter and ha-mcp for device control. Can be set as Assist pipeline conversation agent | HA Component |
| **homeassistant skill** (ClawHub) | HA Community forum thread #990911 | OpenClaw skill for HA via **REST API** — full entity control (lights, switches, locks, thermostats), automation triggers, sensor queries, historical data. Supports MQTT for real-time events | REST + MQTT |
| **omarshahine/HomeClaw** | https://github.com/omarshahine/HomeClaw | **Apple HomeKit** MCP control — lights, locks, thermostats, scenes. Mac Catalyst app with direct HomeKit API. OpenClaw plugin registers HomeKit skill on gateway. CLI tool + TestFlight | HomeKit/MCP |
| **M64GitHub/WireClaw** | https://github.com/M64GitHub/WireClaw | AI agent on a **$5 ESP32** — 20 tools, Telegram bot, NATS interfaces. Rule engine with multi-step chains. Persistent local automation. Cross-device messaging via NATS. 30ms latency. Web flash at wireclaw.io | ESP32 + NATS |
| **memovai/mimiclaw** | https://github.com/memovai/mimiclaw | "MimiClaw" — OpenClaw-like agent on **$5 ESP32-S3** in pure C. No OS/Linux/Node.js required. ReAct agent loop, persistent memory (SOUL.md, USER.md, MEMORY.md), Telegram integration | ESP32-S3 |
| **hrwtech/openclaw-esp32** | https://github.com/hrwtech/openclaw-esp32 | OpenClaw on **Arduino ESP32** devices — port/adaptation for embedded hardware | ESP32/Arduino |

Home Assistant is the primary integration path, with three distinct approaches: an all-in-one HAOS add-on, a conversation component for Assist pipelines, and a REST API skill. Through Home Assistant, OpenClaw gains access to **2,500+ device types** across Zigbee, Z-Wave, MQTT, Wi-Fi, Bluetooth, Thread, and Matter protocols. The ESP32 ecosystem (WireClaw, MimiClaw, HeyClawy) enables **$5 hardware agents** that can run autonomously on bare metal.

**One Mind OS value**: IoT integrations directly serve the **Holistic Performance** pillar (automated environment, health monitoring, circadian lighting) and **Legacy Evolution** (home management that adapts to family routines). The WireClaw NATS integration enables multi-device coordination across an entire home — exactly the infrastructure needed for an AI life operating system.

---

## NATS messaging and agent coordination repos

NATS adoption in the OpenClaw ecosystem is **emerging but architecturally significant**, with two key implementations and OpenClaw's own built-in multi-agent coordination system:

| Repo | URL | Description |
|------|-----|-------------|
| **jbold/exoclaw** | https://github.com/jbold/exoclaw | WASM-sandboxed AI agent runtime in **Rust** — uses **NATS Message Bus** as core routing infrastructure with subject pattern `exoclaw.channel.account.peer`. WebSocket → axum Gateway → Auth → JSON-RPC → Session Router → Agent Runner. Graceful fallback to local-only routing if NATS unavailable |
| **M64GitHub/WireClaw** | https://github.com/M64GitHub/WireClaw | ESP32 devices communicate over **NATS** for cross-device messaging. OpenClaw skill drops in to discover/control all WireClaws on the network. Pattern: OpenClaw → NATS → WireClaw → hardware. **30ms latency** |
| **OpenClaw built-in** | github.com/openclaw/openclaw | Native multi-agent: `sessions_send` (message any session), `sessions_spawn` (child agents with restricted tools), hierarchical sessions, cross-context messaging (#22725). RFC-0001 (Discussion #10036) proposes formal Agent Teams with SendMessage/broadcast/shutdown |
| **openclaw/lobster** | github.com/openclaw/lobster | Workflow shell supporting `openclaw.invoke` for agent-send and LLM-task orchestration. Sub-lobster steps with loop support. "GitHub Actions for OpenClaw" |
| **Agent Protocol Skill** | lobehub.com/skills/openclaw-skills-agent-protocol | Event-bus communication layer — pub/sub agent-to-agent messaging, automated workflows, multi-step orchestration chains (Research→Summary→Notification) |

A DEV Community post by @ggondim documented considering NATS/Redis Streams for external event bus orchestration but found OpenClaw's built-in `sessions_send` and agent coordination sufficient for most use cases. The key insight: **NATS is most valuable for cross-device physical coordination** (WireClaw) and **high-security sandboxed runtimes** (ExoClaw), while OpenClaw's native system handles digital agent-to-agent communication well.

**One Mind OS value**: NATS serves the **Innovative Technology** pillar as the coordination backbone for distributed agents across devices. ExoClaw's WASM-sandboxed architecture addresses security concerns. WireClaw demonstrates the **30ms latency** pub/sub pattern needed for real-time IoT agent coordination. For a life operating system spanning multiple rooms, devices, and contexts, NATS provides the reliable message bus.

---

## LLM orchestration and multi-agent frameworks

OpenClaw positions itself not as a framework (like LangChain or CrewAI) but as a **runtime** — making it complementary rather than competitive with orchestration libraries. The community has built explicit bridges:

| Repo | URL | Description |
|------|-----|-------------|
| **tisu19021997/langclaw** | https://github.com/tisu19021997/langclaw | "LangChain + OpenClaw = Langclaw" — Python framework compiling to real **LangGraph CompiledStateGraph**. RBAC, scheduled tasks, persistent memory, subagent delegation. Multi-channel from day one (Telegram, Discord, WebSocket). FastAPI-style decorators |
| **jovanSAPFIONEER/Network-AI** | https://github.com/jovanSAPFIONEER/Network-AI | Framework-agnostic orchestrator connecting **12 frameworks**: LangChain, AutoGen, CrewAI, OpenAI Assistants, LlamaIndex, Semantic Kernel, Haystack, DSPy, Agno, MCP, OpenClaw, custom adapters. Shared blackboard, AES-256 security, hallucination detection |
| **CrewClaw** | https://www.crewclaw.com/ | Management platform for teams running 3+ agents, combining OpenClaw with CrewAI or LangChain agents. Visual dashboard, workflow builder, monitoring |
| **freema/openclaw-mcp** | https://github.com/freema/openclaw-mcp | MCP server bridging **Claude.ai to self-hosted OpenClaw** with OAuth 2.1 — "AI assistant orchestrating another AI assistant" |
| **lunarpulse/openclaw-mcp-plugin** | https://github.com/lunarpulse/openclaw-mcp-plugin | Plugin enabling OpenClaw to connect to any MCP server via Streamable HTTP transport |
| **Virtual-Protocol/openclaw-acp** | https://github.com/Virtual-Protocol/openclaw-acp | Virtuals Protocol ACP integration — agent-to-agent **commerce** with auto-provisioned wallets on Base chain. Browse marketplace, create bounties, tokenize agents |
| **a3tai/openclaw-go** | https://pkg.go.dev/github.com/a3tai/openclaw-go/acp | Go package implementing ACP for bridging IDE clients ↔ OpenClaw Gateway. JSON-RPC 2.0 over NDJSON |
| **caprihan/openclaw-n8n-stack** | https://github.com/caprihan/openclaw-n8n-stack | "Personal AI operating system in a box" — Docker Compose stack combining **OpenClaw + n8n**. Shared network, webhook-based communication. User → OpenClaw (reasoning) → webhook → n8n (execution) → External APIs |
| **Mem0 integration** | https://docs.mem0.ai/integrations | Mem0 (persistent memory layer) lists OpenClaw as first-class integration alongside CrewAI, AutoGen, Agno, Camel AI, OpenAI Agents SDK |

OpenClaw supports **MCP natively** (@modelcontextprotocol/sdk@1.25.3) and **ACP** (@agentclientprotocol/sdk 0.13.x), giving it protocol-level interoperability with the broader agent ecosystem. The n8n integration is particularly well-documented, with multiple community guides describing OpenClaw as the "reasoning layer" while n8n handles "deterministic execution." Issue #28511 proposes expanding acpx to support **any standard ACP agent** from a registry of 19+ agents, which would make OpenClaw the first messaging-platform-based universal ACP client.

**One Mind OS value**: The orchestration layer is the nervous system of a life operating system. Langclaw bridges LangGraph's powerful graph-based workflows. Network-AI's 12-framework support means One Mind OS can dispatch tasks to the best agent framework for each domain. The n8n stack provides the "personal AI OS in a box" pattern that maps directly to One Mind OS's architecture. Virtual Protocol's agent commerce integration enables the **Generational Entrepreneurship** pillar through autonomous economic agents.

---

## OpenMind OM1 and robotics

| Repo | URL | Description | Stars |
|------|-----|-------------|-------|
| **OpenMind/OM1** | https://github.com/OpenMind/OM1 | Open-source modular AI runtime for robots — supports quadrupeds, humanoids, wheeled robots, drones. ROS2, Zenoh, CycloneDDS. Preconfigured for Unitree Go2/G1, TurtleBot, Ubtech. SLAM, LIDAR, Nav2, Gazebo. Python/MIT | **~2,500** |
| **OM1-video-processor** | OpenMind org | Docker-based real-time video streaming, face recognition, audio for OM1 robots | — |
| **OM1-ros2-sdk** | OpenMind org | ROS 2 SLAM package for Unitree Go2 | — |
| **om-yolo-ws** | OpenMind org | WebSocket-based remote YOLO object detection | — |
| **humanoid_teleoperate** | OpenMind org | XR device teleoperation for Unitree humanoid robots | — |
| **FABRIC** | OpenMind org | Decentralized coordination layer for robot identity, trust, and data sharing | — |

OM1 was founded by Stanford professor Jan Liphardt with **$20M funding** from Pantera Capital, Ribbit, and Coinbase Ventures. It recently launched a **Robot App Store** with partners including UBtech, Agibot, Deep Robotics, Fourier, Booster, Dobot, LimX, and Magic Lab.

**No direct OpenClaw ↔ OM1 integration exists yet.** Both projects share architectural parallels: open-source, MIT-licensed, modular, agent-oriented, LLM-powered. Both support similar providers (OpenAI, Gemini, Claude, Ollama). A theoretical bridge via MCP or ACP is architecturally plausible — OpenClaw handles digital orchestration and messaging while OM1 handles physical robotics execution. The OpenClaw `robot` skill (openclaw/skills/robot) provides Arduino/ESP32/ROS2 assistance, creating a conceptual connection point.

**One Mind OS value**: OM1 represents the **physical embodiment** layer for a life operating system. Once an OpenClaw → MCP/ACP → OM1 bridge exists, One Mind OS could coordinate digital agents and physical robots under a single orchestration layer. This directly serves **Innovative Technology** (autonomous systems) and **Legacy Evolution** (physical assistants for family members).

---

## Community-built tools, dashboards, and extensions

Beyond skills, the community has built substantial standalone tools:

| Repo | URL | Description |
|------|-----|-------------|
| **grp06/openclaw-studio** | https://github.com/grp06/openclaw-studio | Clean web dashboard — agent creation, cron jobs, exec approvals, tool call visualization, thinking traces |
| **AidanPark/openclaw-android** | https://github.com/AidanPark/openclaw-android | Run on Android with single command — no proot, no Linux. Patches Termux compatibility directly |
| **ComposioHQ/secure-openclaw** | https://github.com/ComposioHQ/secure-openclaw | Hardened 24/7 assistant with Docker isolation, persistent memory, cron scheduling, 500+ app integrations via Composio |
| **cloudflare/moltworker** | https://github.com/cloudflare/moltworker | Run on **Cloudflare Workers/Sandbox** — fully managed, always-on. ~$5/month |
| **MemTensor/MemOS-Cloud-OpenClaw-Plugin** | https://github.com/MemTensor/MemOS-Cloud-OpenClaw-Plugin | MemOS Cloud plugin — long-term memory with recall/capture strategies, global or per-conversation |
| **supermemoryai/openclaw-supermemory** | https://github.com/supermemoryai/openclaw-supermemory | Supermemory — auto-recall injects relevant memories before every AI turn, auto-capture stores conversations |
| **lekt9/openclaw-foundry** | https://github.com/lekt9/openclaw-foundry | "The forge that forges itself" — self-writing meta-extension that observes workflows → crystallizes patterns → upgrades itself |
| **prompt-security/clawsec** | https://github.com/prompt-security/clawsec | 8-skill security suite: NVD advisory feed, ClawHub reputation checker, audit watchdog, soul-guardian integrity |
| **Next-Kick/openclaw-hardened-ansible** | https://github.com/Next-Kick/openclaw-hardened-ansible | Tier 3+ security — rootless Podman, Squid proxy egress filtering, LiteLLM credential brokering (OpenClaw never sees real API keys) |
| **clawdotnet/openclaw.net** | https://github.com/clawdotnet/openclaw.net | .NET NativeAOT implementation with Node.js JSON-RPC bridge for JS plugin compatibility. 18+ native C# tools |
| **bfzli/clawhost** | https://github.com/bfzli/clawhost | Open-source hosting platform — multi-cloud (Hetzner, DigitalOcean, Vultr), web + mobile + desktop apps, drag-and-drop workflows |
| **sunnoy/openclaw-plugin-wecom** | https://github.com/sunnoy/openclaw-plugin-wecom | Enterprise WeChat (WeCom) plugin — streaming output, per-user agents, group chat with @mention, AES-256 image decryption |

Additional noteworthy projects mentioned in awesome lists but not fully detailed: **ClawBridge** (mobile dashboard for real-time agent monitoring), **Clawdstrike** (runtime security/swarm detection for agent fleets), **Vault-0** (agent security + x402 wallet), **ClawRouter** (smart LLM router with 78% cost savings), and **CEO Desk** (local-first AI agent office simulator).

---

## Deployment and infrastructure repos

| Repo | URL | Description |
|------|-----|-------------|
| **phioranex/openclaw-docker** | https://github.com/phioranex/openclaw-docker | Pre-built Docker image — auto-builds daily, checks releases every 6 hours |
| **coollabsio/openclaw** | https://github.com/coollabsio/openclaw | Multi-arch (amd64/arm64) Docker images with nginx reverse proxy, basic auth, 10+ AI providers |
| **khal3d/openclaw** | https://github.com/khal3d/openclaw | Docker Compose + **Kubernetes Helm chart** with persistent volumes and health checks |
| **andreesg/openclaw-docker-config** | https://github.com/andreesg/openclaw-docker-config | Docker config with workspace backup sidecar (auto-push to private GitHub via cron) |
| **andreesg/openclaw-terraform-hetzner** | https://github.com/andreesg/openclaw-terraform-hetzner | **Terraform** modules for Hetzner Cloud — remote S3 state, cloud-init, VPS provisioning. Referenced in official docs |
| **jeffweisbein/openclaw-starter-kit** | https://github.com/jeffweisbein/openclaw-starter-kit | "Battle-tested" workspace template — SOUL.md personality, multi-agent squad, safety policies with daily caps |
| **zach-highley/openclaw-starter-kit** | https://github.com/zach-highley/openclaw-starter-kit | "The 24-Hour AI Butler" — self-healing, multi-model system with 355-item security audit, launchd/systemd reliability |
| **MaximusCarapax/openclaw-starter-kit** | https://github.com/MaximusCarapax/openclaw-starter-kit | Bootstrap template with Aider for coding, RAG with free Gemini embeddings, DeepSeek, Gmail integration |

---

## Research papers and benchmarks

The OpenClaw ecosystem generated **8+ academic papers within three weeks of Moltbook's launch** — an extraordinary pace. Key papers and benchmarks:

| Paper | URL | Focus |
|-------|-----|-------|
| "From Assistant to Double Agent" (Wang et al.) | https://arxiv.org/abs/2602.08412 | **PASB benchmark** — end-to-end security evaluation across personalized scenarios, attack types, tool capabilities |
| "A Trajectory-Based Safety Audit" (Chen et al.) | https://arxiv.org/abs/2602.14364 | Safety evaluation across 6 risk dimensions using AgentDoG-Qwen3-4B judge, 34 canonical test cases |
| "OpenClaw Agents on Moltbook: Risky Instruction Sharing" (Manik, Wang) | https://arxiv.org/abs/2602.02625 | First large-scale study: 39,026 posts, 5,712 comments, 14,490 agents. Introduces Action-Inducing Risk Score |
| "Collective Behavior of AI Agents" | https://arxiv.org/html/2602.09270 | 369,209 posts, 3M+ comments, 46,690 agents over 12 days — collective dynamics |
| "When OpenClaw Agents Teach Each Other" (Chen et al.) | https://arxiv.org/abs/2602.14477 | Peer learning among 2.4M+ agents — taxonomy of validation (22%), knowledge extension (18%), application (12%) |
| "OpenClaw, Moltbook, and ClawdLab" (Weidener et al.) | https://arxiv.org/abs/2602.19810 | Multivocal review of 6 papers in 14 days. Documents 5 recurring architectural patterns, introduces ClawdLab |
| "ClawMobile" | https://arxiv.org/html/2602.22942 | Smartphone-native agentic systems — proposes unified mobile runtime from 3 community extensions |
| **HKUDS/ClawWork benchmark** | https://github.com/HKUDS/ClawWork | Economic survival: agents start with $10, complete 220 professional tasks across 44 sectors. Top models hit **$1,500+/hr** |

Security research from major firms includes Cisco's malicious skill analysis, Bitdefender's enterprise advisory (documenting ~900 malicious skills), Kaspersky's audit (512 vulnerabilities, 8 critical), and Bitsight's scan of **42,900+ exposed instances** across 82 countries. Five CVEs have been assigned and patched (CVE-2026-25253 CVSS 8.8 being the most severe).

---

## Community platforms and developer resources

The community is massive and highly active across multiple platforms:

| Platform | URL | Size | Description |
|----------|-----|------|-------------|
| **Discord** ("Friends of the Crustacean 🦞🤝") | https://discord.gg/clawd | **~115,000 members** | Primary real-time community. Support, announcements, PR discussions, bot "Krill" |
| **Reddit** (r/openclaw) | https://reddit.com/r/openclaw | Active subreddit | Guides, troubleshooting, project showcases. Origin of "OpenClaw 101" guide |
| **X Community** | https://x.com/i/communities/2013441068562325602 | **24,400+ members** | Builder community — configs, workflow snippets, results |
| **Moltbook** | Moltbook forum/social | **1.6M+ AI agents** | AI-only social network — Reddit-style, 15,800+ subcommunities. Has its own Wikipedia article |
| **GitHub Discussions** | github.com/openclaw/openclaw/discussions | Thousands of threads | Official Q&A, skill showcases, feature requests, RFCs |
| **Clawstin** (Austin meetup) | https://clawstin.com/ | In-person | OpenClaw meetups in Austin, TX — demos, hacks, workshops |
| **Silicon Valley Meetup** | meetup.com/meetup-group-jrvvvlrp | In-person | "The only OpenClaw meetup in the heart of Silicon Valley" |
| **AI Tinkerers SF** (Unhackathon) | sf.aitinkerers.org/p/openclaw-unhackathon | 40-person build sessions | Curated heads-down build sessions. No vendors, no prizes — just building |
| **Luma Events Calendar** | https://luma.com/claw | Global | Community events calendar for meetups worldwide |

---

## Media, newsletters, and educational content

| Resource | URL | Type |
|----------|-----|------|
| **Lex Fridman Podcast #491** | https://lexfridman.com/peter-steinberger-transcript/ | 3-hour interview with creator — "the age of the lobster" |
| **OpenClaw Newsletter** (official) | https://buttondown.com/openclaw-newsletter/ | Daily ecosystem updates via Buttondown |
| **OpenClaws.io Blog** | https://openclaws.io/blog/ | Third-party comprehensive news — releases, tutorials, security, trading |
| **OpenClaw Radar** | https://openclawradar.com/ | Community articles and guides aggregator |
| **freeCodeCamp Tutorial** | https://www.freecodecamp.org/news/openclaw-full-tutorial-for-beginners/ | 1-hour video course on YouTube — installation through advanced skills |
| **DataCamp Tutorial** | https://www.datacamp.com/tutorial/moltbot-clawdbot-tutorial | Written guide — setup, WhatsApp, file searching, cost estimates |
| **Creator Economy / Nat Eliason** | creatoreconomy.so | "Felix" bot given $1K, made $14,718 launching products |
| **Jason Calacanis / TWIST** | calacanis.substack.com | 20-day OpenClaw obsession — claims 20% task offload at his 20-person firm |
| **Nature** | nature.com/articles/d41586-026-00370-w | "OpenClaw AI chatbots are running amok" — Moltbook scientific interest |
| **CNBC** | cnbc.com | Comprehensive overview of rise and global adoption |
| **CoClaw** | https://coclaw.com/ | Community support — free guides, config generator, platform checklists |
| **LearnOpenClaw** | https://learnopenclaw.com/ | Learning resource with channel-specific setup guides |
| **centminmod/explain-openclaw** | https://github.com/centminmod/explain-openclaw | Multi-AI documentation with architecture deep dives and security audits |

---

## How the ecosystem maps to One Mind OS pillars

The four pillars of One Mind OS each have clear integration points across this ecosystem:

**Holistic Performance** draws from voice agents (Jupiter Voice, openclaw-voice for hands-free interaction), IoT integrations (Home Assistant for environment optimization, ESP32 sensors for health data), memory plugins (MemOS, Supermemory for persistent context), and daily briefing use-cases from the awesome-openclaw-usecases collection. The family calendar aggregation and health tracking projects documented in hesamsheikh's list map directly here.

**Legacy Evolution** leverages the multi-agent coordination layer (Lobster workflows, sessions_spawn for specialized family-member agents), the starter kits (SOUL.md personality templates that encode personal values), persistent memory systems, and the n8n automation stack for systematizing family operations. The Moltbook research on peer learning between agents provides architectural patterns for knowledge transfer.

**Generational Entrepreneurship** is served by the trading assistant (molt-bot/openclaw-trading-assistant), the ClawWork economic benchmark (proving agents can generate $1,500+/hr equivalent), Virtual Protocol's agent commerce layer (tokenizing agent services on-chain), the n8n stack for business process automation, and Nat Eliason's documented case of an agent earning $14,718 from a $1,000 starting position.

**Innovative Technology** encompasses everything: the core gateway architecture, MCP/ACP protocol bridges to 12+ AI frameworks (Network-AI), ExoClaw's WASM-sandboxed NATS runtime, OM1's robotics layer (pending bridge), the .NET implementation (openclaw.net), Terraform/Kubernetes deployment infrastructure, and the security hardening stack (clawsec, hardened-ansible) that makes production deployment viable.

## Conclusion

The OpenClaw ecosystem is **not a single tool but a platform economy** — 236K GitHub stars, 5,700+ skills, 115K Discord members, 1.6M Moltbook agents, and 8+ academic papers in three weeks. The most strategically valuable repos for One Mind OS are the **n8n stack** (reasoning + execution in one box), **WireClaw** (NATS-based IoT agent coordination at 30ms latency), **VisionClaw** (smart glasses perception layer), **Network-AI** (12-framework orchestration bridge), and **openclaw-voice-avatar** (LiveKit embodied agents). The ecosystem's biggest gap — and One Mind OS's biggest opportunity — is the **missing OpenClaw ↔ OM1 bridge** that would unify digital and physical agents under a single life operating system. Building that MCP/ACP connector would position One Mind OS as the first project to span messaging, IoT, voice/vision, LLM orchestration, and robotics in a single coherent architecture.