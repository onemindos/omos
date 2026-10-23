#!/usr/bin/env bash
# Give AI operators (Claude Code, Codex) read-only access to the cluster and the portal.
# Runs on the admin Mac. See docs/operations/ai-operator-access.md.
#
# 1. Builds ~/.kube/ai-operator.kubeconfig from the `kube/ai-operator` ServiceAccount
#    (platform/kube/ai-operator.yaml) and checks it cannot read Secrets or write.
# 2. Writes ~/.config/onemind/openchoreo-mcp.cred (mode 600) from the Secret
#    auth/ai-operator-openchoreo and checks the read-only OpenChoreo account works.
# 3. Registers three MCP servers for Claude Code (user scope) and Codex (~/.codex/config.toml):
#    onemind-kubernetes (kubernetes-mcp-server --read-only), onemind-backstage and
#    onemind-openchoreo. Claude Code mints an OpenChoreo token per session through
#    scripts/ai-operator/openchoreo-token.sh; Codex reads OPENCHOREO_MCP_TOKEN.
#
# Re-run it to refresh tokens. Credentials stay on this machine; nothing is printed.
set -euo pipefail

ADMIN_CONTEXT="${ADMIN_CONTEXT:-do-nyc3-omos-cloud}"
KUBECONFIG_OUT="${KUBECONFIG_OUT:-$HOME/.kube/ai-operator.kubeconfig}"
BACKSTAGE_TOKEN_FILE="${BACKSTAGE_TOKEN_FILE:-$HOME/.config/onemind/backstage-mcp.token}"
BACKSTAGE_MCP_URL="https://onemind-preview.tail717ba0.ts.net/api/mcp-actions/v1"
K8S_MCP_VERSION="0.0.67"
CODEX_CONFIG="${CODEX_CONFIG:-$HOME/.codex/config.toml}"
OPENCHOREO_CRED_FILE="${OPENCHOREO_CRED_FILE:-$HOME/.config/onemind/openchoreo-mcp.cred}"
OPENCHOREO_MCP_URL="http://onemind-openchoreo-api.tail717ba0.ts.net/mcp"
TOKEN_HELPER="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/openchoreo-token.sh"

log() { printf '==> %s\n' "$*"; }
fail() { printf 'error: %s\n' "$*" >&2; exit 1; }

build_kubeconfig() {
  local cluster server ca token
  cluster=$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$ADMIN_CONTEXT\")].context.cluster}")
  [[ -n "$cluster" ]] || fail "kube context $ADMIN_CONTEXT not found"
  server=$(kubectl config view --raw -o jsonpath="{.clusters[?(@.name==\"$cluster\")].cluster.server}")
  ca=$(kubectl config view --raw -o jsonpath="{.clusters[?(@.name==\"$cluster\")].cluster.certificate-authority-data}")
  token=$(kubectl --context "$ADMIN_CONTEXT" -n kube get secret ai-operator-token -o jsonpath='{.data.token}' | base64 -d)
  [[ -n "$token" ]] || fail "secret kube/ai-operator-token has no token yet (is platform-kube synced?)"

  mkdir -p "$(dirname "$KUBECONFIG_OUT")"
  (
    umask 077
    cat >"$KUBECONFIG_OUT" <<EOF
apiVersion: v1
kind: Config
clusters:
  - name: omos
    cluster:
      server: $server
      certificate-authority-data: $ca
users:
  - name: ai-operator
    user:
      token: $token
contexts:
  - name: omos-ai-operator
    context:
      cluster: omos
      user: ai-operator
current-context: omos-ai-operator
EOF
  )
  log "wrote $KUBECONFIG_OUT"
}

verify_kubeconfig() {
  local k=(kubectl --kubeconfig "$KUBECONFIG_OUT")
  [[ $("${k[@]}" auth can-i list pods --all-namespaces) == yes ]] || fail "cannot list pods"
  [[ $("${k[@]}" auth can-i list applications.argoproj.io -n argocd) == yes ]] || fail "cannot list Argo applications"
  for denied in "get secrets --all-namespaces" "delete pods --all-namespaces" "create pods/exec --all-namespaces" "patch deployments --all-namespaces"; do
    # shellcheck disable=SC2086
    if [[ $("${k[@]}" auth can-i $denied 2>/dev/null || true) == yes ]]; then
      fail "ai-operator is allowed to: $denied"
    fi
  done
  log "kubeconfig verified: reads allowed; secrets, delete, exec and patch denied"
}

write_openchoreo_credential() {
  local user pass
  user=$(kubectl --context "$ADMIN_CONTEXT" -n auth get secret ai-operator-openchoreo -o jsonpath='{.data.username}' | base64 -d)
  pass=$(kubectl --context "$ADMIN_CONTEXT" -n auth get secret ai-operator-openchoreo -o jsonpath='{.data.password}' | base64 -d)
  [[ -n "$user" && -n "$pass" ]] || fail "secret auth/ai-operator-openchoreo is empty (did the identity Job run?)"
  mkdir -p "$(dirname "$OPENCHOREO_CRED_FILE")"
  ( umask 077; printf '%s\n%s\n' "$user" "$pass" >"$OPENCHOREO_CRED_FILE" )
  log "wrote $OPENCHOREO_CRED_FILE"
}

verify_openchoreo() {
  local token reply
  token=$(OPENCHOREO_MCP_CRED="$OPENCHOREO_CRED_FILE" "$TOKEN_HELPER" --raw) || fail "could not mint an OpenChoreo token"
  call() {
    curl -fsS --max-time 20 "$OPENCHOREO_MCP_URL" -H "Authorization: Bearer $token" \
      -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' -d "$1"
  }
  reply=$(call '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"list_namespaces","arguments":{}}}')
  grep -q '"namespaces"' <<<"$reply" || fail "OpenChoreo list_namespaces did not return namespaces"
  reply=$(call '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"create_namespace","arguments":{"name":"ai-operator-write-check"}}}' || true)
  grep -q 'not authorized' <<<"$reply" || fail "OpenChoreo did not refuse create_namespace: the account is not read-only"
  log "OpenChoreo verified: reads allowed, create_namespace refused"
}

register_claude() {
  command -v claude >/dev/null || { log "claude CLI not found, skipping Claude Code"; return; }
  local token
  token=$(<"$BACKSTAGE_TOKEN_FILE")
  claude mcp remove --scope user onemind-kubernetes >/dev/null 2>&1 || true
  claude mcp remove --scope user onemind-backstage >/dev/null 2>&1 || true
  claude mcp add --scope user onemind-kubernetes -- \
    npx -y "kubernetes-mcp-server@$K8S_MCP_VERSION" --read-only --disable-multi-cluster \
    --kubeconfig "$KUBECONFIG_OUT" >/dev/null
  claude mcp add --scope user --transport http onemind-backstage "$BACKSTAGE_MCP_URL" \
    --header "Authorization: Bearer $token" >/dev/null
  claude mcp remove --scope user onemind-openchoreo >/dev/null 2>&1 || true
  claude mcp add-json --scope user onemind-openchoreo \
    "$(python3 -c 'import json,sys; print(json.dumps({"type":"http","url":sys.argv[1],"headersHelper":sys.argv[2]}))' "$OPENCHOREO_MCP_URL" "$TOKEN_HELPER")" >/dev/null
  log "registered Kubernetes, portal and OpenChoreo MCP servers in Claude Code (user scope)"
}

register_codex() {
  [[ -f "$CODEX_CONFIG" ]] || { log "no $CODEX_CONFIG, skipping Codex"; return; }
  python3 - "$CODEX_CONFIG" "$KUBECONFIG_OUT" "$BACKSTAGE_TOKEN_FILE" "$BACKSTAGE_MCP_URL" "$K8S_MCP_VERSION" "$OPENCHOREO_MCP_URL" <<'PY'
import re, sys, pathlib
path, kubeconfig, token_file, url, version, oc_url = sys.argv[1:]
token = pathlib.Path(token_file).read_text().strip()
cfg = pathlib.Path(path)
text = cfg.read_text()
blocks = {
    "onemind-kubernetes": (
        '[mcp_servers.onemind-kubernetes]\n'
        'command = "npx"\n'
        f'args = ["-y", "kubernetes-mcp-server@{version}", "--read-only", '
        f'"--disable-multi-cluster", "--kubeconfig", "{kubeconfig}"]\n'
        'startup_timeout_sec = 60\n'
    ),
    "onemind-backstage": (
        '[mcp_servers.onemind-backstage]\n'
        f'url = "{url}"\n'
        f'http_headers = {{ Authorization = "Bearer {token}" }}\n'
    ),
    "onemind-openchoreo": (
        '[mcp_servers.onemind-openchoreo]\n'
        f'url = "{oc_url}"\n'
        'bearer_token_env_var = "OPENCHOREO_MCP_TOKEN"\n'
    ),
}
for name, block in blocks.items():
    pattern = re.compile(r'^\[mcp_servers\.' + re.escape(name) + r'\]\n(?:(?!^\[).*\n?)*', re.M)
    text = pattern.sub('', text).rstrip('\n') + '\n\n' + block
cfg.write_text(text)
PY
  chmod 600 "$CODEX_CONFIG"
  log "registered Kubernetes, portal and OpenChoreo MCP servers in $CODEX_CONFIG"
  log "Codex: export OPENCHOREO_MCP_TOKEN=\"\$($TOKEN_HELPER --raw)\" before starting it (tokens last one hour)"
}

[[ -r "$BACKSTAGE_TOKEN_FILE" ]] || fail "missing $BACKSTAGE_TOKEN_FILE"
build_kubeconfig
verify_kubeconfig
write_openchoreo_credential
verify_openchoreo
register_claude
register_codex
log "done. Start a new Claude Code or Codex session to load the servers."
