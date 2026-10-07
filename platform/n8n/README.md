# n8n

Self-hosted workflow capability, managed by Argo application platform-n8n.

- n8n and external distroless runners pinned to 2.42.4 and image digests.
- Dedicated PostgreSQL 17 with an application role without cluster administration privileges.
- 896 MiB aggregate requested memory and 400m requested CPU. This is a bounded initial deployment, not a validated full-scale capacity estimate.
- Persistent n8n state (5 GiB) and database (10 GiB) use Retain storage classes.
- Independent encryption key, database passwords and runner token are SealedSecrets. Back up the SealedSecrets recovery key separately.
- Code runners have no application secret mounts, no shared data volume, a read-only root, and non-root identity. Containers share pod networking; this is not a complete untrusted-code boundary.
- No Kubernetes service-account token, shell command node, local file node, or community packages by default.
- Public HTTP/HTTPS egress allowed; private APIs need explicitly reviewed NetworkPolicy additions. No cluster-wide AI credentials are installed.
- Seven-day execution history and five concurrent production executions initially.

## Access and owner enrollment

Private editor URL: https://n8n.tail717ba0.ts.net/. Requires access to the existing tailnet; the ingress is not public and does not use Funnel.

Owner enrollment completed for zeus@zeusdelacruz.com through a localhost-only port-forward before ingress was enabled. The generated initial password is stored outside Git in /Users/zeus/.config/omos/n8n-owner-20261007.json with permissions 0600. Move it to the operator's password manager, change it after first login and enable MFA. Do not place it in Git or logs. Native n8n login is used; enterprise SSO is not assumed available in the free edition.

The URL is initially tied to the existing tailnet. A future mesh migration needs a new private DNS/TLS plan and editor/webhook URL changes. Webhooks are private; external SaaS callbacks need a separately reviewed ingress design.

## Recovery and scaling

Retained volumes protect against automatic volume deletion, not node/provider disaster or database corruption. External encrypted database backups and a restore exercise are still required before critical workflows depend on this instance. Preserve N8N_ENCRYPTION_KEY with the database; losing it makes stored credentials unreadable.

Queue mode is a future scale step: add supported Redis, worker deployments and one matching runner per worker, and account for shared binary storage. The current single instance is not HA.

Export reviewed workflows into a Git-managed workflow directory, preserving credentials as references. Argo reconciles application infrastructure; it does not automatically reconcile workflows edited inside n8n.

Sources: https://docs.n8n.io/deploy/host-n8n/configure-n8n/set-up-task-runners.md and https://docs.n8n.io/deploy/host-n8n/configure-n8n/security/harden-task-runners.md.
