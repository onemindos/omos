#!/bin/bash
# process-book.sh — Create book folder structure and place cover images
# Usage: ./process-book.sh <category> <genre> <book-slug> <front-img> [back-img]
# Example: ./process-book.sh non-fiction business-finance the-art-of-war IMG_3401.jpg IMG_3402.jpg

set -euo pipefail

VAULT_ROOT="$(cd "$(dirname "$0")/../../../.." && pwd)"
LIB="$VAULT_ROOT/_library"
STG="$LIB/_staging"

if [[ $# -lt 4 ]]; then
  echo "Usage: $0 <category> <genre> <book-slug> <front-img> [back-img]"
  echo "  category:  non-fiction | fiction | reference | children"
  echo "  genre:     genre folder name (lowercase-hyphens)"
  echo "  book-slug: book folder name (lowercase-hyphens)"
  echo "  front-img: filename in _staging/ (e.g., IMG_3401.jpg)"
  echo "  back-img:  (optional) filename in _staging/"
  exit 1
fi

CATEGORY="$1"
GENRE="$2"
SLUG="$3"
FRONT_IMG="$4"
BACK_IMG="${5:-}"

BOOK_DIR="$LIB/$CATEGORY/$GENRE/$SLUG"

# Validate category
if [[ ! "$CATEGORY" =~ ^(non-fiction|fiction|reference|children)$ ]]; then
  echo "Error: Invalid category '$CATEGORY'. Must be: non-fiction, fiction, reference, children"
  exit 1
fi

# Create book directory
mkdir -p "$BOOK_DIR"
echo "CREATED: $BOOK_DIR"

# Copy front cover
if [[ -f "$STG/$FRONT_IMG" ]]; then
  cp "$STG/$FRONT_IMG" "$BOOK_DIR/cover-front.jpg"
  echo "PLACED: cover-front.jpg"
elif [[ -f "$FRONT_IMG" ]]; then
  cp "$FRONT_IMG" "$BOOK_DIR/cover-front.jpg"
  echo "PLACED: cover-front.jpg (absolute path)"
else
  echo "WARNING: Front image not found: $FRONT_IMG"
fi

# Copy back cover (if provided)
if [[ -n "$BACK_IMG" ]]; then
  if [[ -f "$STG/$BACK_IMG" ]]; then
    cp "$STG/$BACK_IMG" "$BOOK_DIR/cover-back.jpg"
    echo "PLACED: cover-back.jpg"
  elif [[ -f "$BACK_IMG" ]]; then
    cp "$BACK_IMG" "$BOOK_DIR/cover-back.jpg"
    echo "PLACED: cover-back.jpg (absolute path)"
  else
    echo "WARNING: Back image not found: $BACK_IMG"
  fi
fi

echo "--- Book folder ready: $BOOK_DIR"
echo "    Next: create $SLUG.md with frontmatter"
