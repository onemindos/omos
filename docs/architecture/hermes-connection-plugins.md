# Legacy connection layer: Hermes plugins for NATS, TAK and the rest

Status: **proposed design, nothing deployed.** 2026-10-30. Builds on the open PRs #560, #562, #564, #551 and #559 and on the merged plugin catalog (#582); section 12 gives a verdict for each. Owner direction: no NATS sidecar if
avoidable; connect through proper Hermes plugins; TAK-native; plugins are "the connection system"
for TAK, NATS, Home Assistant, VR-TAK, Obsidian and later systems.

Scope: the agent "Legacy" running as `onemind-hermes` (Hermes Agent v0.22.5, namespace `hermes`).
Related: [`platform/hermes/onemind-hermes/MIGRATION.md`](../../platform/hermes/onemind-hermes/MIGRATION.md)
(step 7 "NATS cutover" is replaced by this design), [`docs/NATS-DOCTRINE.md`](../NATS-DOCTRINE.md),
`CLAUDE.md` (Bus Authority, subject grammar, TAK Server), and ADR 0003 and the TAK mesh gateway on
branch `feat/tak-mesh-gateway` (commit `977b7eb2`, not on `main`).

## 1. Decisions in one screen

1. **Connections are Hermes plugins, built from one Git monorepo, shipped in a digest-pinned derived
   image, and copied into `$HERMES_HOME/plugins` by an init container at pod start.** No `patch_*.py`
   hot-patch, no hand-installed plugin on the volume. Plugin enablement and isolation are pinned by a
   read-only managed config (`/etc/hermes/config.yaml`), not by the copy-once `config.defaults.yaml`.
2. **NATS: direct connection, no sidecar.** The plugin connects to
   `nats://nats.nats.svc.cluster.local:4222` as its own user `hermes` in account `AGENTS` (PR #562, connection type widened from leaf-only), with
   explicit publish/subscribe permissions and its own inbox prefix (`_INBOX.hermes`). The Legacy leaf sidecar, `hermes-leaf-conf`
   and the shared `leaf_tak` password are retired at cutover. This also closes a real gap: the sidecar,
   the TAK leaf and the Mac leaf all authenticate as the same full-permission `leaf_tak` user today.
3. **TAK: Legacy is a TAK client first, a bus subscriber second.** A `tak` gateway-platform plugin holds
   a mutual-TLS streaming connection to `takserver-messaging:8089` with its own client certificate
   (callsign `Legacy`). GeoChat becomes conversation, PLI is presence, CoT it sends is attributed to its
   certificate and bounded by its TAK groups. The common operating picture is read from the bus
   (`tak.cot.>`), not by reading the whole firehose a second time.
4. **The TAK bridge to the bus is broken today and has to be fixed regardless** (section 5.1): the
   `takserver-plugins` pod has no plugin JAR, and even with it TAK Server in cluster mode never hands
   CoT to a plugin, so `tak.cot.out` cannot be fed by the NatsBridgePlugin alone. Fix: merge PR #560
   (the Go `tak-client` streaming bridge carries both legs and the plugin source is repaired) instead of
   resurrecting the Java plugin plus `cot-tap`.
5. **Plugin types by job:** gateway platform for anything that is a conversation or an inbound event
   source (NATS agent protocol, TAK, Matrix/Telegram stay bundled); standalone plugin with tools and
   hooks for request/response systems (bus tools, Marti API, Home Assistant, Obsidian); approval
   transport plus `pre_tool_call` policy hook for the safety layer. Memory-provider plugins are not
   used in phase 1.
6. **Every plugin ships a `plugin.yaml`** that is a normal Hermes manifest plus an `onemind:` block
   (id, type, version, owner, systems, subjects, endpoints, credential Secret names, runtime
   location). Hermes ignores unknown manifest fields, so one file serves the loader and Backstage
   (section 6).
7. **Physical actions are not granted to Legacy at the NATS layer.** The `hermes` user cannot publish
   `cmd.drone.>`, `cmd.robot.>`, `cmd.fleet.>`, `cmd.mission.>` or `cmd.intercept.>`. Legacy proposes
   through the policy service (#554, #564); a human confirms from a CLIENT identity; the policy service
   publishes (section 8). `execute` and `override` stay human-typed words.
8. **Credentials are SealedSecrets** (`onemind-hermes-nats`, `onemind-hermes-tak`, ...). Because the
   terminal backend is `local`, anything mounted in the pod is readable by the agent's shell; the
   design therefore relies on server-side NATS permissions and TAK group ACLs, not on hiding secrets.
9. **Phased:** (0) fix TAK bus bridge, (1) image + plugin pipeline with a no-op plugin, (2) NATS agent
   service + guard + registry, (3) TAK platform, (4) tool plugins (Marti, Home Assistant, Obsidian),
   (5) actuation through the policy service. Legacy (the old pod) keeps running until phase 3 passes its checks.

## 2. What is true today (verified read-only, 2026-10-29/30)

| Fact | Evidence |
|---|---|
| `onemind-hermes` runs Hermes v0.22.5, Python 3.14.9, as uid 0, writable root filesystem, **no `nats-py` installed**; `protobuf 6.33.6`, `aiohttp`, `httpx`, `websockets`, `cryptography`, `mcp 2.3.0` are present; `lxml` is not | `pip` metadata inside the pod |
| `/opt/data/plugins` is empty; `plugins.enabled: []`; no service-account token; egress only DNS, LiteLLM 4000, Synapse 8008, public 443 | pod filesystem; `platform/hermes/onemind-hermes/network.yaml` |
| Legacy (old pod) runs three plugins I could read: `nats-platform` (Synadia `hermes-nats-gateway` 0.1.2, 2,400 lines, patched at start by `nats-platform-patch-cm.yaml`), `nats-context` (tools + `pre_llm_call` picture), `platforms/tak-platform` (469 lines, subscribes `tak.cot.out`, injects via `tak.cot.inject`) | `hermes-legacy-*:/opt/data/plugins` |
| NATS auth is config based, one account `ONEMIND`, **one client user `onemind` for every in-cluster service**, plus `leaf_tak` for leaf nodes; no per-service permissions. TLS on 4222 | `kubectl -n nats get cm nats-config`; `connz` shows `authorized_user: onemind` for all 5 python/js clients on `nats-0` |
| The Legacy leaf sidecar (`hermes-legacy-leaf`) connected to `nats-2` as `leaf_tak` into `ONEMIND` (679 msgs in, 24,553 out since 01:52 UTC) | `leafz` on `nats-2`; `platform/hermes/agent/SECRETS.md` |
| `nats-clients` ingress policy on 4222 does not include namespace `hermes`; TAK's ingress policy already admits the whole `hermes` namespace to `takserver-messaging:8089` | `platform/nats/network-policy.yaml`, `platform/takserver/deploy/networkpolicy-takserver-ingress.yaml` |
| TAK plugin manager has no plugins (details in 5.1) | pod `takserver-plugins-*` |

Limits of this evidence: `connz`/`leafz` counters were read from individual servers (`nats-0` for
clients), the NATS account has no per-subject breakdown I could read without credentials, and nothing
below was tested by running a plugin.

## 3. Hermes v0.22.5 plugin surface

Read from the running image (`hermes_cli/plugins*.py`, `plugins/AGENTS.md`, `gateway/platforms/ADDING_A_PLATFORM.md`).

### 3.1 Plugin kinds

| Kind | How it registers | Loaded when | Fits |
|---|---|---|---|
| **Gateway platform** (`kind: platform`) | `ctx.register_platform(name, label, adapter_factory, check_fn, ...)`; adapter subclasses `BasePlatformAdapter` (`connect`, `disconnect`, `send`, `get_chat_info`, ...). Inbound becomes a `MessageEvent`; the gateway supplies authz (`*_ALLOWED_USERS`), sessions, approvals, cron delivery, `send_message` routing | Non-bundled: only if listed in `plugins.enabled`; adapter connects when the gateway starts | **NATS agent protocol** (prompts in, replies out), **TAK** (GeoChat in/out, presence) |
| **Standalone plugin** | `register(ctx)` calls `ctx.register_tool(...)`, `ctx.register_hook(...)`, `ctx.register_command(...)`, `ctx.register_cli_command(...)`, `ctx.register_system_prompt_section(...)`, `ctx.spawn_task(...)` (supervised background task), `ctx.on_unload(...)` | `plugins.enabled` | **Bus tools** (peek, request, KV, JetStream fetch), **Marti API**, **Home Assistant**, **Obsidian**, **guard/audit**, **registry publisher** |
| **Hooks** (`VALID_HOOKS`) | `pre_tool_call`, `post_tool_call`, `pre_llm_call`, `post_llm_call`, `on_session_start/end`, `pre_gateway_dispatch` (can skip/rewrite an inbound event), `post_gateway_admission`, `subagent_start/stop`, streaming observers, ... | with the owning plugin | policy (block a tool call), audit events, live-picture injection |
| **Approval transport** | `ctx.register_approval_transport(name, present_fn)`; selected by `security.approval.transport: <name>`; receives an immutable redacted `ApprovalRequest` and returns a request-bound decision; policy stays host-owned | with the owning plugin | Where humans approve flagged terminal commands (Telegram/Matrix today; a TAK or omos-ui transport later) |
| **Memory provider** (`kind: exclusive`) | `MemoryProvider` ABC, chosen by `memory.provider`; bundled first, then `$HERMES_HOME/plugins/`, then entry points | by config, not `plugins.enabled` | Not used now. Obsidian-as-memory is a possible later step; Legacy's `memory-wiki` is not carried over |
| Model provider, context engine, image/video/web, dashboard auth, cron providers | their own ABCs | own discovery | Not connection plugins |

Capabilities a plugin must declare and an operator must grant (`plugins.entries.<id>.granted_capabilities`):
`tools.override`, `llm.*_override`, `gateway.platform_actions`. Gateway injection
(`ctx.inject_message`) needs `plugins.entries.<id>.allow_gateway_injection: true`. Our connection plugins
need none of the `llm.*` or `tools.override` capabilities, and must not ask for them.

`plugins.isolation: host` runs every non-bundled plugin in a per-profile host process; plugins cannot
use live gateway/adapter objects there (`HOST_UNSUPPORTED_CTX_METHODS`), so platform plugins stay
`in_process`. Isolation is not a security boundary against the agent's shell (same uid, same pod).

### 3.2 Discovery sources and the enable gate

Sources, later wins: bundled `<repo>/plugins/`, user `$HERMES_HOME/plugins/<name>/` (needs `plugin.yaml`
and `__init__.py` with `register(ctx)`), project `./.hermes/plugins/` (opt-in env var), pip entry
points in group `hermes_agent.plugins`. Gate (`gate_manifest`): `plugins.disabled` wins; **non-bundled
plugins load only when named in `plugins.enabled`** (platforms bundled with Hermes are deferred-lazy
instead). The catalog installer (`hermes plugins install <name>`) needs a 40-hex SHA pin and writes
into the plugin directory at runtime; it is the wrong tool for a GitOps pod.

Compatibility contract: documented `ctx` surfaces are additive, hook callbacks are signature-inspected,
unknown manifest fields are ignored, `requires_hermes` can pin a range. This is what makes in-repo
plugins safe to maintain **without** the startup patch Legacy needs: the old `nats-platform` patch
exists only because upstream's adapter imported `set_current_session_key` from `tools.approval`, which
moved to `tools.approval_context`. We vendor the adapter and fix the import as a normal reviewed commit.

## 4. Packaging: versioned in Git, installed at pod start

### 4.1 Options compared

| Option | Versioned | Deps (`nats-py`) | At pod start | Verdict |
|---|---|---|---|---|
| A. Pip wheel + entry points in a derived image | yes (wheel version, image digest) | declared in `pyproject.toml` | none | Works, but entry-point plugins have no `plugin.yaml` for Hermes, so no `kind: platform`, no `requires_env` UI, no per-plugin manifest for Backstage |
| **B. Directory plugins baked into a derived image, copied to `$HERMES_HOME/plugins` by an init container** | yes (Git SHA in image label, image digest) | pip layer in the same image | init container syncs, no network | **Recommended.** Keeps the Hermes manifest, reproducible, no runtime install |
| C. ConfigMap-mounted plugin files | yes | cannot add wheels | kubelet | 1 MiB limit, no dependencies; fine only for a single-file plugin |
| D. `hermes plugins install <git url@sha>` at start | yes | pip at runtime | needs GitHub egress, writes PVC | Rejected: mutable PVC, egress, partial-failure start |

### 4.2 Layout

```
platform/hermes/agent/plugins/           # existing home of onemind-vision, and of onemind-bus in PR #564: keep it
  onemind-bus/      plugin.yaml  __init__.py  bus.py tools.py  # PR #564: entity/layer reads, proposals, map actions
  onemind-agent/    plugin.yaml  __init__.py  adapter.py       # NATS agent protocol (platform), vendored Synadia adapter
  onemind-tak/      plugin.yaml  __init__.py  adapter.py       # TAK streaming client (platform): chat + presence
  onemind-guard/    ...                                        # pre_tool_call policy + audit (standalone)
  onemind-marti/  onemind-homeassistant/  onemind-obsidian/    # phase 4
platform/hermes/onemind-hermes/image/    # new
  Dockerfile                             # FROM nousresearch/hermes-agent:v0.22.5@sha256:9c556157...
                                         # pip install nats-py; COPY platform/hermes/agent/plugins /opt/onemind/plugins
  requirements.txt                       # pinned, hashed
scripts/ci: contract tests run each plugin against the pinned Hermes image
```

CI builds `ghcr.io/onemindos/onemind-hermes:<hermes-version>-<git-sha>` on the existing GHCR path,
runs the contract tests inside the base image, and the Deployment pins the **digest**. Renovate bumps
the base image; a bump that breaks a `ctx` call fails CI before it reaches the cluster.

### 4.3 Pod wiring (additions to `platform/hermes/onemind-hermes/deployment.yaml`)

```yaml
      initContainers:
        - name: install-plugins                       # runs from the same derived image
          image: ghcr.io/onemindos/onemind-hermes@sha256:<digest>
          command: [/bin/sh, -c, "cp -a /opt/onemind/plugins/. /opt/data/plugins/"]
          securityContext: {allowPrivilegeEscalation: false, readOnlyRootFilesystem: true}
          resources: {requests: {cpu: 10m, memory: 32Mi}, limits: {cpu: 100m, memory: 64Mi}}
          volumeMounts: [{name: plugins, mountPath: /opt/data/plugins}]
      containers:
        - name: hermes
          image: ghcr.io/onemindos/onemind-hermes@sha256:<digest>
          volumeMounts:
            - {name: data,    mountPath: /opt/data}
            - {name: plugins, mountPath: /opt/data/plugins}   # emptyDir over the PVC dir: always exactly the image's plugins
            - {name: managed, mountPath: /etc/hermes, readOnly: true}
      volumes:
        - {name: plugins, emptyDir: {}}
        - {name: managed, configMap: {name: onemind-hermes-managed}}
```

`onemind-hermes-managed` carries `config.yaml` for Hermes' managed scope (`hermes_cli/managed_scope.py`,
directory `/etc/hermes`, immutable to the profile):

```yaml
plugins:
  enabled: [onemind-agent, onemind-bus, onemind-guard, onemind-tak]   # list grows by PR
  isolation: in_process
```

This replaces "edit the live `/opt/data/config.yaml` as well" for plugin policy. **Verify in phase 1**
that the managed file overrides a profile value and how it merges lists (read in the image:
`load_config_readonly`); if lists replace, keep the full list in the managed file.

Not changed: `automountServiceAccountToken: false`, `approvals.mode: manual`, `redact_secrets: true`.

## 5. TAK-native

### 5.1 Finding: the NatsBridgePlugin is not running, and cannot carry TAK to NATS

Read-only evidence, pod `takserver-plugins-6fb56b48c7-kbrcp`:

- Image is `ghcr.io/onemindos/takserver-base@sha256:f4824fba65a2...`. `/opt/tak/lib` is **empty**,
  `/opt/tak/conf/plugins` **does not exist**, no file named `*nats*` anywhere in the filesystem.
- `PluginLoader` logs `starting PluginLoader`, then nothing; no `NatsBridgePlugin config` line (the
  plugin logs one on construction).
- The process classpath: `JDK_JAVA_OPTIONS=-Dloader.path=WEB-INF/lib-provided,WEB-INF/lib,WEB-INF/classes`
  (no `lib/`). So even a JAR dropped into `/opt/tak/lib` would not be loaded.
- Cause in Git: `platform/takserver/values-onemind.yaml` and
  `platform/takserver/deploy/takserver-plugins-deployment.yaml` use tag/digest of
  `takserver-base:plugins-upstream-5.8-onemind` (upstream, no plugin). The workflow
  `.github/workflows/build-takserver.yml` builds the overlay from
  `services/takserver/docker/Dockerfile.takserver-plugins.onemind` (adds `lib/` to `loader.path`,
  copies the JAR and `conf/plugins/*.yaml`) and tags it `takserver-base:plugins-5.8-onemind`, but
  nothing deploys that tag. The values comment "Plugin is baked into this image" is wrong.
- Bus side agrees: `tak-nats-leaf` (the plugin's `natsUrl`) shows `in_msgs 0 / out_msgs 0` and no
  interest on its hub connection; nothing on the TAK side subscribes `tak.cot.inject`.
- The `COT` stream is healthy and growing (9,512,792 messages, last at 2026-10-30T03:39:39Z), but the
  writers are OneMind feeds: `feed-satellites` (publishes every 30 s) and `omos-tak-projector`
  publish **`tak.cot.inject`**, which nobody consumes, so the map clients never see it. No
  `tak.cot.out` producer exists. I could not read per-subject counts without credentials.
- Upstream limitation already recorded in `services/takserver/cot-tap/README.md`: in cluster mode TAK's
  `PluginClusterMessenger` is send-only. A plugin can inject (bus to TAK) but never receives the CoT
  firehose (TAK to bus), so the JAR alone cannot produce `tak.cot.out`. `cot-tap` was written for that
  and its Secret `cot-tap-certs` exists (116 days old), but there is **no `cot-tap` Deployment** in
  the cluster and it is not an Argo Application.
- The JAR cannot be rebuilt from the repo: `plugin-src/build.gradle` compiles against
  `../tak-server-sdk`, which is not in Git. The committed JAR (1.0.0, 2026-10-17) is the only artifact.
- Consequences: Legacy's old `tak-platform` adapter (listens on `tak.cot.out`) had nothing to hear;
  the portal catalog entry and `NATS-DOCTRINE.md` section 9.8 describe the bridge as working.

### 5.2 Fix (Phase 0): PR #560 is the route

PR #560 reached the same diagnosis independently (JDK options line, no plugin log line, no `cot-tap`
Deployment, overlay never deployed) and fixes it differently from the first draft of this document:

1. A Go `tak-client` bridge (`services/bridge/bridges/tak-client`, Argo `platform-takserver-tak-client`,
   manual sync) connects to `takserver-messaging:8089` with mutual TLS, negotiates TAK Protocol v1
   (`t-x-takp-v`/`q`/`r`) with XML fallback, publishes `tak.cot.out` with a JetStream ack and
   `Nats-Msg-Id`, and subscribes `tak.cot.inject` in queue group `tak-cot-inject` only while TAK is
   connected, with an echo filter. It carries **both legs**, which removes the cluster-mode
   send-only limit and the need for `cot-tap` (marked deprecated there).
2. `NatsBridgePlugin` 1.1.0 source is repaired (async connect, shared queue group, validation, status
   line) and the overlay Dockerfile now fails the build if the `loader.path` patch or the JAR is missing.
   Its JAR is not rebuilt (TAK SDK needs tak.gov), so the committed JAR stays 1.0.0.
3. Identity: `tak-client` in #562 (serves `tak.cot.inject`); only the projectors publish `tak.cot.inject`.

Under this design: **merge #560 and run `tak-client` as the single bus-to-TAK path; do not redeploy the
Java plugin or `cot-tap`.** If the plugin stays deployed at all, it must share the queue group (it does
in #560) so a message is injected once. `tak-nats-leaf` (plaintext, unauthenticated hop into the bus
with `leaf_tak`) can then be removed. Owner actions from #560: pin the image digest after the first CI
build, rotate the CA key passphrase and TAK user password that `cot-tap`'s README had in plain text
(A21-P4, still in Git history), and issue the `tak-client` certificate (user registered in the cluster DB).

Phase 0 checks: `tak-client` `/readyz` true; a CoT published to `tak.cot.inject` appears on a connected
EUD; a CoT sent by an EUD appears on `tak.cot.out` and the `COT` stream grows by it; the stale claims
(values comment, portal catalog text, NATS-DOCTRINE 9.8 "works") are corrected as #560 does.

### 5.3 How Legacy sees and sends CoT

| Option | Receive | Send | Identity and control | Verdict |
|---|---|---|---|---|
| **NATS** (`tak.cot.>` / `tak.cot.inject`) | Whatever the bridge publishes (all groups the bridge is in); history from stream `COT` | Anonymous inject to `__ANON__` groups | Plugin identity, not Legacy's; no per-contact chat | Right for the **picture**; wrong for conversation |
| **TAK streaming client** (8089 mTLS, TAK Protocol XML v0, v1 protobuf when negotiated) | Everything in the client's groups, including GeoChat addressed to it | CoT as `Legacy` with its own UID/callsign, GeoChat replies, PLI | Cert CN maps to a TAK user; groups restrict it; revocable; visible as a contact | **Recommended for conversation, presence and map writes** |
| **Marti API** (`takserver-api:8443`, client cert) | Missions, data packages, contacts, video | Create/alter missions, upload packages | Cert-based, REST | Tools for structured operations (phase 4), not a chat channel |

Design: `onemind-tak` (platform) opens the 8089 connection with a cert from Secret
`onemind-hermes-tak`, sends a PLI/presence CoT every 60 s, parses inbound `b-t-f` GeoChat addressed to
the Legacy UID (or its chat rooms) into `MessageEvent`s, and returns replies as GeoChat. It ignores the
rest of the firehose. It never writes to the bus. The COP for reasoning comes from `onemind-bus`
(`tak.cot.>` subscription feeding a bounded live picture). The direct connection carries **chat and presence only**. Map writes (markers, routes, overlays) are the
tool `tak_send_cot`, which calls `svc.cot.send` (the policy-gated projector of PR #563, as #562 assigns to
agents) so geofence and tasking rules apply and the bus stays the authority; Legacy does not write
map CoT through its own certificate. A GeoChat-to-`agents.prompt.hermes.>` bridge inside `tak-client` is
the alternative that removes the direct connection entirely (owner decision 6); the direct connection is
recommended because it makes Legacy a real contact with presence and per-group control.

Protocol notes: start with XML v0 (all servers accept it); add TAK Protocol v1 negotiation (`t-x-takp-v` /
`t-x-takp-r`, framing `0xbf 0x01 0xbf`) once the codec from the mesh gateway is vendored; `protobuf 6.33`
is already in the image, `lxml` is not (use `xml.etree`). Transport is in-cluster `takserver-messaging:8089`
(no public listener exists; TAK is tailnet-only for EUDs). Reconnect with jitter; one connection per
process; the cert's user must be registered in the cluster DB as `cot-tap` was.

VR-TAK: the OneMind VR client is a Commo/TAK client; it needs no Legacy plugin. Legacy sees VR users the
same way as any contact (PLI/chat), and a NATS panel in the XR client (WebSocket `wss://ws.onemindos.dev`
or tailnet) talks to Legacy through `agents.prompt`. A `onemind-vrtak` plugin is deferred until XR needs
something the TAK and bus plugins do not give.

## 6. NATS: direct, least-privilege

### 6.1 Direct vs sidecar leaf

| | Leaf sidecar (today) | **Direct client** |
|---|---|---|
| Identity | `leaf_tak`, shared with the Mac and the TAK leaf, full `ONEMIND` rights | Own user `legacy`, own permissions, revocable alone |
| Credentials | `hermes-leaf-conf` created out of band with the plaintext password | SealedSecret in Git |
| Hops / moving parts | Extra container, local 4222, TLS leaf link to 7422 | One TLS connection to the cluster service |
| Buffering | A plain leaf does not store-and-forward core pub/sub (no JetStream on it); gives no resilience to a hub outage | Same; reconnect handled by the client. Durable needs use JetStream consumers server side |
| Subject isolation | Leaf user rules would be needed per leaf | Native per-user permissions |
| Rationale for a leaf | Edge nodes behind NAT that must keep working offline | Hermes lives next to the hub; the argument does not apply |

Recommendation: direct. The leaf pattern stays for the Mac, edge hosts and the TAK mesh gateway, where
offline operation is the point. `CLAUDE.md` already says not to copy the Hermes sidecar pattern.

### 6.2 NATS identity: adopt PR #562's `hermes` user, change its connection type

PR #562 (per-service identities, A19) already defines the right shape: Hermes in the isolated account
`AGENTS`, no raw feeds, no `cmd.*` except `cmd.map.>`, entities through `svc.entity.*`, actions through
`svc.policy.{command,task}.propose`, never `*.confirm`, `agents.*.hermes.>` served, replies via
`allow_responses`, and a checker that enforces "only `policy` publishes `cmd.*`" and "only the
projectors publish `tak.cot.inject`". This design therefore does not give Hermes `tak.cot.>`, `ent.>`,
`evt.>`, `intel.>` or `det.>`: raw feeds are the prompt-injection surface #562 closes; Legacy reads entities and layers through services
(`onemind-bus`, PR #564) and TAK chat through its own TAK connection (section 5.3).

What this design changes in #562, and nothing else:

| #562 today | Change | Why |
|---|---|---|
| `hermes` is **leaf-only** (`allowed_connection_types: [LEAFNODE]`); `leaf-conf` renders the sidecar `leaf.conf` whose local user carries the allow-lists | Add `STANDARD` to the connection types of `hermes` (account `AGENTS`, same allow/deny lists). Drop `leaf-conf --user hermes` and the Hermes part of stage A2. Keep leaf users `leaf-tak` and `leaf-site-mac` | Owner direction: no sidecar. #562 enforces refusals at the sidecar because a hub only filters a leaf user; a direct client gets the server's own `Permissions Violation`, which is stronger and simpler |
| Identity `hermes`, served subjects `agents.*.hermes.>` | Keep. The Synadia agent token stays `HERMES_NATS_AGENT=hermes`, owner `zeus`, session name per MIGRATION step 7, so omos-ui discovery and the `AGENTS` exports to ONEMIND and CLIENT keep working. The persona name "Legacy" lives in the SOUL and in `agent-caps` keys | One subject family, no second identity |
| Replies on `_R_.>` | Client sets `inbox_prefix=_INBOX.hermes`; subscribe allow `_INBOX.hermes.>` | Without a per-user prefix a client could listen to other users' replies |
| `agents.prompt.hermes.>` only (A19-P6: no cross-account delegation) | No change | Delegation to other frameworks waits for them to join `AGENTS` |
| Liveness records | Add publish `$KV.agent-caps.hermes.>` and stream-info for `KV_agent-caps` (TTL 1 min bucket exists) | Plugin heartbeats for the registry (section 7.2) |

Not granted, on purpose (all consistent with #562): `tak.cot.inject`, `tak.cot.mesh.tx` (only the
projectors publish to TAK; Legacy's CoT goes through `svc.cot.send`), `cmd.drone/robot/fleet/mission/
intercept/service/nex/integration.>`, stream admin (`$JS.API.STREAM.{CREATE,UPDATE,DELETE,PURGE}.>`),
`tel.>`.

Server-side permissions are what survives a prompt-injected agent; plugin allowlists (#564's
`ALLOWED_REQUESTS`) are defense in depth. Activation, stages and rollbacks are #562's `IDENTITIES.md`
(config reload, never `helm upgrade`); it is an owner-approved step, not part of this document.

### 6.3 Secret, trust and network

```yaml
# platform/hermes/onemind-hermes/nats.sealedsecret.yaml (shape; values come from kubeseal --scope namespace-wide)
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  name: onemind-hermes-nats
  namespace: hermes
  annotations: {sealedsecrets.bitnami.com/namespace-wide: "true"}
spec:
  encryptedData: {NATS_USER: <sealed>, NATS_PASSWORD: <sealed>}
---
# CA public cert (not secret), same pattern as agent-fabric/nats-ca.yaml
apiVersion: v1
kind: ConfigMap
metadata: {name: nats-ca, namespace: hermes}
data: {ca.crt: "<copy of secret nats/nats-ca ca.crt>"}
```

Deployment env: `NATS_URL=tls://nats.nats.svc.cluster.local:4222`, `NATS_CA_FILE=/etc/nats-ca/ca.crt`,
`NATS_USER`/`NATS_PASSWORD` from `onemind-hermes-nats`, `NATS_INBOX_PREFIX=_INBOX.hermes`,
`HERMES_NATS_OWNER=zeus`, `HERMES_NATS_AGENT=hermes`, `HERMES_NATS_SESSION_NAME=<per MIGRATION step 7>`.
Other in-cluster Python clients connect with plain `nats://` and no CA bundle; the plugin should verify
against the CA (as `agent-fabric` does). Whether `nats-py` on Python 3.14 behaves is a phase 1 check.

NetworkPolicy additions (the server side must be changed in the same PR, otherwise the connection
times out silently, per `platform/nats/OPERATIONS.md`):

```yaml
# platform/hermes/onemind-hermes/network.yaml -> spec.egress (add)
    - to:
        - namespaceSelector: {matchLabels: {kubernetes.io/metadata.name: nats}}
          podSelector:
            matchLabels:
              app.kubernetes.io/name: nats
              app.kubernetes.io/component: nats
              app.kubernetes.io/instance: nats
      ports: [{protocol: TCP, port: 4222}]
    - to:                                       # TAK streaming (phase 3) and Marti API (phase 4)
        - namespaceSelector: {matchLabels: {kubernetes.io/metadata.name: takserver}}
          podSelector: {matchLabels: {app: takserver-messaging}}
      ports: [{protocol: TCP, port: 8089}]
    - to:
        - namespaceSelector: {matchLabels: {kubernetes.io/metadata.name: takserver}}
          podSelector: {matchLabels: {app: takserver-api}}
      ports: [{protocol: TCP, port: 8443}]
---
# platform/nats/network-policy.yaml -> nats-clients ingress.from (add)
        - namespaceSelector: {matchLabels: {kubernetes.io/metadata.name: hermes}}
          podSelector: {matchLabels: {app.kubernetes.io/name: onemind-hermes}}
```

The pod selector labels above come from the Deployment manifests in Git (`app: takserver-messaging`
and `app: takserver-api` are the Helm chart's labels; confirm with `kubectl get pod --show-labels`
before merging). The existing TAK policy rule admits the whole `hermes` namespace to 8089; tighten it to
`app.kubernetes.io/name: onemind-hermes` when Legacy (old pod) retires.

Retired at cutover: container `nats-leaf`, Secret `hermes-leaf-conf`, `nats-platform-patch-cm.yaml`,
the `omos-hermes` rule in `nats-leafnodes` for `onemind.io/agent: legacy`, and the `leaf_tak` password
rotation that follows (the Mac leaf and TAK leaf still use it).

## 7. Plugin set and the registry contract

### 7.1 Plugins

| Plugin id | Type | Connects to | Subjects / endpoints | Credentials Secret | Runs in | Phase |
|---|---|---|---|---|---|---|
| `onemind-agent` | platform | NATS | serves `agents.{prompt,status,attachments}.hermes.zeus.*`, `$SRV.*.agents`, heartbeat `agents.hb.hermes.zeus.*` | `onemind-hermes-nats` | `onemind-hermes` pod | 2 |
| `onemind-bus` (PR #564, as written) | standalone (16 tools + `pre_tool_call` retiring `nats_publish`) | NATS | requests `svc.entity.{get,query}`, `svc.layer.{list,meta,snapshot}`, `svc.policy.{command,task}.propose`, `svc.policy.action.get`; publishes `cmd.map.*` with client ack; **no confirm tool, no generic publish**. The old `nats-context` live-picture hook is not recreated; a bounded picture, if wanted, is built from `svc.entity.query` | `onemind-hermes-nats` | pod | 2 |
| `onemind-guard` | standalone (`pre_tool_call`, `post_tool_call`, approval transport) | NATS (audit only) | audit as log lines first; a bus audit subject needs an `AGENTS` export request to #562 (`evt.sys.legacy.*` is not granted there) | `onemind-hermes-nats` | pod | 2 |
| `onemind-tak` | platform | TAK Server 8089 mTLS | GeoChat and PLI as `Legacy` (map CoT via `svc.cot.send`); env `TAK_HOST`, `TAK_PORT`, `TAK_UID` | `onemind-hermes-tak` (cert, key, CA) | pod | 3 |
| `onemind-marti` | standalone (tools) | TAK API 8443 mTLS | missions, data packages, contacts; writes need approval | `onemind-hermes-tak` | pod | 4 |
| `onemind-homeassistant` | standalone (tools) | NATS only (the HA bridge at the edge owns the HA token) | reads via `svc.entity.*`; actions are policy proposals (`svc.policy.command.propose`), because #562 removes `cmd.ha.*` from Hermes (A19-P2) | none beyond NATS | pod | 4 |
| `onemind-obsidian` | standalone (tools, optional memory provider later) | Nextcloud WebDAV (`platform/nextcloud`) | vault read/search/append; no delete | `onemind-hermes-obsidian` | pod | 4 |
| `onemind-vrtak` | deferred | (XR client side) | none yet | none | none | later |

Telegram and Matrix stay bundled platforms and already run on `onemind-hermes`.

### 7.2 `plugin.yaml` contract (field names aligned with the Backstage catalog)

The merged catalog (PR #582, `platform/backstage/catalog/plugins/inventory.yaml` and its `README.md`)
already names the fields and says a future generator will read `**/plugin.yaml`. This contract uses
those names so the generator is a straight copy. One file per plugin, next to its code under
`platform/hermes/agent/plugins/<id>/plugin.yaml`. The top level is a normal Hermes manifest (the loader
reads it); the `onemind:` block is ignored by Hermes and carries the catalog fields.

```yaml
name: onemind-tak                       # Hermes key; catalog entity name = hermes-plugin-<name minus prefix>
label: TAK
kind: platform                          # Hermes kind (catalog `subtype`): platform | standalone
version: 0.1.0                          # semver; bumped by PR
description: Legacy as a TAK client: GeoChat and presence over a mutual-TLS streaming connection.
author: OneMind OS
requires_hermes: ">=0.22.5,<0.23"
requires_env:
  - {name: TAK_HOST, description: TAK Server messaging host}
  - {name: TAK_CERT_FILE, description: client certificate path}
  - {name: TAK_KEY_FILE, description: client key path, password: true}
provides_tools: [tak_send_cot, tak_contacts]
onemind:
  schema: 1
  title: onemind-tak                    # inventory: title
  type: hermes-plugin                   # inventory: type (spec.type). Catalog types: hermes-plugin | tak-server-plugin | bridge | mcp-adapter | ...
  subtype: platform                     # inventory: subtype (= Hermes kind)
  status: planned                       # inventory: status live | parked | planned (live only when a pod/volume was read)
  origin: onemind                       # inventory: origin onemind | upstream | third-party
  source: platform/hermes/agent/plugins/onemind-tak     # inventory: source
  build: "ghcr.io/onemindos/onemind-hermes image, plugin copied by init container install-plugins"   # inventory: build
  owner: group:default/platform-ai      # catalog owner ref
  runsOn:                               # inventory: runsOn {kind, location, host}
    kind: cluster
    location: "ns hermes / deployment onemind-hermes (container hermes)"
    host: resource:default/host-hermes-plugin-volume      # replaced by a new plugin-host Resource for the image-baked set
  connectsTo: [tak-server, nats]        # inventory: connectsTo (catalogued system slugs)
  publishes: []                         # inventory: subject slugs from `subjects:` in inventory.yaml
  subscribes: []
  providesApis: []
  consumesApis: [tak-cot-bridge]
  dependsOn: [component:default/onemind-hermes]
  subcomponentOf: component:default/onemind-hermes
  verified: "not yet"                   # inventory: verified (date + what was read)
  # fields the catalog does not hold yet; proposed additions, ignored by the current generator:
  nats: {user: hermes, account: AGENTS, requests: [svc.cot.send], serves: []}
  credentials:                          # names only, never values
    - {secret: onemind-hermes-tak, namespace: hermes, keys: [tls.crt, tls.key, ca.crt]}
  risk: {reads: [chat], writes: [chat, map-via-policy], physical: false, approval: guard-policy:tak-write}
  heartbeat: {kv: agent-caps, key: "hermes.onemind-tak", ttl: 60s}
```

Rules: `status: live` only when read from a running pod; every `publishes`/`subscribes` slug must exist in
`subjects:` of the inventory (the generator already fails dangling references); a plugin that uses a
subject outside the locked grammar is rejected by the catalog's `off-grammar` check. The catalog
notes that `hermes-project-stewardship`, `nats-context` and `tak-platform` have no source in the repo;
`onemind-agent`, `onemind-bus` and `onemind-tak` replace the latter two and the Synadia plugin, and
their inventory records move to `status: live` only after the phases below.

Runtime record (liveness, not config), written by the plugins' shared helper every 30 s to KV
`agent-caps` (TTL 1 min, bucket exists) under `hermes.<plugin-id>`:
`{"id","version","state","since","last_error","connected_to":[...],"hermes":"0.22.5","image":"sha256:..."}`.
Backstage joins static (Git) and live (KV) data; no record past its TTL means not running.

## 8. Security and approvals

Tiers of what a connection plugin can do:

| Tier | Examples | Control |
|---|---|---|
| T0 read | `tak.cot.>`, `ent.>`, Marti read, vault read | server-side subscribe permission only |
| T1 reversible digital write | chat reply, `evt.sys.legacy.*`, vault append, `cmd.map.>` | plugin allowlist plus `pre_tool_call` audit |
| T2 consequential digital write | Marti mission change, Home Assistant service call, TAK map marker to a broad group | Hermes manual approval (`approvals.mode: manual`) through a registered approval transport; request-bound decision |
| T3 physical or destructive | drone/robot/fleet/mission/intercept, door locks, deletes | **not available to the pod.** Section below |

T3 design. The agent's shell shares the pod, so anything mounted or in the environment can be read by
it. The control must live outside the pod, and the repo already has it in progress:

1. Legacy calls `onemind_command_propose` / `onemind_task_propose` (PR #564) which request
   `svc.policy.command.propose` (policy service, #554). The request carries `requested_by {kind: agent}`;
   `hermes` cannot publish `cmd.fleet.>` or `svc.policy.*.confirm` (server-enforced by #562, proved in
   #564's `test_policy_integration.py`). The proposal waits.
2. A human confirms from a CLIENT identity (`zeus`, Authentik/keychain on the device) through
   `svc.policy.*.confirm`; the policy service rejects confirms from any agent identity
   (`agent_cannot_confirm`) and publishes the real `cmd.*` with the approver recorded.
3. `execute` and `override` stay words only Zeus can supply: `execute` is the confirm of one proposal;
   `override` is a separate confirm that also bypasses a named interlock and requires a typed reason.
   Both are policy-service concepts, logged as `evt.so.action.*` facts. The chat `execute` the SOUL
   describes becomes an instruction to use the confirm surface, not a message the pod acts on.
4. Home Assistant actions use the same path (#562 A19-P2 removes `cmd.ha.*` from Hermes).

No new gate service is needed: #554, #562 and #564 already provide it.
Until those are live, T3 stays off.

Other rules:

- Credentials: SealedSecrets only, namespace-wide scope as the other `onemind-hermes-*` Secrets; no
  plaintext in Git; rotate the TAK cert with the CA. Do not copy `tak_admin.*`, `tokens.txt` or
  `mcp-tokens/` from Legacy (see MIGRATION.md).
- `allow_gateway_injection` stays false except for `onemind-agent` if proven necessary; `tools.override`
  never granted.
- Untrusted content: CoT remarks, chat text and NATS payloads are data. Tool results from these plugins
  are wrapped and length-capped; the live-picture injection is a bounded, structured summary
  (a few KB), not raw payloads, which also fixes the per-turn token cost of the old `nats-context`.
- Allowlists: `onemind-tak` honors `TAK_ALLOWED_USERS` (callsign/UID/cert CN) like other platforms;
  unknown senders get nothing. `onemind-agent` allowlists NATS-authenticated users.
- Audit: `onemind-guard` publishes every tool call on the three plugin families as
  `evt.sys.legacy.tool.called` (needs an AGENTS export; log lines until then) (name, plugin id, redacted args digest); the `AGENTS`/`EVT` streams keep it.

## 9. Phased plan with checks

| Phase | Work (each a reviewed PR; no step applies Terraform) | Check before moving on |
|---|---|---|
| **0 TAK bus bridge** | Merge #560; issue the `tak-client` cert; pin its digest; rotate the secrets #560 flags; remove `tak-nats-leaf` afterwards | `tak-client` `/readyz`; CoT on `tak.cot.inject` reaches an EUD; EUD CoT reaches `tak.cot.out` and `COT` grows by it; no duplicate injection (queue group) |
| **1 Pipeline** | derived image under `platform/hermes/onemind-hermes/image/`, init-container copy, managed config, a no-op `onemind-guard`; `nats-py` import test on Python 3.14 | `kubectl exec ... hermes plugins list` shows the plugin enabled from the managed list; a wrong `plugins.enabled` in the live profile cannot override it; image built from Git SHA, digest pinned; Legacy (old) untouched |
| **2 NATS** | Merge #562 with `hermes` widened to STANDARD, activate stage A1 (and the Hermes move) by the owner runbook; SealedSecret + CA ConfigMap + both NetworkPolicy edits; vendor and fix the Synadia adapter as `onemind-agent`; merge #564 `onemind-bus` and enable it; registry heartbeat | `connz` shows user `hermes` in `AGENTS` (not `onemind`); a denied publish to `cmd.fleet.drone.x.goto` returns a permissions violation (negative test, as #562's permtest); `onemind_command_propose` waits for a human confirm; `$SRV.PING.agents` lists the new instance; omos-ui discovery and heartbeats on `agents.*` identical to Legacy's; `agent-caps` shows the plugin records |
| **3 TAK** | Issue client cert + TAK user for `Legacy`; `onemind-tak`; Secret; egress rule | GeoChat to `Legacy` from an EUD gets a model reply; presence shows on the map; revoking the cert drops the connection within a minute; no `tak.cot.inject` use by Legacy |
| **3b Cutover** | Scale `hermes-legacy` to 0 (MIGRATION steps 7 and 8), remove the leaf, `hermes-leaf-conf`, patch ConfigMap; tighten TAK policy to the new pod; rotate `leaf_tak` | 24 h of heartbeats on `agents.hb.legacy.zeus.*`; no leaf named `hermes-legacy-leaf` in `leafz`; VolSync snapshot taken before retire |
| **4 Tools** | `onemind-marti`, `onemind-homeassistant`, `onemind-obsidian` (one PR each), approval transport | each tool listed in registry with live heartbeat; T2 action pauses for approval and executes only after it |
| **5 Actuation** | Policy service (#554) live; `onemind-bus` proposal tools enabled; a CLIENT confirm surface | a proposal without human confirm never reaches `cmd.*` (negative test); an agent-identity confirm is refused; expiry works; replay of an old confirm fails |

Rollback at every phase is Git revert plus Argo sync; Legacy (old) keeps running until 3b so there is no
conversation downtime. Phase 0 is independent of the Hermes work and unblocks everyone who publishes
`tak.cot.inject`.

## 10. What Legacy loses and gains versus the old instance

| Loses | Gains |
|---|---|
| Upstream `nats-platform` updates for free (we own a vendored fork) | Fixed import is a commit, not a startup patch; tests run in CI against the pinned image |
| `nats_publish` as an unrestricted tool and the `nats-context` firehose subscription | Server-side permissions per subject; bounded picture; `tel.>` excluded |
| The leaf's local `localhost:4222` and cluster-wide `ONEMIND` rights | Own identity `hermes` in `AGENTS`, revocable and auditable; one fewer container and out-of-band secret |
| Anonymous inject path for chat | Real TAK contact with a certificate, groups and presence |
| Plugins hand-installed on the volume | Reproducible install at every start; the volume holds state only |
| `memory-wiki`, `hermes-project-stewardship`, `onemind-maven` | Not carried over; reconsider individually (MIGRATION.md open decisions) |
| Ability to command hardware from a chat message | Proposals with human confirm enforced by the policy service (phase 5) |

## 11. Not verified and open decisions

Not verified (do these as phase checks, not assumptions):

- GHCR tag `takserver-base:plugins-5.8-onemind` existence and digest; that the overlay's `setenv.sh`
  edit makes the cluster launcher (`takserver-plugins-cluster.sh`, which sources `setenv.sh`) pick up `lib/`.
- `nats-py` on Python 3.14.9 and TLS verification behavior against the NATS CA.
- Managed-scope file name, precedence and list merge semantics in v0.22.5.
- JetStream permission strings for ordered consumers and KV against nats-server 2.10 (negative and
  positive tests in phase 2); labels `app: takserver-messaging` / `app: takserver-api` on live pods.
- Whether the `cot-tap-certs` Secret still chains to the CA rotated on 2026-10-29.
- Per-subject traffic of `COT` (no credentials used).

Owner decisions needed:

1. Phase 0 route: merge #560 (recommended) rather than redeploy the Java plugin.
2. `tak.cot.out` vs `evt.tak.cot.{uid}` (NATS-DOCTRINE 9.8) stays open and does not block this design.
3. Which TAK groups the `Legacy` cert may join (default: not `__ANON__` broadcast groups).
4. Where the first human confirm surface lives: omos-ui/XR client or Telegram/TAK relayed to the policy service.
5. Whether `onemind-obsidian` should also be Legacy's memory provider.
6. Direct TAK connection for chat (recommended) or a GeoChat-to-`agents.prompt.hermes.>` bridge in `tak-client`.

## 12. Verdict on the open PRs

| PR | Verdict | Change needed under this design |
|---|---|---|
| **#560** repair bus-to-TAK leg (`tak-client` + plugin source) | **Merge** (owner-approved, manual sync) | None to the design. Make it the only bus-to-TAK path; retire `tak-nats-leaf` and `cot-tap` after it works. Owner must rotate the two plaintext secrets it flags. Later merge its v1 codec with #551's (it says so) into one module |
| **#562** per-service NATS identities | **Merge, with one change** | Allow `hermes` `STANDARD` connections (account `AGENTS`, same lists) and drop the Hermes sidecar rendering (`leaf-conf --user hermes`, Hermes part of A2). Add `_INBOX.hermes.>`, `$KV.agent-caps.hermes.>`. Its Hermes lists (no raw feeds, no `cmd.*` but `cmd.map.>`) replace the broader draft in this document. Decide whether a bus audit subject needs an `AGENTS` export |
| **#564** `onemind-bus` plugin | **Merge, with changes** | Keep its tools and allowlist. Replace the "nats-py installed into /gateway-deps by the init container" install path with the derived image; default `ONEMIND_BUS_URL` to the direct `tls://nats.nats...:4222`, not `nats://localhost:4222` (sidecar); add the shared `plugin.yaml` `onemind:` block; verify the `pre_tool_call` retire hook on v0.22.5 (it was checked against 0.14.0) |
| **#551** TAK mesh gateway v0.3.0 | **Merge** | Unchanged; it is the edge piece (leaf at the site, multicast does not cross DOKS). Its Go codec is shared with #560 and with a future server GeoChat bridge |
| **#559** video: MediaMTX, `video-registry`, plugin manifests | **Merge; align manifests** | Its plugin manifests should use the `onemind:` field names of section 7.2 and register in the catalog inventory. No Hermes plugin is needed for video: Legacy reaches stream records through `svc.registry.video.*` once #562 grants it (not in Hermes' list today) |
| **#582** Backstage plugin catalog (merged) | Keep | Add records for the new plugins when each reaches `live`; later generate from `plugin.yaml` as its README plans |

Interlocks: #562 needs #560's `tak-client` user and #564's `hermes` lists to agree (they do, per its
second pass); #564's tests that use the A19 layout assume the sidecar-fed `hermes` user, so rerun them
with a STANDARD connection after the #562 change.
