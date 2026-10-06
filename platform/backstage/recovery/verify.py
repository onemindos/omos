import base64, io, json, os, re, subprocess, tarfile
import boto3
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

# Recovery-only target. Never accept the production database host.
if os.environ.get('PGHOST') != 'backstage-recovery':
    raise RuntimeError('Recovery requires the dedicated temporary instance')
client = boto3.client('s3', endpoint_url=os.environ['S3_ENDPOINT'], aws_access_key_id=os.environ['S3_ACCESS_KEY'], aws_secret_access_key=os.environ['S3_SECRET_KEY'], region_name='nyc3')
bucket = os.environ['S3_BUCKET']
objects = []
for page in client.get_paginator('list_objects_v2').paginate(Bucket=bucket, Prefix='backups/backstage/'):
    objects.extend(x for x in page.get('Contents', []) if x['Key'].endswith('/databases.tar.gz.aesgcm'))
if not objects:
    raise RuntimeError('No encrypted backup found')
latest = max(objects, key=lambda x: x['LastModified'])
payload = client.get_object(Bucket=bucket, Key=latest['Key'])['Body'].read()
if payload[:9] != b'OMBACKUP1':
    raise RuntimeError('Unexpected backup format')
key = base64.b64decode(os.environ['BACKUP_ENCRYPTION_KEY'], validate=True)
clear = AESGCM(key).decrypt(payload[9:21], payload[21:], b'OneMind Backstage')
results = []
with tarfile.open(fileobj=io.BytesIO(clear), mode='r:gz') as archive:
    for member in archive.getmembers():
        if not member.isfile() or not re.fullmatch(r'backstage_plugin_[a-z0-9_-]+\.dump', member.name):
            raise RuntimeError('Unexpected backup member')
        name = member.name[:-5]
        subprocess.run(['createdb', name], check=True)
        subprocess.run(['pg_restore', '--exit-on-error', '--no-owner', '--no-acl', '-d', name], input=archive.extractfile(member).read(), check=True)
        tables = subprocess.check_output(['psql', '-d', name, '-At', '-v', 'ON_ERROR_STOP=1', '-c', "SELECT count(*) FROM information_schema.tables WHERE table_schema='public'"], text=True).strip()
        results.append({'database': name, 'tables': int(tables)})
catalog = next(x for x in results if x['database'] == 'backstage_plugin_catalog')
if catalog['tables'] == 0:
    raise RuntimeError('Empty restored catalog')
entities = subprocess.check_output(['psql', '-d', 'backstage_plugin_catalog', '-At', '-v', 'ON_ERROR_STOP=1', '-c', 'SELECT count(*) FROM final_entities'], text=True).strip()
if int(entities) == 0:
    raise RuntimeError('No recovered catalog entities')
print(json.dumps({'object': latest['Key'], 'backupTime': latest['LastModified'].isoformat(), 'isolatedInstance': True, 'restoreVerified': True, 'databases': results, 'catalogEntities': int(entities)}))
