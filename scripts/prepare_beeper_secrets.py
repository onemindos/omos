#!/usr/bin/env python3
"""Prepare initial Beeper SealedSecrets; never apply resources or print credentials."""
import argparse
import base64
import json
from pathlib import Path
import re
import secrets
import subprocess
from urllib.parse import urlparse

ROOT = Path(__file__).resolve().parents[1]
OUTPUT = ROOT / "platform/beeper-bridges/sealed"
BRIDGES = ("whatsapp", "telegram", "gmessages", "instagram", "twitter")
NAMESPACE = "beeper-bridges"


def read_private(path):
    value = path.read_text().strip()
    if not value:
        raise ValueError("A required private input file is empty")
    return value


def prepare_payloads(config, api_id, api_hash, backup):
    token = config.get("environments", {}).get("prod", {}).get("access_token", "")
    if not isinstance(token, str) or not token.startswith(("syt_", "bat_")):
        raise ValueError("Use a completed production bbctl login, not a Desktop MCP token")
    if not re.fullmatch(r"[1-9][0-9]*", api_id):
        raise ValueError("Telegram API ID must be a positive integer")
    if not re.fullmatch(r"[a-fA-F0-9]{32}", api_hash):
        raise ValueError("Telegram API hash must have 32 hexadecimal characters")
    required = ("S3_ACCESS_KEY", "S3_SECRET_KEY", "S3_BUCKET", "S3_ENDPOINT")
    if any(not backup.get(key) for key in required):
        raise ValueError("The established backup credential lacks required S3 fields")
    endpoint = backup["S3_ENDPOINT"].rstrip("/")
    if "://" not in endpoint:
        endpoint = "https://" + endpoint
    if urlparse(endpoint).scheme not in ("https", "http"):
        raise ValueError("Unsupported backup endpoint scheme")
    payloads = {
        "beeper-bridge-auth": {"MATRIX_ACCESS_TOKEN": token},
        "beeper-telegram-auth": {"TELEGRAM_API_ID": api_id, "TELEGRAM_API_HASH": api_hash},
    }
    for bridge in BRIDGES:
        payloads[f"beeper-{bridge}-restic"] = {
            "RESTIC_REPOSITORY": f"s3:{endpoint}/{backup['S3_BUCKET']}/backups/restic/beeper-bridges/{bridge}",
            "RESTIC_PASSWORD": secrets.token_urlsafe(48),
            "AWS_ACCESS_KEY_ID": backup["S3_ACCESS_KEY"],
            "AWS_SECRET_ACCESS_KEY": backup["S3_SECRET_KEY"],
            "AWS_DEFAULT_REGION": "nyc3",
        }
    return payloads


def seal(name, values, cert):
    plain = {
        "apiVersion": "v1", "kind": "Secret", "type": "Opaque",
        "metadata": {"name": name, "namespace": NAMESPACE},
        "data": {key: base64.b64encode(value.encode()).decode() for key, value in values.items()},
    }
    result = subprocess.run(
        ["kubeseal", "--cert", str(cert), "--scope", "strict", "--format", "json"],
        input=json.dumps(plain), text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
        check=False,
    )
    if result.returncode:
        raise RuntimeError("Sealing failed; no plaintext output was written")
    encrypted = json.loads(result.stdout)
    data = encrypted.get("spec", {}).get("encryptedData", {})
    if encrypted.get("kind") != "SealedSecret" or set(data) != set(values):
        raise RuntimeError("Sealing output did not contain the expected encrypted resource")
    encrypted["metadata"]["annotations"] = {"argocd.argoproj.io/sync-wave": "-1"}
    return encrypted


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--context", default="do-nyc3-omos-cloud")
    parser.add_argument("--cert", type=Path, required=True)
    parser.add_argument("--bbctl-config", type=Path, required=True)
    parser.add_argument("--telegram-api-id-file", type=Path, required=True)
    parser.add_argument("--telegram-api-hash-file", type=Path, required=True)
    args = parser.parse_args()
    if args.context != "do-nyc3-omos-cloud":
        raise ValueError("This initial rollout targets the verified omos-cloud context")
    if list(OUTPUT.glob("*.json")) or list(OUTPUT.glob("*.sealedsecret.yaml")):
        raise RuntimeError("Encrypted resources already exist; initial setup does not rotate them")
    config = json.loads(args.bbctl_config.read_text())
    api_id = read_private(args.telegram_api_id_file)
    api_hash = read_private(args.telegram_api_hash_file)
    result = subprocess.run(
        ["kubectl", "--context", args.context, "-n", "ops", "get", "secret",
         "platform-backup-credentials", "-o", "json"],
        text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=False,
    )
    if result.returncode:
        raise RuntimeError("Could not read the established backup credential")
    existing = json.loads(result.stdout).get("data", {})
    backup = {key: base64.b64decode(value).decode() for key, value in existing.items()}
    payloads = prepare_payloads(config, api_id, api_hash, backup)
    # Complete every encryption before creating an output. Only ciphertext is saved.
    pending = {name: seal(name, values, args.cert) for name, values in payloads.items()}
    OUTPUT.mkdir(parents=True, exist_ok=True)
    for name, obj in pending.items():
        path = OUTPUT / f"{name}.json"
        with path.open("x") as handle:
            handle.write(json.dumps(obj, indent=2) + "\n")
    kustomization = [
        "apiVersion: kustomize.config.k8s.io/v1beta1",
        "kind: Kustomization", "resources:",
        *[f"- {name}.json" for name in sorted(pending)], "",
    ]
    (OUTPUT / "kustomization.yaml").write_text("\n".join(kustomization))
    print(f"Prepared {len(pending)} encrypted resources; nothing was applied")


if __name__ == "__main__":
    try:
        main()
    except (ValueError, RuntimeError, OSError, json.JSONDecodeError) as error:
        # Diagnostics identify the failed operation, never the private input values.
        raise SystemExit(str(error)) from None
