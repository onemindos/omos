# Hermes runbook

Operating and recovery procedures for Legacy (the `hermes-legacy` pod) and its edge nodes. Distilled from the July 2026 runbook, integration and fabric notes, and the operational lines in Legacy's own memory. Anything that is only an address, a credential or a dated status has been dropped.

**How far to trust each part.** Each section says whether it was **verified** against the live pod on 2026-10-16 or is carried over from July and **not re-verified**. Fix this file when you find something wrong.

For stack facts (NATS accounts and ports, TAK, ClickHouse) see `FABRIC-REFERENCE.md`. For the agent list see `AGENTS.yaml`. For the recovery design see `docs/operations/hermes-recoverability-plan.md`.

## 1. What runs, and how to check it (verified)

- The pod `hermes-legacy` (namespace `hermes`, image `nousresearch/hermes-agent`, pinned tag in `deployment.yaml`) runs the container under s6: the gateway (`main-hermes`) and the `dashboard`. A `nats-leaf` sidecar sits beside it.
- State is the PVC `hermes-data`, mounted at `/opt/data`. A pod restart does not lose it.
- Tailscale is **not** run inside the pod any more. Access to other machines goes through Tailscale egress services (section 2).
- The pod has **no Kubernetes RBAC** (`kubectl auth can-i get pods` answers no). Cluster commands are run from an admin workstation, not from the pod.

From an admin workstation:

```bash
kubectl -n hermes get pods
kubectl -n hermes logs deploy/hermes-legacy -c hermes --tail=50
kubectl -n hermes rollout restart deployment hermes-legacy     # restart; the volume persists
```

Check the agent is on the bus (from inside the pod, or any host with NATS access):

```bash
nats micro ls --server nats://localhost:4222     # expect an "agents" service for legacy
nats pub legacy.health.alive "$(date -u +%s)"    # heartbeat (July practice)
```

## 2. Reaching other machines (verified that the services exist)

The `hermes` namespace has `ExternalName` services that point at Tailscale egress proxies: `omos-mbp` (Zeus's Mac), `omos-dell-edge`, `omos-dell-zeus`, `omos-ha-edge` (Home Assistant), `omos-intel`, `omos-s24`. From the pod they resolve as `<name>.hermes.svc.cluster.local`.

```bash
ssh -i <pod ssh key> zeus@omos-mbp.hermes.svc.cluster.local "<command>"
```

The pod's SSH key lives under `/opt/data/home/.ssh/`. Non-interactive SSH to the Mac does not load the shell profile, so prefix commands that need Homebrew tools with `export PATH=$PATH:/opt/homebrew/bin:/usr/local/bin`.

Not re-verified since July: SSH to Home Assistant as the `hassio` user (background processes need `setsid cmd >/tmp/out.log 2>&1 </dev/null &`), and ADB to the phone through a tunnel on the Mac.

## 3. NATS leaf connection

**The pod (verified).** The sidecar runs a leaf node from the secret `hermes-leaf-conf` (CA from `hermes-nats-ca`). It listens on `localhost:4222`, has one local account, and connects to the cluster's leaf port over TLS using the in-cluster service address. Hermes always connects to `nats://localhost:4222`; never to the cluster directly.

```bash
nats account info        # a connected leaf shows the account; "no responders" means the sidecar is down: restart the pod
```

**Edge nodes outside the cluster (carried over from July, not re-verified).** Common failures:

| Symptom | Cause | Fix |
|---------|-------|-----|
| TLS handshake failure | Pointing at a node address | Use the `omos-nats` Tailscale load balancer address on the leaf port, never a node address |
| `cannot find local account` | An `account:` field in the leaf remote when no local account is defined | Remove `account:`, or define the local account |
| TLS verification errors with a self-signed hub | `tls { insecure: true }` in the wrong place | It goes inside the `remote` block |

## 4. Routing work between Legacy and her edge nodes

Legacy runs in the cloud and is always on; an edge node has local hands. Route by capability:

| Task | Route to | How |
|------|----------|-----|
| Anything needing macOS: Apple Notes, Reminders, Find My, iMessage, local files, camera, USB or serial devices, Homebrew | Edge node (Mac) | NATS request or SSH |
| Direct cluster access with kubectl | Edge node or an admin workstation | The pod has no RBAC |
| Git push over SSH remotes | Edge node (it holds the SSH keys) | SSH |
| Cron scheduling, always-on presence, authoritative decisions, GitHub API, NATS operations, cross-agent orchestration | Legacy | Direct |
| Code implementation across repos | Delegate (Forge once it runs; until then `delegate_task`) | NATS request |
| Home Assistant control | Legacy | Through the Home Assistant NATS bridge (`cmd.ha.*`, `events.ha.*`), never direct HTTP from the pod |

Prompt an edge node: `nats req agents.prompt.hermes.zeus.legacy-edge "<task>"`.

## 5. GitHub and other credentials

- The pod is logged in to GitHub with `gh`. That token is far more powerful than it needs to be and is being replaced with a narrow one (tracked as a follow-up). Until then, use it only for the repos Hermes actually works on.
- **Never put a token in a git remote URL**, in this repo, or in Legacy's memory. Credentials come from the pod's environment or `gh auth`.
- Notion is accessed with its REST API using a token from the environment. The Notion MCP and Composio were not reliable for it in July.
- Home Assistant's token is a file on the pod volume. It is never committed.

## 6. Standing rules from Zeus

These were operating rules kept in Legacy's memory. They belong with her behaviour, so they are recorded here as well.

1. **Never destroy data.** Deduplication means byte-verify and quarantine, never blind delete.
2. **Never delete files or directories without Zeus's explicit verbal go-ahead**, even after verification. The one standing exception is Kubernetes namespace deletion for cluster cleanup, which he has authorised.
3. **Patch skills immediately** when a loophole surfaces. Do not wait to be asked.
4. **Notion:** a PATCH on a select or multi-select property replaces the whole options array, so send the full set plus the new option and GET the schema afterwards.
5. **Never squash merges** when absorbing or archiving a repository; keep history.

## 7. Where things live

- TAK voice and video services are in the `takserver` namespace (MediaMTX `:8554`, Murmur `:64738`). TAK access is by client certificate. CloudTAK was retired; the TAK portal replaced it.
- Cloudflare for agents is the official MCP (`mcp.cloudflare.com/mcp`), authorised with OAuth or a scoped API token.
- The Mac's synced Dropbox folder is a working copy; the OneMind folders there are the primary one-way backup mirror. The `mac-storage-ops` skill covers the operations.
- Agent memory is not in this repo (it holds personal and business detail). The authoritative backup is the nightly VolSync snapshot; a readable copy is kept in the vault.

## 8. Recovery

| Situation | What to do |
|-----------|------------|
| Pod crashed or restarted | Nothing: `hermes-data` persists. Check section 1; if the leaf sidecar is down, restart the pod |
| The volume is lost or corrupted | Restore the latest VolSync snapshot of `hermes-data` into a new PVC, verify it, then point the Deployment at it. See `platform/backups/volsync/README.md`. Restores were verified twice: files on 2026-10-08 (319,496 files, matching checksum) and application level on 2026-10-16 (a copy of Hermes booted against the restored volume with the network cut: config loads, soul, cron and skills match live, `state.db` passes `quick_check`). See `docs/operations/hermes-restore-drill-2026-10-16.md` |
| The cluster is lost | Rebuild the cluster from `infra/`, let Argo CD sync `platform/`, then restore the volume as above |
| Rebuild identity from Git alone | **Not available yet.** The soul, skills and config are not exported to Git today. This is Step 2 of the recoverability plan |
| Re-authenticating after a restore | Log `gh` in again with a narrow token, confirm the SSH keys under `/opt/data/home/.ssh/`, and check the leaf connection (section 3) |

## 9. Other known pitfalls (carried over, not re-verified)

| Symptom | Cause and fix |
|---------|---------------|
| `fatal: could not read Username` on git push from the Mac | The Mac has no HTTPS git credentials: use SSH remotes, or push from the pod with `gh` |
| Home Assistant API answers "unauthorized" | The token expired or is wrong: check the token file on the pod volume, then test against the HA bridge, not the REST port |
| The pod cannot reach Home Assistant over HTTP | Expected. All Home Assistant traffic goes through the NATS bridge running on the Home Assistant host |
