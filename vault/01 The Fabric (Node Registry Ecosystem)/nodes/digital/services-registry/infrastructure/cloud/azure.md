---
title: "Azure"
type: reference
status: active
tags: [platform, cloud, openai, api]
---
# Azure

> OpenAI API endpoint (GPT-5.x, GPT-4.1, o4-mini).

## Endpoint
- **Base URL:** https://onemind-openai.openai.azure.com/openai/v1
- **Models:** gpt-5.x, gpt-4.1, o4-mini, gpt-4o
- **Note:** gpt-5.x and o4-mini require `max_completion_tokens` (not `max_tokens`)

## CLI
```bash
az login
az cognitiveservices account list
az cognitiveservices account deployment list --name onemind-openai -g <rg>
```

## Cost
~$30/mo (usage-based)
