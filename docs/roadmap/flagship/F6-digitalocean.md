# F6 — DigitalOcean at a glance (custom plugin)

Work order for one implementing agent. Branch: `flagship/f6-digitalocean`.
Program index: [README](README.md).
**Blocker:** needs a read-only DigitalOcean API token that Zeus creates (see step 1).

## Outcome Zeus sees

- **DigitalOcean page** (`/digitalocean`), one screen:
  - month-to-date spend, projected month total, last invoices
  - DOKS cluster and node pools: size, count, version, upgrade available
  - Droplets: status, size, region, tags
  - Volumes: size, attached to which PVC/app, unattached (orphans)
  - Snapshots: count, size, age, the oldest and the largest
  - Load balancers, firewalls, VPCs
  - Spaces buckets: size, object count, if the API allows
  - Domains and DNS records, uptime checks and their state
  - Monitoring alert policies
- **In the catalog**: every DO resource is a catalog entity under
  `resource:default/do-*`, linked to the app that uses it. The Argo CD app page
  shows "runs on DOKS node pool X; data on volume Y; backed up to Space Z".
- **Change log**: a "What changed" list (new or removed resources since the
  previous refresh).

## Evidence

- No DigitalOcean plugin exists in the Backstage plugin directory (checked
  2026-10-21). This must be custom.
- `/today` shows a DO cost card waiting for a `billing:read` token. Secret
  `backstage-billing` does not exist (checked 2026-11-04).
- Hatch credits hide the real bill: October 2026 usage was $776.67 and the invoice
  was $0.03. The page must show usage before credits as well as the amount due.
- Read-only inventory on 2026-11-04 (`doctl`): 1 DOKS cluster `omos-cloud` with
  4 droplets (3× s-8vcpu-16gb, 1× s-8vcpu-32gb-amd), 73 volumes, 53 snapshots,
  1 load balancer, 3 firewalls, 4 VPCs and 1 uptime check. No managed databases
  or DOCR registry.

## Build

1. **Token (Zeus).** Create a DigitalOcean custom-scoped token with read-only
   scopes only (account, billing, droplet, kubernetes, block_storage, snapshot,
   load_balancer, firewall, vpc, domain, monitoring, uptime, spaces/read where
   offered). Zeus seals it as SealedSecret `backstage-digitalocean`,
   key `DIGITALOCEAN_TOKEN`, namespace `backstage`. The agent writes the
   manifest reference and waits; it never handles the token value.
2. **Backend plugin** `plugins/plugin-onemind-digitalocean-backend`:
   - a typed DO API client (fetch, pagination, rate-limit backoff, 30s timeouts)
   - a catalog `EntityProvider` (refresh every 10 minutes) emitting Resources
     with `metadata.title`, `spec.type` (`do-droplet`, `do-volume`,
     `do-snapshot`, `do-load-balancer`, `do-firewall`, `do-vpc`, `do-bucket`,
     `do-kubernetes-cluster`, `do-node-pool`), stable names from DO IDs, and
     links to cloud.digitalocean.com
   - PVC linking: map `pvc-<uid>` volume names to Kubernetes PVCs (read-only
     API) and then to the owning app (F2 mapping)
   - a summary route `/api/onemind-digitalocean/summary` for the page and the
     billing data
   - permission-guarded (admins). The token never reaches the browser
     (negative test).
3. **Frontend plugin** `plugins/plugin-onemind-digitalocean`: page plus an
   entity card "Runs on DigitalOcean" for apps. Spend appears on Home (F1)
   status strip and the Cost page.
4. **Orphan/retention signals** feed F8: unattached volumes, snapshots older
   than the retention policy, idle LBs. These are reports only, with no deletes.

## Out of scope

Creating, resizing or deleting any DO resource. The page is read-only.

## Files

`plugins/plugin-onemind-digitalocean{,-backend}/**` (new),
`packages/{app,backend}` registration, `app-config.yaml` (config schema, no
secret values), `platform/backstage/deployment.yaml` (env from the SealedSecret,
after Zeus provides it), tests and `docs/operations/digitalocean.md`.

## Acceptance

- [ ] With a recorded fixture of real (redacted) API responses, the provider
      emits the expected entity counts and stable names, and is idempotent
      across refreshes.
- [ ] Live, after the token exists: counts match `doctl` for droplets, volumes,
      snapshots and LBs; spend matches the DO billing page.
- [ ] Every PVC-backed volume links to its PVC and app.
- [ ] Token absent → page shows "not connected", nothing crashes.
- [ ] Required gates pass.

## Status

2026-11-05, branch `feat/f6-digitalocean` (not yet merged). Built and tested with the recorded fixture;
not yet checked against the live account.

- Done: backend plugin `plugin-onemind-digitalocean-backend` (read-only client with pagination, 429
  backoff, size and time limits; summary route for administrators; catalog provider with PVC and app
  links; "What changed" list), frontend plugin `plugin-onemind-digitalocean` (page at `/digitalocean`
  and the "Runs on DigitalOcean" card), registration in `packages/app` and `packages/backend`,
  `onemind.features.digitalocean: true` so the sidebar entry shows, `DIGITALOCEAN_TOKEN` and
  `DO_BILLING_READ_TOKEN` in `platform/backstage/deployment.yaml`, and
  [docs/operations/digitalocean.md](../../operations/digitalocean.md). Owner sealed
  `backstage-digitalocean`.
- Spaces buckets: the DigitalOcean API v2 cannot list buckets with an API token, so the page says
  "not available".
- Open: live acceptance (counts match `doctl`, spend matches the DigitalOcean billing page), node pool
  per app (reached through the cluster Resource, not per pod), spend on the Home status strip (F1).
