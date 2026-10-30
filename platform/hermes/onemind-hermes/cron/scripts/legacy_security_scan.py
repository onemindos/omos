"""`security-scan` monitor source: prints `OK` or one stable line per finding.

Used with `--monitor-script` once a day. The agent only runs when this output changes (see
legacy_infra_health.py for the contract: stable, sorted, no timestamps). What it looks at, all
read-only through the onemind MCP tools:

- alerts that look security related (Falco, Kyverno, Authentik, certificates, auth failures);
- Kyverno policy failures in the last 24 h;
- running containers on a floating `:latest` tag (CLAUDE.md says there should be none);
- Argo CD applications that drifted from Git (a manual change shows up as OutOfSync).

Certificate expiry is covered only through alerts (for example KubeClientCertificateExpiration):
cert-manager metrics are not scraped today, and this pod has no Kubernetes token.

Known accepted findings go in ACCEPTED so they do not wake the agent every time.
"""

import os
import re
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import legacy_common as lc  # noqa: E402

SECURITY_ALERT = re.compile(r"falco|kyverno|authentik|cert|tls|auth|unauthori|forbidden|intrusion|vulnerab|secret", re.I)

# Findings the owner has accepted. Match on the exact finding text prefix.
ACCEPTED = (
    "latest-tag onetak/webodm",  # upstream WebODM images are only published as :latest
)


def main():
    lines = []

    for a in lc.safe(lc.alerts, [], lines, "alertmanager"):
        labels = a.get("labels", {})
        name = lc.alert_name(labels)
        if (labels.get("source") == "falco" or SECURITY_ALERT.search(name)) and labels.get("severity") in ("warning", "critical", "Critical"):
            lines.append("alert %s %s" % (name, lc.alert_ns(labels)))

    def kyverno():
        rows = lc.prom('sum by (policy_name, resource_namespace) (increase(kyverno_policy_results_total{rule_result="fail"}[24h])) > 0.5')
        return ["policy-fail %s %s" % (m.get("policy_name"), m.get("resource_namespace", "")) for m, _ in rows]

    lines += lc.safe(kyverno, [], lines, "kyverno")

    def latest_tag():
        rows = lc.prom('count by (namespace, pod) (kube_pod_container_info{image=~".*:latest"})')
        names = {}
        for m, _ in rows:
            # collapse replica-set suffixes so a pod restart does not change the line
            base = re.sub(r"-[a-f0-9]{8,10}-[a-z0-9]{5}$|-[a-z0-9]{5}$", "", m.get("pod", ""))
            names["latest-tag %s/%s" % (m.get("namespace"), base)] = 1
        return list(names)

    lines += lc.safe(latest_tag, [], lines, "image tags")

    def drift():
        return ["argocd-drift %s %s" % (a.get("name"), a.get("sync")) for a in lc.argo_apps() if a.get("sync") != "Synced"]

    lines += lc.safe(drift, [], lines, "argocd")

    findings = sorted(l for l in set(lines) if not l.startswith(ACCEPTED))
    print("\n".join(findings) if findings else "OK")
    return 0


if __name__ == "__main__":
    sys.exit(main())
