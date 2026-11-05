"""Prepare the dedicated Mac credentials in RAM; never log credential data."""
import os
from pathlib import Path
import pwd
import stat
import sys

SOURCE = Path('/var/run/legacy-mac-ssh-input')
DESTINATION = Path('/var/run/legacy-mac-ssh')
FILES = ('ssh-privatekey', 'ssh-publickey', 'known_hosts')
MAX_BYTES = 65536


def checked_directory(path, owners, *, initial_volume=False):
    info = path.lstat()
    if (not stat.S_ISDIR(info.st_mode) or info.st_uid not in owners
            or (not initial_volume and info.st_mode & 0o022)):
        raise ValueError('unsafe directory')


def read_file(path, owners, modes, *, allow_empty=False):
    fd = os.open(path, os.O_RDONLY | os.O_NOFOLLOW | os.O_NONBLOCK)
    with os.fdopen(fd, 'rb') as handle:
        info = os.fstat(handle.fileno())
        if (not stat.S_ISREG(info.st_mode) or info.st_uid not in owners
                or stat.S_IMODE(info.st_mode) not in modes or info.st_nlink != 1
                or not (0 if allow_empty else 1) <= info.st_size <= MAX_BYTES):
            raise ValueError('unsafe credential metadata')
        value = handle.read(MAX_BYTES + 1)
        if len(value) != info.st_size:
            raise ValueError('credential changed')
        return value


def prepare(source, destination, uid, gid, *, trusted_uid=0):
    """Tests may emulate root ownership; production trusts root only."""
    if uid <= 0 or gid <= 0:
        raise ValueError('nonroot account required')
    source, destination = Path(source), Path(destination)
    # A fresh Kubernetes tmpfs can initially be 0777/01777. Establish ownership
    # and absence of payloads before tightening it, rather than rejecting it.
    checked_directory(destination, {trusted_uid, uid}, initial_volume=True)
    entries = {p.name for p in destination.iterdir()}
    allowed = set(FILES) | {'.' + name + '.new' for name in FILES}
    if not entries <= allowed:
        raise ValueError('unexpected destination entries')
    if entries:
        checked_directory(destination, {trusted_uid, uid})
    os.chmod(destination, 0o700)
    # A read-only Kubernetes Secret volume may retain writable root mode bits;
    # permit those only when the kernel reports the source filesystem read-only.
    readonly = bool(os.statvfs(source).f_flag & os.ST_RDONLY)
    checked_directory(source, {trusted_uid}, initial_volume=readonly)
    present = [(source / name).exists() or (source / name).is_symlink() for name in FILES]
    if not any(present):
        source_entries = list(source.iterdir())
        if source_entries:
            # Optional absent Secrets still receive Kubernetes AtomicWriter's
            # empty timestamp directory and ..data link. Metadata alone is not
            # a partial credential payload. Validate its shape before skipping.
            root = source.resolve(strict=True)
            data = source / '..data'
            if not data.is_symlink():
                raise ValueError('unrecognized empty projection')
            version = data.resolve(strict=True)
            if version.parent != root:
                raise ValueError('empty projection escapes Secret volume')
            checked_directory(version, {trusted_uid})
            if any(version.iterdir()):
                raise ValueError('projection contains unexpected payload')
            for entry in source_entries:
                if entry.name == '..data':
                    continue
                if not entry.name.startswith('..'):
                    raise ValueError('projection contains unexpected entry')
                checked_directory(entry, {trusted_uid})
        if entries:
            raise ValueError('cannot verify existing credentials without source')
        return False
    if not all(present):
        raise ValueError('partial credential projection')
    root = source.resolve(strict=True)
    payloads = {}
    for name in FILES:
        resolved = (source / name).resolve(strict=True)
        if root not in resolved.parents:
            raise ValueError('source escapes Secret volume')
        parent = resolved.parent
        while parent != root:
            checked_directory(parent, {trusted_uid})
            parent = parent.parent
        payloads[name] = read_file(resolved, {trusted_uid}, {0o400})

    # Repeated init execution must succeed after later startup steps fail. A
    # complete destination is accepted only when metadata AND every byte match
    # the current trusted source. Never replace an unknown or different key.
    for name, payload in payloads.items():
        target = destination / name
        staged = destination / ('.' + name + '.new')
        if target.exists() or target.is_symlink():
            if read_file(target, {uid}, {0o400}) != payload:
                raise ValueError('existing credential differs')
            if staged.exists() or staged.is_symlink():
                raise ValueError('ambiguous staged credential')
            continue
        if staged.exists() or staged.is_symlink():
            previous = read_file(staged, {trusted_uid, uid}, {0o600, 0o400}, allow_empty=True)
            if not payload.startswith(previous):
                raise ValueError('staged credential differs')
            if previous != payload:
                # Only incomplete, init-owned 0600 staging is resumable. Append
                # the missing suffix; never truncate, unlink or expose it.
                info = staged.lstat()
                if info.st_uid != trusted_uid or stat.S_IMODE(info.st_mode) != 0o600:
                    raise ValueError('unsafe partial credential')
                fd = os.open(staged, os.O_WRONLY | os.O_APPEND | os.O_NOFOLLOW)
                with os.fdopen(fd, 'wb') as handle:
                    handle.write(payload[len(previous):])
                    handle.flush()
                    os.fsync(handle.fileno())
        else:
            fd = os.open(staged, os.O_WRONLY | os.O_CREAT | os.O_EXCL | os.O_NOFOLLOW, 0o600)
            with os.fdopen(fd, 'wb') as handle:
                handle.write(payload)
                handle.flush()
                os.fsync(handle.fileno())
        fd = os.open(staged, os.O_RDONLY | os.O_NOFOLLOW)
        try:
            os.fchown(fd, uid, gid)
            os.fchmod(fd, 0o400)
            os.fsync(fd)
        finally:
            os.close(fd)
        os.rename(staged, target)
    directory_fd = os.open(destination, os.O_RDONLY | os.O_DIRECTORY | os.O_NOFOLLOW)
    try:
        os.fsync(directory_fd)
        os.fchown(directory_fd, uid, gid)
    finally:
        os.close(directory_fd)
    return True


def main():
    try:
        if os.geteuid() != 0:
            raise ValueError('root init required')
        owner = pwd.getpwnam('hermes')
        ready = prepare(SOURCE, DESTINATION, owner.pw_uid, owner.pw_gid)
    except (OSError, ValueError, KeyError, RuntimeError):
        print('Mac SSH preparation refused: check account and volume safety.', file=sys.stderr)
        return 1
    print('Mac SSH credentials prepared in RAM.' if ready else
          'Mac SSH unavailable: optional Secret absent.')
    return 0


if __name__ == '__main__':
    sys.exit(main())
