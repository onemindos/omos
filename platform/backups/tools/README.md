# Independent encrypted recovery copy

`mirror.py` copies only completed Restic repositories and checksum-marked encrypted
archives produced since the selected date. It downloads credentials in memory from
`backstage/backstage-backup` in the explicit source context; no values are logged.
It never deletes remote or local data. Defaults cap the selected copy at 50 GiB and
preserve 40 GiB free disk space. Capacity failure requires another destination;
raise limits only after checking the target disk. Run it after backup jobs complete
so new repositories are included. A changing repository is not a consistent mirror.

```sh
python3 mirror.py --context do-nyc3-omos-cloud --output /path/to/private/recovery
restic --repo /path/to/private/recovery/backups/restic/auth/pgdata-authentik-postgres-0 \
  --password-file "$HOME/.config/onemind/recovery/restic-password" check --read-data
```

Verify each repository listed in mirror-manifest.json with the same command. The
Restic password is distinct from the logical archive key. Copy both keys independently
and retain the sealing-key recovery archive. Keep the mirror on a separate physical
medium for disaster recovery; a copy on the active laptop does not survive its loss.
Provider credential rotation must preserve existing encryption keys and archives.

`decrypt.py` authenticates logical archives OFFLINE and writes a new mode-0600 file.
It supports the current platform format, previous Backstage format and initial
cluster-secrets archive. It refuses to replace an existing output and deletes
partial cleartext on authentication failure. It does not extract or apply resources.

```sh
python3 decrypt.py /path/to/archive.tar.gz.aesgcm --output /private/path/archive.tar.gz
```

Current platform archives contain data/manifest.json plus PostgreSQL custom dumps,
NATS stream snapshots or data/recovery.json. Restore into an isolated target using
matching database versions/extensions. Select only needed recovery Secrets and
ConfigMaps; exclude old service-account tokens, provider-specific credentials and
Helm release records. Restore sealing-key Secrets into kube-system before expecting
existing ciphertext to decrypt. Verify identity, registry access and actual app
workflows before cutover. Native ClickHouse S3 backups use a separate format; these
archive tools do not decrypt or validate them.

VolSync reports an empty filesystem as successful without initializing a Restic
repository. The mirror records those volumes separately as confirmedEmptyVolumes;
it does not claim an encrypted archive or run Restic checks against a nonexistent
repository. Their Git-owned PVC declarations and initial retained snapshots remain.
Do not delete dormant application volumes simply because an export found no files.

`export_helm.py` recovers the deployed revisions' chart versions, values and rendered
manifests from Helm release Secrets inside the verified recovery archive. It needs
only the local archive and key, not DigitalOcean or a live Kubernetes API. Outputs
may contain plaintext app credentials: directories are mode 0700, files 0600, and
an existing output directory is never overwritten. Store these outside Git.

```sh
python3 export_helm.py /path/to/recovery/archive.tar.gz.aesgcm --output /private/helm-bootstrap
```

Review inventory.json and obtain matching charts before recreating manual releases.
Rendered manifests preserve the archived Helm revision, which can predate Argo
adoption and later live repairs. NATS's archived StatefulSet is superseded by the
current Git-owned workload. Do not blindly replay old manifests, provider storage
classes, exposure annotations or overlapping Argo resources. Review namespace,
storage/ingress overlays, selected secrets and chart inputs before a home drill.
