#!/usr/bin/env python3
"""Seal private IoT credentials without printing plaintext or changing the cluster."""
from __future__ import annotations
import argparse
import base64
import json
from pathlib import Path
import re
import subprocess
import sys
import yaml

ROOT = Path(__file__).resolve().parents[1]
REQUIRED = {
    'openremote': {
        'openremote-env': {'DB_PASSWORD', 'ADMIN_PASSWORD'},
        'openremote-database-restic': {'RESTIC_REPOSITORY', 'RESTIC_PASSWORD', 'AWS_ACCESS_KEY_ID', 'AWS_SECRET_ACCESS_KEY'},
        'openremote-storage-restic': {'RESTIC_REPOSITORY', 'RESTIC_PASSWORD', 'AWS_ACCESS_KEY_ID', 'AWS_SECRET_ACCESS_KEY'},
    },
    'thingsboard': {
        'thingsboard-env': {'DB_PASSWORD', 'SYSADMIN_EMAIL', 'SYSADMIN_PASSWORD', 'JWT_TOKEN_SIGNING_KEY'},
        'thingsboard-database-restic': {'RESTIC_REPOSITORY', 'RESTIC_PASSWORD', 'AWS_ACCESS_KEY_ID', 'AWS_SECRET_ACCESS_KEY'},
    },
}


def validate(app: str, supplied: dict) -> dict:
    if set(supplied) != set(REQUIRED[app]):
        raise ValueError('Private input must contain exactly the documented Secret names')
    payloads = {}
    for name, required in REQUIRED[app].items():
        values = supplied[name]
        if set(values) != required or not all(isinstance(v, str) and v.strip() for v in values.values()):
            raise ValueError('Private input has missing, extra or empty keys')
        for key in required & {'DB_PASSWORD', 'ADMIN_PASSWORD', 'SYSADMIN_PASSWORD', 'RESTIC_PASSWORD'}:
            if len(values[key]) < 24:
                raise ValueError('Passwords require at least 24 characters')
        payloads[name] = dict(values)
    repos = [v['RESTIC_REPOSITORY'] for v in payloads.values() if 'RESTIC_REPOSITORY' in v]
    if len(repos) != len(set(repos)) or any(not r.startswith('s3:https://') for r in repos):
        raise ValueError('Each volume requires its own HTTPS S3 Restic repository prefix')
    if app == 'thingsboard':
        values = payloads['thingsboard-env']
        if not re.fullmatch(r'[^\s@]+@[^\s@]+', values['SYSADMIN_EMAIL']):
            raise ValueError('System administrator email is invalid')
        try:
            key = base64.b64decode(values['JWT_TOKEN_SIGNING_KEY'], validate=True)
        except (ValueError, TypeError) as exc:
            raise ValueError('JWT key must be base64 encoded') from exc
        if len(key) < 32:
            raise ValueError('JWT key requires at least 32 random bytes')
        # BCrypt is the encoder used by CE4.3. No plaintext administrator password enters a Kubernetes Secret.
        import bcrypt
        password = values.pop('SYSADMIN_PASSWORD').encode()
        if len(password) > 72:
            raise ValueError('Administrator password must fit BCrypt 72-byte maximum')
        values['SYSADMIN_PASSWORD_HASH'] = bcrypt.hashpw(password, bcrypt.gensalt(rounds=12)).decode()
    return payloads


def seal(app: str, name: str, values: dict, cert: Path) -> str:
    resource = {'apiVersion': 'v1', 'kind': 'Secret', 'metadata': {'name': name, 'namespace': app},
                'type': 'Opaque', 'stringData': values}
    result = subprocess.run(['kubeseal', '--cert', str(cert), '--scope', 'strict', '--format', 'yaml'],
                            input=json.dumps(resource), text=True, capture_output=True)
    if result.returncode:
        raise RuntimeError('Sealing failed; private diagnostics intentionally suppressed')
    try:
        sealed = yaml.safe_load(result.stdout)
        meta = sealed['metadata']
        encrypted = sealed['spec']['encryptedData']
        template = sealed['spec'].get('template', {})
        if (sealed['kind'] != 'SealedSecret' or meta['name'] != name or meta['namespace'] != app
                or set(encrypted) != set(values)
                or not all(isinstance(v, str) and v for v in encrypted.values())
                or any(k in sealed or k in template for k in ('data', 'stringData'))):
            raise ValueError('Invalid sealed result')
    except (ValueError, TypeError, KeyError, yaml.YAMLError) as exc:
        raise RuntimeError('Invalid sealing output; private diagnostics intentionally suppressed') from exc
    return result.stdout


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('app', choices=REQUIRED)
    parser.add_argument('--private-json', type=Path, required=True)
    parser.add_argument('--cert', type=Path, required=True)
    args = parser.parse_args()
    if args.private_json.stat().st_mode & 0o077:
        raise ValueError('Private input must be readable only by its owner (mode0600)')
    if args.private_json.resolve().is_relative_to(ROOT):
        raise ValueError('Private input must live outside this repository')
    values = validate(args.app, json.loads(args.private_json.read_text()))
    target = ROOT / 'platform' / args.app / 'sealed'
    outputs = {target / (name + '.yaml'): seal(args.app, name, data, args.cert) for name, data in values.items()}
    if any(p.exists() for p in outputs):
        raise ValueError('Refuse overwriting existing SealedSecrets; rotation needs separate review')
    for path, text in outputs.items():
        path.write_text(text)
    # Existing directory is deliberately empty until this helper succeeds for every Secret.
    text = 'apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\nresources:\n'
    text += ''.join('  - ' + p.name + '\n' for p in outputs)
    (target / 'kustomization.yaml').write_text(text)
    print('Wrote encrypted SealedSecrets only; no cluster changes. Preserve private input in recovery storage.')


if __name__ == '__main__':
    try:
        main()
    except (ValueError, RuntimeError, OSError, ImportError, json.JSONDecodeError) as exc:
        # Never echo input JSON, subprocess output, usernames, or credential values.
        print(type(exc).__name__ + ': secret preparation failed; check documented prerequisites.', file=sys.stderr)
        raise SystemExit(1)
