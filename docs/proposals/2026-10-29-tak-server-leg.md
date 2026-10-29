# Owner proposals: TAK Server leg (A21)

Date: 2026-10-29. PR: `feat/tak-server-leg`. Context: onemind-xr FEATURES A21, BUS-DESIGN
§3.1/§7.2, `docs/research/omos-live-state.md` §6 (in onemind-xr).

**No new subject is introduced.** `tak.cot.out` and `tak.cot.inject` are in the locked
grammar and in stream `COT` (`platform/nats/streams/cot.yaml`). Message headers on
`tak.cot.out` (`Nats-Msg-Id`, `Tak-Client`, `Tak-Proto`, `Tak-Server`,
`Tak-Undecoded-Extensions`) follow the mesh gateway's header style (PR #551).

| ID | Proposal | Why | Default if not decided |
|---|---|---|---|
| A21-P1 | Rebuild the NatsBridgePlugin JAR 1.1.0 from `services/takserver/plugin-src` with the TAK 5.8 SDK (tak.gov), commit it to `services/takserver/plugin/`, run `build-takserver.yml`, and pin the resulting `plugins-5.8-onemind` digest in `platform/takserver/deploy/takserver-plugins-deployment.yaml` | The running plugins image has no overlay (`-Dloader.path` without `lib/`, no plugin log line), so nothing injects `tak.cot.inject` into TAK. The SDK is not in Git, so this needs someone with tak.gov access. `platform-takserver` auto-syncs, so the pin PR restarts the plugins pod | tak-client alone injects (it does both legs); the plugin stays unloaded |
| A21-P2 | A dedicated NATS user `tak_client` in account `ONEMIND` (NATS Helm `values.secret.yaml`): publish `tak.cot.out` plus the JetStream API for `COT` pub acks (`$JS.API.>` scoped as the account split allows), subscribe `tak.cot.inject` and `_INBOX.>` | BUS-DESIGN §7.2 gives the TAK identity exactly that (`TAK leaf`, pub `tak.cot.out`, sub `tak.cot.inject`). Until the account split, the Secret carries the shared `onemind` user, as `omos-nats-creds` does for the entity service | shared user |
| A21-P3 | `tak-client` is the only producer of `tak.cot.out`: set `natsProducerEnabled: false` in the plugin config once 1.1.0 runs, and change the BUS-DESIGN §3.1 / §7.2 publisher of `tak.cot.out` from "NatsBridgePlugin" to "tak-client" | Cluster-mode TAK delivers nothing to plugin `onMessage` (cot-tap README); one producer keeps one `Nats-Msg-Id` scheme and lets the echo filter work | both enabled; plugin `onMessage` stays silent in cluster mode |
| A21-P4 | Rotate the TAK CA key passphrase and the `omos-cot-tap` TAK user password that `services/takserver/cot-tap/README.md` carried in plain text (replaced by placeholders in this PR; still in Git history) | Exposed credentials; rotation, not suppression (omos `CLAUDE.md` Secrets) | — |
| A21-P5 | Delete `services/takserver/cot-tap/` after tak-client has published `tak.cot.out` live for a week; revoke the `omos-cot-tap` user and certificate | cot-tap was never deployed through Git, runs an unpinned image with `pip install` at start and server verification off | keep the deprecated directory |
| A21-P6 | TAK groups for `omos-tak-client` | Its certificate's groups decide which TAK CoT reaches the bus and to whom its injects go; the plugin injects into `__ANON__` and `Red` today | same as the plugin: `__ANON__`, `Red` |

Follow-ups, not in this PR: Data Sync missions (FEATURES C7) and video feed registration
(`/Marti/api/video`) through tak-client on 8443; merge the TAK v1 codec copies of
`bridges/tak` and `bridges/tak-client` into one module after #551 and this PR merge.
