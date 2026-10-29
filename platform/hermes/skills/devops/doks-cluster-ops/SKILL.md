---
name: doks-cluster-ops
description: "Use for DOKS cluster workups, crashloops, disk and PVC triage, right-sizing, Argo CD drift and decommission sweeps."
category: devops
version: 2.0.0
tags: [doks, kubectl, kubernetes, argocd, pvc, ops]
triggers:
  - cluster health
  - crashloop
  - Released PV
  - right-size
  - Argo degraded
---

# DOKS Cluster Ops

Cluster `omos-cloud` (nyc3). Facts that drift (node counts, Application counts,
which Helm releases are adopted) live in `docs/STACK.md` and `CLAUDE.md`; verify
against the cluster before saying "live".

## Standing gate

No manifest edit, resource change, deletion or cluster mutation without explicit
authorization. A direct instruction to deploy a named application authorizes that
scoped, non-destructive deployment and its ordinary prerequisites; do not ask for
the same approval again. Ask before broadening scope (a different broker or
backend, migrating other workloads, anything destructive). For unrequested changes,
present the plan (what, why, savings, risk) first. After every authorized removal,
verify live state and GitOps ownership so nothing respawns through self-heal.

This agent's in-cluster account is read-only. Reads are always fine; changes go
through a PR to `platform/` and Argo CD.

## Git is the only source of truth

- Never `kubectl apply` a change to an Argo-managed resource: the next sync
  reverts it silently. Edit the manifest under `platform/`, open a PR, let Argo sync.
- `kubectl create secret` in an Argo-managed namespace is pruned. Secrets are
  SealedSecrets committed beside the workload.
- Three Helm releases are still CLI-only: `nats`, `takserver`, `argocd`. A Git
  commit does not deploy them. Check `helm list -A` before concluding "the fix
  didn't work". Never `helm upgrade` an adopted release, and do not run
  `platform/install.sh`. See `docs/gitops/HELM-ADOPTION.md`.
- Do not trigger a sync by patching an empty revision; patch the operation
  with the revision from `.status.sync.revision`, or use the Argo CD UI/CLI.
- `RepeatedResourceWarning` means a resource is declared twice under one source
  path (a full manifest plus a patch). Remove the duplicate. A Degraded Application
  whose resources are all `health=None, sync=Synced` is a controller cache
  artifact; verify pods directly.

## Right-sizing

Every workload needs requests and limits. The hub-core pool is memory-saturated, so
do not add requests casually.

- Measure first (`kubectl top pod --containers`), then set the request near 60 to
  90 percent of measured and the limit at measured plus 15 to 25 percent.
- Never set a memory limit at or below measured usage: the container is OOM-killed
  at once and the crashloop looks like an app bug. After lowering a limit, watch for
  `OOMKilled` in `lastState`.
- Over-requesting is invisible in `kubectl top` but starves the scheduler. A
  request-to-actual ratio above 2 is a right-size candidate.
- Change values only in Git.

## Storage

- Stateful claims use `do-block-storage-retain`. Temporary claims, VolSync caches
  and restore drills use `do-block-storage`. A claim's class cannot change in place.
- Postgres on a fresh block volume needs `PGDATA` in a subdirectory (the mount
  root contains `lost+found`): `references/k8s-stateful-app-triage.md`.
- Released PVs still bill. After a decommission list them (`kubectl get pv | grep Released`,
  `doctl compute volume list`) and report; deleting a volume needs a fresh snapshot
  and the owner's approval.
- Single-replica Deployments holding an RWO claim should use `strategy: Recreate`;
  see `references/rwo-pvc-stale-rs-deadlock.md`.

## Autonomous repair CronJob

The `ops` namespace health-keeper (Argo `platform-ops`) deletes non-running and
CrashLoop pods so controllers respawn them, and publishes `ops.health.escalate`
for conditions it must not touch (a Released PV, a node above 90 percent). Check
`kubectl -n ops logs -l app=health-keeper --tail=200` before debugging a
crashloop by hand, and treat escalate messages as work items for the owner.

## Public and admin access

Public hostnames use ingress class `nginx` with the Cloudflare allowlist; admin UIs
use class `tailscale` and are tailnet-only. Real tailnet hostnames come from the
Ingress status, never from the proxy Service name:

```bash
kubectl get ingress -A -o json | jq -r '.items[] | select(.spec.ingressClassName=="tailscale") | "\(.metadata.namespace)/\(.metadata.name) -> https://\(.status.loadBalancer.ingress[0].hostname)"'
```

Never mint long-lived cluster-admin tokens or write tokens into pod files.

## Load balancer consolidation

Each `type: LoadBalancer` Service with a DO LB bills about $12 a month. To merge
many into one nginx `stream` (raw TCP, no TLS termination): add the upstreams and
ports in Git first, let Argo sync, verify every upstream from inside the hub
(`nc -z`) before deleting the old Services, and tell the owner which client
addresses change. Deleting any Service in a consolidation is an owner-confirm
action; present the merge plan and wait.

## Decommissioning a component

1. Establish provenance first (image, env fingerprints, source repo, upstream
   compose list); an image under `ghcr.io/onemindos/` can still be someone else's
   code, and an empty database in a "their" pod means it is ours. Never cut a data
   store on inference.
2. Sweep in the same commit: Ingresses pointing at it (a dead backend behind a live
   Ingress is a silent 503), env variables in kept workloads that name its DNS,
   ExternalName Services, SealedSecrets, ConfigMaps, PDBs and install-script lines.
3. Verify after sync: pods pruned, kept pods Ready, Application Synced/Healthy,
   re-routed hosts answer.

## Authorized deployment: acceptance

1. Complete a backend and feature eligibility check before scaffolding: pin the
   release, identify the real required services and license or feature gates, and
   probe the documented API on the intended backend. Cloud docs, a changelog
   entry or a landing page do not prove self-hosted compatibility. Never add a
   dummy backend or suppress a dependency to make health green.
2. Limit discovery output to identity, image, ports, ownership and readiness.
   Redact credential-valued arguments and env values before they reach stdout.
3. Deploy through the real GitOps owner and run the repository's schema and
   inventory checks. Verify the namespace against live admission policies.
4. Read back the reconciled revision, readiness, endpoints and ingress hostname,
   then load the UI and exercise the requested feature. A pushed commit and a root
   HTTP 200 are intermediate states, not "done".
5. Report the verified URL and feature status, or the exact blocker and last
   verified state. Distinguish manifest, merged desired state, reconciled
   resources, rendered UI and working feature.

## References

- `references/k8s-stateful-app-triage.md`: Postgres, Redis and read-only ConfigMap crashloops.
- `references/rwo-pvc-stale-rs-deadlock.md`: Multi-Attach deadlocks on rollout.
- `references/ghcr-image-rollout.md`: build, pull secret and verify a new GHCR image.
