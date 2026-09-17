---
id: ui-10-20260111053926
title: "Legacy AI - Tool Upload Checklist"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/10]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Legacy AI - Tool Upload Checklist"
---

# Legacy AI - Tool Upload Checklist

All Python tools that need to be in Letta Cloud for Legacy AI.

## Quick Upload Steps

1. Go to: https://app.letta.com/
2. Select: Legacy AI (`agent-da585f6f-189f-472a-974c-29b9578f9b55`)
3. Navigate to: Tools section
4. For each tool: Add Tool → paste function code → Save

---

## Tool Status

| Tool | File | Status | Env Vars Needed |
|------|------|--------|-----------------|
| `github()` | legacy_tools.py | Verify | `GITHUB_TOKEN`, `GITHUB_REPO` |
| `todoist()` | legacy_tools.py | Verify | `TODOIST_API_TOKEN` |
| `comms()` | legacy_tools.py | Verify | `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID` |
| `zoho_crm()` | zoho_tools.py | **Upload** | See Zoho section |
| `zoho_desk()` | zoho_tools.py | **Upload** | See Zoho section |
| `zoho_calendar()` | zoho_tools.py | **Upload** | See Zoho section |
| `circle()` | circle_tools.py | **Upload** | `CIRCLE_API_TOKEN`, `CIRCLE_COMMUNITY_ID` |
| `letta_admin()` | letta_admin_tools.py | **Upload** | `LETTA_API_KEY` (already set) |
| `delegate_to_hp()` | delegation_tools.py | Verify | `LETTA_API_KEY` |
| `delegate_to_le()` | delegation_tools.py | Verify | `LETTA_API_KEY` |
| `delegate_to_ge()` | delegation_tools.py | Verify | `LETTA_API_KEY` |
| `delegate_to_it()` | delegation_tools.py | Verify | `LETTA_API_KEY` |
| `delegate_to_all()` | delegation_tools.py | Verify | `LETTA_API_KEY` |

---

## Environment Variables to Set

### Already Configured (verify)
```
GITHUB_TOKEN=<github_pat>
GITHUB_REPO=Zeus-Delacruz/OneMind-Codex
TODOIST_API_TOKEN=<todoist_token>
TELEGRAM_BOT_TOKEN=<bot_token>
TELEGRAM_CHAT_ID=<chat_id>
LETTA_API_KEY=<letta_key>
```

### Zoho (setup required)
```
ZOHO_CLIENT_ID=<your_client_id>
ZOHO_CLIENT_SECRET=<your_client_secret>
ZOHO_CRM_REFRESH_TOKEN=<crm_refresh_token>
ZOHO_DESK_REFRESH_TOKEN=<desk_refresh_token>
ZOHO_DESK_ORG_ID=<your_org_id>
ZOHO_CALENDAR_REFRESH_TOKEN=<calendar_refresh_token>
```

### Circle.so (setup required)
```
CIRCLE_API_TOKEN=<api_token>
CIRCLE_COMMUNITY_ID=<community_id>
```

---

## Upload Order (Recommended)

### Phase 1: Admin Tool (enables programmatic uploads)
1. [ ] `letta_admin()` - Upload manually first

### Phase 2: Core Tools (verify existing)
2. [ ] `github()` - Verify attached
3. [ ] `todoist()` - Verify attached
4. [ ] `comms()` - Verify attached

### Phase 3: Delegation Tools (for conductor architecture)
5. [ ] `delegate_to_hp()`
6. [ ] `delegate_to_le()`
7. [ ] `delegate_to_ge()`
8. [ ] `delegate_to_it()`
9. [ ] `delegate_to_all()`

### Phase 4: New Integrations (upload now, configure later)
10. [ ] `zoho_crm()` - Upload, env vars later
11. [ ] `zoho_desk()` - Upload, env vars later
12. [ ] `zoho_calendar()` - Upload, env vars later
13. [ ] `circle()` - Upload, env vars later

---

## File Locations

```
07.4 letta/tools/
├── legacy_tools.py        # github(), todoist(), comms()
├── zoho_tools.py          # zoho_crm(), zoho_desk(), zoho_calendar()
├── circle_tools.py        # circle()
├── delegation_tools.py    # delegate_to_*()
├── letta_admin_tools.py   # letta_admin()
├── todoist_github_tools.py # (older individual tools - deprecated)
└── discord_tools.py       # (if needed)
```

---

## Testing Commands

After upload, test each tool:

```python
# Test admin
letta_admin(action="agents_tools")  # List all attached tools

# Test core tools
github(action="list", path="docs")
todoist(action="get", filter="today")
comms(action="send", message="Test from Legacy", channel="telegram")

# Test Zoho (after env vars set)
zoho_crm(action="modules")
zoho_desk(action="departments")
zoho_calendar(action="calendars_list")

# Test Circle (after env vars set)
circle(action="spaces_list")

# Test delegation
delegate_to_hp(task="What tools do you have?")
```

---

## Notes

1. **Tools uploaded without env vars** will return errors when called - that's expected until you configure credentials.

2. **letta_admin() is special** - once uploaded, Legacy can use it to manage her own tools.

3. **Source files in vault** are the source of truth. If you update code locally, re-upload to Letta Cloud.

4. **Letta dashboard** at https://app.letta.com/ shows all attached tools and lets you edit them directly.

---

*Last updated: 2026-01-09*
