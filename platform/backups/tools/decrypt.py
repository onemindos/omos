"""Offline archive decryption; create a new protected file, never apply resources."""
import argparse
import base64
import json
import os
from pathlib import Path
import tempfile
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes

FORMATS = {b'OMBKP002': b'OneMindPlatformBackupV1',
           b'OMBACKUP1': b'OneMindBackstage',
           b'OMBSECR1': b'OneMindClusterSecretRecoveryV1'}


def decrypt(source, output, key):
    if len(key) != 32:
        raise ValueError('Expected a 32-byte key')
    if output.exists():
        raise FileExistsError('Refusing to replace an existing output')
    descriptor, temporary = tempfile.mkstemp(dir=output.parent, prefix='.recovery-decrypt-')
    try:
        with source.open('rb') as src, os.fdopen(descriptor, 'wb') as dst:
            prefix = src.read(9)
            header = next((candidate for candidate in FORMATS if prefix.startswith(candidate)), None)
            if header is None:
                raise ValueError('Unknown archive format')
            src.seek(len(header))
            nonce = src.read(12)
            remaining = source.stat().st_size - len(header) - 12 - 16
            if remaining < 0:
                raise ValueError('Truncated archive')
            src.seek(-16, 2)
            tag = src.read(16)
            src.seek(len(header) + 12)
            cipher = Cipher(algorithms.AES(key), modes.GCM(nonce, tag)).decryptor()
            cipher.authenticate_additional_data(FORMATS[header])
            while remaining:
                chunk = src.read(min(1024 * 1024, remaining))
                if not chunk:
                    raise ValueError('Truncated archive')
                dst.write(cipher.update(chunk))
                remaining -= len(chunk)
            dst.write(cipher.finalize())
        os.link(temporary, output)
    finally:
        Path(temporary).unlink(missing_ok=True)


def run():
    parser = argparse.ArgumentParser()
    parser.add_argument('archive', type=Path)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--key-file', type=Path, default=Path.home() / '.config/onemind/recovery/backup-encryption-key')
    args = parser.parse_args()
    decrypt(args.archive, args.output, base64.b64decode(args.key_file.read_text().strip(), validate=True))
    print(json.dumps({'archiveAuthenticated': True, 'output': str(args.output), 'resourcesApplied': False}))


if __name__ == '__main__':
    try:
        run()
    except Exception as error:
        print(json.dumps({'decryptionFailed': True, 'errorType': type(error).__name__}))
        raise SystemExit(1)
