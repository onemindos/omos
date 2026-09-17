# Zoho Bigin

> **Lightweight CRM for Small Business**

## Overview

Bigin is a simplified CRM designed for small businesses and startups who need basic pipeline management without the complexity of full Zoho CRM.

## Status

- **Active**: ✅ Yes
- **Integration**: API
- **Last Updated**: January 25, 2026

## Use Cases

1. **Simple Pipeline** - Manage deals visually
2. **Contact Management** - Store customer info
3. **Activity Tracking** - Log interactions
4. **Mobile-First** - Manage on the go
5. **Quick Setup** - No complex configuration

## Key Features Used

- [ ] Kanban pipeline view
- [ ] Contact management
- [ ] Activity timeline
- [ ] Email integration
- [ ] Mobile app
- [ ] Custom fields
- [ ] Task management

## Bigin vs. Zoho CRM

| Feature           | Bigin           | Zoho CRM        |
| ----------------- | --------------- | --------------- |
| **Pipeline View** | Kanban only     | Multiple views  |
| **Modules**       | Contacts, Deals | 10+ modules     |
| **Automation**    | Basic           | Advanced        |
| **Price**         | $7/user/month   | $14+/user/month |
| **Setup Time**    | 5 minutes       | 1-2 hours       |
| **API**           | Limited         | Full REST API   |

## When to Use Bigin

- **Use Bigin for**:
  - Side projects
  - Freelance work
  - Simple sales tracking
  - Quick MVP

- **Use CRM for**:
  - Enterprise clients
  - Complex workflows
  - Advanced automation
  - Reporting needs

## API Integration

### Authentication

```python
# API Key based
api_key = "bigin_api_key"
base_url = "https://www.zohoapis.com/bigin/v1"
```

### Endpoints

- `GET /Contacts` - List contacts
- `POST /Pipelines` - Create deal
- `PUT /Pipelines/{id}` - Update deal stage

### Rate Limits

- **API Calls**: 1,000/day (free tier)

## Pipeline Stages

### Default Pipeline

1. **New** - Initial contact
2. **Qualified** - Interested prospect
3. **Proposal** - Quote sent
4. **Negotiation** - Discussing terms
5. **Won** - Deal closed
6. **Lost** - Opportunity lost

## Custom Fields

| Field Name     | Type     | Purpose             |
| -------------- | -------- | ------------------- |
| `Project_Size` | Picklist | Small/Medium/Large  |
| `Lead_Source`  | Text     | Where they found us |
| `Budget`       | Currency | Estimated budget    |

## Integrations

- **Zoho Mail** - Email sync
- **Google Calendar** - Meeting sync
- **Stripe** - Payment links
- **Zapier** - Connect other apps

## Mobile App

- Available on iOS and Android
- Offline mode for viewing contacts
- Quick add deal feature
- Activity logging on the go

## Notes

- Much simpler than Zoho CRM
- Limited to 2 custom pipelines on free tier
- No advanced automation (use Zoho CRM if needed)
- Good for personal projects or side businesses
- Can migrate data to Zoho CRM later if needed

---

**Related**: [CRM](./CRM.md) | [SalesIQ](./SalesIQ.md)
