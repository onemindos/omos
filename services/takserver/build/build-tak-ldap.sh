#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ZIP_PATH="${HOME}/Downloads/takserver-ldap-5.8-RELEASE-79.zip"
IMAGE_BASE="ghcr.io/onemindos/tak-ldap"
TAG="5.8"

echo "==> Checking prerequisites..."
if [ ! -f "${ZIP_PATH}" ]; then
  echo "ERROR: TAK LDAP zip not found at ${ZIP_PATH}"
  exit 1
fi

if ! docker info &>/dev/null; then
  echo "ERROR: Docker daemon is not running."
  exit 1
fi

if ! docker info 2>/dev/null | grep -q "ghcr.io"; then
  echo "WARNING: You may not be logged in to ghcr.io."
  echo "  Run: docker login ghcr.io -u <github-username> -p <GHCR_PAT>"
fi

TMPDIR=$(mktemp -d)
trap 'echo "==> Cleaning up ${TMPDIR}..."; rm -rf "${TMPDIR}"' EXIT

echo "==> Extracting ${ZIP_PATH} to ${TMPDIR}..."
unzip -q "${ZIP_PATH}" -d "${TMPDIR}"
BUILD_CTX="${TMPDIR}/takserver-ldap-5.8-RELEASE-79"

# Patch the upstream Dockerfile to use our base DN defaults.
# The original Dockerfile.takserver-ldap already has the correct structure;
# we overlay ENV defaults so the K8s manifest can override via envFrom.
echo "==> Patching Dockerfile for GHCR target..."
# The original references slapd / bootstrap.sh / ldifs — keep those intact,
# just set our default base DN and expose both LDAP ports.
ORIG_DF="${BUILD_CTX}/Dockerfile.takserver-ldap"

if [ ! -f "${ORIG_DF}" ]; then
  echo "ERROR: Dockerfile.takserver-ldap not found inside zip at expected path."
  echo "  Expected: ${ORIG_DF}"
  exit 1
fi

# Append our environment defaults to the upstream Dockerfile.
cat >> "${ORIG_DF}" <<'EXTRA'

# OneMind overrides — K8s manifest can override these via ConfigMap envFrom
ENV LDAP_BASE_DN="dc=onemindos,dc=dev"
ENV LDAP_ADMIN="ldapadm"
EXTRA

echo "==> Building Docker image ${IMAGE_BASE}:${TAG}..."
docker build \
  -t "${IMAGE_BASE}:${TAG}" \
  -t "${IMAGE_BASE}:latest" \
  -f "${ORIG_DF}" \
  "${BUILD_CTX}"

echo "==> Pushing ${IMAGE_BASE}:${TAG}..."
docker push "${IMAGE_BASE}:${TAG}"
docker push "${IMAGE_BASE}:latest"

echo ""
echo "SUCCESS: ${IMAGE_BASE}:${TAG} pushed to GHCR."
echo ""
echo "Next steps:"
echo "  1. Create the LDAP admin password SealedSecret:"
echo "       kubectl -n takserver create secret generic tak-ldap-creds \\"
echo "         --from-literal=LDAP_ADMIN_PASSWORD='<strong-password>' \\"
echo "         --dry-run=client -o yaml | kubeseal -o yaml > \\"
echo "         platform/takserver/deploy/tak-ldap-creds.sealedsecret.yaml"
echo "  2. kubectl -n takserver apply -f platform/takserver/deploy/tak-ldap.yaml"
