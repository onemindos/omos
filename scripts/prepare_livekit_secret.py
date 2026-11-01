#!/usr/bin/env python3
"""Seal an existing private LiveKit API-key map; never apply it or print credentials."""
import argparse
import base64
import json
from pathlib import Path
import subprocess
import yaml

ROOT = Path(__file__).resolve().parents[1]
OUTPUT = ROOT / 'platform/livekit/sealed/livekit-api.sealedsecret.yaml'


def payload(path):
    mapping = yaml.safe_load(path.read_text())
    if not isinstance(mapping, dict) or not mapping:
        raise ValueError('Expected a nonempty API-key map')
    for key, secret in mapping.items():
        if not isinstance(key, str) or not key or not isinstance(secret, str) or len(secret) < 32:
            raise ValueError('Expected string API keys and secrets of at least 32 characters')
    return {'apiVersion': 'v1', 'kind': 'Secret', 'type': 'Opaque',
            'metadata': {'name': 'livekit-api', 'namespace': 'livekit'},
            'data': {'keys.yaml': base64.b64encode(yaml.safe_dump(mapping).encode()).decode()}}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--keys-file', required=True, type=Path)
    parser.add_argument('--cert', required=True, type=Path)
    args = parser.parse_args()
    if OUTPUT.exists():
        raise SystemExit('Output already exists; review rotation separately')
    if args.keys_file.resolve().is_relative_to(ROOT):
        raise SystemExit('Private input must stay outside the repository')
    if args.keys_file.stat().st_mode & 0o077:
        raise SystemExit('Private input must not be accessible to group or others')
    try:
        data = payload(args.keys_file)
    except Exception:
        raise SystemExit('Invalid private key map; no output written') from None
    result = subprocess.run(['kubeseal', '--cert', str(args.cert), '--scope', 'strict', '--format', 'yaml'],
                            input=json.dumps(data), capture_output=True, text=True)
    if result.returncode:
        raise SystemExit('Sealing failed; private subprocess diagnostics withheld')
    sealed = yaml.safe_load(result.stdout)
    if (sealed.get('kind') != 'SealedSecret' or
            sealed.get('metadata', {}).get('name') != 'livekit-api' or
            sealed.get('metadata', {}).get('namespace') != 'livekit' or
            set(sealed.get('spec', {}).get('encryptedData', {})) != {'keys.yaml'}):
        raise SystemExit('Unexpected sealed output; no file written')
    with OUTPUT.open('x') as output:
        output.write(result.stdout)
    print('Encrypted SealedSecret written; add its filename to the sealed kustomization after review.')


if __name__ == '__main__':
    main()
