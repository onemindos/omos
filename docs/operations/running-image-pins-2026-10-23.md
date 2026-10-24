# Running-image pins — 2026-10-23

Task 2 pins the mutable references with a uniquely observed running imageID in
`do-nyc3-omos-cloud`. No application version upgrade is selected. The Pod status
imageIDs were captured read-only; registry OCI metadata was inspected without
pulling or running images. Where an application release tag could not be established,
the digest alone identifies the exact running version; base-image labels were not
misrepresented as application versions. NodeODM 3.6.2 and Mumble v1.5.915 come from
application OCI version labels; BusyBox 1.38.0 comes from OCI build history.

| File | Old | New |
| --- | --- | --- |
| `platform/kube/headlamp.yaml` | `ghcr.io/headlamp-k8s/headlamp:latest` | `ghcr.io/headlamp-k8s/headlamp@sha256:8a0ef97527544baa2b35825940ddf4db6e2d09c30cd5aa31fd269ff604cc8b75` |
| `platform/monitoring/deployment.yaml` | `natsio/nats-surveyor:latest` | `natsio/nats-surveyor@sha256:7d0ca8d7ec52645b8366aa601762aa64d9690444584eb7311c9c4922216f7e2b` |
| `platform/nats-tools/heynats/deployment.yaml` | `ghcr.io/onemindos/heynats:latest` | `ghcr.io/onemindos/heynats@sha256:65f1d5aa46f21b14baba5942e2e5df3ec2e3d8700796cab9a7ac9d0361dffe92` |
| `platform/nats-tools/nats-dashboard/deployment.yaml` | `ghcr.io/onemindos/nats-dashboard:latest` | `ghcr.io/onemindos/nats-dashboard@sha256:bccc7a25725540d7db4147544bf25d45716fdce86b269ec3c9aa87ca3e00da38` |
| `platform/nats-tools/nats-ui/deployment.yaml` | `ghcr.io/gastbob40/nats-ui:latest` | `ghcr.io/gastbob40/nats-ui@sha256:19f13ef3454bea03c8673ce7ac6702bce7bdd917ae82506f82e609027d136615` |
| `platform/onetak/tak-portal.yaml` | `ghcr.io/onemindos/tak-portal:latest` | `ghcr.io/onemindos/tak-portal@sha256:eaa1472927823165df82c715c9cd77da36c5d503f75c4172d1215c679fe27f6b` |
| `platform/onetak/tile38.yaml` | `tile38/tile38:latest` | `tile38/tile38@sha256:5b61c6f3b32e508bf447b61b53320b3e524a6d9679782730fd8c930451b53e36` |
| `platform/onetak/tileserver-gl.yaml` | `maptiler/tileserver-gl:latest` | `maptiler/tileserver-gl@sha256:3a9ccdb24820b6814c8119bcc8a4376c39867cb0ffe69d62919ef898b90c2427` |
| `platform/onetak/webodm.yaml` | `webodm/webodm_db:latest` | `webodm/webodm_db@sha256:03f18ec0089325ec2b8975ba1781277bcdc212b6d86c9de33f6409d313e90430` |
| `platform/onetak/webodm.yaml` | `webodm/webodm_webapp:latest` | `webodm/webodm_webapp@sha256:875cf97ee3278d9dfc517feba862ebe83aae1ccf76c398de96e993402f181659` |
| `platform/onetak/webodm.yaml` | `webodm/webodm_webapp:latest` | `webodm/webodm_webapp@sha256:875cf97ee3278d9dfc517feba862ebe83aae1ccf76c398de96e993402f181659` |
| `platform/onetak/webodm.yaml` | `opendronemap/nodeodm:latest` | `opendronemap/nodeodm:3.6.2@sha256:fcd99eb23d8db194a7d1dc95f029f518283b3d6be32dd64422395f4e7af56720` |
| `platform/takserver/deploy/murmur.yaml` | `mumblevoip/mumble-server:latest` | `mumblevoip/mumble-server:v1.5.915@sha256:018ad3515932e513d8fdc970df918373a2664c4da32b09a89a4aaee28eb7507a` |
| `platform/matrix/argo/postgres.yaml` | `busybox` | `busybox:1.38.0@sha256:fd7dc98638c8e305f4dc34e979f1c0fdfdcaeb0fbf8fcff77ae834b6da3d7e6e` |
| `platform/matrix/argo/redis.yaml` | `busybox` | `busybox:1.38.0@sha256:fd7dc98638c8e305f4dc34e979f1c0fdfdcaeb0fbf8fcff77ae834b6da3d7e6e` |
| `platform/matrix/argo/synapse.yaml` | `busybox` | `busybox:1.38.0@sha256:fd7dc98638c8e305f4dc34e979f1c0fdfdcaeb0fbf8fcff77ae834b6da3d7e6e` |
| `platform/matrix/base/synapse/manifests/deployment.yaml` | `busybox:latest` | `busybox:1.38.0@sha256:fd7dc98638c8e305f4dc34e979f1c0fdfdcaeb0fbf8fcff77ae834b6da3d7e6e` |
| `platform/onetak/webodm.yaml` | `busybox` | `busybox:1.38.0@sha256:fd7dc98638c8e305f4dc34e979f1c0fdfdcaeb0fbf8fcff77ae834b6da3d7e6e` |

## Remaining references and exclusions

The following references have no verified runtime bytes for their target environment.
They were not replaced with invented versions or upgraded to an unrelated release:

- `platform/monitoring/scraper-deployment.yaml:22` — `ghcr.io/onemindos/nats-metrics-scrape:latest`.
- `platform/parked/omos-mission-engine/mission-engine-deployment.yaml:24` — `ghcr.io/onemindos/omos-mission-engine:latest`.
- `infra/edge/takserver/deploy/docker-compose.yml:12` — `tak-database:latest`.
- `infra/edge/takserver/deploy/docker-compose.yml:33` — `takserver:latest`.
- `infra/edge/takserver/deploy/docker-compose.yml:67` — edge `mumblevoip/mumble-server:latest` (edge runtime not inspected).

The edge `takserver` and `tak-database` images are local build outputs, not registry
releases. Their Dockerfiles/compose build configuration need a separate reproducible
edge build before a trustworthy pin exists. The parked mission engine is not deployed,
and the obsolete metrics scraper has no matching pod. These are explicit residual
Task 2 findings; this change does not claim a zero-mutable-reference audit.
The edge Mumble example stays unchanged: cloud runtime evidence does not establish
which bytes the edge is running. This PR is limited to 18 verified cloud image
references, including the matching legacy Matrix base reference.

The running OpenSearch `fsgroup-volume` BusyBox reference belongs to an off-limits
ecosystem OpenSearch application and is skipped. Excluded paths stay unchanged.
`run-latest` and existing version-tag references are outside the literal latest/untagged
audit scope; no registry-wide tag refresh was performed. The policy-report API returned
no PolicyReports at observation time, so the earlier count of 17 was not recertified.

## Current-main integration and evidence refresh

This branch is rebased onto `85b9688ffba68fa596295db1ccd5d9b0b6d0455e`, preserving
merged resource budgets and the portal cleanup. The targeted running imageIDs were
rechecked after the resource rollout, including the replacement WebODM database
pod. Its digest remains the same as this pin. BusyBox, NodeODM and Mumble release
tags were also verified to resolve to their captured digests. No version upgrade
or rollback is intended; no edge version is inferred from cloud evidence.

The original scan scope is platform/ and infra/. Development compose references
under services/bridge/ are outside that scope and unchanged. Policies/comments,
parked images and local build outputs must not be mistaken for deployed workloads.

## Reviewer steps after merge

Allow only the owning Argo Applications to reconcile. `platform-matrix` has no
automated sync: use the separately reviewed Matrix drift runbook before syncing it.
Do not blindly sync all Matrix drift merely to pick up these init-container pins.
Observe ready counts and compare each new Pod's runtime imageID against this table;
all replicas for a pinned image must run the captured bytes. A rollout can occur even
when the underlying bytes are unchanged, and a currently mutable tag may have moved
since capture. Recheck live digests before merge if another release lands.
Do not deploy edge compose or parked manifests as part of this PR. No cleanup or
secret operation is required. Any failed image pull or regression should be handled
through a reviewed Git revert, not an imperative rollout or a switch back to latest.
