# Paperclip → existing Hermes gateway

This is source for private HTTPS transport, not evidence that the agent is connected
or that a model invocation has succeeded. Paperclip's native `hermes_gateway`
adapter rejects remote plaintext HTTP. No Codex device login is required for a
Hermes agent.

## Transport and ownership

`platform-legacy` owns an official nginx 1.30.6 Alpine sidecar, private CA Issuer,
root/leaf Certificates and ClusterIP Service in Hermes. The verified official image
index is pinned by digest; nginx stable release details are at
https://nginx.org/ and image sources at https://hub.docker.com/_/nginx.

Endpoint: `https://hermes-paperclip-api.hermes.svc.cluster.local:8643`.
The proxy sends only the installed adapter's endpoints to `127.0.0.1:8642`:

- `GET /health`, `GET /v1/capabilities`
- `POST /v1/runs`, `GET /v1/runs/{id}`
- `GET /v1/runs/{id}/events`, `POST /v1/runs/{id}/stop`

Authorization, Idempotency-Key and X-Hermes-Session-Key are preserved. nginx does
not inject an API key or bypass Hermes authorization. Request access logs are
disabled and event streaming is unbuffered. No dashboard, skills, arbitrary control
or model-provider endpoints are proxied. A TLS-ready TCP probe establishes only
listener readiness, not bearer validation, provider health or end-to-end success.

nginx runs nonroot with a read-only root filesystem and a bounded emptyDir. Its
certificate path variables and one-minute certificate cache load leaf rotations
without a Hermes restart. No service or ingress exposes port 8642 directly.

Paperclip's existing default-deny policy adds only TCP/8643 egress to the precise
Hermes pod identity. Callback ingress is only TCP/3100 from that same identity.
There is deliberately no new Hermes ingress-isolation policy: an ingress policy
would isolate the entire existing pod, including its dashboard and NATS leaf.
Inventory their current consumers before introducing a default-deny rule. The API
still requires its existing bearer authorization, and the Service is cluster-only.

## Ordered activation (coordinator only)

1. Review the `platform-legacy` diff and preserve the current Hermes model, persona,
   PVC, NATS identity and runtime API configuration. Confirm the native API is
   actually bound to loopback; transport does not change that runtime setting.
2. Selectively sync the two Issuers and two Certificates, then wait for the CA,
   issuer and leaf to report Ready before rolling the Recreate Hermes Deployment.
   The leaf Secret is mandatory. A failed or absent issuer must not trigger an
   unprepared Hermes restart.
3. Read **only the public certificate** from `hermes/hermes-paperclip-ca` (`tls.crt`),
   verify it against the issued leaf/SAN, and seal that certificate as a strict
   SealedSecret named `paperclip/hermes-gateway-ca`, key `ca.crt`. Commit this new
   SealedSecret and register it in the Paperclip kustomization before activation.
   Do not copy either private key or the full TLS Secret to Paperclip. The trust
   volume is optional so staging transport does not disable Paperclip's UI; HTTPS
   agent calls fail certificate validation until the trust Secret is present.
4. Sync the public trust Secret and the reviewed Paperclip changes, then restart
   Paperclip through Argo after the file exists: Node reads NODE_EXTRA_CA_CERTS at
   process startup. Never use NODE_TLS_REJECT_UNAUTHORIZED=0 or insecureSkipVerify.
5. Through the signed-in owner's normal controls, create a dedicated Hermes agent
   using the endpoint above and the current Hermes bearer via Paperclip's secret
   input. Keep Paperclip's heartbeat scheduler disabled until one bounded test
   succeeds. Do not put bearer tokens into URLs, logs, manifests or browser notes.
6. If the adapter requires Hermes callbacks, issue a dedicated **Paperclip agent
   credential**, seal it strictly to `hermes/paperclip-agent`, key `credentials.json`, and mount it
   read-only at `/run/paperclip/credentials.json` in Hermes. Configure only this agent's
   callback to `http://paperclip.paperclip.svc.cluster.local:3100`; this narrow
   in-cluster callback is governed by Paperclip authorization and the scoped ingress
   rule. A callback credential is a distinct identity, never the Hermes gateway key
   or a human session. The optional read-only mount is prepared by this PR; the coordinator must add its
   strict SealedSecret before syncing the Hermes rollout. Its JSON format must
   match the installed adapter contract; never invent fields from memory.
7. Verify trusted HTTPS, health/capability detection, unauthorized API rejection,
   one bounded run, event stream, result persistence and cancellation. Do not send
   unsolicited Matrix/NATS/TAK messages or enable recurring autonomous work as a
   side effect. Record actual test outcomes separately from source readiness.

## Recovery and renewal

All existing persistent data and backups are unchanged. Removing this integration
means removing the Paperclip agent first, then reverting the sidecar and dedicated
transport resources through Argo; do not delete a namespace or PVC. Hermes uses
Recreate, so any sidecar rollout briefly restarts the existing agent.

cert-manager renews the 90-day leaf 30 days before expiry; nginx's variable-path
cache reloads material within one minute. Root rollover is a separate operation:
seal a bundle containing both public roots into Paperclip, roll it, reissue the
leaf against the new root, verify, then retire the old trust. CA private-key rotation
is disabled to avoid silently invalidating the consumer's pinned public trust.
The 10-year root itself still needs this explicit rollover before its renewal.

## Activation record, 24 October 2026

Zeus approved Paperclip board CLI access. The Hermes employee was created paused
with an encrypted company-secret reference, timers disabled, and agent/skill creation
disabled. Its callback key uses `task_bridge`, bounded to the integration parent
issue, with no additional assignees. Callbacks must use the assigned issue ID
directly; company-wide search and secret administration are not permitted.
The current Hermes provider passed an authenticated capability check and a real
no-tool response before rollout. Preserve that provider configuration: the optional
`hermes-litellm` Secret is absent, so the existing merge init container is a no-op.
`provider-config-pre-paperclip-v1` preserves the current configuration separately
from the earlier provider backup. New callback credentials are strictly sealed.
The exact nginx image and configuration passed a nonroot, read-only, capability-free
`nginx -t` test in a disposable container on the existing build runner.
These checks do not yet establish the full Paperclip-to-Hermes task round trip.

The corrected proxy also passed an actual TLS handshake with SAN verification
in the same restricted container. Wrong-hostname verification failed closed,
unlisted paths returned 404, and a disallowed HTTP method returned 403.
