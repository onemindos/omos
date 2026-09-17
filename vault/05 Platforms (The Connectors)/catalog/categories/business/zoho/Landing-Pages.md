# Zoho Landing Pages (PageSense)

> **Landing Page Builder & Optimization**

## Overview

Zoho Landing Pages (part of PageSense/Marketing Hub) is a drag-and-drop landing page builder for creating conversion-focused pages without coding.

## Status

- **Active**: ✅ Yes
- **Integration**: Embed + API
- **Last Updated**: January 25, 2026

## Use Cases

1. **Product Launches** - New feature announcements
2. **Lead Generation** - Capture emails/contact info
3. **Event Registration** - Workshop/webinar sign-ups
4. **Sales Pages** - Product/service offerings
5. **Thank You Pages** - Post-conversion confirmations

## Key Features Used

- [ ] Drag-and-drop builder
- [ ] Mobile responsive (auto-adapt)
- [ ] A/B testing
- [ ] Form integration (Zoho Forms)
- [ ] Custom domain mapping
- [ ] Analytics & heatmaps
- [ ] Pop-ups & sticky bars
- [ ] Templates library

## Active Landing Pages

| Page Name                 | URL                      | Purpose       | Conversion Rate |
| ------------------------- | ------------------------ | ------------- | --------------- |
| **OneMind Home**          | onemind.zeusdelacruz.com | Main page     | 12%             |
| **Free Trial**            | onemind.com/trial        | Sign-up page  | 28%             |
| **Workshop Registration** | onemind.com/workshop     | Event sign-up | 35%             |
| **Beta Access**           | onemind.com/beta         | Early access  | 22%             |
| **Case Study**            | onemind.com/case-study   | Social proof  | 8%              |

## Page Elements

### Common Components

- **Hero Section**: Headline, subheadline, CTA button, hero image
- **Features Section**: 3-column grid of features with icons
- **Social Proof**: Testimonials, logos, stats
- **Lead Capture Form**: Email, name, company (Zoho Forms)
- **FAQ Section**: Accordion-style questions
- **Footer**: Links, social media, contact info

### Advanced Components

- **Video Embed**: YouTube/Vimeo
- **Countdown Timer**: Create urgency
- **Pricing Table**: Compare plans
- **Image Gallery**: Showcase product
- **Trust Badges**: Security, certifications

## Template Structure

### Standard Landing Page Template

```
┌─────────────────────────────────┐
│  Navigation Bar (Logo + CTA)    │
├─────────────────────────────────┤
│  Hero Section                   │
│  - Headline                     │
│  - Subheadline                  │
│  - CTA Button                   │
│  - Hero Image                   │
├─────────────────────────────────┤
│  Social Proof (Logos)           │
├─────────────────────────────────┤
│  Features (3 columns)           │
├─────────────────────────────────┤
│  Testimonials                   │
├─────────────────────────────────┤
│  Lead Capture Form              │
├─────────────────────────────────┤
│  FAQ                            │
├─────────────────────────────────┤
│  Footer                         │
└─────────────────────────────────┘
```

## Custom Domain Setup

### DNS Configuration

```
CNAME Record:
  Name: onemind (or www)
  Value: pages.zoho.com
  TTL: 3600
```

### SSL Certificate

- Automatically provisioned by Zoho
- Free Let's Encrypt certificate
- Auto-renewal every 90 days

## A/B Testing

### Example Test: Hero CTA Button

**Variation A (Control)**

- Button Text: "Start Free Trial"
- Button Color: Blue (#667eea)
- Conversion: 25%

**Variation B**

- Button Text: "Get Started Free"
- Button Color: Green (#10b981)
- Conversion: 28%

**Winner**: Variation B (+12% improvement)

### Test Ideas

- Headline variations
- CTA button text/color
- Hero image vs. hero video
- Form length (3 fields vs. 5 fields)
- Social proof placement
- Pricing display (show vs. hide)

## Form Integration

### Connected Forms (Zoho Forms)

- Contact form
- Trial sign-up form
- Demo request form
- Newsletter subscription
- Event registration form

### Form Actions on Submit

1. Send confirmation email (auto-responder)
2. Create CRM lead
3. Add to email campaign (Zoho Campaigns)
4. Redirect to thank you page
5. Trigger webhook to Legacy AI

## Analytics

### Tracked Metrics

- **Page Views**: Total visitors
- **Unique Visitors**: Distinct visitors
- **Conversion Rate**: % who complete form
- **Bounce Rate**: % who leave without interaction
- **Time on Page**: Average session duration
- **Traffic Sources**: Organic, direct, referral, social

### Heatmaps (PageSense Integration)

- Click heatmap: Where users click
- Scroll heatmap: How far users scroll
- Attention heatmap: Areas of interest

## Pop-ups & Overlays

### Pop-up Types

**Exit Intent Pop-up**

- Trigger: Mouse moves to close tab
- Message: "Wait! Get 10% off your first month"
- CTA: "Claim Discount"

**Time-Based Pop-up**

- Trigger: After 30 seconds on page
- Message: "Join 500+ users already using OneMind"
- CTA: "Start Free Trial"

**Scroll-Based Pop-up**

- Trigger: Scrolled 50% of page
- Message: "Want to learn more?"
- CTA: "Download Free Guide"

### Sticky Bar

- Fixed to top or bottom of page
- Message: "Early Bird Special: 20% Off Ends Tonight!"
- Always visible while scrolling

## Mobile Optimization

### Responsive Features

- Auto-adjust layout for mobile
- Touch-friendly buttons (min 44x44px)
- Simplified forms (less fields)
- Optimized images (WebP format)
- Faster load times (< 3 seconds)

### Mobile-Specific Changes

- Stack columns vertically
- Larger fonts (min 16px)
- Hamburger menu navigation
- Click-to-call buttons
- Shorter headlines

## SEO Optimization

### On-Page SEO

- Custom meta title (50-60 characters)
- Meta description (150-160 characters)
- H1 tag (one per page)
- Alt text for images
- Clean URL structure
- Internal linking
- Schema markup (JSON-LD)

### Example Meta Tags

```html
<title>OneMind OS - AI-Powered Personal Operating System</title>
<meta
  name="description"
  content="OneMind OS combines human intelligence with AI to create a seamless personal operating system. Start your free trial today."
/>
<meta property="og:title" content="OneMind OS - Your AI-Powered Second Brain" />
<meta property="og:image" content="https://onemind.com/og-image.jpg" />
```

## Integrations

### Connected Apps

- **Zoho Forms**: Lead capture forms
- **Zoho CRM**: Auto-create leads
- **Zoho Campaigns**: Add to email list
- **Zoho Analytics**: Advanced reporting
- **Zoho PageSense**: Heatmaps & session recordings
- **Google Analytics**: Traffic tracking
- **Facebook Pixel**: Retargeting ads

### Webhook Integration

```json
{
  "event": "form_submission",
  "url": "https://onemind.zeusdelacruz.com/webhooks/landing-page",
  "data": {
    "name": "{{name}}",
    "email": "{{email}}",
    "company": "{{company}}",
    "page": "{{page_url}}"
  }
}
```

## Page Speed Optimization

### Best Practices

- ✅ Image optimization (compress, lazy load)
- ✅ Minify CSS/JavaScript
- ✅ Enable browser caching
- ✅ Use CDN for assets
- ✅ Minimize HTTP requests
- ✅ Defer non-critical JavaScript

### Target Metrics

- **Load Time**: < 3 seconds
- **First Contentful Paint**: < 1.5 seconds
- **Time to Interactive**: < 3.5 seconds
- **Lighthouse Score**: > 90

## Conversion Rate Optimization (CRO)

### CRO Checklist

- [ ] Clear, benefit-driven headline
- [ ] Single, prominent CTA
- [ ] Minimal form fields (3-5 max)
- [ ] Social proof (testimonials, logos, stats)
- [ ] Trust indicators (security badges, guarantees)
- [ ] Mobile-friendly design
- [ ] Fast load time (< 3 sec)
- [ ] A/B test key elements
- [ ] Remove navigation distractions
- [ ] Use contrasting CTA button color

### Conversion Funnel

1. **Ad/Email** (100%) → Click
2. **Landing Page** (50%) → View
3. **Form Started** (30%) → Begin
4. **Form Completed** (25%) → Submit
5. **Thank You Page** (25%) → Converted

## Legacy AI Integration

### Automated Actions

1. **Form Submit** → Create CRM lead, send to Todoist
2. **High Engagement** (5+ min on page) → Trigger live chat
3. **Exit Without Converting** → Add to retargeting campaign
4. **A/B Test Winner** → Notify team in Cliq
5. **Low Conversion Rate** (<10%) → Alert for page review

## Templates Used

### Free Trial Page Template

- **Headline**: "Start Your 14-Day Free Trial"
- **Features**: 3 key benefits
- **Form**: Email + Name only
- **Social Proof**: "Join 500+ users"
- **Guarantee**: "No credit card required"

### Workshop Registration Template

- **Headline**: "Join Our Live Workshop"
- **Details**: Date, time, topics, speaker
- **Form**: Name + Email + Company
- **Countdown**: "Only 50 spots left!"
- **Bonus**: "Get free workbook PDF"

## Notes

- Part of Zoho PageSense or Marketing Automation
- Unlimited landing pages on paid plans
- A/B testing requires Professional+ plan
- Custom domain included on all plans
- Free SSL certificates
- Mobile-responsive by default
- Template library with 50+ pre-built pages

---

**Related**: [PageSense](PageSense.md) | [Forms](Forms.md) | [Campaigns](Campaigns.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
