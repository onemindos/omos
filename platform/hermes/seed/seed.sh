#!/bin/sh
# Seed Hermes' identity into an EMPTY data volume from the identity files kept in Git.
#
# Meant to run as an init container before Hermes starts (see README.md in this directory; it is not wired into the
# Deployment yet). It never overwrites anything that already exists, so a normal restart does nothing.
#
#   SEED_DIR   where the identity files are mounted, laid out like platform/hermes/agent (default /seed)
#   DATA_DIR   the Hermes data volume (default /opt/data)
#   HERMES_RESEED=1  overwrite the seeded files that already exist (a deliberate reset)
set -eu

SEED_DIR="${SEED_DIR:-/seed}"
DATA_DIR="${DATA_DIR:-/opt/data}"
RESEED="${HERMES_RESEED:-0}"

copied=0
skipped=0

# copy_file <source relative to SEED_DIR> <destination relative to DATA_DIR>
copy_file() {
  src="$SEED_DIR/$1"
  dst="$DATA_DIR/$2"
  [ -f "$src" ] || return 0
  if [ -e "$dst" ] && [ "$RESEED" != "1" ]; then
    echo "[seed] keep existing $2"
    skipped=$((skipped + 1))
    return 0
  fi
  mkdir -p "$(dirname "$dst")"
  cp "$src" "$dst"
  echo "[seed] wrote $2"
  copied=$((copied + 1))
}

# Persona: the edited soul lives at agents/legacy/SOUL.md in the repo.
copy_file agents/legacy/SOUL.md SOUL.md
copy_file cron/jobs.json cron/jobs.json

# Skills: copy each file, never overwriting an existing one.
if [ -d "$SEED_DIR/skills" ]; then
  (cd "$SEED_DIR/skills" && find . -type f) | while IFS= read -r f; do
    f="${f#./}"
    src="$SEED_DIR/skills/$f"
    dst="$DATA_DIR/skills/$f"
    if [ -e "$dst" ] && [ "$RESEED" != "1" ]; then
      continue
    fi
    mkdir -p "$(dirname "$dst")"
    cp "$src" "$dst"
  done
  echo "[seed] skills synced (existing files kept)"
fi

# Config: the exported template has ${PLACEHOLDERS} where credentials were. Installing it unresolved would break
# Hermes, so it is only used when it has none left; otherwise it is left for a human to fill in.
tmpl="$SEED_DIR/config/config.template.yaml"
if [ -f "$tmpl" ]; then
  if grep -q '[$]{[A-Z0-9_]*}' "$tmpl"; then
    echo "[seed] config.template.yaml still has \${PLACEHOLDERS}; not installed. Fill them in and copy it to config.yaml."
  else
    copy_file config/config.template.yaml config.yaml
  fi
fi

echo "[seed] done: $copied written, $skipped kept"
