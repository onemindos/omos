"""Create only strict-bound sealed payloads; never print or write plaintext keys.

The backup input is a private JSON file with RESTIC_REPOSITORY, RESTIC_PASSWORD,
AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY. Keep its recovery copy offline.
No Kubernetes mutation, existing credential lookup or account enrollment occurs.
"""
import argparse
import json
from pathlib import Path
import secrets
import subprocess
import sys


def payloads(backup):
    required = {"RESTIC_REPOSITORY", "RESTIC_PASSWORD", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"}
    if not required <= backup.keys() or not all(isinstance(backup[k], str) and backup[k] for k in required):
        raise ValueError("Backup file requires the four documented non-empty string fields")
    if not backup["RESTIC_REPOSITORY"].startswith("s3:https://") or not backup["RESTIC_REPOSITORY"].rstrip("/").endswith("/backups/restic/forgejo/data"):
        raise ValueError("Use a distinct HTTPS S3 repository ending /backups/restic/forgejo/data")
    if len(backup["RESTIC_PASSWORD"]) < 32:
        raise ValueError("Preserve a dedicated recovery password of at least32 characters")
    keys = {"FORGEJO__security__SECRET_KEY": secrets.token_urlsafe(48),
            "FORGEJO__security__INTERNAL_TOKEN": secrets.token_hex(64),
            "FORGEJO__oauth2__JWT_SECRET": secrets.token_urlsafe(32),
            "FORGEJO__server__LFS_JWT_SECRET": secrets.token_urlsafe(32)}
    return {"forgejo-keys": keys, "forgejo-data-restic": {key: backup[key] for key in required}}


def seal(name, values, cert):
    secret = {"apiVersion": "v1", "kind": "Secret", "metadata": {"name": name, "namespace": "forgejo"},
              "type": "Opaque", "stringData": values}
    result = subprocess.run(["kubeseal", "--cert", str(cert), "--scope", "strict", "--format", "yaml"],
                            input=json.dumps(secret), text=True, capture_output=True)
    if result.returncode:
        raise RuntimeError("Sealing failed; no private subprocess diagnostics exposed")
    return result.stdout


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--cert", type=Path, required=True)
    parser.add_argument("--backup-json", type=Path, required=True)
    parser.add_argument("--output", type=Path, default=Path("platform/forgejo/sealed"))
    args = parser.parse_args()
    try:
        if args.backup_json.stat().st_mode & 0o077:
            raise ValueError("Recovery input must not be group/world accessible")
        values = payloads(json.loads(args.backup_json.read_text()))
        names = [name + ".yaml" for name in values]
        if any((args.output / name).exists() for name in names):
            raise ValueError("Refusing to overwrite existing sealed keys; preserve key identity")
        # Finish every encryption before writing any output.
        encrypted = {name + ".yaml": seal(name, data, args.cert) for name, data in values.items()}
        args.output.mkdir(parents=True, exist_ok=True)
        for name, data in encrypted.items():
            with (args.output / name).open("x") as handle:
                handle.write(data)
        (args.output / "kustomization.yaml").write_text(
            "apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\nresources:\n" +
            "".join("- " + name + "\n" for name in names))
        print("Two namespace-bound SealedSecrets written; review, commit and use Argo for activation")
    except Exception as error:
        print("Secret preparation failed (" + type(error).__name__ + "); no private inputs printed", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
