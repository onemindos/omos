#!/usr/bin/env bash
# Synthetic local services only. Never inherits production connection variables.
set -euo pipefail
backup_image="${1:-platform-backup:test}"
test_dir="$(mktemp -d)"
cleanup() {
  docker rm -f onemind-backup-test-pg onemind-backup-test-nats >/dev/null 2>&1 || true
  if [[ -n "${moto_pid:-}" ]]; then kill "$moto_pid" || true; fi
  rm -rf "$test_dir"
}
trap cleanup EXIT
moto_server -H 127.0.0.1 -p 15000 >"$test_dir/moto.log" 2>&1 &
moto_pid=$!
for ((attempt=1; attempt<=60; attempt++)); do
  if curl --silent --fail http://127.0.0.1:15000/ >/dev/null; then break; fi
  sleep 1
done
python - <<'PY'
import boto3
boto3.client('s3', endpoint_url='http://127.0.0.1:15000',
             aws_access_key_id='synthetic', aws_secret_access_key='synthetic',
             region_name='us-east-1').create_bucket(Bucket='backup-tests')
PY
cat >"$test_dir/common.env" <<'ENV'
S3_ENDPOINT=http://127.0.0.1:15000
S3_ACCESS_KEY=synthetic
S3_SECRET_KEY=synthetic
S3_BUCKET=backup-tests
BACKUP_ENCRYPTION_KEY=a2tra2tra2tra2tra2tra2tra2tra2tra2tra2tra2s=
ENV
docker run -d --name onemind-backup-test-pg -p 127.0.0.1:15432:5432 \
  -e POSTGRES_PASSWORD=synthetic -e POSTGRES_DB=fixture \
  postgres:15-alpine@sha256:f7d23353e1b15400d22ebe31189f4d314b87a4c129cc400c8c2d8d4ca127bf81 >/dev/null
for ((attempt=1; attempt<=60; attempt++)); do
  if docker exec onemind-backup-test-pg pg_isready -U postgres >/dev/null; then break; fi
  sleep 1
done
docker exec onemind-backup-test-pg psql -U postgres -d fixture -v ON_ERROR_STOP=1 \
  -c 'CREATE TABLE fixture(id integer primary key); INSERT INTO fixture VALUES (1),(2);' >/dev/null
docker run --rm --network host --env-file "$test_dir/common.env" \
  -e BACKUP_KIND=postgres -e BACKUP_PREFIX=tests/postgres \
  -e PGHOST=127.0.0.1 -e PGPORT=15432 -e PGUSER=postgres -e PGPASSWORD=synthetic \
  -e PGDATABASE=fixture "$backup_image"
docker run -d --name onemind-backup-test-nats -p 127.0.0.1:15222:4222 \
  nats:2.14.2-alpine@sha256:952d157e28d5394a211229bd57a7b37ff9f184e58e2c8486a08fa909fd254e32 \
  -js --user fixture --pass synthetic >/dev/null
sleep 2
docker run --rm --network host --entrypoint nats \
  -e NATS_URL=nats://127.0.0.1:15222 -e NATS_USER=fixture -e NATS_PASSWORD=synthetic \
  "$backup_image" stream add EVENTS --subjects fixture.events --storage file --retention limits \
  --discard old --max-age 0 --max-bytes=-1 --max-msgs=-1 --max-msg-size=-1 --replicas 1 --defaults >/dev/null
docker run --rm --network host --entrypoint nats \
  -e NATS_URL=nats://127.0.0.1:15222 -e NATS_USER=fixture -e NATS_PASSWORD=synthetic \
  "$backup_image" pub fixture.events 'synthetic fixture' >/dev/null
docker run --rm --network host --env-file "$test_dir/common.env" \
  -e BACKUP_KIND=nats -e BACKUP_PREFIX=tests/nats \
  -e NATS_URL=nats://127.0.0.1:15222 -e NATS_USER=fixture -e NATS_PASSWORD=synthetic "$backup_image"
python - <<'PY'
import boto3
s3 = boto3.client('s3', endpoint_url='http://127.0.0.1:15000',
                 aws_access_key_id='synthetic', aws_secret_access_key='synthetic')
objects = s3.list_objects_v2(Bucket='backup-tests')['Contents']
assert sum(x['Key'].endswith('.verified.json') for x in objects) == 2
print('PostgreSQL and NATS encrypted remote restore checks passed')
PY
