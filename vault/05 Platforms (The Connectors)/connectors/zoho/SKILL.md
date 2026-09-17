---
name: zoho
description: Manage Zoho suite including CRM (leads, contacts, deals), Desk (support tickets), and Calendar (events, scheduling). Use when you need to create leads, manage support tickets, schedule meetings, or query customer data.
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Read
  - Write
---

# Zoho Integration Skill

Interact with Zoho CRM, Desk, and Calendar APIs.

## Environment Variables

```bash
ZOHO_CLIENT_ID         # OAuth client ID
ZOHO_CLIENT_SECRET     # OAuth client secret
ZOHO_CRM_REFRESH_TOKEN # CRM refresh token
ZOHO_DESK_REFRESH_TOKEN # Desk refresh token
ZOHO_CALENDAR_REFRESH_TOKEN # Calendar refresh token
```

## Base URLs

```
CRM:      https://www.zohoapis.com/crm/v2
Desk:     https://desk.zoho.com/api/v1
Calendar: https://calendar.zoho.com/api/v1
```

## Authentication

All Zoho APIs use OAuth2 Bearer tokens. Refresh access token:

```bash
curl -X POST "https://accounts.zoho.com/oauth/v2/token" \
  -d "refresh_token=$ZOHO_CRM_REFRESH_TOKEN" \
  -d "client_id=$ZOHO_CLIENT_ID" \
  -d "client_secret=$ZOHO_CLIENT_SECRET" \
  -d "grant_type=refresh_token"
```

---

## Zoho CRM Operations

### List Records
```bash
curl -s "https://www.zohoapis.com/crm/v2/{module}" \
  -H "Authorization: Zoho-oauthtoken $ACCESS_TOKEN" | jq '.data'
```

Modules: `Leads`, `Contacts`, `Accounts`, `Deals`, `Tasks`, `Events`

### Get Record
```bash
curl -s "https://www.zohoapis.com/crm/v2/{module}/{record_id}" \
  -H "Authorization: Zoho-oauthtoken $ACCESS_TOKEN"
```

### Create Lead
```bash
curl -s -X POST "https://www.zohoapis.com/crm/v2/Leads" \
  -H "Authorization: Zoho-oauthtoken $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": [{
      "Last_Name": "Smith",
      "First_Name": "John",
      "Email": "john@example.com",
      "Company": "Acme Inc",
      "Lead_Source": "Website"
    }]
  }'
```

### Create Deal
```bash
curl -s -X POST "https://www.zohoapis.com/crm/v2/Deals" \
  -H "Authorization: Zoho-oauthtoken $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": [{
      "Deal_Name": "Enterprise Contract",
      "Stage": "Qualification",
      "Amount": 50000,
      "Account_Name": "Acme Inc"
    }]
  }'
```

### Search Records
```bash
curl -s "https://www.zohoapis.com/crm/v2/{module}/search?criteria=(Email:equals:john@example.com)" \
  -H "Authorization: Zoho-oauthtoken $ACCESS_TOKEN"
```

---

## Zoho Desk Operations

### List Tickets
```bash
curl -s "https://desk.zoho.com/api/v1/tickets" \
  -H "Authorization: Zoho-oauthtoken $ACCESS_TOKEN" \
  -H "orgId: $ZOHO_ORG_ID"
```

### Create Ticket
```bash
curl -s -X POST "https://desk.zoho.com/api/v1/tickets" \
  -H "Authorization: Zoho-oauthtoken $ACCESS_TOKEN" \
  -H "orgId: $ZOHO_ORG_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "Support Request",
    "description": "Issue details...",
    "contactId": "12345",
    "priority": "High",
    "status": "Open"
  }'
```

### Update Ticket Status
```bash
curl -s -X PATCH "https://desk.zoho.com/api/v1/tickets/{ticket_id}" \
  -H "Authorization: Zoho-oauthtoken $ACCESS_TOKEN" \
  -H "orgId: $ZOHO_ORG_ID" \
  -H "Content-Type: application/json" \
  -d '{"status": "Closed"}'
```

---

## Zoho Calendar Operations

### List Events
```bash
curl -s "https://calendar.zoho.com/api/v1/calendars/{calendar_id}/events" \
  -H "Authorization: Zoho-oauthtoken $ACCESS_TOKEN"
```

### Create Event
```bash
curl -s -X POST "https://calendar.zoho.com/api/v1/calendars/{calendar_id}/events" \
  -H "Authorization: Zoho-oauthtoken $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Meeting",
    "start": {"dateTime": "2025-01-15T10:00:00"},
    "end": {"dateTime": "2025-01-15T11:00:00"},
    "attendees": [{"email": "user@example.com"}]
  }'
```

---

## One Mind Integration

Zoho CRM stages map to GE pillar:
- **Leads** → `75 Zeus Delacruz Enterprise`
- **Deals** → `80 Treasury (Finance)`
- **Support** → `83 Partnerships (Network)`

## Tips

- Access tokens expire in 1 hour; always refresh before operations
- Use MCP server for native access (see .mcp.json)
- Desk requires `orgId` header on all requests


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
