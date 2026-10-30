---
name: tailscale-continuous-mesh-access
description: "Use to check whether the OneMind tailnet and the cluster's Tailscale proxies are healthy, and to turn that into a cheap scheduled watchdog. Not for reaching a device (see tailscale-pod-access)."
version: 2.0.0
author: Legacy
metadata:
  hermes:
    tags: [tailscale, tailnet, mesh, health, cron, egress]
    related_skills: [tailscale-pod-access, doks-cluster-ops, legacy-self-awareness]
triggers:
  - "is the tailnet healthy"
  - "tailscale is down"
  - "mesh health check"
  - "which devices are offline"
---

# Tailnet and mesh health

Rewritten from the Legacy-pod version. That one kept a **userspace `tailscaled`**
alive inside the agent pod (state under `/opt/data/home/.local`), restarted it with
`pkill`, and took five-minute snapshots of every peer. None of that exists now,
and `tailscale-pod-access` says not to rebuild it. What survives is the useful
part: a clear picture of what "healthy" means, and a cheap check you can schedule.

## How the tailnet is wired today (verified 2026-10-30)

- The cluster joins the tailnet through the **Tailscale Kubernetes operator**
  (Argo Application `edge-tailscale-operator`, namespace `tailscale`). The
  operator, not a script, owns daemon lifecycle: each exposed service has a proxy
  pod `ts-<name>-<id>-0`, and its node state lives in a Secret. A restarted pod
  re-joins by itself.
- Machines the agent manages are reached through the **ProxyGroup `egress`**
  (pod `egress-0`, `ProxyGroupReady`) and three `ExternalName` Services in
  namespace `hermes`, defined in `platform/hermes/agent/egress.yaml`:
  `omos-mbp` (22), `omos-s24` (8022), `omos-ha-edge` (8123). The Home Assistant
  box has been offline since about 2026-09-19, so a failing `omos-ha-edge` is
  expected, not an incident.
- Who may talk to what is the tailnet policy `platform/core/policy.hujson`
  (applied from `main` by `.github/workflows/tailscale.yml`). Cluster proxies
  (`tag:k8s`) may reach only the listed machine and port; nothing is public
  (Funnel is off).
- `onemind-hermes` itself has **no `tailscaled`, no `tailscale` binary, no
  `kubectl` and no service-account token**. Its NetworkPolicy
  (`platform/hermes/onemind-hermes/network.yaml`) allows only DNS, LiteLLM
  (4000), Matrix Synapse (8008) and public HTTPS, and excludes `100.64.0.0/10`.
  So from this pod the egress Services are unreachable until the owner adds a
  rule in Git. Check; do not assume.

## What "healthy" means, by layer

| Layer | Healthy | Where to look (read-only) |
|---|---|---|
| Operator | Pod `operator-*` Running in `tailscale` | `kubectl -n tailscale get pods` |
| Egress | `proxygroup/egress` `ProxyGroupReady`, `egress-0` Running | `kubectl -n hermes get proxygroup egress` (needs cluster access; this pod has none) |
| Per-device proxy | `ts-omos-mbp-*`, `ts-omos-s24-*` Running | same |
| Service target | Egress Service has a non-empty `externalName` | `kubectl -n hermes get svc` |
| Policy | Grant exists for that `tag:k8s` to machine:port | `platform/core/policy.hujson` |
| Pod network | A TCP connect to the Service succeeds | `scripts/mesh-health.py probe` |
| Peer | `Online` true, direct (`CurAddr` set) or relayed | `tailscale status --json` on a host that has the CLI |

Read the table bottom-up when something fails: a `probe` timeout with all upper
layers green means the pod's NetworkPolicy (or the tailnet grant), not the
device. A refused connection means the path works and the service on the device
is down. A DNS failure means the Service does not exist.

## The helper

`scripts/mesh-health.py` is stdlib-only Python (no `pip`, no extra binaries).

```bash
# 1. Can this pod open the egress paths? (runs anywhere in the cluster network)
python3 scripts/mesh-health.py probe                 # the three default Services
python3 scripts/mesh-health.py probe --target omos-mbp.hermes.svc.cluster.local:22

# 2. Summarize a peer list. Run `tailscale status --json` on a machine that has the
#    CLI (the owner's Mac) and pipe it in. Names and OS only; addresses are not printed.
tailscale status --json | python3 scripts/mesh-health.py peers -
python3 scripts/mesh-health.py peers status.json --offline-threshold 5
```

Exit code 0 means everything checked was reachable or online (peers named with
`--expect-offline NAME` are allowed to be offline). Exit 1 means a finding. Output
is one line per finding and empty on green, so it works as a silent watchdog.

## Scheduling it

Use a script-only cron job so no inference credit is spent while nothing is wrong
(empty output is silent):

```bash
# After the owner has opened egress to the Services and the script is under ~/.hermes/scripts/
hermes cron create "every 30m" --name tailnet-egress-probe --no-agent \
  --script mesh-health.py --deliver <owner channel> --failure-deliver local
```

Notes:

- The script goes in `~/.hermes/scripts/` (the pod's `/opt/data/scripts/`) and
  runs with Hermes' Python. Put the job definition in the PR that asks for it;
  the cron store lives on the pod volume, not in Git.
- A `--no-agent` job runs the script with no arguments, so the defaults apply. To
  check only some paths, wrap the call in a small `.sh` script that passes
  `--target`.
- Do not schedule a probe that is known to fail (policy still closed, HA box
  offline). It will page the owner with noise.
- The old five-minute cadence was for a daemon that died; the operator-managed
  proxies do not need it. Half an hour is plenty.
- Cron sessions run with `approvals.cron_mode: deny`: the job cannot run anything
  flagged as dangerous, which is right for a read-only probe.

## When it is down

1. Run the probe. Note which targets fail and how (timeout, refused, DNS).
2. All three time out: the NetworkPolicy has no rule for the egress Services, or
   the egress proxy is down. Say which you checked. Do not start a userspace
   `tailscaled`, a tunnel through another machine or an exit node.
3. One target refuses: the device is up but its service is not (Termux not
   running on the S24; Remote Login off on the Mac). That is the owner's device.
4. A pod-side fix is a Git change (egress rule, policy grant). Cluster-side
   inspection (`kubectl -n tailscale get pods`, `get proxygroup`) belongs to a
   session with cluster access, not this pod.
5. A proxy pod that is stuck: report it with the pod name and age. Deleting a
   pod is not yours to do without "execute" or "override".

## Dropped on purpose

The userspace-daemon lifecycle (`pkill`, `--tun=userspace-networking`, state-file
surgery, `tailscale up --authkey`), the REST wrapper on 127.0.0.1:8765, the
snapshot rotation, hard-coded device hostnames and tailnet addresses, and the SSH
test that embedded the Mac's address and a key path. A fresh auth key would be a
credential; it comes from the operator's SealedSecret, never from a skill.
