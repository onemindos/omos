---
id: so-01-ai-models-aws
title: "AI Models — AWS (Amazon Bedrock)"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/aws]
ai_summary: "AWS Amazon Bedrock model reference — managed foundation model API (Claude, Llama, Titan, Nova, Mistral). In stack for enterprise workloads, AWS-integrated pipelines, and S3/Lambda-connected inference."
---

# AWS (Amazon Bedrock)

> STACK provider. Amazon Bedrock gives managed API access to multiple foundation models including Claude (Anthropic), Llama (Meta), Nova (Amazon), and Mistral — all through one AWS API with IAM auth.
> Key advantage: integrates natively with S3, Lambda, API Gateway — no separate API key management when already in AWS ecosystem.

---

## Overview

Amazon Bedrock is AWS's managed foundation model service. Rather than being a model creator, AWS provides a unified API to access models from Anthropic, Meta, Mistral, Amazon, and others — all billed through AWS and secured via IAM roles/policies instead of per-provider API keys.

- **API Base URL:** `https://bedrock-runtime.{region}.amazonaws.com`
- **SDK:** `@aws-sdk/client-bedrock-runtime` (Node.js) | `boto3` + `botocore` (Python)
- **Auth:** IAM credentials (AWS_ACCESS_KEY_ID + AWS_SECRET_ACCESS_KEY) or STS role assumption
- **Key Location:** 1Password → SO vault → `AWS Bedrock Credentials`
- **Console:** https://console.aws.amazon.com/bedrock
- **Pricing:** https://aws.amazon.com/bedrock/pricing

---

## Models Available on Bedrock

### Anthropic Claude (via Bedrock)
| Model ID | Context | Notes |
|----------|---------|-------|
| `anthropic.claude-3-7-sonnet-20250219-v1:0` | 200K | Same as direct Anthropic API |
| `anthropic.claude-3-5-haiku-20241022-v1:0` | 200K | Fast/cheap |
| `anthropic.claude-3-5-sonnet-20241022-v2:0` | 200K | Stable |

### Amazon Nova
| Model ID | Context | Input ($/1M) | Output ($/1M) | Notes |
|----------|---------|-------------|--------------|-------|
| `amazon.nova-pro-v1:0` | 300K | $0.80 | $3.20 | Amazon's flagship |
| `amazon.nova-lite-v1:0` | 300K | $0.06 | $0.24 | Fast, low cost |
| `amazon.nova-micro-v1:0` | 128K | $0.035 | $0.14 | Cheapest text-only |

### Meta Llama (via Bedrock)
| Model ID | Context | Notes |
|----------|---------|-------|
| `meta.llama3-3-70b-instruct-v1:0` | 128K | Latest Llama 3.3 70B |
| `meta.llama3-1-8b-instruct-v1:0` | 128K | Lightweight |

### Mistral (via Bedrock)
| Model ID | Context | Notes |
|----------|---------|-------|
| `mistral.mistral-large-2402-v1:0` | 32K | Mistral Large |
| `mistral.mixtral-8x7b-instruct-v0:1` | 32K | MoE model |

### Amazon Titan (embeddings)
| Model ID | Dimensions | Notes |
|----------|-----------|-------|
| `amazon.titan-embed-text-v2:0` | 1024 | AWS-native embeddings for S3-connected pipelines |

---

## API Setup

```bash
npm install @aws-sdk/client-bedrock-runtime    # Node.js
pip install boto3                               # Python

export AWS_ACCESS_KEY_ID="AKIA..."
export AWS_SECRET_ACCESS_KEY="..."
export AWS_REGION="us-east-1"
```

**Node.js (Converse API — recommended):**
```typescript
import { BedrockRuntimeClient, ConverseCommand } from '@aws-sdk/client-bedrock-runtime';

const client = new BedrockRuntimeClient({ region: 'us-east-1' });

const response = await client.send(new ConverseCommand({
  modelId: 'anthropic.claude-3-5-haiku-20241022-v1:0',
  messages: [{ role: 'user', content: [{ text: 'Hello' }] }],
}));

const text = response.output?.message?.content?.[0]?.text;
```

**Python (boto3):**
```python
import boto3, json

client = boto3.client('bedrock-runtime', region_name='us-east-1')
response = client.converse(
    modelId='anthropic.claude-3-5-haiku-20241022-v1:0',
    messages=[{'role': 'user', 'content': [{'text': 'Hello'}]}]
)
print(response['output']['message']['content'][0]['text'])
```

**OpenAI-compatible endpoint (Bedrock):**
```typescript
// Bedrock supports OpenAI-compatible format for select models
import OpenAI from 'openai';
const client = new OpenAI({
  baseURL: 'https://bedrock-runtime.us-east-1.amazonaws.com/openai/v1',
  // Auth via AWS Signature V4 — use aws4 or @aws-sdk/signature-v4
});
```

---

## IAM Policy (Minimum Required)

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "bedrock:InvokeModel",
      "bedrock:InvokeModelWithResponseStream",
      "bedrock:Converse",
      "bedrock:ConverseStream"
    ],
    "Resource": "arn:aws:bedrock:us-east-1::foundation-model/*"
  }]
}
```

---

## Usage in One Mind Stack

| Use Case | Model | Notes |
|----------|-------|-------|
| **AWS-integrated pipelines** | claude-3-5-haiku (Bedrock) | When S3/Lambda trigger is the source |
| **Enterprise/compliance** | Any Bedrock model | Data stays within AWS VPC |
| **Native AWS embeddings** | amazon.titan-embed-text-v2 | For Bedrock Knowledge Bases |
| **Bedrock Agents** | claude-3-7-sonnet | AWS-native agent with S3 knowledge base |
| **Cost optimization** | amazon.nova-lite | Cheapest managed option per token |

---

## Notes

- **Model access:** Must enable each model in Bedrock console before first use (Model Access page)
- **Regions:** Not all models available in all regions — us-east-1 has widest selection
- **Pricing:** Same as direct provider for Claude; Amazon Nova is cheaper than most alternatives
- **Bedrock Agents:** Full agent system with S3 knowledge bases, action groups (Lambda), guardrails — separate from OpenClaw
- **vs Direct Anthropic:** Bedrock Claude costs the same per token but routes billing through AWS; useful when existing AWS credits or consolidated billing is needed
- **Latency:** Similar to direct API calls; no additional penalty for Bedrock routing
- **VPC support:** Can run Bedrock within a VPC endpoint — zero internet exposure for sensitive inference


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
