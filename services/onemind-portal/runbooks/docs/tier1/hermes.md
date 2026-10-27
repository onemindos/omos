# Hermes

Hermes "Legacy" is the platform's AI agent: Deployment `hermes-legacy` in namespace `hermes`, on the `nousresearch/hermes-agent` image pinned in Git, with a `nats-leaf` sidecar that gives the agent its own bus identity. All of its identity and state is the 25 GiB PVC `hermes-data`, mounted at `/opt/data`. A separate, vanilla `hermes-fresh` instance launched on October 26, 2026 runs beside it with its own volume and no imported Legacy data.

Sources: `platform/hermes/agent/RUNBOOK.md`, `platform/hermes/agent/README.md`, `platform/hermes/fresh/README.md`, `docs/operations/hermes-recoverability-plan.md`, `docs/operations/hermes-restore-drill-2026-10-16.md`.

## Ownership and access

- `platform-legacy` applies the Legacy agent and its sealed bridge secrets; `platform-hermes-fresh` applies `platform/hermes/fresh`.
- Hermes always connects to `nats://localhost:4222`, where the sidecar runs a leaf node configured from the Secret `hermes-leaf-conf`. Do not copy this pattern into new services; they connect to the cluster NATS directly.
- The pod cannot change the cluster: at most it holds the read-only `hermes-observer` identity, which cannot read Secrets, ConfigMaps or logs. Run cluster commands from an admin workstation.
- The dashboards are tailnet-only.

## Investigate (read-only)

```sh
kubectl --context do-nyc3-omos-cloud -n argocd get application platform-legacy platform-hermes-fresh
kubectl --context do-nyc3-omos-cloud -n hermes get deployments,pods,pvc
kubectl --context do-nyc3-omos-cloud -n hermes logs deploy/hermes-legacy -c hermes --tail=50
```

Inside the pod, `nats account info` shows the leaf's account; "no responders" means the sidecar is down.

## Recovery

| Situation | Documented action |
| --- | --- |
| Pod crashed or restarted | Nothing is lost, because `hermes-data` persists. If the leaf sidecar is down, restart the pod, but never while Hermes' state-database recovery or integrity scan is running |
| Volume lost or corrupted | Restore the latest VolSync snapshot of `hermes-data` into a new PVC, verify it, then point the Deployment at it through Git. An application-level restore drill passed on October 16, 2026 |
| Cluster lost | Rebuild the cluster from `infra/`, let Argo CD sync `platform/`, then restore the volume as above |
| Rebuild identity from Git alone | Not available yet: the soul, skills and configuration are not exported to Git (recoverability plan, step 2) |
| After a restore | Log `gh` in again with a narrow token, confirm the SSH keys under `/opt/data/home/.ssh/` and check the leaf connection |

Never put a token in a git remote URL, in this repository or in Hermes' memory. Hermes' standing rules forbid destroying data or deleting files without the owner's explicit go-ahead.

## Verify completion

The pod is Ready, the leaf connection shows its account, and a real prompt gets a real answer; a Ready pod is not proof of working chat.

Generated workload list: [Hermes (Legacy) app runbook](../apps/hermes-legacy.md).
