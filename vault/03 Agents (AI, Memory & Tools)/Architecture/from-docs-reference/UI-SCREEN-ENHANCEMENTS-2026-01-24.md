# Flutter UI Screen Enhancements

> **Date:** 2026-01-24
> **Total Lines Added:** ~4,000 lines across 4 screens
> **Deployment:** http://100.102.21.44:3005

---

## Summary

| Screen | Route | Before | After | Status |
|--------|-------|--------|-------|--------|
| **MCP** | `/mcp` | 27 lines (stub) | 1,177 lines | ✅ Complete |
| **Activity** | `/activity` | 214 lines (mock) | 999 lines | ✅ Complete |
| **Workflows** | `/workflows` | 66 lines (mock) | 1,430 lines | ✅ Complete |
| **Memory** | `/memory` | 374 lines | 690 lines | ✅ Enhanced |

**Screen Status Updated:**
- ✅ Complete: 11 → **14** (+3)
- 🔄 Partial: 13 → **10** (-3)
- ⏳ Stub: 1 → **0** (-1)

---

## 1. MCP Screen (`/mcp`)

**File:** `ui/lib/features/mcp/screens/mcp_screen.dart`
**Lines:** 27 → 1,177

### Features Implemented

1. **Search Bar**
   - Filter servers by name
   - Real-time filtering as you type
   - Clear button

2. **Server List**
   - Card-based layout with server details
   - Server name, type (stdio/sse), and description
   - Tool count badge
   - Resource count badge

3. **Status Indicators**
   - Connected: Green badge with checkmark
   - Connecting: Orange badge with sync spinner
   - Disconnected: Grey badge
   - Error: Red badge with error icon

4. **Server Actions**
   - Connect/Disconnect toggle button
   - Enable/Disable switch
   - Delete with confirmation dialog
   - Reconnect action

5. **Add Server Dialog**
   - Name text field
   - Transport dropdown (stdio, sse, streamable-http)
   - Command field (for stdio transport)
   - URL field (for sse/streamable-http)
   - Arguments text field
   - Environment variables field
   - Create button

6. **Server Details Sheet**
   - Bottom sheet with full server info
   - Tools list with descriptions
   - Resources list with URIs
   - Connection history
   - Error messages if any

7. **Stats Card**
   - Total servers count
   - Connected servers count
   - Total tools available
   - Total resources available

8. **State Management**
   - Connected to `mcpProvider` from Riverpod
   - Loading state with spinner
   - Error state with retry button
   - Empty state with "Add Server" prompt

### Provider Integration

```dart
// Uses these providers from mcp_provider.dart:
ref.watch(mcpProvider)                    // Main state
ref.watch(connectedMcpServersProvider)    // Connected servers
ref.watch(mcpToolsCountProvider)          // Tool count
```

### Color Scheme

| Element | Color |
|---------|-------|
| Background | #121212 |
| Cards | #1A1A1A |
| Connected | #4CAF50 (green) |
| Connecting | #FF9800 (orange) |
| Disconnected | #9E9E9E (grey) |
| Error | #F44336 (red) |
| Accent | #FF0000 |

---

## 2. Activity Screen (`/activity`)

**File:** `ui/lib/features/activity/screens/activity_screen.dart`
**Lines:** 214 → 999

### Features Implemented

1. **Filter Chips**
   - All (default)
   - Agent Run
   - Tool Call
   - Memory
   - Workflow
   - Awareness
   - Session
   - Inbox
   - Approval
   - System
   - Integration
   - Error

2. **Unread Toggle**
   - Filter chip to show only unread items
   - Unread count badge

3. **Activity Stats Card**
   - Total activities
   - Unread count
   - Today's count
   - Breakdown by type

4. **Activity List**
   - Type icon with color
   - Title (bold if unread)
   - Description
   - Time ago (e.g., "5m ago", "2h ago")
   - Actor badge
   - Red dot for unread items

5. **Pull-to-Refresh**
   - RefreshIndicator wrapper
   - Calls `loadActivity()`

6. **Mark as Read**
   - Tap item to mark as read
   - "Mark All as Read" button in app bar
   - Badge showing unread count

7. **Activity Detail Sheet**
   - Full-screen bottom sheet
   - Activity ID and timestamp
   - Actor information
   - Session ID, Run ID, Workflow ID (if present)
   - Duration and tokens used (if present)
   - Error message in red box (if present)
   - Metadata as JSON (if present)

8. **State Handling**
   - Loading state with spinner
   - Error state with retry button
   - Empty state with icon

### Provider Integration

```dart
// Uses these providers from activity_provider.dart:
ref.watch(activityProvider)              // Main state
ref.watch(activityUnreadCountProvider)   // Unread count
ref.watch(activityByTypeProvider)        // Count by type
```

### Activity Types & Colors

| Type | Color | Icon |
|------|-------|------|
| Agent Run | Blue | smart_toy |
| Tool Call | Purple | build |
| Memory | Green | psychology |
| Workflow | Orange | account_tree |
| Awareness | Cyan | visibility |
| Session | Teal | forum |
| Inbox | Amber | inbox |
| Approval | Indigo | approval |
| System | Grey | settings |
| Integration | Pink | sync |
| Error | Red | error |

---

## 3. Workflows Screen (`/workflows`)

**File:** `ui/lib/features/workflows/screens/workflows_screen.dart`
**Lines:** 66 → 1,430

### Features Implemented

1. **Filter Chips**
   - All (default)
   - Active (enabled workflows)
   - Scheduled (cron-based)
   - Manual (trigger-based)

2. **Workflow Stats Card**
   - Total workflows
   - Active count
   - Scheduled count
   - Total runs across all workflows

3. **Workflow Cards**
   - Dynamic icon based on workflow type
   - Name and description
   - Status badge (idle/running/completed/failed/paused)
   - Steps count
   - Run count
   - Schedule badge (for scheduled workflows)
   - Last run timestamp
   - Success rate percentage

4. **Status Indicators**
   - Idle: Grey badge
   - Running: Blue badge with animated spinner
   - Completed: Green badge with checkmark
   - Failed: Red badge with X
   - Paused: Orange badge

5. **Workflow Actions**
   - Play button (trigger workflow)
   - Stop button (cancel running workflow)
   - Enable/Disable toggle switch

6. **Running Workflow Banner**
   - Blue highlighted banner at top
   - Shows workflow name
   - Step count with spinner
   - Appears when any workflow is running

7. **Workflow Details Sheet**
   - Full-screen draggable bottom sheet
   - Header with name, status, description
   - Run count in circle badge
   - Info chips (steps, schedule, last run)
   - Timeline-style step list showing:
     - Step type (agent/tool/condition/parallel)
     - Step status
     - Step name
     - Duration
     - Output preview

8. **Create Workflow Dialog**
   - Name text field (required)
   - Description text field
   - Trigger type dropdown (Manual/Scheduled)
   - Cron expression field (for scheduled)
   - Quick preset chips:
     - Daily 7am: `0 7 * * *`
     - Daily 5pm: `0 17 * * *`
     - Hourly: `0 * * * *`
     - Friday 5pm: `0 17 * * 5`
   - Create button (disabled until name entered)

9. **Cron Display**
   - Human-readable cron descriptions
   - Examples: "Daily at 7:00 AM", "Every hour"

### Provider Integration

```dart
// Uses these providers from workflows_provider.dart:
ref.watch(workflowsProvider)            // Main state
ref.watch(runningWorkflowProvider)      // Currently running
ref.watch(scheduledWorkflowsProvider)   // Scheduled workflows
```

### Workflow Icons

| Workflow Type | Icon |
|---------------|------|
| Morning/Briefing | wb_sunny |
| Task/Triage | assignment |
| Focus | center_focus_strong |
| Review/Weekly | rate_review |
| Summary/Daily | summarize |
| Scheduled | schedule |
| Default | account_tree |

---

## 4. Memory Screen (`/memory`)

**File:** `ui/lib/features/memory/screens/memory_screen.dart`
**Lines:** 374 → 690

### Features Implemented

1. **Search Bar**
   - Filter memories by content
   - Real-time filtering
   - Search icon prefix

2. **Category Filter Chips**
   - All (default)
   - Dynamic chips from existing categories
   - Chip count badges

3. **Memory Stats Card**
   - Total memories count
   - Category breakdown
   - Recent memories count

4. **Memory List**
   - Card-based layout
   - Memory content preview
   - Category badge
   - Created date
   - Edit and delete actions

5. **Create Memory Dialog**
   - Content text field (multiline)
   - Category dropdown or text field
   - Metadata JSON field (optional)
   - Create button

6. **Edit Memory Dialog**
   - Pre-filled content
   - Category selection
   - Update button

7. **Delete Confirmation**
   - Confirmation dialog
   - Undo snackbar option

8. **Enhanced UI**
   - Loading state with spinner
   - Error state with retry
   - Empty state with "Add Memory" prompt
   - Pull-to-refresh

### Provider Integration

```dart
// Uses these providers from app_providers.dart:
ref.watch(memoriesProvider)        // Memory list
ref.watch(memoryStatsProvider)     // Stats
ref.watch(memoryTopicsProvider)    // Categories/topics
```

---

## Files Modified

| File | Changes |
|------|---------|
| `ui/lib/features/mcp/screens/mcp_screen.dart` | Complete rewrite (27 → 1,177 lines) |
| `ui/lib/features/activity/screens/activity_screen.dart` | Complete rewrite (214 → 999 lines) |
| `ui/lib/features/workflows/screens/workflows_screen.dart` | Complete rewrite (66 → 1,430 lines) |
| `ui/lib/features/memory/screens/memory_screen.dart` | Enhanced (374 → 690 lines) |
| `ui/lib/features/SCREEN_STATUS.md` | Updated status counts |
| `ui/lib/features/settings/providers/notifications_provider.dart` | Fixed API call syntax |

---

## Bug Fixes During Build

1. **notifications_provider.dart** - Fixed API client method calls
   - Changed `_client.put(path, data: body)` to `_client.put(path, body)`
   - Changed `_client.post(path, data: body)` to `_client.post(path, body)`
   - Removed response.statusCode check (API client returns parsed JSON)

---

## Deployment

```bash
# Build
cd ui && flutter build web --release

# Deploy to server
scp -r -i ~/.ssh/onemind-key.pem ui/build/web/* ubuntu@100.102.21.44:/home/ubuntu/onemind-os/ui/build/web/

# Update version and restart
ssh -i ~/.ssh/onemind-key.pem ubuntu@100.102.21.44 \
  "echo '{\"version\": \"$(date +%s)\"}' > /home/ubuntu/onemind-os/ui/build/web/version.json && \
   docker restart onemind-flutter-web"
```

---

## Architecture Pattern

All screens follow the same pattern:

```dart
class ScreenName extends ConsumerStatefulWidget {
  @override
  ConsumerState<ScreenName> createState() => _ScreenNameState();
}

class _ScreenNameState extends ConsumerState<ScreenName> {
  @override
  Widget build(BuildContext context) {
    final state = ref.watch(screenProvider);

    return BaseScreen(
      title: 'Screen Title',
      actions: [...],
      floatingActionButton: FloatingActionButton(...),
      child: state.when(
        loading: () => _buildLoading(),
        error: (e, s) => _buildError(e),
        data: (data) => _buildContent(data),
      ),
    );
  }
}
```

### Key Patterns Used

1. **ConsumerStatefulWidget** - For stateful screens with Riverpod
2. **ref.watch()** - Reactive state watching
3. **ref.read()** - One-time reads for actions
4. **ref.invalidate()** - Force refresh after mutations
5. **BaseScreen** - Consistent layout wrapper
6. **FilterChip** - Type/status filtering
7. **BottomSheet** - Detail views
8. **AlertDialog** - Create/edit forms

---

## Related Documentation

- Provider files: `ui/lib/features/*/providers/*_provider.dart`
- Screen status: `ui/lib/features/SCREEN_STATUS.md`
- Backend endpoints: `core/orchestrator/main.py`

---

## Next Steps

Remaining partial screens to enhance:

1. **Knowledge** (OMOS-81) - RAG with file upload, URL ingestion
2. **Approvals** (OMOS-83) - HITL confirmation queue
3. **Tools** (OMOS-82) - Built-in tools management
4. **Traces** - Phoenix/Langfuse integration
5. **Metrics** - Usage dashboard
6. **Inbox** - Notifications
7. **Presence** - Real-time presence
8. **Databases** - Database management
9. **Frame/Watch** - Wearables integration
