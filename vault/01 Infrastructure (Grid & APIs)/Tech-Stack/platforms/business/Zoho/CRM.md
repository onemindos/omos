# Zoho CRM

> **Customer Relationship Management & Sales Pipeline**

## Overview

Zoho CRM helps manage customer relationships, sales pipelines, and track deals from lead to close.

## Status

- **Active**: ✅ Yes
- **Integration**: API + OAuth
- **Last Updated**: January 25, 2026

## Use Cases

1. **Lead Management** - Track potential customers
2. **Deal Pipeline** - Manage sales opportunities
3. **Contact Management** - Store customer information
4. **Activity Tracking** - Log calls, meetings, tasks
5. **Reporting** - Sales analytics and forecasting

## Key Features Used

- [ ] Leads module
- [ ] Contacts module
- [ ] Deals/Opportunities
- [ ] Accounts
- [ ] Tasks & Activities
- [ ] Reports & Dashboards
- [ ] Email integration
- [ ] Workflow automation

## API Integration

### Authentication

```python
# OAuth 2.0
client_id = "your_client_id"
client_secret = "your_client_secret"
redirect_uri = "https://onemind.zeusdelacruz.com/oauth/zoho"
```

### Endpoints Used

- `GET /crm/v2/Leads` - List leads
- `POST /crm/v2/Leads` - Create lead
- `GET /crm/v2/Deals` - List deals
- `PUT /crm/v2/Deals/{id}` - Update deal

### Rate Limits

- **API Calls**: 5,000/day (free tier)
- **Bulk Operations**: 200 records/request

## Workflows

### Lead to Customer Flow

1. Lead captured (form, SalesIQ, manual)
2. Lead qualified by Legacy AI
3. Lead converted to Deal
4. Deal moved through stages
5. Deal won → Customer created

## Zia AI Integration

- **Lead Scoring** - Predict deal likelihood
- **Best Time to Contact** - Optimal outreach timing
- **Email Sentiment** - Analyze customer emails
- **Anomaly Detection** - Flag unusual patterns

## Custom Fields

| Field Name        | Type         | Purpose                |
| ----------------- | ------------ | ---------------------- |
| `Legacy_AI_Score` | Number       | AI qualification score |
| `Source_Platform` | Picklist     | Where lead came from   |
| `Project_Type`    | Multi-select | Type of engagement     |

## Automation Rules

1. **New Lead Alert** - Notify Legacy AI of new lead
2. **Inactive Lead** - Flag leads with no activity for 30 days
3. **Deal Stage Change** - Update Todoist task
4. **Won Deal** - Send to billing system

## Reports & Dashboards

- **Sales Pipeline** - Visual deal stages
- **Lead Conversion Rate** - Lead → Deal conversion
- **Revenue Forecast** - Projected monthly revenue
- **Activity Summary** - Calls, meetings, emails

## Integrations

- **Zoho Mail** - Email sync
- **Zoho Desk** - Support ticket linking
- **Google Calendar** - Meeting sync
- **Todoist** - Task management
- **Legacy AI** - Automated actions

## Notes

- Free tier has limited API calls (5K/day)
- Bulk operations more efficient than individual calls
- Zia AI requires higher-tier plan for advanced features
- Custom modules available on Enterprise plan

---

**Related**: [Desk](./Desk.md) | [Bigin](./Bigin.md) | [SalesIQ](./SalesIQ.md)
