"""Read-only deployment preflight, executed in the existing Legacy Hermes pod.

Never reads runtime config, environment credentials, private keys, or user files.
Only the SealedSecrets PUBLIC certificate is included in the report.
"""
import importlib.metadata
import importlib.util
import importlib
import json
import os
import re
import sys
from pathlib import Path
import urllib.error
import urllib.request

import httpx


def fetch(url):
    opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))
    with opener.open(url, timeout=10) as response:
        return response.read(1_048_577)


def main():
    report = {"upstreams": {}}
    if len(sys.argv) > 1:
        report["hermes_image_id"] = sys.argv[1]
    identity = Path("/var/run/secrets/kubernetes.io/serviceaccount")
    try:
        with httpx.Client(verify=str(identity / "ca.crt"), trust_env=False,
                          headers={"Authorization": "Bearer " + (identity / "token").read_text()}, timeout=10) as kube:
            response = kube.get("https://kubernetes.default.svc/api/v1/namespaces/hermes/pods",
                                params={"labelSelector": "app.kubernetes.io/name=onemind-mcp"})
            response.raise_for_status()
            pods = []
            for pod in response.json()["items"]:
                row = {"name": pod["metadata"]["name"], "phase": pod["status"].get("phase"), "containers": []}
                for container in pod["status"].get("containerStatuses", []):
                    state = container.get("state", {})
                    waiting = state.get("waiting", {})
                    message = waiting.get("message", "").lower()
                    row["containers"].append({
                        "name": container["name"], "ready": container["ready"],
                        "restarts": container["restartCount"], "waiting_reason": waiting.get("reason"),
                        "diagnostic_flags": [word for word in ("unauthorized", "denied", "not found", "secret", "permission", "pull") if word in message],
                        "last_exit_code": container.get("lastState", {}).get("terminated", {}).get("exitCode"),
                    })
                events = kube.get("https://kubernetes.default.svc/api/v1/namespaces/hermes/events",
                                  params={"fieldSelector": "involvedObject.name=" + row["name"]})
                row["pull_events"] = []
                if events.status_code == 200:
                    for event in events.json().get("items", []):
                        message = event.get("message", "").lower()
                        row["pull_events"].append({
                            "reason": event.get("reason"),
                            "http_status_codes": re.findall(r"\b(?:401|403|404|429|500|502|503)\b", message),
                            "diagnostic_flags": [word for word in (
                                "unauthorized", "forbidden", "denied", "not found", "manifest unknown",
                                "i/o timeout", "no such host", "x509", "connection refused", "secret",
                            ) if word in message],
                        })
                pods.append(row)
            report["mcp_pods"] = pods
    except Exception as error:
        report["pod_status_error_type"] = type(error).__name__
    for name, url in {
        "nats": "http://nats-0.nats-headless.nats.svc.cluster.local:8222/varz",
        "prometheus": "http://monitoring-kube-prometheus-prometheus.monitoring.svc.cluster.local:9090/api/v1/query?query=up&limit=1",
    }.items():
        try:
            data = json.loads(fetch(url))
            report["upstreams"][name] = {
                "reachable": True,
                "application": data.get("server_name") if name == "nats" else data.get("status"),
            }
        except Exception as error:
            # Avoid exception strings that could contain response data.
            report["upstreams"][name] = {"reachable": False, "error_type": type(error).__name__}
    for package in ("mcp", "hermes-agent"):
        try:
            report[package + "_version"] = importlib.metadata.version(package)
        except importlib.metadata.PackageNotFoundError:
            report[package + "_version"] = None
    report["tool_search_module"] = importlib.util.find_spec("tools.tool_search") is not None
    # Check source capability, never config contents.
    candidates = [Path("/opt/hermes/tools/mcp_tool.py"), Path("/opt/hermes/tools/mcp_tool_config.py")]
    report["mcp_env_expansion_source"] = any(
        path.is_file() and ("expandvars" in path.read_text() or "_expand_env" in path.read_text())
        for path in candidates)
    report["mcp_header_env_expansion_verified"] = False
    for module_name in ("tools.mcp_tool_config", "tools.mcp_tool"):
        try:
            module = importlib.import_module(module_name)
            interpolate = getattr(module, "_interpolate_env_vars", None)
            if interpolate:
                value = interpolate({"headers": {"Authorization": "Bearer ${PATH}"}})
                report["mcp_header_env_expansion_verified"] = value["headers"]["Authorization"] == "Bearer " + os.environ["PATH"]
                if report["mcp_header_env_expansion_verified"]:
                    break
        except Exception:
            continue
    try:
        certificate = fetch("http://sealed-secrets-controller.kube-system.svc.cluster.local:8080/v1/cert.pem").decode()
        if "-----BEGIN CERTIFICATE-----" not in certificate or "PRIVATE KEY" in certificate:
            raise ValueError("Unexpected public certificate response")
        report["sealed_secrets_public_certificate"] = certificate
    except Exception as error:
        report["certificate_error_type"] = type(error).__name__
    print(json.dumps(report, indent=2))


if __name__ == "__main__":
    main()
