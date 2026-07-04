#!/usr/bin/env bash
# OneMind hub — Step 2 platform layer.
# Reproducible record of the platform components installed on the cluster.
# These are the load-bearing primitives every service depends on.
# (Follow-up: adopt these as Argo CD Applications for full GitOps management.)
set -euo pipefail

# --- Argo CD: GitOps engine. Deploy = commit to Git; cluster self-reconciles.
helm repo add argo https://argoproj.github.io/argo-helm
helm upgrade --install argocd argo/argo-cd -n argocd --create-namespace --wait

# --- cert-manager: automatic TLS certificate issuance + renewal.
helm repo add jetstack https://charts.jetstack.io
helm upgrade --install cert-manager jetstack/cert-manager -n cert-manager \
  --create-namespace --set crds.enabled=true --wait

# --- kube-prometheus-stack: metrics (Prometheus) + dashboards (Grafana)
#     + alerting (Alertmanager) + node/kube-state exporters.
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm upgrade --install monitoring prometheus-community/kube-prometheus-stack \
  -n monitoring --create-namespace --wait

helm repo update
