# Home cluster plan

Audience: the owner, who is new to Kubernetes. Aim: move from DigitalOcean
(list price about $828/month, paid by credits for now) to a home server, first as a
hybrid (cloud and home together), without losing data and with a way back at every step.

Companion documents: [DO-LOCK-IN.md](DO-LOCK-IN.md) (what must change),
[REBUILD-FROM-ZERO.md](REBUILD-FROM-ZERO.md) (how data and secrets get restored),
[../operations/home-migration.md](../operations/home-migration.md) (earlier outline),
draft PRs #142 (k3s, Rook, Forgejo, Tekton, Akri, Envoy; all dormant, target cluster
`onemind-home`) and #139 (GPU serving, OCM, Crossplane; dormant).

## 1. Recommendation in short

| Question | Recommendation | Why (beginner terms) |
|---|---|---|
| Kubernetes distribution | **k3s on a normal Linux install (Ubuntu Server LTS or Debian)** to start. Revisit **Talos** when you reach 3 nodes | You can SSH in, install Tailscale and tools like any Linux box, search any forum answer, and fix things by hand. k3s is one binary and is real, certified Kubernetes. Talos is excellent but has no SSH or shell at all, which is hard to learn on while also running a live system |
| Number of nodes at first | **1 node**, but design for 3 | One machine is cheap and simple. Be honest that it is a single point of failure: backups, not redundancy, protect you |
| Storage | **`local-path` (built into k3s) on NVMe for phase 1.** Longhorn when you have 3 nodes. Rook-Ceph only if you want to learn Ceph | Replicated storage across one machine protects nothing; local disks are faster and simpler |
| Power / network | **UPS (mandatory), dual WAN (recommended for TAK)** | Power cuts corrupt databases; TAK users on the field depend on being reachable |
| Hybrid link | **Tailscale** between cloud and home, Cloudflare Tunnel for web hosts | Already in use; no router ports needed |
| Cloud shutdown | Not before the exit criteria in section 7 | Credits last; moving slowly costs nothing |

### Talos vs k3s: reconciling with PR #142

PR #142 assumed k3s. That was a reasonable default and this plan keeps it for the
first home node, with the reasons above. Talos has real advantages we want
eventually: no SSH or package manager means a much smaller attack surface, the entire
node is configured by one YAML file you can keep in Git, upgrades are atomic, and rebuilding a node is
a 10-minute reinstall. Its costs for you today: no shell (every fix goes through `talosctl`
and its API), a new mental model, and extensions required for some storage and VPN
software (Longhorn, Tailscale, Ceph need system extensions). k3s also runs PR #142's
Rook (needs raw disks and 3 nodes for sensible use), Forgejo and Tekton unchanged.
Decision rule: **start with k3s on 1 node; if you add nodes and want a hardened,
appliance-like cluster, build the next cluster on Talos and migrate using
REBUILD-FROM-ZERO**. Because the platform is plain Kubernetes + Git + backups, switching
distribution later is a rebuild, not a rewrite. Keep it that way: do not use k3s-only
features (the bundled Traefik, Helm controller CRDs, k3s Klipper-specific annotations) in Git.

k3s install choices that match the current platform:
- Disable bundled Traefik (`--disable traefik`): we use ingress-nginx and Cloudflare Tunnel.
- Current cluster uses Cilium (Hubble UI, a Cilium Gateway for OpenChoreo). Either install
  Cilium on k3s (`--flannel-backend=none --disable-network-policy`, then Cilium) to stay
  identical, or accept flannel and re-test the NetworkPolicies and OpenChoreo gateway.
  Recommended: Cilium, for parity.
- Enable etcd snapshots to S3 (`etcd-s3` options) and test a restore: the control plane is now yours.
- Use SQLite (default) for 1 node; embedded etcd when moving to 3 server nodes.

## 2. Single node or three

| | 1 node | 3 nodes |
|---|---|---|
| Survives disk or board failure | No (restore from backup, hours) | Yes if storage is replicated |
| Survives maintenance reboots | Downtime of minutes | Rolling |
| Cost, noise, power | Lowest | About 3x |
| Learning load | Low | High (Longhorn/Ceph, quorum) |
| Fits current sizing? | Yes, see below | Yes |

Current load (7-day peak): 5.7 vCPU and 42 GiB RAM used; Kubernetes requests 27 cores and 64 GiB
(mostly over-requested, review H1). One strong machine (16+ threads, 64 GiB, ideally 128 GiB
with Ollama/WebODM headroom) runs everything. Three workloads are heavy:
TAK server (about 6 GiB), Ollama/LiteLLM (about 7 GiB), WebODM and geo stack (onetak, about 450 GiB of disk).

Recommendation: **1 strong node for the first 6-12 months**, plus the cloud as the second
site until trust is built. Add 2 more small nodes only once the single node has run
a quarter without trouble. A UPS plus tested restores gives more safety than three nodes in one room.
Three replicas on one physical computer do not survive losing that computer.

## 3. Power and network

- **UPS**: line-interactive or online, sized for the server + router + modem + switch for at
  least 15-30 min, with USB signaling to shut the server down cleanly (NUT or apcupsd). Test by
  pulling the plug.
- **Dual WAN**: two different ISPs (e.g. cable + cellular/5G or fibre), router with failover.
  TAK clients need a stable endpoint: use a DNS name that you can repoint, a Tailscale-based
  path, or a small cloud VPS front door with a static IP forwarding over Tailscale. Avoid CGNAT
  connections for any directly published port.
- Router: default deny inbound. Cloudflare Tunnel and Tailscale need only outbound access.
- Time: reliable NTP (TAK certificates and NATS auth depend on correct time).
- Disks: use NVMe (database), enable SMART monitoring, and keep backups on a separate physical device.

## 4. Storage choice

| Option | Use when | Notes |
|---|---|---|
| **local-path** (k3s) | 1 node | Zero setup; data lives on that node's disk; protected only by backups |
| **Longhorn** | 3 nodes | Replicated block storage, UI, snapshots, backup to S3; easiest HA option |
| **Rook-Ceph** (PR #142) | 3+ nodes with spare raw disks, and appetite to learn | Most powerful, most operational burden; needs explicit disk selection |
| OpenEBS LVM / ZFS-localpv | 1-2 nodes wanting snapshots | Good middle ground; needs LVM/ZFS knowledge |

Create StorageClasses named `do-block-storage` and `do-block-storage-retain` first (so
Git works unchanged), rename later. VolSync keeps working with `copyMethod: Direct`.
Capacity: provisioned 1,182 GiB on DO, of which about 246 GiB is idle. Plan 2-4 TB usable NVMe
(or 1-2 TB NVMe for databases plus a larger SATA/HDD pool for geo/media), and a separate backup disk.

## 5. How hybrid works

```
  Cloud (DOKS)                         Home (k3s, "onemind-home")
  Argo CD  --------- Tailscale ------> API server (registered cluster)
  NATS cluster <---- leaf/mirror ----- home NATS leaf + JetStream
  Authentik (IdP)  <-- OIDC ---------- home apps log in against it
  Backups -> S3 (off-DO) <------------ home pulls/pushes the same buckets
```

### Connectivity
Both clusters join the same tailnet. Home exposes its Kubernetes API only on the tailnet
(Tailscale operator or `tailscale serve`); no router ports. Web hosts through Cloudflare Tunnel.
TAK raw ports: see DO-LOCK-IN section 3.

### Argo CD managing both clusters
1. Keep the existing Argo CD in the cloud as the single control point for now.
2. Register home with `argocd cluster add <home-context> --name onemind-home` (the name PR #142's
   Applications already expect). The Argo service account token is stored as a Secret in the
   cloud cluster; treat as a credential.
3. Use a separate app-of-apps path for home (for example `platform/argocd-home/`) rather than
   editing the cloud tree, with home overlays: storage class names, no DO annotations, home node
   labels. Existing Applications do not change. PR #142/#139 apps (Rook, Forgejo, GPU serving,
   OCM) stay dormant until a "home activation" PR.
4. **Risk:** if the cloud dies, home loses its Argo. Before cloud decommission, install a second
   Argo CD at home (same repo, same overlays) and make it the owner of `onemind-home`.
   Argo CD reachable only from its own cluster is the safe end state.
5. Argo CD itself is a Helm-CLI release today (gap G3): adopt it first so the same definition
   is reused at home.

### NATS: leaf, mirror and the authority handover
Terms: a **leaf node** is a NATS server that connects outward to another and shares
subjects. A **JetStream stream** (stored messages) has exactly one home; other sites hold a
**mirror** (read-only copy, same sequence numbers) or a **source** (merged copy). See
docs/NATS-DOCTRINE.md (the Consume/Mirror/Merge table). The cloud cluster is currently the authority for all 39 streams.

Plan (validate in a drill; this is a design, not yet tested here):
1. Stand up a home NATS server with its own JetStream **domain** (e.g. `home`) as a leaf
   to the cloud cluster over the tailnet (`omos-nats` already exposes 4222/7422/8080 there).
   Use separate credentials for the leaf; never reuse operator/system credentials.
2. Create **mirrors** at home of the critical streams (CoT, events). Home now has a live replica.
   Writers stay on cloud. Measure lag and size.
3. Move consumers (projectors, ClickHouse sink, Hermes) to read from home mirrors. Compare outputs.
4. **Handover** in a short window (minutes): stop writers/feeds, wait until the mirror is caught
   up (sequence numbers equal), remove the mirror configuration at home so the stream becomes
   writable (promote), re-point publishers (TAK NATS leaf sidecar, feeds) at home, and change the
   cloud side to a mirror of home or leave it idle. Keep the old stream untouched for rollback.
5. Only one site is ever the writer for a stream. Never run two writers: that is the split-brain
   we avoid.
6. Backups: NATS stream archives (`backup-nats-streams`) keep running at the authority site.

### Databases
Each database has one primary site. Move one at a time, using a logical archive restore at
home, a read-only soak test, then cutover (stop writes, final dump, restore, repoint, start).
CNPG could later add a replica cluster at home for a smoother handover; not needed first.

### Identity
Keep Authentik (cloud) as the IdP in phases 1-3, with the issuer URL on a DNS name you control;
move it only with its database and signing keys (REBUILD-FROM-ZERO step 8). A home outage must
not lock you out: keep a break-glass local admin.

## 6. Hardware shopping guidance

Prices change; **check current price** at purchase time. Insist on amd64 (several custom images
are amd64-only; ARM boards will not run TAK or the geo stack unchanged).

| Tier | Typical class | CPU/RAM | Storage | Network | Good for | Risk |
|---|---|---|---|---|---|---|
| Budget | Used small-form-factor business PC or mini-PC (Dell/HP/Lenovo tiny, Intel NUC class) | 6-8 cores, 32 GiB (check max RAM) | 1 NVMe 1-2 TB + SATA backup disk | 1-2.5 GbE | Learning, game days, dev, a staging cluster | Not enough RAM for the full platform; one failure point |
| Mid (recommended first production node) | Workstation or rack/tower server (Ryzen 9 / Xeon-W / EPYC entry class) | 12-16+ cores, 64-128 GiB ECC if possible | 2 NVMe (mirror) 2-4 TB + backup HDD | 2.5/10 GbE | The whole current platform on one node | Single machine; needs UPS and tested restores |
| Mid+ (3-node) | 3 identical mini-servers or small towers | 8-12 cores, 32-64 GiB each | NVMe each (+ extra disks for Ceph) | 2.5-10 GbE + switch | Longhorn/Ceph, rolling maintenance | Cost, noise, power, learning load |
| GPU (later, PR #139) | Same mid node + a GPU | by model size: choose GPU VRAM from the model, not the reverse | NVMe for model files | | Local inference | Cannot be validated until hardware exists; check power supply and cooling |

Also buy: UPS, a managed or at least reliable switch, a second ISP link or LTE/5G router, a
separate backup drive or NAS (not the same machine), a small KVM/IPMI or a way to reach the console
remotely, and spare NVMe. Avoid consumer SSDs with tiny endurance for database workloads (check TBW).

## 7. Phased migration with exit criteria

Rule: no phase starts until the previous exit criteria are written down and met. Each phase
has a rollback that leaves the cloud untouched.

| Phase | What happens | Exit criteria | Rollback |
|---|---|---|---|
| **0. Today: secure the backups** | Enable the off-DO copy (OFFSITE-COPY.md) to Backblaze B2 or Cloudflare R2; copy the keys offline (K1-K4); wire Alertmanager (review A17); rotate exposed vault tokens (review C1) | A backup of every kind lands off DO within 24 h; `_offsite/last-success.json` fresh; alerts reach your phone; keys verified readable from offline copies | Disable the new job |
| **1. Lab: prove the runbook** | Buy budget or mid hardware; install k3s; run the REBUILD-FROM-ZERO game day into it, fully isolated (no production NATS, no live DNS, no EUD certs) | Checklist passes; time recorded; every improvisation turned into a doc/script PR; the 11 Helm-CLI releases adopted or captured | Power it off |
| **2. Hybrid: home joins** | Home on the tailnet; register in Argo; deploy low-risk, stateless/ecosystem workloads (Ollama, Harbor mirror, Forgejo, monitoring replica); Longhorn/local-path chosen; decide TAK front-door design | Home runs real workloads for 4 weeks with no data loss; UPS test done; dual WAN failover tested; home second Argo installed | Scale the home workloads to 0 |
| **3. Replica: data follows** | NATS leaf and mirrors; CNPG/Postgres restored regularly to home; home holds a warm copy of everything (read-only) | Home can serve a full read-only copy; a timed failover rehearsal succeeds; RPO and RTO measured and acceptable | Remove mirrors |
| **4. Cutover by service** | One service at a time becomes primary at home, in this order: monitoring/observability, ecosystem apps, n8n/Backstage/Matrix, Authentik, NATS authority, TAK last | Each service: 2 weeks stable at home, backups running at home and restore-tested, cloud copy kept idle | Re-point DNS and writers back; cloud copy is intact |
| **5. Cloud shrink** | Right-size cloud to the minimum (review: option 1 saves about $168/month; stop idle pools) while home is primary; cloud becomes a backup/outpost | A month with home primary and zero incidents; game day passes at home; offsite copy verified | Scale up cloud again |
| **6. Decommission** | Take final encrypted backups and volume exports, verify restores, then and only then shut down DO resources, with explicit written approval from the owner for each deletion (cluster, Spaces buckets, volumes, snapshots, LBs). Keep one cheap off-site backup and a cold-standby plan | Two successful consecutive game days from backups only; the cloud shutdown rehearsed with the cloud cluster scaled to zero for a week without impact | Re-create from backups (REBUILD-FROM-ZERO) |

Budget guard: DO costs stop only when resources are destroyed. Credits (about $2,500) cover roughly
3 months at list price; keep a calendar reminder at 60 days remaining to decide Phase 5/6. Verify
the credit balance and expiry in the DO dashboard (review noted billing shows $0, i.e. credits).

## 8. Open decisions (need the owner)

1. Where does the off-DO backup copy live (Backblaze B2, Cloudflare R2, both)?
2. How do TAK clients reach home (Tailscale only, VPS front door, or router forwarding)?
3. 1 strong node or start with 3 small nodes?
4. Which ISP pair, and is a static IP available?
5. Is GPU hardware planned in the first purchase (PR #139)?
6. Who is the second person with access to keys and this runbook if you are unavailable?
