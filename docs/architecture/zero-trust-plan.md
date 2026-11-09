# Zero-trust plan for OneMind

Status: proposal, docs only. Nothing here is applied. Facts were read from the repo and, read-only, from the
`do-nyc3-omos-cloud` cluster on 2026-10-23. No secret values appear in this document; only user names,
subject names and object names.

Provider rule: DigitalOcean only. No AWS tooling or services anywhere in this plan.

## 1. Goal and approach

Zero trust here means: no request is trusted because of where it comes from. Every caller (person, device,
service, agent, AI action) proves who it is, and is allowed only what it needs, for as long as it needs it.

It has to be native to the OneMind stack, not a new product on top of it. The OODA loop (Sense, Make Sense,
Decide, Act) runs on NATS, so the bus is the main enforcement point. The pieces already in the cluster map to
layers like this:

| Layer | Native piece | Role |
|---|---|---|
| People | Authentik | SSO, groups, MFA |
| Devices and private access | Tailscale operator, tailnet policy | Device identity, who reaches what |
| Services (network) | Cilium NetworkPolicy, Dapr mTLS, NATS TLS | Default deny, encrypted service calls |
| Agents | NATS accounts and users, Synadia agent protocol | Per-agent identity and subject permissions |
| AI actions | OpenFGA plus the permissioned action layer (PR #339) | Approve, then act |
| Workloads and supply chain | Kyverno, GHCR images, cosign | Admission rules, signed images |
| Secrets | Sealed Secrets (External Secrets is being introduced) | No plaintext in Git |

SPIRE is being removed, so this plan does not depend on SPIFFE identities. Workload identity comes from
Kubernetes ServiceAccounts, Cilium identities and (for NATS) per-client nkeys.

## 2. Research: current state

### 2.1 NATS authentication

Server: nats-server 2.14.2, 3 replicas (`nats-0..2`) in namespace `nats`, TLS required on 4222, JetStream on.

Mode: **static, config-file accounts with password users**. No operator, no JWT, no resolver, no auth callout.
The live `nats-config` ConfigMap is edited outside Helm (see `platform/nats/OPERATIONS.md`); the checked-in
`platform/nats/values.template.yaml` is a sanitized template of it. Passwords are stored as bcrypt hashes.
No user has a `permissions` block, so every user can publish and subscribe to everything in its account.

| Account | User | Used by (from repo plus live connections) | Permissions today |
|---|---|---|---|
| ONEMIND | `onemind` | Everything on the cloud side: omos feeds (feed-intel*, omos-tak-projector, sink-clickhouse, redpanda-connect), platform-backup, takml vision worker, TAK mumble-bridge, Agent Fabric (agent-console, agent-tap), `testui` and other ad-hoc clients, nats-box contexts | None (full access, JetStream admin included) |
| ONEMIND | `leaf_tak` (leaf listener 7422) | Mac leaf (`zeus-mac`), Hermes cloud leaf sidecar (`hermes-legacy-leaf`), TAK bridge leaf | None (leaf users bind into ONEMIND) |
| SYS | `sys` | nats-surveyor, nats-dashboard | System account |
| SYS | `monitor` | Monitoring | System account |

Other facts:

- The ONEMIND account imports `$SYS.REQ.>` from SYS, so `onemind` can use system requests.
- Live `connz` on the three servers (13 client and leaf connections at that moment): every ONEMIND connection is `onemind` or `leaf_tak`. There is no way to tell one agent from another by identity.
- WebSocket on 8080 is `no_tls` on the server. It is published on the tailnet (`omos-nats`, `nats-ws`) and, through nginx, on the internet as `ws.onemindos.dev` (see 2.6).
- NACK (JetStream controller) reconciles the 39 streams and KV buckets in `platform/nats/streams/`.
- Agent protocol subjects follow `agents.{verb}.{framework}.{owner}.{session}` (`docs/NATS-DOCTRINE.md`,
  `platform/hermes/agent/AGENTS.yaml`), for example `agents.prompt.hermes.zeus.legacy`. The verb comes first,
  so "each agent only gets `agents.<name>.>`" cannot be written literally; permissions have to be written per
  verb with the framework, owner and session tokens fixed (section 5, step 1).
- The Doctrine (section 5) already plans accounts SYS, CORE, FABRIC, INTERFACE and says the cutover to JWT
  should be "a config change". Section 5 step 1 notes why that is not fully true for JetStream data.

### 2.2 NetworkPolicy and Cilium

- Cilium v1.19.8, managed by DOKS (`doks.digitalocean.com/managed=true`), kube-proxy replacement, native routing,
  `enable-policy=default`, `enable-k8s-networkpolicy=true`, Hubble on, policy verdict events on. So standard
  NetworkPolicy and CiliumNetworkPolicy objects are enforced.
- Cluster-wide: one CiliumClusterwideNetworkPolicy, `deny-imds-egress`. Namespaced CiliumNetworkPolicies:
  `github-runner/hermes-runner-api`, `ops/cluster-recovery-api`, one generated in a `dp-*` namespace.
- **Namespaces with a real default-deny (pod selector `{}`): `agent-fabric` (ingress and egress) and `n8n`
  (ingress and egress).** `openchoreo-observability-plane` has a select-all ingress policy with allow rules, but
  no egress policy.
- Namespaces with partial policies (they select a few pods, so everything else is still open): argocd, auth,
  backstage, clickhouse, dapr-system, github-runner, hermes, matrix, monitoring, nats, onemind-ai, onemind-core,
  onemind-data, onemind-registry, onemind-search, onetak, ops, spire-system, takml, takserver, and one `dp-*`.
  Many of these only hold the shared `platform-volume-backup-egress` policy.
- **Namespaces with no policy at all:** cert-manager, cloudflare, cnpg-system, default, external-dns,
  external-secrets, falco, ingress-nginx, keda, kube, kube-system, kyverno, omos, onemind-jobs, onemind-media,
  opencost, openchoreo-control-plane, openchoreo-data-plane, openchoreo-secrets, openchoreo-workflow-plane,
  tailscale, velero, and five of the six `dp-*` namespaces.
- Notable: `nats`, `hermes`, `omos` and `takserver` hold the most sensitive traffic and are not default-deny.

### 2.3 WireGuard transparent encryption

- **Off.** `cilium-dbg status` reports `Encryption: Disabled`; `cilium-config` has no `enable-wireguard`,
  `encryption` or IPsec keys.
- The Cilium DaemonSet and ConfigMap are DOKS-managed. DigitalOcean has an open feature request for WireGuard on
  DOKS (ideas.digitalocean.com, "Cilium WireGuard Encryption for DOKS"), which indicates it is not a supported
  toggle. Editing `cilium-config` on a managed install is expected to be reverted; the only route is replacing
  the managed CNI with a self-managed one, which is out of scope and risky. Confirm with a DigitalOcean
  support ticket before ever trying it.
- Consequence: pod-to-pod traffic on the node network relies on application-layer encryption. Today that is
  NATS TLS (client, cluster and leaf), Dapr mTLS for `onemind-core`, and Tailscale (WireGuard) for anything that
  crosses sites. Everything else in cluster is plaintext. See gap G8.

### 2.4 Tailscale

- The tailnet policy **is in the repo**: `platform/core/policy.hujson`, applied by `.github/workflows/tailscale.yml`
  (diff on PR, tests, apply on merge to main).
- Tags: `tag:k8s-operator` (operator) and `tag:k8s` (every proxy the operator creates, set by `PROXY_TAGS`).
  One ProxyGroup, `egress`, plus one proxy per exposed Service or Ingress (about 25 in namespace `tailscale`).
- Grants:
  1. `group:admins` (one person, GitHub login) reaches everything on every port (`dst *`, `ip *`).
  2. `tag:k8s` reaches `omos-mbp:22`, `omos-ha:8123`, `omos-s24:8022`.
  3. `tag:k8s` reaches `argocd:443` and `onemind-auth:443` (Backstage mesh).
- Tailscale SSH: `autogroup:member` to `autogroup:self`. Funnel is not granted, and no Service or Ingress carries a
  funnel annotation (checked live).
- What this means: all proxies share one tag, so the Hermes egress proxies and, for example, the Grafana ingress
  proxy have the same reach. A compromised proxy pod can use all three tag:k8s grants. The tests cover the main
  deny cases. Every device signed in as the admin is fully trusted, with no device posture check.
- Tailnet-only surface: argocd, Authentik mirror, Grafana, Headlamp, Kubernetes dashboard, Hermes dashboard, n8n,
  Hubble, OpenChoreo observer and gateways, OpenCost, NATS UIs and monitor, Agent Console, Agent tap,
  Backstage preview, the `omos-nats` LB (4222, 7422, 8080) and the `tak` LB (8089, 8443, 8446, 64738).

### 2.5 Authentik

Evidence is repo-only (no admin token was used).

- Deployed in `auth` (2 server replicas, worker, Postgres, Redis). Public at `auth.onemindos.dev` (Cloudflare
  CIDR allowlist on the ingress) and on the tailnet as `onemind-auth`.
- OIDC applications bootstrapped by Sync-hook Jobs (no blueprints in the repo): Backstage (`onemind-backstage`),
  Argo CD (`argocd`, Dex disabled in PR #392), and OpenChoreo, workflows, Langfuse and registry clients
  (`platform/ecosystem/resources/auth-provider-bootstrap.yaml`). TAK Fedhub also references Authentik.
- Groups: a `groups` scope mapping, an `argocd-admins` group, and `onemind-admins` (Backstage notifications).
  The owner account `akadmin` is bound to the app policy bindings. Everything else is per-user bindings.
- MFA: **no MFA policy, authenticator stage or flow is defined in the repo.** Whether the live default flows
  enforce it cannot be told from the repo. Treat as unknown, to be checked in the Authentik admin UI.
- Service accounts: an AI-operator service account is admitted to the OpenChoreo app (policy mode `any`).
- Bootstrap Jobs authenticate to the Authentik API with an admin token Secret, so that token is a high-value
  secret.

### 2.6 OpenFGA

- Deployed: chart 0.3.16, app v1.22.0, one replica in `onemind-core`, Postgres datastore, **preshared-key auth**
  (key in a sealed secret), playground off. Reachable only from `onemind-core` and `backstage` pods
  (`openfga-clients` NetworkPolicy).
- Used by one caller: the Backstage Missions backend plugin, which calls `/stores/{store}/check` with a fixed
  store id and model id mounted from a Secret. Live logs show Write, Read and ListStores calls. No agent, Hermes,
  NATS or action-layer code calls it. `onemind-core` has a `openfga-auth` Secret but its use beyond Missions is
  not wired.
- The authorization model itself is created by a bootstrap step outside the files reviewed here; its contents
  are not in this document. Step 4 below starts by exporting it into Git.

### 2.7 Public exposure

Source: live Ingress and Service objects, 2026-10-23. The Cloudflare tunnel (`cloudflared`, `cloudflared-tunnel`)
is scaled to 0, so public HTTP goes Cloudflare proxy to the DigitalOcean LB `143.244.220.155` (ingress-nginx).
Ingresses have no TLS block, so the Cloudflare-to-origin hop is HTTP on port 80.

Public (internet reachable):

| Endpoint | What | Restriction | Intended per `docs/networking/NETWORKING.md` |
|---|---|---|---|
| `auth.onemindos.dev` | Authentik | Cloudflare CIDR allowlist | Public |
| `ws.onemindos.dev` | NATS WebSocket (8080, no TLS on the server) | Cloudflare CIDR allowlist | Tailnet only |
| `onemindos.dev` | CloudTAK | Cloudflare CIDR allowlist | Tailnet only |
| `tileserver.onemindos.dev` | Geo tiles | Cloudflare CIDR allowlist | Tailnet only |
| `matrix.onemindos.dev`, `element.onemindos.dev`, `onemindos.dev/.well-known` | Synapse, Element | none | Public |
| `enroll.onemindos.dev` | TAK enrollment portal | **none** | Tailnet only |
| `fedhub.onemindos.dev` | TAK Federation Hub admin UI | **none** | Tailnet only |
| LB `152.42.155.95` (`tak-external`) | TAK 8089, 8443 (WebTAK), 8446 (enrollment), 9101 and 9102 (federation), 64738 (Mumble) | no source ranges; no NetworkPolicy in `takserver` beyond backup egress | 8089/8443/8446/64738 to move to the tailnet `tak` node; 9101/9102 open decision |

The Cloudflare CIDR allowlist filters by source IP, so it only works while the LB keeps the client source IP
(`externalTrafficPolicy: Local` is set on the ingress LB). Anyone going through Cloudflare itself can still
reach those hosts, so it is not authentication.

Tailnet only: everything in 2.4. LoadBalancer Services: `ingress-nginx-controller` (public), `tak-external`
(public), `omos-nats`, `tak-tailnet` and two Cilium gateway Services (tailnet). No Tailscale Funnel anywhere.

## 3. Principles per layer

| Layer | Principle | Current state | Gap |
|---|---|---|---|
| People | Authentik is the one identity provider. MFA (WebAuthn first) for every person; group-based access; no shared admin | SSO for Argo CD, Backstage and four platform apps; MFA unknown; only one human; many admin UIs are tailnet-gated only | G6, G10 |
| Devices | Only enrolled, tagged devices reach private services; least-privilege grants by tag | Tailnet policy in Git, tests, no Funnel. Admin group reaches everything. One shared `tag:k8s` | G5 |
| Services | Default deny; explicit allow; encrypt in transit; no implicit trust between namespaces | 2 of ~50 namespaces default-deny; no WireGuard (not available on DOKS); NATS TLS, Dapr mTLS | G2, G8 |
| Agents | Every agent has its own NATS identity limited to its own subjects; observers are read-only; identity is revocable | One shared `onemind` login and one shared `leaf_tak`; no permissions anywhere | **G1** |
| AI actions | Separate observer, requester, approver and executor; approvals bound to exact action; every decision audited | PR #339 draft, offline contract tests only. OpenFGA deployed but used only by Missions | G3 |
| Workloads and supply chain | Admission rules enforced; images pinned by digest and signed; least-privilege pods | Kyverno installed with 0 policies (audit set proposed in #395); images on GHCR, many digest-pinned; no signing | G4, G7 |
| Secrets | Secrets sealed in Git or synced from a store; narrowly scoped; rotated | Sealed Secrets; one shared NATS password in ~10 Secrets; static bcrypt in the live ConfigMap | G1, G9 |

## 4. Gaps, prioritized

| ID | Priority | Gap | Why it matters |
|---|---|---|---|
| G1 | P0 | One `onemind` NATS login for every agent, feed and tool; `leaf_tak` shared by Mac, Hermes and TAK; no permissions | Any compromised pod or prompt-injected agent can publish as Legacy, read every subject, and administer JetStream |
| G0 | P0 | `fedhub.onemindos.dev` and `enroll.onemindos.dev` public with no IP restriction; TAK LB ports public with no policy | Admin and enrollment surfaces are internet-reachable; docs say they should be tailnet only |
| G2 | P1 | Default deny in 2 of ~50 namespaces; `nats`, `hermes`, `omos`, `takserver` open | Lateral movement is unconstrained inside the cluster |
| G3 | P1 | OpenFGA not wired into the action layer; no model in Git | Approvals for AI actions have no policy engine |
| G4 | P1 | Kyverno has no policies; #395 is audit-only | No admission guardrails |
| G5 | P2 | Tailscale: one tag for all proxies; admin group has `*`; no device posture | Blast radius of one proxy or one device |
| G6 | P2 | MFA and group model not codified for Authentik | Depends on UI state not in Git |
| G7 | P2 | No image signing or verification | Anything pulled from GHCR is trusted |
| G8 | P3 | No pod-to-pod encryption (WireGuard unavailable) | Plaintext inside the node network except NATS and Dapr |
| G9 | P3 | Websocket on NATS without TLS; Cloudflare to origin over HTTP | In-flight credentials exposure on those paths |
| G10 | P3 | Hermes webhooks have no designated public path | Needed later; Funnel is the planned way |

## 5. Implementation plan

Small PRs, in order. Each step lists what changes, how it is verified, the risk, the rollback, and how it fits the
cloud plus the future home cluster. Steps that touch the live cluster are applied by Zeus through Argo or a
reviewed manual action; this plan does not change anything by itself.

### Step 0 (P0, quick): close the public admin surfaces

- Change: add the Cloudflare CIDR allowlist annotation (as on `auth` and `ws`) to `fedhub-ingress` and
  `tak-portal-ingress`, as an interim; then follow the existing cut-over runbook in `docs/networking/NETWORKING.md`
  to move them to tailnet-only. Add a `takserver` default-deny with allows for the TAK ports (see Step 2).
  Decide 9101/9102 (federation peers).
- Verify: from outside Cloudflare ranges, requests to the LB with those Host headers return 403; tailnet access
  to `tak` still works for EUDs.
- Risk: an enrolled device that is not on the tailnet loses enrollment. Rollback: revert the annotation PR.
- Cloud/home: the same allowlist idea becomes the Tailscale Funnel rule (Step 6) once a home cluster exists.

### Step 1: per-agent NATS credentials (top priority)

Design:

- Identities (one per principal, never shared): `agent_legacy` (Hermes cloud), `agent_legacy_edge` (Mac), the four
  overseers `agent_so/hp/le/ge`, `fabric_tap` (read only), `fabric_console`, `leaf_mac`, `leaf_hermes`,
  `leaf_tak`, `feed_*` (one per feed group), `sink_clickhouse`, `backup`, `core` (onemind-core), `nack`, and a
  break-glass `admin` kept offline.
- Credential type: **nkeys**. Public keys go in the NATS config in Git; the seed (private key) lives only in the
  owner's SealedSecret. A leaked config reveals nothing; a rotation is a new keypair.
- Permissions are explicit allow lists. For a Hermes agent with owner `zeus` and session `legacy`:
  - publish: `agents.hb.hermes.zeus.legacy`, `agents.status.hermes.zeus.legacy`,
    `agents.attachments.hermes.zeus.legacy.>`, `_INBOX.>`, plus the agent service reply subjects, plus only the
    subjects it truly needs (for example `tak.cot.inject` for the TAK gateway, `world.>` reads)
  - subscribe: `agents.prompt.hermes.zeus.legacy` (and its `.>` if used), `$SRV.PING`, `$SRV.INFO`, `$SRV.STATS`
    and the `.agents` variants for discovery, `_INBOX.>`
  - Legacy additionally publishes `agents.prompt.>` (the sovereign prompts others); overseers publish only to
    their own servants; servants publish to nobody else's prompt subject.
  - The user wrote `agents.<name>.>`; because the grammar is verb-first, the plan fixes the framework, owner and
    session tokens instead.
- Agent Fabric is read only: `fabric_tap` gets subscribe `agents.>` and `$SRV.>` responses it asked for, and
  publish `_INBOX.>` only. `fabric_console` adds publish on `agents.prompt.hermes.zeus.legacy` only (chat with
  Legacy). If chat must be strictly human-gated later, that publish moves behind an Authentik-checked path.
- JetStream: only `nack`, `backup` and `core` get `$JS.API.>` write subjects; everyone else gets the specific
  `$JS.API.CONSUMER.*` they need.
- Leafnodes: three separate leaf users instead of `leaf_tak`, each limited to the subjects its site owns
  (Mac: `agents.*.hermes.zeus.legacy-edge`, `edge.>`; Hermes cloud sidecar: its agent subjects; TAK bridge:
  `tak.>`). Confirm in a rehearsal that leaf authorization accepts nkeys on 2.14.2; if it does not, use long
  random passwords per leaf, hashed, which still gives per-site revocation.

Migration without downtime, in two phases:

1. **Phase A: nkey users with permissions in the existing config-mode account (no restart, no data move).**
   a. Inventory first. Run a permissive-observation window: add the new users with wide permissions and the old
      `onemind` user still present; use `nats-surveyor`/`connz` to list which subjects each client really uses
      (the cluster already exposes `connz` and `subsz` on 8222). Write the allow lists from observed subjects.
   b. Land the new users in a Git-managed ConfigMap (finish the Helm-to-Argo move the OPERATIONS note says is
      unfinished) and reload with a config reload (`nats-server --signal reload`, or `kubectl exec ... -- nats-server
      -sl reload`), one server at a time. A reload applies new users without dropping existing connections.
   c. Move clients one at a time from `onemind` to their own nkey: update the SealedSecret, roll the workload,
      confirm the connection appears under the new name in `connz`. Order: Fabric (read-only), feeds, sinks,
      backup, takml/mumble, Hermes legacy and its leaf, Mac leaf, TAK, last the Core services.
   d. When `connz` shows no `onemind` or `leaf_tak` connections for 7 days, remove both users and reload. Their
      connections drop at that moment, so only do this when nobody is using them.
2. **Phase B (optional, decide after A): operator/account/user JWT with the built-in full resolver.** This is the
   decentralized model asked for. Two caveats the Doctrine does not mention, to confirm in a rehearsal cluster
   before committing:
   - A server runs in operator mode or with config-defined accounts, not both, so Phase B is a server-wide
     switch, and a cluster with mixed auth modes is not a supported intermediate state.
   - In operator mode the account is named by its public key, and JetStream stores data under the account name,
     so the existing `ONEMIND` streams would not be found under the new account. The safe path is a new account
     on a second NATS StatefulSet, with streams restored from the encrypted backup (the restore of all 39
     streams is already tested) or sourced across accounts, then clients switched. Plan for a cut-over window
     with dual publish for the streams that matter (COT, intel).
   - The alternative that keeps the ONEMIND account and its JetStream intact is auth callout (the Doctrine's
     Phase 7), which can mint per-client permissions from Authentik or OpenFGA. Choose B or callout after A.
- Verify: for each identity, a negative test with `nats pub` and `nats sub` on a forbidden subject returns a
  permissions violation; the agent still appears in `nats micro ls`; heartbeats still flow; streams unaffected
  (`nats stream ls` equal before and after).
- Risk: an allow list that is too tight silently breaks an agent (it keeps running but drops messages). Mitigate
  with the observation window and by adding `allow_responses` and `_INBOX.>` where request-reply is used.
  Config reload with a syntax error: the server keeps the old config; validate with `nats-server -t` in CI.
- Rollback: re-add the `onemind` user (kept until step 1d) and point the SealedSecret back; one reload.
- Cloud/home: leaf users are the home-cluster story. A home NATS joins as a leaf with its own nkey, limited to
  the subjects the home site owns, over the Tailscale `omos-nats` address. The home cluster's own agents use
  local users and never hold cloud credentials.

### Step 2: default-deny NetworkPolicies, namespace by namespace

One PR per namespace group, smallest and quietest first. Every PR adds `default-deny` (ingress and egress) plus
named allow policies in the same file, in `platform/` next to the namespace.

Allow rules every namespace needs: DNS to `kube-system/kube-dns` (UDP and TCP 53); ingress from the
`tailscale` namespace proxies for Ingresses; Prometheus scrape from `monitoring`; egress to the apiserver where
a pod uses the API (Cilium `toEntities: kube-apiserver`); webhook-serving pods need ingress from
`kube-apiserver` and `remote-node` entities (`CiliumNetworkPolicy`), or admission calls fail.

Order:
1. Already done: `agent-fabric`, `n8n`.
2. Quiet and low blast radius: `onemind-jobs`, `onemind-media`, `opencost`, `velero`, `external-dns`.
3. Platform tools: `monitoring`, `argocd`, `backstage`, `auth` (Authentik needs Postgres, Redis, ingress-nginx).
4. Data: `clickhouse`, `onemind-data`, `onemind-search`, `onemind-registry`, `onemind-ai`, `onemind-core`
   (keep Dapr sidecar ports 50001/50002/3500, allow `dapr-system`).
5. `hermes` and `omos`: egress only to NATS 4222, DNS, the model provider on 443, tailnet egress proxies in
   `tailscale`, and DigitalOcean inference. This is where agent containment becomes real.
6. `nats`: ingress 4222 from named client namespaces, 7422 from the tailnet proxy and cluster, 6222 only between
   NATS pods, 8222 only from monitoring and the NATS UIs. Last of the sensitive ones because a mistake breaks
   the whole bus.
7. `takserver` and `onetak`: ingress from the public LB path and tailscale; Postgres only from TAK.
8. Not default-deny, use targeted policies and entity rules: `kube-system`, `cert-manager`, `kyverno`,
   `external-secrets`, `cnpg-system`, `ingress-nginx`, `tailscale`, `falco`, `openchoreo-*`. Locking
   `kube-system` carelessly can take down DNS and the CNI.
- Method: before enforcing, observe with Hubble (`hubble observe --namespace X --verdict DROPPED` after a
  test policy in a canary, and `--to-fqdn` for egress) and write allow rules from the flows. Cilium policy
  verdicts are already enabled in the config.
- Verify: per PR, a smoke list (pod Ready, Argo Healthy, an expected request works, an unexpected one is
  denied from a debug pod); Hubble shows no unexpected DROPPED verdicts for 24 hours.
- Risk: broken egress to an unlisted dependency. Rollback: delete the `default-deny` object (Argo revert); the
  allow rules are harmless alone.
- Cloud/home: write rules by label and namespace, not by DigitalOcean IP. The same manifests then apply to a
  home Cilium cluster. FQDN egress rules need Cilium (also on the home cluster).

### Step 3: Tailscale ACL review

- Change (policy.hujson): split `tag:k8s` into `tag:k8s-ingress` (proxies that expose services) and
  `tag:k8s-egress` (the ProxyGroup, Hermes reach to devices); update `PROXY_TAGS`/ProxyClass tags and
  `tagOwners`; scope `group:admins` to the ports actually used instead of `*`; add device posture if the plan
  allows; add `tests` for the new deny cases; consider a `tag:agent` for Hermes-managed edge devices and a
  `tag:home` for the home cluster; keep `ssh` rule narrow.
- Verify: the `tailscale.yml` workflow diff and its `tests` block on the PR; manual check that Argo, Grafana and
  Agent Console still load from the admin laptop and the phone.
- Risk: lock-out of the admin from a UI. Mitigation: keep a break-glass grant for one device. Rollback: revert
  the PR; the workflow reapplies the previous policy on merge.
- Cloud/home: tags are the portability layer. Home nodes join with a tag and get only the NATS leaf port
  (7422) and the services they need.

### Step 4: wire OpenFGA into the #339 approval layer

- Change, in small PRs against #339's contracts (`docs/roadmap/permissioned-ai-actions.md` on the #339 branch):
  1. Export the current store and authorization model into Git (`platform/ecosystem/openfga/model.fga`) and
     load it with a reviewed job; add model tests (`fga model test`) to CI.
  2. Model: `user`, `agent`, `action_type`, `target` (app, namespace), with relations `requester`, `approver`,
     `executor`, `observer`; approvals as short-lived tuples bound to the exact head SHA and expiry.
  3. The action layer calls `check` before executing and writes the audit tuple; deny by default.
  4. Replace the preshared key with a per-caller key (Authentik service accounts later), and grant the call
     path only to the action layer and Backstage through the existing `openfga-clients` policy.
- Verify: contract tests in #339 gain allow and deny fixtures; a staged run where the check result is logged
  but not enforced (shadow mode) for a week.
- Risk: OpenFGA outage blocks actions; this is the intended fail-closed behavior for the executor but must
  never block the read-only observer. Rollback: toggle shadow mode.
- Cloud/home: OpenFGA is stateless apart from its Postgres; the home cluster can run its own store seeded
  from the same model file.

### Step 5: Kyverno from audit to enforce

- Prerequisite: merge #395 (audit set), read PolicyReports for a week, fix or exempt offenders.
- Change, per policy: flip to `Enforce` one at a time, cheapest first: disallow privileged/hostPath (one known
  offender, `github-runner`, exempt it explicitly), then `latest` tags (17 known workloads must be pinned first),
  then requests/limits, then layer label (decide the `security` value first). Raise Kyverno admission to 3
  replicas and keep `failurePolicy: Ignore` until then, because it currently runs 1 replica and a failure under
  `Fail` would block pod creation cluster-wide.
- Add `ValidatingPolicy` (CEL) versions as Kyverno 1.19 deprecates `ClusterPolicy`.
- Verify: `kyverno apply` against exported live manifests; a canary namespace first.
- Risk: blocks rollouts. Rollback: set `validationFailureAction` back to Audit (one-line PR).
- Cloud/home: policies are cluster-agnostic; the home cluster installs the same Kyverno chart and policies.

### Step 6: public endpoints through Tailscale Funnel only

- Change: for each endpoint that must be public and has no better home (Hermes webhooks first), publish it as
  a Tailscale Ingress with `tailscale.com/funnel: "true"`, and in `policy.hujson` add `nodeAttrs` granting
  `funnel` to the specific tag only. The Funnel proxy forwards to a minimal receiver that verifies the sender's
  signature and publishes to a single NATS subject using its own limited nkey. Then remove the matching
  nginx Ingress and Cloudflare record. Matrix and Authentik stay on nginx/Cloudflare until they have a home.
  Decide whether Authentik keeps a public path or moves behind Funnel as well.
- Verify: the endpoint answers from the internet; the policy tests assert only the intended tag has `funnel`;
  Hubble shows the proxy talking only to its backend.
- Risk: Funnel hostnames are on `*.ts.net`, with limited ports (443, 8443, 10000) and bandwidth limits.
  Rollback: remove the annotation and the `nodeAttrs` grant.
- Cloud/home: Funnel works from any tailnet node, so the same receiver can run on the home cluster with no
  inbound port forwarding or public IP.

### Step 7: image signing and verification

- Change: in the GHCR build workflows (`build-*.yml`) add keyless cosign signing (GitHub OIDC) of the digest the
  build already records; add SBOM attestation if cheap. In Kyverno add an `ImageValidatingPolicy` for
  `ghcr.io/onemindos/*` requiring a signature from the OneMind repository workflows, in Audit first. Third-party
  images (Tailscale, NATS, Cilium, Authentik) are verified by digest pin until their publishers' signatures can be
  checked.
- Verify: `cosign verify` in CI; Kyverno reports show signed/unsigned counts; then Enforce for the namespaces
  whose images are all first-party and signed.
- Risk: a rebuild that skips signing blocks a deploy. Rollback: policy back to Audit. Do not require signatures for
  the Argo CD image-update path until the workflow signs every build.
- Cloud/home: the policy trusts a signer identity, not a registry host, so a home registry mirror still works
  provided it preserves signatures.

### Step 8: smaller follow-ups

- Codify Authentik: blueprints (or the bootstrap Jobs) for WebAuthn-required flows, groups
  `onemind-admins/operators/observers`, and an audit export; check MFA state in the UI now.
- Turn off the NATS server WebSocket `no_tls` path in favor of the tailnet `nats-ws` Ingress (TLS terminated at
  the proxy) and delete the public `ws.onemindos.dev` ingress.
- Encrypted transport inside the cluster: since WireGuard is not available, add TLS where plaintext remains
  (Postgres `sslmode=verify-full` for CNPG clients, Redis TLS for Authentik) rather than replacing the CNI.
- Remove the SPIRE objects and README text when the removal PR lands.

## 6. Order and dependencies

| # | PR | Depends on | Size |
|---|---|---|---|
| 0 | Allowlist or retire public fedhub and enroll | none | XS |
| 1a | NATS observation window and nkey users, permissive | none | M |
| 1b | Move clients to own nkeys, Fabric first | 1a | M, several PRs |
| 1c | Remove `onemind` and `leaf_tak` | 1b | S |
| 2 | Default-deny, namespace by namespace | 0 for takserver | several S PRs |
| 3 | Tailscale tags and grants | none | S |
| 4 | OpenFGA model in Git, then wire to #339 | #339, step 1 for agent identity | M |
| 5 | Kyverno enforce | #395 | S per policy |
| 6 | Funnel for webhooks | 1 (receiver identity), 3 | M |
| 7 | cosign plus verifyImages | 5 | M |

Steps 0, 1a, 3 and 5's prerequisite can run in parallel.

## 7. Verification principles

- Every step has a negative test (the forbidden thing fails), not only a positive one.
- Hubble and `connz` are the evidence for network and NATS changes; PolicyReports for Kyverno.
- No production change without a rollback that is one reviewed Git change or one reload.
- NATS backups (the encrypted nightly stream backup) are checked before steps 1c and any Phase B work.

## 8. Related documents

- `docs/NATS-DOCTRINE.md` (security model section)
- `docs/networking/NETWORKING.md` (public endpoint inventory and cut-over)
- `docs/roadmap/permissioned-ai-actions.md` (added by draft PR #339, not on main yet)
- `platform/hermes/agent-fabric/README.md`
- `platform/security/README.md`
- `platform/core/policy.hujson`
- `docs/operations/home-migration.md`
