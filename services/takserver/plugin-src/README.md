# NatsBridgePlugin (TAK Server ↔ NATS)

TAK Server plugin that injects CoT from `tak.cot.inject` into TAK Server and, where TAK
delivers CoT to plugins, publishes it to `tak.cot.out`. The built JAR and its runtime
config live in `../plugin/`; `docker/Dockerfile.takserver-plugins.onemind` bakes both into
the plugins image.

| File | What | Needs the TAK SDK |
|---|---|---|
| `src/main/java/.../NatsBridgePlugin.java` | TAK glue: `MessageSenderReceiverBase`, CoT conversion, `send()` | yes |
| `src/main/java/.../BridgeCore.java` | Config parsing, inject validation, counters | no |
| `src/main/java/.../NatsLink.java` | NATS connection, inject queue subscription, forward, status line | no |
| `src/test/java/...` | JUnit 4 tests of the two SDK-free classes (NatsLink against a real nats-server) | no |
| `test.sh` | Runs those tests with plain `javac` (JDK 17) and checksum-pinned JARs; CI: `.github/workflows/test-takserver-plugin.yml` | no |

## Is it loaded? (A21 diagnosis, 2026-10-29, read-only)

Live `takserver-plugins` pod (`kubectl -n takserver logs deploy/takserver-plugins`, pod
started 2026-10-18, 379 lines):

1. Line 1: `NOTE: Picked up JDK_JAVA_OPTIONS: -Dloader.path=WEB-INF/lib-provided,WEB-INF/lib,WEB-INF/classes ...`.
   The OneMind overlay (`Dockerfile.takserver-plugins.onemind`) rewrites exactly that string
   to `...,WEB-INF/classes,lib/` so the plugin manager loads `/opt/tak/lib/*.jar`. The live
   value has no `lib/`, so **the running image does not carry the overlay** and the plugin
   JAR is never on the class path. This matches `platform/takserver/values-onemind.yaml`,
   which names `plugins-upstream-5.8-onemind`, the tag `build-takserver.yml` builds without
   the plugin. The running digest (`sha256:f4824fba…` in
   `platform/takserver/deploy/takserver-plugins-deployment.yaml`) could not be mapped to a tag:
   GHCR returns 403 without `read:packages`.
2. `PluginService`, `PluginLoader` and `PluginMessenger` lines do reach stdout, but no
   `NatsBridgePlugin` line does, although the 1.0.0 constructor logs its config at INFO.
3. On the hub, the `tak-nats-leaf` leaf shows no `tak.cot.inject` interest
   (onemind-xr `docs/research/omos-live-state.md` §6).

Other defects found in the 1.0.0 source, fixed in 1.1.0:

| 1.0.0 | Effect | 1.1.0 |
|---|---|---|
| `Nats.connect` once in `start()`; on failure it logs and never retries | If `tak-nats-leaf` is not up when TAK starts, the bridge is dead until the next pod restart | `Nats.connectAsynchronously(options, true)`: `start()` returns at once, connects in the background, subscribes on the first (re)connect, replays after reconnects |
| jnats and its eddsa dependency shaded **unrelocated**, slf4j-api bundled | TAK 5.8 ships its own jnats (it uses NATS internally) and SLF4J on the same class path; duplicate classes, first one wins | `io.nats` and `net.i2p` relocated under `tak.server.plugins.natsbridge.shaded`; slf4j-api `compileOnly` |
| Plain subscription on `tak.cot.inject` | Running tak-client next to it would inject every message twice | Queue group `natsQueueGroup` (default `tak-cot-inject`), shared with tak-client |
| Any payload passed to `cotStringToDataMessage` | A message with two events, or protocol-control CoT, reaches TAK | `BridgeCore.rejectReason`: one event, has a type, no `t-x-takp-*`/`t-x-c-t*` |
| Only a start-up log line | "Loaded but idle" and "not loaded" look the same | `NatsBridgePlugin v1.1.0 loaded: …` at construction, `subscribed to …`, and every 60 s `NatsBridgePlugin v1.1.0 status nats=CONNECTED injected=N injectRejected=N injectFailed=N forwarded=N forwardDropped=N` |
| Overlay `sed` silently does nothing if `setenv.sh` changes | An image without the plugin builds green | The Dockerfile fails the build if the patch did not apply or if there is not exactly one plugin JAR |

New config keys (`../plugin/tak.server.plugins.NatsBridgePlugin.yaml`): `natsQueueGroup`
(`none` = plain subscription), `natsProducerEnabled` (set `false` once tak-client publishes
`tak.cot.out`), `statusIntervalSeconds` (`0` = off). Unknown types fail the plugin start
with a clear message instead of being ignored.

## Build (needs the TAK Server SDK; not possible in CI or on the A21 build machine)

The committed `../plugin/takserver-nats-plugin-1.0.0.jar` was built from the 1.0.0 source.
The 1.1.0 source compiles only against the TAK Server 5.8 plugin SDK JAR from tak.gov:

```sh
# from services/takserver/plugin-src, with the SDK JARs in ../tak-server-sdk/tak-server-sdk/
gradle clean shadowJar test           # Gradle 8, JDK 17
unzip -l build/libs/takserver-nats-plugin-1.1.0.jar | grep -c 'tak/server/plugins/natsbridge/shaded/io/nats'   # > 0
unzip -l build/libs/takserver-nats-plugin-1.1.0.jar | grep -E ' (io/nats|org/slf4j)/' && echo "NOT RELOCATED"   # prints nothing
git rm ../plugin/takserver-nats-plugin-1.0.0.jar
cp build/libs/takserver-nats-plugin-1.1.0.jar ../plugin/
```

Then the plugins image: run `build-takserver.yml` (workflow_dispatch, self-hosted runner),
take the digest of `ghcr.io/onemindos/takserver-base:plugins-5.8-onemind` from its log, and
pin it in `platform/takserver/deploy/takserver-plugins-deployment.yaml` in a PR. That
Deployment is Argo-managed with automatic sync, so merging the PR restarts the plugins pod.

## Verify live (after the image change)

```sh
kubectl -n takserver logs deploy/takserver-plugins | head -1          # loader.path ends in ,lib/
kubectl -n takserver logs deploy/takserver-plugins | grep NatsBridgePlugin
#   NatsBridgePlugin v1.1.0 loaded: url=nats://tak-nats-leaf... consumer=tak.cot.inject queue=tak-cot-inject ...
#   NatsBridgePlugin subscribed to tak.cot.inject (queue tak-cot-inject)
#   NatsBridgePlugin v1.1.0 status nats=CONNECTED injected=... (every 60 s)
nats server request leafz    # the tak-nats-leaf connection now lists a tak.cot.inject subscription
```
