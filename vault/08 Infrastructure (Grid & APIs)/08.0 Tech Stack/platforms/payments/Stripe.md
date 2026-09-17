# Stripe

> **Payment Processing Platform**

## Overview

Stripe provides payment processing, subscription management, and financial infrastructure for OneMind OS monetization.

## Status

- **Active**: ✅ Yes
- **Mode**: Live + Test
- **Last Updated**: January 25, 2026

## Use Cases

1. **Subscriptions** - Recurring billing
2. **One-Time Payments** - Product purchases
3. **Invoicing** - Business billing
4. **Connect** - Marketplace payouts
5. **Billing Portal** - Customer self-service

## Products

| Product                | Purpose              | Status    |
| ---------------------- | -------------------- | --------- |
| **OneMind Pro**        | Monthly subscription | ✅ Active |
| **OneMind Enterprise** | Annual subscription  | ✅ Active |
| **Consulting**         | One-time payment     | ✅ Active |

## Installation

```bash
pip install stripe
```

## Configuration

### Environment Variables

```env
STRIPE_SECRET_KEY=sk_live_...
STRIPE_PUBLISHABLE_KEY=pk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...
```

### Client Setup

```python
import stripe

stripe.api_key = os.getenv("STRIPE_SECRET_KEY")
```

## Checkout Sessions

### Create Checkout

```python
def create_checkout(price_id: str, customer_email: str, success_url: str, cancel_url: str):
    session = stripe.checkout.Session.create(
        payment_method_types=['card'],
        line_items=[{
            'price': price_id,
            'quantity': 1,
        }],
        mode='subscription',  # or 'payment'
        success_url=success_url + '?session_id={CHECKOUT_SESSION_ID}',
        cancel_url=cancel_url,
        customer_email=customer_email,
    )
    return session.url
```

### FastAPI Endpoint

```python
from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

class CheckoutRequest(BaseModel):
    price_id: str
    email: str

@app.post("/create-checkout")
async def create_checkout_endpoint(request: CheckoutRequest):
    session = stripe.checkout.Session.create(
        payment_method_types=['card'],
        line_items=[{'price': request.price_id, 'quantity': 1}],
        mode='subscription',
        success_url='https://onemind.io/success?session_id={CHECKOUT_SESSION_ID}',
        cancel_url='https://onemind.io/cancel',
        customer_email=request.email,
    )
    return {"url": session.url}
```

## Subscriptions

### Create Subscription

```python
def create_subscription(customer_id: str, price_id: str):
    subscription = stripe.Subscription.create(
        customer=customer_id,
        items=[{'price': price_id}],
        payment_behavior='default_incomplete',
        expand=['latest_invoice.payment_intent'],
    )
    return subscription
```

### Cancel Subscription

```python
def cancel_subscription(subscription_id: str, at_period_end: bool = True):
    return stripe.Subscription.modify(
        subscription_id,
        cancel_at_period_end=at_period_end
    )
```

### Get Subscription Status

```python
def get_subscription_status(subscription_id: str):
    sub = stripe.Subscription.retrieve(subscription_id)
    return {
        'status': sub.status,
        'current_period_end': sub.current_period_end,
        'cancel_at_period_end': sub.cancel_at_period_end
    }
```

## Customers

### Create Customer

```python
def create_customer(email: str, name: str = None, metadata: dict = None):
    return stripe.Customer.create(
        email=email,
        name=name,
        metadata=metadata or {}
    )
```

### Get Customer

```python
def get_customer(customer_id: str):
    return stripe.Customer.retrieve(customer_id)

def get_customer_by_email(email: str):
    customers = stripe.Customer.list(email=email)
    return customers.data[0] if customers.data else None
```

## Webhooks

### Webhook Handler

```python
from fastapi import Request, HTTPException

@app.post("/webhooks/stripe")
async def stripe_webhook(request: Request):
    payload = await request.body()
    sig_header = request.headers.get('stripe-signature')

    try:
        event = stripe.Webhook.construct_event(
            payload, sig_header, os.getenv("STRIPE_WEBHOOK_SECRET")
        )
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid payload")
    except stripe.error.SignatureVerificationError:
        raise HTTPException(status_code=400, detail="Invalid signature")

    # Handle events
    if event['type'] == 'checkout.session.completed':
        session = event['data']['object']
        await handle_checkout_completed(session)

    elif event['type'] == 'customer.subscription.updated':
        subscription = event['data']['object']
        await handle_subscription_updated(subscription)

    elif event['type'] == 'customer.subscription.deleted':
        subscription = event['data']['object']
        await handle_subscription_cancelled(subscription)

    elif event['type'] == 'invoice.payment_failed':
        invoice = event['data']['object']
        await handle_payment_failed(invoice)

    return {"status": "success"}
```

### Event Handlers

```python
async def handle_checkout_completed(session):
    customer_email = session.get('customer_email')
    subscription_id = session.get('subscription')

    # Update user in database
    await db.users.update_one(
        {'email': customer_email},
        {'$set': {
            'stripe_customer_id': session['customer'],
            'subscription_id': subscription_id,
            'subscription_status': 'active'
        }}
    )

async def handle_subscription_cancelled(subscription):
    await db.users.update_one(
        {'subscription_id': subscription['id']},
        {'$set': {'subscription_status': 'cancelled'}}
    )
```

## Billing Portal

### Create Portal Session

```python
def create_portal_session(customer_id: str, return_url: str):
    session = stripe.billing_portal.Session.create(
        customer=customer_id,
        return_url=return_url,
    )
    return session.url
```

### FastAPI Endpoint

```python
@app.post("/create-portal")
async def create_portal(user_id: str):
    user = await get_user(user_id)

    session = stripe.billing_portal.Session.create(
        customer=user.stripe_customer_id,
        return_url='https://onemind.io/account',
    )

    return {"url": session.url}
```

## Products & Prices

### Create Product

```python
product = stripe.Product.create(
    name="OneMind Pro",
    description="AI-powered personal OS",
)

price = stripe.Price.create(
    product=product.id,
    unit_amount=2900,  # $29.00
    currency="usd",
    recurring={"interval": "month"},
)
```

### List Products

```python
products = stripe.Product.list(active=True)
prices = stripe.Price.list(active=True)
```

## Testing

### Test Cards

| Number             | Description  |
| ------------------ | ------------ |
| `4242424242424242` | Succeeds     |
| `4000000000000002` | Decline      |
| `4000002500003155` | Requires 3DS |

### Test Webhooks

```bash
# Install Stripe CLI
brew install stripe/stripe-cli/stripe

# Forward webhooks to local
stripe listen --forward-to localhost:8000/webhooks/stripe

# Trigger test events
stripe trigger checkout.session.completed
```

## Price IDs

| Product     | Interval | Price ID  |
| ----------- | -------- | --------- |
| OneMind Pro | Monthly  | price_xxx |
| OneMind Pro | Annual   | price_yyy |
| Enterprise  | Annual   | price_zzz |

## Best Practices

- ✅ Always verify webhook signatures
- ✅ Handle all relevant events
- ✅ Use idempotency keys for mutations
- ✅ Store Stripe IDs in your database
- ✅ Test with Stripe CLI locally
- ✅ Handle failed payments gracefully
- ✅ Use test mode during development

## Resources

- [Stripe Documentation](https://stripe.com/docs)
- [Stripe API Reference](https://stripe.com/docs/api)
- [Stripe CLI](https://stripe.com/docs/stripe-cli)

---

**Related**: [Whop](./Whop.md) | [Circle](./Circle.md)
