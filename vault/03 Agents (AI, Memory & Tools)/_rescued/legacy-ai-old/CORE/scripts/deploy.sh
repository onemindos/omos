#!/bin/bash
# Phase 2 Deployment Script
# Deploy Legacy AI hybrid cognitive architecture

set -e  # Exit on error

echo "🚀 Legacy AI - Phase 2 Deployment"
echo "=================================="

# ============================================================
# 1. Environment Validation
# ============================================================

echo ""
echo "📋 Step 1: Validating environment..."

if [ ! -f .env ]; then
    echo "❌ .env file not found!"
    echo "   Please copy .env.example to .env and fill in your API keys."
    exit 1
fi

# Check required environment variables
required_vars=(
    "LIVEKIT_API_KEY"
    "LIVEKIT_API_SECRET"
    "ASSEMBLYAI_API_KEY"
    "CARTESIA_API_KEY"
    "OPENAI_API_KEY"
    "SUPABASE_URL"
    "SUPABASE_SERVICE_KEY"
)

missing_vars=()
for var in "${required_vars[@]}"; do
    if ! grep -q "^${var}=" .env || grep -q "^${var}=your_" .env; then
        missing_vars+=("$var")
    fi
done

if [ ${#missing_vars[@]} -gt 0 ]; then
    echo "❌ Missing or incomplete environment variables:"
    printf '   - %s\n' "${missing_vars[@]}"
    echo ""
    echo "   Please update your .env file with valid API keys."
    exit 1
fi

echo "✅ Environment validation passed"

# ============================================================
# 2. Generate LiveKit Keys (if not set)
# ============================================================

echo ""
echo "🔑 Step 2: Checking LiveKit keys..."

if grep -q "^LIVEKIT_API_KEY=your_" .env; then
    echo "⚙️  Generating LiveKit API keys..."
    API_KEY=$(openssl rand -hex 32)
    API_SECRET=$(openssl rand -hex 32)
    
    # Update .env file
    sed -i.bak "s/^LIVEKIT_API_KEY=.*/LIVEKIT_API_KEY=${API_KEY}/" .env
    sed -i.bak "s/^LIVEKIT_API_SECRET=.*/LIVEKIT_API_SECRET=${API_SECRET}/" .env
    rm .env.bak
    
    echo "✅ Generated new LiveKit keys"
else
    echo "✅ LiveKit keys already configured"
fi

# ============================================================
# 3. Initialize Supabase Schema
# ============================================================

echo ""
echo "🗄️  Step 3: Initializing Supabase schema..."

SUPABASE_URL=$(grep "^SUPABASE_URL=" .env | cut -d '=' -f2)
SUPABASE_KEY=$(grep "^SUPABASE_SERVICE_KEY=" .env | cut -d '=' -f2)

if [ -z "$SUPABASE_URL" ] || [ -z "$SUPABASE_KEY" ]; then
    echo "⚠️  Skipping Supabase init (credentials not set)"
else
    echo "   Applying schema to Supabase..."
    # TODO: Use Supabase CLI or psql to apply schema
    # For now, manual step required
    echo "⚠️  Manual step: Apply supabase/schema.sql in Supabase dashboard"
fi

# ============================================================
# 4. Pull Docker Images
# ============================================================

echo ""
echo "🐳 Step 4: Pulling Docker images..."

docker-compose pull

echo "✅ Images pulled"

# ============================================================
# 5. Build Custom Images
# ============================================================

echo ""
echo "🔨 Step 5: Building custom images..."

docker-compose build

echo "✅ Build complete"

# ============================================================
# 6. Start Services
# ============================================================

echo ""
echo "▶️  Step 6: Starting services..."

docker-compose up -d

echo "✅ Services started"

# ============================================================
# 7. Wait for Services
# ============================================================

echo ""
echo "⏳ Step 7: Waiting for services to be ready..."

# Wait for LiveKit
echo "   Checking LiveKit..."
for i in {1..30}; do
    if curl -s http://localhost:7880 > /dev/null 2>&1; then
        echo "   ✅ LiveKit ready"
        break
    fi
    sleep 1
done

# Wait for Letta
echo "   Checking Letta..."
for i in {1..30}; do
    if curl -s http://localhost:8283/health > /dev/null 2>&1; then
        echo "   ✅ Letta ready"
        break
    fi
    sleep 1
done

# ============================================================
# 8. Create Letta Agent
# ============================================================

echo ""
echo "🧠 Step 8: Creating Letta agent..."

# TODO: Use Letta API to create agent from letta/agent-config.json
# For now, manual step required
echo "⚠️  Manual step: Create Letta agent via UI or API"
echo "   1. Visit http://localhost:8283"
echo "   2. Create new agent using letta/agent-config.json"
echo "   3. Copy agent ID to .env (LETTA_AGENT_ID)"

# ============================================================
# 9. Display Status
# ============================================================

echo ""
echo "=============================================="
echo "✅ Deployment Complete!"
echo "=============================================="
echo ""
echo "Services:"
echo "  - LiveKit:  http://localhost:7880"
echo "  - Letta:    http://localhost:8283"
echo "  - Agent:    Running in background"
echo ""
echo "Next Steps:"
echo "  1. Complete Letta agent setup (see above)"
echo "  2. Apply Supabase schema (supabase/schema.sql)"
echo "  3. Test with: ./scripts/test-latency.sh"
echo ""
echo "View logs:"
echo "  docker-compose logs -f agent"
echo ""
echo "Stop services:"
echo "  docker-compose down"
echo ""
