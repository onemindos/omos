---
id: ui-02-20260111053925
title: "Telegram Bot Integration"
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
ai_summary: "Telegram Bot Integration"
---

# Telegram Bot Integration

> **Bot:** @onemind_approval_bot
> **Status:** Active
> **Last Updated:** 2024-12-24

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                      LETTA AGENTS                               │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │  Legacy-AI  │ │HP-Conductor │ │LE-Conductor │  ...          │
│  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘               │
│         │               │               │                       │
│         └───────────────┼───────────────┘                       │
│                         │                                       │
│              ┌──────────▼──────────┐                           │
│              │   TELEGRAM TOOLS    │                           │
│              │  (Letta Custom)     │                           │
│              └──────────┬──────────┘                           │
└─────────────────────────┼───────────────────────────────────────┘
                          │
               ┌──────────▼──────────┐
               │   TELEGRAM BOT API  │
               │  @onemind_approval  │
               └──────────┬──────────┘
                          │
               ┌──────────▼──────────┐
               │     ZEUS PHONE      │
               │   (Notifications)   │
               └─────────────────────┘
```

---

## Credentials

| Item | Value |
|------|-------|
| Bot Username | @onemind_approval_bot |
| Bot Token | `8437988392:AAHsEyJFsOr2Qdqyp1B17c4nQgagLWDI8LM` |
| Chat ID (Zeus) | `8297858681` |
| Telegram Username | @one_mind_OS |

### API Endpoints

```
Base URL: https://api.telegram.org/bot{TOKEN}

Send Message:    POST /sendMessage
Get Updates:     GET  /getUpdates
Answer Callback: POST /answerCallbackQuery
Edit Message:    POST /editMessageText
```

---

## Tools Built (COMPLETED)

### 1. telegram_notify_approval
**Tool ID:** `tool-8fbb64ee-3c85-400d-9836-8804d64fb9f8`

**Purpose:** Request human approval via Telegram (HITL)

**Parameters:**
| Param | Type | Required | Description |
|-------|------|----------|-------------|
| action | string | Yes | Short action name |
| description | string | Yes | Detailed description |
| context | string | No | Additional context |
| urgency | string | No | "low", "normal", "high", "critical" |

**Returns:** `{"request_id": "abc123", "message_id": 10, "sent": true, "status": "pending"}`

**Behavior:**
- Sends notification with Approve/Deny buttons
- Returns immediately (async, no timeout)
- Use `telegram_check_approval` to check response

---

### 2. telegram_check_approval
**Tool ID:** `tool-8f8074ef-4f94-4a92-b330-9cfa6976d978`

**Purpose:** Check if Zeus responded to approval request

**Parameters:**
| Param | Type | Required | Description |
|-------|------|----------|-------------|
| request_id | string | No | Specific request to check |

**Returns:**
- If found: `{"request_id": "abc123", "approved": true, "decision": "approve"}`
- If pending: `{"request_id": "abc123", "status": "pending"}`

---

### 3. telegram_send_alert
**Tool ID:** `tool-6ff1ffe8-975b-44a5-99da-0b7f011192d4`

**Purpose:** Send alerts/notifications (no approval needed)

**Parameters:**
| Param | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | Yes | Alert title |
| message | string | Yes | Alert details |
| category | string | No | "health", "finance", "security", "system", "task", "general" |
| priority | string | No | "low", "normal", "high", "critical" |

**Message Format:**
```
⚠️ *API Key Expiring Soon*

🔐 _SECURITY_

Your Todoist API token expires in 7 days...
```

---

### 4. telegram_send_insight
**Tool ID:** `tool-aac31317-99bd-4252-8495-8969cd50a955`

**Purpose:** Send proactive observations/insights

**Parameters:**
| Param | Type | Required | Description |
|-------|------|----------|-------------|
| observation | string | Yes | What was noticed |
| data_source | string | Yes | Where insight came from |
| recommendation | string | No | Suggested action |
| pillar | string | No | "HP", "LE", "GE", "IT", "CMD" |

**Message Format:**
```
💡 *INSIGHT*

💚 *HP*

*Observation:*
You've completed 3 workouts this week but sleep score dropped 15%.

*Source:* Pattern analysis

*Recommendation:*
Consider earlier bedtime tonight.
```

---

### 5. telegram_send_daily_summary
**Tool ID:** `tool-13fc5fa9-4641-4735-b29a-531ff824a877`

**Purpose:** End-of-day summary

**Parameters:**
| Param | Type | Required | Description |
|-------|------|----------|-------------|
| completed_tasks | int | Yes | Tasks completed today |
| pending_tasks | int | Yes | Tasks still pending |
| highlights | string | Yes | Key accomplishments |
| tomorrow_focus | string | No | Tomorrow's priorities |

**Message Format:**
```
📊 *DAILY SUMMARY*
_Tuesday, December 24_

✅ Completed: *7*
⏳ Pending: *12*

*Highlights:*
• Telegram HITL integration complete
• Letta agents configured

*Tomorrow's Focus:*
• Review GE Q1 planning
```

---

## Agent Tool Assignments

All 5 agents have these Telegram tools attached:

| Agent | Agent ID | Tools |
|-------|----------|-------|
| Legacy-AI | `agent-da585f6f-189f-472a-974c-29b9578f9b55` | All 5 |
| HP-Conductor | `agent-1e1155fa-de36-475f-90c7-59e1761cb063` | All 5 |
| LE-Conductor | `agent-3ef12371-a11f-48c9-93f3-109dcf87e224` | All 5 |
| GE-Conductor | `agent-052cf675-bcb6-468e-bbec-79e5be465c7a` | All 5 |
| IT-Conductor | `agent-2655ec04-d529-478e-b2b2-56b63b6f5f37` | All 5 |

---

## Development Roadmap

### Phase 1: Quick Wins (Not Started)
**Estimated Time:** ~1 hour total

| Feature | Time | Description |
|---------|------|-------------|
| `/today` command | 15 min | Get today's tasks from Todoist |
| `/hp`, `/le`, `/ge`, `/it` | 15 min | Pillar task summaries |
| `/quick [task]` | 15 min | Quick capture to Todoist inbox |
| Persistent keyboard | 15 min | Always-visible action buttons |

**Implementation:** Requires polling loop or webhook to receive commands

---

### Phase 2: Medium Complexity (Not Started)
**Estimated Time:** ~3-4 hours total

| Feature | Time | Description |
|---------|------|-------------|
| Two-way chat with Legacy-AI | 1-2 hrs | Route messages to Letta, return responses |
| Voice transcription | 1 hr | Whisper API → text → process |
| Image/chart reports | 1-2 hrs | Generate visual summaries |
| Location-based context | 1 hr | Context-aware notifications |

**Dependencies:**
- Two-way chat needs polling loop or webhook
- Voice needs OpenAI Whisper API key
- Charts need matplotlib/pillow

---

### Phase 3: Advanced (Not Started)
**Estimated Time:** ~5-8 hours total
**Requires:** Always-on server (n8n, Railway, VPS)

| Feature | Time | Description |
|---------|------|-------------|
| Webhook receiver | 2-3 hrs | Real-time message handling |
| Multi-turn conversations | 3-4 hrs | Context-aware dialogue |
| File upload to vault | 2-3 hrs | Send docs → GitHub vault |
| Inline queries | 2 hrs | Search OneMind from any chat |

---

## Technical Notes

### Polling vs Webhook

| Aspect | Polling (Current) | Webhook (Future) |
|--------|-------------------|------------------|
| How | Bot checks for updates | Telegram pushes instantly |
| Latency | 1-5 seconds | Instant |
| Requires | Nothing extra | Always-on server |
| Best for | Notifications, approvals | Real-time chat |

### Current Limitations

1. **No incoming message handling** - Bot can send but not actively listen for commands
2. **Polling only when tool called** - `telegram_check_approval` must be called to check responses
3. **No conversation state** - Each message is independent

### To Enable Commands/Chat

Need one of:
1. **Background polling script** - Python script running continuously
2. **n8n workflow** - Telegram trigger node
3. **Webhook server** - Railway/Render/VPS receiving Telegram updates

---

## Code Reference

### Send Message (Python)
```python
import requests

BOT_TOKEN = "8437988392:AAHsEyJFsOr2Qdqyp1B17c4nQgagLWDI8LM"
CHAT_ID = "8297858681"
BASE_URL = f"https://api.telegram.org/bot{BOT_TOKEN}"

def send_message(text, parse_mode="Markdown"):
    return requests.post(
        f"{BASE_URL}/sendMessage",
        json={
            "chat_id": CHAT_ID,
            "text": text,
            "parse_mode": parse_mode
        }
    )
```

### Send with Inline Keyboard
```python
def send_with_buttons(text, buttons):
    keyboard = {
        "inline_keyboard": [
            [{"text": btn["text"], "callback_data": btn["data"]} for btn in row]
            for row in buttons
        ]
    }
    return requests.post(
        f"{BASE_URL}/sendMessage",
        json={
            "chat_id": CHAT_ID,
            "text": text,
            "parse_mode": "Markdown",
            "reply_markup": keyboard
        }
    )

# Usage
send_with_buttons(
    "Choose an option:",
    [[{"text": "✅ Yes", "data": "yes"}, {"text": "❌ No", "data": "no"}]]
)
```

### Poll for Updates
```python
def get_updates(offset=None):
    params = {"timeout": 30}
    if offset:
        params["offset"] = offset
    return requests.get(f"{BASE_URL}/getUpdates", params=params).json()
```

### Handle Callback Query
```python
def handle_callback(callback_query):
    callback_id = callback_query["id"]
    data = callback_query["data"]

    # Acknowledge
    requests.post(
        f"{BASE_URL}/answerCallbackQuery",
        json={"callback_query_id": callback_id, "text": "Received!"}
    )

    return data
```

---

## Letta API Reference

### Create Tool
```python
import requests

LETTA_API_KEY = "sk-let-YjhiY2I3ZjItODZmOC00NTQ2LTlhYTYtNjg2NmNhNDA0ODMzOmM5NjlkODhhLTE0MmEtNGUzMy1iNzMwLTRkYjYzODI2MGZmYQ=="

response = requests.post(
    "https://api.letta.com/v1/tools",
    headers={"Authorization": f"Bearer {LETTA_API_KEY}"},
    json={
        "source_code": "def my_tool(param: str) -> str:\n    ...",
        "tags": ["telegram"]
    }
)
tool_id = response.json()["id"]
```

### Attach Tool to Agent
```python
response = requests.patch(
    f"https://api.letta.com/v1/agents/{agent_id}",
    headers={"Authorization": f"Bearer {LETTA_API_KEY}"},
    json={"tool_ids": [tool_id]}
)
```

### Get Agent Tools
```python
response = requests.get(
    f"https://api.letta.com/v1/agents/{agent_id}",
    headers={"Authorization": f"Bearer {LETTA_API_KEY}"}
)
tools = response.json()["tools"]
```

---

## Quick Test Commands

### Test Bot Connection
```bash
curl "https://api.telegram.org/bot8437988392:AAHsEyJFsOr2Qdqyp1B17c4nQgagLWDI8LM/getMe"
```

### Send Test Message
```bash
curl -X POST "https://api.telegram.org/bot8437988392:AAHsEyJFsOr2Qdqyp1B17c4nQgagLWDI8LM/sendMessage" \
  -H "Content-Type: application/json" \
  -d '{"chat_id": "8297858681", "text": "Test from CLI"}'
```

### Check for Updates
```bash
curl "https://api.telegram.org/bot8437988392:AAHsEyJFsOr2Qdqyp1B17c4nQgagLWDI8LM/getUpdates"
```

---

## Related Files

- [[../Tool Registry/External Tools/todoist_tools.py]] - Todoist integration
- [[../../Task-Engine - Todoist Sync/06-Agent-Integration]] - Agent-Todoist integration
- [[../HITL]] - Human-in-the-loop protocols

---

## Changelog

### 2024-12-24
- Created @onemind_approval_bot
- Built 5 Telegram tools (notify, check, alert, insight, summary)
- Attached tools to all 5 Letta agents
- Tested HITL approval flow successfully
- Documented integration architecture
