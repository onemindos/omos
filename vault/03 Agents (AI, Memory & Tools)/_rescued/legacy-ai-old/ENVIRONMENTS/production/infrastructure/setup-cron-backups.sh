#!/bin/bash
# setup-cron-backups.sh - Setup automated database backups via cron
# Usage: ./setup-cron-backups.sh [--email user@example.com] [--s3-bucket bucket-name]

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Default configuration
BACKUP_TIME="02:00"  # 2 AM daily
NOTIFICATION_EMAIL=""
S3_BUCKET=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --email)
            NOTIFICATION_EMAIL="$2"
            shift 2
            ;;
        --s3-bucket)
            S3_BUCKET="$2"
            shift 2
            ;;
        --time)
            BACKUP_TIME="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [--email user@example.com] [--s3-bucket bucket-name] [--time HH:MM]"
            echo ""
            echo "Options:"
            echo "  --email EMAIL         Email address for backup notifications"
            echo "  --s3-bucket BUCKET    S3 bucket name for offsite backups"
            echo "  --time HH:MM         Backup time in 24-hour format (default: 02:00)"
            echo ""
            echo "Examples:"
            echo "  $0"
            echo "  $0 --email admin@zeusops.dev --s3-bucket legacy-ai-backups"
            echo "  $0 --time 03:30 --email alerts@zeusops.dev"
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🔧 Setting up automated database backups...${NC}"
echo ""

# Get script directory (absolute path)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_SCRIPT="$SCRIPT_DIR/backup-databases.sh"

# Verify backup script exists
if [ ! -f "$BACKUP_SCRIPT" ]; then
    echo -e "${RED}❌ Backup script not found: $BACKUP_SCRIPT${NC}"
    exit 1
fi

# Make backup script executable
chmod +x "$BACKUP_SCRIPT"
echo -e "${GREEN}✓${NC} Backup script executable: $BACKUP_SCRIPT"

# Create log directory
LOG_DIR="/var/log/legacy-ai-backups"
sudo mkdir -p "$LOG_DIR"
echo -e "${GREEN}✓${NC} Log directory created: $LOG_DIR"

# Build backup command
BACKUP_CMD="$BACKUP_SCRIPT"
if [ -n "$S3_BUCKET" ]; then
    BACKUP_CMD="$BACKUP_CMD --s3-bucket $S3_BUCKET"
    echo -e "${BLUE}→${NC} S3 backups enabled: $S3_BUCKET"
fi

# Create wrapper script for logging and notifications
WRAPPER_SCRIPT="/usr/local/bin/legacy-ai-backup-wrapper.sh"
echo -e "${BLUE}→${NC} Creating wrapper script..."

sudo tee "$WRAPPER_SCRIPT" > /dev/null <<EOF
#!/bin/bash
# Automated backup wrapper with logging and notifications

LOG_FILE="$LOG_DIR/backup-\$(date +%Y%m%d-%H%M%S).log"
BACKUP_CMD="$BACKUP_CMD"
EMAIL="$NOTIFICATION_EMAIL"

# Run backup with logging
echo "=== Backup started at \$(date) ===" > "\$LOG_FILE"
if \$BACKUP_CMD >> "\$LOG_FILE" 2>&1; then
    echo "=== Backup completed successfully at \$(date) ===" >> "\$LOG_FILE"
    EXIT_CODE=0
else
    echo "=== Backup failed at \$(date) ===" >> "\$LOG_FILE"
    EXIT_CODE=1
fi

# Send email notification on failure (if email configured)
if [ \$EXIT_CODE -ne 0 ] && [ -n "\$EMAIL" ]; then
    if command -v mail &> /dev/null; then
        cat "\$LOG_FILE" | mail -s "⚠️ Legacy AI Backup Failed" "\$EMAIL"
    elif command -v sendmail &> /dev/null; then
        echo "Subject: ⚠️ Legacy AI Backup Failed" | cat - "\$LOG_FILE" | sendmail "\$EMAIL"
    fi
fi

# Keep only last 30 days of logs
find "$LOG_DIR" -name "backup-*.log" -type f -mtime +30 -delete

exit \$EXIT_CODE
EOF

sudo chmod +x "$WRAPPER_SCRIPT"
echo -e "${GREEN}✓${NC} Wrapper script created: $WRAPPER_SCRIPT"

# Parse backup time
HOUR=$(echo "$BACKUP_TIME" | cut -d':' -f1)
MINUTE=$(echo "$BACKUP_TIME" | cut -d':' -f2)

# Create cron job
CRON_ENTRY="$MINUTE $HOUR * * * $WRAPPER_SCRIPT"
echo -e "${BLUE}→${NC} Installing cron job..."

# Remove existing legacy-ai-backup cron jobs
(crontab -l 2>/dev/null | grep -v "legacy-ai-backup-wrapper.sh" || true) | sudo crontab -

# Add new cron job
(sudo crontab -l 2>/dev/null || true; echo "$CRON_ENTRY") | sudo crontab -

echo -e "${GREEN}✓${NC} Cron job installed"

# Display configuration summary
echo ""
echo -e "${GREEN}✅ Automated backups configured successfully!${NC}"
echo ""
echo "📋 Configuration Summary:"
echo "  • Backup Time: Daily at $BACKUP_TIME"
echo "  • Backup Script: $BACKUP_SCRIPT"
echo "  • Wrapper Script: $WRAPPER_SCRIPT"
echo "  • Log Directory: $LOG_DIR"
if [ -n "$S3_BUCKET" ]; then
    echo "  • S3 Bucket: $S3_BUCKET (offsite backups enabled)"
fi
if [ -n "$NOTIFICATION_EMAIL" ]; then
    echo "  • Notifications: $NOTIFICATION_EMAIL (on failure)"
fi
echo ""
echo "📊 View cron jobs:"
echo "   sudo crontab -l"
echo ""
echo "📁 View logs:"
echo "   ls -lh $LOG_DIR"
echo "   tail -f $LOG_DIR/backup-*.log"
echo ""
echo "🧪 Test backup manually:"
echo "   $WRAPPER_SCRIPT"
echo ""
echo "💡 Tip: First backup will run tomorrow at $BACKUP_TIME"
echo "         Run the test command above to verify configuration"
echo ""
