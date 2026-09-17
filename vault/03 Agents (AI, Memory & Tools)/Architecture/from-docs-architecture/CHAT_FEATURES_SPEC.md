# OneMind Chat - Comprehensive Features Specification

> Research compiled from: Claude, ChatGPT, Grok, LobeChat, Lindy.ai, Relevance AI, Motion

---

## Executive Summary

This document outlines all features discovered across 7 leading AI chat platforms to guide the development of OneMind's chat system. The goal is to build a best-in-class chat experience that combines the strongest elements from each platform.

---

## Platform Feature Matrix

| Feature | Claude | ChatGPT | Grok | LobeChat | Lindy | Relevance | Motion |
|---------|--------|---------|------|----------|-------|-----------|--------|
| Multi-model support | Yes | Yes | No | Yes | Yes | Yes | No |
| Artifacts/Canvas | Yes | Yes | No | Yes | No | No | No |
| Voice chat | Yes | Yes | Yes | Yes | Yes | Yes | No |
| Agent marketplace | No | GPTs | No | Yes | Yes | Yes | No |
| Branching conversations | No | No | No | Yes | No | No | No |
| Memory/Context | Yes | Yes | No | Yes | No | No | Yes |
| Plugin system | Skills | Plugins | No | MCP | 1000+ | Tools | Apps |
| Projects/Workspaces | Yes | Yes | No | Yes | No | Yes | Yes |
| Knowledge base | Yes | Yes | No | Yes | Yes | Yes | Yes |
| Scheduling/Automation | No | No | No | Yes | Yes | Yes | Yes |
| Chain of Thought viz | No | No | Think | Yes | No | No | No |
| Self-hostable | No | No | No | Yes | No | No | No |

---

## 1. CLAUDE (Anthropic)

### Core Features
- **Multi-model Selection**: Opus (powerful), Sonnet (balanced), Haiku (fast)
- **Artifacts**: Interactive outputs - code, documents, diagrams, checklists
- **Projects**: Organize conversations by context with shared knowledge
- **Cowork**: Local file integration (macOS app)
- **Skills Marketplace**: Pre-built capabilities

### Chat UI Elements
- Personalized greeting with user name
- Quick action chips: Code, Learn, Strategize, Write, Life stuff
- Starred conversations sidebar
- Model selector dropdown
- Clean dark theme with warm accents

### Unique Features
- Extended thinking (chain of thought for complex reasoning)
- Google Drive integration
- Web search capability
- Style/voice customization
- Task-specific prompt templates

---

## 2. CHATGPT (OpenAI)

### Core Features
- **GPTs**: Custom AI assistants with specific capabilities
- **Canvas**: Side-by-side editing for code and documents
- **Memory**: Persistent user preferences and context
- **Projects**: Workspace organization
- **Voice Mode**: Real-time voice conversations
- **Vision**: Image analysis and generation (DALL-E)

### Chat UI Elements
- Minimal center area, conversation-focused
- Projects in left sidebar
- GPTs section for custom assistants
- Simple input bar with attachments
- Dark purple/black theme

### Unique Features
- Custom GPT builder
- Plugin ecosystem
- Advanced Data Analysis
- Browse with Bing
- DALL-E image generation
- Code Interpreter

---

## 3. GROK (xAI)

### Core Features
- **DeepSearch**: Advanced web research with source synthesis
- **Think Mode**: Step-by-step reasoning visualization
- **Aurora**: Image generation
- **Expert Mode**: Enhanced capabilities
- **Real-time X Integration**: Access to live social data

### Chat UI Elements
- Ultra-minimal design
- Single prominent input field
- Mode selector dropdown (Expert)
- X account connection
- SuperGrok upgrade CTA

### Unique Features
- Real-time information from X/Twitter
- Humor/personality in responses
- Uncensored responses (compared to competitors)
- Image understanding and generation

---

## 4. LOBECHAT (Open Source) - PRIMARY REFERENCE

### Core Architecture
```
┌─────────────────────────────────────────────────────────────────┐
│                        LOBECHAT HUB                              │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   AGENTS    │  │   PLUGINS   │  │   MODELS    │              │
│  │  Marketplace│  │  (MCP/Tools)│  │  Providers  │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
├─────────────────────────────────────────────────────────────────┤
│                     CONVERSATION ENGINE                          │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Branching   │  Chain of   │  Artifacts  │  File/KB     │   │
│  │  Threads     │  Thought    │  Rendering  │  Integration │   │
│  └──────────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────────┤
│                      AGENT SYSTEM                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   Agent     │  │   Agent     │  │   Personal  │              │
│  │   Builder   │  │   Groups    │  │   Memory    │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
├─────────────────────────────────────────────────────────────────┤
│                    COLLABORATION LAYER                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   Pages     │  │  Schedule   │  │  Projects   │              │
│  │(Multi-agent)│  │  (Cron)     │  │ (Workspaces)│              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

### Agent System
- **Agents as Unit of Work**: Each agent = specific capability/persona
- **Agent Builder**: Auto-configuration from natural language
- **Agent Groups**: Multi-agent collaboration, parallel execution
- **Personal Memory**: Structured, editable, continual learning
- **10,000+ Skills**: Extensive tool library via MCP

### Conversation Features
- **Branching Conversations**: Tree-like discussion structures
  - Continue mode: Extend from any point
  - Standalone mode: Independent branch
- **Chain of Thought (CoT)**: Step-by-step reasoning visualization
- **Artifacts Support**:
  - SVG graphics rendering
  - Interactive HTML pages
  - Live documents
  - Code execution previews

### Plugin System (MCP)
- **Model Context Protocol**: Standardized tool integration
- **One-click Installation**: From marketplace
- **10,000+ Tools**: Functions, APIs, automations
- **MCP Marketplace**: Curated, verified integrations
- **Custom Plugins**: Build your own

### Model Support
- **Multi-Provider**: OpenAI, Anthropic, Google, Mistral, etc.
- **Local LLMs**: Ollama integration
- **Visual Recognition**: GPT-4 Vision, Claude Vision
- **Text-to-Image**: DALL-E 3, Midjourney, Pollinations
- **TTS/STT**: Voice conversations

### Data & Knowledge
- **File Upload**: Documents, images, code
- **Knowledge Base**: Searchable document storage
- **Local Database**: CRDT for offline sync
- **Remote Database**: PostgreSQL for server
- **Vector Search**: Semantic retrieval

### Collaboration
- **Pages**: Multi-agent content creation with shared context
- **Schedule**: Automated task execution (cron-like)
- **Projects**: Organized workspaces
- **Workspaces**: Team collaboration with ownership

### UI Components
- **Custom Themes**: Light/dark, color customization
- **Chat Modes**: Bubble or document style
- **PWA**: Offline-capable web app
- **Desktop App**: Native performance
- **Mobile Optimized**: Responsive design

---

## 5. LINDY.AI

### Core Features
- **No-Code Agent Builder**: Build agents via natural language
- **Prompt-to-Agent**: Describe needs, get functional agents
- **Voice Agents**: State-of-the-art phone AI
- **Computer Use**: Virtual machine for complex tasks
- **1000+ Integrations**: Third-party app connections

### Agent Types
- Support Agent (tickets, chat, email)
- Inbound SDR (lead qualification)
- Creative Agent (marketing campaigns)
- Document Processing Agent
- Meeting Recorder
- Email Assistant
- Phone Call Agent

### Unique Features
- Meeting joins and recording
- Lead research and enrichment
- Ticket routing and resolution
- Knowledge base search
- Human escalation triggers

---

## 6. RELEVANCE AI

### Core Features
- **Multi-Agent System (MAS)**: Coordinated agent operations
- **Version Control**: Track, test, restore agents
- **Scheduling**: Control execution timing
- **Approvals/Escalation**: Human-in-the-loop workflows
- **Phone Agents**: Voice-calling capabilities

### Agent Templates
- AI BDR Agent (sales pipeline)
- Research Agent (account research)
- Sales Notetaker
- CRM Enrichment Agent
- Lifecycle Marketer
- SEO Agent
- Customer Support Agent

### Tool Builder
- Custom API integrations
- LLM prompt chains
- Traditional automations
- Custom code execution
- Bulk dataset operations

---

## 7. MOTION

### Core Features
- **Context-Aware Chat**: Knows your docs, tasks, projects, calendar
- **AI Calendar Assistant**: Auto-schedules and optimizes
- **AI Project Manager**: Generates projects from descriptions
- **AI Task Planner**: Creates and prioritizes tasks

### Unique Value
- Deep integration with work management tools
- Hyper-personalized without prompt engineering
- Real-time re-planning
- Meeting notes become docs and tasks
- Cross-app task consolidation

---

## ONEMIND CHAT ARCHITECTURE

Based on the research, here's the proposed architecture:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           ONEMIND CHAT SYSTEM                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         CHAT INTERFACE LAYER                         │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                      │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │    │
│  │  │   Welcome    │  │   Model      │  │   Context    │               │    │
│  │  │   Screen     │  │   Selector   │  │   Indicator  │               │    │
│  │  │  (greeting,  │  │  (Opus,      │  │  (project,   │               │    │
│  │  │  quick acts) │  │  Sonnet...)  │  │  team, mode) │               │    │
│  │  └──────────────┘  └──────────────┘  └──────────────┘               │    │
│  │                                                                      │    │
│  │  ┌──────────────────────────────────────────────────────────────┐   │    │
│  │  │                    MESSAGE DISPLAY AREA                       │   │    │
│  │  │  ┌────────────────────────────────────────────────────────┐  │   │    │
│  │  │  │  Branching Navigator (tree view, branch switch)        │  │   │    │
│  │  │  ├────────────────────────────────────────────────────────┤  │   │    │
│  │  │  │  Message Bubbles                                       │  │   │    │
│  │  │  │  - User messages (with attachments)                    │  │   │    │
│  │  │  │  - Assistant messages (with actions)                   │  │   │    │
│  │  │  │  - Chain of Thought cards (collapsible)                │  │   │    │
│  │  │  │  - Tool call cards (with results)                      │  │   │    │
│  │  │  │  - Artifact renders (code, HTML, SVG, docs)            │  │   │    │
│  │  │  │  - HITL approval cards                                 │  │   │    │
│  │  │  └────────────────────────────────────────────────────────┘  │   │    │
│  │  └──────────────────────────────────────────────────────────────┘   │    │
│  │                                                                      │    │
│  │  ┌──────────────────────────────────────────────────────────────┐   │    │
│  │  │                      INPUT BAR                                │   │    │
│  │  │  ┌─────┐ ┌──────────────────────────────┐ ┌─────┐ ┌───────┐  │   │    │
│  │  │  │ +   │ │  Type message...             │ │ 🎤  │ │ Send  │  │   │    │
│  │  │  │Attach│ │  @mention  /slash           │ │Voice│ │       │  │   │    │
│  │  │  └─────┘ └──────────────────────────────┘ └─────┘ └───────┘  │   │    │
│  │  │  [Quick prompts: Code | Research | Write | Analyze | Plan]   │   │    │
│  │  └──────────────────────────────────────────────────────────────┘   │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                         SIDEBAR / NAVIGATION                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────┐                                                    │
│  │  + New Chat         │                                                    │
│  │  🔍 Search          │                                                    │
│  ├─────────────────────┤                                                    │
│  │  📌 STARRED         │                                                    │
│  │  - Important conv 1 │                                                    │
│  │  - Important conv 2 │                                                    │
│  ├─────────────────────┤                                                    │
│  │  📁 PROJECTS        │                                                    │
│  │  > Project Alpha    │                                                    │
│  │  > Project Beta     │                                                    │
│  ├─────────────────────┤                                                    │
│  │  🤖 AGENTS          │                                                    │
│  │  - Code Assistant   │                                                    │
│  │  - Research Bot     │                                                    │
│  │  - [+ Create Agent] │                                                    │
│  ├─────────────────────┤                                                    │
│  │  💬 RECENT          │                                                    │
│  │  - Chat 1 (2m ago)  │                                                    │
│  │  - Chat 2 (1h ago)  │                                                    │
│  │  - Chat 3 (yesterday│                                                    │
│  ├─────────────────────┤                                                    │
│  │  ⚡ QUICK ACTIONS   │                                                    │
│  │  [Hub] [Plugins]    │                                                    │
│  │  [Settings]         │                                                    │
│  └─────────────────────┘                                                    │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                           AGENT HUB (MARKETPLACE)                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  🏪 AGENT MARKETPLACE                              [Search...]       │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                      │    │
│  │  Categories: [All] [Productivity] [Code] [Research] [Creative]      │    │
│  │              [Business] [Personal] [Custom]                         │    │
│  │                                                                      │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │    │
│  │  │ 🧑‍💻 Code    │  │ 📊 Data     │  │ ✍️ Writer   │  │ 🔬 Research │  │    │
│  │  │ Assistant   │  │ Analyst     │  │ Pro         │  │ Agent      │  │    │
│  │  │             │  │             │  │             │  │            │  │    │
│  │  │ Reviews,    │  │ Analyzes    │  │ Writes any  │  │ Deep dives │  │    │
│  │  │ debugs,     │  │ data, makes │  │ content     │  │ into any   │  │    │
│  │  │ writes code │  │ charts      │  │ type        │  │ topic      │  │    │
│  │  │             │  │             │  │             │  │            │  │    │
│  │  │ ⭐ 4.9 (2k) │  │ ⭐ 4.8 (1k) │  │ ⭐ 4.7 (3k) │  │ ⭐ 4.9 (800)│  │    │
│  │  │ [Install]   │  │ [Install]   │  │ [Install]   │  │ [Install]  │  │    │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘  │    │
│  │                                                                      │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │    │
│  │  │ 📅 Calendar │  │ 📧 Email    │  │ 🏠 Home     │  │ 🎨 Design  │  │    │
│  │  │ Manager     │  │ Assistant   │  │ Assistant   │  │ Helper     │  │    │
│  │  │ ...         │  │ ...         │  │ ...         │  │ ...        │  │    │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘  │    │
│  │                                                                      │    │
│  │  [Load More...]                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                           PLUGIN SYSTEM (MCP)                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  🔌 PLUGINS & TOOLS                                [Search...]       │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                      │    │
│  │  INSTALLED                          AVAILABLE                        │    │
│  │  ┌─────────────────────┐           ┌─────────────────────┐          │    │
│  │  │ ✅ Web Search       │           │ GitHub Integration  │          │    │
│  │  │ ✅ File System      │           │ Notion Connector    │          │    │
│  │  │ ✅ Code Execution   │           │ Slack Integration   │          │    │
│  │  │ ✅ Home Assistant   │           │ Linear Sync         │          │    │
│  │  │ ✅ Calendar         │           │ Database Query      │          │    │
│  │  └─────────────────────┘           └─────────────────────┘          │    │
│  │                                                                      │    │
│  │  MCP SERVERS                                                         │    │
│  │  ┌─────────────────────────────────────────────────────────────┐    │    │
│  │  │ Server Name        │ Status    │ Tools │ Actions            │    │    │
│  │  │ filesystem         │ 🟢 Active │ 12    │ [Configure] [Stop] │    │    │
│  │  │ home-assistant     │ 🟢 Active │ 45    │ [Configure] [Stop] │    │    │
│  │  │ tavily             │ 🟡 Limited│ 5     │ [Configure] [Stop] │    │    │
│  │  │ memory             │ 🟢 Active │ 8     │ [Configure] [Stop] │    │    │
│  │  └─────────────────────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                           CONVERSATION ENGINE                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                                                                      │    │
│  │  BRANCHING SYSTEM                                                    │    │
│  │  ────────────────                                                    │    │
│  │                                                                      │    │
│  │  Main Thread ──●──●──●──●──●──●──●                                  │    │
│  │                      │                                               │    │
│  │                      └──●──●──● Branch A                            │    │
│  │                           │                                          │    │
│  │                           └──●──● Branch A.1                        │    │
│  │                                                                      │    │
│  │  Features:                                                           │    │
│  │  - Fork from any message                                            │    │
│  │  - Continue mode (carries context)                                  │    │
│  │  - Standalone mode (fresh start)                                    │    │
│  │  - Branch navigation (tree view)                                    │    │
│  │  - Merge branches back to main                                      │    │
│  │                                                                      │    │
│  │  CHAIN OF THOUGHT                                                    │    │
│  │  ─────────────────                                                   │    │
│  │                                                                      │    │
│  │  ┌─────────────────────────────────────────────────────────────┐    │    │
│  │  │ 🧠 Thinking...                                    [Collapse] │    │    │
│  │  ├─────────────────────────────────────────────────────────────┤    │    │
│  │  │ Step 1: Analyzing the request...                            │    │    │
│  │  │ Step 2: Identifying relevant context...                     │    │    │
│  │  │ Step 3: Formulating approach...                             │    │    │
│  │  │ Step 4: Generating response...                              │    │    │
│  │  └─────────────────────────────────────────────────────────────┘    │    │
│  │                                                                      │    │
│  │  ARTIFACTS SYSTEM                                                    │    │
│  │  ────────────────                                                    │    │
│  │                                                                      │    │
│  │  Supported Types:                                                    │    │
│  │  - 📝 Documents (Markdown, rich text)                               │    │
│  │  - 💻 Code (syntax highlighted, executable)                         │    │
│  │  - 🖼️ SVG Graphics (rendered inline)                                │    │
│  │  - 🌐 HTML Pages (sandboxed iframe)                                 │    │
│  │  - 📊 Charts/Diagrams (Mermaid, Chart.js)                           │    │
│  │  - 📋 Checklists (interactive)                                      │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                           MEMORY & CONTEXT SYSTEM                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                                                                      │    │
│  │  PERSONAL MEMORY                                                     │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │    │
│  │  │ User Prefs   │  │ Learned      │  │ Persistent   │               │    │
│  │  │ (style,tone) │  │ Facts        │  │ Context      │               │    │
│  │  └──────────────┘  └──────────────┘  └──────────────┘               │    │
│  │                                                                      │    │
│  │  CONTEXTUAL AWARENESS (Motion-style)                                 │    │
│  │  ┌──────────────────────────────────────────────────────────────┐   │    │
│  │  │  📅 Calendar: Next meeting in 30min - "Design Review"        │   │    │
│  │  │  📋 Tasks: 3 high priority items due today                   │   │    │
│  │  │  📁 Project: Currently in "OneMind v2" workspace             │   │    │
│  │  │  🏠 Home: Living room lights on, 72°F                        │   │    │
│  │  │  💪 Wellness: Energy level moderate, 6h sleep                │   │    │
│  │  └──────────────────────────────────────────────────────────────┘   │    │
│  │                                                                      │    │
│  │  KNOWLEDGE BASE                                                      │    │
│  │  ┌──────────────────────────────────────────────────────────────┐   │    │
│  │  │  Uploaded Documents: 45                                       │   │    │
│  │  │  Indexed Pages: 1,234                                         │   │    │
│  │  │  Vector Embeddings: Active                                    │   │    │
│  │  │  [Upload] [Manage] [Search]                                   │   │    │
│  │  └──────────────────────────────────────────────────────────────┘   │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                           VOICE & MULTIMODAL                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                                                                      │    │
│  │  VOICE CHAT                          IMAGE CAPABILITIES              │    │
│  │  ┌──────────────────────┐           ┌──────────────────────┐        │    │
│  │  │ 🎤 Speech-to-Text    │           │ 👁️ Vision Analysis    │        │    │
│  │  │ 🔊 Text-to-Speech    │           │ 🎨 Image Generation   │        │    │
│  │  │ 🗣️ Real-time Voice   │           │ 📸 Screenshot Input   │        │    │
│  │  │ 📞 Phone Calls       │           │ 📊 Chart Recognition  │        │    │
│  │  └──────────────────────┘           └──────────────────────┘        │    │
│  │                                                                      │    │
│  │  FILE ATTACHMENTS                                                    │    │
│  │  Supported: PDF, DOCX, XLSX, Images, Code files, Audio, Video       │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## COMPLETE FEATURE CHECKLIST FOR ONEMIND

### Phase 1: Core Chat (MVP)
- [ ] Multi-model selector (route to different providers)
- [ ] Message display with markdown rendering
- [ ] Code blocks with syntax highlighting
- [ ] Input bar with send/voice buttons
- [ ] Attachment support (images, files)
- [ ] Conversation history (sidebar)
- [ ] New chat creation
- [ ] Chat search

### Phase 2: Enhanced UX
- [ ] Personalized greeting screen
- [ ] Quick action chips (Code, Research, Write, etc.)
- [ ] Starred/pinned conversations
- [ ] Projects/workspaces organization
- [ ] Model selector dropdown
- [ ] Token counter
- [ ] Typing indicators
- [ ] Message reactions/actions

### Phase 3: Advanced Conversation
- [ ] Branching conversations (fork from any message)
- [ ] Branch navigator (tree visualization)
- [ ] Chain of Thought visualization
- [ ] Reasoning step cards (collapsible)
- [ ] Artifacts rendering:
  - [ ] Code artifacts (editable, runnable)
  - [ ] Document artifacts (Markdown)
  - [ ] SVG/diagram artifacts
  - [ ] HTML artifacts (sandboxed)
  - [ ] Chart artifacts

### Phase 4: Agent System
- [ ] Agent marketplace/hub
- [ ] Agent cards with ratings
- [ ] One-click agent installation
- [ ] Agent builder (no-code)
- [ ] Agent configuration editor
- [ ] Agent groups (multi-agent)
- [ ] @mention agents in chat
- [ ] Slash commands for agents

### Phase 5: Plugin System (MCP)
- [ ] Plugin marketplace
- [ ] One-click plugin installation
- [ ] MCP server management
- [ ] Tool call visualization
- [ ] Plugin configuration UI
- [ ] Custom plugin builder
- [ ] Plugin permissions

### Phase 6: Memory & Context
- [ ] Personal memory (user preferences)
- [ ] Learned facts storage
- [ ] Contextual awareness panel:
  - [ ] Calendar integration
  - [ ] Task integration
  - [ ] Project context
  - [ ] Home state (HA)
  - [ ] Wellness data
- [ ] Knowledge base management
- [ ] File upload & indexing
- [ ] Vector search

### Phase 7: Voice & Multimodal
- [ ] Speech-to-text input
- [ ] Text-to-speech output
- [ ] Real-time voice mode
- [ ] Image analysis (vision)
- [ ] Image generation
- [ ] Screenshot input
- [ ] Video analysis

### Phase 8: Collaboration
- [ ] Pages (multi-agent docs)
- [ ] Scheduled tasks (cron)
- [ ] Shared workspaces
- [ ] Team permissions
- [ ] Export/share conversations

---

## UI/UX PATTERNS BY PLATFORM

### Welcome Screen Patterns

**Claude Style:**
```
┌─────────────────────────────────────────┐
│  🌟 Good evening, Master Delacruz      │
│                                         │
│  [____________________________]         │
│                                         │
│  [Code] [Learn] [Strategize] [Write]   │
└─────────────────────────────────────────┘
```

**ChatGPT Style:**
```
┌─────────────────────────────────────────┐
│                                         │
│                                         │
│         (empty, minimal)                │
│                                         │
│  [Ask ChatGPT___________________]       │
└─────────────────────────────────────────┘
```

**LobeChat Style:**
```
┌─────────────────────────────────────────┐
│  ✨ Everything's ready                  │
│                                         │
│  [Agent] [Group] [Write] [Analyze]     │
│                                         │
│  Recent Topics:                         │
│  [Card 1] [Card 2] [Card 3] [Card 4]   │
└─────────────────────────────────────────┘
```

**Recommended for OneMind:**
```
┌─────────────────────────────────────────┐
│  🧠 Welcome back, Zeus                  │
│     Ready to think together             │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │ What's on your mind?            │   │
│  │ @agent /command                 │   │
│  └─────────────────────────────────┘   │
│                                         │
│  [💻 Code] [🔍 Research] [✍️ Write]    │
│  [📊 Analyze] [🏠 Home] [📅 Plan]      │
│                                         │
│  ── Context ──────────────────────     │
│  📅 Design Review in 30min             │
│  ✓ 2 tasks completed today             │
│  🏠 Home: All systems normal           │
└─────────────────────────────────────────┘
```

---

## TECHNICAL IMPLEMENTATION NOTES

### State Management
- Use Riverpod for reactive state
- Separate providers for:
  - Chat history
  - Current conversation
  - Branch state
  - Agent state
  - Plugin state
  - Memory/context

### Data Models
```dart
// Conversation with branching
class Conversation {
  String id;
  String title;
  List<Branch> branches;
  String activeBranchId;
  DateTime created;
  DateTime updated;
  String? projectId;
  bool starred;
}

class Branch {
  String id;
  String? parentBranchId;
  int forkFromMessageIndex;
  List<Message> messages;
}

class Message {
  String id;
  String role; // user, assistant, system, tool
  String content;
  List<Attachment>? attachments;
  List<ToolCall>? toolCalls;
  ThinkingSteps? thinking;
  List<Artifact>? artifacts;
  DateTime timestamp;
}

class Artifact {
  String id;
  String type; // code, document, svg, html, chart
  String title;
  String content;
  String? language;
  bool editable;
}
```

### API Integration
- Support multiple providers via adapter pattern
- Streaming responses for real-time display
- Tool calling with MCP protocol
- Context injection from memory system

---

## LOBE-UI COMPONENT LIBRARY

LobeChat uses a dedicated UI library (`@lobehub/ui`) with 100+ components optimized for AI chat interfaces. Here's the full component mapping for Flutter implementation:

### Layout & Structure Components
| Lobe-UI (React) | Flutter Equivalent | Purpose |
|-----------------|-------------------|---------|
| `Layout` | `Scaffold` | Page structure |
| `Header` | `AppBar` / Custom | Top navigation |
| `Footer` | `BottomAppBar` / Custom | Bottom section |
| `Sidebar` | `Drawer` / `NavigationRail` | Side navigation |
| `Grid` | `GridView` | Grid layouts |
| `Flexbox` | `Row` / `Column` / `Flex` | Flexible layouts |
| `Block` | `Container` | Content blocks |

### Forms & Input Components
| Lobe-UI (React) | Flutter Equivalent | Purpose |
|-----------------|-------------------|---------|
| `Input` | `TextField` | Text input |
| `TextArea` | `TextField(maxLines: null)` | Multi-line input |
| `InputPassword` | `TextField(obscureText)` | Password field |
| `InputNumber` | Custom `TextField` | Number input |
| `Checkbox` | `Checkbox` | Toggle option |
| `Radio` | `Radio` | Single selection |
| `Select` | `DropdownButton` | Dropdown select |
| `LobeSelect` | Custom styled dropdown | AI model selector |
| `DatePicker` | `showDatePicker` | Date selection |
| `AutoComplete` | `Autocomplete` | Suggestions |
| `Form` | `Form` | Form wrapper |
| `FormModal` | `Dialog` with `Form` | Modal form |

### Navigation Components
| Lobe-UI (React) | Flutter Equivalent | Purpose |
|-----------------|-------------------|---------|
| `Menu` | `PopupMenuButton` | Menu options |
| `Dropdown` | `DropdownButton` | Dropdown |
| `DropdownMenu` | `PopupMenuButton` | Context menu |
| `SideNav` | `NavigationRail` | Side navigation |
| `DraggableSideNav` | Custom with `Draggable` | Resizable sidebar |
| `Tabs` | `TabBar` / `TabBarView` | Tab navigation |
| `Breadcrumb` | Custom widget | Path breadcrumbs |

### Feedback & Modal Components
| Lobe-UI (React) | Flutter Equivalent | Purpose |
|-----------------|-------------------|---------|
| `Modal` | `showDialog` | Modal dialogs |
| `Drawer` | `Drawer` | Slide panels |
| `Popover` | `PopupMenuButton` | Popovers |
| `Tooltip` | `Tooltip` | Hover tips |
| `Toast` | `SnackBar` | Notifications |
| `Alert` | `AlertDialog` | Alerts |
| `Empty` | Custom empty state | No data view |
| `Skeleton` | `Shimmer` / `Skeleton` | Loading state |

### Data Display Components
| Lobe-UI (React) | Flutter Equivalent | Purpose |
|-----------------|-------------------|---------|
| `Table` | `DataTable` | Data tables |
| `List` | `ListView` | List display |
| `Avatar` | `CircleAvatar` | User avatars |
| `Badge` | `Badge` | Count badges |
| `Tag` | `Chip` | Tags/labels |
| `Card` | `Card` | Content cards |
| `Image` | `Image` | Image display |
| `Toc` | Custom widget | Table of contents |

### Action Components (CRITICAL FOR CHAT)
| Lobe-UI (React) | Flutter Equivalent | Purpose |
|-----------------|-------------------|---------|
| `Button` | `ElevatedButton` / `TextButton` | Actions |
| `ActionIcon` | `IconButton` | Icon actions |
| `CopyButton` | Custom with clipboard | Copy to clipboard |
| `DownloadButton` | Custom | Download files |
| `Hotkey` | `Shortcuts` widget | Keyboard shortcuts |

### Specialized AI/Chat Components (MOST IMPORTANT)
| Lobe-UI (React) | Flutter Equivalent | Purpose |
|-----------------|-------------------|---------|
| `CodeEditor` | `code_editor` package | Code editing |
| `Markdown` | `flutter_markdown` | Markdown rendering |
| `Highlighter` | `flutter_highlight` | Syntax highlighting |
| `EmojiPicker` | `emoji_picker_flutter` | Emoji selection |
| `ColorSwatches` | Custom color picker | Theme colors |
| `Mermaid` | `mermaid` package / WebView | Diagram rendering |
| `Video` | `video_player` | Video playback |
| `EditableText` | `TextField` / `EditableText` | Inline editing |
| `SearchBar` | `SearchBar` (M3) | Search input |
| `Segmented` | `SegmentedButton` | Segmented control |
| `Snippet` | Custom code block | Code snippets |
| `SortableList` | `ReorderableListView` | Drag-to-reorder |
| `SliderWithInput` | Custom compound widget | Slider + input |

### Provider Components
| Lobe-UI (React) | Flutter Equivalent | Purpose |
|-----------------|-------------------|---------|
| `ThemeProvider` | `Theme` / `ThemeData` | Theming |
| `MotionProvider` | Animation controllers | Animations |
| `LobeUIProvider` | Custom `InheritedWidget` | Global config |
| `I18nProvider` | `Localizations` | i18n |
| `ConfigProvider` | `InheritedWidget` | Settings |
| `IconProvider` | `IconTheme` | Icon theming |
| `ErrorBoundary` | `ErrorWidget` | Error handling |

---

## ONEMIND FLUTTER CHAT COMPONENTS TO BUILD

Based on lobe-ui, here are the custom components we need to create:

### Core Chat Components
```
lib/shared/widgets/chat/
├── chat_input_bar.dart          ✅ EXISTS
├── message_bubble.dart          ✅ EXISTS (enhanced_message_bubble.dart)
├── message_actions.dart         ✅ EXISTS
├── branch_navigator.dart        ✅ EXISTS
├── reasoning_step_card.dart     ✅ EXISTS
├── tool_call_card.dart          ✅ EXISTS
├── artifact_renderer.dart       🔨 NEEDS WORK
│   ├── code_artifact.dart
│   ├── document_artifact.dart
│   ├── svg_artifact.dart
│   ├── html_artifact.dart
│   └── chart_artifact.dart
├── chain_of_thought.dart        🔨 NEEDS WORK
├── agent_mention.dart           ❌ NEW
├── slash_command_menu.dart      ✅ EXISTS
├── quick_prompts.dart           ✅ EXISTS
├── model_selector.dart          ❌ NEW
├── context_indicator.dart       ❌ NEW
└── welcome_screen.dart          ❌ NEW
```

### Agent Hub Components
```
lib/agno/hub/
├── agent_marketplace.dart       ❌ NEW
├── agent_card.dart              ❌ NEW
├── agent_detail_sheet.dart      ❌ NEW
├── agent_builder.dart           ❌ NEW
├── agent_config_editor.dart     ❌ NEW
├── plugin_marketplace.dart      ❌ NEW
├── plugin_card.dart             ❌ NEW
└── mcp_server_manager.dart      ❌ NEW
```

### Sidebar Components
```
lib/shared/widgets/sidebar/
├── chat_sidebar.dart            ❌ NEW
├── conversation_list.dart       ❌ NEW
├── project_list.dart            ❌ NEW
├── agent_list.dart              ❌ NEW
├── starred_section.dart         ❌ NEW
└── quick_actions.dart           ❌ NEW
```

---

## DESIGN SYSTEM TOKENS

Based on LobeChat's design language:

### Colors (Dark Theme)
```dart
class OneMindColors {
  // Background
  static const background = Color(0xFF0A0A0A);
  static const surface = Color(0xFF141414);
  static const surfaceElevated = Color(0xFF1F1F1F);

  // Primary (Accent)
  static const primary = Color(0xFF6366F1);      // Indigo
  static const primaryLight = Color(0xFF818CF8);
  static const primaryDark = Color(0xFF4F46E5);

  // Semantic
  static const success = Color(0xFF22C55E);
  static const warning = Color(0xFFF59E0B);
  static const error = Color(0xFFEF4444);
  static const info = Color(0xFF3B82F6);

  // Text
  static const textPrimary = Color(0xFFFAFAFA);
  static const textSecondary = Color(0xFFA1A1AA);
  static const textTertiary = Color(0xFF71717A);

  // Border
  static const border = Color(0xFF27272A);
  static const borderLight = Color(0xFF3F3F46);
}
```

### Typography
```dart
class OneMindTypography {
  static const heading1 = TextStyle(
    fontSize: 32,
    fontWeight: FontWeight.w700,
    height: 1.2,
  );

  static const heading2 = TextStyle(
    fontSize: 24,
    fontWeight: FontWeight.w600,
    height: 1.3,
  );

  static const body = TextStyle(
    fontSize: 14,
    fontWeight: FontWeight.w400,
    height: 1.5,
  );

  static const code = TextStyle(
    fontFamily: 'JetBrains Mono',
    fontSize: 13,
    height: 1.6,
  );
}
```

### Spacing
```dart
class OneMindSpacing {
  static const xs = 4.0;
  static const sm = 8.0;
  static const md = 16.0;
  static const lg = 24.0;
  static const xl = 32.0;
  static const xxl = 48.0;
}
```

### Border Radius
```dart
class OneMindRadius {
  static const sm = 4.0;
  static const md = 8.0;
  static const lg = 12.0;
  static const xl = 16.0;
  static const full = 9999.0;
}
```

---

## REFERENCES

- LobeChat GitHub: https://github.com/lobehub/lobe-chat
- LobeChat Docs: https://lobehub.com/docs
- Lobe UI: https://github.com/lobehub/lobe-ui
- Lobe Icons: https://github.com/lobehub/lobe-icons
- Claude: https://claude.ai
- ChatGPT: https://chat.openai.com
- Grok: https://grok.com
- Lindy: https://lindy.ai
- Relevance AI: https://relevanceai.com
- Motion: https://usemotion.com
