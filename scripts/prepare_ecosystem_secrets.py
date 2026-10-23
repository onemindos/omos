"""Prepare initial namespace-bound SealedSecrets without writing plaintext secrets.

Run once with the cluster's sealing certificate. Existing outputs are never
overwritten: rotation must be a separate reviewed operation.
"""
import argparse
import base64
import json
from pathlib import Path
import secrets
import subprocess

ROOT = Path(__file__).resolve().parents[1]
OUTPUT = ROOT / "platform/ecosystem/resources/secrets"
NAMESPACES = ["onemind-data", "onemind-core", "onemind-ai",
              "onemind-search", "onemind-registry", "velero"]


def existing(namespace, name):
    obj = json.loads(subprocess.check_output([
        "kubectl", "get", "secret", name, "-n", namespace, "-o", "json",
    ]))
    return {k: base64.b64decode(v).decode() for k, v in obj["data"].items()}


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--cert", type=Path, required=True)
    args = parser.parse_args()
    context = subprocess.check_output(["kubectl", "config", "current-context"], text=True).strip()
    if context != "do-nyc3-omos-cloud":
        raise RuntimeError("Select the verified cloud context")
    if OUTPUT.exists() and list(OUTPUT.glob("*.json")):
        raise RuntimeError("Encrypted outputs already exist; do not rotate them by rerunning initial setup")
    OUTPUT.mkdir(parents=True, exist_ok=True)
    pending = []

    def add(ns, name, data, kind="Opaque", labels=None):
        secret = {"apiVersion": "v1", "kind": "Secret", "type": kind,
                  "metadata": {"name": name, "namespace": ns, "labels": labels or {}},
                  "data": {k: base64.b64encode(v.encode()).decode() for k, v in data.items()}}
        result = subprocess.run(["kubeseal", "--cert", str(args.cert), "--format", "json"],
                                input=json.dumps(secret), text=True, stdout=subprocess.PIPE,
                                stderr=subprocess.PIPE)
        if result.returncode:
            raise RuntimeError("Could not seal " + ns + "/" + name)
        sealed = json.loads(result.stdout)
        sealed["metadata"]["annotations"] = {"argocd.argoproj.io/sync-wave": "-1"}
        pending.append((OUTPUT / (ns + "--" + name + ".json"), sealed))

    roles = {n: secrets.token_urlsafe(40) for n in
             ["onemind", "dapr_workflows", "litellm", "openfga", "langfuse"]}
    app_namespaces = {"onemind": "onemind-core", "dapr_workflows": "onemind-core",
                      "litellm": "onemind-ai", "openfga": "onemind-core",
                      "langfuse": "onemind-ai"}
    uris = {}
    for role, password in roles.items():
        pair = {"username": role, "password": password}
        add("onemind-data", "db-" + role.replace("_", "-"), pair,
            "kubernetes.io/basic-auth", {"cnpg.io/reload": "true"})
        uri = f"postgresql://{role}:{password}@onemind-data-rw.onemind-data.svc.cluster.local:5432/{role}?sslmode=require"
        uris[role] = uri
        add(app_namespaces[role], role.replace("_", "-") + "-db", {**pair, "uri": uri})

    master = "sk-" + secrets.token_urlsafe(40)
    clients = {n: "sk-" + secrets.token_urlsafe(40) for n in ["core", "openchoreo"]}
    add("onemind-ai", "litellm-auth", {"master-key": master})
    add("onemind-ai", "litellm-client-keys", clients)
    # Provider credentials come from the existing LiteLLM provider secret (OpenClaw, their former source, was retired).
    provider = existing("onemind-ai", "litellm-providers")
    keys = {k: provider[k] for k in ["OPENROUTER_API_KEY", "OPENAI_API_KEY", "ANTHROPIC_API_KEY"] if provider.get(k)}
    if not keys.get("OPENROUTER_API_KEY"):
        raise RuntimeError("The selected model provider has no configured credential")
    add("onemind-ai", "litellm-providers", keys)
    gateway = secrets.token_urlsafe(48)
    nats = existing("hermes", "omos-nats-creds")
    add("onemind-core", "core-runtime", {
        "DATABASE_URL": uris["onemind"], "UI_GATEWAY_KEY": gateway,
        "READ_API_KEY": secrets.token_urlsafe(48),
        "MODEL_GATEWAY_URL": "http://litellm.onemind-ai.svc.cluster.local:4000/v1",
        "MODEL_GATEWAY_KEY": clients["core"],
        "NATS_URL": "nats://nats.nats.svc.cluster.local:4222",
        "NATS_USER": nats["user"], "NATS_PASSWORD": nats["password"],
    })
    add("backstage", "onemind-core-gateway", {"ONEMIND_CORE_GATEWAY_KEY": gateway})
    add("onemind-core", "openfga-auth", {"keys": secrets.token_urlsafe(48)})
    add("onemind-search", "opensearch-auth", {"password": "Om1!" + secrets.token_urlsafe(40)})
    add("onemind-registry", "harbor-admin", {"password": secrets.token_urlsafe(40)})
    backup = existing("ops", "platform-backup-credentials")
    add("velero", "velero-credentials", {"cloud": "[default]\naws_access_key_id=" + backup["S3_ACCESS_KEY"] +
                                         "\naws_secret_access_key=" + backup["S3_SECRET_KEY"] + "\n"})
    add("onemind-ai", "langfuse-s3", {"access-key": backup["S3_ACCESS_KEY"], "secret-key": backup["S3_SECRET_KEY"]})
    add("onemind-ai", "langfuse-auth", {"salt": secrets.token_urlsafe(48), "encryption-key": secrets.token_hex(32),
                                        "nextauth-secret": secrets.token_urlsafe(48),
                                        "redis-password": secrets.token_urlsafe(40),
                                        "clickhouse-password": secrets.token_urlsafe(40)})
    pull = existing("backstage", "ghcr")
    for ns in NAMESPACES:
        add(ns, "ghcr", pull, "kubernetes.io/dockerconfigjson")
    # Write only after every item has been successfully encrypted.
    for path, obj in pending:
        path.write_text(json.dumps(obj, indent=2) + "\n")
    print(f"Prepared {len(pending)} namespace-bound encrypted resources; plaintext was not written")


if __name__ == "__main__":
    main()
