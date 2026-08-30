# webhook-bridge

> **omos-bridge** component — HTTP → NATS  
> Receives inbound HTTP webhooks and publishes them as **One Mind envelopes** to NATS.

---

## What it does

`webhook-bridge` spins up a lightweight HTTP server.  Any `POST`, `GET`,
`PUT`, or `PATCH` request to `/{source}` is transformed into a structured
One Mind envelope and published to NATS subject
`omos.{domain}.{source}.telemetry`.

```
GitHub / Stripe / Twilio / Any HTTP client
        │  POST /my-sensor
        ▼
  webhook-bridge  (:8090)
        │  omos.net.my-sensor.telemetry
        ▼
      NATS  →  any consumer / JetStream
```

---

## Quick start

```bash
# 1. clone & enter
cd omos-bridge/bridges/webhook

# 2. copy env
cp .env.example .env
# edit .env as needed

# 3. install deps & build
go get github.com/nats-io/nats.go@latest
go build -o webhook-bridge .

# 4. run (NATS must be reachable)
source .env && ./webhook-bridge
```

---

## Environment variables

| Variable         | Default                  | Description                                                  |
|------------------|--------------------------|--------------------------------------------------------------|
| `NATS_URL`       | `nats://localhost:4222`  | NATS server URL                                              |
| `BRIDGE_NAME`    | `webhook-bridge`         | Logical bridge name (used in health/error subjects)          |
| `WEBHOOK_PORT`   | `8090`                   | HTTP listen port                                             |
| `WEBHOOK_SECRET` | *(empty)*                | If set, validate `X-Hub-Signature-256` HMAC-SHA256 header    |

---

## URL pattern

```
POST  http://host:8090/{source}
GET   http://host:8090/{source}
PUT   http://host:8090/{source}
PATCH http://host:8090/{source}
```

`{source}` becomes the entity ID. An optional `X-Omos-Domain` header
overrides the domain (default `net`).

---

## Envelope schema

```json
{
  "entity":    "net.my-sensor",
  "domain":    "net",
  "msg_type":  "webhook",
  "timestamp": 1718000000000,
  "data": {
    "source_url":   "http://host:8090/my-sensor",
    "method":       "POST",
    "headers":      { "Content-Type": "application/json" },
    "body_raw":     "{\"temp\":21.5}",
    "body_parsed":  { "temp": 21.5 },
    "received_at":  1718000000000
  }
}
```

---

## NATS subjects

| Subject                                  | Purpose                        |
|------------------------------------------|--------------------------------|
| `omos.{domain}.{source}.telemetry`       | Inbound webhook payload        |
| `omos.bridge.webhook-bridge.health`      | Heartbeat every 30 s           |
| `omos.bridge.webhook-bridge.errors`      | Error reports                  |

---

## Example curl

```bash
# Basic POST
curl -X POST http://localhost:8090/temperature \
  -H "Content-Type: application/json" \
  -d '{"value": 21.5, "unit": "C"}'

# With custom domain
curl -X POST http://localhost:8090/door-sensor \
  -H "X-Omos-Domain: home" \
  -H "Content-Type: application/json" \
  -d '{"state": "open"}'

# Response
# {"status":"ok","subject":"omos.home.door-sensor.telemetry"}
```

---

## Integration examples

### GitHub webhooks

1. In your repo → **Settings → Webhooks → Add webhook**
2. Payload URL: `https://your-host:8090/github-events`
3. Content type: `application/json`
4. Secret: matches `WEBHOOK_SECRET`
5. Bridge validates `X-Hub-Signature-256` automatically.

```bash
# Published to: omos.net.github-events.telemetry
```

### Stripe webhooks

1. Dashboard → **Developers → Webhooks → Add endpoint**
2. Endpoint URL: `https://your-host:8090/stripe`
3. (Stripe uses `Stripe-Signature`; add custom verification if needed.)

```bash
# Published to: omos.net.stripe.telemetry
```

### Twilio status callbacks

1. Twilio Console → phone number → **Status Callback URL**
2. Set to `https://your-host:8090/twilio`
3. Twilio sends `application/x-www-form-urlencoded` — body is captured in
   `body_raw`; `body_parsed` will be null (not JSON).

```bash
# Published to: omos.net.twilio.telemetry
```

---

## HMAC signature verification

Set `WEBHOOK_SECRET` to enable GitHub-style HMAC-SHA256 validation.  
The bridge reads `X-Hub-Signature-256: sha256=<hex>` and rejects requests
with an invalid or missing signature with HTTP **401**.

---

## Roadmap

- [ ] TLS / HTTPS listener
- [ ] JetStream publish with ack
- [ ] Per-source routing rules (YAML)
- [ ] Replay / dead-letter queue for failed publishes
- [ ] Prometheus `/metrics` endpoint
- [ ] Docker image & Helm chart
