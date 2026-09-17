---
title: "{{Machine Name}}"
type: machine
status: active
created: {{DATE}}
updated: {{DATE}}
owner: zeus
access_level: L0
subtype: server
tags: [type/machine, domain/SO, status/active]
ai_summary: "{{Machine description}} — {{VPS / Mac / server / NAS}}. {{primary role}}."
---

# {{Machine Name}}

## Identity

| Field | Value |
|-------|-------|
| Name / hostname | {{e.g. "onemind-main", "zeus-macbook"}} |
| Type | vps / mac / linux-server / nas / raspberry-pi |
| Make / Model | {{DigitalOcean Droplet / MacBook Pro M4 / etc.}} |
| OS | {{Ubuntu 24.04 / macOS 15 / etc.}} |
| CPU | {{specs}} |
| RAM | {{GB}} |
| Storage | {{GB / TB}} |
| Owned / Leased | {{owned / DigitalOcean / Hetzner / etc.}} |

## Network

| Field | Value |
|-------|-------|
| Public IP | {{IP address}} |
| Tailscale IP | {{100.x.x.x}} |
| Tailscale hostname | {{hostname.tail...ts.net}} |
| Local IP | {{192.168.x.x if applicable}} |
| SSH access | `ssh root@{{IP}}` |
| SSH key | In 1Password: `{{path}}` |

## Services Running

| Service | Port | Purpose |
|---------|------|---------|
| {{e.g. openclaw}} | {{3000}} | {{agent gateway}} |
| {{e.g. nats}} | {{4222}} | {{messaging}} |
| {{e.g. postgres}} | {{5432}} | {{database}} |

## Access

| Method | Details |
|--------|---------|
| SSH | `ssh root@{{IP}}` |
| Tailscale | `ssh root@{{tailscale-hostname}}` |
| Web UI | `http://{{IP}}:{{port}}` |
| Credentials | In 1Password: `{{path}}` |

## Cost (if cloud)

| Field | Value |
|-------|-------|
| Provider | {{DigitalOcean / Hetzner / etc.}} |
| Plan | {{Droplet name / tier}} |
| Monthly cost | ${{amount}} |
| Billing account | In 1Password: `{{path}}` |
| Renewal date | {{date if annual}} |

## Notes

{{What runs here, backups, important configs}}


---

## Related

- [[00-24 SO (Sovereign Operations)/08 OM Assets (Entity Registry) (The Registry) (Hive Mind)/08 OM Assets|08 OM Assets]]
