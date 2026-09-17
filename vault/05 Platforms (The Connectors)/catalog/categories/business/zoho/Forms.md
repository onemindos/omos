# Zoho Forms

> **Online Form Builder & Data Collection**

## Overview

Zoho Forms is a drag-and-drop form builder for creating surveys, contact forms, registration forms, and collecting data online.

## Status

- **Active**: ✅ Yes
- **Integration**: API + Webhooks + Embed
- **Last Updated**: January 25, 2026

## Use Cases

1. **Lead Capture** - Contact/inquiry forms
2. **Event Registration** - Workshop sign-ups
3. **Surveys** - Customer feedback
4. **Order Forms** - Product requests
5. **Applications** - Job/project applications

## Key Features Used

- [ ] Drag-and-drop builder
- [ ] Conditional logic (show/hide fields)
- [ ] File uploads
- [ ] Payment integration (Stripe, PayPal)
- [ ] Email notifications
- [ ] Webhooks to Legacy AI
- [ ] Multi-page forms
- [ ] Mobile responsive

## Active Forms

| Form Name               | Purpose               | Submissions/Month |
| ----------------------- | --------------------- | ----------------- |
| **Contact Form**        | General inquiries     | ~20               |
| **Project Application** | New project requests  | ~5                |
| **Bug Report**          | Issue reporting       | ~10               |
| **Event Registration**  | Workshop sign-ups     | Varies            |
| **Feedback Survey**     | Post-service feedback | ~15               |

## API Integration

### Webhooks

```json
{
  "trigger": "submission",
  "url": "https://onemind.zeusdelacruz.com/webhooks/zoho/forms",
  "method": "POST",
  "headers": {
    "X-Zoho-Form": "contact_form"
  }
}
```

### Get Submissions

```python
GET https://forms.zoho.com/api/v1/{form_id}/submissions
```

### Rate Limits

- **Submissions**: Unlimited (free tier has 500/month limit)
- **API Calls**: 1,000/day

## Form Flow

### Standard Submission Flow

1. User fills form
2. Form validation (client-side)
3. Submit to Zoho Forms
4. Webhook fires to Legacy AI
5. Legacy AI processes submission
6. Create CRM lead / Desk ticket / Todoist task
7. Send confirmation email

## Form Fields

### Common Field Types

- **Text** - Short answer
- **Textarea** - Long answer
- **Email** - Email validation
- **Phone** - Phone number
- **Dropdown** - Select from list
- **Checkbox** - Multiple selection
- **Radio** - Single selection
- **File Upload** - Document/image
- **Date** - Date picker
- **Rating** - Star rating
- **Payment** - Stripe integration

## Conditional Logic Example

```
IF Budget = "Enterprise ($50K+)"
  THEN Show "Enterprise Features" section
  AND Show "Schedule Call" field
ELSE
  Show "Get Started" button
```

## Automation Rules

1. **New Submission** → Create CRM lead
2. **Bug Report** → Create Desk ticket (Priority: High)
3. **Event Registration** → Add to Google Calendar
4. **Payment Form** → Send receipt email
5. **Feedback Score < 3** → Alert team in Cliq

## Integrations

- **Zoho CRM** - Auto-create leads/contacts
- **Zoho Desk** - Create support tickets
- **Zoho Campaigns** - Add to mailing list
- **Google Sheets** - Export submissions
- **Stripe/PayPal** - Payment collection
- **Zapier** - Connect to 1000+ apps

## Embedding Forms

### Inline Embed

```html
<iframe
  src="https://forms.zohopublic.com/username/form/ContactForm"
  width="100%"
  height="600px"
></iframe>
```

### Popup/Modal

```html
<script src="https://forms.zohopublic.com/js/popup.js"></script>
<a href="javascript:void(0)" onclick="openForm()">Open Form</a>
```

### Script Embed (AJAX)

- Form loads within page without iframe
- Better user experience
- Maintains page context

## Email Notifications

### To Admin

- Instant notification on submission
- Include all field values
- Link to full submission

### To User (Auto-responder)

- Confirmation of submission
- Thank you message
- Next steps
- Estimated response time

## GDPR & Privacy

- [ ] Add GDPR consent checkbox
- [ ] Include privacy policy link
- [ ] Data retention settings (auto-delete after X days)
- [ ] Export user data on request
- [ ] IP address tracking toggle

## Analytics

- **Submission Rate** - How many complete the form
- **Drop-off Points** - Where users abandon
- **Average Time** - Time to complete
- **Device Breakdown** - Mobile vs. desktop
- **Source Tracking** - UTM parameters

## Custom CSS/Branding

```css
/* Custom form styling */
.zf-form-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.zf-submit-button {
  background: #667eea;
  border-radius: 8px;
}
```

## Notes

- Free tier: Up to 500 submissions/month
- Unlimited forms on all plans
- File uploads limited to 10MB per file (free tier)
- Payment forms require paid plan ($12+/month)
- Multi-page forms improve completion rate for long forms
- Mobile optimization is automatic

---

**Related**: [CRM](CRM.md) | [Desk](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/Tech-Stack/platforms/business/Zoho/Desk.md) | [Campaigns](Campaigns.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
