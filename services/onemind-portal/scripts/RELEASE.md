# Portable portal image release

The production Dockerfile performs immutable Yarn install, TypeScript checking,
pinned TechDocs generation and backend compilation in Linux. Its final stage
checks the common directory JSON and TechDocs files are actually in the image.
GitHub Actions uses this same Dockerfile; it is not the only build executor.

From a committed, clean portal checkout with a running Docker/BuildKit engine:

    services/onemind-portal/scripts/release-image.sh

This builds linux/amd64 and loads it locally, without publishing or changing Git
or Kubernetes. Host-native node_modules and generated bundles are never reused.
JavaScript is compiled on the Linux builder architecture, avoiding AMD64 Go/esbuild
emulation crashes on ARM Macs. The final stage remains linux/amd64 and installs
its own native dependencies; an in-memory SQLite query checks the target binary.
Both common-package JSON exports are required from the final image.
The source context is an archive of tracked files from the exact commit, not a copy
of your working directory, so local credentials and untracked files stay out.
Commit the intended changes before running; dirty portal/workflow changes fail closed.

After authenticating Docker to the existing approved registry, publish explicitly:

    services/onemind-portal/scripts/release-image.sh --push --output-dir /absolute/new-release-directory

The default registry is ghcr.io/onemindos/onemind-portal. IMAGE_REPOSITORY can select
an operator-approved registry when intentionally changing destinations. No login,
credential creation, billing change, provisioning or deployment is automated here.
A remote Linux builder can run the same command independently of GitHub Actions;
do not put heavy image builds on the production Kubernetes nodes.

A successful push verifies the immutable remote digest before emitting
image-digest.txt and release.json (commit, image, platform, gates and digest).
Use a fresh output directory. A failed build or remote verification emits no new
deployment digest. Without --push, no deployable digest file is emitted.

Next, review the release record, run the image's runtime checks, and open a separate
GitOps PR pinning that verified digest in platform/backstage/deployment.yaml.
Argo CD owns rollout. Verify login, existing home/catalog/MCP permissions, 2/2
replicas and new behavior. This script never merges a PR,
patches a Deployment or marks unavailable GitHub checks green.

Contract tests (fake Docker, no network or production changes):

    python3 services/onemind-portal/scripts/test-release-image.py

Real Linux image build, registry access and runtime checks are additional gates;
contract-test success alone does not prove a release image works.
