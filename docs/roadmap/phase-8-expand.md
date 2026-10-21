# Phase 8 brief: expand (home server, Hermes)

Revised preparation: 2026-10-20, America/New_York. **Draft; implementation and
production acceptance remain gated.** This brief belongs to existing
[#284](https://github.com/onemindos/omos/pull/284), not a new expansion PR.
Tracker: [`PLATFORM-ROADMAP.md`](../../PLATFORM-ROADMAP.md); read
[`HANDOFF.md`](HANDOFF.md) and [`MISSION.md`](../MISSION.md).
This branch's tracker/handoff contain older phase status. The coordinator must
reconcile them against current main before updating #284; this docs preparation
does not advance any phase. Phase 6 remains unfinished in the latest supplied
handoff; later execution follows the shared tracker and approvals.

## Goal and scope

Home becomes the primary operational authority, with cloud resources serving
selected workloads, reachability and recovery. A cloud-dependent dashboard or
leaf connection alone does not satisfy local operation during WAN loss.

Propose a dedicated **ARM64 Linux VM running K3s on the existing Mac** for a
bounded pilot. K3s runs inside Linux, not natively on macOS. This is neither an
approved VM allocation nor a production runtime decision. Talos remains a
candidate for a future dedicated host; replacing macOS is not required here.
Preserve Docker Desktop's builder role, existing Home Assistant and user data.
A single Mac/VM is one failure domain, not HA.

This revision authorizes no installation, storage reuse, power-setting change,
credential/provider/billing change, migration or controller deployment. The
coordinator owns tracker updates, tests and publication on the existing PR.

## Evidence ledger: measurements are not allocation approval

These are previously completed observations, not a fresh runtime/network test.
The original record is the 2026-10-20 preflight carried by
[#332](https://github.com/onemindos/omos/pull/332).

| Evidence date / source | Recorded observation | Remaining production gate |
|---|---|---|
| 2026-10-20, 16:15–16:27 UTC preflight | M2 Max, 32 GiB RAM, macOS 27.0.1; internal Data volume only 13–14 GiB free | Approve runtime, vCPU/RAM/disk sizing and host reserves against actual workload demand. |
| 2026-10-20 preparation, America/New_York; earlier read-only `uname -m`, `sysctl -n hw.memsize`, `df -h /System/Volumes/Data` (exact UTC time not captured) | ARM64, 32 GiB RAM; internal volume later showed 26 GiB available, 98% capacity | Capacity varies. Neither observation approves a VM; budget guest OS, images, database/log growth, snapshots and host safety margin before allocation. No cleanup or resizing was performed. |
| 2026-10-20 preflight | External `OM PRIMARY` holds existing data; `MBP TM` is a backup disk | Neither disk is approved spare capacity. No repurpose, file move or backup change. |
| 2026-10-20 preflight and later preparation settings read | AC sleep/disksleep/standby/display sleep set to zero; earlier FileVault on | Settings do not prove plugged-in closed-lid cooling, continuous service or unattended unlock/recovery. |
| 2026-10-20 preflight | Lima installed with no instances; stopped ARM64 Colima builder profile, 4 CPU / 8 GiB / 24 GiB | A build profile is not a production cluster or capacity reservation. Later command discovery confirms installed tools only. |
| 2026-10-20, 16:17:38–39 UTC Legacy preflight | Configured DigitalOcean provider returned HTTP402; HTTP200 SSE ended in error without an answer | Inference remains failed in recorded evidence. No new probe or repair is claimed. |

## Decisions before implementation

- Approve host/runtime and CPU, memory, disk and growth reserves; storage location,
  retention and independent backup destination; service inventory and owner.
- Set outage duration, RPO/RTO, essential local functions and safe degraded modes.
  Decide which changes must remain possible without the cloud control plane.
- Define local authentication/authorization and recovery when cloud Authentik,
  DNS, Tailscale coordination or inference are unreachable. Existing sessions
  alone do not prove new local login or identity recovery.
- Decide whether/when control-plane authority moves home. A second remote data
  plane does not itself transfer authority or remove cloud dependencies.
- Approve later model/provider or budget decisions separately. Preserve existing
  Legacy configuration until then; see the dependent #332 gates.

## A. Mac production acceptance before cluster registration

For every later test record approval, exact versions, UTC start/end, target,
expected result, actual evidence, rollback and PASS/FAIL/BLOCKED. All rows below
are **NOT RUN** in this preparation. Fault/reboot tests require their own bounded
window and must exclude existing Home Assistant and unrelated workloads.
If a host-wide reboot would interrupt them, stop until that impact has its own
explicit maintenance approval; a pilot-only fault window cannot cover it.

| Gate | Later operator procedure | Pass evidence |
|---|---|---|
| Capacity / ARM64 | With approved sizing, run selected workloads in the dedicated Linux VM; measure host/guest pressure, disk growth and image compatibility under representative load. | Approved reserves maintained; pinned ARM64 images actually execute; storage/IO and recovery budgets fit. Multiarch manifests alone are insufficient. |
| Closed lid / power | With approved physical setup, operate plugged in with lid closed for an agreed duration under representative load; observe services, temperatures and guest health. | No unplanned sleep or thermal/service failure; record duration and interruptions. Do not change power settings as part of this brief. |
| Logout / reboot / FileVault | In the approved window, test logout, host reboot and a separately approved cold-start/power-recovery scenario. Follow the documented unlock and VM/service startup procedure. | Measure downtime and human intervention; preserve FileVault and prove authorized recovery without cloud access. If unlock requires a person, document that availability limit rather than calling it unattended. |
| Local identity / access | From an approved local client, verify local DNS/TLS, new login, authorization denial and recovery access with cloud identity unavailable in the isolated test scenario. | Essential local services remain reachable to intended roles; denied roles stay denied; emergency credentials have approved independent custody without disclosure. |
| WAN / cloud loss | In an approved isolated scenario, remove only the pilot's cloud/WAN dependency. Test local reading and permitted control separately, local DNS/identity, event handling and required local inference. Restore connectivity and inspect replay/conflicts. | Essential functions work for the outage target; unavailable cloud functions are explicit; no unsafe commands, duplicate replay or lost authoritative state. Existing-session access is separate from new login. |
| Data / identity recovery | Restore an application-consistent encrypted independent backup into an isolated target through the approved recovery procedure and identity/key recovery references. | Verify known synthetic data, authorization and measured RPO/RTO without overwriting production. An extra copy on the same Mac is not independent recovery. |

## B. Distinct integration gates

Execute only after approved foundation prerequisites, through GitOps-owned
resources. Preserve SealedSecrets as the repository credential pattern; controller
presence alone is no reason to add External Secrets or another foundation owner.
Review private exposure against `platform/core/policy.hujson`; mesh membership
does not grant equipment or infrastructure authorization. AI access above
Layer 3 remains read-only under the roadmap; this brief expands no permissions.

1. **Real home cluster:** record stable cluster identity, actual ARM64 nodes,
   private network path, local recovery procedure and exact versions. A context,
   catalog location or Docker builder is not a registered second cluster.
2. **Argo CD registration:** review scoped cluster identity/access and GitOps
   ownership before registration. Subsequently prove the actual home destination
   is connected and an approved disposable workload reconciles there, without
   targeting cloud resources. Record Backstage visibility against that real
   destination. This proves Argo reachability/reconciliation, not OpenChoreo
   connectivity or local offline autonomy.
3. **OpenChoreo plane:** separately verify the pinned release's supported
   multicluster contract against current official documentation and actual gateway
   state at execution time. The old brief's port/path, Service type, replica and
   manual-CA assumptions are historical proposals, not current observations.
   Review private gateway exposure, TLS trust, per-plane identity and certificate
   rotation/revocation. Prove a distinct home `ClusterDataPlane`/plane ID and
   Environment actually connect; then place an approved disposable Component
   using the intended ReleaseBinding. Record pod location, portal management,
   disconnect/rejoin and continued workload operation independently of
   control-plane mutation availability. Argo registration is insufficient evidence.
4. **Local NATS:** follow [`NATS-DOCTRINE.md`](../NATS-DOCTRINE.md) and existing
   `infra/edge/mac-leaf/` contracts in a separately reviewed change. A forwarding
   leaf is not independent JetStream authority. Verify local outage behavior,
   retention, rejoin/replay and identity boundaries before proposing authority
   transfer; do not invent subjects or redirect cloud workloads to the laptop.
5. **Fleet / edge pilots:** only after two independently recoverable clusters,
   consider a bounded Karmada placement/override pilot and a separate KubeEdge
   edge-node/device lifecycle pilot if each has a demonstrated need. Neither
   substitutes for Argo registration or OpenChoreo acceptance. Define one
   reconciler per resource, placement/storage rules, disconnect/reboot/rejoin,
   revocation and rollback; no stretched WAN etcd or duplicate ownership.

## C. Hermes and TAK dependencies

Keep the conversation pilot and acceptance matrix on existing
[#332](https://github.com/onemindos/omos/pull/332). Metadata/discovery success is
not working inference. Provider402 resolution, dedicated identity/credentials,
two-conversation/session isolation, database persistence, encrypted independent
restore and approval/cancellation/spend controls remain separate gates. A cloud
provider repair would not establish offline local inference.

The proposed read-only MCP route needs dedicated agent identities and enforced
per-tool scopes; the existing privileged Legacy harness must not be described as
read-only because a UI or prompt says so. Evaluate the model/tool gateway under
the roadmap without changing providers, credentials or budgets here. Hermes
moves last in Phase 7 only after its migration gates are satisfied.

Revisit TAK placement after the home failure domain is accepted and current
Helm/Argo ownership and data recovery are verified. No TAK move is authorized.

## Completion gate

Keep #284 draft until the coordinator records actual hardware/runtime, recovery
and local WAN-loss acceptance; real Argo registration; separate home OpenChoreo
plane/component and portal acceptance; and required Hermes read-only tool evidence.
Record partial results per gate rather than marking the phase complete.
Karmada/KubeEdge remain separately gated optional pilots, not required installed
services. Documentation, inventory and source review certify no deployment.
