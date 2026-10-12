# vault — Obsidian task engine sync bake-off (staged, not deployed)

Status: **staged on 2026-10-12, paused.** No Argo CD Application references this
directory, so merging it deploys nothing. Full research, evidence and test plan:
[`docs/research/2026-10-12-obsidian-task-engine-sync-bakeoff.md`](../../docs/research/2026-10-12-obsidian-task-engine-sync-bakeoff.md).

| File | Purpose |
|---|---|
| `namespace.yaml` | `vault` namespace, Pod Security `restricted` |
| `syncthing.yaml` | Candidate A: Syncthing 2.1.6 peer + `vault-git` sidecar (commit every 60 s, optional push); tailnet Service `vault-sync:22000` |
| `couchdb.yaml` | Candidate B: CouchDB 3.5.2 for Self-hosted LiveSync; tailnet HTTPS ingress `vault-couchdb` |
| `network.yaml` | Default deny; ingress only from the Tailscale operator; DNS + HTTPS egress for git push |
| `kustomization.yaml` | Bundles the above into namespace `vault` |

Images are pinned by multi-arch digest (verified on Docker Hub 2026-10-12):
`syncthing/syncthing:2.1.6`, `couchdb:3.5.2`, `alpine/git:v2.54.0`. Both PVCs use
`do-block-storage-retain` so deleting the StatefulSet never deletes vault data.

## Before enabling

1. Decide the git remote (Forgejo in-cluster, private GitHub, or both) — open question in the research doc.
2. Create secrets out of band (never commit plaintext; follow [`docs/secrets-eso.md`](../../docs/secrets-eso.md) or SealedSecrets as used elsewhere):
   - `vault-couchdb` with keys `username`, `password` (only if candidate B is tested).
   - `vault-git-remote` with key `url` (optional; without it the sidecar commits locally only).
3. Add `platform/argocd/platform-vault.yaml` (path `platform/vault`, namespace `vault`), modelled on `platform-n8n.yaml`.
4. Add the `data-vault-sync-0` PVC to the VolSync inventory in `platform/backups/volsync/` so snapshots go to DO Spaces.

## Candidate A: Syncthing setup

1. GUI (kept on pod localhost): `kubectl -n vault port-forward statefulset/vault-sync 8384:8384`, open `http://127.0.0.1:8384`, set a GUI user/password.
2. Add folder path `/var/syncthing/vault`; enable **Staggered File Versioning** on this peer.
3. Devices:
   - macOS: Syncthing ([syncthing.net](https://syncthing.net/downloads/)).
   - Android: Syncthing-Fork from [F-Droid](https://f-droid.org/packages/com.github.catfriend1.syncthingfork) or [GitHub releases](https://github.com/researchxxl/syncthing-android/releases) (the official Android app is [discontinued](https://github.com/syncthing/syncthing-android)). Exempt it from battery optimization.
   - On each device, add the cluster peer with address `tcp://vault-sync:22000` (tailnet MagicDNS).
4. Open the synced folder as a vault in Obsidian; install **TaskNotes**; optionally **Syncthing Integration**.
5. Check commits: `kubectl -n vault exec statefulset/vault-sync -c vault-git -- sh -c 'GIT_DIR=/var/syncthing/git/vault.git git log --oneline | head'`.

The sidecar's git dir lives outside the synced folder, so `.git` is never replicated to
phones. Excluded from commits: `.stfolder`, `.stversions/`, `.stignore`, `.obsidian/workspace*.json`.
`.sync-conflict-*` files **are** committed on purpose — they are data until resolved.

## Candidate B: LiveSync setup

1. After CouchDB is Ready, apply LiveSync's required CouchDB settings with the upstream tool
   ([setup_own_server.md](https://github.com/vrtmrz/obsidian-livesync/blob/main/docs/setup_own_server.md),
   [couchdb-init.sh](https://github.com/vrtmrz/obsidian-livesync/blob/main/utils/couchdb/couchdb-init.sh), requires Deno):
   `hostname=https://vault-couchdb.<tailnet>.ts.net username=… password=… database=obsidiannotes ./couchdb-init.sh`
   (download and read the script first; do not pipe it from the network).
2. Configure the plugin per [quick_setup.md](https://github.com/vrtmrz/obsidian-livesync/blob/main/docs/quick_setup.md); enable E2E encryption.
3. Back up the test vault before installing (upstream README instruction).

## Not included yet

vault-engine, Forgejo, ntfy, the Argo Application and VolSync entries — see the resume
checklist in the research doc.
