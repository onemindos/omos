# Zoho Desk

> **Customer Support Ticketing System**

## Overview

Zoho Desk is a help desk and customer support platform for managing support tickets, knowledge base, and customer service operations.

## Status

- **Active**: ✅ Yes
- **Integration**: API + Webhooks
- **Last Updated**: January 25, 2026

## Use Cases

1. **Support Tickets** - Manage customer inquiries
2. **Knowledge Base** - Self-service articles
3. **Multi-Channel Support** - Email, chat, social, phone
4. **SLA Management** - Track response times
5. **Customer Satisfaction** - CSAT surveys

## Key Features Used

- [ ] Ticket management
- [ ] Email-to-ticket conversion
- [ ] Knowledge base articles
- [ ] Canned responses (templates)
- [ ] Ticket routing rules
- [ ] SLA policies
- [ ] Customer portal
- [ ] Reporting & analytics

## API Integration

### Endpoints Used

- `GET /api/v1/tickets` - List tickets
- `POST /api/v1/tickets` - Create ticket
- `PUT /api/v1/tickets/{id}` - Update ticket
- `GET /api/v1/articles` - Knowledge base articles

### Webhooks

```json
{
  "event": "ticket.created",
  "url": "https://onemind.zeusdelacruz.com/webhooks/zoho/desk",
  "events": ["ticket.created", "ticket.updated", "ticket.closed"]
}
```

### Rate Limits

- **API Calls**: 3,000/day (free tier)
- **Webhooks**: Unlimited

## Ticket Workflow

### Standard Flow

1. Customer submits ticket (email, form, chat)
2. Auto-assignment based on category
3. Agent/Legacy AI responds
4. Customer replies (updates ticket)
5. Resolution & close
6. CSAT survey sent

### Automation Rules

1. **High Priority Alert** - Notify immediately
2. **No Response** - Escalate after 24 hours
3. **Auto-Close** - Close after 7 days inactive (solved status)
4. **Knowledge Base Suggestion** - Legacy AI suggests articles

## Ticket Categories

| Category          | Department  | SLA      |
| ----------------- | ----------- | -------- |
| Technical Support | Engineering | 4 hours  |
| Billing           | Finance     | 8 hours  |
| General Inquiry   | Sales       | 24 hours |
| Bug Report        | Product     | 6 hours  |

## Zia AI Assistant

- **Answer Bot** - Auto-suggest responses
- **Sentiment Analysis** - Detect frustrated customers
- **Tag Suggestions** - Auto-categorize tickets
- **Similar Tickets** - Find related issues

## Knowledge Base Structure

```
Knowledge Base/
├── Getting Started/
├── Technical Documentation/
├── Billing & Payments/
├── Troubleshooting/
└── FAQs/
```

## SLA Policies

| Priority | Response Time | Resolution Time |
| -------- | ------------- | --------------- |
| Critical | 1 hour        | 4 hours         |
| High     | 4 hours       | 1 day           |
| Medium   | 8 hours       | 3 days          |
| Low      | 24 hours      | 7 days          |

## Integrations

- **Zoho CRM** - Link tickets to customers
- **Zoho Cliq** - Team notifications
- **Zoho Mail** - Email support
- **SalesIQ** - Live chat to ticket
- **Legacy AI** - Automated responses & routing

## Canned Responses

1. **Welcome Message** - Initial acknowledgment
2. **Working On It** - Investigation update
3. **Resolution Provided** - Solution sent
4. **Additional Info Needed** - Request more details
5. **Closing Message** - Thank you & feedback request

## Custom Fields

| Field Name               | Type       | Purpose                    |
| ------------------------ | ---------- | -------------------------- |
| `AI_Response_Confidence` | Percentage | Legacy AI confidence score |
| `Escalation_Level`       | Number     | How many times escalated   |
| `Related_Project`        | Lookup     | Link to CRM deal           |

## Reports

- **Ticket Volume** - Daily/weekly ticket count
- **Response Time** - Average first response time
- **Resolution Time** - Average time to close
- **Customer Satisfaction** - CSAT scores
- **Agent Performance** - Individual metrics

## Notes

- Email-to-ticket creates tickets from support@ emails
- Knowledge base reduces ticket volume by 30-40%
- Zia AI requires Professional+ plan
- Multi-brand support available (different portals per brand)

---

**Related**: [CRM](./CRM.md) | [SalesIQ](./SalesIQ.md) | [Cliq](./Cliq.md)
