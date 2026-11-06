#!/usr/bin/env python3
"""Seal a read-only DigitalOcean token as SealedSecret backstage-digitalocean.

Run on the operator's Mac with kubectl access to do-nyc3-omos-cloud and kubeseal.
The token is read hidden, checked against read endpoints, sealed with the cluster's
public certificate and written only in encrypted form. Nothing secret is printed or
saved in plain text.

Usage: python3 services/onemind-portal/scripts/seal-digitalocean.py
"""
import getpass
import json
import pathlib
import subprocess
import sys
import urllib.error
import urllib.request

ROOT = pathlib.Path(__file__).resolve().parents[3]
OUT = ROOT / "platform/backstage/digitalocean.sealedsecret.yaml"
CONTEXT = "do-nyc3-omos-cloud"
# Read endpoints the portal uses. Each must answer 200 with this token.
CHECKS = {
    "account": "/v2/account",
    "billing": "/v2/customers/my/balance",
    "kubernetes": "/v2/kubernetes/clusters?per_page=1",
    "droplets": "/v2/droplets?per_page=1",
    "volumes": "/v2/volumes?per_page=1",
    "snapshots": "/v2/snapshots?per_page=1",
    "load balancers": "/v2/load_balancers?per_page=1",
    "firewalls": "/v2/firewalls?per_page=1",
    "vpcs": "/v2/vpcs?per_page=1",
    "domains": "/v2/domains?per_page=1",
    "monitoring": "/v2/monitoring/alerts?per_page=1",
    "uptime": "/v2/uptime/checks?per_page=1",
}


def status(token: str, path: str) -> int:
    request = urllib.request.Request(
        "https://api.digitalocean.com" + path,
        headers={"Authorization": "Bearer " + token},
    )
    try:
        with urllib.request.urlopen(request, timeout=15) as response:
            response.read(1 << 20)
            return response.status
    except urllib.error.HTTPError as error:
        return error.code
    except Exception:
        return 0


def main() -> None:
    print("Create the token in DigitalOcean: API > Tokens > Generate New Token > Custom Scopes,")
    print("read scopes only. Do not use your doctl or full-access token.")
    token = getpass.getpass("Read-only DigitalOcean token (hidden): ").strip()
    if not token:
        sys.exit("No token entered.")
    failed = []
    for name, path in CHECKS.items():
        code = status(token, path)
        print(f"  {name:15} {'ok' if code == 200 else 'missing (' + str(code) + ')'}")
        if code != 200:
            failed.append(name)
    if "account" in failed or "billing" in failed:
        sys.exit("The token cannot read the account or billing. Nothing was sealed.")
    if failed:
        print("Some read scopes are missing; those sections will show as unavailable.")
    cert = subprocess.run(
        ["kubeseal", "--context", CONTEXT, "--fetch-cert",
         "--controller-name=sealed-secrets-controller", "--controller-namespace=kube-system"],
        check=True, capture_output=True,
    ).stdout
    secret = {
        "apiVersion": "v1", "kind": "Secret",
        "metadata": {"name": "backstage-digitalocean", "namespace": "backstage"},
        "type": "Opaque", "stringData": {"DIGITALOCEAN_TOKEN": token},
    }
    cert_path = pathlib.Path.home() / ".config/onemind/.sealed-secrets-cert.pem"
    cert_path.parent.mkdir(mode=0o700, parents=True, exist_ok=True)
    cert_path.write_bytes(cert)
    sealed = subprocess.run(
        ["kubeseal", "--cert", str(cert_path), "--format", "yaml", "--scope", "strict"],
        input=json.dumps(secret).encode(), check=True, capture_output=True,
    ).stdout
    cert_path.unlink()
    OUT.write_bytes(sealed)
    print(f"Sealed into {OUT.relative_to(ROOT)}. No token was printed or saved in plain text.")


if __name__ == "__main__":
    main()
