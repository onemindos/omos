#!/bin/bash
# audit-library.sh — Audit library for common issues
# Usage: ./audit-library.sh

set -euo pipefail

VAULT_ROOT="$(cd "$(dirname "$0")/../../../.." && pwd)"
LIB="$VAULT_ROOT/_library"

echo "=== OneMind Library Audit ==="
echo ""

# Count total books
TOTAL=$(grep -rl "^type: book" "$LIB/" 2>/dev/null | wc -l | tr -d ' ')
echo "Total books: $TOTAL"
echo ""

# Books missing ai_summary
echo "--- Missing ai_summary ---"
MISSING_SUMMARY=0
while IFS= read -r file; do
  if ! grep -q "^ai_summary:" "$file" 2>/dev/null || grep -q '^ai_summary: ""' "$file" 2>/dev/null; then
    echo "  $(basename "$(dirname "$file")")/$(basename "$file")"
    ((MISSING_SUMMARY++))
  fi
done < <(grep -rl "^type: book" "$LIB/" 2>/dev/null)
echo "  Count: $MISSING_SUMMARY"
echo ""

# Books missing cover-front.jpg
echo "--- Missing cover-front.jpg ---"
MISSING_COVER=0
while IFS= read -r file; do
  DIR=$(dirname "$file")
  if [[ ! -f "$DIR/cover-front.jpg" ]]; then
    echo "  $(basename "$DIR")"
    ((MISSING_COVER++))
  fi
done < <(grep -rl "^type: book" "$LIB/" 2>/dev/null)
echo "  Count: $MISSING_COVER"
echo ""

# Duplicate SKUs
echo "--- Duplicate SKUs ---"
DUPES=$(grep -rh "^sku:" "$LIB/" 2>/dev/null | sort | uniq -d)
if [[ -z "$DUPES" ]]; then
  echo "  None found"
else
  echo "$DUPES" | while read -r line; do
    echo "  DUPLICATE: $line"
    grep -rl "$line" "$LIB/" 2>/dev/null | while read -r f; do
      echo "    → $f"
    done
  done
fi
echo ""

# Genre distribution
echo "--- Genre Distribution ---"
grep -rh "^  - " "$LIB/" 2>/dev/null | grep -v "type/" | grep -v "status/" | grep -v "format/" | grep -v "category/" | grep -v "location/" | grep -v "read/" | head -30
echo ""

# SKU range
echo "--- SKU Range ---"
grep -rh "^sku:" "$LIB/" 2>/dev/null | sort -t'-' -k3 -n | head -1
grep -rh "^sku:" "$LIB/" 2>/dev/null | sort -t'-' -k3 -n | tail -1
echo ""

# Orphan images in staging
STAGING_COUNT=$(find "$LIB/_staging" -name "*.jpg" -o -name "*.JPG" -o -name "*.HEIC" -o -name "*.heic" 2>/dev/null | wc -l | tr -d ' ')
echo "Staging images remaining: $STAGING_COUNT"
echo ""

echo "=== Audit Complete ==="
