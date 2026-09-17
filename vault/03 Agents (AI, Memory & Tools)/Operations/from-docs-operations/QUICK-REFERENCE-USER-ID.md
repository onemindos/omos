# User ID Flow - Quick Reference

## What Was Implemented

✅ **Full user_id flow from UI → API → Agent**

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         FRONTEND (UI)                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Settings Page (/settings)                                       │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  👤 User Identity                                          │ │
│  │     ┌──────────────────────┐  ┌──────┐                    │ │
│  │     │ zeus               ▼ │  │ Save │                    │ │
│  │     └──────────────────────┘  └──────┘                    │ │
│  └────────────────────────────────────────────────────────────┘ │
│                          ↓                                        │
│  Zustand Store (userId state) ← → localStorage                  │
│                          ↓                                        │
│  useAIStreamHandler (appends user_id to FormData)               │
│                                                                   │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            │ POST /api/agents/{id}/runs
                            │ FormData: { message, user_id, ... }
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│                         BACKEND (API)                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  FastAPI Endpoint: run_agent_with_model()                       │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  1. Extract: user_id = form.get("user_id") or "zeus"      │ │
│  │  2. Log: f"Agent run for user: {user_id}"                 │ │
│  │  3. Set: agent.user_id = user_id                          │ │
│  └────────────────────────────────────────────────────────────┘ │
│                          ↓                                        │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│                         AGENT (Agno)                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Agent.user_id = "alice"                                         │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  → Sessions stored per user                                │ │
│  │  → Memories stored per user                                │ │
│  │  → Database operations use user_id                         │ │
│  └────────────────────────────────────────────────────────────┘ │
│                          ↓                                        │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│                    DATABASE (PostgreSQL)                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  legacy_sessions (user_id: alice, session_id, ...)              │
│  legacy_memories (user_id: alice, memory, ...)                  │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

## Key Features

### 1. UI Configuration
- Settings page has new "User Identity" section
- Input field for custom user ID
- Save button with visual feedback
- Persists to localStorage
- Default: "zeus"

### 2. State Management
- Zustand store: `userId` state
- Synced to localStorage automatically
- Available globally via `useStore()`

### 3. Request Flow
- Every chat message includes `user_id` in FormData
- Extracted in `useAIStreamHandler` hook
- Sent to backend API

### 4. Backend Processing
- FastAPI endpoint extracts `user_id`
- Sets on agent instance
- Logs for debugging
- Handles edge cases (missing, empty, None)

### 5. Agent Usage
- Agent uses `user_id` for all operations
- Sessions isolated per user
- Memories isolated per user
- Database queries filtered by user

## Usage Example

### Setting User ID
1. Go to Settings (`/settings`)
2. Find "User Identity" section
3. Enter "alice"
4. Click Save
5. user_id stored in localStorage

### Chatting
1. Type message: "Hello Legacy"
2. Click Send
3. Request sent with `user_id: "alice"`
4. Agent processes with user context
5. Response stored in alice's session
6. Memory saved to alice's memory store

### Multi-User Scenario
```
User: alice
  → Sessions: alice_session_1, alice_session_2
  → Memories: "Alice prefers coffee", "Alice's birthday: Jan 1"

User: bob
  → Sessions: bob_session_1, bob_session_2
  → Memories: "Bob likes tea", "Bob works at NASA"
```

## Code Locations

### Frontend
- `ui/src/lib/utils.ts` - getUserId() helper
- `ui/src/store.ts` - userId state
- `ui/src/hooks/useAIStreamHandler.tsx` - FormData append
- `ui/src/components/settings/UserSection.tsx` - Settings UI
- `ui/src/app/settings/page.tsx` - Settings page

### Backend
- `agents/legacy/main.py` - user_id extraction

### Documentation
- `docs/architecture/user-id-flow.md` - Technical docs
- `docs/architecture/user-id-flow-ui.md` - UI guide
- `IMPLEMENTATION-SUMMARY.md` - Complete summary

## Testing

All edge cases tested and passing:
- ✅ No user_id → defaults to "zeus"
- ✅ Empty user_id → defaults to "zeus"
- ✅ Custom user_id → uses custom value
- ✅ localStorage sync → persists correctly
- ✅ End-to-end flow → works as expected

## Benefits

1. **Privacy**: User data isolated
2. **Multi-User**: Multiple users can use same instance
3. **Debugging**: Easy to trace by user_id
4. **Memory**: Per-user memory storage
5. **Sessions**: Per-user session tracking
6. **Scalability**: Ready for user authentication

## Next Steps

When deploying to production:
1. Test in live environment
2. Verify localStorage persistence
3. Check Network tab for user_id in requests
4. Verify agent logs show user_id
5. Query database for per-user data

## Support

See comprehensive documentation:
- Technical: `docs/architecture/user-id-flow.md`
- UI Guide: `docs/architecture/user-id-flow-ui.md`
- Summary: `IMPLEMENTATION-SUMMARY.md`
