#!/usr/bin/env python3
"""Install only the checksum-pinned public CLI into a pod-scoped volume."""
import argparse
import hashlib
import io
from pathlib import Path
import platform
import os
import pwd
import urllib.request
import zipfile

VERSION = '2.41.1'
SHA256 = 'e52476700cba9202ac7469bdfb2648ecdbf3438ac85822010f77230139e74966'
URL = f'https://cache.agilebits.com/dist/1P/op2/pkg/v{VERSION}/op_linux_amd64_v{VERSION}.zip'


def install(destination, arch=None, fetch=None):
    # No silent amd64-on-arm fallback. Add another architecture only with a
    # separately verified official artifact/checksum and an execution test.
    if (arch or platform.machine()) not in ('x86_64', 'amd64'):
        raise RuntimeError('Pinned 1Password CLI supports linux/amd64 only')
    if fetch is None:
        with urllib.request.urlopen(URL, timeout=60) as response:
            archive = response.read()
    else:
        archive = fetch()
    if hashlib.sha256(archive).hexdigest() != SHA256:
        raise RuntimeError('1Password CLI archive checksum mismatch')
    # Extract exactly one known member, never arbitrary archive paths.
    with zipfile.ZipFile(io.BytesIO(archive)) as bundle:
        binary = bundle.read('op')
    destination = Path(destination)
    destination.mkdir(parents=True, exist_ok=True)
    staged = destination / 'op.new'
    staged.write_bytes(binary)
    staged.chmod(0o555)
    staged.replace(destination / 'op')
    print(f'1Password CLI {VERSION}: verified and installed (linux/amd64)')


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--destination', type=Path, default=Path('/opt/onepassword'))
    parser.add_argument('--state-dir', type=Path, default=Path('/var/run/onepassword'))
    args = parser.parse_args()
    install(args.destination)
    # The official image launcher runs the gateway as hermes, not root.
    args.state_dir.mkdir(parents=True, exist_ok=True)
    args.state_dir.chmod(0o700)
    if os.geteuid() == 0:
        owner = pwd.getpwnam('hermes')
        os.chown(args.state_dir, owner.pw_uid, owner.pw_gid)
