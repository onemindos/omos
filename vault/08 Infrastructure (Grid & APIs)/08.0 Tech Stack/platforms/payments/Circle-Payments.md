# Circle.so Payments

> **Stripe-Powered Community Payments**

## Overview

| Attribute | Value                          |
| --------- | ------------------------------ |
| Type      | Community Payment Integration  |
| Provider  | Circle.so (via Stripe Connect) |
| Currency  | Multi-currency support         |
| Status    | ✅ Active                      |

## How Circle Payments Work

Circle.so uses **Stripe Connect** as its payment backbone. When members purchase memberships, courses, or community access, payments flow through Stripe and are deposited to your connected Stripe account.

## Payment Features

### Membership Billing

| Feature               | Description                         |
| --------------------- | ----------------------------------- |
| Monthly Subscriptions | Recurring monthly payments          |
| Annual Subscriptions  | Yearly billing with discount option |
| One-Time Payments     | Single purchase for lifetime access |
| Free Trials           | Trial periods before billing        |
| Proration             | Automatic proration on plan changes |

### Pricing Options

| Model             | Use Case                     |
| ----------------- | ---------------------------- |
| Free Community    | Open access, no payment      |
| Paid Monthly      | Recurring subscription model |
| Paid Annually     | Discounted annual commitment |
| One-Time Purchase | Lifetime access payment      |
| Tiered Access     | Multiple membership levels   |

### Course Sales

| Feature                 | Description                |
| ----------------------- | -------------------------- |
| Course Pricing          | Set prices per course      |
| Bundles                 | Package multiple courses   |
| Drip Content            | Unlock content over time   |
| Completion Certificates | Issue on course completion |

## Stripe Connect Integration

### Setup Flow

```
1. Connect Stripe account in Circle settings
2. Configure pricing tiers
3. Set up payment pages
4. Enable checkout
5. Payments flow to your Stripe
```

### Revenue Split

| Party   | Typical Split                   |
| ------- | ------------------------------- |
| Creator | Majority of revenue             |
| Stripe  | ~2.9% + $0.30 per transaction   |
| Circle  | Included in Circle subscription |

## Payment Pages

### Customization

- Custom checkout pages
- Brand colors and logo
- Custom thank you messages
- Email confirmation sequences
- Invite-only payment links

### Member Experience

```
1. Member clicks join/purchase
2. Redirected to Stripe Checkout
3. Enters payment details
4. Automatic account creation
5. Access granted immediately
```

## Automation Integration

### Triggered by Payments

| Trigger               | Action                         |
| --------------------- | ------------------------------ |
| New Subscription      | Welcome email, onboarding flow |
| Payment Failed        | Dunning emails, access pause   |
| Subscription Canceled | Exit survey, retention offer   |
| Trial Ending          | Reminder emails                |
| Annual Renewal        | Thank you, loyalty perks       |

### Webhook Events

```json
{
  "event": "member.subscription.created",
  "data": {
    "member_id": "123",
    "plan": "pro_monthly",
    "amount": 49.0,
    "currency": "usd"
  }
}
```

## OneMind Community Setup

### Current Configuration

| Setting           | Value                 |
| ----------------- | --------------------- |
| Community Plan    | Pro                   |
| Payment Processor | Stripe Connect        |
| Tiers             | Free, Pro, Enterprise |
| Currency          | USD                   |

### Member Tiers

| Tier       | Price  | Access                |
| ---------- | ------ | --------------------- |
| Free       | $0/mo  | Basic community       |
| Pro        | $XX/mo | Full access + courses |
| Enterprise | Custom | Custom integrations   |

## Best Practices

1. **Offer Annual Discount** - 2 months free encourages commitment
2. **Use Trials** - 7-14 day trials reduce friction
3. **Automate Dunning** - Recover failed payments
4. **Segment by Tier** - Different content per level
5. **Track Churn** - Monitor cancellation reasons

## Resources

- [Circle Payments Docs](https://community.circle.so/c/help-center)
- [Stripe Connect](https://stripe.com/connect)
- [Circle Pricing](https://circle.so/pricing)

---

Last Updated: January 25, 2026
