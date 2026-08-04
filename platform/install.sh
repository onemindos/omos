#!/usr/bin/env bash
# =============================================================================
# OneMind Hub — Platform Install Script
# Reproduces the full DOKS cloud hub from omos-infra git alone.
# Run from repo root. Requires: helm, kubectl, doctl (DO CLI), op (1Password CLI)
# =============================================================================
set -euo pipefail

CTX="${KUBE_CONTEXT:-tailscale-operator.tail717ba0.ts.net}"
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PLATFORM="$REPO_ROOT/platform"

echo "==> Using cluster context: $CTX"

# -----------------------------------------------------------------------------
# 0. EXTERNAL SECRETS OPERATOR + 1PASSWORD CONNECT
# All K8s secrets are synced from 1Password via ESO — no manual kubectl create secret needed.
# Requires two bootstrap secrets (only needed once on fresh cluster):
#   op-credentials    — 1password-credentials.json file (from 1Password Connect server setup)
#   onepassword-connect-token  — Connect API JWT token
#
# Bootstrap (one-time, from 1Password vault "00-24 SO"):
#   kubectl create namespace external-secrets
#   kubectl create secret generic op-credentials \
#     -n external-secrets --from-file=1password-credentials.json=<path-to-file>
#   kubectl create secret generic onepassword-connect-token \
#     -n external-secrets --from-literal=token=<connect-jwt>
# -----------------------------------------------------------------------------
helm repo add external-secrets https://charts.external-secrets.io --force-update
helm upgrade --install external-secrets external-secrets/external-secrets \
  --kube-context "$CTX" -n external-secrets --create-namespace \
  --set installCRDs=true

helm repo add 1password https://1password.github.io/connect-helm-charts --force-update
helm upgrade --install onepassword-connect 1password/connect \
  --kube-context "$CTX" -n external-secrets \
  --set connect.credentials_base64="" \
  --set operator.create=false
# Note: helm chart reads op-credentials secret created above

kubectl --context "$CTX" apply -f "$PLATFORM/external-secrets/cluster-secret-store.yaml"

# Wait for Connect pod to be ready before services that depend on synced secrets
kubectl --context "$CTX" wait --for=condition=available deployment/onepassword-connect \
  -n external-secrets --timeout=120s

# Apply all ExternalSecrets — they will sync automatically
kubectl --context "$CTX" apply -f "$PLATFORM/legacy/external-secret.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/external-secret.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/wwv/external-secret.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos/external-secret.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/orp/external-secret.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/hermes-webui/external-secret.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudflare/external-secret.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/ingress/external-secret.yaml"
kubectl --context "$CTX" apply -f "$REPO_ROOT/tak-cluster/deploy/external-secret.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/authentik/external-secret.yaml"

# Wait for secrets to sync before deploying services
sleep 15

# -----------------------------------------------------------------------------
# 1. CORE PLATFORM — cert-manager, ingress-nginx, ArgoCD
# -----------------------------------------------------------------------------
helm repo add jetstack https://charts.jetstack.io --force-update
helm upgrade --install cert-manager jetstack/cert-manager \
  --kube-context "$CTX" -n cert-manager --create-namespace \
  -f "$PLATFORM/values/cert-manager.yaml"

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx --force-update
helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --kube-context "$CTX" -n ingress-nginx --create-namespace \
  -f "$PLATFORM/values/ingress-nginx.yaml"

helm repo add argo https://argoproj.github.io/argo-helm --force-update
helm upgrade --install argocd argo/argo-cd \
  --kube-context "$CTX" -n argocd --create-namespace
# ArgoCD uses all defaults — no custom values needed

# -----------------------------------------------------------------------------
# 2. TLS ISSUERS (cert-manager ClusterIssuers)
# cloudflare-api-token secret is synced by ESO (step 0) — no manual action needed
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/ingress/issuer-dns.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/ingress/issuers.yaml"

# -----------------------------------------------------------------------------
# 3. CLOUDFLARE TUNNEL (zero public inbound)
# Routes are declared in platform/cloudflare/config.yaml — edit + apply to add new routes.
# Requires tunnel-creds-json secret (NOT synced by ESO — derived from TUNNEL_TOKEN, one-time):
#   TOKEN=$(op read "op://00-24 SO/tunnel-credentials/tunnel-token")
#   python3 -c "
#     import sys,json,base64; t='$TOKEN'; p=t+'='*(4-len(t)%4)
#     d=json.loads(base64.urlsafe_b64decode(p))
#     print(json.dumps({'AccountTag':d['a'],'TunnelID':d['t'],'TunnelSecret':d['s']}))
#   " | kubectl create secret generic tunnel-creds-json -n cloudflare --from-file=credentials.json=/dev/stdin
# -----------------------------------------------------------------------------
kubectl --context "$CTX" create namespace cloudflare --dry-run=client -o yaml | kubectl --context "$CTX" apply -f -
kubectl --context "$CTX" apply -f "$PLATFORM/cloudflare/config.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudflare/cloudflared.yaml"

# -----------------------------------------------------------------------------
# 4. NATS — 3-node HA JetStream cluster
# Requires: nats/values.secret.yaml (has ONEMIND/SYS passwords — inject from 1Password)
# -----------------------------------------------------------------------------
helm repo add nats https://nats-io.github.io/k8s/helm/charts --force-update
kubectl --context "$CTX" create namespace nats --dry-run=client -o yaml | kubectl --context "$CTX" apply -f -
kubectl --context "$CTX" apply -f "$PLATFORM/nats/certs.yaml"
helm upgrade --install nats nats/nats \
  --kube-context "$CTX" -n nats \
  -f "$PLATFORM/nats/values.secret.yaml"  # inject from 1Password on fresh deploy
kubectl --context "$CTX" apply -f "$PLATFORM/nats/omos-nats-tailnet-service.yaml"
# Note: NATS WebSocket must be enabled in nats/values.secret.yaml:
#   websocket:
#     enabled: true
#     port: 8080
# The nats-ws-route ConfigMap (platform/omos-ui/cloudflare-route.yaml)
# exposes it as wss://nats-ws.onemindos.dev via Cloudflare tunnel

# -----------------------------------------------------------------------------
# 4b. NACK — JetStream CRD controller (Stream + KeyValue declarative mgmt)
# Requires: nack-values.secret.yaml (gitignored) — see nack-values.template.yaml
# Generate:
#   ONEMIND_PW=$(grep password "$PLATFORM/nats/values.secret.yaml" | awk '{print $2}' | head -1)
#   # or: ONEMIND_PW=$(op read "op://00-24 SO/nats-onemind/password")
#   cat > "$PLATFORM/nats/nack-values.secret.yaml" <<EOF
#   jetstream:
#     nats:
#       url: "nats://onemind:${ONEMIND_PW}@nats.nats.svc.cluster.local:4222"
#   EOF
# -----------------------------------------------------------------------------
helm upgrade --install nack nats/nack \
  --kube-context "$CTX" -n nats \
  -f "$PLATFORM/nats/nack-values.yaml" \
  -f "$PLATFORM/nats/nack-values.secret.yaml"  # inject ONEMIND_PW (see template)

# Apply all Stream + KeyValue CRD manifests — controller reconciles them to NATS
kubectl --context "$CTX" apply -f "$PLATFORM/nats/streams/"

# -----------------------------------------------------------------------------
# 5. TAILSCALE OPERATOR (kubectl-over-tailnet + service exposure)
# Requires: tailscale-operator.yaml with real clientSecret from 1Password
# op read "op://00-24 SO/tailscale-oauth-operator/credential"
# -----------------------------------------------------------------------------
helm repo add tailscale https://pkgs.tailscale.com/helmcharts --force-update
helm upgrade --install tailscale-operator tailscale/tailscale-operator \
  --kube-context "$CTX" -n tailscale --create-namespace \
  -f "$PLATFORM/values/tailscale-operator.yaml"  # clientSecret injected at deploy time
kubectl --context "$CTX" apply -f "$PLATFORM/tailscale/apiserver-rbac.yaml"

# -----------------------------------------------------------------------------
# 6. OBSERVABILITY — kube-prometheus-stack, Loki, Tempo, OTel Collector
# -----------------------------------------------------------------------------
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts --force-update
helm upgrade --install monitoring prometheus-community/kube-prometheus-stack \
  --kube-context "$CTX" -n monitoring --create-namespace \
  -f "$PLATFORM/values/kube-prometheus-stack.yaml"

helm repo add grafana https://grafana.github.io/helm-charts --force-update
helm upgrade --install loki grafana/loki \
  --kube-context "$CTX" -n monitoring \
  -f "$PLATFORM/values/loki.yaml"

helm upgrade --install tempo grafana/tempo \
  --kube-context "$CTX" -n monitoring \
  -f "$PLATFORM/values/tempo.yaml"

helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts --force-update
helm upgrade --install otel-collector open-telemetry/opentelemetry-collector \
  --kube-context "$CTX" -n monitoring \
  -f "$PLATFORM/values/otel-collector.yaml"

# Apply custom datasource ConfigMap (Loki + Tempo linked in Grafana)
kubectl --context "$CTX" apply -f "$PLATFORM/grafana-datasources.yaml"

# -----------------------------------------------------------------------------
# 6b. NATS SURVEYOR — Prometheus exporter for the NATS control plane
# Requires: nats-server-tls secret (CA cert from cert-manager) and $SYS creds
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/nats-surveyor/deployment.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/nats-surveyor/servicemonitor.yaml"

# -----------------------------------------------------------------------------
# 7. ORP — Open Reality Protocol fusion engine
# Secrets (orp-nats-creds, orp-bridge-key) synced from 1Password by ESO (step 0)
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/orp/namespace.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/orp/pvc.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/orp/configmap.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/orp/bridge-script.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/orp/deployment.yaml"

# -----------------------------------------------------------------------------
# 8. OMOS SERVICES — shared omos namespace (tak-projector, feed-world, etc.)
# Secrets (omos-nats-creds) synced from 1Password by ESO (step 0)
# -----------------------------------------------------------------------------
# -----------------------------------------------------------------------------
# 7b. AUTHENTIK — IdP for onemindos.dev, Forward Auth for ATOC + ATLAS
# Secrets (authentik-env) synced from 1Password by ESO (step 0)
# Bootstrap DB (one-time):
#   kubectl exec -n cloudtak deploy/cloudtak-postgis -- \
#     psql -U postgres -c "CREATE USER authentik WITH PASSWORD '<from 1Password>';"
#   kubectl exec -n cloudtak deploy/cloudtak-postgis -- \
#     psql -U postgres -c "CREATE DATABASE authentik OWNER authentik;"
# After deploy: visit https://auth.onemindos.dev/if/flow/initial-setup/
# to set the admin password, then configure the Embedded Outpost URL:
#   Authentik Admin → System → Outposts → embedded → Edit
#   Set Authentik URL: https://auth.onemindos.dev
# -----------------------------------------------------------------------------
helm repo add authentik https://charts.goauthentik.io --force-update
kubectl --context "$CTX" apply -f "$PLATFORM/authentik/namespace.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/authentik/external-secret.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/authentik/cloudflare-route.yaml"

# Wait for secret to sync before Helm reads it
kubectl --context "$CTX" wait --for=condition=Ready externalsecret/authentik-env \
  -n authentik --timeout=60s 2>/dev/null || true

helm upgrade --install authentik authentik/authentik \
  --kube-context "$CTX" -n authentik \
  -f "$PLATFORM/authentik/helm-values.yaml"

kubectl --context "$CTX" apply -f "$PLATFORM/omos/namespace.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos/tak-projector.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos-ui/deployment.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos-ui/cloudflare-route.yaml"
# Note: Cesium ion vars (VITE_CESIUM_ION_TOKEN, VITE_CESIUM_ION_SERVER) are baked
# at build time — set them as GitHub Actions secrets in the omos-ui repo, not K8s secrets.
# Self-hosted ion: VITE_CESIUM_ION_SERVER=https://ion.onemindos.dev
# Cesium.com cloud: leave VITE_CESIUM_ION_SERVER empty, set token from cesium.com/ion
kubectl --context "$CTX" apply -f "$PLATFORM/omos-mission-engine/external-secret.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos-mission-engine/deployment.yaml"

# -----------------------------------------------------------------------------
# 8b. OMOS WORLD FEEDS — live global data → NATS → TAK → OMOS ATOC map
# Sources: AIS ships, ADS-B aircraft, APRS trackers, satellites, hazards
# Requires:
#   omos-nats-creds secret in omos namespace (already created above)
#   omos-feed-secrets (AISSTREAM_API_KEY, OPENSKY_USERNAME, OPENSKY_PASSWORD):
#     kubectl create secret generic omos-feed-secrets -n omos \
#       --from-literal=AISSTREAM_API_KEY='<key from aisstream.io — free signup>' \
#       --from-literal=OPENSKY_USERNAME='<opensky.org free account>' \
#       --from-literal=OPENSKY_PASSWORD='<opensky password>'
# Toggle feeds: kubectl scale deployment feed-ships feed-aircraft feed-aprs \
#   feed-satellites feed-hazards -n omos --replicas=0|1
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/omos/feed-scripts.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos/feeds.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos/feed-mta-script.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos/feed-mta.yaml"

# -----------------------------------------------------------------------------
# 9. LEGACY AGENT (Hermes — mother) + TAK-AGENT + OneMind AI Bridge
# Secrets (hermes-env, hermes-leaf-conf, hermes-nats-ca) synced by ESO (step 0)
# cloudtak-agent-token: still out-of-band — generate in CloudTAK UI → Settings → Tokens
#   kubectl create secret generic cloudtak-agent-token -n hermes \
#     --from-literal=token='<token from CloudTAK UI: Settings → Tokens → New>'
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/legacy/namespace.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/legacy/pvc.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/legacy/deployment.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/legacy/tak-agent.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/hermes/onemind-ai-bridge.yaml"

# -----------------------------------------------------------------------------
# 10. CLOUDTAK — API + Events + Media + Tiles + MinIO + PostGIS
# Secrets (cloudtak-secrets) synced from 1Password by ESO (step 0)
# -----------------------------------------------------------------------------
kubectl --context "$CTX" create namespace cloudtak --dry-run=client -o yaml | kubectl --context "$CTX" apply -f -
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/api.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/events.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/media.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/tiles.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/minio.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/postgis.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/cloudflare-route.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/ingress.yaml"

# -----------------------------------------------------------------------------
# 11. Synadia Bridge (Legacy AI full-screen chat UI at legacy-ai.onemindos.dev)
# Requires: omos-nats-creds in synadia-bridge namespace + ghcr-onemindos pull secret
#   kubectl create secret generic omos-nats-creds -n synadia-bridge \
#     --from-literal=user='...' --from-literal=password='...'
#   kubectl create secret docker-registry ghcr-onemindos -n synadia-bridge \
#     --docker-server=ghcr.io --docker-username=zeus --docker-password='<PAT>'
# After deploy: CloudTAK Admin → Config → External Applications → legacy-ai.onemindos.dev
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/synadia-bridge/deployment.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/synadia-bridge/cloudflare-route.yaml"

echo ""
echo "==> Platform install complete."
echo "    TAK cluster:  see tak-cluster/DEPLOY-RUNBOOK.md"
echo "    Streams:      managed as CRDs in platform/nats/streams/ (nack controller)"
echo "    tak-agent:    needs cloudtak-agent-token secret (see step 9 above)"
echo "    feeds:        needs omos-feed-secrets (see step 8b above)"
echo "    AISSTREAM:    sign up free at https://aisstream.io to get API key"
