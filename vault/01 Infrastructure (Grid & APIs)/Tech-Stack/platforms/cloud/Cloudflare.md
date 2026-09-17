# Cloudflare

> **CDN, DNS & Edge Services**

## Overview

Cloudflare provides CDN, DNS, DDoS protection, and edge computing services for OneMind OS infrastructure.

## Status

- **Active**: ✅ Yes
- **Plan**: Free / Pro
- **Last Updated**: January 25, 2026

## Services in Use

| Service             | Purpose               | Status     |
| ------------------- | --------------------- | ---------- |
| **DNS**             | Domain management     | ✅ Active  |
| **CDN**             | Static asset delivery | ✅ Active  |
| **SSL/TLS**         | HTTPS certificates    | ✅ Active  |
| **DDoS Protection** | Security              | ✅ Active  |
| **Workers**         | Edge functions        | 🔄 Planned |
| **R2**              | Object storage        | 🔄 Planned |
| **Tunnels**         | Secure tunnels        | ✅ Active  |

## DNS Configuration

### Domain Setup

1. Add domain to Cloudflare
2. Update nameservers at registrar
3. Configure DNS records

### Common Records

| Type  | Name | Content            | Proxied |
| ----- | ---- | ------------------ | ------- |
| A     | @    | 123.45.67.89       | ✅      |
| A     | api  | 123.45.67.89       | ✅      |
| CNAME | www  | @                  | ✅      |
| CNAME | n8n  | server.example.com | ✅      |
| MX    | @    | mail.zoho.com      | ❌      |
| TXT   | @    | v=spf1 ...         | ❌      |

### API Usage

```python
import CloudFlare

cf = CloudFlare.CloudFlare(
    email="your@email.com",
    key=os.getenv("CLOUDFLARE_API_KEY")
)

# List zones
zones = cf.zones.get()

# Create DNS record
zone_id = "your_zone_id"
cf.zones.dns_records.post(zone_id, data={
    "type": "A",
    "name": "api",
    "content": "123.45.67.89",
    "proxied": True
})
```

## SSL/TLS

### Settings

- **SSL Mode**: Full (strict)
- **Minimum TLS Version**: 1.2
- **Always Use HTTPS**: Enabled
- **Automatic HTTPS Rewrites**: Enabled

### Origin Certificates

```bash
# Generate origin certificate in Cloudflare dashboard
# Install on origin server
```

## Cloudflare Tunnels

### Installation

```bash
# Install cloudflared
brew install cloudflared

# Login
cloudflared tunnel login

# Create tunnel
cloudflared tunnel create onemind
```

### Configuration

```yaml
# ~/.cloudflared/config.yml
tunnel: onemind-tunnel-id
credentials-file: /path/to/credentials.json

ingress:
  - hostname: api.onemind.io
    service: http://localhost:8080
  - hostname: n8n.onemind.io
    service: http://localhost:5678
  - hostname: home.onemind.io
    service: http://192.168.1.100:8123
  - service: http_status:404
```

### Run Tunnel

```bash
# Run manually
cloudflared tunnel run onemind

# Or as service
sudo cloudflared service install
sudo systemctl start cloudflared
```

### Docker Compose

```yaml
services:
  cloudflared:
    image: cloudflare/cloudflared
    command: tunnel run onemind
    environment:
      - TUNNEL_TOKEN=${CLOUDFLARE_TUNNEL_TOKEN}
    restart: unless-stopped
```

## Cloudflare Workers

### Basic Worker

```javascript
// worker.js
export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    // Add custom header
    const response = await fetch(request);
    const newResponse = new Response(response.body, response);
    newResponse.headers.set("X-Custom-Header", "value");

    return newResponse;
  },
};
```

### Deploy

```bash
npm install -g wrangler
wrangler login
wrangler deploy
```

### wrangler.toml

```toml
name = "onemind-worker"
main = "src/worker.js"
compatibility_date = "2024-01-01"

[vars]
ENVIRONMENT = "production"

[[kv_namespaces]]
binding = "CACHE"
id = "abc123"
```

## Page Rules

### Caching

| URL Pattern            | Rule                                |
| ---------------------- | ----------------------------------- |
| `*onemind.io/static/*` | Cache Everything, Edge TTL: 1 month |
| `*onemind.io/api/*`    | Bypass Cache                        |

### Security

| URL Pattern           | Rule                 |
| --------------------- | -------------------- |
| `*onemind.io/admin/*` | Security Level: High |
| `*onemind.io/*`       | SSL: Full (strict)   |

## R2 Storage

### Setup

```python
import boto3

s3 = boto3.client(
    's3',
    endpoint_url=f"https://{account_id}.r2.cloudflarestorage.com",
    aws_access_key_id=os.getenv("R2_ACCESS_KEY"),
    aws_secret_access_key=os.getenv("R2_SECRET_KEY")
)

# Upload file
s3.upload_file('local.txt', 'bucket-name', 'remote.txt')

# Generate presigned URL
url = s3.generate_presigned_url(
    'get_object',
    Params={'Bucket': 'bucket', 'Key': 'file.txt'},
    ExpiresIn=3600
)
```

## Zero Trust

### Access Policies

1. Create access policy in dashboard
2. Add allowed emails/groups
3. Apply to hostnames

### Tunnel + Access

```yaml
# config.yml
ingress:
  - hostname: admin.onemind.io
    service: http://localhost:8080
    originRequest:
      access:
        required: true
        teamName: onemind
```

## Security Features

### WAF Rules

- Enable Managed Rules
- Block known bad bots
- Rate limiting for API endpoints

### Rate Limiting

```
URL: /api/*
Rate: 100 requests per 10 seconds
Action: Block for 1 hour
```

## Analytics

### Key Metrics

- Requests
- Bandwidth saved
- Threats blocked
- Cache hit ratio

### API

```python
# Get zone analytics
analytics = cf.zones.analytics.dashboard.get(
    zone_id,
    params={'since': '-1d', 'until': 'now'}
)
```

## Environment Variables

```env
CLOUDFLARE_API_KEY=your_api_key
CLOUDFLARE_EMAIL=your@email.com
CLOUDFLARE_ZONE_ID=zone_id
CLOUDFLARE_TUNNEL_TOKEN=tunnel_token
R2_ACCESS_KEY=r2_access_key
R2_SECRET_KEY=r2_secret_key
```

## Best Practices

- ✅ Enable "Always Use HTTPS"
- ✅ Use Full (strict) SSL
- ✅ Enable DDoS protection
- ✅ Set up Page Rules for caching
- ✅ Use Tunnels for internal services
- ✅ Enable bot protection
- ✅ Monitor analytics regularly

## Resources

- [Cloudflare Documentation](https://developers.cloudflare.com/)
- [Workers Documentation](https://developers.cloudflare.com/workers/)
- [Tunnels Guide](https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/)

---

**Related**: [Docker](./Docker.md) | [AWS](./AWS.md)
