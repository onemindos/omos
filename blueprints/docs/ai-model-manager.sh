#!/bin/bash
# AI Model Manager for Hermes
# Manages LM Studio and Paperspace integrations

set -e

HERMES_CONFIG="$HOME/.hermes/config.yaml"
HERMES_ENV="$HOME/.hermes/.env"

show_help() {
    cat << EOF
AI Model Manager for Hermes

Usage: $0 [command]

Commands:
    status          Show status of all AI providers
    lmstudio-start  Configure and start LM Studio server
    lmstudio-test   Test LM Studio connection
    lmstudio-use    Set LM Studio as default provider
    paperspace-setup Setup Paperspace integration
    bedrock-use     Set AWS Bedrock as default provider
    list-models     List available models in LM Studio

EOF
}

show_status() {
    echo "=== 🤖 AI Provider Status ==="
    echo ""

    # Check LM Studio
    echo "📦 LM Studio:"
    if pgrep -q "LM Studio"; then
        echo "  ✅ App running"
        if curl -s http://localhost:1234/v1/models > /dev/null 2>&1; then
            echo "  ✅ Server responding on :1234"
            MODELS=$(curl -s http://localhost:1234/v1/models 2>/dev/null | grep -o '"id": "[^"]*"' | wc -l)
            echo "  📊 Models loaded: $MODELS"
        else
            echo "  ⚠️  Server not responding (start in LM Studio UI → Developer → Start Server)"
        fi
    else
        echo "  ❌ App not running"
    fi
    echo ""

    # Check current Hermes provider
    echo "🎯 Hermes Configuration:"
    CURRENT_PROVIDER=$(grep "provider:" "$HERMES_CONFIG" 2>/dev/null | head -1 | cut -d: -f2 | tr -d ' "' || echo "unknown")
    CURRENT_MODEL=$(grep "default:" "$HERMES_CONFIG" 2>/dev/null | head -1 | cut -d: -f2 | tr -d ' "' || echo "unknown")
    echo "  Provider: $CURRENT_PROVIDER"
    echo "  Model: $CURRENT_MODEL"
    echo ""

    # Check Paperspace
    echo "☁️  Paperspace (DigitalOcean):"
    if grep -q "PAPERSPACE_API_KEY" "$HERMES_ENV" 2>/dev/null; then
        echo "  ✅ API key configured"
    else
        echo "  ❌ Not configured (needs PAPERSPACE_API_KEY in ~/.hermes/.env)"
    fi
}

lmstudio_start() {
    echo "=== 🔧 Configuring LM Studio ==="

    # Ensure settings file has server enabled
    if [ -f "$HOME/.lmstudio/settings.json" ]; then
        python3 << PYEOF
import json
with open('$HOME/.lmstudio/settings.json', 'r') as f:
    settings = json.load(f)

settings['apiServer'] = {
    "enabled": True,
    "port": 1234,
    "cors": True,
    "title": "LM Studio Server"
}

with open('$HOME/.lmstudio/settings.json', 'w') as f:
    json.dump(settings, f, indent=2)
print("✅ Server settings updated")
PYEOF
    fi

    # Launch LM Studio if not running
    if ! pgrep -q "LM Studio"; then
        echo "🚀 Launching LM Studio..."
        open -a "LM Studio"
        sleep 3
    fi

    echo ""
    echo "⚠️  IMPORTANT: Manual step required!"
    echo ""
    echo "1. LM Studio app is now open"
    echo "2. Click 'Developer' tab in left sidebar"
    echo "3. Click 'Start Server' button"
    echo "4. Wait for server to start"
    echo ""
    echo "Then run: $0 lmstudio-test"
}

lmstudio_test() {
    echo "=== 🧪 Testing LM Studio Connection ==="

    if ! curl -s http://localhost:1234/v1/models > /dev/null 2>&1; then
        echo "❌ LM Studio server not responding on :1234"
        echo "   Make sure to start the server in LM Studio UI"
        exit 1
    fi

    echo "✅ Server responding!"
    echo ""
    echo "Available models:"
    curl -s http://localhost:1234/v1/models 2>/dev/null | python3 -c "
import json, sys
try:
    data = json.load(sys.stdin)
    for m in data.get('data', []):
        print(f'  - {m[\"id\"]}')
except:
    pass
" || echo "  (could not parse response)"
}

lmstudio_use() {
    echo "=== 🔄 Setting LM Studio as Default ==="

    hermes config set model.provider lmstudio

    # Get first available model
    MODEL=$(curl -s http://localhost:1234/v1/models 2>/dev/null | python3 -c "
import json, sys
try:
    data = json.load(sys.stdin)
    if data.get('data'):
        print(data['data'][0]['id'])
except:
    pass
" || echo "lmstudio-community/gemma-4b-it")

    hermes config set model.default "$MODEL"
    echo "✅ LM Studio set as default provider"
    echo "   Model: $MODEL"
}

paperspace_setup() {
    echo "=== ☁️  Paperspace Setup ==="
    echo ""
    echo "DigitalOcean AI is through their Paperspace acquisition."
    echo ""
    echo "Steps:"
    echo "1. Sign up at https://www.paperspace.com"
    echo "2. Get API key from console"
    echo "3. echo 'PAPERSPACE_API_KEY=your_key' >> ~/.hermes/.env"
    echo "4. Deploy a model endpoint in Paperspace"
    echo "5. Update base_url in ~/.hermes/config.yaml"
    echo ""
    echo "Recommended models: Llama-3.1-70B, Mixtral-8x7B, Qwen2.5-72B"
}

bedrock_use() {
    echo "=== 🔄 Setting AWS Bedrock as Default ==="
    hermes config set model.provider bedrock
    hermes config set model.default us.anthropic.claude-sonnet-4-6
    echo "✅ AWS Bedrock (Claude Sonnet 4) set as default"
}

list_models() {
    echo "=== 📋 LM Studio Models ==="
    if curl -s http://localhost:1234/v1/models > /dev/null 2>&1; then
        curl -s http://localhost:1234/v1/models | python3 -m json.tool 2>/dev/null | grep '"id"' || echo "No models loaded"
    else
        echo "⚠️  Server not running. Start with: $0 lmstudio-start"
    fi

    echo ""
    echo "=== 📋 Recommended Coding Models ==="
    cat << 'MODELEOF'
Download these in LM Studio → Discover tab:

1. Qwen2.5-Coder-32B-Instruct (⭐ Best, ~20GB)
   - Outperforms GPT-3.5 on coding
   - Download: Qwen/Qwen2.5-Coder-32B-Instruct-GGUF

2. DeepSeek-Coder-V2-Lite (⭐, ~16GB)
   - 87% code training data
   - Download: deepseek-ai/DeepSeek-Coder-V2-Lite

3. Mistral-Codestral-22B (⭐, ~13GB)
   - Fill-in-the-middle specialist
   - Download: mistralai/Codestral-22B-v0.1

4. CodeLlama-70B-Instruct (~40GB)
   - Meta's coding model
   - Download: codellama/CodeLlama-70b-Instruct

Use quantization: Q4_K_M for best speed/quality
MODELEOF
}

# Main
case "${1:-status}" in
    status) show_status ;;
    lmstudio-start) lmstudio_start ;;
    lmstudio-test) lmstudio_test ;;
    lmstudio-use) lmstudio_use ;;
    paperspace-setup) paperspace_setup ;;
    bedrock-use) bedrock_use ;;
    list-models) list_models ;;
    help|--help|-h) show_help ;;
    *) echo "Unknown command: $1"; show_help; exit 1 ;;
esac
