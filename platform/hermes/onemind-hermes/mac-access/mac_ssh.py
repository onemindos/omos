"""Strict Mac SSH client. Never copies or prints private-key contents."""
import argparse
import base64
import fcntl
import hashlib
import hmac
import os
from pathlib import Path
import re
import stat
import sys

HOST = 'omos-mbp.hermes.svc.cluster.local'
MOUNT = Path('/var/run/legacy-mac-ssh')


def ssh_command(credentials, user, remote_command):
    if not re.fullmatch(r'[a-z_][a-z0-9_-]*', user):
        raise ValueError('invalid Mac account name')
    credentials = Path(credentials).resolve()
    key = credentials / 'ssh-privatekey'
    if stat.S_IMODE(key.stat().st_mode) != 0o400:
        raise ValueError('private key must have mode 0400')
    for name in ('known_hosts', 'ssh-publickey'):
        if not (credentials / name).is_file():
            raise ValueError('credential mount is incomplete')
    options = [
        'BatchMode=yes', 'IdentitiesOnly=yes', 'IdentityAgent=none',
        'StrictHostKeyChecking=yes', 'UserKnownHostsFile=' + str(credentials / 'known_hosts'),
        'GlobalKnownHostsFile=/dev/null', 'HostKeyAlgorithms=ssh-ed25519',
        'UpdateHostKeys=no', 'VerifyHostKeyDNS=no', 'CheckHostIP=no',
        'PasswordAuthentication=no', 'KbdInteractiveAuthentication=no',
        'PreferredAuthentications=publickey', 'ForwardAgent=no',
        'ForwardX11=no', 'ClearAllForwardings=yes', 'RequestTTY=no',
        'PermitLocalCommand=no', 'ProxyCommand=none', 'ProxyJump=none',
        'ControlMaster=no', 'ControlPath=none', 'ConnectTimeout=10',
        'ConnectionAttempts=1', 'ServerAliveInterval=15', 'ServerAliveCountMax=2',
    ]
    command = ['ssh', '-F', '/dev/null', '-T', '-i', str(key)]
    for option in options:
        command.extend(['-o', option])
    return command + [user + '@' + HOST] + list(remote_command)


RESTRICTIONS = 'restrict,no-agent-forwarding,no-port-forwarding,no-pty,no-user-rc,no-X11-forwarding'


def public_key(text):
    parts = text.strip().split()
    if len(parts) < 2 or parts[0] != 'ssh-ed25519' or '\n' in text.strip():
        raise ValueError('expected a single ED25519 public key')
    raw = base64.b64decode(parts[1], validate=True)
    if len(raw) != 51 or raw[:19] != b'\x00\x00\x00\x0bssh-ed25519\x00\x00\x00\x20':
        raise ValueError('invalid ED25519 public key')
    return 'ssh-ed25519 ' + base64.b64encode(raw).decode()


def fingerprint(text):
    raw = base64.b64decode(public_key(text).split()[1])
    return 'SHA256:' + base64.b64encode(hashlib.sha256(raw).digest()).decode().rstrip('=')


def enroll(public, expected, host_public, home):
    # The CLI fixes host_public to the Mac's local sshd key. No network TOFU.
    dedicated = public_key(public)
    if not hmac.compare_digest(fingerprint(Path(host_public).read_text()), expected):
        raise ValueError('Mac host fingerprint mismatch; no enrollment performed')
    home = Path(home)
    ssh_dir = home / '.ssh'
    if ssh_dir.is_symlink():
        raise ValueError('refusing symlink SSH directory')
    ssh_dir.mkdir(mode=0o700, exist_ok=True)
    directory_stat = ssh_dir.stat()
    if directory_stat.st_uid != os.getuid() or directory_stat.st_mode & 0o022:
        raise ValueError('unsafe SSH directory ownership or permissions')
    target = ssh_dir / 'authorized_keys'
    flags = os.O_RDWR | os.O_APPEND | os.O_CREAT | os.O_NOFOLLOW
    fd = os.open(target, flags, 0o600)
    with os.fdopen(fd, 'r+b') as handle:
        fcntl.flock(handle, fcntl.LOCK_EX)
        info = os.fstat(handle.fileno())
        if not stat.S_ISREG(info.st_mode) or info.st_uid != os.getuid() or info.st_mode & 0o077:
            raise ValueError('unsafe authorized_keys ownership or permissions')
        contents = handle.read()
        entry = RESTRICTIONS + ' ' + dedicated + ' legacy-mac-access'
        for line in contents.decode('utf-8').splitlines():
            if line.lstrip().startswith('#'):
                continue
            if dedicated in line:
                if line == entry:
                    return 'already enrolled'
                raise ValueError('dedicated key already present with different restrictions; inspect locally')
        handle.write((('' if not contents or contents.endswith(b'\n') else '\n') + entry + '\n').encode())
        handle.flush()
        os.fsync(handle.fileno())
    return 'enrolled'


def main():
    if sys.argv[1:2] == ['enroll']:
        parser = argparse.ArgumentParser(description='Run locally as the Mac account, never sudo.')
        parser.add_argument('--public-key-file', type=Path, required=True)
        parser.add_argument('--expected-host-fingerprint', required=True)
        args = parser.parse_args(sys.argv[2:])
        try:
            status = enroll(args.public_key_file.read_text(), args.expected_host_fingerprint,
                            Path('/etc/ssh/ssh_host_ed25519_key.pub'), Path.home())
            print(status)
            return 0
        except (OSError, ValueError, UnicodeError):
            print('Enrollment refused: verify local host fingerprint, public key and SSH file safety.', file=sys.stderr)
            return 1

    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--credentials', type=Path, default=MOUNT)
    parser.add_argument('--user', default='zeus')
    parser.add_argument('command', nargs=argparse.REMAINDER)
    args = parser.parse_args()
    try:
        command = args.command
        if command[:1] == ['--']:
            command = command[1:]
        os.execvp('ssh', ssh_command(args.credentials, args.user, command))
    except (OSError, ValueError):
        print('Mac SSH refused: check mount, permissions and account.', file=sys.stderr)
        return 1


if __name__ == '__main__':
    sys.exit(main())
