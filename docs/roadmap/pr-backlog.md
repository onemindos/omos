# OneMind PR backlog

Written 2026-11-05 after the full review of open PRs. Each entry stays open with the `backlog` label. Hermes PRs #640 and #646 belong to the Hermes sessions and are not listed.

## Do next

| PR | What it does | Waiting for |
|---|---|---|
| [#562](https://github.com/onemindos/omos/pull/562) | Gives every service its own NATS login and permissions instead of one shared password, with staged activation and tests. | Per #585: allow hermes STANDARD connections, drop the leaf.conf/sidecar rendering and add _INBOX.hermes.>; also confirm Argo ignores the three new values.identities.*.template.yaml in the non-recursive platform/nats path (unverified). |
| [#553](https://github.com/onemindos/omos/pull/553) | Adds a reliable feed collector (leader election, dedup, dead letters, freshness alerts) and moves the National Weather Service alerts onto it. | Owner ruling BUS-P1 (tel.feed.>), CI image build and pinned digest, then manual sync of platform-ingest. |
| [#547](https://github.com/onemindos/omos/pull/547) | Moves the NATS system password out of the plain-text deployment file into a sealed secret. | Needs the NATS $SYS password to be rotated in the NATS config to match the committed SealedSecret (last commit says rotated; unverified; the Secret nats-surveyor-creds does not exist in the cluster yet). Then merge and check surveyor. |
| [#587](https://github.com/onemindos/omos/pull/587) | Replaces TAK's default passwords and unencrypted key files with new random ones stored in sealed secrets. | Six SealedSecrets not generated yet; needs the window (TAK down about 5 to 8 min, FedHub/TAK-ML 1 to 2 min) and an Argo pause. Rebase check: no conflicts. |

## On hold by the owner

| PR | What it does | Waiting for |
|---|---|---|
| [#533](https://github.com/onemindos/omos/pull/533) | Makes the monitoring system send critical and warning alerts to Telegram; today alerts go nowhere. | Sealed secret is already committed and CI passes; I did not verify the token/chat id decrypt or that Telegram delivery works. No conflicts. |

## Cluster-changing: owner decides

| PR | What it does | Waiting for |
|---|---|---|
| [#555](https://github.com/onemindos/omos/pull/555) | Puts size caps on every NATS stream and cuts the TAK message history (COT) from 30 days to 72 hours, which deletes about 4.4 GB per node. | When NATS disk or a runaway stream becomes a problem, or the COT history is no longer wanted. |
| [#558](https://github.com/onemindos/omos/pull/558) | Replaces the MTA train and satellite feeds with versions that publish plain observations on grammar subjects, leaving TAK output to the entity service, and removes the unused feed-hazards. | After #552 can read rail and TLE data (A1-P4 follow-ups) and BUS-P6 is ruled. |
| [#563](https://github.com/onemindos/omos/pull/563) | Builds a new TAK projector that only sends safe, relevant entities (no people, marking-checked, area-gated) to TAK and serves a checked send for web and Unreal clients. | When the entity and layer services run and TAK delivery (#560) works. |

## Waits for a feature or device

| PR | What it does | Waiting for |
|---|---|---|
| [#561](https://github.com/onemindos/omos/pull/561) | Builds the layer service so a map client can ask for one snapshot and then a stream of changes for an area. | When the web or Unreal client needs a map view (onemind-xr #2/#5). |
| [#564](https://github.com/onemindos/omos/pull/564) | Adds a Hermes plugin that lets the agent read entities and layers, propose commands (never confirm) and drive the map, only through an allowlist. | After #554, #552, #561 and the #562 hermes identity run, and the derived Hermes image pipeline exists. |
| [#565](https://github.com/onemindos/omos/pull/565) | Lets people log in through Authentik and gets a role (newcomer to admin) that decides what they can do on the bus, via a NATS auth callout. | When remote humans (web, headset) connect to the bus. |
| [#566](https://github.com/onemindos/omos/pull/566) | Adds a ROS 2 robot bridge (via zenoh) that sends robot position and takes the same approved, signed commands as the drone bridge. | A ROS 2 robot or simulator to attach. |
| [#567](https://github.com/onemindos/omos/pull/567) | Makes agents appear as map contacts, runs deterministic watch jobs that raise alerts and propose actions, and writes area briefs that work offline. | After #552, #560 and the agent protocol on onemind-hermes work. |
| [#568](https://github.com/onemindos/omos/pull/568) | Builds the task service that applies approved mission and task commands, tracks task state and answers who is doing what. | The first task-bearing command flow (policy plus a bridge). |
| [#569](https://github.com/onemindos/omos/pull/569) | Builds the geofence and alert service: stores keep-in/keep-out zones and raises alerts when assets cross them or go stale. | The first zone you want enforced, after #552 and the D1-P4 policy patch. |
| [#559](https://github.com/onemindos/omos/pull/559) | Adds an authenticated video server (MediaMTX), a video registry that lists streams as entities, a plugin-manifest validator, and an edge relay script. | First camera or drone video to carry, or the client video panel (K4). |

## New apps not chosen yet

| PR | What it does | Waiting for |
|---|---|---|
| [#605](https://github.com/onemindos/omos/pull/605) | Runs five chat bridges (WhatsApp, Telegram, Google Messages, Instagram, X) so those chats appear in the Beeper app. | When Beeper login, Telegram API keys and backup credentials are ready and a home cluster or spare capacity exists. |
| [#609](https://github.com/onemindos/omos/pull/609) | Adds a private audio/video relay server meant for a future voice conversation with the AI and robot camera feeds. | When there is a client app for talking to onemind-hermes (Android/Mac) or a robot camera to relay. |

## Plans to refresh

| PR | What it does | Waiting for |
|---|---|---|
| [#334](https://github.com/onemindos/omos/pull/334) | A plan (no changes) for replacing ingress-nginx, which upstream stopped maintaining in March 2026, with the Cilium Gateway API. | Before adding any new public service, or when the home tunnel/public-exposure design is chosen; ingress-nginx still serves 4 public ingresses (authentik, nats-ws, geo-services, onetak). |
| [#339](https://github.com/onemindos/omos/pull/339) | A plan, plus a small offline test script, for letting AI request changes that a human approves and a limited identity executes, instead of giving the AI admin rights. | When onemind-hermes or the tool gateway first needs to write (execute gate, tool-gateway phase 4). |
| [#399](https://github.com/onemindos/omos/pull/399) | A written plan to lock down the platform: per-agent NATS credentials first, then default-deny network rules, Tailscale rules, signed images. | Start when NATS per-service identities (#562) are decided; network policies can start any time. |
| [#441](https://github.com/onemindos/omos/pull/441) | A plan and cluster snapshot for signing in once through Authentik across all OneMind apps. | When the app set settles after the slim-down; start with Backstage, Argo CD, Grafana. |
| [#585](https://github.com/onemindos/omos/pull/585) | A design for how onemind-hermes connects to NATS directly, to TAK and to the plugin registry. | After #560, #562, #564 (in the 550-569 range, not reviewed here) are decided; then merge as the record. |

## Dependency work

| PR | What it does | Waiting for |
|---|---|---|
| [#187](https://github.com/onemindos/omos/pull/187) | Renovate batch that pins GitHub Action versions to exact commits and also bumps Python 3.12 to 3.15, Terraform 1.9 to 1.17, Go 1.26 to 1.27 and Helm. | After Renovate is told to skip the Python/Terraform/Go version bumps (or those versions exist on the runners). |

## Known follow-ups on what is already live

- **Trivy:** some scans fail on a cache lock error, `onetak` has no results yet, and the operator logs a forbidden error for `clustercompliancereports`. The first findings include critical items in n8n and Authentik.
- **Off-site backups and Telegram alerts:** on hold by the owner.
- **Bus-to-TAK:** #560 is merged but not started. Start it only after checking what reaches TAK.
