# Cal.com

> **Open-Source Scheduling & Calendar Booking**

## Overview

Cal.com is an open-source, self-hostable scheduling platform (like Calendly). It lets anyone book time on your calendar based on your availability. Fully customizable, privacy-first, and can be self-hosted on the VPS.

## Status

- **Active**: ⏸ Pending setup
- **Plan**: Free (self-hosted) or Cloud
- **Last Updated**: 2026-05-01

## Use Cases

1. **Consultation Booking** — Let clients/students book 1-on-1 time with Zeus
2. **Coaching Sessions** — Layer 5 (Architect) intake scheduling
3. **Community Calls** — Layer 4 scheduling for OneMind community
4. **Agent Scheduling** — Agents can create/read booking links through API

## Key Features

- Fully open-source (GitHub: calcom/cal.com)
- Self-hostable (Docker) or cloud-hosted
- Integrates with Google Calendar, Outlook, iCloud
- Webhooks on booking/cancellation events
- Embeddable booking widget
- Zapier / n8n integration

## Integration Notes

- Can be self-hosted on VPS alongside existing Docker stack
- Webhook events → n8n → Todoist task creation on new bookings
- API: `https://cal.com/api/v1/` (cloud) or local if self-hosted

## Links

- Website: https://cal.com
- GitHub: https://github.com/calcom/cal.com
- API Docs: https://cal.com/docs/enterprise-features/api

## Related

- [[Zoho Calendar]] — existing calendar infrastructure
- `75-99 GE/` — client-facing scheduling
- `_system/infrastructure/` — VPS self-hosting option
