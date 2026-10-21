#!/bin/sh
# Checks seed.sh against temporary directories. Run: sh platform/hermes/seed/test-seed.sh
set -eu
HERE="$(cd "$(dirname "$0")" && pwd)"
T="$(mktemp -d)"
trap 'rm -rf "$T"' EXIT
fail() { echo "FAIL: $1"; exit 1; }

mkdir -p "$T/seed/agents/legacy" "$T/seed/cron" "$T/seed/skills/devops/one" "$T/seed/config"
echo "soul from git" > "$T/seed/agents/legacy/SOUL.md"
echo "[]" > "$T/seed/cron/jobs.json"
echo "skill from git" > "$T/seed/skills/devops/one/SKILL.md"
# shellcheck disable=SC2016  # the ${...} is literal template text
printf 'token: ${MATRIX_TOKEN}\n' > "$T/seed/config/config.template.yaml"

run() { SEED_DIR="$T/seed" DATA_DIR="$T/data" sh "$HERE/seed.sh" >"$T/out" 2>&1; }

# 1. Empty volume: identity is seeded; a template with placeholders is NOT installed as config.yaml.
mkdir -p "$T/data"
run
[ "$(cat "$T/data/SOUL.md")" = "soul from git" ] || fail "soul not seeded"
[ -f "$T/data/cron/jobs.json" ] || fail "cron not seeded"
[ -f "$T/data/skills/devops/one/SKILL.md" ] || fail "skill not seeded"
[ ! -e "$T/data/config.yaml" ] || fail "config with placeholders was installed"
grep -q "not installed" "$T/out" || fail "no warning about placeholders"

# 2. Existing files are never overwritten on a normal start.
echo "live soul" > "$T/data/SOUL.md"
echo "live skill" > "$T/data/skills/devops/one/SKILL.md"
run
[ "$(cat "$T/data/SOUL.md")" = "live soul" ] || fail "soul was overwritten"
[ "$(cat "$T/data/skills/devops/one/SKILL.md")" = "live skill" ] || fail "skill was overwritten"

# 3. A template without placeholders is installed, once.
echo "provider: digitalocean" > "$T/seed/config/config.template.yaml"
run
[ "$(cat "$T/data/config.yaml")" = "provider: digitalocean" ] || fail "clean config not installed"
echo "edited" > "$T/data/config.yaml"
run
[ "$(cat "$T/data/config.yaml")" = "edited" ] || fail "config was overwritten"

# 4. HERMES_RESEED=1 is the deliberate reset.
HERMES_RESEED=1 run
[ "$(cat "$T/data/SOUL.md")" = "soul from git" ] || fail "reseed did not reset the soul"

# 5. Nothing to seed from: harmless.
rm -rf "$T/seed"; mkdir -p "$T/seed"
run
grep -q "0 written" "$T/out" || fail "empty seed should write nothing"

echo "PASS: seed.sh behaves as designed"
