# Rolling a New GHCR Image onto the Cluster

Taking a new `ghcr.io/onemindos/*` image from nothing to a running pod. Every step
here has failed at least once.

## 1. Build in CI

There is no docker daemon on the agent pod and no in-cluster registry. Build with a
GitHub Actions workflow (`.github/workflows/build-<name>.yml`) using
`docker/build-push-action` and the `GITHUB_TOKEN` (`permissions: contents: read,
packages: write`). Keep `workflow_dispatch` so a build can be re-run without a commit.
Pin actions by commit SHA like the existing workflows.

- Tag with a version or digest, never `:latest`, and pin the running manifest to it.
- `pip --index-url` applies to the whole invocation: a CPU-only PyTorch index hosts
  only torch wheels, so install other packages in a separate `pip install`.
- Bake dependencies into the image. A ConfigMap-as-source plus pip-at-boot init
  container re-downloads everything on every start and runs unpinned code.

## 2. New packages are private

A newly pushed GHCR package is private, so pulls fail with `ImagePullBackOff` and a
401 until the namespace has a pull secret. The convention here: every namespace that
pulls our images carries a committed `ghcr` SealedSecret
(`platform/<ns>/ghcr.sealedsecret.yaml`), and the pod spec has
`imagePullSecrets: [{name: ghcr}]`.

## 3. Sealing a pull secret (owner-led)

Creating the SealedSecret needs the cluster's public certificate and the pull
credential, which this agent should not handle. Ask the owner to run the standard
flow in `CLAUDE.md` ("Secrets") with `--scope namespace-wide`. A SealedSecret sealed
for one namespace cannot be relabeled for another; it is re-sealed from the source.
Do not `kubectl create secret` and leave it: Argo prunes it on the next sync.

## 4. Verify

1. The workflow run is green (`gh run list --workflow build-<name>.yml`).
2. The SealedSecret decrypted: `kubectl -n <ns> get secret ghcr`.
3. The pod leaves ImagePullBackOff, passes ContainerCreating (a large image takes a
   few minutes to pull, which is progress, not a stall) and reaches Running.
4. The app starts in seconds from the baked image. A start that takes minutes means
   it still installs dependencies at runtime.
