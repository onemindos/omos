---
id: ui-02-20260111053925
title: "Voice integration for Letta: LiveKit vs Vapi technical analysis"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Voice integration for Letta: LiveKit vs Vapi technical analysis"
---

# Voice integration for Letta: LiveKit vs Vapi technical analysis

**Both LiveKit and Vapi can achieve sub-500ms latency with Letta agents, but they represent fundamentally different architectural philosophies.** LiveKit provides an open-source, self-hosted framework with superior Mac/iOS native support and extensive multi-modal capabilities, requiring 3-5 days to production but offering maximum control. Vapi delivers a managed service with faster deployment (1-2 days) and simpler telephony integration, trading customization depth for implementation speed. Letta officially supports both platforms with production-ready integrations.

## Platform architecture and official Letta support

Letta (formerly MemGPT) has evolved into a stateful AI agent platform with dedicated voice capabilities as of late 2024. The platform introduced a specialized `voice_convo_agent` architecture optimized for low-latency interactions, featuring background memory management through "sleep-time" agents that handle context summarization while the primary agent maintains conversational flow.

Both integration paths are officially documented and production-ready. Letta's LiveKit integration uses a native plugin within the LiveKit Agents framework, exposed through `openai.LLM.with_letta()` that connects to either Letta Cloud (`api.letta.com`) or self-hosted instances. The reference implementation at `letta-ai/letta-voice` demonstrates the complete stack using Deepgram for speech-to-text, Cartesia for text-to-speech, and LiveKit's WebRTC infrastructure for transport. Vapi integration operates through Vapi's Custom LLM provider feature, where you configure the endpoint as `api.letta.com/v1/chat/completions` with the model field set to `agent-{AGENT-ID}` format. This approach treats Letta as an OpenAI-compatible endpoint within Vapi's orchestration layer.

The architectural difference is significant. LiveKit positions you as the orchestrator managing STT, LLM, and TTS components explicitly in code, while Vapi abstracts the orchestration into its managed platform where you configure providers through a dashboard. For Letta specifically, this means LiveKit gives you direct access to configure message buffer lengths, enable sleep-time compute, and tune memory management parameters, whereas Vapi treats Letta as a black-box LLM endpoint with less granular control over its stateful features.

## Real-time latency performance and benchmarks

LiveKit achieves sub-500ms end-to-end latency in optimized configurations. The Cerebrium case study demonstrates 492ms total latency using Deepgram STT (110ms), Llama-3-8b LLM (300ms), and Rime TTS (80ms), with only 2ms network overhead when services are co-located in the same VPC. A comprehensive 30-stack benchmark study found 730-800ms first-byte latency with Deepgram, GPT-4.1-nano, and Cartesia Sonic-Turbo, with subsequent turns benefiting from KV-cache optimization. The first turn adds 300-800ms overhead for system prompt processing, while non-English languages consistently add 300-500ms across all configurations.

Vapi achieves 465ms optimized latency in the AssemblyAI study using AssemblyAI Universal-Streaming STT (90ms), Groq Llama 4 Maverick (200ms), and ElevenLabs Flash v2.5 TTS (75ms), plus 100ms network overhead for web clients. However, real-world performance shows significant variance. The typical range is 550-800ms, with Vapi's official documentation advertising 700ms. Critical issues emerge with default turn detection settings that add 1,500ms+ delays if not optimized. European developers report dramatic discrepancies, with dashboard metrics showing 1,100ms while actual calls measure 3,500-4,000ms, attributed to geographic distance from Vapi's infrastructure.

Both platforms demonstrate that **network co-location is more impactful than model selection for latency optimization**. Services in the same VPC communicate with 2ms latency versus 150ms+ for separate cloud services. The component breakdown consistently shows that LLM time-to-first-token plus TTS time-to-first-byte account for over 90% of total latency, making STT effectively negligible when using streaming recognizers.

For Mac OS and iOS specifically, LiveKit's WebRTC-based transport delivers lower network latency than Vapi's WebSocket approach, particularly important on cellular connections. WebRTC's UDP optimization and adaptive bitrate handling provide superior performance under variable network conditions common in mobile deployments.

## Mac OS and iOS deployment capabilities

LiveKit offers mature, production-grade iOS and macOS support through the `client-sdk-swift` package. The SDK is in its 93rd release (v2.7.1+) with full Swift 6 strict concurrency support, covering iOS, iPadOS, macOS, tvOS, and visionOS platforms. It has accumulated 1,463+ commits from 28+ contributors over four years of active development. The architecture uses Sendable protocol compliance for thread safety, weak reference patterns for memory management, and comprehensive delegate systems for real-time events. Notable features include automatic AVAudioSession management, CallKit integration for background calls, simulcast support for adaptive streaming, and native SwiftUI components.

Apple platform developers access LiveKit through Swift Package Manager or CocoaPods, with minimal examples and comprehensive documentation at `docs.livekit.io/home/quickstarts/swift/`. The `agent-starter-swift` template repository provides a complete starting point featuring voice, text, and video input support, virtual avatar rendering, and screen sharing capabilities. This template is actively maintained with 2024-2025 commits and offers configurable features through straightforward code modifications.

Vapi's iOS support exists through multiple pathways but with less maturity. The native iOS SDK at `VapiAI/client-sdk-ios` requires iOS 12.0+ (iOS 14.0+ for screen sharing) and integrates via Swift Package Manager with microphone permissions configuration. The React Native SDK `@vapi-ai/react-native` version 0.2.4 provides cross-platform deployment but adds an additional abstraction layer. For macOS specifically, Vapi primarily relies on its web SDK rather than native macOS support, meaning Mac applications typically embed WebRTC-based voice calls through web views rather than native Swift implementations.

The practical difference manifests in development experience. **Building native iOS apps with LiveKit follows standard iOS development patterns** with SwiftUI views, native audio session handling, and direct WebRTC control. Estimated implementation time is 2-4 hours for setup and 1-2 days for UI customization. Vapi's iOS approach works adequately for voice-only applications but requires navigating web-to-native bridges, with less seamless integration into native iOS UI patterns. For Mac applications specifically, **LiveKit's architecture naturally extends to macOS development while Vapi treats Mac as a web platform**.

## Self-hosting versus managed service trade-offs

LiveKit's open-source Apache 2.0 license grants complete infrastructure sovereignty. Self-hosting requires a Linux VM (Ubuntu or Amazon Linux tested), 16-core compute-optimized instances recommended for production, and SSL certificates for your domain. Port configuration includes TCP 7880 (signaling), TCP 7881 (RTC), and UDP 50000-60000 (media). The simplest deployment uses LiveKit's configuration generator (`docker run livekit/generate`) that produces cloud-init YAML for one-command VM provisioning on AWS, Azure, or Digital Ocean. Development mode runs with `livekit-server --dev` binding to localhost:7880.

Production deployments require Redis for ingress/egress services and multi-replica setups. Kubernetes deployments use official Helm charts with host networking mode (one pod per node) for optimal performance. Benchmark testing on 16-core c2-standard-16 instances shows capacity for 10 publishers with 3,000 subscribers at 80% CPU for large audio rooms, or 150 publishers plus 150 subscribers at 85% CPU for video meetings. Agent workers run efficiently, supporting 30 concurrent voice agents on a 4-core, 8GB machine.

**Deployment complexity is moderate for teams with DevOps experience** but represents a significant commitment for teams without existing infrastructure expertise. The advantage is complete data sovereignty, no per-minute costs beyond infrastructure, and freedom to optimize for specific latency requirements by deploying in regions closest to users. The disadvantage is ongoing maintenance burden including scaling management, monitoring setup, security patching, and handling traffic spikes.

Vapi eliminates infrastructure management entirely. The platform handles 150M+ calls with 400,000+ daily calls currently, automatically scaling to 1M+ concurrent calls. Setup involves creating an account, configuring an assistant through the dashboard, and obtaining API keys. Deployment time is measured in minutes rather than days. The managed model includes built-in call analytics, automatic transcription storage, and simulated test suites within the dashboard.

However, Vapi's managed approach introduces vendor lock-in concerns that are only partially mitigated. The **core orchestration layer—turn-taking logic, interruption handling, VAD coordination, and streaming pipeline management—is proprietary**. Bring-your-own-keys (BYOK) support allows direct billing from STT/LLM/TTS providers, and custom LLM endpoints preserve agent logic portability, but the sophisticated orchestration that makes voice agents feel natural is Vapi's intellectual property. Developers report this orchestration represents 3+ months of development effort to replicate elsewhere.

On-premise deployment exists for Vapi through enterprise agreements, offering a Dockerized deployment on any cloud provider where audio and text data remain entirely within your infrastructure. This addresses data sovereignty concerns but requires contacting Vapi's enterprise sales, suggesting significant cost implications beyond the standard per-minute pricing.

## Multi-modal capabilities analysis

LiveKit provides comprehensive multi-modal support as a core architectural feature. The platform natively handles audio, video, text, and data streams through its WebRTC transport layer with H.264/VP9 codec support for video and Opus for audio. The agent framework supports vision capabilities through multiple patterns: frame sampling on conversation turns where agents capture the latest video frame and add it to LLM context, continuous video processing for realtime multimodal models like Gemini Live, and computer vision pipelines for object detection or environment analysis.

Production examples demonstrate these capabilities extensively. The Gemini Live with Vision example includes a full iOS app processing camera input alongside voice interactions. Avatar integration examples work with Anam, Bey, BitHuman, Hedra, Simli, and Tavus for live avatar rendering synchronized with voice output. Multi-agent workflows enable voice agents that hand off to specialized vision agents for analyzing shared screens or uploaded images. The underlying architecture supports this naturally because **LiveKit was designed from inception as a multi-participant, multi-modal platform** originally built for video conferencing.

Performance considerations show vision adds context richness at the cost of latency. Frame sampling at conversation turns minimizes overhead versus continuous video processing. Modern multimodal LLMs handle this reasonably well, though most models aren't trained for motion understanding. Gemini Live's native video support handles continuous streams better than cascaded vision+language approaches. Recommended practice is strategic frame sampling—capture frames when users ask questions about their environment rather than continuous processing.

Vapi operates primarily as an audio-only platform. Its WebSocket-based transport uses compressed low-bandwidth audio codecs optimized for telephony rather than the multi-stream architecture needed for video. Video capabilities exist exclusively through third-party integration with Tavus Conversational Video Interface (CVI), announced in 2024. This partnership enables digital twin avatars with video presence, but it's an add-on rather than native functionality.

**For voice plus vision use cases, the architectural gap is substantial.** LiveKit enables developers to build agents that simultaneously process audio and video with sub-second latency, while Vapi requires external service integration with additional latency implications. Voice plus text works adequately in Vapi through its web SDK with message sending capabilities, but screen sharing functionality is limited compared to LiveKit's native support. RAG integration exists in both platforms for voice agents querying documents, but Vapi's implementation doesn't extend to visual RAG scenarios.

Letta's multi-modal support with each platform follows the underlying infrastructure capabilities. With LiveKit, Letta agents can maintain conversational state while processing visual inputs through custom tools that capture frames or integrate vision APIs. With Vapi, Letta agents remain audio-only unless you architect external vision processing pipelines that feed results back through custom function calls.

## Implementation speed and starter templates

LiveKit's starter template ecosystem is comprehensive and actively maintained. The `letta-ai/letta-voice` repository serves as the primary integration reference, featuring complete working code for Letta agents with Deepgram STT, Cartesia TTS, and LiveKit infrastructure. Last updated in Q4 2024 through Q1 2025, this repository includes both Letta Cloud and self-hosted server configurations with step-by-step setup instructions requiring Python 3.10+, relevant API keys, and ngrok for local development.

The `agent-starter-python` template is production-grade with Docker support, automated testing via pytest, CI/CD pipelines, and comprehensive evaluation suites. It features the LiveKit Turn Detector for contextually-aware speech detection and background noise cancellation. Setup takes 1-3 hours using uv for dependency management, with full customization requiring 1-2 days. The template includes both console testing mode (`uv run python src/agent.py console`) and web development mode with frontend integration.

Frontend developers access `agent-starter-react` built on Next.js 14+ with TypeScript, offering real-time voice interaction, camera video streaming, screen sharing, audio visualization, virtual avatar integration, and customizable branding through configuration files. One-command setup via `lk app create --template agent-starter-react` followed by `pnpm install` and `pnpm dev` gets developers running in 2-4 hours. For iOS/macOS native applications, `agent-starter-swift` provides SwiftUI implementations supporting iOS, iPadOS, macOS, and visionOS with voice, text, and video input, plus virtual avatar rendering and screen sharing on physical devices.

All LiveKit starter templates are actively maintained with 2024-2025 commits, production-ready with 9-10/10 quality ratings, and accompanied by extensive documentation. The platform offers additional templates for Flutter and React Native covering cross-platform mobile scenarios. The agents repository examples directory contains 70+ working implementations including avatar integrations, Gemini Live with Vision, supervisor escalation patterns, and multi-agent workflows.

Vapi's implementation approach centers on the official Letta integration guide at `docs.letta.com/guides/voice/vapi` rather than separate boilerplate repositories. This comprehensive documentation walks through account setup, Python SDK usage, web and phone testing, cost optimization, and production deployment with troubleshooting sections. Setup time is 1-2 hours for basic configuration, 4-6 hours for production readiness.

Supporting repositories include `vapi-express-starter` for Node.js/Express webhook servers handling function calls and status updates (8/10 quality, 1-2 hours setup), and `client-side-example-javascript-react` demonstrating React integration with message handling, audio visualization, and event management (2-3 hours integration time). The official Web SDK at `VapiAI/client-sdk-web` provides straightforward integration with event listeners and call lifecycle management, achieving basic integration in 1-2 hours.

Vapi's dashboard offers 100+ pre-built templates for customer support, appointment scheduling, sales, and information gathering with 15-30 minute no-code setup times. However, **these templates are dashboard configurations rather than code repositories**, limiting version control integration and collaborative development workflows compared to LiveKit's GitHub-first approach.

**Time-to-production comparison shows Vapi achieving MVP status faster (1-2 days versus 3-5 days for LiveKit) but LiveKit reaching full production readiness with comparable timelines (2-3 weeks for LiveKit versus 1 week for Vapi) when accounting for the additional customization and optimization typically required.** For multi-modal applications, LiveKit offers 2-3 week timelines while Vapi lacks native support.

## Pricing implications and long-term costs

LiveKit's self-hosted cost model centers on infrastructure rather than per-minute fees. The Cerebrium deployment case study breaks down costs as $1.44/hour for STT (Deepgram on A10 GPU supporting 180 concurrent connections), $4.32/hour for LLM inference (Llama-3-8b on H100 supporting 9 concurrent requests), and $1.44/hour for TTS (A10 GPU for 180 concurrent connections), plus $0.03/hour per conversation for media transport. This yields approximately $0.03 per minute per call for infrastructure, excluding licensing fees for proprietary models.

At scale, these costs become highly competitive. A deployment handling 1,000 concurrent calls with appropriate GPU provisioning might cost $500-1,000 monthly in infrastructure versus $13,000-31,000 for the same volume on a per-minute managed platform charging $0.13-0.31 per minute. The break-even point typically occurs around 5,000-10,000 minutes monthly, after which self-hosting becomes significantly cheaper. However, this calculation excludes DevOps labor—a dedicated engineer managing the infrastructure represents $100,000+ annually in labor costs, amortized across all projects.

LiveKit Cloud offers a managed alternative with 1,000 free minutes for testing and paid plans scaling with usage. Pricing is higher than self-hosted but competitive with other managed platforms. The service eliminates infrastructure management while maintaining LiveKit's flexibility and multi-modal capabilities. This option suits teams wanting LiveKit's features without deployment complexity.

Vapi's pricing structure appears simple at $0.05/minute base platform fee but accumulates costs across the stack. Complete per-minute breakdown includes: $0.05 Vapi platform fee, $0.01-0.0126 for speech-to-text (Deepgram, Gladia, Talkscriber), $0.02-0.20 for language models (GPT-4o-mini at ~$0.02, GPT-4 at ~$0.20, Claude 3 Opus at ~$0.09), $0.005-0.07 for text-to-speech (Vapi voices $0.0216, ElevenLabs $0.036, PlayHT $0.0648), and $0.0055-$0.014 for telephony through Twilio/Vonage/Telnyx.

**Realistic total cost per minute ranges from $0.13 to $0.31 depending on provider selections**, meaning 1,000 minutes costs $130-310 monthly, not the advertised $50. Enterprise features add overhead: HIPAA/SOC2 compliance costs $1,000 monthly, additional concurrent lines are $10/line/month beyond the 10 included, phone numbers cost ~$2/month each, and custom enterprise agreements typically range $40,000-70,000 annually.

Hidden complexity emerges from managing multiple vendor relationships. You receive separate invoices from Vapi, OpenAI (or your chosen LLM provider), Deepgram (or chosen STT provider), ElevenLabs (or chosen TTS provider), and Twilio (if using telephony). Budget forecasting becomes challenging when LLM costs scale unpredictably with conversation length and complexity, and premium voice options dramatically increase TCO.

Bring-your-own-keys (BYOK) mitigates some lock-in by establishing direct billing relationships with STT/LLM/TTS providers rather than paying Vapi's markup. You still pay Vapi's $0.05/minute orchestration fee but control provider relationships and can negotiate enterprise pricing directly. This reduces Vapi's portion to about 20-40% of total costs rather than 100%.

Long-term cost trajectory differs fundamentally. **LiveKit's costs scale with infrastructure (mostly linear) while Vapi's costs scale with usage (purely variable)**. For proof-of-concept and early-stage projects with low volume, Vapi's variable costs are lower than maintaining LiveKit infrastructure. For established products with predictable high volume, LiveKit's infrastructure investment amortizes favorably. The crossover point varies but generally favors LiveKit beyond 10,000 minutes monthly.

## Decision framework and recommendations

**Choose LiveKit for Letta integration when:**

You require infrastructure sovereignty for compliance reasons like HIPAA with full data control or government/enterprise deployments with strict security requirements. Your application roadmap includes multi-modal capabilities such as voice plus vision for visual assistance scenarios, screen sharing for collaborative applications, or video avatars for enhanced presence. The iOS and macOS native experience matters significantly—building App Store-ready applications with native UI patterns or requiring deep integration with iOS CallKit and system audio management.

Engineering capacity exists to manage infrastructure—you have DevOps expertise for self-hosting and monitoring or development teams comfortable with infrastructure-as-code. Call volume justifies infrastructure investment, typically exceeding 10,000 minutes monthly where per-minute costs become prohibitive. You need maximum customization of STT/LLM/TTS provider combinations, specialized latency optimization for specific use cases, or unique conversation flow patterns that managed platforms don't support adequately.

The open-source requirement matters for your organization, whether for code auditability and security review, modification rights for custom features, or philosophical preference for open platforms. Project timeline accommodates 3-5 days to MVP and 2-3 weeks to production, and you're building a voice AI platform rather than just using one.

**Implementation path:** Start with `letta-ai/letta-voice` repository, deploy Letta server via Docker (self-hosted) or use Letta Cloud, configure LiveKit Cloud for initial testing with 1,000 free minutes, use `agent-starter-swift` for iOS/macOS prototyping, implement custom tools in Letta for business logic, transition to self-hosted LiveKit for production scale, and deploy agent workers with 4-core, 8GB instances supporting 30 concurrent agents.

**Choose Vapi for Letta integration when:**

Speed to market is critical—you need proof-of-concept in 1-2 days or MVP in less than one week. Team composition lacks dedicated DevOps engineers or infrastructure expertise. Telephony is primary delivery channel for phone-first applications like customer support or appointment scheduling where phone number provisioning must be turnkey.

Use cases fit standard patterns—appointment booking, customer support, lead qualification, information gathering—scenarios covered by Vapi's pre-built templates. Call volume remains moderate, generally under 10,000 minutes monthly where per-minute pricing is economical. Development resources are constrained and you need dashboard-driven configuration rather than code-first development.

Geographic distribution matches Vapi's infrastructure (primarily US-based), though you must verify latency in your target regions before committing. The application is audio-only with no vision requirements or multi-modal features needed within the next 12-18 months. Budget can accommodate variable costs scaling with usage rather than fixed infrastructure investment.

**Implementation path:** Create Letta agent via SDK or dashboard with `voice_convo_agent` type, configure Vapi assistant using Custom LLM provider, set endpoint to `https://api.letta.com/v1/chat/completions`, set model field to `agent-{YOUR_AGENT_ID}` format, optimize turn detection settings (critical for latency), implement webhooks for custom business logic using `vapi-express-starter`, test thoroughly in production-like conditions (Vapi's dashboard metrics may not reflect actual performance), and consider BYOK for all providers to reduce vendor lock-in.

## Critical technical considerations

Latency optimization requires specific attention regardless of platform choice. **Default configurations rarely achieve sub-500ms performance**—Vapi's default turn detection settings add 1,500ms+ delays that require explicit optimization. Co-location of services in the same VPC saves 150ms+ versus separate cloud providers. Model selection drives 90%+ of latency with LLM time-to-first-token plus TTS time-to-first-byte as the critical path. First conversation turns always incur 300-800ms overhead for prompt processing. Non-English languages add 300-500ms consistently across all configurations.

For Mac and iOS deployment, **LiveKit's native Swift SDK provides significantly better developer experience and performance than Vapi's web-based approach**. Native CallKit integration, background audio handling, and SwiftUI components work seamlessly with LiveKit while requiring workarounds with Vapi. Screen sharing, video capture, and multi-modal features integrate naturally in LiveKit but require external services with Vapi.

Geographic testing is mandatory before production commitment, especially for Vapi where European developers report 3-4x higher latency than US deployments. Always test in your target deployment regions with production-like network conditions. Dashboard metrics may not reflect actual user experience, particularly for Vapi where significant discrepancies have been reported.

Infrastructure sovereignty considerations extend beyond data privacy to long-term strategic flexibility. LiveKit's open-source license means you can modify, extend, or fork the codebase if needed. Vapi's proprietary orchestration layer represents vendor lock-in that persists even with BYOK and custom LLM endpoints. The orchestration logic that handles turn-taking, interruption management, and streaming coordination is Vapi's core intellectual property, not easily replicated.

Cost modeling must account for all expenses. For LiveKit, include infrastructure costs, DevOps labor, monitoring tools, and scaling management time. For Vapi, include the platform fee, all provider costs (STT/LLM/TTS/telephony), potential HIPAA compliance fees if needed, and enterprise agreement minimums for custom features. Calculate total cost of ownership over 12-24 months at expected volume rather than comparing per-minute rates.

## Performance metrics summary

**LiveKit proven performance:**
- 492ms minimum latency (Cerebrium case study with optimal stack)
- 730-800ms typical with GPT-4.1-nano and Cartesia Sonic
- Sub-500ms achievable with Llama-3-8b and Rime TTS
- Consistent performance across geographic regions with proper deployment
- 30 concurrent agents per 4-core, 8GB worker
- Powers ChatGPT Advanced Voice Mode at massive scale

**Vapi proven performance:**
- 465ms minimum latency (AssemblyAI study with Groq Llama 4 and ElevenLabs Flash)
- 550-800ms typical with standard configurations
- 700ms advertised in official documentation
- Geographic variance: US deployments typically meet targets; European deployments report 3,500-4,000ms in some cases
- Turn detection adds 1,500ms+ if not optimized
- Scales to 1M+ concurrent calls on managed infrastructure

Both platforms demonstrate that sub-500ms latency is achievable, meeting the natural conversation threshold. **The difference lies not in peak performance but in consistency, ease of optimization, and architectural flexibility.** LiveKit requires more initial effort but delivers predictable performance with full control. Vapi offers faster deployment but requires careful configuration and performs variably across different deployment scenarios.

For Letta specifically, both integrations are production-ready and officially supported. The `voice_convo_agent` architecture with sleep-time compute works effectively on both platforms. LiveKit's integration provides more direct access to Letta's stateful features and memory management parameters. Vapi's integration is simpler to configure but treats Letta more abstractly as an LLM endpoint. Neither integration is inherently superior—the choice depends on your broader architectural requirements, team capabilities, and strategic priorities around infrastructure sovereignty versus managed services.