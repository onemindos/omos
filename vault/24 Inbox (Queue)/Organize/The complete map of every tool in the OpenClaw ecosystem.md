# The complete map of every tool in the OpenClaw ecosystem

OpenClaw — the open-source personal AI agent framework with **312,000+ GitHub stars** — has spawned one of the most explosive tool ecosystems in open-source history. In barely 60 days since going viral in late January 2026, **100+ distinct repositories, dashboards, plugins, and commercial products** have emerged to extend, monitor, secure, and simplify OpenClaw deployments. This report catalogs every known enhancement tool across 12 categories, from full command centers to single-purpose plugins, covering what exists, what works, and what's still missing.

Created by Austrian developer Peter Steinberger and now transitioning to an OpenAI-backed open-source foundation, OpenClaw runs a local Node.js Gateway that connects LLMs to 25+ messaging platforms. Its architecture — Gateway control plane, WebSocket network, plugin lifecycle hooks, and Markdown-based memory — has become a substrate on which a massive third-party ecosystem now builds. The official OpenClaw GitHub organization hosts **23 repositories**, the ClawHub skill registry contains **13,700+ published skills**, and community projects span dashboards, monitoring tools, memory engines, security suites, trading bots, and hardware appliances.

---

## Full dashboards and command centers now number over fifteen

The most crowded category in the ecosystem is web-based dashboards. At least **15 independent projects** provide visual command surfaces for OpenClaw, ranging from lightweight monitoring panels to full enterprise orchestration platforms.

**OpenClaw Mission Control** (abhi1693/openclaw-mission-control) is the most ambitious. It provides centralized operations and governance for teams: organizations, board groups, tasks with tags, approval-driven governance workflows, gateway management for distributed environments, and a unified API. It deploys via Docker Compose or local Node.js 22+ with a one-line install script. The project explicitly targets platform teams running OpenClaw in self-hosted or internal environments and is under active development with the caveat to "validate and harden your configuration before production use."

**SwarmClaw** (swarmclawai/swarmclaw | swarmclaw.ai) goes further into orchestration territory. It connects **15 built-in providers** (Claude Code CLI, Codex CLI, OpenCode, Anthropic, OpenAI, Gemini, DeepSeek, Groq, Ollama, and more), offers LangGraph-powered multi-agent workflows with fan-out/merge, a Kanban task board, runtime skills with SKILL.md import, agent-linked Solana/Ethereum crypto wallets, and connectors for Discord, Slack, Telegram, WhatsApp, Teams, and Matrix. It authenticates to OpenClaw gateways via ed25519 device identity and manages multiple gateways from a single control plane. Built with Next.js and SQLite (WAL mode), it runs self-hosted on port 3456.

**Paperclip** (paperclipai/paperclip | paperclip.ing) has attracted **23,354 GitHub stars** and takes a radically different angle: orchestrating teams of AI agents to run entire virtual companies. It features org charts, budgets, governance, goal alignment, hierarchical task management, and full tool-call tracing. Its **first-class OpenClaw gateway adapter** (v0.3.0) supports strict SSE streaming, device-key pairing, and invite-based onboarding. It also ships adapters for Claude Code, Codex, Cursor, and Pi. The system uses embedded PostgreSQL and supports multi-company deployments with data isolation.

**EasyClaw** (gaoyangz77/easyclaw | easyclaw.com) targets a fundamentally different audience: non-engineers. Described as "OpenClaw is the engine; EasyClaw is the cockpit," it wraps OpenClaw into an Electron desktop app with a React management panel. Users write natural-language rules that compile to policies, guards, or skills without restarts. It supports **20+ LLM providers**, visual channel configuration, token usage tracking, speech-to-text, a built-in skills marketplace, zero-restart hot reload for API key changes, and system keychain secret storage. The trade-off is unsigned app warnings on macOS/Windows and potential lag behind upstream OpenClaw releases.

| Dashboard | GitHub | Standout Feature | Tech Stack |
|-----------|--------|------------------|------------|
| **OpenClaw Mission Control** (abhi1693) | abhi1693/openclaw-mission-control | Enterprise governance + approvals | Docker/Node.js |
| **SwarmClaw** | swarmclawai/swarmclaw | 15 providers + crypto wallets | Next.js/SQLite |
| **Paperclip** | paperclipai/paperclip (23.3k ⭐) | Virtual company orchestration | Node.js/PostgreSQL |
| **EasyClaw** | gaoyangz77/easyclaw | Natural-language rules, desktop app | Electron/React/SQLite |
| **ClawKitchen** (JiggAI) | JIGGAI/ClawKitchen | Agile task board + cron admin | OpenClaw plugin |
| **OpenClaw Studio** | grp06/openclaw-studio | One-command NPX install | Next.js/SQLite |
| **TenacitOS** | carlosazaustre/tenacitOS | **3D office view** (React Three Fiber) | Next.js/React 19 |
| **Claworc** | gluk-w/claworc | Multi-instance isolated containers | Go binary (20MB)/Docker |
| **VidClaw** | madrzak/vidclaw | Soul Editor (SOUL.md versioning) | Self-hosted |
| **LobsterBoard** | Curbob/LobsterBoard | 60+ drag-and-drop widgets | Single Node.js file |
| **Claw Empire** | GreenSheep01201/claw-empire (390 ⭐) | Pixel-art office simulator | TypeScript/SQLite |
| **actionagentai dashboard** | actionagentai/openclaw-dashboard | Speech-to-text on every page | Next.js/WebSocket |
| **tugcantopaloglu dashboard** | tugcantopaloglu/openclaw-dashboard | TOTP MFA + security dashboard | Node.js (no build step) |
| **mudrii dashboard** | mudrii/openclaw-dashboard | Zero-dependency, 6 themes | Go + Python |
| **builderz-labs mission-control** | builderz-labs/mission-control | Multi-framework (CrewAI, LangGraph) | Multi-adapter |

All dashboards are **fully self-hostable** with no cloud dependencies, consistent with OpenClaw's local-first philosophy. Most deploy via Docker Compose or a single Node.js command. The ecosystem's main weakness: fragmentation — with 15+ dashboards, no single project has achieved critical mass as the de facto standard.

---

## Observability ranges from built-in OpenTelemetry to kernel-level eBPF

OpenClaw v2026.2+ ships a **built-in OpenTelemetry plugin** (`diagnostics-otel`) that exports traces, metrics (token usage, cost, context size, run duration), and logs over OTLP/HTTP protobuf. Enable it with `openclaw plugins enable diagnostics-otel` and point it at any OTLP-compatible backend — Grafana, SigNoz, Jaeger, or others. Current limitation: log ingestion has known issues (PR#22478), and only `http/protobuf` protocol is supported (gRPC is ignored).

**ClawMetry** (vivekchand/clawmetry) provides the most complete visual observability experience: a real-time dashboard with live animated message flow diagrams, health check heatmaps, token/cost tracking, session inspection, color-coded log streaming, memory browsing (SOUL.md, MEMORY.md), and chat-bubble transcript replay. Install with `pip install clawmetry && clawmetry` — zero config, auto-detects your OpenClaw setup, runs at localhost:8900.

**Knostic's openclaw-telemetry** (knostic/openclaw-telemetry) takes a security-audit angle: it captures every tool call, LLM request, and agent session with built-in redaction, tamper-proof hash chains, and syslog/SIEM forwarding. It also includes network detection scripts (`openclaw-detect`) to find exposed OpenClaw instances. The **openclaw-observability-plugin** (henrikrexed) goes even deeper, combining OpenTelemetry span tracing of the full agent lifecycle with **Tetragon eBPF-based kernel-level monitoring** for tamper-proof security observability on Linux.

**ClawWatch** (0xdsqr/clawwatch) bills itself as "CloudWatch for AI agents" — monitoring performance, tracking costs, setting budgets, and triggering alerts. The **openclaw-logger** project (massaindustries/openclaw-logger) provides a dashboard for exploring OpenClaw logs with AI-powered search. And **openclaw-memvid-logger** captures dual-output conversation logs to both JSONL (backup) and Memvid (semantic search) formats.

---

## Memory management spans seven distinct approaches

OpenClaw's built-in memory system uses **LanceDB** (`memory-lancedb`) with auto-recall and auto-capture, storing embeddings in `~/.openclaw/memory/lancedb/`. It works but offers only basic vector search — no hybrid retrieval, no reranking, no multi-scope isolation.

**memory-lancedb-pro** (win4r/memory-lancedb-pro, **1,100+ stars**) is the most popular upgrade. It adds hybrid retrieval (Vector + BM25 with Reciprocal Rank Fusion), cross-encoder reranking (Jina, SiliconFlow, Pinecone compatible), multi-scope isolation (global, agent, project, user), a management CLI, L0/L1/L2 layered storage, Weibull decay curves for memory aging, two-stage deduplication, and session memory distillation. Multiple community forks exist.

**PostClaw** replaces the entire memory backend with **PostgreSQL + pgvector**. It provides semantic memory (RAG), episodic memory (auto-logged prompts and tool calls), a knowledge graph with typed directed edges, dynamic persona injection, a "sleep cycle" for memory consolidation, and row-level security. Install with `openclaw plugins install @postclaw/postclaw` and a setup command that creates the database, installs extensions, and configures the schema.

**QMD** (Query Markdown Documents) is an official alternative backend combining BM25 + vectors + reranking, running fully locally via Bun + node-llama-cpp. It's ideal for indexing Obsidian vaults, project docs, and session transcripts, though CPU-only VPS deployments suffer from timeouts. The **Cognee** integration adds a **knowledge graph layer** that extracts entities and relationships from Markdown memory and enables graph traversal search beyond pure vector similarity. **openclaw-memory-brain** takes the simplest approach: JSONL-based storage with hash-based embeddings, slash commands for CRUD, and configurable memory caps with oldest-first eviction. And the **Supermemory** plugin syncs OpenClaw workspace memory with Supermemory's cloud service for cross-device access.

Two open feature requests (GitHub issues #22927 and #18595) propose **native PostgreSQL + pgvector support** in OpenClaw core. A proof-of-concept by IrriVision Technologies indexed 1,145 documents with ~200ms search latency using Voyage AI embeddings — but neither proposal has been merged yet.

---

## The skills marketplace and its security crisis

**ClawHub** (github.com/openclaw/clawhub, **5,900 stars**) is the official skill registry with **13,700+ published skills** as of March 2026. Skills cover everything from Gmail automation to Spotify control to smart home management to cryptocurrency trading. The registry supports semver versioning, vector-powered semantic search, and npm-like CLI tooling (`clawhub install`, `clawhub sync`).

However, ClawHub suffered the **ClawHavoc supply-chain attack**: security researchers discovered **1,184 malicious skills** uploaded to the registry, with roughly 20% of all skills being malicious at peak. ClawHub responded by purging 2,419 suspicious skills and partnering with VirusTotal for security scanning. The **awesome-openclaw-skills** curated list (VoltAgent, **3,900 stars**) provides a safer alternative by cataloging 5,400+ pre-vetted skills. Other browsers include AgentSkill.work, OpenClaw Map (openclawmap.com), and doanbactam/agent-skills-directory.

**ClawRecipes** (JIGGAI/ClawRecipes) takes a different approach: pre-built multi-agent team templates. Run `openclaw recipes scaffold development-team` to spin up a Lead, Dev, DevOps, and QA agent with an agile workflow board. Its companion **ClawKitchen** provides the UI dashboard. The project launched on Product Hunt 14 days ago with 94 upvotes.

---

## Database persistence, n8n integrations, and mobile interfaces

For **Supabase/PostgreSQL persistence**, the primary options are PostClaw (PostgreSQL + pgvector plugin), the Supabase Skill by stopmoclay (raw SQL, CRUD, pgvector similarity search via Supabase API), and IronClaw's native pgvector backend (though IronClaw is a separate Rust reimplementation, not a plugin).

**n8n integration** has crystallized into three patterns. The **openclaw-n8n-stack** (caprihan) ships a pre-configured Docker Compose with OpenClaw + n8n on a shared network — OpenClaw calls n8n via webhook URLs. **n8n-claw** (freddy-schuetz) is more ambitious: a full OpenClaw-inspired autonomous agent rebuilt entirely in n8n, with memory via Supabase vector search, MCP client, 15-minute heartbeat, and daily memory consolidation. A **Railway template** bundles OpenClaw + n8n + Langfuse + Tailscale for one-click deployment. Several ClawHub skills (`n8n-api`, `n8n-automation`) also allow OpenClaw to trigger n8n workflows directly.

For **mobile access**, OpenClaw natively supports iOS and Android as "nodes" that pair over the Gateway WebSocket, with Voice Wake and push-to-talk on iOS and voice trigger forwarding on Android. No TestFlight distribution exists — you must build from Xcode source. The official macOS menu bar companion app provides Gateway health monitoring. **ClawFace** is a third-party companion app for monitoring from mobile. In practice, most mobile users simply interact via their messaging apps (WhatsApp, Telegram, etc.), which serve as the primary mobile interface by design.

---

## Multi-agent orchestration and configuration tools

Beyond the dashboards already covered, dedicated multi-agent orchestration tools include **openclaw-agents** (shenhao-stu), which deploys 9 specialized agents (Main, Planner, Ideator, Critic, Scout, Surveyor, Coder, Writer, Reviewer) with one command and group chat routing. **ClawSwarm** (The-Swarm-Corporation) compiles to Rust and runs specialist workers on a unified gRPC gateway. **Antfarm** (snarktank) defines agent teams in YAML with security-reviewed, curated workflows only. An **RFC (#27943)** in the main repo proposes fleet management primitives with declarative YAML manifests and drift detection — not yet merged.

For **agent configuration UIs**, **ClawWizard** (OpenKrab/ClawWizard) provides a full GUI wizard for configuring every aspect of the Gateway — 20+ LLM providers, 20+ messaging channels, persona templates, and one-click deployment — running as a Vite frontend + Node.js bridge. **AlphaClaw** (chrysb/alphaclaw) wraps OpenClaw with a setup wizard, self-healing watchdog, Git-backed rollback, and anti-drift prompt hardening, though it's Docker/Linux only. **VidClaw** uniquely includes a **Soul Editor** for editing SOUL.md, IDENTITY.md, USER.md, and AGENTS.md with version history tracking.

---

## Security tools respond to an ecosystem under active attack

The security ecosystem has grown rapidly in response to real threats. Cisco found prompt injection and data exfiltration in third-party skills. Knownsec's 404 Team documented **63,026 identifiable OpenClaw instances** and 245 recorded vulnerabilities.

**SecureClaw** (adversa-ai/secureclaw) provides 56 audit checks, 15 behavioral rules, 9 scripts, and 4 pattern databases mapping to 7 security frameworks (OWASP, MITRE ATLAS, NIST AI). **ClawSec** (prompt-security/clawsec) offers a complete skill suite including SOUL.md drift detection, live NVD CVE polling, ClawHub reputation checks, and prompt injection defense. **secure-openclaw** (ComposioHQ) provides a Docker-hardened deployment with allowlisted DMs, platform-level access control, and 500+ app integrations via Composio. **ClawdStrike** (backbay-labs/clawdstrike, 54 stars) builds a Swarm Detection & Response (SDR) platform. And **SlowMist's security practice guide** is uniquely designed to be fed directly to the agent — not just read by humans.

---

## Alternative implementations compete on size, speed, and security

The "Claw" family now includes **seven alternative implementations**, each making different trade-offs:

| Implementation | Language | RAM Usage | Startup | Key Differentiator |
|---------------|----------|-----------|---------|-------------------|
| **OpenClaw** (core) | TypeScript/Node.js | 390–700 MB | 4–8 sec | 13,700+ skills, 25+ channels |
| **nanobot** (33.4k ⭐) | Python | ~100 MB | Fast | 99% smaller codebase, 9+ channels |
| **ZeroClaw** | Rust | ~7.8 MB | <10 ms | 400x faster startup, WASM plugins |
| **NanoClaw** | TypeScript | Low | Fast | Container isolation per chat group |
| **PicoClaw** | Go | <1 MB | <1 sec | Runs on $10 embedded devices |
| **IronClaw** | Rust | Low | Fast | WASM sandbox + TEE support, pgvector |
| **NemoClaw** (NVIDIA) | Unknown | Enterprise | TBD | Enterprise RBAC, launching GTC March 16 |

Nanobot from HKU is the most successful alternative with **33,400 stars**, delivering core functionality in ~4,000 lines of Python with MCP support, ClawHub integration, and 9 messaging channels. ZeroClaw targets security-conscious users who want a single static binary with 99% smaller memory footprint. IronClaw from NEAR AI focuses on regulated industries with Trusted Execution Environments. NemoClaw from NVIDIA — not yet released as of March 15, 2026 — targets enterprise with RBAC, audit logging, and partnerships with Salesforce, Cisco, and CrowdStrike.

---

## Commercial products and hosted services

The commercial layer includes **Autonomous Intern** (autonomous.ai) — a dedicated hardware appliance running the full OpenClaw stack with no terminal required — and **ClawBox**, built on NVIDIA Jetson Orin Nano for MCP agents. **OpenClawd AI** (openclawd.ai) offers managed cloud hosting as a separate commercial entity. **Kimi Claw** from Moonshot AI runs OpenClaw natively inside the Kimi web interface with 5,000+ pre-loaded skills and zero setup.

Managed hosting is available from **DigitalOcean** (1-Click Deploy with hardened security image), **xCloud** ($24/month), **Hostinger**, **OVHcloud**, **Contabo**, **LumaDock** ($1.99/month), and **NEAR AI Cloud** (with TEE-backed execution). **Cognio Labs** documents a path to running OpenClaw on Oracle Cloud's Always Free tier (4 ARM CPUs, 24GB RAM).

---

## What's still missing from the ecosystem

Despite the explosion of tools, several gaps remain. There is **no unified standard dashboard** — the fragmentation across 15+ projects means no single solution has become authoritative. **Visual workflow builders** for designing agent pipelines graphically remain primitive. **Dedicated cost management** tools exist but lack the depth of cloud cost optimization platforms. The **mobile experience** remains messaging-app-dependent with no polished native app available via app stores. And critically, **enterprise governance tooling** lags behind adoption — NemoClaw may address this when it launches, but as of today, most governance is ad hoc. The ecosystem's most urgent need is **security standardization**: with 13,700+ skills and documented supply-chain attacks, the gap between what agents can do and what operators can verify remains OpenClaw's greatest vulnerability.

## Conclusion

The OpenClaw enhancement ecosystem has grown from zero to over 100 distinct projects in under 60 days — a pace that mirrors the main project's own unprecedented star growth. The richest categories are **dashboards** (15+ options ranging from single-file scripts to enterprise orchestration platforms), **memory tools** (7 distinct approaches from simple JSONL to full PostgreSQL + pgvector backends), and **security suites** (5+ projects responding to documented attacks). Paperclip at 23,300 stars and nanobot at 33,400 stars have themselves become significant open-source projects in their own right. For operators choosing tools today, the most production-ready combination is the **built-in diagnostics-otel plugin** for monitoring, **memory-lancedb-pro** for enhanced memory, and **either EasyClaw** (for individual users wanting a GUI) **or abhi1693's Mission Control** (for team deployments needing governance). The ecosystem's trajectory suggests consolidation is coming — but for now, the OpenClaw tooling landscape remains as fast-moving and decentralized as the lobster emoji that represents it.