#!/bin/bash
# ============================================
# LEGACY AI - DEPLOYMENT SCRIPT
# ============================================
# Deploys all Legacy AI services to your VM
# Date: October 22, 2025

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
VM_HOST="${VM_HOST:-3.215.29.128}"
VM_USER="${VM_USER:-ubuntu}"
VM_SSH_KEY="${VM_SSH_KEY:-$HOME/.ssh/legacyai.pem}"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   LEGACY AI - DEPLOYMENT SCRIPT${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Function to run command on VM
run_on_vm() {
    ssh -i "$VM_SSH_KEY" "$VM_USER@$VM_HOST" "$@"
}

# Function to copy files to VM
copy_to_vm() {
    scp -i "$VM_SSH_KEY" -r "$1" "$VM_USER@$VM_HOST:$2"
}

# Step 1: Create directory structure
echo -e "${YELLOW}📁 Creating directory structure on VM...${NC}"
run_on_vm "sudo mkdir -p /opt/legacy-ai/{stacks,data,logs,config}"
run_on_vm "sudo chown -R $VM_USER:$VM_USER /opt/legacy-ai"

# Step 2: Create Docker network
echo -e "${YELLOW}🌐 Creating Docker network...${NC}"
run_on_vm "sudo docker network create legacy_network || true"

# Step 3: Copy stack files
echo -e "${YELLOW}📦 Copying stack files to VM...${NC}"
copy_to_vm "./stacks" "/opt/legacy-ai/"

# Step 4: Copy environment file
echo -e "${YELLOW}🔐 Setting up environment variables...${NC}"
if [ -f ".env" ]; then
    copy_to_vm ".env" "/opt/legacy-ai/.env"
else
    echo -e "${RED}❌ .env file not found! Please create it from .env.example${NC}"
    exit 1
fi

# Step 5: Deploy Zep Memory
echo -e "${YELLOW}🧠 Deploying Zep Memory...${NC}"
run_on_vm "cd /opt/legacy-ai/stacks/memory && sudo docker-compose up -d"
sleep 5

# Step 6: Deploy Ollama (Local LLM)
echo -e "${YELLOW}🤖 Deploying Ollama...${NC}"
run_on_vm "cd /opt/legacy-ai/stacks/llm && sudo docker-compose up -d"
sleep 10

# Step 7: Deploy Voice Stack (Whisper + Piper)
echo -e "${YELLOW}🎤 Deploying Voice Stack...${NC}"
run_on_vm "cd /opt/legacy-ai/stacks/voice && sudo docker-compose up -d"
sleep 5

# Step 8: Deploy Vision Stack (Frigate + CompreFace)
echo -e "${YELLOW}📹 Deploying Vision Stack...${NC}"
run_on_vm "cd /opt/legacy-ai/stacks/vision && sudo docker-compose up -d"
sleep 5

# Step 9: Deploy Legacy AI Agent
echo -e "${YELLOW}🚀 Deploying Legacy AI Agent...${NC}"
run_on_vm "cd /opt/legacy-ai/stacks/legacy-ai && sudo docker-compose up -d"

# Step 10: Check status
echo -e "${GREEN}✅ Deployment complete! Checking status...${NC}"
echo ""
run_on_vm "sudo docker ps --format 'table {{.Names}}\t{{.Status}}' | grep -E 'legacy|zep|ollama|whisper|piper|frigate|compreface'"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   DEPLOYMENT SUMMARY${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "🧠 Zep Memory: http://$VM_HOST:8000"
echo -e "🤖 Ollama: http://$VM_HOST:11434"
echo -e "🎤 Whisper STT: http://$VM_HOST:10300"
echo -e "🔊 Piper TTS: http://$VM_HOST:10200"
echo -e "📹 Frigate NVR: http://$VM_HOST:5000"
echo -e "👤 CompreFace: http://$VM_HOST:8000"
echo -e "🚀 Legacy AI Agent: Check logs with: ssh -i \"$VM_SSH_KEY\" $VM_USER@$VM_HOST 'sudo docker logs -f legacy-ai-agent'"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo -e "1. Configure your Supabase and LiveKit Cloud credentials in .env"
echo -e "2. Add your Azure/AWS API keys for LLM access"
echo -e "3. Test the voice pipeline: ssh into VM and check container logs"
echo -e "4. Configure Frigate cameras in /opt/legacy-ai/stacks/vision/frigate/config/config.yml"
echo ""
echo -e "${GREEN}🎉 Legacy AI is now running!${NC}"
