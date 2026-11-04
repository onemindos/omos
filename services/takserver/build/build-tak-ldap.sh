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

# Patch the upstream Dockerfile:
# 1. Compile a tiny LD_PRELOAD interceptor that blocks the 56 GB anonymous
#    mmap that Kerberos libraries create at startup. Without it, libkrb5 and
#    libgssapi_krb5 commit ~1 GB/s of physical pages via that mapping and
#    trigger cgroup OOMKill before slapd ever starts listening on port 389.
#    File-backed mmaps (LMDB data.mdb, ~64 MB) have fd >= 0 and are unaffected.
# 2. Set LD_PRELOAD so the interceptor is loaded before any library constructor.
# 3. Add OneMind default environment variables.
cat >> "${ORIG_DF}" <<'EXTRA'

# Build a mmap interceptor that blocks huge anonymous Kerberos allocations.
# Uses the builder's gcc — no extra packages needed in the final image.
RUN apt-get update -qq && \
    apt-get install -y --no-install-recommends gcc libc-dev && \
    printf '%s\n' \
      '#define _GNU_SOURCE' \
      '#include <sys/mman.h>' \
      '#include <dlfcn.h>' \
      '#include <errno.h>' \
      'void *mmap(void *a,size_t n,int p,int fl,int fd,off_t o){' \
      '  static void*(*r)(void*,size_t,int,int,int,off_t);' \
      '  if(!r)r=dlsym(RTLD_NEXT,"mmap");' \
      '  if(fd<0&&n>=(size_t)8<<30){errno=ENOMEM;return MAP_FAILED;}' \
      '  return r(a,n,p,fl,fd,o);' \
      '}' > /tmp/nommap_krb5.c && \
    gcc -shared -fPIC -O2 -o /usr/local/lib/nommap_krb5.so /tmp/nommap_krb5.c -ldl && \
    rm /tmp/nommap_krb5.c && \
    apt-get remove -y --purge gcc libc-dev cpp cpp-12 && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*

ENV LD_PRELOAD=/usr/local/lib/nommap_krb5.so

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
