# DigitalOcean page in the portal

The portal page `/digitalocean` shows what the DigitalOcean account runs and what it costs. It is
**read-only**: the code only sends GET requests, and the page has no delete, resize or create
button. Only members of `group:default/onemind-admins` can open it.

## What the page shows

- **Spend**: usage before credits (DigitalOcean `month_to_date_usage`), amount due after credits
  (`month_to_date_balance`), the projected month (usage before credits, scaled by days in the
  month over days elapsed) and the last invoices, with the open invoice for this month. Credits can
  hide real usage (October 2026: usage $776.67, invoice $0.03), so both numbers are shown.
- **Kubernetes cluster and node pools**: version, size, count, autoscale range, and whether DOKS
  offers an upgrade.
- **Droplets**, **load balancers** (idle ones flagged), **firewalls** (rule counts only, no
  addresses), **VPCs**, **domains** (with record counts), **uptime checks** (with state) and
  **monitoring alert policies**.
- **Volumes**: size, the Kubernetes claim and app that use it, and an "Unattached" filter and badge.
  Unattached volumes are a report for review. Nothing deletes them.
- **Snapshots**: count, total size, the oldest, the largest and how many are older than 90 days.
- **Spaces buckets**: shown as "not available". The DigitalOcean API cannot list buckets with an API
  token, so the portal does not try.
- **What changed**: resources added or removed since the previous refresh (kept in memory since the
  portal started).

Each section fails on its own. A token without one scope gives "Unavailable" for that section only.
Without a token the page says **Not connected: needs the read-only DigitalOcean token**.

Data is cached for 10 minutes. Reload the page for the cached copy; the catalog provider refreshes
the data every 10 minutes.

## Catalog

When the token exists, every droplet, volume, snapshot, load balancer, firewall, VPC, DOKS cluster
and node pool is a catalog `Resource` named `do-<type>-<id>` (`spec.type` is `do-droplet`,
`do-volume`, `do-snapshot`, `do-load-balancer`, `do-firewall`, `do-vpc`, `do-kubernetes-cluster`,
`do-node-pool`), with a link to cloud.digitalocean.com. They come from a catalog provider, not from
files, so they do not use the `backstage-catalog` ConfigMap. Without the token nothing is published.
If any resource section cannot be read, the provider keeps the entities it already has.

A volume named `pvc-<uid>` is matched to its Kubernetes claim (read-only claim and pod lists), then
to the app whose Deployment or StatefulSet mounts it. The volume `dependencyOf` that app, so the
component page shows a **Runs on DigitalOcean** card with its data volumes and, through the cluster
Resource, the cluster and node pools. Claims whose workload is not in the catalog show the claim but
no app.

## Token

Create a **custom-scoped** token in DigitalOcean (API, Tokens) with read scopes only:
`account:read`, `billing:read`, `kubernetes:read`, `droplet:read`, `block_storage:read`,
`image:read` (snapshots), `load_balancer:read`, `firewall:read`, `vpc:read`, `domain:read`,
`monitoring:read`, `uptime:read`. Do not use the `doctl` token or a full-access token. Scope names are
as DigitalOcean shows them in the token form.

Seal it on your Mac (needs `kubectl` access to `do-nyc3-omos-cloud` and `kubeseal`):

```bash
python3 services/onemind-portal/scripts/seal-digitalocean.py
```

The script reads the token hidden, checks the read endpoints, and writes only the encrypted
`platform/backstage/digitalocean.sealedsecret.yaml` (Secret `backstage-digitalocean`, key
`DIGITALOCEAN_TOKEN`). Commit that file; Argo CD applies it. The Deployment reads the key as
`DIGITALOCEAN_TOKEN` and also as `DO_BILLING_READ_TOKEN`, which feeds the cost card on `/today`. Both
are optional, so the portal starts without the Secret. The token is read from the environment only,
is never returned by any route, and is never logged.

## Rotate

1. Create a new read-only token, then run the seal script again with the new value.
2. Commit the changed sealed secret and let Argo CD sync.
3. Restart the portal pod (`kubectl -n backstage rollout restart deploy/backstage`): the token is
   read from the environment, so a pod restart picks it up.
4. Delete the old token in DigitalOcean.

## Tests

`services/onemind-portal/plugins/plugin-onemind-digitalocean-backend/src/__fixtures__/do-api.json` is a
recording of the real account (read-only `doctl` list and get commands, 2026-11-05) with ids replaced
by stable fake ids and IP addresses, email, team, firewall sources and the uptime target removed.
Uptime state is written to the DigitalOcean API shape because `doctl` has no command for it. The
domain and alert-policy lists are empty on the live account.
