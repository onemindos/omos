#!/bin/bash
# next-sku.sh — Find the next available book SKU
# Usage: ./next-sku.sh [PHY|DIG|BTH]
# Default format: PHY

set -euo pipefail

VAULT_ROOT="$(cd "$(dirname "$0")/../../../.." && pwd)"
LIB="$VAULT_ROOT/_library"
FORMAT="${1:-PHY}"

# Validate format
if [[ ! "$FORMAT" =~ ^(PHY|DIG|BTH)$ ]]; then
  echo "Error: Invalid format '$FORMAT'. Must be: PHY, DIG, BTH"
  exit 1
fi

# Find highest SKU across all formats (sequence is global)
HIGHEST=$(grep -rh "^sku:" "$LIB/" 2>/dev/null | sed 's/sku: //' | sort -t'-' -k3 -n | tail -1)

if [[ -z "$HIGHEST" ]]; then
  NEXT_NUM="0001"
else
  CURRENT_NUM=$(echo "$HIGHEST" | grep -o '[0-9]*$')
  NEXT_NUM=$(printf "%04d" $((10#$CURRENT_NUM + 1)))
fi

NEXT_SKU="BK-$FORMAT-$NEXT_NUM"

echo "Current highest: ${HIGHEST:-none}"
echo "Next SKU:        $NEXT_SKU"
echo ""
echo "$NEXT_SKU"
