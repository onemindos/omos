# The Complete OpenClaw Ecosystem Inventory

> Original research document from multi-agent ecosystem analysis.
> Preserved as reference — see ECOSYSTEM-42.md for actionable specs.

**Source**: OpenClaw ecosystem inventory research (2026-03)
**Method**: Multi-agent research across GitHub, awesome lists, community platforms, academic papers

---

## 1. Official OpenClaw Repositories (~20 repos)

| Repo | Purpose | Stars | Status |
| --- | --- | --- | --- |
| openclaw/openclaw | Core runtime | 100K+ | Active |
| openclaw/lobster | Workflow engine (.lob files) | — | Active |
| openclaw/openclaw-mcp | MCP integration | — | Active |
| openclaw/acpx | Agent Client Protocol bridge | — | Active |
| openclaw/openclaw-ansible | Ansible deployment | — | Active |
| openclaw/desktop | Desktop client (Electron) | — | Active |
| openclaw/agent-search | Skill search API | — | Active |
| openclaw/ai-toolkit | SDK for skill development | — | Active |
| openclaw/docs | Official documentation | — | Active |
| openclaw/gateway-configs | Example configs | — | Active |
| openclaw/openclaw-plugins | Plugin templates | — | Active |
| openclaw/community-skills | Community skill contributions | — | Active |
| openclaw/registry | Skill registry (ClawHub) | — | Active |
| openclaw/mcp-servers | Official MCP server implementations | — | Active |
| openclaw/openclaw-studio | Web-based admin dashboard | — | New |
| openclaw/infra-examples | Docker/K8s deployment examples | — | Active |
| openclaw/benchmarks | Performance benchmarks | — | Active |
| openclaw/security-advisories | CVE tracking | — | Active |
| openclaw/rfcs | Architecture RFCs | — | Active |
| openclaw/openclaw-voice | Voice integration module | — | Active |

## 2. Awesome Lists (6 curated collections)

| Collection | Focus | Size |
| --- | --- | --- |
| rohitg00/awesome-openclaw | THE definitive collection | 195K★ |
| VoltAgent/awesome-openclaw-skills | 2,868 curated from 5,705+ | Security-rated |
| hesamsheikh/awesome-openclaw-usecases | 34 use cases with code | Already imported |
| awesome-selfhosted/awesome-selfhosted | General self-hosted (includes OC) | — |
| f/awesome-chatgpt-plugins | MCP-adjacent tools | — |
| korchasa/awesome-kubernetes | Deployment patterns | — |

## 3. Voice & Vision Agents (7 projects)

| Project | What | Tech |
| --- | --- | --- |
| VisionClaw | Meta Ray-Ban smart glasses | WebRTC + Gemini Live |
| Jupiter Voice | Fully local STT/TTS | MLX Whisper + Kokoro |
| openclaw-voice | Browser voice chat | faster-whisper + ElevenLabs |
| openclaw-voice-avatar | Lip-synced face | LiveKit + LemonSlice |
| HeyClawy | ESP32 voice companion | WakeNet + speaker |
| MimiClaw | ESP32-S3 pure C agent | Bare metal firmware |
| HomePod/Echo integration | Voice via HA pipelines | Assist Pipeline |

## 4. IoT & Home Automation (7 projects)

| Project | What | Connection |
| --- | --- | --- |
| WireClaw | $5 ESP32 NATS agents | Native NATS pub/sub |
| OpenClawHomeAssistant | HA add-on | Runs OC inside HA |
| openclaw_conversation | HA custom component | HACS install |
| HomeClaw | Apple HomeKit via MCP | Mac Catalyst |
| HA REST skill | Control HA from OC | Existing skill |
| Matter bridge | Matter protocol devices | Via HA |
| Thread sensor mesh | Low-power sensor network | Via WireClaw |

## 5. NATS Messaging & Infrastructure (5 areas)

| Area | What | Status |
| --- | --- | --- |
| JetStream hub config | Hub with leaf node auth | Configured |
| Leaf node configs | Per-location NATS leafs | Template ready |
| MQTT bridge | :1883 for IoT devices | Running |
| WebSocket bridge | :9222 for browser clients | Running |
| Subject hierarchy | digital/physical/nexus/comms | Documented |

## 6. LLM Orchestration & Memory (9 projects)

| Project | What | Differentiator |
| --- | --- | --- |
| ClawRouter | Smart LLM routing | 78% cost savings |
| supermemory | Auto-recall memory | Plugin for OC |
| Mem0 | Persistent memory layer | Dedup + decay |
| openclaw-foundry | Self-upgrading agent | Pattern → workflow auto |
| Composio | 500+ app integrations | OAuth handled |
| Agent Protocol | Pub/sub agent chains | Event-based |
| sessions_send/spawn | Built-in multi-agent | Zero install |
| Langclaw | LangGraph workflows | Python state machines |
| Network-AI | 12-framework orchestrator | Multi-framework bridge |

## 7. OM1 Robotics (6 repos)

| Repo | What | When |
| --- | --- | --- |
| OM1 core | Main robotics runtime | Strategic |
| OM1-video-processor | Real-time video + face rec | Strategic |
| OM1-ros2-sdk | ROS 2 SLAM for Unitree Go2 | Strategic |
| om-yolo-ws | WebSocket YOLO detection | Strategic |
| humanoid_teleoperate | XR robot teleoperation | Strategic |
| FABRIC | Decentralized robot coordination | Strategic |

## 8. Community Tools & Frameworks (12+ projects)

| Project | What | Category |
| --- | --- | --- |
| prompt-security/clawsec | 8 security skills | Security |
| ComposioHQ/secure-openclaw | Hardened Docker shell | Security + integrations |
| Next-Kick/hardened-ansible | Rootless Podman + Squid | Security |
| openclaw-mcp-plugin | Universal MCP client | Integration |
| openclaw-mcp (freema) | Claude → OC bridge | Bridge |
| CrewClaw | Visual multi-agent dashboard | Management |
| ClawWork | Economic capability bench | Benchmarking |
| ExoClaw | WASM-sandboxed runtime | Security |
| BankrBot skills | Crypto/DeFi/Polymarket | Finance |
| Virtual Protocol ACP | Agent commerce | Commerce |
| Clawdstrike | Fleet security monitoring | Security |
| CrewClaw | Multi-agent visualization | Observability |

## 9. Deployment & Infrastructure (8 repos)

| Repo | What | Approach |
| --- | --- | --- |
| openclaw/openclaw-ansible | Official Ansible | IaC |
| Next-Kick/hardened-ansible | Hardened variant | Security + IaC |
| andreesg/openclaw-docker-config | Docker patterns | Containers |
| Terraform modules | Cloud provisioning | IaC |
| workspace backup patterns | Git cron backup | Recovery |
| starter kits (3) | Battle-tested configs | Patterns |
| DigitalOcean templates | DO-specific deployment | Cloud |
| Coolify templates | Self-hosted PaaS | Alternative |

## 10. Research Papers & Studies (8+)

| Paper/Study | Finding | Relevance |
| --- | --- | --- |
| Bitsight scan | 42,900+ exposed instances | Security baseline |
| Kaspersky audit | 512 vulnerabilities (8 critical) | Vulnerability awareness |
| Cisco analysis | ~900 malicious skills | Skill vetting |
| Bitdefender report | Supply chain risks | clawsec motivation |
| ClawWork benchmark | $1,500+/hr economic capability | GE pillar pricing |
| RFC-0001 (Agent Teams) | Native multi-agent spec | sessions_spawn evolution |
| OpenClaw heartbeat whitepaper | Internal agent architecture | System design |
| OM1 research papers | Physical embodiment theory | Long-term vision |

## 11. Community Platforms (9+)

| Platform | Size | URL |
| --- | --- | --- |
| Discord | 115,000 members | discord.gg/clawd |
| Reddit | Active | r/openclaw |
| X Community | 24,400 | twitter communities |
| GitHub Discussions | Thousands | github.com/openclaw |
| Clawstin meetup | Austin-based | clawstin.com |
| Silicon Valley meetup | West coast | meetup.com |
| AI Tinkerers SF | 40 slots | sf.aitinkerers.org |
| Luma Calendar | Worldwide | luma.com/claw |
| Newsletter | Daily | buttondown.com/openclaw-newsletter |

## 12. Media & Learning Resources (12+)

| Resource | Type | Best For |
| --- | --- | --- |
| Official docs | Documentation | Architecture, API reference |
| YouTube channels (10+) | Video tutorials | Visual learners |
| Podcast appearances | Audio | Background listening |
| Blog posts ecosystem | Articles | Deep dives |
| Conference talks | Presentation slides | Architecture patterns |
| Community showcases | Project galleries | Inspiration |
| Skill development guides | Tutorials | Building custom skills |
| Security audit guides | Checklists | Hardening |
| Performance tuning guides | Benchmarks | Optimization |
| Migration guides | Step-by-step | From other platforms |
| Multi-model setup guides | Configuration | Cost optimization |
| Self-hosting cost analysis | Spreadsheets | Budget planning |

---

## Mapping to OneMind 42

All 42 actionable items were extracted from this inventory and documented in ECOSYSTEM-42.md with full specs, install instructions, and OneMind OS integration details.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/platform-catalog (Platform Catalog)/inventory (Platform Inventory)/inventory|inventory]]
