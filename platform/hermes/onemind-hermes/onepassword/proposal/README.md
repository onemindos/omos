# Non-active wiring snapshot — NOT DEPLOYED

The three YAML files preserve the earlier proposed parent-directory wiring for
review. The active parent Kustomization does not reference this directory or any
1Password source file. Paths in this snapshot are relative to the original parent,
not this directory: it is not a runnable Kustomize overlay. Do not apply it directly.

The snapshot intentionally preserves an unsafe-to-activate image reference and a
bootstrap Secret reference whose real encrypted resource is not yet supplied.
Do not copy it into production unchanged. Before activation, follow the corrected
image, rotated strict-scoped SealedSecret, reviewed wiring and authorized runtime/
recovery gates in [the proposal README](../README.md). No image build, publication,
bootstrap injection or deployment is established by these source files.

## Substitute only a verified build digest

After a successful manual runtime build, download its `hermes-runtime-image.txt`
artifact into operator scratch. Verify the run's source commit, synthetic test
results, remote digest and image review/scan. This file must contain the real
`ghcr.io/onemindos/hermes-agent@sha256:...` reference; do not construct a digest
from a commit SHA or copy a synthetic test fixture.

```sh
proposal=platform/hermes/onemind-hermes/onepassword/proposal
python "$proposal/prepare_image.py"   --image "$(<"$TMPDIR/hermes-runtime-image.txt")"   > "$TMPDIR/hermes-proposed-deployment.yaml"
```

`prepare_image.py` validates the exact corrected-image repository and full digest
format, rejects the known upstream digest and snapshot drift, and changes **only**
the resolving `hermes` container's image. Init containers have no bootstrap and
keep their existing upstream pins. It emits the review candidate on stdout; it
does not write active files, create a Secret, verify image provenance, or make the
snapshot runnable. Digest format alone is not evidence the image exists or is safe.

Before a separate activation change, rebase the snapshot onto current production
wiring so concurrent changes are preserved; integrate the real rotated,
strict-scoped `hermes-onepassword-bootstrap` SealedSecret and parent-relative
ConfigMap paths/resources. Replace source-only byte-hash acceptance tests with
reviewed activation checks. Render/validate the actual parent Kustomization,
review all bootstrap references (main only), resource limits, CLI checksum and
memory-volume ownership. Do not enable this proposal while image or ciphertext is
missing. Authorized review/merge/Argo reconciliation and native authentication,
replacement-pod and empty-volume recovery gates remain separate, unexecuted work.
