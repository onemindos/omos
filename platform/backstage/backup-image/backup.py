import base64, datetime, io, json, os, subprocess, tarfile
import boto3
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

def sql(query):
    return subprocess.check_output(['psql', '-d', 'postgres', '-At', '-v', 'ON_ERROR_STOP=1', '-c', query], text=True).strip()

def run():
    databases = sql("SELECT datname FROM pg_database WHERE datdba=(SELECT oid FROM pg_roles WHERE rolname=current_user) AND NOT datistemplate ORDER BY datname").splitlines()
    if not databases:
        raise RuntimeError('No application databases found; refusing empty backup')
    archive = io.BytesIO()
    with tarfile.open(fileobj=archive, mode='w:gz') as tar:
        for database in databases:
            dump = subprocess.check_output(['pg_dump', '-Fc', '-d', database])
            member = tarfile.TarInfo(database + '.dump')
            member.size = len(dump)
            tar.addfile(member, io.BytesIO(dump))
    key = base64.b64decode(os.environ['BACKUP_ENCRYPTION_KEY'], validate=True)
    nonce = os.urandom(12)
    payload = b'OMBACKUP1' + nonce + AESGCM(key).encrypt(nonce, archive.getvalue(), b'OneMind Backstage')
    stamp = datetime.datetime.now(datetime.timezone.utc).strftime('%Y%m%dT%H%M%SZ')
    object_key = 'backups/backstage/' + stamp + '/databases.tar.gz.aesgcm'
    client = boto3.client('s3', endpoint_url=os.environ['S3_ENDPOINT'], aws_access_key_id=os.environ['S3_ACCESS_KEY'], aws_secret_access_key=os.environ['S3_SECRET_KEY'], region_name='nyc3')
    bucket = os.environ['S3_BUCKET']
    client.put_object(Bucket=bucket, Key=object_key, Body=payload)
    # Verify the remote bytes, authentication tag, and a real restore into an isolated temporary database.
    remote = client.get_object(Bucket=bucket, Key=object_key)['Body'].read()
    clear = AESGCM(key).decrypt(remote[9:21], remote[21:], b'OneMind Backstage')
    with tarfile.open(fileobj=io.BytesIO(clear), mode='r:gz') as tar:
        member = next(m for m in tar.getmembers() if 'catalog' in m.name)
        restored = tar.extractfile(member).read()
    temporary = 'onemind_restore_' + stamp.lower()
    subprocess.run(['createdb', temporary], check=True)
    try:
        subprocess.run(['pg_restore', '--exit-on-error', '--no-owner', '-d', temporary], input=restored, check=True)
        count = subprocess.check_output(['psql', '-d', temporary, '-At', '-c', "SELECT count(*) FROM information_schema.tables WHERE table_schema='public'"], text=True).strip()
        if int(count) == 0:
            raise RuntimeError('Restored catalog contains no tables')
    finally:
        subprocess.run(['dropdb', '--if-exists', temporary], check=True)
    print(json.dumps({'backup': object_key, 'databases': len(databases), 'restoreVerified': True, 'catalogTables': int(count)}))

if __name__ == '__main__':
    run()
