# Whop

> **Digital Products Marketplace**

## Overview

Whop is a digital products marketplace for selling memberships, software access, and digital products. It powers the OneMind product distribution.

## Status

- **Active**: ✅ Yes
- **Products**: OneMind Pro, Templates
- **Last Updated**: January 25, 2026

## Use Cases

1. **Memberships** - Recurring subscriptions
2. **Digital Products** - One-time purchases
3. **Software Access** - License keys
4. **Affiliates** - Partner program
5. **Community Integration** - Discord/Circle access

## Key Features

- Built-in checkout
- Subscription management
- Affiliate system
- Discord/community integration
- Analytics dashboard
- Webhook integrations
- License key generation

## Products

| Product         | Type         | Price   | Status    |
| --------------- | ------------ | ------- | --------- |
| OneMind Pro     | Subscription | $29/mo  | ✅ Active |
| OneMind Annual  | Subscription | $249/yr | ✅ Active |
| Templates Pack  | One-time     | $49     | ✅ Active |
| Consulting Hour | One-time     | $250    | ✅ Active |

## API Integration

### Configuration

```python
import requests

WHOP_API_KEY = os.getenv("WHOP_API_KEY")
WHOP_COMPANY_ID = os.getenv("WHOP_COMPANY_ID")

headers = {
    "Authorization": f"Bearer {WHOP_API_KEY}",
    "Content-Type": "application/json"
}

base_url = "https://api.whop.com/api/v2"
```

### Get Memberships

```python
def get_memberships(user_id: str = None):
    params = {}
    if user_id:
        params["user_id"] = user_id

    response = requests.get(
        f"{base_url}/memberships",
        headers=headers,
        params=params
    )
    return response.json()
```

### Validate Membership

```python
def validate_membership(license_key: str) -> dict:
    response = requests.post(
        f"{base_url}/memberships/{license_key}/validate",
        headers=headers
    )
    return response.json()

def check_access(license_key: str) -> bool:
    result = validate_membership(license_key)
    return result.get("valid", False) and result.get("status") == "active"
```

### Get User

```python
def get_user(user_id: str):
    response = requests.get(
        f"{base_url}/users/{user_id}",
        headers=headers
    )
    return response.json()
```

## Webhooks

### Webhook Events

| Event                     | Description          |
| ------------------------- | -------------------- |
| `membership.went_valid`   | New purchase/renewal |
| `membership.went_invalid` | Cancelled/expired    |
| `payment.succeeded`       | Payment successful   |
| `payment.failed`          | Payment failed       |

### Webhook Handler

```python
from fastapi import FastAPI, Request
import hmac
import hashlib

WHOP_WEBHOOK_SECRET = os.getenv("WHOP_WEBHOOK_SECRET")

def verify_signature(payload: bytes, signature: str) -> bool:
    expected = hmac.new(
        WHOP_WEBHOOK_SECRET.encode(),
        payload,
        hashlib.sha256
    ).hexdigest()
    return hmac.compare_digest(signature, expected)

@app.post("/webhooks/whop")
async def whop_webhook(request: Request):
    payload = await request.body()
    signature = request.headers.get("whop-signature", "")

    if not verify_signature(payload, signature):
        raise HTTPException(status_code=401)

    data = await request.json()
    action = data.get("action")
    membership = data.get("data")

    if action == "membership.went_valid":
        await handle_new_membership(membership)
    elif action == "membership.went_invalid":
        await handle_cancelled_membership(membership)

    return {"status": "ok"}
```

### Event Handlers

```python
async def handle_new_membership(membership: dict):
    user_email = membership.get("user", {}).get("email")
    product_id = membership.get("plan", {}).get("product_id")
    license_key = membership.get("license_key")

    # Create user in system
    await create_user(email=user_email, license_key=license_key)

    # Add to Circle community
    add_member_to_circle(user_email, get_space_groups(product_id))

    # Send welcome email
    await send_welcome_email(user_email)

async def handle_cancelled_membership(membership: dict):
    user_email = membership.get("user", {}).get("email")

    # Downgrade access
    await downgrade_user(user_email)

    # Send retention email
    await send_cancellation_email(user_email)
```

## License Key Validation

### In OneMind App

```python
async def authenticate_user(license_key: str):
    # Validate with Whop
    result = validate_membership(license_key)

    if not result.get("valid"):
        raise AuthenticationError("Invalid license key")

    if result.get("status") != "active":
        raise AuthenticationError("Subscription not active")

    return {
        "user_id": result.get("user", {}).get("id"),
        "email": result.get("user", {}).get("email"),
        "plan": result.get("plan", {}).get("name"),
        "expires_at": result.get("renewal_period_end")
    }
```

### Middleware

```python
from fastapi import Depends, HTTPException, Header

async def validate_whop_license(authorization: str = Header(...)):
    if not authorization.startswith("License "):
        raise HTTPException(status_code=401, detail="Invalid auth format")

    license_key = authorization.replace("License ", "")

    if not check_access(license_key):
        raise HTTPException(status_code=403, detail="Invalid or expired license")

    return license_key

@app.get("/api/protected")
async def protected_endpoint(license: str = Depends(validate_whop_license)):
    return {"message": "Access granted"}
```

## Checkout Links

### Generate Checkout URL

```python
def create_checkout_link(product_id: str, redirect_url: str = None):
    params = {}
    if redirect_url:
        params["redirect"] = redirect_url

    return f"https://whop.com/checkout/{product_id}?" + urlencode(params)
```

### Embed Checkout

```html
<!-- Whop checkout button -->
<a
  href="https://whop.com/checkout/plan_xxx"
  class="whop-checkout-button"
  target="_blank"
>
  Get OneMind Pro
</a>
```

## Affiliate Program

### Get Affiliates

```python
def get_affiliates():
    response = requests.get(
        f"{base_url}/affiliates",
        headers=headers
    )
    return response.json()
```

### Affiliate Stats

```python
def get_affiliate_stats(affiliate_id: str):
    response = requests.get(
        f"{base_url}/affiliates/{affiliate_id}/stats",
        headers=headers
    )
    return response.json()
```

## Legacy AI Integration

### Check User Status

```python
@Tool
def check_subscription_status(email: str) -> dict:
    """Check a user's Whop subscription status."""
    memberships = get_memberships()
    user_membership = next(
        (m for m in memberships.get("data", [])
         if m.get("user", {}).get("email") == email),
        None
    )

    if not user_membership:
        return {"status": "no_subscription", "email": email}

    return {
        "status": user_membership.get("status"),
        "plan": user_membership.get("plan", {}).get("name"),
        "expires": user_membership.get("renewal_period_end"),
        "email": email
    }
```

## Analytics

### Key Metrics

- Monthly Recurring Revenue (MRR)
- New subscriptions
- Churn rate
- Affiliate revenue
- Conversion rate

### Dashboard

- `https://dash.whop.com/company/{company_id}/analytics`

## vs Stripe

| Feature      | Whop        | Stripe     |
| ------------ | ----------- | ---------- |
| Checkout     | Built-in    | Self-built |
| Affiliates   | Built-in    | External   |
| License keys | Built-in    | Custom     |
| Community    | Built-in    | External   |
| Discovery    | Marketplace | No         |
| Fees         | Higher      | Lower      |

## Best Practices

- ✅ Validate licenses on every request
- ✅ Cache validation results (short TTL)
- ✅ Handle webhook retries
- ✅ Set up affiliate program early
- ✅ Use Whop branding for trust
- ✅ Provide clear upgrade paths

## Resources

- [Whop API Documentation](https://dev.whop.com/)
- [Webhook Events](https://dev.whop.com/webhooks)
- [Dashboard](https://dash.whop.com)

---

**Related**: [Stripe](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/Tech-Stack/platforms/payments/Stripe.md) | [Circle](./Circle.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
