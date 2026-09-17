# Zoho Marketing Automation (MarketingHub)

> **Marketing Workflow Automation & Lead Nurturing**

## Overview

Zoho Marketing Automation (formerly MarketingHub) is an advanced marketing platform for creating multi-channel campaigns, lead scoring, and automated customer journeys.

## Status

- **Active**: ✅ Yes
- **Integration**: API + Webhooks
- **Last Updated**: January 25, 2026

## Use Cases

1. **Lead Nurturing** - Multi-touch campaigns
2. **Lead Scoring** - Qualify prospects automatically
3. **Behavior-Based Campaigns** - Trigger emails based on actions
4. **Multi-Channel Marketing** - Email, SMS, social, web
5. **Customer Journey Mapping** - Visual workflow builder

## Key Features Used

- [ ] Journey builder (visual automation)
- [ ] Lead scoring models
- [ ] Website visitor tracking
- [ ] Progressive profiling
- [ ] Multi-channel campaigns
- [ ] A/B testing
- [ ] Predictive analytics
- [ ] Account-based marketing (ABM)

## Marketing Automation vs. Campaigns

| Feature               | Marketing Automation  | Campaigns             |
| --------------------- | --------------------- | --------------------- |
| **Complexity**        | Advanced workflows    | Simple email blasts   |
| **Lead Scoring**      | ✅ Yes                | ❌ No                 |
| **Multi-Channel**     | Email, SMS, social    | Email only            |
| **Behavior Tracking** | Full website tracking | Email engagement only |
| **Price**             | $30+/month            | $3+/month             |
| **Use Case**          | Complex nurturing     | Newsletters           |

## Lead Scoring

### Scoring Model

#### Demographic Score (Max: 50 points)

- Job Title (C-level = 20, Manager = 10)
- Company Size (Enterprise = 20, SMB = 10)
- Industry (Target industries = 10)

#### Behavioral Score (Max: 50 points)

- Email Open = +2
- Email Click = +5
- Visit Pricing Page = +10
- Download Resource = +8
- Attend Webinar = +15
- Request Demo = +20

#### Total Score Classification

- **0-30**: Cold Lead
- **31-60**: Warm Lead
- **61-80**: Hot Lead
- **81-100**: Sales-Ready

### Negative Scoring

- Unsubscribe = -50
- Spam Complaint = -100
- Bounce = -20
- No Activity (90 days) = -10

## Customer Journeys

### Welcome Journey (New Lead)

```
Trigger: Form Submission
↓
Email 1: Welcome (Immediate)
↓
Wait 2 days
↓
Email 2: Feature Introduction
↓
IF Email Opened → Send Email 3: Case Study
IF Not Opened → Wait 3 days → Resend with different subject
↓
IF Link Clicked → Add to Hot Leads, Notify Sales
```

### Product Trial Journey

```
Trigger: Trial Started
↓
Day 1: Onboarding Email + Success Checklist
↓
Day 3: Check Progress
  IF Feature Used → Encourage Next Step
  IF No Activity → Send Help Resources
↓
Day 7: Mid-Trial Check-in + Offer Demo Call
↓
Day 12: Upgrade Prompt (3 days before trial ends)
↓
Day 14: Final Day Reminder
  IF Converted → Move to Customer Onboarding
  IF Not Converted → Add to Re-engagement Campaign
```

### Re-engagement Journey

```
Trigger: No Activity (60 days)
↓
Email 1: "We miss you" + What's New
↓
Wait 5 days
↓
IF Engaged → Move to Active Lead Flow
IF No Response → Email 2: Limited-Time Offer
↓
Wait 7 days
↓
IF Still No Response → SMS: Personal Message from Founder
↓
IF No Response → Move to Disengaged List
```

## Multi-Channel Campaigns

### Channels Available

- **Email** - Primary channel
- **SMS** - Time-sensitive messages
- **Social Media** - LinkedIn, Twitter, Facebook
- **Web Banners** - On-site personalization
- **Push Notifications** - Browser/mobile

### Example Multi-Channel Flow

1. **Email**: Product update announcement
2. **Social**: Share blog post about update
3. **Web Banner**: Show promo on website
4. **SMS** (if no email open in 3 days): Brief reminder
5. **Push** (app users): In-app notification

## Website Visitor Tracking

### Tracked Actions

- Page visits
- Time on site
- Resources downloaded
- Forms submitted
- Videos watched
- Chat initiated

### Tracking Code

```javascript
<!-- Zoho Marketing Automation Tracking -->
<script>
  var ma = document.createElement('script');
  ma.type = 'text/javascript';
  ma.async = true;
  ma.src = 'https://marketinghub.zoho.com/js/ma.js?org=12345';
  var s = document.getElementsByTagName('script')[0];
  s.parentNode.insertBefore(ma, s);
</script>
```

## Progressive Profiling

Gradually collect information about leads over time instead of long forms upfront.

### Example:

- **Form 1 (First visit)**: Name, Email
- **Form 2 (Second visit)**: Company, Role
- **Form 3 (Third visit)**: Company Size, Industry
- **Form 4 (Fourth visit)**: Phone, Use Case

## Account-Based Marketing (ABM)

Target specific high-value accounts with personalized campaigns.

### ABM Workflow

1. **Identify Target Accounts** - Enterprise prospects
2. **Create Account Lists** - Group by industry/size
3. **Personalize Content** - Tailor messaging per account
4. **Multi-Touch Campaign** - Reach multiple stakeholders
5. **Measure Engagement** - Account-level analytics

## Automation Rules

### Email Engagement

- **Rule**: If contact opens 5+ emails in 30 days
- **Action**: Add tag "Highly Engaged", notify sales

### Website Behavior

- **Rule**: If contact visits pricing page 3+ times
- **Action**: Add to Hot Leads, create CRM Task

### Inactivity

- **Rule**: If no activity for 90 days
- **Action**: Move to Re-engagement campaign

### Form Submission

- **Rule**: If downloads whitepaper
- **Action**: Send follow-up email series (5 emails over 10 days)

## API Integration

### Endpoints

```python
# Add contact to journey
POST /api/v1/journeys/{id}/contacts

# Update lead score
PUT /api/v1/contacts/{id}/score

# Track custom event
POST /api/v1/events
{
  "contact_email": "user@example.com",
  "event_name": "feature_used",
  "event_data": {"feature": "ai_agent"}
}
```

## Integrations

- **Zoho CRM** - Bi-directional sync
- **Zoho Campaigns** - Share contacts & templates
- **Zoho SalesIQ** - Trigger campaigns from chat
- **Zoho Analytics** - Advanced reporting
- **Google Ads** - Sync audiences
- **Facebook Ads** - Custom audiences
- **Salesforce** - Alternative CRM sync

## Analytics & Reporting

### Journey Performance

- Conversion rate per step
- Drop-off points
- Average time to convert
- Revenue attributed

### Channel Effectiveness

- Open rates by channel
- Click rates by channel
- Conversion rates by channel
- Cost per acquisition

### Lead Scoring Report

- Distribution of scores
- Hot leads pipeline
- Score changes over time

## Legacy AI Integration

### Automated Actions

1. **High Lead Score** → Create Todoist task for sales call
2. **Journey Completion** → Update CRM deal stage
3. **Trial Expiring** → Send alert to Cliq
4. **Custom Event** → Log to PostgreSQL for analysis
5. **Account Engagement Spike** → Flag for ABM campaign

## Best Practices

- ✅ Start with simple journeys, add complexity later
- ✅ Use lead scoring to prioritize sales outreach
- ✅ Personalize content based on behavior
- ✅ Test different wait times between emails
- ✅ Monitor journey drop-off points
- ✅ Don't over-communicate (frequency caps)
- ✅ Always provide value in each touchpoint
- ✅ Clean lists regularly (remove disengaged)

## Notes

- More powerful than Zoho Campaigns but more complex
- Requires higher-tier plan ($30+/month)
- Website tracking requires JavaScript snippet on all pages
- Lead scoring models should be reviewed quarterly
- Journey builder has learning curve, invest time in training
- Multi-channel requires additional credits (SMS, social ads)

---

**Related**: [Campaigns](Campaigns.md) | [SalesIQ](SalesIQ.md) | [PageSense](PageSense.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
