#!/usr/bin/env bash
# Owner key session for the tool gateway: create each credential in its own UI, paste it here, and this script seals it.
#
#   scripts/gateway/keys-session.sh [--dry-run [--out DIR]] list
#   scripts/gateway/keys-session.sh [--dry-run] all                 every step in order, each can be skipped
#   scripts/gateway/keys-session.sh [--dry-run] grafana argocd ...  only these steps (re-run one to replace its key)
#
# Steps: ghcr nats gate grafana argocd github digitalocean cloudflare n8n homeassistant cloudtak karakeep takportal takserver
#        legacy-ssh legacy-export
#
# Rules this script keeps:
#   * A value is read with `read -s` (no echo, never in the shell history, never in argv, never exported). It is passed to
#     `kubectl create secret ... --dry-run=client -o json | kubeseal` through a pipe/descriptor. Only the file name and the
#     length of each value are printed.
#   * It only writes sealed files (and, for a target, moves its pending/ manifests into manifests/, re-renders the gateway
#     config, and sets the execute-gate replicas to 1 once its three sealed files exist). It commits nothing: review `git diff`.
#   * --dry-run needs no cluster: it works on a temporary copy of the repo files (--out DIR, default a mktemp dir), writes
#     placeholder "sealed" files that list key names and lengths only, and never calls kubeseal or the cluster.
#
# Environment: KUBE_CONTEXT (default do-nyc3-omos-cloud), GATE_USER (default gate; use another name to retry after a failed run).
set -euo pipefail
umask 077

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_ROOT="$(cd "$HERE/../.." && pwd)"
CTX="${KUBE_CONTEXT:-do-nyc3-omos-cloud}"
GATE_USER="${GATE_USER:-gate}"
OWNER_MXID="${OWNER_MXID:-@zeus:onemindos.dev}"
DRY=0
OUT=""
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
CERT=""
TOUCHED=()

die() { printf 'error: %s\n' "$*" >&2; exit 1; }
say() { printf '%s\n' "$*" >&2; }

# ---------------------------------------------------------------------------------------------------------------- setup

setup_repo() {
  if [ "$DRY" = 1 ]; then
    OUT="${OUT:-$(mktemp -d)}"
    REPO="$OUT/repo"
    mkdir -p "$REPO/platform/hermes/onemind-hermes" "$REPO/scripts"
    cp -R "$SRC_ROOT/platform/tool-gateway" "$REPO/platform/"
    cp "$SRC_ROOT/platform/hermes/onemind-hermes/kustomization.yaml" "$REPO/platform/hermes/onemind-hermes/"
    cp -R "$SRC_ROOT/scripts/gateway" "$REPO/scripts/"
    say "DRY RUN: working on a copy in $REPO; no cluster, no kubeseal, repo untouched."
  else
    REPO="$SRC_ROOT"
    command -v kubectl >/dev/null || die "kubectl not found"
    command -v kubeseal >/dev/null || die "kubeseal not found"
    kubectl --context "$CTX" get ns tool-gateway >/dev/null 2>&1 || die "cannot reach context $CTX (or namespace tool-gateway)"
  fi
  TG="$REPO/platform/tool-gateway"
  HM="$REPO/platform/hermes/onemind-hermes"
}

fetch_cert() {
  [ "$DRY" = 1 ] && return 0
  if [ -z "$CERT" ]; then
    CERT="$TMP/pub-cert.pem"
    kubeseal --context "$CTX" --fetch-cert --controller-name=sealed-secrets-controller --controller-namespace=kube-system >"$CERT" \
      || die "could not fetch the sealed-secrets certificate"
  fi
}

# ------------------------------------------------------------------------------------------------------------- helpers

# Reads a secret value without echo into the named variable. Empty input is refused; whitespace and control characters too.
ask_secret() { # ask_secret VAR "label" MINLEN
  local __v=""
  printf '  %s (input hidden): ' "$2" >&2
  IFS= read -r -s __v || true
  printf '\n' >&2
  [ -n "$__v" ] || { say "  nothing entered"; return 1; }
  if [[ "$__v" =~ [[:space:][:cntrl:]] ]]; then say "  the value contains whitespace or control characters; paste only the token"; return 1; fi
  if [ "${#__v}" -lt "$3" ]; then say "  too short (${#__v} characters, expected at least $3)"; return 1; fi
  printf -v "$1" '%s' "$__v"
}

ask_line() { # ask_line VAR "label" [default]
  local __v=""
  printf '  %s%s: ' "$2" "${3:+ [$3]}" >&2
  IFS= read -r __v || true
  printf -v "$1" '%s' "${__v:-${3:-}}"
}

# Shows where to create the credential and offers a skip. Returns 1 when skipped.
begin_step() { # begin_step name "title" "instructions"
  printf '\n== %s: %s ==\n%s\n' "$1" "$2" "$3" >&2
  local a=""
  printf '  Press Enter when you have it, or s to skip: ' >&2
  IFS= read -r a || true
  case "$a" in s|S) say "  skipped $1"; return 1;; esac
  return 0
}

# Echoes key lengths to stderr (names and counts only), passes the Secret JSON on unchanged.
length_tap() {
  python3 -c '
import base64, json, sys
raw = sys.stdin.read()
d = json.loads(raw)
print("  " + ", ".join("%s=%d bytes" % (k, len(base64.b64decode(v))) for k, v in sorted(d.get("data", {}).items())), file=sys.stderr)
sys.stdout.write(raw)'
}

dry_seal() { # reads Secret JSON on stdin, writes a placeholder that has key names and lengths only
  python3 -c '
import base64, json, sys
d = json.load(sys.stdin)
print("# DRY RUN: not a real SealedSecret. Key names and value lengths only; no value is stored.")
print("apiVersion: bitnami.com/v1alpha1\nkind: SealedSecret\nmetadata:\n  name: %s\n  namespace: %s\nspec:\n  encryptedData:" % (d["metadata"]["name"], d["metadata"]["namespace"]))
for k, v in sorted(d.get("data", {}).items()):
    print("    %s: DRYRUN-%d-bytes" % (k, len(base64.b64decode(v))))'
}

# seal NAME NAMESPACE OUTFILE [--type T] (-v KEY VALUE | -f KEY PATH)...
# -v: a value from a shell variable (token, password), passed to kubectl through a process substitution (printf is a builtin,
#     so the value is never in an argument list); -f: a file (keys, certificates, binary), only its path is an argument.
seal() {
  local name=$1 ns=$2 out=$3
  shift 3
  local i=0 k
  local cmd
  cmd="kubectl create secret generic $(printf '%q' "$name") -n $(printf '%q' "$ns") --dry-run=client -o json"
  while [ $# -gt 0 ]; do
    case "$1" in
      --type) cmd+=" --type $(printf '%q' "$2")"; shift 2;;
      -v|-f)
        k=$2
        [[ "$k" =~ ^[A-Za-z0-9._-]+$ ]] || die "seal: bad key name $k"
        if [ "$1" = -v ]; then
          printf -v "__sv$i" '%s' "$3"
          # The eval string names the variable, not its value; the key name was validated above.
          cmd+=" --from-file=$k=<(printf '%s' \"\$__sv$i\")"
          i=$((i + 1))
        else
          cmd+=" --from-file=$k=$(printf '%q' "$3")"
        fi
        shift 3;;
      *) die "seal: bad argument $1";;
    esac
  done
  fetch_cert
  mkdir -p "$(dirname "$out")"
  local tmp="$out.tmp"
  if [ "$DRY" = 1 ]; then
    eval "$cmd" | length_tap | dry_seal >"$tmp"
  else
    eval "$cmd" | length_tap | kubeseal --cert "$CERT" --scope namespace-wide --format yaml >"$tmp"
    grep -q '^  encryptedData:' "$tmp" || { rm -f "$tmp"; die "sealing $name failed (no encryptedData in the output)"; }
  fi
  while [ "$i" -gt 0 ]; do i=$((i - 1)); unset "__sv$i"; done
  mv "$tmp" "$out"
  TOUCHED+=("${out#"$REPO"/}")
  say "  wrote ${out#"$REPO"/}"
}

# Copy a live Secret key to a variable without printing it. Dry run: a placeholder of the right shape.
live_value() { # live_value VAR namespace secret key
  if [ "$DRY" = 1 ]; then printf -v "$1" '%s' "dryrun-placeholder-value-0123456789"; return 0; fi
  local __v
  __v="$(kubectl --context "$CTX" -n "$2" get secret "$3" -o "jsonpath={.data.$4}" | base64 -d)" || die "cannot read $2/$3 key $4"
  [ -n "$__v" ] || die "$2/$3 key $4 is empty"
  printf -v "$1" '%s' "$__v"
}

# Move the pending/<name>/manifests tree (backend, egress, namespace rules) into manifests/, once.
activate() {
  local name=$1 pend="$TG/pending/$1"
  if [ -d "$pend/manifests" ]; then
    cp -R "$pend/manifests/." "$TG/manifests/"
    rm -rf "$pend"
    say "  activated $name: moved pending/$name/manifests into manifests/"
  else
    say "  $name: nothing pending (already active, or no manifests needed)"
  fi
  render
}

render() {
  (cd "$REPO" && python3 scripts/gateway/render_gateway_config.py >/dev/null) || die "render_gateway_config.py failed"
  say "  re-rendered platform/tool-gateway/values.generated.yaml"
}

kustomize_add() { # kustomize_add FILE  (adds "  - FILE" to the Hermes kustomization after gateway-client)
  python3 - "$HM/kustomization.yaml" "$1" <<'PY'
import sys
path, name = sys.argv[1], sys.argv[2]
lines = open(path).read().split("\n")
entry = "  - " + name
if entry not in lines:
    i = next(i for i, l in enumerate(lines) if l.strip() == "- gateway-client.sealedsecret.yaml")
    lines.insert(i + 1, entry)
    open(path, "w").write("\n".join(lines))
PY
}

maybe_enable_gate() {
  local f
  for f in gate-matrix gate-nats ghcr; do
    [ -f "$TG/manifests/$f.sealedsecret.yaml" ] || { say "  execute-gate stays off: $f.sealedsecret.yaml is not there yet"; return 0; }
  done
  python3 - "$TG/manifests/execute-gate/deployment.yaml" <<'PY'
import re, sys
p = sys.argv[1]
s = open(p).read()
s2 = re.sub(r"(?m)^  replicas: 0$", "  replicas: 1", s, count=1)
open(p, "w").write(s2)
PY
  say "  gate-matrix, gate-nats and ghcr are all sealed: execute-gate replicas set to 1"
}

# --------------------------------------------------------------------------------------------------------------- steps

step_ghcr() {
  begin_step ghcr "image pull secret for ghcr.io/onemindos (execute-gate image is private)" \
"  Default: reuse the pull secret that already works in the cluster (matrix/ghcr). Nothing to create.
  Alternative: a classic PAT with scope read:packages (GitHub > Settings > Developer settings > Personal access tokens >
  Tokens (classic) > Generate new token (classic) > tick read:packages > Configure SSO > Authorize onemindos)." || return 0
  local mode=""
  ask_line mode "Reuse matrix/ghcr (r) or paste a PAT (p)" r
  if [ "$mode" = p ]; then
    local pat=""
    ask_secret pat "classic PAT with read:packages" 20 || return 0
    local json
    json="$(printf '%s' "$pat" | python3 -c '
import base64, json, sys
t = sys.stdin.read()
print(json.dumps({"auths": {"ghcr.io": {"username": "onemindos", "password": t, "auth": base64.b64encode(("onemindos:" + t).encode()).decode()}}}))')"
    seal ghcr tool-gateway "$TG/manifests/ghcr.sealedsecret.yaml" --type kubernetes.io/dockerconfigjson -v .dockerconfigjson "$json"
  else
    local cfg
    if [ "$DRY" = 1 ]; then cfg='{"auths":{"ghcr.io":{"auth":"ZHJ5OnJ1bg=="}}}'
    else cfg="$(kubectl --context "$CTX" -n matrix get secret ghcr -o 'jsonpath={.data.\.dockerconfigjson}' | base64 -d)" || die "cannot read matrix/ghcr"; fi
    seal ghcr tool-gateway "$TG/manifests/ghcr.sealedsecret.yaml" --type kubernetes.io/dockerconfigjson -v .dockerconfigjson "$cfg"
  fi
  maybe_enable_gate
}

step_nats() {
  begin_step nats "NATS credentials and CA for the execute-gate" \
"  Nothing to create. No dedicated gate user exists (the nats release is Helm-managed), so this reseals the bus user that
  onemind-core already uses (onemind-core/core-runtime NATS_USER, NATS_PASSWORD) and the bus CA (nats/nats-ca ca.crt).
  The gate only writes the KV bucket exec-approvals and evt.toolgw.approval.*. A narrower user needs a change to the Helm values." || return 0
  local user pass ca
  live_value user onemind-core core-runtime NATS_USER
  live_value pass onemind-core core-runtime NATS_PASSWORD
  if [ "$DRY" = 1 ]; then printf '%s' "-----BEGIN CERTIFICATE-----" >"$TMP/ca.crt"
  else kubectl --context "$CTX" -n nats get secret nats-ca -o 'jsonpath={.data.ca\.crt}' | base64 -d >"$TMP/ca.crt" || die "cannot read nats/nats-ca"; fi
  ca="$TMP/ca.crt"
  seal gate-nats tool-gateway "$TG/manifests/gate-nats.sealedsecret.yaml" -v user "$user" -v password "$pass" -f ca.crt "$ca"
  rm -f "$ca"
  maybe_enable_gate
}

step_gate() {
  begin_step gate "execute-gate Matrix bot and the private approvals room" \
"  Nothing to create by hand. This registers @$GATE_USER:onemindos.dev (not an admin, random password that is never shown),
  logs in, creates the private unencrypted room 'approvals', invites $OWNER_MXID, and seals the bot token and room id.
  Afterwards: open Element (https://element.tail717ba0.ts.net) and ACCEPT the invitation to 'approvals'.
  If a previous run registered the bot but failed later, retry with GATE_USER=gate2 (the password of the old bot is gone)." || return 0
  local pw tok room out
  pw="$(openssl rand -base64 36 | tr -d '/+=\n' | cut -c1-32)"
  if [ "$DRY" = 1 ]; then
    tok="dryrun-matrix-token-0123456789abcdef"; room='!dryrun:onemindos.dev'
  else
    out="$(printf '{"user":"%s","password":"%s","owner":"%s","room_name":"approvals"}' "$GATE_USER" "$pw" "$OWNER_MXID" \
      | kubectl --context "$CTX" -n matrix exec -i deploy/synapse -c synapse -- python3 -c "$(cat "$HERE/gate_bootstrap.py")")" \
      || die "gate bootstrap failed (see the messages above; no secret was printed)"
    tok="$(printf '%s' "$out" | python3 -c 'import json,sys; print(json.load(sys.stdin)["token"])')"
    room="$(printf '%s' "$out" | python3 -c 'import json,sys; print(json.load(sys.stdin)["room_id"])')"
  fi
  pw=""
  seal gate-matrix tool-gateway "$TG/manifests/gate-matrix.sealedsecret.yaml" -v token "$tok" -v room_id "$room"
  say "  room id sealed (not printed). Accept the invite in Element."
  maybe_enable_gate
}

step_grafana() {
  begin_step grafana "Grafana service-account token (target grafana, group ops)" \
"  Grafana (https://grafana.tail717ba0.ts.net) > Administration > Users and access > Service accounts > Add service account
    name legacy-gateway, role Viewer (enough while the gateway backend runs read-only; use Editor only when writes are enabled).
  Then 'Add service account token': display name gateway, expiration 90 days > Generate token > copy it." || return 0
  local t=""
  ask_secret t "service account token (glsa_...)" 20 || return 0
  seal gateway-target-grafana tool-gateway "$TG/manifests/gateway-target-grafana.sealedsecret.yaml" -v token "$t"
  activate grafana
}

step_argocd() {
  begin_step argocd "Argo CD account 'legacy' API token (target argocd, group ops)" \
"  Prerequisite: the PR that adds account 'legacy' to platform/values/argocd.yaml is merged and platform-argocd is Synced.
  Argo CD (https://argocd.tail717ba0.ts.net) > Settings > Accounts > legacy > Generate New (token) > copy it. (Needs the admin login.)
  CLI alternative: argocd account generate-token --account legacy
  The account can only read applications, logs, projects and clusters (no sync, no delete). The script also generates the inbound
  bearer that the Argo CD MCP server insists on (key mcp_auth) and never shows it." || return 0
  local t="" m
  ask_secret t "Argo CD token for account legacy" 20 || return 0
  m="$(openssl rand -hex 24)"
  seal gateway-target-argocd tool-gateway "$TG/manifests/gateway-target-argocd.sealedsecret.yaml" -v token "$t" -v mcp_auth "$m"
  activate argocd
}

step_github() {
  begin_step github "GitHub fine-grained token, read only (target github, group dev)" \
"  GitHub > profile picture > Settings > Developer settings > Personal access tokens > Fine-grained tokens > Generate new token
    Token name legacy-gateway, Expiration 90 days, Resource owner: onemindos
    Repository access: Only select repositories (omos and the repos Legacy should read)
    Repository permissions, all 'Read-only': Contents, Issues, Pull requests, Actions, Discussions, Commit statuses, Dependabot alerts,
    Code scanning alerts, Projects (Metadata is added automatically). Nothing with Write or Administration.
  If the organisation requires approval: Organization > Settings > Personal access tokens > Pending requests > Approve." || return 0
  local t=""
  ask_secret t "fine-grained token (github_pat_...)" 40 || return 0
  seal gateway-target-github tool-gateway "$TG/manifests/gateway-target-github.sealedsecret.yaml" -v token "$t"
  activate github
}

step_digitalocean() {
  begin_step digitalocean "DigitalOcean token, no delete scopes (targets doapps, dodoks, donetworking, group cloud)" \
"  cloud.digitalocean.com > API (left menu) > Tokens > Generate New Token
    Name legacy-gateway, Expiration 90 days, Scopes: Custom Scopes. Tick Read for: app, kubernetes, domain, firewall, load_balancer,
    vpc, droplet, database, account (Create/Update only for what Legacy may change). Do NOT tick any Delete scope.
  Copy the token (shown once)." || return 0
  local t=""
  ask_secret t "DigitalOcean token (dop_v1_...)" 30 || return 0
  seal gateway-target-digitalocean tool-gateway "$TG/manifests/gateway-target-digitalocean.sealedsecret.yaml" -v token "$t"
  activate digitalocean
}

step_cloudflare() {
  begin_step cloudflare "Cloudflare API token, read only (target cloudflare, group cloud)" \
"  dash.cloudflare.com > profile icon > My Profile > API Tokens > Create Token > Create Custom Token (Get started)
    Token name legacy-gateway
    Permissions (all Read): Account > Account Settings, Account > Workers Scripts, Zone > Zone, Zone > DNS
    Account Resources: Include > OneMind account. Zone Resources: Include > All zones. TTL: end date in 90 days. Client IP filtering: leave off.
  Continue to summary > Create Token > copy it. The 'execute' tool runs any API call this token allows, so keep every permission Read." || return 0
  local t=""
  ask_secret t "Cloudflare API token" 30 || return 0
  seal gateway-target-cloudflare tool-gateway "$TG/manifests/gateway-target-cloudflare.sealedsecret.yaml" -v token "$t"
  activate cloudflare
}

step_n8n() {
  begin_step n8n "n8n MCP access token and REST API key (target n8n, group work)" \
"  n8n (your n8n URL) > Settings > Instance-level MCP > enable it > copy the Access Token. Then open each workflow Legacy may use >
    workflow menu > Settings > turn on 'Available in MCP' (only those workflows are visible).
  n8n > Settings > n8n API > Create an API key (label legacy-gateway, expiration 90 days) > copy it. (Stored for a later REST target;
  nothing reads it yet. Press Enter at that prompt to skip it.)" || return 0
  local mcp="" api=""
  ask_secret mcp "MCP access token" 20 || return 0
  if ask_secret api "n8n REST API key (optional)" 10; then
    seal gateway-target-n8n tool-gateway "$TG/manifests/gateway-target-n8n.sealedsecret.yaml" -v mcp_token "$mcp" -v api_key "$api"
  else
    seal gateway-target-n8n tool-gateway "$TG/manifests/gateway-target-n8n.sealedsecret.yaml" -v mcp_token "$mcp"
  fi
  activate n8n
}

step_homeassistant() {
  begin_step homeassistant "Home Assistant long-lived token (target homeassistant, group home)" \
"  1. Settings > People > Users tab (turn on Advanced mode in your profile first) > Add user: name legacy, NOT administrator, local access only.
  2. Settings > Devices & services > Add integration > 'Model Context Protocol Server'.
  3. Settings > Voice assistants > Expose: expose only non-physical entities. Expose NO lock, cover, alarm panel or siren yet.
  4. Sign in as user 'legacy' > profile (bottom left) > Security tab > Long-lived access tokens > Create token, name legacy-gateway > copy it.
  Spike H1 is open: if /api/mcp rejects a long-lived token (OAuth only), the fragment must point at the community ha-mcp instead." || return 0
  local t=""
  ask_secret t "long-lived access token" 100 || return 0
  seal gateway-target-homeassistant tool-gateway "$TG/manifests/gateway-target-homeassistant.sealedsecret.yaml" -v token "$t"
  activate homeassistant
}

step_cloudtak() {
  begin_step cloudtak "CloudTAK API token (target cloudtak, group tak)" \
"  Default: reuse the agent token Legacy already has (hermes/cloudtak-agent-token, key token). Nothing to create.
  New token instead: CloudTAK > profile menu (your name) > Tokens > New Token (name legacy-gateway) > copy it." || return 0
  local mode="" t=""
  ask_line mode "Reuse hermes/cloudtak-agent-token (r) or paste a new token (p)" r
  if [ "$mode" = p ]; then ask_secret t "CloudTAK API token" 20 || return 0
  else live_value t hermes cloudtak-agent-token token; fi
  seal gateway-target-cloudtak tool-gateway "$TG/manifests/gateway-target-cloudtak.sealedsecret.yaml" -v token "$t"
  activate cloudtak
}

step_karakeep() {
  begin_step karakeep "Karakeep API key (stored only: no gateway target yet)" \
"  Karakeep (your Karakeep URL) > profile menu > Settings > API Keys > New API Key, name legacy-gateway > copy it.
  The MCP server is stdio-only and no image turns it into HTTP yet (platform/tool-gateway/pending/karakeep/README.md), so the key is
  sealed and kept, but nothing uses it and no target is activated." || return 0
  local t=""
  ask_secret t "Karakeep API key" 20 || return 0
  seal gateway-target-karakeep tool-gateway "$TG/manifests/gateway-target-karakeep.sealedsecret.yaml" -v api_key "$t"
}

step_takportal() {
  begin_step takportal "TAK Portal credential (stored only: no gateway target yet)" \
"  The portal's API is enroll, download and health only; nothing for a read tool. This step only stores a credential for a future adapter:
  TAK Portal (https://enroll.tail717ba0.ts.net) > your admin account password or API credential. Skip unless you have a use for it." || return 0
  local t=""
  ask_secret t "TAK Portal credential" 8 || return 0
  seal gateway-target-takportal tool-gateway "$TG/manifests/gateway-target-takportal.sealedsecret.yaml" -v credential "$t"
}

step_takserver() {
  begin_step takserver "TAK Server client certificate legacy-gateway (stored only: no adapter yet)" \
"  Generate a client certificate for a dedicated TAK user 'legacy-gateway' with the TAK Server cert tooling or the Certificate
  Authority page of the admin UI (https://tak.tail717ba0.ts.net:8446 > Configuration > Manage Certificates), exported as a .p12 file.
  Give the file path and the p12 password. Nothing consumes it until a TAK Server adapter exists." || return 0
  local path="" pw=""
  ask_line path "path of the .p12 file"
  [ -f "$path" ] || { say "  no such file: $path"; return 0; }
  ask_secret pw "p12 password" 4 || return 0
  seal gateway-target-takserver tool-gateway "$TG/manifests/gateway-target-takserver.sealedsecret.yaml" -f client.p12 "$path" -v password "$pw"
}

step_legacy-ssh() {
  begin_step legacy-ssh "SSH private key Legacy uses to reach the Mac (optional)" \
"  Generate on the Mac: ssh-keygen -t ed25519 -N '' -C legacy@onemind -f ~/.ssh/legacy_ed25519
  Add legacy_ed25519.pub to ~/.ssh/authorized_keys on the Mac (restrict it with a command= or from= option if you can).
  Give the PRIVATE key file path (and optionally a known_hosts file). Sealed as Secret hermes/legacy-ssh keys id_ed25519 and
  known_hosts. The pod mount is not part of this change: mount the Secret read-only (defaultMode 0400) where the Hermes terminal tool
  expects ~/.ssh. Delete the local private key copy afterwards if you do not need it." || return 0
  local key="" kh=""
  ask_line key "path of the private key file"
  [ -f "$key" ] || { say "  no such file: $key"; return 0; }
  ask_line kh "path of a known_hosts file (Enter for none)" ""
  if [ -n "$kh" ] && [ -f "$kh" ]; then
    seal legacy-ssh hermes "$HM/legacy-ssh.sealedsecret.yaml" -f id_ed25519 "$key" -f known_hosts "$kh"
  else
    seal legacy-ssh hermes "$HM/legacy-ssh.sealedsecret.yaml" -f id_ed25519 "$key"
  fi
  kustomize_add legacy-ssh.sealedsecret.yaml
}

step_legacy-export() {
  begin_step legacy-export "credential for the nightly config export (optional)" \
"  The job pushes snapshots to a private repo. Choose one:
  (s) SSH deploy key (works with the job as written). GitHub > repo onemindos/onemind-hermes-config > Settings > Deploy keys > Add deploy key >
      paste the .pub of: ssh-keygen -t ed25519 -N '' -C legacy-config-export -f ./legacy-export-key > tick 'Allow write access'.
      BLOCKED while the organisation has deploy keys disabled (Organization > Settings > Repository > Deploy keys policy).
  (t) Fine-grained token. GitHub > Settings > Developer settings > Fine-grained tokens > Generate new token > Resource owner onemindos >
      Only select repositories: onemind-hermes-config > Repository permissions: Contents = Read and write > 90 days.
      The export script reads only ssh-privatekey today, so a token also needs a small change to legacy_config_export.py (HTTPS push)." || return 0
  local mode="" repo="" key="" t=""
  ask_line mode "SSH key (s) or token (t)" s
  if [ "$mode" = t ]; then
    ask_line repo "HTTPS repo URL" "https://github.com/onemindos/onemind-hermes-config.git"
    ask_secret t "fine-grained token" 40 || return 0
    seal onemind-hermes-export hermes "$HM/export.sealedsecret.yaml" -v token "$t" -v repo "$repo"
    say "  NOTE: legacy_config_export.py must learn to push over HTTPS with this token before the job is unpaused."
  else
    ask_line repo "SSH repo URL" "ssh://git@ssh.github.com:443/onemindos/onemind-hermes-config.git"
    ask_line key "path of the private deploy key file"
    [ -f "$key" ] || { say "  no such file: $key"; return 0; }
    seal onemind-hermes-export hermes "$HM/export.sealedsecret.yaml" -f ssh-privatekey "$key" -v repo "$repo"
  fi
  kustomize_add export.sealedsecret.yaml
}

# ------------------------------------------------------------------------------------------------------------- driver

STEPS=(ghcr nats gate grafana argocd github digitalocean cloudflare n8n homeassistant cloudtak karakeep takportal takserver legacy-ssh legacy-export)

usage() {
  sed -n '2,22p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'
}

list_steps() {
  cat <<'EOF'
step           output file                                                       needs
ghcr           platform/tool-gateway/manifests/ghcr.sealedsecret.yaml            existing matrix/ghcr (or a read:packages PAT)
nats           platform/tool-gateway/manifests/gate-nats.sealedsecret.yaml       existing onemind-core/core-runtime and nats/nats-ca
gate           platform/tool-gateway/manifests/gate-matrix.sealedsecret.yaml     nothing (script creates bot and room)
grafana        .../gateway-target-grafana.sealedsecret.yaml                      Grafana service-account token
argocd         .../gateway-target-argocd.sealedsecret.yaml                       Argo CD account legacy token
github         .../gateway-target-github.sealedsecret.yaml                       fine-grained PAT, read only
digitalocean   .../gateway-target-digitalocean.sealedsecret.yaml                 PAT without delete scopes
cloudflare     .../gateway-target-cloudflare.sealedsecret.yaml                   read-only API token
n8n            .../gateway-target-n8n.sealedsecret.yaml                          MCP access token (+ optional API key)
homeassistant  .../gateway-target-homeassistant.sealedsecret.yaml                long-lived token of a non-admin user
cloudtak       .../gateway-target-cloudtak.sealedsecret.yaml                     existing agent token, or a new one
karakeep       .../gateway-target-karakeep.sealedsecret.yaml                     API key (stored only)
takportal      .../gateway-target-takportal.sealedsecret.yaml                    credential (stored only)
takserver      .../gateway-target-takserver.sealedsecret.yaml                    .p12 file path and password (stored only)
legacy-ssh     platform/hermes/onemind-hermes/legacy-ssh.sealedsecret.yaml       private key file path (optional)
legacy-export  platform/hermes/onemind-hermes/export.sealedsecret.yaml           deploy key file or token (optional)
EOF
}

main() {
  local -a want=()
  while [ $# -gt 0 ]; do
    case "$1" in
      --dry-run) DRY=1;;
      --out) OUT="$2"; shift;;
      -h|--help|help) usage; exit 0;;
      list) list_steps; exit 0;;
      all) want=("${STEPS[@]}");;
      *) want+=("$1");;
    esac
    shift
  done
  [ "${#want[@]}" -gt 0 ] || { usage; exit 1; }
  local s ok
  for s in "${want[@]}"; do
    ok=0
    for k in "${STEPS[@]}"; do [ "$k" = "$s" ] && ok=1; done
    [ "$ok" = 1 ] || die "unknown step '$s' (try: list)"
  done
  setup_repo
  for s in "${want[@]}"; do "step_$s"; done
  printf '\n== done ==\n' >&2
  if [ "${#TOUCHED[@]}" -gt 0 ]; then
    printf 'files written or changed:\n' >&2
    printf '  %s\n' "${TOUCHED[@]}" >&2
  fi
  if [ "$DRY" = 1 ]; then say "dry run: nothing in the repository changed (copy in $REPO)."
  else say "review with git status / git diff, run the gateway tests, then commit the sealed files and open the PR (docs/operations/tool-gateway.md, owner session)."; fi
}

main "$@"
