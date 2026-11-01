# Tier 1 services

Tier 1 is the owner's list of services without which the platform cannot operate or cannot be recovered: TAK Server, NATS, Authentik, ingress, cert-manager, Sealed Secrets, backups and ops, Argo CD and Hermes (`docs/STACK.md`, "Tier definitions"). A tier 1 change goes behind a backup check, dependency bumps are not auto-merged, and each service needs a restore test on record.

Each page condenses the Git-owned documents named under "Sources" and was checked against `main` and read-only cluster observations on October 27, 2026. The pages record configuration and evidence, not live health: check the app page's Live tab and Argo CD before acting. Some cited sections are older than that evidence, for example parts of the root `CLAUDE.md` and `docs/STACK.md`. Where a page and a source disagree, the newer verified evidence wins: confirm it read-only, then correct whichever document is stale.

| Service | Runbook | Catalog entity | Argo CD applications |
| --- | --- | --- | --- |
| TAK Server | [TAK Server](tak-server.md) | `component:default/tak-server` | `platform-takserver` (plus the legacy Helm release `takserver`) |
| NATS | [NATS](nats.md) | `component:default/nats` | `platform-nats`, `edge-nack` (plus the Helm release `nats`) |
| Authentik | [Authentik](authentik.md) | `component:default/authentik` | `platform-auth-helm`, `platform-auth-infra` |
| Ingress | [Ingress (ingress-nginx)](ingress-nginx.md) | `component:default/ingress-nginx` | `edge-ingress-nginx` |
| cert-manager | [cert-manager](cert-manager.md) | `component:default/cert-manager` | `edge-cert-manager` |
| Sealed Secrets | [Sealed Secrets](sealed-secrets.md) | `component:default/platform-core` | `platform-core` |
| Backups and ops | [Backups and ops](backups-ops.md) | `component:default/platform-backup-recovery`, `volsync`, `velero`, `health-keeper` | `platform-backups`, `platform-volume-backups`, `platform-volsync`, `ecosystem-velero`, `platform-ops` |
| Argo CD | [Argo CD](argocd.md) | `component:default/argocd` | `platform-argocd`, `app-of-apps` |
| Hermes | [Hermes](hermes.md) | `component:default/hermes-legacy` | `platform-legacy`, `platform-onemind-hermes`, `platform-hermes-base` |

`docs/STACK.md` also proposes the Tailscale operator for tier 1, because every admin UI is reachable only through it. It has no page here until that proposal is accepted; see its generated page, [Tailscale operator and proxies](../apps/tailscale-operator.md).

## Rules for every tier 1 service

- Select and confirm context `do-nyc3-omos-cloud` first. The commands on these pages only read.
- Change the cluster through a reviewed pull request and Argo CD. Never run `helm upgrade`, `helm rollback` or `helm uninstall` on a release that Argo CD has adopted, and never run `platform/install.sh` on this cluster (`docs/gitops/HELM-ADOPTION.md`).
- Removing or renaming an Application file can delete everything that Application tracks, including a whole namespace and its volumes. Check the live finalizer and tracked resources first ([Argo CD](argocd.md)).
- Before a change that touches storage, confirm the latest verified backup ([Backups and ops](backups-ops.md)). Never delete PVCs, namespaces or backup data without a fresh snapshot and explicit approval.
- Never copy Secret values, keys or passwords into chat, Git, tickets or logs.
- Recovery is complete only when the affected user workflow works, not when pods are Ready.
