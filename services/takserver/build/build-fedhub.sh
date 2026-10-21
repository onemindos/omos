#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ZIP_PATH="${HOME}/Downloads/takserver-fedhub-docker-5.8-RELEASE-79.zip"
IMAGE_BASE="ghcr.io/onemindos/fedhub"
TAG="5.8"

echo "==> Checking prerequisites..."
if [ ! -f "${ZIP_PATH}" ]; then
  echo "ERROR: FedHub zip not found at ${ZIP_PATH}"
  exit 1
fi

if ! docker info &>/dev/null; then
  echo "ERROR: Docker daemon is not running."
  exit 1
fi

if ! docker info 2>/dev/null | grep -q "ghcr.io"; then
  echo "WARNING: You may not be logged in to ghcr.io."
  echo "  Run: docker login ghcr.io -u <github-username> -p <GHCR_PAT>"
  echo "  Continuing anyway — push will fail if not authenticated."
fi

TMPDIR=$(mktemp -d)
trap 'echo "==> Cleaning up ${TMPDIR}..."; rm -rf "${TMPDIR}"' EXIT

echo "==> Extracting ${ZIP_PATH} to ${TMPDIR}..."
unzip -q "${ZIP_PATH}" -d "${TMPDIR}"
BUILD_CTX="${TMPDIR}/takserver-fedhub-docker-5.8-RELEASE-79"

echo "==> Writing custom Dockerfile..."
cat > "${BUILD_CTX}/Dockerfile.onemind" <<'DOCKERFILE'
FROM eclipse-temurin:17.0.19_10-jdk-ubi10-minimal

RUN microdnf update -y && \
    microdnf install -y net-tools nmap-ncat vim-minimal procps && \
    microdnf clean all && \
    useradd -m -u 1000 tak

COPY tak/ /opt/tak/

RUN chown -R tak:tak /opt/tak && \
    chmod +x /opt/tak/federation-hub/scripts/*.sh && \
    mkdir -p /opt/tak/federation-hub/logs \
             /opt/tak/federation-hub/policies \
             /opt/tak/federation-hub/jetty-temp && \
    chown -R tak:tak /opt/tak/federation-hub/logs \
                     /opt/tak/federation-hub/policies \
                     /opt/tak/federation-hub/jetty-temp

USER tak

EXPOSE 9100 9101 9102

ENTRYPOINT ["/bin/bash", "-c", "/opt/tak/federation-hub/scripts/configureInDocker.sh init"]
DOCKERFILE

echo "==> Building Docker image ${IMAGE_BASE}:${TAG}..."
docker build \
  -t "${IMAGE_BASE}:${TAG}" \
  -t "${IMAGE_BASE}:latest" \
  -f "${BUILD_CTX}/Dockerfile.onemind" \
  "${BUILD_CTX}"

echo "==> Pushing ${IMAGE_BASE}:${TAG}..."
docker push "${IMAGE_BASE}:${TAG}"
docker push "${IMAGE_BASE}:latest"

echo ""
echo "SUCCESS: ${IMAGE_BASE}:${TAG} pushed to GHCR."
echo ""
echo "Next steps:"
echo "  1. Generate FedHub certs (see services/takserver/build/README-BUILDS.md)"
echo "  2. Seal certs as K8s Secret and apply platform/takserver/fedhub/ manifests"
echo "  3. kubectl -n takserver apply -f platform/takserver/fedhub/"
