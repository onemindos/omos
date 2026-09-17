# Droplet Migration — Old Resize Ticket Cancelled

**Status**: ✅ New droplet provisioned 2026-03-13 — migration in progress
**Old Ticket**: https://cloud.digitalocean.com/support/tickets/11758225 (cancelled — resize impossible without vCPU downgrade)

---

## Decision: Fresh Droplet Instead of Resize

The 240GB disk lock-in forced Storage-Optimized tier (2 vCPU / 16GB at $131/mo) — a CPU downgrade.
Chose fresh droplet for right-tier, no lock-in, clean migration.

---

## Old Droplet (decommission after migration)

| Field | Value |
|-------|-------|
| Name | `onemind-gateway` |
| ID | `555777688` |
| IP | `134.122.8.135` |
| Region | nyc3 |
| Size | `s-4vcpu-8gb-240gb-intel` |
| Specs | 4 vCPUs / 8GB RAM / 240GB disk |
| Cost | $64/mo |
| Status | **Running — keep alive until migration complete** |

## New Droplet (onemind-main)

| Field | Value |
|-------|-------|
| Name | `onemind-main` |
| ID | `558147996` |
| IP | `104.248.4.212` |
| Region | nyc3 |
| Size | `s-4vcpu-16gb-amd` |
| Specs | 4 vCPUs / 16GB RAM / 200GB disk |
| Cost | $84/mo |
| Transfer | 8TB |
| VPC | `2e617f7b-2721-4371-b34c-b4f5d91aec76` (same as old) |
| Status | **Active — awaiting migration** |

**Savings vs resize option**: $12/mo less than Storage-Optimized, keeps 4 vCPUs

## Future Scale Target

- 8 vCPU / 32GB RAM — once production agent load increases

## Migration Checklist

### Phase 1 — Setup new droplet
- [ ] Run `openclaw/deploy/scripts/setup.sh` on new droplet
- [ ] Install Tailscale and join tailnet
- [ ] Install 1Password CLI (`op`)
- [ ] Copy `/opt/onemind/.env.tpl` and `/opt/onemind/.env.sa` from old droplet
- [ ] Restore Docker stack (NATS, PostgreSQL, n8n, LiveKit)

### Phase 2 — Data migration
- [ ] Dump PostgreSQL from old: `pg_dump -h localhost -U onemind onemind > /tmp/onemind.sql`
- [ ] Copy dump to new: `rsync /tmp/onemind.sql root@104.248.4.212:/tmp/`
- [ ] Restore on new: `psql -h localhost -U onemind onemind < /tmp/onemind.sql`
- [ ] Rsync codex + workspace: `rsync -avz /home/onemind/onemind-os/ root@104.248.4.212:/home/onemind/onemind-os/`

### Phase 3 — OpenClaw cutover
- [ ] Install OpenClaw 2026.3.12 on new droplet
- [ ] Start `openclaw.service` on new — verify health: `curl http://localhost:18789/health`
- [ ] Update Tailscale proxy to point to new droplet IP
- [ ] Update `infra/README.md` with new IP and specs
- [ ] Verify all 9 agents responding in Discord

### Phase 4 — Cleanup
- [ ] Monitor new droplet for 48h
- [ ] Destroy old droplet: `doctl compute droplet delete 555777688`
- [ ] Cancel DO support ticket #11758225
