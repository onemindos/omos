#!/bin/bash
# Setup Caddy with command.zeusops.dev for Portainer
# Run this script on the VM via Lightsail browser SSH

set -e

echo "🚀 Setting up Caddy for Legacy AI Infrastructure..."
echo ""

# Create directory structure
echo "📁 Creating directories..."
sudo mkdir -p /opt/legacy-ai/infrastructure/caddy
cd /opt/legacy-ai/infrastructure/caddy

# Create Caddyfile
echo "📝 Creating Caddyfile..."
sudo tee Caddyfile > /dev/null <<'EOF'
{
    # Enable automatic HTTPS
    email admin@zeusops.dev
}

# EMQX MQTT Broker Dashboard
emqx.zeusops.dev {
    reverse_proxy emqx:18083
}

# Home Assistant
ha.zeusops.dev {
    reverse_proxy home-assistant:8123
}

# Node-RED
nodered.zeusops.dev {
    reverse_proxy node-red:1880
}

# n8n Workflow Automation
n8n.zeusops.dev {
    reverse_proxy n8n-main:5678
}

# Dify AI Platform
dify.zeusops.dev {
    reverse_proxy dify-nginx:80
}

# Letta AI Agent
letta.zeusops.dev {
    reverse_proxy letta:8283
}

# Portainer - Command Center (primary)
command.zeusops.dev {
    reverse_proxy https://portainer:9443 {
        transport http {
            tls_insecure_skip_verify
        }
    }
}

# Portainer - Dashboard (alternative)
dashboard.zeusops.dev {
    reverse_proxy portainer:9000
}
EOF

echo "✅ Caddyfile created"

# Create docker-compose.yml
echo "📝 Creating docker-compose.yml..."
sudo tee docker-compose.yml > /dev/null <<'EOF'
version: '3.8'

services:
  caddy:
    image: caddy:latest
    container_name: caddy
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy_data:/data
      - caddy_config:/config
    networks:
      - legacy_network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:80"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 5s

volumes:
  caddy_data:
  caddy_config:

networks:
  legacy_network:
    external: true
EOF

echo "✅ docker-compose.yml created"

# Create network if it doesn't exist
echo "🌐 Setting up Docker network..."
docker network create legacy_network 2>/dev/null && echo "Created legacy_network" || echo "legacy_network already exists"

# Connect Portainer to network
echo "🔗 Connecting Portainer to network..."
docker network connect legacy_network portainer 2>/dev/null && echo "Connected Portainer" || echo "Portainer already connected"

# Deploy Caddy
echo "🚀 Deploying Caddy..."
sudo docker-compose up -d

# Wait a moment for startup
echo "⏳ Waiting for Caddy to start..."
sleep 3

# Check status
echo ""
echo "📊 Container Status:"
docker ps | grep -E "caddy|portainer"

echo ""
echo "📋 Caddy Logs:"
docker logs caddy --tail 20

echo ""
echo "✅ Setup complete!"
echo ""
echo "📍 Your Tailscale IP: $(tailscale ip -4)"
echo ""
echo "🎯 Configured Domains:"
echo "  • https://emqx.zeusops.dev      → EMQX MQTT Dashboard"
echo "  • https://ha.zeusops.dev        → Home Assistant"
echo "  • https://nodered.zeusops.dev   → Node-RED"
echo "  • https://n8n.zeusops.dev       → n8n Workflow Automation"
echo "  • https://dify.zeusops.dev      → Dify AI Platform"
echo "  • https://letta.zeusops.dev     → Letta AI Agent"
echo "  • https://command.zeusops.dev   → Portainer (primary)"
echo "  • https://dashboard.zeusops.dev → Portainer (alternative)"
echo ""
echo "🔐 Next steps:"
echo "1. Ensure DNS A records point to VM public IP: 3.215.29.128"
echo "2. Wait for Let's Encrypt SSL certificates (automatic)"
echo "3. Test each domain in browser"
echo ""
echo "📊 Check SSL status:"
echo "   docker logs caddy | grep certificate"
echo ""
