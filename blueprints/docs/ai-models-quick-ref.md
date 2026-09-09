# 🤖 AI Models Quick Reference

## LM Studio (Local - Ready to Use)

### Current Status
✅ LM Studio app: Running  
✅ Hermes config: Configured  
⚠️ Server: Needs manual start in LM Studio UI  
⚠️ Models: Only Gemma 4B downloaded (general purpose)

### Start LM Studio Server
1. Open LM Studio app (click its icon in Dock)
2. Click **Developer** tab (left sidebar)
3. Click **Start Server** button
4. Verify: `curl http://localhost:1234/v1/models`

### Download Coding Models (Sonnet/Haiku Level)

| Model | Size | Speed | Best For |
|-------|------|-------|----------|
| **Qwen2.5-Coder-32B** | ~20GB | Medium | Code generation, debugging |
| **DeepSeek-Coder-V2-Lite** | ~16GB | Fast | Code completion, infill |
| **Mistral-Codestral-22B** | ~13GB | Fast | Fill-in-the-middle |
| **CodeLlama-70B** | ~40GB | Slow | Complex reasoning |

### Download Steps
1. In LM Studio: **Discover** tab → Search model name
2. Select **Q4_K_M** quantization (best speed/quality tradeoff)
3. Click Download

### Use in Hermes
```bash
# One-time
hermes --provider lmstudio --model lmstudio-community/gemma-4b-it

# Set default
hermes config set model.provider lmstudio
hermes config set model.default lmstudio-community/gemma-4b-it

# Check available models
curl -s http://localhost:1234/v1/models | jq '.data[].id'
```

---

## DigitalOcean / Paperspace (Cloud - Requires Setup)

### Status
⏳ Needs: Paperspace account + API key + deployed endpoint

### Quick Setup
1. Sign up: https://www.paperspace.com
2. Get API key from console
3. Add to `~/.hermes/.env`:
   ```
   PAPERSPACE_API_KEY=your_key_here
   ```
4. Deploy endpoint in Paperspace with vLLM
5. Update `~/.hermes/config.yaml` with endpoint URL

### Usage
```bash
hermes --provider paperspace --model llama-3.1-70b-instruct
```

---

## Current Hermes Providers

```yaml
# Default: AWS Bedrock (Claude Sonnet 4)
# Fallback: LM Studio (once server started)
# Optional: Paperspace (once configured)
```

```bash
# View config
hermes config get

# Switch provider
hermes config set model.provider lmstudio

# Test connection
hermes --provider lmstudio "Hello from local model!"
```
