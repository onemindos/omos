#!/bin/bash
# backup-databases.sh - Automated database backup script
# Backs up PostgreSQL and Redis databases to local and optionally S3
# Usage: ./backup-databases.sh [--s3-bucket bucket-name]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BACKUP_DIR="/opt/legacy-ai/backups"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=7
S3_BUCKET=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --s3-bucket)
            S3_BUCKET="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [--s3-bucket bucket-name]"
            echo ""
            echo "Options:"
            echo "  --s3-bucket    Upload backups to S3 bucket"
            echo "  --help         Show this help message"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🔒 Legacy AI Database Backup${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"
cd "$BACKUP_DIR"

# Function to check if container is running
check_container() {
    if ! docker ps --format '{{.Names}}' | grep -q "^$1$"; then
        echo -e "${RED}❌ Container $1 is not running${NC}"
        return 1
    fi
    return 0
}

# Function to backup PostgreSQL database
backup_postgres() {
    local container=$1
    local db_name=$2
    local db_user=$3
    local backup_file="${container}_${DATE}.dump"
    
    echo -e "${YELLOW}→ Backing up $container (PostgreSQL)...${NC}"
    
    if ! check_container "$container"; then
        return 1
    fi
    
    # Create backup
    if docker exec "$container" pg_dump -U "$db_user" -d "$db_name" -F c -f "/tmp/$backup_file"; then
        # Copy from container to host
        docker cp "$container:/tmp/$backup_file" "$BACKUP_DIR/$backup_file"
        
        # Remove from container
        docker exec "$container" rm "/tmp/$backup_file"
        
        # Compress backup
        gzip "$backup_file"
        
        local size=$(du -h "${backup_file}.gz" | cut -f1)
        echo -e "${GREEN}✅ $container backed up: ${backup_file}.gz ($size)${NC}"
        return 0
    else
        echo -e "${RED}❌ Failed to backup $container${NC}"
        return 1
    fi
}

# Function to backup Redis
backup_redis() {
    local container=$1
    local password=$2
    local backup_file="${container}_${DATE}.rdb"
    
    echo -e "${YELLOW}→ Backing up $container (Redis)...${NC}"
    
    if ! check_container "$container"; then
        return 1
    fi
    
    # Trigger BGSAVE
    if [ -z "$password" ]; then
        docker exec "$container" redis-cli BGSAVE > /dev/null
    else
        docker exec "$container" redis-cli -a "$password" BGSAVE > /dev/null 2>&1
    fi
    
    # Wait for save to complete (check every second for up to 30 seconds)
    for i in {1..30}; do
        if [ -z "$password" ]; then
            status=$(docker exec "$container" redis-cli LASTSAVE)
        else
            status=$(docker exec "$container" redis-cli -a "$password" LASTSAVE 2>/dev/null)
        fi
        
        if [ ! -z "$status" ]; then
            sleep 1
            break
        fi
        sleep 1
    done
    
    # Copy RDB file
    if docker cp "$container:/data/dump.rdb" "$BACKUP_DIR/$backup_file"; then
        # Compress backup
        gzip "$backup_file"
        
        local size=$(du -h "${backup_file}.gz" | cut -f1)
        echo -e "${GREEN}✅ $container backed up: ${backup_file}.gz ($size)${NC}"
        return 0
    else
        echo -e "${RED}❌ Failed to backup $container${NC}"
        return 1
    fi
}

# Backup PostgreSQL databases
echo -e "${BLUE}📊 PostgreSQL Backups${NC}"
echo "---"

# n8n PostgreSQL
backup_postgres "n8n-postgres" "n8n" "n8n"

# Dify PostgreSQL
backup_postgres "dify-postgres" "dify" "dify"

echo ""

# Backup Redis instances
echo -e "${BLUE}💾 Redis Backups${NC}"
echo "---"

# n8n Redis (no password)
backup_redis "n8n-redis" ""

# Dify Redis (with password)
backup_redis "dify-redis" "2MehsKm00ldzVYo6uG06pkn4J8hr5Mnw"

echo ""

# Clean up old backups (older than RETENTION_DAYS)
echo -e "${BLUE}🧹 Cleanup Old Backups${NC}"
echo "---"
echo -e "${YELLOW}→ Removing backups older than $RETENTION_DAYS days...${NC}"

deleted_count=0
for file in "$BACKUP_DIR"/*.gz; do
    if [ -f "$file" ] && [ $(find "$file" -mtime +$RETENTION_DAYS 2>/dev/null | wc -l) -gt 0 ]; then
        rm "$file"
        deleted_count=$((deleted_count + 1))
    fi
done

if [ $deleted_count -gt 0 ]; then
    echo -e "${GREEN}✅ Removed $deleted_count old backup(s)${NC}"
else
    echo -e "${GREEN}✅ No old backups to remove${NC}"
fi

echo ""

# List current backups
echo -e "${BLUE}📁 Current Backups${NC}"
echo "---"
ls -lh "$BACKUP_DIR"/*.gz 2>/dev/null | awk '{print $9, "-", $5}' || echo "No backups found"

echo ""

# Calculate total size
total_size=$(du -sh "$BACKUP_DIR" 2>/dev/null | cut -f1)
backup_count=$(ls -1 "$BACKUP_DIR"/*.gz 2>/dev/null | wc -l)
echo -e "${BLUE}Total: $backup_count backup(s), $total_size${NC}"

# Upload to S3 if bucket specified
if [ ! -z "$S3_BUCKET" ]; then
    echo ""
    echo -e "${BLUE}☁️  S3 Upload${NC}"
    echo "---"
    
    # Check if AWS CLI is installed
    if ! command -v aws &> /dev/null; then
        echo -e "${YELLOW}⚠️  AWS CLI not installed, skipping S3 upload${NC}"
    else
        echo -e "${YELLOW}→ Uploading to s3://$S3_BUCKET/backups/${DATE}/...${NC}"
        
        # Upload today's backups
        for file in "$BACKUP_DIR"/*_${DATE}*.gz; do
            if [ -f "$file" ]; then
                filename=$(basename "$file")
                if aws s3 cp "$file" "s3://$S3_BUCKET/backups/${DATE}/$filename"; then
                    echo -e "${GREEN}✅ Uploaded: $filename${NC}"
                else
                    echo -e "${RED}❌ Failed to upload: $filename${NC}"
                fi
            fi
        done
        
        echo -e "${GREEN}✅ S3 upload complete${NC}"
    fi
fi

echo ""
echo -e "${GREEN}✅ Backup complete!${NC}"
echo ""
echo "Backup location: $BACKUP_DIR"
echo "Retention: $RETENTION_DAYS days"

# Exit with success
exit 0
