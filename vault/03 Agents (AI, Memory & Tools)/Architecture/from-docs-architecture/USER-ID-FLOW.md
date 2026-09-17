# User ID Flow Implementation

## Overview

This implementation adds user_id flow through the UI → API → Agent pipeline for proper per-user memory and session management.

## Architecture

```
UI (localStorage)
  ↓ getUserId()
Form Data (user_id field)
  ↓ POST /api/agents/{id}/runs
API Endpoint (extract user_id)
  ↓ agent.user_id = user_id
Agent (uses user_id for memory/sessions)
  ↓
PostgreSQL (per-user data)
```

## Implementation Details

### 1. UI Layer (`ui/`)

#### State Management (`src/store.ts`)
- Added `userId: string` state
- Added `setUserId: (userId: string) => void` action
- Persisted to localStorage via Zustand middleware
- Default value: `'zeus'`

#### Utilities (`src/lib/utils.ts`)
```typescript
export const getUserId = (): string => {
  if (typeof window === 'undefined') return 'zeus' // SSR fallback
  return localStorage.getItem('onemind_user_id') || 'zeus'
}

export const setUserId = (userId: string): void => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('onemind_user_id', userId)
  }
}
```

#### Stream Handler (`src/hooks/useAIStreamHandler.tsx`)
- Reads `userId` from Zustand store
- Appends to FormData: `formData.append('user_id', userId || 'zeus')`
- Sends with every agent/team run request

#### Settings UI (`src/components/settings/UserSection.tsx`)
- New component for user ID configuration
- Input field with real-time validation
- Save button with confirmation feedback
- Integrated into Settings page

### 2. API Layer (`agents/legacy/main.py`)

#### Endpoint: `POST /api/agents/{agent_id}/runs`

```python
form = await request.form()
user_id = form.get("user_id") or "zeus"  # Extract with default

agent = get_or_create_agent(model_id=model_id, agent_id=agent_id)
agent.user_id = user_id  # Set user_id on agent
```

**Edge Cases Handled:**
- Missing `user_id` field → defaults to `"zeus"`
- Empty string `""` → defaults to `"zeus"`
- `None` value → defaults to `"zeus"`

### 3. Agent Layer

The Agno Agent class already supports `user_id` parameter:
- Used for session tracking
- Used for memory storage (per-user)
- Passed to database operations

## Testing

### Manual Test Results
All test cases pass:

**Frontend:**
- ✅ No user_id in localStorage → returns "zeus"
- ✅ User 'alice' in localStorage → returns "alice"
- ✅ Empty string in localStorage → returns "zeus"

**Backend:**
- ✅ Form with user_id='alice' → extracts "alice"
- ✅ Form without user_id key → defaults to "zeus"
- ✅ Form with empty user_id → defaults to "zeus"
- ✅ Form with None user_id → defaults to "zeus"

**End-to-End:**
- ✅ UI → API → Agent flow maintains user_id correctly

## Usage

### Setting User ID

1. Navigate to Settings page
2. Find "User Identity" section
3. Enter desired user ID (e.g., "alice", "bob", "zeus")
4. Click Save or press Enter
5. User ID persisted to localStorage and Zustand store

### Verification

Check browser localStorage:
```javascript
localStorage.getItem('onemind_user_id')
```

Check agent logs:
```
INFO:     Agent run for user: alice
```

Check database:
```sql
SELECT session_id, user_id FROM legacy_sessions WHERE user_id = 'alice';
```

## Benefits

1. **Per-User Memory**: Each user gets isolated memory storage
2. **Session Management**: Sessions tracked per-user
3. **Multi-User Support**: Single instance can serve multiple users
4. **Privacy**: User data remains separate
5. **Debugging**: Easy to trace requests by user_id in logs

## Files Modified

- `ui/src/lib/utils.ts` - Added getUserId() helpers
- `ui/src/store.ts` - Added userId state
- `ui/src/hooks/useAIStreamHandler.tsx` - Sends user_id in requests
- `ui/src/components/settings/UserSection.tsx` - NEW: User ID setting UI
- `ui/src/app/settings/page.tsx` - Added UserSection
- `agents/legacy/main.py` - Extracts and uses user_id

## Future Enhancements

- [ ] User authentication integration
- [ ] User profile management
- [ ] Multi-user admin interface
- [ ] User activity dashboard
- [ ] User-specific agent configurations
