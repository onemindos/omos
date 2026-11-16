#!/usr/bin/env python3
"""Put the DOKS cluster on standby (Hermes only) and wake it again.

Standby keeps every volume, Application and manifest. It only:
  1. pauses Argo CD auto-sync on every Application that touches a standby
     namespace (and on the parent Applications, so they don't re-enable it),
  2. suspends CronJobs and pauses VolSync ReplicationSources in those namespaces,
  3. scales Deployments, then StatefulSets, in those namespaces to 0.

Every change records the original value in an annotation on the object
itself (onemind.io/standby-*), so `wake` restores exactly what was there and
needs no other state. Nodes are a separate, explicit step (nodes-down/nodes-up).

Runbook: docs/operations/cluster-standby.md

Usage:
  cluster_standby.py plan                      # show what standby would change
  cluster_standby.py standby --yes             # apps off, nodes untouched
  cluster_standby.py status
  cluster_standby.py nodes-down --yes          # drain + remove the hub pool
  cluster_standby.py nodes-up --yes            # bring the hub pool back
  cluster_standby.py wake --yes                # apps back on
  --keep NS  (repeatable) keeps an extra namespace running during standby
"""
import argparse
import json
import subprocess
import sys
import time

EXPECTED_CONTEXT = "do-nyc3-omos-cloud"
CLUSTER = "omos-cloud"
HUB_POOL = "onemind-hub-core-v2"
HUB_POOL_COUNT = 3
ECO_POOL = "onemind-ecosystem-cpu"
ECO_POOL_MAX = 2  # live setting on 2026-11-16 (autoscale 1-2)

# Namespaces that stay running: Hermes, what it and its login need, how it is
# reached, GitOps, admission webhooks (failurePolicy Fail) and backups.
KEEP_NAMESPACES = {
    "kube-system", "kube-public", "kube-node-lease", "default",
    "hermes",          # Hermes + onemind-mcp
    "tool-gateway",    # Hermes egress: tool calls
    "tailscale",       # Hermes ingress (*.ts.net) and SSH to devices
    "auth",            # Authentik: Hermes dashboard OIDC login
    "onemind-data",    # CNPG Postgres that holds the Authentik database
    "cnpg-system",     # CNPG operator + webhook (failurePolicy Fail)
    "kyverno",         # admission webhooks (failurePolicy Fail)
    "cert-manager",    # webhook (failurePolicy Fail) + TLS renewals
    "ingress-nginx",   # auth.onemindos.dev; webhook (failurePolicy Fail)
    "argocd",          # needed to wake everything back up
    "github-runner",   # github-runner-hermes; platform-core keeps it synced
    "ops", "velero",   # backups (VolSync controller, Velero)
    "falco",           # DaemonSet only: runtime security
}
# Workloads in kept namespaces that still go to 0 (operators for apps that are off).
EXTRA_WORKLOADS = [("kube-system", "deployment", "platform-clickhouse-operator")]
# Parents that (re)write child Applications: pause them first, restore them last.
PARENT_APPS = ["app-of-apps", "platform-ecosystem-applications",
               "platform-ecosystem-foundation", "platform-ecosystem-resources"]
EXTRA_APPS = ["platform-clickhouse-operator"]

A_SYNC = "onemind.io/standby-syncpolicy"
A_REPLICAS = "onemind.io/standby-replicas"
A_SUSPEND = "onemind.io/standby-suspend"
A_PAUSED = "onemind.io/standby-paused"


def kubectl(*args, capture=True):
    r = subprocess.run(["kubectl", *args], text=True, capture_output=capture)
    if r.returncode != 0:
        raise RuntimeError(f"kubectl {' '.join(args)}: {r.stderr.strip()}")
    return r.stdout


def kget(kind, ns=None):
    args = ["get", kind, "-o", "json"] + (["-n", ns] if ns else ["-A"])
    return json.loads(kubectl(*args))["items"]


def ann(obj):
    return obj["metadata"].get("annotations") or {}


def check_context():
    ctx = kubectl("config", "current-context").strip()
    if ctx != EXPECTED_CONTEXT:
        sys.exit(f"refusing: kubectl context is {ctx!r}, expected {EXPECTED_CONTEXT!r}")


def standby_namespaces(keep):
    names = {n["metadata"]["name"] for n in kget("namespaces")}
    return sorted(names - KEEP_NAMESPACES - set(keep))


def is_parent(app):
    """An Application that manages other Applications (app-of-apps pattern)."""
    return app["metadata"]["name"] in PARENT_APPS or any(
        r.get("kind") == "Application" for r in app.get("status", {}).get("resources", []))


def apps_to_pause(drop):
    drop = set(drop)
    out = []
    for a in kget("applications.argoproj.io", "argocd"):
        name = a["metadata"]["name"]
        touched = {r.get("namespace") for r in a.get("status", {}).get("resources", [])}
        if a["spec"]["destination"].get("namespace") in drop:
            touched.add(a["spec"]["destination"]["namespace"])
        if is_parent(a) or name in EXTRA_APPS or touched & drop:
            out.append(a)
    # Every parent first, so none can re-enable a child we just paused
    # (2026-11-16 dry run: platform-security re-enabled ecosystem-openfga).
    out.sort(key=lambda a: (not is_parent(a), a["metadata"]["name"]))
    return out


def workloads(drop):
    """Deployments first (operators stop before their StatefulSets), then StatefulSets."""
    out = []
    for kind in ("deployment", "statefulset"):
        for ns in drop:
            for w in kget(kind + "s", ns):
                out.append((ns, kind, w))
    for ns, kind, name in EXTRA_WORKLOADS:
        w = json.loads(kubectl("get", kind, name, "-n", ns, "-o", "json"))
        out.insert(0, (ns, kind, w))
    return out


def cronjobs(drop):
    return [(ns, c) for ns in drop for c in kget("cronjobs", ns)]


def replication_sources(drop):
    try:
        items = kget("replicationsources.volsync.backube")
    except RuntimeError:
        return []
    return [r for r in items if r["metadata"]["namespace"] in set(drop)]


# ---------------------------------------------------------------- standby

def do_standby(drop, apply):
    verb = "APPLY" if apply else "PLAN"
    print(f"[{verb}] standby namespaces ({len(drop)}): {' '.join(drop)}\n")

    print("1. Pause Argo CD auto-sync")
    for a in apps_to_pause(drop):
        name, sp = a["metadata"]["name"], a["spec"].get("syncPolicy") or {}
        recorded = A_SYNC in ann(a)
        if "automated" not in sp:
            print(f"   {'=' if recorded else '-'} {name} "
                  f"({'already paused' if recorded else 'manual sync already'})"); continue
        # Recorded but automated again: a parent re-applied it before it was
        # paused. Keep the original record; just remove automated again.
        print(f"   ~ {name}{' (re-enabled by its parent; pausing again)' if recorded else ''}")
        if apply:
            if not recorded:
                kubectl("annotate", "application", name, "-n", "argocd",
                        f"{A_SYNC}={json.dumps(sp, separators=(',', ':'))}")
            kubectl("patch", "application", name, "-n", "argocd", "--type", "json",
                    "-p", '[{"op":"remove","path":"/spec/syncPolicy/automated"}]')

    print("2. Suspend CronJobs")
    for ns, c in cronjobs(drop):
        name, susp = c["metadata"]["name"], bool(c["spec"].get("suspend"))
        if A_SUSPEND in ann(c) or susp:
            print(f"   = {ns}/{name} (already suspended)"); continue
        print(f"   ~ {ns}/{name}")
        if apply:
            kubectl("annotate", "cronjob", name, "-n", ns, "--overwrite", f"{A_SUSPEND}=false")
            kubectl("patch", "cronjob", name, "-n", ns, "--type", "merge",
                    "-p", '{"spec":{"suspend":true}}')

    print("3. Pause VolSync backups of standby volumes (data is not changing)")
    for r in replication_sources(drop):
        ns, name = r["metadata"]["namespace"], r["metadata"]["name"]
        if A_PAUSED in ann(r) or r["spec"].get("paused"):
            print(f"   = {ns}/{name} (already paused)"); continue
        print(f"   ~ {ns}/{name}")
        if apply:
            kubectl("annotate", "replicationsource", name, "-n", ns, "--overwrite", f"{A_PAUSED}=false")
            kubectl("patch", "replicationsource", name, "-n", ns, "--type", "merge",
                    "-p", '{"spec":{"paused":true}}')

    print("4. Scale workloads to 0 (Deployments, then StatefulSets)")
    for ns, kind, w in workloads(drop):
        name, rep = w["metadata"]["name"], w["spec"].get("replicas", 1)
        recorded = A_REPLICAS in ann(w)
        if rep == 0:
            print(f"   {'=' if recorded else '-'} {ns}/{kind}/{name} "
                  f"({'already at standby' if recorded else 'already 0'})"); continue
        print(f"   ~ {ns}/{kind}/{name} {rep} -> 0{' (was scaled back up; again)' if recorded else ''}")
        if apply:
            if not recorded:  # never overwrite the original count with a reverted one
                kubectl("annotate", kind, name, "-n", ns, f"{A_REPLICAS}={rep}")
            kubectl("scale", kind, name, "-n", ns, "--replicas=0")

    if not apply:
        print("\nNothing changed. Run `standby --yes` to apply.")


# ---------------------------------------------------------------- wake

def annotated(kind, key):
    out = []
    for o in kget(kind):
        if key in ann(o):
            out.append(o)
    return out


def do_wake(apply):
    verb = "APPLY" if apply else "PLAN"
    print(f"[{verb}] wake\n")

    print("1. Scale workloads back (StatefulSets first: databases before apps)")
    for kind in ("statefulset", "deployment"):
        for w in annotated(kind + "s", A_REPLICAS):
            ns, name, rep = w["metadata"]["namespace"], w["metadata"]["name"], ann(w)[A_REPLICAS]
            print(f"   ~ {ns}/{kind}/{name} -> {rep}")
            if apply:
                kubectl("scale", kind, name, "-n", ns, f"--replicas={rep}")
                kubectl("annotate", kind, name, "-n", ns, f"{A_REPLICAS}-")

    print("2. Resume VolSync backups")
    for r in annotated("replicationsources.volsync.backube", A_PAUSED):
        ns, name = r["metadata"]["namespace"], r["metadata"]["name"]
        print(f"   ~ {ns}/{name}")
        if apply:
            kubectl("patch", "replicationsource", name, "-n", ns, "--type", "merge",
                    "-p", '{"spec":{"paused":false}}')
            kubectl("annotate", "replicationsource", name, "-n", ns, f"{A_PAUSED}-")

    print("3. Resume CronJobs")
    for c in annotated("cronjobs", A_SUSPEND):
        ns, name = c["metadata"]["namespace"], c["metadata"]["name"]
        print(f"   ~ {ns}/{name}")
        if apply:
            kubectl("patch", "cronjob", name, "-n", ns, "--type", "merge",
                    "-p", '{"spec":{"suspend":false}}')
            kubectl("annotate", "cronjob", name, "-n", ns, f"{A_SUSPEND}-")

    print("4. Restore Argo CD auto-sync (children first, parents last)")
    apps = annotated("applications.argoproj.io", A_SYNC)
    apps.sort(key=lambda a: (is_parent(a), a["metadata"]["name"]))
    for a in apps:
        name, sp = a["metadata"]["name"], ann(a)[A_SYNC]
        print(f"   ~ {name}")
        if apply:
            patch = json.dumps([{"op": "add", "path": "/spec/syncPolicy", "value": json.loads(sp)}])
            kubectl("patch", "application", name, "-n", "argocd", "--type", "json", "-p", patch)
            kubectl("annotate", "application", name, "-n", "argocd", f"{A_SYNC}-")

    # Hermes' NATS heartbeat publisher exits when the bus is gone and does not
    # come back by itself (seen 2026-11-16), so restart Hermes once NATS is up.
    print("5. Restart Hermes after NATS is ready (re-registers agents.*.hermes.zeus.legacy)")
    if apply:
        kubectl("rollout", "status", "statefulset/nats", "-n", "nats", "--timeout=600s", capture=False)
        kubectl("rollout", "restart", "deployment/onemind-hermes", "-n", "hermes")

    if not apply:
        print("\nNothing changed. Run `wake --yes` to apply.")


# ---------------------------------------------------------------- nodes

def doctl(*args):
    r = subprocess.run(["doctl", *args], text=True, capture_output=True)
    if r.returncode != 0:
        raise RuntimeError(f"doctl {' '.join(args)}: {r.stderr.strip()}")
    return r.stdout


def hub_nodes():
    return [n["metadata"]["name"] for n in kget("nodes")
            if n["metadata"]["name"].startswith(HUB_POOL + "-")]


def do_nodes_down(apply):
    if not annotated("applications.argoproj.io", A_SYNC):
        sys.exit("refusing: run `standby --yes` first (no paused Applications found)")
    nodes = hub_nodes()
    print(f"[{'APPLY' if apply else 'PLAN'}] cap {ECO_POOL} autoscaling at 1 node, "
          f"drain {nodes}, then set {HUB_POOL} to 0 nodes")
    if not apply:
        return
    doctl("kubernetes", "cluster", "node-pool", "update", CLUSTER, ECO_POOL,
          "--auto-scale", "--min-nodes", "1", "--max-nodes", "1")
    for n in nodes:
        kubectl("cordon", n)
    for n in nodes:
        print(f"   draining {n} ...", flush=True)
        kubectl("drain", n, "--ignore-daemonsets", "--delete-emptydir-data",
                "--timeout=600s", capture=False)
    doctl("kubernetes", "cluster", "node-pool", "update", CLUSTER, HUB_POOL, "--count", "0")
    print("   hub pool set to 0. `status` shows when the nodes are gone.")


def do_nodes_up(apply):
    print(f"[{'APPLY' if apply else 'PLAN'}] set {HUB_POOL} to {HUB_POOL_COUNT} nodes, "
          f"restore {ECO_POOL} autoscaling 1-{ECO_POOL_MAX}")
    if not apply:
        return
    doctl("kubernetes", "cluster", "node-pool", "update", CLUSTER, HUB_POOL,
          "--count", str(HUB_POOL_COUNT))
    doctl("kubernetes", "cluster", "node-pool", "update", CLUSTER, ECO_POOL,
          "--auto-scale", "--min-nodes", "1", "--max-nodes", str(ECO_POOL_MAX))
    want = HUB_POOL_COUNT
    for _ in range(60):
        ready = [n for n in kget("nodes") if n["metadata"]["name"].startswith(HUB_POOL + "-")
                 and any(c["type"] == "Ready" and c["status"] == "True" for c in n["status"]["conditions"])]
        print(f"   hub nodes ready: {len(ready)}/{want}", flush=True)
        if len(ready) >= want:
            return
        time.sleep(20)
    print("   still waiting; check `kubectl get nodes` before running wake")


# ---------------------------------------------------------------- status

def do_status():
    apps = annotated("applications.argoproj.io", A_SYNC)
    ws = annotated("deployments", A_REPLICAS) + annotated("statefulsets", A_REPLICAS)
    print(f"mode: {'STANDBY' if apps or ws else 'normal'}")
    print(f"paused Applications: {len(apps)}   workloads at standby: {len(ws)}")
    nodes = kget("nodes")
    print(f"nodes: {len(nodes)} ({', '.join(n['metadata']['name'] for n in nodes)})")
    bad = [p for p in kget("pods") if p["status"].get("phase") not in ("Running", "Succeeded")]
    print(f"pods not running: {len(bad)}")
    for p in bad[:15]:
        print(f"   {p['metadata']['namespace']}/{p['metadata']['name']} {p['status'].get('phase')}")


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("command", choices=["plan", "standby", "wake", "status", "nodes-down", "nodes-up"])
    ap.add_argument("--yes", action="store_true", help="apply (otherwise print the plan only)")
    ap.add_argument("--keep", action="append", default=[], help="extra namespace to keep running")
    a = ap.parse_args()
    check_context()
    if a.command == "status":
        return do_status()
    if a.command in ("plan", "standby"):
        return do_standby(standby_namespaces(a.keep), a.command == "standby" and a.yes)
    if a.command == "wake":
        return do_wake(a.yes)
    if a.command == "nodes-down":
        return do_nodes_down(a.yes)
    if a.command == "nodes-up":
        return do_nodes_up(a.yes)


if __name__ == "__main__":
    main()
