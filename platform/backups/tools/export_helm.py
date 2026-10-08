"""Recover protected Helm bootstrap inputs from a verified offline archive."""
import argparse
import base64
import gzip
import json
import os
from pathlib import Path
import re
import tarfile
import tempfile
import yaml
from decrypt import decrypt


def export(record, output):
    releases = {}
    for item in record['items']:
        if item.get('kind') != 'Secret' or item.get('type') != 'helm.sh/release.v1':
            continue
        payload = base64.b64decode(base64.b64decode(item['data']['release']))
        if payload.startswith(b'\x1f\x8b'):
            payload = gzip.decompress(payload)
        release = json.loads(payload)
        if release.get('info', {}).get('status') != 'deployed':
            continue
        identity = (item['metadata']['namespace'], release['name'])
        if not all(re.fullmatch(r'[a-z0-9][a-z0-9.-]*', value) for value in identity):
            raise ValueError('Invalid release identity')
        if identity not in releases or release['version'] > releases[identity]['version']:
            releases[identity] = release
    if output.exists():
        raise FileExistsError('Refusing to overwrite bootstrap inputs')
    output.mkdir(mode=0o700, parents=True)
    inventory = []
    for (namespace, name), release in sorted(releases.items()):
        destination = output / namespace / name
        destination.mkdir(mode=0o700, parents=True)
        (destination / 'values.yaml').write_text(yaml.safe_dump(release.get('config', {}), sort_keys=False))
        (destination / 'manifest.yaml').write_text(release['manifest'])
        for filename in ['values.yaml', 'manifest.yaml']:
            (destination / filename).chmod(0o600)
        chart = release['chart']['metadata']
        inventory.append({'namespace': namespace, 'release': name, 'revision': release['version'],
                          'chart': chart['name'], 'chartVersion': chart['version']})
    (output / 'inventory.json').write_text(json.dumps(inventory, indent=2) + '\n')
    (output / 'inventory.json').chmod(0o600)
    return inventory


def run():
    parser = argparse.ArgumentParser()
    parser.add_argument('archive', type=Path)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--key-file', type=Path, default=Path.home() / '.config/onemind/recovery/backup-encryption-key')
    args = parser.parse_args()
    os.umask(0o077)
    with tempfile.TemporaryDirectory() as temporary:
        clear = Path(temporary) / 'recovery.tar.gz'
        decrypt(args.archive, clear, base64.b64decode(args.key_file.read_text().strip(), validate=True))
        with tarfile.open(clear) as tar:
            record = json.load(tar.extractfile('data/recovery.json'))
        inventory = export(record, args.output)
    print(json.dumps({'offlineHelmInputsExported': len(inventory), 'output': str(args.output), 'resourcesApplied': False}))


if __name__ == '__main__':
    try:
        run()
    except Exception as error:
        print(json.dumps({'exportFailed': True, 'errorType': type(error).__name__}))
        raise SystemExit(1)
