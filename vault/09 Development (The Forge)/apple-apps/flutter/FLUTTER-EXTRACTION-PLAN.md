# Feature Extraction Plan: onemind-web + onemind-mobile → onemind-app

> Generated 2026-04-17. Complete audit of both source repos mapped to Flutter app targets.
> **Rule: DO NOT DELETE source repos until Zeus confirms.**

---

## Audit Summary

| Source | Scale | Language | Role |
|--------|-------|----------|------|
| **onemind-web** | 1,391 files · 34 pages · 185 API routes · 200 components · 55 DB tables | TypeScript/React | Full-stack web control plane |
| **onemind-mobile** | 3 native apps · 100+ source files | Swift/Kotlin | Native reference clients |
| **onemind-app** | 13 Dart files (skeleton) | Dart/Flutter | Target: unified mobile client |

### Architecture Decision

**onemind-web STAYS as the web control plane.** It has a full server backend (SQLite, LangGraph, 60+ agent tools, extensions, connectors). The Flutter app is the **mobile client** — it talks to OpenClaw gateway directly, same as the native apps do. The web app and Flutter app complement each other; they don't replace each other.

**onemind-mobile becomes reference-only.** All portable patterns get extracted to Flutter. The native apps stay in the repo for reference but get removed from the active workspace.

---

## What Already Exists in onemind-app (13 files)

| File | What It Does | Port Status |
|------|-------------|-------------|
| `core/gateway_client.dart` | WebSocket JSON-RPC (basic) | ⚠️ Needs Chowder protocol v3 upgrade |
| `core/rest_client.dart` | Dio REST + SSE | ⚠️ Needs VisionClaw headers |
| `core/auth.dart` | flutter_secure_storage | ✅ Done |
| `core/config.dart` | Gateway URLs | ✅ Done |
| `theme/onemind_colors.dart` | Domain/agent colors | ✅ Done |
| `theme/onemind_theme.dart` | Dark theme | ✅ Done |
| `app.dart` | 5-tab GoRouter shell | ✅ Done |
| `main.dart` | Riverpod entry | ✅ Done |
| `pages/*/` | 5 stub pages | ⚠️ Stubs only |

---

## TIER 1 — MUST PORT (Core functionality that exists in both sources)

### 1.1 WebSocket Protocol v3 (from Chowder ChatService.swift — 1062 lines)

**Current gap**: Our `gateway_client.dart` is simplified. Chowder has the full protocol.

| Feature | Source | Priority |
|---------|--------|----------|
| Connect handshake (challenge → connect req with protocol v3) | ChatService.swift:190-230 | 🔴 Critical |
| Device identity (stable UUID, device name, OS) | ChatService.swift:148-165 | 🔴 Critical |
| `chat.send` with idempotencyKey + deliver flag | ChatService.swift:280-310 | 🔴 Critical |
| History polling (1s interval during active runs) | ChatService.swift:350-400 | 🔴 Critical |
| Triple dedup (sequence + toolCallId + content hash) | ChatService.swift:420-480 | 🔴 Critical |
| Event routing (agent events: assistant/thinking/tool/lifecycle streams) | ChatService.swift:500-650 | 🔴 Critical |
| Workspace sync (detect IDENTITY.md/USER.md writes → refresh) | ChatService.swift:700-750 | 🟡 Important |
| Auto-reconnect with 3s delay | ChatService.swift:800-830 | 🟡 Important |
| TLS trust for .ts.net (Tailscale certs) | ChatService.swift (URLSession delegate) | 🟡 Important |
| Delegate pattern → Riverpod state notifiers | ChatService delegate methods | 🔴 Critical |

### 1.2 Chat UI (from Chowder Views + onemind-web components)

| Feature | Source | Priority |
|---------|--------|----------|
| Message bubble (user blue / assistant markdown) | Chowder: MessageBubbleView.swift | 🔴 Critical |
| Markdown rendering (headings, code blocks, lists, quotes, inline) | Chowder: MarkdownContentView.swift | 🔴 Critical |
| Code block with language label + copy button | Chowder + onemind-web: code-block.tsx | 🔴 Critical |
| Thinking shimmer indicator | Chowder: ThinkingShimmerView.swift | 🟡 Important |
| Agent activity step timeline | Chowder: AgentActivityCard.swift + ActivityStepRow.swift | 🟡 Important |
| Tool call bubble (step-by-step tool execution display) | onemind-web: tool-call-bubble.tsx, tool-events-section.tsx | 🟡 Important |
| Scroll to bottom FAB | Chowder + onemind-web | 🟢 Nice |
| Message pagination (50 per page, "Load Earlier") | Chowder: ChatViewModel.swift:displayLimit | 🟢 Nice |
| Haptic feedback on first streaming delta | Chowder: ChatViewModel.swift | 🟢 Nice |
| Copy message action | Chowder context menu | 🟢 Nice |
| Suggestions bar | onemind-web: suggestions-bar.tsx | 🟢 Nice |
| Reply quote | onemind-web: reply-quote.tsx | 🟢 Nice |

### 1.3 Data Models (from both sources)

| Model | Source | Fields |
|-------|--------|--------|
| **Message** | Chowder | id, role (user/assistant), content, timestamp |
| **BotIdentity** | Chowder | name, creature, vibe, emoji, avatar |
| **UserProfile** | Chowder | name, callName, pronouns, timezone, notes, context |
| **ConnectionConfig** | Chowder | gatewayURL, sessionKey, token |
| **AgentActivity** | Chowder | currentLabel, thinkingText, steps[] |
| **ActivityStep** | Chowder | type, label, detail, status, elapsedTime, toolCategory |
| **ToolCategory** | Chowder | thinking/terminal/browser/network/fileSystem/search/generic |
| **Agent** | onemind-web | id, name, emoji, domain, model, status, personality, skills |
| **Session** | onemind-web | id, agentId, sessionKey, lastMessage, unread, pinned |
| **ChatMessage** | onemind-web | id, role, content, toolCalls, toolEvents, usage, timestamp |
| **Provider** | onemind-web | id, type, name, apiKey, baseUrl, models |
| **Schedule** | onemind-web | id, agentId, cron, method, enabled |
| **Task** (Kanban) | onemind-web | id, title, status, priority, agentId, column |
| **Memory** | onemind-web | id, content, agentId, scope, tier, embedding |
| **Connector** | onemind-web | id, platform, name, config, health |
| **Extension** | onemind-web | id, name, enabled, hooks, tools |
| **MCP Server** | onemind-web | id, name, command, tools |
| **Wallet** | onemind-web | id, chain, address, balance |
| **Skill** | onemind-web | id, name, source, installed |
| **Webhook** | onemind-web | id, url, events, secret |
| **Secret** | onemind-web | id, key, value (encrypted reference) |
| **UsageRecord** | onemind-web | sessionId, tokens, cost, provider, model |

### 1.4 Settings & Configuration (from Chowder + onemind-web)

| Screen | Source | Fields |
|--------|--------|--------|
| Connection settings | Chowder SettingsView | Gateway URL, Token, Session Key |
| Agent identity | Chowder SettingsView | Name, Creature, Vibe, Emoji, Avatar photo |
| User profile | Chowder SettingsView | Name, Call Name, Pronouns, Timezone, Notes |
| Location sharing | Chowder SettingsView | Toggle, authorization, current address |
| Clear chat history | Chowder SettingsView | Destructive action |
| Provider config | onemind-web provider-sheet | API keys, base URLs, model selection |
| App settings | onemind-web settings | Theme, notifications, sync preferences |

---

## TIER 2 — SHOULD PORT (Valuable features from onemind-web)

### 2.1 Agent Management (from onemind-web — 16 components)

| Feature | Source Component | What It Does |
|---------|-----------------|-------------|
| Agent fleet list | agent-list.tsx | Grid/list of all agents with status |
| Agent detail card | agent-card.tsx | Avatar, name, domain, model, status |
| Agent editor sheet | agent-sheet.tsx | Full CRUD: persona, model, skills, permissions |
| Agent avatar | agent-avatar.tsx | Emoji + color badge |
| Agent chat list | agent-chat-list.tsx | Sessions for a specific agent |
| Personality builder | personality-builder.tsx | Soul/personality config |
| Soul library | soul-library-picker.tsx | Pre-built personality templates |
| Trash/archive | trash-list.tsx | Trashed agents recovery |
| Inspector panel | inspector-panel.tsx | Debug view: model, context, token usage |

### 2.2 Session Management (from onemind-web)

| Feature | Source | What It Does |
|---------|--------|-------------|
| Session list | sidebar session list | All active sessions with agent, last msg, timestamp |
| Session CRUD | /api/chats/* | Create, rename, pin, archive, delete sessions |
| Multi-agent chatrooms | chatroom-*.tsx (11 components) | Group chat with multiple agents |
| Chatroom reactions | reaction-picker.tsx | Emoji reactions on messages |
| Chatroom moderation | moderation actions | Mute, delete, role management |
| Message queue | queued-message-queue.ts | Offline message queuing |

### 2.3 Memory Viewer (from onemind-web — 6 components)

| Feature | Source | What It Does |
|---------|--------|-------------|
| Memory list | memory-list.tsx | All memories with search/filter |
| Memory cards | memory-card.tsx | Individual memory display |
| Memory detail | memory-detail.tsx | Full memory content + metadata |
| Memory graph | memory-graph-view.tsx | Entity relationship visualization |
| Memory by agent | memory-agent-list.tsx | Filter by agent |
| Memory editor | memory-sheet.tsx | Create/edit memory entries |

### 2.4 Task Board (from onemind-web — 4 components)

| Feature | Source | What It Does |
|---------|--------|-------------|
| Kanban board | task-list.tsx, task-column.tsx | Column-based task management |
| Task cards | task-card.tsx | Title, status, priority, agent assignment |
| Task editor | task-sheet.tsx | Full task CRUD |
| GitHub import | /api/tasks/github-import | Import tasks from GitHub issues |

### 2.5 Knowledge Base (from onemind-web — 2 components)

| Feature | Source | What It Does |
|---------|--------|-------------|
| Knowledge list | knowledge-list.tsx | RAG document library |
| Knowledge upload | knowledge-sheet.tsx | Upload + manage documents |

### 2.6 Skill Management (from onemind-web — 4 components)

| Feature | Source | What It Does |
|---------|--------|-------------|
| Skill list | skill-list.tsx | Installed skills |
| ClawhHub browser | clawhub-browser.tsx | Browse + install skills |
| Skill workspace | skills-workspace.tsx | Skill files + config |
| Skill editor | skill-sheet.tsx | Edit skill definition |

### 2.7 MCP Server Management (from onemind-web — 2 components)

| Feature | Source | What It Does |
|---------|--------|-------------|
| MCP server list | mcp-server-list.tsx | Configured MCP servers |
| MCP server editor | mcp-server-sheet.tsx | Add/configure/test MCP servers |

### 2.8 Schedule/Cron Management (from onemind-web)

| Feature | Source | What It Does |
|---------|--------|-------------|
| Schedule list | via data slice | All scheduled jobs |
| Cron form | cron-job-form.tsx | Create/edit cron expressions |
| Manual trigger | /api/schedules/trigger | Fire schedule now |

### 2.9 Usage & Cost Analytics (from onemind-web)

| Feature | Source | What It Does |
|---------|--------|-------------|
| Usage page | /usage page | Token usage + cost breakdown |
| Cost trend chart | cost-trend-chart.tsx | Historical cost visualization |
| Per-session usage | usage records | Tokens/cost per conversation |

### 2.10 Run History (from onemind-web)

| Feature | Source | What It Does |
|---------|--------|-------------|
| Run list | run-list.tsx | Execution history with events |
| Run events | runtime_run_events table | Step-by-step execution trace |

---

## TIER 3 — NICE TO PORT (Advanced features, port later)

### 3.1 Connectors (from onemind-web — 13 platforms)

Discord, Slack, Telegram, WhatsApp, Signal, Teams, Matrix, Google Chat, BlueBubbles (iMessage), Email, OpenClaw, Voice Note, Media. Mobile app could manage these through gateway API calls — doesn't need the server-side connector code.

### 3.2 Extensions System (from onemind-web)

30+ lifecycle hooks, custom tools, providers, connectors, UI injection. The Flutter app could manage extensions via the gateway API but doesn't need the extension runtime.

### 3.3 Protocols/Workflows (from onemind-web — DAG engine)

Template-based multi-step workflows. The Flutter app would view/trigger these, not execute them.

### 3.4 Missions (from onemind-web)

Objective-driven execution. Same — view/trigger via API.

### 3.5 Wallet System (from onemind-web — ETH/SOL)

Crypto wallet management. Could be a mobile-friendly interface for approving transactions.

### 3.6 Project Management (from onemind-web)

Project CRUD, agent assignment, tabs. Could be a simpler mobile view.

### 3.7 Webhooks (from onemind-web)

Webhook endpoint management. Manage via API.

### 3.8 Secrets Vault (from onemind-web)

Credential management. Manage via API.

### 3.9 Voice Pipeline (from VisionClaw)

| Feature | Source | Complexity |
|---------|--------|-----------|
| Gemini Live WebSocket audio | GeminiLiveService.swift | 🔴 High — needs native audio plugins |
| Audio capture/playback (PCM 16kHz/24kHz) | AudioManager.swift | 🔴 High — platform-specific |
| Echo cancellation | AVAudioSession config | 🟡 Medium |
| Transcription display | StreamSessionView.swift | 🟢 Low |
| OpenClaw REST bridge for tool calling | OpenClawBridge.swift | 🟡 Medium |

### 3.10 Vision/Camera (from VisionClaw)

| Feature | Source | Complexity |
|---------|--------|-----------|
| Camera capture → JPEG frames | IPhoneCameraManager.swift | 🟡 Medium |
| Video → Gemini at 1fps | GeminiSessionViewModel.swift | 🟡 Medium |
| WebRTC live streaming | WebRTCClient.swift + signaling | 🔴 High |

### 3.11 Gamification (from onemind-web — types only, no UI)

XP, levels, ranks, achievements, skill tree. Types and constants defined but zero UI built. Flutter could be the first to actually implement this.

### 3.12 Office Visualization (from onemind-web)

2D SVG floor plan + 3D R3F scene. This is very web-specific (React Three Fiber). Could be reimagined as a simpler 2D agent map in Flutter.

---

## NOT PORTABLE (iOS-only / Web-only)

| Feature | Why Not Portable |
|---------|------------------|
| Live Activities / Dynamic Island | iOS-only (ActivityKit) — no Flutter equivalent |
| Apple FoundationModels (on-device LLM) | iOS 26+ only |
| Meta DAT SDK (Ray-Ban glasses) | Native SDK only |
| App Group (widget data sharing) | iOS-only |
| React Three Fiber 3D scene | Web-only (Three.js) |
| LangGraph server-side execution | Server-side — stays in onemind-web |
| 55 SQLite tables + server runtime | Server-side — stays in onemind-web |
| Browser sandbox (Playwright/Docker) | Server-only |
| 60+ server-side agent tools | Server-only — accessed via gateway |

---

## UNIQUE PATTERNS TO PRESERVE (Don't Lose These)

| Pattern | Source | What It Does | Port To |
|---------|--------|-------------|---------|
| **Triple dedup** | ChatService.swift:420-480 | sequence + toolCallId + content hash prevents duplicates | gateway_client.dart |
| **Verbose tool summary regex** | ChatService.swift:650-700 | Fallback for models without structured events | services/tool_parser.dart |
| **Background instruction** (`deliver: false`) | ChatService.swift:280-310 | Silent system updates without chat bubbles | gateway_client.dart |
| **Tool intent classification** (30+ patterns) | ChatViewModel.swift | Regex → friendly tool name + category + icon | models/tool_category.dart |
| **LOCATION.md template** | ChatViewModel.swift | Formatted location with freshness metadata | services/location_service.dart |
| **BotIdentity ↔ markdown parser** | BotIdentity.swift | `- **Key:** Value` format for IDENTITY.md | models/bot_identity.dart |
| **UserProfile ↔ markdown parser** | UserProfile.swift | Same format for USER.md | models/user_profile.dart |
| **Collect-mode coalescing** (1500ms) | session-run-manager.ts | Merges rapid messages into single LLM turn | N/A (server-side) |
| **setIfChanged fingerprint** | set-if-changed.ts | Skip re-renders when data unchanged | Riverpod equivalent |
| **hmrSingleton** | shared-utils.ts | Survive hot reload | N/A (Flutter has hot reload built-in) |
| **Prompt budget** | prompt-budget.ts | Token-aware prompt assembly | N/A (server-side) |
| **Memory temporal decay** | temporal-decay.ts | Older memories decay in relevance | N/A (server-side) |
| **Extension hook lifecycle** | extensions.ts | 30+ lifecycle hooks | N/A (server-side) |
| **Ed25519 device pairing** | openclaw.ts | Keypair-based device identity | Consider for Flutter (device_info_plus) |

---

## EXTRACTION ROADMAP

### Phase 1: Upgrade Core Client (Week 1)
- [ ] Upgrade `gateway_client.dart` → full protocol v3 (from ChatService.swift)
  - Connect handshake with challenge response
  - Device identity (stable UUID via device_info_plus)
  - `chat.send` with idempotencyKey + deliver flag
  - Full event routing (assistant/thinking/tool/lifecycle streams)
  - History polling with triple dedup
  - Auto-reconnect with exponential backoff
- [ ] Upgrade `rest_client.dart` → add VisionClaw headers (session key, channel)
- [ ] Create all data models (Message, Agent, Session, BotIdentity, UserProfile, etc.)
- [ ] Create Riverpod providers for connection state, agents, sessions

### Phase 2: Real Chat UI (Week 2)
- [ ] Chat thread view with streaming message display
- [ ] Markdown renderer (headings, code blocks, lists, quotes, inline)
- [ ] Code block with syntax label + copy
- [ ] Thinking shimmer + agent activity timeline
- [ ] Tool call display (step-by-step)
- [ ] Chat input bar with send button
- [ ] Session list (sidebar/drawer)
- [ ] New session creation (agent picker)
- [ ] Message pagination + scroll-to-bottom

### Phase 3: Agent Management (Week 3)
- [ ] Agent roster grid/list
- [ ] Agent detail view (status, model, domain, skills)
- [ ] Agent chat sessions list
- [ ] Session CRUD (create, rename, pin, delete)
- [ ] Agent activity step tracker with tool categorization (30+ patterns)

### Phase 4: Settings & Identity (Week 3-4)
- [ ] Connection settings (URL, token, session key)
- [ ] First-launch setup flow
- [ ] Agent identity viewer (parsed from IDENTITY.md)
- [ ] User profile viewer (parsed from USER.md)
- [ ] BotIdentity/UserProfile markdown ↔ model parsers

### Phase 5: Command Center (Week 4)
- [ ] Gateway connection status + health
- [ ] Model list (via gateway API)
- [ ] Schedule/cron viewer + manual trigger
- [ ] Usage/cost analytics view
- [ ] Run history viewer
- [ ] System logs viewer

### Phase 6: Codex Tab (Week 5) — THE DIFFERENTIATOR
- [ ] 3+1 domain health matrix
- [ ] Domain browser (00-99 numbered domains)
- [ ] Vault search (via memory.search RPC)
- [ ] Document viewer (markdown rendering)
- [ ] Life scores display (PI/WI/LI)
- [ ] Review history

### Phase 7: Advanced Features (Week 6+)
- [ ] Multi-agent chatrooms
- [ ] Memory viewer/editor
- [ ] Task board (Kanban)
- [ ] Knowledge base browser
- [ ] Skill manager + ClawhHub
- [ ] MCP server management
- [ ] Connector management
- [ ] Extension management
- [ ] Webhook management
- [ ] Secret vault
- [ ] Wallet management

### Phase 8: Voice & Vision (Week 8+)
- [ ] Voice chat (Gemini Live WebSocket + native audio)
- [ ] Camera capture → AI vision
- [ ] WebRTC streaming
- [ ] Proactive notifications via OpenClaw events

### Phase 9: Gamification (Week 10+)
- [ ] XP system + levels + ranks
- [ ] Achievement tracking
- [ ] Skill tree visualization
- [ ] This would be the FIRST implementation — web has types only

---

## FILE COUNT COMPARISON

| Category | onemind-web | onemind-mobile | onemind-app (current) | onemind-app (target) |
|----------|-------------|----------------|----------------------|---------------------|
| Pages/Screens | 34 | 15 | 5 stubs | ~25 |
| Components | 200 | ~30 | 0 | ~80 |
| API clients | 185 routes | 3 clients | 2 basic | 1 gateway + models |
| State stores | 5 (Zustand) | 3 (ViewModels) | 0 | ~10 (Riverpod) |
| Data models | 55 tables | 15 models | 0 | ~25 |
| Utilities | 50+ | 10 | 0 | ~15 |
| Theme | 300+ CSS vars | 2 color sets | 2 files | 2 files (done) |
| **Total Dart files** | — | — | **13** | **~180** |

---

## WHAT ONEMIND-WEB KEEPS (Server-Side Only)

These features STAY in onemind-web because they require a server runtime:

- Full SQLite database (55 tables)
- LangGraph agent execution pipeline
- 60+ server-side agent tools (shell, files, browser, etc.)
- Extension runtime with 30+ lifecycle hooks
- Connector server processes (Discord bot, Slack bot, etc.)
- Memory consolidation + embedding
- Knowledge/RAG pipeline
- Browser sandbox (Playwright + Docker)
- Protocol/workflow DAG engine
- Mission planner + controller
- Cost calculation engine
- Prompt builder + token budgeting
- CLI (onemind command)
- Background workers (scheduler, heartbeat, daemon)

The Flutter app accesses ALL of this through the OpenClaw gateway API — it doesn't need to run any of it locally.

---

## RECOMMENDATION

1. **Keep onemind-web.** It's the server + full web control plane. The Flutter app is its mobile companion.
2. **Keep onemind-mobile as reference.** Don't delete until all patterns are ported to Flutter.
3. **Remove onemind-mobile from active workspace** once Phase 2 (chat UI) is working.
4. **The Flutter app is simpler by design** — it's a gateway client, not a full server. ~180 files vs 1,391.
5. **Start with Phase 1** (protocol upgrade) — everything else depends on a solid gateway client.


---

## Related

- [[00-24 SO (Sovereign Operations)/05 Development (Code & Testing) (The Forge)/05 Development|05 Development]]
