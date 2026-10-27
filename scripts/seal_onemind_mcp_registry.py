"""Verify the dedicated read-only GHCR credential and seal it without plaintext files."""
import base64
import json
import os
from pathlib import Path
import subprocess
import urllib.error
import urllib.request

import yaml


def request(url, authorization, method="GET", accept=None):
    headers = {"Authorization": authorization, "User-Agent": "onemind-mcp-registry-setup"}
    if accept:
        headers["Accept"] = accept
    try:
        with urllib.request.urlopen(urllib.request.Request(url, headers=headers, method=method), timeout=30) as response:
            return response.read()
    except urllib.error.HTTPError as error:
        raise RuntimeError(f"Registry credential verification failed: HTTP {error.code}") from None
    except urllib.error.URLError:
        raise RuntimeError("Registry credential verification network failure") from None


def main():
    token = os.environ.get("ONEMIND_MCP_GHCR_PULL_TOKEN", "")
    if not token:
        raise RuntimeError("Configure repository secret ONEMIND_MCP_GHCR_PULL_TOKEN first")
    username = json.loads(request("https://api.github.com/user", "Bearer " + token))["login"]
    auth = base64.b64encode(f"{username}:{token}".encode()).decode()
    registry_token = json.loads(request(
        "https://ghcr.io/token?service=ghcr.io&scope=repository:onemindos/onemind-mcp:pull", "Basic " + auth))["token"]
    image = yaml.safe_load(Path("platform/hermes/mcp/deployment.yaml").read_text())["spec"]["template"]["spec"]["containers"][0]["image"]
    digest = image.split("@", 1)[1]
    request("https://ghcr.io/v2/onemindos/onemind-mcp/manifests/" + digest,
            "Bearer " + registry_token, method="HEAD",
            accept="application/vnd.oci.image.index.v1+json, application/vnd.docker.distribution.manifest.v2+json")
    docker_config = json.dumps({"auths": {"ghcr.io": {"auth": auth}}})
    secret = {"apiVersion": "v1", "kind": "Secret",
              "metadata": {"name": "onemind-mcp-registry", "namespace": "hermes"},
              "type": "kubernetes.io/dockerconfigjson",
              "data": {".dockerconfigjson": base64.b64encode(docker_config.encode()).decode()}}
    result = subprocess.run([os.environ["KUBESEAL"], "--cert", "scripts/data/onemind-mcp-sealing.crt",
                             "--format", "json", "--scope", "namespace-wide"],
                            input=json.dumps(secret), capture_output=True, text=True)
    if result.returncode:
        raise RuntimeError("Credential sealing failed; no plaintext was written")
    sealed = json.loads(result.stdout)
    sealed["metadata"].setdefault("annotations", {})["argocd.argoproj.io/sync-wave"] = "-2"
    output = Path(os.environ["SEALED_OUTPUT"])
    output.write_text(yaml.safe_dump(sealed, sort_keys=False))
    print("Private image pull authorized; dedicated registry credential sealed")


if __name__ == "__main__":
    main()
