"""Facts for the daily `legacy-sovereign-brief` (stdout is injected into the agent's prompt).

Gathers everything cheap first so the model only has to write a short brief: cluster and alert
health, what happened in the last 24 h, backups, cost, bus health, and Codex usage for every
credential. Read-only, no cluster token: the data comes from the onemind MCP tools (see
legacy_common.py). A failing probe is listed under "Probe problems" instead of aborting.
"""

import os
import subprocess
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import legacy_common as lc  # noqa: E402

problems = []


def section(title, lines):
    print("## " + title)
    for line in lines or ["none"]:
        print("- " + line)
    print()


def nodes():
    ready = lc.prom('sum(kube_node_status_condition{condition="Ready",status="true"})')
    total = lc.prom("count(kube_node_info)")
    out = ["%d of %d nodes Ready" % (ready[0][1] if ready else 0, total[0][1] if total else 0)]
    mem = lc.prom("100 * (1 - sum(node_memory_MemAvailable_bytes) / sum(node_memory_MemTotal_bytes))")
    cpu = lc.prom('100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[30m])))')
    if mem and cpu:
        out.append("cluster memory used %.0f%%, CPU %.0f%% (30m average)" % (mem[0][1], cpu[0][1]))
    hot = lc.prom("topk(3, 100 * (1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) > 85")
    out += ["node %s memory %.0f%%" % (m.get("nodename") or m.get("instance", "?"), v) for m, v in hot]
    return out


def workloads():
    out = []
    bad = lc.prom('kube_pod_container_status_waiting_reason{reason=~"CrashLoopBackOff|ImagePullBackOff|ErrImagePull|CreateContainerConfigError"} == 1')
    out += ["%s %s" % (lc.ref(m), m.get("reason")) for m, _ in bad]
    pend = lc.prom('kube_pod_status_phase{phase=~"Pending|Failed|Unknown"} == 1')
    out += ["%s phase %s" % (lc.ref(m), m.get("phase")) for m, _ in pend]
    restarts = lc.prom("topk(5, sum by (namespace, pod) (increase(kube_pod_container_status_restarts_total[24h]))) > 0")
    out += ["%s restarted %d time(s) in 24h" % (lc.ref(m), v) for m, v in restarts]
    oom = lc.prom('sum by (namespace, pod, container) (kube_pod_container_status_last_terminated_reason{reason="OOMKilled"}'
                  ' * on (namespace, pod, container) group_left () (changes(kube_pod_container_status_restarts_total[24h]) > 0))')
    out += ["%s/%s OOMKilled in 24h" % (lc.ref(m), m.get("container")) for m, _ in oom]
    return out


def pvcs():
    rows = lc.prom("max by (namespace, persistentvolumeclaim) (kubelet_volume_stats_used_bytes / kubelet_volume_stats_capacity_bytes) > 0.8")
    return ["%s %.0f%% full" % (lc.ref(m), v * 100) for m, v in rows]


def current_alerts():
    groups = {}
    for a in lc.alerts():
        labels = a.get("labels", {})
        key = (labels.get("severity", "?"), lc.alert_name(labels))
        summary = (a.get("annotations", {}).get("summary") or "")[:140]
        count, first = groups.get(key, (0, summary))
        groups[key] = (count + 1, first)
    order = {"critical": 0, "warning": 1}
    return ["[%s] %s%s - %s" % (sev, name, " x%d" % n if n > 1 else "", first)
            for (sev, name), (n, first) in sorted(groups.items(), key=lambda kv: (order.get(kv[0][0], 2), kv[0][1]))]


def fired_24h():
    rows = lc.prom('count by (alertname, severity) (max_over_time(ALERTS{alertstate="firing", alertname!~"Watchdog|InfoInhibitor"}[24h]))')
    return sorted("%s (%s)" % (m.get("alertname"), m.get("severity", "?")) for m, _ in rows)


def argo():
    apps = lc.argo_apps()
    odd = ["%s %s" % (a.get("name"), a.get("sync")) for a in apps if a.get("sync") != "Synced"]
    return ["%d Argo applications, %d not Synced" % (len(apps), len(odd))] + odd


def backups():
    rows = lc.prom("time() - max by (obj_namespace, obj_name) (volsync_last_sync_time_seconds{role=\"source\"}) > 129600")
    return ["volume backup %s/%s last synced %.0fh ago" % (m.get("obj_namespace"), m.get("obj_name"), v / 3600) for m, v in rows]


def cost():
    data = lc.call_json("onemind-opencost_allocation", {"window": "24h"})
    total, top = 0.0, []

    def walk(node):
        nonlocal total
        if isinstance(node, dict):
            if "totalCost" in node and "name" in node:
                total += float(node["totalCost"])
                top.append((float(node["totalCost"]), str(node["name"])))
                return
            for v in node.values():
                walk(v)
        elif isinstance(node, list):
            for v in node:
                walk(v)

    walk(data)
    top.sort(reverse=True)
    return ["estimated last 24h: $%.2f; top: %s" % (total, ", ".join("%s $%.2f" % (n, c) for c, n in top[:3]))] if top else []


def bus():
    out = []
    js = lc.call_json("onemind-nats_jetstream_stats")
    if isinstance(js, dict):
        out.append("JetStream: %s streams, %s consumers, %.1f GB stored" % (js.get("streams"), js.get("consumers"), float(js.get("bytes", 0)) / 1e9))
    for node in ("nats-0", "nats-1", "nats-2"):
        try:
            st = lc.call_json("onemind-nats_server_stats", {"node": node})
            slow = st.get("slow_consumers", 0) if isinstance(st, dict) else 0
            if slow:
                out.append("%s: %s slow consumer(s)" % (node, slow))
        except lc.McpError as exc:
            out.append("%s unreachable (%s)" % (node, exc))
    return out


def codex():
    here = os.path.dirname(os.path.abspath(__file__))
    done = subprocess.run([sys.executable, os.path.join(here, "legacy_codex_usage.py")],
                          capture_output=True, text=True, timeout=90, check=False)
    return [ln for ln in done.stdout.splitlines() if ln.strip()] or ["codex usage produced no output"]


def main():
    print("# Legacy sovereign brief facts (collected by script; no model call yet)\n")
    section("Cluster", lc.safe(nodes, [], problems, "nodes"))
    section("Workloads (bad states and 24h restarts)", lc.safe(workloads, [], problems, "workloads"))
    section("Volumes over 80%", lc.safe(pvcs, [], problems, "pvcs"))
    section("Alerts firing now (excluding Watchdog)", lc.safe(current_alerts, [], problems, "alertmanager"))
    section("Alerts that fired in the last 24h", lc.safe(fired_24h, [], problems, "alert history"))
    section("GitOps", lc.safe(argo, [], problems, "argocd"))
    section("Backups older than 36h", lc.safe(backups, [], problems, "backups"))
    section("Cost", lc.safe(cost, [], problems, "opencost"))
    section("Bus (NATS)", lc.safe(bus, [], problems, "nats"))
    section("Codex usage per credential", lc.safe(codex, ["codex usage unavailable"], problems, "codex"))
    if problems:
        section("Probe problems (data missing, not necessarily a fault)", problems)
    return 0


if __name__ == "__main__":
    sys.exit(main())
