"""Read-only deployment preflight, executed in the existing Legacy Hermes pod.

Never reads runtime config, environment credentials, private keys, or user files.
Only the SealedSecrets PUBLIC certificate is included in the report.
"""
import importlib.metadata
import importlib.util
import importlib
import json
import os
import sys
from pathlib import Path
import urllib.error
import urllib.request


def fetch(url):
    opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))
    with opener.open(url, timeout=10) as response:
        return response.read(1_048_577)


def main():
    report = {"upstreams": {}}
    if len(sys.argv) > 1:
        report["hermes_image_id"] = sys.argv[1]
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
