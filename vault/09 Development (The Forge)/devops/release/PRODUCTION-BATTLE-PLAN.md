# OneMind OS Production Battle Plan

> **Goal:** Remove all mock data, fix all broken features, and bring the platform to life
> **Created:** 2026-01-26
> **Status:** Planning

---

## Executive Summary

The OneMind OS platform has extensive backend capabilities but the frontend is riddled with:
- **14+ providers with mock data fallbacks**
- **12+ screens with hardcoded fake data**
- **18 broken features (empty handlers, TODOs)**
- **Missing CRUD operations** for workflows and team persistence

This document outlines a systematic approach to fix everything using parallel agent teams.

---

## Part 1: Current State Assessment

### Backend Status (Core APIs)

| Feature | CRUD Status | Notes |
|---------|-------------|-------|
| **Agents** | ✅ Complete | Full create/read/update/delete, DB persistence |
| **Teams** | ⚠️ Partial | CRUD works but **in-memory only** - lost on restart |
| **Workflows** | ⚠️ Partial | Can list/run, **cannot create/update/delete** |
| **Sessions** | ✅ Complete | Full CRUD |
| **Memory** | ✅ Complete | Full CRUD with PostgreSQL |
| **Knowledge** | ✅ Complete | Upload, search, RAG |
| **Approvals** | ✅ Complete | Full CRUD with HITL |

### Frontend Status (Flutter UI)

| Screen | API Connected | Mock Data | Broken Features |
|--------|---------------|-----------|-----------------|
| **Chat** | ✅ Yes | None | Edit/Regenerate TODOs |
| **Agents** | ✅ Yes | None | Edit goes to chat (wrong route) |
| **Teams** | ⚠️ Partial | Fallback | Create works, no DB persist |
| **Workflows** | ❌ No | Heavy | Create button fake |
| **Approvals** | ⚠️ Partial | Fallback | May work if API returns data |
| **Memory** | ✅ Yes | None | Working |
| **Knowledge** | ✅ Yes | None | Working |
| **Sessions** | ✅ Yes | None | Working |
| **Models** | ❌ No | Heavy | All stats hardcoded |
| **Presence** | ❌ No | Heavy | All locations/devices fake |
| **Extensions** | ❌ No | Heavy | 100% mock data |
| **Consciousness** | ❌ No | Heavy | All data mocked |
| **Intel** | ❌ No | Heavy | Dashboard fake |
| **Skills** | ⚠️ Partial | Fallback | Depends on backend |
| **Activity** | ⚠️ Partial | Fallback | Depends on backend |
| **Events** | ❌ No | Heavy | All mocked |
| **Context** | ❌ No | Heavy | Weather/location fake |
| **Sight** | ⚠️ Partial | Some | Refresh buttons empty |
| **Pillars** | ❌ No | Heavy | All mocked |
| **Multimedia** | ⚠️ Partial | Fallback | Video placeholder |

---

## Part 2: What Needs to Be Fixed

### Priority 1: Backend Gaps (MUST FIX FIRST)

| Task | Effort | Files |
|------|--------|-------|
| **Team DB Persistence** | Medium | `core/teams/persistence.py`, `init-db.sql` |
| **Workflow Create Endpoint** | Medium | `core/orchestrator/main.py` |
| **Workflow Update Endpoint** | Small | `core/orchestrator/main.py` |
| **Workflow Delete Endpoint** | Small | `core/orchestrator/main.py` |

### Priority 2: Frontend - Remove Mock Data (CRITICAL)

| File | Mock Data Location | Lines | Fix |
|------|-------------------|-------|-----|
| `extensions_screen.dart` | `_mockCaptures`, `_mockRules` | 796-854 | Wire to `/api/extension/*` |
| `extension_settings_screen.dart` | `_ExtensionStats.mock()`, `_CaptureRule.mockRules()` | 1341-1495 | Wire to `/api/extension/*` |
| `models_screen.dart` | `_providers` list, `_LiteLLM Status` | 24-80 | Wire to `/api/models/status` |
| `presence_screen.dart` | `_getMockLocations()`, `_buildMockDevicesList()` | 360-931 | Wire to `/api/presence/*` |
| `consciousness_provider.dart` | `mockConsciousnessStats`, etc. | 141-150 | Wire to `/api/consciousness/*` |
| `workflows_provider.dart` | `_mockWorkflows` | 283-350 | Remove fallback, fix create |
| `activity_provider.dart` | `_mockActivity` | 549-650 | Remove fallback |
| `approvals_provider.dart` | `_mockApprovals` | 369-420 | Remove fallback |
| `skills_provider.dart` | `_mockSkills` | 330-400 | Remove fallback |
| `intel_provider.dart` | `_getMockDashboardState()` | 50+ | Remove, return error state |
| `events_provider.dart` | `mockEvents` | 73-78 | Wire to `/api/events/*` |
| `context_provider.dart` | All of it | 30-75 | Wire to `/api/context/*` |
| `pillars_provider.dart` | `mockPillars`, `mockPillarState` | 107-139 | Wire to `/api/pillars/*` |
| `mcp_provider.dart` | `_mockServers` | 306-350 | Should be OK (API exists) |
| `notifications_provider.dart` | `mockNotifications` | 97-112 | Wire to `/api/notifications/*` |
| `multimedia_provider.dart` | `mockMediaItems` | 191-206 | Wire to `/api/multimedia/*` |

### Priority 3: Fix Broken Features

| Feature | File | Issue | Fix |
|---------|------|-------|-----|
| **Agent Edit** | `agents_screen.dart` | Goes to chat, not edit dialog | Change route/show dialog |
| **Workflow Create** | `workflows_screen.dart` L1387-1397 | Fake snackbar, no actual create | Wire to POST `/api/workflows` |
| **Chat Edit** | `chat_screen.dart` L364 | TODO comment | Implement edit |
| **Chat Regenerate** | `chat_screen.dart` L359 | TODO comment | Implement regenerate |
| **File Attachments** | `chat_screen.dart` L150 | TODO comment | Wire upload |
| **Extension Sync** | `extensions_screen.dart` L43 | Empty handler | Wire to sync endpoint |
| **Extension Edit Rule** | `extensions_screen.dart` L390 | Empty handler | Wire to PUT |
| **Extension Delete Rule** | `extensions_screen.dart` L403 | Empty handler | Wire to DELETE |
| **Sight Refresh** | `sight_*_screen.dart` | Empty handlers | Wire refresh |
| **Share Message** | `message_actions.dart` L129 | TODO | Implement share |

---

## Part 3: Agent Team Deployment Strategy

### Recommended Team Structure

Deploy **4 parallel agent teams** for maximum velocity:

| Team | Focus | Agents | Est. Time |
|------|-------|--------|-----------|
| **Backend Team** | Fix backend gaps (team persist, workflow CRUD) | 2 | 4-6 hours |
| **Mock Removal Team** | Remove all mock data from providers | 3 | 6-8 hours |
| **Feature Fix Team** | Fix broken buttons and features | 2 | 4-6 hours |
| **Testing Team** | Validate all fixes work end-to-end | 1 | 2-3 hours |

### Agent Prompts

#### Backend Team Agent 1 - Team Persistence
```
You are fixing team persistence for OneMind OS.

TASK: Add database persistence for teams (currently in-memory only)

1. Create `core/teams/persistence.py`:
   - TeamRecord dataclass
   - save_team(), get_team(), list_teams(), delete_team()
   - Use asyncpg like agent/persistence.py

2. Update `infra/configs/postgres/init-db.sql`:
   - Add `teams` table with columns: id, name, description, instructions, members (JSONB), routing_mode, created_at, updated_at

3. Update `core/orchestrator/main.py`:
   - Replace _team_configs dict with DB calls
   - Load persisted teams on startup

4. Test: Create team via API, restart container, verify team still exists
```

#### Backend Team Agent 2 - Workflow CRUD
```
You are completing workflow CRUD for OneMind OS.

TASK: Add create/update/delete endpoints for workflows

Current state:
- GET /api/workflows works (lists built-in + custom)
- POST /api/workflows/{id}/run works
- Database table `custom_workflows` exists but unused

1. In `core/orchestrator/main.py`, add:
   - POST /api/workflows - Create custom workflow (use custom_workflows table)
   - PUT /api/workflows/{id} - Update workflow
   - DELETE /api/workflows/{id} - Delete workflow (only custom, not built-in)

2. Model for workflow creation:
   - name, description, steps (List[Dict]), trigger (optional)

3. Test all endpoints with curl
```

#### Mock Removal Team Agent 1 - High Priority Providers
```
You are removing mock data from OneMind OS Flutter providers.

TASK: Remove mock data from these HIGH PRIORITY providers:

1. `ui/lib/features/workflows/providers/workflows_provider.dart`:
   - Remove `_mockWorkflows` (L283-350)
   - Don't fallback to mock on error - set error state instead
   - Add createWorkflow(), updateWorkflow(), deleteWorkflow() methods
   - Wire to POST/PUT/DELETE /api/workflows

2. `ui/lib/features/consciousness/providers/consciousness_provider.dart`:
   - Remove mock* variables (L141-150)
   - Fetch real data from /api/consciousness/* on init
   - Handle loading/error states properly

3. `ui/lib/features/pillars/providers/pillars_provider.dart`:
   - Remove mockPillars, mockPillarState, _generateMockHistory()
   - Wire to /api/pillars/* or remove if no backend

Rules:
- If backend endpoint doesn't exist, show "Feature not available" not mock data
- Always handle loading and error states
- Test by checking the screen in browser
```

#### Mock Removal Team Agent 2 - Extension & Models
```
You are removing mock data from OneMind OS Flutter screens.

TASK: Remove mock data from Extensions and Models screens:

1. `ui/lib/features/extensions/screens/extensions_screen.dart`:
   - Remove `_mockCaptures` (L796-824)
   - Remove `_mockRules` (L826-854)
   - Wire to GET /api/extension/stats, GET /api/extension/rules
   - Fix empty button handlers (Sync Now, Edit, Delete)

2. `ui/lib/features/settings/screens/extension_settings_screen.dart`:
   - Remove `_ExtensionStats.mock()` (L1341-1375)
   - Remove `_CaptureRule.mockRules()` (L1450-1495)
   - Wire to same /api/extension/* endpoints

3. `ui/lib/features/models/screens/models_screen.dart`:
   - Remove hardcoded `_providers` list (L35-80)
   - Remove hardcoded LiteLLM status (L24-33)
   - Wire to GET /api/models/status or /api/litellm/status

Rules:
- Show loading spinner while fetching
- Show error state if API fails
- Don't show fake data ever
```

#### Mock Removal Team Agent 3 - Activity, Approvals, Events
```
You are removing mock data from OneMind OS Flutter providers.

TASK: Remove fallback mock data from these providers:

1. `ui/lib/features/activity/providers/activity_provider.dart`:
   - Remove `_mockActivity` (L549-650)
   - Remove fallback at L418, L426
   - Show error state instead of mock data

2. `ui/lib/features/approvals/providers/approvals_provider.dart`:
   - Remove `_mockApprovals` (L369-420)
   - Remove fallback at L281-283
   - Show "No pending approvals" if empty

3. `ui/lib/features/events/providers/events_provider.dart`:
   - Remove `mockEvents`
   - Remove `_loadMockEvents()` (L72-78)
   - Wire to /api/events/* or show empty state

4. `ui/lib/features/skills/providers/skills_provider.dart`:
   - Remove `_mockSkills` (L330-400)
   - Remove fallback at L287-289
   - Wire to /api/skills/* or show empty state

Rules:
- Empty state is OK, fake data is NOT OK
- Show clear error messages
- Handle loading states
```

#### Feature Fix Team Agent 1 - Agent & Workflow Screens
```
You are fixing broken features in OneMind OS Flutter UI.

TASK: Fix broken buttons and navigation

1. `ui/lib/features/agents/screens/agents_screen.dart`:
   - Find where Edit button routes to chat
   - Change to show edit dialog or go to /agents/{id}/edit
   - The edit dialog should use PUT /api/agents/{id}

2. `ui/lib/features/workflows/screens/workflows_screen.dart`:
   - Fix `_createWorkflow()` (L1387-1397)
   - Replace fake snackbar with actual API call
   - Call ref.read(workflowsProvider.notifier).createWorkflow(...)
   - Show success/error from actual response

3. Test:
   - Agent Edit button opens edit dialog (not chat)
   - Workflow Create actually creates and shows in list
```

#### Feature Fix Team Agent 2 - Chat & Extensions
```
You are fixing broken features in OneMind OS Flutter UI.

TASK: Fix chat and extension features

1. `ui/lib/features/chat/screens/chat_screen.dart`:
   - Implement regenerate logic (L359 TODO)
   - Implement edit logic (L364 TODO)
   - Handle file attachments (L150 TODO) - wire to upload

2. `ui/lib/features/extensions/screens/extensions_screen.dart`:
   - Fix Sync Now button (L43) - call /api/extension/sync or similar
   - Fix Edit Rule button (L390) - show edit dialog, call PUT
   - Fix Delete Rule button (L403) - confirm dialog, call DELETE
   - Fix Delete rule in settings (L576)

3. `ui/lib/features/sight/screens/*`:
   - Fix refresh buttons in sight_overview_screen.dart (L19)
   - Fix refresh buttons in sight_pipeline_screen.dart (L19)
   - Call provider refresh methods

Rules:
- Every button must do something
- Show loading states during operations
- Show success/error feedback
```

---

## Part 4: Testing Checklist

### Backend Tests
```bash
# Agents
curl -X POST http://100.102.21.44:8081/api/agents -H "Content-Type: application/json" -d '{"name":"Test Agent","model_id":"claude-sonnet"}'
curl http://100.102.21.44:8081/api/agents
curl -X PUT http://100.102.21.44:8081/api/agents/{id} -H "Content-Type: application/json" -d '{"name":"Updated"}'
curl -X DELETE http://100.102.21.44:8081/api/agents/{id}

# Teams (after persistence fix)
curl -X POST http://100.102.21.44:8081/api/teams -H "Content-Type: application/json" -d '{"name":"Test Team","members":["legacy-mother-ai"]}'
# Restart container, verify team persists
docker restart onemind-agno
curl http://100.102.21.44:8081/api/teams

# Workflows (after CRUD fix)
curl -X POST http://100.102.21.44:8081/api/workflows -H "Content-Type: application/json" -d '{"name":"Test Workflow","steps":[]}'
curl http://100.102.21.44:8081/api/workflows
curl -X PUT http://100.102.21.44:8081/api/workflows/{id} -d '{"name":"Updated"}'
curl -X DELETE http://100.102.21.44:8081/api/workflows/{id}
```

### Frontend Tests (Manual in Browser)

| Test | Steps | Expected |
|------|-------|----------|
| Create Agent | Agents screen → + button → Fill form → Save | Agent appears in list |
| Edit Agent | Agents screen → Agent row → Edit button | Edit dialog opens (NOT chat) |
| Delete Agent | Agents screen → Agent row → Delete → Confirm | Agent removed from list |
| Create Team | Teams screen → + button → Fill form → Save | Team appears in list |
| Create Workflow | Workflows screen → + button → Fill form → Save | Workflow appears in list |
| No Mock Data | Open any screen | No fake data visible |
| Error States | Disconnect backend, reload | Error messages, not mock data |

---

## Part 5: Repository Cleanup

### Folders to Archive/Remove
```bash
# Move old planning docs to archive
mv planning-old archive/
mv docs-old archive/

# Remove any temp files
find . -name "*.pyc" -delete
find . -name "__pycache__" -type d -exec rm -rf {} +
find . -name ".DS_Store" -delete
```

### Consolidate Documentation
```
docs/
├── architecture/      # System architecture
├── features/          # Feature documentation
├── operations/        # Deployment, monitoring
├── api/               # API reference
└── guides/            # How-to guides
```

---

## Part 6: Execution Plan

### Phase 1: Backend Fixes (Day 1)
1. Deploy Backend Team (2 agents)
2. Fix team persistence
3. Fix workflow CRUD
4. Test with curl commands
5. Restart container, verify persistence

### Phase 2: Mock Data Removal (Day 1-2)
1. Deploy Mock Removal Team (3 agents)
2. Remove mock data from all 16 providers
3. Wire to real APIs or show proper empty/error states
4. Test each screen in browser

### Phase 3: Feature Fixes (Day 2)
1. Deploy Feature Fix Team (2 agents)
2. Fix broken buttons and navigation
3. Fix chat edit/regenerate
4. Fix extension CRUD buttons
5. Test all interactions

### Phase 4: Testing & Polish (Day 2-3)
1. Run full test checklist
2. Fix any remaining issues
3. Clean up repository
4. Deploy to production
5. Final verification

---

## Success Criteria

- [ ] **Zero mock data** - No fake data anywhere in the app
- [ ] **Full CRUD** - Can create/edit/delete agents, teams, workflows
- [ ] **All buttons work** - No empty handlers
- [ ] **Proper error states** - Clear messages when things fail
- [ ] **Data persists** - Teams/workflows survive container restart
- [ ] **Clean repo** - No old/duplicate documentation

---

## Quick Start

To begin, copy-paste these prompts to spawn agents:

1. **Open 4 terminals**
2. **Paste prompts from Part 3** into each terminal
3. **Monitor progress** in Linear or via git commits
4. **Run tests from Part 4** after each phase
5. **Deploy** when all tests pass

---

*This document is your single source of truth for bringing OneMind OS to production.*


---

## Related

- [[00-24 SO (Sovereign Operations)/05 Development (Code & Testing) (The Forge)/05 Development|05 Development]]
