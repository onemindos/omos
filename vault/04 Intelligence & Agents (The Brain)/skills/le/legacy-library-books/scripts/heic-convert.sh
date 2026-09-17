#!/bin/bash
# heic-convert.sh — Convert HEIC images to JPG using macOS sips
# Usage: ./heic-convert.sh [file.HEIC | directory]
# If no argument, converts all HEIC in _library/_staging/

set -euo pipefail

VAULT_ROOT="$(cd "$(dirname "$0")/../../../.." && pwd)"
STAGING="$VAULT_ROOT/_library/_staging"

convert_file() {
  local input="$1"
  local output="${input%.HEIC}.jpg"
  output="${output%.heic}.jpg"

  if [[ -f "$output" ]]; then
    echo "SKIP (exists): $(basename "$output")"
    return 0
  fi

  sips -s format jpeg "$input" --out "$output" >/dev/null 2>&1
  echo "CONVERTED: $(basename "$input") → $(basename "$output")"
}

if [[ $# -eq 0 ]]; then
  # Default: convert all HEIC in staging
  count=0
  for f in "$STAGING"/*.HEIC "$STAGING"/*.heic; do
    [[ -f "$f" ]] || continue
    convert_file "$f"
    ((count++))
  done
  echo "--- Done: $count files converted in _staging/"
elif [[ -d "$1" ]]; then
  # Convert all HEIC in specified directory
  count=0
  for f in "$1"/*.HEIC "$1"/*.heic; do
    [[ -f "$f" ]] || continue
    convert_file "$f"
    ((count++))
  done
  echo "--- Done: $count files converted in $1"
elif [[ -f "$1" ]]; then
  # Convert single file
  convert_file "$1"
else
  echo "Error: '$1' is not a file or directory"
  exit 1
fi
