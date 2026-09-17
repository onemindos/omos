---
id: ui-02-20260111053925
title: "Building Legacy: Production-Ready AI Agent Architecture"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Building Legacy: Production-Ready AI Agent Architecture"
---

# Building Legacy: Production-Ready AI Agent Architecture

**The landscape for multi-modal AI agent systems reached production maturity in 2025, with clear winners emerging across memory, orchestration, and real-time communication layers.** For your Legacy system requiring voice commands, persistent personality, video awareness, and offline capability, the optimal architecture combines LangGraph for orchestration, Letta or Mem0 for memory, LiveKit for real-time communication, and a carefully selected offline stack—all while avoiding vendor lock-in through open standards and self-hosting options.

This report evaluates 8 major frameworks, 5 memory systems, 4 offline AI platforms, and provides a complete integration blueprint based on October 2025 production deployments. Key finding: no single platform dominates, but combining best-in-class components creates enterprise-grade systems already powering applications at massive scale.

## Memory systems have evolved beyond simple context windows

The memory layer represents the most critical architectural decision for Legacy's personality persistence and proactive behavior. As of October 2025, four production-ready frameworks dominate, each excelling in different dimensions.

**Letta (formerly MemGPT) remains the most advanced option for complex multi-agent systems** with true self-editing memory and unique cross-agent memory sharing. Version 0.12.1 released October 2025 introduces hierarchical memory with in-context core blocks and out-of-context archival storage, allowing agents to actively modify their own memory through tool calls. Founded by UC Berkeley's Sky Computing Lab (the same team behind Databricks and Anyscale), Letta raised $10M seed funding from Felicis Ventures in September 2024 with elite angel investors including Google's Jeff Dean and Hugging Face's CEO. With 18,600 GitHub stars and Apache 2.0 licensing, **Letta provides maximum transparency and zero vendor lock-in**. The framework supports personality state through editable memory blocks, procedural memory through tool definitions, and cross-agent sharing where a single memory block feeds multiple agents—critical for Legacy's distributed architecture. However, the v1 architecture introduces complexity with a steeper learning curve than alternatives.

**Mem0 emerged as the production performance leader** in 2025 benchmarks, achieving 26% higher accuracy than OpenAI's memory system with 91% faster responses and 90% lower token usage according to LOCOMO evaluations. YC-backed with $500K seed funding and 22,000 GitHub stars, Mem0 implements a hybrid architecture combining graph stores for relationships, vector stores for semantic search, and key-value stores for metadata. The platform offers the simplest API (just 4 core functions) while supporting multi-level memory across users, sessions, and agents. **For production deployments prioritizing reliability and performance, Mem0 represents the safest choice** with clear business backing and proven benchmarks. The managed Mem0 Cloud service handles infrastructure complexity, though this creates moderate vendor lock-in compared to Letta's pure open-source approach.

**Zep delivers state-of-the-art temporal reasoning** through its Graphiti temporal knowledge graph that tracks facts with valid_at and invalid_at timestamps, enabling Legacy to understand how information changes over time. YC W2024-backed, Zep dominates the LongMemEval benchmark with graph RAG capabilities that excel at dynamic data tracking. The framework automatically extracts facts, classifies dialog intent and emotion, and assembles context for LLMs with real-time incremental updates. However, the community edition was deprecated in 2025, pushing users toward Zep Cloud—introducing vendor lock-in concerns despite technical excellence.

**LangGraph with persistence provides the strongest ecosystem integration** for teams already invested in LangChain. The framework introduced cross-thread memory stores in 2025, complementing its checkpointer system for conversation threads. While memory features are newer and less specialized than dedicated frameworks, LangGraph's graph-based workflow design excels at complex multi-agent orchestration. With 11,700 GitHub stars and 4.2M monthly downloads, LangGraph benefits from the massive LangChain ecosystem (100,000+ stars, 2,000+ contributors). Production deployments include Klarna serving 85M users with 80% faster resolution and AppFolio achieving 2x response accuracy.

**For Legacy's requirements of personality persistence, cross-agent memory sharing, and proactive behavior, the recommended stack is Letta for memory management** paired with Qdrant for vector storage. This combination provides maximum control, complete transparency, and proven multi-agent capabilities while maintaining zero vendor lock-in. Teams prioritizing simplicity over control should choose Mem0, while those needing temporal reasoning should evaluate Zep despite cloud dependency.

## Vector databases split between managed convenience and self-hosted performance

The vector database layer underlies all memory operations, with clear segmentation between enterprise managed services and high-performance open-source alternatives.

**Pinecone dominates the managed vector database market** with serverless auto-scaling, sub-50ms latency handling billions of vectors, and enterprise compliance including SOC 2 Type II, ISO 27001, GDPR, and HIPAA certifications. Zero operational overhead and excellent documentation make Pinecone the default choice for teams without dedicated DevOps resources. However, costs scale rapidly from ~$70/month to $500+ at production volumes, and proprietary architecture creates complete vendor lock-in with no self-hosting option.

**Qdrant represents the open-source performance champion** with Rust-based architecture delivering sub-100ms latency while consuming less memory than alternatives. Advanced pre-filtering, quantization support, and ACID compliance rival Pinecone's capabilities while supporting both self-hosting and managed Qdrant Cloud (~$102/month for 1M vectors). SOC 2 certification provides enterprise credibility. The single-binary deployment simplifies self-hosting dramatically compared to complex distributed systems. For Legacy's requirement to avoid vendor lock-in, **Qdrant offers the optimal balance of performance, cost-efficiency, and operational simplicity**.

**Weaviate excels at hybrid search** combining BM25 keyword search with vector similarity, making it ideal for knowledge graph applications with its GraphQL API and multiple vector space support. SOC 2 Type II and HIPAA compliance (AWS, 2025) provide enterprise validation. However, pricing complexity and steeper learning curves make it better suited for advanced use cases requiring multi-modal or complex search patterns rather than standard agent memory.

**ChromaDB provides the fastest path to prototyping** with Python-native embedding, simple APIs, and no server requirements for development. However, performance degrades significantly beyond 100K vectors, making it unsuitable for production at Legacy's anticipated scale. Use ChromaDB for initial development, then migrate to Qdrant for production deployment.

**pgvector extends PostgreSQL with vector capabilities**, enabling Legacy to consolidate vector and relational data in a single database. This approach leverages existing PostgreSQL expertise, provides ACID compliance, and integrates seamlessly with Supabase. While slower than specialized vector databases at massive scale, pgvector's SQL interface and mature ecosystem make it compelling for moderate-scale deployments under 10M vectors.

**Recommended vector database strategy: Deploy Qdrant self-hosted** on the same infrastructure as other Legacy components, using Docker containers managed by Portainer. Configure Q4_K_M quantization for 4x compression with minimal accuracy loss. Implement pgvector in Supabase as a secondary store for structured data requiring vector search. Budget 4-6GB VRAM for effective 7-8B model embeddings with Qdrant.

## LiveKit Agents framework powers production voice AI at massive scale

Real-time communication forms Legacy's nervous system, with LiveKit Agents v1.2.15 (October 2025) emerging as the production-proven framework for stateful, multimodal AI agents. **LiveKit powers ChatGPT's Advanced Voice Mode** and handles 3 billion calls annually for 100,000+ developers, providing unmatched validation for production reliability.

The v1.0 release in April 2025 brought fundamental improvements including enhanced workflows, native Model Context Protocol support, semantic turn detection across 13 languages under 25ms CPU, and AI avatar integration. Series B funding of $45M (April 2025) brought total funding to $83M with backing from Altimeter Capital, validating market leadership in real-time AI infrastructure.

**Stateful agents represent LiveKit's core strength** through the UserData pattern enabling session state persistence across conversation turns. The typed dataclass approach provides type safety while maintaining flexibility for complex state management including user preferences, conversation history, and personality configurations. Context preservation across agent handoffs enables sophisticated multi-agent workflows where specialist agents (billing, technical support, sales) maintain conversation continuity as users transfer between capabilities.

**Voice integration supports two primary patterns**: traditional STT-LLM-TTS pipelines with providers like Deepgram, OpenAI, and ElevenLabs, or OpenAI's Realtime API for native speech-to-speech processing. Semantic turn detection with configurable interruption handling provides natural conversation flow. The VoiceActivityVideoSampler enables vision capabilities by sampling video frames at different rates during speech (1 FPS) versus silence (0.3 FPS), reducing processing costs while maintaining awareness.

**Multi-agent coordination via tool-based handoffs** enables complex workflows. A routing agent uses function tools to transfer conversations to specialist agents while preserving the ChatContext for continuity. The experimental Task-based coordination (AgentTask API) provides structured phase-based workflows like consent collection, information gathering, and task completion. Pre-built tasks including GetEmailTask demonstrate production-ready patterns for common operations.

**Model Context Protocol integration reached production maturity** in v1.2+ with native support for HTTP/SSE transports. MCP servers enable agents to access external systems through standardized interfaces—Home Assistant for smart home control, Supabase for database operations, Kubernetes for SRE tasks (VoiceOps project), and Zapier for workflow automation. Authentication supports secret keys, OAuth tokens, and A2A protocol for agent-to-agent communication. This positions LiveKit as the orchestration layer for Legacy's external integrations.

**Wake word detection lacks native support**, requiring integration with external libraries. The recommended approach uses openWakeWord, an open-source library providing high-accuracy detection under 50ms with Apache 2.0 licensing. Implementation involves audio frame processing before the LiveKit pipeline activates, with detection triggering the active listening state. Alternative commercial options like Picovoice Porcupine offer higher accuracy under 30ms but introduce licensing costs and vendor dependencies. **For Legacy's "hey legacy" wake word, implement openWakeWord** integrated before the LiveKit STT pipeline with a threshold above 0.5 confidence.

**Production deployment with LiveKit Cloud** provides stateful load balancing, elastic scaling, graceful draining during updates, and zero-downtime deployments. Self-hosting remains viable for Kubernetes deployments with custom worker pools, though this sacrifices the noise-cancellation plugin that requires Cloud infrastructure. For Legacy, **deploy LiveKit self-hosted** on infrastructure managed by Portainer to maintain zero vendor lock-in while accepting manual scaling responsibilities.

Documentation and examples demonstrate production-ready patterns: drive-thru ordering systems, medical triage agents, front desk assistants, and MCP integration projects. The agents-playground repository provides interactive testing environments. Enterprise deployments at OpenAI, Hello Patient, Salient, and Podium validate production readiness at scale.

**LiveKit limitation: Video processing for awareness** requires integration with external systems like Frigate for camera feeds and face recognition. LiveKit handles video streams well but lacks built-in computer vision capabilities beyond basic frame sampling. This necessitates external video AI integration covered in the video processing section below.

## Offline AI reached production readiness with comprehensive tooling

Local AI operation without internet dependency became fully viable in 2025 with mature platforms, high-quality models, and proven hardware configurations spanning edge devices to enterprise servers.

**Ollama emerged as the production API leader** with OpenAI-compatible endpoints, built-in model management, and 34% faster inference than LM Studio on the same hardware (141.59 tokens/second on Qwen 1.5B). The Go-based server backed by llama.cpp provides request batching for concurrent users—essential for multi-agent systems. Simple CLI operations (ollama run llama3) enable rapid deployment, while the Docker container integrates seamlessly with Portainer orchestration. GitHub repository at github.com/ollama/ollama provides full transparency. **For Legacy's offline LLM requirements, Ollama represents the optimal choice** with its combination of performance, simplicity, and API compatibility.

**LM Studio provides the best GUI experience** for non-technical users with an Electron-based interface, integrated Hugging Face model browser, and superior performance on Apple Silicon through MLX engine optimization (26-30% faster than alternatives). While less suitable for production server deployment, LM Studio excels for development, testing, and demonstration scenarios. Download from lmstudio.ai.

**llama.cpp offers maximum control** for embedded systems and resource-constrained environments with its pure C/C++ implementation consuming just 90MB. Support for CUDA, Metal, Vulkan, OpenCL, SYCL, and ROCm enables optimization across diverse hardware. Advanced users seeking absolute performance control or deploying to specialized hardware should use llama.cpp directly, though this requires deeper technical expertise. Repository: github.com/ggml-org/llama.cpp.

**LocalAI provides a unified interface** for multiple AI tasks beyond text generation, including image generation and audio processing (Whisper, TTS) with a drop-in OpenAI API replacement. Support for multiple backends (llama.cpp, vLLM, Transformers) enables flexibility across deployment scenarios. For systems requiring diverse AI capabilities through a single API, LocalAI simplifies integration. GitHub: github.com/mudler/LocalAI, documentation: localai.io.

**Voice model selection prioritizes faster-whisper for STT** with 4x speed improvement over OpenAI Whisper while maintaining identical accuracy. CTranslate2-based implementation with 8-bit quantization processes 16-minute audio in 14 seconds versus 46 seconds for whisper.cpp and 60+ seconds for baseline Whisper. Repository: github.com/SYSTRAN/faster-whisper. Alternative whisper.cpp optimizes for CPU and edge devices with Apple Silicon ANE support (3x speedup via Core ML), making it ideal for Raspberry Pi deployments. Repository: github.com/ggml-org/whisper.cpp.

**Piper TTS dominates edge deployment** with fast neural text-to-speech using ONNX and VITS models that achieve real-time synthesis on Raspberry Pi 4. Multi-language support with natural quality at under 200ms latency makes Piper the obvious choice for voice assistant applications. Installation via pip install piper-tts, models at github.com/rhasspy/piper/blob/master/VOICES.md. **For Legacy's offline TTS, deploy Piper** across all edge devices while reserving higher-quality options for cloud deployments.

**Coqui TTS/XTTS-v2 enables voice cloning** with just 6-30 seconds of audio across 17 languages at realtime factor ~1/3 on GPU. Despite the company shutdown, the open-source codebase remains production-ready. Repository: github.com/coqui-ai/TTS, PyPI: pypi.org/project/coqui-tts/. Alternative StyleTTS 2 provides better licensing terms than XTTS while maintaining high quality. For Legacy's personality features, **voice cloning enables consistent voice identity** across agents and sessions.

**Hardware requirements scale systematically** from edge to enterprise. Entry-level 1-3B models with Q4 quantization require just 2-4GB RAM, running on Raspberry Pi 4 and laptops. Consumer hardware with 8-16GB RAM and RTX 3060 12GB GPU handles 7-8B models at Q4_K_M quantization for 20-50 tokens/second. Enthusiast configurations with RTX 4090 24GB support 30-34B models at Q4 or 7-13B models at Q5-Q6 for 40-100 tokens/second. Professional deployments with A6000 48GB+ can run 70B models at Q4-Q5 or 120B models at Q4 for 30-140 tokens/second. The recently released 671B models require 160GB (2x H100 80GB) with 1.58-bit quantization.

**Quantization format selection determines portability and performance**. GGUF provides the best cross-platform compatibility with CPU efficiency and layer offloading, making Q4_K_M the recommended default (balanced quality/size) and Q5_K_M for quality-prioritized deployments. AWQ delivers the highest quality for 4-bit quantization with 20-30% faster inference than GPTQ while avoiding calibration overfitting. GPTQ maintains mature ecosystem support but underperforms AWQ in recent benchmarks. **For Legacy, standardize on GGUF Q4_K_M** for portability across deployment targets, with Q5_K_M for primary inference nodes.

**Offline sync with Supabase leverages PowerSync** as the recommended production solution. PowerSync provides a Postgres-SQLite sync layer reading from PostgreSQL WAL and syncing to local SQLite with declarative sync rules supporting RLS. Conflict resolution defaults to Last Write Wins with support for custom strategies. Platforms span Flutter, React Native, and JavaScript. Documentation: docs.powersync.com/integration-guides/supabase-+-powersync, website: powersync.com. Alternative WatermelonDB requires manual sync implementation but provides full control, while ElectricSQL offers active-active replication with CRDT conflict resolution for collaborative scenarios. **For Legacy's offline-first architecture, implement PowerSync** for reliable bidirectional sync when connectivity permits.

**Production recommendations for consumer hardware** (16GB RAM, 8-12GB GPU) combine Ollama for LLM serving, faster-whisper for speech recognition (small or medium model), Piper or XTTS-v2 for synthesis, and PowerSync for data synchronization. Deploy Llama 3.1 8B or Mistral 7B at Q4_K_M to Q5_K_M quantization. This configuration provides production-grade performance while maintaining reasonable resource consumption. For edge devices like Raspberry Pi, downgrade to Phi-3 Mini 3.8B (Q4_K_M) via llama.cpp, whisper.cpp (tiny or small), and Piper TTS with aggressive data scoping in PowerSync.

## Multi-agent frameworks matured with distinct architectural philosophies

The multi-agent orchestration landscape consolidated around four primary approaches in 2025, with the global AI agents market projected to grow from $5.29B (2023) to $216.8B by 2035 at 40.15% CAGR.

**LangGraph dominates enterprise deployments** with graph-based architecture providing precise control over agent workflows. The directed acyclic graph approach enables cyclical workflows, stateful checkpointing, time-travel debugging, and parallel execution. LangGraph Studio IDE provides visual debugging with state editing and replay functionality. Version 1.0 maturity brought cross-thread memory stores, Redis checkpoint support, and enhanced persistence APIs. **Unicorn valuation of $1.1B** (July 2025) and $135-160M total funding from IVP validate market leadership. With 11,700 GitHub stars, 4.2M monthly downloads, and proven deployments at Klarna (85M users), AppFolio (2x accuracy), and Elastic, LangGraph represents the production-proven choice for complex multi-agent systems.

The framework excels at precise orchestration control, superior state management through checkpointer systems, and excellent debugging capabilities via LangGraph Studio. Integration with the broader LangChain ecosystem (100,000+ GitHub stars, 130M+ monthly downloads) provides access to 70+ LLM providers and thousands of pre-built components. However, steep learning curve, complex setup requirements, and occasional documentation gaps create barriers for teams without dedicated AI engineering resources. **For Legacy's complex multi-agent coordination requirements across voice, video, IoT, and automation, LangGraph provides the architectural sophistication needed** for production reliability.

**CrewAI emerged as the rapid development champion** with role-based agent teams using specific goals and backstories. Independent architecture (not LangChain-based) with 39,400 GitHub stars and ~1M monthly downloads demonstrates rapid adoption since January 2024 launch. Insight Partners and Boldstart led the $18M Series A (October 2024) for a ~$100M valuation. CrewAI Studio visual interface and YAML/Python configuration enable fast prototyping. The framework reports 150 enterprise customers, 10M+ agents executed monthly, usage by nearly half of Fortune 500, and 100K+ developers certified through community courses.

Built-in memory types (short/long-term, entity) with SQLite3 backend provide adequate capabilities for most use cases without external memory systems. The crew metaphor (agents with roles like captain, engineer, designer) creates intuitive mental models for team-based workflows. 1,200+ tool integrations cover diverse external systems. However, logging issues, less flexibility than LangGraph for complex orchestration, and lack of streaming function calls limit advanced use cases. **Teams prioritizing development velocity over architectural sophistication should choose CrewAI**, accepting the trade-off of less control for faster implementation.

**AutoGen (Microsoft) transitioned to maintenance mode** with the framework being merged into the new "Agent Framework" with only bug fixes for existing AutoGen codebases. Despite this sunset status, AutoGen pioneered conversational multi-agent patterns with 48,077 GitHub stars and enterprise deployments at Novo Nordisk and major banking institutions. Event-driven asynchronous architecture with AutoGen Studio no-code tool and support for both Python and .NET (C#) demonstrated Microsoft's enterprise infrastructure commitment. Conversation-based coordination, multiple topology patterns, and robust human-in-the-loop capabilities established patterns now being integrated into Microsoft's unified agent framework. **Teams should not start new projects with AutoGen**, but existing deployments remain supported with migration paths to the successor framework.

**OpenAI Agents SDK launched March 2025** as a lightweight Python framework built on three primitives: Agents, Handoffs, and Guardrails. Provider-agnostic support for 100+ LLMs via LiteLLM removes OpenAI lock-in concerns. Excellent tracing and visualization combined with minimal abstractions enable 10-minute setup for basic agents. 9,000+ GitHub stars in just 7 months demonstrate rapid adoption. However, the new framework status means smaller community, limited case studies, and fewer production deployments compared to mature alternatives. Documentation: platform.openai.com/docs/agents, GitHub: github.com/openai/agents-sdk. **For new projects prioritizing simplicity and modern APIs**, OpenAI SDK provides compelling developer experience despite youth.

**Semantic Kernel (Microsoft) introduces Azure vendor lock-in risk** despite technical strengths. Enterprise SDK for C#, Python, and Java with plugin/skills architecture powers Microsoft 365 Copilot. Model-agnostic design claims broad compatibility, but Azure optimization creates soft lock-in through performance incentives and feature availability. 21,000+ GitHub stars and Microsoft backing ensure long-term viability, but teams avoiding cloud vendor dependencies should look elsewhere. The plugin architecture and multi-language support excel for organizations already committed to Azure and Microsoft ecosystems.

**Haystack (deepset) focuses on RAG and semantic search** rather than multi-agent orchestration, making it better suited as a complementary component than primary framework. Component-based pipeline architecture with Haystack 2.0 flexibility and Deepset Studio drag-and-drop interface excel at document processing and semantic search. However, limited multi-agent capabilities and pipeline-based rather than graph-based design restrict use for complex orchestration scenarios.

**AutoGPT pioneered autonomous agents** in March 2023 with viral success achieving 177,350 GitHub stars—one of the most-starred repositories on GitHub. However, unclear business model, uncertain direction, known issues with agent looping, and just $12M funding (October 2023) with no follow-on rounds create viability concerns. While historically significant and maintaining brand recognition, AutoGPT lacks the production readiness and business backing of alternatives. **Avoid AutoGPT for production deployments** until business model clarity and technical reliability improve.

**AgentOps provides essential observability** rather than orchestration, occupying a complementary niche in the ecosystem. The platform integrates with CrewAI, AutoGen, LangGraph, OpenAI SDK, LlamaIndex, and CAMEL through 2-line code additions. Session replays, cost tracking, behavior monitoring, and real-time alerts address critical production needs. Pre-seed funding of $2.6M (August 2024) from 645 Ventures and Afore Capital provides early validation, but very early stage creates risk. Competitive alternatives include LangSmith (LangChain), Helicone, and Langfuse. **For production Legacy deployments, implement LangSmith** (mature, proven) rather than AgentOps (early stage) for observability needs.

**Integration with n8n, Node-RED, Home Assistant, EMQX, and Portainer** varies by framework. LangGraph provides deep integration with n8n through LangChain nodes and orchestration layers, enabling visual workflows to trigger complex agent operations. CrewAI integrates via HTTP nodes and n8n workflows as tools. All frameworks support HTTP Request node integration for basic connectivity. Node-RED excels at IoT/hardware integration with MQTT support ideal for EMQX broker communication. Home Assistant achieved leading-edge AI integration in 2025 with native OpenAI, Google Gemini, Anthropic Claude, and Ollama support, plus full Model Context Protocol integration. Custom integrations like home-generative-agent demonstrate LangGraph integration patterns. Portainer manages Docker deployments for all frameworks with stack templates, GPU support for local LLMs, and microservices architecture recommendations.

**For Legacy's comprehensive requirements, adopt a hybrid architecture**: **LangGraph for core multi-agent orchestration**, n8n for cloud/AI workflow automation, Node-RED for IoT/edge device coordination, AgentOps or LangSmith for observability, and integrate via EMQX MQTT broker for event-driven communication. This combination provides production-proven components with minimal vendor lock-in while supporting the full feature spectrum from voice interaction to IoT device control.

## Video awareness requires external computer vision integration

Real-time video processing for agent awareness necessitates integration between camera systems and AI models, with Frigate NVR emerging as the leading open-source solution.

**Frigate provides production-ready video management** with real-time object detection, person tracking, and event-based recording. The platform integrates natively with Home Assistant for automation triggers and notifications. Face recognition requires integration with external services like CompreFace, DeepStack, or Facebox since Frigate focuses on object detection rather than facial identification.

**CompreFace offers open-source facial recognition** with REST API for adding faces, recognizing faces in images, and managing face collections. Self-hosted deployment maintains privacy while avoiding cloud dependencies. Integration involves Frigate detecting person events, sending frame snapshots to CompreFace for identification, and returning identity information to agents via MQTT or webhooks.

**Stateful video processing** tracks people across camera views, maintains appearance history, and correlates events with other agent context. Implement through Supabase storing person detections, CompreFace identities, and location tracking, with EMQX publishing real-time events to agents. Legacy agents subscribe to video events via MQTT topics like `video/{location}/{event_type}`, enabling contextual responses like "Welcome home [name]" when recognized faces enter specific zones.

**Hardware requirements for real-time video AI** depend on camera count and inference complexity. Frigate recommends Google Coral TPU for 5-10 cameras at 10fps or NVIDIA GPU for 10+ cameras. Model selection balances accuracy and speed—YOLO variants provide real-time performance while larger models improve accuracy at latency cost.

**Privacy considerations mandate local processing** for facial recognition data, encrypted storage of face embeddings, clear consent mechanisms, configurable retention policies, and zone-based processing restrictions. Frigate's local-first architecture aligns with these requirements when paired with self-hosted CompreFace.

The video research subagent experienced technical issues preventing comprehensive coverage of this topic. Key resources include Frigate documentation at docs.frigate.video, CompreFace at github.com/exadel-inc/CompreFace, and Home Assistant integration at home-assistant.io/integrations/frigate. For Legacy, **deploy Frigate with Coral TPU acceleration** for object detection, integrate CompreFace for facial recognition, publish events to EMQX, and subscribe agents to video events for contextual awareness.

## Platform stability analysis reveals clear winners and risks

Financial viability and vendor lock-in risk vary dramatically across the AI agent ecosystem, with funding, community metrics, and open-source commitment determining long-term sustainability.

**Tier 1: Excellent viability (85-100 points)** includes LangChain (95/100), Supabase (92/100), and Semantic Kernel (90/100). LangChain achieved unicorn status with $1.1B valuation, $135-160M funding, and 1M+ monthly active developers building 40-60% of LLM applications. LangSmith generates estimated $12-16M ARR with enterprise customers including Elastic, Robinhood, Rakuten, Moody's, JPMorgan, and BlackRock. Model-agnostic support for 70+ providers and MIT open-source licensing create very low vendor lock-in (2/10 score).

Supabase reached $5.12B valuation with $196-520M funding through Series E (October 2025), serving 900,000+ signups with 40% of YC startups using the platform. Enterprise traction includes Meta, Netflix, Microsoft, and GitHub. PostgreSQL foundation with pgvector support for AI use cases provides standards-based portability. Open-core model enables self-hosting while managed Supabase Cloud offers convenience. Some proprietary hosted services create moderate lock-in (5/10 score) but PostgreSQL compatibility enables migration paths.

Semantic Kernel benefits from Microsoft backing providing infinite runway and guaranteed support, with 21,000+ GitHub stars and MIT open-source license ensuring transparency. Cross-platform support (C#, Python, Java) and model-agnostic design provide flexibility, though Azure integration creates soft lock-in (3/10 score) through performance optimization and feature availability.

**Tier 2: Strong viability (70-84 points)** includes CrewAI (80/100), LiveKit (78/100), and Letta (75/100). CrewAI raised $18M Series A (October 2024) from Insight Partners and Boldstart for ~$100M valuation, with 30,000 GitHub stars and 10M+ agents executed monthly demonstrating rapid adoption. 150 enterprise customers and 100K+ certified developers validate product-market fit. Fully open-source framework creates low lock-in (4/10 score) despite proprietary enterprise features.

LiveKit secured $83M total funding with $45M Series B (April 2025) from Altimeter Capital, reaching ~$400M+ valuation. Powering OpenAI's Voice Mode and Character AI at 3 billion calls annually provides unmatched production validation. Open-core model with WebRTC standards enables self-hosting, though cloud features create moderate lock-in (5/10 score).

Letta raised $10M seed (September 2024) from Felicis at $70M valuation despite very early stage and small 9-person team. UC Berkeley pedigree and elite angel investors (Google's Jeff Dean, Hugging Face CEO) provide credibility. 18,600 GitHub stars and MIT open-source licensing ensure complete transparency (2/10 lock-in score). Early stage creates moderate viability risk, but academic foundation and quality investors suggest strong trajectory.

**Tier 3: Moderate viability (50-69 points)** includes AgentOps (62/100) and AutoGPT (58/100). AgentOps raised just $2.6M pre-seed (August 2024) with small team and competitive observability space creating higher risk. Clear product-market fit addressing critical needs provides hope, but early stage and potential acquisition risk by larger platforms warrant caution.

AutoGPT achieved viral success with 169,000 GitHub stars but suffers from unclear business model, single $12M funding round (October 2023) with no follow-on, and uncertain direction. High GitHub stars demonstrate brand strength, but lack of clear product-market fit and revenue strategy create significant viability concerns (58/100 score).

**Vendor lock-in risk scoring** (lower is better) places Letta and LangChain at 2/10 (very low) due to full open source and model-agnostic design. Semantic Kernel scores 3/10 (low) despite Azure soft lock-in. CrewAI at 4/10 (moderate-low) balances open framework with proprietary enterprise features. Supabase and LiveKit both score 5/10 (moderate) with open cores but proprietary cloud features. AgentOps reaches 7/10 (higher risk) as SaaS-focused with uncertain future.

**Development activity metrics** show LangChain and Supabase with very active daily commits and multiple weekly releases. CrewAI demonstrates rapid iteration since January 2024 launch. Letta maintains active development with regular releases. AutoGPT shows moderate activity focused on platform rebuild. Community metrics place LangChain with 10,000+ Discord members and 30,000+ Twitter followers leading engagement.

**Funding history reveals clear trends**: Supabase ($520M, $5.12B valuation) and LangChain ($135M, $1.1B valuation) as established unicorns provide maximum stability. LiveKit ($83M) and CrewAI ($18M) show strong growth trajectory. Letta ($10M seed) and AgentOps ($2.6M pre-seed) represent higher-risk early-stage bets. AutoGPT ($12M, single round) creates viability concerns.

**Platform-specific risk assessment**: LangChain, Semantic Kernel, and Supabase present LOW overall risk with very low financial risk, low technical risk, and low lock-in risk. Letta, CrewAI, and LiveKit present MEDIUM risk with moderate financial risk but low technical and lock-in risk. AgentOps and AutoGPT present HIGH risk across multiple dimensions.

**For Legacy production deployment, build on LangGraph (LangChain ecosystem)** for orchestration with very low vendor lock-in, proven enterprise adoption, and excellent financial backing. **Use Letta for memory management** accepting early-stage risk for maximum control and transparency. **Deploy on self-hosted infrastructure** managed by Portainer using Qdrant, Ollama, Supabase (self-hosted), and LiveKit (self-hosted) to minimize vendor dependencies. This architecture provides production-grade capabilities while maintaining complete control and migration optionality.

## Integration architecture requires event-driven, distributed design

Connecting LiveKit + Memory Layer + Supabase + n8n/Node-RED + Home Assistant + EMQX into a cohesive system demands careful architectural planning with proven patterns from production deployments.

**Four-layer architecture** separates concerns while enabling coordination. Layer 1 (Communication & Real-Time) includes LiveKit for WebRTC audio/video/data and EMQX for MQTT 5.0 IoT messaging, forming the nervous system. Layer 2 (Orchestration & Workflow) combines n8n for cloud/AI workflows and Node-RED for IoT/edge automation as the coordination brain. Layer 3 (Data & State) integrates Supabase for PostgreSQL + real-time + edge functions + auth with dedicated memory layer (Letta/Mem0) for agent context, providing persistent and real-time memory. Layer 4 (Integration & Edge) leverages Home Assistant for 2,000+ device integrations and Portainer for container orchestration, enabling physical world interaction.

**Orchestrator-worker pattern eliminates point-to-point complexity** by designating n8n/Node-RED as orchestrators with agents and devices as workers. EMQX distributes work via topic subscriptions. Communication flows from voice input through LiveKit to STT/LLM via Supabase edge function, webhook to n8n for orchestration querying Supabase context and publishing EMQX commands triggering Home Assistant device actions, with device responses via MQTT updating state in Supabase and UI via LiveKit real-time channels.

**Hybrid n8n/Node-RED strategy** optimizes for capabilities. Use n8n for cloud services, SaaS integrations, AI/LLM workflows (OpenAI, LangGraph), complex business logic, and external API orchestration. Use Node-RED for IoT device integration, MQTT message routing, edge deployment (lower footprint), Home Assistant automation, and local/offline operations. Integration via webhooks or MQTT bridge enables best-of-both-worlds functionality.

**EMQX message queue implementation** provides masterless clustering for no single point of failure, supporting 100M+ concurrent connections with sub-millisecond latency. MQTT 5.0, MQTT-SN, and MQTT over QUIC protocols enable diverse connectivity. Topic structure organizes messages: `home/{location}/{device_type}/{device_id}/{command|status|telemetry}` for devices, `agents/{agent_id}/{input|output|state}` for agents, `livekit/{room_id}/{participant_id}/{event}` for real-time, and `n8n/{workflow_id}/{trigger|result}` for workflows. QoS strategy uses QoS 0 (fire-and-forget) for telemetry, QoS 1 (at-least-once) for commands, and QoS 2 (exactly-once) for critical state.

**Supabase real-time synchronization** leverages PostgreSQL CDC (Change Data Capture) with WebSocket subscriptions and row-level security for secure updates. Event-driven architecture publishes state changes through EMQX retained messages for immutable event logs providing replay capability. CRDT-based sync enables conflict-free offline operation with automatic conflict resolution. State management layers separate ephemeral (Redis, in-memory), session (time-limited Supabase sessions), and persistent (PostgreSQL, backed up) data with last-write-wins conflict resolution for non-critical data and business logic for domain-specific cases.

**Portainer edge deployment** supports standard Edge Agent (periodic 5-60s polling) and Edge Agent Async (on-demand tunnel) without requiring inbound ports through TLS-secured tunnels managing 25,000+ devices. GitOps integration enables infrastructure-as-code deployments. Edge stack pattern deploys central Portainer (cloud/datacenter) managing remote locations running Edge Agent with Docker Swarm/Standalone hosting n8n/Node-RED, EMQX (local), Home Assistant, AI inference, and optional Supabase. Deployment strategy implements staged rollouts, pre-pulls images, includes retry logic (hourly × 7 days), and uses environment-specific configurations.

**Scalability patterns** differ by component. LiveKit scales through distributed mesh architecture supporting millions of users with edge-to-edge routing and auto-scaling. EMQX achieves linear scaling to 100M connections by adding nodes dynamically with DNS/HAProxy load balancing and geo-distributed deployment. Supabase implements connection pooling (pgBouncer), read replicas for queries, auto-scaling edge functions, and CDN for static assets. n8n scales via queue mode for background jobs, multiple worker nodes, database-backed queues, and webhook async triggers.

**Security and authentication** implements OAuth 2.0/OIDC framework with centralized identity provider (Supabase Auth, Auth0, Okta), JWT tokens for stateless authentication, short-lived access tokens (15 minutes) with refresh tokens, and scope-based authorization. Multi-factor authentication uses FIDO2/WebAuthn (passwordless, strongest), TOTP authenticator apps, SMS/email OTP (fallback), and biometric (mobile) with adaptive MFA based on risk scoring. Row-level security in Supabase enforces user data isolation through PostgreSQL policies. MQTT security requires TLS encryption (port 8883), client certificates (mTLS), topic-level ACL, and payload encryption. API security implements rate limiting (100 req/min), input validation, HTTPS/TLS minimum 1.2, CORS configuration, and API versioning.

**Robotics and AR integration future-proofing** follows Internet of Robotic Things (IoRT) architecture with sensors/actuators, connectivity (MQTT, 5G, WiFi 6), edge computing (local AI inference), cloud intelligence (training, updates), and HMI (AR glasses, voice, gesture). Real-time data overlay through LiveKit streams AR video to AI analysis (edge/cloud) with overlay info via WebSocket requiring under 100ms latency. Digital twin patterns synchronize physical robots to Supabase via MQTT with VR/AR twin (Unity/Unreal) enabling remote operation. Edge AI deployment uses ONNX runtime (cross-platform), TensorFlow Lite (mobile/edge), or PyTorch (research) with versioning and A/B testing. ROS2, WebRTC, MQTT, and gRPC provide standard protocols for robotics communication.

**Recommended technology stack** combines LiveKit (open-source, self-hostable SFU, $0 self-hosted), EMQX (100M+ connections, free OSS or enterprise), n8n (400+ integrations, free self-hosted or $20/user/month), Node-RED (IoT/edge, free), Supabase (PostgreSQL + real-time, free tier or $25/month pro), Home Assistant (2000+ integrations, free), Portainer (user-friendly edge, free CE or $10/node/month), LangGraph/LangChain (orchestration, MIT open-source), Letta (memory, MIT open-source), Qdrant (vector DB, open-source), and Ollama (local LLM, open-source). Total self-hosted cost ranges $0-50/month for small deployments to $100-300/month for hybrid recommended stack balancing cost, maintenance, and reliability.

**Implementation roadmap** follows five phases over 20 weeks. Phase 1 (weeks 1-4) establishes foundation with Portainer, EMQX, Supabase, OAuth 2.0, and monitoring. Phase 2 (weeks 5-8) implements core integration deploying n8n, Node-RED, Home Assistant, basic workflows, API gateway patterns, and edge deployment. Phase 3 (weeks 9-12) adds real-time layer with LiveKit, voice AI agents, state sync, and multi-agent workflows. Phase 4 (weeks 13-16) focuses on scaling through EMQX clustering, Supabase replicas, load balancing, caching, and performance testing. Phase 5 (weeks 17-20) delivers advanced features including AR/VR integration, robotics integration, advanced AI workflows, multi-region deployment, and disaster recovery testing.

## Production recommendations for Legacy architecture

Based on comprehensive research across memory, orchestration, real-time, offline, video, stability, and integration dimensions, the following architecture delivers production-grade capabilities while avoiding vendor lock-in.

**Core technology stack:**

- **Orchestration Layer**: LangGraph (v1.0+) for complex multi-agent workflows with stateful checkpointing, time-travel debugging, and graph-based coordination. Self-hosted or LangChain Cloud for managed deployment.

- **Memory Management**: Letta (v0.12.1+) for self-editing memory with cross-agent sharing, episodic/semantic/procedural memory types, and personality state persistence. Vector storage via Qdrant (self-hosted) with Q4_K_M quantization.

- **Real-Time Communication**: LiveKit Agents (v1.2.15+) for voice/video/data streams with stateful session management, native MCP support, and tool-based agent handoffs. Self-hosted deployment via Docker.

- **Workflow Automation**: n8n for cloud/AI/SaaS integration and Node-RED for IoT/edge/local automation with MQTT coordination.

- **Message Broker**: EMQX for MQTT 5.0 broker with clustering, data bridges, and sub-millisecond latency. Self-hosted with masterless architecture.

- **Data Layer**: Supabase (self-hosted PostgreSQL + Supabase stack) for real-time database, authentication, edge functions, and storage with pgvector for supplementary vector operations.

- **Smart Home Integration**: Home Assistant for 2,000+ device integrations with native AI support including MCP protocol, Assist with custom LLMs, and automation creation.

- **Container Orchestration**: Portainer for Docker/Kubernetes management with edge deployment support, stack templates, and GPU passthrough for local AI.

- **Offline AI**: Ollama for local LLM serving (Llama 3.1 8B Q4_K_M), faster-whisper for STT (medium model), Piper TTS for voice synthesis, and PowerSync for Supabase offline sync.

- **Video Processing**: Frigate NVR for object detection with Coral TPU acceleration, CompreFace for facial recognition (self-hosted), integration via MQTT events to agents.

- **Observability**: LangSmith for LangGraph monitoring, session replay, cost tracking, and debugging. Prometheus + Grafana for infrastructure metrics.

**Architecture diagram:**

```
┌─────────────────────────────────────────────────────────────┐
│                        LEGACY SYSTEM                         │
│                     Multi-Modal AI Agents                    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                   USER INTERACTION LAYER                     │
├─────────────────────────────────────────────────────────────┤
│  Voice: "Hey Legacy" (openWakeWord) → LiveKit Agents        │
│  Video: Frigate NVR + CompreFace → MQTT Events              │
│  Smart Home: Home Assistant → Automation Triggers            │
│  Manual: Web/Mobile UI → WebSocket/REST APIs                │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                 COMMUNICATION & REAL-TIME LAYER              │
├─────────────────────────────────────────────────────────────┤
│  LiveKit: WebRTC for voice/video/data (self-hosted)         │
│  EMQX: MQTT 5.0 broker for IoT messaging (clustered)        │
│  Supabase Realtime: PostgreSQL CDC via WebSocket            │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│            ORCHESTRATION & WORKFLOW LAYER                    │
├─────────────────────────────────────────────────────────────┤
│  LangGraph: Multi-agent coordination (self-hosted)          │
│  n8n: Cloud/AI workflow automation (self-hosted)            │
│  Node-RED: IoT/edge automation (edge deployed)              │
│  MCP: Model Context Protocol for tool integration           │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                AGENT MEMORY & INTELLIGENCE LAYER             │
├─────────────────────────────────────────────────────────────┤
│  Letta: Self-editing memory (episodic/semantic/procedural)  │
│  Qdrant: Vector database (self-hosted, Q4_K_M quant)        │
│  Ollama: Local LLM inference (Llama 3.1 8B)                 │
│  faster-whisper: Speech-to-Text (medium model)              │
│  Piper TTS: Text-to-Speech (multi-language)                 │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    DATA & STATE LAYER                        │
├─────────────────────────────────────────────────────────────┤
│  Supabase PostgreSQL: Persistent storage + auth + RLS       │
│  PowerSync: Offline sync (Postgres ↔ SQLite)               │
│  Redis: In-memory cache + pub/sub                           │
│  pgvector: Supplementary vector operations                  │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│              INTEGRATION & EDGE LAYER                        │
├─────────────────────────────────────────────────────────────┤
│  Home Assistant: Device control (2000+ integrations)        │
│  Frigate + CompreFace: Video AI + face recognition          │
│  Portainer: Container orchestration (edge + cloud)          │
│  Edge Nodes: Local processing, offline capability           │
└─────────────────────────────────────────────────────────────┘
```

**Communication flow example (voice command):**

1. User says "Hey Legacy" → openWakeWord detects wake phrase
2. LiveKit activates → faster-whisper transcribes speech
3. Transcript → Ollama LLM processes with Letta memory context
4. LLM response → n8n workflow orchestrates multi-step actions:
   - Query Supabase for user preferences and state
   - Publish MQTT command to EMQX topic
   - Trigger Home Assistant automation
   - Update Letta memory with interaction
5. Devices respond via MQTT → state updates in Supabase
6. Response generated by Ollama → Piper TTS synthesizes voice
7. Audio streamed via LiveKit → user hears response
8. All events logged to EMQX retained messages for replay

**Key features enabled:**

- **Voice Commands**: "Hey Legacy" wake word via openWakeWord integration before LiveKit pipeline, with full conversation context from Letta memory

- **Standby Mode**: LiveKit user_away_timeout with state monitoring, proactive notifications via EMQX scheduled tasks

- **Personality Persistence**: Letta editable memory blocks storing personality traits, preferences, and communication style across all agent instances

- **Schedule Management**: LangGraph workflows with time-based triggers stored in Supabase, proactive reminders published via EMQX to multiple interfaces

- **Face Recognition**: Frigate detects person → CompreFace identifies → MQTT event to agents → contextual greeting with user history from Letta

- **Cross-Agent Memory**: Letta shared memory blocks enable multiple specialized agents (home, calendar, entertainment) to access unified user context

- **Offline Operation**: Ollama + faster-whisper + Piper enable full voice AI without internet, PowerSync syncs data when connectivity returns

- **Proactive Behavior**: Scheduled tasks in Supabase trigger n8n workflows analyzing context via Letta, publishing notifications based on user patterns

- **Robotics Ready**: MQTT protocol with ROS2 bridge enables robot control, digital twin in Supabase, telemetry via EMQX

- **AR Glasses Ready**: LiveKit streams video from glasses, AI analysis at edge via Ollama, overlay commands via WebSocket, under 100ms latency

**Hardware requirements:**

- **Central Server** (if deploying centrally): 32GB RAM, RTX 4090 24GB or equivalent for Ollama, 1TB NVMe SSD, 10Gbps network
- **Edge Nodes** (per location): 16GB RAM, 8-12GB GPU (RTX 3060), 512GB SSD, sufficient for Ollama 8B + voice models
- **IoT Gateway** (per site): Raspberry Pi 4 8GB or equivalent, Coral TPU for Frigate, 128GB SD card
- **Network**: Gigabit LAN, WiFi 6 for devices, 5G/fiber WAN for cloud connectivity

**Cost estimate (self-hosted):**

- Hardware: $2,000-5,000 (central) + $800-1,500 per edge node + $150 per IoT gateway
- Software: $0 (all open-source)
- Cloud services (optional): $25/month Supabase Pro (if not self-hosting), $20/month n8n (if not self-hosting)
- Electricity: ~$50-100/month for always-on servers
- **Total first year**: $3,000-7,000 + $600-1,200 ongoing annually

**Deployment timeline:**

- **Weeks 1-4**: Infrastructure setup (Portainer, EMQX, Supabase, Qdrant)
- **Weeks 5-8**: Basic integration (n8n, Node-RED, Home Assistant, first workflows)
- **Weeks 9-12**: AI layer (Ollama, Letta, LiveKit, voice agents)
- **Weeks 13-16**: Optimization (clustering, caching, performance tuning)
- **Weeks 17-20**: Advanced features (video AI, face recognition, robotics prep)
- **Production launch**: Week 20+

**Risk mitigation:**

- **Vendor Lock-in**: Zero proprietary dependencies, all components self-hostable with standard protocols (MQTT, WebRTC, PostgreSQL, REST)
- **Single Points of Failure**: EMQX clustering, Supabase read replicas, multiple edge nodes, health checks and auto-restart
- **Performance Bottlenecks**: Redis caching, connection pooling, load balancing, horizontal scaling for stateless services
- **Data Loss**: Automated PostgreSQL backups (daily + WAL archiving), EMQX retained messages, edge node local persistence
- **Security Breaches**: OAuth 2.0 with JWT, mTLS for service-to-service, network segmentation, regular security audits
- **Component Failure**: Graceful degradation (offline mode continues working), circuit breakers, retry logic, monitoring alerts

**Compliance and privacy:**

- **Data Residency**: Full control through self-hosting, no data leaves infrastructure
- **GDPR Compliance**: Right to erasure implemented via Supabase policies, audit logs, consent management
- **Face Recognition**: Local processing only, encrypted storage, configurable retention, zone-based restrictions
- **Voice Data**: Processed locally via faster-whisper, not sent to external services, ephemeral by default
- **Audit Trail**: All agent actions logged to Supabase with timestamps, user IDs, and action details

This architecture provides production-grade reliability, complete vendor independence, comprehensive feature coverage for Legacy's requirements, and proven scalability from small home deployments to enterprise installations. All components have active development, strong communities, and long-term viability based on funding analysis.

Total research compiled from official documentation, GitHub repositories, company blogs, funding databases, production case studies, and developer community discussions as of October 2025.

---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
