# Resend

> **Developer-first email API for transactional and marketing emails**

## Overview

Resend is a modern email API built for developers. Clean REST API, React Email for templates, and excellent deliverability. Built by ex-Vercel engineers. The go-to choice for sending emails from Next.js apps and agent-driven notifications.

## Status

| Attribute   | Value                     |
| ----------- | ------------------------- |
| Status      | 🔄 Evaluating             |
| Category    | Communication             |
| Type        | Email API                 |
| Integration | REST API, Node.js SDK     |

## Use Cases in OneMind

1. **Agent Notifications** — Send email alerts from OpenClaw agents
2. **Transactional Email** — onemind-web user emails (signup, alerts)
3. **Digest Reports** — Weekly/monthly review reports via email
4. **Marketing** — One Mind newsletter delivery to community members

## Key Features

- React Email templates
- Custom domain sending
- Webhooks for delivery events
- Free tier: 3,000 emails/month
- SDKs for Node.js, Python, Go, Ruby

## Quick Start

```bash
npm install resend
```

```typescript
import { Resend } from 'resend';
const resend = new Resend('RESEND_API_KEY');

await resend.emails.send({
  from: 'Zeus <zeus@onemindos.com>',
  to: 'user@example.com',
  subject: 'One Mind Weekly',
  html: '<p>Report content</p>'
});
```

## Links

- Website: [resend.com](https://resend.com)
- Docs: [resend.com/docs](https://resend.com/docs)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
