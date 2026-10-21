# Legacy / LibreChat and Mac preflight — 2026-10-20

Scope: authorized read-only capability/production-readiness checks and staged
integration preparation. These observations do **not** mean LibreChat, Nextcloud,
a home cluster, Karmada or KubeEdge has been deployed. No credentials, billing,
provider configuration, files, storage or AI permissions were changed.

## Legacy: metadata works, inference fails

Observed 2026-10-20 around 16:15–16:18 UTC using the existing Hermes container
and its existing credential in memory; no secret values were printed or copied.

| Check | Actual result | Interpretation |
|---|---|---|
| `hermes/hermes-legacy` | Generation/observed generation 95; 1 ready replica; API reports Hermes 0.21.5 | Actual runtime, not a claim based on the older image tag |
| `GET /health` on loopback port 8642 | HTTP 200, status `ok`, platform `hermes-agent` | API process health only |
| Anonymous `/v1/models` and `/v1/capabilities` | HTTP 401 | Protected metadata refuses anonymous callers |
| Authenticated `/v1/models` | HTTP 200, model `hermes-agent` | Correct agent endpoint, not a direct hosted-model endpoint |
| Authenticated `/v1/capabilities` | HTTP 200 | Advertises chat/Responses streaming, runs, events, stop/steer/approval, sessions and durable idempotency; not executed acceptance of those operations |
| Session headers | Advertises `X-Hermes-Session-Id` / `X-Hermes-Session-Key` | LibreChat session isolation still requires runtime verification |
| Detailed health / model options | HTTP 200; configured provider `digitalocean`, default `nemotron-3-ultra-550b` | No provider or model change performed |
| Harmless streaming chat | HTTP 200 SSE, **no answer**, terminal `finish_reason: error`, no tool events observed, about 3.93 seconds | **Failed** inference; HTTP 200 is not a successful chat |

The single chat probe asked only for `ONEMIND_LEGACY_API_OK`, no tools, file reads,
commands, delegation or changes. Runtime logs at **16:17:38–39 UTC** show the
configured DigitalOcean inference provider returning **402 Payment Required**:
“This account has been blocked.” This diagnoses the observed failure; it does not
establish the account's underlying billing cause. No further inference retries,
billing changes, provider switches or key changes were attempted. The request
instruction and `tool_choice: none` are not asserted to enforce harness permissions.

Port 9119 returned dashboard SPA HTML even for `/v1/*`. A 200 HTML fallback is not
an API compatibility check. Temporary local forwarding was stopped after checks.

## NATS: discovery is not full end-to-end acceptance

The existing local leaf at `nats://127.0.0.1:4222` returned an actual microservice
discovery response: `agents`, protocol 0.3 / advertised SDK 0.4.1, metadata
agent `hermes`, owner `zeus`, session `legacy`. Discovered endpoint subjects:

- `agents.prompt.hermes.zeus.legacy`
- `agents.status.hermes.zeus.legacy`

Attachments were advertised as supported with a 1 MB max payload. Only Legacy
was discovered in that bounded sweep; it does not prove other agents are absent
at all times. A status request returned `{"stream":"AGENTS","seq":308529}`,
which is a JetStream acknowledgement, **not an agent health result**. No NATS
inference prompt, approval, cancellation, attachment or delegated spend test was
performed. Discovery subjects/contracts were not changed.

## Mac production readiness

Read-only observations around **16:15 UTC**, with capacity rechecked at 16:27 UTC:

| Item | Observed | Gate |
|---|---|---|
| Machine | M2 Max, 32 GiB RAM, macOS 27.0.1 | Linux ARM64 VM sizing and app images must be tested |
| Internal Data volume | About 869 GiB used, **13–14 GiB free**, 99% capacity | **Blocks a new production VM**; no deletion/cleanup authorized |
| `OM PRIMARY` external volume | About 3.6 TiB total / 739 GiB free | Existing data, not spare storage approved for repurposing |
| `MBP TM` external volume | About 1.8 TiB total / 100 GiB free | Existing backup disk; do not reuse as app data |
| Power | AC attached, battery 100%; AC sleep/disksleep/standby/display sleep 0 | Settings do not prove closed-lid operation, cooling, or recovery |
| FileVault | On | Cold boot/unattended unlock and recovery procedure still untested |
| Docker Desktop | Running; existing image builder | Not silently replaced or converted into production hosting |
| Lima / Colima | Lima installed with no instances; Colima builder profile stopped (4 CPU / 8 GiB / 24 GiB, ARM64) | Existing build profile is not a production cluster |

No VM was created, disk moved/erased, backup configuration changed, or Home
Assistant VM touched. User explicitly wants the Mac on 24/7, plugged in and
closed; validate that operating mode and startup after logout/reboot/power loss
before promising production availability. A single Mac is one failure domain,
not HA. K3s would run **inside Linux**, not natively on macOS.

## Prepared source, not a deployed service

The [pilot configuration/runbook](../../platform/librechat/README.md) pins chart
2.0.16 / app v0.8.8 and an image manifest covering amd64 and arm64. It uses private
ingress, required OIDC group, mandatory per-key secret references, no Kubernetes
token, default-deny networking and a user/conversation session-key template.
The bundled unauthenticated MongoDB is disabled. **Database, backup/restore,
sealed credentials and the dedicated identity provider are still prerequisites.**
Do not sync it, call persistence proven, or declare the full stack production-ready.

Local tests: 22 offline source/stream regression tests passed (2 rendered tests
intentionally skipped in offline mode); a separate real pinned-chart render
passed all 12 source/render tests. The render caught and fixed the chart-specific
NetworkPolicy label mismatch. Argo skips the upstream unrestricted Helm test pod.
No browser login/session/conversation test has passed for this undeployed pilot.
