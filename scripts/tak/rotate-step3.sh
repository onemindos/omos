#!/usr/bin/env bash
# TAK CA/keystore rotation, step 3: replace the last stock-password and unencrypted key material.
#
# Run by the owner in the maintenance window (see platform/takserver/DEPLOY-RUNBOOK.md, "STEP 3"). It never
# prints a password or key. Passwords are generated in this process (openssl rand), kept in a private mktemp
# directory (mode 0700, removed on exit), handed to keytool/openssl through environment variables (never argv)
# and end up only inside the SealedSecret files this script writes.
#
# What it does, per consumer (all re-encryptions keep the certificates and private keys; only the password
# changes, and the script checks that certificate fingerprints and public keys are identical before and after):
#
#   leaf     cert-migration ConfigMap (takserver.key admin.key user.key, takserver.p12 admin.p12 user.p12
#            admin.jks user.jks truststore-root.p12)  ->  Secret tak-leaf-identities (takserver ns)
#   fedhub   Secret fedhub-certs (own CA)             ->  fedhub-certs-v2 + fedhub-store-pass + fedhub-ca-keys
#   takml    Secret tak-keystores-copy (takml ns)     ->  takml-tak-client (only the two stores TAK-ML reads)
#   portal   (new)                                    ->  tak-portal-package-pass (onetak ns)
#
# Usage:  OLD_PASS=<stock TAK password> scripts/tak/rotate-step3.sh [options] <command>
#
#   preflight   read-only: tools, cluster reach, controller certificate, the stock password opens every source,
#               target files and Secrets do not exist yet
#   backup      create Secret tak-step3-backup-<date> (takserver ns) holding every source file; refuses to
#               overwrite an existing backup
#   generate    re-encrypt and write the SealedSecret files into the working tree (nothing is applied)
#   stage       kubectl apply the generated SealedSecrets so the Secrets exist before the roll (nothing mounts
#               them yet) and wait until the controller has unsealed them
#   verify      from the cluster: new passwords open the new stores, the stock password does not, fingerprints
#               and public keys equal the originals
#
# Options:
#   --out DIR     write SealedSecret files under DIR (same relative paths) instead of the repository root
#   --only LIST   comma list of leaf,fedhub,takml,portal (default: all)
#   --force       overwrite existing SealedSecret files (this draws NEW passwords: restage afterwards)
#
# Environment: OLD_PASS (required, the stock TAK store password), KUBE_CONTEXT (default do-nyc3-omos-cloud),
#              OPENSSL (default openssl; OpenSSL 3 or newer, not LibreSSL).
# Tools: kubectl, kubeseal, keytool (JDK 11+), openssl, jq.
set -euo pipefail
umask 077

CTX="${KUBE_CONTEXT:-do-nyc3-omos-cloud}"
OPENSSL="${OPENSSL:-openssl}"
STAMP="20261029"
BACKUP_SECRET="tak-step3-backup-${STAMP}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
OUT_ROOT="$REPO_ROOT"
ONLY="leaf,fedhub,takml,portal"
FORCE=0
W=""

# Sources and targets -------------------------------------------------------------------------------------------
LEAF_KEYS="admin.key takserver.key user.key"
LEAF_STORES="admin.jks admin.p12 takserver.p12 truststore-root.p12 user.jks user.p12"
FILE_LEAF="platform/takserver/deploy/tak-leaf-identities.sealedsecret.yaml"
FILE_FED_CERTS="platform/takserver/fedhub/fedhub-certs-v2.sealedsecret.yaml"
FILE_FED_PASS="platform/takserver/fedhub/fedhub-store-pass.sealedsecret.yaml"
FILE_FED_CA="platform/takserver/fedhub/fedhub-ca-keys.sealedsecret.yaml"
FILE_TAKML="platform/takml/takml-tak-client.sealedsecret.yaml"
FILE_PORTAL="platform/onetak/tak-portal-package-pass.sealedsecret.yaml"

log() { printf '%s\n' "$*"; }
die() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }
K() { kubectl --context "$CTX" "$@"; }
kt() { keytool -J-Dkeystore.pkcs12.legacy "$@"; }
want() { case ",$ONLY," in *",$1,"*) return 0 ;; *) return 1 ;; esac; }

cleanup() {
  if [ -n "$W" ] && [ -d "$W" ]; then rm -rf "$W"; fi
  unset OLD_PASS NEW_PASS
}
trap cleanup EXIT
trap 'exit 130' INT TERM

need() { command -v "$1" >/dev/null 2>&1 || die "$1 not found in PATH"; }

init_work() {
  W="$(mktemp -d "${TMPDIR:-/tmp}/tak-step3.XXXXXX")"
  mkdir -p "$W/src" "$W/dst" "$W/pw"
}

require_tools() {
  need kubectl; need kubeseal; need keytool; need jq; need "$OPENSSL"
  if "$OPENSSL" version | grep -qi libressl; then die "OPENSSL points at LibreSSL; set OPENSSL to an OpenSSL 3+ binary"; fi
  [ -n "${OLD_PASS:-}" ] || die "set OLD_PASS to the stock TAK store password (docs/operations/tak-ca-password-rotation-step1-2026-10-29.md)"
  export OLD_PASS
  K get namespace takserver >/dev/null || die "cannot reach context $CTX"
}

outpath() { printf '%s/%s' "$OUT_ROOT" "$1"; }

# Fetching ------------------------------------------------------------------------------------------------------
# secret_file NS SECRET KEY OUT   (decoded value; fails if the key is missing)
secret_file() {
  local ns=$1 name=$2 key=$3 out=$4 j
  j="$(K -n "$ns" get secret "$name" -o json)"
  printf '%s' "$j" | jq -e --arg k "$key" '.data | has($k)' >/dev/null || die "Secret $ns/$name has no key $key"
  printf '%s' "$j" | jq -r --arg k "$key" '.data[$k]' | tr -d '\n' | "$OPENSSL" base64 -d -A > "$out"
}

# cm_file NS CM KEY OUT   (binaryData first, then data)
cm_file() {
  local ns=$1 name=$2 key=$3 out=$4 j
  j="$(K -n "$ns" get configmap "$name" -o json)"
  if printf '%s' "$j" | jq -e --arg k "$key" '(.binaryData // {}) | has($k)' >/dev/null; then
    printf '%s' "$j" | jq -r --arg k "$key" '.binaryData[$k]' | tr -d '\n' | "$OPENSSL" base64 -d -A > "$out"
    return 0
  fi
  printf '%s' "$j" | jq -e --arg k "$key" '(.data // {}) | has($k)' >/dev/null || die "ConfigMap $ns/$name has no key $key"
  printf '%s\n' "$(printf '%s' "$j" | jq -r --arg k "$key" '.data[$k]')" > "$out"
}

secret_keys() { K -n "$1" get secret "$2" -o json | jq -r '.data | keys[]'; }

# Crypto helpers (passwords only through the environment variable NAMES passed in) -------------------------------
store_type() { # file passvar
  kt -list -keystore "$1" -storepass:env "$2" 2>/dev/null | sed -n 's/^Keystore type: *//p' | head -1
}

store_opens() { # file passvar
  kt -list -keystore "$1" -storepass:env "$2" >/dev/null 2>&1
}

# Certificate-only PKCS#12 files (made by openssl, for example truststore-root.p12) show no entries to keytool, so
# they are read and rewritten with openssl instead.
p12_pem() { # file passvar : certificates as PEM on stdout
  "$OPENSSL" pkcs12 -legacy -in "$1" -passin "env:$2" -nokeys 2>/dev/null \
    || "$OPENSSL" pkcs12 -in "$1" -passin "env:$2" -nokeys 2>/dev/null
}

p12_fps() { # file passvar
  local n=0 line
  rm -f "$W"/c.*.pem
  while IFS= read -r line; do
    case "$line" in
      "-----BEGIN CERTIFICATE-----") n=$((n + 1)); printf '%s\n' "$line" > "$W/c.$n.pem" ;;
      "-----END CERTIFICATE-----") printf '%s\n' "$line" >> "$W/c.$n.pem" ;;
      "-----"* | "") : ;;
      *) if [ "$n" -gt 0 ]; then printf '%s\n' "$line" >> "$W/c.$n.pem"; fi ;;
    esac
  done < <(p12_pem "$1" "$2")
  local i
  [ "$n" -gt 0 ] || return 0
  for i in $(seq 1 "$n"); do "$OPENSSL" x509 -in "$W/c.$i.pem" -noout -fingerprint -sha256 | sed 's/^.*=//'; done | sort
  rm -f "$W"/c.*.pem
}

store_entries() { # file passvar
  kt -list -keystore "$1" -storepass:env "$2" 2>/dev/null | sed -n 's/^Your keystore contains \([0-9]*\) entr.*/\1/p' | head -1
}

store_fps() { # file passvar : sorted SHA-256 fingerprints of every certificate in the store
  local out
  out="$(kt -list -v -keystore "$1" -storepass:env "$2" 2>/dev/null | sed -n 's/^.*SHA256: *//p' | sort)"
  if [ -z "$out" ]; then out="$(p12_fps "$1" "$2")"; fi
  printf '%s\n' "$out"
}

key_opens() { # file passvar
  "$OPENSSL" pkey -in "$1" -passin "env:$2" -noout >/dev/null 2>&1
}

key_pub_hash() { # file passvar
  "$OPENSSL" pkey -in "$1" -passin "env:$2" -pubout 2>/dev/null | "$OPENSSL" dgst -sha256 | sed 's/^.* //'
}

rewrap_store() { # src dst  : OLD_PASS -> NEW_PASS, same type and aliases
  local src=$1 dst=$2 type
  type="$(store_type "$src" OLD_PASS)"
  [ -n "$type" ] || die "$(basename "$src"): the stock password does not open it"
  rm -f "$dst"
  if [ "$(store_entries "$src" OLD_PASS)" = "0" ]; then
    p12_pem "$src" OLD_PASS > "$W/certs-only.pem"
    [ -s "$W/certs-only.pem" ] || die "$(basename "$src"): no certificates found"
    "$OPENSSL" pkcs12 -export -legacy -nokeys -in "$W/certs-only.pem" -passout env:NEW_PASS -out "$dst" \
      || die "re-encrypting $(basename "$src") failed"
    rm -f "$W/certs-only.pem"
  else
    # -destkeypass is only valid for JKS key entries: keytool silently skips trusted-certificate entries when it
    # is given, so the entry count is checked below.
    local keypass=()
    if [ "$type" = "JKS" ]; then keypass=(-destkeypass:env NEW_PASS); fi
    kt -importkeystore -noprompt \
      -srckeystore "$src" -srcstoretype "$type" -srcstorepass:env OLD_PASS \
      -destkeystore "$dst" -deststoretype "$type" -deststorepass:env NEW_PASS "${keypass[@]+"${keypass[@]}"}" \
      >"$W/keytool.err" 2>&1 || { cat "$W/keytool.err" >&2; die "re-encrypting $(basename "$src") failed"; }
    if grep -qi "problem importing" "$W/keytool.err"; then cat "$W/keytool.err" >&2; die "re-encrypting $(basename "$src") skipped entries"; fi
    [ "$(store_entries "$src" OLD_PASS)" = "$(store_entries "$dst" NEW_PASS)" ] \
      || die "$(basename "$src"): entry count changed"
  fi
  [ -s "$dst" ] || die "re-encrypting $(basename "$src") produced no file"
  store_fps "$src" OLD_PASS > "$W/fp.old"
  store_fps "$dst" NEW_PASS > "$W/fp.new"
  [ -s "$W/fp.old" ] || die "$(basename "$src"): no certificate fingerprints read"
  cmp -s "$W/fp.old" "$W/fp.new" || die "$(basename "$src"): certificate fingerprints changed"
  if store_opens "$dst" OLD_PASS; then die "$(basename "$src"): stock password still opens the new store"; fi
}

rewrap_key() { # src dst : PKCS#8, AES-256-CBC, PBKDF2 100000 iterations, same key
  local src=$1 dst=$2
  key_opens "$src" OLD_PASS || die "$(basename "$src"): cannot read with the stock password (or unencrypted)"
  "$OPENSSL" pkcs8 -topk8 -v2 aes-256-cbc -iter 100000 -in "$src" -passin env:OLD_PASS \
    -passout env:NEW_PASS -out "$dst" || die "re-encrypting $(basename "$src") failed"
  [ "$(key_pub_hash "$src" OLD_PASS)" = "$(key_pub_hash "$dst" NEW_PASS)" ] \
    || die "$(basename "$src"): public key changed"
  if key_opens "$dst" OLD_PASS && head -1 "$dst" | grep -vq ENCRYPTED; then die "$(basename "$src"): output is not encrypted"; fi
}

new_password() { # name : generates, exports NEW_PASS, stores a copy for sealing
  NEW_PASS="$("$OPENSSL" rand -hex 24)"
  export NEW_PASS
  printf '%s' "$NEW_PASS" > "$W/pw/$1"
}

# Sealing -------------------------------------------------------------------------------------------------------
fetch_cert() {
  [ -s "$W/pub-cert.pem" ] || kubeseal --context "$CTX" --fetch-cert --controller-name=sealed-secrets-controller \
    --controller-namespace=kube-system > "$W/pub-cert.pem"
  [ -s "$W/pub-cert.pem" ] || die "could not fetch the sealed-secrets certificate"
}

# seal NS NAME REL_OUT_PATH --from-file=... (files only: no secret on any command line)
seal() {
  local ns=$1 name=$2 rel=$3 out
  shift 3
  out="$(outpath "$rel")"
  if [ -e "$out" ] && [ "$FORCE" -ne 1 ]; then die "$rel exists (use --force to draw new passwords and overwrite)"; fi
  mkdir -p "$(dirname "$out")"
  fetch_cert
  K create secret generic "$name" -n "$ns" "$@" --dry-run=client -o json \
    | kubeseal --cert "$W/pub-cert.pem" --format yaml --scope namespace-wide \
    | awk '
        /^[[:space:]]*creationTimestamp: null[[:space:]]*$/ { next }
        { print }
        !done && /sealedsecrets.bitnami.com\/namespace-wide:/ { print "    argocd.argoproj.io/sync-wave: \"-1\""; done = 1 }
      ' > "$out.new"
  grep -q '^kind: SealedSecret' "$out.new" || die "sealing $name produced no SealedSecret"
  mv "$out.new" "$out"
  log "  wrote $rel"
}

# Components ----------------------------------------------------------------------------------------------------
gen_leaf() {
  log "leaf: cert-migration -> tak-leaf-identities"
  local d="$W/src/leaf" o="$W/dst/leaf" f args=()
  mkdir -p "$d" "$o"
  new_password leaf
  for f in $LEAF_KEYS; do cm_file takserver cert-migration "$f" "$d/$f"; rewrap_key "$d/$f" "$o/$f"; args+=("--from-file=$f=$o/$f"); done
  for f in $LEAF_STORES; do cm_file takserver cert-migration "$f" "$d/$f"; rewrap_store "$d/$f" "$o/$f"; args+=("--from-file=$f=$o/$f"); done
  args+=("--from-file=leaf-password=$W/pw/leaf")
  seal takserver tak-leaf-identities "$FILE_LEAF" "${args[@]}"
}

gen_fedhub() {
  log "fedhub: fedhub-certs -> fedhub-certs-v2, fedhub-store-pass, fedhub-ca-keys"
  local d="$W/src/fed" o="$W/dst/fed" f certs=() cas=()
  mkdir -p "$d" "$o"
  new_password fedhub
  while IFS= read -r f; do
    secret_file takserver fedhub-certs "$f" "$d/$f"
    case "$f" in
      ca-do-not-share.key | root-ca-do-not-share.key)
        rewrap_key "$d/$f" "$o/$f"; cas+=("--from-file=$f=$o/$f") ;;
      *.key)
        rewrap_key "$d/$f" "$o/$f"; certs+=("--from-file=$f=$o/$f") ;;
      *.jks | *.p12)
        rewrap_store "$d/$f" "$o/$f"; certs+=("--from-file=$f=$o/$f") ;;
      *)
        certs+=("--from-file=$f=$d/$f") ;;
    esac
  done < <(secret_keys takserver fedhub-certs)
  [ "${#cas[@]}" -eq 2 ] || die "expected both FedHub CA keys in fedhub-certs"
  seal takserver fedhub-certs-v2 "$FILE_FED_CERTS" "${certs[@]}"
  seal takserver fedhub-store-pass "$FILE_FED_PASS" "--from-file=keystore-password=$W/pw/fedhub"
  seal takserver fedhub-ca-keys "$FILE_FED_CA" "${cas[@]}"
}

gen_takml() {
  log "takml: tak-keystores-copy -> takml-tak-client"
  local d="$W/src/takml" o="$W/dst/takml" f args=()
  mkdir -p "$d" "$o"
  new_password takml
  for f in takserver.jks truststore-root.jks; do
    secret_file takml tak-keystores-copy "$f" "$d/$f"; rewrap_store "$d/$f" "$o/$f"; args+=("--from-file=$f=$o/$f")
  done
  args+=("--from-file=store-password=$W/pw/takml")
  seal takml takml-tak-client "$FILE_TAKML" "${args[@]}"
}

gen_portal() {
  log "portal: tak-portal-package-pass"
  new_password portal
  seal onetak tak-portal-package-pass "$FILE_PORTAL" "--from-file=package-password=$W/pw/portal"
}

# Commands ------------------------------------------------------------------------------------------------------
cmd_preflight() {
  log "context: $CTX    OpenSSL: $("$OPENSSL" version | cut -d' ' -f1-2)    keytool: $(keytool -J-version 2>&1 | head -1)"
  fetch_cert
  log "sealed-secrets certificate: ok"
  local f s
  if want leaf; then
    mkdir -p "$W/src/leaf"
    for f in $LEAF_KEYS; do cm_file takserver cert-migration "$f" "$W/src/leaf/$f"; key_opens "$W/src/leaf/$f" OLD_PASS || die "$f unreadable"; done
    for f in $LEAF_STORES; do cm_file takserver cert-migration "$f" "$W/src/leaf/$f"; store_opens "$W/src/leaf/$f" OLD_PASS || die "$f: stock password does not open it"; done
    log "leaf: 3 keys and 6 stores readable with the stock password"
  fi
  if want fedhub; then
    mkdir -p "$W/src/fed"
    while IFS= read -r f; do
      secret_file takserver fedhub-certs "$f" "$W/src/fed/$f"
      case "$f" in
        *.key) key_opens "$W/src/fed/$f" OLD_PASS || die "fedhub $f unreadable" ;;
        *.jks | *.p12) store_opens "$W/src/fed/$f" OLD_PASS || die "fedhub $f: stock password does not open it" ;;
        *) : ;;
      esac
    done < <(secret_keys takserver fedhub-certs)
    log "fedhub: all stores and keys readable with the stock password"
  fi
  if want takml; then
    mkdir -p "$W/src/takml"
    for f in takserver.jks truststore-root.jks; do secret_file takml tak-keystores-copy "$f" "$W/src/takml/$f"; store_opens "$W/src/takml/$f" OLD_PASS || die "takml $f: stock password does not open it"; done
    log "takml: both stores readable with the stock password"
  fi
  for f in "$FILE_LEAF:leaf" "$FILE_FED_CERTS:fedhub" "$FILE_FED_PASS:fedhub" "$FILE_FED_CA:fedhub" "$FILE_TAKML:takml" "$FILE_PORTAL:portal"; do
    want "${f##*:}" || continue
    if [ -e "$(outpath "${f%%:*}")" ]; then log "note: ${f%%:*} already exists (generate needs --force)"; fi
  done
  for s in "takserver tak-leaf-identities" "takserver fedhub-certs-v2" "takserver fedhub-store-pass" "takserver fedhub-ca-keys" "takml takml-tak-client" "onetak tak-portal-package-pass"; do
    # shellcheck disable=SC2086
    set -- $s
    if K -n "$1" get secret "$2" >/dev/null 2>&1; then log "note: Secret $1/$2 already exists in the cluster"; fi
  done
  grep -q tak-leaf-identities "$REPO_ROOT/platform/takserver/deploy/takserver-config.yaml" \
    || log "warning: this working tree does not contain the step-3 manifests"
  log "preflight ok"
}

cmd_backup() {
  if K -n takserver get secret "$BACKUP_SECRET" >/dev/null 2>&1; then die "$BACKUP_SECRET already exists; not overwriting"; fi
  local b="$W/backup" f args=()
  mkdir -p "$b"
  for f in $LEAF_KEYS $LEAF_STORES; do cm_file takserver cert-migration "$f" "$b/cert-migration.$f"; args+=("--from-file=cert-migration.$f=$b/cert-migration.$f"); done
  while IFS= read -r f; do secret_file takserver fedhub-certs "$f" "$b/fedhub-certs.$f"; args+=("--from-file=fedhub-certs.$f=$b/fedhub-certs.$f"); done < <(secret_keys takserver fedhub-certs)
  while IFS= read -r f; do secret_file takml tak-keystores-copy "$f" "$b/tak-keystores-copy.$f"; args+=("--from-file=tak-keystores-copy.$f=$b/tak-keystores-copy.$f"); done < <(secret_keys takml tak-keystores-copy)
  K -n takserver create secret generic "$BACKUP_SECRET" "${args[@]}" >/dev/null
  log "created Secret takserver/$BACKUP_SECRET (${#args[@]} files). Delete it after the soak period."
}

cmd_generate() {
  fetch_cert
  want leaf && gen_leaf
  want fedhub && gen_fedhub
  want takml && gen_takml
  want portal && gen_portal
  log "done. Review with: git status; do not commit anything else. Nothing was applied to the cluster."
}

apply_file() { # rel ns name
  local f
  f="$(outpath "$1")"
  [ -s "$f" ] || die "$1 has not been generated"
  K apply -f "$f" >/dev/null
  log "  applied $1"
}

wait_secret() { # ns name
  for _ in $(seq 1 30); do
    if K -n "$1" get secret "$2" >/dev/null 2>&1; then log "  Secret $1/$2 unsealed"; return 0; fi
    sleep 2
  done
  die "Secret $1/$2 did not appear (sealed-secrets controller error?)"
}

cmd_stage() {
  if want leaf; then apply_file "$FILE_LEAF"; wait_secret takserver tak-leaf-identities; fi
  if want fedhub; then
    apply_file "$FILE_FED_CERTS"; apply_file "$FILE_FED_PASS"; apply_file "$FILE_FED_CA"
    wait_secret takserver fedhub-certs-v2; wait_secret takserver fedhub-store-pass; wait_secret takserver fedhub-ca-keys
  fi
  if want takml; then apply_file "$FILE_TAKML"; wait_secret takml takml-tak-client; fi
  if want portal; then apply_file "$FILE_PORTAL"; wait_secret onetak tak-portal-package-pass; fi
  log "staged. Nothing mounts these Secrets until the manifests are merged."
}

pass_of() { # ns secret key -> password in NEW_PASS (exported)
  NEW_PASS="$(K -n "$1" get secret "$2" -o json | jq -r --arg k "$3" '.data[$k] // empty' | tr -d '\n' | "$OPENSSL" base64 -d -A)"
  [ -n "$NEW_PASS" ] || die "Secret $1/$2 has no $3"
  export NEW_PASS
}

verify_store() { # label oldfile newfile
  store_opens "$3" NEW_PASS || die "$1: the new password does not open the new store"
  if store_opens "$3" OLD_PASS; then die "$1: the stock password still opens the new store"; fi
  store_fps "$2" OLD_PASS > "$W/fp.old"
  store_fps "$3" NEW_PASS > "$W/fp.new"
  if [ ! -s "$W/fp.old" ] || ! cmp -s "$W/fp.old" "$W/fp.new"; then die "$1: certificate fingerprints differ from the original"; fi
  log "  ok $1 (fingerprints unchanged)"
}

verify_key() { # label oldfile newfile
  key_opens "$3" NEW_PASS || die "$1: the new password does not open the new key"
  if key_opens "$3" OLD_PASS; then die "$1: the stock password still opens the new key"; fi
  [ "$(key_pub_hash "$2" OLD_PASS)" = "$(key_pub_hash "$3" NEW_PASS)" ] || die "$1: public key differs from the original"
  log "  ok $1 (public key unchanged)"
}

cmd_verify() {
  local d="$W/src" o="$W/dst" f
  if want leaf; then
    log "leaf"; mkdir -p "$d/leaf" "$o/leaf"; pass_of takserver tak-leaf-identities leaf-password
    for f in $LEAF_KEYS; do cm_file takserver cert-migration "$f" "$d/leaf/$f"; secret_file takserver tak-leaf-identities "$f" "$o/leaf/$f"; verify_key "$f" "$d/leaf/$f" "$o/leaf/$f"; done
    for f in $LEAF_STORES; do cm_file takserver cert-migration "$f" "$d/leaf/$f"; secret_file takserver tak-leaf-identities "$f" "$o/leaf/$f"; verify_store "$f" "$d/leaf/$f" "$o/leaf/$f"; done
  fi
  if want fedhub; then
    log "fedhub"; mkdir -p "$d/fed" "$o/fed"; pass_of takserver fedhub-store-pass keystore-password
    while IFS= read -r f; do
      secret_file takserver fedhub-certs "$f" "$d/fed/$f"
      case "$f" in
        ca-do-not-share.key | root-ca-do-not-share.key)
          secret_file takserver fedhub-ca-keys "$f" "$o/fed/$f"; verify_key "$f" "$d/fed/$f" "$o/fed/$f" ;;
        *.key) secret_file takserver fedhub-certs-v2 "$f" "$o/fed/$f"; verify_key "$f" "$d/fed/$f" "$o/fed/$f" ;;
        *.jks | *.p12) secret_file takserver fedhub-certs-v2 "$f" "$o/fed/$f"; verify_store "$f" "$d/fed/$f" "$o/fed/$f" ;;
        *) secret_file takserver fedhub-certs-v2 "$f" "$o/fed/$f"; cmp -s "$d/fed/$f" "$o/fed/$f" || die "$f differs from the original" ;;
      esac
    done < <(secret_keys takserver fedhub-certs)
    if secret_keys takserver fedhub-certs-v2 | grep -q 'do-not-share'; then die "fedhub-certs-v2 must not hold CA private keys"; fi
  fi
  if want takml; then
    log "takml"; mkdir -p "$d/takml" "$o/takml"; pass_of takml takml-tak-client store-password
    for f in takserver.jks truststore-root.jks; do secret_file takml tak-keystores-copy "$f" "$d/takml/$f"; secret_file takml takml-tak-client "$f" "$o/takml/$f"; verify_store "$f" "$d/takml/$f" "$o/takml/$f"; done
  fi
  if want portal; then
    log "portal"; pass_of onetak tak-portal-package-pass package-password
    [ "${#NEW_PASS}" -ge 32 ] || die "portal package password is too short"
    log "  ok package-password present (${#NEW_PASS} characters)"
  fi
  log "verify ok"
}

# Main ----------------------------------------------------------------------------------------------------------
CMD=""
while [ $# -gt 0 ]; do
  case "$1" in
    --out) [ $# -ge 2 ] || die "--out needs a directory"; OUT_ROOT="$2"; shift 2 ;;
    --only) [ $# -ge 2 ] || die "--only needs a list"; ONLY="$2"; shift 2 ;;
    --force) FORCE=1; shift ;;
    -h | --help) sed -n '2,45p' "${BASH_SOURCE[0]}"; exit 0 ;;
    preflight | backup | generate | stage | verify) CMD="$1"; shift ;;
    *) die "unknown argument: $1" ;;
  esac
done
[ -n "$CMD" ] || die "no command (preflight, backup, generate, stage, verify); see --help"

require_tools
init_work
"cmd_$CMD"
