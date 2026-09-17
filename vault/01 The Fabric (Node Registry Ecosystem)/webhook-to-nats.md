---
title: "Webhook to NATS Bridge"
type: architecture
status: active
created: 2026-05-29
tags: [webhooks, nats, cloudflare-workers, saas, events]
---
# Webhook → NATS Bridge

## Architecture

```
SaaS Platform (Stripe, GitHub, Linear, etc.)
  → sends POST to https://webhooks.onemindos.dev/<service>/<event>
    → Cloudflare Worker (free tier, global edge)
      → validates signature (per-service)
      → publishes to NATS via WebSocket
        → subject: events.webhook.<service>.<event>
          → NATS subscribers react (rules KV)
```

## Why Cloudflare Workers

- **Free tier**: 100K requests/day (more than enough)
- **Global edge**: <50ms response to any SaaS platform (prevents timeout retries)
- **Zero infrastructure**: No server to maintain, auto-scales
- **WebSocket to NATS**: NATS supports WebSocket connections natively on port :8222
- **Signature validation**: Each SaaS has its own HMAC/signature — validate at the edge
- **Dead simple**: One Worker handles ALL webhooks via URL path routing

## NATS WebSocket Setup

Enable WebSocket in NATS config (`/opt/homebrew/etc/nats-server.conf`):

```hcl
websocket {
  port: 8223
  no_tls: true  # TLS handled by Cloudflare tunnel
}
```

Expose via Cloudflare Tunnel (already running):
```
cloudflared tunnel route dns onemind-tunnel nats-ws.onemindos.dev
```

## Cloudflare Worker Code

```javascript
// webhooks.onemindos.dev — Universal Webhook → NATS Bridge
// Deploy: wrangler deploy

import { connect } from 'nats.ws'  // @nats-io/nats.ws

export default {
  async fetch(request, env) {
    if (request.method !== 'POST') {
      return new Response('Method not allowed', { status: 405 })
    }

    const url = new URL(request.url)
    // URL pattern: /webhooks/<service>/<event>
    const [, , service, ...eventParts] = url.pathname.split('/')
    const event = eventParts.join('.') || 'default'

    if (!service) {
      return new Response('Missing service', { status: 400 })
    }

    // Validate webhook signature (per-service)
    const isValid = await validateSignature(request, service, env)
    if (!isValid) {
      return new Response('Invalid signature', { status: 401 })
    }

    // Parse body
    const body = await request.text()
    let payload
    try {
      payload = JSON.parse(body)
    } catch {
      payload = { raw: body }
    }

    // Publish to NATS
    const nc = await connect({
      servers: env.NATS_WS_URL,  // wss://nats-ws.onemindos.dev
      token: env.NATS_TOKEN,
    })

    const subject = `events.webhook.${service}.${event}`
    nc.publish(subject, JSON.stringify({
      service,
      event,
      payload,
      headers: Object.fromEntries(request.headers),
      received_at: new Date().toISOString(),
    }))

    await nc.flush()
    await nc.close()

    return new Response('OK', { status: 200 })
  }
}

async function validateSignature(request, service, env) {
  // Per-service signature validation
  const signatures = {
    stripe: async (req) => {
      const sig = req.headers.get('stripe-signature')
      const body = await req.clone().text()
      // Verify using env.STRIPE_WEBHOOK_SECRET
      return verifyStripeSignature(body, sig, env.STRIPE_WEBHOOK_SECRET)
    },
    github: async (req) => {
      const sig = req.headers.get('x-hub-signature-256')
      const body = await req.clone().text()
      return verifyHMAC(body, sig, env.GITHUB_WEBHOOK_SECRET, 'sha256')
    },
    linear: async (req) => {
      const sig = req.headers.get('linear-signature')
      const body = await req.clone().text()
      return verifyHMAC(body, sig, env.LINEAR_WEBHOOK_SECRET, 'sha256')
    },
    // Add more services as needed
    default: async () => true,  // No validation (development)
  }

  const validator = signatures[service] || signatures.default
  return validator(request)
}

async function verifyHMAC(body, signature, secret, algorithm) {
  const encoder = new TextEncoder()
  const key = await crypto.subtle.importKey(
    'raw', encoder.encode(secret),
    { name: 'HMAC', hash: algorithm.toUpperCase() },
    false, ['sign']
  )
  const signed = await crypto.subtle.sign('HMAC', key, encoder.encode(body))
  const expected = `sha256=${Buffer.from(signed).toString('hex')}`
  return expected === signature
}
```

## Worker Secrets (wrangler)

```bash
wrangler secret put NATS_WS_URL      # wss://nats-ws.onemindos.dev
wrangler secret put NATS_TOKEN        # NATS auth token
wrangler secret put STRIPE_WEBHOOK_SECRET
wrangler secret put GITHUB_WEBHOOK_SECRET
wrangler secret put LINEAR_WEBHOOK_SECRET
```

## NATS Subject Convention

```
events.webhook.<service>.<event>

# Examples:
events.webhook.stripe.invoice.paid
events.webhook.github.push
events.webhook.linear.issue.created
events.webhook.todoist.item.completed
```

## Subscriber Rules (NATS KV)

```json
{
  "subject": "events.webhook.stripe.invoice.*",
  "actions": [
    {"type": "telegram", "message": "💰 Stripe: {payload.type} — ${payload.amount}"},
    {"type": "kv_put", "bucket": "zeus-context", "key": "last_payment", "value": "{payload.id}"}
  ]
}
```

## Deployment Steps

1. Enable NATS WebSocket (add to nats-server.conf)
2. Expose via Cloudflare Tunnel (`nats-ws.onemindos.dev`)
3. `wrangler init webhooks && wrangler deploy`
4. Set secrets via `wrangler secret put`
5. Register webhook URLs in each SaaS: `https://webhooks.onemindos.dev/<service>/<event>`

## Why NOT Redpanda/Kafka

- Redpanda = 100MB Kafka clone designed for millions of events/sec throughput
- Your webhook volume = maybe 100/day
- NATS WebSocket handles this with zero additional infrastructure
- Cloudflare Worker is free, global, and validates signatures at the edge
