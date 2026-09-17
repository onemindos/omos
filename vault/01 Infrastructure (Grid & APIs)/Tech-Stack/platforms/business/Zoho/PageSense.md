# Zoho PageSense

> **Website Analytics & Conversion Optimization**

## Overview

Zoho PageSense is a website analytics and conversion rate optimization (CRO) platform with heatmaps, session recordings, A/B testing, and funnel analysis.

## Status

- **Active**: ✅ Yes
- **Integration**: JavaScript Tracking
- **Last Updated**: January 25, 2026

## Use Cases

1. **Heatmaps** - Visualize clicks, scrolls, mouse movement
2. **Session Recordings** - Watch user sessions
3. **A/B Testing** - Test page variations
4. **Funnel Analysis** - Track conversion funnels
5. **Form Analytics** - Identify form abandonment

## Key Features Used

- [ ] Heatmaps (click, scroll, attention)
- [ ] Session recordings
- [ ] A/B testing
- [ ] Funnel tracking
- [ ] Form analytics
- [ ] Polls & surveys
- [ ] Personalization (dynamic content)
- [ ] Goal tracking

## Tracking Setup

### Installation

```html
<!-- PageSense Tracking Code -->
<script>
  (function (w, d, t, r, u) {
    w[u] = w[u] || [];
    w[u].push({ projectId: "10000123456" });
    var s = d.createElement(t);
    s.async = 1;
    s.src = r;
    var h = d.getElementsByTagName(t)[0];
    h.parentNode.insertBefore(s, h);
  })(
    window,
    document,
    "script",
    "https://cdn-in.pagesense.io/js/zeusorg/abc123.js",
    "_ps",
  );
</script>
```

Place before `</head>` tag on all pages.

## Heatmaps

### Types

**Click Heatmap**

- Shows where users click
- Identifies dead clicks (no action)
- Reveals popular CTAs

**Scroll Heatmap**

- Shows how far users scroll
- Identifies fold line
- Helps place important content

**Attention Heatmap**

- Mouse movement tracking
- Shows areas of interest
- Indicates reading patterns

**Mobile Heatmap**

- Touch-specific tracking
- Pinch/zoom tracking

### Example Insights

- 60% of users don't scroll below fold → Move CTA up
- Users clicking non-clickable image → Make it clickable
- High attention on testimonials → Add more social proof

## Session Recordings

### What's Recorded

- Mouse movements
- Clicks and taps
- Scrolling behavior
- Form interactions
- Page navigation
- Time on page

### Filter Recordings By

- **Device**: Desktop, mobile, tablet
- **Browser**: Chrome, Safari, Firefox
- **Location**: Country, city
- **User Type**: New visitor, returning, logged in
- **Actions**: Rage clicks, dead clicks, error occurred

### Privacy Controls

- ✅ Mask sensitive inputs (passwords, credit cards)
- ✅ Exclude pages (e.g., /checkout, /account)
- ✅ Disable for GDPR regions (if needed)
- ✅ IP anonymization

### Example Use Cases

1. User stuck on form → Identify confusing field
2. Cart abandonment → See where hesitation occurs
3. Feature discovery → Watch how users find features
4. Error diagnosis → See steps leading to error

## A/B Testing

### Test Types

**Page Variation**

- Test completely different page designs
- Split traffic 50/50

**Element Test**

- Test specific elements (headline, CTA, image)
- Multiple variations possible

**Multi-Variate Test**

- Test multiple elements simultaneously
- Requires more traffic

### Example Tests

**Homepage CTA Button**

- Variation A: "Get Started Free"
- Variation B: "Start Your Trial"
- Variation C: "Try OneMind Now"
- **Winner**: Variation B (+18% conversion)

**Pricing Page Layout**

- Variation A: 3 columns
- Variation B: Table format
- **Winner**: Variation A (+24% sign-ups)

### Statistical Significance

- Minimum 100 conversions per variation
- 95% confidence level
- Run until significance achieved (typically 2-4 weeks)

## Funnel Analysis

### Example Funnel: Sign-Up Flow

1. **Landing Page** (1000 visitors) → 100%
2. **Pricing Page** (700 visitors) → 70%
3. **Sign-Up Form** (400 visitors) → 40%
4. **Email Verification** (320 visitors) → 32%
5. **Onboarding Complete** (280 visitors) → 28%

**Drop-off Points**:

- 30% leave after landing page (need better CTA)
- 43% leave after pricing (price objections?)
- 20% abandon form (too long?)
- 13% don't verify email (better copy needed)

## Form Analytics

### Tracked Metrics

- **Field Interaction Time** - How long on each field
- **Drop-off Rate** - Which field causes abandonment
- **Error Rate** - Validation errors per field
- **Refills** - How many times field is changed

### Example Insights

- "Company Size" field ignored 40% → Make optional
- "Phone" field causes 30% drop-off → Move to later
- "Password" field 5+ refills → Add strength indicator
- Email validation error 15% → Show format example

## Polls & Surveys

### On-Site Polls

- Pop-up or slide-in
- Single/multiple questions
- Targeted by page, device, user segment

### Example Poll Questions

1. "How would you rate your experience?" (1-5 stars)
2. "What brought you to our site today?" (multiple choice)
3. "Did you find what you were looking for?" (Yes/No)
4. "What feature would you like to see next?" (open text)

### Trigger Conditions

- **Time on Site**: After 30 seconds
- **Exit Intent**: When mouse leaves viewport
- **Scroll Depth**: At 50% scroll
- **Page Count**: After 3 pages viewed

## Personalization

### Dynamic Content

Show different content based on:

- **Location**: Show local testimonials
- **Device**: Optimize for mobile
- **Referral Source**: Different messaging per channel
- **User Behavior**: Returning visitor vs. new
- **Customer Status**: Free vs. paid user

### Example: Homepage Personalization

- **New Visitor**: "Welcome! Start Your Free Trial"
- **Returning Visitor** (not signed up): "Welcome Back! Pick Up Where You Left Off"
- **Trial User**: "Your Trial Expires in 5 Days - Upgrade Now"
- **Paid User**: "Welcome Back, [Name]! Go to Dashboard →"

## Goal Tracking

### Goals

- **URL Visit**: User reaches thank-you page
- **Click**: CTA button clicked
- **Form Submit**: Contact form submitted
- **Custom Event**: JavaScript event fired

### Example Goals

- **Goal 1**: Sign-up completed (`/welcome` page)
- **Goal 2**: Demo requested (form submit)
- **Goal 3**: Video watched (custom event)
- **Goal 4**: Pricing page visited (URL match)

## Integration with Other Zoho Apps

- **Zoho CRM**: Track lead behavior on website
- **Zoho Campaigns**: Trigger emails based on behavior
- **Zoho SalesIQ**: See visitor activity in live chat
- **Zoho Analytics**: Advanced custom reports

## Legacy AI Integration

### Automated Actions

1. **High Engagement** (5+ page views) → Create hot lead in CRM
2. **Stuck on Page** (3+ min) → Trigger live chat offer
3. **Error Occurred** → Log to Sentry for debugging
4. **Goal Achieved** → Update user journey stage
5. **A/B Test Winner** → Notify team in Cliq

## Privacy & Compliance

- ✅ GDPR compliant
- ✅ Cookie consent integration
- ✅ Data anonymization options
- ✅ Opt-out mechanism
- ✅ Data retention controls (30-90 days)

## Analytics Dashboard

### Key Metrics

- **Total Visitors**: Daily/weekly/monthly
- **Bounce Rate**: % who leave without interaction
- **Avg. Session Duration**: Time on site
- **Pages per Session**: Page views per visit
- **Top Pages**: Most visited pages
- **Traffic Sources**: Organic, direct, referral, social

## Notes

- Free tier: 10K page views/month
- Paid plans: $10+/month (higher traffic limits)
- Heatmaps require 50+ interactions for reliability
- Session recordings don't capture sensitive data (PCI compliance)
- A/B tests should run at least 7 days (account for weekly patterns)
- Combine with Google Analytics for complete picture

---

**Related**: [Marketing Automation](./Marketing-Automation.md) | [SalesIQ](./SalesIQ.md) | [Forms](./Forms.md)
