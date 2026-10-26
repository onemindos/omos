# Fresh Hermes beside Legacy

User-directed fresh launch, October 26, 2026. This installation is independent of
Legacy; it does not import any existing agent data.

- Argo Application: `platform-hermes-fresh`, path `platform/hermes/fresh`.
- Deployment/Service/Ingress: `hermes-fresh`, namespace `hermes`.
- Private dashboard: `https://hermes-fresh.tail717ba0.ts.net` once the operator
  provisions it. Use the existing Hermes dashboard username/password.
- Image: official Hermes v0.22.5, pinned to the verified multi-architecture digest.
- New 25 GiB PVC: `hermes-fresh-data`. No old volume is mounted or cloned.
- Initial config selects the existing preferred model/provider, but **no model
  credentials are imported**. Sign in/configure the model on the fresh instance.
- Native bundled skills/default features remain available. No custom plugins,
  MCP servers, channels, schedules, Legacy memory or conversation history are copied.

## Isolation and persistence

The distinct app label keeps Legacy's selectors separate. The pod has no
Kubernetes service-account token, NATS sidecar, observer credential or Paperclip
credential. Network policy permits dashboard ingress from the Tailscale namespace,
DNS, and public HTTPS; private cluster/tailnet destinations are excluded from
HTTPS egress. Integrations will be enabled explicitly as migration proceeds.

The API is keyed independently with a SealedSecret and is not published by the
Service/Ingress. Dashboard session signing also uses an independent secret. Only
the existing dashboard login Secret is shared. Plaintext credentials are not in Git.

Config defaults are copied once by the init container, using exclusive creation.
Subsequent pod starts keep configuration authored through Hermes. Deployment uses
one replica and `Recreate`. PVC annotations prevent Argo prune/delete of its data.
Deleting or repurposing either installation's PVC remains a separate decision.

## Backup

`ReplicationSource hermes-fresh-data` snapshots the new PVC nightly at 01:10 UTC
using the existing VolSync mover/namespace policy. Retention follows Legacy's
7 daily / 4 weekly / 3 monthly / 1 yearly policy.

Its independently initialized Restic repository uses a distinct `hermes-fresh`
child prefix below the existing repository URI, outside the original `data/`,
`snapshots/` and `locks/` trees. It shares the existing storage credentials and
recovery password; the namespace-bound copy is sealed as `hermes-fresh-restic`.
This does not point both sources at the same Restic repository. Preserve the
existing off-cluster recovery password. First-backup completion and a restore test
must be observed; this source configuration alone does not establish recovery.

## Launch acceptance and gradual migration

Before merge, the exact image was started locally with a separate empty home and
network disabled. Dashboard and gateway health returned 200; unauthenticated
configuration access returned 401. No model inference was performed.

After Argo reconciles, record the actual rollout, gateway health, dashboard auth
and tailnet access in the deployment PR. A ready pod is not proof of working chat;
model sign-in and a real response are a separate acceptance step.

Add selected memory, skill packages, profiles and integrations incrementally after
the fresh instance works. Inspect credential-bearing assets before importing them.
The earlier migration helper remains separate unfinished work and is not used by
this launch. Keep Legacy running until the owner chooses the final cutover.
