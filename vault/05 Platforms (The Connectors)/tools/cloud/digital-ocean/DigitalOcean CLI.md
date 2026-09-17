---
title: "CLI — doctl (DigitalOcean CLI)"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/digitalocean, topic/infra]
ai_summary: "DigitalOcean CLI (doctl) quick reference — droplet management for onemind-main (104.248.4.212), snapshots, firewall, and domains."
---

# doctl — DigitalOcean CLI

> Manages the VPS `onemind-main` (104.248.4.212). Install: `brew install doctl`.

---

## Auth

```bash
# Authenticate (uses DO API token)
doctl auth init

# Switch context
doctl auth switch --context onemind

# Verify auth
doctl account get
```

---

## Droplets (VPS)

```bash
# List all droplets
doctl compute droplet list --format Name,PublicIPv4,Status,Region,Size

# Get info on onemind-main
doctl compute droplet list --format Name,PublicIPv4,Status,Memory,VCPUs,Disk | grep onemind

# Reboot droplet
doctl compute droplet-action reboot --droplet-id <id>

# Power off / on
doctl compute droplet-action power-off --droplet-id <id>
doctl compute droplet-action power-on --droplet-id <id>

# Resize droplet (power off first)
doctl compute droplet-action resize --droplet-id <id> --size s-4vcpu-8gb --resize-disk false
```

### onemind-main specs

| Field | Value |
|-------|-------|
| Name | onemind-main |
| IP | 104.248.4.212 |
| Tailscale IP | 100.77.52.93 |
| CPU | 4 vCPUs |
| RAM | 16 GB |
| Disk | 193 GB |
| Provider | DigitalOcean |

---

## Snapshots & Backups

```bash
# Create snapshot (droplet must be off for consistent snapshot)
doctl compute droplet-action snapshot --droplet-id <id> --snapshot-name "onemind-backup-$(date +%Y%m%d)"

# List snapshots
doctl compute snapshot list

# Restore from snapshot
doctl compute droplet create --name restored-droplet --snapshot-id <snapshot-id> --region nyc3 --size s-4vcpu-8gb
```

---

## Firewall

```bash
# List firewalls
doctl compute firewall list

# Get rules for a firewall
doctl compute firewall get <firewall-id>

# Add inbound rule
doctl compute firewall add-rules <firewall-id> \
  --inbound-rules "protocol:tcp,ports:443,address:0.0.0.0/0"
```

---

## SSH Keys

```bash
# List SSH keys in account
doctl compute ssh-key list

# Add new SSH key
doctl compute ssh-key create "new-key" --public-key-file ~/.ssh/id_rsa.pub
```

---

## Domains & DNS

```bash
# List domains
doctl compute domain list

# List DNS records for domain
doctl compute domain records list onemindos.com

# Create DNS record
doctl compute domain records create onemindos.com \
  --record-type A \
  --record-name api \
  --record-data 104.248.4.212 \
  --record-ttl 300
```

---

## Databases (Managed)

```bash
# List managed databases
doctl databases list

# Get connection info
doctl databases connection <db-id>
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
