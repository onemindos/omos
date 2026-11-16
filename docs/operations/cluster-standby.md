# Cluster standby (Hermes only) and wake

**Why:** during the move to OneMind OS, only Hermes (Legacy) is used. The cluster cost about $510 a month on
2026-11-16 (3 x `s-8vcpu-16gb` at $96, 1 x `s-8vcpu-32gb-amd` at $168, about $42 of volumes, $12 load balancer).
Standby keeps Hermes running and everything else **paused, not deleted**, on one node: about **$220 a month**.
Wake brings it all back in about 20 minutes.

**Tool:** `scripts/ops/cluster_standby.py` (tests: `scripts/tests/test_gitops_cluster_standby.py`).
Every change records the original value in an `onemind.io/standby-*` annotation on the object itself, so wake
restores exactly what was there. Without `--yes`, each command only prints its plan.

## What stays running

| Namespace | Why |
|-----------|-----|
| `hermes` | Hermes (Legacy) and `onemind-mcp` |
| `tool-gateway` | Hermes tool calls |
| `tailscale` | Hermes dashboard (`onemind-hermes.tail717ba0.ts.net`) and SSH to devices |
| `auth`, `onemind-data` | Authentik login for the Hermes dashboard; its database is in the CNPG cluster `onemind-data` |
| `cnpg-system`, `kyverno`, `cert-manager`, `ingress-nginx` | Admission webhooks with `failurePolicy: Fail`: if these stop, no pod can be created anywhere |
| `argocd` | Needed to wake the rest |
| `github-runner`, `ops`, `velero`, `falco`, `kube-system` | Runner, backups, runtime security, the cluster itself |

Everything else (27 namespaces on 2026-11-16: TAK, OneTAK, monitoring, NATS, ClickHouse, Langfuse, n8n, Karakeep,
Vaultwarden, Forgejo, OpenRemote, ThingsBoard and others) goes to 0. Use `--keep <namespace>` to keep one more.

## What you lose while on standby

- **Hermes on NATS.** Telegram, the dashboard, Hermes Desktop and cron keep working, but `agents.*.hermes.zeus.legacy`
  is offline and Hermes logs a NATS reconnect error every 2 seconds. NATS needs 3 nodes (its pods require separate
  hosts), so it cannot run on one node. Wake restarts Hermes so it registers again.
- Grafana, Prometheus, Loki and alerts. `onemind-mcp` tools that read them, and the cron brief's facts about them,
  return errors.
- `tool-gateway/mcp-clickhouse` restarts in a loop because ClickHouse is off. This is harmless.
- Argo CD shows the paused apps as OutOfSync. Git changes to them are not applied until wake. Hermes' own apps
  (`platform-onemind-hermes`, `platform-hermes-base`, `platform-onemind-mcp`, `platform-core`) keep auto-sync.
- Backups (VolSync, CronJobs) of the paused apps are paused; their data does not change. Hermes, `onemind-data`
  and the cluster recovery export keep their backups.

## Steps

```sh
python3 scripts/ops/cluster_standby.py plan            # what standby would change
python3 scripts/ops/cluster_standby.py standby --yes   # apps off; all nodes stay
python3 scripts/ops/cluster_standby.py status          # wait 2 minutes, run it again: nothing should come back
```

Check that Hermes answers on Telegram and that the dashboard still asks you to log in. Then remove the nodes:

```sh
python3 scripts/ops/cluster_standby.py nodes-down --yes
```

`nodes-down` refuses to run unless standby is applied. It caps the `onemind-ecosystem-cpu` pool's autoscaler at
1 node: otherwise any pending pod would make DigitalOcean add a second $168 node. Then it drains the three
`onemind-hub-core-v2` nodes and sets that pool to 0.

**Before `nodes-down`, the offsite backup must work** (`platform/backups/OFFSITE-COPY.md` steps 1 to 5). Until it
does, every backup is stored inside DigitalOcean.

## Wake

```sh
python3 scripts/ops/cluster_standby.py nodes-up --yes  # 3 hub nodes back; ecosystem autoscaling 1-2 again
python3 scripts/ops/cluster_standby.py wake --yes      # databases, then apps, backups, CronJobs, Argo auto-sync
python3 scripts/ops/cluster_standby.py status
```

Wake order: StatefulSets (databases) first, then Deployments, then VolSync and CronJobs, then Argo auto-sync with
child apps before parent apps, then a Hermes restart once NATS is ready.

## Dry run, 2026-11-16 (apps only, all 4 nodes kept)

- Applied to 27 namespaces: 36 Applications paused, 84 workloads scaled to 0, 9 CronJobs suspended and 26 VolSync
  sources paused. Nothing was deleted.
- Bug found and fixed: `platform-security` is a parent of `ecosystem-openfga` and `ecosystem-trivy-operator`.
  Because it was paused after them, it switched their auto-sync back on and Argo restarted both pods. The script
  now finds every Application that manages Applications, pauses those first, and a second run repairs anything a
  parent reverted without overwriting the recorded original.
- After the fix, 2 minutes later: nothing had come back.
  - Hermes, `onemind-mcp`, Authentik and `onemind-data` were Running.
  - The Hermes dashboard on the tailnet redirected to login (302), and `auth.onemindos.dev` was live (200).
  - Node memory fell from 51-59% to 18-38%.
- Not done yet: `nodes-down` and a full wake. Hermes on Telegram is not checked yet: it needs a message from
  the owner.

## Notes

- `nats` and the TAK Server core are Helm releases, not Argo. Argo does not restore them, so the script records
  and restores them itself.
- The hub pool was created outside Terraform (`infra/main.tf` declares a `-core` pool). `nodes-down` and
  `nodes-up` use `doctl`, so no Terraform drift is added.
