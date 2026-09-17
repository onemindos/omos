#!/bin/bash
# setup-symlinks.sh — Create workspace symlinks for OpenClaw multi-agent
# Run as: sudo -u onemind bash setup-symlinks.sh
#
# OpenClaw expects workspaces at ~/.openclaw/workspace-{agent}/
# This script symlinks them to the monorepo so git sync works automatically.

set -euo pipefail

REPO=/home/onemind/onemind-os
OCDIR=/home/onemind/.openclaw

echo "=== OneMind — Workspace Symlink Setup ==="
echo "Repo:     $REPO"
echo "OpenClaw: $OCDIR"
echo ""

# Verify repo exists
if [ ! -d "$REPO/openclaw/workspaces" ]; then
  echo "ERROR: $REPO/openclaw/workspaces not found. Clone the repo first."
  exit 1
fi

# Verify .openclaw dir exists
if [ ! -d "$OCDIR" ]; then
  echo "ERROR: $OCDIR not found. Run 'openclaw gateway' once first to initialize."
  exit 1
fi

# Create symlinks
declare -A WORKSPACES=(
  ["workspace"]="legacy"
  ["workspace-legacy"]="legacy"
  ["workspace-hp"]="hp"
  ["workspace-le"]="le"
  ["workspace-ge"]="ge"
  ["workspace-community"]="community"
)

for link_name in "${!WORKSPACES[@]}"; do
  agent="${WORKSPACES[$link_name]}"
  target="$REPO/openclaw/workspaces/$agent"
  link="$OCDIR/$link_name"

  if [ ! -d "$target" ]; then
    echo "SKIP: $target does not exist"
    continue
  fi

  if [ -L "$link" ]; then
    current=$(readlink "$link")
    if [ "$current" = "$target" ]; then
      echo "OK:   $link_name → $agent (already correct)"
      continue
    fi
    echo "FIX:  $link_name → $agent (was: $current)"
    rm "$link"
  elif [ -d "$link" ]; then
    echo "WARN: $link is a real directory, backing up to ${link}.bak"
    mv "$link" "${link}.bak"
  fi

  ln -s "$target" "$link"
  echo "LINK: $link_name → $agent"
done

echo ""
echo "=== Symlinks complete ==="
ls -la "$OCDIR"/workspace*
