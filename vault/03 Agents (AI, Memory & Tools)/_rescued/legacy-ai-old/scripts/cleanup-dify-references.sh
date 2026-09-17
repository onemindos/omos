#!/bin/bash
###############################################################################
# Remove Dify References from Local Repository
###############################################################################

set -e

REPO_DIR="/Users/zeuslegacy/Documents/Legacy-AI"

echo "╔════════════════════════════════════════════════════════════╗"
echo "║    Cleaning Dify References from Local Repository         ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

cd "$REPO_DIR"

# Create backup of files before modification
BACKUP_DIR="$HOME/repo-cleanup-backup-$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "📦 Creating backup of files to be modified..."
echo "   Backup location: $BACKUP_DIR"
echo ""

# Remove Dify directory
echo "🗑️  Removing Dify capability directory..."
if [ -d "$REPO_DIR/CAPABILITIES/digital-actions/workflow-automation/dify" ]; then
    mkdir -p "$BACKUP_DIR/CAPABILITIES/digital-actions/workflow-automation/"
    mv "$REPO_DIR/CAPABILITIES/digital-actions/workflow-automation/dify" "$BACKUP_DIR/CAPABILITIES/digital-actions/workflow-automation/"
    echo "   ✅ Moved to backup: CAPABILITIES/digital-actions/workflow-automation/dify"
fi

# Remove deploy-dify.sh
echo ""
echo "🗑️  Removing Dify deployment script..."
if [ -f "$REPO_DIR/ENVIRONMENTS/production/deploy-dify.sh" ]; then
    mkdir -p "$BACKUP_DIR/ENVIRONMENTS/production/"
    mv "$REPO_DIR/ENVIRONMENTS/production/deploy-dify.sh" "$BACKUP_DIR/ENVIRONMENTS/production/"
    echo "   ✅ Moved to backup: ENVIRONMENTS/production/deploy-dify.sh"
fi

# Remove symlink if it exists
if [ -L "$REPO_DIR/scripts/deploy-dify.sh" ]; then
    rm "$REPO_DIR/scripts/deploy-dify.sh"
    echo "   ✅ Removed symlink: scripts/deploy-dify.sh"
fi

# Remove Dify from Caddyfile.optimized
echo ""
echo "📝 Updating Caddyfile.optimized..."
if [ -f "$REPO_DIR/ENVIRONMENTS/production/infrastructure/caddy/Caddyfile.optimized" ]; then
    cp "$REPO_DIR/ENVIRONMENTS/production/infrastructure/caddy/Caddyfile.optimized" "$BACKUP_DIR/"
    
    # Create temp file without Dify block
    awk '
        /^# Dify AI Platform/ { skip=1 }
        /^dify\.zeusops\.dev {/ { skip=1 }
        skip && /^}/ { skip=0; next }
        !skip { print }
    ' "$REPO_DIR/ENVIRONMENTS/production/infrastructure/caddy/Caddyfile.optimized" > /tmp/caddyfile.tmp
    
    mv /tmp/caddyfile.tmp "$REPO_DIR/ENVIRONMENTS/production/infrastructure/caddy/Caddyfile.optimized"
    echo "   ✅ Removed Dify domain block from Caddyfile.optimized"
fi

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║         LOCAL REPOSITORY CLEANUP COMPLETE                  ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "✅ Dify directory removed from CAPABILITIES/"
echo "✅ deploy-dify.sh removed from ENVIRONMENTS/production/"
echo "✅ dify.zeusops.dev removed from Caddyfile.optimized"
echo ""
echo "📦 All removed files backed up to:"
echo "   $BACKUP_DIR"
echo ""
echo "⚠️  NOTE: Some documentation files still contain Dify references."
echo "   These should be manually reviewed and updated as needed:"
echo "   • README.md"
echo "   • VM-UPGRADE-*.md files"
echo "   • .github/workflows/*.yml files"
echo "   • KNOWLEDGE/docs/**/*.md files"
echo ""
echo "💡 To find remaining references, run:"
echo "   grep -r 'dify' --exclude-dir=.git --exclude-dir=node_modules ."
echo ""
echo "📝 Summary saved to: $BACKUP_DIR/CLEANUP_SUMMARY.txt"

# Save summary
cat > "$BACKUP_DIR/CLEANUP_SUMMARY.txt" << EOF
Local Repository Dify Cleanup
Date: $(date)

Removed:
1. CAPABILITIES/digital-actions/workflow-automation/dify/ (entire directory)
2. ENVIRONMENTS/production/deploy-dify.sh
3. scripts/deploy-dify.sh (symlink)
4. Dify domain block from Caddyfile.optimized

Backup Location: $BACKUP_DIR

Next Steps:
- Review and update documentation files that mention Dify
- Update GitHub workflows to remove Dify-related jobs
- Consider updating README.md to reflect Dify removal
EOF

echo ""
