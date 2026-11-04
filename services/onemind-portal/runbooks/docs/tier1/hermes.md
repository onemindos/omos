# Hermes

Hermes "Legacy" is the platform's AI agent: Deployment `onemind-hermes` in namespace `hermes`, on the official Hermes image pinned by digest in Git. It is a direct client of the NATS bus (no leaf sidecar) and holds the fabric identity `agents.*.hermes.zeus.legacy`. All of its identity and state is the 25 GiB PVC `onemind-hermes-data` (`do-block-storage-retain`), mounted at `/opt/data`. The older `hermes-legacy` instance, with its `nats-leaf` sidecar and the volume `hermes-data`, was removed from Git on 2026-11-03.

Sources: `platform/hermes/onemind-hermes/README.md`, `platform/hermes/onemind-hermes/MIGRATION.md` (historical), `platform/hermes/onemind-hermes/RETIREMENT-AUDIT.md`, `platform/hermes/base/README.md`, `platform/hermes/doctrine/` (history only).

## Ownership and access

- `platform-hermes-base` (`platform/hermes/base`) owns the `hermes` Namespace, the Tailscale egress Services and the Secrets `onemind-hermes` shares (`hermes-dashboard-auth`, `omos-nats-creds`, `hermes-nats-ca`); `platform-onemind-mcp` (`platform/hermes/mcp`) owns `onemind-mcp`; `platform-onemind-hermes` applies `platform/hermes/onemind-hermes`. All three have no finalizer and `prune: false` (the MCP Application prunes only its hashed ConfigMaps), and the long-lived objects carry `argocd.argoproj.io/sync-options: Prune=false,Delete=false`. See [Argo CD](argocd.md).
- Hermes connects to the cluster bus with `NATS_USER` / `NATS_PASSWORD` from `omos-nats-creds` and the CA from `hermes-nats-ca`. Do not add a sidecar or a leaf.
- The pod has no Kubernetes API token (`automountServiceAccountToken: false`). Cluster reads and writes go through the tool gateway groups. Argo CD deploys whatever reaches `main`, and GitHub does not enforce branch protection on this repository (`hermes-recoverability-plan.md`, section 2.3), so keep Git write credentials off the pod.
- The dashboard is tailnet-only at `onemind-hermes.tail717ba0.ts.net`.

## Investigate (read-only)

```sh
kubectl --context do-nyc3-omos-cloud -n argocd get application platform-onemind-hermes platform-hermes-base platform-onemind-mcp
kubectl --context do-nyc3-omos-cloud -n hermes get deployments,pods,pvc
kubectl --context do-nyc3-omos-cloud -n hermes logs deploy/onemind-hermes -c hermes --tail=50
```

`platform/hermes/onemind-hermes/nats/apply.sh --status` reports the plugin files, the enabled list and which NATS variables the pod has (names only).

## Recovery

| Situation | Documented action |
| --- | --- |
| Pod crashed or restarted | Nothing is lost, because `onemind-hermes-data` persists |
| Volume lost or corrupted | Restore the latest VolSync snapshot of `onemind-hermes-data` (its own restic prefix) into a new PVC, verify it, then point the Deployment at it through Git |
| Cluster lost | Rebuild the cluster from `infra/`, let Argo CD sync `platform/`, then restore the volume as above |
| Rebuild identity from Git alone | The soul and the OneMind skills are seeded from `platform/hermes/seed/` and `platform/hermes/skills/`; config defaults and cron jobs are in `platform/hermes/onemind-hermes/` |
| NATS identity shows "live elsewhere" | Another process holds `agents.*.hermes.zeus.legacy`. Find it with `/connz`; do not start a second copy |

Never put a token in a git remote URL, in this repository or in Hermes' memory. Hermes' standing rules forbid destroying data or deleting files without the owner's explicit go-ahead.

## Verify completion

The pod is Ready, the NATS identity is registered (Agent Console shows Legacy live), and a real prompt gets a real answer; a Ready pod is not proof of working chat.

Generated workload list: [Legacy (onemind-hermes) app runbook](../apps/onemind-hermes.md).
