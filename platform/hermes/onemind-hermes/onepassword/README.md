# Native 1Password → ClickUp — source-only draft, NOT DEPLOYED

This is a **non-active source proposal**, not installed or deployed integration.
Production `deployment.yaml`, `config.defaults.yaml` and `kustomization.yaml` are
restored byte-for-byte to base `85a52b599b1b18b8f4b6a38cee9040c195ee3c9b`.
They do not install `op`, inject bootstrap, reconcile this source or enable native
1Password defaults. The runtime image is unchanged; its known child-filter defect
is not exposed by new bootstrap wiring in this draft.

The proposed wiring is retained under [`proposal/`](proposal/README.md), outside
the active Kustomize graph. It is not a runnable overlay or safe activation payload.
All startup behavior below describes that future proposal, **not current production**.
It would provide authenticated app access without a new MCP/Connect server. Only
the cloud **default** Hermes home `/opt/data` would be configured; named worker
profiles would not be rewritten. `settings.json` binds the approved immutable vault/item IDs to
`CLICKUP_API_TOKEN`. The item title is deliberately not used. No token values,
credential files, session database or live cluster state are needed by this code.

## Startup contract and image requirements

The current Git manifest uses the official image for the gateway and all init
containers, unchanged:

```
nousresearch/hermes-agent:v0.22.5@sha256:9c5561579f3ffdf333401797094ac3a743738d6c25d40f4331c7bf77b694238d
```

The manual-only [runtime build workflow](../../../../../.github/workflows/build-hermes-runtime.yml)
now tracks the corrected-image build, synthetic tests and GHCR publication; it has
not been dispatched by this local change. The follow-up
[`runtime-fix/`](runtime-fix/README.md) now contains a hash-locked derived-image
recipe and actual-filter regressions: **the current upstream digest leaks the
bootstrap token into generic terminal children**. Removing its passthrough entry
is insufficient. Do not roll out this branch until a corrected image is built,
verified and pinned by its real digest; no successful build or new image is
claimed here. The non-active proposal also defines a bounded init container to
install official **op 2.41.1** into a pod-scoped `emptyDir`.
It verifies the amd64 ZIP SHA256 before extracting only `op`:

```
https://cache.agilebits.com/dist/1P/op2/pkg/v2.41.1/op_linux_amd64_v2.41.1.zip
e52476700cba9202ac7469bdfb2648ecdbf3438ac85822010f77230139e74966
```

The checksum was supplied from the maintained Homebrew CLI manifest and verified
against the downloaded archive. The deployment explicitly selects `amd64` and the
installer rejects other architectures. Add arm64 only with a reviewed checksum
and execution test. Existing image/resource limits remain; the installer adds
10m/32Mi requests, 500m/192Mi limits. Init containers run sequentially, so this
request is below the existing init and gateway requests.

The main container mounts `/opt/onepassword` read-only and the native source pins
`binary_path: /opt/onepassword/op`. CLI configuration uses a separate 16Mi memory
volume at `/var/run/onepassword`, prepared mode 0700 for the image's `hermes` user.
No bootstrap or resolved value is written onto the retained PVC by this setup.
Native `cache_ttl_seconds: 0` disables both Hermes secret-cache layers; it does not
delete any historical cache files. If caching was ever enabled before rollout,
treat pre-existing caches/backups as sensitive and handle separately with owner
authorization, not automatic cleanup.

On **every** pod start `initial-config` first seeds an empty PVC using the existing
exclusive-create logic, then reconciles through the supported Hermes config CLI:

- only `secrets.onepassword` is owned, not other sources or model/MCP settings;
- only `CLICKUP_API_TOKEN` is added to `terminal.env_passthrough`, preserving
  unrelated existing entries and removing `OP_SERVICE_ACCOUNT_TOKEN` if present;
- no bootstrap Secret is injected into any init container;
- no credential lookup, login or native setup command writes `.env`/`.op.env`.

The hashed ConfigMap changes the pod template when the installer/settings change.
A normal replacement pod re-downloads the archive; public 1Password CDN access is
therefore a cold-start prerequisite. Download/checksum failure blocks this init
rather than running an unverified binary. An alternative future offline-ready
build must bake this exact CLI into a reviewed derived image, publish it, and pin
its **real** digest in Git; this change does not invent such an image. Existing
NetworkPolicy already permits DNS and public TCP/443 for CDN, 1Password and
ClickUp. No egress broadening, service account, sidecar or bus change is included.

## Required integration before activation (intentionally missing here)

The parent/authorized operator must produce **only an encrypted SealedSecret**:

- file: `onepassword/bootstrap.sealedsecret.yaml` beside these files;
- namespace: `hermes`; Secret name: `hermes-onepassword-bootstrap`;
- exactly the bootstrap key `OP_SERVICE_ACCOUNT_TOKEN`;
- strict name/namespace scope, following the existing Zoho convention;
- include `onepassword/bootstrap.sealedsecret.yaml` in the parent Kustomization.

No dummy SealedSecret or placeholder ciphertext is included. The proposed Secret
reference is required, not optional, but is **absent from active manifests**. Native
1Password activation is blocked until the real encrypted resource is integrated. Keep the public certificate outside Git. No certificate
was found in this repository. An authorized operator with existing Kubernetes
access can retrieve the controller's **public** certificate through its API proxy
(without access to signing-key Secrets):

```sh
kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > "$TMPDIR/hermes-sealed-secrets-public.pem"
```

Run that from an authorized operator environment, not by searching for kubeconfig
or credentials in this pod. This uses API service-proxy access and avoids relying
on pod-to-controller HTTP reachability, which can be denied by NetworkPolicy. The
controller's public endpoint is `/v1/cert.pem`; the parent reported direct HTTP to
`sealed-secrets-controller.kube-system.svc.cluster.local:8080` failed. That failure
does not justify opening new egress or reading private controller keys. The fetch
has **not** been verified by this worker; public-certificate retrieval/sealing is
still a blocker. Have the authorized operator validate sealing against the target
controller before integrating the ciphertext.

The owner must rotate the exposed bootstrap token in 1Password and reseal the
replacement. Never repeat the previous value/history. Scope the service account
to read access to the approved vault; the token can read every item granted to
that account, not just the configured reference. Never allow the bootstrap token
through terminal passthrough or a skill's required environment list.

## Verification, recovery and rotation gates

Local tests use synthetic **configuration**, not authentication responses:

```sh
python platform/hermes/onemind-hermes/onepassword/test_onepassword.py -v
python platform/hermes/onemind-hermes/mcp/whop/test_whop.py -v
python platform/hermes/onemind-hermes/mcp/zoho/test_reconcile.py -v
make validate
```

The first suite includes real isolated Hermes CLI runs against seeded-empty and
existing test homes, repeated reconciliation, preservation of unrelated settings,
bootstrap exclusion, no cache file, checksum failure and unsupported architecture.
Actual archive download/hash verification plus `op --version` is a separate smoke
check. These do not establish runtime authentication, secret isolation in the
pinned image, Kubernetes unsealing, or a successful recovery drill.

Before activation, build/scan/exercise a corrected runtime image and pin its real
digest for the resolving gateway; verify bootstrap denial in real terminal and
execute_code children and bootstrap retention only by the native resolver. Then
review and deliberately reintroduce the proposal wiring, its resource limits,
CLI checksum/amd64 selection, ephemeral writable state and startup reconciliation
with a real rotated bootstrap SealedSecret. Keep bootstrap out of every init
container. Source-only regression hashes must be deliberately replaced by new
activation acceptance tests in that future change. After owner approval, review,
merge and Argo reconciliation,
an authorized operator must verify separately:

1. Init completion, expected image digest, `op --version`, gateway readiness and
   native-source configuration (do not print environment or credentials).
2. Read-only ClickUp `GET https://api.clickup.com/api/v2/user` from a default-profile
   Hermes terminal child, using the resolved token in the Authorization header.
   Report only HTTP status and confirmation of the intended account, never token
   or full identity data. Check bootstrap absence as a boolean in the same child;
   do not dump the environment. CLI startup alone is not authentication evidence:
   native-source errors warn and continue, possibly leaving an older env value.
3. Gateway restart/replacement with the same PVC, then repeat the identity check.
4. Authorized isolated empty-volume recovery: Git reseeds/reconciles config,
   SealedSecrets restore bootstrap, CDN reinstalls CLI and the vault resolves the
   app token. Do not replace or delete production claims. Other model/channel/OAuth
   credentials require their own encrypted recovery; Git alone cannot restore them.
5. Rotate the ClickUp credential in the same vault field, start a new Hermes
   process and repeat the read-only identity probe. Rotation is startup-based, not
   a promise that a running process refreshes automatically. Rotate bootstrap by
   resealing and changing the pod template through Git to reload Secret env vars.
6. Verify cache remains disabled, revoked/expired token behavior is visible, and
   recovery of independently protected SealedSecrets signing keys is available
   through the existing recovery runbook. This worker never accesses those keys.

This implementation is **NOT DEPLOYED**. No successful backup/restore is claimed. Service-account expiry,
API quotas, upstream availability and revocation still apply; disabling the cache
increases lookup volume. Rollback requires reverting both defaults and startup
reconciliation through Git; reverting defaults alone cannot change an existing
PVC's reconciled config. Preserve data and encrypted recovery material.

Sources: [native integration](https://hermes-agent.nousresearch.com/docs/user-guide/secrets/onepassword),
[secret isolation](https://hermes-agent.nousresearch.com/docs/user-guide/secrets),
installed `/opt/hermes/agent/secret_sources/onepassword.py`, existing
[PVC recovery](../../../../docs/portability/REBUILD-FROM-ZERO.md).
