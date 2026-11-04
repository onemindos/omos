# Local verification receipt

Base inspected: public `onemindos/omos` main
`85a52b599b1b18b8f4b6a38cee9040c195ee3c9b` (Whop startup/recovery #622).
Branch: `feat/hermes-native-onepassword`.
Worktree: `/opt/data/cache/scratch/omos-onepassword-worktree`.
No credentials were read, supplied, resolved or written for this work. No cluster
mutations, publication or merge were performed by this worker.

## Source-only safety correction — NOT DEPLOYED

Production `deployment.yaml`, `config.defaults.yaml` and `kustomization.yaml`
were restored byte-for-byte to the original base above. The previous wiring is
retained under `onepassword/proposal/` as a non-active snapshot (not an overlay).
No proposal/source/test files were deleted. The upstream runtime image is unchanged;
its known child-filter flaw remains, but no new bootstrap injection or native
1Password activation is present in the production graph.

Current local execution:

- Native suite: **6 tests passed**, 76.800s, including real isolated Hermes CLI
  reconciliation and new original-base SHA256 / no-activation regressions. The
  former wiring test now explicitly targets only the non-active proposal snapshot.
- `git diff --exit-code 85a52b599b1b18b8f4b6a38cee9040c195ee3c9b --`
  the three production files: **exit 0**, exact original bytes.
- Active `kustomize build`: **exit 0**. Render assertions confirm no 1Password,
  bootstrap variable or installer reference.
- `kubeconform -strict -summary -ignore-missing-schemas`: **exit 0**,
  **23 resources; Valid: 10, Invalid: 0, Errors: 0, Skipped: 13**. Skipped custom
  resources are not schema-validated.
- Active/proposal YAML lint and `git diff --check`: **exit 0**.
- `make validate`: **exit 2**, repository-wide pre-existing YAML lint errors.
  Re-executed against the unchanged base using identical tooling: **exit 2** with
  byte-identical output (`cmp` exit 0). Later gates were not reached.
- Read-only remote main remains the original base; no main advance/overlap found.
- No cluster changes, secret/session database reads, image build/publication, push
  or merge were performed. Activation gates are the corrected real image digest,
  rotated real strictly scoped SealedSecret, reviewed GitOps reactivation and
  authorized target-runtime/native-authentication/recovery tests in the README.

Evidence: `$TMPDIR/source-only-make-validate.log`,
`$TMPDIR/source-only-base-make-validate.log`, `$TMPDIR/source-only-render.yaml`.

The receipts below are historical execution results from the earlier active-wiring
candidate, **not the current source-only graph or deployed state**.

## Checker follow-up: bootstrap child leakage (local correction only)

The original receipt below did **not** establish child isolation. Checker failure
was reproduced against actual installed Hermes v0.22.5 filters; removing bootstrap
from passthrough does not strip it from terminals. The follow-up correction is
tracked in [`runtime-fix/`](runtime-fix/README.md), not applied to running source.

- Added hash-locked build-time changes to the shared internal-secret predicate
  and the execute_code interpreter-probe environment (which otherwise bypassed
  the main execute_code scrub). Native `op` environment construction is unchanged.
- Actual filters, real synthetic-only Python children: unchanged source **10
  tests, 12 assertion failures**; patched scratch source **10 tests passed**.
  Covers config/skill denial, force/case variants, foreground/background terminal,
  execute_code and its interpreter probes, app profile A/B/missing scope, and
  native bootstrap retention without authenticating.
- Public pinned amd64 image source-layer hash verified; relevant source files
  match installed `/opt/hermes` byte-for-byte. Exact digests/hashes are in the
  runtime-fix README. Running-source hashes remained unchanged after testing.
- Existing suites rerun: native **5 passed** (78.144s), Whop **6 passed** (16.036s),
  Zoho **8 passed** (0.600s).
- `make validate`: **exit 2**, existing repository-wide YAML lint errors.
  Unchanged base `85a52b599b1b18b8f4b6a38cee9040c195ee3c9b` gave **byte-identical
  output** with identical tooling. Later checks were not reached.
- `docker build --network=none`: **exit 1**, Docker daemon socket unavailable.
  No corrected image was built, published or pinned. Current Deployment still
  uses vulnerable upstream source: **DO NOT ROLL OUT**, even after SealedSecret
  integration, until a corrected real image is exercised and pinned.
- Read-only remote-main check before commit still returned
  `85a52b599b1b18b8f4b6a38cee9040c195ee3c9b`; no main advance to reconcile.
- No live source/cluster mutations, credential/session DB reads, plaintext
  persistence, publication or merge. No real encrypted bootstrap/public cert
  was obtained. Native/app authentication and Kubernetes recovery remain untested.

Scratch evidence under `/opt/data/profiles/builder/cache/scratch/`:
`bootstrap-policy-red.log`, `bootstrap-policy-green.log`,
`bootstrap-policy-make-validate.log`, `bootstrap-policy-base-validate.log`.

## Original preparation execution results

- `onepassword/test_onepassword.py -v`: **5 tests passed**, including four real
  Hermes CLI reconciliation runs (seeded-empty and existing homes, twice each).
  Final rerun: `Ran 5 tests in 80.003s`, `OK`.
- Existing `mcp/whop/test_whop.py -v`: **6 passed**, including actual Hermes CLI.
- Existing `mcp/zoho/test_reconcile.py -v`: **8 passed**.
- Actual `onepassword/install_cli.py --destination <scratch>/bin --state-dir
  <scratch>/state`: downloaded the official artifact and verified the pinned
  SHA256, output `1Password CLI 2.41.1: verified and installed (linux/amd64)`.
  Running that binary with `--version` returned `2.41.1`. No authentication
  commands were run.
- `kustomize v5.7.1 build platform/hermes/onemind-hermes`: **exit 0**. Installed
  solely in scratch after checking its release archive checksum:
  `ea375e7372f9aa029129d4b2d16c66b7750b7f1213c4f66f910d981c895818d8`.
- `kubeconform v0.7.0 -strict -summary -ignore-missing-schemas` on that render:
  **exit 0**, `24 resources found; Valid: 11, Invalid: 0, Errors: 0, Skipped: 13`.
  Skipped custom resources are not schema-validated by that command.
- Changed-manifest `yamllint -c .yamllint.yml` and `git diff --check`: **exit 0**.
- `make validate`: **exit 2**, stops at repository-wide YAML lint. Running the
  same command against unchanged base main produced **byte-identical output**.
  Examples: pre-existing missing final newline in `platform/takml/ghcr.sealedsecret.yaml`,
  Matrix trailing whitespace, Grafana/Redpanda long lines and NATS flow-brace
  spacing. No unrelated lint repairs were made. Terraform and ShellCheck were not
  reached; those binaries are also absent in this worker environment.

Tests use scratch-only Python venv dependencies PyYAML 6.0.3 and yamllint 1.38.0.
The deployed Hermes CLI is v0.22.5. The actual pinned Docker image was not started:
Docker CLI exists but its daemon is unavailable. Isolated CLI tests exercise this
installed v0.22.5 runtime, not an independently launched Kubernetes container.

Scratch evidence (not committed generated output):

- `/opt/data/cache/scratch/onepassword-final-tests.log`
- `/opt/data/cache/scratch/onepassword-branch-validate.log`
- `/opt/data/cache/scratch/onepassword-baseline-validate.log`
- `/opt/data/cache/scratch/onepassword-rendered.yaml`

## Still blocked / deliberately not claimed

Real encrypted bootstrap SealedSecret integration (including the parent's public
certificate retrieval/sealing), token rotation, checker review, publication, owner
merge, Argo rollout, authenticated ClickUp access from a Hermes child, bootstrap
exclusion in the actual image, replacement-pod and isolated empty-PVC recovery.
The unresolved required Secret reference exists only in the non-active proposal
snapshot, not the production graph. Integrate real ciphertext and its Kustomize
entry only as part of a reviewed future activation change after the corrected
image gate. See README for exact name/key, image contract and certificate path.
