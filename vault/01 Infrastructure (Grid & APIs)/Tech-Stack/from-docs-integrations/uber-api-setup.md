# Uber API Setup Guide

## Overview

This guide walks through setting up Uber API access for OneMind OS personal use. The Uber Rides API allows Legacy AI to request rides, get estimates, and track trips on your behalf.

## Prerequisites

- Active Uber rider account
- Credit card linked to Uber account
- 1Password CLI installed (`op` command)
- Email access for verification

## Step 1: Register Developer Account

1. **Go to Uber Developer Portal**
   - Visit: https://developer.uber.com/
   - Click "Get Started" or "Sign In"

2. **Sign Up with Uber Account**
   - Use your existing Uber rider account
   - If you don't have one, create at https://uber.com first

3. **Verify Email**
   - Check email for verification link
   - Click to confirm developer account

4. **Complete Profile**
   - Name: Zeus Delacruz
   - Company: Personal Use
   - Country: United States
   - Phone: (your number)

## Step 2: Create Application

1. **Navigate to Dashboard**
   - Go to https://developer.uber.com/dashboard
   - Click "Create App"

2. **App Details**

   ```
   App Name: OneMind OS - Personal AI Assistant

   Description:
   Personal AI assistant that automates transportation scheduling
   based on calendar events, location changes, and daily routines.
   For single-user personal use only.

   App Type: Rides API

   Redirect URIs:
   - http://localhost:8000/oauth/uber/callback
   - https://onemind.zeusdelacruz.com/oauth/uber/callback
   ```

3. **Save App**
   - Click "Create"
   - App created in "Sandbox" mode

4. **Get Credentials**
   - Client ID: `xxxxxxxxxxxxxxxx`
   - Client Secret: `yyyyyyyyyyyyyyyy`
   - Save these immediately

## Step 3: Store Credentials Securely

```bash
# Store in 1Password using CLI
op item create \
  --category=login \
  --title="Uber API Credentials" \
  --vault="Development" \
  username="Client ID from dashboard" \
  password="Client Secret from dashboard" \
  url="https://developer.uber.com/dashboard"

# Add custom fields
op item edit "Uber API Credentials" \
  app_name="OneMind OS - Personal AI Assistant" \
  environment="sandbox" \
  redirect_uri="http://localhost:8000/oauth/uber/callback"
```

## Step 4: Test Sandbox Access

```python
# Test sandbox connection
import requests

CLIENT_ID = "your_client_id"
CLIENT_SECRET = "your_client_secret"

# Get sandbox access token (server token for testing)
response = requests.post(
    "https://login.uber.com/oauth/v2/token",
    data={
        "client_id": CLIENT_ID,
        "client_secret": CLIENT_SECRET,
        "grant_type": "client_credentials",
        "scope": "rides.request"
    }
)

token = response.json()["access_token"]
print(f"✅ Sandbox access working! Token: {token[:20]}...")
```

## Step 5: Apply for Production Access

### When to Apply

- After testing sandbox thoroughly
- When ready to use with real rides
- Typically 1-2 weeks before launch

### Production Access Form

1. **Navigate to Production Settings**
   - Dashboard → Your App → Settings
   - Click "Request Production Access"

2. **Fill Out Application**

**Section 1: Use Case**

```
Use Case: Personal Transportation Automation

Detailed Description:
I'm building a personal AI assistant (OneMind OS) that helps me
manage my daily schedule and transportation needs. The assistant
will request Uber rides on my behalf based on:

1. Calendar Events - Auto-suggest rides before meetings
2. Location Changes - Suggest return trips after drop-offs
3. Emergency Situations - Quick ride requests to hospital/urgent care
4. Airport Travel - Coordinate rides with flight schedules

This is for PERSONAL USE ONLY. I am the only user. The assistant
will order rides for me using my Uber account and payment method.
```

**Section 2: Expected Usage**

```
Expected Monthly Volume: 20-50 rides
Number of Users: 1 (myself only)
Geographic Region: United States (primarily San Francisco Bay Area)
```

**Section 3: OAuth Scopes Needed**

```
Required Scopes:
- rides.request (order rides)
- rides.read (track ride status)
- request_receipt (get receipts)
- history (view past rides)
- profile (user info)

Why These Scopes:
I need full ride management capabilities to automate my personal
transportation based on my AI assistant's suggestions.
```

**Section 4: Security & Privacy**

```
Data Storage:
- OAuth tokens stored in encrypted 1Password vault
- No user data shared with third parties
- Ride history stored locally for personal expense tracking only

Security Measures:
- OAuth 2.0 with PKCE
- Tokens refreshed automatically
- API keys never in code/logs
- HTTPS only for all API calls
```

3. **Submit Application**
   - Review all information
   - Click "Submit for Review"
   - Confirmation email sent

## Step 6: Wait for Approval

### Timeline

- **Day 1-2**: Initial review by Uber team
- **Day 3-7**: Possible follow-up questions
- **Day 7-14**: Final approval
- **Total**: 1-2 weeks typical

### What to Expect

1. **Automated Response**
   - "We've received your application"
   - Reference number provided

2. **Possible Follow-up Questions**
   - May ask for clarification on use case
   - Respond within 24 hours
   - Be honest about personal use

3. **Approval Email**
   - "Your app is approved for production"
   - Production credentials activated

### If Denied

- Most denials are for unclear use cases
- Re-apply with more detailed explanation
- Emphasize personal, non-commercial use
- Mention low volume (single user)

## Step 7: Production Configuration

Once approved:

1. **Update Environment**

   ```bash
   # Update 1Password with production status
   op item edit "Uber API Credentials" \
     environment="production" \
     approved_date="2026-01-23"
   ```

2. **Update OneMind OS Config**

   ```bash
   # agents/legacy/.env
   UBER_CLIENT_ID="production_client_id"
   UBER_CLIENT_SECRET="production_client_secret"
   UBER_ENVIRONMENT="production"
   ```

3. **Test Production API**

   ```python
   # First production test
   from core.integrations.uber import UberIntegration

   uber = UberIntegration()

   # Get ride estimate (no charge)
   estimate = uber.get_estimate(
       pickup="current",
       destination="1 Market St, San Francisco"
   )

   print(f"✅ Production API working!")
   print(f"Estimate: ${estimate['price']}, {estimate['duration']} min")
   ```

## Step 8: OAuth Flow Implementation

See `docs/integrations/ride-api-oauth.md` for complete OAuth implementation details.

## Troubleshooting

### Issue: "Invalid Client ID"

**Solution:**

- Verify Client ID copied correctly (no spaces)
- Check if using sandbox vs production credentials
- Regenerate credentials in dashboard

### Issue: "Scope not authorized"

**Solution:**

- Wait for production approval
- Sandbox has limited scopes
- Request additional scopes in dashboard

### Issue: Production access denied

**Solution:**

- Re-read denial reason carefully
- Common issue: unclear use case
- Rewrite application emphasizing personal use
- Mention single user, low volume
- Provide more detail on automation scenarios

## API Rate Limits

**Sandbox:**

- 1,000 requests/hour
- Unlimited fake rides

**Production (Personal Use):**

- 10 requests/minute
- 1,000 requests/hour
- More than enough for personal use

## Cost Summary

- Developer Account: **FREE**
- Sandbox Testing: **FREE**
- Production API Access: **FREE**
- API Requests: **FREE**
- Rides: **Normal Uber pricing**

## Resources

- **Uber Developer Docs**: https://developer.uber.com/docs/riders/ride-requests/introduction
- **OAuth Guide**: https://developer.uber.com/docs/riders/guides/authentication
- **Support**: https://developer.uber.com/support

## Next Steps

After Uber setup is complete:

1. Follow `lyft-api-setup.md` for Lyft integration
2. Implement OAuth flow (see `ride-api-oauth.md`)
3. Test both APIs in sandbox
4. Apply for production access for both
5. Wait for approvals (can run in parallel)

---

**Setup Time:** 1-2 hours  
**Approval Time:** 1-2 weeks  
**Total:** ~2 weeks to production
