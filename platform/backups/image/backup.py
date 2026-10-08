"""Encrypted, fail-closed application backups with remote restore verification."""
import base64
import datetime
import hashlib
import json
import os
from pathlib import Path
import re
import socket
import subprocess
import tarfile
import tempfile
import time
import uuid

import boto3
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes

HEADER = b"OMBKP002"
AAD = b"OneMindPlatformBackupV1"
CHUNK = 1024 * 1024


def command(args, env=None):
    # Never log argv/stderr: provider errors and URLs can contain credentials.
    result = subprocess.run(args, env=env, capture_output=True, check=False)
    if result.returncode:
        raise RuntimeError(f"{Path(args[0]).name} failed (exit {result.returncode})")
    return result.stdout


def digest(path):
    value = hashlib.sha256()
    with path.open("rb") as source:
        for chunk in iter(lambda: source.read(CHUNK), b""):
            value.update(chunk)
    return value.hexdigest()


def encrypt(source, target, key):
    nonce = os.urandom(12)
    cipher = Cipher(algorithms.AES(key), modes.GCM(nonce)).encryptor()
    cipher.authenticate_additional_data(AAD)
    with source.open("rb") as src, target.open("wb") as dst:
        dst.write(HEADER + nonce)
        for chunk in iter(lambda: src.read(CHUNK), b""):
            dst.write(cipher.update(chunk))
        dst.write(cipher.finalize())
        dst.write(cipher.tag)


def decrypt(source, target, key):
    try:
        with source.open("rb") as src, target.open("wb") as dst:
            if src.read(len(HEADER)) != HEADER:
                raise ValueError("Invalid backup format")
            nonce = src.read(12)
            src.seek(-16, 2)
            tag = src.read(16)
            remaining = source.stat().st_size - len(HEADER) - 12 - 16
            if remaining < 0:
                raise ValueError("Truncated backup")
            src.seek(len(HEADER) + 12)
            cipher = Cipher(algorithms.AES(key), modes.GCM(nonce, tag)).decryptor()
            cipher.authenticate_additional_data(AAD)
            while remaining:
                chunk = src.read(min(CHUNK, remaining))
                if not chunk:
                    raise ValueError("Truncated backup")
                dst.write(cipher.update(chunk))
                remaining -= len(chunk)
            dst.write(cipher.finalize())
    except Exception:
        target.unlink(missing_ok=True)
        raise


def postgres_backup(directory):
    database = os.environ["PGDATABASE"]
    command(["pg_dump", "--format=custom", "--file", str(directory / "database.dump"), database])
    return {"kind": "postgres", "database": database}


def postgres_restore(directory, metadata):
    # A separate loopback-only sidecar supplies restore storage and extensions.
    # Never consume production PGHOST, PGHOSTADDR, PGPASSWORD or PGDATABASE here.
    env = {k: v for k, v in os.environ.items() if not k.startswith("PG")}
    env.update(PGHOST="127.0.0.1", PGPORT=os.environ.get("PG_RESTORE_PORT", "15432"),
               PGUSER="postgres", PGCONNECT_TIMEOUT="15")
    name = "onemind_restore_" + uuid.uuid4().hex
    command(["createdb", "--maintenance-db=postgres", name], env=env)
    try:
        command(["pg_restore", "--exit-on-error", "--no-owner", "--no-privileges",
                 "--dbname", name, str(directory / "database.dump")], env=env)
        tables = int(command(["psql", "--dbname", name, "-At", "-v", "ON_ERROR_STOP=1",
                              "-c", "SELECT count(*) FROM information_schema.tables "
                              "WHERE table_schema NOT IN ('pg_catalog','information_schema')"], env=env))
        if tables < 1:
            raise RuntimeError("Restored database has no application tables")
        return {"tables": tables}
    finally:
        command(["dropdb", "--maintenance-db=postgres", "--if-exists", name], env=env)


def nats_backup(directory):
    streams = command(["nats", "--no-context", "stream", "ls", "--names", "--all"]).decode().splitlines()
    if not streams:
        raise RuntimeError("Refusing an empty NATS backup")
    for stream in streams:
        if not re.fullmatch(r"[A-Za-z0-9_-]+", stream):
            raise ValueError("Unexpected stream name")
        command(["nats", "--no-context", "backup", "stream", "--check", "--consumers",
                 "--no-progress", stream, str(directory / stream)])
    return {"kind": "nats", "streams": streams}


def nats_restore(directory, metadata):
    # No production credentials or routing in the restore process. Bind loopback.
    env = {k: v for k, v in os.environ.items() if not k.startswith("NATS_")}
    env["NATS_URL"] = "nats://127.0.0.1:14222"
    with tempfile.TemporaryDirectory() as store:
        server = subprocess.Popen(["nats-server", "-js", "-a", "127.0.0.1", "-p", "14222", "-sd", store],
                                  env=env, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        try:
            for _ in range(100):
                if server.poll() is not None:
                    raise RuntimeError("Isolated restore server exited")
                try:
                    with socket.create_connection(("127.0.0.1", 14222), timeout=0.2):
                        break
                except OSError:
                    time.sleep(0.1)
            else:
                raise RuntimeError("Isolated restore server did not start")
            for stream in metadata["streams"]:
                command(["nats", "--no-context", "backup", "restore", "stream", "--replicas=1",
                         "--no-progress", str(directory / stream)], env=env)
                info = json.loads(command(["nats", "--no-context", "stream", "info", stream, "--json"], env=env))
                if info["config"]["name"] != stream:
                    raise RuntimeError("Restored stream identity mismatch")
            return {"streams": len(metadata["streams"])}
        finally:
            server.terminate()
            try:
                server.wait(timeout=10)
            except subprocess.TimeoutExpired:
                server.kill()
                server.wait()


def run():
    kind = os.environ["BACKUP_KIND"]
    if kind not in ("postgres", "nats"):
        raise ValueError("Unknown backup kind")
    key = base64.b64decode(os.environ["BACKUP_ENCRYPTION_KEY"], validate=True)
    if len(key) != 32:
        raise ValueError("Backup encryption key must be 32 bytes")
    prefix = os.environ["BACKUP_PREFIX"].strip("/")
    if not re.fullmatch(r"[A-Za-z0-9/_-]+", prefix):
        raise ValueError("Invalid backup prefix")
    stamp = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    object_key = f"{prefix}/{stamp}-{uuid.uuid4().hex}/archive.tar.gz.aesgcm"
    with tempfile.TemporaryDirectory() as temporary:
        root = Path(temporary)
        data = root / "data"
        data.mkdir()
        metadata = postgres_backup(data) if kind == "postgres" else nats_backup(data)
        (data / "manifest.json").write_text(json.dumps(metadata))
        archive = root / "archive.tar.gz"
        with tarfile.open(archive, "w:gz") as tar:
            tar.add(data, arcname="data")
        encrypted = root / "encrypted"
        encrypt(archive, encrypted, key)
        client = boto3.client("s3", endpoint_url=os.environ["S3_ENDPOINT"],
                              aws_access_key_id=os.environ["S3_ACCESS_KEY"],
                              aws_secret_access_key=os.environ["S3_SECRET_KEY"],
                              region_name=os.environ.get("S3_REGION", "us-east-1"))
        bucket = os.environ["S3_BUCKET"]
        client.upload_file(str(encrypted), bucket, object_key)
        remote = root / "remote"
        client.download_file(bucket, object_key, str(remote))
        if digest(remote) != digest(encrypted):
            raise RuntimeError("Remote backup checksum mismatch")
        verified = root / "verified.tar.gz"
        decrypt(remote, verified, key)
        restored = root / "restored"
        restored.mkdir()
        with tarfile.open(verified, "r:gz") as tar:
            tar.extractall(restored, filter="data")
        metadata = json.loads((restored / "data" / "manifest.json").read_text())
        result = postgres_restore(restored / "data", metadata) if kind == "postgres" else nats_restore(restored / "data", metadata)
        # This marker is written only after download, authentication and real restore.
        client.put_object(Bucket=bucket, Key=object_key + ".verified.json",
                          Body=json.dumps({"kind": kind, "sha256": digest(remote), "restore": result}).encode())
        print(json.dumps({"backup": object_key, "kind": kind, "restoreVerified": True, **result}))


if __name__ == "__main__":
    try:
        run()
    except Exception as error:
        print(json.dumps({"backupFailed": True, "errorType": type(error).__name__}), flush=True)
        raise SystemExit(1)
