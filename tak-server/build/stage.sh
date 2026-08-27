#!/usr/bin/env bash
# Verify the 5.8 distro has all files required by the upstream Dockerfiles.
# 5.8 ships pre-assembled in tak/ — no staging needed. This script just confirms
# the build context is complete before running docker build.
#
# Usage:
#   DISTRO=/path/to/takserver-docker-hardened-5.8-RELEASE-69 ./build/stage.sh
set -euo pipefail

DISTRO="${DISTRO:-/Users/zeus/OneMind Guardian (TAK)/server/5.8/takserver-docker-hardened-5.8-extracted/takserver-docker-hardened-5.8-RELEASE-69}"
HERE="$(cd "$(dirname "$0")" && pwd)"

miss=0
check(){ # $1 = relative path under distro
  if [[ -f "$DISTRO/$1" ]]; then echo "OK   $1"
  else echo "MISS $1"; miss=$((miss+1)); fi
}

echo "=== Verifying 5.8 distro: $DISTRO ==="
check tak/takserver.war
check tak/takserver-pm.jar
check tak/takserver-plugins-cluster.sh
check tak/setenv.sh
check tak/setenv-cluster.sh
check tak/config/takserver-config-cluster.sh
check tak/API/takserver-api-cluster.sh
check tak/messaging/takserver-messaging-cluster.sh
check tak/api-readiness.sh
check tak/config-readiness.sh
check tak/messaging-readiness.sh
check tak/db-utils/SchemaManager.jar
check tak/db-utils/generic-cluster-database-configuration.sh
check tak/utils/UserManager.jar

echo ""
echo "=== Plugin JAR ==="
if [[ -f "$HERE/../plugin/takserver-nats-plugin-1.0.0.jar" ]]; then
  echo "OK   plugin/takserver-nats-plugin-1.0.0.jar"
else
  echo "NOTE plugin jar not present — rebuild against 5.8 SDK first:"
  echo "     cd plugin-src && gradle build && cp build/libs/*.jar ../plugin/"
fi

echo ""
echo "MISSING: $miss"
[[ $miss -eq 0 ]] && echo "==> DISTRO COMPLETE — ready to build images"

echo ""
echo "=== Build commands (run from $DISTRO) ==="
GHCR="ghcr.io/onemindos"
HERE_REL="$(cd "$HERE/.." && pwd)"
cat <<EOF
cd "$DISTRO"
docker build --platform linux/amd64 -f "$HERE_REL/docker/Dockerfile.takserver-base"     -t $GHCR/takserver-base:5.8-onemind         .
docker build --platform linux/amd64 -f "$HERE_REL/docker/Dockerfile.takserver-config"   -t $GHCR/takserver-config:5.8-onemind        --build-arg TAKSERVER_IMAGE_REPO=$GHCR/takserver-base --build-arg TAKSERVER_IMAGE_TAG=5.8-onemind .
docker build --platform linux/amd64 -f "$HERE_REL/docker/Dockerfile.takserver-api"      -t $GHCR/takserver-api:5.8-onemind           --build-arg TAKSERVER_IMAGE_REPO=$GHCR/takserver-base --build-arg TAKSERVER_IMAGE_TAG=5.8-onemind .
docker build --platform linux/amd64 -f "$HERE_REL/docker/Dockerfile.takserver-messaging" -t $GHCR/takserver-messaging:5.8-onemind    --build-arg TAKSERVER_IMAGE_REPO=$GHCR/takserver-base --build-arg TAKSERVER_IMAGE_TAG=5.8-onemind .
docker build --platform linux/amd64 -f "$HERE_REL/docker/Dockerfile.takserver-plugins"  -t $GHCR/takserver-plugins:5.8-onemind       --build-arg TAKSERVER_IMAGE_REPO=$GHCR/takserver-base --build-arg TAKSERVER_IMAGE_TAG=5.8-onemind .
docker build --platform linux/amd64 -f "$HERE_REL/docker/Dockerfile.database-setup"     -t $GHCR/takserver-database-setup:5.8-onemind .
# Push all upstream images first, then build the OneMind plugins overlay:
docker push $GHCR/takserver-base:5.8-onemind
docker push $GHCR/takserver-config:5.8-onemind
docker push $GHCR/takserver-api:5.8-onemind
docker push $GHCR/takserver-messaging:5.8-onemind
docker push $GHCR/takserver-plugins:5.8-onemind
docker push $GHCR/takserver-database-setup:5.8-onemind
# OneMind plugins overlay (adds NatsBridgePlugin on top of takserver-plugins):
cd "$HERE_REL"
docker build --platform linux/amd64 -f docker/Dockerfile.takserver-plugins.onemind -t $GHCR/takserver-plugins:5.8-onemind .
docker push $GHCR/takserver-plugins:5.8-onemind
EOF
