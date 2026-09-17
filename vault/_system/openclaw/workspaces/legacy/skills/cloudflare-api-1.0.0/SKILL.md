---
name: cloudflare-api
description: Cloudflare DNS and zone management. Legacy can list domains, manage DNS records, and inspect zone settings.
---

# Cloudflare API Skill

6 tools for managing Cloudflare DNS via the v4 API.

## Tools

| Tool               | Purpose                                  |
| ------------------ | ---------------------------------------- |
| `cf_zones`         | List all domains in the account          |
| `cf_dns_list`      | List DNS records (filter by type/name)   |
| `cf_dns_create`    | Create a DNS record                      |
| `cf_dns_update`    | Update an existing record                |
| `cf_dns_delete`    | Delete a record                          |
| `cf_zone_settings` | Get zone SSL, caching, security settings |

## Environment

| Var                    | Required | Description                                     |
| ---------------------- | -------- | ----------------------------------------------- |
| `CLOUDFLARE_API_TOKEN` | Yes      | API token with Zone:Read + DNS:Edit permissions |

## Token Permissions

Create at [dash.cloudflare.com/profile/api-tokens](https://dash.cloudflare.com/profile/api-tokens):

- **Zone:Read** — list domains
- **DNS:Edit** — manage records

## Usage Examples

```
# List all domains
cf_zones()

# List DNS for onemindos.dev
cf_dns_list({ zone: "onemindos.dev" })

# Create A record pointing to Tailscale IP
cf_dns_create({ zone: "onemindos.dev", type: "A", name: "gateway", content: "100.77.52.93" })

# Create wildcard
cf_dns_create({ zone: "onemindos.dev", type: "A", name: "*", content: "100.77.52.93" })
```
