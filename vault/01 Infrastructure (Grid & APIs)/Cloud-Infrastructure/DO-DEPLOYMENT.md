# DigitalOcean Deployment Guide

> **OneMind OS — Moving from Lightsail to DigitalOcean.**
> Template security + full customization. Single Droplet.

---

## Architecture: One Droplet, Everything

```
┌─────────────────────────────────────────────────────────────────────┐
│                 DIGITALOCEAN DROPLET                                  │
│                 s-2vcpu-4gb ($24/mo) · Ubuntu 24.04                  │
│                 Tailscale + UFW (zero public ports except SSH)        │
│                                                                     │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  OPENCLAW GATEWAY (bare metal, systemd)                        │ │
│  │  • Port 18789 (Tailscale-only)                                │ │
│  │  • SOUL.md, AGENTS.md, skills/ at ~/.openclaw/                │ │
│  │  • Codex vault at /home/onemind/codex/ (git pull cron)       │ │
│  │  • BM25+vector hybrid search over entire Codex (native)       │ │
│  │  • GRADIENT_API_KEY → LLM inference from DO $5K credit        │ │
│  │  • Channels: Discord, WhatsApp, Telegram, iMessage, Voice     │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  ┌──── Docker Compose ────────────────────────────────────────────┐ │
│  │  NATS JetStream    (:4222 client, :7422 leaf, :1883 MQTT)     │ │
│  │  TimescaleDB       (:5432 PostgreSQL + time-series)           │ │
│  │  NATS Ingester     (NATS → TimescaleDB pipeline)              │ │
│  │  Legacy Bridge     (NATS → OpenClaw webhook relay)            │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  Security:                                                          │
│  • UFW: deny all incoming, allow SSH, Tailscale-only for services  │
│  • fail2ban: SSH brute-force protection                            │
│  • Non-root 'openclaw' user for all services                       │
│  • Docker isolation for NATS/TimescaleDB                           │
│  • Tailscale mesh — no public ports for NATS/DB/OpenClaw           │
└─────────────────────────────────────────────────────────────────────┘
         │
         │ NATS Leaf Node (port 7422 via Tailscale)
         │
┌────────▼────────────────────────────────────────┐
│           ZEUS'S LOCAL OFFICE                    │
│                                                  │
│  Mac: Obsidian + dev (Codex editor)              │
│  NATS Leaf Node → DO NATS hub                    │
│  Home Assistant → MQTT → NATS (:1883)            │
│  Sensors, cameras, smart devices                 │
└──────────────────────────────────────────────────┘
```

---

## Why One Droplet (Not App Platform)

| Concern | Why Droplet Wins |
|---------|-----------------|
| **Persistent storage** | Codex vault, skills, memory SQLite all persist across restarts. App Platform has ephemeral storage — wiped on every deploy. |
| **Tailscale** | Install directly, join mesh. App Platform has limited Tailscale support. |
| **Custom ports** | NATS needs 4222, 7422, 1883, 8222. App Platform only exposes one HTTP port. |
| **Codex sync** | Simple `git pull` cron job. App Platform would need hacky build-time cloning. |
| **Skill installs** | `cd ~/.openclaw/workspace/skills/ && npm install` — native filesystem. |
| **Security** | UFW + fail2ban + Tailscale-only ports = same security as DO templates, but YOU control it. |
| **Cost** | $24/mo for everything. App Platform + Droplet would be $36+. |
| **Flexibility** | Full root access. Add anything later — more Docker services, local scripts, cron jobs. |

---

## Q4: Yes, TimescaleDB IS PostgreSQL

TimescaleDB is a **PostgreSQL extension**, not a separate database engine:

```sql
-- It's literally this:
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Then use PostgreSQL exactly like normal, but with superpowers:
SELECT create_hypertable('events', 'timestamp');  -- automatic time partitioning
SELECT add_compression_policy('events', INTERVAL '7 days');  -- 95%+ compression
SELECT time_bucket('1 hour', timestamp), avg(value) FROM events GROUP BY 1;
```

The Docker image `timescale/timescaledb:latest-pg16` is PostgreSQL 16 + TimescaleDB extension pre-installed. You get:
- Full PostgreSQL (SQL, JSONB, schemas, roles, pg_dump, pgvector)
- TimescaleDB (hypertables, continuous aggregates, compression, retention policies)
- Regular PostgreSQL tools work: `psql`, DBeaver, pgAdmin, any ORM

**Your `init.sql` already handles all of this.** The `docker-compose.yml` mounts it as `/docker-entrypoint-initdb.d/01-init.sql` — PostgreSQL auto-runs it on first boot.

---

## Step-by-Step Deployment

### Phase 1: NATS + TimescaleDB Droplet (Do This NOW)

**Time: ~30 minutes**

#### 1. Create Droplet

```bash
# Via doctl CLI
doctl compute droplet create onemind-nats \
  --image ubuntu-24-04-x64 \
  --size s-2vcpu-4gb \
  --region nyc1 \
  --ssh-keys <your-key-id> \
  --tag-names onemind,infrastructure

# Or via DigitalOcean dashboard:
# → Droplets → Create → Ubuntu 24.04 → s-2vcpu-4gb → NYC1
```

#### 2. Run Setup Script

```bash
# SSH into the new Droplet
ssh root@<droplet-ip>

# Clone repo and run setup
git clone https://github.com/Zeus-Delacruz/onemindos-v0.4.git /opt/onemind
cd /opt/onemind
bash .do/setup-droplet.sh
```

#### 3. Join Tailscale

```bash
sudo tailscale up
# Follow the link to authorize
# Note the Tailscale IP: tailscale ip -4
```

#### 4. Open Ports on Tailscale Only

```bash
bash /opt/onemind/.do/setup-firewall.sh
```

#### 5. Edit Environment

```bash
nano /opt/onemind/.env
# Set TSDB_PASSWORD to a strong password
# Set OPENCLAW_GATEWAY_URL to OpenClaw's Tailscale IP
# Set OPENCLAW_HOOKS_TOKEN
```

#### 6. Restart with Real Config

```bash
cd /opt/onemind
docker compose --profile data down
docker compose --profile data up -d
docker compose --profile data logs -f  # Watch for "Ingester running"
```

#### 7. Verify

```bash
# Check NATS
docker exec onemind-nats nats-server --signal ldm 2>&1 || echo "NATS is running"

# Check TimescaleDB
docker exec onemind-tsdb psql -U onemind -d onemind_events -c "SELECT count(*) FROM events;"

# Check from Mac (via Tailscale)
nats sub '>' --server nats://$(tailscale ip -4 onemind-nats):4222
```

---

### Phase 2: Move OpenClaw to App Platform (After NATS is stable)

**Time: ~20 minutes**

#### 1. Deploy via App Spec

```bash
doctl apps create --spec .do/app.yaml
```

Or use the DigitalOcean dashboard:
- Apps → Create App → GitHub → onemindos-v0.4 → Select Dockerfile

#### 2. Set Environment Variables

In DO dashboard → App → Settings → Environment Variables:
- `GRADIENT_API_KEY` → your DO API key (LLM inference from $5K credit)
- `OPENCLAW_GATEWAY_TOKEN` → generate a strong token
- `OPENCLAW_HOOKS_TOKEN` → match the Droplet's `.env`
- `NATS_URL` → `nats://onemind-nats.tail-xxxxx.ts.net:4222`
- All channel tokens (Discord, Telegram, WhatsApp, etc.)

#### 3. Add Tailscale to App Platform

App Platform supports Tailscale via built-in integration:
```bash
# In DO dashboard → App → Settings → Internal Services
# Enable Tailscale → paste auth key
```

#### 4. Clone Codex Inside App Platform

Since App Platform is ephemeral storage, the Codex sync needs a build step:

```dockerfile
# Add to openclaw/Dockerfile:
RUN apt-get update && apt-get install -y git cron

# Clone Codex at build time (snapshot)
RUN git clone https://github.com/Zeus-Delacruz/onemindcodex-v0.4.git /home/onemind/codex

# Cron pulls updates at runtime
COPY .do/codex-sync-cron /etc/cron.d/codex-sync
```

**Alternative:** Use DO Spaces to host the Codex files + Gradient KB for search. This avoids git-in-Docker complexity.

#### 5. Verify

```bash
# Health check
curl -H "Authorization: Bearer $TOKEN" https://openclaw-xxxxx.ondigitalocean.app/health

# Or via Tailscale
curl -H "Authorization: Bearer $TOKEN" http://openclaw.tail-xxxxx.ts.net:18789/health
```

---

### Phase 3: Local NATS Leaf Node (Home Assistant)

**Time: ~1 hour**

#### On your Mac Mini or local server:

```bash
# Install NATS server
brew install nats-server

# Create leaf config
cat > /usr/local/etc/nats/leaf.conf << 'EOF'
server_name: onemind-office-leaf
port: 4222
http_port: 8222

leafnodes {
  remotes [
    {
      url: "nats://edge-office:$EDGE_OFFICE_PASSWORD@onemind-nats.tail-xxxxx.ts.net:7422"
    }
  ]
}

jetstream {
  store_dir: /usr/local/var/nats
  max_file_store: 256MB
  max_memory_store: 64MB
}
EOF

# Start
nats-server -c /usr/local/etc/nats/leaf.conf
```

#### Home Assistant → NATS

Two options:

**Option A: MQTT (simplest)**
Home Assistant has native MQTT support. NATS has native MQTT support.
Just point Home Assistant's MQTT broker at your local NATS port 1883:

```yaml
# Home Assistant configuration.yaml
mqtt:
  broker: 127.0.0.1
  port: 1883
  # NATS handles MQTT natively — no Mosquitto needed
```

**Option B: Custom Python bridge** (more control)
See `openclaw/nats/nats_bridge.py` for the pattern.

---

## Cost Summary

| Phase | Service | Cost/mo |
|-------|---------|---------|
| 1 | NATS+TimescaleDB Droplet (s-2vcpu-4gb) | $24 |
| 2 | OpenClaw App Platform (professional-xs) | $12 |
| 2 | LLM inference (Gradient, from $5K credit) | $0* |
| 3 | Tailscale (free tier, up to 100 devices) | $0 |
| 3 | Local NATS leaf (your hardware) | $0 |
| **Total** | | **$36/mo** |

*$5K DO credit covers ~$100-200/mo of LLM inference for months.

---

## File Map (What Was Created/Modified)

| File | What | Status |
|------|------|--------|
| `docker-compose.yml` | Added TimescaleDB + NATS ingester services | ✅ Updated |
| `.do/app.yaml` | DO App Platform spec for OpenClaw | ✅ Created |
| `.do/setup-droplet.sh` | Automated Droplet setup script | ✅ Created |
| `openclaw/nats/nats_ingester.py` | NATS → TimescaleDB pipeline | ✅ Created |
| `openclaw/nats/timescaledb/init.sql` | TimescaleDB schema | ✅ Already existed |
| `openclaw/nats/hub.conf` | NATS cloud hub config | ✅ Already existed |
| `openclaw/nats/leaf.conf` | NATS edge leaf config | ✅ Already existed |
| `openclaw/Dockerfile` | OpenClaw container image | ✅ Already existed |
| `openclaw/docs/HYBRID-SEARCH-SETUP.md` | BM25+vector search reference | ✅ Created |
| `openclaw/docs/DO-DEPLOYMENT.md` | This file | ✅ Created |

---

*Last updated: 2026-03-01*
