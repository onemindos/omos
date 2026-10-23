# Rebuild the platform from zero (disaster recovery runbook)

Use this when the cluster is gone (provider outage, account problem, deleted
cluster) or when you are moving to a new cluster (home, other cloud). It rebuilds
everything from **Git + backups + a small set of keys you keep outside the cluster**.

Honest status first: this procedure is assembled from the documented backup tooling and
the individual restore tests already on record. **It has never been run end to end on
an empty cluster.** The first full run is the "game day" at the bottom. Until it has
passed once, treat the time estimates as guesses and read "What is NOT covered" before
trusting it.

New to Kubernetes? Three ideas are enough:
- **Git** holds the *desired state* (what should run). Argo CD makes the cluster match Git.
- **Backups** hold the *data* (databases, files, message streams). Git has no data.
- **Keys** decrypt things. Without the keys, Git secrets and backups are unreadable noise.
  This is why keys must live somewhere that survives losing the cluster AND the laptop.

## 0. Prerequisites: what you must hold OUTSIDE the cluster

Do not start without all of these. Check them every quarter (game day).

| # | Item | Why | Where it should be |
|---|---|---|---|
| K1 | **Sealed Secrets private keys** (2 active keys exist live, in `kube-system`) | Every `*.sealed*.yaml` in Git (about 150 SealedSecrets live) is unreadable without them | Inside the encrypted cluster-recovery export (`backups/cluster-recovery/`), AND decrypted-and-stored offline in your password manager/offline medium |
| K2 | `BACKUP_ENCRYPTION_KEY` (base64, 32 bytes) | Decrypts every `archive.tar.gz.aesgcm` (TAK Postgres, NATS, Authentik, cluster-recovery, ...) | Offline: `~/.config/onemind/recovery` plus a second physical copy |
| K3 | **Restic password** | Decrypts the 46 VolSync volume repositories | `~/.config/onemind/recovery/restic-password` plus an independent offline copy |
| K4 | **S3 access key and endpoint** for the bucket(s) holding backups | To download anything | Password manager. The in-cluster copy is useless in a disaster. If you only have the sealed copy you are locked out |
| K5 | GitHub access to `onemindos/omos` (private repo) | Argo CD reads Git | Deploy key or token in password manager. Also keep the encrypted repo archive (`repository-recovery.json` bundle, see docs/operations/home-migration.md) |
| K6 | Tailscale OAuth client, Cloudflare API token and tunnel credentials, GHCR pull token | Networking and image pulls | Password manager (also inside the cluster export) |
| K7 | NATS `values.secret.yaml` contents (ONEMIND/SYS passwords) | Helm-CLI NATS install is not fully in Git (gap G3) | Password manager; also inside the cluster export (Helm release Secrets) |

Tools on your computer: `kubectl`, `helm`, `git`, `python3`, `restic`, `kubeseal`,
`jq`, `nats` CLI, `argocd` CLI (optional). The decrypt/export tools are in
`platform/backups/tools/` (`decrypt.py`, `export_helm.py`, `mirror.py`).

Where backups live: bucket `onemind-backups` (DO Spaces nyc3). A second copy exists
ONLY once `backup-offsite-copy` is enabled (live status 2026-10-10: suspended, never run;
see gap G1). A laptop copy exists at `~/OneMindRecovery/<date>` (54 GiB, made by
`tools/mirror.py`).

## 1. Order of operations

Rough rule: foundations first, then identity/secrets, then data stores, then apps.
Do not skip ahead; later steps assume earlier ones. Never point a rebuilt cluster at
production NATS or production DNS until step 12.

### Step 1. Create the empty cluster (about 1-3 h)

- Home or other provider, per HOME-CLUSTER-PLAN. Result: `kubectl get nodes` shows Ready nodes.
- Create the StorageClasses `do-block-storage` (default) and `do-block-storage-retain`
  on the new storage provisioner (see DO-LOCK-IN section 2). Why: manifests reference
  these names; keeping the names means no edits for the rebuild.
- Label nodes `onemind.io/compute-role=ecosystem` (DO-LOCK-IN section 7).
- Install a VolumeSnapshot-capable CSI (or plan VolSync `Direct` copy; see step 9).
- Install CRDs the Git tree expects only via Argo later (do not hand-apply).

### Step 2. Retrieve and decrypt the recovery material (about 30 min)

On a trusted machine (NOT in the cluster):
1. Download the newest `backups/cluster-recovery/<timestamp>/archive.tar.gz.aesgcm`
   (and its `.verified.json`) from the bucket (or the offsite mirror / laptop copy).
2. `python3 platform/backups/tools/decrypt.py <archive> --output <private path>` using
   K2. Output is mode 0600 and contains PLAINTEXT secrets. Keep it on an encrypted disk,
   never in Git or chat.
3. The decrypted `recovery.json` (`OneMindClusterRecoveryV2`) is a list of Secrets and
   ConfigMaps (336 Secrets, 213 ConfigMaps, 2 sealing keys at last check).
4. `python3 platform/backups/tools/export_helm.py <archive> --output <private dir>`
   extracts the chart version, values and rendered manifests of all 11 Helm-CLI releases
   (gap G3). Review `inventory.json`.

### Step 3. Restore the Sealed Secrets keys BEFORE anything else (about 15 min)

Why: the Sealed Secrets controller generates a NEW key if it starts empty, and then every
SealedSecret in Git fails to decrypt ("no key could decrypt secret").
1. Extract only the Secrets labelled `sealedsecrets.bitnami.com/sealed-secrets-key=active`
   (type `kubernetes.io/tls`, `kube-system`) from `recovery.json` into manifests (mode 0600).
   There is **no script for this yet** (gap G2); a 15-line `jq` filter on
   `.items[] | select(.kind=="Secret" and .metadata.labels["sealedsecrets.bitnami.com/sealed-secrets-key"]=="active")`
   is enough. Do not print the values.
2. `kubectl create namespace kube-system` exists already; `kubectl apply -f` the key files.
3. Install the controller (`platform/core/controller.yaml`, image
   `bitnami/sealed-secrets-controller:0.38.4` at last check) and confirm it logs that it
   loaded existing keys. If it already started and generated a fresh key, apply the old
   keys and restart it (`kubectl -n kube-system rollout restart deploy/sealed-secrets-controller`).
4. Test: apply ONE small SealedSecret from Git and confirm the Secret appears.
5. Delete the plaintext key files.

### Step 4. Bootstrap Argo CD (about 30 min)

Argo CD itself is a Helm-CLI release, not self-managed (gap G3). Chart `argo-cd` 10.1.2
(app v3.4.4); values in `platform/values/argocd.yaml` (note `platform/install.sh` says
"defaults": confirm against the `export_helm.py` output).
1. `helm upgrade --install argocd argo/argo-cd -n argocd --create-namespace --version 10.1.2 -f platform/values/argocd.yaml`.
2. Apply the repo credential Secret `argocd-repo-omos-infra` (and the three
   `ecosystem-helm-*` OCI repo Secrets) from the decrypted export, or recreate with K5.
3. Do **not** apply the root app yet if you are restoring to different storage/networking.
   First prepare the overlay (storage class names, LB annotations, remove cloud-only
   Applications). Current `main` describes DigitalOcean (see docs/operations/home-migration.md).
4. When ready: `kubectl apply -f platform/argocd/app-of-apps.yaml`
   (repo `onemindos/omos`, path `platform/argocd`, branch `main`).
5. Expect `platform-core` to sync first. The rest will start but data-dependent apps will
   sit empty or crash-loop until their data is restored (steps 6-10). That is fine.
   To avoid surprises, you can set `spec.syncPolicy.automated` off on data apps while restoring,
   or sync in the order below.

### Step 5. platform-core, cert-manager, ingress, Tailscale (about 1 h)

- `platform-core`: Sealed Secrets, Let's Encrypt DNS issuers, `cloudflared` (0 replicas),
  GitHub runner, Tailscale ACL policy.
- Helm-CLI releases still needed (gap G3, install from `export_helm.py` inventory or
  `platform/values/*.yaml`): `cert-manager` v1.20.3, `ingress-nginx` 4.15.1,
  `tailscale-operator` 1.98.4 (needs a NEW Tailscale OAuth client or the exported one; new
  proxy nodes will appear under new names, and old tailnet node entries must be removed
  by hand in the Tailscale admin console).
- Ingress at home: drop DO annotations; use Cloudflare Tunnel (`platform-cloudflared`,
  staged) for web hosts.
- Verify: `kubectl get certificate -A`, `tailscale status` shows operator proxies.

### Step 6. Storage and operators (about 1 h)

Argo installs these from Git: `platform-volsync` (0.16.0), `ecosystem-velero`,
`ecosystem-cnpg`, `ecosystem-barman`, `platform-clickhouse-operator`, `platform-backups`
(CronJobs), `platform-volume-backups` (ReplicationSources). Point `ecosystem-velero` and
CNPG's object store at a bucket you can read (K4); change endpoint/region only.
Pause ReplicationSources for now (step 9): a new empty volume must not back up OVER a
good repository. Set `spec.trigger` aside or scale the volsync operator to 0 until each
PVC is restored. This is a data-loss risk; see "Danger list" below.

### Step 7. NATS + JetStream (about 1-2 h)

1. NATS is Helm release `nats` 2.14.2 with out-of-band `values.secret.yaml` (K7, gap G3).
   Install it from the `export_helm.py` values, 3 replicas, same cluster name and accounts.
2. Install `nack` 0.34.0 (needs `nack-values.secret.yaml`, see `platform/nats/nack-values.template.yaml`);
   it recreates the 39 Stream/KeyValue objects from `platform/nats/streams/` (empty).
3. Restore stream DATA: decrypt the newest `backups/nats-onemind/<ts>/archive.tar.gz.aesgcm`
   and run, per stream, `nats backup restore <dir>/<stream>` against the new server
   (this is what the backup job's restore test does, into a throwaway server). The
   stream must NOT already exist with data; if nack created it empty, delete the empty
   stream first (or restore before applying nack's CRs). 39 streams; replicas=3 in prod.
4. Verify: `nats stream ls`, message counts vs `.verified.json`, consumers present.
5. If a home cluster runs next to the cloud one, do NOT peer it into the production
   NATS until step 12 (NATS leaf/authority rules, HOME-CLUSTER-PLAN).

### Step 8. Databases (about 2-4 h)

Order: Postgres first, then the apps that use them.

- **TAK Postgres** (PostGIS): install `takserver` (Helm release rev 34, chart
  `platform/takserver/chart`, values `platform/takserver/values-onemind.yaml`) with the database empty,
  then restore the newest `backups/takserver-postgres/<ts>/archive.tar.gz.aesgcm`
  (decrypt, `pg_restore --no-owner` into the TAK database with matching major version
  and PostGIS). Last backup took 86 min and the restore check passed (64 tables).
  Keep the quota override (18 PVCs / 100 GiB). TAK certs/CoreConfig come from the
  restored Secrets/ConfigMaps (step 2). Start the TAK Postgres PVC from a fresh volume
  only; never overwrite an existing one.
- **Authentik** (`platform-auth-infra`, `platform-auth-helm`): restore the
  `backups/postgres-auth/...` dump into Authentik Postgres 17 before starting the server;
  signing keys and issuer identity are in that DB and in Secrets (step 2) so existing
  OIDC clients keep working. Check login after restore.
- **n8n, matrix, takml, backstage**: same pattern with their logical
  archives (`postgres-<name>`). Backstage needs ALL plugin databases (owned-database archive).
- **CloudNativePG `onemind-data`** (3 instances, 120 GiB: core, openfga, litellm DBs):
  Git contains only `bootstrap.initdb` (no recovery stanza; gap G4). Create a restore
  overlay with `bootstrap.recovery` + `externalClusters` pointing at the Barman object
  store (`cnpg/onemind-data/`), as proven in `docs/operations/cnpg-restore-2026-10-09.json`
  (7 databases, 114 tables matched). Apply the overlay in a drill; in a real DR, apply it
  instead of the plain cluster.
- **ClickHouse**: restored by VolSync (step 9); the native backup job is retired.

### Step 9. Volume restores with VolSync (about 2-6 h; depends on bandwidth)

46 PVCs have restic repos (42 non-empty; 4 confirmed empty). Per PVC:
1. Create the destination PVC (same name/size/class as in Git; Argo may already have made it).
2. Apply a `ReplicationDestination` with `copyMethod: Direct`, `repository: <secret name>`,
   `trigger.manual: restore-1`, `moverSecurityContext` per the examples in
   `platform/backups/volsync/restore-tests/*.yaml` (Authentik, Hermes, ClickHouse, CNPG
   tested). Use `Direct` on a non-DO cluster; `Snapshot` needs a snapshot class.
3. Wait for `status.lastManualSync == restore-1`, then start the app that uses the PVC.
4. Check ownership and data-directory modes before starting any database.
5. Only after the app is verified, re-enable that PVC's ReplicationSource.
Raw copies of database or NATS volumes are crash-consistent only; prefer logical archives
for Postgres/NATS (steps 7-8) and use VolSync for file data (Hermes state, Harbor, OpenSearch,
Grafana/Prometheus, Valhalla/geo data, Matrix media, etc.).
Per-PVC order of priority: Hermes (25 GiB), Authentik, TAK volumes, Matrix, n8n, Backstage,
Harbor (61 GiB), OpenSearch, Dapr, ClickHouse (20 GiB), monitoring (Prometheus/Grafana/
Loki/Tempo), onetak geo data (about 450 GiB; may be rebuilt from public sources instead).

### Step 10. Velero and the remaining apps (about 1 h)

- Velero restores Kubernetes OBJECTS only (volume snapshots are off). Use it for namespaces
  `onemind-core`, `onemind-data`, `onemind-ai` if Git is missing something:
  new BSL, `accessMode: ReadOnly`, then `velero restore create --from-backup <name>`.
  Exclude Secrets that Sealed Secrets owns to avoid conflicts.
- Let Argo finish syncing every remaining Application. Resolve the known OutOfSync items
  (`SealedSecret auth/ghcr` owned twice; `litellm-migrations` Job).
- Apps with no Argo owner (gap G5): matrix (`platform/matrix/`, deploy script), NATS UIs
  (`platform/nats-tools/`), openclaw. Apply by hand from Git, using the restored data.

### Step 11. Networking cutover prep (about 1 h)

Cloudflare Tunnel, DNS names, Tailscale nodes, TAK public endpoint (DO-LOCK-IN section 3).
Do this on a test hostname first.

### Step 12. Go live (about 1 h plus a monitored window)

Only after the verification checklist passes: move DNS/tunnel/TAK endpoint, point EUDs,
enable ReplicationSources and backup CronJobs on the new cluster, and THEN retire the old
one. If the old cluster is alive, scale its duplicate workers (feeds, projectors, schedulers)
to 0 first so two clusters do not both consume or write the same streams.

## 2. Verification checklist

Tick every line; record date, who, and times.

- [ ] `kubectl get nodes`: all Ready; storage classes present; default class set.
- [ ] Sealed Secrets: 2 old keys loaded; a SealedSecret from Git decrypts.
- [ ] `kubectl -n argocd get applications`: all Synced/Healthy (compare with 47 expected,
      allowing the known OutOfSync three).
- [ ] No Pending PVCs; PV count and sizes match the inventory (95 PV / 1,182 GiB on DO; expect fewer if you skip the idle geo volumes).
- [ ] NATS: 39 streams, message counts close to the backup marker, consumers attached.
- [ ] TAK Postgres: 64 tables; TAK API pod Ready; a test EUD connects over Tailscale and a CoT message round-trips.
- [ ] Authentik: admin login works, an OIDC client (Argo, Grafana, Backstage) logs in.
- [ ] CNPG: Cluster Ready, 7 databases, 114 tables.
- [ ] Hermes: starts, `state.db` present, answers a prompt.
- [ ] Harbor/OpenSearch/ClickHouse: UI/query works (ClickHouse: 9 tables, tens of millions of rows).
- [ ] Monitoring: Prometheus targets up; Alertmanager has a real receiver (gap G6).
- [ ] Backups on the new cluster: run one job of each kind (`kubectl create job --from=cronjob/...`), confirm `.verified.json`, confirm VolSync sources `Successful`.
- [ ] Offsite copy ran at least once and `_offsite/last-success.json` is fresh.
- [ ] Record measured recovery time and recovered timestamp (data loss window = time of backup).

## 3. Time estimate

Estimates for one experienced operator with all keys at hand, about 1 Gbit/s link, data
about 100 GiB compressed. Add 2x for a first run.

| Phase | Estimate |
|---|---|
| Steps 1-5 (cluster, keys, Argo, core) | 3-5 h |
| Steps 6-8 (operators, NATS, databases) | 4-8 h |
| Step 9 (46 volumes) | 2-6 h |
| Steps 10-12 + verification | 3-4 h |
| **Total** | **about 1.5-2.5 working days; budget 3 days for the first run** |

Recovery point: up to 24 h of data loss (daily backups). Streams written between the
03:30 UTC NATS backup and the failure are lost; TAK Postgres up to its 02:00 UTC backup.
CNPG has WAL archiving, so it can recover to a point in time.

## 4. Danger list (read before touching anything)

1. Never let an empty restored volume back up over a good restic repo. Keep VolSync
   sources paused until each restore is verified.
2. Never restore into existing production PVCs or point a test cluster at production NATS.
3. Never replay ALL exported Secrets/ConfigMaps. Skip service-account tokens, Helm release
   records and provider-specific secrets (DO CSI token, LB IDs).
4. Do not run `terraform apply` from `infra/` (state does not match the live cluster).
5. Nothing is deleted from DigitalOcean until the new cluster passes the checklist.
6. Do not delete the `-retain` volumes or VolumeSnapshots while they are the only recent copy.

## 5. What is NOT yet covered (honest list)

| ID | Gap | Impact | Fix |
|---|---|---|---|
| G1 | **Every backup is in one DO Spaces bucket.** The offsite copy exists in Git but is `suspend: true` live and has never run; its sealed credentials are not wired in. The laptop mirror is the only independent copy and the laptop is a single point of failure | Losing the DO account loses all backups | Finish OFFSITE-COPY.md steps 1-6 to an OFF-DO provider (B2/R2); keep a second physical copy of the laptop mirror and keys |
| G2 | **Sealing-key restore is manual and untested.** The export exists (2 keys, decrypt-equality verified) but nothing has ever replayed it onto a new cluster, and no script extracts and applies the keys | Rebuild blocks at step 3 | Write a small extraction script in `platform/backups/tools/`, drill it |
| G3 | **11 Helm-CLI releases are not in Git as Applications** (argocd, cert-manager, ingress-nginx, loki, monitoring, nack, nats, otel-collector, tailscale-operator, takserver, tempo). NATS needs a gitignored `values.secret.yaml`; whether Git values match deployed values was never diffed | Most important workloads (TAK, NATS) rebuild only from `export_helm.py` output, not from Git | Adopt releases into Argo (review finding H2), TAK and NATS first |
| G4 | **CNPG recovery stanza is not in Git.** Restore is proven (2026-10-09) but via a temporary cluster | DR takes manual YAML authoring | Commit a `restore/` overlay |
| G5 | **Workloads with no Argo owner**: matrix (Synapse/Element/Postgres/Redis), NATS UIs, openclaw | Rebuilt by hand | Create Applications |
| G6 | **Alertmanager delivers nowhere** (review A17) | You would not know a backup stopped | Wire a receiver |
| G7 | **Secrets not in SealedSecret form**: NATS passwords (out-of-band), Helm values with credentials (monitoring has a plaintext `--password=` arg, review/STACK.md), Tailscale OAuth, Cloudflare tokens; the 1Password SA token and other tokens sit in plaintext in a Dropbox-synced vault file (review C1) | A rebuild depends on a password manager that must itself be safe and rotated | Rotate; move to sealed/ESO |
| G8 | **Hermes local state**: the in-cluster agent home (25 GiB PVC, `state.db`) is covered by VolSync; any Hermes configuration kept on the laptop (local config, the `hermes()` shell function, local skills; location not verified here) is not backed up by this system | Laptop loss loses local agent config | Inventory those files and back them up with the same restic scheme |
| G9 | **Mutable `:latest` images (31) and no air-gapped registry**: a rebuild may pull different images than today, and GHCR/Docker Hub must be reachable. Custom images are amd64-only | Behaviour drift; ARM hardware will fail | Pin digests; mirror images into Harbor; choose amd64 hardware |
| G10 | **Tailscale state**: proxy node keys are regenerated; tailnet hostnames and ACL references change. EUD profiles hardcode `152.42.155.95` | EUDs need new profiles | Use DNS names |
| G11 | **etcd/Kubernetes object state** is not backed up on DOKS (managed), and at home it will be our job. Velero covers 3 namespaces only | Cluster-scoped objects come from Git only | Add etcd snapshots at home; widen Velero |
| G12 | **Not yet drilled end to end**; ClickHouse rows, TAK app workflows, Authentik login after restore are checked only at database level | Unknown unknowns | Run the game day |
| G13 | **Droplet (node) backups are off** and CSI snapshots are DO-only | Irrelevant for a rebuild, but do not count them as DR | none |
| G14 | **Large geo volumes** (Valhalla 50 GiB at 99%, Nominatim/Overpass 2x100 GiB idle) | Restore time and cost | Decide rebuild-from-source vs restore |

## 6. Quarterly game day

Purpose: prove the runbook works before you need it. Do it quarterly and after any major
change (new database, new Helm release, key rotation). First one: before any cloud
shutdown step in HOME-CLUSTER-PLAN. Duration: one day (2 people is better than 1).

Rules: production is never touched and never given write access. Use a throwaway cluster
(a spare PC with k3s, a kind/k3d cluster for the logic parts, or a cheap short-lived
cloud cluster). Use a separate Tailscale tailnet tag or no tailnet. Disable cloudflared
and anything with credentials to live systems (feeds, bridges, automations). Do not give the
drill cluster the production NATS URL, GHCR push rights or TAK EUD certs for live units.

Procedure:
1. Prepare (before the day): confirm K1-K7 are present and readable from the OFFLINE
   copies, not from the cluster. Run the read-only freshness checks: latest object of each
   `backups/*` prefix is under 30 h old; `.verified.json` markers exist; offsite
   `_offsite/last-success.json` is under 36 h old; VolSync sources `Successful`.
2. Start the clock. Create an empty cluster (step 1).
3. Execute Steps 2-10 as written, in a scratch namespace set where possible
   (e.g. suffix `-drill`). Write down every point where you had to improvise.
4. Run the verification checklist. Data comparisons: row counts for TAK, CNPG,
   Authentik; NATS stream message counts; file counts on a restored PVC compared to a
   VolSync status.
5. Stop the clock. Record: elapsed time, recovered backup timestamps, failures, fixes.
6. Fix the runbook (this file) the same week; file Git issues for every gap discovered.
7. Tear down the drill cluster (this is your own throwaway cluster; delete nothing on production).
8. Rotate anything you had to expose during the drill.

Pass criteria: every checklist item ticked, total time under 2 working days, no step
needed a secret that was only inside the cluster, and the runbook changes are merged.
Once per year also test a full cutover rehearsal (HOME-CLUSTER-PLAN phase exit criteria).
