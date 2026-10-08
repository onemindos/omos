"""Mirror verified encrypted recovery objects without deleting remote/local data."""
import argparse
import base64
import datetime
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import boto3

GIB = 1024 ** 3


def safe_path(root, key):
    target = (root / key).resolve()
    if not target.is_relative_to(root.resolve()) or key.startswith('/'):
        raise ValueError('Object path escapes recovery directory')
    return target


def digest(path):
    value = hashlib.sha256()
    with path.open('rb') as source:
        for chunk in iter(lambda: source.read(1024 * 1024), b''):
            value.update(chunk)
    return value.hexdigest()


def run():
    parser = argparse.ArgumentParser()
    parser.add_argument('--context', default='do-nyc3-omos-cloud')
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--since', default='2026-10-08')
    parser.add_argument('--max-gib', type=int, default=50)
    parser.add_argument('--reserve-gib', type=int, default=40)
    args = parser.parse_args()
    os.umask(0o077)
    k = ['kubectl', '--context', args.context, '--request-timeout=30s']
    secret = json.loads(subprocess.check_output(k + ['-n', 'backstage', 'get', 'secret',
                                                   'backstage-backup', '-o', 'json']))['data']
    env = {key: base64.b64decode(value).decode() for key, value in secret.items()}
    client = boto3.client('s3', endpoint_url=env['S3_ENDPOINT'],
                          aws_access_key_id=env['S3_ACCESS_KEY'],
                          aws_secret_access_key=env['S3_SECRET_KEY'], region_name='nyc3')
    bucket = env['S3_BUCKET']
    sources = json.loads(subprocess.check_output(k + ['get', 'replicationsources', '-A', '-o', 'json']))['items']
    completed = [d for d in sources if d.get('status', {}).get('latestMoverStatus', {}).get('result') == 'Successful']
    empty_volumes = [{'namespace': d['metadata']['namespace'], 'pvc': d['spec']['sourcePVC']}
                     for d in completed if 'Directory is empty' in d.get('status', {}).get('latestMoverStatus', {}).get('logs', '')]
    repositories = sorted('backups/restic/' + d['metadata']['namespace'] + '/' + d['spec']['sourcePVC']
                          for d in completed if {'namespace': d['metadata']['namespace'], 'pvc': d['spec']['sourcePVC']} not in empty_volumes)
    objects = []
    for page in client.get_paginator('list_objects_v2').paginate(Bucket=bucket, Prefix='backups/'):
        objects.extend(page.get('Contents', []))
    selected = {item['Key']: item for item in objects if any(
        item['Key'].startswith(repo + '/') and '/locks/' not in item['Key'] for repo in repositories)}
    expected = {}
    threshold = datetime.datetime.fromisoformat(args.since).replace(tzinfo=datetime.timezone.utc)
    indexed = {item['Key']: item for item in objects}
    for marker in objects:
        if not marker['Key'].endswith('.verified.json') or marker['LastModified'] < threshold:
            continue
        metadata = json.loads(client.get_object(Bucket=bucket, Key=marker['Key'])['Body'].read())
        archive = marker['Key'][:-len('.verified.json')]
        if archive not in indexed or not isinstance(metadata.get('sha256'), str):
            continue
        selected[marker['Key']] = marker
        selected[archive] = indexed[archive]
        expected[archive] = metadata['sha256']
    total = sum(item['Size'] for item in selected.values())
    args.output.mkdir(parents=True, exist_ok=True, mode=0o700)
    missing = sum(item['Size'] for key, item in selected.items()
                  if not safe_path(args.output, key).exists() or safe_path(args.output, key).stat().st_size != item['Size'])
    if total > args.max_gib * GIB or shutil.disk_usage(args.output).free - missing < args.reserve_gib * GIB:
        print(json.dumps({'mirrorBlocked': 'capacity', 'plannedGiB': round(total / GIB, 2),
                          'maxGiB': args.max_gib, 'reservedFreeGiB': args.reserve_gib}))
        return 2
    for key, item in sorted(selected.items()):
        target = safe_path(args.output, key)
        target.parent.mkdir(parents=True, exist_ok=True, mode=0o700)
        if not target.exists() or target.stat().st_size != item['Size']:
            temporary = target.with_name(target.name + '.partial')
            client.download_file(bucket, key, str(temporary))
            if temporary.stat().st_size != item['Size']:
                raise RuntimeError('Object size mismatch')
            temporary.chmod(0o600)
            temporary.replace(target)
        if key in expected and digest(target) != expected[key]:
            raise RuntimeError('Verified archive checksum mismatch')
    result = {'context': args.context, 'createdAt': datetime.datetime.now(datetime.timezone.utc).isoformat(),
              'objects': len(selected), 'encryptedGiB': round(total / GIB, 2),
              'archiveChecksumsVerified': len(expected), 'repositories': repositories,
              'filesystemIntegrityVerified': False, 'confirmedEmptyVolumes': empty_volumes,
              'note': 'Run restic check --read-data on each local repository; mirror equality alone is not an application restore.'}
    (args.output / 'mirror-manifest.json').write_text(json.dumps(result, indent=2) + '\n')
    print(json.dumps({key: value for key, value in result.items() if key != 'repositories'}))
    return 0


if __name__ == '__main__':
    try:
        sys.exit(run())
    except Exception as error:
        print(json.dumps({'mirrorFailed': True, 'errorType': type(error).__name__}))
        sys.exit(1)
