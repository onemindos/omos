# Composio Setup Checklist

> Step-by-step execution plan for integrating Composio with OneMind OS

## Pre-Flight Checks

- [ ] Python 3.10+ installed
- [ ] OneMind OS repository cloned and up to date
- [ ] NATS server running (for consciousness layer)
- [ ] Redis running (for awareness state)

---

## Phase 1: Composio Account Setup

### 1.1 Create Composio Account
- [ ] Go to [composio.dev](https://composio.dev)
- [ ] Sign up for free account
- [ ] Verify email

### 1.2 Get API Credentials
- [ ] Navigate to Settings → API Keys
- [ ] Create new API key
- [ ] Copy API key (save securely!)
- [ ] Note your Entity ID (usually "default")

### 1.3 Set Up Webhook Secret (Optional but Recommended)
- [ ] Navigate to Settings → Webhooks
- [ ] Generate webhook secret
- [ ] Copy webhook secret

---

## Phase 2: Install Dependencies

### 2.1 Install Composio SDK
```bash
cd /path/to/onemind-os

# Install Composio
pip install composio

# Or add to requirements.txt
echo "composio>=0.5.0" >> requirements.txt
pip install -r requirements.txt
```

### 2.2 Verify Installation
```python
python -c "from composio import ComposioToolSet; print('Composio installed!')"
```

---

## Phase 3: Configure Environment

### 3.1 Add Environment Variables

Add to your `.env` file or environment:

```bash
# Composio Configuration
COMPOSIO_API_KEY=your_api_key_here
COMPOSIO_WEBHOOK_SECRET=your_webhook_secret_here
COMPOSIO_ENTITY_ID=default

# Webhook URL (your gateway's public URL)
COMPOSIO_WEBHOOK_URL=https://your-domain.com/webhooks/composio
```

### 3.2 For Local Development
```bash
# .env.local
COMPOSIO_API_KEY=your_api_key_here
COMPOSIO_WEBHOOK_SECRET=optional_for_local
COMPOSIO_ENTITY_ID=default
COMPOSIO_WEBHOOK_URL=http://localhost:8000/webhooks/composio
```

### 3.3 For Production (Azure VM)
```bash
# SSH into server
ssh zeus@20.121.38.186

# Add to environment
export COMPOSIO_API_KEY=your_api_key_here
export COMPOSIO_WEBHOOK_SECRET=your_webhook_secret_here
export COMPOSIO_WEBHOOK_URL=https://onemindos.dev/webhooks/composio
```

---

## Phase 4: Connect External Services

### 4.1 Connect via Composio Dashboard

For each service, go to Composio Dashboard → Integrations:

#### Todoist
- [ ] Click "Connect" on Todoist
- [ ] Authorize OAuth
- [ ] Verify connection status: "Connected"

#### GitHub
- [ ] Click "Connect" on GitHub
- [ ] Authorize OAuth (select repositories to grant access)
- [ ] Verify connection status: "Connected"

#### Google Calendar
- [ ] Click "Connect" on Google Calendar
- [ ] Authorize OAuth
- [ ] Verify connection status: "Connected"

#### Slack (Optional)
- [ ] Click "Connect" on Slack
- [ ] Authorize OAuth for your workspace
- [ ] Verify connection status: "Connected"

#### Discord (Optional)
- [ ] Click "Connect" on Discord
- [ ] Authorize OAuth
- [ ] Verify connection status: "Connected"

### 4.2 Verify Connections via CLI
```bash
# Install Composio CLI
pip install composio-cli

# Login
composio login

# List connections
composio connections list

# Should show connected apps
```

---

## Phase 5: Configure Triggers (Inbound Events)

### 5.1 Set Up Triggers in Composio Dashboard

Navigate to Triggers section for each connected app:

#### Todoist Triggers
- [ ] Enable `TODOIST_TASK_CREATED`
- [ ] Enable `TODOIST_TASK_COMPLETED`
- [ ] Set webhook URL: `https://your-domain/webhooks/composio`

#### GitHub Triggers
- [ ] Enable `GITHUB_PULL_REQUEST_EVENT`
- [ ] Enable `GITHUB_ISSUE_EVENT`
- [ ] Enable `GITHUB_PUSH_EVENT`
- [ ] Set webhook URL: `https://your-domain/webhooks/composio`

#### Calendar Triggers
- [ ] Enable `GOOGLECALENDAR_EVENT_CREATED`
- [ ] Enable `GOOGLECALENDAR_EVENT_STARTING`
- [ ] Set webhook URL: `https://your-domain/webhooks/composio`

### 5.2 Or Use Python Script
```python
# scripts/setup_composio_triggers.py
import asyncio
from agents.legacy.consciousness import initialize_composio, setup_recommended_triggers

async def main():
    # Initialize bridge
    success = await initialize_composio()
    if not success:
        print("Failed to initialize Composio")
        return

    # Setup triggers for specific apps
    results = await setup_recommended_triggers(["todoist", "github", "calendar"])

    for trigger, success in results.items():
        status = "✓" if success else "✗"
        print(f"{status} {trigger}")

asyncio.run(main())
```

---

## Phase 6: Test the Integration

### 6.1 Test Outbound (Actions)

```python
# test_composio_actions.py
import asyncio
from agents.legacy.tools.composio_tools import composio_todoist, composio_github

async def test_actions():
    # Test Todoist
    result = composio_todoist(
        action="list"
    )
    print(f"Todoist list: {result}")

    # Test GitHub (use your repo)
    result = composio_github(
        action="list_prs",
        repo="your-username/onemind-os"
    )
    print(f"GitHub PRs: {result}")

asyncio.run(test_actions())
```

### 6.2 Test Inbound (Triggers)

#### Start the Gateway
```bash
# Terminal 1: Start gateway
cd agents/legacy/consciousness
python gateway.py

# Should see: "Connected to NATS - Gateway ready!"
```

#### Send Test Webhook
```bash
# Terminal 2: Send test webhook
curl -X POST http://localhost:8000/webhooks/composio \
  -H "Content-Type: application/json" \
  -d '{
    "trigger_name": "TODOIST_TASK_CREATED",
    "trigger_id": "test_123",
    "payload": {
      "content": "Test task from Composio",
      "priority": 2
    },
    "metadata": {
      "timestamp": "2024-01-15T10:00:00Z"
    }
  }'

# Should return: {"status": "ok", "event_id": 1, ...}
```

#### Verify in NATS
```bash
# Check NATS stream
nats stream info ONEMIND
nats consumer info ONEMIND legacy-agent
```

### 6.3 End-to-End Test

1. Create a task in Todoist
2. Watch gateway logs for incoming trigger
3. Check NATS for event
4. Verify Legacy receives it in `legacy.inbox`

---

## Phase 7: Deploy to Production

### 7.1 Update Docker Configuration

Add to `docker-compose.yml`:
```yaml
services:
  gateway:
    environment:
      - COMPOSIO_API_KEY=${COMPOSIO_API_KEY}
      - COMPOSIO_WEBHOOK_SECRET=${COMPOSIO_WEBHOOK_SECRET}
      - COMPOSIO_WEBHOOK_URL=https://onemindos.dev/webhooks/composio
```

### 7.2 Expose Webhook Endpoint

Ensure your gateway is accessible at:
```
https://onemindos.dev/webhooks/composio
```

Via Cloudflare Tunnel or direct exposure.

### 7.3 Update Composio Dashboard

- [ ] Update all trigger webhook URLs to production URL
- [ ] Test each trigger with production endpoint

---

## Phase 8: Monitor & Maintain

### 8.1 Check Composio Dashboard
- [ ] Monitor trigger delivery status
- [ ] Check for failed deliveries
- [ ] Review action execution logs

### 8.2 Check Gateway Logs
```bash
# View gateway logs
docker logs onemind-gateway -f

# Look for:
# - "Composio trigger received: ..."
# - "Composio event published: ..."
```

### 8.3 Monitor NATS Stream
```bash
# Stream statistics
nats stream report

# Recent events
curl http://localhost:8000/api/stream/recent?count=10
```

---

## Troubleshooting

### Triggers Not Arriving

1. **Check webhook URL is correct**
   ```bash
   curl -X POST https://your-domain/webhooks/composio \
     -H "Content-Type: application/json" \
     -d '{"trigger_name": "TEST"}'
   ```

2. **Check Composio Dashboard**
   - Is the trigger enabled?
   - Is the webhook URL correct?
   - Any failed deliveries?

3. **Check Gateway Logs**
   - Is the gateway receiving requests?
   - Any signature verification failures?

### Actions Failing

1. **Check API Key**
   ```python
   import os
   print(os.environ.get("COMPOSIO_API_KEY", "NOT SET"))
   ```

2. **Check Connection Status**
   ```python
   from agents.legacy.consciousness import get_composio_bridge
   bridge = get_composio_bridge()
   print(bridge.get_status())
   ```

3. **Re-authenticate App**
   - Go to Composio Dashboard
   - Disconnect and reconnect the app

### NATS Connection Issues

1. **Check NATS is running**
   ```bash
   docker ps | grep nats
   ```

2. **Check connection string**
   ```bash
   echo $NATS_URL  # Should be nats://localhost:4222 or nats://nats:4222
   ```

---

## Quick Reference

### Environment Variables
| Variable | Required | Description |
|----------|----------|-------------|
| `COMPOSIO_API_KEY` | Yes | Your Composio API key |
| `COMPOSIO_WEBHOOK_SECRET` | No | For webhook signature verification |
| `COMPOSIO_ENTITY_ID` | No | Default: "default" |
| `COMPOSIO_WEBHOOK_URL` | No | Your gateway's public URL |

### Webhook Endpoint
```
POST /webhooks/composio
```

### Important Files
| File | Purpose |
|------|---------|
| `agents/legacy/tools/composio_tools.py` | Outbound action tools |
| `agents/legacy/consciousness/gateway.py` | Inbound webhook handler |
| `agents/legacy/consciousness/composio_bridge.py` | Configuration & management |
| `docs/COMPOSIO_INTEGRATION.md` | Full documentation |

### Useful Commands
```bash
# Install Composio
pip install composio

# Test gateway locally
python agents/legacy/consciousness/gateway.py

# Send test webhook
curl -X POST http://localhost:8000/webhooks/composio -d '...'

# Check stream
curl http://localhost:8000/api/stream/recent
```

---

## Next Steps After Setup

1. [ ] Test with real Todoist task creation
2. [ ] Test with real GitHub PR
3. [ ] Verify Legacy responds to events appropriately
4. [ ] Add more services as needed (Notion, Linear, Gmail, etc.)
5. [ ] Consider setting up Composio MCP server for advanced use cases
