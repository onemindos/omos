# OneMind OS - UI Design Plan

> **Stack:** Next.js 15 + shadcn/ui + Tailwind CSS + TypeScript
> **Base:** Forked Agno UI (`ui/`)
> **Last Updated:** 2026-01-18

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     ONEMIND UI ARCHITECTURE                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   AGNO UI (Base)              ONEMIND EXTENSIONS                │
│   ├── Chat Interface          ├── AwarenessBar                  │
│   ├── Sessions/History        ├── NotificationsPanel            │
│   ├── Agent Selector          ├── EventStream (NATS SSE)        │
│   ├── Tool Visualization      ├── MemoryTimeline                │
│   ├── Markdown Rendering      ├── MCPPanel                      │
│   └── Multi-modal Support     ├── ToolsPanel                    │
│                               ├── CommandPalette (Cmd+K)        │
│                               └── OneMind Branding              │
│                                                                  │
│   BACKEND: AgentOS API (Port 8080)                              │
│   ├── /agents         - List/manage agents                      │
│   ├── /sessions       - Chat sessions                           │
│   ├── /health         - Health check                            │
│   └── /agents/{id}/runs - Execute agent                         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Component Inventory

### From Agno UI (Already Have)

| Component | Location | Purpose |
|-----------|----------|---------|
| `ChatArea` | `components/chat/ChatArea/` | Main chat interface |
| `ChatInput` | `components/chat/ChatArea/ChatInput/` | Message input |
| `Messages` | `components/chat/ChatArea/Messages/` | Message display |
| `MessageItem` | `components/chat/ChatArea/Messages/` | Single message |
| `AgentThinkingLoader` | `components/chat/ChatArea/Messages/` | Loading state |
| `Multimedia` | `components/chat/ChatArea/Messages/Multimedia/` | Images/Audio/Video |
| `Sidebar` | `components/chat/Sidebar/` | Navigation sidebar |
| `Sessions` | `components/chat/Sidebar/Sessions/` | Session history |
| `EntitySelector` | `components/chat/Sidebar/` | Agent picker |
| `ModeSelector` | `components/chat/Sidebar/` | Mode picker |
| `MarkdownRenderer` | `components/ui/typography/` | Markdown display |
| `AuthPrompt` | `components/` | Authentication |
| `ProtectedRoute` | `components/` | Route protection |

### shadcn/ui Components (Already Have)

| Component | Status |
|-----------|--------|
| `button` | ✅ |
| `dialog` | ✅ |
| `select` | ✅ |
| `skeleton` | ✅ |
| `textarea` | ✅ |
| `tooltip` | ✅ |

### To Migrate from ui_portal_gateway

| Component | Priority | Source | Destination |
|-----------|----------|--------|-------------|
| `AwarenessBar` | HIGH | `components/awareness/` | `components/onemind/` |
| `NotificationsPanel` | HIGH | `components/shared/` | `components/onemind/` |
| `NotificationItem` | HIGH | `components/shared/` | `components/onemind/` |
| `NotificationBell` | HIGH | `components/shared/` | `components/onemind/` |
| `EventStream` | HIGH | `components/ops/` | `components/onemind/` |
| `EventCard` | HIGH | `components/ops/` | `components/onemind/` |
| `MemoryTimeline` | MEDIUM | `components/intel/` | `components/memory/` |
| `MemoryCard` | MEDIUM | `components/intel/` | `components/memory/` |
| `MemoryNode` | MEDIUM | `components/intel/` | `components/memory/` |
| `ToolsPanel` | MEDIUM | `components/tools/` | `components/tools/` |
| `ToolCard` | MEDIUM | `components/tools/` | `components/tools/` |
| `MCPPanel` | MEDIUM | `components/mcp/` | `components/mcp/` |
| `MCPServerCard` | MEDIUM | `components/mcp/` | `components/mcp/` |
| `CommandPalette` | LOW | `components/shared/` | `components/ui/` |
| `OneMindLogo` | LOW | `components/onemind/` | `components/branding/` |

### Hooks to Migrate

| Hook | Purpose | Priority |
|------|---------|----------|
| `useAwareness` | Awareness mode API | HIGH |
| `useNatsEvents` | NATS SSE stream | HIGH |
| `useNotifications` | Notification system | HIGH |
| `useAudioAnalyzer` | Voice waveform | LOW |

### shadcn/ui Components to Add

```bash
# Run these to add missing components
npx shadcn@latest add card
npx shadcn@latest add badge
npx shadcn@latest add accordion
npx shadcn@latest add tabs
npx shadcn@latest add switch
npx shadcn@latest add input
npx shadcn@latest add label
npx shadcn@latest add sheet
npx shadcn@latest add command  # For CommandPalette
```

---

## Page Structure

### Current Agno UI Pages

```
src/app/
└── page.tsx          # Main chat interface
```

### Proposed OneMind Pages

```
src/app/
├── page.tsx          # Main chat (enhanced with AwarenessBar)
├── ops/
│   └── page.tsx      # Operations dashboard (EventStream)
├── intel/
│   └── page.tsx      # Intelligence (MemoryTimeline)
├── tools/
│   └── page.tsx      # Tools management
├── mcp/
│   └── page.tsx      # MCP server management
└── settings/
    └── page.tsx      # Configuration
```

---

## Design System

### Colors (OneMind Theme)

```css
/* Add to globals.css */
:root {
  --onemind-purple: #8B5CF6;
  --onemind-purple-dark: #7C3AED;
  --onemind-red: #EF4444;
  --onemind-black: #0F0F0F;
  --onemind-black-card: #1A1A1A;
  --onemind-gray: #6B7280;
  --onemind-gray-light: #9CA3AF;
  --onemind-gray-dark: #374151;
  --onemind-white: #F9FAFB;
}
```

### Typography

- **Headings:** Inter (default from Agno)
- **Body:** Inter
- **Code:** JetBrains Mono (for tool outputs)

### Component Patterns

1. **Cards:** Use `bg-onemind-black-card` with `border-onemind-gray-dark`
2. **Buttons:** Primary = `onemind-purple`, Destructive = `onemind-red`
3. **Status indicators:** Green/Yellow/Orange/Red awareness levels
4. **Animations:** Framer Motion for transitions (already in Agno UI)

---

## Implementation Phases

### Phase 1: Core Integration (Current)
- [x] Fork Agno UI
- [x] Configure for AgentOS API
- [ ] Add OneMind branding/colors
- [ ] Migrate AwarenessBar
- [ ] Migrate NotificationsPanel

### Phase 2: Dashboard Pages
- [ ] Create /ops page with EventStream
- [ ] Create /intel page with MemoryTimeline
- [ ] Create /tools page
- [ ] Create /mcp page
- [ ] Add navigation sidebar

### Phase 3: Enhanced Features
- [ ] CommandPalette (Cmd+K)
- [ ] Voice integration (LiveKit)
- [ ] PWA support
- [ ] Mobile responsive

### Phase 4: Polish
- [ ] Dark/light theme toggle
- [ ] Animation refinements
- [ ] Performance optimization
- [ ] Accessibility audit

---

## API Integration Points

### AgentOS API (Backend)

```typescript
// src/api/routes.ts - Already configured
const API = {
  agents: '/agents',
  sessions: '/sessions',
  health: '/health',
  agentRun: '/agents/{id}/runs',
}
```

### Custom OneMind APIs (to add)

```typescript
// New endpoints needed
const OneMindAPI = {
  awareness: '/awareness',           // GET/PUT awareness mode
  events: '/events',                 // SSE stream from NATS
  notifications: '/notifications',   // GET/POST notifications
  memory: '/memory',                 // GET memory entries
  tools: '/tools',                   // GET available tools
  toolExecute: '/tools/execute',     // POST execute tool
  mcp: '/mcp/servers',              // GET/POST MCP servers
}
```

---

## File Migration Checklist

### High Priority

```bash
# From ui_portal_gateway to ui
□ components/awareness/AwarenessBar.tsx
□ components/awareness/index.ts
□ components/shared/NotificationsPanel.tsx
□ components/shared/NotificationItem.tsx
□ components/shared/NotificationBell.tsx
□ components/ops/EventStream.tsx
□ components/ops/EventCard.tsx
□ hooks/useAwareness.ts
□ hooks/useNatsEvents.ts
□ hooks/useNotifications.ts
```

### Medium Priority

```bash
□ components/intel/MemoryTimeline.tsx
□ components/intel/MemoryCard.tsx
□ components/intel/MemoryNode.tsx
□ components/tools/ToolsPanel.tsx
□ components/tools/ToolCard.tsx
□ components/mcp/MCPPanel.tsx
□ components/mcp/MCPServerCard.tsx
```

### Low Priority

```bash
□ components/shared/CommandPalette.tsx
□ components/onemind/OneMindLogo.tsx
□ hooks/useAudioAnalyzer.ts
```

---

## Development Workflow

### Local Development

```bash
# Terminal 1: Start AgentOS backend
cd /path/to/onemind-os
docker compose up agno -d

# Terminal 2: Start UI with hot reload
cd ui
pnpm install
pnpm watch
```

### Build for Production

```bash
cd ui
pnpm build
# Static files output to: out/
```

### Docker Deployment

The Agno container serves the UI from `/ui` path after building.

---

## Notes

- **State Management:** Zustand (already in Agno UI)
- **Styling:** Tailwind CSS with shadcn/ui components
- **Icons:** Lucide React
- **Animations:** Framer Motion
- **Forms:** React Hook Form (add if needed)
- **Data Fetching:** Native fetch with SWR pattern (consider adding SWR or TanStack Query)

---

## References

- [Agno UI Repository](https://github.com/agno-agi/agent-ui)
- [shadcn/ui Documentation](https://ui.shadcn.com)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [Next.js 15 Docs](https://nextjs.org/docs)
