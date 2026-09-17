# Lob

> **Physical mail API — send postcards, letters, and checks programmatically**

## Overview

Lob provides a REST API for physical mail. Send postcards, letters, self-mailers, and checks at scale. Handles printing, addressing, and USPS delivery. Useful for outreach automation, direct mail campaigns, and physical document delivery.

## Status

| Attribute   | Value                     |
| ----------- | ------------------------- |
| Status      | 🔄 Evaluating             |
| Category    | Communication             |
| Type        | Physical Mail API         |
| Integration | REST API, Node.js SDK     |

## Use Cases in OneMind

1. **Outreach Automation** — Agent-triggered physical mail for GE domain
2. **Direct Mail** — Marketing postcards for One Mind product launches
3. **Document Delivery** — Send physical contracts or documents
4. **Brand Touch** — Surprise physical mail to community members

## Key Features

- Postcard, letter, self-mailer, check APIs
- USPS Intelligent Mail barcode tracking
- Address verification API
- HTML or PDF template support
- Test mode for development
- Webhooks for delivery status

## Quick Start

```bash
curl -X POST https://api.lob.com/v1/postcards \
  -u LOB_API_KEY: \
  -d "to[name]=John Smith" \
  -d "to[address_line1]=123 Main St" \
  -d "from=adr_123" \
  -d "front=tmpl_456" \
  -d "back=tmpl_789"
```

## Links

- Website: [lob.com](https://lob.com)
- API Docs: [docs.lob.com](https://docs.lob.com)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
