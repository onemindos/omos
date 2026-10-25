# Vaultwarden — personal private vault

Argo Application `platform-vaultwarden` deploys one Vaultwarden 1.37.4 replica,
SQLite on a protected 5Gi `data-pvc`, and private Tailscale HTTPS at
**https://vaultwarden.tail717ba0.ts.net**. Source preparation is not live acceptance.
The official image index digest is pinned in `deployment.yaml`. No community chart,
PostgreSQL dependency, public ingress or agent vault access is introduced.

## Access and account enrollment

Public registration and user invitations are disabled, including domain allowlists
and SSO signup. SMTP is unset. The separate administrator token is stored only as
an Argon2id hash in a strict namespace/name-bound SealedSecret. Its recovery token
is kept privately outside Git. It is not a vault master password.

Once Zeus confirms his account email, an administrator can reserve that exact email
using `/admin/invite` without SMTP. Zeus then opens `/#/signup` privately with the
invited email and chooses his own master password. The invitation permits this one
email despite closed registration. Without SMTP this reserves an email but does
not verify mailbox ownership; complete enrollment promptly on the private route.
Never open general registration or create a master password for the owner.

Use official Bitwarden clients with their self-hosted server setting pointing to
this HTTPS address. Website passkey storage, client autofill, desktop/mobile and
browser acceptance require the owner's first login and actual device checks.
Vaultwarden 1.37.4 does not implement logging into the vault with a passkey.
No 1Password export/import, AI credential grant, CLI/MCP enrollment, mobile push
registration or Secrets Manager compatibility is claimed by this deployment.

## Runtime and network

One non-root writer, Recreate updates and RWO storage protect SQLite. The container
has a read-only root, no API token, dropped capabilities, UID/GID 1000 and seccomp.
Startup/readiness/liveness use `/alive`, whose tagged source acquires a DB connection.
Requests are 100m CPU/128Mi memory/128Mi ephemeral storage; limits are 1 CPU/512Mi
memory/2Gi ephemeral storage. `/tmp` has a 1Gi disk-backed upload staging allowance.
These are initial budgets, not usage measurements. The 5Gi claim also holds online
backup copies; watch free space before importing large attachments.

Default-deny policy admits only this Ingress's Tailscale proxy on 8080; DNS and public
IPv4 HTTPS are allowed. Private, metadata and tailnet egress are excluded. This
keeps optional favicon/HIBP HTTPS working without allowing private targets. No
SMTP, external push or IPv6 egress is configured. Trusted forwarded-header proxies
are restricted to the observed Cilium pod CIDRs' aggregate 10.109.2.0/23; update via
Git if cluster addressing changes. Network policy further restricts actual senders.

`/data/config.json`, if later saved through the admin UI, overrides editable env
settings. Keep configuration in Git and do not save admin UI configuration. Audit
for this file before diagnosing configuration drift.

## Acceptance, backup and recovery

Verify Argo Synced/Healthy, sealed-secret synchronization, bound claim, exact
image, private TLS, `/alive`, and `/api/config` reporting disabled registration.
A healthy server does not prove owner enrollment, real client sync or recovery of
populated user data. Record those checks separately after owner login.

See [BACKUP.md](BACKUP.md) for native SQLite online copies, encrypted off-cluster
snapshots and isolated restore. Never delete namespace/PVC or restore over live
SQLite. Before upgrades, make and restore-check a fresh backup, review upstream
migrations, and preserve keys/attachments with the database. An image rollback does
not undo schema migrations. Restore first to an isolated new claim.

## Verified upstream references

- [Tagged official Docker build](https://github.com/dani-garcia/vaultwarden/blob/1.37.4/docker/Dockerfile.debian)
- [Tagged configuration](https://github.com/dani-garcia/vaultwarden/blob/1.37.4/.env.template)
- [Native backup command](https://github.com/dani-garcia/vaultwarden/wiki/Backing-up-your-vault)
- [Admin invitations with closed registration](https://github.com/dani-garcia/vaultwarden/wiki/Disable-registration-of-new-users)
- [Tailscale Kubernetes HTTPS ingress](https://tailscale.com/docs/kubernetes-operator/ingress/expose-workload-to-tailnet-l7)

Local preflight ran the exact amd64 digest as UID1000, read-only root with dropped
capabilities; `/alive`, closed-registration config, and `/vaultwarden backup` passed.
SQLite backup integrity is checked independently. No real owner credentials or
vault items are needed for this preflight.
