#!/usr/bin/env bash
set -euo pipefail

TAK_ADMIN_P12="${TAK_ADMIN_P12:-/Users/zeus/ATAK/server/standard/takserver-docker-5.7-RELEASE-43/tak/certs/files/admin.p12}"
: "${TAK_ADMIN_P12_PASS:?Set TAK_ADMIN_P12_PASS first}"

/usr/bin/curl -sk --cert-type P12 --cert "$TAK_ADMIN_P12:$TAK_ADMIN_P12_PASS" \
  https://localhost:8446/Marti/api/clientEndPoints | python3 -m json.tool
