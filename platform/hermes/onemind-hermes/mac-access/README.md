# Dedicated Mac SSH bootstrap

Source and offline contracts only; **not deployed or authenticated end-to-end**.
Pending gates: owner-local Mac host verification and enrollment, owner merge,
Argo reconciliation, controller decryption and an authenticated read-only probe.
The encrypted resource is the real dedicated identity supplied for this bootstrap,
not an example. It is strict-scoped to `hermes/onemind-hermes-mac-ssh`.
A time-valid repository controller certificate was used by the parent; equality
with the live controller certificate and successful decryption remain unverified.

## Owner: one-time local enrollment (safe to repeat)

Transfer only `mac_ssh.py` and `legacy-mac-access.pub` from this directory to the
Mac. Review both locally, then run **as the target Mac account, not sudo**:

```sh
python3 mac_ssh.py enroll \
  --public-key-file ./legacy-mac-access.pub \
  --expected-host-fingerprint 'SHA256:q8Al9q/LiMEPgwtN+VMoekF+8oB+TNYqjF3xDHzE7d8'
```

This first reads `/etc/ssh/ssh_host_ed25519_key.pub` on the Mac and compares its
SHA256 fingerprint with the network-observed fingerprint above. A mismatch or
missing local key fails closed **before creating `.ssh` or editing anything**.
Do not replace the expected fingerprint with an unverified scan result. If they
disagree, stop and investigate. Optional independent human-readable check:

```sh
ssh-keygen -lf /etc/ssh/ssh_host_ed25519_key.pub -E sha256
```

The script accepts one ED25519 public key, takes an exclusive lock, and appends
only that key with `restrict,no-agent-forwarding,no-port-forwarding,no-pty,no-user-rc,no-X11-forwarding`.
It preserves existing bytes, adds a separating newline only if needed, and returns
`already enrolled` without writing on repeat. Symlink paths, unsafe ownership or
permissions, and a pre-existing copy of the dedicated key with different options
are refused, not repaired. Existing `authorized_keys` must be owner-only (0600 or
stricter); `.ssh` must not be group/world writable. The script never reads a client
private key or emits key contents. It does not enable Remote Login or modify sshd.

`restrict` disables forwarding, PTY and user RC; it **does not restrict remote shell
commands or filesystem access**. The identity has the target account's permissions.
Only run commands authorized by the owner. It is not an interactive terminal,
forwarding identity, root escalation grant, or permission to use other credentials.

## Active source wiring

The actual SealedSecret JSON is referenced by the owning Kustomization. Only the
existing `initial-config` init container receives the optional raw Secret at
`/var/run/legacy-mac-ssh-input`, read-only with root-owned **0400** files. Its
existing pinned image Python directly runs `/etc/legacy-mac-ssh/prepare_mac_ssh.py`
as root; that script resolves `pwd.getpwnam('hermes')` and refuses a root target.
The actual invoking terminal was observed as UID/GID **10000/10000**, not root.

Preparation copies only `ssh-privatekey`, `ssh-publickey`, and `known_hosts` into
a new `emptyDir` with `medium: Memory`, `sizeLimit: 1Mi`. The directory becomes
hermes-owned **0700**, each file hermes-owned **0400**. Exclusive staging, flush,
fsync and atomic rename precede directory ownership transfer. A repeated init run
accepts only byte-identical prepared credentials and can finish approved partial
0600 staging without truncation or deletion. A fresh empty Kubernetes tmpfs's
initial 0777/01777 mode is tightened before any payload is written. Read-only
Secret volume roots may retain writable mode bits; actual files and internal
projection directories still require trusted ownership and strict permissions.
Unexpected or different destination state, unsafe source metadata, partial input
and escaping symlinks fail closed with sanitized errors. Normal Kubernetes
`..data` symlinks resolving
inside the trusted Secret volume are supported. An empty optional input, including
AtomicWriter's empty timestamp directory and `..data` symlink, logs
`Mac SSH unavailable: optional Secret absent.` and succeeds without copying, so
unrelated gateway startup continues; the SSH wrapper still refuses missing keys.

The main `hermes` container receives only the RAM credential mount read-only at
`/var/run/legacy-mac-ssh`, never the raw Secret. Nothing is copied to `/opt/data`
or any PVC. There is no environment export, new container, image change, resource
request/limit change or NATS change. Do not weaken permissions to 0444. Profiles
are not a security boundary: same-UID processes and root can read the credentials.
No other profile's configuration is modified.

A generated, hash-named ConfigMap mounts only Python source read-only at
`/etc/legacy-mac-ssh` in main and initial-config. Durable terminal invocation:

```sh
python3 /etc/legacy-mac-ssh/mac_ssh.py --user zeus -- 'id -un; uname -s'
```

Replace `zeus` only with the locally enrolled account. The wrapper fixes the egress
hostname to `omos-mbp.hermes.svc.cluster.local`, ignores SSH config (`-F /dev/null`),
uses only the mounted identity, pins ED25519 with strict host checking and the
mounted `known_hosts`, disables password/interactive auth, agents, forwarding,
PTY, local commands, proxies and connection multiplexing. It never copies or
prints private-key contents. `--credentials` is available for offline tests, not
an instruction to create persistent plaintext copies. For Homebrew commands,
explicitly export `/opt/homebrew/bin:/usr/local/bin:$PATH` in the remote command.

## Verification and recovery

```sh
python -m unittest discover -s platform/hermes/onemind-hermes/mac-access -v
git diff --check
make validate
```

Tests exercise real `ssh -G` option parsing and real `ssh-keygen` fingerprinting,
wrong-fingerprint rejection, repeat enrollment and preserved existing bytes,
unsafe permissions/symlinks, missing mount with sanitized stderr, RAM preparation
with actual nonroot reads of 0400 synthetic files and real `ssh -G`, projected
symlinks, optional absence, synthetic UID/account dispatch, fail-closed metadata,
and source wiring, encrypted fields, image and resource invariants. Root ownership
transfer is not exercised by this nonroot test runner. They do not prove Mac enrollment,
network authentication or controller decryption.

After owner merge, read the exact Argo revision and SealedSecret decryption status,
verify UID and 0400 without printing contents, then run the read-only invocation
above. If init saw an empty optional Secret, preparation is not live-updated:
after decryption, a separately authorized pod recreation is needed for a fresh RAM
volume and another init run. Do not try to overwrite an existing prepared volume.
Do not call access working until the remote account and Darwin result are observed.
Key rotation requires a fresh dedicated identity and owner-local append; never
reuse broad keys. Revocation/removal of an authorized key or Kubernetes resources
requires separate owner authorization; this bootstrap deletes nothing.
