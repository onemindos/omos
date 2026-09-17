# Zoho Payments

> **Integrated Payment Processing Across Zoho Apps**

## Overview

| Attribute   | Value                     |
| ----------- | ------------------------- |
| Type        | Business Payment Platform |
| Provider    | Zoho Corporation          |
| Integration | Native across Zoho One    |
| Status      | ✅ Active                 |

## What is Zoho Payments?

Zoho Payments is Zoho's native payment processing solution deeply integrated across the Zoho ecosystem. It allows businesses to accept payments directly through Zoho apps like Invoice, Books, Commerce, and Checkout.

## Payment Capabilities

### Payment Methods Supported

| Method             | Availability                        |
| ------------------ | ----------------------------------- |
| Credit/Debit Cards | ✅ Visa, Mastercard, Amex, Discover |
| ACH Bank Transfers | ✅ US bank accounts                 |
| PayPal             | ✅ PayPal integration               |
| Apple Pay          | ✅ Where supported                  |
| Google Pay         | ✅ Where supported                  |
| Stripe Connect     | ✅ As alternate processor           |

### Payment Gateways Integrated

| Gateway                | Use Case              |
| ---------------------- | --------------------- |
| Zoho Payments (Native) | Default for Zoho apps |
| Stripe                 | Alternative processor |
| PayPal                 | Consumer payments     |
| Authorize.Net          | Enterprise payments   |
| Square                 | POS integration       |

## Zoho Apps with Payments

### Zoho Invoice

| Feature            | Description                  |
| ------------------ | ---------------------------- |
| Online Payments    | Accept payments on invoices  |
| Payment Links      | Shareable payment URLs       |
| Recurring Invoices | Automatic billing cycles     |
| Partial Payments   | Accept deposits/installments |
| Multi-Currency     | Bill in customer's currency  |

### Zoho Books

| Feature             | Description                    |
| ------------------- | ------------------------------ |
| Integrated Payments | Payments sync to accounting    |
| Bank Reconciliation | Auto-match payments            |
| Tax Handling        | Automatic tax calculations     |
| Multi-Entity        | Handle multiple business units |

### Zoho Commerce

| Feature               | Description                 |
| --------------------- | --------------------------- |
| E-commerce Checkout   | Full shopping cart payments |
| Subscription Commerce | Recurring product sales     |
| Digital Products      | Instant delivery on payment |
| Inventory Updates     | Stock adjusts on purchase   |

### Zoho Checkout

| Feature              | Description                 |
| -------------------- | --------------------------- |
| Custom Payment Pages | Branded checkout experience |
| Payment Links        | Share anywhere              |
| Embed on Website     | Checkout widget             |
| QR Code Payments     | Scan to pay                 |

## Zoho Payment Links

### Create Payment Links

```python
from zoho_checkout import ZohoCheckout

checkout = ZohoCheckout()

# Create a payment link
payment_link = checkout.create_link(
    amount=99.00,
    currency="USD",
    description="OneMind Pro Subscription",
    customer_email="customer@example.com",
    success_url="https://onemind.ai/success",
    cancel_url="https://onemind.ai/cancel"
)

print(payment_link.url)  # Share this URL
```

### Recurring Payments

```python
# Create subscription plan
plan = checkout.create_plan(
    name="OneMind Pro Monthly",
    amount=49.00,
    interval="month",
    interval_count=1
)

# Subscribe customer
subscription = checkout.subscribe(
    customer_id="cust_123",
    plan_id=plan.id,
    payment_method_id="pm_456"
)
```

## Fee Structure

| Transaction Type  | Typical Fee       |
| ----------------- | ----------------- |
| Credit/Debit Card | ~2.9% + $0.30     |
| ACH Transfer      | ~0.8% (capped)    |
| International     | Additional 1-1.5% |

## Integration with OneMind

### Current Usage

| App          | Payment Use      |
| ------------ | ---------------- |
| Zoho Invoice | Client billing   |
| Zoho Books   | Revenue tracking |
| Zoho CRM     | Quote-to-cash    |
| Zoho Desk    | Service payments |

### Automation Workflows

```yaml
trigger: invoice_payment_received
actions:
  - update_crm_deal:
      status: "Closed Won"
  - send_welcome_email:
      template: "onboarding"
  - create_project:
      from: invoice_details
  - notify_team:
      channel: "sales"
```

## Security & Compliance

| Standard        | Status       |
| --------------- | ------------ |
| PCI DSS Level 1 | ✅ Compliant |
| SOC 2 Type II   | ✅ Certified |
| GDPR            | ✅ Compliant |
| 3D Secure       | ✅ Supported |

## Best Practices

1. **Enable Auto-Pay** - Reduce missed payments
2. **Send Reminders** - Automate payment reminders
3. **Offer Multiple Methods** - Card, ACH, PayPal
4. **Use Payment Links** - Simplify checkout
5. **Sync with Books** - Automatic reconciliation

## Resources

- [Zoho Payments](https://www.zoho.com/payments/)
- [Zoho Checkout](https://www.zoho.com/checkout/)
- [Payment Gateway Setup](https://www.zoho.com/books/help/payment-gateways.html)
- [Zoho Invoice Payments](https://www.zoho.com/invoice/online-payments.html)

---

Last Updated: January 25, 2026


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
