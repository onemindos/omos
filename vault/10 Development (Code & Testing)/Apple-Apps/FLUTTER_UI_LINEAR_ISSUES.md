# Flutter UI Modernization - Linear Issues

## Summary
4 Claude Opus agents completed Phase 1 (Chat Modernization) in parallel. Creating Linear issues for remaining Phase 2 work (Screen Implementations).

---

## ISSUE 1: Integrate ChatInput into ChatScreen ⚡ HIGH PRIORITY

**Title:** Integrate enhanced ChatInputBar into chat screen

**Priority:** High (P0)  
**Estimate:** 30 minutes  
**Labels:** flutter, ui, chat, phase-1

**Description:**
All ChatInput widgets are implemented (1,773 lines) but need integration into chat_screen.dart.

**Tasks:**
- [ ] Update imports (chat_input_bar.dart, attached_file.dart)
- [ ] Remove old TextEditingController
- [ ] Replace TextField with ChatInputBar widget
- [ ] Update _sendMessage() to accept files parameter
- [ ] Test all features: Plus menu, dropdowns, quick prompts, slash commands

**Files:**
- `mobile/lib/features/chat/screens/chat_screen.dart` (lines 175-206)

**Reference:**
- `mobile/CHATINPUT-IMPLEMENTATION-STATUS.md` - Complete integration guide

**Acceptance Criteria:**
- ChatInputBar displays with awareness border
- Plus menu opens with camera/gallery/files/location
- Agent/Team/Model selectors work
- Quick prompts show when input empty
- Slash commands overlay appears on "/"
- Token counter shows when >50 chars
- File attachments preview correctly
- Send button works with streaming

---

## ISSUE 2: Knowledge Screen Implementation 📚 HIGH PRIORITY

**Title:** Implement Knowledge screen with RAG and file upload

**Priority:** High (P1)  
**Estimate:** 8 hours  
**Labels:** flutter, ui, screens, knowledge, rag, phase-2, OMOS-81

**Description:**
Implement Knowledge screen with document list, upload, and RAG search functionality.

**Backend Integration:**
- `GET /api/knowledge` - List knowledge documents
- `POST /api/knowledge/upload` - Upload document

**Tasks:**
- [ ] Create knowledge_provider.dart with Riverpod state management
- [ ] Update agno_client.dart with knowledge endpoints (already done ✅)
- [ ] Implement document list with cards (name, type, size, date)
- [ ] Add upload FAB with file picker
- [ ] Add search bar for RAG queries
- [ ] Loading states with shimmer_loading widget
- [ ] Error states with error_widget
- [ ] Empty state with helpful message

**Files:**
- `mobile/lib/features/knowledge/screens/knowledge_screen.dart` (currently scaffolded)
- `mobile/lib/core/providers/knowledge_providers.dart` (new)

**Reference:**
- `ui/src/app/knowledge/page.tsx` - Next.js UI pattern
- `mobile/lib/shared/widgets/` - Use shared widgets

**Acceptance Criteria:**
- Documents load from backend
- Upload works with file picker
- Search queries RAG system
- Loading/error/empty states display
- Mobile-optimized cards layout

---

## ISSUE 3: Settings Screen Implementation ⚙️ MEDIUM PRIORITY

**Title:** Implement Settings screen with 5 categories

**Priority:** Medium (P2)  
**Estimate:** 6 hours  
**Labels:** flutter, ui, screens, settings, phase-2

**Description:**
Implement Settings screen with 5 collapsible categories: Agent, Awareness, Context, Tools, Voice.

**Backend Integration:**
- `GET /api/settings` - Get all settings
- `PUT /api/settings` - Update settings

**Tasks:**
- [ ] Create settings_provider.dart with Riverpod state management
- [ ] Update agno_client.dart with settings endpoints (already done ✅)
- [ ] Create 5 expandable sections (ExpansionPanel)
- [ ] Agent settings: model selection, temperature, max_tokens
- [ ] Awareness settings: heartbeat interval per mode
- [ ] Context settings: max_messages, context_window
- [ ] Tools settings: enable/disable per tool
- [ ] Voice settings: voice provider, speed, pitch
- [ ] Save button with loading state
- [ ] Reset to defaults button

**Files:**
- `mobile/lib/features/settings/screens/settings_screen.dart` (currently scaffolded)
- `mobile/lib/core/providers/settings_providers.dart` (new)

**Reference:**
- `ui/src/app/settings/page.tsx` - Next.js UI pattern
- Use app_tabs.dart for category navigation

**Acceptance Criteria:**
- Settings load from backend
- All 5 categories display correctly
- Changes save to backend
- Reset to defaults works
- Validation for numeric inputs

---

## ISSUE 4: Workflows Screen Implementation 🔄 MEDIUM PRIORITY

**Title:** Implement Workflows screen with trigger functionality

**Priority:** Medium (P2)  
**Estimate:** 6 hours  
**Labels:** flutter, ui, screens, workflows, phase-2

**Description:**
Implement Workflows screen with workflow list and trigger functionality.

**Backend Integration:**
- `GET /api/workflows` - List workflows
- `POST /api/workflows/{name}/run` - Trigger workflow

**Tasks:**
- [ ] Create workflows_provider.dart with Riverpod state management
- [ ] Update agno_client.dart with workflows endpoints (already done ✅)
- [ ] Implement workflow list with cards (name, description, status)
- [ ] Add trigger button per workflow
- [ ] Show execution status (running, completed, failed)
- [ ] Display last run timestamp
- [ ] Add parameters dialog for parameterized workflows
- [ ] Loading/error/empty states

**Files:**
- `mobile/lib/features/workflows/screens/workflows_screen.dart` (currently scaffolded)
- `mobile/lib/core/providers/workflows_providers.dart` (new)

**Reference:**
- `ui/src/app/workflows/page.tsx` - Next.js UI pattern
- Use app_badge.dart for status badges

**Acceptance Criteria:**
- Workflows load from backend
- Trigger button starts workflow
- Status updates in real-time
- Parameters dialog for parameterized workflows
- Error handling for failed workflows

---

## ISSUE 5: Approvals Screen Implementation ✅ MEDIUM PRIORITY

**Title:** Implement Approvals screen with HITL confirmation queue

**Priority:** Medium (P2)  
**Estimate:** 6 hours  
**Labels:** flutter, ui, screens, approvals, hitl, phase-2, OMOS-83

**Description:**
Implement Approvals screen with Human-In-The-Loop confirmation queue.

**Backend Integration:**
- `GET /api/approvals` - Get approval queue
- `POST /api/approvals/{id}/approve` - Approve action
- `POST /api/approvals/{id}/reject` - Reject action (if exists)

**Tasks:**
- [ ] Create approvals_provider.dart with Riverpod state management
- [ ] Update agno_client.dart with approvals endpoints (already done ✅)
- [ ] Implement approval queue list with cards
- [ ] Show action details (tool, args, risk level)
- [ ] Add approve/reject buttons
- [ ] Add reason dialog for rejections
- [ ] Real-time updates via NATS/SSE
- [ ] Loading/error/empty states

**Files:**
- `mobile/lib/features/approvals/screens/approvals_screen.dart` (currently scaffolded)
- `mobile/lib/core/providers/approvals_providers.dart` (new)

**Reference:**
- `ui/src/app/approvals/page.tsx` - Next.js UI pattern
- Use app_dialog.dart for confirmation

**Acceptance Criteria:**
- Approvals load from backend
- Approve button executes action
- Reject button cancels action
- Real-time updates when new approvals arrive
- Risk level badges (low, medium, high, critical)

---

## ISSUE 6: Complete Remaining 12 Screens 📱 LOW PRIORITY

**Title:** Implement remaining 12 screens with backend integration

**Priority:** Low (P3)  
**Estimate:** 40 hours (3-4 hours each)  
**Labels:** flutter, ui, screens, phase-2

**Description:**
Complete the remaining 12 scaffolded screens with backend integration.

**Screens:**
1. Tools (5 tools) - 3h
2. Activity (activity feed) - 3h
3. Presence (location tracking) - 3h
4. Traces (execution traces) - 3h
5. Evals (evaluation results) - 3h
6. MCP (MCP server management) - 5h (OMOS-86)
7. Metrics (system metrics) - 4h
8. Databases (database connections) - 3h
9. Inbox (notifications) - 3h
10. Frame (BLE wearable) - 5h
11. Watch (MQTT biometrics) - 5h
12. Voice Panel (already complete via overlay) - 0h

**Tasks:**
- [ ] Create Riverpod providers for each screen
- [ ] Add missing API endpoints to agno_client.dart
- [ ] Implement list/detail views with cards
- [ ] Use shared widget library throughout
- [ ] Loading/error/empty states for all
- [ ] Mobile-optimized layouts

**Reference:**
- `ui/src/app/` - Next.js UI patterns
- `mobile/lib/features/teams/screens/teams_screen.dart` - Example implementation

**Acceptance Criteria:**
- All 12 screens load data from backend
- Consistent UI using shared widgets
- Mobile-first responsive design
- Error handling throughout

---

## Implementation Timeline

### Sprint 1 (Week 1) - HIGH PRIORITY ⚡
- **Day 1:** Issue 1 - ChatInput integration (30 min)
- **Day 1-2:** Issue 2 - Knowledge screen (8h)
- **Day 3-4:** Issue 3 - Settings screen (6h)

### Sprint 2 (Week 2) - MEDIUM PRIORITY
- **Day 5-6:** Issue 4 - Workflows screen (6h)
- **Day 7-8:** Issue 5 - Approvals screen (6h)

### Sprint 3 (Weeks 3-4) - LOW PRIORITY
- **Days 9-20:** Issue 6 - Remaining 12 screens (40h)

**Total Estimate:** 66.5 hours over 4 weeks

---

## Related Linear Issues

These Flutter UI issues relate to existing backend Linear issues:
- **OMOS-81:** Knowledge/RAG with LanceDb (backend) → Issue 2 (frontend)
- **OMOS-83:** HITL confirmation (backend) → Issue 5 (frontend)
- **OMOS-86:** MCP support (backend) → Issue 6 (MCP screen frontend)

---

## Success Metrics

After all issues complete:
- ✅ 24/24 screens functional with backend
- ✅ Shared widget library in use
- ✅ Mobile-first responsive design
- ✅ Black/white/red color scheme
- ✅ 100% feature parity with Next.js UI
- ✅ Ready for App Store/Play Store submission
