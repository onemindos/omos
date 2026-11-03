# Bootstrap child-environment correction (local source; rollout BLOCKED)

## Root cause and smallest correction

Removing `OP_SERVICE_ACCOUNT_TOKEN` from `terminal.env_passthrough` does **not**
strip it from v0.22.5 terminals. `_filter_secret_env` otherwise copies names absent
from Hermes's managed-credential policy. A plain static provider-blocklist entry
would still permit `_HERMES_FORCE_` overrides and credential-inheriting children.

`apply_source_patch.py` adds one exact, case-folded name to
`_is_hermes_internal_secret`. Existing shared code then denies it in foreground
and background/PTY terminal filtering, skill/config passthrough acceptance and
scoped additions, execute_code passthrough, and `hermes_subprocess_env` even with
`inherit_credentials=True`. The native source's `_op_child_env` deliberately does
not use that generic scrub: the main process still supplies bootstrap only to
its minimal `op read` environment. No resolver, broad vault API, runtime monkey
patch, plaintext persistence or retained-PVC source overlay is introduced.

An additional execute_code path, `_probe_python`, spawned interpreter probes
with an unsanitized inherited env. The second small source edit routes these
through `hermes_subprocess_env()` too; interpreter usability/prefix probes are
regression-tested.

The patcher requires the **entire original SHA256 of both edited files**, checks
each replacement occurs once and compiles both before writing either. It refuses
unknown/already patched files. Source upgrades require a fresh review/hash, not fuzzy application.
The patch does not change `CLICKUP_API_TOKEN` registration or profile resolution.
Terminal's existing handling of other non-managed application env names is not
changed; this is not a general terminal allowlist redesign.

## Exact pinned-image assessment

Read-only Docker Hub registry inspection resolved the manifest-list digest used
in `deployment.yaml`:

- index: `sha256:9c5561579f3ffdf333401797094ac3a743738d6c25d40f4331c7bf77b694238d`
- linux/amd64 manifest: `sha256:fd822bae0a61bf023450737edbeaeea5fc6460f17cb4ee240a7e67b306ed63fb`
- public source layer: `sha256:fcf720a342af7b5020f7a08992d0fd6437bd319f66b1cd3a76ad735c06a04899`
- public source install stamp revision: `bf0203a794a07aac3ff764972ab9087baa4e4f4e`

Downloaded the source layer, verified its SHA256 before opening it, and compared
these extracted public files byte-for-byte with `/opt/hermes` (all identical):

| File | SHA256 |
| --- | --- |
| `tools/environments/local_env_policy.py` | `4429af03406ead5e79e5cfe711787030f40d42bf50c6253b4abc6ee0154a77d4` |
| `tools/environments/local.py` | `35aa5b88d5252a1fceb8a1b13537069e75bd41d3a7e4af5ba02216b298275169` |
| `tools/code_execution_env.py` | `1c6e0c42a7f4c59f0c65a02aec397c92c87860226ab181de64751e00fecf2785` |
| `tools/env_passthrough.py` | `0acb584cfc8e336bac794991829a10fdee5e70dcb3f3430219e19af57a4ab61f` |
| `agent/secret_sources/onepassword.py` | `765fefb33be2e6eeb667f837c58c5fc66e09cf5045e0151194e214251401fdc3` |

This establishes relevant-source compatibility, **not** an independent full image
execution or verification of all later layer contents. Only amd64 is assessed,
matching the workload's node selector. Do not infer compatibility from the version
string alone.

## Reproduction and build gate

Use a new interpreter with the installed Hermes dependencies. The test harness
clears inherited environment before importing Hermes and uses a fresh scratch home
with synthetic marker values only. It imports the **actual** local and execute_code
filters; `--policy-root` substitutes the scratch policy and execute_code env
modules before those imports. There are no mocked authentication responses or credential lookups.

```sh
fix=platform/hermes/onemind-hermes/onepassword/runtime-fix
# First run without --policy-root: expected RED on unchanged v0.22.5.
PYTHONPATH=/opt/hermes:/opt/hermes/.venv/lib/python3.14/site-packages \
  /opt/hermes/.venv/bin/python "$fix/test_child_policy.py" \
  --runtime-root /opt/hermes -v
# Copy only public policy source into a NEW scratch root, then patch that copy.
mkdir -p "$TMPDIR/policy-copy/tools/environments"
cp /opt/hermes/tools/environments/local_env_policy.py \
  "$TMPDIR/policy-copy/tools/environments/local_env_policy.py"
cp /opt/hermes/tools/code_execution_env.py \
  "$TMPDIR/policy-copy/tools/code_execution_env.py"
python3 "$fix/apply_source_patch.py" --source-root "$TMPDIR/policy-copy"
PYTHONPATH=/opt/hermes:/opt/hermes/.venv/lib/python3.14/site-packages \
  /opt/hermes/.venv/bin/python "$fix/test_child_policy.py" \
  --runtime-root /opt/hermes --policy-root "$TMPDIR/policy-copy" -v
```

Actual results: unchanged runtime **10 tests, 12 assertion failures**; patched
scratch sources **10 tests passed**. Tests spawn real Python subprocesses with the
filtered environments, asserting bootstrap absence and intended app-marker
presence; cover force/case variants, hostile skill/config registration, scoped
app A/B/missing values including scope-only additions, and native minimal-env
bootstrap retention. Native retention is an environment-construction test, not
an authenticated vault read or successful app authentication.

`Dockerfile` is a local build recipe using the existing digest, hash-locked patch
and same synthetic regression suite as a build gate. It inherits upstream s6's
root entrypoint and UID drop. No image name/digest is invented or wired into the
Deployment. The actual local build attempt:

```sh
docker build --network=none -t hermes-bootstrap-policy-local "$fix"
```

failed: **Cannot connect to the Docker daemon at unix:///var/run/docker.sock**.
No image was built or published. An authorized build environment must build and
exercise this recipe, review/scan the result, then publish and pin its **real**
digest for the resolving gateway before any bootstrap injection/rollout. Generic
child Hermes invocations will no longer inherit bootstrap for fresh resolution;
only the designated resolving main process is provisioned. Separately scheduled
resolvers need explicit reviewed bootstrap provisioning, not terminal passthrough.

## Limits / remaining gates

- The current Deployment still points to the unchanged upstream image with the
  known child-filter defect. This correction is **NOT DEPLOYED**. The source-only
  draft removes all new bootstrap/source activation from active manifests; the
  non-active proposal must not be enabled until the corrected-image and real
  SealedSecret gates pass.
- This exact-name fix covers the branch's fixed `OP_SERVICE_ACCOUNT_TOKEN`, not
  arbitrary alternate token variable names, interactive `OP_SESSION_*` or Connect
  credentials. Do not change authentication mode/name without extending policy/tests.
- Env scrubbing is not an OS sandbox: code with same-UID access may inspect process
  memory/procfs, read credential files if present, or deliberately construct new env
  values. Generic tool inheritance is fixed; hostile-code vault isolation requires
  a separate process/security boundary. No such broader guarantee is claimed.
- No source under running `/opt/hermes` was modified. No session DB, credential
  files, real secret values or cluster state were read. No deployment/publication.
- Real bootstrap SealedSecret/public certificate, rotation, target-image execution,
  native authentication, replacement-pod and empty-PVC checks remain blocked.

## Tracked manual CI build (not run by this change)

`.github/workflows/build-hermes-runtime.yml` is **workflow_dispatch only**. After
separate authorized publication makes it available to Actions, dispatch the
reviewed branch/ref (no arbitrary checkout-ref input). It uses standard
`ubuntu-24.04` hosted capacity; no runner purchase or infrastructure changes.
A private repository may consume its existing Actions allowance; confirm that
allowance with the owner rather than purchasing capacity if execution is blocked.

Actions are commit-pinned; checkout does not persist credentials. Permissions are
only `contents: read` and `packages: write`; the only credential is the ephemeral
`GITHUB_TOKEN` for GHCR login **after tests**. There is no 1Password credential in
CI. The narrow build context contains only public correction/test source. The
Dockerfile's synthetic gate must pass, then CI runs the loaded linux/amd64 image
as UID/GID 10000 with read-only root, writable scratch tmpfs, no network and no
injected credentials. This bypasses the s6 entrypoint to test the actual filters;
it does not establish full gateway startup or authenticated native resolution.

CI pushes that same tested local image (not a second rebuild) to
`ghcr.io/onemindos/hermes-agent:sha-<full-commit>-run-<run-id>-<attempt>`.
Run/attempt suffixes prevent ordinary reruns from overwriting earlier tags; GHCR
tags are not registry-enforced immutable, so activation must use the real digest.
CI obtains the pushed RepoDigest, reads that exact reference back from GHCR and
records it in job outputs, the summary and `hermes-runtime-image.txt` artifact.
No digest exists until this succeeds. Review/scan the image and retain its build
record before using [proposal preparation](../proposal/README.md).

Local contract checks (no build/publication):

```sh
python -m unittest discover -s platform/hermes/onemind-hermes/onepassword -p 'test_*.py' -v
yamllint -c .yamllint.yml .github/workflows/build-hermes-runtime.yml
```
