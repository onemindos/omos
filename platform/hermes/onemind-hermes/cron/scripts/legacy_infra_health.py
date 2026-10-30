"""`infra-health` monitor source: prints `OK` when healthy, else one stable line per problem.

Used with `--monitor-script` (every 15 minutes). Hermes hashes this script's exact output: an
unchanged hash suppresses the agent run, so a healthy or unchanged-broken cluster costs no
tokens. The output must therefore be stable: no timestamps, no changing counters, sorted lines.

Only conditions worth waking someone for are listed (critical alerts, NotReady nodes, stuck
or crash-looping pods, full volumes, deployments short of replicas for 15 minutes, a NATS node
down, the MCP gateway itself unreachable). Warnings and drift belong to the daily brief and
the security scan. If the gateway is unreachable the only line is the gateway line, because
nothing else can be known.
"""

import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import legacy_common as lc  # noqa: E402


def main():
    lines = []
    try:
        ready = lc.prom('kube_node_status_condition{condition="Ready",status!="true"} == 1')
    except lc.McpError as exc:
        print("monitoring unreachable: %s" % exc)
        return 0
    lines += ["node %s NotReady" % m.get("node", "?") for m, _ in ready]

    def probe(query, fmt):
        for m, v in lc.safe(lambda: lc.prom(query), [], lines, "query"):
            lines.append(fmt(m, v))

    probe('kube_pod_container_status_waiting_reason{reason=~"CrashLoopBackOff|ImagePullBackOff|ErrImagePull|CreateContainerConfigError"} == 1',
          lambda m, v: "pod %s %s" % (lc.ref(m), m.get("reason")))
    probe('min_over_time((kube_pod_status_phase{phase=~"Pending|Failed"} == 1)[15m:1m]) == 1',
          lambda m, v: "pod %s stuck %s" % (lc.ref(m), m.get("phase")))
    probe("min_over_time((kube_deployment_spec_replicas - kube_deployment_status_replicas_available)[15m:1m]) > 0",
          lambda m, v: "deployment %s short of replicas" % lc.ref(m))
    probe("max by (namespace, persistentvolumeclaim) (kubelet_volume_stats_used_bytes / kubelet_volume_stats_capacity_bytes) > 0.92",
          lambda m, v: "volume %s over 92%% full" % lc.ref(m))

    for a in lc.safe(lc.alerts, [], lines, "alertmanager"):
        labels = a.get("labels", {})
        if labels.get("severity") == "critical" and labels.get("source") != "falco":  # Falco belongs to the security scan
            who = labels.get("pod") or labels.get("deployment") or labels.get("obj_name") or labels.get("job") or ""
            lines.append("alert %s %s" % (lc.alert_name(labels), "/".join(x for x in (lc.alert_ns(labels), who) if x)))

    for node in ("nats-0", "nats-1", "nats-2"):
        try:
            lc.call_json("onemind-nats_server_stats", {"node": node})
        except lc.McpError:
            lines.append("nats node %s unreachable" % node)

    print("\n".join(sorted(set(lines))) if lines else "OK")
    return 0


if __name__ == "__main__":
    sys.exit(main())
