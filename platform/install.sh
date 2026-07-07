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
# Requires: cloudflare-api-token secret (inject from 1Password)
# op read "op://00-24 SO/cloudflare-api-token/credential" | \
#   kubectl create secret generic cloudflare-api-token -n cert-manager --from-literal=api-token=-
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/ingress/issuer-dns.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/ingress/issuers.yaml"

# -----------------------------------------------------------------------------
# 3. CLOUDFLARE TUNNEL (zero public inbound)
# Routes are declared in platform/cloudflare/config.yaml — edit + apply to add new routes.
# Requires two secrets (inject from 1Password — both derived from the tunnel token):
#   tunnel-credentials   (original TUNNEL_TOKEN — kept for reference)
#   tunnel-creds-json    (credentials.json format — what cloudflared reads from local config)
#     TOKEN=$(op read "op://00-24 SO/cloudflare-tunnel-token/credential")
#     python3 -c "
#       import sys,json,base64; t='$TOKEN'; p=t+'='*(4-len(t)%4)
#       d=json.loads(base64.urlsafe_b64decode(p))
#       print(json.dumps({'AccountTag':d['a'],'TunnelID':d['t'],'TunnelSecret':d['s']}))
#     " | kubectl create secret generic tunnel-creds-json -n cloudflare --from-file=credentials.json=/dev/stdin
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
  --kube-context "$CTX" -n monitoring --create-namespace
# kube-prometheus-stack uses all defaults

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
# Requires: orp-nats-creds secret (out-of-band via 1Password)
#   kubectl create secret generic orp-nats-creds -n orp \
#     --from-literal=user=onemind \
#     --from-literal=password='<ONEMIND_PW from 1Password>'
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/orp/namespace.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/orp/pvc.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/orp/configmap.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/orp/bridge-script.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/orp/deployment.yaml"

# -----------------------------------------------------------------------------
# 8. OMOS SERVICES — shared omos namespace (tak-projector, feed-world, etc.)
# Requires: omos-nats-creds secret (out-of-band)
#   kubectl create secret generic omos-nats-creds -n omos \
#     --from-literal=user=onemind \
#     --from-literal=password='<ONEMIND_PW from 1Password>'
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/omos/namespace.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos/tak-projector.yaml"

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

# -----------------------------------------------------------------------------
# 9. LEGACY AGENT (Hermes — mother) + TAK-AGENT + OneMind AI Bridge
# Requires:
#   hermes-env secret (AWS creds + NATS from 1Password)
#   hermes-leaf-conf secret (NATS leaf.conf)
#   hermes-nats-ca secret (NATS CA cert)
#   omos-nats-creds secret in hermes namespace:
#     kubectl create secret generic omos-nats-creds -n hermes \
#       --from-literal=user=onemind \
#       --from-literal=password='<ONEMIND_PW from 1Password>'
#   cloudtak-agent-token secret (CloudTAK profile API token):
#     kubectl create secret generic cloudtak-agent-token -n hermes \
#       --from-literal=token='<token from CloudTAK UI: Settings → Tokens → New>'
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/legacy/namespace.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/legacy/pvc.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/legacy/deployment.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/legacy/tak-agent.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/hermes/onemind-ai-bridge.yaml"

# -----------------------------------------------------------------------------
# 10. CLOUDTAK — API + Events + Media + Tiles + MinIO + PostGIS
# Requires: cloudtak-secrets (signing secret, MinIO keys, Postgres DSN)
#   kubectl create secret generic cloudtak-secrets -n cloudtak \
#     --from-literal=SIGNING_SECRET='...' \
#     --from-literal=AWS_S3_AccessKeyId='...' \
#     --from-literal=AWS_S3_SecretAccessKey='...' \
#     --from-literal=POSTGRES='postgresql://...'
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
echo "    Streams:      nats-migration/provision.sh"
echo "    tak-agent:    needs cloudtak-agent-token secret (see step 9 above)"
echo "    feeds:        needs omos-feed-secrets (see step 8b above)"
echo "    AISSTREAM:    sign up free at https://aisstream.io to get API key"
