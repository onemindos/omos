#!/usr/bin/env bash
# Rung 1: stage Zeus's single-node 5.7 distro into the source layout the upstream
# cluster Dockerfiles expect (build-root/), so we build CURRENT 5.7 cluster images
# from the distro (Option B). Docker-independent: pure file staging + verification.
# If this prints "MISSING: 0", the docker build context is provably complete.
set -euo pipefail
DISTRO="${DISTRO:-/Users/zeus/ATAK/server/standard/takserver-docker-5.7-RELEASE-43/tak}"
HERE="$(cd "$(dirname "$0")" && pwd)"
UPSTREAM="$HERE/../upstream"
OUT="$HERE/../build-root"

rm -rf "$OUT"
mkdir -p "$OUT/takserver-core/config" "$OUT/takserver-core/API" \
         "$OUT/takserver-core/messaging" "$OUT/takserver-core/plugins" \
         "$OUT/takserver-usermanager" "$OUT/takserver-schemamanager"

miss=0
stage(){ # $1 = target rel path under build-root ; $2 = filename to find in distro, or UPSTREAM:relpath
  local target="$OUT/$1" src=""
  if [[ "$2" == UPSTREAM:* ]]; then src="$UPSTREAM/${2#UPSTREAM:}"
  else src="$(find "$DISTRO" -name "$2" 2>/dev/null | head -1)"; fi
  if [[ -n "$src" && -f "$src" ]]; then cp "$src" "$target"; echo "OK   $1"
  else echo "MISS $1  (looked for: $2)"; miss=$((miss+1)); fi
}

stage takserver-core/takserver-core.war                        takserver.war
stage takserver-core/setenv-cluster.sh                         setenv-cluster.sh
stage takserver-core/logging-restrictsize.xml                  logging-restrictsize.xml
stage takserver-core/UserAuthenticationFile.xml                UserAuthenticationFile.xml
stage takserver-core/enable_admin.sh                           enable_admin.sh
stage takserver-core/messaging-readiness.sh                    messaging-readiness.sh
stage takserver-core/config-readiness.sh                       config-readiness.sh
stage takserver-core/api-readiness.sh                          api-readiness.sh
stage takserver-core/takserver-pm.jar                          takserver-pm.jar
stage takserver-core/config/takserver-config-cluster.sh        takserver-config-cluster.sh
stage takserver-core/API/takserver-api-cluster.sh              takserver-api-cluster.sh
stage takserver-core/messaging/takserver-messaging-cluster.sh  takserver-messaging-cluster.sh
stage takserver-core/plugins/takserver-plugins-cluster.sh      takserver-plugins-cluster.sh
stage takserver-usermanager/UserManager.jar                    UserManager.jar
stage takserver-schemamanager/SchemaManager.jar                SchemaManager.jar
stage takserver-schemamanager/db-connection-configuration.sh   UPSTREAM:utilities/db-connection-configuration.sh
stage CoreConfig.xml                                           CoreConfig.example.xml

# takserver.sh: single-node launcher, NOT used by cluster entrypoints (config/api/
# messaging run their *-cluster.sh). Stub so the upstream base COPY+chmod succeeds.
printf '#!/bin/sh\n# stubbed for cluster build (cluster uses *-cluster.sh launchers)\nexit 0\n' > "$OUT/takserver-core/takserver.sh"
echo "STUB takserver-core/takserver.sh"

# NATS bridge plugin (Zeus's) into the plugins build context
if [[ -f "$HERE/../plugin/takserver-nats-plugin-1.0.0.jar" ]]; then
  echo "OK   plugin jar present"
else
  echo "NOTE plugin jar not yet copied -> cp .../tak/lib/takserver-nats-plugin-1.0.0.jar plugin/"
fi

echo ""
echo "MISSING: $miss"
[[ $miss -eq 0 ]] && echo "==> BUILD CONTEXT COMPLETE (docker build will resolve every COPY)"
