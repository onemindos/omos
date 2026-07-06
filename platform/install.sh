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
# Requires: tunnel-credentials secret (inject from 1Password)
# -----------------------------------------------------------------------------
kubectl --context "$CTX" create namespace cloudflare --dry-run=client -o yaml | kubectl --context "$CTX" apply -f -
kubectl --context "$CTX" apply -f "$PLATFORM/cloudflare/cloudflared.yaml"

# -----------------------------------------------------------------------------
# 4. NATS — 3-node HA JetStream cluster
# Requires: nats/values.secret.yaml (has ONEMIND/SYS passwords — inject from 1Password)
# -----------------------------------------------------------------------------
helm repo add nats https://nats-io.github.io/k8s/helm/charts --force-update
kubectl --context "$CTX" create namespace nats --dry-run=client -o yaml | kubectl --context "$CTX" apply -f -
kubectl --context "$CTX" apply -f "$PLATFORM/nats/tls.yaml"
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
# 7. LEGACY AGENT (Hermes — mother)
# Requires: hermes-env secret (AWS creds from 1Password)
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/hermes/namespace.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/hermes/pvc.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/hermes/deployment.yaml"

echo ""
echo "==> Platform install complete."
echo "    TAK cluster: see tak-cluster/DEPLOY-RUNBOOK.md"
echo "    GeoServer:   kubectl apply -f geoserver/deployment.yaml"
echo "    Streams:     nats-migration/provision.sh"
