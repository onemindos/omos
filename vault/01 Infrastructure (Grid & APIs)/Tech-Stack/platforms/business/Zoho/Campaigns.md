# Zoho Campaigns

> **Email Marketing & Newsletter Management**

## Overview

Zoho Campaigns is an email marketing platform for creating, sending, and tracking email campaigns and newsletters.

## Status

- **Active**: ✅ Yes
- **Integration**: API + Webhooks
- **Last Updated**: January 25, 2026

## Use Cases

1. **Newsletter** - Monthly updates
2. **Product Launches** - New feature announcements
3. **Drip Campaigns** - Automated email sequences
4. **Event Invites** - Workshop/webinar invitations
5. **Re-engagement** - Win back inactive users

## Key Features Used

- [ ] Email template builder (drag-and-drop)
- [ ] Contact list management
- [ ] Segmentation & targeting
- [ ] A/B testing
- [ ] Automation workflows
- [ ] Analytics & tracking
- [ ] Social media integration
- [ ] RSS to email

## Active Campaigns

| Campaign Name           | Type            | Frequency     | Open Rate |
| ----------------------- | --------------- | ------------- | --------- |
| **OneMind Monthly**     | Newsletter      | Monthly       | ~35%      |
| **New User Onboarding** | Drip (5 emails) | Auto-trigger  | ~50%      |
| **Feature Updates**     | One-time        | As needed     | ~40%      |
| **Abandoned Cart**      | Automation      | Trigger-based | ~25%      |

## API Integration

### Send Campaign

```python
POST https://campaigns.zoho.com/api/v1.1/sendcampaign
{
  "list_key": "subscriber_list",
  "campaign_key": "monthly_newsletter",
  "send_time": "2026-02-01T10:00:00Z"
}
```

### Manage Subscribers

```python
# Add subscriber
POST /api/v1.1/{list_key}/addcontact

# Update contact
POST /api/v1.1/{list_key}/updatecontact

# Unsubscribe
POST /api/v1.1/{list_key}/unsubscribe
```

### Webhooks

- `subscriber.added` - New subscriber
- `email.opened` - Email opened
- `link.clicked` - Link clicked
- `subscriber.unsubscribed` - Unsubscribe event
- `email.bounced` - Bounce notification

### Rate Limits

- **API Calls**: 3,000/day
- **Email Sends**: Based on plan (12K free/month)

## Contact Lists

| List Name           | Subscribers | Segment              |
| ------------------- | ----------- | -------------------- |
| **All Contacts**    | ~500        | Everyone             |
| **Active Users**    | ~300        | Last login < 30 days |
| **Trial Users**     | ~100        | In trial period      |
| **Enterprise**      | ~50         | Paid customers       |
| **Newsletter Only** | ~200        | Opted in for news    |

## Email Templates

### Newsletter Template

- Header with logo
- Hero image/CTA
- 3-column feature section
- Blog post roundup
- Footer with social links

### Transactional Template

- Simple, clean design
- Clear CTA button
- Minimal distractions
- Mobile-optimized

### Event Invite Template

- Event details (date, time, location)
- Speaker/agenda preview
- Register button (prominent)
- Add to calendar link

## Automation Workflows

### Welcome Series (5 emails)

1. **Day 0**: Welcome + Get started guide
2. **Day 2**: Feature highlight #1
3. **Day 5**: Feature highlight #2
4. **Day 7**: Success stories
5. **Day 10**: Feedback request

### Re-engagement Campaign

1. **Trigger**: No activity for 60 days
2. **Email 1**: "We miss you" + What's new
3. **Wait 5 days**
4. **Email 2**: Special offer/incentive
5. **Wait 7 days**
6. **Email 3**: Final reminder

## Segmentation

### By Engagement

- **Highly Engaged** - Opened last 3 emails
- **Moderately Engaged** - Opened 1 of last 3
- **Disengaged** - No opens in 6 months

### By User Type

- **Free Users** - No payment
- **Trial Users** - In trial period
- **Paid Users** - Active subscription
- **Churned Users** - Cancelled subscription

### By Interest

- **Tech/Development** - Clicked dev articles
- **Business/Strategy** - Clicked business content
- **Personal Use** - Individual users

## A/B Testing

### Test Elements

- **Subject Line** - Different variations
- **Send Time** - Morning vs. afternoon
- **CTA Button** - Color, text, placement
- **Images** - Different hero images
- **Content Length** - Long-form vs. short

### Sample Size

- Minimum 100 subscribers per variant
- Test 20% of list, send winner to remaining 80%

## Email Deliverability

### Best Practices

- ✅ SPF record configured
- ✅ DKIM signature enabled
- ✅ DMARC policy set
- ✅ Custom domain for sending
- ✅ Double opt-in confirmation
- ✅ Easy unsubscribe link
- ✅ Regular list cleaning (remove bounces)

### Bounce Management

- **Hard Bounce** - Remove from list immediately
- **Soft Bounce** - Retry 3 times, then remove
- **Complaint** - Unsubscribe and flag

## Analytics & Reports

### Campaign Performance

- **Open Rate** - % who opened email
- **Click Rate** - % who clicked links
- **Bounce Rate** - % bounced
- **Unsubscribe Rate** - % unsubscribed
- **Spam Complaints** - Marked as spam

### Engagement Over Time

- Track opens/clicks by hour
- Best send time analysis
- Geographic data (where opens occur)
- Device breakdown (mobile vs. desktop)

## Integrations

- **Zoho CRM** - Sync contacts & leads
- **Zoho Forms** - Add form submissions to list
- **Zoho SalesIQ** - Trigger emails based on website behavior
- **Zoho Analytics** - Advanced reporting
- **WordPress** - Embed signup forms
- **Social Media** - Share to Facebook, Twitter, LinkedIn

## Compliance (GDPR/CAN-SPAM)

- ✅ Explicit consent required (double opt-in)
- ✅ Clear unsubscribe link in every email
- ✅ Physical address in footer
- ✅ Process unsubscribe requests within 10 days
- ✅ Never buy/rent email lists
- ✅ Maintain suppression list

## Legacy AI Integration

### Automated Actions

1. **New Subscriber** → Add to CRM
2. **High Engagement** → Flag as hot lead
3. **Unsubscribe** → Update CRM status
4. **Link Click (pricing page)** → Create Todoist task for follow-up
5. **Campaign Performance** → Weekly summary to Cliq

## Notes

- Free tier: 2,000 subscribers, 12,000 emails/month
- Paid plans start at $3/month (2,500 subscribers)
- Custom domain required for deliverability
- Avoid spam trigger words in subject lines
- Test emails before sending to full list
- Keep lists clean (remove inactive/bounced)
- Optimal send time: Tuesday-Thursday, 10 AM - 2 PM

---

**Related**: [Marketing Automation](./Marketing-Automation.md) | [Forms](./Forms.md) | [CRM](./CRM.md)
