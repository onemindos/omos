# Zoho SalesIQ

> **Live Chat & Website Visitor Tracking**

## Overview

Zoho SalesIQ is a live chat and visitor tracking platform that helps engage website visitors in real-time, qualify leads, and convert more prospects.

## Status

- **Active**: ✅ Yes
- **Integration**: JavaScript Widget + API + Webhooks
- **Last Updated**: January 25, 2026

## Use Cases

1. **Live Chat** - Real-time visitor support
2. **Visitor Tracking** - See who's on your site
3. **Lead Qualification** - Chat-based qualification
4. **Proactive Engagement** - Trigger chat based on behavior
5. **Bot Automation** - Handle common questions

## Key Features Used

- [ ] Live chat widget
- [ ] Visitor tracking & history
- [ ] Chatbots (Zobot)
- [ ] Screen sharing & co-browsing
- [ ] Audio/video calls
- [ ] Triggered messages
- [ ] Lead scoring
- [ ] Mobile apps (iOS/Android)

## Chat Widget Setup

### Installation

```html
<!-- SalesIQ Widget -->
<script>
  var $zoho = $zoho || {};
  $zoho.salesiq = $zoho.salesiq || {
    widgetcode: "siq123abc",
    values: {},
    ready: function () {},
  };
  var d = document;
  s = d.createElement("script");
  s.type = "text/javascript";
  s.id = "zsiqscript";
  s.defer = true;
  s.src = "https://salesiq.zoho.com/widget";
  t = d.getElementsByTagName("script")[0];
  t.parentNode.insertBefore(s, t);
</script>
```

### Widget Customization

```javascript
// Custom branding
$zoho.salesiq.ready = function () {
  $zoho.salesiq.floatbutton.color("#667eea");
  $zoho.salesiq.floatbutton.position("bottomright");
  $zoho.salesiq.language("en");
};
```

## Visitor Tracking

### Information Captured

- **Page Views**: Pages visited, time on each
- **Referral Source**: How they found you
- **Location**: Country, city
- **Device**: Desktop, mobile, OS, browser
- **Company Info**: IP-based company name
- **Previous Visits**: Returning visitor data
- **Chat History**: Past conversations

### Visitor Classification

- **Unique Visitors**: First-time visitors
- **Returning Visitors**: Been here before
- **Tracked Visitors**: Known contacts (CRM sync)
- **Hot Leads**: High engagement score

## Live Chat

### Chat Features

- **Canned Responses**: Quick replies for common questions
- **File Sharing**: Send documents/images
- **Screen Sharing**: See visitor's screen
- **Co-Browsing**: Navigate pages together
- **Audio/Video Calls**: Escalate to call
- **Chat Transfer**: Transfer to teammate
- **Typing Indicators**: See when visitor is typing

### Chat Routing

- **Round-Robin**: Distribute chats evenly
- **Department-Based**: Route by topic (Sales, Support, Billing)
- **Skills-Based**: Route by agent expertise
- **VIP Routing**: Hot leads to senior agents

### Offline Behavior

- Show contact form when agents offline
- Auto-respond with business hours
- Send email notification of message

## Chatbots (Zobot)

### Bot Capabilities

- Answer FAQs
- Qualify leads
- Book meetings
- Collect contact info
- Route to human agent
- Multi-language support

### Example Bot Flow: Lead Qualification

```
Bot: Hi! 👋 What brings you to OneMind today?

Options:
1. I want to learn more
2. I need help with something
3. I'm ready to sign up

IF "Learn more":
  → Send resource links
  → Capture email for follow-up

IF "Need help":
  → Ask problem description
  → Route to Support team

IF "Sign up":
  → Capture: Name, Email, Company
  → Calculate lead score
  → IF score > 70: Transfer to Sales
  → IF score < 70: Send self-service resource
```

### Bot Triggers

- **Time on Site**: After 30 seconds
- **Scroll Depth**: At 50% scroll
- **Exit Intent**: When mouse leaves page
- **Specific Pages**: On pricing or demo pages
- **Returning Visitor**: "Welcome back!"

## Proactive Engagement

### Triggered Messages

**Scenario 1: Pricing Page Visitor**

- **Trigger**: User on pricing page for 20+ seconds
- **Message**: "Hi! Have questions about pricing? I'm here to help!"

**Scenario 2: Multiple Page Views**

- **Trigger**: Visited 5+ pages
- **Message**: "I see you're exploring OneMind. Want a quick tour?"

**Scenario 3: Abandoned Cart**

- **Trigger**: Added to cart but leaving site
- **Message**: "Wait! Complete your order and get 10% off today!"

**Scenario 4: Returning Visitor**

- **Trigger**: Visited 3+ times, not signed up
- **Message**: "Welcome back! Ready to start your free trial?"

## Lead Scoring

### Scoring Criteria

- **Page Views** = +5 per page
- **Time on Site** = +2 per minute
- **Chat Initiated** = +15
- **Demo Requested** = +30
- **Pricing Page Viewed** = +20
- **Whitepaper Downloaded** = +10

### Score Thresholds

- **0-30**: Cold (auto-respond with resources)
- **31-60**: Warm (assign to sales queue)
- **61-100**: Hot (notify sales immediately)

## Integration with Zoho CRM

### Automatic Actions

1. **Chat Start** → Create CRM Lead
2. **Email Captured** → Update CRM Contact
3. **Meeting Booked** → Create CRM Event
4. **Deal Interest** → Create CRM Deal
5. **Support Question** → Create Desk Ticket

### CRM Data in Chat

- See contact's CRM profile during chat
- View deal history
- Check open support tickets
- Add notes to CRM record

## Mobile Apps

### Agent App Features

- Receive chat notifications
- Respond from mobile
- See visitor info
- Transfer chats
- Set availability status

### Use Cases

- On-the-go support
- After-hours coverage
- Field team access

## Analytics & Reports

### Chat Metrics

- **Total Chats**: Daily/weekly/monthly
- **Avg. Response Time**: Time to first reply
- **Avg. Chat Duration**: Length of conversation
- **Chat Rating**: Customer satisfaction
- **Conversion Rate**: % chats → leads

### Agent Performance

- **Chats Handled**: Per agent
- **Response Time**: Per agent
- **Customer Rating**: Per agent
- **Availability**: Online time

### Visitor Metrics

- **Total Visitors**: Tracked visitors
- **Conversion Rate**: Visitors → leads
- **Avg. Time on Site**: Session duration
- **Top Pages**: Most visited pages

## API Integration

### Track Custom Events

```javascript
// Track feature usage
$zoho.salesiq.visitor.customaction("feature_used", {
  feature: "ai_agent",
  timestamp: Date.now(),
});
```

### Set Visitor Info

```javascript
// Identify known visitors
$zoho.salesiq.visitor.info({
  name: "Zeus Delacruz",
  email: "zeus@onemind.com",
  company: "OneMind",
});
```

### Webhooks

- `chat.started` - New chat initiated
- `chat.ended` - Chat completed
- `chat.missed` - No agent available
- `message.received` - New message from visitor
- `lead.captured` - Email captured

## Legacy AI Integration

### Automated Actions

1. **High Score Visitor** → Alert in Cliq: "Hot lead on site!"
2. **Chat Transcript** → Analyze sentiment, create summary
3. **Question Pattern** → Update FAQ/knowledge base
4. **Missed Chat** → Create Todoist task for follow-up
5. **Feature Request** → Create Linear issue

### AI-Powered Responses

- Legacy AI can monitor chats
- Suggest responses to agents
- Auto-respond to common questions (when approved)
- Escalate complex questions to human

## Best Practices

- ✅ Respond within 30 seconds (increases conversion by 40%)
- ✅ Use visitor name in chat (if known)
- ✅ Keep responses concise and friendly
- ✅ Use canned responses for speed
- ✅ Always ask "Is there anything else?"
- ✅ End with clear next step (CTA)
- ✅ Follow up missed chats via email
- ✅ Review chat transcripts for improvement

## Privacy & Compliance

- ✅ GDPR compliant
- ✅ Chat data encryption
- ✅ Visitor opt-out available
- ✅ Data retention settings
- ✅ Cookie consent integration

## Notes

- Free tier: 1 operator, 100 chats/month
- Paid plans: $7+/month per operator
- Chatbots require Professional+ plan
- Screen sharing requires enterprise plan
- Mobile apps free for all plans
- Integrates with Zoho Desk (chat → ticket)

---

**Related**: [Desk](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/Tech-Stack/platforms/business/Zoho/Desk.md) | [CRM](CRM.md) | [PageSense](PageSense.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
