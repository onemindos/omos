# OneMind OS — N8N Setup

> Workflow automation · Docker profile: `full` · **NOT currently running**

---

## Status

N8N is defined in the Docker Compose stack but **not running** on the VPS. It is in the `full` profile alongside LiveKit. The default stack (NATS + OmniDB + ingester + bridge) runs separately.

---

## Architecture Role

```
Legacy (OpenClaw) ──fires webhook──▶ N8N workflow ──executes──▶ external APIs / services
                                         │
                                   (complex multi-step automations)
                                   Slack, email, Sheets, Notion, Zoho, etc.
```

OpenClaw handles **reasoning** — N8N handles **multi-step execution**. When a task is too complex to wire manually in a skill, Legacy delegates to an N8N workflow.

---

## Docker Config

```yaml
# From /opt/onemind/docker-compose.yml (profile: full)
onemind-n8n:
  image: n8nio/n8n:latest
  container_name: onemind-n8n
  profiles: [full]
  restart: unless-stopped
  ports:
    - "5678:5678"
  volumes:
    - n8n_data:/home/node/.n8n
  environment:
    N8N_HOST: "0.0.0.0"
    N8N_PORT: "5678"
    N8N_PROTOCOL: "http"
    N8N_ENCRYPTION_KEY: ${N8N_ENCRYPTION_KEY}
    WEBHOOK_URL: ${N8N_WEBHOOK_URL:-http://host.docker.internal:5678}
  # Optional: OmniDB persistence (commented out)
  # DB_TYPE: postgresdb
  # DB_POSTGRESDB_HOST: omnidb
  # DB_POSTGRESDB_PORT: 5432
```

---

## Ports

| Port | Use |
|------|-----|
| 5678 | N8N web UI + webhook receiver |

---

## Required Secrets

| Var | Description |
|-----|-------------|
| `N8N_ENCRYPTION_KEY` | Encrypts stored credentials in N8N — set in `/opt/onemind/.env.sa` or via 1Password |
| `N8N_WEBHOOK_URL` | External URL for webhook triggers (e.g. `https://n8n.onemind.ai` or Tailscale URL) |

---

## Start / Manage

```bash
# Start n8n only
cd /opt/onemind
docker compose --profile full up -d n8n

# Stop n8n
docker compose stop n8n

# View logs
docker logs -f onemind-n8n

# Access UI (via Tailscale or SSH tunnel)
# http://104.248.4.212:5678
```

---

## Data Persistence

N8N stores workflows, credentials, and execution history in:
- Docker volume: `onemind_n8n_data` → `/home/node/.n8n` inside container
- Optional Postgres persistence using OmniDB (commented out in compose)

---

## Legacy Integration Pattern

1. Legacy decides an action needs complex multi-step execution
2. Legacy calls `web_fetch` or `nexus_dispatch` to fire a webhook to N8N
3. N8N workflow runs (Slack DM, email send, Sheets append, etc.)
4. N8N can POST back a result to Legacy's webhook endpoint if needed

---

## Access Behind Tailscale

N8N port 5678 is not exposed publicly. To access:
- Tailscale: `http://onemind-main.tail717ba0.ts.net:5678`
- SSH tunnel: `ssh -L 5678:localhost:5678 root@104.248.4.212`
