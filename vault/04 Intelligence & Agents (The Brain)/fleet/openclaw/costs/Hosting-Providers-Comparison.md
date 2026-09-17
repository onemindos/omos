# OpenClaw Hosting Providers Comparison

> Extracted from [awesome-openclaw](https://github.com/rohitg00/awesome-openclaw) — March 2026

## Free Tier ($0/month)

| Provider | Free Resources | Limits | Permanent? | Setup Time | Best For |
|----------|---------------|--------|------------|------------|----------|
| **[Oracle Cloud](https://www.oracle.com/cloud/free/)** | 4 ARM CPUs, 24 GB RAM, 200 GB storage, 10 TB/mo | ARM architecture only | Yes (forever) | ~3 hours | Power users wanting $0 hosting |
| **[AMD Developer Cloud](https://www.amd.com/en/developer/resources/cloud-access.html)** | MI300X GPU (192 GB memory), $100 credits | ~50 hours compute | No (credits) | 30 min | GPU-accelerated local LLMs |
| **[Google Cloud Run](https://cloud.google.com/run)** | 180K vCPU-sec, 360K GiB-sec, 2M requests/mo | Cold starts, stateless | Yes | 30 min | Serverless testing |
| **[Azure Container Apps](https://azure.microsoft.com/en-us/products/container-apps)** | 180K vCPU-sec, 360K GiB-sec, 2M requests/mo | Cold starts | Yes | 30 min | Enterprise users |
| **[Render](https://render.com/)** | Web service (spins down after 15 min) | Slow wake-up (~60s) | Yes | 5 min | Testing only |
| **[AWS Free Tier](https://aws.amazon.com/free/)** | t2.micro (1 vCPU, 1 GB RAM) | 12 months only | No (12 mo) | 15 min | AWS-familiar users |
| **[Railway](https://railway.com/deploy/openclaw)** | $5 one-time credit | 30 days trial | No | 2 min | Quick testing |

## Budget VPS ($2-8/month)

| Provider | Plan | Price/mo | vCPU | RAM | Storage | Bandwidth | Region |
|----------|------|----------|------|-----|---------|-----------|--------|
| **[LumaDock](https://lumadock.com/)** | Basic | $1.99 | 1 | 1 GB | 20 GB SSD | 1 TB | US/EU |
| **[Vultr](https://www.vultr.com/)** | Regular Cloud | $2.50 | 1 | 512 MB | 10 GB | 0.5 TB | 32 locations |
| **[Vultr](https://www.vultr.com/)** | Standard | $3.50 | 1 | 1 GB | 25 GB NVMe | 1 TB | 32 locations |
| **[Hetzner](https://www.hetzner.com/cloud/)** | CX22 | €3.79 (~$4.15) | 2 | 4 GB | 40 GB SSD | 20 TB | EU (DE/FI) |
| **[Contabo](https://contabo.com/en/openclaw-hosting/)** | Cloud VPS S | $4.95 | 4 | 8 GB | 50 GB SSD | Unlimited | US/EU/Asia |
| **[Hostinger](https://www.hostinger.com/uk/vps/docker/openclaw)** | KVM 1 | $4.99 | 1 | 4 GB | 50 GB NVMe | 4 TB | Global |
| **[DigitalOcean](https://marketplace.digitalocean.com/apps/openclaw)** | Basic Droplet | $6 | 1 | 1 GB | 25 GB | 1 TB | Global |
| **[Hetzner](https://www.hetzner.com/cloud/)** | CX32 | €6.80 (~$7.45) | 4 | 8 GB | 80 GB SSD | 20 TB | EU (DE/FI) |

## Mid-Range ($5-25/month)

| Provider | Plan | Price/mo | vCPU | RAM | Storage | Notes |
|----------|------|----------|------|-----|---------|-------|
| **[DigitalOcean](https://marketplace.digitalocean.com/apps/openclaw)** | Premium Droplet | $7 | 1 | 1 GB | 25 GB NVMe | 1-Click Deploy available |
| **[DigitalOcean](https://marketplace.digitalocean.com/apps/openclaw)** | Standard 2GB | $12 | 1 | 2 GB | 50 GB | Recommended for OpenClaw |
| **[Hostinger](https://www.hostinger.com/uk/vps/docker/openclaw)** | KVM 2 | $6.99 | 2 | 8 GB | 100 GB NVMe | Best Hostinger value |
| **[GCP](https://cloud.google.com/compute)** | e2-medium | ~$12 | 2 | 4 GB | 30 GB SSD | $300 free credit (90 days) |
| **[Azure](https://azure.microsoft.com/en-us/products/virtual-machines)** | B1ms VM | ~$15 | 1 | 2 GB | 32 GB | Enterprise compliance |

## Serverless & PaaS

| Platform | Starting Price | Free Tier? | Persistent Storage | Auto-Scale | Notes |
|----------|---------------|------------|-------------------|------------|-------|
| **[Cloudflare Workers](https://workers.cloudflare.com/)** | $5/mo (paid plan) | 100K req/day (free) | R2 ($0.015/GB/mo) | Yes | Moltworker PoC |
| **[Railway](https://railway.com/deploy/openclaw)** | $5/mo (Hobby) | $5 credit trial | Yes (volumes) | Yes | Best PaaS UX |
| **[Render](https://render.com/)** | $7/mo (web service) | Free w/ limits | Yes (disks) | Yes | YAML IaC |
| **[Fly.io](https://fly.io/)** | Pay-as-you-go | None | Yes (volumes) | Yes | Global edge |
| **[AWS Lightsail](https://aws.amazon.com/lightsail/)** | $3.50/mo | None | Yes | No | Simple AWS VPS |

## Managed Hosting Services

| Provider | Price/mo | Setup Time | What's Included |
|----------|----------|------------|-----------------|
| **[Agent37](https://www.agent37.com/openclaw)** | $0.99-3.99 | 30 sec | Isolated instance, full UI, SSL |
| **[MyClaw.ai](https://myclaw.ai/pricing)** | $9 (was $29) | Instant | Auto-updates, backups, web terminal |
| **[ClawCloud](https://www.clawcloud.sh/)** | $29-129 | < 1 min | All AI models (BYOK), auto-updates |
| **[Contabo OpenClaw](https://contabo.com/en/openclaw-hosting/)** | $4.95+ | Minutes | Unlimited workflows, full data ownership |
| **[Molty by Finna](https://molty.finna.ai)** | $19-99 | 30 sec | Multi-Molty, Mission Control, GitHub sync |

## Master Cost Comparison

| Provider | Monthly Cost | Setup Time | Difficulty | Best For |
|----------|-------------|------------|------------|----------|
| Oracle Cloud | **$0** | 3 hours | Hard | Free-forever hosting |
| Agent37 | **$0.99** | 30 sec | None | Cheapest managed |
| LumaDock | **$1.99** | 15 min | Medium | Cheapest VPS |
| Hetzner CX22 | **$4.15** | 10 min | Medium | Best price/performance |
| DigitalOcean | **$6** | 2 min | Easy | Best docs + 1-Click |
| MyClaw.ai Lite | **$9** | Instant | None | Budget managed |
| DigitalOcean 2GB | **$12** | 2 min | Easy | Recommended production |

> **OneMind runs on**: DigitalOcean droplet `s-4vcpu-8gb-240gb-intel` ($64/mo). If migrating, Hetzner CX32 ($7.45) or Contabo VPS M ($7.95) offer similar specs at ~90% less.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
