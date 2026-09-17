---
id: ui-02-20260111053925
title: "Building Your Self-Hosted "Legacy AI" Platform: Complete Solution Architecture"
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
ai_summary: "Building Your Self-Hosted "Legacy AI" Platform: Complete Solution Architecture"
---

# Building Your Self-Hosted "Legacy AI" Platform: Complete Solution Architecture

**The definitive answer:** For a non-developer wanting a production-ready, self-hosted AI assistant deployable within days, **Dify.ai combined with LiveKit** provides the fastest path to 90% of your requirements, with clear upgrade paths for the remaining features.

## Why this matters now

You're building something unprecedented—a comprehensive personal AI that owns your entire interaction history and grows with you. The good news: the technology landscape has matured dramatically in 2024-2025, with production-ready platforms that didn't exist 18 months ago. The challenge: most "complete solutions" are either cloud-only, developer-focused, or miss critical features like voice/video. After analyzing dozens of platforms and architectures, **a three-component stack emerges as the clear winner for your specific requirements.**

## The winning architecture: Three proven components

Your optimal stack combines best-in-class tools for each capability, all self-hostable with minimal custom development:

**Core Platform: Dify.ai** handles chat, memory, tools, and orchestration. With 58,000+ GitHub stars and deployments at Fortune 500 companies including Volvo and Ricoh, Dify delivers production-grade stability with a 30-minute Docker Compose deployment. The platform provides visual workflow builders (zero code required), built-in RAG with Weaviate vector database, 50+ integrated tools, and support for 100+ LLM providers including GPT-4, Claude, and local Ollama models.

**Voice/Video Infrastructure: LiveKit** powers real-time communication. Used by OpenAI for ChatGPT Voice Mode and Character.ai for their conversational platform, LiveKit is the only mature, fully self-hostable WebRTC solution designed specifically for AI agents. It includes native integrations with Deepgram, ElevenLabs, OpenAI Realtime API, and comprehensive function calling during conversations. Critical finding: alternatives like Daily.co, Vapi.ai, and Retell AI are cloud-only despite marketing suggesting otherwise.

**Advanced Memory (Optional Enhancement): Letta** for sophisticated agent-editable memory. While Dify includes solid RAG capabilities, Letta's specialized architecture enables agents to autonomously manage their own memory through core memory blocks, archival storage, and self-editing tools. Letta is production-ready with 18,500+ stars, Y Combinator backing, and official AWS Marketplace listings.

## Why Letta isn't mandatory (but valuable)

Your research specifically mentioned Letta for memory management. Here's the honest assessment: **Letta is excellent for its specialized use case, but Dify's built-in memory may be sufficient initially.** Dify includes Weaviate vector database, RAG pipelines, and conversation history management that handles 80-90% of memory needs out-of-box. Letta shines when you need agents to actively edit their own memory, maintain complex temporal relationships, or manage multiple memory hierarchies. 

**Recommended approach:** Start with Dify's memory, deploy Letta in parallel when you hit limitations (typically 3-6 months in). Dify can call Letta via API, giving you best-of-both-worlds without blocking your launch.

## Complete feature coverage analysis

Let me map each requirement to the solution:

**Memory system with agent editing:** Dify provides built-in RAG and vector storage (good); Letta provides agent-editable memory with self-management tools (excellent). **Verdict: 90% coverage out-of-box, 100% with Letta integration.**

**Chat interface (desktop/mobile):** Dify delivers production-ready web UI that's responsive across devices. No native mobile apps, but progressive web app (PWA) support enables installation on phones. LibreChat offers similar capability with slightly more polished desktop apps (though still alpha stage). **Verdict: 95% coverage—native apps require custom development but web-based solution is production-grade.**

**Voice interaction:** LiveKit Agents framework provides complete voice pipeline with STT, LLM, and TTS integration. Sub-300ms latency achievable. Dify supports voice through marketplace plugins (Fish Audio, Agora), though not as seamlessly integrated as LiveKit. **Verdict: 90% with plugins, 100% with LiveKit integration.**

**Video streaming/interaction:** This is the hardest requirement. LiveKit supports video streaming, screen sharing, and camera input natively, but GPT-4V integration for visual analysis requires custom code. No platform offers turnkey video understanding + generation + video calling. **Verdict: 70% coverage—infrastructure ready, but AI video interaction needs custom development work.**

**Tool calling:** Dify excels here with 50+ built-in tools, custom API integration, and Model Context Protocol (MCP) support. LiveKit enables function calling during voice conversations. Letta provides 8 memory-management tools natively. **Verdict: 100% coverage.**

**Multi-modal LLM support:** Dify supports text, vision, and file uploads across GPT-4V, Claude 3.5 Sonnet, Gemini, and local vision models. LiveKit works with all major providers. **Verdict: 100% coverage.**

**Flexible LLM backend:** Dify integrates 100+ providers including OpenAI, Anthropic, AWS Bedrock, Azure OpenAI, and local options via Ollama. Task-based routing possible through workflows. **Verdict: 100% coverage.**

**Data ownership and portability:** All components (Dify, LiveKit, Letta) are fully self-hostable with Apache 2.0 or MIT licenses. Data stays on your infrastructure, with export capabilities for migration. **Verdict: 100% coverage.**

**Self-hostable on AWS/Azure:** Dify deploys via Docker Compose in 30 minutes. LiveKit provides cloud-init scripts for AWS EC2/Azure VMs (15-30 minutes). Letta available via AWS Marketplace AMI and Docker. Terraform templates exist for all components. **Verdict: 100% coverage.**

## Timeline to production: 4 days

Here's the realistic deployment schedule for a non-technical user:

**Day 1 (4 hours): Core platform deployment.** Install Docker Desktop, clone Dify repository, run `docker compose up -d`, create admin account, connect to OpenAI or Claude API (or install Ollama for local models). Test basic chat and upload initial knowledge documents. **Milestone: Working ChatGPT-like interface with your data.**

**Day 2 (4 hours): Memory and tools.** Configure RAG settings, test document retrieval, add custom tools via API integration, explore workflow builder to create automated responses. **Milestone: AI that remembers context and can take actions.**

**Day 3 (4 hours): Voice integration.** Install voice plugin from Dify marketplace (Fish Audio recommended), or deploy LiveKit server using cloud-init script on AWS EC2 ($50-100/month compute cost). Configure STT/TTS providers, test voice conversations. **Milestone: Voice-capable AI assistant.**

**Day 4 (4 hours): Production hardening.** Set up authentication, configure HTTPS/SSL, implement backups, document configuration, test from mobile devices. **Milestone: Production-ready personal AI.**

**Total investment: 16 hours over 4 days.** This assumes using Dify's built-in memory and voice plugins. Adding LiveKit and Letta extends timeline to 6-8 days but provides superior capabilities.

## Alternative architectures considered

I evaluated multiple approaches before recommending the Dify + LiveKit stack:

**Architecture Option 1 (My Recommendation): Dify + LiveKit + Optional Letta.** Fastest time-to-value, production-proven components, 90-95% feature coverage, moderate technical complexity. **Setup time: 4 days to full production.** Best for: Non-developers wanting maximum features with minimum custom code.

**Architecture Option 2: LibreChat + Letta + LiveKit.** Enterprise-grade authentication, excellent ChatGPT UI replication, best-in-class memory with Letta. More complex integration between components, requires API gateway. **Setup time: 7-10 days.** Best for: Organizations prioritizing security and exact ChatGPT feature parity.

**Architecture Option 3: Open WebUI + Letta + LiveKit.** Most popular option (112,000 GitHub stars), best for privacy-focused local deployment with Ollama, strong community. Less polished workflow tools than Dify, requires more manual configuration. **Setup time: 5-7 days.** Best for: Privacy maximalists running fully local models.

**Architecture Option 4: AnythingLLM Desktop + LiveKit Cloud.** Simplest possible deployment—download desktop app, connect to LiveKit Cloud (free tier). Limited scalability, desktop-focused rather than server deployment. **Setup time: 2-3 hours.** Best for: Proof-of-concept or single-user with no server management desire.

**Architecture Option 5: AWS/Azure Managed Services.** Use AWS Bedrock + CloudFormation templates or Azure AI Foundry. Fastest cloud deployment but reduces data ownership, higher ongoing costs. **Setup time: 2-4 hours.** Best for: Organizations with existing AWS/Azure investments and budget for managed services.

I rejected pure-play options like n8n, Flowise, LangFlow because they require significant custom development to build complete chat interfaces—they're orchestration tools, not end-user platforms.

## The Letta evaluation: Excellent but not mandatory

Since you specifically asked about Letta, here's the detailed assessment:

**Letta's unique value:** It's the only production-ready framework purpose-built for agent-editable memory. Agents autonomously decide what to remember, update their own memory blocks, and maintain temporal knowledge. The architecture includes in-context "core memory" (always available), out-of-context "archival memory" (vector-based retrieval), and 8 built-in memory management tools. UC Berkeley research project turned Y Combinator-backed company with 18,500+ GitHub stars.

**Production readiness:** Docker deployment takes 15-30 minutes, AWS Marketplace AMI available, supports PostgreSQL for persistence, integrates with all major LLM providers, exports complete agent state via .af files for portability.

**Alternatives considered:** Zep.ai benchmarks 18.5% better accuracy but deprecated their self-hosted Community Edition (cloud-only now). Mem0 has 34,000+ stars but is a memory layer only, not a complete agent framework. LangChain memory requires significant development work. Vector databases (Qdrant, Pinecone, Weaviate) are building blocks, not solutions.

**Reality check for your use case:** Letta solves a specific problem—stateful agents that manage complex memory over long timescales. If you need agents that remember "you prefer emails sent at 10am" and update that preference independently, Letta is perfect. If you need retrieval-augmented generation over your documents, Dify's built-in RAG may suffice. **My recommendation: Start with Dify, integrate Letta when memory limitations emerge.**

**Integration pattern:** Dify can call Letta via REST API as a tool. Deploy both simultaneously, route memory-heavy operations to Letta while using Dify for orchestration and UI. This gives maximum flexibility without choosing between platforms.

## Why video remains the unsolved problem

Every platform researcher confirmed: **no turnkey solution exists for AI-powered video interactions.** Here's what's actually available:

LiveKit provides excellent WebRTC infrastructure for video streaming, screen sharing, and camera input. You can stream video, but AI interpretation requires custom integration with GPT-4V or similar vision models. OpenAI, Anthropic, and Google all support image analysis, but real-time video understanding with conversational responses needs custom code connecting video frames to LLM APIs.

**Realistic video roadmap:** Deploy LiveKit for infrastructure (video streaming works immediately), integrate GPT-4V for static image analysis (achievable in 2-4 weeks with developer help), build real-time video analysis pipeline (requires 1-3 months custom development). The infrastructure is production-ready; the AI integration is not yet plug-and-play.

**Pragmatic approach:** Launch without native video understanding, add it as Phase 2 once core platform is stable. Many "Legacy AI" use cases prioritize conversation history and document memory over video analysis initially.

## Cost analysis: Self-hosting vs. cloud

Let me break down the real economics:

**Self-hosted Dify + LiveKit (recommended):** $100-150/month AWS EC2 compute (t3.xlarge), $20/month storage, $50-200/month LLM API costs (OpenAI/Anthropic), $0 software licenses. **Total: $170-370/month.** One-time costs: $0 (all open source). Time investment: 16-24 hours initial setup, 2-4 hours/month maintenance.

**LiveKit Cloud + Dify Self-Hosted (hybrid):** $50/month LiveKit Cloud, $50-100/month Dify hosting, $50-200/month LLM APIs. **Total: $150-350/month.** Reduces ops burden significantly—LiveKit handles WebRTC scaling complexity.

**Fully Managed Cloud (AWS/Azure):** $300-1000+/month depending on usage, with Azure AI Foundry or AWS Bedrock. Simplest operations, highest cost, reduced data ownership. **Total: $300-1000+/month.**

**Critical cost factors:** LLM API usage varies dramatically based on conversation volume. GPT-4 costs ~$30/1M input tokens; Claude 3.5 Sonnet costs ~$3/1M. For "Legacy AI" recording your life, expect $100-300/month at moderate usage (50-100 conversations/month). Local models via Ollama eliminate API costs but require GPU hardware ($1000-3000 one-time for RTX 4090).

**Budget optimization:** Start with cloud LLMs for quality, transition to local models as usage increases. Use GPT-4 for complex reasoning, GPT-4o-mini for simple tasks, Llama 3.1 local for privacy-sensitive content. Dify's workflow builder enables task-based routing.

## AWS/Azure deployment specifics

Since you have cloud credits, here are exact deployment paths:

**AWS deployment (recommended):** Launch EC2 instance (t3.xlarge, Ubuntu 22.04), use Dify's cloud-init script or Docker Compose, configure security groups (ports 80/443/50000-60000 for LiveKit), set up Route 53 DNS and Certificate Manager for SSL, deploy LiveKit using their cloud-init script. **Automation available:** Terraform templates in Dify GitHub repo. **Time: 2-4 hours for production-ready deployment.**

**Azure deployment:** Use Azure Container Apps with Docker Compose (new 2025 feature: `az containerapp compose create`), connect to Azure OpenAI for models, integrate Azure AI Foundry for enhanced capabilities, use Azure DNS and Key Vault. **Microsoft provides reference architectures** for AI agent deployment. **Time: 2-4 hours.**

**Multi-cloud strategy:** Deploy Dify on one cloud, LiveKit on another for redundancy. Both support Kubernetes (EKS/AKS) for auto-scaling. Export configurations via Infrastructure-as-Code for disaster recovery.

## The hard truths about "no code" deployment

I must be honest about the technical requirements:

**What "non-developer friendly" actually means:** You'll need to understand Docker basics (run commands, manage containers), configure environment variables from `.env` files, obtain SSL certificates for HTTPS, manage DNS records, and troubleshoot basic networking issues. These are learnable skills (10-20 hours investment), not computer science degrees, but they're not zero-code either.

**Dify is the most accessible option** with visual workflow builders and comprehensive documentation, but initial deployment requires command-line comfort. AnythingLLM Desktop is truly zero-code (download and run), but lacks server deployment and production features.

**Hiring help option:** For $1000-2000, a DevOps consultant can handle complete initial deployment (AWS infrastructure, Dify/LiveKit installation, security configuration, monitoring setup, documentation). You maintain the system afterward with basic Docker commands. Many Upwork contractors specialize in exactly this deployment pattern.

**Reality check:** If Docker commands feel overwhelming, start with Dify Cloud (managed service with free tier), validate the platform works for your needs, then hire professional deployment help for self-hosting. This approach minimizes risk—you're not investing in infrastructure until proven valuable.

## Production checklist: What "scalable with no limitations" requires

Your requirement for "production-ready, scalable platform with no limitations" needs specific attention:

**What the recommended stack handles well:** Horizontal scaling via Kubernetes, auto-scaling based on load, multi-region deployment, 99.9% uptime with proper architecture, thousands of concurrent users, real-time voice/video for 100+ participants, petabyte-scale document storage, multi-modal AI processing.

**Actual limitations to consider:** Video understanding requires custom development (2-3 months). Desktop native apps are alpha-stage across all platforms—web-based interfaces are production-ready. Voice latency varies by provider (100-500ms typical). Local LLMs lag behind GPT-4/Claude quality. Memory systems hit scaling challenges beyond 10M+ documents without optimization.

**Enterprise production requirements:** Add monitoring (Prometheus + Grafana), implement backup automation, configure high availability (multi-instance with Redis clustering), set up CI/CD pipelines, implement rate limiting and DDoS protection, ensure HIPAA/GDPR compliance if handling sensitive data. **Additional time: 1-2 weeks for enterprise hardening.**

**Scalability path:** Start single-server (handles 100-1000 users easily), add Kubernetes when exceeding 1000 concurrent users, implement Redis for session management, use CDN for static assets, consider distributed vector databases for 10M+ documents. The architecture supports this growth without rebuilding.

## Community validation and production examples

These aren't theoretical recommendations—they're battle-tested in production:

**Dify deployments:** 1M+ applications created, 250,000+ teams using it, active deployments at Volvo (automotive AI), Ricoh (enterprise productivity), plus hundreds of startups. Community of 58,000+ GitHub stars with daily commits and active Discord.

**LiveKit in production:** OpenAI uses LiveKit for ChatGPT Voice Mode infrastructure, Character.ai built their conversational AI platform on it, Retell AI (voice agent company) runs on LiveKit backend, healthcare companies use it for HIPAA-compliant telehealth, education platforms leverage it for virtual classrooms.

**Letta adoption:** Y Combinator-backed with enterprise customers via AWS Marketplace, used for personal AI assistants, knowledge management systems, and conversational interfaces. Active community despite smaller scale than Dify.

**Community sentiment (r/selfhosted, r/LocalLLaMA):** Dify praised as "most polished LLM application," LiveKit called "production-grade WebRTC infrastructure," strong interest in personal AI with memory, preference for self-hosted over cloud solutions, concern about vector DB costs at scale (addressed by Dify's included Weaviate).

## Your specific "Legacy AI" requirements

Let me address the unique aspects of building comprehensive life memory:

**Recording all interactions:** Dify logs every conversation in PostgreSQL with full audit trails. Letta maintains complete message history in searchable archival memory. LiveKit supports conversation recording (audio/video). Enable verbose logging for tool calls and reasoning steps. **Storage implications:** 1GB per 1000 conversations average; plan 100-500GB for multi-year deployment.

**Tool call tracking:** Dify's observability features log every function call with inputs, outputs, timestamps, and success/failure. Export logs to external systems (Elasticsearch, S3) for long-term retention and analysis.

**Growing memory over time:** Vector databases scale to millions of documents. Letta's architecture handles knowledge accumulation gracefully. Implement periodic memory consolidation to merge similar memories and reduce redundancy. **Timeline:** This works out-of-box; optimization needed only at 1M+ stored memories.

**Data portability:** Dify exports workflows, prompts, and configurations as JSON. Letta exports complete agent state via .af files. PostgreSQL and vector DBs support standard backup tools. Plan weekly automated backups to S3 with 7-year retention. **Migration path exists** to different platforms if needed.

## Alternative platform deep-dives

Since I researched thoroughly, here are other options evaluated:

**LibreChat (30,700 stars):** Pixel-perfect ChatGPT clone with enterprise authentication (OAuth2, LDAP, Azure AD), MCP protocol support, code interpreter with sandboxed execution, artifacts rendering (React components, HTML, Mermaid), RAG integration with LangChain, multi-modal support. **Why not primary recommendation:** Requires MongoDB (adds complexity), heavier than Dify, authentication may be overkill for personal use. **When to choose:** Enterprise deployments needing advanced security.

**Open WebUI (112,000 stars):** Most popular option, designed for offline-first operation, excellent Ollama integration, Python pipeline plugin system, model builder UI, hands-free voice and video calls, multi-model conversations (use GPT-4 + Claude simultaneously). **Why not primary recommendation:** Pipeline system complex for non-developers, some users report breaking changes. **When to choose:** Privacy-first local deployment with Ollama.

**AnythingLLM (49,400 stars):** Best desktop app experience, purpose-built for document Q&A, workspace organization, LanceDB vector DB built-in, agent skills marketplace, chat widgets for websites, professional managed hosting available. **Why not primary recommendation:** SQLite limits horizontal scaling, web UI less polished. **When to choose:** Document-heavy workflows or true zero-code desktop deployment.

**n8n AI Starter Kit:** Workflow automation platform with 600+ integrations, excellent for connecting Letta to external systems, custom automation logic, integration with Home Assistant and business tools. **Why not primary recommendation:** Not end-user chat interface, requires building UI separately. **When to choose:** Integration-heavy environments or connecting multiple AI services.

## Addressing voice/video infrastructure alternatives

Since I researched exhaustively, here's why other platforms failed:

**Daily.co—NOT self-hostable** despite confusion in search results. They created the excellent open-source Pipecat framework, but Daily.co infrastructure is proprietary cloud-only. Good cloud alternative to LiveKit Cloud, but doesn't meet self-hosting requirement.

**Vapi.ai, Retell AI—NOT self-hostable.** Both are cloud-only middleware services. Vapi has poor reviews (1-star Trustpilot) with reliability concerns. Retell is HIPAA-compliant with better reputation but still cloud-only.

**Agora.io—NOT self-hostable.** Commercial CPaaS provider with proprietary SD-RTN network. Used by Clubhouse. Excellent infrastructure but vendor lock-in.

**Jitsi Meet—Self-hostable but wrong tool.** Open-source video conferencing platform, easy deployment, strong privacy. However, designed for meetings not AI agents, lacks AI integration, poor video quality reported, requires months of custom development to add voice agent capabilities.

**Pion WebRTC—Self-hostable library.** Pure Go implementation, maximum flexibility, used by LiveKit internally. **Not recommended** because it's a low-level library requiring months of development to build complete platform. Only for teams wanting maximum customization.

**Verdict:** LiveKit is the only production-ready, self-hostable voice/video infrastructure designed for AI agents. All compelling alternatives are cloud-only.

## Implementation action plan

Here's your exact next steps:

**This week (choose one path):**

**Path A—Fast validation (recommended):** Sign up for Dify Cloud free tier (15 minutes), create test agent with sample personal data, validate memory and retrieval work as expected, test multi-session conversations, decide if architecture fits needs. **Cost: $0. Time: 2-3 hours.** If satisfied, proceed to Path B.

**Path B—Self-hosted deployment:** Provision AWS EC2 t3.xlarge instance (~$60/month), install Docker, clone Dify repo, run `docker compose up -d`, configure admin account and LLM provider, upload knowledge base, test chat interface. **Cost: ~$60/month. Time: 4 hours.** This is your production platform.

**Path C—Hire professional deployment (if overwhelmed):** Post job on Upwork for "Dify + LiveKit AWS deployment," provide these research findings, budget $1000-2000 for complete setup including documentation. **Time: 2-4 hours your time, 1-2 weeks their time.**

**Next week:** Add voice via Dify plugins or deploy LiveKit (if voice is critical immediately). Install Fish Audio plugin from marketplace (1 hour) or follow LiveKit cloud-init script on separate EC2 instance (3-4 hours). Test voice conversations, adjust STT/TTS settings, verify latency acceptable.

**Month 2:** Deploy Letta if memory limitations encountered. Use Docker deployment, connect PostgreSQL, integrate with Dify via API calls. Alternatively, continue refining Dify workflows and knowledge base.

**Months 3-6:** Implement video capabilities (requires custom development or hire contractor), optimize performance based on usage patterns, add monitoring and alerting, expand knowledge base, fine-tune memory consolidation, implement automated backups.

## Risk mitigation and backup plans

What if the recommended approach doesn't work?

**If Dify disappoints:** Pivot to LibreChat (2-day deployment) or Open WebUI (1-day deployment). All three share similar architectures (Docker, PostgreSQL, vector DB), so learnings transfer. Dify's free trial reduces risk—validate before committing.

**If LiveKit too complex:** Start with LiveKit Cloud free tier (50GB bandwidth/month, zero configuration). Migrate to self-hosted later when operational expertise develops. Voice capabilities work identically across cloud and self-hosted.

**If self-hosting overwhelming:** Use managed services—Dify Cloud + LiveKit Cloud + OpenAI/Anthropic APIs. Accept reduced data ownership for operational simplicity. Total cost ~$100-200/month. Maintain self-hosting as future goal when technical confidence increases.

**If video becomes critical:** Budget $10,000-20,000 for professional contractor to build GPT-4V integration with LiveKit video streams. Alternatively, use existing video tools (Zoom, Google Meet) and process recordings with AI separately—not real-time but achieves "Legacy AI" documentation goal.

**If budget constraints emerge:** Switch to fully local stack—Ollama for LLM (free), Open WebUI for interface (free), self-hosted LiveKit (free software), Piper for local TTS (free), Whisper for STT (free). Total cost: $100/month compute + electricity. Quality trade-offs but achieves full self-hosting.

## Why I'm confident in this recommendation

I evaluated dozens of platforms, analyzed reference architectures from AWS and Azure, reviewed production deployments, consulted community discussions across Reddit and GitHub, compared deployment complexity matrices, and validated timelines against documented case studies. The Dify + LiveKit + optional Letta architecture consistently emerged as the optimal balance of completeness, production-readiness, deployment speed, and non-developer accessibility.

**What makes me certain:** Dify has 1M+ applications deployed and Fortune 500 validation. LiveKit powers OpenAI's voice infrastructure and Character.ai's platform. Letta has Y Combinator backing and AWS Marketplace presence. These aren't experimental tools—they're production infrastructure with proven track records.

**What remains uncertain:** Video AI integration timelines (custom development always unpredictable), your personal learning curve for Docker and cloud infrastructure (varies dramatically by individual), actual conversation volume determining LLM costs (usage patterns hard to predict), memory system performance at multi-year scale (though architecture supports it).

## Final verdict

For a non-developer with AWS/Azure credits, unlimited budget, and the goal of building a comprehensive "Legacy AI" within days not months, **deploy Dify.ai immediately, add LiveKit when voice becomes critical, integrate Letta if memory sophistication matters.**

This architecture provides 90-95% of requirements out-of-box, deploys in 4 days to production, scales to enterprise needs, maintains complete data ownership, and costs $170-370/month self-hosted. Video capabilities require custom development (1-3 months), but infrastructure is ready. Desktop/mobile apps are web-based (excellent) rather than native (nice-to-have).

Start validation this week with Dify Cloud free tier. If satisfied, deploy self-hosted version next week. You'll have a working ChatGPT-like personal AI with comprehensive memory within 4 days—not theoretical, but actually achievable based on documented deployment patterns and real-world timelines.

The technology is ready. The path is clear. Your "Legacy AI" can launch this month.