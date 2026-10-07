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

## Initial owner enrollment

No ingress is enabled while the owner account is unclaimed. Use an authenticated local Kubernetes port-forward:

```sh
kubectl -n n8n port-forward service/n8n 5678:5678
```

The intended final editor/webhook URL is https://n8n.tail717ba0.ts.net/. Browser setup should use HTTPS; complete owner enrollment through a local API call with the owner's chosen email and a generated secret before adding the Tailscale ingress. Do not turn off authentication or secure cookies to work around enrollment.

Owner email is pending user input. Credentials must not be placed in Git, logs, or tool output. Native n8n login is initially used; do not assume enterprise SSO is included in the free edition.

## Recovery and scaling

Retained volumes protect against automatic volume deletion, not node/provider disaster or database corruption. External encrypted database backups and a restore exercise are still required before critical workflows depend on this instance. Preserve N8N_ENCRYPTION_KEY with the database; losing it makes stored credentials unreadable.

Queue mode is a future scale step: add supported Redis, worker deployments and one matching runner per worker, and account for shared binary storage. The current single instance is not HA.

Export reviewed workflows into a Git-managed workflow directory, preserving credentials as references. Argo reconciles application infrastructure; it does not automatically reconcile workflows edited inside n8n.

Sources: https://docs.n8n.io/deploy/host-n8n/configure-n8n/set-up-task-runners.md and https://docs.n8n.io/deploy/host-n8n/configure-n8n/security/harden-task-runners.md.
